# Application Config Data Source

Data source to retrieve details about an existing application configuration (application perspective) in Instana. This allows you to reference existing application configurations by `id` or `name` (label) without hardcoding IDs in other resources.

API Documentation: <https://instana.github.io/openapi/#operation/getApplicationConfig>

## Example Usage

### Lookup by Name (Label)

When searching by `name`, if multiple application configurations share the same name/label, the first matching result is returned.

```hcl
data "instana_application_config" "example" {
  name = "My Application Perspective"
}

# Reference in other resources
output "app_id" {
  value = data.instana_application_config.example.id
}
```

### Lookup by ID

```hcl
data "instana_application_config" "example" {
  id = "60845e4e5e6b9cf8fc2868da"
}
```

## Argument Reference

* `id` - (Optional) The unique identifier of the application configuration. Exactly one of `id` or `name` must be specified.
* `name` - (Optional) The name (label) of the application configuration. If multiple configurations have the same name, the first match will be returned. Exactly one of `id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the application configuration.
* `name` - The name/label of the application configuration.
* `label` - The label of the application configuration.
* `scope` - The scope of the application configuration (e.g., `INCLUDE_NO_DOWNSTREAM`, `INCLUDE_ALL_DOWNSTREAM`, `INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING`).
* `boundary_scope` - The boundary scope of the application configuration (`INBOUND`, `ALL`, `DEFAULT`).
* `tag_filter` - The normalized tag filter expression string.
* `access_rules` - List of access rules applied to the application configuration:
  * `access_type` - The access type (`READ_WRITE`, `READ_ONLY`).
  * `relation_type` - The relation type (`GLOBAL`, etc.).
  * `related_id` - The ID of the related entity (user, group, API token, etc.).
