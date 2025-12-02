package websitealertconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewWebsiteAlertConfigResourceHandle creates the resource handle for Website Alert Configs
func NewWebsiteAlertConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.WebsiteAlertConfig] {
	return &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaWebsiteAlertConfig,
			Schema: schema.Schema{
				Description: WebsiteAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					WebsiteAlertConfigFieldID: schema.StringAttribute{
						Computed:    true,
						Description: WebsiteAlertConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					WebsiteAlertConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: WebsiteAlertConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(WebsiteAlertConfigMinNameLength, WebsiteAlertConfigMaxNameLength),
						},
					},
					WebsiteAlertConfigFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: WebsiteAlertConfigDescDescription,
						Validators: []validator.String{
							stringvalidator.LengthBetween(WebsiteAlertConfigMinDescriptionLength, WebsiteAlertConfigMaxDescriptionLength),
						},
					},
					WebsiteAlertConfigFieldTriggering: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: WebsiteAlertConfigDescTriggering,
						Default:     booldefault.StaticBool(WebsiteAlertConfigDefaultTriggering),
					},
					WebsiteAlertConfigFieldWebsiteID: schema.StringAttribute{
						Required:    true,
						Description: WebsiteAlertConfigDescWebsiteID,
						Validators: []validator.String{
							stringvalidator.LengthBetween(WebsiteAlertConfigMinWebsiteIDLength, WebsiteAlertConfigMaxWebsiteIDLength),
						},
					},
					WebsiteAlertConfigFieldTagFilter: schema.StringAttribute{
						Optional:    true,
						Description: WebsiteAlertConfigDescTagFilter,
					},
					WebsiteAlertConfigFieldAlertChannelIDs: schema.SetAttribute{
						Optional:    true,
						Description: WebsiteAlertConfigDescAlertChannelIDs,
						ElementType: types.StringType,
					},
					WebsiteAlertConfigFieldGranularity: schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: WebsiteAlertConfigDescGranularity,
						Default:     int64default.StaticInt64(WebsiteAlertConfigDefaultGranularity),
					},
					WebsiteAlertConfigFieldRules: schema.ListNestedAttribute{
						Description: WebsiteAlertConfigDescRules,
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								WebsiteAlertConfigFieldRule: schema.SingleNestedAttribute{
									Description: WebsiteAlertConfigDescRule,
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										WebsiteAlertConfigFieldRuleSlowness: schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleSlowness,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												WebsiteAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												WebsiteAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf(WebsiteAlertConfigAggregationSUM, WebsiteAlertConfigAggregationMEAN, WebsiteAlertConfigAggregationMAX, WebsiteAlertConfigAggregationMIN, WebsiteAlertConfigAggregationP25, WebsiteAlertConfigAggregationP50, WebsiteAlertConfigAggregationP75, WebsiteAlertConfigAggregationP90, WebsiteAlertConfigAggregationP95, WebsiteAlertConfigAggregationP98, WebsiteAlertConfigAggregationP99),
													},
												},
											},
										},
										WebsiteAlertConfigFieldRuleSpecificJsError: schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleSpecificJsError,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												WebsiteAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												WebsiteAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf(WebsiteAlertConfigAggregationSUM, WebsiteAlertConfigAggregationMEAN, WebsiteAlertConfigAggregationMAX, WebsiteAlertConfigAggregationMIN, WebsiteAlertConfigAggregationP25, WebsiteAlertConfigAggregationP50, WebsiteAlertConfigAggregationP75, WebsiteAlertConfigAggregationP90, WebsiteAlertConfigAggregationP95, WebsiteAlertConfigAggregationP98, WebsiteAlertConfigAggregationP99),
													},
												},
												WebsiteAlertConfigFieldRuleValue: schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleValueJsError,
												},
												WebsiteAlertConfigFieldRuleOperator: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: WebsiteAlertConfigDescRuleOperator,
													Validators: []validator.String{
														stringvalidator.OneOf(shared.SupportedOparators...),
													},
												},
											},
										},
										WebsiteAlertConfigFieldRuleStatusCode: schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleStatusCode,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												WebsiteAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												WebsiteAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf(WebsiteAlertConfigAggregationSUM, WebsiteAlertConfigAggregationMEAN, WebsiteAlertConfigAggregationMAX, WebsiteAlertConfigAggregationMIN, WebsiteAlertConfigAggregationP25, WebsiteAlertConfigAggregationP50, WebsiteAlertConfigAggregationP75, WebsiteAlertConfigAggregationP90, WebsiteAlertConfigAggregationP95, WebsiteAlertConfigAggregationP98, WebsiteAlertConfigAggregationP99),
													},
												},
												WebsiteAlertConfigFieldRuleOperator: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleOperatorEval,
													Validators: []validator.String{
														stringvalidator.OneOf(shared.SupportedOparators...),
													},
												},
												WebsiteAlertConfigFieldRuleValue: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleValueStatusCode,
												},
											},
										},
										WebsiteAlertConfigFieldRuleThroughput: schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleThroughput,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												WebsiteAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												WebsiteAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf(WebsiteAlertConfigAggregationSUM, WebsiteAlertConfigAggregationMEAN, WebsiteAlertConfigAggregationMAX, WebsiteAlertConfigAggregationMIN, WebsiteAlertConfigAggregationP25, WebsiteAlertConfigAggregationP50, WebsiteAlertConfigAggregationP75, WebsiteAlertConfigAggregationP90, WebsiteAlertConfigAggregationP95, WebsiteAlertConfigAggregationP98, WebsiteAlertConfigAggregationP99),
													},
												},
											},
										},
									},
								},
								WebsiteAlertConfigFieldRuleOperator: schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The operator to apply for threshold comparison",
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "<", "<="),
									},
								},
								WebsiteAlertConfigFieldThreshold: schema.SingleNestedAttribute{
									Description: WebsiteAlertConfigDescThreshold,
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										shared.LogAlertConfigFieldWarning:  shared.AllThresholdAttributeSchema(),
										shared.LogAlertConfigFieldCritical: shared.AllThresholdAttributeSchema(),
									},
								},
							},
						},
					},
					WebsiteAlertConfigFieldCustomPayloadFields: shared.GetCustomPayloadFieldsSchema(),
					WebsiteAlertConfigFieldTimeThreshold: schema.SingleNestedAttribute{
						Description: WebsiteAlertConfigDescTimeThreshold,
						Required:    true,
						Attributes: map[string]schema.Attribute{
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: schema.SingleNestedAttribute{
								Description: WebsiteAlertConfigDescTimeThresholdUserImpact,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									WebsiteAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdTimeWindow,
									},
									WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod: schema.StringAttribute{
										Required:    true,
										Description: WebsiteAlertConfigDescTimeThresholdImpactMethod,
										Validators: []validator.String{
											stringvalidator.OneOf(WebsiteAlertConfigImpactMethodAggregated, WebsiteAlertConfigImpactMethodPerWindow),
										},
									},
									WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage: schema.Float64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdUserPercentage,
									},
									WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdUsers,
									},
								},
							},
							WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod: schema.SingleNestedAttribute{
								Description: WebsiteAlertConfigDescTimeThresholdViolationsInPeriod,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									WebsiteAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdTimeWindow,
									},
									WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations: schema.Int64Attribute{
										Required:    true,
										Description: WebsiteAlertConfigDescTimeThresholdViolations,
									},
								},
							},
							WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: schema.SingleNestedAttribute{
								Description: WebsiteAlertConfigDescTimeThresholdViolationsInSequence,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									WebsiteAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdTimeWindow,
										Default:     int64default.StaticInt64(600000),
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

type websiteAlertConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *websiteAlertConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *websiteAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return api.WebsiteAlertConfig()
}

