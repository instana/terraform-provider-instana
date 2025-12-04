package automationpolicy

// ResourceInstanaAutomationPolicy the name of the terraform-provider-instana resource to manage automation policies
const ResourceInstanaAutomationPolicy = "automation_policy"

// runnable types
const actionRunnable = "action"

// Resource description
const AutomationPolicyDescResource = "This resource manages automation policies in Instana."

// Field descriptions - ID
const AutomationPolicyDescID = "The ID of the automation policy."

// Field descriptions - Basic fields
const AutomationPolicyDescName = "The name of the automation policy."
const AutomationPolicyDescDescription = "The description of the automation policy."
const AutomationPolicyDescTags = "The tags of the automation policy."

// Field descriptions - Trigger
const AutomationPolicyDescTrigger = "The trigger for the automation policy."
const AutomationPolicyDescTriggerID = "Trigger (Instana event or Smart Alert) identifier."
const AutomationPolicyDescTriggerType = "Instana event or Smart Alert type."
const AutomationPolicyDescTriggerName = "The name of the trigger."
const AutomationPolicyDescTriggerDescription = "The description of the trigger."
const AutomationPolicyDescTriggerScheduling = "The scheduling configuration for the trigger."
const AutomationPolicyDescTriggerSchedulingStartTime = "Unix timestamp in milliseconds for the start time."
const AutomationPolicyDescTriggerSchedulingDuration = "Duration value."
const AutomationPolicyDescTriggerSchedulingDurationUnit = "Duration unit (MINUTE, HOUR, DAY)."
const AutomationPolicyDescTriggerSchedulingRecurrentRule = "Recurrent rule for scheduling."
const AutomationPolicyDescTriggerSchedulingRecurrent = "Whether the scheduling is recurrent."

// Field descriptions - Type Configuration
const AutomationPolicyDescTypeConfiguration = "A list of configurations with the list of actions to run and the mode (automatic or manual)."
const AutomationPolicyDescTypeConfigurationName = "The policy type."
const AutomationPolicyDescCondition = "The condition that selects a list of entities on which the policy is run. Only for automatic policy type."
const AutomationPolicyDescConditionQuery = "Dynamic Focus Query string that selects a list of entities on which the policy is run."

// Field descriptions - Action
const AutomationPolicyDescAction = "The configuration for the automation action."
const AutomationPolicyDescActionAction = "The automation action configuration."
const AutomationPolicyDescActionAgentID = "The identifier of the agent host."

// Field names
const (
	AutomationPolicyFieldID            = "id"
	AutomationPolicyFieldScheduling    = "scheduling"
	AutomationPolicyFieldStartTime     = "start_time"
	AutomationPolicyFieldDuration      = "duration"
	AutomationPolicyFieldDurationUnit  = "duration_unit"
	AutomationPolicyFieldRecurrentRule = "recurrent_rule"
	AutomationPolicyFieldRecurrent     = "recurrent"
	AutomationPolicyFieldAction        = "action"
)

// Duration unit values
const (
	DurationUnitMinute = "MINUTE"
	DurationUnitHour   = "HOUR"
	DurationUnitDay    = "DAY"
)

// Error messages
const AutomationPolicyErrMappingTags = "Error mapping tags"
const AutomationPolicyErrTagNotString = "Tag at index %d is not a string"
const AutomationPolicyErrTagsFormat = "Tags are not in the expected format"
