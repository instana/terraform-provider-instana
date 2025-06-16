package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

func TestInfraAlertConfig(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaInfraAlertConfig + ".example"
	inst := &infraAlertConfigTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewInfraAlertConfigResourceHandle(),
	}
	inst.run(t)
}

type infraAlertConfigTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle[*restapi.InfraAlertConfig]
}

var infraAlertConfigTerraformTemplate = `
resource "instana_infra_alert_config" "example" {
	name              = "name %d"
    description       = "test-alert-description"
    alert_channels {
		warning = [ "alert-channel-id-1" ]
		critical = [ "alert-channel-id-2" ]
	}
    group_by          = [ "metricId" ]
    granularity       = 600000
	tag_filter        = "host.fqdn@na STARTS_WITH 'fooBar'"

	rules {
		generic_rule {
			metric_name 			 = "cpu\\.(nice|user|sys|wait)"
			entity_type 			 = "host"
			aggregation 			 = "MEAN"
			cross_series_aggregation = "SUM"
			regex 					 = true
			threshold_operator       = ">="

			threshold {
				critical {
					static {
						value = 5.0
					}	
				}
			}
		}
    }

	time_threshold {
		violations_in_sequence {
			time_window = 600000
		}
    }

	custom_payload_field {
		key    = "test1"
		value  = "foo"
	}

	custom_payload_field {
		key = "test2"
		dynamic_value {
			key      = "dynamic-value-key"
			tag_name = "dynamic-value-tag-name"
		}
	}
}
`

var infraAlertConfigServerResponseTemplate = `
	{
		"id": "%s",
		"name": "name %d",
		"description": "test-alert-description",
		"alertChannels": {
			"WARNING": ["alert-channel-id-1"],
			"CRITICAL": ["alert-channel-id-2"]
		},
		"tagFilterExpression": {
			"type": "TAG_FILTER",
			"name": "host.fqdn",
			"stringValue": "fooBar",
			"numberValue": null,
			"booleanValue": null,
			"key": null,
			"value": "fooBar",
			"operator": "STARTS_WITH",
			"entity": "NOT_APPLICABLE"
		},
	    "rules":[
		   {
			  "thresholdOperator":">=",
			  "rule":{
			     "alertType":"genericRule",
			     "metricName":"cpu\\.(nice|user|sys|wait)",
			     "entityType":"host",
			     "aggregation":"MEAN",
			     "crossSeriesAggregation":"SUM",
			     "regex":true
			  },
			  "thresholds":{
			     "CRITICAL":{
				    "type":"staticThreshold",
				    "value":5.0
			     }
			  }
		   }
	    ],
		"groupBy": [ "metricId" ],
		"granularity": 600000,
		"timeThreshold": {
		  "type": "violationsInSequence",
		  "timeWindow": 600000
		},
		"customPayloadFields": [
			{
				"type": "staticString",
				"key": "test1",
				"value": "foo"
			},
			{
				"type": "dynamic",
				"key": "test2",
				"value": {
					"key": "dynamic-value-key",
					"tagName": "dynamic-value-tag-name"
				}
			}
		],
		"created": 1647679325301,
		"readOnly": false,
		"enabled": true
  }
`

func (test *infraAlertConfigTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("CRUD integration test of %s", ResourceInstanaInfraAlertConfig), test.createIntegrationTest())
	t.Run(fmt.Sprintf("%s should have schema version one", ResourceInstanaInfraAlertConfig), test.createTestResourceShouldHaveSchemaVersionOne())
	t.Run(fmt.Sprintf("%s should have one state upgrader", ResourceInstanaInfraAlertConfig), test.createTestResourceShouldHaveOneStateUpgrader())
	t.Run(fmt.Sprintf("%s should migrate fullname to name when executing first state migration and fullname is available", ResourceInstanaInfraAlertConfig), test.createTestInfraAlertConfigShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable())
	t.Run(fmt.Sprintf("%s should do nothing when executing first state migration and fullname is not available", ResourceInstanaInfraAlertConfig), test.createTestInfraAlertConfigShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsAvailable())
	t.Run(fmt.Sprintf("%s should have correct resouce name", ResourceInstanaInfraAlertConfig), test.createTestResourceShouldHaveCorrectResourceName())
	test.createTestCasesForUpdatesOfTerraformResourceStateFromModel(t)
	t.Run(fmt.Sprintf("%s should fail to update state from model when tag filter expression is invalid", ResourceInstanaInfraAlertConfig), test.createTestCasesShouldFailToUpdateTerraformResourceStateFromModeWhenTagFilterExpressionIsNotValid())
	test.createTestCasesForMappingOfTerraformResourceStateToModel(t)
	t.Run(fmt.Sprintf("%s should fail to map state to model when tag filter expression is invalid", ResourceInstanaInfraAlertConfig), test.createTestCaseShouldFailToMapTerraformResourceStateToModelWhenTagFilterIsNotValid())
	t.Run(fmt.Sprintf("%s should return errr when converting state to data model and custom field is not valid", ResourceInstanaInfraAlertConfig), test.shouldReturnErrorWhenConvertingStateToDataModelAndCustomFieldIsNotValid)
}

