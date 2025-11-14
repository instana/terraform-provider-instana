package slocorrectionconfig

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

func TestNewSloCorrectionConfigResourceHandleFramework(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewSloCorrectionConfigResourceHandleFramework()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSloCorrectionConfigFramework, metadata.ResourceName)
		assert.Equal(t, int64(1), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewSloCorrectionConfigResourceHandleFramework()
		metadata := handle.MetaData()

		schema := metadata.Schema
		assert.NotNil(t, schema.Attributes["id"])
		assert.NotNil(t, schema.Attributes["name"])
		assert.NotNil(t, schema.Attributes["description"])
		assert.NotNil(t, schema.Attributes["active"])
		assert.NotNil(t, schema.Attributes["slo_ids"])
		assert.NotNil(t, schema.Attributes["tags"])
		assert.NotNil(t, schema.Attributes["scheduling"])
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &sloCorrectionConfigResourceFramework{
			metaData: resourcehandle.ResourceMetaDataFramework{
				ResourceName:  ResourceInstanaSloCorrectionConfigFramework,
				SchemaVersion: 1,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSloCorrectionConfigFramework, metadata.ResourceName)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return SLO correction config rest resource", func(t *testing.T) {
		resource := &sloCorrectionConfigResourceFramework{}
		
		// Create a mock API to test the GetRestResource method
		mockAPI := &mockInstanaAPI{}
		restResource := resource.GetRestResource(mockAPI)
		
		assert.NotNil(t, restResource)
	})
}

// Mock API for testing - implements all required methods from InstanaAPI interface
type mockInstanaAPI struct{}

func (m *mockInstanaAPI) CustomEventSpecifications() restapi.RestResource[*restapi.CustomEventSpecification] { return nil }
func (m *mockInstanaAPI) BuiltinEventSpecifications() restapi.ReadOnlyRestResource[*restapi.BuiltinEventSpecification] { return nil }
func (m *mockInstanaAPI) APITokens() restapi.RestResource[*restapi.APIToken] { return nil }
func (m *mockInstanaAPI) ApplicationConfigs() restapi.RestResource[*restapi.ApplicationConfig] { return nil }
func (m *mockInstanaAPI) ApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] { return nil }
func (m *mockInstanaAPI) GlobalApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] { return nil }
func (m *mockInstanaAPI) AlertingChannels() restapi.RestResource[*restapi.AlertingChannel] { return nil }
func (m *mockInstanaAPI) AlertingConfigurations() restapi.RestResource[*restapi.AlertingConfiguration] { return nil }
func (m *mockInstanaAPI) SliConfigs() restapi.RestResource[*restapi.SliConfig] { return nil }
func (m *mockInstanaAPI) SloConfigs() restapi.RestResource[*restapi.SloConfig] { return nil }
func (m *mockInstanaAPI) SloAlertConfig() restapi.RestResource[*restapi.SloAlertConfig] { return nil }
func (m *mockInstanaAPI) SloCorrectionConfig() restapi.RestResource[*restapi.SloCorrectionConfig] {
	return &mockSloCorrectionConfigRestResource{}
}
func (m *mockInstanaAPI) WebsiteMonitoringConfig() restapi.RestResource[*restapi.WebsiteMonitoringConfig] { return nil }
func (m *mockInstanaAPI) WebsiteAlertConfig() restapi.RestResource[*restapi.WebsiteAlertConfig] { return nil }
func (m *mockInstanaAPI) InfraAlertConfig() restapi.RestResource[*restapi.InfraAlertConfig] { return nil }
func (m *mockInstanaAPI) Groups() restapi.RestResource[*restapi.Group] { return nil }
func (m *mockInstanaAPI) CustomDashboards() restapi.RestResource[*restapi.CustomDashboard] { return nil }
func (m *mockInstanaAPI) SyntheticTest() restapi.RestResource[*restapi.SyntheticTest] { return nil }
func (m *mockInstanaAPI) SyntheticLocation() restapi.ReadOnlyRestResource[*restapi.SyntheticLocation] { return nil }
func (m *mockInstanaAPI) SyntheticAlertConfigs() restapi.RestResource[*restapi.SyntheticAlertConfig] { return nil }
func (m *mockInstanaAPI) AutomationActions() restapi.RestResource[*restapi.AutomationAction] { return nil }
func (m *mockInstanaAPI) AutomationPolicies() restapi.RestResource[*restapi.AutomationPolicy] { return nil }
func (m *mockInstanaAPI) HostAgents() restapi.ReadOnlyRestResource[*restapi.HostAgent] { return nil }
func (m *mockInstanaAPI) LogAlertConfig() restapi.RestResource[*restapi.LogAlertConfig] { return nil }

