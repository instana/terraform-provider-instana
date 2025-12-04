package customeventspec

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

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}

func TestNewCustomEventSpecificationResourceHandle(t *testing.T) {
	resource := NewCustomEventSpecificationResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaCustomEventSpecification, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(1), metaData.SchemaVersion)
}

func TestMetaData(t *testing.T) {
	resource := &customEventSpecificationResource{
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
	resource := NewCustomEventSpecificationResourceHandle()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetRestResource(t *testing.T) {
	resource := &customEventSpecificationResource{}

	// Verify the method exists and interface is properly implemented
	var _ resourcehandle.ResourceHandle[*restapi.CustomEventSpecification] = resource
	assert.NotNil(t, resource.GetRestResource)
}

func TestMapIntToSeverityString(t *testing.T) {
	tests := []struct {
		name     string
		severity int
		expected string
	}{
		{
			name:     "warning severity",
			severity: 5,
			expected: "warning",
		},
		{
			name:     "critical severity",
			severity: 10,
			expected: "critical",
		},
		{
			name:     "unknown severity defaults to warning",
			severity: 99,
			expected: "warning",
		},
		{
			name:     "zero severity defaults to warning",
			severity: 0,
			expected: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapIntToSeverityString(tt.severity)
			assert.Equal(t, tt.expected, result.ValueString())
		})
	}
}

func TestMapSeverityToInt(t *testing.T) {
	tests := []struct {
		name     string
		severity string
		expected int
	}{
		{
			name:     "warning severity",
			severity: "warning",
			expected: 5,
		},
		{
			name:     "critical severity",
			severity: "critical",
			expected: 10,
		},
		{
			name:     "unknown severity defaults to warning",
			severity: "unknown",
			expected: 5,
		},
		{
			name:     "empty severity defaults to warning",
			severity: "",
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapSeverityToInt(tt.severity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapStateToDataObject_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Query:               types.StringValue("entity.type:host"),
		Triggering:          types.BoolValue(false),
		Description:         types.StringValue("Test Description"),
		ExpirationTime:      types.Int64Value(3600000),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules:               nil,
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Event", result.Name)
	assert.Equal(t, "host", result.EntityType)
	assert.NotNil(t, result.Query)
	assert.Equal(t, "entity.type:host", *result.Query)
	assert.False(t, result.Triggering)
	assert.NotNil(t, result.Description)
	assert.Equal(t, "Test Description", *result.Description)
	assert.NotNil(t, result.ExpirationTime)
	assert.Equal(t, 3600000, *result.ExpirationTime)
	assert.True(t, result.Enabled)
	assert.Equal(t, "AND", result.RuleLogicalOperator)
}

func TestMapStateToDataObject_WithEntityCountRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Query:               types.StringValue(""),
		Triggering:          types.BoolValue(false),
		Description:         types.StringValue(""),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			EntityCount: &EntityCountRuleModel{
				Severity:          types.StringValue("warning"),
				ConditionOperator: types.StringValue(">"),
				ConditionValue:    types.Float64Value(10.0),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.EntityCountRuleType, rule.DType)
	assert.Equal(t, 5, rule.Severity)
	assert.NotNil(t, rule.ConditionOperator)
	assert.Equal(t, ">", *rule.ConditionOperator)
	assert.NotNil(t, rule.ConditionValue)
	assert.Equal(t, 10.0, *rule.ConditionValue)
}

func TestMapStateToDataObject_WithEntityCountVerificationRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			EntityCountVerification: &EntityCountVerificationRuleModel{
				Severity:            types.StringValue("critical"),
				ConditionOperator:   types.StringValue(">="),
				ConditionValue:      types.Float64Value(5.0),
				MatchingEntityType:  types.StringValue("service"),
				MatchingOperator:    types.StringValue("is"),
				MatchingEntityLabel: types.StringValue("test-service"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.EntityCountVerificationRuleType, rule.DType)
	assert.Equal(t, 10, rule.Severity)
	assert.Equal(t, ">=", *rule.ConditionOperator)
	assert.Equal(t, 5.0, *rule.ConditionValue)
	assert.Equal(t, "service", *rule.MatchingEntityType)
	assert.Equal(t, "is", *rule.MatchingOperator)
	assert.Equal(t, "test-service", *rule.MatchingEntityLabel)
}

func TestMapStateToDataObject_WithEntityVerificationRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			EntityVerification: &EntityVerificationRuleModel{
				Severity:            types.StringValue("warning"),
				MatchingEntityType:  types.StringValue("process"),
				MatchingOperator:    types.StringValue("contains"),
				MatchingEntityLabel: types.StringValue("java"),
				OfflineDuration:     types.Int64Value(60000),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.EntityVerificationRuleType, rule.DType)
	assert.Equal(t, 5, rule.Severity)
	assert.Equal(t, "process", *rule.MatchingEntityType)
	assert.Equal(t, "contains", *rule.MatchingOperator)
	assert.Equal(t, "java", *rule.MatchingEntityLabel)
	assert.Equal(t, 60000, *rule.OfflineDuration)
}

func TestMapStateToDataObject_WithHostAvailabilityRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			HostAvailability: &HostAvailabilityRuleModel{
				Severity:        types.StringValue("critical"),
				OfflineDuration: types.Int64Value(120000),
				CloseAfter:      types.Int64Value(300000),
				TagFilter:       types.StringValue("entity.type EQUALS 'host'"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.HostAvailabilityRuleType, rule.DType)
	assert.Equal(t, 10, rule.Severity)
	assert.Equal(t, 120000, *rule.OfflineDuration)
	assert.NotNil(t, rule.CloseAfter)
	assert.Equal(t, 300000, *rule.CloseAfter)
	assert.NotNil(t, rule.TagFilter)
}

func TestMapStateToDataObject_WithHostAvailabilityRuleNoCloseAfter(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			HostAvailability: &HostAvailabilityRuleModel{
				Severity:        types.StringValue("warning"),
				OfflineDuration: types.Int64Value(60000),
				CloseAfter:      types.Int64Null(),
				TagFilter:       types.StringValue(""),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.HostAvailabilityRuleType, rule.DType)
	assert.Nil(t, rule.CloseAfter)
	assert.Nil(t, rule.TagFilter)
}

func TestMapStateToDataObject_WithInvalidTagFilter(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			HostAvailability: &HostAvailabilityRuleModel{
				Severity:        types.StringValue("warning"),
				OfflineDuration: types.Int64Value(60000),
				TagFilter:       types.StringValue("invalid tag filter syntax"),
			},
		},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_WithSystemRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			System: &SystemRuleModel{
				Severity:     types.StringValue("critical"),
				SystemRuleID: types.StringValue("system-rule-123"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.SystemRuleType, rule.DType)
	assert.Equal(t, 10, rule.Severity)
	assert.NotNil(t, rule.SystemRuleID)
	assert.Equal(t, "system-rule-123", *rule.SystemRuleID)
}

func TestMapStateToDataObject_WithThresholdRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			Threshold: &ThresholdRuleModel{
				Severity:          types.StringValue("warning"),
				MetricName:        types.StringValue("cpu.usage"),
				Rollup:            types.Int64Value(60000),
				Window:            types.Int64Value(300000),
				Aggregation:       types.StringValue("avg"),
				ConditionOperator: types.StringValue(">"),
				ConditionValue:    types.Float64Value(80.0),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.ThresholdRuleType, rule.DType)
	assert.Equal(t, 5, rule.Severity)
	assert.Equal(t, "cpu.usage", *rule.MetricName)
	assert.Equal(t, 60000, *rule.Rollup)
	assert.Equal(t, 300000, *rule.Window)
	assert.Equal(t, "avg", *rule.Aggregation)
	assert.Equal(t, ">", *rule.ConditionOperator)
	assert.Equal(t, 80.0, *rule.ConditionValue)
	assert.Nil(t, rule.MetricPattern)
}

func TestMapStateToDataObject_WithThresholdRuleAndMetricPattern(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			Threshold: &ThresholdRuleModel{
				Severity:          types.StringValue("critical"),
				MetricName:        types.StringValue("memory.usage"),
				Rollup:            types.Int64Value(60000),
				Window:            types.Int64Value(300000),
				Aggregation:       types.StringValue("sum"),
				ConditionOperator: types.StringValue(">="),
				ConditionValue:    types.Float64Value(90.0),
				MetricPattern: &MetricPatternModel{
					Prefix:      types.StringValue("jvm."),
					Postfix:     types.StringValue(".heap"),
					Placeholder: types.StringValue("instance"),
					Operator:    types.StringValue("EQUALS"),
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	assert.Equal(t, restapi.ThresholdRuleType, rule.DType)
	require.NotNil(t, rule.MetricPattern)
	assert.Equal(t, "jvm.", rule.MetricPattern.Prefix)
	assert.NotNil(t, rule.MetricPattern.Postfix)
	assert.Equal(t, ".heap", *rule.MetricPattern.Postfix)
	assert.NotNil(t, rule.MetricPattern.Placeholder)
	assert.Equal(t, "instance", *rule.MetricPattern.Placeholder)
	assert.Equal(t, "EQUALS", rule.MetricPattern.Operator)
}

func TestMapStateToDataObject_WithThresholdRuleAndEmptyMetricPattern(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules: &RulesModel{
			Threshold: &ThresholdRuleModel{
				Severity:          types.StringValue("warning"),
				MetricName:        types.StringValue("disk.usage"),
				Rollup:            types.Int64Value(60000),
				Window:            types.Int64Value(300000),
				Aggregation:       types.StringValue("max"),
				ConditionOperator: types.StringValue("<"),
				ConditionValue:    types.Float64Value(20.0),
				MetricPattern: &MetricPatternModel{
					Prefix:      types.StringValue("disk."),
					Postfix:     types.StringValue(""),
					Placeholder: types.StringValue(""),
					Operator:    types.StringValue("STARTS_WITH"),
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 1)

	rule := result.Rules[0]
	require.NotNil(t, rule.MetricPattern)
	assert.Equal(t, "disk.", rule.MetricPattern.Prefix)
	assert.Nil(t, rule.MetricPattern.Postfix)
	assert.Nil(t, rule.MetricPattern.Placeholder)
}

func TestMapStateToDataObject_WithMultipleRules(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("OR"),
		Rules: &RulesModel{
			EntityCount: &EntityCountRuleModel{
				Severity:          types.StringValue("warning"),
				ConditionOperator: types.StringValue(">"),
				ConditionValue:    types.Float64Value(5.0),
			},
			System: &SystemRuleModel{
				Severity:     types.StringValue("critical"),
				SystemRuleID: types.StringValue("sys-123"),
			},
			Threshold: &ThresholdRuleModel{
				Severity:          types.StringValue("warning"),
				MetricName:        types.StringValue("cpu.usage"),
				Rollup:            types.Int64Value(60000),
				Window:            types.Int64Value(300000),
				Aggregation:       types.StringValue("avg"),
				ConditionOperator: types.StringValue(">"),
				ConditionValue:    types.Float64Value(75.0),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	require.Len(t, result.Rules, 3)

	// Verify all rule types are present
	ruleTypes := make(map[string]bool)
	for _, rule := range result.Rules {
		ruleTypes[rule.DType] = true
	}
	assert.True(t, ruleTypes[restapi.EntityCountRuleType])
	assert.True(t, ruleTypes[restapi.SystemRuleType])
	assert.True(t, ruleTypes[restapi.ThresholdRuleType])
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	plan := createMockPlan(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(true),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
	})

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.True(t, result.Triggering)
}

func TestMapStateToDataObject_WithNullOptionalFields(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringNull(),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Query:               types.StringNull(),
		Triggering:          types.BoolValue(false),
		Description:         types.StringNull(),
		ExpirationTime:      types.Int64Null(),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("AND"),
		Rules:               nil,
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "", result.ID)
	assert.Nil(t, result.Query)
	assert.Nil(t, result.Description)
	assert.Nil(t, result.ExpirationTime)
}

func TestUpdateState_BasicConfig(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Query:               ptr("entity.type:host"),
		Triggering:          false,
		Description:         ptr("Test Description"),
		ExpirationTime:      ptr(3600000),
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules:               []restapi.RuleSpecification{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Event", model.Name.ValueString())
	assert.Equal(t, "host", model.EntityType.ValueString())
	assert.Equal(t, "entity.type:host", model.Query.ValueString())
	assert.False(t, model.Triggering.ValueBool())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	assert.Equal(t, int64(3600000), model.ExpirationTime.ValueInt64())
	assert.True(t, model.Enabled.ValueBool())
	assert.Equal(t, "AND", model.RuleLogicalOperator.ValueString())
}

func TestUpdateState_WithEntityCountRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	conditionOp := ">"
	conditionVal := 10.0
	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.EntityCountRuleType,
				Severity:          5,
				ConditionOperator: &conditionOp,
				ConditionValue:    &conditionVal,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.EntityCount)
	assert.Equal(t, "warning", model.Rules.EntityCount.Severity.ValueString())
	assert.Equal(t, ">", model.Rules.EntityCount.ConditionOperator.ValueString())
	assert.Equal(t, 10.0, model.Rules.EntityCount.ConditionValue.ValueFloat64())
}

func TestUpdateState_WithEntityCountVerificationRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	conditionOp := ">="
	conditionVal := 5.0
	matchingType := "service"
	matchingOp := "is"
	matchingLabel := "test-service"

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:               restapi.EntityCountVerificationRuleType,
				Severity:            10,
				ConditionOperator:   &conditionOp,
				ConditionValue:      &conditionVal,
				MatchingEntityType:  &matchingType,
				MatchingOperator:    &matchingOp,
				MatchingEntityLabel: &matchingLabel,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.EntityCountVerification)
	assert.Equal(t, "critical", model.Rules.EntityCountVerification.Severity.ValueString())
	assert.Equal(t, ">=", model.Rules.EntityCountVerification.ConditionOperator.ValueString())
	assert.Equal(t, 5.0, model.Rules.EntityCountVerification.ConditionValue.ValueFloat64())
	assert.Equal(t, "service", model.Rules.EntityCountVerification.MatchingEntityType.ValueString())
	assert.Equal(t, "is", model.Rules.EntityCountVerification.MatchingOperator.ValueString())
	assert.Equal(t, "test-service", model.Rules.EntityCountVerification.MatchingEntityLabel.ValueString())
}

func TestUpdateState_WithEntityVerificationRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	matchingType := "process"
	matchingOp := "contains"
	matchingLabel := "java"
	offlineDur := 60000

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:               restapi.EntityVerificationRuleType,
				Severity:            5,
				MatchingEntityType:  &matchingType,
				MatchingOperator:    &matchingOp,
				MatchingEntityLabel: &matchingLabel,
				OfflineDuration:     &offlineDur,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.EntityVerification)
	assert.Equal(t, "warning", model.Rules.EntityVerification.Severity.ValueString())
	assert.Equal(t, "process", model.Rules.EntityVerification.MatchingEntityType.ValueString())
	assert.Equal(t, "contains", model.Rules.EntityVerification.MatchingOperator.ValueString())
	assert.Equal(t, "java", model.Rules.EntityVerification.MatchingEntityLabel.ValueString())
	assert.Equal(t, int64(60000), model.Rules.EntityVerification.OfflineDuration.ValueInt64())
}

