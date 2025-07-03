package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

//nolint:gosec
const resourceAPITokenDefinitionTemplate = `
resource "instana_api_token" "example" {
  name = "name %d"
  can_configure_service_mapping = true
  can_configure_eum_applications = true
  can_configure_mobile_app_monitoring = true
  can_configure_users = true
  can_install_new_agents = true
  can_configure_integrations = true
  can_configure_events_and_alerts = true
  can_configure_maintenance_windows = true
  can_configure_application_smart_alerts = true
  can_configure_website_smart_alerts = true
  can_configure_mobile_app_smart_alerts = true
  can_configure_api_tokens = true
  can_configure_agent_run_mode = true
  can_view_audit_log = true
  can_configure_agents = true
  can_configure_authentication_methods = true
  can_configure_applications = true
  can_configure_teams = true
  can_configure_releases = true
  can_configure_log_management = true
  can_create_public_custom_dashboards = true
  can_view_logs = true
  can_view_trace_details = true
  can_configure_session_settings = true
  can_configure_service_level_indicators = true
  can_configure_global_alert_payload = true
  can_configure_global_application_smart_alerts = true
  can_configure_global_synthetic_smart_alerts = true
  can_configure_global_infra_smart_alerts = true
  can_configure_global_log_smart_alerts = true
  can_view_account_and_billing_information = true
  can_edit_all_accessible_custom_dashboards = true
  
  # Scope limitations
  limited_applications_scope = true
  limited_biz_ops_scope = true
  limited_websites_scope = true
  limited_kubernetes_scope = true
  limited_mobile_apps_scope = true
  limited_infrastructure_scope = true
  limited_synthetics_scope = true
  limited_vsphere_scope = true
  limited_phmc_scope = true
  limited_pvc_scope = true
  limited_zhmc_scope = true
  limited_pcf_scope = true
  limited_openstack_scope = true
  limited_automation_scope = true
  limited_logs_scope = true
  limited_nutanix_scope = true
  limited_xen_server_scope = true
  limited_windows_hypervisor_scope = true
  limited_alert_channels_scope = true
  limited_linux_kvm_hypervisor_scope = true
  
  # Additional permissions
  can_configure_personal_api_tokens = true
  can_configure_database_management = true
  can_configure_automation_actions = true
  can_configure_automation_policies = true
  can_run_automation_actions = true
  can_delete_automation_action_history = true
  can_configure_synthetic_tests = true
  can_configure_synthetic_locations = true
  can_configure_synthetic_credentials = true
  can_view_synthetic_tests = true
  can_view_synthetic_locations = true
  can_view_synthetic_test_results = true
  can_use_synthetic_credentials = true
  can_configure_bizops = true
  can_view_business_processes = true
  can_view_business_process_details = true
  can_view_business_activities = true
  can_view_biz_alerts = true
  can_delete_logs = true
  can_create_heap_dump = true
  can_create_thread_dump = true
  can_manually_close_issue = true
  can_view_log_volume = true
  can_configure_log_retention_period = true
  can_configure_subtraces = true
  can_invoke_alert_channel = true
  can_configure_llm = true
}
`

//nolint:gosec
const (
	apiTokenApiPath             = restapi.APITokensResourcePath + "/{internal-id}"
	testAPITokenDefinition      = "instana_api_token.example"
	apiTokenID                  = "api-token-id"
	apiTokenNameFieldValue      = resourceName
	apiTokenAccessGrantingToken = "api-token-access-granting-token"
	apiTokenInternalID          = "api-token-internal-id"
)

