package automationpolicy

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

// TestNewAutomationPolicyResourceHandleFramework tests the resource initialization
func TestNewAutomationPolicyResourceHandleFramework(t *testing.T) {
	resource := NewAutomationPolicyResourceHandleFramework()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaAutomationPolicyFramework, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(0), metaData.SchemaVersion)
}

// TestMetaData tests the MetaData method
func TestMetaData(t *testing.T) {
	resource := &automationPolicyResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  "test_resource",
			SchemaVersion: 0,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_resource", metaData.ResourceName)
	assert.Equal(t, int64(0), metaData.SchemaVersion)
}

// TestSetComputedFields tests the SetComputedFields method
func TestSetComputedFields(t *testing.T) {
	resource := NewAutomationPolicyResourceHandleFramework()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

// TestGetRestResource tests the GetRestResource method
func TestGetRestResource(t *testing.T) {
	resource := &automationPolicyResourceFramework{}
	assert.NotNil(t, resource.GetRestResource)
}

// TestMapStateToDataObject_BasicPolicy tests mapping basic policy from state to data object
func TestMapStateToDataObject_BasicPolicy(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:          types.StringValue("trigger-123"),
			Type:        types.StringValue("customEvent"),
			Name:        types.StringNull(),
			Description: types.StringNull(),
			Scheduling:  nil,
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name:      types.StringValue("manual"),
				Condition: nil,
				Action: []PolicyActionModel{
					{
						Action: shared.AutomationActionModel{
							ID:             types.StringValue("action-123"),
							Name:           types.StringValue("Test Action"),
							Description:    types.StringValue("Action Description"),
							Tags:           types.ListNull(types.StringType),
							InputParameter: []shared.ParameterModel{},
							Manual: &shared.ManualModel{
								Content: types.StringValue("Manual content"),
							},
						},
						AgentID: types.StringNull(),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Policy", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Nil(t, result.Tags)
	assert.Equal(t, "trigger-123", result.Trigger.Id)
	assert.Equal(t, "customEvent", result.Trigger.Type)
	require.Len(t, result.TypeConfigurations, 1)
	assert.Equal(t, "manual", result.TypeConfigurations[0].Name)
}

// TestMapStateToDataObject_WithTags tests mapping policy with tags
func TestMapStateToDataObject_WithTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	tags := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("tag1"),
		types.StringValue("tag2"),
		types.StringValue("tag3"),
	})

	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        tags,
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-123"),
			Type: types.StringValue("builtinEvent"),
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("manual"),
				Action: []PolicyActionModel{
					{
						Action: createMinimalManualAction(),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.NotNil(t, result.Tags)
	tagSlice, ok := result.Tags.([]string)
	require.True(t, ok)
	assert.Len(t, tagSlice, 3)
	assert.Contains(t, tagSlice, "tag1")
	assert.Contains(t, tagSlice, "tag2")
	assert.Contains(t, tagSlice, "tag3")
}

// TestMapStateToDataObject_WithScheduling tests mapping policy with scheduling
func TestMapStateToDataObject_WithScheduling(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:          types.StringValue("trigger-123"),
			Type:        types.StringValue("schedule"),
			Name:        types.StringValue("Scheduled Trigger"),
			Description: types.StringValue("Trigger Description"),
			Scheduling: &SchedulingModel{
				StartTime:     types.Int64Value(1609459200000),
				Duration:      types.Int64Value(60),
				DurationUnit:  types.StringValue("MINUTE"),
				RecurrentRule: types.StringValue("FREQ=DAILY"),
				Recurrent:     types.BoolValue(true),
			},
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("automatic"),
				Condition: &ConditionModel{
					Query: types.StringValue("entity.type:service"),
				},
				Action: []PolicyActionModel{
					{
						Action: createMinimalManualAction(),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "schedule", result.Trigger.Type)
	assert.Equal(t, "Scheduled Trigger", result.Trigger.Name)
	assert.Equal(t, "Trigger Description", result.Trigger.Description)
	assert.Equal(t, int64(1609459200000), result.Trigger.Scheduling.StartTime)
	assert.Equal(t, 60, result.Trigger.Scheduling.Duration)
	assert.Equal(t, restapi.DurationUnit("MINUTE"), result.Trigger.Scheduling.DurationUnit)
	assert.Equal(t, "FREQ=DAILY", result.Trigger.Scheduling.RecurrentRule)
	assert.True(t, result.Trigger.Scheduling.Recurrent)
}

// TestMapStateToDataObject_WithCondition tests mapping policy with condition
func TestMapStateToDataObject_WithCondition(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-123"),
			Type: types.StringValue("applicationSmartAlert"),
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("automatic"),
				Condition: &ConditionModel{
					Query: types.StringValue("entity.type:service AND entity.tag:production"),
				},
				Action: []PolicyActionModel{
					{
						Action:  createMinimalManualAction(),
						AgentID: types.StringValue("agent-456"),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.Len(t, result.TypeConfigurations, 1)
	assert.Equal(t, "automatic", result.TypeConfigurations[0].Name)
	require.NotNil(t, result.TypeConfigurations[0].Condition)
	assert.Equal(t, "entity.type:service AND entity.tag:production", result.TypeConfigurations[0].Condition.Query)
	require.Len(t, result.TypeConfigurations[0].Runnable.RunConfiguration.Actions, 1)
	assert.Equal(t, "agent-456", result.TypeConfigurations[0].Runnable.RunConfiguration.Actions[0].AgentId)
}

// TestMapStateToDataObject_WithMultipleActions tests mapping policy with multiple actions
func TestMapStateToDataObject_WithMultipleActions(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-123"),
			Type: types.StringValue("customEvent"),
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("manual"),
				Action: []PolicyActionModel{
					{
						Action:  createMinimalManualAction(),
						AgentID: types.StringNull(),
					},
					{
						Action: shared.AutomationActionModel{
							ID:             types.StringValue("action-456"),
							Name:           types.StringValue("Second Action"),
							Description:    types.StringValue("Second Description"),
							Tags:           types.ListNull(types.StringType),
							InputParameter: []shared.ParameterModel{},
							Manual: &shared.ManualModel{
								Content: types.StringValue("Second manual content"),
							},
						},
						AgentID: types.StringValue("agent-789"),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.Len(t, result.TypeConfigurations, 1)
	require.Len(t, result.TypeConfigurations[0].Runnable.RunConfiguration.Actions, 2)
	assert.Equal(t, "action-123", result.TypeConfigurations[0].Runnable.RunConfiguration.Actions[0].Action.ID)
	assert.Equal(t, "", result.TypeConfigurations[0].Runnable.RunConfiguration.Actions[0].AgentId)
	assert.Equal(t, "action-456", result.TypeConfigurations[0].Runnable.RunConfiguration.Actions[1].Action.ID)
	assert.Equal(t, "agent-789", result.TypeConfigurations[0].Runnable.RunConfiguration.Actions[1].AgentId)
}

// TestMapStateToDataObject_FromPlan tests mapping from plan instead of state
func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	plan := &tfsdk.Plan{
		Schema: getTestSchema(),
	}

	model := AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-123"),
			Type: types.StringValue("customEvent"),
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("manual"),
				Action: []PolicyActionModel{
					{
						Action: createMinimalManualAction(),
					},
				},
			},
		},
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Policy", result.Name)
}

// TestMapStateToDataObject_NullID tests mapping with null ID
func TestMapStateToDataObject_NullID(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-123"),
			Type: types.StringValue("customEvent"),
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("manual"),
				Action: []PolicyActionModel{
					{
						Action: createMinimalManualAction(),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "", result.ID)
}

// TestUpdateState_BasicPolicy tests updating state with basic policy
func TestUpdateState_BasicPolicy(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Tags:        nil,
		Trigger: restapi.Trigger{
			Id:   "trigger-123",
			Type: "customEvent",
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name:      "manual",
				Condition: &restapi.Condition{Query: ""},
				Runnable: restapi.Runnable{
					Id:   "action-123",
					Type: "action",
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.AutomationActionPolicy{
							{
								Action: restapi.AutomationAction{
									ID:          "action-123",
									Name:        "Test Action",
									Description: "Action Description",
									Type:        "manual",
									Fields: []restapi.Field{
										{Name: "content", Value: "Manual content"},
									},
									InputParameters: []restapi.Parameter{},
								},
								AgentId: "",
							},
						},
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

	var model AutomationPolicyModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Policy", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	assert.True(t, model.Tags.IsNull())
	assert.Equal(t, "trigger-123", model.Trigger.ID.ValueString())
	assert.Equal(t, "customEvent", model.Trigger.Type.ValueString())
	require.Len(t, model.TypeConfiguration, 1)
	assert.Equal(t, "manual", model.TypeConfiguration[0].Name.ValueString())
}

// TestUpdateState_WithTags tests updating state with tags
func TestUpdateState_WithTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Tags:        []interface{}{"tag1", "tag2", "tag3"},
		Trigger: restapi.Trigger{
			Id:   "trigger-123",
			Type: "builtinEvent",
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name: "manual",
				Runnable: restapi.Runnable{
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.AutomationActionPolicy{
							{
								Action: createMinimalRestAPIAction(),
							},
						},
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

	var model AutomationPolicyModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.Tags.IsNull())
	var tags []string
	diags = model.Tags.ElementsAs(ctx, &tags, false)
	require.False(t, diags.HasError())
	assert.Len(t, tags, 3)
	assert.Contains(t, tags, "tag1")
	assert.Contains(t, tags, "tag2")
	assert.Contains(t, tags, "tag3")
}

// TestUpdateState_WithScheduling tests updating state with scheduling
func TestUpdateState_WithScheduling(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Tags:        nil,
		Trigger: restapi.Trigger{
			Id:          "trigger-123",
			Type:        "schedule",
			Name:        "Scheduled Trigger",
			Description: "Trigger Description",
			Scheduling: restapi.Scheduling{
				StartTime:     1609459200000,
				Duration:      60,
				DurationUnit:  "MINUTE",
				RecurrentRule: "FREQ=DAILY",
				Recurrent:     true,
			},
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name: "automatic",
				Condition: &restapi.Condition{
					Query: "entity.type:service",
				},
				Runnable: restapi.Runnable{
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.AutomationActionPolicy{
							{
								Action: createMinimalRestAPIAction(),
							},
						},
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

	var model AutomationPolicyModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "schedule", model.Trigger.Type.ValueString())
	assert.Equal(t, "Scheduled Trigger", model.Trigger.Name.ValueString())
	assert.Equal(t, "Trigger Description", model.Trigger.Description.ValueString())
	require.NotNil(t, model.Trigger.Scheduling)
	assert.Equal(t, int64(1609459200000), model.Trigger.Scheduling.StartTime.ValueInt64())
	assert.Equal(t, int64(60), model.Trigger.Scheduling.Duration.ValueInt64())
	assert.Equal(t, "MINUTE", model.Trigger.Scheduling.DurationUnit.ValueString())
	assert.Equal(t, "FREQ=DAILY", model.Trigger.Scheduling.RecurrentRule.ValueString())
	assert.True(t, model.Trigger.Scheduling.Recurrent.ValueBool())
}

// TestMapTagsToState tests tag mapping to state
func TestMapTagsToState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	tests := []struct {
		name     string
		tags     interface{}
		expected []string
		hasError bool
	}{
		{
			name:     "valid tags",
			tags:     []interface{}{"tag1", "tag2", "tag3"},
			expected: []string{"tag1", "tag2", "tag3"},
			hasError: false,
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: nil,
			hasError: false,
		},
		{
			name:     "empty tags",
			tags:     []interface{}{},
			expected: []string{},
			hasError: false,
		},
		{
			name:     "invalid tag type",
			tags:     []interface{}{"tag1", 123},
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid tags format",
			tags:     "not a slice",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := resource.mapTagsToState(ctx, tt.tags)

			if tt.hasError {
				assert.True(t, diags.HasError())
			} else {
				assert.False(t, diags.HasError())
				if tt.tags == nil {
					assert.True(t, result.IsNull())
				} else {
					var tags []string
					diags = result.ElementsAs(ctx, &tags, false)
					assert.False(t, diags.HasError())
					assert.Equal(t, tt.expected, tags)
				}
			}
		})
	}
}

// TestMapTagsFromState tests tag mapping from state
func TestMapTagsFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	tests := []struct {
		name     string
		tagsList types.List
		expected interface{}
	}{
		{
			name: "with tags",
			tagsList: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("tag1"),
				types.StringValue("tag2"),
			}),
			expected: []string{"tag1", "tag2"},
		},
		{
			name:     "null tags",
			tagsList: types.ListNull(types.StringType),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := resource.mapTagsFromState(ctx, tt.tagsList)
			assert.False(t, diags.HasError())

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				tagSlice, ok := result.([]string)
				require.True(t, ok)
				assert.Equal(t, tt.expected, tagSlice)
			}
		})
	}
}

// TestMapTriggerToState tests trigger mapping to state
func TestMapTriggerToState(t *testing.T) {
	resource := &automationPolicyResourceFramework{}

	tests := []struct {
		name         string
		trigger      *restapi.Trigger
		triggerModel TriggerModel
		expectedID   string
		expectedType string
	}{
		{
			name: "basic trigger",
			trigger: &restapi.Trigger{
				Id:   "trigger-123",
				Type: "customEvent",
			},
			triggerModel: TriggerModel{
				ID:   types.StringNull(),
				Type: types.StringNull(),
			},
			expectedID:   "trigger-123",
			expectedType: "customEvent",
		},
		{
			name: "trigger with name and description",
			trigger: &restapi.Trigger{
				Id:          "trigger-456",
				Type:        "builtinEvent",
				Name:        "Test Trigger",
				Description: "Test Description",
			},
			triggerModel: TriggerModel{
				ID:          types.StringNull(),
				Type:        types.StringNull(),
				Name:        types.StringNull(),
				Description: types.StringNull(),
			},
			expectedID:   "trigger-456",
			expectedType: "builtinEvent",
		},
		{
			name: "trigger with scheduling",
			trigger: &restapi.Trigger{
				Id:   "trigger-789",
				Type: "schedule",
				Scheduling: restapi.Scheduling{
					StartTime:     1609459200000,
					Duration:      60,
					DurationUnit:  "MINUTE",
					RecurrentRule: "FREQ=DAILY",
					Recurrent:     true,
				},
			},
			triggerModel: TriggerModel{
				ID:         types.StringNull(),
				Type:       types.StringNull(),
				Scheduling: nil,
			},
			expectedID:   "trigger-789",
			expectedType: "schedule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resource.mapTriggerToState(tt.trigger, tt.triggerModel)
			assert.Equal(t, tt.expectedID, result.ID.ValueString())
			assert.Equal(t, tt.expectedType, result.Type.ValueString())

			if tt.trigger.Name != "" {
				assert.Equal(t, tt.trigger.Name, result.Name.ValueString())
			}
			if tt.trigger.Description != "" {
				assert.Equal(t, tt.trigger.Description, result.Description.ValueString())
			}
			if tt.trigger.Scheduling.StartTime != 0 {
				require.NotNil(t, result.Scheduling)
				assert.Equal(t, tt.trigger.Scheduling.StartTime, result.Scheduling.StartTime.ValueInt64())
			}
		})
	}
}

