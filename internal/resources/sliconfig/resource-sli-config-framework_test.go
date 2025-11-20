package sliconfig

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSliConfigResourceHandleFramework(t *testing.T) {
	resource := NewSliConfigResourceHandleFramework()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaSliConfigFramework, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(1), metaData.SchemaVersion)
	assert.True(t, metaData.CreateOnly)
}

func TestMetaData(t *testing.T) {
	resource := &sliConfigResourceFramework{
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
	resource := NewSliConfigResourceHandleFramework()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestGetRestResource(t *testing.T) {
	resource := &sliConfigResourceFramework{}
	assert.NotNil(t, resource.GetRestResource)
}

// UpdateState Tests - Application Time Based

func TestUpdateState_ApplicationTimeBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:          "application",
			ApplicationID: &appID,
			BoundaryScope: &boundaryScope,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "sli-id-1", model.ID.ValueString())
	assert.Equal(t, "Test SLI", model.Name.ValueString())
	assert.Equal(t, int64(1234567890), model.InitialEvaluationTimestamp.ValueInt64())
	require.NotNil(t, model.SliEntity)
	require.NotNil(t, model.SliEntity.ApplicationTimeBased)
	assert.Equal(t, "app-123", model.SliEntity.ApplicationTimeBased.ApplicationID.ValueString())
	assert.Equal(t, "ALL", model.SliEntity.ApplicationTimeBased.BoundaryScope.ValueString())
}

func TestUpdateState_ApplicationTimeBased_WithServiceAndEndpoint(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "INBOUND"
	serviceID := "service-456"
	endpointID := "endpoint-789"

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:          "application",
			ApplicationID: &appID,
			BoundaryScope: &boundaryScope,
			ServiceID:     &serviceID,
			EndpointID:    &endpointID,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.ApplicationTimeBased)
	assert.Equal(t, "service-456", model.SliEntity.ApplicationTimeBased.ServiceID.ValueString())
	assert.Equal(t, "endpoint-789", model.SliEntity.ApplicationTimeBased.EndpointID.ValueString())
}

func TestUpdateState_ApplicationTimeBased_WithMetricConfiguration(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI",
		InitialEvaluationTimestamp: 1234567890,
		MetricConfiguration: &restapi.MetricConfiguration{
			Name:        "latency",
			Aggregation: "P95",
			Threshold:   500.5,
		},
		SliEntity: restapi.SliEntity{
			Type:          "application",
			ApplicationID: &appID,
			BoundaryScope: &boundaryScope,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.MetricConfiguration)
	assert.Equal(t, "latency", model.MetricConfiguration.MetricName.ValueString())
	assert.Equal(t, "P95", model.MetricConfiguration.Aggregation.ValueString())
	assert.Equal(t, 500.5, model.MetricConfiguration.Threshold.ValueFloat64())
}

// UpdateState Tests - Application Event Based

func TestUpdateState_ApplicationEventBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"
	
	goodFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "200"),
	})
	
	badFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "500"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI Availability",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "availability",
			ApplicationID:             &appID,
			BoundaryScope:             &boundaryScope,
			GoodEventFilterExpression: goodFilter,
			BadEventFilterExpression:  badFilter,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity)
	require.NotNil(t, model.SliEntity.ApplicationEventBased)
	assert.Equal(t, "app-123", model.SliEntity.ApplicationEventBased.ApplicationID.ValueString())
	assert.Equal(t, "ALL", model.SliEntity.ApplicationEventBased.BoundaryScope.ValueString())
	assert.False(t, model.SliEntity.ApplicationEventBased.GoodEventFilterExpression.IsNull())
	assert.False(t, model.SliEntity.ApplicationEventBased.BadEventFilterExpression.IsNull())
}

func TestUpdateState_ApplicationEventBased_WithIncludeFlags(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "INBOUND"
	includeInternal := true
	includeSynthetic := false
	
	goodFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "200"),
	})
	
	badFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "500"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI Availability",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "availability",
			ApplicationID:             &appID,
			BoundaryScope:             &boundaryScope,
			GoodEventFilterExpression: goodFilter,
			BadEventFilterExpression:  badFilter,
			IncludeInternal:           &includeInternal,
			IncludeSynthetic:          &includeSynthetic,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.ApplicationEventBased)
	assert.True(t, model.SliEntity.ApplicationEventBased.IncludeInternal.ValueBool())
	assert.False(t, model.SliEntity.ApplicationEventBased.IncludeSynthetic.ValueBool())
}

