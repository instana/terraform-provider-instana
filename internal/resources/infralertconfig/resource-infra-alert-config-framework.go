package infralertconfig

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewInfraAlertConfigResourceHandleFramework creates a new instance of the infrastructure alert configuration resource
func NewInfraAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.InfraAlertConfig] {
	return &infraAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaInfraAlertConfigFramework,
			Schema:        buildInfraAlertConfigSchema(),
			SchemaVersion: 0,
		},
	}
}

// buildInfraAlertConfigSchema constructs the Terraform schema for the infrastructure alert config resource
func buildInfraAlertConfigSchema() schema.Schema {
	return schema.Schema{
		Description: InfraAlertConfigDescResource,
		Attributes: map[string]schema.Attribute{
			InfraAlertConfigFieldID: schema.StringAttribute{
				Description: InfraAlertConfigDescID,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			InfraAlertConfigFieldName: schema.StringAttribute{
				Description: InfraAlertConfigDescName,
				Required:    true,
			},
			InfraAlertConfigFieldDescription: schema.StringAttribute{
				Description: InfraAlertConfigDescDescription,
				Optional:    true,
			},
			InfraAlertConfigFieldTagFilter: schema.StringAttribute{
				Description: InfraAlertConfigDescTagFilter,
				Optional:    true,
			},
			InfraAlertConfigFieldGroupBy: schema.ListAttribute{
				Description: InfraAlertConfigDescGroupBy,
				Optional:    true,
				ElementType: types.StringType,
			},
			InfraAlertConfigFieldGranularity: schema.Int64Attribute{
				Description: InfraAlertConfigDescGranularity,
				Required:    true,
			},
			InfraAlertConfigFieldEvaluationType: schema.StringAttribute{
				Description: InfraAlertConfigDescEvaluationType,
				Required:    true,
			},
			InfraAlertConfigFieldCustomPayloadField: shared.GetCustomPayloadFieldsSchema(),
			InfraAlertConfigFieldRules:              buildRulesSchema(),
			InfraAlertConfigFieldAlertChannels:      buildAlertChannelsSchema(),
			InfraAlertConfigFieldTimeThreshold:      buildTimeThresholdSchema(),
		},
	}
}

// buildRulesSchema constructs the schema for rules configuration
func buildRulesSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: InfraAlertConfigDescRules,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			InfraAlertConfigFieldGenericRule: buildGenericRuleSchema(),
		},
	}
}

// buildGenericRuleSchema constructs the schema for generic rule configuration
func buildGenericRuleSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: InfraAlertConfigDescGenericRule,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			InfraAlertConfigFieldMetricName: schema.StringAttribute{
				Description: InfraAlertConfigDescMetricName,
				Required:    true,
			},
			InfraAlertConfigFieldEntityType: schema.StringAttribute{
				Description: InfraAlertConfigDescEntityType,
				Required:    true,
			},
			InfraAlertConfigFieldAggregation: schema.StringAttribute{
				Description: InfraAlertConfigDescAggregation,
				Required:    true,
			},
			InfraAlertConfigFieldCrossSeriesAggregation: schema.StringAttribute{
				Description: InfraAlertConfigDescCrossSeriesAggregation,
				Required:    true,
			},
			InfraAlertConfigFieldRegex: schema.BoolAttribute{
				Description: InfraAlertConfigDescRegex,
				Required:    true,
			},
			InfraAlertConfigFieldThresholdOperator: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The operator to apply for threshold comparison",
				Validators: []validator.String{
					stringvalidator.OneOf(">", ">=", "<", "<="),
				},
			},
			InfraAlertConfigFieldThreshold: schema.SingleNestedAttribute{
				Description: InfraAlertConfigDescThreshold,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					ResourceFieldThresholdRuleWarningSeverity:  shared.StaticAndAdaptiveThresholdAttributeSchema(),
					ResourceFieldThresholdRuleCriticalSeverity: shared.StaticAndAdaptiveThresholdAttributeSchema(),
				},
			},
		},
	}
}

