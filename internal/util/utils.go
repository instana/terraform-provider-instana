package util

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/google/go-cmp/cmp"

	"github.com/instana/terraform-provider-instana/internal/restapi"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/xid"
)

// DecimalPrecisionMultiplier is used for rounding floats to 2 decimal places
const DecimalPrecisionMultiplier = 100

// RoundFloat64To2Decimals rounds a float64 value to 2 decimal places
// This is more efficient than string formatting and parsing
func RoundFloat64To2Decimals(value float64) float64 {
	return math.Round(value*DecimalPrecisionMultiplier) / DecimalPrecisionMultiplier
}

// RoundFloat32To2Decimals rounds a float32 value to 2 decimal places and returns float32
func RoundFloat32To2Decimals(value float32) float32 {
	return float32(math.Round(float64(value)*DecimalPrecisionMultiplier) / DecimalPrecisionMultiplier)
}

// ConvertInt64ToInt32WithValidation converts int64 to int32 with overflow checking
// Returns the converted value and a boolean indicating if the conversion was successful
func ConvertInt64ToInt32WithValidation(value int64) (int32, bool) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, false
	}
	return int32(value), true
}

// RandomID generates a random ID for a resource
func RandomID() string {
	id := xid.New()
	return id.String()
}

// ReadArrayParameterFromMap reads an array parameter from a resource
func ReadArrayParameterFromMap[T any](data map[string]interface{}, key string) []T {
	if attr, ok := data[key]; ok {
		var array []T
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(T)
			array = append(array, item)
		}
		return array
	}
	return nil
}

// ReadStringSetParameterFromResource reads a string set parameter from a resource and returns it as a slice of strings
func ReadStringSetParameterFromResource(d *schema.ResourceData, key string) []string {
	if attr, ok := d.GetOk(key); ok {
		var array []string
		set := attr.(*schema.Set)
		for _, x := range set.List() {
			item := x.(string)
			array = append(array, item)
		}
		return array
	}
	return nil
}

// ReadArrayParameterFromResource reads a string array parameter from a resource and returns it as a slice of strings
func ReadArrayParameterFromResource[T any](d *schema.ResourceData, key string) []T {
	if attr, ok := d.GetOk(key); ok {
		var array []T
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(T)
			array = append(array, item)
		}
		return array
	}
	return []T{}
}

// ReadSetParameterFromMap reads a set parameter from a map and returns it as a slice
func ReadSetParameterFromMap[T any](data map[string]interface{}, key string) []T {
	if attr, ok := data[key]; ok {
		var array []T
		set := attr.(*schema.Set)
		for _, x := range set.List() {
			item := x.(T)
			array = append(array, item)
		}
		return array
	}
	return nil
}

// ConvertSeverityFromInstanaAPIToTerraformRepresentation converts the integer representation of the Instana API to the string representation of the Terraform provider
func ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity int) (string, error) {
	if severity == restapi.SeverityWarning.GetAPIRepresentation() {
		return restapi.SeverityWarning.GetTerraformRepresentation(), nil
	} else if severity == restapi.SeverityCritical.GetAPIRepresentation() {
		return restapi.SeverityCritical.GetTerraformRepresentation(), nil
	} else {
		return "INVALID", fmt.Errorf("%d is not a valid severity", severity)
	}
}

// ConvertSeverityFromTerraformToInstanaAPIRepresentation converts the string representation of the Terraform to the int representation of the Instana API provider
func ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity string) (int, error) {
	if severity == restapi.SeverityWarning.GetTerraformRepresentation() {
		return restapi.SeverityWarning.GetAPIRepresentation(), nil
	} else if severity == restapi.SeverityCritical.GetTerraformRepresentation() {
		return restapi.SeverityCritical.GetAPIRepresentation(), nil
	} else {
		return -1, fmt.Errorf("%s is not a valid severity", severity)
	}
}

// GetIntPointerFromResourceData gets a int value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetIntPointerFromResourceData(d *schema.ResourceData, key string) *int {
	val, ok := d.GetOk(key)
	if ok {
		intValue := val.(int)
		return &intValue
	}
	return nil
}

// GetInt32PointerFromResourceData gets a int32 value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetInt32PointerFromResourceData(d *schema.ResourceData, key string) *int32 {
	val, ok := d.GetOk(key)
	if ok {
		intValue := int32(val.(int))
		return &intValue
	}
	return nil
}

// GetFloat64PointerFromResourceData gets a float64 value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetFloat64PointerFromResourceData(d *schema.ResourceData, key string) *float64 {
	val, ok := d.GetOk(key)
	if ok {
		floatValue := val.(float64)
		return &floatValue
	}
	return nil
}

// GetFloat32PointerFromResourceData gets a float32 value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetFloat32PointerFromResourceData(d *schema.ResourceData, key string) *float32 {
	val, ok := d.GetOk(key)
	if ok {
		floatValue := float32(val.(float64))
		return &floatValue
	}
	return nil
}

// GetStringPointerFromResourceData gets a string value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetStringPointerFromResourceData(d *schema.ResourceData, key string) *string {
	val, ok := d.GetOk(key)
	if ok {
		stringValue := val.(string)
		return &stringValue
	}
	return nil
}

