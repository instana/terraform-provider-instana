package release

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	instanaClient "github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// applicationAttrTypes are the attribute types for a single {name} application entry used
// both at the top-level applications list and inside scoped_to.applications.
var applicationAttrTypes = map[string]attr.Type{
	ReleaseFieldName: types.StringType,
}

// scopedToAttrTypes are the attribute types for the scoped_to object.
var scopedToAttrTypes = map[string]attr.Type{
	ReleaseFieldApplications: types.ListType{ElemType: types.ObjectType{AttrTypes: applicationAttrTypes}},
}

// serviceAttrTypes are the attribute types for a single service list element.
var serviceAttrTypes = map[string]attr.Type{
	ReleaseFieldName:     types.StringType,
	ReleaseFieldScopedTo: types.ObjectType{AttrTypes: scopedToAttrTypes},
}

// NewReleaseResourceHandle creates the resource handle for releases
func NewReleaseResourceHandle() resourcehandle.ResourceHandle[*api.ReleaseWithMetadata] {
	return &releaseResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaRelease,
			Schema:        buildReleaseSchema(),
			SchemaVersion: 0,
		},
	}
}

func buildReleaseSchema() schema.Schema {
	// Reusable nested attribute for {name} entries
	nameAttr := schema.StringAttribute{
		Required:    true,
		Description: ReleaseDescApplicationName,
		Validators:  []validator.String{stringvalidator.LengthBetween(0, ReleaseNameMaxLength)},
	}

	return schema.Schema{
		Description: ReleaseDescResource,
		Attributes: map[string]schema.Attribute{
			ReleaseFieldID: schema.StringAttribute{
				Computed:    true,
				Description: ReleaseDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			ReleaseFieldName: schema.StringAttribute{
				Required:    true,
				Description: ReleaseDescName,
				Validators:  []validator.String{stringvalidator.LengthBetween(0, ReleaseNameMaxLength)},
			},
			ReleaseFieldStart: schema.Int64Attribute{
				Required:    true,
				Description: ReleaseDescStart,
				Validators:  []validator.Int64{int64validator.AtLeast(ReleaseStartMinValue)},
			},
			ReleaseFieldLastUpdated: schema.Int64Attribute{
				Computed:    true,
				Description: ReleaseDescLastUpdated,
			},
			// Top-level applications: list of {name}
			ReleaseFieldApplications: schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: ReleaseDescApplications,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(0, ReleaseApplicationsMaxItems),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ReleaseFieldName: nameAttr,
					},
				},
			},
			// Services: list of {name, scoped_to?}
			ReleaseFieldServices: schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: ReleaseDescServices,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(0, ReleaseServicesMaxItems),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						ReleaseFieldName: schema.StringAttribute{
							Required:    true,
							Description: ReleaseDescServiceName,
							Validators:  []validator.String{stringvalidator.LengthBetween(0, ReleaseNameMaxLength)},
						},
						// scoped_to: optional, contains applications list (1–10)
						ReleaseFieldScopedTo: schema.SingleNestedAttribute{
							Optional:    true,
							Description: ReleaseDescScopedTo,
							Attributes: map[string]schema.Attribute{
								ReleaseFieldApplications: schema.ListNestedAttribute{
									Required:    true,
									Description: ReleaseDescScopedToApplications,
									Validators: []validator.List{
										listvalidator.SizeBetween(ReleaseScopedToApplicationsMinItems, ReleaseScopedToApplicationsMaxItems),
									},
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											ReleaseFieldName: nameAttr,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *releaseResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *releaseResource) GetRestResource(instana instanaClient.InstanaAPI) rest.RestResource[*api.ReleaseWithMetadata] {
	return instana.Releases()
}

func (r *releaseResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ─── State → Terraform ───────────────────────────────────────────────────────

func (r *releaseResource) UpdateState(ctx context.Context, state *tfsdk.State, _ *tfsdk.Plan, release *api.ReleaseWithMetadata) diag.Diagnostics {
	var diags diag.Diagnostics

	applications, d := buildApplicationList(release.Applications)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	services, d := buildServiceList(ctx, release.Services)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	model := ReleaseModel{
		ID:           types.StringValue(release.ID),
		Name:         types.StringValue(release.Name),
		Start:        types.Int64Value(release.Start),
		LastUpdated:  types.Int64Value(release.LastUpdated),
		Applications: applications,
		Services:     services,
	}

	diags.Append(state.Set(ctx, model)...)
	return diags
}

// buildApplicationList converts []*ReleaseApplicationScope into a types.List of {name} objects.
func buildApplicationList(apps []*api.ReleaseApplicationScope) (types.List, diag.Diagnostics) {
	elemType := types.ObjectType{AttrTypes: applicationAttrTypes}
	if len(apps) == 0 {
		return types.ListValueMust(elemType, []attr.Value{}), nil
	}
	elems := make([]attr.Value, len(apps))
	for i, app := range apps {
		obj, diags := types.ObjectValue(applicationAttrTypes, map[string]attr.Value{
			ReleaseFieldName: types.StringValue(app.Name),
		})
		if diags.HasError() {
			return types.ListNull(elemType), diags
		}
		elems[i] = obj
	}
	return types.ListValue(elemType, elems)
}

// buildServiceList converts []*ReleaseServiceScope into a types.List of service objects.
func buildServiceList(ctx context.Context, services []*api.ReleaseServiceScope) (types.List, diag.Diagnostics) {
	elemType := types.ObjectType{AttrTypes: serviceAttrTypes}
	if len(services) == 0 {
		return types.ListValueMust(elemType, []attr.Value{}), nil
	}
	elems := make([]attr.Value, len(services))
	for i, svc := range services {
		scopedToVal, diags := buildScopedToObject(svc.ScopedTo)
		if diags.HasError() {
			return types.ListNull(elemType), diags
		}
		obj, diags := types.ObjectValue(serviceAttrTypes, map[string]attr.Value{
			ReleaseFieldName:     types.StringValue(svc.Name),
			ReleaseFieldScopedTo: scopedToVal,
		})
		if diags.HasError() {
			return types.ListNull(elemType), diags
		}
		elems[i] = obj
	}
	return types.ListValue(elemType, elems)
}

// buildScopedToObject converts *ReleaseServiceScopedTo into an attr.Value (object or null).
func buildScopedToObject(scopedTo *api.ReleaseServiceScopedTo) (attr.Value, diag.Diagnostics) {
	if scopedTo == nil {
		return types.ObjectNull(scopedToAttrTypes), nil
	}
	appsList, diags := buildApplicationList(scopedTo.Applications)
	if diags.HasError() {
		return types.ObjectNull(scopedToAttrTypes), diags
	}
	return types.ObjectValue(scopedToAttrTypes, map[string]attr.Value{
		ReleaseFieldApplications: appsList,
	})
}

// ─── Terraform → API ─────────────────────────────────────────────────────────

func (r *releaseResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.ReleaseWithMetadata, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model ReleaseModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	if diags.HasError() {
		return nil, diags
	}

	release := &api.ReleaseWithMetadata{
		ID:    model.ID.ValueString(),
		Name:  model.Name.ValueString(),
		Start: model.Start.ValueInt64(),
	}

	apps, d := extractApplicationsFromList(ctx, model.Applications)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	release.Applications = apps

	svcs, d := extractServicesFromList(ctx, model.Services)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	release.Services = svcs

	return release, diags
}

// extractApplicationsFromList converts a types.List of {name} objects to []*ReleaseApplicationScope.
func extractApplicationsFromList(ctx context.Context, list types.List) ([]*api.ReleaseApplicationScope, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() || len(list.Elements()) == 0 {
		return nil, nil
	}
	var models []ApplicationModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]*api.ReleaseApplicationScope, len(models))
	for i, m := range models {
		result[i] = &api.ReleaseApplicationScope{Name: m.Name.ValueString()}
	}
	return result, nil
}

// extractServicesFromList converts a types.List of service objects to []*ReleaseServiceScope.
func extractServicesFromList(ctx context.Context, list types.List) ([]*api.ReleaseServiceScope, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() || len(list.Elements()) == 0 {
		return nil, nil
	}
	var models []ServiceModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]*api.ReleaseServiceScope, len(models))
	for i, svc := range models {
		scope := &api.ReleaseServiceScope{Name: svc.Name.ValueString()}
		if svc.ScopedTo != nil {
			apps, d := extractApplicationsFromList(ctx, svc.ScopedTo.Applications)
			if d.HasError() {
				return nil, d
			}
			scope.ScopedTo = &api.ReleaseServiceScopedTo{Applications: apps}
		}
		result[i] = scope
	}
	return result, nil
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *releaseResource) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return nil
}
