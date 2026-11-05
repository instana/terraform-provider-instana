package businessperspectiveconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	tag "github.com/instana/instana-go-client/shared/tagfilter"
	sharedtypes "github.com/instana/instana-go-client/shared/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func emptyRbacTagsList() types.List {
	return types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			BusinessPerspectiveConfigFieldRbacTagID:          types.StringType,
			BusinessPerspectiveConfigFieldRbacTagDisplayName: types.StringType,
		},
	})
}

func createMockPlan(t *testing.T, model BusinessPerspectiveConfigModel) tfsdk.Plan {
	t.Helper()
	if model.RbacTags.IsUnknown() || (model.RbacTags.IsNull() && model.RbacTags.ElementType(context.Background()) == nil) {
		model.RbacTags = emptyRbacTagsList()
	}
	s := createResourceSchema()
	plan := tfsdk.Plan{
		Schema: s,
	}
	diags := plan.Set(context.Background(), &model)
	require.False(t, diags.HasError(), "Failed to create mock plan: %v", diags)
	return plan
}

func createMockState(t *testing.T, model BusinessPerspectiveConfigModel) tfsdk.State {
	t.Helper()
	if model.RbacTags.IsUnknown() || (model.RbacTags.IsNull() && model.RbacTags.ElementType(context.Background()) == nil) {
		model.RbacTags = emptyRbacTagsList()
	}
	s := createResourceSchema()
	state := tfsdk.State{
		Schema: s,
	}
	diags := state.Set(context.Background(), &model)
	require.False(t, diags.HasError(), "Failed to create mock state: %v", diags)
	return state
}

func TestNewBusinessPerspectiveConfigResourceHandle(t *testing.T) {
	handle := NewBusinessPerspectiveConfigResourceHandle()

	require.NotNil(t, handle)
	require.NotNil(t, handle.MetaData())
	assert.Equal(t, ResourceInstanaBusinessPerspectiveConfig, handle.MetaData().ResourceName)
	assert.Equal(t, int64(0), handle.MetaData().SchemaVersion)

	attrs := handle.MetaData().Schema.Attributes
	assert.Contains(t, attrs, BusinessPerspectiveConfigFieldID)
	assert.Contains(t, attrs, BusinessPerspectiveConfigFieldName)
	assert.Contains(t, attrs, BusinessPerspectiveConfigFieldDescription)
	assert.Contains(t, attrs, BusinessPerspectiveConfigFieldTagFilter)
	assert.Contains(t, attrs, BusinessPerspectiveConfigFieldRbacTags)

	assert.True(t, attrs[BusinessPerspectiveConfigFieldID].IsComputed())
	assert.True(t, attrs[BusinessPerspectiveConfigFieldName].IsRequired())
	assert.True(t, attrs[BusinessPerspectiveConfigFieldDescription].IsOptional())
	assert.True(t, attrs[BusinessPerspectiveConfigFieldDescription].IsComputed())
	assert.True(t, attrs[BusinessPerspectiveConfigFieldTagFilter].IsOptional())
	assert.True(t, attrs[BusinessPerspectiveConfigFieldRbacTags].IsOptional())
	assert.True(t, attrs[BusinessPerspectiveConfigFieldRbacTags].IsComputed())
}

