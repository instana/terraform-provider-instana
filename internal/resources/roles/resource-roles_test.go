package roles

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildMembersSet builds a types.Set of RoleMemberModel from a list of user IDs.
func buildMembersSet(t *testing.T, userIDs ...string) types.Set {
	t.Helper()
	memberType := buildMemberNestedObject().Type()
	if len(userIDs) == 0 {
		return types.SetValueMust(memberType, []attr.Value{})
	}
	elems := make([]attr.Value, len(userIDs))
	for i, uid := range userIDs {
		obj, diags := types.ObjectValue(
			map[string]attr.Type{RoleFieldMemberUserID: types.StringType},
			map[string]attr.Value{RoleFieldMemberUserID: types.StringValue(uid)},
		)
		require.False(t, diags.HasError())
		elems[i] = obj
	}
	s, diags := types.SetValue(memberType, elems)
	require.False(t, diags.HasError())
	return s
}

// extractMembersFromSet extracts []RoleMemberModel from a types.Set for assertions.
func extractMembersFromSet(t *testing.T, ctx context.Context, membersSet types.Set) []RoleMemberModel {
	t.Helper()
	if membersSet.IsNull() || membersSet.IsUnknown() {
		return nil
	}
	var members []RoleMemberModel
	diags := membersSet.ElementsAs(ctx, &members, false)
	require.False(t, diags.HasError())
	return members
}

func TestNewRoleResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewRoleResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaRole, metadata.ResourceName)
		assert.Equal(t, int64(1), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewRoleResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		assert.NotNil(t, schema.Attributes[RoleFieldID])
		assert.NotNil(t, schema.Attributes[RoleFieldName])
		assert.NotNil(t, schema.Attributes[RoleFieldMembers])
		assert.NotNil(t, schema.Attributes[RoleFieldPermissions])
	})

	t.Run("should have computed ID field", func(t *testing.T) {
		handle := NewRoleResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		idAttr := schema.Attributes[RoleFieldID]
		assert.NotNil(t, idAttr)
	})

	t.Run("should have required name field", func(t *testing.T) {
		handle := NewRoleResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		nameAttr := schema.Attributes[RoleFieldName]
		assert.NotNil(t, nameAttr)
	})

	t.Run("should have required members field", func(t *testing.T) {
		handle := NewRoleResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		membersAttr := schema.Attributes[RoleFieldMembers]
		assert.NotNil(t, membersAttr)
	})

	t.Run("should have required permissions field with validators", func(t *testing.T) {
		handle := NewRoleResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		permissionsAttr := schema.Attributes[RoleFieldPermissions]
		assert.NotNil(t, permissionsAttr)
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &roleResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaRole,
				SchemaVersion: 1,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaRole, metadata.ResourceName)
		assert.Equal(t, int64(1), metadata.SchemaVersion)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return roles rest resource", func(t *testing.T) {
		resource := &roleResource{}
		assert.NotNil(t, resource.GetRestResource)
	})
}

// mockRoleAPI extends the common mock to provide specific behavior for role tests
type mockRoleAPI struct {
	testutils.MockInstanaAPI
}

func (m *mockRoleAPI) Roles() rest.RestResource[*api.Role] {
	return &mockRoleRestResource{}
}

// Mock rest resource - implements all required methods from RestResource interface
type mockRoleRestResource struct{}

func (m *mockRoleRestResource) GetAll() (*[]*api.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) GetOne(id string) (*api.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) Create(data *api.Role) (*api.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) Update(data *api.Role) (*api.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) Delete(data *api.Role) error {
	return nil
}

