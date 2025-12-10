package sloalertconfig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
)

// NewSloAlertConfigResourceHandle creates the resource handle for SLO Alert configuration
func NewSloAlertConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.SloAlertConfig] {
	return &sloAlertConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSloAlertConfig,
			Schema:        buildSloAlertConfigSchema(),
			SchemaVersion: 2,
			CreateOnly:    false,
		},
	}
}

// buildSloAlertConfigSchema constructs the complete schema for SLO Alert configuration resource
func buildSloAlertConfigSchema() schema.Schema {
	return schema.Schema{
		Description: SloAlertConfigDescResource,
		Attributes: map[string]schema.Attribute{
			SchemaFieldID:                  buildIDAttribute(),
			SchemaFieldName:                buildNameAttribute(),
			SchemaFieldDescription:         buildDescriptionAttribute(),
			SchemaFieldSeverity:            buildSeverityAttribute(),
			SchemaFieldTriggering:          buildTriggeringAttribute(),
			SchemaFieldAlertType:           buildAlertTypeAttribute(),
			SchemaFieldSloIds:              buildSloIdsAttribute(),
			SchemaFieldAlertChannelIds:     buildAlertChannelIdsAttribute(),
			SchemaFieldCustomPayloadFields: shared.GetCustomPayloadFieldsSchema(),
			SchemaFieldThreshold:           buildThresholdAttribute(),
			SchemaFieldTimeThreshold:       buildTimeThresholdAttribute(),
			SchemaFieldBurnRateConfig:      buildBurnRateConfigAttribute(),
		},
	}
}

// buildIDAttribute creates the ID field schema attribute
func buildIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: SloAlertConfigDescID,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// buildNameAttribute creates the name field schema attribute
func buildNameAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescName,
		Validators: []validator.String{
			stringvalidator.LengthBetween(NameMinLength, NameMaxLength),
		},
	}
}

// buildDescriptionAttribute creates the description field schema attribute
func buildDescriptionAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescDescription,
	}
}

// buildSeverityAttribute creates the severity field schema attribute
func buildSeverityAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Required:    true,
		Description: SloAlertConfigDescSeverity,
	}
}

// buildTriggeringAttribute creates the triggering field schema attribute
func buildTriggeringAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:    true,
		Description: SloAlertConfigDescTriggering,
	}
}

// buildAlertTypeAttribute creates the alert_type field schema attribute
func buildAlertTypeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescAlertType,
		Validators: []validator.String{
			stringvalidator.OneOf(SloAlertConfigStatus, SloAlertConfigErrorBudget, SloAlertConfigBurnRateV2),
		},
	}
}

// buildSloIdsAttribute creates the slo_ids field schema attribute
func buildSloIdsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		Required:    true,
		Description: SloAlertConfigDescSloIds,
		ElementType: types.StringType,
	}
}

// buildAlertChannelIdsAttribute creates the alert_channel_ids field schema attribute
func buildAlertChannelIdsAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		Required:    true,
		Description: SloAlertConfigDescAlertChannelIds,
		ElementType: types.StringType,
	}
}

// buildThresholdAttribute creates the threshold nested attribute
func buildThresholdAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SloAlertConfigDescThreshold,
		Attributes: map[string]schema.Attribute{
			SchemaFieldThresholdType:     buildThresholdTypeAttribute(),
			SchemaFieldThresholdOperator: buildThresholdOperatorAttribute(),
			SchemaFieldThresholdValue:    buildThresholdValueAttribute(),
		},
	}
}

// buildThresholdTypeAttribute creates the threshold type field schema attribute
func buildThresholdTypeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Default:     stringdefault.StaticString(ThresholdTypeStaticThreshold),
		Description: SloAlertConfigDescThresholdType,
		Validators: []validator.String{
			stringvalidator.OneOf(ThresholdTypeStaticThreshold),
		},
	}
}

// buildThresholdOperatorAttribute creates the threshold operator field schema attribute
func buildThresholdOperatorAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescThresholdOperator,
		Validators: []validator.String{
			stringvalidator.OneOf(
				OperatorGreaterThan,
				OperatorGreaterThanOrEqual,
				OperatorEqual,
				OperatorLessThanOrEqual,
				OperatorLessThan,
			),
		},
	}
}

