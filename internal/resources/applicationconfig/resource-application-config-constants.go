package applicationconfig

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

	// Error message constants
	ApplicationConfigErrConvertingTagFilter  = "Error converting tag filter"
	ApplicationConfigErrParsingTagFilter     = "Error parsing tag filter"
	ApplicationConfigErrInvalidAttributeType = "Invalid attribute type"
	ApplicationConfigErrMissingAttribute     = "Missing attribute"
)

// Made with Bob
