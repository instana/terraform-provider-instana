# Custom Event Specification Resource

Configuration of custom event specifications for monitoring infrastructure and application metrics. Custom events allow you to define rules that trigger incidents based on various conditions like entity counts, thresholds, host availability, and more.

API Documentation: <https://instana.github.io/openapi/#operation/putCustomEventSpecification>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.

## Migration Guide (v5 to v6)

### Syntax Changes Overview

The main change is in how the `rules` block is defined. In v6, nested rule configurations use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_custom_event_specification" "example" {
  name = "threshold-alert"
  
  rules {  
    threshold {
      severity = "critical"
      metric_name = "cpu.usage"
      # ... other fields
    }
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_custom_event_specification" "example" {
  name = "threshold-alert"
  
  rules = {
    threshold = {
      severity = "critical"
      metric_name = "cpu.usage"
      # ... other fields
    }
  }
}
```


### Complete Migration Examples

#### Entity Count Rule Migration

**OLD (v5.x):**
```hcl
rules {  
  entity_count {
    severity = "warning"
    condition_operator = "="
    condition_value = 100
  }
}
```

**NEW (v6.x):**
```hcl
rules = {
  entity_count = {
    severity = "warning"
    condition_operator = "="
    condition_value = 100
  }
}
```

#### Threshold Rule Migration

**OLD (v5.x):**
```hcl
rules { 
  threshold {
    severity = "critical"
    metric_name = "nomad.client.allocations.pending"
    window = 60000
    aggregation = "avg"
    condition_operator = ">"
    condition_value = 0
  }
}
```

**NEW (v6.x):**
```hcl
rules = {
  threshold = {
    severity = "critical"
    metric_name = "nomad.client.allocations.pending"
    window = 60000
    aggregation = "avg"
    condition_operator = ">"
    condition_value = 0
  }
}
```

#### Host Availability Rule Migration

**OLD (v5.x):**
```hcl
rules {  
  host_availability {
    severity = "warning"
    offline_duration = 60000
    close_after = 120000
    tag_filter = "tag@na EQUALS 'foo_bar'"
  }
}
```

**NEW (v6.x):**
```hcl
rules = {
  host_availability = {
    severity = "warning"
    offline_duration = 60000
    close_after = 120000
    tag_filter = "tag@na EQUALS 'foo_bar'"
  }
}
```

## Example Usage

### Basic Entity Count Rule

Monitor the count of entities matching specific criteria:

```hcl
resource "instana_custom_event_specification" "entity_count" {
  name = "High Agent Count Alert"
  description = "Alert when agent count exceeds threshold"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "instanaAgent"

  rules = {
    entity_count = {
      severity = "warning"
      condition_operator = ">"
      condition_value = 100
    }
  }
}
```

### Entity Count Verification Rule

Verify entity counts with matching criteria:

```hcl
resource "instana_custom_event_specification" "entity_count_verification" {
  name = "Process Count Verification"
  description = "Verify process count on hosts"
  query = "entity.host.name:prod-*"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "host"

  rules = {
    entity_count_verification = {
      severity = "critical"
      matching_entity_type = "process"
      matching_operator = "is"
      matching_entity_label = "nginx"
      offline_duration = 60000
      condition_operator = "<"
      condition_value = 2
    }
  }
}
```

### Entity Verification Rule

Check for missing entities on hosts:

```hcl
resource "instana_custom_event_specification" "entity_verification" {
  name = "Missing Process Alert"
  description = "Alert when required process is missing"
  query = "entity.host.environment:production"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "host"

  rules = {
    entity_verification = {
      severity = "critical"
      matching_entity_type = "process"
      matching_operator = "contains"
      matching_entity_label = "critical-service"
      offline_duration = 300000
    }
  }
}
```

### Host Availability Rule - Scoped by Tag

Monitor host availability for specific tagged hosts:

```hcl
resource "instana_custom_event_specification" "host_availability_tagged" {
  name = "Production Host Offline"
  description = "Alert when production hosts go offline"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "host"

  rules = {
    host_availability = {
      severity = "critical"
      offline_duration = 60000
      close_after = 120000
      tag_filter = "host.environment@na EQUALS 'production'"
    }
  }
}
```

### Host Availability Rule - All Hosts

Monitor all hosts without tag filtering:

```hcl
resource "instana_custom_event_specification" "host_availability_all" {
  name = "Any Host Offline"
  description = "Alert when any host goes offline"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "host"

  rules = {
    host_availability = {
      severity = "warning"
      offline_duration = 120000
      close_after = 300000
    }
  }
}
```

### System Rule

Use built-in Instana system rules:

```hcl
resource "instana_custom_event_specification" "system_rule" {
  name = "System Rule Alert"
  description = "Alert based on Instana system rule"
  query = "entity.type:any"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "any"

  rules = {
    system = {
      severity = "critical"
      system_rule_id = "builtin-system-rule-id"
    }
  }
}
```

### Single Threshold Rule

Monitor a single metric threshold:

```hcl
resource "instana_custom_event_specification" "single_threshold" {
  name = "High CPU Usage"
  description = "Alert on high CPU usage"
  query = "entity.host.name:prod-*"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "host"

  rules = {
    threshold = {
      severity = "critical"
      metric_name = "cpu.usage"
      rollup = 60000
      window = 300000
      aggregation = "avg"
      condition_operator = ">"
      condition_value = 80
    }
  }
}
```

### Multiple Threshold Rules with OR Logic

Combine multiple threshold rules:

```hcl
resource "instana_custom_event_specification" "multiple_thresholds" {
  name = "Nomad Scheduler Issues"
  description = "Alert on blocked or pending allocations"
  query = "entity.type:nomadScheduler"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "nomadScheduler"
  rule_logical_operator = "OR"
  
  rules = {
    threshold = {
      severity = "critical"
      metric_name = "nomad.client.allocations.blocked"
      rollup = 60000
      window = 300000
      aggregation = "avg"
      condition_operator = ">"
      condition_value = 0
    }
  }
}

