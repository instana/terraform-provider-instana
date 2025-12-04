package team

import "github.com/hashicorp/terraform-plugin-framework/types"

// TeamModel represents the data model for RBAC Team
type TeamModel struct {
	ID      types.String      `tfsdk:"id"`
	Tag     types.String      `tfsdk:"tag"`
	Info    *TeamInfoModel    `tfsdk:"info"`
	Members []TeamMemberModel `tfsdk:"member"`
	Scope   *TeamScopeModel   `tfsdk:"scope"`
}

// TeamInfoModel represents additional information about the team
type TeamInfoModel struct {
	Description types.String `tfsdk:"description"`
}

// TeamMemberModel represents a member in the team
type TeamMemberModel struct {
	UserID types.String     `tfsdk:"user_id"`
	Roles  []TeamMemberRole `tfsdk:"roles"`
}

// TeamMemberRole represents a role assigned to a team member
type TeamMemberRole struct {
	RoleID types.String `tfsdk:"role_id"`
}

// TeamScopeModel represents the scope configuration for the team
type TeamScopeModel struct {
	AccessPermissions           []string                              `tfsdk:"access_permissions"`
	Applications                []string                              `tfsdk:"applications"`
	KubernetesClusters          []string                              `tfsdk:"kubernetes_clusters"`
	KubernetesNamespaces        []string                              `tfsdk:"kubernetes_namespaces"`
	MobileApps                  []string                              `tfsdk:"mobile_apps"`
	Websites                    []string                              `tfsdk:"websites"`
	InfraDFQFilter              types.String                          `tfsdk:"infra_dfq_filter"`
	ActionFilter                types.String                          `tfsdk:"action_filter"`
	LogFilter                   types.String                          `tfsdk:"log_filter"`
	BusinessPerspectives        []string                              `tfsdk:"business_perspectives"`
	SloIDs                      []string                              `tfsdk:"slo_ids"`
	SyntheticTests              []string                              `tfsdk:"synthetic_tests"`
	SyntheticCredentials        []string                              `tfsdk:"synthetic_credentials"`
	TagIDs                      []string                              `tfsdk:"tag_ids"`
	RestrictedApplicationFilter *TeamRestrictedApplicationFilterModel `tfsdk:"restricted_application_filter"`
}

// TeamRestrictedApplicationFilterModel represents the restricted application filter configuration
type TeamRestrictedApplicationFilterModel struct {
	Label               types.String `tfsdk:"label"`
	Scope               types.String `tfsdk:"scope"`
	TagFilterExpression types.String `tfsdk:"tag_filter_expression"`
}
