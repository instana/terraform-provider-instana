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
					APITokenFieldID: schema.StringAttribute{
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

// ============================================================================
// Resource Implementation
// ============================================================================

type apiTokenResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

// MetaData returns the resource metadata
func (r *apiTokenResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

// GetRestResource returns the REST resource for API tokens
func (r *apiTokenResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.APIToken] {
	return api.APITokens()
}

// SetComputedFields sets computed fields (internal_id and access_granting_token) in the plan
func (r *apiTokenResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(plan.SetAttribute(ctx, path.Root(APITokenFieldInternalID), types.StringValue(util.RandomID()))...)
	diags.Append(plan.SetAttribute(ctx, path.Root(APITokenFieldAccessGrantingToken), types.StringValue(util.RandomID()))...)
	return diags
}

// ============================================================================
// State Management
// ============================================================================

// UpdateState updates the Terraform state with the API token data from the API
func (r *apiTokenResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiToken *restapi.APIToken) diag.Diagnostics {
	// Create base model with core fields
	model := APITokenModel{
		ID:                  types.StringValue(apiToken.ID),
		AccessGrantingToken: types.StringValue(apiToken.AccessGrantingToken),
		InternalID:          types.StringValue(apiToken.InternalID),
		Name:                types.StringValue(apiToken.Name),
	}

	// Map permissions
	r.mapPermissionsToModel(apiToken, &model)

	// Map scope limitations
	r.mapScopeLimitationsToModel(apiToken, &model)

	// Map additional permissions
	r.mapAdditionalPermissionsToModel(apiToken, &model)

	// Set the state with our populated model
	return state.Set(ctx, &model)
}

// mapPermissionsToModel maps basic permissions from API to model
func (r *apiTokenResourceFramework) mapPermissionsToModel(apiToken *restapi.APIToken, model *APITokenModel) {
	model.CanConfigureServiceMapping = types.BoolValue(apiToken.CanConfigureServiceMapping)
	model.CanConfigureEumApplications = types.BoolValue(apiToken.CanConfigureEumApplications)
	model.CanConfigureMobileAppMonitoring = types.BoolValue(apiToken.CanConfigureMobileAppMonitoring)
	model.CanConfigureUsers = types.BoolValue(apiToken.CanConfigureUsers)
	model.CanInstallNewAgents = types.BoolValue(apiToken.CanInstallNewAgents)
	model.CanConfigureIntegrations = types.BoolValue(apiToken.CanConfigureIntegrations)
	model.CanConfigureEventsAndAlerts = types.BoolValue(apiToken.CanConfigureEventsAndAlerts)
	model.CanConfigureMaintenanceWindows = types.BoolValue(apiToken.CanConfigureMaintenanceWindows)
	model.CanConfigureApplicationSmartAlerts = types.BoolValue(apiToken.CanConfigureApplicationSmartAlerts)
	model.CanConfigureWebsiteSmartAlerts = types.BoolValue(apiToken.CanConfigureWebsiteSmartAlerts)
	model.CanConfigureMobileAppSmartAlerts = types.BoolValue(apiToken.CanConfigureMobileAppSmartAlerts)
	model.CanConfigureAPITokens = types.BoolValue(apiToken.CanConfigureAPITokens)
	model.CanConfigureAgentRunMode = types.BoolValue(apiToken.CanConfigureAgentRunMode)
	model.CanViewAuditLog = types.BoolValue(apiToken.CanViewAuditLog)
	model.CanConfigureAgents = types.BoolValue(apiToken.CanConfigureAgents)
	model.CanConfigureAuthenticationMethods = types.BoolValue(apiToken.CanConfigureAuthenticationMethods)
	model.CanConfigureApplications = types.BoolValue(apiToken.CanConfigureApplications)
	model.CanConfigureTeams = types.BoolValue(apiToken.CanConfigureTeams)
	model.CanConfigureReleases = types.BoolValue(apiToken.CanConfigureReleases)
	model.CanConfigureLogManagement = types.BoolValue(apiToken.CanConfigureLogManagement)
	model.CanCreatePublicCustomDashboards = types.BoolValue(apiToken.CanCreatePublicCustomDashboards)
	model.CanViewLogs = types.BoolValue(apiToken.CanViewLogs)
	model.CanViewTraceDetails = types.BoolValue(apiToken.CanViewTraceDetails)
	model.CanConfigureSessionSettings = types.BoolValue(apiToken.CanConfigureSessionSettings)
	model.CanConfigureGlobalAlertPayload = types.BoolValue(apiToken.CanConfigureGlobalAlertPayload)
	model.CanConfigureGlobalApplicationSmartAlerts = types.BoolValue(apiToken.CanConfigureGlobalApplicationSmartAlerts)
	model.CanConfigureGlobalSyntheticSmartAlerts = types.BoolValue(apiToken.CanConfigureGlobalSyntheticSmartAlerts)
	model.CanConfigureGlobalInfraSmartAlerts = types.BoolValue(apiToken.CanConfigureGlobalInfraSmartAlerts)
	model.CanConfigureGlobalLogSmartAlerts = types.BoolValue(apiToken.CanConfigureGlobalLogSmartAlerts)
	model.CanViewAccountAndBillingInformation = types.BoolValue(apiToken.CanViewAccountAndBillingInformation)
	model.CanEditAllAccessibleCustomDashboards = types.BoolValue(apiToken.CanEditAllAccessibleCustomDashboards)
}

// mapScopeLimitationsToModel maps scope limitations from API to model
func (r *apiTokenResourceFramework) mapScopeLimitationsToModel(apiToken *restapi.APIToken, model *APITokenModel) {
	model.LimitedApplicationsScope = types.BoolValue(apiToken.LimitedApplicationsScope)
	model.LimitedBizOpsScope = types.BoolValue(apiToken.LimitedBizOpsScope)
	model.LimitedWebsitesScope = types.BoolValue(apiToken.LimitedWebsitesScope)
	model.LimitedKubernetesScope = types.BoolValue(apiToken.LimitedKubernetesScope)
	model.LimitedMobileAppsScope = types.BoolValue(apiToken.LimitedMobileAppsScope)
	model.LimitedInfrastructureScope = types.BoolValue(apiToken.LimitedInfrastructureScope)
	model.LimitedSyntheticsScope = types.BoolValue(apiToken.LimitedSyntheticsScope)
	model.LimitedVsphereScope = types.BoolValue(apiToken.LimitedVsphereScope)
	model.LimitedPhmcScope = types.BoolValue(apiToken.LimitedPhmcScope)
	model.LimitedPvcScope = types.BoolValue(apiToken.LimitedPvcScope)
	model.LimitedZhmcScope = types.BoolValue(apiToken.LimitedZhmcScope)
	model.LimitedPcfScope = types.BoolValue(apiToken.LimitedPcfScope)
	model.LimitedOpenstackScope = types.BoolValue(apiToken.LimitedOpenstackScope)
	model.LimitedAutomationScope = types.BoolValue(apiToken.LimitedAutomationScope)
	model.LimitedLogsScope = types.BoolValue(apiToken.LimitedLogsScope)
	model.LimitedNutanixScope = types.BoolValue(apiToken.LimitedNutanixScope)
	model.LimitedXenServerScope = types.BoolValue(apiToken.LimitedXenServerScope)
	model.LimitedWindowsHypervisorScope = types.BoolValue(apiToken.LimitedWindowsHypervisorScope)
	model.LimitedAlertChannelsScope = types.BoolValue(apiToken.LimitedAlertChannelsScope)
	model.LimitedLinuxKvmHypervisorScope = types.BoolValue(apiToken.LimitedLinuxKvmHypervisorScope)
	model.LimitedServiceLevelScope = types.BoolValue(apiToken.LimitedServiceLevelScope)
	model.LimitedAiGatewayScope = types.BoolValue(apiToken.LimitedAiGatewayScope)
}

// mapAdditionalPermissionsToModel maps additional permissions from API to model
func (r *apiTokenResourceFramework) mapAdditionalPermissionsToModel(apiToken *restapi.APIToken, model *APITokenModel) {
	model.CanConfigurePersonalAPITokens = types.BoolValue(apiToken.CanConfigurePersonalAPITokens)
	model.CanConfigureDatabaseManagement = types.BoolValue(apiToken.CanConfigureDatabaseManagement)
	model.CanConfigureAutomationActions = types.BoolValue(apiToken.CanConfigureAutomationActions)
	model.CanConfigureAutomationPolicies = types.BoolValue(apiToken.CanConfigureAutomationPolicies)
	model.CanRunAutomationActions = types.BoolValue(apiToken.CanRunAutomationActions)
	model.CanDeleteAutomationActionHistory = types.BoolValue(apiToken.CanDeleteAutomationActionHistory)
	model.CanConfigureSyntheticTests = types.BoolValue(apiToken.CanConfigureSyntheticTests)
	model.CanConfigureSyntheticLocations = types.BoolValue(apiToken.CanConfigureSyntheticLocations)
	model.CanConfigureSyntheticCredentials = types.BoolValue(apiToken.CanConfigureSyntheticCredentials)
	model.CanViewSyntheticTests = types.BoolValue(apiToken.CanViewSyntheticTests)
	model.CanViewSyntheticLocations = types.BoolValue(apiToken.CanViewSyntheticLocations)
	model.CanViewSyntheticTestResults = types.BoolValue(apiToken.CanViewSyntheticTestResults)
	model.CanUseSyntheticCredentials = types.BoolValue(apiToken.CanUseSyntheticCredentials)
	model.CanConfigureBizops = types.BoolValue(apiToken.CanConfigureBizops)
	model.CanViewBusinessProcesses = types.BoolValue(apiToken.CanViewBusinessProcesses)
	model.CanViewBusinessProcessDetails = types.BoolValue(apiToken.CanViewBusinessProcessDetails)
	model.CanViewBusinessActivities = types.BoolValue(apiToken.CanViewBusinessActivities)
	model.CanViewBizAlerts = types.BoolValue(apiToken.CanViewBizAlerts)
	model.CanDeleteLogs = types.BoolValue(apiToken.CanDeleteLogs)
	model.CanCreateHeapDump = types.BoolValue(apiToken.CanCreateHeapDump)
	model.CanCreateThreadDump = types.BoolValue(apiToken.CanCreateThreadDump)
	model.CanManuallyCloseIssue = types.BoolValue(apiToken.CanManuallyCloseIssue)
	model.CanViewLogVolume = types.BoolValue(apiToken.CanViewLogVolume)
	model.CanConfigureLogRetentionPeriod = types.BoolValue(apiToken.CanConfigureLogRetentionPeriod)
	model.CanConfigureSubtraces = types.BoolValue(apiToken.CanConfigureSubtraces)
	model.CanInvokeAlertChannel = types.BoolValue(apiToken.CanInvokeAlertChannel)
	model.CanConfigureLlm = types.BoolValue(apiToken.CanConfigureLlm)
	model.CanConfigureAiAgents = types.BoolValue(apiToken.CanConfigureAiAgents)
	model.CanConfigureApdex = types.BoolValue(apiToken.CanConfigureApdex)
	model.CanConfigureServiceLevelCorrectionWindows = types.BoolValue(apiToken.CanConfigureServiceLevelCorrectionWindows)
	model.CanConfigureServiceLevelSmartAlerts = types.BoolValue(apiToken.CanConfigureServiceLevelSmartAlerts)
	model.CanConfigureServiceLevels = types.BoolValue(apiToken.CanConfigureServiceLevels)
}

// ============================================================================
// State to API Mapping
// ============================================================================

// MapStateToDataObject converts Terraform state to API object
func (r *apiTokenResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.APIToken, diag.Diagnostics) {
	// Get model from plan or state
	model, diags := r.getAPITokenModelFromPlanOrState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	// Create base API token with core fields
	apiToken := &restapi.APIToken{
		ID:                  model.ID.ValueString(),
		AccessGrantingToken: model.AccessGrantingToken.ValueString(),
		InternalID:          model.InternalID.ValueString(),
		Name:                model.Name.ValueString(),
	}

	// Map permissions
	r.mapPermissionsFromModel(model, apiToken)

	// Map scope limitations
	r.mapScopeLimitationsFromModel(model, apiToken)

	// Map additional permissions
	r.mapAdditionalPermissionsFromModel(model, apiToken)

	return apiToken, diags
}

// getAPITokenModelFromPlanOrState retrieves the model from either plan or state
func (r *apiTokenResourceFramework) getAPITokenModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (APITokenModel, diag.Diagnostics) {
	var model APITokenModel
	var diags diag.Diagnostics

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// mapPermissionsFromModel maps basic permissions from model to API object
func (r *apiTokenResourceFramework) mapPermissionsFromModel(model APITokenModel, apiToken *restapi.APIToken) {
	apiToken.CanConfigureServiceMapping = model.CanConfigureServiceMapping.ValueBool()
	apiToken.CanConfigureEumApplications = model.CanConfigureEumApplications.ValueBool()
	apiToken.CanConfigureMobileAppMonitoring = model.CanConfigureMobileAppMonitoring.ValueBool()
	apiToken.CanConfigureUsers = model.CanConfigureUsers.ValueBool()
	apiToken.CanInstallNewAgents = model.CanInstallNewAgents.ValueBool()
	apiToken.CanConfigureIntegrations = model.CanConfigureIntegrations.ValueBool()
	apiToken.CanConfigureEventsAndAlerts = model.CanConfigureEventsAndAlerts.ValueBool()
	apiToken.CanConfigureMaintenanceWindows = model.CanConfigureMaintenanceWindows.ValueBool()
	apiToken.CanConfigureApplicationSmartAlerts = model.CanConfigureApplicationSmartAlerts.ValueBool()
	apiToken.CanConfigureWebsiteSmartAlerts = model.CanConfigureWebsiteSmartAlerts.ValueBool()
	apiToken.CanConfigureMobileAppSmartAlerts = model.CanConfigureMobileAppSmartAlerts.ValueBool()
	apiToken.CanConfigureAPITokens = model.CanConfigureAPITokens.ValueBool()
	apiToken.CanConfigureAgentRunMode = model.CanConfigureAgentRunMode.ValueBool()
	apiToken.CanViewAuditLog = model.CanViewAuditLog.ValueBool()
	apiToken.CanConfigureAgents = model.CanConfigureAgents.ValueBool()
	apiToken.CanConfigureAuthenticationMethods = model.CanConfigureAuthenticationMethods.ValueBool()
	apiToken.CanConfigureApplications = model.CanConfigureApplications.ValueBool()
	apiToken.CanConfigureTeams = model.CanConfigureTeams.ValueBool()
	apiToken.CanConfigureReleases = model.CanConfigureReleases.ValueBool()
	apiToken.CanConfigureLogManagement = model.CanConfigureLogManagement.ValueBool()
	apiToken.CanCreatePublicCustomDashboards = model.CanCreatePublicCustomDashboards.ValueBool()
	apiToken.CanViewLogs = model.CanViewLogs.ValueBool()
	apiToken.CanViewTraceDetails = model.CanViewTraceDetails.ValueBool()
	apiToken.CanConfigureSessionSettings = model.CanConfigureSessionSettings.ValueBool()
	apiToken.CanConfigureGlobalAlertPayload = model.CanConfigureGlobalAlertPayload.ValueBool()
	apiToken.CanConfigureGlobalApplicationSmartAlerts = model.CanConfigureGlobalApplicationSmartAlerts.ValueBool()
	apiToken.CanConfigureGlobalSyntheticSmartAlerts = model.CanConfigureGlobalSyntheticSmartAlerts.ValueBool()
	apiToken.CanConfigureGlobalInfraSmartAlerts = model.CanConfigureGlobalInfraSmartAlerts.ValueBool()
	apiToken.CanConfigureGlobalLogSmartAlerts = model.CanConfigureGlobalLogSmartAlerts.ValueBool()
	apiToken.CanViewAccountAndBillingInformation = model.CanViewAccountAndBillingInformation.ValueBool()
	apiToken.CanEditAllAccessibleCustomDashboards = model.CanEditAllAccessibleCustomDashboards.ValueBool()
}

// mapScopeLimitationsFromModel maps scope limitations from model to API object
func (r *apiTokenResourceFramework) mapScopeLimitationsFromModel(model APITokenModel, apiToken *restapi.APIToken) {
	apiToken.LimitedApplicationsScope = model.LimitedApplicationsScope.ValueBool()
	apiToken.LimitedBizOpsScope = model.LimitedBizOpsScope.ValueBool()
	apiToken.LimitedWebsitesScope = model.LimitedWebsitesScope.ValueBool()
	apiToken.LimitedKubernetesScope = model.LimitedKubernetesScope.ValueBool()
	apiToken.LimitedMobileAppsScope = model.LimitedMobileAppsScope.ValueBool()
	apiToken.LimitedInfrastructureScope = model.LimitedInfrastructureScope.ValueBool()
	apiToken.LimitedSyntheticsScope = model.LimitedSyntheticsScope.ValueBool()
	apiToken.LimitedVsphereScope = model.LimitedVsphereScope.ValueBool()
	apiToken.LimitedPhmcScope = model.LimitedPhmcScope.ValueBool()
	apiToken.LimitedPvcScope = model.LimitedPvcScope.ValueBool()
	apiToken.LimitedZhmcScope = model.LimitedZhmcScope.ValueBool()
	apiToken.LimitedPcfScope = model.LimitedPcfScope.ValueBool()
	apiToken.LimitedOpenstackScope = model.LimitedOpenstackScope.ValueBool()
	apiToken.LimitedAutomationScope = model.LimitedAutomationScope.ValueBool()
	apiToken.LimitedLogsScope = model.LimitedLogsScope.ValueBool()
	apiToken.LimitedNutanixScope = model.LimitedNutanixScope.ValueBool()
	apiToken.LimitedXenServerScope = model.LimitedXenServerScope.ValueBool()
	apiToken.LimitedWindowsHypervisorScope = model.LimitedWindowsHypervisorScope.ValueBool()
	apiToken.LimitedAlertChannelsScope = model.LimitedAlertChannelsScope.ValueBool()
	apiToken.LimitedLinuxKvmHypervisorScope = model.LimitedLinuxKvmHypervisorScope.ValueBool()
	apiToken.LimitedServiceLevelScope = model.LimitedServiceLevelScope.ValueBool()
	apiToken.LimitedAiGatewayScope = model.LimitedAiGatewayScope.ValueBool()
}

// mapAdditionalPermissionsFromModel maps additional permissions from model to API object
func (r *apiTokenResourceFramework) mapAdditionalPermissionsFromModel(model APITokenModel, apiToken *restapi.APIToken) {
	apiToken.CanConfigurePersonalAPITokens = model.CanConfigurePersonalAPITokens.ValueBool()
	apiToken.CanConfigureDatabaseManagement = model.CanConfigureDatabaseManagement.ValueBool()
	apiToken.CanConfigureAutomationActions = model.CanConfigureAutomationActions.ValueBool()
	apiToken.CanConfigureAutomationPolicies = model.CanConfigureAutomationPolicies.ValueBool()
	apiToken.CanRunAutomationActions = model.CanRunAutomationActions.ValueBool()
	apiToken.CanDeleteAutomationActionHistory = model.CanDeleteAutomationActionHistory.ValueBool()
	apiToken.CanConfigureSyntheticTests = model.CanConfigureSyntheticTests.ValueBool()
	apiToken.CanConfigureSyntheticLocations = model.CanConfigureSyntheticLocations.ValueBool()
	apiToken.CanConfigureSyntheticCredentials = model.CanConfigureSyntheticCredentials.ValueBool()
	apiToken.CanViewSyntheticTests = model.CanViewSyntheticTests.ValueBool()
	apiToken.CanViewSyntheticLocations = model.CanViewSyntheticLocations.ValueBool()
	apiToken.CanViewSyntheticTestResults = model.CanViewSyntheticTestResults.ValueBool()
	apiToken.CanUseSyntheticCredentials = model.CanUseSyntheticCredentials.ValueBool()
	apiToken.CanConfigureBizops = model.CanConfigureBizops.ValueBool()
	apiToken.CanViewBusinessProcesses = model.CanViewBusinessProcesses.ValueBool()
	apiToken.CanViewBusinessProcessDetails = model.CanViewBusinessProcessDetails.ValueBool()
	apiToken.CanViewBusinessActivities = model.CanViewBusinessActivities.ValueBool()
	apiToken.CanViewBizAlerts = model.CanViewBizAlerts.ValueBool()
	apiToken.CanDeleteLogs = model.CanDeleteLogs.ValueBool()
	apiToken.CanCreateHeapDump = model.CanCreateHeapDump.ValueBool()
	apiToken.CanCreateThreadDump = model.CanCreateThreadDump.ValueBool()
	apiToken.CanManuallyCloseIssue = model.CanManuallyCloseIssue.ValueBool()
	apiToken.CanViewLogVolume = model.CanViewLogVolume.ValueBool()
	apiToken.CanConfigureLogRetentionPeriod = model.CanConfigureLogRetentionPeriod.ValueBool()
	apiToken.CanConfigureSubtraces = model.CanConfigureSubtraces.ValueBool()
	apiToken.CanInvokeAlertChannel = model.CanInvokeAlertChannel.ValueBool()
	apiToken.CanConfigureLlm = model.CanConfigureLlm.ValueBool()
	apiToken.CanConfigureAiAgents = model.CanConfigureAiAgents.ValueBool()
	apiToken.CanConfigureApdex = model.CanConfigureApdex.ValueBool()
	apiToken.CanConfigureServiceLevelCorrectionWindows = model.CanConfigureServiceLevelCorrectionWindows.ValueBool()
	apiToken.CanConfigureServiceLevelSmartAlerts = model.CanConfigureServiceLevelSmartAlerts.ValueBool()
	apiToken.CanConfigureServiceLevels = model.CanConfigureServiceLevels.ValueBool()
}
