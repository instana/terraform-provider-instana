package syntheticalertconfig

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewSyntheticAlertConfigResourceHandleFramework creates the resource handle for Synthetic Alert Configuration
func NewSyntheticAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.SyntheticAlertConfig] {
	resource := &syntheticAlertConfigResourceFramework{}
	return resource.initialize()
}

// initialize sets up the resource with metadata and schema
func (r *syntheticAlertConfigResourceFramework) initialize() *syntheticAlertConfigResourceFramework {
	r.metaData = resourcehandle.ResourceMetaDataFramework{
		ResourceName:  ResourceInstanaSyntheticAlertConfigFramework,
		Schema:        r.buildSchema(),
		SchemaVersion: 1,
	}
	return r
}

// buildSchema constructs the complete schema for the resource
func (r *syntheticAlertConfigResourceFramework) buildSchema() schema.Schema {
	return schema.Schema{
		Description: SyntheticAlertConfigDescResource,
		Attributes:  r.buildSchemaAttributes(),
	}
}

// buildSchemaAttributes constructs the top-level schema attributes
func (r *syntheticAlertConfigResourceFramework) buildSchemaAttributes() map[string]schema.Attribute {
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
func (r *syntheticAlertConfigResourceFramework) buildIDAttribute() schema.Attribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: SyntheticAlertConfigDescID,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// buildNameAttribute creates the name attribute schema
func (r *syntheticAlertConfigResourceFramework) buildNameAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescName,
		Validators:  r.buildNameValidators(),
	}
}

// buildNameValidators creates validators for the name field
func (r *syntheticAlertConfigResourceFramework) buildNameValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthBetween(1, 256),
	}
}

// buildDescriptionAttribute creates the description attribute schema
func (r *syntheticAlertConfigResourceFramework) buildDescriptionAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescDescription,
		Validators:  r.buildDescriptionValidators(),
	}
}

// buildDescriptionValidators creates validators for the description field
func (r *syntheticAlertConfigResourceFramework) buildDescriptionValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthBetween(0, 1024),
	}
}

// buildSyntheticTestIdsAttribute creates the synthetic_test_ids attribute schema
func (r *syntheticAlertConfigResourceFramework) buildSyntheticTestIdsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescSyntheticTestIds,
		ElementType: types.StringType,
	}
}

// buildSeverityAttribute creates the severity attribute schema
func (r *syntheticAlertConfigResourceFramework) buildSeverityAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescSeverity,
		Validators:  r.buildSeverityValidators(),
	}
}

// buildSeverityValidators creates validators for the severity field
func (r *syntheticAlertConfigResourceFramework) buildSeverityValidators() []validator.Int64 {
	return []validator.Int64{
		int64validator.OneOf(5, 10),
	}
}

// buildTagFilterAttribute creates the tag_filter attribute schema
func (r *syntheticAlertConfigResourceFramework) buildTagFilterAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescTagFilter,
	}
}

// buildAlertChannelIdsAttribute creates the alert_channel_ids attribute schema
func (r *syntheticAlertConfigResourceFramework) buildAlertChannelIdsAttribute() schema.Attribute {
	return schema.SetAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescAlertChannelIds,
		ElementType: types.StringType,
	}
}

// buildGracePeriodAttribute creates the grace_period attribute schema
func (r *syntheticAlertConfigResourceFramework) buildGracePeriodAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescGracePeriod,
	}
}

// buildRuleAttribute creates the rule nested attribute schema
func (r *syntheticAlertConfigResourceFramework) buildRuleAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: SyntheticAlertConfigDescRule,
		Required:    true,
		Attributes:  r.buildRuleNestedAttributes(),
	}
}

// buildRuleNestedAttributes constructs the nested rule attributes
func (r *syntheticAlertConfigResourceFramework) buildRuleNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticAlertRuleFieldAlertType:   r.buildRuleAlertTypeAttribute(),
		SyntheticAlertRuleFieldMetricName:  r.buildRuleMetricNameAttribute(),
		SyntheticAlertRuleFieldAggregation: r.buildRuleAggregationAttribute(),
	}
}

// buildRuleAlertTypeAttribute creates the rule alert_type attribute schema
func (r *syntheticAlertConfigResourceFramework) buildRuleAlertTypeAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescRuleAlertType,
		Validators:  r.buildRuleAlertTypeValidators(),
	}
}

// buildRuleAlertTypeValidators creates validators for the rule alert_type field
func (r *syntheticAlertConfigResourceFramework) buildRuleAlertTypeValidators() []validator.String {
	return []validator.String{
		stringvalidator.OneOf(SyntheticAlertConfigValidAlertType),
	}
}

// buildRuleMetricNameAttribute creates the rule metric_name attribute schema
func (r *syntheticAlertConfigResourceFramework) buildRuleMetricNameAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescRuleMetricName,
		Validators:  r.buildRuleMetricNameValidators(),
	}
}

