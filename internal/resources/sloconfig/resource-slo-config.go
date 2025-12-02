package sloconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

// NewSloConfigResourceHandle creates the resource handle for SLO Config
func NewSloConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.SloConfig] {
	return &sloConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:     ResourceInstanaSloConfig,
			Schema:           buildSloConfigSchema(),
			SchemaVersion:    1,
			SkipIDGeneration: true,
		},
	}
}

type sloConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *sloConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *sloConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloConfig] {
	return api.SloConfigs()
}

func (r *sloConfigResource) SetComputedFields(_ context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(plan.SetAttribute(context.Background(), path.Root(SchemaFieldID), types.StringValue(SloConfigFromTerraformIdPrefix+util.RandomID()))...)
	return diags
}

// buildSloConfigSchema constructs the complete schema for SLO configuration resource
func buildSloConfigSchema() schema.Schema {
	return schema.Schema{
		Description: SloConfigDescResource,
		Attributes: map[string]schema.Attribute{
			SchemaFieldID:               buildIDAttribute(),
			SloConfigFieldName:          buildNameAttribute(),
			SloConfigFieldTarget:        buildTargetAttribute(),
			SloConfigFieldTags:          buildTagsAttribute(),
			SloConfigFieldRbacTags:      buildRbacTagsAttribute(),
			SloConfigFieldSloEntity:     buildEntityAttribute(),
			SloConfigFieldSloIndicator:  buildIndicatorAttribute(),
			SloConfigFieldSloTimeWindow: buildTimeWindowAttribute(),
		},
	}
}

// buildIDAttribute creates the ID field schema attribute
func buildIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: SloConfigDescID,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// buildNameAttribute creates the name field schema attribute
func buildNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloConfigDescName,
	}
}

// buildTargetAttribute creates the target field schema attribute
func buildTargetAttribute() schema.Float64Attribute {
	return schema.Float64Attribute{
		Required:    true,
		Description: SloConfigDescTarget,
	}
}

// buildTagsAttribute creates the tags field schema attribute
func buildTagsAttribute() schema.ListAttribute {
	return schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Description: SloConfigDescTags,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	}
}

// buildRbacTagsAttribute creates the RBAC tags nested attribute
func buildRbacTagsAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: SloConfigDescRbacTags,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				SchemaFieldDisplayName: schema.StringAttribute{
					Required:    true,
					Description: SloConfigDescRbacTagDisplayName,
				},
				SchemaFieldID: schema.StringAttribute{
					Required:    true,
					Description: SloConfigDescRbacTagID,
				},
			},
		},
	}
}

// buildEntityAttribute creates the entity nested attribute
func buildEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigApplicationEntity:    buildApplicationEntityAttribute(),
			SloConfigWebsiteEntity:        buildWebsiteEntityAttribute(),
			SloConfigSyntheticEntity:      buildSyntheticEntityAttribute(),
			SloConfigInfrastructureEntity: buildInfrastructureEntityAttribute(),
		},
	}
}

// buildApplicationEntityAttribute creates the application entity nested attribute
func buildApplicationEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescApplicationEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldApplicationID: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescApplicationID,
			},
			SloConfigFieldBoundaryScope: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescBoundaryScope,
			},
			SloConfigFieldFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescEntityFilter,
			},
			SloConfigFieldIncludeInternal: schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: SloConfigDescIncludeInternal,
			},
			SloConfigFieldIncludeSynthetic: schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: SloConfigDescIncludeSynthetic,
			},
			SloConfigFieldServiceID: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescServiceID,
			},
			SloConfigFieldEndpointID: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescEndpointID,
			},
		},
	}
}

// buildWebsiteEntityAttribute creates the website entity nested attribute
func buildWebsiteEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescWebsiteEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldWebsiteID: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescWebsiteID,
			},
			SloConfigFieldFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescEntityFilter,
			},
			SloConfigFieldBeaconType: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescBeaconType,
			},
		},
	}
}

// buildSyntheticEntityAttribute creates the synthetic entity nested attribute
func buildSyntheticEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescSyntheticEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldSyntheticTestIDs: schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: SloConfigDescSyntheticTestIDs,
			},
			SloConfigFieldFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescEntityFilter,
			},
		},
	}
}

// buildInfrastructureEntityAttribute creates the infrastructure entity nested attribute
func buildInfrastructureEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescInfrastructureEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldInfraType: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescInfraType,
			},
			SloConfigFieldFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescEntityFilter,
			},
		},
	}
}

// buildIndicatorAttribute creates the indicator nested attribute
func buildIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescIndicator,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldTimeBasedLatency:       buildTimeBasedLatencyIndicatorAttribute(),
			SchemaFieldEventBasedLatency:      buildEventBasedLatencyIndicatorAttribute(),
			SchemaFieldTimeBasedAvailability:  buildTimeBasedAvailabilityIndicatorAttribute(),
			SchemaFieldEventBasedAvailability: buildEventBasedAvailabilityIndicatorAttribute(),
			SchemaFieldTraffic:                buildTrafficIndicatorAttribute(),
			SchemaFieldCustom:                 buildCustomIndicatorAttribute(),
			SchemaFieldSaturation:             buildSaturationIndicatorAttribute(),
		},
	}
}

// buildTimeBasedLatencyIndicatorAttribute creates the time-based latency indicator attribute
func buildTimeBasedLatencyIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescTimeBasedLatency,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldThreshold: schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: SloConfigDescThreshold,
			},
			SloConfigFieldAggregation: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: SloConfigDescAggregation,
			},
		},
	}
}

