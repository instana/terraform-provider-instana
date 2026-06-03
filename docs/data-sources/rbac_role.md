# Data Source: instana_rbac_role


Use this data source to get information about an Instana RBAC Role by its ID or name.

## Example Usage


```hcl
# Fetch by ID
data "instana_rbac_role" "by_id" {
  id = "<role-id>"
}

# Fetch by name
data "instana_rbac_role" "by_name" {
  name = "My Role Name"
}
```


## Argument Reference

- `id` (Optional) - The ID of the RBAC Role. Exactly one of `id` or `name` must be specified.
- `name` (Optional) - The name of the RBAC Role. Exactly one of `id` or `name` must be specified.

## Attribute Reference

- `id` - The ID of the RBAC Role.
- `name` - The name of the RBAC Role.
- `permissions` - The set of permissions assigned to the RBAC Role.