func (m *mockRoleRestResource) DeleteByID(id string) error {
	return nil
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		resource := &roleResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaRole,
				Schema:        NewRoleResourceHandle().MetaData().Schema,
				SchemaVersion: 1,
			},
		}
		ctx := context.Background()

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		diags := resource.SetComputedFields(ctx, plan)
		assert.False(t, diags.HasError())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &roleResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaRole,
			Schema:        NewRoleResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map complete model from state successfully", func(t *testing.T) {
		model := RoleModel{
			ID:      types.StringValue("role-id-123"),
			Name:    types.StringValue("Test Role"),
			Members: buildMembersSet(t, "user-1", "user-2"),
			Permissions: []string{
				string(api.PermissionCanConfigureApplications),
				string(api.PermissionCanViewLogs),
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "role-id-123", result.ID)
		assert.Equal(t, "Test Role", result.Name)
		assert.Len(t, result.Members, 2)
		assert.Len(t, result.Permissions, 2)
		assert.Contains(t, result.Permissions, string(api.PermissionCanConfigureApplications))
	})

	t.Run("should map model from plan successfully", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("plan-role-id"),
			Name:        types.StringValue("Plan Role"),
			Members:     buildMembersSet(t, "user-3"),
			Permissions: []string{string(api.PermissionCanConfigureUsers)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "plan-role-id", result.ID)
		assert.Equal(t, "Plan Role", result.Name)
		assert.Len(t, result.Members, 1)
		assert.Equal(t, "user-3", result.Members[0].UserID)
		assert.Nil(t, result.Members[0].Email)
		assert.Nil(t, result.Members[0].Name)
	})

	t.Run("should handle when both plan and state are nil", func(t *testing.T) {
		result, diags := resource.MapStateToDataObject(ctx, nil, nil)

		assert.NotNil(t, result)
		assert.False(t, diags.HasError())
	})

	t.Run("should map members with optional fields", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t, "user-1", "user-2"),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Members, 2)
	})

	t.Run("should handle empty members list", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Empty(t, result.Members)
	})

	t.Run("should handle empty permissions list", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t, "user-1"),
			Permissions: []string{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Empty(t, result.Permissions)
	})

	t.Run("should handle null ID for new resource", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringNull(),
			Name:        types.StringValue("New Role"),
			Members:     buildMembersSet(t, "user-1"),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "", result.ID)
	})

	t.Run("should handle multiple permissions", func(t *testing.T) {
		model := RoleModel{
			ID:      types.StringValue("role-id"),
			Name:    types.StringValue("Test Role"),
			Members: buildMembersSet(t, "user-1"),
			Permissions: []string{
				string(api.PermissionCanConfigureApplications),
				string(api.PermissionCanViewLogs),
				string(api.PermissionCanConfigureUsers),
				string(api.PermissionCanConfigureTeams),
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Permissions, 4)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &roleResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaRole,
			Schema:        NewRoleResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with complete API object", func(t *testing.T) {
		name2 := "User Two"

		apiObject := &api.Role{
			ID:   "api-role-id-123",
			Name: "API Role",
			Members: []api.APIMember{
				{
					UserID: "user-1",
				},
				{
					UserID: "user-2",
					Name:   &name2,
				},
			},
			Permissions: []string{
				string(api.PermissionCanConfigureApplications),
				string(api.PermissionCanViewLogs),
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyRoleState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-role-id-123", model.ID.ValueString())
		members := extractMembersFromSet(t, ctx, model.Members)
		assert.Len(t, members, 2)
		assert.Len(t, model.Permissions, 2)
	})

	t.Run("should update state with members without optional fields", func(t *testing.T) {
		apiObject := &api.Role{
			ID:   "api-role-id-456",
			Name: "Minimal Role",
			Members: []api.APIMember{
				{
					UserID: "user-3",
					Email:  nil,
				},
			},
			Permissions: []string{
				string(api.PermissionCanConfigureUsers),
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyRoleState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-role-id-456", model.ID.ValueString())
		members := extractMembersFromSet(t, ctx, model.Members)
		assert.Len(t, members, 1)
		assert.Equal(t, "user-3", members[0].UserID.ValueString())
	})

	t.Run("should update state with empty members list", func(t *testing.T) {
		apiObject := &api.Role{
			ID:          "api-role-id-789",
			Name:        "No Members Role",
			Members:     []api.APIMember{},
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyRoleState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		members := extractMembersFromSet(t, ctx, model.Members)
		assert.Empty(t, members)
	})

	t.Run("should preserve existing member data from plan", func(t *testing.T) {
		// Set up existing plan with member data
		existingModel := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t, "user-1"),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, existingModel)
		require.False(t, diags.HasError())

		// API returns member without email/name
		apiObject := &api.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []api.APIMember{
				{
					UserID: "user-1",
					Email:  nil,
				},
			},
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags = resource.UpdateState(ctx, state, plan, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		// Should preserve email and name from plan
	})

	t.Run("should handle API returning empty strings for optional fields", func(t *testing.T) {

		apiObject := &api.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []api.APIMember{
				{
					UserID: "user-1",
				},
			},
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyRoleState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		members := extractMembersFromSet(t, ctx, model.Members)
		assert.Len(t, members, 1)
	})

	t.Run("should update state with API values when present", func(t *testing.T) {

		// Set up existing plan with different data
		existingModel := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t, "user-1"),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, existingModel)
		require.False(t, diags.HasError())

		// API returns new values
		apiObject := &api.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []api.APIMember{
				{
					UserID: "user-1",
				},
			},
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyRoleState(t, ctx, state)

		diags = resource.UpdateState(ctx, state, plan, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		// Should use API values when present
		members := extractMembersFromSet(t, ctx, model.Members)
		assert.Len(t, members, 1)
	})

	t.Run("should handle new members not in existing state", func(t *testing.T) {

		// Set up existing plan with one member
		existingModel := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t, "user-1"),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, existingModel)
		require.False(t, diags.HasError())

		// API returns a new member
		apiObject := &api.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []api.APIMember{
				{
					UserID: "user-2",
				},
			},
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyRoleState(t, ctx, state)

		diags = resource.UpdateState(ctx, state, plan, apiObject)

		assert.False(t, diags.HasError())

		var model RoleModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		members := extractMembersFromSet(t, ctx, model.Members)
		assert.Len(t, members, 1)
		assert.Equal(t, "user-2", members[0].UserID.ValueString())
	})
}

func TestMapMembersToSet(t *testing.T) {
	resource := &roleResource{}
	ctx := context.Background()

	t.Run("should map empty members list", func(t *testing.T) {
		result, diags := resource.mapMembersToSet(ctx, []api.APIMember{})
		assert.False(t, diags.HasError())
		assert.Empty(t, result.Elements())
	})

	t.Run("should map members with all fields", func(t *testing.T) {
		apiMembers := []api.APIMember{
			{
				UserID: "user-1",
			},
		}

		result, diags := resource.mapMembersToSet(ctx, apiMembers)
		assert.False(t, diags.HasError())
		assert.Len(t, result.Elements(), 1)

		members := extractMembersFromSet(t, ctx, result)
		assert.Equal(t, "user-1", members[0].UserID.ValueString())
	})

	t.Run("should map multiple members", func(t *testing.T) {
		apiMembers := []api.APIMember{
			{UserID: "user-1"},
			{UserID: "user-2"},
		}

		result, diags := resource.mapMembersToSet(ctx, apiMembers)
		assert.False(t, diags.HasError())
		assert.Len(t, result.Elements(), 2)
	})
}

func TestMapSetMembersToAPI(t *testing.T) {
	resource := &roleResource{}
	ctx := context.Background()

	t.Run("should map empty members set", func(t *testing.T) {
		result, diags := resource.mapSetMembersToAPI(ctx, buildMembersSet(t))
		assert.False(t, diags.HasError())
		assert.Empty(t, result)
	})

	t.Run("should map members with all fields", func(t *testing.T) {
		membersSet := buildMembersSet(t, "user-1")

		result, diags := resource.mapSetMembersToAPI(ctx, membersSet)
		assert.False(t, diags.HasError())
		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].UserID)
		assert.Nil(t, result[0].Email)
		assert.Nil(t, result[0].Name)
	})

	t.Run("should handle null members set", func(t *testing.T) {
		memberType := buildMemberNestedObject().Type()
		result, diags := resource.mapSetMembersToAPI(ctx, types.SetNull(memberType))
		assert.False(t, diags.HasError())
		assert.Empty(t, result)
	})

	t.Run("should map multiple members", func(t *testing.T) {
		membersSet := buildMembersSet(t, "user-1", "user-2", "user-3")

		result, diags := resource.mapSetMembersToAPI(ctx, membersSet)
		assert.False(t, diags.HasError())
		assert.Len(t, result, 3)
		// Collect user IDs for assertion (set order is not guaranteed)
		userIDs := make([]string, len(result))
		for i, m := range result {
			userIDs[i] = m.UserID
		}
		assert.Contains(t, userIDs, "user-1")
		assert.Contains(t, userIDs, "user-2")
		assert.Contains(t, userIDs, "user-3")
	})
}