// buildEventBasedLatencyIndicatorAttribute creates the event-based latency indicator attribute
func buildEventBasedLatencyIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescEventBasedLatency,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldThreshold: schema.Float64Attribute{
				Optional:    true,
				Description: SloConfigDescThreshold,
			},
		},
	}
}

// buildTimeBasedAvailabilityIndicatorAttribute creates the time-based availability indicator attribute
func buildTimeBasedAvailabilityIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescTimeBasedAvailability,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldThreshold: schema.Float64Attribute{
				Optional:    true,
				Description: SloConfigDescThreshold,
			},
			SloConfigFieldAggregation: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescAggregation,
			},
		},
	}
}

// buildEventBasedAvailabilityIndicatorAttribute creates the event-based availability indicator attribute
func buildEventBasedAvailabilityIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescEventBasedAvailability,
		Optional:    true,
		Attributes:  map[string]schema.Attribute{},
	}
}

// buildTrafficIndicatorAttribute creates the traffic indicator attribute
func buildTrafficIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescTraffic,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldTrafficType: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescTrafficType,
			},
			SloConfigFieldThreshold: schema.Float64Attribute{
				Optional:    true,
				Description: SloConfigDescThreshold,
			},
			SchemaFieldOperator: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: SloConfigDescOperator,
				Validators: []validator.String{
					stringvalidator.OneOf(OperatorGreaterThan, OperatorGreaterThanOrEqual, OperatorLessThan, OperatorLessThanOrEqual),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// buildCustomIndicatorAttribute creates the custom indicator attribute
func buildCustomIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescCustom,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldGoodEventFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescGoodEventFilterExpression,
			},
			SloConfigFieldBadEventFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescBadEventFilterExpression,
			},
		},
	}
}

// buildSaturationIndicatorAttribute creates the saturation indicator attribute
func buildSaturationIndicatorAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescSaturation,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldMetricName: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescMetricName,
			},
			SloConfigFieldThreshold: schema.Float64Attribute{
				Optional:    true,
				Description: SloConfigDescThreshold,
			},
			SloConfigFieldAggregation: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescAggregation,
			},
			SchemaFieldOperator: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescOperator,
				Validators: []validator.String{
					stringvalidator.OneOf(OperatorGreaterThan, OperatorGreaterThanOrEqual, OperatorLessThan, OperatorLessThanOrEqual),
				},
			},
		},
	}
}

// buildTimeWindowAttribute creates the time window nested attribute
func buildTimeWindowAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescTimeWindow,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldRolling: buildRollingTimeWindowAttribute(),
			SchemaFieldFixed:   buildFixedTimeWindowAttribute(),
		},
	}
}

// buildRollingTimeWindowAttribute creates the rolling time window attribute
func buildRollingTimeWindowAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescRollingTimeWindow,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldDuration: schema.Int64Attribute{
				Optional:    true,
				Description: SloConfigDescDuration,
			},
			SloConfigFieldDurationUnit: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescDurationUnit,
			},
			SloConfigFieldTimezone: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescTimezone,
			},
		},
	}
}

// buildFixedTimeWindowAttribute creates the fixed time window attribute
func buildFixedTimeWindowAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: SloConfigDescFixedTimeWindow,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			SloConfigFieldDuration: schema.Int64Attribute{
				Optional:    true,
				Description: SloConfigDescDuration,
			},
			SloConfigFieldDurationUnit: schema.StringAttribute{
				Optional:    true,
				Description: SloConfigDescDurationUnit,
			},
			SloConfigFieldTimezone: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: SloConfigDescTimezone,
			},
			SloConfigFieldStartTimestamp: schema.Float64Attribute{
				Optional:    true,
				Description: SloConfigDescStartTimestamp,
			},
		},
	}
}

func (r *sloConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SloConfigModel

	diags.Append(r.extractModelFromPlanOrState(ctx, plan, state, &model)...)
	if diags.HasError() {
		return nil, diags
	}

	id := r.extractIDFromModel(model)

	entityData, entityDiags := r.mapEntityFromState(model.Entity)
	diags.Append(entityDiags...)

	indicator, indicatorDiags := r.mapIndicatorFromState(model.Indicator)
	diags.Append(indicatorDiags...)

	timeWindowData, timeWindowDiags := r.mapTimeWindowFromState(model.TimeWindow)
	diags.Append(timeWindowDiags...)

	if diags.HasError() {
		return nil, diags
	}
	return &restapi.SloConfig{
		ID:         id,
		Name:       model.Name.ValueString(),
		Target:     model.Target.ValueFloat64(),
		Entity:     entityData,
		Indicator:  indicator,
		TimeWindow: timeWindowData,
		Tags:       r.mapTagsFromState(model.Tags),
		RbacTags:   r.mapRbacTagsFromState(model.RbacTags),
	}, diags
}

// extractModelFromPlanOrState retrieves the model from plan or state
func (r *sloConfigResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State, model *SloConfigModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if plan != nil {
		diags.Append(plan.Get(ctx, model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, model)...)
	} else {
		diags.AddError(SloConfigErrMappingState, SloConfigErrBothPlanStateNil)
	}

	return diags
}

// extractIDFromModel extracts ID from model, returning empty string if null
func (r *sloConfigResource) extractIDFromModel(model SloConfigModel) string {
	if model.ID.IsNull() {
		return EmptyString
	}
	return model.ID.ValueString()
}

// mapTagsFromState converts tags from state to string slice
func (r *sloConfigResource) mapTagsFromState(tags []types.String) []string {
	tagsList := make([]string, 0, len(tags))
	for _, tag := range tags {
		if !tag.IsNull() && !tag.IsUnknown() {
			tagsList = append(tagsList, tag.ValueString())
		}
	}
	return tagsList
}

