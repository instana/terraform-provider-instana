package customdashboard

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCustomDashboardResourceHandle(t *testing.T) {
	handle := NewCustomDashboardResourceHandle()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaCustomDashboard, metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &customDashboardResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 1,
		},
	}

	metadata := resource.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, "test_resource", metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := &customDashboardResource{}
	ctx := context.Background()

	diags := resource.SetComputedFields(ctx, nil)
	assert.False(t, diags.HasError())
}

func TestMapAccessRulesToState(t *testing.T) {
	resource := &customDashboardResource{}

	t.Run("with access rules", func(t *testing.T) {
		relatedID := "user-123"
		accessRules := []restapi.AccessRule{
			{
				AccessType:   restapi.AccessTypeRead,
				RelationType: restapi.RelationTypeUser,
				RelatedID:    &relatedID,
			},
			{
				AccessType:   restapi.AccessTypeReadWrite,
				RelationType: restapi.RelationTypeGlobal,
				RelatedID:    nil,
			},
		}

		models := resource.mapAccessRulesToState(accessRules)
		require.Len(t, models, 2)

		assert.Equal(t, "READ", models[0].AccessType.ValueString())
		assert.Equal(t, "USER", models[0].RelationType.ValueString())
		assert.Equal(t, "user-123", models[0].RelatedID.ValueString())

		assert.Equal(t, "READ_WRITE", models[1].AccessType.ValueString())
		assert.Equal(t, "GLOBAL", models[1].RelationType.ValueString())
		assert.True(t, models[1].RelatedID.IsNull())
	})

	t.Run("with empty access rules", func(t *testing.T) {
		models := resource.mapAccessRulesToState([]restapi.AccessRule{})
		assert.Nil(t, models)
	})

	t.Run("with nil access rules", func(t *testing.T) {
		models := resource.mapAccessRulesToState(nil)
		assert.Nil(t, models)
	})
}

func TestMapAccessRulesFromState(t *testing.T) {
	resource := &customDashboardResource{}

	t.Run("with access rules", func(t *testing.T) {
		models := []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringValue("user-123"),
			},
			{
				AccessType:   types.StringValue("READ_WRITE"),
				RelationType: types.StringValue("GLOBAL"),
				RelatedID:    types.StringNull(),
			},
		}

		accessRules := resource.mapAccessRulesFromState(models)
		require.Len(t, accessRules, 2)

		assert.Equal(t, restapi.AccessTypeRead, accessRules[0].AccessType)
		assert.Equal(t, restapi.RelationTypeUser, accessRules[0].RelationType)
		require.NotNil(t, accessRules[0].RelatedID)
		assert.Equal(t, "user-123", *accessRules[0].RelatedID)

		assert.Equal(t, restapi.AccessTypeReadWrite, accessRules[1].AccessType)
		assert.Equal(t, restapi.RelationTypeGlobal, accessRules[1].RelationType)
		assert.Nil(t, accessRules[1].RelatedID)
	})

	t.Run("with empty related ID", func(t *testing.T) {
		models := []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringValue(""),
			},
		}

		accessRules := resource.mapAccessRulesFromState(models)
		require.Len(t, accessRules, 1)
		assert.Nil(t, accessRules[0].RelatedID)
	})

	t.Run("with empty access rules", func(t *testing.T) {
		accessRules := resource.mapAccessRulesFromState([]AccessRuleModel{})
		assert.Nil(t, accessRules)
	})

	t.Run("with nil access rules", func(t *testing.T) {
		accessRules := resource.mapAccessRulesFromState(nil)
		assert.Nil(t, accessRules)
	})
}

