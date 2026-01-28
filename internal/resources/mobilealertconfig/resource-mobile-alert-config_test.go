package mobilealertconfig

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
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}

// mockMobileAlertAPI is a mock implementation for testing
type mockMobileAlertAPI struct {
	testutils.MockInstanaAPI
}

func (m *mockMobileAlertAPI) MobileAlertConfig() restapi.RestResource[*restapi.MobileAlertConfig] {
	return nil
}

// Helper function to create a test schema
func getTestSchema() schema.Schema {
	handle := NewMobileAlertConfigResourceHandle()
	return handle.MetaData().Schema
}

// Helper function to create mock state
func createMockState(t *testing.T, model MobileAlertConfigModel) *tfsdk.State {
	state := &tfsdk.State{
		Schema: getTestSchema(),
	}
	diags := state.Set(context.Background(), &model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)
	return state
}

// Helper function to create mock plan
func createMockPlan(t *testing.T, model MobileAlertConfigModel) *tfsdk.Plan {
	plan := &tfsdk.Plan{
		Schema: getTestSchema(),
	}
	diags := plan.Set(context.Background(), &model)
	require.False(t, diags.HasError(), "Failed to set plan: %v", diags)
	return plan
}

func TestNewMobileAlertConfigResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewMobileAlertConfigResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaMobileAlertConfig, metadata.ResourceName)
		assert.Equal(t, int64(0), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewMobileAlertConfigResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldID])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldName])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldDescription])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldMobileAppID])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldSeverity])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldTriggering])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldTagFilter])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldAlertChannels])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldGranularity])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldGracePeriod])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldCustomPayloadFields])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldRules])
		assert.NotNil(t, schema.Attributes[MobileAlertConfigFieldTimeThreshold])
	})

	t.Run("should have computed ID field", func(t *testing.T) {
		handle := NewMobileAlertConfigResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		idAttr := schema.Attributes[MobileAlertConfigFieldID]
		assert.NotNil(t, idAttr)
	})

	t.Run("should have required name field", func(t *testing.T) {
		handle := NewMobileAlertConfigResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		nameAttr := schema.Attributes[MobileAlertConfigFieldName]
		assert.NotNil(t, nameAttr)
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &mobileAlertConfigResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaMobileAlertConfig,
				SchemaVersion: 0,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaMobileAlertConfig, metadata.ResourceName)
		assert.Equal(t, int64(0), metadata.SchemaVersion)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return mobile alert config rest resource", func(t *testing.T) {
		resource := &mobileAlertConfigResource{}

		mockAPI := &mockMobileAlertAPI{}
		restResource := resource.GetRestResource(mockAPI)

		assert.Nil(t, restResource) // Mock returns nil
	})
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should not return errors", func(t *testing.T) {
		resource := NewMobileAlertConfigResourceHandle()
		ctx := context.Background()

		plan := &tfsdk.Plan{
			Schema: resource.MetaData().Schema,
		}

		diags := resource.SetComputedFields(ctx, plan)
		assert.False(t, diags.HasError())
	})
}

func TestMapStateToDataObject_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Mobile Alert", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, "mobile-app-123", result.MobileAppID)
	assert.NotNil(t, result.Severity)
	assert.Equal(t, 5, *result.Severity)
	assert.False(t, result.Triggering)
	assert.Equal(t, restapi.Granularity(600000), result.Granularity)
}

func TestMapStateToDataObject_WithGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),
		GracePeriod: types.Int64Value(300000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.GracePeriod)
	assert.Equal(t, int64(300000), *result.GracePeriod)
}

func TestMapStateToDataObject_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),
		TagFilter:     types.StringValue("entity.type EQUALS 'mobileApp'"),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)
	require.NotNil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	alertChannelsMap := map[string]attr.Value{
		"5": types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-1"),
			types.StringValue("channel-2"),
		}),
		"10": types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("channel-3"),
		}),
	}

	state := createMockState(t, MobileAlertConfigModel{
		ID:            types.StringValue("test-id"),
		Name:          types.StringValue("Test Mobile Alert"),
		Description:   types.StringValue("Test Description"),
		MobileAppID:   types.StringValue("mobile-app-123"),
		Triggering:    types.BoolValue(false),
		Granularity:   types.Int64Value(600000),
		AlertChannels: types.MapValueMust(types.SetType{ElemType: types.StringType}, alertChannelsMap),

		TagFilter:   types.StringNull(),
		GracePeriod: types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.AlertChannels)
	assert.Len(t, result.AlertChannels, 2)
}

