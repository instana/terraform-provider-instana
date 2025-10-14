package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//nolint:gosec
const (
	// ResourceInstanaAPIToken the name of the terraform-provider-instana resource to manage API tokens
	ResourceInstanaAPIToken = "instana_api_token"

	//APITokenFieldAccessGrantingToken constant value for the schema field access_granting_token
	APITokenFieldAccessGrantingToken = "access_granting_token"
	//APITokenFieldInternalID constant value for the schema field internal_id
	APITokenFieldInternalID = "internal_id"
	//APITokenFieldName constant value for the schema field name
	APITokenFieldName = "name"
	//APITokenFieldFullName constant value for the schema field full_name
	APITokenFieldFullName = "full_name"
	//APITokenFieldCanConfigureServiceMapping constant value for the schema field can_configure_service_mapping
	APITokenFieldCanConfigureServiceMapping = "can_configure_service_mapping"
	//APITokenFieldCanConfigureEumApplications constant value for the schema field can_configure_eum_applications
	APITokenFieldCanConfigureEumApplications = "can_configure_eum_applications"
	//APITokenFieldCanConfigureMobileAppMonitoring constant value for the schema field can_configure_mobile_app_monitoring
	APITokenFieldCanConfigureMobileAppMonitoring = "can_configure_mobile_app_monitoring"
	//APITokenFieldCanConfigureUsers constant value for the schema field can_configure_users
	APITokenFieldCanConfigureUsers = "can_configure_users"
	//APITokenFieldCanInstallNewAgents constant value for the schema field can_install_new_agents
	APITokenFieldCanInstallNewAgents = "can_install_new_agents"
	//APITokenFieldCanConfigureIntegrations constant value for the schema field can_configure_integrations
	APITokenFieldCanConfigureIntegrations           = "can_configure_integrations"
	APITokenFieldCanConfigureEventsAndAlerts        = "can_configure_events_and_alerts"
	APITokenFieldCanConfigureMaintenanceWindows     = "can_configure_maintenance_windows"
	APITokenFieldCanConfigureApplicationSmartAlerts = "can_configure_application_smart_alerts"
	APITokenFieldCanConfigureWebsiteSmartAlerts     = "can_configure_website_smart_alerts"
	APITokenFieldCanConfigureMobileAppSmartAlerts   = "can_configure_mobile_app_smart_alerts"
	//APITokenFieldCanConfigureAPITokens constant value for the schema field can_configure_api_tokens
	APITokenFieldCanConfigureAPITokens = "can_configure_api_tokens"
	//APITokenFieldCanConfigureAgentRunMode constant value for the schema field can_configure_agent_run_mode
	APITokenFieldCanConfigureAgentRunMode = "can_configure_agent_run_mode"
	//APITokenFieldCanViewAuditLog constant value for the schema field can_view_audit_log
	APITokenFieldCanViewAuditLog = "can_view_audit_log"
	//APITokenFieldCanConfigureAgents constant value for the schema field can_configure_agents
	APITokenFieldCanConfigureAgents = "can_configure_agents"
	//APITokenFieldCanConfigureAuthenticationMethods constant value for the schema field can_configure_authentication_methods
	APITokenFieldCanConfigureAuthenticationMethods = "can_configure_authentication_methods"
	//APITokenFieldCanConfigureApplications constant value for the schema field can_configure_applications
	APITokenFieldCanConfigureApplications = "can_configure_applications"
	//APITokenFieldCanConfigureTeams constant value for the schema field can_configure_teams
	APITokenFieldCanConfigureTeams = "can_configure_teams"
	//APITokenFieldCanConfigureReleases constant value for the schema field can_configure_releases
	APITokenFieldCanConfigureReleases = "can_configure_releases"
	//APITokenFieldCanConfigureLogManagement constant value for the schema field can_configure_log_management
	APITokenFieldCanConfigureLogManagement = "can_configure_log_management"
	//APITokenFieldCanCreatePublicCustomDashboards constant value for the schema field can_create_public_custom_dashboards
	APITokenFieldCanCreatePublicCustomDashboards = "can_create_public_custom_dashboards"
	//APITokenFieldCanViewLogs constant value for the schema field can_view_logs
	APITokenFieldCanViewLogs = "can_view_logs"
	//APITokenFieldCanViewTraceDetails constant value for the schema field can_view_trace_details
	APITokenFieldCanViewTraceDetails = "can_view_trace_details"
	//APITokenFieldCanConfigureSessionSettings constant value for the schema field can_configure_session_settings
	APITokenFieldCanConfigureSessionSettings = "can_configure_session_settings"
	//APITokenFieldCanConfigureGlobalAlertPayload constant value for the schema field can_configure_global_alert_payload
	APITokenFieldCanConfigureGlobalAlertPayload           = "can_configure_global_alert_payload"
	APITokenFieldCanConfigureGlobalApplicationSmartAlerts = "can_configure_global_application_smart_alerts"
	APITokenFieldCanConfigureGlobalSyntheticSmartAlerts   = "can_configure_global_synthetic_smart_alerts"
	APITokenFieldCanConfigureGlobalInfraSmartAlerts       = "can_configure_global_infra_smart_alerts"
	APITokenFieldCanConfigureGlobalLogSmartAlerts         = "can_configure_global_log_smart_alerts"
	//APITokenFieldCanViewAccountAndBillingInformation constant value for the schema field can_view_account_and_billing_information
	APITokenFieldCanViewAccountAndBillingInformation = "can_view_account_and_billing_information"
	//APITokenFieldCanEditAllAccessibleCustomDashboards constant value for the schema field can_edit_all_accessible_custom_dashboards
	APITokenFieldCanEditAllAccessibleCustomDashboards = "can_edit_all_accessible_custom_dashboards"

	// New permission fields
	APITokenFieldLimitedApplicationsScope       = "limited_applications_scope"
	APITokenFieldLimitedBizOpsScope             = "limited_biz_ops_scope"
	APITokenFieldLimitedWebsitesScope           = "limited_websites_scope"
	APITokenFieldLimitedKubernetesScope         = "limited_kubernetes_scope"
	APITokenFieldLimitedMobileAppsScope         = "limited_mobile_apps_scope"
	APITokenFieldLimitedInfrastructureScope     = "limited_infrastructure_scope"
	APITokenFieldLimitedSyntheticsScope         = "limited_synthetics_scope"
	APITokenFieldLimitedVsphereScope            = "limited_vsphere_scope"
	APITokenFieldLimitedPhmcScope               = "limited_phmc_scope"
	APITokenFieldLimitedPvcScope                = "limited_pvc_scope"
	APITokenFieldLimitedZhmcScope               = "limited_zhmc_scope"
	APITokenFieldLimitedPcfScope                = "limited_pcf_scope"
	APITokenFieldLimitedOpenstackScope          = "limited_openstack_scope"
	APITokenFieldLimitedAutomationScope         = "limited_automation_scope"
	APITokenFieldLimitedLogsScope               = "limited_logs_scope"
	APITokenFieldLimitedNutanixScope            = "limited_nutanix_scope"
	APITokenFieldLimitedXenServerScope          = "limited_xen_server_scope"
	APITokenFieldLimitedWindowsHypervisorScope  = "limited_windows_hypervisor_scope"
	APITokenFieldLimitedAlertChannelsScope      = "limited_alert_channels_scope"
	APITokenFieldLimitedLinuxKvmHypervisorScope = "limited_linux_kvm_hypervisor_scope"

	APITokenFieldLimitedServiceLevelScope = "limited_service_level_scope"
	APITokenFieldLimitedAiGatewayScope    = "limited_ai_gateway_scope"
	APITokenFieldLimitedGenAIScope        = "limited_gen_ai_scope"

	APITokenFieldCanConfigurePersonalAPITokens             = "can_configure_personal_api_tokens"
	APITokenFieldCanConfigureDatabaseManagement            = "can_configure_database_management"
	APITokenFieldCanConfigureAutomationActions             = "can_configure_automation_actions"
	APITokenFieldCanConfigureAutomationPolicies            = "can_configure_automation_policies"
	APITokenFieldCanRunAutomationActions                   = "can_run_automation_actions"
	APITokenFieldCanDeleteAutomationActionHistory          = "can_delete_automation_action_history"
	APITokenFieldCanConfigureSyntheticTests                = "can_configure_synthetic_tests"
	APITokenFieldCanConfigureSyntheticLocations            = "can_configure_synthetic_locations"
	APITokenFieldCanConfigureSyntheticCredentials          = "can_configure_synthetic_credentials"
	APITokenFieldCanViewSyntheticTests                     = "can_view_synthetic_tests"
	APITokenFieldCanViewSyntheticLocations                 = "can_view_synthetic_locations"
	APITokenFieldCanViewSyntheticTestResults               = "can_view_synthetic_test_results"
	APITokenFieldCanUseSyntheticCredentials                = "can_use_synthetic_credentials"
	APITokenFieldCanConfigureBizops                        = "can_configure_bizops"
	APITokenFieldCanViewBusinessProcesses                  = "can_view_business_processes"
	APITokenFieldCanViewBusinessProcessDetails             = "can_view_business_process_details"
	APITokenFieldCanViewBusinessActivities                 = "can_view_business_activities"
	APITokenFieldCanViewBizAlerts                          = "can_view_biz_alerts"
	APITokenFieldCanDeleteLogs                             = "can_delete_logs"
	APITokenFieldCanCreateHeapDump                         = "can_create_heap_dump"
	APITokenFieldCanCreateThreadDump                       = "can_create_thread_dump"
	APITokenFieldCanManuallyCloseIssue                     = "can_manually_close_issue"
	APITokenFieldCanViewLogVolume                          = "can_view_log_volume"
	APITokenFieldCanConfigureLogRetentionPeriod            = "can_configure_log_retention_period"
	APITokenFieldCanConfigureSubtraces                     = "can_configure_subtraces"
	APITokenFieldCanInvokeAlertChannel                     = "can_invoke_alert_channel"
	APITokenFieldCanConfigureLlm                           = "can_configure_llm"
	APITokenFieldCanConfigureAiAgents                      = "can_configure_ai_agents"
	APITokenFieldCanConfigureApdex                         = "can_configure_apdex"
	APITokenFieldCanConfigureServiceLevelCorrectionWindows = "can_configure_service_level_correction_windows"
	APITokenFieldCanConfigureServiceLevelSmartAlerts       = "can_configure_service_level_smart_alerts"
	APITokenFieldCanConfigureServiceLevels                 = "can_configure_service_levels"
)

