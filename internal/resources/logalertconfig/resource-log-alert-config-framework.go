package logalertconfig

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
					LogAlertConfigFieldAlertChannels: schema.ListNestedAttribute{
						Description: LogAlertConfigDescAlertChannels,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								shared.ThresholdFieldWarning: schema.ListAttribute{
									Optional:    true,
									Description: LogAlertConfigDescAlertChannelIDs,
									ElementType: types.StringType,
								},
								shared.ThresholdFieldCritical: schema.ListAttribute{
									Optional:    true,
									Description: LogAlertConfigDescAlertChannelIDs,
									ElementType: types.StringType,
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
		groupByList, groupByDiags := r.mapGroupByToState(ctx, config.GroupBy)
		if groupByDiags.HasError() {
			diags.Append(groupByDiags...)
			return diags
		}
		model.GroupBy = groupByList
	} else {
		model.GroupBy = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldGroupByTagName: types.StringType,
				LogAlertConfigFieldGroupByKey:     types.StringType,
			},
		})
	}

	// Set alert channels
	alertChannelsList, alertChannelsDiags := shared.MapAlertChannelsToState(ctx, config.AlertChannels)
	if alertChannelsDiags.HasError() {
		diags.Append(alertChannelsDiags...)
		return diags
	}
	model.AlertChannels = alertChannelsList

	// // Set time threshold
	// timeThresholdList, timeThresholdDiags := r.mapTimeThresholdToState(ctx, config.TimeThreshold)
	// if timeThresholdDiags.HasError() {
	// 	diags.Append(timeThresholdDiags...)
	// 	return diags
	// }
	model.TimeThreshold = r.mapTimeThresholdToState(ctx, config.TimeThreshold)

	// Set rules
	rulesList, rulesDiags := r.mapRulesToState(ctx, config.Rules)
	if rulesDiags.HasError() {
		diags.Append(rulesDiags...)
		return diags
	}
	model.Rules = rulesList

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
	if !model.GroupBy.IsNull() {
		var groupByDiags diag.Diagnostics
		groupBy, groupByDiags = r.mapGroupByFromState(ctx, model.GroupBy)
		if groupByDiags.HasError() {
			diags.Append(groupByDiags...)
			return nil, diags
		}
	}

	// Map alert channels
	alertChannels := make(map[restapi.AlertSeverity][]string)
	if !model.AlertChannels.IsNull() {
		var alertChannelsDiags diag.Diagnostics
		alertChannels, alertChannelsDiags = shared.MapAlertChannelsFromState(ctx, model.AlertChannels)
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
	if !model.Rules.IsNull() {
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

func (r *logAlertConfigResourceFramework) mapGroupByToState(ctx context.Context, groupBy []restapi.GroupByTag) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(groupBy) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldGroupByTagName: types.StringType,
				LogAlertConfigFieldGroupByKey:     types.StringType,
			},
		}), diags
	}

	groupByElements := make([]attr.Value, len(groupBy))
	for i, g := range groupBy {
		groupByObj := map[string]attr.Value{
			LogAlertConfigFieldGroupByTagName: types.StringValue(g.TagName),
		}

		if g.Key != "" {
			groupByObj[LogAlertConfigFieldGroupByKey] = types.StringValue(g.Key)
		} else {
			groupByObj[LogAlertConfigFieldGroupByKey] = types.StringNull()
		}

		objVal, objDiags := types.ObjectValue(
			map[string]attr.Type{
				LogAlertConfigFieldGroupByTagName: types.StringType,
				LogAlertConfigFieldGroupByKey:     types.StringType,
			},
			groupByObj,
		)
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		groupByElements[i] = objVal
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldGroupByTagName: types.StringType,
				LogAlertConfigFieldGroupByKey:     types.StringType,
			},
		},
		groupByElements,
	)
}

func (r *logAlertConfigResourceFramework) mapGroupByFromState(ctx context.Context, groupByList types.List) ([]restapi.GroupByTag, diag.Diagnostics) {
	var diags diag.Diagnostics

	if groupByList.IsNull() || groupByList.IsUnknown() {
		return []restapi.GroupByTag{}, diags
	}

	var groupByElements []types.Object
	diags.Append(groupByList.ElementsAs(ctx, &groupByElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]restapi.GroupByTag, len(groupByElements))
	for i, element := range groupByElements {
		var groupBy struct {
			TagName types.String `tfsdk:"tag_name"`
			Key     types.String `tfsdk:"key"`
		}

		diags.Append(element.As(ctx, &groupBy, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		groupByTag := restapi.GroupByTag{
			TagName: groupBy.TagName.ValueString(),
		}

		if !groupBy.Key.IsNull() {
			groupByTag.Key = groupBy.Key.ValueString()
		}

		result[i] = groupByTag
	}

	return result, diags
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

func (r *logAlertConfigResourceFramework) mapTimeThresholdToState1(ctx context.Context, timeThreshold *restapi.LogTimeThreshold) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if timeThreshold == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldTimeThresholdTimeWindow: types.Int64Type,
			},
		}), diags
	}

	// Create time threshold object
	timeThresholdObj := map[string]attr.Value{
		LogAlertConfigFieldTimeThresholdTimeWindow: types.Int64Value(timeThreshold.TimeWindow),
	}

	timeThresholdObjVal, timeThresholdObjDiags := types.ObjectValue(
		map[string]attr.Type{
			LogAlertConfigFieldTimeThresholdTimeWindow: types.Int64Type,
		},
		timeThresholdObj,
	)
	diags.Append(timeThresholdObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldTimeThresholdTimeWindow: types.Int64Type,
			},
		},
		[]attr.Value{timeThresholdObjVal},
	)
}

