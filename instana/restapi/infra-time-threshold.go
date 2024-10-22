package restapi

type InfraTimeThreshold struct {
	Type       string `json:"type"`
	TimeWindow *int64 `json:"timeWindow"`
}
