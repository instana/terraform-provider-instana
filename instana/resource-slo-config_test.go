package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

const (
	sloConfigDefinition = "instana_slo_config.example_slo_config"
	sloConfigID         = "id"
	sloName             = "slo-config-example-name"
	sloTarget           = 0.90
	sloThreshold        = 20
	sloAggregation      = "MEAN"
   )

func TestSloConfig(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaSloConfig + ".example"
	inst := &sloConfigTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewSloConfigResourceHandle(),
	}
	inst.run(t)
}

type sloConfigTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle[*restapi.SloConfig]
}


var sloConfigTerraformTemplate = `
resource "instana_slo_config" "app_1" {
	name = "tfslo_app_timebased_latency_fixed"
	target = 0.91
	tags = ["terraform", "app", "timebased", "latency", "fixed"]
	entity {
	  application {
		application_id = instana_application_config.myAllServices.id
		boundary_scope = "ALL"
		include_internal = false
		include_synthetic = false
		filter_expression = "AND"
	  }
	}
	indicator {
	   time_based_latency {
		 threshold = 13.1
		 aggregation = "MEAN"
	   }
	}
	time_window {
	  fixed {
		duration = 1
		duration_unit = "day"
		start_timestamp = var.fixed_timewindow_start_timestamp
	  }
	}
  }  
`

var sloConfigServerResponseTemplate = `
{
  "id": "SLOTFcv9c6mtbd0vlovhcn72g",
  "name": "tfslo_app_timebased_latency_fixed",
  "target": 0.91,
  "tags": [
    "terraform",
    "app",
    "timebased",
    "latency",
    "fixed"
  ],
  "entity": {
    "type": "application",
    "applicationId": "cv9c6mlbd0vlovhcn720",
    "serviceId": null,
    "endpointId": null,
    "boundaryScope": "ALL",
    "includeSynthetic": null,
    "includeInternal": null,
    "tagFilterExpression": null
  },
  "indicator": {
    "blueprint": "latency",
    "type": "timeBased",
    "threshold": 13.1,
    "aggregation": "MEAN"
  },
  "timeWindow": {
    "type": "fixed",
    "duration": 1,
    "durationUnit": "day",
    "startTimestamp": 1698552000000
  }
}
`

func (test *sloConfigTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("%s should have correct resouce name", ResourceInstanaSloConfig), test.createTestResourceShouldHaveResourceName())
	t.Run(fmt.Sprintf("%s should have one state upgrader", ResourceInstanaSloConfig), test.createTestResourceShouldHaveOneStateUpgrader())
	t.Run(fmt.Sprintf("%s should have schema version one", ResourceInstanaSloConfig), test.createTestResourceShouldHaveSchemaVersionOne())
	t.Run("SLO Threshold Should Be Greater Than Zero", test.thresholdShouldBeGreaterThanZero())
	t.Run("Application ID Should Be Required For Application Entity", test.applicationIDShouldBeRequiredForApplicationEntity())
	t.Run("Should Map Valid Slo Application Config To API Object", test.shouldMapApplicationSloConfigToAPIObject())
	t.Run("Should Map Valid Slo Website Config To API Object", test.shouldMapWebsiteSloConfigToAPIObject())


	
}

func (test *sloConfigTest) createTestResourceShouldHaveResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, test.resourceHandle.MetaData().ResourceName, "instana_slo_config")
	}
}

func (test *sloConfigTest) createTestResourceShouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Len(t, test.resourceHandle.StateUpgraders(), 1)
	}
}

func (test *sloConfigTest) createTestResourceShouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, test.resourceHandle.MetaData().SchemaVersion)
	}
}

func (r *sloConfigTest) thresholdShouldBeGreaterThanZero() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloConfig](t)
		resourceHandle := NewSloConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloConfigID)
		setValueOnResourceData(t, resourceData, SloConfigFieldName, sloName)
		setValueOnResourceData(t, resourceData, SloConfigFieldTarget, sloTarget)

		sloConfigIndicatorStateObject := []interface{}{
			map[string]interface{}{
				SloConfigTimeBasedLatencyIndicator: []interface{}{
					map[string]interface{}{
						SloConfigFieldThreshold  : 0,
						SloConfigFieldAggregation: sloAggregation,
					},
				},
			},
		}

		setValueOnResourceData(t, resourceData, SloConfigFieldSloIndicator, sloConfigIndicatorStateObject)

		_, metricThresholdIsOK:= resourceData.GetOk("indicator.0.time_based_latency.0.threshold")
		require.False(t, metricThresholdIsOK)
	}
}

