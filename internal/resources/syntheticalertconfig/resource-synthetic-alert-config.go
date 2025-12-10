package syntheticalertconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewSyntheticAlertConfigResourceHandle creates the resource handle for Synthetic Alert Configuration
func NewSyntheticAlertConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.SyntheticAlertConfig] {
	resource := &syntheticAlertConfigResource{}
	return resource.initialize()
}

// initialize sets up the resource with metadata and schema
func (r *syntheticAlertConfigResource) initialize() *syntheticAlertConfigResource {
	r.metaData = resourcehandle.ResourceMetaData{
		ResourceName:  ResourceInstanaSyntheticAlertConfig,
		Schema:        r.buildSchema(),
		SchemaVersion: 2,
	}
	return r
}

// buildSchema constructs the complete schema for the resource
func (r *syntheticAlertConfigResource) buildSchema() schema.Schema {
	return schema.Schema{
		Description: SyntheticAlertConfigDescResource,
		Attributes:  r.buildSchemaAttributes(),
	}
}

// buildSchemaAttributes constructs the top-level schema attributes
func (r *syntheticAlertConfigResource) buildSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticAlertConfigFieldID:                 r.buildIDAttribute(),
		SyntheticAlertConfigFieldName:               r.buildNameAttribute(),
		SyntheticAlertConfigFieldDescription:        r.buildDescriptionAttribute(),
		SyntheticAlertConfigFieldSyntheticTestIds:   r.buildSyntheticTestIdsAttribute(),
		SyntheticAlertConfigFieldSeverity:           r.buildSeverityAttribute(),
		SyntheticAlertConfigFieldTagFilter:          r.buildTagFilterAttribute(),
		SyntheticAlertConfigFieldAlertChannelIds:    r.buildAlertChannelIdsAttribute(),
		SyntheticAlertConfigFieldGracePeriod:        r.buildGracePeriodAttribute(),
		SyntheticAlertConfigFieldCustomPayloadField: shared.GetCustomPayloadFieldsSchema(),
		SyntheticAlertConfigFieldRule:               r.buildRuleAttribute(),
		SyntheticAlertConfigFieldTimeThreshold:      r.buildTimeThresholdAttribute(),
	}
}

// buildIDAttribute creates the ID attribute schema
func (r *syntheticAlertConfigResource) buildIDAttribute() schema.Attribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: SyntheticAlertConfigDescID,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// buildNameAttribute creates the name attribute schema
func (r *syntheticAlertConfigResource) buildNameAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescName,
		Validators:  r.buildNameValidators(),
	}
}

// buildNameValidators creates validators for the name field
func (r *syntheticAlertConfigResource) buildNameValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthBetween(1, 256),
	}
}

// buildDescriptionAttribute creates the description attribute schema
func (r *syntheticAlertConfigResource) buildDescriptionAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescDescription,
		Validators:  r.buildDescriptionValidators(),
	}
}

// buildDescriptionValidators creates validators for the description field
func (r *syntheticAlertConfigResource) buildDescriptionValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthBetween(0, 1024),
	}
}

// buildSyntheticTestIdsAttribute creates the synthetic_test_ids attribute schema
func (r *syntheticAlertConfigResource) buildSyntheticTestIdsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescSyntheticTestIds,
		ElementType: types.StringType,
	}
}

// buildSeverityAttribute creates the severity attribute schema
func (r *syntheticAlertConfigResource) buildSeverityAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescSeverity,
		Validators:  r.buildSeverityValidators(),
	}
}

// buildSeverityValidators creates validators for the severity field
func (r *syntheticAlertConfigResource) buildSeverityValidators() []validator.Int64 {
	return []validator.Int64{
		int64validator.OneOf(5, 10),
	}
}

// buildTagFilterAttribute creates the tag_filter attribute schema
func (r *syntheticAlertConfigResource) buildTagFilterAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescTagFilter,
	}
}

// buildAlertChannelIdsAttribute creates the alert_channel_ids attribute schema
func (r *syntheticAlertConfigResource) buildAlertChannelIdsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescAlertChannelIds,
		ElementType: types.StringType,
	}
}

// buildGracePeriodAttribute creates the grace_period attribute schema
func (r *syntheticAlertConfigResource) buildGracePeriodAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescGracePeriod,
	}
}

