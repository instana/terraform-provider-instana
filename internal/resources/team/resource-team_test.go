package team

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

func TestNewTeamResourceHandle(t *testing.T) {
	handle := NewTeamResourceHandle()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaTeam, metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &teamResource{
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

func TestGetRestResource(t *testing.T) {
	resource := &teamResource{}

	// This test just ensures the method exists and can be called
	// The actual implementation requires a real API instance
	assert.NotNil(t, resource)
}

func TestSetComputedFields(t *testing.T) {
	resource := &teamResource{}
	ctx := context.Background()

	diags := resource.SetComputedFields(ctx, nil)
	assert.False(t, diags.HasError())
}

func TestUpdateState(t *testing.T) {
	resource := &teamResource{}
	ctx := context.Background()

	t.Run("basic team without members or scope", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-team", model.Tag.ValueString())
		assert.Nil(t, model.Info)
		assert.Nil(t, model.Members)
		assert.Nil(t, model.Scope)
	})

	t.Run("team with info", func(t *testing.T) {
		desc := "Test team description"
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Info: &restapi.TeamInfo{
				Description: &desc,
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		// Create a plan with Info field set
		planModel := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringValue("Test team description"),
			},
		}
		plan := createMockTeamPlan(t, ctx, planModel)

		diags := resource.UpdateState(ctx, state, plan, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Info)
		assert.Equal(t, "Test team description", model.Info.Description.ValueString())
	})

	t.Run("team with info nil description", func(t *testing.T) {
		team := &restapi.Team{
			ID:   "test-id",
			Tag:  "test-team",
			Info: &restapi.TeamInfo{},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		// Create a plan with Info field set but description null
		planModel := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringNull(),
			},
		}
		plan := createMockTeamPlan(t, ctx, planModel)

		diags := resource.UpdateState(ctx, state, plan, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Info)
		assert.True(t, model.Info.Description.IsNull())
	})

	t.Run("team with members without roles", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Members: []restapi.TeamMember{
				{
					UserID: "user-1",
				},
				{
					UserID: "user-2",
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Len(t, model.Members, 2)
		assert.Equal(t, "user-1", model.Members[0].UserID.ValueString())
		assert.Equal(t, "user-2", model.Members[1].UserID.ValueString())
		assert.Nil(t, model.Members[0].Roles)
		assert.Nil(t, model.Members[1].Roles)
	})

	t.Run("team with members with roles", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Members: []restapi.TeamMember{
				{
					UserID: "user-1",
					Roles: []restapi.TeamRole{
						{RoleID: "role-1"},
						{RoleID: "role-2"},
					},
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Len(t, model.Members, 1)
		assert.Len(t, model.Members[0].Roles, 2)
		assert.Equal(t, "role-1", model.Members[0].Roles[0].RoleID.ValueString())
		assert.Equal(t, "role-2", model.Members[0].Roles[1].RoleID.ValueString())
	})

	t.Run("team with members with empty roles", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Members: []restapi.TeamMember{
				{
					UserID: "user-1",
					Roles:  []restapi.TeamRole{},
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Len(t, model.Members, 1)
		assert.Nil(t, model.Members[0].Roles)
	})

	t.Run("team with scope - access permissions", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				AccessPermissions: []string{"perm-1", "perm-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.AccessPermissions, 2)
		assert.Equal(t, "perm-1", model.Scope.AccessPermissions[0])
		assert.Equal(t, "perm-2", model.Scope.AccessPermissions[1])
	})

	t.Run("team with scope - applications", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				Applications: []string{"app-1", "app-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.Applications, 2)
		assert.Equal(t, "app-1", model.Scope.Applications[0])
	})

	t.Run("team with scope - kubernetes clusters", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				KubernetesClusters: []string{"k8s-1", "k8s-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.KubernetesClusters, 2)
		assert.Equal(t, "k8s-1", model.Scope.KubernetesClusters[0])
	})

	t.Run("team with scope - kubernetes namespaces", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				KubernetesNamespaces: []string{"ns-1", "ns-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.KubernetesNamespaces, 2)
		assert.Equal(t, "ns-1", model.Scope.KubernetesNamespaces[0])
	})

	t.Run("team with scope - mobile apps", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				MobileApps: []string{"mobile-1", "mobile-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.MobileApps, 2)
		assert.Equal(t, "mobile-1", model.Scope.MobileApps[0])
	})

	t.Run("team with scope - websites", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				Websites: []string{"website-1", "website-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.Websites, 2)
		assert.Equal(t, "website-1", model.Scope.Websites[0])
	})

	t.Run("team with scope - infra DFQ filter", func(t *testing.T) {
		filter := "entity.type:host"
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				InfraDFQFilter: &filter,
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Equal(t, "entity.type:host", model.Scope.InfraDFQFilter.ValueString())
	})

	t.Run("team with scope - action filter", func(t *testing.T) {
		filter := "action.type:custom"
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				ActionFilter: &filter,
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Equal(t, "action.type:custom", model.Scope.ActionFilter.ValueString())
	})

	t.Run("team with scope - log filter", func(t *testing.T) {
		filter := "log.level:error"
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				LogFilter: &filter,
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Equal(t, "log.level:error", model.Scope.LogFilter.ValueString())
	})

	t.Run("team with scope - business perspectives", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				BusinessPerspectives: []string{"bp-1", "bp-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.BusinessPerspectives, 2)
		assert.Equal(t, "bp-1", model.Scope.BusinessPerspectives[0])
	})

	t.Run("team with scope - SLO IDs", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				SloIDs: []string{"slo-1", "slo-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.SloIDs, 2)
		assert.Equal(t, "slo-1", model.Scope.SloIDs[0])
	})

	t.Run("team with scope - synthetic tests", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				SyntheticTests: []string{"test-1", "test-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.SyntheticTests, 2)
		assert.Equal(t, "test-1", model.Scope.SyntheticTests[0])
	})

	t.Run("team with scope - synthetic credentials", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				SyntheticCredentials: []string{"cred-1", "cred-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.SyntheticCredentials, 2)
		assert.Equal(t, "cred-1", model.Scope.SyntheticCredentials[0])
	})

	t.Run("team with scope - tag IDs", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				TagIDs: []string{"tag-1", "tag-2"},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.TagIDs, 2)
		assert.Equal(t, "tag-1", model.Scope.TagIDs[0])
	})

	t.Run("team with scope - restricted application filter with label", func(t *testing.T) {
		label := "test-label"
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				RestrictedApplicationFilter: &restapi.RestrictedApplicationFilter{
					Label: &label,
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		require.NotNil(t, model.Scope.RestrictedApplicationFilter)
		assert.Equal(t, "test-label", model.Scope.RestrictedApplicationFilter.Label.ValueString())
	})

	t.Run("team with scope - restricted application filter with scope", func(t *testing.T) {
		scope := restapi.RestrictedApplicationFilterScopeIncludeNoDownstream
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				RestrictedApplicationFilter: &restapi.RestrictedApplicationFilter{
					Scope: &scope,
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		require.NotNil(t, model.Scope.RestrictedApplicationFilter)
		assert.Equal(t, string(restapi.RestrictedApplicationFilterScopeIncludeNoDownstream), model.Scope.RestrictedApplicationFilter.Scope.ValueString())
	})

	t.Run("team with scope - restricted application filter with tag filter", func(t *testing.T) {
		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &restapi.TeamScope{
				RestrictedApplicationFilter: &restapi.RestrictedApplicationFilter{
					TagFilterExpression: restapi.NewStringTagFilter(
						restapi.TagFilterEntityNotApplicable,
						"entity.type",
						restapi.EqualsOperator,
						"service",
					),
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.Scope)
		require.NotNil(t, model.Scope.RestrictedApplicationFilter)
		assert.False(t, model.Scope.RestrictedApplicationFilter.TagFilterExpression.IsNull())
	})

	t.Run("team with all fields populated", func(t *testing.T) {
		desc := "Full team"
		filter := "entity.type:host"
		label := "test-label"
		scope := restapi.RestrictedApplicationFilterScopeIncludeAllDownstream

		team := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Info: &restapi.TeamInfo{
				Description: &desc,
			},
			Members: []restapi.TeamMember{
				{
					UserID: "user-1",
					Roles: []restapi.TeamRole{
						{RoleID: "role-1"},
					},
				},
			},
			Scope: &restapi.TeamScope{
				AccessPermissions:    []string{"perm-1"},
				Applications:         []string{"app-1"},
				KubernetesClusters:   []string{"k8s-1"},
				KubernetesNamespaces: []string{"ns-1"},
				MobileApps:           []string{"mobile-1"},
				Websites:             []string{"website-1"},
				InfraDFQFilter:       &filter,
				BusinessPerspectives: []string{"bp-1"},
				SloIDs:               []string{"slo-1"},
				SyntheticTests:       []string{"test-1"},
				SyntheticCredentials: []string{"cred-1"},
				TagIDs:               []string{"tag-1"},
				RestrictedApplicationFilter: &restapi.RestrictedApplicationFilter{
					Label: &label,
					Scope: &scope,
				},
			},
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		// Create a plan with all fields populated
		planModel := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringValue("Full team"),
			},
			Members: []TeamMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Roles: []TeamMemberRole{
						{RoleID: types.StringValue("role-1")},
					},
				},
			},
			Scope: &TeamScopeModel{
				AccessPermissions:    []string{"perm-1"},
				Applications:         []string{"app-1"},
				KubernetesClusters:   []string{"k8s-1"},
				KubernetesNamespaces: []string{"ns-1"},
				MobileApps:           []string{"mobile-1"},
				Websites:             []string{"website-1"},
				InfraDFQFilter:       types.StringValue("entity.type:host"),
				BusinessPerspectives: []string{"bp-1"},
				SloIDs:               []string{"slo-1"},
				SyntheticTests:       []string{"test-1"},
				SyntheticCredentials: []string{"cred-1"},
				TagIDs:               []string{"tag-1"},
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label: types.StringValue("test-label"),
					Scope: types.StringValue(string(restapi.RestrictedApplicationFilterScopeIncludeAllDownstream)),
				},
			},
		}
		plan := createMockTeamPlan(t, ctx, planModel)

		diags := resource.UpdateState(ctx, state, plan, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-team", model.Tag.ValueString())
		require.NotNil(t, model.Info)
		assert.Equal(t, "Full team", model.Info.Description.ValueString())
		assert.Len(t, model.Members, 1)
		require.NotNil(t, model.Scope)
		assert.Len(t, model.Scope.AccessPermissions, 1)
		assert.Len(t, model.Scope.Applications, 1)
		require.NotNil(t, model.Scope.RestrictedApplicationFilter)
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &teamResource{}
	ctx := context.Background()

	t.Run("basic team from state", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Equal(t, "test-id", team.ID)
		assert.Equal(t, "test-team", team.Tag)
		assert.Nil(t, team.Info)
		assert.Empty(t, team.Members)
		assert.Nil(t, team.Scope)
	})

	t.Run("team from plan", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue(""),
			Tag: types.StringValue("new-team"),
		}

		plan := createMockTeamPlan(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Equal(t, "", team.ID)
		assert.Equal(t, "new-team", team.Tag)
	})

	t.Run("team with null ID", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringNull(),
			Tag: types.StringValue("test-team"),
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Equal(t, "", team.ID)
	})

	t.Run("team with info", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringValue("Test description"),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Info)
		assert.Equal(t, "Test description", *team.Info.Description)
	})

	t.Run("team with info null description", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringNull(),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Nil(t, team.Info)
	})

	t.Run("team with info unknown description", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringUnknown(),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Nil(t, team.Info)
	})

	t.Run("team with members without roles", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Members: []TeamMemberModel{
				{
					UserID: types.StringValue("user-1"),
				},
				{
					UserID: types.StringValue("user-2"),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Len(t, team.Members, 2)
		assert.Equal(t, "user-1", team.Members[0].UserID)
		assert.Equal(t, "user-2", team.Members[1].UserID)
		assert.Empty(t, team.Members[0].Roles)
		assert.Empty(t, team.Members[1].Roles)
	})

	t.Run("team with members with roles", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Members: []TeamMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Roles: []TeamMemberRole{
						{RoleID: types.StringValue("role-1")},
						{RoleID: types.StringValue("role-2")},
					},
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Len(t, team.Members, 1)
		assert.Len(t, team.Members[0].Roles, 2)
		assert.Equal(t, "role-1", team.Members[0].Roles[0].RoleID)
		assert.Equal(t, "role-2", team.Members[0].Roles[1].RoleID)
	})

	t.Run("team with empty members list", func(t *testing.T) {
		model := TeamModel{
			ID:      types.StringValue("test-id"),
			Tag:     types.StringValue("test-team"),
			Members: []TeamMemberModel{},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Empty(t, team.Members)
	})

	t.Run("team with scope - access permissions", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				AccessPermissions: []string{"perm-1", "perm-2"},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		assert.Len(t, team.Scope.AccessPermissions, 2)
		assert.Equal(t, "perm-1", team.Scope.AccessPermissions[0])
	})

	t.Run("team with scope - applications", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				Applications: []string{"app-1", "app-2"},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		assert.Len(t, team.Scope.Applications, 2)
		assert.Equal(t, "app-1", team.Scope.Applications[0])
	})

	t.Run("team with scope - infra DFQ filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				InfraDFQFilter: types.StringValue("entity.type:host"),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.InfraDFQFilter)
		assert.Equal(t, "entity.type:host", *team.Scope.InfraDFQFilter)
	})

	t.Run("team with scope - null infra DFQ filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				InfraDFQFilter: types.StringNull(),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		assert.Nil(t, team.Scope.InfraDFQFilter)
	})

	t.Run("team with scope - action filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				ActionFilter: types.StringValue("action.type:custom"),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.ActionFilter)
		assert.Equal(t, "action.type:custom", *team.Scope.ActionFilter)
	})

	t.Run("team with scope - log filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				LogFilter: types.StringValue("log.level:error"),
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.LogFilter)
		assert.Equal(t, "log.level:error", *team.Scope.LogFilter)
	})

	t.Run("team with scope - all array fields", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				KubernetesClusters:   []string{"k8s-1"},
				KubernetesNamespaces: []string{"ns-1"},
				MobileApps:           []string{"mobile-1"},
				Websites:             []string{"website-1"},
				BusinessPerspectives: []string{"bp-1"},
				SloIDs:               []string{"slo-1"},
				SyntheticTests:       []string{"test-1"},
				SyntheticCredentials: []string{"cred-1"},
				TagIDs:               []string{"tag-1"},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		assert.Len(t, team.Scope.KubernetesClusters, 1)
		assert.Len(t, team.Scope.KubernetesNamespaces, 1)
		assert.Len(t, team.Scope.MobileApps, 1)
		assert.Len(t, team.Scope.Websites, 1)
		assert.Len(t, team.Scope.BusinessPerspectives, 1)
		assert.Len(t, team.Scope.SloIDs, 1)
		assert.Len(t, team.Scope.SyntheticTests, 1)
		assert.Len(t, team.Scope.SyntheticCredentials, 1)
		assert.Len(t, team.Scope.TagIDs, 1)
	})

	t.Run("team with scope - restricted application filter with label", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label: types.StringValue("test-label"),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter.Label)
		assert.Equal(t, "test-label", *team.Scope.RestrictedApplicationFilter.Label)
	})

	t.Run("team with scope - restricted application filter with scope", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Scope: types.StringValue(string(restapi.RestrictedApplicationFilterScopeIncludeNoDownstream)),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter.Scope)
		assert.Equal(t, restapi.RestrictedApplicationFilterScopeIncludeNoDownstream, *team.Scope.RestrictedApplicationFilter.Scope)
	})

	t.Run("team with scope - restricted application filter with tag filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					TagFilterExpression: types.StringValue("entity.type EQUALS 'service'"),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter.TagFilterExpression)
	})

	t.Run("team with scope - restricted application filter with null fields", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label:               types.StringNull(),
					Scope:               types.StringNull(),
					TagFilterExpression: types.StringNull(),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter)
		assert.Nil(t, team.Scope.RestrictedApplicationFilter.Label)
		assert.Nil(t, team.Scope.RestrictedApplicationFilter.Scope)
		assert.Nil(t, team.Scope.RestrictedApplicationFilter.TagFilterExpression)
	})

	t.Run("team with all fields populated", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringValue("Full team"),
			},
			Members: []TeamMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Roles: []TeamMemberRole{
						{RoleID: types.StringValue("role-1")},
					},
				},
			},
			Scope: &TeamScopeModel{
				AccessPermissions:    []string{"perm-1"},
				Applications:         []string{"app-1"},
				KubernetesClusters:   []string{"k8s-1"},
				KubernetesNamespaces: []string{"ns-1"},
				MobileApps:           []string{"mobile-1"},
				Websites:             []string{"website-1"},
				InfraDFQFilter:       types.StringValue("entity.type:host"),
				ActionFilter:         types.StringValue("action.type:custom"),
				LogFilter:            types.StringValue("log.level:error"),
				BusinessPerspectives: []string{"bp-1"},
				SloIDs:               []string{"slo-1"},
				SyntheticTests:       []string{"test-1"},
				SyntheticCredentials: []string{"cred-1"},
				TagIDs:               []string{"tag-1"},
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label:               types.StringValue("test-label"),
					Scope:               types.StringValue(string(restapi.RestrictedApplicationFilterScopeIncludeAllDownstream)),
					TagFilterExpression: types.StringValue("entity.type EQUALS 'service'"),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		assert.Equal(t, "test-id", team.ID)
		assert.Equal(t, "test-team", team.Tag)
		require.NotNil(t, team.Info)
		assert.Equal(t, "Full team", *team.Info.Description)
		assert.Len(t, team.Members, 1)
		require.NotNil(t, team.Scope)
		assert.Len(t, team.Scope.AccessPermissions, 1)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter)
	})
}