func TestMapStateToDataObject_WithRules(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	aggregation := restapi.Aggregation("MEAN")
	operator := "EQUALS"
	value := "test-value"

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules: []MobileRuleWithThresholdModel{
			{
				Rule: &MobileAlertRuleModel{
					AlertType:   types.StringValue("httpError"),
					MetricName:  types.StringValue("errors"),
					Aggregation: types.StringValue(string(aggregation)),
					Operator:    types.StringValue(operator),
					Value:       types.StringValue(value),
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &shared.ThresholdAllPluginModel{
					Warning: &shared.ThresholdAllTypeModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(100.0),
						},
					},
					Critical: &shared.ThresholdAllTypeModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(200.0),
						},
					},
				},
			},
		},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.NotNil(t, rule.Rule)
	assert.Equal(t, "httpError", rule.Rule.AlertType)
	assert.Equal(t, "errors", rule.Rule.MetricName)
	assert.NotNil(t, rule.Rule.Aggregation)
	assert.Equal(t, aggregation, *rule.Rule.Aggregation)
	assert.NotNil(t, rule.Rule.Operator)
	assert.Equal(t, operator, *rule.Rule.Operator)
	assert.NotNil(t, rule.Rule.Value)
	assert.Equal(t, value, *rule.Rule.Value)
	assert.Equal(t, ">", rule.ThresholdOperator)
}

func TestMapStateToDataObject_TimeThresholdTypes(t *testing.T) {
	tests := []struct {
		name          string
		timeThreshold *MobileAlertTimeThresholdModel
		expectedType  string
	}{
		{
			name: "violations_in_sequence",
			timeThreshold: &MobileAlertTimeThresholdModel{
				ViolationsInSequence: &MobileViolationsInSequenceModel{
					TimeWindow: types.Int64Value(600000),
				},
			},
			expectedType: MobileAlertConfigTimeThresholdTypeViolationsInSequence,
		},
		{
			name: "user_impact_of_violations_in_sequence",
			timeThreshold: &MobileAlertTimeThresholdModel{
				UserImpactOfViolationsInSequence: &MobileUserImpactOfViolationsInSequenceModel{
					TimeWindow: types.Int64Value(600000),
					Users:      types.Int64Value(100),
					Percentage: types.Float64Value(50.0),
				},
			},
			expectedType: MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence,
		},
		{
			name: "violations_in_period",
			timeThreshold: &MobileAlertTimeThresholdModel{
				ViolationsInPeriod: &MobileViolationsInPeriodModel{
					TimeWindow: types.Int64Value(600000),
					Violations: types.Int64Value(5),
				},
			},
			expectedType: MobileAlertConfigTimeThresholdTypeViolationsInPeriod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &mobileAlertConfigResource{}

			state := createMockState(t, MobileAlertConfigModel{
				ID:          types.StringValue("test-id"),
				Name:        types.StringValue("Test Mobile Alert"),
				Description: types.StringValue("Test Description"),
				MobileAppID: types.StringValue("mobile-app-123"),
				Triggering:  types.BoolValue(false),
				Granularity: types.Int64Value(600000),

				TagFilter:     types.StringNull(),
				AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
				GracePeriod:         types.Int64Null(),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Rules:               []MobileRuleWithThresholdModel{},
				TimeThreshold:       tt.timeThreshold,
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			require.NotNil(t, result.TimeThreshold)
			assert.Equal(t, tt.expectedType, result.TimeThreshold.Type)
		})
	}
}

func TestMapStateToDataObject_MissingTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold:       nil,
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	assert.True(t, diags.HasError())
	assert.Nil(t, result)
}

func TestMapStateToDataObject_InvalidTimeThresholdConfig(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold:       &MobileAlertTimeThresholdModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	assert.True(t, diags.HasError())
	assert.Nil(t, result)
}

func TestUpdateState_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Triggering:  false,
		Granularity: 600000,
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Mobile Alert", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	assert.Equal(t, "mobile-app-123", model.MobileAppID.ValueString())
	assert.False(t, model.Triggering.ValueBool())
	assert.Equal(t, int64(600000), model.Granularity.ValueInt64())
}

func TestUpdateState_WithSeverity(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	severity := 10
	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Severity:    &severity,
		Triggering:  false,
		Granularity: 600000,
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())
}

