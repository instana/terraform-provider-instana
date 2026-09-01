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

// ─── helpers ─────────────────────────────────────────────────────────────────

func emptyReleaseModel() ReleaseModel {
	appElemType := types.ObjectType{AttrTypes: applicationAttrTypes}
	svcElemType := types.ObjectType{AttrTypes: serviceAttrTypes}
	return ReleaseModel{
		ID:           types.StringValue(""),
		Name:         types.StringValue(""),
		Start:        types.Int64Value(0),
		LastUpdated:  types.Int64Value(0),
		Applications: types.ListValueMust(appElemType, []attr.Value{}),
		Services:     types.ListValueMust(svcElemType, []attr.Value{}),
	}
}

func setState(t *testing.T, state *tfsdk.State, model ReleaseModel) {
	t.Helper()
	require.False(t, state.Set(context.Background(), model).HasError())
}

func setPlan(t *testing.T, plan *tfsdk.Plan, model ReleaseModel) {
	t.Helper()
	require.False(t, plan.Set(context.Background(), model).HasError())
}

func getModel(t *testing.T, state *tfsdk.State) ReleaseModel {
	t.Helper()
	var m ReleaseModel
	require.False(t, state.Get(context.Background(), &m).HasError())
	return m
}

// buildScopedToValue creates a types.Object for scoped_to from a list of app names.
func buildScopedToValue(t *testing.T, appNames []string) attr.Value {
	t.Helper()
	appElemType := types.ObjectType{AttrTypes: applicationAttrTypes}
	appElems := make([]attr.Value, len(appNames))
	for i, n := range appNames {
		obj, diags := types.ObjectValue(applicationAttrTypes, map[string]attr.Value{
			ReleaseFieldName: types.StringValue(n),
		})
		require.False(t, diags.HasError())
		appElems[i] = obj
	}
	appsList, diags := types.ListValue(appElemType, appElems)
	require.False(t, diags.HasError())
	scopedTo, diags := types.ObjectValue(scopedToAttrTypes, map[string]attr.Value{
		ReleaseFieldApplications: appsList,
	})
	require.False(t, diags.HasError())
	return scopedTo
}

// buildServiceModel creates a ReleaseModel with one service (optionally scoped).
func buildServiceModel(t *testing.T, svcName string, scopedToAppNames []string) ReleaseModel {
	t.Helper()
	appElemType := types.ObjectType{AttrTypes: applicationAttrTypes}
	svcElemType := types.ObjectType{AttrTypes: serviceAttrTypes}

	var scopedToVal attr.Value
	if len(scopedToAppNames) > 0 {
		scopedToVal = buildScopedToValue(t, scopedToAppNames)
	} else {
		scopedToVal = types.ObjectNull(scopedToAttrTypes)
	}

	svcObj, diags := types.ObjectValue(serviceAttrTypes, map[string]attr.Value{
		ReleaseFieldName:     types.StringValue(svcName),
		ReleaseFieldScopedTo: scopedToVal,
	})
	require.False(t, diags.HasError())

	svcList, diags := types.ListValue(svcElemType, []attr.Value{svcObj})
	require.False(t, diags.HasError())

	return ReleaseModel{
		ID:           types.StringValue("id-1"),
		Name:         types.StringValue("release"),
		Start:        types.Int64Value(1742349976000),
		LastUpdated:  types.Int64Value(0),
		Applications: types.ListValueMust(appElemType, []attr.Value{}),
		Services:     svcList,
	}
}

// ─── NewReleaseResourceHandle ─────────────────────────────────────────────────

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
		metaData: resourcehandle.ResourceMetaData{ResourceName: "test_release"},
	}
	assert.Equal(t, "test_release", r.MetaData().ResourceName)
}

func TestReleaseSetComputedFields(t *testing.T) {
	r := NewReleaseResourceHandle()
	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	diags := r.SetComputedFields(context.Background(), plan)
	assert.False(t, diags.HasError())
}

func TestReleaseGetStateUpgraders(t *testing.T) {
	assert.Nil(t, NewReleaseResourceHandle().GetStateUpgraders(context.Background()))
}

