package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Constants for threshold field names
const (
	ThresholdFieldWarning  = "warning"
	ThresholdFieldCritical = "critical"
	ThresholdFieldStatic   = "static"
)

// StaticThresholdBlockSchema returns the schema for static threshold configuration
func StaticThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Static threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				LogAlertConfigFieldValue: schema.Int64Attribute{
					Required:    true,
					Description: "The value of the threshold",
				},
			},
		},
	}
}

// WarningThresholdBlockSchema returns the schema for warning threshold configuration
func WarningThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Warning threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ThresholdFieldStatic: StaticThresholdBlockSchema(),
			},
		},
	}
}

// CriticalThresholdBlockSchema returns the schema for critical threshold configuration
func CriticalThresholdBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Critical threshold configuration",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				ThresholdFieldStatic: StaticThresholdBlockSchema(),
			},
		},
	}
}

// MapThresholdToState maps a threshold rule to a Terraform state representation
func MapThresholdToState(ctx context.Context, isThresholdPresent bool, threshold *restapi.ThresholdRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if isThresholdPresent {
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
		staticList, staticListDiags := types.ListValue(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					LogAlertConfigFieldValue: types.Int64Type,
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
			ThresholdFieldStatic: staticList,
		}

		thresholdObjVal, thresholdObjDiags := types.ObjectValue(
			map[string]attr.Type{
				ThresholdFieldStatic: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							LogAlertConfigFieldValue: types.Int64Type,
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
					ThresholdFieldStatic: types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								LogAlertConfigFieldValue: types.Int64Type,
							},
						},
					},
				},
			},
			[]attr.Value{thresholdObjVal},
		)
	} else {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				ThresholdFieldStatic: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							LogAlertConfigFieldValue: types.Int64Type,
						},
					},
				},
			},
		}), diags
	}
}

// Made with Bob
