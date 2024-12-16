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
	thresholdConfig := make(map[string]interface{})

	if threshold.Value != nil {
		thresholdConfig[ResourceFieldThresholdRuleStaticValue] = *threshold.Value
	}
	if threshold.Baseline != nil {
		thresholdConfig[ResourceFieldThresholdRuleHistoricBaselineBaseline] = *threshold.Baseline
	}
	if threshold.DeviationFactor != nil {
		thresholdConfig[ResourceFieldThresholdRuleHistoricBaselineDeviationFactor] = float64(*threshold.DeviationFactor)
	}
	if threshold.Seasonality != nil {
		thresholdConfig[ResourceFieldThresholdRuleHistoricBaselineSeasonality] = *threshold.Seasonality
	}

	thresholdType := t.mapThresholdTypeToSchema(threshold.Type)
	thresholdRule := make(map[string]interface{})
	thresholdRule[thresholdType] = []interface{}{thresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = thresholdRule

	return result
}

func (t *thresholdRuleMapperImpl) mapThresholdTypeToSchema(input string) string {
	if input == "historicBaseline" {
		return ResourceFieldThresholdRuleHistoricBaseline
	} else if input == "staticThreshold" {
		return ResourceFieldThresholdRuleStatic
	}
	return input
}

func (t *thresholdRuleMapperImpl) mapThresholdTypeFromSchema(input string) string {
	if input == ResourceFieldThresholdRuleHistoricBaseline {
		return "historicBaseline"
	} else if input == ResourceFieldThresholdRuleStatic {
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
	if v, ok := config[ResourceFieldThresholdHistoricBaselineSeasonality]; ok {
		seasonality := restapi.ThresholdSeasonality(v.(string))
		seasonalityPtr = &seasonality
	}
	var valuePtr *float64
	if v, ok := config[ResourceFieldThresholdStaticValue]; ok {
		value := v.(float64)
		valuePtr = &value
	}
	var deviationFactorPtr *float32
	if v, ok := config[ResourceFieldThresholdHistoricBaselineDeviationFactor]; ok {
		deviationFactor := float32(v.(float64))
		deviationFactorPtr = &deviationFactor
	}
	var baselinePtr *[][]float64
	if v, ok := config[ResourceFieldThresholdHistoricBaselineBaseline]; ok {
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
