package sloalertconfig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

// NewSloAlertConfigResourceHandleFramework creates the resource handle for SLO Alert configuration
func NewSloAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.SloAlertConfig] {
	return &sloAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSloAlertConfigFramework,
			Schema: schema.Schema{
				Description: SloAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: SloAlertConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: SloAlertConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					"description": schema.StringAttribute{
						Required:    true,
						Description: SloAlertConfigDescDescription,
					},
					"severity": schema.Int64Attribute{
						Required:    true,
						Description: SloAlertConfigDescSeverity,
					},
					"triggering": schema.BoolAttribute{
						Optional:    true,
						Description: SloAlertConfigDescTriggering,
					},
					"enabled": schema.BoolAttribute{
						Optional:    true,
						Description: SloAlertConfigDescEnabled,
					},
					"alert_type": schema.StringAttribute{
						Required:    true,
						Description: SloAlertConfigDescAlertType,
						Validators: []validator.String{
							stringvalidator.OneOf("status", "error_budget", "burn_rate_v2"),
						},
					},
					"slo_ids": schema.SetAttribute{
						Required:    true,
						Description: SloAlertConfigDescSloIds,
						ElementType: types.StringType,
					},
					"alert_channel_ids": schema.SetAttribute{
						Required:    true,
						Description: SloAlertConfigDescAlertChannelIds,
						ElementType: types.StringType,
					},
					"custom_payload_fields": shared.GetCustomPayloadFieldsSchema(),
					"threshold": schema.SingleNestedAttribute{
						Optional:    true,
						Description: SloAlertConfigDescThreshold,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Optional:    true,
								Description: SloAlertConfigDescThresholdType,
								Validators: []validator.String{
									stringvalidator.OneOf("staticThreshold"),
								},
							},
							"operator": schema.StringAttribute{
								Required:    true,
								Description: SloAlertConfigDescThresholdOperator,
								Validators: []validator.String{
									stringvalidator.OneOf(">", ">=", "=", "<=", "<"),
								},
							},
							"value": schema.Float64Attribute{
								Required:    true,
								Description: SloAlertConfigDescThresholdValue,
							},
						},
					},
					"time_threshold": schema.SingleNestedAttribute{
						Required:    true,
						Description: SloAlertConfigDescTimeThreshold,
						Attributes: map[string]schema.Attribute{
							"warm_up": schema.Int64Attribute{
								Required:    true,
								Description: SloAlertConfigDescTimeThresholdWarmUp,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
							},
							"cool_down": schema.Int64Attribute{
								Required:    true,
								Description: SloAlertConfigDescTimeThresholdCoolDown,
								Validators: []validator.Int64{
									int64validator.AtLeast(1),
								},
							},
						},
					},
					"burn_rate_config": schema.ListNestedAttribute{
						Optional:    true,
						Description: SloAlertConfigDescBurnRateConfig,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"alert_window_type": schema.StringAttribute{
									Required:    true,
									Description: SloAlertConfigDescBurnRateAlertWindowType,
								},
								"duration": schema.StringAttribute{
									Required:    true,
									Description: SloAlertConfigDescBurnRateDuration,
								},
								"duration_unit_type": schema.StringAttribute{
									Required:    true,
									Description: SloAlertConfigDescBurnRateDurationUnitType,
								},
								"threshold_operator": schema.StringAttribute{
									Required:    true,
									Description: SloAlertConfigDescBurnRateThresholdOperator,
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "=", "<=", "<"),
									},
								},
								"threshold_value": schema.StringAttribute{
									Required:    true,
									Description: SloAlertConfigDescBurnRateThresholdValue,
								},
							},
						},
					},
				},
			},
			SchemaVersion: 1,
			CreateOnly:    false,
		},
	}
}

func (r *sloAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *sloAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloAlertConfig] {
	return api.SloAlertConfig()
}

