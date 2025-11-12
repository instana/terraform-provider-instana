package cutomeventspec

import "github.com/hashicorp/terraform-plugin-framework/types"

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
	Rules               *RulesModel  `tfsdk:"rules"`
}

// RulesModel represents the rules container in the custom event specification
type RulesModel struct {
	EntityCount             *EntityCountRuleModel             `tfsdk:"entity_count"`
	EntityCountVerification *EntityCountVerificationRuleModel `tfsdk:"entity_count_verification"`
	EntityVerification      *EntityVerificationRuleModel      `tfsdk:"entity_verification"`
	HostAvailability        *HostAvailabilityRuleModel        `tfsdk:"host_availability"`
	System                  *SystemRuleModel                  `tfsdk:"system"`
	Threshold               *ThresholdRuleModel               `tfsdk:"threshold"`
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
	Severity          types.String        `tfsdk:"severity"`
	MetricName        types.String        `tfsdk:"metric_name"`
	MetricPattern     *MetricPatternModel `tfsdk:"metric_pattern"`
	Rollup            types.Int64         `tfsdk:"rollup"`
	Window            types.Int64         `tfsdk:"window"`
	Aggregation       types.String        `tfsdk:"aggregation"`
	ConditionOperator types.String        `tfsdk:"condition_operator"`
	ConditionValue    types.Float64       `tfsdk:"condition_value"`
}

// MetricPatternModel represents a metric pattern in a threshold rule
type MetricPatternModel struct {
	Prefix      types.String `tfsdk:"prefix"`
	Postfix     types.String `tfsdk:"postfix"`
	Placeholder types.String `tfsdk:"placeholder"`
	Operator    types.String `tfsdk:"operator"`
}