// buildRuleMetricNameValidators creates validators for the rule metric_name field
func (r *syntheticAlertConfigResourceFramework) buildRuleMetricNameValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthBetween(1, 256),
	}
}

// buildRuleAggregationAttribute creates the rule aggregation attribute schema
func (r *syntheticAlertConfigResourceFramework) buildRuleAggregationAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: SyntheticAlertConfigDescRuleAggregation,
		Validators:  r.buildRuleAggregationValidators(),
	}
}

// buildRuleAggregationValidators creates validators for the rule aggregation field
func (r *syntheticAlertConfigResourceFramework) buildRuleAggregationValidators() []validator.String {
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
func (r *syntheticAlertConfigResourceFramework) buildTimeThresholdAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: SyntheticAlertConfigDescTimeThreshold,
		Required:    true,
		Attributes:  r.buildTimeThresholdNestedAttributes(),
	}
}

// buildTimeThresholdNestedAttributes constructs the nested time threshold attributes
func (r *syntheticAlertConfigResourceFramework) buildTimeThresholdNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticAlertTimeThresholdFieldType:            r.buildTimeThresholdTypeAttribute(),
		SyntheticAlertTimeThresholdFieldViolationsCount: r.buildTimeThresholdViolationsCountAttribute(),
	}
}

// buildTimeThresholdTypeAttribute creates the time threshold type attribute schema
func (r *syntheticAlertConfigResourceFramework) buildTimeThresholdTypeAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SyntheticAlertConfigDescTimeThresholdType,
		Validators:  r.buildTimeThresholdTypeValidators(),
	}
}

// buildTimeThresholdTypeValidators creates validators for the time threshold type field
func (r *syntheticAlertConfigResourceFramework) buildTimeThresholdTypeValidators() []validator.String {
	return []validator.String{
		stringvalidator.OneOf(SyntheticAlertConfigValidTimeThresholdType),
	}
}

// buildTimeThresholdViolationsCountAttribute creates the time threshold violations_count attribute schema
func (r *syntheticAlertConfigResourceFramework) buildTimeThresholdViolationsCountAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Required:    true,
		Description: SyntheticAlertConfigDescTimeThresholdViolationsCount,
		Validators:  r.buildTimeThresholdViolationsCountValidators(),
	}
}

// buildTimeThresholdViolationsCountValidators creates validators for the time threshold violations_count field
func (r *syntheticAlertConfigResourceFramework) buildTimeThresholdViolationsCountValidators() []validator.Int64 {
	return []validator.Int64{
		int64validator.Between(1, 12),
	}
}

type syntheticAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *syntheticAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *syntheticAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SyntheticAlertConfig] {
	return api.SyntheticAlertConfigs()
}