# Note: For multiple threshold rules, you need to create separate resources
# and combine them using alert configurations
resource "instana_custom_event_specification" "multiple_thresholds_2" {
  name = "Nomad Pending Allocations"
  description = "Alert on pending allocations"
  query = "entity.type:nomadScheduler"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "nomadScheduler"
  
  rules = {
    threshold = {
      severity = "critical"
      metric_name = "nomad.client.allocations.pending"
      rollup = 60000
      window = 300000
      aggregation = "avg"
      condition_operator = ">"
      condition_value = 0
    }
  }
}
```

### Threshold Rule with Metric Pattern

Use dynamic metric patterns for flexible matching:

```hcl
resource "instana_custom_event_specification" "metric_pattern" {
  name = "Dynamic Metric Alert"
  description = "Alert on metrics matching pattern"
  query = "entity.type:custom"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "custom"

  rules = {
    threshold = {
      severity = "warning"
      rollup = 60000
      window = 300000
      aggregation = "sum"
      condition_operator = ">"
      condition_value = 1000
      
      metric_pattern = {
        prefix = "custom.metric"
        postfix = ".count"
        placeholder = "*"
        operator = "contains"
      }
    }
  }
}
```

### Complex Tag Filter Example

Use advanced tag filtering with multiple conditions:

```hcl
resource "instana_custom_event_specification" "complex_filter" {
  name = "Complex Tag Filter Alert"
  description = "Alert with complex tag filtering"
  enabled = true
  triggering = true
  entity_type = "host"
  
  # Complex tag filter with AND/OR logic
  query = "(host.environment@na EQUALS 'production' OR host.environment@na EQUALS 'staging') AND host.region@na STARTS_WITH 'us-'"

  rules = {
    threshold = {
      severity = "critical"
      metric_name = "memory.usage"
      rollup = 60000
      window = 300000
      aggregation = "avg"
      condition_operator = ">"
      condition_value = 90
    }
  }
}
```

### Kubernetes-Specific Monitoring

Monitor Kubernetes resources:

```hcl
resource "instana_custom_event_specification" "k8s_pod_restarts" {
  name = "High Pod Restart Rate"
  description = "Alert on frequent pod restarts"
  enabled = true
  triggering = true
  expiration_time = 300000
  entity_type = "kubernetesDeployment"
  query = "kubernetes.namespace.name:production"

  rules = {
    threshold = {
      severity = "warning"
      metric_name = "kubernetes.pod.restarts"
      rollup = 300000
      window = 600000
      aggregation = "sum"
      condition_operator = ">"
      condition_value = 5
    }
  }
}
```

### Database Performance Monitoring

Monitor database metrics:

```hcl
resource "instana_custom_event_specification" "db_slow_queries" {
  name = "Slow Database Queries"
  description = "Alert on slow query execution"
  enabled = true
  triggering = true
  expiration_time = 60000
  entity_type = "postgresql"
  query = "postgresql.database.name:production_db"

  rules = {
    threshold = {
      severity = "critical"
      metric_name = "postgresql.query.duration"
      rollup = 60000
      window = 300000
      aggregation = "p95"
      condition_operator = ">"
      condition_value = 5000
    }
  }
}
```

### Multi-Environment Setup

Create alerts for different environments:

```hcl
locals {
  environments = {
    production = {
      cpu_threshold = 80
      memory_threshold = 85
      severity = "critical"
    }
    staging = {
      cpu_threshold = 90
      memory_threshold = 90
      severity = "warning"
    }
  }
}