// buildRuleAttribute creates the rule nested attribute schema
func (r *syntheticAlertConfigResource) buildRuleAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: SyntheticAlertConfigDescRule,
		Required:    true,
		Attributes:  r.buildRuleNestedAttributes(),
	}
}

// buildRuleNestedAttributes constructs the nested rule attributes
func (r *syntheticAlertConfigResource) buildRuleNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticAlertRuleFieldAlertType:   r.buildRuleAlertTypeAttribute(),
		SyntheticAlertRuleFieldMetricName:  r.buildRuleMetricNameAttribute(),
		SyntheticAlertRuleFieldAggregation: r.buildRuleAggregationAttribute(),
	}
}

// buildRuleAlertTypeAttribute creates the rule alert_type attribute schema
func (r *syntheticAlertConfigResource) buildRuleAlertTypeAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescRuleAlertType,
		Validators:  r.buildRuleAlertTypeValidators(),
	}
}

// buildRuleAlertTypeValidators creates validators for the rule alert_type field
func (r *syntheticAlertConfigResource) buildRuleAlertTypeValidators() []validator.String {
	return []validator.String{
		stringvalidator.OneOf(SyntheticAlertConfigValidAlertType),
	}
}

// buildRuleMetricNameAttribute creates the rule metric_name attribute schema
func (r *syntheticAlertConfigResource) buildRuleMetricNameAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescRuleMetricName,
		Validators:  r.buildRuleMetricNameValidators(),
	}
}

// buildRuleMetricNameValidators creates validators for the rule metric_name field
func (r *syntheticAlertConfigResource) buildRuleMetricNameValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthBetween(1, 256),
	}
}

// buildRuleAggregationAttribute creates the rule aggregation attribute schema
func (r *syntheticAlertConfigResource) buildRuleAggregationAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: SyntheticAlertConfigDescRuleAggregation,
		Validators:  r.buildRuleAggregationValidators(),
	}
}

// buildRuleAggregationValidators creates validators for the rule aggregation field
func (r *syntheticAlertConfigResource) buildRuleAggregationValidators() []validator.String {
	return []validator.String{
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
			AggregationTypeDistinctCount,
			AggregationTypeSumPositive,
			AggregationTypePerSecond,
			AggregationTypeIncrease,
		),
	}
}

// buildTimeThresholdAttribute creates the time_threshold nested attribute schema
func (r *syntheticAlertConfigResource) buildTimeThresholdAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: SyntheticAlertConfigDescTimeThreshold,
		Required:    true,
		Attributes:  r.buildTimeThresholdNestedAttributes(),
	}
}

// buildTimeThresholdNestedAttributes constructs the nested time threshold attributes
func (r *syntheticAlertConfigResource) buildTimeThresholdNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticAlertTimeThresholdFieldType:            r.buildTimeThresholdTypeAttribute(),
		SyntheticAlertTimeThresholdFieldViolationsCount: r.buildTimeThresholdViolationsCountAttribute(),
	}
}

// buildTimeThresholdTypeAttribute creates the time threshold type attribute schema
func (r *syntheticAlertConfigResource) buildTimeThresholdTypeAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescTimeThresholdType,
		Validators:  r.buildTimeThresholdTypeValidators(),
	}
}

// buildTimeThresholdTypeValidators creates validators for the time threshold type field
func (r *syntheticAlertConfigResource) buildTimeThresholdTypeValidators() []validator.String {
	return []validator.String{
		stringvalidator.OneOf(SyntheticAlertConfigValidTimeThresholdType),
	}
}

// buildTimeThresholdViolationsCountAttribute creates the time threshold violations_count attribute schema
func (r *syntheticAlertConfigResource) buildTimeThresholdViolationsCountAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Required:    true,
		Description: SyntheticAlertConfigDescTimeThresholdViolationsCount,
		Validators:  r.buildTimeThresholdViolationsCountValidators(),
	}
}

// buildTimeThresholdViolationsCountValidators creates validators for the time threshold violations_count field
func (r *syntheticAlertConfigResource) buildTimeThresholdViolationsCountValidators() []validator.Int64 {
	return []validator.Int64{
		int64validator.Between(1, 12),
	}
}

type syntheticAlertConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *syntheticAlertConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *syntheticAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SyntheticAlertConfig] {
	return api.SyntheticAlertConfigs()
}

