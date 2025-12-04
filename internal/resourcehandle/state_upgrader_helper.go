package resourcehandle

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// CreatePassThroughStateUpgrader creates a state upgrader that passes through the state unchanged.
// This is useful when migrating from SDK v2 to Plugin Framework where the schema structure hasn't changed,
// but the framework requires explicit state upgraders for version transitions.
//
// IMPORTANT: When migrating from SDK v2 to Plugin Framework, the state representation changes from
// block schema to attribute schema. This means:
// 1. req.State will be nil because the old state format cannot be parsed into the new schema
// 2. req.RawState contains the raw state data from the state file
// 3. The framework will automatically handle the conversion from RawState to the new schema format
//
// This upgrader simply validates that RawState exists and lets the framework handle the conversion.
func CreatePassThroughStateUpgrader() func(context.Context, resource.UpgradeStateRequest, *resource.UpgradeStateResponse) {
	return func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
		// When migrating from SDK v2 to Plugin Framework:
		// - req.State is nil (old block schema cannot be parsed into new attribute schema)
		// - req.RawState contains the raw state data
		// - The framework will automatically convert RawState to the new schema format

		if req.State == nil {
			// This is expected during SDK v2 to Plugin Framework migration
			// Check if RawState exists - if it does, the framework will handle the conversion
			if req.RawState == nil {
				resp.Diagnostics.AddError(
					"State Upgrade Error",
					"If you are updating from SDK v2 to the Plugin Framework, you need to provide an explicit state upgrade for the current version. Please refer the migration guide.",
				)
				return
			}
			resp.Diagnostics.AddError(
				"State Upgrade Error",
				"If you are updating from SDK v2 to the Plugin Framework, you need to provide an explicit state upgrade for the current version. Please refer the migration guide.",
			)
			// RawState exists - let the framework handle the conversion
			// No need to set resp.State as the framework will populate it from RawState
			return
		}

		// If req.State is not nil (e.g., upgrading between Plugin Framework versions),
		// copy the state from request to response
		resp.State = *req.State
	}
}

// CreateStateUpgraderForVersion creates a state upgrader for a specific version
// that performs a pass-through upgrade (no data transformation needed)
func CreateStateUpgraderForVersion(fromVersion int64) resource.StateUpgrader {
	return resource.StateUpgrader{
		StateUpgrader: CreatePassThroughStateUpgrader(),
	}
}