func TestUpdateState_WithGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	gracePeriod := int64(300000)
	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		GracePeriod: &gracePeriod,
		Triggering:  false,
		Granularity: 600000,
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.GracePeriod.IsNull())
	assert.Equal(t, int64(300000), model.GracePeriod.ValueInt64())
}

func TestUpdateState_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Triggering:  false,
		Granularity: 600000,
		AlertChannels: map[string][]string{
			"5":  {"channel-1", "channel-2"},
			"10": {"channel-3"},
		},
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.AlertChannels.IsNull())
}

func TestUpdateState_WithRules(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	aggregation := restapi.Aggregation("MEAN")
	operator := "EQUALS"
	value := "test-value"
	operatorLogical := restapi.LogicalOperatorType("AND")

	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Triggering:  false,
		Granularity: 600000,
		Rules: []restapi.MobileAppAlertRuleWithThresholds{
			{
				Rule: &restapi.MobileAppAlertRule{
					AlertType:   "httpError",
					MetricName:  "errors",
					Aggregation: &aggregation,
					Operator:    &operator,
					Value:       &value,
				},
				ThresholdOperator: ">",
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:     "staticThreshold",
						Operator: ptr(">"),
						Value:    ptr(100.0),
					},
					restapi.CriticalSeverity: {
						Type:     "staticThreshold",
						Operator: ptr(">"),
						Value:    ptr(200.0),
					},
				},
			},
		},
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operatorLogical,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Len(t, model.Rules, 1)
	assert.NotNil(t, model.Rules[0].Rule)
	assert.Equal(t, "httpError", model.Rules[0].Rule.AlertType.ValueString())
	assert.Equal(t, "errors", model.Rules[0].Rule.MetricName.ValueString())
}

func TestUpdateState_TimeThresholdTypes(t *testing.T) {
	tests := []struct {
		name          string
		timeThreshold *restapi.MobileAppTimeThreshold
		expectedType  string
	}{
		{
			name: "violations_in_sequence",
			timeThreshold: &restapi.MobileAppTimeThreshold{
				Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
				TimeWindow: ptr(int64(600000)),
			},
			expectedType: MobileAlertConfigTimeThresholdTypeViolationsInSequence,
		},
		{
			name: "user_impact_of_violations_in_sequence",
			timeThreshold: &restapi.MobileAppTimeThreshold{
				Type:       MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence,
				TimeWindow: ptr(int64(600000)),
				Users:      ptr(int32(100)),
				Percentage: ptr(50.0),
			},
			expectedType: MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence,
		},
		{
			name: "violations_in_period",
			timeThreshold: &restapi.MobileAppTimeThreshold{
				Type:       MobileAlertConfigTimeThresholdTypeViolationsInPeriod,
				TimeWindow: ptr(int64(600000)),
				Violations: ptr(int32(5)),
			},
			expectedType: MobileAlertConfigTimeThresholdTypeViolationsInPeriod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &mobileAlertConfigResource{}

			operator := restapi.LogicalOperatorType("AND")
			data := &restapi.MobileAlertConfig{
				ID:          "test-id",
				Name:        "Test Mobile Alert",
				Description: "Test Description",
				MobileAppID: "mobile-app-123",
				Triggering:  false,
				Granularity: 600000,
				TagFilterExpression: &restapi.TagFilter{
					Type:            "EXPRESSION",
					LogicalOperator: &operator,
					Elements:        []*restapi.TagFilter{},
				},
				TimeThreshold:         tt.timeThreshold,
				CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
				Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
			}

			state := &tfsdk.State{
				Schema: getTestSchema(),
			}

			diags := resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError())

			var model MobileAlertConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			require.NotNil(t, model.TimeThreshold)
			switch tt.expectedType {
			case MobileAlertConfigTimeThresholdTypeViolationsInSequence:
				assert.NotNil(t, model.TimeThreshold.ViolationsInSequence)
			case MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence:
				assert.NotNil(t, model.TimeThreshold.UserImpactOfViolationsInSequence)
			case MobileAlertConfigTimeThresholdTypeViolationsInPeriod:
				assert.NotNil(t, model.TimeThreshold.ViolationsInPeriod)
			}
		})
	}
}