// buildAlertChannelsSchema constructs the schema for alert channels configuration
func buildAlertChannelsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: InfraAlertConfigDescAlertChannels,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ResourceFieldThresholdRuleWarningSeverity: schema.ListAttribute{
				Optional:    true,
				Description: InfraAlertConfigDescAlertChannelIDs,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			ResourceFieldThresholdRuleCriticalSeverity: schema.ListAttribute{
				Optional:    true,
				Description: InfraAlertConfigDescAlertChannelIDs,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

// buildTimeThresholdSchema constructs the schema for time threshold configuration
func buildTimeThresholdSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: InfraAlertConfigDescTimeThreshold,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			InfraAlertConfigFieldTimeThresholdViolationsInSequence: schema.SingleNestedAttribute{
				Description: InfraAlertConfigDescViolationsInSequence,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					InfraAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
						Required:    true,
						Description: InfraAlertConfigDescTimeWindow,
					},
				},
			},
		},
	}
}

type infraAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *infraAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *infraAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.InfraAlertConfig] {
	return api.InfraAlertConfig()
}

func (r *infraAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *infraAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, _ *tfsdk.Plan, resource *restapi.InfraAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	model, modelDiags := r.buildInfraAlertConfigModelFromAPIResponse(ctx, resource)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return diags
	}

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// buildInfraAlertConfigModelFromAPIResponse constructs an InfraAlertConfigModel from the API response
func (r *infraAlertConfigResourceFramework) buildInfraAlertConfigModelFromAPIResponse(ctx context.Context, resource *restapi.InfraAlertConfig) (InfraAlertConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := InfraAlertConfigModel{
		ID:             types.StringValue(resource.ID),
		Name:           types.StringValue(resource.Name),
		Description:    types.StringValue(resource.Description),
		Granularity:    types.Int64Value(int64(resource.Granularity)),
		EvaluationType: types.StringValue(string(resource.EvaluationType)),
	}

	tagFilter, tagFilterDiags := r.mapTagFilterToModel(resource.TagFilterExpression)
	diags.Append(tagFilterDiags...)
	model.TagFilter = tagFilter

	model.GroupBy = r.mapGroupByToModel(resource.GroupBy)

	alertChannels, alertChannelsDiags := r.mapAlertChannelsToModel(ctx, resource.AlertChannels)
	diags.Append(alertChannelsDiags...)
	model.AlertChannels = alertChannels

	model.TimeThreshold = r.mapTimeThresholdToModel(resource.TimeThreshold)

	customPayloadFields, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, resource.CustomerPayloadFields)
	diags.Append(payloadDiags...)
	model.CustomPayloadField = customPayloadFields

	rules, rulesDiags := r.mapRulesToModel(ctx, resource.Rules)
	diags.Append(rulesDiags...)
	model.Rules = rules

	return model, diags
}

// mapTagFilterToModel converts API tag filter to model representation
func (r *infraAlertConfigResourceFramework) mapTagFilterToModel(tagFilterExpression *restapi.TagFilter) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagFilterExpression == nil {
		return types.StringNull(), diags
	}

	tagFilterString, err := tagfilter.MapTagFilterToNormalizedString(tagFilterExpression)
	if err != nil {
		diags.AddError(
			InfraAlertConfigErrMappingTagFilter,
			fmt.Sprintf(InfraAlertConfigErrMappingTagFilterMsg, err),
		)
		return types.StringNull(), diags
	}

	return util.SetStringPointerToState(tagFilterString), diags
}

// mapGroupByToModel converts API group by slice to model representation
func (r *infraAlertConfigResourceFramework) mapGroupByToModel(groupBy []string) types.List {
	if len(groupBy) == 0 {
		return types.ListNull(types.StringType)
	}

	groupByElements := make([]attr.Value, len(groupBy))
	for i, groupByValue := range groupBy {
		groupByElements[i] = types.StringValue(groupByValue)
	}
	return types.ListValueMust(types.StringType, groupByElements)
}

