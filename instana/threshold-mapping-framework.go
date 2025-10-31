package instana

import (
	"context"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Constants for threshold field names
const (
	ThresholdFieldWarning  = "warning"
	ThresholdFieldCritical = "critical"
	ThresholdFieldStatic   = "static"
)

type AdaptiveBaselineModel struct {
	Operator        types.String  `tfsdk:"operator"`
	DeviationFactor types.Float32 `tfsdk:"deviation_factor"`
	Adaptability    types.Float32 `tfsdk:"adaptability"`
	Seasonality     types.String  `tfsdk:"seasonality"`
}

type StaticTypeModel struct {
	Operator types.String `tfsdk:"operator"`
	Value    types.Int64  `tfsdk:"value"`
}

type ThresholdPluginModel struct {
	Warning  *ThresholdTypeModel `tfsdk:"warning"`
	Critical *ThresholdTypeModel `tfsdk:"critical"`
}
type ThresholdTypeModel struct {
	Static           *StaticTypeModel       `tfsdk:"static"`
	AdaptiveBaseline *AdaptiveBaselineModel `tfsdk:"adaptive_baseline"`
}

// StaticThresholdBlockSchema returns the schema for static block configuration
func StaticBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Static threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				LogAlertConfigFieldValue: schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "The value of the threshold",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeBetween(0, 1),
		},
	}
}

func StaticAttributeSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "Static threshold configuration",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				LogAlertConfigFieldValue: schema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "The value of the threshold",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeBetween(0, 1),
		},
	}
}

func AdaptiveAttributeSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "Threshold configuration",
		NestedObject: schema.NestedAttributeObject{
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
			},
		},
	}
}

func AdaptiveBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Threshold configuration",
		NestedObject: schema.NestedBlockObject{
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
			},
		},
	}
}

// StaticThresholdBlockSchema returns the schema for static threshold configuration
func StaticThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Warning threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ThresholdFieldStatic: StaticBlockSchema(),
			},
		},
	}
}

// define a static and adaptive schema for threshold configuration
func StaticAndAdaptiveThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ThresholdFieldStatic:           StaticBlockSchema(),
				ThresholdFieldAdaptiveBaseline: AdaptiveBlockSchema(),
			},
		},
	}
}

func StaticAndAdaptiveThresholdAttributeSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: "Threshold configuration",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				ThresholdFieldStatic:           StaticAttributeSchema(),
				ThresholdFieldAdaptiveBaseline: AdaptiveAttributeSchema(),
			},
		},
	}
}

// Constants for threshold field names related to historic baseline
const (
	ThresholdFieldHistoricBaseline            = "historic_baseline"
	ThresholdFieldHistoricBaselineBaseline    = "baseline"
	ThresholdFieldHistoricBaselineDeviation   = "deviation_factor"
	ThresholdFieldHistoricBaselineSeasonality = "seasonality"
)

// Constants for threshold field names related to adaptive baseline
const (
	ThresholdFieldAdaptiveBaseline             = "adaptive_baseline"
	ThresholdFieldAdaptiveBaselineDeviation    = "deviation_factor"
	ThresholdFieldAdaptiveBaselineAdaptability = "adaptability"
	ThresholdFieldAdaptiveBaselineSeasonality  = "seasonality"
)

// HistoricBaselineBlockSchema returns the schema for historic baseline configuration
func HistoricBaselineBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Historic baseline threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				ThresholdFieldHistoricBaselineDeviation: schema.Float32Attribute{
					Optional:    true,
					Description: "The deviation factor for the historic baseline threshold",
				},
				ThresholdFieldHistoricBaselineSeasonality: schema.StringAttribute{
					Required:    true,
					Description: "The seasonality of the historic baseline threshold (DAILY or WEEKLY)",
				},
				ThresholdFieldHistoricBaselineBaseline: schema.SetAttribute{
					Optional:    true,
					Description: "The baseline of the historic baseline threshold",
					ElementType: types.SetType{
						ElemType: types.Float64Type,
					},
				},
			},
		},
	}
}

