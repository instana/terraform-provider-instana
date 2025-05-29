package restapi

// AutomationPolicyResourcePath path to Automation Policies resource of Instana RESTful API
const AutomationPolicyResourcePath = AutomationBasePath + "/policies"

// AutomationPolicy is the representation of an automation policy in Instana
type AutomationPolicy struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Tags               interface{}         `json:"tags"`
	Trigger            Trigger             `json:"trigger"`
	TypeConfigurations []TypeConfiguration `json:"typeConfigurations"`
}

type Trigger struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type TypeConfiguration struct {
	Name      string     `json:"name"`
	Condition *Condition `json:"condition"`
	Runnable  Runnable   `json:"runnable"`
}

type Condition struct {
	Query string `json:"query"`
}

type Runnable struct {
	Id               string           `json:"id"`
	Type             string           `json:"type"`
	RunConfiguration RunConfiguration `json:"runConfiguration"`
}

type RunConfiguration struct {
	Actions []ActionConfiguration `json:"actions"`
}

type ActionConfiguration struct {
	Action               Action                `json:"action"`
	AgentId              string                `json:"agentId"`
	InputParameterValues []InputParameterValue `json:"inputParameterValues"`
}

type Action struct {
	Id string `json:"id"`
}

type InputParameterValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *AutomationPolicy) GetIDForResourcePath() string {
	return spec.ID
}
