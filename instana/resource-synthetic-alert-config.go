package instana

// import (
// 	"context"

// 	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
// 	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
// 	"github.com/gessnerfl/terraform-provider-instana/tfutils"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
// )

// // ResourceInstanaSyntheticAlertConfig the name of the terraform-provider-instana resource to manage synthetic alert configurations
// const ResourceInstanaSyntheticAlertConfig = "instana_synthetic_alert_config"

// const (
// 	// SyntheticAlertConfigFieldName constant value for the schema field name
// 	SyntheticAlertConfigFieldName = "name"
// 	// SyntheticAlertConfigFieldName constant value for the schema field full name
// 	SyntheticAlertConfigFieldFullName = "full_name"
// 	// SyntheticAlertConfigFieldDescription constant value for the schema field description
// 	SyntheticAlertConfigFieldDescription = "description"
// 	// SyntheticAlertConfigFieldSyntheticTestIds constant value for the schema field synthetic_test_ids
// 	SyntheticAlertConfigFieldSyntheticTestIds = "synthetic_test_ids"
// 	// SyntheticAlertConfigFieldSeverity constant value for the schema field severity
// 	SyntheticAlertConfigFieldSeverity = "severity"
// 	// SyntheticAlertConfigFieldTagFilter constant value for the schema field tag_filter
// 	SyntheticAlertConfigFieldTagFilter = "tag_filter"
// 	// SyntheticAlertConfigFieldRule constant value for the schema field rule
// 	SyntheticAlertConfigFieldRule = "rule"
// 	// SyntheticAlertConfigFieldAlertChannelIds constant value for the schema field alert_channel_ids
// 	SyntheticAlertConfigFieldAlertChannelIds = "alert_channel_ids"
// 	// SyntheticAlertConfigFieldTimeThreshold constant value for the schema field time_threshold
// 	SyntheticAlertConfigFieldTimeThreshold = "time_threshold"
// 	// SyntheticAlertConfigFieldGracePeriod constant value for the schema field grace_period
// 	SyntheticAlertConfigFieldGracePeriod = "grace_period"

// 	// Rule fields
// 	SyntheticAlertRuleFieldAlertType   = "alert_type"
// 	SyntheticAlertRuleFieldMetricName  = "metric_name"
// 	SyntheticAlertRuleFieldAggregation = "aggregation"

// 	// TimeThreshold fields
// 	SyntheticAlertTimeThresholdFieldType            = "type"
// 	SyntheticAlertTimeThresholdFieldViolationsCount = "violations_count"
// )

// // SyntheticAlertConfigSchemaName schema field definition of instana_synthetic_alert_config field name
// var SyntheticAlertConfigSchemaName = &schema.Schema{
// 	Type:         schema.TypeString,
// 	Required:     true,
// 	Description:  "Configures the name of the synthetic alert configuration",
// 	ValidateFunc: validation.StringLenBetween(1, 256),
// }

// // SyntheticAlertConfigSchemaDescription schema field definition of instana_synthetic_alert_config field description
// var SyntheticAlertConfigSchemaDescription = &schema.Schema{
// 	Type:         schema.TypeString,
// 	Required:     true,
// 	Description:  "Configures the description of the synthetic alert configuration",
// 	ValidateFunc: validation.StringLenBetween(0, 1024),
// }

// // SyntheticAlertConfigSchemaSyntheticTestIds schema field definition of instana_synthetic_alert_config field synthetic_test_ids
// var SyntheticAlertConfigSchemaSyntheticTestIds = &schema.Schema{
// 	Type:     schema.TypeSet,
// 	MinItems: 0,
// 	MaxItems: 1024,
// 	Elem: &schema.Schema{
// 		Type: schema.TypeString,
// 	},
// 	Required:    true,
// 	Description: "Configures the list of Synthetic Test IDs to monitor.",
// }

// // SyntheticAlertConfigSchemaSeverity schema field definition of instana_synthetic_alert_config field severity
// var SyntheticAlertConfigSchemaSeverity = &schema.Schema{
// 	Type:         schema.TypeInt,
// 	Optional:     true,
// 	Description:  "Configures the severity of the alert (5=critical, 10=warning)",
// 	ValidateFunc: validation.IntInSlice([]int{5, 10}),
// }

