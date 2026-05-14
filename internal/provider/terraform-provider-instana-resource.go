package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// CorrelationIDHeader is the HTTP header name for correlation ID
const CorrelationIDHeader = "X-Correlation-ID"

// attributeGetter is an interface for types that can get attributes (Plan and State)
type attributeGetter interface {
	GetAttribute(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics
}

// NewTerraformResource creates a new terraform resource for the given handle
func NewTerraformResource[T client.InstanaDataObject](handle resourcehandle.ResourceHandle[T]) TerraformResource {
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

type terraformResourceImpl[T client.InstanaDataObject] struct {
	resourceHandle resourcehandle.ResourceHandle[T]
	providerMeta   *shared.ProviderMeta
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

	providerMeta, ok := req.ProviderData.(*shared.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.providerMeta = providerMeta
}

// Create defines the create operation for the terraform resource
func (r *terraformResourceImpl[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Generate correlation ID for this operation
	correlationID := util.GenerateCorrelationID()
	
	// Add correlation ID to client config headers
	r.addCorrelationIDToClient(correlationID)
	
	// Extract resource name from plan
	resourceName := r.getResourceName(ctx, &req.Plan)
	
	tflog.Debug(ctx, "Starting resource creation", map[string]interface{}{
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})

	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		tflog.Error(ctx, "Provider not configured for resource creation")
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
		tflog.Debug(ctx, "Generated resource ID", map[string]interface{}{
			"resource_id":    id,
			"correlation_id": correlationID,
			"resource_name":  resourceName,
		})
		resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set computed fields
	diags := r.resourceHandle.SetComputedFields(ctx, &req.Plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to set computed fields")
		return
	}

	// Map state to data object
	createRequest, diags := r.resourceHandle.MapStateToDataObject(ctx, &req.Plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to map state to data object")
		return
	}

	// Create the resource
	tflog.Debug(ctx, "Calling Instana API to create resource", map[string]interface{}{
		"resource_id":    createRequest.GetIDForResourcePath(),
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})
	createdObject, err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).Create(createRequest)
	if err != nil {
		tflog.Error(ctx, "Failed to create resource via API", map[string]interface{}{
			"resource_id":    createRequest.GetIDForResourcePath(),
			"correlation_id": correlationID,
			"resource_name":  resourceName,
			"error":          err.Error(),
		})
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Could not create resource: %s", err),
		)
		return
	}

	// Update state with created object
	diags = r.resourceHandle.UpdateState(ctx, &resp.State, &req.Plan, createdObject)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to update state after creation", map[string]interface{}{
			"resource_id":    createdObject.GetIDForResourcePath(),
			"correlation_id": correlationID,
			"resource_name":  resourceName,
		})
		return
	}

	tflog.Debug(ctx, "Successfully created resource", map[string]interface{}{
		"resource_id":    createdObject.GetIDForResourcePath(),
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})
}

// Read defines the read operation for the terraform resource
func (r *terraformResourceImpl[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Generate correlation ID for this operation
	correlationID := util.GenerateCorrelationID()
	
	// Add correlation ID to client config headers
	r.addCorrelationIDToClient(correlationID)
	
	// Extract resource name from state
	resourceName := r.getResourceName(ctx, &req.State)
	
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		tflog.Error(ctx, "Provider not configured for resource read")
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
		tflog.Error(ctx, "Resource ID is missing")
		resp.Diagnostics.AddError(
			"Resource ID is missing",
			fmt.Sprintf("Resource ID of %s is missing", r.resourceHandle.MetaData().ResourceName),
		)
		return
	}

	tflog.Debug(ctx, "Reading resource from API", map[string]interface{}{
		"resource_id":    resourceID,
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})

	// Get the resource from the API
	obj, err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).GetOne(resourceID)
	if err != nil {
		if errors.Is(err, client.ErrEntityNotFound) {
			tflog.Warn(ctx, "Resource not found, removing from state", map[string]interface{}{
				"resource_id":    resourceID,
				"correlation_id": correlationID,
				"resource_name":  resourceName,
			})
			// Resource no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
		tflog.Error(ctx, "Failed to read resource from API", map[string]interface{}{
			"resource_id":    resourceID,
			"correlation_id": correlationID,
			"resource_name":  resourceName,
			"error":          err.Error(),
		})
		resp.Diagnostics.AddError(
			"Error reading resource",
			fmt.Sprintf("Could not read resource: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Successfully read resource from API", map[string]interface{}{
		"resource_id":    resourceID,
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})

	// Update state with the current object
	diags := r.resourceHandle.UpdateState(ctx, &resp.State, nil, obj)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to update state after read", map[string]interface{}{
			"resource_id":    resourceID,
			"correlation_id": correlationID,
			"resource_name":  resourceName,
		})
	}
}

