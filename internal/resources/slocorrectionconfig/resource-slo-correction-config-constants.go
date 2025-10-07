package slocorrectionconfig

// ResourceInstanaSloCorrectionConfigFramework the name of the terraform-provider-instana resource to manage SLO correction configurations
const ResourceInstanaSloCorrectionConfigFramework = "slo_correction_config"

const (
	// Slo Correction Config Field names for Terraform
	SloCorrectionConfigFieldName                    = "name"
	SloCorrectionConfigFieldFullName                = "full_name"
	SloCorrectionConfigFieldDescription             = "description"
	SloCorrectionConfigFieldActive                  = "active"
	SloCorrectionConfigFieldScheduling              = "scheduling"
	SloCorrectionConfigFieldSloIds                  = "slo_ids"
	SloCorrectionConfigFieldTags                    = "tags"
	SloCorrectionConfigFieldSchedulingStartTime     = "start_time"
	SloCorrectionConfigFieldSchedulingDuration      = "duration"
	SloCorrectionConfigFieldSchedulingDurationUnit  = "duration_unit"
	SloCorrectionConfigFieldSchedulingRecurrentRule = "recurrent_rule"
)

const (
	// Resource description constants
	SloCorrectionConfigDescResource      = "This resource manages SLO Correction Configurations in Instana."
	SloCorrectionConfigDescID            = "The ID of the SLO Correction Config."
	SloCorrectionConfigDescName          = "The name of the SLO Correction Config."
	SloCorrectionConfigDescDescription   = "The description of the SLO Correction Config."
	SloCorrectionConfigDescActive        = "Indicates whether the Correction Config is active."
	SloCorrectionConfigDescSloIds        = "A set of SLO IDs that this correction config applies to."
	SloCorrectionConfigDescTags          = "A list of tags to be associated with the SLO Correction Config."
	SloCorrectionConfigDescScheduling    = "Scheduling configuration for the SLO Correction Config."
	SloCorrectionConfigDescStartTime     = "The start time of the scheduling in Unix timestamp in milliseconds."
	SloCorrectionConfigDescDuration      = "The duration of the scheduling in the specified unit."
	SloCorrectionConfigDescDurationUnit  = "The unit of the duration (e.g.,'minute' 'hour', 'day')."
	SloCorrectionConfigDescRecurrentRule = "Recurrent rule for scheduling, if applicable."
	SloCorrectionConfigDescRecurrent     = "Indicates whether the Rule is reccurrent"
)
