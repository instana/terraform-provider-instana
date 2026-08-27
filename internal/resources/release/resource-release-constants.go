package release

// ResourceInstanaRelease the name of the terraform-provider-instana resource to manage releases
const ResourceInstanaRelease = "release"

const (
	// Field name constants

	// ReleaseFieldID constant value for the schema field id
	ReleaseFieldID = "id"
	// ReleaseFieldName constant value for the schema field name
	ReleaseFieldName = "name"
	// ReleaseFieldStart constant value for the schema field start
	ReleaseFieldStart = "start"
	// ReleaseFieldLastUpdated constant value for the schema field last_updated
	ReleaseFieldLastUpdated = "last_updated"
	// ReleaseFieldApplications constant value for the schema field applications
	ReleaseFieldApplications = "applications"
	// ReleaseFieldServices constant value for the schema field services
	ReleaseFieldServices = "services"
	// ReleaseFieldScopedTo constant value for the schema field scoped_to
	ReleaseFieldScopedTo = "scoped_to"
	// ReleaseFieldApplicationName constant value for the schema field application_name
	ReleaseFieldApplicationName = "application_name"
	// ReleaseFieldEnvironmentName constant value for the schema field environment_name
	ReleaseFieldEnvironmentName = "environment_name"

	// Description constants

	// ReleaseDescResource description for the resource
	ReleaseDescResource = "This resource manages release markers in Instana."
	// ReleaseDescID description for the ID field
	ReleaseDescID = "The unique ID of the release."
	// ReleaseDescName description for the name field
	ReleaseDescName = "The name of the release. For example: `frontend/release-2000`."
	// ReleaseDescStart description for the start field
	ReleaseDescStart = "The timestamp (in milliseconds since epoch) for when the release is created. For example: `1742349976000`."
	// ReleaseDescLastUpdated description for the last_updated field
	ReleaseDescLastUpdated = "The timestamp (in milliseconds since epoch) of the last update to this release."
	// ReleaseDescApplications description for the applications field
	ReleaseDescApplications = "The list of application perspectives where the release can be viewed (max 10)."
	// ReleaseDescApplicationName description for the application name field
	ReleaseDescApplicationName = "Name of the Application Perspective. For example: `app1`."
	// ReleaseDescServices description for the services field
	ReleaseDescServices = "The list of services where the release can be viewed (max 10)."
	// ReleaseDescServiceName description for the service name field
	ReleaseDescServiceName = "Name of the Service. For example: `payment`."
	// ReleaseDescScopedTo description for the scoped_to block
	ReleaseDescScopedTo = "Optional scope restriction for the service."
	// ReleaseDescApplicationNameScope description for the application_name scope field
	ReleaseDescApplicationNameScope = "The application name that scopes the service."
	// ReleaseDescEnvironmentName description for the environment_name scope field
	ReleaseDescEnvironmentName = "The environment name that scopes the service."

	// Validation constants

	// ReleaseNameMaxLength maximum length for the name field
	ReleaseNameMaxLength = 256
	// ReleaseApplicationsMaxItems maximum number of application scopes
	ReleaseApplicationsMaxItems = 10
	// ReleaseServicesMaxItems maximum number of service scopes
	ReleaseServicesMaxItems = 10
	// ReleaseStartMinValue minimum value for the start timestamp
	ReleaseStartMinValue = 1
)
