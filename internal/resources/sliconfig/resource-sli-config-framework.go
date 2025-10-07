package sliconfig

import (
	"context"
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewSliConfigResourceHandleFramework creates the resource handle for SLI configuration
func NewSliConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.SliConfig] {
	return &sliConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSliConfigFramework,
			Schema: schema.Schema{
				Description: SliConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: SliConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: SliConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					"initial_evaluation_timestamp": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: SliConfigDescInitialEvaluationTimestamp,
					},
					"metric_configuration": schema.SingleNestedAttribute{
						Description: SliConfigDescMetricConfiguration,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"metric_name": schema.StringAttribute{
								Required:    true,
								Description: SliConfigDescMetricName,
							},
							"aggregation": schema.StringAttribute{
								Required:    true,
								Description: SliConfigDescAggregation,
								Validators: []validator.String{
									stringvalidator.OneOf(
										"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTRIBUTION", "DISTINCT_COUNT", "SUM_POSITIVE", "PER_SECOND",
									),
								},
							},
							"threshold": schema.Float64Attribute{
								Required:    true,
								Description: SliConfigDescThreshold,
								Validators: []validator.Float64{
									float64validator.AtLeast(0.000001),
								},
							},
						},
					},
					"sli_entity": schema.SingleNestedAttribute{
						Description: SliConfigDescSliEntity,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"application_time_based": schema.SingleNestedAttribute{
								Description: SliConfigDescApplicationTimeBased,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"application_id": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescApplicationID,
									},
									"service_id": schema.StringAttribute{
										Optional:    true,
										Description: SliConfigDescServiceID,
									},
									"endpoint_id": schema.StringAttribute{
										Optional:    true,
										Description: SliConfigDescEndpointID,
									},
									"boundary_scope": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescBoundaryScope,
										Validators: []validator.String{
											stringvalidator.OneOf("ALL", "INBOUND"),
										},
									},
								},
							},
							"application_event_based": schema.SingleNestedAttribute{
								Description: SliConfigDescApplicationEventBased,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"application_id": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescApplicationID,
									},
									"boundary_scope": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescBoundaryScope,
										Validators: []validator.String{
											stringvalidator.OneOf("ALL", "INBOUND"),
										},
									},
									"bad_event_filter_expression": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescBadEventFilterExpression,
									},
									"good_event_filter_expression": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescGoodEventFilterExpression,
									},
									"include_internal": schema.BoolAttribute{
										Optional:    true,
										Description: SliConfigDescIncludeInternal,
									},
									"include_synthetic": schema.BoolAttribute{
										Optional:    true,
										Description: SliConfigDescIncludeSynthetic,
									},
									"endpoint_id": schema.StringAttribute{
										Optional:    true,
										Description: SliConfigDescEndpointIDAvailability,
									},
									"service_id": schema.StringAttribute{
										Optional:    true,
										Description: SliConfigDescServiceIDAvailability,
									},
								},
							},
							"website_event_based": schema.SingleNestedAttribute{
								Description: SliConfigDescWebsiteEventBased,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"website_id": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescWebsiteID,
									},
									"bad_event_filter_expression": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescBadEventFilterExpression,
									},
									"good_event_filter_expression": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescGoodEventFilterExpression,
									},
									"beacon_type": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescBeaconType,
										Validators: []validator.String{
											stringvalidator.OneOf("pageLoad", "resourceLoad", "httpRequest", "error", "custom", "pageChange"),
										},
									},
								},
							},
							"website_time_based": schema.SingleNestedAttribute{
								Description: SliConfigDescWebsiteTimeBased,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"website_id": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescWebsiteID,
									},
									"filter_expression": schema.StringAttribute{
										Optional:    true,
										Description: SliConfigDescFilterExpression,
									},
									"beacon_type": schema.StringAttribute{
										Required:    true,
										Description: SliConfigDescBeaconType,
										Validators: []validator.String{
											stringvalidator.OneOf("pageLoad", "resourceLoad", "httpRequest", "error", "custom", "pageChange"),
										},
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 1,
			CreateOnly:    true,
		},
	}
}

func (r *sliConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *sliConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SliConfig] {
	return api.SliConfigs()
}

