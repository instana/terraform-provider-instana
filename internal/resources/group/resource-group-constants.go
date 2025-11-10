package group

// Resource description constants
const (
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

// Made with Bob
