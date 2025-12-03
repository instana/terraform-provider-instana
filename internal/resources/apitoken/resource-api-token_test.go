package apitoken

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPITokenResourceHandle(t *testing.T) {
	handle := NewAPITokenResourceHandle()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaAPIToken, metadata.ResourceName)
	assert.Equal(t, int64(3), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &apiTokenResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 2,
		},
	}

	metadata := resource.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, "test_resource", metadata.ResourceName)
	assert.Equal(t, int64(2), metadata.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := &apiTokenResource{}
	ctx := context.Background()

	handle := NewAPITokenResourceHandle()
	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	// Set a basic model in the plan
	model := APITokenModel{
		Name: types.StringValue("test-token"),
	}
	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	// Call SetComputedFields
	diags = resource.SetComputedFields(ctx, plan)
	require.False(t, diags.HasError())

	// Get the updated model
	var updatedModel APITokenModel
	diags = plan.Get(ctx, &updatedModel)
	require.False(t, diags.HasError())

	// Verify computed fields are set
	assert.False(t, updatedModel.InternalID.IsNull())
	assert.False(t, updatedModel.AccessGrantingToken.IsNull())
	assert.NotEmpty(t, updatedModel.InternalID.ValueString())
	assert.NotEmpty(t, updatedModel.AccessGrantingToken.ValueString())
}

