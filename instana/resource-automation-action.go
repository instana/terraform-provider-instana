package instana

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInstanaAutomationAction the name of the terraform-provider-instana resource to manage automation actions
const ResourceInstanaAutomationAction = "instana_automation_action"

// action types
const ActionTypeScript = "SCRIPT"
const ActionTypeHttp = "HTTP"

// encodings
const AsciiEncoding = "ascii"
const Base64Encoding = "base64"
const UTF8Encoding = "UTF8"

// NewAutomationActionResourceHandle creates the resource handle for Automation Actions
func NewAutomationActionResourceHandle() ResourceHandle[*restapi.AutomationAction] {
	return &AutomationActionResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAutomationAction,
			Schema: map[string]*schema.Schema{
				AutomationActionFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the automation action.",
				},
				AutomationActionFieldDescription: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The description of the automation action.",
				},
				AutomationActionFieldTags: {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "The tags of the automation action.",
				},
				AutomationActionFieldScript:         automationActionScriptSchema,
				AutomationActionFieldHttp:           automationActionHttpSchema,
				AutomationActionFieldInputParameter: automationActionInputParameterSchema,
			},
			SchemaVersion: 0,
		},
	}
}

type AutomationActionResource struct {
	metaData ResourceMetaData
}

func (r *AutomationActionResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *AutomationActionResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *AutomationActionResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AutomationAction] {
	return api.AutomationActions()
}

func (r *AutomationActionResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *AutomationActionResource) UpdateState(d *schema.ResourceData, automationAction *restapi.AutomationAction) error {
	// convert input parameters to schema
	inputParameters, err := r.mapInputParametersToSchema(automationAction)
	if err != nil {
		return err
	}

	// convert script configuration to schema
	scriptConfig, err := r.mapScriptFieldsToSchema(automationAction)
	if err != nil {
		return err
	}

	// convert http configuration to schema
	httpConfig, err := r.mapHttpFieldsToSchema(automationAction)
	if err != nil {
		return err
	}

	d.SetId(automationAction.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AutomationActionFieldName:           automationAction.Name,
		AutomationActionFieldDescription:    automationAction.Description,
		AutomationActionFieldTags:           automationAction.Tags,
		AutomationActionFieldInputParameter: inputParameters,
		AutomationActionFieldScript:         scriptConfig,
		AutomationActionFieldHttp:           httpConfig,
	})
}

func (r *AutomationActionResource) mapInputParametersToSchema(action *restapi.AutomationAction) ([]interface{}, error) {
	result := make([]interface{}, len(action.InputParameters))

	i := 0
	for _, v := range action.InputParameters {
		val := v

		item := make(map[string]interface{})
		item[AutomationActionParameterFieldName] = val.Name
		item[AutomationActionParameterFieldDescription] = val.Description
		item[AutomationActionParameterFieldLabel] = val.Label
		item[AutomationActionParameterFieldRequired] = val.Required
		item[AutomationActionParameterFieldHidden] = val.Hidden
		item[AutomationActionParameterFieldType] = val.Type
		item[AutomationActionParameterFieldValue] = val.Value

		result[i] = item
		i++
	}
	return result, nil
}

func (r *AutomationActionResource) mapScriptFieldsToSchema(action *restapi.AutomationAction) ([]interface{}, error) {
	if action.Type == ActionTypeScript {
		result := make(map[string]interface{})

		// script content is required field
		result[AutomationActionFieldContent] = r.getFieldValue(action, restapi.ScriptSshFieldName)
		// interpreter is optional field
		interpreter := r.getFieldValue(action, restapi.SubtypeFieldName)
		if len(interpreter) > 0 {
			result[AutomationActionFieldInterpreter] = interpreter
		}
		// timeout is optional field
		timeout := r.getFieldValue(action, restapi.TimeoutFieldName)
		if len(timeout) > 0 {
			result[AutomationActionFieldTimeout] = timeout
		}

		return []interface{}{result}, nil
	}

	return []interface{}{}, nil
}

func (r *AutomationActionResource) getFieldValue(action *restapi.AutomationAction, fieldName string) string {
	for _, v := range action.Fields {
		if v.Name == fieldName {
			return v.Value
		}
	}
	return ""
}