// mapRbacTagsFromState converts RBAC tags from state to API model
func (r *sloConfigResource) mapRbacTagsFromState(rbacTags []RbacTagModel) []restapi.RbacTag {
	rbacTagsList := make([]restapi.RbacTag, 0, len(rbacTags))
	for _, t := range rbacTags {
		rbacTagsList = append(rbacTagsList, restapi.RbacTag{
			DisplayName: t.DisplayName.ValueString(),
			ID:          t.ID.ValueString(),
		})
	}
	return rbacTagsList
}

func (r *sloConfigResource) mapEntityFromState(entityObj EntityModel) (restapi.SloEntity, diag.Diagnostics) {
	if entityObj.ApplicationEntityModel != nil {
		return r.validateAndMapApplicationEntity(entityObj.ApplicationEntityModel)
	}

	if entityObj.WebsiteEntityModel != nil {
		return r.validateAndMapWebsiteEntity(entityObj.WebsiteEntityModel)
	}

	if entityObj.SyntheticEntityModel != nil {
		return r.validateAndMapSyntheticEntity(entityObj.SyntheticEntityModel)
	}

	if entityObj.InfrastructureEntityModel != nil {
		return r.validateAndMapInfrastructureEntity(entityObj.InfrastructureEntityModel)
	}

	var diags diag.Diagnostics
	diags.AddError(SloConfigErrMissingEntity, SloConfigErrExactlyOneEntity)
	return restapi.SloEntity{}, diags
}

// validateAndMapApplicationEntity validates and maps application entity from state
func (r *sloConfigResource) validateAndMapApplicationEntity(model *ApplicationEntityModel) (restapi.SloEntity, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err := r.validateApplicationEntityFields(model); err.HasError() {
		return restapi.SloEntity{}, err
	}

	entity := r.buildApplicationEntity(model)

	filterExpr, filterDiags := r.mapFilterExpressionToEntity(model.FilterExpression)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return restapi.SloEntity{}, diags
	}

	entity.FilterExpression = filterExpr
	return entity, diags
}

// validateApplicationEntityFields validates required fields for application entity
func (r *sloConfigResource) validateApplicationEntityFields(model *ApplicationEntityModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.ApplicationID.IsUnknown() || model.ApplicationID.IsNull() ||
		model.BoundaryScope.IsNull() || model.BoundaryScope.IsUnknown() {
		diags.AddError(SloConfigErrApplicationIDRequired, SloConfigErrApplicationIDRequired)
	}

	return diags
}

// buildApplicationEntity constructs application entity from model
func (r *sloConfigResource) buildApplicationEntity(model *ApplicationEntityModel) restapi.SloEntity {
	applicationID := model.ApplicationID.ValueString()
	boundaryScope := util.SetStringPointerFromState(model.BoundaryScope)
	includeInternal := model.IncludeInternal.ValueBool()
	includeSynthetic := model.IncludeSynthetic.ValueBool()

	return restapi.SloEntity{
		Type:             SloConfigApplicationEntity,
		ApplicationID:    &applicationID,
		ServiceID:        util.SetStringPointerFromState(model.ServiceID),
		EndpointID:       util.SetStringPointerFromState(model.EndpointID),
		BoundaryScope:    boundaryScope,
		IncludeInternal:  &includeInternal,
		IncludeSynthetic: &includeSynthetic,
	}
}

// validateAndMapWebsiteEntity validates and maps website entity from state
func (r *sloConfigResource) validateAndMapWebsiteEntity(model *WebsiteEntityModel) (restapi.SloEntity, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err := r.validateWebsiteEntityFields(model); err.HasError() {
		return restapi.SloEntity{}, err
	}

	entity := r.buildWebsiteEntity(model)

	filterExpr, filterDiags := r.mapFilterExpressionToEntity(model.FilterExpression)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return restapi.SloEntity{}, diags
	}

	entity.FilterExpression = filterExpr
	return entity, diags
}

// validateWebsiteEntityFields validates required fields for website entity
func (r *sloConfigResource) validateWebsiteEntityFields(model *WebsiteEntityModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.WebsiteID.IsNull() || model.WebsiteID.IsUnknown() ||
		model.BeaconType.IsNull() || model.BeaconType.IsUnknown() {
		diags.AddError(SloConfigErrWebsiteIDRequired, SloConfigErrWebsiteIDRequired)
	}

	return diags
}

// buildWebsiteEntity constructs website entity from model
func (r *sloConfigResource) buildWebsiteEntity(model *WebsiteEntityModel) restapi.SloEntity {
	return restapi.SloEntity{
		Type:       SloConfigWebsiteEntity,
		WebsiteId:  util.SetStringPointerFromState(model.WebsiteID),
		BeaconType: util.SetStringPointerFromState(model.BeaconType),
	}
}

// validateAndMapSyntheticEntity validates and maps synthetic entity from state
func (r *sloConfigResource) validateAndMapSyntheticEntity(model *SyntheticEntityModel) (restapi.SloEntity, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err := r.validateSyntheticEntityFields(model); err.HasError() {
		return restapi.SloEntity{}, err
	}

	entity := r.buildSyntheticEntity(model)

	filterExpr, filterDiags := r.mapFilterExpressionToEntity(model.FilterExpression)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return restapi.SloEntity{}, diags
	}

	entity.FilterExpression = filterExpr
	return entity, diags
}

