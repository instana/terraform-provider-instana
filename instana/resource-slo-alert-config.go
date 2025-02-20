package instana

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"

)

const ResourceInstanaSloAlertConfig = "instana_slo_alert_config"

const (
	//Slo Alert Config Field names for Terraform

	SloAlertConfigFieldName                    = "name"
	SloAlertConfigFieldDescription             = "description"
	SloAlertConfigFieldSeverity 			   = "severity"
	SloAlertConfigFieldTriggering 		   	   = "triggering"
	SloAlertConfigFieldAlertType 			   = "alert_type"
	SloAlertConfigFieldThreshold 			   = "threshold"
	SloAlertConfigFieldThresholdOperator       = "threshold_operator"
	SloAlertConfigFieldThresholdValue    	   = "threshold_value"
	SloAlertConfigFieldSloIds                  = "slo_ids"
	SloAlertConfigFieldAlertChannelIds         = "alert_channel_ids"
	SloAlertConfigFieldTimeThreshold           = "time_threshold"
	SloAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	SloAlertConfigFieldTimeThresholdExpiry     = "expiry"
	SloAlertConfigFieldCustomPayloadFields     = "custom_payload_fields"
	SloAlertConfigFieldEnabled                 = "enabled"

	SloAlertConfigFieldBurnRateTimeWindows      = "burn_rate_time_windows"
	SloAlertConfigFieldLongTimeWindow			= "long_time_window"
	SloAlertConfigFieldShortTimeWindow			= "short_time_window"
	SloAlertConfigFieldTimeWindowDuration		= "time_window_duration"
	SloAlertConfigFieldTimeWindowDurationType	= "time_window_duration_type"


	// Slo Alert Types for Terraform

	SloAlertConfigStatus            = "status"
	SloAlertConfigErrorBudget       = "error_budget"
	SloAlertConfigBurnRate          = "burn_rate"


)

var sloAlertConfigAlertTypeKeys = []string{
	"alert.0." + SloAlertConfigStatus,
	"alert.0." + SloAlertConfigErrorBudget,
	"alert.0." + SloAlertConfigBurnRate,
}



var (
	//SloAlertConfigName schema field definition of instana_slo_alert_config field name
	SloAlertConfigName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(0, 256),
		Description:  "The name of the SLO Alert config",
	}

	//SloAlertConfigDescription schema field definition of instana_slo_alert_config field description
	SloAlertConfigDescription = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The full name of the SLI config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}

	//SloAlertConfigSeverity schema field definition of instana_slo_alert_config field severity
	SloAlertConfigSeverity = &schema.Schema{
		Type:        schema.TypeInt,
		Required:     true,
		Description: "The severity of the alert when triggered",
	}

	SloAlertConfigTriggering = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also an Incident is triggered or not. The default is false",
	}

	SloAlertConfigAlertType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"status", "errorBudget", "burnRate"}, true),
		Description:  "What do you want to be alerted on? (Type of Smart Alert)",
	}

	SloAlertConfigThreshold = &schema.Schema{
		Type:        schema.TypeList,  // Represents the "threshold" block
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"operator": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The operator used to evaluate this rule.",
					ValidateFunc: validation.StringInSlice([]string{">=", "<=", ">", "<", "==", "!="}, false),
				},
				"value": {
					Type:        schema.TypeFloat,
					Required:    true,
					Description: "The threshold value for the alert condition.",
				},
			},
		},
	}
	

	SloAlertConfigSloIds = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "The SLO IDs that are monitored",
		MinItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	SloAlertConfigAlertChannelIds = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "The IDs of the Alert Channels",
		MinItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	SloAlertConfigTimeThreshold = &schema.Schema{
		Type:        schema.TypeList,  
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Defines the time threshold for triggering and suppressing alerts.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"time_window": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The duration for which the condition must be violated for the alert to be triggered (in ms).",
				},
				"expiry": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The duration for which the condition must remain suppressed for the alert to end (in ms).",
				},
			},
		},
	}

	// SloAlertConfigCustomPayloadFields = &schema.Schema{
	// 	Type:        schema.TypeList,  
	// 	Optional:    true,
	// 	Description: "Custom payload fields to include additional metadata in the alert.",
	// 	Elem: &schema.Resource{
	// 		Schema: map[string]*schema.Schema{
	// 			"key": {
	// 				Type:        schema.TypeString,
	// 				Required:    true,
	// 				Description: "The key name for the custom payload field.",
	// 			},
	// 			"value": {
	// 				Type:        schema.TypeString,
	// 				Required:    true,
	// 				Description: "The value associated with the custom payload field.",
	// 			},
	// 		},
	// 	},
	// }
	

	SloAlertConfigEnabled = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether this Alert is Enabled",
	}

	SloAlertConfigBurnRateTimeWindows = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Defines the burn rate time windows for evaluating alert conditions.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"long_time_window": {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Required: true,
					Description: "Defines the long time window duration and type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"duration": {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "The duration for the long time window.",
							},
							"duration_type": {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The unit of time for the long time window duration (e.g., 'minute', 'hour').",
								ValidateFunc: validation.StringInSlice([]string{"minute", "hour"}, false),
							},
						},
					},
				},
				"short_time_window": {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Required: true,
					Description: "Defines the short time window duration and type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"duration": {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "The duration for the short time window.",
							},
							"duration_type": {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The unit of time for the short time window duration (e.g., 'minute', 'hour').",
								ValidateFunc: validation.StringInSlice([]string{"minute", "hour"}, false),
							},
						},
					},
				},
			},
		},
	}
	
)