// // SyntheticAlertConfigSchemaAlertChannelIds schema field definition of instana_synthetic_alert_config field alert_channel_ids
// var SyntheticAlertConfigSchemaAlertChannelIds = &schema.Schema{
// 	Type:     schema.TypeSet,
// 	MinItems: 0,
// 	MaxItems: 1024,
// 	Elem: &schema.Schema{
// 		Type: schema.TypeString,
// 	},
// 	Required:    true,
// 	Description: "Configures the list of Alert Channel IDs.",
// }

// // SyntheticAlertConfigSchemaGracePeriod schema field definition of instana_synthetic_alert_config field grace_period
// var SyntheticAlertConfigSchemaGracePeriod = &schema.Schema{
// 	Type:        schema.TypeInt,
// 	Optional:    true,
// 	Description: "The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires",
// }

// // SyntheticAlertConfigSchemaRule schema field definition of instana_synthetic_alert_config field rule
// var SyntheticAlertConfigSchemaRule = &schema.Schema{
// 	Type:        schema.TypeList,
// 	Required:    true,
// 	MaxItems:    1,
// 	Description: "Configures the rule for the synthetic alert",
// 	Elem: &schema.Resource{
// 		Schema: map[string]*schema.Schema{
// 			SyntheticAlertRuleFieldAlertType: {
// 				Type:         schema.TypeString,
// 				Required:     true,
// 				Description:  "The type of the alert rule (e.g., failure)",
// 				ValidateFunc: validation.StringInSlice([]string{"failure"}, false),
// 			},
// 			SyntheticAlertRuleFieldMetricName: {
// 				Type:         schema.TypeString,
// 				Required:     true,
// 				Description:  "The metric name to monitor (e.g., status)",
// 				ValidateFunc: validation.StringLenBetween(1, 256),
// 			},
// 			SyntheticAlertRuleFieldAggregation: {
// 				Type:         schema.TypeString,
// 				Optional:     true,
// 				Description:  "The aggregation method {SUM,MEAN,MAX,MIN,P25,P50,P75,P90,P95,P98,P99,P99_9,P99_99,DISTINCT_COUNT,SUM_POSITIVE,PER_SECOND,INCREASE}",
// 				ValidateFunc: validation.StringInSlice([]string{"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTINCT_COUNT", "SUM_POSITIVE", "PER_SECOND", "INCREASE"}, false),
// 			},
// 		},
// 	},
// }

// // SyntheticAlertConfigSchemaTimeThreshold schema field definition of instana_synthetic_alert_config field time_threshold
// var SyntheticAlertConfigSchemaTimeThreshold = &schema.Schema{
// 	Type:        schema.TypeList,
// 	Required:    true,
// 	MaxItems:    1,
// 	Description: "Configures the time threshold for the synthetic alert",
// 	Elem: &schema.Resource{
// 		Schema: map[string]*schema.Schema{
// 			SyntheticAlertTimeThresholdFieldType: {
// 				Type:         schema.TypeString,
// 				Required:     true,
// 				Description:  "The type of the time threshold (only violationsInSequence is supported)",
// 				ValidateFunc: validation.StringInSlice([]string{"violationsInSequence"}, false),
// 			},
// 			SyntheticAlertTimeThresholdFieldViolationsCount: {
// 				Type:         schema.TypeInt,
// 				Optional:     true,
// 				Description:  "The number of violations required to trigger the alert (value between 1 and 12)",
// 				ValidateFunc: validation.IntBetween(1, 12),
// 			},
// 		},
// 	},
// }

