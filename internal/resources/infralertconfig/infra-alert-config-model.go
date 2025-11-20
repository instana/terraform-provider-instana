package infralertconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/shared"
)

// InfraAlertConfigModel represents the data model for infrastructure alert configuration
type InfraAlertConfigModel struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Description        types.String             `tfsdk:"description"`
	TagFilter          types.String             `tfsdk:"tag_filter"`
	GroupBy            types.List               `tfsdk:"group_by"`
	AlertChannels      *InfraAlertChannelsModel `tfsdk:"alert_channels"`
	Granularity        types.Int64              `tfsdk:"granularity"`
	TimeThreshold      *InfraTimeThresholdModel `tfsdk:"time_threshold"`
	CustomPayloadField types.List               `tfsdk:"custom_payload_field"`
	Rules              *InfraRulesModel         `tfsdk:"rules"`
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
	MetricName             types.String                 `tfsdk:"metric_name"`
	EntityType             types.String                 `tfsdk:"entity_type"`
	Aggregation            types.String                 `tfsdk:"aggregation"`
	CrossSeriesAggregation types.String                 `tfsdk:"cross_series_aggregation"`
	Regex                  types.Bool                   `tfsdk:"regex"`
	ThresholdOperator      types.String                 `tfsdk:"threshold_operator"`
	ThresholdRule          *shared.ThresholdPluginModel `tfsdk:"threshold"`
}
