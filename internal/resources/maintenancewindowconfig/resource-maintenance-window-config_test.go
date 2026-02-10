package maintenancewindowconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMaintenanceWindowConfigResourceHandle(t *testing.T) {
	handle := NewMaintenanceWindowConfigResourceHandle()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaMaintenanceWindowConfig, metadata.ResourceName)
	assert.Equal(t, int64(0), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &maintenanceWindowConfigResource{
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

func TestSetComputedFields(t *testing.T) {
	resource := &maintenanceWindowConfigResource{}
	ctx := context.Background()

	diags := resource.SetComputedFields(ctx, nil)
	assert.False(t, diags.HasError())
}

func TestUpdateState(t *testing.T) {
	resource := &maintenanceWindowConfigResource{}
	ctx := context.Background()

	t.Run("one-time maintenance window", func(t *testing.T) {
		apiConfig := &restapi.MaintenanceWindowConfig{
			ID:    "test-id",
			Name:  "test-maintenance",
			Query: "entity.type:host",
			Scheduling: &restapi.MaintenanceScheduling{
				Start: 1698938631036,
				Duration: &restapi.MaintenanceDuration{
					Amount: 2,
					Unit:   "HOURS",
				},
				Type: "ONE_TIME",
			},
		}

		handle := NewMaintenanceWindowConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model MaintenanceWindowConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-maintenance", model.Name.ValueString())
		assert.Equal(t, "entity.type:host", model.Query.ValueString())
		assert.False(t, model.Scheduling.IsNull())
	})

	t.Run("recurrent maintenance window", func(t *testing.T) {
		rrule := "FREQ=WEEKLY;INTERVAL=2;BYDAY=SA;COUNT=10"
		timezoneId := "America/New_York"
		apiConfig := &restapi.MaintenanceWindowConfig{
			ID:    "test-id",
			Name:  "test-maintenance",
			Query: "",
			Scheduling: &restapi.MaintenanceScheduling{
				Start: 1683827571245,
				Duration: &restapi.MaintenanceDuration{
					Amount: 2,
					Unit:   "HOURS",
				},
				Type:       "RECURRENT",
				Rrule:      &rrule,
				TimezoneId: &timezoneId,
			},
		}

		handle := NewMaintenanceWindowConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model MaintenanceWindowConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-maintenance", model.Name.ValueString())
		assert.False(t, model.Scheduling.IsNull())
	})

	t.Run("with tag filter expression enabled", func(t *testing.T) {
		enabled := true
		apiConfig := &restapi.MaintenanceWindowConfig{
			ID:    "test-id",
			Name:  "test-maintenance",
			Query: "",
			Scheduling: &restapi.MaintenanceScheduling{
				Start: 1698938631036,
				Duration: &restapi.MaintenanceDuration{
					Amount: 1,
					Unit:   "DAYS",
				},
				Type: "ONE_TIME",
			},
			TagFilterExpressionEnabled: &enabled,
		}

		handle := NewMaintenanceWindowConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model MaintenanceWindowConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.True(t, model.TagFilterExpressionEnabled.ValueBool())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &maintenanceWindowConfigResource{}
	ctx := context.Background()

	t.Run("one-time maintenance window", func(t *testing.T) {
		durationAttrs := map[string]attr.Value{
			DurationFieldAmount: types.Int64Value(2),
			DurationFieldUnit:   types.StringValue("HOURS"),
		}
		durationObj, _ := types.ObjectValue(
			map[string]attr.Type{
				DurationFieldAmount: types.Int64Type,
				DurationFieldUnit:   types.StringType,
			},
			durationAttrs,
		)

		schedulingAttrs := map[string]attr.Value{
			SchedulingFieldStart:      types.Int64Value(1698938631036),
			SchedulingFieldDuration:   durationObj,
			SchedulingFieldType:       types.StringValue("ONE_TIME"),
			SchedulingFieldRrule:      types.StringNull(),
			SchedulingFieldTimezoneId: types.StringNull(),
		}
		schedulingObj, _ := types.ObjectValue(
			map[string]attr.Type{
				SchedulingFieldStart: types.Int64Type,
				SchedulingFieldDuration: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						DurationFieldAmount: types.Int64Type,
						DurationFieldUnit:   types.StringType,
					},
				},
				SchedulingFieldType:       types.StringType,
				SchedulingFieldRrule:      types.StringType,
				SchedulingFieldTimezoneId: types.StringType,
			},
			schedulingAttrs,
		)

		model := MaintenanceWindowConfigModel{
			ID:                         types.StringValue("test-id"),
			Name:                       types.StringValue("test-maintenance"),
			Query:                      types.StringValue("entity.type:host"),
			Scheduling:                 schedulingObj,
			TagFilterExpressionEnabled: types.BoolNull(),
			TagFilterExpression:        types.StringNull(),
		}

		handle := NewMaintenanceWindowConfigResourceHandle()
		plan := &tfsdk.Plan{
			Schema: handle.MetaData().Schema,
		}
		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiConfig, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, apiConfig)
		assert.Equal(t, "test-id", apiConfig.ID)
		assert.Equal(t, "test-maintenance", apiConfig.Name)
		assert.Equal(t, "entity.type:host", apiConfig.Query)
		require.NotNil(t, apiConfig.Scheduling)
		assert.Equal(t, int64(1698938631036), apiConfig.Scheduling.Start)
		assert.Equal(t, "ONE_TIME", apiConfig.Scheduling.Type)
		require.NotNil(t, apiConfig.Scheduling.Duration)
		assert.Equal(t, int64(2), apiConfig.Scheduling.Duration.Amount)
		assert.Equal(t, "HOURS", apiConfig.Scheduling.Duration.Unit)
	})

	t.Run("recurrent maintenance window", func(t *testing.T) {
		durationAttrs := map[string]attr.Value{
			DurationFieldAmount: types.Int64Value(3),
			DurationFieldUnit:   types.StringValue("DAYS"),
		}
		durationObj, _ := types.ObjectValue(
			map[string]attr.Type{
				DurationFieldAmount: types.Int64Type,
				DurationFieldUnit:   types.StringType,
			},
			durationAttrs,
		)

		schedulingAttrs := map[string]attr.Value{
			SchedulingFieldStart:      types.Int64Value(1683827571245),
			SchedulingFieldDuration:   durationObj,
			SchedulingFieldType:       types.StringValue("RECURRENT"),
			SchedulingFieldRrule:      types.StringValue("FREQ=WEEKLY;INTERVAL=2;BYDAY=SA;COUNT=10"),
			SchedulingFieldTimezoneId: types.StringValue("America/New_York"),
		}
		schedulingObj, _ := types.ObjectValue(
			map[string]attr.Type{
				SchedulingFieldStart: types.Int64Type,
				SchedulingFieldDuration: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						DurationFieldAmount: types.Int64Type,
						DurationFieldUnit:   types.StringType,
					},
				},
				SchedulingFieldType:       types.StringType,
				SchedulingFieldRrule:      types.StringType,
				SchedulingFieldTimezoneId: types.StringType,
			},
			schedulingAttrs,
		)

		model := MaintenanceWindowConfigModel{
			ID:                         types.StringValue("test-id"),
			Name:                       types.StringValue("test-maintenance"),
			Query:                      types.StringValue(""),
			Scheduling:                 schedulingObj,
			TagFilterExpressionEnabled: types.BoolValue(false),
			TagFilterExpression:        types.StringNull(),
		}

		handle := NewMaintenanceWindowConfigResourceHandle()
		plan := &tfsdk.Plan{
			Schema: handle.MetaData().Schema,
		}
		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiConfig, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, apiConfig)
		assert.Equal(t, "RECURRENT", apiConfig.Scheduling.Type)
		require.NotNil(t, apiConfig.Scheduling.Rrule)
		assert.Equal(t, "FREQ=WEEKLY;INTERVAL=2;BYDAY=SA;COUNT=10", *apiConfig.Scheduling.Rrule)
		require.NotNil(t, apiConfig.Scheduling.TimezoneId)
		assert.Equal(t, "America/New_York", *apiConfig.Scheduling.TimezoneId)
	})

	t.Run("with tag filter expression", func(t *testing.T) {
		durationAttrs := map[string]attr.Value{
			DurationFieldAmount: types.Int64Value(1),
			DurationFieldUnit:   types.StringValue("MINUTES"),
		}
		durationObj, _ := types.ObjectValue(
			map[string]attr.Type{
				DurationFieldAmount: types.Int64Type,
				DurationFieldUnit:   types.StringType,
			},
			durationAttrs,
		)

		schedulingAttrs := map[string]attr.Value{
			SchedulingFieldStart:      types.Int64Value(1698938631036),
			SchedulingFieldDuration:   durationObj,
			SchedulingFieldType:       types.StringValue("ONE_TIME"),
			SchedulingFieldRrule:      types.StringNull(),
			SchedulingFieldTimezoneId: types.StringNull(),
		}
		schedulingObj, _ := types.ObjectValue(
			map[string]attr.Type{
				SchedulingFieldStart: types.Int64Type,
				SchedulingFieldDuration: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						DurationFieldAmount: types.Int64Type,
						DurationFieldUnit:   types.StringType,
					},
				},
				SchedulingFieldType:       types.StringType,
				SchedulingFieldRrule:      types.StringType,
				SchedulingFieldTimezoneId: types.StringType,
			},
			schedulingAttrs,
		)

		model := MaintenanceWindowConfigModel{
			ID:                         types.StringValue("test-id"),
			Name:                       types.StringValue("test-maintenance"),
			Query:                      types.StringValue(""),
			Scheduling:                 schedulingObj,
			TagFilterExpressionEnabled: types.BoolValue(true),
			TagFilterExpression:        types.StringValue("synthetic.locationLabelAggregated@na EQUALS 'us-east'"),
		}

		handle := NewMaintenanceWindowConfigResourceHandle()
		plan := &tfsdk.Plan{
			Schema: handle.MetaData().Schema,
		}
		diags := plan.Set(ctx, &model)
		require.False(t, diags.HasError())

		apiConfig, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, apiConfig)
		require.NotNil(t, apiConfig.TagFilterExpressionEnabled)
		assert.True(t, *apiConfig.TagFilterExpressionEnabled)
		require.NotNil(t, apiConfig.TagFilterExpression)
	})
}

func TestGetStateUpgraders(t *testing.T) {
	resource := &maintenanceWindowConfigResource{}
	ctx := context.Background()

	upgraders := resource.GetStateUpgraders(ctx)
	assert.Nil(t, upgraders)
}

// Made with Bob
