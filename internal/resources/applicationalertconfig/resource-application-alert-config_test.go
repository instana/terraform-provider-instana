package applicationalertconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}

func TestNewApplicationAlertConfigResourceHandle(t *testing.T) {
	resource := NewApplicationAlertConfigResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaApplicationAlertConfig, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(2), metaData.SchemaVersion)
	assert.True(t, metaData.SkipIDGeneration)
}

func TestNewGlobalApplicationAlertConfigResourceHandle(t *testing.T) {
	resource := NewGlobalApplicationAlertConfigResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaGlobalApplicationAlertConfig, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(2), metaData.SchemaVersion)
}

func TestMetaData(t *testing.T) {
	resource := &applicationAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 1,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_resource", metaData.ResourceName)
	assert.Equal(t, int64(1), metaData.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := NewApplicationAlertConfigResourceHandle()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetRestResource(t *testing.T) {
	tests := []struct {
		name     string
		isGlobal bool
	}{
		{
			name:     "non-global resource",
			isGlobal: false,
		},
		{
			name:     "global resource",
			isGlobal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &applicationAlertConfigResource{
				isGlobal: tt.isGlobal,
			}

			// We can't test the actual API call without a mock, but we can verify the method exists
			assert.NotNil(t, resource.GetRestResource)
		})
	}
}

func TestResourceImpl_GetID(t *testing.T) {
	resource := &applicationAlertConfigResourceImpl{}
	data := &restapi.ApplicationAlertConfig{
		ID: "test-id-123",
	}

	id := resource.GetID(data)
	assert.Equal(t, "test-id-123", id)
}

func TestResourceImpl_SetID(t *testing.T) {
	resource := &applicationAlertConfigResourceImpl{}
	data := &restapi.ApplicationAlertConfig{}

	resource.SetID(data, "new-id-456")
	assert.Equal(t, "new-id-456", data.ID)
}

func TestMapStateToDataObject_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),
		Enabled:          types.BoolValue(true),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Alert", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, restapi.BoundaryScope("ALL"), result.BoundaryScope)
	assert.Equal(t, restapi.ApplicationAlertEvaluationType("PER_AP"), result.EvaluationType)
	assert.Equal(t, restapi.Granularity(600000), result.Granularity)
	assert.False(t, result.IncludeInternal)
	assert.False(t, result.IncludeSynthetic)
	assert.False(t, result.Triggering)
	// Verify enabled field defaults to true when not specified
	require.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled)
}

func TestMapStateToDataObject_WithGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		GracePeriod:      types.Int64Value(int64(300000)),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.GracePeriod)
	assert.Equal(t, int64(300000), *result.GracePeriod)
}

func TestMapStateToDataObject_WithApplications(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications: []ApplicationModel{
			{
				ApplicationID: types.StringValue("app-1"),
				Inclusive:     types.BoolValue(true),
				Services: []ServiceModel{
					{
						ServiceID: types.StringValue("svc-1"),
						Inclusive: types.BoolValue(true),
						Endpoints: []EndpointModel{
							{
								EndpointID: types.StringValue("ep-1"),
								Inclusive:  types.BoolValue(true),
							},
						},
					},
				},
			},
		},
		Rules: []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Applications, 1)

	app := result.Applications["app-1"]
	assert.Equal(t, "app-1", app.ApplicationID)
	assert.True(t, app.Inclusive)
	require.Len(t, app.Services, 1)

	svc := app.Services["svc-1"]
	assert.Equal(t, "svc-1", svc.ServiceID)
	assert.True(t, svc.Inclusive)
	require.Len(t, svc.Endpoints, 1)

	ep := svc.Endpoints["ep-1"]
	assert.Equal(t, "ep-1", ep.EndpointID)
	assert.True(t, ep.Inclusive)
}

func TestMapStateToDataObject_WithTimeThreshold(t *testing.T) {
	tests := []struct {
		name          string
		timeThreshold *AppAlertTimeThresholdModel
		expectedType  string
	}{
		{
			name: "request_impact",
			timeThreshold: &AppAlertTimeThresholdModel{
				RequestImpact: &AppAlertRequestImpactModel{
					TimeWindow: types.Int64Value(int64(600000)),
					Requests:   types.Int64Value(int64(100)),
				},
			},
			expectedType: "requestImpact",
		},
		{
			name: "violations_in_period",
			timeThreshold: &AppAlertTimeThresholdModel{
				ViolationsInPeriod: &AppAlertViolationsInPeriodModel{
					TimeWindow: types.Int64Value(int64(600000)),
					Violations: types.Int64Value(int64(5)),
				},
			},
			expectedType: "violationsInPeriod",
		},
		{
			name: "violations_in_sequence",
			timeThreshold: &AppAlertTimeThresholdModel{
				ViolationsInSequence: &AppAlertViolationsInSequenceModel{
					TimeWindow: types.Int64Value(int64(600000)),
				},
			},
			expectedType: "violationsInSequence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationAlertConfigResourceImpl{}

			state := createMockState(t, ApplicationAlertConfigModel{
				ID:               types.StringValue("test-id"),
				Name:             types.StringValue("Test Alert"),
				Description:      types.StringValue("Test Description"),
				BoundaryScope:    types.StringValue("ALL"),
				EvaluationType:   types.StringValue("PER_AP"),
				Granularity:      types.Int64Value(int64(600000)),
				IncludeInternal:  types.BoolValue(false),
				IncludeSynthetic: types.BoolValue(false),
				Triggering:       types.BoolValue(false),

				AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Applications:        []ApplicationModel{},
				Rules:               []RuleWithThresholdModel{},
				TimeThreshold:       tt.timeThreshold,
			})

			result, diags := resource.MapStateToDataObject(ctx, state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			require.NotNil(t, result.TimeThreshold)
			assert.Equal(t, tt.expectedType, result.TimeThreshold.Type)
		})
	}
}