func (r *websiteAlertConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *websiteAlertConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.WebsiteAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model WebsiteAlertConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map tag filter
	tagFilter, tagFilterDiags := r.mapTagFilterToAPI(model)
	diags.Append(tagFilterDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map alert channel IDs
	alertChannelIDs, channelDiags := r.mapAlertChannelIDsToAPI(ctx, model)
	diags.Append(channelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map custom payload fields
	customPayloadFields, payloadDiags := r.mapCustomPayloadFieldsToAPI(ctx, model)
	diags.Append(payloadDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map time threshold
	timeThreshold, timeThresholdDiags := r.mapTimeThresholdFromModel(ctx, model)
	diags.Append(timeThresholdDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map rules collection
	rules, rulesDiags := r.mapRulesCollectionToAPI(ctx, model)
	diags.Append(rulesDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Create API object
	return &restapi.WebsiteAlertConfig{
		ID:                    model.ID.ValueString(),
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		Severity:              nil,
		Triggering:            model.Triggering.ValueBool(),
		WebsiteID:             model.WebsiteID.ValueString(),
		TagFilterExpression:   tagFilter,
		AlertChannelIDs:       alertChannelIDs,
		Granularity:           restapi.Granularity(model.Granularity.ValueInt64()),
		CustomerPayloadFields: customPayloadFields,
		TimeThreshold:         *timeThreshold,
		Rules:                 rules,
	}, diags
}

// mapTagFilterToAPI maps tag filter expression from model to API
func (r *websiteAlertConfigResource) mapTagFilterToAPI(model WebsiteAlertConfigModel) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError(WebsiteAlertConfigErrParseTagFilter, err.Error())
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		return mapper.ToAPIModel(expr), diags
	}

	// Return default tag filter
	operator := restapi.LogicalOperatorType(WebsiteAlertConfigLogicalOperatorAND)
	return &restapi.TagFilter{
		Type:            WebsiteAlertConfigTagFilterTypeExpression,
		LogicalOperator: &operator,
		Elements:        []*restapi.TagFilter{},
	}, diags
}

// mapAlertChannelIDsToAPI maps alert channel IDs from model to API
func (r *websiteAlertConfigResource) mapAlertChannelIDsToAPI(ctx context.Context, model WebsiteAlertConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var alertChannelIDs []string

	if !model.AlertChannelIDs.IsNull() && !model.AlertChannelIDs.IsUnknown() {
		diags.Append(model.AlertChannelIDs.ElementsAs(ctx, &alertChannelIDs, false)...)
	}

	return alertChannelIDs, diags
}

// mapCustomPayloadFieldsToAPI maps custom payload fields from model to API
func (r *websiteAlertConfigResource) mapCustomPayloadFieldsToAPI(ctx context.Context, model WebsiteAlertConfigModel) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics
	customPayloadFields := make([]restapi.CustomPayloadField[any], 0)

	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		var payloadDiags diag.Diagnostics
		customPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		diags.Append(payloadDiags...)
	}

	return customPayloadFields, diags
}

// mapRulesCollectionToAPI maps the rules collection from model to API
func (r *websiteAlertConfigResource) mapRulesCollectionToAPI(ctx context.Context, model WebsiteAlertConfigModel) ([]restapi.WebsiteAlertRuleWithThresholds, diag.Diagnostics) {
	var diags diag.Diagnostics
	rules := make([]restapi.WebsiteAlertRuleWithThresholds, 0)

	if len(model.Rules) == 0 {
		return rules, diags
	}

	// Process each rule
	for _, ruleModel := range model.Rules {
		websiteAlertRule, ruleDiags := r.mapSingleRuleFromCollection(ctx, ruleModel)
		diags.Append(ruleDiags...)
		if diags.HasError() {
			return nil, diags
		}

		// Only process if we have a valid rule
		if websiteAlertRule != nil {
			thresholdMap, thresholdDiags := shared.MapThresholdsAllPluginFromState(ctx, ruleModel.Thresholds)
			diags.Append(thresholdDiags...)
			if diags.HasError() {
				return nil, diags
			}

			rules = append(rules, restapi.WebsiteAlertRuleWithThresholds{
				Rule:              websiteAlertRule,
				ThresholdOperator: ruleModel.ThresholdOperator.ValueString(),
				Thresholds:        thresholdMap,
			})
		}
	}

	return rules, diags
}

// mapSingleRuleFromCollection maps a single rule from the rules collection
func (r *websiteAlertConfigResource) mapSingleRuleFromCollection(ctx context.Context, ruleModel RuleWithThresholdPluginModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleModel.Rule == nil {
		return nil, diags
	}

	// Check each rule type and only process if not null/unknown
	if ruleModel.Rule.Slowness != nil && !ruleModel.Rule.Slowness.MetricName.IsNull() && !ruleModel.Rule.Slowness.MetricName.IsUnknown() {
		return r.mapSlownessRule(*ruleModel.Rule.Slowness), diags
	} else if ruleModel.Rule.SpecificJsError != nil && !ruleModel.Rule.SpecificJsError.MetricName.IsNull() && !ruleModel.Rule.SpecificJsError.MetricName.IsUnknown() {
		return r.mapSpecificJsErrorRule(*ruleModel.Rule.SpecificJsError), diags
	} else if ruleModel.Rule.StatusCode != nil && !ruleModel.Rule.StatusCode.MetricName.IsNull() && !ruleModel.Rule.StatusCode.MetricName.IsUnknown() {
		return r.mapStatusCodeRule(*ruleModel.Rule.StatusCode), diags
	} else if ruleModel.Rule.Throughput != nil && !ruleModel.Rule.Throughput.MetricName.IsNull() && !ruleModel.Rule.Throughput.MetricName.IsUnknown() {
		return r.mapThroughputRule(*ruleModel.Rule.Throughput), diags
	}

	return nil, diags
}

func (r *websiteAlertConfigResource) mapThroughputRule(throughputModel WebsiteAlertRuleConfigModel) *restapi.WebsiteAlertRule {
	var aggregationPtr *restapi.Aggregation
	if !throughputModel.Aggregation.IsNull() && !throughputModel.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(throughputModel.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	return &restapi.WebsiteAlertRule{
		AlertType:   WebsiteAlertConfigAlertTypeThroughput,
		MetricName:  throughputModel.MetricName.ValueString(),
		Aggregation: aggregationPtr,
	}
}

func (r *websiteAlertConfigResource) mapStatusCodeRule(statusCodeModel WebsiteAlertRuleConfigCompleteModel) *restapi.WebsiteAlertRule {

	var aggregationPtr *restapi.Aggregation
	if !statusCodeModel.Aggregation.IsNull() && !statusCodeModel.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(statusCodeModel.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	operator := restapi.ExpressionOperator(statusCodeModel.Operator.ValueString())
	value := statusCodeModel.Value.ValueString()

	return &restapi.WebsiteAlertRule{
		AlertType:   WebsiteAlertConfigAlertTypeStatusCode,
		MetricName:  statusCodeModel.MetricName.ValueString(),
		Aggregation: aggregationPtr,
		Operator:    &operator,
		Value:       &value,
	}
}

func (r *websiteAlertConfigResource) mapSpecificJsErrorRule(specificJsErrorModel WebsiteAlertRuleConfigCompleteModel) *restapi.WebsiteAlertRule {

	var aggregationPtr *restapi.Aggregation
	if !specificJsErrorModel.Aggregation.IsNull() && !specificJsErrorModel.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(specificJsErrorModel.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	operator := restapi.ExpressionOperator(specificJsErrorModel.Operator.ValueString())
	var valuePtr *string
	if !specificJsErrorModel.Value.IsNull() && !specificJsErrorModel.Value.IsUnknown() {
		value := specificJsErrorModel.Value.ValueString()
		valuePtr = &value
	}

	return &restapi.WebsiteAlertRule{
		AlertType:   WebsiteAlertConfigAlertTypeSpecificJsError,
		MetricName:  specificJsErrorModel.MetricName.ValueString(),
		Aggregation: aggregationPtr,
		Operator:    &operator,
		Value:       valuePtr,
	}
}

func (r *websiteAlertConfigResource) mapSlownessRule(slownessModel WebsiteAlertRuleConfigModel) *restapi.WebsiteAlertRule {
	aggregation := restapi.Aggregation(slownessModel.Aggregation.ValueString())

	return &restapi.WebsiteAlertRule{
		AlertType:   WebsiteAlertConfigAlertTypeSlowness,
		MetricName:  slownessModel.MetricName.ValueString(),
		Aggregation: &aggregation,
	}
}

// mapTimeThresholdFromModel converts the time threshold configuration from the Terraform model to the API representation.
// It validates that exactly one time threshold type is configured and returns appropriate diagnostics on error.
func (r *websiteAlertConfigResource) mapTimeThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate time threshold is provided
	if model.TimeThreshold == nil {
		diags.AddError(WebsiteAlertConfigErrTimeThresholdRequired, WebsiteAlertConfigErrTimeThresholdRequiredMsg)
		return nil, diags
	}

	timeThresholdModel := *model.TimeThreshold

	// Map to appropriate time threshold type based on which field is set
	switch {
	case timeThresholdModel.UserImpactOfViolationsInSequence != nil:
		return r.mapUserImpactTimeThreshold(timeThresholdModel.UserImpactOfViolationsInSequence, diags)
	case timeThresholdModel.ViolationsInPeriod != nil:
		return r.mapViolationsInPeriodTimeThreshold(timeThresholdModel.ViolationsInPeriod, diags)
	case timeThresholdModel.ViolationsInSequence != nil:
		return r.mapViolationsInSequenceTimeThreshold(timeThresholdModel.ViolationsInSequence, diags)
	default:
		diags.AddError(WebsiteAlertConfigErrInvalidTimeThresholdConfig, WebsiteAlertConfigErrInvalidTimeThresholdConfigMsg)
		return nil, diags
	}
}

// mapUserImpactTimeThreshold maps the user impact of violations in sequence time threshold configuration.
func (r *websiteAlertConfigResource) mapUserImpactTimeThreshold(
	userImpactModel *WebsiteUserImpactOfViolationsInSequenceModel,
	diags diag.Diagnostics,
) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	threshold := &restapi.WebsiteTimeThreshold{
		Type: WebsiteAlertConfigTimeThresholdTypeUserImpact,
	}

	// Map time window
	threshold.TimeWindow = r.mapOptionalInt64Field(userImpactModel.TimeWindow)

	// Map impact measurement method
	if !userImpactModel.ImpactMeasurementMethod.IsNull() && !userImpactModel.ImpactMeasurementMethod.IsUnknown() {
		method := restapi.WebsiteImpactMeasurementMethod(userImpactModel.ImpactMeasurementMethod.ValueString())
		threshold.ImpactMeasurementMethod = &method
	}

	// Map user percentage
	threshold.UserPercentage = r.mapOptionalFloat64Field(userImpactModel.UserPercentage)

	// Map users count
	threshold.Users = r.mapOptionalInt64ToInt32Field(userImpactModel.Users)

	return threshold, diags
}

// mapViolationsInPeriodTimeThreshold maps the violations in period time threshold configuration.
func (r *websiteAlertConfigResource) mapViolationsInPeriodTimeThreshold(
	violationsInPeriodModel *WebsiteViolationsInPeriodModel,
	diags diag.Diagnostics,
) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	threshold := &restapi.WebsiteTimeThreshold{
		Type: WebsiteAlertConfigTimeThresholdTypeViolationsInPeriod,
	}

	// Map time window
	threshold.TimeWindow = r.mapOptionalInt64Field(violationsInPeriodModel.TimeWindow)

	// Map violations count
	threshold.Violations = r.mapOptionalInt64ToInt32Field(violationsInPeriodModel.Violations)

	return threshold, diags
}

// mapViolationsInSequenceTimeThreshold maps the violations in sequence time threshold configuration.
func (r *websiteAlertConfigResource) mapViolationsInSequenceTimeThreshold(
	violationsInSequenceModel *WebsiteViolationsInSequenceModel,
	diags diag.Diagnostics,
) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	threshold := &restapi.WebsiteTimeThreshold{
		Type: WebsiteAlertConfigTimeThresholdTypeViolationsInSequence,
	}

	// Map time window
	threshold.TimeWindow = r.mapOptionalInt64Field(violationsInSequenceModel.TimeWindow)

	return threshold, diags
}

