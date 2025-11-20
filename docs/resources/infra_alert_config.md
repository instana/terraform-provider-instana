# Infrastructure Alert Configuration Resource

Management of Infrastructure alert configurations (Infrastructure Smart Alerts). These alerts monitor infrastructure metrics and trigger notifications based on defined thresholds and conditions.

API Documentation: <https://instana.github.io/openapi/#tag/Infrastructure-Alert-Configuration>


## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.


## Migration Guide (v5 to v6)

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
      metric_name = "cpu.usage"
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
  
  alert_channels = {
    warning = ["channel-1"]
    critical = ["channel-2"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu.usage"
      threshold = {
        critical = {
          static = {
            value = 90
          }
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
2. **Rules**: `rules { }` → `rules = { }`
3. **Time Threshold**: `time_threshold { }` → `time_threshold = { }`
4. **Custom Payload Fields**: Multiple `custom_payload_field { }` blocks → Single `custom_payload_field = [{ }, { }]` list
5. **All nested objects**: Use `= { }` syntax

## Example Usage

### Basic CPU Alert with Static Thresholds

```hcl
resource "instana_infra_alert_config" "cpu_alert" {
  name = "High CPU Usage Alert"
  description = "Alert when CPU usage exceeds thresholds"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  alert_channels = {
    warning = ["warning-channel-id"]
    critical = ["critical-channel-id"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu.usage"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "SUM"
      regex = false
      threshold_operator = ">="
      
      threshold = {
        critical = {
          static = {
            value = 90.0
          }
        }
        warning = {
          static = {
            value = 75.0
          }
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

### Memory Alert with Tag Filter

```hcl
resource "instana_infra_alert_config" "memory_alert" {
  name = "High Memory Usage"
  description = "Alert on high memory usage for production hosts"
  granularity = 600000
  evaluation_type = "CUSTOM"
  tag_filter = "host.environment@na EQUALS 'production'"
  
  alert_channels = {
    critical = ["ops-team-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "memory.used"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 8589934592  # 8GB in bytes
          }
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

### Regex Pattern Matching for Multiple Metrics

```hcl
resource "instana_infra_alert_config" "cpu_components" {
  name = "CPU Component Alert"
  description = "Monitor multiple CPU metrics using regex"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  alert_channels = {
    warning = ["monitoring-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu\\.(nice|user|sys|wait)"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "SUM"
      regex = true
      threshold_operator = ">="
      
      threshold = {
        warning = {
          static = {
            value = 80.0
          }
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

### Per-Entity Evaluation

```hcl
resource "instana_infra_alert_config" "per_entity_disk" {
  name = "Disk Usage Per Host"
  description = "Monitor disk usage for each host individually"
  granularity = 600000
  evaluation_type = "PER_ENTITY"
  
  alert_channels = {
    critical = ["storage-team-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "disk.used.percent"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "MEAN"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 90.0
          }
        }
        warning = {
          static = {
            value = 80.0
          }
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

### Alert with Group By

```hcl
resource "instana_infra_alert_config" "grouped_alert" {
  name = "CPU Alert Grouped by Region"
  description = "Monitor CPU grouped by region tag"
  granularity = 600000
  evaluation_type = "CUSTOM"
  group_by = ["host.region", "host.environment"]
  
  alert_channels = {
    critical = ["regional-ops-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu.usage"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 85.0
          }
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

### Alert with Custom Payload Fields

```hcl
resource "instana_infra_alert_config" "with_custom_payload" {
  name = "Alert with Custom Context"
  description = "Alert with additional context in notifications"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  alert_channels = {
    critical = ["enriched-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu.usage"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "SUM"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 90.0
          }
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
      value = "https://wiki.example.com/runbooks/high-cpu"
    }
  ]
}
```

### Alert with Dynamic Custom Payload

```hcl
resource "instana_infra_alert_config" "dynamic_payload" {
  name = "Alert with Dynamic Tags"
  description = "Include dynamic tag values in alert payload"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  alert_channels = {
    critical = ["ops-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "memory.used.percent"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 90.0
          }
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
      key = "static_field"
      value = "static_value"
    },
    {
      key = "host_name"
      dynamic_value = {
        key = "name"
        tag_name = "host.name"
      }
    },
    {
      key = "region"
      dynamic_value = {
        tag_name = "host.region"
      }
    }
  ]
}
```

### Complex Tag Filter with Multiple Conditions

```hcl
resource "instana_infra_alert_config" "complex_filter" {
  name = "Production Critical Services"
  description = "Monitor critical services in production"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  # Complex tag filter with AND/OR logic
  tag_filter = "(host.environment@na EQUALS 'production' AND host.tier@na EQUALS 'critical') OR host.name@na STARTS_WITH 'prod-critical-'"
  
  alert_channels = {
    critical = ["critical-ops-channel", "pagerduty-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu.usage"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 80.0
          }
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

### Network Metrics Alert

```hcl
resource "instana_infra_alert_config" "network_traffic" {
  name = "High Network Traffic"
  description = "Alert on high network throughput"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  alert_channels = {
    warning = ["network-team-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "network.tx.bytes"
      entity_type = "host"
      aggregation = "SUM"
      cross_series_aggregation = "SUM"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        warning = {
          static = {
            value = 1073741824  # 1GB
          }
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

### Different Aggregation Methods

```hcl
resource "instana_infra_alert_config" "p95_latency" {
  name = "P95 Latency Alert"
  description = "Monitor 95th percentile latency"
  granularity = 600000
  evaluation_type = "CUSTOM"
  
  alert_channels = {
    critical = ["performance-team"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "response.time"
      entity_type = "service"
      aggregation = "P95"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 1000.0  # 1 second
          }
        }
        warning = {
          static = {
            value = 500.0  # 500ms
          }
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

### Multi-Environment Setup

```hcl
locals {
  environments = {
    production = {
      cpu_threshold = 80
      memory_threshold = 85
      granularity = 300000
      channels = ["prod-ops", "pagerduty"]
    }
    staging = {
      cpu_threshold = 90
      memory_threshold = 90
      granularity = 600000
      channels = ["staging-ops"]
    }
  }
}

resource "instana_infra_alert_config" "env_cpu_alert" {
  for_each = local.environments

  name = "${each.key} CPU Alert"
  description = "CPU monitoring for ${each.key} environment"
  granularity = each.value.granularity
  evaluation_type = "CUSTOM"
  tag_filter = "host.environment@na EQUALS '${each.key}'"
  
  alert_channels = {
    critical = each.value.channels
  }
  
  rules = {
    generic_rule = {
      metric_name = "cpu.usage"
      entity_type = "host"
      aggregation = "MEAN"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = each.value.cpu_threshold
          }
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

### Kubernetes Node Monitoring

```hcl
resource "instana_infra_alert_config" "k8s_node_pressure" {
  name = "Kubernetes Node Pressure"
  description = "Alert on node resource pressure"
  granularity = 300000
  evaluation_type = "PER_ENTITY"
  tag_filter = "kubernetes.node.name@na NOT_EMPTY"
  
  alert_channels = {
    critical = ["k8s-ops-channel"]
  }
  
  rules = {
    generic_rule = {
      metric_name = "kubernetes.node.memory.pressure"
      entity_type = "kubernetesNode"
      aggregation = "MAX"
      cross_series_aggregation = "MAX"
      regex = false
      threshold_operator = ">"
      
      threshold = {
        critical = {
          static = {
            value = 1.0
          }
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
      key = "cluster"
      dynamic_value = {
        tag_name = "kubernetes.cluster.name"
      }
    },
    {
      key = "node"
      dynamic_value = {
        tag_name = "kubernetes.node.name"
      }
    }
  ]
}
```

## Argument Reference

* `name` - Required - The name for the infrastructure alert configuration
* `description` - Required - The description text of the infrastructure alert config
* `granularity` - Required - The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used. Allowed values: `300000`, `600000`, `900000`, `1200000`, `1800000` (milliseconds)
* `evaluation_type` - Required - The evaluation type of the infrastructure alert config. Allowed values:
  * `CUSTOM` - Combine all metrics in scope into a single metric per group (default)
  * `PER_ENTITY` - Monitor each metric individually and trigger alerts for each individual entity
* `alert_channels` - Optional - Set of alert channel IDs associated with the severity [Details](#alert-channels-reference)
* `group_by` - Optional - List of grouping tags used to group the metric results
* `tag_filter` - Optional - The tag filter of the infrastructure alert config [Details](#tag-filter-argument-reference)
* `rules` - Required - A rule configuration with thresholds and their corresponding severity levels [Details](#rules-argument-reference)
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

### Rules Argument Reference

* `generic_rule` - Required - A generic rule based on custom aggregated metric [Details](#generic-rule-argument-reference)

#### Generic Rule Argument Reference 

* `metric_name` - Required - The metric name of the infrastructure alert rule
* `entity_type` - Required - The entity type of the infrastructure alert rule
* `aggregation` - Required - The aggregation function of the infra alert rule. Supported values: `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `SUM`, `PER_SECOND`
* `cross_series_aggregation` - Required - Cross-series aggregation function of the infra alert rule. Supported values: `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `SUM`
* `regex` - Required - Boolean indicating if the given metric name follows regex pattern or not
* `threshold_operator` - Required - The operator which will be applied to evaluate the threshold. Supported values: `>`, `>=`, `<`, `<=`
* `threshold` - Required - Indicates the type of threshold associated with given severity this alert rule is evaluated on [Details](#threshold-rule-argument-reference)

#### Threshold Rule Argument Reference

At least one of the elements below must be configured:

* `warning` - Optional - Threshold associated with the warning severity [Details](#threshold-argument-reference)
* `critical` - Optional - Threshold associated with the critical severity [Details](#threshold-argument-reference)

##### Threshold Argument Reference

* `static` - Required - Static threshold definition [Details](#static-threshold-argument-reference)

###### Static Threshold Argument Reference

* `value` - Required - The value of the static threshold

### Time Threshold Argument Reference

* `violations_in_sequence` - Required - Time threshold based on violations in sequence [Details](#violations-in-sequence-time-threshold-argument-reference)

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Required - The time window of the time threshold in milliseconds

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
* Regex patterns in `metric_name` allow monitoring multiple related metrics with a single rule
* Tag filters support complex expressions for precise entity selection
* Custom payload fields can include both static values and dynamic tag values
* The `granularity` determines how frequently the alert condition is evaluated
* Time threshold defines how many consecutive violations are required before triggering an alert
