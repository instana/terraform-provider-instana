package customdashboard

// ResourceInstanaCustomDashboardFramework the name of the terraform-provider-instana resource to manage custom dashboards
const ResourceInstanaCustomDashboardFramework = "custom_dashboard"

const (
	//CustomDashboardFieldTitle constant value for the schema field title
	CustomDashboardFieldTitle = "title"
	//CustomDashboardFieldFullTitle constant value for the computed schema field full_title
	CustomDashboardFieldFullTitle = "full_title"
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

// Error messages
const CustomDashboardErrMarshalWidgets = "Error marshaling widgets"
const CustomDashboardErrMarshalWidgetsFailed = "Failed to marshal widgets: %s"
