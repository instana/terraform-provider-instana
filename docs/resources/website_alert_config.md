# Website Alert Configuration Resource

> **⚠️ BREAKING CHANGES in v3.0.0**
> 
> This resource has been migrated from Terraform SDK v2 to the Plugin Framework. While most configurations remain compatible, there are important syntax changes you need to be aware of.
>
> **Key Changes:**
> - `rule` block now uses attribute syntax `= { }` instead of block syntax
> - `threshold` block now uses attribute syntax `= { }` instead of block syntax
> - `time_threshold` block now uses attribute syntax `= { }` instead of block syntax
> - `custom_payload_field` blocks now use list syntax `custom_payload_fields = [{ }]` instead of multiple blocks
> - `rules` (multiple alert rules) now use list syntax `= [{ }]`
> - All nested rule types use attribute syntax
> - See [Migration Guide](#migration-guide-v2-to-v3) below for detailed examples

Manages website alert configurations (Website Smart Alerts) in Instana. Website alerts monitor the performance and availability of your websites from the end-user perspective.

API Documentation: <https://instana.github.io/openapi/#operation/findActiveWebsiteAlertConfigs>

## Migration Guide (v2 to v3)

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v3, they use attribute syntax instead of block syntax.

#### OLD (v2.x) Syntax:
```hcl
resource "instana_website_alert_config" "example" {
  name = "Website Alert"
  
  rule {
    slowness {
      metric_name = "onLoadTime"
      aggregation = "P90"
    }
  }
  
  threshold {
    static {
      operator = ">="
      value = 5.0
    }
  }
  
  time_threshold {
    violations_in_sequence {
      time_window = 600000
    }
  }
  
  custom_payload_field {
    key = "environment"
    value = "production"
  }
  custom_payload_field {
    key = "team"
    value = "frontend"
  }
}
```

#### NEW (v3.x) Syntax:
```hcl
resource "instana_website_alert_config" "example" {
  name = "Website Alert"
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value = 5
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_fields = [
    {
      key = "environment"
      value = "production"
    },
    {
      key = "team"
      value = "frontend"
    }
  ]
}
```

### Key Syntax Changes

1. **Rule**: `rule { }` → `rule = { }`
2. **Threshold**: `threshold { }` → `threshold = { }`
3. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
4. **Custom Payload Fields**: `custom_payload_field { }` (multiple) → `custom_payload_fields = [{ }]` (list)
5. **Nested Rule Types**: `slowness { }` → `slowness = { }`
6. **Aggregation**: Case-insensitive but lowercase recommended: `P90` → `p90`

## Example Usage

### Basic Slowness Alert

Monitor page load time:

```hcl
resource "instana_website_alert_config" "slowness_basic" {
  name        = "Page Load Time Alert"
  description = "Alert when page load time exceeds threshold"
  severity    = "warning"
  triggering  = false
  website_id  = instana_website_monitoring_config.example.id
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
  granularity       = 600000
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 3000
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### Slowness Alert with Tag Filter

Monitor specific pages or user segments:

```hcl
resource "instana_website_alert_config" "slowness_filtered" {
  name        = "Checkout Page Slowness"
  description = "Monitor checkout page performance"
  severity    = "critical"
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "beacon.page.name@na EQUALS '/checkout'"
  
  rule = {
    slowness = {
      metric_name = "domContentLoadedTime"
      aggregation = "p95"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 2000
    }
  }
  
  time_threshold = {
    violations_in_period = {
      time_window = 600000
      violations  = 3
    }
  }
  
  alert_channel_ids = [
    instana_alerting_channel_pagerduty.oncall.id,
    instana_alerting_channel_slack.frontend.id
  ]
}
```

### JavaScript Error Alert

Monitor specific JavaScript errors:

```hcl
resource "instana_website_alert_config" "js_error" {
  name        = "Critical JS Error Alert"
  description = "Alert on specific JavaScript errors"
  severity    = "critical"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    specific_js_error = {
      metric_name = "jsErrors"
      aggregation = "sum"
      operator    = "CONTAINS"
      value       = "TypeError"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 10
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.dev_team.id]
}
```

### Status Code Alert

Monitor HTTP status codes:

```hcl
resource "instana_website_alert_config" "status_code" {
  name        = "4xx Error Alert"
  description = "Alert on client errors"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    status_code = {
      metric_name = "httpStatusCode"
      aggregation = "sum"
      operator    = "GREATER_OR_EQUAL_THAN"
      value       = "400"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 50
    }
  }
  
  time_threshold = {
    violations_in_period = {
      time_window = 600000
      violations  = 2
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.ops.id]
}
```

### Throughput Alert

Monitor page view throughput:

```hcl
resource "instana_website_alert_config" "throughput" {
  name        = "Low Traffic Alert"
  description = "Alert when traffic drops significantly"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    throughput = {
      metric_name = "pageViews"
      aggregation = "sum"
    }
  }
  
  threshold = {
    static = {
      operator = "<="
      value    = 100
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.ops.id]
}
```

### Alert with Adaptive Baseline

Use adaptive baseline for dynamic thresholds:

```hcl
resource "instana_website_alert_config" "adaptive_slowness" {
  name        = "Adaptive Page Load Alert"
  description = "Alert on abnormal page load times"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "mean"
    }
  }
  
  threshold = {
    adaptive_baseline = {
      operator         = ">="
      deviation_factor = 2.0
      adaptability     = 0.5
      seasonality      = "WEEKLY"
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
}
```

### Alert with Historic Baseline

Use historic baseline for comparison:

```hcl
resource "instana_website_alert_config" "historic_baseline" {
  name        = "Historic Baseline Alert"
  description = "Compare against historical performance"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    historic_baseline = {
      deviation   = 1.5
      seasonality = "DAILY"
      baseline = [
        {
          day_of_week = "MONDAY"
          start       = "09:00"
          end         = "17:00"
          baseline    = 2000
        },
        {
          day_of_week = "TUESDAY"
          start       = "09:00"
          end         = "17:00"
          baseline    = 1800
        }
      ]
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
}
```

### User Impact Alert

Alert based on user impact:

```hcl
resource "instana_website_alert_config" "user_impact" {
  name        = "High User Impact Alert"
  description = "Alert when many users are affected"
  severity    = "critical"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 5000
    }
  }
  
  time_threshold = {
    user_impact_of_violations_in_sequence = {
      time_window              = 600000
      impact_measurement_method = "AGGREGATED"
      user_percentage          = 0.1  # 10% of users
    }
  }
  
  alert_channel_ids = [
    instana_alerting_channel_pagerduty.oncall.id
  ]
}
```

### Alert with Custom Payload

Add custom fields to alert notifications:

```hcl
resource "instana_website_alert_config" "with_custom_payload" {
  name        = "Alert with Custom Payload"
  description = "Alert with additional context"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 3000
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_fields = [
    {
      key   = "environment"
      value = "production"
    },
    {
      key   = "team"
      value = "frontend"
    },
    {
      key   = "runbook"
      value = "https://wiki.example.com/runbooks/website-slowness"
    },
    {
      key = "user_segment"
      dynamic_value = {
        key      = "segment"
        tag_name = "beacon.user.segment"
      }
    }
  ]
  
  alert_channel_ids = [instana_alerting_channel_slack.frontend.id]
}
```

### Complex Tag Filter Alert

Use complex tag filters for precise monitoring:

```hcl
resource "instana_website_alert_config" "complex_filter" {
  name        = "Premium User Slowness"
  description = "Monitor premium users on critical pages"
  severity    = "critical"
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "(beacon.user.tier@na EQUALS 'premium' OR beacon.user.tier@na EQUALS 'enterprise') AND (beacon.page.name@na EQUALS '/dashboard' OR beacon.page.name@na EQUALS '/reports')"
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p95"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 2000
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  
  alert_channel_ids = [
    instana_alerting_channel_pagerduty.oncall.id,
    instana_alerting_channel_slack.vip_support.id
  ]
}
```

### Multiple Metrics Alert

Monitor multiple metrics with different thresholds:

```hcl
resource "instana_website_alert_config" "multi_metric" {
  name        = "Comprehensive Performance Alert"
  description = "Monitor multiple performance metrics"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  granularity = 300000
  
  rules = [
    {
      rule = {
        slowness = {
          metric_name = "onLoadTime"
          aggregation = "p90"
        }
      }
      threshold_operator = ">="
      thresholds = {
        warning = {
          static = {
            operator = ">="
            value    = 3000
          }
        }
        critical = {
          static = {
            operator = ">="
            value    = 5000
          }
        }
      }
    },
    {
      rule = {
        slowness = {
          metric_name = "domContentLoadedTime"
          aggregation = "p90"
        }
      }
      threshold_operator = ">="
      thresholds = {
        warning = {
          static = {
            operator = ">="
            value    = 2000
          }
        }
      }
    }
  ]
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
}
```

### Mobile Website Alert

Monitor mobile website performance:

```hcl
resource "instana_website_alert_config" "mobile_performance" {
  name        = "Mobile Performance Alert"
  description = "Monitor mobile user experience"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "beacon.device.type@na EQUALS 'mobile'"
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 4000  # Higher threshold for mobile
    }
  }
  
  time_threshold = {
    violations_in_period = {
      time_window = 600000
      violations  = 3
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.mobile_team.id]
}
```

### Geographic-Specific Alert

Monitor performance in specific regions:

```hcl
resource "instana_website_alert_config" "geo_specific" {
  name        = "APAC Region Performance"
  description = "Monitor performance for APAC users"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "beacon.geo.region@na EQUALS 'APAC'"
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 3500
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.apac_team.id]
}
```

### Browser-Specific Alert

Monitor specific browser performance:

```hcl
resource "instana_website_alert_config" "browser_specific" {
  name        = "IE11 Performance Alert"
  description = "Monitor legacy browser performance"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "beacon.browser.name@na EQUALS 'Internet Explorer' AND beacon.browser.version@na STARTS_WITH '11'"
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 6000  # Higher threshold for legacy browsers
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.legacy_support.id]
}
```

### Error Rate Alert

Monitor JavaScript error rate:

```hcl
resource "instana_website_alert_config" "error_rate" {
  name        = "High Error Rate Alert"
  description = "Alert on elevated error rates"
  severity    = "critical"
  website_id  = instana_website_monitoring_config.example.id
  
  rule = {
    specific_js_error = {
      metric_name = "jsErrors"
      aggregation = "sum"
      operator    = "IS_NOT_EMPTY"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 100
    }
  }
  
  time_threshold = {
    violations_in_period = {
      time_window = 300000
      violations  = 2
    }
  }
  
  alert_channel_ids = [
    instana_alerting_channel_pagerduty.oncall.id,
    instana_alerting_channel_slack.dev_team.id
  ]
}
```

### Environment-Specific Alerts

Create alerts for different environments:

```hcl
locals {
  environments = {
    production = {
      severity    = "critical"
      threshold   = 3000
      time_window = 300000
      channels    = [
        instana_alerting_channel_pagerduty.prod_oncall.id,
        instana_alerting_channel_slack.prod_alerts.id
      ]
    }
    staging = {
      severity    = "warning"
      threshold   = 5000
      time_window = 600000
      channels    = [
        instana_alerting_channel_email.staging_team.id
      ]
    }
  }
}

