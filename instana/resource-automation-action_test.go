package instana_test

import (
	"encoding/json"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"

	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

const (
	actionId          = "actionId1"
	actionName        = "action4test"
	actionDescription = "action for unit test"
	actionTag1        = "tag1"
	actionTag2        = "tag2"
	actionTimeout     = "10"

	actionScriptContent     = "echo \"Hello world\""
	actionScriptInterpreter = "bash"

	actionHttpHost             = "http://localhost"
	actionHttpMethod           = "POST"
	actionHttpIgnoreCertErrors = "true"
	actionHttpBody             = "{\"name\":\"test\"}"
	actionHttpHeaderKey        = "Authentication"
	actionHttpHeaderValue      = "Bearer bearerToken"

	actionParamName        = "testParam"
	actionParamLabel       = "Parameter test"
	actionParamDescription = "Parameter for unit test"
	actionParamRequired    = true
	actionParamSecured     = true
	actionParamHidden      = false
	actionParamValue       = "testValue"
	actionParamType        = "static"
)

func TestAutomationActionResource(t *testing.T) {
	unitTest := &automationActionResourceUnitTest{}
	t.Run("schema should be valid", unitTest.resourceSchemaShouldBeValid)
	t.Run("schema version should be 0", unitTest.schemaShouldHaveVersion0)
	t.Run("should have no state upgraders", unitTest.shouldHaveNoStateUpgraders)
	t.Run("should return correct schema name", unitTest.shouldReturnCorrectResourceNameForAutomationAction)
	t.Run("should map script action to state", unitTest.shouldMapScriptActionToState)
	t.Run("should map http action to state", unitTest.shouldMapHttpActionToState)
	t.Run("should map script action from state", unitTest.shouldMapScriptActionFromState)
	t.Run("should map http action from state", unitTest.shouldMapHttpActionFromState)
}

type automationActionResourceUnitTest struct{}

func (r *automationActionResourceUnitTest) resourceSchemaShouldBeValid(t *testing.T) {
	resourceHandle := NewAutomationActionResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfStrings(AutomationActionFieldTags)

	r.validateScriptSchema(t, schemaMap[AutomationActionFieldScript].Elem.(*schema.Resource).Schema)
	r.validateHttpSchema(t, schemaMap[AutomationActionFieldHttp].Elem.(*schema.Resource).Schema)
	r.validateInputParameterSchema(t, schemaMap[AutomationActionFieldInputParameter].Elem.(*schema.Resource).Schema)
}

func (r *automationActionResourceUnitTest) schemaShouldHaveVersion0(t *testing.T) {
	require.Equal(t, 0, NewAutomationActionResourceHandle().MetaData().SchemaVersion)
}

func (r *automationActionResourceUnitTest) shouldHaveNoStateUpgraders(t *testing.T) {
	resourceHandler := NewAlertingChannelResourceHandle()

	require.Equal(t, 0, len(resourceHandler.StateUpgraders()))
}

func (r *automationActionResourceUnitTest) validateScriptSchema(t *testing.T, scriptSchema map[string]*schema.Schema) {
	require.Len(t, scriptSchema, 3)

	schemaAssert := testutils.NewTerraformSchemaAssert(scriptSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldContent)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationActionFieldInterpreter)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationActionFieldTimeout)
}

func (r *automationActionResourceUnitTest) validateHttpSchema(t *testing.T, httpSchema map[string]*schema.Schema) {
	require.Len(t, httpSchema, 6)

	schemaAssert := testutils.NewTerraformSchemaAssert(httpSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldHost)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldMethod)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationActionFieldBody)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(AutomationActionFieldIgnoreCertErrors, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeMapOfStrings(AutomationActionFieldHeaders)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationActionFieldTimeout)
}

