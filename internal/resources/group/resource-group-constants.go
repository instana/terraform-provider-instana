package group

// ResourceInstanaGroupFramework the name of the terraform-provider-instana resource to manage groups for role based access control
const ResourceInstanaGroupFramework = "rbac_group"

//nolint:gosec
const (
	//GroupFieldName constant value for the schema field name
	GroupFieldName = "name"
	//GroupFieldFullName constant value for the schema field full_name
	GroupFieldFullName = "full_name"
	//GroupFieldMembers constant value for the schema field members
	GroupFieldMembers = "member"
	//GroupFieldMemberEmail constant value for the schema field email
	GroupFieldMemberEmail = "email"
	//GroupFieldMemberUserID constant value for the schema field user_id
	GroupFieldMemberUserID = "user_id"
	//GroupFieldPermissionSet constant value for the schema field permission_set
	GroupFieldPermissionSet = "permission_set"
	//GroupFieldPermissionSetApplicationIDs constant value for the schema field application_ids
	GroupFieldPermissionSetApplicationIDs = "application_ids"
	//GroupFieldPermissionSetInfraDFQFilter constant value for the schema field infra_dfq_filter
	GroupFieldPermissionSetInfraDFQFilter = "infra_dfq_filter"
	//GroupFieldPermissionSetKubernetesClusterUUIDs constant value for the schema field kubernetes_cluster_uuids
	GroupFieldPermissionSetKubernetesClusterUUIDs = "kubernetes_cluster_uuids"
	//GroupFieldPermissionSetKubernetesNamespaceUIDs constant value for the schema field kubernetes_namespaces_uuids
	GroupFieldPermissionSetKubernetesNamespaceUIDs = "kubernetes_namespaces_uuids"
	//GroupFieldPermissionSetMobileAppIDs constant value for the schema field mobile_app_ids
	GroupFieldPermissionSetMobileAppIDs = "mobile_app_ids"
	//GroupFieldPermissionSetWebsiteIDs constant value for the schema field website_ids
	GroupFieldPermissionSetWebsiteIDs = "website_ids"
	//GroupFieldPermissionSetPermissions constant value for the schema field permissions
	GroupFieldPermissionSetPermissions = "permissions"

	// Resource description constants

	GroupDescResource                             = "This resource manages RBAC groups in Instana."
	GroupDescID                                   = "The ID of the group."
	GroupDescName                                 = "The name of the Group"
	GroupDescMembers                              = "The members of the group"
	GroupDescMemberUserID                         = "The user id of the group member"
	GroupDescMemberEmail                          = "The email address of the group member"
	GroupDescPermissionSet                        = "The permission set of the group"
	GroupDescPermissionSetApplicationIDs          = "The scope bindings to restrict access to applications"
	GroupDescPermissionSetInfraDFQFilter          = "The scope binding for the dynamic filter query to restrict access to infrastructure assets"
	GroupDescPermissionSetKubernetesClusterUUIDs  = "The scope bindings to restrict access to Kubernetes Clusters"
	GroupDescPermissionSetKubernetesNamespaceUIDs = "The scope bindings to restrict access to Kubernetes namespaces"
	GroupDescPermissionSetMobileAppIDs            = "The scope bindings to restrict access to mobile apps"
	GroupDescPermissionSetWebsiteIDs              = "The scope bindings to restrict access to websites"
	GroupDescPermissionSetPermissions             = "The permissions assigned which should be assigned to the users of the group"
)
