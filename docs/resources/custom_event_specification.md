# Custom Event Specification Resource

Configuration of custom event specifications for monitoring infrastructure and application metrics. Custom events allow you to define rules that trigger incidents based on various conditions like entity counts, thresholds, host availability, and more.

API Documentation: <https://instana.github.io/openapi/#operation/putCustomEventSpecification>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

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
  enabled               = true
  entity_type           = "host"
  expiration_time       = 600000
  query                 = "entity.application.id:\"application-id\""
  rule_logical_operator = "AND"
  rules = {
    entity_verification = {
      matching_entity_label = "CustomQueMg"
      matching_entity_type  = "ibmMqQueueManager"
      matching_operator     = "contains"
      offline_duration      = 60000
      severity              = "warning"
    }
  }
}
```

## Generating Configuration from Existing Resources

If you have already created a custom event specification in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your custom event specification in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_custom_event_specification.custom_event_type
  id = "resource_id"
}
```

Replace:
- `custom_event_type` with your desired terraform block name
- `resource_id` with your actual custom event specification ID from Instana

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
