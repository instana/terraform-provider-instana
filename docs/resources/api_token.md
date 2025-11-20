# API Token Resource

Management of API Tokens for programmatic access to Instana.

API Documentation: <https://instana.github.io/openapi/#operation/getApiToken

 **⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)**

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block struture.

## Migration Guide (v5 to v6)

### Syntax Changes Overview

 **Major Changes:**
 - All attributes are now top-level attributes (no nested blocks)
 - Boolean attributes now use explicit `true`/`false` values with defaults
 - The `id`, `access_granting_token`, and `internal_id` attributes are computed
 - All permission and scope attributes have default values of `false`
 - Attribute syntax remains the same (key = value), but schema validation is stricter

 **Migration Example:**
 ```hcl
 # OLD (SDK v2)
 resource "instana_api_token" "example" {
   name                          = "my-token"
   can_configure_service_mapping = true
   # Omitted attributes defaulted to false
 }

 # NEW (Plugin Framework - Same Syntax, Enhanced Validation)
 resource "instana_api_token" "example" {
   name                          = "my-token"
   can_configure_service_mapping = true
   # Omitted attributes explicitly default to false
   # All boolean attributes now have computed defaults
 }
 ```

 The syntax remains largely compatible, but the framework provides better validation and clearer default behavior.

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

### Syntax Changes Overview

- All attributes are now top-level attributes (no nested blocks)
- Boolean attributes now use explicit `true`/`false` values with defaults
- The `id`, `access_granting_token`, and `internal_id` attributes are computed
- All permission and scope attributes have default values of `false`
- Attribute syntax remains the same (key = value), but schema validation is stricter

#### OLD (v5.x) Syntax:

### Basic API Token

#### Minimal Configuration
```hcl
resource "instana_api_token" "basic" {
  name = "basic-api-token"
}
```

#### Read-Only Token
```hcl
resource "instana_api_token" "readonly" {
  name                = "readonly-token"
  can_view_logs       = true
  can_view_trace_details = true
}
```

### Service Mapping and EUM Configuration

#### Service Mapping Token
```hcl
resource "instana_api_token" "service_mapping" {
  name                          = "service-mapping-token"
  can_configure_service_mapping = true
  can_configure_applications    = true
}
```

#### EUM Applications Token
```hcl
resource "instana_api_token" "eum_config" {
  name                           = "eum-configuration-token"
  can_configure_eum_applications = true
  can_configure_applications     = true
}
```

#### Mobile App Monitoring Token
```hcl
resource "instana_api_token" "mobile_monitoring" {
  name                               = "mobile-monitoring-token"
  can_configure_mobile_app_monitoring = true
  can_configure_applications          = true
}
```

### Agent Management Tokens

#### Agent Installation Token
```hcl
resource "instana_api_token" "agent_install" {
  name                   = "agent-installation-token"
  can_install_new_agents = true
  can_configure_agents   = true
}
```

#### Agent Configuration Token
```hcl
resource "instana_api_token" "agent_config" {
  name                      = "agent-configuration-token"
  can_configure_agents      = true
  can_configure_agent_run_mode = true
}
```

### Alert and Event Management Tokens

#### Events and Alerts Token
```hcl
resource "instana_api_token" "events_alerts" {
  name                            = "events-alerts-token"
  can_configure_events_and_alerts = true
  can_invoke_alert_channel        = true
}
```

#### Smart Alerts Configuration Token
```hcl
resource "instana_api_token" "smart_alerts" {
  name                                   = "smart-alerts-token"
  can_configure_application_smart_alerts = true
  can_configure_website_smart_alerts     = true
  can_configure_mobile_app_smart_alerts  = true
}
```

#### Global Smart Alerts Token
```hcl
resource "instana_api_token" "global_smart_alerts" {
  name                                           = "global-smart-alerts-token"
  can_configure_global_application_smart_alerts  = true
  can_configure_global_synthetic_smart_alerts    = true
  can_configure_global_infra_smart_alerts        = true
  can_configure_global_log_smart_alerts          = true
  can_configure_global_alert_payload             = true
}
```

#### Maintenance Windows Token
```hcl
resource "instana_api_token" "maintenance" {
  name                                = "maintenance-windows-token"
  can_configure_maintenance_windows   = true
  can_configure_events_and_alerts     = true
}
```

### Integration and User Management Tokens