func NewSloAlertConfigResourceHandle() ResourceHandle[*restapi.SloAlertConfig] {
	Resource := &sloAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSloAlertConfig,
			Schema: map[string]*schema.Schema{
				SloAlertConfigFieldName:          	SloAlertConfigName,
				SloAlertConfigFieldDescription:     SloAlertConfigDescription,
				SloAlertConfigFieldSeverity:        SloAlertConfigSeverity,
				SloAlertConfigFieldTriggering:   	SloAlertConfigTriggering,
				SloAlertConfigFieldAlertType:   	SloAlertConfigAlertType,
				SloAlertConfigFieldThreshold:    	SloAlertConfigThreshold,
				SloAlertConfigFieldSloIds:  		SloAlertConfigSloIds,
				SloAlertConfigFieldAlertChannelIds: SloAlertConfigAlertChannelIds,
				SloAlertConfigFieldTimeThreshold:   SloAlertConfigTimeThreshold,
				DefaultCustomPayloadFieldsName:  	buildCustomPayloadFields(),
				SloAlertConfigFieldEnabled: 		SloAlertConfigEnabled,
			},
			SchemaVersion:    1,
			CreateOnly:       false,
			SkipIDGeneration: true,
		},
	}

	return Resource
}

type sloAlertConfigResource struct {
	metaData ResourceMetaData
}


func (r *sloAlertConfigResource) MetaData() *ResourceMetaData {
	resourceData := &r.metaData
	return resourceData
}

func (r *sloAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloAlertConfig] {
	x := api.SloAlertConfigs()
	return x
}

func (r *sloAlertConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *sloAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.sloAlertConfigSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.sloAlertConfigStateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *sloAlertConfigResource) sloAlertConfigStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
    if _, ok := state[SloAlertConfigFieldName]; ok {
        state[SloAlertConfigFieldName] = state[SloAlertConfigFieldName]
        delete(state, SloAlertConfigFieldName)
    }

    if _, ok := state[SloAlertConfigFieldThreshold]; ok {
        oldThreshold, isOldFormat := state[SloAlertConfigFieldThreshold].(string)
        if isOldFormat {
            // Convert old threshold format (string) to new list format
            state[SloAlertConfigFieldThreshold] = []interface{}{
                map[string]interface{}{
                    SloAlertConfigFieldThresholdOperator: ">=",
                    SloAlertConfigFieldThresholdValue:    oldThreshold,
                },
            }
        }
    }

    if _, ok := state[SloAlertConfigFieldTimeThreshold]; ok {
        oldTimeThreshold, isOldFormat := state[SloAlertConfigFieldTimeThreshold].(string)
        if isOldFormat {
            // Convert old time threshold format (string) to new list format
            state[SloAlertConfigFieldTimeThreshold] = []interface{}{
                map[string]interface{}{
                    SloAlertConfigFieldTimeThresholdExpiry:     60000,
                    SloAlertConfigFieldTimeThresholdTimeWindow: oldTimeThreshold,
                },
            }
        }
    }

    return state, nil
}


// tf state -> api
func (r *sloAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SloAlertConfig, error) {
    debug(">> MapStateToDataObject")
    debug(obj2json(d))
    sid := d.Id()
    if len(sid) == 0 {
        sid = RandomID()
    }
    // Convert threshold
    threshold, err := r.mapThresholdToAPIModel(d)
    if err != nil {
        return nil, err
    }
    // Convert time threshold
    timeThreshold, err := r.mapTimeThresholdToAPIModel(d)
    if err != nil {
        return nil, err
    }

	// customPayloadFields
	customPayloadFields, err := mapDefaultCustomPayloadFieldsFromSchema(d)
	if err != nil {
		return &restapi.SloAlertConfig{}, err
	}
    payload := &restapi.SloAlertConfig{
        ID:                     sid,
        Name:                   d.Get(SloAlertConfigFieldName).(string),
        Description:            d.Get(SloAlertConfigFieldDescription).(string),
        Severity:               d.Get(SloAlertConfigFieldSeverity).(int),
        Triggering:             d.Get(SloAlertConfigFieldTriggering).(bool),
        AlertType:              d.Get(SloAlertConfigFieldAlertType).(string),
        Threshold:              threshold,
        SloIds:                 convertInterfaceSliceToStringSlice(d.Get(SloAlertConfigFieldSloIds).([]interface{})),
        AlertChannelIds:        convertInterfaceSliceToStringSlice(d.Get(SloAlertConfigFieldAlertChannelIds).([]interface{})),
        TimeThreshold:          timeThreshold,
        CustomerPayloadFields:  customPayloadFields,
        Enabled:                d.Get(SloAlertConfigFieldEnabled).(bool),
    }
    return payload, nil
}

