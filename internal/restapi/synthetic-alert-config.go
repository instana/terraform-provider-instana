package restapi

// SyntheticAlertConfigsResourcePath is the path to the synthetic alert configs resource of the Instana API
const SyntheticAlertConfigsResourcePath = "/api/events/settings/global-alert-configs/synthetics"

// SyntheticAlertConfig represents the API model for synthetic alert configurations
type SyntheticAlertConfig struct {
	ID                    string                      `json:"id,omitempty"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description,omitempty"`
	SyntheticTestIds      []string                    `json:"syntheticTestIds"`
	Severity              int                         `json:"severity"`
	TagFilterExpression   *TagFilter                  `json:"tagFilterExpression,omitempty"`
	Rule                  SyntheticAlertRule          `json:"rule"`
	AlertChannelIds       []string                    `json:"alertChannelIds"`
	TimeThreshold         SyntheticAlertTimeThreshold `json:"timeThreshold"`
	CustomerPayloadFields []CustomPayloadField[any]   `json:"customPayloadFields,omitempty"`
	GracePeriod           int64                       `json:"gracePeriod,omitempty"`
}

// SyntheticAlertRule represents the rule configuration for synthetic alerts
type SyntheticAlertRule struct {
	AlertType   string `json:"alertType"`
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation"`
}

// SyntheticAlertTimeThreshold represents the time threshold configuration for synthetic alerts
type SyntheticAlertTimeThreshold struct {
	Type            string `json:"type"`
	ViolationsCount int    `json:"violationsCount"`
}

// GetIDForResourcePath returns the ID to be used in the resource path when calling the API
func (r *SyntheticAlertConfig) GetIDForResourcePath() string {
	return r.ID
}

// GetCustomerPayloadFields implements the CustomPayloadFieldsAware interface
func (r *SyntheticAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return r.CustomerPayloadFields
}

// SetCustomerPayloadFields implements the CustomPayloadFieldsAware interface
func (r *SyntheticAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	r.CustomerPayloadFields = fields
}
