package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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

// ResourceInstanaLogAlertConfigFramework the name of the terraform-provider-instana resource to manage log alert configurations
const ResourceInstanaLogAlertConfigFramework = "log_alert_config"

// LogAlertConfigModel represents the data model for the log alert configuration resource
type LogAlertConfigModel struct {
	ID                  types.String        `tfsdk:"id"`
	Name                types.String        `tfsdk:"name"`
	Description         types.String        `tfsdk:"description"`
	AlertChannels       types.List          `tfsdk:"alert_channels"`
	GracePeriod         types.Int64         `tfsdk:"grace_period"`
	GroupBy             types.List          `tfsdk:"group_by"`
	Granularity         types.Int64         `tfsdk:"granularity"`
	TagFilter           types.String        `tfsdk:"tag_filter"`
	Rules               types.List          `tfsdk:"rules"`
	TimeThreshold       *TimeThresholdModel `tfsdk:"time_threshold"`
	CustomPayloadFields types.List          `tfsdk:"custom_payload_field"`
}

type TimeThresholdModel struct {
	ViolationsInSequence *ViolationsInSequenceModel `tfsdk:"violations_in_sequence"`
}

type ViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}

// GroupByModel represents a group by tag in the Terraform model
type GroupByModel struct {
	TagName types.String `tfsdk:"tag_name"`
	Key     types.String `tfsdk:"key"`
}

// AlertChannelsModel represents alert channels in the Terraform model
type AlertChannelsModel struct {
	Warning  types.List `tfsdk:"warning"`
	Critical types.List `tfsdk:"critical"`
}

// LogAlertRuleModel represents a log alert rule in the Terraform model
type LogAlertRuleModel struct {
	MetricName        types.String `tfsdk:"metric_name"`
	AlertType         types.String `tfsdk:"alert_type"`
	Aggregation       types.String `tfsdk:"aggregation"`
	ThresholdOperator types.String `tfsdk:"threshold_operator"`
	Threshold         types.List   `tfsdk:"threshold"`
}

// ThresholdModel represents a threshold in the Terraform model
type ThresholdModel struct {
	Warning  types.List `tfsdk:"warning"`
	Critical types.List `tfsdk:"critical"`
}