func (test *infraAlertConfigTest) shouldReturnErrorWhenConvertingStateToDataModelAndCustomFieldIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
	sut := test.resourceHandle
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
	setValueOnResourceData(t, resourceData, InfraAlertConfigFieldName, "infra-alert-config-name")
	setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTagFilter, "host.fqdn@na EQUALS 'fooBar'")
	setValueOnResourceData(t, resourceData, DefaultCustomPayloadFieldsName, []interface{}{
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "dynamic-key",
			CustomPayloadFieldsFieldStaticStringValue: "invalid",
			CustomPayloadFieldsFieldDynamicValue: []interface{}{
				map[string]interface{}{
					CustomPayloadFieldsFieldDynamicKey:     "dynamic-value-key",
					CustomPayloadFieldsFieldDynamicTagName: "dynamic-value-tag-name",
				},
			},
		},
	})

	_, err := sut.MapStateToDataObject(resourceData)

	require.Error(t, err)
	require.ErrorContains(t, err, "either a static string value or a dynamic value must")
}

func (test *infraAlertConfigTest) createTestCaseShouldFailToMapTerraformResourceStateToModelWhenTagFilterIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldName, "infra-alert-config-name")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTagFilter, "invalid invalid invalid")

		_, err := sut.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected token")
	}
}

func (test *infraAlertConfigTest) createTestCasesForMappingOfTerraformResourceStateToModel(t *testing.T) {
	metricName := "cpu\\\\.(nice|user|sys|wait)"
	host := "host"
	aggregation := restapi.MeanAggregation

	thresholdValue := 12.3
	rules := []testPair[[]map[string]interface{}, restapi.RuleWithThreshold[restapi.InfraAlertRule]]{
		{
			name: "genericRule",
			expected: restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ThresholdOperator: ">",
				Rule: restapi.InfraAlertRule{
					AlertType:              "genericRule",
					MetricName:             metricName,
					EntityType:             host,
					Aggregation:            aggregation,
					CrossSeriesAggregation: aggregation,
					Regex:                  true,
				},
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &thresholdValue,
					},
				},
			},
			input: []map[string]interface{}{
				{
					InfraAlertConfigFieldGenericRule: []interface{}{
						map[string]interface{}{
							InfraAlertConfigFieldMetricName:             metricName,
							InfraAlertConfigFieldEntityType:             host,
							InfraAlertConfigFieldAggregation:            aggregation,
							InfraAlertConfigFieldCrossSeriesAggregation: aggregation,
							InfraAlertConfigFieldRegex:                  true,

							InfraAlertConfigFieldThresholdOperator: ">",

							ResourceFieldThresholdRule: []interface{}{
								map[string]interface{}{
									ResourceFieldThresholdRuleWarningSeverity: []interface{}{
										map[string]interface{}{
											ResourceFieldThresholdRuleHistoricBaseline: []interface{}{},
											ResourceFieldThresholdRuleStatic: []interface{}{
												map[string]interface{}{
													ResourceFieldThresholdRuleStaticValue: thresholdValue,
												},
											},
										},
									},
									ResourceFieldThresholdRuleCriticalSeverity: []interface{}{},
								},
							},
						},
					},
				},
			},
		},
	}

	timeThresholdWindow := int64(300000)
	timeThresholds := []testPair[[]map[string]interface{}, restapi.InfraTimeThreshold]{
		{
			name: "ViolationsInSequence",
			expected: restapi.InfraTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: timeThresholdWindow,
			},
			input: []map[string]interface{}{
				{
					InfraAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{
						map[string]interface{}{
							InfraAlertConfigFieldTimeThresholdTimeWindow: int(timeThresholdWindow),
						},
					},
				},
			},
		},
	}

	for _, rule := range rules {
		for _, timeThreshold := range timeThresholds {
			t.Run(fmt.Sprintf("Should update terraform state of %s from REST response with %s and %s", ResourceInstanaInfraAlertConfig, rule.name, timeThreshold.name),
				test.createTestShouldMapTerraformResourceStateToModelCase(rule, timeThreshold))
			t.Run(fmt.Sprintf("Should update terraform state of %s from REST response with %s and %s", ResourceInstanaInfraAlertConfig, rule.name, timeThreshold.name),
				test.createTestWithSingleSeverityAlertChannelsShouldMapTerraformResourceStateToModelCase(rule, timeThreshold))
			t.Run(fmt.Sprintf("Should update terraform state of %s from REST response with %s and %s", ResourceInstanaInfraAlertConfig, rule.name, timeThreshold.name),
				test.createTestWithNoAlertChannelsShouldMapTerraformResourceStateToModelCase(rule, timeThreshold))
		}
	}
}

