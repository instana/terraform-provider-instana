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

// TeamScopeModel represents the scope configuration for the team.
// All set-typed scope fields use types.Set (not []string) because the schema marks them
// Optional+Computed — the Terraform Framework requires types.Set to correctly carry the
// unknown value during plan when a Computed field has not yet been resolved.
type TeamScopeModel struct {
	AccessPermissions           types.Set                             `tfsdk:"access_permissions"`
	Applications                types.Set                             `tfsdk:"applications"`
	KubernetesClusters          types.Set                             `tfsdk:"kubernetes_clusters"`
	KubernetesNamespaces        types.Set                             `tfsdk:"kubernetes_namespaces"`
	MobileApps                  types.Set                             `tfsdk:"mobile_apps"`
	Websites                    types.Set                             `tfsdk:"websites"`
	InfraDFQFilter              types.String                          `tfsdk:"infra_dfq_filter"`
	ActionFilter                types.String                          `tfsdk:"action_filter"`
	LogFilter                   types.String                          `tfsdk:"log_filter"`
	BusinessPerspectives        types.Set                             `tfsdk:"business_perspectives"`
	SloIDs                      types.Set                             `tfsdk:"slo_ids"`
	SyntheticTests              types.Set                             `tfsdk:"synthetic_tests"`
	SyntheticCredentials        types.Set                             `tfsdk:"synthetic_credentials"`
	TagIDs                      types.Set                             `tfsdk:"tag_ids"`
	RestrictedApplicationFilter *TeamRestrictedApplicationFilterModel `tfsdk:"restricted_application_filter"`
}

// TeamRestrictedApplicationFilterModel represents the restricted application filter configuration
type TeamRestrictedApplicationFilterModel struct {
	Label               types.String `tfsdk:"label"`
	Scope               types.String `tfsdk:"scope"`
	TagFilterExpression types.String `tfsdk:"tag_filter_expression"`
}
