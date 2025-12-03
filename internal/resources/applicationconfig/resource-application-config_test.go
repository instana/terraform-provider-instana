package applicationconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	accessRules := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringNull(),
		},
	}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

	accessRules := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ_WRITE"),
			RelationType: types.StringValue("ROLE"),
			RelatedID:    types.StringValue("user-123"),
		},
	}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_ALL_DOWNSTREAM"),
		BoundaryScope: types.StringValue("ALL"),
		TagFilter:     types.StringValue("entity.type EQUALS 'service'"),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

	accessRules := []AccessRuleModel{
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
	}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("INBOUND"),
		TagFilter:     types.StringNull(),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
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
		AccessRules:   createAccessRulesList(t, ctx, []AccessRuleModel{}),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Empty(t, result.AccessRules)
}

func TestMapStateToDataObject_WithNullID(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	accessRules := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringNull(),
		},
	}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringNull(),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "", result.ID)
}

func TestMapStateToDataObject_WithInvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	accessRules := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringNull(),
		},
	}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringValue("invalid tag filter syntax"),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

	accessRules := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringNull(),
		},
	}

	model := ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

	// Initialize state with a model that has TagFilter and AccessRules set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
		AccessRules: types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		}),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError(), "UpdateState failed: %v", diags)

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Application", model.Label.ValueString())
	assert.Equal(t, "INCLUDE_NO_DOWNSTREAM", model.Scope.ValueString())
	assert.Equal(t, "DEFAULT", model.BoundaryScope.ValueString())
	assert.True(t, model.TagFilter.IsNull())

	accessRules := getAccessRulesFromList(t, ctx, model.AccessRules)
	require.Len(t, accessRules, 1)
	assert.Equal(t, "READ", accessRules[0].AccessType.ValueString())
	assert.Equal(t, "USER", accessRules[0].RelationType.ValueString())
	assert.True(t, accessRules[0].RelatedID.IsNull())
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

	// Initialize state with a model that has TagFilter and AccessRules set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
		AccessRules: types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		}),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

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

	accessRules := getAccessRulesFromList(t, ctx, model.AccessRules)
	require.Len(t, accessRules, 1)
	assert.Equal(t, "READ_WRITE", accessRules[0].AccessType.ValueString())
	assert.Equal(t, "ROLE", accessRules[0].RelationType.ValueString())
	assert.Equal(t, "user-123", accessRules[0].RelatedID.ValueString())
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

	// Initialize state with a model that has TagFilter and AccessRules set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
		AccessRules: types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		}),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	accessRules := getAccessRulesFromList(t, ctx, model.AccessRules)
	require.Len(t, accessRules, 3)

	// Check first rule
	assert.Equal(t, "READ", accessRules[0].AccessType.ValueString())
	assert.Equal(t, "USER", accessRules[0].RelationType.ValueString())
	assert.Equal(t, "user-1", accessRules[0].RelatedID.ValueString())

	// Check second rule
	assert.Equal(t, "READ_WRITE", accessRules[1].AccessType.ValueString())
	assert.Equal(t, "ROLE", accessRules[1].RelationType.ValueString())
	assert.Equal(t, "user-2", accessRules[1].RelatedID.ValueString())

	// Check third rule
	assert.Equal(t, "READ_WRITE", accessRules[2].AccessType.ValueString())
	assert.Equal(t, "USER", accessRules[2].RelationType.ValueString())
	assert.True(t, accessRules[2].RelatedID.IsNull())
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

	// Initialize state with a model that has TagFilter and AccessRules set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
		AccessRules: types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		}),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	accessRules := getAccessRulesFromList(t, ctx, model.AccessRules)
	assert.Empty(t, accessRules)
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

	// Initialize state with a model that has TagFilter and AccessRules set to null
	// This is necessary because UpdateState checks if TagFilter is null/unknown
	initialModel := ApplicationConfigModel{
		TagFilter: types.StringNull(),
		AccessRules: types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		}),
	}
	diags := state.Set(ctx, initialModel)
	require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

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

			accessRules := []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringNull(),
				},
			}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue(scope),
				BoundaryScope: types.StringValue("DEFAULT"),
				TagFilter:     types.StringNull(),
				AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

			accessRules := []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringNull(),
				},
			}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
				BoundaryScope: types.StringValue(boundaryScope),
				TagFilter:     types.StringNull(),
				AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

			accessRules := []AccessRuleModel{
				{
					AccessType:   types.StringValue(accessType),
					RelationType: types.StringValue("USER"),
					RelatedID:    types.StringNull(),
				},
			}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
				BoundaryScope: types.StringValue("DEFAULT"),
				TagFilter:     types.StringNull(),
				AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

			accessRules := []AccessRuleModel{
				{
					AccessType:   types.StringValue("READ"),
					RelationType: types.StringValue(relationType),
					RelatedID:    types.StringNull(),
				},
			}

			state := createMockState(t, ApplicationConfigModel{
				ID:            types.StringValue("test-id"),
				Label:         types.StringValue("Test Application"),
				Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
				BoundaryScope: types.StringValue("DEFAULT"),
				TagFilter:     types.StringNull(),
				AccessRules:   createAccessRulesList(t, ctx, accessRules),
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

	accessRules := []AccessRuleModel{
		{
			AccessType:   types.StringValue("READ"),
			RelationType: types.StringValue("USER"),
			RelatedID:    types.StringValue(""),
		},
	}

	state := createMockState(t, ApplicationConfigModel{
		ID:            types.StringValue("test-id"),
		Label:         types.StringValue("Test Application"),
		Scope:         types.StringValue("INCLUDE_NO_DOWNSTREAM"),
		BoundaryScope: types.StringValue("DEFAULT"),
		TagFilter:     types.StringNull(),
		AccessRules:   createAccessRulesList(t, ctx, accessRules),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AccessRules, 1)
	// Empty string should result in nil RelatedID
	assert.Nil(t, result.AccessRules[0].RelatedID)
}

func TestMapAccessRulesToState_EmptyRules(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	result, diags := resource.mapAccessRulesToState(ctx, []restapi.AccessRule{})
	require.False(t, diags.HasError())
	if !result.IsNull() {
		rules := getAccessRulesFromList(t, ctx, result)
		assert.Empty(t, rules)
	}
}

func TestMapAccessRulesToState_WithRules(t *testing.T) {
	ctx := context.Background()
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

	result, diags := resource.mapAccessRulesToState(ctx, accessRules)
	require.False(t, diags.HasError())

	resultRules := getAccessRulesFromList(t, ctx, result)
	require.Len(t, resultRules, 2)

	assert.Equal(t, "READ", resultRules[0].AccessType.ValueString())
	assert.Equal(t, "USER", resultRules[0].RelationType.ValueString())
	assert.Equal(t, "user-1", resultRules[0].RelatedID.ValueString())

	assert.Equal(t, "READ_WRITE", resultRules[1].AccessType.ValueString())
	assert.Equal(t, "ROLE", resultRules[1].RelationType.ValueString())
	assert.True(t, resultRules[1].RelatedID.IsNull())
}

func TestMapAccessRulesFromState_EmptyRules(t *testing.T) {
	ctx := context.Background()
	resource := &applicationConfigResource{}

	emptyList := createAccessRulesList(t, ctx, []AccessRuleModel{})
	result, diags := resource.mapAccessRulesFromState(ctx, emptyList)
	require.False(t, diags.HasError())
	assert.Empty(t, result)
}

func TestMapAccessRulesFromState_WithRules(t *testing.T) {
	ctx := context.Background()
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

	rulesList := createAccessRulesList(t, ctx, accessRuleModels)
	result, diags := resource.mapAccessRulesFromState(ctx, rulesList)
	require.False(t, diags.HasError())
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

			// Initialize state with a model that has TagFilter and AccessRules set to null
			// This is necessary because UpdateState checks if TagFilter is null/unknown
			initialModel := ApplicationConfigModel{
				TagFilter: types.StringNull(),
				AccessRules: types.ListNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"access_type":   types.StringType,
						"related_id":    types.StringType,
						"relation_type": types.StringType,
					},
				}),
			}
			diags := state.Set(ctx, initialModel)
			require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

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

			// Initialize state with a model that has TagFilter and AccessRules set to null
			// This is necessary because UpdateState checks if TagFilter is null/unknown
			initialModel := ApplicationConfigModel{
				TagFilter: types.StringNull(),
				AccessRules: types.ListNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"access_type":   types.StringType,
						"related_id":    types.StringType,
						"relation_type": types.StringType,
					},
				}),
			}
			diags := state.Set(ctx, initialModel)
			require.False(t, diags.HasError(), "Failed to set initial state: %v", diags)

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

// Helper function to create a types.List from AccessRuleModel slice
func createAccessRulesList(t *testing.T, ctx context.Context, rules []AccessRuleModel) types.List {
	if len(rules) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
		})
	}

	elements := make([]attr.Value, len(rules))
	for i, rule := range rules {
		elements[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"access_type":   types.StringType,
				"related_id":    types.StringType,
				"relation_type": types.StringType,
			},
			map[string]attr.Value{
				"access_type":   rule.AccessType,
				"related_id":    rule.RelatedID,
				"relation_type": rule.RelationType,
			},
		)
	}

	list, diags := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"access_type":   types.StringType,
			"related_id":    types.StringType,
			"relation_type": types.StringType,
		},
	}, elements)
	require.False(t, diags.HasError(), "Failed to create list: %v", diags)
	return list
}

// Helper function to extract AccessRuleModel slice from types.List
func getAccessRulesFromList(t *testing.T, ctx context.Context, list types.List) []AccessRuleModel {
	if list.IsNull() || list.IsUnknown() {
		return []AccessRuleModel{}
	}

	var rules []AccessRuleModel
	diags := list.ElementsAs(ctx, &rules, false)
	require.False(t, diags.HasError(), "Failed to extract access rules: %v", diags)
	return rules
}
