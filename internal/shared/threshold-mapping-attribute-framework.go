package shared

import (
	"context"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	LogAlertConfigFieldValue = "value"
)

type ThresholdPluginModel struct {
	Warning  *ThresholdTypeModel `tfsdk:"warning"`
	Critical *ThresholdTypeModel `tfsdk:"critical"`
}

type ThresholdTypeModel struct {
	Static           *StaticTypeModel       `tfsdk:"static"`
	AdaptiveBaseline *AdaptiveBaselineModel `tfsdk:"adaptive_baseline"`
}

type ThresholdStaticPluginModel struct {
	Warning  *ThresholdStaticTypeModel `tfsdk:"warning"`
	Critical *ThresholdStaticTypeModel `tfsdk:"critical"`
}

type ThresholdStaticTypeModel struct {
	Static *StaticTypeModel `tfsdk:"static"`
}

type ThresholdAllPluginModel struct {
	Warning  *ThresholdAllTypeModel `tfsdk:"warning"`
	Critical *ThresholdAllTypeModel `tfsdk:"critical"`
}

type ThresholdAllTypeModel struct {
	Static           *StaticTypeModel       `tfsdk:"static"`
	AdaptiveBaseline *AdaptiveBaselineModel `tfsdk:"adaptive_baseline"`
	HistoricBaseline *HistoricBaselineModel `tfsdk:"historic_baseline"`
}

func StaticAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Static threshold configuration",
		Attributes: map[string]schema.Attribute{
			"operator": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The operator for the static threshold",
			},
			LogAlertConfigFieldValue: schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The value of the threshold",
			},
		},
	}
}

func AdaptiveAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldAdaptiveBaselineDeviation: schema.Float32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The deviation factor for the adaptive baseline threshold",
			},
			ThresholdFieldAdaptiveBaselineAdaptability: schema.Float32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The adaptability for the adaptive baseline threshold",
			},
			ThresholdFieldAdaptiveBaselineSeasonality: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The seasonality for the adaptive baseline threshold",
			},
			"operator": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The operator for the adaptive baseline threshold",
			},
		},
	}
}

func HistoricAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldHistoricBaselineDeviation: schema.Float32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The deviation factor for the adaptive baseline threshold",
			},
			ThresholdFieldHistoricBaselineBaseline: schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Historic baseline as list of [timestamp, mean, sd] entries",
				ElementType: types.ListType{ElemType: types.Float64Type},
			},
			ThresholdFieldHistoricBaselineSeasonality: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The seasonality for the adaptive baseline threshold",
			},
		},
	}
}

func StaticAndAdaptiveThresholdAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldStatic:           StaticAttributeSchema(),
			ThresholdFieldAdaptiveBaseline: AdaptiveAttributeSchema(),
		},
	}
}

func StaticThresholdAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldStatic: StaticAttributeSchema(),
		},
	}
}

func AllThresholdAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldStatic:           StaticAttributeSchema(),
			ThresholdFieldAdaptiveBaseline: AdaptiveAttributeSchema(),
			ThresholdFieldHistoricBaseline: HistoricAttributeSchema(),
		},
	}
}

func MapThresholdsPluginFromState(ctx context.Context, thresholdStruct *ThresholdPluginModel) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)
	if thresholdStruct != nil {
		// Process warning threshold
		if thresholdStruct.Warning != nil {
			warningThreshold, warningDiags := MapThresholdRulePluginFromState(ctx, thresholdStruct.Warning)
			diags.Append(warningDiags...)
			if diags.HasError() {
				return nil, diags
			}

			if warningThreshold != nil {
				thresholdMap[restapi.WarningSeverity] = *warningThreshold
			}
		}

		// Process critical threshold
		if thresholdStruct.Critical != nil {
			criticalThreshold, criticalDiags := MapThresholdRulePluginFromState(ctx, thresholdStruct.Critical)
			diags.Append(criticalDiags...)
			if diags.HasError() {
				return nil, diags
			}

			if criticalThreshold != nil {
				thresholdMap[restapi.CriticalSeverity] = *criticalThreshold
			}
		}
	}
	return thresholdMap, diags
}