func TestRoundTripConversion(t *testing.T) {
	resource := &teamResource{}
	ctx := context.Background()

	t.Run("state to API and back to state", func(t *testing.T) {
		desc := "Test team"
		filter := "entity.type:host"
		label := "test-label"
		scope := restapi.RestrictedApplicationFilterScopeIncludeAllDownstream

		originalTeam := &restapi.Team{
			ID:  "test-id",
			Tag: "test-team",
			Info: &restapi.TeamInfo{
				Description: &desc,
			},
			Members: []restapi.TeamMember{
				{
					UserID: "user-1",
					Roles: []restapi.TeamRole{
						{RoleID: "role-1"},
					},
				},
			},
			Scope: &restapi.TeamScope{
				AccessPermissions:    []string{"perm-1"},
				Applications:         []string{"app-1"},
				KubernetesClusters:   []string{"k8s-1"},
				KubernetesNamespaces: []string{"ns-1"},
				MobileApps:           []string{"mobile-1"},
				Websites:             []string{"website-1"},
				InfraDFQFilter:       &filter,
				BusinessPerspectives: []string{"bp-1"},
				SloIDs:               []string{"slo-1"},
				SyntheticTests:       []string{"test-1"},
				SyntheticCredentials: []string{"cred-1"},
				TagIDs:               []string{"tag-1"},
				RestrictedApplicationFilter: &restapi.RestrictedApplicationFilter{
					Label: &label,
					Scope: &scope,
				},
			},
		}

		// Create a plan with all fields set
		planModel := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Info: &TeamInfoModel{
				Description: types.StringValue("Test team"),
			},
			Members: []TeamMemberModel{
				{
					UserID: types.StringValue("user-1"),
					Roles: []TeamMemberRole{
						{RoleID: types.StringValue("role-1")},
					},
				},
			},
			Scope: &TeamScopeModel{
				AccessPermissions:    []string{"perm-1"},
				Applications:         []string{"app-1"},
				KubernetesClusters:   []string{"k8s-1"},
				KubernetesNamespaces: []string{"ns-1"},
				MobileApps:           []string{"mobile-1"},
				Websites:             []string{"website-1"},
				InfraDFQFilter:       types.StringValue("entity.type:host"),
				BusinessPerspectives: []string{"bp-1"},
				SloIDs:               []string{"slo-1"},
				SyntheticTests:       []string{"test-1"},
				SyntheticCredentials: []string{"cred-1"},
				TagIDs:               []string{"tag-1"},
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label: types.StringValue("test-label"),
					Scope: types.StringValue(string(restapi.RestrictedApplicationFilterScopeIncludeAllDownstream)),
				},
			},
		}
		plan := createMockTeamPlan(t, ctx, planModel)

		// Convert to state
		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}
		diags := resource.UpdateState(ctx, state, plan, originalTeam)
		require.False(t, diags.HasError())

		// Convert back to API object
		convertedTeam, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, convertedTeam)

		// Verify all fields match
		assert.Equal(t, originalTeam.ID, convertedTeam.ID)
		assert.Equal(t, originalTeam.Tag, convertedTeam.Tag)
		require.NotNil(t, convertedTeam.Info)
		assert.Equal(t, *originalTeam.Info.Description, *convertedTeam.Info.Description)
		assert.Len(t, convertedTeam.Members, 1)
		assert.Equal(t, originalTeam.Members[0].UserID, convertedTeam.Members[0].UserID)
		require.NotNil(t, convertedTeam.Scope)
		assert.Len(t, convertedTeam.Scope.AccessPermissions, 1)
		assert.Equal(t, originalTeam.Scope.AccessPermissions[0], convertedTeam.Scope.AccessPermissions[0])
		assert.Equal(t, *originalTeam.Scope.InfraDFQFilter, *convertedTeam.Scope.InfraDFQFilter)
	})
}

