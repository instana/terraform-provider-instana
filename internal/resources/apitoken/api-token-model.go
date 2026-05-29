package apitoken

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// APITokenModel represents the data model for the API token resource
type APITokenModel struct {
	// This struct is now aligned exactly with the schema attributes for the API token resource.
	ID                                       types.String `tfsdk:"id"`
	AccessGrantingToken                      types.String `tfsdk:"access_granting_token"`
	InternalID                               types.String `tfsdk:"internal_id"`
	Name                                     types.String `tfsdk:"name"`
	CanConfigureEventsAndAlerts              types.Bool   `tfsdk:"can_configure_events_and_alerts"`
	CanConfigureMaintenanceWindows           types.Bool   `tfsdk:"can_configure_maintenance_windows"`
	CanConfigureApplicationSmartAlerts       types.Bool   `tfsdk:"can_configure_application_smart_alerts"`
	CanConfigureWebsiteSmartAlerts           types.Bool   `tfsdk:"can_configure_website_smart_alerts"`
	CanConfigureMobileAppSmartAlerts         types.Bool   `tfsdk:"can_configure_mobile_app_smart_alerts"`
	CanConfigureServiceLevelCorrectionWindows types.Bool  `tfsdk:"can_configure_service_level_correction_windows"`
	CanConfigureServiceLevelSmartAlerts      types.Bool   `tfsdk:"can_configure_service_level_smart_alerts"`
	CanConfigureServiceLevels                types.Bool   `tfsdk:"can_configure_service_levels"`
	CanViewTraceDetails                      types.Bool   `tfsdk:"can_view_trace_details"`
	CanConfigureSessionSettings              types.Bool   `tfsdk:"can_configure_session_settings"`
	CanConfigureGlobalAlertPayload           types.Bool   `tfsdk:"can_configure_global_alert_payload"`
	CanConfigureGlobalApplicationSmartAlerts types.Bool   `tfsdk:"can_configure_global_application_smart_alerts"`
	CanConfigureGlobalSyntheticSmartAlerts   types.Bool   `tfsdk:"can_configure_global_synthetic_smart_alerts"`
	CanConfigureGlobalInfraSmartAlerts       types.Bool   `tfsdk:"can_configure_global_infra_smart_alerts"`
	CanConfigureGlobalLogSmartAlerts         types.Bool   `tfsdk:"can_configure_global_log_smart_alerts"`
	CanViewAccountAndBillingInformation      types.Bool   `tfsdk:"can_view_account_and_billing_information"`
	CanEditAllAccessibleCustomDashboards     types.Bool   `tfsdk:"can_edit_all_accessible_custom_dashboards"`
	LimitedApplicationsScope                 types.Bool   `tfsdk:"limited_applications_scope"`
	LimitedBizOpsScope                       types.Bool   `tfsdk:"limited_biz_ops_scope"`
	LimitedWebsitesScope                     types.Bool   `tfsdk:"limited_websites_scope"`
	LimitedKubernetesScope                   types.Bool   `tfsdk:"limited_kubernetes_scope"`
	LimitedMobileAppsScope                   types.Bool   `tfsdk:"limited_mobile_apps_scope"`
	LimitedInfrastructureScope               types.Bool   `tfsdk:"limited_infrastructure_scope"`
	LimitedSyntheticsScope                   types.Bool   `tfsdk:"limited_synthetics_scope"`
	LimitedVsphereScope                      types.Bool   `tfsdk:"limited_vsphere_scope"`
	LimitedPhmcScope                         types.Bool   `tfsdk:"limited_phmc_scope"`
	LimitedPvcScope                          types.Bool   `tfsdk:"limited_pvc_scope"`
	LimitedZhmcScope                         types.Bool   `tfsdk:"limited_zhmc_scope"`
	LimitedPcfScope                          types.Bool   `tfsdk:"limited_pcf_scope"`
	LimitedOpenstackScope                    types.Bool   `tfsdk:"limited_openstack_scope"`
	LimitedAutomationScope                   types.Bool   `tfsdk:"limited_automation_scope"`
	LimitedLogsScope                         types.Bool   `tfsdk:"limited_logs_scope"`
	LimitedNutanixScope                      types.Bool   `tfsdk:"limited_nutanix_scope"`
	LimitedXenServerScope                    types.Bool   `tfsdk:"limited_xen_server_scope"`
	LimitedWindowsHypervisorScope            types.Bool   `tfsdk:"limited_windows_hypervisor_scope"`
	LimitedAlertChannelsScope                types.Bool   `tfsdk:"limited_alert_channels_scope"`
	LimitedLinuxKvmHypervisorScope           types.Bool   `tfsdk:"limited_linux_kvm_hypervisor_scope"`
	LimitedServiceLevelScope                 types.Bool   `tfsdk:"limited_service_level_scope"`
	LimitedAiGatewayScope                    types.Bool   `tfsdk:"limited_ai_gateway_scope"`
	CanConfigurePersonalAPITokens            types.Bool   `tfsdk:"can_configure_personal_api_tokens"`
	CanConfigureDatabaseManagement           types.Bool   `tfsdk:"can_configure_database_management"`
	CanConfigureAutomationActions            types.Bool   `tfsdk:"can_configure_automation_actions"`
	CanConfigureAutomationPolicies           types.Bool   `tfsdk:"can_configure_automation_policies"`
	CanRunAutomationActions                  types.Bool   `tfsdk:"can_run_automation_actions"`
	CanDeleteAutomationActionHistory         types.Bool   `tfsdk:"can_delete_automation_action_history"`
	CanConfigureSyntheticTests               types.Bool   `tfsdk:"can_configure_synthetic_tests"`
	CanConfigureSyntheticLocations           types.Bool   `tfsdk:"can_configure_synthetic_locations"`
	CanConfigureSyntheticCredentials         types.Bool   `tfsdk:"can_configure_synthetic_credentials"`
	CanViewSyntheticTests                    types.Bool   `tfsdk:"can_view_synthetic_tests"`
	CanViewSyntheticLocations                types.Bool   `tfsdk:"can_view_synthetic_locations"`
	CanViewSyntheticTestResults              types.Bool   `tfsdk:"can_view_synthetic_test_results"`
	CanUseSyntheticCredentials               types.Bool   `tfsdk:"can_use_synthetic_credentials"`
	CanConfigureBizops                       types.Bool   `tfsdk:"can_configure_bizops"`
	CanViewBusinessProcesses                 types.Bool   `tfsdk:"can_view_business_processes"`
	CanViewBusinessProcessDetails            types.Bool   `tfsdk:"can_view_business_process_details"`
	CanViewBusinessActivities                types.Bool   `tfsdk:"can_view_business_activities"`
	CanViewBizAlerts                         types.Bool   `tfsdk:"can_view_biz_alerts"`
	CanDeleteLogs                            types.Bool   `tfsdk:"can_delete_logs"`
	CanCreateHeapDump                        types.Bool   `tfsdk:"can_create_heap_dump"`
	CanCreateThreadDump                      types.Bool   `tfsdk:"can_create_thread_dump"`
	CanManuallyCloseIssue                    types.Bool   `tfsdk:"can_manually_close_issue"`
	CanViewLogVolume                         types.Bool   `tfsdk:"can_view_log_volume"`
	CanConfigureLogRetentionPeriod           types.Bool   `tfsdk:"can_configure_log_retention_period"`
	CanConfigureSubtraces                    types.Bool   `tfsdk:"can_configure_subtraces"`
	CanInvokeAlertChannel                    types.Bool   `tfsdk:"can_invoke_alert_channel"`
	CanConfigureLlm                          types.Bool   `tfsdk:"can_configure_llm"`
	CanConfigureAiAgents                     types.Bool   `tfsdk:"can_configure_ai_agents"`
	CanConfigureApdex                        types.Bool   `tfsdk:"can_configure_apdex"`
}
