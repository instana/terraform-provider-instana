package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaCustomEventSpecificationFramework the name of the terraform-provider-instana data source to read custom event specifications
const DataSourceInstanaCustomEventSpecificationFramework = "custom_event_spec"

// CustomEventSpecificationDataSourceModel represents the data model for the custom event specification data source
type CustomEventSpecificationDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	EntityType     types.String `tfsdk:"entity_type"`
	Triggering     types.Bool   `tfsdk:"triggering"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Query          types.String `tfsdk:"query"`
	ExpirationTime types.Int64  `tfsdk:"expiration_time"`
}

// NewCustomEventSpecificationDataSourceFramework creates a new data source for custom event specifications
func NewCustomEventSpecificationDataSourceFramework() datasource.DataSource {
	return &customEventSpecificationDataSourceFramework{}
}

type customEventSpecificationDataSourceFramework struct {
	instanaAPI restapi.InstanaAPI
}

func (d *customEventSpecificationDataSourceFramework) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaCustomEventSpecificationFramework
}

func (d *customEventSpecificationDataSourceFramework) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for an Instana custom event specification. Custom events are user-defined events in Instana.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the custom event specification.",
				Computed:    true,
			},
			CustomEventSpecificationFieldName: schema.StringAttribute{
				Description: "The name of the custom event specification.",
				Required:    true,
			},
			CustomEventSpecificationFieldDescription: schema.StringAttribute{
				Description: "The description of the custom event specification.",
				Computed:    true,
			},
			CustomEventSpecificationFieldEntityType: schema.StringAttribute{
				Description: "The entity type for which the custom event specification is created.",
				Required:    true,
			},
			CustomEventSpecificationFieldTriggering: schema.BoolAttribute{
				Description: "Indicates if an incident is triggered the custom event or not.",
				Computed:    true,
			},
			CustomEventSpecificationFieldEnabled: schema.BoolAttribute{
				Description: "Indicates if the custom event is enabled or not.",
				Computed:    true,
			},
			CustomEventSpecificationFieldQuery: schema.StringAttribute{
				Description: "Dynamic focus query for the custom event specification.",
				Computed:    true,
			},
			CustomEventSpecificationFieldExpirationTime: schema.Int64Attribute{
				Description: "The expiration time (grace period) to wait before the issue is closed.",
				Computed:    true,
			},
		},
	}
}

func (d *customEventSpecificationDataSourceFramework) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *customEventSpecificationDataSourceFramework) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CustomEventSpecificationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the name and entity type from the configuration
	name := data.Name.ValueString()
	entityType := data.EntityType.ValueString()

	// Get all custom event specifications
	specs, err := d.instanaAPI.CustomEventSpecifications().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading custom event specifications",
			fmt.Sprintf("Could not read custom event specifications: %s", err),
		)
		return
	}

	// Find the specification with the matching name and entity type
	var matchingSpec *restapi.CustomEventSpecification
	for _, spec := range *specs {
		if spec.Name == name && spec.EntityType == entityType {
			matchingSpec = spec
			break
		}
	}

	if matchingSpec == nil {
		resp.Diagnostics.AddError(
			"Custom event specification not found",
			fmt.Sprintf("No custom event specification found for name '%s' and entity type '%s'", name, entityType),
		)
		return
	}

	// Update the data model with the specification details
	data.ID = types.StringValue(matchingSpec.ID)

	// Handle description which is a pointer
	if matchingSpec.Description != nil {
		data.Description = types.StringValue(*matchingSpec.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Handle query which is a pointer
	if matchingSpec.Query != nil {
		data.Query = types.StringValue(*matchingSpec.Query)
	} else {
		data.Query = types.StringNull()
	}

	// Handle expiration time which is a pointer
	if matchingSpec.ExpirationTime != nil {
		data.ExpirationTime = types.Int64Value(int64(*matchingSpec.ExpirationTime))
	} else {
		data.ExpirationTime = types.Int64Null()
	}

	data.Triggering = types.BoolValue(matchingSpec.Triggering)
	data.Enabled = types.BoolValue(matchingSpec.Enabled)

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
