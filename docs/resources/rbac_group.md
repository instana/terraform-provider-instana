# RBAC Group Resource

> **⚠️ DEPRECATION NOTICE**: RBAC Group is deprecated. Please use [RBAC Roles](rbac_role.md) and [RBAC Teams](rbac_team.md) instead.


Management of Groups for role-based access control (RBAC). RBAC groups allow you to define permissions and scope access to specific resources for team members.

API Documentation: <https://instana.github.io/openapi/#tag/Groups>

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.


## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

The main changes are in how `permission_set` and `member` blocks are defined. In v6, these use attribute syntax instead of block syntax.

#### OLD (v5.x) Syntax:
```hcl
resource "instana_rbac_group" "example" {
  name = "DevOps Team"

  permission_set {
    application_ids = ["app1", "app2"]
    permissions = ["CAN_CONFIGURE_APPLICATIONS"]
  }
  
  member {
    user_id = "user1"
    email = "user1@example.com"
  }
  
  member {
    user_id = "user2"
    email = "user2@example.com"
  }
}
```

#### NEW (v6.x) Syntax:
```hcl
resource "instana_rbac_group" "dev_team" {
  name = "Development Team"

  permission_set = {
    permissions = [
      "CAN_CONFIGURE_APPLICATIONS",
      "CAN_VIEW_TRACE_DETAILS",
      "CAN_VIEW_LOGS"
    ]
  }

  member = [
    {
      user_id = "68ca71ee952dd200010ecc04" # replace with valid user_id
      email = "abhishek.kumar144@ibm.com" # replace with valid user email
    }
  ]
}
```

### Key Syntax Changes

1. **Permission Set**: `permission_set { }` → `permission_set = { }`
2. **Members**: Multiple `member { }` blocks → Single `member = [{ }, { }]` list
3. **All nested attributes**: Use `= { }` or `= [{ }]` syntax

## Example Usage

```hcl
resource "instana_rbac_group" "dev_team" {
  name = "Development Team"

  permission_set = {
    permissions = [
      "CAN_CONFIGURE_APPLICATIONS",
      "CAN_VIEW_TRACE_DETAILS",
      "CAN_VIEW_LOGS"
    ]
  }

  member = [
    {
      user_id = "68ca71ee952dd200010ecc04" # replace with valid user_id
      email = "abhishek.kumar144@ibm.com" # replace with valid user email
    }
  ]
}
```

## Generating Configuration from Existing Resources

If you have already created a RBAC group in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your RBAC group in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_rbac_group.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual RBAC group ID from Instana

3. **Generate the Configuration**: Run the following Terraform command:

```bash
terraform plan -generate-config-out=generated.tf
```

This will:
- Import the existing resource state
- Generate the complete Terraform configuration in `generated.tf`
- Show you what will be imported

4. **Review and Apply**: Review the generated configuration in `generated.tf` and make any necessary adjustments.

   - **To import the existing resource**: Keep the import block and run `terraform apply`. This will import the resource into your Terraform state and link it to the existing Instana resource.
   
   - **To create a new resource**: If you only need the configuration structure as a template, remove the import block from your configuration. Modify the generated configuration as needed, and when you run `terraform apply`, it will create a new resource in Instana instead of importing the existing one.

```bash
terraform apply
```

## Argument Reference

