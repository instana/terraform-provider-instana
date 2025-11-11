package infralertconfig

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResourceInstanaInfraAlertConfigFramework the name of the terraform-provider-instana resource to manage infrastructure alert configurations
const ResourceInstanaInfraAlertConfigFramework = "infrastructure_alert_config"

const (
	InfraAlertConfigFieldName                  = "name"
	InfraAlertConfigFieldFullName              = "full_name"
	InfraAlertConfigFieldDescription           = "description"
	InfraAlertConfigFieldAlertChannels         = "alert_channels"
	ResourceFieldThresholdRuleWarningSeverity  = "warning"
	ResourceFieldThresholdRuleCriticalSeverity = "critical"
	InfraAlertConfigFieldGroupBy               = "group_by"
	InfraAlertConfigFieldGranularity           = "granularity"
	InfraAlertConfigFieldTagFilter             = "tag_filter"
	InfraAlertConfigFieldEvaluationType        = "evaluation_type"

	InfraAlertConfigFieldRules       = "rules"
	InfraAlertConfigFieldGenericRule = "generic_rule"

	InfraAlertConfigFieldMetricName             = "metric_name"
	InfraAlertConfigFieldEntityType             = "entity_type"
	InfraAlertConfigFieldAggregation            = "aggregation"
	InfraAlertConfigFieldCrossSeriesAggregation = "cross_series_aggregation"
	InfraAlertConfigFieldRegex                  = "regex"
	InfraAlertConfigFieldThresholdOperator      = "threshold_operator"

	InfraAlertConfigFieldTimeThreshold                     = "time_threshold"
	InfraAlertConfigFieldTimeThresholdTimeWindow           = "time_window"
	InfraAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
)

// InfraAlertConfigModel represents the data model for infrastructure alert configuration
type InfraAlertConfigModel struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Description        types.String             `tfsdk:"description"`
	TagFilter          types.String             `tfsdk:"tag_filter"`
	GroupBy            types.List               `tfsdk:"group_by"`
	AlertChannels      types.Object             `tfsdk:"alert_channels"`
	Granularity        types.Int64              `tfsdk:"granularity"`
	TimeThreshold      *InfraTimeThresholdModel `tfsdk:"time_threshold"`
	CustomPayloadField types.List               `tfsdk:"custom_payload_field"`
	Rules              types.Object             `tfsdk:"rules"`
	EvaluationType     types.String             `tfsdk:"evaluation_type"`
}

// InfraAlertChannelsModel represents the alert channels model
type InfraAlertChannelsModel struct {
	Warning  types.List `tfsdk:"warning"`
	Critical types.List `tfsdk:"critical"`
}

// InfraTimeThresholdModel represents the time threshold model
type InfraTimeThresholdModel struct {
	ViolationsInSequence *InfraViolationsInSequenceModel `tfsdk:"violations_in_sequence"`
}

// InfraViolationsInSequenceModel represents the violations in sequence model
type InfraViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}

// InfraCustomPayloadFieldModel represents the custom payload field model
type InfraCustomPayloadFieldModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// InfraRulesModel represents the rules model
type InfraRulesModel struct {
	GenericRule *InfraGenericRuleModel `tfsdk:"generic_rule"`
}

// InfraGenericRuleModel represents the generic rule model
type InfraGenericRuleModel struct {
	MetricName             types.String `tfsdk:"metric_name"`
	EntityType             types.String `tfsdk:"entity_type"`
	Aggregation            types.String `tfsdk:"aggregation"`
	CrossSeriesAggregation types.String `tfsdk:"cross_series_aggregation"`
	Regex                  types.Bool   `tfsdk:"regex"`
	ThresholdOperator      types.String `tfsdk:"threshold_operator"`
	ThresholdRule          types.Object `tfsdk:"threshold"`
}

// InfraThresholdRuleModel represents the threshold rule model
type InfraThresholdRuleModel struct {
	Warning  types.List `tfsdk:"warning"`
	Critical types.List `tfsdk:"critical"`
}

// InfraStaticThresholdModel represents the static threshold model
type InfraStaticThresholdModel struct {
	Value types.Float64 `tfsdk:"value"`
}

