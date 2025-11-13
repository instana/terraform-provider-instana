package logalertconfig

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}

func TestNewLogAlertConfigResourceHandleFramework(t *testing.T) {
	resource := NewLogAlertConfigResourceHandleFramework()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaLogAlertConfigFramework, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(1), metaData.SchemaVersion)
}

func TestMetaData(t *testing.T) {
	resource := &logAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  "test_resource",
			SchemaVersion: 1,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_resource", metaData.ResourceName)
	assert.Equal(t, int64(1), metaData.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := NewLogAlertConfigResourceHandleFramework()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetRestResource(t *testing.T) {
	resource := &logAlertConfigResourceFramework{}
	assert.NotNil(t, resource.GetRestResource)
}

func TestUpdateState_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Log Alert", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	assert.Equal(t, int64(600000), model.Granularity.ValueInt64())
}

func TestUpdateState_WithGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		GracePeriod:   300000,
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.GracePeriod.IsNull())
	assert.Equal(t, int64(300000), model.GracePeriod.ValueInt64())
}

func TestUpdateState_WithZeroGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		GracePeriod:   0,
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.GracePeriod.IsNull())
}

func TestUpdateState_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	// Create a tag filter using AND logic with two expressions
	tagFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "entity.type", restapi.EqualsOperator, "log"),
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "service.name", restapi.EqualsOperator, "test-service"),
	})

	data := &restapi.LogAlertConfig{
		ID:                  "test-id",
		Name:                "Test Log Alert",
		Description:         "Test Description",
		Granularity:         600000,
		TagFilterExpression: tagFilter,
		AlertChannels:       map[restapi.AlertSeverity][]string{},
		Rules:               []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError(), "UpdateState should not have errors: %v", diags)

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError(), "Getting state should not have errors: %v", diags)

	assert.False(t, model.TagFilter.IsNull(), "TagFilter should not be null")
	if !model.TagFilter.IsNull() {
		assert.NotEmpty(t, model.TagFilter.ValueString(), "TagFilter should have a value")
	}
}

func TestUpdateState_WithGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		GroupBy: []restapi.GroupByTag{
			{TagName: "host.name", Key: ""},
			{TagName: "service.name", Key: "key1"},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.GroupBy, 2)
	assert.Equal(t, "host.name", model.GroupBy[0].TagName.ValueString())
	assert.True(t, model.GroupBy[0].Key.IsNull())
	assert.Equal(t, "service.name", model.GroupBy[1].TagName.ValueString())
	assert.Equal(t, "key1", model.GroupBy[1].Key.ValueString())
}

func TestUpdateState_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		AlertChannels: map[restapi.AlertSeverity][]string{
			restapi.WarningSeverity:  {"channel-1", "channel-2"},
			restapi.CriticalSeverity: {"channel-3"},
		},
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
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
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		TimeThreshold: &restapi.LogTimeThreshold{
			Type:       "violationsInSequence",
			TimeWindow: 300000,
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.TimeThreshold)
	require.NotNil(t, model.TimeThreshold.ViolationsInSequence)
	assert.Equal(t, int64(300000), model.TimeThreshold.ViolationsInSequence.TimeWindow.ValueInt64())
}

func TestUpdateState_WithRulesAndThresholds(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningValue := float64(100)
	criticalValue := float64(200)

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{
			{
				Rule: restapi.LogAlertRule{
					AlertType:   "logCount",
					MetricName:  "log.count",
					Aggregation: restapi.SumAggregation,
				},
				ThresholdOperator: restapi.ThresholdOperatorGreaterThan,
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
					restapi.CriticalSeverity: {
						Type:  "staticThreshold",
						Value: &criticalValue,
					},
				},
			},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	assert.Equal(t, "log.count", model.Rules.MetricName.ValueString())
	assert.Equal(t, LogAlertTypeLogCount, model.Rules.AlertType.ValueString())
	assert.Equal(t, string(restapi.SumAggregation), model.Rules.Aggregation.ValueString())
	assert.Equal(t, string(restapi.ThresholdOperatorGreaterThan), model.Rules.ThresholdOperator.ValueString())

	require.NotNil(t, model.Rules.Threshold)
	require.NotNil(t, model.Rules.Threshold.Warning)
	require.NotNil(t, model.Rules.Threshold.Warning.Static)
	assert.Equal(t, int64(100), model.Rules.Threshold.Warning.Static.Value.ValueInt64())

	require.NotNil(t, model.Rules.Threshold.Critical)
	require.NotNil(t, model.Rules.Threshold.Critical.Static)
	assert.Equal(t, int64(200), model.Rules.Threshold.Critical.Static.Value.ValueInt64())
}

