package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestSyntheticAlertConfig(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaSyntheticAlertConfig + ".example"
	inst := &syntheticAlertConfigTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewSyntheticAlertConfigResourceHandle(),
	}
	inst.run(t)
}

type syntheticAlertConfigTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle[*restapi.SyntheticAlertConfig]
}

var syntheticAlertConfigTerraformTemplate = `
resource "instana_synthetic_alert_config" "example" {
	name = "name %d"
	description = "description %d"
	synthetic_test_ids = [ "test-1", "test-2" ]
	severity = 5
	tag_filter = "synthetic.testId@na EQUALS 'test-1'"
	rule {
		alert_type = "failure"
		metric_name = "status"
		aggregation = "SUM"
	}
	alert_channel_ids = [ "channel-1", "channel-2" ]
	time_threshold {
		type = "violationsInSequence"
		violations_count = 2
	}

	custom_payload_field {
		key    = "static-key"
		value  = "static-value"
	}
}
`

var syntheticAlertConfigServerResponseTemplate = `
{
	"id" : "%s",
	"name" : "name %d",
	"description" : "description %d",
	"syntheticTestIds" : [ "test-2", "test-1" ],
	"severity" : 5,
	"tagFilterExpression" : {
		"type": "TAG_FILTER",
		"name": "synthetic.testId",
		"stringValue": "test-1",
		"operator": "EQUALS",
		"entity": "NOT_APPLICABLE"
	},
	"rule" : {
		"alertType" : "failure",
		"metricName" : "status",
		"aggregation" : "SUM"
	},
	"alertChannelIds" : [ "channel-2", "channel-1" ],
	"timeThreshold" : {
		"type" : "violationsInSequence",
		"violationsCount" : 2
	},
	   "customPayloadFields": [
		{
			"type": "staticString",
			"key": "static-key",
			"value": "static-value"
      	}
	]
}
`

func (test *syntheticAlertConfigTest) run(t *testing.T) {
	t.Run("CRUD integration test", test.createIntegrationTest())
	t.Run("schema should be valid", test.resourceSchemaShouldBeValid)
	t.Run("should return correct schema name", test.shouldReturnCorrectResourceNameForSyntheticAlertConfig)
	t.Run("should have schema version one", test.shouldHaveSchemaVersionOne)
	t.Run("should have one state upgrader for version zero", test.shouldHaveOneStateUpgraderForVersionZero)
	t.Run("should update resource state", test.shouldUpdateResourceState)
	t.Run("should convert state to data model", test.shouldConvertStateToDataModel)
	t.Run("should migrate full_name to name when executing state upgrader and full_name is available", test.shouldMigrateFullNameToNameWhenExecutingStateUpgraderAndFullNameIsAvailable)
	t.Run("should do nothing when executing state upgrader and full_name is not available", test.shouldDoNothingWhenExecutingStateUpgraderAndFullNameIsNotAvailable)
}

func (test *syntheticAlertConfigTest) createIntegrationTest() func(t *testing.T) {
	return func(t *testing.T) {
		id := RandomID()
		resourceRestAPIPath := restapi.SyntheticAlertConfigsResourcePath
		resourceInstanceRestAPIPath := resourceRestAPIPath + "/{internal-id}"

		httpServer := testutils.NewTestHTTPServer()
		httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			config := &restapi.SyntheticAlertConfig{}
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
			jsonData := fmt.Sprintf(syntheticAlertConfigServerResponseTemplate, id, modCount, modCount)
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

func (test *syntheticAlertConfigTest) createIntegrationTestStep(httpPort int, iteration int, id string) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(syntheticAlertConfigTerraformTemplate, iteration, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, "id", id),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, SyntheticAlertConfigFieldName, fmt.Sprintf("name %d", iteration)),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, SyntheticAlertConfigFieldDescription, fmt.Sprintf("description %d", iteration)),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.#", SyntheticAlertConfigFieldSyntheticTestIds), "2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, SyntheticAlertConfigFieldSeverity, "5"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, SyntheticAlertConfigFieldTagFilter, "synthetic.testId@na EQUALS 'test-1'"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", SyntheticAlertConfigFieldRule, SyntheticAlertRuleFieldAlertType), "failure"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", SyntheticAlertConfigFieldRule, SyntheticAlertRuleFieldMetricName), "status"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", SyntheticAlertConfigFieldRule, SyntheticAlertRuleFieldAggregation), "SUM"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.#", SyntheticAlertConfigFieldAlertChannelIds), "2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", SyntheticAlertConfigFieldTimeThreshold, SyntheticAlertTimeThresholdFieldType), "violationsInSequence"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", SyntheticAlertConfigFieldTimeThreshold, SyntheticAlertTimeThresholdFieldViolationsCount), "2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldKey), "static-key"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, fmt.Sprintf("%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldStaticStringValue), "static-value"),
		),
	}
}

