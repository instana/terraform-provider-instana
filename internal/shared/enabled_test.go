package shared

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type enabledTestModel struct {
	ID      types.String `tfsdk:"id"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func enabledTestSchema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"id":      schema.StringAttribute{Computed: true},
		"enabled": schema.BoolAttribute{Optional: true, Computed: true},
	}}
}

func enabledTestPlan(t *testing.T, enabled types.Bool) *tfsdk.Plan {
	plan := &tfsdk.Plan{Schema: enabledTestSchema()}
	require.False(t, plan.Set(context.Background(), enabledTestModel{ID: types.StringValue("abc"), Enabled: enabled}).HasError())
	return plan
}

func enabledTestState(t *testing.T, enabled types.Bool) *tfsdk.State {
	state := &tfsdk.State{Schema: enabledTestSchema()}
	require.False(t, state.Set(context.Background(), enabledTestModel{ID: types.StringValue("abc"), Enabled: enabled}).HasError())
	return state
}

func boolPtr(b bool) *bool { return &b }

func TestEnabledToModel(t *testing.T) {
	assert.Equal(t, types.BoolValue(true), EnabledToModel(boolPtr(true), types.BoolValue(false)))
	assert.Equal(t, types.BoolValue(false), EnabledToModel(boolPtr(false), types.BoolValue(true)))
	assert.Equal(t, types.BoolValue(false), EnabledToModel(nil, types.BoolValue(false)))
	assert.Equal(t, types.BoolValue(true), EnabledToModel(nil, types.BoolNull()))
	assert.Equal(t, types.BoolValue(true), EnabledToModel(nil, types.BoolUnknown()))
}

func TestSyncEnabled(t *testing.T) {
	put := rest.ActionMethodPut
	tests := []struct {
		name          string
		desired       types.Bool
		prior         *types.Bool
		apiValue      *bool
		expectedCalls []testutils.ActionCall
		expected      bool
	}{
		{name: "create default enabled makes no call", desired: types.BoolValue(true), expected: true},
		{name: "create null desired falls back to default", desired: types.BoolNull(), expected: true},
		{name: "create unknown desired falls back to default", desired: types.BoolUnknown(), expected: true},
		{name: "create disabled calls disable", desired: types.BoolValue(false), expectedCalls: []testutils.ActionCall{{Method: put, ID: "abc", Action: "disable"}}, expected: false},
		{name: "update unchanged makes no call", desired: types.BoolValue(false), prior: ptr(types.BoolValue(false)), expected: false},
		{name: "update enabled to disabled calls disable", desired: types.BoolValue(false), prior: ptr(types.BoolValue(true)), expectedCalls: []testutils.ActionCall{{Method: put, ID: "abc", Action: "disable"}}, expected: false},
		{name: "update disabled to enabled calls enable", desired: types.BoolValue(true), prior: ptr(types.BoolValue(false)), expectedCalls: []testutils.ActionCall{{Method: put, ID: "abc", Action: "enable"}}, expected: true},
		{name: "null prior state is treated as enabled", desired: types.BoolValue(false), prior: ptr(types.BoolNull()), expectedCalls: []testutils.ActionCall{{Method: put, ID: "abc", Action: "disable"}}, expected: false},
		{name: "api value wins over prior state", desired: types.BoolValue(true), prior: ptr(types.BoolValue(true)), apiValue: boolPtr(false), expectedCalls: []testutils.ActionCall{{Method: put, ID: "abc", Action: "enable"}}, expected: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fake := testutils.NewFakeActionInstanaAPI(nil)
			params := EnabledSyncParams{
				APIClient: fake, BasePath: "/api/test", ID: "abc", Label: "test", Operation: "update",
				Plan: enabledTestPlan(t, tc.desired), APIValue: tc.apiValue,
			}
			if tc.prior != nil {
				params.PriorState = enabledTestState(t, *tc.prior)
			}

			result, diags := SyncEnabled(context.Background(), params)

			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedCalls, fake.Actions.Calls)
			if len(tc.expectedCalls) > 0 {
				assert.Equal(t, "/api/test", fake.BasePath)
			}
		})
	}
}

func TestSyncEnabled_ActionErrorKeepsCurrentValue(t *testing.T) {
	fake := testutils.NewFakeActionInstanaAPI(errors.New("boom"))

	result, diags := SyncEnabled(context.Background(), EnabledSyncParams{
		APIClient: fake, BasePath: "/api/test", ID: "abc", Label: "test thing", Operation: "create",
		Plan: enabledTestPlan(t, types.BoolValue(false)),
	})

	require.True(t, diags.HasError())
	assert.Contains(t, diags.Errors()[0].Summary(), "disable test thing")
	assert.Contains(t, diags.Errors()[0].Detail(), "boom")
	assert.True(t, result, "state must reflect the unchanged remote value when the call failed")
}

func ptr[T any](v T) *T { return &v }
