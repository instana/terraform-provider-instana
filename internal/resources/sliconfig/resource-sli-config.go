package sliconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewSliConfigResourceHandle creates the resource handle for SLI configuration
func NewSliConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.SliConfig] {
	return &sliConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSliConfig,
			Schema:        buildSliConfigSchema(),
			SchemaVersion: 1,
			CreateOnly:    true,
		},
	}
}

// buildSliConfigSchema constructs the complete schema for SLI configuration resource
func buildSliConfigSchema() schema.Schema {
	return schema.Schema{
		Description: SliConfigDescResource,
		Attributes: map[string]schema.Attribute{
			SchemaFieldID:                         buildIDAttribute(),
			SchemaFieldName:                       buildNameAttribute(),
			SchemaFieldInitialEvaluationTimestamp: buildInitialEvaluationTimestampAttribute(),
			SchemaFieldMetricConfiguration:        buildMetricConfigurationAttribute(),
			SchemaFieldSliEntity:                  buildSliEntityAttribute(),
		},
	}
}

// buildIDAttribute creates the ID field schema attribute
func buildIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: SliConfigDescID,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// buildNameAttribute creates the name field schema attribute
func buildNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescName,
		Validators: []validator.String{
			stringvalidator.LengthBetween(NameMinLength, NameMaxLength),
		},
	}
}

// buildInitialEvaluationTimestampAttribute creates the initial evaluation timestamp field schema attribute
func buildInitialEvaluationTimestampAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: SliConfigDescInitialEvaluationTimestamp,
	}
}

// buildMetricConfigurationAttribute creates the metric configuration nested attribute
func buildMetricConfigurationAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SliConfigDescMetricConfiguration,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldMetricName:  buildMetricNameAttribute(),
			SchemaFieldAggregation: buildAggregationAttribute(),
			SchemaFieldThreshold:   buildThresholdAttribute(),
		},
	}
}

// buildMetricNameAttribute creates the metric name field schema attribute
func buildMetricNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescMetricName,
	}
}

// buildAggregationAttribute creates the aggregation field schema attribute with all valid values
func buildAggregationAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescAggregation,
		Validators: []validator.String{
			stringvalidator.OneOf(
				AggregationTypeSum,
				AggregationTypeMean,
				AggregationTypeMax,
				AggregationTypeMin,
				AggregationTypeP25,
				AggregationTypeP50,
				AggregationTypeP75,
				AggregationTypeP90,
				AggregationTypeP95,
				AggregationTypeP98,
				AggregationTypeP99,
				AggregationTypeP99_9,
				AggregationTypeP99_99,
				AggregationTypeDistribution,
				AggregationTypeDistinctCount,
				AggregationTypeSumPositive,
				AggregationTypePerSecond,
			),
		},
	}
}

// buildThresholdAttribute creates the threshold field schema attribute
func buildThresholdAttribute() schema.Float64Attribute {
	return schema.Float64Attribute{
		Required:    true,
		Description: SliConfigDescThreshold,
		Validators: []validator.Float64{
			float64validator.AtLeast(ThresholdMinValue),
		},
	}
}

// buildSliEntityAttribute creates the SLI entity nested attribute
func buildSliEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SliConfigDescSliEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldApplicationTimeBased:  buildApplicationTimeBasedAttribute(),
			SchemaFieldApplicationEventBased: buildApplicationEventBasedAttribute(),
			SchemaFieldWebsiteEventBased:     buildWebsiteEventBasedAttribute(),
			SchemaFieldWebsiteTimeBased:      buildWebsiteTimeBasedAttribute(),
		},
	}
}

// buildApplicationTimeBasedAttribute creates the application time-based entity nested attribute
func buildApplicationTimeBasedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SliConfigDescApplicationTimeBased,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldApplicationID: buildApplicationIDAttribute(),
			SchemaFieldServiceID:     buildServiceIDAttribute(),
			SchemaFieldEndpointID:    buildEndpointIDAttribute(),
			SchemaFieldBoundaryScope: buildBoundaryScopeAttribute(),
		},
	}
}

// buildApplicationEventBasedAttribute creates the application event-based entity nested attribute
func buildApplicationEventBasedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SliConfigDescApplicationEventBased,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldApplicationID:             buildApplicationIDAttribute(),
			SchemaFieldBoundaryScope:             buildBoundaryScopeAttribute(),
			SchemaFieldBadEventFilterExpression:  buildBadEventFilterExpressionAttribute(),
			SchemaFieldGoodEventFilterExpression: buildGoodEventFilterExpressionAttribute(),
			SchemaFieldIncludeInternal:           buildIncludeInternalAttribute(),
			SchemaFieldIncludeSynthetic:          buildIncludeSyntheticAttribute(),
			SchemaFieldEndpointID:                buildEndpointIDAvailabilityAttribute(),
			SchemaFieldServiceID:                 buildServiceIDAvailabilityAttribute(),
		},
	}
}

