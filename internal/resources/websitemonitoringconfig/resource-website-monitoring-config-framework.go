package websitemonitoringconfig

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tf_framework"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewWebsiteMonitoringConfigResourceHandleFramework creates the resource handle for Website Monitoring Configurations
func NewWebsiteMonitoringConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.WebsiteMonitoringConfig] {
	return &websiteMonitoringConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaWebsiteMonitoringConfigFramework,
			Schema: schema.Schema{
				Description: WebsiteMonitoringConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: WebsiteMonitoringConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: WebsiteMonitoringConfigDescName,
					},
					"app_name": schema.StringAttribute{
						Computed:    true,
						Description: WebsiteMonitoringConfigDescAppName,
					},
				},
			},
			SchemaVersion: 1,
		},
	}
}

type websiteMonitoringConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *websiteMonitoringConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *websiteMonitoringConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.WebsiteMonitoringConfig] {
	return api.WebsiteMonitoringConfig()
}

func (r *websiteMonitoringConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *websiteMonitoringConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.WebsiteMonitoringConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model tf_framework.WebsiteMonitoringConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Create API object
	return &restapi.WebsiteMonitoringConfig{
		ID:   model.ID.ValueString(),
		Name: model.Name.ValueString(),
	}, diags
}

func (r *websiteMonitoringConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.WebsiteMonitoringConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API object
	model := tf_framework.WebsiteMonitoringConfigModel{
		ID:      types.StringValue(apiObject.ID),
		Name:    types.StringValue(apiObject.Name),
		AppName: types.StringValue(apiObject.AppName),
	}

	// Set state
	diags.Append(state.Set(ctx, &model)...)
	return diags
}
