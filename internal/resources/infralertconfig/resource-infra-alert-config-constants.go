package infralertconfig

// ResourceInstanaInfraAlertConfig the name of the terraform-provider-instana resource to manage infrastructure alert configurations
const ResourceInstanaInfraAlertConfig = "infra_alert_config"

// Schema field names
const (
	// InfraAlertConfigFieldID constant value for the schema field id
	InfraAlertConfigFieldID = "id"
	// InfraAlertConfigFieldName constant value for the schema field name
	InfraAlertConfigFieldName = "name"
	// InfraAlertConfigFieldDescription constant value for the schema field description
	InfraAlertConfigFieldDescription = "description"
	// InfraAlertConfigFieldTagFilter constant value for the schema field tag_filter
	InfraAlertConfigFieldTagFilter = "tag_filter"
	// InfraAlertConfigFieldGroupBy constant value for the schema field group_by
	InfraAlertConfigFieldGroupBy = "group_by"
	// InfraAlertConfigFieldGranularity constant value for the schema field granularity
	InfraAlertConfigFieldGranularity = "granularity"
	// InfraAlertConfigFieldEvaluationType constant value for the schema field evaluation_type
	InfraAlertConfigFieldEvaluationType = "evaluation_type"
	// InfraAlertConfigFieldCustomPayloadField constant value for the schema field custom_payload_field
	InfraAlertConfigFieldCustomPayloadField = "custom_payload_field"
	// InfraAlertConfigFieldAlertChannels constant value for the schema field alert_channels
	InfraAlertConfigFieldAlertChannels = "alert_channels"
	// InfraAlertConfigFieldTimeThreshold constant value for the schema field time_threshold
	InfraAlertConfigFieldTimeThreshold = "time_threshold"
	// InfraAlertConfigFieldRules constant value for the schema field rules
	InfraAlertConfigFieldRules = "rules"

	// Alert channel severity fields
	// ResourceFieldThresholdRuleWarningSeverity constant value for warning severity
	ResourceFieldThresholdRuleWarningSeverity = "warning"
	// ResourceFieldThresholdRuleCriticalSeverity constant value for critical severity
	ResourceFieldThresholdRuleCriticalSeverity = "critical"

	// Time threshold fields
	// InfraAlertConfigFieldTimeThresholdViolationsInSequence constant value for violations_in_sequence
	InfraAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	// InfraAlertConfigFieldTimeThresholdTimeWindow constant value for time_window
	InfraAlertConfigFieldTimeThresholdTimeWindow = "time_window"

	// Rules fields
	// InfraAlertConfigFieldGenericRule constant value for generic_rule
	InfraAlertConfigFieldGenericRule = "generic_rule"
	// InfraAlertConfigFieldMetricName constant value for metric_name
	InfraAlertConfigFieldMetricName = "metric_name"
	// InfraAlertConfigFieldEntityType constant value for entity_type
	InfraAlertConfigFieldEntityType = "entity_type"
	// InfraAlertConfigFieldAggregation constant value for aggregation
	InfraAlertConfigFieldAggregation = "aggregation"
	// InfraAlertConfigFieldCrossSeriesAggregation constant value for cross_series_aggregation
	InfraAlertConfigFieldCrossSeriesAggregation = "cross_series_aggregation"
	// InfraAlertConfigFieldRegex constant value for regex
	InfraAlertConfigFieldRegex = "regex"
	// InfraAlertConfigFieldThresholdOperator constant value for threshold_operator
	InfraAlertConfigFieldThresholdOperator = "threshold_operator"
	// InfraAlertConfigFieldThreshold constant value for threshold
	InfraAlertConfigFieldThreshold = "threshold"
)

