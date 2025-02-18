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