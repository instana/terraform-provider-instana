package customdashboard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ============================================================================
// Resource Factory
// ============================================================================

// NewCustomDashboardResourceHandleFramework creates the resource handle for Custom Dashboards
func NewCustomDashboardResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.CustomDashboard] {
	return &customDashboardResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaCustomDashboardFramework,
			Schema: schema.Schema{
				Description: CustomDashboardDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: CustomDashboardDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					CustomDashboardFieldTitle: schema.StringAttribute{
						Required:    true,
						Description: CustomDashboardDescTitle,
					},
					CustomDashboardFieldWidgets: schema.StringAttribute{
						Required:    true,
						Description: CustomDashboardDescWidgets,
						CustomType:  jsontypes.NormalizedType{},
						PlanModifiers: []planmodifier.String{
							CanonicalizeJSONPlanModifier{},
						},
					},
					CustomDashboardFieldAccessRule: schema.ListNestedAttribute{
						Description: CustomDashboardDescAccessRule,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								CustomDashboardFieldAccessRuleAccessType: schema.StringAttribute{
									Required:    true,
									Description: CustomDashboardDescAccessRuleAccessType,
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedAccessTypes.ToStringSlice()...),
									},
								},
								CustomDashboardFieldAccessRuleRelatedID: schema.StringAttribute{
									Optional:    true,
									Description: CustomDashboardDescAccessRuleRelatedID,
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 64),
									},
								},
								CustomDashboardFieldAccessRuleRelationType: schema.StringAttribute{
									Required:    true,
									Description: CustomDashboardDescAccessRuleRelationType,
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedRelationTypes.ToStringSlice()...),
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 1,
		},
	}
}

// ============================================================================
// Resource Implementation
// ============================================================================

type customDashboardResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

// MetaData returns the resource metadata
func (r *customDashboardResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

// GetRestResource returns the REST resource for custom dashboards
func (r *customDashboardResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomDashboard] {
	return api.CustomDashboards()
}

// SetComputedFields sets computed fields in the plan (none for this resource)
func (r *customDashboardResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ============================================================================
// API to State Mapping
// ============================================================================

// UpdateState converts API data object to Terraform state
func (r *customDashboardResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, dashboard *restapi.CustomDashboard) diag.Diagnostics {
	var diags diag.Diagnostics

	var model CustomDashboardModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	} else {
		model = CustomDashboardModel{}
	}
	// Create a model and populate it with values from the dashboard
	// model = CustomDashboardModel{
	// 	ID:    types.StringValue(dashboard.ID),
	// 	Title: types.StringValue(dashboard.Title),
	// }
	model.ID = types.StringValue(dashboard.ID)
	model.Title = types.StringValue(dashboard.Title)

	// Handle widgets
	if model.Widgets.IsNull() || model.Widgets.IsUnknown() {

		widgetsBytes, err := dashboard.Widgets.MarshalJSON()
		if err != nil {
			diags.AddError(
				CustomDashboardErrMarshalWidgets,
				fmt.Sprintf(CustomDashboardErrMarshalWidgetsFailed, err),
			)
			return diags
		}
		json, _ := util.CanonicalizeJSON(string(widgetsBytes))
		model.Widgets = jsontypes.NewNormalizedValue(json)
	} 
	// else we keep the existing values

	// Map access rules
	model.AccessRules = r.mapAccessRulesToState(dashboard.AccessRules)

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapAccessRulesToState converts access rules from API format to state models
func (r *customDashboardResourceFramework) mapAccessRulesToState(accessRules []restapi.AccessRule) []AccessRuleModel {
	if len(accessRules) == 0 {
		return nil
	}

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
func (r *customDashboardResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.CustomDashboard, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model CustomDashboardModel

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

	// Map access rules
	accessRules := r.mapAccessRulesFromState(model.AccessRules)

	// Map widgets - normalize the JSON
	var widgets json.RawMessage
	if !model.Widgets.IsNull() {
		normalizedWidgets, _ := util.CanonicalizeJSON(model.Widgets.ValueString())
		widgets = json.RawMessage(normalizedWidgets)
	}

	return &restapi.CustomDashboard{
		ID:          id,
		Title:       model.Title.ValueString(),
		AccessRules: accessRules,
		Widgets:     widgets,
	}, diags
}

// ============================================================================
// Helper Methods
// ============================================================================

// mapAccessRulesFromState converts access rule models from state to API format
func (r *customDashboardResourceFramework) mapAccessRulesFromState(accessRuleModels []AccessRuleModel) []restapi.AccessRule {
	if len(accessRuleModels) == 0 {
		return nil
	}

	accessRules := make([]restapi.AccessRule, len(accessRuleModels))
	for i, ruleModel := range accessRuleModels {
		rule := restapi.AccessRule{
			AccessType:   restapi.AccessType(ruleModel.AccessType.ValueString()),
			RelationType: restapi.RelationType(ruleModel.RelationType.ValueString()),
		}

		// Handle related ID
		if !ruleModel.RelatedID.IsNull() && !utils.IsBlank(ruleModel.RelatedID.ValueString()) {
			relatedID := ruleModel.RelatedID.ValueString()
			rule.RelatedID = &relatedID
		}

		accessRules[i] = rule
	}

	return accessRules
}

type CanonicalizeJSONPlanModifier struct{}

func (m CanonicalizeJSONPlanModifier) Description(ctx context.Context) string {
	return "Canonicalize JSON plan value into deterministic normalized JSON"
}
func (m CanonicalizeJSONPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m CanonicalizeJSONPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Nothing to do if config not set / unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get the string value - works for both types.String and jsontypes.Normalized
	raw := req.ConfigValue.ValueString()

	canon, err := util.CanonicalizeJSON(raw)
	if err != nil {
		resp.Diagnostics.AddError("JSON canonicalization", fmt.Sprintf("failed to canonicalize config JSON: %s", err))
		return
	}

	resp.PlanValue = types.StringValue(canon)
}
