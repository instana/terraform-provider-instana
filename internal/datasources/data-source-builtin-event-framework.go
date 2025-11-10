package datasources

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Constants are now defined in data-source-builtin-event-constants.go

// BuiltinEventDataSourceModel represents the data model for the builtin event data source
type BuiltinEventDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ShortPluginID types.String `tfsdk:"short_plugin_id"`
	Severity      types.String `tfsdk:"severity"`
	SeverityCode  types.Int64  `tfsdk:"severity_code"`
	Triggering    types.Bool   `tfsdk:"triggering"`
	Enabled       types.Bool   `tfsdk:"enabled"`
}

// NewBuiltinEventDataSourceFramework creates a new data source for builtin events
func NewBuiltinEventDataSourceFramework() datasource.DataSource {
	return &builtinEventDataSourceFramework{}
}

type builtinEventDataSourceFramework struct {
	instanaAPI restapi.InstanaAPI
}

func (d *builtinEventDataSourceFramework) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaBuiltinEventFramework
}

func (d *builtinEventDataSourceFramework) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: BuiltinEventDescDataSource,
		Attributes: map[string]schema.Attribute{
			BuiltinEventFieldID: schema.StringAttribute{
				Description: BuiltinEventDescID,
				Computed:    true,
			},
			BuiltinEventSpecificationFieldName: schema.StringAttribute{
				Description: BuiltinEventDescName,
				Required:    true,
			},
			BuiltinEventSpecificationFieldDescription: schema.StringAttribute{
				Description: BuiltinEventDescDescription,
				Computed:    true,
			},
			BuiltinEventSpecificationFieldShortPluginID: schema.StringAttribute{
				Description: BuiltinEventDescShortPluginID,
				Required:    true,
			},
			BuiltinEventSpecificationFieldSeverity: schema.StringAttribute{
				Description: BuiltinEventDescSeverity,
				Computed:    true,
			},
			BuiltinEventSpecificationFieldSeverityCode: schema.Int64Attribute{
				Description: BuiltinEventDescSeverityCode,
				Computed:    true,
			},
			BuiltinEventSpecificationFieldTriggering: schema.BoolAttribute{
				Description: BuiltinEventDescTriggering,
				Computed:    true,
			},
			BuiltinEventSpecificationFieldEnabled: schema.BoolAttribute{
				Description: BuiltinEventDescEnabled,
				Computed:    true,
			},
		},
	}
}

func (d *builtinEventDataSourceFramework) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			BuiltinEventErrUnexpectedConfigureType,
			fmt.Sprintf(BuiltinEventErrUnexpectedConfigureTypeDetail, req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *builtinEventDataSourceFramework) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BuiltinEventDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the name and short plugin ID from the configuration
	name := data.Name.ValueString()
	shortPluginID := data.ShortPluginID.ValueString()

	// Get all builtin events
	events, err := d.instanaAPI.BuiltinEventSpecifications().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			BuiltinEventErrReadingEvents,
			fmt.Sprintf(BuiltinEventErrReadingEventsDetail, err),
		)
		return
	}

	// Find the event with the matching name and short plugin ID
	var matchingEvent *restapi.BuiltinEventSpecification
	for _, event := range *events {
		if event.Name == name && event.ShortPluginID == shortPluginID {
			matchingEvent = event
			break
		}
	}

	if matchingEvent == nil {
		resp.Diagnostics.AddError(
			BuiltinEventErrNotFound,
			fmt.Sprintf(BuiltinEventErrNotFoundDetail, name, shortPluginID),
		)
		return
	}

	// Update the data model with the event details
	data.ID = types.StringValue(matchingEvent.ID)

	// Handle description which is a pointer
	data.Description = util.setStringPointerToState(matchingEvent.Description)

	// Convert severity from API representation to Terraform representation
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(matchingEvent.Severity)
	if err != nil {
		resp.Diagnostics.AddError(
			BuiltinEventErrConvertingSeverity,
			fmt.Sprintf(BuiltinEventErrConvertingSeverityDetail, err),
		)
		return
	}

	data.Severity = types.StringValue(severity)
	data.SeverityCode = types.Int64Value(int64(matchingEvent.Severity))
	data.Triggering = types.BoolValue(matchingEvent.Triggering)
	data.Enabled = types.BoolValue(matchingEvent.Enabled)

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
