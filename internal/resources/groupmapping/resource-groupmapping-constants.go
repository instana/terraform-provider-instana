package groupmapping

// ResourceInstanaGroupMapping the name of the terraform-provider-instana resource to manage group mappings
const ResourceInstanaGroupMapping = "rbac_group_mapping"

const (
	// Schema field names

	// GroupMappingFieldID constant value for the schema field id
	GroupMappingFieldID = "id"
	// GroupMappingFieldKey constant value for the schema field key
	GroupMappingFieldKey = "key"
	// GroupMappingFieldValue constant value for the schema field value
	GroupMappingFieldValue = "value"
	// GroupMappingFieldGroupID constant value for the schema field group_id
	GroupMappingFieldGroupID = "group_id"
	// GroupMappingFieldTeamID constant value for the schema field team_id
	GroupMappingFieldTeamID = "team_id"

	// Resource description constants

	// GroupMappingDescResource description for the group mapping resource
	GroupMappingDescResource = "This resource manages RBAC group mappings in Instana. A group mapping maps an IdP (LDAP, OIDC, SAML) attribute key/value pair to an Instana group so that users whose IdP token contains that pair are automatically assigned to the corresponding group on login."
	// GroupMappingDescID description for the ID field
	GroupMappingDescID = "The ID of the group mapping."
	// GroupMappingDescKey description for the key field
	GroupMappingDescKey = "The IdP attribute key used for the group mapping."
	// GroupMappingDescValue description for the value field
	GroupMappingDescValue = "The IdP attribute value used for the group mapping."
	// GroupMappingDescGroupID description for the group_id field
	GroupMappingDescGroupID = "The ID of the Instana group that matched users will be assigned to."
	// GroupMappingDescTeamID description for the team_id field
	GroupMappingDescTeamID = "The ID of the Instana team to additionally scope this mapping to. Optional."
)
