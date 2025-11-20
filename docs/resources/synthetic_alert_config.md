# Synthetic Alert Configuration Resource

Manages a synthetic alert configuration in Instana. Synthetic alerts monitor the health and performance of synthetic tests, triggering notifications when tests fail or performance degrades.

API Documentation: <https://instana.github.io/openapi/#tag/Synthetic-Alert-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.


## Migration Guide (v5 to v6)

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
  name = "Synthetic Test Alert"
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "sum"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  custom_payload_fields = [
    {
      key = "environment"
      value = "production"
    },
    {
      key = "team"
      value = "platform"
    }
  ]
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
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "sum"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  alert_channel_ids = ["channel-id-1"]
}
```

### Alert with Tag Filter

Use tag filters to monitor specific synthetic tests:

```hcl
resource "instana_synthetic_alert_config" "location_based" {
  name = "US Region Synthetic Alerts"
  description = "Monitor synthetic tests in US locations"
  
  tag_filter = "synthetic.locationLabel@na EQUALS 'us-east-1'"
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "sum"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 1
  }
  
  alert_channel_ids = ["us-team-channel"]
}
```

### Alert with Grace Period

Add a grace period before auto-closing alerts:

```hcl
resource "instana_synthetic_alert_config" "with_grace_period" {
  name = "Synthetic Alert with Grace Period"
  description = "Alert with 5-minute grace period"
  
  synthetic_test_ids = ["critical-test-1"]
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  grace_period = 300000  # 5 minutes in milliseconds
  alert_channel_ids = ["ops-channel"]
}
```

### Alert with Static Custom Payload

Add static custom fields to alert notifications:

```hcl
resource "instana_synthetic_alert_config" "with_static_payload" {
  name = "Synthetic Alert with Custom Payload"
  description = "Alert with custom static fields"
  
  synthetic_test_ids = ["api-test-1"]
  severity = 10  # Warning
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 1
  }
  
  custom_payload_fields = [
    {
      key = "environment"
      value = "production"
    },
    {
      key = "team"
      value = "platform"
    },
    {
      key = "runbook"
      value = "https://wiki.example.com/runbooks/synthetic-failures"
    }
  ]
  
  alert_channel_ids = ["platform-team-channel"]
}
```

### Alert with Dynamic Custom Payload

Add dynamic fields from tags to alert notifications:

```hcl
resource "instana_synthetic_alert_config" "with_dynamic_payload" {
  name = "Synthetic Alert with Dynamic Payload"
  description = "Alert with dynamic tag-based fields"
  
  synthetic_test_ids = ["test-1", "test-2"]
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  custom_payload_fields = [
    {
      key = "test_location"
      dynamic_value = {
        key = "location"
        tag_name = "synthetic.locationLabel"
      }
    },
    {
      key = "test_type"
      dynamic_value = {
        tag_name = "synthetic.testType"
      }
    }
  ]
  
  alert_channel_ids = ["monitoring-channel"]
}
```

### Alert with Mixed Custom Payload

Combine static and dynamic custom fields:

```hcl
resource "instana_synthetic_alert_config" "mixed_payload" {
  name = "Synthetic Alert with Mixed Payload"
  description = "Alert with both static and dynamic fields"
  
  synthetic_test_ids = ["critical-test"]
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 1
  }
  
  custom_payload_fields = [
    {
      key = "severity_level"
      value = "critical"
    },
    {
      key = "alert_source"
      value = "synthetic_monitoring"
    },
    {
      key = "test_location"
      dynamic_value = {
        tag_name = "synthetic.locationLabel"
      }
    },
    {
      key = "test_name"
      dynamic_value = {
        key = "name"
        tag_name = "synthetic.testName"
      }
    }
  ]
  
  alert_channel_ids = ["ops-channel", "slack-channel"]
}
```

### Violations in Period Alert

Alert based on violations within a time period:

```hcl
resource "instana_synthetic_alert_config" "violations_in_period" {
  name = "Multiple Failures in Period"
  description = "Alert when multiple failures occur in a time window"
  
  synthetic_test_ids = ["test-1", "test-2", "test-3"]
  severity = 10
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInPeriod"
    violations_count = 3
  }
  
  alert_channel_ids = ["team-channel"]
}
```

### Multi-Location Monitoring

Monitor synthetic tests across multiple locations:

```hcl
resource "instana_synthetic_alert_config" "multi_location" {
  name = "Multi-Location Synthetic Monitoring"
  description = "Monitor tests across all regions"
  
  tag_filter = "synthetic.locationLabel@na CONTAINS 'region'"
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "sum"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  custom_payload_fields = [
    {
      key = "alert_type"
      value = "multi_location_failure"
    },
    {
      key = "location"
      dynamic_value = {
        tag_name = "synthetic.locationLabel"
      }
    }
  ]
  
  alert_channel_ids = ["global-ops-channel"]
}
```

### API Endpoint Monitoring

Monitor specific API endpoints via synthetic tests:

```hcl
resource "instana_synthetic_alert_config" "api_monitoring" {
  name = "API Endpoint Monitoring"
  description = "Monitor critical API endpoints"
  
  tag_filter = "synthetic.testType@na EQUALS 'api'"
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 1
  }
  
  grace_period = 60000  # 1 minute
  
  custom_payload_fields = [
    {
      key = "endpoint_type"
      value = "api"
    },
    {
      key = "priority"
      value = "high"
    }
  ]
  
  alert_channel_ids = ["api-team-channel", "pagerduty-channel"]
}
```

### Browser Test Monitoring

Monitor browser-based synthetic tests:

```hcl
resource "instana_synthetic_alert_config" "browser_monitoring" {
  name = "Browser Test Monitoring"
  description = "Monitor browser synthetic tests"
  
  tag_filter = "synthetic.testType@na EQUALS 'browser'"
  severity = 10
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  custom_payload_fields = [
    {
      key = "test_type"
      value = "browser"
    },
    {
      key = "browser"
      dynamic_value = {
        tag_name = "synthetic.browser"
      }
    }
  ]
  
  alert_channel_ids = ["frontend-team-channel"]
}
```

### Critical Service Monitoring

High-priority monitoring for critical services:

```hcl
resource "instana_synthetic_alert_config" "critical_service" {
  name = "Critical Service Monitoring"
  description = "Immediate alerts for critical service failures"
  
  synthetic_test_ids = ["payment-api-test", "auth-api-test"]
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 1  # Alert immediately
  }
  
  grace_period = 0  # No grace period
  
  custom_payload_fields = [
    {
      key = "priority"
      value = "P1"
    },
    {
      key = "service_tier"
      value = "critical"
    },
    {
      key = "escalation_policy"
      value = "immediate"
    }
  ]
  
  alert_channel_ids = ["pagerduty-critical", "slack-critical", "email-oncall"]
}
```

### Environment-Specific Monitoring

Monitor tests in specific environments:

```hcl
locals {
  environments = {
    production = {
      tag_filter = "synthetic.environment@na EQUALS 'production'"
      severity = 5
      violations = 1
      channels = ["prod-ops-channel", "pagerduty"]
    }
    staging = {
      tag_filter = "synthetic.environment@na EQUALS 'staging'"
      severity = 10
      violations = 2
      channels = ["staging-team-channel"]
    }
  }
}

