package maintenancewindowconfig

// Resource name constants
const (
	// ResourceInstanaMaintenanceWindowConfig the name of the terraform-provider-instana resource to manage maintenance window configurations
	ResourceInstanaMaintenanceWindowConfig = "maintenance_window_config"

	// Field name constants for maintenance window config

	// MaintenanceWindowConfigFieldID constant value for the schema field id
	MaintenanceWindowConfigFieldID = "id"
	// MaintenanceWindowConfigFieldName constant value for the schema field name
	MaintenanceWindowConfigFieldName = "name"
	// MaintenanceWindowConfigFieldQuery constant value for the schema field query
	MaintenanceWindowConfigFieldQuery = "query"
	// MaintenanceWindowConfigFieldScheduling constant value for the schema field scheduling
	MaintenanceWindowConfigFieldScheduling = "scheduling"
	// MaintenanceWindowConfigFieldTagFilterExpressionEnabled constant value for the schema field tag_filter_expression_enabled
	MaintenanceWindowConfigFieldTagFilterExpressionEnabled = "tag_filter_expression_enabled"
	// MaintenanceWindowConfigFieldTagFilterExpression constant value for the schema field tag_filter_expression
	MaintenanceWindowConfigFieldTagFilterExpression = "tag_filter_expression"

	// Scheduling field constants
	// SchedulingFieldStart constant value for the schema field start
	SchedulingFieldStart = "start"
	// SchedulingFieldDuration constant value for the schema field duration
	SchedulingFieldDuration = "duration"
	// SchedulingFieldType constant value for the schema field type
	SchedulingFieldType = "type"
	// SchedulingFieldRrule constant value for the schema field rrule
	SchedulingFieldRrule = "rrule"
	// SchedulingFieldTimezoneId constant value for the schema field timezone_id
	SchedulingFieldTimezoneId = "timezone_id"

	// Duration field constants
	// DurationFieldAmount constant value for the schema field amount
	DurationFieldAmount = "amount"
	// DurationFieldUnit constant value for the schema field unit
	DurationFieldUnit = "unit"

	// Description constants

	// MaintenanceWindowConfigDescResource description for the resource
	MaintenanceWindowConfigDescResource = "This resource manages maintenance window configurations in Instana."
	// MaintenanceWindowConfigDescID description for the ID field
	MaintenanceWindowConfigDescID = "The ID of the maintenance window configuration."
	// MaintenanceWindowConfigDescName description for the name field
	MaintenanceWindowConfigDescName = "The name of the maintenance window configuration."
	// MaintenanceWindowConfigDescQuery description for the query field
	MaintenanceWindowConfigDescQuery = "Dynamic Focus Query that determines the scope of the maintenance window configuration."
	// MaintenanceWindowConfigDescScheduling description for the scheduling field
	MaintenanceWindowConfigDescScheduling = "Time scheduling of the maintenance window configuration."
	// MaintenanceWindowConfigDescTagFilterExpressionEnabled description for the tag_filter_expression_enabled field
	MaintenanceWindowConfigDescTagFilterExpressionEnabled = "Boolean flag to determine if the tagFilterExpression is enabled."
	// MaintenanceWindowConfigDescTagFilterExpression description for the tag_filter_expression field
	MaintenanceWindowConfigDescTagFilterExpression = "Tag filter expression used to filter alert notifications that will be muted."

	// SchedulingDescStart description for the start field
	SchedulingDescStart = "Start time in milliseconds from epoch."
	// SchedulingDescDuration description for the duration field
	SchedulingDescDuration = "Duration of each maintenance window occurrence."
	// SchedulingDescType description for the type field
	SchedulingDescType = "Type of maintenance window: ONE_TIME or RECURRENT."
	// SchedulingDescRrule description for the rrule field
	SchedulingDescRrule = "For RECURRENT maintenance configurations, the RRULE standard from the iCalendar Spec."
	// SchedulingDescTimezoneId description for the timezone_id field
	SchedulingDescTimezoneId = "Timezone ID for recurrent maintenance windows."

	// DurationDescAmount description for the amount field
	DurationDescAmount = "The amount of time for the duration."
	// DurationDescUnit description for the unit field
	DurationDescUnit = "The unit of time for the duration: MINUTES, HOURS, or DAYS."
)

// Supported scheduling types
var SupportedSchedulingTypes = []string{"ONE_TIME", "RECURRENT"}

// Supported duration units
var SupportedDurationUnits = []string{"MINUTES", "HOURS", "DAYS"}

// Made with Bob
