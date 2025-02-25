package instana

import (
	"fmt"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
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
	SloAlertConfigFieldThresholdOperator       = "operator"
	SloAlertConfigFieldThresholdValue    	   = "value"
	SloAlertConfigFieldSloIds                  = "slo_ids"
	SloAlertConfigFieldAlertChannelIds         = "alert_channel_ids"
	SloAlertConfigFieldTimeThreshold           = "time_threshold"
	SloAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	SloAlertConfigFieldTimeThresholdExpiry     = "expiry"
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
		ValidateFunc: validation.StringInSlice([]string{"status", "error_budget", "burn_rate"}, false),
		Description:  "What do you want to be alerted on? (Type of Smart Alert)",
	}

	SloAlertConfigThreshold = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The type of threshold (should be staticThreshold).",
					ValidateFunc: validation.StringInSlice([]string{"staticThreshold"}, false),
				},
				"operator": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The operator used to evaluate this rule.",
					ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdOperators.ToStringSlice(), true),
				},
				"value": {
					Type:        schema.TypeFloat,
					Required:    true,
					Description: "The threshold value for the alert condition.",
				},
			},
		},
	}

	SloAlertConfigBurnRateTimeWindows = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Optional:    true,
		Description: "Defines the burn rate time windows for evaluating alert conditions.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SloAlertConfigFieldLongTimeWindow: {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Required: true,
					Description: "Defines the long time window duration and type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SloAlertConfigFieldTimeWindowDuration: {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "The duration for the long time window.",
							},
							SloAlertConfigFieldTimeWindowDurationType: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The unit of time for the long time window duration (e.g., 'MINUTE', 'HOUR', 'DAY').",
								ValidateFunc: validation.StringInSlice([]string{"MINUTE", "HOUR", "DAY"}, false), // Case-sensitive validation
							},
						},
					},
				},
				SloAlertConfigFieldShortTimeWindow: {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Required: true,
					Description: "Defines the short time window duration and type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SloAlertConfigFieldTimeWindowDuration: {
								Type:        schema.TypeInt,
								Required:    true,
								Description: "The duration for the short time window.",
							},
							SloAlertConfigFieldTimeWindowDurationType: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The unit of time for the short time window duration (e.g., 'MINUTE', 'HOUR', 'DAY').",
								ValidateFunc: validation.StringInSlice([]string{"MINUTE", "HOUR", "DAY"}, false), // Case-sensitive validation
							},
						},
					},
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

	SloAlertConfigEnabled = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether this Alert is Enabled",
	}
)

// wrapSloCustomPayloadFields ensures CustomerPayloadFields values are correctly typed
func wrapSloCustomPayloadFields(fields []restapi.CustomPayloadField[any]) []restapi.CustomPayloadField[any] {
    for i, field := range fields {
        if field.Type != restapi.DynamicCustomPayloadType {
            if str, ok := field.Value.(string); ok {
                fields[i].Value = restapi.StaticStringCustomPayloadFieldValue(str)
            }
        }
    }
    return fields
}

func NewSloAlertConfigResourceHandle() ResourceHandle[*restapi.SloAlertConfig] {
	Resource := &sloAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSloAlertConfig,
			Schema: map[string]*schema.Schema{
				SloAlertConfigFieldName:          			SloAlertConfigName,
				SloAlertConfigFieldDescription:     		SloAlertConfigDescription,
				SloAlertConfigFieldSeverity:        		SloAlertConfigSeverity,
				SloAlertConfigFieldTriggering:   			SloAlertConfigTriggering,
				SloAlertConfigFieldAlertType:   			SloAlertConfigAlertType,
				SloAlertConfigFieldThreshold:    			SloAlertConfigThreshold,
				SloAlertConfigFieldSloIds:  				SloAlertConfigSloIds,
				SloAlertConfigFieldAlertChannelIds: 		SloAlertConfigAlertChannelIds,
				SloAlertConfigFieldTimeThreshold:   		SloAlertConfigTimeThreshold,
				DefaultCustomPayloadFieldsName: 			buildCustomPayloadFields(),
				SloAlertConfigFieldEnabled: 				SloAlertConfigEnabled,
			},
			SchemaVersion:    1,
			CreateOnly:       false,
			SkipIDGeneration: true,
		},
	}

	return Resource
}