// buildWebsiteEventBasedAttribute creates the website event-based entity nested attribute
func buildWebsiteEventBasedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SliConfigDescWebsiteEventBased,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldWebsiteID:                 buildWebsiteIDAttribute(),
			SchemaFieldBadEventFilterExpression:  buildBadEventFilterExpressionAttribute(),
			SchemaFieldGoodEventFilterExpression: buildGoodEventFilterExpressionAttribute(),
			SchemaFieldBeaconType:                buildBeaconTypeAttribute(),
		},
	}
}

// buildWebsiteTimeBasedAttribute creates the website time-based entity nested attribute
func buildWebsiteTimeBasedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SliConfigDescWebsiteTimeBased,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldWebsiteID:        buildWebsiteIDAttribute(),
			SchemaFieldFilterExpression: buildFilterExpressionAttribute(),
			SchemaFieldBeaconType:       buildBeaconTypeAttribute(),
		},
	}
}

// buildApplicationIDAttribute creates the application ID field schema attribute
func buildApplicationIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescApplicationID,
	}
}

// buildServiceIDAttribute creates the service ID field schema attribute
func buildServiceIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SliConfigDescServiceID,
	}
}

// buildEndpointIDAttribute creates the endpoint ID field schema attribute
func buildEndpointIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SliConfigDescEndpointID,
	}
}

// buildEndpointIDAvailabilityAttribute creates the endpoint ID field schema attribute for availability context
func buildEndpointIDAvailabilityAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SliConfigDescEndpointIDAvailability,
	}
}

// buildServiceIDAvailabilityAttribute creates the service ID field schema attribute for availability context
func buildServiceIDAvailabilityAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SliConfigDescServiceIDAvailability,
	}
}

// buildBoundaryScopeAttribute creates the boundary scope field schema attribute
func buildBoundaryScopeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescBoundaryScope,
		Validators: []validator.String{
			stringvalidator.OneOf(BoundaryScopeAll, BoundaryScopeInbound),
		},
	}
}

// buildBadEventFilterExpressionAttribute creates the bad event filter expression field schema attribute
func buildBadEventFilterExpressionAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescBadEventFilterExpression,
	}
}

// buildGoodEventFilterExpressionAttribute creates the good event filter expression field schema attribute
func buildGoodEventFilterExpressionAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescGoodEventFilterExpression,
	}
}

// buildIncludeInternalAttribute creates the include internal field schema attribute
func buildIncludeInternalAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:    true,
		Description: SliConfigDescIncludeInternal,
	}
}

// buildIncludeSyntheticAttribute creates the include synthetic field schema attribute
func buildIncludeSyntheticAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:    true,
		Description: SliConfigDescIncludeSynthetic,
	}
}

// buildWebsiteIDAttribute creates the website ID field schema attribute
func buildWebsiteIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescWebsiteID,
	}
}

// buildBeaconTypeAttribute creates the beacon type field schema attribute
func buildBeaconTypeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SliConfigDescBeaconType,
		Validators: []validator.String{
			stringvalidator.OneOf(
				BeaconTypePageLoad,
				BeaconTypeResourceLoad,
				BeaconTypeHttpRequest,
				BeaconTypeError,
				BeaconTypeCustom,
				BeaconTypePageChange,
			),
		},
	}
}

// buildFilterExpressionAttribute creates the filter expression field schema attribute
func buildFilterExpressionAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SliConfigDescFilterExpression,
	}
}

func (r *sliConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *sliConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SliConfig] {
	return api.SliConfigs()
}

