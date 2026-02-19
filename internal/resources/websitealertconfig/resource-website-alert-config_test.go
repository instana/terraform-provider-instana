package websitealertconfig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/util"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebsiteAlertConfigResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewWebsiteAlertConfigResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaWebsiteAlertConfig, metadata.ResourceName)
		assert.Equal(t, int64(2), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewWebsiteAlertConfigResourceHandle()
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
		resource := &websiteAlertConfigResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaWebsiteAlertConfig,
				SchemaVersion: 1,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaWebsiteAlertConfig, metadata.ResourceName)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return website alert config rest resource", func(t *testing.T) {
		resource := &websiteAlertConfigResource{}

		mockAPI := &mockWebsiteAlertAPI{}
		restResource := resource.GetRestResource(mockAPI)

		assert.NotNil(t, restResource)
	})
}

// mockWebsiteAlertAPI extends the common mock to provide specific behavior for website alert config tests
type mockWebsiteAlertAPI struct {
	testutils.MockInstanaAPI
}

func (m *mockWebsiteAlertAPI) WebsiteAlertConfig() restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return &mockWebsiteAlertConfigRestResource{}
}

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
		resource := &websiteAlertConfigResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaWebsiteAlertConfig,
				Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
								Value: types.Float64Value(1000),
							},
						},
						Critical: &shared.ThresholdAllTypeModel{
							Static: &shared.StaticTypeModel{
								Value: types.Float64Value(2000),
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := NewWebsiteAlertConfigResourceHandle()
	return resource.MetaData().Schema
}

func TestMapStateToDataObject_InvalidSeverity(t *testing.T) {
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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
								Value: types.Float64Value(100),
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
								Value: types.Float64Value(10),
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
								Value: types.Float64Value(1),
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
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

// Helper function to create bool pointer
func ptrBool(v bool) *bool {
	return &v
}

// Helper function to create int64 pointer
func ptrInt64(v int64) *int64 {
	return &v
}

func TestExtractEnabledFlag(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Bool
		expected *bool
	}{
		{
			name:     "explicit true",
			input:    types.BoolValue(true),
			expected: ptrBool(true),
		},
		{
			name:     "explicit false",
			input:    types.BoolValue(false),
			expected: ptrBool(false),
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	tests := []struct {
		name            string
		enabledValue    types.Bool
		expectedEnabled *bool
		description     string
	}{
		{
			name:            "enabled explicitly set to true",
			enabledValue:    types.BoolValue(true),
			expectedEnabled: ptrBool(true),
			description:     "When user explicitly sets enabled=true, API should receive true",
		},
		{
			name:            "enabled explicitly set to false",
			enabledValue:    types.BoolValue(false),
			expectedEnabled: ptrBool(false),
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
			model := WebsiteAlertConfigModel{
				ID:                  types.StringValue("test-id"),
				Name:                types.StringValue("Test Alert"),
				Description:         types.StringValue("Test Description"),
				Triggering:          types.BoolValue(false),
				Enabled:             tt.enabledValue,
				WebsiteID:           types.StringValue("website-1"),
				TagFilter:           types.StringNull(),
				Granularity:         types.Int64Value(600000),
				AlertChannelIDs:     types.SetNull(types.StringType),
				CustomPayloadFields: types.ListNull(shared.GetCustomPayloadFieldType()),
				Rules:               []RuleWithThresholdPluginModel{},
				TimeThreshold: &WebsiteTimeThresholdModel{
					ViolationsInSequence: &WebsiteViolationsInSequenceModel{
						TimeWindow: types.Int64Value(300000),
					},
				},
			}

			state := &tfsdk.State{
				Schema: resource.metaData.Schema,
			}
			diags := state.Set(ctx, model)
			require.False(t, diags.HasError())

			result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
			require.False(t, resultDiags.HasError(), "Expected no errors for %s", tt.name)
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
	resource := &websiteAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaWebsiteAlertConfig,
			Schema:        NewWebsiteAlertConfigResourceHandle().MetaData().Schema,
			SchemaVersion: 1,
		},
	}
	ctx := context.Background()

	tests := []struct {
		name            string
		apiEnabled      *bool
		expectedEnabled bool
		description     string
	}{
		{
			name:            "API returns enabled=true",
			apiEnabled:      ptrBool(true),
			expectedEnabled: true,
			description:     "When API returns enabled=true, state should have enabled=true",
		},
		{
			name:            "API returns enabled=false",
			apiEnabled:      ptrBool(false),
			expectedEnabled: false,
			description:     "When API returns enabled=false, state should have enabled=false",
		},
		{
			name:            "API returns nil - should use default true",
			apiEnabled:      nil,
			expectedEnabled: WebsiteAlertConfigDefaultEnabled,
			description:     "When API returns nil, state should use default value (true)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiObject := &restapi.WebsiteAlertConfig{
				ID:          "test-id",
				Name:        "Test Alert",
				Description: "Test Description",
				Triggering:  false,
				Enabled:     tt.apiEnabled,
				WebsiteID:   "website-1",
				Granularity: 600000,
				TimeThreshold: restapi.WebsiteTimeThreshold{
					Type:       "violationsInSequence",
					TimeWindow: ptrInt64(300000),
				},
				CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
				Rules:                 []restapi.WebsiteAlertRuleWithThresholds{},
			}

			state := &tfsdk.State{
				Schema: resource.metaData.Schema,
			}

			// Initialize state with empty model
			initializeEmptyState(t, ctx, state)

			diags := resource.UpdateState(ctx, state, nil, apiObject)
			require.False(t, diags.HasError(), "Expected no errors for %s", tt.name)

			var model WebsiteAlertConfigModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError(), "Expected no errors getting state for %s", tt.name)

			assert.Equal(t, tt.expectedEnabled, model.Enabled.ValueBool(), "%s: %s", tt.name, tt.description)
		})
	}
}
