package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInstanaAutomationPolicy the name of the terraform-provider-instana resource to manage automation policies
const ResourceInstanaAutomationPolicy = "instana_automation_policy"

// runnable types
const actionRunnable = "action"

// NewAutomationPolicyResourceHandle creates the resource handle for Automation Policies
func NewAutomationPolicyResourceHandle() ResourceHandle[*restapi.AutomationPolicy] {
	return &AutomationPolicyResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAutomationPolicy,
			Schema: map[string]*schema.Schema{
				AutomationPolicyFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the automation policy.",
				},
				AutomationPolicyFieldDescription: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The description of the automation policy.",
				},
				AutomationPolicyFieldTags: {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "The tags of the automation policy.",
				},
				AutomationPolicyFieldTrigger:           automationPolicyTriggerSchema,
				AutomationPolicyFieldTypeConfiguration: automationPolicyTypeConfigurationSchema,
			},
			SchemaVersion: 0,
		},
	}
}

type AutomationPolicyResource struct {
	metaData ResourceMetaData
}

func (r *AutomationPolicyResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *AutomationPolicyResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *AutomationPolicyResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AutomationPolicy] {
	return api.AutomationPolicies()
}

func (r *AutomationPolicyResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *AutomationPolicyResource) UpdateState(d *schema.ResourceData, policy *restapi.AutomationPolicy) error {
	d.SetId(policy.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AutomationPolicyFieldName:              policy.Name,
		AutomationPolicyFieldDescription:       policy.Description,
		AutomationPolicyFieldTags:              policy.Tags,
		AutomationPolicyFieldTrigger:           r.mapTriggerToSchema(policy),
		AutomationPolicyFieldTypeConfiguration: r.mapTypeConfigurationsToSchema(policy),
	})
}

func (r *AutomationPolicyResource) mapTriggerToSchema(policy *restapi.AutomationPolicy) []interface{} {
	return []interface{}{map[string]interface{}{
		AutomationPolicyFieldId:   policy.Trigger.Id,
		AutomationPolicyFieldType: policy.Trigger.Type,
	}}
}

func (r *AutomationPolicyResource) mapTypeConfigurationsToSchema(policy *restapi.AutomationPolicy) []interface{} {
	result := make([]interface{}, len(policy.TypeConfigurations))

	i := 0
	for _, v := range policy.TypeConfigurations {
		typeConfig := v

		item := make(map[string]interface{})
		item[AutomationPolicyFieldName] = typeConfig.Name
		item[AutomationPolicyFieldAction] = r.mapActionToSchema(&typeConfig)
		item[AutomationPolicyFieldCondition] = r.mapConditionToSchema(&typeConfig)
		result[i] = item
		i++
	}
	return result
}

func (r *AutomationPolicyResource) mapActionToSchema(typeConfiguration *restapi.TypeConfiguration) []interface{} {
	result := make([]interface{}, len(typeConfiguration.Runnable.RunConfiguration.Actions))

	i := 0
	for _, v := range typeConfiguration.Runnable.RunConfiguration.Actions {
		actionConfiguration := v

		item := make(map[string]interface{})
		item[AutomationPolicyFieldActionId] = actionConfiguration.Action.Id
		item[AutomationPolicyFieldAgentId] = actionConfiguration.AgentId
		item[AutomationPolicyFieldInputParameters] = r.mapParameterValuesToSchema(&actionConfiguration)
		result[i] = item

		i++
	}
	return result
}

func (r *AutomationPolicyResource) mapParameterValuesToSchema(ac *restapi.ActionConfiguration) map[string]string {
	parameterMap := make(map[string]string)
	for _, ipv := range ac.InputParameterValues {
		parameterMap[ipv.Name] = ipv.Value
	}
	return parameterMap
}

func (r *AutomationPolicyResource) mapConditionToSchema(typeConfiguration *restapi.TypeConfiguration) []interface{} {
	// FIXME: should we check if the condition is set?
	return []interface{}{map[string]interface{}{
		AutomationPolicyFieldQuery: typeConfiguration.Condition.Query,
	}}
}

