package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ResourceInstanaSloConfigFramework the name of the terraform-provider-instana resource to manage SLO configurations
const ResourceInstanaSloConfigFramework = "slo_config"

// NewSloConfigResourceHandleFramework creates the resource handle for SLO Config
func NewSloConfigResourceHandleFramework() ResourceHandleFramework[*restapi.SloConfig] {
	return &sloConfigResourceFramework{}
}

type sloConfigResourceFramework struct{}

func (r *sloConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &ResourceMetaDataFramework{
		ResourceName: ResourceInstanaSloConfigFramework,
		// Use a simple schema for now, the actual implementation will use the tf_framework resource
		Schema: schema.Schema{
			Description: "This resource manages SLO Configurations in Instana.",
		},
		SchemaVersion: 1,
	}
}

func (r *sloConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloConfig] {
	return api.SloConfigs()
}

func (r *sloConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sloConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloConfig, diag.Diagnostics) {
	// This is a wrapper for the existing SloConfigResource
	// In a real implementation, we would need to properly map the state to the API object
	// For now, we'll return a placeholder
	return &restapi.SloConfig{}, nil
}

func (r *sloConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SloConfig) diag.Diagnostics {
	// This is a wrapper for the existing SloConfigResource
	// In a real implementation, we would need to properly update the state from the API object
	return nil
}

// Made with Bob
