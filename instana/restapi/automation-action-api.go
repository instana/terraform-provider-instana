package restapi

// AutomationActionResourcePath path to Automation Actions resource of Instana RESTful API
const AutomationActionResourcePath = AutomationBasePath + "/actions"

// AutomationAction is the representation of a automation action in Instana
type AutomationAction struct {
	ID            string  	`json:"id"`
	Name          string  	`json:"name"`
	Description   *string 	`json:"description"`
	Type      	  string  	`json:"type"`
	Tags		  []string	`json:"tags"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *AutomationAction) GetIDForResourcePath() string {
	return spec.ID
}