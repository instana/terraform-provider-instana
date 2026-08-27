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

// applicationAttrTypes defines the attribute types for ApplicationModel elements in a types.List.
var applicationAttrTypes = map[string]attr.Type{
	ReleaseFieldName: types.StringType,
}

// scopedToAttrTypes defines the attribute types for ScopedToModel elements in a types.Object.
var scopedToAttrTypes = map[string]attr.Type{
	ReleaseFieldApplicationName: types.StringType,
	ReleaseFieldEnvironmentName: types.StringType,
}

// serviceAttrTypes defines the attribute types for ServiceModel elements in a types.List.
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
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, ReleaseNameMaxLength),
				},
			},
			ReleaseFieldStart: schema.Int64Attribute{
				Required:    true,
				Description: ReleaseDescStart,
				Validators: []validator.Int64{
					int64validator.AtLeast(ReleaseStartMinValue),
				},
			},
			ReleaseFieldLastUpdated: schema.Int64Attribute{
				Computed:    true,
				Description: ReleaseDescLastUpdated,
			},
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
						ReleaseFieldName: schema.StringAttribute{
							Required:    true,
							Description: ReleaseDescApplicationName,
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, ReleaseNameMaxLength),
							},
						},
					},
				},
			},
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
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, ReleaseNameMaxLength),
							},
						},
						ReleaseFieldScopedTo: schema.SingleNestedAttribute{
							Optional:    true,
							Description: ReleaseDescScopedTo,
							Attributes: map[string]schema.Attribute{
								ReleaseFieldApplicationName: schema.StringAttribute{
									Optional:    true,
									Description: ReleaseDescApplicationNameScope,
								},
								ReleaseFieldEnvironmentName: schema.StringAttribute{
									Optional:    true,
									Description: ReleaseDescEnvironmentName,
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

func (r *releaseResource) UpdateState(ctx context.Context, state *tfsdk.State, _ *tfsdk.Plan, release *api.ReleaseWithMetadata) diag.Diagnostics {
	var diags diag.Diagnostics

	applications, d := mapApplicationsToState(ctx, release.Applications)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	services, d := mapServicesToState(ctx, release.Services)
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

func mapApplicationsToState(ctx context.Context, apps []*api.ReleaseApplicationScope) (types.List, diag.Diagnostics) {
	elementType := types.ObjectType{AttrTypes: applicationAttrTypes}

	if len(apps) == 0 {
		return types.ListValueMust(elementType, []attr.Value{}), nil
	}

	elems := make([]attr.Value, len(apps))
	for i, app := range apps {
		obj, diags := types.ObjectValue(applicationAttrTypes, map[string]attr.Value{
			ReleaseFieldName: types.StringValue(app.Name),
		})
		if diags.HasError() {
			return types.ListNull(elementType), diags
		}
		elems[i] = obj
	}

	return types.ListValue(elementType, elems)
}

func mapServicesToState(ctx context.Context, services []*api.ReleaseServiceScope) (types.List, diag.Diagnostics) {
	elementType := types.ObjectType{AttrTypes: serviceAttrTypes}

	if len(services) == 0 {
		return types.ListValueMust(elementType, []attr.Value{}), nil
	}

	elems := make([]attr.Value, len(services))
	for i, svc := range services {
		scopedToVal, diags := buildScopedToObject(svc.ScopedTo)
		if diags.HasError() {
			return types.ListNull(elementType), diags
		}

		obj, diags := types.ObjectValue(serviceAttrTypes, map[string]attr.Value{
			ReleaseFieldName:     types.StringValue(svc.Name),
			ReleaseFieldScopedTo: scopedToVal,
		})
		if diags.HasError() {
			return types.ListNull(elementType), diags
		}
		elems[i] = obj
	}

	return types.ListValue(elementType, elems)
}

func buildScopedToObject(scopedTo *api.ReleaseServiceScopedTo) (attr.Value, diag.Diagnostics) {
	if scopedTo == nil {
		return types.ObjectNull(scopedToAttrTypes), nil
	}
	return types.ObjectValue(scopedToAttrTypes, map[string]attr.Value{
		ReleaseFieldApplicationName: types.StringValue(scopedTo.ApplicationName),
		ReleaseFieldEnvironmentName: types.StringValue(scopedTo.EnvironmentName),
	})
}

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

	apps, d := mapApplicationsFromState(ctx, model.Applications)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	release.Applications = apps

	svcs, d := mapServicesFromState(ctx, model.Services)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	release.Services = svcs

	return release, diags
}

func mapApplicationsFromState(ctx context.Context, list types.List) ([]*api.ReleaseApplicationScope, diag.Diagnostics) {
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

func mapServicesFromState(ctx context.Context, list types.List) ([]*api.ReleaseServiceScope, diag.Diagnostics) {
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
			scope.ScopedTo = &api.ReleaseServiceScopedTo{
				ApplicationName: svc.ScopedTo.ApplicationName.ValueString(),
				EnvironmentName: svc.ScopedTo.EnvironmentName.ValueString(),
			}
		}
		result[i] = scope
	}
	return result, nil
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *releaseResource) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return nil
}
