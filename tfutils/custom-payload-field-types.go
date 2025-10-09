package tfutils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomPayloadFieldAttributeTypes returns a map of attribute types for custom payload fields
// This can be used by any resource that needs to define custom payload fields
// The structure follows the pattern defined in custom-payload-field-schema.go
func CustomPayloadFieldAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":           types.StringType,
		"value":         types.StringType,
		"dynamic_value": types.ListType{ElemType: types.MapType{ElemType: types.StringType}},
	}
}

// Made with Bob
