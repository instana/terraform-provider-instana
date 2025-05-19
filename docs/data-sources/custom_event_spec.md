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

* `name` - Required - the name of the custom event.
* `entity_type` - Required - the entity type of the custom event.
