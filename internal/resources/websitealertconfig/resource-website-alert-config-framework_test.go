package websitealertconfig

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebsiteAlertConfigResourceHandleFramework(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewWebsiteAlertConfigResourceHandleFramework()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaWebsiteAlertConfigFramework, metadata.ResourceName)
		assert.Equal(t, int64(1), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewWebsiteAlertConfigResourceHandleFramework()
		metadata := handle.MetaData()

		schema := metadata.Schema
		assert.NotNil(t, schema.Attributes["id"])
		assert.NotNil(t, schema.Attributes["name"])
		assert.NotNil(t, schema.Attributes["description"])
		assert.NotNil(t, schema.Attributes["triggering"])
		assert.NotNil(t, schema.Attributes["website_id"])
		assert.NotNil(t, schema.Attributes["tag_filter"])
		assert.NotNil(t, schema.Attributes["alert_channel_ids"])
		assert.NotNil(t, schema.Attributes["granularity"])
		assert.NotNil(t, schema.Attributes["rules"])
		assert.NotNil(t, schema.Attributes["custom_payload_fields"])
		assert.NotNil(t, schema.Attributes["time_threshold"])
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &websiteAlertConfigResourceFramework{
			metaData: resourcehandle.ResourceMetaDataFramework{
				ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
				SchemaVersion: 1,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaWebsiteAlertConfigFramework, metadata.ResourceName)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return website alert config rest resource", func(t *testing.T) {
		resource := &websiteAlertConfigResourceFramework{}

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
	return &mockWebsiteAlertConfigRestResource{}
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
	return nil
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
type mockWebsiteAlertConfigRestResource struct{}

func (m *mockWebsiteAlertConfigRestResource) GetAll() (*[]*restapi.WebsiteAlertConfig, error) {
	return nil, nil
}

func (m *mockWebsiteAlertConfigRestResource) GetOne(id string) (*restapi.WebsiteAlertConfig, error) {
	return nil, nil
}

func (m *mockWebsiteAlertConfigRestResource) Create(data *restapi.WebsiteAlertConfig) (*restapi.WebsiteAlertConfig, error) {
	return nil, nil
}

func (m *mockWebsiteAlertConfigRestResource) Update(data *restapi.WebsiteAlertConfig) (*restapi.WebsiteAlertConfig, error) {
	return nil, nil
}

func (m *mockWebsiteAlertConfigRestResource) Delete(data *restapi.WebsiteAlertConfig) error {
	return nil
}

func (m *mockWebsiteAlertConfigRestResource) DeleteByID(id string) error {
	return nil
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		resource := &websiteAlertConfigResourceFramework{
			metaData: resourcehandle.ResourceMetaDataFramework{
				ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
				Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
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

func TestMapStateToDataObject_WithRules(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map multiple rules with thresholds", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Multi Rule Alert"),
			Description: types.StringValue("Multi Rule Description"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules: []RuleWithThresholdPluginModel{
				{
					Rule: &WebsiteAlertRuleModel{
						Slowness: &WebsiteAlertRuleConfigModel{
							MetricName:  types.StringValue("latency"),
							Aggregation: types.StringValue("MEAN"),
						},
					},
					ThresholdOperator: types.StringValue(">="),
					Thresholds: &shared.ThresholdAllPluginModel{
						Warning: &shared.ThresholdAllTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float32Value(float32(1000)),
							},
						},
						Critical: &shared.ThresholdAllTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float32Value(float32(2000)),
							},
						},
					},
				},
			},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInSequence: &WebsiteViolationsInSequenceModel{
					TimeWindow: types.Int64Value(300000),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Rules, 1)
		assert.Equal(t, "slowness", result.Rules[0].Rule.AlertType)
		assert.Equal(t, ">=", result.Rules[0].ThresholdOperator)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with multiple rules", func(t *testing.T) {
		severity := 5
		operator := restapi.LogicalOperatorType("AND")
		aggregation := restapi.Aggregation("MEAN")
		apiObject := &restapi.WebsiteAlertConfig{
			ID:          "api-id-multi",
			Name:        "Multi Rule Alert",
			Description: "Multi Rule Description",
			Severity:    &severity,
			Triggering:  true,
			WebsiteID:   "website-multi",
			TagFilterExpression: &restapi.TagFilter{
				Type:            "EXPRESSION",
				LogicalOperator: &operator,
				Elements:        []*restapi.TagFilter{},
			},
			AlertChannelIDs: []string{"channel-f"},
			Granularity:     600000,
			TimeThreshold: restapi.WebsiteTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: func() *int64 { v := int64(300000); return &v }(),
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
			Rules: []restapi.WebsiteAlertRuleWithThresholds{
				{
					Rule: &restapi.WebsiteAlertRule{
						AlertType:   "slowness",
						MetricName:  "latency",
						Aggregation: &aggregation,
					},
					ThresholdOperator: ">=",
					Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
						restapi.WarningSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(1000); return &v }(),
						},
						restapi.CriticalSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(2000); return &v }(),
						},
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model to ensure Rules is an empty slice, not nil
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.Equal(t, "api-id-multi", model.ID.ValueString())
		assert.Equal(t, "Multi Rule Alert", model.Name.ValueString())
		assert.NotNil(t, model.Rules)
		assert.Len(t, model.Rules, 1)
		assert.NotNil(t, model.Rules[0].Rule)
		assert.NotNil(t, model.Rules[0].Rule.Slowness)
		assert.Equal(t, "latency", model.Rules[0].Rule.Slowness.MetricName.ValueString())
		assert.Equal(t, ">=", model.Rules[0].ThresholdOperator.ValueString())
	})

	t.Run("should update state with null tag filter", func(t *testing.T) {
		severity := 5
		aggregation := restapi.Aggregation("MEAN")
		apiObject := &restapi.WebsiteAlertConfig{
			ID:                  "api-id-no-tag",
			Name:                "No Tag Filter",
			Description:         "No Tag Filter Description",
			Severity:            &severity,
			Triggering:          true,
			WebsiteID:           "website-no-tag",
			TagFilterExpression: nil,
			AlertChannelIDs:     []string{"channel-g"},
			Granularity:         600000,
			TimeThreshold: restapi.WebsiteTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: func() *int64 { v := int64(300000); return &v }(),
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
			Rules: []restapi.WebsiteAlertRuleWithThresholds{
				{
					Rule: &restapi.WebsiteAlertRule{
						AlertType:   "slowness",
						MetricName:  "latency",
						Aggregation: &aggregation,
					},
					ThresholdOperator: ">=",
					Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
						restapi.WarningSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(1000); return &v }(),
						},
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)

		assert.False(t, diags.HasError())

		var model WebsiteAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())

		assert.True(t, model.TagFilter.IsNull())
	})
}