func (r *sliConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sliConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, sliConfig *restapi.SliConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	model := SliConfigModel{
		ID:                         types.StringValue(sliConfig.ID),
		Name:                       types.StringValue(sliConfig.Name),
		InitialEvaluationTimestamp: types.Int64Value(int64(sliConfig.InitialEvaluationTimestamp)),
	}

	// Map metric configuration if present
	if sliConfig.MetricConfiguration != nil {
		metricConfigModel := MetricConfigurationModel{
			MetricName:  types.StringValue(sliConfig.MetricConfiguration.Name),
			Aggregation: types.StringValue(sliConfig.MetricConfiguration.Aggregation),
			Threshold:   types.Float64Value(sliConfig.MetricConfiguration.Threshold),
		}

		metricConfigObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"metric_name": types.StringType,
			"aggregation": types.StringType,
			"threshold":   types.Float64Type,
		}, metricConfigModel)
		if diags.HasError() {
			return diags
		}

		model.MetricConfiguration = metricConfigObj
	} else {
		model.MetricConfiguration = types.ObjectNull(map[string]attr.Type{
			"metric_name": types.StringType,
			"aggregation": types.StringType,
			"threshold":   types.Float64Type,
		})
	}

	// Map SLI entity
	sliEntityModel := SliEntityModel{}
	var entityDiags diag.Diagnostics

	sliEntityModel.ApplicationTimeBased = types.ObjectNull(applicationTimeBasedObjectType.AttrTypes)
	sliEntityModel.ApplicationEventBased = types.ObjectNull(applicationEventBasedObjectType.AttrTypes)
	sliEntityModel.WebsiteEventBased = types.ObjectNull(websiteEventBasedObjectType.AttrTypes)
	sliEntityModel.WebsiteTimeBased = types.ObjectNull(websiteTimeBasedObjectType.AttrTypes)

	switch sliConfig.SliEntity.Type {
	case "application":
		sliEntityModel.ApplicationTimeBased, entityDiags = r.mapApplicationTimeBasedToStateObject(ctx, sliConfig.SliEntity)
	case "availability":
		sliEntityModel.ApplicationEventBased, entityDiags = r.mapApplicationEventBasedToStateObject(ctx, sliConfig.SliEntity)
	case "websiteEventBased":
		sliEntityModel.WebsiteEventBased, entityDiags = r.mapWebsiteEventBasedToStateObject(ctx, sliConfig.SliEntity)
	case "websiteTimeBased":
		sliEntityModel.WebsiteTimeBased, entityDiags = r.mapWebsiteTimeBasedToStateObject(ctx, sliConfig.SliEntity)
	default:
		diags.AddError(
			SliConfigErrUnsupportedEntityType,
			fmt.Sprintf(SliConfigErrUnsupportedEntityTypeMsg, sliConfig.SliEntity.Type),
		)
		return diags
	}

	if entityDiags.HasError() {
		diags.Append(entityDiags...)
		return diags
	}

	// Create SLI entity object
	model.SliEntity = types.ObjectValueMust(
		map[string]attr.Type{
			"application_time_based":  types.ObjectType{AttrTypes: applicationTimeBasedObjectType.AttrTypes},
			"application_event_based": types.ObjectType{AttrTypes: applicationEventBasedObjectType.AttrTypes},
			"website_event_based":     types.ObjectType{AttrTypes: websiteEventBasedObjectType.AttrTypes},
			"website_time_based":      types.ObjectType{AttrTypes: websiteTimeBasedObjectType.AttrTypes},
		},
		map[string]attr.Value{
			"application_time_based":  sliEntityModel.ApplicationTimeBased,
			"application_event_based": sliEntityModel.ApplicationEventBased,
			"website_event_based":     sliEntityModel.WebsiteEventBased,
			"website_time_based":      sliEntityModel.WebsiteTimeBased,
		},
	)

	// Set the state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *sliConfigResourceFramework) mapApplicationTimeBasedToStateObject(ctx context.Context, sliEntity restapi.SliEntity) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	appTimeBasedModel := ApplicationTimeBasedModel{
		ApplicationID: util.SetStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: util.SetStringPointerToState(sliEntity.BoundaryScope),
	}

	appTimeBasedModel.ServiceID = util.SetStringPointerToState(sliEntity.ServiceID)

	appTimeBasedModel.EndpointID = util.SetStringPointerToState(sliEntity.EndpointID)

	appTimeBasedObj, objDiags := types.ObjectValueFrom(ctx, applicationTimeBasedObjectType.AttrTypes, appTimeBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ObjectNull(applicationTimeBasedObjectType.AttrTypes), diags
	}

	return appTimeBasedObj, diags
}

