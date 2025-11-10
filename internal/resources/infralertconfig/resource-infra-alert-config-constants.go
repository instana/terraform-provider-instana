package instana

// Resource description constants
const (
	InfraAlertConfigDescResource               = "This resource represents an infrastructure alert configuration in Instana"
	InfraAlertConfigDescID                     = "The ID of the infrastructure alert configuration"
	InfraAlertConfigDescName                   = "The name of the infrastructure alert configuration"
	InfraAlertConfigDescDescription            = "The description of the infrastructure alert configuration"
	InfraAlertConfigDescTagFilter              = "The tag filter expression for the infrastructure alert configuration"
	InfraAlertConfigDescGroupBy                = "The list of tags to group by"
	InfraAlertConfigDescGranularity            = "The granularity of the infrastructure alert configuration"
	InfraAlertConfigDescEvaluationType         = "The evaluation type of the infrastructure alert configuration"
	InfraAlertConfigDescRules                  = "The rules configuration"
	InfraAlertConfigDescGenericRule            = "The generic rule configuration"
	InfraAlertConfigDescMetricName             = "The metric name for the generic rule"
	InfraAlertConfigDescEntityType             = "The entity type for the generic rule"
	InfraAlertConfigDescAggregation            = "The aggregation for the generic rule"
	InfraAlertConfigDescCrossSeriesAggregation = "The cross series aggregation for the generic rule"
	InfraAlertConfigDescRegex                  = "Whether regex is enabled for the generic rule"
	InfraAlertConfigDescThresholdOperator      = "The threshold operator for the generic rule"
	InfraAlertConfigDescThreshold              = "Threshold configuration for different severity levels"
	InfraAlertConfigDescTimeThreshold          = "Indicates the type of violation of the defined threshold."
	InfraAlertConfigDescViolationsInSequence   = "Time threshold base on violations in sequence"
	InfraAlertConfigDescTimeWindow             = "The time window if the time threshold"
	InfraAlertConfigDescAlertChannels          = "Set of alert channel IDs associated with the severity."
	InfraAlertConfigDescAlertChannelIDs        = "List of IDs of alert channels defined in Instana."
)

// Error message constants
const (
	InfraAlertConfigErrMappingTagFilter    = "Error mapping tag filter"
	InfraAlertConfigErrMappingTagFilterMsg = "Failed to map tag filter: %s"
	InfraAlertConfigErrParsingTagFilter    = "Error parsing tag filter"
	InfraAlertConfigErrParsingTagFilterMsg = "Failed to parse tag filter: %s"
)

// Made with Bob