// validateSyntheticEntityFields validates required fields for synthetic entity
func (r *sloConfigResource) validateSyntheticEntityFields(model *SyntheticEntityModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(model.SyntheticTestIDs) == 0 {
		diags.AddError(SloConfigErrSyntheticTestIDsRequired, SloConfigErrSyntheticTestIDsRequired)
	}

	return diags
}

// buildSyntheticEntity constructs synthetic entity from model
func (r *sloConfigResource) buildSyntheticEntity(model *SyntheticEntityModel) restapi.SloEntity {
	testIDs := make([]interface{}, 0, len(model.SyntheticTestIDs))
	for _, id := range model.SyntheticTestIDs {
		if !id.IsNull() && !id.IsUnknown() {
			testIDs = append(testIDs, id.ValueString())
		}
	}

	return restapi.SloEntity{
		Type:             SloConfigSyntheticEntity,
		SyntheticTestIDs: testIDs,
	}
}

// validateAndMapInfrastructureEntity validates and maps infrastructure entity from state
func (r *sloConfigResource) validateAndMapInfrastructureEntity(model *InfrastructureEntityModel) (restapi.SloEntity, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err := r.validateInfrastructureEntityFields(model); err.HasError() {
		return restapi.SloEntity{}, err
	}

	entity := r.buildInfrastructureEntity(model)

	filterExpr, filterDiags := r.mapFilterExpressionToEntity(model.FilterExpression)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return restapi.SloEntity{}, diags
	}

	entity.FilterExpression = filterExpr
	return entity, diags
}

// validateInfrastructureEntityFields validates required fields for infrastructure entity
func (r *sloConfigResource) validateInfrastructureEntityFields(model *InfrastructureEntityModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.InfraType.IsNull() || model.InfraType.IsUnknown() {
		diags.AddError(SloConfigErrInfraTypeRequired, SloConfigErrInfraTypeRequiredMsg)
	}

	return diags
}

// buildInfrastructureEntity constructs infrastructure entity from model
func (r *sloConfigResource) buildInfrastructureEntity(model *InfrastructureEntityModel) restapi.SloEntity {
	infraType := model.InfraType.ValueString()
	return restapi.SloEntity{
		Type:      SloConfigInfrastructureEntity,
		InfraType: &infraType,
	}
}

// mapFilterExpressionToEntity converts filter expression to API model
func (r *sloConfigResource) mapFilterExpressionToEntity(filterExpression types.String) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !filterExpression.IsNull() && !filterExpression.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(filterExpression.ValueString())
		if err != nil {
			diags.AddError(SloConfigErrParsingFilterExpression, fmt.Sprintf(SloConfigErrParsingFilterExpressionMsg, err))
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		return mapper.ToAPIModel(expr), diags
	}

	return r.createDefaultTagFilter(), diags
}

// createDefaultTagFilter creates a default empty tag filter
func (r *sloConfigResource) createDefaultTagFilter() *restapi.TagFilter {
	operator := restapi.LogicalOperatorType(LogicalOperatorAnd)
	return &restapi.TagFilter{
		Type:            TagFilterTypeExpression,
		LogicalOperator: &operator,
		Elements:        []*restapi.TagFilter{},
	}
}

func (r *sloConfigResource) mapIndicatorFromState(indicatorModel IndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	if indicatorModel.TimeBasedLatencyIndicatorModel != nil {
		return r.mapTimeBasedLatencyIndicator(indicatorModel.TimeBasedLatencyIndicatorModel)
	}

	if indicatorModel.EventBasedLatencyIndicatorModel != nil {
		return r.mapEventBasedLatencyIndicator(indicatorModel.EventBasedLatencyIndicatorModel)
	}

	if indicatorModel.TimeBasedAvailabilityIndicatorModel != nil {
		return r.mapTimeBasedAvailabilityIndicator(indicatorModel.TimeBasedAvailabilityIndicatorModel)
	}

	if indicatorModel.EventBasedAvailabilityIndicatorModel != nil {
		return r.mapEventBasedAvailabilityIndicator()
	}

	if indicatorModel.TrafficIndicatorModel != nil {
		return r.mapTrafficIndicator(indicatorModel.TrafficIndicatorModel)
	}

	if indicatorModel.CustomIndicatorModel != nil {
		return r.mapCustomIndicator(indicatorModel.CustomIndicatorModel)
	}

	if indicatorModel.SaturationIndicatorModel != nil {
		return r.mapSaturationIndicator(indicatorModel.SaturationIndicatorModel)
	}

	var diags diag.Diagnostics
	diags.AddError(SloConfigErrMissingIndicator, SloConfigErrExactlyOneIndicator)
	return restapi.SloIndicator{}, diags
}

// mapTimeBasedLatencyIndicator maps time-based latency indicator from state
func (r *sloConfigResource) mapTimeBasedLatencyIndicator(model *TimeBasedLatencyIndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
		model.Aggregation.IsNull() || model.Aggregation.IsUnknown() {
		diags.AddError(SloConfigErrTimeBasedLatencyRequired, SloConfigErrTimeBasedLatencyRequired)
		return restapi.SloIndicator{}, diags
	}

	aggregation := model.Aggregation.ValueString()
	return restapi.SloIndicator{
		Blueprint:   SloConfigAPIIndicatorBlueprintLatency,
		Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
		Threshold:   model.Threshold.ValueFloat64(),
		Aggregation: &aggregation,
	}, diags
}

