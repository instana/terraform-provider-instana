package customeventspec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/instana/instana-go-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// upgradeRawJSON runs the state upgrader on the given raw state and returns the upgraded model
func upgradeRawJSON(t *testing.T, version int64, rawJSON string) CustomEventSpecificationModel {
	t.Helper()
	ctx := context.Background()
	handle := &customEventSpecificationResource{}

	upgrader, ok := handle.GetStateUpgraders(ctx)[version]
	require.True(t, ok, "no state upgrader registered for version %d", version)

	resp := &resource.UpgradeStateResponse{}
	upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{RawState: &tfprotov6.RawState{JSON: []byte(rawJSON)}}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.NotNil(t, resp.DynamicValue)

	schema := createCustomEventSpecificationSchema()
	value, err := resp.DynamicValue.Unmarshal(schema.Type().TerraformType(ctx))
	require.NoError(t, err)

	state := tfsdk.State{Schema: schema, Raw: value}
	var model CustomEventSpecificationModel
	diags := state.Get(ctx, &model)
	require.False(t, diags.HasError(), "%v", diags)
	return model
}

func TestStateUpgrader_Version0_ThresholdListWithMetricPattern(t *testing.T) {
	model := upgradeRawJSON(t, 0, `{
		"id": "abc",
		"name": "multi",
		"entity_type": "host",
		"triggering": false,
		"enabled": true,
		"rule_logical_operator": "OR",
		"rules": [{
			"entity_count": [],
			"entity_count_verification": [],
			"entity_verification": [],
			"host_availability": [],
			"system": [],
			"threshold": [
				{"severity": "warning", "metric_name": "cpu.used", "metric_pattern": [], "rollup": 0, "window": 60000,
				 "aggregation": "avg", "condition_operator": ">", "condition_value": 80},
				{"severity": "warning", "metric_name": "", "rollup": 0, "window": 60000, "aggregation": "max",
				 "condition_operator": ">=", "condition_value": 10.5,
				 "metric_pattern": [{"prefix": "disk", "postfix": null, "placeholder": "sda", "operator": "is"}]}
			]
		}]
	}`)

	assert.Equal(t, "abc", model.ID.ValueString())
	assert.Equal(t, "OR", model.RuleLogicalOperator.ValueString())
	require.NotNil(t, model.Rules)
	assert.Nil(t, model.Rules.EntityCount)
	assert.Nil(t, model.Rules.System)
	require.Len(t, model.Rules.Threshold, 2)

	first := model.Rules.Threshold[0]
	assert.Equal(t, "cpu.used", first.MetricName.ValueString())
	assert.Equal(t, int64(0), first.Rollup.ValueInt64())
	assert.Nil(t, first.MetricPattern)

	second := model.Rules.Threshold[1]
	assert.Equal(t, 10.5, second.ConditionValue.ValueFloat64())
	require.NotNil(t, second.MetricPattern)
	assert.Equal(t, "disk", second.MetricPattern.Prefix.ValueString())
	assert.Equal(t, "sda", second.MetricPattern.Placeholder.ValueString())
	assert.Equal(t, "is", second.MetricPattern.Operator.ValueString())
	assert.Equal(t, "", second.MetricPattern.Postfix.ValueString())
}

func TestStateUpgrader_Version0_NonThresholdRule(t *testing.T) {
	model := upgradeRawJSON(t, 0, `{
		"id": "abc", "name": "sys", "entity_type": "any", "triggering": true, "enabled": true,
		"rule_logical_operator": "AND",
		"rules": [{
			"entity_count": [], "entity_count_verification": [], "entity_verification": [],
			"host_availability": [], "threshold": [],
			"system": [{"severity": "critical", "system_rule_id": "sys-1"}]
		}]
	}`)

	require.NotNil(t, model.Rules)
	require.NotNil(t, model.Rules.System)
	assert.Equal(t, "critical", model.Rules.System.Severity.ValueString())
	assert.Equal(t, "sys-1", model.Rules.System.SystemRuleID.ValueString())
	assert.Empty(t, model.Rules.Threshold)
}

