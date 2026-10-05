# Infrastructure Alert Configuration Resource

Management of Infrastructure alert configurations (Infrastructure Smart Alerts). These alerts monitor infrastructure metrics and trigger notifications based on defined thresholds and conditions.

API Documentation: <https://instana.github.io/openapi/#tag/Infrastructure-Alert-Configuration>


## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how nested blocks are defined. In v6, all nested configurations use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_infra_alert_config" "example" {
  name = "test-alert"
  
  alert_channels {
    warning = ["channel-1"]
    critical = ["channel-2"]
  }
  
  rules {
    generic_rule {
      aggregation              = "MIN"
      cross_series_aggregation = "MIN"
      entity_type              = "kubernetesPod"
      metric_name              = "cpuUsageToLimitRatio"
      threshold {
        critical {
          static {
            value = 90
          }
        }
      }
    }
  }
  
  time_threshold {
    violations_in_sequence {
      time_window = 600000
    }
  }
  tag_filter = "kubernetes.namespace.name@na EQUALS 'otel'"

  custom_payload_field {
    key = "env"
    value = "prod"
  }
  
  custom_payload_field {
    key = "region"
    value = "us-east"
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_infra_alert_config" "example" {
  name = "test-alert"
  description = "Description"

  evaluation_type      = "PER_ENTITY"
  granularity          = 600000
  alert_channels = {
    warning = ["channel-1"]
    critical = ["channel-2"]
  }
  tag_filter = "kubernetes.namespace.name@na EQUALS 'otel'"

  rules = [
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesHorizontalPodAutoscaler"
        metric_name              = "maxReplicas"
        threshold = {
          critical = {
            static = {
              value = 90
            }
          }
        }
        threshold_operator = ">="
      }
    }
  ]
  
  time_threshold = {
    violations_in_sequence = {
      time_window = 600000
    }
  }
  
  custom_payload_field = [
    {
      key = "env"
      value = "prod"
    },
    {
      key = "region"
      value = "us-east"
    }
  ]
}
```

### Key Syntax Changes

1. **Alert Channels**: `alert_channels { }` → `alert_channels = { }`
2. **Rules**: `rules { }` → `rules = [{ }]` (now a **list** of rule objects)
3. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
4. **Custom Payload Fields**: Multiple `custom_payload_field { }` blocks → Single `custom_payload_field = [{ }, { }]` list
5. **All nested objects**: Use `= { }` syntax

## Example Usage

### CPU Alert with Static Thresholds

```hcl
resource "instana_infra_alert_config" "cpu_alert" {
  name = "High CPU Usage Alert - $${severity}" # Use double $$ to define placeholders
  description = "Alert when CPU usage exceeds thresholds"
  alert_channels = {
    critical = ["channel-id-1"]
    warning  = ["channel-id-1"]
  }
  evaluation_type      = "CUSTOM"
  granularity          = 60000
  triggering           = true
  rules = [
    {
      generic_rule = {
        aggregation              = "MIN"
        cross_series_aggregation = "MIN"
        entity_type              = "kubernetesPod"
        metric_name              = "cpuUsageToLimitRatio"
        threshold = {
          warning = {
            static = {
              value = 1
            }
          }
        }
        threshold_operator = ">="
      }
    }
  ]
  tag_filter = "kubernetes.namespace.name@na EQUALS 'otel-demo'"
  time_threshold = {
    violations_in_sequence = {
      time_window = 60000
    }
  }
}
```

### Multiple Rules with Rule Logical Operator

```hcl
resource "instana_infra_alert_config" "multi_rule_alert" {
  name            = "Multi-Rule Infrastructure Alert"
  description     = "Alert combining CPU and memory rules with OR logic"
  evaluation_type = "PER_ENTITY"
  granularity     = 600000

  rule_logical_operator = "OR"

  rules = [
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesPod"
        metric_name              = "cpuUsageToLimitRatio"
        threshold_operator       = ">="
        threshold = {
          critical = {
            static = {
              value = 90
            }
          }
        }
      }
    },
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesPod"
        metric_name              = "memoryUsageToLimitRatio"
        threshold_operator       = ">="
        threshold = {
          critical = {
            static = {
              value = 85
            }
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
}
```

### Alert with BETWEEN Threshold Operator

```hcl
resource "instana_infra_alert_config" "range_alert" {
  name            = "CPU Usage Out-of-Range Alert"
  description     = "Alert when CPU usage falls outside normal operating range"
  evaluation_type = "PER_ENTITY"
  granularity     = 600000

  rules = [
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesPod"
        metric_name              = "cpuUsageToLimitRatio"
        threshold_operator       = "OUTSIDE"
        threshold = {
          warning = {
            static = {
              lower_bound = 10
              upper_bound = 80
            }
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
}
```

### Alert with Metric Group By and Metric Tag Filter

```hcl
resource "instana_infra_alert_config" "grouped_alert" {
  name            = "Grouped Metric Alert"
  description     = "Alert with metric grouping and scoped metric tag filter"
  evaluation_type = "CUSTOM"
  granularity     = 600000

  rules = [
    {
      generic_rule = {
        aggregation              = "MEAN"
        cross_series_aggregation = "MEAN"
        entity_type              = "kubernetesPod"
        metric_name              = "cpuUsageToLimitRatio"
        threshold_operator       = ">="
        metric_group_by          = ["namespace", "pod"]
        metric_tag_filter_expression = {
          logical_operator = "AND"
        }
        threshold = {
          critical = {
            static = {
              value = 90
            }
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
}
```

## Generating Configuration from Existing Resources

If you have already created an infrastructure alert configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your infrastructure alert configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_infra_alert_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual infrastructure alert configuration ID from Instana

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

* `name` - Required - The name for the infrastructure alert configuration
* `description` - Required - The description text of the infrastructure alert config
* `granularity` - Required - The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used. Allowed values: `300000`, `600000`, `900000`, `1200000`, `1800000` (milliseconds)
* `grace_period` - Optional - The duration in milliseconds for which an alert remains open after conditions are no longer violated. The alert automatically closes once the grace period expires.
* `evaluation_type` - Required - The evaluation type of the infrastructure alert config. Allowed values:
  * `CUSTOM` - Combine all metrics in scope into a single metric per group (default)
  * `PER_ENTITY` - Monitor each metric individually and trigger alerts for each individual entity
* `triggering` - Optional - Indicates whether the alert should be triggered. Default: `false`
* `alert_channels` - Optional - Set of alert channel IDs associated with the severity [Details](#alert-channels-reference)
* `group_by` - Optional - List of grouping tags used to group the metric results
* `tag_filter` - Optional - The tag filter of the infrastructure alert config [Details](#tag-filter-argument-reference)
* `rule_logical_operator` - Optional - Logical operator used to combine multiple rules. Allowed values: `AND`, `OR`. Default: `AND`
* `rules` - Optional - A list of rule configurations (1–5 entries). Each entry contains a single `generic_rule` [Details](#rules-argument-reference)
* `time_threshold` - Required - Indicates the type of violation of the defined threshold [Details](#time-threshold-argument-reference)
* `custom_payload_field` - Optional - A list of custom payload fields (static key/value pairs or dynamic tag values added to the event) [Details](#custom-payload-field-argument-reference)

### Alert Channels Reference

* `warning` - Optional - List of alert channel IDs associated with the warning severity 
* `critical` - Optional - List of alert channel IDs associated with the critical severity

### Tag Filter Argument Reference

The **tag_filter** defines which entities should be included in the alert scope. It supports:

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

### Understanding the name, key, origin Structure

A tag filter string is a compact DSL that maps to **5 distinct fields** in the Instana API payload. Understanding this split is essential to writing filter expressions that the API will correctly recognise.

```
otel.attribute:'aws.rds.cluster.id'@na EQUALS 'my-cluster'
│──────────────│─────────────────────│────│──────│──────────│
     name              key           origin  op     value
```

| Part | Delimiter | API Field | Example |
|---|---|---|---|
| `otel.attribute` | everything before `:` | `name` | `"otel.attribute"` |
| `aws.rds.cluster.id` | after `:` | `key` | `"aws.rds.cluster.id"` |
| `na` | after `@` | `entity` | `"NOT_APPLICABLE"` |
| `EQUALS` | space-separated keyword | `operator` | `"EQUALS"` |
| `my-cluster` | quoted string | `value` / `stringValue` | `"my-cluster"` |

**Not all filters have a `key`.** For standard Instana tags like `entity.name` or `call.http.method`, the dots are simply part of the tag name and the entire dotted string maps to the `name` field alone — no `:` is needed:

```hcl
tag_filter = "entity.name@dest EQUALS 'payments'"
# API receives: name="entity.name", key=null
```

**When the API requires a `key`**, you must use the `:` separator to split the tag name from the key. It applies to any tag type that carries a dynamic key portion:

```hcl
# OTel attribute tag — name=otel.attribute, key=aws.rds.cluster.id
tag_filter = "otel.attribute:'aws.rds.cluster.id'@na EQUALS 'my-cluster'"

# Website beacon metadata — name=beacon.meta, key=stage
tag_filter = "beacon.meta:'stage'@na EQUALS 'production'"
```

**The `@` suffix** specifies the entity origin. If omitted it defaults to `@dest` (`DESTINATION`). Always verify the origin required by the specific tag type.

| Suffix | API value | When to use |
|---|---|---|
| `@src` | `SOURCE` | Source-side entity attributes |
| `@dest` | `DESTINATION` | Default — destination-side attributes |
| `@na` | `NOT_APPLICABLE` | Tags not tied to call direction (OTel, beacon meta, agent tags) |

> **Tip:** If you are unsure about the correct structure for a particular filter, configure it once through the Instana UI, then use Terraform's `import` block to generate the equivalent Terraform configuration. The imported output will show you exactly how the filter string should be written for that tag type — including whether `:` and `@` are needed.

### Rules Argument Reference

`rules` is a **list** of rule objects with a minimum of 1 and a maximum of 5 entries. Multiple rules are combined using the top-level `rule_logical_operator` (`AND` by default).

Each list item contains:

* `generic_rule` - Optional - A generic rule based on a custom aggregated metric [Details](#generic-rule-argument-reference)

#### Generic Rule Argument Reference 

* `metric_name` - Required - The metric name of the infrastructure alert rule (refer to [/api/infrastructure-monitoring/catalog/metrics/{plugin}](https://developer.ibm.com/apis/catalog/instana--instana-rest-api/api/API--instana--instana-rest-api-documentation#getInfrastructureCatalogMetrics) for valid metric names per plugin, e.g. `host`, `instanaAgent`)
* `entity_type` - Required - The entity type of the infrastructure alert rule
* `aggregation` - Required - The aggregation function of the infra alert rule. Supported values: `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `SUM`, `PER_SECOND`
* `cross_series_aggregation` - Required - Cross-series aggregation function of the infra alert rule. Supported values: `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `SUM`
* `regex` - Optional/Computed - Boolean indicating if the given metric name follows a regex pattern
* `threshold_operator` - Optional/Computed - The operator which will be applied to evaluate the threshold. Supported values: `>`, `>=`, `<`, `<=`, `OUTSIDE`, `BETWEEN`
  * Use `>`, `>=`, `<`, `<=` with `static.value`
  * Use `BETWEEN` or `OUTSIDE` with `static.lower_bound` and `static.upper_bound`
* `threshold` - Optional - Indicates the type of threshold associated with the given severity [Details](#threshold-rule-argument-reference)
* `metric_group_by` - Optional - A set of metric tag names used to group the metric within this rule
* `metric_tag_filter_expression` - Optional/Computed - A metric-scoped tag filter expression to narrow which metric time series are evaluated by this rule [Details](#metric-tag-filter-expression-argument-reference)

#### Threshold Rule Argument Reference

At least one of the elements below must be configured:

* `warning` - Optional - Threshold associated with the warning severity [Details](#threshold-argument-reference)
* `critical` - Optional - Threshold associated with the critical severity [Details](#threshold-argument-reference)

##### Threshold Argument Reference

* `static` - Optional - Static threshold definition [Details](#static-threshold-argument-reference)
* `adaptive_baseline` - Optional - Adaptive baseline threshold definition [Details](#adaptive-baseline-threshold-argument-reference)

###### Static Threshold Argument Reference

* `value` - Optional - The point threshold value. Required for `>`, `>=`, `<`, `<=` operators. Must be omitted for `BETWEEN`/`OUTSIDE` operators.
* `lower_bound` - Optional - The lower bound of the threshold range. Required for `BETWEEN` and `OUTSIDE` operators. Must be omitted for point-value operators.
* `upper_bound` - Optional - The upper bound of the threshold range. Required for `BETWEEN` and `OUTSIDE` operators. Must be omitted for point-value operators.

###### Adaptive Baseline Threshold Argument Reference

* `deviation_factor` - Optional - The deviation factor for the adaptive baseline threshold
* `adaptability` - Optional - The adaptability for the adaptive baseline threshold
* `seasonality` - Optional - The seasonality for the adaptive baseline threshold

#### Metric Tag Filter Expression Argument Reference

* `logical_operator` - Optional/Computed - The logical operator for combining filter elements within the metric tag filter. Supported values: `AND`, `OR`. Default: `AND`

### Time Threshold Argument Reference

* `violations_in_sequence` - Required - Time threshold based on violations in sequence [Details](#violations-in-sequence-time-threshold-argument-reference)

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Optional/Computed - The time window of the time threshold in milliseconds. Default: `600000`

### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Optional - The static string value of the custom payload field. Either `value` or `dynamic_value` must be defined
* `dynamic_value` - Optional - The dynamic value of the custom payload field [Details](#dynamic-custom-payload-field-value). Either `value` or `dynamic_value` must be defined

#### Dynamic Custom Payload Field Value

* `key` - Optional - The key of the tag which should be added to the payload
* `tag_name` - Required - The name of the tag which should be added to the payload

## Attributes Reference

* `id` - The ID of the infrastructure alert configuration

## Import

Infrastructure alert configurations can be imported using the `id`, e.g.:

```bash
$ terraform import instana_infra_alert_config.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Use `CUSTOM` evaluation type to aggregate metrics across entities
* Use `PER_ENTITY` evaluation type to monitor each entity individually
* `rules` is a list supporting 1–5 rule entries; use `rule_logical_operator` to control how multiple rules are combined (`AND` by default)
* `metric_group_by` groups the metric time series within a rule before aggregation
* `metric_tag_filter_expression` scopes a rule to specific metric tags; defaults to an empty `AND` expression when omitted
* Regex patterns in `metric_name` allow monitoring multiple related metrics with a single rule
* For range-based threshold operators (`BETWEEN`, `OUTSIDE`), use `lower_bound` and `upper_bound` instead of `value` in the `static` block
* Tag filters support complex expressions for precise entity selection
* Custom payload fields can include both static values and dynamic tag values
* The `granularity` determines how frequently the alert condition is evaluated
* Time threshold defines how many consecutive violations are required before triggering an alert