// NewLogAlertConfigResourceHandleFramework creates the resource handle for Log Alert Configuration
func NewLogAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.LogAlertConfig] {
	return &logAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaLogAlertConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages log alert configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the log alert configuration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					LogAlertConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: "Name for the Log alert configuration",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					LogAlertConfigFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: "The description text of the Log alert config",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 65536),
						},
					},
					LogAlertConfigFieldGracePeriod: schema.Int64Attribute{
						Optional:    true,
						Description: "The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.",
					},
					LogAlertConfigFieldGranularity: schema.Int64Attribute{
						Optional:    true,
						Description: "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
						Validators: []validator.Int64{
							int64validator.OneOf(int64(restapi.Granularity60000), int64(restapi.Granularity300000), int64(restapi.Granularity600000), int64(restapi.Granularity900000), int64(restapi.Granularity1200000), int64(restapi.Granularity1800000)),
						},
					},
					LogAlertConfigFieldTagFilter: schema.StringAttribute{
						Required:    true,
						Description: "The tag filter expression used for this log alert",
					},
				},
				Blocks: map[string]schema.Block{
					LogAlertConfigFieldAlertChannels: schema.ListNestedBlock{
						Description: "Set of alert channel IDs associated with the severity.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								ResourceFieldThresholdRuleWarningSeverity: schema.ListAttribute{
									Optional:    true,
									Description: "List of IDs of alert channels defined in Instana.",
									ElementType: types.StringType,
								},
								ResourceFieldThresholdRuleCriticalSeverity: schema.ListAttribute{
									Optional:    true,
									Description: "List of IDs of alert channels defined in Instana.",
									ElementType: types.StringType,
								},
							},
						},
					},
					LogAlertConfigFieldGroupBy: schema.ListNestedBlock{
						Description: "The grouping tags used to group the metric results.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								LogAlertConfigFieldGroupByTagName: schema.StringAttribute{
									Required:    true,
									Description: "The tag name used for grouping",
								},
								LogAlertConfigFieldGroupByKey: schema.StringAttribute{
									Optional:    true,
									Description: "The key used for grouping",
								},
							},
						},
					},
					LogAlertConfigFieldRules: schema.ListNestedBlock{
						Description: "Configuration for the log alert rule",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								LogAlertConfigFieldMetricName: schema.StringAttribute{
									Required:    true,
									Description: "The metric name of the log alert rule",
								},
								LogAlertConfigFieldAlertType: schema.StringAttribute{
									Optional:    true,
									Description: "The type of the log alert rule (only 'log.count' is supported)",
									Validators: []validator.String{
										stringvalidator.OneOf(LogAlertTypeLogCount),
									},
								},
								LogAlertConfigFieldAggregation: schema.StringAttribute{
									Optional:    true,
									Description: "The aggregation method to use for the log alert (only 'SUM' is supported)",
									Validators: []validator.String{
										stringvalidator.OneOf(string(restapi.SumAggregation)),
									},
								},
								LogAlertConfigFieldThresholdOperator: schema.StringAttribute{
									Required:    true,
									Description: "The operator which will be applied to evaluate the threshold",
									Validators: []validator.String{
										stringvalidator.OneOf(restapi.SupportedThresholdOperators.ToStringSlice()...),
									},
								},
							},
							Blocks: map[string]schema.Block{
								LogAlertConfigFieldThreshold: schema.ListNestedBlock{
									Description: "Threshold configuration for different severity levels",
									NestedObject: schema.NestedBlockObject{
										Blocks: map[string]schema.Block{
											LogAlertConfigFieldWarning:  StaticThresholdBlockSchema(),
											LogAlertConfigFieldCritical: StaticThresholdBlockSchema(),
										},
									},
									Validators: []validator.List{
										listvalidator.SizeAtMost(1),
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtMost(1),
						},
					},
					LogAlertConfigFieldTimeThreshold: schema.SingleNestedBlock{
						Description: "Indicates the type of violation of the defined threshold.",
						Blocks: map[string]schema.Block{
							"violations_in_sequence": schema.SingleNestedBlock{
								Description: "Time threshold base on violations in sequence",
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Description: "Time window in milliseconds.",
										Required:    true,
									},
								},
							},
						},
					},
					DefaultCustomPayloadFieldsName: GetCustomPayloadFieldsSchema(),
				},
			},
			SchemaVersion: 1,
		},
	}
}

type logAlertConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *logAlertConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
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
				"Error normalizing tag filter",
				"Could not normalize tag filter: "+err.Error(),
			)
			return diags
		}
		model.TagFilter = types.StringValue(*normalizedTagFilterString)
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
	alertChannelsList, alertChannelsDiags := MapAlertChannelsToState(ctx, config.AlertChannels)
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
	customPayloadFieldsList, payloadDiags := CustomPayloadFieldsToTerraform(ctx, config.CustomerPayloadFields)
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
				"Error parsing tag filter",
				"Could not parse tag filter: "+err.Error(),
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
		alertChannels, alertChannelsDiags = MapAlertChannelsFromState(ctx, model.AlertChannels)
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
		customerPayloadFields, payloadDiags = MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
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

