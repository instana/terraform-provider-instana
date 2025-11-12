package resourcehandle

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ResourceMetaDataFramework the metadata of a terraform ResourceHandleFramework
type ResourceMetaDataFramework struct {
	ResourceName       string
	Schema             schema.Schema
	SchemaVersion      int64
	SkipIDGeneration   bool
	ResourceIDField    *string
	CreateOnly         bool
	DeprecationMessage string
}

// ResourceHandleFramework resource specific implementation which provides metadata and maps data from/to terraform state.
// Together with TerraformResourceFramework terraform schema resources can be created
type ResourceHandleFramework[T restapi.InstanaDataObject] interface {
	// MetaData returns the metadata of this ResourceHandleFramework
	MetaData() *ResourceMetaDataFramework

	// GetRestResource provides the restapi.RestResource used by the ResourceHandleFramework
	GetRestResource(api restapi.InstanaAPI) restapi.RestResource[T]

	// UpdateState updates the state of the resource with the input data from the Instana API
	UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, obj T) diag.Diagnostics

	// MapStateToDataObject maps the current state to the API model of the Instana API
	MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (T, diag.Diagnostics)

	// SetComputedFields calculate and set the calculated value of computed fields of the given resource
	SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics
}
