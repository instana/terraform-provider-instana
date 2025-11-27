package roles

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

		mockAPI := &mockRoleAPI{}
		restResource := resource.GetRestResource(mockAPI)

		assert.NotNil(t, restResource)
	})
}

// mockRoleAPI extends the common mock to provide specific behavior for role tests
type mockRoleAPI struct {
	testutils.MockInstanaAPI
}

func (m *mockRoleAPI) Roles() restapi.RestResource[*restapi.Role] {
	return &mockRoleRestResource{}
}

// Mock rest resource - implements all required methods from RestResource interface
type mockRoleRestResource struct{}

func (m *mockRoleRestResource) GetAll() (*[]*restapi.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) GetOne(id string) (*restapi.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) Create(data *restapi.Role) (*restapi.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) Update(data *restapi.Role) (*restapi.Role, error) {
	return nil, nil
}

func (m *mockRoleRestResource) Delete(data *restapi.Role) error {
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
			ID:   types.StringValue("role-id-123"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("user1@example.com"),
					Name:   types.StringValue("User One"),
				},
				{
					UserID: types.StringValue("user-2"),
					Email:  types.StringValue("user2@example.com"),
					Name:   types.StringValue("User Two"),
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureApplications),
				string(restapi.PermissionCanViewLogs),
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
		assert.Equal(t, "user-1", result.Members[0].UserID)
		assert.NotNil(t, result.Members[0].Email)
		assert.Equal(t, "user1@example.com", *result.Members[0].Email)
		assert.NotNil(t, result.Members[0].Name)
		assert.Equal(t, "User One", *result.Members[0].Name)
		assert.Len(t, result.Permissions, 2)
		assert.Contains(t, result.Permissions, string(restapi.PermissionCanConfigureApplications))
	})

	t.Run("should map model from plan successfully", func(t *testing.T) {
		model := RoleModel{
			ID:   types.StringValue("plan-role-id"),
			Name: types.StringValue("Plan Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-3"),
					Email:  types.StringNull(),
					Name:   types.StringNull(),
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureUsers),
			},
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
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("user1@example.com"),
					Name:   types.StringNull(),
				},
				{
					UserID: types.StringValue("user-2"),
					Email:  types.StringNull(),
					Name:   types.StringValue("User Two"),
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureApplications),
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
		assert.Len(t, result.Members, 2)
		assert.NotNil(t, result.Members[0].Email)
		assert.Nil(t, result.Members[0].Name)
		assert.Nil(t, result.Members[1].Email)
		assert.NotNil(t, result.Members[1].Name)
	})

	t.Run("should handle empty members list", func(t *testing.T) {
		model := RoleModel{
			ID:          types.StringValue("role-id"),
			Name:        types.StringValue("Test Role"),
			Members:     []RoleMemberModel{},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringNull(),
					Name:   types.StringNull(),
				},
			},
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
			ID:   types.StringNull(),
			Name: types.StringValue("New Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringNull(),
					Name:   types.StringNull(),
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringNull(),
					Name:   types.StringNull(),
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureApplications),
				string(restapi.PermissionCanViewLogs),
				string(restapi.PermissionCanConfigureUsers),
				string(restapi.PermissionCanConfigureTeams),
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
		email1 := "user1@example.com"
		name1 := "User One"
		email2 := "user2@example.com"
		name2 := "User Two"

		apiObject := &restapi.Role{
			ID:   "api-role-id-123",
			Name: "API Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  &email1,
					Name:   &name1,
				},
				{
					UserID: "user-2",
					Email:  &email2,
					Name:   &name2,
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureApplications),
				string(restapi.PermissionCanViewLogs),
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
		assert.Equal(t, "API Role", model.Name.ValueString())
		assert.Len(t, model.Members, 2)
		assert.Equal(t, "user-1", model.Members[0].UserID.ValueString())
		assert.Equal(t, "user1@example.com", model.Members[0].Email.ValueString())
		assert.Equal(t, "User One", model.Members[0].Name.ValueString())
		assert.Len(t, model.Permissions, 2)
	})

	t.Run("should update state with members without optional fields", func(t *testing.T) {
		apiObject := &restapi.Role{
			ID:   "api-role-id-456",
			Name: "Minimal Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-3",
					Email:  nil,
					Name:   nil,
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureUsers),
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
		assert.Len(t, model.Members, 1)
		assert.Equal(t, "user-3", model.Members[0].UserID.ValueString())
		assert.True(t, model.Members[0].Email.IsNull())
		assert.True(t, model.Members[0].Name.IsNull())
	})

	t.Run("should update state with empty members list", func(t *testing.T) {
		apiObject := &restapi.Role{
			ID:          "api-role-id-789",
			Name:        "No Members Role",
			Members:     []restapi.APIMember{},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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

		assert.Empty(t, model.Members)
	})

	t.Run("should preserve existing member data from plan", func(t *testing.T) {
		// Set up existing plan with member data
		existingModel := RoleModel{
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("user1@example.com"),
					Name:   types.StringValue("User One"),
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, existingModel)
		require.False(t, diags.HasError())

		// API returns member without email/name
		apiObject := &restapi.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  nil,
					Name:   nil,
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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
		assert.Equal(t, "user1@example.com", model.Members[0].Email.ValueString())
		assert.Equal(t, "User One", model.Members[0].Name.ValueString())
	})

	t.Run("should handle API returning empty strings for optional fields", func(t *testing.T) {
		emptyEmail := ""
		emptyName := ""

		apiObject := &restapi.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  &emptyEmail,
					Name:   &emptyName,
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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

		assert.Len(t, model.Members, 1)
		assert.True(t, model.Members[0].Email.IsNull())
		assert.True(t, model.Members[0].Name.IsNull())
	})

	t.Run("should update state with API values when present", func(t *testing.T) {
		email := "api@example.com"
		name := "API User"

		// Set up existing plan with different data
		existingModel := RoleModel{
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("old@example.com"),
					Name:   types.StringValue("Old Name"),
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, existingModel)
		require.False(t, diags.HasError())

		// API returns new values
		apiObject := &restapi.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  &email,
					Name:   &name,
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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
		assert.Len(t, model.Members, 1)
		assert.Equal(t, "api@example.com", model.Members[0].Email.ValueString())
		assert.Equal(t, "API User", model.Members[0].Name.ValueString())
	})

	t.Run("should handle new members not in existing state", func(t *testing.T) {
		email := "new@example.com"
		name := "New User"

		// Set up existing plan with one member
		existingModel := RoleModel{
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("user1@example.com"),
					Name:   types.StringValue("User One"),
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, existingModel)
		require.False(t, diags.HasError())

		// API returns a new member
		apiObject := &restapi.Role{
			ID:   "role-id",
			Name: "Test Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-2",
					Email:  &email,
					Name:   &name,
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
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

		assert.Len(t, model.Members, 1)
		assert.Equal(t, "user-2", model.Members[0].UserID.ValueString())
		assert.Equal(t, "new@example.com", model.Members[0].Email.ValueString())
		assert.Equal(t, "New User", model.Members[0].Name.ValueString())
	})
}

func TestMapMembersToModel(t *testing.T) {
	resource := &roleResource{}

	t.Run("should map empty members list", func(t *testing.T) {
		result := resource.mapMembersToModel([]restapi.APIMember{}, []RoleMemberModel{})
		assert.Empty(t, result)
	})

	t.Run("should map members with all fields", func(t *testing.T) {
		email := "user@example.com"
		name := "User Name"

		apiMembers := []restapi.APIMember{
			{
				UserID: "user-1",
				Email:  &email,
				Name:   &name,
			},
		}

		result := resource.mapMembersToModel(apiMembers, []RoleMemberModel{})

		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].UserID.ValueString())
		assert.Equal(t, "user@example.com", result[0].Email.ValueString())
		assert.Equal(t, "User Name", result[0].Name.ValueString())
	})

	t.Run("should preserve existing member data when API returns nil", func(t *testing.T) {
		existingMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("existing@example.com"),
				Name:   types.StringValue("Existing Name"),
			},
		}

		apiMembers := []restapi.APIMember{
			{
				UserID: "user-1",
				Email:  nil,
				Name:   nil,
			},
		}

		result := resource.mapMembersToModel(apiMembers, existingMembers)

		assert.Len(t, result, 1)
		assert.Equal(t, "existing@example.com", result[0].Email.ValueString())
		assert.Equal(t, "Existing Name", result[0].Name.ValueString())
	})

	t.Run("should handle empty strings from API", func(t *testing.T) {
		emptyEmail := ""
		emptyName := ""

		existingMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("existing@example.com"),
				Name:   types.StringValue("Existing Name"),
			},
		}

		apiMembers := []restapi.APIMember{
			{
				UserID: "user-1",
				Email:  &emptyEmail,
				Name:   &emptyName,
			},
		}

		result := resource.mapMembersToModel(apiMembers, existingMembers)

		assert.Len(t, result, 1)
		// Should preserve existing values when API returns empty strings
		assert.Equal(t, "existing@example.com", result[0].Email.ValueString())
		assert.Equal(t, "Existing Name", result[0].Name.ValueString())
	})

	t.Run("should use API values when present and non-empty", func(t *testing.T) {
		newEmail := "new@example.com"
		newName := "New Name"

		existingMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("old@example.com"),
				Name:   types.StringValue("Old Name"),
			},
		}

		apiMembers := []restapi.APIMember{
			{
				UserID: "user-1",
				Email:  &newEmail,
				Name:   &newName,
			},
		}

		result := resource.mapMembersToModel(apiMembers, existingMembers)

		assert.Len(t, result, 1)
		assert.Equal(t, "new@example.com", result[0].Email.ValueString())
		assert.Equal(t, "New Name", result[0].Name.ValueString())
	})

	t.Run("should handle multiple members", func(t *testing.T) {
		email1 := "user1@example.com"
		name1 := "User One"
		email2 := "user2@example.com"

		apiMembers := []restapi.APIMember{
			{
				UserID: "user-1",
				Email:  &email1,
				Name:   &name1,
			},
			{
				UserID: "user-2",
				Email:  &email2,
				Name:   nil,
			},
		}

		result := resource.mapMembersToModel(apiMembers, []RoleMemberModel{})

		assert.Len(t, result, 2)
		assert.Equal(t, "user-1", result[0].UserID.ValueString())
		assert.Equal(t, "user1@example.com", result[0].Email.ValueString())
		assert.Equal(t, "User One", result[0].Name.ValueString())
		assert.Equal(t, "user-2", result[1].UserID.ValueString())
		assert.Equal(t, "user2@example.com", result[1].Email.ValueString())
		assert.True(t, result[1].Name.IsNull())
	})
}