func (r *logAlertConfigResourceFramework) mapRulesToState(ctx context.Context, rules []restapi.RuleWithThreshold[restapi.LogAlertRule]) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(rules) == 0 {
		return types.ListNull(types.ObjectType{}), diags
	}

	ruleElements := make([]attr.Value, len(rules))

	for i, ruleWithThreshold := range rules {
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

		// Map thresholds
		thresholdObj := map[string]attr.Value{}

		// Map warning threshold
		warningThreshold, isWarningThresholdPresent := ruleWithThreshold.Thresholds[restapi.WarningSeverity]
		warningThresholdList, warningDiags := MapThresholdToState(ctx, isWarningThresholdPresent, &warningThreshold, []string{"static"})
		diags.Append(warningDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		thresholdObj[LogAlertConfigFieldWarning] = warningThresholdList

		// Map critical threshold
		criticalThreshold, isCriticalThresholdPresent := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]
		criticalThresholdList, criticalDiags := MapThresholdToState(ctx, isCriticalThresholdPresent, &criticalThreshold, []string{"static"})
		diags.Append(criticalDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		thresholdObj[LogAlertConfigFieldCritical] = criticalThresholdList

		// Create threshold object value
		thresholdObjVal, thresholdObjDiags := types.ObjectValue(
			map[string]attr.Type{
				LogAlertConfigFieldWarning:  GetStaticThresholdAttrListTypes(),
				LogAlertConfigFieldCritical: GetStaticThresholdAttrListTypes(),
			},
			thresholdObj,
		)
		diags.Append(thresholdObjDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		// Add threshold to rule
		thresholdList, thresholdListDiags := types.ListValue(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					LogAlertConfigFieldWarning:  GetStaticThresholdAttrListTypes(),
					LogAlertConfigFieldCritical: GetStaticThresholdAttrListTypes(),
				},
			},
			[]attr.Value{thresholdObjVal},
		)
		diags.Append(thresholdListDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		ruleObj[LogAlertConfigFieldThreshold] = thresholdList

		// Create rule object value
		ruleObjVal, ruleObjDiags := types.ObjectValue(
			map[string]attr.Type{
				LogAlertConfigFieldMetricName:        types.StringType,
				LogAlertConfigFieldAlertType:         types.StringType,
				LogAlertConfigFieldAggregation:       types.StringType,
				LogAlertConfigFieldThresholdOperator: types.StringType,
				LogAlertConfigFieldThreshold: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							LogAlertConfigFieldWarning:  GetStaticThresholdAttrListTypes(),
							LogAlertConfigFieldCritical: GetStaticThresholdAttrListTypes(),
						},
					},
				},
			},
			ruleObj,
		)
		diags.Append(ruleObjDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		ruleElements[i] = ruleObjVal
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldMetricName:        types.StringType,
				LogAlertConfigFieldAlertType:         types.StringType,
				LogAlertConfigFieldAggregation:       types.StringType,
				LogAlertConfigFieldThresholdOperator: types.StringType,
				LogAlertConfigFieldThreshold: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							LogAlertConfigFieldWarning:  GetStaticThresholdAttrListTypes(),
							LogAlertConfigFieldCritical: GetStaticThresholdAttrListTypes(),
						},
					},
				},
			},
		},
		ruleElements,
	)
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
		GetStaticThresholdAttrTypes(),
		thresholdObj,
	)
	diags.Append(thresholdObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: GetStaticThresholdAttrTypes(),
		},
		[]attr.Value{thresholdObjVal},
	)
}

func (r *logAlertConfigResourceFramework) mapRulesFromState(ctx context.Context, rulesList types.List) ([]restapi.RuleWithThreshold[restapi.LogAlertRule], diag.Diagnostics) {
	var diags diag.Diagnostics

	if rulesList.IsNull() || rulesList.IsUnknown() {
		return []restapi.RuleWithThreshold[restapi.LogAlertRule]{}, diags
	}

	var rulesElements []types.Object
	diags.Append(rulesList.ElementsAs(ctx, &rulesElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(rulesElements) == 0 {
		return []restapi.RuleWithThreshold[restapi.LogAlertRule]{}, diags
	}

	result := make([]restapi.RuleWithThreshold[restapi.LogAlertRule], len(rulesElements))

	for i, ruleObj := range rulesElements {
		var rule struct {
			MetricName        types.String `tfsdk:"metric_name"`
			AlertType         types.String `tfsdk:"alert_type"`
			Aggregation       types.String `tfsdk:"aggregation"`
			ThresholdOperator types.String `tfsdk:"threshold_operator"`
			Threshold         types.List   `tfsdk:"threshold"`
		}

		diags.Append(ruleObj.As(ctx, &rule, basetypes.ObjectAsOptions{})...)
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

		// Map thresholds
		var thresholdMap map[restapi.AlertSeverity]restapi.ThresholdRule
		var thresholdDiags diag.Diagnostics

		if !rule.Threshold.IsNull() && !rule.Threshold.IsUnknown() {
			thresholdMap, thresholdDiags = MapThresholdsFromState(ctx, rule.Threshold)
			diags.Append(thresholdDiags...)
			if diags.HasError() {
				return nil, diags
			}
		} else {
			thresholdMap = make(map[restapi.AlertSeverity]restapi.ThresholdRule)
		}

		// Create rule with threshold
		result[i] = restapi.RuleWithThreshold[restapi.LogAlertRule]{
			ThresholdOperator: thresholdOperator,
			Rule:              logAlertRule,
			Thresholds:        thresholdMap,
		}
	}

	return result, diags
}

// mapThresholdRuleFromState and MapThresholdsFromState have been moved to threshold-mapping-framework.go

// Made with Bob
