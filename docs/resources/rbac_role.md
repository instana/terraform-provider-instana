# RBAC Role Resource

Management of Roles for role-based access control (RBAC). RBAC roles allow you to define permissions and assign them to team members, providing fine-grained access control across your Instana environment.

API Documentation: <https://instana.github.io/openapi/#tag/Roles>

## Example Usage

### Basic Role with Permissions

```hcl
resource "instana_rbac_role" "developer_role" {
  name = "Developer Role"

  permissions = [
    "CAN_CONFIGURE_APPLICATIONS",
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS"
  ]

  member = [
    {
      user_id = "user-id-1"
      email   = "developer1@example.com"
    },
    {
      user_id = "user-id-2"
      email   = "developer2@example.com"
      name    = "Jane Developer"
    }
  ]
}
```

### Read-Only Role

```hcl
resource "instana_rbac_role" "readonly_role" {
  name = "Read-Only Users"

  permissions = [
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_VIEW_AUDIT_LOG"
  ]

  member = [
    {
      user_id = "viewer-1"
      email   = "viewer1@example.com"
    }
  ]
}
```

### Administrator Role

```hcl
resource "instana_rbac_role" "admin_role" {
  name = "Administrators"

  permissions = [
    "CAN_CONFIGURE_APPLICATIONS",
    "CAN_CONFIGURE_EUM_APPLICATIONS",
    "CAN_CONFIGURE_AGENTS",
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_CONFIGURE_SESSION_SETTINGS",
    "CAN_CONFIGURE_INTEGRATIONS",
    "CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS",
    "CAN_CONFIGURE_GLOBAL_SYNTHETIC_SMART_ALERTS",
    "CAN_CONFIGURE_GLOBAL_INFRA_SMART_ALERTS",
    "CAN_CONFIGURE_GLOBAL_LOG_SMART_ALERTS",
    "CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD",
    "CAN_CONFIGURE_MOBILE_APP_MONITORING",
    "CAN_CONFIGURE_API_TOKENS",
    "CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS",
    "CAN_CONFIGURE_AUTHENTICATION_METHODS",
    "CAN_CONFIGURE_RELEASES",
    "CAN_VIEW_AUDIT_LOG",
    "CAN_CONFIGURE_EVENTS_AND_ALERTS",
    "CAN_CONFIGURE_MAINTENANCE_WINDOWS",
    "CAN_CONFIGURE_APPLICATION_SMART_ALERTS",
    "CAN_CONFIGURE_WEBSITE_SMART_ALERTS",
    "CAN_CONFIGURE_MOBILE_APP_SMART_ALERTS",
    "CAN_CONFIGURE_AGENT_RUN_MODE",
    "CAN_CONFIGURE_SERVICE_MAPPING",
    "CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS",
    "CAN_CONFIGURE_USERS",
    "CAN_INSTALL_NEW_AGENTS",
    "CAN_CONFIGURE_TEAMS",
    "CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS",
    "CAN_CONFIGURE_LOG_MANAGEMENT",
    "CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION"
  ]

  member = [
    {
      user_id = "admin-1"
      email   = "admin1@example.com"
      name    = "Admin User"
    }
  ]
}
```

### SRE Role

```hcl
resource "instana_rbac_role" "sre_role" {
  name = "Site Reliability Engineering"

  permissions = [
    "CAN_CONFIGURE_APPLICATIONS",
    "CAN_CONFIGURE_AGENTS",
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_CONFIGURE_EVENTS_AND_ALERTS",
    "CAN_CONFIGURE_APPLICATION_SMART_ALERTS",
    "CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS",
    "CAN_CONFIGURE_MAINTENANCE_WINDOWS",
    "CAN_INVOKE_ALERT_CHANNEL",
    "CAN_MANUALLY_CLOSE_ISSUE",
    "CAN_INSTALL_NEW_AGENTS",
    "ACCESS_INFRASTRUCTURE_ANALYZE"
  ]

  member = [
    {
      user_id = "sre-1"
      email   = "sre1@example.com"
    },
    {
      user_id = "sre-2"
      email   = "sre2@example.com"
    }
  ]
}
```