// mapAlertChannelsToModel converts API alert channels to model representation
func (r *infraAlertConfigResourceFramework) mapAlertChannelsToModel(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string) (*InfraAlertChannelsModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Always return a model structure with null lists to maintain consistency
	// with Terraform's expected schema structure, even when the API returns
	// an empty map. This prevents "inconsistent result after apply" errors
	// when the plan has alert_channels defined but empty.
	if len(alertChannels) == 0 {
		return nil, diags
	}
	alertChannelsModel := &InfraAlertChannelsModel{
		Warning:  types.ListNull(types.StringType),
		Critical: types.ListNull(types.StringType),
	}

	// Only populate the lists if there are actual channels in the API response
	if len(alertChannels) > 0 {
		warningChannels, warningDiags := r.mapSeverityChannelsToModel(ctx, alertChannels, restapi.WarningSeverity)
		diags.Append(warningDiags...)
		alertChannelsModel.Warning = warningChannels

		criticalChannels, criticalDiags := r.mapSeverityChannelsToModel(ctx, alertChannels, restapi.CriticalSeverity)
		diags.Append(criticalDiags...)
		alertChannelsModel.Critical = criticalChannels
	}

	return alertChannelsModel, diags
}

// mapSeverityChannelsToModel converts alert channels for a specific severity to model representation
func (r *infraAlertConfigResourceFramework) mapSeverityChannelsToModel(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string, severity restapi.AlertSeverity) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	channels, exists := alertChannels[severity]
	if !exists || len(channels) == 0 {
		return types.ListNull(types.StringType), diags
	}

	channelList, listDiags := types.ListValueFrom(ctx, types.StringType, channels)
	diags.Append(listDiags...)
	return channelList, diags
}

// mapTimeThresholdToModel converts API time threshold to model representation
func (r *infraAlertConfigResourceFramework) mapTimeThresholdToModel(apiTimeThreshold *restapi.InfraTimeThreshold) *InfraTimeThresholdModel {
	if apiTimeThreshold == nil {
		return nil
	}

	if apiTimeThreshold.Type != TimeThresholdTypeViolationsInSequence {
		return nil
	}

	return &InfraTimeThresholdModel{
		ViolationsInSequence: &InfraViolationsInSequenceModel{
			TimeWindow: types.Int64Value(apiTimeThreshold.TimeWindow),
		},
	}
}

// mapRulesToModel converts API rules to model representation
func (r *infraAlertConfigResourceFramework) mapRulesToModel(ctx context.Context, rules []restapi.RuleWithThreshold[restapi.InfraAlertRule]) (*InfraRulesModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(rules) == 0 {
		return nil, diags
	}

	firstRule := rules[0]
	thresholdRuleModel := r.mapThresholdsToModel(ctx, firstRule.Thresholds)

	genericRuleModel := &InfraGenericRuleModel{
		MetricName:             types.StringValue(firstRule.Rule.MetricName),
		EntityType:             types.StringValue(firstRule.Rule.EntityType),
		Aggregation:            types.StringValue(string(firstRule.Rule.Aggregation)),
		CrossSeriesAggregation: types.StringValue(string(firstRule.Rule.CrossSeriesAggregation)),
		Regex:                  types.BoolValue(firstRule.Rule.Regex),
		ThresholdOperator:      types.StringValue(string(firstRule.ThresholdOperator)),
		ThresholdRule:          thresholdRuleModel,
	}

	return &InfraRulesModel{
		GenericRule: genericRuleModel,
	}, diags
}

// mapThresholdsToModel converts API thresholds to model representation
func (r *infraAlertConfigResourceFramework) mapThresholdsToModel(ctx context.Context, thresholds map[restapi.AlertSeverity]restapi.ThresholdRule) *shared.ThresholdPluginModel {
	thresholdRuleModel := &shared.ThresholdPluginModel{}

	if warningThreshold, hasWarning := thresholds[restapi.WarningSeverity]; hasWarning {
		thresholdRuleModel.Warning = shared.MapThresholdPluginToState(ctx, &warningThreshold, true)
	}

	if criticalThreshold, hasCritical := thresholds[restapi.CriticalSeverity]; hasCritical {
		thresholdRuleModel.Critical = shared.MapThresholdPluginToState(ctx, &criticalThreshold, true)
	}

	return thresholdRuleModel
}

