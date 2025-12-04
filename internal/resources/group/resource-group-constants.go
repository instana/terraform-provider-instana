package group

// ResourceInstanaGroup the name of the terraform-provider-instana resource to manage groups for role based access control
const ResourceInstanaGroup = "rbac_group"

//nolint:gosec
const (
	// Schema field names

	// GroupFieldID constant value for the schema field id
	GroupFieldID = "id"
	// GroupFieldName constant value for the schema field name
	GroupFieldName = "name"
	// GroupFieldMembers constant value for the schema field members
	GroupFieldMembers = "member"
	// GroupFieldMemberEmail constant value for the schema field email
	GroupFieldMemberEmail = "email"
	// GroupFieldMemberUserID constant value for the schema field user_id
	GroupFieldMemberUserID = "user_id"
	// GroupFieldPermissionSet constant value for the schema field permission_set
	GroupFieldPermissionSet = "permission_set"
	// GroupFieldPermissionSetApplicationIDs constant value for the schema field application_ids
	GroupFieldPermissionSetApplicationIDs = "application_ids"
	// GroupFieldPermissionSetInfraDFQFilter constant value for the schema field infra_dfq_filter
	GroupFieldPermissionSetInfraDFQFilter = "infra_dfq_filter"
	// GroupFieldPermissionSetKubernetesClusterUUIDs constant value for the schema field kubernetes_cluster_uuids
	GroupFieldPermissionSetKubernetesClusterUUIDs = "kubernetes_cluster_uuids"
	// GroupFieldPermissionSetKubernetesNamespaceUIDs constant value for the schema field kubernetes_namespaces_uuids
	GroupFieldPermissionSetKubernetesNamespaceUIDs = "kubernetes_namespaces_uuids"
	// GroupFieldPermissionSetMobileAppIDs constant value for the schema field mobile_app_ids
	GroupFieldPermissionSetMobileAppIDs = "mobile_app_ids"
	// GroupFieldPermissionSetWebsiteIDs constant value for the schema field website_ids
	GroupFieldPermissionSetWebsiteIDs = "website_ids"
	// GroupFieldPermissionSetPermissions constant value for the schema field permissions
	GroupFieldPermissionSetPermissions = "permissions"

	// Resource description constants

	// GroupDescResource description for the group resource
	GroupDescResource = "This resource manages RBAC groups in Instana."
	// GroupDescID description for the ID field
	GroupDescID = "The ID of the group."
	// GroupDescName description for the name field
	GroupDescName = "The name of the Group"
	// GroupDescMembers description for the members field
	GroupDescMembers = "The members of the group"
	// GroupDescMemberUserID description for the member user_id field
	GroupDescMemberUserID = "The user id of the group member"
	// GroupDescMemberEmail description for the member email field
	GroupDescMemberEmail = "The email address of the group member"
	// GroupDescPermissionSet description for the permission_set field
	GroupDescPermissionSet = "The permission set of the group"
	// GroupDescPermissionSetApplicationIDs description for the application_ids field
	GroupDescPermissionSetApplicationIDs = "The scope bindings to restrict access to applications"
	// GroupDescPermissionSetInfraDFQFilter description for the infra_dfq_filter field
	GroupDescPermissionSetInfraDFQFilter = "The scope binding for the dynamic filter query to restrict access to infrastructure assets"
	// GroupDescPermissionSetKubernetesClusterUUIDs description for the kubernetes_cluster_uuids field
	GroupDescPermissionSetKubernetesClusterUUIDs = "The scope bindings to restrict access to Kubernetes Clusters"
	// GroupDescPermissionSetKubernetesNamespaceUIDs description for the kubernetes_namespaces_uuids field
	GroupDescPermissionSetKubernetesNamespaceUIDs = "The scope bindings to restrict access to Kubernetes namespaces"
	// GroupDescPermissionSetMobileAppIDs description for the mobile_app_ids field
	GroupDescPermissionSetMobileAppIDs = "The scope bindings to restrict access to mobile apps"
	// GroupDescPermissionSetWebsiteIDs description for the website_ids field
	GroupDescPermissionSetWebsiteIDs = "The scope bindings to restrict access to websites"
	// GroupDescPermissionSetPermissions description for the permissions field
	GroupDescPermissionSetPermissions = "The permissions assigned which should be assigned to the users of the group"

	// Default values and special constants

	// DefaultScopeRoleID represents the default scope role ID for empty infra DFQ filters
	DefaultScopeRoleID = "-1"
	// EmptyScopeID represents an empty scope ID
	EmptyScopeID = ""
)