### Security Team Role

```hcl
resource "instana_rbac_role" "security_role" {
  name = "Security Team"

  permissions = [
    "CAN_VIEW_AUDIT_LOG",
    "CAN_CONFIGURE_AUTHENTICATION_METHODS",
    "CAN_CONFIGURE_API_TOKENS",
    "CAN_CONFIGURE_USERS",
    "CAN_CONFIGURE_TEAMS",
    "CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION",
    "CAN_CONFIGURE_IP_FILTERING"
  ]

  member = [
    {
      user_id = "security-lead"
      email   = "security.lead@example.com"
      name    = "Security Lead"
    }
  ]
}
```

### Monitoring Team Role

```hcl
resource "instana_rbac_role" "monitoring_role" {
  name = "Monitoring Team"

  permissions = [
    "CAN_VIEW_TRACE_DETAILS",
    "CAN_VIEW_LOGS",
    "CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS",
    "CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS",
    "CAN_CONFIGURE_EVENTS_AND_ALERTS",
    "CAN_VIEW_SYNTHETIC_TESTS",
    "CAN_VIEW_SYNTHETIC_TEST_RESULTS"
  ]

  member = [
    {
      user_id = "monitor-1"
      email   = "monitor1@example.com"
    }
  ]
}
```

### Synthetic Monitoring Role

```hcl
resource "instana_rbac_role" "synthetic_role" {
  name = "Synthetic Monitoring Team"

  permissions = [
    "CAN_VIEW_SYNTHETIC_TESTS",
    "CAN_CONFIGURE_SYNTHETIC_TESTS",
    "CAN_VIEW_SYNTHETIC_LOCATIONS",
    "CAN_CONFIGURE_SYNTHETIC_LOCATIONS",
    "CAN_USE_SYNTHETIC_CREDENTIALS",
    "CAN_CONFIGURE_SYNTHETIC_CREDENTIALS",
    "CAN_VIEW_SYNTHETIC_TEST_RESULTS",
    "CAN_CONFIGURE_GLOBAL_SYNTHETIC_SMART_ALERTS"
  ]

  member = [
    {
      user_id = "synthetic-user"
      email   = "synthetic@example.com"
    }
  ]
}
```

### Log Management Role

```hcl
resource "instana_rbac_role" "log_management_role" {
  name = "Log Management Team"

  permissions = [
    "CAN_VIEW_LOGS",
    "CAN_CONFIGURE_LOG_MANAGEMENT",
    "CAN_CONFIGURE_LOG_RETENTION_PERIOD",
    "CAN_VIEW_LOG_VOLUME",
    "CAN_DELETE_LOGS",
    "CAN_CONFIGURE_GLOBAL_LOG_SMART_ALERTS"
  ]

  member = [
    {
      user_id = "log-admin"
      email   = "log.admin@example.com"
    }
  ]
}
```

### Automation Team Role

```hcl
resource "instana_rbac_role" "automation_role" {
  name = "Automation Team"

  permissions = [
    "CAN_RUN_AUTOMATION_ACTIONS",
    "CAN_CONFIGURE_AUTOMATION_ACTIONS",
    "CAN_CONFIGURE_AUTOMATION_POLICIES",
    "CAN_DELETE_AUTOMATION_ACTION_HISTORY"
  ]

  member = [
    {
      user_id = "automation-user"
      email   = "automation@example.com"
    }
  ]
}
```

### Business Operations Role