func (test *infraAlertConfigTest) createTestShouldMapTerraformResourceStateToModelCase(
	ruleTestPair testPair[[]map[string]interface{}, restapi.RuleWithThreshold[restapi.InfraAlertRule]],
	timeThresholdTestPair testPair[[]map[string]interface{}, restapi.InfraTimeThreshold]) func(t *testing.T) {

	return func(t *testing.T) {
		infraAlertConfigID := "infra-alert-config-id"
		name := "infra-alert-config-name"
		dynamicValueKey := "dynamic-value-key"
		dynamicValueTagName := "dynamic-value-tag-name"
		expectedInfraConfig := restapi.InfraAlertConfig{
			ID:                  infraAlertConfigID,
			Name:                name,
			Description:         "infra-alert-config-description",
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "host.fqdn", restapi.EqualsOperator, "fooBar"),
			GroupBy:             []string{"metricId"},
			AlertChannels: map[restapi.AlertSeverity][]string{
				restapi.WarningSeverity:  {"channel-1"},
				restapi.CriticalSeverity: {"channel-2"},
			},
			Granularity:   restapi.Granularity300000,
			TimeThreshold: timeThresholdTestPair.expected,
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type:  restapi.StaticStringCustomPayloadType,
					Key:   "static-key",
					Value: restapi.StaticStringCustomPayloadFieldValue("static-value"),
				},
				{
					Type:  restapi.DynamicCustomPayloadType,
					Key:   "dynamic-key",
					Value: restapi.DynamicCustomPayloadFieldValue{Key: &dynamicValueKey, TagName: dynamicValueTagName},
				},
			},
			Rules: []restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ruleTestPair.expected,
			},
		}

		testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldName, "infra-alert-config-name")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldDescription, "infra-alert-config-description")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldAlertChannels, []interface{}{
			map[string]interface{}{
				ResourceFieldThresholdRuleWarningSeverity:  []interface{}{"channel-1"},
				ResourceFieldThresholdRuleCriticalSeverity: []interface{}{"channel-2"},
			},
		})
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldGroupBy, []interface{}{"metricId"})
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldGranularity, restapi.Granularity300000)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTagFilter, "host.fqdn@na EQUALS 'fooBar'")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldRules, ruleTestPair.input)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTimeThreshold, timeThresholdTestPair.input)
		setValueOnResourceData(t, resourceData, DefaultCustomPayloadFieldsName, []interface{}{
			map[string]interface{}{
				CustomPayloadFieldsFieldKey:               "static-key",
				CustomPayloadFieldsFieldStaticStringValue: "static-value",
				CustomPayloadFieldsFieldDynamicValue:      []interface{}{},
			},
			map[string]interface{}{
				CustomPayloadFieldsFieldKey:          "dynamic-key",
				CustomPayloadFieldsFieldDynamicValue: []interface{}{map[string]interface{}{CustomPayloadFieldsFieldDynamicKey: dynamicValueKey, CustomPayloadFieldsFieldDynamicTagName: dynamicValueTagName}},
			},
		})

		resourceData.SetId(infraAlertConfigID)

		result, err := sut.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.Equal(t, &expectedInfraConfig, result)
	}
}

