package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaInfraAlertConfigFramework the name of the terraform-provider-instana resource to manage infrastructure alert configurations
const ResourceInstanaInfraAlertConfigFramework = "infrastructure_alert_config"

// InfraAlertConfigModel represents the data model for infrastructure alert configuration
type InfraAlertConfigModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	TagFilter           types.String `tfsdk:"tag_filter"`
	GroupBy             types.List   `tfsdk:"group_by"`
	AlertChannels       types.List   `tfsdk:"alert_channels"`
	Granularity         types.Int64  `tfsdk:"granularity"`
	TimeThreshold       types.List   `tfsdk:"time_threshold"`
	CustomPayloadFields types.List   `tfsdk:"custom_payload_fields"`
	Rules               types.List   `tfsdk:"rules"`
	EvaluationType      types.String `tfsdk:"evaluation_type"`
}

// InfraAlertChannelsModel represents the alert channels model
type InfraAlertChannelsModel struct {
	Warning  types.List `tfsdk:"warning"`
	Critical types.List `tfsdk:"critical"`
}

// InfraTimeThresholdModel represents the time threshold model
type InfraTimeThresholdModel struct {
	ViolationsInSequence types.List `tfsdk:"violations_in_sequence"`
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
	GenericRule types.List `tfsdk:"generic_rule"`
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
	Warning  types.Object `tfsdk:"warning"`
	Critical types.Object `tfsdk:"critical"`
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
					"alert_channels": schema.ListNestedAttribute{
						Description: "The alert channels configuration",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"warning": schema.ListAttribute{
									Description: "The list of warning alert channels",
									Optional:    true,
									ElementType: types.StringType,
								},
								"critical": schema.ListAttribute{
									Description: "The list of critical alert channels",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
					"granularity": schema.Int64Attribute{
						Description: "The granularity of the infrastructure alert configuration",
						Required:    true,
					},
					"time_threshold": schema.ListNestedAttribute{
						Description: "The time threshold configuration",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"violations_in_sequence": schema.ListNestedAttribute{
									Description: "The violations in sequence configuration",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"time_window": schema.Int64Attribute{
												Description: "The time window for violations in sequence",
												Required:    true,
											},
										},
									},
								},
							},
						},
					},
					"custom_payload_fields": schema.ListNestedAttribute{
						Description: "The custom payload fields",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: "The key of the custom payload field",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value of the custom payload field",
									Required:    true,
								},
							},
						},
					},
					"rules": schema.ListNestedAttribute{
						Description: "The rules configuration",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"generic_rule": schema.ListNestedAttribute{
									Description: "The generic rule configuration",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
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
											"threshold": schema.ListNestedAttribute{
												Description: "The threshold configuration for the generic rule",
												Required:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"warning": schema.ObjectAttribute{
															Description: "The warning threshold configuration",
															Optional:    true,
															AttributeTypes: map[string]attr.Type{
																"static": types.ListType{
																	ElemType: types.ObjectType{
																		AttrTypes: map[string]attr.Type{
																			"value": types.Float64Type,
																		},
																	},
																},
																"historic_baseline": types.ListType{
																	ElemType: types.ObjectType{
																		AttrTypes: map[string]attr.Type{
																			"deviation_factor": types.Float64Type,
																			"seasonality":      types.StringType,
																			"baseline": types.ListType{
																				ElemType: types.ObjectType{
																					AttrTypes: map[string]attr.Type{
																						"values": types.ListType{
																							ElemType: types.Float64Type,
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
														"critical": schema.ObjectAttribute{
															Description: "The critical threshold configuration",
															Optional:    true,
															AttributeTypes: map[string]attr.Type{
																"static": types.ListType{
																	ElemType: types.ObjectType{
																		AttrTypes: map[string]attr.Type{
																			"value": types.Float64Type,
																		},
																	},
																},
																"historic_baseline": types.ListType{
																	ElemType: types.ObjectType{
																		AttrTypes: map[string]attr.Type{
																			"deviation_factor": types.Float64Type,
																			"seasonality":      types.StringType,
																			"baseline": types.ListType{
																				ElemType: types.ObjectType{
																					AttrTypes: map[string]attr.Type{
																						"values": types.ListType{
																							ElemType: types.Float64Type,
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"evaluation_type": schema.StringAttribute{
						Description: "The evaluation type of the infrastructure alert configuration",
						Required:    true,
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

func (r *infraAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, resource *restapi.InfraAlertConfig) diag.Diagnostics {
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
		model.TagFilter = types.StringValue(*tagFilterString)
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
		alertChannelsModel := InfraAlertChannelsModel{}

		// Map warning alert channels
		if warningChannels, ok := resource.AlertChannels[restapi.WarningSeverity]; ok && len(warningChannels) > 0 {
			warningElements := make([]attr.Value, len(warningChannels))
			for i, channel := range warningChannels {
				warningElements[i] = types.StringValue(channel)
			}
			alertChannelsModel.Warning = types.ListValueMust(types.StringType, warningElements)
		} else {
			alertChannelsModel.Warning = types.ListNull(types.StringType)
		}

		// Map critical alert channels
		if criticalChannels, ok := resource.AlertChannels[restapi.CriticalSeverity]; ok && len(criticalChannels) > 0 {
			criticalElements := make([]attr.Value, len(criticalChannels))
			for i, channel := range criticalChannels {
				criticalElements[i] = types.StringValue(channel)
			}
			alertChannelsModel.Critical = types.ListValueMust(types.StringType, criticalElements)
		} else {
			alertChannelsModel.Critical = types.ListNull(types.StringType)
		}

		// Convert alert channels model to object
		alertChannelsObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"warning":  types.ListType{ElemType: types.StringType},
			"critical": types.ListType{ElemType: types.StringType},
		}, alertChannelsModel)
		if diags.HasError() {
			return diags
		}

		model.AlertChannels = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"warning":  types.ListType{ElemType: types.StringType},
				"critical": types.ListType{ElemType: types.StringType},
			},
		}, []attr.Value{alertChannelsObj})
	} else {
		model.AlertChannels = types.ListNull(types.ObjectType{})
	}

	// Map time threshold if present
	if resource.TimeThreshold.Type != "" {
		timeThresholdModel := InfraTimeThresholdModel{}

		if resource.TimeThreshold.Type == "violationsInSequence" {
			violationsInSequenceModel := InfraViolationsInSequenceModel{
				TimeWindow: types.Int64Value(resource.TimeThreshold.TimeWindow),
			}

			// Convert violations in sequence model to object
			violationsInSequenceObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"time_window": types.Int64Type,
			}, violationsInSequenceModel)
			if diags.HasError() {
				return diags
			}

			timeThresholdModel.ViolationsInSequence = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"time_window": types.Int64Type,
				},
			}, []attr.Value{violationsInSequenceObj})
		} else {
			timeThresholdModel.ViolationsInSequence = types.ListNull(types.ObjectType{})
		}

		// Convert time threshold model to object
		timeThresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"violations_in_sequence": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"time_window": types.Int64Type,
				},
			}},
		}, timeThresholdModel)
		if diags.HasError() {
			return diags
		}

		model.TimeThreshold = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"violations_in_sequence": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"time_window": types.Int64Type,
					},
				}},
			},
		}, []attr.Value{timeThresholdObj})
	} else {
		model.TimeThreshold = types.ListNull(types.ObjectType{})
	}

	// Map custom payload fields if present
	if len(resource.CustomerPayloadFields) > 0 {
		customPayloadFieldElements := make([]attr.Value, len(resource.CustomerPayloadFields))
		for i, field := range resource.CustomerPayloadFields {
			customPayloadFieldModel := InfraCustomPayloadFieldModel{
				Key:   types.StringValue(field.Key),
				Value: types.StringValue(fmt.Sprintf("%v", field.Value)),
			}

			// Convert custom payload field model to object
			customPayloadFieldObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			}, customPayloadFieldModel)
			if diags.HasError() {
				return diags
			}

			customPayloadFieldElements[i] = customPayloadFieldObj
		}

		model.CustomPayloadFields = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}, customPayloadFieldElements)
	} else {
		model.CustomPayloadFields = types.ListNull(types.ObjectType{})
	}

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

		// Map warning threshold if present
		if warningThreshold, ok := resource.Rules[0].Thresholds[restapi.WarningSeverity]; ok {
			if warningThreshold.Type == "staticThreshold" {
				staticThresholdModel := InfraStaticThresholdModel{
					Value: types.Float64Value(*warningThreshold.Value),
				}

				// Convert static threshold model to object
				staticThresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
					"value": types.Float64Type,
				}, staticThresholdModel)
				if diags.HasError() {
					return diags
				}

				warningThresholdObj := types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"value": types.Float64Type,
					},
				}, []attr.Value{staticThresholdObj})

				thresholdRuleModel.Warning = types.ObjectValueMust(map[string]attr.Type{
					"static": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"value": types.Float64Type,
						},
					}},
				}, map[string]attr.Value{
					"static": warningThresholdObj,
				})
			} else if warningThreshold.Type == "historicBaseline" {
				historicBaselineModel := InfraHistoricBaselineThresholdModel{
					DeviationFactor: types.Float64Value(float64(*warningThreshold.DeviationFactor)),
					Seasonality:     types.StringValue(string(*warningThreshold.Seasonality)),
				}

				// Map baseline if present
				var baselineList []attr.Value
				if warningThreshold.Baseline != nil {
					for _, baseline := range *warningThreshold.Baseline {
						valuesElements := make([]attr.Value, len(baseline))
						for i, value := range baseline {
							valuesElements[i] = types.Float64Value(value)
						}

						baselineObj, baselineDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						}, map[string]attr.Value{
							"values": types.ListValueMust(types.Float64Type, valuesElements),
						})
						if baselineDiags.HasError() {
							diags.Append(baselineDiags...)
							return diags
						}

						baselineList = append(baselineList, baselineObj)
					}
					historicBaselineModel.Baseline = types.ListValueMust(types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						},
					}, baselineList)
				} else {
					historicBaselineModel.Baseline = types.ListNull(types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						},
					})
				}

				historicBaselineObj, criticalDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
					"deviation_factor": types.Float64Type,
					"seasonality":      types.StringType,
					"baseline": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						},
					}},
				}, historicBaselineModel)
				if criticalDiags.HasError() {
					diags.Append(criticalDiags...)
					return diags
				}

				warningThresholdObj := types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"deviation_factor": types.Float64Type,
						"seasonality":      types.StringType,
						"baseline": types.ListType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"values": types.ListType{ElemType: types.Float64Type},
							},
						}},
					},
				}, []attr.Value{historicBaselineObj})

				thresholdRuleModel.Warning = types.ObjectValueMust(map[string]attr.Type{
					"historic_baseline": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"deviation_factor": types.Float64Type,
							"seasonality":      types.StringType,
							"baseline": types.ListType{ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"values": types.ListType{ElemType: types.Float64Type},
								},
							}},
						},
					}},
				}, map[string]attr.Value{
					"historic_baseline": warningThresholdObj,
				})
			} else {
				thresholdRuleModel.Warning = types.ObjectNull(map[string]attr.Type{})
			}
		} else {
			thresholdRuleModel.Warning = types.ObjectNull(map[string]attr.Type{})
		}

		// Map critical threshold if present
		if criticalThreshold, ok := resource.Rules[0].Thresholds[restapi.CriticalSeverity]; ok {
			var criticalThresholdObj types.List

			if criticalThreshold.Type == "staticThreshold" {
				staticThresholdModel := InfraStaticThresholdModel{
					Value: types.Float64Value(*criticalThreshold.Value),
				}

				// Convert static threshold model to object
				staticThresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
					"value": types.Float64Type,
				}, staticThresholdModel)
				if diags.HasError() {
					return diags
				}

				criticalThresholdObj = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"value": types.Float64Type,
					},
				}, []attr.Value{staticThresholdObj})
			} else if criticalThreshold.Type == "historicBaseline" {
				historicBaselineModel := InfraHistoricBaselineThresholdModel{
					DeviationFactor: types.Float64Value(float64(*criticalThreshold.DeviationFactor)),
					Seasonality:     types.StringValue(string(*criticalThreshold.Seasonality)),
				}

				// Map baseline if present
				var baselineList []attr.Value
				if criticalThreshold.Baseline != nil {
					for _, baseline := range *criticalThreshold.Baseline {
						valuesElements := make([]attr.Value, len(baseline))
						for i, value := range baseline {
							valuesElements[i] = types.Float64Value(value)
						}

						baselineObj, baselineDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						}, map[string]attr.Value{
							"values": types.ListValueMust(types.Float64Type, valuesElements),
						})
						if baselineDiags.HasError() {
							diags.Append(baselineDiags...)
							return diags
						}

						baselineList = append(baselineList, baselineObj)
					}
					historicBaselineModel.Baseline = types.ListValueMust(types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						},
					}, baselineList)
				} else {
					historicBaselineModel.Baseline = types.ListNull(types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						},
					})
				}

				historicBaselineObj, criticalDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
					"deviation_factor": types.Float64Type,
					"seasonality":      types.StringType,
					"baseline": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"values": types.ListType{ElemType: types.Float64Type},
						},
					}},
				}, historicBaselineModel)
				if criticalDiags.HasError() {
					diags.Append(criticalDiags...)
					return diags
				}

				criticalThresholdObj = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"deviation_factor": types.Float64Type,
						"seasonality":      types.StringType,
						"baseline": types.ListType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"values": types.ListType{ElemType: types.Float64Type},
							},
						}},
					},
				}, []attr.Value{historicBaselineObj})
			} else {
				criticalThresholdObj = types.ListNull(types.ObjectType{})
			}

			if criticalThreshold.Type == "staticThreshold" {
				thresholdRuleModel.Critical = types.ObjectValueMust(map[string]attr.Type{
					"static": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"value": types.Float64Type,
						},
					}},
				}, map[string]attr.Value{
					"static": criticalThresholdObj,
				})
			} else if criticalThreshold.Type == "historicBaseline" {
				thresholdRuleModel.Critical = types.ObjectValueMust(map[string]attr.Type{
					"historic_baseline": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"deviation_factor": types.Float64Type,
							"seasonality":      types.StringType,
							"baseline": types.ListType{ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"values": types.ListType{ElemType: types.Float64Type},
								},
							}},
						},
					}},
				}, map[string]attr.Value{
					"historic_baseline": criticalThresholdObj,
				})
			} else {
				thresholdRuleModel.Critical = types.ObjectNull(map[string]attr.Type{})
			}
		} else {
			thresholdRuleModel.Critical = types.ObjectNull(map[string]attr.Type{})
		}

		// Convert threshold rule model to object
		thresholdRuleObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"warning": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"static":            types.ListType{ElemType: types.ObjectType{}},
					"historic_baseline": types.ListType{ElemType: types.ObjectType{}},
				},
			},
			"critical": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"static":            types.ListType{ElemType: types.ObjectType{}},
					"historic_baseline": types.ListType{ElemType: types.ObjectType{}},
				},
			},
		}, thresholdRuleModel)
		if diags.HasError() {
			return diags
		}

		// Set threshold rule in generic rule model
		genericRuleModel.ThresholdRule = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"warning": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"static":            types.ListType{ElemType: types.ObjectType{}},
						"historic_baseline": types.ListType{ElemType: types.ObjectType{}},
					},
				},
				"critical": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"static":            types.ListType{ElemType: types.ObjectType{}},
						"historic_baseline": types.ListType{ElemType: types.ObjectType{}},
					},
				},
			},
		}, []attr.Value{thresholdRuleObj})

		// Convert generic rule model to object
		genericRuleObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"metric_name":              types.StringType,
			"entity_type":              types.StringType,
			"aggregation":              types.StringType,
			"cross_series_aggregation": types.StringType,
			"regex":                    types.BoolType,
			"threshold_operator":       types.StringType,
			"threshold":                types.ListType{ElemType: types.ObjectType{}},
		}, genericRuleModel)
		if diags.HasError() {
			return diags
		}

		// Create rules model
		rulesModel := InfraRulesModel{
			GenericRule: types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":              types.StringType,
					"entity_type":              types.StringType,
					"aggregation":              types.StringType,
					"cross_series_aggregation": types.StringType,
					"regex":                    types.BoolType,
					"threshold_operator":       types.StringType,
					"threshold":                types.ListType{ElemType: types.ObjectType{}},
				},
			}, []attr.Value{genericRuleObj}),
		}

		// Convert rules model to object
		rulesObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"generic_rule": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":              types.StringType,
					"entity_type":              types.StringType,
					"aggregation":              types.StringType,
					"cross_series_aggregation": types.StringType,
					"regex":                    types.BoolType,
					"threshold_operator":       types.StringType,
					"threshold":                types.ListType{ElemType: types.ObjectType{}},
				},
			}},
		}, rulesModel)
		if diags.HasError() {
			return diags
		}

		// Set rules in model
		model.Rules = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"generic_rule": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name":              types.StringType,
						"entity_type":              types.StringType,
						"aggregation":              types.StringType,
						"cross_series_aggregation": types.StringType,
						"regex":                    types.BoolType,
						"threshold_operator":       types.StringType,
						"threshold":                types.ListType{ElemType: types.ObjectType{}},
					},
				}},
			},
		}, []attr.Value{rulesObj})
	} else {
		model.Rules = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"generic_rule": types.ListType{ElemType: types.ObjectType{}},
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
	if !model.AlertChannels.IsNull() && !model.AlertChannels.IsUnknown() {
		var alertChannelsModels []InfraAlertChannelsModel
		diags.Append(model.AlertChannels.ElementsAs(ctx, &alertChannelsModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(alertChannelsModels) > 0 {
			alertChannelsModel := alertChannelsModels[0]

			// Map warning alert channels
			if !alertChannelsModel.Warning.IsNull() && !alertChannelsModel.Warning.IsUnknown() {
				var warningChannels []string
				diags.Append(alertChannelsModel.Warning.ElementsAs(ctx, &warningChannels, false)...)
				if diags.HasError() {
					return nil, diags
				}
				alertChannels[restapi.WarningSeverity] = warningChannels
			}

			// Map critical alert channels
			if !alertChannelsModel.Critical.IsNull() && !alertChannelsModel.Critical.IsUnknown() {
				var criticalChannels []string
				diags.Append(alertChannelsModel.Critical.ElementsAs(ctx, &criticalChannels, false)...)
				if diags.HasError() {
					return nil, diags
				}
				alertChannels[restapi.CriticalSeverity] = criticalChannels
			}
		}
	}

	// Map time threshold
	var timeThreshold restapi.InfraTimeThreshold
	if !model.TimeThreshold.IsNull() && !model.TimeThreshold.IsUnknown() {
		var timeThresholdModels []InfraTimeThresholdModel
		diags.Append(model.TimeThreshold.ElementsAs(ctx, &timeThresholdModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(timeThresholdModels) > 0 {
			timeThresholdModel := timeThresholdModels[0]

			if !timeThresholdModel.ViolationsInSequence.IsNull() && !timeThresholdModel.ViolationsInSequence.IsUnknown() {
				var violationsInSequenceModels []InfraViolationsInSequenceModel
				diags.Append(timeThresholdModel.ViolationsInSequence.ElementsAs(ctx, &violationsInSequenceModels, false)...)
				if diags.HasError() {
					return nil, diags
				}

				if len(violationsInSequenceModels) > 0 {
					violationsInSequenceModel := violationsInSequenceModels[0]
					timeThreshold = restapi.InfraTimeThreshold{
						Type:       "violationsInSequence",
						TimeWindow: violationsInSequenceModel.TimeWindow.ValueInt64(),
					}
				}
			}
		}
	}

	// Map custom payload fields if present
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		var customPayloadFieldModels []InfraCustomPayloadFieldModel
		diags.Append(model.CustomPayloadFields.ElementsAs(ctx, &customPayloadFieldModels, false)...)
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
			if !rulesModel.GenericRule.IsNull() && !rulesModel.GenericRule.IsUnknown() {
				var genericRuleModels []InfraGenericRuleModel
				diags.Append(rulesModel.GenericRule.ElementsAs(ctx, &genericRuleModels, false)...)
				if diags.HasError() {
					return nil, diags
				}

				if len(genericRuleModels) > 0 {
					genericRuleModel := genericRuleModels[0]

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

					// Map threshold rules
					if !genericRuleModel.ThresholdRule.IsNull() && !genericRuleModel.ThresholdRule.IsUnknown() {
						var thresholdRuleModels []InfraThresholdRuleModel
						diags.Append(genericRuleModel.ThresholdRule.ElementsAs(ctx, &thresholdRuleModels, false)...)
						if diags.HasError() {
							return nil, diags
						}

						if len(thresholdRuleModels) > 0 {
							thresholdRuleModel := thresholdRuleModels[0]

							// Map warning threshold
							if !thresholdRuleModel.Warning.IsNull() && !thresholdRuleModel.Warning.IsUnknown() {
								var warningThresholdMap map[string]attr.Value
								diags.Append(tfsdk.ValueAs(ctx, thresholdRuleModel.Warning, &warningThresholdMap)...)
								if diags.HasError() {
									return nil, diags
								}

								// Check if static threshold
								if staticThresholdValue, ok := warningThresholdMap["static"]; ok && !staticThresholdValue.IsNull() {
									var staticThresholdList []map[string]attr.Value
									diags.Append(tfsdk.ValueAs(ctx, staticThresholdValue, &staticThresholdList)...)
									if diags.HasError() {
										return nil, diags
									}

									if len(staticThresholdList) > 0 {
										var staticThresholdModel InfraStaticThresholdModel
										diags.Append(tfsdk.ValueAs(ctx, staticThresholdList[0]["value"], &staticThresholdModel.Value)...)
										if diags.HasError() {
											return nil, diags
										}

										value := staticThresholdModel.Value.ValueFloat64()
										ruleWithThreshold.Thresholds[restapi.WarningSeverity] = restapi.ThresholdRule{
											Type:  "staticThreshold",
											Value: &value,
										}
									}
								}

								// Check if historic baseline threshold
								if historicBaselineValue, ok := warningThresholdMap["historic_baseline"]; ok && !historicBaselineValue.IsNull() {
									var historicBaselineList []map[string]attr.Value
									diags.Append(tfsdk.ValueAs(ctx, historicBaselineValue, &historicBaselineList)...)
									if diags.HasError() {
										return nil, diags
									}

									if len(historicBaselineList) > 0 {
										var historicBaselineModel InfraHistoricBaselineThresholdModel
										diags.Append(tfsdk.ValueAs(ctx, historicBaselineList[0]["deviation_factor"], &historicBaselineModel.DeviationFactor)...)
										diags.Append(tfsdk.ValueAs(ctx, historicBaselineList[0]["seasonality"], &historicBaselineModel.Seasonality)...)
										if diags.HasError() {
											return nil, diags
										}

										deviationFactor := float32(historicBaselineModel.DeviationFactor.ValueFloat64())
										seasonality := restapi.ThresholdSeasonality(historicBaselineModel.Seasonality.ValueString())

										thresholdRule := restapi.ThresholdRule{
											Type:            "historicBaseline",
											DeviationFactor: &deviationFactor,
											Seasonality:     &seasonality,
										}

										// Map baseline if present
										if baselineValue, ok := historicBaselineList[0]["baseline"]; ok && !baselineValue.IsNull() {
											var baselineList []map[string]attr.Value
											diags.Append(tfsdk.ValueAs(ctx, baselineValue, &baselineList)...)
											if diags.HasError() {
												return nil, diags
											}

											if len(baselineList) > 0 {
												var baselineValues [][]float64
												for _, baseline := range baselineList {
													var values []float64
													diags.Append(tfsdk.ValueAs(ctx, baseline["values"], &values)...)
													if diags.HasError() {
														return nil, diags
													}
													baselineValues = append(baselineValues, values)
												}
												thresholdRule.Baseline = &baselineValues
											}
										}

										ruleWithThreshold.Thresholds[restapi.WarningSeverity] = thresholdRule
									}
								}
							}

							// Map critical threshold
							if !thresholdRuleModel.Critical.IsNull() && !thresholdRuleModel.Critical.IsUnknown() {
								var criticalThresholdMap map[string]attr.Value
								diags.Append(tfsdk.ValueAs(ctx, thresholdRuleModel.Critical, &criticalThresholdMap)...)
								if diags.HasError() {
									return nil, diags
								}

								// Check if static threshold
								if staticThresholdValue, ok := criticalThresholdMap["static"]; ok && !staticThresholdValue.IsNull() {
									var staticThresholdList []map[string]attr.Value
									diags.Append(tfsdk.ValueAs(ctx, staticThresholdValue, &staticThresholdList)...)
									if diags.HasError() {
										return nil, diags
									}

									if len(staticThresholdList) > 0 {
										var staticThresholdModel InfraStaticThresholdModel
										diags.Append(tfsdk.ValueAs(ctx, staticThresholdList[0]["value"], &staticThresholdModel.Value)...)
										if diags.HasError() {
											return nil, diags
										}

										value := staticThresholdModel.Value.ValueFloat64()
										ruleWithThreshold.Thresholds[restapi.CriticalSeverity] = restapi.ThresholdRule{
											Type:  "staticThreshold",
											Value: &value,
										}
									}
								}

								// Check if historic baseline threshold
								if historicBaselineValue, ok := criticalThresholdMap["historic_baseline"]; ok && !historicBaselineValue.IsNull() {
									var historicBaselineList []map[string]attr.Value
									diags.Append(tfsdk.ValueAs(ctx, historicBaselineValue, &historicBaselineList)...)
									if diags.HasError() {
										return nil, diags
									}

									if len(historicBaselineList) > 0 {
										var historicBaselineModel InfraHistoricBaselineThresholdModel
										diags.Append(tfsdk.ValueAs(ctx, historicBaselineList[0]["deviation_factor"], &historicBaselineModel.DeviationFactor)...)
										diags.Append(tfsdk.ValueAs(ctx, historicBaselineList[0]["seasonality"], &historicBaselineModel.Seasonality)...)
										if diags.HasError() {
											return nil, diags
										}

										deviationFactor := float32(historicBaselineModel.DeviationFactor.ValueFloat64())
										seasonality := restapi.ThresholdSeasonality(historicBaselineModel.Seasonality.ValueString())

										thresholdRule := restapi.ThresholdRule{
											Type:            "historicBaseline",
											DeviationFactor: &deviationFactor,
											Seasonality:     &seasonality,
										}

										// Map baseline if present
										if baselineValue, ok := historicBaselineList[0]["baseline"]; ok && !baselineValue.IsNull() {
											var baselineList []map[string]attr.Value
											diags.Append(tfsdk.ValueAs(ctx, baselineValue, &baselineList)...)
											if diags.HasError() {
												return nil, diags
											}

											if len(baselineList) > 0 {
												var baselineValues [][]float64
												for _, baseline := range baselineList {
													var values []float64
													diags.Append(tfsdk.ValueAs(ctx, baseline["values"], &values)...)
													if diags.HasError() {
														return nil, diags
													}
													baselineValues = append(baselineValues, values)
												}
												thresholdRule.Baseline = &baselineValues
											}
										}

										ruleWithThreshold.Thresholds[restapi.CriticalSeverity] = thresholdRule
									}
								}
							}
						}
					}

					rules = append(rules, ruleWithThreshold)
				}
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

// Made with Bob

// Made with Bob
