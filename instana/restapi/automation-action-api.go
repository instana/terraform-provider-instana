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
	SubtypeFieldName        = "subtype"
	SubtypeFieldDescription = "script subtype"

	ScriptSshFieldName        = "script_ssh"
	ScriptSshFieldDescription = "script content"

	TimeoutFieldName        = "timeout"
	TimeoutFieldDescription = "timeout of the action execution in seconds"

	HttpHostFieldName        = "host"
	HttpHostFieldDescription = "url of the https request"

	HttpBodyFieldName        = "body"
	HttpBodyFieldDescription = "body of the https request"

	HttpMethodFieldName        = "method"
	HttpMethodFieldDescription = "HTTP method"

	HttpHeaderFieldName        = "header"
	HttpHeaderFieldDescription = "header of the https request"

	HttpIgnoreCertErrorsFieldName        = "ignoreCertErrors"
	HttpIgnoreCertErrorsFieldDescription = "ignore certificate errors for request"
)

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *AutomationAction) GetIDForResourcePath() string {
	return spec.ID
}
