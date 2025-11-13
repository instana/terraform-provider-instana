package logalertconfig

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewLogAlertConfigResourceHandleFramework creates the resource handle for Log Alert Configuration
func NewLogAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.LogAlertConfig] {
	return &logAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaLogAlertConfigFramework,
			Schema: schema.Schema{
				Description: LogAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: LogAlertConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					LogAlertConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: LogAlertConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					LogAlertConfigFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: LogAlertConfigDescDescription,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 65536),
						},
					},
					LogAlertConfigFieldGracePeriod: schema.Int64Attribute{
						Optional:    true,
						Description: LogAlertConfigDescGracePeriod,
					},
					LogAlertConfigFieldGranularity: schema.Int64Attribute{
						Optional:    true,
						Description: LogAlertConfigDescGranularity,
						Validators: []validator.Int64{
							int64validator.OneOf(int64(restapi.Granularity60000), int64(restapi.Granularity300000), int64(restapi.Granularity600000), int64(restapi.Granularity900000), int64(restapi.Granularity1200000), int64(restapi.Granularity1800000)),
						},
					},
					LogAlertConfigFieldTagFilter: schema.StringAttribute{
						Required:    true,
						Description: LogAlertConfigDescTagFilter,
					},
					shared.DefaultCustomPayloadFieldsName: shared.GetCustomPayloadFieldsSchema(),
					LogAlertConfigFieldAlertChannels: schema.SingleNestedAttribute{
						Description: LogAlertConfigDescAlertChannels,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							shared.ThresholdFieldWarning: schema.ListAttribute{
								Optional:    true,
								Description: LogAlertConfigDescAlertChannelIDs,
								ElementType: types.StringType,
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
								},
							},
							shared.ThresholdFieldCritical: schema.ListAttribute{
								Optional:    true,
								Description: LogAlertConfigDescAlertChannelIDs,
								ElementType: types.StringType,
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
								},
							},
						},
					},
					LogAlertConfigFieldGroupBy: schema.ListNestedAttribute{
						Description: LogAlertConfigDescGroupBy,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								LogAlertConfigFieldGroupByTagName: schema.StringAttribute{
									Required:    true,
									Description: LogAlertConfigDescGroupByTagName,
								},
								LogAlertConfigFieldGroupByKey: schema.StringAttribute{
									Optional:    true,
									Description: LogAlertConfigDescGroupByKey,
								},
							},
						},
					},
					LogAlertConfigFieldRules: schema.SingleNestedAttribute{
						Description: LogAlertConfigDescRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							LogAlertConfigFieldMetricName: schema.StringAttribute{
								Required:    true,
								Description: LogAlertConfigDescMetricName,
							},
							LogAlertConfigFieldAlertType: schema.StringAttribute{
								Optional:    true,
								Description: LogAlertConfigDescAlertType,
								Validators: []validator.String{
									stringvalidator.OneOf(LogAlertTypeLogCount),
								},
							},
							LogAlertConfigFieldAggregation: schema.StringAttribute{
								Optional:    true,
								Description: LogAlertConfigDescAggregation,
								Validators: []validator.String{
									stringvalidator.OneOf(string(restapi.SumAggregation)),
								},
							},
							LogAlertConfigFieldThresholdOperator: schema.StringAttribute{
								Required:    true,
								Description: LogAlertConfigDescThresholdOperator,
								Validators: []validator.String{
									stringvalidator.OneOf(restapi.SupportedThresholdOperators.ToStringSlice()...),
								},
							},
							LogAlertConfigFieldThreshold: schema.SingleNestedAttribute{
								Description: LogAlertConfigDescThreshold,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									LogAlertConfigFieldWarning:  shared.StaticThresholdAttributeSchema(),
									LogAlertConfigFieldCritical: shared.StaticThresholdAttributeSchema(),
								},
							},
						},
					},
					LogAlertConfigFieldTimeThreshold: schema.SingleNestedAttribute{
						Description: LogAlertConfigDescTimeThreshold,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"violations_in_sequence": schema.SingleNestedAttribute{
								Description: LogAlertConfigDescViolationsInSequence,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Description: LogAlertConfigDescTimeWindow,
										Required:    true,
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

type logAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *logAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *logAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.LogAlertConfig] {
	return api.LogAlertConfig()
}

func (r *logAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *logAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.LogAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the config
	model := LogAlertConfigModel{
		ID:          types.StringValue(config.ID),
		Name:        types.StringValue(config.Name),
		Description: types.StringValue(config.Description),
		Granularity: types.Int64Value(int64(config.Granularity)),
	}

	// Set grace period if present
	if config.GracePeriod > 0 {
		model.GracePeriod = types.Int64Value(config.GracePeriod)
	} else {
		model.GracePeriod = types.Int64Null()
	}

	// Set tag filter
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			diags.AddError(
				LogAlertConfigErrNormalizingTagFilter,
				LogAlertConfigErrNormalizingTagFilterMsg+err.Error(),
			)
			return diags
		}
		model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
	} else {
		model.TagFilter = types.StringNull()
	}

	// Set group by
	if len(config.GroupBy) > 0 {
		groupByModels, groupByDiags := r.mapGroupByToState(ctx, config.GroupBy)
		if groupByDiags.HasError() {
			diags.Append(groupByDiags...)
			return diags
		}
		model.GroupBy = groupByModels
	} else {
		model.GroupBy = []GroupByModel{}
	}

	// Set alert channels
	alertChannelsModel, alertChannelsDiags := r.mapAlertChannelsToState(ctx, config.AlertChannels)
	if alertChannelsDiags.HasError() {
		diags.Append(alertChannelsDiags...)
		return diags
	}
	model.AlertChannels = alertChannelsModel

	// // Set time threshold
	// timeThresholdList, timeThresholdDiags := r.mapTimeThresholdToState(ctx, config.TimeThreshold)
	// if timeThresholdDiags.HasError() {
	// 	diags.Append(timeThresholdDiags...)
	// 	return diags
	// }
	model.TimeThreshold = r.mapTimeThresholdToState(ctx, config.TimeThreshold)

	// Set rules
	rulesModel, rulesDiags := r.mapRulesToState(ctx, config.Rules)
	if rulesDiags.HasError() {
		diags.Append(rulesDiags...)
		return diags
	}
	model.Rules = rulesModel

	// Set custom payload fields
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, config.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *logAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.LogAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model LogAlertConfigModel

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
		var err error
		tagFilter, err = r.mapTagFilterExpressionFromState(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError(
				LogAlertConfigErrParsingTagFilter,
				LogAlertConfigErrParsingTagFilterMsg+err.Error(),
			)
			return nil, diags
		}
	}

	// Map group by
	groupBy := []restapi.GroupByTag{}
	if len(model.GroupBy) > 0 {
		var groupByDiags diag.Diagnostics
		groupBy, groupByDiags = r.mapGroupByFromState(ctx, model.GroupBy)
		if groupByDiags.HasError() {
			diags.Append(groupByDiags...)
			return nil, diags
		}
	}

	// Map alert channels
	alertChannels := make(map[restapi.AlertSeverity][]string)
	if model.AlertChannels != nil {
		var alertChannelsDiags diag.Diagnostics
		alertChannels, alertChannelsDiags = r.mapAlertChannelsFromState(ctx, model.AlertChannels)
		if alertChannelsDiags.HasError() {
			diags.Append(alertChannelsDiags...)
			return nil, diags
		}
	}

	// Map time threshold
	// var timeThreshold *restapi.LogTimeThreshold
	// if model.TimeThreshold == nil {
	// 	var timeThresholdDiags diag.Diagnostics
	// 	timeThreshold, timeThresholdDiags = r.mapTimeThresholdFromState(ctx, model.TimeThreshold)
	// 	if timeThresholdDiags.HasError() {
	// 		diags.Append(timeThresholdDiags...)
	// 		return nil, diags
	// 	}
	// } else {
	// 	timeThreshold = &restapi.LogTimeThreshold{}
	// }

	//timeThreshold = r.mapTimeThresholdFromState(ctx, model.TimeThreshold)
	// Map rules
	var rules []restapi.RuleWithThreshold[restapi.LogAlertRule]
	if model.Rules != nil {
		var rulesDiags diag.Diagnostics
		rules, rulesDiags = r.mapRulesFromState(ctx, model.Rules)
		if rulesDiags.HasError() {
			diags.Append(rulesDiags...)
			return nil, diags
		}
	}

	// Map custom payload fields
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() {
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
	}

	// Create the log alert config
	config := &restapi.LogAlertConfig{
		ID:                    id,
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		TagFilterExpression:   tagFilter,
		GroupBy:               groupBy,
		AlertChannels:         alertChannels,
		CustomerPayloadFields: customerPayloadFields,
		Rules:                 rules,
		TimeThreshold:         r.mapTimeThresholdFromState(ctx, model.TimeThreshold),
	}

	// Set granularity
	if !model.Granularity.IsNull() {
		config.Granularity = restapi.Granularity(model.Granularity.ValueInt64())
	} else {
		config.Granularity = restapi.Granularity600000
	}

	// Set grace period
	if !model.GracePeriod.IsNull() {
		config.GracePeriod = model.GracePeriod.ValueInt64()
	}

	return config, diags
}