func TestMapStateToDataObject_NullSeverity(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should handle null severity", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules:       []RuleWithThresholdPluginModel{},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInSequence: &WebsiteViolationsInSequenceModel{
					TimeWindow: types.Int64Value(300000),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Nil(t, result.Severity)
	})
}

func TestMapStateToDataObject_MissingTimeThreshold(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should return error when time threshold is missing", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:            types.StringValue("test-id"),
			Name:          types.StringValue("Test"),
			Description:   types.StringValue("Desc"),
			Triggering:    types.BoolValue(true),
			WebsiteID:     types.StringValue("website-1"),
			TagFilter:     types.StringNull(),
			Granularity:   types.Int64Value(600000),
			Rules:         []RuleWithThresholdPluginModel{},
			TimeThreshold: nil,
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
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
}

// Helper functions

func createMockState(t *testing.T, model WebsiteAlertConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewWebsiteAlertConfigResourceHandleFramework()
	return resource.MetaData().Schema
}

func TestMapStateToDataObject_InvalidSeverity(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should handle missing severity field", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules:       []RuleWithThresholdPluginModel{},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInSequence: &WebsiteViolationsInSequenceModel{
					TimeWindow: types.Int64Value(300000),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Nil(t, result.Severity)
	})
}

func TestUpdateState_InvalidSeverity(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should return error for invalid severity from API", func(t *testing.T) {
		invalidSeverity := 999
		aggregation := restapi.Aggregation("MEAN")
		apiObject := &restapi.WebsiteAlertConfig{
			ID:          "api-id-invalid",
			Name:        "Invalid Severity",
			Description: "Invalid Severity Description",
			Severity:    &invalidSeverity,
			Triggering:  true,
			WebsiteID:   "website-invalid",
			Granularity: 600000,
			TimeThreshold: restapi.WebsiteTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: func() *int64 { v := int64(300000); return &v }(),
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
			Rules: []restapi.WebsiteAlertRuleWithThresholds{
				{
					Rule: &restapi.WebsiteAlertRule{
						AlertType:   "slowness",
						MetricName:  "latency",
						Aggregation: &aggregation,
					},
					ThresholdOperator: ">=",
					Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
						restapi.WarningSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(1000); return &v }(),
						},
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.True(t, diags.HasError())
	})
}

