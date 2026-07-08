package sloconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SloConfigModel represents the data model for the SLO configuration resource
type SloConfigModel struct {
	ID         types.String     `tfsdk:"id"`
	Name       types.String     `tfsdk:"name"`
	Target     types.Float64    `tfsdk:"target"`
	Tags       types.Set        `tfsdk:"tags"`
	RbacTags   types.List       `tfsdk:"rbac_tags"`
	Entity     *EntityModel     `tfsdk:"entity"`
	Indicator  *IndicatorModel  `tfsdk:"indicator"`
	TimeWindow *TimeWindowModel `tfsdk:"time_window"`
}

// RbacTagModel represents an RBAC tag in the Terraform model
type RbacTagModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
}
type TimeWindowModel struct {
	FixedTimeWindowModel   *FixedTimeWindowModel   `tfsdk:"fixed"`
	RollingTimeWindowModel *RollingTimeWindowModel `tfsdk:"rolling"`
}
type EntityModel struct {
	ApplicationEntityModel    *ApplicationEntityModel    `tfsdk:"application"`
	WebsiteEntityModel        *WebsiteEntityModel        `tfsdk:"website"`
	SyntheticEntityModel      *SyntheticEntityModel      `tfsdk:"synthetic"`
	InfrastructureEntityModel *InfrastructureEntityModel `tfsdk:"infrastructure"`
	MobileEntityModel         *MobileEntityModel         `tfsdk:"mobile"`
}
type IndicatorModel struct {
	TimeBasedLatencyIndicatorModel       *TimeBasedLatencyIndicatorModel       `tfsdk:"time_based_latency"`
	EventBasedLatencyIndicatorModel      *EventBasedLatencyIndicatorModel      `tfsdk:"event_based_latency"`
	TimeBasedAvailabilityIndicatorModel  *TimeBasedAvailabilityIndicatorModel  `tfsdk:"time_based_availability"`
	EventBasedAvailabilityIndicatorModel *EventBasedAvailabilityIndicatorModel `tfsdk:"event_based_availability"`
	TrafficIndicatorModel                *TrafficIndicatorModel                `tfsdk:"traffic"`
	CustomIndicatorModel                 *CustomIndicatorModel                 `tfsdk:"custom"`
	TimeBasedSaturationIndicatorModel    *TimeBasedSaturationIndicatorModel    `tfsdk:"time_based_saturation"`
	EventBasedSaturationIndicatorModel   *EventBasedSaturationIndicatorModel   `tfsdk:"event_based_saturation"`
	AdvancedCustomIndicatorModel         *AdvancedCustomIndicatorModel         `tfsdk:"advanced_custom"`
}

