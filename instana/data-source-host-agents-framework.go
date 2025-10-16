package instana

import (
	"context"
	"fmt"
	"time"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaHostAgentsFramework the name of the terraform-provider-instana data source to read host agents
const DataSourceInstanaHostAgentsFramework = "host_agents"

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
		Description: "Data source for Instana host agents. Host agents are the Instana agents installed on hosts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the data source.",
				Computed:    true,
			},
			HostAgentFieldFilter: schema.StringAttribute{
				Description: "Dynamic Focus Query filter.",
				Required:    true,
			},
			HostAgentFieldItems: schema.ListNestedAttribute{
				Description: "A list of host agents.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						HostAgentFieldSnapshotId: schema.StringAttribute{
							Description: "The snapshot ID of the host agent.",
							Computed:    true,
						},
						HostAgentFieldLabel: schema.StringAttribute{
							Description: "The label of the host agent.",
							Computed:    true,
						},
						HostAgentFieldHost: schema.StringAttribute{
							Description: "The host identifier of the host agent.",
							Computed:    true,
						},
						HostAgentFieldPlugin: schema.StringAttribute{
							Description: "The plugin of the host agent.",
							Computed:    true,
						},
						HostAgentFieldTags: schema.ListAttribute{
							Description: "The tags of the host agent.",
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

	providerMeta, ok := req.ProviderData.(*ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
	queryParams := map[string]string{"query": filter}

	// Get host agents by query
	hostAgents, err := d.instanaAPI.HostAgents().GetByQuery(queryParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading host agents",
			fmt.Sprintf("Could not read host agents: %s", err),
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

// Made with Bob
