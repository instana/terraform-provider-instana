package release

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyApplicationsList returns an empty types.List with the correct element type.
func emptyApplicationsList() types.List {
	return types.ListValueMust(types.ObjectType{AttrTypes: applicationAttrTypes}, []attr.Value{})
}

// emptyServicesList returns an empty types.List with the correct element type.
func emptyServicesList() types.List {
	return types.ListValueMust(types.ObjectType{AttrTypes: serviceAttrTypes}, []attr.Value{})
}

func TestNewReleaseResourceHandle(t *testing.T) {
	handle := NewReleaseResourceHandle()
	require.NotNil(t, handle)

	meta := handle.MetaData()
	assert.Equal(t, ResourceInstanaRelease, meta.ResourceName)
	assert.NotNil(t, meta.Schema)
	assert.Equal(t, int64(0), meta.SchemaVersion)
	assert.False(t, meta.CreateOnly)
}

func TestReleaseMetaData(t *testing.T) {
	r := &releaseResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_release",
			SchemaVersion: 0,
		},
	}
	meta := r.MetaData()
	assert.Equal(t, "test_release", meta.ResourceName)
}

func TestReleaseSetComputedFields(t *testing.T) {
	r := NewReleaseResourceHandle()
	ctx := context.Background()
	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	diags := r.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

func TestReleaseGetStateUpgraders(t *testing.T) {
	r := NewReleaseResourceHandle()
	upgraders := r.GetStateUpgraders(context.Background())
	assert.Nil(t, upgraders)
}

func TestReleaseUpdateState_BasicRelease(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	release := &api.ReleaseWithMetadata{
		ID:          "Tiu16hLCTniHDtHb_uDV1w",
		Name:        "demo-app/main-**",
		Start:       1709091782000,
		LastUpdated: 1709091782533,
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	initialModel := ReleaseModel{
		ID:           types.StringValue(""),
		Name:         types.StringValue(""),
		Start:        types.Int64Value(0),
		LastUpdated:  types.Int64Value(0),
		Applications: emptyApplicationsList(),
		Services:     emptyServicesList(),
	}
	require.False(t, state.Set(ctx, initialModel).HasError())

	diags := r.UpdateState(ctx, state, nil, release)
	require.False(t, diags.HasError())

	var model ReleaseModel
	require.False(t, state.Get(ctx, &model).HasError())

	assert.Equal(t, "Tiu16hLCTniHDtHb_uDV1w", model.ID.ValueString())
	assert.Equal(t, "demo-app/main-**", model.Name.ValueString())
	assert.Equal(t, int64(1709091782000), model.Start.ValueInt64())
	assert.Equal(t, int64(1709091782533), model.LastUpdated.ValueInt64())
	assert.Empty(t, model.Applications.Elements())
	assert.Empty(t, model.Services.Elements())
}

func TestReleaseUpdateState_WithApplicationsAndServices(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	release := &api.ReleaseWithMetadata{
		ID:          "XK1e1TF3T9SHKugndn_soQ",
		Name:        "frontend/release-2000",
		Start:       1706674621000,
		LastUpdated: 1706674621604,
		Applications: []*api.ReleaseApplicationScope{
			{Name: "app1"},
		},
		Services: []*api.ReleaseServiceScope{
			{
				Name: "payment",
				ScopedTo: &api.ReleaseServiceScopedTo{
					ApplicationName: "checkout-app",
				},
			},
		},
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	initialModel := ReleaseModel{
		ID:           types.StringValue(""),
		Name:         types.StringValue(""),
		Start:        types.Int64Value(0),
		LastUpdated:  types.Int64Value(0),
		Applications: emptyApplicationsList(),
		Services:     emptyServicesList(),
	}
	require.False(t, state.Set(ctx, initialModel).HasError())

	diags := r.UpdateState(ctx, state, nil, release)
	require.False(t, diags.HasError())

	var model ReleaseModel
	require.False(t, state.Get(ctx, &model).HasError())

	var apps []ApplicationModel
	require.False(t, model.Applications.ElementsAs(ctx, &apps, false).HasError())
	require.Len(t, apps, 1)
	assert.Equal(t, "app1", apps[0].Name.ValueString())

	var svcs []ServiceModel
	require.False(t, model.Services.ElementsAs(ctx, &svcs, false).HasError())
	require.Len(t, svcs, 1)
	assert.Equal(t, "payment", svcs[0].Name.ValueString())
	require.NotNil(t, svcs[0].ScopedTo)
	assert.Equal(t, "checkout-app", svcs[0].ScopedTo.ApplicationName.ValueString())
}

func TestReleaseMapStateToDataObject_BasicRelease(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := ReleaseModel{
		ID:           types.StringValue("release-id-1"),
		Name:         types.StringValue("frontend/release-1000"),
		Start:        types.Int64Value(1742349976000),
		LastUpdated:  types.Int64Value(0),
		Applications: emptyApplicationsList(),
		Services:     emptyServicesList(),
	}

	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	require.False(t, plan.Set(ctx, model).HasError())

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, release)

	assert.Equal(t, "release-id-1", release.ID)
	assert.Equal(t, "frontend/release-1000", release.Name)
	assert.Equal(t, int64(1742349976000), release.Start)
	assert.Nil(t, release.Applications)
	assert.Nil(t, release.Services)
}

func TestReleaseMapStateToDataObject_WithApplicationsAndServices(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	appObj, d := types.ObjectValue(applicationAttrTypes, map[string]attr.Value{
		ReleaseFieldName: types.StringValue("my-app"),
	})
	require.False(t, d.HasError())

	scopedToObj, d := types.ObjectValue(scopedToAttrTypes, map[string]attr.Value{
		ReleaseFieldApplicationName: types.StringValue("my-app"),
		ReleaseFieldEnvironmentName: types.StringValue("production"),
	})
	require.False(t, d.HasError())

	svcObj, d := types.ObjectValue(serviceAttrTypes, map[string]attr.Value{
		ReleaseFieldName:     types.StringValue("my-service"),
		ReleaseFieldScopedTo: scopedToObj,
	})
	require.False(t, d.HasError())

	appsList, d := types.ListValue(types.ObjectType{AttrTypes: applicationAttrTypes}, []attr.Value{appObj})
	require.False(t, d.HasError())

	svcsList, d := types.ListValue(types.ObjectType{AttrTypes: serviceAttrTypes}, []attr.Value{svcObj})
	require.False(t, d.HasError())

	model := ReleaseModel{
		ID:           types.StringValue("release-id-2"),
		Name:         types.StringValue("backend/release-5"),
		Start:        types.Int64Value(1742349976000),
		LastUpdated:  types.Int64Value(0),
		Applications: appsList,
		Services:     svcsList,
	}

	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	require.False(t, plan.Set(ctx, model).HasError())

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, release)

	require.Len(t, release.Applications, 1)
	assert.Equal(t, "my-app", release.Applications[0].Name)

	require.Len(t, release.Services, 1)
	assert.Equal(t, "my-service", release.Services[0].Name)
	require.NotNil(t, release.Services[0].ScopedTo)
	assert.Equal(t, "my-app", release.Services[0].ScopedTo.ApplicationName)
	assert.Equal(t, "production", release.Services[0].ScopedTo.EnvironmentName)
}

func TestReleaseMapStateToDataObject_WithStateWhenNoPlan(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := ReleaseModel{
		ID:           types.StringValue("release-id-3"),
		Name:         types.StringValue("state-read-release"),
		Start:        types.Int64Value(1742349976000),
		LastUpdated:  types.Int64Value(0),
		Applications: emptyApplicationsList(),
		Services:     emptyServicesList(),
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	require.False(t, state.Set(ctx, model).HasError())

	release, diags := r.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	require.NotNil(t, release)

	assert.Equal(t, "state-read-release", release.Name)
}

func TestMapApplicationsToState_Empty(t *testing.T) {
	ctx := context.Background()
	result, diags := mapApplicationsToState(ctx, nil)
	require.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestMapServicesToState_Empty(t *testing.T) {
	ctx := context.Background()
	result, diags := mapServicesToState(ctx, nil)
	require.False(t, diags.HasError())
	assert.False(t, result.IsNull())
	assert.Empty(t, result.Elements())
}

func TestMapApplicationsFromState_Empty(t *testing.T) {
	ctx := context.Background()
	result, diags := mapApplicationsFromState(ctx, emptyApplicationsList())
	require.False(t, diags.HasError())
	assert.Nil(t, result)
}

func TestMapServicesFromState_Empty(t *testing.T) {
	ctx := context.Background()
	result, diags := mapServicesFromState(ctx, emptyServicesList())
	require.False(t, diags.HasError())
	assert.Nil(t, result)
}

func TestMapServicesToState_NoScopedTo(t *testing.T) {
	ctx := context.Background()
	services := []*api.ReleaseServiceScope{
		{Name: "plain-service"},
	}
	result, diags := mapServicesToState(ctx, services)
	require.False(t, diags.HasError())
	require.Len(t, result.Elements(), 1)

	var svcs []ServiceModel
	require.False(t, result.ElementsAs(ctx, &svcs, false).HasError())
	assert.Equal(t, "plain-service", svcs[0].Name.ValueString())
	assert.Nil(t, svcs[0].ScopedTo)
}
