package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaInfraAlertConfigFramework the name of the terraform-provider-instana resource to manage infrastructure alert configurations
const ResourceInstanaInfraAlertConfigFramework = "infrastructure_alert_config"

// InfraAlertConfigModel represents the data model for infrastructure alert configuration
type InfraAlertConfigModel struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Description        types.String             `tfsdk:"description"`
	TagFilter          types.String             `tfsdk:"tag_filter"`
	GroupBy            types.List               `tfsdk:"group_by"`
	AlertChannels      types.List               `tfsdk:"alert_channels"`
	Granularity        types.Int64              `tfsdk:"granularity"`
	TimeThreshold      *InfraTimeThresholdModel `tfsdk:"time_threshold"`
	CustomPayloadField types.List               `tfsdk:"custom_payload_field"`
	Rules              types.List               `tfsdk:"rules"`
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
	ThresholdRule          types.List   `tfsdk:"threshold"`
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
func NewInfraAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.InfraAlertConfig] {
	return &infraAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaInfraAlertConfigFramework,
			Schema: schema.Schema{
				Description: "This resource represents an infrastructure alert configuration in Instana",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the infrastructure alert configuration",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description: "The name of the infrastructure alert configuration",
						Required:    true,
					},
					"description": schema.StringAttribute{
						Description: "The description of the infrastructure alert configuration",
						Optional:    true,
					},
					"tag_filter": schema.StringAttribute{
						Description: "The tag filter expression for the infrastructure alert configuration",
						Optional:    true,
					},
					"group_by": schema.ListAttribute{
						Description: "The list of tags to group by",
						Optional:    true,
						ElementType: types.StringType,
					},
					"granularity": schema.Int64Attribute{
						Description: "The granularity of the infrastructure alert configuration",
						Required:    true,
					},

					"evaluation_type": schema.StringAttribute{
						Description: "The evaluation type of the infrastructure alert configuration",
						Required:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"rules": schema.ListNestedBlock{
						Description: "The rules configuration",
						//Optional:    true,
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"generic_rule": schema.SingleNestedBlock{
									Description: "The generic rule configuration",

									Attributes: map[string]schema.Attribute{
										"metric_name": schema.StringAttribute{
											Description: "The metric name for the generic rule",
											Required:    true,
										},
										"entity_type": schema.StringAttribute{
											Description: "The entity type for the generic rule",
											Required:    true,
										},
										"aggregation": schema.StringAttribute{
											Description: "The aggregation for the generic rule",
											Required:    true,
										},
										"cross_series_aggregation": schema.StringAttribute{
											Description: "The cross series aggregation for the generic rule",
											Required:    true,
										},
										"regex": schema.BoolAttribute{
											Description: "Whether regex is enabled for the generic rule",
											Required:    true,
										},
										"threshold_operator": schema.StringAttribute{
											Description: "The threshold operator for the generic rule",
											Required:    true,
										},
									},
									Blocks: map[string]schema.Block{
										"threshold": schema.ListNestedBlock{
											Description: "Threshold configuration for different severity levels",
											NestedObject: schema.NestedBlockObject{
												Blocks: map[string]schema.Block{
													"warning":  StaticAndAdaptiveThresholdBlockSchema(),
													"critical": StaticAndAdaptiveThresholdBlockSchema(),
												},
											},
											Validators: []validator.List{
												listvalidator.SizeAtMost(1),
											},
										},
									},
								},
							},
						},
					},
					"custom_payload_field": GetCustomPayloadFieldsSchema(),
					"time_threshold": schema.SingleNestedBlock{
						Description: "Indicates the type of violation of the defined threshold.",
						Blocks: map[string]schema.Block{
							"violations_in_sequence": schema.SingleNestedBlock{
								Description: "Time threshold base on violations in sequence",
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Required:    true,
										Description: "The time window if the time threshold",
									},
								},
							},
						},
					},
					"alert_channels": schema.ListNestedBlock{
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
				},
			},
			SchemaVersion: 0,
		},
	}
}

type infraAlertConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *infraAlertConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
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
				"Error mapping tag filter",
				fmt.Sprintf("Failed to map tag filter: %s", err),
			)
			return diags
		}
		model.TagFilter = setStringPointerToState(tagFilterString)

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
	alertChannelsList, alertChannelsDiags := MapAlertChannelsToState(ctx, resource.AlertChannels)
	if alertChannelsDiags.HasError() {
		diags.Append(alertChannelsDiags...)
		return diags
	}
	model.AlertChannels = alertChannelsList

	// Map time threshold if present

	model.TimeThreshold = r.mapTimeThresholdToState(ctx, resource.TimeThreshold)

	// Map custom payload fields if present
	customPayloadFieldsList, payloadDiags := CustomPayloadFieldsToTerraform(ctx, resource.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayloadField = customPayloadFieldsList

	// Map rules if present
	if len(resource.Rules) > 0 {
		// Create generic rule model
		genericRuleModel := InfraGenericRuleModel{
			MetricName:             types.StringValue(resource.Rules[0].Rule.MetricName),
			EntityType:             types.StringValue(resource.Rules[0].Rule.EntityType),
			Aggregation:            types.StringValue(string(resource.Rules[0].Rule.Aggregation)),
			CrossSeriesAggregation: types.StringValue(string(resource.Rules[0].Rule.CrossSeriesAggregation)),
			Regex:                  types.BoolValue(resource.Rules[0].Rule.Regex),
			ThresholdOperator:      types.StringValue(string(resource.Rules[0].ThresholdOperator)),
		}

		// Create threshold rule model
		thresholdRuleModel := InfraThresholdRuleModel{}

		// Map warning threshold
		warningThreshold, isWarningThresholdPresent := resource.Rules[0].Thresholds[restapi.WarningSeverity]
		warningThresholdList, warningDiags := MapThresholdToState(ctx, isWarningThresholdPresent, &warningThreshold, []string{"static", "adaptiveBaseline"})
		diags.Append(warningDiags...)
		if diags.HasError() {
			return diags
		}
		thresholdRuleModel.Warning = warningThresholdList

		// Map critical threshold
		criticalThreshold, isCriticalThresholdPresent := resource.Rules[0].Thresholds[restapi.CriticalSeverity]
		criticalThresholdList, criticalDiags := MapThresholdToState(ctx, isCriticalThresholdPresent, &criticalThreshold, []string{"static", "adaptiveBaseline"})
		diags.Append(criticalDiags...)
		if diags.HasError() {
			return diags
		}

		thresholdRuleModel.Critical = criticalThresholdList

		// Convert threshold rule model to object
		thresholdRuleObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
			LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
		}, thresholdRuleModel)
		if diags.HasError() {
			return diags
		}

		// Set threshold rule in generic rule model
		genericRuleModel.ThresholdRule = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
				LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
			},
		}, []attr.Value{thresholdRuleObj})

		// Convert generic rule model to object
		thresholdListType := types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
				LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
			},
		}}
		// genericRuleObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		// 	"metric_name":              types.StringType,
		// 	"entity_type":              types.StringType,
		// 	"aggregation":              types.StringType,
		// 	"cross_series_aggregation": types.StringType,
		// 	"regex":                    types.BoolType,
		// 	"threshold_operator":       types.StringType,
		// 	"threshold":                thresholdListType,
		// }, genericRuleModel)
		// if diags.HasError() {
		// 	return diags
		// }

		// Create rules model

		rulesModel := InfraRulesModel{
			GenericRule: &genericRuleModel,
		}

		// rulesModel := InfraRulesModel{
		// 	GenericRule: types.ListValueMust(types.ObjectType{
		// 		AttrTypes: map[string]attr.Type{
		// 			"metric_name":              types.StringType,
		// 			"entity_type":              types.StringType,
		// 			"aggregation":              types.StringType,
		// 			"cross_series_aggregation": types.StringType,
		// 			"regex":                    types.BoolType,
		// 			"threshold_operator":       types.StringType,
		// 			"threshold":                thresholdListType,
		// 		},
		// 	}, []attr.Value{genericRuleObj}),
		// }

		// Convert rules model to object
		rulesObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"generic_rule": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":              types.StringType,
					"entity_type":              types.StringType,
					"aggregation":              types.StringType,
					"cross_series_aggregation": types.StringType,
					"regex":                    types.BoolType,
					"threshold_operator":       types.StringType,
					"threshold":                thresholdListType,
				},
			},
		}, rulesModel)
		if diags.HasError() {
			return diags
		}

		// Set rules in model
		model.Rules = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"generic_rule": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name":              types.StringType,
						"entity_type":              types.StringType,
						"aggregation":              types.StringType,
						"cross_series_aggregation": types.StringType,
						"regex":                    types.BoolType,
						"threshold_operator":       types.StringType,
						"threshold":                thresholdListType,
					},
				},
			},
		}, []attr.Value{rulesObj})
	} else {
		model.Rules = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"generic_rule": types.ObjectType{},
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
				"Error parsing tag filter",
				fmt.Sprintf("Failed to parse tag filter: %s", err),
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
	if !model.AlertChannels.IsNull() {
		var alertChannelsDiags diag.Diagnostics
		alertChannels, alertChannelsDiags = MapAlertChannelsFromState(ctx, model.AlertChannels)
		if alertChannelsDiags.HasError() {
			diags.Append(alertChannelsDiags...)
			return nil, diags
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
		var rulesModels []InfraRulesModel
		diags.Append(model.Rules.ElementsAs(ctx, &rulesModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(rulesModels) > 0 {
			rulesModel := rulesModels[0]

			// Map generic rule
			if rulesModel.GenericRule != nil {

				genericRuleModel := rulesModel.GenericRule

				//if len(genericRuleModels) > 0 {
				//	genericRuleModel := genericRuleModels[0]

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
				var thresholdDiags diag.Diagnostics
				if !genericRuleModel.ThresholdRule.IsNull() && !genericRuleModel.ThresholdRule.IsUnknown() {
					ruleWithThreshold.Thresholds, thresholdDiags = MapThresholdsFromState(ctx, genericRuleModel.ThresholdRule)
					diags.Append(thresholdDiags...)
					if diags.HasError() {
						return nil, diags
					}
				}
				rules = append(rules, ruleWithThreshold)
				//}
			}
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
