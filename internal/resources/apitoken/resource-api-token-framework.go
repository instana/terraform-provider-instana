package apitoken

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaAPITokenFramework the name of the terraform-provider-instana resource to manage API tokens
const ResourceInstanaAPITokenFramework = "api_token"

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

// APITokenModel represents the data model for the API token resource
type APITokenModel struct {
	ID                                       types.String `tfsdk:"id"`
	AccessGrantingToken                      types.String `tfsdk:"access_granting_token"`
	InternalID                               types.String `tfsdk:"internal_id"`
	Name                                     types.String `tfsdk:"name"`
	CanConfigureServiceMapping               types.Bool   `tfsdk:"can_configure_service_mapping"`
	CanConfigureEumApplications              types.Bool   `tfsdk:"can_configure_eum_applications"`
	CanConfigureMobileAppMonitoring          types.Bool   `tfsdk:"can_configure_mobile_app_monitoring"`
	CanConfigureUsers                        types.Bool   `tfsdk:"can_configure_users"`
	CanInstallNewAgents                      types.Bool   `tfsdk:"can_install_new_agents"`
	CanConfigureIntegrations                 types.Bool   `tfsdk:"can_configure_integrations"`
	CanConfigureEventsAndAlerts              types.Bool   `tfsdk:"can_configure_events_and_alerts"`
	CanConfigureMaintenanceWindows           types.Bool   `tfsdk:"can_configure_maintenance_windows"`
	CanConfigureApplicationSmartAlerts       types.Bool   `tfsdk:"can_configure_application_smart_alerts"`
	CanConfigureWebsiteSmartAlerts           types.Bool   `tfsdk:"can_configure_website_smart_alerts"`
	CanConfigureMobileAppSmartAlerts         types.Bool   `tfsdk:"can_configure_mobile_app_smart_alerts"`
	CanConfigureAPITokens                    types.Bool   `tfsdk:"can_configure_api_tokens"`
	CanConfigureAgentRunMode                 types.Bool   `tfsdk:"can_configure_agent_run_mode"`
	CanViewAuditLog                          types.Bool   `tfsdk:"can_view_audit_log"`
	CanConfigureAgents                       types.Bool   `tfsdk:"can_configure_agents"`
	CanConfigureAuthenticationMethods        types.Bool   `tfsdk:"can_configure_authentication_methods"`
	CanConfigureApplications                 types.Bool   `tfsdk:"can_configure_applications"`
	CanConfigureTeams                        types.Bool   `tfsdk:"can_configure_teams"`
	CanConfigureReleases                     types.Bool   `tfsdk:"can_configure_releases"`
	CanConfigureLogManagement                types.Bool   `tfsdk:"can_configure_log_management"`
	CanCreatePublicCustomDashboards          types.Bool   `tfsdk:"can_create_public_custom_dashboards"`
	CanViewLogs                              types.Bool   `tfsdk:"can_view_logs"`
	CanViewTraceDetails                      types.Bool   `tfsdk:"can_view_trace_details"`
	CanConfigureSessionSettings              types.Bool   `tfsdk:"can_configure_session_settings"`
	CanConfigureGlobalAlertPayload           types.Bool   `tfsdk:"can_configure_global_alert_payload"`
	CanConfigureGlobalApplicationSmartAlerts types.Bool   `tfsdk:"can_configure_global_application_smart_alerts"`
	CanConfigureGlobalSyntheticSmartAlerts   types.Bool   `tfsdk:"can_configure_global_synthetic_smart_alerts"`
	CanConfigureGlobalInfraSmartAlerts       types.Bool   `tfsdk:"can_configure_global_infra_smart_alerts"`
	CanConfigureGlobalLogSmartAlerts         types.Bool   `tfsdk:"can_configure_global_log_smart_alerts"`
	CanViewAccountAndBillingInformation      types.Bool   `tfsdk:"can_view_account_and_billing_information"`
	CanEditAllAccessibleCustomDashboards     types.Bool   `tfsdk:"can_edit_all_accessible_custom_dashboards"`

	// Scope limitations
	LimitedApplicationsScope       types.Bool `tfsdk:"limited_applications_scope"`
	LimitedBizOpsScope             types.Bool `tfsdk:"limited_biz_ops_scope"`
	LimitedWebsitesScope           types.Bool `tfsdk:"limited_websites_scope"`
	LimitedKubernetesScope         types.Bool `tfsdk:"limited_kubernetes_scope"`
	LimitedMobileAppsScope         types.Bool `tfsdk:"limited_mobile_apps_scope"`
	LimitedInfrastructureScope     types.Bool `tfsdk:"limited_infrastructure_scope"`
	LimitedSyntheticsScope         types.Bool `tfsdk:"limited_synthetics_scope"`
	LimitedVsphereScope            types.Bool `tfsdk:"limited_vsphere_scope"`
	LimitedPhmcScope               types.Bool `tfsdk:"limited_phmc_scope"`
	LimitedPvcScope                types.Bool `tfsdk:"limited_pvc_scope"`
	LimitedZhmcScope               types.Bool `tfsdk:"limited_zhmc_scope"`
	LimitedPcfScope                types.Bool `tfsdk:"limited_pcf_scope"`
	LimitedOpenstackScope          types.Bool `tfsdk:"limited_openstack_scope"`
	LimitedAutomationScope         types.Bool `tfsdk:"limited_automation_scope"`
	LimitedLogsScope               types.Bool `tfsdk:"limited_logs_scope"`
	LimitedNutanixScope            types.Bool `tfsdk:"limited_nutanix_scope"`
	LimitedXenServerScope          types.Bool `tfsdk:"limited_xen_server_scope"`
	LimitedWindowsHypervisorScope  types.Bool `tfsdk:"limited_windows_hypervisor_scope"`
	LimitedAlertChannelsScope      types.Bool `tfsdk:"limited_alert_channels_scope"`
	LimitedLinuxKvmHypervisorScope types.Bool `tfsdk:"limited_linux_kvm_hypervisor_scope"`
	LimitedServiceLevelScope       types.Bool `tfsdk:"limited_service_level_scope"`
	LimitedAiGatewayScope          types.Bool `tfsdk:"limited_ai_gateway_scope"`

	// Additional permissions
	CanConfigurePersonalAPITokens             types.Bool `tfsdk:"can_configure_personal_api_tokens"`
	CanConfigureDatabaseManagement            types.Bool `tfsdk:"can_configure_database_management"`
	CanConfigureAutomationActions             types.Bool `tfsdk:"can_configure_automation_actions"`
	CanConfigureAutomationPolicies            types.Bool `tfsdk:"can_configure_automation_policies"`
	CanRunAutomationActions                   types.Bool `tfsdk:"can_run_automation_actions"`
	CanDeleteAutomationActionHistory          types.Bool `tfsdk:"can_delete_automation_action_history"`
	CanConfigureSyntheticTests                types.Bool `tfsdk:"can_configure_synthetic_tests"`
	CanConfigureSyntheticLocations            types.Bool `tfsdk:"can_configure_synthetic_locations"`
	CanConfigureSyntheticCredentials          types.Bool `tfsdk:"can_configure_synthetic_credentials"`
	CanViewSyntheticTests                     types.Bool `tfsdk:"can_view_synthetic_tests"`
	CanViewSyntheticLocations                 types.Bool `tfsdk:"can_view_synthetic_locations"`
	CanViewSyntheticTestResults               types.Bool `tfsdk:"can_view_synthetic_test_results"`
	CanUseSyntheticCredentials                types.Bool `tfsdk:"can_use_synthetic_credentials"`
	CanConfigureBizops                        types.Bool `tfsdk:"can_configure_bizops"`
	CanViewBusinessProcesses                  types.Bool `tfsdk:"can_view_business_processes"`
	CanViewBusinessProcessDetails             types.Bool `tfsdk:"can_view_business_process_details"`
	CanViewBusinessActivities                 types.Bool `tfsdk:"can_view_business_activities"`
	CanViewBizAlerts                          types.Bool `tfsdk:"can_view_biz_alerts"`
	CanDeleteLogs                             types.Bool `tfsdk:"can_delete_logs"`
	CanCreateHeapDump                         types.Bool `tfsdk:"can_create_heap_dump"`
	CanCreateThreadDump                       types.Bool `tfsdk:"can_create_thread_dump"`
	CanManuallyCloseIssue                     types.Bool `tfsdk:"can_manually_close_issue"`
	CanViewLogVolume                          types.Bool `tfsdk:"can_view_log_volume"`
	CanConfigureLogRetentionPeriod            types.Bool `tfsdk:"can_configure_log_retention_period"`
	CanConfigureSubtraces                     types.Bool `tfsdk:"can_configure_subtraces"`
	CanInvokeAlertChannel                     types.Bool `tfsdk:"can_invoke_alert_channel"`
	CanConfigureLlm                           types.Bool `tfsdk:"can_configure_llm"`
	CanConfigureAiAgents                      types.Bool `tfsdk:"can_configure_ai_agents"`
	CanConfigureApdex                         types.Bool `tfsdk:"can_configure_apdex"`
	CanConfigureServiceLevelCorrectionWindows types.Bool `tfsdk:"can_configure_service_level_correction_windows"`
	CanConfigureServiceLevelSmartAlerts       types.Bool `tfsdk:"can_configure_service_level_smart_alerts"`
	CanConfigureServiceLevels                 types.Bool `tfsdk:"can_configure_service_levels"`
}

