package customdashboard

// ResourceInstanaCustomDashboard the name of the terraform-provider-instana resource to manage custom dashboards
const ResourceInstanaCustomDashboard = "custom_dashboard"

const (
	//CustomDashboardFieldID constant value for the schema field id
	CustomDashboardFieldID = "id"
	//CustomDashboardFieldTitle constant value for the schema field title
	CustomDashboardFieldTitle = "title"
	//CustomDashboardFieldAccessRule constant value for the schema field access_rule
	CustomDashboardFieldAccessRule = "access_rule"
	//CustomDashboardFieldAccessRuleAccessType constant value for the schema field access_rule.access_type
	CustomDashboardFieldAccessRuleAccessType = "access_type"
	//CustomDashboardFieldAccessRuleRelatedID constant value for the schema field access_rule.related_id
	CustomDashboardFieldAccessRuleRelatedID = "related_id"
	//CustomDashboardFieldAccessRuleRelationType constant value for the schema field access_rule.relation_type
	CustomDashboardFieldAccessRuleRelationType = "relation_type"
	//CustomDashboardFieldWidgets constant value for the schema field widgets
	CustomDashboardFieldWidgets = "widgets"
	//CustomDashboardFieldRbacTags constant value for the schema field rbac_tags
	CustomDashboardFieldRbacTags = "rbac_tags"
	//CustomDashboardFieldRbacTagDisplayName constant value for the schema field rbac_tags.display_name
	CustomDashboardFieldRbacTagDisplayName = "display_name"
	//CustomDashboardFieldRbacTagID constant value for the schema field rbac_tags.id
	CustomDashboardFieldRbacTagID = "id"
)

// Resource description
const CustomDashboardDescResource = "This resource manages custom dashboards in Instana."

// Field descriptions - ID
const CustomDashboardDescID = "The ID of the custom dashboard."

// Field descriptions - Basic fields
const CustomDashboardDescTitle = "The title of the custom dashboard."
const CustomDashboardDescWidgets = "The json array containing the widgets configured for the custom dashboard."

// Field descriptions - Access Rule
const CustomDashboardDescAccessRule = "The access rules applied to the custom dashboard."
const CustomDashboardDescAccessRuleAccessType = "The access type of the given access rule."
const CustomDashboardDescAccessRuleRelatedID = "The id of the related entity (user, api_token, etc.) of the given access rule."
const CustomDashboardDescAccessRuleRelationType = "The relation type of the given access rule."

// Field descriptions - RBAC Tags
const CustomDashboardDescRbacTags = "RBAC tags (teams) the custom dashboard is assigned to."
const CustomDashboardDescRbacTagDisplayName = "Display name of the RBAC tag (team)."
const CustomDashboardDescRbacTagID = "ID of the RBAC tag (team)."

// Error messages
const CustomDashboardErrMarshalWidgets = "Error marshaling widgets"
const CustomDashboardErrMarshalWidgetsFailed = "Failed to marshal widgets: %s"