var apiTokenPermissionFields = []string{
	APITokenFieldCanConfigureServiceMapping,
	APITokenFieldCanConfigureEumApplications,
	APITokenFieldCanConfigureMobileAppMonitoring,
	APITokenFieldCanConfigureUsers,
	APITokenFieldCanInstallNewAgents,
	APITokenFieldCanConfigureIntegrations,
	APITokenFieldCanConfigureEventsAndAlerts,
	APITokenFieldCanConfigureMaintenanceWindows,
	APITokenFieldCanConfigureApplicationSmartAlerts,
	APITokenFieldCanConfigureWebsiteSmartAlerts,
	APITokenFieldCanConfigureMobileAppSmartAlerts,
	APITokenFieldCanConfigureAPITokens,
	APITokenFieldCanConfigureAgentRunMode,
	APITokenFieldCanViewAuditLog,
	APITokenFieldCanConfigureAgents,
	APITokenFieldCanConfigureAuthenticationMethods,
	APITokenFieldCanConfigureApplications,
	APITokenFieldCanConfigureTeams,
	APITokenFieldCanConfigureReleases,
	APITokenFieldCanConfigureLogManagement,
	APITokenFieldCanCreatePublicCustomDashboards,
	APITokenFieldCanViewLogs,
	APITokenFieldCanViewTraceDetails,
	APITokenFieldCanConfigureSessionSettings,
	APITokenFieldCanConfigureServiceLevelIndicators,
	APITokenFieldCanConfigureGlobalAlertPayload,
	APITokenFieldCanConfigureGlobalApplicationSmartAlerts,
	APITokenFieldCanConfigureGlobalSyntheticSmartAlerts,
	APITokenFieldCanConfigureGlobalInfraSmartAlerts,
	APITokenFieldCanConfigureGlobalLogSmartAlerts,
	APITokenFieldCanViewAccountAndBillingInformation,
	APITokenFieldCanEditAllAccessibleCustomDashboards,
	
	// Scope limitations
	APITokenFieldLimitedApplicationsScope,
	APITokenFieldLimitedBizOpsScope,
	APITokenFieldLimitedWebsitesScope,
	APITokenFieldLimitedKubernetesScope,
	APITokenFieldLimitedMobileAppsScope,
	APITokenFieldLimitedInfrastructureScope,
	APITokenFieldLimitedSyntheticsScope,
	APITokenFieldLimitedVsphereScope,
	APITokenFieldLimitedPhmcScope,
	APITokenFieldLimitedPvcScope,
	APITokenFieldLimitedZhmcScope,
	APITokenFieldLimitedPcfScope,
	APITokenFieldLimitedOpenstackScope,
	APITokenFieldLimitedAutomationScope,
	APITokenFieldLimitedLogsScope,
	APITokenFieldLimitedNutanixScope,
	APITokenFieldLimitedXenServerScope,
	APITokenFieldLimitedWindowsHypervisorScope,
	APITokenFieldLimitedAlertChannelsScope,
	APITokenFieldLimitedLinuxKvmHypervisorScope,
	
	// Additional permissions
	APITokenFieldCanConfigurePersonalAPITokens,
	APITokenFieldCanConfigureDatabaseManagement,
	APITokenFieldCanConfigureAutomationActions,
	APITokenFieldCanConfigureAutomationPolicies,
	APITokenFieldCanRunAutomationActions,
	APITokenFieldCanDeleteAutomationActionHistory,
	APITokenFieldCanConfigureSyntheticTests,
	APITokenFieldCanConfigureSyntheticLocations,
	APITokenFieldCanConfigureSyntheticCredentials,
	APITokenFieldCanViewSyntheticTests,
	APITokenFieldCanViewSyntheticLocations,
	APITokenFieldCanViewSyntheticTestResults,
	APITokenFieldCanUseSyntheticCredentials,
	APITokenFieldCanConfigureBizops,
	APITokenFieldCanViewBusinessProcesses,
	APITokenFieldCanViewBusinessProcessDetails,
	APITokenFieldCanViewBusinessActivities,
	APITokenFieldCanViewBizAlerts,
	APITokenFieldCanDeleteLogs,
	APITokenFieldCanCreateHeapDump,
	APITokenFieldCanCreateThreadDump,
	APITokenFieldCanManuallyCloseIssue,
	APITokenFieldCanViewLogVolume,
	APITokenFieldCanConfigureLogRetentionPeriod,
	APITokenFieldCanConfigureSubtraces,
	APITokenFieldCanInvokeAlertChannel,
	APITokenFieldCanConfigureLlm,
}