// NewAPITokenResourceHandleFramework creates the resource handle for API Tokens
func NewAPITokenResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.APIToken] {
	return &apiTokenResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAPITokenFramework,
			Schema: schema.Schema{
				Description: APITokenDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: APITokenDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					APITokenFieldAccessGrantingToken: schema.StringAttribute{
						Computed:    true,
						Description: APITokenDescAccessGrantingToken,
					},
					APITokenFieldInternalID: schema.StringAttribute{
						Computed:    true,
						Description: APITokenDescInternalID,
					},
					APITokenFieldName: schema.StringAttribute{
						Required:    true,
						Description: APITokenDescName,
					},
					APITokenFieldCanConfigureServiceMapping: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureServiceMapping,
					},
					APITokenFieldCanConfigureEumApplications: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureEumApplications,
					},
					APITokenFieldCanConfigureMobileAppMonitoring: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureMobileAppMonitoring,
					},
					APITokenFieldCanConfigureUsers: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureUsers,
					},
					APITokenFieldCanInstallNewAgents: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanInstallNewAgents,
					},
					APITokenFieldCanConfigureIntegrations: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureIntegrations,
					},
					APITokenFieldCanConfigureEventsAndAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureEventsAndAlerts,
					},
					APITokenFieldCanConfigureMaintenanceWindows: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureMaintenanceWindows,
					},
					APITokenFieldCanConfigureApplicationSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureApplicationSmartAlerts,
					},
					APITokenFieldCanConfigureWebsiteSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureWebsiteSmartAlerts,
					},
					APITokenFieldCanConfigureMobileAppSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureMobileAppSmartAlerts,
					},
					APITokenFieldCanConfigureAPITokens: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAPITokens,
					},
					APITokenFieldCanConfigureAgentRunMode: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAgentRunMode,
					},
					APITokenFieldCanViewAuditLog: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewAuditLog,
					},
					APITokenFieldCanConfigureAgents: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAgents,
					},
					APITokenFieldCanConfigureAuthenticationMethods: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAuthenticationMethods,
					},
					APITokenFieldCanConfigureApplications: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureApplications,
					},
					APITokenFieldCanConfigureTeams: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureTeams,
					},
					APITokenFieldCanConfigureReleases: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureReleases,
					},
					APITokenFieldCanConfigureLogManagement: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureLogManagement,
					},
					APITokenFieldCanCreatePublicCustomDashboards: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanCreatePublicCustomDashboards,
					},
					APITokenFieldCanViewLogs: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewLogs,
					},
					APITokenFieldCanViewTraceDetails: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewTraceDetails,
					},
					APITokenFieldCanConfigureSessionSettings: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureSessionSettings,
					},
					APITokenFieldCanConfigureGlobalAlertPayload: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureGlobalAlertPayload,
					},
					APITokenFieldCanConfigureGlobalApplicationSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureGlobalApplicationSmartAlerts,
					},
					APITokenFieldCanConfigureGlobalSyntheticSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureGlobalSyntheticSmartAlerts,
					},
					APITokenFieldCanConfigureGlobalInfraSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureGlobalInfraSmartAlerts,
					},
					APITokenFieldCanConfigureGlobalLogSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureGlobalLogSmartAlerts,
					},
					APITokenFieldCanViewAccountAndBillingInformation: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewAccountAndBillingInformation,
					},
					APITokenFieldCanEditAllAccessibleCustomDashboards: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanEditAllAccessibleCustomDashboards,
					},

					// Scope limitations
					APITokenFieldLimitedApplicationsScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedApplicationsScope,
					},
					APITokenFieldLimitedBizOpsScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedBizOpsScope,
					},
					APITokenFieldLimitedWebsitesScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedWebsitesScope,
					},
					APITokenFieldLimitedKubernetesScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedKubernetesScope,
					},
					APITokenFieldLimitedMobileAppsScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedMobileAppsScope,
					},
					APITokenFieldLimitedInfrastructureScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedInfrastructureScope,
					},
					APITokenFieldLimitedSyntheticsScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedSyntheticsScope,
					},
					APITokenFieldLimitedVsphereScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedVsphereScope,
					},
					APITokenFieldLimitedPhmcScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedPhmcScope,
					},
					APITokenFieldLimitedPvcScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedPvcScope,
					},
					APITokenFieldLimitedZhmcScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedZhmcScope,
					},
					APITokenFieldLimitedPcfScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedPcfScope,
					},
					APITokenFieldLimitedOpenstackScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedOpenstackScope,
					},
					APITokenFieldLimitedAutomationScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedAutomationScope,
					},
					APITokenFieldLimitedLogsScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedLogsScope,
					},
					APITokenFieldLimitedNutanixScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedNutanixScope,
					},
					APITokenFieldLimitedXenServerScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedXenServerScope,
					},
					APITokenFieldLimitedWindowsHypervisorScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedWindowsHypervisorScope,
					},
					APITokenFieldLimitedAlertChannelsScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedAlertChannelsScope,
					},
					APITokenFieldLimitedLinuxKvmHypervisorScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedLinuxKvmHypervisorScope,
					},
					APITokenFieldLimitedServiceLevelScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedServiceLevelScope,
					},
					APITokenFieldLimitedAiGatewayScope: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescLimitedAiGatewayScope,
					},

					// Additional permissions
					APITokenFieldCanConfigurePersonalAPITokens: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigurePersonalAPITokens,
					},
					APITokenFieldCanConfigureDatabaseManagement: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureDatabaseManagement,
					},
					APITokenFieldCanConfigureAutomationActions: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAutomationActions,
					},
					APITokenFieldCanConfigureAutomationPolicies: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAutomationPolicies,
					},
					APITokenFieldCanRunAutomationActions: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanRunAutomationActions,
					},
					APITokenFieldCanDeleteAutomationActionHistory: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanDeleteAutomationActionHistory,
					},
					APITokenFieldCanConfigureSyntheticTests: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureSyntheticTests,
					},
					APITokenFieldCanConfigureSyntheticLocations: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureSyntheticLocations,
					},
					APITokenFieldCanConfigureSyntheticCredentials: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureSyntheticCredentials,
					},
					APITokenFieldCanViewSyntheticTests: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewSyntheticTests,
					},
					APITokenFieldCanViewSyntheticLocations: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewSyntheticLocations,
					},
					APITokenFieldCanViewSyntheticTestResults: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewSyntheticTestResults,
					},
					APITokenFieldCanUseSyntheticCredentials: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanUseSyntheticCredentials,
					},
					APITokenFieldCanConfigureBizops: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureBizops,
					},
					APITokenFieldCanViewBusinessProcesses: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewBusinessProcesses,
					},
					APITokenFieldCanViewBusinessProcessDetails: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewBusinessProcessDetails,
					},
					APITokenFieldCanViewBusinessActivities: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewBusinessActivities,
					},
					APITokenFieldCanViewBizAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewBizAlerts,
					},
					APITokenFieldCanDeleteLogs: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanDeleteLogs,
					},
					APITokenFieldCanCreateHeapDump: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanCreateHeapDump,
					},
					APITokenFieldCanCreateThreadDump: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanCreateThreadDump,
					},
					APITokenFieldCanManuallyCloseIssue: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanManuallyCloseIssue,
					},
					APITokenFieldCanViewLogVolume: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanViewLogVolume,
					},
					APITokenFieldCanConfigureLogRetentionPeriod: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureLogRetentionPeriod,
					},
					APITokenFieldCanConfigureSubtraces: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureSubtraces,
					},
					APITokenFieldCanInvokeAlertChannel: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanInvokeAlertChannel,
					},
					APITokenFieldCanConfigureLlm: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureLlm,
					},
					APITokenFieldCanConfigureAiAgents: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureAiAgents,
					},
					APITokenFieldCanConfigureApdex: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureApdex,
					},
					APITokenFieldCanConfigureServiceLevelCorrectionWindows: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureServiceLevelCorrectionWindows,
					},
					APITokenFieldCanConfigureServiceLevelSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureServiceLevelSmartAlerts,
					},
					APITokenFieldCanConfigureServiceLevels: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: APITokenDescCanConfigureServiceLevels,
					},
				},
			},
			SchemaVersion: 2,
		},
	}
}

type apiTokenResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *apiTokenResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *apiTokenResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.APIToken] {
	return api.APITokens()
}

func (r *apiTokenResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(plan.SetAttribute(ctx, path.Root("internal_id"), types.StringValue(util.RandomID()))...)
	diags.Append(plan.SetAttribute(ctx, path.Root("access_granting_token"), types.StringValue(util.RandomID()))...)
	return diags
}

func (r *apiTokenResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiToken *restapi.APIToken) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API token
	model := APITokenModel{
		ID:                                       types.StringValue(apiToken.ID),
		AccessGrantingToken:                      types.StringValue(apiToken.AccessGrantingToken),
		InternalID:                               types.StringValue(apiToken.InternalID),
		Name:                                     types.StringValue(apiToken.Name),
		CanConfigureServiceMapping:               types.BoolValue(apiToken.CanConfigureServiceMapping),
		CanConfigureEumApplications:              types.BoolValue(apiToken.CanConfigureEumApplications),
		CanConfigureMobileAppMonitoring:          types.BoolValue(apiToken.CanConfigureMobileAppMonitoring),
		CanConfigureUsers:                        types.BoolValue(apiToken.CanConfigureUsers),
		CanInstallNewAgents:                      types.BoolValue(apiToken.CanInstallNewAgents),
		CanConfigureIntegrations:                 types.BoolValue(apiToken.CanConfigureIntegrations),
		CanConfigureEventsAndAlerts:              types.BoolValue(apiToken.CanConfigureEventsAndAlerts),
		CanConfigureMaintenanceWindows:           types.BoolValue(apiToken.CanConfigureMaintenanceWindows),
		CanConfigureApplicationSmartAlerts:       types.BoolValue(apiToken.CanConfigureApplicationSmartAlerts),
		CanConfigureWebsiteSmartAlerts:           types.BoolValue(apiToken.CanConfigureWebsiteSmartAlerts),
		CanConfigureMobileAppSmartAlerts:         types.BoolValue(apiToken.CanConfigureMobileAppSmartAlerts),
		CanConfigureAPITokens:                    types.BoolValue(apiToken.CanConfigureAPITokens),
		CanConfigureAgentRunMode:                 types.BoolValue(apiToken.CanConfigureAgentRunMode),
		CanViewAuditLog:                          types.BoolValue(apiToken.CanViewAuditLog),
		CanConfigureAgents:                       types.BoolValue(apiToken.CanConfigureAgents),
		CanConfigureAuthenticationMethods:        types.BoolValue(apiToken.CanConfigureAuthenticationMethods),
		CanConfigureApplications:                 types.BoolValue(apiToken.CanConfigureApplications),
		CanConfigureTeams:                        types.BoolValue(apiToken.CanConfigureTeams),
		CanConfigureReleases:                     types.BoolValue(apiToken.CanConfigureReleases),
		CanConfigureLogManagement:                types.BoolValue(apiToken.CanConfigureLogManagement),
		CanCreatePublicCustomDashboards:          types.BoolValue(apiToken.CanCreatePublicCustomDashboards),
		CanViewLogs:                              types.BoolValue(apiToken.CanViewLogs),
		CanViewTraceDetails:                      types.BoolValue(apiToken.CanViewTraceDetails),
		CanConfigureSessionSettings:              types.BoolValue(apiToken.CanConfigureSessionSettings),
		CanConfigureGlobalAlertPayload:           types.BoolValue(apiToken.CanConfigureGlobalAlertPayload),
		CanConfigureGlobalApplicationSmartAlerts: types.BoolValue(apiToken.CanConfigureGlobalApplicationSmartAlerts),
		CanConfigureGlobalSyntheticSmartAlerts:   types.BoolValue(apiToken.CanConfigureGlobalSyntheticSmartAlerts),
		CanConfigureGlobalInfraSmartAlerts:       types.BoolValue(apiToken.CanConfigureGlobalInfraSmartAlerts),
		CanConfigureGlobalLogSmartAlerts:         types.BoolValue(apiToken.CanConfigureGlobalLogSmartAlerts),
		CanViewAccountAndBillingInformation:      types.BoolValue(apiToken.CanViewAccountAndBillingInformation),
		CanEditAllAccessibleCustomDashboards:     types.BoolValue(apiToken.CanEditAllAccessibleCustomDashboards),

		// Scope limitations
		LimitedApplicationsScope:       types.BoolValue(apiToken.LimitedApplicationsScope),
		LimitedBizOpsScope:             types.BoolValue(apiToken.LimitedBizOpsScope),
		LimitedWebsitesScope:           types.BoolValue(apiToken.LimitedWebsitesScope),
		LimitedKubernetesScope:         types.BoolValue(apiToken.LimitedKubernetesScope),
		LimitedMobileAppsScope:         types.BoolValue(apiToken.LimitedMobileAppsScope),
		LimitedInfrastructureScope:     types.BoolValue(apiToken.LimitedInfrastructureScope),
		LimitedSyntheticsScope:         types.BoolValue(apiToken.LimitedSyntheticsScope),
		LimitedVsphereScope:            types.BoolValue(apiToken.LimitedVsphereScope),
		LimitedPhmcScope:               types.BoolValue(apiToken.LimitedPhmcScope),
		LimitedPvcScope:                types.BoolValue(apiToken.LimitedPvcScope),
		LimitedZhmcScope:               types.BoolValue(apiToken.LimitedZhmcScope),
		LimitedPcfScope:                types.BoolValue(apiToken.LimitedPcfScope),
		LimitedOpenstackScope:          types.BoolValue(apiToken.LimitedOpenstackScope),
		LimitedAutomationScope:         types.BoolValue(apiToken.LimitedAutomationScope),
		LimitedLogsScope:               types.BoolValue(apiToken.LimitedLogsScope),
		LimitedNutanixScope:            types.BoolValue(apiToken.LimitedNutanixScope),
		LimitedXenServerScope:          types.BoolValue(apiToken.LimitedXenServerScope),
		LimitedWindowsHypervisorScope:  types.BoolValue(apiToken.LimitedWindowsHypervisorScope),
		LimitedAlertChannelsScope:      types.BoolValue(apiToken.LimitedAlertChannelsScope),
		LimitedLinuxKvmHypervisorScope: types.BoolValue(apiToken.LimitedLinuxKvmHypervisorScope),
		LimitedServiceLevelScope:       types.BoolValue(apiToken.LimitedServiceLevelScope),
		LimitedAiGatewayScope:          types.BoolValue(apiToken.LimitedAiGatewayScope),

		// Additional permissions
		CanConfigurePersonalAPITokens:             types.BoolValue(apiToken.CanConfigurePersonalAPITokens),
		CanConfigureDatabaseManagement:            types.BoolValue(apiToken.CanConfigureDatabaseManagement),
		CanConfigureAutomationActions:             types.BoolValue(apiToken.CanConfigureAutomationActions),
		CanConfigureAutomationPolicies:            types.BoolValue(apiToken.CanConfigureAutomationPolicies),
		CanRunAutomationActions:                   types.BoolValue(apiToken.CanRunAutomationActions),
		CanDeleteAutomationActionHistory:          types.BoolValue(apiToken.CanDeleteAutomationActionHistory),
		CanConfigureSyntheticTests:                types.BoolValue(apiToken.CanConfigureSyntheticTests),
		CanConfigureSyntheticLocations:            types.BoolValue(apiToken.CanConfigureSyntheticLocations),
		CanConfigureSyntheticCredentials:          types.BoolValue(apiToken.CanConfigureSyntheticCredentials),
		CanViewSyntheticTests:                     types.BoolValue(apiToken.CanViewSyntheticTests),
		CanViewSyntheticLocations:                 types.BoolValue(apiToken.CanViewSyntheticLocations),
		CanViewSyntheticTestResults:               types.BoolValue(apiToken.CanViewSyntheticTestResults),
		CanUseSyntheticCredentials:                types.BoolValue(apiToken.CanUseSyntheticCredentials),
		CanConfigureBizops:                        types.BoolValue(apiToken.CanConfigureBizops),
		CanViewBusinessProcesses:                  types.BoolValue(apiToken.CanViewBusinessProcesses),
		CanViewBusinessProcessDetails:             types.BoolValue(apiToken.CanViewBusinessProcessDetails),
		CanViewBusinessActivities:                 types.BoolValue(apiToken.CanViewBusinessActivities),
		CanViewBizAlerts:                          types.BoolValue(apiToken.CanViewBizAlerts),
		CanDeleteLogs:                             types.BoolValue(apiToken.CanDeleteLogs),
		CanCreateHeapDump:                         types.BoolValue(apiToken.CanCreateHeapDump),
		CanCreateThreadDump:                       types.BoolValue(apiToken.CanCreateThreadDump),
		CanManuallyCloseIssue:                     types.BoolValue(apiToken.CanManuallyCloseIssue),
		CanViewLogVolume:                          types.BoolValue(apiToken.CanViewLogVolume),
		CanConfigureLogRetentionPeriod:            types.BoolValue(apiToken.CanConfigureLogRetentionPeriod),
		CanConfigureSubtraces:                     types.BoolValue(apiToken.CanConfigureSubtraces),
		CanInvokeAlertChannel:                     types.BoolValue(apiToken.CanInvokeAlertChannel),
		CanConfigureLlm:                           types.BoolValue(apiToken.CanConfigureLlm),
		CanConfigureAiAgents:                      types.BoolValue(apiToken.CanConfigureAiAgents),
		CanConfigureApdex:                         types.BoolValue(apiToken.CanConfigureApdex),
		CanConfigureServiceLevelCorrectionWindows: types.BoolValue(apiToken.CanConfigureServiceLevelCorrectionWindows),
		CanConfigureServiceLevelSmartAlerts:       types.BoolValue(apiToken.CanConfigureServiceLevelSmartAlerts),
		CanConfigureServiceLevels:                 types.BoolValue(apiToken.CanConfigureServiceLevels),
	}
	// Set the state with our populated model
	diags = state.Set(ctx, &model)
	return diags
}