func (r *sloConfigTest) shouldMapApplicationSloConfigToAPIObject() func(t *testing.T) {
    return func(t *testing.T) {
        testHelper := NewTestHelper[*restapi.SloConfig](t)
        resourceHandle := NewSloConfigResourceHandle()
        resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
        resourceData.SetId("test-slo-id")

        setValueOnResourceData(t, resourceData, instana.SloConfigFieldName, "Test SLO Config")
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldTarget, 99.9)
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldTags, []interface{}{"tag1", "tag2"})

        setValueOnResourceData(t, resourceData, instana.SloConfigFieldSloEntity, []interface{}{
            map[string]interface{}{
                instana.SloConfigApplicationEntity: []interface{}{
                    map[string]interface{}{
                        instana.SloConfigFieldApplicationID:       "app-123",
                        instana.SloConfigFieldBoundaryScope:       "ALL",
                        instana.SloConfigFieldFilterExpression:    "tag:env=prod",
                        instana.SloConfigFieldIncludeInternal:     false,
                        instana.SloConfigFieldIncludeSynthetic:    true,
                        instana.SloConfigFieldServiceID:           "service-id",
                        instana.SloConfigFieldEndpointID:          "endpoint-id",
                    },
                },
            },
        })

        // Indicator (time_based_latency)
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldSloIndicator, []interface{}{
            map[string]interface{}{
                instana.SloConfigTimeBasedLatencyIndicator: []interface{}{
                    map[string]interface{}{
                        "threshold":   20.0,
                        "aggregation": "MEAN",
                    },
                },
            },
        })

        // Time window (rolling)
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldSloTimeWindow, []interface{}{
            map[string]interface{}{
                "rolling": []interface{}{
                    map[string]interface{}{
                        "duration":      7,
                        "duration_unit": "day",
                    },
                },
            },
        })

        apiObject, err := resourceHandle.MapStateToDataObject(resourceData)
        require.NoError(t, err, "MapStateToDataObject should not return an error for valid input")

        // Log for debugging
        // t.Log("API Object", apiObject)

        require.Equal(t, "test-slo-id", apiObject.ID, "ID should match")
        require.Equal(t, "Test SLO Config", apiObject.Name, "Name should match")
        require.Equal(t, 99.9, apiObject.Target, "Target should match")
        require.Equal(t, []interface{}{"tag1", "tag2"}, apiObject.Tags, "Tags should match")

		// application entity
		require.NotNil(t, apiObject.Entity, "Entity should not be nil")
		entity, ok := apiObject.Entity.(restapi.SloApplicationEntity)
		require.True(t, ok, "Entity should be a SloApplicationEntity")
		require.Equal(t, instana.SloConfigApplicationEntity, entity.Type, "Entity Type should be application")
		require.NotNil(t, entity.ApplicationID, "ApplicationID should not be nil")
		require.Equal(t, "app-123", *entity.ApplicationID, "ApplicationID should match")
		require.NotNil(t, entity.BoundaryScope, "BoundaryScope should not be nil")
		require.Equal(t, "ALL", *entity.BoundaryScope, "BoundaryScope should match")
		require.NotNil(t, entity.ServiceID, "ServiceID should not be nil")
		require.Equal(t, "service-id", *entity.ServiceID, "ServiceID should match")
		require.NotNil(t, entity.EndpointID, "EndpointID should not be nil")
		require.Equal(t, "endpoint-id", *entity.EndpointID, "EndpointID should match")

        // Validate indicator
        require.NotNil(t, apiObject.Indicator, "Indicator should not be nil")
        indicator, ok := apiObject.Indicator.(restapi.SloTimeBasedLatencyIndicator)
        require.True(t, ok, "Indicator should be a SloTimeBasedLatencyIndicator")
        require.Equal(t, instana.SloConfigAPIIndicatorBlueprintLatency, indicator.Blueprint, "Blueprint should match")
        require.Equal(t, instana.SloConfigAPIIndicatorMeasurementTypeTimeBased, indicator.Type, "Type should match")
        require.Equal(t, 20.0, indicator.Threshold, "Threshold should match")
        require.Equal(t, "MEAN", indicator.Aggregation, "Aggregation should match")

        // Validate time window
        require.NotNil(t, apiObject.TimeWindow, "TimeWindow should not be nil")
        timeWindow, ok := apiObject.TimeWindow.(restapi.SloRollingTimeWindow)
        require.True(t, ok, "TimeWindow should be a SloRollingTimeWindow")
        require.Equal(t, instana.SloConfigRollingTimeWindow, timeWindow.Type, "Type should match")
        require.Equal(t, 7, timeWindow.Duration, "Duration should match")
        require.Equal(t, "day", timeWindow.DurationUnit, "DurationUnit should match")
    }
}