// Function to convert interface slice to string slice
func convertInterfaceSliceToStringSlice(input []interface{}) []string {
    result := make([]string, len(input))
    for i, v := range input {
        result[i] = v.(string)
    }
    return result
}
// Function to map threshold to API model
func (r *sloAlertConfigResource) mapThresholdToAPIModel(d *schema.ResourceData) (*restapi.Threshold, error) {
    thresholdStateObject := d.Get(SloAlertConfigFieldThreshold).([]interface{})
    if len(thresholdStateObject) == 1 {
        thresholdObject := thresholdStateObject[0].(map[string]interface{})
        return &restapi.Threshold{
            Operator: thresholdObject[SloAlertConfigFieldThresholdOperator].(string),
            Value:    thresholdObject[SloAlertConfigFieldThresholdValue].(float64),
        }, nil
    }
    return nil, errors.New("exactly one threshold configuration is required")
}
// Function to map time threshold to API model
func (r *sloAlertConfigResource) mapTimeThresholdToAPIModel(d *schema.ResourceData) (*restapi.TimeThreshold, error) {
    timeThresholdStateObject := d.Get(SloAlertConfigFieldTimeThreshold).([]interface{})
    if len(timeThresholdStateObject) == 1 {
        timeThresholdObject := timeThresholdStateObject[0].(map[string]interface{})
        return &restapi.TimeThreshold{
            Expiry:      timeThresholdObject[SloAlertConfigFieldTimeThresholdExpiry].(int),
            TimeWindow:  timeThresholdObject[SloAlertConfigFieldTimeThresholdTimeWindow].(int),
        }, nil
    }
    return nil, errors.New("exactly one time threshold configuration is required")
}

// root schema
func (r *sloAlertConfigResource) sloAlertConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SloAlertConfigFieldName:          	SloAlertConfigName,
			SloAlertConfigFieldDescription:     SloAlertConfigDescription,
			SloAlertConfigFieldSeverity:        SloAlertConfigSeverity,
			SloAlertConfigFieldTriggering:   	SloAlertConfigTriggering,
			SloAlertConfigFieldAlertType:   	SloAlertConfigAlertType,
			SloAlertConfigFieldThreshold:    	SloAlertConfigThreshold,
			SloAlertConfigFieldSloIds:  		SloAlertConfigSloIds,
			SloAlertConfigFieldAlertChannelIds: SloAlertConfigAlertChannelIds,
			SloAlertConfigFieldTimeThreshold:   SloAlertConfigTimeThreshold,
			DefaultCustomPayloadFieldsName:  	buildCustomPayloadFields(),
			SloAlertConfigFieldEnabled: 		SloAlertConfigEnabled,
		},
	}
}

func (r *sloAlertConfigResource) UpdateState(d *schema.ResourceData, sloAlertConfig *restapi.SloAlertConfig) error {
    debug(">> UpdateState")
    
    // Convert Threshold to Terraform state
    threshold := map[string]interface{}{
        SloAlertConfigFieldThresholdOperator: sloAlertConfig.Threshold.Operator,
        SloAlertConfigFieldThresholdValue:    sloAlertConfig.Threshold.Value,
    }

    // Convert TimeThreshold to Terraform state
    timeThreshold := map[string]interface{}{
        SloAlertConfigFieldTimeThresholdExpiry:     sloAlertConfig.TimeThreshold.Expiry,
        SloAlertConfigFieldTimeThresholdTimeWindow: sloAlertConfig.TimeThreshold.TimeWindow,
    }

    // Prepare Terraform state data
    tfData := map[string]interface{}{
        SloAlertConfigFieldName:          sloAlertConfig.Name,
        SloAlertConfigFieldDescription:   sloAlertConfig.Description,
        SloAlertConfigFieldSeverity:      sloAlertConfig.Severity,
        SloAlertConfigFieldTriggering:    sloAlertConfig.Triggering,
        SloAlertConfigFieldAlertType:     sloAlertConfig.AlertType,
        SloAlertConfigFieldThreshold:     []interface{}{threshold},
        SloAlertConfigFieldSloIds:        sloAlertConfig.SloIds,
        SloAlertConfigFieldAlertChannelIds: sloAlertConfig.AlertChannelIds,
        SloAlertConfigFieldTimeThreshold: []interface{}{timeThreshold},
        SloAlertConfigFieldEnabled:       sloAlertConfig.Enabled,
    }

    d.SetId(sloAlertConfig.ID)

    debug(">> UpdateState with: " + obj2json(tfData))

    return tfutils.UpdateState(d, tfData)
}
