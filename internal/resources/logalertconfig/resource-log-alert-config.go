package logalertconfig

import (
	"context"
	"math"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewLogAlertConfigResourceHandle creates the resource handle for Log Alert Configuration
func NewLogAlertConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.LogAlertConfig] {
	return &logAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaLogAlertConfig,
			Schema:        buildLogAlertConfigSchema(),
			SchemaVersion: 1,
		},
	}
}

// buildLogAlertConfigSchema constructs the Terraform schema for the log alert config resource
func buildLogAlertConfigSchema() schema.Schema {
	return schema.Schema{
		Description: LogAlertConfigDescResource,
		Attributes: map[string]schema.Attribute{
			LogAlertConfigFieldID: schema.StringAttribute{
				Computed:    true,
				Description: LogAlertConfigDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			LogAlertConfigFieldName: schema.StringAttribute{
				Required:    true,
				Description: LogAlertConfigDescName,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 256),
				},
			},
			LogAlertConfigFieldDescription: schema.StringAttribute{
				Required:    true,
				Description: LogAlertConfigDescDescription,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 65536),
				},
			},
			LogAlertConfigFieldGracePeriod: schema.Int64Attribute{
				Optional:    true,
				Description: LogAlertConfigDescGracePeriod,
			},
			LogAlertConfigFieldGranularity: schema.Int64Attribute{
				Optional:    true,
				Description: LogAlertConfigDescGranularity,
				Validators: []validator.Int64{
					int64validator.OneOf(
						int64(restapi.Granularity60000),
						int64(restapi.Granularity300000),
						int64(restapi.Granularity600000),
						int64(restapi.Granularity900000),
						int64(restapi.Granularity1200000),
						int64(restapi.Granularity1800000),
					),
				},
			},
			LogAlertConfigFieldTagFilter: schema.StringAttribute{
				Optional:    true,
				Description: LogAlertConfigDescTagFilter,
			},
			shared.DefaultCustomPayloadFieldsName: shared.GetCustomPayloadFieldsSchema(),
			LogAlertConfigFieldAlertChannels:      buildAlertChannelsSchema(),
			LogAlertConfigFieldGroupBy:            buildGroupBySchema(),
			LogAlertConfigFieldRules:              buildRulesSchema(),
			LogAlertConfigFieldTimeThreshold:      buildTimeThresholdSchema(),
		},
	}
}

// buildAlertChannelsSchema constructs the schema for alert channels configuration
func buildAlertChannelsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: LogAlertConfigDescAlertChannels,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			shared.ThresholdFieldWarning: schema.ListAttribute{
				Optional:    true,
				Description: LogAlertConfigDescAlertChannelIDs,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			shared.ThresholdFieldCritical: schema.ListAttribute{
				Optional:    true,
				Description: LogAlertConfigDescAlertChannelIDs,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

// buildGroupBySchema constructs the schema for group by configuration
func buildGroupBySchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: LogAlertConfigDescGroupBy,
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				LogAlertConfigFieldGroupByTagName: schema.StringAttribute{
					Required:    true,
					Description: LogAlertConfigDescGroupByTagName,
				},
				LogAlertConfigFieldGroupByKey: schema.StringAttribute{
					Optional:    true,
					Description: LogAlertConfigDescGroupByKey,
				},
			},
		},
	}
}

// buildRulesSchema constructs the schema for rules configuration
func buildRulesSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: LogAlertConfigDescRules,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			LogAlertConfigFieldMetricName: schema.StringAttribute{
				Required:    true,
				Description: LogAlertConfigDescMetricName,
			},
			LogAlertConfigFieldAlertType: schema.StringAttribute{
				Optional:    true,
				Description: LogAlertConfigDescAlertType,
				Validators: []validator.String{
					stringvalidator.OneOf(LogAlertTypeLogCount),
				},
			},
			LogAlertConfigFieldAggregation: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: LogAlertConfigDescAggregation,
				Validators: []validator.String{
					stringvalidator.OneOf(string(restapi.SumAggregation)),
				},
			},
			LogAlertConfigFieldThresholdOperator: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The operator to apply for threshold comparison",
				Validators: []validator.String{
					stringvalidator.OneOf(">", ">=", "<", "<="),
				},
			},
			LogAlertConfigFieldThreshold: schema.SingleNestedAttribute{
				Description: LogAlertConfigDescThreshold,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					LogAlertConfigFieldWarning:  shared.StaticAndAdaptiveThresholdAttributeSchema(),
					LogAlertConfigFieldCritical: shared.StaticAndAdaptiveThresholdAttributeSchema(),
				},
			},
		},
	}
}