// mapOptionalInt64Field extracts an optional int64 field value, returning nil if null or unknown.
func (r *websiteAlertConfigResource) mapOptionalInt64Field(field types.Int64) *int64 {
	if field.IsNull() || field.IsUnknown() {
		return nil
	}
	value := field.ValueInt64()
	return &value
}

// mapOptionalFloat64Field extracts an optional float64 field value, returning nil if null or unknown.
func (r *websiteAlertConfigResource) mapOptionalFloat64Field(field types.Float64) *float64 {
	if field.IsNull() || field.IsUnknown() {
		return nil
	}
	value := field.ValueFloat64()
	return &value
}

// mapOptionalInt64ToInt32Field extracts an optional int64 field and converts it to int32, returning nil if null or unknown.
func (r *websiteAlertConfigResource) mapOptionalInt64ToInt32Field(field types.Int64) *int32 {
	if field.IsNull() || field.IsUnknown() {
		return nil
	}
	value := int32(field.ValueInt64())
	return &value
}

func (r *websiteAlertConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.WebsiteAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model WebsiteAlertConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	} else {
		model = WebsiteAlertConfigModel{}
	}
	// Create base model
	model.ID = types.StringValue(apiObject.ID)
	model.Name = types.StringValue(apiObject.Name)
	model.Description = types.StringValue(apiObject.Description)
	model.Triggering = types.BoolValue(apiObject.Triggering)
	model.WebsiteID = types.StringValue(apiObject.WebsiteID)
	model.Granularity = types.Int64Value(int64(apiObject.Granularity))

	// Map tag filter expression
	tagFilterDiags := r.mapTagFilterExpressionToState(&model, apiObject)
	diags.Append(tagFilterDiags...)
	if diags.HasError() {
		return diags
	}

	// Map alert channel IDs
	r.mapAlertChannelIDsToState(&model, apiObject)

	// Map rules collection (preserve the plan values / update only if empty)
	if len(model.Rules) == 0 {
		model.Rules = r.mapRulesToState(ctx, apiObject)
	}

	// Map custom payload fields
	customPayloadDiags := r.mapCustomPayloadFieldsToState(ctx, &model, apiObject)
	diags.Append(customPayloadDiags...)
	if diags.HasError() {
		return diags
	}

	// Map time threshold
	model.TimeThreshold = r.mapTimeThresholdToState(apiObject.TimeThreshold)

	// Set state
	diags.Append(state.Set(ctx, &model)...)
	return diags
}