// mapEventBasedLatencyIndicator maps event-based latency indicator from state
func (r *sloConfigResource) mapEventBasedLatencyIndicator(model *EventBasedLatencyIndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Threshold.IsNull() || model.Threshold.IsUnknown() {
		diags.AddError(SloConfigErrEventBasedLatencyRequired, SloConfigErrEventBasedLatencyRequired)
		return restapi.SloIndicator{}, diags
	}

	defaultAgg := r.getDefaultAggregation()
	return restapi.SloIndicator{
		Blueprint:   SloConfigAPIIndicatorBlueprintLatency,
		Type:        SloConfigAPIIndicatorMeasurementTypeEventBased,
		Threshold:   model.Threshold.ValueFloat64(),
		Aggregation: defaultAgg,
	}, diags
}

// mapTimeBasedAvailabilityIndicator maps time-based availability indicator from state
func (r *sloConfigResource) mapTimeBasedAvailabilityIndicator(model *TimeBasedAvailabilityIndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
		model.Aggregation.IsNull() || model.Aggregation.IsUnknown() {
		diags.AddError(SloConfigErrTimeBasedAvailabilityRequired, SloConfigErrTimeBasedAvailabilityRequired)
		return restapi.SloIndicator{}, diags
	}

	aggregation := model.Aggregation.ValueString()
	return restapi.SloIndicator{
		Blueprint:   SloConfigAPIIndicatorBlueprintAvailability,
		Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
		Threshold:   model.Threshold.ValueFloat64(),
		Aggregation: &aggregation,
	}, diags
}

// mapEventBasedAvailabilityIndicator maps event-based availability indicator from state
func (r *sloConfigResource) mapEventBasedAvailabilityIndicator() (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics
	defaultAgg := r.getDefaultAggregation()

	return restapi.SloIndicator{
		Blueprint:   SloConfigAPIIndicatorBlueprintAvailability,
		Type:        SloConfigAPIIndicatorMeasurementTypeEventBased,
		Aggregation: defaultAgg,
	}, diags
}

// mapTrafficIndicator maps traffic indicator from state
func (r *sloConfigResource) mapTrafficIndicator(model *TrafficIndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Threshold.IsNull() || model.Threshold.IsUnknown() {
		diags.AddError(SloConfigErrTrafficRequired, SloConfigErrTrafficRequired)
		return restapi.SloIndicator{}, diags
	}

	trafficType := model.TrafficType.ValueString()
	operator := model.Operator.ValueString()
	defaultAgg := r.getDefaultAggregation()

	return restapi.SloIndicator{
		Blueprint:   SloConfigAPIIndicatorBlueprintTraffic,
		Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
		TrafficType: &trafficType,
		Threshold:   model.Threshold.ValueFloat64(),
		Operator:    &operator,
		Aggregation: defaultAgg,
	}, diags
}

// mapCustomIndicator maps custom indicator from state
func (r *sloConfigResource) mapCustomIndicator(model *CustomIndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.GoodEventFilterExpression.IsNull() || model.GoodEventFilterExpression.IsUnknown() {
		diags.AddError(SloConfigErrCustomRequired, SloConfigErrCustomRequired)
		return restapi.SloIndicator{}, diags
	}

	goodEventFilter, goodDiags := r.mapFilterExpressionToEntity(model.GoodEventFilterExpression)
	diags.Append(goodDiags...)
	if diags.HasError() {
		return restapi.SloIndicator{}, diags
	}

	badEventFilter, badDiags := r.mapFilterExpressionToEntity(model.BadEventFilterExpression)
	diags.Append(badDiags...)
	if diags.HasError() {
		return restapi.SloIndicator{}, diags
	}

	defaultAgg := r.getDefaultAggregation()
	return restapi.SloIndicator{
		Type:                      SloConfigAPIIndicatorMeasurementTypeEventBased,
		Blueprint:                 SloConfigAPIIndicatorBlueprintCustom,
		GoodEventFilterExpression: goodEventFilter,
		BadEventFilterExpression:  badEventFilter,
		Aggregation:               defaultAgg,
	}, diags
}

// mapSaturationIndicator maps saturation indicator from state
func (r *sloConfigResource) mapSaturationIndicator(model *SaturationIndicatorModel) (restapi.SloIndicator, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
		model.Operator.IsNull() || model.Operator.IsUnknown() {
		diags.AddError(SloConfigErrSaturationRequired, SloConfigErrSaturationRequiredMsg)
		return restapi.SloIndicator{}, diags
	}

	aggregation := model.Aggregation.ValueString()
	operator := model.Operator.ValueString()
	metricName := model.MetricName.ValueString()

	return restapi.SloIndicator{
		Blueprint:   SloConfigAPIIndicatorBlueprintSaturation,
		Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
		Threshold:   model.Threshold.ValueFloat64(),
		Aggregation: &aggregation,
		Operator:    &operator,
		MetricName:  &metricName,
	}, diags
}

// getDefaultAggregation returns the default aggregation value
func (r *sloConfigResource) getDefaultAggregation() *string {
	defaultAgg := DefaultAggregation
	return &defaultAgg
}

func (r *sloConfigResource) mapTimeWindowFromState(timeWindowModel TimeWindowModel) (restapi.SloTimeWindow, diag.Diagnostics) {
	if timeWindowModel.RollingTimeWindowModel != nil {
		return r.mapRollingTimeWindow(timeWindowModel.RollingTimeWindowModel)
	}

	if timeWindowModel.FixedTimeWindowModel != nil {
		return r.mapFixedTimeWindow(timeWindowModel.FixedTimeWindowModel)
	}

	var diags diag.Diagnostics
	diags.AddError(SloConfigErrMissingTimeWindow, SloConfigErrExactlyOneTimeWindow)
	return restapi.SloTimeWindow{}, diags
}

