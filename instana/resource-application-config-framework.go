package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaApplicationConfigFramework the name of the terraform-provider-instana resource to manage application config
const ResourceInstanaApplicationConfigFramework = "application_config"

// ApplicationConfigModel represents the data model for the application configuration resource
type ApplicationConfigModel struct {
	ID            types.String `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	Scope         types.String `tfsdk:"scope"`
	BoundaryScope types.String `tfsdk:"boundary_scope"`
	TagFilter     types.String `tfsdk:"tag_filter"`
}

// NewApplicationConfigResourceHandleFramework creates the resource handle for Application Configuration
func NewApplicationConfigResourceHandleFramework() ResourceHandleFramework[*restapi.ApplicationConfig] {
	return &applicationConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaApplicationConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages application configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the application configuration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					ApplicationConfigFieldLabel: schema.StringAttribute{
						Required:    true,
						Description: "The label of the application config",
					},
					ApplicationConfigFieldScope: schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(string(restapi.ApplicationConfigScopeIncludeNoDownstream)),
						Description: "The scope of the application config",
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedApplicationConfigScopes.ToStringSlice()...),
						},
					},
					ApplicationConfigFieldBoundaryScope: schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(string(restapi.BoundaryScopeDefault)),
						Description: "The boundary scope of the application config",
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedApplicationConfigBoundaryScopes.ToStringSlice()...),
						},
					},
					ApplicationConfigFieldTagFilter: schema.StringAttribute{
						Optional:    true,
						Description: "The tag filter expression",
					},
				},
			},
			SchemaVersion: 4,
		},
	}
}

type applicationConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *applicationConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *applicationConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationConfig] {
	return api.ApplicationConfigs()
}

func (r *applicationConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *applicationConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, config *restapi.ApplicationConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the config
	model := ApplicationConfigModel{
		ID:            types.StringValue(config.ID),
		Label:         types.StringValue(config.Label),
		Scope:         types.StringValue(string(config.Scope)),
		BoundaryScope: types.StringValue(string(config.BoundaryScope)),
	}

	// Set tag filter
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Error converting tag filter",
					fmt.Sprintf("Failed to convert tag filter: %s", err),
				),
			}
		}
		if normalizedTagFilterString != nil {
			model.TagFilter = types.StringValue(*normalizedTagFilterString)
		} else {
			model.TagFilter = types.StringNull()
		}
	} else {
		model.TagFilter = types.StringNull()
	}

	// Set the entire model to state
	diags = state.Set(ctx, model)
	return diags
}

func (r *applicationConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.ApplicationConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model ApplicationConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map tag filter
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() {
		tagFilterStr := model.TagFilter.ValueString()
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(tagFilterStr)
		if err != nil {
			diags.AddError(
				"Error parsing tag filter",
				fmt.Sprintf("Failed to parse tag filter: %s", err),
			)
			return nil, diags
		}

		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
	}

	return &restapi.ApplicationConfig{
		ID:                  id,
		Label:               model.Label.ValueString(),
		Scope:               restapi.ApplicationConfigScope(model.Scope.ValueString()),
		BoundaryScope:       restapi.BoundaryScope(model.BoundaryScope.ValueString()),
		TagFilterExpression: tagFilter,
	}, diags
}

// Made with Bob