func TestCRUDOfAPITokenResourceWithMockServer(t *testing.T) {
	id := RandomID()
	accessGrantingToken := RandomID()
	internalID := RandomID()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, restapi.APITokensResourcePath, func(w http.ResponseWriter, r *http.Request) {
		apiToken := &restapi.APIToken{}
		err := json.NewDecoder(r.Body).Decode(apiToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = r.Write(bytes.NewBufferString("Failed to get request"))
			if err != nil {
				fmt.Printf("failed to write response; %s\n", err)
			}
		} else {
			apiToken.ID = id
			apiToken.AccessGrantingToken = accessGrantingToken
			apiToken.InternalID = internalID
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(apiToken)
			if err != nil {
				fmt.Printf("failed to decode json; %s\n", err)
			}
		}
	})
	httpServer.AddRoute(http.MethodPut, apiTokenApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, apiTokenApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, apiTokenApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		modCount := httpServer.GetCallCount(http.MethodPut, restapi.APITokensResourcePath+"/"+internalID)
		jsonData := fmt.Sprintf(`
		{
			"id" : "%s",
			"accessGrantingToken": "%s",
			"internalId" : "%s",
			"name" : "name %d",
			"canConfigureServiceMapping" : true,
			"canConfigureEumApplications" : true,
			"canConfigureMobileAppMonitoring" : true,
			"canConfigureUsers" : true,
			"canInstallNewAgents" : true,
			"canConfigureIntegrations" : true,
			"canConfigureEventsAndAlerts" : true,
			"canConfigureMaintenanceWindows" : true,
			"canConfigureApplicationSmartAlerts" : true,
			"canConfigureWebsiteSmartAlerts" : true,
			"canConfigureMobileAppSmartAlerts" : true,
			"canConfigureApiTokens" : true,
			"canConfigureAgentRunMode" : true,
			"canViewAuditLog" : true,
			"canConfigureAgents" : true,
			"canConfigureAuthenticationMethods" : true,
			"canConfigureApplications" : true,
			"canConfigureTeams" : true,
			"canConfigureReleases" : true,
			"canConfigureLogManagement" : true,
			"canCreatePublicCustomDashboards" : true,
			"canViewLogs" : true,
			"canViewTraceDetails" : true,
			"canConfigureSessionSettings" : true,
			"canConfigureServiceLevelIndicators" : true,
			"canConfigureGlobalAlertPayload" : true,
			"canConfigureGlobalApplicationSmartAlerts" : true,
			"canConfigureGlobalSyntheticSmartAlerts" : true,
			"canConfigureGlobalInfraSmartAlerts" : true,
			"canConfigureGlobalLogSmartAlerts" : true,
			"canViewAccountAndBillingInformation" : true,
			"canEditAllAccessibleCustomDashboards" : true,
			
			"limitedApplicationsScope" : true,
			"limitedBizOpsScope" : true,
			"limitedWebsitesScope" : true,
			"limitedKubernetesScope" : true,
			"limitedMobileAppsScope" : true,
			"limitedInfrastructureScope" : true,
			"limitedSyntheticsScope" : true,
			"limitedVsphereScope" : true,
			"limitedPhmcScope" : true,
			"limitedPvcScope" : true,
			"limitedZhmcScope" : true,
			"limitedPcfScope" : true,
			"limitedOpenstackScope" : true,
			"limitedAutomationScope" : true,
			"limitedLogsScope" : true,
			"limitedNutanixScope" : true,
			"limitedXenServerScope" : true,
			"limitedWindowsHypervisorScope" : true,
			"limitedAlertChannelsScope" : true,
			"limitedLinuxKvmHypervisorScope" : true,
			
			"canConfigurePersonalApiTokens" : true,
			"canConfigureDatabaseManagement" : true,
			"canConfigureAutomationActions" : true,
			"canConfigureAutomationPolicies" : true,
			"canRunAutomationActions" : true,
			"canDeleteAutomationActionHistory" : true,
			"canConfigureSyntheticTests" : true,
			"canConfigureSyntheticLocations" : true,
			"canConfigureSyntheticCredentials" : true,
			"canViewSyntheticTests" : true,
			"canViewSyntheticLocations" : true,
			"canViewSyntheticTestResults" : true,
			"canUseSyntheticCredentials" : true,
			"canConfigureBizops" : true,
			"canViewBusinessProcesses" : true,
			"canViewBusinessProcessDetails" : true,
			"canViewBusinessActivities" : true,
			"canViewBizAlerts" : true,
			"canDeleteLogs" : true,
			"canCreateHeapDump" : true,
			"canCreateThreadDump" : true,
			"canManuallyCloseIssue" : true,
			"canViewLogVolume" : true,
			"canConfigureLogRetentionPeriod" : true,
			"canConfigureSubtraces" : true,
			"canInvokeAlertChannel" : true,
			"canConfigureLlm" : true
		}
		`, id, accessGrantingToken, vars["internal-id"], modCount)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(jsonData))
		if err != nil {
			fmt.Printf("failed to write response; %s\n", err)
		}
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAPITokenConfigResourceTestStep(httpServer.GetPort(), 0, id, accessGrantingToken, internalID),
			testStepImportWithCustomID(testAPITokenDefinition, internalID),
			createAPITokenConfigResourceTestStep(httpServer.GetPort(), 1, id, accessGrantingToken, internalID),
			testStepImportWithCustomID(testAPITokenDefinition, internalID),
		},
	})
}