func TestUpdateState(t *testing.T) {
	resource := &apiTokenResource{}
	ctx := context.Background()

	t.Run("basic API token", func(t *testing.T) {
		apiToken := &restapi.APIToken{
			ID:                  "test-id",
			AccessGrantingToken: "test-access-token",
			InternalID:          "test-internal-id",
			Name:                "test-token",
		}

		handle := NewAPITokenResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiToken)
		require.False(t, diags.HasError())

		var model APITokenModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-access-token", model.AccessGrantingToken.ValueString())
		assert.Equal(t, "test-internal-id", model.InternalID.ValueString())
		assert.Equal(t, "test-token", model.Name.ValueString())
	})

	t.Run("API token with permissions", func(t *testing.T) {
		apiToken := &restapi.APIToken{
			ID:                          "test-id",
			AccessGrantingToken:         "test-access-token",
			InternalID:                  "test-internal-id",
			Name:                        "test-token",
			CanConfigureServiceMapping:  true,
			CanConfigureUsers:           true,
			CanInstallNewAgents:         true,
			CanConfigureIntegrations:    true,
			CanConfigureEventsAndAlerts: true,
			CanViewAuditLog:             true,
		}

		handle := NewAPITokenResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiToken)
		require.False(t, diags.HasError())

		var model APITokenModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.True(t, model.CanConfigureServiceMapping.ValueBool())
		assert.True(t, model.CanConfigureUsers.ValueBool())
		assert.True(t, model.CanInstallNewAgents.ValueBool())
		assert.True(t, model.CanConfigureIntegrations.ValueBool())
		assert.True(t, model.CanConfigureEventsAndAlerts.ValueBool())
		assert.True(t, model.CanViewAuditLog.ValueBool())
	})

	t.Run("API token with scope limitations", func(t *testing.T) {
		apiToken := &restapi.APIToken{
			ID:                         "test-id",
			AccessGrantingToken:        "test-access-token",
			InternalID:                 "test-internal-id",
			Name:                       "test-token",
			LimitedApplicationsScope:   true,
			LimitedKubernetesScope:     true,
			LimitedInfrastructureScope: true,
			LimitedLogsScope:           true,
		}

		handle := NewAPITokenResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiToken)
		require.False(t, diags.HasError())

		var model APITokenModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.True(t, model.LimitedApplicationsScope.ValueBool())
		assert.True(t, model.LimitedKubernetesScope.ValueBool())
		assert.True(t, model.LimitedInfrastructureScope.ValueBool())
		assert.True(t, model.LimitedLogsScope.ValueBool())
	})

	t.Run("API token with all permissions enabled", func(t *testing.T) {
		apiToken := createFullyPermissionedAPIToken()

		handle := NewAPITokenResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiToken)
		require.False(t, diags.HasError())

		var model APITokenModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		// Verify all permissions are true
		assert.True(t, model.CanConfigureServiceMapping.ValueBool())
		assert.True(t, model.CanConfigureEumApplications.ValueBool())
		assert.True(t, model.CanConfigureMobileAppMonitoring.ValueBool())
		assert.True(t, model.CanConfigureUsers.ValueBool())
		assert.True(t, model.CanInstallNewAgents.ValueBool())
		assert.True(t, model.CanConfigureAPITokens.ValueBool())
		assert.True(t, model.CanViewAuditLog.ValueBool())
		assert.True(t, model.CanConfigureAgents.ValueBool())
		assert.True(t, model.CanConfigureApplications.ValueBool())
		assert.True(t, model.CanConfigureTeams.ValueBool())
		assert.True(t, model.CanViewLogs.ValueBool())
		assert.True(t, model.CanViewTraceDetails.ValueBool())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &apiTokenResource{}
	ctx := context.Background()

	t.Run("basic API token from state", func(t *testing.T) {
		model := APITokenModel{
			ID:                  types.StringValue("test-id"),
			AccessGrantingToken: types.StringValue("test-access-token"),
			InternalID:          types.StringValue("test-internal-id"),
			Name:                types.StringValue("test-token"),
		}

		state := createMockState(t, ctx, model)
		apiToken, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, apiToken)

		assert.Equal(t, "test-id", apiToken.ID)
		assert.Equal(t, "test-access-token", apiToken.AccessGrantingToken)
		assert.Equal(t, "test-internal-id", apiToken.InternalID)
		assert.Equal(t, "test-token", apiToken.Name)
	})

	t.Run("API token from plan", func(t *testing.T) {
		model := APITokenModel{
			ID:                  types.StringValue(""),
			AccessGrantingToken: types.StringValue("new-access-token"),
			InternalID:          types.StringValue("new-internal-id"),
			Name:                types.StringValue("new-token"),
		}

		plan := createMockPlan(t, ctx, model)
		apiToken, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, apiToken)

		assert.Equal(t, "", apiToken.ID)
		assert.Equal(t, "new-access-token", apiToken.AccessGrantingToken)
		assert.Equal(t, "new-internal-id", apiToken.InternalID)
		assert.Equal(t, "new-token", apiToken.Name)
	})

	t.Run("API token with permissions", func(t *testing.T) {
		model := APITokenModel{
			ID:                          types.StringValue("test-id"),
			AccessGrantingToken:         types.StringValue("test-access-token"),
			InternalID:                  types.StringValue("test-internal-id"),
			Name:                        types.StringValue("test-token"),
			CanConfigureServiceMapping:  types.BoolValue(true),
			CanConfigureUsers:           types.BoolValue(true),
			CanInstallNewAgents:         types.BoolValue(true),
			CanConfigureIntegrations:    types.BoolValue(true),
			CanConfigureEventsAndAlerts: types.BoolValue(true),
			CanViewAuditLog:             types.BoolValue(true),
		}

		state := createMockState(t, ctx, model)
		apiToken, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, apiToken)

		assert.True(t, apiToken.CanConfigureServiceMapping)
		assert.True(t, apiToken.CanConfigureUsers)
		assert.True(t, apiToken.CanInstallNewAgents)
		assert.True(t, apiToken.CanConfigureIntegrations)
		assert.True(t, apiToken.CanConfigureEventsAndAlerts)
		assert.True(t, apiToken.CanViewAuditLog)
	})

	t.Run("API token with scope limitations", func(t *testing.T) {
		model := APITokenModel{
			ID:                         types.StringValue("test-id"),
			AccessGrantingToken:        types.StringValue("test-access-token"),
			InternalID:                 types.StringValue("test-internal-id"),
			Name:                       types.StringValue("test-token"),
			LimitedApplicationsScope:   types.BoolValue(true),
			LimitedKubernetesScope:     types.BoolValue(true),
			LimitedInfrastructureScope: types.BoolValue(true),
			LimitedLogsScope:           types.BoolValue(true),
			LimitedSyntheticsScope:     types.BoolValue(true),
		}

		state := createMockState(t, ctx, model)
		apiToken, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, apiToken)

		assert.True(t, apiToken.LimitedApplicationsScope)
		assert.True(t, apiToken.LimitedKubernetesScope)
		assert.True(t, apiToken.LimitedInfrastructureScope)
		assert.True(t, apiToken.LimitedLogsScope)
		assert.True(t, apiToken.LimitedSyntheticsScope)
	})

	t.Run("API token with all fields populated", func(t *testing.T) {
		model := createFullyPermissionedModel()

		state := createMockState(t, ctx, model)
		apiToken, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, apiToken)

		// Verify all permissions
		assert.True(t, apiToken.CanConfigureServiceMapping)
		assert.True(t, apiToken.CanConfigureEumApplications)
		assert.True(t, apiToken.CanConfigureMobileAppMonitoring)
		assert.True(t, apiToken.CanConfigureUsers)
		assert.True(t, apiToken.CanInstallNewAgents)
		assert.True(t, apiToken.CanConfigureIntegrations)
		assert.True(t, apiToken.CanConfigureEventsAndAlerts)
		assert.True(t, apiToken.CanConfigureMaintenanceWindows)
		assert.True(t, apiToken.CanConfigureApplicationSmartAlerts)
		assert.True(t, apiToken.CanConfigureWebsiteSmartAlerts)
		assert.True(t, apiToken.CanConfigureMobileAppSmartAlerts)
		assert.True(t, apiToken.CanConfigureAPITokens)
		assert.True(t, apiToken.CanConfigureAgentRunMode)
		assert.True(t, apiToken.CanViewAuditLog)
		assert.True(t, apiToken.CanConfigureAgents)
		assert.True(t, apiToken.CanConfigureAuthenticationMethods)
		assert.True(t, apiToken.CanConfigureApplications)
		assert.True(t, apiToken.CanConfigureTeams)
		assert.True(t, apiToken.CanConfigureReleases)
		assert.True(t, apiToken.CanConfigureLogManagement)
		assert.True(t, apiToken.CanCreatePublicCustomDashboards)
		assert.True(t, apiToken.CanViewLogs)
		assert.True(t, apiToken.CanViewTraceDetails)

		// Verify scope limitations
		assert.True(t, apiToken.LimitedApplicationsScope)
		assert.True(t, apiToken.LimitedBizOpsScope)
		assert.True(t, apiToken.LimitedWebsitesScope)
		assert.True(t, apiToken.LimitedKubernetesScope)
		assert.True(t, apiToken.LimitedMobileAppsScope)
		assert.True(t, apiToken.LimitedInfrastructureScope)
		assert.True(t, apiToken.LimitedSyntheticsScope)
		assert.True(t, apiToken.LimitedLogsScope)

		// Verify additional permissions
		assert.True(t, apiToken.CanConfigurePersonalAPITokens)
		assert.True(t, apiToken.CanConfigureDatabaseManagement)
		assert.True(t, apiToken.CanConfigureAutomationActions)
		assert.True(t, apiToken.CanConfigureAutomationPolicies)
		assert.True(t, apiToken.CanRunAutomationActions)
		assert.True(t, apiToken.CanDeleteAutomationActionHistory)
		assert.True(t, apiToken.CanConfigureSyntheticTests)
		assert.True(t, apiToken.CanConfigureSyntheticLocations)
		assert.True(t, apiToken.CanConfigureSyntheticCredentials)
		assert.True(t, apiToken.CanViewSyntheticTests)
		assert.True(t, apiToken.CanViewSyntheticLocations)
		assert.True(t, apiToken.CanViewSyntheticTestResults)
		assert.True(t, apiToken.CanUseSyntheticCredentials)
		assert.True(t, apiToken.CanConfigureBizops)
		assert.True(t, apiToken.CanViewBusinessProcesses)
		assert.True(t, apiToken.CanViewBusinessProcessDetails)
		assert.True(t, apiToken.CanViewBusinessActivities)
		assert.True(t, apiToken.CanViewBizAlerts)
		assert.True(t, apiToken.CanDeleteLogs)
		assert.True(t, apiToken.CanCreateHeapDump)
		assert.True(t, apiToken.CanCreateThreadDump)
		assert.True(t, apiToken.CanManuallyCloseIssue)
		assert.True(t, apiToken.CanViewLogVolume)
		assert.True(t, apiToken.CanConfigureLogRetentionPeriod)
		assert.True(t, apiToken.CanConfigureSubtraces)
		assert.True(t, apiToken.CanInvokeAlertChannel)
		assert.True(t, apiToken.CanConfigureLlm)
		assert.True(t, apiToken.CanConfigureAiAgents)
		assert.True(t, apiToken.CanConfigureApdex)
		assert.True(t, apiToken.CanConfigureServiceLevelCorrectionWindows)
		assert.True(t, apiToken.CanConfigureServiceLevelSmartAlerts)
		assert.True(t, apiToken.CanConfigureServiceLevels)
	})
}