// InfraHistoricBaselineThresholdModel represents the historic baseline threshold model
type InfraHistoricBaselineThresholdModel struct {
	DeviationFactor types.Float64 `tfsdk:"deviation_factor"`
	Seasonality     types.String  `tfsdk:"seasonality"`
	Baseline        types.List    `tfsdk:"baseline"`
}

// NewInfraAlertConfigResourceHandleFramework creates a new instance of the infrastructure alert configuration resource
func NewInfraAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.InfraAlertConfig] {
	return &infraAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaInfraAlertConfigFramework,
			Schema: schema.Schema{
				Description: InfraAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: InfraAlertConfigDescID,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description: InfraAlertConfigDescName,
						Required:    true,
					},
					"description": schema.StringAttribute{
						Description: InfraAlertConfigDescDescription,
						Optional:    true,
					},
					"tag_filter": schema.StringAttribute{
						Description: InfraAlertConfigDescTagFilter,
						Optional:    true,
					},
					"group_by": schema.ListAttribute{
						Description: InfraAlertConfigDescGroupBy,
						Optional:    true,
						ElementType: types.StringType,
					},
					"granularity": schema.Int64Attribute{
						Description: InfraAlertConfigDescGranularity,
						Required:    true,
					},

					"evaluation_type": schema.StringAttribute{
						Description: InfraAlertConfigDescEvaluationType,
						Required:    true,
					},
					"custom_payload_field": shared.GetCustomPayloadFieldsSchema(),
					"rules": schema.SingleNestedAttribute{
						Description: InfraAlertConfigDescRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"generic_rule": schema.SingleNestedAttribute{
								Description: InfraAlertConfigDescGenericRule,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Description: InfraAlertConfigDescMetricName,
										Required:    true,
									},
									"entity_type": schema.StringAttribute{
										Description: InfraAlertConfigDescEntityType,
										Required:    true,
									},
									"aggregation": schema.StringAttribute{
										Description: InfraAlertConfigDescAggregation,
										Required:    true,
									},
									"cross_series_aggregation": schema.StringAttribute{
										Description: InfraAlertConfigDescCrossSeriesAggregation,
										Required:    true,
									},
									"regex": schema.BoolAttribute{
										Description: InfraAlertConfigDescRegex,
										Required:    true,
									},
									"threshold_operator": schema.StringAttribute{
										Description: InfraAlertConfigDescThresholdOperator,
										Required:    true,
									},
									"threshold": schema.SingleNestedAttribute{
										Description: InfraAlertConfigDescThreshold,
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"warning":  shared.StaticAndAdaptiveThresholdAttributeSchema(),
											"critical": shared.StaticAndAdaptiveThresholdAttributeSchema(),
										},
									},
								},
							},
						},
					},
					"alert_channels": schema.SingleNestedAttribute{
						Description: InfraAlertConfigDescAlertChannels,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							ResourceFieldThresholdRuleWarningSeverity: schema.ListAttribute{
								Optional:    true,
								Description: InfraAlertConfigDescAlertChannelIDs,
								ElementType: types.StringType,
							},
							ResourceFieldThresholdRuleCriticalSeverity: schema.ListAttribute{
								Optional:    true,
								Description: InfraAlertConfigDescAlertChannelIDs,
								ElementType: types.StringType,
							},
						},
					},
					"time_threshold": schema.SingleNestedAttribute{
						Description: InfraAlertConfigDescTimeThreshold,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"violations_in_sequence": schema.SingleNestedAttribute{
								Description: InfraAlertConfigDescViolationsInSequence,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Required:    true,
										Description: InfraAlertConfigDescTimeWindow,
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 0,
		},
	}
}

type infraAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *infraAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *infraAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.InfraAlertConfig] {
	return api.InfraAlertConfig()
}

