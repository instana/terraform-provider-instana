package instana

import (
	"context"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const ResourceInstanaInfraAlertConfig = "instana_infra_alert_config"

const (
	InfraAlertConfigFieldName                  = "name"
	InfraAlertConfigFieldFullName              = "full_name"
	InfraAlertConfigFieldDescription           = "description"
	InfraAlertConfigFieldAlertChannels         = "alert_channels"
	ResourceFieldThresholdRuleWarningSeverity  = "warning"
	ResourceFieldThresholdRuleCriticalSeverity = "critical"
	InfraAlertConfigFieldGroupBy               = "group_by"
	InfraAlertConfigFieldGranularity           = "granularity"
	InfraAlertConfigFieldTagFilter             = "tag_filter"
	InfraAlertConfigFieldEvaluationType        = "evaluation_type"

	InfraAlertConfigFieldRules       = "rules"
	InfraAlertConfigFieldGenericRule = "generic_rule"

	InfraAlertConfigFieldMetricName             = "metric_name"
	InfraAlertConfigFieldEntityType             = "entity_type"
	InfraAlertConfigFieldAggregation            = "aggregation"
	InfraAlertConfigFieldCrossSeriesAggregation = "cross_series_aggregation"
	InfraAlertConfigFieldRegex                  = "regex"
	InfraAlertConfigFieldThresholdOperator      = "threshold_operator"

	InfraAlertConfigFieldTimeThreshold                     = "time_threshold"
	InfraAlertConfigFieldTimeThresholdTimeWindow           = "time_window"
	InfraAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
)

var (
	infraAlertConfigSchemaAlertChannels = &schema.Schema{
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
	infraAlertConfigSchemaGroupBy = &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 0,
		MaxItems: 5,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "The grouping tags used to group the metric results.",
	}
	infraAlertConfigSchemaEvaluationType = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  string(restapi.EvaluationTypeCustom),
		ValidateFunc: validation.StringInSlice(
			restapi.SupportedInfraAlertEvaluationTypes.ToStringSlice(), false,
		),
		Description: "Defines how the alert is evaluated. Possible values: 'PER_ENTITY', 'CUSTOM'. Default is 'CUSTOM'.",
	}
	infraAlertConfigSchemaRules = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels. This enables more complex alert configurations with validations to ensure consistent and logical threshold-severity combinations.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				InfraAlertConfigFieldGenericRule: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							InfraAlertConfigFieldMetricName: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The metric name of the infrastructure alert rule",
							},
							InfraAlertConfigFieldEntityType: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The entity type of the infrastructure alert rule",
							},
							InfraAlertConfigFieldAggregation: {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(restapi.InfraExploreAggregations.ToStringSlice(), false),
							},
							InfraAlertConfigFieldCrossSeriesAggregation: {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(restapi.TimeAggregations.ToStringSlice(), false),
							},
							InfraAlertConfigFieldRegex: {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Indicates if the given metric name follows regex pattern or not",
							},
							InfraAlertConfigFieldThresholdOperator: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The operator which will be applied to evaluate the threshold",
								ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdOperators.ToStringSlice(), true),
							},

							ResourceFieldThresholdRule: {
								Type:        schema.TypeList,
								MinItems:    1,
								MaxItems:    1,
								Required:    true,
								Description: "",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										ResourceFieldThresholdRuleWarningSeverity:  thresholdRuleSchema,
										ResourceFieldThresholdRuleCriticalSeverity: thresholdRuleSchema,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	infraAlertConfigSchemaName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Name for the Infrastructure alert configuration",
		ValidateFunc: validation.StringLenBetween(0, 256),
	}
	infraAlertConfigSchemaDescription = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The description text of the Infrastructure alert config",
		ValidateFunc: validation.StringLenBetween(0, 65536),
	}
	infraAlertConfigSchemaGranularity = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      restapi.Granularity600000,
		ValidateFunc: validation.IntInSlice(restapi.SupportedSmartAlertGranularities.ToIntSlice()),
		Description:  "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
	}
	infraAlertConfigSchemaTimeThreshold = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				InfraAlertConfigFieldTimeThresholdViolationsInSequence: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Time threshold base on violations in sequence",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							InfraAlertConfigFieldTimeThresholdTimeWindow: {
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