func (r *automationActionResourceUnitTest) validateInputParameterSchema(t *testing.T, inputParamSchema map[string]*schema.Schema) {
	require.Len(t, inputParamSchema, 8)

	schemaAssert := testutils.NewTerraformSchemaAssert(inputParamSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionParameterFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionParameterFieldType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionParameterFieldValue)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationActionParameterFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AutomationActionParameterFieldLabel)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(AutomationActionParameterFieldSecured, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(AutomationActionParameterFieldHidden, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(AutomationActionParameterFieldRequired, true)
}

func (ut *automationActionResourceUnitTest) shouldReturnCorrectResourceNameForAutomationAction(t *testing.T) {
	name := NewAutomationActionResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_automation_action", name, "Expected resource name to be instana_automation_action")
}

func (r *automationActionResourceUnitTest) shouldMapScriptActionToState(t *testing.T) {
	data := restapi.AutomationAction{
		ID:          actionId,
		Name:        actionName,
		Description: actionDescription,
		Tags:        []string{actionTag1, actionTag2},
		Type:        "SCRIPT",
		Fields: []restapi.Field{
			{
				Name:        restapi.SCRIPT_SSH_FIELD_NAME,
				Description: restapi.SCRIPT_SSH_FIELD_NAME,
				Value:       actionScriptContent,
			},
			{
				Name:        restapi.SUBTYPE_FIELD_NAME,
				Description: restapi.SUBTYPE_FIELD_DESCRIPTION,
				Value:       actionScriptInterpreter,
			},
			{
				Name:        restapi.TIMEOUT_FIELD_NAME,
				Description: restapi.TIMEOUT_FIELD_DESCRIPTION,
				Value:       actionTimeout,
			},
		},
		InputParameters: []restapi.Parameter{
			{
				Name:        actionParamName,
				Value:       actionParamValue,
				Description: actionParamDescription,
				Label:       actionParamLabel,
				Secured:     actionParamSecured,
				Hidden:      actionParamHidden,
				Required:    actionParamRequired,
				Type:        actionParamType,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	sut := NewAutomationActionResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data)

	require.Nil(t, err)
	// verify the common fields
	r.assertActionResourceData(t, resourceData)
	// verify the script configuration
	r.assertScriptResourceData(t, resourceData)
	// verify the input parameters
	r.assertInputParameterResourceData(t, resourceData)
}

func (r *automationActionResourceUnitTest) shouldMapHttpActionToState(t *testing.T) {
	data := restapi.AutomationAction{
		ID:          actionId,
		Name:        actionName,
		Description: actionDescription,
		Tags:        []string{actionTag1, actionTag2},
		Type:        "HTTP",
		Fields: []restapi.Field{
			{
				Name:        restapi.HTTP_HOST_FIELD_NAME,
				Description: restapi.HTTP_HOST_FIELD_DESCRIPTION,
				Value:       actionHttpHost,
			},
			{
				Name:        restapi.HTTP_METHOD_FIELD_NAME,
				Description: restapi.HTTP_METHOD_FIELD_DESCRIPTION,
				Value:       actionHttpMethod,
			},
			{
				Name:        restapi.HTTP_BODY_FIELD_NAME,
				Description: restapi.HTTP_BODY_FIELD_DESCRIPTION,
				Value:       actionHttpBody,
			},
			{
				Name:        restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_NAME,
				Description: restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_DESCRIPTION,
				Value:       actionHttpIgnoreCertErrors,
			},
			{
				Name:        restapi.HTTP_HEADER_FIELD_NAME,
				Description: restapi.HTTP_HEADER_FIELD_DESCRIPTION,
				Value:       r.buildHeadersString(),
			},
			{
				Name:        restapi.TIMEOUT_FIELD_NAME,
				Description: restapi.TIMEOUT_FIELD_DESCRIPTION,
				Value:       actionTimeout,
			},
		},
		InputParameters: []restapi.Parameter{
			{
				Name:        actionParamName,
				Value:       actionParamValue,
				Description: actionParamDescription,
				Label:       actionParamLabel,
				Secured:     actionParamSecured,
				Hidden:      actionParamHidden,
				Required:    actionParamRequired,
				Type:        actionParamType,
			},
		},
	}

	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	sut := NewAutomationActionResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data)

	require.Nil(t, err)
	// verify common action fields
	r.assertActionResourceData(t, resourceData)
	// verify the http configuration
	r.assertHttpResourceData(t, resourceData)
	// verify the input parameters
	r.assertInputParameterResourceData(t, resourceData)
}

func (r *automationActionResourceUnitTest) shouldMapScriptActionFromState(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	resourceHandle := NewAutomationActionResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(actionId)
	setValueOnResourceData(t, resourceData, AutomationActionFieldName, actionName)
	setValueOnResourceData(t, resourceData, AutomationActionFieldDescription, actionDescription)
	setValueOnResourceData(t, resourceData, AutomationActionFieldTags, []string{actionTag1, actionTag2})
	setValueOnResourceData(t, resourceData, AutomationActionFieldScript, []interface{}{
		map[string]interface{}{
			AutomationActionFieldContent:     actionScriptContent,
			AutomationActionFieldInterpreter: actionScriptInterpreter,
			AutomationActionFieldTimeout:     actionTimeout,
		},
	})
	setValueOnResourceData(t, resourceData, AutomationActionFieldInputParameter, []interface{}{
		map[string]interface{}{
			AutomationActionParameterFieldName:        actionParamName,
			AutomationActionParameterFieldType:        actionParamType,
			AutomationActionParameterFieldDescription: actionParamDescription,
			AutomationActionParameterFieldLabel:       actionParamLabel,
			AutomationActionParameterFieldSecured:     actionParamSecured,
			AutomationActionParameterFieldHidden:      actionParamHidden,
			AutomationActionParameterFieldRequired:    actionParamRequired,
			AutomationActionParameterFieldValue:       actionParamValue,
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	// verify the action common attributes
	r.assertActionDataModel(t, result, "SCRIPT")
	// verify the action has the right `fields`
	r.assertScriptActionDataModelFields(t, result)
	// verify the action has the correct input parameters
	r.assertDataModelInputParameters(t, result)
}

func (r *automationActionResourceUnitTest) shouldMapHttpActionFromState(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	resourceHandle := NewAutomationActionResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(actionId)
	setValueOnResourceData(t, resourceData, AutomationActionFieldName, actionName)
	setValueOnResourceData(t, resourceData, AutomationActionFieldDescription, actionDescription)
	setValueOnResourceData(t, resourceData, AutomationActionFieldTags, []string{actionTag1, actionTag2})
	setValueOnResourceData(t, resourceData, AutomationActionFieldHttp, []interface{}{
		map[string]interface{}{
			AutomationActionFieldHost:             actionHttpHost,
			AutomationActionFieldMethod:           actionHttpMethod,
			AutomationActionFieldBody:             actionHttpBody,
			AutomationActionFieldIgnoreCertErrors: true,
			AutomationActionFieldTimeout:          actionTimeout,
			AutomationActionFieldHeaders: map[string]interface{}{
				actionHttpHeaderKey: actionHttpHeaderValue,
			},
		},
	})
	setValueOnResourceData(t, resourceData, AutomationActionFieldInputParameter, []interface{}{
		map[string]interface{}{
			AutomationActionParameterFieldName:        actionParamName,
			AutomationActionParameterFieldType:        actionParamType,
			AutomationActionParameterFieldDescription: actionParamDescription,
			AutomationActionParameterFieldLabel:       actionParamLabel,
			AutomationActionParameterFieldSecured:     actionParamSecured,
			AutomationActionParameterFieldHidden:      actionParamHidden,
			AutomationActionParameterFieldRequired:    actionParamRequired,
			AutomationActionParameterFieldValue:       actionParamValue,
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	// verify the action common attributes
	r.assertActionDataModel(t, result, "HTTP")
	// verify the action has the right `fields`
	r.assertHttpActionDataModelFields(t, result)
	// verify the action has the correct input parameters
	r.assertDataModelInputParameters(t, result)
}

func (r *automationActionResourceUnitTest) assertActionResourceData(t *testing.T, resourceData *schema.ResourceData) {
	require.Equal(t, actionId, resourceData.Id())
	require.Equal(t, actionName, resourceData.Get(AutomationActionFieldName))
	require.Equal(t, actionDescription, resourceData.Get(AutomationActionFieldDescription))

	// assert action tags
	tags := resourceData.Get(AutomationActionFieldTags).([]interface{})
	require.Len(t, tags, 2)
	require.Contains(t, tags, actionTag1)
	require.Contains(t, tags, actionTag2)
}

func (r *automationActionResourceUnitTest) assertScriptResourceData(t *testing.T, resourceData *schema.ResourceData) {
	require.IsType(t, []interface{}{}, resourceData.Get(AutomationActionFieldScript))
	require.Len(t, resourceData.Get(AutomationActionFieldScript).([]interface{}), 1)

	script := resourceData.Get(AutomationActionFieldScript).([]interface{})[0].(map[string]interface{})
	require.Len(t, script, 3)
	require.Equal(t, actionScriptContent, script[AutomationActionFieldContent])
	require.Equal(t, actionScriptInterpreter, script[AutomationActionFieldInterpreter])
	require.Equal(t, actionTimeout, script[AutomationActionFieldTimeout])
}

func (r *automationActionResourceUnitTest) assertHttpResourceData(t *testing.T, resourceData *schema.ResourceData) {
	require.IsType(t, []interface{}{}, resourceData.Get(AutomationActionFieldScript))
	require.Len(t, resourceData.Get(AutomationActionFieldHttp).([]interface{}), 1)

	http := resourceData.Get(AutomationActionFieldHttp).([]interface{})[0].(map[string]interface{})
	require.Len(t, http, 6)
	require.Equal(t, actionHttpHost, http[AutomationActionFieldHost])
	require.Equal(t, actionHttpMethod, http[AutomationActionFieldMethod])
	require.Equal(t, actionHttpBody, http[AutomationActionFieldBody])
	require.Equal(t, true, http[AutomationActionFieldIgnoreCertErrors])
	require.Equal(t, actionTimeout, http[AutomationActionFieldTimeout])

	headers := http[AutomationActionFieldHeaders].(map[string]interface{})
	require.Len(t, headers, 1)
	require.Equal(t, actionHttpHeaderValue, headers[actionHttpHeaderKey])
}

func (r *automationActionResourceUnitTest) assertInputParameterResourceData(t *testing.T, resourceData *schema.ResourceData) {
	inputParameter := resourceData.Get(AutomationActionFieldInputParameter).(*schema.Set).List()[0].(map[string]interface{})

	require.Len(t, inputParameter, 8)
	require.Equal(t, actionParamName, inputParameter[AutomationActionParameterFieldName])
	require.Equal(t, actionParamDescription, inputParameter[AutomationActionParameterFieldDescription])
	require.Equal(t, actionParamLabel, inputParameter[AutomationActionParameterFieldLabel])
	require.Equal(t, actionParamValue, inputParameter[AutomationActionParameterFieldValue])
	require.Equal(t, actionParamSecured, inputParameter[AutomationActionParameterFieldSecured])
	require.Equal(t, actionParamRequired, inputParameter[AutomationActionParameterFieldRequired])
	require.Equal(t, actionParamHidden, inputParameter[AutomationActionParameterFieldHidden])
	require.Equal(t, actionParamType, inputParameter[AutomationActionParameterFieldType])
}

func (r *automationActionResourceUnitTest) assertActionDataModel(t *testing.T, dataModel *restapi.AutomationAction, expectedType string) {
	require.Equal(t, actionId, dataModel.GetIDForResourcePath())
	require.Equal(t, actionName, dataModel.Name)
	require.Equal(t, actionDescription, dataModel.Description)
	require.Equal(t, expectedType, dataModel.Type)
	require.Len(t, dataModel.Tags, 2)
	require.Contains(t, dataModel.Tags, actionTag1)
	require.Contains(t, dataModel.Tags, actionTag2)
}

func (r *automationActionResourceUnitTest) assertScriptActionDataModelFields(t *testing.T, dataModel *restapi.AutomationAction) {
	require.Len(t, dataModel.Fields, 3)
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.SCRIPT_SSH_FIELD_NAME, restapi.SCRIPT_SSH_FIELD_DESCRIPTION, actionScriptContent, "base64")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.SUBTYPE_FIELD_NAME, restapi.SUBTYPE_FIELD_DESCRIPTION, actionScriptInterpreter, "ascii")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.TIMEOUT_FIELD_NAME, restapi.TIMEOUT_FIELD_DESCRIPTION, actionTimeout, "ascii")

}

func (r *automationActionResourceUnitTest) assertHttpActionDataModelFields(t *testing.T, dataModel *restapi.AutomationAction) {
	require.Len(t, dataModel.Fields, 6)
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.HTTP_HOST_FIELD_NAME, restapi.HTTP_HOST_FIELD_DESCRIPTION, actionHttpHost, "ascii")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.HTTP_METHOD_FIELD_NAME, restapi.HTTP_METHOD_FIELD_DESCRIPTION, actionHttpMethod, "ascii")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.HTTP_BODY_FIELD_NAME, restapi.HTTP_BODY_FIELD_DESCRIPTION, actionHttpBody, "ascii")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.HTTP_HEADER_FIELD_NAME, restapi.HTTP_HEADER_FIELD_DESCRIPTION, r.buildHeadersString(), "ascii")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_NAME, restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_DESCRIPTION, actionHttpIgnoreCertErrors, "ascii")
	r.assertFieldsContains(t, dataModel.Fields,
		restapi.TIMEOUT_FIELD_NAME, restapi.TIMEOUT_FIELD_DESCRIPTION, actionTimeout, "ascii")
}

