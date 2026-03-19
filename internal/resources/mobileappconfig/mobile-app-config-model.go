package mobileappconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppConfigModel represents the Terraform model for a Mobile App Configuration
type MobileAppConfigModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