// buildThresholdValueAttribute creates the threshold value field schema attribute
func buildThresholdValueAttribute() schema.Float64Attribute {
	return schema.Float64Attribute{
		Required:    true,
		Description: SloAlertConfigDescThresholdValue,
	}
}

// buildTimeThresholdAttribute creates the time_threshold nested attribute
func buildTimeThresholdAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required:    true,
		Description: SloAlertConfigDescTimeThreshold,
		Attributes: map[string]schema.Attribute{
			SchemaFieldTimeThresholdWarmUp:   buildTimeThresholdWarmUpAttribute(),
			SchemaFieldTimeThresholdCoolDown: buildTimeThresholdCoolDownAttribute(),
		},
	}
}

// buildTimeThresholdWarmUpAttribute creates the warm_up field schema attribute
func buildTimeThresholdWarmUpAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Required:    true,
		Description: SloAlertConfigDescTimeThresholdWarmUp,
		Validators: []validator.Int64{
			int64validator.AtLeast(TimeThresholdMinValue),
		},
	}
}

// buildTimeThresholdCoolDownAttribute creates the cool_down field schema attribute
func buildTimeThresholdCoolDownAttribute() schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Description: SloAlertConfigDescTimeThresholdCoolDown,
		Validators: []validator.Int64{
			int64validator.AtLeast(0),
		},
	}
}

// buildBurnRateConfigAttribute creates the burn_rate_config list nested attribute
func buildBurnRateConfigAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: SloAlertConfigDescBurnRateConfig,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				SchemaFieldBurnRateAlertWindowType:   buildBurnRateAlertWindowTypeAttribute(),
				SchemaFieldBurnRateDuration:          buildBurnRateDurationAttribute(),
				SchemaFieldBurnRateDurationUnitType:  buildBurnRateDurationUnitTypeAttribute(),
				SchemaFieldBurnRateThresholdOperator: buildBurnRateThresholdOperatorAttribute(),
				SchemaFieldBurnRateThresholdValue:    buildBurnRateThresholdValueAttribute(),
			},
		},
	}
}

// buildBurnRateAlertWindowTypeAttribute creates the alert_window_type field schema attribute
func buildBurnRateAlertWindowTypeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescBurnRateAlertWindowType,
	}
}

// buildBurnRateDurationAttribute creates the duration field schema attribute
func buildBurnRateDurationAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescBurnRateDuration,
	}
}

// buildBurnRateDurationUnitTypeAttribute creates the duration_unit_type field schema attribute
func buildBurnRateDurationUnitTypeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescBurnRateDurationUnitType,
	}
}

// buildBurnRateThresholdOperatorAttribute creates the threshold_operator field schema attribute
func buildBurnRateThresholdOperatorAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescBurnRateThresholdOperator,
		Validators: []validator.String{
			stringvalidator.OneOf(
				OperatorGreaterThan,
				OperatorGreaterThanOrEqual,
				OperatorEqual,
				OperatorLessThanOrEqual,
				OperatorLessThan,
			),
		},
	}
}

// buildBurnRateThresholdValueAttribute creates the threshold_value field schema attribute
func buildBurnRateThresholdValueAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: SloAlertConfigDescBurnRateThresholdValue,
	}
}

func (r *sloAlertConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *sloAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloAlertConfig] {
	return api.SloAlertConfig()
}

