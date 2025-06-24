# SLO Smart Alert Configuration Resource

Management of Smart Alerts for Service Level Objectives (SLOs) to trigger notifications when specific thresholds are surpassed, 
such as SLO status, the percentage of error budget consumed, or the error budget burn rate. Additionally, you can customize 
the alert by selecting one or more alert channels, adding custom payloads, or setting time-based thresholds.

API Documentation: <https://instana.github.io/openapi/#tag/Service-Levels-Alert-Configuration>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage
Creating a `Status Alert`

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
Creating an `Error Budget Alert`

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

Creating a `Burn Rate alert with single alerting window`

```hcl
resource "instana_slo_alert_config" "burn_rate_alert" {
    name = "terraform_burn_rate_v2_alert_slo"
    description = "Burn rate is exceeded"
    severity = 10
    triggering = true
    alert_type = "burn_rate_v2"

    burn_rate_config = [
      {
        alert_window_type   = "SINGLE"
        duration            = "6"
        duration_unit_type  = "hour"
        threshold_operator  = ">="
        threshold_value     = "1"
      }
    ]

    slo_ids = ["SLOjRzdfwyLQ6KDWoq1pTvoeA", "SLOXj71OScUQ9GZWGTInq15Pg "]
    alert_channel_ids = ["orhurggugksjfgh "]

    time_threshold {
        warm_up = 60000
        cool_down = 60000
    }

    enabled = true
}
``` 

Creating a `Burn Rate alert with multiple alerting windows`

```hcl
resource "instana_slo_alert_config" "burn_rate_alert" {
    name = "terraform_burn_rate_v2_alert_slo"
    description = "Burn rate is exceeded"
    severity = 10
    triggering = true
    alert_type = "burn_rate_v2"

    burn_rate_config = [
      {
        alert_window_type   = "LONG"
        duration            = "6"
        duration_unit_type  = "hour"
        threshold_operator  = ">="
        threshold_value     = "1"
      }, 
    {
        alert_window_type   = "SHORT"
        duration            = "1"
        duration_unit_type  = "hour"
        threshold_operator  = ">="
        threshold_value     = "4"
      }
    ]

    slo_ids = ["SLOjRzdfwyLQ6KDWoq1pTvoeA", "SLOXj71OScUQ9GZWGTInq15Pg "]
    alert_channel_ids = ["orhurggugksjfgh "]

    time_threshold {
        warm_up = 60000
        cool_down = 60000
    }

    enabled = true
}
``` 

## Argument Reference
Ths SLO smart alert could be configured with the following arguments:
* `name` - Required - The name of the SLO Alert configuration.
* `description` - Required - The description of the SLO Alert configuration.
* `severity` - Required - The severity of the alert when triggered. Must be set to `5` for a warning alert level or `10` for a critical alert level.
* `alert_type` - Required - The type of Smart Alert. Allowed values: `status`, `error_budget`, `burn_rate_v2`. Defines what to alert on (e.g., SLO status, error budget percentage, or burn rate).
* `slo_ids` - Required - A set of SLO IDs to monitor. Must contain at least one ID.
* `triggering` - Optional - Flag to indicate whether an incident. Must be a boolean. Defaults to `false`.
* `enabled` - Optional - Flag to indicate whether the alert is enabled. Must be a boolean. Defaults to `false`.
* `alert_channel_ids` - Required - A set of alert channel IDs to send notifications to.
* `custom_payload_fields` - Optional - A list of custom payload fields to include in the alert notification.
* `threshold` - Required - A resource block defining the threshold for the alert condition. [Details](#threshold-reference)
* `time_threshold` - Required - A resource block defining the time threshold for triggering and suppressing alerts. [Details](#time-threshold-reference)
* `burn_rate_config` - Optional - A resource block defining the burn rate config and alerting windows for evaluating alert conditions. Required for `alert_type` set to `burn_rate_v2`. [Details](#burn-rate-config-reference)

### Threshold Reference
The alert is triggered when the threshold is evaluated by the value and operator. 
* `value` - Required - The threshold value for the alert condition. Must be a float.
* `operator` - Required - The operator used to evaluate the threshold. Supported operators depend on the Instana API (e.g., `>`, `<`, `>=`, `<=`, `=`).

### Time Threshold Reference
If the alert is triggered, after the warm up period, a notification will be generated and transmit to the alert channel. After notification generated, Instana event remains until violation ends for a cool down period. The time threshold is to configurate the warm up and cool down period settings.
* `warm_up` - Required - The duration (in milliseconds) for which the condition must be violated before the alert is triggered.
* `cool_down` - Required - The duration (in milliseconds) for which the condition must remain suppressed before the alert ends.

### Burn Rate Config Reference
The burn_rate_config block is applicable only for burn rate alerts (i.e., when `alert_type` = `burn_rate_v2`). This setting is required in such cases.
Currently, two types of burn rate alert configurations are supported:
- **Single Window, Single Threshold**: Uses a single threshold and a single alerting window. [Details](#burn-rate-config-fields-reference)
- **Multiple Windows, Multiple Thresholds**: Uses both short and long alerting windows with respective thresholds. An alert is triggered if *both* threshold is breached (AND condition).[Details](#burn-rate-config-fields-reference)
  
### Burn Rate Config Fields Reference

- `alert_window_type` – Required – Determines the type of burn rate alert. Allowed values:
  - `SINGLE`: Defines single alerting window.
  - `SHORT`: Defines short alerting window.
  - `LONG`: Defines long alerting window. 

- `duration` – Required – Duration of the alerting window in integer format. Must be greater than 0.
- `duration_unit_type` – Required – The unit of time for the duration. Allowed values: `minute`, `hour`, `day`.
- `threshold_value` - Required - Numeric threshold value of the alerting window.
- `threshold_operator` - Required – Comparison operator. E.g., `>=`, `>`, `<` and `<=`.