var (
	apiTokenSchemaAccessGrantingToken = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The token used for the api Client used in the Authorization header to authenticate the client",
	}
	apiTokenSchemaInternalID = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The internal ID of the access token from the Instana platform",
	}
	apiTokenSchemaName = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the API token",
	}
	apiTokenSchemaFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the API token including prefix in suffix",
	}
	apiTokenSchemaCanConfigureServiceMapping = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure service mappings",
	}
	apiTokenSchemaCanConfigureEumApplications = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure End User Monitoring applications",
	}
	apiTokenSchemaCanConfigureMobileAppMonitoring = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Mobile App Monitoring",
	}
	apiTokenSchemaCanConfigureUsers = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure users",
	}
	apiTokenSchemaCanInstallNewAgents = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to install new agents",
	}
	apiTokenSchemaCanConfigureIntegrations = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure integrations",
	}
	apiTokenSchemaCanConfigureEventsAndAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Instana Events and Alerts",
	}
	apiTokenSchemaCanConfigureMaintenanceWindows = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Instana Maintenance Windows",
	}
	apiTokenSchemaCanConfigureApplicationSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Instana Application Smart Alerts",
	}
	apiTokenSchemaCanConfigureWebsiteSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Instana Website Smart Alerts",
	}
	apiTokenSchemaCanConfigureMobileAppSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Instana MobileApp Smart Alerts",
	}
	apiTokenSchemaCanConfigureAPITokens = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure API tokens",
	}
	apiTokenSchemaCanConfigureAgentRunMode = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure agent run mode",
	}
	apiTokenSchemaCanViewAuditLog = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view the audit log",
	}
	apiTokenSchemaCanConfigureAgents = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure agents",
	}
	apiTokenSchemaCanConfigureAuthenticationMethods = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure authentication methods",
	}
	apiTokenSchemaCanConfigureApplications = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure applications",
	}
	apiTokenSchemaCanConfigureTeams = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure teams (Groups)",
	}
	apiTokenSchemaCanConfigureReleases = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure releases",
	}
	apiTokenSchemaCanConfigureLogManagement = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure log management",
	}
	apiTokenSchemaCanCreatePublicCustomDashboards = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to create public custom dashboards",
	}
	apiTokenSchemaCanViewLogs = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view logs",
	}
	apiTokenSchemaCanViewTraceDetails = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view trace details",
	}
	apiTokenSchemaCanConfigureSessionSettings = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure session settings",
	}
	apiTokenSchemaCanConfigureServiceLevelIndicators = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure service level indicators",
	}
	apiTokenSchemaCanConfigureGlobalAlertPayload = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure global alert payload",
	}
	apiTokenSchemaCanConfigureGlobalAlertConfigs = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure global alert configs",
	}
	apiTokenSchemaCanConfigureGlobalApplicationSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Global Application Smart Alerts",
	}
	apiTokenSchemaCanConfigureGlobalSyntheticSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Global Synthetic Smart Alerts",
	}
	apiTokenSchemaCanConfigureGlobalInfraSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Global Infra Smart Alerts",
	}
	apiTokenSchemaCanConfigureGlobalLogSmartAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Global Log Smart Alerts",
	}
	apiTokenSchemaCanViewAccountAndBillingInformation = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view account and billing information",
	}
	apiTokenSchemaCanEditAllAccessibleCustomDashboards = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to edit all accessible custom dashboards",
	}

	// New schema definitions for scope limitations
	apiTokenSchemaLimitedApplicationsScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited applications scope",
	}
	apiTokenSchemaLimitedBizOpsScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited business operations scope",
	}
	apiTokenSchemaLimitedWebsitesScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited websites scope",
	}
	apiTokenSchemaLimitedKubernetesScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited kubernetes scope",
	}
	apiTokenSchemaLimitedMobileAppsScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited mobile apps scope",
	}
	apiTokenSchemaLimitedInfrastructureScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited infrastructure scope",
	}
	apiTokenSchemaLimitedSyntheticsScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited synthetics scope",
	}
	apiTokenSchemaLimitedVsphereScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited vsphere scope",
	}
	apiTokenSchemaLimitedPhmcScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited phmc scope",
	}
	apiTokenSchemaLimitedPvcScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited pvc scope",
	}
	apiTokenSchemaLimitedZhmcScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited zhmc scope",
	}
	apiTokenSchemaLimitedPcfScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited pcf scope",
	}
	apiTokenSchemaLimitedOpenstackScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited openstack scope",
	}
	apiTokenSchemaLimitedAutomationScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited automation scope",
	}
	apiTokenSchemaLimitedLogsScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited logs scope",
	}
	apiTokenSchemaLimitedNutanixScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited nutanix scope",
	}
	apiTokenSchemaLimitedXenServerScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited xen server scope",
	}
	apiTokenSchemaLimitedWindowsHypervisorScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited windows hypervisor scope",
	}
	apiTokenSchemaLimitedAlertChannelsScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited alert channels scope",
	}
	apiTokenSchemaLimitedLinuxKvmHypervisorScope = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token has limited linux kvm hypervisor scope",
	}

	// New schema definitions for additional permissions
	apiTokenSchemaCanConfigurePersonalAPITokens = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure personal API tokens",
	}
	apiTokenSchemaCanConfigureDatabaseManagement = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure database management",
	}
	apiTokenSchemaCanConfigureAutomationActions = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure automation actions",
	}
	apiTokenSchemaCanConfigureAutomationPolicies = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure automation policies",
	}
	apiTokenSchemaCanRunAutomationActions = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to run automation actions",
	}
	apiTokenSchemaCanDeleteAutomationActionHistory = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to delete automation action history",
	}
	apiTokenSchemaCanConfigureSyntheticTests = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure synthetic tests",
	}
	apiTokenSchemaCanConfigureSyntheticLocations = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure synthetic locations",
	}
	apiTokenSchemaCanConfigureSyntheticCredentials = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure synthetic credentials",
	}
	apiTokenSchemaCanViewSyntheticTests = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view synthetic tests",
	}
	apiTokenSchemaCanViewSyntheticLocations = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view synthetic locations",
	}
	apiTokenSchemaCanViewSyntheticTestResults = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view synthetic test results",
	}
	apiTokenSchemaCanUseSyntheticCredentials = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to use synthetic credentials",
	}
	apiTokenSchemaCanConfigureBizops = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure business operations",
	}
	apiTokenSchemaCanViewBusinessProcesses = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view business processes",
	}
	apiTokenSchemaCanViewBusinessProcessDetails = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view business process details",
	}
	apiTokenSchemaCanViewBusinessActivities = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view business activities",
	}
	apiTokenSchemaCanViewBizAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view business alerts",
	}
	apiTokenSchemaCanDeleteLogs = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to delete logs",
	}
	apiTokenSchemaCanCreateHeapDump = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to create heap dumps",
	}
	apiTokenSchemaCanCreateThreadDump = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to create thread dumps",
	}
	apiTokenSchemaCanManuallyCloseIssue = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to manually close issues",
	}
	apiTokenSchemaCanViewLogVolume = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view log volume",
	}
	apiTokenSchemaCanConfigureLogRetentionPeriod = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure log retention period",
	}
	apiTokenSchemaCanConfigureSubtraces = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure subtraces",
	}
	apiTokenSchemaCanInvokeAlertChannel = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to invoke alert channels",
	}
	apiTokenSchemaCanConfigureLlm = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure LLM",
	}
)