func (r *infraAlertConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *infraAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, resource *restapi.InfraAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	model := InfraAlertConfigModel{
		ID:             types.StringValue(resource.ID),
		Name:           types.StringValue(resource.Name),
		Description:    types.StringValue(resource.Description),
		Granularity:    types.Int64Value(int64(resource.Granularity)),
		EvaluationType: types.StringValue(string(resource.EvaluationType)),
	}

	// Map tag filter expression if present
	if resource.TagFilterExpression != nil {
		tagFilterString, err := tagfilter.MapTagFilterToNormalizedString(resource.TagFilterExpression)
		if err != nil {
			diags.AddError(
				InfraAlertConfigErrMappingTagFilter,
				fmt.Sprintf(InfraAlertConfigErrMappingTagFilterMsg, err),
			)
			return diags
		}
		model.TagFilter = util.SetStringPointerToState(tagFilterString)

	} else {
		model.TagFilter = types.StringNull()
	}

	// Map group by if present
	if len(resource.GroupBy) > 0 {
		groupByElements := make([]attr.Value, len(resource.GroupBy))
		for i, groupBy := range resource.GroupBy {
			groupByElements[i] = types.StringValue(groupBy)
		}
		model.GroupBy = types.ListValueMust(types.StringType, groupByElements)
	} else {
		model.GroupBy = types.ListNull(types.StringType)
	}

	// Map alert channels if present
	if len(resource.AlertChannels) > 0 {
		alertChannelsObj := map[string]attr.Value{}

		// Map warning severity
		if warningChannels, ok := resource.AlertChannels[restapi.WarningSeverity]; ok && len(warningChannels) > 0 {
			warningList, warningDiags := types.ListValueFrom(ctx, types.StringType, warningChannels)
			diags.Append(warningDiags...)
			if diags.HasError() {
				return diags
			}
			alertChannelsObj[ResourceFieldThresholdRuleWarningSeverity] = warningList
		} else {
			alertChannelsObj[ResourceFieldThresholdRuleWarningSeverity] = types.ListNull(types.StringType)
		}

		// Map critical severity
		if criticalChannels, ok := resource.AlertChannels[restapi.CriticalSeverity]; ok && len(criticalChannels) > 0 {
			criticalList, criticalDiags := types.ListValueFrom(ctx, types.StringType, criticalChannels)
			diags.Append(criticalDiags...)
			if diags.HasError() {
				return diags
			}
			alertChannelsObj[ResourceFieldThresholdRuleCriticalSeverity] = criticalList
		} else {
			alertChannelsObj[ResourceFieldThresholdRuleCriticalSeverity] = types.ListNull(types.StringType)
		}

		objVal, objDiags := types.ObjectValue(
			map[string]attr.Type{
				ResourceFieldThresholdRuleWarningSeverity:  types.ListType{ElemType: types.StringType},
				ResourceFieldThresholdRuleCriticalSeverity: types.ListType{ElemType: types.StringType},
			},
			alertChannelsObj,
		)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}
		model.AlertChannels = objVal
	} else {
		model.AlertChannels = types.ObjectNull(map[string]attr.Type{
			ResourceFieldThresholdRuleWarningSeverity:  types.ListType{ElemType: types.StringType},
			ResourceFieldThresholdRuleCriticalSeverity: types.ListType{ElemType: types.StringType},
		})
	}

	// Map time threshold if present

	model.TimeThreshold = r.mapTimeThresholdToState(ctx, resource.TimeThreshold)

	// Map custom payload fields if present
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, resource.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayloadField = customPayloadFieldsList

	// Map rules if present
	if len(resource.Rules) > 0 {
		// Create threshold rule model
		thresholdRuleModel := InfraThresholdRuleModel{}

		// Map warning threshold
		warningThreshold, isWarningThresholdPresent := resource.Rules[0].Thresholds[restapi.WarningSeverity]
		warningThresholdList, warningDiags := shared.MapThresholdToState(ctx, isWarningThresholdPresent, &warningThreshold, []string{"static", "adaptiveBaseline"})
		diags.Append(warningDiags...)
		if diags.HasError() {
			return diags
		}
		thresholdRuleModel.Warning = warningThresholdList

		// Map critical threshold
		criticalThreshold, isCriticalThresholdPresent := resource.Rules[0].Thresholds[restapi.CriticalSeverity]
		criticalThresholdList, criticalDiags := shared.MapThresholdToState(ctx, isCriticalThresholdPresent, &criticalThreshold, []string{"static", "adaptiveBaseline"})
		diags.Append(criticalDiags...)
		if diags.HasError() {
			return diags
		}
		thresholdRuleModel.Critical = criticalThresholdList

		// Convert threshold rule model to object
		thresholdRuleObj, thresholdDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrListTypes(),
			shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrListTypes(),
		}, thresholdRuleModel)
		diags.Append(thresholdDiags...)
		if diags.HasError() {
			return diags
		}

		// Create generic rule model
		genericRuleModel := InfraGenericRuleModel{
			MetricName:             types.StringValue(resource.Rules[0].Rule.MetricName),
			EntityType:             types.StringValue(resource.Rules[0].Rule.EntityType),
			Aggregation:            types.StringValue(string(resource.Rules[0].Rule.Aggregation)),
			CrossSeriesAggregation: types.StringValue(string(resource.Rules[0].Rule.CrossSeriesAggregation)),
			Regex:                  types.BoolValue(resource.Rules[0].Rule.Regex),
			ThresholdOperator:      types.StringValue(string(resource.Rules[0].ThresholdOperator)),
			ThresholdRule:          thresholdRuleObj,
		}

		// Create rules model and convert to object
		thresholdObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrListTypes(),
				shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrListTypes(),
			},
		}

		rulesModel := InfraRulesModel{
			GenericRule: &genericRuleModel,
		}

		rulesObj, rulesDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"generic_rule": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":              types.StringType,
					"entity_type":              types.StringType,
					"aggregation":              types.StringType,
					"cross_series_aggregation": types.StringType,
					"regex":                    types.BoolType,
					"threshold_operator":       types.StringType,
					"threshold":                thresholdObjType,
				},
			},
		}, rulesModel)
		diags.Append(rulesDiags...)
		if diags.HasError() {
			return diags
		}

		model.Rules = rulesObj
	} else {
		model.Rules = types.ObjectNull(map[string]attr.Type{
			"generic_rule": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":              types.StringType,
					"entity_type":              types.StringType,
					"aggregation":              types.StringType,
					"cross_series_aggregation": types.StringType,
					"regex":                    types.BoolType,
					"threshold_operator":       types.StringType,
					"threshold": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrListTypes(),
							shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrListTypes(),
						},
					},
				},
			},
		})
	}

	// Set the state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *infraAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.InfraAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model InfraAlertConfigModel

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

	// Map tag filter if present
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() && model.TagFilter.ValueString() != "" {
		tagFilterStr := model.TagFilter.ValueString()
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(tagFilterStr)
		if err != nil {
			diags.AddError(
				InfraAlertConfigErrParsingTagFilter,
				fmt.Sprintf(InfraAlertConfigErrParsingTagFilterMsg, err),
			)
			return nil, diags
		}

		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
	}

	// Map group by if present
	var groupBy []string
	if !model.GroupBy.IsNull() && !model.GroupBy.IsUnknown() {
		diags.Append(model.GroupBy.ElementsAs(ctx, &groupBy, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map alert channels if present
	alertChannels := make(map[restapi.AlertSeverity][]string)
	if !model.AlertChannels.IsNull() && !model.AlertChannels.IsUnknown() {
		var alertChannelsModel struct {
			Warning  types.List `tfsdk:"warning"`
			Critical types.List `tfsdk:"critical"`
		}

		diags.Append(model.AlertChannels.As(ctx, &alertChannelsModel, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		// Map warning severity
		if !alertChannelsModel.Warning.IsNull() && !alertChannelsModel.Warning.IsUnknown() {
			var warningChannels []string
			diags.Append(alertChannelsModel.Warning.ElementsAs(ctx, &warningChannels, false)...)
			if diags.HasError() {
				return nil, diags
			}
			if len(warningChannels) > 0 {
				alertChannels[restapi.WarningSeverity] = warningChannels
			}
		}

		// Map critical severity
		if !alertChannelsModel.Critical.IsNull() && !alertChannelsModel.Critical.IsUnknown() {
			var criticalChannels []string
			diags.Append(alertChannelsModel.Critical.ElementsAs(ctx, &criticalChannels, false)...)
			if diags.HasError() {
				return nil, diags
			}
			if len(criticalChannels) > 0 {
				alertChannels[restapi.CriticalSeverity] = criticalChannels
			}
		}
	}

	// Map time threshold
	var timeThreshold *restapi.InfraTimeThreshold
	if model.TimeThreshold != nil && model.TimeThreshold.ViolationsInSequence != nil {

		v := model.TimeThreshold.ViolationsInSequence
		if !v.TimeWindow.IsNull() && !v.TimeWindow.IsUnknown() {
			timeThreshold = &restapi.InfraTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: v.TimeWindow.ValueInt64(),
			}
		}

	}

	// Map custom payload fields if present
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadField.IsNull() && !model.CustomPayloadField.IsUnknown() {
		var customPayloadFieldModels []InfraCustomPayloadFieldModel
		diags.Append(model.CustomPayloadField.ElementsAs(ctx, &customPayloadFieldModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, field := range customPayloadFieldModels {
			customerPayloadFields = append(customerPayloadFields, restapi.CustomPayloadField[any]{
				Key:   field.Key.ValueString(),
				Value: field.Value.ValueString(),
			})
		}
	}

	// Map rules if present
	var rules []restapi.RuleWithThreshold[restapi.InfraAlertRule]
	if !model.Rules.IsNull() && !model.Rules.IsUnknown() {
		var rulesModel InfraRulesModel
		diags.Append(model.Rules.As(ctx, &rulesModel, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		// Map generic rule
		if rulesModel.GenericRule != nil {
			genericRuleModel := rulesModel.GenericRule

			// Create rule with threshold
			ruleWithThreshold := restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ThresholdOperator: restapi.ThresholdOperator(genericRuleModel.ThresholdOperator.ValueString()),
				Rule: restapi.InfraAlertRule{
					AlertType:              "genericRule",
					MetricName:             genericRuleModel.MetricName.ValueString(),
					EntityType:             genericRuleModel.EntityType.ValueString(),
					Aggregation:            restapi.Aggregation(genericRuleModel.Aggregation.ValueString()),
					CrossSeriesAggregation: restapi.Aggregation(genericRuleModel.CrossSeriesAggregation.ValueString()),
					Regex:                  genericRuleModel.Regex.ValueBool(),
				},
				Thresholds: make(map[restapi.AlertSeverity]restapi.ThresholdRule),
			}

			// Map thresholds
			if !genericRuleModel.ThresholdRule.IsNull() && !genericRuleModel.ThresholdRule.IsUnknown() {
				var thresholdRuleModel InfraThresholdRuleModel
				diags.Append(genericRuleModel.ThresholdRule.As(ctx, &thresholdRuleModel, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}

				// Map warning threshold
				if !thresholdRuleModel.Warning.IsNull() && !thresholdRuleModel.Warning.IsUnknown() {
					warningThresholds, warningDiags := shared.MapThresholdRuleFromState(ctx, thresholdRuleModel.Warning)
					diags.Append(warningDiags...)
					if diags.HasError() {
						return nil, diags
					}
					if warningThresholds != nil {
						ruleWithThreshold.Thresholds[restapi.WarningSeverity] = *warningThresholds
					}
				}

				// Map critical threshold
				if !thresholdRuleModel.Critical.IsNull() && !thresholdRuleModel.Critical.IsUnknown() {
					criticalThresholds, criticalDiags := shared.MapThresholdRuleFromState(ctx, thresholdRuleModel.Critical)
					diags.Append(criticalDiags...)
					if diags.HasError() {
						return nil, diags
					}
					if criticalThresholds != nil {
						ruleWithThreshold.Thresholds[restapi.CriticalSeverity] = *criticalThresholds
					}
				}
			}
			rules = append(rules, ruleWithThreshold)
		}
	}

	// Create the API object
	return &restapi.InfraAlertConfig{
		ID:                    id,
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		TagFilterExpression:   tagFilter,
		GroupBy:               groupBy,
		AlertChannels:         alertChannels,
		Granularity:           restapi.Granularity(model.Granularity.ValueInt64()),
		TimeThreshold:         timeThreshold,
		CustomerPayloadFields: customerPayloadFields,
		Rules:                 rules,
		EvaluationType:        restapi.InfraAlertEvaluationType(model.EvaluationType.ValueString()),
	}, diags
}

func (r *infraAlertConfigResourceFramework) mapTimeThresholdToState(ctx context.Context, api *restapi.InfraTimeThreshold) *InfraTimeThresholdModel {
	if api == nil {
		return nil
	}

	if api.Type != "violationsInSequence" {
		// unsupported type — ignore or handle others
		return nil
	}

	return &InfraTimeThresholdModel{
		ViolationsInSequence: &InfraViolationsInSequenceModel{
			TimeWindow: types.Int64Value(api.TimeWindow),
		},
	}
}
