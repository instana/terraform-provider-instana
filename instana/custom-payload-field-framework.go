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
		"dynamic_value": types.ListType{ElemType: GetDynamicValueType()},
	}
}

type customPayloadFieldTFModel struct {
	Key          types.String `tfsdk:"key"`
	Value        types.String `tfsdk:"value"`
	DynamicValue types.List   `tfsdk:"dynamic_value"`
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
				CustomPayloadFieldsFieldDynamicValue: schema.ListNestedAttribute{
					Optional:    true,
					Description: "The value of a dynamic custom payload field",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							CustomPayloadFieldsFieldDynamicKey: schema.StringAttribute{
								Optional:    true,
								Description: "The key of the dynamic custom payload field",
							},
							CustomPayloadFieldsFieldDynamicTagName: schema.StringAttribute{
								Required:    true,
								Description: "The name of the tag of the dynamic custom payload field",
							},
						},
					},
				},
			},
		},
	}
}

// GetCustomPayloadFieldsSetAttribute returns the schema for custom payload fields as a set nested attribute
// This can be used by any resource that needs to define custom payload fields as a set
func GetCustomPayloadFieldsSetAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Description: "Custom payload fields for the configuration.",
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				CustomPayloadFieldsFieldKey: schema.StringAttribute{
					Required:    true,
					Description: "The key of the custom payload field",
				},
				CustomPayloadFieldsFieldStaticStringValue: schema.StringAttribute{
					Optional:    true,
					Description: "The value of a static string custom payload field",
				},
				CustomPayloadFieldsFieldDynamicValue: schema.ListNestedAttribute{
					Optional:    true,
					Description: "The value of a dynamic custom payload field",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							CustomPayloadFieldsFieldDynamicKey: schema.StringAttribute{
								Optional:    true,
								Description: "The key of the dynamic custom payload field",
							},
							CustomPayloadFieldsFieldDynamicTagName: schema.StringAttribute{
								Required:    true,
								Description: "The name of the tag of the dynamic custom payload field",
							},
						},
					},
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
			// Expect dynamic_value to be a list of maps[string]string
			var dynList []map[string]string
			diags = append(diags, e.DynamicValue.ElementsAs(ctx, &dynList, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(dynList) == 0 {
				diags.AddError(
					"custom_payload_field.dynamic_value empty",
					fmt.Sprintf("element index %d: dynamic_value list is empty", idx),
				)
				return nil, diags
			}

			// take first map (adapt if your schema allows multiples)
			dynMap := dynList[0]

			// build DynamicCustomPayloadFieldValue
			var keyPtr *string
			if v, ok := dynMap["key"]; ok && v != "" {
				tmp := v
				keyPtr = &tmp
			}
			tagName, ok := dynMap["tagName"]
			if !ok {
				diags.AddError(
					"custom_payload_field.dynamic_value missing tagName",
					fmt.Sprintf("element index %d: dynamic_value map missing required 'tagName' key", idx),
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

// Made with Bob
