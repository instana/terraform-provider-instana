package applicationconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// ApplicationConfigModel represents the data model for the application configuration resource
type ApplicationConfigModel struct {
	ID            types.String `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	Scope         types.String `tfsdk:"scope"`
	BoundaryScope types.String `tfsdk:"boundary_scope"`
	TagFilter     types.String `tfsdk:"tag_filter"`
	AccessRules   types.List   `tfsdk:"access_rules"`
}
