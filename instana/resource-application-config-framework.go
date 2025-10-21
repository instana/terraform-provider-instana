package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// ApplicationConfigFieldAccessRules field name for access rules
const ApplicationConfigFieldAccessRules = "access_rules"

// ApplicationConfigModel represents the data model for the application configuration resource
type ApplicationConfigModel struct {
	ID            types.String `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	Scope         types.String `tfsdk:"scope"`
	BoundaryScope types.String `tfsdk:"boundary_scope"`
	TagFilter     types.String `tfsdk:"tag_filter"`
	AccessRules   types.List   `tfsdk:"access_rules"`
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
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "The tag filter expression",
					},
					ApplicationConfigFieldAccessRules: schema.ListNestedAttribute{
						Required:    true,
						Description: "The access rules applied to the application config",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_type": schema.StringAttribute{
									Required:    true,
									Description: "The access type of the given access rule",
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedAccessTypes.ToStringSlice()...),
									},
								},
								"related_id": schema.StringAttribute{
									Optional:    true,
									Description: "The id of the related entity (user, api_token, etc.) of the given access rule",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 64),
									},
								},
								"relation_type": schema.StringAttribute{
									Required:    true,
									Description: "The relation type of the given access rule",
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

	// Map access rules
	accessRules, d := r.mapAccessRulesToState(ctx, config.AccessRules)
	diags.Append(d...)
	if !diags.HasError() {
		model.AccessRules = accessRules
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *applicationConfigResourceFramework) mapAccessRulesToState(ctx context.Context, accessRules []restapi.AccessRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// If there are no access rules, return an empty list
	if len(accessRules) == 0 {
		return types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"access_type":   types.StringType,
					"related_id":    types.StringType,
					"relation_type": types.StringType,
				},
			},
			[]attr.Value{},
		), diags
	}

	// Create elements for each access rule
	elements := make([]attr.Value, len(accessRules))
	for i, rule := range accessRules {
		// Create a map for the rule attributes
		ruleMap := map[string]attr.Value{
			"access_type":   types.StringValue(string(rule.AccessType)),
			"relation_type": types.StringValue(string(rule.RelationType)),
		}

		// Handle related ID
		if rule.RelatedID != nil {
			ruleMap["related_id"] = types.StringValue(*rule.RelatedID)
		} else {
			ruleMap["related_id"] = types.StringNull()
		}

		// Create object value
		objValue, d := types.ObjectValue(
			map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
			ruleMap,
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		elements[i] = objValue
	}

	// Create list value
	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		},
		elements,
	), diags
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

	// Map access rules
	accessRules, d := r.mapAccessRulesFromState(ctx, model.AccessRules)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &restapi.ApplicationConfig{
		ID:                  id,
		Label:               model.Label.ValueString(),
		Scope:               restapi.ApplicationConfigScope(model.Scope.ValueString()),
		BoundaryScope:       restapi.BoundaryScope(model.BoundaryScope.ValueString()),
		TagFilterExpression: tagFilter,
		AccessRules:         accessRules,
	}, diags
}

func (r *applicationConfigResourceFramework) mapAccessRulesFromState(ctx context.Context, accessRulesList types.List) ([]restapi.AccessRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	var accessRules []restapi.AccessRule

	if accessRulesList.IsNull() {
		return accessRules, diags
	}

	// Get the list of objects
	var accessRuleObjects []types.Object
	diags.Append(accessRulesList.ElementsAs(ctx, &accessRuleObjects, false)...)
	if diags.HasError() {
		return accessRules, diags
	}

	accessRules = make([]restapi.AccessRule, len(accessRuleObjects))
	for i, obj := range accessRuleObjects {
		var accessType types.String
		var relationType types.String
		var relatedID types.String

		// Extract values from the object
		attrMap := obj.Attributes()

		// Get access_type
		if v, ok := attrMap["access_type"]; ok {
			if str, ok := v.(types.String); ok {
				accessType = str
			} else {
				diags.AddError(
					"Invalid attribute type",
					"access_type must be a string",
				)
				return accessRules, diags
			}
		} else {
			diags.AddError(
				"Missing attribute",
				"access_type attribute is required for access rule",
			)
			return accessRules, diags
		}

		// Get relation_type
		if v, ok := attrMap["relation_type"]; ok {
			if str, ok := v.(types.String); ok {
				relationType = str
			} else {
				diags.AddError(
					"Invalid attribute type",
					"relation_type must be a string",
				)
				return accessRules, diags
			}
		} else {
			diags.AddError(
				"Missing attribute",
				"relation_type attribute is required for access rule",
			)
			return accessRules, diags
		}

		// Get related_id (optional)
		if v, ok := attrMap["related_id"]; ok {
			if str, ok := v.(types.String); ok {
				relatedID = str
			}
		}

		rule := restapi.AccessRule{
			AccessType:   restapi.AccessType(accessType.ValueString()),
			RelationType: restapi.RelationType(relationType.ValueString()),
		}

		// Handle related ID
		if !relatedID.IsNull() && relatedID.ValueString() != "" {
			relatedIDValue := relatedID.ValueString()
			rule.RelatedID = &relatedIDValue
		}

		accessRules[i] = rule
	}

	return accessRules, diags
}

// Made with Bob
