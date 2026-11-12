package infralertconfig

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

type actionCall struct {
	method rest.ActionMethod
	id     string
	action string
}

type fakeActionRestResource struct {
	calls []actionCall
	err   error
}

func (f *fakeActionRestResource) Execute(method rest.ActionMethod, id string, action string) error {
	f.calls = append(f.calls, actionCall{method: method, id: id, action: action})
	return f.err
}

type fakeInstanaAPI struct {
	testutils.MockInstanaAPI
	actions  *fakeActionRestResource
	basePath string
}

func (f *fakeInstanaAPI) ResourceActions(basePath string) rest.ActionRestResource {
	f.basePath = basePath
	return f.actions
}

func enabledModel(enabled types.Bool) InfraAlertConfigModel {
	return InfraAlertConfigModel{
		ID:                 types.StringValue("test-id"),
		Name:               types.StringValue("Test Infra Alert"),
		Description:        types.StringValue("Test Description"),
		Granularity:        types.Int64Value(600000),
		EvaluationType:     types.StringValue(string(api.EvaluationTypePerEntity)),
		TagFilter:          types.StringNull(),
		GroupBy:            types.SetNull(types.StringType),
		CustomPayloadField: types.ListNull(shared.GetCustomPayloadFieldType()),
		Enabled:            enabled,
	}
}

func planWithEnabled(t *testing.T, enabled types.Bool) *tfsdk.Plan {
	plan := &tfsdk.Plan{Schema: getTestSchema()}
	diags := plan.Set(context.Background(), enabledModel(enabled))
	require.False(t, diags.HasError(), "%v", diags)
	return plan
}

func stateWithEnabled(t *testing.T, enabled types.Bool) *tfsdk.State {
	state := createMockState(t, enabledModel(enabled))
	return &state
}

func boolPtr(b bool) *bool { return &b }

func TestSchema_EnabledAttribute(t *testing.T) {
	attr, ok := getTestSchema().Attributes[InfraAlertConfigFieldEnabled]
	require.True(t, ok, "enabled attribute must exist in the schema")

	boolAttr, ok := attr.(schema.BoolAttribute)
	require.True(t, ok)
	assert.True(t, boolAttr.Optional)
	assert.True(t, boolAttr.Computed)
	assert.NotNil(t, boolAttr.Default)
}

func TestMapStateToDataObject_DoesNotSendEnabled(t *testing.T) {
	ctx := context.Background()
	resource := NewInfraAlertConfigResourceHandle()
	plan := planWithEnabled(t, types.BoolValue(false))

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	assert.Nil(t, result.Enabled, "enabled must not be part of the create/update payload")
}

func TestUpdateState_Enabled(t *testing.T) {
	tests := []struct {
		name     string
		apiValue *bool
		initial  types.Bool
		expected bool
	}{
		{"api returns true", boolPtr(true), types.BoolValue(false), true},
		{"api returns false", boolPtr(false), types.BoolValue(true), false},
		{"api returns nothing keeps planned value", nil, types.BoolValue(false), false},
		{"api returns nothing and null value defaults to true", nil, types.BoolNull(), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			resource := NewInfraAlertConfigResourceHandle()
			state := stateWithEnabled(t, tc.initial)

			data := &api.InfraAlertConfig{
				ID:             "test-id",
				Name:           "Test Infra Alert",
				Description:    "Test Description",
				Granularity:    600000,
				EvaluationType: api.EvaluationTypePerEntity,
				Enabled:        tc.apiValue,
			}

			diags := resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError(), "%v", diags)

			var model InfraAlertConfigModel
			require.False(t, state.Get(ctx, &model).HasError())
			assert.Equal(t, tc.expected, model.Enabled.ValueBool())
		})
	}
}