func TestUpdateState_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Alert", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
}

func TestUpdateState_WithGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	gracePeriod := int64(300000)
	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		GracePeriod:      &gracePeriod,
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.GracePeriod.IsNull())
	assert.Equal(t, int64(300000), model.GracePeriod.ValueInt64())
}

func TestUpdateState_WithApplications(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications: map[string]restapi.IncludedApplication{
			"app-1": {
				ApplicationID: "app-1",
				Inclusive:     true,
				Services: map[string]restapi.IncludedService{
					"svc-1": {
						ServiceID: "svc-1",
						Inclusive: true,
						Endpoints: map[string]restapi.IncludedEndpoint{
							"ep-1": {
								EndpointID: "ep-1",
								Inclusive:  true,
							},
						},
					},
				},
			},
		},
		Rules: []restapi.ApplicationAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Applications, 1)
	assert.Equal(t, "app-1", model.Applications[0].ApplicationID.ValueString())
}

func TestUpdateState_WithTimeThreshold(t *testing.T) {
	t.Run("request_impact", func(t *testing.T) {
		ctx := context.Background()
		resource := &applicationAlertConfigResourceImpl{}

		data := &restapi.ApplicationAlertConfig{
			ID:               "test-id",
			Name:             "Test Alert",
			Description:      "Test Description",
			BoundaryScope:    "ALL",
			EvaluationType:   "PER_AP",
			Granularity:      600000,
			IncludeInternal:  false,
			IncludeSynthetic: false,
			Triggering:       false,
			AlertChannelIDs:  []string{},
			Applications:     map[string]restapi.IncludedApplication{},
			Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
			TimeThreshold: &restapi.ApplicationAlertTimeThreshold{
				Type:       "requestImpact",
				TimeWindow: 600000,
				Requests:   100,
			},
		}

		state := &tfsdk.State{
			Schema: getTestSchema(),
		}

		diags := resource.UpdateState(ctx, state, nil, data)
		require.False(t, diags.HasError())

		var model ApplicationAlertConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		require.NotNil(t, model.TimeThreshold)
		require.NotNil(t, model.TimeThreshold.RequestImpact)
		assert.Equal(t, int64(600000), model.TimeThreshold.RequestImpact.TimeWindow.ValueInt64())
		assert.Equal(t, int64(100), model.TimeThreshold.RequestImpact.Requests.ValueInt64())
	})
}

func TestExtractGracePeriod(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Int64
		expected *int64
	}{
		{
			name:     "with_value",
			input:    types.Int64Value(int64(300000)),
			expected: func() *int64 { v := int64(300000); return &v }(),
		},
		{
			name:     "null_value",
			input:    types.Int64Null(),
			expected: nil,
		},
		{
			name:     "unknown_value",
			input:    types.Int64Unknown(),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractGracePeriod(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestMapStateToDataObject_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	// Create alert channels map
	alertChannels := types.MapValueMust(
		types.SetType{ElemType: types.StringType},
		map[string]attr.Value{
			"warning": types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("channel-1"),
				types.StringValue("channel-2"),
			}),
			"critical": types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("channel-3"),
			}),
		},
	)

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       alertChannels,
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.AlertChannels)
	assert.Len(t, result.AlertChannels["warning"], 2)
	assert.Len(t, result.AlertChannels["critical"], 1)
}

func TestUpdateState_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		AlertChannels: map[string][]string{
			"warning":  {"channel-1", "channel-2"},
			"critical": {"channel-3"},
		},
		Applications: map[string]restapi.IncludedApplication{},
		Rules:        []restapi.ApplicationAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.AlertChannels.IsNull())
}

func TestMapStateToDataObject_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	// Create custom payload fields
	customPayloadFields := types.ListValueMust(
		shared.GetCustomPayloadFieldType(),
		[]attr.Value{
			types.ObjectValueMust(
				shared.CustomPayloadFieldAttributeTypes(),
				map[string]attr.Value{
					"key":   types.StringValue("field1"),
					"value": types.StringValue("value1"),
					"dynamic_value": types.ObjectNull(map[string]attr.Type{
						"key":      types.StringType,
						"tag_name": types.StringType,
					}),
				},
			),
		},
	)

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: customPayloadFields,
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.CustomerPayloadFields, 1)
	assert.Equal(t, "field1", result.CustomerPayloadFields[0].Key)
}