// NewAPITokenResourceHandle creates a ResourceHandle instance for the terraform resource API token
func NewAPITokenResourceHandle() ResourceHandle[*restapi.APIToken] {
	internalIDFieldName := APITokenFieldInternalID
	return &apiTokenResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAPIToken,
			Schema: map[string]*schema.Schema{
				APITokenFieldAccessGrantingToken:                      apiTokenSchemaAccessGrantingToken,
				APITokenFieldInternalID:                               apiTokenSchemaInternalID,
				APITokenFieldName:                                     apiTokenSchemaName,
				APITokenFieldCanConfigureServiceMapping:               apiTokenSchemaCanConfigureServiceMapping,
				APITokenFieldCanConfigureEumApplications:              apiTokenSchemaCanConfigureEumApplications,
				APITokenFieldCanConfigureMobileAppMonitoring:          apiTokenSchemaCanConfigureMobileAppMonitoring,
				APITokenFieldCanConfigureUsers:                        apiTokenSchemaCanConfigureUsers,
				APITokenFieldCanInstallNewAgents:                      apiTokenSchemaCanInstallNewAgents,
				APITokenFieldCanConfigureIntegrations:                 apiTokenSchemaCanConfigureIntegrations,
				APITokenFieldCanConfigureEventsAndAlerts:              apiTokenSchemaCanConfigureEventsAndAlerts,
				APITokenFieldCanConfigureMaintenanceWindows:           apiTokenSchemaCanConfigureMaintenanceWindows,
				APITokenFieldCanConfigureApplicationSmartAlerts:       apiTokenSchemaCanConfigureApplicationSmartAlerts,
				APITokenFieldCanConfigureWebsiteSmartAlerts:           apiTokenSchemaCanConfigureWebsiteSmartAlerts,
				APITokenFieldCanConfigureMobileAppSmartAlerts:         apiTokenSchemaCanConfigureMobileAppSmartAlerts,
				APITokenFieldCanConfigureAPITokens:                    apiTokenSchemaCanConfigureAPITokens,
				APITokenFieldCanConfigureAgentRunMode:                 apiTokenSchemaCanConfigureAgentRunMode,
				APITokenFieldCanViewAuditLog:                          apiTokenSchemaCanViewAuditLog,
				APITokenFieldCanConfigureAgents:                       apiTokenSchemaCanConfigureAgents,
				APITokenFieldCanConfigureAuthenticationMethods:        apiTokenSchemaCanConfigureAuthenticationMethods,
				APITokenFieldCanConfigureApplications:                 apiTokenSchemaCanConfigureApplications,
				APITokenFieldCanConfigureTeams:                        apiTokenSchemaCanConfigureTeams,
				APITokenFieldCanConfigureReleases:                     apiTokenSchemaCanConfigureReleases,
				APITokenFieldCanConfigureLogManagement:                apiTokenSchemaCanConfigureLogManagement,
				APITokenFieldCanCreatePublicCustomDashboards:          apiTokenSchemaCanCreatePublicCustomDashboards,
				APITokenFieldCanViewLogs:                              apiTokenSchemaCanViewLogs,
				APITokenFieldCanViewTraceDetails:                      apiTokenSchemaCanViewTraceDetails,
				APITokenFieldCanConfigureSessionSettings:              apiTokenSchemaCanConfigureSessionSettings,
				APITokenFieldCanConfigureGlobalAlertPayload:           apiTokenSchemaCanConfigureGlobalAlertPayload,
				APITokenFieldCanConfigureGlobalApplicationSmartAlerts: apiTokenSchemaCanConfigureGlobalApplicationSmartAlerts,
				APITokenFieldCanConfigureGlobalSyntheticSmartAlerts:   apiTokenSchemaCanConfigureGlobalSyntheticSmartAlerts,
				APITokenFieldCanConfigureGlobalInfraSmartAlerts:       apiTokenSchemaCanConfigureGlobalInfraSmartAlerts,
				APITokenFieldCanConfigureGlobalLogSmartAlerts:         apiTokenSchemaCanConfigureGlobalLogSmartAlerts,
				APITokenFieldCanViewAccountAndBillingInformation:      apiTokenSchemaCanViewAccountAndBillingInformation,
				APITokenFieldCanEditAllAccessibleCustomDashboards:     apiTokenSchemaCanEditAllAccessibleCustomDashboards,

				// Scope limitations
				APITokenFieldLimitedApplicationsScope:       apiTokenSchemaLimitedApplicationsScope,
				APITokenFieldLimitedBizOpsScope:             apiTokenSchemaLimitedBizOpsScope,
				APITokenFieldLimitedWebsitesScope:           apiTokenSchemaLimitedWebsitesScope,
				APITokenFieldLimitedKubernetesScope:         apiTokenSchemaLimitedKubernetesScope,
				APITokenFieldLimitedMobileAppsScope:         apiTokenSchemaLimitedMobileAppsScope,
				APITokenFieldLimitedInfrastructureScope:     apiTokenSchemaLimitedInfrastructureScope,
				APITokenFieldLimitedSyntheticsScope:         apiTokenSchemaLimitedSyntheticsScope,
				APITokenFieldLimitedVsphereScope:            apiTokenSchemaLimitedVsphereScope,
				APITokenFieldLimitedPhmcScope:               apiTokenSchemaLimitedPhmcScope,
				APITokenFieldLimitedPvcScope:                apiTokenSchemaLimitedPvcScope,
				APITokenFieldLimitedZhmcScope:               apiTokenSchemaLimitedZhmcScope,
				APITokenFieldLimitedPcfScope:                apiTokenSchemaLimitedPcfScope,
				APITokenFieldLimitedOpenstackScope:          apiTokenSchemaLimitedOpenstackScope,
				APITokenFieldLimitedAutomationScope:         apiTokenSchemaLimitedAutomationScope,
				APITokenFieldLimitedLogsScope:               apiTokenSchemaLimitedLogsScope,
				APITokenFieldLimitedNutanixScope:            apiTokenSchemaLimitedNutanixScope,
				APITokenFieldLimitedXenServerScope:          apiTokenSchemaLimitedXenServerScope,
				APITokenFieldLimitedWindowsHypervisorScope:  apiTokenSchemaLimitedWindowsHypervisorScope,
				APITokenFieldLimitedAlertChannelsScope:      apiTokenSchemaLimitedAlertChannelsScope,
				APITokenFieldLimitedLinuxKvmHypervisorScope: apiTokenSchemaLimitedLinuxKvmHypervisorScope,

				// Additional permissions
				APITokenFieldCanConfigurePersonalAPITokens:    apiTokenSchemaCanConfigurePersonalAPITokens,
				APITokenFieldCanConfigureDatabaseManagement:   apiTokenSchemaCanConfigureDatabaseManagement,
				APITokenFieldCanConfigureAutomationActions:    apiTokenSchemaCanConfigureAutomationActions,
				APITokenFieldCanConfigureAutomationPolicies:   apiTokenSchemaCanConfigureAutomationPolicies,
				APITokenFieldCanRunAutomationActions:          apiTokenSchemaCanRunAutomationActions,
				APITokenFieldCanDeleteAutomationActionHistory: apiTokenSchemaCanDeleteAutomationActionHistory,
				APITokenFieldCanConfigureSyntheticTests:       apiTokenSchemaCanConfigureSyntheticTests,
				APITokenFieldCanConfigureSyntheticLocations:   apiTokenSchemaCanConfigureSyntheticLocations,
				APITokenFieldCanConfigureSyntheticCredentials: apiTokenSchemaCanConfigureSyntheticCredentials,
				APITokenFieldCanViewSyntheticTests:            apiTokenSchemaCanViewSyntheticTests,
				APITokenFieldCanViewSyntheticLocations:        apiTokenSchemaCanViewSyntheticLocations,
				APITokenFieldCanViewSyntheticTestResults:      apiTokenSchemaCanViewSyntheticTestResults,
				APITokenFieldCanUseSyntheticCredentials:       apiTokenSchemaCanUseSyntheticCredentials,
				APITokenFieldCanConfigureBizops:               apiTokenSchemaCanConfigureBizops,
				APITokenFieldCanViewBusinessProcesses:         apiTokenSchemaCanViewBusinessProcesses,
				APITokenFieldCanViewBusinessProcessDetails:    apiTokenSchemaCanViewBusinessProcessDetails,
				APITokenFieldCanViewBusinessActivities:        apiTokenSchemaCanViewBusinessActivities,
				APITokenFieldCanViewBizAlerts:                 apiTokenSchemaCanViewBizAlerts,
				APITokenFieldCanDeleteLogs:                    apiTokenSchemaCanDeleteLogs,
				APITokenFieldCanCreateHeapDump:                apiTokenSchemaCanCreateHeapDump,
				APITokenFieldCanCreateThreadDump:              apiTokenSchemaCanCreateThreadDump,
				APITokenFieldCanManuallyCloseIssue:            apiTokenSchemaCanManuallyCloseIssue,
				APITokenFieldCanViewLogVolume:                 apiTokenSchemaCanViewLogVolume,
				APITokenFieldCanConfigureLogRetentionPeriod:   apiTokenSchemaCanConfigureLogRetentionPeriod,
				APITokenFieldCanConfigureSubtraces:            apiTokenSchemaCanConfigureSubtraces,
				APITokenFieldCanInvokeAlertChannel:            apiTokenSchemaCanInvokeAlertChannel,
				APITokenFieldCanConfigureLlm:                  apiTokenSchemaCanConfigureLlm,
			},
			SchemaVersion:    2,
			SkipIDGeneration: true,
			ResourceIDField:  &internalIDFieldName,
		},
	}
}

