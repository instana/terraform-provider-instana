package restapi

const (
	// LogAlertConfigResourcePath path to Log Alert Config resource of Instana RESTful API
	LogAlertConfigResourcePath = EventSettingsBasePath + "/global-alert-configs/logs"
)

// LogAlertRule represents the rule configuration for log alerts
type LogAlertRule struct {
	AlertType   string      `json:"alertType"`
	MetricName  string      `json:"metricName"`
	Aggregation Aggregation `json:"aggregation,omitempty"`
}

// LogTimeThreshold represents the time threshold configuration for log alerts
type LogTimeThreshold struct {
	Type       string `json:"type"`
	TimeWindow int64  `json:"timeWindow"`
}

// GroupByTag represents a tag used for grouping in log alerts
type GroupByTag struct {
	TagName string `json:"tagName"`
	Key     string `json:"key,omitempty"`
}

// LogAlertConfig represents the Instana API model for log alert configurations
type LogAlertConfig struct {
	ID                    string                            `json:"id,omitempty"`
	Name                  string                            `json:"name"`
	Description           string                            `json:"description"`
	TagFilterExpression   *TagFilter                        `json:"tagFilterExpression"`
	AlertChannels         map[AlertSeverity][]string        `json:"alertChannels,omitempty"`
	AlertChannelIds       []string                          `json:"alertChannelIds,omitempty"`
	Severity              int32                             `json:"severity,omitempty"`
	Granularity           Granularity                       `json:"granularity"`
	TimeThreshold         *LogTimeThreshold                 `json:"timeThreshold"`
	Threshold             *Threshold                        `json:"threshold,omitempty"`
	GracePeriod           int64                             `json:"gracePeriod,omitempty"`
	CustomerPayloadFields []CustomPayloadField[any]         `json:"customPayloadFields,omitempty"`
	Rules                 []RuleWithThreshold[LogAlertRule] `json:"rules"`
	GroupBy               []GroupByTag                      `json:"groupBy,omitempty"`
	Enabled               bool                              `json:"enabled,omitempty"`
	Created               int64                             `json:"created,omitempty"`
	ReadOnly              bool                              `json:"readOnly,omitempty"`
}

// GetIDForResourcePath implementation of the InstanaDataObject interface
func (r *LogAlertConfig) GetIDForResourcePath() string {
	return r.ID
}

// GetCustomerPayloadFields implementation of the customPayloadFieldsAwareInstanaDataObject interface
func (r *LogAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return r.CustomerPayloadFields
}

// SetCustomerPayloadFields implementation of the customPayloadFieldsAwareInstanaDataObject interface
func (r *LogAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	r.CustomerPayloadFields = fields
}

// Made with Bob