// mapRollingTimeWindow maps rolling time window from state
func (r *sloConfigResource) mapRollingTimeWindow(model *RollingTimeWindowModel) (restapi.SloTimeWindow, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err := r.validateRollingTimeWindowFields(model); err.HasError() {
		return restapi.SloTimeWindow{}, err
	}

	timeWindow := r.buildRollingTimeWindow(model)
	r.setTimezoneIfPresent(&timeWindow, model.Timezone)

	return timeWindow, diags
}

// validateRollingTimeWindowFields validates required fields for rolling time window
func (r *sloConfigResource) validateRollingTimeWindowFields(model *RollingTimeWindowModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.Duration.IsNull() || model.Duration.IsUnknown() ||
		model.DurationUnit.IsNull() || model.DurationUnit.IsUnknown() {
		diags.AddError(SloConfigErrRollingTimeWindowRequired, SloConfigErrRollingTimeWindowRequired)
	}

	return diags
}

// buildRollingTimeWindow constructs rolling time window from model
func (r *sloConfigResource) buildRollingTimeWindow(model *RollingTimeWindowModel) restapi.SloTimeWindow {
	return restapi.SloTimeWindow{
		Type:         SloConfigRollingTimeWindow,
		Duration:     int(model.Duration.ValueInt64()),
		DurationUnit: model.DurationUnit.ValueString(),
	}
}

// mapFixedTimeWindow maps fixed time window from state
func (r *sloConfigResource) mapFixedTimeWindow(model *FixedTimeWindowModel) (restapi.SloTimeWindow, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err := r.validateFixedTimeWindowFields(model); err.HasError() {
		return restapi.SloTimeWindow{}, err
	}

	timeWindow := r.buildFixedTimeWindow(model)
	r.setTimezoneIfPresent(&timeWindow, model.Timezone)

	return timeWindow, diags
}

// validateFixedTimeWindowFields validates required fields for fixed time window
func (r *sloConfigResource) validateFixedTimeWindowFields(model *FixedTimeWindowModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if model.Duration.IsNull() || model.Duration.IsUnknown() ||
		model.DurationUnit.IsNull() || model.DurationUnit.IsUnknown() ||
		model.StartTimestamp.IsNull() || model.StartTimestamp.IsUnknown() {
		diags.AddError(SloConfigErrFixedTimeWindowRequired, SloConfigErrFixedTimeWindowRequired)
	}

	return diags
}

// buildFixedTimeWindow constructs fixed time window from model
func (r *sloConfigResource) buildFixedTimeWindow(model *FixedTimeWindowModel) restapi.SloTimeWindow {
	return restapi.SloTimeWindow{
		Type:         SloConfigFixedTimeWindow,
		Duration:     int(model.Duration.ValueInt64()),
		DurationUnit: model.DurationUnit.ValueString(),
		StartTime:    model.StartTimestamp.ValueFloat64(),
	}
}

// setTimezoneIfPresent sets timezone on time window if present in model
func (r *sloConfigResource) setTimezoneIfPresent(timeWindow *restapi.SloTimeWindow, timezone types.String) {
	if !timezone.IsNull() && !timezone.IsUnknown() {
		timeWindow.Timezone = timezone.ValueString()
	}
}

