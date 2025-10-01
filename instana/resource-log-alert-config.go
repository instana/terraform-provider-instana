package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const ResourceInstanaLogAlertConfig = "instana_log_alert_config"

const (
	LogAlertConfigFieldName           = "name"
	LogAlertConfigFieldDescription    = "description"
	LogAlertConfigFieldAlertChannels  = "alert_channels"
	LogAlertConfigFieldGracePeriod    = "grace_period"
	LogAlertConfigFieldGroupBy        = "group_by"
	LogAlertConfigFieldGroupByTagName = "tag_name"
	LogAlertConfigFieldGroupByKey     = "key"
	LogAlertConfigFieldGranularity    = "granularity"
	LogAlertConfigFieldTagFilter      = "tag_filter"

	LogAlertConfigFieldRules             = "rules"
	LogAlertConfigFieldAlertType         = "alert_type"
	LogAlertConfigFieldMetricName        = "metric_name"
	LogAlertConfigFieldAggregation       = "aggregation"
	LogAlertConfigFieldThresholdOperator = "threshold_operator"
	LogAlertConfigFieldThreshold         = "threshold"
	LogAlertConfigFieldWarning           = "warning"
	LogAlertConfigFieldCritical          = "critical"
	LogAlertConfigFieldType              = "type"
	LogAlertConfigFieldValue             = "value"

	// LogAlertTypeLogCount is the constant for the log count alert type
	LogAlertTypeLogCount = "log.count"

	LogAlertConfigFieldTimeThreshold                     = "time_threshold"
	LogAlertConfigFieldTimeThresholdTimeWindow           = "time_window"
	LogAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
)

var (
	logAlertConfigSchemaAlertChannels = &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 0,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				ResourceFieldThresholdRuleWarningSeverity: {
					Type:     schema.TypeList,
					MinItems: 0,
					MaxItems: 1024,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "List of IDs of alert channels defined in Instana.",
				},
				ResourceFieldThresholdRuleCriticalSeverity: {
					Type:     schema.TypeList,
					MinItems: 0,
					MaxItems: 1024,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "List of IDs of alert channels defined in Instana.",
				},
			},
		},
		Description: "Set of alert channel IDs associated with the severity.",
	}

	logAlertConfigSchemaGroupBy = &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 0,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				LogAlertConfigFieldGroupByTagName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The tag name used for grouping",
				},
				LogAlertConfigFieldGroupByKey: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The key used for grouping",
				},
			},
		},
		Description: "The grouping tags used to group the metric results.",
	}

	logAlertConfigSchemaRules = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Configuration for the log alert rule",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				LogAlertConfigFieldMetricName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The metric name of the log alert rule",
				},
				LogAlertConfigFieldAlertType: {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      LogAlertTypeLogCount,
					ValidateFunc: validation.StringInSlice([]string{LogAlertTypeLogCount}, false),
					Description:  "The type of the log alert rule (only 'log.count' is supported)",
				},
				LogAlertConfigFieldAggregation: {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      string(restapi.SumAggregation),
					ValidateFunc: validation.StringInSlice([]string{string(restapi.SumAggregation)}, false),
					Description:  "The aggregation method to use for the log alert (only 'SUM' is supported)",
				},
				LogAlertConfigFieldThresholdOperator: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The operator which will be applied to evaluate the threshold",
					ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdOperators.ToStringSlice(), true),
				},
				LogAlertConfigFieldThreshold: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Threshold configuration for different severity levels",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							LogAlertConfigFieldWarning: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Warning severity threshold configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										ResourceFieldThresholdRuleStatic: {
											Type:        schema.TypeList,
											MinItems:    1,
											MaxItems:    1,
											Required:    true,
											Description: "Static threshold configuration",
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													ResourceFieldThresholdRuleStaticValue: {
														Type:        schema.TypeFloat,
														Required:    true,
														Description: "The static threshold value to compare against",
													},
												},
											},
										},
									},
								},
							},
							LogAlertConfigFieldCritical: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Critical severity threshold configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										ResourceFieldThresholdRuleStatic: {
											Type:        schema.TypeList,
											MinItems:    1,
											MaxItems:    1,
											Required:    true,
											Description: "Static threshold configuration",
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													ResourceFieldThresholdRuleStaticValue: {
														Type:        schema.TypeFloat,
														Required:    true,
														Description: "The static threshold value to compare against",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	logAlertConfigSchemaName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Name for the Log alert configuration",
		ValidateFunc: validation.StringLenBetween(0, 256),
	}

	logAlertConfigSchemaDescription = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The description text of the Log alert config",
		ValidateFunc: validation.StringLenBetween(0, 65536),
	}

	logAlertConfigSchemaGranularity = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      restapi.Granularity600000,
		ValidateFunc: validation.IntInSlice(restapi.SupportedSmartAlertGranularities.ToIntSlice()),
		Description:  "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
	}

	logAlertConfigSchemaTimeThreshold = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				LogAlertConfigFieldTimeThresholdViolationsInSequence: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Time threshold base on violations in sequence",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							LogAlertConfigFieldTimeThresholdTimeWindow: {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
)

