package restapi

// InstanaPermission data type representing an Instana permission string
type InstanaPermission string

const (
	//PermissionCanConfigureApplications const for Instana permission CAN_CONFIGURE_APPLICATIONS
	PermissionCanConfigureApplications = InstanaPermission("CAN_CONFIGURE_APPLICATIONS")
	//PermissionCanConfigureEumApplications const for Instana permission CAN_CONFIGURE_EUM_APPLICATIONS
	PermissionCanConfigureEumApplications = InstanaPermission("CAN_CONFIGURE_EUM_APPLICATIONS")
	//PermissionCanConfigureAgents const for Instana permission CAN_CONFIGURE_AGENTS
	PermissionCanConfigureAgents = InstanaPermission("CAN_CONFIGURE_AGENTS")
	//PermissionCanViewTraceDetails const for Instana permission CAN_VIEW_TRACE_DETAILS
	PermissionCanViewTraceDetails = InstanaPermission("CAN_VIEW_TRACE_DETAILS")
	//PermissionCanViewLogs const for Instana permission CAN_VIEW_LOGS
	PermissionCanViewLogs = InstanaPermission("CAN_VIEW_LOGS")
	//PermissionCanConfigureSessionSettings const for Instana permission CAN_CONFIGURE_SESSION_SETTINGS
	PermissionCanConfigureSessionSettings = InstanaPermission("CAN_CONFIGURE_SESSION_SETTINGS")
	//PermissionCanConfigureIntegrations const for Instana permission CAN_CONFIGURE_INTEGRATIONS
	PermissionCanConfigureIntegrations = InstanaPermission("CAN_CONFIGURE_INTEGRATIONS")
	// PermissionCanConfigureGlobalApplicationSmartAlerts Permission to configure Global Application Smart Alerts
	PermissionCanConfigureGlobalApplicationSmartAlerts = InstanaPermission("CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS")
	// PermissionCanConfigureGlobalSyntheticSmartAlerts Permission to configure Global Synthetic Smart Alerts
	PermissionCanConfigureGlobalSyntheticSmartAlerts = InstanaPermission("CAN_CONFIGURE_GLOBAL_SYNTHETIC_SMART_ALERTS")
	// PermissionCanConfigureGlobalInfraSmartAlerts Permission to configure Global Infrastructure Smart Alerts
	PermissionCanConfigureGlobalInfraSmartAlerts = InstanaPermission("CAN_CONFIGURE_GLOBAL_INFRA_SMART_ALERTS")
	// PermissionCanConfigureGlobalLogSmartAlerts Permission to configure Global Log Smart Alerts
	PermissionCanConfigureGlobalLogSmartAlerts = InstanaPermission("CAN_CONFIGURE_GLOBAL_LOG_SMART_ALERTS")
	//PermissionCanConfigureGlobalAlertPayload const for Instana permission CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD
	PermissionCanConfigureGlobalAlertPayload = InstanaPermission("CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD")
	//PermissionCanConfigureMobileAppMonitoring const for Instana permission CAN_CONFIGURE_MOBILE_APP_MONITORING
	PermissionCanConfigureMobileAppMonitoring = InstanaPermission("CAN_CONFIGURE_MOBILE_APP_MONITORING")
	//PermissionCanConfigureAPITokens const for Instana permission CAN_CONFIGURE_API_TOKENS
	PermissionCanConfigureAPITokens = InstanaPermission("CAN_CONFIGURE_API_TOKENS")
	//PermissionCanConfigureServiceLevelIndicators const for Instana permission CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS
	PermissionCanConfigureServiceLevelIndicators = InstanaPermission("CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS")
	//PermissionCanConfigureAuthenticationMethods const for Instana permission CAN_CONFIGURE_AUTHENTICATION_METHODS
	PermissionCanConfigureAuthenticationMethods = InstanaPermission("CAN_CONFIGURE_AUTHENTICATION_METHODS")
	//PermissionCanConfigureReleases const for Instana permission CAN_CONFIGURE_RELEASES
	PermissionCanConfigureReleases = InstanaPermission("CAN_CONFIGURE_RELEASES")
	//PermissionCanViewAuditLog const for Instana permission CAN_VIEW_AUDIT_LOG
	PermissionCanViewAuditLog = InstanaPermission("CAN_VIEW_AUDIT_LOG")
	// PermissionCanConfigureEventsAndAlerts Permission to configure Instana Events and Alerts
	PermissionCanConfigureEventsAndAlerts = InstanaPermission("CAN_CONFIGURE_EVENTS_AND_ALERTS")
	// PermissionCanConfigureMaintenanceWindows Permission to configure Instana Maintenance Windows
	PermissionCanConfigureMaintenanceWindows = InstanaPermission("CAN_CONFIGURE_MAINTENANCE_WINDOWS")
	// PermissionCanConfigureApplicationSmartAlerts Permission to configure Instana Application Smart Alerts
	PermissionCanConfigureApplicationSmartAlerts = InstanaPermission("CAN_CONFIGURE_APPLICATION_SMART_ALERTS")
	// PermissionCanConfigureWebsiteSmartAlerts Permission to configure Instana Website Smart Alerts
	PermissionCanConfigureWebsiteSmartAlerts = InstanaPermission("CAN_CONFIGURE_WEBSITE_SMART_ALERTS")
	// PermissionCanConfigureMobileAppSmartAlerts Permission to configure Instana MobileApp Smart Alerts
	PermissionCanConfigureMobileAppSmartAlerts = InstanaPermission("CAN_CONFIGURE_MOBILE_APP_SMART_ALERTS")
	//PermissionCanConfigureAgentRunMode const for Instana permission CAN_CONFIGURE_AGENT_RUN_MODE
	PermissionCanConfigureAgentRunMode = InstanaPermission("CAN_CONFIGURE_AGENT_RUN_MODE")
	//PermissionCanConfigureServiceMapping const for Instana permission CAN_CONFIGURE_SERVICE_MAPPING
	PermissionCanConfigureServiceMapping = InstanaPermission("CAN_CONFIGURE_SERVICE_MAPPING")
	//PermissionCanEditAllAccessibleCustomDashboards const for Instana permission CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS
	PermissionCanEditAllAccessibleCustomDashboards = InstanaPermission("CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS")
	//PermissionCanConfigureUsers const for Instana permission CAN_CONFIGURE_USERS
	PermissionCanConfigureUsers = InstanaPermission("CAN_CONFIGURE_USERS")
	//PermissionCanInstallNewAgents const for Instana permission CAN_INSTALL_NEW_AGENTS
	PermissionCanInstallNewAgents = InstanaPermission("CAN_INSTALL_NEW_AGENTS")
	//PermissionCanConfigureTeams const for Instana permission CAN_CONFIGURE_TEAMS
	PermissionCanConfigureTeams = InstanaPermission("CAN_CONFIGURE_TEAMS")
	//PermissionCanCreatePublicCustomDashboards const for Instana permission CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS
	PermissionCanCreatePublicCustomDashboards = InstanaPermission("CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS")
	//PermissionCanConfigureLogManagement const for Instana permission CAN_CONFIGURE_LOG_MANAGEMENT
	PermissionCanConfigureLogManagement = InstanaPermission("CAN_CONFIGURE_LOG_MANAGEMENT")
	//PermissionCanViewAccountAndBillingInformation const for Instana permission CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION
	PermissionCanViewAccountAndBillingInformation = InstanaPermission("CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION")
	//PermissionCanViewSyntheticTests const for Instana permission CAN_VIEW_SYNTHETIC_TESTS
	PermissionCanViewSyntheticTests = InstanaPermission("CAN_VIEW_SYNTHETIC_TESTS")
	//PermissionCanViewSyntheticLocations const for Instana permission CAN_VIEW_SYNTHETIC_LOCATIONS
	PermissionCanViewSyntheticLocations = InstanaPermission("CAN_VIEW_SYNTHETIC_LOCATIONS")
	//PermissionCanCreateThreadDump const for Instana permission CAN_CREATE_THREAD_DUMP
	PermissionCanCreateThreadDump = InstanaPermission("CAN_CREATE_THREAD_DUMP")
	//PermissionCanCreateHeapDump const for Instana permission CAN_CREATE_HEAP_DUMP
	PermissionCanCreateHeapDump = InstanaPermission("CAN_CREATE_HEAP_DUMP")
	//PermissionCanConfigureDatabaseManagement const for Instana permission CAN_CONFIGURE_DATABASE_MANAGEMENT
	PermissionCanConfigureDatabaseManagement = InstanaPermission("CAN_CONFIGURE_DATABASE_MANAGEMENT")
	//PermissionCanConfigureLogRetentionPeriod const for Instana permission CAN_CONFIGURE_LOG_RETENTION_PERIOD
	PermissionCanConfigureLogRetentionPeriod = InstanaPermission("CAN_CONFIGURE_LOG_RETENTION_PERIOD")
	//PermissionCanConfigurePersonalAPITokens const for Instana permission CAN_CONFIGURE_PERSONAL_API_TOKENS
	PermissionCanConfigurePersonalAPITokens = InstanaPermission("CAN_CONFIGURE_PERSONAL_API_TOKENS")
	//PermissionAccessInfrastructureAnalyze const for Instana permission ACCESS_INFRASTRUCTURE_ANALYZE
	PermissionAccessInfrastructureAnalyze = InstanaPermission("ACCESS_INFRASTRUCTURE_ANALYZE")
	//PermissionCanViewLogVolume const for Instana permission CAN_VIEW_LOG_VOLUME
	PermissionCanViewLogVolume = InstanaPermission("CAN_VIEW_LOG_VOLUME")
	//PermissionCanRunAutomationActions const for Instana permission CAN_RUN_AUTOMATION_ACTIONS
	PermissionCanRunAutomationActions = InstanaPermission("CAN_RUN_AUTOMATION_ACTIONS")
	//PermissionCanViewSyntheticTestResults const for Instana permission CAN_VIEW_SYNTHETIC_TEST_RESULTS
	PermissionCanViewSyntheticTestResults = InstanaPermission("CAN_VIEW_SYNTHETIC_TEST_RESULTS")
	//PermissionCanInvokeAlertChannel const for Instana permission CAN_INVOKE_ALERT_CHANNEL
	PermissionCanInvokeAlertChannel = InstanaPermission("CAN_INVOKE_ALERT_CHANNEL")
	//PermissionCanManuallyCloseIssue const for Instana permission CAN_MANUALLY_CLOSE_ISSUE
	PermissionCanManuallyCloseIssue = InstanaPermission("CAN_MANUALLY_CLOSE_ISSUE")
	//PermissionCanDeleteLogs const for Instana permission CAN_DELETE_LOGS
	PermissionCanDeleteLogs = InstanaPermission("CAN_DELETE_LOGS")
	//PermissionCanConfigureSyntheticTests const for Instana permission CAN_CONFIGURE_SYNTHETIC_TESTS
	PermissionCanConfigureSyntheticTests = InstanaPermission("CAN_CONFIGURE_SYNTHETIC_TESTS")
	//PermissionCanViewBusinessProcessDetails const for Instana permission CAN_VIEW_BUSINESS_PROCESS_DETAILS
	PermissionCanViewBusinessProcessDetails = InstanaPermission("CAN_VIEW_BUSINESS_PROCESS_DETAILS")
	//PermissionCanViewBizOpsAlerts const for Instana permission CAN_VIEW_BIZOPS_ALERTS
	PermissionCanViewBizOpsAlerts = InstanaPermission("CAN_VIEW_BIZOPS_ALERTS")
	//PermissionCanUseSyntheticCredentials const for Instana permission CAN_USE_SYNTHETIC_CREDENTIALS
	PermissionCanUseSyntheticCredentials = InstanaPermission("CAN_USE_SYNTHETIC_CREDENTIALS")
	//PermissionCanDeleteAutomationActionHistory const for Instana permission CAN_DELETE_AUTOMATION_ACTION_HISTORY
	PermissionCanDeleteAutomationActionHistory = InstanaPermission("CAN_DELETE_AUTOMATION_ACTION_HISTORY")
	//PermissionCanConfigureSyntheticLocations const for Instana permission CAN_CONFIGURE_SYNTHETIC_LOCATIONS
	PermissionCanConfigureSyntheticLocations = InstanaPermission("CAN_CONFIGURE_SYNTHETIC_LOCATIONS")
	//PermissionCanConfigureSyntheticCredentials const for Instana permission CAN_CONFIGURE_SYNTHETIC_CREDENTIALS
	PermissionCanConfigureSyntheticCredentials = InstanaPermission("CAN_CONFIGURE_SYNTHETIC_CREDENTIALS")
	//PermissionCanConfigureSubtraces const for Instana permission CAN_CONFIGURE_SUBTRACES
	PermissionCanConfigureSubtraces = InstanaPermission("CAN_CONFIGURE_SUBTRACES")
	//PermissionCanConfigureLLM const for Instana permission CAN_CONFIGURE_LLM
	PermissionCanConfigureLLM = InstanaPermission("CAN_CONFIGURE_LLM")
	//PermissionCanConfigureBizOps const for Instana permission CAN_CONFIGURE_BIZOPS
	PermissionCanConfigureBizOps = InstanaPermission("CAN_CONFIGURE_BIZOPS")
	//PermissionCanConfigureAutomationPolicies const for Instana permission CAN_CONFIGURE_AUTOMATION_POLICIES
	PermissionCanConfigureAutomationPolicies = InstanaPermission("CAN_CONFIGURE_AUTOMATION_POLICIES")
	//PermissionCanConfigureAutomationActions const for Instana permission CAN_CONFIGURE_AUTOMATION_ACTIONS
	PermissionCanConfigureAutomationActions = InstanaPermission("CAN_CONFIGURE_AUTOMATION_ACTIONS")
)