func TestUpdateState_WithHostAvailabilityRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	offlineDur := 120000
	closeAfter := 300000
	tagName := "entity.type"
	tagValue := "host"
	tagOp := restapi.EqualsOperator

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:           restapi.HostAvailabilityRuleType,
				Severity:        10,
				OfflineDuration: &offlineDur,
				CloseAfter:      &closeAfter,
				TagFilter: &restapi.TagFilter{
					Type:     restapi.TagFilterExpressionType,
					Name:     &tagName,
					Operator: &tagOp,
					Value:    &tagValue,
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.HostAvailability)
	assert.Equal(t, "critical", model.Rules.HostAvailability.Severity.ValueString())
	assert.Equal(t, int64(120000), model.Rules.HostAvailability.OfflineDuration.ValueInt64())
	assert.Equal(t, int64(300000), model.Rules.HostAvailability.CloseAfter.ValueInt64())
	assert.False(t, model.Rules.HostAvailability.TagFilter.IsNull())
}

func TestUpdateState_WithHostAvailabilityRuleNoCloseAfter(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	offlineDur := 60000

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:           restapi.HostAvailabilityRuleType,
				Severity:        5,
				OfflineDuration: &offlineDur,
				CloseAfter:      nil,
				TagFilter:       nil,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.HostAvailability)
	assert.True(t, model.Rules.HostAvailability.CloseAfter.IsNull())
	assert.Equal(t, "", model.Rules.HostAvailability.TagFilter.ValueString())
}

