package restapi

const InfraAlertConfigResourcePath = EventSettingsBasePath + "/infra-alert-configs"

type InfraAlertConfig struct {
	ID                    string                              `json:"id"`
	Name                  string                              `json:"name"`
	Description           string                              `json:"description"`
	TagFilterExpression   *TagFilter                          `json:"tagFilterExpression"`
	GroupBy               []string                            `json:"groupBy"`
	Granularity           Granularity                         `json:"granularity"`
	TimeThreshold         *InfraTimeThreshold                 `json:"timeThreshold"`
	CustomerPayloadFields []CustomPayloadField[any]           `json:"customPayloadFields"`
	Rules                 []RuleWithThreshold[InfraAlertRule] `json:"rules"`
	AlertChannels         map[AlertSeverity][]string          `json:"alertChannels"`
	EvaluationType        InfraAlertEvaluationType            `json:"evaluationType"`
}

func (config *InfraAlertConfig) GetIDForResourcePath() string {
	return config.ID
}

func (config *InfraAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return config.CustomerPayloadFields
}

func (config *InfraAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	config.CustomerPayloadFields = fields
}
