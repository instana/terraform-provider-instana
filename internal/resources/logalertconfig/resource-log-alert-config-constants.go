package logalertconfig

// ResourceInstanaLogAlertConfigFramework the name of the terraform-provider-instana resource to manage log alert configurations
const ResourceInstanaLogAlertConfigFramework = "log_alert_config"

const (
	LogAlertConfigFieldName           = "name"
	LogAlertConfigFieldDescription    = "description"
	LogAlertConfigFieldAlertChannels  = "alert_channels"
	LogAlertConfigFieldGracePeriod    = "grace_period"
	LogAlertConfigFieldGroupBy        = "group_by"
	LogAlertConfigFieldGroupByTagName = "tag_name"
	LogAlertConfigFieldGroupByKey     = "key"
	LogAlertConfigFieldGranularity    = "granularity"
	LogAlertConfigFieldTagFilter      = "tag_filter"

	LogAlertConfigFieldRules             = "rules"
	LogAlertConfigFieldAlertType         = "alert_type"
	LogAlertConfigFieldMetricName        = "metric_name"
	LogAlertConfigFieldAggregation       = "aggregation"
	LogAlertConfigFieldThresholdOperator = "threshold_operator"
	LogAlertConfigFieldThreshold         = "threshold"
	LogAlertConfigFieldWarning           = "warning"
	LogAlertConfigFieldCritical          = "critical"
	LogAlertConfigFieldType              = "type"
	LogAlertConfigFieldValue             = "value"

	// LogAlertTypeLogCount is the constant for the log count alert type
	LogAlertTypeLogCount = "log.count"

	LogAlertConfigFieldTimeThreshold                     = "time_threshold"
	LogAlertConfigFieldTimeThresholdTimeWindow           = "time_window"
	LogAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
)

// Resource description constants
const (
	LogAlertConfigDescResource             = "This resource manages log alert configurations in Instana."
	LogAlertConfigDescID                   = "The ID of the log alert configuration."
	LogAlertConfigDescName                 = "Name for the Log alert configuration"
	LogAlertConfigDescDescription          = "The description text of the Log alert config"
	LogAlertConfigDescGracePeriod          = "The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires."
	LogAlertConfigDescGranularity          = "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used"
	LogAlertConfigDescTagFilter            = "The tag filter expression used for this log alert"
	LogAlertConfigDescAlertChannels        = "Set of alert channel IDs associated with the severity."
	LogAlertConfigDescAlertChannelIDs      = "List of IDs of alert channels defined in Instana."
	LogAlertConfigDescGroupBy              = "The grouping tags used to group the metric results."
	LogAlertConfigDescGroupByTagName       = "The tag name used for grouping"
	LogAlertConfigDescGroupByKey           = "The key used for grouping"
	LogAlertConfigDescRules                = "Configuration for the log alert rule"
	LogAlertConfigDescMetricName           = "The metric name of the log alert rule"
	LogAlertConfigDescAlertType            = "The type of the log alert rule (only 'log.count' is supported)"
	LogAlertConfigDescAggregation          = "The aggregation method to use for the log alert (only 'SUM' is supported)"
	LogAlertConfigDescThresholdOperator    = "The operator which will be applied to evaluate the threshold"
	LogAlertConfigDescThreshold            = "Threshold configuration for different severity levels"
	LogAlertConfigDescTimeThreshold        = "Indicates the type of violation of the defined threshold."
	LogAlertConfigDescViolationsInSequence = "Time threshold base on violations in sequence"
	LogAlertConfigDescTimeWindow           = "Time window in milliseconds."
)

// Error message constants
const (
	LogAlertConfigErrNormalizingTagFilter    = "Error normalizing tag filter"
	LogAlertConfigErrNormalizingTagFilterMsg = "Could not normalize tag filter: %s"
	LogAlertConfigErrParsingTagFilter        = "Error parsing tag filter"
	LogAlertConfigErrParsingTagFilterMsg     = "Could not parse tag filter: %s"
)