func createAPITokenConfigResourceTestStep(httpPort int, iteration int, id string, accessGrantingToken string, internalID string) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(resourceAPITokenDefinitionTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(testAPITokenDefinition, "id", id),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldAccessGrantingToken, accessGrantingToken),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldInternalID, internalID),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureServiceMapping, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureEumApplications, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureMobileAppMonitoring, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureUsers, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanInstallNewAgents, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureIntegrations, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureEventsAndAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureMaintenanceWindows, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureApplicationSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureWebsiteSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureMobileAppSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAPITokens, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAgentRunMode, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewAuditLog, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAgents, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAuthenticationMethods, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureApplications, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureTeams, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureReleases, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureLogManagement, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanCreatePublicCustomDashboards, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewLogs, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewTraceDetails, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureSessionSettings, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureServiceLevelIndicators, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalAlertPayload, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalApplicationSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalSyntheticSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalInfraSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalLogSmartAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewAccountAndBillingInformation, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanEditAllAccessibleCustomDashboards, trueAsString),
			
			// Scope limitations
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedApplicationsScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedBizOpsScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedWebsitesScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedKubernetesScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedMobileAppsScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedInfrastructureScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedSyntheticsScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedVsphereScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedPhmcScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedPvcScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedZhmcScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedPcfScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedOpenstackScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedAutomationScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedLogsScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedNutanixScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedXenServerScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedWindowsHypervisorScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedAlertChannelsScope, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldLimitedLinuxKvmHypervisorScope, trueAsString),
			
			// Additional permissions
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigurePersonalAPITokens, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureDatabaseManagement, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAutomationActions, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAutomationPolicies, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanRunAutomationActions, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanDeleteAutomationActionHistory, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureSyntheticTests, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureSyntheticLocations, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureSyntheticCredentials, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewSyntheticTests, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewSyntheticLocations, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewSyntheticTestResults, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanUseSyntheticCredentials, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureBizops, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewBusinessProcesses, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewBusinessProcessDetails, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewBusinessActivities, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewBizAlerts, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanDeleteLogs, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanCreateHeapDump, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanCreateThreadDump, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanManuallyCloseIssue, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewLogVolume, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureLogRetentionPeriod, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureSubtraces, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanInvokeAlertChannel, trueAsString),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureLlm, trueAsString),
		),
	}
}

func TestAPITokenSchemaDefinitionIsValid(t *testing.T) {
	schema := NewAPITokenResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(APITokenFieldAccessGrantingToken)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(APITokenFieldInternalID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(APITokenFieldName)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureServiceMapping, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureEumApplications, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureMobileAppMonitoring, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureUsers, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanInstallNewAgents, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureIntegrations, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureEventsAndAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureMaintenanceWindows, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureApplicationSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureWebsiteSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureMobileAppSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureAPITokens, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureAgentRunMode, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewAuditLog, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureAgents, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureAuthenticationMethods, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureApplications, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureTeams, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureReleases, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureLogManagement, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanCreatePublicCustomDashboards, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewLogs, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewTraceDetails, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureSessionSettings, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureServiceLevelIndicators, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureGlobalAlertPayload, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureGlobalApplicationSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureGlobalSyntheticSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureGlobalInfraSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureGlobalLogSmartAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewAccountAndBillingInformation, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanEditAllAccessibleCustomDashboards, false)
	
	// Scope limitations
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedApplicationsScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedBizOpsScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedWebsitesScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedKubernetesScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedMobileAppsScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedInfrastructureScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedSyntheticsScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedVsphereScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedPhmcScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedPvcScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedZhmcScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedPcfScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedOpenstackScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedAutomationScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedLogsScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedNutanixScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedXenServerScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedWindowsHypervisorScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedAlertChannelsScope, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldLimitedLinuxKvmHypervisorScope, false)
	
	// Additional permissions
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigurePersonalAPITokens, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureDatabaseManagement, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureAutomationActions, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureAutomationPolicies, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanRunAutomationActions, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanDeleteAutomationActionHistory, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureSyntheticTests, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureSyntheticLocations, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureSyntheticCredentials, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewSyntheticTests, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewSyntheticLocations, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewSyntheticTestResults, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanUseSyntheticCredentials, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureBizops, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewBusinessProcesses, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewBusinessProcessDetails, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewBusinessActivities, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewBizAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanDeleteLogs, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanCreateHeapDump, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanCreateThreadDump, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanManuallyCloseIssue, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanViewLogVolume, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureLogRetentionPeriod, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureSubtraces, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanInvokeAlertChannel, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(APITokenFieldCanConfigureLlm, false)
}

func TestAPITokenResourceShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 2, NewAPITokenResourceHandle().MetaData().SchemaVersion)
}

func TestAPITokenResourceShouldHaveOneStateMigrators(t *testing.T) {
	require.Equal(t, 1, len(NewAPITokenResourceHandle().StateUpgraders()))
}

func TestAPITokenResourceShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAPITokenResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, APITokenFieldFullName)
	require.Contains(t, result, APITokenFieldName)
	require.Equal(t, "test", result[APITokenFieldName])
}

func TestAPITokenResourceShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAPITokenResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForUserroleResource(t *testing.T) {
	name := NewAPITokenResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_api_token")
}

func TestShouldSetCalculateAccessGrantingTokenAndInternal(t *testing.T) {
	testHelper := NewTestHelper[*restapi.APIToken](t)
	sut := NewAPITokenResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.SetComputedFields(resourceData)

	require.NoError(t, err)
	require.NotEmpty(t, resourceData.Get(APITokenFieldInternalID))
	require.NotEmpty(t, resourceData.Get(APITokenFieldAccessGrantingToken))
}

func TestShouldUpdateBasicFieldsOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	testHelper := NewTestHelper[*restapi.APIToken](t)
	sut := NewAPITokenResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		InternalID:          apiTokenInternalID,
	}

	err := sut.UpdateState(resourceData, &apiToken)

	require.Nil(t, err)
	require.Equal(t, apiTokenID, resourceData.Id())
	require.Equal(t, apiTokenAccessGrantingToken, resourceData.Get(APITokenFieldAccessGrantingToken))
	require.Equal(t, apiTokenInternalID, resourceData.Get(APITokenFieldInternalID))
	require.Equal(t, apiTokenNameFieldValue, resourceData.Get(APITokenFieldName))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureServiceMapping).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureEumApplications).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureMobileAppMonitoring).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureUsers).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanInstallNewAgents).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureIntegrations).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureEventsAndAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureMaintenanceWindows).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureApplicationSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureWebsiteSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureMobileAppSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureAPITokens).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureAgentRunMode).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanViewAuditLog).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureAgents).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureAuthenticationMethods).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureApplications).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureTeams).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureReleases).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureLogManagement).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanCreatePublicCustomDashboards).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanViewLogs).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanViewTraceDetails).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureSessionSettings).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureServiceLevelIndicators).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureGlobalAlertPayload).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureGlobalApplicationSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureGlobalSyntheticSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureGlobalInfraSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanConfigureGlobalLogSmartAlerts).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanViewAccountAndBillingInformation).(bool))
	require.False(t, resourceData.Get(APITokenFieldCanEditAllAccessibleCustomDashboards).(bool))
}

func TestShouldUpdateCanConfigureServiceMappingPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                         apiTokenID,
		InternalID:                 apiTokenInternalID,
		AccessGrantingToken:        apiTokenAccessGrantingToken,
		Name:                       apiTokenNameFieldValue,
		CanConfigureServiceMapping: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureServiceMapping)
}

func TestShouldUpdateCanConfigureEumApplicationsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                          apiTokenID,
		InternalID:                  apiTokenInternalID,
		AccessGrantingToken:         apiTokenAccessGrantingToken,
		Name:                        apiTokenNameFieldValue,
		CanConfigureEumApplications: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureEumApplications)
}