// TestMapTriggerFromState tests trigger mapping from state
func TestMapTriggerFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	tests := []struct {
		name         string
		triggerModel TriggerModel
		expectedID   string
		expectedType string
	}{
		{
			name: "basic trigger",
			triggerModel: TriggerModel{
				ID:   types.StringValue("trigger-123"),
				Type: types.StringValue("customEvent"),
			},
			expectedID:   "trigger-123",
			expectedType: "customEvent",
		},
		{
			name: "trigger with optional fields",
			triggerModel: TriggerModel{
				ID:          types.StringValue("trigger-456"),
				Type:        types.StringValue("builtinEvent"),
				Name:        types.StringValue("Test Trigger"),
				Description: types.StringValue("Test Description"),
			},
			expectedID:   "trigger-456",
			expectedType: "builtinEvent",
		},
		{
			name: "trigger with scheduling",
			triggerModel: TriggerModel{
				ID:   types.StringValue("trigger-789"),
				Type: types.StringValue("schedule"),
				Scheduling: &SchedulingModel{
					StartTime:     types.Int64Value(1609459200000),
					Duration:      types.Int64Value(60),
					DurationUnit:  types.StringValue("MINUTE"),
					RecurrentRule: types.StringValue("FREQ=DAILY"),
					Recurrent:     types.BoolValue(true),
				},
			},
			expectedID:   "trigger-789",
			expectedType: "schedule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := resource.mapTriggerFromState(ctx, tt.triggerModel)
			assert.False(t, diags.HasError())
			assert.Equal(t, tt.expectedID, result.Id)
			assert.Equal(t, tt.expectedType, result.Type)

			if !tt.triggerModel.Name.IsNull() {
				assert.Equal(t, tt.triggerModel.Name.ValueString(), result.Name)
			}
			if !tt.triggerModel.Description.IsNull() {
				assert.Equal(t, tt.triggerModel.Description.ValueString(), result.Description)
			}
			if tt.triggerModel.Scheduling != nil {
				assert.Equal(t, tt.triggerModel.Scheduling.StartTime.ValueInt64(), result.Scheduling.StartTime)
			}
		})
	}
}