#### Integrations Token
```hcl
resource "instana_api_token" "integrations" {
  name                       = "integrations-token"
  can_configure_integrations = true
}
```

#### User Management Token
```hcl
resource "instana_api_token" "user_management" {
  name                                 = "user-management-token"
  can_configure_users                  = true
  can_configure_teams                  = true
  can_configure_authentication_methods = true
}
```

### Synthetic Monitoring Tokens

#### Synthetic Tests Configuration
```hcl
resource "instana_api_token" "synthetic_config" {
  name                            = "synthetic-config-token"
  can_configure_synthetic_tests   = true
  can_configure_synthetic_locations = true
  can_configure_synthetic_credentials = true
}
```

#### Synthetic Tests Read-Only
```hcl
resource "instana_api_token" "synthetic_readonly" {
  name                           = "synthetic-readonly-token"
  can_view_synthetic_tests       = true
  can_view_synthetic_locations   = true
  can_view_synthetic_test_results = true
}
```

#### Synthetic Tests Full Access
```hcl
resource "instana_api_token" "synthetic_full" {
  name                                = "synthetic-full-access-token"
  can_configure_synthetic_tests       = true
  can_configure_synthetic_locations   = true
  can_configure_synthetic_credentials = true
  can_view_synthetic_tests            = true
  can_view_synthetic_locations        = true
  can_view_synthetic_test_results     = true
  can_use_synthetic_credentials       = true
}
```

### Automation Tokens

#### Automation Actions Token
```hcl
resource "instana_api_token" "automation_actions" {
  name                               = "automation-actions-token"
  can_configure_automation_actions   = true
  can_run_automation_actions         = true
  can_delete_automation_action_history = true
}
```

#### Automation Policies Token
```hcl
resource "instana_api_token" "automation_policies" {
  name                              = "automation-policies-token"
  can_configure_automation_policies = true
  can_configure_automation_actions  = true
}
```

### Log Management Tokens

#### Log Configuration Token
```hcl
resource "instana_api_token" "log_config" {
  name                            = "log-configuration-token"
  can_configure_log_management    = true
  can_configure_log_retention_period = true
  can_view_logs                   = true
}
```

#### Log Operations Token
```hcl
resource "instana_api_token" "log_operations" {
  name                = "log-operations-token"
  can_view_logs       = true
  can_delete_logs     = true
  can_view_log_volume = true
}
```

### Business Operations (BizOps) Tokens

#### BizOps Configuration Token
```hcl
resource "instana_api_token" "bizops_config" {
  name                 = "bizops-config-token"
  can_configure_bizops = true
}
```

#### BizOps Read-Only Token
```hcl
resource "instana_api_token" "bizops_readonly" {
  name                          = "bizops-readonly-token"
  can_view_business_processes   = true
  can_view_business_process_details = true
  can_view_business_activities  = true
  can_view_biz_alerts           = true
}
```

### Service Level Management Tokens

#### SLO Configuration Token
```hcl
resource "instana_api_token" "slo_config" {
  name                                           = "slo-configuration-token"
  can_configure_service_levels                   = true
  can_configure_service_level_smart_alerts       = true
  can_configure_service_level_correction_windows = true
}
```

### Debugging and Diagnostics Tokens

#### Diagnostics Token
```hcl
resource "instana_api_token" "diagnostics" {
  name                   = "diagnostics-token"
  can_create_heap_dump   = true
  can_create_thread_dump = true
  can_view_trace_details = true
  can_configure_subtraces = true
}
```

### Dashboard Management Tokens

#### Dashboard Creation Token
```hcl
resource "instana_api_token" "dashboard_create" {
  name                            = "dashboard-creation-token"
  can_create_public_custom_dashboards = true
}
```

#### Dashboard Edit Token
```hcl
resource "instana_api_token" "dashboard_edit" {
  name                                  = "dashboard-edit-token"
  can_edit_all_accessible_custom_dashboards = true
  can_create_public_custom_dashboards       = true
}
```

### Scope-Limited Tokens

#### Application-Scoped Token
```hcl
resource "instana_api_token" "app_scoped" {
  name                       = "application-scoped-token"
  can_configure_applications = true
  limited_applications_scope = true
}
```

