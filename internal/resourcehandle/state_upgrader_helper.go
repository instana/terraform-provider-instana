package resourcehandle

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// CreatePassThroughStateUpgrader creates a state upgrader that passes through the state unchanged.
// This is useful when migrating from SDK v2 to Plugin Framework where the schema structure hasn't changed,
// but the framework requires explicit state upgraders for version transitions.
func CreatePassThroughStateUpgrader() func(context.Context, resource.UpgradeStateRequest, *resource.UpgradeStateResponse) {
	return func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
		// Copy the state from request to response
		// The framework will handle the conversion to the new schema format
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

// Made with Bob
