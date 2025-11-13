# Builtin Event Specification Data Source

Data source to get the specification of builtin events from Instana API. This allows you to retrieve the specification
by UI name and Plugin ID and reference it in other resources such as Alerting Configurations.

API Documentation: <https://instana.github.io/openapi/#operation/getBuiltInEventSpecifications>

## Example Usage

```hcl
data "instana_builtin_event_spec" "host_system_load_too_high" {
  name = "System load too high"
  short_plugin_id = "host"
}
```

## Argument Reference

* `name` - Required - the name of the builtin event
* `short_plugin_id` - Required - the short plugin ID of the builtin event (can be retrieved from <https://instana.github.io/openapi/#operation/getInfrastructureCatalogPlugins>)

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - The builtin event specification identifier
* `description` - The description text of the builtin event specification
* `severity` - The severity of the event (e.g., "warning", "critical")
* `severity_code` - The numeric severity code of the event
* `triggering` - Boolean flag indicating if the event should trigger an incident
* `enabled` - Boolean flag indicating if the event specification is enabled
