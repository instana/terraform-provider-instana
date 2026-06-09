package apitoken

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestAlwaysFalseAPITokenPlanModifier_TrueValue_AddsError(t *testing.T) {
	modifier := AlwaysFalseAPITokenPlanModifier{Field: "can_configure_personal_api_tokens"}
	ctx := context.Background()
	request := planmodifier.BoolRequest{
		PlanValue: types.BoolValue(true),
	}
	response := &planmodifier.BoolResponse{}

	modifier.PlanModifyBool(ctx, request, response)

	assert.True(t, response.Diagnostics.HasError(), "Expected error when PlanValue is true")
	if len(response.Diagnostics) > 0 {
		assert.Contains(t, response.Diagnostics[0].Summary(), "cannot be set to true")
	}
}

func TestAlwaysFalseAPITokenPlanModifier_FalseValue_NoError(t *testing.T) {
	modifier := AlwaysFalseAPITokenPlanModifier{Field: "can_configure_personal_api_tokens"}
	ctx := context.Background()
	request := planmodifier.BoolRequest{
		PlanValue: types.BoolValue(false),
	}
	response := &planmodifier.BoolResponse{}

	modifier.PlanModifyBool(ctx, request, response)

	assert.False(t, response.Diagnostics.HasError(), "Expected no error when PlanValue is false")
}

func TestAlwaysFalseAPITokenPlanModifier_UnknownValue(t *testing.T) {
	modifier := AlwaysFalseAPITokenPlanModifier{Field: "can_configure_personal_api_tokens"}
	ctx := context.Background()
	request := planmodifier.BoolRequest{
		PlanValue: types.BoolUnknown(),
	}
	response := &planmodifier.BoolResponse{}

	modifier.PlanModifyBool(ctx, request, response)

	assert.False(t, response.Diagnostics.HasError(), "Expected no error for unknown value")
}

func TestAlwaysFalseAPITokenPlanModifier_NilPlanValue(t *testing.T) {
	modifier := AlwaysFalseAPITokenPlanModifier{Field: "can_configure_personal_api_tokens"}
	ctx := context.Background()
	request := planmodifier.BoolRequest{}
	response := &planmodifier.BoolResponse{}

	modifier.PlanModifyBool(ctx, request, response)

	assert.False(t, response.Diagnostics.HasError(), "Expected no error for nil PlanValue (zero value)")
}

func TestAlwaysFalseAPITokenPlanModifier_AppendDiagnostics(t *testing.T) {
	modifier := AlwaysFalseAPITokenPlanModifier{Field: "can_configure_personal_api_tokens"}
	ctx := context.Background()
	request := planmodifier.BoolRequest{
		PlanValue: types.BoolValue(true),
	}
	response := &planmodifier.BoolResponse{}
	response.Diagnostics.AddWarning("existing warning", "existing warning detail")

	modifier.PlanModifyBool(ctx, request, response)

	assert.True(t, response.Diagnostics.HasError(), "Expected error when PlanValue is true")
	// Check that the warning is still present (should be 2 diagnostics: 1 warning, 1 error)
	warningCount := 0
	errorCount := 0
	for _, d := range response.Diagnostics {
		if d.Severity() == 1 { // 1 = Warning
			warningCount++
		}
		if d.Severity() == 2 { // 2 = Error
			errorCount++
		}
	}
	assert.Equal(t, 1, warningCount, "Expected 1 warning diagnostic")
	assert.Equal(t, 1, errorCount, "Expected 1 error diagnostic")
}
