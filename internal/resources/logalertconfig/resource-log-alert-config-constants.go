package logalertconfig

import "github.com/instana/terraform-provider-instana/internal/restapi"

// ResourceInstanaLogAlertConfig the name of the terraform-provider-instana resource to manage log alert configurations
const ResourceInstanaLogAlertConfig = "log_alert_config"

// Schema field names
const (
	// LogAlertConfigFieldID constant value for the schema field id
	LogAlertConfigFieldID = "id"
	// LogAlertConfigFieldName constant value for the schema field name
	LogAlertConfigFieldName = "name"
	// LogAlertConfigFieldDescription constant value for the schema field description
	LogAlertConfigFieldDescription = "description"
	// LogAlertConfigFieldAlertChannels constant value for the schema field alert_channels
	LogAlertConfigFieldAlertChannels = "alert_channels"
	// LogAlertConfigFieldGracePeriod constant value for the schema field grace_period
	LogAlertConfigFieldGracePeriod = "grace_period"
	// LogAlertConfigFieldGroupBy constant value for the schema field group_by
	LogAlertConfigFieldGroupBy = "group_by"
	// LogAlertConfigFieldGroupByTagName constant value for the schema field tag_name
	LogAlertConfigFieldGroupByTagName = "tag_name"
	// LogAlertConfigFieldGroupByKey constant value for the schema field key
	LogAlertConfigFieldGroupByKey = "key"
	// LogAlertConfigFieldGranularity constant value for the schema field granularity
	LogAlertConfigFieldGranularity = "granularity"
	// LogAlertConfigFieldTagFilter constant value for the schema field tag_filter
	LogAlertConfigFieldTagFilter = "tag_filter"

	// LogAlertConfigFieldRules constant value for the schema field rules
	LogAlertConfigFieldRules = "rules"
	// LogAlertConfigFieldAlertType constant value for the schema field alert_type
	LogAlertConfigFieldAlertType = "alert_type"
	// LogAlertConfigFieldMetricName constant value for the schema field metric_name
	LogAlertConfigFieldMetricName = "metric_name"
	// LogAlertConfigFieldAggregation constant value for the schema field aggregation
	LogAlertConfigFieldAggregation = "aggregation"
	// LogAlertConfigFieldThresholdOperator constant value for the schema field threshold_operator
	LogAlertConfigFieldThresholdOperator = "threshold_operator"
	// LogAlertConfigFieldThreshold constant value for the schema field threshold
	LogAlertConfigFieldThreshold = "threshold"
	// LogAlertConfigFieldWarning constant value for the schema field warning
	LogAlertConfigFieldWarning = "warning"
	// LogAlertConfigFieldCritical constant value for the schema field critical
	LogAlertConfigFieldCritical = "critical"
	// LogAlertConfigFieldType constant value for the schema field type
	LogAlertConfigFieldType = "type"
	// LogAlertConfigFieldValue constant value for the schema field value
	LogAlertConfigFieldValue = "value"

	// LogAlertConfigFieldTimeThreshold constant value for the schema field time_threshold
	LogAlertConfigFieldTimeThreshold = "time_threshold"
	// LogAlertConfigFieldTimeThresholdTimeWindow constant value for the schema field time_window
	LogAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	// LogAlertConfigFieldTimeThresholdViolationsInSequence constant value for the schema field violations_in_sequence
	LogAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
)

// API-related constants
const (
	// LogAlertTypeLogCount is the constant for the log count alert type in schema
	LogAlertTypeLogCount = "log.count"
	// LogAlertTypeLogCountAPI is the constant for the log count alert type in API
	LogAlertTypeLogCountAPI = "logCount"
	// TimeThresholdTypeViolationsInSequence represents the violations in sequence time threshold type
	TimeThresholdTypeViolationsInSequence = "violationsInSequence"
	// ThresholdTypeStatic represents the static threshold type
	ThresholdTypeStatic = "staticThreshold"
	// EmptyString represents an empty string value
	EmptyString = ""
	// DefaultGranularity represents the default granularity value
	DefaultGranularity = int64(restapi.Granularity600000)
)

