# Global Application Alert Configuration Resource

Manages global application alert configurations (Global Application Smart Alerts) in Instana. Global application alerts monitor application performance across all applications or specific application scopes.

API Documentation: <https://instana.github.io/openapi/#tag/Global-Application-Alert-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.


## Migration Guide (v5 to v6)

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v6, they use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_global_application_alert_config" "example" {
  name = "Global Alert"
  
  application {
    application_id = "app-id"
    inclusive = true
  }
  
  rule {
    slowness {
      metric_name = "latency"
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
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_global_application_alert_config" "example" {
  name = "Global Alert"
  
  application = {
    application_id = "app-id"
    inclusive = true
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
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
    }
  ]
}
```

### Key Syntax Changes

1. **Application**: `application { }` → `application = { }`
2. **Service**: `service { }` → `service = { }`
3. **Endpoint**: `endpoint { }` → `endpoint = { }`
4. **Rule**: `rule { }` → `rule = { }`
5. **Threshold**: `threshold { }` → `threshold = { }`
6. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
7. **Custom Payload Fields**: `custom_payload_field { }` (multiple) → `custom_payload_fields = [{ }]` (list)
8. **Aggregation**: Case-insensitive but lowercase recommended: `P90` → `p90`

## Example Usage

### Basic Slowness Alert

Monitor application latency globally:

```hcl
resource "instana_global_application_alert_config" "slowness_basic" {
  name            = "Global Latency Alert"
  description     = "Alert on high latency across all applications"
  boundary_scope  = "ALL"
  severity        = "warning"
  triggering      = false
  evaluation_type = "PER_AP"
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
  granularity       = 600000
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
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
}
```

### Error Rate Alert

Monitor error rates across applications:

```hcl
resource "instana_global_application_alert_config" "error_rate" {
  name            = "Global Error Rate Alert"
  description     = "Alert on elevated error rates"
  boundary_scope  = "INBOUND"
  severity        = "critical"
  triggering      = true
  evaluation_type = "PER_AP"
  
  include_internal  = false
  include_synthetic = false
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    error_rate = {
      metric_name = "errors"
      aggregation = "sum"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 5
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
    instana_alerting_channel_slack.ops.id
  ]
}
```

### Service-Level Alert

Monitor specific services within an application:

```hcl
resource "instana_global_application_alert_config" "service_level" {
  name            = "Payment Service Alert"
  description     = "Monitor payment service performance"
  boundary_scope  = "ALL"
  severity        = "critical"
  evaluation_type = "PER_AP_SERVICE"
  
  application = {
    application_id = instana_application_config.ecommerce.id
    inclusive      = true
    service = {
      service_id = "payment-service-id"
      inclusive  = true
    }
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
      aggregation = "p95"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 500
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_pagerduty.payment_team.id]
}
```

### Endpoint-Level Alert

Monitor specific endpoints:

```hcl
resource "instana_global_application_alert_config" "endpoint_level" {
  name            = "Checkout Endpoint Alert"
  description     = "Monitor checkout endpoint performance"
  boundary_scope  = "ALL"
  severity        = "critical"
  evaluation_type = "PER_AP_ENDPOINT"
  
  application = {
    application_id = instana_application_config.ecommerce.id
    inclusive      = true
    service = {
      service_id = "api-service-id"
      inclusive  = true
      endpoint = {
        endpoint_id = "checkout-endpoint-id"
        inclusive   = true
      }
    }
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
      aggregation = "p99"
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
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.checkout_team.id]
}
```

### Alert with Tag Filter

Use tag filters to scope monitoring:

```hcl
resource "instana_global_application_alert_config" "with_tag_filter" {
  name            = "Production API Alert"
  description     = "Monitor production API calls"
  boundary_scope  = "INBOUND"
  severity        = "warning"
  evaluation_type = "PER_AP"
  
  tag_filter = "call.type@na EQUALS 'HTTP' AND call.http.status@na GREATER_OR_EQUAL_THAN 500"
  
  application = {
    application_id = instana_application_config.api.id
    inclusive      = true
  }
  
  rule = {
    errors = {
      metric_name = "errors"
      aggregation = "sum"
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
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_email.api_team.id]
}
```

### Status Code Alert

Monitor specific HTTP status codes:

```hcl
resource "instana_global_application_alert_config" "status_code" {
  name            = "5xx Error Alert"
  description     = "Alert on server errors"
  boundary_scope  = "ALL"
  severity        = "critical"
  evaluation_type = "PER_AP"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    status_code = {
      metric_name      = "calls"
      aggregation      = "sum"
      status_code_start = 500
      status_code_end   = 599
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 20
    }
  }
  
  time_threshold = {
    violations_in_period = {
      time_window = 300000
      violations  = 2
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_pagerduty.oncall.id]
}
```

### Throughput Alert

Monitor call throughput:

```hcl
resource "instana_global_application_alert_config" "throughput" {
  name            = "Low Throughput Alert"
  description     = "Alert when traffic drops"
  boundary_scope  = "INBOUND"
  severity        = "warning"
  evaluation_type = "PER_AP"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    throughput = {
      metric_name = "calls"
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

### Log-Based Alert

Monitor application logs:

```hcl
resource "instana_global_application_alert_config" "log_alert" {
  name            = "Error Log Alert"
  description     = "Alert on error logs"
  boundary_scope  = "ALL"
  severity        = "warning"
  evaluation_type = "PER_AP"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    logs = {
      metric_name = "logs"
      aggregation = "sum"
      level       = "ERROR"
      message     = "OutOfMemoryError"
      operator    = "CONTAINS"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 5
    }
  }
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_slack.dev_team.id]
}
```

### Alert with Historic Baseline

Use historic baseline for comparison:

```hcl
resource "instana_global_application_alert_config" "historic_baseline" {
  name            = "Historic Baseline Alert"
  description     = "Compare against historical performance"
  boundary_scope  = "ALL"
  severity        = "warning"
  evaluation_type = "PER_AP"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
      aggregation = "p90"
    }
  }
  
  threshold = {
    historic_baseline = {
      operator        = ">="
      deviation_factor = 2.0
      seasonality     = "WEEKLY"
      baseline = [
        {
          day_of_week = "MONDAY"
          start       = "09:00"
          end         = "17:00"
          baseline    = 500
        },
        {
          day_of_week = "FRIDAY"
          start       = "09:00"
          end         = "17:00"
          baseline    = 600
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

### Request Impact Alert

Alert based on request impact:

```hcl
resource "instana_global_application_alert_config" "request_impact" {
  name            = "High Request Impact Alert"
  description     = "Alert when many requests are affected"
  boundary_scope  = "ALL"
  severity        = "critical"
  evaluation_type = "PER_AP"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
      aggregation = "p90"
    }
  }
  
  threshold = {
    static = {
      operator = ">="
      value    = 2000
    }
  }
  
  time_threshold = {
    request_impact = {
      time_window = 600000
      request     = 1000
    }
  }
  
  alert_channel_ids = [instana_alerting_channel_pagerduty.oncall.id]
}
```

### Alert with Custom Payload

Add custom fields to alert notifications:

```hcl
resource "instana_global_application_alert_config" "with_custom_payload" {
  name            = "Alert with Custom Payload"
  description     = "Alert with additional context"
  boundary_scope  = "ALL"
  severity        = "warning"
  evaluation_type = "PER_AP"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
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
  
  custom_payload_fields = [
    {
      key   = "environment"
      value = "production"
    },
    {
      key   = "team"
      value = "backend"
    },
    {
      key   = "runbook"
      value = "https://wiki.example.com/runbooks/latency"
    },
    {
      key = "region"
      dynamic_value = {
        key      = "region"
        tag_name = "aws.region"
      }
    }
  ]
  
  alert_channel_ids = [instana_alerting_channel_slack.backend.id]
}
```

### Multi-Application Alert

Monitor multiple applications:

```hcl
locals {
  critical_apps = {
    payment = {
      app_id    = instana_application_config.payment.id
      threshold = 500
      channels  = [instana_alerting_channel_pagerduty.payment.id]
    }
    auth = {
      app_id    = instana_application_config.auth.id
      threshold = 300
      channels  = [instana_alerting_channel_pagerduty.auth.id]
    }
  }
}

resource "instana_global_application_alert_config" "critical_apps" {
  for_each = local.critical_apps

  name            = "${title(each.key)} Service Alert"
  description     = "Monitor ${each.key} service performance"
  boundary_scope  = "ALL"
  severity        = "critical"
  triggering      = true
  evaluation_type = "PER_AP"
  
  application = {
    application_id = each.value.app_id
    inclusive      = true
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
      aggregation = "p95"
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
      time_window = 300000
    }
  }
  
  custom_payload_fields = [
    {
      key   = "service"
      value = each.key
    }
  ]
  
  alert_channel_ids = each.value.channels
}
```

### Exclude Specific Services

Exclude services from monitoring:

```hcl
resource "instana_global_application_alert_config" "exclude_services" {
  name            = "Application Alert (Excluding Test Services)"
  description     = "Monitor all services except test services"
  boundary_scope  = "ALL"
  severity        = "warning"
  evaluation_type = "PER_AP_SERVICE"
  
  application = {
    application_id = instana_application_config.example.id
    inclusive      = true
    service = {
      service_id = "test-service-id"
      inclusive  = false  # Exclude this service
    }
  }
  
  rule = {
    slowness = {
      metric_name = "latency"
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
  
  alert_channel_ids = [instana_alerting_channel_email.example.id]
}
```

### Complex Scope Configuration

Monitor with complex application/service/endpoint scope:

```hcl
resource "instana_global_application_alert_config" "complex_scope" {
  name            = "Complex Scope Alert"
  description     = "Monitor specific endpoints in specific services"
  boundary_scope  = "ALL"
  severity        = "critical"
  evaluation_type = "PER_AP_ENDPOINT"
  
  application = {
    application_id = instana_application_config.ecommerce.id
    inclusive      = true
    service = {
      service_id = "api-gateway-id"
      inclusive  = true
      endpoint = {
        endpoint_id = "critical-endpoint-id"
        inclusive   = true
      }
    }
  }
  
  tag_filter = "call.http.method@na EQUALS 'POST'"
  
  rule = {
    slowness = {
      metric_name = "latency"
      aggregation = "p99"
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
      time_window = 300000
    }
  }
  
  alert_channel_ids = [
    instana_alerting_channel_pagerduty.oncall.id,
    instana_alerting_channel_slack.critical.id
  ]
}
```

## Argument Reference

* `name` - Required - Name of the global application alert configuration
* `description` - Required - Description of the alert configuration
* `severity` - Required - Severity of the alert. Values: `critical`, `warning`
* `boundary_scope` - Required - Boundary scope of the alert. Values: `INBOUND`, `ALL`, `DEFAULT`
* `triggering` - Optional - Boolean flag to trigger incidents. Default: `false`
* `include_internal` - Optional - Include internal calls in scope. Default: `false`
* `include_synthetic` - Optional - Include synthetic calls in scope. Default: `false`
* `alert_channel_ids` - Optional - Set of alert channel IDs to notify
* `granularity` - Optional - Evaluation granularity in milliseconds. Values: `300000`, `600000`, `900000`, `1200000`, `800000`. Default: `600000`
* `evaluation_type` - Required - Evaluation type. Values: `PER_AP`, `PER_AP_SERVICE`, `PER_AP_ENDPOINT`
* `tag_filter` - Optional - Tag filter expression to limit monitoring scope [Details](#tag-filter-reference)
* `application` - Required - Application scope configuration [Details](#application-reference)
* `rule` - Required - Alert rule configuration [Details](#rule-reference)
* `threshold` - Required - Threshold configuration [Details](#threshold-reference)
* `time_threshold` - Required - Time threshold configuration [Details](#time-threshold-reference)
* `custom_payload_fields` - Optional - List of custom payload fields for alert notifications [Details](#custom-payload-fields-reference)

### Application Reference

* `application_id` - Required - ID of the application
* `inclusive` - Required - Boolean flag to include (true) or exclude (false) this application
* `service` - Optional - Service scope configuration [Details](#service-reference)

#### Service Reference

* `service_id` - Required - ID of the service
* `inclusive` - Required - Boolean flag to include (true) or exclude (false) this service
* `endpoint` - Optional - Endpoint scope configuration [Details](#endpoint-reference)

##### Endpoint Reference

* `endpoint_id` - Required - ID of the endpoint
* `inclusive` - Required - Boolean flag to include (true) or exclude (false) this endpoint

### Rule Reference

Exactly one of the following rule types must be configured:

* `error_rate` - Optional - Rule based on error rate [Details](#error-rate-rule-reference)
* `errors` - Optional - Rule based on number of errors [Details](#errors-rule-reference)
* `logs` - Optional - Rule based on logs [Details](#logs-rule-reference)
* `slowness` - Optional - Rule based on slowness [Details](#slowness-rule-reference)
* `status_code` - Optional - Rule based on HTTP status code [Details](#status-code-rule-reference)
* `throughput` - Optional - Rule based on throughput [Details](#throughput-rule-reference)

#### Error Rate Rule Reference

* `metric_name` - Required - Metric name
* `aggregation` - Optional - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distribution`, `distinct_count`, `sum_positive`

#### Errors Rule Reference

* `metric_name` - Required - Metric name
* `aggregation` - Optional - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distribution`, `distinct_count`, `sum_positive`

#### Logs Rule Reference

* `metric_name` - Required - Metric name
* `aggregation` - Required - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distribution`, `distinct_count`, `sum_positive`
* `level` - Required - Log level. Values: `WARN`, `ERROR`, `ANY`
* `message` - Optional - Log message to match
* `operator` - Required - Comparison operator. Values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`

#### Slowness Rule Reference

* `metric_name` - Required - Metric name
* `aggregation` - Required - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distribution`, `distinct_count`, `sum_positive`

#### Status Code Rule Reference

* `metric_name` - Required - Metric name
* `aggregation` - Required - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distribution`, `distinct_count`, `sum_positive`
* `status_code_start` - Optional - Minimum HTTP status code
* `status_code_end` - Optional - Maximum HTTP status code

#### Throughput Rule Reference

* `metric_name` - Required - Metric name
* `aggregation` - Required - Aggregation function. Values: `sum`, `mean`, `max`, `min`, `p25`, `p50`, `p75`, `p90`, `p95`, `p98`, `p99`, `p99_9`, `p99_99`, `distribution`, `distinct_count`, `sum_positive`

### Threshold Reference

Exactly one of the following threshold types must be configured:

* `static` - Optional - Static threshold [Details](#static-threshold-reference)
* `historic_baseline` - Optional - Historic baseline threshold [Details](#historic-baseline-threshold-reference)

#### Static Threshold Reference

* `operator` - Required - Comparison operator. Values: `>`, `>=`, `<`, `<=`
* `value` - Optional - Threshold value (float)
* `last_updated` - Optional - Last updated timestamp

#### Historic Baseline Threshold Reference

* `operator` - Required - Comparison operator. Values: `>`, `>=`, `<`, `<=`
* `deviation_factor` - Optional - Deviation factor (float)
* `seasonality` - Required - Seasonality pattern. Values: `WEEKLY`, `DAILY`
* `baseline` - Optional - List of baseline configurations [Details](#baseline-reference)
* `last_updated` - Optional - Last updated timestamp

##### Baseline Reference

* `day_of_week` - Required - Day of week. Values: `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY`
* `start` - Required - Start time (HH:MM format)
* `end` - Required - End time (HH:MM format)
* `baseline` - Required - Baseline value (float)

### Time Threshold Reference

Exactly one of the following time threshold types must be configured:

* `violations_in_sequence` - Optional - Violations in sequence [Details](#violations-in-sequence-reference)
* `violations_in_period` - Optional - Violations in period [Details](#violations-in-period-reference)
* `request_impact` - Optional - Request impact based [Details](#request-impact-reference)

#### Violations In Sequence Reference

* `time_window` - Optional - Time window in milliseconds

#### Violations In Period Reference

* `time_window` - Optional - Time window in milliseconds
* `violations` - Optional - Number of violations required

#### Request Impact Reference

* `time_window` - Optional - Time window in milliseconds
* `request` - Optional - Number of requests threshold

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

## Attributes Reference

* `id` - The ID of the global application alert configuration

## Import

Global application alert configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_global_application_alert_config.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Global alerts monitor across all applications or specific application scopes
* Use `evaluation_type` to control alert granularity (per application, service, or endpoint)
* `boundary_scope` controls which call boundaries are monitored
* Set `triggering = true` to create incidents for critical alerts
* Use `include_internal` and `include_synthetic` to control call scope
* Tag filters provide fine-grained control over monitored calls
* Historic baselines compare against past performance patterns
* Request impact thresholds focus on affected request count
* Custom payload fields enhance alert notifications with context
* Use `inclusive = false` to exclude specific applications, services, or endpoints
* Evaluation types: `PER_AP` (per application), `PER_AP_SERVICE` (per service), `PER_AP_ENDPOINT` (per endpoint)