func TestStateUpgrader_Version1_ThresholdObjectToList(t *testing.T) {
	model := upgradeRawJSON(t, 1, `{
		"id": "abc", "name": "single", "entity_type": "host", "triggering": false, "enabled": true,
		"rule_logical_operator": "AND",
		"rules": {
			"entity_count": null, "entity_count_verification": null, "entity_verification": null,
			"host_availability": null, "system": null,
			"threshold": {"severity": "critical", "metric_name": "cpu.used", "metric_pattern": null,
				"rollup": 1000, "window": 60000, "aggregation": "sum", "condition_operator": "<", "condition_value": 5}
		}
	}`)

	require.NotNil(t, model.Rules)
	require.Len(t, model.Rules.Threshold, 1)
	assert.Equal(t, "critical", model.Rules.Threshold[0].Severity.ValueString())
	assert.Equal(t, int64(1000), model.Rules.Threshold[0].Rollup.ValueInt64())
	assert.Nil(t, model.Rules.Threshold[0].MetricPattern)
}

func TestStateUpgrader_Version1_NoRules(t *testing.T) {
	model := upgradeRawJSON(t, 1, `{
		"id": "abc", "name": "none", "entity_type": "host", "triggering": false, "enabled": true,
		"rule_logical_operator": "AND", "rules": null
	}`)

	assert.Nil(t, model.Rules)
	assert.Equal(t, types.StringValue("abc"), model.ID)
}

func TestStateUpgrader_MissingRawState(t *testing.T) {
	ctx := context.Background()
	handle := &customEventSpecificationResource{}

	resp := &resource.UpgradeStateResponse{}
	handle.GetStateUpgraders(ctx)[0].StateUpgrader(ctx, resource.UpgradeStateRequest{}, resp)
	assert.True(t, resp.Diagnostics.HasError())
}

func TestMapStateToDataObject_WithMultipleThresholdRules(t *testing.T) {
	ctx := context.Background()
	handle := &customEventSpecificationResource{}

	threshold := func(metric string) ThresholdRuleModel {
		return ThresholdRuleModel{
			Severity:          types.StringValue("warning"),
			MetricName:        types.StringValue(metric),
			Rollup:            types.Int64Value(0),
			Window:            types.Int64Value(60000),
			Aggregation:       types.StringValue("avg"),
			ConditionOperator: types.StringValue(">"),
			ConditionValue:    types.Float64Value(1),
		}
	}

	state := createMockState(t, CustomEventSpecificationModel{
		ID:                  types.StringValue("test-id"),
		Name:                types.StringValue("Test Event"),
		EntityType:          types.StringValue("host"),
		Triggering:          types.BoolValue(false),
		Enabled:             types.BoolValue(true),
		RuleLogicalOperator: types.StringValue("OR"),
		Rules: &RulesModel{
			Threshold: []ThresholdRuleModel{threshold("cpu.used"), threshold("memory.used")},
		},
	})

	result, diags := handle.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.Len(t, result.Rules, 2)
	assert.Equal(t, "cpu.used", *result.Rules[0].MetricName)
	assert.Equal(t, "memory.used", *result.Rules[1].MetricName)
	assert.Equal(t, "OR", result.RuleLogicalOperator)
}

func TestUpdateState_WithMultipleThresholdRules(t *testing.T) {
	ctx := context.Background()
	handle := &customEventSpecificationResource{}

	rollup, window := 0, 60000
	aggregation, operator, value := "avg", ">", 1.0
	metrics := []string{"cpu.used", "memory.used"}
	spec := &api.CustomEventSpecification{
		ID: "test-id", Name: "Test Event", EntityType: "host", Enabled: true, RuleLogicalOperator: "AND",
	}
	for i := range metrics {
		spec.Rules = append(spec.Rules, api.RuleSpecification{
			DType: api.ThresholdRuleType, Severity: 5, MetricName: &metrics[i], Rollup: &rollup, Window: &window,
			Aggregation: &aggregation, ConditionOperator: &operator, ConditionValue: &value,
		})
	}

	state := &tfsdk.State{Schema: getTestSchema()}
	initializeEmptyState(t, ctx, state)
	diags := handle.UpdateState(ctx, state, nil, spec)
	require.False(t, diags.HasError(), "%v", diags)

	var model CustomEventSpecificationModel
	require.False(t, state.Get(ctx, &model).HasError())
	require.NotNil(t, model.Rules)
	require.Len(t, model.Rules.Threshold, 2)
	assert.Equal(t, "cpu.used", model.Rules.Threshold[0].MetricName.ValueString())
	assert.Equal(t, "memory.used", model.Rules.Threshold[1].MetricName.ValueString())
}
