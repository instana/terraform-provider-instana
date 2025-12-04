package customdashboard

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomDashboardModel represents the data model for the custom dashboard resource
type CustomDashboardModel struct {
	ID          types.String         `tfsdk:"id"`
	Title       types.String         `tfsdk:"title"`
	AccessRules []AccessRuleModel    `tfsdk:"access_rule"`
	Widgets     jsontypes.Normalized `tfsdk:"widgets"`
}

// AccessRuleModel represents an access rule in the custom dashboard
type AccessRuleModel struct {
	AccessType   types.String `tfsdk:"access_type"`
	RelatedID    types.String `tfsdk:"related_id"`
	RelationType types.String `tfsdk:"relation_type"`
}
