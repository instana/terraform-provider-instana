package tf_framework

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SloConfigModel represents the Terraform resource data model for the SLO Config resource
type SloConfigModel struct {
	ID         types.String  `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Target     types.Float64 `tfsdk:"target"`
	Tags       types.List    `tfsdk:"tags"`
	Entity     types.List    `tfsdk:"entity"`
	Indicator  types.List    `tfsdk:"indicator"`
	TimeWindow types.List    `tfsdk:"time_window"`
}

// SloApplicationEntityModel represents the application entity model for SLO Config
type SloApplicationEntityModel struct {
	ApplicationID    types.String `tfsdk:"application_id"`
	ServiceID        types.String `tfsdk:"service_id"`
	EndpointID       types.String `tfsdk:"endpoint_id"`
	BoundaryScope    types.String `tfsdk:"boundary_scope"`
	IncludeSynthetic types.Bool   `tfsdk:"include_synthetic"`
	IncludeInternal  types.Bool   `tfsdk:"include_internal"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// SloWebsiteEntityModel represents the website entity model for SLO Config
type SloWebsiteEntityModel struct {
	WebsiteID        types.String `tfsdk:"website_id"`
	BeaconType       types.String `tfsdk:"beacon_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// SloSyntheticEntityModel represents the synthetic entity model for SLO Config
type SloSyntheticEntityModel struct {
	SyntheticTestIDs types.List   `tfsdk:"synthetic_test_ids"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// SloEntityModel represents the entity model for SLO Config
type SloEntityModel struct {
	Application *SloApplicationEntityModel `tfsdk:"application"`
	Website     *SloWebsiteEntityModel     `tfsdk:"website"`
	Synthetic   *SloSyntheticEntityModel   `tfsdk:"synthetic"`
}

// SloTimeBasedLatencyIndicatorModel represents the time-based latency indicator model for SLO Config
type SloTimeBasedLatencyIndicatorModel struct {
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// SloEventBasedLatencyIndicatorModel represents the event-based latency indicator model for SLO Config
type SloEventBasedLatencyIndicatorModel struct {
	Threshold types.Float64 `tfsdk:"threshold"`
}

// SloTimeBasedAvailabilityIndicatorModel represents the time-based availability indicator model for SLO Config
type SloTimeBasedAvailabilityIndicatorModel struct {
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// SloEventBasedAvailabilityIndicatorModel represents the event-based availability indicator model for SLO Config
type SloEventBasedAvailabilityIndicatorModel struct {
	// No fields required for this indicator type
}

// SloTrafficIndicatorModel represents the traffic indicator model for SLO Config
type SloTrafficIndicatorModel struct {
	TrafficType types.String  `tfsdk:"traffic_type"`
	Threshold   types.Float64 `tfsdk:"threshold"`
	Aggregation types.String  `tfsdk:"aggregation"`
}

// SloCustomIndicatorModel represents the custom indicator model for SLO Config
type SloCustomIndicatorModel struct {
	GoodEventFilterExpression types.String `tfsdk:"good_event_filter_expression"`
	BadEventFilterExpression  types.String `tfsdk:"bad_event_filter_expression"`
}

// SloIndicatorModel represents the indicator model for SLO Config
type SloIndicatorModel struct {
	TimeBasedLatency       *SloTimeBasedLatencyIndicatorModel       `tfsdk:"time_based_latency"`
	EventBasedLatency      *SloEventBasedLatencyIndicatorModel      `tfsdk:"event_based_latency"`
	TimeBasedAvailability  *SloTimeBasedAvailabilityIndicatorModel  `tfsdk:"time_based_availability"`
	EventBasedAvailability *SloEventBasedAvailabilityIndicatorModel `tfsdk:"event_based_availability"`
	Traffic                *SloTrafficIndicatorModel                `tfsdk:"traffic"`
	Custom                 *SloCustomIndicatorModel                 `tfsdk:"custom"`
}

// SloRollingTimeWindowModel represents the rolling time window model for SLO Config
type SloRollingTimeWindowModel struct {
	Duration     types.Int64  `tfsdk:"duration"`
	DurationUnit types.String `tfsdk:"duration_unit"`
	Timezone     types.String `tfsdk:"timezone"`
}

// SloFixedTimeWindowModel represents the fixed time window model for SLO Config
type SloFixedTimeWindowModel struct {
	Duration     types.Int64   `tfsdk:"duration"`
	DurationUnit types.String  `tfsdk:"duration_unit"`
	Timezone     types.String  `tfsdk:"timezone"`
	StartTime    types.Float64 `tfsdk:"start_timestamp"`
}

// SloTimeWindowModel represents the time window model for SLO Config
type SloTimeWindowModel struct {
	Rolling *SloRollingTimeWindowModel `tfsdk:"rolling"`
	Fixed   *SloFixedTimeWindowModel   `tfsdk:"fixed"`
}

// Made with Bob
