# User Data Source

Data source to retrieve details about existing Instana users by their email address. This is useful for referencing users in roles and teams.

API Documentation: <https://instana.github.io/openapi/#operation/getUsers>

## Example Usage

```hcl
data "instana_user" "example" {
  email = "user@example.com"
}

# Use the user ID in a role resource
resource "instana_rbac_role" "example_role" {
  name = "Example Role"
  
  members {
    user_id = data.instana_user.example.id
  }
  
  permissions = ["CAN_CONFIGURE_APPLICATIONS"]
}

# Use the user ID in a team resource
resource "instana_rbac_team" "example_team" {
  tag = "example-team"
  
  members {
    user_id = data.instana_user.example.id
    
    roles {
      role_id = instana_rbac_role.example_role.id
    }
  }
}
```

## Argument Reference

* `email` - Required - the email address of the user to look up

## Attribute Reference

* `id` - the unique identifier of the user
* `full_name` - the full name of the user