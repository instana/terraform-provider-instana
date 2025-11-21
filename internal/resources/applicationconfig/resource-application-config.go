package applicationconfig

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
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

// ============================================================================
// Resource Factory
// ============================================================================

// NewApplicationConfigResourceHandle creates the resource handle for Application Configuration
func NewApplicationConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.ApplicationConfig] {
	return &applicationConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaApplicationConfig,
			Schema: schema.Schema{
				Description: ApplicationConfigDescResource,
				Attributes: map[string]schema.Attribute{
					ApplicationConfigFieldID: schema.StringAttribute{
						Computed:    true,
						Description: ApplicationConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					ApplicationConfigFieldLabel: schema.StringAttribute{
						Required:    true,
						Description: ApplicationConfigDescLabel,
					},
					ApplicationConfigFieldScope: schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(string(restapi.ApplicationConfigScopeIncludeNoDownstream)),
						Description: ApplicationConfigDescScope,
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedApplicationConfigScopes.ToStringSlice()...),
						},
					},
					ApplicationConfigFieldBoundaryScope: schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(string(restapi.BoundaryScopeDefault)),
						Description: ApplicationConfigDescBoundaryScope,
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedApplicationConfigBoundaryScopes.ToStringSlice()...),
						},
					},
					ApplicationConfigFieldTagFilter: schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: ApplicationConfigDescTagFilter,
					},
					ApplicationConfigFieldAccessRules: schema.ListNestedAttribute{
						Required:    true,
						Description: ApplicationConfigDescAccessRules,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								ApplicationConfigFieldAccessType: schema.StringAttribute{
									Required:    true,
									Description: ApplicationConfigDescAccessType,
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedAccessTypes.ToStringSlice()...),
									},
								},
								ApplicationConfigFieldRelatedID: schema.StringAttribute{
									Optional:    true,
									Description: ApplicationConfigDescRelatedID,
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 64),
									},
								},
								ApplicationConfigFieldRelationType: schema.StringAttribute{
									Required:    true,
									Description: ApplicationConfigDescRelationType,
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedRelationTypes.ToStringSlice()...),
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 4,
		},
	}
}

// ============================================================================
// Resource Implementation
// ============================================================================

type applicationConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *applicationConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for application configs
func (r *applicationConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationConfig] {
	return api.ApplicationConfigs()
}

// SetComputedFields sets computed fields in the plan (none for this resource)
func (r *applicationConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ============================================================================
// API to State Mapping
// ============================================================================

// UpdateState converts API data object to Terraform state
func (r *applicationConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.ApplicationConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model ApplicationConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	model.ID = types.StringValue(config.ID)
	model.Label = types.StringValue(config.Label)
	model.Scope = types.StringValue(string(config.Scope))
	model.BoundaryScope = types.StringValue(string(config.BoundaryScope))

	// Map tag filter with error handling
	if model.TagFilter.IsNull() || model.TagFilter.IsUnknown() {
		if err := r.mapTagFilterToState(config, &model); err != nil {
			return err
		}
	}

	// Map access rules
	model.AccessRules = r.mapAccessRulesToState(config.AccessRules)

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapTagFilterToState handles tag filter normalization and mapping to state
func (r *applicationConfigResource) mapTagFilterToState(config *restapi.ApplicationConfig, model *ApplicationConfigModel) diag.Diagnostics {
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic(
					ApplicationConfigErrConvertingTagFilter,
					fmt.Sprintf(ApplicationConfigErrFailedToConvert, err),
				),
			}
		}
		model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
	} else {
		model.TagFilter = types.StringNull()
	}
	return nil
}

// mapAccessRulesToState converts access rules from API to state format
func (r *applicationConfigResource) mapAccessRulesToState(accessRules []restapi.AccessRule) []AccessRuleModel {
	// If there are no access rules, return an empty slice
	if len(accessRules) == 0 {
		return []AccessRuleModel{}
	}

	// Create slice for access rule models
	models := make([]AccessRuleModel, len(accessRules))
	for i, rule := range accessRules {
		models[i] = AccessRuleModel{
			AccessType:   types.StringValue(string(rule.AccessType)),
			RelationType: types.StringValue(string(rule.RelationType)),
			RelatedID:    util.SetStringPointerToState(rule.RelatedID),
		}
	}

	return models
}

// ============================================================================
// State to API Mapping
// ============================================================================

// MapStateToDataObject converts Terraform state to API data object
func (r *applicationConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.ApplicationConfig, diag.Diagnostics) {
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

	// Extract ID
	id := r.extractID(&model)

	// Parse and map tag filter
	tagFilter, tagDiags := r.mapTagFilterFromState(&model)
	if tagDiags.HasError() {
		diags.Append(tagDiags...)
		return nil, diags
	}

	// Map access rules
	accessRules := r.mapAccessRulesFromState(model.AccessRules)

	return &restapi.ApplicationConfig{
		ID:                  id,
		Label:               model.Label.ValueString(),
		Scope:               restapi.ApplicationConfigScope(model.Scope.ValueString()),
		BoundaryScope:       restapi.BoundaryScope(model.BoundaryScope.ValueString()),
		TagFilterExpression: tagFilter,
		AccessRules:         accessRules,
	}, diags
}

// extractID extracts the ID from the model, returning empty string if null
func (r *applicationConfigResource) extractID(model *ApplicationConfigModel) string {
	if !model.ID.IsNull() {
		return model.ID.ValueString()
	}
	return ""
}

// mapTagFilterFromState parses and converts tag filter from state to API format
func (r *applicationConfigResource) mapTagFilterFromState(model *ApplicationConfigModel) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.TagFilter.IsNull() {
		return nil, diags
	}

	tagFilterStr := model.TagFilter.ValueString()
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(tagFilterStr)
	if err != nil {
		diags.AddError(
			ApplicationConfigErrParsingTagFilter,
			fmt.Sprintf(ApplicationConfigErrFailedToParse, err),
		)
		return nil, diags
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), diags
}

// mapAccessRulesFromState converts access rules from state to API format
func (r *applicationConfigResource) mapAccessRulesFromState(accessRulesModels []AccessRuleModel) []restapi.AccessRule {
	// If there are no access rules, return an empty slice
	if len(accessRulesModels) == 0 {
		return []restapi.AccessRule{}
	}

	accessRules := make([]restapi.AccessRule, len(accessRulesModels))
	for i, model := range accessRulesModels {
		rule := restapi.AccessRule{
			AccessType:   restapi.AccessType(model.AccessType.ValueString()),
			RelationType: restapi.RelationType(model.RelationType.ValueString()),
		}

		// Handle related ID (optional)
		if !model.RelatedID.IsNull() && model.RelatedID.ValueString() != "" {
			relatedIDValue := model.RelatedID.ValueString()
			rule.RelatedID = &relatedIDValue
		}

		accessRules[i] = rule
	}

	return accessRules
}
