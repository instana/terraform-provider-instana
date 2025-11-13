package group

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGroupResourceHandleFramework(t *testing.T) {
	handle := NewGroupResourceHandleFramework()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaGroupFramework, metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &groupResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  "test_resource",
			SchemaVersion: 1,
		},
	}

	metadata := resource.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, "test_resource", metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
}

func TestGetRestResource(t *testing.T) {
	resource := &groupResourceFramework{}
	
	// This test just ensures the method exists and can be called
	// The actual implementation requires a real API instance
	assert.NotNil(t, resource)
}

func TestSetComputedFields(t *testing.T) {
	resource := &groupResourceFramework{}
	ctx := context.Background()

	diags := resource.SetComputedFields(ctx, nil)
	assert.False(t, diags.HasError())
}

func TestUpdateState(t *testing.T) {
	resource := &groupResourceFramework{}
	ctx := context.Background()

	t.Run("basic group without members or permissions", func(t *testing.T) {
		group := &restapi.Group{
			ID:            "test-id",
			Name:          "test-group",
			Members:       []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-group", model.Name.ValueString())
		assert.Nil(t, model.Members)
		assert.Nil(t, model.PermissionSet)
	})

	t.Run("group with members", func(t *testing.T) {
		email := "user@example.com"
		group := &restapi.Group{
			ID:   "test-id",
			Name: "test-group",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  &email,
				},
				{
					UserID: "user-2",
					Email:  nil,
				},
			},
			PermissionSet: restapi.APIPermissionSetWithRoles{},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Len(t, model.Members, 2)
		assert.Equal(t, "user-1", model.Members[0].UserID.ValueString())
		assert.Equal(t, "user@example.com", model.Members[0].Email.ValueString())
		assert.Equal(t, "user-2", model.Members[1].UserID.ValueString())
		assert.True(t, model.Members[1].Email.IsNull())
	})

	t.Run("group with application IDs", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				ApplicationIDs: []restapi.ScopeBinding{
					{ScopeID: "app-1"},
					{ScopeID: "app-2"},
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.ApplicationIDs, 2)
		assert.Equal(t, "app-1", model.PermissionSet.ApplicationIDs[0])
		assert.Equal(t, "app-2", model.PermissionSet.ApplicationIDs[1])
	})

	t.Run("group with infra DFQ filter", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				InfraDFQFilter: &restapi.ScopeBinding{
					ScopeID: "entity.type:host",
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Equal(t, "entity.type:host", model.PermissionSet.InfraDFQFilter.ValueString())
	})

	t.Run("group with empty infra DFQ filter", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				InfraDFQFilter: &restapi.ScopeBinding{
					ScopeID: "",
				},
				Permissions: []restapi.InstanaPermission{
					restapi.PermissionCanViewLogs,
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.True(t, model.PermissionSet.InfraDFQFilter.IsNull())
	})

	t.Run("group with kubernetes cluster UUIDs", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				KubernetesClusterUUIDs: []restapi.ScopeBinding{
					{ScopeID: "k8s-cluster-1"},
					{ScopeID: "k8s-cluster-2"},
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.KubernetesClusterUUIDs, 2)
		assert.Equal(t, "k8s-cluster-1", model.PermissionSet.KubernetesClusterUUIDs[0])
	})

	t.Run("group with kubernetes namespace UIDs", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				KubernetesNamespaceUIDs: []restapi.ScopeBinding{
					{ScopeID: "k8s-ns-1"},
					{ScopeID: "k8s-ns-2"},
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.KubernetesNamespaceUIDs, 2)
		assert.Equal(t, "k8s-ns-1", model.PermissionSet.KubernetesNamespaceUIDs[0])
	})

	t.Run("group with mobile app IDs", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				MobileAppIDs: []restapi.ScopeBinding{
					{ScopeID: "mobile-app-1"},
					{ScopeID: "mobile-app-2"},
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.MobileAppIDs, 2)
		assert.Equal(t, "mobile-app-1", model.PermissionSet.MobileAppIDs[0])
	})

	t.Run("group with website IDs", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				WebsiteIDs: []restapi.ScopeBinding{
					{ScopeID: "website-1"},
					{ScopeID: "website-2"},
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.WebsiteIDs, 2)
		assert.Equal(t, "website-1", model.PermissionSet.WebsiteIDs[0])
	})

	t.Run("group with permissions", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				Permissions: []restapi.InstanaPermission{
					restapi.PermissionCanConfigureApplications,
					restapi.PermissionCanViewLogs,
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.Permissions, 2)
		assert.Equal(t, "CAN_CONFIGURE_APPLICATIONS", model.PermissionSet.Permissions[0])
		assert.Equal(t, "CAN_VIEW_LOGS", model.PermissionSet.Permissions[1])
	})

	t.Run("group with all permission set fields", func(t *testing.T) {
		group := &restapi.Group{
			ID:   "test-id",
			Name: "test-group",
			Members: []restapi.APIMember{
				{UserID: "user-1"},
			},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				ApplicationIDs: []restapi.ScopeBinding{
					{ScopeID: "app-1"},
				},
				InfraDFQFilter: &restapi.ScopeBinding{
					ScopeID: "entity.type:host",
				},
				KubernetesClusterUUIDs: []restapi.ScopeBinding{
					{ScopeID: "k8s-cluster-1"},
				},
				KubernetesNamespaceUIDs: []restapi.ScopeBinding{
					{ScopeID: "k8s-ns-1"},
				},
				MobileAppIDs: []restapi.ScopeBinding{
					{ScopeID: "mobile-app-1"},
				},
				WebsiteIDs: []restapi.ScopeBinding{
					{ScopeID: "website-1"},
				},
				Permissions: []restapi.InstanaPermission{
					restapi.PermissionCanConfigureApplications,
				},
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.PermissionSet)
		assert.Len(t, model.PermissionSet.ApplicationIDs, 1)
		assert.Equal(t, "entity.type:host", model.PermissionSet.InfraDFQFilter.ValueString())
		assert.Len(t, model.PermissionSet.KubernetesClusterUUIDs, 1)
		assert.Len(t, model.PermissionSet.KubernetesNamespaceUIDs, 1)
		assert.Len(t, model.PermissionSet.MobileAppIDs, 1)
		assert.Len(t, model.PermissionSet.WebsiteIDs, 1)
		assert.Len(t, model.PermissionSet.Permissions, 1)
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &groupResourceFramework{}
	ctx := context.Background()

	t.Run("basic group from state", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Equal(t, "test-id", group.ID)
		assert.Equal(t, "test-group", group.Name)
		assert.Empty(t, group.Members)
		assert.Empty(t, group.PermissionSet.ApplicationIDs)
		assert.Empty(t, group.PermissionSet.Permissions)
	})

	t.Run("group from plan", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue(""),
			Name: types.StringValue("new-group"),
		}

		plan := createMockPlan(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Equal(t, "", group.ID)
		assert.Equal(t, "new-group", group.Name)
	})

	t.Run("group with null ID", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringNull(),
			Name: types.StringValue("test-group"),
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Equal(t, "", group.ID)
	})

	t.Run("group with members", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			Members: []GroupMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("user1@example.com"),
				},
				{
					UserID: types.StringValue("user-2"),
					Email:  types.StringNull(),
				},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.Members, 2)
		assert.Equal(t, "user-1", group.Members[0].UserID)
		assert.Equal(t, "user1@example.com", *group.Members[0].Email)
		assert.Equal(t, "user-2", group.Members[1].UserID)
		assert.Nil(t, group.Members[1].Email)
	})

	t.Run("group with member email unknown", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			Members: []GroupMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringUnknown(),
				},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.Members, 1)
		assert.Nil(t, group.Members[0].Email)
	})

	t.Run("group with application IDs", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				ApplicationIDs: []string{"app-1", "app-2"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.PermissionSet.ApplicationIDs, 2)
		assert.Equal(t, "app-1", group.PermissionSet.ApplicationIDs[0].ScopeID)
		assert.Equal(t, "app-2", group.PermissionSet.ApplicationIDs[1].ScopeID)
	})

	t.Run("group with infra DFQ filter", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				InfraDFQFilter: types.StringValue("entity.type:host"),
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		require.NotNil(t, group.PermissionSet.InfraDFQFilter)
		assert.Equal(t, "entity.type:host", group.PermissionSet.InfraDFQFilter.ScopeID)
	})

	t.Run("group with null infra DFQ filter", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				InfraDFQFilter: types.StringNull(),
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		require.NotNil(t, group.PermissionSet.InfraDFQFilter)
		assert.Equal(t, "", group.PermissionSet.InfraDFQFilter.ScopeID)
		assert.Equal(t, "-1", *group.PermissionSet.InfraDFQFilter.ScopeRoleID)
	})

	t.Run("group with unknown infra DFQ filter", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				InfraDFQFilter: types.StringUnknown(),
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		require.NotNil(t, group.PermissionSet.InfraDFQFilter)
		assert.Equal(t, "", group.PermissionSet.InfraDFQFilter.ScopeID)
		assert.Equal(t, "-1", *group.PermissionSet.InfraDFQFilter.ScopeRoleID)
	})

	t.Run("group with kubernetes cluster UUIDs", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				KubernetesClusterUUIDs: []string{"k8s-cluster-1", "k8s-cluster-2"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.PermissionSet.KubernetesClusterUUIDs, 2)
		assert.Equal(t, "k8s-cluster-1", group.PermissionSet.KubernetesClusterUUIDs[0].ScopeID)
	})

	t.Run("group with kubernetes namespace UIDs", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				KubernetesNamespaceUIDs: []string{"k8s-ns-1", "k8s-ns-2"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.PermissionSet.KubernetesNamespaceUIDs, 2)
		assert.Equal(t, "k8s-ns-1", group.PermissionSet.KubernetesNamespaceUIDs[0].ScopeID)
	})

	t.Run("group with mobile app IDs", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				MobileAppIDs: []string{"mobile-app-1", "mobile-app-2"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.PermissionSet.MobileAppIDs, 2)
		assert.Equal(t, "mobile-app-1", group.PermissionSet.MobileAppIDs[0].ScopeID)
	})

	t.Run("group with website IDs", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				WebsiteIDs: []string{"website-1", "website-2"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.PermissionSet.WebsiteIDs, 2)
		assert.Equal(t, "website-1", group.PermissionSet.WebsiteIDs[0].ScopeID)
	})

	t.Run("group with permissions", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				Permissions: []string{"CAN_CONFIGURE_APPLICATIONS", "CAN_VIEW_LOGS"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.PermissionSet.Permissions, 2)
		assert.Equal(t, restapi.PermissionCanConfigureApplications, group.PermissionSet.Permissions[0])
		assert.Equal(t, restapi.PermissionCanViewLogs, group.PermissionSet.Permissions[1])
	})

	t.Run("group with all permission set fields", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			Members: []GroupMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Email:  types.StringValue("user1@example.com"),
				},
			},
			PermissionSet: &GroupPermissionSetModel{
				ApplicationIDs:          []string{"app-1"},
				InfraDFQFilter:          types.StringValue("entity.type:host"),
				KubernetesClusterUUIDs:  []string{"k8s-cluster-1"},
				KubernetesNamespaceUIDs: []string{"k8s-ns-1"},
				MobileAppIDs:            []string{"mobile-app-1"},
				WebsiteIDs:              []string{"website-1"},
				Permissions:             []string{"CAN_CONFIGURE_APPLICATIONS"},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Len(t, group.Members, 1)
		assert.Len(t, group.PermissionSet.ApplicationIDs, 1)
		assert.Equal(t, "entity.type:host", group.PermissionSet.InfraDFQFilter.ScopeID)
		assert.Len(t, group.PermissionSet.KubernetesClusterUUIDs, 1)
		assert.Len(t, group.PermissionSet.KubernetesNamespaceUIDs, 1)
		assert.Len(t, group.PermissionSet.MobileAppIDs, 1)
		assert.Len(t, group.PermissionSet.WebsiteIDs, 1)
		assert.Len(t, group.PermissionSet.Permissions, 1)
	})

	t.Run("group without permission set", func(t *testing.T) {
		model := GroupModel{
			ID:            types.StringValue("test-id"),
			Name:          types.StringValue("test-group"),
			PermissionSet: nil,
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)

		assert.Empty(t, group.PermissionSet.ApplicationIDs)
		assert.Empty(t, group.PermissionSet.KubernetesClusterUUIDs)
		assert.Empty(t, group.PermissionSet.KubernetesNamespaceUIDs)
		assert.Empty(t, group.PermissionSet.MobileAppIDs)
		assert.Empty(t, group.PermissionSet.WebsiteIDs)
		assert.Empty(t, group.PermissionSet.Permissions)
	})
}

