package syntheticalertconfig

// Resource name constants
const (
	// ResourceInstanaSyntheticAlertConfig the name of the terraform-provider-instana resource to manage synthetic alert configurations
	ResourceInstanaSyntheticAlertConfig = "synthetic_alert_config"

	// Field name constants for synthetic alert config

	// SyntheticAlertConfigFieldName constant value for the schema field name
	SyntheticAlertConfigFieldName = "name"
	// SyntheticAlertConfigFieldDescription constant value for the schema field description
	SyntheticAlertConfigFieldDescription = "description"
	// SyntheticAlertConfigFieldSyntheticTestIds constant value for the schema field synthetic_test_ids
	SyntheticAlertConfigFieldSyntheticTestIds = "synthetic_test_ids"
	// SyntheticAlertConfigFieldSeverity constant value for the schema field severity
	SyntheticAlertConfigFieldSeverity = "severity"
	// SyntheticAlertConfigFieldTagFilter constant value for the schema field tag_filter
	SyntheticAlertConfigFieldTagFilter = "tag_filter"
	// SyntheticAlertConfigFieldRule constant value for the schema field rule
	SyntheticAlertConfigFieldRule = "rule"
	// SyntheticAlertConfigFieldAlertChannelIds constant value for the schema field alert_channel_ids
	SyntheticAlertConfigFieldAlertChannelIds = "alert_channel_ids"
	// SyntheticAlertConfigFieldTimeThreshold constant value for the schema field time_threshold
	SyntheticAlertConfigFieldTimeThreshold = "time_threshold"
	// SyntheticAlertConfigFieldGracePeriod constant value for the schema field grace_period
	SyntheticAlertConfigFieldGracePeriod = "grace_period"
	// SyntheticAlertConfigFieldID constant value for the schema field id
	SyntheticAlertConfigFieldID = "id"
	// SyntheticAlertConfigFieldCustomPayloadField constant value for the schema field custom_payload_field
	SyntheticAlertConfigFieldCustomPayloadField = "custom_payload_field"

	// Rule field constants

	// SyntheticAlertRuleFieldAlertType constant value for the rule field alert_type
	SyntheticAlertRuleFieldAlertType = "alert_type"
	// SyntheticAlertRuleFieldMetricName constant value for the rule field metric_name
	SyntheticAlertRuleFieldMetricName = "metric_name"
	// SyntheticAlertRuleFieldAggregation constant value for the rule field aggregation
	SyntheticAlertRuleFieldAggregation = "aggregation"

	// TimeThreshold field constants

	// SyntheticAlertTimeThresholdFieldType constant value for the time threshold field type
	SyntheticAlertTimeThresholdFieldType = "type"
	// SyntheticAlertTimeThresholdFieldViolationsCount constant value for the time threshold field violations_count
	SyntheticAlertTimeThresholdFieldViolationsCount = "violations_count"

	// Description constants

	// SyntheticAlertConfigDescResource description for the resource
	SyntheticAlertConfigDescResource = "This resource manages Synthetic Alert Configurations in Instana."
	// SyntheticAlertConfigDescID description for the ID field
	SyntheticAlertConfigDescID = "The ID of the Synthetic Alert Config."
	// SyntheticAlertConfigDescName description for the name field
	SyntheticAlertConfigDescName = "The name of the Synthetic Alert Config."
	// SyntheticAlertConfigDescDescription description for the description field
	SyntheticAlertConfigDescDescription = "The description of the Synthetic Alert Config."
	// SyntheticAlertConfigDescSyntheticTestIds description for the synthetic_test_ids field
	SyntheticAlertConfigDescSyntheticTestIds = "A set of Synthetic Test IDs that this alert config applies to."
	// SyntheticAlertConfigDescSeverity description for the severity field
	SyntheticAlertConfigDescSeverity = "The severity of the alert (5=critical, 10=warning)."
	// SyntheticAlertConfigDescTagFilter description for the tag_filter field
	SyntheticAlertConfigDescTagFilter = "The tag filter expression used for this synthetic alert."
	// SyntheticAlertConfigDescAlertChannelIds description for the alert_channel_ids field
	SyntheticAlertConfigDescAlertChannelIds = "A set of Alert Channel IDs."
	// SyntheticAlertConfigDescGracePeriod description for the grace_period field
	SyntheticAlertConfigDescGracePeriod = "The duration in milliseconds for which an alert remains open after conditions are no longer violated."
	// SyntheticAlertConfigDescRule description for the rule block
	SyntheticAlertConfigDescRule = "Configuration for the synthetic alert rule."
	// SyntheticAlertConfigDescRuleAlertType description for the rule alert_type field
	SyntheticAlertConfigDescRuleAlertType = "The type of the alert rule (e.g., failure)."
	// SyntheticAlertConfigDescRuleMetricName description for the rule metric_name field
	SyntheticAlertConfigDescRuleMetricName = "The metric name to monitor (e.g., status)."
	// SyntheticAlertConfigDescRuleAggregation description for the rule aggregation field
	SyntheticAlertConfigDescRuleAggregation = "The aggregation method {SUM,MEAN,MAX,MIN,P25,P50,P75,P90,P95,P98,P99,P99_9,P99_99,DISTINCT_COUNT,SUM_POSITIVE,PER_SECOND,INCREASE}."
	// SyntheticAlertConfigDescTimeThreshold description for the time_threshold block
	SyntheticAlertConfigDescTimeThreshold = "Configuration for the time threshold."
	// SyntheticAlertConfigDescTimeThresholdType description for the time threshold type field
	SyntheticAlertConfigDescTimeThresholdType = "The type of the time threshold (only violationsInSequence is supported)."
	// SyntheticAlertConfigDescTimeThresholdViolationsCount description for the time threshold violations_count field
	SyntheticAlertConfigDescTimeThresholdViolationsCount = "The number of violations required to trigger the alert (value between 1 and 12)."

	// Error message constants

	// SyntheticAlertConfigErrParsingTagFilter error message for parsing tag filter
	SyntheticAlertConfigErrParsingTagFilter = "Error parsing tag filter"
	// SyntheticAlertConfigErrParsingTagFilterDetail error message detail for parsing tag filter
	SyntheticAlertConfigErrParsingTagFilterDetail = "Could not parse tag filter: "
	// SyntheticAlertConfigErrNormalizingTagFilter error message for normalizing tag filter
	SyntheticAlertConfigErrNormalizingTagFilter = "Error normalizing tag filter"
	// SyntheticAlertConfigErrNormalizingTagFilterDetail error message detail for normalizing tag filter
	SyntheticAlertConfigErrNormalizingTagFilterDetail = "Could not normalize tag filter: "

	// Validation constants

	// SyntheticAlertConfigValidAlertType valid alert type value
	SyntheticAlertConfigValidAlertType = "failure"
	// SyntheticAlertConfigValidTimeThresholdType valid time threshold type value
	SyntheticAlertConfigValidTimeThresholdType = "violationsInSequence"

	// Aggregation type constants

	AggregationTypeSum           = "SUM"
	AggregationTypeMean          = "MEAN"
	AggregationTypeMax           = "MAX"
	AggregationTypeMin           = "MIN"
	AggregationTypeP25           = "P25"
	AggregationTypeP50           = "P50"
	AggregationTypeP75           = "P75"
	AggregationTypeP90           = "P90"
	AggregationTypeP95           = "P95"
	AggregationTypeP98           = "P98"
	AggregationTypeP99           = "P99"
	AggregationTypeP99_9         = "P99_9"
	AggregationTypeP99_99        = "P99_99"
	AggregationTypeDistinctCount = "DISTINCT_COUNT"
	AggregationTypeSumPositive   = "SUM_POSITIVE"
	AggregationTypePerSecond     = "PER_SECOND"
	AggregationTypeIncrease      = "INCREASE"

	// Tag filter constants

	// TagFilterTypeExpression tag filter type value
	TagFilterTypeExpression = "EXPRESSION"
	// TagFilterLogicalOperatorAnd logical operator value
	TagFilterLogicalOperatorAnd = "AND"
)
