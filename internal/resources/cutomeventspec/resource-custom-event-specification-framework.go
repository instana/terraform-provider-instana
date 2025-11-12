package cutomeventspec

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
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
)

// NewCustomEventSpecificationResourceHandleFramework creates the resource handle for Custom Event Specification
func NewCustomEventSpecificationResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.CustomEventSpecification] {
	return &customEventSpecificationResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaCustomEventSpecificationFramework,
			Schema:        createCustomEventSpecificationSchema(),
			SchemaVersion: 1,
		},
	}
}

type customEventSpecificationResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *customEventSpecificationResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *customEventSpecificationResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

// createCustomEventSpecificationSchema creates the schema for the custom event specification resource
func createCustomEventSpecificationSchema() schema.Schema {
	return schema.Schema{
		Description: CustomEventSpecificationResourceDescResource,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescID,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescName,
				Required:    true,
			},
			"entity_type": schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescEntityType,
				Required:    true,
			},
			"query": schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescQuery,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"triggering": schema.BoolAttribute{
				Description: CustomEventSpecificationResourceDescTriggering,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescDescription,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"expiration_time": schema.Int64Attribute{
				Description: CustomEventSpecificationResourceDescExpirationTime,
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: CustomEventSpecificationResourceDescEnabled,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"rule_logical_operator": schema.StringAttribute{
				Description: CustomEventSpecificationResourceDescRuleLogicalOperator,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("AND"),
				Validators: []validator.String{
					stringvalidator.OneOf("AND", "OR"),
				},
			},
			"rules": schema.SingleNestedAttribute{
				Description: CustomEventSpecificationResourceDescRules,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"entity_count": schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescEntityCountRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("warning", "critical"),
								},
							},
							"condition_operator": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescConditionOperator,
								Required:    true,
							},
							"condition_value": schema.Float64Attribute{
								Description: CustomEventSpecificationResourceDescConditionValue,
								Required:    true,
							},
						},
					},
					"entity_count_verification": schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescEntityCountVerification,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("warning", "critical"),
								},
							},
							"condition_operator": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescConditionOperator,
								Required:    true,
							},
							"condition_value": schema.Float64Attribute{
								Description: CustomEventSpecificationResourceDescConditionValue,
								Required:    true,
							},
							"matching_entity_type": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityType,
								Required:    true,
							},
							"matching_operator": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingOperator,
								Required:    true,
							},
							"matching_entity_label": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityLabel,
								Required:    true,
							},
						},
					},
					"entity_verification": schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescEntityVerification,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("warning", "critical"),
								},
							},
							"matching_entity_type": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityType,
								Required:    true,
							},
							"matching_operator": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingOperator,
								Required:    true,
							},
							"matching_entity_label": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMatchingEntityLabel,
								Required:    true,
							},
							"offline_duration": schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescOfflineDuration,
								Required:    true,
							},
						},
					},
					"host_availability": schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescHostAvailability,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("warning", "critical"),
								},
							},
							"offline_duration": schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescOfflineDuration,
								Required:    true,
							},
							"close_after": schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescCloseAfter,
								Optional:    true,
							},
							"tag_filter": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescTagFilter,
								Optional:    true,
							},
						},
					},
					"system": schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescSystemRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("warning", "critical"),
								},
							},
							"system_rule_id": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSystemRuleID,
								Required:    true,
							},
						},
					},
					"threshold": schema.SingleNestedAttribute{
						Description: CustomEventSpecificationResourceDescThresholdRules,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"severity": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescSeverity,
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("warning", "critical"),
								},
							},
							"metric_name": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescMetricName,
								Required:    true,
							},
							"rollup": schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescRollup,
								Required:    true,
							},
							"window": schema.Int64Attribute{
								Description: CustomEventSpecificationResourceDescWindow,
								Required:    true,
							},
							"aggregation": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescAggregation,
								Required:    true,
							},
							"condition_operator": schema.StringAttribute{
								Description: CustomEventSpecificationResourceDescConditionOperator,
								Required:    true,
							},
							"condition_value": schema.Float64Attribute{
								Description: CustomEventSpecificationResourceDescConditionValue,
								Required:    true,
							},
							"metric_pattern": schema.SingleNestedAttribute{
								Description: CustomEventSpecificationResourceDescMetricPattern,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"prefix": schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternPrefix,
										Required:    true,
									},
									"postfix": schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternPostfix,
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString(""),
									},
									"placeholder": schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternPlaceholder,
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString(""),
									},
									"operator": schema.StringAttribute{
										Description: CustomEventSpecificationResourceDescMetricPatternOperator,
										Optional:    true,
										Computed:    true,
										Default:     stringdefault.StaticString("EQUALS"),
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
func (r *customEventSpecificationResourceFramework) SetComputedFields(_ context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	// No computed fields to set
	return nil
}

func (r *customEventSpecificationResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, spec *restapi.CustomEventSpecification) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the spec
	model := CustomEventSpecificationModel{
		ID:                  types.StringValue(spec.ID),
		Name:                types.StringValue(spec.Name),
		EntityType:          types.StringValue(spec.EntityType),
		Triggering:          types.BoolValue(spec.Triggering),
		Enabled:             types.BoolValue(spec.Enabled),
		RuleLogicalOperator: types.StringValue(spec.RuleLogicalOperator),
	}

	// Set optional fields
	model.Query = util.SetStringPointerToState(spec.Query)

	model.Description = util.SetStringPointerToState(spec.Description)

	if spec.ExpirationTime != nil {
		model.ExpirationTime = util.SetInt64PointerToState(spec.ExpirationTime)
	} else {
		model.ExpirationTime = types.Int64Null()
	}

	// Process rules
	if len(spec.Rules) > 0 {
		// Create rule objects (single instances, not lists)
		var entityCountRule *EntityCountRuleModel
		var entityCountVerificationRule *EntityCountVerificationRuleModel
		var entityVerificationRule *EntityVerificationRuleModel
		var hostAvailabilityRule *HostAvailabilityRuleModel
		var systemRule *SystemRuleModel
		var thresholdRule *ThresholdRuleModel

		// Process each rule based on its type (take first occurrence of each type)
		for _, rule := range spec.Rules {
			switch rule.DType {
			case restapi.EntityCountRuleType:
				if entityCountRule == nil && rule.ConditionOperator != nil && rule.ConditionValue != nil {
					entityCountRule = &EntityCountRuleModel{
						Severity:          mapIntToSeverityString(rule.Severity),
						ConditionOperator: util.SetStringPointerToState(rule.ConditionOperator),
						ConditionValue:    util.SetFloat64PointerToState(rule.ConditionValue),
					}
				}
			case restapi.EntityCountVerificationRuleType:
				if entityCountVerificationRule == nil && rule.ConditionOperator != nil && rule.ConditionValue != nil &&
					rule.MatchingEntityType != nil && rule.MatchingOperator != nil && rule.MatchingEntityLabel != nil {
					entityCountVerificationRule = &EntityCountVerificationRuleModel{
						Severity:            mapIntToSeverityString(rule.Severity),
						ConditionOperator:   util.SetStringPointerToState(rule.ConditionOperator),
						ConditionValue:      util.SetFloat64PointerToState(rule.ConditionValue),
						MatchingEntityType:  util.SetStringPointerToState(rule.MatchingEntityType),
						MatchingOperator:    util.SetStringPointerToState(rule.MatchingOperator),
						MatchingEntityLabel: util.SetStringPointerToState(rule.MatchingEntityLabel),
					}
				}
			case restapi.EntityVerificationRuleType:
				if entityVerificationRule == nil && rule.MatchingEntityType != nil && rule.MatchingOperator != nil &&
					rule.MatchingEntityLabel != nil && rule.OfflineDuration != nil {
					entityVerificationRule = &EntityVerificationRuleModel{
						Severity:            mapIntToSeverityString(rule.Severity),
						MatchingEntityType:  util.SetStringPointerToState(rule.MatchingEntityType),
						MatchingOperator:    util.SetStringPointerToState(rule.MatchingOperator),
						MatchingEntityLabel: util.SetStringPointerToState(rule.MatchingEntityLabel),
						OfflineDuration:     util.SetInt64PointerToState(rule.OfflineDuration),
					}
				}
			case restapi.HostAvailabilityRuleType:
				if hostAvailabilityRule == nil && rule.OfflineDuration != nil {
					hr := HostAvailabilityRuleModel{
						Severity:        mapIntToSeverityString(rule.Severity),
						OfflineDuration: util.SetInt64PointerToState(rule.OfflineDuration),
						TagFilter:       types.StringValue(""), // Default empty string
					}

					if rule.CloseAfter != nil {
						hr.CloseAfter = util.SetInt64PointerToState(rule.CloseAfter)
					} else {
						hr.CloseAfter = types.Int64Null()
					}

					// Handle tag filter conversion
					if rule.TagFilter != nil {
						// Convert tag filter to string representation
						normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(rule.TagFilter)
						if err == nil && normalizedTagFilterString != nil {
							hr.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
						}
					}

					hostAvailabilityRule = &hr
				}
			case restapi.SystemRuleType:
				if systemRule == nil && rule.SystemRuleID != nil {
					systemRule = &SystemRuleModel{
						Severity:     mapIntToSeverityString(rule.Severity),
						SystemRuleID: util.SetStringPointerToState(rule.SystemRuleID),
					}
				}
			case restapi.ThresholdRuleType:
				if thresholdRule == nil && rule.MetricName != nil && rule.Rollup != nil && rule.Window != nil &&
					rule.Aggregation != nil && rule.ConditionOperator != nil && rule.ConditionValue != nil {
					tr := ThresholdRuleModel{
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
						tr.MetricPattern = &MetricPatternModel{
							Prefix:      types.StringValue(rule.MetricPattern.Prefix),
							Operator:    types.StringValue(rule.MetricPattern.Operator),
							Postfix:     util.SetStringPointerToState(rule.MetricPattern.Postfix),
							Placeholder: util.SetStringPointerToState(rule.MetricPattern.Placeholder),
						}
					}
					// If nil, leave tr.MetricPattern as nil

					thresholdRule = &tr
				}
			}
		}

		// Create the rules model with pointer fields
		model.Rules = &RulesModel{
			EntityCount:             entityCountRule,
			EntityCountVerification: entityCountVerificationRule,
			EntityVerification:      entityVerificationRule,
			HostAvailability:        hostAvailabilityRule,
			System:                  systemRule,
			Threshold:               thresholdRule,
		}
	}
	// If no rules, leave model.Rules as nil

	// Set the entire model to state
	diags = state.Set(ctx, model)
	return diags
}

