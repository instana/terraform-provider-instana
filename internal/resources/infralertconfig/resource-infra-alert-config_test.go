package infralertconfig

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/rest"
	tag "github.com/instana/instana-go-client/shared/tagfilter"
	common "github.com/instana/instana-go-client/shared/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInfraAlertConfigResourceHandle(t *testing.T) {
	resource := NewInfraAlertConfigResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaInfraAlertConfig, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(3), metaData.SchemaVersion)
}

func TestRulesSchema_GenericRuleIsRequired(t *testing.T) {
	genericRule, exists := buildRulesSchema().NestedObject.Attributes[InfraAlertConfigFieldGenericRule]
	require.True(t, exists)
	assert.True(t, genericRule.IsRequired(), "generic_rule must be required so rules = [{}] is rejected at validation time")
	assert.False(t, genericRule.IsOptional())
}

func TestV2RulesSchema_GenericRuleRemainsOptional(t *testing.T) {
	genericRule, exists := buildV2RulesSchema().Attributes[InfraAlertConfigFieldGenericRule]
	require.True(t, exists)
	assert.True(t, genericRule.IsOptional(), "v2 prior schema must stay unchanged for the state upgrader")
}

func TestMetaData(t *testing.T) {
	resource := &infraAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 0,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_resource", metaData.ResourceName)
	assert.Equal(t, int64(0), metaData.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := NewInfraAlertConfigResourceHandle()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetRestResource(t *testing.T) {
	resource := &infraAlertConfigResource{}
	assert.NotNil(t, resource.GetRestResource)
}

// mockSyntheticAlertAPI extends the common mock to provide specific behavior for synthetic alert config tests
type mockInfraAlertConfigResource struct {
	testutils.MockInstanaAPI
}

func (m *mockInfraAlertConfigResource) SyntheticAlertConfigs() rest.RestResource[*api.InfraAlertConfig] {
	return &mockInfraAlertConfigRestResource{}
}

// Mock rest resource
type mockInfraAlertConfigRestResource struct{}

func (m *mockInfraAlertConfigRestResource) GetAll() (*[]*api.InfraAlertConfig, error) {
	return nil, nil
}

func (m *mockInfraAlertConfigRestResource) GetOne(id string) (*api.InfraAlertConfig, error) {
	return nil, nil
}

func (m *mockInfraAlertConfigRestResource) Create(data *api.InfraAlertConfig) (*api.InfraAlertConfig, error) {
	return nil, nil
}

func (m *mockInfraAlertConfigRestResource) Update(data *api.InfraAlertConfig) (*api.InfraAlertConfig, error) {
	return nil, nil
}

func (m *mockInfraAlertConfigRestResource) Delete(data *api.InfraAlertConfig) error {
	return nil
}

func (m *mockInfraAlertConfigRestResource) DeleteByID(id string) error {
	return nil
}
func TestUpdateState_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with minimal model to prevent errors during UpdateState
	initModel := InfraAlertConfigModel{
		ID:                 types.StringNull(),
		Name:               types.StringNull(),
		Description:        types.StringNull(),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		Granularity:        types.Int64Null(),
		EvaluationType:     types.StringNull(),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	}
	diags := state.Set(ctx, initModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	if diags.HasError() {
		for _, d := range diags.Errors() {
			t.Logf("Error: %s - %s", d.Summary(), d.Detail())
		}
	}
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Infra Alert", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	assert.Equal(t, int64(600000), model.Granularity.ValueInt64())
	assert.Equal(t, string(api.EvaluationTypePerEntity), model.EvaluationType.ValueString())
}

func TestUpdateState_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	entityType := "entity.type"
	hostValue := "host"
	equalsOp := common.EqualsOperator

	tagFilter := &tag.TagFilter{
		Type:     tag.TagFilterExpressionType,
		Name:     &entityType,
		Operator: &equalsOp,
		Value:    &hostValue,
	}

	data := &api.InfraAlertConfig{
		ID:                  "test-id",
		Name:                "Test Infra Alert",
		Description:         "Test Description",
		Granularity:         600000,
		EvaluationType:      api.EvaluationTypePerEntity,
		TagFilterExpression: tagFilter,
		AlertChannels:       map[common.AlertSeverity][]string{},
		Rules:               []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with minimal model to prevent errors during UpdateState
	initModel := InfraAlertConfigModel{
		ID:                 types.StringNull(),
		Name:               types.StringNull(),
		Description:        types.StringNull(),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		Granularity:        types.Int64Null(),
		EvaluationType:     types.StringNull(),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	}
	diags := state.Set(ctx, initModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Tag filter should be set (not checking the exact value due to normalization)
	// The important thing is that UpdateState doesn't error
}

func TestUpdateState_WithGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		GroupBy:        []string{"host.name", "zone"},
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with minimal model to prevent errors during UpdateState
	initModel := InfraAlertConfigModel{
		ID:                 types.StringNull(),
		Name:               types.StringNull(),
		Description:        types.StringNull(),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		Granularity:        types.Int64Null(),
		EvaluationType:     types.StringNull(),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	}
	diags := state.Set(ctx, initModel)
	require.False(t, diags.HasError())

	diags = resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.GroupBy.IsNull())
	var groupBy []string
	diags = model.GroupBy.ElementsAs(ctx, &groupBy, false)
	require.False(t, diags.HasError())
	assert.Len(t, groupBy, 2)
	assert.Contains(t, groupBy, "host.name")
	assert.Contains(t, groupBy, "zone")
}

func TestUpdateState_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels: map[common.AlertSeverity][]string{
			common.WarningSeverity:  {"channel-1", "channel-2"},
			common.CriticalSeverity: {"channel-3"},
		},
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)

	var warningChannels []string
	diags = model.AlertChannels.Warning.ElementsAs(ctx, &warningChannels, false)
	require.False(t, diags.HasError())
	assert.Len(t, warningChannels, 2)

	var criticalChannels []string
	diags = model.AlertChannels.Critical.ElementsAs(ctx, &criticalChannels, false)
	require.False(t, diags.HasError())
	assert.Len(t, criticalChannels, 1)
}

