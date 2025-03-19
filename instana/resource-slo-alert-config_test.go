package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

const (
 sloAlertConfigDefinition                  = "instana_slo_alert_config.example_slo_alert_config"
 sloAlertID                                = "id"
 sloAlertName                              = "slo-alert-name"
 sloAlertDescription                       = "slo-alert-description"
 sloAlertSeverity                          = 10
 SloAlertTriggering                        = true
 SloAlertConfigFieldRule 				   = "rule"
 sloAlertAlertType                         = "status"
 sloAlertAlertEnabled                      = true
 SloAlertThresholdType                     = "staticThreshold"
 SloAlertThresholdOperator                 = ">"
 SloAlertThresholdValue                    = 0.3
 SloAlertTimeThresholdWarmUp               = 60000
 SloAlertTimeThresholdCoolDown             = 30000
 SloAlertBurnRateShortDuration             = 5
 SloAlertBurnRateShortDurationType         = "minute"
 SloAlertBurnRateLongDuration              = 24
 SloAlertBurnRateLongDurationType          = "hour"
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
	t.Run("Should require Burn Rate Time Windows Values to Match the Assigned Values", test.shouldRequireBurnRateTimeWindowsToMatchAssignedValues())
	t.Run("Should Fail When Burn Rate Time Windows Is Missing For Burn Rate Alert", test.shouldFailWhenBurnRateTimeWindowsIsMissingForBurnRateAlert())
	t.Run("Should Fail When Short Time Window Duration Is Missing", test.shouldFailWhenShortTimeWindowDurationIsMissing())
	t.Run("Should Fail When Long Time Window Duration Is Missing", test.shouldFailWhenLongTimeWindowDurationIsMissing())
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
				SloAlertConfigFieldThresholdType:       SloAlertThresholdType,
				SloAlertConfigFieldThresholdOperator: 	SloAlertThresholdOperator,
				SloAlertConfigFieldThresholdValue:    	SloAlertThresholdValue,
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
				SloAlertConfigFieldTimeThresholdWarmUp  : SloAlertTimeThresholdWarmUp,
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

        // Assert burn_rate_time_windows is not set for status alert
        require.Nil(t, apiObject.BurnRateTimeWindows, "burn_rate_time_windows should be nil for status alert")
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

func (test *sloAlertConfigTest) shouldRequireBurnRateTimeWindowsToMatchAssignedValues() func(t *testing.T) {
    return func(t *testing.T) {
        testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
        resourceHandle := NewSloAlertConfigResourceHandle()
        resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
        resourceData.SetId(sloAlertID)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "burn_rate")

        // Set burn_rate_time_windows
        burnRateTimeWindowsStateObject := []map[string]interface{}{
            {
                SloAlertConfigFieldShortTimeWindow: []map[string]interface{}{
                    {
                        SloAlertConfigFieldTimeWindowDuration:     SloAlertBurnRateShortDuration,
                        SloAlertConfigFieldTimeWindowDurationType: SloAlertBurnRateShortDurationType,
                    },
                },
                SloAlertConfigFieldLongTimeWindow: []map[string]interface{}{
                    {
                        SloAlertConfigFieldTimeWindowDuration:     SloAlertBurnRateLongDuration,
                        SloAlertConfigFieldTimeWindowDurationType: SloAlertBurnRateLongDurationType,
                    },
                },
            },
        }
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldBurnRateTimeWindows, burnRateTimeWindowsStateObject)

        shortDuration, shortDurOk := resourceData.GetOk("burn_rate_time_windows.0.short_time_window.0.duration")
        shortDurationType, shortTypeOk := resourceData.GetOk("burn_rate_time_windows.0.short_time_window.0.duration_type")

        longDuration, longDurOk := resourceData.GetOk("burn_rate_time_windows.0.long_time_window.0.duration")
        longDurationType, longTypeOk := resourceData.GetOk("burn_rate_time_windows.0.long_time_window.0.duration_type")

        require.True(t, shortDurOk, "burn_rate_time_windows.0.short_time_window.0.duration should exist")
        require.True(t, shortTypeOk, "burn_rate_time_windows.0.short_time_window.0.duration_type should exist")
        require.True(t, longDurOk, "burn_rate_time_windows.0.long_time_window.0.duration should exist")
        require.True(t, longTypeOk, "burn_rate_time_windows.0.long_time_window.0.duration_type should exist")

        require.Equal(t, SloAlertBurnRateShortDuration, shortDuration.(int), "burn_rate_time_windows.0.short_time_window.0.duration should match set value")
        require.Equal(t, SloAlertBurnRateShortDurationType, shortDurationType.(string), "burn_rate_time_windows.0.short_time_window.0.duration_type should match set value")
        require.Equal(t, SloAlertBurnRateLongDuration, longDuration.(int), "burn_rate_time_windows.0.long_time_window.0.duration should match set value")
        require.Equal(t, SloAlertBurnRateLongDurationType, longDurationType.(string), "burn_rate_time_window.0.long_time_window.0.duration_type should match set value")
    }
}