func (r *syntheticAlertConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *syntheticAlertConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SyntheticAlertConfig, diag.Diagnostics) {
	model, diags := r.extractModelFromState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	rule, ruleDiags := r.mapRuleFromModel(model)
	diags.Append(ruleDiags...)
	if diags.HasError() {
		return nil, diags
	}

	timeThreshold, thresholdDiags := r.mapTimeThresholdFromModel(model)
	diags.Append(thresholdDiags...)
	if diags.HasError() {
		return nil, diags
	}

	syntheticTestIds, testIdsDiags := r.extractSyntheticTestIdsFromModel(ctx, model)
	diags.Append(testIdsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	alertChannelIds, channelIdsDiags := r.extractAlertChannelIdsFromModel(ctx, model)
	diags.Append(channelIdsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	tagFilter, tagFilterDiags := r.mapTagFilterFromModel(model)
	diags.Append(tagFilterDiags...)
	if diags.HasError() {
		return nil, diags
	}

	customPayloadFields, payloadDiags := r.extractCustomPayloadFieldsFromModel(ctx, model)
	diags.Append(payloadDiags...)
	if diags.HasError() {
		return nil, diags
	}

	syntheticAlertConfig := r.buildAPIObjectFromModel(model, rule, timeThreshold, syntheticTestIds, alertChannelIds, tagFilter, customPayloadFields)
	return syntheticAlertConfig, diags
}

// extractModelFromState retrieves the model from plan or state
func (r *syntheticAlertConfigResource) extractModelFromState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*SyntheticAlertConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SyntheticAlertConfigModel

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

// mapRuleFromModel converts rule model to API rule object
func (r *syntheticAlertConfigResource) mapRuleFromModel(model *SyntheticAlertConfigModel) (restapi.SyntheticAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	var rule restapi.SyntheticAlertRule

	if model.Rule == nil {
		return rule, diags
	}

	rule = restapi.SyntheticAlertRule{
		AlertType:  model.Rule.AlertType.ValueString(),
		MetricName: model.Rule.MetricName.ValueString(),
	}

	if !model.Rule.Aggregation.IsNull() {
		rule.Aggregation = model.Rule.Aggregation.ValueString()
	}

	return rule, diags
}

// mapTimeThresholdFromModel converts time threshold model to API time threshold object
func (r *syntheticAlertConfigResource) mapTimeThresholdFromModel(model *SyntheticAlertConfigModel) (restapi.SyntheticAlertTimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics
	var timeThreshold restapi.SyntheticAlertTimeThreshold

	if model.TimeThreshold == nil {
		return timeThreshold, diags
	}

	timeThreshold = restapi.SyntheticAlertTimeThreshold{
		Type:            model.TimeThreshold.Type.ValueString(),
		ViolationsCount: int(model.TimeThreshold.ViolationsCount.ValueInt64()),
	}

	return timeThreshold, diags
}

// extractSyntheticTestIdsFromModel extracts synthetic test IDs from the model
func (r *syntheticAlertConfigResource) extractSyntheticTestIdsFromModel(ctx context.Context, model *SyntheticAlertConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	syntheticTestIds := []string{}

	if !model.SyntheticTestIds.IsNull() && !model.SyntheticTestIds.IsUnknown() {
		diags.Append(model.SyntheticTestIds.ElementsAs(ctx, &syntheticTestIds, false)...)
	}

	// Ensure we always return an initialized empty slice, never nil
	if syntheticTestIds == nil {
		syntheticTestIds = []string{}
	}

	return syntheticTestIds, diags
}

// extractAlertChannelIdsFromModel extracts alert channel IDs from the model
func (r *syntheticAlertConfigResource) extractAlertChannelIdsFromModel(ctx context.Context, model *SyntheticAlertConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannelIds := []string{}

	if !model.AlertChannelIds.IsNull() && !model.AlertChannelIds.IsUnknown() {
		diags.Append(model.AlertChannelIds.ElementsAs(ctx, &alertChannelIds, false)...)
	}

	// Ensure we always return an initialized empty slice, never nil
	if alertChannelIds == nil {
		alertChannelIds = []string{}
	}

	return alertChannelIds, diags
}

// mapTagFilterFromModel converts tag filter model to API tag filter object
func (r *syntheticAlertConfigResource) mapTagFilterFromModel(model *SyntheticAlertConfigModel) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if r.hasTagFilterValue(model) {
		return r.parseTagFilterExpression(model.TagFilter.ValueString(), &diags)
	}

	return r.createDefaultTagFilter(), diags
}

// hasTagFilterValue checks if the model has a tag filter value
func (r *syntheticAlertConfigResource) hasTagFilterValue(model *SyntheticAlertConfigModel) bool {
	return !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown()
}

// parseTagFilterExpression parses the tag filter expression string
func (r *syntheticAlertConfigResource) parseTagFilterExpression(expression string, diags *diag.Diagnostics) (*restapi.TagFilter, diag.Diagnostics) {
	tagFilter, err := mapTagFilterExpressionFromSchema(expression)
	if err != nil {
		diags.AddError(
			SyntheticAlertConfigErrParsingTagFilter,
			SyntheticAlertConfigErrParsingTagFilterDetail+err.Error(),
		)
		return nil, *diags
	}
	return tagFilter, *diags
}

// createDefaultTagFilter creates a default tag filter with AND operator
func (r *syntheticAlertConfigResource) createDefaultTagFilter() *restapi.TagFilter {
	operator := restapi.LogicalOperatorType(TagFilterLogicalOperatorAnd)
	return &restapi.TagFilter{
		Type:            TagFilterTypeExpression,
		LogicalOperator: &operator,
		Elements:        []*restapi.TagFilter{},
	}
}

// extractCustomPayloadFieldsFromModel extracts custom payload fields from the model
func (r *syntheticAlertConfigResource) extractCustomPayloadFieldsFromModel(ctx context.Context, model *SyntheticAlertConfigModel) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics
	var customPayloadFields []restapi.CustomPayloadField[any]

	if !model.CustomPayloadFields.IsNull() {
		customPayloadFields, diags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
	}

	return customPayloadFields, diags
}

