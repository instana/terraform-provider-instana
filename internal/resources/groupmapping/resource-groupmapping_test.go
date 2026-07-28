package groupmapping

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string { return &s }

func TestNewGroupMappingResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewGroupMappingResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaGroupMapping, metadata.ResourceName)
		assert.Equal(t, int64(0), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewGroupMappingResourceHandle()
		s := handle.MetaData().Schema

		assert.NotNil(t, s.Attributes[GroupMappingFieldID])
		assert.NotNil(t, s.Attributes[GroupMappingFieldKey])
		assert.NotNil(t, s.Attributes[GroupMappingFieldValue])
		assert.NotNil(t, s.Attributes[GroupMappingFieldGroupID])
		assert.NotNil(t, s.Attributes[GroupMappingFieldTeamID])
	})
}

func TestGroupMappingMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		r := &groupMappingResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaGroupMapping,
				SchemaVersion: 0,
			},
		}
		metadata := r.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaGroupMapping, metadata.ResourceName)
	})
}

func TestGroupMappingGetRestResource(t *testing.T) {
	t.Run("should return group mappings rest resource", func(t *testing.T) {
		r := &groupMappingResource{}
		assert.NotNil(t, r.GetRestResource)
	})
}

// mockGroupMappingAPI extends the common mock to provide specific behavior for group mapping tests
type mockGroupMappingAPI struct {
	testutils.MockInstanaAPI
}

func (m *mockGroupMappingAPI) GroupMappings() rest.RestResource[*api.GroupMapping] {
	return &mockGroupMappingRestResource{}
}

// mockGroupMappingRestResource implements all required methods from RestResource interface
type mockGroupMappingRestResource struct{}

func (m *mockGroupMappingRestResource) GetAll() (*[]*api.GroupMapping, error)          { return nil, nil }
func (m *mockGroupMappingRestResource) GetOne(id string) (*api.GroupMapping, error)    { return nil, nil }
func (m *mockGroupMappingRestResource) Create(data *api.GroupMapping) (*api.GroupMapping, error) {
	return nil, nil
}
func (m *mockGroupMappingRestResource) Update(data *api.GroupMapping) (*api.GroupMapping, error) {
	return nil, nil
}
func (m *mockGroupMappingRestResource) Delete(data *api.GroupMapping) error   { return nil }
func (m *mockGroupMappingRestResource) DeleteByID(id string) error             { return nil }

func TestGroupMappingSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		r := &groupMappingResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaGroupMapping,
				Schema:        NewGroupMappingResourceHandle().MetaData().Schema,
				SchemaVersion: 0,
			},
		}
		ctx := context.Background()
		plan := &tfsdk.Plan{Schema: r.metaData.Schema}
		diags := r.SetComputedFields(ctx, plan)
		assert.False(t, diags.HasError())
	})
}

func TestGroupMappingMapStateToDataObject(t *testing.T) {
	r := &groupMappingResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaGroupMapping,
			Schema:        NewGroupMappingResourceHandle().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should map complete model with team_id from plan successfully", func(t *testing.T) {
		model := GroupMappingModel{
			ID:      types.StringValue("mapping-id-123"),
			Key:     types.StringValue("department"),
			Value:   types.StringValue("engineering"),
			GroupID: types.StringValue("group-id-456"),
			TeamID:  types.StringValue("team-id-789"),
		}

		plan := &tfsdk.Plan{Schema: r.metaData.Schema}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := r.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		require.NotNil(t, result)
		assert.Equal(t, "mapping-id-123", result.ID)
		assert.Equal(t, "department", result.Key)
		assert.Equal(t, "engineering", result.Value)
		assert.Equal(t, "group-id-456", result.GroupID)
		require.NotNil(t, result.TeamID)
		assert.Equal(t, "team-id-789", *result.TeamID)
	})

	t.Run("should map model without team_id from plan successfully", func(t *testing.T) {
		model := GroupMappingModel{
			ID:      types.StringValue("mapping-id-123"),
			Key:     types.StringValue("department"),
			Value:   types.StringValue("engineering"),
			GroupID: types.StringValue("group-id-456"),
			TeamID:  types.StringNull(),
		}

		plan := &tfsdk.Plan{Schema: r.metaData.Schema}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := r.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		require.NotNil(t, result)
		assert.Equal(t, "group-id-456", result.GroupID)
		assert.Nil(t, result.TeamID)
	})

	t.Run("should map model from state when plan is nil", func(t *testing.T) {
		model := GroupMappingModel{
			ID:      types.StringValue("mapping-id-789"),
			Key:     types.StringValue("team"),
			Value:   types.StringValue("platform"),
			GroupID: types.StringValue("group-id-000"),
			TeamID:  types.StringValue("team-id-111"),
		}

		state := &tfsdk.State{Schema: r.metaData.Schema}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := r.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		require.NotNil(t, result)
		assert.Equal(t, "mapping-id-789", result.ID)
		assert.Equal(t, "team", result.Key)
		assert.Equal(t, "platform", result.Value)
		assert.Equal(t, "group-id-000", result.GroupID)
		require.NotNil(t, result.TeamID)
		assert.Equal(t, "team-id-111", *result.TeamID)
	})

	t.Run("should use empty string for null ID", func(t *testing.T) {
		model := GroupMappingModel{
			ID:      types.StringNull(),
			Key:     types.StringValue("role"),
			Value:   types.StringValue("admin"),
			GroupID: types.StringValue("group-id-111"),
			TeamID:  types.StringNull(),
		}

		plan := &tfsdk.Plan{Schema: r.metaData.Schema}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := r.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		require.NotNil(t, result)
		assert.Equal(t, "", result.ID)
		assert.Nil(t, result.TeamID)
	})
}

func TestGroupMappingUpdateState(t *testing.T) {
	r := &groupMappingResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaGroupMapping,
			Schema:        NewGroupMappingResourceHandle().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should set state from API response with team_id", func(t *testing.T) {
		mapping := &api.GroupMapping{
			ID:      "mapping-id-123",
			Key:     "department",
			Value:   "engineering",
			GroupID: "group-id-456",
			TeamID:  strPtr("team-id-789"),
		}

		state := &tfsdk.State{Schema: r.metaData.Schema}
		diags := r.UpdateState(ctx, state, nil, mapping)

		assert.False(t, diags.HasError())

		var model GroupMappingModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.Equal(t, "mapping-id-123", model.ID.ValueString())
		assert.Equal(t, "department", model.Key.ValueString())
		assert.Equal(t, "engineering", model.Value.ValueString())
		assert.Equal(t, "group-id-456", model.GroupID.ValueString())
		assert.Equal(t, "team-id-789", model.TeamID.ValueString())
	})

	t.Run("should set state from API response without team_id", func(t *testing.T) {
		mapping := &api.GroupMapping{
			ID:      "mapping-id-123",
			Key:     "department",
			Value:   "engineering",
			GroupID: "group-id-456",
			TeamID:  nil,
		}

		state := &tfsdk.State{Schema: r.metaData.Schema}
		diags := r.UpdateState(ctx, state, nil, mapping)

		assert.False(t, diags.HasError())

		var model GroupMappingModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.True(t, model.TeamID.IsNull())
	})
}

func TestGroupMappingGetStateUpgraders(t *testing.T) {
	t.Run("should return empty state upgraders for version 0", func(t *testing.T) {
		r := &groupMappingResource{}
		upgraders := r.GetStateUpgraders(context.Background())
		assert.Empty(t, upgraders)
	})
}
