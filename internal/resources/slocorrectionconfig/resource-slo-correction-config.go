package slocorrectionconfig

import (
	"context"
	"strings"

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
)

// NewSloCorrectionConfigResourceHandle creates the resource handle for SLO Correction Config
func NewSloCorrectionConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.SloCorrectionConfig] {
	resource := &sloCorrectionConfigResource{}
	return resource.initialize()
}

// initialize sets up the resource with metadata and schema
func (r *sloCorrectionConfigResource) initialize() *sloCorrectionConfigResource {
	r.metaData = resourcehandle.ResourceMetaData{
		ResourceName:  ResourceInstanaSloCorrectionConfig,
		Schema:        r.buildSchema(),
		SchemaVersion: 1,
	}
	return r
}

// buildSchema constructs the complete schema for the resource
func (r *sloCorrectionConfigResource) buildSchema() schema.Schema {
	return schema.Schema{
		Description: SloCorrectionConfigDescResource,
		Attributes:  r.buildSchemaAttributes(),
	}
}

// buildSchemaAttributes constructs the top-level schema attributes
func (r *sloCorrectionConfigResource) buildSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SloCorrectionConfigFieldID:          r.buildIDAttribute(),
		SloCorrectionConfigFieldName:        r.buildNameAttribute(),
		SloCorrectionConfigFieldDescription: r.buildDescriptionAttribute(),
		SloCorrectionConfigFieldActive:      r.buildActiveAttribute(),
		SloCorrectionConfigFieldSloIds:      r.buildSloIdsAttribute(),
		SloCorrectionConfigFieldTags:        r.buildTagsAttribute(),
		SloCorrectionConfigFieldScheduling:  r.buildSchedulingAttribute(),
	}
}

// buildIDAttribute creates the ID attribute schema
func (r *sloCorrectionConfigResource) buildIDAttribute() schema.Attribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: SloCorrectionConfigDescID,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// buildNameAttribute creates the name attribute schema
func (r *sloCorrectionConfigResource) buildNameAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloCorrectionConfigDescName,
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
		},
	}
}

// buildDescriptionAttribute creates the description attribute schema
func (r *sloCorrectionConfigResource) buildDescriptionAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloCorrectionConfigDescDescription,
	}
}

// buildActiveAttribute creates the active attribute schema
func (r *sloCorrectionConfigResource) buildActiveAttribute() schema.Attribute {
	return schema.BoolAttribute{
		Required:    true,
		Description: SloCorrectionConfigDescActive,
	}
}

// buildSloIdsAttribute creates the slo_ids attribute schema
func (r *sloCorrectionConfigResource) buildSloIdsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Required:    true,
		Description: SloCorrectionConfigDescSloIds,
		ElementType: types.StringType,
	}
}

// buildTagsAttribute creates the tags attribute schema
func (r *sloCorrectionConfigResource) buildTagsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Optional:    true,
		Description: SloCorrectionConfigDescTags,
		ElementType: types.StringType,
	}
}

// buildSchedulingAttribute creates the scheduling nested attribute schema
func (r *sloCorrectionConfigResource) buildSchedulingAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Required:    true,
		Description: SloCorrectionConfigDescScheduling,
		Attributes:  r.buildSchedulingNestedAttributes(),
	}
}

// buildSchedulingNestedAttributes constructs the nested scheduling attributes
func (r *sloCorrectionConfigResource) buildSchedulingNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SloCorrectionConfigFieldSchedulingStartTime:     r.buildStartTimeAttribute(),
		SloCorrectionConfigFieldSchedulingDuration:      r.buildDurationAttribute(),
		SloCorrectionConfigFieldSchedulingDurationUnit:  r.buildDurationUnitAttribute(),
		SloCorrectionConfigFieldSchedulingRecurrentRule: r.buildRecurrentRuleAttribute(),
		SloCorrectionConfigFieldSchedulingRecurrent:     r.buildRecurrentAttribute(),
	}
}

// buildStartTimeAttribute creates the start_time attribute schema
func (r *sloCorrectionConfigResource) buildStartTimeAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Required:    true,
		Description: SloCorrectionConfigDescStartTime,
	}
}

// buildDurationAttribute creates the duration attribute schema
func (r *sloCorrectionConfigResource) buildDurationAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Required:    true,
		Description: SloCorrectionConfigDescDuration,
	}
}

// buildDurationUnitAttribute creates the duration_unit attribute schema with validators
func (r *sloCorrectionConfigResource) buildDurationUnitAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloCorrectionConfigDescDurationUnit,
		Validators:  r.buildDurationUnitValidators(),
	}
}

// buildDurationUnitValidators creates validators for duration unit
func (r *sloCorrectionConfigResource) buildDurationUnitValidators() []validator.String {
	return []validator.String{
		stringvalidator.OneOf(
			DurationUnitMillisecond,
			DurationUnitSecond,
			DurationUnitMinute,
			DurationUnitHour,
			DurationUnitDay,
			DurationUnitWeek,
			DurationUnitMonth,
		),
	}
}

// buildRecurrentRuleAttribute creates the recurrent_rule attribute schema
func (r *sloCorrectionConfigResource) buildRecurrentRuleAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SloCorrectionConfigDescRecurrentRule,
	}
}

// buildRecurrentAttribute creates the recurrent attribute schema
func (r *sloCorrectionConfigResource) buildRecurrentAttribute() schema.Attribute {
	return schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: SloCorrectionConfigDescRecurrent,
	}
}

type sloCorrectionConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *sloCorrectionConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *sloCorrectionConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloCorrectionConfig] {
	return api.SloCorrectionConfig()
}