#### Infrastructure-Scoped Token
```hcl
resource "instana_api_token" "infra_scoped" {
  name                        = "infrastructure-scoped-token"
  can_configure_agents        = true
  limited_infrastructure_scope = true
}
```

#### Kubernetes-Scoped Token
```hcl
resource "instana_api_token" "k8s_scoped" {
  name                   = "kubernetes-scoped-token"
  can_configure_agents   = true
  limited_kubernetes_scope = true
}
```

#### Multi-Scope Token
```hcl
resource "instana_api_token" "multi_scope" {
  name                        = "multi-scope-token"
  can_configure_applications  = true
  can_configure_agents        = true
  limited_applications_scope  = true
  limited_infrastructure_scope = true
  limited_kubernetes_scope    = true
  limited_websites_scope      = true
}
```

### AI and LLM Tokens

#### AI Configuration Token
```hcl
resource "instana_api_token" "ai_config" {
  name                  = "ai-configuration-token"
  can_configure_llm     = true
  can_configure_ai_agents = true
  limited_ai_gateway_scope = true
}
```

### Administrative Tokens

#### Full Admin Token
```hcl
resource "instana_api_token" "admin" {
  name                                           = "admin-token"
  can_configure_service_mapping                  = true
  can_configure_eum_applications                 = true
  can_configure_mobile_app_monitoring            = true
  can_configure_users                            = true
  can_install_new_agents                         = true
  can_configure_integrations                     = true
  can_configure_events_and_alerts                = true
  can_configure_maintenance_windows              = true
  can_configure_application_smart_alerts         = true
  can_configure_website_smart_alerts             = true
  can_configure_mobile_app_smart_alerts          = true
  can_configure_api_tokens                       = true
  can_configure_agent_run_mode                   = true
  can_view_audit_log                             = true
  can_configure_agents                           = true
  can_configure_authentication_methods           = true
  can_configure_applications                     = true
  can_configure_teams                            = true
  can_configure_releases                         = true
  can_configure_log_management                   = true
  can_create_public_custom_dashboards            = true
  can_view_logs                                  = true
  can_view_trace_details                         = true
  can_configure_session_settings                 = true
  can_configure_global_alert_payload             = true
  can_configure_global_application_smart_alerts  = true
  can_configure_global_synthetic_smart_alerts    = true
  can_configure_global_infra_smart_alerts        = true
  can_configure_global_log_smart_alerts          = true
  can_view_account_and_billing_information       = true
  can_edit_all_accessible_custom_dashboards      = true
}
```

#### Security Auditor Token
```hcl
resource "instana_api_token" "security_auditor" {
  name                                     = "security-auditor-token"
  can_view_audit_log                       = true
  can_view_account_and_billing_information = true
  can_view_logs                            = true
  can_view_trace_details                   = true
}
```

### Platform-Specific Tokens

#### VMware vSphere Token
```hcl
resource "instana_api_token" "vsphere" {
  name                 = "vsphere-token"
  can_configure_agents = true
  limited_vsphere_scope = true
}
```

#### OpenStack Token
```hcl
resource "instana_api_token" "openstack" {
  name                   = "openstack-token"
  can_configure_agents   = true
  limited_openstack_scope = true
}
```

#### Nutanix Token
```hcl
resource "instana_api_token" "nutanix" {
  name                  = "nutanix-token"
  can_configure_agents  = true
  limited_nutanix_scope = true
}
```

### Complete Production Token Example

```hcl
resource "instana_api_token" "production_full" {
  name = "production-full-access-token"
  
  # Application and Service Configuration
  can_configure_service_mapping       = true
  can_configure_eum_applications      = true
  can_configure_mobile_app_monitoring = true
  can_configure_applications          = true
  
  # Agent Management
  can_install_new_agents       = true
  can_configure_agents         = true
  can_configure_agent_run_mode = true
  
  # Alerts and Events
  can_configure_events_and_alerts        = true
  can_configure_application_smart_alerts = true
  can_configure_website_smart_alerts     = true
  can_configure_maintenance_windows      = true
  can_invoke_alert_channel               = true
  
  # Integrations and Automation
  can_configure_integrations         = true
  can_configure_automation_actions   = true
  can_configure_automation_policies  = true
  can_run_automation_actions         = true
  
  # Synthetic Monitoring
  can_configure_synthetic_tests       = true
  can_configure_synthetic_locations   = true
  can_view_synthetic_tests            = true
  can_view_synthetic_test_results     = true
  
  # Log Management
  can_configure_log_management = true
  can_view_logs                = true
  can_view_log_volume          = true
  
  # Service Levels
  can_configure_service_levels                   = true
  can_configure_service_level_smart_alerts       = true
  can_configure_service_level_correction_windows = true
  
  # Dashboards
  can_create_public_custom_dashboards       = true
  can_edit_all_accessible_custom_dashboards = true
  
  # Observability
  can_view_trace_details = true
  can_configure_releases = true
  
  # Scope Limitations
  limited_applications_scope   = false
  limited_infrastructure_scope = false
  limited_kubernetes_scope     = false
}
```