func (r *sloAlertConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sloAlertConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, sloAlertConfig *restapi.SloAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model SloAlertConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	} else {
		model = SloAlertConfigModel{}
	}

	model.ID = types.StringValue(sloAlertConfig.ID)
	model.Name = types.StringValue(sloAlertConfig.Name)
	model.Description = types.StringValue(sloAlertConfig.Description)
	model.Severity = types.Int64Value(int64(sloAlertConfig.Severity))
	model.Triggering = types.BoolValue(sloAlertConfig.Triggering)

	terraformAlertType := r.mapAPIAlertTypeToTerraform(sloAlertConfig.Rule)
	model.AlertType = types.StringValue(terraformAlertType)

	model.Threshold = r.mapThresholdToState(sloAlertConfig.Threshold)
	model.TimeThreshold = r.mapTimeThresholdToState(sloAlertConfig.TimeThreshold)
	model.SloIds = r.mapStringSliceToSet(sloAlertConfig.SloIds)
	model.AlertChannelIds = r.mapStringSliceToSet(sloAlertConfig.AlertChannelIds)

	if len(model.BurnRateConfig) == 0 {
		burnRateConfigs, burnRateDiags := r.mapBurnRateConfigsToState(sloAlertConfig.BurnRateConfigs)
		diags.Append(burnRateDiags...)
		if diags.HasError() {
			return diags
		}
		model.BurnRateConfig = burnRateConfigs
	}

	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, sloAlertConfig.CustomerPayloadFields)
	diags.Append(payloadDiags...)
	if diags.HasError() {
		return diags
	}
	model.CustomPayload = customPayloadFieldsList

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapAPIAlertTypeToTerraform converts API alert type and metric to Terraform alert type
func (r *sloAlertConfigResource) mapAPIAlertTypeToTerraform(rule restapi.SloAlertRule) string {
	if rule.AlertType == APIAlertTypeServiceLevelsObjective && rule.Metric == APIMetricStatus {
		return SloAlertConfigStatus
	}

	if rule.AlertType == APIAlertTypeErrorBudget {
		if rule.Metric == APIMetricBurnedPercentage {
			return SloAlertConfigErrorBudget
		}
		if rule.Metric == APIMetricBurnRateV2 {
			return SloAlertConfigBurnRateV2
		}
		// Handle legacy BURN_RATE metric as error_budget
		if rule.Metric == APIMetricBurnRate {
			return SloAlertConfigErrorBudget
		}
	}

	return ""
}

// mapThresholdToState converts API threshold to state model
func (r *sloAlertConfigResource) mapThresholdToState(threshold *restapi.SloAlertThreshold) *SloAlertThresholdModel {
	if threshold == nil {
		return nil
	}

	thresholdType := threshold.Type
	if thresholdType == ThresholdTypeStatic {
		thresholdType = ThresholdTypeStaticThreshold
	}

	return &SloAlertThresholdModel{
		Type:     types.StringValue(thresholdType),
		Operator: types.StringValue(threshold.Operator),
		Value:    types.Float64Value(threshold.Value),
	}
}

// mapTimeThresholdToState converts API time threshold to state model
func (r *sloAlertConfigResource) mapTimeThresholdToState(timeThreshold restapi.SloAlertTimeThreshold) *SloAlertTimeThresholdModel {
	return &SloAlertTimeThresholdModel{
		WarmUp:   types.Int64Value(int64(timeThreshold.TimeWindow)),
		CoolDown: types.Int64Value(int64(timeThreshold.Expiry)),
	}
}

// mapStringSliceToSet converts a string slice to a Terraform set
func (r *sloAlertConfigResource) mapStringSliceToSet(values []string) types.Set {
	attrValues := make([]attr.Value, 0, len(values))
	for _, value := range values {
		attrValues = append(attrValues, types.StringValue(value))
	}
	return types.SetValueMust(types.StringType, attrValues)
}

// mapBurnRateConfigsToState converts API burn rate configs to state models
func (r *sloAlertConfigResource) mapBurnRateConfigsToState(burnRateConfigs *[]restapi.BurnRateConfig) ([]SloAlertBurnRateConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if burnRateConfigs == nil || len(*burnRateConfigs) == 0 {
		return nil, diags
	}

	configs := make([]SloAlertBurnRateConfigModel, 0, len(*burnRateConfigs))
	for _, cfg := range *burnRateConfigs {
		configs = append(configs, SloAlertBurnRateConfigModel{
			AlertWindowType:   types.StringValue(cfg.AlertWindowType),
			Duration:          types.StringValue(fmt.Sprintf("%d", cfg.Duration)),
			DurationUnitType:  types.StringValue(cfg.DurationUnitType),
			ThresholdOperator: types.StringValue(cfg.Threshold.Operator),
			ThresholdValue:    types.StringValue(fmt.Sprintf("%.2f", cfg.Threshold.Value)),
		})
	}

	return configs, diags
}