func (r *logAlertConfigResourceFramework) mapRulesToState(ctx context.Context, rules []restapi.RuleWithThreshold[restapi.LogAlertRule]) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	thresholdAttrTypes := map[string]attr.Type{
		LogAlertConfigFieldWarning:  types.ObjectType{AttrTypes: shared.GetStaticThresholdAttrTypes()},
		LogAlertConfigFieldCritical: types.ObjectType{AttrTypes: shared.GetStaticThresholdAttrTypes()},
	}

	ruleAttrTypes := map[string]attr.Type{
		LogAlertConfigFieldMetricName:        types.StringType,
		LogAlertConfigFieldAlertType:         types.StringType,
		LogAlertConfigFieldAggregation:       types.StringType,
		LogAlertConfigFieldThresholdOperator: types.StringType,
		LogAlertConfigFieldThreshold:         types.ObjectType{AttrTypes: thresholdAttrTypes},
	}

	if len(rules) == 0 {
		return types.ObjectNull(ruleAttrTypes), diags
	}

	// Since rules is now SingleNestedAttribute, we only handle the first rule
	ruleWithThreshold := rules[0]
	rule := ruleWithThreshold.Rule

	// Convert "logCount" to "log.count" for the schema
	alertType := rule.AlertType
	if alertType == "logCount" {
		alertType = LogAlertTypeLogCount
	}

	// Create rule object
	ruleObj := map[string]attr.Value{
		LogAlertConfigFieldMetricName:        types.StringValue(rule.MetricName),
		LogAlertConfigFieldAlertType:         types.StringValue(alertType),
		LogAlertConfigFieldThresholdOperator: types.StringValue(string(ruleWithThreshold.ThresholdOperator)),
	}

	if rule.Aggregation != "" {
		ruleObj[LogAlertConfigFieldAggregation] = types.StringValue(string(rule.Aggregation))
	} else {
		ruleObj[LogAlertConfigFieldAggregation] = types.StringValue(string(restapi.SumAggregation))
	}

	// Map thresholds - create nested objects for warning and critical
	thresholdObj := map[string]attr.Value{}

	// Map warning threshold - convert to object
	warningThreshold, isWarningThresholdPresent := ruleWithThreshold.Thresholds[restapi.WarningSeverity]
	if isWarningThresholdPresent {
		staticInnerObj, staticInnerDiags := types.ObjectValue(
			map[string]attr.Type{
				"operator": types.StringType,
				"value":    types.Int64Type,
			},
			map[string]attr.Value{
				"operator": types.StringNull(),
				"value":    types.Int64Value(int64(*warningThreshold.Value)),
			},
		)
		diags.Append(staticInnerDiags...)
		if diags.HasError() {
			return types.ObjectNull(ruleAttrTypes), diags
		}

		warningStaticObj := map[string]attr.Value{
			"static": staticInnerObj,
		}
		warningObj, warnDiags := types.ObjectValue(shared.GetStaticThresholdAttrTypes(), warningStaticObj)
		diags.Append(warnDiags...)
		if diags.HasError() {
			return types.ObjectNull(ruleAttrTypes), diags
		}
		thresholdObj[LogAlertConfigFieldWarning] = warningObj
	} else {
		thresholdObj[LogAlertConfigFieldWarning] = types.ObjectNull(shared.GetStaticThresholdAttrTypes())
	}

	// Map critical threshold - convert to object
	criticalThreshold, isCriticalThresholdPresent := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]
	if isCriticalThresholdPresent {
		staticInnerObj, staticInnerDiags := types.ObjectValue(
			map[string]attr.Type{
				"operator": types.StringType,
				"value":    types.Int64Type,
			},
			map[string]attr.Value{
				"operator": types.StringNull(),
				"value":    types.Int64Value(int64(*criticalThreshold.Value)),
			},
		)
		diags.Append(staticInnerDiags...)
		if diags.HasError() {
			return types.ObjectNull(ruleAttrTypes), diags
		}

		criticalStaticObj := map[string]attr.Value{
			"static": staticInnerObj,
		}
		criticalObj, critDiags := types.ObjectValue(shared.GetStaticThresholdAttrTypes(), criticalStaticObj)
		diags.Append(critDiags...)
		if diags.HasError() {
			return types.ObjectNull(ruleAttrTypes), diags
		}
		thresholdObj[LogAlertConfigFieldCritical] = criticalObj
	} else {
		thresholdObj[LogAlertConfigFieldCritical] = types.ObjectNull(shared.GetStaticThresholdAttrTypes())
	}

	// Create threshold object value
	thresholdObjVal, thresholdObjDiags := types.ObjectValue(thresholdAttrTypes, thresholdObj)
	diags.Append(thresholdObjDiags...)
	if diags.HasError() {
		return types.ObjectNull(ruleAttrTypes), diags
	}

	ruleObj[LogAlertConfigFieldThreshold] = thresholdObjVal

	// Create rule object value
	return types.ObjectValue(ruleAttrTypes, ruleObj)
}