func (r *AutomationActionResource) getBoolFieldValueOrDefault(action *restapi.AutomationAction, fieldName string, defaultValue bool) bool {
	for _, v := range action.Fields {
		if v.Name == fieldName {
			boolValue, err := strconv.ParseBool(v.Value)
			if err != nil {
				fmt.Println("Error parsing value ", v.Value, " defaulting to ", defaultValue)
				return defaultValue
			} else {
				return boolValue
			}
		}
	}
	return defaultValue
}

func (r *AutomationActionResource) mapHttpFieldsToSchema(action *restapi.AutomationAction) ([]interface{}, error) {
	if action.Type == ActionTypeHttp {
		result := make(map[string]interface{})

		// required fields
		result[AutomationActionFieldHost] = r.getFieldValue(action, restapi.HttpHostFieldName)
		result[AutomationActionFieldMethod] = r.getFieldValue(action, restapi.HttpMethodFieldName)

		// optional fields
		body := r.getFieldValue(action, restapi.HttpBodyFieldName)
		if len(body) > 0 {
			result[AutomationActionFieldBody] = body
		}
		ignoreCertErrors := r.getBoolFieldValueOrDefault(action, restapi.HttpIgnoreCertErrorsFieldName, false)
		result[AutomationActionFieldIgnoreCertErrors] = ignoreCertErrors

		headersData := r.getFieldValue(action, restapi.HttpHeaderFieldName)
		headers, err := r.mapHttpHeadersToSchema(headersData)
		if err != nil {
			return nil, err
		}
		result[AutomationActionFieldHeaders] = headers

		timeout := r.getFieldValue(action, restapi.TimeoutFieldName)
		if len(timeout) > 0 {
			result[AutomationActionFieldTimeout] = timeout
		}

		return []interface{}{result}, nil
	}
	return []interface{}{}, nil
}

func (r *AutomationActionResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AutomationAction, error) {
	actionType, err := r.mapActionTypeFromSchema(d)
	if err != nil {
		return nil, err
	}

	inputParameters, err := r.mapInputParametersFromSchema(d)
	if err != nil {
		return nil, err
	}

	fields, err := r.mapActionFieldsFromSchema(d)
	if err != nil {
		return nil, err
	}

	return &restapi.AutomationAction{
		ID:              d.Id(),
		Name:            d.Get(AutomationActionFieldName).(string),
		Description:     d.Get(AutomationActionFieldDescription).(string),
		Tags:            d.Get(AutomationActionFieldTags),
		Type:            actionType,
		InputParameters: inputParameters,
		Fields:          fields,
	}, nil
}

func (r *AutomationActionResource) mapActionTypeFromSchema(d *schema.ResourceData) (string, error) {
	val, ok := d.GetOk(AutomationActionFieldScript)
	if ok && val != nil {
		return ActionTypeScript, nil
	}

	val, ok = d.GetOk(AutomationActionFieldHttp)
	if ok && val != nil {
		return ActionTypeHttp, nil
	}

	return "", errors.New("cannot determine the action type, invalid action configuration")
}

func (r *AutomationActionResource) mapInputParametersFromSchema(d *schema.ResourceData) ([]restapi.Parameter, error) {
	val, ok := d.GetOk(AutomationActionFieldInputParameter)

	if ok && val != nil {
		schemaParameters := val.([]interface{})
		result := make([]restapi.Parameter, len(schemaParameters))

		i := 0
		for _, v := range schemaParameters {
			param := v.(map[string]interface{})

			result[i] = restapi.Parameter{
				Name:        param[AutomationActionParameterFieldName].(string),
				Description: param[AutomationActionParameterFieldDescription].(string),
				Label:       param[AutomationActionParameterFieldLabel].(string),
				Required:    param[AutomationActionParameterFieldRequired].(bool),
				Hidden:      param[AutomationActionParameterFieldHidden].(bool),
				Type:        param[AutomationActionParameterFieldType].(string),
				Value:       param[AutomationActionParameterFieldValue].(string),
			}
			i++
		}
		return result, nil
	}

	return []restapi.Parameter{}, nil
}

