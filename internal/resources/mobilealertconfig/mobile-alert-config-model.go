package mobilealertconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/shared"
)

// MobileAlertConfigModel represents the data model for the mobile alert configuration resource
type MobileAlertConfigModel struct {
	ID                  types.String                   `tfsdk:"id"`
	Name                types.String                   `tfsdk:"name"`
	Description         types.String                   `tfsdk:"description"`
	MobileAppID         types.String                   `tfsdk:"mobile_app_id"`
	Severity            types.Int64                    `tfsdk:"severity"`
	Triggering          types.Bool                     `tfsdk:"triggering"`
	TagFilter           types.String                   `tfsdk:"tag_filter"`
	CompleteTagFilter   types.String                   `tfsdk:"complete_tag_filter"`
	AlertChannels       types.Map                      `tfsdk:"alert_channels"`
	Granularity         types.Int64                    `tfsdk:"granularity"`
	GracePeriod         types.Int64                    `tfsdk:"grace_period"`
	CustomPayloadFields types.List                     `tfsdk:"custom_payload_field"`
	Rules               []MobileRuleWithThresholdModel `tfsdk:"rules"`
	TimeThreshold       *MobileAlertTimeThresholdModel `tfsdk:"time_threshold"`
}

// MobileRuleWithThresholdModel represents a rule with multiple thresholds and severity levels
type MobileRuleWithThresholdModel struct {
	Rule              *MobileAlertRuleModel           `tfsdk:"rule"`
	ThresholdOperator types.String                    `tfsdk:"threshold_operator"`
	Thresholds        *shared.ThresholdAllPluginModel `tfsdk:"threshold"`
}

// MobileAlertRuleModel represents a mobile app alert rule
type MobileAlertRuleModel struct {
	AlertType   types.String `tfsdk:"alert_type"`
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
	Operator    types.String `tfsdk:"operator"`
	Value       types.String `tfsdk:"value"`
}

// MobileAlertTimeThresholdModel represents the time threshold configuration for mobile app alerts
type MobileAlertTimeThresholdModel struct {
	UserImpactOfViolationsInSequence *MobileUserImpactOfViolationsInSequenceModel `tfsdk:"user_impact_of_violations_in_sequence"`
	ViolationsInPeriod               *MobileViolationsInPeriodModel               `tfsdk:"violations_in_period"`
	ViolationsInSequence             *MobileViolationsInSequenceModel             `tfsdk:"violations_in_sequence"`
}

// MobileUserImpactOfViolationsInSequenceModel represents the user impact configuration for time threshold
type MobileUserImpactOfViolationsInSequenceModel struct {
	TimeWindow types.Int64   `tfsdk:"time_window"`
	Users      types.Int64   `tfsdk:"users"`
	Percentage types.Float64 `tfsdk:"percentage"`
}

// MobileViolationsInPeriodModel represents the violations in period configuration for time threshold
type MobileViolationsInPeriodModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
	Violations types.Int64 `tfsdk:"violations"`
}

// MobileViolationsInSequenceModel represents the violations in sequence configuration for time threshold
type MobileViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}