// mapTagFilterExpressionToState maps tag filter expression from API to state
func (r *websiteAlertConfigResource) mapTagFilterExpressionToState(model *WebsiteAlertConfigModel, apiObject *restapi.WebsiteAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject.TagFilterExpression != nil {
		filterExprStr, err := tagfilter.MapTagFilterToNormalizedString(apiObject.TagFilterExpression)
		if err != nil {
			diags.AddError(
				WebsiteAlertConfigErrMapFilterExpression,
				fmt.Sprintf(WebsiteAlertConfigErrMapFilterExpressionMsg, err),
			)
			return diags
		}
		model.TagFilter = util.SetStringPointerToState(filterExprStr)
	} else {
		model.TagFilter = types.StringNull()
	}

	return diags
}

// mapAlertChannelIDsToState maps alert channel IDs from API to state
func (r *websiteAlertConfigResource) mapAlertChannelIDsToState(model *WebsiteAlertConfigModel, apiObject *restapi.WebsiteAlertConfig) {
	if len(apiObject.AlertChannelIDs) > 0 {
		alertChannelIDs := make([]attr.Value, len(apiObject.AlertChannelIDs))
		for i, id := range apiObject.AlertChannelIDs {
			alertChannelIDs[i] = types.StringValue(id)
		}
		model.AlertChannelIDs = types.SetValueMust(types.StringType, alertChannelIDs)
	} else {
		model.AlertChannelIDs = types.SetNull(types.StringType)
	}
}