func TestMapModelMembersToAPI(t *testing.T) {
	resource := &roleResource{}

	t.Run("should map empty members list", func(t *testing.T) {
		result := resource.mapModelMembersToAPI([]RoleMemberModel{})
		assert.Empty(t, result)
	})

	t.Run("should map members with all fields", func(t *testing.T) {
		modelMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("user@example.com"),
				Name:   types.StringValue("User Name"),
			},
		}

		result := resource.mapModelMembersToAPI(modelMembers)

		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].UserID)
		assert.NotNil(t, result[0].Email)
		assert.Equal(t, "user@example.com", *result[0].Email)
		assert.NotNil(t, result[0].Name)
		assert.Equal(t, "User Name", *result[0].Name)
	})

	t.Run("should handle null email and name", func(t *testing.T) {
		modelMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringNull(),
				Name:   types.StringNull(),
			},
		}

		result := resource.mapModelMembersToAPI(modelMembers)

		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].UserID)
		assert.Nil(t, result[0].Email)
		assert.Nil(t, result[0].Name)
	})

	t.Run("should handle unknown email and name", func(t *testing.T) {
		modelMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringUnknown(),
				Name:   types.StringUnknown(),
			},
		}

		result := resource.mapModelMembersToAPI(modelMembers)

		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].UserID)
		assert.Nil(t, result[0].Email)
		assert.Nil(t, result[0].Name)
	})

	t.Run("should handle mixed null and non-null fields", func(t *testing.T) {
		modelMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("user1@example.com"),
				Name:   types.StringNull(),
			},
			{
				UserID: types.StringValue("user-2"),
				Email:  types.StringNull(),
				Name:   types.StringValue("User Two"),
			},
		}

		result := resource.mapModelMembersToAPI(modelMembers)

		assert.Len(t, result, 2)
		assert.NotNil(t, result[0].Email)
		assert.Nil(t, result[0].Name)
		assert.Nil(t, result[1].Email)
		assert.NotNil(t, result[1].Name)
	})

	t.Run("should handle multiple members", func(t *testing.T) {
		modelMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("user1@example.com"),
				Name:   types.StringValue("User One"),
			},
			{
				UserID: types.StringValue("user-2"),
				Email:  types.StringValue("user2@example.com"),
				Name:   types.StringValue("User Two"),
			},
			{
				UserID: types.StringValue("user-3"),
				Email:  types.StringNull(),
				Name:   types.StringNull(),
			},
		}

		result := resource.mapModelMembersToAPI(modelMembers)

		assert.Len(t, result, 3)
		assert.Equal(t, "user-1", result[0].UserID)
		assert.Equal(t, "user-2", result[1].UserID)
		assert.Equal(t, "user-3", result[2].UserID)
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

