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
    webhook_url = "https://chat.googleapis.com/v1/spaces/AAAA/messages?key=YOUR_KEY&token=YOUR_TOKEN"
  }
}
```

### Office 365 Alerting Channel

#### Basic Office 365 Configuration
```hcl
resource "instana_alerting_channel" "office365_basic" {
  name = "office365-alerts"
  
  office_365 = {
    webhook_url = "https://outlook.office.com/webhook/YOUR_WEBHOOK_ID"
  }
}
```

### OpsGenie Alerting Channel

#### Basic OpsGenie Configuration
```hcl
resource "instana_alerting_channel" "opsgenie_eu" {
  name = "opsgenie-eu-alerts"
  
  ops_genie = {
    api_key = var.opsgenie_api_key
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
    api_key = var.opsgenie_api_key
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
    service_integration_key = var.pagerduty_integration_key
  }
}
```


### Slack Alerting Channel

#### Basic Slack Configuration
```hcl
resource "instana_alerting_channel" "slack_channel" {
  name                    = "slack alert"
  slack = {
    webhook_url = "https://hooks.slack.com/services/XXXX"
  }
}
```

#### Slack with Custom Icon
```hcl
resource "instana_alerting_channel" "slack_custom_icon" {
  name = "slack-custom-alerts"
  
  slack = {
    webhook_url = "https://hooks.slack.com/services/XXX"
    icon_url    = "https://example.com/instana-icon.png"
  }
}
```

#### Slack with Specific Channel
```hcl
resource "instana_alerting_channel" "slack_specific_channel" {
  name = "slack-production-alerts"
  
  slack = {
    webhook_url = "https://hooks.slack.com/services/XXX"
    icon_url    = "https://example.com/alert-icon.png"
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
    url   = "https://splunk.example.com:8088/services/collector"
    token = var.splunk_hec_token
  }
}
```

### VictorOps Alerting Channel

#### Basic VictorOps Configuration
```hcl
resource "instana_alerting_channel" "victorops_basic" {
  name = "victorops-alerts"
  
  victor_ops = {
    api_key     = var.victorops_api_key
    routing_key = "instana-alerts"
  }
}
```

#### VictorOps with Team-Specific Routing
```hcl
resource "instana_alerting_channel" "victorops_platform_team" {
  name = "victorops-platform-team"
  
  victor_ops = {
    api_key     = var.victorops_api_key
    routing_key = "platform-team"
  }
}

resource "instana_alerting_channel" "victorops_database_team" {
  name = "victorops-database-team"
  
  victor_ops = {
    api_key     = var.victorops_api_key
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
    webhook_urls = ["https://api.example.com/instana/alerts"]
  }
}
```

#### Multiple Webhook URLs
```hcl
resource "instana_alerting_channel" "webhook_multiple" {
  name = "webhook-multiple-endpoints"
  
  webhook = {
    http_headers = {}
    webhook_urls = [
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
    webhook_urls = ["https://api.example.com/instana/alerts"]
    http_headers = {
      "Authorization" = "Bearer ${var.webhook_api_token}"
      "X-API-Key"     = var.webhook_api_key
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
    service_now_url      = "https://your-instance.service-now.com"
    username             = var.servicenow_username
    password             = var.servicenow_password
    auto_close_incidents = false
  }
}
```

#### ServiceNow with Auto-Close
```hcl
resource "instana_alerting_channel" "servicenow_autoclose" {
  name = "servicenow-autoclose-alerts"
  
  service_now = {
    service_now_url      = "https://your-instance.service-now.com"
    username             = var.servicenow_username
    password             = var.servicenow_password
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
    instana_url                        = null
    manually_closed_incidents          = true
    username                           = var.servicenow_username
    password                           = var.servicenow_password
    resolution_of_incident             = true
    service_now_url                    = "https://service-now.com"
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
    webhook_url = "https://prometheus.example.com/api/v1/alerts"
  }
}
```

### Webex Teams Webhook Alerting Channel

#### Basic Webex Teams Configuration
```hcl
resource "instana_alerting_channel" "webex_basic" {
  name = "webex-teams-alerts"
  
  webex_teams_webhook = {
    webhook_url = "https://webexapis.com/v1/webhooks/incoming/YOUR_WEBHOOK_ID"
  }
}
```

### Watson AIOps Webhook Alerting Channel

#### Basic Watson AIOps Configuration
```hcl
resource "instana_alerting_channel" "watson_basic" {
  name = "watson-aiops-alerts"
  
  watson_aiops_webhook = {
    webhook_url = "https://watson-aiops.example.com/webhook"
  }
}
```

#### Watson AIOps with Custom Headers
```hcl
resource "instana_alerting_channel" "watson_with_headers" {
  name = "watson-aiops-authenticated"
  
  watson_aiops_webhook = {
    webhook_url = "https://watson-aiops.example.com/webhook"
    http_headers = [
      "Authorization: Bearer ${var.watson_api_token}",
      "X-API-Key: ${var.watson_api_key}",
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
    app_id       = "A01234567"
    team_id      = "T01234567"
    team_name    = "My Team"
    channel_id   = "C01234567"
    channel_name = "#alerts"
  }
}
```

#### Slack App with Emoji Rendering
```hcl
resource "instana_alerting_channel" "slack_app_with_emoji" {
  name = "slack-app-emoji-alerts"
  
  slack_app = {
    app_id          = "A01234567"
    team_id         = "T01234567"
    team_name       = "My Team"
    channel_id      = "C01234567"
    channel_name    = "#production-alerts"
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
    api_token_id = var.msteams_api_token
    channel_id   = var.msteams_channel_id
    channel_name = var.msteams_channel_name
    instana_url  = var.instana_base_url
    service_url  = var.msteams_service_url
    team_id      = var.msteams_team_id
    team_name    = var.msteams_name
    tenant_id    = var.msteams_tenant_id
    tenant_name  = var.organization_name
  }
}
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

### Password Handling for ServiceNow Channels

For ServiceNow and ServiceNow Enhanced channels, the `password` attribute uses a plan modifier that preserves the state value when the password is not included in the plan. This means:

- The password **must** be provided when creating the resource
- The password can be omitted in subsequent updates, and the existing value will be preserved
- If you need to update the password, explicitly include it in your configuration

### Channel Type Exclusivity

Only one channel type can be configured per alerting channel resource. If you need to send alerts to multiple destinations, create separate alerting channel resources for each destination.

### Webhook URL Security

When using webhook-based channels (webhook, google_chat, office_365, slack, prometheus_webhook, webex_teams_webhook, watson_aiops_webhook), ensure that:

- Webhook URLs are stored securely (use Terraform variables or secrets management)
- URLs use HTTPS for secure transmission
- Authentication tokens/keys are not hardcoded in configurations