func TestMapOptionalFields(t *testing.T) {
	resource := &mobileAlertConfigResource{}

	t.Run("mapOptionalInt64Field with value", func(t *testing.T) {
		field := types.Int64Value(12345)
		result := resource.mapOptionalInt64Field(field)
		require.NotNil(t, result)
		assert.Equal(t, int64(12345), *result)
	})

	t.Run("mapOptionalInt64Field with null", func(t *testing.T) {
		field := types.Int64Null()
		result := resource.mapOptionalInt64Field(field)
		assert.Nil(t, result)
	})

	t.Run("mapOptionalInt64Field with unknown", func(t *testing.T) {
		field := types.Int64Unknown()
		result := resource.mapOptionalInt64Field(field)
		assert.Nil(t, result)
	})

	t.Run("mapOptionalFloat64Field with value", func(t *testing.T) {
		field := types.Float64Value(123.45)
		result := resource.mapOptionalFloat64Field(field)
		require.NotNil(t, result)
		assert.Equal(t, 123.45, *result)
	})

	t.Run("mapOptionalFloat64Field with null", func(t *testing.T) {
		field := types.Float64Null()
		result := resource.mapOptionalFloat64Field(field)
		assert.Nil(t, result)
	})

	t.Run("mapOptionalInt64ToInt32Field with value", func(t *testing.T) {
		field := types.Int64Value(12345)
		result := resource.mapOptionalInt64ToInt32Field(field)
		require.NotNil(t, result)
		assert.Equal(t, int32(12345), *result)
	})

	t.Run("mapOptionalInt64ToInt32Field with null", func(t *testing.T) {
		field := types.Int64Null()
		result := resource.mapOptionalInt64ToInt32Field(field)
		assert.Nil(t, result)
	})
}

func TestGetStateUpgraders(t *testing.T) {
	t.Run("should return empty map", func(t *testing.T) {
		resource := &mobileAlertConfigResource{}
		ctx := context.Background()

		upgraders := resource.GetStateUpgraders(ctx)
		assert.NotNil(t, upgraders)
		assert.Empty(t, upgraders)
	})
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	plan := createMockPlan(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
}

func TestMapStateToDataObject_WithNullRule(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules: []MobileRuleWithThresholdModel{
			{
				Rule:              nil,
				ThresholdOperator: types.StringValue(">"),
			},
		},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Len(t, result.Rules, 0)
}

func TestUpdateState_WithPlan(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Triggering:  false,
		Granularity: 600000,
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	plan := createMockPlan(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, plan, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
}

func TestMapStateToDataObject_WithNullSeverity(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules:               []MobileRuleWithThresholdModel{},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Nil(t, result.Severity)
}

func TestUpdateState_WithNullSeverity(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Severity:    nil,
		Triggering:  false,
		Granularity: 600000,
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())
}

func TestMapStateToDataObject_WithRulesNullOptionalFields(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	state := createMockState(t, MobileAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Mobile Alert"),
		Description: types.StringValue("Test Description"),
		MobileAppID: types.StringValue("mobile-app-123"),
		Triggering:  types.BoolValue(false),
		Granularity: types.Int64Value(600000),

		TagFilter:     types.StringNull(),
		AlertChannels: types.MapNull(types.SetType{ElemType: types.StringType}),
		GracePeriod:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rules: []MobileRuleWithThresholdModel{
			{
				Rule: &MobileAlertRuleModel{
					AlertType:   types.StringValue("httpError"),
					MetricName:  types.StringValue("errors"),
					Aggregation: types.StringNull(),
					Operator:    types.StringNull(),
					Value:       types.StringNull(),
				},
				ThresholdOperator: types.StringValue(">"),
				Thresholds: &shared.ThresholdAllPluginModel{
					Warning: &shared.ThresholdAllTypeModel{
						Static: &shared.StaticTypeModel{
							Value: types.Float64Value(100.0),
						},
					},
				},
			},
		},
		TimeThreshold: &MobileAlertTimeThresholdModel{
			ViolationsInSequence: &MobileViolationsInSequenceModel{
				TimeWindow: types.Int64Value(600000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.NotNil(t, rule.Rule)
	assert.Nil(t, rule.Rule.Aggregation)
	assert.Nil(t, rule.Rule.Operator)
	assert.Nil(t, rule.Rule.Value)
}

func TestUpdateState_WithRulesNullOptionalFields(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:          "test-id",
		Name:        "Test Mobile Alert",
		Description: "Test Description",
		MobileAppID: "mobile-app-123",
		Triggering:  false,
		Granularity: 600000,
		Rules: []restapi.MobileAppAlertRuleWithThresholds{
			{
				Rule: &restapi.MobileAppAlertRule{
					AlertType:   "httpError",
					MetricName:  "errors",
					Aggregation: nil,
					Operator:    nil,
					Value:       nil,
				},
				ThresholdOperator: ">",
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:     "staticThreshold",
						Operator: ptr(">"),
						Value:    ptr(100.0),
					},
				},
			},
		},
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Len(t, model.Rules, 1)
	assert.NotNil(t, model.Rules[0].Rule)
	assert.True(t, model.Rules[0].Rule.Aggregation.IsNull())
	assert.True(t, model.Rules[0].Rule.Operator.IsNull())
	assert.True(t, model.Rules[0].Rule.Value.IsNull())
}

func TestUpdateState_EmptyAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	operator := restapi.LogicalOperatorType("AND")
	data := &restapi.MobileAlertConfig{
		ID:            "test-id",
		Name:          "Test Mobile Alert",
		Description:   "Test Description",
		MobileAppID:   "mobile-app-123",
		Triggering:    false,
		Granularity:   600000,
		AlertChannels: map[string][]string{},
		TagFilterExpression: &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		},
		TimeThreshold: &restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		Rules:                 []restapi.MobileAppAlertRuleWithThresholds{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model MobileAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.AlertChannels.IsNull())
}

func TestMapTimeThresholdToState_AllTypes(t *testing.T) {
	resource := &mobileAlertConfigResource{}

	t.Run("violations_in_sequence", func(t *testing.T) {
		threshold := restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
		}

		result := resource.mapTimeThresholdToState(threshold)
		require.NotNil(t, result)
		require.NotNil(t, result.ViolationsInSequence)
		assert.Equal(t, int64(600000), result.ViolationsInSequence.TimeWindow.ValueInt64())
	})

	t.Run("user_impact_of_violations_in_sequence", func(t *testing.T) {
		threshold := restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence,
			TimeWindow: ptr(int64(600000)),
			Users:      ptr(int32(100)),
			Percentage: ptr(50.0),
		}

		result := resource.mapTimeThresholdToState(threshold)
		require.NotNil(t, result)
		require.NotNil(t, result.UserImpactOfViolationsInSequence)
		assert.Equal(t, int64(600000), result.UserImpactOfViolationsInSequence.TimeWindow.ValueInt64())
		assert.Equal(t, int64(100), result.UserImpactOfViolationsInSequence.Users.ValueInt64())
		assert.Equal(t, 50.0, result.UserImpactOfViolationsInSequence.Percentage.ValueFloat64())
	})

	t.Run("violations_in_period", func(t *testing.T) {
		threshold := restapi.MobileAppTimeThreshold{
			Type:       MobileAlertConfigTimeThresholdTypeViolationsInPeriod,
			TimeWindow: ptr(int64(600000)),
			Violations: ptr(int32(5)),
		}

		result := resource.mapTimeThresholdToState(threshold)
		require.NotNil(t, result)
		require.NotNil(t, result.ViolationsInPeriod)
		assert.Equal(t, int64(600000), result.ViolationsInPeriod.TimeWindow.ValueInt64())
		assert.Equal(t, int64(5), result.ViolationsInPeriod.Violations.ValueInt64())
	})
}