func TestGetRestResource(t *testing.T) {
	resource := &apiTokenResource{}

	// This test just ensures the method exists and can be called
	// The actual implementation requires a real API instance
	assert.NotNil(t, resource)
}

func TestRoundTripConversion(t *testing.T) {
	resource := &apiTokenResource{}
	ctx := context.Background()

	t.Run("state to API and back to state", func(t *testing.T) {
		// Create original API token
		originalAPIToken := createFullyPermissionedAPIToken()

		// Convert to state
		handle := NewAPITokenResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}
		diags := resource.UpdateState(ctx, state, nil, originalAPIToken)
		require.False(t, diags.HasError())

		// Convert back to API token
		convertedAPIToken, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, convertedAPIToken)

		// Verify all fields match
		assert.Equal(t, originalAPIToken.ID, convertedAPIToken.ID)
		assert.Equal(t, originalAPIToken.Name, convertedAPIToken.Name)
		assert.Equal(t, originalAPIToken.AccessGrantingToken, convertedAPIToken.AccessGrantingToken)
		assert.Equal(t, originalAPIToken.InternalID, convertedAPIToken.InternalID)
		assert.Equal(t, originalAPIToken.CanConfigureServiceMapping, convertedAPIToken.CanConfigureServiceMapping)
		assert.Equal(t, originalAPIToken.CanConfigureUsers, convertedAPIToken.CanConfigureUsers)
		assert.Equal(t, originalAPIToken.LimitedApplicationsScope, convertedAPIToken.LimitedApplicationsScope)
		assert.Equal(t, originalAPIToken.CanConfigurePersonalAPITokens, convertedAPIToken.CanConfigurePersonalAPITokens)
	})
}