func (test *infraAlertConfigTest) createTestCasesShouldFailToUpdateTerraformResourceStateFromModeWhenTagFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		value := "fooBar"
		operator := restapi.EqualsOperator
		tagFilterName := "host.fqdn"
		tagFilterEntity := restapi.TagFilterEntityNotApplicable
		infraAlertConfig := restapi.InfraAlertConfig{
			Name: "test",
			TagFilterExpression: &restapi.TagFilter{
				Entity:      &tagFilterEntity,
				Name:        &tagFilterName,
				Operator:    &operator,
				StringValue: &value,
				Value:       value,
				Type:        restapi.TagFilterExpressionElementType("invalid"),
			},
		}

		testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &infraAlertConfig)

		require.Error(t, err)
		require.Equal(t, "unsupported tag filter expression of type invalid", err.Error())
	}
}

func (test *infraAlertConfigTest) createTestCasesForUpdatesOfTerraformResourceStateFromModel(t *testing.T) {
	metricName := "cpu\\\\.(nice|user|sys|wait)"
	host := "host"
	aggregation := restapi.MeanAggregation

	thresholdValue := 12.3
	rules := []testPair[restapi.RuleWithThreshold[restapi.InfraAlertRule], []interface{}]{
		{
			name: "genericRule",
			input: restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ThresholdOperator: ">",
				Rule: restapi.InfraAlertRule{
					AlertType:              "genericRule",
					MetricName:             metricName,
					EntityType:             host,
					Aggregation:            aggregation,
					CrossSeriesAggregation: aggregation,
					Regex:                  true,
				},
				Thresholds: map[restapi.AlertSeverity]restapi.ThresholdRule{
					restapi.WarningSeverity: {
						Type:  "staticThreshold",
						Value: &thresholdValue,
					},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					InfraAlertConfigFieldGenericRule: []interface{}{
						map[string]interface{}{
							InfraAlertConfigFieldMetricName:             metricName,
							InfraAlertConfigFieldEntityType:             host,
							InfraAlertConfigFieldAggregation:            "MEAN",
							InfraAlertConfigFieldCrossSeriesAggregation: "MEAN",
							InfraAlertConfigFieldRegex:                  true,

							InfraAlertConfigFieldThresholdOperator: ">",

							ResourceFieldThresholdRule: []interface{}{
								map[string]interface{}{
									ResourceFieldThresholdRuleWarningSeverity: []interface{}{
										map[string]interface{}{
											ResourceFieldThresholdRuleHistoricBaseline: []interface{}{},
											ResourceFieldThresholdRuleStatic: []interface{}{
												map[string]interface{}{
													ResourceFieldThresholdRuleStaticValue: thresholdValue,
												},
											},
										},
									},
									ResourceFieldThresholdRuleCriticalSeverity: []interface{}{},
								},
							},
						},
					},
				},
			},
		},
	}

	timeThresholdWindow := int64(300000)
	timeThresholds := []testPair[restapi.InfraTimeThreshold, []interface{}]{
		{
			name: "ViolationsInSequence",
			input: restapi.InfraTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: timeThresholdWindow,
			},
			expected: []interface{}{
				map[string]interface{}{
					InfraAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{
						map[string]interface{}{
							InfraAlertConfigFieldTimeThresholdTimeWindow: int(timeThresholdWindow),
						},
					},
				},
			},
		},
	}

	for _, rule := range rules {
		for _, timeThreshold := range timeThresholds {
			t.Run(fmt.Sprintf("Should update terraform state of %s from REST response with %s and %s", ResourceInstanaInfraAlertConfig, rule.name, timeThreshold.name),
				test.createTestShouldUpdateTerraformResourceStateFromModelCase(rule, timeThreshold))
		}
	}
}

