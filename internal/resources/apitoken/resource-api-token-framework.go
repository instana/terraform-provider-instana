package apitoken

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
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

// NewAPITokenResourceHandleFramework creates the resource handle for API Tokens
func NewAPITokenResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.APIToken] {
	internalIDFieldName := APITokenFieldInternalID
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
			SkipIDGeneration: true,
			SchemaVersion:    2,
			ResourceIDField:  &internalIDFieldName,
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
