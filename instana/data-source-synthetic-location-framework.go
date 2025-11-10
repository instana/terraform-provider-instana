package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaSyntheticLocationFramework the name of the terraform-provider-instana data source to read synthetic locations
const DataSourceInstanaSyntheticLocationFramework = "synthetic_location"

const (
	//SyntheticLocationFieldLabel constant value for the schema field label
	SyntheticLocationFieldLabel = "label"
	//SyntheticLocationFieldDescription constant value for the computed schema field description
	SyntheticLocationFieldDescription = "description"
	//SyntheticLocationFieldLocationType constant value for the schema field location_type
	SyntheticLocationFieldLocationType = "location_type"
	//DataSourceSyntheticLocation the name of the terraform-provider-instana data sourcefor synthetic location specifications
	DataSourceSyntheticLocation = "instana_synthetic_location"
)

// SyntheticLocationDataSourceModel represents the data model for the synthetic location data source
type SyntheticLocationDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Description  types.String `tfsdk:"description"`
	LocationType types.String `tfsdk:"location_type"`
}

// NewSyntheticLocationDataSourceFramework creates a new data source for synthetic locations
func NewSyntheticLocationDataSourceFramework() datasource.DataSource {
	return &syntheticLocationDataSourceFramework{}
}

type syntheticLocationDataSourceFramework struct {
	instanaAPI restapi.InstanaAPI
}

func (d *syntheticLocationDataSourceFramework) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaSyntheticLocationFramework
}

func (d *syntheticLocationDataSourceFramework) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Instana synthetic locations. Synthetic locations are the locations from which synthetic tests are executed.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the synthetic location.",
				Computed:    true,
			},
			SyntheticLocationFieldLabel: schema.StringAttribute{
				Description: "Friendly name of the Synthetic Location",
				Optional:    true,
			},
			SyntheticLocationFieldDescription: schema.StringAttribute{
				Description: "The description of the Synthetic location",
				Optional:    true,
				Computed:    true,
			},
			SyntheticLocationFieldLocationType: schema.StringAttribute{
				Description: "Indicates if the location is public or private",
				Optional:    true,
				Validators:  []validator.String{
					// We'll implement custom validators in the Configure method
				},
			},
		},
	}
}

func (d *syntheticLocationDataSourceFramework) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *syntheticLocationDataSourceFramework) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SyntheticLocationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all synthetic locations
	syntheticLocations, err := d.instanaAPI.SyntheticLocation().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading synthetic locations",
			fmt.Sprintf("Could not read synthetic locations: %s", err),
		)
		return
	}

	// Find the matching synthetic location
	label := data.Label.ValueString()
	locationType := data.LocationType.ValueString()

	var matchingLocation *restapi.SyntheticLocation
	for _, location := range *syntheticLocations {
		if location.Label == label && location.LocationType == locationType {
			matchingLocation = location
			break
		}
	}

	if matchingLocation == nil {
		resp.Diagnostics.AddError(
			"No matching synthetic location found",
			fmt.Sprintf("No synthetic location found with label '%s' and location type '%s'", label, locationType),
		)
		return
	}

	// Map the synthetic location to the model
	data.ID = types.StringValue(matchingLocation.ID)
	data.Label = types.StringValue(matchingLocation.Label)
	data.Description = types.StringValue(matchingLocation.Description)
	data.LocationType = types.StringValue(matchingLocation.LocationType)

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