// buildTimeThresholdSchema constructs the schema for time threshold configuration
func buildTimeThresholdSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: LogAlertConfigDescTimeThreshold,
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			LogAlertConfigFieldTimeThresholdViolationsInSequence: schema.SingleNestedAttribute{
				Description: LogAlertConfigDescViolationsInSequence,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					LogAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
						Description: LogAlertConfigDescTimeWindow,
						Required:    true,
					},
				},
			},
		},
	}
}

type logAlertConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *logAlertConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *logAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.LogAlertConfig] {
	return api.LogAlertConfig()
}

func (r *logAlertConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *logAlertConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.LogAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	var model LogAlertConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	model.ID = types.StringValue(config.ID)
	model.Name = types.StringValue(config.Name)
	model.Description = types.StringValue(config.Description)
	model.Granularity = types.Int64Value(int64(config.Granularity))
	model.GracePeriod = r.mapGracePeriodToModel(config.GracePeriod)

	// to preserve the existing value in plan/state to handle the value drift
	if model.TagFilter.IsNull() || model.TagFilter.IsUnknown() {
		tagFilter, tagFilterDiags := r.mapTagFilterToModel(config.TagFilterExpression)
		diags.Append(tagFilterDiags...)
		model.TagFilter = tagFilter
	}

	groupBy, groupByDiags := r.mapGroupByToModel(config.GroupBy)
	diags.Append(groupByDiags...)
	model.GroupBy = groupBy

	alertChannels, alertChannelsDiags := r.mapAlertChannelsToModel(ctx, config.AlertChannels)
	diags.Append(alertChannelsDiags...)
	model.AlertChannels = alertChannels

	model.TimeThreshold = r.mapTimeThresholdToModel(config.TimeThreshold)

	rules, rulesDiags := r.mapRulesToModel(config.Rules)
	diags.Append(rulesDiags...)
	model.Rules = rules

	customPayloadFields, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, config.CustomerPayloadFields)
	diags.Append(payloadDiags...)
	model.CustomPayloadFields = customPayloadFields

	if diags.HasError() {
		return diags
	}

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapGracePeriodToModel converts grace period to model representation
func (r *logAlertConfigResource) mapGracePeriodToModel(gracePeriod int64) types.Int64 {
	if gracePeriod > 0 {
		return types.Int64Value(gracePeriod)
	}
	return types.Int64Null()
}

// mapTagFilterToModel converts API tag filter to model representation
func (r *logAlertConfigResource) mapTagFilterToModel(tagFilterExpression *restapi.TagFilter) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagFilterExpression == nil {
		return types.StringNull(), diags
	}

	normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(tagFilterExpression)
	if err != nil {
		diags.AddError(
			LogAlertConfigErrNormalizingTagFilter,
			LogAlertConfigErrNormalizingTagFilterMsg+err.Error(),
		)
		return types.StringNull(), diags
	}

	return util.SetStringPointerToState(normalizedTagFilterString), diags
}

// mapGroupByToModel converts API group by to model representation
func (r *logAlertConfigResource) mapGroupByToModel(groupBy []restapi.GroupByTag) ([]GroupByModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(groupBy) == 0 {
		return []GroupByModel{}, diags
	}

	result := make([]GroupByModel, len(groupBy))
	for i, g := range groupBy {
		result[i] = GroupByModel{
			TagName: types.StringValue(g.TagName),
			Key:     r.mapGroupByKeyToModel(g.Key),
		}
	}

	return result, diags
}

// mapGroupByKeyToModel converts group by key to model representation
func (r *logAlertConfigResource) mapGroupByKeyToModel(key string) types.String {
	if key != EmptyString {
		return types.StringValue(key)
	}
	return types.StringNull()
}

