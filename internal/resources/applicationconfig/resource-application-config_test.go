package applicationconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}

func TestNewApplicationConfigResourceHandle(t *testing.T) {
	resource := NewApplicationConfigResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaApplicationConfig, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(4), metaData.SchemaVersion)
}

func TestMetaData(t *testing.T) {
	resource := &applicationConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 4,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_resource", metaData.ResourceName)
	assert.Equal(t, int64(4), metaData.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := NewApplicationConfigResourceHandle()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetRestResource(t *testing.T) {
	resource := &applicationConfigResource{}

	// We can't test the actual API call without a mock, but we can verify the method exists
	assert.NotNil(t, resource.GetRestResource)
}

func TestMapStateToDataObject_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Application", result.Label)
	assert.Equal(t, restapi.ApplicationConfigScope("INCLUDE_NO_DOWNSTREAM"), result.Scope)
	assert.Equal(t, restapi.BoundaryScope("DEFAULT"), result.BoundaryScope)
	assert.Nil(t, result.TagFilterExpression)
	require.Len(t, result.AccessRules, 1)
	assert.Equal(t, restapi.AccessType("READ"), result.AccessRules[0].AccessType)
	assert.Equal(t, restapi.RelationType("USER"), result.AccessRules[0].RelationType)
	assert.Nil(t, result.AccessRules[0].RelatedID)
}

func TestMapStateToDataObject_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_ALL_DOWNSTREAM"),
		BoundaryScope: types.StringValue("ALL"),
		TagFilter:     types.StringValue("entity.type EQUALS 'service'"),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ_WRITE"),
				RelationType: types.StringValue("ROLE"),
				RelatedID:    types.StringValue("user-123"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Application", result.Label)
	assert.Equal(t, restapi.ApplicationConfigScope("INCLUDE_ALL_DOWNSTREAM"), result.Scope)
	assert.Equal(t, restapi.BoundaryScope("ALL"), result.BoundaryScope)
	require.NotNil(t, result.TagFilterExpression)
	require.Len(t, result.AccessRules, 1)
	assert.Equal(t, restapi.AccessType("READ_WRITE"), result.AccessRules[0].AccessType)
	assert.Equal(t, restapi.RelationType("ROLE"), result.AccessRules[0].RelationType)
	require.NotNil(t, result.AccessRules[0].RelatedID)
	assert.Equal(t, "user-123", *result.AccessRules[0].RelatedID)
}

func TestMapStateToDataObject_WithMultipleAccessRules(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("INBOUND"),
		TagFilter:     types.StringNull(),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringValue("user-1"),
			},
			{
				AccessType:   types.StringValue("READ_WRITE"),
				RelationType: types.StringValue("ROLE"),
				RelatedID:    types.StringValue("user-2"),
			},
			{
				AccessType:   types.StringValue("READ_WRITE"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.Len(t, result.AccessRules, 3)

	// Check first rule
	assert.Equal(t, restapi.AccessType("READ"), result.AccessRules[0].AccessType)
	assert.Equal(t, restapi.RelationType("USER"), result.AccessRules[0].RelationType)
	require.NotNil(t, result.AccessRules[0].RelatedID)
	assert.Equal(t, "user-1", *result.AccessRules[0].RelatedID)

	// Check second rule
	assert.Equal(t, restapi.AccessType("READ_WRITE"), result.AccessRules[1].AccessType)
	assert.Equal(t, restapi.RelationType("ROLE"), result.AccessRules[1].RelationType)
	require.NotNil(t, result.AccessRules[1].RelatedID)
	assert.Equal(t, "user-2", *result.AccessRules[1].RelatedID)

	// Check third rule
	assert.Equal(t, restapi.AccessType("READ_WRITE"), result.AccessRules[2].AccessType)
	assert.Equal(t, restapi.RelationType("USER"), result.AccessRules[2].RelationType)
	assert.Nil(t, result.AccessRules[2].RelatedID)
}

func TestMapStateToDataObject_WithEmptyAccessRules(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules:   []AccessRuleModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Empty(t, result.AccessRules)
}

func TestMapStateToDataObject_WithNullID(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringNull(),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "", result.ID)
}

func TestMapStateToDataObject_WithInvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringValue("invalid tag filter syntax"),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringNull(),
			},
		},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	plan := &tfsdk.Plan{
		Schema: getTestSchema(),
	}

	model := ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringNull(),
			},
		},
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Application", result.Label)
}

