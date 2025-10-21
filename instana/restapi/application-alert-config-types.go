package restapi

// ApplicationAlertRuleWithThresholds represents a rule with multiple thresholds and severity levels
type ApplicationAlertRuleWithThresholds struct {
	Rule              *ApplicationAlertRule     `json:"rule"`
	ThresholdOperator string                    `json:"thresholdOperator"`
	Thresholds        map[string]ThresholdValue `json:"thresholds"`
}

// ThresholdValue represents a threshold value for a specific severity level
type ThresholdValue struct {
	Value float64 `json:"value"`
}

// Made with Bob
