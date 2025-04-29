package restapi

// AutomationActionResourcePath path to Automation Actions resource of Instana RESTful API
const AutomationActionResourcePath = AutomationBasePath + "/actions"

// AutomationAction is the representation of an automation action in Instana
type AutomationAction struct {
	ID            	string  		`json:"id"`
	Name          	string  		`json:"name"`
	Description   	string 			`json:"description"`
	Type      	  	string  		`json:"type"`
	Tags       	  	interface{} 	`json:"tags"`
	Timeout		  	int       		`json:"timeout"`
	InputParameters []Parameter		`json:"inputParameters"`
}

type Parameter struct {
	Name          string  		`json:"name"`
	Label		  string		`json:"label"`
	Description   string 		`json:"description"`
	Type      	  string  		`json:"type"`
	Value 		  string		`json:"value"`
	Required	  bool			`json:"required"`
	Hidden		  bool			`json:"hidden"`
	Secured       bool			`json:"secured"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *AutomationAction) GetIDForResourcePath() string {
	return spec.ID
}