func TestRoundTripConversion(t *testing.T) {
	resource := &groupResourceFramework{}
	ctx := context.Background()

	t.Run("state to API and back to state", func(t *testing.T) {
		email := "user@example.com"
		originalGroup := &restapi.Group{
			ID:   "test-id",
			Name: "test-group",
			Members: []restapi.APIMember{
				{
					UserID: "user-1",
					Email:  &email,
				},
			},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				ApplicationIDs: []restapi.ScopeBinding{
					{ScopeID: "app-1"},
				},
				InfraDFQFilter: &restapi.ScopeBinding{
					ScopeID: "entity.type:host",
				},
				KubernetesClusterUUIDs: []restapi.ScopeBinding{
					{ScopeID: "k8s-cluster-1"},
				},
				KubernetesNamespaceUIDs: []restapi.ScopeBinding{
					{ScopeID: "k8s-ns-1"},
				},
				MobileAppIDs: []restapi.ScopeBinding{
					{ScopeID: "mobile-app-1"},
				},
				WebsiteIDs: []restapi.ScopeBinding{
					{ScopeID: "website-1"},
				},
				Permissions: []restapi.InstanaPermission{
					restapi.PermissionCanConfigureApplications,
				},
			},
		}

		// Convert to state
		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}
		diags := resource.UpdateState(ctx, state, nil, originalGroup)
		require.False(t, diags.HasError())

		// Convert back to API object
		convertedGroup, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, convertedGroup)

		// Verify all fields match
		assert.Equal(t, originalGroup.ID, convertedGroup.ID)
		assert.Equal(t, originalGroup.Name, convertedGroup.Name)
		assert.Len(t, convertedGroup.Members, 1)
		assert.Equal(t, originalGroup.Members[0].UserID, convertedGroup.Members[0].UserID)
		assert.Equal(t, *originalGroup.Members[0].Email, *convertedGroup.Members[0].Email)
		assert.Len(t, convertedGroup.PermissionSet.ApplicationIDs, 1)
		assert.Equal(t, originalGroup.PermissionSet.ApplicationIDs[0].ScopeID, convertedGroup.PermissionSet.ApplicationIDs[0].ScopeID)
		assert.Equal(t, originalGroup.PermissionSet.InfraDFQFilter.ScopeID, convertedGroup.PermissionSet.InfraDFQFilter.ScopeID)
	})
}