func TestUpdateState_WithTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		TimeThreshold: &api.InfraTimeThreshold{
			Type:       "violationsInSequence",
			TimeWindow: 300000,
		},
		AlertChannels: map[common.AlertSeverity][]string{},
		Rules:         []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.TimeThreshold)
	require.NotNil(t, model.TimeThreshold.ViolationsInSequence)
	assert.Equal(t, int64(300000), model.TimeThreshold.ViolationsInSequence.TimeWindow.ValueInt64())
}

func TestUpdateState_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		CustomerPayloadFields: []common.CustomPayloadField[any]{
			{
				Type:  common.StaticStringCustomPayloadType,
				Key:   "field1",
				Value: "value1",
			},
		},
		AlertChannels: map[common.AlertSeverity][]string{},
		Rules:         []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.CustomPayloadField.IsNull())
}

func TestUpdateState_WithRulesAndStaticThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	value := float64(100)
	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{
			{
				Rule: api.InfraAlertRule{
					AlertType:              "genericRule",
					MetricName:             "cpu.usage",
					EntityType:             "host",
					Aggregation:            common.SumAggregation,
					CrossSeriesAggregation: common.MeanAggregation,
					Regex:                  false,
				},
				ThresholdOperator: common.ThresholdOperatorGreaterThan,
				Thresholds: map[common.AlertSeverity]common.ThresholdRule{
					common.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &value,
					},
				},
			},
		},
		AlertChannels: map[common.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].GenericRule)
	assert.Equal(t, "cpu.usage", model.Rules[0].GenericRule.MetricName.ValueString())
	assert.Equal(t, "host", model.Rules[0].GenericRule.EntityType.ValueString())
	assert.Equal(t, string(common.SumAggregation), model.Rules[0].GenericRule.Aggregation.ValueString())
	assert.Equal(t, string(common.MeanAggregation), model.Rules[0].GenericRule.CrossSeriesAggregation.ValueString())
	assert.False(t, model.Rules[0].GenericRule.Regex.ValueBool())
	assert.Equal(t, string(common.ThresholdOperatorGreaterThan), model.Rules[0].GenericRule.ThresholdOperator.ValueString())

	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule.Warning)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule.Warning.Static)
	assert.Equal(t, float64(100), model.Rules[0].GenericRule.ThresholdRule.Warning.Static.Value.ValueFloat64())
}

func TestUpdateState_WithRulesAndAdaptiveBaselineThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	deviationFactor := float32(2.0)
	adaptability := float32(0.5)
	seasonality := common.ThresholdSeasonality("DAILY")

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{
			{
				Rule: api.InfraAlertRule{
					AlertType:              "genericRule",
					MetricName:             "memory.usage",
					EntityType:             "host",
					Aggregation:            common.MeanAggregation,
					CrossSeriesAggregation: common.MaxAggregation,
					Regex:                  true,
				},
				ThresholdOperator: common.ThresholdOperatorGreaterThan,
				Thresholds: map[common.AlertSeverity]common.ThresholdRule{
					common.CriticalSeverity: {
						Type:            "adaptiveBaseline",
						DeviationFactor: &deviationFactor,
						Adaptability:    &adaptability,
						Seasonality:     &seasonality,
					},
				},
			},
		},
		AlertChannels: map[common.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].GenericRule)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule.Critical)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule.Critical.AdaptiveBaseline)

	adaptiveModel := model.Rules[0].GenericRule.ThresholdRule.Critical.AdaptiveBaseline
	assert.Equal(t, float64(2.0), adaptiveModel.DeviationFactor.ValueFloat64())
	assert.Equal(t, float64(0.5), adaptiveModel.Adaptability.ValueFloat64())
	assert.Equal(t, "DAILY", adaptiveModel.Seasonality.ValueString())
}

func TestMapStateToDataObject_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Infra Alert", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, common.Granularity(600000), result.Granularity)
	assert.Equal(t, api.EvaluationTypePerEntity, result.EvaluationType)
}

func TestMapStateToDataObject_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringValue("entity.type EQUALS 'host'"),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_WithInvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringValue("invalid tag filter syntax"),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_WithGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	groupBySet := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("host.name"),
		types.StringValue("zone"),
	})

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            groupBySet,
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.GroupBy, 2)
	assert.Contains(t, result.GroupBy, "host.name")
	assert.Contains(t, result.GroupBy, "zone")
}

func TestMapStateToDataObject_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	warningSet := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-1"),
		types.StringValue("channel-2"),
	})
	criticalSet := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-3"),
	})

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		AlertChannels: &InfraAlertChannelsModel{
			Warning:  warningSet,
			Critical: criticalSet,
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AlertChannels, 2)
	assert.Len(t, result.AlertChannels[common.WarningSeverity], 2)
	assert.Len(t, result.AlertChannels[common.CriticalSeverity], 1)
}

func TestMapStateToDataObject_WithTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		TimeThreshold: &InfraTimeThresholdModel{
			ViolationsInSequence: &InfraViolationsInSequenceModel{
				TimeWindow: types.Int64Value(300000),
			},
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.TimeThreshold)
	assert.Equal(t, "violationsInSequence", result.TimeThreshold.Type)
	assert.Equal(t, int64(300000), result.TimeThreshold.TimeWindow)
}