func TestMapStateToDataObject_InvalidTimeThresholdConfig(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should return error for invalid time threshold configuration", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules:       []RuleWithThresholdPluginModel{},
			TimeThreshold: &WebsiteTimeThresholdModel{
				// All time threshold types are nil - invalid configuration
				ViolationsInSequence:             nil,
				ViolationsInPeriod:               nil,
				UserImpactOfViolationsInSequence: nil,
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		_, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.True(t, resultDiags.HasError())
		assert.Contains(t, resultDiags[0].Summary(), "Invalid time threshold configuration")
	})
}

func TestMapStateToDataObject_WithTagFilter(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should parse valid tag filter expression", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringValue("entity.type EQUALS 'website'"),
			Granularity: types.Int64Value(600000),
			Rules:       []RuleWithThresholdPluginModel{},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInSequence: &WebsiteViolationsInSequenceModel{
					TimeWindow: types.Int64Value(300000),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
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
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringValue("invalid{{{syntax"),
			Granularity: types.Int64Value(600000),
			Rules:       []RuleWithThresholdPluginModel{},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInSequence: &WebsiteViolationsInSequenceModel{
					TimeWindow: types.Int64Value(300000),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
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
}

func TestMapStateToDataObject_AllRuleTypes(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should map throughput rule", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules: []RuleWithThresholdPluginModel{
				{
					Rule: &WebsiteAlertRuleModel{
						Throughput: &WebsiteAlertRuleConfigModel{
							MetricName:  types.StringValue("beaconCount"),
							Aggregation: types.StringValue("SUM"),
						},
					},
					ThresholdOperator: types.StringValue(">="),
					Thresholds: &shared.ThresholdAllPluginModel{
						Warning: &shared.ThresholdAllTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float32Value(float32(100)),
							},
						},
					},
				},
			},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInSequence: &WebsiteViolationsInSequenceModel{
					TimeWindow: types.Int64Value(300000),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Rules, 1)
		assert.Equal(t, "throughput", result.Rules[0].Rule.AlertType)
	})

	t.Run("should map status code rule", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules: []RuleWithThresholdPluginModel{
				{
					Rule: &WebsiteAlertRuleModel{
						StatusCode: &WebsiteAlertRuleConfigCompleteModel{
							MetricName:  types.StringValue("httpStatusCode"),
							Aggregation: types.StringValue("SUM"),
							Operator:    types.StringValue("EQUALS"),
							Value:       types.StringValue("500"),
						},
					},
					ThresholdOperator: types.StringValue(">"),
					Thresholds: &shared.ThresholdAllPluginModel{
						Critical: &shared.ThresholdAllTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float32Value(float32(10)),
							},
						},
					},
				},
			},
			TimeThreshold: &WebsiteTimeThresholdModel{
				ViolationsInPeriod: &WebsiteViolationsInPeriodModel{
					TimeWindow: types.Int64Value(600000),
					Violations: types.Int64Value(5),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Rules, 1)
		assert.Equal(t, "statusCode", result.Rules[0].Rule.AlertType)
		assert.Equal(t, "500", *result.Rules[0].Rule.Value)
	})

	t.Run("should map specific js error rule", func(t *testing.T) {
		model := WebsiteAlertConfigModel{
			ID:          types.StringValue("test-id"),
			Name:        types.StringValue("Test"),
			Description: types.StringValue("Desc"),
			Triggering:  types.BoolValue(true),
			WebsiteID:   types.StringValue("website-1"),
			TagFilter:   types.StringNull(),
			Granularity: types.Int64Value(600000),
			Rules: []RuleWithThresholdPluginModel{
				{
					Rule: &WebsiteAlertRuleModel{
						SpecificJsError: &WebsiteAlertRuleConfigCompleteModel{
							MetricName:  types.StringValue("jsErrors"),
							Aggregation: types.StringValue("SUM"),
							Operator:    types.StringValue("CONTAINS"),
							Value:       types.StringValue("TypeError"),
						},
					},
					ThresholdOperator: types.StringValue(">="),
					Thresholds: &shared.ThresholdAllPluginModel{
						Warning: &shared.ThresholdAllTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float32Value(float32(1)),
							},
						},
					},
				},
			},
			TimeThreshold: &WebsiteTimeThresholdModel{
				UserImpactOfViolationsInSequence: &WebsiteUserImpactOfViolationsInSequenceModel{
					TimeWindow:              types.Int64Value(300000),
					ImpactMeasurementMethod: types.StringValue("AGGREGATED"),
					UserPercentage:          types.Float64Value(10.5),
					Users:                   types.Int64Value(100),
				},
			},
		}

		alertChannelIds, _ := types.SetValueFrom(ctx, types.StringType, []string{"channel-1"})
		model.AlertChannelIDs = alertChannelIds
		model.CustomPayloadFields = types.ListNull(shared.GetCustomPayloadFieldType())

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)

		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Len(t, result.Rules, 1)
		assert.Equal(t, "specificJsError", result.Rules[0].Rule.AlertType)
		assert.Equal(t, "TypeError", *result.Rules[0].Rule.Value)
		assert.Equal(t, "userImpactOfViolationsInSequence", result.TimeThreshold.Type)
		assert.NotNil(t, result.TimeThreshold.UserPercentage)
		assert.NotNil(t, result.TimeThreshold.Users)
	})
}

