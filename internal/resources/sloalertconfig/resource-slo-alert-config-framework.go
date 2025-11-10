package sloalertconfig

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
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

// ResourceInstanaSloAlertConfigFramework the name of the terraform-provider-instana resource to manage SLO Alert configurations
const ResourceInstanaSloAlertConfigFramework = "slo_alert_config"

const (
	//Slo Alert Config Field names for Terraform
	SloAlertConfigFieldName                            = "name"
	SloAlertConfigFieldFullName                        = "full_name"
	SloAlertConfigFieldDescription                     = "description"
	SloAlertConfigFieldSeverity                        = "severity"
	SloAlertConfigFieldTriggering                      = "triggering"
	SloAlertConfigFieldAlertType                       = "alert_type"
	SloAlertConfigFieldThreshold                       = "threshold"
	SloAlertConfigFieldThresholdType                   = "type"
	SloAlertConfigFieldThresholdOperator               = "operator"
	SloAlertConfigFieldThresholdValue                  = "value"
	SloAlertConfigFieldSloIds                          = "slo_ids"
	SloAlertConfigFieldAlertChannelIds                 = "alert_channel_ids"
	SloAlertConfigFieldTimeThreshold                   = "time_threshold"
	SloAlertConfigFieldTimeThresholdWarmUp             = "warm_up"
	SloAlertConfigFieldTimeThresholdCoolDown           = "cool_down"
	SloAlertConfigFieldEnabled                         = "enabled"
	SloAlertConfigFieldBurnRateConfig                  = "burn_rate_config"
	SloAlertConfigFieldBurnRateConfigDuration          = "duration"
	SloAlertConfigFieldBurnRateConfigThresholdValue    = "threshold_value"
	SloAlertConfigFieldBurnRateConfigThresholdOperator = "threshold_operator"
	SloAlertConfigFieldBurnRateConfigDurationUnitType  = "duration_unit_type"
	SloAlertConfigFieldBurnRateConfigAlertWindowType   = "alert_window_type"

	// Slo Alert Types for Terraform
	SloAlertConfigStatus      = "status"
	SloAlertConfigErrorBudget = "error_budget"
	SloAlertConfigBurnRateV2  = "burn_rate_v2"
)