func (r *logAlertConfigResourceFramework) mapTimeThresholdFromState(ctx context.Context, tf *TimeThresholdModel) *restapi.LogTimeThreshold {
	if tf == nil || tf.ViolationsInSequence == nil {
		return nil
	}

	v := tf.ViolationsInSequence
	if v.TimeWindow.IsNull() || v.TimeWindow.IsUnknown() {
		return nil
	}

	return &restapi.LogTimeThreshold{
		Type:       "violationsInSequence",
		TimeWindow: v.TimeWindow.ValueInt64(),
	}
}

// Helper methods for mapping between state and API models

func (r *logAlertConfigResourceFramework) mapTagFilterExpressionFromState(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *logAlertConfigResourceFramework) mapGroupByToState(ctx context.Context, groupBy []restapi.GroupByTag) ([]GroupByModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(groupBy) == 0 {
		return []GroupByModel{}, diags
	}

	result := make([]GroupByModel, len(groupBy))
	for i, g := range groupBy {
		model := GroupByModel{
			TagName: types.StringValue(g.TagName),
		}

		if g.Key != "" {
			model.Key = types.StringValue(g.Key)
		} else {
			model.Key = types.StringNull()
		}

		result[i] = model
	}

	return result, diags
}

func (r *logAlertConfigResourceFramework) mapGroupByFromState(ctx context.Context, groupByModels []GroupByModel) ([]restapi.GroupByTag, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(groupByModels) == 0 {
		return []restapi.GroupByTag{}, diags
	}

	result := make([]restapi.GroupByTag, len(groupByModels))
	for i, model := range groupByModels {
		groupByTag := restapi.GroupByTag{
			TagName: model.TagName.ValueString(),
		}

		if !model.Key.IsNull() {
			groupByTag.Key = model.Key.ValueString()
		}

		result[i] = groupByTag
	}

	return result, diags
}