// Resource description constants
const (
	// LogAlertConfigDescResource description for the log alert config resource
	LogAlertConfigDescResource = "This resource manages log alert configurations in Instana."
	// LogAlertConfigDescID description for the ID field
	LogAlertConfigDescID = "The ID of the log alert configuration."
	// LogAlertConfigDescName description for the name field
	LogAlertConfigDescName = "Name for the Log alert configuration"
	// LogAlertConfigDescDescription description for the description field
	LogAlertConfigDescDescription = "The description text of the Log alert config"
	// LogAlertConfigDescGracePeriod description for the grace_period field
	LogAlertConfigDescGracePeriod = "The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires."
	// LogAlertConfigDescGranularity description for the granularity field
	LogAlertConfigDescGranularity = "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used"
	// LogAlertConfigDescTagFilter description for the tag_filter field
	LogAlertConfigDescTagFilter = "The tag filter expression used for this log alert"
	// LogAlertConfigDescAlertChannels description for the alert_channels field
	LogAlertConfigDescAlertChannels = "Set of alert channel IDs associated with the severity."
	// LogAlertConfigDescAlertChannelIDs description for alert channel ID lists
	LogAlertConfigDescAlertChannelIDs = "List of IDs of alert channels defined in Instana."
	// LogAlertConfigDescGroupBy description for the group_by field
	LogAlertConfigDescGroupBy = "The grouping tags used to group the metric results."
	// LogAlertConfigDescGroupByTagName description for the tag_name field
	LogAlertConfigDescGroupByTagName = "The tag name used for grouping"
	// LogAlertConfigDescGroupByKey description for the key field
	LogAlertConfigDescGroupByKey = "The key used for grouping"
	// LogAlertConfigDescRules description for the rules field
	LogAlertConfigDescRules = "Configuration for the log alert rule"
	// LogAlertConfigDescMetricName description for the metric_name field
	LogAlertConfigDescMetricName = "The metric name of the log alert rule"
	// LogAlertConfigDescAlertType description for the alert_type field
	LogAlertConfigDescAlertType = "The type of the log alert rule (only 'log.count' is supported)"
	// LogAlertConfigDescAggregation description for the aggregation field
	LogAlertConfigDescAggregation = "The aggregation method to use for the log alert (only 'SUM' is supported)"
	// LogAlertConfigDescThresholdOperator description for the threshold_operator field
	LogAlertConfigDescThresholdOperator = "The operator which will be applied to evaluate the threshold"
	// LogAlertConfigDescThreshold description for the threshold field
	LogAlertConfigDescThreshold = "Threshold configuration for different severity levels"
	// LogAlertConfigDescTimeThreshold description for the time_threshold field
	LogAlertConfigDescTimeThreshold = "Indicates the type of violation of the defined threshold."
	// LogAlertConfigDescViolationsInSequence description for the violations_in_sequence field
	LogAlertConfigDescViolationsInSequence = "Time threshold base on violations in sequence"
	// LogAlertConfigDescTimeWindow description for the time_window field
	LogAlertConfigDescTimeWindow = "Time window in milliseconds."
)

// Error message constants
const (
	// LogAlertConfigErrNormalizingTagFilter error title for tag filter normalization failures
	LogAlertConfigErrNormalizingTagFilter = "Error normalizing tag filter"
	// LogAlertConfigErrNormalizingTagFilterMsg error message template for tag filter normalization failures
	LogAlertConfigErrNormalizingTagFilterMsg = "Could not normalize tag filter: "
	// LogAlertConfigErrParsingTagFilter error title for tag filter parsing failures
	LogAlertConfigErrParsingTagFilter = "Error parsing tag filter"
	// LogAlertConfigErrParsingTagFilterMsg error message template for tag filter parsing failures
	LogAlertConfigErrParsingTagFilterMsg = "Could not parse tag filter: "
)

// Made with Bob
