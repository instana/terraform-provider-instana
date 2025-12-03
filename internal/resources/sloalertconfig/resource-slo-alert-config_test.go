package sloalertconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateState_StatusAlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id",
		Name:        "Test SLO Alert",
		Description: "Test Description",
		Severity:    5,
		Triggering:  true,
		Enabled:     true,
		Rule: restapi.SloAlertRule{
			AlertType: "SERVICE_LEVELS_OBJECTIVE",
			Metric:    "STATUS",
		},
		Threshold: &restapi.SloAlertThreshold{
			Type:     "static",
			Operator: ">=",
			Value:    95.0,
		},
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 300000,
			Expiry:     600000,
		},
		SloIds:                []string{"slo-1", "slo-2"},
		AlertChannelIds:       []string{"channel-1"},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		BurnRateConfigs:       &[]restapi.BurnRateConfig{},
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test SLO Alert", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	assert.Equal(t, int64(5), model.Severity.ValueInt64())
	assert.True(t, model.Triggering.ValueBool())
	assert.Equal(t, "status", model.AlertType.ValueString())
	assert.NotNil(t, model.Threshold)
	assert.Equal(t, "staticThreshold", model.Threshold.Type.ValueString())
	assert.Equal(t, ">=", model.Threshold.Operator.ValueString())
	assert.Equal(t, 95.0, model.Threshold.Value.ValueFloat64())
	assert.Equal(t, int64(300000), model.TimeThreshold.WarmUp.ValueInt64())
	assert.Equal(t, int64(600000), model.TimeThreshold.CoolDown.ValueInt64())
}

func TestUpdateState_ErrorBudgetAlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id-2",
		Name:        "Error Budget Alert",
		Description: "Error Budget Description",
		Severity:    8,
		Triggering:  false,
		Enabled:     true,
		Rule: restapi.SloAlertRule{
			AlertType: "ERROR_BUDGET",
			Metric:    "BURNED_PERCENTAGE",
		},
		Threshold: &restapi.SloAlertThreshold{
			Type:     "staticThreshold",
			Operator: "<=",
			Value:    10.0,
		},
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 60000,
			Expiry:     120000,
		},
		SloIds:                []string{"slo-3"},
		AlertChannelIds:       []string{"channel-2", "channel-3"},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		BurnRateConfigs:       &[]restapi.BurnRateConfig{},
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id-2", model.ID.ValueString())
	assert.Equal(t, "Error Budget Alert", model.Name.ValueString())
	assert.Equal(t, "error_budget", model.AlertType.ValueString())
	assert.Equal(t, int64(8), model.Severity.ValueInt64())
	assert.False(t, model.Triggering.ValueBool())
	assert.NotNil(t, model.Threshold)
	assert.Equal(t, "<=", model.Threshold.Operator.ValueString())
	assert.Equal(t, 10.0, model.Threshold.Value.ValueFloat64())
}

