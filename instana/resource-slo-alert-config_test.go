package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
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
