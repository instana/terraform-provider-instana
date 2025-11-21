# instana Provider

Terraform provider implementation of the Instana Web REST API. The provider can be used to configure different
assents in Instana. The provider is aligned with the REST API and links to the endpoint is provided for each 
resource. 

**NOTE:** Starting with version 0.6.0 Terraform version 0.12.x or later is required.

## Supported Resources:

* Application Settings
  * Application Configuration - `instana_application_config`
  * Application Alert Configuration - `instana_application_alert_config`
  * Global Application Alert Configuration - `instana_global_application_alert_config`
* Automation
  * Automation Action - `instana_automation_action`
  * Automation Policy - `instana_automation_policy`
* Custom Dashboard - `instana_custom_dashboard`
* Event Settings
  * Custom Event Specification - `instana_custom_event_specification`
  * Alerting Channels - `instana_alerting_channel`
  * Alerting Config - `instana_alerting_config`
* Infrastructure Monitoring
  * Infrastructure Alert Config - `instana_infra_alert_config`
* Log monitoring
  * Log Alert Config - `instana_log_alert_config`
* Service Levels
  * Service Level Objective Config - `instana_slo_config`
  * Service Level Objective (SLO) Alert Config - `instana_slo_alert_config`
* Settings
  * API Tokens - `instana_api_token`
  * Groups - `instana_rbac_group`
  * Team - `instana_rbac_team`
* SLI Settings
  * SLI Config - `instana_sli_config`
* Synthetic Settings
  * Synthetic Test - `instana_synthetic_test`
  * Synthetic Alert Config - `instana_synthetic_alert_config`
* Website Monitoring
  * Website Monitoring Config - `instana_website_monitoring_config`
  * Website Alert Config - `instana_website_alert_config`

## Supported Data Source:

* Automation
  * Automation Action - `instana_automation_action`
* Event Settings
  * Alerting Channel - `instana_alerting_channel`
  * Builtin Event Specifications - `instana_builtin_event_spec`
  * Custom Event Specifications - `instana_custom_event_spec`
* Host Agent - `instana_host_agents`
* Synthetic Settings
  * Synthetic Location - `instana_synthetic_location`

## Example Usage

```hcl
provider "instana" {
  api_token = "secure-api-token"  
  endpoint = "<tenant>-<org>.instana.io"
  tls_skip_verify     = false
}
```

## Argument Reference

* `api_token` - Required - The API token which is created in the Settings area of Instana for remote access through 
the REST API. You have to make sure that you assign the proper permissions for this token to configure the desired 
resources with this provider. E.g. when User Roles should be provisioned by terraform using this provider implementation 
then the permission 'Access role configuration' must be activated. (Defaults to the environment variable `INSTANA_API_TOKEN`).
* `endpoint` - Required - The endpoint of the instana backend. For SaaS the endpoint URL has the pattern
`<tenant>-<organization>.instana.io`. For onPremise installation the endpoint URL depends on your local setup. (Defaults to the environment variable `INSTANA_ENDPOINT`).
* `tls_skip_verify` - `Òptional` - Default `false` - If set to true, TLS verification will be skipped when calling Instana API

## Import support

All resources of the terraform provider instana support resource import.