func TestUpdateState_BurnRateV2AlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	burnRateConfigs := []restapi.BurnRateConfig{
		{
			AlertWindowType:  "SHORT",
			Duration:         60,
			DurationUnitType: "MINUTES",
			Threshold: restapi.ServiceLevelsStaticThresholdConfig{
				Operator: ">",
				Value:    2.5,
			},
		},
		{
			AlertWindowType:  "LONG",
			Duration:         360,
			DurationUnitType: "MINUTES",
			Threshold: restapi.ServiceLevelsStaticThresholdConfig{
				Operator: ">=",
				Value:    1.5,
			},
		},
	}

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id-3",
		Name:        "Burn Rate Alert",
		Description: "Burn Rate Description",
		Severity:    10,
		Triggering:  true,
		Enabled:     true,
		Rule: restapi.SloAlertRule{
			AlertType: "ERROR_BUDGET",
			Metric:    "BURN_RATE_V2",
		},
		Threshold: nil,
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 180000,
			Expiry:     360000,
		},
		SloIds:                []string{"slo-4", "slo-5"},
		AlertChannelIds:       []string{"channel-4"},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		BurnRateConfigs:       &burnRateConfigs,
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id-3", model.ID.ValueString())
	assert.Equal(t, "Burn Rate Alert", model.Name.ValueString())
	assert.Equal(t, "burn_rate_v2", model.AlertType.ValueString())
	assert.Nil(t, model.Threshold)
	assert.NotNil(t, model.BurnRateConfig)
	assert.Len(t, model.BurnRateConfig, 2)

	// Check first burn rate config
	assert.Equal(t, "SHORT", model.BurnRateConfig[0].AlertWindowType.ValueString())
	assert.Equal(t, "60", model.BurnRateConfig[0].Duration.ValueString())
	assert.Equal(t, "MINUTES", model.BurnRateConfig[0].DurationUnitType.ValueString())
	assert.Equal(t, ">", model.BurnRateConfig[0].ThresholdOperator.ValueString())
	assert.Equal(t, "2.50", model.BurnRateConfig[0].ThresholdValue.ValueString())

	// Check second burn rate config
	assert.Equal(t, "LONG", model.BurnRateConfig[1].AlertWindowType.ValueString())
	assert.Equal(t, "360", model.BurnRateConfig[1].Duration.ValueString())
	assert.Equal(t, ">=", model.BurnRateConfig[1].ThresholdOperator.ValueString())
}

func TestUpdateState_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	customPayloadFields := []restapi.CustomPayloadField[any]{
		{
			Key:   "environment",
			Value: "production",
			Type:  restapi.StaticStringCustomPayloadType,
		},
		{
			Key:   "team",
			Value: "platform",
			Type:  restapi.StaticStringCustomPayloadType,
		},
	}

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id-4",
		Name:        "Alert with Custom Payload",
		Description: "Custom Payload Description",
		Severity:    7,
		Triggering:  false,
		Enabled:     true,
		Rule: restapi.SloAlertRule{
			AlertType: "SERVICE_LEVELS_OBJECTIVE",
			Metric:    "STATUS",
		},
		Threshold: &restapi.SloAlertThreshold{
			Type:     "static",
			Operator: ">=",
			Value:    99.0,
		},
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 300000,
			Expiry:     600000,
		},
		SloIds:                []string{"slo-6"},
		AlertChannelIds:       []string{"channel-5"},
		CustomerPayloadFields: customPayloadFields,
		BurnRateConfigs:       &[]restapi.BurnRateConfig{},
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id-4", model.ID.ValueString())
	assert.False(t, model.CustomPayload.IsNull())
}

func TestUpdateState_WithNullThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id-5",
		Name:        "Alert without Threshold",
		Description: "No Threshold Description",
		Severity:    6,
		Triggering:  true,
		Enabled:     false,
		Rule: restapi.SloAlertRule{
			AlertType: "ERROR_BUDGET",
			Metric:    "BURN_RATE_V2",
		},
		Threshold: nil,
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 120000,
			Expiry:     240000,
		},
		SloIds:                []string{"slo-7"},
		AlertChannelIds:       []string{"channel-6"},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		BurnRateConfigs:       &[]restapi.BurnRateConfig{},
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id-5", model.ID.ValueString())
	assert.Nil(t, model.Threshold)
}

func TestUpdateState_WithEmptyBurnRateConfigs(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id-6",
		Name:        "Alert with Empty Burn Rate",
		Description: "Empty Burn Rate Description",
		Severity:    5,
		Triggering:  false,
		Enabled:     true,
		Rule: restapi.SloAlertRule{
			AlertType: "SERVICE_LEVELS_OBJECTIVE",
			Metric:    "STATUS",
		},
		Threshold: &restapi.SloAlertThreshold{
			Type:     "static",
			Operator: ">=",
			Value:    98.0,
		},
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 300000,
			Expiry:     600000,
		},
		SloIds:                []string{"slo-8"},
		AlertChannelIds:       []string{"channel-7"},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		BurnRateConfigs:       &[]restapi.BurnRateConfig{},
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id-6", model.ID.ValueString())
	assert.Nil(t, model.BurnRateConfig)
}

