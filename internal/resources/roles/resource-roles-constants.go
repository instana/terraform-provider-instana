package roles

// ResourceInstanaRole the name of the terraform-provider-instana resource to manage roles for role based access control
const ResourceInstanaRole = "rbac_role"

//nolint:gosec
const (
	// Schema field names

	// RoleFieldID constant value for the schema field id
	RoleFieldID = "id"
	// RoleFieldName constant value for the schema field name
	RoleFieldName = "name"
	// RoleFieldMembers constant value for the schema field members
	RoleFieldMembers = "member"
	// RoleFieldMemberUserID constant value for the schema field user_id
	RoleFieldMemberUserID = "user_id"
	// RoleFieldPermissions constant value for the schema field permissions
	RoleFieldPermissions = "permissions"

	// Resource description constants

	// RoleDescResource description for the role resource
	RoleDescResource = "This resource manages RBAC roles in Instana."
	// RoleDescID description for the ID field
	RoleDescID = "The ID of the role."
	// RoleDescName description for the name field
	RoleDescName = "The name of the role"
	// RoleDescMembers description for the members field
	RoleDescMembers = "The members of the role"
	// RoleDescMemberUserID description for the member user_id field
	RoleDescMemberUserID = "The user id of the role member"
	// RoleDescPermissions description for the permissions field
	RoleDescPermissions = "The permissions assigned to the role"
)