var infraAlertConfigResourceSchema = map[string]*schema.Schema{
	InfraAlertConfigFieldName:           infraAlertConfigSchemaName,
	InfraAlertConfigFieldDescription:    infraAlertConfigSchemaDescription,
	InfraAlertConfigFieldAlertChannels:  infraAlertConfigSchemaAlertChannels,
	InfraAlertConfigFieldGroupBy:        infraAlertConfigSchemaGroupBy,
	InfraAlertConfigFieldGranularity:    infraAlertConfigSchemaGranularity,
	InfraAlertConfigFieldTagFilter:      OptionalTagFilterExpressionSchema,
	InfraAlertConfigFieldRules:          infraAlertConfigSchemaRules,
	DefaultCustomPayloadFieldsName:      buildCustomPayloadFields(),
	InfraAlertConfigFieldTimeThreshold:  infraAlertConfigSchemaTimeThreshold,
	InfraAlertConfigFieldEvaluationType: infraAlertConfigSchemaEvaluationType,
}

func (c *infraAlertConfigResource) stateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[InfraAlertConfigFieldFullName]; ok {
		state[InfraAlertConfigFieldName] = state[InfraAlertConfigFieldFullName]
		delete(state, InfraAlertConfigFieldFullName)
	}
	return state, nil
}

func (c *infraAlertConfigResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			InfraAlertConfigFieldName:          infraAlertConfigSchemaName,
			InfraAlertConfigFieldDescription:   infraAlertConfigSchemaDescription,
			InfraAlertConfigFieldAlertChannels: infraAlertConfigSchemaAlertChannels,
			InfraAlertConfigFieldGranularity:   infraAlertConfigSchemaGranularity,
			InfraAlertConfigFieldTagFilter:     OptionalTagFilterExpressionSchema,
			InfraAlertConfigFieldRules:         infraAlertConfigSchemaRules,
			DefaultCustomPayloadFieldsName:     buildCustomPayloadFields(),
			InfraAlertConfigFieldTimeThreshold: infraAlertConfigSchemaTimeThreshold,
		},
	}
}

// NewInfraAlertConfigResourceHandle creates the resource handle for Website Alert Configs
func NewInfraAlertConfigResourceHandle() ResourceHandle[*restapi.InfraAlertConfig] {
	return &infraAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName:     ResourceInstanaInfraAlertConfig,
			Schema:           infraAlertConfigResourceSchema,
			SkipIDGeneration: true,
			SchemaVersion:    1,
		},
	}
}

type infraAlertConfigResource struct {
	metaData ResourceMetaData
}

func (c *infraAlertConfigResource) MetaData() *ResourceMetaData {
	return &c.metaData
}

func (c *infraAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    c.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: c.stateUpgradeV0,
			Version: 0,
		},
	}
}

func (c *infraAlertConfigResource) UpdateState(d *schema.ResourceData, config *restapi.InfraAlertConfig) error {
	var normalizedTagFilterString *string
	var err error
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err = tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			return err
		}
	}

	d.SetId(config.ID)

	return tfutils.UpdateState(d, map[string]interface{}{
		InfraAlertConfigFieldName:           config.Name,
		InfraAlertConfigFieldDescription:    config.Description,
		InfraAlertConfigFieldTagFilter:      normalizedTagFilterString,
		InfraAlertConfigFieldGroupBy:        config.GroupBy,
		InfraAlertConfigFieldAlertChannels:  c.mapAlertChannelsToSchema(config),
		InfraAlertConfigFieldGranularity:    config.Granularity,
		InfraAlertConfigFieldTimeThreshold:  c.mapTimeThresholdToSchema(config),
		DefaultCustomPayloadFieldsName:      mapCustomPayloadFieldsToSchema(config),
		InfraAlertConfigFieldRules:          c.mapRulesToSchema(config),
		InfraAlertConfigFieldEvaluationType: string(config.EvaluationType),
	})
}

