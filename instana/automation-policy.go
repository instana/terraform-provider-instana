package instana

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var supportedTriggerTypes = []string{
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
}

var supportedPolicyTypes = []string{
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

var (
	automationPolicyTriggerSchema = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "The trigger for the automation policy.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationPolicyFieldId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Trigger (Instana event or Smart Alert) identifier.",
				},
				AutomationPolicyFieldType: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Instana event or Smart Alert type.",
					ValidateFunc: validation.StringInSlice(supportedTriggerTypes, false),
				},
			},
		},
	}
	automationPolicyConditionSchema = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    0,
		MaxItems:    1,
		Optional:    true,
		Description: "The condition that selects a list of entities on which the policy is run. Only for automatic policy type.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationPolicyFieldQuery: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Dynamic Focus Query string that selects a list of entities on which the policy is run.",
				},
			},
		},
	}
	automationPolicyActionSchema = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		Required:    true,
		Description: "The configuration for the automation action.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationPolicyFieldActionId: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The identifier for the automation action.",
				},
				AutomationPolicyFieldAgentId: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The identifier of the agent host.",
				},
				AutomationPolicyFieldInputParameters: {
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "Optional map with input parameters name and value.",
				},
			},
		},
	}
	automationPolicyTypeConfigurationSchema = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		Required:    true,
		Description: "A list of configurations with the list of actions to run and the mode (automatic or manual).",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				AutomationPolicyFieldName: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The policy type.",
					ValidateFunc: validation.StringInSlice(supportedPolicyTypes, false),
				},
				AutomationPolicyFieldCondition: automationPolicyConditionSchema,
				AutomationPolicyFieldAction:    automationPolicyActionSchema,
			},
		},
	}
)
