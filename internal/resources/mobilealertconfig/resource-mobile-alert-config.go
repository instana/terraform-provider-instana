package mobilealertconfig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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

// NewMobileAlertConfigResourceHandle creates the resource handle for Mobile Alert Configs
func NewMobileAlertConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.MobileAlertConfig] {
	return &mobileAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaMobileAlertConfig,
			Schema: schema.Schema{
				Description: MobileAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					MobileAlertConfigFieldID: schema.StringAttribute{
						Computed:    true,
						Description: MobileAlertConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					MobileAlertConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: MobileAlertConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(MobileAlertConfigMinNameLength, MobileAlertConfigMaxNameLength),
						},
					},
					MobileAlertConfigFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: MobileAlertConfigDescDescription,
						Validators: []validator.String{
							stringvalidator.LengthBetween(MobileAlertConfigMinDescriptionLength, MobileAlertConfigMaxDescriptionLength),
						},
					},
					MobileAlertConfigFieldMobileAppID: schema.StringAttribute{
						Required:    true,
						Description: MobileAlertConfigDescMobileAppID,
						Validators: []validator.String{
							stringvalidator.LengthBetween(MobileAlertConfigMinMobileAppIDLength, MobileAlertConfigMaxMobileAppIDLength),
						},
					},
					MobileAlertConfigFieldTriggering: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: MobileAlertConfigDescTriggering,
						Default:     booldefault.StaticBool(MobileAlertConfigDefaultTriggering),
					},
					MobileAlertConfigFieldTagFilter: schema.StringAttribute{
						Optional:    true,
						Description: MobileAlertConfigDescTagFilter,
					},
					MobileAlertConfigFieldAlertChannels: schema.MapAttribute{
						Optional:    true,
						Description: MobileAlertConfigDescAlertChannels,
						ElementType: types.SetType{ElemType: types.StringType},
					},
					MobileAlertConfigFieldGranularity: schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: MobileAlertConfigDescGranularity,
						Default:     int64default.StaticInt64(MobileAlertConfigDefaultGranularity),
					},
					MobileAlertConfigFieldGracePeriod: schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: MobileAlertConfigDescGracePeriod,
					},
					MobileAlertConfigFieldCustomPayloadFields: shared.GetCustomPayloadFieldsSchema(),
					MobileAlertConfigFieldRules: schema.ListNestedAttribute{
						Description: MobileAlertConfigDescRules,
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								MobileAlertConfigFieldRule: schema.SingleNestedAttribute{
									Description: MobileAlertConfigDescRule,
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										MobileAlertConfigFieldRuleAlertType: schema.StringAttribute{
											Required:    true,
											Description: MobileAlertConfigDescRuleAlertType,
										},
										MobileAlertConfigFieldRuleMetricName: schema.StringAttribute{
											Required:    true,
											Description: MobileAlertConfigDescRuleMetricName,
										},
										MobileAlertConfigFieldRuleAggregation: schema.StringAttribute{
											Optional:    true,
											Computed:    true,
											Description: MobileAlertConfigDescRuleAggregation,
											Validators: []validator.String{
												stringvalidator.OneOf(
													MobileAlertConfigAggregationSUM,
													MobileAlertConfigAggregationMEAN,
													MobileAlertConfigAggregationMAX,
													MobileAlertConfigAggregationMIN,
													MobileAlertConfigAggregationP25,
													MobileAlertConfigAggregationP50,
													MobileAlertConfigAggregationP75,
													MobileAlertConfigAggregationP90,
													MobileAlertConfigAggregationP95,
													MobileAlertConfigAggregationP98,
													MobileAlertConfigAggregationP99,
												),
											},
										},
										MobileAlertConfigFieldRuleOperator: schema.StringAttribute{
											Optional:    true,
											Computed:    true,
											Description: MobileAlertConfigDescRuleOperator,
											Validators: []validator.String{
												stringvalidator.OneOf("STARTS_WITH", "EQUALS"),
											},
										},
										MobileAlertConfigFieldRuleValue: schema.StringAttribute{
											Optional:    true,
											Computed:    true,
											Description: MobileAlertConfigDescRuleValue,
										},
									},
								},
								MobileAlertConfigFieldThresholdOperator: schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: MobileAlertConfigDescThresholdOperator,
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "<", "<="),
									},
								},
								MobileAlertConfigFieldThreshold: schema.SingleNestedAttribute{
									Description: MobileAlertConfigDescThreshold,
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										shared.LogAlertConfigFieldWarning:  shared.AllThresholdAttributeSchema(),
										shared.LogAlertConfigFieldCritical: shared.AllThresholdAttributeSchema(),
									},
								},
							},
						},
					},
					MobileAlertConfigFieldTimeThreshold: schema.SingleNestedAttribute{
						Description: MobileAlertConfigDescTimeThreshold,
						Required:    true,
						Attributes: map[string]schema.Attribute{
							MobileAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: schema.SingleNestedAttribute{
								Description: MobileAlertConfigDescTimeThresholdUserImpact,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									MobileAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: MobileAlertConfigDescTimeThresholdTimeWindow,
									},
									MobileAlertConfigFieldTimeThresholdUserImpactUsers: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: MobileAlertConfigDescTimeThresholdUsers,
									},
									MobileAlertConfigFieldTimeThresholdUserImpactPercentage: schema.Float64Attribute{
										Optional:    true,
										Computed:    true,
										Description: MobileAlertConfigDescTimeThresholdPercentage,
									},
								},
							},
							MobileAlertConfigFieldTimeThresholdViolationsInPeriod: schema.SingleNestedAttribute{
								Description: MobileAlertConfigDescTimeThresholdViolationsInPeriod,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									MobileAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: MobileAlertConfigDescTimeThresholdTimeWindow,
									},
									MobileAlertConfigFieldTimeThresholdViolationsInPeriodViolations: schema.Int64Attribute{
										Required:    true,
										Description: MobileAlertConfigDescTimeThresholdViolations,
									},
								},
							},
							MobileAlertConfigFieldTimeThresholdViolationsInSequence: schema.SingleNestedAttribute{
								Description: MobileAlertConfigDescTimeThresholdViolationsInSequence,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									MobileAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: MobileAlertConfigDescTimeThresholdTimeWindow,
										Default:     int64default.StaticInt64(600000),
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 0,
		},
	}
}

type mobileAlertConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *mobileAlertConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *mobileAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.MobileAlertConfig] {
	return api.MobileAlertConfig()
}

func (r *mobileAlertConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *mobileAlertConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.MobileAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model MobileAlertConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map tag filter
	tagFilter, tagFilterDiags := r.mapTagFilterToAPI(model)
	diags.Append(tagFilterDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map alert channels
	alertChannels, channelDiags := r.mapAlertChannelsToAPI(ctx, model)
	diags.Append(channelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map custom payload fields
	customPayloadFields, payloadDiags := r.mapCustomPayloadFieldsToAPI(ctx, model)
	diags.Append(payloadDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map time threshold
	timeThreshold, timeThresholdDiags := r.mapTimeThresholdFromModel(ctx, model)
	diags.Append(timeThresholdDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map rules collection
	rules, rulesDiags := r.mapRulesCollectionToAPI(ctx, model)
	diags.Append(rulesDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map grace period
	var gracePeriod *int64
	if !model.GracePeriod.IsNull() && !model.GracePeriod.IsUnknown() {
		gp := model.GracePeriod.ValueInt64()
		gracePeriod = &gp
	}

	// Create API object
	return &restapi.MobileAlertConfig{
		ID:                  model.ID.ValueString(),
		Name:                model.Name.ValueString(),
		Description:         model.Description.ValueString(),
		MobileAppID:         model.MobileAppID.ValueString(),
		Severity:            nil,
		Triggering:          model.Triggering.ValueBool(),
		TagFilterExpression:   tagFilter,
		AlertChannels:         alertChannels,
		Granularity:           restapi.Granularity(model.Granularity.ValueInt64()),
		GracePeriod:           gracePeriod,
		CustomerPayloadFields: customPayloadFields,
		Rules:                 rules,
		TimeThreshold:         timeThreshold,
	}, diags
}

// mapTagFilterToAPI maps tag filter expression from model to API
func (r *mobileAlertConfigResource) mapTagFilterToAPI(model MobileAlertConfigModel) (*restapi.TagFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError(MobileAlertConfigErrParseTagFilter, err.Error())
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		return mapper.ToAPIModel(expr), diags
	}

	// Return default tag filter
	operator := restapi.LogicalOperatorType(MobileAlertConfigLogicalOperatorAND)
	return &restapi.TagFilter{
		Type:            MobileAlertConfigTagFilterTypeExpression,
		LogicalOperator: &operator,
		Elements:        []*restapi.TagFilter{},
	}, diags
}

// mapAlertChannelsToAPI maps alert channels from model to API
func (r *mobileAlertConfigResource) mapAlertChannelsToAPI(ctx context.Context, model MobileAlertConfigModel) (map[string][]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannels := make(map[string][]string)

	if !model.AlertChannels.IsNull() && !model.AlertChannels.IsUnknown() {
		diags.Append(model.AlertChannels.ElementsAs(ctx, &alertChannels, false)...)
	}

	return alertChannels, diags
}

// mapCustomPayloadFieldsToAPI maps custom payload fields from model to API
func (r *mobileAlertConfigResource) mapCustomPayloadFieldsToAPI(ctx context.Context, model MobileAlertConfigModel) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics
	customPayloadFields := make([]restapi.CustomPayloadField[any], 0)

	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		var payloadDiags diag.Diagnostics
		customPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		diags.Append(payloadDiags...)
	}

	return customPayloadFields, diags
}

// mapRulesCollectionToAPI maps the rules collection from model to API
func (r *mobileAlertConfigResource) mapRulesCollectionToAPI(ctx context.Context, model MobileAlertConfigModel) ([]restapi.MobileAppAlertRuleWithThresholds, diag.Diagnostics) {
	var diags diag.Diagnostics
	rules := make([]restapi.MobileAppAlertRuleWithThresholds, 0)

	if len(model.Rules) == 0 {
		return rules, diags
	}

	// Process each rule
	for _, ruleModel := range model.Rules {
		mobileAlertRule, ruleDiags := r.mapSingleRuleFromCollection(ctx, ruleModel)
		diags.Append(ruleDiags...)
		if diags.HasError() {
			return nil, diags
		}

		// Only process if we have a valid rule
		if mobileAlertRule != nil {
			thresholdMap, thresholdDiags := shared.MapThresholdsAllPluginFromState(ctx, ruleModel.Thresholds)
			diags.Append(thresholdDiags...)
			if diags.HasError() {
				return nil, diags
			}

			rules = append(rules, restapi.MobileAppAlertRuleWithThresholds{
				Rule:              mobileAlertRule,
				ThresholdOperator: ruleModel.ThresholdOperator.ValueString(),
				Thresholds:        thresholdMap,
			})
		}
	}

	return rules, diags
}

// mapSingleRuleFromCollection maps a single rule from the rules collection
func (r *mobileAlertConfigResource) mapSingleRuleFromCollection(ctx context.Context, ruleModel MobileRuleWithThresholdModel) (*restapi.MobileAppAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ruleModel.Rule == nil {
		return nil, diags
	}

	var aggregationPtr *restapi.Aggregation
	if !ruleModel.Rule.Aggregation.IsNull() && !ruleModel.Rule.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(ruleModel.Rule.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	var operatorPtr *string
	if !ruleModel.Rule.Operator.IsNull() && !ruleModel.Rule.Operator.IsUnknown() {
		operator := ruleModel.Rule.Operator.ValueString()
		operatorPtr = &operator
	}

	var valuePtr *string
	if !ruleModel.Rule.Value.IsNull() && !ruleModel.Rule.Value.IsUnknown() {
		value := ruleModel.Rule.Value.ValueString()
		valuePtr = &value
	}

	return &restapi.MobileAppAlertRule{
		AlertType:   ruleModel.Rule.AlertType.ValueString(),
		MetricName:  ruleModel.Rule.MetricName.ValueString(),
		Aggregation: aggregationPtr,
		Operator:    operatorPtr,
		Value:       valuePtr,
	}, diags
}

// mapTimeThresholdFromModel converts the time threshold configuration from the Terraform model to the API representation
func (r *mobileAlertConfigResource) mapTimeThresholdFromModel(ctx context.Context, model MobileAlertConfigModel) (*restapi.MobileAppTimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate time threshold is provided
	if model.TimeThreshold == nil {
		diags.AddError(MobileAlertConfigErrTimeThresholdRequired, MobileAlertConfigErrTimeThresholdRequiredMsg)
		return nil, diags
	}

	timeThresholdModel := *model.TimeThreshold

	// Map to appropriate time threshold type based on which field is set
	switch {
	case timeThresholdModel.UserImpactOfViolationsInSequence != nil:
		return r.mapUserImpactTimeThreshold(timeThresholdModel.UserImpactOfViolationsInSequence, diags)
	case timeThresholdModel.ViolationsInPeriod != nil:
		return r.mapViolationsInPeriodTimeThreshold(timeThresholdModel.ViolationsInPeriod, diags)
	case timeThresholdModel.ViolationsInSequence != nil:
		return r.mapViolationsInSequenceTimeThreshold(timeThresholdModel.ViolationsInSequence, diags)
	default:
		diags.AddError(MobileAlertConfigErrInvalidTimeThresholdConfig, MobileAlertConfigErrInvalidTimeThresholdConfigMsg)
		return nil, diags
	}
}

// mapUserImpactTimeThreshold maps the user impact of violations in sequence time threshold configuration
func (r *mobileAlertConfigResource) mapUserImpactTimeThreshold(
	userImpactModel *MobileUserImpactOfViolationsInSequenceModel,
	diags diag.Diagnostics,
) (*restapi.MobileAppTimeThreshold, diag.Diagnostics) {
	threshold := &restapi.MobileAppTimeThreshold{
		Type: MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence,
	}

	// Map time window
	threshold.TimeWindow = r.mapOptionalInt64Field(userImpactModel.TimeWindow)

	// Map users count
	threshold.Users = r.mapOptionalInt64ToInt32Field(userImpactModel.Users)

	// Map percentage
	threshold.Percentage = r.mapOptionalFloat64Field(userImpactModel.Percentage)

	return threshold, diags
}

// mapViolationsInPeriodTimeThreshold maps the violations in period time threshold configuration
func (r *mobileAlertConfigResource) mapViolationsInPeriodTimeThreshold(
	violationsInPeriodModel *MobileViolationsInPeriodModel,
	diags diag.Diagnostics,
) (*restapi.MobileAppTimeThreshold, diag.Diagnostics) {
	threshold := &restapi.MobileAppTimeThreshold{
		Type: MobileAlertConfigTimeThresholdTypeViolationsInPeriod,
	}

	// Map time window
	threshold.TimeWindow = r.mapOptionalInt64Field(violationsInPeriodModel.TimeWindow)

	// Map violations count
	threshold.Violations = r.mapOptionalInt64ToInt32Field(violationsInPeriodModel.Violations)

	return threshold, diags
}

// mapViolationsInSequenceTimeThreshold maps the violations in sequence time threshold configuration
func (r *mobileAlertConfigResource) mapViolationsInSequenceTimeThreshold(
	violationsInSequenceModel *MobileViolationsInSequenceModel,
	diags diag.Diagnostics,
) (*restapi.MobileAppTimeThreshold, diag.Diagnostics) {
	threshold := &restapi.MobileAppTimeThreshold{
		Type: MobileAlertConfigTimeThresholdTypeViolationsInSequence,
	}

	// Map time window
	threshold.TimeWindow = r.mapOptionalInt64Field(violationsInSequenceModel.TimeWindow)

	return threshold, diags
}

// mapOptionalInt64Field extracts an optional int64 field value, returning nil if null or unknown
func (r *mobileAlertConfigResource) mapOptionalInt64Field(field types.Int64) *int64 {
	if field.IsNull() || field.IsUnknown() {
		return nil
	}
	value := field.ValueInt64()
	return &value
}

// mapOptionalFloat64Field extracts an optional float64 field value, returning nil if null or unknown
func (r *mobileAlertConfigResource) mapOptionalFloat64Field(field types.Float64) *float64 {
	if field.IsNull() || field.IsUnknown() {
		return nil
	}
	value := field.ValueFloat64()
	return &value
}

// mapOptionalInt64ToInt32Field extracts an optional int64 field and converts it to int32, returning nil if null or unknown
func (r *mobileAlertConfigResource) mapOptionalInt64ToInt32Field(field types.Int64) *int32 {
	if field.IsNull() || field.IsUnknown() {
		return nil
	}
	value := int32(field.ValueInt64())
	return &value
}

func (r *mobileAlertConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.MobileAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model MobileAlertConfigModel

	// Get current state from plan or state to preserve optional fields
	// Don't check for errors here as state might be empty on first create
	if plan != nil {
		plan.Get(ctx, &model)
	} else if state != nil {
		state.Get(ctx, &model)
	}

	// Create base model
	model.ID = types.StringValue(apiObject.ID)
	model.Name = types.StringValue(apiObject.Name)
	model.Description = types.StringValue(apiObject.Description)
	model.MobileAppID = types.StringValue(apiObject.MobileAppID)
	model.Triggering = types.BoolValue(apiObject.Triggering)
	model.Granularity = types.Int64Value(int64(apiObject.Granularity))

	// Map grace period
	if apiObject.GracePeriod != nil {
		model.GracePeriod = types.Int64Value(*apiObject.GracePeriod)
	} else {
		model.GracePeriod = types.Int64Null()
	}

	// Map tag filter expression
	tagFilterDiags := r.mapTagFilterExpressionToState(&model, apiObject)
	diags.Append(tagFilterDiags...)
	if diags.HasError() {
		return diags
	}

	// Map alert channels
	r.mapAlertChannelsToState(ctx, &model, apiObject)

	// Map rules collection (preserve the plan values / update only if empty)
	if len(model.Rules) == 0 {
		model.Rules = r.mapRulesToState(ctx, apiObject)
	}

	// Map custom payload fields
	customPayloadDiags := r.mapCustomPayloadFieldsToState(ctx, &model, apiObject)
	diags.Append(customPayloadDiags...)
	if diags.HasError() {
		return diags
	}

	// Map time threshold
	if apiObject.TimeThreshold != nil {
		model.TimeThreshold = r.mapTimeThresholdToState(*apiObject.TimeThreshold)
	}

	// Set state
	diags.Append(state.Set(ctx, &model)...)
	return diags
}

// mapTagFilterExpressionToState maps tag filter expression from API to state
func (r *mobileAlertConfigResource) mapTagFilterExpressionToState(model *MobileAlertConfigModel, apiObject *restapi.MobileAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject.TagFilterExpression != nil {
		filterExprStr, err := tagfilter.MapTagFilterToNormalizedString(apiObject.TagFilterExpression)
		if err != nil {
			diags.AddError(
				MobileAlertConfigErrMapFilterExpression,
				fmt.Sprintf(MobileAlertConfigErrMapFilterExpressionMsg, err),
			)
			return diags
		}
		model.TagFilter = util.SetStringPointerToState(filterExprStr)
	} else {
		model.TagFilter = types.StringNull()
	}

	return diags
}

// mapAlertChannelsToState maps alert channels from API to state
func (r *mobileAlertConfigResource) mapAlertChannelsToState(ctx context.Context, model *MobileAlertConfigModel, apiObject *restapi.MobileAlertConfig) {
	if len(apiObject.AlertChannels) > 0 {
		alertChannelsMap := make(map[string]attr.Value)
		for severity, channelIDs := range apiObject.AlertChannels {
			channelElements := make([]attr.Value, len(channelIDs))
			for i, id := range channelIDs {
				channelElements[i] = types.StringValue(id)
			}
			alertChannelsMap[severity] = types.SetValueMust(types.StringType, channelElements)
		}
		model.AlertChannels = types.MapValueMust(types.SetType{ElemType: types.StringType}, alertChannelsMap)
	} else {
		model.AlertChannels = types.MapNull(types.SetType{ElemType: types.StringType})
	}
}

// mapCustomPayloadFieldsToState maps custom payload fields from API to state
func (r *mobileAlertConfigResource) mapCustomPayloadFieldsToState(ctx context.Context, model *MobileAlertConfigModel, apiObject *restapi.MobileAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, apiObject.CustomerPayloadFields)
	diags.Append(payloadDiags...)
	if !diags.HasError() {
		model.CustomPayloadFields = customPayloadFieldsList
	}

	return diags
}

func (r *mobileAlertConfigResource) mapTimeThresholdToState(timeThreshold restapi.MobileAppTimeThreshold) *MobileAlertTimeThresholdModel {
	mobileTimeThresholdModel := MobileAlertTimeThresholdModel{}
	switch timeThreshold.Type {
	case MobileAlertConfigTimeThresholdTypeViolationsInSequence:
		mobileTimeThresholdModel.ViolationsInSequence = &MobileViolationsInSequenceModel{
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
		}
	case MobileAlertConfigTimeThresholdTypeUserImpactOfViolationsInSequence:
		mobileTimeThresholdModel.UserImpactOfViolationsInSequence = &MobileUserImpactOfViolationsInSequenceModel{
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
			Users:      util.SetInt64PointerToState(timeThreshold.Users),
			Percentage: util.SetFloat64PointerToState(timeThreshold.Percentage),
		}
	case MobileAlertConfigTimeThresholdTypeViolationsInPeriod:
		mobileTimeThresholdModel.ViolationsInPeriod = &MobileViolationsInPeriodModel{
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
			Violations: util.SetInt64PointerToState(timeThreshold.Violations),
		}
	}
	return &mobileTimeThresholdModel
}

func (r *mobileAlertConfigResource) mapRuleToState(ctx context.Context, rule *restapi.MobileAppAlertRule) *MobileAlertRuleModel {
	var aggregation types.String
	if rule.Aggregation != nil {
		aggregation = types.StringValue(string(*rule.Aggregation))
	} else {
		aggregation = types.StringNull()
	}

	var operator types.String
	if rule.Operator != nil {
		operator = types.StringValue(*rule.Operator)
	} else {
		operator = types.StringNull()
	}

	var value types.String
	if rule.Value != nil {
		value = types.StringValue(*rule.Value)
	} else {
		value = types.StringNull()
	}

	return &MobileAlertRuleModel{
		AlertType:   types.StringValue(rule.AlertType),
		MetricName:  types.StringValue(rule.MetricName),
		Aggregation: aggregation,
		Operator:    operator,
		Value:       value,
	}
}

func (r *mobileAlertConfigResource) mapRulesToState(ctx context.Context, apiObject *restapi.MobileAlertConfig) []MobileRuleWithThresholdModel {
	rules := apiObject.Rules
	var rulesModel []MobileRuleWithThresholdModel
	for _, i := range rules {
		warningThreshold, isWarningThresholdPresent := i.Thresholds[restapi.WarningSeverity]
		criticalThreshold, isCriticalThresholdPresent := i.Thresholds[restapi.CriticalSeverity]

		thresholdPluginModel := shared.ThresholdAllPluginModel{
			Warning:  shared.MapAllThresholdPluginToState(ctx, &warningThreshold, isWarningThresholdPresent),
			Critical: shared.MapAllThresholdPluginToState(ctx, &criticalThreshold, isCriticalThresholdPresent),
		}
		ruleModel := MobileRuleWithThresholdModel{
			Rule:              r.mapRuleToState(ctx, i.Rule),
			ThresholdOperator: types.StringValue(i.ThresholdOperator),
			Thresholds:        &thresholdPluginModel,
		}
		rulesModel = append(rulesModel, ruleModel)
	}
	return rulesModel
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *mobileAlertConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
