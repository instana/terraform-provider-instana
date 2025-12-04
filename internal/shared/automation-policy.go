package shared

var SupportedTriggerTypes = []string{
	"customEvent",
	"builtinEvent",
	"applicationSmartAlert",
	"globalApplicationSmartAlert",
	"websiteSmartAlert",
	"infraSmartAlert",
	"mobileAppSmartAlert",
	"syntheticsSmartAlert",
	"logSmartAlert",
	"sloSmartAlert",
	"schedule",
}

var SupportedPolicyTypes = []string{
	"manual",
	"automatic",
}

const (
	AutomationPolicyFieldId          = "id"
	AutomationPolicyFieldType        = "type"
	AutomationPolicyFieldName        = "name"
	AutomationPolicyFieldDescription = "description"
	AutomationPolicyFieldTags        = "tags"

	AutomationPolicyFieldTrigger           = "trigger"
	AutomationPolicyFieldTypeConfiguration = "type_configuration"
	AutomationPolicyFieldCondition         = "condition"
	AutomationPolicyFieldQuery             = "query"
	AutomationPolicyFieldAction            = "action"
	AutomationPolicyFieldActionId          = "action_id"
	AutomationPolicyFieldAgentId           = "agent_id"
	AutomationPolicyFieldInputParameters   = "input_parameters"
)
