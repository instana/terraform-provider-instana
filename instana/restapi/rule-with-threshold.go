package restapi

type RuleWithThreshold[R InfraAlertRule | ApplicationAlertRule | WebsiteAlertRule] struct {
	ThresholdOperator ThresholdOperator               `json:"thresholdOperator"`
	Rule              R                               `json:"rule"`
	Thresholds        map[AlertSeverity]ThresholdRule `json:"thresholds"`
}

type AlertSeverity string
type AlertSeverities []AlertSeverity

const (
	WarningSeverity  = AlertSeverity("WARNING")
	CriticalSeverity = AlertSeverity("CRITICAL")
)

// SupportedAlertSeverities : will be used as part of validation in a follow-up.
var SupportedAlertSeverities = AlertSeverities{WarningSeverity, CriticalSeverity}
