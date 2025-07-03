# API Token Resource

Management of API Tokens.

API Documentation: <https://instana.github.io/openapi/#operation/getApiToken>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

## Example Usage

```hcl
resource "instana_api_token" "example" {
  name                                           = "name"
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
  can_configure_service_level_indicators         = true
  can_configure_global_alert_payload             = true
  can_configure_global_application_smart_alerts  = true
  can_configure_global_synthetic_smart_alerts    = true
  can_configure_global_infra_smart_alerts        = true
  can_configure_global_log_smart_alerts          = true
  can_view_account_and_billing_information       = true
  can_edit_all_accessible_custom_dashboards      = true
  
  # Scope limitations
  limited_applications_scope                     = true
  limited_biz_ops_scope                          = true
  limited_websites_scope                         = true
  limited_kubernetes_scope                       = true
  limited_mobile_apps_scope                      = true
  limited_infrastructure_scope                   = true
  limited_synthetics_scope                       = true
  limited_vsphere_scope                          = true
  limited_phmc_scope                             = true
  limited_pvc_scope                              = true
  limited_zhmc_scope                             = true
  limited_pcf_scope                              = true
  limited_openstack_scope                        = true
  limited_automation_scope                       = true
  limited_logs_scope                             = true
  limited_nutanix_scope                          = true
  limited_xen_server_scope                       = true
  limited_windows_hypervisor_scope               = true
  limited_alert_channels_scope                   = true
  limited_linux_kvm_hypervisor_scope             = true
  
  # Additional permissions
  can_configure_personal_api_tokens              = true
  can_configure_database_management              = true
  can_configure_automation_actions               = true
  can_configure_automation_policies              = true
  can_run_automation_actions                     = true
  can_delete_automation_action_history           = true
  can_configure_synthetic_tests                  = true
  can_configure_synthetic_locations              = true
  can_configure_synthetic_credentials            = true
  can_view_synthetic_tests                       = true
  can_view_synthetic_locations                   = true
  can_view_synthetic_test_results                = true
  can_use_synthetic_credentials                  = true
  can_configure_bizops                           = true
  can_view_business_processes                    = true
  can_view_business_process_details              = true
  can_view_business_activities                   = true
  can_view_biz_alerts                            = true
  can_delete_logs                                = true
  can_create_heap_dump                           = true
  can_create_thread_dump                         = true
  can_manually_close_issue                       = true
  can_view_log_volume                            = true
  can_configure_log_retention_period             = true
  can_configure_subtraces                        = true
  can_invoke_alert_channel                       = true
  can_configure_llm                              = true
}
```

## Argument Reference

* `access_granting_token`-  Calculated - The token used for the api Client used in the Authorization header to authenticate the client
* `name` - Required - the name of the alerting channel
* `can_configure_service_mapping` - Optional - default false - enables permission to configure service mappings
* `can_configure_eum_applications` - Optional - default false - enables permission to configure EUM applications
* `can_configure_mobile_app_monitoring` - Optional - default false - enables permission to configure mobile app monitoring
* `can_configure_users` - Optional - default false - enables permission to configure users
* `can_install_new_agents` - Optional - default false - enables permission to install new agents
* `can_configure_integrations` - Optional - default false - enables permission to configure integrations
* `can_configure_events_and_alerts` - Optional - default false - enables permission to configure Events and Alerts
* `can_configure_maintenance_windows` - Optional - default false - enables permission to configure Maintenance Windows
* `can_configure_application_smart_alerts` - Optional - default false - enables permission to configure Application Smart Alerts
* `can_configure_website_smart_alerts` - Optional - default false - enables permission to configure Website Smart Alerts
* `can_configure_mobile_app_smart_alerts` - Optional - default false - enables permission to configure MobileApp Smart Alerts
* `can_configure_api_tokens` - Optional - default false - enables permission to configure api tokes
* `can_configure_agent_run_mode` - Optional - default false - enables permission to configure agent run mode
* `can_view_audit_log` - Optional - default false - enables permission to view audit logs
* `can_configure_agents` - Optional - default false - enables permission to configure agents
* `can_configure_authentication_methods` - Optional - default false - enables permission to configure authentication methods
* `can_configure_applications` - Optional - default false - enables permission to configure applications
* `can_configure_teams` - Optional - default false - enables permission to configure teams (groups)
* `can_configure_releases` - Optional - default false - enables permission to configure releases
* `can_configure_log_management` - Optional - default false - enables permission to configure log management
* `can_create_public_custom_dashboards` - Optional - default false - enables permission to create public custom dashboards
* `can_view_logs` - Optional - default false - enables permission to view logs
* `can_view_trace_details` - Optional - default false - enables permission to view trace details
* `can_configure_session_settings` - Optional - default false - enables permission to configure session settings
* `can_configure_service_level_indicators` - Optional - default false - enables permission to configure service level indicators
* `can_configure_global_alert_payload` - Optional - default false - enables permission to configure global alert payload
* `can_configure_global_application_smart_alerts` - Optional - default false - enables permission to configure global Application Smart Alerts
* `can_configure_global_synthetic_smart_alerts` - Optional - default false - enables permission to configure global Synthetic Smart Alerts
* `can_configure_global_infra_smart_alerts` - Optional - default false - enables permission to configure global Infrastructure Smart Alerts
* `can_configure_global_log_smart_alerts` - Optional - default false - enables permission to configure global Log Smart Alerts
* `can_view_account_and_billing_information` - Optional - default false - enables permission to view account and billing information
* `can_edit_all_accessible_custom_dashboards` - Optional - default false - enables permission to edit all accessible custom dashboards

