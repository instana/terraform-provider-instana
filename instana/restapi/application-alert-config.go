package restapi

// ApplicationAlertConfigsResourcePath the base path of the Instana REST API for application alert configs
const ApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/application-alert-configs"

// GlobalApplicationAlertConfigsResourcePath the base path of the Instana REST API for global application alert configs
const GlobalApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/global-alert-configs/applications"

// ApplicationAlertConfig is the representation of an application alert configuration in Instana
type ApplicationAlertConfig struct {
	ID                    string                               `json:"id"`
	Name                  string                               `json:"name"`
	Description           string                               `json:"description"`
	Severity              int                                  `json:"severity"`
	Triggering            bool                                 `json:"triggering"`
	Applications          map[string]IncludedApplication       `json:"applications"`
	BoundaryScope         BoundaryScope                        `json:"boundaryScope"`
	TagFilterExpression   *TagFilter                           `json:"tagFilterExpression"`
	IncludeInternal       bool                                 `json:"includeInternal"`
	IncludeSynthetic      bool                                 `json:"includeSynthetic"`
	EvaluationType        ApplicationAlertEvaluationType       `json:"evaluationType"`
	AlertChannelIDs       []string                             `json:"alertChannelIds"`
	AlertChannels         map[string][]string                  `json:"alertChannels"`
	Granularity           Granularity                          `json:"granularity"`
	GracePeriod           *int64                               `json:"gracePeriod"`
	CustomerPayloadFields []CustomPayloadField[any]            `json:"customPayloadFields"`
	Rule                  *ApplicationAlertRule                `json:"rule"`
	Rules                 []ApplicationAlertRuleWithThresholds `json:"rules"`
	Threshold             *Threshold                           `json:"threshold"`
	TimeThreshold         *ApplicationAlertTimeThreshold       `json:"timeThreshold"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (a *ApplicationAlertConfig) GetIDForResourcePath() string {
	return a.ID
}

// GetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (a *ApplicationAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return a.CustomerPayloadFields
}

// SetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (a *ApplicationAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	a.CustomerPayloadFields = fields
}

// ApplicationScope represents an application in the application alert config
type ApplicationScope struct {
	ApplicationID string         `json:"applicationId"`
	Inclusive     bool           `json:"inclusive"`
	Services      []ServiceScope `json:"services,omitempty"`
}

// ServiceScope represents a service in the application alert config
type ServiceScope struct {
	ServiceID string          `json:"serviceId"`
	Inclusive bool            `json:"inclusive"`
	Endpoints []EndpointScope `json:"endpoints,omitempty"`
}

// EndpointScope represents an endpoint in the application alert config
type EndpointScope struct {
	EndpointID string `json:"endpointId"`
	Inclusive  bool   `json:"inclusive"`
}

// ApplicationAlertTimeThreshold represents the time threshold configuration for application alerts
type ApplicationAlertTimeThreshold struct {
	Type                 string                                             `json:"type"`
	TimeWindow           int64                                              `json:"timeWindow"`
	RequestImpact        *ApplicationAlertTimeThresholdRequestImpact        `json:"-"`
	ViolationsInPeriod   *ApplicationAlertTimeThresholdViolationsInPeriod   `json:"-"`
	ViolationsInSequence *ApplicationAlertTimeThresholdViolationsInSequence `json:"-"`
}

// ApplicationAlertTimeThresholdRequestImpact represents the request impact time threshold configuration
type ApplicationAlertTimeThresholdRequestImpact struct {
	TimeWindow int `json:"timeWindow"`
	Requests   int `json:"requests"`
}

// ApplicationAlertTimeThresholdViolationsInPeriod represents the violations in period time threshold configuration
type ApplicationAlertTimeThresholdViolationsInPeriod struct {
	TimeWindow int `json:"timeWindow"`
	Violations int `json:"violations"`
}

// ApplicationAlertTimeThresholdViolationsInSequence represents the violations in sequence time threshold configuration
type ApplicationAlertTimeThresholdViolationsInSequence struct {
	TimeWindow int `json:"timeWindow"`
}

// ApplicationAlertRuleErrorRate represents an error rate rule
type ApplicationAlertRuleErrorRate struct {
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation,omitempty"`
}

// ApplicationAlertRuleErrors represents an errors rule
type ApplicationAlertRuleErrors struct {
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation,omitempty"`
}

// ApplicationAlertRuleLogs represents a logs rule
type ApplicationAlertRuleLogs struct {
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation,omitempty"`
	Level       string `json:"level"`
	Message     string `json:"message,omitempty"`
	Operator    string `json:"operator"`
}

// ApplicationAlertRuleSlowness represents a slowness rule
type ApplicationAlertRuleSlowness struct {
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation"`
}

// ApplicationAlertRuleStatusCode represents a status code rule
type ApplicationAlertRuleStatusCode struct {
	MetricName      string `json:"metricName"`
	Aggregation     string `json:"aggregation,omitempty"`
	StatusCodeStart int    `json:"statusCodeStart,omitempty"`
	StatusCodeEnd   int    `json:"statusCodeEnd,omitempty"`
}

// ApplicationAlertRuleThroughput represents a throughput rule
type ApplicationAlertRuleThroughput struct {
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation,omitempty"`
}