func (r *sliConfigResourceFramework) mapApplicationEventBasedToStateObject(ctx context.Context, sliEntity restapi.SliEntity) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	appEventBasedModel := ApplicationEventBasedModel{
		ApplicationID: util.SetStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: util.SetStringPointerToState(sliEntity.BoundaryScope),
		EndpointID:    util.SetStringPointerToState(sliEntity.EndpointID),
		ServiceID:     util.SetStringPointerToState(sliEntity.ServiceID),
	}

	// Map good event filter expression
	if sliEntity.GoodEventFilterExpression != nil {
		goodEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.GoodEventFilterExpression)
		if err != nil {
			diags.AddError(
				SliConfigErrMappingGoodEventFilter,
				fmt.Sprintf(SliConfigErrMappingGoodEventFilterMsg, err),
			)
			return types.ObjectNull(applicationEventBasedObjectType.AttrTypes), diags
		}
		appEventBasedModel.GoodEventFilterExpression = util.SetStringPointerToState(goodEventFilterStr)
	} else {
		appEventBasedModel.GoodEventFilterExpression = types.StringNull()
	}

	// Map bad event filter expression
	if sliEntity.BadEventFilterExpression != nil {
		badEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.BadEventFilterExpression)
		if err != nil {
			diags.AddError(
				SliConfigErrMappingBadEventFilter,
				fmt.Sprintf(SliConfigErrMappingBadEventFilterMsg, err),
			)
			return types.ObjectNull(applicationEventBasedObjectType.AttrTypes), diags
		}
		appEventBasedModel.BadEventFilterExpression = util.SetStringPointerToState(badEventFilterStr)
	} else {
		appEventBasedModel.BadEventFilterExpression = types.StringNull()
	}

	// Map include internal and synthetic flags
	if sliEntity.IncludeInternal != nil {
		appEventBasedModel.IncludeInternal = types.BoolValue(*sliEntity.IncludeInternal)
	} else {
		appEventBasedModel.IncludeInternal = types.BoolValue(false)
	}

	if sliEntity.IncludeSynthetic != nil {
		appEventBasedModel.IncludeSynthetic = types.BoolValue(*sliEntity.IncludeSynthetic)
	} else {
		appEventBasedModel.IncludeSynthetic = types.BoolValue(false)
	}

	appEventBasedObj, objDiags := types.ObjectValueFrom(ctx, applicationEventBasedObjectType.AttrTypes, appEventBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ObjectNull(applicationEventBasedObjectType.AttrTypes), diags
	}

	return appEventBasedObj, diags
}

func (r *sliConfigResourceFramework) mapWebsiteEventBasedToStateObject(ctx context.Context, sliEntity restapi.SliEntity) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	websiteEventBasedModel := WebsiteEventBasedModel{
		WebsiteID:  util.SetStringPointerToState(sliEntity.WebsiteId),
		BeaconType: util.SetStringPointerToState(sliEntity.BeaconType),
	}

	// Map good event filter expression
	if sliEntity.GoodEventFilterExpression != nil {
		goodEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.GoodEventFilterExpression)
		if err != nil {
			diags.AddError(
				SliConfigErrMappingGoodEventFilter,
				fmt.Sprintf(SliConfigErrMappingGoodEventFilterMsg, err),
			)
			return types.ObjectNull(websiteEventBasedObjectType.AttrTypes), diags
		}
		websiteEventBasedModel.GoodEventFilterExpression = util.SetStringPointerToState(goodEventFilterStr)
	} else {
		websiteEventBasedModel.GoodEventFilterExpression = types.StringNull()
	}

	// Map bad event filter expression
	if sliEntity.BadEventFilterExpression != nil {
		badEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.BadEventFilterExpression)
		if err != nil {
			diags.AddError(
				SliConfigErrMappingBadEventFilter,
				fmt.Sprintf(SliConfigErrMappingBadEventFilterMsg, err),
			)
			return types.ObjectNull(websiteEventBasedObjectType.AttrTypes), diags
		}
		websiteEventBasedModel.BadEventFilterExpression = util.SetStringPointerToState(badEventFilterStr)
	} else {
		websiteEventBasedModel.BadEventFilterExpression = types.StringNull()
	}

	websiteEventBasedObj, objDiags := types.ObjectValueFrom(ctx, websiteEventBasedObjectType.AttrTypes, websiteEventBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ObjectNull(websiteEventBasedObjectType.AttrTypes), diags
	}

	return websiteEventBasedObj, diags
}

