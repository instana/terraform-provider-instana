# Log Alert Configuration Resource

This resource manages log alert configurations in Instana. Log alerts monitor log data and trigger notifications based on log counts and patterns.

API Documentation: <https://instana.github.io/openapi/#tag/Log-Alert-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v6, all nested configurations use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_log_alert_config" "example" {
  name = "Log Alert"
  
  alert_channels {
    warning = ["channel-1"]
    critical = ["channel-2"]
  }
  
  group_by {
    tag_name = "kubernetes.namespace.name"
  }
  
  rules {
    metric_name = "log.count"
    threshold {
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
  
  custom_payload_field {
    key = "env"
    value = "prod"
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_log_alert_config" "example" {
  name = "Log Alert"
  
  alert_channels = {
    warning = ["channel-1"]
    critical = ["channel-2"]
  }
  
  group_by = [
    {
      tag_name = "kubernetes.namespace.name"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    threshold = {
      critical = {
        static = {
          value = 500
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_field = [
    {
      key = "env"
      value = "prod"
    }
  ]
}
```

### Key Syntax Changes

1. **Alert Channels**: `alert_channels { }` → `alert_channels = { }`
2. **Group By**: Multiple `group_by { }` blocks → Single `group_by = [{ }]` list
3. **Rules**: `rules { }` → `rules = { }`
4. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
5. **Custom Payload Fields**: Multiple `custom_payload_field { }` blocks → Single `custom_payload_field = [{ }]` list
6. **All nested objects**: Use `= { }` syntax

## Example Usage

### Basic Log Count Alert

```hcl
resource "instana_log_alert_config" "error_logs" {
  name = "High Error Log Count"
  description = "Alert when error logs exceed threshold"
  tag_filter = "log.level@na EQUALS 'ERROR'"
  granularity = 600000
  
  alert_channels = {
    critical = ["ops-team-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 100
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### Log Alert with Warning and Critical Thresholds

```hcl
resource "instana_log_alert_config" "application_errors" {
  name = "Application Error Logs"
  description = "Monitor application error log volume"
  tag_filter = "log.message@na CONTAINS 'error' AND application.name@na EQUALS 'my-app'"
  granularity = 300000
  
  alert_channels = {
    warning = ["dev-team-channel"]
    critical = ["ops-team-channel", "pagerduty-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      warning = {
        static = {
          value = 50
        }
      }
      critical = {
        static = {
          value = 200
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
}
```

### Log Alert with Group By

```hcl
resource "instana_log_alert_config" "grouped_by_namespace" {
  name = "Errors by Namespace"
  description = "Monitor error logs grouped by Kubernetes namespace"
  tag_filter = "log.level@na EQUALS 'ERROR'"
  granularity = 600000
  
  alert_channels = {
    critical = ["k8s-ops-channel"]
  }
  
  group_by = [
    {
      tag_name = "kubernetes.namespace.name"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 100
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### Multiple Group By Tags

```hcl
resource "instana_log_alert_config" "multi_group" {
  name = "Errors by Service and Environment"
  description = "Group errors by service and environment"
  tag_filter = "log.level@na EQUALS 'ERROR'"
  granularity = 600000
  
  alert_channels = {
    critical = ["ops-channel"]
  }
  
  group_by = [
    {
      tag_name = "service.name"
    },
    {
      tag_name = "environment"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 50
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### Log Alert with Custom Payload

```hcl
resource "instana_log_alert_config" "with_context" {
  name = "Errors with Context"
  description = "Error logs with additional context"
  tag_filter = "log.level@na EQUALS 'ERROR'"
  granularity = 600000
  
  alert_channels = {
    critical = ["enriched-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 100
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_field = [
    {
      key = "environment"
      value = "production"
    },
    {
      key = "team"
      value = "platform-ops"
    },
    {
      key = "runbook"
      value = "https://wiki.example.com/runbooks/error-logs"
    }
  ]
}
```

### Log Alert with Dynamic Payload

```hcl
resource "instana_log_alert_config" "dynamic_context" {
  name = "Errors with Dynamic Tags"
  description = "Include dynamic tag values in alerts"
  tag_filter = "log.level@na EQUALS 'ERROR'"
  granularity = 600000
  
  alert_channels = {
    critical = ["ops-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 100
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_field = [
    {
      key = "static_env"
      value = "production"
    },
    {
      key = "host_fqdn"
      dynamic_value = {
        tag_name = "host.fqdn"
      }
    },
    {
      key = "service_name"
      dynamic_value = {
        key = "name"
        tag_name = "service.name"
      }
    }
  ]
}
```

### Complex Tag Filter

```hcl
resource "instana_log_alert_config" "complex_filter" {
  name = "Production Critical Errors"
  description = "Monitor critical errors in production"
  granularity = 300000
  
  # Complex tag filter with multiple conditions
  tag_filter = "(log.level@na EQUALS 'ERROR' OR log.level@na EQUALS 'FATAL') AND environment@na EQUALS 'production' AND service.tier@na EQUALS 'critical'"
  
  alert_channels = {
    critical = ["critical-ops-channel", "pagerduty-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 10
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
}
```

### Kubernetes Pod Errors

```hcl
resource "instana_log_alert_config" "k8s_pod_errors" {
  name = "Kubernetes Pod Errors"
  description = "Monitor pod error logs"
  tag_filter = "log.level@na EQUALS 'ERROR' AND kubernetes.pod.name@na NOT_EMPTY"
  granularity = 600000
  
  alert_channels = {
    warning = ["k8s-team-channel"]
    critical = ["k8s-ops-channel"]
  }
  
  group_by = [
    {
      tag_name = "kubernetes.namespace.name"
    },
    {
      tag_name = "kubernetes.pod.name"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      warning = {
        static = {
          value = 20
        }
      }
      critical = {
        static = {
          value = 100
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_field = [
    {
      key = "cluster"
      dynamic_value = {
        tag_name = "kubernetes.cluster.name"
      }
    }
  ]
}
```

### Application-Specific Log Monitoring

```hcl
resource "instana_log_alert_config" "app_exceptions" {
  name = "Application Exceptions"
  description = "Monitor application exception logs"
  tag_filter = "log.message@na CONTAINS 'Exception' AND application.name@na EQUALS 'payment-service'"
  granularity = 300000
  grace_period = 600000
  
  alert_channels = {
    critical = ["payment-team-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 5
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  
  custom_payload_field = [
    {
      key = "application"
      value = "payment-service"
    },
    {
      key = "severity"
      value = "high"
    }
  ]
}
```

### Security Log Monitoring

```hcl
resource "instana_log_alert_config" "security_events" {
  name = "Security Events"
  description = "Monitor security-related log events"
  tag_filter = "log.message@na CONTAINS 'authentication failed' OR log.message@na CONTAINS 'unauthorized access'"
  granularity = 300000
  
  alert_channels = {
    critical = ["security-team-channel", "soc-channel"]
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = 3
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  
  custom_payload_field = [
    {
      key = "alert_type"
      value = "security"
    },
    {
      key = "priority"
      value = "immediate"
    }
  ]
}
```

### Multi-Environment Setup

```hcl
locals {
  environments = {
    production = {
      threshold = 50
      granularity = 300000
      channels = ["prod-ops", "pagerduty"]
    }
    staging = {
      threshold = 200
      granularity = 600000
      channels = ["staging-ops"]
    }
    development = {
      threshold = 500
      granularity = 900000
      channels = ["dev-team"]
    }
  }
}

resource "instana_log_alert_config" "env_errors" {
  for_each = local.environments

  name = "${each.key} Error Logs"
  description = "Error log monitoring for ${each.key}"
  tag_filter = "log.level@na EQUALS 'ERROR' AND environment@na EQUALS '${each.key}'"
  granularity = each.value.granularity
  
  alert_channels = {
    critical = each.value.channels
  }
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      critical = {
        static = {
          value = each.value.threshold
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = each.value.granularity
    }
  }
  
  custom_payload_field = [
    {
      key = "environment"
      value = each.key
    }
  ]
}
```

### Database Error Logs

```hcl
resource "instana_log_alert_config" "db_errors" {
  name = "Database Error Logs"
  description = "Monitor database connection and query errors"
  tag_filter = "log.message@na CONTAINS 'database' AND (log.message@na CONTAINS 'error' OR log.message@na CONTAINS 'timeout')"
  granularity = 600000
  
  alert_channels = {
    warning = ["dba-team-channel"]
    critical = ["dba-team-channel", "ops-channel"]
  }
  
  group_by = [
    {
      tag_name = "database.name"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      warning = {
        static = {
          value = 10
        }
      }
      critical = {
        static = {
          value = 50
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### API Rate Limit Logs

```hcl
resource "instana_log_alert_config" "rate_limit" {
  name = "API Rate Limit Exceeded"
  description = "Alert on rate limit violations"
  tag_filter = "log.message@na CONTAINS 'rate limit exceeded'"
  granularity = 300000
  
  alert_channels = {
    warning = ["api-team-channel"]
  }
  
  group_by = [
    {
      tag_name = "api.endpoint"
    },
    {
      tag_name = "client.id"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    alert_type = "log.count"
    aggregation = "SUM"
    threshold_operator = ">"
    
    threshold = {
      warning = {
        static = {
          value = 100
        }
      }
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  
  custom_payload_field = [
    {
      key = "action"
      value = "review_rate_limits"
    }
  ]
}
```

## Argument Reference

* `name` - Required - The name of the log alert configuration (max 256 characters)
* `description` - Required - A description of the log alert configuration (max 65536 characters)
* `tag_filter` - Required - The tag filter expression used to scope the alert [Details](#tag-filter-argument-reference)
* `granularity` - Optional - The evaluation granularity in milliseconds. Default: `600000` (10 minutes). Possible values: `60000` (1 minute), `300000` (5 minutes), `600000` (10 minutes), `900000` (15 minutes), `1200000` (20 minutes), `1800000` (30 minutes)
* `grace_period` - Optional - The duration in milliseconds for which an alert remains open after conditions are no longer violated. The alert auto-closes once the grace period expires
* `alert_channels` - Optional - Configuration for alert notification channels [Details](#alert-channels-reference)
* `group_by` - Optional - List of tags to group results by [Details](#group-by-reference)
* `rules` - Required - Configuration for alert rules [Details](#rules-reference)
* `time_threshold` - Required - Configuration for time threshold [Details](#time-threshold-reference)
* `custom_payload_field` - Optional - List of custom payload fields in alert notifications [Details](#custom-payload-field-reference)

### Alert Channels Reference

* `warning` - Optional - List of alert channel IDs to notify for warning severity alerts
* `critical` - Optional - List of alert channel IDs to notify for critical severity alerts

### Group By Reference

* `tag_name` - Required - The tag name to group by
* `key` - Optional - The key to group by (for nested tags)

### Rules Reference

* `metric_name` - Required - The metric name to monitor (typically `log.count`)
* `alert_type` - Optional - The type of the log alert rule. Only `log.count` is supported. Default: `log.count`
* `aggregation` - Optional - The aggregation method to use for the log alert. Only `SUM` is supported. Default: `SUM`
* `threshold_operator` - Required - The operator to use for threshold comparison. Possible values: `>`, `>=`, `<`, `<=`
* `threshold` - Required - Configuration for thresholds [Details](#threshold-reference)

#### Threshold Reference

At least one of the following must be configured:

* `warning` - Optional - Configuration for warning severity threshold [Details](#threshold-severity-reference)
* `critical` - Optional - Configuration for critical severity threshold [Details](#threshold-severity-reference)

##### Threshold Severity Reference

* `static` - Required - Configuration for static threshold
  * `value` - Required - The threshold value (integer)

### Time Threshold Reference

* `violations_in_sequence` - Required - Configuration for violations in sequence
  * `time_window` - Required - The time window in milliseconds

### Custom Payload Field Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - The dynamic value of the custom payload field. Either `value` or `dynamic_value` must be defined [Details](#dynamic-custom-payload-field-value)

#### Dynamic Custom Payload Field Value

* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload

### Tag Filter Argument Reference

The **tag_filter** defines which log entries should be included in the alert scope. It supports:

* Logical AND and/or logical OR conjunctions (AND has higher precedence than OR)
* Comparison operators: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `LESS_THAN`, `GREATER_THAN`
* Unary operators: `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`

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

## Attributes Reference

* `id` - The ID of the log alert configuration

## Import

Log alert configurations can be imported using the ID, e.g.:

```bash
$ terraform import instana_log_alert_config.example 12345678-1234-1234-1234-123456789012
```

## Notes

* Log alerts are evaluated at the specified `granularity` interval
* The `grace_period` prevents alert flapping by keeping alerts open for a specified duration after conditions are no longer met
* Use `group_by` to create separate alerts for different tag values (e.g., per namespace, per service)
* Tag filters support complex expressions for precise log selection
* Custom payload fields can include both static values and dynamic tag values from the logs
* Only `log.count` metric and `SUM` aggregation are currently supported
* The `time_threshold` defines how many consecutive violations are required before triggering an alert