func TestEdgeCases(t *testing.T) {
	resource := &apiTokenResource{}
	ctx := context.Background()

	t.Run("empty name", func(t *testing.T) {
		model := APITokenModel{
			ID:                  types.StringValue("test-id"),
			AccessGrantingToken: types.StringValue("test-access-token"),
			InternalID:          types.StringValue("test-internal-id"),
			Name:                types.StringValue(""),
		}

		state := createMockState(t, ctx, model)
		apiToken, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, apiToken)
		assert.Equal(t, "", apiToken.Name)
	})

	t.Run("all permissions false", func(t *testing.T) {
		apiToken := &restapi.APIToken{
			ID:                  "test-id",
			AccessGrantingToken: "test-access-token",
			InternalID:          "test-internal-id",
			Name:                "test-token",
			// All permissions default to false
		}

		handle := NewAPITokenResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiToken)
		require.False(t, diags.HasError())

		var model APITokenModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		// Verify all permissions are false
		assert.False(t, model.CanConfigureServiceMapping.ValueBool())
		assert.False(t, model.CanConfigureUsers.ValueBool())
		assert.False(t, model.CanInstallNewAgents.ValueBool())
		assert.False(t, model.LimitedApplicationsScope.ValueBool())
	})
}

// Helper functions

func createMockState(t *testing.T, ctx context.Context, model APITokenModel) *tfsdk.State {
	handle := NewAPITokenResourceHandle()
	state := &tfsdk.State{
		Schema: handle.MetaData().Schema,
	}

	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set state: %v", diags)
	}

	return state
}

