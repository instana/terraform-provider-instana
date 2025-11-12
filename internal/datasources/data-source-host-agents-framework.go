package datasources

import (
	"context"
	"fmt"
	"time"

	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Constants are now defined in data-source-host-agents-constants.go

// HostAgentDataSourceModel represents the data model for a single host agent
type HostAgentDataSourceModel struct {
	SnapshotID types.String   `tfsdk:"snapshot_id"`
	Label      types.String   `tfsdk:"label"`
	Host       types.String   `tfsdk:"host"`
	Plugin     types.String   `tfsdk:"plugin"`
	Tags       []types.String `tfsdk:"tags"`
}

// HostAgentsDataSourceModel represents the data model for the host agents data source
type HostAgentsDataSourceModel struct {
	ID     types.String               `tfsdk:"id"`
	Filter types.String               `tfsdk:"filter"`
	Items  []HostAgentDataSourceModel `tfsdk:"items"`
}

// NewHostAgentsDataSourceFramework creates a new data source for host agents
func NewHostAgentsDataSourceFramework() datasource.DataSource {
	return &hostAgentsDataSourceFramework{}
}

type hostAgentsDataSourceFramework struct {
	instanaAPI restapi.InstanaAPI
}

func (d *hostAgentsDataSourceFramework) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaHostAgentsFramework
}

func (d *hostAgentsDataSourceFramework) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: HostAgentDescDataSource,
		Attributes: map[string]schema.Attribute{
			HostAgentFieldID: schema.StringAttribute{
				Description: HostAgentDescID,
				Computed:    true,
			},
			HostAgentFieldFilter: schema.StringAttribute{
				Description: HostAgentDescFilter,
				Required:    true,
			},
			HostAgentFieldItems: schema.ListNestedAttribute{
				Description: HostAgentDescItems,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						HostAgentFieldSnapshotId: schema.StringAttribute{
							Description: HostAgentDescSnapshotID,
							Computed:    true,
						},
						HostAgentFieldLabel: schema.StringAttribute{
							Description: HostAgentDescLabel,
							Computed:    true,
						},
						HostAgentFieldHost: schema.StringAttribute{
							Description: HostAgentDescHost,
							Computed:    true,
						},
						HostAgentFieldPlugin: schema.StringAttribute{
							Description: HostAgentDescPlugin,
							Computed:    true,
						},
						HostAgentFieldTags: schema.ListAttribute{
							Description: HostAgentDescTags,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *hostAgentsDataSourceFramework) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*restapi.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			HostAgentErrUnexpectedConfigureType,
			fmt.Sprintf(HostAgentErrUnexpectedConfigureTypeDetail, req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *hostAgentsDataSourceFramework) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostAgentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the filter from the configuration
	filter := data.Filter.ValueString()
	queryParams := map[string]string{HostAgentQueryParamQuery: filter}

	// Get host agents by query
	hostAgents, err := d.instanaAPI.HostAgents().GetByQuery(queryParams)
	if err != nil {
		resp.Diagnostics.AddError(
			HostAgentErrReadingAgents,
			fmt.Sprintf(HostAgentErrReadingAgentsDetail, err),
		)
		return
	}

	// Set ID to current timestamp
	data.ID = types.StringValue(time.Now().UTC().String())

	// Map host agents to model
	items := make([]HostAgentDataSourceModel, len(*hostAgents))
	for i, hostAgent := range *hostAgents {
		// Map tags
		tags := make([]types.String, len(hostAgent.Tags))
		for j, tag := range hostAgent.Tags {
			tags[j] = types.StringValue(tag)
		}

		items[i] = HostAgentDataSourceModel{
			SnapshotID: types.StringValue(hostAgent.SnapshotID),
			Label:      types.StringValue(hostAgent.Label),
			Host:       types.StringValue(hostAgent.Host),
			Plugin:     types.StringValue(hostAgent.Plugin),
			Tags:       tags,
		}
	}
	data.Items = items

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
