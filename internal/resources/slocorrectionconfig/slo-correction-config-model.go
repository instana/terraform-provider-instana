package slocorrectionconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// SloCorrectionConfigModel represents the data model for the SLO Correction Config resource
type SloCorrectionConfigModel struct {
	ID          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
	Active      types.Bool       `tfsdk:"active"`
	Scheduling  *SchedulingModel `tfsdk:"scheduling"`
	SloIds      types.Set        `tfsdk:"slo_ids"`
	Tags        types.Set        `tfsdk:"tags"`
}

// SchedulingModel represents the scheduling configuration for SLO Correction Config
type SchedulingModel struct {
	StartTime     types.Int64  `tfsdk:"start_time"`
	Duration      types.Int64  `tfsdk:"duration"`
	DurationUnit  types.String `tfsdk:"duration_unit"`
	RecurrentRule types.String `tfsdk:"recurrent_rule"`
	Recurrent     types.Bool   `tfsdk:"recurrent"`
}
