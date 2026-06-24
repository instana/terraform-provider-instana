---
page_title: "instana_rbac_team Data Source - terraform-provider-instana"
subcategory: ""
description: |-
  Data source for an Instana RBAC team
---

# instana_rbac_team (Data Source)

Data source for an Instana RBAC team. RBAC teams group users and define their access scope within Instana.

## Example Usage

```terraform
# Lookup team by tag/name
data "instana_rbac_team" "example" {
  tag = "my-team"
}

# Lookup team by ID
data "instana_rbac_team" "example_by_id" {
  id = "team-id-123"
}

# Use the team data in other resources
resource "instana_rbac_group" "example" {
  name = "My Group"
  
  member {
    user_id = "user-123"
  }
}
```

## Schema

### Required

Exactly one of the following must be specified:

- `id` (String) The ID of the RBAC team
- `tag` (String) The tag/name of the RBAC team

### Read-Only

- `id` (String) The ID of the RBAC team
- `tag` (String) The tag/name of the RBAC team
- `info` (Block, Optional) Additional information about the team (see [below for nested schema](#nestedblock--info))
- `member` (Block Set) The members of the team (see [below for nested schema](#nestedblock--member))
- `scope` (Block, Optional) The scope configuration for the team (see [below for nested schema](#nestedblock--scope))

<a id="nestedblock--info"></a>
### Nested Schema for `info`

Read-Only:

- `description` (String) The description of the team

<a id="nestedblock--member"></a>
### Nested Schema for `member`

Read-Only:

- `user_id` (String) The user ID of the team member
- `roles` (Block Set) The roles assigned to the team member (see [below for nested schema](#nestedblock--member--roles))

<a id="nestedblock--member--roles"></a>
### Nested Schema for `member.roles`

Read-Only:

- `role_id` (String) The ID of the role

<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Read-Only:

- `access_permissions` (Set of String) The access permissions for the team
- `applications` (Set of String) The application IDs accessible to the team
- `kubernetes_clusters` (Set of String) The Kubernetes cluster IDs accessible to the team
- `kubernetes_namespaces` (Set of String) The Kubernetes namespace IDs accessible to the team
- `mobile_apps` (Set of String) The mobile app IDs accessible to the team
- `websites` (Set of String) The website IDs accessible to the team
- `infra_dfq_filter` (String) The infrastructure DFQ filter for the team
- `action_filter` (String) The action filter for the team
- `log_filter` (String) The log filter for the team
- `business_perspectives` (Set of String) The business perspective IDs accessible to the team
- `slo_ids` (Set of String) The SLO IDs accessible to the team
- `synthetic_tests` (Set of String) The synthetic test IDs accessible to the team
- `synthetic_credentials` (Set of String) The synthetic credential IDs accessible to the team
- `tag_ids` (Set of String) The tag IDs accessible to the team
- `restricted_application_filter` (Block, Optional) The restricted application filter configuration (see [below for nested schema](#nestedblock--scope--restricted_application_filter))

<a id="nestedblock--scope--restricted_application_filter"></a>
### Nested Schema for `scope.restricted_application_filter`

Read-Only:

- `label` (String) The label for the restricted application filter
- `scope` (String) The scope of the restricted application filter. Valid values: `INCLUDE_NO_DOWNSTREAM`, `INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING`, `INCLUDE_ALL_DOWNSTREAM`
- `tag_filter_expression` (String) The tag filter expression for the restricted application filter