func (test *infraAlertConfigTest) createTestShouldUpdateTerraformResourceStateFromModelCase(
	ruleTestPair testPair[restapi.RuleWithThreshold[restapi.InfraAlertRule], []interface{}],
	timeThresholdTestPair testPair[restapi.InfraTimeThreshold, []interface{}]) func(t *testing.T) {

	return func(t *testing.T) {
		infraAlertConfigID := "infra-alert-config-id"
		name := "infra-alert-config-name"
		dynamicValueKey := "dynamic-value-key"
		dynamicValueTagName := "dynamic-value-tag-name"
		alertChannelsMap := make(map[restapi.AlertSeverity][]string)
		alertChannelsMap[restapi.WarningSeverity] = []string{"channel-1"}
		alertChannelsMap[restapi.CriticalSeverity] = []string{"channel-2"}

		infraConfig := restapi.InfraAlertConfig{
			ID:                  infraAlertConfigID,
			Name:                name,
			Description:         "infra-alert-config-description",
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "host.fqdn", restapi.EqualsOperator, "fooBar"),
			GroupBy:             []string{"metricId"},
			AlertChannels:       alertChannelsMap,
			Granularity:         restapi.Granularity300000,
			TimeThreshold:       timeThresholdTestPair.input,
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type:  restapi.StaticStringCustomPayloadType,
					Key:   "static-key",
					Value: restapi.StaticStringCustomPayloadFieldValue("static-value"),
				},
				{
					Type:  restapi.DynamicCustomPayloadType,
					Key:   "dynamic-key",
					Value: restapi.DynamicCustomPayloadFieldValue{Key: &dynamicValueKey, TagName: dynamicValueTagName},
				},
			},
			Rules: []restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ruleTestPair.input,
			},
		}

		testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &infraConfig)
		require.NoError(t, err)
		require.Equal(t, infraAlertConfigID, resourceData.Id())

		require.Equal(t, "infra-alert-config-name", resourceData.Get(InfraAlertConfigFieldName))
		require.Equal(t, "infra-alert-config-description", resourceData.Get(InfraAlertConfigFieldDescription))
		require.Equal(t, []interface{}{
			map[string]interface{}{
				ResourceFieldThresholdRuleWarningSeverity:  []interface{}{"channel-1"},
				ResourceFieldThresholdRuleCriticalSeverity: []interface{}{"channel-2"},
			},
		}, resourceData.Get(InfraAlertConfigFieldAlertChannels).([]interface{}))
		require.Equal(t, []interface{}{"metricId"}, resourceData.Get(InfraAlertConfigFieldGroupBy).([]interface{}))
		require.Equal(t, ruleTestPair.expected, resourceData.Get(InfraAlertConfigFieldRules).([]interface{}))

		require.Equal(t, []interface{}{
			map[string]interface{}{
				CustomPayloadFieldsFieldKey:               "static-key",
				CustomPayloadFieldsFieldDynamicValue:      []interface{}{},
				CustomPayloadFieldsFieldStaticStringValue: "static-value",
			},
			map[string]interface{}{
				CustomPayloadFieldsFieldKey:               "dynamic-key",
				CustomPayloadFieldsFieldDynamicValue:      []interface{}{map[string]interface{}{CustomPayloadFieldsFieldDynamicKey: dynamicValueKey, CustomPayloadFieldsFieldDynamicTagName: dynamicValueTagName}},
				CustomPayloadFieldsFieldStaticStringValue: "",
			},
		}, resourceData.Get(DefaultCustomPayloadFieldsName).(*schema.Set).List())
		require.Equal(t, "host.fqdn@na EQUALS 'fooBar'", resourceData.Get(InfraAlertConfigFieldTagFilter))
		require.Equal(t, timeThresholdTestPair.expected, resourceData.Get(InfraAlertConfigFieldTimeThreshold))
	}
}

func (test *infraAlertConfigTest) createIntegrationTest() func(t *testing.T) {
	return func(t *testing.T) {
		id := RandomID()
		resourceRestAPIPath := restapi.InfraAlertConfigResourcePath
		resourceInstanceRestAPIPath := resourceRestAPIPath + "/{internal-id}"

		httpServer := testutils.NewTestHTTPServer()
		httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			config := &restapi.InfraAlertConfig{}
			err := json.NewDecoder(r.Body).Decode(config)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				err = r.Write(bytes.NewBufferString("Failed to get request"))
				if err != nil {
					fmt.Printf("failed to write response; %s\n", err)
				}
			} else {
				config.ID = id
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(config)
				if err != nil {
					fmt.Printf("failed to encode json; %s\n", err)
				}
			}
		})
		httpServer.AddRoute(http.MethodPost, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			testutils.EchoHandlerFunc(w, r)
		})
		httpServer.AddRoute(http.MethodDelete, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
		httpServer.AddRoute(http.MethodGet, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			modCount := httpServer.GetCallCount(http.MethodPost, resourceRestAPIPath+"/"+id)
			jsonData := fmt.Sprintf(infraAlertConfigServerResponseTemplate, id, modCount)
			w.Header().Set(contentType, r.Header.Get(contentType))
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(jsonData))
			if err != nil {
				fmt.Printf("failed to write response; %s\n", err)
			}
		})
		httpServer.Start()
		defer httpServer.Close()

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: testProviderFactory,
			Steps: []resource.TestStep{
				test.createIntegrationTestStep(httpServer.GetPort(), 0, id),
				testStepImportWithCustomID(test.terraformResourceInstanceName, id),
				test.createIntegrationTestStep(httpServer.GetPort(), 1, id),
				testStepImportWithCustomID(test.terraformResourceInstanceName, id),
			},
		})
	}
}