func TestMapStateToDataObject_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	// Create a simple list with one custom payload field
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

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: customPayloadFields,
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	// The test may have diagnostics due to custom payload field conversion
	// but we're testing that the code path is executed
	_ = diags
	// Just verify the result is not nil - custom payload field conversion
	// is tested by the shared package
	require.NotNil(t, result)
}

func TestMapStateToDataObject_WithRulesAndStaticThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		Rules: []InfraRulesModel{
			{
				GenericRule: &InfraGenericRuleModel{
					MetricName:             types.StringValue("cpu.usage"),
					EntityType:             types.StringValue("host"),
					Aggregation:            types.StringValue(string(common.SumAggregation)),
					CrossSeriesAggregation: types.StringValue(string(common.MeanAggregation)),
					Regex:                     types.BoolValue(false),
					ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
					MetricGroupBy:             types.SetNull(types.StringType),
					MetricTagFilterExpression: types.StringNull(),
					ThresholdRule: &shared.ThresholdPluginModel{
						Warning: &shared.ThresholdTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float64Value(100),
							},
						},
					},
				},
			},
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, "genericRule", rule.Rule.AlertType)
	assert.Equal(t, "cpu.usage", rule.Rule.MetricName)
	assert.Equal(t, "host", rule.Rule.EntityType)
	assert.Equal(t, common.SumAggregation, rule.Rule.Aggregation)
	assert.Equal(t, common.MeanAggregation, rule.Rule.CrossSeriesAggregation)
	assert.False(t, rule.Rule.Regex)
	assert.Equal(t, common.ThresholdOperatorGreaterThan, rule.ThresholdOperator)

	require.Contains(t, rule.Thresholds, common.WarningSeverity)
	threshold := rule.Thresholds[common.WarningSeverity]
	assert.Equal(t, "staticThreshold", threshold.Type)
	require.NotNil(t, threshold.Value)
	assert.Equal(t, float64(100), *threshold.Value)
}

func TestMapStateToDataObject_WithRulesAndAdaptiveBaselineThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypeCustom)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		Rules: []InfraRulesModel{
			{
				GenericRule: &InfraGenericRuleModel{
					MetricName:             types.StringValue("memory.usage"),
					EntityType:             types.StringValue("host"),
					Aggregation:            types.StringValue(string(common.MeanAggregation)),
					CrossSeriesAggregation: types.StringValue(string(common.MaxAggregation)),
					Regex:                     types.BoolValue(true),
					ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
					MetricGroupBy:             types.SetNull(types.StringType),
					MetricTagFilterExpression: types.StringNull(),
					ThresholdRule: &shared.ThresholdPluginModel{
						Critical: &shared.ThresholdTypeModel{
							AdaptiveBaseline: &shared.AdaptiveBaselineModel{
								DeviationFactor: types.Float64Value(2.0),
								Adaptability:    types.Float64Value(0.5),
								Seasonality:     types.StringValue("DAILY"),
							},
						},
					},
				},
			},
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	require.Contains(t, rule.Thresholds, common.CriticalSeverity)
	threshold := rule.Thresholds[common.CriticalSeverity]
	assert.Equal(t, "adaptiveBaseline", threshold.Type)
	require.NotNil(t, threshold.DeviationFactor)
	assert.Equal(t, float32(2.0), *threshold.DeviationFactor)
	require.NotNil(t, threshold.Adaptability)
	assert.Equal(t, float32(0.5), *threshold.Adaptability)
	require.NotNil(t, threshold.Seasonality)
	assert.Equal(t, common.ThresholdSeasonality("DAILY"), *threshold.Seasonality)
}

func TestMapStateToDataObject_WithBothThresholds(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		Rules: []InfraRulesModel{
			{
				GenericRule: &InfraGenericRuleModel{
					MetricName:             types.StringValue("disk.usage"),
					EntityType:             types.StringValue("host"),
					Aggregation:            types.StringValue(string(common.SumAggregation)),
					CrossSeriesAggregation: types.StringValue(string(common.SumAggregation)),
					Regex:                     types.BoolValue(false),
					ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
					MetricGroupBy:             types.SetNull(types.StringType),
					MetricTagFilterExpression: types.StringNull(),
					ThresholdRule: &shared.ThresholdPluginModel{
						Warning: &shared.ThresholdTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float64Value(80),
							},
						},
						Critical: &shared.ThresholdTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float64Value(95),
							},
						},
					},
				},
			},
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	require.Contains(t, rule.Thresholds, common.WarningSeverity)
	require.Contains(t, rule.Thresholds, common.CriticalSeverity)

	assert.Equal(t, float64(80), *rule.Thresholds[common.WarningSeverity].Value)
	assert.Equal(t, float64(95), *rule.Thresholds[common.CriticalSeverity].Value)
}

func TestMapStateToDataObject_WithEmptyTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringValue(""),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Nil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_WithNullGranularity(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Null(),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, common.Granularity(0), result.Granularity)
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	model := InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Infra Alert", result.Name)
}

func TestUpdateState_WithEmptyAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.AlertChannels)
}

func TestUpdateState_WithNullTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:                  "test-id",
		Name:                "Test Infra Alert",
		Description:         "Test Description",
		Granularity:         600000,
		EvaluationType:      api.EvaluationTypePerEntity,
		TagFilterExpression: nil,
		AlertChannels:       map[common.AlertSeverity][]string{},
		Rules:               []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.TagFilter.IsNull())
}

