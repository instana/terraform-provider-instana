package syntheticalertconfig

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSyntheticAlertConfigResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewSyntheticAlertConfigResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSyntheticAlertConfig, metadata.ResourceName)
		assert.Equal(t, int64(1), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewSyntheticAlertConfigResourceHandle()
		metadata := handle.MetaData()

		schema := metadata.Schema
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldID])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldName])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldDescription])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldSyntheticTestIds])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldSeverity])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldTagFilter])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldRule])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldAlertChannelIds])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldTimeThreshold])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldGracePeriod])
		assert.NotNil(t, schema.Attributes[SyntheticAlertConfigFieldCustomPayloadField])
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &syntheticAlertConfigResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaSyntheticAlertConfig,
				SchemaVersion: 1,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSyntheticAlertConfig, metadata.ResourceName)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return synthetic alert config rest resource", func(t *testing.T) {
		resource := &syntheticAlertConfigResource{}

		mockAPI := &mockInstanaAPI{}
		restResource := resource.GetRestResource(mockAPI)

		assert.NotNil(t, restResource)
	})
}

// Mock API for testing
type mockInstanaAPI struct{}

func (m *mockInstanaAPI) CustomEventSpecifications() restapi.RestResource[*restapi.CustomEventSpecification] {
	return nil
}
func (m *mockInstanaAPI) BuiltinEventSpecifications() restapi.ReadOnlyRestResource[*restapi.BuiltinEventSpecification] {
	return nil
}
func (m *mockInstanaAPI) APITokens() restapi.RestResource[*restapi.APIToken] { return nil }
func (m *mockInstanaAPI) ApplicationConfigs() restapi.RestResource[*restapi.ApplicationConfig] {
	return nil
}
func (m *mockInstanaAPI) ApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) GlobalApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) AlertingChannels() restapi.RestResource[*restapi.AlertingChannel] {
	return nil
}
func (m *mockInstanaAPI) AlertingConfigurations() restapi.RestResource[*restapi.AlertingConfiguration] {
	return nil
}
func (m *mockInstanaAPI) SliConfigs() restapi.RestResource[*restapi.SliConfig]          { return nil }
func (m *mockInstanaAPI) SloConfigs() restapi.RestResource[*restapi.SloConfig]          { return nil }
func (m *mockInstanaAPI) SloAlertConfig() restapi.RestResource[*restapi.SloAlertConfig] { return nil }
func (m *mockInstanaAPI) SloCorrectionConfig() restapi.RestResource[*restapi.SloCorrectionConfig] {
	return nil
}
func (m *mockInstanaAPI) WebsiteMonitoringConfig() restapi.RestResource[*restapi.WebsiteMonitoringConfig] {
	return nil
}
func (m *mockInstanaAPI) WebsiteAlertConfig() restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) InfraAlertConfig() restapi.RestResource[*restapi.InfraAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) Teams() restapi.RestResource[*restapi.Team]   { return nil }
func (m *mockInstanaAPI) Groups() restapi.RestResource[*restapi.Group] { return nil }
func (m *mockInstanaAPI) CustomDashboards() restapi.RestResource[*restapi.CustomDashboard] {
	return nil
}
func (m *mockInstanaAPI) SyntheticTest() restapi.RestResource[*restapi.SyntheticTest] { return nil }
func (m *mockInstanaAPI) SyntheticLocation() restapi.ReadOnlyRestResource[*restapi.SyntheticLocation] {
	return nil
}
func (m *mockInstanaAPI) SyntheticAlertConfigs() restapi.RestResource[*restapi.SyntheticAlertConfig] {
	return &mockSyntheticAlertConfigRestResource{}
}
func (m *mockInstanaAPI) AutomationActions() restapi.RestResource[*restapi.AutomationAction] {
	return nil
}
func (m *mockInstanaAPI) AutomationPolicies() restapi.RestResource[*restapi.AutomationPolicy] {
	return nil
}
func (m *mockInstanaAPI) HostAgents() restapi.ReadOnlyRestResource[*restapi.HostAgent]  { return nil }
func (m *mockInstanaAPI) LogAlertConfig() restapi.RestResource[*restapi.LogAlertConfig] { return nil }

