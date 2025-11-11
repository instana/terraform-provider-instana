package applicationalertconfig

import (
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplicationAlertConfigModel represents the data model for the application alert configuration resource
type ApplicationAlertConfigModel struct {
	ID                  types.String `tfsdk:"id"`
	AlertChannelIDs     types.Set    `tfsdk:"alert_channel_ids"`
	AlertChannels       types.Map    `tfsdk:"alert_channels"`
	Applications        types.Set    `tfsdk:"application"`
	BoundaryScope       types.String `tfsdk:"boundary_scope"`
	CustomPayloadFields types.List   `tfsdk:"custom_payload_field"`
	Description         types.String `tfsdk:"description"`
	EvaluationType      types.String `tfsdk:"evaluation_type"`
	GracePeriod         types.Int64  `tfsdk:"grace_period"`
	Granularity         types.Int64  `tfsdk:"granularity"`
	IncludeInternal     types.Bool   `tfsdk:"include_internal"`
	IncludeSynthetic    types.Bool   `tfsdk:"include_synthetic"`
	Name                types.String `tfsdk:"name"`
	Rules               types.List   `tfsdk:"rules"`
	Severity            types.String `tfsdk:"severity"`
	TagFilter           types.String `tfsdk:"tag_filter"`
	TimeThreshold       types.Object `tfsdk:"time_threshold"`
	Triggering          types.Bool   `tfsdk:"triggering"`
}

// ApplicationModel represents an application in the application alert config
type ApplicationModel struct {
	ApplicationID types.String `tfsdk:"application_id"`
	Inclusive     types.Bool   `tfsdk:"inclusive"`
	Services      types.Set    `tfsdk:"service"`
}

// ApplicationThresholdModel represents a threshold in the application alert config
type ApplicationThresholdModel struct {
	Static           *shared.StaticTypeModel       `tfsdk:"static"`
	AdaptiveBaseline *shared.AdaptiveBaselineModel `tfsdk:"adaptive_baseline"`
}

// ServiceModel represents a service in the application alert config
type ServiceModel struct {
	ServiceID types.String `tfsdk:"service_id"`
	Inclusive types.Bool   `tfsdk:"inclusive"`
	Endpoints types.Set    `tfsdk:"endpoint"`
}

// EndpointModel represents an endpoint in the application alert config
type EndpointModel struct {
	EndpointID types.String `tfsdk:"endpoint_id"`
	Inclusive  types.Bool   `tfsdk:"inclusive"`
}

// RuleModel represents a rule in the application alert config
type RuleModel struct {
	ErrorRate  types.Object `tfsdk:"error_rate"`
	Errors     types.Object `tfsdk:"errors"`
	Logs       types.Object `tfsdk:"logs"`
	Slowness   types.Object `tfsdk:"slowness"`
	StatusCode types.Object `tfsdk:"status_code"`
	Throughput types.Object `tfsdk:"throughput"`
}

// RuleConfigModel represents the common configuration for rules
type RuleConfigModel struct {
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
}

// LogsRuleModel represents the logs rule configuration
type LogsRuleModel struct {
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
	Level       types.String `tfsdk:"level"`
	Message     types.String `tfsdk:"message"`
	Operator    types.String `tfsdk:"operator"`
}

// StatusCodeRuleModel represents the status code rule configuration
type StatusCodeRuleModel struct {
	MetricName      types.String `tfsdk:"metric_name"`
	Aggregation     types.String `tfsdk:"aggregation"`
	StatusCodeStart types.Int64  `tfsdk:"status_code_start"`
	StatusCodeEnd   types.Int64  `tfsdk:"status_code_end"`
}

type AppAlertTimeThresholdModel struct {
	RequestImpact        types.Object `tfsdk:"request_impact"`
	ViolationsInPeriod   types.Object `tfsdk:"violations_in_period"`
	ViolationsInSequence types.Object `tfsdk:"violations_in_sequence"`
}

// AppAlertRequestImpactModel represents the request impact time threshold configuration
type AppAlertRequestImpactModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
	Requests   types.Int64 `tfsdk:"requests"`
}

// AppAlertViolationsInPeriodModel represents the violations in period time threshold configuration
type AppAlertViolationsInPeriodModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
	Violations types.Int64 `tfsdk:"violations"`
}

// AppAlertViolationsInSequenceModel represents the violations in sequence time threshold configuration
type AppAlertViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}

// RuleWithThresholdModel represents a rule with multiple thresholds and severity levels
type RuleWithThresholdModel struct {
	Rule              types.Object `tfsdk:"rule"`
	ThresholdOperator types.String `tfsdk:"threshold_operator"`
	Thresholds        types.Object `tfsdk:"threshold"`
}

// ThresholdConfigRuleModel represents a threshold configuration for a rule
type ThresholdConfigRuleModel struct {
	Value types.Float64 `tfsdk:"value"`
}