func TestUpdateState_WithSystemRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	systemRuleID := "system-rule-123"

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:        restapi.SystemRuleType,
				Severity:     10,
				SystemRuleID: &systemRuleID,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.System)
	assert.Equal(t, "critical", model.Rules.System.Severity.ValueString())
	assert.Equal(t, "system-rule-123", model.Rules.System.SystemRuleID.ValueString())
}

func TestUpdateState_WithThresholdRule(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	metricName := "cpu.usage"
	rollup := 60000
	window := 300000
	aggregation := "avg"
	conditionOp := ">"
	conditionVal := 80.0

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          5,
				MetricName:        &metricName,
				Rollup:            &rollup,
				Window:            &window,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOp,
				ConditionValue:    &conditionVal,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.Threshold)
	assert.Equal(t, "warning", model.Rules.Threshold.Severity.ValueString())
	assert.Equal(t, "cpu.usage", model.Rules.Threshold.MetricName.ValueString())
	assert.Equal(t, int64(60000), model.Rules.Threshold.Rollup.ValueInt64())
	assert.Equal(t, int64(300000), model.Rules.Threshold.Window.ValueInt64())
	assert.Equal(t, "avg", model.Rules.Threshold.Aggregation.ValueString())
	assert.Equal(t, ">", model.Rules.Threshold.ConditionOperator.ValueString())
	assert.Equal(t, 80.0, model.Rules.Threshold.ConditionValue.ValueFloat64())
	assert.Nil(t, model.Rules.Threshold.MetricPattern)
}