type apiTokenResource struct {
	metaData ResourceMetaData
}

func (r *apiTokenResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *apiTokenResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV1().CoreConfigSchema().ImpliedType(),
			Upgrade: r.stateUpgradeV1,
			Version: 1,
		},
	}
}

func (r *apiTokenResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.APIToken] {
	return api.APITokens()
}

func (r *apiTokenResource) SetComputedFields(d *schema.ResourceData) error {
	return tfutils.UpdateState(d, map[string]interface{}{
		APITokenFieldInternalID:          RandomID(),
		APITokenFieldAccessGrantingToken: RandomID(),
	})
}

func (r *apiTokenResource) UpdateState(d *schema.ResourceData, apiToken *restapi.APIToken) error {
	d.SetId(apiToken.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		APITokenFieldAccessGrantingToken:                      apiToken.AccessGrantingToken,
		APITokenFieldInternalID:                               apiToken.InternalID,
		APITokenFieldName:                                     apiToken.Name,
		APITokenFieldCanConfigureServiceMapping:               apiToken.CanConfigureServiceMapping,
		APITokenFieldCanConfigureEumApplications:              apiToken.CanConfigureEumApplications,
		APITokenFieldCanConfigureMobileAppMonitoring:          apiToken.CanConfigureMobileAppMonitoring,
		APITokenFieldCanConfigureUsers:                        apiToken.CanConfigureUsers,
		APITokenFieldCanInstallNewAgents:                      apiToken.CanInstallNewAgents,
		APITokenFieldCanConfigureIntegrations:                 apiToken.CanConfigureIntegrations,
		APITokenFieldCanConfigureEventsAndAlerts:              apiToken.CanConfigureEventsAndAlerts,
		APITokenFieldCanConfigureMaintenanceWindows:           apiToken.CanConfigureMaintenanceWindows,
		APITokenFieldCanConfigureApplicationSmartAlerts:       apiToken.CanConfigureApplicationSmartAlerts,
		APITokenFieldCanConfigureWebsiteSmartAlerts:           apiToken.CanConfigureWebsiteSmartAlerts,
		APITokenFieldCanConfigureMobileAppSmartAlerts:         apiToken.CanConfigureMobileAppSmartAlerts,
		APITokenFieldCanConfigureAPITokens:                    apiToken.CanConfigureAPITokens,
		APITokenFieldCanConfigureAgentRunMode:                 apiToken.CanConfigureAgentRunMode,
		APITokenFieldCanViewAuditLog:                          apiToken.CanViewAuditLog,
		APITokenFieldCanConfigureAgents:                       apiToken.CanConfigureAgents,
		APITokenFieldCanConfigureAuthenticationMethods:        apiToken.CanConfigureAuthenticationMethods,
		APITokenFieldCanConfigureApplications:                 apiToken.CanConfigureApplications,
		APITokenFieldCanConfigureTeams:                        apiToken.CanConfigureTeams,
		APITokenFieldCanConfigureReleases:                     apiToken.CanConfigureReleases,
		APITokenFieldCanConfigureLogManagement:                apiToken.CanConfigureLogManagement,
		APITokenFieldCanCreatePublicCustomDashboards:          apiToken.CanCreatePublicCustomDashboards,
		APITokenFieldCanViewLogs:                              apiToken.CanViewLogs,
		APITokenFieldCanViewTraceDetails:                      apiToken.CanViewTraceDetails,
		APITokenFieldCanConfigureSessionSettings:              apiToken.CanConfigureSessionSettings,
		APITokenFieldCanConfigureGlobalAlertPayload:           apiToken.CanConfigureGlobalAlertPayload,
		APITokenFieldCanConfigureGlobalApplicationSmartAlerts: apiToken.CanConfigureGlobalApplicationSmartAlerts,
		APITokenFieldCanConfigureGlobalSyntheticSmartAlerts:   apiToken.CanConfigureGlobalSyntheticSmartAlerts,
		APITokenFieldCanConfigureGlobalInfraSmartAlerts:       apiToken.CanConfigureGlobalInfraSmartAlerts,
		APITokenFieldCanConfigureGlobalLogSmartAlerts:         apiToken.CanConfigureGlobalLogSmartAlerts,
		APITokenFieldCanViewAccountAndBillingInformation:      apiToken.CanViewAccountAndBillingInformation,
		APITokenFieldCanEditAllAccessibleCustomDashboards:     apiToken.CanEditAllAccessibleCustomDashboards,

		// Scope limitations
		APITokenFieldLimitedApplicationsScope:       apiToken.LimitedApplicationsScope,
		APITokenFieldLimitedBizOpsScope:             apiToken.LimitedBizOpsScope,
		APITokenFieldLimitedWebsitesScope:           apiToken.LimitedWebsitesScope,
		APITokenFieldLimitedKubernetesScope:         apiToken.LimitedKubernetesScope,
		APITokenFieldLimitedMobileAppsScope:         apiToken.LimitedMobileAppsScope,
		APITokenFieldLimitedInfrastructureScope:     apiToken.LimitedInfrastructureScope,
		APITokenFieldLimitedSyntheticsScope:         apiToken.LimitedSyntheticsScope,
		APITokenFieldLimitedVsphereScope:            apiToken.LimitedVsphereScope,
		APITokenFieldLimitedPhmcScope:               apiToken.LimitedPhmcScope,
		APITokenFieldLimitedPvcScope:                apiToken.LimitedPvcScope,
		APITokenFieldLimitedZhmcScope:               apiToken.LimitedZhmcScope,
		APITokenFieldLimitedPcfScope:                apiToken.LimitedPcfScope,
		APITokenFieldLimitedOpenstackScope:          apiToken.LimitedOpenstackScope,
		APITokenFieldLimitedAutomationScope:         apiToken.LimitedAutomationScope,
		APITokenFieldLimitedLogsScope:               apiToken.LimitedLogsScope,
		APITokenFieldLimitedNutanixScope:            apiToken.LimitedNutanixScope,
		APITokenFieldLimitedXenServerScope:          apiToken.LimitedXenServerScope,
		APITokenFieldLimitedWindowsHypervisorScope:  apiToken.LimitedWindowsHypervisorScope,
		APITokenFieldLimitedAlertChannelsScope:      apiToken.LimitedAlertChannelsScope,
		APITokenFieldLimitedLinuxKvmHypervisorScope: apiToken.LimitedLinuxKvmHypervisorScope,

		// Additional permissions
		APITokenFieldCanConfigurePersonalAPITokens:    apiToken.CanConfigurePersonalAPITokens,
		APITokenFieldCanConfigureDatabaseManagement:   apiToken.CanConfigureDatabaseManagement,
		APITokenFieldCanConfigureAutomationActions:    apiToken.CanConfigureAutomationActions,
		APITokenFieldCanConfigureAutomationPolicies:   apiToken.CanConfigureAutomationPolicies,
		APITokenFieldCanRunAutomationActions:          apiToken.CanRunAutomationActions,
		APITokenFieldCanDeleteAutomationActionHistory: apiToken.CanDeleteAutomationActionHistory,
		APITokenFieldCanConfigureSyntheticTests:       apiToken.CanConfigureSyntheticTests,
		APITokenFieldCanConfigureSyntheticLocations:   apiToken.CanConfigureSyntheticLocations,
		APITokenFieldCanConfigureSyntheticCredentials: apiToken.CanConfigureSyntheticCredentials,
		APITokenFieldCanViewSyntheticTests:            apiToken.CanViewSyntheticTests,
		APITokenFieldCanViewSyntheticLocations:        apiToken.CanViewSyntheticLocations,
		APITokenFieldCanViewSyntheticTestResults:      apiToken.CanViewSyntheticTestResults,
		APITokenFieldCanUseSyntheticCredentials:       apiToken.CanUseSyntheticCredentials,
		APITokenFieldCanConfigureBizops:               apiToken.CanConfigureBizops,
		APITokenFieldCanViewBusinessProcesses:         apiToken.CanViewBusinessProcesses,
		APITokenFieldCanViewBusinessProcessDetails:    apiToken.CanViewBusinessProcessDetails,
		APITokenFieldCanViewBusinessActivities:        apiToken.CanViewBusinessActivities,
		APITokenFieldCanViewBizAlerts:                 apiToken.CanViewBizAlerts,
		APITokenFieldCanDeleteLogs:                    apiToken.CanDeleteLogs,
		APITokenFieldCanCreateHeapDump:                apiToken.CanCreateHeapDump,
		APITokenFieldCanCreateThreadDump:              apiToken.CanCreateThreadDump,
		APITokenFieldCanManuallyCloseIssue:            apiToken.CanManuallyCloseIssue,
		APITokenFieldCanViewLogVolume:                 apiToken.CanViewLogVolume,
		APITokenFieldCanConfigureLogRetentionPeriod:   apiToken.CanConfigureLogRetentionPeriod,
		APITokenFieldCanConfigureSubtraces:            apiToken.CanConfigureSubtraces,
		APITokenFieldCanInvokeAlertChannel:            apiToken.CanInvokeAlertChannel,
		APITokenFieldCanConfigureLlm:                  apiToken.CanConfigureLlm,
	})
}

