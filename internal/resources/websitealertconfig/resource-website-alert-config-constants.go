package websitealertconfig

// ResourceInstanaWebsiteAlertConfig the name of the terraform-provider-instana resource to manage website alert configs
const ResourceInstanaWebsiteAlertConfig = "website_alert_config"

const (
	//WebsiteAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_website_alert_config
	WebsiteAlertConfigFieldAlertChannelIDs = "alert_channel_ids"
	//WebsiteAlertConfigFieldWebsiteID constant value for field websites.website_id of resource instana_website_alert_config
	WebsiteAlertConfigFieldWebsiteID = "website_id"
	//WebsiteAlertConfigFieldDescription constant value for field description of resource instana_website_alert_config
	WebsiteAlertConfigFieldDescription = "description"
	//WebsiteAlertConfigFieldGranularity constant value for field granularity of resource instana_website_alert_config
	WebsiteAlertConfigFieldGranularity = "granularity"
	//WebsiteAlertConfigFieldName constant value for field name of resource instana_website_alert_config
	WebsiteAlertConfigFieldName = "name"
	//WebsiteAlertConfigFieldEnabled constant value for field enabled of resource instana_website_alert_config
	WebsiteAlertConfigFieldEnabled = "enabled"

	//WebsiteAlertConfigFieldRule constant value for field rule of resource instana_website_alert_config
	WebsiteAlertConfigFieldRule = "rule"
	//WebsiteAlertConfigFieldRuleMetricName constant value for field rule.*.metric_name of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleMetricName = "metric_name"
	//WebsiteAlertConfigFieldRuleAggregation constant value for field rule.*.aggregation of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleAggregation = "aggregation"
	//WebsiteAlertConfigFieldRuleOperator constant value for field rule.*.operator of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleOperator = "operator"
	//WebsiteAlertConfigFieldRuleValue constant value for field rule.*.value of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleValue = "value"
	//WebsiteAlertConfigFieldRuleSlowness constant value for field rule.slowness of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleSlowness = "slowness"
	//WebsiteAlertConfigFieldRuleStatusCode constant value for field rule.status_code of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleStatusCode = "status_code"
	//WebsiteAlertConfigFieldRuleThroughput constant value for field rule.throughput of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleThroughput = "throughput"
	//WebsiteAlertConfigFieldRuleSpecificJsError constant value for field rule.specific_js_error of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleSpecificJsError = "specific_js_error"

	//WebsiteAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_website_alert_config
	WebsiteAlertConfigFieldTagFilter = "tag_filter"

	//WebsiteAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThreshold = "time_threshold"
	//WebsiteAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence constant value for field time_threshold.user_impact_of_violations_in_sequence of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence = "user_impact_of_violations_in_sequence"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod constant value for field time_threshold.user_impact_of_violations_in_sequence.impact_measurement_method of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod = "impact_measurement_method"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage constant value for field time_threshold.user_impact_of_violations_in_sequence.user_percentage of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage = "user_percentage"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers constant value for field time_threshold.user_impact_of_violations_in_sequence.users of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers = "users"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	//WebsiteAlertConfigFieldTriggering constant value for field triggering of resource instana_website_alert_config
	WebsiteAlertConfigFieldTriggering = "triggering"
	WebsiteAlertConfigFieldRules      = "rules"
	WebsiteAlertConfigFieldThreshold  = "threshold"

	// Resource description constants
	WebsiteAlertConfigDescResource        = "This resource manages Website Alert Configurations in Instana."
	WebsiteAlertConfigDescID              = "The ID of the Website Alert Configuration."
	WebsiteAlertConfigDescName            = "The name of the Website Alert Configuration."
	WebsiteAlertConfigDescDescription     = "The description of the Website Alert Configuration."
	WebsiteAlertConfigDescTriggering      = "Flag to indicate whether also an Incident is triggered or not."
	WebsiteAlertConfigDescEnabled         = "Flag to enable or disable the smart alert."
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

// Field name constants
const (
	WebsiteAlertConfigFieldID                  = "id"
	WebsiteAlertConfigFieldStatic              = "static"
	WebsiteAlertConfigFieldAdaptiveBaseline    = "adaptive_baseline"
	WebsiteAlertConfigFieldHistoricBaseline    = "historic_baseline"
	WebsiteAlertConfigFieldDeviationFactor     = "deviation_factor"
	WebsiteAlertConfigFieldAdaptability        = "adaptability"
	WebsiteAlertConfigFieldSeasonality         = "seasonality"
	WebsiteAlertConfigFieldCustomPayloadFields = "custom_payload_fields"
)

// Operator value constants
const (
	WebsiteAlertConfigOperatorGreaterThan       = ">"
	WebsiteAlertConfigOperatorGreaterThanEquals = ">="
	WebsiteAlertConfigOperatorLessThan          = "<"
	WebsiteAlertConfigOperatorLessThanEquals    = "<="
	WebsiteAlertConfigOperatorEquals            = "=="
)

// Aggregation value constants
const (
	WebsiteAlertConfigAggregationSUM  = "SUM"
	WebsiteAlertConfigAggregationMEAN = "MEAN"
	WebsiteAlertConfigAggregationMAX  = "MAX"
	WebsiteAlertConfigAggregationMIN  = "MIN"
	WebsiteAlertConfigAggregationP25  = "P25"
	WebsiteAlertConfigAggregationP50  = "P50"
	WebsiteAlertConfigAggregationP75  = "P75"
	WebsiteAlertConfigAggregationP90  = "P90"
	WebsiteAlertConfigAggregationP95  = "P95"
	WebsiteAlertConfigAggregationP98  = "P98"
	WebsiteAlertConfigAggregationP99  = "P99"
)

// Impact measurement method constants
const (
	WebsiteAlertConfigImpactMethodAggregated = "AGGREGATED"
	WebsiteAlertConfigImpactMethodPerWindow  = "PER_WINDOW"
)

// Alert type constants
const (
	WebsiteAlertConfigAlertTypeThroughput      = "throughput"
	WebsiteAlertConfigAlertTypeSlowness        = "slowness"
	WebsiteAlertConfigAlertTypeStatusCode      = "statusCode"
	WebsiteAlertConfigAlertTypeSpecificJsError = "specificJsError"
)

// Threshold type constants
const (
	WebsiteAlertConfigThresholdTypeStatic   = "staticThreshold"
	WebsiteAlertConfigThresholdTypeAdaptive = "adaptiveBaseline"
	WebsiteAlertConfigThresholdTypeHistoric = "historicBaseline"
)

// Time threshold type constants
const (
	WebsiteAlertConfigTimeThresholdTypeViolationsInSequence = "violationsInSequence"
	WebsiteAlertConfigTimeThresholdTypeUserImpact           = "userImpactOfViolationsInSequence"
	WebsiteAlertConfigTimeThresholdTypeViolationsInPeriod   = "violationsInPeriod"
)

// Tag filter constants
const (
	WebsiteAlertConfigTagFilterTypeExpression = "EXPRESSION"
	WebsiteAlertConfigLogicalOperatorAND      = "AND"
)

// Default value constants
const (
	WebsiteAlertConfigDefaultTriggering  = false
	WebsiteAlertConfigDefaultEnabled     = true
	WebsiteAlertConfigDefaultGranularity = 600000
)

// Validation limit constants
const (
	WebsiteAlertConfigMinNameLength        = 0
	WebsiteAlertConfigMaxNameLength        = 256
	WebsiteAlertConfigMinDescriptionLength = 0
	WebsiteAlertConfigMaxDescriptionLength = 65536
	WebsiteAlertConfigMinWebsiteIDLength   = 0
	WebsiteAlertConfigMaxWebsiteIDLength   = 64
)
