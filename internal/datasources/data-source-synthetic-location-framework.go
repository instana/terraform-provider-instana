package datasources

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Constants are now defined in data-source-synthetic-location-constants.go

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
		Description: SyntheticLocationDescDataSource,
		Attributes: map[string]schema.Attribute{
			SyntheticLocationFieldID: schema.StringAttribute{
				Description: SyntheticLocationDescID,
				Computed:    true,
			},
			SyntheticLocationFieldLabel: schema.StringAttribute{
				Description: SyntheticLocationDescLabel,
				Optional:    true,
			},
			SyntheticLocationFieldDescription: schema.StringAttribute{
				Description: SyntheticLocationDescDescription,
				Optional:    true,
				Computed:    true,
			},
			SyntheticLocationFieldLocationType: schema.StringAttribute{
				Description: SyntheticLocationDescLocationType,
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

	providerMeta, ok := req.ProviderData.(*restapi.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			SyntheticLocationErrUnexpectedConfigureType,
			fmt.Sprintf(SyntheticLocationErrUnexpectedConfigureTypeDetail, req.ProviderData),
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
			SyntheticLocationErrReadingLocations,
			fmt.Sprintf(SyntheticLocationErrReadingLocationsDetail, err),
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
			SyntheticLocationErrNotFound,
			fmt.Sprintf(SyntheticLocationErrNotFoundDetail, label, locationType),
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