// GetPointerFromMap gets a pointer of the given type from the map or nil if the value is not defined
func GetPointerFromMap[T any](d map[string]interface{}, key string) *T {
	var defaultValue T
	val, ok := d[key]
	if ok && val != nil && !cmp.Equal(defaultValue, val) {
		value := val.(T)
		return &value
	}
	return nil
}

// MergeSchemaMap merges the provided maps into a single map
func MergeSchemaMap(mapA map[string]*schema.Schema, mapB map[string]*schema.Schema) map[string]*schema.Schema {
	mergedMap := make(map[string]*schema.Schema)

	for k, v := range mapA {
		mergedMap[k] = v
	}
	for k, v := range mapB {
		mergedMap[k] = v
	}

	return mergedMap
}

// ConvertInterfaceSlice converts the given interface slice to the desired target slice
func ConvertInterfaceSlice[T any](input []interface{}) []T {
	result := make([]T, len(input))
	for i, v := range input {
		result[i] = v.(T)
	}
	return result
}

func NormalizeJSONString(jsonString string) string {
	var raw []map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &raw)
	if err != nil {
		return jsonString
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return jsonString
	}
	return string(bytes)
}

// NormalizeJSONObjectString normalizes a JSON object string (not an array)
// This is useful for webpage scripts and other single JSON objects
func NormalizeJSONObjectString(jsonString string) string {
	var raw map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &raw)
	if err != nil {
		return jsonString
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return jsonString
	}
	return string(bytes)
}

// this interface to handle the conversion of differnt numeric types
type numericPtr interface {
	~*int32 | ~*int64 | ~*float32 | ~*float64 | ~*int
}

func SetStringPointerToState(i *string) types.String {
	if i == nil || *i == "" {
		return types.StringNull()
	} else {
		return types.StringValue(*i)
	}
}

func SetInt64PointerToState[T numericPtr](i T) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}
	switch v := any(i).(type) {
	case *int:
		return types.Int64Value(int64(*v))
	case *int32:
		return types.Int64Value(int64(*v))
	case *int64:
		return types.Int64Value(*v)
	case *float32:
		return types.Int64Value(int64(*v))
	case *float64:
		return types.Int64Value(int64(*v))
	default:
		// unreachable because of the constraint, but keep safe fallback
		return types.Int64Null()
	}
}

func SetInt32PointerToState[T numericPtr](i T) types.Int32 {
	if i == nil {
		return types.Int32Null()
	}
	switch v := any(i).(type) {
	case *int:
		return types.Int32Value(int32(*v))
	case *int32:
		return types.Int32Value(int32(*v))
	case *int64:
		return types.Int32Value(int32(*v))
	case *float32:
		return types.Int32Value(int32(*v))
	case *float64:
		return types.Int32Value(int32(*v))
	default:
		// unreachable because of the constraint, but keep safe fallback
		return types.Int32Null()
	}
}
func SetFloat32PointerToState[T numericPtr](i T) types.Float32 {
	if i == nil {
		return types.Float32Null()
	}
	switch v := any(i).(type) {
	case *int:
		return types.Float32Value(float32(*v))
	case *int32:
		return types.Float32Value(float32(*v))
	case *int64:
		return types.Float32Value(float32(*v))
	case *float32:
		return types.Float32Value(RoundFloat32To2Decimals(*v))
	case *float64:
		return types.Float32Value(float32(RoundFloat64To2Decimals(*v)))
	default:
		// unreachable because of the constraint, but keep safe fallback
		return types.Float32Null()
	}
}

func SetFloat64PointerToState[T numericPtr](i T) types.Float64 {
	if i == nil {
		return types.Float64Null()
	}
	switch v := any(i).(type) {
	case *int:
		return types.Float64Value(float64(*v))
	case *int32:
		return types.Float64Value(float64(*v))
	case *int64:
		return types.Float64Value(float64(*v))
	case *float32:
		return types.Float64Value(RoundFloat64To2Decimals(float64(*v)))
	case *float64:
		return types.Float64Value(RoundFloat64To2Decimals(*v))
	default:
		// unreachable because of the constraint, but keep safe fallback
		return types.Float64Null()
	}
}

func SetBoolPointerToState(i *bool) types.Bool {
	if i == nil {
		return types.BoolNull()
	} else {
		return types.BoolValue(*i)
	}
}

func SetStringPointerFromState(s types.String) *string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	v := s.ValueString()
	return &v
}

// SetBoolPointerFromState converts types.Bool to *bool for API calls
// Returns nil if the value is null or unknown, otherwise returns a pointer to the boolean value
func SetBoolPointerFromState(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	val := v.ValueBool()
	return &val
}

// SetInt32PointerToInt64State converts an int32 pointer to types.Int64
func SetInt32PointerToInt64State(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*i))
}

// canonicalizeJSON returns a compact, deterministic JSON string or an error.
func CanonicalizeJSON(input string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	// json.Marshal produces deterministic ordering for map keys in Go
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