// mapCustomPayloadFieldsToState maps custom payload fields from API to state
func (r *websiteAlertConfigResource) mapCustomPayloadFieldsToState(ctx context.Context, model *WebsiteAlertConfigModel, apiObject *restapi.WebsiteAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, apiObject.CustomerPayloadFields)
	diags.Append(payloadDiags...)
	if !diags.HasError() {
		model.CustomPayloadFields = customPayloadFieldsList
	}

	return diags
}

func (r *websiteAlertConfigResource) mapTimeThresholdToState(timeThreshold restapi.WebsiteTimeThreshold) *WebsiteTimeThresholdModel {
	websiteTimeThresholdModel := WebsiteTimeThresholdModel{}
	switch timeThreshold.Type {
	case WebsiteAlertConfigTimeThresholdTypeViolationsInSequence:
		websiteTimeThresholdModel.ViolationsInSequence = &WebsiteViolationsInSequenceModel{
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
		}
	case WebsiteAlertConfigTimeThresholdTypeUserImpact:
		websiteTimeThresholdModel.UserImpactOfViolationsInSequence = &WebsiteUserImpactOfViolationsInSequenceModel{
			TimeWindow:              util.SetInt64PointerToState(timeThreshold.TimeWindow),
			ImpactMeasurementMethod: types.StringValue(string(*timeThreshold.ImpactMeasurementMethod)),
			UserPercentage:          util.SetFloat64PointerToState(timeThreshold.UserPercentage),
			Users:                   util.SetInt64PointerToState(timeThreshold.Users),
		}
	case WebsiteAlertConfigTimeThresholdTypeViolationsInPeriod:
		websiteTimeThresholdModel.ViolationsInPeriod = &WebsiteViolationsInPeriodModel{
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
			Violations: util.SetInt64PointerToState(timeThreshold.Violations),
		}
	}
	return &websiteTimeThresholdModel
}