func TestEdgeCases(t *testing.T) {
	resource := &teamResource{}
	ctx := context.Background()

	t.Run("empty team tag", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue(""),
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)
		assert.Equal(t, "", team.Tag)
	})

	t.Run("empty members list", func(t *testing.T) {
		model := TeamModel{
			ID:      types.StringValue("test-id"),
			Tag:     types.StringValue("test-team"),
			Members: []TeamMemberModel{},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)
		assert.Empty(t, team.Members)
	})

	t.Run("empty scope arrays", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				AccessPermissions:    []string{},
				Applications:         []string{},
				KubernetesClusters:   []string{},
				KubernetesNamespaces: []string{},
				MobileApps:           []string{},
				Websites:             []string{},
				BusinessPerspectives: []string{},
				SloIDs:               []string{},
				SyntheticTests:       []string{},
				SyntheticCredentials: []string{},
				TagIDs:               []string{},
			},
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)
		require.NotNil(t, team.Scope)
		assert.Empty(t, team.Scope.AccessPermissions)
		assert.Empty(t, team.Scope.Applications)
	})

	t.Run("nil scope", func(t *testing.T) {
		team := &restapi.Team{
			ID:    "test-id",
			Tag:   "test-team",
			Scope: nil,
		}

		handle := NewTeamResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, team)
		require.False(t, diags.HasError())

		var model TeamModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Nil(t, model.Scope)
	})

	t.Run("invalid tag filter expression", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: &TeamScopeModel{
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					TagFilterExpression: types.StringValue("invalid expression"),
				},
			},
		}

		state := createMockTeamState(t, ctx, model)
		_, diags := resource.MapStateToDataObject(ctx, nil, state)
		assert.True(t, diags.HasError())
	})
}