func (c *infraAlertConfigResource) mapAlertChannelsToSchema(config *restapi.InfraAlertConfig) []map[string]interface{} {
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

func (c *infraAlertConfigResource) mapTimeThresholdToSchema(config *restapi.InfraAlertConfig) []map[string]interface{} {
	timeThresholdConfig := make(map[string]interface{})
	timeThresholdConfig[InfraAlertConfigFieldTimeThresholdTimeWindow] = config.TimeThreshold.TimeWindow

	timeThresholdType := c.mapTimeThresholdTypeToSchema(config.TimeThreshold.Type)
	timeThreshold := make(map[string]interface{})
	timeThreshold[timeThresholdType] = []interface{}{timeThresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = timeThreshold

	return result
}

func (c *infraAlertConfigResource) mapTimeThresholdTypeToSchema(input string) string {
	if input == "violationsInSequence" {
		return InfraAlertConfigFieldTimeThresholdViolationsInSequence
	}

	return input
}

func (c *infraAlertConfigResource) mapRulesToSchema(config *restapi.InfraAlertConfig) []map[string]interface{} {
	if len(config.Rules) > 0 {
		firstRule := config.Rules[0].Rule

		if firstRule.AlertType == "genericRule" {
			rule := make(map[string]interface{})
			ruleAttribute := c.mapGenericRuleToSchema(&config.Rules[0])
			rule[InfraAlertConfigFieldGenericRule] = []interface{}{ruleAttribute}
			result := make([]map[string]interface{}, 1)
			result[0] = rule

			return result
		}
	}

	return []map[string]interface{}{}
}

func (c *infraAlertConfigResource) mapGenericRuleToSchema(ruleWithThreshold *restapi.RuleWithThreshold[restapi.InfraAlertRule]) map[string]interface{} {
	var rule = ruleWithThreshold.Rule

	return map[string]interface{}{
		InfraAlertConfigFieldMetricName:             rule.MetricName,
		InfraAlertConfigFieldEntityType:             rule.EntityType,
		InfraAlertConfigFieldAggregation:            rule.Aggregation,
		InfraAlertConfigFieldCrossSeriesAggregation: rule.CrossSeriesAggregation,
		InfraAlertConfigFieldRegex:                  rule.Regex,
		InfraAlertConfigFieldThresholdOperator:      ruleWithThreshold.ThresholdOperator,
		ResourceFieldThresholdRule:                  c.mapThresholdToSchema(ruleWithThreshold),
	}
}

func (c *infraAlertConfigResource) mapThresholdToSchema(ruleWithThreshold *restapi.RuleWithThreshold[restapi.InfraAlertRule]) []map[string]interface{} {
	var thresholds = ruleWithThreshold.Thresholds
	warningThresholdRule, isWarningThresholdPresent := thresholds[restapi.WarningSeverity]
	criticalThresholdRule, isCriticalThresholdPresent := thresholds[restapi.CriticalSeverity]

	result := make([]map[string]interface{}, 1)
	thresholdMap := make(map[string]interface{})

	if isWarningThresholdPresent {
		thresholdMap[ResourceFieldThresholdRuleWarningSeverity] = newThresholdRuleMapper().toState(&warningThresholdRule)
	}
	if isCriticalThresholdPresent {
		thresholdMap[ResourceFieldThresholdRuleCriticalSeverity] = newThresholdRuleMapper().toState(&criticalThresholdRule)
	}

	result[0] = thresholdMap

	return result
}

func (c *infraAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.InfraAlertConfig, error) {
	var tagFilter *restapi.TagFilter
	var err error
	tagFilterStr, ok := d.GetOk(InfraAlertConfigFieldTagFilter)
	if ok {
		tagFilter, err = c.mapTagFilterExpressionFromSchema(tagFilterStr.(string))
		if err != nil {
			return &restapi.InfraAlertConfig{}, err
		}
	}

	customPayloadFields, err := mapDefaultCustomPayloadFieldsFromSchema(d)
	if err != nil {
		return &restapi.InfraAlertConfig{}, err
	}

	evaluationTypeStr, ok := d.Get(InfraAlertConfigFieldEvaluationType).(string)
	if !ok || evaluationTypeStr == "" {
		evaluationTypeStr = "CUSTOM"
	}

	return &restapi.InfraAlertConfig{
		ID:                    d.Id(),
		Name:                  d.Get(InfraAlertConfigFieldName).(string),
		Description:           d.Get(InfraAlertConfigFieldDescription).(string),
		TagFilterExpression:   tagFilter,
		GroupBy:               ReadArrayParameterFromResource[string](d, InfraAlertConfigFieldGroupBy),
		AlertChannels:         c.mapAlertChannelsFromSchema(d),
		Granularity:           restapi.Granularity(d.Get(InfraAlertConfigFieldGranularity).(int)),
		TimeThreshold:         c.mapTimeThresholdFromSchema(d),
		CustomerPayloadFields: customPayloadFields,
		Rules:                 c.mapRuleFromSchema(d),
		EvaluationType:        restapi.InfraAlertEvaluationType(evaluationTypeStr),
	}, nil
}

func (c *infraAlertConfigResource) mapTagFilterExpressionFromSchema(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (c *infraAlertConfigResource) mapTimeThresholdFromSchema(d *schema.ResourceData) restapi.InfraTimeThreshold {
	timeThresholdSlice := d.Get(InfraAlertConfigFieldTimeThreshold).([]interface{})
	timeThreshold := timeThresholdSlice[0].(map[string]interface{})

	for timeThresholdType, v := range timeThreshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})

			return restapi.InfraTimeThreshold{
				Type:       c.mapTimeThresholdTypeFromSchema(timeThresholdType),
				TimeWindow: int64(config[InfraAlertConfigFieldTimeThresholdTimeWindow].(int)),
			}
		}
	}
	return restapi.InfraTimeThreshold{}
}