func TestUpdateState_ApplicationEventBased_WithNullIncludeFlags(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"
	
	goodFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "200"),
	})
	
	badFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "500"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI Availability",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "availability",
			ApplicationID:             &appID,
			BoundaryScope:             &boundaryScope,
			GoodEventFilterExpression: goodFilter,
			BadEventFilterExpression:  badFilter,
			IncludeInternal:           nil,
			IncludeSynthetic:          nil,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.ApplicationEventBased)
	// Null flags should default to false
	assert.False(t, model.SliEntity.ApplicationEventBased.IncludeInternal.ValueBool())
	assert.False(t, model.SliEntity.ApplicationEventBased.IncludeSynthetic.ValueBool())
}

// UpdateState Tests - Website Event Based

func TestUpdateState_WebsiteEventBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	websiteID := "website-123"
	beaconType := "pageLoad"
	
	goodFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "beacon.page.name", restapi.EqualsOperator, "home"),
	})
	
	badFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "beacon.error", restapi.EqualsOperator, "true"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test Website SLI",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "websiteEventBased",
			WebsiteId:                 &websiteID,
			BeaconType:                &beaconType,
			GoodEventFilterExpression: goodFilter,
			BadEventFilterExpression:  badFilter,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity)
	require.NotNil(t, model.SliEntity.WebsiteEventBased)
	assert.Equal(t, "website-123", model.SliEntity.WebsiteEventBased.WebsiteID.ValueString())
	assert.Equal(t, "pageLoad", model.SliEntity.WebsiteEventBased.BeaconType.ValueString())
	assert.False(t, model.SliEntity.WebsiteEventBased.GoodEventFilterExpression.IsNull())
	assert.False(t, model.SliEntity.WebsiteEventBased.BadEventFilterExpression.IsNull())
}

// UpdateState Tests - Website Time Based

func TestUpdateState_WebsiteTimeBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	websiteID := "website-123"
	beaconType := "pageLoad"
	
	filterExpr := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "beacon.page.name", restapi.EqualsOperator, "home"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test Website Time SLI",
		InitialEvaluationTimestamp: 1234567890,
		MetricConfiguration: &restapi.MetricConfiguration{
			Name:        "beacon.page.load.time",
			Aggregation: "P95",
			Threshold:   2000.0,
		},
		SliEntity: restapi.SliEntity{
			Type:             "websiteTimeBased",
			WebsiteId:        &websiteID,
			BeaconType:       &beaconType,
			FilterExpression: filterExpr,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity)
	require.NotNil(t, model.SliEntity.WebsiteTimeBased)
	assert.Equal(t, "website-123", model.SliEntity.WebsiteTimeBased.WebsiteID.ValueString())
	assert.Equal(t, "pageLoad", model.SliEntity.WebsiteTimeBased.BeaconType.ValueString())
	assert.False(t, model.SliEntity.WebsiteTimeBased.FilterExpression.IsNull())
}

func TestUpdateState_WebsiteTimeBased_WithoutFilterExpression(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	websiteID := "website-123"
	beaconType := "httpRequest"

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test Website Time SLI",
		InitialEvaluationTimestamp: 1234567890,
		MetricConfiguration: &restapi.MetricConfiguration{
			Name:        "beacon.http.duration",
			Aggregation: "MEAN",
			Threshold:   1000.0,
		},
		SliEntity: restapi.SliEntity{
			Type:             "websiteTimeBased",
			WebsiteId:        &websiteID,
			BeaconType:       &beaconType,
			FilterExpression: nil,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.WebsiteTimeBased)
	assert.True(t, model.SliEntity.WebsiteTimeBased.FilterExpression.IsNull())
}

// UpdateState Tests - Error Cases

func TestUpdateState_UnsupportedEntityType(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type: "unsupported_type",
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	assert.True(t, diags.HasError())
}

func TestUpdateState_WithoutMetricConfiguration(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI",
		InitialEvaluationTimestamp: 1234567890,
		MetricConfiguration:        nil,
		SliEntity: restapi.SliEntity{
			Type:          "application",
			ApplicationID: &appID,
			BoundaryScope: &boundaryScope,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.MetricConfiguration)
}

// MapStateToDataObject Tests - Application Time Based