func (r *logAlertConfigResourceFramework) mapAlertChannelsToState(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string) (*AlertChannelsModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(alertChannels) == 0 {
		return nil, diags
	}

	model := &AlertChannelsModel{}

	// Map warning channels
	if warningChannels, ok := alertChannels[restapi.WarningSeverity]; ok && len(warningChannels) > 0 {
		warningList, warningDiags := types.ListValueFrom(ctx, types.StringType, warningChannels)
		diags.Append(warningDiags...)
		if diags.HasError() {
			return nil, diags
		}
		model.Warning = warningList
	} else {
		model.Warning = types.ListNull(types.StringType)
	}

	// Map critical channels
	if criticalChannels, ok := alertChannels[restapi.CriticalSeverity]; ok && len(criticalChannels) > 0 {
		criticalList, criticalDiags := types.ListValueFrom(ctx, types.StringType, criticalChannels)
		diags.Append(criticalDiags...)
		if diags.HasError() {
			return nil, diags
		}
		model.Critical = criticalList
	} else {
		model.Critical = types.ListNull(types.StringType)
	}

	return model, diags
}

func (r *logAlertConfigResourceFramework) mapAlertChannelsFromState(ctx context.Context, model *AlertChannelsModel) (map[restapi.AlertSeverity][]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannels := make(map[restapi.AlertSeverity][]string)

	if model == nil {
		return alertChannels, diags
	}

	// Map warning channels
	if !model.Warning.IsNull() && !model.Warning.IsUnknown() {
		var warningChannels []string
		diags.Append(model.Warning.ElementsAs(ctx, &warningChannels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(warningChannels) > 0 {
			alertChannels[restapi.WarningSeverity] = warningChannels
		}
	}

	// Map critical channels
	if !model.Critical.IsNull() && !model.Critical.IsUnknown() {
		var criticalChannels []string
		diags.Append(model.Critical.ElementsAs(ctx, &criticalChannels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(criticalChannels) > 0 {
			alertChannels[restapi.CriticalSeverity] = criticalChannels
		}
	}

	return alertChannels, diags
}

func (r *logAlertConfigResourceFramework) mapTimeThresholdToState(ctx context.Context, api *restapi.LogTimeThreshold) *TimeThresholdModel {
	if api == nil {
		return nil
	}

	if api.Type != "violationsInSequence" {
		// unsupported type — ignore or handle others
		return nil
	}

	return &TimeThresholdModel{
		ViolationsInSequence: &ViolationsInSequenceModel{
			TimeWindow: types.Int64Value(api.TimeWindow),
		},
	}
}

func (r *logAlertConfigResourceFramework) mapRulesToState(ctx context.Context, rules []restapi.RuleWithThreshold[restapi.LogAlertRule]) (*LogAlertRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(rules) == 0 {
		return nil, diags
	}

	// Since rules is now SingleNestedAttribute, we only handle the first rule
	ruleWithThreshold := rules[0]
	rule := ruleWithThreshold.Rule

	// Convert "logCount" to "log.count" for the schema
	alertType := rule.AlertType
	if alertType == "logCount" {
		alertType = LogAlertTypeLogCount
	}

	// Create rule model
	ruleModel := &LogAlertRuleModel{
		MetricName:        types.StringValue(rule.MetricName),
		AlertType:         types.StringValue(alertType),
		ThresholdOperator: types.StringValue(string(ruleWithThreshold.ThresholdOperator)),
	}

	if rule.Aggregation != "" {
		ruleModel.Aggregation = types.StringValue(string(rule.Aggregation))
	} else {
		ruleModel.Aggregation = types.StringValue(string(restapi.SumAggregation))
	}

	// Map thresholds
	thresholdModel := &ThresholdModel{}

	// Map warning threshold
	if warningThreshold, ok := ruleWithThreshold.Thresholds[restapi.WarningSeverity]; ok {
		thresholdModel.Warning = &shared.ThresholdStaticTypeModel{
			Static: &shared.StaticTypeModel{
				Operator: types.StringNull(),
				Value:    types.Int64Value(int64(*warningThreshold.Value)),
			},
		}
	}

	// Map critical threshold
	if criticalThreshold, ok := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]; ok {
		thresholdModel.Critical = &shared.ThresholdStaticTypeModel{
			Static: &shared.StaticTypeModel{
				Operator: types.StringNull(),
				Value:    types.Int64Value(int64(*criticalThreshold.Value)),
			},
		}
	}

	ruleModel.Threshold = thresholdModel

	return ruleModel, diags
}

func (r *logAlertConfigResourceFramework) mapRulesFromState(ctx context.Context, ruleModel *LogAlertRuleModel) ([]restapi.RuleWithThreshold[restapi.LogAlertRule], diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleModel == nil {
		return []restapi.RuleWithThreshold[restapi.LogAlertRule]{}, diags
	}

	// Convert "log.count" to "logCount" for the API
	alertType := ruleModel.AlertType.ValueString()
	if alertType == LogAlertTypeLogCount {
		alertType = "logCount"
	}

	// Create log alert rule
	logAlertRule := restapi.LogAlertRule{
		MetricName: ruleModel.MetricName.ValueString(),
		AlertType:  alertType,
	}

	if !ruleModel.Aggregation.IsNull() {
		logAlertRule.Aggregation = restapi.Aggregation(ruleModel.Aggregation.ValueString())
	}

	// Get threshold operator
	thresholdOperator := restapi.ThresholdOperator(ruleModel.ThresholdOperator.ValueString())

	// Map thresholds
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	if ruleModel.Threshold != nil {
		// Map warning threshold
		if ruleModel.Threshold.Warning != nil &&
			ruleModel.Threshold.Warning.Static != nil && !ruleModel.Threshold.Warning.Static.Value.IsNull() {
			valueFloat := float64(ruleModel.Threshold.Warning.Static.Value.ValueInt64())
			thresholdMap[restapi.WarningSeverity] = restapi.ThresholdRule{
				Type:  "staticThreshold",
				Value: &valueFloat,
			}
		}

		// Map critical threshold
		if ruleModel.Threshold.Critical != nil && ruleModel.Threshold.Critical.Static != nil &&
			!ruleModel.Threshold.Critical.Static.Value.IsNull() {
			valueFloat := float64(ruleModel.Threshold.Critical.Static.Value.ValueInt64())
			thresholdMap[restapi.CriticalSeverity] = restapi.ThresholdRule{
				Type:  "staticThreshold",
				Value: &valueFloat,
			}
		}
	}

	// Create rule with threshold - return as single-element array
	result := []restapi.RuleWithThreshold[restapi.LogAlertRule]{
		{
			ThresholdOperator: thresholdOperator,
			Rule:              logAlertRule,
			Thresholds:        thresholdMap,
		},
	}

	return result, diags
}

// mapSingleThresholdFromState maps a single threshold (static) to API model
func (r *logAlertConfigResourceFramework) mapSingleThresholdFromState(ctx context.Context, thresholdObj types.Object) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	var threshold struct {
		Static types.Object `tfsdk:"static"`
	}

	diags.Append(thresholdObj.As(ctx, &threshold, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	if threshold.Static.IsNull() || threshold.Static.IsUnknown() {
		return nil, diags
	}

	var static struct {
		Operator types.String `tfsdk:"operator"`
		Value    types.Int64  `tfsdk:"value"`
	}

	diags.Append(threshold.Static.As(ctx, &static, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	valueFloat := float64(static.Value.ValueInt64())
	return &restapi.ThresholdRule{
		Type:  "staticThreshold",
		Value: &valueFloat,
	}, diags
}

// mapThresholdRuleFromState and MapThresholdsFromState have been moved to threshold-mapping-framework.go
