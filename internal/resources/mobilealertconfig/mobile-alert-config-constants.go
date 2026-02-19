package mobilealertconfig

// ResourceInstanaMobileAlertConfig the name of the terraform-provider-instana resource to manage mobile alert configs
const ResourceInstanaMobileAlertConfig = "mobile_alert_config"

const (
	// MobileAlertConfigFieldID constant value for field id of resource instana_mobile_alert_config
	MobileAlertConfigFieldID = "id"
	// MobileAlertConfigFieldName constant value for field name of resource instana_mobile_alert_config
	MobileAlertConfigFieldName = "name"
	// MobileAlertConfigFieldDescription constant value for field description of resource instana_mobile_alert_config
	MobileAlertConfigFieldDescription = "description"
	// MobileAlertConfigFieldMobileAppID constant value for field mobile_app_id of resource instana_mobile_alert_config
	MobileAlertConfigFieldMobileAppID = "mobile_app_id"
	// MobileAlertConfigFieldSeverity constant value for field severity of resource instana_mobile_alert_config
	MobileAlertConfigFieldSeverity = "severity"
	// MobileAlertConfigFieldTriggering constant value for field triggering of resource instana_mobile_alert_config
	MobileAlertConfigFieldTriggering = "triggering"
	// MobileAlertConfigFieldEnabled constant value for field enabled of resource instana_mobile_alert_config
	MobileAlertConfigFieldEnabled = "enabled"
	// MobileAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_mobile_alert_config
	MobileAlertConfigFieldTagFilter = "tag_filter"
	// MobileAlertConfigFieldAlertChannels constant value for field alert_channels of resource instana_mobile_alert_config
	MobileAlertConfigFieldAlertChannels = "alert_channels"
	// MobileAlertConfigFieldGranularity constant value for field granularity of resource instana_mobile_alert_config
	MobileAlertConfigFieldGranularity = "granularity"
	// MobileAlertConfigFieldGracePeriod constant value for field grace_period of resource instana_mobile_alert_config
	MobileAlertConfigFieldGracePeriod = "grace_period"
	// MobileAlertConfigFieldCustomPayloadFields constant value for field custom_payload_field of resource instana_mobile_alert_config
	MobileAlertConfigFieldCustomPayloadFields = "custom_payload_field"
	// MobileAlertConfigFieldRules constant value for field rules of resource instana_mobile_alert_config
	MobileAlertConfigFieldRules = "rules"
	// MobileAlertConfigFieldRule constant value for field rule of resource instana_mobile_alert_config
	MobileAlertConfigFieldRule = "rule"
	// MobileAlertConfigFieldThresholdOperator constant value for field threshold_operator of resource instana_mobile_alert_config
	MobileAlertConfigFieldThresholdOperator = "threshold_operator"
	// MobileAlertConfigFieldThreshold constant value for field threshold of resource instana_mobile_alert_config
	MobileAlertConfigFieldThreshold = "threshold"
	// MobileAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_mobile_alert_config
	MobileAlertConfigFieldTimeThreshold = "time_threshold"
	// MobileAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_mobile_alert_config
	MobileAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	// MobileAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence constant value for field time_threshold.user_impact_of_violations_in_sequence
	MobileAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence = "user_impact_of_violations_in_sequence"
	// MobileAlertConfigFieldTimeThresholdUserImpactUsers constant value for field time_threshold.user_impact_of_violations_in_sequence.users
	MobileAlertConfigFieldTimeThresholdUserImpactUsers = "users"
	// MobileAlertConfigFieldTimeThresholdUserImpactPercentage constant value for field time_threshold.user_impact_of_violations_in_sequence.percentage
	MobileAlertConfigFieldTimeThresholdUserImpactPercentage = "percentage"
	// MobileAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period
	MobileAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	// MobileAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations
	MobileAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	// MobileAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence
	MobileAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	// MobileAlertConfigFieldRuleAlertType constant value for field rule.alert_type
	MobileAlertConfigFieldRuleAlertType = "alert_type"
	// MobileAlertConfigFieldRuleMetricName constant value for field rule.metric_name
	MobileAlertConfigFieldRuleMetricName = "metric_name"
	// MobileAlertConfigFieldRuleAggregation constant value for field rule.aggregation
	MobileAlertConfigFieldRuleAggregation = "aggregation"
	// MobileAlertConfigFieldRuleOperator constant value for field rule.operator
	MobileAlertConfigFieldRuleOperator = "operator"
	// MobileAlertConfigFieldRuleValue constant value for field rule.value
	MobileAlertConfigFieldRuleValue = "value"
	// MobileAlertConfigFieldRuleCustomEventName constant value for field rule.custom_event_name
	MobileAlertConfigFieldRuleCustomEventName = "custom_event_name"
)