// MapThresholdRulePluginFromState maps a threshold rule from Terraform state to API model - used for new plugin model using nestedAttribute
func MapThresholdRulePluginFromState(ctx context.Context, thresholdObj *ThresholdTypeModel) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if thresholdObj == nil {
		return nil, diags
	}

	// Check for static threshold
	if thresholdObj.Static != nil {
		staticVal := thresholdObj.Static
		valueFloat := float64(staticVal.Value.ValueInt64())
		return &restapi.ThresholdRule{
			Type:     "staticThreshold",
			Value:    &valueFloat,
			Operator: staticVal.Operator.ValueStringPointer(),
		}, diags
	}

	// Check for adaptive baseline threshold
	if thresholdObj.AdaptiveBaseline != nil {
		adaptiveVal := thresholdObj.AdaptiveBaseline
		seasonality := restapi.ThresholdSeasonality(adaptiveVal.Seasonality.ValueString())
		deviationFactor := float32(adaptiveVal.DeviationFactor.ValueFloat32())
		adaptability := adaptiveVal.Adaptability.ValueFloat32()
		operator := util.SetStringPointerFromState(adaptiveVal.Operator)
		return &restapi.ThresholdRule{
			Type:            "adaptiveBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Adaptability:    &adaptability,
			Operator:        operator,
		}, diags
	}
	return nil, diags
}

func MapThresholdsAllPluginFromState(ctx context.Context, thresholdStruct *ThresholdAllPluginModel) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)
	if thresholdStruct != nil {
		// Process warning threshold
		if thresholdStruct.Warning != nil {
			log.Printf("inside warning thresh")

			warningThreshold, warningDiags := MapThresholdRuleAllPluginFromState(ctx, thresholdStruct.Warning)
			diags.Append(warningDiags...)
			if diags.HasError() {
				return nil, diags
			}

			if warningThreshold != nil {
				thresholdMap[restapi.WarningSeverity] = *warningThreshold
			}
		}

		// Process critical threshold
		if thresholdStruct.Critical != nil {
			log.Printf("inside warning thresh")

			criticalThreshold, criticalDiags := MapThresholdRuleAllPluginFromState(ctx, thresholdStruct.Critical)
			diags.Append(criticalDiags...)
			if diags.HasError() {
				return nil, diags
			}

			if criticalThreshold != nil {
				thresholdMap[restapi.CriticalSeverity] = *criticalThreshold
			}
		}
	}
	return thresholdMap, diags
}

// MapThresholdRulePluginFromState maps a threshold rule from Terraform state to API model - used for new plugin model using nestedAttribute
func MapThresholdRuleAllPluginFromState(ctx context.Context, thresholdObj *ThresholdAllTypeModel) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if thresholdObj == nil {
		return nil, diags
	}

	// Check for static threshold
	if thresholdObj.Static != nil {
		log.Printf("inside static thresh")

		staticVal := thresholdObj.Static
		valueFloat := float64(staticVal.Value.ValueInt64())
		return &restapi.ThresholdRule{
			Type:  "staticThreshold",
			Value: &valueFloat,
		}, diags
	}

	// Check for adaptive baseline threshold
	if thresholdObj.AdaptiveBaseline != nil {
		log.Printf("inside adaptive thresh")

		adaptiveVal := thresholdObj.AdaptiveBaseline
		seasonality := restapi.ThresholdSeasonality(adaptiveVal.Seasonality.ValueString())
		deviationFactor := float32(adaptiveVal.DeviationFactor.ValueFloat32())
		adaptability := adaptiveVal.Adaptability.ValueFloat32()
		operator := util.SetStringPointerFromState(adaptiveVal.Operator)
		return &restapi.ThresholdRule{
			Type:            "adaptiveBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Adaptability:    &adaptability,
			Operator:        operator,
		}, diags
	}

	// Check for historic baseline threshold
	if thresholdObj.HistoricBaseline != nil {
		log.Printf("inside histori thresh")

		baselineVal := thresholdObj.HistoricBaseline
		seasonality := restapi.ThresholdSeasonality(baselineVal.Seasonality.ValueString())
		deviationFactor := float32(baselineVal.Deviation.ValueFloat32())
		baseline, _ := MapBaselineFromState(ctx, baselineVal.Baseline)
		return &restapi.ThresholdRule{
			Type:            "historicBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Baseline:        baseline,
		}, diags
	}
	return nil, diags
}

