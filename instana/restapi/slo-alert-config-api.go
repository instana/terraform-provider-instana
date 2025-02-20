package restapi

import (
	"fmt"
	"os"
)

const (
	//SloAlertConfigResourcePath path to slo alert config resource of Instana RESTful API
	SloAlertConfigResourcePath = EventSettingsBasePath + "/global-alert-configs/service-levels"
)

// SloAlertConfig represents the REST resource of SLO Alert Configuration at Instana
type SloAlertConfig struct {
	ID          		string      `json:"id"`
	Name        		string      `json:"name"`
	Description 		string     	`json:"description"`
	Severity    		int      	`json:"severity"`
	Triggering  		bool      	`json:"triggering"`
	Enabled    			bool     	`json:"enabled"`
	AlertType			string		`json:"alertType"`
	Threshold   		interface{} `json:"threshold"`
	TimeThreshold  		interface{} `json:"timeThreshold"`
	SloIds      		[]string    `json:"sloIds"`
	AlertChannelIds 	[]string 	`json:"alertChannelIds"`
	CustomPayloadFields []string 	`json:"customPayloadFields"`
}

type SloAlertThreshold struct {
	Operator     string        `json:"operator"`
	Value 		 float64 	   `json:"value"`
}

type SloAlertTimeThreshold struct {
	Timewindow      int        `json:"timewindow"`
	expiry			int    	   `json:"expiry"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (s *SloAlertConfig) GetIDForResourcePath() string {
	fmt.Fprintln(os.Stderr, ">> GetIDForResourcePath: "+s.ID)
	return s.ID
}