func TestMetaData(t *testing.T) {
	resource := &businessPerspectiveConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_business_perspective",
			SchemaVersion: 0,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_business_perspective", metaData.ResourceName)
	assert.Equal(t, int64(0), metaData.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	handle := NewBusinessPerspectiveConfigResourceHandle()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	diags := handle.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetStateUpgraders(t *testing.T) {
	handle := NewBusinessPerspectiveConfigResourceHandle()
	ctx := context.Background()

	upgraders := handle.GetStateUpgraders(ctx)
	assert.Nil(t, upgraders)
}

func TestMapStateToDataObject_FullConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	rbacList, rbacDiags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           types.StringType,
			"display_name": types.StringType,
		},
	}, []RbacTagModel{
		{ID: types.StringValue("team-1"), DisplayName: types.StringValue("Team One")},
	})
	require.False(t, rbacDiags.HasError())

	model := BusinessPerspectiveConfigModel{
		ID:          types.StringValue("bp-id-123"),
		Name:        types.StringValue("Checkout Flow"),
		Description: types.StringValue("Tracks checkout operations"),
		TagFilter:   types.StringValue("service.name EQUALS 'checkout-service'"),
		RbacTags:    rbacList,
	}

	state := createMockState(t, model)

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Equal(t, "bp-id-123", result.ID)
	assert.Equal(t, "Checkout Flow", result.Name)
	assert.Equal(t, "Tracks checkout operations", result.Description)
	require.NotNil(t, result.TagFilterExpression)
	assert.Equal(t, tag.TagFilterType, result.TagFilterExpression.Type)
	require.Len(t, result.RbacTags, 1)
	assert.Equal(t, "team-1", result.RbacTags[0].ID)
	assert.Equal(t, "Team One", result.RbacTags[0].DisplayName)
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	model := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Order Processing"),
		Description: types.StringValue("Orders perspective"),
		TagFilter:   types.StringNull(),
	}

	plan := createMockPlan(t, model)

	result, diags := resource.MapStateToDataObject(ctx, &plan, nil)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Empty(t, result.ID)
	assert.Equal(t, "Order Processing", result.Name)
	assert.Equal(t, "Orders perspective", result.Description)
	assert.Nil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_MinimalConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	model := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Minimal Perspective"),
		Description: types.StringNull(),
		TagFilter:   types.StringNull(),
	}

	state := createMockState(t, model)

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Empty(t, result.ID)
	assert.Equal(t, "Minimal Perspective", result.Name)
	assert.Empty(t, result.Description)
	assert.Nil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_InvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	model := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Invalid Filter BP"),
		Description: types.StringNull(),
		TagFilter:   types.StringValue("invalid tag filter expression !!!"),
	}

	state := createMockState(t, model)

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
	assert.Nil(t, result)
}

func TestUpdateState_FullConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	tagFilter := tag.NewStringTagFilter(
		tag.TagFilterEntityDestination,
		"service.name",
		sharedtypes.EqualsOperator,
		"cart-service",
	)

	apiConfig := &api.BusinessPerspectiveConfig{
		ID:                  "bp-456",
		Name:                "Cart Perspective",
		Description:         "Cart operations",
		TagFilterExpression: tagFilter,
		RbacTags: []api.RbacTag{
			{ID: "team-cart", DisplayName: "Cart Team"},
		},
	}

	emptyModel := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringNull(),
		Description: types.StringNull(),
		TagFilter:   types.StringNull(),
		RbacTags:    types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType, "display_name": types.StringType}}),
	}

	state := createMockState(t, emptyModel)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)

	var updatedModel BusinessPerspectiveConfigModel
	diags = state.Get(ctx, &updatedModel)
	require.False(t, diags.HasError())

	assert.Equal(t, "bp-456", updatedModel.ID.ValueString())
	assert.Equal(t, "Cart Perspective", updatedModel.Name.ValueString())
	assert.Equal(t, "Cart operations", updatedModel.Description.ValueString())
	assert.False(t, updatedModel.TagFilter.IsNull())
	assert.Contains(t, updatedModel.TagFilter.ValueString(), "service.name")
	assert.False(t, updatedModel.RbacTags.IsNull())
	var rbacTags []RbacTagModel
	diags = updatedModel.RbacTags.ElementsAs(ctx, &rbacTags, false)
	require.False(t, diags.HasError())
	require.Len(t, rbacTags, 1)
	assert.Equal(t, "team-cart", rbacTags[0].ID.ValueString())
	assert.Equal(t, "Cart Team", rbacTags[0].DisplayName.ValueString())
}

func TestUpdateState_MinimalConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	apiConfig := &api.BusinessPerspectiveConfig{
		ID:   "bp-minimal-id",
		Name: "Minimal BP",
	}

	emptyModel := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringNull(),
		Description: types.StringNull(),
		TagFilter:   types.StringNull(),
	}

	state := createMockState(t, emptyModel)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)

	var updatedModel BusinessPerspectiveConfigModel
	diags = state.Get(ctx, &updatedModel)
	require.False(t, diags.HasError())

	assert.Equal(t, "bp-minimal-id", updatedModel.ID.ValueString())
	assert.Equal(t, "Minimal BP", updatedModel.Name.ValueString())
	assert.True(t, updatedModel.Description.IsNull())
	assert.True(t, updatedModel.TagFilter.IsNull())
}

