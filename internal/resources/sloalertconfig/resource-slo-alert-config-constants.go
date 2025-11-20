package sloalertconfig

// ResourceInstanaSloAlertConfigFramework the name of the terraform-provider-instana resource to manage SLO Alert configurations
const ResourceInstanaSloAlertConfigFramework = "slo_alert_config"

const (
	//Slo Alert Config Field names for Terraform
	SloAlertConfigFieldName        = "name"
	SloAlertConfigFieldDescription = "description"
	SloAlertConfigFieldSeverity                        = "severity"
	SloAlertConfigFieldTriggering                      = "triggering"
	SloAlertConfigFieldAlertType                       = "alert_type"
	SloAlertConfigFieldThreshold                       = "threshold"
	SloAlertConfigFieldThresholdType                   = "type"
	SloAlertConfigFieldThresholdOperator               = "operator"
	SloAlertConfigFieldThresholdValue                  = "value"
	SloAlertConfigFieldSloIds                          = "slo_ids"
	SloAlertConfigFieldAlertChannelIds                 = "alert_channel_ids"
	SloAlertConfigFieldTimeThreshold                   = "time_threshold"
	SloAlertConfigFieldTimeThresholdWarmUp             = "warm_up"
	SloAlertConfigFieldTimeThresholdCoolDown           = "cool_down"
	SloAlertConfigFieldBurnRateConfig                  = "burn_rate_config"
	SloAlertConfigFieldBurnRateConfigDuration          = "duration"
	SloAlertConfigFieldBurnRateConfigThresholdValue    = "threshold_value"
	SloAlertConfigFieldBurnRateConfigThresholdOperator = "threshold_operator"
	SloAlertConfigFieldBurnRateConfigDurationUnitType  = "duration_unit_type"
	SloAlertConfigFieldBurnRateConfigAlertWindowType   = "alert_window_type"

	// Slo Alert Types for Terraform
	SloAlertConfigStatus      = "status"
	SloAlertConfigErrorBudget = "error_budget"
	SloAlertConfigBurnRateV2  = "burn_rate_v2"

	// Resource description constants

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

	// Error message constants

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
	// API Alert Type constants

	// APIAlertTypeServiceLevelsObjective represents the SERVICE_LEVELS_OBJECTIVE API alert type
	APIAlertTypeServiceLevelsObjective = "SERVICE_LEVELS_OBJECTIVE"
	// APIAlertTypeErrorBudget represents the ERROR_BUDGET API alert type
	APIAlertTypeErrorBudget = "ERROR_BUDGET"

	// API Metric constants

	// APIMetricStatus represents the STATUS API metric
	APIMetricStatus = "STATUS"
	// APIMetricBurnedPercentage represents the BURNED_PERCENTAGE API metric
	APIMetricBurnedPercentage = "BURNED_PERCENTAGE"
	// APIMetricBurnRate represents the legacy BURN_RATE API metric
	APIMetricBurnRate = "BURN_RATE"
	// APIMetricBurnRateV2 represents the BURN_RATE_V2 API metric
	APIMetricBurnRateV2 = "BURN_RATE_V2"

	// Threshold operator constants

	// OperatorGreaterThan represents the > operator
	OperatorGreaterThan = ">"
	// OperatorGreaterThanOrEqual represents the >= operator
	OperatorGreaterThanOrEqual = ">="
	// OperatorEqual represents the = operator
	OperatorEqual = "="
	// OperatorLessThanOrEqual represents the <= operator
	OperatorLessThanOrEqual = "<="
	// OperatorLessThan represents the < operator
	OperatorLessThan = "<"

	// Threshold type constants

	// ThresholdTypeStatic represents the static threshold type from API
	ThresholdTypeStatic = "static"
	// ThresholdTypeStaticThreshold represents the staticThreshold type for Terraform
	ThresholdTypeStaticThreshold = "staticThreshold"

	// Schema field identifier constants

	// SchemaFieldID represents the id field identifier
	SchemaFieldID = "id"
	// SchemaFieldName represents the name field identifier
	SchemaFieldName = "name"
	// SchemaFieldDescription represents the description field identifier
	SchemaFieldDescription = "description"
	// SchemaFieldSeverity represents the severity field identifier
	SchemaFieldSeverity = "severity"
	// SchemaFieldTriggering represents the triggering field identifier
	SchemaFieldTriggering = "triggering"
	// SchemaFieldAlertType represents the alert_type field identifier
	SchemaFieldAlertType = "alert_type"
	// SchemaFieldSloIds represents the slo_ids field identifier
	SchemaFieldSloIds = "slo_ids"
	// SchemaFieldAlertChannelIds represents the alert_channel_ids field identifier
	SchemaFieldAlertChannelIds = "alert_channel_ids"
	// SchemaFieldCustomPayloadFields represents the custom_payload_fields field identifier
	SchemaFieldCustomPayloadFields = "custom_payload_fields"
	// SchemaFieldThreshold represents the threshold field identifier
	SchemaFieldThreshold = "threshold"
	// SchemaFieldThresholdType represents the type field identifier within threshold
	SchemaFieldThresholdType = "type"
	// SchemaFieldThresholdOperator represents the operator field identifier within threshold
	SchemaFieldThresholdOperator = "operator"
	// SchemaFieldThresholdValue represents the value field identifier within threshold
	SchemaFieldThresholdValue = "value"
	// SchemaFieldTimeThreshold represents the time_threshold field identifier
	SchemaFieldTimeThreshold = "time_threshold"
	// SchemaFieldTimeThresholdWarmUp represents the warm_up field identifier
	SchemaFieldTimeThresholdWarmUp = "warm_up"
	// SchemaFieldTimeThresholdCoolDown represents the cool_down field identifier
	SchemaFieldTimeThresholdCoolDown = "cool_down"
	// SchemaFieldBurnRateConfig represents the burn_rate_config field identifier
	SchemaFieldBurnRateConfig = "burn_rate_config"
	// SchemaFieldBurnRateAlertWindowType represents the alert_window_type field identifier
	SchemaFieldBurnRateAlertWindowType = "alert_window_type"
	// SchemaFieldBurnRateDuration represents the duration field identifier
	SchemaFieldBurnRateDuration = "duration"
	// SchemaFieldBurnRateDurationUnitType represents the duration_unit_type field identifier
	SchemaFieldBurnRateDurationUnitType = "duration_unit_type"
	// SchemaFieldBurnRateThresholdOperator represents the threshold_operator field identifier
	SchemaFieldBurnRateThresholdOperator = "threshold_operator"
	// SchemaFieldBurnRateThresholdValue represents the threshold_value field identifier
	SchemaFieldBurnRateThresholdValue = "threshold_value"

	// Validation constants

	// NameMaxLength represents the maximum length for the name field
	NameMaxLength = 256
	// NameMinLength represents the minimum length for the name field
	NameMinLength = 0
	// TimeThresholdMinValue represents the minimum value for time threshold fields
	TimeThresholdMinValue = 1
	// FloatFormatPrecision represents the precision for float formatting
	FloatFormatPrecision = 2

	// Alternative alert type naming constants for normalization

	// AlertTypeErrorBudgetAlt1 represents alternative naming for error_budget
	AlertTypeErrorBudgetAlt1 = "errorBudget"
	// AlertTypeErrorBudgetAlt2 represents alternative naming for error_budget
	AlertTypeErrorBudgetAlt2 = "ErrorBudget"
	// AlertTypeStatusAlt1 represents alternative naming for status
	AlertTypeStatusAlt1 = "Status"
	// AlertTypeBurnRateV2Alt1 represents alternative naming for burn_rate_v2
	AlertTypeBurnRateV2Alt1 = "burnRateV2"
	// AlertTypeBurnRateV2Alt2 represents alternative naming for burn_rate_v2
	AlertTypeBurnRateV2Alt2 = "BurnRateV2"
)
