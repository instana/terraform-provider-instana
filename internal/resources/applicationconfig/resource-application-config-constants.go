package applicationconfig

// ResourceInstanaApplicationConfigFramework the name of the terraform-provider-instana resource to manage application config
const ResourceInstanaApplicationConfigFramework = "application_config"

// Description constants for Application Config resource
const (
	ApplicationConfigDescResource      = "This resource manages application configurations in Instana."
	ApplicationConfigDescID            = "The ID of the application configuration."
	ApplicationConfigDescLabel         = "The label of the application config"
	ApplicationConfigDescScope         = "The scope of the application config"
	ApplicationConfigDescBoundaryScope = "The boundary scope of the application config"
	ApplicationConfigDescTagFilter     = "The tag filter expression"
	ApplicationConfigDescAccessRules   = "The access rules applied to the application config"
	ApplicationConfigDescAccessType    = "The access type of the given access rule"
	ApplicationConfigDescRelatedID     = "The id of the related entity (user, api_token, etc.) of the given access rule"
	ApplicationConfigDescRelationType  = "The relation type of the given access rule"
)

// Error message constants
const (
	ApplicationConfigErrConvertingTagFilter  = "Error converting tag filter"
	ApplicationConfigErrParsingTagFilter     = "Error parsing tag filter"
	ApplicationConfigErrInvalidAttributeType = "Invalid attribute type"
	ApplicationConfigErrMissingAttribute     = "Missing attribute"
	ApplicationConfigErrFailedToConvert      = "Failed to convert tag filter: %s"
	ApplicationConfigErrFailedToParse        = "Failed to parse tag filter: %s"
)

// Field name constants
const (
	//ApplicationConfigFieldID field name for id
	ApplicationConfigFieldID = "id"
	//ApplicationConfigFieldAccessRules field name for access rules
	ApplicationConfigFieldAccessRules = "access_rules"
	//ApplicationConfigFieldAccessType field name for access type
	ApplicationConfigFieldAccessType = "access_type"
	//ApplicationConfigFieldRelatedID field name for related ID
	ApplicationConfigFieldRelatedID = "related_id"
	//ApplicationConfigFieldRelationType field name for relation type
	ApplicationConfigFieldRelationType = "relation_type"
	//ApplicationConfigFieldLabel const for the label field of the application config
	ApplicationConfigFieldLabel = "label"
	//ApplicationConfigFieldScope const for the scope field of the application config
	ApplicationConfigFieldScope = "scope"
	//ApplicationConfigFieldBoundaryScope const for the boundary_scope field of the application config
	ApplicationConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationConfigFieldTagFilter const for the tag_filter field of the application config
	ApplicationConfigFieldTagFilter = "tag_filter"
)
