package datasources

// DataSourceInstanaApplicationConfig the name of the terraform-provider-instana data source to read application config
const DataSourceInstanaApplicationConfig = "application_config"

// Description constants for Application Config data source
const (
	ApplicationConfigDataSourceDesc = "Data source to retrieve details about an existing application configuration (application perspective) in Instana."
	ApplicationConfigDataSourceDescID = "The unique identifier of the application configuration. Exactly one of 'id' or 'name' must be specified."
	ApplicationConfigDataSourceDescName = "The name (label) of the application configuration. If multiple configurations have the same name, the first match will be returned. Exactly one of 'id' or 'name' must be specified."
	ApplicationConfigDataSourceDescLabel = "The label of the application configuration."
	ApplicationConfigDataSourceDescScope = "The scope of the application configuration."
	ApplicationConfigDataSourceDescBoundaryScope = "The boundary scope of the application configuration."
	ApplicationConfigDataSourceDescTagFilter = "The tag filter expression of the application configuration."
	ApplicationConfigDataSourceDescAccessRules = "The access rules applied to the application configuration."
	ApplicationConfigDataSourceDescAccessType = "The access type of the given access rule."
	ApplicationConfigDataSourceDescRelatedID = "The id of the related entity (user, api_token, etc.) of the given access rule."
	ApplicationConfigDataSourceDescRelationType = "The relation type of the given access rule."
)

// Field name constants
const (
	ApplicationConfigDataSourceFieldID            = "id"
	ApplicationConfigDataSourceFieldName          = "name"
	ApplicationConfigDataSourceFieldLabel         = "label"
	ApplicationConfigDataSourceFieldScope         = "scope"
	ApplicationConfigDataSourceFieldBoundaryScope = "boundary_scope"
	ApplicationConfigDataSourceFieldTagFilter     = "tag_filter"
	ApplicationConfigDataSourceFieldAccessRules   = "access_rules"
	ApplicationConfigDataSourceFieldAccessType    = "access_type"
	ApplicationConfigDataSourceFieldRelatedID     = "related_id"
	ApplicationConfigDataSourceFieldRelationType  = "relation_type"
)

// Error message constants
const (
	ApplicationConfigDataSourceErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	ApplicationConfigDataSourceErrMissingLookupAttribute  = "Exactly one of 'id' or 'name' must be set."
	ApplicationConfigDataSourceErrReadByID                = "Could not read application configuration with ID '%s': %s"
	ApplicationConfigDataSourceErrReadAll                 = "Could not read application configurations: %s"
	ApplicationConfigDataSourceErrNotFoundByID            = "No application configuration found with ID '%s'"
	ApplicationConfigDataSourceErrNotFoundByName          = "No application configuration found with name '%s'"
	ApplicationConfigDataSourceErrConvertingTagFilter     = "Error converting tag filter"
	ApplicationConfigDataSourceErrFailedToConvert         = "Failed to convert tag filter: %s"
)