// mapIntToSeverityString maps the severity integer to a string value
func mapIntToSeverityString(severity int) types.String {
	switch severity {
	case 5:
		return types.StringValue("warning")
	case 10:
		return types.StringValue("critical")
	default:
		return types.StringValue("warning") // Default to warning
	}
}

func (r *customEventSpecificationResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.CustomEventSpecification, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model CustomEventSpecificationModel

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

	// Map optional fields
	var query *string
	if !model.Query.IsNull() && model.Query.ValueString() != "" {
		queryStr := model.Query.ValueString()
		query = &queryStr
	}

	var description *string
	if !model.Description.IsNull() && model.Description.ValueString() != "" {
		descStr := model.Description.ValueString()
		description = &descStr
	}

	var expirationTime *int
	if !model.ExpirationTime.IsNull() {
		expTime := int(model.ExpirationTime.ValueInt64())
		expirationTime = &expTime
	}

	// Map rules
	var rules []restapi.RuleSpecification

	// Check if rules are defined
	if model.Rules != nil {
		// Process entity count rule
		if model.Rules.EntityCount != nil {
			severity := mapSeverityToInt(model.Rules.EntityCount.Severity.ValueString())
			conditionOperator := model.Rules.EntityCount.ConditionOperator.ValueString()
			conditionValue := model.Rules.EntityCount.ConditionValue.ValueFloat64()

			rules = append(rules, restapi.RuleSpecification{
				DType:             restapi.EntityCountRuleType,
				Severity:          severity,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
			})
		}

		// Process entity count verification rule
		if model.Rules.EntityCountVerification != nil {
			severity := mapSeverityToInt(model.Rules.EntityCountVerification.Severity.ValueString())
			conditionOperator := model.Rules.EntityCountVerification.ConditionOperator.ValueString()
			conditionValue := model.Rules.EntityCountVerification.ConditionValue.ValueFloat64()
			matchingEntityType := model.Rules.EntityCountVerification.MatchingEntityType.ValueString()
			matchingOperator := model.Rules.EntityCountVerification.MatchingOperator.ValueString()
			matchingEntityLabel := model.Rules.EntityCountVerification.MatchingEntityLabel.ValueString()

			rules = append(rules, restapi.RuleSpecification{
				DType:               restapi.EntityCountVerificationRuleType,
				Severity:            severity,
				ConditionOperator:   &conditionOperator,
				ConditionValue:      &conditionValue,
				MatchingEntityType:  &matchingEntityType,
				MatchingOperator:    &matchingOperator,
				MatchingEntityLabel: &matchingEntityLabel,
			})
		}

		// Process entity verification rule
		if model.Rules.EntityVerification != nil {
			severity := mapSeverityToInt(model.Rules.EntityVerification.Severity.ValueString())
			matchingEntityType := model.Rules.EntityVerification.MatchingEntityType.ValueString()
			matchingOperator := model.Rules.EntityVerification.MatchingOperator.ValueString()
			matchingEntityLabel := model.Rules.EntityVerification.MatchingEntityLabel.ValueString()
			offlineDuration := int(model.Rules.EntityVerification.OfflineDuration.ValueInt64())

			rules = append(rules, restapi.RuleSpecification{
				DType:               restapi.EntityVerificationRuleType,
				Severity:            severity,
				MatchingEntityType:  &matchingEntityType,
				MatchingOperator:    &matchingOperator,
				MatchingEntityLabel: &matchingEntityLabel,
				OfflineDuration:     &offlineDuration,
			})
		}

		// Process host availability rule
		if model.Rules.HostAvailability != nil {
			severity := mapSeverityToInt(model.Rules.HostAvailability.Severity.ValueString())
			offlineDuration := int(model.Rules.HostAvailability.OfflineDuration.ValueInt64())

			var closeAfter *int
			if !model.Rules.HostAvailability.CloseAfter.IsNull() {
				ca := int(model.Rules.HostAvailability.CloseAfter.ValueInt64())
				closeAfter = &ca
			}

			// Parse tag filter if provided
			var tagFilter *restapi.TagFilter
			if !model.Rules.HostAvailability.TagFilter.IsNull() && model.Rules.HostAvailability.TagFilter.ValueString() != "" {
				tagFilterStr := model.Rules.HostAvailability.TagFilter.ValueString()
				parser := tagfilter.NewParser()
				expr, err := parser.Parse(tagFilterStr)
				if err != nil {
					diags.AddError(
						CustomEventSpecificationResourceErrParseTagFilter,
						fmt.Sprintf(CustomEventSpecificationResourceErrParseTagFilterMsg, err),
					)
					return nil, diags
				}

				mapper := tagfilter.NewMapper()
				tagFilter = mapper.ToAPIModel(expr)
			}

			rules = append(rules, restapi.RuleSpecification{
				DType:           restapi.HostAvailabilityRuleType,
				Severity:        severity,
				OfflineDuration: &offlineDuration,
				CloseAfter:      closeAfter,
				TagFilter:       tagFilter,
			})
		}

		// Process system rule
		if model.Rules.System != nil {
			severity := mapSeverityToInt(model.Rules.System.Severity.ValueString())
			systemRuleID := model.Rules.System.SystemRuleID.ValueString()

			rules = append(rules, restapi.RuleSpecification{
				DType:        restapi.SystemRuleType,
				Severity:     severity,
				SystemRuleID: &systemRuleID,
			})
		}

		// Process threshold rule
		if model.Rules.Threshold != nil {
			severity := mapSeverityToInt(model.Rules.Threshold.Severity.ValueString())
			metricName := model.Rules.Threshold.MetricName.ValueString()
			rollup := int(model.Rules.Threshold.Rollup.ValueInt64())
			window := int(model.Rules.Threshold.Window.ValueInt64())
			aggregation := model.Rules.Threshold.Aggregation.ValueString()
			conditionOperator := model.Rules.Threshold.ConditionOperator.ValueString()
			conditionValue := model.Rules.Threshold.ConditionValue.ValueFloat64()

			var metricPattern *restapi.MetricPattern
			if model.Rules.Threshold.MetricPattern != nil {
				prefix := model.Rules.Threshold.MetricPattern.Prefix.ValueString()

				var postfix *string
				if !model.Rules.Threshold.MetricPattern.Postfix.IsNull() && model.Rules.Threshold.MetricPattern.Postfix.ValueString() != "" {
					p := model.Rules.Threshold.MetricPattern.Postfix.ValueString()
					postfix = &p
				}

				var placeholder *string
				if !model.Rules.Threshold.MetricPattern.Placeholder.IsNull() && model.Rules.Threshold.MetricPattern.Placeholder.ValueString() != "" {
					p := model.Rules.Threshold.MetricPattern.Placeholder.ValueString()
					placeholder = &p
				}

				operator := model.Rules.Threshold.MetricPattern.Operator.ValueString()

				metricPattern = &restapi.MetricPattern{
					Prefix:      prefix,
					Postfix:     postfix,
					Placeholder: placeholder,
					Operator:    operator,
				}
			}

			rules = append(rules, restapi.RuleSpecification{
				DType:             restapi.ThresholdRuleType,
				Severity:          severity,
				MetricName:        &metricName,
				Rollup:            &rollup,
				Window:            &window,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
				MetricPattern:     metricPattern,
			})
		}
	}

	// Create the API object
	return &restapi.CustomEventSpecification{
		ID:                  id,
		Name:                model.Name.ValueString(),
		EntityType:          model.EntityType.ValueString(),
		Query:               query,
		Triggering:          model.Triggering.ValueBool(),
		Description:         description,
		ExpirationTime:      expirationTime,
		Enabled:             model.Enabled.ValueBool(),
		RuleLogicalOperator: model.RuleLogicalOperator.ValueString(),
		Rules:               rules,
	}, diags
}

// mapSeverityToInt maps the severity string to an integer value
func mapSeverityToInt(severity string) int {
	switch severity {
	case "warning":
		return 5
	case "critical":
		return 10
	default:
		return 5 // Default to warning
	}
}