func TestUpdateState_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "field1",
				Value: "value1",
			},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.CustomPayloadFields.IsNull())
}

func TestMapStateToDataObject_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Log Alert", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, restapi.Granularity(600000), result.Granularity)
}

func TestMapStateToDataObject_WithGracePeriod(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		GracePeriod:         types.Int64Value(300000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, int64(300000), result.GracePeriod)
}

func TestMapStateToDataObject_WithNullGranularity(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Null(),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, restapi.Granularity600000, result.Granularity)
}

func TestMapStateToDataObject_WithTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_WithInvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("invalid tag filter syntax"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_WithGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		GroupBy: []GroupByModel{
			{
				TagName: types.StringValue("host.name"),
				Key:     types.StringNull(),
			},
			{
				TagName: types.StringValue("service.name"),
				Key:     types.StringValue("key1"),
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.GroupBy, 2)
	assert.Equal(t, "host.name", result.GroupBy[0].TagName)
	assert.Equal(t, "", result.GroupBy[0].Key)
	assert.Equal(t, "service.name", result.GroupBy[1].TagName)
	assert.Equal(t, "key1", result.GroupBy[1].Key)
}

func TestMapStateToDataObject_WithAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningList := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-1"),
		types.StringValue("channel-2"),
	})
	criticalList := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-3"),
	})

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		AlertChannels: &AlertChannelsModel{
			Warning:  warningList,
			Critical: criticalList,
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AlertChannels, 2)
	assert.Len(t, result.AlertChannels[restapi.WarningSeverity], 2)
	assert.Len(t, result.AlertChannels[restapi.CriticalSeverity], 1)
}

func TestMapStateToDataObject_WithTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		TimeThreshold: &TimeThresholdModel{
			ViolationsInSequence: &ViolationsInSequenceModel{
				TimeWindow: types.Int64Value(300000),
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.NotNil(t, result.TimeThreshold)
	assert.Equal(t, "violationsInSequence", result.TimeThreshold.Type)
	assert.Equal(t, int64(300000), result.TimeThreshold.TimeWindow)
}

func TestMapStateToDataObject_WithRules(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		Rules: &LogAlertRuleModel{
			MetricName:        types.StringValue("log.count"),
			AlertType:         types.StringValue(LogAlertTypeLogCount),
			Aggregation:       types.StringValue(string(restapi.SumAggregation)),
			ThresholdOperator: types.StringValue(string(restapi.ThresholdOperatorGreaterThan)),
			Threshold: &ThresholdModel{
				Warning: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Value(100),
					},
				},
				Critical: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Value(200),
					},
				},
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, "logCount", rule.Rule.AlertType)
	assert.Equal(t, "log.count", rule.Rule.MetricName)
	assert.Equal(t, restapi.SumAggregation, rule.Rule.Aggregation)
	assert.Equal(t, restapi.ThresholdOperatorGreaterThan, rule.ThresholdOperator)

	require.Contains(t, rule.Thresholds, restapi.WarningSeverity)
	assert.Equal(t, float64(100), *rule.Thresholds[restapi.WarningSeverity].Value)

	require.Contains(t, rule.Thresholds, restapi.CriticalSeverity)
	assert.Equal(t, float64(200), *rule.Thresholds[restapi.CriticalSeverity].Value)
}

func TestUpdateState_WithEmptyAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.AlertChannels)
}