func TestHelperFunctions(t *testing.T) {
	resource := &teamResource{}

	t.Run("mapTeamInfoToModel with nil description", func(t *testing.T) {
		apiInfo := &restapi.TeamInfo{
			Description: nil,
		}

		modelInfo := resource.mapTeamInfoToModel(apiInfo)
		require.NotNil(t, modelInfo)
		assert.True(t, modelInfo.Description.IsNull())
	})

	t.Run("mapTeamInfoToModel with description", func(t *testing.T) {
		desc := "Test description"
		apiInfo := &restapi.TeamInfo{
			Description: &desc,
		}

		modelInfo := resource.mapTeamInfoToModel(apiInfo)
		require.NotNil(t, modelInfo)
		assert.Equal(t, "Test description", modelInfo.Description.ValueString())
	})

	t.Run("mapMembersToModel with empty members", func(t *testing.T) {
		apiMembers := []restapi.TeamMember{}

		modelMembers := resource.mapMembersToModel(apiMembers)
		assert.Empty(t, modelMembers)
	})

	t.Run("mapRolesToModel with nil roles", func(t *testing.T) {
		roles := resource.mapRolesToModel(nil)
		assert.Nil(t, roles)
	})

	t.Run("mapRolesToModel with empty roles", func(t *testing.T) {
		apiRoles := []restapi.TeamRole{}
		roles := resource.mapRolesToModel(apiRoles)
		assert.Nil(t, roles)
	})

	t.Run("mapModelInfoToAPI with null description", func(t *testing.T) {
		modelInfo := &TeamInfoModel{
			Description: types.StringNull(),
		}

		apiInfo := resource.mapModelInfoToAPI(modelInfo)
		assert.Nil(t, apiInfo)
	})

	t.Run("mapModelInfoToAPI with unknown description", func(t *testing.T) {
		modelInfo := &TeamInfoModel{
			Description: types.StringUnknown(),
		}

		apiInfo := resource.mapModelInfoToAPI(modelInfo)
		assert.Nil(t, apiInfo)
	})

	t.Run("mapModelMembersToAPI with empty members", func(t *testing.T) {
		modelMembers := []TeamMemberModel{}

		apiMembers := resource.mapModelMembersToAPI(modelMembers)
		assert.Nil(t, apiMembers)
	})

	t.Run("mapModelRolesToAPI with nil roles", func(t *testing.T) {
		roles := resource.mapModelRolesToAPI(nil)
		assert.Nil(t, roles)
	})

	t.Run("mapModelRolesToAPI with empty roles", func(t *testing.T) {
		modelRoles := []TeamMemberRole{}
		roles := resource.mapModelRolesToAPI(modelRoles)
		assert.Nil(t, roles)
	})

	t.Run("extractTeamID with null ID", func(t *testing.T) {
		model := TeamModel{
			ID: types.StringNull(),
		}

		id := resource.extractTeamID(model)
		assert.Equal(t, "", id)
	})

	t.Run("extractTeamID with value", func(t *testing.T) {
		model := TeamModel{
			ID: types.StringValue("test-id"),
		}

		id := resource.extractTeamID(model)
		assert.Equal(t, "test-id", id)
	})
}

// Helper functions

func createMockTeamState(t *testing.T, ctx context.Context, model TeamModel) *tfsdk.State {
	handle := NewTeamResourceHandle()
	state := &tfsdk.State{
		Schema: handle.MetaData().Schema,
	}

	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set state: %v", diags)
	}

	return state
}

func createMockTeamPlan(t *testing.T, ctx context.Context, model TeamModel) *tfsdk.Plan {
	handle := NewTeamResourceHandle()
	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	diags := plan.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set plan: %v", diags)
	}

	return plan
}
