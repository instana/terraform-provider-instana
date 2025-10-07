package applicationconfig

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
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

// NewApplicationConfigResourceHandleFramework creates the resource handle for Application Configuration
func NewApplicationConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.ApplicationConfig] {
	return &applicationConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaApplicationConfigFramework,
			Schema: schema.Schema{
				Description: ApplicationConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
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
								"access_type": schema.StringAttribute{
									Required:    true,
									Description: ApplicationConfigDescAccessType,
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedAccessTypes.ToStringSlice()...),
									},
								},
								"related_id": schema.StringAttribute{
									Optional:    true,
									Description: ApplicationConfigDescRelatedID,
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 64),
									},
								},
								"relation_type": schema.StringAttribute{
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

type applicationConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *applicationConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *applicationConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationConfig] {
	return api.ApplicationConfigs()
}

func (r *applicationConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *applicationConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.ApplicationConfig) diag.Diagnostics {
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
					ApplicationConfigErrConvertingTagFilter,
					fmt.Sprintf("Failed to convert tag filter: %s", err),
				),
			}
		}
		model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)

	} else {
		model.TagFilter = types.StringNull()
	}

	// Map access rules
	model.AccessRules = r.mapAccessRulesToState(config.AccessRules)

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *applicationConfigResourceFramework) mapAccessRulesToState(accessRules []restapi.AccessRule) []AccessRuleModel {
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
				ApplicationConfigErrParsingTagFilter,
				fmt.Sprintf("Failed to parse tag filter: %s", err),
			)
			return nil, diags
		}

		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
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

func (r *applicationConfigResourceFramework) mapAccessRulesFromState(accessRulesModels []AccessRuleModel) []restapi.AccessRule {
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
