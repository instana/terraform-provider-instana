# API Token Resource

Management of API Tokens for programmatic access to Instana.

API Documentation: <https://instana.github.io/openapi/#operation/getApiToken

 **⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)**

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- All attributes are now top-level attributes (no nested blocks)
- Boolean attributes now use explicit `true`/`false` values with defaults
- The `id`, `access_granting_token`, and `internal_id` attributes are computed
- All permission and scope attributes have default values of `false`
- Attribute syntax remains the same (key = value), but schema validation is stricter

#### OLD (v5.x) Syntax:

```hcl
resource "instana_api_token" "example" {
  name                          = "my-token"
  can_configure_service_mapping = true
  # Omitted attributes defaulted to false
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_api_token" "example" {
  name                          = "my-token"
  can_configure_service_mapping = true
  # Omitted attributes explicitly default to false
  # All boolean attributes now have computed defaults
}
```

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.

## Example Usage

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

####  Access Token with More Permissions

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

## Generating Configuration from Existing Resources

If you have already created a API token in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your API token in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_api_token.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual API token ID from Instana

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