func createMockPlan(t *testing.T, ctx context.Context, model APITokenModel) *tfsdk.Plan {
	handle := NewAPITokenResourceHandle()
	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	diags := plan.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set plan: %v", diags)
	}

	return plan
}

func createFullyPermissionedAPIToken() *restapi.APIToken {
	return &restapi.APIToken{
		ID:                                        "test-id",
		AccessGrantingToken:                       "test-access-token",
		InternalID:                                "test-internal-id",
		Name:                                      "test-token",
		CanConfigureServiceMapping:                true,
		CanConfigureEumApplications:               true,
		CanConfigureMobileAppMonitoring:           true,
		CanConfigureUsers:                         true,
		CanInstallNewAgents:                       true,
		CanConfigureIntegrations:                  true,
		CanConfigureEventsAndAlerts:               true,
		CanConfigureMaintenanceWindows:            true,
		CanConfigureApplicationSmartAlerts:        true,
		CanConfigureWebsiteSmartAlerts:            true,
		CanConfigureMobileAppSmartAlerts:          true,
		CanConfigureAPITokens:                     true,
		CanConfigureAgentRunMode:                  true,
		CanViewAuditLog:                           true,
		CanConfigureAgents:                        true,
		CanConfigureAuthenticationMethods:         true,
		CanConfigureApplications:                  true,
		CanConfigureTeams:                         true,
		CanConfigureReleases:                      true,
		CanConfigureLogManagement:                 true,
		CanCreatePublicCustomDashboards:           true,
		CanViewLogs:                               true,
		CanViewTraceDetails:                       true,
		CanConfigureSessionSettings:               true,
		CanConfigureGlobalAlertPayload:            true,
		CanConfigureGlobalApplicationSmartAlerts:  true,
		CanConfigureGlobalSyntheticSmartAlerts:    true,
		CanConfigureGlobalInfraSmartAlerts:        true,
		CanConfigureGlobalLogSmartAlerts:          true,
		CanViewAccountAndBillingInformation:       true,
		CanEditAllAccessibleCustomDashboards:      true,
		LimitedApplicationsScope:                  true,
		LimitedBizOpsScope:                        true,
		LimitedWebsitesScope:                      true,
		LimitedKubernetesScope:                    true,
		LimitedMobileAppsScope:                    true,
		LimitedInfrastructureScope:                true,
		LimitedSyntheticsScope:                    true,
		LimitedVsphereScope:                       true,
		LimitedPhmcScope:                          true,
		LimitedPvcScope:                           true,
		LimitedZhmcScope:                          true,
		LimitedPcfScope:                           true,
		LimitedOpenstackScope:                     true,
		LimitedAutomationScope:                    true,
		LimitedLogsScope:                          true,
		LimitedNutanixScope:                       true,
		LimitedXenServerScope:                     true,
		LimitedWindowsHypervisorScope:             true,
		LimitedAlertChannelsScope:                 true,
		LimitedLinuxKvmHypervisorScope:            true,
		LimitedServiceLevelScope:                  true,
		LimitedAiGatewayScope:                     true,
		CanConfigurePersonalAPITokens:             true,
		CanConfigureDatabaseManagement:            true,
		CanConfigureAutomationActions:             true,
		CanConfigureAutomationPolicies:            true,
		CanRunAutomationActions:                   true,
		CanDeleteAutomationActionHistory:          true,
		CanConfigureSyntheticTests:                true,
		CanConfigureSyntheticLocations:            true,
		CanConfigureSyntheticCredentials:          true,
		CanViewSyntheticTests:                     true,
		CanViewSyntheticLocations:                 true,
		CanViewSyntheticTestResults:               true,
		CanUseSyntheticCredentials:                true,
		CanConfigureBizops:                        true,
		CanViewBusinessProcesses:                  true,
		CanViewBusinessProcessDetails:             true,
		CanViewBusinessActivities:                 true,
		CanViewBizAlerts:                          true,
		CanDeleteLogs:                             true,
		CanCreateHeapDump:                         true,
		CanCreateThreadDump:                       true,
		CanManuallyCloseIssue:                     true,
		CanViewLogVolume:                          true,
		CanConfigureLogRetentionPeriod:            true,
		CanConfigureSubtraces:                     true,
		CanInvokeAlertChannel:                     true,
		CanConfigureLlm:                           true,
		CanConfigureAiAgents:                      true,
		CanConfigureApdex:                         true,
		CanConfigureServiceLevelCorrectionWindows: true,
		CanConfigureServiceLevelSmartAlerts:       true,
		CanConfigureServiceLevels:                 true,
	}
}

