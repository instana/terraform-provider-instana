package slocorrectionconfig

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
