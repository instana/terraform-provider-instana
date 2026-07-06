package apdexconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/instana/instana-go-client/api"
	client "github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	tag "github.com/instana/instana-go-client/shared/tagfilter"
	common "github.com/instana/instana-go-client/shared/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewApdexConfigResourceHandle creates a new resource handle for Apdex configurations
func NewApdexConfigResourceHandle() resourcehandle.ResourceHandle[*api.ApdexConfig] {
	return &apdexConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaApdexConfig,
			Schema:        buildApdexConfigSchema(),
			SchemaVersion: 0,
		},
	}
}

type apdexConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *apdexConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *apdexConfigResource) GetRestResource(api client.InstanaAPI) rest.RestResource[*api.ApdexConfig] {
	return api.ApdexConfigs()
}

func (r *apdexConfigResource) SetComputedFields(_ context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	// No computed fields to set for Apdex config
	return diag.Diagnostics{}
}

// buildApdexConfigSchema builds the schema for the Apdex configuration resource
func buildApdexConfigSchema() schema.Schema {
	return schema.Schema{
		Description: ApdexConfigDescResource,
		Attributes: map[string]schema.Attribute{
			ApdexConfigFieldID:          buildIDAttribute(),
			ApdexConfigFieldApdexName:   buildApdexNameAttribute(),
			ApdexConfigFieldTags:        buildTagsAttribute(),
			ApdexConfigFieldRbacTags:    buildRbacTagsAttribute(),
			ApdexConfigFieldApdexEntity: buildApdexEntityAttribute(),
		},
	}
}

func buildIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: ApdexConfigDescID,
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func buildApdexNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: ApdexConfigDescApdexName,
		Required:    true,
	}
}

func buildTagsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		Description: ApdexConfigDescTags,
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
}

func buildRbacTagsAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Computed:    true,
		Description: ApdexConfigDescRbacTags,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				SchemaFieldDisplayName: schema.StringAttribute{
					Description: ApdexConfigDescRbacTagDisplayName,
					Required:    true,
				},
				SchemaFieldID: schema.StringAttribute{
					Description: ApdexConfigDescRbacTagID,
					Required:    true,
				},
			},
		},
	}
}

// buildApdexEntityAttribute builds the polymorphic entity attribute
func buildApdexEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: ApdexConfigDescApdexEntity,
		Required:    true,
		Attributes: map[string]schema.Attribute{
			SchemaFieldApplication: buildApplicationEntityAttribute(),
			SchemaFieldWebsite:     buildWebsiteEntityAttribute(),
		},
	}
}

// buildApplicationEntityAttribute builds the application entity nested attribute
func buildApplicationEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: ApdexConfigDescApplicationEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ApdexConfigFieldEntityID: schema.StringAttribute{
				Required:    true,
				Description: ApdexConfigDescEntityID,
			},
			ApdexConfigFieldThreshold: schema.Int64Attribute{
				Required:    true,
				Description: ApdexConfigDescThreshold,
			},
			ApdexConfigFieldBoundaryScope: schema.StringAttribute{
				Required:    true,
				Description: ApdexConfigDescBoundaryScope,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ALL",
						"INBOUND",
					),
				},
			},
			ApdexConfigFieldIncludeInternal: schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: ApdexConfigDescIncludeInternal,
			},
			ApdexConfigFieldIncludeSynthetic: schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: ApdexConfigDescIncludeSynthetic,
			},
			ApdexConfigFieldFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: ApdexConfigDescFilterExpression,
			},
		},
	}
}

// buildWebsiteEntityAttribute builds the website entity nested attribute
func buildWebsiteEntityAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: ApdexConfigDescWebsiteEntity,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ApdexConfigFieldEntityID: schema.StringAttribute{
				Required:    true,
				Description: ApdexConfigDescEntityID,
			},
			ApdexConfigFieldThreshold: schema.Int64Attribute{
				Required:    true,
				Description: ApdexConfigDescThreshold,
			},
			ApdexConfigFieldBeaconType: schema.StringAttribute{
				Required:    true,
				Description: ApdexConfigDescBeaconType,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"pageLoad",
						"httpRequest",
						"custom",
					),
				},
			},
			ApdexConfigFieldFilterExpression: schema.StringAttribute{
				Optional:    true,
				Description: ApdexConfigDescFilterExpression,
			},
		},
	}
}

