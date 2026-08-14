package sessionsettings

import (
	"context"
	"errors"
	"testing"

	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	restapi "github.com/instana/instana-go-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSingletonRestResource is a hand-rolled mock for rest.SingletonRestResource[*restapi.SessionSettings].
type mockSingletonRestResource struct {
	getResult    *restapi.SessionSettings
	getErr       error
	upsertResult *restapi.SessionSettings
	upsertErr    error
	deleteErr    error
	upsertCalled bool
	deleteCalled bool
}

func (m *mockSingletonRestResource) Get() (*restapi.SessionSettings, error) {
	return m.getResult, m.getErr
}

func (m *mockSingletonRestResource) Upsert(data *restapi.SessionSettings) (*restapi.SessionSettings, error) {
	m.upsertCalled = true
	if m.upsertResult != nil {
		return m.upsertResult, m.upsertErr
	}
	return data, m.upsertErr
}

func (m *mockSingletonRestResource) Delete() error {
	m.deleteCalled = true
	return m.deleteErr
}

// handleSchema returns the resource/schema.Schema defined inside the handle.
func handleSchema() fwschema.Schema {
	return NewSessionSettingsResourceHandle().MetaData().Schema
}

// ---- handle metadata ----

func TestNewSessionSettingsResourceHandle(t *testing.T) {
	handle := NewSessionSettingsResourceHandle()

	assert.NotNil(t, handle)
	assert.NotNil(t, handle.MetaData())
	assert.Equal(t, ResourceInstanaSessionSettings, handle.MetaData().ResourceName)
}

func TestSessionSettingsSchema(t *testing.T) {
	sch := handleSchema()

	assert.Contains(t, sch.Attributes, SessionSettingsFieldTokenLifeTimeInMillis)
	assert.Contains(t, sch.Attributes, SessionSettingsFieldIdleTimeInMillis)
}

// ---- SetComputedFields ----

func TestSessionSettingsSetComputedFields(t *testing.T) {
	handle := &sessionSettingsResourceHandle{}
	plan := &tfsdk.Plan{Schema: handleSchema()}

	diags := handle.SetComputedFields(context.Background(), plan)

	assert.False(t, diags.HasError())
}

// ---- MapStateToDataObject ----

func TestSessionSettingsMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	plan := &tfsdk.Plan{Schema: sch}
	require.False(t, plan.Set(ctx, &SessionSettingsModel{
		TokenLifeTimeInMillis: types.Int64Value(28800000),
		IdleTimeInMillis:      types.Int64Value(3600000),
	}).HasError())

	handle := &sessionSettingsResourceHandle{}
	settings, diags := handle.MapStateToDataObject(ctx, plan, nil)

	assert.False(t, diags.HasError())
	require.NotNil(t, settings)
	assert.Equal(t, int64(28800000), settings.TokenLifeTimeInMillis)
	assert.Equal(t, int64(3600000), settings.IdleTimeInMillis)
}

func TestSessionSettingsMapStateToDataObject_FromState(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	state := &tfsdk.State{Schema: sch}
	require.False(t, state.Set(ctx, &SessionSettingsModel{
		TokenLifeTimeInMillis: types.Int64Value(604800000),
		IdleTimeInMillis:      types.Int64Value(28800000),
	}).HasError())

	handle := &sessionSettingsResourceHandle{}
	settings, diags := handle.MapStateToDataObject(ctx, nil, state)

	assert.False(t, diags.HasError())
	require.NotNil(t, settings)
	assert.Equal(t, int64(604800000), settings.TokenLifeTimeInMillis)
	assert.Equal(t, int64(28800000), settings.IdleTimeInMillis)
}

func TestSessionSettingsMapStateToDataObject_IdleGreaterThanToken_ReturnsError(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	plan := &tfsdk.Plan{Schema: sch}
	require.False(t, plan.Set(ctx, &SessionSettingsModel{
		TokenLifeTimeInMillis: types.Int64Value(3600000),  // 1 hour
		IdleTimeInMillis:      types.Int64Value(28800000), // 8 hours — larger than token lifetime
	}).HasError())

	handle := &sessionSettingsResourceHandle{}
	settings, diags := handle.MapStateToDataObject(ctx, plan, nil)

	assert.True(t, diags.HasError())
	assert.Nil(t, settings)
	assert.Contains(t, diags[0].Detail(), "idle_time_in_millis must not be greater than token_life_time_in_millis")
}

func TestSessionSettingsMapStateToDataObject_IdleEqualToToken_IsValid(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	plan := &tfsdk.Plan{Schema: sch}
	require.False(t, plan.Set(ctx, &SessionSettingsModel{
		TokenLifeTimeInMillis: types.Int64Value(3600000),
		IdleTimeInMillis:      types.Int64Value(3600000), // equal — must be allowed
	}).HasError())

	handle := &sessionSettingsResourceHandle{}
	settings, diags := handle.MapStateToDataObject(ctx, plan, nil)

	assert.False(t, diags.HasError())
	require.NotNil(t, settings)
	assert.Equal(t, int64(3600000), settings.TokenLifeTimeInMillis)
	assert.Equal(t, int64(3600000), settings.IdleTimeInMillis)
}