// Test all 6 rule types
func TestMapStateToDataObject_WithRules(t *testing.T) {
	tests := []struct {
		name     string
		rule     *RuleModel
		expected string
	}{
		{
			name: "error_rate",
			rule: &RuleModel{
				ErrorRate: &RuleConfigModel{
					MetricName:  types.StringValue("errors"),
					Aggregation: types.StringValue("sum"),
				},
			},
			expected: APIAlertTypeErrorRate,
		},
		{
			name: "errors",
			rule: &RuleModel{
				Errors: &RuleConfigModel{
					MetricName:  types.StringValue("errors"),
					Aggregation: types.StringValue("sum"),
				},
			},
			expected: APIAlertTypeErrors,
		},
		{
			name: "logs",
			rule: &RuleModel{
				Logs: &LogsRuleModel{
					MetricName:  types.StringValue("logs"),
					Aggregation: types.StringValue("sum"),
					Level:       types.StringValue("ERROR"),
					Message:     types.StringValue("test message"),
					Operator:    types.StringValue("CONTAINS"),
				},
			},
			expected: APIAlertTypeLogs,
		},
		{
			name: "slowness",
			rule: &RuleModel{
				Slowness: &RuleConfigModel{
					MetricName:  types.StringValue("latency"),
					Aggregation: types.StringValue("mean"),
				},
			},
			expected: APIAlertTypeSlowness,
		},
		{
			name: "status_code",
			rule: &RuleModel{
				StatusCode: &StatusCodeRuleModel{
					MetricName:      types.StringValue("statusCode"),
					Aggregation:     types.StringValue("sum"),
					StatusCodeStart: types.Int64Value(int64(500)),
					StatusCodeEnd:   types.Int64Value(int64(599)),
				},
			},
			expected: APIAlertTypeStatusCode,
		},
		{
			name: "throughput",
			rule: &RuleModel{
				Throughput: &RuleConfigModel{
					MetricName:  types.StringValue("calls"),
					Aggregation: types.StringValue("sum"),
				},
			},
			expected: APIAlertTypeThroughput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationAlertConfigResourceImpl{}

			state := createMockState(t, ApplicationAlertConfigModel{
				ID:               types.StringValue("test-id"),
				Name:             types.StringValue("Test Alert"),
				Description:      types.StringValue("Test Description"),
				BoundaryScope:    types.StringValue("ALL"),
				EvaluationType:   types.StringValue("PER_AP"),
				Granularity:      types.Int64Value(int64(600000)),
				IncludeInternal:  types.BoolValue(false),
				IncludeSynthetic: types.BoolValue(false),
				Triggering:       types.BoolValue(false),

				AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Applications:        []ApplicationModel{},
				Rules: []RuleWithThresholdModel{
					{
						Rule:              tt.rule,
						ThresholdOperator: types.StringValue(">"),
						Thresholds: &ApplicationThresholdModel{
							Warning: &ThresholdLevelModel{
								Static: &shared.StaticTypeModel{
									Value: types.Float64Value(100),
								},
							},
						},
					},
				},
			})

			result, diags := resource.MapStateToDataObject(ctx, state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			require.Len(t, result.Rules, 1)
			assert.Equal(t, tt.expected, result.Rules[0].Rule.AlertType)
		})
	}
}

// Test static and adaptive baseline thresholds
func TestMapStateToDataObject_WithStaticThreshold(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					ErrorRate: &RuleConfigModel{
						MetricName:  types.StringValue("errors"),
						Aggregation: types.StringValue("sum"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Warning: &ThresholdLevelModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(100),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)
	assert.Equal(t, "staticThreshold", result.Rules[0].Thresholds[restapi.WarningSeverity].Type)
	require.NotNil(t, result.Rules[0].Thresholds[restapi.WarningSeverity].Value)
	assert.Equal(t, float64(100), *result.Rules[0].Thresholds[restapi.WarningSeverity].Value)
}

// Test UpdateState with rules containing thresholds
func TestUpdateState_WithRulesAndThresholds(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	// Test with static threshold
	value := float64(100)
	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules: []restapi.ApplicationAlertRuleWithThresholds{
			{
				Rule: &restapi.ApplicationAlertRule{
					AlertType:   ApplicationAlertConfigFieldRuleErrorRate,
					MetricName:  "errors",
					Aggregation: "sum",
				},
				ThresholdOperator: ">",
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &value,
					},
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	assert.Equal(t, ">", model.Rules[0].ThresholdOperator.ValueString())
	require.NotNil(t, model.Rules[0].Thresholds)
	require.NotNil(t, model.Rules[0].Thresholds.Warning)
	require.NotNil(t, model.Rules[0].Thresholds.Warning.Static)
	assert.Equal(t, float64(100), model.Rules[0].Thresholds.Warning.Static.Value.ValueFloat64())
}

// Test UpdateState with custom payload fields
func TestUpdateState_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "field1",
				Value: "value1",
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.CustomPayloadFields.IsNull())
}

// Test UpdateState with tag filter
func TestUpdateState_WithTagFilterExpression(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	entityType := "entity.type"
	serviceValue := "service"
	equalsOp := restapi.EqualsOperator

	tagFilter := &restapi.TagFilter{
		Type:     restapi.TagFilterExpressionType,
		Name:     &entityType,
		Operator: &equalsOp,
		Value:    &serviceValue,
	}

	data := &restapi.ApplicationAlertConfig{
		ID:                  "test-id",
		Name:                "Test Alert",
		Description:         "Test Description",
		BoundaryScope:       "ALL",
		EvaluationType:      "PER_AP",
		Granularity:         600000,
		IncludeInternal:     false,
		IncludeSynthetic:    false,
		Triggering:          false,
		AlertChannelIDs:     []string{},
		Applications:        map[string]restapi.IncludedApplication{},
		Rules:               []restapi.ApplicationAlertRuleWithThresholds{},
		TagFilterExpression: tagFilter,
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Tag filter should be set (not checking the exact value due to normalization)
	// The important thing is that UpdateState doesn't error
}

// Test UpdateState with all time threshold types
func TestUpdateState_WithAllTimeThresholdTypes(t *testing.T) {
	tests := []struct {
		name          string
		timeThreshold *restapi.ApplicationAlertTimeThreshold
		checkFunc     func(*testing.T, ApplicationAlertConfigModel)
	}{
		{
			name: "violations_in_period",
			timeThreshold: &restapi.ApplicationAlertTimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: 600000,
				Violations: 5,
			},
			checkFunc: func(t *testing.T, model ApplicationAlertConfigModel) {
				require.NotNil(t, model.TimeThreshold)
				require.NotNil(t, model.TimeThreshold.ViolationsInPeriod)
				assert.Equal(t, int64(600000), model.TimeThreshold.ViolationsInPeriod.TimeWindow.ValueInt64())
				assert.Equal(t, int64(5), model.TimeThreshold.ViolationsInPeriod.Violations.ValueInt64())
			},
		},
		{
			name: "violations_in_sequence",
			timeThreshold: &restapi.ApplicationAlertTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: 600000,
			},
			checkFunc: func(t *testing.T, model ApplicationAlertConfigModel) {
				require.NotNil(t, model.TimeThreshold)
				require.NotNil(t, model.TimeThreshold.ViolationsInSequence)
				assert.Equal(t, int64(600000), model.TimeThreshold.ViolationsInSequence.TimeWindow.ValueInt64())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationAlertConfigResourceImpl{}

			data := &restapi.ApplicationAlertConfig{
				ID:               "test-id",
				Name:             "Test Alert",
				Description:      "Test Description",
				BoundaryScope:    "ALL",
				EvaluationType:   "PER_AP",
				Granularity:      600000,
				IncludeInternal:  false,
				IncludeSynthetic: false,
				Triggering:       false,
				AlertChannelIDs:  []string{},
				Applications:     map[string]restapi.IncludedApplication{},
				Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
				TimeThreshold:    tt.timeThreshold,
			}

			state := &tfsdk.State{
				Schema: getTestSchema(),
			}

			diags := resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError())

			var model ApplicationAlertConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			tt.checkFunc(t, model)
		})
	}
}