resource "instana_website_alert_config" "env_alerts" {
  for_each = local.environments

  name        = "${title(each.key)} Website Performance"
  description = "Monitor ${each.key} website performance"
  severity    = each.value.severity
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "beacon.environment@na EQUALS '${each.key}'"
  
  rule = {
    slowness = {
      metric_name = "onLoadTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = each.value.threshold
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = each.value.time_window
    }
  }
  
  custom_payload_fields = [
    {
      key   = "environment"
      value = each.key
    }
  ]
  
  alert_channel_ids = each.value.channels
}
```

### SPA Performance Alert

Monitor Single Page Application metrics:

```hcl
resource "instana_website_alert_config" "spa_performance" {
  name        = "SPA Route Change Performance"
  description = "Monitor SPA route transitions"
  severity    = "warning"
  website_id  = instana_website_monitoring_config.example.id
  
  tag_filter = "beacon.page.type@na EQUALS 'spa'"
  
  rule = {
    slowness = {
      metric_name = "routeChangeTime"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 1000
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.spa_team.id]
}
```

## Argument Reference

* `name` - Required - Name of the website alert configuration (max 256 characters)
* `description` - Required - Description of the alert configuration (max 65536 characters)
* `severity` - Optional - Severity of the alert. Values: `warning`, `critical`. Default: `warning`
* `triggering` - Optional - Boolean flag to trigger incidents. Default: `false`
* `website_id` - Required - Unique ID of the website to monitor (max 64 characters)
* `tag_filter` - Optional - Tag filter expression to limit monitoring scope [Details](#tag-filter-reference)
* `alert_channel_ids` - Optional - Set of alert channel IDs to notify
* `granularity` - Optional - Evaluation granularity in milliseconds. Values: `300000`, `600000`, `900000`, `1200000`, `800000`. Default: `600000`
* `rule` - Optional - Single alert rule configuration [Details](#rule-reference)
* `rules` - Optional - Multiple alert rules with individual thresholds [Details](#rules-reference)
* `threshold` - Optional - Threshold configuration for single rule [Details](#threshold-reference)
* `time_threshold` - Required - Time threshold configuration [Details](#time-threshold-reference)
* `custom_payload_fields` - Optional - List of custom payload fields for alert notifications [Details](#custom-payload-fields-reference)

**Note**: Either `rule` with `threshold` OR `rules` must be specified, but not both.

### Rule Reference

Exactly one of the following rule types must be configured:

* `slowness` - Optional - Rule based on page load slowness [Details](#slowness-rule-reference)
* `specific_js_error` - Optional - Rule based on specific JavaScript errors [Details](#specific-js-error-rule-reference)
* `status_code` - Optional - Rule based on HTTP status codes [Details](#status-code-rule-reference)
* `throughput` - Optional - Rule based on page view throughput [Details](#throughput-rule-reference)

#### Slowness Rule Reference

* `metric_name` - Required - Metric name (e.g., `onLoadTime`, `domContentLoadedTime`, `routeChangeTime`)
* `aggregation` - Required - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`

#### Specific JS Error Rule Reference

* `metric_name` - Required - Metric name (e.g., `jsErrors`)
* `aggregation` - Optional - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`
* `operator` - Required - Comparison operator. Values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`
* `value` - Optional - Value to identify the specific JavaScript error

#### Status Code Rule Reference

* `metric_name` - Required - Metric name (e.g., `httpStatusCode`)
* `aggregation` - Optional - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`
* `operator` - Required - Comparison operator. Values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`
* `value` - Required - HTTP status code value

#### Throughput Rule Reference

* `metric_name` - Required - Metric name (e.g., `pageViews`)
* `aggregation` - Optional - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`

### Rules Reference

For multiple alert rules with individual thresholds:

* `rule` - Required - Rule configuration (same structure as single rule)
* `threshold_operator` - Optional - Threshold operator
* `thresholds` - Required - Threshold configurations [Details](#thresholds-reference)

#### Thresholds Reference

* `warning` - Optional - Warning threshold configuration
* `critical` - Optional - Critical threshold configuration

Each threshold can have:
* `static` - Static threshold [Details](#static-threshold-reference)
* `adaptive_baseline` - Adaptive baseline threshold [Details](#adaptive-baseline-threshold-reference)
* `historic_baseline` - Historic baseline threshold [Details](#historic-baseline-threshold-reference)

### Threshold Reference

Exactly one of the following threshold types must be configured:

* `static` - Optional - Static threshold [Details](#static-threshold-reference)
* `adaptive_baseline` - Optional - Adaptive baseline threshold [Details](#adaptive-baseline-threshold-reference)
* `historic_baseline` - Optional - Historic baseline threshold [Details](#historic-baseline-threshold-reference)

#### Static Threshold Reference

* `operator` - Required - Comparison operator. Values: `>=`, `>`, `<=`, `<`, `==`
* `value` - Required - Threshold value (integer)

#### Adaptive Baseline Threshold Reference

* `operator` - Required - Comparison operator. Values: `>=`, `>`, `<=`, `<`, `==`
* `deviation_factor` - Required - Deviation factor (float)
* `adaptability` - Required - Adaptability factor (float)
* `seasonality` - Required - Seasonality pattern. Values: `WEEKLY`, `DAILY`

#### Historic Baseline Threshold Reference

* `deviation` - Required - Deviation factor (float)
* `seasonality` - Required - Seasonality pattern. Values: `WEEKLY`, `DAILY`
* `baseline` - Required - List of baseline configurations [Details](#baseline-reference)

##### Baseline Reference

* `day_of_week` - Required - Day of week. Values: `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY`
* `start` - Required - Start time (HH:MM format)
* `end` - Required - End time (HH:MM format)
* `baseline` - Required - Baseline value (float)

### Time Threshold Reference

Exactly one of the following time threshold types must be configured:

* `violations_in_sequence` - Optional - Violations in sequence [Details](#violations-in-sequence-reference)
* `violations_in_period` - Optional - Violations in period [Details](#violations-in-period-reference)
* `user_impact_of_violations_in_sequence` - Optional - User impact based [Details](#user-impact-reference)

#### Violations In Sequence Reference

* `time_window` - Optional - Time window in milliseconds

#### Violations In Period Reference

* `time_window` - Optional - Time window in milliseconds
* `violations` - Optional - Number of violations required

#### User Impact Reference

* `time_window` - Optional - Time window in milliseconds
* `impact_measurement_method` - Required - Impact measurement method. Values: `AGGREGATED`, `PER_WINDOW`
* `user_percentage` - Optional - Percentage of impacted users (0.0 to 1.0)
* `users` - Optional - Number of impacted users (integer > 0)

### Custom Payload Fields Reference

* `key` - Required - Key of the custom payload field
* `value` - Optional - Static string value. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - Dynamic value from tags [Details](#dynamic-value-reference)

#### Dynamic Value Reference

* `key` - Optional - Key to use in the payload
* `tag_name` - Required - Name of the tag to extract value from

### Tag Filter Reference

The **tag_filter** defines which entities should be monitored. It supports:

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

**Common Tag Filter Examples:**
* `beacon.page.name@na EQUALS '/checkout'` - Specific page
* `beacon.user.tier@na EQUALS 'premium'` - Premium users
* `beacon.device.type@na EQUALS 'mobile'` - Mobile devices
* `beacon.geo.region@na EQUALS 'APAC'` - Geographic region
* `beacon.browser.name@na EQUALS 'Chrome'` - Specific browser

## Attributes Reference

* `id` - The ID of the website alert configuration

## Import

Website alert configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_website_alert_config.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Severity `critical` triggers incidents when `triggering = true`
* Use tag filters to monitor specific pages, users, or segments
* Granularity defines the evaluation window size
* Adaptive baselines automatically adjust to traffic patterns
* Historic baselines compare against past performance
* User impact thresholds focus on affected user count
* Custom payload fields enhance alert notifications with context
* Multiple rules allow monitoring different metrics with individual thresholds
* Website alerts monitor real user monitoring (RUM) data
* Alerts evaluate at the configured granularity interval