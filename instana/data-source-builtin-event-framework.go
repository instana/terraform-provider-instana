package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaBuiltinEventFramework the name of the terraform-provider-instana data source to read builtin events
const DataSourceInstanaBuiltinEventFramework = "builtin_event_spec"

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
		Description: "Data source for an Instana builtin event specification. Builtin events are predefined events in Instana.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the builtin event.",
				Computed:    true,
			},
			BuiltinEventSpecificationFieldName: schema.StringAttribute{
				Description: "The name of the builtin event.",
				Required:    true,
			},
			BuiltinEventSpecificationFieldDescription: schema.StringAttribute{
				Description: "The description text of the builtin event.",
				Computed:    true,
			},
			BuiltinEventSpecificationFieldShortPluginID: schema.StringAttribute{
				Description: "The plugin id for which the builtin event is created.",
				Required:    true,
			},
			BuiltinEventSpecificationFieldSeverity: schema.StringAttribute{
				Description: "The severity (WARNING, CRITICAL, etc.) of the builtin event.",
				Computed:    true,
			},
			BuiltinEventSpecificationFieldSeverityCode: schema.Int64Attribute{
				Description: "The severity code used by Instana API (5, 10, etc.) of the builtin event.",
				Computed:    true,
			},
			BuiltinEventSpecificationFieldTriggering: schema.BoolAttribute{
				Description: "Indicates if an incident is triggered the builtin event or not.",
				Computed:    true,
			},
			BuiltinEventSpecificationFieldEnabled: schema.BoolAttribute{
				Description: "Indicates if the builtin event is enabled or not.",
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
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
			"Error reading builtin events",
			fmt.Sprintf("Could not read builtin events: %s", err),
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
			"Builtin event not found",
			fmt.Sprintf("No built in event found for name '%s' and short plugin ID '%s'", name, shortPluginID),
		)
		return
	}

	// Update the data model with the event details
	data.ID = types.StringValue(matchingEvent.ID)

	// Handle description which is a pointer
	if matchingEvent.Description != nil {
		data.Description = types.StringValue(*matchingEvent.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Convert severity from API representation to Terraform representation
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(matchingEvent.Severity)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting severity",
			fmt.Sprintf("Could not convert severity: %s", err),
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