```hcl
resource "instana_rbac_role" "bizops_role" {
  name = "Business Operations"

  permissions = [
    "CAN_VIEW_BUSINESS_PROCESS_DETAILS",
    "CAN_VIEW_BIZOPS_ALERTS",
    "CAN_CONFIGURE_BIZOPS",
    "CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS",
    "CAN_CONFIGURE_SERVICE_LEVELS",
    "CAN_CONFIGURE_SERVICE_LEVEL_CORRECTION_WINDOWS",
    "CAN_CONFIGURE_SERVICE_LEVEL_SMART_ALERTS"
  ]

  member = [
    {
      user_id = "bizops-user"
      email   = "bizops@example.com"
    }
  ]
}
```

### Multi-Environment Setup with for_each

```hcl
locals {
  team_roles = {
    development = {
      permissions = [
        "CAN_CONFIGURE_APPLICATIONS",
        "CAN_CONFIGURE_AGENTS",
        "CAN_VIEW_TRACE_DETAILS",
        "CAN_VIEW_LOGS"
      ]
      members = [
        { user_id = "dev-1", email = "dev1@example.com" },
        { user_id = "dev-2", email = "dev2@example.com" }
      ]
    }
    staging = {
      permissions = [
        "CAN_CONFIGURE_APPLICATIONS",
        "CAN_VIEW_TRACE_DETAILS",
        "CAN_VIEW_LOGS",
        "CAN_CONFIGURE_EVENTS_AND_ALERTS"
      ]
      members = [
        { user_id = "stage-1", email = "stage1@example.com" }
      ]
    }
    production = {
      permissions = [
        "CAN_VIEW_TRACE_DETAILS",
        "CAN_VIEW_LOGS",
        "CAN_CONFIGURE_EVENTS_AND_ALERTS",
        "CAN_MANUALLY_CLOSE_ISSUE"
      ]
      members = [
        { user_id = "prod-1", email = "prod1@example.com" }
      ]
    }
  }
}

resource "instana_rbac_role" "team_roles" {
  for_each = local.team_roles

  name        = "${title(each.key)} Team Role"
  permissions = each.value.permissions
  member      = each.value.members
}
```

### Role with Scope-Based Permissions

```hcl
resource "instana_rbac_role" "scoped_role" {
  name = "Scoped Access Role"

  permissions = [
    "CAN_CONFIGURE_APPLICATIONS",
    "CAN_VIEW_TRACE_DETAILS",
    "LIMITED_APPLICATIONS_SCOPE",
    "LIMITED_KUBERNETES_SCOPE",
    "ACCESS_APPLICATIONS"
  ]

  member = [
    {
      user_id = "scoped-user"
      email   = "scoped@example.com"
    }
  ]
}
```

## Argument Reference

