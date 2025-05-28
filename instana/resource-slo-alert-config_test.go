package instana_test

import (
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

const (
	sloAlertConfigDefinition               = "instana_slo_alert_config.example_slo_alert_config"
	sloAlertID                             = "id"
	sloAlertName                           = "slo-alert-name"
	sloAlertDescription                    = "slo-alert-description"
	sloAlertSeverity                       = 10
	SloAlertTriggering                     = true
	SloAlertConfigFieldRule                = "rule"
	sloAlertAlertType                      = "status"
	sloAlertAlertEnabled                   = true
	SloAlertThresholdType                  = "staticThreshold"
	SloAlertThresholdOperator              = ">"
	SloAlertThresholdValue                 = 0.3
	SloAlertTimeThresholdWarmUp            = 60000
	SloAlertTimeThresholdCoolDown          = 30000
	SloAlertBurnRateConfigAlertWindowType  = "SINGLE"
	SloAlertBurnRateConfigDuration         = 1
	SloAlertBurnRateConfigDurationUnitType = "hour"
	SloAlertBurnRateConfigOperator         = ">="
	SloAlertBurnRateConfigValue            = 1.0
	SloAlertBurnRateV2                     = "burn_rate_v2"
)

func TestSloAlertConfig(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaSloAlertConfig + ".example"
	inst := &sloAlertConfigTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewSloAlertConfigResourceHandle(),
	}
	inst.run(t)
}

type sloAlertConfigTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle[*restapi.SloAlertConfig]
}

var sloAlertConfigTerraformTemplate = `
resource "instana_slo_alert_config" "status_alert" {
  name         = "terraform_status_alert"
  description  = "terraform_status_alert testing"
  severity     = 5
  triggering   = true
  slo_ids           = [instana_slo_config.slo4Alert_1.id, instana_slo_config.slo4Alert_2.id]
  alert_channel_ids = ["orhurugksjfgh"]

  alert_type   = "status"
  threshold {
    operator = ">"
    value    = 0.3
  }
  time_threshold {
    warm_up     = 60000
    cool_down   = 60000
  }

  custom_payload_field {
    key    = "test1"
    value  = "foo"
  }

  enabled  = true
}
`

var sloAlertConfigServerResponseTemplate = `
{
  "id": "cv43l6lbd0vp6adbjmug",
  "name": "terraform_status_alert",
  "description": "terraform_status_alert testing",
  "severity": 5,
  "triggering": true,
  "enabled": true,
  "rule": {
    "alertType": "SERVICE_LEVELS_OBJECTIVE",
    "metric": "STATUS"
  },
  "threshold": {
    "type": "staticThreshold",
    "operator": "\u003e",
    "value": 0.3
  },
  "timeThreshold": {
    "timeWindow": 60000,
    "expiry": 60000
  },
  "sloIds": [
    "SLOTFcv43l6dbd0vp6adbjmtg",
    "SLOTFcv43l6dbd0vp6adbjmu0"
  ],
  "alertChannelIds": [
    "orhurugksjfgh"
  ],
  "customPayloadFields": [
    {
      "type": "staticString",
      "key": "test1",
      "value": "foo"
    }
  ]
}
`

func (test *sloAlertConfigTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("%s should have correct resouce name", ResourceInstanaSloAlertConfig), test.createTestResourceShouldHaveResourceName())
	t.Run(fmt.Sprintf("%s should have one state upgrader", ResourceInstanaSloAlertConfig), test.createTestResourceShouldHaveOneStateUpgrader())
	t.Run(fmt.Sprintf("%s should have schema version one", ResourceInstanaSloAlertConfig), test.createTestResourceShouldHaveSchemaVersionOne())
	t.Run("Should require Threshold Values to Match the Assigned Values", test.shouldRequireThresholdMetricsToMatchTheAssignedValues())
	t.Run("Should require Time Threshold Values to Match the Assigned Values", test.shouldRequireTimeThresholdMetricsToMatchTheAssignedValues())
	t.Run("Should require Alert Type Value to Match the Assigned Value", test.shouldFailToMapAlertTypeWhenNoSupportedValueIsProvided())
	t.Run("should map Smart Alert values correctly", test.shouldMapSmartAlertValuesCorrectly())
	t.Run("Should require Burn Rate Config Values to Match the Assigned Values", test.shouldRequireBurnRateConfigToMatchAssignedValues())
	t.Run("Should Fail When Burn Rate Config Is Missing For Burn Rate Alert", test.shouldFailWhenBurnRateConfigIsMissingForBurnRateV2Alert())

}

