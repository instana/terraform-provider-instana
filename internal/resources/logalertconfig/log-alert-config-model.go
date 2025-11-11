package logalertconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// LogAlertConfigModel represents the data model for the log alert configuration resource
type LogAlertConfigModel struct {
	ID                  types.String        `tfsdk:"id"`
	Name                types.String        `tfsdk:"name"`
	Description         types.String        `tfsdk:"description"`
	AlertChannels       types.List          `tfsdk:"alert_channels"`
	GracePeriod         types.Int64         `tfsdk:"grace_period"`
	GroupBy             types.List          `tfsdk:"group_by"`
	Granularity         types.Int64         `tfsdk:"granularity"`
	TagFilter           types.String        `tfsdk:"tag_filter"`
	Rules               types.Object        `tfsdk:"rules"`
	TimeThreshold       *TimeThresholdModel `tfsdk:"time_threshold"`
	CustomPayloadFields types.List          `tfsdk:"custom_payload_field"`
}

type TimeThresholdModel struct {
	ViolationsInSequence *ViolationsInSequenceModel `tfsdk:"violations_in_sequence"`
}

type ViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}

// GroupByModel represents a group by tag in the Terraform model
type GroupByModel struct {
	TagName types.String `tfsdk:"tag_name"`
	Key     types.String `tfsdk:"key"`
}

// AlertChannelsModel represents alert channels in the Terraform model
type AlertChannelsModel struct {
	Warning  types.List `tfsdk:"warning"`
	Critical types.List `tfsdk:"critical"`
}

// LogAlertRuleModel represents a log alert rule in the Terraform model
type LogAlertRuleModel struct {
	MetricName        types.String `tfsdk:"metric_name"`
	AlertType         types.String `tfsdk:"alert_type"`
	Aggregation       types.String `tfsdk:"aggregation"`
	ThresholdOperator types.String `tfsdk:"threshold_operator"`
	Threshold         types.Object `tfsdk:"threshold"`
}

// ThresholdModel represents a threshold in the Terraform model
type ThresholdModel struct {
	Warning  types.Object `tfsdk:"warning"`
	Critical types.Object `tfsdk:"critical"`
}