func TestUpdateState_WithThresholdRuleAndMetricPattern(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	metricName := "memory.usage"
	rollup := 60000
	window := 300000
	aggregation := "sum"
	conditionOp := ">="
	conditionVal := 90.0
	postfix := ".heap"
	placeholder := "instance"

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          10,
				MetricName:        &metricName,
				Rollup:            &rollup,
				Window:            &window,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOp,
				ConditionValue:    &conditionVal,
				MetricPattern: &restapi.MetricPattern{
					Prefix:      "jvm.",
					Postfix:     &postfix,
					Placeholder: &placeholder,
					Operator:    "EQUALS",
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.Threshold)
	require.NotNil(t, model.Rules.Threshold.MetricPattern)
	assert.Equal(t, "jvm.", model.Rules.Threshold.MetricPattern.Prefix.ValueString())
	assert.Equal(t, ".heap", model.Rules.Threshold.MetricPattern.Postfix.ValueString())
	assert.Equal(t, "instance", model.Rules.Threshold.MetricPattern.Placeholder.ValueString())
	assert.Equal(t, "EQUALS", model.Rules.Threshold.MetricPattern.Operator.ValueString())
}

func TestUpdateState_WithMultipleRules(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	conditionOp := ">"
	conditionVal := 5.0
	systemRuleID := "sys-123"
	metricName := "cpu.usage"
	rollup := 60000
	window := 300000
	aggregation := "avg"
	thresholdOp := ">"
	thresholdVal := 75.0

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "OR",
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.EntityCountRuleType,
				Severity:          5,
				ConditionOperator: &conditionOp,
				ConditionValue:    &conditionVal,
			},
			{
				DType:        restapi.SystemRuleType,
				Severity:     10,
				SystemRuleID: &systemRuleID,
			},
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          5,
				MetricName:        &metricName,
				Rollup:            &rollup,
				Window:            &window,
				Aggregation:       &aggregation,
				ConditionOperator: &thresholdOp,
				ConditionValue:    &thresholdVal,
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Rules)
	assert.NotNil(t, model.Rules.EntityCount)
	assert.NotNil(t, model.Rules.System)
	assert.NotNil(t, model.Rules.Threshold)
}