// InstanaPermissions data type representing a slice of Instana permissions
type InstanaPermissions []InstanaPermission

// ToStringSlice converts the slice of InstanaPermissions to its string representation
func (permissions InstanaPermissions) ToStringSlice() []string {
	result := make([]string, len(permissions))
	for i, v := range permissions {
		result[i] = string(v)
	}
	return result
}

// SupportedInstanaPermissions slice of all supported Permissions of the Instana API
var SupportedInstanaPermissions = InstanaPermissions{
	PermissionCanConfigureApplications,
	PermissionCanConfigureEumApplications,
	PermissionCanConfigureAgents,
	PermissionCanViewTraceDetails,
	PermissionCanViewLogs,
	PermissionCanConfigureSessionSettings,
	PermissionCanConfigureIntegrations,
	PermissionCanConfigureGlobalApplicationSmartAlerts,
	PermissionCanConfigureGlobalSyntheticSmartAlerts,
	PermissionCanConfigureGlobalInfraSmartAlerts,
	PermissionCanConfigureGlobalLogSmartAlerts,
	PermissionCanConfigureGlobalAlertPayload,
	PermissionCanConfigureMobileAppMonitoring,
	PermissionCanConfigureAPITokens,
	PermissionCanConfigureServiceLevelIndicators,
	PermissionCanConfigureAuthenticationMethods,
	PermissionCanConfigureReleases,
	PermissionCanViewAuditLog,
	PermissionCanConfigureEventsAndAlerts,
	PermissionCanConfigureMaintenanceWindows,
	PermissionCanConfigureApplicationSmartAlerts,
	PermissionCanConfigureWebsiteSmartAlerts,
	PermissionCanConfigureMobileAppSmartAlerts,
	PermissionCanConfigureAgentRunMode,
	PermissionCanConfigureServiceMapping,
	PermissionCanEditAllAccessibleCustomDashboards,
	PermissionCanConfigureUsers,
	PermissionCanInstallNewAgents,
	PermissionCanConfigureTeams,
	PermissionCanCreatePublicCustomDashboards,
	PermissionCanConfigureLogManagement,
	PermissionCanViewAccountAndBillingInformation,
	PermissionCanViewSyntheticTests,
	PermissionCanViewSyntheticLocations,
	PermissionCanCreateThreadDump,
	PermissionCanCreateHeapDump,
	PermissionCanConfigureDatabaseManagement,
	PermissionCanConfigureLogRetentionPeriod,
	PermissionCanConfigurePersonalAPITokens,
	PermissionAccessInfrastructureAnalyze,
	PermissionCanViewLogVolume,
	PermissionCanRunAutomationActions,
	PermissionCanViewSyntheticTestResults,
	PermissionCanInvokeAlertChannel,
	PermissionCanManuallyCloseIssue,
	PermissionCanDeleteLogs,
	PermissionCanConfigureSyntheticTests,
	PermissionCanViewBusinessProcessDetails,
	PermissionCanViewBizOpsAlerts,
	PermissionCanUseSyntheticCredentials,
	PermissionCanDeleteAutomationActionHistory,
	PermissionCanConfigureSyntheticLocations,
	PermissionCanConfigureSyntheticCredentials,
	PermissionCanConfigureSubtraces,
	PermissionCanConfigureLLM,
	PermissionCanConfigureBizOps,
	PermissionCanConfigureAutomationPolicies,
	PermissionCanConfigureAutomationActions,
}

