package syntheticalertconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// SyntheticAlertConfigModel represents the data model for the Synthetic Alert Config resource
type SyntheticAlertConfigModel struct {
	ID                  types.String                      `tfsdk:"id"`
	Name                types.String                      `tfsdk:"name"`
	Description         types.String                      `tfsdk:"description"`
	SyntheticTestIds    types.Set                         `tfsdk:"synthetic_test_ids"`
	Severity            types.Int64                       `tfsdk:"severity"`
	TagFilter           types.String                      `tfsdk:"tag_filter"`
	Rule                *SyntheticAlertRuleModel          `tfsdk:"rule"`
	AlertChannelIds     types.Set                         `tfsdk:"alert_channel_ids"`
	TimeThreshold       *SyntheticAlertTimeThresholdModel `tfsdk:"time_threshold"`
	GracePeriod         types.Int64                       `tfsdk:"grace_period"`
	CustomPayloadFields types.List                        `tfsdk:"custom_payload_field"`
}

// SyntheticAlertRuleModel represents the rule configuration for synthetic alerts
type SyntheticAlertRuleModel struct {
	AlertType   types.String `tfsdk:"alert_type"`
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
}

// SyntheticAlertTimeThresholdModel represents the time threshold configuration for synthetic alerts
type SyntheticAlertTimeThresholdModel struct {
	Type            types.String `tfsdk:"type"`
	ViolationsCount types.Int64  `tfsdk:"violations_count"`
}
