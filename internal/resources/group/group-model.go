package group

import "github.com/hashicorp/terraform-plugin-framework/types"

// GroupModel represents the data model for RBAC Group
type GroupModel struct {
	ID            types.String             `tfsdk:"id"`
	Name          types.String             `tfsdk:"name"`
	Members       []GroupMemberModel       `tfsdk:"member"`
	PermissionSet *GroupPermissionSetModel `tfsdk:"permission_set"`
}

// GroupMemberModel represents a member in the group
type GroupMemberModel struct {
	UserID types.String `tfsdk:"user_id"`
	Email  types.String `tfsdk:"email"`
}

// GroupPermissionSetModel represents the permission set for the group
type GroupPermissionSetModel struct {
	ApplicationIDs          []string     `tfsdk:"application_ids"`
	InfraDFQFilter          types.String `tfsdk:"infra_dfq_filter"`
	KubernetesClusterUUIDs  []string     `tfsdk:"kubernetes_cluster_uuids"`
	KubernetesNamespaceUIDs []string     `tfsdk:"kubernetes_namespaces_uuids"`
	MobileAppIDs            []string     `tfsdk:"mobile_app_ids"`
	WebsiteIDs              []string     `tfsdk:"website_ids"`
	Permissions             []string     `tfsdk:"permissions"`
}
