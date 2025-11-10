package restapi

// ApplicationAlertRuleWithThresholds represents a rule with multiple thresholds and severity levels
type ApplicationAlertRuleWithThresholds struct {
	Rule              *ApplicationAlertRule           `json:"rule"`
	ThresholdOperator string                          `json:"thresholdOperator"`
	Thresholds        map[AlertSeverity]ThresholdRule `json:"thresholds"`
}

// ThresholdValue represents a threshold value for a specific severity level
type ThresholdValue struct {
	Value float64 `json:"value"`
}