func (r *AutomationActionResource) mapActionFieldsFromSchema(d *schema.ResourceData) ([]restapi.Field, error) {
	val, ok := d.GetOk(AutomationActionFieldScript)
	if ok && len(val.([]interface{})) == 1 {
		scriptData := val.([]interface{})[0].(map[string]interface{})
		return r.mapScriptFieldsFromSchema(scriptData)
	}

	val, ok = d.GetOk(AutomationActionFieldHttp)
	if ok && len(val.([]interface{})) == 1 {
		httpData := val.([]interface{})[0].(map[string]interface{})
		return r.mapHttpFieldsFromSchema(httpData)
	}

	return []restapi.Field{}, nil
}

func (r *AutomationActionResource) mapScriptFieldsFromSchema(scriptData map[string]interface{}) ([]restapi.Field, error) {
	result := make([]restapi.Field, len(scriptData))

	i := 0
	for k, v := range scriptData {
		var fieldName, fieldDescription, encoding string
		if k == AutomationActionFieldContent {
			fieldName = restapi.ScriptSshFieldName
			fieldDescription = restapi.ScriptSshFieldDescription
			encoding = Base64Encoding
		} else if k == AutomationActionFieldInterpreter {
			fieldName = restapi.SubtypeFieldName
			fieldDescription = restapi.SubtypeFieldDescription
			encoding = AsciiEncoding
		} else if k == AutomationActionFieldTimeout {
			fieldName = restapi.TimeoutFieldName
			fieldDescription = restapi.TimeoutFieldDescription
			encoding = AsciiEncoding
		}

		result[i] = restapi.Field{
			Name:        fieldName,
			Description: fieldDescription,
			Value:       v.(string),
			Encoding:    encoding,
			Secured:     false,
		}
		i++
	}
	return result, nil
}

func (r *AutomationActionResource) mapHttpFieldsFromSchema(httpData map[string]interface{}) ([]restapi.Field, error) {
	result := make([]restapi.Field, len(httpData))

	i := 0
	for k, v := range httpData {
		var fieldName, fieldDescription, fieldValue string
		if k == AutomationActionFieldHost {
			fieldName = restapi.HttpHostFieldName
			fieldDescription = restapi.HttpHostFieldDescription
			fieldValue = v.(string)
		} else if k == AutomationActionFieldMethod {
			fieldName = restapi.HttpMethodFieldName
			fieldDescription = restapi.HttpMethodFieldDescription
			fieldValue = v.(string)
		} else if k == AutomationActionFieldBody {
			fieldName = restapi.HttpBodyFieldName
			fieldDescription = restapi.HttpBodyFieldDescription
			fieldValue = v.(string)
		} else if k == AutomationActionFieldIgnoreCertErrors {
			fieldName = restapi.HttpIgnoreCertErrorsFieldName
			fieldDescription = restapi.HttpIgnoreCertErrorsFieldDescription
			fieldValue = strconv.FormatBool(v.(bool))
		} else if k == AutomationActionFieldTimeout {
			fieldName = restapi.TimeoutFieldName
			fieldDescription = restapi.TimeoutFieldDescription
			fieldValue = v.(string)
		} else if k == AutomationActionFieldHeaders {
			headers, err := r.mapHttpHeadersFromSchema(v.(map[string]interface{}))
			if err != nil {
				return nil, err
			}

			fieldName = restapi.HttpHeaderFieldName
			fieldDescription = restapi.HttpHeaderFieldDescription
			fieldValue = headers
		}

		result[i] = restapi.Field{
			Name:        fieldName,
			Description: fieldDescription,
			Value:       fieldValue,
			Encoding:    AsciiEncoding,
			Secured:     false,
		}
		i++
	}

	return result, nil
}

func (r *AutomationActionResource) mapHttpHeadersFromSchema(headers map[string]interface{}) (string, error) {
	headersJson, err := json.Marshal(headers)
	if err != nil {
		return "", fmt.Errorf("error marshaling HTTP headers. Caused by: %v", err)
	}

	return string(headersJson), nil
}

func (r *AutomationActionResource) mapHttpHeadersToSchema(headers string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(headers), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling HTTP headers. Caused by: %v", err)
	}
	return result, nil
}