func (test *infraAlertConfigTest) createIntegrationTestStep(httpPort int, iteration int, id string) resource.TestStep {
	ruleMetricName := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldRules, 0, InfraAlertConfigFieldGenericRule, 0, InfraAlertConfigFieldMetricName)
	ruleEntityType := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldRules, 0, InfraAlertConfigFieldGenericRule, 0, InfraAlertConfigFieldEntityType)
	ruleAggregation := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldRules, 0, InfraAlertConfigFieldGenericRule, 0, InfraAlertConfigFieldAggregation)
	ruleCrossSeriesAggregation := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldRules, 0, InfraAlertConfigFieldGenericRule, 0, InfraAlertConfigFieldCrossSeriesAggregation)
	ruleRegex := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldRules, 0, InfraAlertConfigFieldGenericRule, 0, InfraAlertConfigFieldRegex)
	thresholdOperator := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldRules, 0, InfraAlertConfigFieldGenericRule, 0, InfraAlertConfigFieldThresholdOperator)

	thresholdStaticValue := fmt.Sprintf("%s.0.%s.0.%s.0.%s.0.%s.0.%s", InfraAlertConfigFieldRules, InfraAlertConfigFieldGenericRule, ResourceFieldThresholdRule, ResourceFieldThresholdRuleCriticalSeverity, ResourceFieldThresholdRuleStatic, ResourceFieldThresholdRuleStaticValue)
	timeThresholdViolationsInSequence := fmt.Sprintf("%s.%d.%s.%d.%s", InfraAlertConfigFieldTimeThreshold, 0, InfraAlertConfigFieldTimeThresholdViolationsInSequence, 0, InfraAlertConfigFieldTimeThresholdTimeWindow)

	customPayloadFieldStaticKey := fmt.Sprintf("%s.1.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldKey)
	customPayloadFieldStaticValue := fmt.Sprintf("%s.1.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldStaticStringValue)
	customPayloadFieldDynamicKey := fmt.Sprintf("%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldKey)
	customPayloadFieldDynamicValueKey := fmt.Sprintf("%s.0.%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldDynamicValue, CustomPayloadFieldsFieldDynamicKey)
	customPayloadFieldDynamicValueTagName := fmt.Sprintf("%s.0.%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldDynamicValue, CustomPayloadFieldsFieldDynamicTagName)

	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(infraAlertConfigTerraformTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, "id", id),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldDescription, "test-alert-description"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldAlertChannels+".0."+ResourceFieldThresholdRuleWarningSeverity+".0", "alert-channel-id-1"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldAlertChannels+".0."+ResourceFieldThresholdRuleCriticalSeverity+".0", "alert-channel-id-2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldGroupBy+".0", "metricId"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldGranularity, "600000"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, InfraAlertConfigFieldTagFilter, "host.fqdn@na STARTS_WITH 'fooBar'"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleMetricName, "cpu\\.(nice|user|sys|wait)"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleEntityType, "host"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleAggregation, "MEAN"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleCrossSeriesAggregation, "SUM"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleRegex, "true"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, thresholdOperator, ">="),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, thresholdStaticValue, "5"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, timeThresholdViolationsInSequence, "600000"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldStaticKey, "test1"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldStaticValue, "foo"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldDynamicKey, "test2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldDynamicValueKey, "dynamic-value-key"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldDynamicValueTagName, "dynamic-value-tag-name"),
		),
	}
}

func (test *infraAlertConfigTest) createTestResourceShouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, test.resourceHandle.MetaData().SchemaVersion)
	}
}

func (test *infraAlertConfigTest) createTestResourceShouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Len(t, test.resourceHandle.StateUpgraders(), 1)
	}
}

