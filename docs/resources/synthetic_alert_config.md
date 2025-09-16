# Synthetic Alert Configuration Resource

This resource manages a synthetic alert configuration in Instana.

API Documentation: <https://instana.github.io/openapi/#tag/Synthetic-Alert-Configuration>

## Example Usage

```hcl
resource  "instana_synthetic_alert_config"  "example" {
  name        = "Synthetic Test Failure Alert"
  description = "Alert when synthetic tests fail"
  
  # Specific synthetic tests to monitor 
  synthetic_test_ids = ["test-id-1", "test-id-2"]
  
  # Alert severity: 5 (critical) or 10 (warning)
  severity = 5
  
  # tag filter to limit which synthetic tests are monitored
  tag_filter = "synthetic.locationLabel@na EQUALS 'location'"
  
  # Rule configuration
  rule {
    alert_type  = "failure"
    metric_name = "status"
    aggregation = "SUM"
  }
  
  # Alert channels to notify
  alert_channel_ids = ["alert-channel-id-1", "alert-channel-id-2"]
  
  # Time threshold configuration
  time_threshold {
    type            = "violationsInSequence"
    violations_count = 2
  }
  
  # Grace period in milliseconds (optional)
  grace_period = 300000
  
  # Custom payload fields for alert notifications
  custom_payload_field {
    key   = "key"
    value = "value"
  }
  custom_payload_field {
    key = "test2"
    dynamic_value {
      key = "dynamic-value-key"
      tag_name = "dynamic-value-tag-name"
    }
  }
}
```

## Argument Reference

* `name` - (Required) The name of the synthetic alert configuration.
* `description` - (Required) A description of the synthetic alert configuration.
* `synthetic_test_ids` - (Required) A list of synthetic test IDs to monitor. If not specified, all synthetic tests matching the tag filter will be monitored.
* `severity` - (Optional) The severity of the alert. Must be either `5` (critical) or `10` (warning).
* `tag_filter` - (Required) A tag filter expression to limit which synthetic tests are monitored.[Details](#tag-filter-argument-reference)
* `rule` - (Required) The rule configuration block. Only one rule block is allowed.
  * `alert_type` - (Required) The type of the alert rule. Currently only `failure` is supported.
  * `metric_name` - (Required) The metric name to monitor (e.g., `status`).
  * `aggregation` - (Optional) The aggregation method. Allowed avalues are `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTINCT_COUNT`, `SUM_POSITIVE`, `PER_SECOND`, `INCREASE`.
* `alert_channel_ids` - (Required) A list of alert channel IDs to notify when the alert is triggered.
* `time_threshold` - (Required) The time threshold configuration block. Only one time threshold block is allowed.
  * `type` - (Required) The type of the time threshold. Must be either `violationsInSequence` or `violationsInPeriod`.
  * `violations_count` - (Optional) The number of violations required to trigger the alert.
* `grace_period` - (Optional) The duration in milliseconds for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.
* `custom_payload_field` - (Optional) Custom payload fields to include in alert notifications. Multiple blocks can be specified.
[Details](#custom-payload-field-argument-reference)


### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined.
* `dynamic_value` - Optional - The dynamic value of the custom payload field [Details](#dynamic-custom-payload-field-value).

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
logical_and               := primary_expression AND logical_and | primary_expression
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

* `id` - The ID of the synthetic alert configuration.