func TestMapRuleToState(t *testing.T) {
	ctx := context.Background()
	resource := &mobileAlertConfigResource{}

	t.Run("with all fields", func(t *testing.T) {
		aggregation := restapi.Aggregation("MEAN")
		operator := "EQUALS"
		value := "test-value"

		rule := &restapi.MobileAppAlertRule{
			AlertType:   "httpError",
			MetricName:  "errors",
			Aggregation: &aggregation,
			Operator:    &operator,
			Value:       &value,
		}

		result := resource.mapRuleToState(ctx, rule)
		require.NotNil(t, result)
		assert.Equal(t, "httpError", result.AlertType.ValueString())
		assert.Equal(t, "errors", result.MetricName.ValueString())
		assert.Equal(t, "MEAN", result.Aggregation.ValueString())
		assert.Equal(t, "EQUALS", result.Operator.ValueString())
		assert.Equal(t, "test-value", result.Value.ValueString())
	})

	t.Run("with null optional fields", func(t *testing.T) {
		rule := &restapi.MobileAppAlertRule{
			AlertType:   "httpError",
			MetricName:  "errors",
			Aggregation: nil,
			Operator:    nil,
			Value:       nil,
		}

		result := resource.mapRuleToState(ctx, rule)
		require.NotNil(t, result)
		assert.Equal(t, "httpError", result.AlertType.ValueString())
		assert.Equal(t, "errors", result.MetricName.ValueString())
		assert.True(t, result.Aggregation.IsNull())
		assert.True(t, result.Operator.IsNull())
		assert.True(t, result.Value.IsNull())
	})
}
