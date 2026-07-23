package groupmapping

import "github.com/hashicorp/terraform-plugin-framework/types"

// GroupMappingModel represents the data model for an RBAC Group Mapping
type GroupMappingModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.String `tfsdk:"value"`
	GroupID types.String `tfsdk:"group_id"`
	TeamID  types.String `tfsdk:"team_id"`
}
