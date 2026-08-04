package team

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/tagfilter"
	tag "github.com/instana/instana-go-client/shared/tagfilter"
	common "github.com/instana/instana-go-client/shared/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stringsToTypesSet converts a []string to a types.Set for use in test struct literals.
func stringsToTypesSet(t *testing.T, ss ...string) types.Set {
	t.Helper()
	if len(ss) == 0 {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}
	elems := make([]attr.Value, len(ss))
	for i, s := range ss {
		elems[i] = types.StringValue(s)
	}
	result, diags := types.SetValue(types.StringType, elems)
	require.False(t, diags.HasError())
	return result
}

// extractStringsFromSet extracts a []string from a types.Set for assertions.
func extractStringsFromSet(t *testing.T, ctx context.Context, s types.Set) []string {
	t.Helper()
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	var result []string
	diags := s.ElementsAs(ctx, &result, false)
	require.False(t, diags.HasError())
	return result
}

// emptyTeamScopeModel returns a TeamScopeModel with all types.Set fields initialised to
// empty, properly-typed sets. Tests that only need to set a subset of scope fields should
// start from this base to avoid the "MISSING TYPE" framework panic.
func emptyTeamScopeModel(t *testing.T) TeamScopeModel {
	t.Helper()
	return TeamScopeModel{
		AccessPermissions:    stringsToTypesSet(t),
		Applications:         stringsToTypesSet(t),
		KubernetesClusters:   stringsToTypesSet(t),
		KubernetesNamespaces: stringsToTypesSet(t),
		MobileApps:           stringsToTypesSet(t),
		Websites:             stringsToTypesSet(t),
		InfraDFQFilter:       types.StringNull(),
		ActionFilter:         types.StringNull(),
		LogFilter:            types.StringNull(),
		BusinessPerspectives: stringsToTypesSet(t),
		SloIDs:               stringsToTypesSet(t),
		SyntheticTests:       stringsToTypesSet(t),
		SyntheticCredentials: stringsToTypesSet(t),
		TagIDs:               stringsToTypesSet(t),
	}
}

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
		team := &api.Team{
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Info: &api.TeamInfo{
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
		team := &api.Team{
			ID:   "test-id",
			Tag:  "test-team",
			Info: &api.TeamInfo{},
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Members: []api.TeamMember{
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Members: []api.TeamMember{
				{
					UserID: "user-1",
					Roles: []api.TeamRole{
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Members: []api.TeamMember{
				{
					UserID: "user-1",
					Roles:  []api.TeamRole{},
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		perms := extractStringsFromSet(t, ctx, model.Scope.AccessPermissions)
		assert.Len(t, perms, 2)
		assert.Contains(t, perms, "perm-1")
		assert.Contains(t, perms, "perm-2")
	})

	t.Run("team with scope - applications", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		apps := extractStringsFromSet(t, ctx, model.Scope.Applications)
		assert.Len(t, apps, 2)
		assert.Contains(t, apps, "app-1")
	})

	t.Run("team with scope - kubernetes clusters", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		clusters := extractStringsFromSet(t, ctx, model.Scope.KubernetesClusters)
		assert.Len(t, clusters, 2)
		assert.Contains(t, clusters, "k8s-1")
	})

	t.Run("team with scope - kubernetes namespaces", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		namespaces := extractStringsFromSet(t, ctx, model.Scope.KubernetesNamespaces)
		assert.Len(t, namespaces, 2)
		assert.Contains(t, namespaces, "ns-1")
	})

	t.Run("team with scope - mobile apps", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		mobileApps := extractStringsFromSet(t, ctx, model.Scope.MobileApps)
		assert.Len(t, mobileApps, 2)
		assert.Contains(t, mobileApps, "mobile-1")
	})

	t.Run("team with scope - websites", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		websites := extractStringsFromSet(t, ctx, model.Scope.Websites)
		assert.Len(t, websites, 2)
		assert.Contains(t, websites, "website-1")
	})

	t.Run("team with scope - infra DFQ filter", func(t *testing.T) {
		filter := "entity.type:host"
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		bp := extractStringsFromSet(t, ctx, model.Scope.BusinessPerspectives)
		assert.Len(t, bp, 2)
		assert.Contains(t, bp, "bp-1")
	})

	t.Run("team with scope - SLO IDs", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		sloIDs := extractStringsFromSet(t, ctx, model.Scope.SloIDs)
		assert.Len(t, sloIDs, 2)
		assert.Contains(t, sloIDs, "slo-1")
	})

	t.Run("team with scope - synthetic tests", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		tests := extractStringsFromSet(t, ctx, model.Scope.SyntheticTests)
		assert.Len(t, tests, 2)
		assert.Contains(t, tests, "test-1")
	})

	t.Run("team with scope - synthetic credentials", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		creds := extractStringsFromSet(t, ctx, model.Scope.SyntheticCredentials)
		assert.Len(t, creds, 2)
		assert.Contains(t, creds, "cred-1")
	})

	t.Run("team with scope - tag IDs", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
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
		tagIDs := extractStringsFromSet(t, ctx, model.Scope.TagIDs)
		assert.Len(t, tagIDs, 2)
		assert.Contains(t, tagIDs, "tag-1")
	})

	t.Run("team with scope - restricted application filter with label", func(t *testing.T) {
		label := "test-label"
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
				RestrictedApplicationFilter: &api.RestrictedApplicationFilter{
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
		scope := api.RestrictedApplicationFilterScopeIncludeNoDownstream
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
				RestrictedApplicationFilter: &api.RestrictedApplicationFilter{
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
		assert.Equal(t, string(api.RestrictedApplicationFilterScopeIncludeNoDownstream), model.Scope.RestrictedApplicationFilter.Scope.ValueString())
	})

	t.Run("team with scope - restricted application filter with tag filter", func(t *testing.T) {
		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Scope: &api.TeamScope{
				RestrictedApplicationFilter: &api.RestrictedApplicationFilter{
					TagFilterExpression: tag.NewStringTagFilter(
						tagfilter.TagFilterEntityNotApplicable,
						"entity.type",
						common.EqualsOperator,
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
		scope := api.RestrictedApplicationFilterScopeIncludeAllDownstream

		team := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Info: &api.TeamInfo{
				Description: &desc,
			},
			Members: []api.TeamMember{
				{
					UserID: "user-1",
					Roles: []api.TeamRole{
						{RoleID: "role-1"},
					},
				},
			},
			Scope: &api.TeamScope{
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
				RestrictedApplicationFilter: &api.RestrictedApplicationFilter{
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
				AccessPermissions:    stringsToTypesSet(t, "perm-1"),
				Applications:         stringsToTypesSet(t, "app-1"),
				KubernetesClusters:   stringsToTypesSet(t, "k8s-1"),
				KubernetesNamespaces: stringsToTypesSet(t, "ns-1"),
				MobileApps:           stringsToTypesSet(t, "mobile-1"),
				Websites:             stringsToTypesSet(t, "website-1"),
				InfraDFQFilter:       types.StringValue("entity.type:host"),
				BusinessPerspectives: stringsToTypesSet(t, "bp-1"),
				SloIDs:               stringsToTypesSet(t, "slo-1"),
				SyntheticTests:       stringsToTypesSet(t, "test-1"),
				SyntheticCredentials: stringsToTypesSet(t, "cred-1"),
				TagIDs:               stringsToTypesSet(t, "tag-1"),
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label: types.StringValue("test-label"),
					Scope: types.StringValue(string(api.RestrictedApplicationFilterScopeIncludeAllDownstream)),
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
		assert.Len(t, model.Scope.AccessPermissions.Elements(), 1)
		assert.Len(t, model.Scope.Applications.Elements(), 1)
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
			Scope: func() *TeamScopeModel {
					s := emptyTeamScopeModel(t)
					s.AccessPermissions = stringsToTypesSet(t, "perm-1", "perm-2")
					return &s
				}(),
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		assert.Len(t, team.Scope.AccessPermissions, 2)
		assert.Contains(t, team.Scope.AccessPermissions, "perm-1")
	})

	t.Run("team with scope - applications", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: func() *TeamScopeModel {
					s := emptyTeamScopeModel(t)
					s.Applications = stringsToTypesSet(t, "app-1", "app-2")
					return &s
				}(),
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		assert.Len(t, team.Scope.Applications, 2)
		assert.Contains(t, team.Scope.Applications, "app-1")
	})

	t.Run("team with scope - infra DFQ filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: func() *TeamScopeModel {
					s := emptyTeamScopeModel(t)
					s.InfraDFQFilter = types.StringValue("entity.type:host")
					return &s
				}(),
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
			Scope: func() *TeamScopeModel {
					s := emptyTeamScopeModel(t)
					s.InfraDFQFilter = types.StringNull()
					return &s
				}(),
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
			Scope: func() *TeamScopeModel {
					s := emptyTeamScopeModel(t)
					s.ActionFilter = types.StringValue("action.type:custom")
					return &s
				}(),
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
			Scope: func() *TeamScopeModel {
					s := emptyTeamScopeModel(t)
					s.LogFilter = types.StringValue("log.level:error")
					return &s
				}(),
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
		scope := emptyTeamScopeModel(t)
		scope.KubernetesClusters = stringsToTypesSet(t, "k8s-1")
		scope.KubernetesNamespaces = stringsToTypesSet(t, "ns-1")
		scope.MobileApps = stringsToTypesSet(t, "mobile-1")
		scope.Websites = stringsToTypesSet(t, "website-1")
		scope.BusinessPerspectives = stringsToTypesSet(t, "bp-1")
		scope.SloIDs = stringsToTypesSet(t, "slo-1")
		scope.SyntheticTests = stringsToTypesSet(t, "test-1")
		scope.SyntheticCredentials = stringsToTypesSet(t, "cred-1")
		scope.TagIDs = stringsToTypesSet(t, "tag-1")
		model := TeamModel{
			ID:    types.StringValue("test-id"),
			Tag:   types.StringValue("test-team"),
			Scope: &scope,
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
			Scope: func() *TeamScopeModel {
				s := emptyTeamScopeModel(t)
				s.RestrictedApplicationFilter = &TeamRestrictedApplicationFilterModel{
					Label: types.StringValue("test-label"),
				}
				return &s
			}(),
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
			Scope: func() *TeamScopeModel {
				s := emptyTeamScopeModel(t)
				s.RestrictedApplicationFilter = &TeamRestrictedApplicationFilterModel{
					Scope: types.StringValue(string(api.RestrictedApplicationFilterScopeIncludeNoDownstream)),
				}
				return &s
			}(),
		}

		state := createMockTeamState(t, ctx, model)
		team, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, team)

		require.NotNil(t, team.Scope)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter)
		require.NotNil(t, team.Scope.RestrictedApplicationFilter.Scope)
		assert.Equal(t, api.RestrictedApplicationFilterScopeIncludeNoDownstream, *team.Scope.RestrictedApplicationFilter.Scope)
	})

	t.Run("team with scope - restricted application filter with tag filter", func(t *testing.T) {
		model := TeamModel{
			ID:  types.StringValue("test-id"),
			Tag: types.StringValue("test-team"),
			Scope: func() *TeamScopeModel {
				s := emptyTeamScopeModel(t)
				s.RestrictedApplicationFilter = &TeamRestrictedApplicationFilterModel{
					TagFilterExpression: types.StringValue("entity.type EQUALS 'service'"),
				}
				return &s
			}(),
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
			Scope: func() *TeamScopeModel {
				s := emptyTeamScopeModel(t)
				s.RestrictedApplicationFilter = &TeamRestrictedApplicationFilterModel{
					Label:               types.StringNull(),
					Scope:               types.StringNull(),
					TagFilterExpression: types.StringNull(),
				}
				return &s
			}(),
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
				AccessPermissions:    stringsToTypesSet(t, "perm-1"),
				Applications:         stringsToTypesSet(t, "app-1"),
				KubernetesClusters:   stringsToTypesSet(t, "k8s-1"),
				KubernetesNamespaces: stringsToTypesSet(t, "ns-1"),
				MobileApps:           stringsToTypesSet(t, "mobile-1"),
				Websites:             stringsToTypesSet(t, "website-1"),
				InfraDFQFilter:       types.StringValue("entity.type:host"),
				ActionFilter:         types.StringValue("action.type:custom"),
				LogFilter:            types.StringValue("log.level:error"),
				BusinessPerspectives: stringsToTypesSet(t, "bp-1"),
				SloIDs:               stringsToTypesSet(t, "slo-1"),
				SyntheticTests:       stringsToTypesSet(t, "test-1"),
				SyntheticCredentials: stringsToTypesSet(t, "cred-1"),
				TagIDs:               stringsToTypesSet(t, "tag-1"),
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label:               types.StringValue("test-label"),
					Scope:               types.StringValue(string(api.RestrictedApplicationFilterScopeIncludeAllDownstream)),
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
		scope := api.RestrictedApplicationFilterScopeIncludeAllDownstream

		originalTeam := &api.Team{
			ID:  "test-id",
			Tag: "test-team",
			Info: &api.TeamInfo{
				Description: &desc,
			},
			Members: []api.TeamMember{
				{
					UserID: "user-1",
					Roles: []api.TeamRole{
						{RoleID: "role-1"},
					},
				},
			},
			Scope: &api.TeamScope{
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
				RestrictedApplicationFilter: &api.RestrictedApplicationFilter{
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
				AccessPermissions:    stringsToTypesSet(t, "perm-1"),
				Applications:         stringsToTypesSet(t, "app-1"),
				KubernetesClusters:   stringsToTypesSet(t, "k8s-1"),
				KubernetesNamespaces: stringsToTypesSet(t, "ns-1"),
				MobileApps:           stringsToTypesSet(t, "mobile-1"),
				Websites:             stringsToTypesSet(t, "website-1"),
				InfraDFQFilter:       types.StringValue("entity.type:host"),
				BusinessPerspectives: stringsToTypesSet(t, "bp-1"),
				SloIDs:               stringsToTypesSet(t, "slo-1"),
				SyntheticTests:       stringsToTypesSet(t, "test-1"),
				SyntheticCredentials: stringsToTypesSet(t, "cred-1"),
				TagIDs:               stringsToTypesSet(t, "tag-1"),
				RestrictedApplicationFilter: &TeamRestrictedApplicationFilterModel{
					Label: types.StringValue("test-label"),
					Scope: types.StringValue(string(api.RestrictedApplicationFilterScopeIncludeAllDownstream)),
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
				AccessPermissions:    stringsToTypesSet(t),
				Applications:         stringsToTypesSet(t),
				KubernetesClusters:   stringsToTypesSet(t),
				KubernetesNamespaces: stringsToTypesSet(t),
				MobileApps:           stringsToTypesSet(t),
				Websites:             stringsToTypesSet(t),
				BusinessPerspectives: stringsToTypesSet(t),
				SloIDs:               stringsToTypesSet(t),
				SyntheticTests:       stringsToTypesSet(t),
				SyntheticCredentials: stringsToTypesSet(t),
				TagIDs:               stringsToTypesSet(t),
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
		team := &api.Team{
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
			Scope: func() *TeamScopeModel {
				s := emptyTeamScopeModel(t)
				s.RestrictedApplicationFilter = &TeamRestrictedApplicationFilterModel{
					TagFilterExpression: types.StringValue("invalid expression"),
				}
				return &s
			}(),
		}

		state := createMockTeamState(t, ctx, model)
		_, diags := resource.MapStateToDataObject(ctx, nil, state)
		assert.True(t, diags.HasError())
	})
}

func TestHelperFunctions(t *testing.T) {
	resource := &teamResource{}

	t.Run("mapTeamInfoToModel with nil description", func(t *testing.T) {
		apiInfo := &api.TeamInfo{
			Description: nil,
		}

		modelInfo := resource.mapTeamInfoToModel(apiInfo)
		require.NotNil(t, modelInfo)
		assert.True(t, modelInfo.Description.IsNull())
	})

	t.Run("mapTeamInfoToModel with description", func(t *testing.T) {
		desc := "Test description"
		apiInfo := &api.TeamInfo{
			Description: &desc,
		}

		modelInfo := resource.mapTeamInfoToModel(apiInfo)
		require.NotNil(t, modelInfo)
		assert.Equal(t, "Test description", modelInfo.Description.ValueString())
	})

	t.Run("mapMembersToModel with empty members", func(t *testing.T) {
		apiMembers := []api.TeamMember{}

		modelMembers := resource.mapMembersToModel(apiMembers)
		assert.Empty(t, modelMembers)
	})

	t.Run("mapRolesToModel with nil roles", func(t *testing.T) {
		roles := resource.mapRolesToModel(nil)
		assert.Nil(t, roles)
	})

	t.Run("mapRolesToModel with empty roles", func(t *testing.T) {
		apiRoles := []api.TeamRole{}
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