func (r *syntheticAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *syntheticAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SyntheticAlertConfig, diag.Diagnostics) {
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
func (r *syntheticAlertConfigResourceFramework) extractModelFromState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*SyntheticAlertConfigModel, diag.Diagnostics) {
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
func (r *syntheticAlertConfigResourceFramework) mapRuleFromModel(model *SyntheticAlertConfigModel) (restapi.SyntheticAlertRule, diag.Diagnostics) {
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
func (r *syntheticAlertConfigResourceFramework) mapTimeThresholdFromModel(model *SyntheticAlertConfigModel) (restapi.SyntheticAlertTimeThreshold, diag.Diagnostics) {
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
func (r *syntheticAlertConfigResourceFramework) extractSyntheticTestIdsFromModel(ctx context.Context, model *SyntheticAlertConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	syntheticTestIds := []string{}

	if !model.SyntheticTestIds.IsNull() {
		diags.Append(model.SyntheticTestIds.ElementsAs(ctx, &syntheticTestIds, false)...)
	}

	return syntheticTestIds, diags
}

// extractAlertChannelIdsFromModel extracts alert channel IDs from the model
func (r *syntheticAlertConfigResourceFramework) extractAlertChannelIdsFromModel(ctx context.Context, model *SyntheticAlertConfigModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var alertChannelIds []string

	if !model.AlertChannelIds.IsNull() {
		diags.Append(model.AlertChannelIds.ElementsAs(ctx, &alertChannelIds, false)...)
	}

	return alertChannelIds, diags
}

// mapTagFilterFromModel converts tag filter model to API tag filter object
func (r *syntheticAlertConfigResourceFramework) mapTagFilterFromModel(model *SyntheticAlertConfigModel) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if r.hasTagFilterValue(model) {
		return r.parseTagFilterExpression(model.TagFilter.ValueString(), &diags)
	}

	return r.createDefaultTagFilter(), diags
}

// hasTagFilterValue checks if the model has a tag filter value
func (r *syntheticAlertConfigResourceFramework) hasTagFilterValue(model *SyntheticAlertConfigModel) bool {
	return !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown()
}

// parseTagFilterExpression parses the tag filter expression string
func (r *syntheticAlertConfigResourceFramework) parseTagFilterExpression(expression string, diags *diag.Diagnostics) (*restapi.TagFilter, diag.Diagnostics) {
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
func (r *syntheticAlertConfigResourceFramework) createDefaultTagFilter() *restapi.TagFilter {
	operator := restapi.LogicalOperatorType(TagFilterLogicalOperatorAnd)
	return &restapi.TagFilter{
		Type:            TagFilterTypeExpression,
		LogicalOperator: &operator,
		Elements:        []*restapi.TagFilter{},
	}
}

// extractCustomPayloadFieldsFromModel extracts custom payload fields from the model
func (r *syntheticAlertConfigResourceFramework) extractCustomPayloadFieldsFromModel(ctx context.Context, model *SyntheticAlertConfigModel) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics
	var customPayloadFields []restapi.CustomPayloadField[any]

	if !model.CustomPayloadFields.IsNull() {
		customPayloadFields, diags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
	}

	return customPayloadFields, diags
}

// buildAPIObjectFromModel constructs the API object from model and extracted data
func (r *syntheticAlertConfigResourceFramework) buildAPIObjectFromModel(
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
func (r *syntheticAlertConfigResourceFramework) setOptionalFieldsOnAPIObject(apiObject *restapi.SyntheticAlertConfig, model *SyntheticAlertConfigModel) {
	if !model.Severity.IsNull() {
		apiObject.Severity = int(model.Severity.ValueInt64())
	}

	if !model.GracePeriod.IsNull() {
		apiObject.GracePeriod = model.GracePeriod.ValueInt64()
	}
}

func (r *syntheticAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SyntheticAlertConfig) diag.Diagnostics {
	model := r.buildModelFromAPIObject(apiObject)
	var diags diag.Diagnostics

	gracePeriodValue := r.mapGracePeriodToState(apiObject.GracePeriod)
	model.GracePeriod = gracePeriodValue

	tagFilterValue, tagFilterDiags := r.mapTagFilterToState(apiObject.TagFilterExpression)
	diags.Append(tagFilterDiags...)
	if diags.HasError() {
		return diags
	}
	model.TagFilter = tagFilterValue

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

// buildModelFromAPIObject creates a model with basic fields from API object
func (r *syntheticAlertConfigResourceFramework) buildModelFromAPIObject(apiObject *restapi.SyntheticAlertConfig) SyntheticAlertConfigModel {
	return SyntheticAlertConfigModel{
		ID:          types.StringValue(apiObject.ID),
		Name:        types.StringValue(apiObject.Name),
		Description: types.StringValue(apiObject.Description),
		Severity:    types.Int64Value(int64(apiObject.Severity)),
	}
}

// mapGracePeriodToState converts grace period to state value
func (r *syntheticAlertConfigResourceFramework) mapGracePeriodToState(gracePeriod int64) types.Int64 {
	if gracePeriod > 0 {
		return types.Int64Value(gracePeriod)
	}
	return types.Int64Null()
}

// mapTagFilterToState converts tag filter expression to state value
func (r *syntheticAlertConfigResourceFramework) mapTagFilterToState(tagFilterExpression *restapi.TagFilter) (types.String, diag.Diagnostics) {
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
func (r *syntheticAlertConfigResourceFramework) mapRuleToModel(rule restapi.SyntheticAlertRule) *SyntheticAlertRuleModel {
	aggregationValue := r.buildAggregationValue(rule.Aggregation)

	return &SyntheticAlertRuleModel{
		AlertType:   types.StringValue(rule.AlertType),
		MetricName:  types.StringValue(rule.MetricName),
		Aggregation: aggregationValue,
	}
}

// buildAggregationValue creates a types.String value for aggregation
func (r *syntheticAlertConfigResourceFramework) buildAggregationValue(aggregation string) types.String {
	if aggregation == "" {
		return types.StringNull()
	}
	return types.StringValue(aggregation)
}

// mapTimeThresholdToModel converts API time threshold to model time threshold
func (r *syntheticAlertConfigResourceFramework) mapTimeThresholdToModel(timeThreshold restapi.SyntheticAlertTimeThreshold) *SyntheticAlertTimeThresholdModel {
	return &SyntheticAlertTimeThresholdModel{
		Type:            types.StringValue(timeThreshold.Type),
		ViolationsCount: types.Int64Value(int64(timeThreshold.ViolationsCount)),
	}
}

// mapSyntheticTestIdsToState converts synthetic test IDs array to Terraform set
func (r *syntheticAlertConfigResourceFramework) mapSyntheticTestIdsToState(ctx context.Context, syntheticTestIds []string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, syntheticTestIds)
}

// mapAlertChannelIdsToState converts alert channel IDs array to Terraform set
func (r *syntheticAlertConfigResourceFramework) mapAlertChannelIdsToState(ctx context.Context, alertChannelIds []string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, alertChannelIds)
}

// mapCustomPayloadFieldsToState converts custom payload fields to Terraform list
func (r *syntheticAlertConfigResourceFramework) mapCustomPayloadFieldsToState(ctx context.Context, customPayloadFields []restapi.CustomPayloadField[any]) (types.List, diag.Diagnostics) {
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