var (
	logAlertConfigSchemaGracePeriod = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.",
	}
)

var logAlertConfigResourceSchema = map[string]*schema.Schema{
	LogAlertConfigFieldName:          logAlertConfigSchemaName,
	LogAlertConfigFieldDescription:   logAlertConfigSchemaDescription,
	LogAlertConfigFieldAlertChannels: logAlertConfigSchemaAlertChannels,
	LogAlertConfigFieldGracePeriod:   logAlertConfigSchemaGracePeriod,
	LogAlertConfigFieldGroupBy:       logAlertConfigSchemaGroupBy,
	LogAlertConfigFieldGranularity:   logAlertConfigSchemaGranularity,
	LogAlertConfigFieldTagFilter:     RequiredTagFilterExpressionSchema,
	LogAlertConfigFieldRules:         logAlertConfigSchemaRules,
	DefaultCustomPayloadFieldsName:   buildCustomPayloadFields(),
	LogAlertConfigFieldTimeThreshold: logAlertConfigSchemaTimeThreshold,
}

// LogAlertConfigResourceHandle creates the resource handle for Log Alert Configs
func LogAlertConfigResourceHandle() ResourceHandle[*restapi.LogAlertConfig] {
	return &logAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName:     ResourceInstanaLogAlertConfig,
			Schema:           logAlertConfigResourceSchema,
			SkipIDGeneration: true,
		},
	}
}

type logAlertConfigResource struct {
	metaData ResourceMetaData
}

func (c *logAlertConfigResource) MetaData() *ResourceMetaData {
	return &c.metaData
}

func (c *logAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (c *logAlertConfigResource) UpdateState(d *schema.ResourceData, config *restapi.LogAlertConfig) error {
	var normalizedTagFilterString *string
	var err error
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err = tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			return err
		}
	}

	d.SetId(config.ID)

	stateMap := map[string]interface{}{
		LogAlertConfigFieldName:          config.Name,
		LogAlertConfigFieldDescription:   config.Description,
		LogAlertConfigFieldTagFilter:     normalizedTagFilterString,
		LogAlertConfigFieldGroupBy:       c.mapGroupByToSchema(config),
		LogAlertConfigFieldAlertChannels: c.mapAlertChannelsToSchema(config),
		LogAlertConfigFieldGranularity:   config.Granularity,
		LogAlertConfigFieldTimeThreshold: c.mapTimeThresholdToSchema(config),
		DefaultCustomPayloadFieldsName:   mapCustomPayloadFieldsToSchema(config),
		LogAlertConfigFieldRules:         c.mapRulesToSchema(config),
	}

	if config.GracePeriod > 0 {
		stateMap[LogAlertConfigFieldGracePeriod] = config.GracePeriod
	}

	return tfutils.UpdateState(d, stateMap)
}

func (c *logAlertConfigResource) mapGroupByToSchema(config *restapi.LogAlertConfig) []map[string]interface{} {
	if len(config.GroupBy) > 0 {
		result := make([]map[string]interface{}, len(config.GroupBy))

		for i, groupBy := range config.GroupBy {
			groupByMap := make(map[string]interface{})
			groupByMap[LogAlertConfigFieldGroupByTagName] = groupBy.TagName
			if groupBy.Key != "" {
				groupByMap[LogAlertConfigFieldGroupByKey] = groupBy.Key
			}
			result[i] = groupByMap
		}
		return result
	}
	return []map[string]interface{}{}
}

