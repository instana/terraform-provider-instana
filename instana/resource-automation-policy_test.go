package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/require"
)

const (
	policyId              = "policyId1"
	policyName            = "policy4test"
	policyActionId        = "actionId1"
	policyDescription     = "policy for unit test"
	policyTag1            = "tag1"
	policyTag2            = "tag2"
	triggerId             = "eventId1"
	triggerType           = "customEvent"
	typeConfigurationName = "manual"
	queryTest             = "entity.agent.capability:action"
	runnableType          = "action"
	agentId               = "agentId1"
	parameterName         = "parameter1"
	parameterValue        = "parameterValue1"
)

func TestAutomationPolicyResource(t *testing.T) {
	unitTest := &automationPolicyResourceUnitTest{}
	t.Run("schema should be valid", unitTest.resourceSchemaShouldBeValid)
	t.Run("schema version should be 0", unitTest.schemaShouldHaveVersion0)
	t.Run("should have no state upgraders", unitTest.shouldHaveNoStateUpgraders)
	t.Run("should return correct schema name", unitTest.shouldReturnCorrectResourceNameForAutomationPolicy)
	t.Run("should map policy to state", unitTest.shouldMapPolicyToState)
	t.Run("should map policy from state", unitTest.shouldMapPolicyFromState)
}

type automationPolicyResourceUnitTest struct{}

func (ut *automationPolicyResourceUnitTest) resourceSchemaShouldBeValid(t *testing.T) {
	resourceHandle := NewAutomationPolicyResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfStrings(AutomationPolicyFieldTags)

	ut.validateTriggerSchema(t, schemaMap[AutomationPolicyFieldTrigger].Elem.(*schema.Resource).Schema)
	ut.validateTypeConfigurationSchema(t, schemaMap[AutomationPolicyFieldTypeConfiguration].Elem.(*schema.Resource).Schema)
}

func (ut *automationPolicyResourceUnitTest) validateTypeConfigurationSchema(t *testing.T, typeConfigSchema map[string]*schema.Schema) {
	require.Len(t, typeConfigSchema, 3)

	schemaAssert := testutils.NewTerraformSchemaAssert(typeConfigSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldName)

	ut.validateConditionSchema(t, typeConfigSchema[AutomationPolicyFieldCondition].Elem.(*schema.Resource).Schema)
	ut.validateActionSchema(t, typeConfigSchema[AutomationPolicyFieldAction].Elem.(*schema.Resource).Schema)
}

func (ut *automationPolicyResourceUnitTest) validateConditionSchema(t *testing.T, conditionSchema map[string]*schema.Schema) {
	require.Len(t, conditionSchema, 1)

	schemaAssert := testutils.NewTerraformSchemaAssert(conditionSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldQuery)
}

func (ut *automationPolicyResourceUnitTest) validateTriggerSchema(t *testing.T, triggerSchema map[string]*schema.Schema) {
	require.Len(t, triggerSchema, 2)

	schemaAssert := testutils.NewTerraformSchemaAssert(triggerSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldId)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldType)
}

func (ut *automationPolicyResourceUnitTest) validateActionSchema(t *testing.T, actionSchema map[string]*schema.Schema) {
	require.Len(t, actionSchema, 3)

	schemaAssert := testutils.NewTerraformSchemaAssert(actionSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationPolicyFieldActionId)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationPolicyFieldAgentId)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeMapOfStrings(AutomationPolicyFieldInputParameters)
}

func (ut *automationPolicyResourceUnitTest) schemaShouldHaveVersion0(t *testing.T) {
	require.Equal(t, 0, NewAutomationPolicyResourceHandle().MetaData().SchemaVersion)
}

func (ut *automationPolicyResourceUnitTest) shouldHaveNoStateUpgraders(t *testing.T) {
	resourceHandler := NewAutomationPolicyResourceHandle()

	require.Equal(t, 0, len(resourceHandler.StateUpgraders()))
}

func (ut *automationPolicyResourceUnitTest) shouldReturnCorrectResourceNameForAutomationPolicy(t *testing.T) {
	name := NewAutomationPolicyResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_automation_policy", name, "Expected resource name to be instana_automation_policy")
}

