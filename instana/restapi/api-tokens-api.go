package restapi

// APITokensResourcePath path to API Tokens resource of Instana RESTful API
const APITokensResourcePath = SettingsBasePath + "/api-tokens"

// APIToken is the representation of a API Token in Instana
type APIToken struct {
	ID                                       string `json:"id"`
	AccessGrantingToken                      string `json:"accessGrantingToken"`
	InternalID                               string `json:"internalId"`
	Name                                     string `json:"name"`
	CanConfigureServiceMapping               bool   `json:"canConfigureServiceMapping"`
	CanConfigureEumApplications              bool   `json:"canConfigureEumApplications"`
	CanConfigureMobileAppMonitoring          bool   `json:"canConfigureMobileAppMonitoring"` //NEW
	CanConfigureUsers                        bool   `json:"canConfigureUsers"`
	CanInstallNewAgents                      bool   `json:"canInstallNewAgents"`
	CanConfigureIntegrations                 bool   `json:"canConfigureIntegrations"`
	CanConfigureEventsAndAlerts              bool   `json:"canConfigureEventsAndAlerts"`
	CanConfigureMaintenanceWindows           bool   `json:"canConfigureMaintenanceWindows"`
	CanConfigureApplicationSmartAlerts       bool   `json:"canConfigureApplicationSmartAlerts"`
	CanConfigureWebsiteSmartAlerts           bool   `json:"canConfigureWebsiteSmartAlerts"`
	CanConfigureMobileAppSmartAlerts         bool   `json:"canConfigureMobileAppSmartAlerts"`
	CanConfigureAPITokens                    bool   `json:"canConfigureApiTokens"`
	CanConfigureAgentRunMode                 bool   `json:"canConfigureAgentRunMode"`
	CanViewAuditLog                          bool   `json:"canViewAuditLog"`
	CanConfigureAgents                       bool   `json:"canConfigureAgents"`
	CanConfigureAuthenticationMethods        bool   `json:"canConfigureAuthenticationMethods"`
	CanConfigureApplications                 bool   `json:"canConfigureApplications"`
	CanConfigureTeams                        bool   `json:"canConfigureTeams"`
	CanConfigureReleases                     bool   `json:"canConfigureReleases"`
	CanConfigureLogManagement                bool   `json:"canConfigureLogManagement"`
	CanCreatePublicCustomDashboards          bool   `json:"canCreatePublicCustomDashboards"`
	CanViewLogs                              bool   `json:"canViewLogs"`
	CanViewTraceDetails                      bool   `json:"canViewTraceDetails"`
	CanConfigureSessionSettings              bool   `json:"canConfigureSessionSettings"`
	CanConfigureServiceLevelIndicators       bool   `json:"canConfigureServiceLevelIndicators"`
	CanConfigureGlobalAlertPayload           bool   `json:"canConfigureGlobalAlertPayload"`
	CanConfigureGlobalApplicationSmartAlerts bool   `json:"canConfigureGlobalApplicationSmartAlerts"`
	CanConfigureGlobalSyntheticSmartAlerts   bool   `json:"canConfigureGlobalSyntheticSmartAlerts"`
	CanConfigureGlobalInfraSmartAlerts       bool   `json:"canConfigureGlobalInfraSmartAlerts"`
	CanConfigureGlobalLogSmartAlerts         bool   `json:"canConfigureGlobalLogSmartAlerts"`
	CanViewAccountAndBillingInformation      bool   `json:"canViewAccountAndBillingInformation"`
	CanEditAllAccessibleCustomDashboards     bool   `json:"canEditAllAccessibleCustomDashboards"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (r *APIToken) GetIDForResourcePath() string {
	return r.InternalID
}
