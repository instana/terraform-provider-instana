package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewTerraformResource creates a new terraform resource for the given handle
func NewTerraformResource[T restapi.InstanaDataObject](handle resourcehandle.ResourceHandle[T]) TerraformResource {
	return &terraformResourceImpl[T]{
		resourceHandle: handle,
	}
}

// TerraformResource internal simplified representation of a Terraform resource
type TerraformResource interface {
	resource.Resource
	resource.ResourceWithConfigure
	resource.ResourceWithImportState
	resource.ResourceWithUpgradeState
}

type terraformResourceImpl[T restapi.InstanaDataObject] struct {
	resourceHandle resourcehandle.ResourceHandle[T]
	providerMeta   *restapi.ProviderMeta
}

// Metadata returns the resource type name
func (r *terraformResourceImpl[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.resourceHandle.MetaData().ResourceName
}

// Schema defines the schema for the resource
func (r *terraformResourceImpl[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.resourceHandle.MetaData().Schema
	resp.Schema.Version = r.resourceHandle.MetaData().SchemaVersion
	resp.Schema.DeprecationMessage = r.resourceHandle.MetaData().DeprecationMessage
}

// Configure stores the provider meta for use by the resource
func (r *terraformResourceImpl[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*restapi.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *restapi.ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.providerMeta = providerMeta
}

func (r *terraformResourceImpl[T]) getResourceID(d *schema.ResourceData) string {
	if r.resourceHandle.MetaData().ResourceIDField != nil {
		return d.Get(*r.resourceHandle.MetaData().ResourceIDField).(string)
	}
	return d.Id()
}

// Create defines the create operation for the terraform resource
func (r *terraformResourceImpl[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource.",
		)
		return
	}

	// Generate ID if needed
	if !r.resourceHandle.MetaData().SkipIDGeneration {
		// Set ID in state
		id := util.RandomID()
		resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set computed fields
	diags := r.resourceHandle.SetComputedFields(ctx, &req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map state to data object
	createRequest, diags := r.resourceHandle.MapStateToDataObject(ctx, &req.Plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	createdObject, err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).Create(createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Could not create resource: %s", err),
		)
		return
	}

	// Update state with created object
	diags = r.resourceHandle.UpdateState(ctx, &resp.State, &req.Plan, createdObject)
	resp.Diagnostics.Append(diags...)
}

// Read defines the read operation for the terraform resource
func (r *terraformResourceImpl[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource.",
		)
		return
	}

	// Get resource ID
	var resourceID string
	if r.resourceHandle.MetaData().ResourceIDField != nil {
		var idValue types.String
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(*r.resourceHandle.MetaData().ResourceIDField), &idValue)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resourceID = idValue.ValueString()
	} else {
		var idValue types.String
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idValue)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resourceID = idValue.ValueString()
	}

	if resourceID == "" {
		resp.Diagnostics.AddError(
			"Resource ID is missing",
			fmt.Sprintf("Resource ID of %s is missing", r.resourceHandle.MetaData().ResourceName),
		)
		return
	}

	// Get the resource from the API
	obj, err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).GetOne(resourceID)
	if err != nil {
		if errors.Is(err, restapi.ErrEntityNotFound) {
			// Resource no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading resource",
			fmt.Sprintf("Could not read resource: %s", err),
		)
		return
	}

	// Update state with the current object
	diags := r.resourceHandle.UpdateState(ctx, &resp.State, nil, obj)
	resp.Diagnostics.Append(diags...)
}

// Update defines the update operation for the terraform resource
func (r *terraformResourceImpl[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource.",
		)
		return
	}

	// Map state to data object
	obj, diags := r.resourceHandle.MapStateToDataObject(ctx, &req.Plan, &req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	updatedObject, err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).Update(obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating resource",
			fmt.Sprintf("Could not update resource: %s", err),
		)
		return
	}

	// Update state with updated object
	diags = r.resourceHandle.UpdateState(ctx, &resp.State, &req.Plan, updatedObject)
	resp.Diagnostics.Append(diags...)
}

// Delete defines the delete operation for the terraform resource
func (r *terraformResourceImpl[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource.",
		)
		return
	}

	// Map state to data object
	object, diags := r.resourceHandle.MapStateToDataObject(ctx, nil, &req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the resource
	err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).DeleteByID(object.GetIDForResourcePath())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting resource",
			fmt.Sprintf("Could not delete resource: %s", err),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource into Terraform
func (r *terraformResourceImpl[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set ID
	if r.resourceHandle.MetaData().ResourceIDField != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(*r.resourceHandle.MetaData().ResourceIDField), types.StringValue(req.ID))...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	}
}

// UpgradeState handles state upgrades for resources that need schema migration
func (r *terraformResourceImpl[T]) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return r.resourceHandle.GetStateUpgraders(ctx)
}