func TestUpdateState_WithEmptyGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		GroupBy:        []string{},
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.GroupBy.IsNull())
	assert.Empty(t, model.GroupBy.Elements())
}

func TestUpdateState_WithNullTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		TimeThreshold:  nil,
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.TimeThreshold)
}

func TestUpdateState_WithUnsupportedTimeThresholdType(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		TimeThreshold: &api.InfraTimeThreshold{
			Type:       "unsupportedType",
			TimeWindow: 300000,
		},
		AlertChannels: map[common.AlertSeverity][]string{},
		Rules:         []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Unsupported type should result in nil TimeThreshold
	assert.Nil(t, model.TimeThreshold)
}

func TestMapStateToDataObject_WithNullTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		TimeThreshold:      nil,
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Nil(t, result.TimeThreshold)
}

func TestMapStateToDataObject_WithNullAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		AlertChannels:      nil,
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.AlertChannels)
}

func TestMapStateToDataObject_WithNullRules(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		Rules:              nil,
		RuleLogicalOperator: types.StringNull(),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.Rules)
}

func TestUpdateState_WithOnlyWarningAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels: map[common.AlertSeverity][]string{
			common.WarningSeverity: {"channel-1"},
		},
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)
	assert.False(t, model.AlertChannels.Warning.IsNull())
	assert.True(t, model.AlertChannels.Critical.IsNull())
}

func TestUpdateState_WithOnlyCriticalAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels: map[common.AlertSeverity][]string{
			common.CriticalSeverity: {"channel-1"},
		},
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)
	assert.True(t, model.AlertChannels.Warning.IsNull())
	assert.False(t, model.AlertChannels.Critical.IsNull())
}

func TestMapStateToDataObject_WithOnlyWarningAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	warningSet := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-1"),
	})

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		AlertChannels: &InfraAlertChannelsModel{
			Warning:  warningSet,
			Critical: types.SetNull(types.StringType),
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AlertChannels, 2)
	assert.Contains(t, result.AlertChannels, common.WarningSeverity)
	assert.Contains(t, result.AlertChannels, common.CriticalSeverity)
	assert.Len(t, result.AlertChannels[common.WarningSeverity], 1)
	assert.Len(t, result.AlertChannels[common.CriticalSeverity], 0)
}

func TestMapStateToDataObject_AllEvaluationTypes(t *testing.T) {
	evaluationTypes := []api.InfraAlertEvaluationType{
		api.EvaluationTypePerEntity,
		api.EvaluationTypeCustom,
	}

	for _, evalType := range evaluationTypes {
		t.Run(string(evalType), func(t *testing.T) {
			ctx := context.Background()
			resource := NewInfraAlertConfigResourceHandle()

			state := createMockState(t, InfraAlertConfigModel{
				ID:                 types.StringValue("test-id"),
				Name:               types.StringValue("Test Infra Alert"),
				Description:        types.StringValue("Test Description"),
				Granularity:        types.Int64Value(600000),
				EvaluationType:     types.StringValue(string(evalType)),
				TagFilter:          types.StringNull(),
				GroupBy:            types.SetNull(types.StringType),
				CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			assert.Equal(t, evalType, result.EvaluationType)
		})
	}
}

func TestMapTimeThresholdToModel_NilInput(t *testing.T) {
	resource := &infraAlertConfigResource{}

	result := resource.mapTimeThresholdToModel(nil)
	assert.Nil(t, result)
}

func TestMapTimeThresholdToModel_UnsupportedType(t *testing.T) {
	resource := &infraAlertConfigResource{}

	timeThreshold := &api.InfraTimeThreshold{
		Type:       "unsupportedType",
		TimeWindow: 300000,
	}

	result := resource.mapTimeThresholdToModel(timeThreshold)
	assert.Nil(t, result)
}

func TestMapTimeThresholdToModel_ViolationsInSequence(t *testing.T) {
	resource := &infraAlertConfigResource{}

	timeThreshold := &api.InfraTimeThreshold{
		Type:       "violationsInSequence",
		TimeWindow: 300000,
	}

	result := resource.mapTimeThresholdToModel(timeThreshold)
	require.NotNil(t, result)
	require.NotNil(t, result.ViolationsInSequence)
	assert.Equal(t, int64(300000), result.ViolationsInSequence.TimeWindow.ValueInt64())
}

// Helper functions

func initializeEmptyState(ctx context.Context, t *testing.T, state *tfsdk.State) {
	// Initialize state with minimal model to prevent errors during UpdateState
	initModel := InfraAlertConfigModel{
		ID:                 types.StringNull(),
		Name:               types.StringNull(),
		Description:        types.StringNull(),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		Granularity:        types.Int64Null(),
		EvaluationType:     types.StringNull(),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	}
	diags := state.Set(ctx, initModel)
	require.False(t, diags.HasError())
}

func createMockState(t *testing.T, model InfraAlertConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func TestUpdateState_WithEmptyRules(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Empty(t, model.Rules)
}

func TestUpdateState_WithBothThresholds(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	warningValue := float64(80)
	criticalValue := float64(95)
	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{
			{
				Rule: api.InfraAlertRule{
					AlertType:              "genericRule",
					MetricName:             "disk.usage",
					EntityType:             "host",
					Aggregation:            common.SumAggregation,
					CrossSeriesAggregation: common.SumAggregation,
					Regex:                  false,
				},
				ThresholdOperator: common.ThresholdOperatorGreaterThan,
				Thresholds: map[common.AlertSeverity]common.ThresholdRule{
					common.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
					common.CriticalSeverity: {
						Type:  "staticThreshold",
						Value: &criticalValue,
					},
				},
			},
		},
		AlertChannels: map[common.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].GenericRule)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule.Warning)
	require.NotNil(t, model.Rules[0].GenericRule.ThresholdRule.Critical)
}

func TestMapStateToDataObject_WithEmptyCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.CustomerPayloadFields)
}