func (r *sloConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SloConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model SloConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	model.ID = types.StringValue(apiObject.ID)
	model.Name = types.StringValue(apiObject.Name)
	model.Target = types.Float64Value(apiObject.Target)
	if len(apiObject.Tags) > 0 {
		model.Tags = r.mapTagsToState(apiObject.Tags)
	}
	if len(apiObject.RbacTags) > 0 {
		model.RbacTags = r.mapRbacTagsToState(apiObject.RbacTags)
	}

	entityData, entityDiags := r.mapEntityToState(apiObject)
	diags.Append(entityDiags...)
	if !diags.HasError() {
		model.Entity = entityData
	}

	indicatorData, indicatorDiags := r.mapIndicatorToState(apiObject)
	diags.Append(indicatorDiags...)
	if !diags.HasError() {
		model.Indicator = indicatorData
	}

	timeWindowData, timeWindowDiags := r.mapTimeWindowToState(apiObject)
	diags.Append(timeWindowDiags...)
	if !diags.HasError() {
		model.TimeWindow = timeWindowData
	}

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapTagsToState converts tags from API to state
func (r *sloConfigResource) mapTagsToState(tags []string) []types.String {
	if tags == nil {
		return nil
	}

	stateTags := make([]types.String, 0, len(tags))
	for _, tag := range tags {
		stateTags = append(stateTags, types.StringValue(tag))
	}
	return stateTags
}

// mapRbacTagsToState converts RBAC tags from API to state
func (r *sloConfigResource) mapRbacTagsToState(rbacTags []restapi.RbacTag) []RbacTagModel {
	if rbacTags == nil {
		return nil
	}

	stateRbacTags := make([]RbacTagModel, 0, len(rbacTags))
	for _, tag := range rbacTags {
		stateRbacTags = append(stateRbacTags, RbacTagModel{
			DisplayName: types.StringValue(tag.DisplayName),
			ID:          types.StringValue(tag.ID),
		})
	}
	return stateRbacTags
}

func (r *sloConfigResource) mapEntityToState(apiObject *restapi.SloConfig) (EntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	entityModel := EntityModel{}

	switch apiObject.Entity.Type {
	case SloConfigApplicationEntity:
		appModel, appDiags := r.mapApplicationEntityToState(apiObject.Entity)
		diags.Append(appDiags...)
		if !diags.HasError() {
			entityModel.ApplicationEntityModel = &appModel
		}
	case SloConfigWebsiteEntity:
		websiteModel, websiteDiags := r.mapWebsiteEntityToState(apiObject.Entity)
		diags.Append(websiteDiags...)
		if !diags.HasError() {
			entityModel.WebsiteEntityModel = &websiteModel
		}
	case SloConfigSyntheticEntity:
		syntheticModel, syntheticDiags := r.mapSyntheticEntityToState(apiObject.Entity)
		diags.Append(syntheticDiags...)
		if !diags.HasError() {
			entityModel.SyntheticEntityModel = &syntheticModel
		}
	case SloConfigInfrastructureEntity:
		infraModel, infraDiags := r.mapInfrastructureEntityToState(apiObject.Entity)
		diags.Append(infraDiags...)
		if !diags.HasError() {
			entityModel.InfrastructureEntityModel = &infraModel
		}
	default:
		diags.AddError(SloConfigErrMappingEntityToState, fmt.Sprintf(SloConfigErrUnsupportedEntityType, apiObject.Entity.Type))
	}

	return entityModel, diags
}

// mapApplicationEntityToState converts application entity from API to state
func (r *sloConfigResource) mapApplicationEntityToState(entity restapi.SloEntity) (ApplicationEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := ApplicationEntityModel{
		ApplicationID:    util.SetStringPointerToState(entity.ApplicationID),
		BoundaryScope:    util.SetStringPointerToState(entity.BoundaryScope),
		IncludeInternal:  util.SetBoolPointerToState(entity.IncludeInternal),
		IncludeSynthetic: util.SetBoolPointerToState(entity.IncludeSynthetic),
		ServiceID:        util.SetStringPointerToState(entity.ServiceID),
		EndpointID:       util.SetStringPointerToState(entity.EndpointID),
	}

	filterExpr, filterDiags := r.mapFilterExpressionToState(entity.FilterExpression)
	diags.Append(filterDiags...)
	if !diags.HasError() {
		model.FilterExpression = filterExpr
	}

	return model, diags
}

// mapWebsiteEntityToState converts website entity from API to state
func (r *sloConfigResource) mapWebsiteEntityToState(entity restapi.SloEntity) (WebsiteEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := WebsiteEntityModel{
		WebsiteID:  util.SetStringPointerToState(entity.WebsiteId),
		BeaconType: util.SetStringPointerToState(entity.BeaconType),
	}

	filterExpr, filterDiags := r.mapFilterExpressionToState(entity.FilterExpression)
	diags.Append(filterDiags...)
	if !diags.HasError() {
		model.FilterExpression = filterExpr
	}

	return model, diags
}

// mapSyntheticEntityToState converts synthetic entity from API to state
func (r *sloConfigResource) mapSyntheticEntityToState(entity restapi.SloEntity) (SyntheticEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	testIDs := make([]types.String, 0, len(entity.SyntheticTestIDs))
	for _, id := range entity.SyntheticTestIDs {
		if idStr, ok := id.(string); ok {
			testIDs = append(testIDs, types.StringValue(idStr))
		}
	}

	model := SyntheticEntityModel{
		SyntheticTestIDs: testIDs,
	}

	filterExpr, filterDiags := r.mapFilterExpressionToState(entity.FilterExpression)
	diags.Append(filterDiags...)
	if !diags.HasError() {
		model.FilterExpression = filterExpr
	}

	return model, diags
}

// mapInfrastructureEntityToState converts infrastructure entity from API to state
func (r *sloConfigResource) mapInfrastructureEntityToState(entity restapi.SloEntity) (InfrastructureEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := InfrastructureEntityModel{
		InfraType: util.SetStringPointerToState(entity.InfraType),
	}

	filterExpr, filterDiags := r.mapFilterExpressionToState(entity.FilterExpression)
	diags.Append(filterDiags...)
	if !diags.HasError() {
		model.FilterExpression = filterExpr
	}

	return model, diags
}

// mapFilterExpressionToState converts filter expression from API to state
func (r *sloConfigResource) mapFilterExpressionToState(filterExpression *restapi.TagFilter) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	if filterExpression == nil {
		return types.StringNull(), diags
	}

	filterExprStr, err := tagfilter.MapTagFilterToNormalizedString(filterExpression)
	if err != nil {
		diags.AddError(SloConfigErrNormalizingFilterExpression, fmt.Sprintf(SloConfigErrNormalizingFilterExpressionMsg, err))
		return types.StringNull(), diags
	}

	return util.SetStringPointerToState(filterExprStr), diags
}

func (r *sloConfigResource) mapIndicatorToState(apiObject *restapi.SloConfig) (IndicatorModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	indicator := apiObject.Indicator

	indicatorModel := IndicatorModel{}

	switch {
	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintLatency:
		indicatorModel.TimeBasedLatencyIndicatorModel = r.createTimeBasedLatencyModel(indicator)

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintLatency:
		indicatorModel.EventBasedLatencyIndicatorModel = r.createEventBasedLatencyModel(indicator)

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintAvailability:
		indicatorModel.TimeBasedAvailabilityIndicatorModel = r.createTimeBasedAvailabilityModel(indicator)

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintAvailability:
		indicatorModel.EventBasedAvailabilityIndicatorModel = r.createEventBasedAvailabilityModel()

	case indicator.Blueprint == SloConfigAPIIndicatorBlueprintTraffic:
		indicatorModel.TrafficIndicatorModel = r.createTrafficModel(indicator)

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintCustom:
		customModel, customDiags := r.createCustomModel(indicator)
		diags.Append(customDiags...)
		if !diags.HasError() {
			indicatorModel.CustomIndicatorModel = customModel
		}

	case indicator.Type == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator.Blueprint == SloConfigAPIIndicatorBlueprintSaturation:
		indicatorModel.SaturationIndicatorModel = r.createSaturationModel(indicator)

	default:
		diags.AddError(SloConfigErrMappingIndicatorToState, fmt.Sprintf(SloConfigErrUnsupportedIndicatorType, indicator.Type, indicator.Blueprint))
	}

	return indicatorModel, diags
}

