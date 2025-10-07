package group

import "github.com/hashicorp/terraform-plugin-framework/types"

// GroupModel represents the data model for RBAC Group
type GroupModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Members       types.Set    `tfsdk:"member"`
	PermissionSet types.Object `tfsdk:"permission_set"`
}

// GroupMemberModel represents a member in the group
type GroupMemberModel struct {
	UserID types.String `tfsdk:"user_id"`
	Email  types.String `tfsdk:"email"`
}

// GroupPermissionSetModel represents the permission set for the group
type GroupPermissionSetModel struct {
	ApplicationIDs          types.Set    `tfsdk:"application_ids"`
	InfraDFQFilter          types.String `tfsdk:"infra_dfq_filter"`
	KubernetesClusterUUIDs  types.Set    `tfsdk:"kubernetes_cluster_uuids"`
	KubernetesNamespaceUIDs types.Set    `tfsdk:"kubernetes_namespaces_uuids"`
	MobileAppIDs            types.Set    `tfsdk:"mobile_app_ids"`
	WebsiteIDs              types.Set    `tfsdk:"website_ids"`
	Permissions             types.Set    `tfsdk:"permissions"`
}