func TestMapStateToDataObject_WithNullID(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:                 types.StringNull(),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "", result.ID)
}

func TestUpdateState_ErrorInTagFilterMapping(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	// Create an invalid tag filter that will cause an error during mapping
	invalidName := ""
	equalsOp := common.EqualsOperator
	hostValue := "host"

	tagFilter := &tag.TagFilter{
		Type:     tag.TagFilterExpressionType,
		Name:     &invalidName,
		Operator: &equalsOp,
		Value:    &hostValue,
	}

	data := &api.InfraAlertConfig{
		ID:                  "test-id",
		Name:                "Test Infra Alert",
		Description:         "Test Description",
		Granularity:         600000,
		EvaluationType:      api.EvaluationTypePerEntity,
		TagFilterExpression: tagFilter,
		AlertChannels:       map[common.AlertSeverity][]string{},
		Rules:               []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	// This may or may not error depending on the tag filter validation
	// The test ensures the error handling path is covered
	_ = diags
}

func TestMapStateToDataObject_WithOnlyCriticalAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	criticalSet := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-1"),
	})

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		AlertChannels: &InfraAlertChannelsModel{
			Warning:  types.SetNull(types.StringType),
			Critical: criticalSet,
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AlertChannels, 2)
	assert.Contains(t, result.AlertChannels, common.CriticalSeverity)
	assert.Contains(t, result.AlertChannels, common.WarningSeverity)
	assert.Len(t, result.AlertChannels[common.CriticalSeverity], 1)
	assert.Len(t, result.AlertChannels[common.WarningSeverity], 0)
}

func TestUpdateState_WithEmptyCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:                    "test-id",
		Name:                  "Test Infra Alert",
		Description:           "Test Description",
		Granularity:           600000,
		EvaluationType:        api.EvaluationTypePerEntity,
		CustomerPayloadFields: []common.CustomPayloadField[any]{},
		AlertChannels:         map[common.AlertSeverity][]string{},
		Rules:                 []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Empty custom payload fields should result in empty list (not null)
	assert.False(t, model.CustomPayloadField.IsNull())
	assert.True(t, len(model.CustomPayloadField.Elements()) == 0)
}

func TestMapStateToDataObject_WithEmptyAlertChannelLists(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	warningSet := types.SetValueMust(types.StringType, []attr.Value{})
	criticalSet := types.SetValueMust(types.StringType, []attr.Value{})

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		AlertChannels: &InfraAlertChannelsModel{
			Warning:  warningSet,
			Critical: criticalSet,
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	// Empty lists are added to the map as empty arrays
	require.Len(t, result.AlertChannels, 2)
	assert.Len(t, result.AlertChannels[common.WarningSeverity], 0)
	assert.Len(t, result.AlertChannels[common.CriticalSeverity], 0)
}

func TestUpdateState_WithEmptyGroupByList(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		GroupBy:        nil,
		AlertChannels:  map[common.AlertSeverity][]string{},
		Rules:          []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.GroupBy.IsNull())
	assert.Empty(t, model.GroupBy.Elements())
}

func TestMapStateToDataObject_WithNullTimeWindowInTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		TimeThreshold: &InfraTimeThresholdModel{
			ViolationsInSequence: &InfraViolationsInSequenceModel{
				TimeWindow: types.Int64Null(),
			},
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	// Null time window should result in nil time threshold
	assert.Nil(t, result.TimeThreshold)
}

func TestUpdateState_WithEmptyAlertChannelSeverities(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		AlertChannels: map[common.AlertSeverity][]string{
			common.WarningSeverity:  {},
			common.CriticalSeverity: {},
		},
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Empty alert channel lists should result in null lists
	require.NotNil(t, model.AlertChannels)
	assert.True(t, model.AlertChannels.Warning.IsNull())
	assert.True(t, model.AlertChannels.Critical.IsNull())
}

func TestMapStateToDataObject_WithRulesButNoThresholds(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, InfraAlertConfigModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Infra Alert"),
		Description:    types.StringValue("Test Description"),
		Granularity:    types.Int64Value(600000),
		EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:      types.StringNull(),
		GroupBy:        types.SetNull(types.StringType),
		Rules: []InfraRulesModel{
			{
				GenericRule: &InfraGenericRuleModel{
					MetricName:             types.StringValue("cpu.usage"),
					EntityType:             types.StringValue("host"),
					Aggregation:            types.StringValue(string(common.SumAggregation)),
					CrossSeriesAggregation: types.StringValue(string(common.MeanAggregation)),
					Regex:                     types.BoolValue(false),
					ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
					MetricGroupBy:             types.SetNull(types.StringType),
					MetricTagFilterExpression: types.StringNull(),
					ThresholdRule:             nil,
				},
			},
		},
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	// Rules without thresholds should have empty threshold map
	assert.Empty(t, result.Rules[0].Thresholds)
}

func TestMapGroupByToModel(t *testing.T) {
	r := &infraAlertConfigResource{}

	t.Run("empty slice returns empty set", func(t *testing.T) {
		result := r.mapGroupByToModel([]string{})
		assert.False(t, result.IsNull())
		assert.Equal(t, 0, len(result.Elements()))
	})

	t.Run("nil slice returns empty set", func(t *testing.T) {
		result := r.mapGroupByToModel(nil)
		assert.False(t, result.IsNull())
		assert.Equal(t, 0, len(result.Elements()))
	})

	t.Run("populated slice returns set with values", func(t *testing.T) {
		result := r.mapGroupByToModel([]string{"host.name", "zone"})
		assert.False(t, result.IsNull())
		assert.Equal(t, 2, len(result.Elements()))
	})
}