func TestShouldUpdateCanConfigureMobileAppMonitoringPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                              apiTokenID,
		InternalID:                      apiTokenInternalID,
		AccessGrantingToken:             apiTokenAccessGrantingToken,
		Name:                            apiTokenNameFieldValue,
		CanConfigureMobileAppMonitoring: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureMobileAppMonitoring)
}

func TestShouldUpdateCanConfigureUsersPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanConfigureUsers:   true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureUsers)
}

func TestShouldUpdateCanInstallNewAgentsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanInstallNewAgents: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanInstallNewAgents)
}

func TestShouldUpdateCanConfigureIntegrationsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                       apiTokenID,
		InternalID:               apiTokenInternalID,
		AccessGrantingToken:      apiTokenAccessGrantingToken,
		Name:                     apiTokenNameFieldValue,
		CanConfigureIntegrations: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureIntegrations)
}

func TestShouldUpdateCanConfigureEventsAndAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                          apiTokenID,
		InternalID:                  apiTokenInternalID,
		AccessGrantingToken:         apiTokenAccessGrantingToken,
		Name:                        apiTokenNameFieldValue,
		CanConfigureEventsAndAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureEventsAndAlerts)
}

func TestShouldUpdateCanConfigureMaintenanceWindowsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                             apiTokenID,
		InternalID:                     apiTokenInternalID,
		AccessGrantingToken:            apiTokenAccessGrantingToken,
		Name:                           apiTokenNameFieldValue,
		CanConfigureMaintenanceWindows: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureMaintenanceWindows)
}

func TestShouldUpdateCanConfigureApplicationSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                 apiTokenID,
		InternalID:                         apiTokenInternalID,
		AccessGrantingToken:                apiTokenAccessGrantingToken,
		Name:                               apiTokenNameFieldValue,
		CanConfigureApplicationSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureApplicationSmartAlerts)
}

func TestShouldUpdateCanConfigureWebsiteSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                             apiTokenID,
		InternalID:                     apiTokenInternalID,
		AccessGrantingToken:            apiTokenAccessGrantingToken,
		Name:                           apiTokenNameFieldValue,
		CanConfigureWebsiteSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureWebsiteSmartAlerts)
}

func TestShouldUpdateCanConfigureMobileAppSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                               apiTokenID,
		InternalID:                       apiTokenInternalID,
		AccessGrantingToken:              apiTokenAccessGrantingToken,
		Name:                             apiTokenNameFieldValue,
		CanConfigureMobileAppSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureMobileAppSmartAlerts)
}

func TestShouldUpdateCanConfigureAPITokensPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                    apiTokenID,
		InternalID:            apiTokenInternalID,
		AccessGrantingToken:   apiTokenAccessGrantingToken,
		Name:                  apiTokenNameFieldValue,
		CanConfigureAPITokens: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureAPITokens)
}

func TestShouldUpdateCanConfigureAgentRunModePermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                       apiTokenID,
		InternalID:               apiTokenInternalID,
		AccessGrantingToken:      apiTokenAccessGrantingToken,
		Name:                     apiTokenNameFieldValue,
		CanConfigureAgentRunMode: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureAgentRunMode)
}

func TestShouldUpdateCanViewAuditLogPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanViewAuditLog:     true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanViewAuditLog)
}

func TestShouldUpdateCanConfigureAgentsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanConfigureAgents:  true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureAgents)
}

func TestShouldUpdateCanConfigureAuthenticationMethodsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                apiTokenID,
		InternalID:                        apiTokenInternalID,
		AccessGrantingToken:               apiTokenAccessGrantingToken,
		Name:                              apiTokenNameFieldValue,
		CanConfigureAuthenticationMethods: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureAuthenticationMethods)
}

func TestShouldUpdateCanConfigureApplicationsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                       apiTokenID,
		InternalID:               apiTokenInternalID,
		AccessGrantingToken:      apiTokenAccessGrantingToken,
		Name:                     apiTokenNameFieldValue,
		CanConfigureApplications: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureApplications)
}

func TestShouldUpdateCanConfigureTeamsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanConfigureTeams:   true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureTeams)
}

func TestShouldUpdateCanConfigureReleasesPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                   apiTokenID,
		InternalID:           apiTokenInternalID,
		AccessGrantingToken:  apiTokenAccessGrantingToken,
		Name:                 apiTokenNameFieldValue,
		CanConfigureReleases: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureReleases)
}

