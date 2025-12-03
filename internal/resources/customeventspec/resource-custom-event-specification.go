package customeventspec

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ============================================================================
// Resource Factory
// ============================================================================

// NewCustomEventSpecificationResourceHandle creates the resource handle for Custom Event Specifications
func NewCustomEventSpecificationResourceHandle() resourcehandle.ResourceHandle[*restapi.CustomEventSpecification] {
	return &customEventSpecificationResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaCustomEventSpecification,
			Schema:        createCustomEventSpecificationSchema(),
			SchemaVersion: 1,
		},
	}
}

// ============================================================================
// Resource Implementation
// ============================================================================

type customEventSpecificationResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *customEventSpecificationResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for custom event specifications
func (r *customEventSpecificationResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

// ============================================================================
// Schema Definition
// ============================================================================

// createCustomEventSpecificationSchema creates the schema for the custom event specification resource
func createCustomEventSpecificationSchema() schema.Schema {
	return schema.Schema{
		Description: CustomEventSpecificationResourceDescResource,
		Attributes: map[string]schema.Attribute{
			CustomEventSpecificationFieldID: schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescID,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			CustomEventSpecificationFieldName: schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescName,
				Required:    true,
			},
			CustomEventSpecificationFieldEntityType: schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescEntityType,
				Required:    true,
			},
			CustomEventSpecificationFieldQuery: schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescQuery,
				Optional:    true,
			},
			CustomEventSpecificationFieldTriggering: schema.BoolAttribute{
				Description: CustomEventSpecificationResourceDescTriggering,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			CustomEventSpecificationFieldDescription: schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescDescription,
				Optional:    true,
			},
			CustomEventSpecificationFieldExpirationTime: schema.Int64Attribute{
				Description: CustomEventSpecificationResourceDescExpirationTime,
				Optional:    true,
			},
			CustomEventSpecificationFieldEnabled: schema.BoolAttribute{
				Description: CustomEventSpecificationResourceDescEnabled,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			CustomEventSpecificationFieldRuleLogicalOperator: schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescRuleLogicalOperator,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(CustomEventSpecificationLogicalOperatorAnd),
				Validators: []validator.String{
					stringvalidator.OneOf(CustomEventSpecificationLogicalOperatorAnd, CustomEventSpecificationLogicalOperatorOr),
				},
			},
			CustomEventSpecificationFieldRules: schema.SingleNestedAttribute{
				Description: CustomEventSpecificationResourceDescRules,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					CustomEventSpecificationFieldEntityCountRule: schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescEntityCountRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							CustomEventSpecificationRuleFieldSeverity: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(CustomEventSpecificationSeverityWarning, CustomEventSpecificationSeverityCritical),
								},
							},
							CustomEventSpecificationRuleFieldConditionOperator: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescConditionOperator,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldConditionValue: schema.Float64Attribute{
								Description: CustomEventSpecificationResourceDescConditionValue,
								Required:    true,
							},
						},
					},
					CustomEventSpecificationFieldEntityCountVerificationRule: schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescEntityCountVerification,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							CustomEventSpecificationRuleFieldSeverity: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(CustomEventSpecificationSeverityWarning, CustomEventSpecificationSeverityCritical),
								},
							},
							CustomEventSpecificationRuleFieldConditionOperator: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescConditionOperator,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldConditionValue: schema.Float64Attribute{
								Description: CustomEventSpecificationResourceDescConditionValue,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldMatchingEntityType: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityType,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldMatchingOperator: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingOperator,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldMatchingEntityLabel: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityLabel,
								Required:    true,
							},
						},
					},
					CustomEventSpecificationFieldEntityVerificationRule: schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescEntityVerification,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							CustomEventSpecificationRuleFieldSeverity: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(CustomEventSpecificationSeverityWarning, CustomEventSpecificationSeverityCritical),
								},
							},
							CustomEventSpecificationRuleFieldMatchingEntityType: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityType,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldMatchingOperator: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingOperator,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldMatchingEntityLabel: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityLabel,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldOfflineDuration: schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescOfflineDuration,
								Required:    true,
							},
						},
					},
					CustomEventSpecificationFieldHostAvailabilityRule: schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescHostAvailability,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							CustomEventSpecificationRuleFieldSeverity: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(CustomEventSpecificationSeverityWarning, CustomEventSpecificationSeverityCritical),
								},
							},
							CustomEventSpecificationRuleFieldOfflineDuration: schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescOfflineDuration,
								Required:    true,
							},
							CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter: schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescCloseAfter,
								Optional:    true,
							},
							CustomEventSpecificationHostAvailabilityRuleFieldTagFilter: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescTagFilter,
								Optional:    true,
							},
						},
					},
					CustomEventSpecificationFieldSystemRule: schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescSystemRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							CustomEventSpecificationRuleFieldSeverity: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(CustomEventSpecificationSeverityWarning, CustomEventSpecificationSeverityCritical),
								},
							},
							CustomEventSpecificationSystemRuleFieldSystemRuleId: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSystemRuleID,
								Required:    true,
							},
						},
					},
					CustomEventSpecificationFieldThresholdRule: schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescThresholdRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							CustomEventSpecificationRuleFieldSeverity: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(CustomEventSpecificationSeverityWarning, CustomEventSpecificationSeverityCritical),
								},
							},
							CustomEventSpecificationThresholdRuleFieldMetricName: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMetricName,
								Required:    true,
							},
							CustomEventSpecificationThresholdRuleFieldRollup: schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescRollup,
								Required:    true,
							},
							CustomEventSpecificationThresholdRuleFieldWindow: schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescWindow,
								Required:    true,
							},
							CustomEventSpecificationThresholdRuleFieldAggregation: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescAggregation,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										CustomEventSpecificationAggregationSum,
										CustomEventSpecificationAggregationAvg,
										CustomEventSpecificationAggregationMin,
										CustomEventSpecificationAggregationMax,
										CustomEventSpecificationAggregationAbsDiff,
										CustomEventSpecificationAggregationRelDiff,
									),
								},
							},
							CustomEventSpecificationRuleFieldConditionOperator: schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescConditionOperator,
								Required:    true,
							},
							CustomEventSpecificationRuleFieldConditionValue: schema.Float64Attribute{
								Description: CustomEventSpecificationResourceDescConditionValue,
								Optional:    true,
							},
							CustomEventSpecificationThresholdRuleFieldMetricPattern: schema.SingleNestedAttribute{
								Description: CustomEventSpecificationResourceDescMetricPattern,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix: schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternPrefix,
										Required:    true,
									},
									CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix: schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternPostfix,
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString(CustomEventSpecificationDefaultEmptyString),
									},
									CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder: schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternPlaceholder,
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString(CustomEventSpecificationDefaultEmptyString),
									},
									CustomEventSpecificationThresholdRuleFieldMetricPatternOperator: schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternOperator,
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString(CustomEventSpecificationMetricPatternOperatorEquals),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// ============================================================================
// API to State Mapping
// ============================================================================