func (r *sloAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sloAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, sloAlertConfig *restapi.SloAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	model := SloAlertConfigModel{
		ID:          types.StringValue(sloAlertConfig.ID),
		Name:        types.StringValue(sloAlertConfig.Name),
		Description: types.StringValue(sloAlertConfig.Description),
		Severity:    types.Int64Value(int64(sloAlertConfig.Severity)),
		Triggering:  types.BoolValue(sloAlertConfig.Triggering),
		Enabled:     types.BoolValue(sloAlertConfig.Enabled),
	}

	// Map alert type from API to Terraform
	var terraformAlertType string
	switch sloAlertConfig.Rule.AlertType {
	case "SERVICE_LEVELS_OBJECTIVE":
		if sloAlertConfig.Rule.Metric == "STATUS" {
			terraformAlertType = "status"
		}
	case "ERROR_BUDGET":
		if sloAlertConfig.Rule.Metric == "BURNED_PERCENTAGE" {
			terraformAlertType = "error_budget"
		} else if sloAlertConfig.Rule.Metric == "BURN_RATE_V2" {
			terraformAlertType = "burn_rate_v2"
		}
	}
	model.AlertType = types.StringValue(terraformAlertType)

	// Map threshold
	if sloAlertConfig.Threshold != nil {
		thresholdType := sloAlertConfig.Threshold.Type
		if thresholdType == "static" {
			thresholdType = "staticThreshold"
		}

		model.Threshold = &SloAlertThresholdModel{
			Type:     types.StringValue(thresholdType),
			Operator: types.StringValue(sloAlertConfig.Threshold.Operator),
			Value:    types.Float64Value(sloAlertConfig.Threshold.Value),
		}
	} else {
		model.Threshold = nil
	}

	// Map time threshold
	model.TimeThreshold = &SloAlertTimeThresholdModel{
		WarmUp:   types.Int64Value(int64(sloAlertConfig.TimeThreshold.TimeWindow)),
		CoolDown: types.Int64Value(int64(sloAlertConfig.TimeThreshold.Expiry)),
	}

	// Map SLO IDs
	sloIds := []attr.Value{}
	for _, id := range sloAlertConfig.SloIds {
		sloIds = append(sloIds, types.StringValue(id))
	}
	model.SloIds = types.SetValueMust(types.StringType, sloIds)

	// Map Alert Channel IDs
	alertChannelIds := []attr.Value{}
	for _, id := range sloAlertConfig.AlertChannelIds {
		alertChannelIds = append(alertChannelIds, types.StringValue(id))
	}
	model.AlertChannelIds = types.SetValueMust(types.StringType, alertChannelIds)

	// Map burn rate configs
	if sloAlertConfig.BurnRateConfigs != nil && len(*sloAlertConfig.BurnRateConfigs) > 0 {
		burnRateConfigs := []SloAlertBurnRateConfigModel{}
		for _, cfg := range *sloAlertConfig.BurnRateConfigs {
			burnRateConfigs = append(burnRateConfigs, SloAlertBurnRateConfigModel{
				AlertWindowType:   types.StringValue(cfg.AlertWindowType),
				Duration:          types.StringValue(fmt.Sprintf("%d", cfg.Duration)),
				DurationUnitType:  types.StringValue(cfg.DurationUnitType),
				ThresholdOperator: types.StringValue(cfg.Threshold.Operator),
				ThresholdValue:    types.StringValue(fmt.Sprintf("%.2f", cfg.Threshold.Value)),
			})
		}
		model.BurnRateConfig = burnRateConfigs
	} else {
		model.BurnRateConfig = nil
	}

	// Map custom payload fields using the reusable function
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, sloAlertConfig.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayload = customPayloadFieldsList

	// Set the state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *sloAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SloAlertConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map alert type
	terraformAlertType := model.AlertType.ValueString()
	var apiAlertType, apiMetric string

	// Normalize alert type
	normalizedType := terraformAlertType
	switch terraformAlertType {
	case "errorBudget", "ErrorBudget":
		normalizedType = "error_budget"
	case "status", "Status":
		normalizedType = "status"
	case "burnRateV2", "BurnRateV2":
		normalizedType = "burn_rate_v2"
	}

	switch normalizedType {
	case "status":
		apiAlertType = "SERVICE_LEVELS_OBJECTIVE"
		apiMetric = "STATUS"
	case "error_budget":
		apiAlertType = "ERROR_BUDGET"
		apiMetric = "BURNED_PERCENTAGE"
	case "burn_rate_v2":
		apiAlertType = "ERROR_BUDGET"
		apiMetric = "BURN_RATE_V2"
	default:
		diags.AddError(
			SloAlertConfigErrMappingAlertType,
			fmt.Sprintf(SloAlertConfigErrInvalidAlertType, terraformAlertType),
		)
		return nil, diags
	}

	rule := restapi.SloAlertRule{
		AlertType: apiAlertType,
		Metric:    apiMetric,
	}

	// Map threshold
	var threshold *restapi.SloAlertThreshold
	if terraformAlertType != "burn_rate_v2" && model.Threshold != nil {
		thresholdType := "staticThreshold"
		if !model.Threshold.Type.IsNull() {
			thresholdType = model.Threshold.Type.ValueString()
		}

		threshold = &restapi.SloAlertThreshold{
			Type:     thresholdType,
			Operator: model.Threshold.Operator.ValueString(),
			Value:    model.Threshold.Value.ValueFloat64(),
		}
	}

	// Map time threshold
	var timeThreshold restapi.SloAlertTimeThreshold
	if model.TimeThreshold != nil {
		timeThreshold = restapi.SloAlertTimeThreshold{
			TimeWindow: int(model.TimeThreshold.WarmUp.ValueInt64()),
			Expiry:     int(model.TimeThreshold.CoolDown.ValueInt64()),
		}
	}

	// Map SLO IDs
	var sloIds []string
	if !model.SloIds.IsNull() && !model.SloIds.IsUnknown() {
		var sloIdValues []string
		diags.Append(model.SloIds.ElementsAs(ctx, &sloIdValues, false)...)
		if diags.HasError() {
			return nil, diags
		}
		sloIds = sloIdValues
	}

	// Map Alert Channel IDs
	var alertChannelIds []string
	if !model.AlertChannelIds.IsNull() && !model.AlertChannelIds.IsUnknown() {
		var alertChannelIdValues []string
		diags.Append(model.AlertChannelIds.ElementsAs(ctx, &alertChannelIdValues, false)...)
		if diags.HasError() {
			return nil, diags
		}
		alertChannelIds = alertChannelIdValues
	}

	// Map burn rate configs
	var burnRateConfigs []restapi.BurnRateConfig
	if terraformAlertType == "burn_rate_v2" && model.BurnRateConfig != nil && len(model.BurnRateConfig) > 0 {
		for _, burnRateConfigModel := range model.BurnRateConfig {
			duration, err := strconv.Atoi(burnRateConfigModel.Duration.ValueString())
			if err != nil {
				diags.AddError(
					SloAlertConfigErrParsingDuration,
					fmt.Sprintf(SloAlertConfigErrParsingDurationMsg, err),
				)
				return nil, diags
			}

			value, err := strconv.ParseFloat(burnRateConfigModel.ThresholdValue.ValueString(), 64)
			if err != nil {
				diags.AddError(
					SloAlertConfigErrParsingThresholdValue,
					fmt.Sprintf(SloAlertConfigErrParsingThresholdValueMsg, err),
				)
				return nil, diags
			}

			burnRateConfigs = append(burnRateConfigs, restapi.BurnRateConfig{
				AlertWindowType:  burnRateConfigModel.AlertWindowType.ValueString(),
				Duration:         duration,
				DurationUnitType: burnRateConfigModel.DurationUnitType.ValueString(),
				Threshold: restapi.ServiceLevelsStaticThresholdConfig{
					Operator: burnRateConfigModel.ThresholdOperator.ValueString(),
					Value:    value,
				},
			})
		}
	}

	// Map custom payload fields using the reusable function
	var customPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayload.IsNull() && !model.CustomPayload.IsUnknown() {
		var payloadDiags diag.Diagnostics
		customPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayload)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
	}

	// Create SLO alert config
	return &restapi.SloAlertConfig{
		ID:                    id,
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		Severity:              int(model.Severity.ValueInt64()),
		Triggering:            model.Triggering.ValueBool(),
		Enabled:               model.Enabled.ValueBool(),
		Rule:                  rule,
		Threshold:             threshold,
		TimeThreshold:         timeThreshold,
		SloIds:                sloIds,
		AlertChannelIds:       alertChannelIds,
		CustomerPayloadFields: customPayloadFields,
		BurnRateConfigs:       &burnRateConfigs,
	}, diags
}

// This section intentionally left empty to remove conflicting functions