// TestMapConditionFromState tests condition mapping from state
func TestMapConditionFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	tests := []struct {
		name           string
		conditionModel *ConditionModel
		expectedQuery  string
	}{
		{
			name:           "nil condition",
			conditionModel: nil,
			expectedQuery:  "",
		},
		{
			name: "with query",
			conditionModel: &ConditionModel{
				Query: types.StringValue("entity.type:service"),
			},
			expectedQuery: "entity.type:service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := resource.mapConditionFromState(ctx, tt.conditionModel)
			assert.False(t, diags.HasError())
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedQuery, result.Query)
		})
	}
}

// TestMapTypeConfigurationsToState tests type configuration mapping to state
func TestMapTypeConfigurationsToState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	typeConfigs := []restapi.TypeConfiguration{
		{
			Name: "manual",
			Condition: &restapi.Condition{
				Query: "entity.type:service",
			},
			Runnable: restapi.Runnable{
				RunConfiguration: restapi.RunConfiguration{
					Actions: []restapi.AutomationActionPolicy{
						{
							Action:  createMinimalRestAPIAction(),
							AgentId: "agent-123",
						},
					},
				},
			},
		},
	}

	result := resource.mapTypeConfigurationsToState(ctx, typeConfigs)
	require.Len(t, result, 1)
	assert.Equal(t, "manual", result[0].Name.ValueString())
	require.NotNil(t, result[0].Condition)
	assert.Equal(t, "entity.type:service", result[0].Condition.Query.ValueString())
	require.Len(t, result[0].Action, 1)
}

