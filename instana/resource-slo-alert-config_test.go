package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

const (
 sloAlertConfigDefinition                  = "instana_sli_config.example_sli_config"
 sloAlertID                                = "id"
 sloAlertName                              = "slo-alert-name"
 sloAlertDescription                       = "slo-alert-description"
 sloAlertSeverity                          = 10
 SloAlertTriggering                        = true
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
	t.Run("Should require Burn Rate Time Windows Values to Match the Assigned Values", test.shouldRequireBurnRateTimeWindowsToMatchAssignedValues())
	// t.Run("Should Fail With Invalid Burn Rate Duration Type", test.shouldFailWhenBurnRateTimeWindowsHasInvalidDurationType())
	t.Run("Should Fail When Burn Rate Time Windows Is Missing For Burn Rate Alert", test.shouldFailWhenBurnRateTimeWindowsIsMissingForBurnRateAlert())
	t.Run("Should Fail When Short Time Window Duration Is Missing", test.shouldFailWhenShortTimeWindowDurationIsMissing())


	
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

// func (test *sloAlertConfigTest) shouldFailWhenBurnRateTimeWindowsHasInvalidDurationType() func(t *testing.T) {
//     return func(t *testing.T) {
//         testHelper := NewTestHelper[*restapi.SloAlertConfig](t)
//         resourceHandle := NewSloAlertConfigResourceHandle()
//         resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
//         resourceData.SetId(sloAlertID)
//         setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)

//         // Add required fields to avoid unrelated errors
//         setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "burn_rate")

//         // Set burn_rate_time_windows with invalid duration_type
//         burnRateTimeWindowsStateObject := []map[string]interface{}{
//             {
//                 SloAlertConfigFieldShortTimeWindow: []map[string]interface{}{
//                     {
//                         SloAlertConfigFieldTimeWindowDuration:     5,
//                         SloAlertConfigFieldTimeWindowDurationType: "invalid", // Invalid value
//                     },
//                 },
//                 SloAlertConfigFieldLongTimeWindow: []map[string]interface{}{
//                     {
//                         SloAlertConfigFieldTimeWindowDuration:     24,
//                         SloAlertConfigFieldTimeWindowDurationType: "invalid", // Invalid value
//                     },
//                 },
//             },
//         }
//         setValueOnResourceData(t, resourceData, SloAlertConfigFieldBurnRateTimeWindows, burnRateTimeWindowsStateObject)

//         // Attempt to map to API object, expect API rejection
//         _, err := resourceHandle.MapStateToDataObject(resourceData)
//         require.Error(t, err)
//         require.Contains(t, err.Error(), "invalid duration_type", "expected API rejection for invalid duration_type")
//     }
// }


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

        setValueOnResourceData(t, resourceData, SloAlertConfigFieldName, sloAlertName)
        setValueOnResourceData(t, resourceData, SloAlertConfigFieldAlertType, "burn_rate")

        burnRateTimeWindowsStateObject := []map[string]interface{}{
            {
                SloAlertConfigFieldShortTimeWindow: []map[string]interface{}{
                    {
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

        // Attempt to map to API object, expect rejection
        _, err := resourceHandle.MapStateToDataObject(resourceData)
        require.Error(t, err)
        require.Contains(t, err.Error(), "duration in short_time_window must be an integer", "expected error when short_time_window duration is missing")
    }
}