package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetThresholdSchema returns the schema for threshold configuration with warning and critical severity levels
func GetThresholdSchema(valueFieldName string, valueType attr.Type) schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Threshold configuration for different severity levels",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ResourceFieldThresholdRuleWarningSeverity: schema.ListNestedBlock{
					Description: "Warning threshold configuration",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"static": schema.ListNestedBlock{
								Description: "Static threshold configuration",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										valueFieldName: getSchemaAttribute(valueType),
									},
								},
							},
						},
					},
				},
				ResourceFieldThresholdRuleCriticalSeverity: schema.ListNestedBlock{
					Description: "Critical threshold configuration",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"static": schema.ListNestedBlock{
								Description: "Static threshold configuration",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										valueFieldName: getSchemaAttribute(valueType),
									},
								},
							},
						},
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

// Helper function to create the appropriate schema attribute based on the type
func getSchemaAttribute(valueType attr.Type) schema.Attribute {
	if _, ok := valueType.(basetypes.Int64Type); ok {
		return schema.Int64Attribute{
			Required:    true,
			Description: "The value of the threshold",
		}
	} else if _, ok := valueType.(basetypes.Float64Type); ok {
		return schema.Float64Attribute{
			Required:    true,
			Description: "The value of the threshold",
		}
	} else {
		// Default to string if type is not recognized
		return schema.StringAttribute{
			Required:    true,
			Description: "The value of the threshold",
		}
	}
}

// MapThresholdRuleToState maps a threshold rule from the API model to the Terraform state
func MapThresholdRuleToState(ctx context.Context, threshold *restapi.ThresholdRule, valueFieldName string, valueType attr.Type) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get the threshold value
	var thresholdValue attr.Value
	if threshold.Value != nil {
		if _, ok := valueType.(basetypes.Int64Type); ok {
			thresholdValue = types.Int64Value(int64(*threshold.Value))
		} else if _, ok := valueType.(basetypes.Float64Type); ok {
			thresholdValue = types.Float64Value(*threshold.Value)
		} else {
			diags.AddError(
				"Unsupported value type",
				"The threshold value type is not supported",
			)
			return types.ListNull(types.ObjectType{}), diags
		}
	} else {
		if _, ok := valueType.(basetypes.Int64Type); ok {
			thresholdValue = types.Int64Value(0)
		} else if _, ok := valueType.(basetypes.Float64Type); ok {
			thresholdValue = types.Float64Value(0.0)
		} else {
			diags.AddError(
				"Unsupported value type",
				"The threshold value type is not supported",
			)
			return types.ListNull(types.ObjectType{}), diags
		}
	}

	// Create static threshold object
	staticObj := map[string]attr.Value{
		valueFieldName: thresholdValue,
	}

	staticObjVal, staticObjDiags := types.ObjectValue(
		map[string]attr.Type{
			valueFieldName: valueType,
		},
		staticObj,
	)
	diags.Append(staticObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create static list
	staticList, staticListDiags := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				valueFieldName: valueType,
			},
		},
		[]attr.Value{staticObjVal},
	)
	diags.Append(staticListDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create threshold object with static block
	thresholdObj := map[string]attr.Value{
		"static": staticList,
	}

	thresholdObjVal, thresholdObjDiags := types.ObjectValue(
		map[string]attr.Type{
			"static": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						valueFieldName: valueType,
					},
				},
			},
		},
		thresholdObj,
	)
	diags.Append(thresholdObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"static": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							valueFieldName: valueType,
						},
					},
				},
			},
		},
		[]attr.Value{thresholdObjVal},
	)
}

// MapThresholdRuleFromState maps a threshold rule from the Terraform state to the API model
func MapThresholdRuleFromState(ctx context.Context, thresholdList types.List, valueFieldName string) (*restapi.ThresholdRule, diag.Diagnostics) {
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

	// Extract the static block using As method
	var thresholdStruct struct {
		Static types.List `tfsdk:"static"`
	}
	diags.Append(thresholdObj.As(ctx, &thresholdStruct, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	if thresholdStruct.Static.IsNull() || thresholdStruct.Static.IsUnknown() {
		diags.AddError(
			"Missing static threshold",
			"The threshold configuration is missing the required 'static' block",
		)
		return nil, diags
	}

	var staticElements []types.Object
	diags.Append(thresholdStruct.Static.ElementsAs(ctx, &staticElements, false)...)
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

	// Extract the value using As method
	staticObj := staticElements[0]

	// Try to get the value as float64 or int64
	var valueFloat float64

	// Try different value types
	var valueStruct struct {
		Value types.Float64 `tfsdk:"value"`
	}
	floatDiags := staticObj.As(ctx, &valueStruct, basetypes.ObjectAsOptions{})
	if !floatDiags.HasError() && !valueStruct.Value.IsNull() && !valueStruct.Value.IsUnknown() {
		valueFloat = valueStruct.Value.ValueFloat64()
	} else {
		// Try as int64
		var intValueStruct struct {
			Value types.Int64 `tfsdk:"value"`
		}
		intDiags := staticObj.As(ctx, &intValueStruct, basetypes.ObjectAsOptions{})
		if !intDiags.HasError() && !intValueStruct.Value.IsNull() && !intValueStruct.Value.IsUnknown() {
			valueFloat = float64(intValueStruct.Value.ValueInt64())
		} else {
			diags.AddError(
				"Invalid threshold value",
				"The threshold value must be a number",
			)
			return nil, diags
		}
	}

	return &restapi.ThresholdRule{
		Type:  "staticThreshold", // Always static for now
		Value: &valueFloat,
	}, diags
}

// GetThresholdAttrTypes returns the attribute types for threshold configuration
func GetThresholdAttrTypes(valueFieldName string, valueType attr.Type) map[string]attr.Type {
	return map[string]attr.Type{
		ResourceFieldThresholdRuleWarningSeverity: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"static": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								valueFieldName: valueType,
							},
						},
					},
				},
			},
		},
		ResourceFieldThresholdRuleCriticalSeverity: types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"static": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								valueFieldName: valueType,
							},
						},
					},
				},
			},
		},
	}
}

// GetThresholdListType returns the ListType for threshold configuration with warning and critical severity levels
func GetThresholdListType() types.ListType {
	return types.ListType{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				ResourceFieldThresholdRuleWarningSeverity: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"static": types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"value": types.Int64Type,
									},
								},
							},
						},
					},
				},
				ResourceFieldThresholdRuleCriticalSeverity: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"static": types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"value": types.Int64Type,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Made with Bob