// MapStateToDataObject maps Terraform state to API data object
func (r *apdexConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.ApdexConfig, diag.Diagnostics) {
	var model ApdexConfigModel
	diags := r.extractModelFromPlanOrState(ctx, plan, state, &model)
	if diags.HasError() {
		return nil, diags
	}

	tags := r.mapTagsFromState(model.Tags)
	RbacTags := r.mapRbacTagsFromState(model.RbacTags)
	entity, entityDiags := r.mapEntityFromState(model.ApdexEntity)
	diags.Append(entityDiags...)
	if diags.HasError() {
		return nil, diags
	}

	apdexConfig := &api.ApdexConfig{
		ApdexName:   model.ApdexName.ValueString(),
		ApdexEntity: entity,
		Tags:        tags,
		RbacTags:    RbacTags,
	}

	if !model.ID.IsNull() {
		apdexConfig.ID = model.ID.ValueString()
	}

	return apdexConfig, diags
}

func (r *apdexConfigResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State, model *ApdexConfigModel) diag.Diagnostics {
	if plan != nil {
		return plan.Get(ctx, model)
	}
	if state != nil {
		return state.Get(ctx, model)
	}
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(ApdexConfigErrMappingState, ApdexConfigErrBothPlanStateNil),
	}
}

func (r *apdexConfigResource) mapTagsFromState(tagsSet types.Set) []string {
	if tagsSet.IsNull() || tagsSet.IsUnknown() {
		return []string{}
	}

	var tags []string
	tagsSet.ElementsAs(context.Background(), &tags, false)
	return tags
}

// mapRbacTagsFromState maps RBAC tags from Terraform state List to API model
func (r *apdexConfigResource) mapRbacTagsFromState(rbacTagsList types.List) []api.RbacTag {
	if rbacTagsList.IsNull() || rbacTagsList.IsUnknown() {
		return []api.RbacTag{}
	}

	var tagModels []RbacTagModel
	rbacTagsList.ElementsAs(context.Background(), &tagModels, false)

	if len(tagModels) == 0 {
		return []api.RbacTag{}
	}

	rbacTags := make([]api.RbacTag, 0, len(tagModels))
	for _, model := range tagModels {
		rbacTags = append(rbacTags, api.RbacTag{
			DisplayName: model.DisplayName.ValueString(),
			ID:          model.ID.ValueString(),
		})
	}
	return rbacTags
}

// mapEntityFromState maps entity model to API entity
func (r *apdexConfigResource) mapEntityFromState(entityModel *ApdexEntityModel) (api.ApdexEntity, diag.Diagnostics) {
	if entityModel == nil {
		var diags diag.Diagnostics
		diags.AddError(
			ApdexConfigErrMissingEntity,
			ApdexConfigErrExactlyOneEntity,
		)
		return api.ApdexEntity{}, diags
	}

	if entityModel.ApplicationEntityModel != nil {
		return r.validateAndMapApplicationEntity(entityModel.ApplicationEntityModel)
	}

	if entityModel.WebsiteEntityModel != nil {
		return r.validateAndMapWebsiteEntity(entityModel.WebsiteEntityModel)
	}

	var diags diag.Diagnostics
	diags.AddError(
		ApdexConfigErrMissingEntity,
		ApdexConfigErrExactlyOneEntity,
	)
	return api.ApdexEntity{}, diags
}

// validateAndMapApplicationEntity validates and maps application entity from state
func (r *apdexConfigResource) validateAndMapApplicationEntity(model *ApplicationApdexEntityModel) (api.ApdexEntity, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate required fields
	if model.EntityID.IsNull() || model.EntityID.IsUnknown() ||
		model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
		model.BoundaryScope.IsNull() || model.BoundaryScope.IsUnknown() {
		diags.AddError(
			ApdexConfigErrMissingEntity,
			ApdexConfigErrApplicationEntityRequired,
		)
		return api.ApdexEntity{}, diags
	}

	entityID := model.EntityID.ValueString()
	threshold := int(model.Threshold.ValueInt64())
	entity := api.ApdexEntity{
		Type:          "application",
		EntityID:      entityID,
		Threshold:     threshold,
		BoundaryScope: util.SetStringPointerFromState(model.BoundaryScope),
	}

	// Optional fields with defaults
	if !model.IncludeInternal.IsNull() {
		includeInternal := model.IncludeInternal.ValueBool()
		entity.IncludeInternal = &includeInternal
	} else {
		// Default to false
		includeInternal := false
		entity.IncludeInternal = &includeInternal
	}

	if !model.IncludeSynthetic.IsNull() {
		includeSynthetic := model.IncludeSynthetic.ValueBool()
		entity.IncludeSynthetic = &includeSynthetic
	} else {
		// Default to false
		includeSynthetic := false
		entity.IncludeSynthetic = &includeSynthetic
	}

	// Handle filter expression - always set TagFilter
	tagFilter, filterDiags := r.mapFilterExpressionToEntity(model.FilterExpression)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return api.ApdexEntity{}, diags
	}
	entity.TagFilter = tagFilter

	return entity, diags
}

