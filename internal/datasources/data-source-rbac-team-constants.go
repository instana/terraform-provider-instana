package datasources

const (
	RbacTeamDataSourceFieldID  = "id"
	RbacTeamDataSourceFieldTag = "tag"

	RbacTeamDescDataSource = "Data source for an Instana RBAC team. RBAC teams group users and define their access scope."
	RbacTeamDescID         = "ID of the RBAC Team."
	RbacTeamDescTag        = "Tag/name of the RBAC Team."

	RbacTeamErrMissingLookupAttribute = "Exactly one of 'id' or 'tag' must be set."
	RbacTeamErrReadByID               = "Could not read RBAC Team with ID '%s': %s"
	RbacTeamErrReadAll                = "Could not read RBAC Teams: %s"
	RbacTeamErrNotFoundByID           = "No RBAC Team found with ID '%s'"
	RbacTeamErrNotFoundByTag          = "No RBAC Team found with tag '%s'"
)