func (test *syntheticAlertConfigTest) resourceSchemaShouldBeValid(t *testing.T) {
	resourceHandle := NewSyntheticAlertConfigResourceHandle()
	schemaMap := resourceHandle.MetaData().Schema

	// Verify required fields
	require.Equal(t, schema.TypeString, schemaMap[SyntheticAlertConfigFieldName].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldName].Required)

	require.Equal(t, schema.TypeString, schemaMap[SyntheticAlertConfigFieldDescription].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldDescription].Required)

	require.Equal(t, schema.TypeSet, schemaMap[SyntheticAlertConfigFieldSyntheticTestIds].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldSyntheticTestIds].Required)

	require.Equal(t, schema.TypeInt, schemaMap[SyntheticAlertConfigFieldSeverity].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldSeverity].Optional)

	require.Equal(t, schema.TypeString, schemaMap[SyntheticAlertConfigFieldTagFilter].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldTagFilter].Optional)

	require.Equal(t, schema.TypeList, schemaMap[SyntheticAlertConfigFieldRule].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldRule].Required)

	require.Equal(t, schema.TypeSet, schemaMap[SyntheticAlertConfigFieldAlertChannelIds].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldAlertChannelIds].Required)

	require.Equal(t, schema.TypeList, schemaMap[SyntheticAlertConfigFieldTimeThreshold].Type)
	require.True(t, schemaMap[SyntheticAlertConfigFieldTimeThreshold].Required)
}

func (test *syntheticAlertConfigTest) shouldReturnCorrectResourceNameForSyntheticAlertConfig(t *testing.T) {
	name := NewSyntheticAlertConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_synthetic_alert_config", name, "Expected resource name to be instana_synthetic_alert_config")
}

func (test *syntheticAlertConfigTest) shouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewSyntheticAlertConfigResourceHandle().MetaData().SchemaVersion)
}

func (test *syntheticAlertConfigTest) shouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewSyntheticAlertConfigResourceHandle()

	require.Equal(t, 1, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
}

func (test *syntheticAlertConfigTest) shouldMigrateFullNameToNameWhenExecutingStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		SyntheticAlertConfigFieldFullName: "test",
	}
	result, err := NewSyntheticAlertConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, SyntheticAlertConfigFieldFullName)
	require.Contains(t, result, SyntheticAlertConfigFieldName)
	require.Equal(t, "test", result[SyntheticAlertConfigFieldName])
}

func (test *syntheticAlertConfigTest) shouldDoNothingWhenExecutingStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		SyntheticAlertConfigFieldName: "test",
	}
	result, err := NewSyntheticAlertConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

const (
	syntheticAlertConfigID          = "synthetic-alert-id"
	syntheticAlertConfigName        = "synthetic-alert-name"
	syntheticAlertConfigDescription = "synthetic-alert-description"
	syntheticAlertConfigTestId1     = "synthetic-test-id1"
	syntheticAlertConfigTestId2     = "synthetic-test-id2"
	syntheticAlertConfigChannelId1  = "synthetic-channel-id1"
	syntheticAlertConfigChannelId2  = "synthetic-channel-id2"
	syntheticAlertConfigTagFilter   = "synthetic.testId@na EQUALS 'test-1'"
)