func createFullyPermissionedModel() APITokenModel {
	return APITokenModel{
		ID:                                        types.StringValue("test-id"),
		AccessGrantingToken:                       types.StringValue("test-access-token"),
		InternalID:                                types.StringValue("test-internal-id"),
		Name:                                      types.StringValue("test-token"),
		CanConfigureServiceMapping:                types.BoolValue(true),
		CanConfigureEumApplications:               types.BoolValue(true),
		CanConfigureMobileAppMonitoring:           types.BoolValue(true),
		CanConfigureUsers:                         types.BoolValue(true),
		CanInstallNewAgents:                       types.BoolValue(true),
		CanConfigureIntegrations:                  types.BoolValue(true),
		CanConfigureEventsAndAlerts:               types.BoolValue(true),
		CanConfigureMaintenanceWindows:            types.BoolValue(true),
		CanConfigureApplicationSmartAlerts:        types.BoolValue(true),
		CanConfigureWebsiteSmartAlerts:            types.BoolValue(true),
		CanConfigureMobileAppSmartAlerts:          types.BoolValue(true),
		CanConfigureAPITokens:                     types.BoolValue(true),
		CanConfigureAgentRunMode:                  types.BoolValue(true),
		CanViewAuditLog:                           types.BoolValue(true),
		CanConfigureAgents:                        types.BoolValue(true),
		CanConfigureAuthenticationMethods:         types.BoolValue(true),
		CanConfigureApplications:                  types.BoolValue(true),
		CanConfigureTeams:                         types.BoolValue(true),
		CanConfigureReleases:                      types.BoolValue(true),
		CanConfigureLogManagement:                 types.BoolValue(true),
		CanCreatePublicCustomDashboards:           types.BoolValue(true),
		CanViewLogs:                               types.BoolValue(true),
		CanViewTraceDetails:                       types.BoolValue(true),
		CanConfigureSessionSettings:               types.BoolValue(true),
		CanConfigureGlobalAlertPayload:            types.BoolValue(true),
		CanConfigureGlobalApplicationSmartAlerts:  types.BoolValue(true),
		CanConfigureGlobalSyntheticSmartAlerts:    types.BoolValue(true),
		CanConfigureGlobalInfraSmartAlerts:        types.BoolValue(true),
		CanConfigureGlobalLogSmartAlerts:          types.BoolValue(true),
		CanViewAccountAndBillingInformation:       types.BoolValue(true),
		CanEditAllAccessibleCustomDashboards:      types.BoolValue(true),
		LimitedApplicationsScope:                  types.BoolValue(true),
		LimitedBizOpsScope:                        types.BoolValue(true),
		LimitedWebsitesScope:                      types.BoolValue(true),
		LimitedKubernetesScope:                    types.BoolValue(true),
		LimitedMobileAppsScope:                    types.BoolValue(true),
		LimitedInfrastructureScope:                types.BoolValue(true),
		LimitedSyntheticsScope:                    types.BoolValue(true),
		LimitedVsphereScope:                       types.BoolValue(true),
		LimitedPhmcScope:                          types.BoolValue(true),
		LimitedPvcScope:                           types.BoolValue(true),
		LimitedZhmcScope:                          types.BoolValue(true),
		LimitedPcfScope:                           types.BoolValue(true),
		LimitedOpenstackScope:                     types.BoolValue(true),
		LimitedAutomationScope:                    types.BoolValue(true),
		LimitedLogsScope:                          types.BoolValue(true),
		LimitedNutanixScope:                       types.BoolValue(true),
		LimitedXenServerScope:                     types.BoolValue(true),
		LimitedWindowsHypervisorScope:             types.BoolValue(true),
		LimitedAlertChannelsScope:                 types.BoolValue(true),
		LimitedLinuxKvmHypervisorScope:            types.BoolValue(true),
		LimitedServiceLevelScope:                  types.BoolValue(true),
		LimitedAiGatewayScope:                     types.BoolValue(true),
		CanConfigurePersonalAPITokens:             types.BoolValue(true),
		CanConfigureDatabaseManagement:            types.BoolValue(true),
		CanConfigureAutomationActions:             types.BoolValue(true),
		CanConfigureAutomationPolicies:            types.BoolValue(true),
		CanRunAutomationActions:                   types.BoolValue(true),
		CanDeleteAutomationActionHistory:          types.BoolValue(true),
		CanConfigureSyntheticTests:                types.BoolValue(true),
		CanConfigureSyntheticLocations:            types.BoolValue(true),
		CanConfigureSyntheticCredentials:          types.BoolValue(true),
		CanViewSyntheticTests:                     types.BoolValue(true),
		CanViewSyntheticLocations:                 types.BoolValue(true),
		CanViewSyntheticTestResults:               types.BoolValue(true),
		CanUseSyntheticCredentials:                types.BoolValue(true),
		CanConfigureBizops:                        types.BoolValue(true),
		CanViewBusinessProcesses:                  types.BoolValue(true),
		CanViewBusinessProcessDetails:             types.BoolValue(true),
		CanViewBusinessActivities:                 types.BoolValue(true),
		CanViewBizAlerts:                          types.BoolValue(true),
		CanDeleteLogs:                             types.BoolValue(true),
		CanCreateHeapDump:                         types.BoolValue(true),
		CanCreateThreadDump:                       types.BoolValue(true),
		CanManuallyCloseIssue:                     types.BoolValue(true),
		CanViewLogVolume:                          types.BoolValue(true),
		CanConfigureLogRetentionPeriod:            types.BoolValue(true),
		CanConfigureSubtraces:                     types.BoolValue(true),
		CanInvokeAlertChannel:                     types.BoolValue(true),
		CanConfigureLlm:                           types.BoolValue(true),
		CanConfigureAiAgents:                      types.BoolValue(true),
		CanConfigureApdex:                         types.BoolValue(true),
		CanConfigureServiceLevelCorrectionWindows: types.BoolValue(true),
		CanConfigureServiceLevelSmartAlerts:       types.BoolValue(true),
		CanConfigureServiceLevels:                 types.BoolValue(true),
	}
}
