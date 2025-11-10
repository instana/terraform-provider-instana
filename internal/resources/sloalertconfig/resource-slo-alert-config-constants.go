package sloalertconfig

// Resource description constants
const (
	// SloAlertConfigDescResource is the description for the SLO Alert config resource
	SloAlertConfigDescResource = "This resource manages SLO Alert configurations in Instana."
	// SloAlertConfigDescID is the description for the ID field
	SloAlertConfigDescID = "The ID of the SLO Alert configuration."
	// SloAlertConfigDescName is the description for the name field
	SloAlertConfigDescName = "The name of the SLO Alert config"
	// SloAlertConfigDescDescription is the description for the description field
	SloAlertConfigDescDescription = "The description of the SLO Alert config"
	// SloAlertConfigDescSeverity is the description for the severity field
	SloAlertConfigDescSeverity = "The severity of the alert when triggered"
	// SloAlertConfigDescTriggering is the description for the triggering field
	SloAlertConfigDescTriggering = "Optional flag to indicate whether also an Incident is triggered or not. The default is false"
	// SloAlertConfigDescEnabled is the description for the enabled field
	SloAlertConfigDescEnabled = "Optional flag to indicate whether this Alert is Enabled"
	// SloAlertConfigDescAlertType is the description for the alert_type field
	SloAlertConfigDescAlertType = "What do you want to be alerted on? (Type of Smart Alert: status, error_budget, burn_rate_v2)"
	// SloAlertConfigDescSloIds is the description for the slo_ids field
	SloAlertConfigDescSloIds = "The SLO IDs that are monitored"
	// SloAlertConfigDescAlertChannelIds is the description for the alert_channel_ids field
	SloAlertConfigDescAlertChannelIds = "The IDs of the Alert Channels"
	// SloAlertConfigDescThreshold is the description for the threshold block
	SloAlertConfigDescThreshold = "Indicates the type of violation of the defined threshold."
	// SloAlertConfigDescThresholdType is the description for the threshold type field
	SloAlertConfigDescThresholdType = "The type of threshold (should be staticThreshold)."
	// SloAlertConfigDescThresholdOperator is the description for the threshold operator field
	SloAlertConfigDescThresholdOperator = "The operator used to evaluate this rule."
	// SloAlertConfigDescThresholdValue is the description for the threshold value field
	SloAlertConfigDescThresholdValue = "The threshold value for the alert condition."
	// SloAlertConfigDescTimeThreshold is the description for the time_threshold block
	SloAlertConfigDescTimeThreshold = "Defines the time threshold for triggering and suppressing alerts."
	// SloAlertConfigDescTimeThresholdWarmUp is the description for the warm_up field
	SloAlertConfigDescTimeThresholdWarmUp = "The duration for which the condition must be violated for the alert to be triggered (in ms)."
	// SloAlertConfigDescTimeThresholdCoolDown is the description for the cool_down field
	SloAlertConfigDescTimeThresholdCoolDown = "The duration for which the condition must remain suppressed for the alert to end (in ms)."
	// SloAlertConfigDescBurnRateConfig is the description for the burn_rate_config block
	SloAlertConfigDescBurnRateConfig = "List of burn rate configs fields."
	// SloAlertConfigDescBurnRateAlertWindowType is the description for the alert_window_type field
	SloAlertConfigDescBurnRateAlertWindowType = "The alert window type for the burn rate config."
	// SloAlertConfigDescBurnRateDuration is the description for the duration field
	SloAlertConfigDescBurnRateDuration = "The duration for the burn rate config."
	// SloAlertConfigDescBurnRateDurationUnitType is the description for the duration_unit_type field
	SloAlertConfigDescBurnRateDurationUnitType = "The duration unit type for the burn rate config."
	// SloAlertConfigDescBurnRateThresholdOperator is the description for the threshold_operator field
	SloAlertConfigDescBurnRateThresholdOperator = "The threshold operator for the burn rate config."
	// SloAlertConfigDescBurnRateThresholdValue is the description for the threshold_value field
	SloAlertConfigDescBurnRateThresholdValue = "The threshold value for the burn rate config."
)

// Error message constants
const (
	// SloAlertConfigErrMappingAlertType is the error title for mapping alert type
	SloAlertConfigErrMappingAlertType = "Error mapping alert type"
	// SloAlertConfigErrInvalidAlertType is the error message for invalid alert type
	SloAlertConfigErrInvalidAlertType = "Invalid alert_type: %s"
	// SloAlertConfigErrParsingDuration is the error title for parsing duration
	SloAlertConfigErrParsingDuration = "Error parsing duration"
	// SloAlertConfigErrParsingDurationMsg is the error message for parsing duration
	SloAlertConfigErrParsingDurationMsg = "Failed to parse duration: %s"
	// SloAlertConfigErrParsingThresholdValue is the error title for parsing threshold value
	SloAlertConfigErrParsingThresholdValue = "Error parsing threshold value"
	// SloAlertConfigErrParsingThresholdValueMsg is the error message for parsing threshold value
	SloAlertConfigErrParsingThresholdValueMsg = "Failed to parse threshold value: %s"
)