func TestUpdateState_WithNullTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:                  "test-id",
		Name:                "Test Log Alert",
		Description:         "Test Description",
		Granularity:         600000,
		TagFilterExpression: nil,
		AlertChannels:       map[restapi.AlertSeverity][]string{},
		Rules:               []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.TagFilter.IsNull())
}

func TestUpdateState_WithEmptyGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		GroupBy:       []restapi.GroupByTag{},
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Empty(t, model.GroupBy)
}

func TestUpdateState_WithNullTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		TimeThreshold: nil,
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.TimeThreshold)
}

func TestUpdateState_WithUnsupportedTimeThresholdType(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		TimeThreshold: &restapi.LogTimeThreshold{
			Type:       "unsupportedType",
			TimeWindow: 300000,
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.TimeThreshold)
}

func TestMapStateToDataObject_WithNullTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		TimeThreshold:       nil,
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Nil(t, result.TimeThreshold)
}

func TestMapStateToDataObject_WithNullAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		AlertChannels:       nil,
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.AlertChannels)
}

func TestMapStateToDataObject_WithNullRules(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		Rules:               nil,
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.Rules)
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	model := LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Log Alert", result.Name)
}

func TestMapStateToDataObject_WithEmptyGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		GroupBy:             []GroupByModel{},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.GroupBy)
}

func TestMapStateToDataObject_WithNullTimeWindowInTimeThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		TimeThreshold: &TimeThresholdModel{
			ViolationsInSequence: &ViolationsInSequenceModel{
				TimeWindow: types.Int64Null(),
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Nil(t, result.TimeThreshold)
}

func TestUpdateState_WithEmptyRules(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:            "test-id",
		Name:          "Test Log Alert",
		Description:   "Test Description",
		Granularity:   600000,
		AlertChannels: map[restapi.AlertSeverity][]string{},
		Rules:         []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.Rules)
}

func TestUpdateState_WithRulesWithoutAggregation(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningValue := float64(100)

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{
			{
				Rule: restapi.LogAlertRule{
					AlertType:  "logCount",
					MetricName: "log.count",
					// No Aggregation specified
				},
				ThresholdOperator: restapi.ThresholdOperatorGreaterThan,
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
				},
			},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	// Should default to SUM
	assert.Equal(t, string(restapi.SumAggregation), model.Rules.Aggregation.ValueString())
}

func TestMapStateToDataObject_WithEmptyAlertChannelLists(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningList := types.ListValueMust(types.StringType, []attr.Value{})
	criticalList := types.ListValueMust(types.StringType, []attr.Value{})

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		AlertChannels: &AlertChannelsModel{
			Warning:  warningList,
			Critical: criticalList,
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.AlertChannels)
}

func TestMapStateToDataObject_WithNullID(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringNull(),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "", result.ID)
}

func TestMapTimeThresholdToState_NilInput(t *testing.T) {
	ctx := context.Background()
	resource := &logAlertConfigResourceFramework{}

	result := resource.mapTimeThresholdToState(ctx, nil)
	assert.Nil(t, result)
}

func TestMapTimeThresholdToState_UnsupportedType(t *testing.T) {
	ctx := context.Background()
	resource := &logAlertConfigResourceFramework{}

	timeThreshold := &restapi.LogTimeThreshold{
		Type:       "unsupportedType",
		TimeWindow: 300000,
	}

	result := resource.mapTimeThresholdToState(ctx, timeThreshold)
	assert.Nil(t, result)
}

func TestMapTimeThresholdToState_ViolationsInSequence(t *testing.T) {
	ctx := context.Background()
	resource := &logAlertConfigResourceFramework{}

	timeThreshold := &restapi.LogTimeThreshold{
		Type:       "violationsInSequence",
		TimeWindow: 300000,
	}

	result := resource.mapTimeThresholdToState(ctx, timeThreshold)
	require.NotNil(t, result)
	require.NotNil(t, result.ViolationsInSequence)
	assert.Equal(t, int64(300000), result.ViolationsInSequence.TimeWindow.ValueInt64())
}

