package restapi

// ApplicationAlertRuleWithThresholds represents a rule with multiple thresholds and severity levels
type ApplicationAlertRuleWithThresholds struct {
	Rule              *ApplicationAlertRule           `json:"rule"`
	ThresholdOperator string                          `json:"thresholdOperator"`
	Thresholds        map[AlertSeverity]ThresholdRule `json:"thresholds"`
}

// Made with Bob
