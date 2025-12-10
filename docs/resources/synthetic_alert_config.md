# Synthetic Alert Configuration Resource

Manages a synthetic alert configuration in Instana. Synthetic alerts monitor the health and performance of synthetic tests, triggering notifications when tests fail or performance degrades.

API Documentation: <https://instana.github.io/openapi/#tag/Synthetic-Alert-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v6, they use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_synthetic_alert_config" "example" {
  name = "Synthetic Test Alert"
  
  rule {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "SUM"
  }
  
  time_threshold {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  custom_payload_field {
    key = "environment"
    value = "production"
  }
  custom_payload_field {
    key = "team"
    value = "platform"
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_synthetic_alert_config" "example" {
  alert_channel_ids = ["var.alert_channel_id"]
  custom_payload_field = [
    {
      key           = "e2e_test"
      value         = "40"
    },
  ]
  description  = "test example"
  grace_period = 604800000
  name         = "test example"
  rule = {
    aggregation = "SUM"
    alert_type  = "failure"
    metric_name = "status"
  }
  severity           = 5
  synthetic_test_ids = ["var.test_id"]
  tag_filter         = "synthetic.locationLabel@na EQUALS 'testlocation' "
  time_threshold = {
    type             = "violationsInSequence"
    violations_count = 1
  }
}
```

### Key Syntax Changes

1. **Rule**: `rule { }` → `rule = { }`
2. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
3. **Custom Payload Fields**: `custom_payload_field { }` (multiple) → `custom_payload_fields = [{ }]` (list)
4. **Aggregation**: Case-insensitive but lowercase recommended: `SUM` → `sum`

## Example Usage

### Basic Synthetic Test Failure Alert

Monitor synthetic tests and alert on failures:

```hcl
resource "instana_synthetic_alert_config" "basic_failure" {
  name = "Synthetic Test Failure Alert"
  description = "Alert when synthetic tests fail"
  
  synthetic_test_ids = ["test-id-1", "test-id-2"]
  severity = 5  # Critical
  grace_period = 604800000

  rule = {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "sum"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  tag_filter         = "synthetic.locationLabel@na EQUALS 'testlocation' "

  alert_channel_ids = ["channel-id-1"]
  
   custom_payload_field = [
    {
      key           = "e2e_test"
      value         = "40"
    },
    {
      key = "test_location"
      dynamic_value = {
        key = "location"
        tag_name = "synthetic.locationLabel"
      }
    }
  ]
}
```

## Argument Reference

* `name` - Required - Name of the synthetic alert configuration
* `description` - Required - Description of the synthetic alert configuration
* `synthetic_test_ids` - Optional - Set of synthetic test IDs to monitor. If not specified, all tests matching the tag filter will be monitored
* `severity` - Optional - Severity of the alert. Must be either `5` (critical) or `10` (warning). Default: `5`
* `tag_filter` - Required - Tag filter expression to limit which synthetic tests are monitored [Details](#tag-filter-reference)
* `rule` - Required - Rule configuration for the alert [Details](#rule-reference)
* `alert_channel_ids` - Required - Set of alert channel IDs to notify when the alert is triggered
* `time_threshold` - Required - Time threshold configuration [Details](#time-threshold-reference)
* `grace_period` - Optional - Duration in milliseconds for which an alert remains open after conditions are no longer violated. The alert auto-closes once the grace period expires
* `custom_payload_fields` - Optional - List of custom payload fields to include in alert notifications [Details](#custom-payload-fields-reference)

### Rule Reference

* `alert_type` - Required - Type of the alert rule. Currently only `failure` is supported
* `metric_name` - Required - Metric name to monitor (e.g., `status`, `responseTime`)
* `aggregation` - Optional - Aggregation method. Allowed values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distinct_count`, `sum_positive`, `per_second`, `increase`

### Time Threshold Reference

* `type` - Required - Type of the time threshold. Must be either `violationsInSequence`
* `violations_count` - Optional - Number of violations required to trigger the alert

### Custom Payload Fields Reference

Custom payload fields allow you to add additional context to alert notifications. Each field can have either a static value or a dynamic value derived from tags.

* `key` - Required - Key of the custom payload field
* `value` - Optional - Static string value of the custom payload field. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - Dynamic value derived from tags [Details](#dynamic-value-reference)

#### Dynamic Value Reference

* `key` - Optional - Key to use in the payload for the tag value
* `tag_name` - Required - Name of the tag whose value should be added to the payload

### Tag Filter Reference

The **tag_filter** defines which synthetic tests should be monitored. It supports:

* Logical AND and/or logical OR conjunctions (AND has higher precedence than OR)
* Comparison operators: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `LESS_THAN`, `GREATER_THAN`
* Unary operators: `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`

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

**Common Tag Filter Examples:**
* `synthetic.locationLabel@na EQUALS 'us-east-1'` - Tests in US East region
* `synthetic.testType@na EQUALS 'api'` - API tests only
* `synthetic.environment@na EQUALS 'production'` - Production tests
* `(synthetic.locationLabel@na EQUALS 'us-east-1' OR synthetic.locationLabel@na EQUALS 'us-west-2') AND synthetic.testType@na EQUALS 'api'` - API tests in US regions

## Attributes Reference

* `id` - The ID of the synthetic alert configuration

## Import

Synthetic alert configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_synthetic_alert_config.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Severity `5` represents critical alerts, `10` represents warning alerts
* Use `synthetic_test_ids` to monitor specific tests or `tag_filter` to monitor tests matching criteria
* The `grace_period` prevents alert flapping by keeping alerts open for a specified duration
* Custom payload fields enhance alert notifications with additional context
* Dynamic payload fields pull values from synthetic test tags at runtime
* Tag filters use the same syntax as application and infrastructure monitoring
* The `alert_type` currently only supports `failure` for synthetic tests
* Multiple alert channels can be notified simultaneously
* Time thresholds help reduce noise by requiring multiple violations before alerting
