package sloalertconfig

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// enabledTestModel returns a model with only the enabled flag set; all other values are null
func enabledTestModel(enabled types.Bool) SloAlertConfigModel {
	return SloAlertConfigModel{
		ID:              types.StringNull(),
		Name:            types.StringNull(),
		Description:     types.StringNull(),
		Severity:        types.Int64Null(),
		Triggering:      types.BoolNull(),
		AlertType:       types.StringNull(),
		SloIds:          types.SetNull(types.StringType),
		ApdexIds:        types.SetNull(types.StringType),
		AlertChannelIds: types.SetNull(types.StringType),
		CustomPayload:   types.ListNull(shared.GetCustomPayloadFieldType()),
		Enabled:         enabled,
	}
}

func enabledTestState(t *testing.T, enabled types.Bool) *tfsdk.State {
	state := &tfsdk.State{Schema: getTestSchema()}
	diags := state.Set(context.Background(), enabledTestModel(enabled))
	require.False(t, diags.HasError(), "%v", diags)
	return state
}

func enabledTestPlan(t *testing.T, enabled types.Bool) *tfsdk.Plan {
	state := enabledTestState(t, enabled)
	return &tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}
}

func enabledPtr(b bool) *bool { return &b }

func TestSchema_EnabledAttribute(t *testing.T) {
	attr, ok := getTestSchema().Attributes["enabled"]
	require.True(t, ok, "enabled attribute must exist in the schema")

	boolAttr, ok := attr.(schema.BoolAttribute)
	require.True(t, ok)
	assert.True(t, boolAttr.Optional)
	assert.True(t, boolAttr.Computed)
	require.NotNil(t, boolAttr.Default)
}

func TestImplementsPostWriteHook(t *testing.T) {
	var handle resourcehandle.ResourceHandle[*api.SloAlertConfig] = NewSloAlertConfigResourceHandle()
	_, ok := handle.(resourcehandle.PostWriteHook[*api.SloAlertConfig])
	assert.True(t, ok)
}

func TestPostWrite_DisableOnCreate(t *testing.T) {
	handle := NewSloAlertConfigResourceHandle()
	hook := handle.(resourcehandle.PostWriteHook[*api.SloAlertConfig])
	fake := testutils.NewFakeActionInstanaAPI(nil)
	written := &api.SloAlertConfig{ID: "test-id"}

	result, diags := hook.PostWrite(context.Background(), fake, resourcehandle.OperationCreate, enabledTestPlan(t, types.BoolValue(false)), nil, written)

	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, api.SloAlertConfigResourcePath, fake.BasePath)
	assert.Equal(t, []testutils.ActionCall{{Method: rest.ActionMethodPut, ID: "test-id", Action: "disable"}}, fake.Actions.Calls)
	require.NotNil(t, result.Enabled)
	assert.False(t, *result.Enabled)
}

func TestPostWrite_DefaultEnabledOnCreateMakesNoCall(t *testing.T) {
	handle := NewSloAlertConfigResourceHandle()
	hook := handle.(resourcehandle.PostWriteHook[*api.SloAlertConfig])
	fake := testutils.NewFakeActionInstanaAPI(nil)
	written := &api.SloAlertConfig{ID: "test-id"}

	result, diags := hook.PostWrite(context.Background(), fake, resourcehandle.OperationCreate, enabledTestPlan(t, types.BoolValue(true)), nil, written)

	require.False(t, diags.HasError(), "%v", diags)
	assert.Empty(t, fake.Actions.Calls)
	require.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled)
}

func TestPostWrite_ToggleOnUpdate(t *testing.T) {
	tests := []struct {
		name     string
		prior    bool
		desired  bool
		action   string
		expected bool
	}{
		{"enabled to disabled", true, false, "disable", false},
		{"disabled to enabled", false, true, "enable", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handle := NewSloAlertConfigResourceHandle()
			hook := handle.(resourcehandle.PostWriteHook[*api.SloAlertConfig])
			fake := testutils.NewFakeActionInstanaAPI(nil)
			written := &api.SloAlertConfig{ID: "test-id"}

			result, diags := hook.PostWrite(context.Background(), fake, resourcehandle.OperationUpdate,
				enabledTestPlan(t, types.BoolValue(tc.desired)), enabledTestState(t, types.BoolValue(tc.prior)), written)

			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, []testutils.ActionCall{{Method: rest.ActionMethodPut, ID: "test-id", Action: tc.action}}, fake.Actions.Calls)
			require.NotNil(t, result.Enabled)
			assert.Equal(t, tc.expected, *result.Enabled)
		})
	}
}

func TestPostWrite_UnchangedOnUpdateMakesNoCall(t *testing.T) {
	handle := NewSloAlertConfigResourceHandle()
	hook := handle.(resourcehandle.PostWriteHook[*api.SloAlertConfig])
	fake := testutils.NewFakeActionInstanaAPI(nil)
	written := &api.SloAlertConfig{ID: "test-id"}

	_, diags := hook.PostWrite(context.Background(), fake, resourcehandle.OperationUpdate,
		enabledTestPlan(t, types.BoolValue(false)), enabledTestState(t, types.BoolValue(false)), written)

	require.False(t, diags.HasError(), "%v", diags)
	assert.Empty(t, fake.Actions.Calls)
}

func TestPostWrite_ActionError(t *testing.T) {
	handle := NewSloAlertConfigResourceHandle()
	hook := handle.(resourcehandle.PostWriteHook[*api.SloAlertConfig])
	fake := testutils.NewFakeActionInstanaAPI(errors.New("boom"))
	written := &api.SloAlertConfig{ID: "test-id"}

	result, diags := hook.PostWrite(context.Background(), fake, resourcehandle.OperationCreate, enabledTestPlan(t, types.BoolValue(false)), nil, written)

	require.True(t, diags.HasError())
	assert.Contains(t, diags.Errors()[0].Detail(), "boom")
	require.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled, "state must reflect the unchanged remote value when the call failed")
}

func TestUpdateState_Enabled(t *testing.T) {
	tests := []struct {
		name     string
		apiValue *bool
		initial  types.Bool
		expected bool
	}{
		{"api returns true", enabledPtr(true), types.BoolValue(false), true},
		{"api returns false", enabledPtr(false), types.BoolValue(true), false},
		{"api returns nothing keeps planned value", nil, types.BoolValue(false), false},
		{"api returns nothing and null value defaults to true", nil, types.BoolNull(), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			handle := NewSloAlertConfigResourceHandle()
			state := enabledTestState(t, tc.initial)
			data := &api.SloAlertConfig{ID: "test-id", Name: "n", Description: "d", Severity: 5, Rule: api.SloAlertRule{AlertType: "SERVICE_LEVELS_OBJECTIVE", Metric: "STATUS"}, BurnRateConfigs: &[]api.BurnRateConfig{}}
			data.Enabled = tc.apiValue

			diags := handle.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError(), "%v", diags)

			var model SloAlertConfigModel
			require.False(t, state.Get(ctx, &model).HasError())
			assert.Equal(t, tc.expected, model.Enabled.ValueBool())
		})
	}
}