// Update defines the update operation for the terraform resource
func (r *terraformResourceImpl[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Generate correlation ID for this operation
	correlationID := util.GenerateCorrelationID()
	
	// Add correlation ID to client config headers
	r.addCorrelationIDToClient(correlationID)
	
	// Extract resource name from plan
	resourceName := r.getResourceName(ctx, &req.Plan)
	
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		tflog.Error(ctx, "Provider not configured for resource update")
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
		tflog.Error(ctx, "Failed to map state to data object for update")
		return
	}

	tflog.Debug(ctx, "Starting resource update", map[string]interface{}{
		"resource_id":    obj.GetIDForResourcePath(),
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})

	// Update the resource
	tflog.Debug(ctx, "Calling Instana API to update resource", map[string]interface{}{
		"resource_id":    obj.GetIDForResourcePath(),
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})
	updatedObject, err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).Update(obj)
	if err != nil {
		tflog.Error(ctx, "Failed to update resource via API", map[string]interface{}{
			"resource_id":    obj.GetIDForResourcePath(),
			"correlation_id": correlationID,
			"resource_name":  resourceName,
			"error":          err.Error(),
		})
		resp.Diagnostics.AddError(
			"Error updating resource",
			fmt.Sprintf("Could not update resource: %s", err),
		)
		return
	}

	// Update state with updated object
	diags = r.resourceHandle.UpdateState(ctx, &resp.State, &req.Plan, updatedObject)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to update state after update", map[string]interface{}{
			"resource_id":    updatedObject.GetIDForResourcePath(),
			"correlation_id": correlationID,
			"resource_name":  resourceName,
		})
		return
	}

	tflog.Debug(ctx, "Successfully updated resource", map[string]interface{}{
		"resource_id":    updatedObject.GetIDForResourcePath(),
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})
}

// Delete defines the delete operation for the terraform resource
func (r *terraformResourceImpl[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Generate correlation ID for this operation
	correlationID := util.GenerateCorrelationID()
	
	// Add correlation ID to client config headers
	r.addCorrelationIDToClient(correlationID)
	
	// Extract resource name from state
	resourceName := r.getResourceName(ctx, &req.State)
	
	if r.providerMeta == nil || r.providerMeta.InstanaAPI == nil {
		tflog.Error(ctx, "Provider not configured for resource deletion")
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
		tflog.Error(ctx, "Failed to map state to data object for deletion")
		return
	}

	resourceID := object.GetIDForResourcePath()
	tflog.Debug(ctx, "Starting resource deletion", map[string]interface{}{
		"resource_id":    resourceID,
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})

	// Delete the resource
	tflog.Debug(ctx, "Calling Instana API to delete resource", map[string]interface{}{
		"resource_id":    resourceID,
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})
	err := r.resourceHandle.GetRestResource(r.providerMeta.InstanaAPI).DeleteByID(resourceID)
	if err != nil {
		tflog.Error(ctx, "Failed to delete resource via API", map[string]interface{}{
			"resource_id":    resourceID,
			"correlation_id": correlationID,
			"resource_name":  resourceName,
			"error":          err.Error(),
		})
		resp.Diagnostics.AddError(
			"Error deleting resource",
			fmt.Sprintf("Could not delete resource: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Successfully deleted resource", map[string]interface{}{
		"resource_id":    resourceID,
		"correlation_id": correlationID,
		"resource_name":  resourceName,
	})

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// ImportState handles importing an existing resource into Terraform
func (r *terraformResourceImpl[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceName := r.resourceHandle.MetaData().ResourceName
	tflog.Info(ctx, "Importing resource", map[string]interface{}{
		"resource_type": resourceName,
		"resource_id":   req.ID,
	})

	// Set ID
	if r.resourceHandle.MetaData().ResourceIDField != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(*r.resourceHandle.MetaData().ResourceIDField), types.StringValue(req.ID))...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	}

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to import resource", map[string]interface{}{
			"resource_type": resourceName,
			"resource_id":   req.ID,
		})
	} else {
		tflog.Info(ctx, "Successfully imported resource", map[string]interface{}{
			"resource_type": resourceName,
			"resource_id":   req.ID,
		})
	}
}

// UpgradeState handles state upgrades for resources that need schema migration
func (r *terraformResourceImpl[T]) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return r.resourceHandle.GetStateUpgraders(ctx)
}

// addCorrelationIDToClient adds the correlation ID header to the client configuration
func (r *terraformResourceImpl[T]) addCorrelationIDToClient(correlationID string) {
	if r.providerMeta != nil && r.providerMeta.ClientConfig != nil {
		if r.providerMeta.ClientConfig.Headers.Custom == nil {
			r.providerMeta.ClientConfig.Headers.Custom = make(map[string]string)
		}
		r.providerMeta.ClientConfig.Headers.Custom[CorrelationIDHeader] = correlationID
	}
}

// getResourceName extracts the resource name from the plan or state
func (r *terraformResourceImpl[T]) getResourceName(ctx context.Context, getter attributeGetter) string {
	if getter == nil {
		return ""
	}
	
	var name types.String
	diags := getter.GetAttribute(ctx, path.Root("name"), &name)
	if diags.HasError() || name.IsNull() || name.IsUnknown() {
		return ""
	}
	
	return name.ValueString()
}
