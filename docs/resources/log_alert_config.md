# Log Alert Configuration Resource

This resource manages a log alert configuration in Instana.

## Example Usage

```hcl
# Example: Static threshold log alert with schema
resource "instana_log_alert_config" "example" {
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
  
  # Define the log count rule with static thresholds
  rules {
    metric_name       = "log.count"
    alert_type        = "log.count"
    aggregation       = "SUM"
    threshold_operator = ">"
    threshold {
      critical {
        static {
          value = 500
        }
      }
      warning {
        static {
          value = 100
        }
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
  # Optional: Add custom dynamic payload fields for alert notifications
  custom_payload_field {
    key   = "host.fqdn"
    dynamic_value {
      tag_name = "host.fqdn"
    }
  }
}
```

## Argument Reference

* `name` - (Required) The name of the log alert configuration.
* `description` - (Required) A description of the log alert configuration.
* `tag_filter` - (Required) The tag filter expression used to scope the alert [Details](#tag-filter-argument-reference).
* `granularity` - (Required) The evaluation granularity in milliseconds. Default is 600000 (10 minutes). Possible values: 60000 (1 minute), 300000 (5 minutes), 600000 (10 minutes), 900000 (15 minutes), 1200000 (20 minutes), 1800000 (30 minutes).
* `grace_period` - (Optional) The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.
* `alert_channels` - (Optional) Configuration block for alert notification channels.
  * `warning` - (Optional) List of alert channel IDs to notify for warning severity alerts.
  * `critical` - (Optional) List of alert channel IDs to notify for critical severity alerts.
* `group_by` - (Optional) Configuration block for grouping results.
  * `tag_name` - (Required) The tag name to group by.
  * `key` - (Optional) The key to group by.
* `rules` - (Required) Configuration block for alert rules.
  * `metric_name` - (Required) The metric name to monitor.
  * `alert_type` - (Optional) The type of the log alert rule. Only "log.count" is supported. Default is "log.count".
  * `aggregation` - (Optional) The aggregation method to use for the log alert. Only 'SUM' is supported. Default is 'SUM'.
  * `threshold_operator` - (Required) The operator to use for threshold comparison. Possible values: >, >=, <, <=.
  * `threshold` - (Required) Configuration block for thresholds.
    * `warning` - (Optional) Configuration block for warning severity threshold.
      * `static` - (Required) Configuration block for static threshold.
        * `value` - (Required) The threshold value.
    * `critical` - (Optional) Configuration block for critical severity threshold.
      * `static` - (Required) Configuration block for static threshold.
        * `value` - (Required) The threshold value.
* `time_threshold` - (Required) Configuration block for time threshold.
  * `violations_in_sequence` - (Required) Configuration block for violations in sequence.
    * `time_window` - (Optional) The time window in milliseconds.
* `custom_payload_field` - (Optional) Configuration block for custom payload fields in alert notifications.
  * `key` - (Required) The key of the custom payload field.
  * `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined.
  * `dynamic_value` - Optional - The dynamic value of the custom payload field.
  Either `value` or `dynamic_value` must be defined [Details](#dynamic-custom-payload-field-value).

#### Dynamic Custom Payload Field Value
* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload


### Tag Filter Argument Reference
The **tag_filter** defines which entities should be included into the application. It supports:

* logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
* comparison operators EQUALS, NOT_EQUAL, CONTAINS | NOT_CONTAIN, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH, NOT_ENDS_WITH, GREATER_OR_EQUAL_THAN, LESS_OR_EQUAL_THAN, LESS_THAN, GREATER_THAN
* unary operators IS_EMPTY, NOT_EMPTY, IS_BLANK, NOT_BLANK.

The **tag_filter** is defined by the following eBNF:

```plain
tag_filter                := logical_or
logical_or                := logical_and OR logical_or | logical_and
logical_and               := primary_expression AND logical_and | bracket_expression
bracket_expression        := ( logical_or ) | primary_expression
primary_expression        := comparison | unary_operator_expression
comparison                := identifier comparison_operator value | identifier@entity_origin comparison_operator value | identifier:tag_key comparison_operator value | identifier:tag_key@entity_origin comparison_operator value
comparison_operator       := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN | STARTS_WITH | ENDS_WITH | NOT_STARTS_WITH | NOT_ENDS_WITH | GREATER_OR_EQUAL_THAN | LESS_OR_EQUAL_THAN | LESS_THAN | GREATER_THAN
unary_operator_expression := identifier unary_operator | identifier@entity_origin unary_operator
unary_operator            := IS_EMPTY | NOT_EMPTY | IS_BLANK | NOT_BLANK
tag_key                   := identifier | string_value
entity_origin             := src | dest | na
value                     := string_value | number_value | boolean_value
string_value              := "'" <string> "'"
number_value              := (+-)?[0-9]+
boolean_value             := TRUE | FALSE
identifier                := [a-zA-Z_][\.a-zA-Z0-9_\-/]*
```

## Attribute Reference

* `id` - The ID of the log alert configuration.

## Import

Log alert configurations can be imported using the ID, e.g.,

```
$ terraform import instana_log_alert_config.example 12345678-1234-1234-1234-123456789012
```
