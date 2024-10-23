package restapi

type InfraAlertRule struct {
	AlertType              string       `json:"alertType"`
	MetricName             string       `json:"metricName"`
	EntityType             string       `json:"entityType"`
	Aggregation            *Aggregation `json:"aggregation"`
	CrossSeriesAggregation *Aggregation `json:"crossSeriesAggregation"`
	Regex                  bool         `json:"regex"`
}
