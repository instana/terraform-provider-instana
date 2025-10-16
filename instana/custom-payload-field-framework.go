package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
		"dynamic_value": GetDynamicValueType(),
	}
}

type customPayloadFieldTFModel struct {
	Key          types.String `tfsdk:"key"`
	Value        types.String `tfsdk:"value"`
	DynamicValue types.Object `tfsdk:"dynamic_value"`
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

// GetCustomPayloadFieldsSchema returns the schema for custom payload fields as a list nested block
// This can be used by any resource that needs to define custom payload fields
func GetCustomPayloadFieldsSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Custom payload fields for the configuration.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				CustomPayloadFieldsFieldKey: schema.StringAttribute{
					Required:    true,
					Description: "The key of the custom payload field",
				},
				CustomPayloadFieldsFieldStaticStringValue: schema.StringAttribute{
					Optional:    true,
					Description: "The value of a static string custom payload field",
				},
			},
			Blocks: map[string]schema.Block{
				CustomPayloadFieldsFieldDynamicValue: schema.SingleNestedBlock{
					Description: "The value of a dynamic custom payload field",
					Attributes: map[string]schema.Attribute{
						CustomPayloadFieldsFieldDynamicKey: schema.StringAttribute{
							Optional:    true,
							Description: "The key of the dynamic custom payload field",
						},
						CustomPayloadFieldsFieldDynamicTagName: schema.StringAttribute{
							Optional:    true,
							Description: "The name of the tag of the dynamic custom payload field",
						},
					},
				},
			},
		},
	}
}

// GetStaticOnlyCustomPayloadFieldsSchema returns the schema for custom payload fields as a list nested block
// This version only supports static string values and does not include the dynamic_value block
func GetStaticOnlyCustomPayloadFieldsSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Custom payload fields for the configuration.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				CustomPayloadFieldsFieldKey: schema.StringAttribute{
					Required:    true,
					Description: "The key of the custom payload field",
				},
				CustomPayloadFieldsFieldStaticStringValue: schema.StringAttribute{
					Required:    true,
					Description: "The value of a static string custom payload field",
				},
			},
		},
	}
}

// MapCustomPayloadFieldsToTerraform is a helper function to map custom payload fields from API objects to Terraform
// This function delegates to the existing CustomPayloadFieldsToTerraform function in tfutils
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
			tfField["dynamic_value"] = types.ObjectNull(GetDynamicValueType().AttrTypes)
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

			// Set the dynamic value object directly
			tfField["dynamic_value"] = dynamicObj
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

// MapCustomPayloadFieldsToAPIObject is a helper function to map custom payload fields from Terraform to API objects
func MapCustomPayloadFieldsToAPIObject(ctx context.Context, customPayloadFieldsList types.List) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics

	// handle null/unknown
	if customPayloadFieldsList.IsNull() || customPayloadFieldsList.IsUnknown() {
		return []restapi.CustomPayloadField[any]{}, diags
	}

	// Decode terraform list elements into our intermediate struct slice
	var elems []customPayloadFieldTFModel
	diags = append(diags, customPayloadFieldsList.ElementsAs(ctx, &elems, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]restapi.CustomPayloadField[any], 0, len(elems))

	for idx, e := range elems {
		// If dynamic_value present -> dynamic
		if !e.DynamicValue.IsNull() && !e.DynamicValue.IsUnknown() {
			// Get attributes from the dynamic_value object
			dynAttrs := e.DynamicValue.Attributes()

			// Extract key and tag_name values
			var keyPtr *string
			if keyAttr, ok := dynAttrs["key"]; ok && !keyAttr.(types.String).IsNull() && !keyAttr.(types.String).IsUnknown() {
				tmp := keyAttr.(types.String).ValueString()
				if tmp != "" {
					keyPtr = &tmp
				}
			}

			tagNameAttr, ok := dynAttrs["tag_name"]
			if !ok || tagNameAttr.(types.String).IsNull() || tagNameAttr.(types.String).IsUnknown() {
				// Only validate tag_name if dynamic_value is present and being used
				if !e.DynamicValue.IsNull() && !e.DynamicValue.IsUnknown() {
					diags.AddError(
						"custom_payload_field.dynamic_value missing tag_name",
						fmt.Sprintf("element index %d: dynamic_value object missing required 'tag_name' attribute", idx),
					)
					return nil, diags
				}
			}
			tagName := ""
			if ok && !tagNameAttr.(types.String).IsNull() && !tagNameAttr.(types.String).IsUnknown() {
				tagName = tagNameAttr.(types.String).ValueString()
			}

			// If tag_name is empty, we can't create a valid dynamic value
			if tagName == "" {
				diags.AddError(
					"custom_payload_field.dynamic_value missing tag_name",
					fmt.Sprintf("element index %d: dynamic_value object requires a 'tag_name' attribute", idx),
				)
				return nil, diags
			}

			dynValue := restapi.DynamicCustomPayloadFieldValue{
				Key:     keyPtr,
				TagName: tagName,
			}

			// top-level key may exist in e.Key
			topKey := ""
			if !e.Key.IsNull() && !e.Key.IsUnknown() {
				topKey = e.Key.ValueString()
			}

			out = append(out, restapi.CustomPayloadField[any]{
				Type:  restapi.DynamicCustomPayloadType,
				Key:   topKey,
				Value: any(dynValue),
			})
			continue
		}

		// else if static value present
		if !e.Value.IsNull() && !e.Value.IsUnknown() {
			out = append(out, restapi.CustomPayloadField[any]{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   e.Key.ValueString(),
				Value: any(e.Value.ValueString()),
			})
			continue
		}

		// neither dynamic nor static -> error
		diags.AddError(
			"custom_payload_field missing value",
			fmt.Sprintf("element index %d: neither 'value' nor 'dynamic_value' present", idx),
		)
		return nil, diags
	}

	return out, diags
}

