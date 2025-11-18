package team

// ResourceInstanaTeamFramework the name of the terraform-provider-instana resource to manage teams for role based access control
const ResourceInstanaTeamFramework = "rbac_team"

//nolint:gosec
const (
	// Schema field names

	// TeamFieldID constant value for the schema field id
	TeamFieldID = "id"
	// TeamFieldTag constant value for the schema field tag
	TeamFieldTag = "tag"
	// TeamFieldInfo constant value for the schema field info
	TeamFieldInfo = "info"
	// TeamFieldInfoDescription constant value for the schema field description
	TeamFieldInfoDescription = "description"
	// TeamFieldMembers constant value for the schema field members
	TeamFieldMembers = "member"
	// TeamFieldMemberUserID constant value for the schema field user_id
	TeamFieldMemberUserID = "user_id"
	// TeamFieldMemberEmail constant value for the schema field email
	TeamFieldMemberEmail = "email"
	// TeamFieldMemberName constant value for the schema field name
	TeamFieldMemberName = "name"
	// TeamFieldMemberRoles constant value for the schema field roles
	TeamFieldMemberRoles = "roles"
	// TeamFieldMemberRoleID constant value for the schema field role_id
	TeamFieldMemberRoleID = "role_id"
	// TeamFieldMemberRoleName constant value for the schema field role_name
	TeamFieldMemberRoleName = "role_name"
	// TeamFieldMemberRoleViaIdP constant value for the schema field via_idp
	TeamFieldMemberRoleViaIdP = "via_idp"
	// TeamFieldScope constant value for the schema field scope
	TeamFieldScope = "scope"
	// TeamFieldScopeAccessPermissions constant value for the schema field access_permissions
	TeamFieldScopeAccessPermissions = "access_permissions"
	// TeamFieldScopeApplications constant value for the schema field applications
	TeamFieldScopeApplications = "applications"
	// TeamFieldScopeKubernetesClusters constant value for the schema field kubernetes_clusters
	TeamFieldScopeKubernetesClusters = "kubernetes_clusters"
	// TeamFieldScopeKubernetesNamespaces constant value for the schema field kubernetes_namespaces
	TeamFieldScopeKubernetesNamespaces = "kubernetes_namespaces"
	// TeamFieldScopeMobileApps constant value for the schema field mobile_apps
	TeamFieldScopeMobileApps = "mobile_apps"
	// TeamFieldScopeWebsites constant value for the schema field websites
	TeamFieldScopeWebsites = "websites"
	// TeamFieldScopeInfraDFQFilter constant value for the schema field infra_dfq_filter
	TeamFieldScopeInfraDFQFilter = "infra_dfq_filter"
	// TeamFieldScopeActionFilter constant value for the schema field action_filter
	TeamFieldScopeActionFilter = "action_filter"
	// TeamFieldScopeLogFilter constant value for the schema field log_filter
	TeamFieldScopeLogFilter = "log_filter"
	// TeamFieldScopeBusinessPerspectives constant value for the schema field business_perspectives
	TeamFieldScopeBusinessPerspectives = "business_perspectives"
	// TeamFieldScopeSloIDs constant value for the schema field slo_ids
	TeamFieldScopeSloIDs = "slo_ids"
	// TeamFieldScopeSyntheticTests constant value for the schema field synthetic_tests
	TeamFieldScopeSyntheticTests = "synthetic_tests"
	// TeamFieldScopeSyntheticCredentials constant value for the schema field synthetic_credentials
	TeamFieldScopeSyntheticCredentials = "synthetic_credentials"
	// TeamFieldScopeTagIDs constant value for the schema field tag_ids
	TeamFieldScopeTagIDs = "tag_ids"
	// TeamFieldScopeRestrictedApplicationFilter constant value for the schema field restricted_application_filter
	TeamFieldScopeRestrictedApplicationFilter = "restricted_application_filter"
	// TeamFieldScopeRestrictedApplicationFilterLabel constant value for the schema field label
	TeamFieldScopeRestrictedApplicationFilterLabel = "label"
	// TeamFieldScopeRestrictedApplicationFilterRestrictingApplicationID constant value for the schema field restricting_application_id
	TeamFieldScopeRestrictedApplicationFilterRestrictingApplicationID = "restricting_application_id"
	// TeamFieldScopeRestrictedApplicationFilterScope constant value for the schema field scope
	TeamFieldScopeRestrictedApplicationFilterScope = "scope"
	// TeamFieldScopeRestrictedApplicationFilterTagFilterExpression constant value for the schema field tag_filter_expression
	TeamFieldScopeRestrictedApplicationFilterTagFilterExpression = "tag_filter_expression"

	// Resource description constants

	// TeamDescResource description for the team resource
	TeamDescResource = "This resource manages RBAC teams in Instana."
	// TeamDescID description for the ID field
	TeamDescID = "The ID of the team."
	// TeamDescTag description for the tag field
	TeamDescTag = "The tag/name of the team"
	// TeamDescInfo description for the info field
	TeamDescInfo = "Additional information about the team"
	// TeamDescInfoDescription description for the info description field
	TeamDescInfoDescription = "The description of the team"
	// TeamDescMembers description for the members field
	TeamDescMembers = "The members of the team"
	// TeamDescMemberUserID description for the member user_id field
	TeamDescMemberUserID = "The user ID of the team member"
	// TeamDescMemberEmail description for the member email field
	TeamDescMemberEmail = "The email address of the team member"
	// TeamDescMemberName description for the member name field
	TeamDescMemberName = "The name of the team member"
	// TeamDescMemberRoles description for the member roles field
	TeamDescMemberRoles = "The roles assigned to the team member"
	// TeamDescMemberRoleID description for the role_id field
	TeamDescMemberRoleID = "The ID of the role"
	// TeamDescMemberRoleName description for the role_name field
	TeamDescMemberRoleName = "The name of the role"
	// TeamDescMemberRoleViaIdP description for the via_idp field
	TeamDescMemberRoleViaIdP = "Whether the role is assigned via IdP"
	// TeamDescScope description for the scope field
	TeamDescScope = "The scope configuration for the team"
	// TeamDescScopeAccessPermissions description for the access_permissions field
	TeamDescScopeAccessPermissions = "The access permissions for the team"
	// TeamDescScopeApplications description for the applications field
	TeamDescScopeApplications = "The application IDs accessible to the team"
	// TeamDescScopeKubernetesClusters description for the kubernetes_clusters field
	TeamDescScopeKubernetesClusters = "The Kubernetes cluster IDs accessible to the team"
	// TeamDescScopeKubernetesNamespaces description for the kubernetes_namespaces field
	TeamDescScopeKubernetesNamespaces = "The Kubernetes namespace IDs accessible to the team"
	// TeamDescScopeMobileApps description for the mobile_apps field
	TeamDescScopeMobileApps = "The mobile app IDs accessible to the team"
	// TeamDescScopeWebsites description for the websites field
	TeamDescScopeWebsites = "The website IDs accessible to the team"
	// TeamDescScopeInfraDFQFilter description for the infra_dfq_filter field
	TeamDescScopeInfraDFQFilter = "The infrastructure DFQ filter for the team"
	// TeamDescScopeActionFilter description for the action_filter field
	TeamDescScopeActionFilter = "The action filter for the team"
	// TeamDescScopeLogFilter description for the log_filter field
	TeamDescScopeLogFilter = "The log filter for the team"
	// TeamDescScopeBusinessPerspectives description for the business_perspectives field
	TeamDescScopeBusinessPerspectives = "The business perspective IDs accessible to the team"
	// TeamDescScopeSloIDs description for the slo_ids field
	TeamDescScopeSloIDs = "The SLO IDs accessible to the team"
	// TeamDescScopeSyntheticTests description for the synthetic_tests field
	TeamDescScopeSyntheticTests = "The synthetic test IDs accessible to the team"
	// TeamDescScopeSyntheticCredentials description for the synthetic_credentials field
	TeamDescScopeSyntheticCredentials = "The synthetic credential IDs accessible to the team"
	// TeamDescScopeTagIDs description for the tag_ids field
	TeamDescScopeTagIDs = "The tag IDs accessible to the team"
	// TeamDescScopeRestrictedApplicationFilter description for the restricted_application_filter field
	TeamDescScopeRestrictedApplicationFilter = "The restricted application filter configuration"
	// TeamDescScopeRestrictedApplicationFilterLabel description for the label field
	TeamDescScopeRestrictedApplicationFilterLabel = "The label for the restricted application filter"
	// TeamDescScopeRestrictedApplicationFilterRestrictingApplicationID description for the restricting_application_id field
	TeamDescScopeRestrictedApplicationFilterRestrictingApplicationID = "The ID of the restricting application"
	// TeamDescScopeRestrictedApplicationFilterScope description for the scope field
	TeamDescScopeRestrictedApplicationFilterScope = "The scope of the restricted application filter (INCLUDE_NO_DOWNSTREAM, INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING, INCLUDE_ALL_DOWNSTREAM)"
	// TeamDescScopeRestrictedApplicationFilterTagFilterExpression description for the tag_filter_expression field
	TeamDescScopeRestrictedApplicationFilterTagFilterExpression = "The tag filter expression for the restricted application filter"
)

// Made with Bob
