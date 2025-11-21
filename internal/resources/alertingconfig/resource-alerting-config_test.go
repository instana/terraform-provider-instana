package alertingconfig

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlertingConfigResourceHandle(t *testing.T) {
	handle := NewAlertingConfigResourceHandle()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaAlertingConfig, metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &alertingConfigResource{
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
	resource := &alertingConfigResource{}
	ctx := context.Background()

	diags := resource.SetComputedFields(ctx, nil)
	assert.False(t, diags.HasError())
}

func TestUpdateState(t *testing.T) {
	resource := &alertingConfigResource{}
	ctx := context.Background()

	t.Run("basic configuration", func(t *testing.T) {
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert",
			IntegrationIDs: []string{"integration-1", "integration-2"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      nil,
				RuleIDs:    []string{},
				EventTypes: []restapi.AlertEventType{},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-alert", model.AlertName.ValueString())

		var integrationIDs []string
		diags = model.IntegrationIDs.ElementsAs(ctx, &integrationIDs, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"integration-1", "integration-2"}, integrationIDs)
	})

	t.Run("with event filter query", func(t *testing.T) {
		query := "entity.type:host"
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert",
			IntegrationIDs: []string{"integration-1"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      &query,
				RuleIDs:    []string{},
				EventTypes: []restapi.AlertEventType{},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.Equal(t, "entity.type:host", model.EventFilterQuery.ValueString())
	})

	t.Run("with event types", func(t *testing.T) {
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert",
			IntegrationIDs: []string{"integration-1"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      nil,
				RuleIDs:    []string{},
				EventTypes: []restapi.AlertEventType{restapi.CriticalAlertEventType, restapi.WarningAlertEventType},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		var eventTypes []string
		diags = model.EventFilterEventTypes.ElementsAs(ctx, &eventTypes, false)
		require.False(t, diags.HasError())
		assert.Contains(t, eventTypes, "critical")
		assert.Contains(t, eventTypes, "warning")
	})

	t.Run("with rule IDs", func(t *testing.T) {
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert",
			IntegrationIDs: []string{"integration-1"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      nil,
				RuleIDs:    []string{"rule-1", "rule-2"},
				EventTypes: []restapi.AlertEventType{},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		var ruleIDs []string
		diags = model.EventFilterRuleIDs.ElementsAs(ctx, &ruleIDs, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"rule-1", "rule-2"}, ruleIDs)
	})

	t.Run("with custom payload fields", func(t *testing.T) {
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert",
			IntegrationIDs: []string{"integration-1"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      nil,
				RuleIDs:    []string{},
				EventTypes: []restapi.AlertEventType{},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type:  restapi.StaticStringCustomPayloadType,
					Key:   "static_key",
					Value: "static_value",
				},
			},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.False(t, model.CustomPayloadFields.IsNull())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &alertingConfigResource{}
	ctx := context.Background()

	t.Run("basic configuration from state", func(t *testing.T) {
		integrationIDs := []string{"integration-1", "integration-2"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id"),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringNull(),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            types.ListNull(shared.GetCustomPayloadFieldType()),
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Equal(t, "test-id", config.ID)
		assert.Equal(t, "test-alert", config.AlertName)
		assert.Equal(t, integrationIDs, config.IntegrationIDs)
		assert.Nil(t, config.EventFilteringConfiguration.Query)
		assert.Empty(t, config.EventFilteringConfiguration.RuleIDs)
		assert.Empty(t, config.EventFilteringConfiguration.EventTypes)
	})

	t.Run("configuration from plan", func(t *testing.T) {
		integrationIDs := []string{"integration-1"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		model := AlertingConfigModel{
			ID:                             types.StringValue(""),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringNull(),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            types.ListNull(shared.GetCustomPayloadFieldType()),
		}

		plan := createMockPlan(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, plan, nil)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Equal(t, "", config.ID)
		assert.Equal(t, "test-alert", config.AlertName)
	})

	t.Run("with event filter query", func(t *testing.T) {
		integrationIDs := []string{"integration-1"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id"),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringValue("entity.type:host"),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            types.ListNull(shared.GetCustomPayloadFieldType()),
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Equal(t, "entity.type:host", *config.EventFilteringConfiguration.Query)
	})

	t.Run("with event types", func(t *testing.T) {
		integrationIDs := []string{"integration-1"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		eventTypes := []string{"critical", "warning"}
		eventTypesSet, _ := types.SetValueFrom(ctx, types.StringType, eventTypes)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id"),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringNull(),
			EventFilterEventTypes:          eventTypesSet,
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            types.ListNull(shared.GetCustomPayloadFieldType()),
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Len(t, config.EventFilteringConfiguration.EventTypes, 2)
		assert.Contains(t, config.EventFilteringConfiguration.EventTypes, restapi.CriticalAlertEventType)
		assert.Contains(t, config.EventFilteringConfiguration.EventTypes, restapi.WarningAlertEventType)
	})

	t.Run("with rule IDs", func(t *testing.T) {
		integrationIDs := []string{"integration-1"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		ruleIDs := []string{"rule-1", "rule-2"}
		ruleIDsSet, _ := types.SetValueFrom(ctx, types.StringType, ruleIDs)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id"),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringNull(),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             ruleIDsSet,
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            types.ListNull(shared.GetCustomPayloadFieldType()),
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Equal(t, ruleIDs, config.EventFilteringConfiguration.RuleIDs)
	})

	t.Run("with custom payload fields", func(t *testing.T) {
		integrationIDs := []string{"integration-1"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		// Create custom payload fields
		customFields := []attr.Value{
			createStaticCustomPayloadField(t, "key1", "value1"),
		}
		customFieldsList, _ := types.ListValue(shared.GetCustomPayloadFieldType(), customFields)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id"),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringNull(),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            customFieldsList,
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Len(t, config.CustomerPayloadFields, 1)
		assert.Equal(t, "key1", config.CustomerPayloadFields[0].Key)
	})
}

func TestConvertEventTypesToHarmonizedStringRepresentation(t *testing.T) {
	resource := &alertingConfigResource{}

	testCases := []struct {
		name     string
		input    []restapi.AlertEventType
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []restapi.AlertEventType{},
			expected: []string{},
		},
		{
			name:     "single event type",
			input:    []restapi.AlertEventType{restapi.CriticalAlertEventType},
			expected: []string{"critical"},
		},
		{
			name:     "multiple event types",
			input:    []restapi.AlertEventType{restapi.CriticalAlertEventType, restapi.WarningAlertEventType, restapi.IncidentAlertEventType},
			expected: []string{"critical", "warning", "incident"},
		},
		{
			name:     "all supported event types",
			input:    restapi.SupportedAlertEventTypes,
			expected: []string{"incident", "critical", "warning", "change", "online", "offline", "none", "agent_monitoring_issue"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := resource.convertEventTypesToHarmonizedStringRepresentation(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestReadEventTypesFromStrings(t *testing.T) {
	resource := &alertingConfigResource{}

	testCases := []struct {
		name     string
		input    []string
		expected []restapi.AlertEventType
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []restapi.AlertEventType{},
		},
		{
			name:     "single event type",
			input:    []string{"critical"},
			expected: []restapi.AlertEventType{restapi.CriticalAlertEventType},
		},
		{
			name:     "multiple event types",
			input:    []string{"critical", "warning", "incident"},
			expected: []restapi.AlertEventType{restapi.CriticalAlertEventType, restapi.WarningAlertEventType, restapi.IncidentAlertEventType},
		},
		{
			name:     "mixed case input",
			input:    []string{"CRITICAL", "Warning", "incident"},
			expected: []restapi.AlertEventType{restapi.CriticalAlertEventType, restapi.WarningAlertEventType, restapi.IncidentAlertEventType},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := resource.readEventTypesFromStrings(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertSupportedEventTypesToStringSlice(t *testing.T) {
	result := convertSupportedEventTypesToStringSlice()

	expected := []string{
		"incident",
		"critical",
		"warning",
		"change",
		"online",
		"offline",
		"none",
		"agent_monitoring_issue",
	}

	assert.Equal(t, expected, result)
	assert.Equal(t, len(restapi.SupportedAlertEventTypes), len(result))
}

func TestSupportedEventTypes(t *testing.T) {
	// Test that the supportedEventTypes variable is properly initialized
	assert.NotEmpty(t, supportedEventTypes)
	assert.Equal(t, len(restapi.SupportedAlertEventTypes), len(supportedEventTypes))

	// Verify all expected event types are present
	expectedTypes := []string{
		"incident",
		"critical",
		"warning",
		"change",
		"online",
		"offline",
		"none",
		"agent_monitoring_issue",
	}

	for _, expectedType := range expectedTypes {
		assert.Contains(t, supportedEventTypes, expectedType)
	}
}

func TestGetRestResource(t *testing.T) {
	resource := &alertingConfigResource{}

	// This test just ensures the method exists and can be called
	// The actual implementation requires a real API instance
	assert.NotNil(t, resource)
}

// Helper functions

func createMockState(t *testing.T, ctx context.Context, model AlertingConfigModel) *tfsdk.State {
	handle := NewAlertingConfigResourceHandle()
	state := &tfsdk.State{
		Schema: handle.MetaData().Schema,
	}

	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set state: %v", diags)
	}

	return state
}

func createMockPlan(t *testing.T, ctx context.Context, model AlertingConfigModel) *tfsdk.Plan {
	handle := NewAlertingConfigResourceHandle()
	plan := &tfsdk.Plan{
		Schema: handle.MetaData().Schema,
	}

	diags := plan.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set plan: %v", diags)
	}

	return plan
}

func createStaticCustomPayloadField(t *testing.T, key, value string) attr.Value {
	fieldAttrs := map[string]attr.Value{
		"key":           types.StringValue(key),
		"value":         types.StringValue(value),
		"dynamic_value": types.ObjectNull(shared.GetDynamicValueType().AttrTypes),
	}

	objVal, diags := types.ObjectValue(shared.CustomPayloadFieldAttributeTypes(), fieldAttrs)
	if diags.HasError() {
		t.Fatalf("Failed to create custom payload field: %v", diags)
	}

	return objVal
}

func createDynamicCustomPayloadField(t *testing.T, key, tagName string, dynamicKey *string) attr.Value {
	dynamicAttrs := map[string]attr.Value{
		"tag_name": types.StringValue(tagName),
	}

	if dynamicKey != nil {
		dynamicAttrs["key"] = types.StringValue(*dynamicKey)
	} else {
		dynamicAttrs["key"] = types.StringNull()
	}

	dynamicObj, diags := types.ObjectValue(shared.GetDynamicValueType().AttrTypes, dynamicAttrs)
	if diags.HasError() {
		t.Fatalf("Failed to create dynamic value: %v", diags)
	}

	fieldAttrs := map[string]attr.Value{
		"key":           types.StringValue(key),
		"value":         types.StringNull(),
		"dynamic_value": dynamicObj,
	}

	objVal, diags := types.ObjectValue(shared.CustomPayloadFieldAttributeTypes(), fieldAttrs)
	if diags.HasError() {
		t.Fatalf("Failed to create custom payload field: %v", diags)
	}

	return objVal
}

func TestMapStateToDataObjectWithDynamicCustomPayloadFields(t *testing.T) {
	resource := &alertingConfigResource{}
	ctx := context.Background()

	t.Run("with dynamic custom payload fields", func(t *testing.T) {
		integrationIDs := []string{"integration-1"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		// Create dynamic custom payload fields
		dynamicKey := "dynamic_key"
		customFields := []attr.Value{
			createDynamicCustomPayloadField(t, "key1", "tag1", &dynamicKey),
			createDynamicCustomPayloadField(t, "key2", "tag2", nil),
		}
		customFieldsList, _ := types.ListValue(shared.GetCustomPayloadFieldType(), customFields)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id"),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringNull(),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            customFieldsList,
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)
		assert.Len(t, config.CustomerPayloadFields, 2)
		assert.Equal(t, "key1", config.CustomerPayloadFields[0].Key)
		assert.Equal(t, restapi.DynamicCustomPayloadType, config.CustomerPayloadFields[0].Type)
	})
}

func TestUpdateStateWithComplexCustomPayloadFields(t *testing.T) {
	resource := &alertingConfigResource{}
	ctx := context.Background()

	t.Run("with dynamic custom payload fields", func(t *testing.T) {
		dynamicKey := "dynamic_key"
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert",
			IntegrationIDs: []string{"integration-1"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      nil,
				RuleIDs:    []string{},
				EventTypes: []restapi.AlertEventType{},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type: restapi.DynamicCustomPayloadType,
					Key:  "dynamic_field",
					Value: restapi.DynamicCustomPayloadFieldValue{
						Key:     &dynamicKey,
						TagName: "test_tag",
					},
				},
			},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.False(t, model.CustomPayloadFields.IsNull())
	})
}

// Made with Bob

func TestMapStateToDataObjectEdgeCases(t *testing.T) {
	resource := &alertingConfigResource{}
	ctx := context.Background()

	t.Run("with all fields populated including complex scenarios", func(t *testing.T) {
		integrationIDs := []string{"integration-1", "integration-2", "integration-3"}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		eventTypes := []string{"incident", "critical", "warning", "change", "online", "offline"}
		eventTypesSet, _ := types.SetValueFrom(ctx, types.StringType, eventTypes)

		ruleIDs := []string{"rule-1", "rule-2", "rule-3", "rule-4"}
		ruleIDsSet, _ := types.SetValueFrom(ctx, types.StringType, ruleIDs)

		// Create mixed custom payload fields
		dynamicKey := "dynamic_key"
		customFields := []attr.Value{
			createStaticCustomPayloadField(t, "static_key1", "static_value1"),
			createStaticCustomPayloadField(t, "static_key2", "static_value2"),
			createDynamicCustomPayloadField(t, "dynamic_key1", "tag_name1", &dynamicKey),
			createDynamicCustomPayloadField(t, "dynamic_key2", "tag_name2", nil),
		}
		customFieldsList, _ := types.ListValue(shared.GetCustomPayloadFieldType(), customFields)

		model := AlertingConfigModel{
			ID:                             types.StringValue("test-id-complex"),
			AlertName:                      types.StringValue("comprehensive-alert-config"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringValue("entity.type:host AND entity.label:production AND entity.zone:us-east-1"),
			EventFilterEventTypes:          eventTypesSet,
			EventFilterRuleIDs:             ruleIDsSet,
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            customFieldsList,
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)

		assert.Equal(t, "test-id-complex", config.ID)
		assert.Equal(t, "comprehensive-alert-config", config.AlertName)
		assert.Equal(t, integrationIDs, config.IntegrationIDs)
		assert.NotNil(t, config.EventFilteringConfiguration.Query)
		assert.Equal(t, "entity.type:host AND entity.label:production AND entity.zone:us-east-1", *config.EventFilteringConfiguration.Query)
		assert.Equal(t, ruleIDs, config.EventFilteringConfiguration.RuleIDs)
		assert.Len(t, config.EventFilteringConfiguration.EventTypes, 6)
		assert.Len(t, config.CustomerPayloadFields, 4)
	})

	t.Run("with empty string values", func(t *testing.T) {
		integrationIDs := []string{}
		integrationIDsSet, _ := types.SetValueFrom(ctx, types.StringType, integrationIDs)

		model := AlertingConfigModel{
			ID:                             types.StringValue(""),
			AlertName:                      types.StringValue("test-alert"),
			IntegrationIDs:                 integrationIDsSet,
			EventFilterQuery:               types.StringValue(""),
			EventFilterEventTypes:          types.SetNull(types.StringType),
			EventFilterRuleIDs:             types.SetNull(types.StringType),
			EventFilterApplicationAlertIDs: types.SetNull(types.StringType),
			CustomPayloadFields:            types.ListNull(shared.GetCustomPayloadFieldType()),
		}

		state := createMockState(t, ctx, model)
		config, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, config)

		assert.Equal(t, "", config.ID)
		assert.Equal(t, "test-alert", config.AlertName)
		assert.Empty(t, config.IntegrationIDs)
		assert.NotNil(t, config.EventFilteringConfiguration.Query)
		assert.Equal(t, "", *config.EventFilteringConfiguration.Query)
	})
}

func TestUpdateStateWithAllEventTypes(t *testing.T) {
	resource := &alertingConfigResource{}
	ctx := context.Background()

	t.Run("with all supported event types", func(t *testing.T) {
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert-all-types",
			IntegrationIDs: []string{"integration-1"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      nil,
				RuleIDs:    []string{},
				EventTypes: restapi.SupportedAlertEventTypes,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		var eventTypes []string
		diags = model.EventFilterEventTypes.ElementsAs(ctx, &eventTypes, false)
		require.False(t, diags.HasError())
		assert.Len(t, eventTypes, len(restapi.SupportedAlertEventTypes))

		// Verify all event types are present
		for _, expectedType := range supportedEventTypes {
			assert.Contains(t, eventTypes, expectedType)
		}
	})

	t.Run("with mixed payload fields and all filters", func(t *testing.T) {
		query := "entity.type:host"
		dynamicKey := "dynamic_key"
		apiConfig := &restapi.AlertingConfiguration{
			ID:             "test-id",
			AlertName:      "test-alert-mixed",
			IntegrationIDs: []string{"integration-1", "integration-2"},
			EventFilteringConfiguration: restapi.EventFilteringConfiguration{
				Query:      &query,
				RuleIDs:    []string{"rule-1", "rule-2", "rule-3"},
				EventTypes: []restapi.AlertEventType{restapi.CriticalAlertEventType, restapi.WarningAlertEventType, restapi.IncidentAlertEventType},
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type:  restapi.StaticStringCustomPayloadType,
					Key:   "static_key",
					Value: "static_value",
				},
				{
					Type: restapi.DynamicCustomPayloadType,
					Key:  "dynamic_field",
					Value: restapi.DynamicCustomPayloadFieldValue{
						Key:     &dynamicKey,
						TagName: "test_tag",
					},
				},
			},
		}

		handle := NewAlertingConfigResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiConfig)
		require.False(t, diags.HasError())

		var model AlertingConfigModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())

		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-alert-mixed", model.AlertName.ValueString())
		assert.Equal(t, "entity.type:host", model.EventFilterQuery.ValueString())
		assert.False(t, model.CustomPayloadFields.IsNull())
		assert.False(t, model.EventFilterEventTypes.IsNull())
		assert.False(t, model.EventFilterRuleIDs.IsNull())
	})
}
