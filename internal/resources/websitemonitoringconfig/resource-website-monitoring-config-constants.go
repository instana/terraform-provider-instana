package websitemonitoringconfig

// ResourceInstanaWebsiteMonitoringConfig the name of the terraform-provider-instana resource to manage website monitoring configurations
const ResourceInstanaWebsiteMonitoringConfig = "website_monitoring_config"

// Schema field name constants
const (
	// WebsiteMonitoringConfigFieldID constant value for the schema field id
	WebsiteMonitoringConfigFieldID = "id"
	// WebsiteMonitoringConfigFieldName constant value for the schema field name
	WebsiteMonitoringConfigFieldName = "name"
	// WebsiteMonitoringConfigFieldAppName constant value for the schema field app_name
	WebsiteMonitoringConfigFieldAppName = "app_name"
)

// Schema configuration constants
const (
	// WebsiteMonitoringConfigSchemaVersion defines the current schema version
	WebsiteMonitoringConfigSchemaVersion = 1
	// WebsiteMonitoringConfigComputedFieldsCount defines the number of computed fields
	WebsiteMonitoringConfigComputedFieldsCount = 2
	// WebsiteMonitoringConfigRequiredFieldsCount defines the number of required fields
	WebsiteMonitoringConfigRequiredFieldsCount = 1
)

// Resource description constants
const (
	// WebsiteMonitoringConfigDescResource describes the resource purpose
	WebsiteMonitoringConfigDescResource = "This resource manages Website Monitoring Configurations in Instana."
	// WebsiteMonitoringConfigDescID describes the ID field
	WebsiteMonitoringConfigDescID = "The ID of the Website Monitoring Configuration."
	// WebsiteMonitoringConfigDescName describes the name field
	WebsiteMonitoringConfigDescName = "Configures the name of the website monitoring configuration."
	// WebsiteMonitoringConfigDescAppName describes the app_name field
	WebsiteMonitoringConfigDescAppName = "Configures the calculated app name of the website monitoring configuration."
)

// Error message constants
const (
	// WebsiteMonitoringConfigErrRetrievingPlan error message for plan retrieval failures
	WebsiteMonitoringConfigErrRetrievingPlan = "Error retrieving plan data"
	// WebsiteMonitoringConfigErrRetrievingState error message for state retrieval failures
	WebsiteMonitoringConfigErrRetrievingState = "Error retrieving state data"
	// WebsiteMonitoringConfigErrMappingToAPI error message for API object mapping failures
	WebsiteMonitoringConfigErrMappingToAPI = "Error mapping to API object"
	// WebsiteMonitoringConfigErrUpdatingState error message for state update failures
	WebsiteMonitoringConfigErrUpdatingState = "Error updating state"
	// WebsiteMonitoringConfigErrInvalidInput error message for invalid input data
	WebsiteMonitoringConfigErrInvalidInput = "Invalid input data provided"
	// WebsiteMonitoringConfigErrNilContext error message for nil context
	WebsiteMonitoringConfigErrNilContext = "Context cannot be nil"
	// WebsiteMonitoringConfigErrNilState error message for nil state
	WebsiteMonitoringConfigErrNilState = "State cannot be nil"
	// WebsiteMonitoringConfigErrNilAPIObject error message for nil API object
	WebsiteMonitoringConfigErrNilAPIObject = "API object cannot be nil"
)

// Validation constants
const (
	// WebsiteMonitoringConfigMinNameLength minimum length for name field
	WebsiteMonitoringConfigMinNameLength = 1
	// WebsiteMonitoringConfigMaxNameLength maximum length for name field
	WebsiteMonitoringConfigMaxNameLength = 256
	// WebsiteMonitoringConfigMinIDLength minimum length for ID field
	WebsiteMonitoringConfigMinIDLength = 1
	// WebsiteMonitoringConfigMaxIDLength maximum length for ID field
	WebsiteMonitoringConfigMaxIDLength = 64
)

// Default values constants
const (
	// WebsiteMonitoringConfigDefaultSchemaVersion default schema version
	WebsiteMonitoringConfigDefaultSchemaVersion = 1
)