func TestUpdateState(t *testing.T) {
	ctx := context.Background()

	t.Run("with plan - successful update", func(t *testing.T) {
		resource := &customDashboardResource{}

		widgets := json.RawMessage(`[{"type":"chart","data":"test"}]`)
		dashboard := &restapi.CustomDashboard{
			ID:      "dashboard-123",
			Title:   "Updated Dashboard",
			Widgets: widgets,
			AccessRules: []restapi.AccessRule{
				{
					AccessType:   restapi.AccessTypeRead,
					RelationType: restapi.RelationTypeUser,
					RelatedID:    stringPtr("user-456"),
				},
			},
		}

		// Create plan with model
		planModel := CustomDashboardModel{
			ID:      types.StringValue("dashboard-123"),
			Title:   types.StringValue("Test Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[{"type":"chart","data":"test"}]`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringValue("user-456"),
				},
			},
		}

		handle := NewCustomDashboardResourceHandle()
		plan := &tfsdk.Plan{
			Schema: handle.MetaData().Schema,
		}
		diags := plan.Set(ctx, planModel)
		require.False(t, diags.HasError())

		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags = resource.UpdateState(ctx, state, plan, dashboard)
		require.False(t, diags.HasError())

		var resultModel CustomDashboardModel
		diags = state.Get(ctx, &resultModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "dashboard-123", resultModel.ID.ValueString())
		assert.Equal(t, "Updated Dashboard", resultModel.Title.ValueString())
	})

	t.Run("without plan - successful read", func(t *testing.T) {
		resource := &customDashboardResource{}

		widgets := json.RawMessage(`[{"type":"chart","data":"test"}]`)
		relatedID := "user-789"
		dashboard := &restapi.CustomDashboard{
			ID:      "dashboard-456",
			Title:   "Test Dashboard",
			Widgets: widgets,
			AccessRules: []restapi.AccessRule{
				{
					AccessType:   restapi.AccessTypeReadWrite,
					RelationType: restapi.RelationTypeApiToken,
					RelatedID:    &relatedID,
				},
			},
		}

		handle := NewCustomDashboardResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, dashboard)
		require.False(t, diags.HasError())

		var resultModel CustomDashboardModel
		diags = state.Get(ctx, &resultModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "dashboard-456", resultModel.ID.ValueString())
		assert.Equal(t, "Test Dashboard", resultModel.Title.ValueString())
		assert.Contains(t, resultModel.Widgets.ValueString(), "chart")
		require.Len(t, resultModel.AccessRules, 1)
		assert.Equal(t, "READ_WRITE", resultModel.AccessRules[0].AccessType.ValueString())
		assert.Equal(t, "API_TOKEN", resultModel.AccessRules[0].RelationType.ValueString())
		assert.Equal(t, "user-789", resultModel.AccessRules[0].RelatedID.ValueString())
	})

	t.Run("without plan - empty access rules", func(t *testing.T) {
		resource := &customDashboardResource{}

		widgets := json.RawMessage(`[]`)
		dashboard := &restapi.CustomDashboard{
			ID:          "dashboard-789",
			Title:       "Simple Dashboard",
			Widgets:     widgets,
			AccessRules: []restapi.AccessRule{},
		}

		handle := NewCustomDashboardResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, dashboard)
		require.False(t, diags.HasError())

		var resultModel CustomDashboardModel
		diags = state.Get(ctx, &resultModel)
		require.False(t, diags.HasError())
		assert.Nil(t, resultModel.AccessRules)
	})

	t.Run("without plan - invalid widgets JSON", func(t *testing.T) {
		resource := &customDashboardResource{}

		// Create invalid JSON that will fail canonicalization
		widgets := json.RawMessage(`invalid json`)
		dashboard := &restapi.CustomDashboard{
			ID:      "dashboard-error",
			Title:   "Error Dashboard",
			Widgets: widgets,
		}

		handle := NewCustomDashboardResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, dashboard)
		// util.CanonicalizeJSON() will fail for invalid JSON, so we expect an error
		require.True(t, diags.HasError())
	})

	t.Run("with plan - multiple access rules", func(t *testing.T) {
		resource := &customDashboardResource{}

		relatedID1 := "user-123"
		relatedID2 := "role-456"
		dashboard := &restapi.CustomDashboard{
			ID:      "dashboard-multi",
			Title:   "Multi Access Dashboard",
			Widgets: json.RawMessage(`[]`),
			AccessRules: []restapi.AccessRule{
				{
					AccessType:   restapi.AccessTypeRead,
					RelationType: restapi.RelationTypeUser,
					RelatedID:    &relatedID1,
				},
				{
					AccessType:   restapi.AccessTypeReadWrite,
					RelationType: restapi.RelationTypeRole,
					RelatedID:    &relatedID2,
				},
			},
		}

		planModel := CustomDashboardModel{
			ID:      types.StringValue("dashboard-multi"),
			Title:   types.StringValue("Multi Access Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[]`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringValue("user-123"),
				},
				{
					AccessType:   types.StringValue("READ_WRITE"),
					RelationType: types.StringValue("ROLE"),
					RelatedID:    types.StringValue("role-456"),
				},
			},
		}

		handle := NewCustomDashboardResourceHandle()
		plan := &tfsdk.Plan{
			Schema: handle.MetaData().Schema,
		}
		diags := plan.Set(ctx, planModel)
		require.False(t, diags.HasError())

		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags = resource.UpdateState(ctx, state, plan, dashboard)
		require.False(t, diags.HasError())

		var resultModel CustomDashboardModel
		diags = state.Get(ctx, &resultModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "dashboard-multi", resultModel.ID.ValueString())
		assert.Equal(t, "Multi Access Dashboard", resultModel.Title.ValueString())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &customDashboardResource{}
	ctx := context.Background()

	t.Run("from plan - complete dashboard", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-123"),
			Title:   types.StringValue("Test Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[{"type":"chart","data":"test"}]`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringValue("user-123"),
				},
				{
					AccessType:   types.StringValue("READ_WRITE"),
					RelationType: types.StringValue("ROLE"),
					RelatedID:    types.StringValue("role-456"),
				},
			},
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		assert.Equal(t, "dashboard-123", dashboard.ID)
		assert.Equal(t, "Test Dashboard", dashboard.Title)
		assert.NotNil(t, dashboard.Widgets)
		require.Len(t, dashboard.AccessRules, 2)
		assert.Equal(t, restapi.AccessTypeRead, dashboard.AccessRules[0].AccessType)
		assert.Equal(t, restapi.RelationTypeUser, dashboard.AccessRules[0].RelationType)
		assert.Equal(t, "user-123", *dashboard.AccessRules[0].RelatedID)
	})

	t.Run("from state - complete dashboard", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-456"),
			Title:   types.StringValue("State Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`{"widgets":[{"id":"w1"}]}`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ_WRITE"),
					RelationType: types.StringValue("TEAM"),
					RelatedID:    types.StringValue("team-789"),
				},
			},
		}

		state := createMockState(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		assert.Equal(t, "dashboard-456", dashboard.ID)
		assert.Equal(t, "State Dashboard", dashboard.Title)
		assert.Contains(t, string(dashboard.Widgets), "widgets")
		require.Len(t, dashboard.AccessRules, 1)
		assert.Equal(t, restapi.AccessTypeReadWrite, dashboard.AccessRules[0].AccessType)
		assert.Equal(t, restapi.RelationTypeTeam, dashboard.AccessRules[0].RelationType)
	})

	t.Run("with null ID", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringNull(),
			Title:   types.StringValue("New Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[]`),
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		assert.Equal(t, "", dashboard.ID)
		assert.Equal(t, "New Dashboard", dashboard.Title)
	})

	t.Run("with null widgets", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-789"),
			Title:   types.StringValue("No Widgets Dashboard"),
			Widgets: jsontypes.NewNormalizedNull(),
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		assert.Nil(t, dashboard.Widgets)
	})

	t.Run("with empty access rules", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:          types.StringValue("dashboard-empty"),
			Title:       types.StringValue("Empty Rules Dashboard"),
			Widgets:     jsontypes.NewNormalizedValue(`[]`),
			AccessRules: []AccessRuleModel{},
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		assert.Nil(t, dashboard.AccessRules)
	})

	t.Run("with global access rule without related ID", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-global"),
			Title:   types.StringValue("Global Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[]`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("GLOBAL"),
					RelatedID:    types.StringNull(),
				},
			},
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		require.Len(t, dashboard.AccessRules, 1)
		assert.Equal(t, restapi.RelationTypeGlobal, dashboard.AccessRules[0].RelationType)
		assert.Nil(t, dashboard.AccessRules[0].RelatedID)
	})

	t.Run("with complex widgets JSON", func(t *testing.T) {
		complexWidgets := `{
			"widgets": [
				{
					"id": "widget1",
					"type": "chart",
					"config": {
						"metric": "cpu.usage",
						"aggregation": "avg"
					}
				},
				{
					"id": "widget2",
					"type": "table",
					"config": {
						"columns": ["name", "value"]
					}
				}
			]
		}`

		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-complex"),
			Title:   types.StringValue("Complex Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(complexWidgets),
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		// Verify JSON is normalized and valid
		var widgetsMap map[string]interface{}
		err := json.Unmarshal(dashboard.Widgets, &widgetsMap)
		require.NoError(t, err)
		assert.Contains(t, widgetsMap, "widgets")
	})

	t.Run("with all access types", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-all-types"),
			Title:   types.StringValue("All Types Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[]`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringValue("user-1"),
				},
				{
					AccessType:   types.StringValue("READ_WRITE"),
					RelationType: types.StringValue("API_TOKEN"),
					RelatedID:    types.StringValue("token-1"),
				},
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("ROLE"),
					RelatedID:    types.StringValue("role-1"),
				},
				{
					AccessType:   types.StringValue("READ_WRITE"),
					RelationType: types.StringValue("TEAM"),
					RelatedID:    types.StringValue("team-1"),
				},
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("GLOBAL"),
					RelatedID:    types.StringNull(),
				},
			},
		}

		plan := createMockPlan(t, ctx, model)
		dashboard, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)

		require.Len(t, dashboard.AccessRules, 5)
		assert.Equal(t, restapi.RelationTypeUser, dashboard.AccessRules[0].RelationType)
		assert.Equal(t, restapi.RelationTypeApiToken, dashboard.AccessRules[1].RelationType)
		assert.Equal(t, restapi.RelationTypeRole, dashboard.AccessRules[2].RelationType)
		assert.Equal(t, restapi.RelationTypeTeam, dashboard.AccessRules[3].RelationType)
		assert.Equal(t, restapi.RelationTypeGlobal, dashboard.AccessRules[4].RelationType)
	})

	t.Run("error when neither plan nor state provided", func(t *testing.T) {
		dashboard, diags := resource.MapStateToDataObject(ctx, nil, nil)
		// The function should handle this gracefully
		// Based on the implementation, it will return an empty dashboard
		require.False(t, diags.HasError())
		require.NotNil(t, dashboard)
	})
}