* `name` - Required - The name of the RBAC group
* `permission_set` - Optional - Configuration block to describe the assigned permissions [Details](#permission-set-reference)
* `member` - Optional - List of group members [Details](#member-reference)

### Permission Set Reference

* `application_ids` - Optional - List of application IDs which are permitted to the given group
* `kubernetes_cluster_uuids` - Optional - List of Kubernetes Cluster UUIDs which are permitted to the given group
* `kubernetes_namespaces_uuids` - Optional - List of Kubernetes Namespace UUIDs which are permitted to the given group
* `mobile_app_ids` - Optional - List of mobile app IDs which are permitted to the given group
* `website_ids` - Optional - List of website IDs which are permitted to the given group
* `infra_dfq_filter` - Optional - A dynamic focus query to restrict access to a limited set of infrastructure resources
* `permissions` - Optional - The list of permissions granted to the given group. Allowed values:
  * `CAN_CONFIGURE_APPLICATIONS` - Configure application monitoring
  * `CAN_CONFIGURE_EUM_APPLICATIONS` - Configure End User Monitoring applications
  * `CAN_CONFIGURE_AGENTS` - Configure agents
  * `CAN_VIEW_TRACE_DETAILS` - View trace details
  * `CAN_VIEW_LOGS` - View logs
  * `CAN_CONFIGURE_SESSION_SETTINGS` - Configure session settings
  * `CAN_CONFIGURE_INTEGRATIONS` - Configure integrations
  * `CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS` - Configure global application smart alerts
  * `CAN_CONFIGURE_GLOBAL_SYNTHETIC_SMART_ALERTS` - Configure global synthetic smart alerts
  * `CAN_CONFIGURE_GLOBAL_INFRA_SMART_ALERTS` - Configure global infrastructure smart alerts
  * `CAN_CONFIGURE_GLOBAL_LOG_SMART_ALERTS` - Configure global log smart alerts
  * `CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD` - Configure global alert payload
  * `CAN_CONFIGURE_MOBILE_APP_MONITORING` - Configure mobile app monitoring
  * `CAN_CONFIGURE_API_TOKENS` - Configure API tokens
  * `CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS` - Configure service level indicators
  * `CAN_CONFIGURE_AUTHENTICATION_METHODS` - Configure authentication methods
  * `CAN_CONFIGURE_RELEASES` - Configure releases
  * `CAN_VIEW_AUDIT_LOG` - View audit log
  * `CAN_CONFIGURE_EVENTS_AND_ALERTS` - Configure events and alerts
  * `CAN_CONFIGURE_MAINTENANCE_WINDOWS` - Configure maintenance windows
  * `CAN_CONFIGURE_APPLICATION_SMART_ALERTS` - Configure application smart alerts
  * `CAN_CONFIGURE_WEBSITE_SMART_ALERTS` - Configure website smart alerts
  * `CAN_CONFIGURE_MOBILE_APP_SMART_ALERTS` - Configure mobile app smart alerts
  * `CAN_CONFIGURE_AGENT_RUN_MODE` - Configure agent run mode
  * `CAN_CONFIGURE_SERVICE_MAPPING` - Configure service mapping
  * `CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS` - Edit all accessible custom dashboards
  * `CAN_CONFIGURE_USERS` - Configure users
  * `CAN_INSTALL_NEW_AGENTS` - Install new agents
  * `CAN_CONFIGURE_TEAMS` - Configure teams
  * `CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS` - Create public custom dashboards
  * `CAN_CONFIGURE_LOG_MANAGEMENT` - Configure log management
  * `CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION` - View account and billing information
  * `CAN_VIEW_SYNTHETIC_TESTS` - View synthetic tests
  * `CAN_VIEW_SYNTHETIC_LOCATIONS` - View synthetic locations
  * `CAN_CREATE_THREAD_DUMP` - Create thread dumps
  * `CAN_CREATE_HEAP_DUMP` - Create heap dumps
  * `CAN_CONFIGURE_DATABASE_MANAGEMENT` - Configure database management
  * `CAN_CONFIGURE_LOG_RETENTION_PERIOD` - Configure log retention period
  * `CAN_CONFIGURE_PERSONAL_API_TOKENS` - Configure personal API tokens
  * `ACCESS_INFRASTRUCTURE_ANALYZE` - Access infrastructure analyze
  * `CAN_VIEW_LOG_VOLUME` - View log volume
  * `CAN_RUN_AUTOMATION_ACTIONS` - Run automation actions
  * `CAN_VIEW_SYNTHETIC_TEST_RESULTS` - View synthetic test results
  * `CAN_INVOKE_ALERT_CHANNEL` - Invoke alert channel
  * `CAN_MANUALLY_CLOSE_ISSUE` - Manually close issues
  * `CAN_DELETE_LOGS` - Delete logs
  * `CAN_CONFIGURE_SYNTHETIC_TESTS` - Configure synthetic tests
  * `CAN_VIEW_BUSINESS_PROCESS_DETAILS` - View business process details
  * `CAN_VIEW_BIZOPS_ALERTS` - View BizOps alerts
  * `CAN_USE_SYNTHETIC_CREDENTIALS` - Use synthetic credentials
  * `CAN_DELETE_AUTOMATION_ACTION_HISTORY` - Delete automation action history
  * `CAN_CONFIGURE_SYNTHETIC_LOCATIONS` - Configure synthetic locations
  * `CAN_CONFIGURE_SYNTHETIC_CREDENTIALS` - Configure synthetic credentials
  * `CAN_CONFIGURE_SUBTRACES` - Configure subtraces
  * `CAN_CONFIGURE_LLM` - Configure LLM
  * `CAN_CONFIGURE_BIZOPS` - Configure BizOps
  * `CAN_CONFIGURE_AUTOMATION_POLICIES` - Configure automation policies
  * `CAN_CONFIGURE_AUTOMATION_ACTIONS` - Configure automation actions

### Member Reference

* `user_id` - Required - The user ID of the group member
* `email` - Optional - The email address of the group member

## Attributes Reference

* `id` - The ID of the RBAC group

## Import

RBAC Groups can be imported using the `id` of the group, e.g.:

```bash
$ terraform import instana_rbac_group.my_group 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Groups can be scoped to specific resources (applications, Kubernetes clusters, websites, etc.)
* The `infra_dfq_filter` uses Dynamic Focus Query syntax to filter infrastructure resources
* Permissions are additive - users with multiple group memberships get the union of all permissions
* Members can be added or removed without affecting the permission set
* Use the `for_each` meta-argument to create multiple groups with similar configurations
