# Website Alert Configuration Resource

Manages website alert configurations (Website Smart Alerts) in Instana. Website alerts monitor the performance and availability of your websites from the end-user perspective.

API Documentation: <https://instana.github.io/openapi/#operation/findActiveWebsiteAlertConfigs>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v6, they use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
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

#### NEW (v6.x) Syntax:
```hcl
resource "instana_website_alert_config" "example" {
  name = "Website Alert"
  rules = [
    {
      operator = ">="
      rule = {
        slowness = {
          aggregation = "P90"
          metric_name = "onLoadTime"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 5
          }
        }
      }
    },
  ]
  
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
  triggering  = false
  website_id  = instana_website_monitoring_config.example.id
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
  granularity       = 600000
  rules = [
    {
      operator = ">="
      rule = {
        slowness = {
          aggregation = "P90"
          metric_name = "onLoadTime"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 5
          }
        }
      }
    },
  ]

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
  website_id  = instana_website_monitoring_config.example.id
  tag_filter = "endpoint.name@dest NOT_EQUAL 'x'"  
   rules = [
    {
      operator = ">="
      rule = {
        slowness = {
          aggregation = "P90"
          metric_name = "onLoadTime"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 5
          }
        }
      }
    },
  ]
  
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
### Status Code Alert

Monitor HTTP status codes:

```hcl
resource "instana_website_alert_config" "status_code" {
  name        = "4xx Error Alert"
  description = "Alert on client errors"
  website_id  = instana_website_monitoring_config.example.id
  rules = [
    {
      operator = ">="
      rule = {
        status_code = {
          aggregation = "SUM"
          metric_name = "httpxxx"
          operator    = "EQUALS"
          value       = "404"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 10
          }
        }
      }
    },
  ]
  
  time_threshold = {
    violations_in_period = {
      time_window = 600000
      violations  = 2
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.ops.id]
}
```

### Alert with Adaptive Baseline

Use adaptive baseline for dynamic thresholds:

```hcl
resource "instana_website_alert_config" "adaptive_slowness" {
  name        = "4xx Error Alert"
  description = "Alert on client errors"
  website_id  = instana_website_monitoring_config.example.id
  rules = [
    {
      operator = ">="
      rule = {
        status_code = {
          aggregation = "SUM"
          metric_name = "httpxxx"
          operator    = "EQUALS"
          value       = "404"
        }
      }
      threshold = {
        warning = {
          adaptive_baseline = {
            adaptability     = 1
            deviation_factor = 3
            seasonality      = "AUTO"
          }
        }
      }
    },
  ]
  
  time_threshold = {
    violations_in_period = {
      time_window = 600000
      violations  = 2
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.ops.id]
}
```

### Alert with Custom Payload

Add custom fields to alert notifications:

```hcl
resource "instana_website_alert_config" "with_custom_payload" {
  name        = "4xx Error Alert"
  description = "Alert on client errors"
  website_id  = instana_website_monitoring_config.example.id
  rules = [
    {
      operator = ">="
      rule = {
        status_code = {
          aggregation = "SUM"
          metric_name = "httpxxx"
          operator    = "EQUALS"
          value       = "404"
        }
      }
      threshold = {
        warning = {
          adaptive_baseline = {
            adaptability     = 1
            deviation_factor = 3
            seasonality      = "AUTO"
          }
        }
      }
    },
  ]
  
  time_threshold = {
    violations_in_period = {
      time_window = 600000
      violations  = 2
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.ops.id]

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
      key = "user_segment"
      dynamic_value = {
        key      = "segment"
        tag_name = "beacon.user.segment"
      }
    }
  ]
  }
```

## Argument Reference

* `name` - Required - Name of the website alert configuration (max 256 characters)
* `description` - Required - Description of the alert configuration (max 65536 characters)
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
* Use tag filters to monitor specific pages, users, or segments
* Granularity defines the evaluation window size
* Adaptive baselines automatically adjust to traffic patterns
* Historic baselines compare against past performance
* User impact thresholds focus on affected user count
* Custom payload fields enhance alert notifications with context
* Multiple rules allow monitoring different metrics with individual thresholds
* Website alerts monitor real user monitoring (RUM) data
* Alerts evaluate at the configured granularity interval