### Scope Limitations
* `limited_applications_scope` - Optional - default false - limits the scope to applications
* `limited_biz_ops_scope` - Optional - default false - limits the scope to business operations
* `limited_websites_scope` - Optional - default false - limits the scope to websites
* `limited_kubernetes_scope` - Optional - default false - limits the scope to kubernetes
* `limited_mobile_apps_scope` - Optional - default false - limits the scope to mobile apps
* `limited_infrastructure_scope` - Optional - default false - limits the scope to infrastructure
* `limited_synthetics_scope` - Optional - default false - limits the scope to synthetics
* `limited_vsphere_scope` - Optional - default false - limits the scope to vsphere
* `limited_phmc_scope` - Optional - default false - limits the scope to phmc
* `limited_pvc_scope` - Optional - default false - limits the scope to pvc
* `limited_zhmc_scope` - Optional - default false - limits the scope to zhmc
* `limited_pcf_scope` - Optional - default false - limits the scope to pcf
* `limited_openstack_scope` - Optional - default false - limits the scope to openstack
* `limited_automation_scope` - Optional - default false - limits the scope to automation
* `limited_logs_scope` - Optional - default false - limits the scope to logs
* `limited_nutanix_scope` - Optional - default false - limits the scope to nutanix
* `limited_xen_server_scope` - Optional - default false - limits the scope to xen server
* `limited_windows_hypervisor_scope` - Optional - default false - limits the scope to windows hypervisor
* `limited_alert_channels_scope` - Optional - default false - limits the scope to alert channels
* `limited_linux_kvm_hypervisor_scope` - Optional - default false - limits the scope to linux kvm hypervisor

### Additional Permissions
* `can_configure_personal_api_tokens` - Optional - default false - enables permission to configure personal API tokens
* `can_configure_database_management` - Optional - default false - enables permission to configure database management
* `can_configure_automation_actions` - Optional - default false - enables permission to configure automation actions
* `can_configure_automation_policies` - Optional - default false - enables permission to configure automation policies
* `can_run_automation_actions` - Optional - default false - enables permission to run automation actions
* `can_delete_automation_action_history` - Optional - default false - enables permission to delete automation action history
* `can_configure_synthetic_tests` - Optional - default false - enables permission to configure synthetic tests
* `can_configure_synthetic_locations` - Optional - default false - enables permission to configure synthetic locations
* `can_configure_synthetic_credentials` - Optional - default false - enables permission to configure synthetic credentials
* `can_view_synthetic_tests` - Optional - default false - enables permission to view synthetic tests
* `can_view_synthetic_locations` - Optional - default false - enables permission to view synthetic locations
* `can_view_synthetic_test_results` - Optional - default false - enables permission to view synthetic test results
* `can_use_synthetic_credentials` - Optional - default false - enables permission to use synthetic credentials
* `can_configure_bizops` - Optional - default false - enables permission to configure business operations
* `can_view_business_processes` - Optional - default false - enables permission to view business processes
* `can_view_business_process_details` - Optional - default false - enables permission to view business process details
* `can_view_business_activities` - Optional - default false - enables permission to view business activities
* `can_view_biz_alerts` - Optional - default false - enables permission to view business alerts
* `can_delete_logs` - Optional - default false - enables permission to delete logs
* `can_create_heap_dump` - Optional - default false - enables permission to create heap dumps
* `can_create_thread_dump` - Optional - default false - enables permission to create thread dumps
* `can_manually_close_issue` - Optional - default false - enables permission to manually close issues
* `can_view_log_volume` - Optional - default false - enables permission to view log volume
* `can_configure_log_retention_period` - Optional - default false - enables permission to configure log retention period
* `can_configure_subtraces` - Optional - default false - enables permission to configure subtraces
* `can_invoke_alert_channel` - Optional - default false - enables permission to invoke alert channels
* `can_configure_llm` - Optional - default false - enables permission to configure LLM

## Import

API Tokens can be imported using the `internal_id`, e.g.:

```
$ terraform import instana_api_token.my_token 60845e4e5e6b9cf8fc2868da
```