## Argument Reference

### Computed Attributes

* `id` - (Computed) The unique identifier of the API token
* `access_granting_token` - (Computed) The token value used in the Authorization header for API authentication
* `internal_id` - (Computed) The internal identifier used by Instana

### Required Attributes

* `name` - (Required) The name of the API token

### Core Permissions

All permission attributes are optional and default to `false`:

* `can_configure_service_mapping` - Enables permission to configure service mappings
* `can_configure_eum_applications` - Enables permission to configure End User Monitoring (EUM) applications
* `can_configure_mobile_app_monitoring` - Enables permission to configure mobile app monitoring
* `can_configure_users` - Enables permission to configure users
* `can_install_new_agents` - Enables permission to install new agents
* `can_configure_integrations` - Enables permission to configure integrations
* `can_configure_events_and_alerts` - Enables permission to configure events and alerts
* `can_configure_maintenance_windows` - Enables permission to configure maintenance windows
* `can_configure_application_smart_alerts` - Enables permission to configure application smart alerts
* `can_configure_website_smart_alerts` - Enables permission to configure website smart alerts
* `can_configure_mobile_app_smart_alerts` - Enables permission to configure mobile app smart alerts
* `can_configure_api_tokens` - Enables permission to configure API tokens
* `can_configure_agent_run_mode` - Enables permission to configure agent run mode
* `can_view_audit_log` - Enables permission to view audit logs
* `can_configure_agents` - Enables permission to configure agents
* `can_configure_authentication_methods` - Enables permission to configure authentication methods
* `can_configure_applications` - Enables permission to configure applications
* `can_configure_teams` - Enables permission to configure teams (groups)
* `can_configure_releases` - Enables permission to configure releases
* `can_configure_log_management` - Enables permission to configure log management
* `can_create_public_custom_dashboards` - Enables permission to create public custom dashboards
* `can_view_logs` - Enables permission to view logs
* `can_view_trace_details` - Enables permission to view trace details
* `can_configure_session_settings` - Enables permission to configure session settings
* `can_configure_global_alert_payload` - Enables permission to configure global alert payload
* `can_configure_global_application_smart_alerts` - Enables permission to configure global application smart alerts
* `can_configure_global_synthetic_smart_alerts` - Enables permission to configure global synthetic smart alerts
* `can_configure_global_infra_smart_alerts` - Enables permission to configure global infrastructure smart alerts
* `can_configure_global_log_smart_alerts` - Enables permission to configure global log smart alerts
* `can_view_account_and_billing_information` - Enables permission to view account and billing information
* `can_edit_all_accessible_custom_dashboards` - Enables permission to edit all accessible custom dashboards

### Scope Limitations

All scope limitation attributes are optional and default to `false`:

* `limited_applications_scope` - Limits the scope to applications
* `limited_biz_ops_scope` - Limits the scope to business operations
* `limited_websites_scope` - Limits the scope to websites
* `limited_kubernetes_scope` - Limits the scope to Kubernetes
* `limited_mobile_apps_scope` - Limits the scope to mobile apps
* `limited_infrastructure_scope` - Limits the scope to infrastructure
* `limited_synthetics_scope` - Limits the scope to synthetics
* `limited_vsphere_scope` - Limits the scope to VMware vSphere
* `limited_phmc_scope` - Limits the scope to PowerHMC
* `limited_pvc_scope` - Limits the scope to PowerVC
* `limited_zhmc_scope` - Limits the scope to z/HMC
* `limited_pcf_scope` - Limits the scope to Pivotal Cloud Foundry
* `limited_openstack_scope` - Limits the scope to OpenStack
* `limited_automation_scope` - Limits the scope to automation
* `limited_logs_scope` - Limits the scope to logs
* `limited_nutanix_scope` - Limits the scope to Nutanix
* `limited_xen_server_scope` - Limits the scope to Xen Server
* `limited_windows_hypervisor_scope` - Limits the scope to Windows Hypervisor
* `limited_alert_channels_scope` - Limits the scope to alert channels
* `limited_linux_kvm_hypervisor_scope` - Limits the scope to Linux KVM Hypervisor
* `limited_service_level_scope` - Limits the scope to service levels
* `limited_ai_gateway_scope` - Limits the scope to AI Gateway

