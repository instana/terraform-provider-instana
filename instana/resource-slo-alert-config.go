package instana

import (
	"fmt"
	"strconv"

	// "log"
	// "encoding/json"
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const ResourceInstanaSloAlertConfig = "instana_slo_alert_config"

const (
	//Slo Alert Config Field names for Terraform
	SloAlertConfigFieldName                            = "name"
	SloAlertConfigFieldFullName                        = "full_name"
	SloAlertConfigFieldDescription                     = "description"
	SloAlertConfigFieldSeverity                        = "severity"
	SloAlertConfigFieldTriggering                      = "triggering"
	SloAlertConfigFieldAlertType                       = "alert_type"
	SloAlertConfigFieldThreshold                       = "threshold"
	SloAlertConfigFieldThresholdType                   = "type"
	SloAlertConfigFieldThresholdOperator               = "operator"
	SloAlertConfigFieldThresholdValue                  = "value"
	SloAlertConfigFieldSloIds                          = "slo_ids"
	SloAlertConfigFieldAlertChannelIds                 = "alert_channel_ids"
	SloAlertConfigFieldTimeThreshold                   = "time_threshold"
	SloAlertConfigFieldTimeThresholdWarmUp             = "warm_up"
	SloAlertConfigFieldTimeThresholdCoolDown           = "cool_down"
	SloAlertConfigFieldEnabled                         = "enabled"
	SloAlertConfigFieldBurnRateConfig                  = "burn_rate_config"
	SloAlertConfigFieldBurnRateConfigDuration          = "duration"
	SloAlertConfigFieldBurnRateConfigThresholdValue    = "threshold_value"
	SloAlertConfigFieldBurnRateConfigThresholdOperator = "threshold_operator"
	SloAlertConfigFieldBurnRateConfigDurationUnitType  = "duration_unit_type"
	SloAlertConfigFieldBurnRateConfigAlertWindowType   = "alert_window_type"

	// Slo Alert Types for Terraform
	SloAlertConfigStatus      = "status"
	SloAlertConfigErrorBudget = "error_budget"
	SloAlertConfigBurnRateV2  = "burn_rate_v2"
)

var (
	//SloAlertConfigName schema field definition of instana_slo_alert_config field name
	SloAlertConfigName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(0, 256),
		Description:  "The name of the SLO Alert config",
	}

	//SloAlertConfigFullName schema field definition of instana_slo_alert_config field full_name
	SloAlertConfigFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the SLO Alert config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
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
		Required:    true,
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
		ValidateFunc: validation.StringInSlice([]string{"status", "error_budget", "burn_rate_v2"}, false),
		Description:  "What do you want to be alerted on? (Type of Smart Alert: status, error_budget, burn_rate_v2)",
	}

	SloAlertConfigThreshold = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SloAlertConfigFieldThresholdType: {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "staticThreshold",
					Description:  "The type of threshold (should be staticThreshold).",
					ValidateFunc: validation.StringInSlice([]string{"staticThreshold"}, false),
				},
				SloAlertConfigFieldThresholdOperator: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The operator used to evaluate this rule.",
					ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdOperators.ToStringSlice(), true),
				},
				SloAlertConfigFieldThresholdValue: {
					Type:        schema.TypeFloat,
					Required:    true,
					Description: "The threshold value for the alert condition.",
				},
			},
		},
	}

	SloAlertConfigBurnRateConfig = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of burn rate configs fields.",
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	SloAlertConfigSloIds = &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Description: "The SLO IDs that are monitored",
		MinItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	SloAlertConfigAlertChannelIds = &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Description: "The IDs of the Alert Channels",
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
				SloAlertConfigFieldTimeThresholdWarmUp: {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The duration for which the condition must be violated for the alert to be triggered (in ms).",
				},
				SloAlertConfigFieldTimeThresholdCoolDown: {
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

func NewSloAlertConfigResourceHandle() ResourceHandle[*restapi.SloAlertConfig] {
	Resource := &sloAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSloAlertConfig,
			Schema: map[string]*schema.Schema{
				SloAlertConfigFieldName:            SloAlertConfigName,
				SloAlertConfigFieldDescription:     SloAlertConfigDescription,
				SloAlertConfigFieldSeverity:        SloAlertConfigSeverity,
				SloAlertConfigFieldTriggering:      SloAlertConfigTriggering,
				SloAlertConfigFieldAlertType:       SloAlertConfigAlertType,
				SloAlertConfigFieldThreshold:       SloAlertConfigThreshold,
				SloAlertConfigFieldSloIds:          SloAlertConfigSloIds,
				SloAlertConfigFieldAlertChannelIds: SloAlertConfigAlertChannelIds,
				SloAlertConfigFieldTimeThreshold:   SloAlertConfigTimeThreshold,
				DefaultCustomPayloadFieldsName:     buildCustomPayloadFields(),
				SloAlertConfigFieldEnabled:         SloAlertConfigEnabled,
				SloAlertConfigFieldBurnRateConfig:  SloAlertConfigBurnRateConfig,
			},
			SchemaVersion:    1,
			CreateOnly:       false,
			SkipIDGeneration: true,
		},
	}
	return Resource
}

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

