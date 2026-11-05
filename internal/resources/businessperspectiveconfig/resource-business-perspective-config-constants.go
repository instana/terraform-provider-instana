package businessperspectiveconfig

// ResourceInstanaBusinessPerspectiveConfig the name of the terraform-provider-instana resource to manage business perspective configs
const ResourceInstanaBusinessPerspectiveConfig = "business_perspective_config"

// Description constants for Business Perspective Config resource
const (
	BusinessPerspectiveConfigDescResource           = "This resource manages business perspective configurations in Instana."
	BusinessPerspectiveConfigDescID                 = "The ID of the business perspective configuration."
	BusinessPerspectiveConfigDescName               = "The name of the business perspective configuration."
	BusinessPerspectiveConfigDescDescription        = "The description of the business perspective configuration."
	BusinessPerspectiveConfigDescTagFilter          = "The tag filter expression for the business perspective configuration."
	BusinessPerspectiveConfigDescRbacTags           = "RBAC tags (teams) the business perspective is assigned to."
	BusinessPerspectiveConfigDescRbacTagID          = "The ID of the RBAC team."
	BusinessPerspectiveConfigDescRbacTagDisplayName = "The display name of the RBAC team."
)

// Error message constants
const (
	BusinessPerspectiveConfigErrConvertingTagFilter  = "Error converting tag filter"
	BusinessPerspectiveConfigErrParsingTagFilter     = "Error parsing tag filter"
	BusinessPerspectiveConfigErrInvalidAttributeType = "Invalid attribute type"
	BusinessPerspectiveConfigErrMissingAttribute     = "Missing attribute"
	BusinessPerspectiveConfigErrFailedToConvert      = "Failed to convert tag filter: %s"
	BusinessPerspectiveConfigErrFailedToParse        = "Failed to parse tag filter: %s"
	BusinessPerspectiveConfigErrMappingToAPI         = "Error mapping business perspective configuration to API"
	BusinessPerspectiveConfigErrExtractingModel      = "Error extracting business perspective configuration model"
	BusinessPerspectiveConfigErrUpdatingState        = "Error updating business perspective configuration state"
)

// Field name constants
const (
	// BusinessPerspectiveConfigFieldID field name for id
	BusinessPerspectiveConfigFieldID = "id"
	// BusinessPerspectiveConfigFieldName field name for name
	BusinessPerspectiveConfigFieldName = "name"
	// BusinessPerspectiveConfigFieldDescription field name for description
	BusinessPerspectiveConfigFieldDescription = "description"
	// BusinessPerspectiveConfigFieldTagFilter field name for tag_filter
	BusinessPerspectiveConfigFieldTagFilter = "tag_filter"
	// BusinessPerspectiveConfigFieldRbacTags field name for rbac_tags
	BusinessPerspectiveConfigFieldRbacTags = "rbac_tags"
	// BusinessPerspectiveConfigFieldRbacTagID field name for rbac_tags.id
	BusinessPerspectiveConfigFieldRbacTagID = "id"
	// BusinessPerspectiveConfigFieldRbacTagDisplayName field name for rbac_tags.display_name
	BusinessPerspectiveConfigFieldRbacTagDisplayName = "display_name"
)

// Schema version
const BusinessPerspectiveConfigSchemaVersion = 0
