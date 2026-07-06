package apdexconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApdexConfigModel represents the data model for the Apdex configuration resource
type ApdexConfigModel struct {
	ID          types.String      `tfsdk:"id"`
	ApdexName   types.String      `tfsdk:"apdex_name"`
	Tags        types.Set         `tfsdk:"tags"`
	RbacTags    types.List        `tfsdk:"rbac_tags"`
	ApdexEntity *ApdexEntityModel `tfsdk:"apdex_entity"`
}

// RbacTagModel represents an RBAC tag in the Terraform model
type RbacTagModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
}

// ApdexEntityModel represents the polymorphic entity configuration
type ApdexEntityModel struct {
	ApplicationEntityModel *ApplicationApdexEntityModel `tfsdk:"application"`
	WebsiteEntityModel     *WebsiteApdexEntityModel     `tfsdk:"website"`
}

// ApplicationApdexEntityModel represents an application entity in the Terraform model
type ApplicationApdexEntityModel struct {
	EntityID         types.String `tfsdk:"entity_id"`
	Threshold        types.Int64  `tfsdk:"threshold"`
	BoundaryScope    types.String `tfsdk:"boundary_scope"`
	IncludeInternal  types.Bool   `tfsdk:"include_internal"`
	IncludeSynthetic types.Bool   `tfsdk:"include_synthetic"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}

// WebsiteApdexEntityModel represents a website entity in the Terraform model
type WebsiteApdexEntityModel struct {
	EntityID         types.String `tfsdk:"entity_id"`
	Threshold        types.Int64  `tfsdk:"threshold"`
	BeaconType       types.String `tfsdk:"beacon_type"`
	FilterExpression types.String `tfsdk:"filter_expression"`
}
