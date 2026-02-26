package shared

import (
	"context"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/util"
)

const (
	LogAlertConfigFieldValue = "value"
	//ResourceFieldThreshold constant value for field threshold
	ResourceFieldThreshold = "threshold"
	//ResourceFieldThresholdLastUpdated constant value for field threshold.*.last_updated
	ResourceFieldThresholdLastUpdated = "last_updated"
	//ResourceFieldThresholdOperator constant value for field threshold.*.operator
	ResourceFieldThresholdOperator = "operator"
	//ResourceFieldThresholdHistoricBaseline constant value for field threshold.historic_baseline
	ResourceFieldThresholdHistoricBaseline = "historic_baseline"
	//ResourceFieldThresholdHistoricBaselineBaseline constant value for field threshold.historic_baseline.baseline
	ResourceFieldThresholdHistoricBaselineBaseline = "baseline"
	//ResourceFieldThresholdHistoricBaselineDeviationFactor constant value for field threshold.historic_baseline.deviation_factor
	ResourceFieldThresholdHistoricBaselineDeviationFactor = "deviation_factor"
	//ResourceFieldThresholdHistoricBaselineSeasonality constant value for field threshold.historic_baseline.seasonality
	ResourceFieldThresholdHistoricBaselineSeasonality = "seasonality"
	//ResourceFieldThresholdStatic constant value for field threshold.static
	ResourceFieldThresholdStatic = "static"
	//ResourceFieldThresholdStaticValue constant value for field threshold.static.value
	ResourceFieldThresholdStaticValue = "value"
	LogAlertConfigFieldWarning        = "warning"
	LogAlertConfigFieldCritical       = "critical"

	// Constants for threshold field names related to historic baseline
	ThresholdFieldHistoricBaseline            = "historic_baseline"
	ThresholdFieldHistoricBaselineBaseline    = "baseline"
	ThresholdFieldHistoricBaselineDeviation   = "deviation_factor"
	ThresholdFieldHistoricBaselineSeasonality = "seasonality"

	// Constants for threshold field names related to adaptive baseline
	ThresholdFieldAdaptiveBaseline             = "adaptive_baseline"
	ThresholdFieldAdaptiveBaselineDeviation    = "deviation_factor"
	ThresholdFieldAdaptiveBaselineAdaptability = "adaptability"
	ThresholdFieldAdaptiveBaselineSeasonality  = "seasonality"

	// Constants for threshold field names
	ThresholdFieldWarning  = "warning"
	ThresholdFieldCritical = "critical"
	ThresholdFieldStatic   = "static"
)

type AdaptiveBaselineModel struct {
	//Operator        types.String  `tfsdk:"operator"`
	DeviationFactor types.Float64 `tfsdk:"deviation_factor"`
	Adaptability    types.Float64 `tfsdk:"adaptability"`
	Seasonality     types.String  `tfsdk:"seasonality"`
}

type HistoricBaselineModel struct {
	//Operator    types.String  `tfsdk:"operator"`
	Baseline    types.List    `tfsdk:"baseline"`
	Deviation   types.Float64 `tfsdk:"deviation_factor"`
	Seasonality types.String  `tfsdk:"seasonality"`
}

type StaticTypeModel struct {
	//Operator types.String  `tfsdk:"operator"`
	Value types.Float64 `tfsdk:"value"`
}
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
			LogAlertConfigFieldValue: schema.Float64Attribute{
				Optional:    true,
				Description: "The value of the threshold",
			},
			// ResourceFieldThresholdOperator: schema.StringAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "The operator for the adaptive baseline threshold",
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf(">", ">=", "<", "<="),
			// 	},
			// },
		},
	}
}

func AdaptiveAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldAdaptiveBaselineDeviation: schema.Float64Attribute{
				Optional:    true,
				Description: "The deviation factor for the adaptive baseline threshold",
			},
			ThresholdFieldAdaptiveBaselineAdaptability: schema.Float64Attribute{
				Optional:    true,
				Description: "The adaptability for the adaptive baseline threshold",
			},
			ThresholdFieldAdaptiveBaselineSeasonality: schema.StringAttribute{
				Optional:    true,
				Description: "The seasonality for the adaptive baseline threshold",
			},
			// ResourceFieldThresholdOperator: schema.StringAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "The operator for the adaptive baseline threshold",
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf(">", ">=", "<", "<="),
			// 	},
			// },
		},
	}
}

func HistoricAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			ThresholdFieldHistoricBaselineDeviation: schema.Float64Attribute{
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
			// ResourceFieldThresholdOperator: schema.StringAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "The operator for the adaptive baseline threshold",
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf(">", ">=", "<", "<="),
			// 	},
			// },
		},
	}
}

func StaticAndAdaptiveThresholdAttributeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Threshold configuration",
		Optional:    true,
		//Computed:    true,
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

// used for static and adaptive types
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
		valueFloat := staticVal.Value.ValueFloat64()
		rounded := math.Round(valueFloat*100) / 100
		return &restapi.ThresholdRule{
			Type:  "staticThreshold",
			Value: &rounded,
			//Operator: staticVal.Operator.ValueStringPointer(),
		}, diags
	}

	// Check for adaptive baseline threshold
	if thresholdObj.AdaptiveBaseline != nil {
		adaptiveVal := thresholdObj.AdaptiveBaseline
		seasonality := restapi.ThresholdSeasonality(adaptiveVal.Seasonality.ValueString())
		// Round to 2 decimal places to avoid floating-point precision issues
		deviationFactorVal := adaptiveVal.DeviationFactor.ValueFloat64()
		formatted := strconv.FormatFloat(deviationFactorVal, 'f', 2, 64)
		parsed, _ := strconv.ParseFloat(formatted, 64)
		deviationFactor := float32(parsed)

		adaptabilityVal := adaptiveVal.Adaptability.ValueFloat64()
		formattedAdapt := strconv.FormatFloat(adaptabilityVal, 'f', 2, 64)
		parsedAdapt, _ := strconv.ParseFloat(formattedAdapt, 64)
		adaptability := float32(parsedAdapt)
		//operator := util.SetStringPointerFromState(adaptiveVal.Operator)
		return &restapi.ThresholdRule{
			Type:            "adaptiveBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Adaptability:    &adaptability,
			//Operator:        operator,
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

// used for static, adaptive and historic baseline types
// MapThresholdRulePluginFromState maps a threshold rule from Terraform state to API model - used for new plugin model using nestedAttribute
func MapThresholdRuleAllPluginFromState(ctx context.Context, thresholdObj *ThresholdAllTypeModel) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if thresholdObj == nil {
		return nil, diags
	}

	// Check for static threshold
	if thresholdObj.Static != nil {

		staticVal := thresholdObj.Static
		//operator := util.SetStringPointerFromState(staticVal.Operator)

		valueFloat := staticVal.Value.ValueFloat64()
		rounded := math.Round(valueFloat*100) / 100
		return &restapi.ThresholdRule{
			Type:  "staticThreshold",
			Value: &rounded,
			//Operator: operator,
		}, diags
	}

	// Check for adaptive baseline threshold
	if thresholdObj.AdaptiveBaseline != nil {

		adaptiveVal := thresholdObj.AdaptiveBaseline
		seasonality := restapi.ThresholdSeasonality(adaptiveVal.Seasonality.ValueString())
		// Round to 2 decimal places to avoid floating-point precision issues
		deviationFactorVal := adaptiveVal.DeviationFactor.ValueFloat64()
		formatted := strconv.FormatFloat(deviationFactorVal, 'f', 2, 64)
		parsed, _ := strconv.ParseFloat(formatted, 64)
		deviationFactor := float32(parsed)

		adaptabilityVal := adaptiveVal.Adaptability.ValueFloat64()
		formattedAdapt := strconv.FormatFloat(adaptabilityVal, 'f', 2, 64)
		parsedAdapt, _ := strconv.ParseFloat(formattedAdapt, 64)
		adaptability := float32(parsedAdapt)
		//operator := util.SetStringPointerFromState(adaptiveVal.Operator)
		return &restapi.ThresholdRule{
			Type:            "adaptiveBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Adaptability:    &adaptability,
			//Operator:        operator,
		}, diags
	}

	// Check for historic baseline threshold
	if thresholdObj.HistoricBaseline != nil {

		baselineVal := thresholdObj.HistoricBaseline
		//operator := util.SetStringPointerFromState(baselineVal.Operator)
		seasonality := restapi.ThresholdSeasonality(baselineVal.Seasonality.ValueString())
		// Round to 2 decimal places to avoid floating-point precision issues
		deviationFactorVal := baselineVal.Deviation.ValueFloat64()
		formatted := strconv.FormatFloat(deviationFactorVal, 'f', 2, 64)
		parsed, _ := strconv.ParseFloat(formatted, 64)
		deviationFactor := float32(parsed)

		// Convert types.List to [][]float64
		var baseline *[][]float64
		if !baselineVal.Baseline.IsNull() && !baselineVal.Baseline.IsUnknown() {
			baseline, diags = MapBaselineFromState(ctx, baselineVal.Baseline)
			if diags.HasError() {
				return nil, diags
			}
		}

		return &restapi.ThresholdRule{
			Type:            "historicBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Baseline:        baseline,
			//Operator:        operator,
		}, diags
	}
	return nil, diags
}

func MapBaselineToState(threshold *restapi.ThresholdRule) (basetypes.ListValue, diag.Diagnostics) {
	var baselineDiags diag.Diagnostics
	baselineListValues := []attr.Value{}
	for _, baselineArray := range *threshold.Baseline {
		innerListValues := []attr.Value{}
		for _, value := range baselineArray {
			// Round to 4 decimal places to avoid floating-point precision issues
			rounded := math.Round(value*10000) / 10000
			innerListValues = append(innerListValues, types.Float64Value(rounded))
		}

		innerSet, innerSetDiags := types.ListValue(types.Float64Type, innerListValues)
		baselineDiags.Append(innerSetDiags...)
		if baselineDiags.HasError() {
			return basetypes.ListValue{}, baselineDiags
		}

		baselineListValues = append(baselineListValues, innerSet)
	}

	baselineSet, baselineSetDiags := types.ListValue(
		types.ListType{ElemType: types.Float64Type},
		baselineListValues,
	)
	baselineDiags.Append(baselineSetDiags...)
	if baselineDiags.HasError() {
		return basetypes.ListValue{}, baselineDiags
	}
	return baselineSet, baselineDiags
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
		// Use the exact values from state without rounding
		// Rounding should only happen when converting FROM API TO state
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
		var seasonality types.String
		if threshold.Seasonality != nil {
			seasonality = types.StringValue(string(*threshold.Seasonality))
		} else {
			seasonality = types.StringNull()
		}
		adaptiveBaselineModel := AdaptiveBaselineModel{
			//Operator:        util.SetStringPointerToState(threshold.Operator),
			DeviationFactor: util.SetFloat64PointerToState(threshold.DeviationFactor),
			Adaptability:    util.SetFloat64PointerToState(threshold.Adaptability),
			Seasonality:     seasonality,
		}
		thresholdTypeModel.AdaptiveBaseline = &adaptiveBaselineModel
	default:
		// Default to static threshold for all other types
		// Round to 2 decimal places to avoid floating-point precision issues
		var roundedValue types.Float64
		if threshold.Value != nil {
			// Format to 2 decimal places and parse back to ensure exact precision
			formatted := strconv.FormatFloat(*threshold.Value, 'f', 2, 64)
			parsed, _ := strconv.ParseFloat(formatted, 64)
			roundedValue = types.Float64Value(parsed)
		} else {
			roundedValue = types.Float64Null()
		}
		static := StaticTypeModel{
			//Operator: util.SetStringPointerToState(threshold.Operator),
			Value: roundedValue,
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
		var seasonality types.String
		if threshold.Seasonality != nil {
			seasonality = types.StringValue(string(*threshold.Seasonality))
		} else {
			seasonality = types.StringNull()
		}
		adaptiveBaselineModel := AdaptiveBaselineModel{
			//Operator:        util.SetStringPointerToState(threshold.Operator),
			DeviationFactor: util.SetFloat64PointerToState(threshold.DeviationFactor),
			Adaptability:    util.SetFloat64PointerToState(threshold.Adaptability),
			Seasonality:     seasonality,
		}
		thresholdTypeModel.AdaptiveBaseline = &adaptiveBaselineModel
	case "historicBaseline":
		var baselineList basetypes.ListValue
		var baselineDiags diag.Diagnostics
		if threshold.Baseline != nil {
			baselineList, baselineDiags = MapBaselineToState(threshold)
			if baselineDiags.HasError() {
				// If there's an error, set to null list
				baselineList = types.ListNull(types.ListType{ElemType: types.Float64Type})
			}
		} else {
			baselineList = types.ListNull(types.ListType{ElemType: types.Float64Type})
		}

		var seasonality types.String
		if threshold.Seasonality != nil {
			seasonality = types.StringValue(string(*threshold.Seasonality))
		} else {
			seasonality = types.StringNull()
		}
		historicBaselineModel := HistoricBaselineModel{
			Baseline:    baselineList,
			Deviation:   util.SetFloat64PointerToState(threshold.DeviationFactor),
			Seasonality: seasonality,
		}
		thresholdTypeModel.HistoricBaseline = &historicBaselineModel
	default:
		// Default to static threshold for all other types
		// Round to 2 decimal places to avoid floating-point precision issues
		var roundedValue types.Float64
		if threshold.Value != nil {
			// Format to 2 decimal places and parse back to ensure exact precision
			formatted := strconv.FormatFloat(*threshold.Value, 'f', 2, 64)
			parsed, _ := strconv.ParseFloat(formatted, 64)
			roundedValue = types.Float64Value(parsed)
		} else {
			roundedValue = types.Float64Null()
		}
		static := StaticTypeModel{
			//Operator: util.SetStringPointerToState(threshold.Operator),
			Value: roundedValue,
		}
		thresholdTypeModel.Static = &static
	}

	return &thresholdTypeModel
}