func (r *logAlertConfigResourceFramework) mapThresholdRuleToState(ctx context.Context, threshold *restapi.ThresholdRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get the threshold value
	var thresholdValue int64
	if threshold.Value != nil {
		thresholdValue = int64(*threshold.Value)
	} else {
		thresholdValue = 0
	}

	// Create static threshold object
	staticObj := map[string]attr.Value{
		LogAlertConfigFieldValue: types.Int64Value(thresholdValue),
	}

	staticObjVal, staticObjDiags := types.ObjectValue(
		map[string]attr.Type{
			LogAlertConfigFieldValue: types.Int64Type,
		},
		staticObj,
	)
	diags.Append(staticObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create static list
	staticList, staticListDiags := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldValue: types.Int64Type,
			},
		},
		[]attr.Value{staticObjVal},
	)
	diags.Append(staticListDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create threshold object with static block
	thresholdObj := map[string]attr.Value{
		"static": staticList,
	}

	thresholdObjVal, thresholdObjDiags := types.ObjectValue(
		shared.GetStaticThresholdAttrTypes(),
		thresholdObj,
	)
	diags.Append(thresholdObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: shared.GetStaticThresholdAttrTypes(),
		},
		[]attr.Value{thresholdObjVal},
	)
}

func (r *logAlertConfigResourceFramework) mapRulesFromState(ctx context.Context, rulesObj types.Object) ([]restapi.RuleWithThreshold[restapi.LogAlertRule], diag.Diagnostics) {
	var diags diag.Diagnostics

	if rulesObj.IsNull() || rulesObj.IsUnknown() {
		return []restapi.RuleWithThreshold[restapi.LogAlertRule]{}, diags
	}

	var rule struct {
		MetricName        types.String `tfsdk:"metric_name"`
		AlertType         types.String `tfsdk:"alert_type"`
		Aggregation       types.String `tfsdk:"aggregation"`
		ThresholdOperator types.String `tfsdk:"threshold_operator"`
		Threshold         types.Object `tfsdk:"threshold"`
	}

	diags.Append(rulesObj.As(ctx, &rule, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert "log.count" to "logCount" for the API
	alertType := rule.AlertType.ValueString()
	if alertType == LogAlertTypeLogCount {
		alertType = "logCount"
	}

	// Create log alert rule
	logAlertRule := restapi.LogAlertRule{
		MetricName: rule.MetricName.ValueString(),
		AlertType:  alertType,
	}

	if !rule.Aggregation.IsNull() {
		logAlertRule.Aggregation = restapi.Aggregation(rule.Aggregation.ValueString())
	}

	// Get threshold operator
	thresholdOperator := restapi.ThresholdOperator(rule.ThresholdOperator.ValueString())

	// Map thresholds from object
	var thresholdMap map[restapi.AlertSeverity]restapi.ThresholdRule
	var thresholdDiags diag.Diagnostics

	if !rule.Threshold.IsNull() && !rule.Threshold.IsUnknown() {
		thresholdMap, thresholdDiags = r.mapThresholdObjectFromState(ctx, rule.Threshold)
		diags.Append(thresholdDiags...)
		if diags.HasError() {
			return nil, diags
		}
	} else {
		thresholdMap = make(map[restapi.AlertSeverity]restapi.ThresholdRule)
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

// mapThresholdObjectFromState maps threshold object (warning/critical) to API model
func (r *logAlertConfigResourceFramework) mapThresholdObjectFromState(ctx context.Context, thresholdObj types.Object) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	var thresholds struct {
		Warning  types.Object `tfsdk:"warning"`
		Critical types.Object `tfsdk:"critical"`
	}

	diags.Append(thresholdObj.As(ctx, &thresholds, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Map warning threshold
	if !thresholds.Warning.IsNull() && !thresholds.Warning.IsUnknown() {
		warningThreshold, warnDiags := r.mapSingleThresholdFromState(ctx, thresholds.Warning)
		diags.Append(warnDiags...)
		if diags.HasError() {
			return nil, diags
		}
		if warningThreshold != nil {
			thresholdMap[restapi.WarningSeverity] = *warningThreshold
		}
	}

	// Map critical threshold
	if !thresholds.Critical.IsNull() && !thresholds.Critical.IsUnknown() {
		criticalThreshold, critDiags := r.mapSingleThresholdFromState(ctx, thresholds.Critical)
		diags.Append(critDiags...)
		if diags.HasError() {
			return nil, diags
		}
		if criticalThreshold != nil {
			thresholdMap[restapi.CriticalSeverity] = *criticalThreshold
		}
	}

	return thresholdMap, diags
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
