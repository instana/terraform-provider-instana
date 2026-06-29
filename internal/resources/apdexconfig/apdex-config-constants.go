package apdexconfig

// ResourceInstanaApdexConfig the name of the terraform-provider-instana resource to manage Apdex configurations
const ResourceInstanaApdexConfig = "apdex_config"

// ApdexConfigFieldRbacTags is the field name for RBAC tags
const ApdexConfigFieldRbacTags = "rbac_tags"

// Resource description constants
const (
	// ApdexConfigDescResource is the description for the Apdex config resource
	ApdexConfigDescResource = "This resource manages Apdex V2 Configurations in Instana."
	// ApdexConfigDescID is the description for the ID field
	ApdexConfigDescID = "The ID of the Apdex configuration"
	// ApdexConfigDescApdexName is the description for the apdex_name field
	ApdexConfigDescApdexName = "The name of the Apdex configuration"
	// ApdexConfigDescTags is the description for the tags field
	ApdexConfigDescTags = "The tags of the Apdex configuration"
	// ApdexConfigDescRbacTags is the description for the rbac_tags field
	ApdexConfigDescRbacTags = "RBAC tags for the Apdex configuration"
	// ApdexConfigDescRbacTagDisplayName is the description for the display_name field in rbac_tags
	ApdexConfigDescRbacTagDisplayName = "Display name of the RBAC tag"
	// ApdexConfigDescRbacTagID is the description for the id field in rbac_tags
	ApdexConfigDescRbacTagID = "ID of the RBAC tag"
	// ApdexConfigDescApdexEntity is the description for the apdex_entity block
	ApdexConfigDescApdexEntity = "The entity configuration for the Apdex"
	// ApdexConfigDescApplicationEntity is the description for the application entity block
	ApdexConfigDescApplicationEntity = "Application entity of Apdex"
	// ApdexConfigDescEntityID is the description for the entity_id field
	ApdexConfigDescEntityID = "The entity ID (Application ID or Website ID)"
	// ApdexConfigDescThreshold is the description for the threshold field
	ApdexConfigDescThreshold = "The Apdex threshold value in milliseconds"
	// ApdexConfigDescFilterExpression is the description for the filter_expression field
	ApdexConfigDescFilterExpression = "Tag filter expression for the entity"
	// ApdexConfigDescBoundaryScope is the description for the boundary_scope field
	ApdexConfigDescBoundaryScope = "The boundary scope for the application entity (ALL, INBOUND)"
	// ApdexConfigDescIncludeInternal is the description for the include_internal field
	ApdexConfigDescIncludeInternal = "Flag to indicate whether internal calls are included"
	// ApdexConfigDescIncludeSynthetic is the description for the include_synthetic field
	ApdexConfigDescIncludeSynthetic = "Flag to indicate whether synthetic calls are included"
	// ApdexConfigDescWebsiteEntity is the description for the website entity block
	ApdexConfigDescWebsiteEntity = "Website entity of Apdex"
	// ApdexConfigDescBeaconType is the description for the beacon_type field
	ApdexConfigDescBeaconType = "The beacon type for the website entity (pageLoad, httpRequest, custom)"

	// Error message constants

	// ApdexConfigErrMappingState is the error title for mapping state to data object
	ApdexConfigErrMappingState = "Error mapping state to data object"
	// ApdexConfigErrBothPlanStateNil is the error message when both plan and state are nil
	ApdexConfigErrBothPlanStateNil = "Both plan and state are nil"
	// ApdexConfigErrApplicationEntityRequired is the error message for missing application entity fields
	ApdexConfigErrApplicationEntityRequired = "entity_id, threshold, and boundary_scope are required for application entity"
	// ApdexConfigErrWebsiteEntityRequired is the error message for missing website entity fields
	ApdexConfigErrWebsiteEntityRequired = "entity_id, threshold, and beacon_type are required for website entity"
	// ApdexConfigErrMissingEntity is the error title for missing entity configuration
	ApdexConfigErrMissingEntity = "Missing entity configuration"
	// ApdexConfigErrExactlyOneEntity is the error message for missing entity configuration
	ApdexConfigErrExactlyOneEntity = "Exactly one entity configuration (application or website) is required"
	// ApdexConfigErrParsingFilterExpression is the error title for parsing filter expression
	ApdexConfigErrParsingFilterExpression = "Error parsing filter expression"
	// ApdexConfigErrParsingFilterExpressionMsg is the error message for parsing filter expression
	ApdexConfigErrParsingFilterExpressionMsg = "Could not parse filter expression: %s"
	// ApdexConfigErrMappingEntityToState is the error title for mapping entity to state
	ApdexConfigErrMappingEntityToState = "Error mapping entity to state"
	// ApdexConfigErrUnsupportedEntityType is the error message for unsupported entity type
	ApdexConfigErrUnsupportedEntityType = "Unsupported entity type: %s"
	// ApdexConfigErrNormalizingFilterExpression is the error title for normalizing filter expression
	ApdexConfigErrNormalizingFilterExpression = "Error normalizing filter expression"
	// ApdexConfigErrNormalizingFilterExpressionMsg is the error message for normalizing filter expression
	ApdexConfigErrNormalizingFilterExpressionMsg = "Could not normalize filter expression: %s"

	// ApdexConfigField names for terraform
	ApdexConfigFieldID               = "id"
	ApdexConfigFieldApdexName        = "apdex_name"
	ApdexConfigFieldTags             = "tags"
	ApdexConfigFieldApdexEntity      = "apdex_entity"
	ApdexConfigFieldApdexType        = "apdex_type"
	ApdexConfigFieldEntityID         = "entity_id"
	ApdexConfigFieldThreshold        = "threshold"
	ApdexConfigFieldFilterExpression = "filter_expression"
	ApdexConfigFieldBoundaryScope    = "boundary_scope"
	ApdexConfigFieldIncludeInternal  = "include_internal"
	ApdexConfigFieldIncludeSynthetic = "include_synthetic"
	ApdexConfigFieldBeaconType       = "beacon_type"

	// Apdex entity types for terraform
	ApdexConfigApplicationEntity = "application"
	ApdexConfigWebsiteEntity     = "website"

	// ApdexConfigFieldNames and values for API
	ApdexConfigAPIFieldApdexType     = "apdexType"
	ApdexConfigAPIFieldEntityID      = "entityId"
	ApdexConfigAPIFieldThreshold     = "threshold"
	ApdexConfigAPIFieldFilter        = "tagFilterExpression"
	ApdexConfigAPIFieldBoundaryScope = "boundaryScope"
	ApdexConfigAPIFieldIncludeInternal  = "includeInternal"
	ApdexConfigAPIFieldIncludeSynthetic = "includeSynthetic"
	ApdexConfigAPIFieldBeaconType       = "beaconType"

	// Schema field identifier constants

	// SchemaFieldID represents the id field identifier
	SchemaFieldID = "id"
	// SchemaFieldDisplayName represents the display_name field identifier
	SchemaFieldDisplayName = "display_name"
	// SchemaFieldApplication represents the application field identifier
	SchemaFieldApplication = "application"
	// SchemaFieldWebsite represents the website field identifier
	SchemaFieldWebsite = "website"

	// Tag filter constants

	// TagFilterTypeExpression represents the EXPRESSION tag filter type
	TagFilterTypeExpression = "EXPRESSION"
	// LogicalOperatorAnd represents the AND logical operator
	LogicalOperatorAnd = "AND"

	// EmptyString represents an empty string constant
	EmptyString = ""
)