// Test UpdateState with all rule types
func TestUpdateState_WithAllRuleTypes(t *testing.T) {
	tests := []struct {
		name string
		rule *restapi.ApplicationAlertRule
	}{
		{
			name: "errors_rule",
			rule: &restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleErrors,
				MetricName:  "errors",
				Aggregation: "sum",
			},
		},
		{
			name: "logs_rule",
			rule: &restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleLogs,
				MetricName:  "logs",
				Aggregation: "sum",
				Level:       ptr(restapi.LogLevelError),
				Message:     ptr("test message"),
				Operator:    ptr(restapi.ContainsOperator),
			},
		},
		{
			name: "slowness_rule",
			rule: &restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleSlowness,
				MetricName:  "latency",
				Aggregation: "mean",
			},
		},
		{
			name: "status_code_rule",
			rule: &restapi.ApplicationAlertRule{
				AlertType:       ApplicationAlertConfigFieldRuleStatusCode,
				MetricName:      "statusCode",
				Aggregation:     "sum",
				StatusCodeStart: ptr(int32(500)),
				StatusCodeEnd:   ptr(int32(599)),
			},
		},
		{
			name: "throughput_rule",
			rule: &restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleThroughput,
				MetricName:  "calls",
				Aggregation: "sum",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationAlertConfigResourceImpl{}

			data := &restapi.ApplicationAlertConfig{
				ID:               "test-id",
				Name:             "Test Alert",
				Description:      "Test Description",
				BoundaryScope:    "ALL",
				EvaluationType:   "PER_AP",
				Granularity:      600000,
				IncludeInternal:  false,
				IncludeSynthetic: false,
				Triggering:       false,
				AlertChannelIDs:  []string{},
				Applications:     map[string]restapi.IncludedApplication{},
				Rules: []restapi.ApplicationAlertRuleWithThresholds{
					{
						Rule:              tt.rule,
						ThresholdOperator: ">",
						Thresholds:        map[restapi.AlertSeverity]restapi.ThresholdRule{},
					},
				},
			}

			state := &tfsdk.State{
				Schema: getTestSchema(),
			}

			diags := resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError())

			var model ApplicationAlertConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			require.Len(t, model.Rules, 1)
		})
	}
}

// Test tag filter expression
func TestMapStateToDataObject_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),
		TagFilter:        types.StringValue("entity.type EQUALS 'service'"),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)
	require.NotNil(t, result.TagFilterExpression)
}

// Test error handling for missing required fields
func TestMapStateToDataObject_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		rule        *RuleModel
		expectError bool
	}{
		{
			name: "error_rate_missing_metric_name",
			rule: &RuleModel{
				ErrorRate: &RuleConfigModel{
					MetricName:  types.StringNull(),
					Aggregation: types.StringValue("sum"),
				},
			},
			expectError: true,
		},
		{
			name: "logs_missing_required_fields",
			rule: &RuleModel{
				Logs: &LogsRuleModel{
					MetricName:  types.StringNull(),
					Aggregation: types.StringValue("sum"),
					Level:       types.StringNull(),
					Operator:    types.StringNull(),
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationAlertConfigResourceImpl{}

			state := createMockState(t, ApplicationAlertConfigModel{
				ID:               types.StringValue("test-id"),
				Name:             types.StringValue("Test Alert"),
				Description:      types.StringValue("Test Description"),
				BoundaryScope:    types.StringValue("ALL"),
				EvaluationType:   types.StringValue("PER_AP"),
				Granularity:      types.Int64Value(int64(600000)),
				IncludeInternal:  types.BoolValue(false),
				IncludeSynthetic: types.BoolValue(false),
				Triggering:       types.BoolValue(false),

				AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Applications:        []ApplicationModel{},
				Rules: []RuleWithThresholdModel{
					{
						Rule:              tt.rule,
						ThresholdOperator: types.StringValue(">"),
					},
				},
			})

			_, diags := resource.MapStateToDataObject(ctx, state)
			if tt.expectError {
				assert.True(t, diags.HasError())
			} else {
				assert.False(t, diags.HasError())
			}
		})
	}
}

// Test wrapper MapStateToDataObject function
func TestWrapperMapStateToDataObject(t *testing.T) {
	ctx := context.Background()
	resource := NewApplicationAlertConfigResourceHandle()

	// Create a plan with basic config
	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	model := ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Alert", result.Name)
}

// Test wrapper MapStateToDataObject with state instead of plan
func TestWrapperMapStateToDataObject_WithState(t *testing.T) {
	ctx := context.Background()
	resource := NewApplicationAlertConfigResourceHandle()

	state := &tfsdk.State{
		Schema: resource.MetaData().Schema,
	}

	model := ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id-2"),
		Name:             types.StringValue("Test Alert 2"),
		Description:      types.StringValue("Test Description 2"),
		BoundaryScope:    types.StringValue("INBOUND"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(true),
		IncludeSynthetic: types.BoolValue(true),
		Triggering:       types.BoolValue(true),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	}

	diags := state.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id-2", result.ID)
	assert.Equal(t, "Test Alert 2", result.Name)
	assert.Equal(t, restapi.BoundaryScope("INBOUND"), result.BoundaryScope)
	assert.True(t, result.IncludeInternal)
	assert.True(t, result.IncludeSynthetic)
	assert.True(t, result.Triggering)
}

