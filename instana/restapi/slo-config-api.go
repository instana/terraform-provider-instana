package restapi

import (
	"fmt"
	"os"
)

const (
	//SloConfigResourcePath path to sli config resource of Instana RESTful API
	SloConfigResourcePath = SettingsBasePath + "/slo"
)

// SloConfig represents the REST resource of slo configuration at Instana
type SloConfig struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Target     float64       `json:"target"`
	Tags       []string      `json:"tags"`
	Entity     SloEntity     `json:"entity"`
	Indicator  SloIndicator  `json:"indicator"`
	TimeWindow SloTimeWindow `json:"timeWindow"`
	RbacTags   []RbacTag     `json:"rbacTags,omitempty"`

	// CreatedDate int         `json:"createdDate"`
	// LastUpdated int         `json:"lastUpdated"`
}

// RbacTag represents a RBAC tag in the SLO configuration
type RbacTag struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

type SloEntity struct {
	Type             string        `json:"type"`
	ApplicationID    *string       `json:"applicationId"`
	ServiceID        *string       `json:"serviceId"`
	EndpointID       *string       `json:"endpointId"`
	BoundaryScope    *string       `json:"boundaryScope"`
	IncludeSynthetic *bool         `json:"includeSynthetic"`
	IncludeInternal  *bool         `json:"includeInternal"`
	FilterExpression *TagFilter    `json:"tagFilterExpression"`
	WebsiteId        *string       `json:"websiteId"`
	BeaconType       *string       `json:"beaconType"`
	SyntheticTestIDs []interface{} `json:"syntheticTestIds"`
	InfraType        *string       `json:"infraType"`
}

// SloEntity represents the nested object sli entity of the sli config REST resource at Instana
type SloApplicationEntity struct {
	Type             string     `json:"type"`
	ApplicationID    *string    `json:"applicationId"`
	ServiceID        *string    `json:"serviceId"`
	EndpointID       *string    `json:"endpointId"`
	BoundaryScope    *string    `json:"boundaryScope"`
	IncludeSynthetic *bool      `json:"includeSynthetic"`
	IncludeInternal  *bool      `json:"includeInternal"`
	FilterExpression *TagFilter `json:"tagFilterExpression"`
}

type SloWebsiteEntity struct {
	Type             string     `json:"type"`
	WebsiteId        *string    `json:"websiteId"`
	BeaconType       *string    `json:"beaconType"`
	FilterExpression *TagFilter `json:"tagFilterExpression"`
}

type SloInfraEntity struct {
	Type      string `json:"type"`
	InfraType string `json:"infraType"`
}

type SloSyntheticEntity struct {
	Type             string        `json:"type"`
	SyntheticTestIDs []interface{} `json:"syntheticTestIds"`
	FilterExpression *TagFilter    `json:"tagFilterExpression"`
}

type SloIndicator struct {
	Blueprint                 string     `json:"blueprint"`
	Type                      string     `json:"type"`
	Threshold                 float64    `json:"threshold"`
	Aggregation               *string    `json:"aggregation"`
	Operator                  *string    `json:"operator"`
	TrafficType               *string    `json:"trafficType"`
	GoodEventFilterExpression *TagFilter `json:"goodEventsFilter"`
	BadEventFilterExpression  *TagFilter `json:"badEventsFilter"`
}

// Blueprints
type SloTimeBasedLatencyIndicator struct {
	Blueprint   string  `json:"blueprint"`
	Type        string  `json:"type"`
	Threshold   float64 `json:"threshold"`
	Aggregation string  `json:"aggregation"`
}

type SloTimeBasedAvailabilityIndicator struct {
	Blueprint   string  `json:"blueprint"`
	Type        string  `json:"type"`
	Threshold   float64 `json:"threshold"`
	Aggregation string  `json:"aggregation"`
}

type SloTrafficIndicator struct {
	Blueprint   string  `json:"blueprint"`
	TrafficType string  `json:"trafficType"`
	Threshold   float64 `json:"threshold"`
	Aggregation string  `json:"aggregation"`
}

type SloEventBasedLatencyIndicator struct {
	Blueprint string  `json:"blueprint"`
	Type      string  `json:"type"`
	Threshold float64 `json:"threshold"`
}

type SloEventBasedAvailabilityIndicator struct {
	Blueprint string `json:"blueprint"`
	Type      string `json:"type"`
}

type SloCustomIndicator struct {
	Type                      string     `json:"type"`
	Blueprint                 string     `json:"blueprint"`
	GoodEventFilterExpression *TagFilter `json:"goodEventsFilter"`
	BadEventFilterExpression  *TagFilter `json:"badEventsFilter"`
}

type SloTimeWindow struct {
	Type         string  `json:"type"`
	Duration     int     `json:"duration"`
	DurationUnit string  `json:"durationUnit"`
	Timezone     string  `json:"timezone,omitempty"`
	StartTime    float64 `json:"startTimestamp"`
}

// time windows
type SloRollingTimeWindow struct {
	Type         string `json:"type"`
	Duration     int    `json:"duration"`
	DurationUnit string `json:"durationUnit"`
	Timezone     string `json:"timezone,omitempty"`
}

type SloFixedTimeWindow struct {
	Type         string  `json:"type"`
	Duration     int     `json:"duration"`
	DurationUnit string  `json:"durationUnit"`
	Timezone     string  `json:"timezone,omitempty"`
	StartTime    float64 `json:"startTimestamp"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (s *SloConfig) GetIDForResourcePath() string {
	fmt.Fprintln(os.Stderr, ">> GetIDForResourcePath: "+s.ID)
	return s.ID
}