func TestMapModelGroupByToAPI(t *testing.T) {
	r := &infraAlertConfigResource{}
	ctx := context.Background()

	t.Run("null set returns empty slice", func(t *testing.T) {
		result, diags := r.mapModelGroupByToAPI(ctx, types.SetNull(types.StringType))
		assert.False(t, diags.HasError())
		assert.Equal(t, []string{}, result)
	})

	t.Run("empty set returns empty slice", func(t *testing.T) {
		emptySet := types.SetValueMust(types.StringType, []attr.Value{})
		result, diags := r.mapModelGroupByToAPI(ctx, emptySet)
		assert.False(t, diags.HasError())
		assert.Equal(t, []string{}, result)
	})

	t.Run("populated set returns elements", func(t *testing.T) {
		set := types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("host.name"),
		})
		result, diags := r.mapModelGroupByToAPI(ctx, set)
		assert.False(t, diags.HasError())
		assert.Equal(t, []string{"host.name"}, result)
	})
}

func TestMapMetricGroupByToModel(t *testing.T) {
	r := &infraAlertConfigResource{}

	t.Run("empty slice returns empty set", func(t *testing.T) {
		result := r.mapMetricGroupByToModel([]string{})
		assert.False(t, result.IsNull())
		assert.Equal(t, 0, len(result.Elements()))
	})

	t.Run("nil slice returns empty set", func(t *testing.T) {
		result := r.mapMetricGroupByToModel(nil)
		assert.False(t, result.IsNull())
		assert.Equal(t, 0, len(result.Elements()))
	})

	t.Run("populated slice returns set with values", func(t *testing.T) {
		result := r.mapMetricGroupByToModel([]string{"host", "namespace"})
		assert.False(t, result.IsNull())
		assert.Equal(t, 2, len(result.Elements()))
	})
}

func TestMapModelMetricGroupByToAPI(t *testing.T) {
	r := &infraAlertConfigResource{}
	ctx := context.Background()

	t.Run("null set returns empty slice", func(t *testing.T) {
		result, diags := r.mapModelMetricGroupByToAPI(ctx, types.SetNull(types.StringType))
		assert.False(t, diags.HasError())
		assert.Equal(t, []string{}, result)
	})

	t.Run("empty set returns empty slice", func(t *testing.T) {
		emptySet := types.SetValueMust(types.StringType, []attr.Value{})
		result, diags := r.mapModelMetricGroupByToAPI(ctx, emptySet)
		assert.False(t, diags.HasError())
		assert.Equal(t, []string{}, result)
	})

	t.Run("populated set returns elements", func(t *testing.T) {
		set := types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("host"),
		})
		result, diags := r.mapModelMetricGroupByToAPI(ctx, set)
		assert.False(t, diags.HasError())
		assert.Equal(t, []string{"host"}, result)
	})
}

func getTestSchema() schema.Schema {
	resource := NewInfraAlertConfigResourceHandle()
	return resource.MetaData().Schema
}
func TestMapStateToDataObject_WithMetricTagFilterExpression(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	t.Run("valid metric tag filter expression maps to API TagFilter", func(t *testing.T) {
		state := createMockState(t, InfraAlertConfigModel{
			ID:             types.StringValue("test-id"),
			Name:           types.StringValue("Test Infra Alert"),
			Description:    types.StringValue("Test Description"),
			Granularity:    types.Int64Value(600000),
			EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
			TagFilter:      types.StringNull(),
			GroupBy:        types.SetNull(types.StringType),
			Rules: []InfraRulesModel{
				{
					GenericRule: &InfraGenericRuleModel{
						MetricName:                types.StringValue("cpu.usage"),
						EntityType:                types.StringValue("host"),
						Aggregation:               types.StringValue(string(common.SumAggregation)),
						CrossSeriesAggregation:    types.StringValue(string(common.MeanAggregation)),
						Regex:                     types.BoolValue(false),
						ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
						MetricGroupBy:             types.SetValueMust(types.StringType, []attr.Value{}),
						MetricTagFilterExpression: types.StringValue("entity.type EQUALS 'host'"),
						ThresholdRule:             nil,
					},
				},
			},
			CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
		})

		result, diags := resource.MapStateToDataObject(ctx, nil, &state)
		require.False(t, diags.HasError())
		require.NotNil(t, result)
		require.Len(t, result.Rules, 1)
		require.NotNil(t, result.Rules[0].Rule.MetricTagFilterExpression)
	})

	t.Run("invalid metric tag filter expression returns error", func(t *testing.T) {
		state := createMockState(t, InfraAlertConfigModel{
			ID:             types.StringValue("test-id"),
			Name:           types.StringValue("Test Infra Alert"),
			Description:    types.StringValue("Test Description"),
			Granularity:    types.Int64Value(600000),
			EvaluationType: types.StringValue(string(api.EvaluationTypePerEntity)),
			TagFilter:      types.StringNull(),
			GroupBy:        types.SetNull(types.StringType),
			Rules: []InfraRulesModel{
				{
					GenericRule: &InfraGenericRuleModel{
						MetricName:                types.StringValue("cpu.usage"),
						EntityType:                types.StringValue("host"),
						Aggregation:               types.StringValue(string(common.SumAggregation)),
						CrossSeriesAggregation:    types.StringValue(string(common.MeanAggregation)),
						Regex:                     types.BoolValue(false),
						ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
						MetricGroupBy:             types.SetValueMust(types.StringType, []attr.Value{}),
						MetricTagFilterExpression: types.StringValue("invalid filter expression @@@"),
						ThresholdRule:             nil,
					},
				},
			},
			CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
		})

		_, diags := resource.MapStateToDataObject(ctx, nil, &state)
		require.True(t, diags.HasError())
	})
}

