# Alerting Channel Resource

Alerting channel configuration for notifications to a specified target channel.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels

## ⚠️ BREAKING CHANGES - Plugin Framework Migration (v6.0.0)

 **This resource has been migrated from Terraform SDK v2 to the Terraform Plugin Framework**. The schema has transitioned from **block structure to attribute format**.While the basic structure remains similar, there are important syntax changes for block structure.
 
## Migration Guide (v5 to v6)

For detailed migration instructions and examples, see the [Plugin Framework Migration Guide](https://github.com/instana/terraform-provider-instana/blob/main/PLUGIN-FRAMEWORK-MIGRATION.md).

### Syntax Changes Overview

- All channel configurations (email, slack, webhook, etc.) are now **single nested attributes** instead of blocks
- Use `attribute = { ... }` syntax instead of `attribute { ... }` block syntax
- Channel configurations must use the equals sign (`=`) before the opening brace
- All nested configurations within channels follow the same attribute pattern
- The `id` attribute is now computed and uses a plan modifier for state management

#### OLD (v5.x) Syntax:

```hcl
resource "instana_alerting_channel" "example" {
  name = "my-channel"
  email {
    emails = ["user@example.com"]
  }
}
```

#### NEW (v6.x) Syntax:

```hcl
resource "instana_alerting_channel" "example" {
  name = "my-channel"
  email = {
    emails = ["user@example.com"]
  }
}
```

Please update your Terraform configurations to use the new attribute-based syntax.


## Example Usage

### Email Alerting Channel

#### Basic Email Configuration
```hcl
resource "instana_alerting_channel" "email_basic" {
  name = "basic-email-alerts"
  
  email = {
    emails = ["alerts@example.com"]
  }
}
```

#### Multiple Email Recipients
```hcl
resource "instana_alerting_channel" "email_multiple" {
  name = "team-email-alerts"
  
  email = {
    emails = [
      "team-lead@example.com",
      "on-call@example.com",
      "devops@example.com"
    ]
  }
}
```

### Google Chat Alerting Channel

#### Basic Google Chat Configuration
```hcl
resource "instana_alerting_channel" "google_chat_basic" {
  name = "google-chat-alerts"
  
  google_chat = {
    webhook_url = "https://chat.googleapis.com/v1/spaces/AAAA/messages?key=YOUR_KEY&token=YOUR_TOKEN" # Replace with your own value
  }
}
```

### Office 365 Alerting Channel

#### Basic Office 365 Configuration
```hcl
resource "instana_alerting_channel" "office365_basic" {
  name = "office365-alerts"
  
  office_365 = {
    webhook_url = "https://outlook.office.com/webhook/a1b2c3d4-e5f6-7890-abcd-ef1234567890@tenant.onmicrosoft.com/IncomingWebhook/xyz123abc456/def789ghi012" # Replace with your own value
  }
}
```

### OpsGenie Alerting Channel

#### Basic OpsGenie Configuration
```hcl
resource "instana_alerting_channel" "opsgenie_eu" {
  name = "opsgenie-eu-alerts"
  
  ops_genie = {
    api_key = "a1b2c3d4-e5f6-7890-abcd" # Replace with your own value
    tags    = ["instana", "production"]
    region  = "EU"
  }
}
```

#### OpsGenie with Multiple Tags
```hcl
resource "instana_alerting_channel" "opsgenie_detailed" {
  name = "opsgenie-detailed-alerts"
  
  ops_genie = {
    api_key = "f9e8d7c6-b5a4-3210-9876" # Replace with your own value
    tags = [
      "environment:production",
      "team:platform",
      "severity:high",
      "source:instana"
    ]
    region = "EU"
  }
}
```

### PagerDuty Alerting Channel

#### Basic PagerDuty Configuration
```hcl
resource "instana_alerting_channel" "pagerduty_basic" {
  name = "pagerduty-alerts"
  
  pager_duty = {
    service_integration_key = "a1b2c3d4e5f67890abcde" # Replace with your own value
  }
}
```


### Slack Alerting Channel

#### Basic Slack Configuration
```hcl
resource "instana_alerting_channel" "slack_channel" {
  name                    = "slack alert"
  slack = {
    webhook_url = "https://hooks.slack.com/services/T01234567/B01234567/abcdefghijklmnopqrstuvwx" # Replace with your own value
  }
}
```

#### Slack with Custom Icon
```hcl
resource "instana_alerting_channel" "slack_custom_icon" {
  name = "slack-custom-alerts"
  
  slack = {
    webhook_url = "https://hooks.slack.com/services/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" # Replace with your own value
    icon_url    = "https://example.com/instana-icon.png" # Replace with your own value
  }
}
```

#### Slack with Specific Channel
```hcl
resource "instana_alerting_channel" "slack_specific_channel" {
  name = "slack-production-alerts"
  
  slack = {
    webhook_url = "https://hooks.slack.com/services/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" # Replace with your own value
    icon_url    = "https://example.com/alert-icon.png" # Replace with your own value
    channel     = "#production-alerts"
  }
}
```

### Splunk Alerting Channel

#### Basic Splunk Configuration
```hcl
resource "instana_alerting_channel" "splunk_basic" {
  name = "splunk-alerts"
  
  splunk = {
    url   = "https://splunk.example.com:8088/services/collector" # Replace with your own value
    token = "a1b2c3d4-e5f6-7890-abcd" # Replace with your own value
  }
}
```

### VictorOps Alerting Channel

#### Basic VictorOps Configuration
```hcl
resource "instana_alerting_channel" "victorops_basic" {
  name = "victorops-alerts"
  
  victor_ops = {
    api_key     = "f9e8d7c6-b5a4-3210-9876-543210fedcba" # Replace with your own value
    routing_key = "instana-alerts" # Replace with your own value
  }
}
```

#### VictorOps with Team-Specific Routing
```hcl
resource "instana_alerting_channel" "victorops_platform_team" {
  name = "victorops-platform-team"
  
  victor_ops = {
    api_key     = "12345678-90ab-cdef-1234-567890abcdef" # Replace with your own value
    routing_key = "platform-team"
  }
}

resource "instana_alerting_channel" "victorops_database_team" {
  name = "victorops-database-team"
  
  victor_ops = {
    api_key     = "abcdef12-3456-7890-abcd-ef1234567890" # Replace with your own value
    routing_key = "database-team"
  }
}
```

### Webhook Alerting Channel

#### Basic Webhook Configuration
```hcl
resource "instana_alerting_channel" "webhook_channel" {
  name                    = "webhook-channel"
  webhook = {
    http_headers = {}
    webhook_urls = ["https://api.example.com/instana/alerts"] # Replace with your own value
  }
}
```

#### Multiple Webhook URLs
```hcl
resource "instana_alerting_channel" "webhook_multiple" {
  name = "webhook-multiple-endpoints"
  
  webhook = {
    http_headers = {}
    webhook_urls = [ # Replace with your own value
      "https://api.example.com/instana/alerts",
      "https://backup-api.example.com/alerts",
      "https://monitoring.example.com/webhooks/instana"
    ]
  }
}
```

#### Webhook with Custom Headers
```hcl
resource "instana_alerting_channel" "webhook_with_headers" {
  name = "webhook-authenticated"
  
  webhook = {
    webhook_urls = ["https://api.example.com/instana/alerts"] # Replace with your own value
    http_headers = {
      "Authorization" = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6Ik" # Replace with your own value
      "X-API-Key"     = "sk_live_51A1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0" # Replace with your own value
      "Content-Type"  = "application/json"
    }
  }
}
```

### ServiceNow Alerting Channel

#### Basic ServiceNow Configuration
```hcl
resource "instana_alerting_channel" "servicenow_alerts" {
  name               = "servicenow-alerts"
  service_now = {
    service_now_url      = "https://your-instance.service-now.com" # Replace with your own value
    username             = "instana_integration_user" # Replace with your own value
    password             = "P@ssw0rd!2024#SecurePass" # Replace with your own value
    auto_close_incidents = false
  }
}
```

#### ServiceNow with Auto-Close
```hcl
resource "instana_alerting_channel" "servicenow_autoclose" {
  name = "servicenow-autoclose-alerts"
  
  service_now = {
    service_now_url      = "https://your-instance.service-now.com" # Replace with your own value
    username             = "snow_admin_user" # Replace with your own value
    password             = "SecureP@ss123!Winter2024" # Replace with your own value
    auto_close_incidents = true
  }
}
```

### ServiceNow Enhanced (ITSM) Alerting Channel

#### Basic ServiceNow Enhanced Configuration
```hcl
resource "instana_alerting_channel" "service_now_channel" {
  name               = "Service-now-App"
  service_now_application = {
    auto_close_incidents               = true
    enable_send_instana_notes          = true
    enable_send_service_now_activities = true
    enable_send_service_now_work_notes = true
    instana_url                        = "https://test-instana.pink.instana.rocks/" # Replace with your own value
    manually_closed_incidents          = true
    username                           = "servicenow_itsm_user" # Replace with your own value
    password                           = "ITSM@Pass2024!Secure" # Replace with your own value
    resolution_of_incident             = true
    service_now_url                    = "https://service-now.com" # Replace with your own value
    snow_status_on_close_event         = -1
    tenant                             = "instana"
    unit                               = "test"
  }
}

```

### Prometheus Webhook Alerting Channel

#### Prometheus Webhook Configuration
```hcl
resource "instana_alerting_channel" "prometheus_channel" {
  name         = "prometheus-webhook-receiver"
  prometheus_webhook = {
    receiver    = "instana-alerts"
    webhook_url = "https://prometheus.example.com/api/v1/alerts" # Replace with your own value
  }
}
```

### Webex Teams Webhook Alerting Channel

#### Basic Webex Teams Configuration
```hcl
resource "instana_alerting_channel" "webex_basic" {
  name = "webex-teams-alerts"
  
  webex_teams_webhook = {
    webhook_url = "https://webexapis.com/v1/webhooks/incoming/Y2lzY29zcGF" # Replace with your own value
  }
}
```

### Watson AIOps Webhook Alerting Channel

#### Basic Watson AIOps Configuration
```hcl
resource "instana_alerting_channel" "watson_basic" {
  name = "watson-aiops-alerts"
  
  watson_aiops_webhook = {
    webhook_url = "https://watson-aiops.example.com/webhook" # Replace with your own value
  }
}
```

#### Watson AIOps with Custom Headers
```hcl
resource "instana_alerting_channel" "watson_with_headers" {
  name = "watson-aiops-authenticated"
  
  watson_aiops_webhook = {
    webhook_url = "https://watson-aiops.example.com/webhook" # Replace with your own value
    http_headers = [
      "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9", # Replace with your own value
      "X-API-Key: wsk_prod_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6", # Replace with your own value
      "X-Environment: production"
    ]
  }
}
```

### Slack App (Bidirectional) Alerting Channel

#### Basic Slack App Configuration
```hcl
resource "instana_alerting_channel" "slack_app_basic" {
  name = "slack-app-alerts"
  
  slack_app = {
    app_id       = "A01234567" # Replace with your own value
    team_id      = "T01234567" # Replace with your own value
    team_name    = "My Team" # Replace with your own value
    channel_id   = "C01234567" # Replace with your own value
    channel_name = "#alerts" # Replace with your own value
  }
}
```

#### Slack App with Emoji Rendering
```hcl
resource "instana_alerting_channel" "slack_app_with_emoji" {
  name = "slack-app-emoji-alerts"
  
  slack_app = {
    app_id          = "A01234567" # Replace with your own value
    team_id         = "T01234567" # Replace with your own value
    team_name       = "My Team" # Replace with your own value
    channel_id      = "C01234567" # Replace with your own value
    channel_name    = "#production-alerts" # Replace with your own value
    emoji_rendering = true
  }
}
```

### MS Teams App (Bidirectional) Alerting Channel

#### Basic MS Teams App Configuration
```hcl
resource "instana_alerting_channel" "ms_teams_bidirect" {
  name           = "MS Teams App Alert Channel"
  ms_teams_app = {
    channel_id   = "19:a1b2c3d4e5f6g7h8i9j0@thread.tacv2" # Replace with your own value
    channel_name = "Instana Alerts" # Replace with your own value
    service_url  = "https://smba.trafficmanager.net/amer/" # Replace with your own value
    team_id      = "a1b2c3d4-e5f6-7890-abcd-ef1234567890" # Replace with your own value
    team_name    = "Platform Engineering" # Replace with your own value
    tenant_id    = "f9e8d7c6-b5a4-3210-9876-543210fedcba" # Replace with your own value
    tenant_name  = "Acme Corporation" # Replace with your own value
  }
}
```

## Generating Configuration from Existing Resources

If you have already created an alerting channel in Instana and want to generate the Terraform configuration for it, you can use Terraform's import block feature with the `-generate-config-out` flag.

This approach is also helpful when you're unsure about the correct Terraform structure for a specific resource configuration. Simply create the resource in Instana first, then use this functionality to automatically generate the corresponding Terraform configuration.

### Steps to Generate Configuration:

1. **Get the Resource ID**: First, locate the ID of your alerting channel in Instana. You can find this in the Instana UI or via the API.

2. **Create an Import Block**: Create a new `.tf` file (e.g., `import.tf`) with an import block:

```hcl
import {
  to = instana_alerting_channel.example
  id = "resource_id"
}
```

Replace:
- `example` with your desired terraform block name
- `resource_id` with your actual alerting channel ID from Instana

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

* `id` - (Computed) The unique identifier of the alerting channel
* `name` - (Required) The name of the alerting channel

**Exactly one of the following channel types must be configured:**

* `email` - (Optional) Configuration of an email alerting channel - [Details](#email-channel-attributes)
* `google_chat` - (Optional) Configuration of a Google Chat alerting channel - [Details](#google-chat-channel-attributes)
* `office_365` - (Optional) Configuration of an Office 365 alerting channel - [Details](#office-365-channel-attributes)
* `ops_genie` - (Optional) Configuration of an OpsGenie alerting channel - [Details](#opsgenie-channel-attributes)
* `pager_duty` - (Optional) Configuration of a PagerDuty alerting channel - [Details](#pagerduty-channel-attributes)
* `slack` - (Optional) Configuration of a Slack alerting channel - [Details](#slack-channel-attributes)
* `splunk` - (Optional) Configuration of a Splunk alerting channel - [Details](#splunk-channel-attributes)
* `victor_ops` - (Optional) Configuration of a VictorOps alerting channel - [Details](#victorops-channel-attributes)
* `webhook` - (Optional) Configuration of a webhook alerting channel - [Details](#webhook-channel-attributes)
* `service_now` - (Optional) Configuration of a ServiceNow alerting channel - [Details](#servicenow-channel-attributes)
* `service_now_application` - (Optional) Configuration of a ServiceNow Enhanced (ITSM) alerting channel - [Details](#servicenow-enhanced-channel-attributes)
* `prometheus_webhook` - (Optional) Configuration of a Prometheus webhook alerting channel - [Details](#prometheus-webhook-channel-attributes)
* `webex_teams_webhook` - (Optional) Configuration of a Webex Teams webhook alerting channel - [Details](#webex-teams-webhook-channel-attributes)
* `watson_aiops_webhook` - (Optional) Configuration of a Watson AIOps webhook alerting channel - [Details](#watson-aiops-webhook-channel-attributes)
* `slack_app` - (Optional) Configuration of a Slack App (bidirectional) alerting channel - [Details](#slack-app-channel-attributes)
* `ms_teams_app` - (Optional) Configuration of a MS Teams App (bidirectional) alerting channel - [Details](#ms-teams-app-channel-attributes)

### Email Channel Attributes

* `emails` - (Required) Set of email addresses to send alerts to. Must contain at least one email address.

**Type:** `set(string)`

### Google Chat Channel Attributes

* `webhook_url` - (Required) The URL of the Google Chat webhook where alerts will be sent.

**Type:** `string`

### Office 365 Channel Attributes

* `webhook_url` - (Required) The URL of the Office 365 webhook where alerts will be sent.

**Type:** `string`

### OpsGenie Channel Attributes

* `api_key` - (Required) The API key for authentication with the OpsGenie API.
* `tags` - (Required) List of tags to attach to alerts in OpsGenie. Must contain at least one tag.
* `region` - (Required) The target OpsGenie region. Valid values: `EU`, `US`.

**Types:**
- `api_key`: `string`
- `tags`: `list(string)`
- `region`: `string`

### PagerDuty Channel Attributes

* `service_integration_key` - (Required) The service integration key for PagerDuty.

**Type:** `string`

### Slack Channel Attributes

* `webhook_url` - (Required) The URL of the Slack webhook to send alerts to.
* `icon_url` - (Optional) The URL to the icon that should be displayed in Slack messages.
* `channel` - (Optional) The target Slack channel where alerts should be posted (e.g., `#alerts`).

**Type:** `string` for all attributes

### Splunk Channel Attributes

* `url` - (Required) The target Splunk HTTP Event Collector (HEC) endpoint URL.
* `token` - (Required) The authentication token for the Splunk HEC API.

**Type:** `string` for both attributes

### VictorOps Channel Attributes

* `api_key` - (Required) The API key to authenticate with the VictorOps API.
* `routing_key` - (Required) The routing key used by VictorOps to route alerts to the desired target.

**Type:** `string` for both attributes

### Webhook Channel Attributes

* `webhook_urls` - (Required) Set of webhook URLs where alerts will be sent. Must contain at least one URL.
* `http_headers` - (Optional) Map of additional HTTP headers to send with webhook requests. Keys are header names, values are header values.

**Types:**
- `webhook_urls`: `set(string)`
- `http_headers`: `map(string)`

### ServiceNow Channel Attributes

* `service_now_url` - (Required) The ServiceNow instance URL.
* `username` - (Required) The username for ServiceNow authentication.
* `password` - (Optional) The password for ServiceNow authentication. Required when creating the resource. Uses state preservation for updates.
* `auto_close_incidents` - (Optional) Whether to automatically close incidents in ServiceNow when alerts are resolved.

**Types:**
- `service_now_url`: `string`
- `username`: `string`
- `password`: `string`
- `auto_close_incidents`: `bool`

### ServiceNow Enhanced Channel Attributes

* `service_now_url` - (Required) The ServiceNow instance URL.
* `username` - (Required) The username for ServiceNow authentication.
* `password` - (Optional) The password for ServiceNow authentication. Required when creating the resource. Uses state preservation for updates.
* `tenant` - (Required) The tenant identifier for ServiceNow Enhanced.
* `unit` - (Required) The unit identifier for ServiceNow Enhanced.
* `instana_url` - (Optional) The Instana URL for linking back from ServiceNow incidents.
* `auto_close_incidents` - (Optional) Whether to automatically close incidents in ServiceNow.
* `enable_send_instana_notes` - (Optional) Whether to send Instana notes to ServiceNow.
* `enable_send_service_now_activities` - (Optional) Whether to send ServiceNow activities.
* `enable_send_service_now_work_notes` - (Optional) Whether to send ServiceNow work notes.
* `manually_closed_incidents` - (Optional) Whether incidents are manually closed.
* `resolution_of_incident` - (Optional) Whether to resolve incidents.
* `snow_status_on_close_event` - (Optional) The ServiceNow status code when closing events.

**Types:**
- String attributes: `string`
- Boolean attributes: `bool`
- `snow_status_on_close_event`: `int64`

### Prometheus Webhook Channel Attributes

* `webhook_url` - (Required) The URL of the Prometheus webhook endpoint.
* `receiver` - (Optional) The receiver name for Prometheus Alertmanager.

**Type:** `string` for both attributes

### Webex Teams Webhook Channel Attributes

* `webhook_url` - (Required) The URL of the Webex Teams webhook.

**Type:** `string`

### Watson AIOps Webhook Channel Attributes

* `webhook_url` - (Required) The URL of the Watson AIOps webhook endpoint.
* `http_headers` - (Optional) List of HTTP headers to send with webhook requests. Each header should be in the format `"Header-Name: value"`.

**Types:**
- `webhook_url`: `string`
- `http_headers`: `list(string)`

### Slack App Channel Attributes

* `app_id` - (Required) The App ID of the Slack App.
* `team_id` - (Required) The Team ID where the Slack App is installed.
* `team_name` - (Required) The Team Name where the Slack App is installed.
* `channel_id` - (Required) The Channel ID where alerts will be sent.
* `channel_name` - (Required) The Channel Name where alerts will be sent.
* `emoji_rendering` - (Optional) Whether to enable emoji rendering in alert messages.

**Types:**
- String attributes: `string`
- `emoji_rendering`: `bool`

### MS Teams App Channel Attributes

* `api_token_id` - (Required) The API Token ID for MS Teams App authentication.
* `team_id` - (Required) The Team ID where the MS Teams App is installed.
* `team_name` - (Required) The Team Name where the MS Teams App is installed.
* `channel_id` - (Required) The Channel ID where alerts will be sent.
* `channel_name` - (Required) The Channel Name where alerts will be sent.
* `instana_url` - (Required) The Instana URL for linking back from MS Teams.
* `service_url` - (Required) The MS Teams service URL.
* `tenant_id` - (Required) The Tenant ID for MS Teams.
* `tenant_name` - (Required) The Tenant Name for MS Teams.

**Type:** `string` for all attributes

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the alerting channel

## Import

Alerting channels can be imported using the `id`, e.g.:

```bash
$ terraform import instana_alerting_channel.my_channel 60845e4e5e6b9cf8fc2868da
```

## Notes

### Channel Type Exclusivity

Only one channel type can be configured per alerting channel resource. If you need to send alerts to multiple destinations, create separate alerting channel resources for each destination.

### Webhook URL Security

When using webhook-based channels (webhook, google_chat, office_365, slack, prometheus_webhook, webex_teams_webhook, watson_aiops_webhook), ensure that:

- Webhook URLs are stored securely (use Terraform variables or secrets management)
- URLs use HTTPS for secure transmission
- Authentication tokens/keys are not hardcoded in configurations