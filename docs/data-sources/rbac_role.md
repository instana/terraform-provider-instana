# Data Source: instana_rbac_role

Use this data source to get information about an Instana RBAC Role by its ID.

## Example Usage

```hcl
 data "instana_rbac_role" "example" {
   id = "<role-id>"
 }
```

## Argument Reference

- `id` (Required) - The ID of the RBAC Role.

## Attribute Reference

- `id` - The ID of the RBAC Role.
- `name` - The name of the RBAC Role.
- `permissions` - The set of permissions assigned to the RBAC Role.
