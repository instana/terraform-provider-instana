package tfutils

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomPayloadFieldAttributeTypes returns a map of attribute types for custom payload fields
// This can be used by any resource that needs to define custom payload fields
// The structure follows the pattern defined in custom-payload-field-schema.go
//
// Returns:
//   - map[string]attr.Type: A map of attribute names to their types for custom payload fields
//
// This function is useful for schema definitions in Terraform resources that need to include
// custom payload fields as part of their configuration.
func CustomPayloadFieldAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":           types.StringType,
		"value":         types.StringType,
		"dynamic_value": types.ListType{ElemType: GetDynamicValueType()},
	}
}

// GetDynamicValueType returns the object type for dynamic values in custom payload fields
// This type definition is used for the nested dynamic_value field in custom payload fields
// It defines the structure with "key" (optional) and "tag_name" (required) attributes
func GetDynamicValueType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":      types.StringType,
			"tag_name": types.StringType,
		},
	}
}

// GetCustomPayloadFieldType returns the object type for custom payload fields
// This type definition is used for the top-level custom payload field structure
// It defines the schema with "key" (required), "value" (for static string values),
// and "dynamic_value" (for dynamic values) attributes
func GetCustomPayloadFieldType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: CustomPayloadFieldAttributeTypes(),
	}
}

// CustomPayloadFieldsToTerraform converts API custom payload fields to Terraform types
// This function handles both static string and dynamic custom payload field types
//
// Parameters:
//   - ctx: The context for the operation
//   - fields: A slice of CustomPayloadField objects from the API
//
// Returns:
//   - types.List: A Terraform list value containing the converted custom payload fields
//   - diag.Diagnostics: Any diagnostics that occurred during the conversion
//
// The function handles both static string and dynamic custom payload field types:
//   - For static string fields, it sets the "value" attribute and sets "dynamic_value" to null
//   - For dynamic fields, it sets the "dynamic_value" attribute and sets "value" to null
//   - If no fields are provided, it returns a null list
//
// This implementation uses types.ListValueFrom and types.ObjectValueFrom for cleaner code
func CustomPayloadFieldsToTerraform(ctx context.Context, fields []restapi.CustomPayloadField[any]) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// If no fields, return null list
	if len(fields) == 0 {
		return types.ListNull(GetCustomPayloadFieldType()), diags
	}

	// Convert API fields to a slice of maps that can be used with ObjectValueFrom
	tfFields := make([]attr.Value, 0, len(fields))

	for _, field := range fields {
		tfField := map[string]attr.Value{}
		tfField["key"] = types.StringValue(field.Key)

		// Handle different field types
		if field.Type == restapi.StaticStringCustomPayloadType {
			// Static string value
			staticValue, ok := field.Value.(string)
			if !ok {
				// Try to convert from the custom type
				if customValue, ok := field.Value.(restapi.StaticStringCustomPayloadFieldValue); ok {
					staticValue = string(customValue)
				} else {
					diags.AddError(
						"Error converting custom payload field",
						fmt.Sprintf("Failed to convert static string value for key %s", field.Key),
					)
					return types.ListNull(GetCustomPayloadFieldType()), diags
				}
			}
			tfField["value"] = types.StringValue(staticValue)
			// Use null for dynamic_value
			tfField["dynamic_value"] = types.ListNull(GetDynamicValueType())
		} else if field.Type == restapi.DynamicCustomPayloadType {
			// Dynamic value
			dynamicValue, ok := field.Value.(restapi.DynamicCustomPayloadFieldValue)
			if !ok {
				diags.AddError(
					"Error converting custom payload field",
					fmt.Sprintf("Failed to convert dynamic value for key %s", field.Key),
				)
				return types.ListNull(GetCustomPayloadFieldType()), diags
			}

			// Create a map for the dynamic value
			dynamicAttrs := make(map[string]attr.Value)

			dynamicAttrs["tag_name"] = types.StringValue(dynamicValue.TagName)

			if dynamicValue.Key != nil {
				dynamicAttrs["key"] = types.StringValue(*dynamicValue.Key)
			} else {
				dynamicAttrs["key"] = types.StringNull()
			}

			// Create the dynamic value object
			dynamicObj, d := types.ObjectValue(GetDynamicValueType().AttrTypes, dynamicAttrs)
			if d.HasError() {
				diags = append(diags, d...)
				return types.List{}, diags
			}

			// Create a list with one dynamic value object
			dynListVal, d := types.ListValue(GetDynamicValueType(), []attr.Value{dynamicObj})
			if d.HasError() {
				diags = append(diags, d...)
				return types.List{}, diags
			}

			tfField["dynamic_value"] = dynListVal
			tfField["value"] = types.StringNull()
		}
		objVal, d := types.ObjectValue(CustomPayloadFieldAttributeTypes(), tfField)
		diags = append(diags, d...)
		if d.HasError() {
			return types.List{}, diags
		}

		tfFields = append(tfFields, objVal)
	}

	// Create the list of custom payload fields
	return types.ListValue(types.ObjectType{
		AttrTypes: CustomPayloadFieldAttributeTypes(),
	}, tfFields)
}

// Made with Bob