func (r *sliConfigResourceFramework) mapWebsiteTimeBasedToStateObject(ctx context.Context, sliEntity restapi.SliEntity) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	websiteTimeBasedModel := WebsiteTimeBasedModel{
		WebsiteID:  util.SetStringPointerToState(sliEntity.WebsiteId),
		BeaconType: util.SetStringPointerToState(sliEntity.BeaconType),
	}

	// Map filter expression
	if sliEntity.FilterExpression != nil {
		filterExprStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.FilterExpression)
		if err != nil {
			diags.AddError(
				SliConfigErrMappingFilterExpression,
				fmt.Sprintf(SliConfigErrMappingFilterExpressionMsg, err),
			)
			return types.ObjectNull(websiteTimeBasedObjectType.AttrTypes), diags
		}
		websiteTimeBasedModel.FilterExpression = util.SetStringPointerToState(filterExprStr)
	} else {
		websiteTimeBasedModel.FilterExpression = types.StringNull()
	}

	websiteTimeBasedObj, objDiags := types.ObjectValueFrom(ctx, websiteTimeBasedObjectType.AttrTypes, websiteTimeBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ObjectNull(websiteTimeBasedObjectType.AttrTypes), diags
	}

	return websiteTimeBasedObj, diags
}

func (r *sliConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SliConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SliConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map name
	name := model.Name.ValueString()

	// Map initial evaluation timestamp
	initialEvaluationTimestamp := 0
	if !model.InitialEvaluationTimestamp.IsNull() {
		initialEvaluationTimestamp = int(model.InitialEvaluationTimestamp.ValueInt64())
	}

	// Map metric configuration
	var metricConfiguration *restapi.MetricConfiguration
	if !model.MetricConfiguration.IsNull() && !model.MetricConfiguration.IsUnknown() {
		var metricConfigModel MetricConfigurationModel
		diags.Append(model.MetricConfiguration.As(ctx, &metricConfigModel, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		metricConfiguration = &restapi.MetricConfiguration{
			Name:        metricConfigModel.MetricName.ValueString(),
			Aggregation: metricConfigModel.Aggregation.ValueString(),
			Threshold:   metricConfigModel.Threshold.ValueFloat64(),
		}
	}

	// Map SLI entity
	var sliEntity restapi.SliEntity
	var entityErr error

	if !model.SliEntity.IsNull() && !model.SliEntity.IsUnknown() {
		var sliEntityModel SliEntityModel
		diags.Append(model.SliEntity.As(ctx, &sliEntityModel, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		// Check which entity type is set
		if !sliEntityModel.ApplicationTimeBased.IsNull() && !sliEntityModel.ApplicationTimeBased.IsUnknown() {
			sliEntity, entityErr = r.mapApplicationTimeBasedFromStateObject(ctx, sliEntityModel.ApplicationTimeBased)
		} else if !sliEntityModel.ApplicationEventBased.IsNull() && !sliEntityModel.ApplicationEventBased.IsUnknown() {
			sliEntity, entityErr = r.mapApplicationEventBasedFromStateObject(ctx, sliEntityModel.ApplicationEventBased)
		} else if !sliEntityModel.WebsiteEventBased.IsNull() && !sliEntityModel.WebsiteEventBased.IsUnknown() {
			sliEntity, entityErr = r.mapWebsiteEventBasedFromStateObject(ctx, sliEntityModel.WebsiteEventBased)
		} else if !sliEntityModel.WebsiteTimeBased.IsNull() && !sliEntityModel.WebsiteTimeBased.IsUnknown() {
			sliEntity, entityErr = r.mapWebsiteTimeBasedFromStateObject(ctx, sliEntityModel.WebsiteTimeBased)
		}

		if entityErr != nil {
			diags.AddError(
				SliConfigErrMappingSliEntity,
				fmt.Sprintf(SliConfigErrMappingSliEntityMsg, entityErr),
			)
			return nil, diags
		}
	}

	// Create SLI config
	return &restapi.SliConfig{
		ID:                         id,
		Name:                       name,
		InitialEvaluationTimestamp: initialEvaluationTimestamp,
		MetricConfiguration:        metricConfiguration,
		SliEntity:                  sliEntity,
	}, diags
}

func (r *sliConfigResourceFramework) mapApplicationTimeBasedFromStateObject(ctx context.Context, appTimeBasedObj types.Object) (restapi.SliEntity, error) {
	if appTimeBasedObj.IsNull() || appTimeBasedObj.IsUnknown() {
		return restapi.SliEntity{}, errors.New("application time based entity is null or unknown")
	}

	var appTimeBasedModel ApplicationTimeBasedModel
	diags := appTimeBasedObj.As(ctx, &appTimeBasedModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse application time based entity: %v", diags)
	}
	applicationID := appTimeBasedModel.ApplicationID.ValueString()
	boundaryScope := appTimeBasedModel.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          "application",
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if !appTimeBasedModel.ServiceID.IsNull() {
		serviceID := appTimeBasedModel.ServiceID.ValueString()
		entity.ServiceID = &serviceID
	}

	if !appTimeBasedModel.EndpointID.IsNull() {
		endpointID := appTimeBasedModel.EndpointID.ValueString()
		entity.EndpointID = &endpointID
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapTagFilterStringToAPIModel(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *sliConfigResourceFramework) mapApplicationEventBasedFromStateObject(ctx context.Context, appEventBasedObj types.Object) (restapi.SliEntity, error) {
	if appEventBasedObj.IsNull() || appEventBasedObj.IsUnknown() {
		return restapi.SliEntity{}, errors.New("application event based entity is null or unknown")
	}

	var appEventBasedModel ApplicationEventBasedModel
	diags := appEventBasedObj.As(ctx, &appEventBasedModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse application event based entity: %v", diags)
	}
	applicationID := appEventBasedModel.ApplicationID.ValueString()
	boundaryScope := appEventBasedModel.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          "availability",
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if !appEventBasedModel.BadEventFilterExpression.IsNull() {
		badEventFilter, err := r.mapTagFilterStringToAPIModel(appEventBasedModel.BadEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse bad event filter expression: %v", err)
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !appEventBasedModel.GoodEventFilterExpression.IsNull() {
		goodEventFilter, err := r.mapTagFilterStringToAPIModel(appEventBasedModel.GoodEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse good event filter expression: %v", err)
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	if !appEventBasedModel.IncludeInternal.IsNull() {
		includeInternal := appEventBasedModel.IncludeInternal.ValueBool()
		entity.IncludeInternal = &includeInternal
	}

	if !appEventBasedModel.IncludeSynthetic.IsNull() {
		includeSynthetic := appEventBasedModel.IncludeSynthetic.ValueBool()
		entity.IncludeSynthetic = &includeSynthetic
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapWebsiteEventBasedFromStateObject(ctx context.Context, websiteEventBasedObj types.Object) (restapi.SliEntity, error) {
	if websiteEventBasedObj.IsNull() || websiteEventBasedObj.IsUnknown() {
		return restapi.SliEntity{}, errors.New("website event based entity is null or unknown")
	}

	var websiteEventBasedModel WebsiteEventBasedModel
	diags := websiteEventBasedObj.As(ctx, &websiteEventBasedModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse website event based entity: %v", diags)
	}
	websiteID := websiteEventBasedModel.WebsiteID.ValueString()
	beaconType := websiteEventBasedModel.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       "websiteEventBased",
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !websiteEventBasedModel.BadEventFilterExpression.IsNull() {
		badEventFilter, err := r.mapTagFilterStringToAPIModel(websiteEventBasedModel.BadEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse bad event filter expression: %v", err)
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !websiteEventBasedModel.GoodEventFilterExpression.IsNull() {
		goodEventFilter, err := r.mapTagFilterStringToAPIModel(websiteEventBasedModel.GoodEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse good event filter expression: %v", err)
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapWebsiteTimeBasedFromStateObject(ctx context.Context, websiteTimeBasedObj types.Object) (restapi.SliEntity, error) {
	if websiteTimeBasedObj.IsNull() || websiteTimeBasedObj.IsUnknown() {
		return restapi.SliEntity{}, errors.New("website time based entity is null or unknown")
	}

	var websiteTimeBasedModel WebsiteTimeBasedModel
	diags := websiteTimeBasedObj.As(ctx, &websiteTimeBasedModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse website time based entity: %v", diags)
	}
	websiteID := websiteTimeBasedModel.WebsiteID.ValueString()
	beaconType := websiteTimeBasedModel.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       "websiteTimeBased",
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !websiteTimeBasedModel.FilterExpression.IsNull() {
		filterExpression, err := r.mapTagFilterStringToAPIModel(websiteTimeBasedModel.FilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse filter expression: %v", err)
		}
		entity.FilterExpression = filterExpression
	}

	return entity, nil
}
