package resourcehandle

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
)

// SingletonResourceHandle is the parallel of ResourceHandle for singleton REST resources
// (resources that have no ID and expose Get / Upsert / Delete at a fixed path).
// T is the API data type; it does NOT need to satisfy client.InstanaDataObject because
// singleton resources carry no ID.
//
// Implementations follow exactly the same shape as ResourceHandle: MetaData, UpdateState,
// MapStateToDataObject, SetComputedFields, and GetStateUpgraders are identical.
// The only difference is GetSingletonRestResource, which returns a SingletonRestResource
// instead of a RestResource.
type SingletonResourceHandle[T any] interface {
	// MetaData returns the metadata of this handle (resource name, schema, etc.)
	MetaData() *ResourceMetaData

	// GetSingletonRestResource provides the singleton REST resource client.
	GetSingletonRestResource(api client.InstanaAPI) rest.SingletonRestResource[T]

	// UpdateState updates the Terraform state with data received from the API.
	UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, obj T) diag.Diagnostics

	// MapStateToDataObject maps the current Terraform plan/state to the API data object.
	MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (T, diag.Diagnostics)

	// SetComputedFields calculates and sets computed field values in the plan.
	SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics

	// GetStateUpgraders returns the state upgraders for schema version migrations.
	// Return nil or an empty map if no state upgrades are needed.
	GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader
}

// ResourceMetaData the metadata of a terraform ResourceHandle
type ResourceMetaData struct {
	ResourceName       string
	Schema             schema.Schema
	SchemaVersion      int64
	SkipIDGeneration   bool
	ResourceIDField    *string
	CreateOnly         bool
	DeprecationMessage string
}

// ResourceHandle resource specific implementation which provides metadata and maps data from/to terraform state.
// Together with TerraformResource terraform schema resources can be created
type ResourceHandle[T client.InstanaDataObject] interface {
	// MetaData returns the metadata of this ResourceHandle
	MetaData() *ResourceMetaData

	// GetRestResource provides the client.RestResource used by the ResourceHandle
	GetRestResource(api client.InstanaAPI) rest.RestResource[T]

	// UpdateState updates the state of the resource with the input data from the Instana API
	UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, obj T) diag.Diagnostics

	// MapStateToDataObject maps the current state to the API model of the Instana API
	MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (T, diag.Diagnostics)

	// SetComputedFields calculate and set the calculated value of computed fields of the given resource
	SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics

	// GetStateUpgraders returns the state upgraders for migrating resource state between schema versions
	// Return nil or empty map if no state upgrades are needed
	GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader
}

// PostCreateUpdater is an optional interface that a ResourceHandle can implement
// to indicate that an additional Update API call should be made immediately after
// the Create API call. This is required for resources (e.g. custom dashboards)
// whose Create endpoint does not accept certain fields (e.g. RBAC tags) but whose
// Update endpoint does.
//
// If the resource handle implements this interface and NeedsPostCreateUpdate
// returns true for the just-created object, the generic Create operation will
// call the Update API with the original request payload (with the ID from the
// created object applied) and use the update response to set the final state.
type PostCreateUpdater[T client.InstanaDataObject] interface {
	// NeedsPostCreateUpdate returns true when the created object requires a
	// follow-up Update call to persist fields that were silently dropped by the
	// Create endpoint.
	NeedsPostCreateUpdate(original T) bool

	// ApplyCreatedID copies the server-assigned ID from the created object into
	// the original request payload so the subsequent Update call targets the
	// correct resource.
	ApplyCreatedID(original T, created T) T
}