resource "instana_custom_event_specification" "env_cpu_alert" {
  for_each = local.environments

  name = "${each.key} High CPU"
  description = "CPU alert for ${each.key} environment"
  enabled = true
  triggering = true
  entity_type = "host"
  query = "host.environment@na EQUALS '${each.key}'"

  rules = {
    threshold = {
      severity = each.value.severity
      metric_name = "cpu.usage"
      rollup = 60000
      window = 300000
      aggregation = "avg"
      condition_operator = ">"
      condition_value = each.value.cpu_threshold
    }
  }
}
```

### Disabled Event for Testing

Create a disabled event specification for testing:

```hcl
resource "instana_custom_event_specification" "test_event" {
  name = "Test Event (Disabled)"
  description = "Test event specification - not triggering"
  enabled = false
  triggering = false
  entity_type = "host"

  rules = {
    threshold = {
      severity = "warning"
      metric_name = "cpu.usage"
      rollup = 60000
      window = 300000
      aggregation = "avg"
      condition_operator = ">"
      condition_value = 95
    }
  }
}
```

## Argument Reference

* `name` - Required - The name of the custom event specification
* `entity_type` - Required - The entity type/plugin for which the rule will be defined. Must be set to:
  * `any` for [System Rules](#system-rule)
  * `host` for [Entity Verification Rules](#entity-verification-rule), [Entity Count Verification Rules](#entity-count-verification-rule), and [Host Availability Rules](#host-availability-rule)
  * `instanaAgent` for [Entity Count Rules](#entity-count-rule)
  * For threshold rules, supported entity types (plugins) can be retrieved from the Instana REST API using `/api/infrastructure-monitoring/catalog/plugins`
* `description` - Optional - The description text of the custom event specification (default: empty string)
* `query` - Optional - The dynamic filter query for which the rule should be applied to (default: empty string)
* `enabled` - Optional - Boolean flag if the rule should be enabled (default: true)
* `triggering` - Optional - Boolean flag if the rule should trigger an incident (default: false)
* `expiration_time` - Optional - The grace period in milliseconds until the issue is closed
* `rule_logical_operator` - Optional - The logical operator which will be applied to combine multiple rules (threshold rules only). Default: `AND`. Allowed values: `AND`, `OR`
* `rules` - Required - The configuration of the specific rule of the custom event [Details](#rules)

### Rules

The `rules` attribute is a single nested object containing exactly one of the following rule types:

* `entity_count` - Optional - Configuration of entity count rules [Details](#entity-count-rule)
* `entity_count_verification` - Optional - Configuration of entity count verification rules [Details](#entity-count-verification-rule)
* `entity_verification` - Optional - Configuration of entity verification rules [Details](#entity-verification-rule)
* `host_availability` - Optional - Configuration of host availability rules [Details](#host-availability-rule)
* `system` - Optional - Configuration of system rules [Details](#system-rule)
* `threshold` - Optional - Configuration of threshold rules [Details](#threshold-rule)

#### Entity Count Rule

Monitor the count of entities matching specific criteria.

* `severity` - Required - The severity of the rule. Allowed values: `warning`, `critical`
* `condition_operator` - Required - The condition operator used to check against the calculated metric value. Supported values: `=`, `!=`, `<=`, `<`, `>`, `>=`
* `condition_value` - Required - The numeric condition value used to check against the calculated metric value

#### Entity Count Verification Rule

Verify that a specific number of matching entities exist on selected hosts.

* `severity` - Required - The severity of the rule. Allowed values: `warning`, `critical`
* `condition_operator` - Required - The condition operator used to check against the calculated metric value. Supported values: `=`, `!=`, `<=`, `<`, `>`, `>=`
* `condition_value` - Required - The numeric condition value used to check against the calculated metric value
* `matching_entity_type` - Required - The entity type used to check for matching entities on the selected hosts. Supported entity types (plugins) can be retrieved from the Instana REST API using `/api/infrastructure-monitoring/catalog/plugins`
* `matching_operator` - Required - The comparison operator used to check for matching entities. Allowed values: `is`, `contains`, `startsWith`, `endsWith`
* `matching_entity_label` - Required - The label/string to check for matching entities on the selected hosts

#### Entity Verification Rule

Verify that specific entities are running on selected hosts.

* `severity` - Required - The severity of the rule. Allowed values: `warning`, `critical`
* `matching_entity_type` - Required - The entity type used to check for matching entities on the selected hosts. Supported entity types (plugins) can be retrieved from the Instana REST API using `/api/infrastructure-monitoring/catalog/plugins`
* `matching_operator` - Required - The comparison operator used to check for matching entities. Allowed values: `is`, `contains`, `startsWith`, `endsWith`
* `matching_entity_label` - Required - The label/string to check for matching entities on the selected hosts
* `offline_duration` - Required - The duration in milliseconds to wait until the entity is considered as offline

#### Host Availability Rule

Monitor host availability and trigger alerts when hosts go offline.

* `severity` - Required - The severity of the rule. Allowed values: `warning`, `critical`
* `offline_duration` - Required - The duration in milliseconds to wait until the host is considered as offline
* `close_after` - Optional - If a host is offline for longer than the defined period, Instana does not expect the host to reappear anymore, and the event will be closed after the grace period
* `tag_filter` - Optional - Tag filter expression to scope which hosts to monitor. Only `tag` is allowed for the tag filter. Example: `tag:my_tag EQUALS 'test'`. If not specified, all hosts are monitored

#### System Rule

Use built-in Instana system rules for monitoring.

* `severity` - Required - The severity of the rule. Allowed values: `warning`, `critical`
* `system_rule_id` - Required - The ID of the Instana system rule

#### Threshold Rule

Monitor metrics against defined thresholds.

* `severity` - Required - The severity of the rule. Allowed values: `warning`, `critical`
* `metric_name` - Optional (required for Built-In and Custom Metrics only) - The name of the built-in or custom metric. Supported built-in metrics can be retrieved from the REST API using `/api/infrastructure-monitoring/catalog/metrics/{plugin}`. Conflicts with `metric_pattern` - exactly one must be provided
* `metric_pattern` - Optional (required for Dynamic Built-In Metrics only) - The metric pattern of the dynamic built-in metric [Details](#metric-pattern). Conflicts with `metric_name` - exactly one must be provided
* `rollup` - Required - The resolution of the monitored metrics in milliseconds
* `window` - Required - The time window in milliseconds within which the rule condition is applied
* `aggregation` - Required - The aggregation used to calculate the metric value for the given time window and/or rollup. Supported values: `sum`, `avg`, `min`, `max`
* `condition_operator` - Required - The condition operator used to check against the calculated metric value. Supported values: `=`, `!=`, `<=`, `<`, `>`, `>=`
* `condition_value` - Required - The numeric condition value used to check against the calculated metric value

##### Metric Pattern

For dynamic built-in metrics, define a pattern to match multiple metrics.

* `prefix` - Required - The prefix of the built-in dynamic metric
* `postfix` - Optional - The postfix of the built-in dynamic metric (default: empty string)
* `placeholder` - Optional - The placeholder string of the dynamic metric (default: empty string)
* `operator` - Optional - The operation used to check for matching placeholder string. Allowed values: `is`, `contains`, `any`, `startsWith`, `endsWith`. Default: `EQUALS`

## Attributes Reference

* `id` - The ID of the custom event specification

## Import

Custom event specifications can be imported using the `id`, e.g.:

```bash
$ terraform import instana_custom_event_specification.example 60845e4e5e6b9cf8fc2868da
```

## Notes

* The `default_name_prefix` and `default_name_suffix` features are **NOT** supported for this resource as they will be removed in future versions
* For threshold rules, you can combine up to 5 rules using the `rule_logical_operator`
* Tag filters support complex expressions with AND/OR logic and various comparison operators
* Entity types and available metrics can be discovered using the Instana REST API
* When `triggering` is set to `false`, the event will be created but won't trigger incidents
* The `expiration_time` defines how long an issue remains open after conditions are no longer met
