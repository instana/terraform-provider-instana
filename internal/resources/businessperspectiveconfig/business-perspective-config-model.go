package businessperspectiveconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// BusinessPerspectiveConfigModel represents the data model for the business perspective configuration resource
type BusinessPerspectiveConfigModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TagFilter   types.String `tfsdk:"tag_filter"`
	RbacTags    types.List   `tfsdk:"rbac_tags"`
}

// RbacTagModel represents a single RBAC tag in Terraform state
type RbacTagModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}
