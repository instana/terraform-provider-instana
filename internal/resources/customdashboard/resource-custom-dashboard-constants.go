package customdashboard

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

// Made with Bob
