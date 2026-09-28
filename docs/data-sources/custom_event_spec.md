# Custom Event Specification Data Source

Data source to get the specification of custom events from Instana API. This allows you to retrieve the specification
by UI name and entity type and reference it in other resources such as Automation Policies.

API Documentation: <https://instana.github.io/openapi/#operation/getCustomEventSpecifications>

## Example Usage

```hcl
data "instana_custom_event_spec" "host_system_load_too_high" {
  name = "System load too high"
  entity_type = "host"
}
```

## Argument Reference

* `name` - Required - Name of the custom event.
* `entity_type` - Required - Entity type of the custom event.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - The custom event specification identifier.
* `description` - The description text of the custom event specification.
* `query` - The dynamic filter query for which the rule should be applied to.
* `enabled` - Boolean flag if the rule should be enabled. Default is `true`.
* `triggering` - Boolean flag if the rule should trigger an incident. Default is `false`.
* `expiration_time` - The grace period in milliseconds until the issue is closed.
* `transient_event_enabled` - Indicates if transient event handling is enabled. Default is `true`.
* `transient_event_threshold` - The duration in milliseconds that an event must persist before it is considered non-transient. Default is `300000`.
* `transient_event_alert_muted` - Indicates if alerts for transient events are muted. Default is `false`.