// buildAPIObjectFromModel constructs the API object from model and extracted data
func (r *syntheticAlertConfigResource) buildAPIObjectFromModel(
	model *SyntheticAlertConfigModel,
	rule restapi.SyntheticAlertRule,
	timeThreshold restapi.SyntheticAlertTimeThreshold,
	syntheticTestIds []string,
	alertChannelIds []string,
	tagFilter *restapi.TagFilter,
	customPayloadFields []restapi.CustomPayloadField[any],
) *restapi.SyntheticAlertConfig {
	syntheticAlertConfig := &restapi.SyntheticAlertConfig{
		ID:                    model.ID.ValueString(),
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		SyntheticTestIds:      syntheticTestIds,
		TagFilterExpression:   tagFilter,
		Rule:                  rule,
		AlertChannelIds:       alertChannelIds,
		TimeThreshold:         timeThreshold,
		CustomerPayloadFields: customPayloadFields,
	}

	r.setOptionalFieldsOnAPIObject(syntheticAlertConfig, model)
	return syntheticAlertConfig
}

// setOptionalFieldsOnAPIObject sets optional fields on the API object
func (r *syntheticAlertConfigResource) setOptionalFieldsOnAPIObject(apiObject *restapi.SyntheticAlertConfig, model *SyntheticAlertConfigModel) {
	if !model.Severity.IsNull() {
		apiObject.Severity = int(model.Severity.ValueInt64())
	}

	if !model.GracePeriod.IsNull() {
		apiObject.GracePeriod = model.GracePeriod.ValueInt64()
	}
}