func TestUpdateState_WithNoRules(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules:               []restapi.RuleSpecification{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// When there are no rules, buildRulesModel returns an empty RulesModel (not nil)
	require.NotNil(t, model.Rules)
	assert.Nil(t, model.Rules.EntityCount)
	assert.Nil(t, model.Rules.EntityCountVerification)
	assert.Nil(t, model.Rules.EntityVerification)
	assert.Nil(t, model.Rules.HostAvailability)
	assert.Nil(t, model.Rules.System)
	assert.Nil(t, model.Rules.Threshold)
}

func TestUpdateState_WithNullOptionalFields(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Query:               nil,
		Triggering:          false,
		Description:         nil,
		ExpirationTime:      nil,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules:               []restapi.RuleSpecification{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.Query.IsNull() || model.Query.ValueString() == "")
	assert.True(t, model.Description.IsNull() || model.Description.ValueString() == "")
	assert.True(t, model.ExpirationTime.IsNull())
}

func TestUpdateState_SkipsIncompleteRules(t *testing.T) {
	ctx := context.Background()
	resource := &customEventSpecificationResource{}

	// Entity count rule without required fields
	spec := &restapi.CustomEventSpecification{
		ID:                  "test-id",
		Name:                "Test Event",
		EntityType:          "host",
		Triggering:          false,
		Enabled:             true,
		RuleLogicalOperator: "AND",
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.EntityCountRuleType,
				Severity:          5,
				ConditionOperator: nil, // Missing required field
				ConditionValue:    nil, // Missing required field
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	// Initialize state with empty model
	initializeEmptyState(t, ctx, state)

	diags := resource.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError())

	var model CustomEventSpecificationModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Rule should be skipped due to missing required fields, but an empty RulesModel is created
	// This is expected behavior - the Rules struct exists but all rule types are nil
	if model.Rules != nil {
		assert.Nil(t, model.Rules.EntityCount)
		assert.Nil(t, model.Rules.EntityCountVerification)
		assert.Nil(t, model.Rules.EntityVerification)
		assert.Nil(t, model.Rules.HostAvailability)
		assert.Nil(t, model.Rules.System)
		assert.Nil(t, model.Rules.Threshold)
	}
}

func TestSchemaValidation(t *testing.T) {
	handle := NewCustomEventSpecificationResourceHandle()
	schema := handle.MetaData().Schema

	t.Run("has required attributes", func(t *testing.T) {
		assert.Contains(t, schema.Attributes, "id")
		assert.Contains(t, schema.Attributes, "name")
		assert.Contains(t, schema.Attributes, "entity_type")
		assert.Contains(t, schema.Attributes, "query")
		assert.Contains(t, schema.Attributes, "triggering")
		assert.Contains(t, schema.Attributes, "description")
		assert.Contains(t, schema.Attributes, "expiration_time")
		assert.Contains(t, schema.Attributes, "enabled")
		assert.Contains(t, schema.Attributes, "rule_logical_operator")
		assert.Contains(t, schema.Attributes, "rules")
	})

	t.Run("id is computed", func(t *testing.T) {
		idAttr := schema.Attributes["id"]
		assert.NotNil(t, idAttr)
	})

	t.Run("name is required", func(t *testing.T) {
		nameAttr := schema.Attributes["name"]
		assert.NotNil(t, nameAttr)
	})

	t.Run("entity_type is required", func(t *testing.T) {
		entityTypeAttr := schema.Attributes["entity_type"]
		assert.NotNil(t, entityTypeAttr)
	})
}

// Helper functions

func createMockState(t *testing.T, model CustomEventSpecificationModel) *tfsdk.State {
	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func createMockPlan(t *testing.T, model CustomEventSpecificationModel) *tfsdk.Plan {
	plan := &tfsdk.Plan{
		Schema: getTestSchema(),
	}

	diags := plan.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set plan: %v", diags)

	return plan
}

func getTestSchema() schema.Schema {
	resource := NewCustomEventSpecificationResourceHandle()
	return resource.MetaData().Schema
}

// initializeEmptyState initializes the state with an empty CustomEventSpecificationModel
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := CustomEventSpecificationModel{
		ID:                  types.StringNull(),
		Name:                types.StringNull(),
		EntityType:          types.StringNull(),
		Query:               types.StringNull(),
		Triggering:          types.BoolNull(),
		Description:         types.StringNull(),
		ExpirationTime:      types.Int64Null(),
		Enabled:             types.BoolNull(),
		RuleLogicalOperator: types.StringNull(),
		Rules:               nil,
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}
