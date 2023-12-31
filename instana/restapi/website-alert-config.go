package restapi

// WebsiteAlertConfigResourcePath path to website alert config resource of Instana RESTful API
const WebsiteAlertConfigResourcePath = EventSettingsBasePath + "/website-alert-configs"

// WebsiteAlertConfig is the representation of an website alert configuration in Instana
type WebsiteAlertConfig struct {
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
	Description           string                    `json:"description"`
	Severity              int                       `json:"severity"`
	Triggering            bool                      `json:"triggering"`
	WebsiteID             string                    `json:"websiteId"`
	TagFilterExpression   *TagFilter                `json:"tagFilterExpression"`
	AlertChannelIDs       []string                  `json:"alertChannelIds"`
	Granularity           Granularity               `json:"granularity"`
	CustomerPayloadFields []CustomPayloadField[any] `json:"customPayloadFields"`
	Rule                  WebsiteAlertRule          `json:"rule"`
	Threshold             Threshold                 `json:"threshold"`
	TimeThreshold         WebsiteTimeThreshold      `json:"timeThreshold"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (r *WebsiteAlertConfig) GetIDForResourcePath() string {
	return r.ID
}

// GetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (a *WebsiteAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return a.CustomerPayloadFields
}

// SetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (a *WebsiteAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	a.CustomerPayloadFields = fields
}