func (test *sloAlertConfigTest) createTestResourceShouldHaveResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, test.resourceHandle.MetaData().ResourceName, "instana_slo_alert_config")
	}
}

func (test *sloAlertConfigTest) createTestResourceShouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Len(t, test.resourceHandle.StateUpgraders(), 1)
	}
}

func (test *sloAlertConfigTest) createTestResourceShouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, test.resourceHandle.MetaData().SchemaVersion)
	}
}

func (test *sloAlertConfigTest) shouldRequireThresholdMetricsToMatchTheAssignedValues() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
		resourceHandle := NewSloAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloAlertID)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)

		thresholdStateObject := []map[string]interface{}{
			{
				SloAlertConfigFieldThresholdType:     SloAlertThresholdType,
				SloAlertConfigFieldThresholdOperator: SloAlertThresholdOperator,
				SloAlertConfigFieldThresholdValue:    SloAlertThresholdValue,
			},
		}
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldThreshold, thresholdStateObject)

		thresholdType, type_ok := resourceData.GetOk("threshold.0.type")
		thresholdOperator, operator_ok := resourceData.GetOk("threshold.0.operator")
		thresholdValue, value_ok := resourceData.GetOk("threshold.0.value")

		require.True(t, type_ok, "threshold.0.type should exist")
		require.True(t, operator_ok, "threshold.0.operator should exist")
		require.True(t, value_ok, "threshold.0.value should exist")

		require.Equal(t, SloAlertThresholdType, thresholdType, "threshold.type should match set value")
		require.Equal(t, SloAlertThresholdOperator, thresholdOperator, "threshold.type should match set value")
		require.Equal(t, SloAlertThresholdValue, thresholdValue, "threshold.type should match set value")

	}
}

func (test *sloAlertConfigTest) shouldRequireTimeThresholdMetricsToMatchTheAssignedValues() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
		resourceHandle := NewSloAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloAlertID)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)

		timeThresholdStateObject := []map[string]interface{}{
			{
				SloAlertConfigFieldTimeThresholdWarmUp:   SloAlertTimeThresholdWarmUp,
				SloAlertConfigFieldTimeThresholdCoolDown: SloAlertTimeThresholdCoolDown,
			},
		}
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldTimeThreshold, timeThresholdStateObject)

		timeThresholdWarmUp, warm_ok := resourceData.GetOk("time_threshold.0.warm_up")
		timeThresholdCoolDown, cool_ok := resourceData.GetOk("time_threshold.0.cool_down")

		require.True(t, warm_ok, "time_threshold.0.warm_up should exist")
		require.True(t, cool_ok, "time_threshold.0.cool_down should exist")

		require.Equal(t, SloAlertTimeThresholdWarmUp, timeThresholdWarmUp, "time_threshold.0.warm_up should match set value")
		require.Equal(t, SloAlertTimeThresholdCoolDown, timeThresholdCoolDown, "time_threshold.0.cool_down should match set value")

	}
}

