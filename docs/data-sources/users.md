# Users Data Source

Data source to retrieve details about multiple Instana users by their email addresses. This is useful for bulk operations when referencing multiple users in roles and teams.

API Documentation: <https://instana.github.io/openapi/#operation/getUsers>

## Example Usage

### Basic Usage - Query Multiple Users

```hcl
data "instana_users" "team_members" {
  emails = [
    "user1@example.com",
    "user2@example.com",
    "user3@example.com"
  ]
}

# Use the users in a role resource with for expression
resource "instana_rbac_role" "example_role" {
  name = "Example Role"
  
  member = [
    for user in data.instana_users.team_members.users : {
      user_id = user.id
    }
  ]
  
  permissions = ["CAN_CONFIGURE_APPLICATIONS"]
}
```

### Advanced Usage - Team with Multiple Members

```hcl
data "instana_users" "dev_team" {
  emails = [
    "dev1@example.com",
    "dev2@example.com",
    "dev3@example.com"
  ]
}

resource "instana_rbac_role" "developer_role" {
  name = "Developer"
  permissions = ["CAN_CONFIGURE_APPLICATIONS", "CAN_VIEW_TRACE_DETAILS"]
}

resource "instana_rbac_team" "development_team" {
  tag = "development-team"
  
  member = [
    for user in data.instana_users.dev_team.users : {
      user_id = user.id
      role = [
        {
          role_id = instana_rbac_role.developer_role.id
        }
      ]
    }
  ]
}
```

### Output User Information

```hcl
data "instana_users" "admins" {
  emails = [
    "admin1@example.com",
    "admin2@example.com"
  ]
}

output "admin_user_ids" {
  value = {
    for user in data.instana_users.admins.users :
    user.email => user.id
  }
}

output "admin_names" {
  value = [for user in data.instana_users.admins.users : user.full_name]
}
```

## Argument Reference

* `emails` - Required - list of email addresses of the users to look up

## Attribute Reference

* `id` - the unique identifier for this datasource query
* `users` - list of user objects matching the provided emails. Each user object contains:
  * `id` - the unique identifier of the user
  * `email` - the email address of the user
  * `full_name` - the full name of the user

## Notes

- Only users with matching email addresses will be returned
- If an email doesn't match any user, it will be silently skipped (no error)
- The order of users in the result may not match the order of emails in the input
- For querying a single user, consider using the `instana_user` data source instead