func (test *syntheticAlertConfigTest) shouldUpdateResourceState(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticAlertConfig](t)
	resourceHandle := NewSyntheticAlertConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	data := restapi.SyntheticAlertConfig{
		ID:               syntheticAlertConfigID,
		Name:             syntheticAlertConfigName,
		Description:      syntheticAlertConfigDescription,
		SyntheticTestIds: []string{syntheticAlertConfigTestId1, syntheticAlertConfigTestId2},
		Severity:         5,
		TagFilterExpression: &restapi.TagFilter{
			Type:        restapi.TagFilterType,
			Name:        stringPtr("synthetic.testId"),
			StringValue: stringPtr("test-1"),
			Operator:    expressionOperatorPtr(restapi.EqualsOperator),
			Entity:      tagFilterEntityPtr(restapi.TagFilterEntityNotApplicable),
		},
		Rule: restapi.SyntheticAlertRule{
			AlertType:   "failure",
			MetricName:  "status",
			Aggregation: "SUM",
		},
		AlertChannelIds: []string{syntheticAlertConfigChannelId1, syntheticAlertConfigChannelId2},
		TimeThreshold: restapi.SyntheticAlertTimeThreshold{
			Type:            "violationsInSequence",
			ViolationsCount: 2,
		},
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "static-key-1",
				Value: restapi.StaticStringCustomPayloadFieldValue("static-value-1"),
			},
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "static-key-2",
				Value: restapi.StaticStringCustomPayloadFieldValue("static-value-2"),
			},
		},
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)
	require.Equal(t, syntheticAlertConfigID, resourceData.Id())
	require.Equal(t, syntheticAlertConfigName, resourceData.Get(SyntheticAlertConfigFieldName))
	require.Equal(t, syntheticAlertConfigDescription, resourceData.Get(SyntheticAlertConfigFieldDescription))
	require.Equal(t, 5, resourceData.Get(SyntheticAlertConfigFieldSeverity))
	require.Equal(t, syntheticAlertConfigTagFilter, resourceData.Get(SyntheticAlertConfigFieldTagFilter))

	// Check synthetic test IDs
	syntheticTestIds := resourceData.Get(SyntheticAlertConfigFieldSyntheticTestIds).(*schema.Set)
	require.Equal(t, 2, syntheticTestIds.Len())
	for _, v := range []string{syntheticAlertConfigTestId1, syntheticAlertConfigTestId2} {
		require.Contains(t, syntheticTestIds.List(), v)
	}

	// Check alert channel IDs
	alertChannelIds := resourceData.Get(SyntheticAlertConfigFieldAlertChannelIds).(*schema.Set)
	require.Equal(t, 2, alertChannelIds.Len())
	for _, v := range []string{syntheticAlertConfigChannelId1, syntheticAlertConfigChannelId2} {
		require.Contains(t, alertChannelIds.List(), v)
	}

	// Check custom payload fields
	fields := resourceData.Get(DefaultCustomPayloadFieldsName)
	require.NotNil(t, fields)
	require.IsType(t, &schema.Set{}, fields)
	fieldList := fields.(*schema.Set).List()
	require.Len(t, fieldList, 2)
}

func (test *syntheticAlertConfigTest) shouldConvertStateToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticAlertConfig](t)
	resourceHandle := NewSyntheticAlertConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(syntheticAlertConfigID)

	testIds := []string{syntheticAlertConfigTestId1, syntheticAlertConfigTestId2}
	channelIds := []string{syntheticAlertConfigChannelId1, syntheticAlertConfigChannelId2}

	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldName, syntheticAlertConfigName)
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldDescription, syntheticAlertConfigDescription)
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldSyntheticTestIds, testIds)
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldSeverity, 5)
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldTagFilter, syntheticAlertConfigTagFilter)
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldRule, []interface{}{
		map[string]interface{}{
			SyntheticAlertRuleFieldAlertType:   "failure",
			SyntheticAlertRuleFieldMetricName:  "status",
			SyntheticAlertRuleFieldAggregation: "SUM",
		},
	})
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldAlertChannelIds, channelIds)
	setValueOnResourceData(t, resourceData, SyntheticAlertConfigFieldTimeThreshold, []interface{}{
		map[string]interface{}{
			SyntheticAlertTimeThresholdFieldType:            "violationsInSequence",
			SyntheticAlertTimeThresholdFieldViolationsCount: 2,
		},
	})
	setValueOnResourceData(t, resourceData, DefaultCustomPayloadFieldsName, []interface{}{
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "static-key-1",
			CustomPayloadFieldsFieldStaticStringValue: "static-value-1",
		},
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "static-key-2",
			CustomPayloadFieldsFieldStaticStringValue: "static-value-2",
		},
	})

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.SyntheticAlertConfig{}, model)
	require.Equal(t, syntheticAlertConfigID, model.GetIDForResourcePath())
	require.Equal(t, syntheticAlertConfigName, model.Name)
	require.Equal(t, syntheticAlertConfigDescription, model.Description)
	require.Equal(t, 5, model.Severity)
	require.Equal(t, "failure", model.Rule.AlertType)
	require.Equal(t, "status", model.Rule.MetricName)
	require.Equal(t, "SUM", model.Rule.Aggregation)
	require.Equal(t, "violationsInSequence", model.TimeThreshold.Type)
	require.Equal(t, 2, model.TimeThreshold.ViolationsCount)

	// Check synthetic test IDs
	require.Equal(t, 2, len(model.SyntheticTestIds))
	for _, v := range testIds {
		require.Contains(t, model.SyntheticTestIds, v)
	}

	// Check alert channel IDs
	require.Equal(t, 2, len(model.AlertChannelIds))
	for _, v := range channelIds {
		require.Contains(t, model.AlertChannelIds, v)
	}

	// Check custom payload fields
	require.Equal(t, 2, len(model.CustomerPayloadFields))
}

func stringPtr(s string) *string {
	return &s
}

func expressionOperatorPtr(op restapi.ExpressionOperator) *restapi.ExpressionOperator {
	return &op
}

func tagFilterEntityPtr(entity restapi.TagFilterEntity) *restapi.TagFilterEntity {
	return &entity
}