// AdaptiveBaselineBlockSchema returns the schema for adaptive baseline configuration
func AdaptiveBaselineBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Adaptive baseline threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				ThresholdFieldAdaptiveBaselineDeviation: schema.Float32Attribute{
					Required:    true,
					Description: "The deviation factor for the adaptive baseline threshold",
				},
				ThresholdFieldAdaptiveBaselineAdaptability: schema.Float32Attribute{
					Required:    true,
					Description: "The adaptability factor for the adaptive baseline threshold",
				},
				ThresholdFieldAdaptiveBaselineSeasonality: schema.StringAttribute{
					Required:    true,
					Description: "The seasonality of the adaptive baseline threshold (DAILY or WEEKLY)",
				},
			},
		},
	}
}

// HistoricThresholdBlockSchema returns the schema for historic threshold configuration
func HistoricThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Historic baseline threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ThresholdFieldHistoricBaseline: HistoricBaselineBlockSchema(),
			},
		},
	}
}

// AdaptiveThresholdBlockSchema returns the schema for adaptive threshold configuration
func AdaptiveThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Adaptive baseline threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ThresholdFieldAdaptiveBaseline: AdaptiveBaselineBlockSchema(),
			},
		},
	}
}

// MapThresholdToState maps a threshold rule to a Terraform state representation
func MapThresholdToState(ctx context.Context, isThresholdPresent bool, threshold *restapi.ThresholdRule, expectedThresholdTypes []string) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Initialize attribute types map with all expected threshold types
	thresholdAttrTypes := map[string]attr.Type{}
	// Create threshold object based on type
	thresholdObj := map[string]attr.Value{}

	// Add expected threshold types to the attribute types map
	for _, thresholdType := range expectedThresholdTypes {
		switch thresholdType {
		case "static":
			thresholdAttrTypes[ThresholdFieldStatic] = types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						LogAlertConfigFieldValue: types.Int64Type,
					},
				},
			}
			thresholdObj[ThresholdFieldStatic] = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					LogAlertConfigFieldValue: types.Int64Type,
				},
			})
		case "historicBaseline":
			thresholdAttrTypes[ThresholdFieldHistoricBaseline] = types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						ThresholdFieldHistoricBaselineBaseline:    types.SetType{ElemType: types.SetType{ElemType: types.Float64Type}},
						ThresholdFieldHistoricBaselineDeviation:   types.Float32Type,
						ThresholdFieldHistoricBaselineSeasonality: types.StringType,
					},
				},
			}
			thresholdObj[ThresholdFieldHistoricBaseline] = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					ThresholdFieldHistoricBaselineBaseline:    types.SetType{ElemType: types.SetType{ElemType: types.Float64Type}},
					ThresholdFieldHistoricBaselineDeviation:   types.Float32Type,
					ThresholdFieldHistoricBaselineSeasonality: types.StringType,
				},
			})
		case "adaptiveBaseline":
			thresholdAttrTypes[ThresholdFieldAdaptiveBaseline] = types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
						ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
						ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
					},
				},
			}
			thresholdObj[ThresholdFieldAdaptiveBaseline] = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
					ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
					ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
				},
			})
		}
	}

	if !isThresholdPresent || threshold == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: thresholdAttrTypes,
		}), diags
	}

	// Initialize all expected threshold types with null values
	for _, thresholdType := range expectedThresholdTypes {
		switch thresholdType {
		case "static":
			if threshold.Type != "staticThreshold" {
				thresholdObj[ThresholdFieldStatic] = types.ListNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						LogAlertConfigFieldValue: types.Int64Type,
					},
				})
			}
		case "historicBaseline":
			if threshold.Type != "historicBaseline" {
				thresholdObj[ThresholdFieldHistoricBaseline] = types.ListNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						ThresholdFieldHistoricBaselineBaseline:    types.SetType{ElemType: types.SetType{ElemType: types.Float64Type}},
						ThresholdFieldHistoricBaselineDeviation:   types.Float32Type,
						ThresholdFieldHistoricBaselineSeasonality: types.StringType,
					},
				})
			}
		case "adaptiveBaseline":
			if threshold.Type != "adaptiveBaseline" {
				thresholdObj[ThresholdFieldAdaptiveBaseline] = types.ListNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
						ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
						ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
					},
				})
			}
		}
	}

	switch threshold.Type {
	case "historicBaseline":
		// Handle historic baseline
		historicList, historicDiags := mapHistoricBaselineToState(ctx, threshold)
		diags.Append(historicDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		thresholdObj[ThresholdFieldHistoricBaseline] = historicList
		thresholdAttrTypes[ThresholdFieldHistoricBaseline] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					ThresholdFieldHistoricBaselineBaseline:    types.SetType{ElemType: types.SetType{ElemType: types.Float64Type}},
					ThresholdFieldHistoricBaselineDeviation:   types.Float32Type,
					ThresholdFieldHistoricBaselineSeasonality: types.StringType,
				},
			},
		}

	case "adaptiveBaseline":
		// Handle adaptive baseline
		adaptiveList, adaptiveDiags := mapAdaptiveBaselineToState(ctx, threshold)
		diags.Append(adaptiveDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		thresholdObj[ThresholdFieldAdaptiveBaseline] = adaptiveList
		thresholdAttrTypes[ThresholdFieldAdaptiveBaseline] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
					ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
					ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
				},
			},
		}
	default:
		// Default to static threshold for all other types
		staticList, staticDiags := mapStaticThresholdToState(ctx, threshold)
		diags.Append(staticDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		thresholdObj[ThresholdFieldStatic] = staticList
		thresholdAttrTypes[ThresholdFieldStatic] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					LogAlertConfigFieldValue: types.Int64Type,
				},
			},
		}
	}

	// Create threshold object value
	thresholdObjVal, thresholdObjDiags := types.ObjectValue(
		thresholdAttrTypes,
		thresholdObj,
	)
	diags.Append(thresholdObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: thresholdAttrTypes,
		},
		[]attr.Value{thresholdObjVal},
	)
}

