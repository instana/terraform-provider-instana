package ipfiltering

import "github.com/hashicorp/terraform-plugin-framework/types"

// IPFilteringRuleModel is the Terraform model for a single IP filtering rule.
type IPFilteringRuleModel struct {
	Target types.String `tfsdk:"target"`
	Block  types.Bool   `tfsdk:"block"`
}

// IPFilteringModel is the Terraform model for IP filtering.
type IPFilteringModel struct {
	Enabled              types.Bool             `tfsdk:"enabled"`
	DenyAll              types.Bool             `tfsdk:"deny_all"`
	SupportAccessEnabled types.Bool             `tfsdk:"support_access_enabled"`
	Rules                []IPFilteringRuleModel `tfsdk:"rules"`
	Active               types.Bool             `tfsdk:"active"`
	LastChangedAt        types.Int64            `tfsdk:"last_changed_at"`
	LastVerifiedAt       types.Int64            `tfsdk:"last_verified_at"`
}
