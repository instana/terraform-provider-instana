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
resource "instana_alerting_config" "example" {
  alert_name            = "Critical Infrastructure Alerts"
  integration_ids       = [
    instana_alerting_channel_email.ops_team.id,
    instana_alerting_channel_slack.incidents.id
  ]
  event_filter_query    = "entity.type:host AND entity.zone:production"
  event_filter_rule_ids = ["rule-1", "rule-2"]
}
```

### Configuration with Event Types

```hcl
resource "instana_alerting_config" "example" {
  alert_name               = "Application Monitoring"
  integration_ids          = [instana_alerting_channel_pagerduty.oncall.id]
  event_filter_query       = "entity.application.id:'my-app-id'"
  event_filter_event_types = ["incident", "critical"]

  custom_payload_field = [{
    key   = "environment"
    value = "production"
  }, {
    key   = "team"
    value = "platform"
  }]
}
```

### Application-Specific Alerts

```hcl
resource "instana_alerting_config" "app_alerts" {
  alert_name         = "My Application Alerts"
  integration_ids    = [instana_alerting_channel_webhook.custom.id]
  event_filter_query = "entity.application.id:'${instana_application_config.my_app.id}'"
  
  event_filter_event_types = ["incident", "critical", "warning"]
}
```

### Multi-Channel Configuration

```hcl
resource "instana_alerting_config" "multi_channel" {
  alert_name = "Critical System Alerts"
  
  integration_ids = [
    instana_alerting_channel_email.ops.id,
    instana_alerting_channel_slack.alerts.id,
    instana_alerting_channel_pagerduty.oncall.id,
    instana_alerting_channel_webhook.monitoring.id
  ]
  
  event_filter_query       = "entity.zone:production AND entity.type:host"
  event_filter_event_types = ["incident", "critical"]
  
  custom_payload_field = [{
    key   = "severity"
    value = "high"
  }, {
    key   = "escalation_policy"
    value = "immediate"
  }]
}
```

### Database Monitoring Alerts

```hcl
resource "instana_alerting_config" "database_alerts" {
  alert_name               = "Database Performance Alerts"
  integration_ids          = [instana_alerting_channel_email.dba_team.id]
  event_filter_query       = "entity.type:database AND entity.database.type:postgresql"
  event_filter_event_types = ["critical", "warning"]
  
  custom_payload_field = [{
    key   = "database_type"
    value = "postgresql"
  }, {
    key   = "alert_category"
    value = "performance"
  }]
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

### Service-Specific Alerts

```hcl
resource "instana_alerting_config" "service_alerts" {
  alert_name         = "Payment Service Alerts"
  integration_ids    = [instana_alerting_channel_slack.payments.id]
  event_filter_query = "entity.service.name:'payment-service' AND entity.zone:production"
  
  event_filter_rule_ids = [
    instana_custom_event_specification.payment_latency.id,
    instana_custom_event_specification.payment_errors.id
  ]
  
  custom_payload_field = [{
    key   = "service"
    value = "payment-service"
  }, {
    key   = "business_impact"
    value = "high"
  }]
}
```

### Kubernetes Cluster Monitoring

```hcl
resource "instana_alerting_config" "k8s_alerts" {
  alert_name      = "Kubernetes Cluster Alerts"
  integration_ids = [instana_alerting_channel_webhook.k8s_monitor.id]
  
  event_filter_query = "entity.kubernetes.cluster.name:'prod-cluster' AND entity.type:kubernetesNode"
  event_filter_event_types = [
    "incident",
    "critical",
    "offline",
    "agent_monitoring_issue"
  ]
  
  custom_payload_field = [{
    key   = "cluster"
    value = "prod-cluster"
  }, {
    key   = "platform"
    value = "kubernetes"
  }]
}
```

### Change Event Monitoring

```hcl
resource "instana_alerting_config" "change_tracking" {
  alert_name               = "Deployment Change Tracking"
  integration_ids          = [instana_alerting_channel_email.devops.id]
  event_filter_query       = "entity.application.id:'${instana_application_config.my_app.id}'"
  event_filter_event_types = ["change"]
  
  custom_payload_field = [{
    key   = "notification_type"
    value = "deployment"
  }]
}
```

### Comprehensive Monitoring Setup

```hcl
resource "instana_alerting_config" "comprehensive" {
  alert_name = "Production Environment Monitoring"
  
  integration_ids = [
    instana_alerting_channel_email.ops.id,
    instana_alerting_channel_slack.prod_alerts.id,
    instana_alerting_channel_pagerduty.oncall.id
  ]
  
  event_filter_query = join(" AND ", [
    "entity.zone:production",
    "(entity.type:host OR entity.type:container OR entity.type:service)",
    "NOT entity.tag:monitoring-excluded"
  ])
  
  event_filter_event_types = [
    "incident",
    "critical",
    "warning",
    "offline",
    "agent_monitoring_issue"
  ]
  
  custom_payload_field = [{
    key   = "environment"
    value = "production"
  }, {
    key   = "severity_level"
    value = "high"
  }, {
    key   = "escalation_required"
    value = "true"
  }, {
    key   = "sla_impact"
    value = "yes"
  }]
}
```

## Argument Reference

* `alert_name` - Required - the name of the alerting configuration
* `integration_ids` - Required - the list of target alerting channel ids (set of strings)
* `event_filter_query` - Optional - a dynamic focus query to restrict the alert configuration to a sub set of entities
* `event_filter_rule_ids` - Optional - list of rule IDs which are included by the alerting config (set of strings)
* `event_filter_event_types` - Optional - list of event types which are included by the alerting config (set of strings). Allowed values: `incident`, `critical`, `warning`, `change`, `online`, `offline`, `agent_monitoring_issue`, `none`
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
