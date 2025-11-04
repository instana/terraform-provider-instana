package instana

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WebsiteAlertConfigModel represents the data model for the Website Alert Config resource
type WebsiteAlertConfigModel struct {
	ID                  types.String               `tfsdk:"id"`
	Name                types.String               `tfsdk:"name"`
	Description         types.String               `tfsdk:"description"`
	Severity            types.String               `tfsdk:"severity"`
	Triggering          types.Bool                 `tfsdk:"triggering"`
	WebsiteID           types.String               `tfsdk:"website_id"`
	TagFilter           types.String               `tfsdk:"tag_filter"`
	AlertChannelIDs     types.Set                  `tfsdk:"alert_channel_ids"`
	Granularity         types.Int64                `tfsdk:"granularity"`
	CustomPayloadFields types.List                 `tfsdk:"custom_payload_fields"`
	Rule                *WebsiteAlertRuleModel     `tfsdk:"rule"`
	Threshold           *WebsiteThresholdModel     `tfsdk:"threshold"`
	TimeThreshold       *WebsiteTimeThresholdModel `tfsdk:"time_threshold"`
	Rules               types.List                 `tfsdk:"rules"`
}

type RuleWithThresholdPluginModel struct {
	Rule              *WebsiteAlertRuleModel   `tfsdk:"rule"`
	ThresholdOperator types.String             `tfsdk:"operator"`
	Thresholds        *ThresholdAllPluginModel `tfsdk:"threshold"`
}

// // WebsiteAlertRuleModel represents the rule configuration for Website Alert Config
// type WebsiteAlertRuleModel struct {
// 	Slowness        types.List `tfsdk:"slowness"`
// 	SpecificJsError types.List `tfsdk:"specific_js_error"`
// 	StatusCode      types.List `tfsdk:"status_code"`
// 	Throughput      types.List `tfsdk:"throughput"`
// }

// WebsiteAlertRuleModel represents the rule configuration for Website Alert Config
type WebsiteAlertRuleModel struct {
	Slowness        *WebsiteAlertRuleConfigModel         `tfsdk:"slowness"`
	SpecificJsError *WebsiteAlertRuleConfigCompleteModel `tfsdk:"specific_js_error"`
	StatusCode      *WebsiteAlertRuleConfigCompleteModel `tfsdk:"status_code"`
	Throughput      *WebsiteAlertRuleConfigModel         `tfsdk:"throughput"`
}

// WebsiteAlertRuleConfigModel represents the common configuration for Website Alert Rules
type WebsiteAlertRuleConfigCompleteModel struct {
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
	Operator    types.String `tfsdk:"operator"`
	Value       types.String `tfsdk:"value"`
}

type WebsiteAlertRuleConfigModel struct {
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
}

// WebsiteTimeThresholdModel represents the time threshold configuration for Website Alert Config
type WebsiteTimeThresholdModel struct {
	UserImpactOfViolationsInSequence *WebsiteUserImpactOfViolationsInSequenceModel `tfsdk:"user_impact_of_violations_in_sequence"`
	ViolationsInPeriod               *WebsiteViolationsInPeriodModel               `tfsdk:"violations_in_period"`
	ViolationsInSequence             *WebsiteViolationsInSequenceModel             `tfsdk:"violations_in_sequence"`
}

// WebsiteUserImpactOfViolationsInSequenceModel represents the user impact configuration for time threshold
type WebsiteUserImpactOfViolationsInSequenceModel struct {
	TimeWindow              types.Int64   `tfsdk:"time_window"`
	ImpactMeasurementMethod types.String  `tfsdk:"impact_measurement_method"`
	UserPercentage          types.Float64 `tfsdk:"user_percentage"`
	Users                   types.Int64   `tfsdk:"users"`
}

// WebsiteViolationsInPeriodModel represents the violations in period configuration for time threshold
type WebsiteViolationsInPeriodModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
	Violations types.Int64 `tfsdk:"violations"`
}

// WebsiteViolationsInSequenceModel represents the violations in sequence configuration for time threshold
type WebsiteViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}

// CustomPayloadFieldModel represents a custom payload field
type CustomPayloadFieldModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// ThresholdModel represents the threshold configuration
type WebsiteThresholdModel struct {
	Static           *StaticTypeModel       `tfsdk:"static"`
	AdaptiveBaseline *AdaptiveBaselineModel `tfsdk:"adaptive_baseline"`
	HistoricBaseline *HistoricBaselineModel `tfsdk:"historic_baseline"`
}

// Made with Bob