func mapAlertTypeToAPI(terraformAlertType string) (string, string, error) {
	normalizedType := normalizeAlertType(terraformAlertType)

	switch normalizedType {
    case "status":
        return "SERVICE_LEVELS_OBJECTIVE", "STATUS", nil
    case "error_budget":
        return "ERROR_BUDGET", "BURNED_PERCENTAGE", nil
    case "burn_rate":
        return "ERROR_BUDGET", "BURN_RATE", nil
    default:
		fmt.Printf("WARNING: Unknown alert type '%s' received from Terraform\n", terraformAlertType)
        return "", "", fmt.Errorf("invalid alert_type: %s", terraformAlertType)
    }
}

// Normalize Terraform input values
func normalizeAlertType(alertType string) string {
	switch alertType {
	case "errorBudget", "ErrorBudget":
		return "error_budget"
	case "burnRate", "BurnRate":
		return "burn_rate"
	case "status", "Status":
		return "status"
	default:
		return alertType
	}
}


type sloAlertConfigResource struct {
	metaData ResourceMetaData
}

func (r *sloAlertConfigResource) MetaData() *ResourceMetaData {
	resourceData := &r.metaData
	return resourceData
}

func (r *sloAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloAlertConfig] {
	return api.SloAlertConfig()
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

func (r *sloAlertConfigResource) UpdateState(d *schema.ResourceData, sloAlertConfig *restapi.SloAlertConfig) error {
    debug(">> UpdateState")

	thresholdType := sloAlertConfig.Threshold.Type
	if thresholdType == "static" {  
		thresholdType = "staticThreshold"
	}

	threshold := map[string]interface{}{
		"type": thresholdType,
		SloAlertConfigFieldThresholdOperator: sloAlertConfig.Threshold.Operator,
		SloAlertConfigFieldThresholdValue:    sloAlertConfig.Threshold.Value,
	}
	
    timeThreshold := map[string]interface{}{
        SloAlertConfigFieldTimeThresholdExpiry:     sloAlertConfig.TimeThreshold.Expiry,
        SloAlertConfigFieldTimeThresholdTimeWindow: sloAlertConfig.TimeThreshold.Timewindow,
    }

	var terraformAlertType string

	// Reverse Map API's "alertType" and "metric" to Terraform's expected values
	switch sloAlertConfig.Rule.AlertType {
	case "SERVICE_LEVELS_OBJECTIVE":
		if sloAlertConfig.Rule.Metric == "STATUS" {
			terraformAlertType = "status"
		}
	case "ERROR_BUDGET":
		if sloAlertConfig.Rule.Metric == "BURNED_PERCENTAGE" {
			terraformAlertType = "error_budget"
		} else if sloAlertConfig.Rule.Metric == "BURN_RATE" {
			terraformAlertType = "burn_rate"
		}
	}
	
	// Ensure consistency before storing in state
	terraformAlertType = normalizeAlertType(terraformAlertType)
	
	if terraformAlertType == "" {
		return fmt.Errorf("unexpected alertType/metric from API: %v", sloAlertConfig.Rule)
	}
	
	// Preprocess CustomerPayloadFields to ensure correct type
	sloAlertConfig.CustomerPayloadFields = wrapSloCustomPayloadFields(sloAlertConfig.CustomerPayloadFields)

	tfData := map[string]interface{}{
		SloAlertConfigFieldName:          sloAlertConfig.Name,
		SloAlertConfigFieldDescription:   sloAlertConfig.Description,
		SloAlertConfigFieldSeverity:      sloAlertConfig.Severity,
		SloAlertConfigFieldTriggering:    sloAlertConfig.Triggering,
		SloAlertConfigFieldAlertType:     terraformAlertType, 
		SloAlertConfigFieldThreshold:     []interface{}{threshold},
		SloAlertConfigFieldSloIds:        sloAlertConfig.SloIds,
		SloAlertConfigFieldAlertChannelIds: sloAlertConfig.AlertChannelIds,
		SloAlertConfigFieldTimeThreshold: []interface{}{timeThreshold},
		DefaultCustomPayloadFieldsName:   mapCustomPayloadFieldsToSchema(sloAlertConfig),
		SloAlertConfigFieldEnabled:       sloAlertConfig.Enabled,
	}	

    d.SetId(sloAlertConfig.ID)

    debug(">> UpdateState with: " + obj2json(tfData))

    return tfutils.UpdateState(d, tfData)
}

// convertToFloat64 safely converts different numeric types to float64.
func convertToFloat64(value interface{}) (float64, error) {
    switch v := value.(type) {
    case float64:
        return v, nil
    case float32:
        return float64(v), nil
    case int:
        return float64(v), nil
    case int64:
        return float64(v), nil
    case string:
        var parsedValue float64
        _, err := fmt.Sscanf(v, "%f", &parsedValue)
        if err != nil {
            return 0, fmt.Errorf("unable to parse float64 from string: %s", v)
        }
        return parsedValue, nil
    default:
        return 0, fmt.Errorf("unexpected type for float64 conversion: %T", value)
    }
}

// tf state -> api
func (r *sloAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SloAlertConfig, error) {
    debug(">> MapStateToDataObject")
    debug(obj2json(d))
    sid := d.Id()
    if len(sid) == 0 {
        sid = RandomID()
    }
	
// Convert threshold from Terraform state
thresholdStateObject := d.Get(SloAlertConfigFieldThreshold).([]interface{})
var threshold restapi.SloAlertThreshold

if len(thresholdStateObject) > 0 {
    thresholdObject, ok := thresholdStateObject[0].(map[string]interface{})
    if ok {
        operatorRaw, opOK := thresholdObject[SloAlertConfigFieldThresholdOperator]
        valueRaw, valOK := thresholdObject[SloAlertConfigFieldThresholdValue]

        if opOK && valOK {
            operator := fmt.Sprintf("%v", operatorRaw)
            value, err := convertToFloat64(valueRaw) 
            if err != nil {
                return nil, fmt.Errorf("threshold value is invalid: %v", err)
            }

            threshold = restapi.SloAlertThreshold{
                Type:     "staticThreshold",
                Operator: operator,
                Value:    value,
            }
        } else {
            return nil, fmt.Errorf("threshold operator or value is missing or incorrect type")
        }
    }
}

	// Convert time threshold from Terraform state
	timeThresholdStateObject := d.Get(SloAlertConfigFieldTimeThreshold).([]interface{})
	var timeThreshold restapi.SloAlertTimeThreshold

	if len(timeThresholdStateObject) > 0 {
		timeThresholdObject, ok := timeThresholdStateObject[0].(map[string]interface{})
		if ok {
			expiry, expiryOK := timeThresholdObject[SloAlertConfigFieldTimeThresholdExpiry].(int)
			timewindow, timeWindowOK := timeThresholdObject[SloAlertConfigFieldTimeThresholdTimeWindow].(int)

			if expiryOK && timeWindowOK {
				timeThreshold = restapi.SloAlertTimeThreshold{
					Expiry:     expiry,
					Timewindow: timewindow,
				}
			} else {
				return nil, fmt.Errorf("time threshold expiry or time window is missing or of incorrect type")
			}
		}
	}

	// Custom Payload Fields
	customPayloadFields, err := mapDefaultCustomPayloadFieldsFromSchema(d)
	if err != nil {
		return nil, fmt.Errorf("error processing custom payload fields: %w", err)
	}

	customPayloadFields = wrapSloCustomPayloadFields(customPayloadFields) 

	terraformAlertType := d.Get(SloAlertConfigFieldAlertType).(string)
    terraformAlertType = normalizeAlertType(terraformAlertType)
	
	apiAlertType, apiMetric, err := mapAlertTypeToAPI(terraformAlertType)
	if err != nil {
		return nil, fmt.Errorf("invalid alert_type: %v", err)
	}

	// Construct the API-compatible Rule object
	rule := restapi.SloAlertRule{
		AlertType: apiAlertType,
		Metric:    apiMetric,
	}

	// Construct API payload
	payload := &restapi.SloAlertConfig{
		ID:          sid,
		Name:        d.Get(SloAlertConfigFieldName).(string),
		Description: d.Get(SloAlertConfigFieldDescription).(string),
		Severity:    d.Get(SloAlertConfigFieldSeverity).(int),
		Triggering:  d.Get(SloAlertConfigFieldTriggering).(bool),
		Enabled:     d.Get(SloAlertConfigFieldEnabled).(bool),
		Rule:        rule,  
		Threshold:   threshold,
		SloIds:      convertInterfaceSliceToStringSlice(d.Get(SloAlertConfigFieldSloIds).([]interface{})),
		AlertChannelIds: convertInterfaceSliceToStringSlice(d.Get(SloAlertConfigFieldAlertChannelIds).([]interface{})),
		TimeThreshold: timeThreshold,
		CustomerPayloadFields: customPayloadFields,
	}
	return payload, nil
	}

func convertInterfaceSliceToStringSlice(input []interface{}) []string {
    result := make([]string, len(input))
    for i, v := range input {
        result[i] = v.(string)
    }
    return result
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
