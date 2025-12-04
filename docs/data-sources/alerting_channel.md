# Alerting Channel Data Source

Data source to retrieve details about existing alerting channels

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
data "instana_alerting_channel" "example" {
  name = "my-alerting-channel"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel

## Attribute Reference

Exactly one of the following items is provided depending on the type of the alerting channel:

* `email` - configuration of a email alerting channel - [Details](#email)
* `google_chat` - configuration of a Google Chat alerting channel - [Details](#google-chat)
* `office_365` - configuration of a Office 365 alerting channel - [Details](#office-365)
* `ops_genie` - configuration of a OpsGenie alerting channel - [Details](#opsgenie)
* `pager_duty` - configuration of a PagerDuty alerting channel - [Details](#pagerduty)
* `prometheus_webhook` - configuration of a Prometheus Webhook alerting channel - [Details](#prometheus-webhook)
* `service_now` - configuration of a ServiceNow alerting channel - [Details](#servicenow)
* `service_now_application` - configuration of a ServiceNow Application alerting channel - [Details](#servicenow-application)
* `slack` - configuration of a Slack alerting channel - [Details](#slack)
* `splunk` - configuration of a Splunk alerting channel - [Details](#splunk)
* `victor_ops` - configuration of a VictorOps alerting channel - [Details](#victorops)
* `watson_aiops_webhook` - configuration of a Watson AIOps Webhook alerting channel - [Details](#watson-aiops-webhook)
* `webex_teams_webhook` - configuration of a Webex Teams Webhook alerting channel - [Details](#webex-teams-webhook)
* `webhook` - configuration of a webhook alerting channel - [Details](#webhook)

### Email

* `emails` - the list of target email addresses

### Google Chat

* `webhook_url` - the URL of the Google Chat Webhook where the alert will be sent to

### Office 365

* `webhook_url` - the URL of the Google Chat Webhook where the alert will be sent to

### OpsGenie

* `api_key` - the API Key for authentication at the Ops Genie API
* `tags` - a list of tags (strings) for the alert in Ops Genie
* `region` - the target Ops Genie region

### PagerDuty

* `service_integration_key` - the key for the service integration in pager duty

### Slack

* `webhook_url` - the URL of the Slack webhook to send alerts to
* `icon_url` - the URL to the icon which should be rendered in the slack message
* `channel` - the target Slack channel where the alert should be posted

### Splunk

* `url` - the target Splunk endpoint URL
* `token` - the authentication token to login at the Splunk API

### VictorOps

* `api_key` - the api key to authenticate at the VictorOps API
* `routing_key` - the routing key used by VictoryOps to route the alert to the desired targe

### Prometheus Webhook

* `webhook_url` - the URL of the Prometheus Webhook where the alert will be sent to
* `receiver` - the receiver name for the Prometheus Webhook

### ServiceNow

* `service_now_url` - the ServiceNow instance URL
* `username` - the username for authentication
* `password` - the password for authentication
* `auto_close_incidents` - boolean flag to automatically close incidents when the alert is resolved

### ServiceNow Application

* `service_now_url` - the ServiceNow instance URL
* `username` - the username for authentication
* `password` - the password for authentication
* `tenant` - the tenant identifier
* `unit` - the unit identifier
* `instana_url` - the Instana URL for linking back to Instana
* `enable_send_instana_notes` - boolean flag to enable sending Instana notes to ServiceNow
* `enable_send_service_now_activities` - boolean flag to enable sending ServiceNow activities
* `enable_send_service_now_work_notes` - boolean flag to enable sending ServiceNow work notes
* `manually_closed_incidents` - boolean flag for manually closed incidents
* `resolution_of_incident` - the resolution text for incidents
* `snow_status_on_close_event` - the ServiceNow status when closing events

### Victor Ops

* `api_key` - the api key to authenticate at the VictorOps API
* `routing_key` - the routing key used by VictoryOps to route the alert to the desired target

### Watson AIOps Webhook

* `webhook_url` - the URL of the Watson AIOps Webhook where the alert will be sent to
* `http_headers` - list of additional http headers which will be sent to the webhook

### Webex Teams Webhook

* `webhook_url` - the URL of the Webex Teams Webhook where the alert will be sent to

### Webhook

* `webhook_urls` - the list of webhook URLs where the alert will be sent to
* `http_headers` - key/value map of additional http headers which will be sent to the webhook
