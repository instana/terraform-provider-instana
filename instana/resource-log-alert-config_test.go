package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

const resourceLogAlertConfigDefinition = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:%d"
  tls_skip_verify = true
}

resource "instana_log_alert_config" "example" {
  name        = "name %d"
  description = "description"
  
  tag_filter = "entity.type@dest EQUALS 'host'"
  granularity = 600000
  
  alert_channels {
    warning  = ["channel-id-1"]
    critical = ["channel-id-2"]
  }
  
  rules {
    metric_name       = "log.count"
    alert_type        = "log.count"
    aggregation       = "SUM"
    threshold_operator = ">"
    
    threshold {
      warning {
        static {
          value = 100
        }
      }
      critical {
        static {
          value = 500
        }
      }
    }
  }
  
  time_threshold {
    violations_in_sequence {
      time_window = 600000
    }
  }
}
`

const logAlertConfigServerResponseTemplate = `
{
  "id": "%s",
  "name": "name %d",
  "description": "description",
  "tagFilterExpression": {
    "type": "TAG_FILTER",
    "name": "entity.type",
    "stringValue": "host",
    "numberValue": null,
    "booleanValue": null,
    "key": null,
    "value": "host",
    "operator": "EQUALS",
    "entity": "DESTINATION"
  },
  "granularity": 600000,
  "alertChannels": {
    "WARNING": ["channel-id-1"],
    "CRITICAL": ["channel-id-2"]
  },
  "rules": [
    {
      "thresholdOperator": ">",
      "rule": {
        "alertType": "logCount",
        "metricName": "log.count",
        "aggregation": "SUM"
      },
      "thresholds": {
        "WARNING": {
          "type": "staticThreshold",
          "value": 100
        },
        "CRITICAL": {
          "type": "staticThreshold",
          "value": 500
        }
      }
    }
  ],
  "timeThreshold": {
    "type": "violationsInSequence",
    "timeWindow": 600000
  },
  "created": 1647679325301,
  "readOnly": false,
  "enabled": true
}
`

func TestLogAlertConfigResource(t *testing.T) {
	id := RandomID()
	resourceRestAPIPath := restapi.LogAlertConfigResourcePath
	resourceInstanceRestAPIPath := resourceRestAPIPath + "/{internal-id}"

	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		config := &restapi.LogAlertConfig{}
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
		jsonData := fmt.Sprintf(logAlertConfigServerResponseTemplate, id, modCount)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(jsonData))
		if err != nil {
			fmt.Printf("failed to write response; %s\n", err)
		}
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceName := "instana_log_alert_config.example"
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(resourceLogAlertConfigDefinition, httpServer.GetPort(), 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", id),
					resource.TestCheckResourceAttr(resourceName, "name", "name 0"),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tag_filter", "entity.type@dest EQUALS 'host'"),
					resource.TestCheckResourceAttr(resourceName, "granularity", "600000"),
					resource.TestCheckResourceAttr(resourceName, "alert_channels.0.warning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alert_channels.0.warning.0", "channel-id-1"),
					resource.TestCheckResourceAttr(resourceName, "alert_channels.0.critical.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alert_channels.0.critical.0", "channel-id-2"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.metric_name", "log.count"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.alert_type", "log.count"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.aggregation", "SUM"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.threshold_operator", ">"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.threshold.0.warning.0.static.0.value", "100"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.threshold.0.critical.0.static.0.value", "500"),
					resource.TestCheckResourceAttr(resourceName, "time_threshold.0.violations_in_sequence.0.time_window", "600000"),
				),
			},
			testStepImportWithCustomID(resourceName, id),
		},
	})
}

func TestLogAlertConfigSchemaDefinition(t *testing.T) {
	resourceHandle := LogAlertConfigResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	assert.NotNil(t, schemaMap[LogAlertConfigFieldName])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldDescription])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldTagFilter])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldGranularity])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldAlertChannels])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldRules])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldTimeThreshold])
	assert.NotNil(t, schemaMap[LogAlertConfigFieldGroupBy])
	assert.NotNil(t, schemaMap[DefaultCustomPayloadFieldsName])
}