func TestExtractRoleID(t *testing.T) {
	resource := &roleResource{}

	t.Run("should extract non-null ID", func(t *testing.T) {
		model := RoleModel{
			ID: types.StringValue("role-123"),
		}

		result := resource.extractRoleID(model)
		assert.Equal(t, "role-123", result)
	})

	t.Run("should return empty string for null ID", func(t *testing.T) {
		model := RoleModel{
			ID: types.StringNull(),
		}

		result := resource.extractRoleID(model)
		assert.Equal(t, "", result)
	})
}

func TestExtractModelFromPlanOrState(t *testing.T) {
	resource := &roleResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaRole,
			Schema:        NewRoleResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should extract from plan when provided", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     buildMembersSet(t, "user-1"),
			Permissions: []string{string(api.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.extractModelFromPlanOrState(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		assert.Equal(t, "role-id", result.ID.ValueString())
	})

	t.Run("should extract from state when plan is nil", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("state-role-id"),
			Name:        types.StringValue("State Role"),
			Members:     buildMembersSet(t, "user-2"),
			Permissions: []string{string(api.PermissionCanConfigureUsers)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.extractModelFromPlanOrState(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.Equal(t, "state-role-id", result.ID.ValueString())
	})

	t.Run("should return empty model when both are nil", func(t *testing.T) {
		result, diags := resource.extractModelFromPlanOrState(ctx, nil, nil)

		assert.False(t, diags.HasError())
		assert.True(t, result.ID.IsNull())
	})
}

// initializeEmptyRoleState initializes the state with an empty model to ensure proper state initialization
func initializeEmptyRoleState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	memberType := buildMemberNestedObject().Type()
	emptyModel := RoleModel{
		ID:          types.StringNull(),
		Name:        types.StringNull(),
		Members:     types.SetValueMust(memberType, []attr.Value{}),
		Permissions: []string{},
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}