func TestShouldUpdateCanConfigureLogManagementPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                        apiTokenID,
		InternalID:                apiTokenInternalID,
		AccessGrantingToken:       apiTokenAccessGrantingToken,
		Name:                      apiTokenNameFieldValue,
		CanConfigureLogManagement: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureLogManagement)
}

func TestShouldUpdateCanCreatePublicCustomDashboardsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                              apiTokenID,
		InternalID:                      apiTokenInternalID,
		AccessGrantingToken:             apiTokenAccessGrantingToken,
		Name:                            apiTokenNameFieldValue,
		CanCreatePublicCustomDashboards: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanCreatePublicCustomDashboards)
}

func TestShouldUpdateCanViewLogsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanViewLogs:         true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanViewLogs)
}

func TestShouldUpdateCanViewTraceDetailsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                  apiTokenID,
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenNameFieldValue,
		CanViewTraceDetails: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanViewTraceDetails)
}

func TestShouldUpdateCanConfigureSessionSettingsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                          apiTokenID,
		InternalID:                  apiTokenInternalID,
		AccessGrantingToken:         apiTokenAccessGrantingToken,
		Name:                        apiTokenNameFieldValue,
		CanConfigureSessionSettings: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureSessionSettings)
}

func TestShouldUpdateCanConfigureServiceLevelIndicatorsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                 apiTokenID,
		InternalID:                         apiTokenInternalID,
		AccessGrantingToken:                apiTokenAccessGrantingToken,
		Name:                               apiTokenNameFieldValue,
		CanConfigureServiceLevelIndicators: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureServiceLevelIndicators)
}

func TestShouldUpdateCanConfigureGlobalAlertPayloadPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                             apiTokenID,
		Name:                           apiTokenNameFieldValue,
		CanConfigureGlobalAlertPayload: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureGlobalAlertPayload)
}

func TestShouldUpdateCanConfigureGlobalApplicationSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                       apiTokenID,
		InternalID:                               apiTokenInternalID,
		AccessGrantingToken:                      apiTokenAccessGrantingToken,
		Name:                                     apiTokenNameFieldValue,
		CanConfigureGlobalApplicationSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureGlobalApplicationSmartAlerts)
}

func TestShouldUpdateCanConfigureGlobalSyntheticSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                     apiTokenID,
		InternalID:                             apiTokenInternalID,
		AccessGrantingToken:                    apiTokenAccessGrantingToken,
		Name:                                   apiTokenNameFieldValue,
		CanConfigureGlobalSyntheticSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureGlobalSyntheticSmartAlerts)
}

func TestShouldUpdateCanConfigureGlobalInfraSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                 apiTokenID,
		InternalID:                         apiTokenInternalID,
		AccessGrantingToken:                apiTokenAccessGrantingToken,
		Name:                               apiTokenNameFieldValue,
		CanConfigureGlobalInfraSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureGlobalInfraSmartAlerts)
}

func TestShouldUpdateCanConfigureGlobalLogSmartAlertsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                               apiTokenID,
		InternalID:                       apiTokenInternalID,
		AccessGrantingToken:              apiTokenAccessGrantingToken,
		Name:                             apiTokenNameFieldValue,
		CanConfigureGlobalLogSmartAlerts: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanConfigureGlobalLogSmartAlerts)
}

func TestShouldUpdateCanViewAccountAndBillingInformationPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                  apiTokenID,
		InternalID:                          apiTokenInternalID,
		AccessGrantingToken:                 apiTokenAccessGrantingToken,
		Name:                                apiTokenNameFieldValue,
		CanViewAccountAndBillingInformation: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanViewAccountAndBillingInformation)
}

func TestShouldUpdateCanEditAllAccessibleCustomDashboardsPermissionOfTerraformResourceStateFromModelForAPIToken(t *testing.T) {
	apiToken := restapi.APIToken{
		ID:                                   apiTokenID,
		InternalID:                           apiTokenInternalID,
		AccessGrantingToken:                  apiTokenAccessGrantingToken,
		Name:                                 apiTokenNameFieldValue,
		CanEditAllAccessibleCustomDashboards: true,
	}

	testSingleAPITokenPermissionSet(t, apiToken, APITokenFieldCanEditAllAccessibleCustomDashboards)
}

