package applicationalertconfig

// ResourceInstanaApplicationAlertConfigFramework the name of the terraform-provider-instana resource to manage application alert configs
const ResourceInstanaApplicationAlertConfigFramework = "application_alert_config"

// ResourceInstanaGlobalApplicationAlertConfigFramework the name of the terraform-provider-instana resource to manage global application alert configs
const ResourceInstanaGlobalApplicationAlertConfigFramework = "global_application_alert_config"

const (
	//ApplicationAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_application_alert_config
	ApplicationAlertConfigFieldAlertChannelIDs = "alert_channel_ids"
	//ApplicationAlertConfigFieldBoundaryScope constant value for field boundary_scope of resource instana_application_alert_config
	ApplicationAlertConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationAlertConfigFieldDescription constant value for field description of resource instana_application_alert_config
	ApplicationAlertConfigFieldDescription = "description"
	//ApplicationAlertConfigFieldEvaluationType constant value for field evaluation_type of resource instana_application_alert_config
	ApplicationAlertConfigFieldEvaluationType = "evaluation_type"
	//ApplicationAlertConfigFieldGranularity constant value for field granularity of resource instana_application_alert_config
	ApplicationAlertConfigFieldGranularity = "granularity"
	//ApplicationAlertConfigFieldIncludeInternal constant value for field include_internal of resource instana_application_alert_config
	ApplicationAlertConfigFieldIncludeInternal = "include_internal"
	//ApplicationAlertConfigFieldIncludeSynthetic constant value for field include_synthetic of resource instana_application_alert_config
	ApplicationAlertConfigFieldIncludeSynthetic = "include_synthetic"
	//ApplicationAlertConfigFieldName constant value for field name of resource instana_application_alert_config
	ApplicationAlertConfigFieldName = "name"
	//ApplicationAlertConfigFieldFullName constant value for field full_name of resource instana_application_alert_config
	ApplicationAlertConfigFieldFullName = "full_name"
	//ApplicationAlertConfigFieldRule constant value for field rule of resource instana_application_alert_config
	ApplicationAlertConfigFieldRule = "rule"
	//ApplicationAlertConfigFieldRuleMetricName constant value for field rule.*.metric_name of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleMetricName = "metric_name"
	//ApplicationAlertConfigFieldRuleAggregation constant value for field rule.*.aggregation of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleAggregation = "aggregation"
	//ApplicationAlertConfigFieldRuleErrorRate constant value for field rule.error_rate of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleErrorRate = "error_rate"
	//ApplicationAlertConfigFieldRuleErrors constant value for field rule.errors of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleErrors = "errors"
	//ApplicationAlertConfigFieldRuleLogs constant value for field rule.logs of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogs = "logs"
	//ApplicationAlertConfigFieldRuleLogsLevel constant value for field rule.logs.level of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsLevel = "level"
	//ApplicationAlertConfigFieldRuleLogsMessage constant value for field rule.logs.message of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsMessage = "message"
	//ApplicationAlertConfigFieldRuleLogsOperator constant value for field rule.logs.operator of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsOperator = "operator"
	//ApplicationAlertConfigFieldRuleSlowness constant value for field rule.slowness of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleSlowness = "slowness"
	//ApplicationAlertConfigFieldRuleStatusCode constant value for field rule.status_code of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCode = "status_code"
	//ApplicationAlertConfigFieldRuleStatusCodeStart constant value for field rule.status_code.status_code_start of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCodeStart = "status_code_start"
	//ApplicationAlertConfigFieldRuleStatusCodeEnd constant value for field rule.status_code.status_code_end of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCodeEnd = "status_code_end"
	//ApplicationAlertConfigFieldRuleThroughput constant value for field rule.throughput of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleThroughput = "throughput"
	//ApplicationAlertConfigFieldSeverity constant value for field severity of resource instana_application_alert_config
	ApplicationAlertConfigFieldSeverity = "severity"
	//ApplicationAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_application_alert_config
	ApplicationAlertConfigFieldTagFilter = "tag_filter"
	//ApplicationAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThreshold = "time_threshold"
	//ApplicationAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	//ApplicationAlertConfigFieldTimeThresholdRequestImpact constant value for field time_threshold.request_impact of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdRequestImpact = "request_impact"
	//ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests constant value for field time_threshold.request_impact.requests of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests = "requests"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	//ApplicationAlertConfigFieldTriggering constant value for field triggering of resource instana_application_alert_config
	ApplicationAlertConfigFieldTriggering = "triggering"
)

// Additional application alert config field names
const (
	ApplicationAlertConfigFieldRules             = "rules"
	ApplicationAlertConfigFieldThreshold         = "threshold"
	ApplicationAlertConfigFieldThresholdOperator = "threshold_operator"
	ApplicationAlertConfigFieldGracePeriod       = "grace_period"
	ApplicationAlertConfigFieldAlertChannels     = "alert_channels"
	ApplicationAlertConfigFieldRuleConfig        = "rule_config"
	ApplicationAlertConfigFieldValue             = "value"

	// Re-define constants from resource-application-alert-config.go for compatibility
	ApplicationAlertConfigFieldApplications                            = "application"
	ApplicationAlertConfigFieldApplicationsApplicationID               = "application_id"
	ApplicationAlertConfigFieldApplicationsInclusive                   = "inclusive"
	ApplicationAlertConfigFieldApplicationsServices                    = "service"
	ApplicationAlertConfigFieldApplicationsServicesServiceID           = "service_id"
	ApplicationAlertConfigFieldApplicationsServicesEndpoints           = "endpoint"
	ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID = "endpoint_id"
)
