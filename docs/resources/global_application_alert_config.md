# Global Application Alert Configuration Resource

Manages global application alert configurations (Global Application Smart Alerts) in Instana. Global application alerts monitor application performance across all applications or specific application scopes.

API Documentation: <https://instana.github.io/openapi/#tag/Global-Application-Alert-Configuration>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

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
  
  application = [ {
    application_id = "app-id"
    inclusive = true
  }]
  boundary_scope       = "INBOUND"
  rules = [
    {
      rule = {
        slowness = {
          aggregation = "P90"
          metric_name = "latency"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 5
          }
        }
      }
      threshold_operator = ">="
    }
  ]
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_field = [
    {
      key = "environment"
      value = "production"
    }
  ]
  # rest of the configuration
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
  alert_channels = {
    CRITICAL = ["critical_email_channel"]
    WARNING  = ["warning_email_channel"]
  }
  application = [
    {
      application_id = "-1E6OCrFTZazfuwo34wUzw"
      inclusive      = true
      service = [
        {
      service_id = "payment-service-id"
      inclusive  = true
      endpoint = [{
        endpoint_id = "checkout-endpoint-id"
        inclusive   = true
      }]
       }
      ]
    }
  ]
  boundary_scope       = "INBOUND"
  description          = "slowness_basic"
  evaluation_type      = "PER_AP"
  grace_period         = 300000
  granularity          = 300000
  name                 = "slowness_basic"
  rules = [
    {
      rule = {
        slowness = {
          aggregation = "P90"
          metric_name = "latency"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 1
          }
        }
      }
      threshold_operator = ">="
    },
  ]
  time_threshold = {
    violations_in_sequence = {
      time_window = 300000
    }
  }
  tag_filter = "(service.name@dest EQUALS 'robotshop-eum-service-2' OR service.name@dest EQUALS 'robotshop-eum-service-1')"
  triggering = false
}
```

## Generating Configuration from Existing Resources

If you have already created a global application alert configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your global application alert configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_global_application_alert_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual global application alert configuration ID from Instana

3. **Generate the Configuration**: Run the following Terraform command:

```bash
terraform plan -generate-config-out=generated.tf
```

This will:
- Import the existing resource state
- Generate the complete Terraform configuration in `generated.tf`
- Show you what will be imported

4. **Review and Apply**: Review the generated configuration in `generated.tf` and make any necessary adjustments.

   - **To import the existing resource**: Keep the import block and run `terraform apply`. This will import the resource into your Terraform state and link it to the existing Instana resource.
   
   - **To create a new resource**: If you only need the configuration structure as a template, remove the import block from your configuration. Modify the generated configuration as needed, and when you run `terraform apply`, it will create a new resource in Instana instead of importing the existing one.

```bash
terraform apply
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
* Request impact thresholds focus on affected request count
* Custom payload fields enhance alert notifications with context
* Use `inclusive = false` to exclude specific applications, services, or endpoints
* Evaluation types: `PER_AP` (per application), `PER_AP_SERVICE` (per service), `PER_AP_ENDPOINT` (per endpoint)