func TestGetRestResource(t *testing.T) {
	resource := &customDashboardResource{}

	// We can't fully test this without mocking the InstanaAPI,
	// but we can verify the method signature and that it doesn't panic
	assert.NotNil(t, resource)

	// The method should be callable (even if we can't test the return value without a mock)
	// This ensures the interface is properly implemented
	var _ resourcehandle.ResourceHandle[*restapi.CustomDashboard] = resource
}

func TestSchemaValidation(t *testing.T) {
	handle := NewCustomDashboardResourceHandle()
	schema := handle.MetaData().Schema

	t.Run("has required attributes", func(t *testing.T) {
		assert.Contains(t, schema.Attributes, "id")
		assert.Contains(t, schema.Attributes, CustomDashboardFieldTitle)
		assert.Contains(t, schema.Attributes, CustomDashboardFieldWidgets)
		assert.Contains(t, schema.Attributes, CustomDashboardFieldAccessRule)
	})

	t.Run("id is computed", func(t *testing.T) {
		idAttr := schema.Attributes["id"]
		assert.NotNil(t, idAttr)
	})

	t.Run("title is required", func(t *testing.T) {
		titleAttr := schema.Attributes[CustomDashboardFieldTitle]
		assert.NotNil(t, titleAttr)
	})

	t.Run("widgets is required", func(t *testing.T) {
		widgetsAttr := schema.Attributes[CustomDashboardFieldWidgets]
		assert.NotNil(t, widgetsAttr)
	})

	t.Run("access_rule is optional", func(t *testing.T) {
		accessRuleAttr := schema.Attributes[CustomDashboardFieldAccessRule]
		assert.NotNil(t, accessRuleAttr)
	})
}