// Resource description constants
const (
	MobileAlertConfigDescResource                          = "This resource manages Mobile App Alert Configurations in Instana."
	MobileAlertConfigDescID                                = "The ID of the Mobile Alert Configuration."
	MobileAlertConfigDescName                              = "The name of the Mobile Alert Configuration."
	MobileAlertConfigDescDescription                       = "The description of the Mobile Alert Configuration."
	MobileAlertConfigDescMobileAppID                       = "ID of the mobile app that this Smart Alert configuration is applied to."
	MobileAlertConfigDescSeverity                          = "The severity of the alert when triggered (5 for Warning, 10 for Critical). Deprecated - use rules with thresholds instead."
	MobileAlertConfigDescTriggering                        = "Flag to indicate whether an Incident is also triggered or not."
	MobileAlertConfigDescEnabled                           = "Flag to enable or disable the alert configuration."
	MobileAlertConfigDescTagFilter                         = "The tag filter expression for the Mobile Alert Configuration."
	MobileAlertConfigDescAlertChannels                     = "Set of alert channel IDs associated with the severity."
	MobileAlertConfigDescGranularity                       = "The evaluation granularity used for detection of violations of the defined threshold."
	MobileAlertConfigDescGracePeriod                       = "The duration for which an alert remains open after conditions are no longer violated."
	MobileAlertConfigDescRules                             = "A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels."
	MobileAlertConfigDescRule                              = "The mobile app alert rule configuration."
	MobileAlertConfigDescThresholdOperator                 = "The operator to apply for threshold comparison."
	MobileAlertConfigDescThreshold                         = "Threshold configuration for different severity levels."
	MobileAlertConfigDescTimeThreshold                     = "The type of threshold to define the criteria when the event and alert triggers and resolves."
	MobileAlertConfigDescTimeThresholdUserImpact           = "Time threshold based on user impact of violations in sequence."
	MobileAlertConfigDescTimeThresholdTimeWindow           = "The time window of the time threshold."
	MobileAlertConfigDescTimeThresholdUsers                = "The number of impacted users."
	MobileAlertConfigDescTimeThresholdPercentage           = "The percentage of impacted users."
	MobileAlertConfigDescTimeThresholdViolationsInPeriod   = "Time threshold based on violations in period."
	MobileAlertConfigDescTimeThresholdViolations           = "The violations appeared in the period."
	MobileAlertConfigDescTimeThresholdViolationsInSequence = "Time threshold based on violations in sequence."
	MobileAlertConfigDescRuleAlertType                     = "The type of alert rule."
	MobileAlertConfigDescRuleMetricName                    = "The metric name of the mobile alert rule."
	MobileAlertConfigDescRuleAggregation                   = "The aggregation function of the mobile alert rule."
	MobileAlertConfigDescRuleOperator                      = "The operator for the rule. Valid values are STARTS_WITH and EQUALS."
	MobileAlertConfigDescRuleValue                         = "The value to compare against."
	MobileAlertConfigDescRuleCustomEventName               = "The name of the custom event to monitor. Required when alert_type is 'customEvent'."
)

// Error message constants
const (
	MobileAlertConfigErrParseTagFilter                = "Error parsing tag filter"
	MobileAlertConfigErrMapFilterExpression           = "Error mapping filter expression"
	MobileAlertConfigErrMapFilterExpressionMsg        = "Failed to map filter expression: %s"
	MobileAlertConfigErrTimeThresholdRequired         = "Time threshold is required"
	MobileAlertConfigErrTimeThresholdRequiredMsg      = "Mobile alert config time threshold is required"
	MobileAlertConfigErrInvalidTimeThresholdConfig    = "Invalid time threshold configuration"
	MobileAlertConfigErrInvalidTimeThresholdConfigMsg = "Exactly one time threshold type configuration is required"
)

// Time threshold type constants
const (
	MobileAlertConfigTimeThresholdTypeViolationsInSequence             = "violationsInSequence"
	MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence = "userImpactOfViolationsInSequence"
	MobileAlertConfigTimeThresholdTypeViolationsInPeriod               = "violationsInPeriod"
)

// Tag filter constants
const (
	MobileAlertConfigTagFilterTypeExpression = "EXPRESSION"
	MobileAlertConfigLogicalOperatorAND      = "AND"
)

// Default value constants
const (
	MobileAlertConfigDefaultTriggering  = false
	MobileAlertConfigDefaultEnabled     = true
	MobileAlertConfigDefaultGranularity = 600000
)

// Validation limit constants
const (
	MobileAlertConfigMinNameLength        = 0
	MobileAlertConfigMaxNameLength        = 256
	MobileAlertConfigMinDescriptionLength = 0
	MobileAlertConfigMaxDescriptionLength = 65536
	MobileAlertConfigMinMobileAppIDLength = 0
	MobileAlertConfigMaxMobileAppIDLength = 64
	MobileAlertConfigMinSeverity          = 5
	MobileAlertConfigMaxSeverity          = 10
)

// Aggregation value constants
const (
	MobileAlertConfigAggregationSUM  = "SUM"
	MobileAlertConfigAggregationMEAN = "MEAN"
	MobileAlertConfigAggregationMAX  = "MAX"
	MobileAlertConfigAggregationMIN  = "MIN"
	MobileAlertConfigAggregationP25  = "P25"
	MobileAlertConfigAggregationP50  = "P50"
	MobileAlertConfigAggregationP75  = "P75"
	MobileAlertConfigAggregationP90  = "P90"
	MobileAlertConfigAggregationP95  = "P95"
	MobileAlertConfigAggregationP98  = "P98"
	MobileAlertConfigAggregationP99  = "P99"
)
