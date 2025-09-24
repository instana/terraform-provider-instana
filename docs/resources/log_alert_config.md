# Log Alert Configuration Resource

This resource manages a log alert configuration in Instana.

## Example Usage

```hcl
# Example: Static threshold log alert with simplified schema
resource "instana_log_alert_config" "static_example" {
  name        = "Static Threshold Log Alert"
  description = "This is an example log alert with static thresholds"
  
  # Define the tag filter to scope the alert
  tag_filter = "log.message@na CONTAINS 'error' AND log.level@na EQUALS 'ERROR'"
  
  # Set the evaluation granularity (in milliseconds)
  granularity = 600000 # 10 minutes
  
  # Configure alert channels for different severity levels
  alert_channels {
    warning  = ["channel-id-1", "channel-id-2"]
    critical = ["channel-id-3"]
  }
  
  # Optional: Group results by a specific tag
  group_by {
    tag_name = "kubernetes.namespace.name"
  }
  
  # Define the log count rule with static thresholds (simplified schema)
  rules {
    metric_name       = "log.count"
    alert_type        = "log.count"
    aggregation       = "SUM"
    threshold_operator = ">"
    
    threshold {
      warning {
        type = "static_threshold"
        value = 100
      }
      critical {
        type = "static_threshold"
        value = 500
      }
    }
  }
  
  # Define the time threshold for violations
  time_threshold {
    violations_in_sequence {
      time_window = 600000 # 10 minutes
    }
  }
  
  # Optional: Add custom payload fields for alert notifications
  custom_payload_field {
    key  = "environment"
    type = "staticString"
    value = "production"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the log alert configuration.
* `description` - (Required) A description of the log alert configuration.
* `tag_filter` - (Optional) The tag filter expression used to scope the alert.
* `granularity` - (Optional) The evaluation granularity in milliseconds. Default is 600000 (10 minutes). Possible values: 60000 (1 minute), 300000 (5 minutes), 600000 (10 minutes), 900000 (15 minutes), 1200000 (20 minutes), 1800000 (30 minutes).
* `grace_period` - (Optional) The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.
* `alert_channels` - (Optional) Configuration block for alert notification channels.
  * `warning` - (Optional) List of alert channel IDs to notify for warning severity alerts.
  * `critical` - (Optional) List of alert channel IDs to notify for critical severity alerts.
* `alert_channel_ids` - (Optional, Deprecated) List of IDs of alert channels defined in Instana. Use `alert_channels` instead.
* `severity` - (Optional, Deprecated) The severity of the alert when triggered, which is either 5 (Warning), or 10 (Critical). Use threshold rules with severity levels instead.
* `threshold` - (Optional, Deprecated) Configuration block for threshold settings. Use `rules` instead.
  * `operator` - (Required) The operator which will be applied to evaluate the threshold. Possible values: >, >=, <, <=.
  * `static_threshold` - (Optional) Configuration block for static threshold.
    * `value` - (Required) The static threshold value to compare against.
* `group_by` - (Optional) Configuration block for grouping results.
  * `tag_name` - (Required) The tag name to group by.
  * `key` - (Optional) The key to group by.
* `rules` - (Required) Configuration block for alert rules (simplified schema).
  * `metric_name` - (Required) The metric name to monitor.
  * `alert_type` - (Optional) The type of the log alert rule. Only "log.count" is supported. Default is "log.count".
  * `aggregation` - (Optional) The aggregation method to use for the log alert. Only 'SUM' is supported. Default is 'SUM'.
  * `threshold_operator` - (Required) The operator to use for threshold comparison. Possible values: >, >=, <, <=.
  * `threshold` - (Required) Configuration block for thresholds.
    * `warning` - (Optional) Configuration block for warning severity threshold.
      * `type` - (Required) The type of threshold. Currently only "static_threshold" is supported.
      * `value` - (Required) The threshold value.
    * `critical` - (Optional) Configuration block for critical severity threshold.
      * `type` - (Required) The type of threshold. Currently only "static_threshold" is supported.
      * `value` - (Required) The threshold value.
* `time_threshold` - (Required) Configuration block for time threshold.
  * `violations_in_sequence` - (Required) Configuration block for violations in sequence.
    * `time_window` - (Optional) The time window in milliseconds.
* `custom_payload_field` - (Optional) Configuration block for custom payload fields in alert notifications.
  * `key` - (Required) The key of the custom payload field.
  * `type` - (Required) The type of the custom payload field. Possible values: staticString, dynamic.
  * `value` - (Required) The value of the custom payload field.

## Attribute Reference

* `id` - The ID of the log alert configuration.

## Import

Log alert configurations can be imported using the ID, e.g.,

```
$ terraform import instana_log_alert_config.example 12345678-1234-1234-1234-123456789012
```

## Note

Adaptive baseline thresholds are not currently active in production and should not be used.