func TestUpdateState_WithMultipleSloIds(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	apiConfig := &restapi.SloAlertConfig{
		ID:          "test-id-7",
		Name:        "Multi SLO Alert",
		Description: "Multiple SLOs",
		Severity:    9,
		Triggering:  true,
		Enabled:     true,
		Rule: restapi.SloAlertRule{
			AlertType: "SERVICE_LEVELS_OBJECTIVE",
			Metric:    "STATUS",
		},
		Threshold: &restapi.SloAlertThreshold{
			Type:     "static",
			Operator: ">",
			Value:    90.0,
		},
		TimeThreshold: restapi.SloAlertTimeThreshold{
			TimeWindow: 180000,
			Expiry:     360000,
		},
		SloIds:                []string{"slo-9", "slo-10", "slo-11", "slo-12"},
		AlertChannelIds:       []string{"channel-8", "channel-9"},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		BurnRateConfigs:       &[]restapi.BurnRateConfig{},
	}

	state := tfsdk.State{
		Schema: resource.MetaData().Schema,
	}
	initializeEmptyState(ctx, &state)

	diags := resource.UpdateState(ctx, &state, nil, apiConfig)
	require.False(t, diags.HasError())

	var model SloAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id-7", model.ID.ValueString())
	assert.Equal(t, 4, len(model.SloIds.Elements()))
	assert.Equal(t, 2, len(model.AlertChannelIds.Elements()))
}

func TestMapStateToDataObject_StatusAlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test SLO Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("status"),
		Threshold: &SloAlertThresholdModel{
			Type:     types.StringValue("staticThreshold"),
			Operator: types.StringValue(">="),
			Value:    types.Float64Value(95.0),
		},
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
			types.StringValue("slo-2"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		BurnRateConfig: nil,
		CustomPayload:  types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.Equal(t, "test-id", apiConfig.ID)
	assert.Equal(t, "Test SLO Alert", apiConfig.Name)
	assert.Equal(t, "Test Description", apiConfig.Description)
	assert.Equal(t, 5, apiConfig.Severity)
	assert.True(t, apiConfig.Triggering)
	assert.True(t, apiConfig.Enabled)
	assert.Equal(t, "SERVICE_LEVELS_OBJECTIVE", apiConfig.Rule.AlertType)
	assert.Equal(t, "STATUS", apiConfig.Rule.Metric)
	assert.NotNil(t, apiConfig.Threshold)
	assert.Equal(t, "staticThreshold", apiConfig.Threshold.Type)
	assert.Equal(t, ">=", apiConfig.Threshold.Operator)
	assert.Equal(t, 95.0, apiConfig.Threshold.Value)
	assert.Equal(t, 300000, apiConfig.TimeThreshold.TimeWindow)
	assert.Equal(t, 600000, apiConfig.TimeThreshold.Expiry)
	assert.Len(t, apiConfig.SloIds, 2)
	assert.Len(t, apiConfig.AlertChannelIds, 1)
}

func TestMapStateToDataObject_ErrorBudgetAlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id-2"),
		Name:        types.StringValue("Error Budget Alert"),
		Description: types.StringValue("Error Budget Description"),
		Severity:    types.Int64Value(8),
		Triggering:  types.BoolValue(false),

		AlertType: types.StringValue("error_budget"),
		Threshold: &SloAlertThresholdModel{
			Type:     types.StringValue("staticThreshold"),
			Operator: types.StringValue("<="),
			Value:    types.Float64Value(10.0),
		},
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(60000),
			CoolDown: types.Int64Value(120000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-3"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-2"),
			types.StringValue("channel-3"),
		}),
		BurnRateConfig: nil,
		CustomPayload:  types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.Equal(t, "test-id-2", apiConfig.ID)
	assert.Equal(t, "Error Budget Alert", apiConfig.Name)
	assert.Equal(t, "ERROR_BUDGET", apiConfig.Rule.AlertType)
	assert.Equal(t, "BURNED_PERCENTAGE", apiConfig.Rule.Metric)
	assert.Equal(t, 8, apiConfig.Severity)
	assert.False(t, apiConfig.Triggering)
	assert.NotNil(t, apiConfig.Threshold)
	assert.Equal(t, "<=", apiConfig.Threshold.Operator)
	assert.Equal(t, 10.0, apiConfig.Threshold.Value)
}

