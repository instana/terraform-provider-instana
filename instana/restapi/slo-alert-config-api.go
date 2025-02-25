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
	ID                  string                    	`json:"id"`
	Name                string                    	`json:"name"`
	Description         string                    	`json:"description"`
	Severity            int                       	`json:"severity"`
	Triggering          bool                      	`json:"triggering"`
	Enabled             bool                      	`json:"enabled"`
	Rule                SloAlertRule              	`json:"rule"`
	Threshold           SloAlertThreshold         	`json:"threshold"`
	TimeThreshold       SloAlertTimeThreshold     	`json:"timeThreshold"`
	SloIds              []string                 	`json:"sloIds"`
	AlertChannelIds     []string                  	`json:"alertChannelIds"`
    CustomerPayloadFields []CustomPayloadField[any] `json:"customPayloadFields"`
}

type SloAlertRule struct {
    AlertType string `json:"alertType"`
    Metric    string `json:"metric"`
}

type SloAlertThreshold struct {
    Type     string  `json:"type"`     
    Operator string  `json:"operator"`
    Value    float64 `json:"value"`
}

type SloAlertTimeThreshold struct {
	Timewindow 	int `json:"timewindow"`
	Expiry     	int `json:"expiry"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (s *SloAlertConfig) GetIDForResourcePath() string {
	fmt.Fprintln(os.Stderr, ">> GetIDForResourcePath: "+s.ID)
	return s.ID
}

func (config *SloAlertConfig) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return config.CustomerPayloadFields
}

func (config *SloAlertConfig) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	config.CustomerPayloadFields = fields
}