func (r *syntheticAlertConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SyntheticAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model SyntheticAlertConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	model.ID = types.StringValue(apiObject.ID)
	model.Name = types.StringValue(apiObject.Name)
	model.Description = types.StringValue(apiObject.Description)
	model.Severity = types.Int64Value(int64(apiObject.Severity))

	if model.GracePeriod.IsNull() || model.GracePeriod.IsUnknown() {
		gracePeriodValue := r.mapGracePeriodToState(apiObject.GracePeriod)
		model.GracePeriod = gracePeriodValue
	}

	// to preserve the existing value in plan/state to handle the value drift
	if model.TagFilter.IsNull() || model.TagFilter.IsUnknown() {
		tagFilterValue, tagFilterDiags := r.mapTagFilterToState(apiObject.TagFilterExpression)
		diags.Append(tagFilterDiags...)
		if diags.HasError() {
			return diags
		}
		model.TagFilter = tagFilterValue
	}
	ruleModel := r.mapRuleToModel(apiObject.Rule)
	model.Rule = ruleModel

	timeThresholdModel := r.mapTimeThresholdToModel(apiObject.TimeThreshold)
	model.TimeThreshold = timeThresholdModel

	syntheticTestIdsSet, testIdsDiags := r.mapSyntheticTestIdsToState(ctx, apiObject.SyntheticTestIds)
	diags.Append(testIdsDiags...)
	if diags.HasError() {
		return diags
	}
	model.SyntheticTestIds = syntheticTestIdsSet

	alertChannelIdsSet, channelIdsDiags := r.mapAlertChannelIdsToState(ctx, apiObject.AlertChannelIds)
	diags.Append(channelIdsDiags...)
	if diags.HasError() {
		return diags
	}
	model.AlertChannelIds = alertChannelIdsSet

	customPayloadFieldsList, payloadDiags := r.mapCustomPayloadFieldsToState(ctx, apiObject.CustomerPayloadFields)
	diags.Append(payloadDiags...)
	if diags.HasError() {
		return diags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapGracePeriodToState converts grace period to state value
func (r *syntheticAlertConfigResource) mapGracePeriodToState(gracePeriod int64) types.Int64 {
	if gracePeriod > 0 {
		return types.Int64Value(gracePeriod)
	}
	return types.Int64Null()
}

// mapTagFilterToState converts tag filter expression to state value
func (r *syntheticAlertConfigResource) mapTagFilterToState(tagFilterExpression *restapi.TagFilter) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagFilterExpression == nil {
		return types.StringNull(), diags
	}

	normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(tagFilterExpression)
	if err != nil {
		diags.AddError(
			SyntheticAlertConfigErrNormalizingTagFilter,
			SyntheticAlertConfigErrNormalizingTagFilterDetail+err.Error(),
		)
		return types.StringNull(), diags
	}

	if normalizedTagFilterString != nil {
		return util.SetStringPointerToState(normalizedTagFilterString), diags
	}

	return types.StringNull(), diags
}

// mapRuleToModel converts API rule to model rule
func (r *syntheticAlertConfigResource) mapRuleToModel(rule restapi.SyntheticAlertRule) *SyntheticAlertRuleModel {
	aggregationValue := r.buildAggregationValue(rule.Aggregation)

	return &SyntheticAlertRuleModel{
		AlertType:   types.StringValue(rule.AlertType),
		MetricName:  types.StringValue(rule.MetricName),
		Aggregation: aggregationValue,
	}
}

// buildAggregationValue creates a types.String value for aggregation
func (r *syntheticAlertConfigResource) buildAggregationValue(aggregation string) types.String {
	if aggregation == "" {
		return types.StringNull()
	}
	return types.StringValue(aggregation)
}

// mapTimeThresholdToModel converts API time threshold to model time threshold
func (r *syntheticAlertConfigResource) mapTimeThresholdToModel(timeThreshold restapi.SyntheticAlertTimeThreshold) *SyntheticAlertTimeThresholdModel {
	return &SyntheticAlertTimeThresholdModel{
		Type:            types.StringValue(timeThreshold.Type),
		ViolationsCount: types.Int64Value(int64(timeThreshold.ViolationsCount)),
	}
}

// mapSyntheticTestIdsToState converts synthetic test IDs array to Terraform set
func (r *syntheticAlertConfigResource) mapSyntheticTestIdsToState(ctx context.Context, syntheticTestIds []string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, syntheticTestIds)
}

// mapAlertChannelIdsToState converts alert channel IDs array to Terraform set
func (r *syntheticAlertConfigResource) mapAlertChannelIdsToState(ctx context.Context, alertChannelIds []string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, alertChannelIds)
}

// mapCustomPayloadFieldsToState converts custom payload fields to Terraform list
func (r *syntheticAlertConfigResource) mapCustomPayloadFieldsToState(ctx context.Context, customPayloadFields []restapi.CustomPayloadField[any]) (types.List, diag.Diagnostics) {
	return shared.CustomPayloadFieldsToTerraform(ctx, customPayloadFields)
}

func mapTagFilterExpressionFromSchema(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *syntheticAlertConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: resourcehandle.CreateStateUpgraderForVersion(1),
	}
}