### Additional Permissions

All additional permission attributes are optional and default to `false`:

* `can_configure_personal_api_tokens` - Enables permission to configure personal API tokens
* `can_configure_database_management` - Enables permission to configure database management
* `can_configure_automation_actions` - Enables permission to configure automation actions
* `can_configure_automation_policies` - Enables permission to configure automation policies
* `can_run_automation_actions` - Enables permission to run automation actions
* `can_delete_automation_action_history` - Enables permission to delete automation action history
* `can_configure_synthetic_tests` - Enables permission to configure synthetic tests
* `can_configure_synthetic_locations` - Enables permission to configure synthetic locations
* `can_configure_synthetic_credentials` - Enables permission to configure synthetic credentials
* `can_view_synthetic_tests` - Enables permission to view synthetic tests
* `can_view_synthetic_locations` - Enables permission to view synthetic locations
* `can_view_synthetic_test_results` - Enables permission to view synthetic test results
* `can_use_synthetic_credentials` - Enables permission to use synthetic credentials
* `can_configure_bizops` - Enables permission to configure business operations
* `can_view_business_processes` - Enables permission to view business processes
* `can_view_business_process_details` - Enables permission to view business process details
* `can_view_business_activities` - Enables permission to view business activities
* `can_view_biz_alerts` - Enables permission to view business alerts
* `can_delete_logs` - Enables permission to delete logs
* `can_create_heap_dump` - Enables permission to create heap dumps
* `can_create_thread_dump` - Enables permission to create thread dumps
* `can_manually_close_issue` - Enables permission to manually close issues
* `can_view_log_volume` - Enables permission to view log volume
* `can_configure_log_retention_period` - Enables permission to configure log retention period
* `can_configure_subtraces` - Enables permission to configure subtraces
* `can_invoke_alert_channel` - Enables permission to invoke alert channels
* `can_configure_llm` - Enables permission to configure Large Language Models (LLM)
* `can_configure_ai_agents` - Enables permission to configure AI agents
* `can_configure_apdex` - Enables permission to configure Apdex thresholds
* `can_configure_service_level_correction_windows` - Enables permission to configure service level correction windows
* `can_configure_service_level_smart_alerts` - Enables permission to configure service level smart alerts
* `can_configure_service_levels` - Enables permission to configure service levels (SLOs/SLIs)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the API token
* `access_granting_token` - The token value to use for API authentication
* `internal_id` - The internal identifier used by Instana

## Import

API Tokens can be imported using the `internal_id`, e.g.:

```bash
$ terraform import instana_api_token.my_token 60845e4e5e6b9cf8fc2868da
```

## Notes

### Token Security

- The `access_granting_token` is sensitive and should be stored securely
- Use Terraform outputs with `sensitive = true` when exposing token values
- Consider using a secrets management system for production tokens
- Rotate tokens regularly according to your security policies

### Permission Model

The Instana API token permission model follows a principle of least privilege:

- All permissions default to `false`
- Grant only the permissions required for the token's intended use
- Use scope limitations to further restrict token access
- Combine permissions logically based on use cases

### Scope Limitations

Scope limitations work in conjunction with permissions:

- Setting a scope limitation restricts the token to specific resource types
- Permissions still apply within the limited scope
- Multiple scope limitations can be combined
- Scope limitations provide an additional security layer

### Best Practices

1. **Use Descriptive Names**: Name tokens based on their purpose (e.g., "ci-cd-deployment-token")
2. **Minimal Permissions**: Grant only necessary permissions
3. **Scope Appropriately**: Use scope limitations when possible
4. **Separate Tokens**: Create separate tokens for different purposes
5. **Regular Audits**: Review and rotate tokens periodically
6. **Secure Storage**: Never commit tokens to version control