// MapStateToDataObject maps Terraform state/plan to API InfraAlertConfig object
func (r *infraAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.InfraAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	infraAlertConfig := &restapi.InfraAlertConfig{
		ID:             r.extractConfigID(model),
		Name:           model.Name.ValueString(),
		Description:    model.Description.ValueString(),
		Granularity:    restapi.Granularity(model.Granularity.ValueInt64()),
		EvaluationType: restapi.InfraAlertEvaluationType(model.EvaluationType.ValueString()),
	}

	tagFilter, tagFilterDiags := r.mapModelTagFilterToAPI(model.TagFilter)
	diags.Append(tagFilterDiags...)
	infraAlertConfig.TagFilterExpression = tagFilter

	groupBy, groupByDiags := r.mapModelGroupByToAPI(ctx, model.GroupBy)
	diags.Append(groupByDiags...)
	infraAlertConfig.GroupBy = groupBy

	alertChannels, alertChannelsDiags := r.mapModelAlertChannelsToAPI(ctx, model.AlertChannels)
	diags.Append(alertChannelsDiags...)
	infraAlertConfig.AlertChannels = alertChannels

	infraAlertConfig.TimeThreshold = r.mapModelTimeThresholdToAPI(model.TimeThreshold)

	customPayloadFields, payloadDiags := r.mapModelCustomPayloadFieldsToAPI(ctx, model.CustomPayloadField)
	diags.Append(payloadDiags...)
	infraAlertConfig.CustomerPayloadFields = customPayloadFields

	rules, rulesDiags := r.mapModelRulesToAPI(ctx, model.Rules)
	diags.Append(rulesDiags...)
	infraAlertConfig.Rules = rules

	return infraAlertConfig, diags
}

// extractModelFromPlanOrState retrieves the InfraAlertConfigModel from plan or state
func (r *infraAlertConfigResourceFramework) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (InfraAlertConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model InfraAlertConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractConfigID extracts the configuration ID from the model
func (r *infraAlertConfigResourceFramework) extractConfigID(model InfraAlertConfigModel) string {
	if model.ID.IsNull() {
		return EmptyString
	}
	return model.ID.ValueString()
}

// mapModelTagFilterToAPI converts model tag filter to API representation
func (r *infraAlertConfigResourceFramework) mapModelTagFilterToAPI(tagFilter types.String) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagFilter.IsNull() || tagFilter.ValueString() == EmptyString {
		return nil, diags
	}

	tagFilterStr := tagFilter.ValueString()
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(tagFilterStr)
	if err != nil {
		diags.AddError(
			InfraAlertConfigErrParsingTagFilter,
			fmt.Sprintf(InfraAlertConfigErrParsingTagFilterMsg, err),
		)
		return nil, diags
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), diags
}

// mapModelGroupByToAPI converts model group by to API representation
func (r *infraAlertConfigResourceFramework) mapModelGroupByToAPI(ctx context.Context, groupBy types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if groupBy.IsNull() || groupBy.IsUnknown() {
		return nil, diags
	}

	var groupBySlice []string
	diags.Append(groupBy.ElementsAs(ctx, &groupBySlice, false)...)
	return groupBySlice, diags
}

// mapModelAlertChannelsToAPI converts model alert channels to API representation
func (r *infraAlertConfigResourceFramework) mapModelAlertChannelsToAPI(ctx context.Context, alertChannelsModel *InfraAlertChannelsModel) (map[restapi.AlertSeverity][]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannels := make(map[restapi.AlertSeverity][]string)

	if alertChannelsModel == nil {
		return alertChannels, diags
	}

	// When alert_channels block is defined, always include both severity levels
	// with empty arrays if no channels are specified, to maintain consistency
	warningChannels, warningDiags := r.extractChannelsForSeverity(ctx, alertChannelsModel.Warning)
	diags.Append(warningDiags...)
	if warningChannels == nil {
		warningChannels = []string{}
	}
	alertChannels[restapi.WarningSeverity] = warningChannels

	criticalChannels, criticalDiags := r.extractChannelsForSeverity(ctx, alertChannelsModel.Critical)
	diags.Append(criticalDiags...)
	if criticalChannels == nil {
		criticalChannels = []string{}
	}
	alertChannels[restapi.CriticalSeverity] = criticalChannels

	return alertChannels, diags
}

// extractChannelsForSeverity extracts channel IDs from a list for a specific severity
func (r *infraAlertConfigResourceFramework) extractChannelsForSeverity(ctx context.Context, channelList types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if channelList.IsNull() || channelList.IsUnknown() {
		return nil, diags
	}

	var channels []string
	diags.Append(channelList.ElementsAs(ctx, &channels, false)...)
	return channels, diags
}