// mapAlertChannelsToModel converts API alert channels to model representation
func (r *logAlertConfigResource) mapAlertChannelsToModel(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string) (*AlertChannelsModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(alertChannels) == 0 {
		return nil, diags
	}

	model := &AlertChannelsModel{}

	warningChannels, warningDiags := r.mapSeverityChannelsToModel(ctx, alertChannels, restapi.WarningSeverity)
	diags.Append(warningDiags...)
	model.Warning = warningChannels

	criticalChannels, criticalDiags := r.mapSeverityChannelsToModel(ctx, alertChannels, restapi.CriticalSeverity)
	diags.Append(criticalDiags...)
	model.Critical = criticalChannels

	return model, diags
}

// mapSeverityChannelsToModel converts alert channels for a specific severity to model representation
func (r *logAlertConfigResource) mapSeverityChannelsToModel(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string, severity restapi.AlertSeverity) (types.List, diag.Diagnostics) {
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
func (r *logAlertConfigResource) mapTimeThresholdToModel(apiTimeThreshold *restapi.LogTimeThreshold) *TimeThresholdModel {
	if apiTimeThreshold == nil {
		return nil
	}

	if apiTimeThreshold.Type != TimeThresholdTypeViolationsInSequence {
		return nil
	}

	return &TimeThresholdModel{
		ViolationsInSequence: &ViolationsInSequenceModel{
			TimeWindow: types.Int64Value(apiTimeThreshold.TimeWindow),
		},
	}
}

// mapRulesToModel converts API rules to model representation
func (r *logAlertConfigResource) mapRulesToModel(rules []restapi.RuleWithThreshold[restapi.LogAlertRule]) (*LogAlertRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(rules) == 0 {
		return nil, diags
	}

	ruleWithThreshold := rules[0]
	rule := ruleWithThreshold.Rule

	ruleModel := &LogAlertRuleModel{
		MetricName:        types.StringValue(rule.MetricName),
		AlertType:         types.StringValue(r.convertAlertTypeToSchema(rule.AlertType)),
		Aggregation:       r.mapAggregationToModel(rule.Aggregation),
		ThresholdOperator: types.StringValue(string(ruleWithThreshold.ThresholdOperator)),
		Threshold:         r.mapThresholdsToModel(ruleWithThreshold.Thresholds),
	}

	return ruleModel, diags
}

// convertAlertTypeToSchema converts API alert type to schema format
func (r *logAlertConfigResource) convertAlertTypeToSchema(alertType string) string {
	if alertType == LogAlertTypeLogCountAPI {
		return LogAlertTypeLogCount
	}
	return alertType
}

// mapAggregationToModel converts aggregation to model representation
func (r *logAlertConfigResource) mapAggregationToModel(aggregation restapi.Aggregation) types.String {
	if aggregation != EmptyString {
		return types.StringValue(string(aggregation))
	}
	return types.StringValue(string(restapi.SumAggregation))
}

// mapThresholdsToModel converts API thresholds to model representation
func (r *logAlertConfigResource) mapThresholdsToModel(thresholds map[restapi.AlertSeverity]restapi.ThresholdRule) *ThresholdModel {
	thresholdModel := &ThresholdModel{}

	if warningThreshold, hasWarning := thresholds[restapi.WarningSeverity]; hasWarning {
		thresholdModel.Warning = r.createThresholdModel(warningThreshold)
	}

	if criticalThreshold, hasCritical := thresholds[restapi.CriticalSeverity]; hasCritical {
		thresholdModel.Critical = r.createThresholdModel(criticalThreshold)
	}

	return thresholdModel
}

// createThresholdModel creates a threshold model from API threshold rule
func (r *logAlertConfigResource) createThresholdModel(threshold restapi.ThresholdRule) *shared.ThresholdTypeModel {
	// Handle adaptive baseline
	if threshold.Type == "adaptiveBaseline" {
		return &shared.ThresholdTypeModel{
			AdaptiveBaseline: &shared.AdaptiveBaselineModel{
				DeviationFactor: util.SetFloat32PointerToState(threshold.DeviationFactor),
				Adaptability:    util.SetFloat32PointerToState(threshold.Adaptability),
				Seasonality:     types.StringValue(string(*threshold.Seasonality)),
			},
		}
	}

	// Default to static threshold
	if threshold.Value == nil {
		return nil
	}
	return &shared.ThresholdTypeModel{
		Static: &shared.StaticTypeModel{
			Value: types.Float32Value(float32(*threshold.Value)),
		},
	}
}

// MapStateToDataObject maps Terraform state/plan to API LogAlertConfig object
func (r *logAlertConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.LogAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	config := &restapi.LogAlertConfig{
		ID:          r.extractConfigID(model),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Granularity: r.extractGranularity(model.Granularity),
		GracePeriod: r.extractGracePeriod(model.GracePeriod),
	}

	tagFilter, tagFilterDiags := r.mapModelTagFilterToAPI(model.TagFilter)
	diags.Append(tagFilterDiags...)
	config.TagFilterExpression = tagFilter

	groupBy, groupByDiags := r.mapModelGroupByToAPI(ctx, model.GroupBy)
	diags.Append(groupByDiags...)
	config.GroupBy = groupBy

	alertChannels, alertChannelsDiags := r.mapModelAlertChannelsToAPI(ctx, model.AlertChannels)
	diags.Append(alertChannelsDiags...)
	config.AlertChannels = alertChannels

	config.TimeThreshold = r.mapModelTimeThresholdToAPI(model.TimeThreshold)

	rules, rulesDiags := r.mapModelRulesToAPI(model.Rules)
	diags.Append(rulesDiags...)
	config.Rules = rules

	customPayloadFields, payloadDiags := shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
	diags.Append(payloadDiags...)
	config.CustomerPayloadFields = customPayloadFields

	return config, diags
}

// extractModelFromPlanOrState retrieves the LogAlertConfigModel from plan or state
func (r *logAlertConfigResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (LogAlertConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model LogAlertConfigModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractConfigID extracts the configuration ID from the model
func (r *logAlertConfigResource) extractConfigID(model LogAlertConfigModel) string {
	if model.ID.IsNull() {
		return EmptyString
	}
	return model.ID.ValueString()
}

// extractGranularity extracts granularity from the model with default value
func (r *logAlertConfigResource) extractGranularity(granularity types.Int64) restapi.Granularity {
	if !granularity.IsNull() {
		return restapi.Granularity(granularity.ValueInt64())
	}
	return restapi.Granularity600000
}

// extractGracePeriod extracts grace period from the model
func (r *logAlertConfigResource) extractGracePeriod(gracePeriod types.Int64) int64 {
	if !gracePeriod.IsNull() {
		return gracePeriod.ValueInt64()
	}
	return 0
}

// mapModelTagFilterToAPI converts model tag filter to API representation
func (r *logAlertConfigResource) mapModelTagFilterToAPI(tagFilter types.String) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagFilter.IsNull() {
		return nil, diags
	}

	parser := tagfilter.NewParser()
	expr, err := parser.Parse(tagFilter.ValueString())
	if err != nil {
		diags.AddError(
			LogAlertConfigErrParsingTagFilter,
			LogAlertConfigErrParsingTagFilterMsg+err.Error(),
		)
		return nil, diags
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), diags
}

// mapModelGroupByToAPI converts model group by to API representation
func (r *logAlertConfigResource) mapModelGroupByToAPI(ctx context.Context, groupByModels []GroupByModel) ([]restapi.GroupByTag, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(groupByModels) == 0 {
		return []restapi.GroupByTag{}, diags
	}

	result := make([]restapi.GroupByTag, len(groupByModels))
	for i, model := range groupByModels {
		result[i] = restapi.GroupByTag{
			TagName: model.TagName.ValueString(),
			Key:     r.extractGroupByKey(model.Key),
		}
	}

	return result, diags
}

// extractGroupByKey extracts the group by key from the model
func (r *logAlertConfigResource) extractGroupByKey(key types.String) string {
	if !key.IsNull() {
		return key.ValueString()
	}
	return EmptyString
}

// mapModelAlertChannelsToAPI converts model alert channels to API representation
func (r *logAlertConfigResource) mapModelAlertChannelsToAPI(ctx context.Context, alertChannelsModel *AlertChannelsModel) (map[restapi.AlertSeverity][]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannels := make(map[restapi.AlertSeverity][]string)

	if alertChannelsModel == nil {
		return alertChannels, diags
	}

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
func (r *logAlertConfigResource) extractChannelsForSeverity(ctx context.Context, channelList types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if channelList.IsNull() || channelList.IsUnknown() {
		return nil, diags
	}

	var channels []string
	diags.Append(channelList.ElementsAs(ctx, &channels, false)...)
	return channels, diags
}

// mapModelTimeThresholdToAPI converts model time threshold to API representation
func (r *logAlertConfigResource) mapModelTimeThresholdToAPI(timeThresholdModel *TimeThresholdModel) *restapi.LogTimeThreshold {
	if timeThresholdModel == nil || timeThresholdModel.ViolationsInSequence == nil {
		return nil
	}

	violationsModel := timeThresholdModel.ViolationsInSequence
	if violationsModel.TimeWindow.IsNull() || violationsModel.TimeWindow.IsUnknown() {
		return nil
	}

	return &restapi.LogTimeThreshold{
		Type:       TimeThresholdTypeViolationsInSequence,
		TimeWindow: violationsModel.TimeWindow.ValueInt64(),
	}
}

// mapModelRulesToAPI converts model rules to API representation
func (r *logAlertConfigResource) mapModelRulesToAPI(ruleModel *LogAlertRuleModel) ([]restapi.RuleWithThreshold[restapi.LogAlertRule], diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleModel == nil {
		return []restapi.RuleWithThreshold[restapi.LogAlertRule]{}, diags
	}

	logAlertRule := restapi.LogAlertRule{
		MetricName:  ruleModel.MetricName.ValueString(),
		AlertType:   r.convertAlertTypeToAPI(ruleModel.AlertType.ValueString()),
		Aggregation: r.extractAggregation(ruleModel.Aggregation),
	}

	thresholdOperator := restapi.ThresholdOperator(ruleModel.ThresholdOperator.ValueString())
	thresholdMap, thresholdDiags := r.mapModelThresholdsToAPI(ruleModel.Threshold)
	diags.Append(thresholdDiags...)

	result := []restapi.RuleWithThreshold[restapi.LogAlertRule]{
		{
			ThresholdOperator: thresholdOperator,
			Rule:              logAlertRule,
			Thresholds:        thresholdMap,
		},
	}

	return result, diags
}

// convertAlertTypeToAPI converts schema alert type to API format
func (r *logAlertConfigResource) convertAlertTypeToAPI(alertType string) string {
	if alertType == LogAlertTypeLogCount {
		return LogAlertTypeLogCountAPI
	}
	return alertType
}

// extractAggregation extracts aggregation from the model
func (r *logAlertConfigResource) extractAggregation(aggregation types.String) restapi.Aggregation {
	if !aggregation.IsNull() {
		return restapi.Aggregation(aggregation.ValueString())
	}
	return EmptyString
}

// mapModelThresholdsToAPI converts model thresholds to API representation
func (r *logAlertConfigResource) mapModelThresholdsToAPI(thresholdModel *ThresholdModel) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	if thresholdModel == nil {
		return thresholdMap, diags
	}

	if thresholdModel.Warning != nil {
		warningRule, warningDiags := r.convertThresholdModelToAPI(thresholdModel.Warning)
		diags.Append(warningDiags...)
		if warningRule != nil {
			thresholdMap[restapi.WarningSeverity] = *warningRule
		}
	}

	if thresholdModel.Critical != nil {
		criticalRule, criticalDiags := r.convertThresholdModelToAPI(thresholdModel.Critical)
		diags.Append(criticalDiags...)
		if criticalRule != nil {
			thresholdMap[restapi.CriticalSeverity] = *criticalRule
		}
	}

	return thresholdMap, diags
}

// convertThresholdModelToAPI converts a threshold model to API representation
func (r *logAlertConfigResource) convertThresholdModelToAPI(threshold *shared.ThresholdTypeModel) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if threshold == nil {
		return nil, diags
	}

	// Handle static threshold
	if threshold.Static != nil && !threshold.Static.Value.IsNull() {
		valueFloat := float64(threshold.Static.Value.ValueFloat32())
		rounded := math.Round(valueFloat*100) / 100
		return &restapi.ThresholdRule{
			Type:  ThresholdTypeStatic,
			Value: &rounded,
		}, diags
	}

	// Handle adaptive baseline threshold
	if threshold.AdaptiveBaseline != nil {
		seasonality := restapi.ThresholdSeasonality(threshold.AdaptiveBaseline.Seasonality.ValueString())
		deviationFactor := threshold.AdaptiveBaseline.DeviationFactor.ValueFloat32()
		adaptability := threshold.AdaptiveBaseline.Adaptability.ValueFloat32()
		return &restapi.ThresholdRule{
			Type:            "adaptiveBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Adaptability:    &adaptability,
		}, diags
	}

	return nil, diags
}

// Made with Bob
