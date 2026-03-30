package mobileappconfig

// ResourceInstanaMobileAppConfig the name of the terraform-provider-instana resource to manage mobile app configurations
const ResourceInstanaMobileAppConfig = "mobile_app_config"

// Schema field name constants
const (
	// MobileAppConfigFieldID constant value for the schema field id
	MobileAppConfigFieldID = "id"
	// MobileAppConfigFieldName constant value for the schema field name
	MobileAppConfigFieldName = "name"
)

// Schema configuration constants
const (
	// MobileAppConfigSchemaVersion defines the current schema version
	MobileAppConfigSchemaVersion = 1
	// MobileAppConfigComputedFieldsCount defines the number of computed fields
	MobileAppConfigComputedFieldsCount = 1
	// MobileAppConfigRequiredFieldsCount defines the number of required fields
	MobileAppConfigRequiredFieldsCount = 1
)

// Resource description constants
const (
	// MobileAppConfigDescResource describes the resource purpose
	MobileAppConfigDescResource = "This resource manages Mobile App Configurations in Instana."
	// MobileAppConfigDescID describes the ID field
	MobileAppConfigDescID = "The ID of the Mobile App Configuration."
	// MobileAppConfigDescName describes the name field
	MobileAppConfigDescName = "Configures the name of the mobile app configuration."
)

// Error message constants
const (
	// MobileAppConfigErrRetrievingPlan error message for plan retrieval failures
	MobileAppConfigErrRetrievingPlan = "Error retrieving plan data"
	// MobileAppConfigErrRetrievingState error message for state retrieval failures
	MobileAppConfigErrRetrievingState = "Error retrieving state data"
	// MobileAppConfigErrMappingToAPI error message for API object mapping failures
	MobileAppConfigErrMappingToAPI = "Error mapping to API object"
	// MobileAppConfigErrUpdatingState error message for state update failures
	MobileAppConfigErrUpdatingState = "Error updating state"
	// MobileAppConfigErrInvalidInput error message for invalid input data
	MobileAppConfigErrInvalidInput = "Invalid input data provided"
	// MobileAppConfigErrNilContext error message for nil context
	MobileAppConfigErrNilContext = "Context cannot be nil"
	// MobileAppConfigErrNilState error message for nil state
	MobileAppConfigErrNilState = "State cannot be nil"
	// MobileAppConfigErrNilAPIObject error message for nil API object
	MobileAppConfigErrNilAPIObject = "API object cannot be nil"
)

// Validation constants
const (
	// MobileAppConfigMinNameLength minimum length for name field
	MobileAppConfigMinNameLength = 1
	// MobileAppConfigMaxNameLength maximum length for name field
	MobileAppConfigMaxNameLength = 128
	// MobileAppConfigMinIDLength minimum length for ID field
	MobileAppConfigMinIDLength = 1
	// MobileAppConfigMaxIDLength maximum length for ID field
	MobileAppConfigMaxIDLength = 128
)

// Default values constants
const (
	// MobileAppConfigDefaultSchemaVersion default schema version
	MobileAppConfigDefaultSchemaVersion = 1
)