func (r *sloAlertConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SloAlertConfigModel

	diags.Append(r.extractModelFromPlanOrState(ctx, plan, state, &model)...)

	rule, ruleDiags := r.mapAlertTypeToAPIRule(model.AlertType.ValueString())
	diags.Append(ruleDiags...)

	threshold := r.mapThresholdFromState(model.AlertType.ValueString(), model.Threshold)
	timeThreshold := r.mapTimeThresholdFromState(model.TimeThreshold)

	sloIds, sloIdsDiags := r.mapSetToStringSlice(ctx, model.SloIds)
	diags.Append(sloIdsDiags...)

	alertChannelIds, channelIdsDiags := r.mapSetToStringSlice(ctx, model.AlertChannelIds)
	diags.Append(channelIdsDiags...)

	burnRateConfigs, burnRateDiags := r.mapBurnRateConfigsFromState(model.AlertType.ValueString(), model.BurnRateConfig)
	diags.Append(burnRateDiags...)

	customPayloadFields, payloadDiags := r.mapCustomPayloadFieldsFromState(ctx, model.CustomPayload)
	diags.Append(payloadDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &restapi.SloAlertConfig{
		ID:                    r.extractIDFromModel(model),
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		Severity:              int(model.Severity.ValueInt64()),
		Triggering:            model.Triggering.ValueBool(),
		Enabled:               true,
		Rule:                  rule,
		Threshold:             threshold,
		TimeThreshold:         timeThreshold,
		SloIds:                sloIds,
		AlertChannelIds:       alertChannelIds,
		CustomerPayloadFields: customPayloadFields,
		BurnRateConfigs:       &burnRateConfigs,
	}, diags
}

// extractModelFromPlanOrState retrieves the model from plan or state
func (r *sloAlertConfigResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State, model *SloAlertConfigModel) diag.Diagnostics {
	if plan != nil {
		return plan.Get(ctx, model)
	}
	if state != nil {
		return state.Get(ctx, model)
	}
	return nil
}

// extractIDFromModel extracts ID from model, returning empty string if null
func (r *sloAlertConfigResource) extractIDFromModel(model SloAlertConfigModel) string {
	if model.ID.IsNull() {
		return ""
	}
	return model.ID.ValueString()
}

// mapAlertTypeToAPIRule converts Terraform alert type to API rule
func (r *sloAlertConfigResource) mapAlertTypeToAPIRule(terraformAlertType string) (restapi.SloAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	normalizedType := r.normalizeAlertType(terraformAlertType)

	switch normalizedType {
	case SloAlertConfigStatus:
		return restapi.SloAlertRule{
			AlertType: APIAlertTypeServiceLevelsObjective,
			Metric:    APIMetricStatus,
		}, diags
	case SloAlertConfigErrorBudget:
		return restapi.SloAlertRule{
			AlertType: APIAlertTypeErrorBudget,
			Metric:    APIMetricBurnedPercentage,
		}, diags
	case SloAlertConfigBurnRateV2:
		return restapi.SloAlertRule{
			AlertType: APIAlertTypeErrorBudget,
			Metric:    APIMetricBurnRateV2,
		}, diags
	default:
		diags.AddError(
			SloAlertConfigErrMappingAlertType,
			fmt.Sprintf(SloAlertConfigErrInvalidAlertType, terraformAlertType),
		)
		return restapi.SloAlertRule{}, diags
	}
}

// normalizeAlertType normalizes alert type to standard format
func (r *sloAlertConfigResource) normalizeAlertType(alertType string) string {
	switch alertType {
	case AlertTypeErrorBudgetAlt1, AlertTypeErrorBudgetAlt2:
		return SloAlertConfigErrorBudget
	case AlertTypeStatusAlt1:
		return SloAlertConfigStatus
	case AlertTypeBurnRateV2Alt1, AlertTypeBurnRateV2Alt2:
		return SloAlertConfigBurnRateV2
	default:
		return alertType
	}
}

// mapThresholdFromState converts state threshold to API threshold
func (r *sloAlertConfigResource) mapThresholdFromState(alertType string, threshold *SloAlertThresholdModel) *restapi.SloAlertThreshold {
	if threshold == nil {
		return nil
	}

	thresholdType := ThresholdTypeStaticThreshold
	if !threshold.Type.IsNull() {
		thresholdType = threshold.Type.ValueString()
	}

	return &restapi.SloAlertThreshold{
		Type:     thresholdType,
		Operator: threshold.Operator.ValueString(),
		Value:    threshold.Value.ValueFloat64(),
	}
}

// mapTimeThresholdFromState converts state time threshold to API time threshold
func (r *sloAlertConfigResource) mapTimeThresholdFromState(timeThreshold *SloAlertTimeThresholdModel) restapi.SloAlertTimeThreshold {
	if timeThreshold == nil {
		return restapi.SloAlertTimeThreshold{}
	}

	return restapi.SloAlertTimeThreshold{
		TimeWindow: int(timeThreshold.WarmUp.ValueInt64()),
		Expiry:     int(timeThreshold.CoolDown.ValueInt64()),
	}
}

// mapSetToStringSlice converts a Terraform set to a string slice
func (r *sloAlertConfigResource) mapSetToStringSlice(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if set.IsNull() || set.IsUnknown() {
		return nil, diags
	}

	var values []string
	diags.Append(set.ElementsAs(ctx, &values, false)...)
	return values, diags
}

// mapBurnRateConfigsFromState converts state burn rate configs to API burn rate configs
func (r *sloAlertConfigResource) mapBurnRateConfigsFromState(alertType string, configs []SloAlertBurnRateConfigModel) ([]restapi.BurnRateConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	if alertType != SloAlertConfigBurnRateV2 || configs == nil || len(configs) == 0 {
		return []restapi.BurnRateConfig{}, diags
	}

	burnRateConfigs := make([]restapi.BurnRateConfig, 0, len(configs))
	for _, configModel := range configs {
		duration, durationErr := r.parseIntFromString(configModel.Duration.ValueString(), SloAlertConfigErrParsingDuration, SloAlertConfigErrParsingDurationMsg)
		if durationErr != nil {
			diags.AddError(SloAlertConfigErrParsingDuration, durationErr.Error())
			return nil, diags
		}

		value, valueErr := r.parseFloatFromString(configModel.ThresholdValue.ValueString(), SloAlertConfigErrParsingThresholdValue, SloAlertConfigErrParsingThresholdValueMsg)
		if valueErr != nil {
			diags.AddError(SloAlertConfigErrParsingThresholdValue, valueErr.Error())
			return nil, diags
		}

		burnRateConfigs = append(burnRateConfigs, restapi.BurnRateConfig{
			AlertWindowType:  configModel.AlertWindowType.ValueString(),
			Duration:         duration,
			DurationUnitType: configModel.DurationUnitType.ValueString(),
			Threshold: restapi.ServiceLevelsStaticThresholdConfig{
				Operator: configModel.ThresholdOperator.ValueString(),
				Value:    value,
			},
		})
	}

	return burnRateConfigs, diags
}

// parseIntFromString parses an integer from a string with error handling
func (r *sloAlertConfigResource) parseIntFromString(value, errorTitle, errorMsgFormat string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf(errorMsgFormat, err)
	}
	return intValue, nil
}

// parseFloatFromString parses a float from a string with error handling
func (r *sloAlertConfigResource) parseFloatFromString(value, errorTitle, errorMsgFormat string) (float64, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(errorMsgFormat, err)
	}
	return floatValue, nil
}

// mapCustomPayloadFieldsFromState converts state custom payload fields to API custom payload fields
func (r *sloAlertConfigResource) mapCustomPayloadFieldsFromState(ctx context.Context, customPayload types.List) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics

	if customPayload.IsNull() || customPayload.IsUnknown() {
		return nil, diags
	}

	customPayloadFields, payloadDiags := shared.MapCustomPayloadFieldsToAPIObject(ctx, customPayload)
	diags.Append(payloadDiags...)
	return customPayloadFields, diags
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *sloAlertConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: resourcehandle.CreateStateUpgraderForVersion(1),
	}
}