// MapStaticOnlyCustomPayloadFieldsToAPIObject is a helper function to map static-only custom payload fields from Terraform to API objects
func MapStaticOnlyCustomPayloadFieldsToAPIObject(ctx context.Context, customPayloadFieldsList types.List) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics

	// handle null/unknown
	if customPayloadFieldsList.IsNull() || customPayloadFieldsList.IsUnknown() {
		return []restapi.CustomPayloadField[any]{}, diags
	}

	// Define a simple model for static-only fields
	type staticOnlyFieldModel struct {
		Key   types.String `tfsdk:"key"`
		Value types.String `tfsdk:"value"`
	}

	// Decode terraform list elements into our intermediate struct slice
	var elems []staticOnlyFieldModel
	diags = append(diags, customPayloadFieldsList.ElementsAs(ctx, &elems, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]restapi.CustomPayloadField[any], 0, len(elems))

	for idx, e := range elems {
		// Check if key and value are present
		if e.Key.IsNull() || e.Key.IsUnknown() {
			diags.AddError(
				"custom_payload_field missing key",
				fmt.Sprintf("element index %d: 'key' is required", idx),
			)
			return nil, diags
		}

		if e.Value.IsNull() || e.Value.IsUnknown() {
			diags.AddError(
				"custom_payload_field missing value",
				fmt.Sprintf("element index %d: 'value' is required", idx),
			)
			return nil, diags
		}

		// Add the static field to the output
		out = append(out, restapi.CustomPayloadField[any]{
			Type:  restapi.StaticStringCustomPayloadType,
			Key:   e.Key.ValueString(),
			Value: any(e.Value.ValueString()),
		})
	}

	return out, diags
}

// StaticOnlyCustomPayloadFieldsToTerraform is a helper function to map custom payload fields from API objects to Terraform
// This version only includes static string values and omits dynamic values
func StaticOnlyCustomPayloadFieldsToTerraform(ctx context.Context, fields []restapi.CustomPayloadField[any]) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// If no fields, return null list
	if len(fields) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		}), diags
	}

	// Convert API fields to a slice of maps that can be used with ObjectValueFrom
	tfFields := make([]attr.Value, 0, len(fields))

	for _, field := range fields {
		// Only process static string fields
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
					continue
				}
			}

			tfField := map[string]attr.Value{
				"key":   types.StringValue(field.Key),
				"value": types.StringValue(staticValue),
			}

			objVal, d := types.ObjectValue(
				map[string]attr.Type{
					"key":   types.StringType,
					"value": types.StringType,
				},
				tfField,
			)
			diags = append(diags, d...)
			if d.HasError() {
				continue
			}

			tfFields = append(tfFields, objVal)
		}
		// Skip dynamic fields
	}

	// Create the list of custom payload fields
	return types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":   types.StringType,
			"value": types.StringType,
		},
	}, tfFields)
}

// Made with Bob