func mapAlertTypeToAPI(terraformAlertType string) (string, string, error) {
	normalizedType := normalizeAlertType(terraformAlertType)

	switch normalizedType {
	case "status":
		return "SERVICE_LEVELS_OBJECTIVE", "STATUS", nil
	case "error_budget":
		return "ERROR_BUDGET", "BURNED_PERCENTAGE", nil
	case "burn_rate_v2":
		return "ERROR_BUDGET", "BURN_RATE_V2", nil
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
	case "status", "Status":
		return "status"
	case "burnRateV2", "BurnRateV2":
		return "burn_rate_v2"
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
	if _, ok := state[SloAlertConfigFieldFullName]; ok {
		state[SloAlertConfigFieldName] = state[SloAlertConfigFieldFullName]
		delete(state, SloAlertConfigFieldFullName)
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
					SloAlertConfigFieldTimeThresholdCoolDown: 60000,
					SloAlertConfigFieldTimeThresholdWarmUp:   oldTimeThreshold,
				},
			}
		}
	}

	return state, nil
}

func (r *sloAlertConfigResource) UpdateState(d *schema.ResourceData, sloAlertConfig *restapi.SloAlertConfig) error {
	debug(">> UpdateState")

	var thresholdType, thresholdOperator string

	var thresholdValue float64

	if sloAlertConfig.Threshold != nil {

		thresholdType = sloAlertConfig.Threshold.Type

		if thresholdType == "static" {

			thresholdType = "staticThreshold"

		}

		thresholdOperator = sloAlertConfig.Threshold.Operator

		thresholdValue = sloAlertConfig.Threshold.Value

	} else {

		thresholdType = ""

		thresholdOperator = ""

		thresholdValue = 0
	}

	threshold := map[string]interface{}{
		SloAlertConfigFieldThresholdType:     thresholdType,
		SloAlertConfigFieldThresholdOperator: thresholdOperator,
		SloAlertConfigFieldThresholdValue:    thresholdValue,
	}

	timeThreshold := map[string]interface{}{
		SloAlertConfigFieldTimeThresholdCoolDown: sloAlertConfig.TimeThreshold.Expiry,
		SloAlertConfigFieldTimeThresholdWarmUp:   sloAlertConfig.TimeThreshold.TimeWindow,
	}

	var terraformAlertType string

	// Reverse Map API - "alertType" and "metric" to Terraform's expected values
	switch sloAlertConfig.Rule.AlertType {
	case "SERVICE_LEVELS_OBJECTIVE":
		if sloAlertConfig.Rule.Metric == "STATUS" {
			terraformAlertType = "status"
		}
	case "ERROR_BUDGET":
		if sloAlertConfig.Rule.Metric == "BURNED_PERCENTAGE" {
			terraformAlertType = "error_budget"
		} else if sloAlertConfig.Rule.Metric == "BURN_RATE_V2" {
			terraformAlertType = "burn_rate_v2"
		}
	}

	terraformAlertType = normalizeAlertType(terraformAlertType)

	if terraformAlertType == "" {
		return fmt.Errorf("unexpected alertType/metric from API: %v", sloAlertConfig.Rule)
	}

	// Preprocess CustomerPayloadFields to ensure correct type
	sloAlertConfig.CustomerPayloadFields = wrapSloCustomPayloadFields(sloAlertConfig.CustomerPayloadFields)

	// Handle burn_rate_config
	var burnRateConfigs []interface{}
	if sloAlertConfig.BurnRateConfigs != nil {
		for _, cfg := range *sloAlertConfig.BurnRateConfigs {
			burnRateConfigs = append(burnRateConfigs, map[string]interface{}{
				SloAlertConfigFieldBurnRateConfigAlertWindowType:   cfg.AlertWindowType,
				SloAlertConfigFieldBurnRateConfigDuration:          fmt.Sprintf("%d", cfg.Duration),
				SloAlertConfigFieldBurnRateConfigDurationUnitType:  cfg.DurationUnitType,
				SloAlertConfigFieldBurnRateConfigThresholdOperator: cfg.Threshold.Operator,
				SloAlertConfigFieldBurnRateConfigThresholdValue:    fmt.Sprintf("%.2f", cfg.Threshold.Value),
			})
		}
	}

	tfData := map[string]interface{}{
		SloAlertConfigFieldName:            sloAlertConfig.Name,
		SloAlertConfigFieldDescription:     sloAlertConfig.Description,
		SloAlertConfigFieldSeverity:        sloAlertConfig.Severity,
		SloAlertConfigFieldTriggering:      sloAlertConfig.Triggering,
		SloAlertConfigFieldAlertType:       terraformAlertType,
		SloAlertConfigFieldThreshold:       []interface{}{threshold},
		SloAlertConfigFieldSloIds:          convertSetToStringSlice(d.Get(SloAlertConfigFieldSloIds).(*schema.Set)),
		SloAlertConfigFieldAlertChannelIds: convertSetToStringSlice(d.Get(SloAlertConfigFieldAlertChannelIds).(*schema.Set)),
		SloAlertConfigFieldTimeThreshold:   []interface{}{timeThreshold},
		DefaultCustomPayloadFieldsName:     mapCustomPayloadFieldsToSchema(sloAlertConfig),
		SloAlertConfigFieldEnabled:         sloAlertConfig.Enabled,
		SloAlertConfigFieldBurnRateConfig:  burnRateConfigs,
	}

	d.SetId(sloAlertConfig.ID)

	debug(">> UpdateState with: " + obj2json(tfData))

	return tfutils.UpdateState(d, tfData)
}