func (c *logAlertConfigResource) mapAlertChannelsToSchema(config *restapi.LogAlertConfig) []map[string]interface{} {
	alertChannels := config.AlertChannels
	alertChannelsMap := make(map[string]interface{})

	if v, ok := alertChannels[restapi.WarningSeverity]; ok {
		alertChannelsMap[ResourceFieldThresholdRuleWarningSeverity] = v
	}

	if v, ok := alertChannels[restapi.CriticalSeverity]; ok {
		alertChannelsMap[ResourceFieldThresholdRuleCriticalSeverity] = v
	}

	// Only add to result if we have something
	if len(alertChannelsMap) > 0 {
		return []map[string]interface{}{alertChannelsMap}
	}

	// Otherwise, return an empty slice
	return []map[string]interface{}{}
}

func (c *logAlertConfigResource) mapTimeThresholdToSchema(config *restapi.LogAlertConfig) []map[string]interface{} {
	timeThresholdConfig := make(map[string]interface{})
	timeThresholdConfig[LogAlertConfigFieldTimeThresholdTimeWindow] = config.TimeThreshold.TimeWindow

	timeThresholdType := c.mapTimeThresholdTypeToSchema(config.TimeThreshold.Type)
	timeThreshold := make(map[string]interface{})
	timeThreshold[timeThresholdType] = []interface{}{timeThresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = timeThreshold

	return result
}

func (c *logAlertConfigResource) mapTimeThresholdTypeToSchema(input string) string {
	if input == "violationsInSequence" {
		return LogAlertConfigFieldTimeThresholdViolationsInSequence
	}

	return input
}

// rules mapping function
func (c *logAlertConfigResource) mapRulesToSchema(config *restapi.LogAlertConfig) []map[string]interface{} {
	if len(config.Rules) > 0 {
		ruleWithThreshold := config.Rules[0]
		rule := ruleWithThreshold.Rule

		// Convert "logCount" to "log.count" for the schema
		alertType := rule.AlertType
		if alertType == "logCount" {
			alertType = LogAlertTypeLogCount
		}

		ruleMap := map[string]interface{}{
			LogAlertConfigFieldMetricName:        rule.MetricName,
			LogAlertConfigFieldAlertType:         alertType,
			LogAlertConfigFieldThresholdOperator: ruleWithThreshold.ThresholdOperator,
		}

		if rule.Aggregation != "" {
			ruleMap[LogAlertConfigFieldAggregation] = rule.Aggregation
		}

		// Map thresholds
		thresholdMap := make(map[string]interface{})
		warningThreshold, isWarningThresholdPresent := ruleWithThreshold.Thresholds[restapi.WarningSeverity]
		criticalThreshold, isCriticalThresholdPresent := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]

		if isWarningThresholdPresent {
			// Create a static threshold structure for warning
			if warningThreshold.Type == "staticThreshold" && warningThreshold.Value != nil {
				warningMap := []map[string]interface{}{
					{
						ResourceFieldThresholdRuleStatic: []map[string]interface{}{
							{
								ResourceFieldThresholdRuleStaticValue: *warningThreshold.Value,
							},
						},
					},
				}
				thresholdMap[LogAlertConfigFieldWarning] = warningMap
			}
		}

		if isCriticalThresholdPresent {
			// Create a static threshold structure for critical
			if criticalThreshold.Type == "staticThreshold" && criticalThreshold.Value != nil {
				criticalMap := []map[string]interface{}{
					{
						ResourceFieldThresholdRuleStatic: []map[string]interface{}{
							{
								ResourceFieldThresholdRuleStaticValue: *criticalThreshold.Value,
							},
						},
					},
				}
				thresholdMap[LogAlertConfigFieldCritical] = criticalMap
			}
		}

		ruleMap[LogAlertConfigFieldThreshold] = []interface{}{thresholdMap}

		return []map[string]interface{}{ruleMap}
	}

	return []map[string]interface{}{}
}