func (r *apiTokenResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.APIToken, error) {
	return &restapi.APIToken{
		ID:                                       d.Id(),
		AccessGrantingToken:                      d.Get(APITokenFieldAccessGrantingToken).(string),
		InternalID:                               d.Get(APITokenFieldInternalID).(string),
		Name:                                     d.Get(APITokenFieldName).(string),
		CanConfigureServiceMapping:               d.Get(APITokenFieldCanConfigureServiceMapping).(bool),
		CanConfigureEumApplications:              d.Get(APITokenFieldCanConfigureEumApplications).(bool),
		CanConfigureMobileAppMonitoring:          d.Get(APITokenFieldCanConfigureMobileAppMonitoring).(bool),
		CanConfigureUsers:                        d.Get(APITokenFieldCanConfigureUsers).(bool),
		CanInstallNewAgents:                      d.Get(APITokenFieldCanInstallNewAgents).(bool),
		CanConfigureIntegrations:                 d.Get(APITokenFieldCanConfigureIntegrations).(bool),
		CanConfigureEventsAndAlerts:              d.Get(APITokenFieldCanConfigureEventsAndAlerts).(bool),
		CanConfigureMaintenanceWindows:           d.Get(APITokenFieldCanConfigureMaintenanceWindows).(bool),
		CanConfigureApplicationSmartAlerts:       d.Get(APITokenFieldCanConfigureApplicationSmartAlerts).(bool),
		CanConfigureWebsiteSmartAlerts:           d.Get(APITokenFieldCanConfigureWebsiteSmartAlerts).(bool),
		CanConfigureMobileAppSmartAlerts:         d.Get(APITokenFieldCanConfigureMobileAppSmartAlerts).(bool),
		CanConfigureAPITokens:                    d.Get(APITokenFieldCanConfigureAPITokens).(bool),
		CanConfigureAgentRunMode:                 d.Get(APITokenFieldCanConfigureAgentRunMode).(bool),
		CanViewAuditLog:                          d.Get(APITokenFieldCanViewAuditLog).(bool),
		CanConfigureAgents:                       d.Get(APITokenFieldCanConfigureAgents).(bool),
		CanConfigureAuthenticationMethods:        d.Get(APITokenFieldCanConfigureAuthenticationMethods).(bool),
		CanConfigureApplications:                 d.Get(APITokenFieldCanConfigureApplications).(bool),
		CanConfigureTeams:                        d.Get(APITokenFieldCanConfigureTeams).(bool),
		CanConfigureReleases:                     d.Get(APITokenFieldCanConfigureReleases).(bool),
		CanConfigureLogManagement:                d.Get(APITokenFieldCanConfigureLogManagement).(bool),
		CanCreatePublicCustomDashboards:          d.Get(APITokenFieldCanCreatePublicCustomDashboards).(bool),
		CanViewLogs:                              d.Get(APITokenFieldCanViewLogs).(bool),
		CanViewTraceDetails:                      d.Get(APITokenFieldCanViewTraceDetails).(bool),
		CanConfigureSessionSettings:              d.Get(APITokenFieldCanConfigureSessionSettings).(bool),
		CanConfigureGlobalAlertPayload:           d.Get(APITokenFieldCanConfigureGlobalAlertPayload).(bool),
		CanConfigureGlobalApplicationSmartAlerts: d.Get(APITokenFieldCanConfigureGlobalApplicationSmartAlerts).(bool),
		CanConfigureGlobalSyntheticSmartAlerts:   d.Get(APITokenFieldCanConfigureGlobalSyntheticSmartAlerts).(bool),
		CanConfigureGlobalInfraSmartAlerts:       d.Get(APITokenFieldCanConfigureGlobalInfraSmartAlerts).(bool),
		CanConfigureGlobalLogSmartAlerts:         d.Get(APITokenFieldCanConfigureGlobalLogSmartAlerts).(bool),
		CanViewAccountAndBillingInformation:      d.Get(APITokenFieldCanViewAccountAndBillingInformation).(bool),
		CanEditAllAccessibleCustomDashboards:     d.Get(APITokenFieldCanEditAllAccessibleCustomDashboards).(bool),

		// Scope limitations
		LimitedApplicationsScope:       d.Get(APITokenFieldLimitedApplicationsScope).(bool),
		LimitedBizOpsScope:             d.Get(APITokenFieldLimitedBizOpsScope).(bool),
		LimitedWebsitesScope:           d.Get(APITokenFieldLimitedWebsitesScope).(bool),
		LimitedKubernetesScope:         d.Get(APITokenFieldLimitedKubernetesScope).(bool),
		LimitedMobileAppsScope:         d.Get(APITokenFieldLimitedMobileAppsScope).(bool),
		LimitedInfrastructureScope:     d.Get(APITokenFieldLimitedInfrastructureScope).(bool),
		LimitedSyntheticsScope:         d.Get(APITokenFieldLimitedSyntheticsScope).(bool),
		LimitedVsphereScope:            d.Get(APITokenFieldLimitedVsphereScope).(bool),
		LimitedPhmcScope:               d.Get(APITokenFieldLimitedPhmcScope).(bool),
		LimitedPvcScope:                d.Get(APITokenFieldLimitedPvcScope).(bool),
		LimitedZhmcScope:               d.Get(APITokenFieldLimitedZhmcScope).(bool),
		LimitedPcfScope:                d.Get(APITokenFieldLimitedPcfScope).(bool),
		LimitedOpenstackScope:          d.Get(APITokenFieldLimitedOpenstackScope).(bool),
		LimitedAutomationScope:         d.Get(APITokenFieldLimitedAutomationScope).(bool),
		LimitedLogsScope:               d.Get(APITokenFieldLimitedLogsScope).(bool),
		LimitedNutanixScope:            d.Get(APITokenFieldLimitedNutanixScope).(bool),
		LimitedXenServerScope:          d.Get(APITokenFieldLimitedXenServerScope).(bool),
		LimitedWindowsHypervisorScope:  d.Get(APITokenFieldLimitedWindowsHypervisorScope).(bool),
		LimitedAlertChannelsScope:      d.Get(APITokenFieldLimitedAlertChannelsScope).(bool),
		LimitedLinuxKvmHypervisorScope: d.Get(APITokenFieldLimitedLinuxKvmHypervisorScope).(bool),

		// Additional permissions
		CanConfigurePersonalAPITokens:    d.Get(APITokenFieldCanConfigurePersonalAPITokens).(bool),
		CanConfigureDatabaseManagement:   d.Get(APITokenFieldCanConfigureDatabaseManagement).(bool),
		CanConfigureAutomationActions:    d.Get(APITokenFieldCanConfigureAutomationActions).(bool),
		CanConfigureAutomationPolicies:   d.Get(APITokenFieldCanConfigureAutomationPolicies).(bool),
		CanRunAutomationActions:          d.Get(APITokenFieldCanRunAutomationActions).(bool),
		CanDeleteAutomationActionHistory: d.Get(APITokenFieldCanDeleteAutomationActionHistory).(bool),
		CanConfigureSyntheticTests:       d.Get(APITokenFieldCanConfigureSyntheticTests).(bool),
		CanConfigureSyntheticLocations:   d.Get(APITokenFieldCanConfigureSyntheticLocations).(bool),
		CanConfigureSyntheticCredentials: d.Get(APITokenFieldCanConfigureSyntheticCredentials).(bool),
		CanViewSyntheticTests:            d.Get(APITokenFieldCanViewSyntheticTests).(bool),
		CanViewSyntheticLocations:        d.Get(APITokenFieldCanViewSyntheticLocations).(bool),
		CanViewSyntheticTestResults:      d.Get(APITokenFieldCanViewSyntheticTestResults).(bool),
		CanUseSyntheticCredentials:       d.Get(APITokenFieldCanUseSyntheticCredentials).(bool),
		CanConfigureBizops:               d.Get(APITokenFieldCanConfigureBizops).(bool),
		CanViewBusinessProcesses:         d.Get(APITokenFieldCanViewBusinessProcesses).(bool),
		CanViewBusinessProcessDetails:    d.Get(APITokenFieldCanViewBusinessProcessDetails).(bool),
		CanViewBusinessActivities:        d.Get(APITokenFieldCanViewBusinessActivities).(bool),
		CanViewBizAlerts:                 d.Get(APITokenFieldCanViewBizAlerts).(bool),
		CanDeleteLogs:                    d.Get(APITokenFieldCanDeleteLogs).(bool),
		CanCreateHeapDump:                d.Get(APITokenFieldCanCreateHeapDump).(bool),
		CanCreateThreadDump:              d.Get(APITokenFieldCanCreateThreadDump).(bool),
		CanManuallyCloseIssue:            d.Get(APITokenFieldCanManuallyCloseIssue).(bool),
		CanViewLogVolume:                 d.Get(APITokenFieldCanViewLogVolume).(bool),
		CanConfigureLogRetentionPeriod:   d.Get(APITokenFieldCanConfigureLogRetentionPeriod).(bool),
		CanConfigureSubtraces:            d.Get(APITokenFieldCanConfigureSubtraces).(bool),
		CanInvokeAlertChannel:            d.Get(APITokenFieldCanInvokeAlertChannel).(bool),
		CanConfigureLlm:                  d.Get(APITokenFieldCanConfigureLlm).(bool),
	}, nil
}

