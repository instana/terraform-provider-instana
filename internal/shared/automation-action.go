package shared

const (
	AutomationActionFieldName           = "name"
	AutomationActionFieldDescription    = "description"
	AutomationActionFieldTags           = "tags"
	AutomationActionFieldTimeout        = "timeout"
	AutomationActionFieldType           = "type"
	AutomationActionFieldInputParameter = "input_parameter"

	// script constants
	AutomationActionFieldScript      = "script"
	AutomationActionFieldContent     = "content"
	AutomationActionFieldInterpreter = "interpreter"

	// http constants
	AutomationActionFieldHttp             = "http"
	AutomationActionFieldMethod           = "method"
	AutomationActionFieldHost             = "host"
	AutomationActionFieldHeaders          = "headers"
	AutomationActionFieldBody             = "body"
	AutomationActionFieldIgnoreCertErrors = "ignore_certificate_errors"

	// input parameter constants
	AutomationActionParameterFieldName        = "name"
	AutomationActionParameterFieldLabel       = "label"
	AutomationActionParameterFieldDescription = "description"
	AutomationActionParameterFieldType        = "type"
	AutomationActionParameterFieldValue       = "value"
	AutomationActionParameterFieldRequired    = "required"
	AutomationActionParameterFieldHidden      = "hidden"
)