* `name` - Required - The name of the RBAC role
* `permissions` - Required - List of permissions assigned to the role. See [Permissions Reference](#permissions-reference) for allowed values
* `member` - Required - List of role members [Details](#member-reference)

### Member Reference

* `user_id` - Required - The user ID of the role member
* `email` - Optional - The email address of the role member
* `name` - Optional - The display name of the role member

## Attributes Reference

* `id` - The ID of the RBAC role

## Permissions Reference

The following permissions can be assigned to roles:

### Application & Service Permissions
* `CAN_CONFIGURE_APPLICATIONS` - Configure application monitoring
* `CAN_CONFIGURE_EUM_APPLICATIONS` - Configure End User Monitoring applications
* `CAN_CONFIGURE_SERVICE_MAPPING` - Configure service mapping
* `CAN_CONFIGURE_APDEX` - Configure Apdex settings
* `CAN_CONFIGURE_CUSTOM_ENTITIES` - Configure custom entities
* `ACCESS_APPLICATIONS` - Access applications

### Agent & Infrastructure Permissions
* `CAN_CONFIGURE_AGENTS` - Configure agents
* `CAN_INSTALL_NEW_AGENTS` - Install new agents
* `CAN_CONFIGURE_AGENT_RUN_MODE` - Configure agent run mode
* `ACCESS_INFRASTRUCTURE_ANALYZE` - Access infrastructure analyze

### Monitoring & Observability Permissions
* `CAN_VIEW_TRACE_DETAILS` - View trace details
* `CAN_VIEW_LOGS` - View logs
* `CAN_CONFIGURE_LOG_MANAGEMENT` - Configure log management
* `CAN_CONFIGURE_LOG_RETENTION_PERIOD` - Configure log retention period
* `CAN_VIEW_LOG_VOLUME` - View log volume
* `CAN_DELETE_LOGS` - Delete logs

### Alert & Event Permissions
* `CAN_CONFIGURE_EVENTS_AND_ALERTS` - Configure events and alerts
* `CAN_CONFIGURE_MAINTENANCE_WINDOWS` - Configure maintenance windows
* `CAN_CONFIGURE_APPLICATION_SMART_ALERTS` - Configure application smart alerts
* `CAN_CONFIGURE_WEBSITE_SMART_ALERTS` - Configure website smart alerts
* `CAN_CONFIGURE_MOBILE_APP_SMART_ALERTS` - Configure mobile app smart alerts
* `CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS` - Configure global application smart alerts
* `CAN_CONFIGURE_GLOBAL_SYNTHETIC_SMART_ALERTS` - Configure global synthetic smart alerts
* `CAN_CONFIGURE_GLOBAL_INFRA_SMART_ALERTS` - Configure global infrastructure smart alerts
* `CAN_CONFIGURE_GLOBAL_LOG_SMART_ALERTS` - Configure global log smart alerts
* `CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD` - Configure global alert payload
* `CAN_INVOKE_ALERT_CHANNEL` - Invoke alert channel
* `CAN_MANUALLY_CLOSE_ISSUE` - Manually close issues
* `CAN_CONFIGURE_SERVICE_LEVEL_SMART_ALERTS` - Configure service level smart alerts

### Synthetic Monitoring Permissions
* `CAN_VIEW_SYNTHETIC_TESTS` - View synthetic tests
* `CAN_CONFIGURE_SYNTHETIC_TESTS` - Configure synthetic tests
* `CAN_VIEW_SYNTHETIC_LOCATIONS` - View synthetic locations
* `CAN_CONFIGURE_SYNTHETIC_LOCATIONS` - Configure synthetic locations
* `CAN_VIEW_SYNTHETIC_TEST_RESULTS` - View synthetic test results
* `CAN_USE_SYNTHETIC_CREDENTIALS` - Use synthetic credentials
* `CAN_CONFIGURE_SYNTHETIC_CREDENTIALS` - Configure synthetic credentials

### Mobile & Website Monitoring Permissions
* `CAN_CONFIGURE_MOBILE_APP_MONITORING` - Configure mobile app monitoring
* `CAN_CONFIGURE_WEBSITE_CONVERSIONS` - Configure website conversions
* `ACCESS_MOBILE_APPS` - Access mobile apps
* `ACCESS_WEBSITES` - Access websites
* `ACCESS_SYNTHETICS` - Access synthetics

### Dashboard & Visualization Permissions
* `CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS` - Edit all accessible custom dashboards
* `CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS` - Create public custom dashboards

### User & Access Management Permissions
* `CAN_CONFIGURE_USERS` - Configure users
* `CAN_CONFIGURE_TEAMS` - Configure teams
* `CAN_CONFIGURE_AUTHENTICATION_METHODS` - Configure authentication methods
* `CAN_CONFIGURE_API_TOKENS` - Configure API tokens
* `CAN_CONFIGURE_PERSONAL_API_TOKENS` - Configure personal API tokens
* `CAN_CONFIGURE_IP_FILTERING` - Configure IP filtering

### Integration & Configuration Permissions
* `CAN_CONFIGURE_INTEGRATIONS` - Configure integrations
* `CAN_CONFIGURE_SESSION_SETTINGS` - Configure session settings
* `CAN_CONFIGURE_RELEASES` - Configure releases
* `CAN_CONFIGURE_DATABASE_MANAGEMENT` - Configure database management
* `CAN_CONFIGURE_SUBTRACES` - Configure subtraces
* `CAN_CONFIGURE_LLM` - Configure LLM

### Service Level & BizOps Permissions
* `CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS` - Configure service level indicators
* `CAN_CONFIGURE_SERVICE_LEVELS` - Configure service levels
* `CAN_CONFIGURE_SERVICE_LEVEL_CORRECTION_WINDOWS` - Configure service level correction windows
* `CAN_VIEW_BUSINESS_PROCESS_DETAILS` - View business process details
* `CAN_VIEW_BIZOPS_ALERTS` - View BizOps alerts
* `CAN_CONFIGURE_BIZOPS` - Configure BizOps

### Automation Permissions
* `CAN_RUN_AUTOMATION_ACTIONS` - Run automation actions
* `CAN_CONFIGURE_AUTOMATION_ACTIONS` - Configure automation actions
* `CAN_CONFIGURE_AUTOMATION_POLICIES` - Configure automation policies
* `CAN_DELETE_AUTOMATION_ACTION_HISTORY` - Delete automation action history

### Diagnostic & Troubleshooting Permissions
* `CAN_CREATE_THREAD_DUMP` - Create thread dumps
* `CAN_CREATE_HEAP_DUMP` - Create heap dumps

### Administrative Permissions
* `CAN_VIEW_AUDIT_LOG` - View audit log
* `CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION` - View account and billing information

### Scope-Based Permissions
* `LIMITED_APPLICATIONS_SCOPE` - Limited applications scope
* `LIMITED_KUBERNETES_SCOPE` - Limited Kubernetes scope
* `LIMITED_LINUX_KVM_HYPERVISOR_SCOPE` - Limited Linux KVM hypervisor scope
* `LIMITED_VSPHERE_SCOPE` - Limited vSphere scope
* `LIMITED_OPENSTACK_SCOPE` - Limited OpenStack scope
* `LIMITED_ZHMC_SCOPE` - Limited z/HMC scope
* `LIMITED_XENSERVER_SCOPE` - Limited XenServer scope
* `LIMITED_POWERVC_SCOPE` - Limited PowerVC scope
* `LIMITED_SAP_SCOPE` - Limited SAP scope
* `LIMITED_PCF_SCOPE` - Limited PCF scope
* `LIMITED_SYNTHETICS_SCOPE` - Limited synthetics scope
* `LIMITED_SERVICE_LEVEL_SCOPE` - Limited service level scope
* `LIMITED_AUTOMATION_SCOPE` - Limited automation scope
* `LIMITED_BIZOPS_SCOPE` - Limited BizOps scope
* `LIMITED_PHMC_SCOPE` - Limited PHMC scope
* `LIMITED_GEN_AI_SCOPE` - Limited Gen AI scope
* `LIMITED_INFRASTRUCTURE_SCOPE` - Limited infrastructure scope
* `LIMITED_NUTANIX_SCOPE` - Limited Nutanix scope
* `LIMITED_WINDOWS_HYPERVISOR_SCOPE` - Limited Windows hypervisor scope
* `LIMITED_AI_GATEWAY_SCOPE` - Limited AI Gateway scope
* `LIMITED_MOBILE_APPS_SCOPE` - Limited mobile apps scope
* `LIMITED_WEBSITES_SCOPE` - Limited websites scope

## Import

RBAC Roles can be imported using the `id` of the role, e.g.:

```bash
$ terraform import instana_rbac_role.my_role 60845e4e5e6b9cf8fc2868da
```

## Notes

* The ID is auto-generated by Instana
* Roles define permissions that can be assigned to users through group memberships
* For scoped access control, use `instana_rbac_group` which combines permissions with scope restrictions
* Permissions are additive - users with multiple role assignments get the union of all permissions
* Members can be added or removed without affecting the permission set
* Use scope-based permissions (e.g., `LIMITED_APPLICATIONS_SCOPE`) in combination with groups to restrict access to specific resources