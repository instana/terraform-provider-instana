package sessionsettings

import "github.com/hashicorp/terraform-plugin-framework/types"

// SessionSettingsModel is the Terraform model for session settings.
type SessionSettingsModel struct {
	TokenLifeTimeInMillis types.Int64 `tfsdk:"token_life_time_in_millis"`
	IdleTimeInMillis      types.Int64 `tfsdk:"idle_time_in_millis"`
}
