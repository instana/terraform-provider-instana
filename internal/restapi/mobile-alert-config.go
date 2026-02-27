package restapi

// MobileAlertConfigResourcePath path to mobile alert config resource of Instana RESTful API
const MobileAlertConfigResourcePath = EventSettingsBasePath + "/mobile-app-alert-configs"

// MobileAlertConfig is the representation of a mobile app alert configuration in Instana
type MobileAlertConfig struct {
	ID                  string                             `json:"id"`
	Name                string                             `json:"name"`
	Description         string                             `json:"description"`
	MobileAppID         string                             `json:"mobileAppId"`
	Triggering          bool                               `json:"triggering"`
	Enabled             *bool                              `json:"enabled,omitempty"`
	TagFilterExpression *TagFilter                         `json:"tagFilterExpression"`
	AlertChannels       map[string][]string                `json:"alertChannels,omitempty"`
	Granularity         Granularity                        `json:"granularity"`
	GracePeriod         *int64                             `json:"gracePeriod,omitempty"`
	CustomPayloadFields []CustomPayloadField[any]          `json:"customPayloadFields"`
	Rules               []MobileAppAlertRuleWithThresholds `json:"rules,omitempty"`
	TimeThreshold       *MobileAppTimeThreshold            `json:"timeThreshold"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (m *MobileAlertConfig) GetIDForResourcePath() string {
	return m.ID
}

// GetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (m *MobileAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return m.CustomPayloadFields
}

// SetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (m *MobileAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	m.CustomPayloadFields = fields
}

// MobileAppAlertRule represents a mobile app alert rule
type MobileAppAlertRule struct {
	AlertType       string       `json:"alertType"`
	MetricName      string       `json:"metricName"`
	Aggregation     *Aggregation `json:"aggregation,omitempty"`
	Operator        *string      `json:"operator,omitempty"`
	Value           *string      `json:"value,omitempty"`
	CustomEventName *string      `json:"customEventName,omitempty"`
}

// MobileAppAlertRuleWithThresholds represents a rule with multiple thresholds and severity levels
type MobileAppAlertRuleWithThresholds struct {
	Rule              *MobileAppAlertRule             `json:"rule"`
	ThresholdOperator string                          `json:"thresholdOperator"`
	Thresholds        map[AlertSeverity]ThresholdRule `json:"thresholds"`
}

// MobileAppTimeThreshold represents the time threshold configuration for mobile app alerts
type MobileAppTimeThreshold struct {
	Type                    string   `json:"type"`
	TimeWindow              *int64   `json:"timeWindow,omitempty"`
	Violations              *int32   `json:"violations,omitempty"`
	Users                   *int32   `json:"users,omitempty"`
	UserPercentage          *float64 `json:"userPercentage,omitempty"`
	ImpactMeasurementMethod string   `json:"impactMeasurementMethod,omitempty"`
}