func TestMapStateToDataObject_BurnRateV2AlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id-3"),
		Name:        types.StringValue("Burn Rate Alert"),
		Description: types.StringValue("Burn Rate Description"),
		Severity:    types.Int64Value(10),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("burn_rate_v2"),
		Threshold: nil,
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(180000),
			CoolDown: types.Int64Value(360000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-4"),
			types.StringValue("slo-5"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-4"),
		}),
		BurnRateConfig: []SloAlertBurnRateConfigModel{
			{
				AlertWindowType:   types.StringValue("SHORT"),
				Duration:          types.StringValue("60"),
				DurationUnitType:  types.StringValue("MINUTES"),
				ThresholdOperator: types.StringValue(">"),
				ThresholdValue:    types.StringValue("2.5"),
			},
			{
				AlertWindowType:   types.StringValue("LONG"),
				Duration:          types.StringValue("360"),
				DurationUnitType:  types.StringValue("MINUTES"),
				ThresholdOperator: types.StringValue(">="),
				ThresholdValue:    types.StringValue("1.5"),
			},
		},
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.Equal(t, "test-id-3", apiConfig.ID)
	assert.Equal(t, "Burn Rate Alert", apiConfig.Name)
	assert.Equal(t, "ERROR_BUDGET", apiConfig.Rule.AlertType)
	assert.Equal(t, "BURN_RATE_V2", apiConfig.Rule.Metric)
	assert.Nil(t, apiConfig.Threshold)
	assert.NotNil(t, apiConfig.BurnRateConfigs)
	assert.Len(t, *apiConfig.BurnRateConfigs, 2)

	// Check first burn rate config
	burnRateConfigs := *apiConfig.BurnRateConfigs
	assert.Equal(t, "SHORT", burnRateConfigs[0].AlertWindowType)
	assert.Equal(t, 60, burnRateConfigs[0].Duration)
	assert.Equal(t, "MINUTES", burnRateConfigs[0].DurationUnitType)
	assert.Equal(t, ">", burnRateConfigs[0].Threshold.Operator)
	assert.Equal(t, 2.5, burnRateConfigs[0].Threshold.Value)

	// Check second burn rate config
	assert.Equal(t, "LONG", burnRateConfigs[1].AlertWindowType)
	assert.Equal(t, 360, burnRateConfigs[1].Duration)
	assert.Equal(t, ">=", burnRateConfigs[1].Threshold.Operator)
	assert.Equal(t, 1.5, burnRateConfigs[1].Threshold.Value)
}

func TestMapStateToDataObject_WithNormalizedAlertTypes(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	testCases := []struct {
		name              string
		inputAlertType    string
		expectedAlertType string
		expectedMetric    string
	}{
		{
			name:              "errorBudget normalized",
			inputAlertType:    "errorBudget",
			expectedAlertType: "ERROR_BUDGET",
			expectedMetric:    "BURNED_PERCENTAGE",
		},
		{
			name:              "ErrorBudget normalized",
			inputAlertType:    "ErrorBudget",
			expectedAlertType: "ERROR_BUDGET",
			expectedMetric:    "BURNED_PERCENTAGE",
		},
		{
			name:              "status normalized",
			inputAlertType:    "status",
			expectedAlertType: "SERVICE_LEVELS_OBJECTIVE",
			expectedMetric:    "STATUS",
		},
		{
			name:              "Status normalized",
			inputAlertType:    "Status",
			expectedAlertType: "SERVICE_LEVELS_OBJECTIVE",
			expectedMetric:    "STATUS",
		},
		{
			name:              "burnRateV2 normalized",
			inputAlertType:    "burnRateV2",
			expectedAlertType: "ERROR_BUDGET",
			expectedMetric:    "BURN_RATE_V2",
		},
		{
			name:              "BurnRateV2 normalized",
			inputAlertType:    "BurnRateV2",
			expectedAlertType: "ERROR_BUDGET",
			expectedMetric:    "BURN_RATE_V2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := SloAlertConfigModel{
				ID:          types.StringValue("test-id"),
				Name:        types.StringValue("Test Alert"),
				Description: types.StringValue("Test Description"),
				Severity:    types.Int64Value(5),
				Triggering:  types.BoolValue(true),

				AlertType: types.StringValue(tc.inputAlertType),
				TimeThreshold: &SloAlertTimeThresholdModel{
					WarmUp:   types.Int64Value(300000),
					CoolDown: types.Int64Value(600000),
				},
				SloIds: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("slo-1"),
				}),
				AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("channel-1"),
				}),
				CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
			}

			if tc.inputAlertType != "burn_rate_v2" && tc.inputAlertType != "burnRateV2" && tc.inputAlertType != "BurnRateV2" {
				model.Threshold = &SloAlertThresholdModel{
					Type:     types.StringValue("staticThreshold"),
					Operator: types.StringValue(">="),
					Value:    types.Float64Value(95.0),
				}
			}

			state := createMockState(t, model)

			apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, apiConfig)

			assert.Equal(t, tc.expectedAlertType, apiConfig.Rule.AlertType)
			assert.Equal(t, tc.expectedMetric, apiConfig.Rule.Metric)
		})
	}
}