func (r *AutomationPolicyResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AutomationPolicy, error) {
	trigger, err := r.mapTriggerFromSchema(d)
	if err != nil {
		return nil, err
	}

	typeConfigurations, err := r.mapTypeConfigurationsFromSchema(d)
	if err != nil {
		return nil, err
	}

	return &restapi.AutomationPolicy{
		ID:                 d.Id(),
		Name:               d.Get(AutomationActionFieldName).(string),
		Description:        d.Get(AutomationActionFieldDescription).(string),
		Tags:               d.Get(AutomationActionFieldTags),
		Trigger:            trigger,
		TypeConfigurations: typeConfigurations,
	}, nil
}

func (r *AutomationPolicyResource) mapTriggerFromSchema(d *schema.ResourceData) (restapi.Trigger, error) {
	val, ok := d.GetOk(AutomationPolicyFieldTrigger)

	if ok && len(val.([]interface{})) == 1 {
		triggerData := val.([]interface{})[0].(map[string]interface{})
		return restapi.Trigger{
			Id:   triggerData[AutomationPolicyFieldId].(string),
			Type: triggerData[AutomationPolicyFieldType].(string),
		}, nil
	}

	// FIXME: raise an error here, no trigger specified
	return restapi.Trigger{}, nil
}

func (r *AutomationPolicyResource) mapTypeConfigurationsFromSchema(d *schema.ResourceData) ([]restapi.TypeConfiguration, error) {
	val, ok := d.GetOk(AutomationPolicyFieldTypeConfiguration)

	if ok && val != nil {
		typeConfigurations := val.([]interface{})
		result := make([]restapi.TypeConfiguration, len(typeConfigurations))

		i := 0
		for _, v := range typeConfigurations {
			typeConfiguration := v.(map[string]interface{})

			runnable := r.getRunnableFromTypeConfigurationSchema(typeConfiguration)

			result[i] = restapi.TypeConfiguration{
				Name:      typeConfiguration[AutomationPolicyFieldName].(string),
				Condition: r.getConditionFromTypeConfigurationSchema(typeConfiguration),
				Runnable:  runnable,
			}
			i++
		}
		return result, nil
	}

	// FIXME: raise an error here, no type configuration specified
	return []restapi.TypeConfiguration{}, nil
}

func (r *AutomationPolicyResource) getConditionFromTypeConfigurationSchema(typeConfiguration map[string]interface{}) restapi.Condition {
	var query string

	// retrieve the condition from the type configuration
	conditions := typeConfiguration[AutomationPolicyFieldCondition].([]interface{})
	if len(conditions) == 1 {
		// extract the query from the condition
		query = conditions[0].(map[string]interface{})[AutomationPolicyFieldQuery].(string)
	}
	return restapi.Condition{
		Query: query,
	}
}

func (r *AutomationPolicyResource) getRunnableFromTypeConfigurationSchema(typeConfiguration map[string]interface{}) restapi.Runnable {
	// retrieve the actions list from type configuration
	actions := typeConfiguration[AutomationPolicyFieldAction].([]interface{})

	// FIXME: if no actions raise an error

	// currently we only support one action
	action := actions[0].(map[string]interface{})

	actionId := action[AutomationPolicyFieldActionId].(string)
	agentId := action[AutomationPolicyFieldAgentId].(string)

	return restapi.Runnable{
		Id:   actionId,
		Type: actionRunnable,
		RunConfiguration: restapi.RunConfiguration{
			Actions: []restapi.ActionConfiguration{
				{
					Action: restapi.Action{
						Id: actionId,
					},
					AgentId:              agentId,
					InputParameterValues: r.getInputParametersFromActionSchema(action),
				},
			},
		},
	}
}

func (r *AutomationPolicyResource) getInputParametersFromActionSchema(action map[string]interface{}) []restapi.InputParameterValue {
	if attr, ok := action[AutomationPolicyFieldInputParameters]; ok {
		parameters := attr.(map[string]interface{})

		result := make([]restapi.InputParameterValue, len(parameters))
		i := 0
		for key, value := range parameters {
			result[i] = restapi.InputParameterValue{
				Name:  key,
				Value: value.(string),
			}
			i++
		}

		return result
	}
	return []restapi.InputParameterValue{}
}
