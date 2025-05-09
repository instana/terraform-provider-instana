package restapi

// AutomationActionResourcePath path to Automation Actions resource of Instana RESTful API
const AutomationActionResourcePath = AutomationBasePath + "/actions"

// AutomationAction is the representation of an automation action in Instana
type AutomationAction struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Type            string      `json:"type"`
	Tags            interface{} `json:"tags"`
	Timeout         int         `json:"timeout"`
	Fields          []Field     `json:"fields"`
	InputParameters []Parameter `json:"inputParameters"`
}

type Parameter struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Required    bool   `json:"required"`
	Hidden      bool   `json:"hidden"`
	Secured     bool   `json:"secured"`
}

type Field struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Encoding    string `json:"encoding"`
	Value       string `json:"value"`
	Secured     bool   `json:"secured"`
}

// constants for field names and descriptions
const (
	SUBTYPE_FIELD_NAME        = "subtype"
	SUBTYPE_FIELD_DESCRIPTION = "script subtype"

	SCRIPT_SSH_FIELD_NAME        = "script_ssh"
	SCRIPT_SSH_FIELD_DESCRIPTION = "script content"

	TIMEOUT_FIELD_NAME        = "timeout"
	TIMEOUT_FIELD_DESCRIPTION = "timeout of the action execution in seconds"

	HTTP_HOST_FIELD_NAME        = "host"
	HTTP_HOST_FIELD_DESCRIPTION = "url of the https request"

	HTTP_BODY_FIELD_NAME        = "body"
	HTTP_BODY_FIELD_DESCRIPTION = "body of the https request"

	HTTP_METHOD_FIELD_NAME        = "method"
	HTTP_METHOD_FIELD_DESCRIPTION = "HTTP method"

	HTTP_HEADER_FIELD_NAME        = "header"
	HTTP_HEADER_FIELD_DESCRIPTION = "header of the https request"

	HTTP_IGNORE_CERT_ERRORS_FIELD_NAME        = "ignoreCertErrors"
	HTTP_IGNORE_CERT_ERRORS_FIELD_DESCRIPTION = "ignore certificate errors for request"
)

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *AutomationAction) GetIDForResourcePath() string {
	return spec.ID
}