// Mock rest resource - implements all required methods from RestResource interface
type mockSloCorrectionConfigRestResource struct{}

func (m *mockSloCorrectionConfigRestResource) GetAll() (*[]*restapi.SloCorrectionConfig, error) {
	return nil, nil
}

func (m *mockSloCorrectionConfigRestResource) GetOne(id string) (*restapi.SloCorrectionConfig, error) {
	return nil, nil
}

func (m *mockSloCorrectionConfigRestResource) Create(data *restapi.SloCorrectionConfig) (*restapi.SloCorrectionConfig, error) {
	return nil, nil
}

func (m *mockSloCorrectionConfigRestResource) Update(data *restapi.SloCorrectionConfig) (*restapi.SloCorrectionConfig, error) {
	return nil, nil
}

func (m *mockSloCorrectionConfigRestResource) Delete(data *restapi.SloCorrectionConfig) error {
	return nil
}

func (m *mockSloCorrectionConfigRestResource) DeleteByID(id string) error {
	return nil
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		resource := &sloCorrectionConfigResourceFramework{
			metaData: resourcehandle.ResourceMetaDataFramework{
				ResourceName:  ResourceInstanaSloCorrectionConfigFramework,
				Schema:        NewSloCorrectionConfigResourceHandleFramework().MetaData().Schema,
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
	resource := &sloCorrectionConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSloCorrectionConfigFramework,
			Schema:        NewSloCorrectionConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map complete model from state successfully", func(t *testing.T) {
		model := SloCorrectionConfigModel{
			ID:          types.StringValue("test-id-123"),
			Name:        types.StringValue("Test Correction"),
			Description: types.StringValue("Test Description"),
			Active:      types.BoolValue(true),
			Scheduling: &SchedulingModel{
				StartTime:     types.Int64Value(1741600800000),
				Duration:      types.Int64Value(60),
				DurationUnit:  types.StringValue("minute"),
				RecurrentRule: types.StringValue("FREQ=DAILY"),
				Recurrent:     types.BoolValue(true),
			},
		}

		sloIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"slo-1", "slo-2"})
		model.SloIds = sloIds

		tags, _ := types.SetValueFrom(ctx, types.StringType, []string{"tag1", "tag2"})
		model.Tags = tags

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "test-id-123", result.ID)
		assert.Equal(t, "Test Correction", result.Name)
		assert.Equal(t, "Test Description", result.Description)
		assert.True(t, result.Active)
		assert.Equal(t, int64(1741600800000), result.Scheduling.StartTime)
		assert.Equal(t, 60, result.Scheduling.Duration)
		assert.Equal(t, restapi.DurationUnit("MINUTE"), result.Scheduling.DurationUnit)
		assert.Equal(t, "FREQ=DAILY", result.Scheduling.RecurrentRule)
		assert.True(t, result.Scheduling.Recurrent)
		assert.ElementsMatch(t, []string{"slo-1", "slo-2"}, result.SloIds)
		assert.ElementsMatch(t, []string{"tag1", "tag2"}, result.Tags)
	})

	t.Run("should map model from plan successfully", func(t *testing.T) {
		model := SloCorrectionConfigModel{
			ID:          types.StringValue("plan-id"),
			Name:        types.StringValue("Plan Correction"),
			Description: types.StringValue("Plan Description"),
			Active:      types.BoolValue(false),
			Scheduling: &SchedulingModel{
				StartTime:     types.Int64Value(1741600900000),
				Duration:      types.Int64Value(120),
				DurationUnit:  types.StringValue("hour"),
				RecurrentRule: types.StringNull(),
				Recurrent:     types.BoolValue(false),
			},
		}

		sloIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"slo-3"})
		model.SloIds = sloIds
		model.Tags = types.SetNull(types.StringType)

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "plan-id", result.ID)
		assert.Equal(t, "Plan Correction", result.Name)
		assert.False(t, result.Active)
		assert.Equal(t, int64(1741600900000), result.Scheduling.StartTime)
		assert.Equal(t, 120, result.Scheduling.Duration)
		assert.Equal(t, restapi.DurationUnit("HOUR"), result.Scheduling.DurationUnit)
		assert.Empty(t, result.Scheduling.RecurrentRule)
		assert.False(t, result.Scheduling.Recurrent)
		assert.ElementsMatch(t, []string{"slo-3"}, result.SloIds)
		assert.Empty(t, result.Tags)
	})

	t.Run("should handle when both plan and state are nil", func(t *testing.T) {
		result, diags := resource.MapStateToDataObject(ctx, nil, nil)

		// The function returns an empty object when both are nil, not an error
		assert.NotNil(t, result)
		assert.False(t, diags.HasError())
	})

	t.Run("should handle null tags", func(t *testing.T) {
		model := SloCorrectionConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Active:      types.BoolValue(true),
			Scheduling: &SchedulingModel{
				StartTime:    types.Int64Value(1741600800000),
				Duration:     types.Int64Value(30),
				DurationUnit: types.StringValue("day"),
				Recurrent:    types.BoolValue(false),
			},
		}

		sloIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"slo-1"})
		model.SloIds = sloIds
		model.Tags = types.SetNull(types.StringType)

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Empty(t, result.Tags)
	})

	t.Run("should handle non-null recurrent rule", func(t *testing.T) {
		model := SloCorrectionConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Active:      types.BoolValue(true),
			Scheduling: &SchedulingModel{
				StartTime:     types.Int64Value(1741600800000),
				Duration:      types.Int64Value(60),
				DurationUnit:  types.StringValue("minute"),
				RecurrentRule: types.StringValue("FREQ=DAILY;INTERVAL=1"),
				Recurrent:     types.BoolValue(true),
			},
		}

		sloIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"slo-1"})
		model.SloIds = sloIds
		model.Tags = types.SetNull(types.StringType)

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "FREQ=DAILY;INTERVAL=1", result.Scheduling.RecurrentRule)
		assert.True(t, result.Scheduling.Recurrent)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &sloCorrectionConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSloCorrectionConfigFramework,
			Schema:        NewSloCorrectionConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with complete API object", func(t *testing.T) {
		apiObject := &restapi.SloCorrectionConfig{
			ID:          "api-id-123",
			Name:        "API Correction",
			Description: "API Description",
			Active:      true,
			Scheduling: restapi.Scheduling{
				StartTime:     1741600800000,
				Duration:      90,
				DurationUnit:  restapi.DurationUnit("HOUR"),
				RecurrentRule: "FREQ=WEEKLY",
				Recurrent:     true,
			},
			SloIds: []string{"slo-a", "slo-b"},
			Tags:   []string{"tag-a", "tag-b"},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SloCorrectionConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-id-123", model.ID.ValueString())
		assert.Equal(t, "API Correction", model.Name.ValueString())
		assert.Equal(t, "API Description", model.Description.ValueString())
		assert.True(t, model.Active.ValueBool())
		assert.Equal(t, int64(1741600800000), model.Scheduling.StartTime.ValueInt64())
		assert.Equal(t, int64(90), model.Scheduling.Duration.ValueInt64())
		assert.Equal(t, "HOUR", model.Scheduling.DurationUnit.ValueString())
		assert.Equal(t, "FREQ=WEEKLY", model.Scheduling.RecurrentRule.ValueString())
		assert.True(t, model.Scheduling.Recurrent.ValueBool())

		var sloIds []string
		diags = model.SloIds.ElementsAs(ctx, &sloIds, false)
		assert.False(t, diags.HasError())
		assert.ElementsMatch(t, []string{"slo-a", "slo-b"}, sloIds)

		var tags []string
		diags = model.Tags.ElementsAs(ctx, &tags, false)
		assert.False(t, diags.HasError())
		assert.ElementsMatch(t, []string{"tag-a", "tag-b"}, tags)
	})

	t.Run("should update state with empty recurrent rule", func(t *testing.T) {
		apiObject := &restapi.SloCorrectionConfig{
			ID:          "api-id-456",
			Name:        "No Recurrence",
			Description: "No Recurrence Description",
			Active:      false,
			Scheduling: restapi.Scheduling{
				StartTime:     1741600900000,
				Duration:      30,
				DurationUnit:  restapi.DurationUnit("DAY"),
				RecurrentRule: "",
				Recurrent:     false,
			},
			SloIds: []string{"slo-c"},
			Tags:   []string{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SloCorrectionConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-id-456", model.ID.ValueString())
		assert.False(t, model.Active.ValueBool())
		assert.True(t, model.Scheduling.RecurrentRule.IsNull())
		assert.False(t, model.Scheduling.Recurrent.ValueBool())
	})

	t.Run("should update state with null tags", func(t *testing.T) {
		apiObject := &restapi.SloCorrectionConfig{
			ID:          "api-id-789",
			Name:        "No Tags",
			Description: "No Tags Description",
			Active:      true,
			Scheduling: restapi.Scheduling{
				StartTime:    1741600800000,
				Duration:     15,
				DurationUnit: restapi.DurationUnit("MINUTE"),
				Recurrent:    false,
			},
			SloIds: []string{"slo-d"},
			Tags:   []string{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SloCorrectionConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.True(t, model.Tags.IsNull())
	})

	t.Run("should update state with different duration units", func(t *testing.T) {
		durationUnits := []restapi.DurationUnit{"MINUTE", "HOUR", "DAY"}

		for _, unit := range durationUnits {
			t.Run("duration_unit_"+string(unit), func(t *testing.T) {
				apiObject := &restapi.SloCorrectionConfig{
					ID:          "test-id",
					Name:        "Test",
					Description: "Desc",
					Active:      true,
					Scheduling: restapi.Scheduling{
						StartTime:    1741600800000,
						Duration:     10,
						DurationUnit: unit,
						Recurrent:    false,
					},
					SloIds: []string{"slo-1"},
				}

				state := &tfsdk.State{
					Schema: resource.metaData.Schema,
				}

				diags := resource.UpdateState(ctx, state, nil, apiObject)

				assert.False(t, diags.HasError())

				var model SloCorrectionConfigModel
				diags = state.Get(ctx, &model)
				assert.False(t, diags.HasError())

				assert.Equal(t, string(unit), model.Scheduling.DurationUnit.ValueString())
			})
		}
	})

	t.Run("should handle large duration values", func(t *testing.T) {
		apiObject := &restapi.SloCorrectionConfig{
			ID:          "test-id",
			Name:        "Large Duration",
			Description: "Desc",
			Active:      true,
			Scheduling: restapi.Scheduling{
				StartTime:    1741600800000,
				Duration:     999999,
				DurationUnit: restapi.DurationUnit("MINUTE"),
				Recurrent:    false,
			},
			SloIds: []string{"slo-1"},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SloCorrectionConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, int64(999999), model.Scheduling.Duration.ValueInt64())
	})

	t.Run("should handle multiple tags", func(t *testing.T) {
		apiObject := &restapi.SloCorrectionConfig{
			ID:          "test-id",
			Name:        "Multiple Tags",
			Description: "Desc",
			Active:      true,
			Scheduling: restapi.Scheduling{
				StartTime:    1741600800000,
				Duration:     60,
				DurationUnit: restapi.DurationUnit("MINUTE"),
				Recurrent:    false,
			},
			SloIds: []string{"slo-1"},
			Tags:   []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SloCorrectionConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		var tags []string
		diags = model.Tags.ElementsAs(ctx, &tags, false)
		assert.False(t, diags.HasError())
		assert.Len(t, tags, 5)
		assert.ElementsMatch(t, []string{"tag1", "tag2", "tag3", "tag4", "tag5"}, tags)
	})
}

func TestTagsHandling(t *testing.T) {
	resource := &sloCorrectionConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSloCorrectionConfigFramework,
			Schema:        NewSloCorrectionConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should handle empty tags list", func(t *testing.T) {
		model := SloCorrectionConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Active:      types.BoolValue(true),
			Scheduling: &SchedulingModel{
				StartTime:    types.Int64Value(1741600800000),
				Duration:     types.Int64Value(60),
				DurationUnit: types.StringValue("minute"),
				Recurrent:    types.BoolValue(false),
			},
		}

		sloIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"slo-1"})
		model.SloIds = sloIds

		tags, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		model.Tags = tags

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Empty(t, result.Tags)
	})

	t.Run("should handle single tag", func(t *testing.T) {
		model := SloCorrectionConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Active:      types.BoolValue(true),
			Scheduling: &SchedulingModel{
				StartTime:    types.Int64Value(1741600800000),
				Duration:     types.Int64Value(60),
				DurationUnit: types.StringValue("minute"),
				Recurrent:    types.BoolValue(false),
			},
		}

		sloIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"slo-1"})
		model.SloIds = sloIds

		tags, _ := types.SetValueFrom(ctx, types.StringType, []string{"single-tag"})
		model.Tags = tags

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Tags, 1)
		assert.Equal(t, "single-tag", result.Tags[0])
	})
}


// Made with Bob
