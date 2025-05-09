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
const ACTION_TYPE_SCRIPT = "SCRIPT"
const ACTION_TYPE_HTTP = "HTTP"

// encodings
const ASCII_ENCODING = "ascii"
const BASE64_ENCODING = "base64"

// NewAutomationActionResourceHandle creates the resource handle for Automation Actions
func NewAutomationActionResourceHandle() ResourceHandle[*restapi.AutomationAction] {
	return &AutomationActionResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAutomationAction,
			Schema: map[string]*schema.Schema{
				AutomationActionFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the automation action",
				},
				AutomationActionFieldDescription: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The description of the automation action",
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
		item[AutomationActionParameterFieldHidden] = val.Hidden
		item[AutomationActionParameterFieldSecured] = val.Secured
		item[AutomationActionParameterFieldType] = val.Type
		item[AutomationActionParameterFieldValue] = val.Value

		result[i] = item
		i++
	}
	return result, nil
}

func (r *AutomationActionResource) mapScriptFieldsToSchema(action *restapi.AutomationAction) ([]interface{}, error) {
	if action.Type == ACTION_TYPE_SCRIPT {
		result := make(map[string]interface{})

		// script content is required field
		result[AutomationActionFieldContent] = r.getFieldValue(action, restapi.SCRIPT_SSH_FIELD_NAME)
		// interpreter is optional field
		interpreter := r.getFieldValue(action, restapi.SUBTYPE_FIELD_NAME)
		if len(interpreter) > 0 {
			result[AutomationActionFieldInterpreter] = interpreter
		}
		// timeout is optional field
		timeout := r.getFieldValue(action, restapi.TIMEOUT_FIELD_NAME)
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
	if action.Type == ACTION_TYPE_HTTP {
		result := make(map[string]interface{})

		// required fields
		result[AutomationActionFieldHost] = r.getFieldValue(action, restapi.HTTP_HOST_FIELD_NAME)
		result[AutomationActionFieldMethod] = r.getFieldValue(action, restapi.HTTP_METHOD_FIELD_NAME)

		// optional fields
		body := r.getFieldValue(action, restapi.HTTP_BODY_FIELD_NAME)
		if len(body) > 0 {
			result[AutomationActionFieldBody] = body
		}
		ignoreCertErrors := r.getBoolFieldValueOrDefault(action, restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_NAME, false)
		result[AutomationActionFieldIgnoreCertErrors] = ignoreCertErrors

		headersData := r.getFieldValue(action, restapi.HTTP_HEADER_FIELD_NAME)
		headers, err := r.mapHttpHeadersToSchema(headersData)
		if err != nil {
			return nil, err
		}
		result[AutomationActionFieldHeaders] = headers

		timeout := r.getFieldValue(action, restapi.TIMEOUT_FIELD_NAME)
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
		return ACTION_TYPE_SCRIPT, nil
	}

	val, ok = d.GetOk(AutomationActionFieldHttp)
	if ok && val != nil {
		return ACTION_TYPE_HTTP, nil
	}

	return "", errors.New("cannot determine the action type, invalid action configuration")
}

func (r *AutomationActionResource) mapInputParametersFromSchema(d *schema.ResourceData) ([]restapi.Parameter, error) {
	val, ok := d.GetOk(AutomationActionFieldInputParameter)

	if ok && val != nil {
		schemaParameters := val.(*schema.Set).List()
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
				Secured:     param[AutomationActionParameterFieldSecured].(bool),
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
			fieldName = restapi.SCRIPT_SSH_FIELD_NAME
			fieldDescription = restapi.SCRIPT_SSH_FIELD_DESCRIPTION
			encoding = BASE64_ENCODING
		} else if k == AutomationActionFieldInterpreter {
			fieldName = restapi.SUBTYPE_FIELD_NAME
			fieldDescription = restapi.SUBTYPE_FIELD_DESCRIPTION
			encoding = ASCII_ENCODING
		} else if k == AutomationActionFieldTimeout {
			fieldName = restapi.TIMEOUT_FIELD_NAME
			fieldDescription = restapi.TIMEOUT_FIELD_DESCRIPTION
			encoding = ASCII_ENCODING
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
			fieldName = restapi.HTTP_HOST_FIELD_NAME
			fieldDescription = restapi.HTTP_HOST_FIELD_DESCRIPTION
			fieldValue = v.(string)
		} else if k == AutomationActionFieldMethod {
			fieldName = restapi.HTTP_METHOD_FIELD_NAME
			fieldDescription = restapi.HTTP_METHOD_FIELD_DESCRIPTION
			fieldValue = v.(string)
		} else if k == AutomationActionFieldBody {
			fieldName = restapi.HTTP_BODY_FIELD_NAME
			fieldDescription = restapi.HTTP_BODY_FIELD_DESCRIPTION
			fieldValue = v.(string)
		} else if k == AutomationActionFieldIgnoreCertErrors {
			fieldName = restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_NAME
			fieldDescription = restapi.HTTP_IGNORE_CERT_ERRORS_FIELD_DESCRIPTION
			fieldValue = strconv.FormatBool(v.(bool))
		} else if k == AutomationActionFieldTimeout {
			fieldName = restapi.TIMEOUT_FIELD_NAME
			fieldDescription = restapi.TIMEOUT_FIELD_DESCRIPTION
			fieldValue = v.(string)
		} else if k == AutomationActionFieldHeaders {
			headers, err := r.mapHttpHeadersFromSchema(v.(map[string]interface{}))
			if err != nil {
				return nil, err
			}

			fieldName = restapi.HTTP_HEADER_FIELD_NAME
			fieldDescription = restapi.HTTP_HEADER_FIELD_DESCRIPTION
			fieldValue = headers
		}

		result[i] = restapi.Field{
			Name:        fieldName,
			Description: fieldDescription,
			Value:       fieldValue,
			Encoding:    ASCII_ENCODING,
			Secured:     false,
		}
		i++
	}

	fmt.Printf(">>>>>>>> HTTP fields are: %v\n", result)
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