func TestBuildRoleModelFromAPIResponse(t *testing.T) {
	resource := &roleResource{}

	t.Run("should build model with all fields", func(t *testing.T) {
		email := "user@example.com"
		name := "User Name"

		apiRole := &restapi.Role{
			ID:   "role-123",
			Name: "Test Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  &email,
					Name:   &name,
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureApplications),
			},
		}

		result := resource.buildRoleModelFromAPIResponse(apiRole, []RoleMemberModel{})

		assert.Equal(t, "role-123", result.ID.ValueString())
		assert.Equal(t, "Test Role", result.Name.ValueString())
		assert.Len(t, result.Members, 1)
		assert.Len(t, result.Permissions, 1)
	})

	t.Run("should preserve existing member data", func(t *testing.T) {
		existingMembers := []RoleMemberModel{
			{
				UserID: types.StringValue("user-1"),
				Email:  types.StringValue("existing@example.com"),
				Name:   types.StringValue("Existing Name"),
			},
		}

		apiRole := &restapi.Role{
			ID:   "role-123",
			Name: "Test Role",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  nil,
					Name:   nil,
				},
			},
			Permissions: []string{
				string(restapi.PermissionCanConfigureApplications),
			},
		}

		result := resource.buildRoleModelFromAPIResponse(apiRole, existingMembers)

		assert.Equal(t, "existing@example.com", result.Members[0].Email.ValueString())
		assert.Equal(t, "Existing Name", result.Members[0].Name.ValueString())
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
			ID:   types.StringValue("role-id"),
			Name: types.StringValue("Test Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringNull(),
					Name:   types.StringNull(),
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureApplications)},
		}

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.extractModelFromPlanOrState(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		assert.Equal(t, "role-id", result.ID.ValueString())
		assert.Equal(t, "Test Role", result.Name.ValueString())
	})

	t.Run("should extract from state when plan is nil", func(t *testing.T) {
		model := RoleModel{
			ID:   types.StringValue("state-role-id"),
			Name: types.StringValue("State Role"),
			Members: []RoleMemberModel{
				{
					UserID: types.StringValue("user-2"),
					Email:  types.StringNull(),
					Name:   types.StringNull(),
				},
			},
			Permissions: []string{string(restapi.PermissionCanConfigureUsers)},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.extractModelFromPlanOrState(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.Equal(t, "state-role-id", result.ID.ValueString())
		assert.Equal(t, "State Role", result.Name.ValueString())
	})

	t.Run("should return empty model when both are nil", func(t *testing.T) {
		result, diags := resource.extractModelFromPlanOrState(ctx, nil, nil)

		assert.False(t, diags.HasError())
		assert.True(t, result.ID.IsNull())
		assert.True(t, result.Name.IsNull())
	})
}

// initializeEmptyRoleState initializes the state with an empty model to ensure proper state initialization
func initializeEmptyRoleState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := RoleModel{
		ID:          types.StringNull(),
		Name:        types.StringNull(),
		Members:     []RoleMemberModel{},
		Permissions: []string{},
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}

// Made with Bob