func TestMapStateToDataObject_InvalidAlertType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("invalid_type"),
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.True(t, diags.HasError())
	require.Nil(t, apiConfig)
	assert.Contains(t, diags[0].Summary(), "Error mapping alert type")
}

func TestMapStateToDataObject_InvalidDuration(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("burn_rate_v2"),
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		BurnRateConfig: []SloAlertBurnRateConfigModel{
			{
				AlertWindowType:   types.StringValue("SHORT"),
				Duration:          types.StringValue("invalid"),
				DurationUnitType:  types.StringValue("MINUTES"),
				ThresholdOperator: types.StringValue(">"),
				ThresholdValue:    types.StringValue("2.5"),
			},
		},
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.True(t, diags.HasError())
	require.Nil(t, apiConfig)
	assert.Contains(t, diags[0].Summary(), "Error parsing duration")
}

func TestMapStateToDataObject_InvalidThresholdValue(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("burn_rate_v2"),
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		BurnRateConfig: []SloAlertBurnRateConfigModel{
			{
				AlertWindowType:   types.StringValue("SHORT"),
				Duration:          types.StringValue("60"),
				DurationUnitType:  types.StringValue("MINUTES"),
				ThresholdOperator: types.StringValue(">"),
				ThresholdValue:    types.StringValue("invalid"),
			},
		},
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.True(t, diags.HasError())
	require.Nil(t, apiConfig)
	assert.Contains(t, diags[0].Summary(), "Error parsing threshold value")
}

func TestMapStateToDataObject_WithNullThresholdType(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("status"),
		Threshold: &SloAlertThresholdModel{
			Type:     types.StringNull(),
			Operator: types.StringValue(">="),
			Value:    types.Float64Value(95.0),
		},
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.NotNil(t, apiConfig.Threshold)
	assert.Equal(t, "staticThreshold", apiConfig.Threshold.Type)
}

func TestMapStateToDataObject_BurnRateV2WithoutBurnRateConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("burn_rate_v2"),
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		BurnRateConfig: nil,
		CustomPayload:  types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.Equal(t, "ERROR_BUDGET", apiConfig.Rule.AlertType)
	assert.Equal(t, "BURN_RATE_V2", apiConfig.Rule.Metric)
	assert.NotNil(t, apiConfig.BurnRateConfigs)
	assert.Len(t, *apiConfig.BurnRateConfigs, 0)
}

func TestMapStateToDataObject_FromState(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("status"),
		Threshold: &SloAlertThresholdModel{
			Type:     types.StringValue("staticThreshold"),
			Operator: types.StringValue(">="),
			Value:    types.Float64Value(95.0),
		},
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.Equal(t, "test-id", apiConfig.ID)
	assert.Equal(t, "Test Alert", apiConfig.Name)
	assert.Equal(t, "SERVICE_LEVELS_OBJECTIVE", apiConfig.Rule.AlertType)
	assert.Equal(t, "STATUS", apiConfig.Rule.Metric)
}

