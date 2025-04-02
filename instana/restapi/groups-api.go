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