// GroupsResourcePath path to Group resource of Instana RESTful API
const GroupsResourcePath = RBACSettingsBasePath + "/groups"

// ScopeBinding data structure for the Instana API model for scope bindings
type ScopeBinding struct {
	ScopeID     string  `json:"scopeId"`
	ScopeRoleID *string `json:"scopeRoleId"`
}

// APIPermissionSetWithRoles data structure for the Instana API model for permissions with roles
type APIPermissionSetWithRoles struct {
	ApplicationIDs          []ScopeBinding      `json:"applicationIds"`
	InfraDFQFilter          *ScopeBinding       `json:"infraDfqFilter"`
	KubernetesClusterUUIDs  []ScopeBinding      `json:"kubernetesClusterUUIDs"`
	KubernetesNamespaceUIDs []ScopeBinding      `json:"kubernetesNamespaceUIDs"`
	MobileAppIDs            []ScopeBinding      `json:"mobileAppIds"`
	WebsiteIDs              []ScopeBinding      `json:"websiteIds"`
	Permissions             []InstanaPermission `json:"permissions"`
}

// IsEmpty returns true when no permission or scope is assigned
func (m *APIPermissionSetWithRoles) IsEmpty() bool {
	if len(m.ApplicationIDs) > 0 {
		return false
	}
	if len(m.KubernetesClusterUUIDs) > 0 {
		return false
	}
	if len(m.KubernetesNamespaceUIDs) > 0 {
		return false
	}
	if len(m.MobileAppIDs) > 0 {
		return false
	}
	if len(m.WebsiteIDs) > 0 {
		return false
	}
	if len(m.Permissions) > 0 {
		return false
	}
	if m.InfraDFQFilter != nil && len(m.InfraDFQFilter.ScopeID) > 0 {
		return false
	}
	return true
}

// APIMember data structure for the Instana API model for group members
type APIMember struct {
	UserID string  `json:"userId"`
	Email  *string `json:"email"`
}

// Group data structure for the Instana API model for groups
type Group struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Members       []APIMember               `json:"members"`
	PermissionSet APIPermissionSetWithRoles `json:"permissionSet"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (c *Group) GetIDForResourcePath() string {
	return c.ID
}
