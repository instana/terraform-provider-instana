package shared

// Field name constants for builtin event specification
const (
	// BuiltinEventSpecificationFieldName constant value for the schema field name
	BuiltinEventSpecificationFieldName = "name"
	// BuiltinEventSpecificationFieldDescription constant value for the schema field description
	BuiltinEventSpecificationFieldDescription = "description"
	// BuiltinEventSpecificationFieldShortPluginID constant value for the schema field short_plugin_id
	BuiltinEventSpecificationFieldShortPluginID = "short_plugin_id"
	// BuiltinEventSpecificationFieldSeverity constant value for the schema field severity
	BuiltinEventSpecificationFieldSeverity = "severity"
	// BuiltinEventSpecificationFieldSeverityCode constant value for the schema field severity_code
	BuiltinEventSpecificationFieldSeverityCode = "severity_code"
	// BuiltinEventSpecificationFieldTriggering constant value for the schema field triggering
	BuiltinEventSpecificationFieldTriggering = "triggering"
	// BuiltinEventSpecificationFieldEnabled constant value for the schema field enabled
	BuiltinEventSpecificationFieldEnabled  = "enabled"
	DataSourceInstanaBuiltinEventFramework = "builtin_event_spec"
)

// Field ID constant
const (
	// BuiltinEventFieldID constant value for the schema field id
	BuiltinEventFieldID = "id"
)

// Description constants for builtin event specification fields
const (
	// BuiltinEventDescDataSource description for the data source
	BuiltinEventDescDataSource = "Data source for an Instana builtin event specification. Builtin events are predefined events in Instana."
	// BuiltinEventDescID description for the ID field
	BuiltinEventDescID = "The ID of the builtin event."
	// BuiltinEventDescName description for the name field
	BuiltinEventDescName = "The name of the builtin event."
	// BuiltinEventDescDescription description for the description field
	BuiltinEventDescDescription = "The description text of the builtin event."
	// BuiltinEventDescShortPluginID description for the short_plugin_id field
	BuiltinEventDescShortPluginID = "The plugin id for which the builtin event is created."
	// BuiltinEventDescSeverity description for the severity field
	BuiltinEventDescSeverity = "The severity (WARNING, CRITICAL, etc.) of the builtin event."
	// BuiltinEventDescSeverityCode description for the severity_code field
	BuiltinEventDescSeverityCode = "The severity code used by Instana API (5, 10, etc.) of the builtin event."
	// BuiltinEventDescTriggering description for the triggering field
	BuiltinEventDescTriggering = "Indicates if an incident is triggered the builtin event or not."
	// BuiltinEventDescEnabled description for the enabled field
	BuiltinEventDescEnabled = "Indicates if the builtin event is enabled or not."
)

// Error message constants
const (
	// BuiltinEventErrUnexpectedConfigureType error message for unexpected configure type
	BuiltinEventErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	// BuiltinEventErrUnexpectedConfigureTypeDetail error message detail for unexpected configure type
	BuiltinEventErrUnexpectedConfigureTypeDetail = "Expected *ProviderMeta, got: %T. Please report this issue to the provider developers."
	// BuiltinEventErrReadingEvents error message for reading builtin events
	BuiltinEventErrReadingEvents = "Error reading builtin events"
	// BuiltinEventErrReadingEventsDetail error message detail for reading builtin events
	BuiltinEventErrReadingEventsDetail = "Could not read builtin events: %s"
	// BuiltinEventErrNotFound error message for builtin event not found
	BuiltinEventErrNotFound = "Builtin event not found"
	// BuiltinEventErrNotFoundDetail error message detail for builtin event not found
	BuiltinEventErrNotFoundDetail = "No built in event found for name '%s' and short plugin ID '%s'"
	// BuiltinEventErrConvertingSeverity error message for converting severity
	BuiltinEventErrConvertingSeverity = "Error converting severity"
	// BuiltinEventErrConvertingSeverityDetail error message detail for converting severity
	BuiltinEventErrConvertingSeverityDetail = "Could not convert severity: %s"
)

// Made with Bob