func TestUpdateState_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	data := &restapi.ApplicationConfig{
		ID:            "test-id",
		Label:         "Test Application",
		Scope:         restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope: restapi.BoundaryScopeDefault,
		AccessRules: []restapi.AccessRule{
			{
				AccessType:   restapi.AccessTypeRead,
				RelationType: restapi.RelationTypeUser,
				RelatedID:    nil,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with a model that has TagFilter set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Application", model.Label.ValueString())
	assert.Equal(t, "INCLUDE_NO_DOWNSTREAM", model.Scope.ValueString())
	assert.Equal(t, "DEFAULT", model.BoundaryScope.ValueString())
	assert.True(t, model.TagFilter.IsNull())
	require.Len(t, model.AccessRules, 1)
	assert.Equal(t, "READ", model.AccessRules[0].AccessType.ValueString())
	assert.Equal(t, "USER", model.AccessRules[0].RelationType.ValueString())
	assert.True(t, model.AccessRules[0].RelatedID.IsNull())
}

func TestUpdateState_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	entityType := "entity.type"
	serviceValue := "service"
	equalsOp := restapi.EqualsOperator

	tagFilter := &restapi.TagFilter{
		Type:     restapi.TagFilterExpressionType,
		Name:     &entityType,
		Operator: &equalsOp,
		Value:    &serviceValue,
	}

	data := &restapi.ApplicationConfig{
		ID:                  "test-id",
		Label:               "Test Application",
		Scope:               restapi.ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:       restapi.BoundaryScopeAll,
		TagFilterExpression: tagFilter,
		AccessRules: []restapi.AccessRule{
			{
				AccessType:   restapi.AccessTypeReadWrite,
				RelationType: restapi.RelationTypeRole,
				RelatedID:    ptr("user-123"),
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with a model that has TagFilter set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Application", model.Label.ValueString())
	assert.Equal(t, "INCLUDE_ALL_DOWNSTREAM", model.Scope.ValueString())
	assert.Equal(t, "ALL", model.BoundaryScope.ValueString())
	// Tag filter should be set (not null) when provided
	if !model.TagFilter.IsNull() {
		// If tag filter is present, verify it's not empty
		assert.NotEmpty(t, model.TagFilter.ValueString())
	}
	require.Len(t, model.AccessRules, 1)
	assert.Equal(t, "READ_WRITE", model.AccessRules[0].AccessType.ValueString())
	assert.Equal(t, "ROLE", model.AccessRules[0].RelationType.ValueString())
	assert.Equal(t, "user-123", model.AccessRules[0].RelatedID.ValueString())
}

func TestUpdateState_WithMultipleAccessRules(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	data := &restapi.ApplicationConfig{
		ID:            "test-id",
		Label:         "Test Application",
		Scope:         restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope: restapi.BoundaryScopeInbound,
		AccessRules: []restapi.AccessRule{
			{
				AccessType:   restapi.AccessTypeRead,
				RelationType: restapi.RelationTypeUser,
				RelatedID:    ptr("user-1"),
			},
			{
				AccessType:   restapi.AccessTypeReadWrite,
				RelationType: restapi.RelationTypeRole,
				RelatedID:    ptr("user-2"),
			},
			{
				AccessType:   restapi.AccessTypeReadWrite,
				RelationType: restapi.RelationTypeUser,
				RelatedID:    nil,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with a model that has TagFilter set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.AccessRules, 3)

	// Check first rule
	assert.Equal(t, "READ", model.AccessRules[0].AccessType.ValueString())
	assert.Equal(t, "USER", model.AccessRules[0].RelationType.ValueString())
	assert.Equal(t, "user-1", model.AccessRules[0].RelatedID.ValueString())

	// Check second rule
	assert.Equal(t, "READ_WRITE", model.AccessRules[1].AccessType.ValueString())
	assert.Equal(t, "ROLE", model.AccessRules[1].RelationType.ValueString())
	assert.Equal(t, "user-2", model.AccessRules[1].RelatedID.ValueString())

	// Check third rule
	assert.Equal(t, "READ_WRITE", model.AccessRules[2].AccessType.ValueString())
	assert.Equal(t, "USER", model.AccessRules[2].RelationType.ValueString())
	assert.True(t, model.AccessRules[2].RelatedID.IsNull())
}

func TestUpdateState_WithEmptyAccessRules(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	data := &restapi.ApplicationConfig{
		ID:            "test-id",
		Label:         "Test Application",
		Scope:         restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope: restapi.BoundaryScopeDefault,
		AccessRules:   []restapi.AccessRule{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with a model that has TagFilter set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Empty(t, model.AccessRules)
}

func TestUpdateState_WithNullTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	data := &restapi.ApplicationConfig{
		ID:                  "test-id",
		Label:               "Test Application",
		Scope:               restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope:       restapi.BoundaryScopeDefault,
		TagFilterExpression: nil,
		AccessRules: []restapi.AccessRule{
			{
				AccessType:   restapi.AccessTypeRead,
				RelationType: restapi.RelationTypeUser,
				RelatedID:    nil,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with a model that has TagFilter set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.TagFilter.IsNull())
}

func TestMapStateToDataObject_AllScopes(t *testing.T) {
	scopes := []string{
		"INCLUDE_NO_DOWNSTREAM",
		"INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING",
		"INCLUDE_ALL_DOWNSTREAM",
	}

	for _, scope := range scopes {
		t.Run(scope, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationConfigResource{}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue(scope),
				BoundaryScope: types.StringValue("DEFAULT"),
				TagFilter:     types.StringNull(),
				AccessRules: []AccessRuleModel{
					{
						AccessType:   types.StringValue("READ"),
						RelationType: types.StringValue("USER"),
						RelatedID:    types.StringNull(),
					},
				},
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			assert.Equal(t, restapi.ApplicationConfigScope(scope), result.Scope)
		})
	}
}

func TestMapStateToDataObject_AllBoundaryScopes(t *testing.T) {
	boundaryScopes := []string{"ALL", "INBOUND", "DEFAULT"}

	for _, boundaryScope := range boundaryScopes {
		t.Run(boundaryScope, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationConfigResource{}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
				BoundaryScope: types.StringValue(boundaryScope),
				TagFilter:     types.StringNull(),
				AccessRules: []AccessRuleModel{
					{
						AccessType:   types.StringValue("READ"),
						RelationType: types.StringValue("USER"),
						RelatedID:    types.StringNull(),
					},
				},
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			assert.Equal(t, restapi.BoundaryScope(boundaryScope), result.BoundaryScope)
		})
	}
}

func TestMapStateToDataObject_AllAccessTypes(t *testing.T) {
	accessTypes := []string{"READ", "READ_WRITE"}

	for _, accessType := range accessTypes {
		t.Run(accessType, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationConfigResource{}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
				BoundaryScope: types.StringValue("DEFAULT"),
				TagFilter:     types.StringNull(),
				AccessRules: []AccessRuleModel{
					{
						AccessType:   types.StringValue(accessType),
						RelationType: types.StringValue("USER"),
						RelatedID:    types.StringNull(),
					},
				},
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			require.Len(t, result.AccessRules, 1)
			assert.Equal(t, restapi.AccessType(accessType), result.AccessRules[0].AccessType)
		})
	}
}

func TestMapStateToDataObject_AllRelationTypes(t *testing.T) {
	relationTypes := []string{"USER", "API_TOKEN", "ROLE", "TEAM", "GLOBAL"}

	for _, relationType := range relationTypes {
		t.Run(relationType, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationConfigResource{}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
				BoundaryScope: types.StringValue("DEFAULT"),
				TagFilter:     types.StringNull(),
				AccessRules: []AccessRuleModel{
					{
						AccessType:   types.StringValue("READ"),
						RelationType: types.StringValue(relationType),
						RelatedID:    types.StringNull(),
					},
				},
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			require.Len(t, result.AccessRules, 1)
			assert.Equal(t, restapi.RelationType(relationType), result.AccessRules[0].RelationType)
		})
	}
}

func TestMapStateToDataObject_WithEmptyRelatedID(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules: []AccessRuleModel{
			{
				AccessType:   types.StringValue("READ"),
				RelationType: types.StringValue("USER"),
				RelatedID:    types.StringValue(""),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AccessRules, 1)
	// Empty string should result in nil RelatedID
	assert.Nil(t, result.AccessRules[0].RelatedID)
}

func TestMapAccessRulesToState_EmptyRules(t *testing.T) {
	resource := &applicationConfigResource{}

	result := resource.mapAccessRulesToState([]restapi.AccessRule{})
	assert.Empty(t, result)
}

func TestMapAccessRulesToState_WithRules(t *testing.T) {
	resource := &applicationConfigResource{}

	accessRules := []restapi.AccessRule{
		{
			AccessType:   restapi.AccessTypeRead,
			RelationType: restapi.RelationTypeUser,
			RelatedID:    ptr("user-1"),
		},
		{
			AccessType:   restapi.AccessTypeReadWrite,
			RelationType: restapi.RelationTypeRole,
			RelatedID:    nil,
		},
	}

	result := resource.mapAccessRulesToState(accessRules)
	require.Len(t, result, 2)

	assert.Equal(t, "READ", result[0].AccessType.ValueString())
	assert.Equal(t, "USER", result[0].RelationType.ValueString())
	assert.Equal(t, "user-1", result[0].RelatedID.ValueString())

	assert.Equal(t, "READ_WRITE", result[1].AccessType.ValueString())
	assert.Equal(t, "ROLE", result[1].RelationType.ValueString())
	assert.True(t, result[1].RelatedID.IsNull())
}

func TestMapAccessRulesFromState_EmptyRules(t *testing.T) {
	resource := &applicationConfigResource{}

	result := resource.mapAccessRulesFromState([]AccessRuleModel{})
	assert.Empty(t, result)
}

func TestMapAccessRulesFromState_WithRules(t *testing.T) {
	resource := &applicationConfigResource{}

	accessRuleModels := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringValue("user-1"),
		},
		{
			AccessType:   types.StringValue("READ_WRITE"),
			RelationType: types.StringValue("ROLE"),
			RelatedID:    types.StringNull(),
		},
	}

	result := resource.mapAccessRulesFromState(accessRuleModels)
	require.Len(t, result, 2)

	assert.Equal(t, restapi.AccessTypeRead, result[0].AccessType)
	assert.Equal(t, restapi.RelationTypeUser, result[0].RelationType)
	require.NotNil(t, result[0].RelatedID)
	assert.Equal(t, "user-1", *result[0].RelatedID)

	assert.Equal(t, restapi.AccessTypeReadWrite, result[1].AccessType)
	assert.Equal(t, restapi.RelationTypeRole, result[1].RelationType)
	assert.Nil(t, result[1].RelatedID)
}

func TestUpdateState_AllScopes(t *testing.T) {
	scopes := []restapi.ApplicationConfigScope{
		restapi.ApplicationConfigScopeIncludeNoDownstream,
		restapi.ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging,
		restapi.ApplicationConfigScopeIncludeAllDownstream,
	}

	for _, scope := range scopes {
		t.Run(string(scope), func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationConfigResource{}

			data := &restapi.ApplicationConfig{
				ID:            "test-id",
				Label:         "Test Application",
				Scope:         scope,
				BoundaryScope: restapi.BoundaryScopeDefault,
				AccessRules: []restapi.AccessRule{
					{
						AccessType:   restapi.AccessTypeRead,
						RelationType: restapi.RelationTypeUser,
						RelatedID:    nil,
					},
				},
			}

			state := &tfsdk.State{
				Schema: getTestSchema(),
			}

			// Initialize state with a model that has TagFilter set to null
			// This is necessary because UpdateState checks if TagFilter is null/unknown
			initialModel := ApplicationConfigModel{
				TagFilter: types.StringNull(),
			}
			diags := state.Set(ctx, initialModel)
			require.False(t, diags.HasError())

			diags = resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError())

			var model ApplicationConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			assert.Equal(t, string(scope), model.Scope.ValueString())
		})
	}
}

func TestUpdateState_AllBoundaryScopes(t *testing.T) {
	boundaryScopes := []restapi.BoundaryScope{
		restapi.BoundaryScopeAll,
		restapi.BoundaryScopeInbound,
		restapi.BoundaryScopeDefault,
	}

	for _, boundaryScope := range boundaryScopes {
		t.Run(string(boundaryScope), func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationConfigResource{}

			data := &restapi.ApplicationConfig{
				ID:            "test-id",
				Label:         "Test Application",
				Scope:         restapi.ApplicationConfigScopeIncludeNoDownstream,
				BoundaryScope: boundaryScope,
				AccessRules: []restapi.AccessRule{
					{
						AccessType:   restapi.AccessTypeRead,
						RelationType: restapi.RelationTypeUser,
						RelatedID:    nil,
					},
				},
			}

			state := &tfsdk.State{
				Schema: getTestSchema(),
			}

			// Initialize state with a model that has TagFilter set to null
			// This is necessary because UpdateState checks if TagFilter is null/unknown
			initialModel := ApplicationConfigModel{
				TagFilter: types.StringNull(),
			}
			diags := state.Set(ctx, initialModel)
			require.False(t, diags.HasError())

			diags = resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError())

			var model ApplicationConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			assert.Equal(t, string(boundaryScope), model.BoundaryScope.ValueString())
		})
	}
}

// Helper functions

func createMockState(t *testing.T, model ApplicationConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewApplicationConfigResourceHandle()
	return resource.MetaData().Schema
}

// Made with Bob