func TestMapTimeThresholdFromState_NilInput(t *testing.T) {
	ctx := context.Background()
	resource := &logAlertConfigResourceFramework{}

	result := resource.mapTimeThresholdFromState(ctx, nil)
	assert.Nil(t, result)
}

func TestMapTimeThresholdFromState_NullTimeWindow(t *testing.T) {
	ctx := context.Background()
	resource := &logAlertConfigResourceFramework{}

	timeThreshold := &TimeThresholdModel{
		ViolationsInSequence: &ViolationsInSequenceModel{
			TimeWindow: types.Int64Null(),
		},
	}

	result := resource.mapTimeThresholdFromState(ctx, timeThreshold)
	assert.Nil(t, result)
}

func TestUpdateState_WithSingleWarningAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		AlertChannels: map[restapi.AlertSeverity][]string{
			restapi.WarningSeverity: {"channel-1"},
		},
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)
	assert.False(t, model.AlertChannels.Warning.IsNull())
	assert.True(t, model.AlertChannels.Critical.IsNull())
}

func TestUpdateState_WithSingleCriticalAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		AlertChannels: map[restapi.AlertSeverity][]string{
			restapi.CriticalSeverity: {"channel-1"},
		},
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)
	assert.True(t, model.AlertChannels.Warning.IsNull())
	assert.False(t, model.AlertChannels.Critical.IsNull())
}

func TestMapStateToDataObject_WithSingleWarningAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningList := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-1"),
	})

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		AlertChannels: &AlertChannelsModel{
			Warning:  warningList,
			Critical: types.ListNull(types.StringType),
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AlertChannels, 1)
	assert.Contains(t, result.AlertChannels, restapi.WarningSeverity)
	assert.NotContains(t, result.AlertChannels, restapi.CriticalSeverity)
}

func TestMapStateToDataObject_WithSingleCriticalAlertChannel(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	criticalList := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("channel-1"),
	})

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		AlertChannels: &AlertChannelsModel{
			Warning:  types.ListNull(types.StringType),
			Critical: criticalList,
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.AlertChannels, 1)
	assert.NotContains(t, result.AlertChannels, restapi.WarningSeverity)
	assert.Contains(t, result.AlertChannels, restapi.CriticalSeverity)
}

func TestMapStateToDataObject_WithRulesWithoutThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		Rules: &LogAlertRuleModel{
			MetricName:        types.StringValue("log.count"),
			AlertType:         types.StringValue(LogAlertTypeLogCount),
			Aggregation:       types.StringValue(string(restapi.SumAggregation)),
			ThresholdOperator: types.StringValue(string(restapi.ThresholdOperatorGreaterThan)),
			Threshold:         nil,
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	assert.Empty(t, result.Rules[0].Thresholds)
}

func TestMapStateToDataObject_WithRulesWithOnlyWarningThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		Rules: &LogAlertRuleModel{
			MetricName:        types.StringValue("log.count"),
			AlertType:         types.StringValue(LogAlertTypeLogCount),
			Aggregation:       types.StringValue(string(restapi.SumAggregation)),
			ThresholdOperator: types.StringValue(string(restapi.ThresholdOperatorGreaterThan)),
			Threshold: &ThresholdModel{
				Warning: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Value(100),
					},
				},
				Critical: nil,
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Len(t, result.Rules[0].Thresholds, 1)
	assert.Contains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)
	assert.NotContains(t, result.Rules[0].Thresholds, restapi.CriticalSeverity)
}

func TestMapStateToDataObject_WithRulesWithOnlyCriticalThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		Rules: &LogAlertRuleModel{
			MetricName:        types.StringValue("log.count"),
			AlertType:         types.StringValue(LogAlertTypeLogCount),
			Aggregation:       types.StringValue(string(restapi.SumAggregation)),
			ThresholdOperator: types.StringValue(string(restapi.ThresholdOperatorGreaterThan)),
			Threshold: &ThresholdModel{
				Warning: nil,
				Critical: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Value(200),
					},
				},
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	require.Len(t, result.Rules[0].Thresholds, 1)
	assert.NotContains(t, result.Rules[0].Thresholds, restapi.WarningSeverity)
	assert.Contains(t, result.Rules[0].Thresholds, restapi.CriticalSeverity)
}