// createTimeBasedLatencyModel creates time-based latency indicator model
func (r *sloConfigResource) createTimeBasedLatencyModel(indicator restapi.SloIndicator) *TimeBasedLatencyIndicatorModel {
	return &TimeBasedLatencyIndicatorModel{
		Threshold:   types.Float64Value(indicator.Threshold),
		Aggregation: util.SetStringPointerToState(indicator.Aggregation),
	}
}

// createEventBasedLatencyModel creates event-based latency indicator model
func (r *sloConfigResource) createEventBasedLatencyModel(indicator restapi.SloIndicator) *EventBasedLatencyIndicatorModel {
	return &EventBasedLatencyIndicatorModel{
		Threshold: types.Float64Value(indicator.Threshold),
	}
}

// createTimeBasedAvailabilityModel creates time-based availability indicator model
func (r *sloConfigResource) createTimeBasedAvailabilityModel(indicator restapi.SloIndicator) *TimeBasedAvailabilityIndicatorModel {
	return &TimeBasedAvailabilityIndicatorModel{
		Threshold:   types.Float64Value(indicator.Threshold),
		Aggregation: util.SetStringPointerToState(indicator.Aggregation),
	}
}

// createEventBasedAvailabilityModel creates event-based availability indicator model
func (r *sloConfigResource) createEventBasedAvailabilityModel() *EventBasedAvailabilityIndicatorModel {
	return &EventBasedAvailabilityIndicatorModel{}
}

// createTrafficModel creates traffic indicator model
func (r *sloConfigResource) createTrafficModel(indicator restapi.SloIndicator) *TrafficIndicatorModel {
	return &TrafficIndicatorModel{
		TrafficType: util.SetStringPointerToState(indicator.TrafficType),
		Threshold:   types.Float64Value(indicator.Threshold),
		Operator:    util.SetStringPointerToState(indicator.Operator),
	}
}

// createCustomModel creates custom indicator model
func (r *sloConfigResource) createCustomModel(indicator restapi.SloIndicator) (*CustomIndicatorModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	model := &CustomIndicatorModel{}

	goodFilterExpr, goodDiags := r.mapFilterExpressionToState(indicator.GoodEventFilterExpression)
	diags.Append(goodDiags...)
	if !diags.HasError() {
		model.GoodEventFilterExpression = goodFilterExpr
	}

	badFilterExpr, badDiags := r.mapFilterExpressionToState(indicator.BadEventFilterExpression)
	diags.Append(badDiags...)
	if !diags.HasError() {
		model.BadEventFilterExpression = badFilterExpr
	}

	return model, diags
}

// createSaturationModel creates saturation indicator model
func (r *sloConfigResource) createSaturationModel(indicator restapi.SloIndicator) *SaturationIndicatorModel {
	return &SaturationIndicatorModel{
		MetricName:  util.SetStringPointerToState(indicator.MetricName),
		Threshold:   types.Float64Value(indicator.Threshold),
		Aggregation: util.SetStringPointerToState(indicator.Aggregation),
		Operator:    util.SetStringPointerToState(indicator.Operator),
	}
}

func (r *sloConfigResource) mapTimeWindowToState(apiObject *restapi.SloConfig) (TimeWindowModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	timeWindow := apiObject.TimeWindow

	timeWindowModel := TimeWindowModel{}

	switch timeWindow.Type {
	case SloConfigRollingTimeWindow:
		timeWindowModel.RollingTimeWindowModel = r.createRollingTimeWindowModel(timeWindow)

	case SloConfigFixedTimeWindow:
		timeWindowModel.FixedTimeWindowModel = r.createFixedTimeWindowModel(timeWindow)

	default:
		diags.AddError(SloConfigErrMappingTimeWindowToState, fmt.Sprintf(SloConfigErrUnsupportedTimeWindowType, timeWindow.Type))
	}

	return timeWindowModel, diags
}

// createRollingTimeWindowModel creates rolling time window model
func (r *sloConfigResource) createRollingTimeWindowModel(timeWindow restapi.SloTimeWindow) *RollingTimeWindowModel {
	return &RollingTimeWindowModel{
		Duration:     types.Int64Value(int64(timeWindow.Duration)),
		DurationUnit: types.StringValue(timeWindow.DurationUnit),
		Timezone:     r.setTimezoneToState(timeWindow.Timezone),
	}
}

// createFixedTimeWindowModel creates fixed time window model
func (r *sloConfigResource) createFixedTimeWindowModel(timeWindow restapi.SloTimeWindow) *FixedTimeWindowModel {
	return &FixedTimeWindowModel{
		Duration:       types.Int64Value(int64(timeWindow.Duration)),
		DurationUnit:   types.StringValue(timeWindow.DurationUnit),
		StartTimestamp: types.Float64Value(timeWindow.StartTime),
		Timezone:       r.setTimezoneToState(timeWindow.Timezone),
	}
}

// setTimezoneToState converts timezone string to state value
func (r *sloConfigResource) setTimezoneToState(timezone string) types.String {
	if timezone != EmptyString {
		return types.StringValue(timezone)
	}
	return types.StringNull()
}

// Made with Bob