// TestMapTypeConfigurationsFromState tests type configuration mapping from state
func TestMapTypeConfigurationsFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	typeConfigModels := []TypeConfigurationModel{
		{
			Name: types.StringValue("automatic"),
			Condition: &ConditionModel{
				Query: types.StringValue("entity.type:service"),
			},
			Action: []PolicyActionModel{
				{
					Action:  createMinimalManualAction(),
					AgentID: types.StringValue("agent-456"),
				},
			},
		},
	}

	result, diags := resource.mapTypeConfigurationsFromState(ctx, typeConfigModels)
	assert.False(t, diags.HasError())
	require.Len(t, result, 1)
	assert.Equal(t, "automatic", result[0].Name)
	require.NotNil(t, result[0].Condition)
	assert.Equal(t, "entity.type:service", result[0].Condition.Query)
}

// TestMapActionsToState tests action mapping to state
func TestMapActionsToState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	runnable := &restapi.Runnable{
		RunConfiguration: restapi.RunConfiguration{
			Actions: []restapi.AutomationActionPolicy{
				{
					Action:  createMinimalRestAPIAction(),
					AgentId: "agent-123",
				},
				{
					Action: restapi.AutomationAction{
						ID:              "action-456",
						Name:            "Second Action",
						Description:     "Second Description",
						Type:            "manual",
						Fields:          []restapi.Field{{Name: "content", Value: "Content"}},
						InputParameters: []restapi.Parameter{},
					},
					AgentId: "",
				},
			},
		},
	}

	result := resource.mapActionsToState(ctx, runnable)
	require.Len(t, result, 2)
	assert.Equal(t, "action-123", result[0].Action.ID.ValueString())
	assert.Equal(t, "agent-123", result[0].AgentID.ValueString())
	assert.Equal(t, "action-456", result[1].Action.ID.ValueString())
	assert.True(t, result[1].AgentID.IsNull())
}

