package restapi

// MaintenanceWindowConfigResourcePath path to maintenance window config resource of Instana RESTful API
const MaintenanceWindowConfigResourcePath = SettingsBasePath + "/v2/maintenance"

// MaintenanceWindowConfig data structure of a Maintenance Window Configuration of the Instana API
type MaintenanceWindowConfig struct {
	ID                         string                 `json:"id"`
	Name                       string                 `json:"name"`
	Query                      string                 `json:"query"`
	Scheduling                 *MaintenanceScheduling `json:"scheduling"`
	TagFilterExpressionEnabled *bool                  `json:"tagFilterExpressionEnabled,omitempty"`
	TagFilterExpression        *TagFilter             `json:"tagFilterExpression,omitempty"`
	Paused                     *bool                  `json:"paused,omitempty"`
	RetriggerOpenAlertsEnabled *bool                  `json:"retriggerOpenAlertsEnabled,omitempty"`
	ValidVersion               *int                   `json:"validVersion,omitempty"`
	LastUpdated                *int64                 `json:"lastUpdated,omitempty"`
	State                      *string                `json:"state,omitempty"`
	Occurrence                 *MaintenanceOccurrence `json:"occurrence,omitempty"`
}

// MaintenanceScheduling represents the scheduling configuration for maintenance windows
type MaintenanceScheduling struct {
	Start      int64                `json:"start"`
	Duration   *MaintenanceDuration `json:"duration"`
	Type       string               `json:"type"` // ONE_TIME or RECURRENT
	Rrule      *string              `json:"rrule,omitempty"`
	TimezoneId *string              `json:"timezoneId,omitempty"`
}

// MaintenanceDuration represents the duration of a maintenance window
type MaintenanceDuration struct {
	Amount int64  `json:"amount"`
	Unit   string `json:"unit"` // MINUTES, HOURS, DAYS
}

// MaintenanceOccurrence represents the occurrence details of a maintenance window
type MaintenanceOccurrence struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (m *MaintenanceWindowConfig) GetIDForResourcePath() string {
	return m.ID
}

// Made with Bob
