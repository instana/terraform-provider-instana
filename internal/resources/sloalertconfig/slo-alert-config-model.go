package sloalertconfig

import (
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SloAlertConfigModel represents the data model for SLO Alert configuration
type SloAlertConfigModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Severity        types.Int64  `tfsdk:"severity"`
	Triggering      types.Bool   `tfsdk:"triggering"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	AlertType       types.String `tfsdk:"alert_type"`
	Threshold       types.Object `tfsdk:"threshold"`
	SloIds          types.Set    `tfsdk:"slo_ids"`
	AlertChannelIds types.Set    `tfsdk:"alert_channel_ids"`
	TimeThreshold   types.Object `tfsdk:"time_threshold"`
	BurnRateConfig  types.List   `tfsdk:"burn_rate_config"`
	CustomPayload   types.List   `tfsdk:"custom_payload_fields"`
}

// SloAlertThresholdModel represents the threshold configuration for SLO Alert
type SloAlertThresholdModel struct {
	Type     types.String  `tfsdk:"type"`
	Operator types.String  `tfsdk:"operator"`
	Value    types.Float64 `tfsdk:"value"`
}

// SloAlertTimeThresholdModel represents the time threshold configuration for SLO Alert
type SloAlertTimeThresholdModel struct {
	WarmUp   types.Int64 `tfsdk:"warm_up"`
	CoolDown types.Int64 `tfsdk:"cool_down"`
}

// SloAlertBurnRateConfigModel represents the burn rate configuration for SLO Alert
type SloAlertBurnRateConfigModel struct {
	AlertWindowType   types.String `tfsdk:"alert_window_type"`
	Duration          types.String `tfsdk:"duration"`
	DurationUnitType  types.String `tfsdk:"duration_unit_type"`
	ThresholdOperator types.String `tfsdk:"threshold_operator"`
	ThresholdValue    types.String `tfsdk:"threshold_value"`
}

// SloAlertCustomPayloadFieldModel represents a custom payload field
type SloAlertCustomPayloadFieldModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
	Type  types.String `tfsdk:"type"`
}

type sloAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}