func TestMapStateToDataObject_WithRulesWithNullThresholdValues(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		Rules: &LogAlertRuleModel{
			MetricName:        types.StringValue("log.count"),
			AlertType:         types.StringValue(LogAlertTypeLogCount),
			Aggregation:       types.StringValue(string(restapi.SumAggregation)),
			ThresholdOperator: types.StringValue(string(restapi.ThresholdOperatorGreaterThan)),
			Threshold: &ThresholdModel{
				Warning: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Null(),
					},
				},
				Critical: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Null(),
					},
				},
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	assert.Empty(t, result.Rules[0].Thresholds)
}

func TestUpdateState_WithRulesWithOnlyWarningThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningValue := float64(100)

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{
			{
				Rule: restapi.LogAlertRule{
					AlertType:   "logCount",
					MetricName:  "log.count",
					Aggregation: restapi.SumAggregation,
				},
				ThresholdOperator: restapi.ThresholdOperatorGreaterThan,
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
				},
			},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.Threshold)
	require.NotNil(t, model.Rules.Threshold.Warning)
	assert.Nil(t, model.Rules.Threshold.Critical)
}

func TestUpdateState_WithRulesWithOnlyCriticalThreshold(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	criticalValue := float64(200)

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{
			{
				Rule: restapi.LogAlertRule{
					AlertType:   "logCount",
					MetricName:  "log.count",
					Aggregation: restapi.SumAggregation,
				},
				ThresholdOperator: restapi.ThresholdOperatorGreaterThan,
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.CriticalSeverity: {
						Type:  "staticThreshold",
						Value: &criticalValue,
					},
				},
			},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.Threshold)
	assert.Nil(t, model.Rules.Threshold.Warning)
	require.NotNil(t, model.Rules.Threshold.Critical)
}

func TestMapStateToDataObject_WithNullAggregation(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		Rules: &LogAlertRuleModel{
			MetricName:        types.StringValue("log.count"),
			AlertType:         types.StringValue(LogAlertTypeLogCount),
			Aggregation:       types.StringNull(),
			ThresholdOperator: types.StringValue(string(restapi.ThresholdOperatorGreaterThan)),
			Threshold: &ThresholdModel{
				Warning: &shared.ThresholdStaticTypeModel{
					Static: &shared.StaticTypeModel{
						Value: types.Int64Value(100),
					},
				},
			},
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)
	// When aggregation is null, it should not be set
	assert.Empty(t, result.Rules[0].Rule.Aggregation)
}

func TestMapStateToDataObject_WithUnknownAlertChannels(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Log Alert"),
		Description: types.StringValue("Test Description"),
		Granularity: types.Int64Value(600000),
		TagFilter:   types.StringValue("entity.type EQUALS 'log'"),
		AlertChannels: &AlertChannelsModel{
			Warning:  types.ListUnknown(types.StringType),
			Critical: types.ListUnknown(types.StringType),
		},
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.AlertChannels)
}

func TestMapStateToDataObject_WithUnknownGroupBy(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	model := LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	}

	diags := state.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.GroupBy)
}

func TestMapStateToDataObject_WithCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	// Create custom payload field objects
	field1, diags := types.ObjectValue(
		shared.CustomPayloadFieldAttributeTypes(),
		map[string]attr.Value{
			"key":           types.StringValue("field1"),
			"value":         types.StringValue("value1"),
			"dynamic_value": types.ObjectNull(shared.GetDynamicValueType().AttrTypes),
		},
	)
	require.False(t, diags.HasError())

	field2, diags := types.ObjectValue(
		shared.CustomPayloadFieldAttributeTypes(),
		map[string]attr.Value{
			"key":           types.StringValue("field2"),
			"value":         types.StringValue("value2"),
			"dynamic_value": types.ObjectNull(shared.GetDynamicValueType().AttrTypes),
		},
	)
	require.False(t, diags.HasError())

	customPayloadFieldsList, diags := types.ListValue(
		shared.GetCustomPayloadFieldType(),
		[]attr.Value{field1, field2},
	)
	require.False(t, diags.HasError())

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: customPayloadFieldsList,
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.CustomerPayloadFields, 2)
	assert.Equal(t, "field1", result.CustomerPayloadFields[0].Key)
	assert.Equal(t, "value1", result.CustomerPayloadFields[0].Value)
}