// // NewSyntheticAlertConfigResourceHandle creates the resource handle for Synthetic Alert Configuration
// func NewSyntheticAlertConfigResourceHandle() ResourceHandle[*restapi.SyntheticAlertConfig] {
// 	return &syntheticAlertConfigResource{
// 		metaData: ResourceMetaData{
// 			ResourceName: ResourceInstanaSyntheticAlertConfig,
// 			Schema: map[string]*schema.Schema{
// 				SyntheticAlertConfigFieldName:             SyntheticAlertConfigSchemaName,
// 				SyntheticAlertConfigFieldDescription:      SyntheticAlertConfigSchemaDescription,
// 				SyntheticAlertConfigFieldSyntheticTestIds: SyntheticAlertConfigSchemaSyntheticTestIds,
// 				SyntheticAlertConfigFieldSeverity:         SyntheticAlertConfigSchemaSeverity,
// 				SyntheticAlertConfigFieldTagFilter:        OptionalTagFilterExpressionSchema,
// 				SyntheticAlertConfigFieldRule:             SyntheticAlertConfigSchemaRule,
// 				SyntheticAlertConfigFieldAlertChannelIds:  SyntheticAlertConfigSchemaAlertChannelIds,
// 				SyntheticAlertConfigFieldTimeThreshold:    SyntheticAlertConfigSchemaTimeThreshold,
// 				SyntheticAlertConfigFieldGracePeriod:      SyntheticAlertConfigSchemaGracePeriod,
// 				DefaultCustomPayloadFieldsName:            buildStaticStringCustomPayloadFields(),
// 			},
// 			SchemaVersion: 1,
// 		},
// 	}
// }

// type syntheticAlertConfigResource struct {
// 	metaData ResourceMetaData
// }

// func (r *syntheticAlertConfigResource) MetaData() *ResourceMetaData {
// 	return &r.metaData
// }

// func (r *syntheticAlertConfigResource) stateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
// 	// Handle field name changes
// 	if _, ok := state[SyntheticAlertConfigFieldFullName]; ok {
// 		state[SyntheticAlertConfigFieldName] = state[SyntheticAlertConfigFieldFullName]
// 		delete(state, SyntheticAlertConfigFieldFullName)
// 	}
// 	return state, nil
// }

// func (r *syntheticAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
// 	return []schema.StateUpgrader{
// 		{
// 			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
// 			Upgrade: r.stateUpgradeV0,
// 			Version: 0,
// 		},
// 	}
// }

// func (r *syntheticAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SyntheticAlertConfig] {
// 	return api.SyntheticAlertConfigs()
// }

// func (r *syntheticAlertConfigResource) SetComputedFields(_ *schema.ResourceData) error {
// 	return nil
// }

// func (r *syntheticAlertConfigResource) UpdateState(d *schema.ResourceData, config *restapi.SyntheticAlertConfig) error {
// 	d.SetId(config.ID)

// 	var normalizedTagFilterString *string
// 	var err error
// 	if config.TagFilterExpression != nil {
// 		normalizedTagFilterString, err = tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return tfutils.UpdateState(d, map[string]interface{}{
// 		SyntheticAlertConfigFieldName:             config.Name,
// 		SyntheticAlertConfigFieldDescription:      config.Description,
// 		SyntheticAlertConfigFieldSyntheticTestIds: config.SyntheticTestIds,
// 		SyntheticAlertConfigFieldSeverity:         config.Severity,
// 		SyntheticAlertConfigFieldTagFilter:        normalizedTagFilterString,
// 		SyntheticAlertConfigFieldRule: []interface{}{
// 			map[string]interface{}{
// 				SyntheticAlertRuleFieldAlertType:   config.Rule.AlertType,
// 				SyntheticAlertRuleFieldMetricName:  config.Rule.MetricName,
// 				SyntheticAlertRuleFieldAggregation: config.Rule.Aggregation,
// 			},
// 		},
// 		SyntheticAlertConfigFieldAlertChannelIds: config.AlertChannelIds,
// 		SyntheticAlertConfigFieldTimeThreshold: []interface{}{
// 			map[string]interface{}{
// 				SyntheticAlertTimeThresholdFieldType:            config.TimeThreshold.Type,
// 				SyntheticAlertTimeThresholdFieldViolationsCount: config.TimeThreshold.ViolationsCount,
// 			},
// 		},
// 		SyntheticAlertConfigFieldGracePeriod: config.GracePeriod,
// 		DefaultCustomPayloadFieldsName:       mapCustomPayloadFieldsToSchema(config),
// 	})
// }

