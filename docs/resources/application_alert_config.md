# Application Alert Configuration Resource

Management of application alert configurations (Application Smart Alerts).

API Documentation: <https://instana.github.io/openapi/#tag/Application-Alert-Configuration>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- `application` → `application` (now a list with `= [{ }]`)
- `service` → `service` (nested list)
- `endpoint` → `endpoint` (nested list)
- `rule` → `rules` (list with new structure)
- `threshold` → `threshold` (nested in rules, supports multiple severity levels)
- New `threshold_operator` field in rules
- `time_threshold` now uses attribute syntax with `= { }`
- Enhanced support for both static and adaptive baseline thresholds
- `alert_channels` now supports severity-based routing (map of severity to channel IDs)

#### OLD (v5.x) Syntax:

```hcl
resource "instana_application_alert_config" "example" {
  name = "test-alert"
  
  application {
    application_id = "app-123"
    inclusive      = true
    
    service {
      service_id = "svc-456"
      inclusive  = true
    }
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
      value    = 5.0
    }
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_application_alert_config" "example" {
  name = "test-alert"
  
  application = [{
    application_id = "app-123"
    inclusive      = true
    
    service = [{
      service_id = "svc-456"
      inclusive  = true
    }]
  }]
  
  rules = [{
    rule = {
      slowness = {
        metric_name = "latency"
        aggregation = "P90"
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
  }]
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  # Rest of the configuration
}
```


### Basic Alert

```hcl
resource "instana_application_alert_config" "application_alert_config" {
  name                 = "Calls are slower than usual - $${severity}" # Use double $$ to define placeholders
  application = [
    {
      application_id = "eQjat8RMR3-5lb376rV5kA" # Replace with your own value
      inclusive      = true
    },
  ]
  boundary_scope       = "INBOUND"
  description          = "Calls are slower or equal to 5 ms based on latency (90th)."
  enabled              = true
  evaluation_type      = "PER_AP"
  granularity          = 600000
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
    },
  ]
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
}
```

### Status code Alert with Tag Filter

```hcl
resource "instana_application_alert_config" "application_alert_config" {
  alert_channels = {
    CRITICAL = ["alert_channel_1"] # Replace with your own value
    WARNING  = ["alert_channel_2"] # Replace with your own value
  }
  application = [
    {
      application_id = "eQjat8RMR3-5lb376rV5kA" # Replace with your own value
      inclusive      = true
      service = [
      ]
    },
  ]
  boundary_scope       = "ALL"
  description          = "Occurrences of HTTP Status Code 5XX (Server Error) is higher or equal to ." # Replace with your own value
  enabled              = true
  evaluation_type      = "PER_AP"
  granularity          = 600000
  name                 = "AN Smart alert" # Replace with your own value
  rules = [
    {
      rule = {
        status_code = {
          aggregation       = "SUM"
          metric_name       = "calls"
          status_code_end   = 599
          status_code_start = 500
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
    },
  ]
  tag_filter = "(call.tag:'accessType'@na EQUALS 'test' AND call.metric:'codeAmount'@na EQUALS '1')" # Replace with your own value
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  triggering = false
}
```


### Alert with Custom Payload

```hcl
resource "instana_application_alert_config" "alert_config_3" {
  alert_channels = {
    WARNING = ["d35rlokgvalloa2fnps0"] # Replace with your own value
  }
  application = [
    {
      application_id = "eQjat8RMR3-5lb376rV5kA" # Replace with your own value
      inclusive      = true
    },
  ]
  boundary_scope = "ALL"
  custom_payload_field = [
    {
      key           = "test" # Replace with your own value
      value         = "test123" # Replace with your own value
    },
    {
      dynamic_value = {
        key      = "stage" # Replace with your own value
        tag_name = "aws.tag" # Replace with your own value
      }
      key   = "stage" # Replace with your own value
    },
  ]
  description       = "test-alert-description" # Replace with your own value
  enabled           = true
  evaluation_type   = "PER_AP"
  granularity       = 600000
  name              = "test-alert" # Replace with your own value
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
            value = 100
          }
        }
      }
      threshold_operator = ">"
    },
  ]
  tag_filter = "call.type@na EQUALS 'HTTP'" # Replace with your own value
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  triggering = false
}

```

### Alert with throughput rule

```hcl
resource "instana_application_alert_config" "throughput_application_alert_config" {
  application = [
    {
      application_id = "hpY3YjbxSm6jlrUqIEcl8w" # replace with valid application Id
      inclusive      = false
      service = [
        {
          endpoint = [
            {
              endpoint_id = "8ZZ5e1ZMRNPqvzQOaFkh-InyULM" # replace with valid endpoint Id
              inclusive   = true
            },
          ]
          inclusive  = false
          service_id = "0cbc4b430e327759f5fc2b02781c13940d0b439e" # replace with valid service Id
        },
      ]
    },
  ]
  boundary_scope       = "ALL"
  description          = "The number of calls is lower or equal to  calls."
  enabled              = true
  evaluation_type      = "PER_AP"
  granularity          = 600000
  include_internal     = false
  include_synthetic    = false
  name                 = "throughput based alert"
  rules = [
    {
      rule = {
        throughput = {
          aggregation = "SUM"
          metric_name = "calls"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 100
          }
        }
      }
      threshold_operator = "<="
    },
  ]
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  triggering = false
}
```

