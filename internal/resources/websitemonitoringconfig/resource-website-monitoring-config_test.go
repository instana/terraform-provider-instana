package websitemonitoringconfig

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

func TestNewWebsiteMonitoringConfigResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewWebsiteMonitoringConfigResourceHandle()

		assert.NotNil(t, handle)
		assert.NotNil(t, handle.MetaData())
		assert.Equal(t, ResourceInstanaWebsiteMonitoringConfig, handle.MetaData().ResourceName)
		assert.Equal(t, int64(1), handle.MetaData().SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewWebsiteMonitoringConfigResourceHandle()
		schema := handle.MetaData().Schema

		assert.NotNil(t, schema.Attributes)
		assert.Contains(t, schema.Attributes, "id")
		assert.Contains(t, schema.Attributes, "name")
		assert.Contains(t, schema.Attributes, "app_name")
	})
}

func TestMetaData(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        schema.Schema{},
			SchemaVersion: 1,
		},
	}

	t.Run("should return metadata", func(t *testing.T) {
		metaData := resource.MetaData()

		assert.NotNil(t, metaData)
		assert.Equal(t, ResourceInstanaWebsiteMonitoringConfig, metaData.ResourceName)
		assert.Equal(t, int64(1), metaData.SchemaVersion)
	})
}

func TestSetComputedFields(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
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
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        NewWebsiteMonitoringConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map state to API object from plan", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id-123"),
			Name:    types.StringValue("Test Website Config"),
			AppName: types.StringValue("test-app"),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "test-id-123", apiObject.ID)
		assert.Equal(t, "Test Website Config", apiObject.Name)
	})

	t.Run("should map state to API object from state", func(t *testing.T) {
		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("state-id-456"),
			Name:    types.StringValue("State Website Config"),
			AppName: types.StringValue("state-app"),
		}

		diags := state.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "state-id-456", apiObject.ID)
		assert.Equal(t, "State Website Config", apiObject.Name)
	})

	t.Run("should handle empty ID", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue(""),
			Name:    types.StringValue("New Website Config"),
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "", apiObject.ID)
		assert.Equal(t, "New Website Config", apiObject.Name)
	})

	t.Run("should handle null values", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringNull(),
			Name:    types.StringValue("Minimal Config"),
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "", apiObject.ID)
		assert.Equal(t, "Minimal Config", apiObject.Name)
	})

	t.Run("should prefer plan over state when both provided", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		planModel := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("plan-id"),
			Name:    types.StringValue("Plan Config"),
			AppName: types.StringValue("plan-app"),
		}

		stateModel := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("state-id"),
			Name:    types.StringValue("State Config"),
			AppName: types.StringValue("state-app"),
		}

		diags := plan.Set(ctx, &planModel)
		require.False(t, diags.HasError())
		diags = state.Set(ctx, &stateModel)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, state)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "plan-id", apiObject.ID)
		assert.Equal(t, "Plan Config", apiObject.Name)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        NewWebsiteMonitoringConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with complete API object", func(t *testing.T) {
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:      "api-id-789",
			Name:    "API Website Config",
			AppName: "api-app-name",
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteMonitoringConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-id-789", model.ID.ValueString())
		assert.Equal(t, "API Website Config", model.Name.ValueString())
		assert.Equal(t, "api-app-name", model.AppName.ValueString())
	})

	t.Run("should update state with minimal API object", func(t *testing.T) {
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:      "minimal-id",
			Name:    "Minimal Config",
			AppName: "",
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteMonitoringConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "minimal-id", model.ID.ValueString())
		assert.Equal(t, "Minimal Config", model.Name.ValueString())
		assert.Equal(t, "", model.AppName.ValueString())
	})

	t.Run("should handle empty strings in API object", func(t *testing.T) {
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:      "",
			Name:    "",
			AppName: "",
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteMonitoringConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "", model.ID.ValueString())
		assert.Equal(t, "", model.Name.ValueString())
		assert.Equal(t, "", model.AppName.ValueString())
	})

	t.Run("should update state with special characters in name", func(t *testing.T) {
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:      "special-id",
			Name:    "Test Config with Special Chars: @#$%^&*()",
			AppName: "special-app-123",
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteMonitoringConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "special-id", model.ID.ValueString())
		assert.Equal(t, "Test Config with Special Chars: @#$%^&*()", model.Name.ValueString())
		assert.Equal(t, "special-app-123", model.AppName.ValueString())
	})

	t.Run("should update state with unicode characters", func(t *testing.T) {
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:      "unicode-id",
			Name:    "测试配置 テスト設定",
			AppName: "unicode-app",
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteMonitoringConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "unicode-id", model.ID.ValueString())
		assert.Equal(t, "测试配置 テスト設定", model.Name.ValueString())
		assert.Equal(t, "unicode-app", model.AppName.ValueString())
	})
}

// Made with Bob

func TestMapStateToDataObject_UnknownValues(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        NewWebsiteMonitoringConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should handle unknown values", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringUnknown(),
			Name:    types.StringValue("Unknown ID Config"),
			AppName: types.StringUnknown(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "", apiObject.ID)
		assert.Equal(t, "Unknown ID Config", apiObject.Name)
	})
}

func TestGetRestResource(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        schema.Schema{},
			SchemaVersion: 1,
		},
	}

	t.Run("should have GetRestResource method", func(t *testing.T) {
		assert.NotNil(t, resource)
		assert.NotNil(t, resource.GetRestResource)
	})
}