// SetComputedFields sets computed fields in the plan (none for this resource)
func (r *customEventSpecificationResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState converts API data object to Terraform state
func (r *customEventSpecificationResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, spec *restapi.CustomEventSpecification) diag.Diagnostics {
	var diags diag.Diagnostics
	var model CustomEventSpecificationModel
	// Validate input
	if spec == nil {
		diags.AddError("Invalid Input", "CustomEventSpecification cannot be nil")
		return diags
	}

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	model.ID = types.StringValue(spec.ID)
	model.Name = types.StringValue(spec.Name)
	model.EntityType = types.StringValue(spec.EntityType)
	model.Triggering = types.BoolValue(spec.Triggering)
	model.Enabled = types.BoolValue(spec.Enabled)
	model.RuleLogicalOperator = types.StringValue(spec.RuleLogicalOperator)
	model.Query = util.SetStringPointerToState(spec.Query)
	model.Description = util.SetStringPointerToState(spec.Description)
	model.ExpirationTime = util.SetInt64PointerToState(spec.ExpirationTime)

	// Process rules if present (preserve the value from plan/model to handle the value drift)
	if model.Rules == nil {
		model.Rules = r.buildRulesModel(spec.Rules, &diags)
	}

	// Set the model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// buildRulesModel processes all rules and returns a RulesModel
func (r *customEventSpecificationResource) buildRulesModel(rules []restapi.RuleSpecification, diags *diag.Diagnostics) *RulesModel {
	rulesModel := &RulesModel{}

	// Use a map for O(1) lookup to process only the first occurrence of each rule type
	processedTypes := make(map[string]bool, 6)

	for i := range rules {
		rule := &rules[i]

		// Skip if this rule type was already processed
		if processedTypes[rule.DType] {
			continue
		}

		switch rule.DType {
		case restapi.EntityCountRuleType:
			if r.isValidEntityCountRule(rule) {
				rulesModel.EntityCount = r.buildEntityCountRule(rule)
				processedTypes[rule.DType] = true
			}
		case restapi.EntityCountVerificationRuleType:
			if r.isValidEntityCountVerificationRule(rule) {
				rulesModel.EntityCountVerification = r.buildEntityCountVerificationRule(rule)
				processedTypes[rule.DType] = true
			}
		case restapi.EntityVerificationRuleType:
			if r.isValidEntityVerificationRule(rule) {
				rulesModel.EntityVerification = r.buildEntityVerificationRule(rule)
				processedTypes[rule.DType] = true
			}
		case restapi.HostAvailabilityRuleType:
			if r.isValidHostAvailabilityRule(rule) {
				rulesModel.HostAvailability = r.buildHostAvailabilityRule(rule, diags)
				processedTypes[rule.DType] = true
			}
		case restapi.SystemRuleType:
			if r.isValidSystemRule(rule) {
				rulesModel.System = r.buildSystemRule(rule)
				processedTypes[rule.DType] = true
			}
		case restapi.ThresholdRuleType:
			if r.isValidThresholdRule(rule) {
				rulesModel.Threshold = r.buildThresholdRule(rule)
				processedTypes[rule.DType] = true
			}
		}
	}

	return rulesModel
}

// Validation methods for each rule type
func (r *customEventSpecificationResource) isValidEntityCountRule(rule *restapi.RuleSpecification) bool {
	return rule.ConditionOperator != nil && rule.ConditionValue != nil
}

func (r *customEventSpecificationResource) isValidEntityCountVerificationRule(rule *restapi.RuleSpecification) bool {
	return rule.ConditionOperator != nil && rule.ConditionValue != nil &&
		rule.MatchingEntityType != nil && rule.MatchingOperator != nil && rule.MatchingEntityLabel != nil
}

func (r *customEventSpecificationResource) isValidEntityVerificationRule(rule *restapi.RuleSpecification) bool {
	return rule.MatchingEntityType != nil && rule.MatchingOperator != nil &&
		rule.MatchingEntityLabel != nil && rule.OfflineDuration != nil
}

func (r *customEventSpecificationResource) isValidHostAvailabilityRule(rule *restapi.RuleSpecification) bool {
	return rule.OfflineDuration != nil
}

func (r *customEventSpecificationResource) isValidSystemRule(rule *restapi.RuleSpecification) bool {
	return rule.SystemRuleID != nil
}

func (r *customEventSpecificationResource) isValidThresholdRule(rule *restapi.RuleSpecification) bool {
	return rule.MetricName != nil && rule.Rollup != nil && rule.Window != nil &&
		rule.Aggregation != nil && rule.ConditionOperator != nil && rule.ConditionValue != nil
}

// Builder methods for each rule type
func (r *customEventSpecificationResource) buildEntityCountRule(rule *restapi.RuleSpecification) *EntityCountRuleModel {
	return &EntityCountRuleModel{
		Severity:          mapIntToSeverityString(rule.Severity),
		ConditionOperator: util.SetStringPointerToState(rule.ConditionOperator),
		ConditionValue:    util.SetFloat64PointerToState(rule.ConditionValue),
	}
}

func (r *customEventSpecificationResource) buildEntityCountVerificationRule(rule *restapi.RuleSpecification) *EntityCountVerificationRuleModel {
	return &EntityCountVerificationRuleModel{
		Severity:            mapIntToSeverityString(rule.Severity),
		ConditionOperator:   util.SetStringPointerToState(rule.ConditionOperator),
		ConditionValue:      util.SetFloat64PointerToState(rule.ConditionValue),
		MatchingEntityType:  util.SetStringPointerToState(rule.MatchingEntityType),
		MatchingOperator:    util.SetStringPointerToState(rule.MatchingOperator),
		MatchingEntityLabel: util.SetStringPointerToState(rule.MatchingEntityLabel),
	}
}

func (r *customEventSpecificationResource) buildEntityVerificationRule(rule *restapi.RuleSpecification) *EntityVerificationRuleModel {
	return &EntityVerificationRuleModel{
		Severity:            mapIntToSeverityString(rule.Severity),
		MatchingEntityType:  util.SetStringPointerToState(rule.MatchingEntityType),
		MatchingOperator:    util.SetStringPointerToState(rule.MatchingOperator),
		MatchingEntityLabel: util.SetStringPointerToState(rule.MatchingEntityLabel),
		OfflineDuration:     util.SetInt64PointerToState(rule.OfflineDuration),
	}
}

func (r *customEventSpecificationResource) buildHostAvailabilityRule(rule *restapi.RuleSpecification, diags *diag.Diagnostics) *HostAvailabilityRuleModel {
	model := &HostAvailabilityRuleModel{
		Severity:        mapIntToSeverityString(rule.Severity),
		OfflineDuration: util.SetInt64PointerToState(rule.OfflineDuration),
		CloseAfter:      util.SetInt64PointerToState(rule.CloseAfter),
		TagFilter:       types.StringValue(CustomEventSpecificationDefaultEmptyString),
	}

	// Handle tag filter conversion with proper error handling
	if rule.TagFilter != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(rule.TagFilter)
		if err != nil {
			diags.AddWarning(
				"Tag Filter Conversion Warning",
				fmt.Sprintf("Failed to convert tag filter to string: %v. Using empty string.", err),
			)
		} else if normalizedTagFilterString != nil {
			model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
		}
	}

	return model
}

func (r *customEventSpecificationResource) buildSystemRule(rule *restapi.RuleSpecification) *SystemRuleModel {
	return &SystemRuleModel{
		Severity:     mapIntToSeverityString(rule.Severity),
		SystemRuleID: util.SetStringPointerToState(rule.SystemRuleID),
	}
}

func (r *customEventSpecificationResource) buildThresholdRule(rule *restapi.RuleSpecification) *ThresholdRuleModel {
	model := &ThresholdRuleModel{
		Severity:          mapIntToSeverityString(rule.Severity),
		MetricName:        util.SetStringPointerToState(rule.MetricName),
		Rollup:            util.SetInt64PointerToState(rule.Rollup),
		Window:            util.SetInt64PointerToState(rule.Window),
		Aggregation:       util.SetStringPointerToState(rule.Aggregation),
		ConditionOperator: util.SetStringPointerToState(rule.ConditionOperator),
		ConditionValue:    util.SetFloat64PointerToState(rule.ConditionValue),
	}

	// Handle metric pattern if present
	if rule.MetricPattern != nil {
		model.MetricPattern = &MetricPatternModel{
			Prefix:      types.StringValue(rule.MetricPattern.Prefix),
			Operator:    types.StringValue(rule.MetricPattern.Operator),
			Postfix:     util.SetStringPointerToState(rule.MetricPattern.Postfix),
			Placeholder: util.SetStringPointerToState(rule.MetricPattern.Placeholder),
		}
	}

	return model
}

// ============================================================================
// Helper Methods
// ============================================================================

// mapIntToSeverityString maps the severity integer to a string value
func mapIntToSeverityString(severity int) types.String {
	switch severity {
	case 5:
		return types.StringValue(CustomEventSpecificationSeverityWarning)
	case 10:
		return types.StringValue(CustomEventSpecificationSeverityCritical)
	default:
		return types.StringValue(CustomEventSpecificationSeverityWarning) // Default to warning
	}
}

// ============================================================================
// State to API Mapping
// ============================================================================

// MapStateToDataObject converts Terraform state to API data object
func (r *customEventSpecificationResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.CustomEventSpecification, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract model from plan or state
	model, extractDiags := r.extractModel(ctx, plan, state)
	diags.Append(extractDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Build API specification from model
	spec := r.buildAPISpecification(model, &diags)
	if diags.HasError() {
		return nil, diags
	}

	return spec, diags
}

// extractModel retrieves the model from plan or state with proper validation
func (r *customEventSpecificationResource) extractModel(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*CustomEventSpecificationModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model CustomEventSpecificationModel

	// Validate input - at least one source must be provided
	if plan == nil && state == nil {
		diags.AddError(
			"Invalid Input",
			"Both plan and state are nil. At least one must be provided.",
		)
		return nil, diags
	}

	// Get model from plan (preferred) or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &model, diags
}

// buildAPISpecification constructs the API specification from the model
func (r *customEventSpecificationResource) buildAPISpecification(model *CustomEventSpecificationModel, diags *diag.Diagnostics) *restapi.CustomEventSpecification {
	return &restapi.CustomEventSpecification{
		ID:                  r.extractID(model),
		Name:                model.Name.ValueString(),
		EntityType:          model.EntityType.ValueString(),
		Query:               r.extractOptionalString(model.Query),
		Triggering:          model.Triggering.ValueBool(),
		Description:         r.extractOptionalString(model.Description),
		ExpirationTime:      r.extractOptionalInt(model.ExpirationTime),
		Enabled:             model.Enabled.ValueBool(),
		RuleLogicalOperator: model.RuleLogicalOperator.ValueString(),
		Rules:               r.buildRulesFromModel(model.Rules, diags),
	}
}

// extractID extracts the ID from the model, defaulting to empty string
func (r *customEventSpecificationResource) extractID(model *CustomEventSpecificationModel) string {
	if model.ID.IsNull() {
		return CustomEventSpecificationDefaultEmptyString
	}
	return model.ID.ValueString()
}

// extractOptionalString converts a types.String to *string, handling null and empty values
func (r *customEventSpecificationResource) extractOptionalString(value types.String) *string {
	if value.IsNull() || value.ValueString() == CustomEventSpecificationDefaultEmptyString {
		return nil
	}
	str := value.ValueString()
	return &str
}

// extractOptionalInt converts a types.Int64 to *int, handling null values
func (r *customEventSpecificationResource) extractOptionalInt(value types.Int64) *int {
	if value.IsNull() {
		return nil
	}
	intVal := int(value.ValueInt64())
	return &intVal
}

// buildRulesFromModel converts model rules to API rule specifications
func (r *customEventSpecificationResource) buildRulesFromModel(rulesModel *RulesModel, diags *diag.Diagnostics) []restapi.RuleSpecification {
	if rulesModel == nil {
		return nil
	}

	// Pre-allocate slice with estimated capacity to reduce allocations
	rules := make([]restapi.RuleSpecification, 0, 6)

	// Process each rule type using dedicated builder methods
	if rulesModel.EntityCount != nil {
		rules = append(rules, r.buildEntityCountRuleSpec(rulesModel.EntityCount))
	}

	if rulesModel.EntityCountVerification != nil {
		rules = append(rules, r.buildEntityCountVerificationRuleSpec(rulesModel.EntityCountVerification))
	}

	if rulesModel.EntityVerification != nil {
		rules = append(rules, r.buildEntityVerificationRuleSpec(rulesModel.EntityVerification))
	}

	if rulesModel.HostAvailability != nil {
		if spec, ok := r.buildHostAvailabilityRuleSpec(rulesModel.HostAvailability, diags); ok {
			rules = append(rules, spec)
		}
	}

	if rulesModel.System != nil {
		rules = append(rules, r.buildSystemRuleSpec(rulesModel.System))
	}

	if rulesModel.Threshold != nil {
		rules = append(rules, r.buildThresholdRuleSpec(rulesModel.Threshold))
	}

	return rules
}

// buildEntityCountRuleSpec creates an entity count rule specification
func (r *customEventSpecificationResource) buildEntityCountRuleSpec(rule *EntityCountRuleModel) restapi.RuleSpecification {
	conditionOperator := rule.ConditionOperator.ValueString()
	conditionValue := rule.ConditionValue.ValueFloat64()

	return restapi.RuleSpecification{
		DType:             restapi.EntityCountRuleType,
		Severity:          mapSeverityToInt(rule.Severity.ValueString()),
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}
}

// buildEntityCountVerificationRuleSpec creates an entity count verification rule specification
func (r *customEventSpecificationResource) buildEntityCountVerificationRuleSpec(rule *EntityCountVerificationRuleModel) restapi.RuleSpecification {
	conditionOperator := rule.ConditionOperator.ValueString()
	conditionValue := rule.ConditionValue.ValueFloat64()
	matchingEntityType := rule.MatchingEntityType.ValueString()
	matchingOperator := rule.MatchingOperator.ValueString()
	matchingEntityLabel := rule.MatchingEntityLabel.ValueString()

	return restapi.RuleSpecification{
		DType:               restapi.EntityCountVerificationRuleType,
		Severity:            mapSeverityToInt(rule.Severity.ValueString()),
		ConditionOperator:   &conditionOperator,
		ConditionValue:      &conditionValue,
		MatchingEntityType:  &matchingEntityType,
		MatchingOperator:    &matchingOperator,
		MatchingEntityLabel: &matchingEntityLabel,
	}
}

// buildEntityVerificationRuleSpec creates an entity verification rule specification
func (r *customEventSpecificationResource) buildEntityVerificationRuleSpec(rule *EntityVerificationRuleModel) restapi.RuleSpecification {
	matchingEntityType := rule.MatchingEntityType.ValueString()
	matchingOperator := rule.MatchingOperator.ValueString()
	matchingEntityLabel := rule.MatchingEntityLabel.ValueString()
	offlineDuration := int(rule.OfflineDuration.ValueInt64())

	return restapi.RuleSpecification{
		DType:               restapi.EntityVerificationRuleType,
		Severity:            mapSeverityToInt(rule.Severity.ValueString()),
		MatchingEntityType:  &matchingEntityType,
		MatchingOperator:    &matchingOperator,
		MatchingEntityLabel: &matchingEntityLabel,
		OfflineDuration:     &offlineDuration,
	}
}

// buildHostAvailabilityRuleSpec creates a host availability rule specification
// Returns the specification and a boolean indicating success
func (r *customEventSpecificationResource) buildHostAvailabilityRuleSpec(rule *HostAvailabilityRuleModel, diags *diag.Diagnostics) (restapi.RuleSpecification, bool) {
	offlineDuration := int(rule.OfflineDuration.ValueInt64())

	spec := restapi.RuleSpecification{
		DType:           restapi.HostAvailabilityRuleType,
		Severity:        mapSeverityToInt(rule.Severity.ValueString()),
		OfflineDuration: &offlineDuration,
		CloseAfter:      r.extractOptionalInt(rule.CloseAfter),
	}

	// Parse tag filter if provided
	if !rule.TagFilter.IsNull() && rule.TagFilter.ValueString() != CustomEventSpecificationDefaultEmptyString {
		tagFilter, err := r.parseTagFilter(rule.TagFilter.ValueString())
		if err != nil {
			diags.AddError(
				CustomEventSpecificationResourceErrParseTagFilter,
				fmt.Sprintf(CustomEventSpecificationResourceErrParseTagFilterMsg, err),
			)
			return restapi.RuleSpecification{}, false
		}
		spec.TagFilter = tagFilter
	}

	return spec, true
}

// parseTagFilter parses a tag filter string into an API model
func (r *customEventSpecificationResource) parseTagFilter(tagFilterStr string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(tagFilterStr)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

// buildSystemRuleSpec creates a system rule specification
func (r *customEventSpecificationResource) buildSystemRuleSpec(rule *SystemRuleModel) restapi.RuleSpecification {
	systemRuleID := rule.SystemRuleID.ValueString()

	return restapi.RuleSpecification{
		DType:        restapi.SystemRuleType,
		Severity:     mapSeverityToInt(rule.Severity.ValueString()),
		SystemRuleID: &systemRuleID,
	}
}

// buildThresholdRuleSpec creates a threshold rule specification
func (r *customEventSpecificationResource) buildThresholdRuleSpec(rule *ThresholdRuleModel) restapi.RuleSpecification {
	metricName := rule.MetricName.ValueString()
	rollup := int(rule.Rollup.ValueInt64())
	window := int(rule.Window.ValueInt64())
	aggregation := rule.Aggregation.ValueString()
	conditionOperator := rule.ConditionOperator.ValueString()
	conditionValue := rule.ConditionValue.ValueFloat64()

	return restapi.RuleSpecification{
		DType:             restapi.ThresholdRuleType,
		Severity:          mapSeverityToInt(rule.Severity.ValueString()),
		MetricName:        &metricName,
		Rollup:            &rollup,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
		MetricPattern:     r.buildMetricPattern(rule.MetricPattern),
	}
}

// buildMetricPattern creates a metric pattern from the model
func (r *customEventSpecificationResource) buildMetricPattern(pattern *MetricPatternModel) *restapi.MetricPattern {
	if pattern == nil {
		return nil
	}

	return &restapi.MetricPattern{
		Prefix:      pattern.Prefix.ValueString(),
		Operator:    pattern.Operator.ValueString(),
		Postfix:     r.extractOptionalString(pattern.Postfix),
		Placeholder: r.extractOptionalString(pattern.Placeholder),
	}
}

// mapSeverityToInt maps the severity string to an integer value
func mapSeverityToInt(severity string) int {
	switch severity {
	case CustomEventSpecificationSeverityWarning:
		return 5
	case CustomEventSpecificationSeverityCritical:
		return 10
	default:
		return 5 // Default to warning
	}
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *customEventSpecificationResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: resourcehandle.CreateStateUpgraderForVersion(0),
	}
}