func TestMapStateToDataObject_ApplicationTimeBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("ALL"),
				ServiceID:     types.StringNull(),
				EndpointID:    types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "sli-id-1", result.ID)
	assert.Equal(t, "Test SLI", result.Name)
	assert.Equal(t, 1234567890, result.InitialEvaluationTimestamp)
	assert.Equal(t, "application", result.SliEntity.Type)
	assert.Equal(t, "app-123", *result.SliEntity.ApplicationID)
	assert.Equal(t, "ALL", *result.SliEntity.BoundaryScope)
}

func TestMapStateToDataObject_ApplicationTimeBased_WithServiceAndEndpoint(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("INBOUND"),
				ServiceID:     types.StringValue("service-456"),
				EndpointID:    types.StringValue("endpoint-789"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "service-456", *result.SliEntity.ServiceID)
	assert.Equal(t, "endpoint-789", *result.SliEntity.EndpointID)
}

func TestMapStateToDataObject_ApplicationTimeBased_WithMetricConfiguration(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		MetricConfiguration: &MetricConfigurationModel{
			MetricName:  types.StringValue("latency"),
			Aggregation: types.StringValue("P95"),
			Threshold:   types.Float64Value(500.5),
		},
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("ALL"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.NotNil(t, result.MetricConfiguration)
	assert.Equal(t, "latency", result.MetricConfiguration.Name)
	assert.Equal(t, "P95", result.MetricConfiguration.Aggregation)
	assert.Equal(t, 500.5, result.MetricConfiguration.Threshold)
}

// MapStateToDataObject Tests - Application Event Based

func TestMapStateToDataObject_ApplicationEventBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("ALL"),
				GoodEventFilterExpression: types.StringValue("call.http.status EQUALS '200'"),
				BadEventFilterExpression:  types.StringValue("call.http.status EQUALS '500'"),
				IncludeInternal:           types.BoolValue(false),
				IncludeSynthetic:          types.BoolValue(false),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "availability", result.SliEntity.Type)
	assert.Equal(t, "app-123", *result.SliEntity.ApplicationID)
	assert.NotNil(t, result.SliEntity.GoodEventFilterExpression)
	assert.NotNil(t, result.SliEntity.BadEventFilterExpression)
}

func TestMapStateToDataObject_ApplicationEventBased_WithIncludeFlags(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("INBOUND"),
				GoodEventFilterExpression: types.StringValue("call.http.status EQUALS '200'"),
				BadEventFilterExpression:  types.StringValue("call.http.status EQUALS '500'"),
				IncludeInternal:           types.BoolValue(true),
				IncludeSynthetic:          types.BoolValue(true),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.NotNil(t, result.SliEntity.IncludeInternal)
	assert.True(t, *result.SliEntity.IncludeInternal)
	require.NotNil(t, result.SliEntity.IncludeSynthetic)
	assert.True(t, *result.SliEntity.IncludeSynthetic)
}

func TestMapStateToDataObject_ApplicationEventBased_WithNullIncludeFlags(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("ALL"),
				GoodEventFilterExpression: types.StringValue("call.http.status EQUALS '200'"),
				BadEventFilterExpression:  types.StringValue("call.http.status EQUALS '500'"),
				IncludeInternal:           types.BoolNull(),
				IncludeSynthetic:          types.BoolNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	// Null flags should not be set in the API object
	assert.Nil(t, result.SliEntity.IncludeInternal)
	assert.Nil(t, result.SliEntity.IncludeSynthetic)
}

// MapStateToDataObject Tests - Website Event Based

func TestMapStateToDataObject_WebsiteEventBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			WebsiteEventBased: &WebsiteEventBasedModel{
				WebsiteID:                 types.StringValue("website-123"),
				BeaconType:                types.StringValue("pageLoad"),
				GoodEventFilterExpression: types.StringValue("beacon.page.name EQUALS 'home'"),
				BadEventFilterExpression:  types.StringValue("beacon.error EQUALS 'true'"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "websiteEventBased", result.SliEntity.Type)
	assert.Equal(t, "website-123", *result.SliEntity.WebsiteId)
	assert.Equal(t, "pageLoad", *result.SliEntity.BeaconType)
	assert.NotNil(t, result.SliEntity.GoodEventFilterExpression)
	assert.NotNil(t, result.SliEntity.BadEventFilterExpression)
}

// MapStateToDataObject Tests - Website Time Based