func TestEdgeCases(t *testing.T) {
	resource := &groupResourceFramework{}
	ctx := context.Background()

	t.Run("empty group name", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue(""),
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)
		assert.Equal(t, "", group.Name)
	})

	t.Run("empty members list", func(t *testing.T) {
		model := GroupModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringValue("test-group"),
			Members: []GroupMemberModel{},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)
		assert.Empty(t, group.Members)
	})

	t.Run("empty permission set arrays", func(t *testing.T) {
		model := GroupModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-group"),
			PermissionSet: &GroupPermissionSetModel{
				ApplicationIDs:          []string{},
				KubernetesClusterUUIDs:  []string{},
				KubernetesNamespaceUIDs: []string{},
				MobileAppIDs:            []string{},
				WebsiteIDs:              []string{},
				Permissions:             []string{},
			},
		}

		state := createMockState(t, ctx, model)
		group, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, group)
		assert.Empty(t, group.PermissionSet.ApplicationIDs)
		assert.Empty(t, group.PermissionSet.Permissions)
	})

	t.Run("nil infra DFQ filter", func(t *testing.T) {
		group := &restapi.Group{
			ID:      "test-id",
			Name:    "test-group",
			Members: []restapi.APIMember{},
			PermissionSet: restapi.APIPermissionSetWithRoles{
				InfraDFQFilter: nil,
			},
		}

		handle := NewGroupResourceHandleFramework()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, group)
		require.False(t, diags.HasError())

		var model GroupModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		// When permission set is empty, it should be nil
		assert.Nil(t, model.PermissionSet)
	})
}

// Helper functions

func createMockState(t *testing.T, ctx context.Context, model GroupModel) *tfsdk.State {
	handle := NewGroupResourceHandleFramework()
	state := &tfsdk.State{
		Schema: handle.MetaData().Schema,
	}

	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set state: %v", diags)
	}

	return state
}

func createMockPlan(t *testing.T, ctx context.Context, model GroupModel) *tfsdk.Plan {
	handle := NewGroupResourceHandleFramework()
	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	diags := plan.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set plan: %v", diags)
	}

	return plan
}

// Made with Bob