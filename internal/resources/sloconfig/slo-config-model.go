package sloconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SloConfigModel represents the data model for the SLO configuration resource
type SloConfigModel struct {
	ID         types.String     `tfsdk:"id"`
	Name       types.String     `tfsdk:"name"`
	Target     types.Float64    `tfsdk:"target"`
	Tags       []types.String   `tfsdk:"tags"`
	RbacTags   []RbacTagModel   `tfsdk:"rbac_tags"`
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
	SyntheticTestIDs types.Set    `tfsdk:"synthetic_test_ids"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// InfrastructureEntityModel represents an infrastructure entity in the Terraform model
type InfrastructureEntityModel struct {
	InfraType        types.String `tfsdk:"infra_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// TimeBasedLatencyIndicatorModel represents a time-based latency indicator in the Terraform model
type TimeBasedLatencyIndicatorModel struct {
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// EventBasedLatencyIndicatorModel represents an event-based latency indicator in the Terraform model
type EventBasedLatencyIndicatorModel struct {
	Threshold types.Float64 `tfsdk:"threshold"`
}

// TimeBasedAvailabilityIndicatorModel represents a time-based availability indicator in the Terraform model
type TimeBasedAvailabilityIndicatorModel struct {
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// EventBasedAvailabilityIndicatorModel represents an event-based availability indicator in the Terraform model
type EventBasedAvailabilityIndicatorModel struct {
	// No fields needed for this indicator type
}

// TrafficIndicatorModel represents a traffic indicator in the Terraform model
type TrafficIndicatorModel struct {
	TrafficType types.String  `tfsdk:"traffic_type"`
	Threshold   types.Float64 `tfsdk:"threshold"`
	Operator    types.String  `tfsdk:"operator"`
}

// CustomIndicatorModel represents a custom indicator in the Terraform model
type CustomIndicatorModel struct {
	GoodEventFilterExpression types.String `tfsdk:"good_event_filter_expression"`
	BadEventFilterExpression  types.String `tfsdk:"bad_event_filter_expression"`
}

// SaturationIndicatorModel represents a saturation indicator in the Terraform model
type TimeBasedSaturationIndicatorModel struct {
	MetricName  types.String  `tfsdk:"metric_name"`
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
	Operator    types.String  `tfsdk:"operator"`
}

// SaturationIndicatorModel represents a saturation indicator in the Terraform model
type EventBasedSaturationIndicatorModel struct {
	MetricName  types.String  `tfsdk:"metric_name"`
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
	Operator    types.String  `tfsdk:"operator"`
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
