package sliconfig

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		model.MetricConfiguration = &MetricConfigurationModel{
			MetricName:  types.StringValue(sliConfig.MetricConfiguration.Name),
			Aggregation: types.StringValue(sliConfig.MetricConfiguration.Aggregation),
			Threshold:   types.Float64Value(sliConfig.MetricConfiguration.Threshold),
		}
	} else {
		model.MetricConfiguration = nil
	}

	// Map SLI entity
	sliEntityModel := &SliEntityModel{}
	var entityDiags diag.Diagnostics

	switch sliConfig.SliEntity.Type {
	case "application":
		sliEntityModel.ApplicationTimeBased, entityDiags = r.mapApplicationTimeBasedToState(ctx, sliConfig.SliEntity)
	case "availability":
		sliEntityModel.ApplicationEventBased, entityDiags = r.mapApplicationEventBasedToState(ctx, sliConfig.SliEntity)
	case "websiteEventBased":
		sliEntityModel.WebsiteEventBased, entityDiags = r.mapWebsiteEventBasedToState(ctx, sliConfig.SliEntity)
	case "websiteTimeBased":
		sliEntityModel.WebsiteTimeBased, entityDiags = r.mapWebsiteTimeBasedToState(ctx, sliConfig.SliEntity)
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

	model.SliEntity = sliEntityModel

	// Set the state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *sliConfigResourceFramework) mapApplicationTimeBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (*ApplicationTimeBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	appTimeBasedModel := &ApplicationTimeBasedModel{
		ApplicationID: util.SetStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: util.SetStringPointerToState(sliEntity.BoundaryScope),
		ServiceID:     util.SetStringPointerToState(sliEntity.ServiceID),
		EndpointID:    util.SetStringPointerToState(sliEntity.EndpointID),
	}

	return appTimeBasedModel, diags
}

func (r *sliConfigResourceFramework) mapApplicationEventBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (*ApplicationEventBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	appEventBasedModel := &ApplicationEventBasedModel{
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
			return nil, diags
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
			return nil, diags
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

	return appEventBasedModel, diags
}

func (r *sliConfigResourceFramework) mapWebsiteEventBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (*WebsiteEventBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	websiteEventBasedModel := &WebsiteEventBasedModel{
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
			return nil, diags
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
			return nil, diags
		}
		websiteEventBasedModel.BadEventFilterExpression = util.SetStringPointerToState(badEventFilterStr)
	} else {
		websiteEventBasedModel.BadEventFilterExpression = types.StringNull()
	}

	return websiteEventBasedModel, diags
}

func (r *sliConfigResourceFramework) mapWebsiteTimeBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (*WebsiteTimeBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	websiteTimeBasedModel := &WebsiteTimeBasedModel{
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
			return nil, diags
		}
		websiteTimeBasedModel.FilterExpression = util.SetStringPointerToState(filterExprStr)
	} else {
		websiteTimeBasedModel.FilterExpression = types.StringNull()
	}

	return websiteTimeBasedModel, diags
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
	if model.MetricConfiguration != nil {
		metricConfiguration = &restapi.MetricConfiguration{
			Name:        model.MetricConfiguration.MetricName.ValueString(),
			Aggregation: model.MetricConfiguration.Aggregation.ValueString(),
			Threshold:   model.MetricConfiguration.Threshold.ValueFloat64(),
		}
	}

	// Map SLI entity
	var sliEntity restapi.SliEntity
	var entityErr error

	if model.SliEntity != nil {
		// Check which entity type is set
		if model.SliEntity.ApplicationTimeBased != nil {
			sliEntity, entityErr = r.mapApplicationTimeBasedFromState(ctx, model.SliEntity.ApplicationTimeBased)
		} else if model.SliEntity.ApplicationEventBased != nil {
			sliEntity, entityErr = r.mapApplicationEventBasedFromState(ctx, model.SliEntity.ApplicationEventBased)
		} else if model.SliEntity.WebsiteEventBased != nil {
			sliEntity, entityErr = r.mapWebsiteEventBasedFromState(ctx, model.SliEntity.WebsiteEventBased)
		} else if model.SliEntity.WebsiteTimeBased != nil {
			sliEntity, entityErr = r.mapWebsiteTimeBasedFromState(ctx, model.SliEntity.WebsiteTimeBased)
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

func (r *sliConfigResourceFramework) mapApplicationTimeBasedFromState(ctx context.Context, model *ApplicationTimeBasedModel) (restapi.SliEntity, error) {
	applicationID := model.ApplicationID.ValueString()
	boundaryScope := model.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          "application",
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if !model.ServiceID.IsNull() {
		serviceID := model.ServiceID.ValueString()
		entity.ServiceID = &serviceID
	}

	if !model.EndpointID.IsNull() {
		endpointID := model.EndpointID.ValueString()
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

func (r *sliConfigResourceFramework) mapApplicationEventBasedFromState(ctx context.Context, model *ApplicationEventBasedModel) (restapi.SliEntity, error) {
	applicationID := model.ApplicationID.ValueString()
	boundaryScope := model.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          "availability",
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if !model.BadEventFilterExpression.IsNull() {
		badEventFilter, err := r.mapTagFilterStringToAPIModel(model.BadEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse bad event filter expression: %v", err)
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !model.GoodEventFilterExpression.IsNull() {
		goodEventFilter, err := r.mapTagFilterStringToAPIModel(model.GoodEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse good event filter expression: %v", err)
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	if !model.IncludeInternal.IsNull() {
		includeInternal := model.IncludeInternal.ValueBool()
		entity.IncludeInternal = &includeInternal
	}

	if !model.IncludeSynthetic.IsNull() {
		includeSynthetic := model.IncludeSynthetic.ValueBool()
		entity.IncludeSynthetic = &includeSynthetic
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapWebsiteEventBasedFromState(ctx context.Context, model *WebsiteEventBasedModel) (restapi.SliEntity, error) {
	websiteID := model.WebsiteID.ValueString()
	beaconType := model.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       "websiteEventBased",
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !model.BadEventFilterExpression.IsNull() {
		badEventFilter, err := r.mapTagFilterStringToAPIModel(model.BadEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse bad event filter expression: %v", err)
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !model.GoodEventFilterExpression.IsNull() {
		goodEventFilter, err := r.mapTagFilterStringToAPIModel(model.GoodEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse good event filter expression: %v", err)
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapWebsiteTimeBasedFromState(ctx context.Context, model *WebsiteTimeBasedModel) (restapi.SliEntity, error) {
	websiteID := model.WebsiteID.ValueString()
	beaconType := model.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       "websiteTimeBased",
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !model.FilterExpression.IsNull() {
		filterExpression, err := r.mapTagFilterStringToAPIModel(model.FilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse filter expression: %v", err)
		}
		entity.FilterExpression = filterExpression
	}

	return entity, nil
}