func (r *sliConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sliConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, sliConfig *restapi.SliConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model SliConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	model.ID = types.StringValue(sliConfig.ID)
	model.Name = types.StringValue(sliConfig.Name)
	if model.InitialEvaluationTimestamp.IsNull() || model.InitialEvaluationTimestamp.IsUnknown() {
		model.InitialEvaluationTimestamp = types.Int64Value(int64(sliConfig.InitialEvaluationTimestamp))
	}

	if sliConfig.MetricConfiguration != nil {
		model.MetricConfiguration = r.mapMetricConfigurationToState(sliConfig.MetricConfiguration)
	}

	sliEntityModel, entityDiags := r.mapSliEntityToState(sliConfig.SliEntity)
	diags.Append(entityDiags...)
	if diags.HasError() {
		return diags
	}

	model.SliEntity = sliEntityModel
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapMetricConfigurationToState converts API metric configuration to state model
func (r *sliConfigResource) mapMetricConfigurationToState(metricConfig *restapi.MetricConfiguration) *MetricConfigurationModel {
	return &MetricConfigurationModel{
		MetricName:  types.StringValue(metricConfig.Name),
		Aggregation: types.StringValue(metricConfig.Aggregation),
		Threshold:   types.Float64Value(metricConfig.Threshold),
	}
}

// mapSliEntityToState routes SLI entity mapping based on entity type
func (r *sliConfigResource) mapSliEntityToState(sliEntity restapi.SliEntity) (*SliEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	sliEntityModel := &SliEntityModel{}

	switch sliEntity.Type {
	case SliEntityTypeApplication:
		sliEntityModel.ApplicationTimeBased = r.mapApplicationTimeBasedToState(sliEntity)
	case SliEntityTypeAvailability:
		model, err := r.mapApplicationEventBasedToState(sliEntity)
		if err.HasError() {
			return nil, err
		}
		sliEntityModel.ApplicationEventBased = model
	case SliEntityTypeWebsiteEventBased:
		model, err := r.mapWebsiteEventBasedToState(sliEntity)
		if err.HasError() {
			return nil, err
		}
		sliEntityModel.WebsiteEventBased = model
	case SliEntityTypeWebsiteTimeBased:
		model, err := r.mapWebsiteTimeBasedToState(sliEntity)
		if err.HasError() {
			return nil, err
		}
		sliEntityModel.WebsiteTimeBased = model
	default:
		diags.AddError(
			SliConfigErrUnsupportedEntityType,
			fmt.Sprintf(SliConfigErrUnsupportedEntityTypeMsg, sliEntity.Type),
		)
		return nil, diags
	}

	return sliEntityModel, diags
}

// mapApplicationTimeBasedToState converts application time-based entity to state model
func (r *sliConfigResource) mapApplicationTimeBasedToState(sliEntity restapi.SliEntity) *ApplicationTimeBasedModel {
	return &ApplicationTimeBasedModel{
		ApplicationID: util.SetStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: util.SetStringPointerToState(sliEntity.BoundaryScope),
		ServiceID:     util.SetStringPointerToState(sliEntity.ServiceID),
		EndpointID:    util.SetStringPointerToState(sliEntity.EndpointID),
	}
}

// mapApplicationEventBasedToState converts application event-based entity to state model
func (r *sliConfigResource) mapApplicationEventBasedToState(sliEntity restapi.SliEntity) (*ApplicationEventBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &ApplicationEventBasedModel{
		ApplicationID: util.SetStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: util.SetStringPointerToState(sliEntity.BoundaryScope),
		EndpointID:    util.SetStringPointerToState(sliEntity.EndpointID),
		ServiceID:     util.SetStringPointerToState(sliEntity.ServiceID),
	}

	goodFilterStr, err := r.mapTagFilterToString(sliEntity.GoodEventFilterExpression, ErrMsgGoodEventFilterContext)
	if err != nil {
		diags.AddError(SliConfigErrMappingGoodEventFilter, err.Error())
		return nil, diags
	}
	model.GoodEventFilterExpression = goodFilterStr

	badFilterStr, err := r.mapTagFilterToString(sliEntity.BadEventFilterExpression, ErrMsgBadEventFilterContext)
	if err != nil {
		diags.AddError(SliConfigErrMappingBadEventFilter, err.Error())
		return nil, diags
	}
	model.BadEventFilterExpression = badFilterStr

	model.IncludeInternal = r.mapBooleanPointerToState(sliEntity.IncludeInternal, DefaultIncludeInternalValue)
	model.IncludeSynthetic = r.mapBooleanPointerToState(sliEntity.IncludeSynthetic, DefaultIncludeSyntheticValue)

	return model, diags
}

// mapWebsiteEventBasedToState converts website event-based entity to state model
func (r *sliConfigResource) mapWebsiteEventBasedToState(sliEntity restapi.SliEntity) (*WebsiteEventBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &WebsiteEventBasedModel{
		WebsiteID:  util.SetStringPointerToState(sliEntity.WebsiteId),
		BeaconType: util.SetStringPointerToState(sliEntity.BeaconType),
	}

	goodFilterStr, err := r.mapTagFilterToString(sliEntity.GoodEventFilterExpression, ErrMsgGoodEventFilterContext)
	if err != nil {
		diags.AddError(SliConfigErrMappingGoodEventFilter, err.Error())
		return nil, diags
	}
	model.GoodEventFilterExpression = goodFilterStr

	badFilterStr, err := r.mapTagFilterToString(sliEntity.BadEventFilterExpression, ErrMsgBadEventFilterContext)
	if err != nil {
		diags.AddError(SliConfigErrMappingBadEventFilter, err.Error())
		return nil, diags
	}
	model.BadEventFilterExpression = badFilterStr

	return model, diags
}

// mapWebsiteTimeBasedToState converts website time-based entity to state model
func (r *sliConfigResource) mapWebsiteTimeBasedToState(sliEntity restapi.SliEntity) (*WebsiteTimeBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &WebsiteTimeBasedModel{
		WebsiteID:  util.SetStringPointerToState(sliEntity.WebsiteId),
		BeaconType: util.SetStringPointerToState(sliEntity.BeaconType),
	}

	filterStr, err := r.mapTagFilterToString(sliEntity.FilterExpression, ErrMsgFilterExpressionContext)
	if err != nil {
		diags.AddError(SliConfigErrMappingFilterExpression, err.Error())
		return nil, diags
	}
	model.FilterExpression = filterStr

	return model, diags
}

// mapTagFilterToString converts tag filter to normalized string representation
func (r *sliConfigResource) mapTagFilterToString(tagFilter *restapi.TagFilter, context string) (types.String, error) {
	if tagFilter == nil {
		return types.StringNull(), nil
	}

	filterStr, err := tagfilter.MapTagFilterToNormalizedString(tagFilter)
	if err != nil {
		return types.StringNull(), fmt.Errorf(ErrMsgFailedToParseFilterExpression, context, err)
	}

	return util.SetStringPointerToState(filterStr), nil
}

// mapBooleanPointerToState converts boolean pointer to state value with default fallback
func (r *sliConfigResource) mapBooleanPointerToState(value *bool, defaultValue bool) types.Bool {
	if value != nil {
		return types.BoolValue(*value)
	}
	return types.BoolValue(defaultValue)
}

func (r *sliConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SliConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SliConfigModel

	diags.Append(r.extractModelFromPlanOrState(ctx, plan, state, &model)...)
	if diags.HasError() {
		return nil, diags
	}

	sliConfig := &restapi.SliConfig{
		ID:                         r.extractIDFromModel(model),
		Name:                       model.Name.ValueString(),
		InitialEvaluationTimestamp: r.extractInitialEvaluationTimestamp(model),
		MetricConfiguration:        r.mapMetricConfigurationFromState(model.MetricConfiguration),
	}

	sliEntity, err := r.mapSliEntityFromState(model.SliEntity)
	if err != nil {
		diags.AddError(SliConfigErrMappingSliEntity, fmt.Sprintf(SliConfigErrMappingSliEntityMsg, err))
		return nil, diags
	}
	sliConfig.SliEntity = sliEntity

	return sliConfig, diags
}

// extractModelFromPlanOrState retrieves the model from plan or state
func (r *sliConfigResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State, model *SliConfigModel) diag.Diagnostics {
	if plan != nil {
		return plan.Get(ctx, model)
	}
	if state != nil {
		return state.Get(ctx, model)
	}
	return nil
}

// extractIDFromModel extracts ID from model, returning empty string if null
func (r *sliConfigResource) extractIDFromModel(model SliConfigModel) string {
	if model.ID.IsNull() {
		return ""
	}
	return model.ID.ValueString()
}

// extractInitialEvaluationTimestamp extracts initial evaluation timestamp from model
func (r *sliConfigResource) extractInitialEvaluationTimestamp(model SliConfigModel) int {
	if model.InitialEvaluationTimestamp.IsNull() {
		return DefaultInitialEvaluationTimestamp
	}
	return int(model.InitialEvaluationTimestamp.ValueInt64())
}

// mapMetricConfigurationFromState converts state metric configuration to API model
func (r *sliConfigResource) mapMetricConfigurationFromState(metricConfig *MetricConfigurationModel) *restapi.MetricConfiguration {
	if metricConfig == nil {
		return nil
	}

	return &restapi.MetricConfiguration{
		Name:        metricConfig.MetricName.ValueString(),
		Aggregation: metricConfig.Aggregation.ValueString(),
		Threshold:   metricConfig.Threshold.ValueFloat64(),
	}
}

// mapSliEntityFromState routes SLI entity mapping from state based on entity type
func (r *sliConfigResource) mapSliEntityFromState(sliEntityModel *SliEntityModel) (restapi.SliEntity, error) {
	if sliEntityModel == nil {
		return restapi.SliEntity{}, nil
	}

	if sliEntityModel.ApplicationTimeBased != nil {
		return r.mapApplicationTimeBasedFromState(sliEntityModel.ApplicationTimeBased)
	}

	if sliEntityModel.ApplicationEventBased != nil {
		return r.mapApplicationEventBasedFromState(sliEntityModel.ApplicationEventBased)
	}

	if sliEntityModel.WebsiteEventBased != nil {
		return r.mapWebsiteEventBasedFromState(sliEntityModel.WebsiteEventBased)
	}

	if sliEntityModel.WebsiteTimeBased != nil {
		return r.mapWebsiteTimeBasedFromState(sliEntityModel.WebsiteTimeBased)
	}

	return restapi.SliEntity{}, nil
}

// mapApplicationTimeBasedFromState converts application time-based state model to API entity
func (r *sliConfigResource) mapApplicationTimeBasedFromState(model *ApplicationTimeBasedModel) (restapi.SliEntity, error) {
	applicationID := model.ApplicationID.ValueString()
	boundaryScope := model.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          SliEntityTypeApplication,
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

// mapApplicationEventBasedFromState converts application event-based state model to API entity
func (r *sliConfigResource) mapApplicationEventBasedFromState(model *ApplicationEventBasedModel) (restapi.SliEntity, error) {
	applicationID := model.ApplicationID.ValueString()
	boundaryScope := model.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          SliEntityTypeAvailability,
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if err := r.mapEventFilterExpressions(&entity, model.BadEventFilterExpression, model.GoodEventFilterExpression); err != nil {
		return restapi.SliEntity{}, err
	}

	r.mapOptionalBooleanFields(&entity, model.IncludeInternal, model.IncludeSynthetic)

	return entity, nil
}

// mapWebsiteEventBasedFromState converts website event-based state model to API entity
func (r *sliConfigResource) mapWebsiteEventBasedFromState(model *WebsiteEventBasedModel) (restapi.SliEntity, error) {
	websiteID := model.WebsiteID.ValueString()
	beaconType := model.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       SliEntityTypeWebsiteEventBased,
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if err := r.mapEventFilterExpressions(&entity, model.BadEventFilterExpression, model.GoodEventFilterExpression); err != nil {
		return restapi.SliEntity{}, err
	}

	return entity, nil
}

// mapWebsiteTimeBasedFromState converts website time-based state model to API entity
func (r *sliConfigResource) mapWebsiteTimeBasedFromState(model *WebsiteTimeBasedModel) (restapi.SliEntity, error) {
	websiteID := model.WebsiteID.ValueString()
	beaconType := model.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       SliEntityTypeWebsiteTimeBased,
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !model.FilterExpression.IsNull() {
		filterExpression, err := r.parseTagFilterString(model.FilterExpression.ValueString(), ErrMsgFilterExpressionContext)
		if err != nil {
			return restapi.SliEntity{}, err
		}
		entity.FilterExpression = filterExpression
	}

	return entity, nil
}

// mapEventFilterExpressions maps bad and good event filter expressions to entity
func (r *sliConfigResource) mapEventFilterExpressions(entity *restapi.SliEntity, badFilter, goodFilter types.String) error {
	if !badFilter.IsNull() {
		badEventFilter, err := r.parseTagFilterString(badFilter.ValueString(), ErrMsgBadEventFilterContext)
		if err != nil {
			return err
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !goodFilter.IsNull() {
		goodEventFilter, err := r.parseTagFilterString(goodFilter.ValueString(), ErrMsgGoodEventFilterContext)
		if err != nil {
			return err
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	return nil
}

// mapOptionalBooleanFields maps optional boolean fields to entity
func (r *sliConfigResource) mapOptionalBooleanFields(entity *restapi.SliEntity, includeInternal, includeSynthetic types.Bool) {
	if !includeInternal.IsNull() {
		value := includeInternal.ValueBool()
		entity.IncludeInternal = &value
	}

	if !includeSynthetic.IsNull() {
		value := includeSynthetic.ValueBool()
		entity.IncludeSynthetic = &value
	}
}

// parseTagFilterString parses a tag filter string into API model
func (r *sliConfigResource) parseTagFilterString(input, context string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf(ErrMsgFailedToParseFilterExpression, context, err)
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

// Made with Bob
