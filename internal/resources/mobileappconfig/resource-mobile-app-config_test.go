package mobileappconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMobileAppConfigResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewMobileAppConfigResourceHandle()

		assert.NotNil(t, handle)
		assert.NotNil(t, handle.MetaData())
		assert.Equal(t, ResourceInstanaMobileAppConfig, handle.MetaData().ResourceName)
		assert.Equal(t, int64(1), handle.MetaData().SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewMobileAppConfigResourceHandle()
		schema := handle.MetaData().Schema

		assert.NotNil(t, schema.Attributes)
		assert.Contains(t, schema.Attributes, "id")
		assert.Contains(t, schema.Attributes, "name")
	})
}

func TestMetaData(t *testing.T) {
	resource := &mobileAppConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaMobileAppConfig,
			Schema:        schema.Schema{},
			SchemaVersion: 1,
		},
	}

	t.Run("should return metadata", func(t *testing.T) {
		metaData := resource.MetaData()

		assert.NotNil(t, metaData)
		assert.Equal(t, ResourceInstanaMobileAppConfig, metaData.ResourceName)
		assert.Equal(t, int64(1), metaData.SchemaVersion)
	})
}

func TestSetComputedFields(t *testing.T) {
	resource := &mobileAppConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaMobileAppConfig,
			Schema:        schema.Schema{},
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should return no diagnostics", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		diags := resource.SetComputedFields(ctx, plan)

		assert.False(t, diags.HasError())
		assert.Len(t, diags, 0)
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &mobileAppConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaMobileAppConfig,
			Schema:        NewMobileAppConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map state to API object from plan", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		model := MobileAppConfigModel{
			ID:   types.StringValue("test-id-123"),
			Name: types.StringValue("Test Mobile App"),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "test-id-123", apiObject.ID)
		assert.Equal(t, "Test Mobile App", apiObject.Name)
	})

	t.Run("should return error when context is nil", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		apiObject, diags := resource.MapStateToDataObject(nil, plan, nil)

		assert.True(t, diags.HasError())
		assert.Nil(t, apiObject)
	})

	t.Run("should return error when both plan and state are nil", func(t *testing.T) {
		apiObject, diags := resource.MapStateToDataObject(ctx, nil, nil)

		assert.True(t, diags.HasError())
		assert.Nil(t, apiObject)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		model := MobileAppConfigModel{
			ID:   types.StringValue("test-id-123"),
			Name: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.True(t, diags.HasError())
		assert.Nil(t, apiObject)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &mobileAppConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaMobileAppConfig,
			Schema:        NewMobileAppConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with API object", func(t *testing.T) {
		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		apiObject := &api.MobileAppConfig{
			ID:   "test-id-456",
			Name: "Updated Mobile App",
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model MobileAppConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Equal(t, "test-id-456", model.ID.ValueString())
		assert.Equal(t, "Updated Mobile App", model.Name.ValueString())
	})

	t.Run("should return error when context is nil", func(t *testing.T) {
		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		apiObject := &api.MobileAppConfig{
			ID:   "test-id-456",
			Name: "Updated Mobile App",
		}

		diags := resource.UpdateState(nil, state, nil, apiObject)

		assert.True(t, diags.HasError())
	})

	t.Run("should return error when state is nil", func(t *testing.T) {
		apiObject := &api.MobileAppConfig{
			ID:   "test-id-456",
			Name: "Updated Mobile App",
		}

		diags := resource.UpdateState(ctx, nil, nil, apiObject)

		assert.True(t, diags.HasError())
	})

	t.Run("should return error when API object is nil", func(t *testing.T) {
		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, nil)

		assert.True(t, diags.HasError())
	})
}

func TestValidateModelFields(t *testing.T) {
	t.Run("should validate valid model", func(t *testing.T) {
		model := &MobileAppConfigModel{
			Name: types.StringValue("Valid Name"),
		}

		err := validateModelFields(model)

		assert.NoError(t, err)
	})

	t.Run("should return error for null name", func(t *testing.T) {
		model := &MobileAppConfigModel{
			Name: types.StringNull(),
		}

		err := validateModelFields(model)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name field is required")
	})

	t.Run("should return error for name too short", func(t *testing.T) {
		model := &MobileAppConfigModel{
			Name: types.StringValue(""),
		}

		err := validateModelFields(model)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name length must be between")
	})

	t.Run("should return error for name too long", func(t *testing.T) {
		longName := string(make([]byte, MobileAppConfigMaxNameLength+1))
		model := &MobileAppConfigModel{
			Name: types.StringValue(longName),
		}

		err := validateModelFields(model)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name length must be between")
	})
}

func TestMapAPIObjectToModel(t *testing.T) {
	t.Run("should map API object to model", func(t *testing.T) {
		apiObject := &api.MobileAppConfig{
			ID:   "api-id-789",
			Name: "API Mobile App",
		}

		model := mapAPIObjectToModel(apiObject)

		assert.NotNil(t, model)
		assert.Equal(t, "api-id-789", model.ID.ValueString())
		assert.Equal(t, "API Mobile App", model.Name.ValueString())
	})

	t.Run("should handle empty values", func(t *testing.T) {
		apiObject := &api.MobileAppConfig{
			ID:   "",
			Name: "",
		}

		model := mapAPIObjectToModel(apiObject)

		assert.NotNil(t, model)
		assert.True(t, model.ID.IsNull())
		assert.True(t, model.Name.IsNull())
	})
}
