package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Using restapi.CustomPayloadType instead of defining our own

// Using restapi.DynamicCustomPayloadFieldValue instead of defining our own

// --- Terraform intermediate model matching schema of each list element ---
// (each element has key, value, dynamic_value)
type customPayloadFieldTFModel struct {
	Key          types.String `tfsdk:"key"`
	Value        types.String `tfsdk:"value"`
	DynamicValue types.List   `tfsdk:"dynamic_value"`
}

func BuildCustomPayloadFieldsTyped(ctx context.Context, customPayloadFieldsList types.List) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
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