func TestMapStateToDataObject_WebsiteTimeBased_Basic(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website Time SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		MetricConfiguration: &MetricConfigurationModel{
			MetricName:  types.StringValue("beacon.page.load.time"),
			Aggregation: types.StringValue("P95"),
			Threshold:   types.Float64Value(2000.0),
		},
		SliEntity: &SliEntityModel{
			WebsiteTimeBased: &WebsiteTimeBasedModel{
				WebsiteID:        types.StringValue("website-123"),
				BeaconType:       types.StringValue("pageLoad"),
				FilterExpression: types.StringValue("beacon.page.name EQUALS 'home'"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "websiteTimeBased", result.SliEntity.Type)
	assert.Equal(t, "website-123", *result.SliEntity.WebsiteId)
	assert.Equal(t, "pageLoad", *result.SliEntity.BeaconType)
	assert.NotNil(t, result.SliEntity.FilterExpression)
}

func TestMapStateToDataObject_WebsiteTimeBased_WithoutFilterExpression(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website Time SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		MetricConfiguration: &MetricConfigurationModel{
			MetricName:  types.StringValue("beacon.http.duration"),
			Aggregation: types.StringValue("MEAN"),
			Threshold:   types.Float64Value(1000.0),
		},
		SliEntity: &SliEntityModel{
			WebsiteTimeBased: &WebsiteTimeBasedModel{
				WebsiteID:        types.StringValue("website-123"),
				BeaconType:       types.StringValue("httpRequest"),
				FilterExpression: types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Nil(t, result.SliEntity.FilterExpression)
}

// MapStateToDataObject Tests - Edge Cases

func TestMapStateToDataObject_WithNullID(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringNull(),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("ALL"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "", result.ID)
}

func TestMapStateToDataObject_WithNullInitialEvaluationTimestamp(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Null(),
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("ALL"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, 0, result.InitialEvaluationTimestamp)
}

func TestMapStateToDataObject_WithNullMetricConfiguration(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		MetricConfiguration:        nil,
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("ALL"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Nil(t, result.MetricConfiguration)
}

func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	model := SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationTimeBased: &ApplicationTimeBasedModel{
				ApplicationID: types.StringValue("app-123"),
				BoundaryScope: types.StringValue("ALL"),
			},
		},
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "sli-id-1", result.ID)
	assert.Equal(t, "Test SLI", result.Name)
}

func TestMapStateToDataObject_ApplicationEventBased_InvalidGoodEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("ALL"),
				GoodEventFilterExpression: types.StringValue("invalid filter syntax"),
				BadEventFilterExpression:  types.StringValue("call.http.status EQUALS '500'"),
			},
		},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_ApplicationEventBased_InvalidBadEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("ALL"),
				GoodEventFilterExpression: types.StringValue("call.http.status EQUALS '200'"),
				BadEventFilterExpression:  types.StringValue("invalid filter syntax"),
			},
		},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_WebsiteEventBased_InvalidGoodEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			WebsiteEventBased: &WebsiteEventBasedModel{
				WebsiteID:                 types.StringValue("website-123"),
				BeaconType:                types.StringValue("pageLoad"),
				GoodEventFilterExpression: types.StringValue("invalid filter syntax"),
				BadEventFilterExpression:  types.StringValue("beacon.error EQUALS 'true'"),
			},
		},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

func TestMapStateToDataObject_WebsiteTimeBased_InvalidFilterExpression(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website Time SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			WebsiteTimeBased: &WebsiteTimeBasedModel{
				WebsiteID:        types.StringValue("website-123"),
				BeaconType:       types.StringValue("pageLoad"),
				FilterExpression: types.StringValue("invalid filter syntax"),
			},
		},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

// Helper functions

func createMockState(t *testing.T, model SliConfigModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewSliConfigResourceHandleFramework()
	return resource.MetaData().Schema
}

// Made with Bob
func TestUpdateState_ApplicationEventBased_WithNullGoodEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"
	
	badFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "500"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI Availability",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "availability",
			ApplicationID:             &appID,
			BoundaryScope:             &boundaryScope,
			GoodEventFilterExpression: nil,
			BadEventFilterExpression:  badFilter,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.ApplicationEventBased)
	assert.True(t, model.SliEntity.ApplicationEventBased.GoodEventFilterExpression.IsNull())
	assert.False(t, model.SliEntity.ApplicationEventBased.BadEventFilterExpression.IsNull())
}

func TestUpdateState_ApplicationEventBased_WithNullBadEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	appID := "app-123"
	boundaryScope := "ALL"
	
	goodFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, "200"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test SLI Availability",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "availability",
			ApplicationID:             &appID,
			BoundaryScope:             &boundaryScope,
			GoodEventFilterExpression: goodFilter,
			BadEventFilterExpression:  nil,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.ApplicationEventBased)
	assert.False(t, model.SliEntity.ApplicationEventBased.GoodEventFilterExpression.IsNull())
	assert.True(t, model.SliEntity.ApplicationEventBased.BadEventFilterExpression.IsNull())
}

