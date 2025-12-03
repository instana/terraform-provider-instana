package roles

import "github.com/hashicorp/terraform-plugin-framework/types"

// RoleModel represents the data model for RBAC Role
type RoleModel struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Members     []RoleMemberModel `tfsdk:"member"`
	Permissions []string          `tfsdk:"permissions"`
}

// RoleMemberModel represents a member in the role
type RoleMemberModel struct {
	UserID types.String `tfsdk:"user_id"`
}