// func (r *syntheticAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SyntheticAlertConfig, error) {
// 	customPayloadFields, err := mapDefaultCustomPayloadFieldsFromSchema(d)
// 	if err != nil {
// 		return &restapi.SyntheticAlertConfig{}, err
// 	}

// 	rule := d.Get(SyntheticAlertConfigFieldRule).([]interface{})[0].(map[string]interface{})
// 	timeThreshold := d.Get(SyntheticAlertConfigFieldTimeThreshold).([]interface{})[0].(map[string]interface{})

// 	var tagFilter *restapi.TagFilter
// 	tagFilterStr, ok := d.GetOk(SyntheticAlertConfigFieldTagFilter)
// 	if ok {
// 		tagFilter, err = mapTagFilterExpressionFromSchema(tagFilterStr.(string))
// 		if err != nil {
// 			return &restapi.SyntheticAlertConfig{}, err
// 		}
// 	}

// 	var gracePeriod int64
// 	if val, ok := d.GetOk(SyntheticAlertConfigFieldGracePeriod); ok {
// 		gracePeriod = int64(val.(int))
// 	}

// 	return &restapi.SyntheticAlertConfig{
// 		ID:                  d.Id(),
// 		Name:                d.Get(SyntheticAlertConfigFieldName).(string),
// 		Description:         d.Get(SyntheticAlertConfigFieldDescription).(string),
// 		SyntheticTestIds:    ReadStringSetParameterFromResource(d, SyntheticAlertConfigFieldSyntheticTestIds),
// 		Severity:            d.Get(SyntheticAlertConfigFieldSeverity).(int),
// 		TagFilterExpression: tagFilter,
// 		Rule: restapi.SyntheticAlertRule{
// 			AlertType:   rule[SyntheticAlertRuleFieldAlertType].(string),
// 			MetricName:  rule[SyntheticAlertRuleFieldMetricName].(string),
// 			Aggregation: rule[SyntheticAlertRuleFieldAggregation].(string),
// 		},
// 		AlertChannelIds: ReadStringSetParameterFromResource(d, SyntheticAlertConfigFieldAlertChannelIds),
// 		TimeThreshold: restapi.SyntheticAlertTimeThreshold{
// 			Type:            timeThreshold[SyntheticAlertTimeThresholdFieldType].(string),
// 			ViolationsCount: timeThreshold[SyntheticAlertTimeThresholdFieldViolationsCount].(int),
// 		},
// 		GracePeriod:           gracePeriod,
// 		CustomerPayloadFields: customPayloadFields,
// 	}, nil
// }

// func mapTagFilterExpressionFromSchema(input string) (*restapi.TagFilter, error) {
// 	parser := tagfilter.NewParser()
// 	expr, err := parser.Parse(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mapper := tagfilter.NewMapper()
// 	return mapper.ToAPIModel(expr), nil
// }

// func (r *syntheticAlertConfigResource) schemaV0() *schema.Resource {
// 	return &schema.Resource{
// 		Schema: map[string]*schema.Schema{
// 			SyntheticAlertConfigFieldName:             SyntheticAlertConfigSchemaName,
// 			SyntheticAlertConfigFieldDescription:      SyntheticAlertConfigSchemaDescription,
// 			SyntheticAlertConfigFieldSyntheticTestIds: SyntheticAlertConfigSchemaSyntheticTestIds,
// 			SyntheticAlertConfigFieldSeverity:         SyntheticAlertConfigSchemaSeverity,
// 			SyntheticAlertConfigFieldTagFilter:        OptionalTagFilterExpressionSchema,
// 			SyntheticAlertConfigFieldRule:             SyntheticAlertConfigSchemaRule,
// 			SyntheticAlertConfigFieldAlertChannelIds:  SyntheticAlertConfigSchemaAlertChannelIds,
// 			SyntheticAlertConfigFieldTimeThreshold:    SyntheticAlertConfigSchemaTimeThreshold,
// 			SyntheticAlertConfigFieldGracePeriod:      SyntheticAlertConfigSchemaGracePeriod,
// 			DefaultCustomPayloadFieldsName:            buildStaticStringCustomPayloadFields(),
// 		},
// 	}
// }