func TestMapStateToDataObject_WithEmptyID(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	model := SloAlertConfigModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Test Alert"),
		Description: types.StringValue("Test Description"),
		Severity:    types.Int64Value(5),
		Triggering:  types.BoolValue(true),

		AlertType: types.StringValue("status"),
		Threshold: &SloAlertThresholdModel{
			Type:     types.StringValue("staticThreshold"),
			Operator: types.StringValue(">="),
			Value:    types.Float64Value(95.0),
		},
		TimeThreshold: &SloAlertTimeThresholdModel{
			WarmUp:   types.Int64Value(300000),
			CoolDown: types.Int64Value(600000),
		},
		SloIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("slo-1"),
		}),
		AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
		}),
		CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	state := createMockState(t, model)

	apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, apiConfig)

	assert.Equal(t, "", apiConfig.ID)
}

func TestMapStateToDataObject_WithDifferentOperators(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	operators := []string{">", ">=", "=", "<=", "<"}

	for _, operator := range operators {
		t.Run("operator_"+operator, func(t *testing.T) {
			model := SloAlertConfigModel{
				ID:          types.StringValue("test-id"),
				Name:        types.StringValue("Test Alert"),
				Description: types.StringValue("Test Description"),
				Severity:    types.Int64Value(5),
				Triggering:  types.BoolValue(true),

				AlertType: types.StringValue("status"),
				Threshold: &SloAlertThresholdModel{
					Type:     types.StringValue("staticThreshold"),
					Operator: types.StringValue(operator),
					Value:    types.Float64Value(95.0),
				},
				TimeThreshold: &SloAlertTimeThresholdModel{
					WarmUp:   types.Int64Value(300000),
					CoolDown: types.Int64Value(600000),
				},
				SloIds: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("slo-1"),
				}),
				AlertChannelIds: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("channel-1"),
				}),
				CustomPayload: types.ListNull(shared.GetCustomPayloadFieldType()),
			}

			state := createMockState(t, model)

			apiConfig, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, apiConfig)

			assert.Equal(t, operator, apiConfig.Threshold.Operator)
		})
	}
}

func TestSetComputedFields(t *testing.T) {
	ctx := context.Background()
	resource := NewSloAlertConfigResourceHandle()

	plan := tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, &plan)
	require.False(t, diags.HasError())
}

func TestMetaData(t *testing.T) {
	resource := NewSloAlertConfigResourceHandle()

	metaData := resource.MetaData()
	require.NotNil(t, metaData)
	assert.Equal(t, ResourceInstanaSloAlertConfig, metaData.ResourceName)
	assert.Equal(t, int64(2), metaData.SchemaVersion)
	assert.False(t, metaData.CreateOnly)
}

func TestNewSloAlertConfigResourceHandle(t *testing.T) {
	resource := NewSloAlertConfigResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	require.NotNil(t, metaData)
	assert.Equal(t, ResourceInstanaSloAlertConfig, metaData.ResourceName)
}

// Helper functions

func initializeEmptyState(ctx context.Context, state *tfsdk.State) {
	emptyModel := SloAlertConfigModel{
		ID:              types.StringNull(),
		Name:            types.StringNull(),
		Description:     types.StringNull(),
		Severity:        types.Int64Null(),
		Triggering:      types.BoolNull(),
		AlertType:       types.StringNull(),
		Threshold:       nil,
		SloIds:          types.SetNull(types.StringType),
		AlertChannelIds: types.SetNull(types.StringType),
		TimeThreshold:   nil,
		BurnRateConfig:  nil,
		CustomPayload:   types.ListNull(shared.GetCustomPayloadFieldType()),
	}
	state.Set(ctx, &emptyModel)
}

func createMockState(t *testing.T, model SloAlertConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewSloAlertConfigResourceHandle()
	return resource.MetaData().Schema
}
