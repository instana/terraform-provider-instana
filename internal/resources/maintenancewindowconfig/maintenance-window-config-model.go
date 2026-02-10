package maintenancewindowconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// MaintenanceWindowConfigModel represents the data model for the maintenance window configuration resource
type MaintenanceWindowConfigModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Query                      types.String `tfsdk:"query"`
	Scheduling                 types.Object `tfsdk:"scheduling"`
	TagFilterExpressionEnabled types.Bool   `tfsdk:"tag_filter_expression_enabled"`
	TagFilterExpression        types.String `tfsdk:"tag_filter_expression"`
}

// MaintenanceSchedulingModel represents the scheduling configuration
type MaintenanceSchedulingModel struct {
	Start      types.Int64  `tfsdk:"start"`
	Duration   types.Object `tfsdk:"duration"`
	Type       types.String `tfsdk:"type"`
	Rrule      types.String `tfsdk:"rrule"`
	TimezoneId types.String `tfsdk:"timezone_id"`
}

// MaintenanceDurationModel represents the duration configuration
type MaintenanceDurationModel struct {
	Amount types.Int64  `tfsdk:"amount"`
	Unit   types.String `tfsdk:"unit"`
}

// Made with Bob
