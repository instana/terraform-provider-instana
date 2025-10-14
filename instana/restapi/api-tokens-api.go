package restapi

// APITokensResourcePath path to API Tokens resource of Instana RESTful API
const APITokensResourcePath = SettingsBasePath + "/api-tokens"

// APIToken is the representation of a API Token in Instana
type APIToken struct {
	ID                                        string `json:"id"`
	AccessGrantingToken                       string `json:"accessGrantingToken"`
	InternalID                                string `json:"internalId"`
	Name                                      string `json:"name"`
	CanConfigureServiceMapping                bool   `json:"canConfigureServiceMapping"`
	CanConfigureEumApplications               bool   `json:"canConfigureEumApplications"`
	CanConfigureMobileAppMonitoring           bool   `json:"canConfigureMobileAppMonitoring"`
	CanConfigureUsers                         bool   `json:"canConfigureUsers"`
	CanInstallNewAgents                       bool   `json:"canInstallNewAgents"`
	CanConfigureIntegrations                  bool   `json:"canConfigureIntegrations"`
	CanConfigureEventsAndAlerts               bool   `json:"canConfigureEventsAndAlerts"`
	CanConfigureMaintenanceWindows            bool   `json:"canConfigureMaintenanceWindows"`
	CanConfigureApplicationSmartAlerts        bool   `json:"canConfigureApplicationSmartAlerts"`
	CanConfigureWebsiteSmartAlerts            bool   `json:"canConfigureWebsiteSmartAlerts"`
	CanConfigureMobileAppSmartAlerts          bool   `json:"canConfigureMobileAppSmartAlerts"`
	CanConfigureAPITokens                     bool   `json:"canConfigureApiTokens"`
	CanConfigureAgentRunMode                  bool   `json:"canConfigureAgentRunMode"`
	CanViewAuditLog                           bool   `json:"canViewAuditLog"`
	CanConfigureAgents                        bool   `json:"canConfigureAgents"`
	CanConfigureAuthenticationMethods         bool   `json:"canConfigureAuthenticationMethods"`
	CanConfigureApplications                  bool   `json:"canConfigureApplications"`
	CanConfigureTeams                         bool   `json:"canConfigureTeams"`
	CanConfigureReleases                      bool   `json:"canConfigureReleases"`
	CanConfigureLogManagement                 bool   `json:"canConfigureLogManagement"`
	CanCreatePublicCustomDashboards           bool   `json:"canCreatePublicCustomDashboards"`
	CanViewLogs                               bool   `json:"canViewLogs"`
	CanViewTraceDetails                       bool   `json:"canViewTraceDetails"`
	CanConfigureSessionSettings               bool   `json:"canConfigureSessionSettings"`
	CanConfigureGlobalAlertPayload            bool   `json:"canConfigureGlobalAlertPayload"`
	CanConfigureGlobalApplicationSmartAlerts  bool   `json:"canConfigureGlobalApplicationSmartAlerts"`
	CanConfigureGlobalSyntheticSmartAlerts    bool   `json:"canConfigureGlobalSyntheticSmartAlerts"`
	CanConfigureGlobalInfraSmartAlerts        bool   `json:"canConfigureGlobalInfraSmartAlerts"`
	CanConfigureGlobalLogSmartAlerts          bool   `json:"canConfigureGlobalLogSmartAlerts"`
	CanViewAccountAndBillingInformation       bool   `json:"canViewAccountAndBillingInformation"`
	CanEditAllAccessibleCustomDashboards      bool   `json:"canEditAllAccessibleCustomDashboards"`
	LimitedApplicationsScope                  bool   `json:"limitedApplicationsScope"`
	LimitedBizOpsScope                        bool   `json:"limitedBizOpsScope"`
	LimitedWebsitesScope                      bool   `json:"limitedWebsitesScope"`
	LimitedKubernetesScope                    bool   `json:"limitedKubernetesScope"`
	LimitedMobileAppsScope                    bool   `json:"limitedMobileAppsScope"`
	LimitedInfrastructureScope                bool   `json:"limitedInfrastructureScope"`
	LimitedSyntheticsScope                    bool   `json:"limitedSyntheticsScope"`
	LimitedVsphereScope                       bool   `json:"limitedVsphereScope"`
	LimitedPhmcScope                          bool   `json:"limitedPhmcScope"`
	LimitedPvcScope                           bool   `json:"limitedPvcScope"`
	LimitedZhmcScope                          bool   `json:"limitedZhmcScope"`
	LimitedPcfScope                           bool   `json:"limitedPcfScope"`
	LimitedOpenstackScope                     bool   `json:"limitedOpenstackScope"`
	LimitedAutomationScope                    bool   `json:"limitedAutomationScope"`
	LimitedLogsScope                          bool   `json:"limitedLogsScope"`
	LimitedNutanixScope                       bool   `json:"limitedNutanixScope"`
	LimitedXenServerScope                     bool   `json:"limitedXenServerScope"`
	LimitedWindowsHypervisorScope             bool   `json:"limitedWindowsHypervisorScope"`
	LimitedAlertChannelsScope                 bool   `json:"limitedAlertChannelsScope"`
	LimitedLinuxKvmHypervisorScope            bool   `json:"limitedLinuxKvmHypervisorScope"`
	LimitedAiGatewayScope                     bool   `json:"limitedAiGatewayScope"`
	LimitedGenAIScope                         bool   `json:"limitedGenAIScope"`
	LimitedServiceLevelScope                  bool   `json:"limitedServiceLevelScope"`
	CanConfigurePersonalAPITokens             bool   `json:"canConfigurePersonalApiTokens"`
	CanConfigureDatabaseManagement            bool   `json:"canConfigureDatabaseManagement"`
	CanConfigureAutomationActions             bool   `json:"canConfigureAutomationActions"`
	CanConfigureAutomationPolicies            bool   `json:"canConfigureAutomationPolicies"`
	CanRunAutomationActions                   bool   `json:"canRunAutomationActions"`
	CanDeleteAutomationActionHistory          bool   `json:"canDeleteAutomationActionHistory"`
	CanConfigureSyntheticTests                bool   `json:"canConfigureSyntheticTests"`
	CanConfigureSyntheticLocations            bool   `json:"canConfigureSyntheticLocations"`
	CanConfigureSyntheticCredentials          bool   `json:"canConfigureSyntheticCredentials"`
	CanViewSyntheticTests                     bool   `json:"canViewSyntheticTests"`
	CanViewSyntheticLocations                 bool   `json:"canViewSyntheticLocations"`
	CanViewSyntheticTestResults               bool   `json:"canViewSyntheticTestResults"`
	CanUseSyntheticCredentials                bool   `json:"canUseSyntheticCredentials"`
	CanConfigureBizops                        bool   `json:"canConfigureBizops"`
	CanViewBusinessProcesses                  bool   `json:"canViewBusinessProcesses"`
	CanViewBusinessProcessDetails             bool   `json:"canViewBusinessProcessDetails"`
	CanViewBusinessActivities                 bool   `json:"canViewBusinessActivities"`
	CanViewBizAlerts                          bool   `json:"canViewBizAlerts"`
	CanDeleteLogs                             bool   `json:"canDeleteLogs"`
	CanCreateHeapDump                         bool   `json:"canCreateHeapDump"`
	CanCreateThreadDump                       bool   `json:"canCreateThreadDump"`
	CanManuallyCloseIssue                     bool   `json:"canManuallyCloseIssue"`
	CanViewLogVolume                          bool   `json:"canViewLogVolume"`
	CanConfigureLogRetentionPeriod            bool   `json:"canConfigureLogRetentionPeriod"`
	CanConfigureSubtraces                     bool   `json:"canConfigureSubtraces"`
	CanInvokeAlertChannel                     bool   `json:"canInvokeAlertChannel"`
	CanConfigureLlm                           bool   `json:"canConfigureLLM"`
	CanConfigureAiAgents                      bool   `json:"canConfigureAiAgents"`
	CanConfigureApdex                         bool   `json:"canConfigureApdex"`
	CanConfigureServiceLevelCorrectionWindows bool   `json:"canConfigureServiceLevelCorrectionWindows"`
	CanConfigureServiceLevelSmartAlerts       bool   `json:"canConfigureServiceLevelSmartAlerts"`
	CanConfigureServiceLevels                 bool   `json:"canConfigureServiceLevels"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (r *APIToken) GetIDForResourcePath() string {
	return r.InternalID
}