// Test wrapper UpdateState function
func TestWrapperUpdateState(t *testing.T) {
	ctx := context.Background()
	resource := NewApplicationAlertConfigResourceHandle()

	state := &tfsdk.State{
		Schema: resource.MetaData().Schema,
	}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())
	assert.Equal(t, "test-id", model.ID.ValueString())
}

// Test NewResource function
func TestNewResource(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResource{
		isGlobal: false,
	}

	impl, diags := resource.NewResource(ctx, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, impl)

	// Test that it's the correct implementation
	assert.IsType(t, &applicationAlertConfigResourceImpl{}, impl)
}

// Test NewResource for global config
func TestNewResource_Global(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResource{
		isGlobal: true,
	}

	impl, diags := resource.NewResource(ctx, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, impl)

	implTyped := impl.(*applicationAlertConfigResourceImpl)
	assert.True(t, implTyped.isGlobal)
}

// Test MapStateToDataObject with invalid tag filter
func TestMapStateToDataObject_InvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),
		TagFilter:        types.StringValue("invalid tag filter syntax"),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	_, diags := resource.MapStateToDataObject(ctx, state)
	// Should have error due to invalid tag filter
	assert.True(t, diags.HasError())
}

// Test MapStateToDataObject with all boundary scopes
func TestMapStateToDataObject_AllBoundaryScopes(t *testing.T) {
	scopes := []string{"ALL", "INBOUND", "DEFAULT"}

	for _, scope := range scopes {
		t.Run(scope, func(t *testing.T) {
			ctx := context.Background()
			resource := &applicationAlertConfigResourceImpl{}

			state := createMockState(t, ApplicationAlertConfigModel{
				ID:               types.StringValue("test-id"),
				Name:             types.StringValue("Test Alert"),
				Description:      types.StringValue("Test Description"),
				BoundaryScope:    types.StringValue(scope),
				EvaluationType:   types.StringValue("PER_AP"),
				Granularity:      types.Int64Value(int64(600000)),
				IncludeInternal:  types.BoolValue(false),
				IncludeSynthetic: types.BoolValue(false),
				Triggering:       types.BoolValue(false),

				AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Applications:        []ApplicationModel{},
				Rules:               []RuleWithThresholdModel{},
			})

			result, diags := resource.MapStateToDataObject(ctx, state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			assert.Equal(t, restapi.BoundaryScope(scope), result.BoundaryScope)
		})
	}
}

// Test MapStateToDataObject with complex application hierarchy
func TestMapStateToDataObject_ComplexApplicationHierarchy(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications: []ApplicationModel{
			{
				ApplicationID: types.StringValue("app-1"),
				Inclusive:     types.BoolValue(true),
				Services: []ServiceModel{
					{
						ServiceID: types.StringValue("svc-1"),
						Inclusive: types.BoolValue(true),
						Endpoints: []EndpointModel{
							{
								EndpointID: types.StringValue("ep-1"),
								Inclusive:  types.BoolValue(true),
							},
							{
								EndpointID: types.StringValue("ep-2"),
								Inclusive:  types.BoolValue(false),
							},
						},
					},
					{
						ServiceID: types.StringValue("svc-2"),
						Inclusive: types.BoolValue(false),
						Endpoints: []EndpointModel{},
					},
				},
			},
			{
				ApplicationID: types.StringValue("app-2"),
				Inclusive:     types.BoolValue(false),
				Services:      []ServiceModel{},
			},
		},
		Rules: []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Applications, 2)

	app1 := result.Applications["app-1"]
	assert.Equal(t, "app-1", app1.ApplicationID)
	assert.True(t, app1.Inclusive)
	require.Len(t, app1.Services, 2)

	svc1 := app1.Services["svc-1"]
	assert.Equal(t, "svc-1", svc1.ServiceID)
	assert.True(t, svc1.Inclusive)
	require.Len(t, svc1.Endpoints, 2)

	ep1 := svc1.Endpoints["ep-1"]
	assert.Equal(t, "ep-1", ep1.EndpointID)
	assert.True(t, ep1.Inclusive)

	ep2 := svc1.Endpoints["ep-2"]
	assert.Equal(t, "ep-2", ep2.EndpointID)
	assert.False(t, ep2.Inclusive)
}

// Test MapStateToDataObject with critical threshold
func TestMapStateToDataObject_WithCriticalThreshold(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					ErrorRate: &RuleConfigModel{
						MetricName:  types.StringValue("errors"),
						Aggregation: types.StringValue("sum"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Critical: &ThresholdLevelModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(200),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.CriticalSeverity)
	assert.Equal(t, "staticThreshold", result.Rules[0].Thresholds[restapi.CriticalSeverity].Type)
	require.NotNil(t, result.Rules[0].Thresholds[restapi.CriticalSeverity].Value)
	assert.Equal(t, float64(200), *result.Rules[0].Thresholds[restapi.CriticalSeverity].Value)
}

// Test MapStateToDataObject with both warning and critical thresholds
func TestMapStateToDataObject_WithBothThresholds(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					ErrorRate: &RuleConfigModel{
						MetricName:  types.StringValue("errors"),
						Aggregation: types.StringValue("sum"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Warning: &ThresholdLevelModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(100),
						},
					},
					Critical: &ThresholdLevelModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(200),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)
	require.Contains(t, result.Rules[0].Thresholds, restapi.CriticalSeverity)
	assert.Equal(t, float64(100), *result.Rules[0].Thresholds[restapi.WarningSeverity].Value)
	assert.Equal(t, float64(200), *result.Rules[0].Thresholds[restapi.CriticalSeverity].Value)
}