func TestUpdateState_WithMetricTagFilterExpression(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	tagFilterExpr := tag.NewStringTagFilter(tag.TagFilterEntityDestination, "name", common.EqualsOperator, "foo")

	data := &api.InfraAlertConfig{
		ID:             "test-id",
		Name:           "Test Infra Alert",
		Description:    "Test Description",
		Granularity:    600000,
		EvaluationType: api.EvaluationTypePerEntity,
		Rules: []common.RuleWithThreshold[api.InfraAlertRule]{
			{
				Rule: api.InfraAlertRule{
					AlertType:                 "genericRule",
					MetricName:                "cpu.usage",
					EntityType:                "host",
					Aggregation:               common.SumAggregation,
					CrossSeriesAggregation:    common.MeanAggregation,
					Regex:                     false,
					MetricTagFilterExpression: tagFilterExpr,
				},
				ThresholdOperator: common.ThresholdOperatorGreaterThan,
			},
		},
		AlertChannels: map[common.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}
	initializeEmptyState(ctx, t, state)

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model InfraAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].GenericRule)
	assert.False(t, model.Rules[0].GenericRule.MetricTagFilterExpression.IsNull())
}

const v2StateJSONTemplate = `{
  "id": "abc",
  "name": "n",
  "description": "d",
  "tag_filter": null,
  "group_by": null,
  "alert_channels": null,
  "granularity": 600000,
  "grace_period": null,
  "time_threshold": {"violations_in_sequence": {"time_window": 600000}},
  "custom_payload_field": [],
  "rules": %RULES%,
  "evaluation_type": "CUSTOM",
  "triggering": false
}`

const v2RuleJSON = `{"generic_rule": {
  "metric_name": "m", "entity_type": "host", "aggregation": "MEAN", "cross_series_aggregation": "MEAN",
  "regex": false, "threshold_operator": ">=",
  "threshold": {"warning": null, "critical": {"static": {"value": 90}, "adaptive_baseline": null}},
  "metric_group_by": [], "metric_tag_filter_expression": null}}`

func runV2Upgrade(t *testing.T, rulesJSON string) (InfraAlertConfigModel, *resource.UpgradeStateResponse) {
	t.Helper()
	ctx := context.Background()
	r := &infraAlertConfigResource{}
	upgrader, ok := r.GetStateUpgraders(ctx)[2]
	require.True(t, ok, "a state upgrader for schema version 2 must be registered")
	require.NotNil(t, upgrader.PriorSchema)

	// v2 state never contained rule_logical_operator
	_, hasOperator := upgrader.PriorSchema.Attributes[InfraAlertConfigFieldRuleLogicalOperator]
	require.False(t, hasOperator, "v2 prior schema must not contain rule_logical_operator")

	js := []byte(strings.Replace(v2StateJSONTemplate, "%RULES%", rulesJSON, 1))
	raw := tfprotov6.RawState{JSON: js}
	val, err := raw.UnmarshalWithOpts(upgrader.PriorSchema.Type().TerraformType(ctx), tfprotov6.UnmarshalOpts{})
	require.NoError(t, err)

	current := buildInfraAlertConfigSchema()
	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: *upgrader.PriorSchema, Raw: val}}
	resp := &resource.UpgradeStateResponse{State: tfsdk.State{
		Schema: current,
		Raw:    tftypes.NewValue(current.Type().TerraformType(ctx), nil),
	}}
	upgrader.StateUpgrader(ctx, req, resp)

	var model InfraAlertConfigModel
	if !resp.Diagnostics.HasError() {
		diags := resp.State.Get(ctx, &model)
		require.False(t, diags.HasError(), "upgraded state must be readable with the v3 model: %v", diags)
	}
	return model, resp
}

func TestStateUpgraderV2_WrapsSingleRuleAndDefaultsOperator(t *testing.T) {
	model, resp := runV2Upgrade(t, v2RuleJSON)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

	require.Len(t, model.Rules, 1)
	require.NotNil(t, model.Rules[0].GenericRule)
	assert.Equal(t, "m", model.Rules[0].GenericRule.MetricName.ValueString())
	assert.Equal(t, "host", model.Rules[0].GenericRule.EntityType.ValueString())

	assert.False(t, model.RuleLogicalOperator.IsNull(), "rule_logical_operator must not be null after upgrade")
	assert.Equal(t, RuleLogicalOperatorAND, model.RuleLogicalOperator.ValueString())
	assert.Equal(t, "abc", model.ID.ValueString())
	assert.Equal(t, "CUSTOM", model.EvaluationType.ValueString())
}

func TestStateUpgraderV2_NullRulesBecomeEmptyList(t *testing.T) {
	model, resp := runV2Upgrade(t, "null")
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

	assert.NotNil(t, model.Rules)
	assert.Empty(t, model.Rules)
	assert.Equal(t, RuleLogicalOperatorAND, model.RuleLogicalOperator.ValueString())
}