// mapModelTimeThresholdToAPI converts model time threshold to API representation
func (r *infraAlertConfigResourceFramework) mapModelTimeThresholdToAPI(timeThresholdModel *InfraTimeThresholdModel) *restapi.InfraTimeThreshold {
	if timeThresholdModel == nil || timeThresholdModel.ViolationsInSequence == nil {
		return nil
	}

	violationsModel := timeThresholdModel.ViolationsInSequence
	if violationsModel.TimeWindow.IsNull() || violationsModel.TimeWindow.IsUnknown() {
		return nil
	}

	return &restapi.InfraTimeThreshold{
		Type:       TimeThresholdTypeViolationsInSequence,
		TimeWindow: violationsModel.TimeWindow.ValueInt64(),
	}
}

// mapModelCustomPayloadFieldsToAPI converts model custom payload fields to API representation
func (r *infraAlertConfigResourceFramework) mapModelCustomPayloadFieldsToAPI(ctx context.Context, customPayloadField types.List) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics

	if customPayloadField.IsNull() || customPayloadField.IsUnknown() {
		return nil, diags
	}

	var customPayloadFieldModels []InfraCustomPayloadFieldModel
	diags.Append(customPayloadField.ElementsAs(ctx, &customPayloadFieldModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	customerPayloadFields := make([]restapi.CustomPayloadField[any], 0, len(customPayloadFieldModels))
	for _, field := range customPayloadFieldModels {
		customerPayloadFields = append(customerPayloadFields, restapi.CustomPayloadField[any]{
			Key:   field.Key.ValueString(),
			Value: field.Value.ValueString(),
		})
	}

	return customerPayloadFields, diags
}

// mapModelRulesToAPI converts model rules to API representation
func (r *infraAlertConfigResourceFramework) mapModelRulesToAPI(ctx context.Context, rulesModel *InfraRulesModel) ([]restapi.RuleWithThreshold[restapi.InfraAlertRule], diag.Diagnostics) {
	var diags diag.Diagnostics

	if rulesModel == nil || rulesModel.GenericRule == nil {
		return nil, diags
	}

	genericRuleModel := rulesModel.GenericRule

	ruleWithThreshold := restapi.RuleWithThreshold[restapi.InfraAlertRule]{
		ThresholdOperator: restapi.ThresholdOperator(genericRuleModel.ThresholdOperator.ValueString()),
		Rule: restapi.InfraAlertRule{
			AlertType:              GenericRuleAlertType,
			MetricName:             genericRuleModel.MetricName.ValueString(),
			EntityType:             genericRuleModel.EntityType.ValueString(),
			Aggregation:            restapi.Aggregation(genericRuleModel.Aggregation.ValueString()),
			CrossSeriesAggregation: restapi.Aggregation(genericRuleModel.CrossSeriesAggregation.ValueString()),
			Regex:                  genericRuleModel.Regex.ValueBool(),
		},
		Thresholds: make(map[restapi.AlertSeverity]restapi.ThresholdRule),
	}

	thresholds, thresholdDiags := r.mapModelThresholdsToAPI(ctx, genericRuleModel.ThresholdRule)
	diags.Append(thresholdDiags...)
	ruleWithThreshold.Thresholds = thresholds

	return []restapi.RuleWithThreshold[restapi.InfraAlertRule]{ruleWithThreshold}, diags
}

// mapModelThresholdsToAPI converts model thresholds to API representation
func (r *infraAlertConfigResourceFramework) mapModelThresholdsToAPI(ctx context.Context, thresholdRuleModel *shared.ThresholdPluginModel) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholds := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	if thresholdRuleModel == nil {
		return thresholds, diags
	}

	if thresholdRuleModel.Warning != nil {
		warningThreshold, warningDiags := shared.MapThresholdRulePluginFromState(ctx, thresholdRuleModel.Warning)
		diags.Append(warningDiags...)
		if warningThreshold != nil {
			thresholds[restapi.WarningSeverity] = *warningThreshold
		}
	}

	if thresholdRuleModel.Critical != nil {
		criticalThreshold, criticalDiags := shared.MapThresholdRulePluginFromState(ctx, thresholdRuleModel.Critical)
		diags.Append(criticalDiags...)
		if criticalThreshold != nil {
			thresholds[restapi.CriticalSeverity] = *criticalThreshold
		}
	}

	return thresholds, diags
}

// Made with Bob