func (r *automationActionResourceUnitTest) assertDataModelInputParameters(t *testing.T, dataModel *restapi.AutomationAction) {
	require.Len(t, dataModel.InputParameters, 1)
	inputParam := dataModel.InputParameters[0]
	require.Equal(t, actionParamName, inputParam.Name)
	require.Equal(t, actionParamDescription, inputParam.Description)
	require.Equal(t, actionParamLabel, inputParam.Label)
	require.Equal(t, actionParamSecured, inputParam.Secured)
	require.Equal(t, actionParamHidden, inputParam.Hidden)
	require.Equal(t, actionParamRequired, inputParam.Required)
	require.Equal(t, actionParamValue, inputParam.Value)
}

func (r *automationActionResourceUnitTest) buildHeadersString() string {
	headers := map[string]string{
		actionHttpHeaderKey: actionHttpHeaderValue,
	}
	headersString, err := json.Marshal(headers)
	if err != nil {
		return ""
	}
	return string(headersString)
}

func (r *automationActionResourceUnitTest) assertFieldsContains(t *testing.T, fields []restapi.Field,
	fieldName string, fieldDescription string, fieldValue string, fieldEncoding string) {

	var field restapi.Field
	for _, f := range fields {
		if f.Name == fieldName {
			field = f
		}
	}
	require.NotNil(t, field)
	require.Equal(t, fieldValue, field.Value)
	require.Equal(t, fieldDescription, field.Description)
	require.Equal(t, fieldEncoding, field.Encoding)
}