func (r *apiTokenResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.APIToken, diag.Diagnostics) {
	var model APITokenModel
	var diags diag.Diagnostics

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	if diags.HasError() {
		return nil, diags
	}

	// Create API token object and populate it from the model
	apiToken := &restapi.APIToken{
		ID:                                       model.ID.ValueString(),
		AccessGrantingToken:                      model.AccessGrantingToken.ValueString(),
		InternalID:                               model.InternalID.ValueString(),
		Name:                                     model.Name.ValueString(),
		CanConfigureServiceMapping:               model.CanConfigureServiceMapping.ValueBool(),
		CanConfigureEumApplications:              model.CanConfigureEumApplications.ValueBool(),
		CanConfigureMobileAppMonitoring:          model.CanConfigureMobileAppMonitoring.ValueBool(),
		CanConfigureUsers:                        model.CanConfigureUsers.ValueBool(),
		CanInstallNewAgents:                      model.CanInstallNewAgents.ValueBool(),
		CanConfigureIntegrations:                 model.CanConfigureIntegrations.ValueBool(),
		CanConfigureEventsAndAlerts:              model.CanConfigureEventsAndAlerts.ValueBool(),
		CanConfigureMaintenanceWindows:           model.CanConfigureMaintenanceWindows.ValueBool(),
		CanConfigureApplicationSmartAlerts:       model.CanConfigureApplicationSmartAlerts.ValueBool(),
		CanConfigureWebsiteSmartAlerts:           model.CanConfigureWebsiteSmartAlerts.ValueBool(),
		CanConfigureMobileAppSmartAlerts:         model.CanConfigureMobileAppSmartAlerts.ValueBool(),
		CanConfigureAPITokens:                    model.CanConfigureAPITokens.ValueBool(),
		CanConfigureAgentRunMode:                 model.CanConfigureAgentRunMode.ValueBool(),
		CanViewAuditLog:                          model.CanViewAuditLog.ValueBool(),
		CanConfigureAgents:                       model.CanConfigureAgents.ValueBool(),
		CanConfigureAuthenticationMethods:        model.CanConfigureAuthenticationMethods.ValueBool(),
		CanConfigureApplications:                 model.CanConfigureApplications.ValueBool(),
		CanConfigureTeams:                        model.CanConfigureTeams.ValueBool(),
		CanConfigureReleases:                     model.CanConfigureReleases.ValueBool(),
		CanConfigureLogManagement:                model.CanConfigureLogManagement.ValueBool(),
		CanCreatePublicCustomDashboards:          model.CanCreatePublicCustomDashboards.ValueBool(),
		CanViewLogs:                              model.CanViewLogs.ValueBool(),
		CanViewTraceDetails:                      model.CanViewTraceDetails.ValueBool(),
		CanConfigureSessionSettings:              model.CanConfigureSessionSettings.ValueBool(),
		CanConfigureGlobalAlertPayload:           model.CanConfigureGlobalAlertPayload.ValueBool(),
		CanConfigureGlobalApplicationSmartAlerts: model.CanConfigureGlobalApplicationSmartAlerts.ValueBool(),
		CanConfigureGlobalSyntheticSmartAlerts:   model.CanConfigureGlobalSyntheticSmartAlerts.ValueBool(),
		CanConfigureGlobalInfraSmartAlerts:       model.CanConfigureGlobalInfraSmartAlerts.ValueBool(),
		CanConfigureGlobalLogSmartAlerts:         model.CanConfigureGlobalLogSmartAlerts.ValueBool(),
		CanViewAccountAndBillingInformation:      model.CanViewAccountAndBillingInformation.ValueBool(),
		CanEditAllAccessibleCustomDashboards:     model.CanEditAllAccessibleCustomDashboards.ValueBool(),

		// Scope limitations
		LimitedApplicationsScope:       model.LimitedApplicationsScope.ValueBool(),
		LimitedBizOpsScope:             model.LimitedBizOpsScope.ValueBool(),
		LimitedWebsitesScope:           model.LimitedWebsitesScope.ValueBool(),
		LimitedKubernetesScope:         model.LimitedKubernetesScope.ValueBool(),
		LimitedMobileAppsScope:         model.LimitedMobileAppsScope.ValueBool(),
		LimitedInfrastructureScope:     model.LimitedInfrastructureScope.ValueBool(),
		LimitedSyntheticsScope:         model.LimitedSyntheticsScope.ValueBool(),
		LimitedVsphereScope:            model.LimitedVsphereScope.ValueBool(),
		LimitedPhmcScope:               model.LimitedPhmcScope.ValueBool(),
		LimitedPvcScope:                model.LimitedPvcScope.ValueBool(),
		LimitedZhmcScope:               model.LimitedZhmcScope.ValueBool(),
		LimitedPcfScope:                model.LimitedPcfScope.ValueBool(),
		LimitedOpenstackScope:          model.LimitedOpenstackScope.ValueBool(),
		LimitedAutomationScope:         model.LimitedAutomationScope.ValueBool(),
		LimitedLogsScope:               model.LimitedLogsScope.ValueBool(),
		LimitedNutanixScope:            model.LimitedNutanixScope.ValueBool(),
		LimitedXenServerScope:          model.LimitedXenServerScope.ValueBool(),
		LimitedWindowsHypervisorScope:  model.LimitedWindowsHypervisorScope.ValueBool(),
		LimitedAlertChannelsScope:      model.LimitedAlertChannelsScope.ValueBool(),
		LimitedLinuxKvmHypervisorScope: model.LimitedLinuxKvmHypervisorScope.ValueBool(),
		LimitedServiceLevelScope:       model.LimitedServiceLevelScope.ValueBool(),
		LimitedAiGatewayScope:          model.LimitedAiGatewayScope.ValueBool(),

		// Additional permissions
		CanConfigurePersonalAPITokens:             model.CanConfigurePersonalAPITokens.ValueBool(),
		CanConfigureDatabaseManagement:            model.CanConfigureDatabaseManagement.ValueBool(),
		CanConfigureAutomationActions:             model.CanConfigureAutomationActions.ValueBool(),
		CanConfigureAutomationPolicies:            model.CanConfigureAutomationPolicies.ValueBool(),
		CanRunAutomationActions:                   model.CanRunAutomationActions.ValueBool(),
		CanDeleteAutomationActionHistory:          model.CanDeleteAutomationActionHistory.ValueBool(),
		CanConfigureSyntheticTests:                model.CanConfigureSyntheticTests.ValueBool(),
		CanConfigureSyntheticLocations:            model.CanConfigureSyntheticLocations.ValueBool(),
		CanConfigureSyntheticCredentials:          model.CanConfigureSyntheticCredentials.ValueBool(),
		CanViewSyntheticTests:                     model.CanViewSyntheticTests.ValueBool(),
		CanViewSyntheticLocations:                 model.CanViewSyntheticLocations.ValueBool(),
		CanViewSyntheticTestResults:               model.CanViewSyntheticTestResults.ValueBool(),
		CanUseSyntheticCredentials:                model.CanUseSyntheticCredentials.ValueBool(),
		CanConfigureBizops:                        model.CanConfigureBizops.ValueBool(),
		CanViewBusinessProcesses:                  model.CanViewBusinessProcesses.ValueBool(),
		CanViewBusinessProcessDetails:             model.CanViewBusinessProcessDetails.ValueBool(),
		CanViewBusinessActivities:                 model.CanViewBusinessActivities.ValueBool(),
		CanViewBizAlerts:                          model.CanViewBizAlerts.ValueBool(),
		CanDeleteLogs:                             model.CanDeleteLogs.ValueBool(),
		CanCreateHeapDump:                         model.CanCreateHeapDump.ValueBool(),
		CanCreateThreadDump:                       model.CanCreateThreadDump.ValueBool(),
		CanManuallyCloseIssue:                     model.CanManuallyCloseIssue.ValueBool(),
		CanViewLogVolume:                          model.CanViewLogVolume.ValueBool(),
		CanConfigureLogRetentionPeriod:            model.CanConfigureLogRetentionPeriod.ValueBool(),
		CanConfigureSubtraces:                     model.CanConfigureSubtraces.ValueBool(),
		CanInvokeAlertChannel:                     model.CanInvokeAlertChannel.ValueBool(),
		CanConfigureLlm:                           model.CanConfigureLlm.ValueBool(),
		CanConfigureAiAgents:                      model.CanConfigureAiAgents.ValueBool(),
		CanConfigureApdex:                         model.CanConfigureApdex.ValueBool(),
		CanConfigureServiceLevelCorrectionWindows: model.CanConfigureServiceLevelCorrectionWindows.ValueBool(),
		CanConfigureServiceLevelSmartAlerts:       model.CanConfigureServiceLevelSmartAlerts.ValueBool(),
		CanConfigureServiceLevels:                 model.CanConfigureServiceLevels.ValueBool(),
	}

	return apiToken, diags
}
