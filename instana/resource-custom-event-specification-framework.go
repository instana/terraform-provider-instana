package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// ResourceInstanaCustomEventSpecificationFramework the name of the terraform-provider-instana resource to manage custom event specifications
const ResourceInstanaCustomEventSpecificationFramework = "custom_event_specification"

// CustomEventSpecificationModel represents the data model for the custom event specification resource
type CustomEventSpecificationModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	EntityType          types.String `tfsdk:"entity_type"`
	Query               types.String `tfsdk:"query"`
	Triggering          types.Bool   `tfsdk:"triggering"`
	Description         types.String `tfsdk:"description"`
	ExpirationTime      types.Int64  `tfsdk:"expiration_time"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	RuleLogicalOperator types.String `tfsdk:"rule_logical_operator"`
	Rules               types.List   `tfsdk:"rules"`
}

// RulesModel represents the rules container in the custom event specification
type RulesModel struct {
	EntityCount             types.List `tfsdk:"entity_count"`
	EntityCountVerification types.List `tfsdk:"entity_count_verification"`
	EntityVerification      types.List `tfsdk:"entity_verification"`
	HostAvailability        types.List `tfsdk:"host_availability"`
	System                  types.List `tfsdk:"system"`
	Threshold               types.List `tfsdk:"threshold"`
}

// EntityCountRuleModel represents an entity count rule
type EntityCountRuleModel struct {
	Severity          types.String  `tfsdk:"severity"`
	ConditionOperator types.String  `tfsdk:"condition_operator"`
	ConditionValue    types.Float64 `tfsdk:"condition_value"`
}

// EntityCountVerificationRuleModel represents an entity count verification rule
type EntityCountVerificationRuleModel struct {
	Severity            types.String  `tfsdk:"severity"`
	ConditionOperator   types.String  `tfsdk:"condition_operator"`
	ConditionValue      types.Float64 `tfsdk:"condition_value"`
	MatchingEntityType  types.String  `tfsdk:"matching_entity_type"`
	MatchingOperator    types.String  `tfsdk:"matching_operator"`
	MatchingEntityLabel types.String  `tfsdk:"matching_entity_label"`
}

// EntityVerificationRuleModel represents an entity verification rule
type EntityVerificationRuleModel struct {
	Severity            types.String `tfsdk:"severity"`
	MatchingEntityType  types.String `tfsdk:"matching_entity_type"`
	MatchingOperator    types.String `tfsdk:"matching_operator"`
	MatchingEntityLabel types.String `tfsdk:"matching_entity_label"`
	OfflineDuration     types.Int64  `tfsdk:"offline_duration"`
}

// HostAvailabilityRuleModel represents a host availability rule
type HostAvailabilityRuleModel struct {
	Severity        types.String `tfsdk:"severity"`
	OfflineDuration types.Int64  `tfsdk:"offline_duration"`
	CloseAfter      types.Int64  `tfsdk:"close_after"`
	TagFilter       types.String `tfsdk:"tag_filter"`
}

// SystemRuleModel represents a system rule
type SystemRuleModel struct {
	Severity     types.String `tfsdk:"severity"`
	SystemRuleID types.String `tfsdk:"system_rule_id"`
}

// ThresholdRuleModel represents a threshold rule
type ThresholdRuleModel struct {
	Severity          types.String  `tfsdk:"severity"`
	MetricName        types.String  `tfsdk:"metric_name"`
	MetricPattern     types.List    `tfsdk:"metric_pattern"`
	Rollup            types.Int64   `tfsdk:"rollup"`
	Window            types.Int64   `tfsdk:"window"`
	Aggregation       types.String  `tfsdk:"aggregation"`
	ConditionOperator types.String  `tfsdk:"condition_operator"`
	ConditionValue    types.Float64 `tfsdk:"condition_value"`
}

// MetricPatternModel represents a metric pattern in a threshold rule
type MetricPatternModel struct {
	Prefix      types.String `tfsdk:"prefix"`
	Postfix     types.String `tfsdk:"postfix"`
	Placeholder types.String `tfsdk:"placeholder"`
	Operator    types.String `tfsdk:"operator"`
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

// NewCustomEventSpecificationResourceHandleFramework creates the resource handle for Custom Event Specification
func NewCustomEventSpecificationResourceHandleFramework() ResourceHandleFramework[*restapi.CustomEventSpecification] {
	return &customEventSpecificationResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaCustomEventSpecificationFramework,
			Schema:        createCustomEventSpecificationSchema(),
			SchemaVersion: 1,
		},
	}
}

type customEventSpecificationResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *customEventSpecificationResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *customEventSpecificationResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

// createCustomEventSpecificationSchema creates the schema for the custom event specification resource
func createCustomEventSpecificationSchema() schema.Schema {
	return schema.Schema{
		Description: "This resource represents a custom event specification in Instana",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the custom event specification",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the custom event specification",
				Required:    true,
			},
			"entity_type": schema.StringAttribute{
				Description: "The entity type of the custom event specification",
				Required:    true,
			},
			"query": schema.StringAttribute{
				Description: "The query of the custom event specification",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"triggering": schema.BoolAttribute{
				Description: "Indicates if the custom event specification is triggering",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: "The description of the custom event specification",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"expiration_time": schema.Int64Attribute{
				Description: "The expiration time of the custom event specification",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates if the custom event specification is enabled",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"rule_logical_operator": schema.StringAttribute{
				Description: "The logical operator for the rules (AND, OR)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("AND"),
				Validators: []validator.String{
					stringvalidator.OneOf("AND", "OR"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rules": schema.SingleNestedBlock{
				Description: "The rules of the custom event specification",
				Blocks: map[string]schema.Block{
					"entity_count": schema.ListNestedBlock{
						Description: "Entity count rules",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Description: "The severity of the rule (warning, critical)",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("warning", "critical"),
									},
								},
								"condition_operator": schema.StringAttribute{
									Description: "The condition operator of the rule",
									Required:    true,
								},
								"condition_value": schema.Float64Attribute{
									Description: "The condition value of the rule",
									Required:    true,
								},
							},
						},
					},
					"entity_count_verification": schema.ListNestedBlock{
						Description: "Entity count verification rules",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Description: "The severity of the rule (warning, critical)",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("warning", "critical"),
									},
								},
								"condition_operator": schema.StringAttribute{
									Description: "The condition operator of the rule",
									Required:    true,
								},
								"condition_value": schema.Float64Attribute{
									Description: "The condition value of the rule",
									Required:    true,
								},
								"matching_entity_type": schema.StringAttribute{
									Description: "The matching entity type of the rule",
									Required:    true,
								},
								"matching_operator": schema.StringAttribute{
									Description: "The matching operator of the rule",
									Required:    true,
								},
								"matching_entity_label": schema.StringAttribute{
									Description: "The matching entity label of the rule",
									Required:    true,
								},
							},
						},
					},
					"entity_verification": schema.ListNestedBlock{
						Description: "Entity verification rules",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Description: "The severity of the rule (warning, critical)",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("warning", "critical"),
									},
								},
								"matching_entity_type": schema.StringAttribute{
									Description: "The matching entity type of the rule",
									Required:    true,
								},
								"matching_operator": schema.StringAttribute{
									Description: "The matching operator of the rule",
									Required:    true,
								},
								"matching_entity_label": schema.StringAttribute{
									Description: "The matching entity label of the rule",
									Required:    true,
								},
								"offline_duration": schema.Int64Attribute{
									Description: "The offline duration of the rule",
									Required:    true,
								},
							},
						},
					},
					"host_availability": schema.ListNestedBlock{
						Description: "Host availability rules",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Description: "The severity of the rule (warning, critical)",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("warning", "critical"),
									},
								},
								"offline_duration": schema.Int64Attribute{
									Description: "The offline duration of the rule",
									Required:    true,
								},
								"close_after": schema.Int64Attribute{
									Description: "The close after duration of the rule",
									Optional:    true,
								},
								"tag_filter": schema.StringAttribute{
									Description: "The tag filter of the rule",
									Required:    true,
								},
							},
						},
					},
					"system": schema.ListNestedBlock{
						Description: "System rules",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Description: "The severity of the rule (warning, critical)",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("warning", "critical"),
									},
								},
								"system_rule_id": schema.StringAttribute{
									Description: "The system rule ID",
									Required:    true,
								},
							},
						},
					},
					"threshold": schema.ListNestedBlock{
						Description: "Threshold rules",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"severity": schema.StringAttribute{
									Description: "The severity of the rule (warning, critical)",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("warning", "critical"),
									},
								},
								"metric_name": schema.StringAttribute{
									Description: "The metric name of the rule",
									Required:    true,
								},
								"rollup": schema.Int64Attribute{
									Description: "The rollup of the rule",
									Required:    true,
								},
								"window": schema.Int64Attribute{
									Description: "The window of the rule",
									Required:    true,
								},
								"aggregation": schema.StringAttribute{
									Description: "The aggregation of the rule",
									Required:    true,
								},
								"condition_operator": schema.StringAttribute{
									Description: "The condition operator of the rule",
									Required:    true,
								},
								"condition_value": schema.Float64Attribute{
									Description: "The condition value of the rule",
									Required:    true,
								},
							},
							Blocks: map[string]schema.Block{
								"metric_pattern": schema.ListNestedBlock{
									Description: "The metric pattern of the rule",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"prefix": schema.StringAttribute{
												Description: "The prefix of the metric pattern",
												Required:    true,
											},
											"postfix": schema.StringAttribute{
												Description: "The postfix of the metric pattern",
												Optional:    true,
												Computed:    true,
												Default:     stringdefault.StaticString(""),
											},
											"placeholder": schema.StringAttribute{
												Description: "The placeholder of the metric pattern",
												Optional:    true,
												Computed:    true,
												Default:     stringdefault.StaticString(""),
											},
											"operator": schema.StringAttribute{
												Description: "The operator of the metric pattern",
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
			},
		},
	}
}
func (r *customEventSpecificationResourceFramework) SetComputedFields(_ context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	// No computed fields to set
	return nil
}

func (r *customEventSpecificationResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, spec *restapi.CustomEventSpecification) diag.Diagnostics {
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
	if spec.Query != nil {
		model.Query = types.StringValue(*spec.Query)
	} else {
		model.Query = types.StringValue("")
	}

	if spec.Description != nil {
		model.Description = types.StringValue(*spec.Description)
	} else {
		model.Description = types.StringValue("")
	}

	if spec.ExpirationTime != nil {
		model.ExpirationTime = types.Int64Value(int64(*spec.ExpirationTime))
	} else {
		model.ExpirationTime = types.Int64Null()
	}

	// Process rules
	if len(spec.Rules) > 0 {
		// Create rule containers
		var entityCountRules []EntityCountRuleModel
		var entityCountVerificationRules []EntityCountVerificationRuleModel
		var entityVerificationRules []EntityVerificationRuleModel
		var hostAvailabilityRules []HostAvailabilityRuleModel
		var systemRules []SystemRuleModel
		var thresholdRules []ThresholdRuleModel

		// Process each rule based on its type
		for _, rule := range spec.Rules {
			switch rule.DType {
			case restapi.EntityCountRuleType:
				if rule.ConditionOperator != nil && rule.ConditionValue != nil {
					entityCountRules = append(entityCountRules, EntityCountRuleModel{
						Severity:          mapIntToSeverityString(rule.Severity),
						ConditionOperator: types.StringValue(*rule.ConditionOperator),
						ConditionValue:    types.Float64Value(*rule.ConditionValue),
					})
				}
			case restapi.EntityCountVerificationRuleType:
				if rule.ConditionOperator != nil && rule.ConditionValue != nil &&
					rule.MatchingEntityType != nil && rule.MatchingOperator != nil && rule.MatchingEntityLabel != nil {
					entityCountVerificationRules = append(entityCountVerificationRules, EntityCountVerificationRuleModel{
						Severity:            mapIntToSeverityString(rule.Severity),
						ConditionOperator:   types.StringValue(*rule.ConditionOperator),
						ConditionValue:      types.Float64Value(*rule.ConditionValue),
						MatchingEntityType:  types.StringValue(*rule.MatchingEntityType),
						MatchingOperator:    types.StringValue(*rule.MatchingOperator),
						MatchingEntityLabel: types.StringValue(*rule.MatchingEntityLabel),
					})
				}
			case restapi.EntityVerificationRuleType:
				if rule.MatchingEntityType != nil && rule.MatchingOperator != nil &&
					rule.MatchingEntityLabel != nil && rule.OfflineDuration != nil {
					entityVerificationRules = append(entityVerificationRules, EntityVerificationRuleModel{
						Severity:            mapIntToSeverityString(rule.Severity),
						MatchingEntityType:  types.StringValue(*rule.MatchingEntityType),
						MatchingOperator:    types.StringValue(*rule.MatchingOperator),
						MatchingEntityLabel: types.StringValue(*rule.MatchingEntityLabel),
						OfflineDuration:     types.Int64Value(int64(*rule.OfflineDuration)),
					})
				}
			case restapi.HostAvailabilityRuleType:
				if rule.OfflineDuration != nil {
					hostRule := HostAvailabilityRuleModel{
						Severity:        mapIntToSeverityString(rule.Severity),
						OfflineDuration: types.Int64Value(int64(*rule.OfflineDuration)),
						TagFilter:       types.StringValue(""), // Default empty string
					}

					if rule.CloseAfter != nil {
						hostRule.CloseAfter = types.Int64Value(int64(*rule.CloseAfter))
					} else {
						hostRule.CloseAfter = types.Int64Null()
					}

					// Handle tag filter conversion
					if rule.TagFilter != nil {
						// Convert tag filter to string representation
						normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(rule.TagFilter)
						if err == nil && normalizedTagFilterString != nil {
							hostRule.TagFilter = types.StringValue(*normalizedTagFilterString)
						}
					}

					hostAvailabilityRules = append(hostAvailabilityRules, hostRule)
				}
			case restapi.SystemRuleType:
				if rule.SystemRuleID != nil {
					systemRules = append(systemRules, SystemRuleModel{
						Severity:     mapIntToSeverityString(rule.Severity),
						SystemRuleID: types.StringValue(*rule.SystemRuleID),
					})
				}
			case restapi.ThresholdRuleType:
				if rule.MetricName != nil && rule.Rollup != nil && rule.Window != nil &&
					rule.Aggregation != nil && rule.ConditionOperator != nil && rule.ConditionValue != nil {
					thresholdRule := ThresholdRuleModel{
						Severity:          mapIntToSeverityString(rule.Severity),
						MetricName:        types.StringValue(*rule.MetricName),
						Rollup:            types.Int64Value(int64(*rule.Rollup)),
						Window:            types.Int64Value(int64(*rule.Window)),
						Aggregation:       types.StringValue(*rule.Aggregation),
						ConditionOperator: types.StringValue(*rule.ConditionOperator),
						ConditionValue:    types.Float64Value(*rule.ConditionValue),
					}

					// Handle metric pattern if present
					if rule.MetricPattern != nil {
						metricPatternModel := MetricPatternModel{
							Prefix:   types.StringValue(rule.MetricPattern.Prefix),
							Operator: types.StringValue(rule.MetricPattern.Operator),
						}

						if rule.MetricPattern.Postfix != nil {
							metricPatternModel.Postfix = types.StringValue(*rule.MetricPattern.Postfix)
						} else {
							metricPatternModel.Postfix = types.StringValue("")
						}

						if rule.MetricPattern.Placeholder != nil {
							metricPatternModel.Placeholder = types.StringValue(*rule.MetricPattern.Placeholder)
						} else {
							metricPatternModel.Placeholder = types.StringValue("")
						}

						// Create a list of metric patterns with this single pattern
						metricPatterns, diags := types.ListValueFrom(ctx, types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"prefix":      types.StringType,
								"postfix":     types.StringType,
								"placeholder": types.StringType,
								"operator":    types.StringType,
							},
						}, []MetricPatternModel{metricPatternModel})

						if diags.HasError() {
							return diags
						}

						thresholdRule.MetricPattern = metricPatterns
					} else {
						// Empty list for metric pattern
						emptyList, diags := types.ListValueFrom(ctx, types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"prefix":      types.StringType,
								"postfix":     types.StringType,
								"placeholder": types.StringType,
								"operator":    types.StringType,
							},
						}, []MetricPatternModel{})

						if diags.HasError() {
							return diags
						}

						thresholdRule.MetricPattern = emptyList
					}

					thresholdRules = append(thresholdRules, thresholdRule)
				}
			}
		}

		// Create lists for each rule type
		entityCountList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"severity":           types.StringType,
				"condition_operator": types.StringType,
				"condition_value":    types.Float64Type,
			},
		}, entityCountRules)
		if diags.HasError() {
			return diags
		}

		entityCountVerificationList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"severity":              types.StringType,
				"condition_operator":    types.StringType,
				"condition_value":       types.Float64Type,
				"matching_entity_type":  types.StringType,
				"matching_operator":     types.StringType,
				"matching_entity_label": types.StringType,
			},
		}, entityCountVerificationRules)
		if diags.HasError() {
			return diags
		}

		entityVerificationList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"severity":              types.StringType,
				"matching_entity_type":  types.StringType,
				"matching_operator":     types.StringType,
				"matching_entity_label": types.StringType,
				"offline_duration":      types.Int64Type,
			},
		}, entityVerificationRules)
		if diags.HasError() {
			return diags
		}

		hostAvailabilityList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"severity":         types.StringType,
				"offline_duration": types.Int64Type,
				"close_after":      types.Int64Type,
				"tag_filter":       types.StringType,
			},
		}, hostAvailabilityRules)
		if diags.HasError() {
			return diags
		}

		systemList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"severity":       types.StringType,
				"system_rule_id": types.StringType,
			},
		}, systemRules)
		if diags.HasError() {
			return diags
		}

		thresholdList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"severity":    types.StringType,
				"metric_name": types.StringType,
				"metric_pattern": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"prefix":      types.StringType,
						"postfix":     types.StringType,
						"placeholder": types.StringType,
						"operator":    types.StringType,
					},
				}},
				"rollup":             types.Int64Type,
				"window":             types.Int64Type,
				"aggregation":        types.StringType,
				"condition_operator": types.StringType,
				"condition_value":    types.Float64Type,
			},
		}, thresholdRules)
		if diags.HasError() {
			return diags
		}

		// Create the rules model
		rulesModel := RulesModel{
			EntityCount:             entityCountList,
			EntityCountVerification: entityCountVerificationList,
			EntityVerification:      entityVerificationList,
			HostAvailability:        hostAvailabilityList,
			System:                  systemList,
			Threshold:               thresholdList,
		}

		// Convert the rules model to a list
		rulesObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"entity_count": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity":           types.StringType,
					"condition_operator": types.StringType,
					"condition_value":    types.Float64Type,
				},
			}},
			"entity_count_verification": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity":              types.StringType,
					"condition_operator":    types.StringType,
					"condition_value":       types.Float64Type,
					"matching_entity_type":  types.StringType,
					"matching_operator":     types.StringType,
					"matching_entity_label": types.StringType,
				},
			}},
			"entity_verification": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity":              types.StringType,
					"matching_entity_type":  types.StringType,
					"matching_operator":     types.StringType,
					"matching_entity_label": types.StringType,
					"offline_duration":      types.Int64Type,
				},
			}},
			"host_availability": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity":         types.StringType,
					"offline_duration": types.Int64Type,
					"close_after":      types.Int64Type,
					"tag_filter":       types.StringType,
				},
			}},
			"system": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity":       types.StringType,
					"system_rule_id": types.StringType,
				},
			}},
			"threshold": types.ListType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"severity":    types.StringType,
					"metric_name": types.StringType,
					"metric_pattern": types.ListType{ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"prefix":      types.StringType,
							"postfix":     types.StringType,
							"placeholder": types.StringType,
							"operator":    types.StringType,
						},
					}},
					"rollup":             types.Int64Type,
					"window":             types.Int64Type,
					"aggregation":        types.StringType,
					"condition_operator": types.StringType,
					"condition_value":    types.Float64Type,
				},
			}},
		}, rulesModel)
		if diags.HasError() {
			return diags
		}

		// Set the rules in the model
		model.Rules = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"entity_count": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity":           types.StringType,
						"condition_operator": types.StringType,
						"condition_value":    types.Float64Type,
					},
				}},
				"entity_count_verification": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity":              types.StringType,
						"condition_operator":    types.StringType,
						"condition_value":       types.Float64Type,
						"matching_entity_type":  types.StringType,
						"matching_operator":     types.StringType,
						"matching_entity_label": types.StringType,
					},
				}},
				"entity_verification": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity":              types.StringType,
						"matching_entity_type":  types.StringType,
						"matching_operator":     types.StringType,
						"matching_entity_label": types.StringType,
						"offline_duration":      types.Int64Type,
					},
				}},
				"host_availability": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity":         types.StringType,
						"offline_duration": types.Int64Type,
						"close_after":      types.Int64Type,
						"tag_filter":       types.StringType,
					},
				}},
				"system": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity":       types.StringType,
						"system_rule_id": types.StringType,
					},
				}},
				"threshold": types.ListType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"severity":    types.StringType,
						"metric_name": types.StringType,
						"metric_pattern": types.ListType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"prefix":      types.StringType,
								"postfix":     types.StringType,
								"placeholder": types.StringType,
								"operator":    types.StringType,
							},
						}},
						"rollup":             types.Int64Type,
						"window":             types.Int64Type,
						"aggregation":        types.StringType,
						"condition_operator": types.StringType,
						"condition_value":    types.Float64Type,
					},
				}},
			},
		}, []attr.Value{rulesObj})
	} else {
		// No rules
		model.Rules = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"entity_count":              types.ListType{ElemType: types.ObjectType{}},
				"entity_count_verification": types.ListType{ElemType: types.ObjectType{}},
				"entity_verification":       types.ListType{ElemType: types.ObjectType{}},
				"host_availability":         types.ListType{ElemType: types.ObjectType{}},
				"system":                    types.ListType{ElemType: types.ObjectType{}},
				"threshold":                 types.ListType{ElemType: types.ObjectType{}},
			},
		})
	}

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
	if !model.Rules.IsNull() && !model.Rules.IsUnknown() {
		var rulesModel RulesModel
		diags.Append(tfsdk.ValueAs(ctx, model.Rules, &rulesModel)...)
		if diags.HasError() {
			return nil, diags
		}

		// Process entity count rules
		if !rulesModel.EntityCount.IsNull() && !rulesModel.EntityCount.IsUnknown() {
			var entityCountRules []EntityCountRuleModel
			diags.Append(rulesModel.EntityCount.ElementsAs(ctx, &entityCountRules, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, rule := range entityCountRules {
				severity := mapSeverityToInt(rule.Severity.ValueString())
				conditionOperator := rule.ConditionOperator.ValueString()
				conditionValue := rule.ConditionValue.ValueFloat64()

				rules = append(rules, restapi.RuleSpecification{
					DType:             restapi.EntityCountRuleType,
					Severity:          severity,
					ConditionOperator: &conditionOperator,
					ConditionValue:    &conditionValue,
				})
			}
		}

		// Process entity count verification rules
		if !rulesModel.EntityCountVerification.IsNull() && !rulesModel.EntityCountVerification.IsUnknown() {
			var entityCountVerificationRules []EntityCountVerificationRuleModel
			diags.Append(rulesModel.EntityCountVerification.ElementsAs(ctx, &entityCountVerificationRules, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, rule := range entityCountVerificationRules {
				severity := mapSeverityToInt(rule.Severity.ValueString())
				conditionOperator := rule.ConditionOperator.ValueString()
				conditionValue := rule.ConditionValue.ValueFloat64()
				matchingEntityType := rule.MatchingEntityType.ValueString()
				matchingOperator := rule.MatchingOperator.ValueString()
				matchingEntityLabel := rule.MatchingEntityLabel.ValueString()

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
		}

		// Process entity verification rules
		if !rulesModel.EntityVerification.IsNull() && !rulesModel.EntityVerification.IsUnknown() {
			var entityVerificationRules []EntityVerificationRuleModel
			diags.Append(rulesModel.EntityVerification.ElementsAs(ctx, &entityVerificationRules, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, rule := range entityVerificationRules {
				severity := mapSeverityToInt(rule.Severity.ValueString())
				matchingEntityType := rule.MatchingEntityType.ValueString()
				matchingOperator := rule.MatchingOperator.ValueString()
				matchingEntityLabel := rule.MatchingEntityLabel.ValueString()
				offlineDuration := int(rule.OfflineDuration.ValueInt64())

				rules = append(rules, restapi.RuleSpecification{
					DType:               restapi.EntityVerificationRuleType,
					Severity:            severity,
					MatchingEntityType:  &matchingEntityType,
					MatchingOperator:    &matchingOperator,
					MatchingEntityLabel: &matchingEntityLabel,
					OfflineDuration:     &offlineDuration,
				})
			}
		}

		// Process host availability rules
		if !rulesModel.HostAvailability.IsNull() && !rulesModel.HostAvailability.IsUnknown() {
			var hostAvailabilityRules []HostAvailabilityRuleModel
			diags.Append(rulesModel.HostAvailability.ElementsAs(ctx, &hostAvailabilityRules, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, rule := range hostAvailabilityRules {
				severity := mapSeverityToInt(rule.Severity.ValueString())
				offlineDuration := int(rule.OfflineDuration.ValueInt64())

				var closeAfter *int
				if !rule.CloseAfter.IsNull() {
					ca := int(rule.CloseAfter.ValueInt64())
					closeAfter = &ca
				}

				// Parse tag filter if provided
				var tagFilter *restapi.TagFilter
				if !rule.TagFilter.IsNull() && rule.TagFilter.ValueString() != "" {
					tagFilterStr := rule.TagFilter.ValueString()
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

				rules = append(rules, restapi.RuleSpecification{
					DType:           restapi.HostAvailabilityRuleType,
					Severity:        severity,
					OfflineDuration: &offlineDuration,
					CloseAfter:      closeAfter,
					TagFilter:       tagFilter,
				})
			}
		}

		// Process system rules
		if !rulesModel.System.IsNull() && !rulesModel.System.IsUnknown() {
			var systemRules []SystemRuleModel
			diags.Append(rulesModel.System.ElementsAs(ctx, &systemRules, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, rule := range systemRules {
				severity := mapSeverityToInt(rule.Severity.ValueString())
				systemRuleID := rule.SystemRuleID.ValueString()

				rules = append(rules, restapi.RuleSpecification{
					DType:        restapi.SystemRuleType,
					Severity:     severity,
					SystemRuleID: &systemRuleID,
				})
			}
		}

		// Process threshold rules
		if !rulesModel.Threshold.IsNull() && !rulesModel.Threshold.IsUnknown() {
			var thresholdRules []ThresholdRuleModel
			diags.Append(rulesModel.Threshold.ElementsAs(ctx, &thresholdRules, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, rule := range thresholdRules {
				severity := mapSeverityToInt(rule.Severity.ValueString())
				metricName := rule.MetricName.ValueString()
				rollup := int(rule.Rollup.ValueInt64())
				window := int(rule.Window.ValueInt64())
				aggregation := rule.Aggregation.ValueString()
				conditionOperator := rule.ConditionOperator.ValueString()
				conditionValue := rule.ConditionValue.ValueFloat64()

				var metricPattern *restapi.MetricPattern
				if !rule.MetricPattern.IsNull() && !rule.MetricPattern.IsUnknown() {
					var metricPatterns []MetricPatternModel
					diags.Append(rule.MetricPattern.ElementsAs(ctx, &metricPatterns, false)...)
					if diags.HasError() {
						return nil, diags
					}

					if len(metricPatterns) > 0 {
						mp := metricPatterns[0]
						prefix := mp.Prefix.ValueString()

						var postfix *string
						if !mp.Postfix.IsNull() && mp.Postfix.ValueString() != "" {
							p := mp.Postfix.ValueString()
							postfix = &p
						}

						var placeholder *string
						if !mp.Placeholder.IsNull() && mp.Placeholder.ValueString() != "" {
							p := mp.Placeholder.ValueString()
							placeholder = &p
						}

						operator := mp.Operator.ValueString()

						metricPattern = &restapi.MetricPattern{
							Prefix:      prefix,
							Postfix:     postfix,
							Placeholder: placeholder,
							Operator:    operator,
						}
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