func (c *infraAlertConfigResource) mapTimeThresholdTypeFromSchema(input string) string {
	if input == InfraAlertConfigFieldTimeThresholdViolationsInSequence {
		return "violationsInSequence"
	}

	return input
}

func (c *infraAlertConfigResource) mapRuleFromSchema(d *schema.ResourceData) []restapi.RuleWithThreshold[restapi.InfraAlertRule] {
	ruleSlice := d.Get(InfraAlertConfigFieldRules).([]interface{})
	// Only one rule definition is allowed for now.
	rule := ruleSlice[0].(map[string]interface{})

	for alertType, v := range rule {
		configSlice := v.([]interface{})
		if alertType == InfraAlertConfigFieldGenericRule && len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			return c.mapRuleWithThresholdFromSchema(config, alertType)
		}
	}

	return []restapi.RuleWithThreshold[restapi.InfraAlertRule]{}
}

func (c *infraAlertConfigResource) mapAlertChannelsFromSchema(d *schema.ResourceData) map[restapi.AlertSeverity][]string {
	alertChannelsMap := make(map[restapi.AlertSeverity][]string)
	alertChannelsSlice, ok := d.Get(InfraAlertConfigFieldAlertChannels).([]interface{})

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

func (c *infraAlertConfigResource) mapAlertTypeFromSchema(alertType string) string {
	if alertType == InfraAlertConfigFieldGenericRule {
		return "genericRule"
	}

	return alertType
}

func (c *infraAlertConfigResource) mapRuleWithThresholdFromSchema(config map[string]interface{}, alertType string) []restapi.RuleWithThreshold[restapi.InfraAlertRule] {
	result := make([]restapi.RuleWithThreshold[restapi.InfraAlertRule], 1)

	infraAlertRule := restapi.InfraAlertRule{
		AlertType:              c.mapAlertTypeFromSchema(alertType),
		MetricName:             config[InfraAlertConfigFieldMetricName].(string),
		EntityType:             config[InfraAlertConfigFieldEntityType].(string),
		Aggregation:            restapi.Aggregation(config[InfraAlertConfigFieldAggregation].(string)),
		CrossSeriesAggregation: restapi.Aggregation(config[InfraAlertConfigFieldCrossSeriesAggregation].(string)),
		Regex:                  config[InfraAlertConfigFieldRegex].(bool),
	}

	thresholdSlice := config[ResourceFieldThresholdRule].([]interface{})
	thresholdConfig := thresholdSlice[0].(map[string]interface{})
	thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

	if v, ok := thresholdConfig[ResourceFieldThresholdRuleWarningSeverity]; ok && len(v.([]interface{})) == 1 {
		warningThresholdSlice := v.([]interface{})
		thresholdRule := newThresholdRuleMapper().fromState(warningThresholdSlice)
		thresholdMap[restapi.WarningSeverity] = *thresholdRule
	}
	if v, ok := thresholdConfig[ResourceFieldThresholdRuleCriticalSeverity]; ok && len(v.([]interface{})) == 1 {
		criticalThresholdSlice := v.([]interface{})
		thresholdRule := newThresholdRuleMapper().fromState(criticalThresholdSlice)
		thresholdMap[restapi.CriticalSeverity] = *thresholdRule
	}

	ruleWithThreshold := restapi.RuleWithThreshold[restapi.InfraAlertRule]{
		ThresholdOperator: restapi.ThresholdOperator(config[InfraAlertConfigFieldThresholdOperator].(string)),
		Rule:              infraAlertRule,
		Thresholds:        thresholdMap,
	}
	result[0] = ruleWithThreshold

	return result
}

func (c *infraAlertConfigResource) SetComputedFields(d *schema.ResourceData) error {
	return nil
}

func (c *infraAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.InfraAlertConfig] {
	return api.InfraAlertConfig()
}