func (c *logAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.LogAlertConfig, error) {
	var tagFilter *restapi.TagFilter
	var err error
	tagFilterStr, ok := d.GetOk(LogAlertConfigFieldTagFilter)
	if ok {
		tagFilter, err = c.mapTagFilterExpressionFromSchema(tagFilterStr.(string))
		if err != nil {
			return &restapi.LogAlertConfig{}, err
		}
	}

	customPayloadFields, err := mapDefaultCustomPayloadFieldsFromSchema(d)
	if err != nil {
		return &restapi.LogAlertConfig{}, err
	}

	config := &restapi.LogAlertConfig{
		ID:                    d.Id(),
		Name:                  d.Get(LogAlertConfigFieldName).(string),
		Description:           d.Get(LogAlertConfigFieldDescription).(string),
		TagFilterExpression:   tagFilter,
		GroupBy:               c.mapGroupByFromSchema(d),
		AlertChannels:         c.mapAlertChannelsFromSchema(d),
		Granularity:           restapi.Granularity(d.Get(LogAlertConfigFieldGranularity).(int)),
		TimeThreshold:         c.mapTimeThresholdFromSchema(d),
		CustomerPayloadFields: customPayloadFields,
		Rules:                 c.mapRuleFromSchema(d),
	}

	if v, ok := d.GetOk(LogAlertConfigFieldGracePeriod); ok {
		config.GracePeriod = int64(v.(int))
	}

	return config, nil
}

func (c *logAlertConfigResource) mapGroupByFromSchema(d *schema.ResourceData) []restapi.GroupByTag {
	groupBySlice, ok := d.Get(LogAlertConfigFieldGroupBy).([]interface{})
	if !ok || len(groupBySlice) == 0 {
		return []restapi.GroupByTag{}
	}

	result := make([]restapi.GroupByTag, len(groupBySlice))

	for i, groupByItem := range groupBySlice {
		groupBy := groupByItem.(map[string]interface{})
		tagName := groupBy[LogAlertConfigFieldGroupByTagName].(string)
		groupByTag := restapi.GroupByTag{
			TagName: tagName,
		}

		if key, ok := groupBy[LogAlertConfigFieldGroupByKey]; ok && key != nil {
			groupByTag.Key = key.(string)
		}

		result[i] = groupByTag
	}

	return result
}

func (c *logAlertConfigResource) mapTagFilterExpressionFromSchema(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (c *logAlertConfigResource) mapTimeThresholdFromSchema(d *schema.ResourceData) *restapi.LogTimeThreshold {
	timeThresholdSlice := d.Get(LogAlertConfigFieldTimeThreshold).([]interface{})
	timeThreshold := timeThresholdSlice[0].(map[string]interface{})

	for timeThresholdType, v := range timeThreshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})

			return &restapi.LogTimeThreshold{
				Type:       c.mapTimeThresholdTypeFromSchema(timeThresholdType),
				TimeWindow: int64(config[LogAlertConfigFieldTimeThresholdTimeWindow].(int)),
			}
		}
	}
	return &restapi.LogTimeThreshold{}
}

func (c *logAlertConfigResource) mapTimeThresholdTypeFromSchema(input string) string {
	if input == LogAlertConfigFieldTimeThresholdViolationsInSequence {
		return "violationsInSequence"
	}

	return input
}

