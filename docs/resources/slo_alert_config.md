# SLO Smart Alert Configuration Resource

Management of Smart Alerts for Service Level Objectives (SLOs) to trigger notifications when specific thresholds are surpassed, such as SLO status, the percentage of error budget consumed, or the error budget burn rate. Additionally, you can customize the alert by selecting one or more alert channels, adding custom payloads, or setting time-based thresholds.

API Documentation: <https://instana.github.io/openapi/#tag/Service-Levels-Alert-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v6, all nested configurations use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_slo_alert_config" "example" {
  name = "SLO Alert"
  alert_type = "status"
  
  threshold {
    operator = ">"
    value = 0.7
  }
  
  time_threshold {
    warm_up = 60000
    cool_down = 60000
  }
  
  custom_payload_field {
    key = "env"
    value = "prod"
  }
  
  custom_payload_field {
    key = "team"
    value = "ops"
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_slo_alert_config" "example" {
  name = "SLO Alert"
  alert_type = "status"
  severity              = 5
  slo_ids               = [var.slo_id]
  threshold = {
    type     = "staticThreshold"
    operator = ">"
    value = 0.7
  }
  
  time_threshold = {
    warm_up = 60000
    cool_down = 60000
  }
  
  custom_payload_fields = [
    {
      key = "env"
      value = "prod"
    },
    {
      key = "team"
      value = "ops"
    }
  ]
}
```

### Key Syntax Changes

1. **Threshold**: `threshold { }` → `threshold = { }`
2. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
3. **Burn Rate Config**: Multiple `burn_rate_config { }` blocks → Single `burn_rate_config = [{ }]` list
4. **Custom Payload Fields**: Multiple `custom_payload_field { }` blocks → Single `custom_payload_fields = [{ }]` list
5. **All nested objects**: Use `= { }` or `= [{ }]` syntax

## Example Usage

### Status Alert

Monitor SLO status and alert when it falls below a threshold:

```hcl
resource "instana_slo_alert_config" "status_alert" {
  name = "SLO Status Alert"
  description = "Alert when SLO status drops below 70%"
  severity = 10
  triggering = true  
  alert_type = "status"
  slo_ids = ["slo-id-1", "slo-id-2"]
  alert_channel_ids = ["channel-id-1"]
  
  threshold = {
    operator = ">"
    value = 0.7
  }
  
  time_threshold = {
    warm_up = 60000
    cool_down = 60000
  }
}
```

### Error Budget Alert

Monitor error budget consumption:

```hcl
resource "instana_slo_alert_config" "error_budget_alert" {
  name = "Error Budget Alert"
  description = "Alert when 50% of error budget is consumed"
  severity = 5
  triggering = true  
  alert_type = "error_budget"
  slo_ids = ["slo-id-1"]
  alert_channel_ids = ["channel-id-1", "channel-id-2"]
  
  threshold = {
    operator = ">="
    value = 0.5
  }
  
  time_threshold = {
    warm_up = 300000
    cool_down = 300000
  }
}
```

### Burn Rate Alert with Single Window

Monitor burn rate with a single alerting window:

```hcl
resource "instana_slo_alert_config" "burn_rate_single" {
  name = "Burn Rate Alert - Single Window"
  description = "Alert when burn rate exceeds threshold"
  severity = 10
  triggering = true  
  alert_type = "burn_rate_v2"
  slo_ids = ["slo-id-1"]
  alert_channel_ids = ["pagerduty-channel"]
  
  burn_rate_config = [
    {
      alert_window_type  = "SINGLE"
      duration           = "1"
      duration_unit_type = "minute"
      threshold_operator = ">"
      threshold_value    = "3.00"
    }
  ]
  
  time_threshold = {
    warm_up = 60000
    cool_down = 60000
  }
}
```
## Argument Reference

* `name` - Required - The name of the SLO Alert configuration (max 256 characters)
* `description` - Required - The description of the SLO Alert configuration
* `severity` - Required - The severity of the alert when triggered. Must be `5` for warning or `10` for critical
* `alert_type` - Required - The type of Smart Alert. Allowed values: `status`, `error_budget`, `burn_rate_v2`
* `slo_ids` - Required - A set of SLO IDs to monitor. Must contain at least one ID
* `alert_channel_ids` - Required - A set of alert channel IDs to send notifications to
* `triggering` - Optional - Flag to indicate whether to trigger an incident. Default: `false`
* `threshold` - Optional - Configuration block defining the threshold for the alert condition. Required for `status` and `error_budget` alert types [Details](#threshold-reference)
* `time_threshold` - Required - Configuration block defining the time threshold for triggering and suppressing alerts [Details](#time-threshold-reference)
* `burn_rate_config` - Optional - List of burn rate configurations and alerting windows. Required for `alert_type` set to `burn_rate_v2` [Details](#burn-rate-config-reference)
* `custom_payload_fields` - Optional - List of custom payload fields to include in the alert notification [Details](#custom-payload-fields-reference)

### Threshold Reference

The alert is triggered when the threshold is evaluated by the value and operator.

* `type` - Optional - The type of threshold. Default: `staticThreshold`. Allowed values: `staticThreshold`
* `operator` - Required - The operator used to evaluate the threshold. Allowed values: `>`, `>=`, `=`, `<=`, `<`
* `value` - Required - The threshold value for the alert condition (float)

### Time Threshold Reference

If the alert is triggered, after the warm-up period, a notification will be generated and transmitted to the alert channel. After notification is generated, the Instana event remains until the violation ends for a cool-down period.

* `warm_up` - Required - The duration (in milliseconds) for which the condition must be violated before the alert is triggered. Must be at least 1
* `cool_down` - Required - The duration (in milliseconds) for which the condition must remain suppressed before the alert ends. Must be at least 1

### Burn Rate Config Reference

The `burn_rate_config` block is applicable only for burn rate alerts (i.e., when `alert_type` = `burn_rate_v2`). This setting is required in such cases.

Currently, two types of burn rate alert configurations are supported:
- **Single Window, Single Threshold**: Uses a single threshold and a single alerting window
- **Multiple Windows, Multiple Thresholds**: Uses both short and long alerting windows with respective thresholds. An alert is triggered if *both* thresholds are breached (AND condition)

Each burn rate configuration object contains:

* `alert_window_type` - Required - Determines the type of burn rate alert. Allowed values:
  * `SINGLE` - Defines single alerting window
  * `SHORT` - Defines short alerting window
  * `LONG` - Defines long alerting window
* `duration` - Required - Duration of the alerting window (string format integer). Must be greater than 0
* `duration_unit_type` - Required - The unit of time for the duration. Allowed values: `minute`, `hour`, `day`
* `threshold_operator` - Required - Comparison operator. Allowed values: `>=`, `>`, `<`, `<=`
* `threshold_value` - Required - Numeric threshold value of the alerting window (string format)

### Custom Payload Fields Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - The dynamic value of the custom payload field. Either `value` or `dynamic_value` must be defined [Details](#dynamic-custom-payload-field-value)

#### Dynamic Custom Payload Field Value

* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload

## Attributes Reference

* `id` - The ID of the SLO alert configuration

## Import

SLO alert configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_slo_alert_config.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Use `status` alert type to monitor SLO compliance status
* Use `error_budget` alert type to monitor error budget consumption percentage
* Use `burn_rate_v2` alert type to monitor the rate at which error budget is being consumed
* Burn rate alerts with multiple windows trigger only when ALL configured windows breach their thresholds (AND logic)
* The `warm_up` period prevents alert flapping by requiring sustained violations
* The `cool_down` period prevents premature alert closure
* Custom payload fields can include both static values and dynamic tag values
* Severity level 5 is for warnings, severity level 10 is for critical alerts
* Multiple SLOs can be monitored with a single alert configuration
* The `triggering` flag controls whether the alert creates incidents in Instana