// mapStaticThresholdToState maps a static threshold to Terraform state
func mapStaticThresholdToState(ctx context.Context, threshold *restapi.ThresholdRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get the threshold value
	var thresholdValue int64
	if threshold.Value != nil {
		thresholdValue = int64(*threshold.Value)
	} else {
		thresholdValue = 0
	}

	// Create static threshold object
	staticObj := map[string]attr.Value{
		LogAlertConfigFieldValue: types.Int64Value(thresholdValue),
	}

	staticObjVal, staticObjDiags := types.ObjectValue(
		map[string]attr.Type{
			LogAlertConfigFieldValue: types.Int64Type,
		},
		staticObj,
	)
	diags.Append(staticObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create static list
	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				LogAlertConfigFieldValue: types.Int64Type,
			},
		},
		[]attr.Value{staticObjVal},
	)
}

// mapHistoricBaselineToState maps a historic baseline threshold to Terraform state
func mapHistoricBaselineToState(ctx context.Context, threshold *restapi.ThresholdRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create historic baseline object
	historicObj := map[string]attr.Value{}

	// Map seasonality
	if threshold.Seasonality != nil {
		historicObj[ThresholdFieldHistoricBaselineSeasonality] = types.StringValue(string(*threshold.Seasonality))
	} else {
		historicObj[ThresholdFieldHistoricBaselineSeasonality] = types.StringValue(string(restapi.ThresholdSeasonalityDaily))
	}

	// Map deviation factor
	if threshold.DeviationFactor != nil {
		historicObj[ThresholdFieldHistoricBaselineDeviation] = types.Float32Value(float32(*threshold.DeviationFactor))
	} else {
		historicObj[ThresholdFieldHistoricBaselineDeviation] = types.Float32Value(1.0)
	}

	// Map baseline
	if threshold.Baseline != nil {
		baselineSetValues := []attr.Value{}
		for _, baselineArray := range *threshold.Baseline {
			innerSetValues := []attr.Value{}
			for _, value := range baselineArray {
				innerSetValues = append(innerSetValues, types.Float64Value(value))
			}

			innerSet, innerSetDiags := types.SetValue(types.Float64Type, innerSetValues)
			diags.Append(innerSetDiags...)
			if diags.HasError() {
				return types.ListNull(types.ObjectType{}), diags
			}

			baselineSetValues = append(baselineSetValues, innerSet)
		}

		baselineSet, baselineSetDiags := types.SetValue(
			types.SetType{ElemType: types.Float64Type},
			baselineSetValues,
		)
		diags.Append(baselineSetDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		historicObj[ThresholdFieldHistoricBaselineBaseline] = baselineSet
	} else {
		historicObj[ThresholdFieldHistoricBaselineBaseline] = types.SetNull(types.SetType{ElemType: types.Float64Type})
	}

	// Create historic baseline object value
	historicObjVal, historicObjDiags := types.ObjectValue(
		map[string]attr.Type{
			ThresholdFieldHistoricBaselineBaseline:    types.SetType{ElemType: types.SetType{ElemType: types.Float64Type}},
			ThresholdFieldHistoricBaselineDeviation:   types.Float32Type,
			ThresholdFieldHistoricBaselineSeasonality: types.StringType,
		},
		historicObj,
	)
	diags.Append(historicObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create historic baseline list
	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				ThresholdFieldHistoricBaselineBaseline:    types.SetType{ElemType: types.SetType{ElemType: types.Float64Type}},
				ThresholdFieldHistoricBaselineDeviation:   types.Float32Type,
				ThresholdFieldHistoricBaselineSeasonality: types.StringType,
			},
		},
		[]attr.Value{historicObjVal},
	)
}