// mapBaselineFromState converts a Terraform List of List (baseline data) to API format
func MapBaselineFromState(ctx context.Context, baselineList types.List) (*[][]float64, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get the outer list elements (each element is itself a list of float64)
	var outerListElements []types.List
	diags.Append(baselineList.ElementsAs(ctx, &outerListElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert to [][]float64
	baseline := make([][]float64, 0, len(outerListElements))
	for _, innerList := range outerListElements {
		// Get the inner list elements (float64 values)
		var innerValues []float64
		diags.Append(innerList.ElementsAs(ctx, &innerValues, false)...)
		if diags.HasError() {
			return nil, diags
		}
		baseline = append(baseline, innerValues)
	}

	return &baseline, diags
}

// MapThresholdToState maps a threshold rule to a Terraform state representation - used for nested attribute instead of block object
func MapThresholdPluginToState(ctx context.Context, threshold *restapi.ThresholdRule, dataPresent bool) *ThresholdTypeModel {
	thresholdTypeModel := ThresholdTypeModel{}
	if dataPresent == false {
		return nil
	}
	switch threshold.Type {
	case "adaptiveBaseline":
		adaptiveBaselineModel := AdaptiveBaselineModel{
			Operator:        util.SetStringPointerToState(threshold.Operator),
			DeviationFactor: util.SetFloat32PointerToState(threshold.DeviationFactor),
			Adaptability:    util.SetFloat32PointerToState(threshold.Adaptability),
			Seasonality:     types.StringValue(string(*threshold.Seasonality)),
		}
		thresholdTypeModel.AdaptiveBaseline = &adaptiveBaselineModel
	default:
		// Default to static threshold for all other types
		static := StaticTypeModel{
			Operator: util.SetStringPointerToState(threshold.Operator),
			Value:    util.SetInt64PointerToState(threshold.Value),
		}
		thresholdTypeModel.Static = &static
	}

	return &thresholdTypeModel
}

// MapThresholdToState maps a threshold rule to a Terraform state representation - used for nested attribute instead of block object
func MapAllThresholdPluginToState(ctx context.Context, threshold *restapi.ThresholdRule, dataPresent bool) *ThresholdAllTypeModel {
	thresholdTypeModel := ThresholdAllTypeModel{}
	if dataPresent == false {
		return nil
	}
	switch threshold.Type {
	case "adaptiveBaseline":
		adaptiveBaselineModel := AdaptiveBaselineModel{
			Operator:        util.SetStringPointerToState(threshold.Operator),
			DeviationFactor: util.SetFloat32PointerToState(threshold.DeviationFactor),
			Adaptability:    util.SetFloat32PointerToState(threshold.Adaptability),
			Seasonality:     types.StringValue(string(*threshold.Seasonality)),
		}
		thresholdTypeModel.AdaptiveBaseline = &adaptiveBaselineModel
	case "historicBaseline":
		historicBaselineModel := HistoricBaselineModel{
			//Baseline:    threshold.Baseline,
			Deviation:   util.SetFloat32PointerToState(threshold.DeviationFactor),
			Seasonality: types.StringValue(string(*threshold.Seasonality)),
		}
		historicBaselineModel.Baseline, _ = MapBaseline(threshold)
		thresholdTypeModel.HistoricBaseline = &historicBaselineModel
	default:
		// Default to static threshold for all other types
		static := StaticTypeModel{
			Operator: util.SetStringPointerToState(threshold.Operator),
			Value:    util.SetInt64PointerToState(threshold.Value),
		}
		thresholdTypeModel.Static = &static
	}

	return &thresholdTypeModel
}