func (test *sloAlertConfigTest) shouldMapSmartAlertValuesCorrectly() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
		resourceHandle := NewSloAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloAlertID)

		// Set test values for a status alert
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, "status-alert-test")
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "status")
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldDescription, "Test status alert description")
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldSeverity, 5)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldTriggering, true)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldEnabled, true)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldSloIds, []interface{}{"slo-1", "slo-2"})
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertChannelIds, []interface{}{"channel-1", "channel-2"})
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldThreshold, []interface{}{
			map[string]interface{}{
				"type":     "staticThreshold",
				"operator": ">=",
				"value":    95.0,
			},
		})
		// Use schema field names (warm_up and cool_down) instead of API field names
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldTimeThreshold, []interface{}{
			map[string]interface{}{
				SloAlertConfigFieldTimeThresholdWarmUp:   300000,
				SloAlertConfigFieldTimeThresholdCoolDown: 60000,
			},
		})
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "status")

		apiObject, err := resourceHandle.MapStateToDataObject(resourceData)
		require.NoError(t, err, "mapping status alert should not fail with valid configuration")

		// Assert all assigned values are correct
		require.Equal(t, sloAlertID, apiObject.ID, "ID should match the set value")
		require.Equal(t, "status-alert-test", apiObject.Name, "name should match the set value")
		require.Equal(t, "Test status alert description", apiObject.Description, "description should match the set value")
		require.Equal(t, 5, apiObject.Severity, "severity should match the set value")
		require.Equal(t, true, apiObject.Triggering, "triggering should match the set value")
		require.Equal(t, true, apiObject.Enabled, "enabled should match the set value")

		require.ElementsMatch(t, []string{"slo-1", "slo-2"}, apiObject.SloIds, "sloIds should match the set values")
		require.ElementsMatch(t, []string{"channel-1", "channel-2"}, apiObject.AlertChannelIds, "alertChannelIds should match the set values")

		require.Equal(t, "staticThreshold", apiObject.Threshold.Type, "threshold type should match")
		require.Equal(t, ">=", apiObject.Threshold.Operator, "threshold operator should match")
		require.Equal(t, 95.0, apiObject.Threshold.Value, "threshold value should match")

		require.Equal(t, 300000, apiObject.TimeThreshold.TimeWindow, "time_threshold time_window should match")
		require.Equal(t, 60000, apiObject.TimeThreshold.Expiry, "time_threshold expiry should match")

		require.Equal(t, "SERVICE_LEVELS_OBJECTIVE", apiObject.Rule.AlertType, "rule alert_type should be 'SERVICE_LEVELS_OBJECTIVE'")
		require.Equal(t, "STATUS", apiObject.Rule.Metric, "rule metric should match")

		// Assert burn_rate_config is not set for status alert
		require.True(t, apiObject.BurnRateConfigs == nil || len(*apiObject.BurnRateConfigs) == 0, "burn_rate_config should be nil or empty for status alert")
	}
}

func (r *sloAlertConfigTest) shouldFailToMapAlertTypeWhenNoSupportedValueIsProvided() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
		resourceHandle := NewSloAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloAlertID)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)

		// Set an invalid alert_type
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "unknown-type")

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid alert_type: unknown-type")
	}
}

func (test *sloAlertConfigTest) shouldRequireBurnRateConfigToMatchAssignedValues() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
		resourceHandle := NewSloAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloAlertID)

		setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, SloAlertBurnRateV2)

		burnRateConfig := []map[string]interface{}{
			{
				"duration":           "1",
				"duration_unit_type": "hour",
				"alert_window_type":  "SINGLE",
				"threshold_operator": ">=",
				"threshold_value":    "1.0",
			},
		}
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldBurnRateConfig, burnRateConfig)

		raw, ok := resourceData.GetOk(SloAlertConfigFieldBurnRateConfig)
		require.True(t, ok, "burn_rate_config should exist")

		burnRateList := raw.([]interface{})
		require.Equal(t, 1, len(burnRateList), "burn_rate_config should contain 1 item")

		burnRateObj, ok := burnRateList[0].(map[string]interface{})
		require.True(t, ok, "burn_rate_config[0] should be a map")

		require.Equal(t, "1", burnRateObj["duration"], "duration should match")
		require.Equal(t, "hour", burnRateObj["duration_unit_type"], "duration_unit_type should match")
		require.Equal(t, "SINGLE", burnRateObj["alert_window_type"], "alert_window_type should match")
		require.Equal(t, ">=", burnRateObj["threshold_operator"], "threshold_operator should match")
		require.Equal(t, "1.0", burnRateObj["threshold_value"], "threshold_value should match")
	}
}

func (test *sloAlertConfigTest) shouldFailWhenBurnRateConfigIsMissingForBurnRateV2Alert() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
		resourceHandle := NewSloAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloAlertID)

		setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
		setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, SloAlertBurnRateV2)

		// burn_rate_config is intentionally NOT set here

		_, err := resourceHandle.MapStateToDataObject(resourceData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "burn_rate_config is required for alert_type 'burn_rate_v2'", "expected error when burn_rate_config is missing for burn_rate_v2 alert")
	}
}