func TestUpdateState_PreserveExistingTagFilterInPlan(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	tagFilter := tag.NewStringTagFilter(
		tag.TagFilterEntityDestination,
		"service.name",
		sharedtypes.EqualsOperator,
		"auth-service",
	)

	apiConfig := &api.BusinessPerspectiveConfig{
		ID:                  "bp-auth",
		Name:                "Auth BP",
		Description:         "Authentication perspective",
		TagFilterExpression: tagFilter,
	}

	customFilterStr := "service.name EQUALS 'auth-service'"
	planModel := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Auth BP"),
		Description: types.StringValue("Authentication perspective"),
		TagFilter:   types.StringValue(customFilterStr),
	}

	plan := createMockPlan(t, planModel)
	emptyState := createMockState(t, BusinessPerspectiveConfigModel{})

	diags := resource.UpdateState(ctx, &emptyState, &plan, apiConfig)
	require.False(t, diags.HasError())

	var updatedModel BusinessPerspectiveConfigModel
	diags = emptyState.Get(ctx, &updatedModel)
	require.False(t, diags.HasError())

	assert.Equal(t, customFilterStr, updatedModel.TagFilter.ValueString())
}

func TestUpdateState_ExplicitEmptyDescription(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()

	apiConfig := &api.BusinessPerspectiveConfig{
		ID:          "bp-empty-desc",
		Name:        "BP with empty desc",
		Description: "",
	}

	planModel := BusinessPerspectiveConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("BP with empty desc"),
		Description: types.StringValue(""),
		TagFilter:   types.StringNull(),
	}

	plan := createMockPlan(t, planModel)
	emptyState := createMockState(t, BusinessPerspectiveConfigModel{})

	diags := resource.UpdateState(ctx, &emptyState, &plan, apiConfig)
	require.False(t, diags.HasError())

	var updatedModel BusinessPerspectiveConfigModel
	diags = emptyState.Get(ctx, &updatedModel)
	require.False(t, diags.HasError())

	assert.Equal(t, "bp-empty-desc", updatedModel.ID.ValueString())
	assert.Equal(t, "", updatedModel.Description.ValueString())
}

func TestUpdateState_NilValidation(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()
	state := createMockState(t, BusinessPerspectiveConfigModel{})
	config := &api.BusinessPerspectiveConfig{ID: "test"}

	t.Run("nil context", func(t *testing.T) {
		diags := resource.UpdateState(nil, &state, nil, config)
		assert.True(t, diags.HasError())
	})

	t.Run("nil state", func(t *testing.T) {
		diags := resource.UpdateState(ctx, nil, nil, config)
		assert.True(t, diags.HasError())
	})

	t.Run("nil config", func(t *testing.T) {
		diags := resource.UpdateState(ctx, &state, nil, nil)
		assert.True(t, diags.HasError())
	})
}

func TestMapStateToDataObject_NilValidation(t *testing.T) {
	ctx := context.Background()
	resource := NewBusinessPerspectiveConfigResourceHandle()
	state := createMockState(t, BusinessPerspectiveConfigModel{Name: types.StringValue("test")})

	t.Run("nil context", func(t *testing.T) {
		result, diags := resource.MapStateToDataObject(nil, nil, &state)
		assert.True(t, diags.HasError())
		assert.Nil(t, result)
	})

	t.Run("nil plan and state", func(t *testing.T) {
		result, diags := resource.MapStateToDataObject(ctx, nil, nil)
		assert.True(t, diags.HasError())
		assert.Nil(t, result)
	})
}

type mockBPAPI struct {
	client.InstanaAPI
}

func (m *mockBPAPI) BusinessPerspectives() rest.RestResource[*api.BusinessPerspectiveConfig] {
	return nil
}

func TestGetRestResource(t *testing.T) {
	resource := NewBusinessPerspectiveConfigResourceHandle()
	mockAPI := &mockBPAPI{}
	restResource := resource.GetRestResource(mockAPI)
	assert.Nil(t, restResource)
}