func convertSetToStringSlice(set *schema.Set) []string {
	var result []string
	for _, v := range set.List() {
		result = append(result, v.(string))
	}
	return result
}

// convert different numeric types to float64.
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

	terraformAlertType := normalizeAlertType(d.Get(SloAlertConfigFieldAlertType).(string))
	var threshold *restapi.SloAlertThreshold
	// Convert threshold from Terraform state
	if terraformAlertType != SloAlertConfigBurnRateV2 {
		thresholdStateObject := d.Get(SloAlertConfigFieldThreshold).([]interface{})
		if len(thresholdStateObject) > 0 {
			thresholdObject := thresholdStateObject[0].(map[string]interface{})
			operatorRaw, opOK := thresholdObject[SloAlertConfigFieldThresholdOperator]
			valueRaw, valOK := thresholdObject[SloAlertConfigFieldThresholdValue]
			if opOK && valOK {
				operator := fmt.Sprintf("%v", operatorRaw)
				value, err := convertToFloat64(valueRaw)
				if err != nil {
					return nil, fmt.Errorf("threshold value is invalid: %v", err)
				}
				thresholdType := "staticThreshold"
				if typeRaw, typeOK := thresholdObject[SloAlertConfigFieldThresholdType]; typeOK && typeRaw != nil {
					thresholdType = typeRaw.(string)
				}
				threshold = &restapi.SloAlertThreshold{
					Type:     thresholdType,
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
		timeThresholdObject := timeThresholdStateObject[0].(map[string]interface{})
		warmUp, warmUpOK := timeThresholdObject[SloAlertConfigFieldTimeThresholdWarmUp].(int)
		coolDown, coolDownOK := timeThresholdObject[SloAlertConfigFieldTimeThresholdCoolDown].(int)
		if warmUpOK && coolDownOK {
			timeThreshold = restapi.SloAlertTimeThreshold{
				TimeWindow: warmUp,
				Expiry:     coolDown,
			}
		} else {
			return nil, fmt.Errorf("time threshold warm_up or cool_down is missing or incorrect type")
		}
	}

	// Custom Payload Fields
	customPayloadFields, err := mapDefaultCustomPayloadFieldsFromSchema(d)
	if err != nil {
		return nil, fmt.Errorf("error processing custom payload fields: %w", err)
	}
	customPayloadFields = wrapSloCustomPayloadFields(customPayloadFields)

	// Alert Type
	apiAlertType, apiMetric, err := mapAlertTypeToAPI(terraformAlertType)
	if err != nil {
		return nil, fmt.Errorf("invalid alert_type: %v", err)
	}
	rule := restapi.SloAlertRule{
		AlertType: apiAlertType,
		Metric:    apiMetric,
	}

	//Burn Rate Config
	var burnRateConfigs []restapi.BurnRateConfig
	if terraformAlertType == SloAlertConfigBurnRateV2 {
		raw, ok := d.GetOk(SloAlertConfigFieldBurnRateConfig)
		if !ok || raw == nil {
			return nil, fmt.Errorf("burn_rate_config is required for alert_type 'burn_rate_v2'")
		}

		rawList := raw.([]interface{})
		if len(rawList) == 0 {
			return nil, fmt.Errorf("burn_rate_config must contain at least one item for alert_type 'burn_rate_v2'")
		}

		for _, item := range rawList {
			obj := item.(map[string]interface{})
			duration, err := strconv.Atoi(obj["duration"].(string))
			if err != nil {
				return nil, fmt.Errorf("invalid duration: %v", err)
			}
			alertWindowType := obj["alert_window_type"].(string)
			durationUnitType := obj["duration_unit_type"].(string)
			operator := obj["threshold_operator"].(string)
			value, err := strconv.ParseFloat(obj["threshold_value"].(string), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid threshold_value: %v", err)
			}
			burnRateConfigs = append(burnRateConfigs, restapi.BurnRateConfig{
				AlertWindowType:  alertWindowType,
				Duration:         duration,
				DurationUnitType: durationUnitType,
				Threshold: restapi.ServiceLevelsStaticThresholdConfig{
					Operator: operator,
					Value:    value,
				},
			})
		}
	}

	// Construct payload
	payload := &restapi.SloAlertConfig{
		ID:                    sid,
		Name:                  d.Get(SloAlertConfigFieldName).(string),
		Description:           d.Get(SloAlertConfigFieldDescription).(string),
		Severity:              d.Get(SloAlertConfigFieldSeverity).(int),
		Triggering:            d.Get(SloAlertConfigFieldTriggering).(bool),
		Enabled:               d.Get(SloAlertConfigFieldEnabled).(bool),
		Rule:                  rule,
		Threshold:             threshold,
		TimeThreshold:         timeThreshold,
		SloIds:                convertSetToStringSlice(d.Get(SloAlertConfigFieldSloIds).(*schema.Set)),
		AlertChannelIds:       convertSetToStringSlice(d.Get(SloAlertConfigFieldAlertChannelIds).(*schema.Set)),
		CustomerPayloadFields: customPayloadFields,
		BurnRateConfigs:       &burnRateConfigs,
	}

	return payload, nil
}

// contains checks if a string is present in a slice of strings
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Schema
func (r *sloAlertConfigResource) sloAlertConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SloAlertConfigFieldName:            SloAlertConfigName,
			SloAlertConfigFieldFullName:        SloAlertConfigFullName,
			SloAlertConfigFieldDescription:     SloAlertConfigDescription,
			SloAlertConfigFieldSeverity:        SloAlertConfigSeverity,
			SloAlertConfigFieldTriggering:      SloAlertConfigTriggering,
			SloAlertConfigFieldAlertType:       SloAlertConfigAlertType,
			SloAlertConfigFieldThreshold:       SloAlertConfigThreshold,
			SloAlertConfigFieldSloIds:          SloAlertConfigSloIds,
			SloAlertConfigFieldAlertChannelIds: SloAlertConfigAlertChannelIds,
			SloAlertConfigFieldTimeThreshold:   SloAlertConfigTimeThreshold,
			DefaultCustomPayloadFieldsName:     buildCustomPayloadFields(),
			SloAlertConfigFieldBurnRateConfig:  SloAlertConfigBurnRateConfig,
			SloAlertConfigFieldEnabled:         SloAlertConfigEnabled,
		},
	}
}