### Alert with errors rule

```hcl
resource "instana_application_alert_config" "errors_application_alert_config" {
  application = [
    {
      application_id = "A2wMMcEyQ36O-KUsDk4eYg" # replace with valid application Id
      inclusive      = true
    },
  ]
  boundary_scope       = "ALL"
  description          = "errors based alert"
  enabled              = true
  evaluation_type      = "PER_AP_SERVICE"
  granularity          = 60000
  include_internal     = false
  include_synthetic    = false
  name                 = "errors based alert"
  rules = [
    {
      rule = {
        errors = {
          aggregation = "MEAN"
          metric_name = "errors"
        }
      }
      threshold = {
        warning = {
          static = {
            value = 0.01
          }
        }
      }
      threshold_operator = ">="
    },
  ]
  time_threshold = {
    violations_in_sequence = {
      time_window = 60000
    }
  }
  triggering = true
}
```

## Generating Configuration from Existing Resources

If you have already created a application alert configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your application alert configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_application_alert_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual application alert configuration ID from Instana

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

* `name` - Required - The name for the application alert configuration
* `description` - Required - The description text of the application alert config
* `boundary_scope` - Required - The boundary scope of the application alert config. Allowed values: `INBOUND`, `ALL`, `DEFAULT`
* `evaluation_type` - Required - The evaluation type of the application alert config. Allowed values: `PER_AP`, `PER_AP_SERVICE`, `PER_AP_ENDPOINT`
* `enabled` - Optional - default `true` - Flag to indicate whether the alert configuration is enabled or not
* `triggering` - Optional - default `false` - Flag to indicate whether also an Incident is triggered or not
* `include_internal` - Optional - default `false` - Flag to indicate whether also internal calls are included in the scope or not
* `include_synthetic` - Optional - default `false` - Flag to indicate whether also synthetic calls are included in the scope or not
* `granularity` - Optional - default `600000` - The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used. Allowed values: `300000`, `600000`, `900000`, `1200000`, `1800000`
* `grace_period` - Optional - The duration (in milliseconds) for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires
* `tag_filter` - Optional - The tag filter of the application alert config. [Details](#tag-filter-argument-reference)
* `application` - Required - Selection/Set of applications in scope (list). [Details](#applications-argument-reference)
* `rules` - Required - List of rules where each rule is associated with multiple thresholds and their corresponding severity levels (list). [Details](#rules-argument-reference)
* `time_threshold` - Required - Indicates the type of violation of the defined threshold (object). [Details](#time-threshold-argument-reference)
* `alert_channels` - Optional - Map of alert channel IDs associated with severity levels (map of sets). Keys: `warning`, `critical`
* `custom_payload_field` - Optional - An optional list of custom payload fields (static key/value pairs or dynamic values added to the event). [Details](#custom-payload-field-argument-reference)

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

### Applications Argument Reference

* `application_id` - Required - ID of the included application
* `inclusive` - Required - Defines whether this node and his child nodes are included (true) or excluded (false)
* `service` - Optional - Selection of services in scope (list). [Details](#services-argument-reference)

#### Services Argument Reference

* `service_id` - Required - ID of the included service
* `inclusive` - Required - Defines whether this node and his child nodes are included (true) or excluded (false)
* `endpoint` - Optional - Selection of endpoints in scope (list). [Details](#endpoints-argument-reference)

##### Endpoints Argument Reference

* `endpoint_id` - Required - ID of the included endpoint
* `inclusive` - Required - Defines whether this node and his child nodes are included (true) or excluded (false)

### Rules Argument Reference

Each rule object contains:

* `rule` - Required - The rule configuration (object). Exactly one of the rule types below must be configured
* `threshold` - Required - Threshold configuration for different severity levels (object). [Details](#thresholds-argument-reference)
* `threshold_operator` - Optional - The operator to apply for threshold comparison. Allowed values: `>`, `>=`, `<`, `<=`

#### Rule Types

Exactly one of the elements below must be configured within the `rule` object:

* `error_rate` - (Deprecated - please use errors instead)Optional - Rule based on the error rate. [Details](#error-rate-rule-argument-reference)
* `errors` - Optional - Rule based on the number of errors. [Details](#errors-rule-argument-reference)
* `logs` - Optional - Rule based on logs. [Details](#logs-rule-argument-reference)
* `slowness` - Optional - Rule based on the slowness. [Details](#slowness-rule-argument-reference)
* `status_code` - Optional - Rule based on the HTTP status code. [Details](#status-code-rule-argument-reference)
* `throughput` - Optional - Rule based on the throughput. [Details](#throughput-rule-argument-reference)

##### Error Rate Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Optional - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`

##### Errors Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Optional - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`

##### Logs Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `level` - Required - The log level for which this rule applies to. Supported values: `WARN`, `ERROR`, `ANY`
* `message` - Optional - The log message for which this rule applies to
* `operator` - Required - The operator which will be applied to evaluate this rule. Supported values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`

##### Slowness Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`

##### Status Code Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `status_code_start` - Optional - minimal HTTP status code applied for this rule
* `status_code_end` - Optional - maximum HTTP status code applied for this rule

##### Throughput Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`

### Thresholds Argument Reference

The thresholds object can contain:

* `warning` - Optional - Warning threshold configuration (object). [Details](#threshold-level-argument-reference)
* `critical` - Optional - Critical threshold configuration (object). [Details](#threshold-level-argument-reference)

At least one of `warning` or `critical` must be specified.

#### Threshold Level Argument Reference

Exactly one of the elements below must be configured:

* `static` - Optional - Static threshold definition (object). [Details](#static-threshold-argument-reference)
* `adaptive_baseline` - Optional - Adaptive baseline threshold definition (object). [Details](#adaptive-baseline-threshold-argument-reference)

##### Static Threshold Argument Reference

* `value` - Required - The value of the static threshold (integer)

##### Adaptive Baseline Threshold Argument Reference

* `deviation_factor` - Optional - The deviation factor for the adaptive baseline (float)
* `adaptability` - Optional - The adaptability factor (float, 0.0 to 1.0)
* `seasonality` - Optional - The seasonality of the adaptive baseline.

### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - The dynamic value of the custom payload field (object). Either `value` or `dynamic_value` must be defined. [Details](#dynamic-custom-payload-field-value)

#### Dynamic Custom Payload Field Value

* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload

### Time Threshold Argument Reference

Exactly one of the elements below must be configured:

* `request_impact` - Optional - Time threshold based on request impact (object). [Details](#request-impact-time-threshold-argument-reference)
* `violations_in_period` - Optional - Time threshold based on violations in period (object). [Details](#violations-in-period-time-threshold-argument-reference)
* `violations_in_sequence` - Optional - Time threshold based on violations in sequence (object). [Details](#violations-in-sequence-time-threshold-argument-reference)

#### Request Impact Time Threshold Argument Reference

* `time_window` - Required - The time window of the time threshold (milliseconds)
* `requests` - Required - The number of requests in the given window

#### Violations In Period Time Threshold Argument Reference

* `time_window` - Required - The time window of the time threshold (milliseconds)
* `violations` - Required - The violations appeared in the period (1-12)

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Required - The time window of the time threshold (milliseconds)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the application alert configuration

## Import

Application Alert Configs can be imported using the `id`, e.g.:

```bash
terraform import instana_application_alert_config.example 60845e4e5e6b9cf8fc2868da
```

## Best Practices

### Evaluation Types

Choose the appropriate evaluation type:

- `PER_AP` - Alert on application-level metrics (overall application health)
- `PER_AP_SERVICE` - Alert per service (identify problematic services)
- `PER_AP_ENDPOINT` - Alert per endpoint (pinpoint specific endpoint issues)

### Threshold Configuration

Use multiple severity levels for better alert management:

```hcl
thresholds = {
  warning = {
    static = {
      value = 1000  # Early warning
    }
  }
  critical = {
    static = {
      value = 3000  # Critical threshold
    }
  }
}
```

### Alert Channel Routing

Route alerts based on severity:

```hcl
alert_channels = {
  warning  = [instana_alerting_channel_email.ops.id]
  critical = [
    instana_alerting_channel_pagerduty.oncall.id,
    instana_alerting_channel_slack.incidents.id
  ]
}
```

### Grace Period

Use grace periods to prevent alert flapping:

```hcl
grace_period = 300000  # 5 minutes - alert auto-closes if condition resolves
```

### Tag Filters

Use tag filters to scope alerts precisely:

```hcl
tag_filter = join(" AND ", [
  "call.type@na EQUALS 'HTTP'",
  "call.http.status@na GREATER_THAN 499",
  "entity.zone EQUALS 'production'"
])
```

### Time Thresholds

Choose appropriate time threshold types:

- `violations_in_sequence` - Alert on consecutive violations (sustained issues)
- `violations_in_period` - Alert on multiple violations within a window (intermittent issues)
- `request_impact` - Alert based on affected request volume (business impact)

## Notes

- The resource ID is auto-generated by Instana upon creation
- Multiple rules can be defined to monitor different metrics simultaneously
- Adaptive baseline thresholds learn from historical data and adjust automatically
- The `boundary_scope` determines which calls are included (inbound only, all, or default)
- Use `include_internal` and `include_synthetic` to control which call types are monitored
- The `granularity` setting affects how quickly alerts are triggered (smaller = faster detection, more sensitive)