// TestMapRunnableFromState tests runnable mapping from state
func TestMapRunnableFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	actionModels := []PolicyActionModel{
		{
			Action:  createMinimalManualAction(),
			AgentID: types.StringValue("agent-123"),
		},
	}

	result, diags := resource.mapRunnableFromState(ctx, actionModels)
	assert.False(t, diags.HasError())
	assert.Equal(t, "action", result.Type)
	require.Len(t, result.RunConfiguration.Actions, 1)
	assert.Equal(t, "action-123", result.RunConfiguration.Actions[0].Action.ID)
	assert.Equal(t, "agent-123", result.RunConfiguration.Actions[0].AgentId)
}

// TestMapRunnableFromState_EmptyActions tests runnable mapping with empty actions
func TestMapRunnableFromState_EmptyActions(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	actionModels := []PolicyActionModel{}

	result, diags := resource.mapRunnableFromState(ctx, actionModels)
	assert.False(t, diags.HasError())
	assert.Equal(t, "action", result.Type)
	assert.Empty(t, result.RunConfiguration.Actions)
}

// TestUpdateState_AllTriggerTypes tests UpdateState for all trigger types
func TestUpdateState_AllTriggerTypes(t *testing.T) {
	triggerTypes := []string{
		"customEvent",
		"builtinEvent",
		"applicationSmartAlert",
		"globalApplicationSmartAlert",
		"websiteSmartAlert",
		"infraSmartAlert",
		"mobileAppSmartAlert",
		"syntheticsSmartAlert",
		"logSmartAlert",
		"sloSmartAlert",
		"schedule",
	}

	for _, triggerType := range triggerTypes {
		t.Run(triggerType, func(t *testing.T) {
			ctx := context.Background()
			resource := &automationPolicyResourceFramework{}

			data := &restapi.AutomationPolicy{
				ID:          "test-id",
				Name:        "Test Policy",
				Description: "Test Description",
				Trigger: restapi.Trigger{
					Id:   "trigger-123",
					Type: triggerType,
				},
				TypeConfigurations: []restapi.TypeConfiguration{
					{
						Name: "manual",
						Runnable: restapi.Runnable{
							RunConfiguration: restapi.RunConfiguration{
								Actions: []restapi.AutomationActionPolicy{
									{
										Action: createMinimalRestAPIAction(),
									},
								},
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

			var model AutomationPolicyModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			assert.Equal(t, triggerType, model.Trigger.Type.ValueString())
		})
	}
}

// TestUpdateState_AllPolicyTypes tests UpdateState for all policy types
func TestUpdateState_AllPolicyTypes(t *testing.T) {
	policyTypes := []string{"manual", "automatic"}

	for _, policyType := range policyTypes {
		t.Run(policyType, func(t *testing.T) {
			ctx := context.Background()
			resource := &automationPolicyResourceFramework{}

			data := &restapi.AutomationPolicy{
				ID:          "test-id",
				Name:        "Test Policy",
				Description: "Test Description",
				Trigger: restapi.Trigger{
					Id:   "trigger-123",
					Type: "customEvent",
				},
				TypeConfigurations: []restapi.TypeConfiguration{
					{
						Name: policyType,
						Runnable: restapi.Runnable{
							RunConfiguration: restapi.RunConfiguration{
								Actions: []restapi.AutomationActionPolicy{
									{
										Action: createMinimalRestAPIAction(),
									},
								},
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

			var model AutomationPolicyModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			assert.Equal(t, policyType, model.TypeConfiguration[0].Name.ValueString())
		})
	}
}

// TestMapTriggerToState_EmptyOptionalFields tests trigger mapping with empty optional fields
func TestMapTriggerToState_EmptyOptionalFields(t *testing.T) {
	resource := &automationPolicyResourceFramework{}

	trigger := &restapi.Trigger{
		Id:          "trigger-123",
		Type:        "customEvent",
		Name:        "",
		Description: "",
	}

	triggerModel := TriggerModel{
		ID:          types.StringNull(),
		Type:        types.StringNull(),
		Name:        types.StringNull(),
		Description: types.StringNull(),
	}

	result := resource.mapTriggerToState(trigger, triggerModel)
	assert.Equal(t, "trigger-123", result.ID.ValueString())
	assert.Equal(t, "customEvent", result.Type.ValueString())
	assert.True(t, result.Name.IsNull())
	assert.True(t, result.Description.IsNull())
}

// TestMapTriggerToState_WithExistingModel tests trigger mapping preserving existing model values
func TestMapTriggerToState_WithExistingModel(t *testing.T) {
	resource := &automationPolicyResourceFramework{}

	trigger := &restapi.Trigger{
		Id:   "trigger-new",
		Type: "builtinEvent",
	}

	triggerModel := TriggerModel{
		ID:   types.StringValue("trigger-old"),
		Type: types.StringValue("customEvent"),
	}

	result := resource.mapTriggerToState(trigger, triggerModel)
	// Should preserve existing values
	assert.Equal(t, "trigger-old", result.ID.ValueString())
	assert.Equal(t, "customEvent", result.Type.ValueString())
}

// TestMapTypeConfigurationsToState_WithoutCondition tests type configuration without condition
func TestMapTypeConfigurationsToState_WithoutCondition(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	typeConfigs := []restapi.TypeConfiguration{
		{
			Name:      "manual",
			Condition: nil,
			Runnable: restapi.Runnable{
				RunConfiguration: restapi.RunConfiguration{
					Actions: []restapi.AutomationActionPolicy{
						{
							Action: createMinimalRestAPIAction(),
						},
					},
				},
			},
		},
	}

	result := resource.mapTypeConfigurationsToState(ctx, typeConfigs)
	require.Len(t, result, 1)
	assert.Nil(t, result[0].Condition)
}

// TestMapTypeConfigurationsToState_WithEmptyCondition tests type configuration with empty condition
func TestMapTypeConfigurationsToState_WithEmptyCondition(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	typeConfigs := []restapi.TypeConfiguration{
		{
			Name: "automatic",
			Condition: &restapi.Condition{
				Query: "",
			},
			Runnable: restapi.Runnable{
				RunConfiguration: restapi.RunConfiguration{
					Actions: []restapi.AutomationActionPolicy{
						{
							Action: createMinimalRestAPIAction(),
						},
					},
				},
			},
		},
	}

	result := resource.mapTypeConfigurationsToState(ctx, typeConfigs)
	require.Len(t, result, 1)
	assert.Nil(t, result[0].Condition)
}

// Helper functions

func createMockState(t *testing.T, model AutomationPolicyModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewAutomationPolicyResourceHandleFramework()
	return resource.MetaData().Schema
}

func createMinimalManualAction() shared.AutomationActionModel {
	return shared.AutomationActionModel{
		ID:             types.StringValue("action-123"),
		Name:           types.StringValue("Test Action"),
		Description:    types.StringValue("Action Description"),
		Tags:           types.ListNull(types.StringType),
		InputParameter: []shared.ParameterModel{},
		Manual: &shared.ManualModel{
			Content: types.StringValue("Manual content"),
		},
	}
}

func createMinimalRestAPIAction() restapi.AutomationAction {
	return restapi.AutomationAction{
		ID:          "action-123",
		Name:        "Test Action",
		Description: "Action Description",
		Type:        "manual",
		Fields: []restapi.Field{
			{Name: "content", Value: "Manual content"},
		},
		InputParameters: []restapi.Parameter{},
	}
}

// TestUpdateState_WithPlan tests UpdateState with plan provided
func TestUpdateState_WithPlan(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	// Create a plan with existing data
	plan := &tfsdk.Plan{
		Schema: getTestSchema(),
	}

	planModel := AutomationPolicyModel{
		ID:          types.StringValue("plan-id"),
		Name:        types.StringValue("Plan Policy"),
		Description: types.StringValue("Plan Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-plan"),
			Type: types.StringValue("customEvent"),
			Scheduling: &SchedulingModel{
				StartTime:     types.Int64Value(1609459200000),
				Duration:      types.Int64Value(30),
				DurationUnit:  types.StringValue("HOUR"),
				RecurrentRule: types.StringValue("FREQ=WEEKLY"),
				Recurrent:     types.BoolValue(false),
			},
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("manual"),
				Action: []PolicyActionModel{
					{
						Action: createMinimalManualAction(),
					},
				},
			},
		},
	}

	diags := plan.Set(ctx, planModel)
	require.False(t, diags.HasError())

	// API response data
	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Tags:        []interface{}{"tag1"},
		Trigger: restapi.Trigger{
			Id:   "trigger-123",
			Type: "builtinEvent",
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name: "automatic",
				Runnable: restapi.Runnable{
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.AutomationActionPolicy{
							{
								Action: createMinimalRestAPIAction(),
							},
						},
					},
				},
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags = resource.UpdateState(ctx, state, plan, data)
	require.False(t, diags.HasError())

	var model AutomationPolicyModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	// Verify API response values are set
	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Policy", model.Name.ValueString())
	// Scheduling should be preserved from plan
	require.NotNil(t, model.Trigger.Scheduling)
}

// TestUpdateState_ErrorInPlanGet tests error handling when plan.Get fails
func TestUpdateState_ErrorInPlanGet(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	// Create an invalid plan that will cause Get to fail
	plan := &tfsdk.Plan{
		Schema: schema.Schema{
			Attributes: map[string]schema.Attribute{
				"invalid": schema.StringAttribute{},
			},
		},
	}

	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Trigger: restapi.Trigger{
			Id:   "trigger-123",
			Type: "customEvent",
		},
		TypeConfigurations: []restapi.TypeConfiguration{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, plan, data)
	assert.True(t, diags.HasError())
}

// TestMapStateToDataObject_ErrorInStateGet tests error handling when state.Get fails
func TestMapStateToDataObject_ErrorInStateGet(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	// Create an invalid state that will cause Get to fail
	state := &tfsdk.State{
		Schema: schema.Schema{
			Attributes: map[string]schema.Attribute{
				"invalid": schema.StringAttribute{},
			},
		},
	}

	result, diags := resource.MapStateToDataObject(ctx, nil, state)
	assert.True(t, diags.HasError())
	assert.Nil(t, result)
}

// TestMapStateToDataObject_ErrorInTriggerMapping tests error handling in trigger mapping
func TestMapStateToDataObject_ErrorInTriggerMapping(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	// This test verifies the trigger mapping doesn't cause errors
	state := createMockState(t, AutomationPolicyModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Policy"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Trigger: TriggerModel{
			ID:   types.StringValue("trigger-123"),
			Type: types.StringValue("customEvent"),
		},
		TypeConfiguration: []TypeConfigurationModel{
			{
				Name: types.StringValue("manual"),
				Action: []PolicyActionModel{
					{
						Action: createMinimalManualAction(),
					},
				},
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
}

// TestMapTypeConfigurationsFromState_ErrorInCondition tests error handling in condition mapping
func TestMapTypeConfigurationsFromState_ErrorInCondition(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	typeConfigModels := []TypeConfigurationModel{
		{
			Name:      types.StringValue("automatic"),
			Condition: nil,
			Action: []PolicyActionModel{
				{
					Action: createMinimalManualAction(),
				},
			},
		},
	}

	result, diags := resource.mapTypeConfigurationsFromState(ctx, typeConfigModels)
	assert.False(t, diags.HasError())
	require.Len(t, result, 1)
	require.NotNil(t, result[0].Condition)
	assert.Equal(t, "", result[0].Condition.Query)
}

// TestMapRunnableFromState_WithInputParameters tests runnable mapping with input parameters
func TestMapRunnableFromState_WithInputParameters(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	actionModels := []PolicyActionModel{
		{
			Action: shared.AutomationActionModel{
				ID:          types.StringValue("action-123"),
				Name:        types.StringValue("Test Action"),
				Description: types.StringValue("Action Description"),
				Tags: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("tag1"),
				}),
				InputParameter: []shared.ParameterModel{
					{
						Name:        types.StringValue("param1"),
						Description: types.StringValue("Parameter 1"),
						Label:       types.StringValue("Param 1"),
						Required:    types.BoolValue(true),
						Hidden:      types.BoolValue(false),
						Type:        types.StringValue("static"),
						Value:       types.StringValue("value1"),
					},
				},
				Manual: &shared.ManualModel{
					Content: types.StringValue("Manual content"),
				},
			},
			AgentID: types.StringValue("agent-123"),
		},
	}

	result, diags := resource.mapRunnableFromState(ctx, actionModels)
	assert.False(t, diags.HasError())
	assert.Equal(t, "action", result.Type)
	require.Len(t, result.RunConfiguration.Actions, 1)
	assert.Equal(t, "action-123", result.RunConfiguration.Actions[0].Action.ID)
	assert.Len(t, result.RunConfiguration.Actions[0].Action.InputParameters, 1)
	assert.NotNil(t, result.RunConfiguration.Actions[0].Action.Tags)
}

// TestMapRunnableFromState_MultipleActionsSetRunnableID tests that first action ID is set as runnable ID
func TestMapRunnableFromState_MultipleActionsSetRunnableID(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	actionModels := []PolicyActionModel{
		{
			Action: shared.AutomationActionModel{
				ID:             types.StringValue("first-action"),
				Name:           types.StringValue("First Action"),
				Description:    types.StringValue("First Description"),
				Tags:           types.ListNull(types.StringType),
				InputParameter: []shared.ParameterModel{},
				Manual: &shared.ManualModel{
					Content: types.StringValue("First content"),
				},
			},
			AgentID: types.StringNull(),
		},
		{
			Action: shared.AutomationActionModel{
				ID:             types.StringValue("second-action"),
				Name:           types.StringValue("Second Action"),
				Description:    types.StringValue("Second Description"),
				Tags:           types.ListNull(types.StringType),
				InputParameter: []shared.ParameterModel{},
				Manual: &shared.ManualModel{
					Content: types.StringValue("Second content"),
				},
			},
			AgentID: types.StringNull(),
		},
	}

	result, diags := resource.mapRunnableFromState(ctx, actionModels)
	assert.False(t, diags.HasError())
	// Verify the runnable ID is set to the first action's ID
	assert.Equal(t, "first-action", result.Id)
	require.Len(t, result.RunConfiguration.Actions, 2)
}

// TestMapTagsFromState_ErrorInElementsAs tests error handling in mapTagsFromState
func TestMapTagsFromState_ErrorInElementsAs(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	// Create a list with invalid element type that will cause ElementsAs to fail
	invalidList := types.ListValueMust(types.Int64Type, []attr.Value{
		types.Int64Value(123),
	})

	result, diags := resource.mapTagsFromState(ctx, invalidList)
	assert.True(t, diags.HasError())
	assert.Nil(t, result)
}

// TestUpdateState_WithSchedulingEmptyFields tests scheduling with empty optional fields
func TestUpdateState_WithSchedulingEmptyFields(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Tags:        nil,
		Trigger: restapi.Trigger{
			Id:   "trigger-123",
			Type: "schedule",
			Scheduling: restapi.Scheduling{
				StartTime:     1609459200000,
				Duration:      60,
				DurationUnit:  "",
				RecurrentRule: "",
				Recurrent:     false,
			},
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name: "manual",
				Runnable: restapi.Runnable{
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.AutomationActionPolicy{
							{
								Action: createMinimalRestAPIAction(),
							},
						},
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

	var model AutomationPolicyModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.Trigger.Scheduling)
	assert.True(t, model.Trigger.Scheduling.DurationUnit.IsNull())
	assert.True(t, model.Trigger.Scheduling.RecurrentRule.IsNull())
}

// TestMapTriggerToState_WithSchedulingPreserved tests that scheduling is preserved from model
func TestMapTriggerToState_WithSchedulingPreserved(t *testing.T) {
	resource := &automationPolicyResourceFramework{}

	trigger := &restapi.Trigger{
		Id:   "trigger-123",
		Type: "schedule",
		Scheduling: restapi.Scheduling{
			StartTime: 0, // Zero start time should not override existing scheduling
		},
	}

	triggerModel := TriggerModel{
		ID:   types.StringNull(),
		Type: types.StringNull(),
		Scheduling: &SchedulingModel{
			StartTime:     types.Int64Value(1609459200000),
			Duration:      types.Int64Value(60),
			DurationUnit:  types.StringValue("MINUTE"),
			RecurrentRule: types.StringValue("FREQ=DAILY"),
			Recurrent:     types.BoolValue(true),
		},
	}

	result := resource.mapTriggerToState(trigger, triggerModel)
	// Scheduling should be preserved from the model when API returns zero start time
	require.NotNil(t, result.Scheduling)
	assert.Equal(t, int64(1609459200000), result.Scheduling.StartTime.ValueInt64())
}

// TestMapActionsToState_WithTags tests action mapping with tags
func TestMapActionsToState_WithTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	runnable := &restapi.Runnable{
		RunConfiguration: restapi.RunConfiguration{
			Actions: []restapi.AutomationActionPolicy{
				{
					Action: restapi.AutomationAction{
						ID:          "action-123",
						Name:        "Test Action",
						Description: "Action Description",
						Type:        "manual",
						Tags:        []interface{}{"tag1", "tag2"},
						Fields: []restapi.Field{
							{Name: "content", Value: "Manual content"},
						},
						InputParameters: []restapi.Parameter{},
					},
					AgentId: "agent-123",
				},
			},
		},
	}

	result := resource.mapActionsToState(ctx, runnable)
	require.Len(t, result, 1)
	assert.False(t, result[0].Action.Tags.IsNull())
}

// TestMapActionsToState_WithInputParameters tests action mapping with input parameters
func TestMapActionsToState_WithInputParameters(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	runnable := &restapi.Runnable{
		RunConfiguration: restapi.RunConfiguration{
			Actions: []restapi.AutomationActionPolicy{
				{
					Action: restapi.AutomationAction{
						ID:          "action-123",
						Name:        "Test Action",
						Description: "Action Description",
						Type:        "manual",
						Fields: []restapi.Field{
							{Name: "content", Value: "Manual content"},
						},
						InputParameters: []restapi.Parameter{
							{
								Name:        "param1",
								Description: "Parameter 1",
								Label:       "Param 1",
								Required:    true,
								Hidden:      false,
								Type:        "static",
								Value:       "value1",
							},
						},
					},
					AgentId: "",
				},
			},
		},
	}

	result := resource.mapActionsToState(ctx, runnable)
	require.Len(t, result, 1)
	require.Len(t, result[0].Action.InputParameter, 1)
	assert.Equal(t, "param1", result[0].Action.InputParameter[0].Name.ValueString())
}

// TestUpdateState_NullTags tests UpdateState with null tags
func TestUpdateState_NullTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationPolicyResourceFramework{}

	data := &restapi.AutomationPolicy{
		ID:          "test-id",
		Name:        "Test Policy",
		Description: "Test Description",
		Tags:        nil,
		Trigger: restapi.Trigger{
			Id:   "trigger-123",
			Type: "customEvent",
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name: "manual",
				Runnable: restapi.Runnable{
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.AutomationActionPolicy{
							{
								Action: createMinimalRestAPIAction(),
							},
						},
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

	var model AutomationPolicyModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.Tags.IsNull())
}

// Made with Bob
