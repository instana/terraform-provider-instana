package instana

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		Computed:    true,
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
		Type:        schema.TypeList,  // Represents the "time_threshold" block
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
	




)