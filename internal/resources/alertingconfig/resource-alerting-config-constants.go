package alertingconfig

// Resource name constants
const (
	// ResourceInstanaAlertingConfigFramework the name of the terraform-provider-instana resource to manage alerting configurations
	ResourceInstanaAlertingConfigFramework = "alerting_config"
)

// Field name constants for alerting config
const (
	// AlertingConfigFieldAlertName constant value for the schema field alert_name
	AlertingConfigFieldAlertName = "alert_name"
	// AlertingConfigFieldFullAlertName constant value for the schema field full_alert_name
	AlertingConfigFieldFullAlertName = "full_alert_name"
	// AlertingConfigFieldIntegrationIds constant value for the schema field integration_ids
	AlertingConfigFieldIntegrationIds = "integration_ids"
	// AlertingConfigFieldEventFilterQuery constant value for the schema field event_filter_query
	AlertingConfigFieldEventFilterQuery = "event_filter_query"
	// AlertingConfigFieldEventFilterEventTypes constant value for the schema field event_filter_event_types
	AlertingConfigFieldEventFilterEventTypes = "event_filter_event_types"
	// AlertingConfigFieldEventFilterRuleIDs constant value for the schema field event_filter_rule_ids
	AlertingConfigFieldEventFilterRuleIDs = "event_filter_rule_ids"
	// AlertingConfigFieldID constant value for the schema field id
	AlertingConfigFieldID = "id"
)

// Description constants
const (
	// AlertingConfigDescResource description for the resource
	AlertingConfigDescResource = "This resource manages alerting configurations in Instana."
	// AlertingConfigDescID description for the ID field
	AlertingConfigDescID = "The ID of the alerting configuration."
	// AlertingConfigDescAlertName description for the alert_name field
	AlertingConfigDescAlertName = "Configures the alert name of the alerting configuration"
	// AlertingConfigDescIntegrationIds description for the integration_ids field
	AlertingConfigDescIntegrationIds = "Configures the list of Integration IDs (Alerting Channels)."
	// AlertingConfigDescEventFilterQuery description for the event_filter_query field
	AlertingConfigDescEventFilterQuery = "Configures a filter query to to filter rules or event types for a limited set of entities"
	// AlertingConfigDescEventFilterEventTypes description for the event_filter_event_types field
	AlertingConfigDescEventFilterEventTypes = "Configures the list of Event Types IDs which should trigger an alert."
	// AlertingConfigDescEventFilterRuleIDs description for the event_filter_rule_ids field
	AlertingConfigDescEventFilterRuleIDs = "Configures the list of Rule IDs which should trigger an alert."
)

// Made with Bob