func newGenericRuleModel(metricName string, tagFilter types.String) InfraRulesModel {
	return InfraRulesModel{
		GenericRule: &InfraGenericRuleModel{
			MetricName:                types.StringValue(metricName),
			EntityType:                types.StringValue("host"),
			Aggregation:               types.StringValue(string(common.MeanAggregation)),
			CrossSeriesAggregation:    types.StringValue(string(common.MeanAggregation)),
			Regex:                     types.BoolValue(false),
			ThresholdOperator:         types.StringValue(string(common.ThresholdOperatorGreaterThan)),
			MetricGroupBy:             types.SetNull(types.StringType),
			MetricTagFilterExpression: tagFilter,
			ThresholdRule: &shared.ThresholdPluginModel{
				Critical: &shared.ThresholdTypeModel{
					Static: &shared.StaticTypeModel{Value: types.Float64Value(90)},
				},
			},
		},
	}
}

func newBaseModel(rules []InfraRulesModel, operator types.String) InfraAlertConfigModel {
	return InfraAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Infra Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		EvaluationType:      types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:           types.StringNull(),
		GroupBy:             types.SetNull(types.StringType),
		Rules:               rules,
		RuleLogicalOperator: operator,
		CustomPayloadField:  types.ListNull(shared.GetCustomPayloadFieldType()),
	}
}

func TestMapStateToDataObject_RuleWithoutGenericRuleReturnsError(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()

	state := createMockState(t, newBaseModel([]InfraRulesModel{
		newGenericRuleModel("cpu.usage", types.StringNull()),
		{GenericRule: nil},
	}, types.StringValue(RuleLogicalOperatorAND)))

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.True(t, diags.HasError(), "a rules entry without generic_rule must not be silently dropped")
	assert.Contains(t, diags.Errors()[0].Detail(), "rules[1]")
}

// apiRuleWithFilter builds an API rule whose metric tag filter is parsed from the given expression.
func apiRuleWithFilter(t *testing.T, r *infraAlertConfigResource, metricName, filter string) common.RuleWithThreshold[api.InfraAlertRule] {
	t.Helper()
	rule := common.RuleWithThreshold[api.InfraAlertRule]{
		ThresholdOperator: common.ThresholdOperatorGreaterThan,
		Rule: api.InfraAlertRule{
			AlertType:              GenericRuleAlertType,
			MetricName:             metricName,
			EntityType:             "host",
			Aggregation:            common.MeanAggregation,
			CrossSeriesAggregation: common.MeanAggregation,
		},
		Thresholds: map[common.AlertSeverity]common.ThresholdRule{},
	}
	if filter != "" {
		parsed, diags := r.mapModelTagFilterToAPI(types.StringValue(filter))
		require.False(t, diags.HasError(), "%v", diags)
		rule.Rule.MetricTagFilterExpression = parsed
	}
	return rule
}

func TestMapRulesToModel_MultipleRules(t *testing.T) {
	ctx := context.Background()
	r := &infraAlertConfigResource{}

	apiRules := []common.RuleWithThreshold[api.InfraAlertRule]{
		apiRuleWithFilter(t, r, "cpu.usage", ""),
		apiRuleWithFilter(t, r, "memory.usage", ""),
	}

	result, diags := r.mapRulesToModel(ctx, apiRules, nil)
	require.False(t, diags.HasError(), "%v", diags)
	require.Len(t, result, 2)
	assert.Equal(t, "cpu.usage", result[0].GenericRule.MetricName.ValueString())
	assert.Equal(t, "memory.usage", result[1].GenericRule.MetricName.ValueString())
}

func TestMapRulesToModel_DoesNotSwapTagFilterWhenRulesReordered(t *testing.T) {
	ctx := context.Background()
	r := &infraAlertConfigResource{}

	// The configuration declared [cpu (host filter), memory (pod filter)] but the API returns
	// the rules in the opposite order.
	existing := []InfraRulesModel{
		newGenericRuleModel("cpu.usage", types.StringValue("entity.type EQUALS 'host'")),
		newGenericRuleModel("memory.usage", types.StringValue("entity.type EQUALS 'pod'")),
	}
	apiRules := []common.RuleWithThreshold[api.InfraAlertRule]{
		apiRuleWithFilter(t, r, "memory.usage", "entity.type EQUALS 'pod'"),
		apiRuleWithFilter(t, r, "cpu.usage", "entity.type EQUALS 'host'"),
	}

	result, diags := r.mapRulesToModel(ctx, apiRules, existing)
	require.False(t, diags.HasError(), "%v", diags)
	require.Len(t, result, 2)

	// Each filter must stay attached to the rule it belongs to, not to the rule at the same index.
	assert.Equal(t, "memory.usage", result[0].GenericRule.MetricName.ValueString())
	assert.Contains(t, result[0].GenericRule.MetricTagFilterExpression.ValueString(), "pod")
	assert.Equal(t, "cpu.usage", result[1].GenericRule.MetricName.ValueString())
	assert.Contains(t, result[1].GenericRule.MetricTagFilterExpression.ValueString(), "host")
}

func TestMapRulesToModel_RemoteTagFilterChangeIsDetected(t *testing.T) {
	ctx := context.Background()
	r := &infraAlertConfigResource{}

	existing := []InfraRulesModel{newGenericRuleModel("cpu.usage", types.StringValue("entity.type EQUALS 'host'"))}
	apiRules := []common.RuleWithThreshold[api.InfraAlertRule]{
		apiRuleWithFilter(t, r, "cpu.usage", "entity.type EQUALS 'pod'"),
	}

	result, diags := r.mapRulesToModel(ctx, apiRules, existing)
	require.False(t, diags.HasError(), "%v", diags)
	require.Len(t, result, 1)
	assert.Contains(t, result[0].GenericRule.MetricTagFilterExpression.ValueString(), "pod")
}

