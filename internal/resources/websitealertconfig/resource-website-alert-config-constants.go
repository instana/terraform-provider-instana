package websitealertconfig

const (
	// Resource description constants
	WebsiteAlertConfigDescResource        = "This resource manages Website Alert Configurations in Instana."
	WebsiteAlertConfigDescID              = "The ID of the Website Alert Configuration."
	WebsiteAlertConfigDescName            = "The name of the Website Alert Configuration."
	WebsiteAlertConfigDescDescription     = "The description of the Website Alert Configuration."
	WebsiteAlertConfigDescSeverity        = "The severity of the alert when triggered."
	WebsiteAlertConfigDescTriggering      = "Flag to indicate whether also an Incident is triggered or not."
	WebsiteAlertConfigDescWebsiteID       = "Unique ID of the website."
	WebsiteAlertConfigDescTagFilter       = "The tag filter expression for the Website Alert Configuration."
	WebsiteAlertConfigDescAlertChannelIDs = "List of IDs of alert channels defined in Instana."
	WebsiteAlertConfigDescGranularity     = "The evaluation granularity used for detection of violations of the defined threshold."
	WebsiteAlertConfigDescRules           = "A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels."

	// Rule descriptions
	WebsiteAlertConfigDescRuleOperator        = "The operator to apply for threshold comparison"
	WebsiteAlertConfigDescRule                = "Indicates the type of rule this alert configuration is about."
	WebsiteAlertConfigDescRuleSlowness        = "Rule based on the slowness of the configured alert configuration target."
	WebsiteAlertConfigDescRuleSpecificJsError = "Rule based on a specific javascript error of the configured alert configuration target."
	WebsiteAlertConfigDescRuleStatusCode      = "Rule based on the HTTP status code of the configured alert configuration target."
	WebsiteAlertConfigDescRuleThroughput      = "Rule based on the throughput of the configured alert configuration target."
	WebsiteAlertConfigDescRuleMetricName      = "The metric name of the website alert rule."
	WebsiteAlertConfigDescRuleAggregation     = "The aggregation function of the website alert rule."
	WebsiteAlertConfigDescRuleOperatorEval    = "The operator which will be applied to evaluate this rule."
	WebsiteAlertConfigDescRuleValueJsError    = "The value identify the specific javascript error."
	WebsiteAlertConfigDescRuleValueStatusCode = "The value identify the specific http status code."

	// Threshold descriptions
	WebsiteAlertConfigDescThreshold                = "Threshold configuration for different severity levels"
	WebsiteAlertConfigDescThresholdStatic          = "Static threshold definition."
	WebsiteAlertConfigDescThresholdOperator        = "Comparison operator for the static threshold."
	WebsiteAlertConfigDescThresholdValue           = "The numeric value for the static threshold."
	WebsiteAlertConfigDescThresholdAdaptive        = "Static threshold definition."
	WebsiteAlertConfigDescThresholdDeviationFactor = "The numeric value for the deviation factor."
	WebsiteAlertConfigDescThresholdAdaptability    = "The numeric value for the adaptability."
	WebsiteAlertConfigDescThresholdSeasonality     = "Value for the seasonality."

	// Time threshold descriptions
	WebsiteAlertConfigDescTimeThreshold                     = "Indicates the type of violation of the defined threshold."
	WebsiteAlertConfigDescTimeThresholdUserImpact           = "Time threshold base on user impact of violations in sequence."
	WebsiteAlertConfigDescTimeThresholdTimeWindow           = "The time window if the time threshold."
	WebsiteAlertConfigDescTimeThresholdImpactMethod         = "The impact method of the time threshold based on user impact of violations in sequence."
	WebsiteAlertConfigDescTimeThresholdUserPercentage       = "The percentage of impacted users of the time threshold based on user impact of violations in sequence."
	WebsiteAlertConfigDescTimeThresholdUsers                = "The number of impacted users of the time threshold based on user impact of violations in sequence."
	WebsiteAlertConfigDescTimeThresholdViolationsInPeriod   = "Time threshold base on violations in period."
	WebsiteAlertConfigDescTimeThresholdViolations           = "The violations appeared in the period."
	WebsiteAlertConfigDescTimeThresholdViolationsInSequence = "Time threshold base on violations in sequence."

	// Error message constants
	WebsiteAlertConfigErrConvertSeverity               = "Error converting severity"
	WebsiteAlertConfigErrParseTagFilter                = "Error parsing tag filter"
	WebsiteAlertConfigErrMapFilterExpression           = "Error mapping filter expression"
	WebsiteAlertConfigErrMapFilterExpressionMsg        = "Failed to map filter expression: %s"
	WebsiteAlertConfigErrTimeThresholdRequired         = "Time threshold is required"
	WebsiteAlertConfigErrTimeThresholdRequiredMsg      = "Website alert config time threshold is required"
	WebsiteAlertConfigErrInvalidRuleConfig             = "Invalid rule configuration"
	WebsiteAlertConfigErrInvalidRuleConfigMsg          = "Exactly one rule type configuration is required"
	WebsiteAlertConfigErrInvalidTimeThresholdConfig    = "Invalid time threshold configuration"
	WebsiteAlertConfigErrInvalidTimeThresholdConfigMsg = "Exactly one time threshold type configuration is required"
)

// Made with Bob
