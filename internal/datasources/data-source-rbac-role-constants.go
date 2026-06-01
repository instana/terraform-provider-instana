package datasources

const (
	RbacRoleDataSourceFieldID          = "id"
	RbacRoleDataSourceFieldName        = "name"
	RbacRoleDataSourceFieldPermissions = "permissions"

	RbacRoleDescDataSource  = "Data source for an Instana RBAC role. RBAC roles define permissions for users and teams."
	RbacRoleDescID          = "ID of the RBAC Role."
	RbacRoleDescName        = "Name of the RBAC Role."
	RbacRoleDescPermissions = "Permissions assigned to the RBAC Role."

	RbacRoleErrMissingLookupAttribute = "Exactly one of 'id' or 'name' must be set."
	RbacRoleErrReadByID               = "Could not read RBAC Role with ID '%s': %s"
	RbacRoleErrReadAll                = "Could not read RBAC Roles: %s"
	RbacRoleErrNotFoundByID           = "No RBAC Role found with ID '%s'"
	RbacRoleErrNotFoundByName         = "No RBAC Role found with name '%s'"
)
