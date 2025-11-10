package cutomeventspec

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
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

const (
	CustomEventSpecificationFieldName           = "name"
	CustomEventSpecificationFieldEntityType     = "entity_type"
	CustomEventSpecificationFieldQuery          = "query"
	CustomEventSpecificationFieldTriggering     = "triggering"
	CustomEventSpecificationFieldDescription    = "description"
	CustomEventSpecificationFieldExpirationTime = "expiration_time"
	CustomEventSpecificationFieldEnabled        = "enabled"

	CustomEventSpecificationFieldRuleLogicalOperator         = "rule_logical_operator"
	CustomEventSpecificationFieldRules                       = "rules"
	CustomEventSpecificationFieldEntityCountRule             = "entity_count"
	CustomEventSpecificationFieldEntityCountVerificationRule = "entity_count_verification"
	CustomEventSpecificationFieldEntityVerificationRule      = "entity_verification"
	CustomEventSpecificationFieldHostAvailabilityRule        = "host_availability"
	CustomEventSpecificationFieldSystemRule                  = "system"
	CustomEventSpecificationFieldThresholdRule               = "threshold"

	CustomEventSpecificationRuleFieldSeverity                          = "severity"
	CustomEventSpecificationRuleFieldMatchingEntityType                = "matching_entity_type"
	CustomEventSpecificationRuleFieldMatchingOperator                  = "matching_operator"
	CustomEventSpecificationRuleFieldMatchingEntityLabel               = "matching_entity_label"
	CustomEventSpecificationRuleFieldOfflineDuration                   = "offline_duration"
	CustomEventSpecificationSystemRuleFieldSystemRuleId                = "system_rule_id"
	CustomEventSpecificationThresholdRuleFieldMetricName               = "metric_name"
	CustomEventSpecificationThresholdRuleFieldRollup                   = "rollup"
	CustomEventSpecificationThresholdRuleFieldWindow                   = "window"
	CustomEventSpecificationThresholdRuleFieldAggregation              = "aggregation"
	CustomEventSpecificationRuleFieldConditionOperator                 = "condition_operator"
	CustomEventSpecificationRuleFieldConditionValue                    = "condition_value"
	CustomEventSpecificationThresholdRuleFieldMetricPattern            = "metric_pattern"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix      = "prefix"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix     = "postfix"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder = "placeholder"
	CustomEventSpecificationThresholdRuleFieldMetricPatternOperator    = "operator"
	CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter  = "close_after"
	CustomEventSpecificationHostAvailabilityRuleFieldTagFilter         = "tag_filter"
)

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
	Rules               types.Object `tfsdk:"rules"`
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
		},
		Blocks: map[string]schema.Block{
			"rules": schema.SingleNestedBlock{
				Description: CustomEventSpecificationResourceDescRules,
				Blocks: map[string]schema.Block{
					"entity_count": schema.ListNestedBlock{
						Description: CustomEventSpecificationResourceDescEntityCountRules,
						NestedObject: schema.NestedBlockObject{
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
					},
					"entity_count_verification": schema.ListNestedBlock{
						Description: CustomEventSpecificationResourceDescEntityCountVerification,
						NestedObject: schema.NestedBlockObject{
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
					},
					"entity_verification": schema.ListNestedBlock{
						Description: CustomEventSpecificationResourceDescEntityVerification,
						NestedObject: schema.NestedBlockObject{
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
					},
					"host_availability": schema.ListNestedBlock{
						Description: CustomEventSpecificationResourceDescHostAvailability,
						NestedObject: schema.NestedBlockObject{
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
					},
					"system": schema.ListNestedBlock{
						Description: CustomEventSpecificationResourceDescSystemRules,
						NestedObject: schema.NestedBlockObject{
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
					},
					"threshold": schema.ListNestedBlock{
						Description: CustomEventSpecificationResourceDescThresholdRules,
						NestedObject: schema.NestedBlockObject{
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
							},
							Blocks: map[string]schema.Block{
								"metric_pattern": schema.ListNestedBlock{
									Description: CustomEventSpecificationResourceDescMetricPattern,
									NestedObject: schema.NestedBlockObject{
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
						ConditionOperator: util.SetStringPointerToState(rule.ConditionOperator),
						ConditionValue:    util.SetFloat64PointerToState(rule.ConditionValue),
					})
				}
			case restapi.EntityCountVerificationRuleType:
				if rule.ConditionOperator != nil && rule.ConditionValue != nil &&
					rule.MatchingEntityType != nil && rule.MatchingOperator != nil && rule.MatchingEntityLabel != nil {
					entityCountVerificationRules = append(entityCountVerificationRules, EntityCountVerificationRuleModel{
						Severity:            mapIntToSeverityString(rule.Severity),
						ConditionOperator:   util.SetStringPointerToState(rule.ConditionOperator),
						ConditionValue:      util.SetFloat64PointerToState(rule.ConditionValue),
						MatchingEntityType:  util.SetStringPointerToState(rule.MatchingEntityType),
						MatchingOperator:    util.SetStringPointerToState(rule.MatchingOperator),
						MatchingEntityLabel: util.SetStringPointerToState(rule.MatchingEntityLabel),
					})
				}
			case restapi.EntityVerificationRuleType:
				if rule.MatchingEntityType != nil && rule.MatchingOperator != nil &&
					rule.MatchingEntityLabel != nil && rule.OfflineDuration != nil {
					entityVerificationRules = append(entityVerificationRules, EntityVerificationRuleModel{
						Severity:            mapIntToSeverityString(rule.Severity),
						MatchingEntityType:  util.SetStringPointerToState(rule.MatchingEntityType),
						MatchingOperator:    util.SetStringPointerToState(rule.MatchingOperator),
						MatchingEntityLabel: util.SetStringPointerToState(rule.MatchingEntityLabel),
						OfflineDuration:     util.SetInt64PointerToState(rule.OfflineDuration),
					})
				}
			case restapi.HostAvailabilityRuleType:
				if rule.OfflineDuration != nil {
					hostRule := HostAvailabilityRuleModel{
						Severity:        mapIntToSeverityString(rule.Severity),
						OfflineDuration: util.SetInt64PointerToState(rule.OfflineDuration),
						TagFilter:       types.StringValue(""), // Default empty string
					}

					if rule.CloseAfter != nil {
						hostRule.CloseAfter = util.SetInt64PointerToState(rule.CloseAfter)
					} else {
						hostRule.CloseAfter = types.Int64Null()
					}

					// Handle tag filter conversion
					if rule.TagFilter != nil {
						// Convert tag filter to string representation
						normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(rule.TagFilter)
						if err == nil && normalizedTagFilterString != nil {
							hostRule.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
						}
					}

					hostAvailabilityRules = append(hostAvailabilityRules, hostRule)
				}
			case restapi.SystemRuleType:
				if rule.SystemRuleID != nil {
					systemRules = append(systemRules, SystemRuleModel{
						Severity:     mapIntToSeverityString(rule.Severity),
						SystemRuleID: util.SetStringPointerToState(rule.SystemRuleID),
					})
				}
			case restapi.ThresholdRuleType:
				if rule.MetricName != nil && rule.Rollup != nil && rule.Window != nil &&
					rule.Aggregation != nil && rule.ConditionOperator != nil && rule.ConditionValue != nil {
					thresholdRule := ThresholdRuleModel{
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
						metricPatternModel := MetricPatternModel{
							Prefix:   types.StringValue(rule.MetricPattern.Prefix),
							Operator: types.StringValue(rule.MetricPattern.Operator),
						}

						metricPatternModel.Postfix = util.SetStringPointerToState(rule.MetricPattern.Postfix)

						metricPatternModel.Placeholder = util.SetStringPointerToState(rule.MetricPattern.Placeholder)

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
		model.Rules = rulesObj
	} else {
		// No rules
		model.Rules = types.ObjectNull(map[string]attr.Type{
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
