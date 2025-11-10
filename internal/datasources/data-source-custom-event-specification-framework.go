package datasources

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Constants are now defined in data-source-custom-event-specification-constants.go

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
		Description: CustomEventSpecificationDescDataSource,
		Attributes: map[string]schema.Attribute{
			CustomEventSpecificationFieldID: schema.StringAttribute{
				Description: CustomEventSpecificationDescID,
				Computed:    true,
			},
			CustomEventSpecificationFieldName: schema.StringAttribute{
				Description: CustomEventSpecificationDescName,
				Required:    true,
			},
			CustomEventSpecificationFieldDescription: schema.StringAttribute{
				Description: CustomEventSpecificationDescDescription,
				Computed:    true,
			},
			CustomEventSpecificationFieldEntityType: schema.StringAttribute{
				Description: CustomEventSpecificationDescEntityType,
				Required:    true,
			},
			CustomEventSpecificationFieldTriggering: schema.BoolAttribute{
				Description: CustomEventSpecificationDescTriggering,
				Computed:    true,
			},
			CustomEventSpecificationFieldEnabled: schema.BoolAttribute{
				Description: CustomEventSpecificationDescEnabled,
				Computed:    true,
			},
			CustomEventSpecificationFieldQuery: schema.StringAttribute{
				Description: CustomEventSpecificationDescQuery,
				Computed:    true,
			},
			CustomEventSpecificationFieldExpirationTime: schema.Int64Attribute{
				Description: CustomEventSpecificationDescExpirationTime,
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
			CustomEventSpecificationErrUnexpectedConfigureType,
			fmt.Sprintf(CustomEventSpecificationErrUnexpectedConfigureTypeDetail, req.ProviderData),
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
			CustomEventSpecificationErrReadingSpecs,
			fmt.Sprintf(CustomEventSpecificationErrReadingSpecsDetail, err),
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
			CustomEventSpecificationErrNotFound,
			fmt.Sprintf(CustomEventSpecificationErrNotFoundDetail, name, entityType),
		)
		return
	}

	// Update the data model with the specification details
	data.ID = types.StringValue(matchingSpec.ID)

	// Handle description which is a pointer
	data.Description = util.setStringPointerToState(matchingSpec.Description)

	// Handle query which is a pointer
	data.Query = util.setStringPointerToState(matchingSpec.Query)

	// Handle expiration time which is a pointer
	if matchingSpec.ExpirationTime != nil {
		data.ExpirationTime = setInt64PointerToState(matchingSpec.ExpirationTime)
	} else {
		data.ExpirationTime = types.Int64Null()
	}

	data.Triggering = types.BoolValue(matchingSpec.Triggering)
	data.Enabled = types.BoolValue(matchingSpec.Enabled)

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
