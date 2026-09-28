package ipfiltering

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	restapi "github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// NewIPFilteringResourceHandle creates the singleton resource handle for IP filtering.
func NewIPFilteringResourceHandle() resourcehandle.SingletonResourceHandle[*restapi.IPFiltering] {
	return &ipFilteringResourceHandle{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaIPFiltering,
			Schema: schema.Schema{
				Description: IPFilteringDescResource,
				Attributes: map[string]schema.Attribute{
					IPFilteringFieldEnabled: schema.BoolAttribute{
						Required:    true,
						Description: IPFilteringDescEnabled,
					},
					IPFilteringFieldDenyAll: schema.BoolAttribute{
						Required:    true,
						Description: IPFilteringDescDenyAll,
					},
					IPFilteringFieldSupportAccessEnabled: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: IPFilteringDescSupportAccessEnabled,
					},
					IPFilteringFieldRules: schema.ListNestedAttribute{
						Required:    true,
						Description: IPFilteringDescRules,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								IPFilteringFieldRuleTarget: schema.StringAttribute{
									Required:    true,
									Description: IPFilteringDescRuleTarget,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								IPFilteringFieldRuleBlock: schema.BoolAttribute{
									Required:    true,
									Description: IPFilteringDescRuleBlock,
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeBetween(IPFilteringMinRules, IPFilteringMaxRules),
						},
					},
					IPFilteringFieldActive: schema.BoolAttribute{
						Computed:    true,
						Description: IPFilteringDescActive,
					},
					IPFilteringFieldLastChangedAt: schema.Int64Attribute{
						Computed:    true,
						Description: IPFilteringDescLastChangedAt,
					},
					IPFilteringFieldLastVerifiedAt: schema.Int64Attribute{
						Computed:    true,
						Description: IPFilteringDescLastVerifiedAt,
					},
				},
			},
		},
	}
}

type ipFilteringResourceHandle struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata.
func (h *ipFilteringResourceHandle) MetaData() *resourcehandle.ResourceMetaData {
	return &h.metaData
}

// GetSingletonRestResource returns the singleton REST client for IP filtering.
func (h *ipFilteringResourceHandle) GetSingletonRestResource(api client.InstanaAPI) rest.SingletonRestResource[*restapi.IPFiltering] {
	return api.IPFiltering()
}

// NeedsPostUpsertVerification returns true when the IP filtering configuration has enabled = true
// and needs to be verified to become permanently active.
func (h *ipFilteringResourceHandle) NeedsPostUpsertVerification(obj *restapi.IPFiltering) bool {
	return obj != nil && obj.Enabled
}

// Verify promotes the IP filtering configuration from temporary to permanently active.
func (h *ipFilteringResourceHandle) Verify(api client.InstanaAPI) (*restapi.IPFiltering, error) {
	return api.IPFiltering().Verify()
}

// SetComputedFields is a no-op — computed fields are populated directly from the API response in UpdateState.
func (h *ipFilteringResourceHandle) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return diag.Diagnostics{}
}

// MapStateToDataObject maps the Terraform plan/state to the API object.
func (h *ipFilteringResourceHandle) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.IPFiltering, diag.Diagnostics) {
	var model IPFilteringModel
	var diags diag.Diagnostics

	if plan != nil {
		diags = plan.Get(ctx, &model)
	} else {
		diags = state.Get(ctx, &model)
	}

	if diags.HasError() {
		return nil, diags
	}

	rules := make([]restapi.IPFilteringRule, len(model.Rules))
	for i, r := range model.Rules {
		rules[i] = restapi.IPFilteringRule{
			Target: r.Target.ValueString(),
			Block:  r.Block.ValueBool(),
		}
	}

	return &restapi.IPFiltering{
		Enabled:              model.Enabled.ValueBool(),
		DenyAll:              model.DenyAll.ValueBool(),
		SupportAccessEnabled: model.SupportAccessEnabled.ValueBool(),
		Rules:                rules,
	}, diags
}

// UpdateState updates the Terraform state with the API object.
func (h *ipFilteringResourceHandle) UpdateState(ctx context.Context, state *tfsdk.State, _ *tfsdk.Plan, filtering *restapi.IPFiltering) diag.Diagnostics {
	rules := make([]IPFilteringRuleModel, len(filtering.Rules))
	for i, r := range filtering.Rules {
		rules[i] = IPFilteringRuleModel{
			Target: types.StringValue(r.Target),
			Block:  types.BoolValue(r.Block),
		}
	}

	var lastChangedAt types.Int64
	if filtering.LastChangedAt != nil {
		lastChangedAt = types.Int64Value(*filtering.LastChangedAt)
	} else {
		lastChangedAt = types.Int64Null()
	}

	var lastVerifiedAt types.Int64
	if filtering.LastVerifiedAt != nil {
		lastVerifiedAt = types.Int64Value(*filtering.LastVerifiedAt)
	} else {
		lastVerifiedAt = types.Int64Null()
	}

	return state.Set(ctx, IPFilteringModel{
		Enabled:              types.BoolValue(filtering.Enabled),
		DenyAll:              types.BoolValue(filtering.DenyAll),
		SupportAccessEnabled: types.BoolValue(filtering.SupportAccessEnabled),
		Rules:                rules,
		Active:               types.BoolValue(filtering.Active),
		LastChangedAt:        lastChangedAt,
		LastVerifiedAt:       lastVerifiedAt,
	})
}

// GetStateUpgraders returns nil — no state schema migrations are needed for this resource.
func (h *ipFilteringResourceHandle) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return nil
}