func (test *infraAlertConfigTest) createTestInfraAlertConfigShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"full_name": "test",
		}
		result, err := NewInfraAlertConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.NotContains(t, result, InfraAlertConfigFieldFullName)
		require.Contains(t, result, InfraAlertConfigFieldName)
		require.Equal(t, "test", result[InfraAlertConfigFieldName])
	}
}

func (test *infraAlertConfigTest) createTestInfraAlertConfigShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"name": "test",
		}
		result, err := NewInfraAlertConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Equal(t, input, result)
	}
}

func (test *infraAlertConfigTest) createTestResourceShouldHaveCorrectResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, test.resourceHandle.MetaData().ResourceName, "instana_infra_alert_config")
	}
}

func (test *infraAlertConfigTest) createTestWithSingleSeverityAlertChannelsShouldMapTerraformResourceStateToModelCase(
	ruleTestPair testPair[[]map[string]interface{}, restapi.RuleWithThreshold[restapi.InfraAlertRule]],
	timeThresholdTestPair testPair[[]map[string]interface{}, restapi.InfraTimeThreshold]) func(t *testing.T) {

	return func(t *testing.T) {
		infraAlertConfigID := "infra-alert-config-id"
		name := "infra-alert-config-name"
		expectedInfraConfig := restapi.InfraAlertConfig{
			ID:                  infraAlertConfigID,
			Name:                name,
			Description:         "infra-alert-config-description",
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "host.fqdn", restapi.EqualsOperator, "fooBar"),
			GroupBy:             []string{"metricId"},
			AlertChannels: map[restapi.AlertSeverity][]string{
				restapi.WarningSeverity: {"channel-1"},
			},
			Granularity:   restapi.Granularity300000,
			TimeThreshold: timeThresholdTestPair.expected,
			Rules: []restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ruleTestPair.expected,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldName, "infra-alert-config-name")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldDescription, "infra-alert-config-description")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldAlertChannels, []interface{}{
			map[string]interface{}{
				ResourceFieldThresholdRuleWarningSeverity: []interface{}{"channel-1"},
			},
		})
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldGroupBy, []interface{}{"metricId"})
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldGranularity, restapi.Granularity300000)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTagFilter, "host.fqdn@na EQUALS 'fooBar'")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldRules, ruleTestPair.input)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTimeThreshold, timeThresholdTestPair.input)

		resourceData.SetId(infraAlertConfigID)

		result, err := sut.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.Equal(t, &expectedInfraConfig, result)
	}
}

func (test *infraAlertConfigTest) createTestWithNoAlertChannelsShouldMapTerraformResourceStateToModelCase(
	ruleTestPair testPair[[]map[string]interface{}, restapi.RuleWithThreshold[restapi.InfraAlertRule]],
	timeThresholdTestPair testPair[[]map[string]interface{}, restapi.InfraTimeThreshold]) func(t *testing.T) {

	return func(t *testing.T) {
		infraAlertConfigID := "infra-alert-config-id"
		name := "infra-alert-config-name"
		expectedInfraConfig := restapi.InfraAlertConfig{
			ID:                  infraAlertConfigID,
			Name:                name,
			Description:         "infra-alert-config-description",
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntityNotApplicable, "host.fqdn", restapi.EqualsOperator, "fooBar"),
			GroupBy:             []string{"metricId"},
			AlertChannels:       map[restapi.AlertSeverity][]string{},
			Granularity:         restapi.Granularity300000,
			TimeThreshold:       timeThresholdTestPair.expected,
			Rules: []restapi.RuleWithThreshold[restapi.InfraAlertRule]{
				ruleTestPair.expected,
			},
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{},
		}

		testHelper := NewTestHelper[*restapi.InfraAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldName, "infra-alert-config-name")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldDescription, "infra-alert-config-description")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldGroupBy, []interface{}{"metricId"})
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldGranularity, restapi.Granularity300000)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTagFilter, "host.fqdn@na EQUALS 'fooBar'")
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldRules, ruleTestPair.input)
		setValueOnResourceData(t, resourceData, InfraAlertConfigFieldTimeThreshold, timeThresholdTestPair.input)

		resourceData.SetId(infraAlertConfigID)

		result, err := sut.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.Equal(t, &expectedInfraConfig, result)
	}
}