// validateAndMapWebsiteEntity validates and maps website entity from state
func (r *apdexConfigResource) validateAndMapWebsiteEntity(model *WebsiteApdexEntityModel) (api.ApdexEntity, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate required fields
	if model.EntityID.IsNull() || model.EntityID.IsUnknown() ||
		model.Threshold.IsNull() || model.Threshold.IsUnknown() ||
		model.BeaconType.IsNull() || model.BeaconType.IsUnknown() {
		diags.AddError(
			ApdexConfigErrMissingEntity,
			ApdexConfigErrWebsiteEntityRequired,
		)
		return api.ApdexEntity{}, diags
	}

	entityID := model.EntityID.ValueString()
	threshold := int(model.Threshold.ValueInt64())
	entity := api.ApdexEntity{
		Type:       "website",
		EntityID:   entityID,
		Threshold:  threshold,
		BeaconType: util.SetStringPointerFromState(model.BeaconType),
	}

	// Handle filter expression - always set TagFilter
	tagFilter, filterDiags := r.mapFilterExpressionToEntity(model.FilterExpression)
	diags.Append(filterDiags...)
	if diags.HasError() {
		return api.ApdexEntity{}, diags
	}
	entity.TagFilter = tagFilter

	return entity, diags
}

// mapFilterExpressionToEntity converts filter expression to API model
func (r *apdexConfigResource) mapFilterExpressionToEntity(filterExpression types.String) (*tag.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !filterExpression.IsNull() && !filterExpression.IsUnknown() && filterExpression.ValueString() != "" {
		parser := tagfilter.NewParser()
		filterExpr, err := parser.Parse(filterExpression.ValueString())
		if err != nil {
			diags.AddError(
				ApdexConfigErrParsingFilterExpression,
				fmt.Sprintf(ApdexConfigErrParsingFilterExpressionMsg, err.Error()),
			)
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		return mapper.ToAPIModel(filterExpr), diags
	}

	return r.createDefaultTagFilter(), diags
}

// createDefaultTagFilter creates a default empty tag filter
func (r *apdexConfigResource) createDefaultTagFilter() *tag.TagFilter {
	operator := common.LogicalOperatorType(LogicalOperatorAnd)
	return &tag.TagFilter{
		Type:            tag.TagFilterExpressionElementType(TagFilterTypeExpression),
		LogicalOperator: &operator,
		Elements:        []*tag.TagFilter{},
	}
}

// UpdateState updates Terraform state from API data object
func (r *apdexConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, data *api.ApdexConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model ApdexConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	model.ID = types.StringValue(data.ID)
	model.ApdexName = types.StringValue(data.ApdexName)
	if len(data.Tags) > 0 {
		model.Tags = r.mapTagsToState(data.Tags)
	}
	if len(data.RbacTags) > 0 {
		model.RbacTags = r.mapRbacTagsToState(data.RbacTags)
	}

	entityModel, entityDiags := r.mapEntityToState(data.ApdexEntity)
	diags.Append(entityDiags...)
	if diags.HasError() {
		return diags
	}
	model.ApdexEntity = &entityModel

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapTagsToState converts tags from API to state
func (r *apdexConfigResource) mapTagsToState(tags []string) types.Set {
	if tags == nil {
		return types.SetNull(types.StringType)
	}

	elements := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		elements = append(elements, types.StringValue(tag))
	}
	return types.SetValueMust(types.StringType, elements)
}

// mapRbacTagsToState converts RBAC tags from API to state List
func (r *apdexConfigResource) mapRbacTagsToState(rbacTags []api.RbacTag) types.List {
	tagAttrTypes := map[string]attr.Type{
		SchemaFieldID:          types.StringType,
		SchemaFieldDisplayName: types.StringType,
	}

	// Always initialize with empty list, even if data is null or empty
	if rbacTags == nil || len(rbacTags) == 0 {
		emptyList, _ := types.ListValue(
			types.ObjectType{AttrTypes: tagAttrTypes},
			[]attr.Value{},
		)
		return emptyList
	}

	tagValues := make([]attr.Value, len(rbacTags))
	for i, tag := range rbacTags {
		tagObj, _ := types.ObjectValue(
			tagAttrTypes,
			map[string]attr.Value{
				SchemaFieldID:          types.StringValue(tag.ID),
				SchemaFieldDisplayName: types.StringValue(tag.DisplayName),
			},
		)
		tagValues[i] = tagObj
	}

	list, _ := types.ListValue(types.ObjectType{AttrTypes: tagAttrTypes}, tagValues)
	return list
}

// mapEntityToState maps API entity to Terraform state
func (r *apdexConfigResource) mapEntityToState(entity api.ApdexEntity) (ApdexEntityModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := ApdexEntityModel{}

	if entity.Type == "application" {
		appModel := r.buildApplicationEntityModel(entity)
		model.ApplicationEntityModel = &appModel
	} else if entity.Type == "website" {
		webModel := r.buildWebsiteEntityModel(entity)
		model.WebsiteEntityModel = &webModel
	} else {
		diags.AddError(
			ApdexConfigErrMappingEntityToState,
			fmt.Sprintf(ApdexConfigErrUnsupportedEntityType, entity.Type),
		)
	}

	return model, diags
}

// buildApplicationEntityModel builds application entity model from API entity
func (r *apdexConfigResource) buildApplicationEntityModel(entity api.ApdexEntity) ApplicationApdexEntityModel {
	model := ApplicationApdexEntityModel{
		EntityID:  types.StringValue(entity.EntityID),
		Threshold: types.Int64Value(int64(entity.Threshold)),
	}

	if entity.BoundaryScope != nil {
		model.BoundaryScope = types.StringValue(*entity.BoundaryScope)
	}

	if entity.IncludeInternal != nil {
		model.IncludeInternal = types.BoolValue(*entity.IncludeInternal)
	} else {
		model.IncludeInternal = types.BoolValue(false)
	}

	if entity.IncludeSynthetic != nil {
		model.IncludeSynthetic = types.BoolValue(*entity.IncludeSynthetic)
	} else {
		model.IncludeSynthetic = types.BoolValue(false)
	}

	// Map filter expression
	if entity.TagFilter != nil {
		expression, err := tagfilter.MapTagFilterToNormalizedString(entity.TagFilter)
		if err == nil && expression != nil {
			model.FilterExpression = types.StringValue(*expression)
		} else {
			model.FilterExpression = types.StringNull()
		}
	} else {
		model.FilterExpression = types.StringNull()
	}

	return model
}

// buildWebsiteEntityModel builds website entity model from API entity
func (r *apdexConfigResource) buildWebsiteEntityModel(entity api.ApdexEntity) WebsiteApdexEntityModel {
	model := WebsiteApdexEntityModel{
		EntityID:  types.StringValue(entity.EntityID),
		Threshold: types.Int64Value(int64(entity.Threshold)),
	}

	if entity.BeaconType != nil {
		model.BeaconType = types.StringValue(*entity.BeaconType)
	}

	// Map filter expression
	if entity.TagFilter != nil {
		expression, err := tagfilter.MapTagFilterToNormalizedString(entity.TagFilter)
		if err == nil && expression != nil {
			model.FilterExpression = types.StringValue(*expression)
		} else {
			model.FilterExpression = types.StringNull()
		}
	} else {
		model.FilterExpression = types.StringNull()
	}

	return model
}

func (r *apdexConfigResource) GetID(data *api.ApdexConfig) (string, diag.Diagnostics) {
	return data.ID, diag.Diagnostics{}
}

func (r *apdexConfigResource) GetData(id string) (*api.ApdexConfig, diag.Diagnostics) {
	return &api.ApdexConfig{ID: id}, diag.Diagnostics{}
}

func (r *apdexConfigResource) GetUpdatedDataFromState(data *api.ApdexConfig, state *tfsdk.State) (*api.ApdexConfig, diag.Diagnostics) {
	var model ApdexConfigModel
	diags := state.Get(context.Background(), &model)
	if diags.HasError() {
		return nil, diags
	}

	rbacTags := r.mapRbacTagsFromState(model.RbacTags)
	data.RbacTags = rbacTags

	return data, diags
}

func (r *apdexConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
