# SLO Smart Alert Configuration Resource

Management of Smart Alerts for Service Level Objectives (SLOs) to trigger notifications when specific thresholds are surpassed, 
such as SLO status, the percentage of error budget consumed, or the error budget burn rate. Additionally, you can customize 
the alert by selecting one or more alert channels, adding custom payloads, or setting time-based thresholds.

API Documentation: <https://instana.github.io/openapi/#tag/Service-Levels-Alert-Configuration>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage
Creating a Status Alert

```hcl
resource "instana_slo_alert_config" "status_alert" {
  name         = "terraform_status_alert"
  description  = "terraform_status_alert testing"
  severity     = 5
  triggering   = true
  slo_ids           = ["SLOjRzdfwyLQ6KDWoq1pTvoeA", "SLOXj71OScUQ9GZWGTInq15Pg"]
  alert_channel_ids = ["orhurugksjfgh"]

  alert_type   = "status"
  threshold {
    operator = ">"
    value    = 0.7
  }
  time_threshold {
    warm_up     = 6000
    cool_down   = 6000
  }

  custom_payload_field {
    key    = "test1"
    value  = "foo"
  }

  enabled  = true
}
``` 
Creating an Error Budget Alert

```hcl
resource "instana_slo_alert_config" "error_budget_alert" {

  name = "terraform_error_budget_alert"
  description = "Consumed >= 3% of the error budget."
  severity = 10
  triggering = true
  slo_ids           = ["SLOjRzdfwyLQ6KDWoq1pTvoeA", "SLOXj71OScUQ9GZWGTInq15Pg"]
  alert_channel_ids = ["orhurugksjfgh"]

  alert_type = "error_budget"
  threshold {
    operator = ">"
    value    = 0.5
  }
  time_threshold {
    warm_up     = 60000
    cool_down   = 60000
  }

  custom_payload_field {
    key    = "test"
    value  = "foo"
  }

  enabled  = true
}
```
Creating a Burn Rate Alert

```hcl
resource "instana_slo_alert_config" "burn_rate_alert" {
  name        = "terraform"
  description = "3% of of the error budget."
  severity    = 5
  triggering  = true   
  slo_ids           = ["SLOjRzdfwyLQ6KDWoq1pTvoeA", "SLOXj71OScUQ9GZWGTInq15Pg"]
  alert_channel_ids = ["orhurggugksjfgh"]

  alert_type  = "burn_rate"
  threshold {
    operator = ">"
    value    = 1
  }
  time_threshold {
    warm_up     = 60000
    cool_down   = 60000
  }

  burn_rate_time_windows {
    long_time_window {
      duration     = 1
      duration_type = "day"
    }

    short_time_window {
      duration     = 30
      duration_type = "minute"
    }
  }

  custom_payload_field {
    key    = "test"
    value  = "foo"
  }

  enabled = true
}
``` 

## Argument Reference

* `name` - Required - The name of the SLO Alert configuration.
* `description` - Required - The description of the SLO Alert configuration.
* `severity` - Required - The severity of the alert when triggered. Must be set to `5` for a warning alert level or `10` for a critical alert level.
* `triggering` - Optional - Flag to indicate whether an incident is also triggered. Must be a boolean. Defaults to `false`.
* `alert_type` - Required - The type of Smart Alert. Allowed values: `status`, `error_budget`, `burn_rate`. Defines what to alert on (e.g., SLO status, error budget percentage, or burn rate).
* `slo_ids` - Required - A set of SLO IDs to monitor. Must contain at least one ID.
* `alert_channel_ids` - Required - A set of alert channel IDs to send notifications to.
* `enabled` - Optional - Flag to indicate whether the alert is enabled. Must be a boolean. Defaults to `false`.
* `custom_payload_fields` - Optional - A list of custom payload fields to include in the alert notification.
* `threshold` - Required - A resource block defining the threshold for the alert condition. [Details](#threshold-reference)
* `time_threshold` - Required - A resource block defining the time threshold for triggering and suppressing alerts. [Details](#time-threshold-reference)
* `burn_rate_time_windows` - Optional - A resource block defining the burn rate time windows for evaluating alert conditions. Required for `alert_type` set to `burn_rate`. [Details](#burn-rate-time-windows-reference)

### Threshold Reference

* `type` - Optional - The type of threshold. Must be `staticThreshold`. Defaults to `staticThreshold`.
* `operator` - Required - The operator used to evaluate the threshold. Supported operators depend on the Instana API (e.g., `>`, `<`, `>=`, `<=`, `=`).
* `value` - Required - The threshold value for the alert condition. Must be a float.

### Time Threshold Reference

* `warm_up` - Required - The duration (in milliseconds) for which the condition must be violated before the alert is triggered.
* `cool_down` - Required - The duration (in milliseconds) for which the condition must remain suppressed before the alert ends.

### Burn Rate Time Windows Reference

* `long_time_window` - Optional - A resource block defining the long time window duration and type. Required if `alert_type` is `burn_rate`. [Details](#long-time-window-reference)
* `short_time_window` - Optional - A resource block defining the short time window duration and type. Required if `alert_type` is `burn_rate`. [Details](#short-time-window-reference)

#### Long Time Window Reference

* `duration` - Required - The duration for the long time window. Must be an integer greater than 0.
* `duration_type` - Required - The unit of time for the long time window duration. Allowed values: `minute`, `hour`, `day`.

#### Short Time Window Reference

* `duration` - Required - The duration for the short time window. Must be an integer greater than 0.
* `duration_type` - Required - The unit of time for the short time window duration. Allowed values: `minute`, `hour`, `day`.

