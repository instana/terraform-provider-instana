package tfutils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UpdateState(d *sdkschema.ResourceData, data map[string]interface{}) error {
	for k, v := range data {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateStatePlugin(ctx context.Context, state *tfsdk.State, data map[string]interface{}) error {
	for k, v := range data {
		if v == nil {
			// Set null value for nil
			diags := state.SetAttribute(ctx, path.Root(k), types.StringNull())
			if diags.HasError() {
				return fmt.Errorf("error setting null attribute %s: %s", k, diagsToString(diags))
			}
			continue
		}

		// Handle different types of values
		switch value := v.(type) {
		case string:
			diags := state.SetAttribute(ctx, path.Root(k), types.StringValue(value))
			if diags.HasError() {
				return fmt.Errorf("error setting string attribute %s: %s", k, diagsToString(diags))
			}
		case bool:
			diags := state.SetAttribute(ctx, path.Root(k), types.BoolValue(value))
			if diags.HasError() {
				return fmt.Errorf("error setting bool attribute %s: %s", k, diagsToString(diags))
			}
		case int:
			diags := state.SetAttribute(ctx, path.Root(k), types.Int64Value(int64(value)))
			if diags.HasError() {
				return fmt.Errorf("error setting int attribute %s: %s", k, diagsToString(diags))
			}
		case int64:
			diags := state.SetAttribute(ctx, path.Root(k), types.Int64Value(value))
			if diags.HasError() {
				return fmt.Errorf("error setting int64 attribute %s: %s", k, diagsToString(diags))
			}
		case float64:
			diags := state.SetAttribute(ctx, path.Root(k), types.Float64Value(value))
			if diags.HasError() {
				return fmt.Errorf("error setting float64 attribute %s: %s", k, diagsToString(diags))
			}
		case []string:
			// Check if the attribute is a set or list
			attrType := getAttributeType(state, k)
			if attrType == "set" {
				setValue, diags := types.SetValueFrom(ctx, types.StringType, value)
				if diags.HasError() {
					return fmt.Errorf("error creating set value for attribute %s: %s", k, diagsToString(diags))
				}
				diags = state.SetAttribute(ctx, path.Root(k), setValue)
				if diags.HasError() {
					return fmt.Errorf("error setting set attribute %s: %s", k, diagsToString(diags))
				}
			} else {
				// Default to list
				listValue, diags := types.ListValueFrom(ctx, types.StringType, value)
				if diags.HasError() {
					return fmt.Errorf("error creating list value for attribute %s: %s", k, diagsToString(diags))
				}
				diags = state.SetAttribute(ctx, path.Root(k), listValue)
				if diags.HasError() {
					return fmt.Errorf("error setting list attribute %s: %s", k, diagsToString(diags))
				}
			}
		case []interface{}:
			// Try to determine the element type from the first element
			if len(value) > 0 {
				firstElem := value[0]
				switch firstElem.(type) {
				case string:
					// Convert all elements to strings
					strSlice := make([]string, len(value))
					for i, v := range value {
						if str, ok := v.(string); ok {
							strSlice[i] = str
						} else {
							return fmt.Errorf("element %d in attribute %s is not a string", i, k)
						}
					}

					// Check if the attribute is a set or list
					attrType := getAttributeType(state, k)
					if attrType == "set" {
						setValue, diags := types.SetValueFrom(ctx, types.StringType, strSlice)
						if diags.HasError() {
							return fmt.Errorf("error creating set value for attribute %s: %s", k, diagsToString(diags))
						}
						diags = state.SetAttribute(ctx, path.Root(k), setValue)
						if diags.HasError() {
							return fmt.Errorf("error setting set attribute %s: %s", k, diagsToString(diags))
						}
					} else {
						// Default to list
						listValue, diags := types.ListValueFrom(ctx, types.StringType, strSlice)
						if diags.HasError() {
							return fmt.Errorf("error creating string list value for attribute %s: %s", k, diagsToString(diags))
						}
						diags = state.SetAttribute(ctx, path.Root(k), listValue)
						if diags.HasError() {
							return fmt.Errorf("error setting string list attribute %s: %s", k, diagsToString(diags))
						}
					}
				default:
					return fmt.Errorf("unsupported list element type %T for attribute %s", firstElem, k)
				}
			} else {
				// Empty list, default to string type
				// Check if the attribute is a set or list
				attrType := getAttributeType(state, k)
				if attrType == "set" {
					setValue, diags := types.SetValueFrom(ctx, types.StringType, []string{})
					if diags.HasError() {
						return fmt.Errorf("error creating empty set value for attribute %s: %s", k, diagsToString(diags))
					}
					diags = state.SetAttribute(ctx, path.Root(k), setValue)
					if diags.HasError() {
						return fmt.Errorf("error setting empty set attribute %s: %s", k, diagsToString(diags))
					}
				} else {
					// Default to list
					listValue, diags := types.ListValueFrom(ctx, types.StringType, []string{})
					if diags.HasError() {
						return fmt.Errorf("error creating empty list value for attribute %s: %s", k, diagsToString(diags))
					}
					diags = state.SetAttribute(ctx, path.Root(k), listValue)
					if diags.HasError() {
						return fmt.Errorf("error setting empty list attribute %s: %s", k, diagsToString(diags))
					}
				}
			}
		case map[string]interface{}:
			mapValue, diags := types.MapValueFrom(ctx, types.StringType, value)
			if diags.HasError() {
				return fmt.Errorf("error creating map value for attribute %s: %s", k, diagsToString(diags))
			}
			diags = state.SetAttribute(ctx, path.Root(k), mapValue)
			if diags.HasError() {
				return fmt.Errorf("error setting map attribute %s: %s", k, diagsToString(diags))
			}
		case map[string]string:
			mapValue, diags := types.MapValueFrom(ctx, types.StringType, value)
			if diags.HasError() {
				return fmt.Errorf("error creating string map value for attribute %s: %s", k, diagsToString(diags))
			}
			diags = state.SetAttribute(ctx, path.Root(k), mapValue)
			if diags.HasError() {
				return fmt.Errorf("error setting string map attribute %s: %s", k, diagsToString(diags))
			}
		case []map[string]interface{}:
			// For list of objects
			objListValue, diags := types.ListValueFrom(ctx, types.ObjectType{
				AttrTypes: getObjectAttrTypes(value),
			}, value)
			if diags.HasError() {
				return fmt.Errorf("error creating object list value for attribute %s: %s", k, diagsToString(diags))
			}
			diags = state.SetAttribute(ctx, path.Root(k), objListValue)
			if diags.HasError() {
				return fmt.Errorf("error setting object list attribute %s: %s", k, diagsToString(diags))
			}
		default:
			// Try to handle sets
			valueType := reflect.TypeOf(v)
			if valueType.Kind() == reflect.Slice || valueType.Kind() == reflect.Array {
				// Check if the attribute is a set or list
				attrType := getAttributeType(state, k)
				if attrType == "set" {
					setValue, diags := types.SetValueFrom(ctx, types.StringType, v)
					if diags.HasError() {
						return fmt.Errorf("error creating set value for attribute %s: %s", k, diagsToString(diags))
					}
					diags = state.SetAttribute(ctx, path.Root(k), setValue)
					if diags.HasError() {
						return fmt.Errorf("error setting set attribute %s: %s", k, diagsToString(diags))
					}
				} else {
					// Default to list
					listValue, diags := types.ListValueFrom(ctx, types.StringType, v)
					if diags.HasError() {
						return fmt.Errorf("error creating list value for attribute %s: %s", k, diagsToString(diags))
					}
					diags = state.SetAttribute(ctx, path.Root(k), listValue)
					if diags.HasError() {
						return fmt.Errorf("error setting list attribute %s: %s", k, diagsToString(diags))
					}
				}
			} else {
				return fmt.Errorf("unsupported type %T for attribute %s", v, k)
			}
		}
	}

	return nil
}

// Helper function to determine if an attribute is a set or list
func getAttributeType(state *tfsdk.State, attrName string) string {
	if state == nil || state.Raw.IsNull() || !state.Raw.IsKnown() {
		return "list" // Default to list if schema is not available
	}

	// Try to get the attribute type from the raw value
	objType, ok := state.Raw.Type().(tftypes.Object)
	if !ok {
		return "list" // Default to list if raw value is not an object
	}

	attrType, ok := objType.AttributeTypes[attrName]
	if !ok {
		return "list" // Default to list if attribute type not found
	}

	// Check if it's a set
	if _, ok := attrType.(tftypes.Set); ok {
		return "set"
	}

	return "list" // Default to list for all other cases
}

// Helper function to determine object attribute types from a slice of maps
func getObjectAttrTypes(maps []map[string]interface{}) map[string]attr.Type {
	if len(maps) == 0 {
		return map[string]attr.Type{}
	}

	// Use the first map to determine attribute types
	firstMap := maps[0]
	attrTypes := make(map[string]attr.Type)

	for k, v := range firstMap {
		switch v.(type) {
		case string:
			attrTypes[k] = types.StringType
		case bool:
			attrTypes[k] = types.BoolType
		case int, int64:
			attrTypes[k] = types.Int64Type
		case float64:
			attrTypes[k] = types.Float64Type
		case []interface{}, []string:
			attrTypes[k] = types.ListType{
				ElemType: types.StringType,
			}
		case map[string]interface{}:
			attrTypes[k] = types.MapType{
				ElemType: types.StringType,
			}
		default:
			// Default to string for unknown types
			attrTypes[k] = types.StringType
		}
	}

	return attrTypes
}

// Helper function to convert diagnostics to a string
func diagsToString(diags diag.Diagnostics) string {
	if !diags.HasError() {
		return ""
	}

	var errorMessages []string
	for _, d := range diags {
		if d.Severity() == diag.SeverityError {
			errorMessages = append(errorMessages, d.Summary()+": "+d.Detail())
		}
	}

	return strings.Join(errorMessages, "; ")
}

func ArrayToStateSet(elements []string) (types.Set, diag.Diagnostics) {
	// Each element in the Set must have an element type
	elems := make([]attr.Value, len(elements))
	for i, id := range elements {
		elems[i] = types.StringValue(id)
	}

	// Build the types.Set from the elements
	setVal, diags := types.SetValue(types.StringType, elems)
	return setVal, diags
}

// Made with Bob