func (r *sloCorrectionConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sloCorrectionConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloCorrectionConfig, diag.Diagnostics) {
	model, diags := r.extractModelFromState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	scheduling, schedulingDiags := r.mapSchedulingFromModel(model)
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return nil, diags
	}

	sloIds, sloIdsDiags := r.extractSloIdsFromModel(ctx, model)
	diags.Append(sloIdsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	tags, tagsDiags := r.extractTagsFromModel(ctx, model)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	sloCorrectionConfig := r.buildAPIObjectFromModel(model, scheduling, sloIds, tags)
	return sloCorrectionConfig, diags
}

// extractModelFromState retrieves the model from plan or state
func (r *sloCorrectionConfigResource) extractModelFromState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*SloCorrectionConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SloCorrectionConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &model, diags
}

// mapSchedulingFromModel converts scheduling model to API scheduling object
func (r *sloCorrectionConfigResource) mapSchedulingFromModel(model *SloCorrectionConfigModel) (restapi.Scheduling, diag.Diagnostics) {
	var diags diag.Diagnostics
	var scheduling restapi.Scheduling

	if model.Scheduling == nil {
		return scheduling, diags
	}

	scheduling = restapi.Scheduling{
		StartTime:    model.Scheduling.StartTime.ValueInt64(),
		Duration:     int(model.Scheduling.Duration.ValueInt64()),
		DurationUnit: r.convertDurationUnitToAPI(model.Scheduling.DurationUnit.ValueString()),
		Recurrent:    model.Scheduling.Recurrent.ValueBool(),
	}

	if !model.Scheduling.RecurrentRule.IsNull() {
		scheduling.RecurrentRule = model.Scheduling.RecurrentRule.ValueString()
	}

	return scheduling, diags
}

// convertDurationUnitToAPI converts duration unit string to API format (uppercase)
func (r *sloCorrectionConfigResource) convertDurationUnitToAPI(unit string) restapi.DurationUnit {
	return restapi.DurationUnit(strings.ToUpper(unit))
}

// extractSloIdsFromModel extracts SLO IDs from the model
func (r *sloCorrectionConfigResource) extractSloIdsFromModel(ctx context.Context, model *SloCorrectionConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var sloIds []string

	if !model.SloIds.IsNull() {
		diags.Append(model.SloIds.ElementsAs(ctx, &sloIds, false)...)
	}

	return sloIds, diags
}

// extractTagsFromModel extracts tags from the model
func (r *sloCorrectionConfigResource) extractTagsFromModel(ctx context.Context, model *SloCorrectionConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var tags []string

	if !model.Tags.IsNull() {
		diags.Append(model.Tags.ElementsAs(ctx, &tags, false)...)
	}

	return tags, diags
}

// buildAPIObjectFromModel constructs the API object from model and extracted data
func (r *sloCorrectionConfigResource) buildAPIObjectFromModel(
	model *SloCorrectionConfigModel,
	scheduling restapi.Scheduling,
	sloIds []string,
	tags []string,
) *restapi.SloCorrectionConfig {
	return &restapi.SloCorrectionConfig{
		ID:          model.ID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Active:      model.Active.ValueBool(),
		Scheduling:  scheduling,
		SloIds:      sloIds,
		Tags:        tags,
	}
}

func (r *sloCorrectionConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SloCorrectionConfig) diag.Diagnostics {
	model := r.buildModelFromAPIObject(apiObject)

	schedulingModel, schedulingDiags := r.mapSchedulingToModel(apiObject.Scheduling)
	var diags diag.Diagnostics
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return diags
	}
	model.Scheduling = schedulingModel

	sloIdsSet, sloIdsDiags := r.mapSloIdsToState(ctx, apiObject.SloIds)
	diags.Append(sloIdsDiags...)
	if diags.HasError() {
		return diags
	}
	model.SloIds = sloIdsSet

	tagsSet, tagsDiags := r.mapTagsToState(ctx, apiObject.Tags)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return diags
	}
	model.Tags = tagsSet

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// buildModelFromAPIObject creates a model with basic fields from API object
func (r *sloCorrectionConfigResource) buildModelFromAPIObject(apiObject *restapi.SloCorrectionConfig) SloCorrectionConfigModel {
	return SloCorrectionConfigModel{
		ID:          types.StringValue(apiObject.ID),
		Name:        types.StringValue(apiObject.Name),
		Description: types.StringValue(apiObject.Description),
		Active:      types.BoolValue(apiObject.Active),
	}
}

// mapSchedulingToModel converts API scheduling to model scheduling
func (r *sloCorrectionConfigResource) mapSchedulingToModel(scheduling restapi.Scheduling) (*SchedulingModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	recurrentRuleValue := r.buildRecurrentRuleValue(scheduling.RecurrentRule)

	schedulingModel := &SchedulingModel{
		StartTime:     types.Int64Value(scheduling.StartTime),
		Duration:      types.Int64Value(int64(scheduling.Duration)),
		DurationUnit:  types.StringValue(string(scheduling.DurationUnit)),
		RecurrentRule: recurrentRuleValue,
		Recurrent:     types.BoolValue(scheduling.Recurrent),
	}

	return schedulingModel, diags
}

// buildRecurrentRuleValue creates a types.String value for recurrent rule
func (r *sloCorrectionConfigResource) buildRecurrentRuleValue(recurrentRule string) types.String {
	if recurrentRule == "" {
		return types.StringNull()
	}
	return types.StringValue(recurrentRule)
}

// mapSloIdsToState converts SLO IDs array to Terraform set
func (r *sloCorrectionConfigResource) mapSloIdsToState(ctx context.Context, sloIds []string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, sloIds)
}

// mapTagsToState converts tags array to Terraform set
func (r *sloCorrectionConfigResource) mapTagsToState(ctx context.Context, tags []string) (types.Set, diag.Diagnostics) {
	if len(tags) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, tags)
}