// Rule mapping function
func (c *logAlertConfigResource) mapRuleFromSchema(d *schema.ResourceData) []restapi.RuleWithThreshold[restapi.LogAlertRule] {
	ruleSlice := d.Get(LogAlertConfigFieldRules).([]interface{})
	if len(ruleSlice) == 0 {
		return []restapi.RuleWithThreshold[restapi.LogAlertRule]{}
	}

	// Get the rule configuration
	ruleConfig := ruleSlice[0].(map[string]interface{})

	// Create the LogAlertRule
	alertType := ruleConfig[LogAlertConfigFieldAlertType].(string)

	// Convert "log.count" to "logCount" for the API
	if alertType == LogAlertTypeLogCount {
		alertType = "logCount"
	}

	logAlertRule := restapi.LogAlertRule{
		AlertType:  alertType,
		MetricName: ruleConfig[LogAlertConfigFieldMetricName].(string),
	}

	if aggregation, ok := ruleConfig[LogAlertConfigFieldAggregation]; ok && aggregation != nil {
		logAlertRule.Aggregation = restapi.Aggregation(aggregation.(string))
	}

	// Get threshold operator
	thresholdOperator := restapi.ThresholdOperator(ruleConfig[LogAlertConfigFieldThresholdOperator].(string))

	// Map thresholds
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	// Get threshold configuration
	thresholdSlice := ruleConfig[LogAlertConfigFieldThreshold].([]interface{})
	if len(thresholdSlice) > 0 {
		thresholdConfig := thresholdSlice[0].(map[string]interface{})

		// Map warning threshold if present
		if warningSlice, ok := thresholdConfig[LogAlertConfigFieldWarning].([]interface{}); ok && len(warningSlice) > 0 {
			warningConfig := warningSlice[0].(map[string]interface{})
			if staticSlice, ok := warningConfig[ResourceFieldThresholdRuleStatic].([]interface{}); ok && len(staticSlice) > 0 {
				staticConfig := staticSlice[0].(map[string]interface{})
				value := staticConfig[ResourceFieldThresholdRuleStaticValue].(float64)
				valuePtr := &value
				thresholdMap[restapi.WarningSeverity] = restapi.ThresholdRule{
					Type:  "staticThreshold",
					Value: valuePtr,
				}
			}
		}

		// Map critical threshold if present
		if criticalSlice, ok := thresholdConfig[LogAlertConfigFieldCritical].([]interface{}); ok && len(criticalSlice) > 0 {
			criticalConfig := criticalSlice[0].(map[string]interface{})
			if staticSlice, ok := criticalConfig[ResourceFieldThresholdRuleStatic].([]interface{}); ok && len(staticSlice) > 0 {
				staticConfig := staticSlice[0].(map[string]interface{})
				value := staticConfig[ResourceFieldThresholdRuleStaticValue].(float64)
				valuePtr := &value
				thresholdMap[restapi.CriticalSeverity] = restapi.ThresholdRule{
					Type:  "staticThreshold",
					Value: valuePtr,
				}
			}
		}
	}

	// Create the RuleWithThreshold
	ruleWithThreshold := restapi.RuleWithThreshold[restapi.LogAlertRule]{
		ThresholdOperator: thresholdOperator,
		Rule:              logAlertRule,
		Thresholds:        thresholdMap,
	}

	return []restapi.RuleWithThreshold[restapi.LogAlertRule]{ruleWithThreshold}
}

func (c *logAlertConfigResource) mapAlertChannelsFromSchema(d *schema.ResourceData) map[restapi.AlertSeverity][]string {
	alertChannelsMap := make(map[restapi.AlertSeverity][]string)
	alertChannelsSlice, ok := d.Get(LogAlertConfigFieldAlertChannels).([]interface{})

	if !ok || len(alertChannelsSlice) == 0 {
		// no alert channels defined
		return alertChannelsMap
	}

	alertChannels := alertChannelsSlice[0].(map[string]interface{})
	if val, ok := alertChannels[ResourceFieldThresholdRuleWarningSeverity]; ok && val != nil {
		if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			alertChannelsMap[restapi.WarningSeverity] = ReadArrayParameterFromMap[string](alertChannels, ResourceFieldThresholdRuleWarningSeverity)
		}
	}
	if val, ok := alertChannels[ResourceFieldThresholdRuleCriticalSeverity]; ok && val != nil {
		if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			alertChannelsMap[restapi.CriticalSeverity] = ReadArrayParameterFromMap[string](alertChannels, ResourceFieldThresholdRuleCriticalSeverity)
		}
	}

	return alertChannelsMap
}

func (c *logAlertConfigResource) SetComputedFields(d *schema.ResourceData) error {
	return nil
}

func (c *logAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.LogAlertConfig] {
	return api.LogAlertConfig()
}
