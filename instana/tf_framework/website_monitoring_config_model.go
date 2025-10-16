package tf_framework

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WebsiteMonitoringConfigModel represents the Terraform model for a Website Monitoring Configuration
type WebsiteMonitoringConfigModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	AppName types.String `tfsdk:"app_name"`
}

// Made with Bob