// mapAdaptiveBaselineToState maps an adaptive baseline threshold to Terraform state
func mapAdaptiveBaselineToState(ctx context.Context, threshold *restapi.ThresholdRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create adaptive baseline object
	adaptiveObj := map[string]attr.Value{}

	// Map seasonality
	if threshold.Seasonality != nil {
		adaptiveObj[ThresholdFieldAdaptiveBaselineSeasonality] = types.StringValue(string(*threshold.Seasonality))
	} else {
		adaptiveObj[ThresholdFieldAdaptiveBaselineSeasonality] = types.StringValue(string(restapi.ThresholdSeasonalityDaily))
	}

	// Map deviation factor
	if threshold.DeviationFactor != nil {
		adaptiveObj[ThresholdFieldAdaptiveBaselineDeviation] = types.Float32Value(float32(*threshold.DeviationFactor))
	} else {
		adaptiveObj[ThresholdFieldAdaptiveBaselineDeviation] = types.Float32Value(1.0)
	}

	adaptiveObj[ThresholdFieldAdaptiveBaselineAdaptability] = types.Float32Value(*threshold.Adaptability)

	// Create adaptive baseline object value
	adaptiveObjVal, adaptiveObjDiags := types.ObjectValue(
		map[string]attr.Type{
			ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
			ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
			ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
		},
		adaptiveObj,
	)
	diags.Append(adaptiveObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create adaptive baseline list
	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
				ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
				ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
			},
		},
		[]attr.Value{adaptiveObjVal},
	)
}

// GetStaticThresholdAttrTypes returns the attribute types map for static thresholds
func GetStaticThresholdAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		ThresholdFieldStatic: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					LogAlertConfigFieldValue: types.Int64Type,
				},
			},
		},
	}
}

