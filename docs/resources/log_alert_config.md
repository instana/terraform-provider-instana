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
    alert_type  = "log.count"
    aggregation = "SUM"
    threshold {
      critical {
        static {
          value = 500
        }
      }
    }
    threshold_operator = ">="
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
  description = "Log Alert description"

  alert_channels = {
    warning = ["channel-1"] # replace with actual channel Ids
    critical = ["channel-2"]
  }
  
  group_by = [
    {
      tag_name = "kubernetes.namespace.name"
    }
  ]
  
  rules = {
    metric_name = "log.count"
    alert_type  = "log.count"
    aggregation = "SUM"
    threshold = {
      critical = {
        static = {
          value = 500
        }
      }
    }
    threshold_operator = ">="
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

### Log Count Alert

```hcl
resource "instana_log_alert_config" "error_logs" {
  name = "High Error Log Count - $${severity}" # Use double $$ to define placeholders
  description = "Alert when error logs exceed threshold"
  granularity = 600000
  
  alert_channels = {
    critical = ["ops-team-channel"] # replace with actual channel Ids
    warning = ["dev-team-channel"]
  }
  rules = {
    aggregation = "SUM"
    alert_type  = "log.count"
    metric_name = "logCount"
    threshold = {
      critical = {
        static = {
          value = 2
        }
      }
      warning = {
        static = {
          value = 1
        }
      }
    }
    threshold_operator = ">="
  }
  tag_filter = "log.exception.type@na EQUALS 'error'"
  group_by = [
    {
      tag_name = "kubernetes.namespace.name"
    }
  ]
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
    custom_payload_field = [
    {
      dynamic_value = {
        tag_name = "log.level"
      }
      key   = "log"
    },
  ]
}
```

## Generating Configuration from Existing Resources

If you have already created a log alert configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your log alert configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_log_alert_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual log alert configuration ID from Instana

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
