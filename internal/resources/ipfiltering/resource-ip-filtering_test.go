package ipfiltering

import (
	"context"
	"testing"

	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	restapi "github.com/instana/instana-go-client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// handleSchema returns the resource/schema.Schema defined inside the handle.
func handleSchema() fwschema.Schema {
	return NewIPFilteringResourceHandle().MetaData().Schema
}

// ---- handle metadata ----

func TestNewIPFilteringResourceHandle(t *testing.T) {
	handle := NewIPFilteringResourceHandle()

	assert.NotNil(t, handle)
	assert.NotNil(t, handle.MetaData())
	assert.Equal(t, ResourceInstanaIPFiltering, handle.MetaData().ResourceName)
}

func TestIPFilteringSchema(t *testing.T) {
	sch := handleSchema()

	assert.Contains(t, sch.Attributes, IPFilteringFieldEnabled)
	assert.Contains(t, sch.Attributes, IPFilteringFieldDenyAll)
	assert.Contains(t, sch.Attributes, IPFilteringFieldSupportAccessEnabled)
	assert.Contains(t, sch.Attributes, IPFilteringFieldRules)
	assert.Contains(t, sch.Attributes, IPFilteringFieldActive)
	assert.Contains(t, sch.Attributes, IPFilteringFieldLastChangedAt)
	assert.Contains(t, sch.Attributes, IPFilteringFieldLastVerifiedAt)
}

// ---- SetComputedFields ----

func TestIPFilteringSetComputedFields(t *testing.T) {
	handle := &ipFilteringResourceHandle{}
	plan := &tfsdk.Plan{Schema: handleSchema()}

	diags := handle.SetComputedFields(context.Background(), plan)

	assert.False(t, diags.HasError())
}

// ---- MapStateToDataObject ----

func TestIPFilteringMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	plan := &tfsdk.Plan{Schema: sch}
	require.False(t, plan.Set(ctx, &IPFilteringModel{
		Enabled:              types.BoolValue(true),
		DenyAll:              types.BoolValue(false),
		SupportAccessEnabled: types.BoolValue(true),
		Rules: []IPFilteringRuleModel{
			{Target: types.StringValue("192.168.1.1/32"), Block: types.BoolValue(false)},
			{Target: types.StringValue("10.0.0.0/8"), Block: types.BoolValue(true)},
		},
	}).HasError())

	handle := &ipFilteringResourceHandle{}
	filtering, diags := handle.MapStateToDataObject(ctx, plan, nil)

	assert.False(t, diags.HasError())
	require.NotNil(t, filtering)
	assert.True(t, filtering.Enabled)
	assert.False(t, filtering.DenyAll)
	assert.True(t, filtering.SupportAccessEnabled)
	require.Len(t, filtering.Rules, 2)
	assert.Equal(t, "192.168.1.1/32", filtering.Rules[0].Target)
	assert.False(t, filtering.Rules[0].Block)
	assert.Equal(t, "10.0.0.0/8", filtering.Rules[1].Target)
	assert.True(t, filtering.Rules[1].Block)
}

func TestIPFilteringMapStateToDataObject_FromState(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	state := &tfsdk.State{Schema: sch}
	require.False(t, state.Set(ctx, &IPFilteringModel{
		Enabled:              types.BoolValue(false),
		DenyAll:              types.BoolValue(true),
		SupportAccessEnabled: types.BoolValue(false),
		Rules: []IPFilteringRuleModel{
			{Target: types.StringValue("127.0.0.1"), Block: types.BoolValue(true)},
		},
	}).HasError())

	handle := &ipFilteringResourceHandle{}
	filtering, diags := handle.MapStateToDataObject(ctx, nil, state)

	assert.False(t, diags.HasError())
	require.NotNil(t, filtering)
	assert.False(t, filtering.Enabled)
	assert.True(t, filtering.DenyAll)
	assert.False(t, filtering.SupportAccessEnabled)
	require.Len(t, filtering.Rules, 1)
	assert.Equal(t, "127.0.0.1", filtering.Rules[0].Target)
	assert.True(t, filtering.Rules[0].Block)
}

// ---- UpdateState ----

func TestIPFilteringUpdateState(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	state := &tfsdk.State{Schema: sch}
	lastChanged := int64(1762356274224)
	lastVerified := int64(1762356280000)
	apiObj := &restapi.IPFiltering{
		Active:               true,
		DenyAll:              false,
		Enabled:              true,
		SupportAccessEnabled: true,
		Rules: []restapi.IPFilteringRule{
			{Target: "127.0.0.1", Block: true},
			{Target: "0:0:0:0:0:0:0:1", Block: false},
		},
		LastChangedAt:  &lastChanged,
		LastVerifiedAt: &lastVerified,
	}

	handle := &ipFilteringResourceHandle{}
	diags := handle.UpdateState(ctx, state, nil, apiObj)

	assert.False(t, diags.HasError())

	var model IPFilteringModel
	require.False(t, state.Get(ctx, &model).HasError())
	assert.True(t, model.Enabled.ValueBool())
	assert.False(t, model.DenyAll.ValueBool())
	assert.True(t, model.SupportAccessEnabled.ValueBool())
	assert.True(t, model.Active.ValueBool())
	assert.Equal(t, lastChanged, model.LastChangedAt.ValueInt64())
	assert.Equal(t, lastVerified, model.LastVerifiedAt.ValueInt64())
	require.Len(t, model.Rules, 2)
	assert.Equal(t, "127.0.0.1", model.Rules[0].Target.ValueString())
	assert.True(t, model.Rules[0].Block.ValueBool())
	assert.Equal(t, "0:0:0:0:0:0:0:1", model.Rules[1].Target.ValueString())
	assert.False(t, model.Rules[1].Block.ValueBool())
}

func TestIPFilteringUpdateState_NilTimestamps(t *testing.T) {
	ctx := context.Background()
	sch := handleSchema()

	state := &tfsdk.State{Schema: sch}
	apiObj := &restapi.IPFiltering{
		Active:               false,
		DenyAll:              true,
		Enabled:              false,
		SupportAccessEnabled: false,
		Rules: []restapi.IPFilteringRule{
			{Target: "10.0.0.0/8", Block: false},
		},
		LastChangedAt:  nil,
		LastVerifiedAt: nil,
	}

	handle := &ipFilteringResourceHandle{}
	diags := handle.UpdateState(ctx, state, nil, apiObj)

	assert.False(t, diags.HasError())

	var model IPFilteringModel
	require.False(t, state.Get(ctx, &model).HasError())
	assert.True(t, model.LastChangedAt.IsNull())
	assert.True(t, model.LastVerifiedAt.IsNull())
}

// ---- PostUpsertVerifier ----

func TestIPFilteringNeedsPostUpsertVerification(t *testing.T) {
	handle := &ipFilteringResourceHandle{}

	assert.False(t, handle.NeedsPostUpsertVerification(nil))
	assert.False(t, handle.NeedsPostUpsertVerification(&restapi.IPFiltering{Enabled: false}))
	assert.True(t, handle.NeedsPostUpsertVerification(&restapi.IPFiltering{Enabled: true}))
}

// ---- GetStateUpgraders ----

func TestIPFilteringGetStateUpgraders(t *testing.T) {
	handle := &ipFilteringResourceHandle{}
	upgraders := handle.GetStateUpgraders(context.Background())
	assert.Nil(t, upgraders)
}