// Resource description constants
const (
	// InfraAlertConfigDescResource description for the infrastructure alert config resource
	InfraAlertConfigDescResource = "This resource represents an infrastructure alert configuration in Instana"
	// InfraAlertConfigDescID description for the ID field
	InfraAlertConfigDescID = "The ID of the infrastructure alert configuration"
	// InfraAlertConfigDescName description for the name field
	InfraAlertConfigDescName = "The name of the infrastructure alert configuration"
	// InfraAlertConfigDescDescription description for the description field
	InfraAlertConfigDescDescription = "The description of the infrastructure alert configuration"
	// InfraAlertConfigDescTagFilter description for the tag_filter field
	InfraAlertConfigDescTagFilter = "The tag filter expression for the infrastructure alert configuration"
	// InfraAlertConfigDescGroupBy description for the group_by field
	InfraAlertConfigDescGroupBy = "The list of tags to group by"
	// InfraAlertConfigDescGranularity description for the granularity field
	InfraAlertConfigDescGranularity = "The granularity of the infrastructure alert configuration"
	// InfraAlertConfigDescEvaluationType description for the evaluation_type field
	InfraAlertConfigDescEvaluationType = "The evaluation type of the infrastructure alert configuration"
	// InfraAlertConfigDescRules description for the rules field
	InfraAlertConfigDescRules = "The rules configuration"
	// InfraAlertConfigDescGenericRule description for the generic_rule field
	InfraAlertConfigDescGenericRule = "The generic rule configuration"
	// InfraAlertConfigDescMetricName description for the metric_name field
	InfraAlertConfigDescMetricName = "The metric name for the generic rule"
	// InfraAlertConfigDescEntityType description for the entity_type field
	InfraAlertConfigDescEntityType = "The entity type for the generic rule"
	// InfraAlertConfigDescAggregation description for the aggregation field
	InfraAlertConfigDescAggregation = "The aggregation for the generic rule"
	// InfraAlertConfigDescCrossSeriesAggregation description for the cross_series_aggregation field
	InfraAlertConfigDescCrossSeriesAggregation = "The cross series aggregation for the generic rule"
	// InfraAlertConfigDescRegex description for the regex field
	InfraAlertConfigDescRegex = "Whether regex is enabled for the generic rule"
	// InfraAlertConfigDescThresholdOperator description for the threshold_operator field
	InfraAlertConfigDescThresholdOperator = "The threshold operator for the generic rule"
	// InfraAlertConfigDescThreshold description for the threshold field
	InfraAlertConfigDescThreshold = "Threshold configuration for different severity levels"
	// InfraAlertConfigDescTimeThreshold description for the time_threshold field
	InfraAlertConfigDescTimeThreshold = "Indicates the type of violation of the defined threshold."
	// InfraAlertConfigDescViolationsInSequence description for the violations_in_sequence field
	InfraAlertConfigDescViolationsInSequence = "Time threshold base on violations in sequence"
	// InfraAlertConfigDescTimeWindow description for the time_window field
	InfraAlertConfigDescTimeWindow = "The time window if the time threshold"
	// InfraAlertConfigDescAlertChannels description for the alert_channels field
	InfraAlertConfigDescAlertChannels = "Set of alert channel IDs associated with the severity."
	// InfraAlertConfigDescAlertChannelIDs description for alert channel ID lists
	InfraAlertConfigDescAlertChannelIDs = "List of IDs of alert channels defined in Instana."
)

// Error message constants
const (
	// InfraAlertConfigErrMappingTagFilter error title for tag filter mapping failures
	InfraAlertConfigErrMappingTagFilter = "Error mapping tag filter"
	// InfraAlertConfigErrMappingTagFilterMsg error message template for tag filter mapping failures
	InfraAlertConfigErrMappingTagFilterMsg = "Failed to map tag filter: %s"
	// InfraAlertConfigErrParsingTagFilter error title for tag filter parsing failures
	InfraAlertConfigErrParsingTagFilter = "Error parsing tag filter"
	// InfraAlertConfigErrParsingTagFilterMsg error message template for tag filter parsing failures
	InfraAlertConfigErrParsingTagFilterMsg = "Failed to parse tag filter: %s"
)

// API-related constants
const (
	// TimeThresholdTypeViolationsInSequence represents the violations in sequence time threshold type
	TimeThresholdTypeViolationsInSequence = "violationsInSequence"
	// GenericRuleAlertType represents the generic rule alert type
	GenericRuleAlertType = "genericRule"
	// EmptyString represents an empty string value
	EmptyString = ""
)

// Made with Bob