// ─── UpdateState ─────────────────────────────────────────────────────────────

func TestUpdateState_BasicRelease(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	release := &api.ReleaseWithMetadata{
		ID: "Tiu16hLCTniHDtHb_uDV1w", Name: "demo-app/main-**",
		Start: 1709091782000, LastUpdated: 1709091782533,
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	setState(t, state, emptyReleaseModel())

	require.False(t, r.UpdateState(ctx, state, nil, release).HasError())

	m := getModel(t, state)
	assert.Equal(t, "Tiu16hLCTniHDtHb_uDV1w", m.ID.ValueString())
	assert.Equal(t, "demo-app/main-**", m.Name.ValueString())
	assert.Equal(t, int64(1709091782000), m.Start.ValueInt64())
	assert.Equal(t, int64(1709091782533), m.LastUpdated.ValueInt64())
	assert.Equal(t, 0, len(m.Applications.Elements()))
	assert.Equal(t, 0, len(m.Services.Elements()))
}

func TestUpdateState_WithApplications(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	release := &api.ReleaseWithMetadata{
		ID: "id-1", Name: "backend/v1", Start: 1709091782000, LastUpdated: 1709091782533,
		Applications: []*api.ReleaseApplicationScope{{Name: "app1"}, {Name: "app2"}},
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	setState(t, state, emptyReleaseModel())

	require.False(t, r.UpdateState(ctx, state, nil, release).HasError())

	m := getModel(t, state)
	assert.Equal(t, 2, len(m.Applications.Elements()))
}

func TestUpdateState_WithServiceNoScopedTo(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	release := &api.ReleaseWithMetadata{
		ID: "id-1", Name: "svc/v1", Start: 1709091782000, LastUpdated: 1709091782533,
		Services: []*api.ReleaseServiceScope{{Name: "payment"}},
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	setState(t, state, emptyReleaseModel())

	require.False(t, r.UpdateState(ctx, state, nil, release).HasError())

	m := getModel(t, state)
	assert.Equal(t, 1, len(m.Services.Elements()))
}

func TestUpdateState_WithServiceAndScopedTo(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	release := &api.ReleaseWithMetadata{
		ID: "id-2", Name: "svc/v2", Start: 1709091782000, LastUpdated: 1709091782533,
		Services: []*api.ReleaseServiceScope{
			{
				Name: "payment",
				ScopedTo: &api.ReleaseServiceScopedTo{
					Applications: []*api.ReleaseApplicationScope{{Name: "checkout-app"}},
				},
			},
		},
	}

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	setState(t, state, emptyReleaseModel())

	require.False(t, r.UpdateState(ctx, state, nil, release).HasError())

	m := getModel(t, state)
	assert.Equal(t, 1, len(m.Services.Elements()))

	// Inspect the service object directly
	var svcs []ServiceModel
	require.False(t, m.Services.ElementsAs(ctx, &svcs, false).HasError())
	require.Len(t, svcs, 1)
	assert.Equal(t, "payment", svcs[0].Name.ValueString())
	require.NotNil(t, svcs[0].ScopedTo)

	var scopedApps []ApplicationModel
	require.False(t, svcs[0].ScopedTo.Applications.ElementsAs(ctx, &scopedApps, false).HasError())
	require.Len(t, scopedApps, 1)
	assert.Equal(t, "checkout-app", scopedApps[0].Name.ValueString())
}

// ─── MapStateToDataObject ─────────────────────────────────────────────────────

func TestMapStateToDataObject_BasicFromPlan(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := emptyReleaseModel()
	model.ID = types.StringValue("release-id-1")
	model.Name = types.StringValue("frontend/release-1000")
	model.Start = types.Int64Value(1742349976000)

	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	setPlan(t, plan, model)

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, release)
	assert.Equal(t, "release-id-1", release.ID)
	assert.Equal(t, "frontend/release-1000", release.Name)
	assert.Equal(t, int64(1742349976000), release.Start)
	assert.Nil(t, release.Applications)
	assert.Nil(t, release.Services)
}

func TestMapStateToDataObject_BasicFromState(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := emptyReleaseModel()
	model.Name = types.StringValue("state-read-release")
	model.Start = types.Int64Value(1742349976000)

	state := &tfsdk.State{Schema: r.MetaData().Schema}
	setState(t, state, model)

	release, diags := r.MapStateToDataObject(ctx, nil, state)
	require.False(t, diags.HasError())
	assert.Equal(t, "state-read-release", release.Name)
}

func TestMapStateToDataObject_WithApplications(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	appElemType := types.ObjectType{AttrTypes: applicationAttrTypes}
	appObj, _ := types.ObjectValue(applicationAttrTypes, map[string]attr.Value{
		ReleaseFieldName: types.StringValue("my-app"),
	})
	appList, _ := types.ListValue(appElemType, []attr.Value{appObj})

	model := emptyReleaseModel()
	model.Name = types.StringValue("app-release")
	model.Start = types.Int64Value(1742349976000)
	model.Applications = appList

	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	setPlan(t, plan, model)

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.Len(t, release.Applications, 1)
	assert.Equal(t, "my-app", release.Applications[0].Name)
}

func TestMapStateToDataObject_WithServiceNoScopedTo(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := buildServiceModel(t, "my-service", nil)
	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	setPlan(t, plan, model)

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.Len(t, release.Services, 1)
	assert.Equal(t, "my-service", release.Services[0].Name)
	assert.Nil(t, release.Services[0].ScopedTo)
}

func TestMapStateToDataObject_WithServiceAndScopedTo(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := buildServiceModel(t, "payment", []string{"checkout-app"})
	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	setPlan(t, plan, model)

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.Len(t, release.Services, 1)
	assert.Equal(t, "payment", release.Services[0].Name)
	require.NotNil(t, release.Services[0].ScopedTo)
	require.Len(t, release.Services[0].ScopedTo.Applications, 1)
	assert.Equal(t, "checkout-app", release.Services[0].ScopedTo.Applications[0].Name)
}

func TestMapStateToDataObject_WithServiceAndMultipleScopedToApps(t *testing.T) {
	ctx := context.Background()
	r := NewReleaseResourceHandle()

	model := buildServiceModel(t, "payment", []string{"app-a", "app-b"})
	plan := &tfsdk.Plan{Schema: r.MetaData().Schema}
	setPlan(t, plan, model)

	release, diags := r.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, release.Services[0].ScopedTo)
	require.Len(t, release.Services[0].ScopedTo.Applications, 2)
	assert.Equal(t, "app-a", release.Services[0].ScopedTo.Applications[0].Name)
	assert.Equal(t, "app-b", release.Services[0].ScopedTo.Applications[1].Name)
}

// ─── helper functions ─────────────────────────────────────────────────────────

func TestBuildApplicationList_Empty(t *testing.T) {
	list, diags := buildApplicationList(nil)
	require.False(t, diags.HasError())
	assert.Equal(t, 0, len(list.Elements()))
}

func TestBuildApplicationList_OneEntry(t *testing.T) {
	list, diags := buildApplicationList([]*api.ReleaseApplicationScope{{Name: "app1"}})
	require.False(t, diags.HasError())
	assert.Equal(t, 1, len(list.Elements()))
}

func TestBuildScopedToObject_Nil(t *testing.T) {
	val, diags := buildScopedToObject(nil)
	require.False(t, diags.HasError())
	obj, ok := val.(types.Object)
	require.True(t, ok)
	assert.True(t, obj.IsNull())
}

func TestBuildScopedToObject_WithApplications(t *testing.T) {
	scopedTo := &api.ReleaseServiceScopedTo{
		Applications: []*api.ReleaseApplicationScope{{Name: "app1"}},
	}
	val, diags := buildScopedToObject(scopedTo)
	require.False(t, diags.HasError())
	obj, ok := val.(types.Object)
	require.True(t, ok)
	assert.False(t, obj.IsNull())
}