func TestUpdateState_AllRuleTypes(t *testing.T) {
	resource := &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaWebsiteAlertConfigFramework,
			Schema:        NewWebsiteAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	t.Run("should update state with throughput rule", func(t *testing.T) {
		severity := 5
		aggregation := restapi.Aggregation("SUM")
		apiObject := &restapi.WebsiteAlertConfig{
			ID:          "api-id",
			Name:        "Throughput Alert",
			Description: "Throughput Description",
			Severity:    &severity,
			Triggering:  true,
			WebsiteID:   "website-1",
			Granularity: 600000,
			TimeThreshold: restapi.WebsiteTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: func() *int64 { v := int64(300000); return &v }(),
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
			Rules: []restapi.WebsiteAlertRuleWithThresholds{
				{
					Rule: &restapi.WebsiteAlertRule{
						AlertType:   "throughput",
						MetricName:  "beaconCount",
						Aggregation: &aggregation,
					},
					ThresholdOperator: ">=",
					Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
						restapi.WarningSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(100); return &v }(),
						},
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model WebsiteAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.Len(t, model.Rules, 1)
		assert.NotNil(t, model.Rules[0].Rule.Throughput)
	})

	t.Run("should update state with status code rule", func(t *testing.T) {
		severity := 10
		aggregation := restapi.Aggregation("SUM")
		operator := restapi.ExpressionOperator("EQUALS")
		value := "500"
		apiObject := &restapi.WebsiteAlertConfig{
			ID:          "api-id",
			Name:        "Status Code Alert",
			Description: "Status Code Description",
			Severity:    &severity,
			Triggering:  false,
			WebsiteID:   "website-1",
			Granularity: 600000,
			TimeThreshold: restapi.WebsiteTimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: func() *int64 { v := int64(600000); return &v }(),
				Violations: func() *int32 { v := int32(5); return &v }(),
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
			Rules: []restapi.WebsiteAlertRuleWithThresholds{
				{
					Rule: &restapi.WebsiteAlertRule{
						AlertType:   "statusCode",
						MetricName:  "httpStatusCode",
						Aggregation: &aggregation,
						Operator:    &operator,
						Value:       &value,
					},
					ThresholdOperator: ">",
					Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
						restapi.CriticalSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(10); return &v }(),
						},
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model WebsiteAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.Len(t, model.Rules, 1)
		assert.NotNil(t, model.Rules[0].Rule.StatusCode)
		assert.NotNil(t, model.TimeThreshold.ViolationsInPeriod)
	})

	t.Run("should update state with specific js error rule", func(t *testing.T) {
		severity := 5
		aggregation := restapi.Aggregation("SUM")
		operator := restapi.ExpressionOperator("CONTAINS")
		value := "TypeError"
		apiObject := &restapi.WebsiteAlertConfig{
			ID:          "api-id",
			Name:        "JS Error Alert",
			Description: "JS Error Description",
			Severity:    &severity,
			Triggering:  true,
			WebsiteID:   "website-1",
			Granularity: 600000,
			TimeThreshold: restapi.WebsiteTimeThreshold{
				Type:       "userImpactOfViolationsInSequence",
				TimeWindow: func() *int64 { v := int64(300000); return &v }(),
				ImpactMeasurementMethod: func() *restapi.WebsiteImpactMeasurementMethod {
					v := restapi.WebsiteImpactMeasurementMethod("AGGREGATED")
					return &v
				}(),
				UserPercentage: func() *float64 { v := float64(10.5); return &v }(),
				Users:          func() *int32 { v := int32(100); return &v }(),
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
			Rules: []restapi.WebsiteAlertRuleWithThresholds{
				{
					Rule: &restapi.WebsiteAlertRule{
						AlertType:   "specificJsError",
						MetricName:  "jsErrors",
						Aggregation: &aggregation,
						Operator:    &operator,
						Value:       &value,
					},
					ThresholdOperator: ">=",
					Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
						restapi.WarningSeverity: {
							Type:  "staticThreshold",
							Value: func() *float64 { v := float64(1); return &v }(),
						},
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model WebsiteAlertConfigModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.Len(t, model.Rules, 1)
		assert.NotNil(t, model.Rules[0].Rule.SpecificJsError)
		assert.NotNil(t, model.TimeThreshold.UserImpactOfViolationsInSequence)
	})
}

// initializeEmptyState initializes the state with an empty model to ensure proper state initialization
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := WebsiteAlertConfigModel{
		ID:                  types.StringNull(),
		Name:                types.StringNull(),
		Description:         types.StringNull(),
		Triggering:          types.BoolNull(),
		WebsiteID:           types.StringNull(),
		TagFilter:           types.StringNull(),
		AlertChannelIDs:     types.SetNull(types.StringType),
		Granularity:         types.Int64Null(),
		CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
		TimeThreshold:       nil,
		Rules:               []RuleWithThresholdPluginModel{},
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}

// Made with Bob