func (ut *automationPolicyResourceUnitTest) shouldMapPolicyToState(t *testing.T) {
	data := restapi.AutomationPolicy{
		ID:          policyId,
		Name:        policyName,
		Description: policyDescription,
		Tags:        []string{policyTag1, policyTag2},
		Trigger: restapi.Trigger{
			Id:   triggerId,
			Type: triggerType,
		},
		TypeConfigurations: []restapi.TypeConfiguration{
			{
				Name: typeConfigurationName,
				Condition: restapi.Condition{
					Query: queryTest,
				},
				Runnable: restapi.Runnable{
					Id:   policyActionId,
					Type: runnableType,
					RunConfiguration: restapi.RunConfiguration{
						Actions: []restapi.ActionConfiguration{
							{
								Action: restapi.Action{
									Id: policyActionId,
								},
								AgentId: agentId,
								InputParameterValues: []restapi.InputParameterValue{
									{
										Name:  parameterName,
										Value: parameterValue,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	testHelper := NewTestHelper[*restapi.AutomationPolicy](t)
	sut := NewAutomationPolicyResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data)

	require.Nil(t, err)
	ut.assertPolicyResourceData(t, resourceData)
}

func (ut *automationPolicyResourceUnitTest) shouldMapPolicyFromState(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AutomationPolicy](t)
	resourceHandle := NewAutomationPolicyResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(policyId)
	setValueOnResourceData(t, resourceData, AutomationPolicyFieldName, policyName)
	setValueOnResourceData(t, resourceData, AutomationPolicyFieldDescription, policyDescription)
	setValueOnResourceData(t, resourceData, AutomationPolicyFieldTags, []string{policyTag1, policyTag2})
	setValueOnResourceData(t, resourceData, AutomationPolicyFieldTrigger, []interface{}{
		map[string]interface{}{
			AutomationPolicyFieldId:   triggerId,
			AutomationPolicyFieldType: triggerType,
		},
	})
	setValueOnResourceData(t, resourceData, AutomationPolicyFieldTypeConfiguration, []interface{}{
		map[string]interface{}{
			AutomationPolicyFieldName: typeConfigurationName,
			AutomationPolicyFieldCondition: []interface{}{
				map[string]interface{}{
					AutomationPolicyFieldQuery: queryTest,
				},
			},
			AutomationPolicyFieldAction: []interface{}{
				map[string]interface{}{
					AutomationPolicyFieldActionId: policyActionId,
					AutomationPolicyFieldAgentId:  agentId,
					AutomationPolicyFieldInputParameters: map[string]string{
						parameterName: parameterValue,
					},
				},
			},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	// verify the policy common attributes
	ut.assertPolicyDataModel(t, result)
}

func (ut *automationPolicyResourceUnitTest) assertPolicyResourceData(t *testing.T, resourceData *schema.ResourceData) {
	// verify common policy fields
	require.Equal(t, policyId, resourceData.Id())
	require.Equal(t, policyName, resourceData.Get(AutomationPolicyFieldName))
	require.Equal(t, policyDescription, resourceData.Get(AutomationPolicyFieldDescription))

	// verify policy tags
	tags := resourceData.Get(AutomationPolicyFieldTags).([]interface{})
	require.Len(t, tags, 2)
	require.Contains(t, tags, policyTag1)
	require.Contains(t, tags, policyTag2)

	// verify trigger
	ut.assertPolicyTriggerResourceData(t, resourceData)

	// verify type configuration
	ut.assertPolicyTypeConfigurationResourceData(t, resourceData)
}

func (ut *automationPolicyResourceUnitTest) assertPolicyTriggerResourceData(t *testing.T, resourceData *schema.ResourceData) {
	triggers := resourceData.Get(AutomationPolicyFieldTrigger).([]interface{})
	require.Len(t, triggers, 1)

	trigger := triggers[0].(map[string]interface{})
	require.Equal(t, triggerId, trigger[AutomationPolicyFieldId])
	require.Equal(t, triggerType, trigger[AutomationPolicyFieldType])
}

func (ut *automationPolicyResourceUnitTest) assertPolicyTypeConfigurationResourceData(t *testing.T, resourceData *schema.ResourceData) {
	typeConfigurations := resourceData.Get(AutomationPolicyFieldTypeConfiguration).([]interface{})
	require.Len(t, typeConfigurations, 1)
	typeConfiguration := typeConfigurations[0].(map[string]interface{})
	require.Equal(t, typeConfigurationName, typeConfiguration[AutomationPolicyFieldName])

	// verify the condition of the type configuration
	conditions := typeConfiguration[AutomationPolicyFieldCondition].([]interface{})
	require.Len(t, conditions, 1)
	condition := conditions[0].(map[string]interface{})
	require.Equal(t, queryTest, condition[AutomationPolicyFieldQuery])

	// verify the action of the type configuration
	actions := typeConfiguration[AutomationPolicyFieldAction].([]interface{})
	require.Len(t, actions, 1)
	action := actions[0].(map[string]interface{})
	require.Equal(t, policyActionId, action[AutomationPolicyFieldActionId])
	require.Equal(t, agentId, action[AutomationPolicyFieldAgentId])
	inputParameters := action[AutomationPolicyFieldInputParameters].(map[string]interface{})
	require.Len(t, inputParameters, 1)
	require.Equal(t, parameterValue, inputParameters[parameterName])
}

func (ut *automationPolicyResourceUnitTest) assertPolicyDataModel(t *testing.T, dataModel *restapi.AutomationPolicy) {
	require.Equal(t, policyId, dataModel.GetIDForResourcePath())
	require.Equal(t, policyName, dataModel.Name)
	require.Equal(t, policyDescription, dataModel.Description)
	require.Len(t, dataModel.Tags, 2)
	require.Contains(t, dataModel.Tags, policyTag1)
	require.Contains(t, dataModel.Tags, policyTag2)

	require.Equal(t, triggerId, dataModel.Trigger.Id)
	require.Equal(t, triggerType, dataModel.Trigger.Type)

	require.Len(t, dataModel.TypeConfigurations, 1)
	typeConfiguration := dataModel.TypeConfigurations[0]
	require.Equal(t, typeConfigurationName, typeConfiguration.Name)
	require.Equal(t, queryTest, typeConfiguration.Condition.Query)

	require.Equal(t, policyActionId, typeConfiguration.Runnable.Id)
	require.Equal(t, "action", typeConfiguration.Runnable.Type)

	require.Len(t, typeConfiguration.Runnable.RunConfiguration.Actions, 1)
	action := typeConfiguration.Runnable.RunConfiguration.Actions[0]
	require.Equal(t, policyActionId, action.Action.Id)
	require.Equal(t, agentId, action.AgentId)
	require.Len(t, action.InputParameterValues, 1)
	require.Equal(t, parameterName, action.InputParameterValues[0].Name)
	require.Equal(t, parameterValue, action.InputParameterValues[0].Value)
}