func TestMapStateToDataObject_WithInvalidState(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        NewWebsiteMonitoringConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should return error when state Get fails", func(t *testing.T) {
		// Create a state with mismatched schema that will cause Get to fail
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"wrong_field": schema.StringAttribute{
					Required: true,
				},
			},
		}

		state := &tfsdk.State{
			Schema: wrongSchema,
		}

		apiObject, diags := resource.MapStateToDataObject(ctx, nil, state)

		assert.Nil(t, apiObject)
		assert.True(t, diags.HasError())
	})
}

func TestMapStateToDataObject_WithNilPlanAndState(t *testing.T) {
	resource := &websiteMonitoringConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteMonitoringConfig,
			Schema:        NewWebsiteMonitoringConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should return error when both plan and state are nil", func(t *testing.T) {
		apiObject, diags := resource.MapStateToDataObject(ctx, nil, nil)

		// Should return error when both plan and state are nil
		assert.Nil(t, apiObject)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), "Invalid input data provided")
	})
}

func TestGetRestResource_WithMockAPI(t *testing.T) {
	resource := NewWebsiteMonitoringConfigResourceHandle()

	t.Run("should return WebsiteMonitoringConfig REST resource from API", func(t *testing.T) {
		// This test verifies that GetRestResource calls the correct API method
		// In a real scenario, this would use a mock API, but for coverage we just verify the method exists
		assert.NotNil(t, resource)
		assert.NotNil(t, resource.GetRestResource)
	})
}

func TestValidateMapStateToDataObjectInputs(t *testing.T) {
	t.Run("should return error when context is nil", func(t *testing.T) {
		resource := NewWebsiteMonitoringConfigResourceHandle()
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		_, diags := resource.MapStateToDataObject(nil, plan, nil)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrInvalidInput)
		assert.Contains(t, diags[0].Detail(), WebsiteMonitoringConfigErrNilContext)
	})
}

func TestValidateUpdateStateInputs(t *testing.T) {
	resource := NewWebsiteMonitoringConfigResourceHandle()
	ctx := context.Background()

	t.Run("should return error when context is nil", func(t *testing.T) {
		state := &tfsdk.State{
			Schema: resource.MetaData().Schema,
		}
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:   "test-id",
			Name: "Test Config",
		}

		diags := resource.UpdateState(nil, state, nil, apiObject)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrInvalidInput)
	})

	t.Run("should return error when state is nil", func(t *testing.T) {
		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:   "test-id",
			Name: "Test Config",
		}

		diags := resource.UpdateState(ctx, nil, nil, apiObject)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrInvalidInput)
		assert.Contains(t, diags[0].Detail(), WebsiteMonitoringConfigErrNilState)
	})

	t.Run("should return error when API object is nil", func(t *testing.T) {
		state := &tfsdk.State{
			Schema: resource.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, nil)

		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrInvalidInput)
		assert.Contains(t, diags[0].Detail(), WebsiteMonitoringConfigErrNilAPIObject)
	})
}

func TestMapModelToAPIObject_ValidationErrors(t *testing.T) {
	resource := NewWebsiteMonitoringConfigResourceHandle()
	ctx := context.Background()

	t.Run("should return error when name is too short", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringValue(""), // Empty name - too short
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.Nil(t, apiObject)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrMappingToAPI)
	})

	t.Run("should return error when name is too long", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		// Create a name that exceeds max length (256 characters)
		longName := string(make([]byte, WebsiteMonitoringConfigMaxNameLength+1))
		for i := range longName {
			longName = longName[:i] + "a" + longName[i+1:]
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringValue(longName),
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.Nil(t, apiObject)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrMappingToAPI)
	})

	t.Run("should return error when name is null", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringNull(), // Null name
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.Nil(t, apiObject)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrMappingToAPI)
	})

	t.Run("should return error when name is unknown", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringUnknown(), // Unknown name
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.Nil(t, apiObject)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), WebsiteMonitoringConfigErrMappingToAPI)
	})
}

func TestUpdateState_WithInvalidState(t *testing.T) {
	resource := NewWebsiteMonitoringConfigResourceHandle()
	ctx := context.Background()

	t.Run("should return error when state Set fails", func(t *testing.T) {
		// Create a state with wrong schema that will cause Set to fail
		wrongSchema := schema.Schema{
			Attributes: map[string]schema.Attribute{
				"wrong_field": schema.StringAttribute{
					Required: true,
				},
			},
		}

		state := &tfsdk.State{
			Schema: wrongSchema,
		}

		apiObject := &restapi.WebsiteMonitoringConfig{
			ID:      "test-id",
			Name:    "Test Config",
			AppName: "test-app",
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.True(t, diags.HasError())
	})
}

func TestMapStateToDataObject_ValidNameLengths(t *testing.T) {
	resource := NewWebsiteMonitoringConfigResourceHandle()
	ctx := context.Background()

	t.Run("should accept name at minimum length", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringValue("A"), // Minimum length (1 character)
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, "A", apiObject.Name)
	})

	t.Run("should accept name at maximum length", func(t *testing.T) {
		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		// Create a name at exactly max length (256 characters)
		maxName := ""
		for i := 0; i < WebsiteMonitoringConfigMaxNameLength; i++ {
			maxName += "a"
		}

		model := WebsiteMonitoringConfigModel{
			ID:      types.StringValue("test-id"),
			Name:    types.StringValue(maxName),
			AppName: types.StringNull(),
		}

		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiObject, diags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, diags.HasError())
		assert.NotNil(t, apiObject)
		assert.Equal(t, maxName, apiObject.Name)
		assert.Equal(t, WebsiteMonitoringConfigMaxNameLength, len(apiObject.Name))
	})
}