func TestPostWrite(t *testing.T) {
	tests := []struct {
		name          string
		op            resourcehandle.Operation
		desired       types.Bool
		prior         *types.Bool
		apiEnabled    *bool
		expectedCalls []actionCall
		expectedState bool
	}{
		{
			name:          "create with default enabled makes no extra call",
			op:            resourcehandle.OperationCreate,
			desired:       types.BoolValue(true),
			expectedState: true,
		},
		{
			name:          "create with null enabled falls back to default and makes no call",
			op:            resourcehandle.OperationCreate,
			desired:       types.BoolNull(),
			expectedState: true,
		},
		{
			name:          "create disabled calls disable",
			op:            resourcehandle.OperationCreate,
			desired:       types.BoolValue(false),
			expectedCalls: []actionCall{{rest.ActionMethodPut, "test-id", "disable"}},
			expectedState: false,
		},
		{
			name:          "update unchanged makes no call",
			op:            resourcehandle.OperationUpdate,
			desired:       types.BoolValue(false),
			prior:         func() *types.Bool { v := types.BoolValue(false); return &v }(),
			expectedState: false,
		},
		{
			name:          "update from enabled to disabled calls disable",
			op:            resourcehandle.OperationUpdate,
			desired:       types.BoolValue(false),
			prior:         func() *types.Bool { v := types.BoolValue(true); return &v }(),
			expectedCalls: []actionCall{{rest.ActionMethodPut, "test-id", "disable"}},
			expectedState: false,
		},
		{
			name:          "update from disabled to enabled calls enable",
			op:            resourcehandle.OperationUpdate,
			desired:       types.BoolValue(true),
			prior:         func() *types.Bool { v := types.BoolValue(false); return &v }(),
			expectedCalls: []actionCall{{rest.ActionMethodPut, "test-id", "enable"}},
			expectedState: true,
		},
		{
			name:          "value returned by api wins over prior state",
			op:            resourcehandle.OperationUpdate,
			desired:       types.BoolValue(true),
			prior:         func() *types.Bool { v := types.BoolValue(true); return &v }(),
			apiEnabled:    boolPtr(false),
			expectedCalls: []actionCall{{rest.ActionMethodPut, "test-id", "enable"}},
			expectedState: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &infraAlertConfigResource{}
			fake := &fakeInstanaAPI{actions: &fakeActionRestResource{}}
			plan := planWithEnabled(t, tc.desired)
			var prior *tfsdk.State
			if tc.prior != nil {
				prior = stateWithEnabled(t, *tc.prior)
			}
			written := &api.InfraAlertConfig{ID: "test-id", Enabled: tc.apiEnabled}

			result, diags := resource.PostWrite(ctx, fake, tc.op, plan, prior, written)

			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, tc.expectedCalls, fake.actions.calls)
			require.NotNil(t, result.Enabled)
			assert.Equal(t, tc.expectedState, *result.Enabled)
			if len(tc.expectedCalls) > 0 {
				assert.Equal(t, api.InfraAlertConfigResourcePath, fake.basePath)
			}
		})
	}
}

func TestPostWrite_ActionError(t *testing.T) {
	ctx := context.Background()
	resource := &infraAlertConfigResource{}
	fake := &fakeInstanaAPI{actions: &fakeActionRestResource{err: errors.New("boom")}}
	plan := planWithEnabled(t, types.BoolValue(false))
	written := &api.InfraAlertConfig{ID: "test-id"}

	result, diags := resource.PostWrite(ctx, fake, resourcehandle.OperationCreate, plan, nil, written)

	require.True(t, diags.HasError())
	assert.Contains(t, diags.Errors()[0].Detail(), "boom")
	require.NotNil(t, result)
	require.NotNil(t, result.Enabled)
	assert.True(t, *result.Enabled, "state must reflect the unchanged remote value when the call failed")
}

func TestInfraAlertConfigResourceImplementsPostWriteHook(t *testing.T) {
	var handle resourcehandle.ResourceHandle[*api.InfraAlertConfig] = NewInfraAlertConfigResourceHandle()
	_, ok := handle.(resourcehandle.PostWriteHook[*api.InfraAlertConfig])
	assert.True(t, ok)
}