// Mock rest resource
type mockSyntheticAlertConfigRestResource struct{}

func (m *mockSyntheticAlertConfigRestResource) GetAll() (*[]*restapi.SyntheticAlertConfig, error) {
	return nil, nil
}

func (m *mockSyntheticAlertConfigRestResource) GetOne(id string) (*restapi.SyntheticAlertConfig, error) {
	return nil, nil
}

func (m *mockSyntheticAlertConfigRestResource) Create(data *restapi.SyntheticAlertConfig) (*restapi.SyntheticAlertConfig, error) {
	return nil, nil
}

func (m *mockSyntheticAlertConfigRestResource) Update(data *restapi.SyntheticAlertConfig) (*restapi.SyntheticAlertConfig, error) {
	return nil, nil
}

func (m *mockSyntheticAlertConfigRestResource) Delete(data *restapi.SyntheticAlertConfig) error {
	return nil
}

func (m *mockSyntheticAlertConfigRestResource) DeleteByID(id string) error {
	return nil
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		resource := &syntheticAlertConfigResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaSyntheticAlertConfig,
				Schema:        NewSyntheticAlertConfigResourceHandle().MetaData().Schema,
				SchemaVersion: 1,
			},
		}
		ctx := context.Background()

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		diags := resource.SetComputedFields(ctx, plan)
		assert.False(t, diags.HasError())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &syntheticAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticAlertConfig,
			Schema:        NewSyntheticAlertConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map complete model from state successfully", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id-123"),
			Name:        types.StringValue("Test Alert"),
			Description: types.StringValue("Test Description"),
			Severity:    types.Int64Value(5),
			TagFilter:   types.StringNull(),
			GracePeriod: types.Int64Value(60000),
			Rule: &SyntheticAlertRuleModel{
				AlertType:   types.StringValue("failure"),
				MetricName:  types.StringValue("status"),
				Aggregation: types.StringValue("SUM"),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(3),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1", "test-2"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "test-id-123", result.ID)
		assert.Equal(t, "Test Alert", result.Name)
		assert.Equal(t, "Test Description", result.Description)
		assert.Equal(t, 5, result.Severity)
		assert.Equal(t, int64(60000), result.GracePeriod)
		assert.Equal(t, "failure", result.Rule.AlertType)
		assert.Equal(t, "status", result.Rule.MetricName)
		assert.Equal(t, "SUM", result.Rule.Aggregation)
		assert.Equal(t, "violationsInSequence", result.TimeThreshold.Type)
		assert.Equal(t, 3, result.TimeThreshold.ViolationsCount)
		assert.ElementsMatch(t, []string{"test-1", "test-2"}, result.SyntheticTestIds)
		assert.ElementsMatch(t, []string{"channel-1"}, result.AlertChannelIds)
	})

	t.Run("should map model from plan successfully", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("plan-id"),
			Name:        types.StringValue("Plan Alert"),
			Description: types.StringValue("Plan Description"),
			Severity:    types.Int64Value(10),
			TagFilter:   types.StringNull(),
			GracePeriod: types.Int64Null(),
			Rule: &SyntheticAlertRuleModel{
				AlertType:   types.StringValue("failure"),
				MetricName:  types.StringValue("status"),
				Aggregation: types.StringNull(),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(1),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-3"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-2", "channel-3"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, plan, nil)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "plan-id", result.ID)
		assert.Equal(t, "Plan Alert", result.Name)
		assert.Equal(t, 10, result.Severity)
		assert.Equal(t, int64(0), result.GracePeriod)
		assert.Equal(t, "", result.Rule.Aggregation)
		assert.ElementsMatch(t, []string{"test-3"}, result.SyntheticTestIds)
		assert.ElementsMatch(t, []string{"channel-2", "channel-3"}, result.AlertChannelIds)
	})

	t.Run("should handle when both plan and state are nil", func(t *testing.T) {
		result, diags := resource.MapStateToDataObject(ctx, nil, nil)

		assert.NotNil(t, result)
		assert.False(t, diags.HasError())
	})

	t.Run("should handle null tag filter", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Severity:    types.Int64Value(5),
			TagFilter:   types.StringNull(),
			Rule: &SyntheticAlertRuleModel{
				AlertType:  types.StringValue("failure"),
				MetricName: types.StringValue("status"),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(2),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.NotNil(t, result.TagFilterExpression)
		assert.Equal(t, restapi.TagFilterExpressionElementType(TagFilterTypeExpression), result.TagFilterExpression.Type)
	})

	t.Run("should parse tag filter expression", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Severity:    types.Int64Value(5),
			TagFilter:   types.StringValue("entity.type EQUALS 'synthetic_test'"),
			Rule: &SyntheticAlertRuleModel{
				AlertType:  types.StringValue("failure"),
				MetricName: types.StringValue("status"),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(2),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.NotNil(t, result.TagFilterExpression)
	})

	t.Run("should handle invalid tag filter expression", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Severity:    types.Int64Value(5),
			TagFilter:   types.StringValue("invalid tag filter syntax"),
			Rule: &SyntheticAlertRuleModel{
				AlertType:  types.StringValue("failure"),
				MetricName: types.StringValue("status"),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(2),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.True(t, resultDiags.HasError())
		assert.Nil(t, result)
	})

	t.Run("should handle null severity", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Severity:    types.Int64Null(),
			TagFilter:   types.StringNull(),
			Rule: &SyntheticAlertRuleModel{
				AlertType:  types.StringValue("failure"),
				MetricName: types.StringValue("status"),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(2),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Severity)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &syntheticAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticAlertConfig,
			Schema:        NewSyntheticAlertConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with complete API object", func(t *testing.T) {
		operator := restapi.LogicalOperatorType("AND")
		apiObject := &restapi.SyntheticAlertConfig{
			ID:               "api-id-123",
			Name:             "API Alert",
			Description:      "API Description",
			SyntheticTestIds: []string{"test-a", "test-b"},
			Severity:         5,
			TagFilterExpression: &restapi.TagFilter{
				Type:            TagFilterTypeExpression,
				LogicalOperator: &operator,
				Elements:        []*restapi.TagFilter{},
			},
			Rule: restapi.SyntheticAlertRule{
				AlertType:   "failure",
				MetricName:  "status",
				Aggregation: "MEAN",
			},
			AlertChannelIds: []string{"channel-a", "channel-b"},
			TimeThreshold: restapi.SyntheticAlertTimeThreshold{
				Type:            "violationsInSequence",
				ViolationsCount: 5,
			},
			GracePeriod:           120000,
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SyntheticAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-id-123", model.ID.ValueString())
		assert.Equal(t, "API Alert", model.Name.ValueString())
		assert.Equal(t, "API Description", model.Description.ValueString())
		assert.Equal(t, int64(5), model.Severity.ValueInt64())
		assert.Equal(t, int64(120000), model.GracePeriod.ValueInt64())
		assert.Equal(t, "failure", model.Rule.AlertType.ValueString())
		assert.Equal(t, "status", model.Rule.MetricName.ValueString())
		assert.Equal(t, "MEAN", model.Rule.Aggregation.ValueString())
		assert.Equal(t, "violationsInSequence", model.TimeThreshold.Type.ValueString())
		assert.Equal(t, int64(5), model.TimeThreshold.ViolationsCount.ValueInt64())

		var syntheticTestIds []string
		diags = model.SyntheticTestIds.ElementsAs(ctx, &syntheticTestIds, false)
		assert.False(t, diags.HasError())
		assert.ElementsMatch(t, []string{"test-a", "test-b"}, syntheticTestIds)

		var alertChannelIds []string
		diags = model.AlertChannelIds.ElementsAs(ctx, &alertChannelIds, false)
		assert.False(t, diags.HasError())
		assert.ElementsMatch(t, []string{"channel-a", "channel-b"}, alertChannelIds)
	})

	t.Run("should update state with null grace period", func(t *testing.T) {
		operator := restapi.LogicalOperatorType("AND")
		apiObject := &restapi.SyntheticAlertConfig{
			ID:               "api-id-456",
			Name:             "No Grace Period",
			Description:      "No Grace Period Description",
			SyntheticTestIds: []string{"test-c"},
			Severity:         10,
			TagFilterExpression: &restapi.TagFilter{
				Type:            TagFilterTypeExpression,
				LogicalOperator: &operator,
				Elements:        []*restapi.TagFilter{},
			},
			Rule: restapi.SyntheticAlertRule{
				AlertType:  "failure",
				MetricName: "status",
			},
			AlertChannelIds: []string{"channel-c"},
			TimeThreshold: restapi.SyntheticAlertTimeThreshold{
				Type:            "violationsInSequence",
				ViolationsCount: 1,
			},
			GracePeriod:           0,
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SyntheticAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-id-456", model.ID.ValueString())
		assert.True(t, model.GracePeriod.IsNull())
	})

	t.Run("should update state with null aggregation", func(t *testing.T) {
		operator := restapi.LogicalOperatorType("AND")
		apiObject := &restapi.SyntheticAlertConfig{
			ID:               "api-id-789",
			Name:             "No Aggregation",
			Description:      "No Aggregation Description",
			SyntheticTestIds: []string{"test-d"},
			Severity:         5,
			TagFilterExpression: &restapi.TagFilter{
				Type:            TagFilterTypeExpression,
				LogicalOperator: &operator,
				Elements:        []*restapi.TagFilter{},
			},
			Rule: restapi.SyntheticAlertRule{
				AlertType:   "failure",
				MetricName:  "status",
				Aggregation: "",
			},
			AlertChannelIds: []string{"channel-d"},
			TimeThreshold: restapi.SyntheticAlertTimeThreshold{
				Type:            "violationsInSequence",
				ViolationsCount: 2,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SyntheticAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.True(t, model.Rule.Aggregation.IsNull())
	})

	t.Run("should update state with different severity values", func(t *testing.T) {
		severities := []int{5, 10}

		for _, severity := range severities {
			t.Run("severity_"+string(rune(severity)), func(t *testing.T) {
				operator := restapi.LogicalOperatorType("AND")
				apiObject := &restapi.SyntheticAlertConfig{
					ID:               "test-id",
					Name:             "Test",
					Description:      "Desc",
					SyntheticTestIds: []string{"test-1"},
					Severity:         severity,
					TagFilterExpression: &restapi.TagFilter{
						Type:            TagFilterTypeExpression,
						LogicalOperator: &operator,
						Elements:        []*restapi.TagFilter{},
					},
					Rule: restapi.SyntheticAlertRule{
						AlertType:  "failure",
						MetricName: "status",
					},
					AlertChannelIds: []string{"channel-1"},
					TimeThreshold: restapi.SyntheticAlertTimeThreshold{
						Type:            "violationsInSequence",
						ViolationsCount: 1,
					},
					CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
				}

				state := &tfsdk.State{
					Schema: resource.metaData.Schema,
				}

				// Initialize state with empty model
				initializeEmptyState(t, ctx, state)

				diags := resource.UpdateState(ctx, state, nil, apiObject)

				assert.False(t, diags.HasError())

				var model SyntheticAlertConfigModel
				diags = state.Get(ctx, &model)
				assert.False(t, diags.HasError())

				assert.Equal(t, int64(severity), model.Severity.ValueInt64())
			})
		}
	})

	t.Run("should handle multiple synthetic test IDs", func(t *testing.T) {
		operator := restapi.LogicalOperatorType("AND")
		apiObject := &restapi.SyntheticAlertConfig{
			ID:               "test-id",
			Name:             "Multiple Tests",
			Description:      "Desc",
			SyntheticTestIds: []string{"test-1", "test-2", "test-3", "test-4", "test-5"},
			Severity:         5,
			TagFilterExpression: &restapi.TagFilter{
				Type:            TagFilterTypeExpression,
				LogicalOperator: &operator,
				Elements:        []*restapi.TagFilter{},
			},
			Rule: restapi.SyntheticAlertRule{
				AlertType:  "failure",
				MetricName: "status",
			},
			AlertChannelIds: []string{"channel-1"},
			TimeThreshold: restapi.SyntheticAlertTimeThreshold{
				Type:            "violationsInSequence",
				ViolationsCount: 1,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SyntheticAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		var syntheticTestIds []string
		diags = model.SyntheticTestIds.ElementsAs(ctx, &syntheticTestIds, false)
		assert.False(t, diags.HasError())
		assert.Len(t, syntheticTestIds, 5)
		assert.ElementsMatch(t, []string{"test-1", "test-2", "test-3", "test-4", "test-5"}, syntheticTestIds)
	})

	t.Run("should handle multiple alert channel IDs", func(t *testing.T) {
		operator := restapi.LogicalOperatorType("AND")
		apiObject := &restapi.SyntheticAlertConfig{
			ID:               "test-id",
			Name:             "Multiple Channels",
			Description:      "Desc",
			SyntheticTestIds: []string{"test-1"},
			Severity:         5,
			TagFilterExpression: &restapi.TagFilter{
				Type:            TagFilterTypeExpression,
				LogicalOperator: &operator,
				Elements:        []*restapi.TagFilter{},
			},
			Rule: restapi.SyntheticAlertRule{
				AlertType:  "failure",
				MetricName: "status",
			},
			AlertChannelIds: []string{"channel-1", "channel-2", "channel-3"},
			TimeThreshold: restapi.SyntheticAlertTimeThreshold{
				Type:            "violationsInSequence",
				ViolationsCount: 1,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SyntheticAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		var alertChannelIds []string
		diags = model.AlertChannelIds.ElementsAs(ctx, &alertChannelIds, false)
		assert.False(t, diags.HasError())
		assert.Len(t, alertChannelIds, 3)
		assert.ElementsMatch(t, []string{"channel-1", "channel-2", "channel-3"}, alertChannelIds)
	})
}

// Made with Bob
func TestUpdateStateWithNullTagFilter(t *testing.T) {
	resource := &syntheticAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticAlertConfig,
			Schema:        NewSyntheticAlertConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with null tag filter", func(t *testing.T) {
		apiObject := &restapi.SyntheticAlertConfig{
			ID:                  "test-id",
			Name:                "No Tag Filter",
			Description:         "Desc",
			SyntheticTestIds:    []string{"test-1"},
			Severity:            5,
			TagFilterExpression: nil,
			Rule: restapi.SyntheticAlertRule{
				AlertType:  "failure",
				MetricName: "status",
			},
			AlertChannelIds: []string{"channel-1"},
			TimeThreshold: restapi.SyntheticAlertTimeThreshold{
				Type:            "violationsInSequence",
				ViolationsCount: 1,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model SyntheticAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.True(t, model.TagFilter.IsNull())
	})
}

// Made with Bob
func TestAdditionalCoverage(t *testing.T) {
	resource := &syntheticAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticAlertConfig,
			Schema:        NewSyntheticAlertConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should handle rule with null aggregation in MapStateToDataObject", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Severity:    types.Int64Value(5),
			TagFilter:   types.StringNull(),
			GracePeriod: types.Int64Null(),
			Rule: &SyntheticAlertRuleModel{
				AlertType:   types.StringValue("failure"),
				MetricName:  types.StringValue("status"),
				Aggregation: types.StringNull(),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(1),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "", result.Rule.Aggregation)
		assert.Equal(t, int64(0), result.GracePeriod)
	})

	t.Run("should handle different violations count in time threshold", func(t *testing.T) {
		model := SyntheticAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Severity:    types.Int64Value(5),
			TagFilter:   types.StringNull(),
			Rule: &SyntheticAlertRuleModel{
				AlertType:  types.StringValue("failure"),
				MetricName: types.StringValue("status"),
			},
			TimeThreshold: &SyntheticAlertTimeThresholdModel{
				Type:            types.StringValue("violationsInSequence"),
				ViolationsCount: types.Int64Value(12),
			},
		}

		syntheticTestIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"test-1"})
		model.SyntheticTestIds = syntheticTestIds

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIds = alertChannelIds

		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, 12, result.TimeThreshold.ViolationsCount)
	})
}

// initializeEmptyState initializes the state with an empty model to ensure proper state initialization
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := SyntheticAlertConfigModel{
		ID:                  types.StringNull(),
		Name:                types.StringNull(),
		Description:         types.StringNull(),
		Severity:            types.Int64Null(),
		TagFilter:           types.StringNull(),
		GracePeriod:         types.Int64Null(),
		SyntheticTestIds:    types.SetNull(types.StringType),
		AlertChannelIds:     types.SetNull(types.StringType),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		Rule:                nil,
		TimeThreshold:       nil,
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}

// Made with Bob