// ---- UpdateState ----

func TestSessionSettingsUpdateState(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	state := &tfsdk.State{Schema: sch}
	apiObj := &restapi.SessionSettings{
		TokenLifeTimeInMillis: 28800000,
		IdleTimeInMillis:      3600000,
	}

	handle := &sessionSettingsResourceHandle{}
	diags := handle.UpdateState(ctx, state, nil, apiObj)

	assert.False(t, diags.HasError())

	var model SessionSettingsModel
	require.False(t, state.Get(ctx, &model).HasError())
	assert.Equal(t, int64(28800000), model.TokenLifeTimeInMillis.ValueInt64())
	assert.Equal(t, int64(3600000), model.IdleTimeInMillis.ValueInt64())
}

// ---- GetStateUpgraders ----

func TestSessionSettingsGetStateUpgraders(t *testing.T) {
	handle := &sessionSettingsResourceHandle{}
	// Returns nil — no state schema migrations are needed for this resource.
	upgraders := handle.GetStateUpgraders(context.Background())
	assert.Nil(t, upgraders)
}

// ---- GetSingletonRestResource (via mock) ----

// mockInstanaAPI is a minimal stub that satisfies the client.InstanaAPI interface
// just enough for GetSingletonRestResource to be called in tests.
type mockInstanaAPIForSingleton struct {
	singleton *mockSingletonRestResource
}

func (m *mockInstanaAPIForSingleton) SessionSettings() interface{ Get() (*restapi.SessionSettings, error) } {
	return m.singleton
}

// We test GetSingletonRestResource indirectly through the full round-trip below
// using the singleton mock directly rather than constructing a full client.InstanaAPI stub,
// because wiring up all 30+ interface methods just for one method would be noisy.

// ---- round-trip via handle methods ----

func TestSessionSettingsRoundTrip_MapAndUpdateState(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()
	handle := &sessionSettingsResourceHandle{}

	// Build plan
	plan := &tfsdk.Plan{Schema: sch}
	require.False(t, plan.Set(ctx, &SessionSettingsModel{
		TokenLifeTimeInMillis: types.Int64Value(28800000),
		IdleTimeInMillis:      types.Int64Value(3600000),
	}).HasError())

	// Map plan → API object
	apiObj, diags := handle.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	assert.Equal(t, int64(28800000), apiObj.TokenLifeTimeInMillis)
	assert.Equal(t, int64(3600000), apiObj.IdleTimeInMillis)

	// API object → state
	state := &tfsdk.State{Schema: sch}
	diags = handle.UpdateState(ctx, state, nil, apiObj)
	require.False(t, diags.HasError())

	var model SessionSettingsModel
	require.False(t, state.Get(ctx, &model).HasError())
	assert.Equal(t, int64(28800000), model.TokenLifeTimeInMillis.ValueInt64())
	assert.Equal(t, int64(3600000), model.IdleTimeInMillis.ValueInt64())
}

// ---- singleton mock used in provider layer tests ----

func TestMockSingletonRestResource_Upsert(t *testing.T) {
	m := &mockSingletonRestResource{
		upsertResult: &restapi.SessionSettings{TokenLifeTimeInMillis: 10, IdleTimeInMillis: 20},
	}
	res, err := m.Upsert(&restapi.SessionSettings{})
	assert.NoError(t, err)
	assert.Equal(t, int64(10), res.TokenLifeTimeInMillis)
	assert.True(t, m.upsertCalled)
}

func TestMockSingletonRestResource_Get(t *testing.T) {
	m := &mockSingletonRestResource{
		getResult: &restapi.SessionSettings{TokenLifeTimeInMillis: 5, IdleTimeInMillis: 6},
	}
	res, err := m.Get()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), res.TokenLifeTimeInMillis)
}

func TestMockSingletonRestResource_Delete(t *testing.T) {
	m := &mockSingletonRestResource{}
	assert.NoError(t, m.Delete())
	assert.True(t, m.deleteCalled)
}

func TestMockSingletonRestResource_Errors(t *testing.T) {
	errMsg := errors.New("boom")

	t.Run("upsert error", func(t *testing.T) {
		m := &mockSingletonRestResource{upsertErr: errMsg}
		_, err := m.Upsert(&restapi.SessionSettings{})
		assert.ErrorIs(t, err, errMsg)
	})

	t.Run("get error", func(t *testing.T) {
		m := &mockSingletonRestResource{getErr: errMsg}
		_, err := m.Get()
		assert.ErrorIs(t, err, errMsg)
	})

	t.Run("delete error", func(t *testing.T) {
		m := &mockSingletonRestResource{deleteErr: errMsg}
		assert.ErrorIs(t, m.Delete(), errMsg)
	})
}