func (r *websiteAlertConfigResource) mapRuleToState(ctx context.Context, rule *restapi.WebsiteAlertRule) *WebsiteAlertRuleModel {
	websiteAlertRuleModel := WebsiteAlertRuleModel{
		Slowness:        nil,
		SpecificJsError: nil,
		StatusCode:      nil,
		Throughput:      nil,
	}
	switch rule.AlertType {
	case WebsiteAlertConfigAlertTypeThroughput:
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
		}
		websiteAlertRuleModel.Throughput = &websiteAlertRuleConfigModel
	case WebsiteAlertConfigAlertTypeSlowness:
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
		}
		websiteAlertRuleModel.Slowness = &websiteAlertRuleConfigModel
	case WebsiteAlertConfigAlertTypeStatusCode:
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigCompleteModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
			Operator:    types.StringValue(string(*rule.Operator)),
			Value:       util.SetStringPointerToState(rule.Value),
		}
		websiteAlertRuleModel.StatusCode = &websiteAlertRuleConfigModel
	case WebsiteAlertConfigAlertTypeSpecificJsError:
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigCompleteModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
			Operator:    types.StringValue(string(*rule.Operator)),
			Value:       util.SetStringPointerToState(rule.Value),
		}
		websiteAlertRuleModel.SpecificJsError = &websiteAlertRuleConfigModel

	}

	return &websiteAlertRuleModel

}

func (r *websiteAlertConfigResource) mapRulesToState(ctx context.Context, apiObject *restapi.WebsiteAlertConfig) []RuleWithThresholdPluginModel {
	rules := apiObject.Rules
	var rulesModel []RuleWithThresholdPluginModel
	for _, i := range rules {
		warningThreshold, isWarningThresholdPresent := i.Thresholds[restapi.WarningSeverity]
		criticalThreshold, isCriticalThresholdPresent := i.Thresholds[restapi.CriticalSeverity]

		thresholdPluginModel := shared.ThresholdAllPluginModel{
			Warning:  shared.MapAllThresholdPluginToState(ctx, &warningThreshold, isWarningThresholdPresent),
			Critical: shared.MapAllThresholdPluginToState(ctx, &criticalThreshold, isCriticalThresholdPresent),
		}
		ruleModel := RuleWithThresholdPluginModel{
			Rule:              r.mapRuleToState(ctx, i.Rule),
			ThresholdOperator: types.StringValue(i.ThresholdOperator),
			Thresholds:        &thresholdPluginModel,
		}
		rulesModel = append(rulesModel, ruleModel)
	}
	return rulesModel

}
