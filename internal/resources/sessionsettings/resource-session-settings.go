package sessionsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	restapi "github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// NewSessionSettingsResourceHandle creates the resource handle for session settings.
func NewSessionSettingsResourceHandle() resourcehandle.SingletonResourceHandle[*restapi.SessionSettings] {
	return &sessionSettingsResourceHandle{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaSessionSettings,
			Schema: schema.Schema{
				Description: SessionSettingsDescResource,
				Attributes: map[string]schema.Attribute{
					SessionSettingsFieldTokenLifeTimeInMillis: schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: SessionSettingsDescTokenLifeTime,
						Default:     int64default.StaticInt64(SessionSettingsDefaultTokenLifeTimeInMillis),
						Validators: []validator.Int64{
							int64validator.Between(
								SessionSettingsMinTokenLifeTimeInMillis,
								SessionSettingsMaxTokenLifeTimeInMillis,
							),
						},
					},
					SessionSettingsFieldIdleTimeInMillis: schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: SessionSettingsDescIdleTime,
						Default:     int64default.StaticInt64(SessionSettingsDefaultIdleTimeInMillis),
						Validators: []validator.Int64{
							int64validator.Between(
								SessionSettingsMinIdleTimeInMillis,
								SessionSettingsMaxIdleTimeInMillis,
							),
						},
					},
				},
			},
		},
	}
}

type sessionSettingsResourceHandle struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata.
func (h *sessionSettingsResourceHandle) MetaData() *resourcehandle.ResourceMetaData {
	return &h.metaData
}

// GetSingletonRestResource returns the singleton REST client for session settings.
func (h *sessionSettingsResourceHandle) GetSingletonRestResource(api client.InstanaAPI) rest.SingletonRestResource[*restapi.SessionSettings] {
	return api.SessionSettings()
}

// SetComputedFields is a no-op — session settings has no computed fields.
func (h *sessionSettingsResourceHandle) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return diag.Diagnostics{}
}

// MapStateToDataObject maps the Terraform plan/state to the API object.
func (h *sessionSettingsResourceHandle) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SessionSettings, diag.Diagnostics) {
	var model SessionSettingsModel
	var diags diag.Diagnostics

	if plan != nil {
		diags = plan.Get(ctx, &model)
	} else {
		diags = state.Get(ctx, &model)
	}

	if diags.HasError() {
		return nil, diags
	}

	return &restapi.SessionSettings{
		TokenLifeTimeInMillis: model.TokenLifeTimeInMillis.ValueInt64(),
		IdleTimeInMillis:      model.IdleTimeInMillis.ValueInt64(),
	}, diags
}

// UpdateState updates the Terraform state with the API object.
func (h *sessionSettingsResourceHandle) UpdateState(ctx context.Context, state *tfsdk.State, _ *tfsdk.Plan, settings *restapi.SessionSettings) diag.Diagnostics {
	return state.Set(ctx, SessionSettingsModel{
		TokenLifeTimeInMillis: types.Int64Value(settings.TokenLifeTimeInMillis),
		IdleTimeInMillis:      types.Int64Value(settings.IdleTimeInMillis),
	})
}

// GetStateUpgraders returns state upgraders — none needed for this resource.
func (h *sessionSettingsResourceHandle) GetStateUpgraders(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{}
}