resource "instana_synthetic_alert_config" "env_monitoring" {
  for_each = local.environments

  name = "${title(each.key)} Synthetic Monitoring"
  description = "Monitor synthetic tests in ${each.key}"
  
  tag_filter = each.value.tag_filter
  severity = each.value.severity
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = each.value.violations
  }
  
  custom_payload_fields = [
    {
      key = "environment"
      value = each.key
    }
  ]
  
  alert_channel_ids = each.value.channels
}
```

### Performance Degradation Alert

Monitor for performance issues:

```hcl
resource "instana_synthetic_alert_config" "performance_degradation" {
  name = "Synthetic Performance Degradation"
  description = "Alert on slow response times"
  
  synthetic_test_ids = ["perf-test-1", "perf-test-2"]
  severity = 10
  
  rule = {
    alert_type = "failure"
    metric_name = "responseTime"
    aggregation = "mean"
  }
  
  time_threshold = {
    type = "violationsInPeriod"
    violations_count = 3
  }
  
  custom_payload_fields = [
    {
      key = "alert_category"
      value = "performance"
    },
    {
      key = "response_time"
      dynamic_value = {
        tag_name = "synthetic.responseTime"
      }
    }
  ]
  
  alert_channel_ids = ["performance-team-channel"]
}
```

### Complex Tag Filter Alert

Use complex tag filters for precise monitoring:

```hcl
resource "instana_synthetic_alert_config" "complex_filter" {
  name = "Complex Filter Synthetic Alert"
  description = "Monitor with complex tag filter conditions"
  
  tag_filter = "(synthetic.locationLabel@na EQUALS 'us-east-1' OR synthetic.locationLabel@na EQUALS 'us-west-2') AND synthetic.testType@na EQUALS 'api'"
  severity = 5
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 2
  }
  
  custom_payload_fields = [
    {
      key = "region"
      dynamic_value = {
        tag_name = "synthetic.locationLabel"
      }
    },
    {
      key = "test_type"
      dynamic_value = {
        tag_name = "synthetic.testType"
      }
    }
  ]
  
  alert_channel_ids = ["us-ops-channel"]
}
```

### Aggregation-Based Alert

Use different aggregation methods:

```hcl
resource "instana_synthetic_alert_config" "aggregation_alert" {
  name = "Aggregated Synthetic Alert"
  description = "Alert using max aggregation"
  
  synthetic_test_ids = ["test-1"]
  severity = 10
  
  rule = {
    alert_type = "failure"
    metric_name = "status"
    aggregation = "max"
  }
  
  time_threshold = {
    type = "violationsInSequence"
    violations_count = 1
  }
  
  alert_channel_ids = ["monitoring-channel"]
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

* `type` - Required - Type of the time threshold. Must be either `violationsInSequence` or `violationsInPeriod`
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