func TestUpdateState_WithMultipleRules(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	warningValue := float64(100)
	criticalValue := float64(200)

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{
			{
				Rule: restapi.LogAlertRule{
					AlertType:   "logCount",
					MetricName:  "log.count",
					Aggregation: restapi.SumAggregation,
				},
				ThresholdOperator: restapi.ThresholdOperatorGreaterThan,
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
					restapi.CriticalSeverity: {
						Type:  "staticThreshold",
						Value: &criticalValue,
					},
				},
			},
			// Second rule should be ignored since we only support single rule
			{
				Rule: restapi.LogAlertRule{
					AlertType:   "logCount",
					MetricName:  "log.count2",
					Aggregation: restapi.SumAggregation,
				},
				ThresholdOperator: restapi.ThresholdOperatorLessThan,
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &warningValue,
					},
				},
			},
		},
		AlertChannels: map[restapi.AlertSeverity][]string{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	// Should only map the first rule
	assert.Equal(t, "log.count", model.Rules.MetricName.ValueString())
}

func TestUpdateState_WithEmptyAlertChannelsMap(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		AlertChannels: map[restapi.AlertSeverity][]string{
			restapi.WarningSeverity:  {},
			restapi.CriticalSeverity: {},
		},
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Empty alert channel lists still create an AlertChannelsModel with null lists
	require.NotNil(t, model.AlertChannels)
	assert.True(t, model.AlertChannels.Warning.IsNull())
	assert.True(t, model.AlertChannels.Critical.IsNull())
}

func TestMapStateToDataObject_WithUnknownCustomPayloadFields(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue("entity.type EQUALS 'log'"),
		CustomPayloadFields: types.ListUnknown(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Empty(t, result.CustomerPayloadFields)
}

func TestMapStateToDataObject_WithNullTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringNull(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Nil(t, result.TagFilterExpression)
}

func TestMapStateToDataObject_WithEmptyTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	state := createMockState(t, LogAlertConfigModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Log Alert"),
		Description:         types.StringValue("Test Description"),
		Granularity:         types.Int64Value(600000),
		TagFilter:           types.StringValue(""),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	// Empty string tag filter will cause a parsing error
	assert.True(t, diags.HasError())
}

func TestUpdateState_WithEmptyWarningChannelList(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		AlertChannels: map[restapi.AlertSeverity][]string{
			restapi.WarningSeverity:  {},
			restapi.CriticalSeverity: {"channel-1"},
		},
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)
	assert.True(t, model.AlertChannels.Warning.IsNull())
	assert.False(t, model.AlertChannels.Critical.IsNull())
}

func TestUpdateState_WithEmptyCriticalChannelList(t *testing.T) {
	ctx := context.Background()
	resource := NewLogAlertConfigResourceHandleFramework()

	data := &restapi.LogAlertConfig{
		ID:          "test-id",
		Name:        "Test Log Alert",
		Description: "Test Description",
		Granularity: 600000,
		AlertChannels: map[restapi.AlertSeverity][]string{
			restapi.WarningSeverity:  {"channel-1"},
			restapi.CriticalSeverity: {},
		},
		Rules: []restapi.RuleWithThreshold[restapi.LogAlertRule]{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model LogAlertConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.AlertChannels)
	assert.False(t, model.AlertChannels.Warning.IsNull())
	assert.True(t, model.AlertChannels.Critical.IsNull())
}

// Helper functions

func createMockState(t *testing.T, model LogAlertConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewLogAlertConfigResourceHandleFramework()
	return resource.MetaData().Schema
}

// Made with Bob
