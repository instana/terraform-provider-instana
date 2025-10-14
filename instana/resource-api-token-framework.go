package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaAPITokenFramework the name of the terraform-provider-instana resource to manage API tokens
const ResourceInstanaAPITokenFramework = "api_token"

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
	LimitedGenAIScope              types.Bool `tfsdk:"limited_gen_ai_scope"`

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
func NewAPITokenResourceHandleFramework() ResourceHandleFramework[*restapi.APIToken] {
	return &apiTokenResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAPITokenFramework,
			Schema: schema.Schema{
				Description: "This resource manages API tokens in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the API token.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					APITokenFieldAccessGrantingToken: schema.StringAttribute{
						Computed:    true,
						Description: "The token used for the api Client used in the Authorization header to authenticate the client",
					},
					APITokenFieldInternalID: schema.StringAttribute{
						Computed:    true,
						Description: "The internal ID of the access token from the Instana platform",
					},
					APITokenFieldName: schema.StringAttribute{
						Required:    true,
						Description: "The name of the API token",
					},
					APITokenFieldCanConfigureServiceMapping: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure service mappings",
					},
					APITokenFieldCanConfigureEumApplications: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure End User Monitoring applications",
					},
					APITokenFieldCanConfigureMobileAppMonitoring: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Mobile App Monitoring",
					},
					APITokenFieldCanConfigureUsers: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure users",
					},
					APITokenFieldCanInstallNewAgents: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to install new agents",
					},
					APITokenFieldCanConfigureIntegrations: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure integrations",
					},
					APITokenFieldCanConfigureEventsAndAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Instana Events and Alerts",
					},
					APITokenFieldCanConfigureMaintenanceWindows: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Instana Maintenance Windows",
					},
					APITokenFieldCanConfigureApplicationSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Instana Application Smart Alerts",
					},
					APITokenFieldCanConfigureWebsiteSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Instana Website Smart Alerts",
					},
					APITokenFieldCanConfigureMobileAppSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Instana MobileApp Smart Alerts",
					},
					APITokenFieldCanConfigureAPITokens: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure API tokens",
					},
					APITokenFieldCanConfigureAgentRunMode: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure agent run mode",
					},
					APITokenFieldCanViewAuditLog: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view the audit log",
					},
					APITokenFieldCanConfigureAgents: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure agents",
					},
					APITokenFieldCanConfigureAuthenticationMethods: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure authentication methods",
					},
					APITokenFieldCanConfigureApplications: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure applications",
					},
					APITokenFieldCanConfigureTeams: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure teams (Groups)",
					},
					APITokenFieldCanConfigureReleases: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure releases",
					},
					APITokenFieldCanConfigureLogManagement: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure log management",
					},
					APITokenFieldCanCreatePublicCustomDashboards: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to create public custom dashboards",
					},
					APITokenFieldCanViewLogs: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view logs",
					},
					APITokenFieldCanViewTraceDetails: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view trace details",
					},
					APITokenFieldCanConfigureSessionSettings: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure session settings",
					},
					APITokenFieldCanConfigureGlobalAlertPayload: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure global alert payload",
					},
					APITokenFieldCanConfigureGlobalApplicationSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Global Application Smart Alerts",
					},
					APITokenFieldCanConfigureGlobalSyntheticSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Global Synthetic Smart Alerts",
					},
					APITokenFieldCanConfigureGlobalInfraSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Global Infra Smart Alerts",
					},
					APITokenFieldCanConfigureGlobalLogSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Global Log Smart Alerts",
					},
					APITokenFieldCanViewAccountAndBillingInformation: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view account and billing information",
					},
					APITokenFieldCanEditAllAccessibleCustomDashboards: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to edit all accessible custom dashboards",
					},

					// Scope limitations
					APITokenFieldLimitedApplicationsScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited applications scope",
					},
					APITokenFieldLimitedBizOpsScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited business operations scope",
					},
					APITokenFieldLimitedWebsitesScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited websites scope",
					},
					APITokenFieldLimitedKubernetesScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited kubernetes scope",
					},
					APITokenFieldLimitedMobileAppsScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited mobile apps scope",
					},
					APITokenFieldLimitedInfrastructureScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited infrastructure scope",
					},
					APITokenFieldLimitedSyntheticsScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited synthetics scope",
					},
					APITokenFieldLimitedVsphereScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited vsphere scope",
					},
					APITokenFieldLimitedPhmcScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited phmc scope",
					},
					APITokenFieldLimitedPvcScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited pvc scope",
					},
					APITokenFieldLimitedZhmcScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited zhmc scope",
					},
					APITokenFieldLimitedPcfScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited pcf scope",
					},
					APITokenFieldLimitedOpenstackScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited openstack scope",
					},
					APITokenFieldLimitedAutomationScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited automation scope",
					},
					APITokenFieldLimitedLogsScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited logs scope",
					},
					APITokenFieldLimitedNutanixScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited nutanix scope",
					},
					APITokenFieldLimitedXenServerScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited xen server scope",
					},
					APITokenFieldLimitedWindowsHypervisorScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited windows hypervisor scope",
					},
					APITokenFieldLimitedAlertChannelsScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited alert channels scope",
					},
					APITokenFieldLimitedLinuxKvmHypervisorScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited linux kvm hypervisor scope",
					},
					APITokenFieldLimitedServiceLevelScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited service level scope",
					},
					APITokenFieldLimitedAiGatewayScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited api gateway scope",
					},
					APITokenFieldLimitedGenAIScope: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token has limited gen gi scope",
					},

					// Additional permissions
					APITokenFieldCanConfigurePersonalAPITokens: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure personal API tokens",
					},
					APITokenFieldCanConfigureDatabaseManagement: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure database management",
					},
					APITokenFieldCanConfigureAutomationActions: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure automation actions",
					},
					APITokenFieldCanConfigureAutomationPolicies: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure automation policies",
					},
					APITokenFieldCanRunAutomationActions: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to run automation actions",
					},
					APITokenFieldCanDeleteAutomationActionHistory: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to delete automation action history",
					},
					APITokenFieldCanConfigureSyntheticTests: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure synthetic tests",
					},
					APITokenFieldCanConfigureSyntheticLocations: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure synthetic locations",
					},
					APITokenFieldCanConfigureSyntheticCredentials: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure synthetic credentials",
					},
					APITokenFieldCanViewSyntheticTests: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view synthetic tests",
					},
					APITokenFieldCanViewSyntheticLocations: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view synthetic locations",
					},
					APITokenFieldCanViewSyntheticTestResults: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view synthetic test results",
					},
					APITokenFieldCanUseSyntheticCredentials: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to use synthetic credentials",
					},
					APITokenFieldCanConfigureBizops: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure business operations",
					},
					APITokenFieldCanViewBusinessProcesses: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view business processes",
					},
					APITokenFieldCanViewBusinessProcessDetails: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view business process details",
					},
					APITokenFieldCanViewBusinessActivities: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view business activities",
					},
					APITokenFieldCanViewBizAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view business alerts",
					},
					APITokenFieldCanDeleteLogs: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to delete logs",
					},
					APITokenFieldCanCreateHeapDump: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to create heap dumps",
					},
					APITokenFieldCanCreateThreadDump: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to create thread dumps",
					},
					APITokenFieldCanManuallyCloseIssue: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to manually close issues",
					},
					APITokenFieldCanViewLogVolume: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to view log volume",
					},
					APITokenFieldCanConfigureLogRetentionPeriod: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure log retention period",
					},
					APITokenFieldCanConfigureSubtraces: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure subtraces",
					},
					APITokenFieldCanInvokeAlertChannel: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to invoke alert channels",
					},
					APITokenFieldCanConfigureLlm: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure LLM",
					},
					APITokenFieldCanConfigureAiAgents: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Ai Agents",
					},
					APITokenFieldCanConfigureApdex: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure Apdex",
					},
					APITokenFieldCanConfigureServiceLevelCorrectionWindows: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure service level Correction Windows",
					},
					APITokenFieldCanConfigureServiceLevelSmartAlerts: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure service level smart alerts",
					},
					APITokenFieldCanConfigureServiceLevels: schema.BoolAttribute{
						Optional:    true,
						Description: "Configures if the API token is allowed to configure service levels",
					},
				},
			},
			SchemaVersion: 2,
		},
	}
}

type apiTokenResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *apiTokenResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *apiTokenResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.APIToken] {
	return api.APITokens()
}

func (r *apiTokenResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(plan.SetAttribute(ctx, path.Root("internal_id"), types.StringValue(RandomID()))...)
	diags.Append(plan.SetAttribute(ctx, path.Root("access_granting_token"), types.StringValue(RandomID()))...)
	return diags
}

func (r *apiTokenResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, apiToken *restapi.APIToken) diag.Diagnostics {
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
		LimitedGenAIScope:              types.BoolValue(apiToken.LimitedGenAIScope),
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
		LimitedGenAIScope:              model.LimitedGenAIScope.ValueBool(),

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

// Made with Bob