func TestUpdateState_WebsiteEventBased_WithNullGoodEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	websiteID := "website-123"
	beaconType := "pageLoad"
	
	badFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "beacon.error", restapi.EqualsOperator, "true"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test Website SLI",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "websiteEventBased",
			WebsiteId:                 &websiteID,
			BeaconType:                &beaconType,
			GoodEventFilterExpression: nil,
			BadEventFilterExpression:  badFilter,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.WebsiteEventBased)
	assert.True(t, model.SliEntity.WebsiteEventBased.GoodEventFilterExpression.IsNull())
	assert.False(t, model.SliEntity.WebsiteEventBased.BadEventFilterExpression.IsNull())
}

func TestUpdateState_WebsiteEventBased_WithNullBadEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	websiteID := "website-123"
	beaconType := "pageLoad"
	
	goodFilter := restapi.NewLogicalAndTagFilter([]*restapi.TagFilter{
		restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "beacon.page.name", restapi.EqualsOperator, "home"),
	})

	data := &restapi.SliConfig{
		ID:                         "sli-id-1",
		Name:                       "Test Website SLI",
		InitialEvaluationTimestamp: 1234567890,
		SliEntity: restapi.SliEntity{
			Type:                      "websiteEventBased",
			WebsiteId:                 &websiteID,
			BeaconType:                &beaconType,
			GoodEventFilterExpression: goodFilter,
			BadEventFilterExpression:  nil,
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model SliConfigModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.NotNil(t, model.SliEntity.WebsiteEventBased)
	assert.False(t, model.SliEntity.WebsiteEventBased.GoodEventFilterExpression.IsNull())
	assert.True(t, model.SliEntity.WebsiteEventBased.BadEventFilterExpression.IsNull())
}

func TestMapStateToDataObject_ApplicationEventBased_WithNullGoodEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("ALL"),
				GoodEventFilterExpression: types.StringNull(),
				BadEventFilterExpression:  types.StringValue("call.http.status EQUALS '500'"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Nil(t, result.SliEntity.GoodEventFilterExpression)
	assert.NotNil(t, result.SliEntity.BadEventFilterExpression)
}

func TestMapStateToDataObject_ApplicationEventBased_WithNullBadEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			ApplicationEventBased: &ApplicationEventBasedModel{
				ApplicationID:             types.StringValue("app-123"),
				BoundaryScope:             types.StringValue("ALL"),
				GoodEventFilterExpression: types.StringValue("call.http.status EQUALS '200'"),
				BadEventFilterExpression:  types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.NotNil(t, result.SliEntity.GoodEventFilterExpression)
	assert.Nil(t, result.SliEntity.BadEventFilterExpression)
}

func TestMapStateToDataObject_WebsiteEventBased_WithNullGoodEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			WebsiteEventBased: &WebsiteEventBasedModel{
				WebsiteID:                 types.StringValue("website-123"),
				BeaconType:                types.StringValue("pageLoad"),
				GoodEventFilterExpression: types.StringNull(),
				BadEventFilterExpression:  types.StringValue("beacon.error EQUALS 'true'"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Nil(t, result.SliEntity.GoodEventFilterExpression)
	assert.NotNil(t, result.SliEntity.BadEventFilterExpression)
}

func TestMapStateToDataObject_WebsiteEventBased_WithNullBadEventFilter(t *testing.T) {
	ctx := context.Background()
	resource := NewSliConfigResourceHandleFramework()

	state := createMockState(t, SliConfigModel{
		ID:                         types.StringValue("sli-id-1"),
		Name:                       types.StringValue("Test Website SLI"),
		InitialEvaluationTimestamp: types.Int64Value(1234567890),
		SliEntity: &SliEntityModel{
			WebsiteEventBased: &WebsiteEventBasedModel{
				WebsiteID:                 types.StringValue("website-123"),
				BeaconType:                types.StringValue("pageLoad"),
				GoodEventFilterExpression: types.StringValue("beacon.page.name EQUALS 'home'"),
				BadEventFilterExpression:  types.StringNull(),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.NotNil(t, result.SliEntity.GoodEventFilterExpression)
	assert.Nil(t, result.SliEntity.BadEventFilterExpression)
}