func (test *sloAlertConfigTest) shouldFailWhenBurnRateTimeWindowsIsMissingForBurnRateAlert() func(t *testing.T) {
    return func(t *testing.T) {
        testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
        resourceHandle := NewSloAlertConfigResourceHandle()
        resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
        resourceData.SetId(sloAlertID)

        setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "burn_rate")

        _, err := resourceHandle.MapStateToDataObject(resourceData)
        require.Error(t, err)
        require.Contains(t, err.Error(), "burn_rate_time_windows is required for alert_type 'burn_rate'", "expected error when burn_rate_time_windows is missing for burn_rate alert")
    }
}

func (test *sloAlertConfigTest) shouldFailWhenShortTimeWindowDurationIsMissing() func(t *testing.T) {
    return func(t *testing.T) {
        testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
        resourceHandle := NewSloAlertConfigResourceHandle()
        resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
        resourceData.SetId(sloAlertID)

        // Set required fields
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "burn_rate")
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldDescription, "Test burn rate alert")
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldSeverity, 5)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldTriggering, true)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldEnabled, true)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldSloIds, []interface{}{"slo-1"})
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertChannelIds, []interface{}{"channel-1"})

        setValueOnResourceData(t, resourceData, SloAlertConfigFieldThreshold, []interface{}{
            map[string]interface{}{
                "type":     "staticThreshold",
                "operator": ">=",
                "value":    95.0,
            },
        })
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldTimeThreshold, []interface{}{
            map[string]interface{}{
                SloAlertConfigFieldTimeThresholdWarmUp:   300000,
                SloAlertConfigFieldTimeThresholdCoolDown: 60000,
            },
        })

        // Set burn_rate_time_windows with missing short_time_window duration
        burnRateTimeWindowsStateObject := []map[string]interface{}{
            {
                SloAlertConfigFieldShortTimeWindow: []map[string]interface{}{ 
                    map[string]interface{}{ 
                        SloAlertConfigFieldTimeWindowDurationType: "minute",
                    },
                },
                SloAlertConfigFieldLongTimeWindow: []map[string]interface{}{
                    {
                        SloAlertConfigFieldTimeWindowDuration:     24,
                        SloAlertConfigFieldTimeWindowDurationType: "hour",
                    },
                },
            },
        }
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldBurnRateTimeWindows, burnRateTimeWindowsStateObject)

        _, err := resourceHandle.MapStateToDataObject(resourceData)
        require.Error(t, err)
        require.Contains(t, err.Error(), "duration in short_time_window must be an integer", "expected error when short_time_window duration is missing")
    }
}

func (test *sloAlertConfigTest) shouldFailWhenLongTimeWindowDurationIsMissing() func(t *testing.T) {
    return func(t *testing.T) {
        testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
        resourceHandle := NewSloAlertConfigResourceHandle()
        resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
        resourceData.SetId(sloAlertID)

        // Set required fields
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "burn_rate")
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldDescription, "Test burn rate alert")
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldSeverity, 5)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldTriggering, true)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldEnabled, true)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldSloIds, []interface{}{"slo-1"})
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertChannelIds, []interface{}{"channel-1"})
		
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldThreshold, []interface{}{
            map[string]interface{}{
                "type":     "staticThreshold",
                "operator": ">=",
                "value":    95.0,
            },
        })
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldTimeThreshold, []interface{}{
            map[string]interface{}{
                SloAlertConfigFieldTimeThresholdWarmUp:   300000,
                SloAlertConfigFieldTimeThresholdCoolDown: 60000,
            },
        })

        // Set burn_rate_time_windows with missing short_time_window duration
        burnRateTimeWindowsStateObject := []map[string]interface{}{
            {
                SloAlertConfigFieldShortTimeWindow: []map[string]interface{}{ 
                    map[string]interface{}{ 
						SloAlertConfigFieldTimeWindowDuration:     60,
                        SloAlertConfigFieldTimeWindowDurationType: "minute",
                    },
                },
                SloAlertConfigFieldLongTimeWindow: []map[string]interface{}{
                    {
                        SloAlertConfigFieldTimeWindowDurationType: "hour",
                    },
                },
            },
        }
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldBurnRateTimeWindows, burnRateTimeWindowsStateObject)

        _, err := resourceHandle.MapStateToDataObject(resourceData)
        require.Error(t, err)
        require.Contains(t, err.Error(), "duration in long_time_window must be an integer", "expected error when long_time_window duration is missing")
    }
}