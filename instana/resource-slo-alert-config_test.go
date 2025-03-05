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


func TestSliConfigTest(t *testing.T) {
	unitTest := &sliConfigUnitTest{}

	}