func (r *sloConfigTest) applicationIDShouldBeRequiredForApplicationEntity() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SloConfig](t)
		resourceHandle := NewSloConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sloConfigID)

		setValueOnResourceData(t, resourceData, SloConfigFieldName, sloName)
		setValueOnResourceData(t, resourceData, SloConfigFieldTarget, sloTarget)

		sloConfigEntityStateObject := []interface{}{
			map[string]interface{}{
				SloConfigApplicationEntity: []interface{}{
					map[string]interface{}{
						SloConfigFieldBoundaryScope: "ALL",
						SloConfigFieldApplicationID: nil,
					},
				},
			},
		}

		setValueOnResourceData(t, resourceData, SloConfigFieldSloEntity, sloConfigEntityStateObject)

		_, applicationIDIsOK := resourceData.GetOk("entity.0.application.0.application_id")

		require.False(t, applicationIDIsOK, "Application ID should be required for application entity")
	}
}

func (r *sloConfigTest) shouldMapWebsiteSloConfigToAPIObject() func(t *testing.T) {
    return func(t *testing.T) {
        testHelper := NewTestHelper[*restapi.SloConfig](t)
        resourceHandle := NewSloConfigResourceHandle()
        resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
        resourceData.SetId("test-slo-id")

        setValueOnResourceData(t, resourceData, instana.SloConfigFieldName, "Test SLO Config")
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldTarget, 99.9)
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldTags, []interface{}{"tag1", "tag2"})

		// website entity
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldSloEntity, []interface{}{
            map[string]interface{}{
                instana.SloConfigWebsiteEntity: []interface{}{
                    map[string]interface{}{

						instana.SloConfigFieldWebsiteID:        "web-123",
						instana.SloConfigFieldFilterExpression: "AND",
						instana.SloConfigFieldBeaconType:       "httpRequest",
                    },
                },
            },
        })

        // Indicator (time_based_latency)
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldSloIndicator, []interface{}{
            map[string]interface{}{
                instana.SloConfigTimeBasedLatencyIndicator: []interface{}{
                    map[string]interface{}{
                        "threshold":   20.0,
                        "aggregation": "MEAN",
                    },
                },
            },
        })

        // Time window (rolling)
        setValueOnResourceData(t, resourceData, instana.SloConfigFieldSloTimeWindow, []interface{}{
            map[string]interface{}{
                "rolling": []interface{}{
                    map[string]interface{}{
                        "duration":      7,
                        "duration_unit": "day",
                    },
                },
            },
        })

        apiObject, err := resourceHandle.MapStateToDataObject(resourceData)
        require.NoError(t, err, "MapStateToDataObject should not return an error for valid input")

        // Log for debugging
        // t.Log("API Object", apiObject)

        require.Equal(t, "test-slo-id", apiObject.ID, "ID should match")
        require.Equal(t, "Test SLO Config", apiObject.Name, "Name should match")
        require.Equal(t, 99.9, apiObject.Target, "Target should match")
        require.Equal(t, []interface{}{"tag1", "tag2"}, apiObject.Tags, "Tags should match")

		// Validate entity
		require.NotNil(t, apiObject.Entity, "Entity should not be nil")
		entity, ok := apiObject.Entity.(restapi.SloWebsiteEntity)
		require.True(t, ok, "Entity should be a SloWebsiteEntity")
		require.Equal(t, instana.SloConfigWebsiteEntity, entity.Type, "Entity Type should be website")
		require.NotNil(t, entity.WebsiteId, "WebsiteID should not be nil")
		require.Equal(t, "web-123", *entity.WebsiteId, "WebsiteID should match")
		require.NotNil(t, entity.BeaconType, "BeaconType should not be nil")
		require.Equal(t, "httpRequest", *entity.BeaconType, "BeaconType should match")

        // Validate indicator
        require.NotNil(t, apiObject.Indicator, "Indicator should not be nil")
        indicator, ok := apiObject.Indicator.(restapi.SloTimeBasedLatencyIndicator)
        require.True(t, ok, "Indicator should be a SloTimeBasedLatencyIndicator")
        require.Equal(t, instana.SloConfigAPIIndicatorBlueprintLatency, indicator.Blueprint, "Blueprint should match")
        require.Equal(t, instana.SloConfigAPIIndicatorMeasurementTypeTimeBased, indicator.Type, "Type should match")
        require.Equal(t, 20.0, indicator.Threshold, "Threshold should match")
        require.Equal(t, "MEAN", indicator.Aggregation, "Aggregation should match")

        // Validate time window
        require.NotNil(t, apiObject.TimeWindow, "TimeWindow should not be nil")
        timeWindow, ok := apiObject.TimeWindow.(restapi.SloRollingTimeWindow)
        require.True(t, ok, "TimeWindow should be a SloRollingTimeWindow")
        require.Equal(t, instana.SloConfigRollingTimeWindow, timeWindow.Type, "Type should match")
        require.Equal(t, 7, timeWindow.Duration, "Duration should match")
        require.Equal(t, "day", timeWindow.DurationUnit, "DurationUnit should match")
    }
}