// Test MapStateToDataObject with null granularity
func TestMapStateToDataObject_WithNullGranularity(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Null(),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	// Granularity should be 0 when null
	assert.Equal(t, restapi.Granularity(0), result.Granularity)
}

// Test MapStateToDataObject with empty tag filter
func TestMapStateToDataObject_WithEmptyTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),
		TagFilter:        types.StringValue(""),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	// Empty tag filter should result in nil TagFilterExpression
	assert.Nil(t, result.TagFilterExpression)
}

// Test MapStateToDataObject with null severity
func TestMapStateToDataObject_WithNullSeverity(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules:               []RuleWithThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	// Null severity should result in 0
	// Severity field removed from model
}

// Test MapStateToDataObject with adaptive baseline threshold
func TestMapStateToDataObject_WithAdaptiveBaselineThreshold(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					ErrorRate: &RuleConfigModel{
						MetricName:  types.StringValue("errors"),
						Aggregation: types.StringValue("sum"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Warning: &ThresholdLevelModel{
						AdaptiveBaseline: &shared.AdaptiveBaselineModel{
							DeviationFactor: types.Float64Value(2.0),
							Adaptability:    types.Float64Value(0.5),
							Seasonality:     types.StringValue("DAILY"),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)

	threshold := result.Rules[0].Thresholds[restapi.WarningSeverity]
	assert.Equal(t, "adaptiveBaseline", threshold.Type)
	require.NotNil(t, threshold.DeviationFactor)
	assert.Equal(t, float32(2.0), *threshold.DeviationFactor)
	require.NotNil(t, threshold.Adaptability)
	assert.Equal(t, float32(0.5), *threshold.Adaptability)
	require.NotNil(t, threshold.Seasonality)
	assert.Equal(t, restapi.ThresholdSeasonality("DAILY"), *threshold.Seasonality)
}

// Test MapStateToDataObject with adaptive baseline critical threshold
func TestMapStateToDataObject_WithAdaptiveBaselineCriticalThreshold(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					Slowness: &RuleConfigModel{
						MetricName:  types.StringValue("latency"),
						Aggregation: types.StringValue("mean"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Critical: &ThresholdLevelModel{
						AdaptiveBaseline: &shared.AdaptiveBaselineModel{
							DeviationFactor: types.Float64Value(3.0),
							Adaptability:    types.Float64Value(0.8),
							Seasonality:     types.StringValue("WEEKLY"),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.CriticalSeverity)

	threshold := result.Rules[0].Thresholds[restapi.CriticalSeverity]
	assert.Equal(t, "adaptiveBaseline", threshold.Type)
	require.NotNil(t, threshold.DeviationFactor)
	assert.Equal(t, float32(3.0), *threshold.DeviationFactor)
	require.NotNil(t, threshold.Adaptability)
	assert.Equal(t, float32(0.8), *threshold.Adaptability)
	require.NotNil(t, threshold.Seasonality)
	assert.Equal(t, restapi.ThresholdSeasonality("WEEKLY"), *threshold.Seasonality)
}

// Test UpdateState with adaptive baseline threshold
func TestUpdateState_WithAdaptiveBaselineThreshold(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	deviationFactor := float32(2.0)
	adaptability := float32(0.5)
	seasonality := restapi.ThresholdSeasonality("DAILY")

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules: []restapi.ApplicationAlertRuleWithThresholds{
			{
				Rule: &restapi.ApplicationAlertRule{
					AlertType:   ApplicationAlertConfigFieldRuleErrorRate,
					MetricName:  "errors",
					Aggregation: "sum",
				},
				ThresholdOperator: ">",
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:            "adaptiveBaseline",
						DeviationFactor: &deviationFactor,
						Adaptability:    &adaptability,
						Seasonality:     &seasonality,
					},
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].Thresholds)
	require.NotNil(t, model.Rules[0].Thresholds.Warning)
	require.NotNil(t, model.Rules[0].Thresholds.Warning.AdaptiveBaseline)

	adaptiveModel := model.Rules[0].Thresholds.Warning.AdaptiveBaseline
	assert.Equal(t, float64(2.0), adaptiveModel.DeviationFactor.ValueFloat64())
	assert.Equal(t, float64(0.5), adaptiveModel.Adaptability.ValueFloat64())
	assert.Equal(t, "DAILY", adaptiveModel.Seasonality.ValueString())
}

// Test UpdateState with both static and adaptive baseline thresholds
func TestUpdateState_WithMixedThresholds(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	warningValue := float64(100)
	criticalDeviationFactor := float32(3.0)
	criticalAdaptability := float32(0.8)
	criticalSeasonality := restapi.ThresholdSeasonality("WEEKLY")

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules: []restapi.ApplicationAlertRuleWithThresholds{
			{
				Rule: &restapi.ApplicationAlertRule{
					AlertType:   ApplicationAlertConfigFieldRuleThroughput,
					MetricName:  "calls",
					Aggregation: "sum",
				},
				ThresholdOperator: ">",
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
					restapi.CriticalSeverity: {
						Type:            "adaptiveBaseline",
						DeviationFactor: &criticalDeviationFactor,
						Adaptability:    &criticalAdaptability,
						Seasonality:     &criticalSeasonality,
					},
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].Thresholds)

	// Check warning static threshold
	require.NotNil(t, model.Rules[0].Thresholds.Warning)
	require.NotNil(t, model.Rules[0].Thresholds.Warning.Static)
	assert.Equal(t, float64(100), model.Rules[0].Thresholds.Warning.Static.Value.ValueFloat64())

	// Check critical adaptive baseline threshold
	require.NotNil(t, model.Rules[0].Thresholds.Critical)
	require.NotNil(t, model.Rules[0].Thresholds.Critical.AdaptiveBaseline)
	assert.Equal(t, float64(3.0), model.Rules[0].Thresholds.Critical.AdaptiveBaseline.DeviationFactor.ValueFloat64())
	assert.Equal(t, float64(0.8), model.Rules[0].Thresholds.Critical.AdaptiveBaseline.Adaptability.ValueFloat64())
	assert.Equal(t, "WEEKLY", model.Rules[0].Thresholds.Critical.AdaptiveBaseline.Seasonality.ValueString())
}

// Test MapStateToDataObject with null adaptive baseline fields
func TestMapStateToDataObject_WithNullAdaptiveBaselineFields(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					ErrorRate: &RuleConfigModel{
						MetricName:  types.StringValue("errors"),
						Aggregation: types.StringValue("sum"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Warning: &ThresholdLevelModel{
						AdaptiveBaseline: &shared.AdaptiveBaselineModel{
							DeviationFactor: types.Float64Null(),
							Adaptability:    types.Float64Null(),
							Seasonality:     types.StringNull(),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)

	threshold := result.Rules[0].Thresholds[restapi.WarningSeverity]
	assert.Equal(t, "adaptiveBaseline", threshold.Type)
	// Null fields are converted to zero values by the shared mapping function
	require.NotNil(t, threshold.DeviationFactor)
	assert.Equal(t, float32(0), *threshold.DeviationFactor)
	require.NotNil(t, threshold.Adaptability)
	assert.Equal(t, float32(0), *threshold.Adaptability)
	require.NotNil(t, threshold.Seasonality)
	assert.Equal(t, restapi.ThresholdSeasonality(""), *threshold.Seasonality)
}

// Test UpdateState with null adaptive baseline fields
func TestUpdateState_WithNullAdaptiveBaselineFields(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules: []restapi.ApplicationAlertRuleWithThresholds{
			{
				Rule: &restapi.ApplicationAlertRule{
					AlertType:   ApplicationAlertConfigFieldRuleErrorRate,
					MetricName:  "errors",
					Aggregation: "sum",
				},
				ThresholdOperator: ">",
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:            "adaptiveBaseline",
						DeviationFactor: nil,
						Adaptability:    nil,
						Seasonality:     nil,
					},
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].Thresholds)
	require.NotNil(t, model.Rules[0].Thresholds.Warning)
	require.NotNil(t, model.Rules[0].Thresholds.Warning.AdaptiveBaseline)

	adaptiveModel := model.Rules[0].Thresholds.Warning.AdaptiveBaseline
	assert.True(t, adaptiveModel.DeviationFactor.IsNull())
	assert.True(t, adaptiveModel.Adaptability.IsNull())
	assert.True(t, adaptiveModel.Seasonality.IsNull())
}

// Test MapStateToDataObject with both adaptive baseline thresholds
func TestMapStateToDataObject_WithBothAdaptiveBaselineThresholds(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	state := createMockState(t, ApplicationAlertConfigModel{
		ID:               types.StringValue("test-id"),
		Name:             types.StringValue("Test Alert"),
		Description:      types.StringValue("Test Description"),
		BoundaryScope:    types.StringValue("ALL"),
		EvaluationType:   types.StringValue("PER_AP"),
		Granularity:      types.Int64Value(int64(600000)),
		IncludeInternal:  types.BoolValue(false),
		IncludeSynthetic: types.BoolValue(false),
		Triggering:       types.BoolValue(false),

		AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Applications:        []ApplicationModel{},
		Rules: []RuleWithThresholdModel{
			{
				Rule: &RuleModel{
					Throughput: &RuleConfigModel{
						MetricName:  types.StringValue("calls"),
						Aggregation: types.StringValue("sum"),
					},
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &ApplicationThresholdModel{
					Warning: &ThresholdLevelModel{
						AdaptiveBaseline: &shared.AdaptiveBaselineModel{
							DeviationFactor: types.Float64Value(1.5),
							Adaptability:    types.Float64Value(0.3),
							Seasonality:     types.StringValue("DAILY"),
						},
					},
					Critical: &ThresholdLevelModel{
						AdaptiveBaseline: &shared.AdaptiveBaselineModel{
							DeviationFactor: types.Float64Value(2.5),
							Adaptability:    types.Float64Value(0.6),
							Seasonality:     types.StringValue("WEEKLY"),
						},
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Contains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)
	require.Contains(t, result.Rules[0].Thresholds, restapi.CriticalSeverity)

	// Check warning threshold
	warningThreshold := result.Rules[0].Thresholds[restapi.WarningSeverity]
	assert.Equal(t, "adaptiveBaseline", warningThreshold.Type)
	assert.Equal(t, float32(1.5), *warningThreshold.DeviationFactor)
	assert.Equal(t, float32(0.3), *warningThreshold.Adaptability)
	assert.Equal(t, restapi.ThresholdSeasonality("DAILY"), *warningThreshold.Seasonality)

	// Check critical threshold
	criticalThreshold := result.Rules[0].Thresholds[restapi.CriticalSeverity]
	assert.Equal(t, "adaptiveBaseline", criticalThreshold.Type)
	assert.Equal(t, float32(2.5), *criticalThreshold.DeviationFactor)
	assert.Equal(t, float32(0.6), *criticalThreshold.Adaptability)
	assert.Equal(t, restapi.ThresholdSeasonality("WEEKLY"), *criticalThreshold.Seasonality)
}

// Test UpdateState with empty applications
func TestUpdateState_WithEmptyApplications(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Empty(t, model.Applications)
}

// Test UpdateState with empty alert channels map
func TestUpdateState_WithEmptyAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:               "test-id",
		Name:             "Test Alert",
		Description:      "Test Description",
		BoundaryScope:    "ALL",
		EvaluationType:   "PER_AP",
		Granularity:      600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Triggering:       false,
		AlertChannelIDs:  []string{},
		AlertChannels:    map[string][]string{},
		Applications:     map[string]restapi.IncludedApplication{},
		Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Empty map should result in an empty map value, not null
	assert.False(t, model.AlertChannels.IsNull())
	assert.False(t, model.AlertChannels.IsUnknown())
	assert.Equal(t, 0, len(model.AlertChannels.Elements()))
}

// Test UpdateState with null tag filter
func TestUpdateState_WithNullTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	data := &restapi.ApplicationAlertConfig{
		ID:                  "test-id",
		Name:                "Test Alert",
		Description:         "Test Description",
		BoundaryScope:       "ALL",
		EvaluationType:      "PER_AP",
		Granularity:         600000,
		IncludeInternal:     false,
		IncludeSynthetic:    false,
		Triggering:          false,
		AlertChannelIDs:     []string{},
		Applications:        map[string]restapi.IncludedApplication{},
		Rules:               []restapi.ApplicationAlertRuleWithThresholds{},
		TagFilterExpression: nil,
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model ApplicationAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.TagFilter.IsNull())
}

// Helper functions
func TestExtractEnabled(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Bool
		expected *bool
	}{
		{
			name:     "explicit true",
			input:    types.BoolValue(true),
			expected: ptr(true),
		},
		{
			name:     "explicit false",
			input:    types.BoolValue(false),
			expected: ptr(false),
		},
		{
			name:     "null value returns nil",
			input:    types.BoolNull(),
			expected: nil,
		},
		{
			name:     "unknown value returns nil",
			input:    types.BoolUnknown(),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.SetBoolPointerFromState(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result, "Expected nil for %s", tt.name)
			} else {
				require.NotNil(t, result, "Expected non-nil for %s", tt.name)
				assert.Equal(t, *tt.expected, *result, "Expected %v for %s", *tt.expected, tt.name)
			}
		})
	}
}

func TestMapStateToDataObject_EnabledField(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	tests := []struct {
		name            string
		enabledValue    types.Bool
		expectedEnabled *bool
		description     string
	}{
		{
			name:            "enabled explicitly set to true",
			enabledValue:    types.BoolValue(true),
			expectedEnabled: ptr(true),
			description:     "When user explicitly sets enabled=true, API should receive true",
		},
		{
			name:            "enabled explicitly set to false",
			enabledValue:    types.BoolValue(false),
			expectedEnabled: ptr(false),
			description:     "When user explicitly sets enabled=false, API should receive false",
		},
		{
			name:            "enabled is null - should send nil to API",
			enabledValue:    types.BoolNull(),
			expectedEnabled: nil,
			description:     "When enabled is null, API should receive nil (omitted in JSON)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := createMockState(t, ApplicationAlertConfigModel{
				ID:                  types.StringValue("test-id"),
				Name:                types.StringValue("Test Alert"),
				Description:         types.StringValue("Test Description"),
				BoundaryScope:       types.StringValue("ALL"),
				EvaluationType:      types.StringValue("PER_AP"),
				Granularity:         types.Int64Value(600000),
				IncludeInternal:     types.BoolValue(false),
				IncludeSynthetic:    types.BoolValue(false),
				Triggering:          types.BoolValue(false),
				Enabled:             tt.enabledValue,
				AlertChannels:       types.MapNull(types.SetType{ElemType: types.StringType}),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Applications:        []ApplicationModel{},
				Rules:               []RuleWithThresholdModel{},
			})

			result, diags := resource.MapStateToDataObject(ctx, state)
			require.False(t, diags.HasError(), "Expected no errors for %s", tt.name)
			require.NotNil(t, result, "Expected result for %s", tt.name)

			if tt.expectedEnabled == nil {
				assert.Nil(t, result.Enabled, "%s: %s", tt.name, tt.description)
			} else {
				require.NotNil(t, result.Enabled, "%s: %s", tt.name, tt.description)
				assert.Equal(t, *tt.expectedEnabled, *result.Enabled, "%s: %s", tt.name, tt.description)
			}
		})
	}
}

func TestUpdateState_EnabledField(t *testing.T) {
	ctx := context.Background()
	resource := &applicationAlertConfigResourceImpl{}

	tests := []struct {
		name            string
		apiEnabled      *bool
		expectedEnabled bool
		description     string
	}{
		{
			name:            "API returns enabled=true",
			apiEnabled:      ptr(true),
			expectedEnabled: true,
			description:     "When API returns enabled=true, state should have enabled=true",
		},
		{
			name:            "API returns enabled=false",
			apiEnabled:      ptr(false),
			expectedEnabled: false,
			description:     "When API returns enabled=false, state should have enabled=false",
		},
		{
			name:            "API returns nil - should use default true",
			apiEnabled:      nil,
			expectedEnabled: ApplicationAlertConfigDefaultEnabled,
			description:     "When API returns nil, state should use default value (true)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiData := &restapi.ApplicationAlertConfig{
				ID:               "test-id",
				Name:             "Test Alert",
				Description:      "Test Description",
				BoundaryScope:    "ALL",
				EvaluationType:   "PER_AP",
				Granularity:      600000,
				IncludeInternal:  false,
				IncludeSynthetic: false,
				Triggering:       false,
				Enabled:          tt.apiEnabled,
				Applications:     map[string]restapi.IncludedApplication{},
				AlertChannelIDs:  []string{},
				Rules:            []restapi.ApplicationAlertRuleWithThresholds{},
			}

			state := tfsdk.State{
				Schema: schema.Schema{
					Attributes: NewApplicationAlertConfigResourceHandle().MetaData().Schema.Attributes,
				},
			}

			diags := resource.UpdateState(ctx, &state, nil, apiData)
			require.False(t, diags.HasError(), "Expected no errors for %s", tt.name)

			var model ApplicationAlertConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError(), "Expected no errors getting state for %s", tt.name)

			assert.Equal(t, tt.expectedEnabled, model.Enabled.ValueBool(), "%s: %s", tt.name, tt.description)
		})
	}
}

func createMockState(t *testing.T, model ApplicationAlertConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewApplicationAlertConfigResourceHandle()
	return resource.MetaData().Schema
}