// ApplicationEntityModel represents an application entity in the Terraform model
type ApplicationEntityModel struct {
	ApplicationID    types.String `tfsdk:"application_id"`
	ServiceID        types.String `tfsdk:"service_id"`
	EndpointID       types.String `tfsdk:"endpoint_id"`
	BoundaryScope    types.String `tfsdk:"boundary_scope"`
	IncludeSynthetic types.Bool   `tfsdk:"include_synthetic"`
	IncludeInternal  types.Bool   `tfsdk:"include_internal"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// WebsiteEntityModel represents a website entity in the Terraform model
type WebsiteEntityModel struct {
	WebsiteID        types.String `tfsdk:"website_id"`
	BeaconType       types.String `tfsdk:"beacon_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// SyntheticEntityModel represents a synthetic entity in the Terraform model
type SyntheticEntityModel struct {
	SyntheticTestIDs              types.Set    `tfsdk:"synthetic_test_ids"`
	IncludeUnscheduledTestResults types.Bool   `tfsdk:"include_unscheduled_test_results"`
	FilterExpression              types.String `tfsdk:"filter_expression"`
}

// InfrastructureEntityModel represents an infrastructure entity in the Terraform model
type InfrastructureEntityModel struct {
	InfraType        types.String `tfsdk:"infra_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// MobileEntityModel represents a mobile app entity in the Terraform model
type MobileEntityModel struct {
	MobileIDs        types.Set    `tfsdk:"mobile_ids"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// EntityMetricScopeModel represents the scope nested inside a metric block
type EntityMetricScopeModel struct {
	ScopeType        types.String `tfsdk:"scope_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// EntityMetricModel represents the metric nested block for threshold-based mobile indicators
type EntityMetricModel struct {
	MetricName types.String            `tfsdk:"metric_name"`
	Scope      *EntityMetricScopeModel `tfsdk:"scope"`
}

// AdvancedFilterModel represents one side (good or bad) of an advanced-custom indicator
type AdvancedFilterModel struct {
	Aggregation types.String       `tfsdk:"aggregation"`
	Threshold   types.Float64      `tfsdk:"threshold"`
	Operator    types.String       `tfsdk:"operator"`
	Metric      *EntityMetricModel `tfsdk:"metric"`
}

// AdvancedCustomIndicatorModel represents the advanced-custom blueprint indicator
type AdvancedCustomIndicatorModel struct {
	Type       types.String         `tfsdk:"type"`
	GoodEvents *AdvancedFilterModel `tfsdk:"good_events"`
	BadEvents  *AdvancedFilterModel `tfsdk:"bad_events"`
}

// TimeBasedLatencyIndicatorModel represents a time-based latency indicator in the Terraform model
type TimeBasedLatencyIndicatorModel struct {
	Threshold   types.Float64      `tfsdk:"threshold"`
	Aggregation types.String       `tfsdk:"aggregation"`
	Metric      *EntityMetricModel `tfsdk:"metric"`
}

// EventBasedLatencyIndicatorModel represents an event-based latency indicator in the Terraform model
type EventBasedLatencyIndicatorModel struct {
	Threshold types.Float64      `tfsdk:"threshold"`
	Metric    *EntityMetricModel `tfsdk:"metric"`
}

// TimeBasedAvailabilityIndicatorModel represents a time-based availability indicator in the Terraform model
type TimeBasedAvailabilityIndicatorModel struct {
	Threshold   types.Float64      `tfsdk:"threshold"`
	Aggregation types.String       `tfsdk:"aggregation"`
	Metric      *EntityMetricModel `tfsdk:"metric"`
}

// EventBasedAvailabilityIndicatorModel represents an event-based availability indicator in the Terraform model
type EventBasedAvailabilityIndicatorModel struct {
	Threshold   types.Float64      `tfsdk:"threshold"`
	Aggregation types.String       `tfsdk:"aggregation"`
	Metric      *EntityMetricModel `tfsdk:"metric"`
}

// TrafficIndicatorModel represents a traffic indicator in the Terraform model
type TrafficIndicatorModel struct {
	TrafficType types.String       `tfsdk:"traffic_type"`
	Threshold   types.Float64      `tfsdk:"threshold"`
	Operator    types.String       `tfsdk:"operator"`
	Metric      *EntityMetricModel `tfsdk:"metric"`
}

// CustomIndicatorModel represents a custom indicator in the Terraform model
type CustomIndicatorModel struct {
	GoodEventFilterExpression types.String `tfsdk:"good_event_filter_expression"`
	BadEventFilterExpression  types.String `tfsdk:"bad_event_filter_expression"`
}

// TimeBasedSaturationIndicatorModel represents a saturation indicator in the Terraform model
type TimeBasedSaturationIndicatorModel struct {
	MetricName  types.String       `tfsdk:"metric_name"`
	Threshold   types.Float64      `tfsdk:"threshold"`
	Aggregation types.String       `tfsdk:"aggregation"`
	Operator    types.String       `tfsdk:"operator"`
	Metric      *EntityMetricModel `tfsdk:"metric"`
}

// EventBasedSaturationIndicatorModel represents a saturation indicator in the Terraform model
type EventBasedSaturationIndicatorModel struct {
	MetricName types.String       `tfsdk:"metric_name"`
	Threshold  types.Float64      `tfsdk:"threshold"`
	Operator   types.String       `tfsdk:"operator"`
	Metric     *EntityMetricModel `tfsdk:"metric"`
}

// RollingTimeWindowModel represents a rolling time window in the Terraform model
type RollingTimeWindowModel struct {
	Duration     types.Int64  `tfsdk:"duration"`
	DurationUnit types.String `tfsdk:"duration_unit"`
	Timezone     types.String `tfsdk:"timezone"`
}

// FixedTimeWindowModel represents a fixed time window in the Terraform model
type FixedTimeWindowModel struct {
	Duration       types.Int64   `tfsdk:"duration"`
	DurationUnit   types.String  `tfsdk:"duration_unit"`
	Timezone       types.String  `tfsdk:"timezone"`
	StartTimestamp types.Float64 `tfsdk:"start_timestamp"`
}
