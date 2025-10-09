package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	//ResourceFieldThresholdRule constant value for field threshold
	ResourceFieldThresholdRule = "threshold"
	//ResourceFieldThresholdRuleHistoricBaseline constant value for field threshold.historic_baseline
	ResourceFieldThresholdRuleHistoricBaseline = "historic_baseline"
	//ResourceFieldThresholdRuleHistoricBaselineBaseline constant value for field threshold.historic_baseline.baseline
	ResourceFieldThresholdRuleHistoricBaselineBaseline = "baseline"
	//ResourceFieldThresholdRuleHistoricBaselineDeviationFactor constant value for field threshold.historic_baseline.deviation_factor
	ResourceFieldThresholdRuleHistoricBaselineDeviationFactor = "deviation_factor"
	//ResourceFieldThresholdRuleHistoricBaselineSeasonality constant value for field threshold.historic_baseline.seasonality
	ResourceFieldThresholdRuleHistoricBaselineSeasonality = "seasonality"
	//ResourceFieldThresholdRuleStatic constant value for field threshold.static
	ResourceFieldThresholdRuleStatic = "static"
	//ResourceFieldThresholdRuleStaticValue constant value for field threshold.static.value
	ResourceFieldThresholdRuleStaticValue = "value"
)

var thresholdRuleSchema = &schema.Schema{
	Type:        schema.TypeList,
	MinItems:    0,
	MaxItems:    1,
	Optional:    true,
	Description: "Indicates the type of threshold this alert rule is evaluated on.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			ResourceFieldThresholdRuleHistoricBaseline: {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "Threshold based on a historic baseline.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ResourceFieldThresholdRuleHistoricBaselineBaseline: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:     schema.TypeSet,
								Optional: false,
								Elem: &schema.Schema{
									Type: schema.TypeFloat,
								},
							},
							Description: "The baseline of the historic baseline threshold",
						},
						ResourceFieldThresholdRuleHistoricBaselineDeviationFactor: {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatBetween(0.5, 16),
							Description:  "The baseline of the historic baseline threshold",
						},
						ResourceFieldThresholdRuleHistoricBaselineSeasonality: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdSeasonalities.ToStringSlice(), false),
							Description:  "The seasonality of the historic baseline threshold",
						},
					},
				},
			},
			ResourceFieldThresholdRuleStatic: {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "Static threshold definition",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ResourceFieldThresholdRuleStaticValue: {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "The value of the static threshold",
						},
					},
				},
			},
		},
	},
}

func newThresholdRuleMapper() thresholdRuleMapper {
	return &thresholdRuleMapperImpl{}
}

type thresholdRuleMapper interface {
	toState(threshold *restapi.ThresholdRule) []map[string]interface{}
	fromState(thresholdSlice []interface{}) *restapi.ThresholdRule
}

type thresholdRuleMapperImpl struct{}

func (t thresholdRuleMapperImpl) toState(threshold *restapi.ThresholdRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	thresholdRule := make(map[string]interface{})
	result[0] = thresholdRule

	switch threshold.Type {
	case "staticThreshold":
		staticConfig := make(map[string]interface{})
		if threshold.Value != nil {
			staticConfig[ResourceFieldThresholdRuleStaticValue] = *threshold.Value
		}
		thresholdRule[ResourceFieldThresholdRuleStatic] = []interface{}{staticConfig}
	case "historicBaseline":
		historicConfig := make(map[string]interface{})
		if threshold.Baseline != nil {
			historicConfig[ResourceFieldThresholdRuleHistoricBaselineBaseline] = *threshold.Baseline
		}
		if threshold.DeviationFactor != nil {
			historicConfig[ResourceFieldThresholdRuleHistoricBaselineDeviationFactor] = float64(*threshold.DeviationFactor)
		}
		if threshold.Seasonality != nil {
			historicConfig[ResourceFieldThresholdRuleHistoricBaselineSeasonality] = *threshold.Seasonality
		}
		thresholdRule[ResourceFieldThresholdRuleHistoricBaseline] = []interface{}{historicConfig}
	}

	return result
}

func (t *thresholdRuleMapperImpl) mapThresholdTypeFromSchema(input string) string {
	switch input {
	case ResourceFieldThresholdRuleHistoricBaseline:
		return "historicBaseline"
	case ResourceFieldThresholdRuleStatic:
		return "staticThreshold"
	}
	return input
}

func (t *thresholdRuleMapperImpl) fromState(thresholdSlice []interface{}) *restapi.ThresholdRule {
	threshold := thresholdSlice[0].(map[string]interface{})

	for thresholdType, v := range threshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			return t.mapThresholdConfigFromSchema(config, thresholdType)
		}
	}

	return &restapi.ThresholdRule{}
}

func (t *thresholdRuleMapperImpl) mapThresholdConfigFromSchema(config map[string]interface{}, thresholdType string) *restapi.ThresholdRule {
	var seasonalityPtr *restapi.ThresholdSeasonality
	var valuePtr *float64
	var deviationFactorPtr *float32
	var baselinePtr *[][]float64

	// Handle static threshold
	if v, ok := config[ResourceFieldThresholdRuleStaticValue]; ok {
		value := v.(float64)
		valuePtr = &value
	}

	// Handle historic baseline
	if v, ok := config[ResourceFieldThresholdRuleHistoricBaselineSeasonality]; ok {
		seasonality := restapi.ThresholdSeasonality(v.(string))
		seasonalityPtr = &seasonality
	}
	if v, ok := config[ResourceFieldThresholdRuleHistoricBaselineDeviationFactor]; ok {
		deviationFactor := float32(v.(float64))
		deviationFactorPtr = &deviationFactor
	}
	if v, ok := config[ResourceFieldThresholdRuleHistoricBaselineBaseline]; ok {
		baselineSet := v.(*schema.Set)
		if baselineSet.Len() > 0 {
			baseline := make([][]float64, baselineSet.Len())
			for i, val := range baselineSet.List() {
				baseline[i] = ConvertInterfaceSlice[float64](val.(*schema.Set).List())
			}
			baselinePtr = &baseline
		}
	}

	return &restapi.ThresholdRule{
		Type:            t.mapThresholdTypeFromSchema(thresholdType),
		Value:           valuePtr,
		DeviationFactor: deviationFactorPtr,
		Baseline:        baselinePtr,
		Seasonality:     seasonalityPtr,
	}
}
