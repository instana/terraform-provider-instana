# Alerting Configuration

Management of alert configurations. Alert configurations define how either event types or 
event (aka rules) are reported to integrated services (Alerting Channels).

API Documentation: <https://instana.github.io/openapi/#operation/putAlert>

---
## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- Nested blocks now use **list/set attribute syntax** with `= [{ }]` instead of block syntax `{ }`
- All nested configurations require the equals sign
- Enhanced validation and better error messages
- Improved state management with computed fields

#### OLD (v5.x) Syntax:

```hcl
resource "instana_alerting_config" "example" {
  alert_name = "name"
  
  custom_payload_field {
    key   = "test"
    value = "test123"
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_alerting_config" "example" {
  alert_name = "name"
  
  custom_payload_field = [{
    key   = "test"
    value = "test123"
  }]
  #rest of the configuration
}
```
Please update your Terraform configurations to use the new attribute-based syntax.
---


- Nested blocks now use **list/set attribute syntax** with `= [{ }]` instead of block syntax `{ }`
- All nested configurations require the equals sign
- Enhanced validation and better error messages
- Improved state management with computed fields

### Basic Configuration with Rule IDs

```hcl
resource "instana_alerting_config" "alert_config {
  alert_name                         = "New Alert Configuration AN test"
  event_filter_query                 = "event.type:issue AND event.severity:critical entity.zone:\"helmrefactoring\""
  event_filter_rule_ids              = ["rule-id1", "rule-id2"]
  integration_ids                    = []
}
```

### Multi-Channel Configuration

```hcl
resource "instana_alerting_config" "alert_config_multi_channel" {
  alert_name = "Critical System Alerts"
  custom_payload_field = [
    {
      key           = var.alert_config_key
      value         = var.alert_config_value
    },
  ]
  event_filter_event_types           = ["critical", "incident", "warning"]
  integration_ids                    = [instana_alerting_channel_email.ops.id, instana_alerting_channel_slack.alerts.id]
}
```

### Infrastructure Monitoring with Custom Payload

```hcl
resource "instana_alerting_config" "infra_monitoring" {
  alert_name      = "Infrastructure Health"
  integration_ids = [instana_alerting_channel_opsgenie.infra.id]
  
  event_filter_query       = "entity.zone:us-east-1 AND (entity.type:host OR entity.type:container)"
  event_filter_event_types = ["incident", "critical", "agent_monitoring_issue"]
  
  custom_payload_field = [{
    key   = "region"
    value = "us-east-1"
  }, {
    key   = "priority"
    value = "P1"
  }, {
    key   = "runbook"
    value = "https://wiki.example.com/runbooks/infrastructure"
  }]
}
```


### Event Monitoring with Different Event Types

```hcl
resource "instana_alerting_config" "alert_config_change_event_monitoring" {
  alert_name                         = "K8s-Pod-Alert"
  event_filter_event_types           = ["agent_monitoring_issue", "change", "critical", "incident", "offline", "online", "warning"]
  event_filter_query                 = "entity.kubernetes.cluster.label:demo-test"
  integration_ids                    = ["id1"]
}
```

## Generating Configuration from Existing Resources

If you have already created an alerting configuration in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your alerting configuration in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_alerting_config.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual alerting configuration ID from Instana

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

* `alert_name` - Required - the name of the alerting configuration
* `integration_ids` - Required - the list of target alerting channel ids (set of strings)
* `event_filter_query` - Optional - a dynamic focus query to restrict the alert configuration to a sub set of entities
* `event_filter_rule_ids` - Optional - list of rule IDs which are included by the alerting config (set of strings)
* `event_filter_event_types` - Optional - list of event types which are included by the alerting config (set of strings).
* `custom_payload_field` - Optional - An optional list of custom payload fields (static key/value pairs added to the event). [Details](#custom-payload-field-argument-reference)

### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Required - The value of the custom payload field

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the alerting configuration

## Import

Alerting configs can be imported using the `id`, e.g.:

```bash
terraform import instana_alerting_config.my_alerting_config 60845e4e5e6b9cf8fc2868da
```

## Best Practices

### Event Filter Queries

Use specific queries to avoid alert fatigue:

```hcl
# Good - Specific targeting
event_filter_query = "entity.type:host AND entity.zone:production AND entity.tag:critical"

# Avoid - Too broad
event_filter_query = "entity.type:host"
```

### Integration IDs

Always reference alerting channels by resource reference:

```hcl
# Good - Using resource reference
integration_ids = [instana_alerting_channel_email.ops.id]

# Avoid - Hardcoded IDs
integration_ids = ["abc123"]
```

### Custom Payload Fields

Use custom payload fields to add context:

```hcl
custom_payload_field = [{
  key   = "team"
  value = "platform-engineering"
}, {
  key   = "runbook_url"
  value = "https://wiki.example.com/runbooks/alerts"
}, {
  key   = "escalation_policy"
  value = "follow-the-sun"
}]
```

### Event Type Selection

Choose appropriate event types for your use case:

- `incident` - For critical issues requiring immediate attention
- `critical` - For severe problems affecting service
- `warning` - For potential issues that need monitoring
- `change` - For tracking deployments and configuration changes
- `offline`/`online` - For entity availability monitoring
- `agent_monitoring_issue` - For Instana agent problems

## Notes

- Either `event_filter_rule_ids` or `event_filter_event_types` should be specified, but not both
- The `event_filter_query` uses Instana's Dynamic Focus query language
- Multiple integration IDs can be specified to send alerts to multiple channels
- Custom payload fields are included in all alert notifications
- The resource ID is auto-generated by Instana upon creation