func (r *apiTokenResource) stateUpgradeV1(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[APITokenFieldFullName]; ok {
		state[APITokenFieldName] = state[APITokenFieldFullName]
		delete(state, APITokenFieldFullName)
	}
	return state, nil
}

func (r *apiTokenResource) schemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			APITokenFieldAccessGrantingToken:                      apiTokenSchemaAccessGrantingToken,
			APITokenFieldInternalID:                               apiTokenSchemaInternalID,
			APITokenFieldName:                                     apiTokenSchemaName,
			APITokenFieldFullName:                                 apiTokenSchemaFullName,
			APITokenFieldCanConfigureServiceMapping:               apiTokenSchemaCanConfigureServiceMapping,
			APITokenFieldCanConfigureEumApplications:              apiTokenSchemaCanConfigureEumApplications,
			APITokenFieldCanConfigureMobileAppMonitoring:          apiTokenSchemaCanConfigureMobileAppMonitoring,
			APITokenFieldCanConfigureUsers:                        apiTokenSchemaCanConfigureUsers,
			APITokenFieldCanInstallNewAgents:                      apiTokenSchemaCanInstallNewAgents,
			APITokenFieldCanConfigureIntegrations:                 apiTokenSchemaCanConfigureIntegrations,
			APITokenFieldCanConfigureEventsAndAlerts:              apiTokenSchemaCanConfigureEventsAndAlerts,
			APITokenFieldCanConfigureMaintenanceWindows:           apiTokenSchemaCanConfigureMaintenanceWindows,
			APITokenFieldCanConfigureApplicationSmartAlerts:       apiTokenSchemaCanConfigureApplicationSmartAlerts,
			APITokenFieldCanConfigureWebsiteSmartAlerts:           apiTokenSchemaCanConfigureWebsiteSmartAlerts,
			APITokenFieldCanConfigureMobileAppSmartAlerts:         apiTokenSchemaCanConfigureMobileAppSmartAlerts,
			APITokenFieldCanConfigureAPITokens:                    apiTokenSchemaCanConfigureAPITokens,
			APITokenFieldCanConfigureAgentRunMode:                 apiTokenSchemaCanConfigureAgentRunMode,
			APITokenFieldCanViewAuditLog:                          apiTokenSchemaCanViewAuditLog,
			APITokenFieldCanConfigureAgents:                       apiTokenSchemaCanConfigureAgents,
			APITokenFieldCanConfigureAuthenticationMethods:        apiTokenSchemaCanConfigureAuthenticationMethods,
			APITokenFieldCanConfigureApplications:                 apiTokenSchemaCanConfigureApplications,
			APITokenFieldCanConfigureTeams:                        apiTokenSchemaCanConfigureTeams,
			APITokenFieldCanConfigureReleases:                     apiTokenSchemaCanConfigureReleases,
			APITokenFieldCanConfigureLogManagement:                apiTokenSchemaCanConfigureLogManagement,
			APITokenFieldCanCreatePublicCustomDashboards:          apiTokenSchemaCanCreatePublicCustomDashboards,
			APITokenFieldCanViewLogs:                              apiTokenSchemaCanViewLogs,
			APITokenFieldCanViewTraceDetails:                      apiTokenSchemaCanViewTraceDetails,
			APITokenFieldCanConfigureSessionSettings:              apiTokenSchemaCanConfigureSessionSettings,
			APITokenFieldCanConfigureGlobalAlertPayload:           apiTokenSchemaCanConfigureGlobalAlertPayload,
			APITokenFieldCanConfigureGlobalApplicationSmartAlerts: apiTokenSchemaCanConfigureGlobalApplicationSmartAlerts,
			APITokenFieldCanConfigureGlobalSyntheticSmartAlerts:   apiTokenSchemaCanConfigureGlobalSyntheticSmartAlerts,
			APITokenFieldCanConfigureGlobalInfraSmartAlerts:       apiTokenSchemaCanConfigureGlobalInfraSmartAlerts,
			APITokenFieldCanConfigureGlobalLogSmartAlerts:         apiTokenSchemaCanConfigureGlobalLogSmartAlerts,
			APITokenFieldCanViewAccountAndBillingInformation:      apiTokenSchemaCanViewAccountAndBillingInformation,
			APITokenFieldCanEditAllAccessibleCustomDashboards:     apiTokenSchemaCanEditAllAccessibleCustomDashboards,
		},
	}
}