func GetStaticThresholdAttrListTypes() types.ListType {
	return types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: GetStaticThresholdAttrTypes(),
		},
	}
}

// GetStaticAndAdaptiveThresholdAttrTypes returns the attribute types map for both static and adaptive thresholds
func GetStaticAndAdaptiveThresholdAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		ThresholdFieldStatic: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					LogAlertConfigFieldValue: types.Int64Type,
				},
			},
		},
		ThresholdFieldAdaptiveBaseline: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					ThresholdFieldAdaptiveBaselineDeviation:    types.Float32Type,
					ThresholdFieldAdaptiveBaselineAdaptability: types.Float32Type,
					ThresholdFieldAdaptiveBaselineSeasonality:  types.StringType,
				},
			},
		},
	}
}

// GetStaticAndAdaptiveThresholdAttrListTypes returns a ListType for both static and adaptive threshold schemas
func GetStaticAndAdaptiveThresholdAttrListTypes() types.ListType {
	return types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: GetStaticAndAdaptiveThresholdAttrTypes(),
		},
	}
}

// MapThresholdRuleFromState maps a threshold rule from Terraform state to API model
func MapThresholdRuleFromState(ctx context.Context, thresholdList types.List) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	if thresholdList.IsNull() || thresholdList.IsUnknown() {
		return nil, diags
	}

	var thresholdElements []types.Object
	diags.Append(thresholdList.ElementsAs(ctx, &thresholdElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(thresholdElements) == 0 {
		return nil, diags
	}

	// Get the threshold object
	thresholdObj := thresholdElements[0]

	log.Printf("Threshold: %+v\n", thresholdObj)

	// Get the attributes as a map to check which fields exist
	attrs := thresholdObj.Attributes()

	// Check for static threshold
	if staticVal, ok := attrs[ThresholdFieldStatic]; ok && !staticVal.IsNull() && !staticVal.IsUnknown() {
		staticList, ok := staticVal.(types.List)
		if !ok {
			diags.AddError(
				"Invalid static threshold type",
				"Expected list type for static threshold",
			)
			return nil, diags
		}
		return mapStaticThresholdFromState(ctx, staticList)
	}

	// Check for historic baseline threshold
	if historicVal, ok := attrs[ThresholdFieldHistoricBaseline]; ok && !historicVal.IsNull() && !historicVal.IsUnknown() {
		historicList, ok := historicVal.(types.List)
		if !ok {
			diags.AddError(
				"Invalid historic baseline threshold type",
				"Expected list type for historic baseline threshold",
			)
			return nil, diags
		}
		return mapHistoricBaselineFromState(ctx, historicList)
	}

	// Check for adaptive baseline threshold
	if adaptiveVal, ok := attrs[ThresholdFieldAdaptiveBaseline]; ok && !adaptiveVal.IsNull() && !adaptiveVal.IsUnknown() {
		adaptiveList, ok := adaptiveVal.(types.List)
		if !ok {
			diags.AddError(
				"Invalid adaptive baseline threshold type",
				"Expected list type for adaptive baseline threshold",
			)
			return nil, diags
		}
		return mapAdaptiveBaselineFromState(ctx, adaptiveList)
	}

	// If we get here, no valid threshold type was found
	diags.AddError(
		"Invalid threshold configuration",
		"The threshold configuration must include one of: static, historic_baseline, or adaptive_baseline",
	)
	return nil, diags
}

// mapStaticThresholdFromState maps a static threshold from Terraform state to API model
func mapStaticThresholdFromState(ctx context.Context, staticList types.List) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	var staticElements []types.Object
	diags.Append(staticList.ElementsAs(ctx, &staticElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(staticElements) == 0 {
		diags.AddError(
			"Empty static threshold",
			"The static threshold block is empty",
		)
		return nil, diags
	}

	// Extract the value
	var staticObj struct {
		Value types.Int64 `tfsdk:"value"`
	}

	diags.Append(staticElements[0].As(ctx, &staticObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}
	if staticObj.Value.IsNull() || staticObj.Value.IsUnknown() {
		diags.AddError(
			"Static Threshold value is required",
			"Static Threshold value is required",
		)
		return nil, diags
	}

	// Convert to float64 pointer
	valueFloat := float64(staticObj.Value.ValueInt64())

	return &restapi.ThresholdRule{
		Type:  "staticThreshold",
		Value: &valueFloat,
	}, diags
}

// mapHistoricBaselineFromState maps a historic baseline threshold from Terraform state to API model
func mapHistoricBaselineFromState(ctx context.Context, historicList types.List) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	var historicElements []types.Object
	diags.Append(historicList.ElementsAs(ctx, &historicElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(historicElements) == 0 {
		diags.AddError(
			"Empty historic baseline threshold",
			"The historic baseline threshold block is empty",
		)
		return nil, diags
	}

	// Extract the historic baseline configuration
	var historicObj struct {
		Baseline    types.Set     `tfsdk:"baseline"`
		Deviation   types.Float32 `tfsdk:"deviation_factor"`
		Seasonality types.String  `tfsdk:"seasonality"`
	}

	diags.Append(historicElements[0].As(ctx, &historicObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create the threshold rule
	thresholdRule := &restapi.ThresholdRule{
		Type: "historicBaseline",
	}

	// Set seasonality
	if !historicObj.Seasonality.IsNull() && !historicObj.Seasonality.IsUnknown() {
		seasonality := restapi.ThresholdSeasonality(historicObj.Seasonality.ValueString())
		thresholdRule.Seasonality = &seasonality
	}

	// Set deviation factor
	if !historicObj.Deviation.IsNull() && !historicObj.Deviation.IsUnknown() {
		deviationFactor := float32(historicObj.Deviation.ValueFloat32())
		thresholdRule.DeviationFactor = &deviationFactor
	}

	// Set baseline
	if !historicObj.Baseline.IsNull() && !historicObj.Baseline.IsUnknown() {
		// For now, we'll leave it as nil
	}

	return thresholdRule, diags
}

// mapAdaptiveBaselineFromState maps an adaptive baseline threshold from Terraform state to API model
func mapAdaptiveBaselineFromState(ctx context.Context, adaptiveList types.List) (*restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	var adaptiveElements []types.Object
	diags.Append(adaptiveList.ElementsAs(ctx, &adaptiveElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(adaptiveElements) == 0 {
		diags.AddError(
			"Empty adaptive baseline threshold",
			"The adaptive baseline threshold block is empty",
		)
		return nil, diags
	}

	// Extract the adaptive baseline configuration
	var adaptiveObj struct {
		Deviation    types.Float32 `tfsdk:"deviation_factor"`
		Adaptability types.Float32 `tfsdk:"adaptability"`
		Seasonality  types.String  `tfsdk:"seasonality"`
	}

	diags.Append(adaptiveElements[0].As(ctx, &adaptiveObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create the threshold rule
	thresholdRule := &restapi.ThresholdRule{
		Type: "adaptiveBaseline",
	}

	// Set seasonality
	if !adaptiveObj.Seasonality.IsNull() && !adaptiveObj.Seasonality.IsUnknown() {
		seasonality := restapi.ThresholdSeasonality(adaptiveObj.Seasonality.ValueString())
		thresholdRule.Seasonality = &seasonality
	}

	// Set deviation factor
	if !adaptiveObj.Deviation.IsNull() && !adaptiveObj.Deviation.IsUnknown() {
		deviationFactor := float32(adaptiveObj.Deviation.ValueFloat32())
		thresholdRule.DeviationFactor = &deviationFactor
	}

	// Set adaptability (stored in Value field)
	if !adaptiveObj.Adaptability.IsNull() && !adaptiveObj.Adaptability.IsUnknown() {
		adaptability := adaptiveObj.Adaptability.ValueFloat32()
		thresholdRule.Adaptability = &adaptability
	}

	return thresholdRule, diags
}

// Made with Bob

// MapThresholdsFromState maps thresholds from Terraform state to a map of AlertSeverity to ThresholdRule
func MapThresholdsFromState(ctx context.Context, thresholdList types.List) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	if !thresholdList.IsNull() && !thresholdList.IsUnknown() {
		var thresholdElements []types.Object
		diags.Append(thresholdList.ElementsAs(ctx, &thresholdElements, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(thresholdElements) > 0 {
			thresholdObj := thresholdElements[0]

			// Use a properly structured type instead of a generic map
			var thresholdStruct struct {
				Warning  types.List `tfsdk:"warning"`
				Critical types.List `tfsdk:"critical"`
			}

			diags.Append(thresholdObj.As(ctx, &thresholdStruct, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}

			// Process warning threshold
			if !thresholdStruct.Warning.IsNull() && !thresholdStruct.Warning.IsUnknown() {
				warningThreshold, warningDiags := MapThresholdRuleFromState(ctx, thresholdStruct.Warning)
				diags.Append(warningDiags...)
				if diags.HasError() {
					return nil, diags
				}

				if warningThreshold != nil {
					thresholdMap[restapi.WarningSeverity] = *warningThreshold
				}
			}

			// Process critical threshold
			if !thresholdStruct.Critical.IsNull() && !thresholdStruct.Critical.IsUnknown() {
				criticalThreshold, criticalDiags := MapThresholdRuleFromState(ctx, thresholdStruct.Critical)
				diags.Append(criticalDiags...)
				if diags.HasError() {
					return nil, diags
				}

				if criticalThreshold != nil {
					thresholdMap[restapi.CriticalSeverity] = *criticalThreshold
				}
			}
		}
	}

	return thresholdMap, diags
}

func MapThresholdsPluginFromState(ctx context.Context, thresholdStruct ThresholdPluginModel) (map[restapi.AlertSeverity]restapi.ThresholdRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

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
			Type:  "staticThreshold",
			Value: &valueFloat,
		}, diags
	}

	// Check for adaptive baseline threshold
	if thresholdObj.AdaptiveBaseline != nil {
		adaptiveVal := thresholdObj.AdaptiveBaseline
		seasonality := restapi.ThresholdSeasonality(adaptiveVal.Seasonality.ValueString())
		deviationFactor := float32(adaptiveVal.DeviationFactor.ValueFloat32())
		adaptability := adaptiveVal.Adaptability.ValueFloat32()

		return &restapi.ThresholdRule{
			Type:            "adaptiveBaseline",
			Seasonality:     &seasonality,
			DeviationFactor: &deviationFactor,
			Adaptability:    &adaptability,
		}, diags
	}
	return nil, diags
}

// MapThresholdToState maps a threshold rule to a Terraform state representation - used for nested attribute instead of block object
func MapThresholdPluginToState(ctx context.Context, threshold *restapi.ThresholdRule) *ThresholdTypeModel {
	thresholdTypeModel := ThresholdTypeModel{}
	switch threshold.Type {
	case "adaptiveBaseline":
		adaptiveBaselineModel := AdaptiveBaselineModel{
			Operator:        types.StringValue(threshold.Operator),
			DeviationFactor: types.Float32Value(*threshold.DeviationFactor),
			Adaptability:    types.Float32Value(*threshold.Adaptability),
			Seasonality:     types.StringValue(string(*threshold.Seasonality)),
		}
		thresholdTypeModel.AdaptiveBaseline = &adaptiveBaselineModel
	default:
		// Default to static threshold for all other types
		static := StaticTypeModel{
			Operator: types.StringValue(threshold.Operator),
			Value:    types.Int64Value(int64(*threshold.Value)),
		}
		thresholdTypeModel.Static = &static
	}

	return &thresholdTypeModel
}
