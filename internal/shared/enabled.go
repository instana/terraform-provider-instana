package shared

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
)

const (
	// EnabledFieldName is the name of the schema attribute used by resources which support enabling and disabling
	EnabledFieldName = "enabled"
	// EnabledDefaultValue is the default value of the enabled attribute
	EnabledDefaultValue = true
)

// EnabledToModel maps the enabled flag returned by the API to the terraform model value.
// If the API did not return the flag, the existing plan/state value is preserved and a null or unknown value falls back to the default.
func EnabledToModel(apiValue *bool, existing types.Bool) types.Bool {
	if apiValue != nil {
		return types.BoolValue(*apiValue)
	}
	if existing.IsNull() || existing.IsUnknown() {
		return types.BoolValue(EnabledDefaultValue)
	}
	return existing
}

// EnabledSyncParams are the parameters of SyncEnabled
type EnabledSyncParams struct {
	// APIClient is used to call the enable/disable endpoints
	APIClient client.InstanaAPI
	// BasePath is the resource path of the alert configuration (e.g. /api/events/settings/infra-alert-configs)
	BasePath string
	// ID is the ID of the created/updated resource
	ID string
	// Label is a human readable name of the resource used in error messages
	Label string
	// Operation is the write operation (create/update) which has been executed before; only used in error messages
	Operation string
	// Plan contains the desired state
	Plan *tfsdk.Plan
	// PriorState is the state before the operation (nil for create)
	PriorState *tfsdk.State
	// APIValue is the enabled flag returned by the create/update API call (nil if the API did not return it)
	APIValue *bool
}

// SyncEnabled makes sure the enabled flag of a resource matches the planned value. The enabled flag is not part of the
// create/update payload of the alert configuration APIs, therefore the dedicated {id}/enable or {id}/disable endpoint
// is called when the desired value differs from the current one. The current value is taken from the API response first,
// then from the prior state and finally defaults to enabled (the default of newly created resources).
//
// It returns the effective enabled value after the operation, which is also correct when the endpoint call failed.
func SyncEnabled(ctx context.Context, params EnabledSyncParams) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var desiredValue types.Bool
	diags.Append(params.Plan.GetAttribute(ctx, path.Root(EnabledFieldName), &desiredValue)...)
	if diags.HasError() {
		return EnabledDefaultValue, diags
	}
	desired := EnabledDefaultValue
	if !desiredValue.IsNull() && !desiredValue.IsUnknown() {
		desired = desiredValue.ValueBool()
	}

	current := EnabledDefaultValue
	if params.APIValue != nil {
		current = *params.APIValue
	} else if params.PriorState != nil {
		var priorValue types.Bool
		diags.Append(params.PriorState.GetAttribute(ctx, path.Root(EnabledFieldName), &priorValue)...)
		if diags.HasError() {
			return current, diags
		}
		if !priorValue.IsNull() && !priorValue.IsUnknown() {
			current = priorValue.ValueBool()
		}
	}

	if desired == current {
		return current, diags
	}

	action := rest.ActionDisable
	if desired {
		action = rest.ActionEnable
	}
	err := params.APIClient.ResourceActions(params.BasePath).Execute(rest.ActionMethodPut, params.ID, action)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error trying to %s %s", action, params.Label),
			fmt.Sprintf("The %s %s was %sd but could not be %sd: %s", params.Label, params.ID, params.Operation, action, err),
		)
		return current, diags
	}
	return desired, diags
}