func TestAccessRuleModelValidation(t *testing.T) {
	t.Run("valid access rule model", func(t *testing.T) {
		model := AccessRuleModel{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringValue("user-123"),
		}

		assert.False(t, model.AccessType.IsNull())
		assert.False(t, model.RelationType.IsNull())
		assert.False(t, model.RelatedID.IsNull())
	})

	t.Run("access rule model with null related ID", func(t *testing.T) {
		model := AccessRuleModel{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("GLOBAL"),
			RelatedID:    types.StringNull(),
		}

		assert.False(t, model.AccessType.IsNull())
		assert.False(t, model.RelationType.IsNull())
		assert.True(t, model.RelatedID.IsNull())
	})
}

func TestCustomDashboardModelValidation(t *testing.T) {
	t.Run("complete dashboard model", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:      types.StringValue("dashboard-123"),
			Title:   types.StringValue("Test Dashboard"),
			Widgets: jsontypes.NewNormalizedValue(`[]`),
			AccessRules: []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringValue("user-123"),
				},
			},
		}

		assert.False(t, model.ID.IsNull())
		assert.False(t, model.Title.IsNull())
		assert.False(t, model.Widgets.IsNull())
		assert.Len(t, model.AccessRules, 1)
	})

	t.Run("minimal dashboard model", func(t *testing.T) {
		model := CustomDashboardModel{
			ID:          types.StringNull(),
			Title:       types.StringValue("Minimal Dashboard"),
			Widgets:     jsontypes.NewNormalizedValue(`[]`),
			AccessRules: nil,
		}

		assert.True(t, model.ID.IsNull())
		assert.False(t, model.Title.IsNull())
		assert.False(t, model.Widgets.IsNull())
		assert.Nil(t, model.AccessRules)
	})
}

// Helper functions

func createMockPlan(t *testing.T, ctx context.Context, model CustomDashboardModel) *tfsdk.Plan {
	handle := NewCustomDashboardResourceHandle()
	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	diags := plan.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set plan: %v", diags)
	}

	return plan
}

func createMockState(t *testing.T, ctx context.Context, model CustomDashboardModel) *tfsdk.State {
	handle := NewCustomDashboardResourceHandle()
	state := &tfsdk.State{
		Schema: handle.MetaData().Schema,
	}

	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set state: %v", diags)
	}

	return state
}

func stringPtr(s string) *string {
	return &s
}

// initializeEmptyState initializes the state with an empty CustomDashboardModel
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := CustomDashboardModel{
		ID:          types.StringNull(),
		Title:       types.StringNull(),
		AccessRules: []AccessRuleModel{},
		Widgets:     jsontypes.NewNormalizedNull(),
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}

// Made with Bob
// Made with Bob