// SloAlertConfigModel represents the data model for SLO Alert configuration
type SloAlertConfigModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Severity        types.Int64  `tfsdk:"severity"`
	Triggering      types.Bool   `tfsdk:"triggering"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	AlertType       types.String `tfsdk:"alert_type"`
	Threshold       types.List   `tfsdk:"threshold"`
	SloIds          types.Set    `tfsdk:"slo_ids"`
	AlertChannelIds types.Set    `tfsdk:"alert_channel_ids"`
	TimeThreshold   types.List   `tfsdk:"time_threshold"`
	BurnRateConfig  types.List   `tfsdk:"burn_rate_config"`
	CustomPayload   types.List   `tfsdk:"custom_payload_fields"`
}

// SloAlertThresholdModel represents the threshold configuration for SLO Alert
type SloAlertThresholdModel struct {
	Type     types.String  `tfsdk:"type"`
	Operator types.String  `tfsdk:"operator"`
	Value    types.Float64 `tfsdk:"value"`
}

// SloAlertTimeThresholdModel represents the time threshold configuration for SLO Alert
type SloAlertTimeThresholdModel struct {
	WarmUp   types.Int64 `tfsdk:"warm_up"`
	CoolDown types.Int64 `tfsdk:"cool_down"`
}

// SloAlertBurnRateConfigModel represents the burn rate configuration for SLO Alert
type SloAlertBurnRateConfigModel struct {
	AlertWindowType   types.String `tfsdk:"alert_window_type"`
	Duration          types.String `tfsdk:"duration"`
	DurationUnitType  types.String `tfsdk:"duration_unit_type"`
	ThresholdOperator types.String `tfsdk:"threshold_operator"`
	ThresholdValue    types.String `tfsdk:"threshold_value"`
}

// SloAlertCustomPayloadFieldModel represents a custom payload field
type SloAlertCustomPayloadFieldModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
	Type  types.String `tfsdk:"type"`
}

type sloAlertConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

// NewSloAlertConfigResourceHandleFramework creates the resource handle for SLO Alert configuration
func NewSloAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.SloAlertConfig] {
	return &sloAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
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
				},
				Blocks: map[string]schema.Block{
					"threshold": schema.ListNestedBlock{
						Description: SloAlertConfigDescThreshold,
						NestedObject: schema.NestedBlockObject{
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
									Validators: []validator.Float64{
										float64validator.AtLeast(0.000001),
									},
								},
							},
						},
					},
					"time_threshold": schema.ListNestedBlock{
						Description: SloAlertConfigDescTimeThreshold,
						NestedObject: schema.NestedBlockObject{
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
					},
					"burn_rate_config": schema.ListNestedBlock{
						Description: SloAlertConfigDescBurnRateConfig,
						NestedObject: schema.NestedBlockObject{
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
					"custom_payload_fields": GetCustomPayloadFieldsSchema(),
				},
			},
			SchemaVersion: 1,
			CreateOnly:    false,
		},
	}
}

func (r *sloAlertConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
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

		thresholdModel := SloAlertThresholdModel{
			Type:     types.StringValue(thresholdType),
			Operator: types.StringValue(sloAlertConfig.Threshold.Operator),
			Value:    types.Float64Value(sloAlertConfig.Threshold.Value),
		}

		thresholdObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"type":     types.StringType,
			"operator": types.StringType,
			"value":    types.Float64Type,
		}, thresholdModel)
		if objDiags.HasError() {
			diags.Append(objDiags...)
			return diags
		}

		model.Threshold = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"operator": types.StringType,
				"value":    types.Float64Type,
			},
		}, []attr.Value{thresholdObj})
	} else {
		model.Threshold = types.ListNull(types.ObjectType{})
	}

	// Map time threshold
	timeThresholdModel := SloAlertTimeThresholdModel{
		WarmUp:   types.Int64Value(int64(sloAlertConfig.TimeThreshold.TimeWindow)),
		CoolDown: types.Int64Value(int64(sloAlertConfig.TimeThreshold.Expiry)),
	}

	timeThresholdObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"warm_up":   types.Int64Type,
		"cool_down": types.Int64Type,
	}, timeThresholdModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return diags
	}

	model.TimeThreshold = types.ListValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"warm_up":   types.Int64Type,
			"cool_down": types.Int64Type,
		},
	}, []attr.Value{timeThresholdObj})

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
		burnRateConfigs := []attr.Value{}
		for _, cfg := range *sloAlertConfig.BurnRateConfigs {
			burnRateConfigModel := SloAlertBurnRateConfigModel{
				AlertWindowType:   types.StringValue(cfg.AlertWindowType),
				Duration:          types.StringValue(fmt.Sprintf("%d", cfg.Duration)),
				DurationUnitType:  types.StringValue(cfg.DurationUnitType),
				ThresholdOperator: types.StringValue(cfg.Threshold.Operator),
				ThresholdValue:    types.StringValue(fmt.Sprintf("%.2f", cfg.Threshold.Value)),
			}

			burnRateConfigObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"alert_window_type":  types.StringType,
				"duration":           types.StringType,
				"duration_unit_type": types.StringType,
				"threshold_operator": types.StringType,
				"threshold_value":    types.StringType,
			}, burnRateConfigModel)
			if objDiags.HasError() {
				diags.Append(objDiags...)
				return diags
			}
			burnRateConfigs = append(burnRateConfigs, burnRateConfigObj)
		}

		model.BurnRateConfig = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"alert_window_type":  types.StringType,
				"duration":           types.StringType,
				"duration_unit_type": types.StringType,
				"threshold_operator": types.StringType,
				"threshold_value":    types.StringType,
			},
		}, burnRateConfigs)
	} else {
		model.BurnRateConfig = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"alert_window_type":  types.StringType,
				"duration":           types.StringType,
				"duration_unit_type": types.StringType,
				"threshold_operator": types.StringType,
				"threshold_value":    types.StringType,
			},
		})
	}

	// Map custom payload fields using the reusable function
	customPayloadFieldsList, payloadDiags := CustomPayloadFieldsToTerraform(ctx, sloAlertConfig.CustomerPayloadFields)
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
	if terraformAlertType != "burn_rate_v2" && !model.Threshold.IsNull() && !model.Threshold.IsUnknown() {
		var thresholdModels []SloAlertThresholdModel
		diags.Append(model.Threshold.ElementsAs(ctx, &thresholdModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(thresholdModels) > 0 {
			thresholdModel := thresholdModels[0]
			thresholdType := "staticThreshold"
			if !thresholdModel.Type.IsNull() {
				thresholdType = thresholdModel.Type.ValueString()
			}

			threshold = &restapi.SloAlertThreshold{
				Type:     thresholdType,
				Operator: thresholdModel.Operator.ValueString(),
				Value:    thresholdModel.Value.ValueFloat64(),
			}
		}
	}

	// Map time threshold
	var timeThreshold restapi.SloAlertTimeThreshold
	if !model.TimeThreshold.IsNull() && !model.TimeThreshold.IsUnknown() {
		var timeThresholdModels []SloAlertTimeThresholdModel
		diags.Append(model.TimeThreshold.ElementsAs(ctx, &timeThresholdModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(timeThresholdModels) > 0 {
			timeThresholdModel := timeThresholdModels[0]
			timeThreshold = restapi.SloAlertTimeThreshold{
				TimeWindow: int(timeThresholdModel.WarmUp.ValueInt64()),
				Expiry:     int(timeThresholdModel.CoolDown.ValueInt64()),
			}
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
	if terraformAlertType == "burn_rate_v2" && !model.BurnRateConfig.IsNull() && !model.BurnRateConfig.IsUnknown() {
		var burnRateConfigModels []SloAlertBurnRateConfigModel
		diags.Append(model.BurnRateConfig.ElementsAs(ctx, &burnRateConfigModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, burnRateConfigModel := range burnRateConfigModels {
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
		customPayloadFields, payloadDiags = MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayload)
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

// Made with Bob