func testSingleAPITokenPermissionSet(t *testing.T, apiToken restapi.APIToken, expectedPermissionField string) {
	testHelper := NewTestHelper[*restapi.APIToken](t)
	sut := NewAPITokenResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &apiToken)

	require.Nil(t, err)
	require.True(t, resourceData.Get(expectedPermissionField).(bool))
	for _, permissionField := range apiTokenPermissionFields {
		if permissionField != expectedPermissionField {
			require.False(t, resourceData.Get(permissionField).(bool))
		}
	}
}

func TestShouldConvertStateOfAPITokenTerraformResourceToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.APIToken](t)
	resourceHandle := NewAPITokenResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(apiTokenID)
	setValueOnResourceData(t, resourceData, APITokenFieldAccessGrantingToken, apiTokenAccessGrantingToken)
	setValueOnResourceData(t, resourceData, APITokenFieldInternalID, apiTokenInternalID)
	setValueOnResourceData(t, resourceData, APITokenFieldName, apiTokenNameFieldValue)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureServiceMapping, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureEumApplications, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureMobileAppMonitoring, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureUsers, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanInstallNewAgents, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureIntegrations, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureEventsAndAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureMaintenanceWindows, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureApplicationSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureWebsiteSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureMobileAppSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureAPITokens, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureAgentRunMode, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanViewAuditLog, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureAgents, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureAuthenticationMethods, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureApplications, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureTeams, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureReleases, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureLogManagement, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanCreatePublicCustomDashboards, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanViewLogs, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanViewTraceDetails, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureSessionSettings, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureServiceLevelIndicators, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureGlobalAlertPayload, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureGlobalApplicationSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureGlobalSyntheticSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureGlobalInfraSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanConfigureGlobalLogSmartAlerts, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanViewAccountAndBillingInformation, true)
	setValueOnResourceData(t, resourceData, APITokenFieldCanEditAllAccessibleCustomDashboards, true)

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.APIToken{}, model, "Model should be an alerting channel")
	require.Equal(t, apiTokenID, model.ID)
	require.Equal(t, apiTokenAccessGrantingToken, model.AccessGrantingToken)
	require.Equal(t, apiTokenInternalID, model.GetIDForResourcePath())
	require.Equal(t, apiTokenInternalID, model.InternalID)
	require.Equal(t, apiTokenNameFieldValue, model.Name)
	require.True(t, model.CanConfigureServiceMapping)
	require.True(t, model.CanConfigureEumApplications)
	require.True(t, model.CanConfigureMobileAppMonitoring)
	require.True(t, model.CanConfigureUsers)
	require.True(t, model.CanInstallNewAgents)
	require.True(t, model.CanConfigureIntegrations)
	require.True(t, model.CanConfigureEventsAndAlerts)
	require.True(t, model.CanConfigureMaintenanceWindows)
	require.True(t, model.CanConfigureApplicationSmartAlerts)
	require.True(t, model.CanConfigureWebsiteSmartAlerts)
	require.True(t, model.CanConfigureMobileAppSmartAlerts)
	require.True(t, model.CanConfigureAPITokens)
	require.True(t, model.CanConfigureAgentRunMode)
	require.True(t, model.CanViewAuditLog)
	require.True(t, model.CanConfigureAgents)
	require.True(t, model.CanConfigureAuthenticationMethods)
	require.True(t, model.CanConfigureApplications)
	require.True(t, model.CanConfigureTeams)
	require.True(t, model.CanConfigureReleases)
	require.True(t, model.CanConfigureLogManagement)
	require.True(t, model.CanCreatePublicCustomDashboards)
	require.True(t, model.CanViewLogs)
	require.True(t, model.CanViewTraceDetails)
	require.True(t, model.CanConfigureSessionSettings)
	require.True(t, model.CanConfigureServiceLevelIndicators)
	require.True(t, model.CanConfigureGlobalAlertPayload)
	require.True(t, model.CanConfigureGlobalApplicationSmartAlerts)
	require.True(t, model.CanConfigureGlobalSyntheticSmartAlerts)
	require.True(t, model.CanConfigureGlobalInfraSmartAlerts)
	require.True(t, model.CanConfigureGlobalLogSmartAlerts)
	require.True(t, model.CanViewAccountAndBillingInformation)
	require.True(t, model.CanEditAllAccessibleCustomDashboards)
}
