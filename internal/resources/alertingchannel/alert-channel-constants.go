package alertingchannel

const (
	// ResourceInstanaAlertingChannelFramework the name of the terraform-provider-instana resource to manage alerting channels
	ResourceInstanaAlertingChannelFramework = "alerting_channel"

	//AlertingChannelFieldName constant value for the schema field name
	AlertingChannelFieldName = "name"

	//AlertingChannelFieldChannelEmail const for schema field of the email channel
	AlertingChannelFieldChannelEmail = "email"
	//AlertingChannelEmailFieldEmails const for the emails field of the alerting channel
	AlertingChannelEmailFieldEmails = "emails"

	//AlertingChannelFieldChannelOpsGenie const for schema field of the OpsGenie channel
	AlertingChannelFieldChannelOpsGenie = "ops_genie"
	//AlertingChannelOpsGenieFieldAPIKey const for the api key field of the alerting channel OpsGenie
	AlertingChannelOpsGenieFieldAPIKey = "api_key"
	//AlertingChannelOpsGenieFieldTags const for the tags field of the alerting channel OpsGenie
	AlertingChannelOpsGenieFieldTags = "tags"
	//AlertingChannelOpsGenieFieldRegion const for the region field of the alerting channel OpsGenie
	AlertingChannelOpsGenieFieldRegion = "region"

	//AlertingChannelFieldChannelPageDuty const for schema field of the PagerDuty channel
	AlertingChannelFieldChannelPageDuty = "pager_duty"
	//AlertingChannelPagerDutyFieldServiceIntegrationKey const for the emails field of the alerting channel
	AlertingChannelPagerDutyFieldServiceIntegrationKey = "service_integration_key"

	//AlertingChannelFieldChannelSlack const for schema field of the Slack channel
	AlertingChannelFieldChannelSlack = "slack"
	//AlertingChannelSlackFieldWebhookURL const for the webhookUrl field of the Slack alerting channel
	AlertingChannelSlackFieldWebhookURL = "webhook_url"
	//AlertingChannelSlackFieldIconURL const for the iconURL field of the Slack alerting channel
	AlertingChannelSlackFieldIconURL = "icon_url"
	//AlertingChannelSlackFieldChannel const for the channel field of the Slack alerting channel
	AlertingChannelSlackFieldChannel = "channel"

	//AlertingChannelFieldChannelSplunk const for schema field of the Splunk channel
	AlertingChannelFieldChannelSplunk = "splunk"
	//AlertingChannelSplunkFieldURL const for the url field of the Splunk alerting channel
	AlertingChannelSplunkFieldURL = "url"
	//AlertingChannelSplunkFieldToken const for the token field of the Splunk alerting channel
	AlertingChannelSplunkFieldToken = "token"

	//AlertingChannelFieldChannelVictorOps const for schema field of the Victor Ops channel
	AlertingChannelFieldChannelVictorOps = "victor_ops"
	//AlertingChannelVictorOpsFieldAPIKey const for the apiKey field of the VictorOps alerting channel
	AlertingChannelVictorOpsFieldAPIKey = "api_key"
	//AlertingChannelVictorOpsFieldRoutingKey const for the routingKey field of the VictorOps alerting channel
	AlertingChannelVictorOpsFieldRoutingKey = "routing_key"

	//AlertingChannelFieldChannelWebhook const for schema field of the Webhook channel
	AlertingChannelFieldChannelWebhook = "webhook"
	//AlertingChannelWebhookFieldWebhookURLs const for the webhooks field of the Webhook alerting channel
	AlertingChannelWebhookFieldWebhookURLs = "webhook_urls"
	//AlertingChannelWebhookFieldHTTPHeaders const for the http headers field of the Webhook alerting channel
	AlertingChannelWebhookFieldHTTPHeaders = "http_headers"

	//AlertingChannelFieldChannelOffice365 const for schema field of the Office 365 channel
	AlertingChannelFieldChannelOffice365 = "office_365"
	//AlertingChannelFieldChannelGoogleChat const for schema field of the Google Chat channel
	AlertingChannelFieldChannelGoogleChat = "google_chat"
	//AlertingChannelWebhookBasedFieldWebhookURL const for the webhookUrl field of the alerting channel
	AlertingChannelWebhookBasedFieldWebhookURL = "webhook_url"

	//AlertingChannelFieldChannelServiceNow const for schema field of the ServiceNow channel
	AlertingChannelFieldChannelServiceNow = "service_now"
	//AlertingChannelServiceNowFieldServiceNowURL const for the serviceNowUrl field of the ServiceNow alerting channel
	AlertingChannelServiceNowFieldServiceNowURL = "service_now_url"
	//AlertingChannelServiceNowFieldUsername const for the username field of the ServiceNow alerting channel
	AlertingChannelServiceNowFieldUsername = "username"
	//AlertingChannelServiceNowFieldPassword const for the password field of the ServiceNow alerting channel
	AlertingChannelServiceNowFieldPassword = "password"
	//AlertingChannelServiceNowFieldAutoCloseIncidents const for the autoCloseIncidents field of the ServiceNow alerting channel
	AlertingChannelServiceNowFieldAutoCloseIncidents = "auto_close_incidents"

	//AlertingChannelFieldChannelServiceNowApplication const for schema field of the ServiceNow Enhanced (ITSM) channel
	AlertingChannelFieldChannelServiceNowApplication = "service_now_application"
	//AlertingChannelServiceNowApplicationFieldTenant const for the tenant field of the ServiceNow Enhanced alerting channel
	AlertingChannelServiceNowApplicationFieldTenant = "tenant"
	//AlertingChannelServiceNowApplicationFieldUnit const for the unit field of the ServiceNow Enhanced alerting channel
	AlertingChannelServiceNowApplicationFieldUnit = "unit"
	//AlertingChannelServiceNowApplicationFieldInstanaURL const for the instanaUrl field of the ServiceNow Enhanced alerting channel
	AlertingChannelServiceNowApplicationFieldInstanaURL = "instana_url"
	//AlertingChannelServiceNowApplicationFieldEnableSendInstanaNotes const for the enableSendInstanaNotes field
	AlertingChannelServiceNowApplicationFieldEnableSendInstanaNotes = "enable_send_instana_notes"
	//AlertingChannelServiceNowApplicationFieldEnableSendServiceNowActivities const for the enableSendServiceNowActivities field
	AlertingChannelServiceNowApplicationFieldEnableSendServiceNowActivities = "enable_send_service_now_activities"
	//AlertingChannelServiceNowApplicationFieldEnableSendServiceNowWorkNotes const for the enableSendServiceNowWorkNotes field
	AlertingChannelServiceNowApplicationFieldEnableSendServiceNowWorkNotes = "enable_send_service_now_work_notes"
	//AlertingChannelServiceNowApplicationFieldManuallyClosedIncidents const for the manuallyClosedIncidents field
	AlertingChannelServiceNowApplicationFieldManuallyClosedIncidents = "manually_closed_incidents"
	//AlertingChannelServiceNowApplicationFieldResolutionOfIncident const for the resolutionOfIncident field
	AlertingChannelServiceNowApplicationFieldResolutionOfIncident = "resolution_of_incident"
	//AlertingChannelServiceNowApplicationFieldSnowStatusOnCloseEvent const for the snowStatusOnCloseEvent field
	AlertingChannelServiceNowApplicationFieldSnowStatusOnCloseEvent = "snow_status_on_close_event"

	//AlertingChannelFieldChannelPrometheusWebhook const for schema field of the Prometheus Webhook channel
	AlertingChannelFieldChannelPrometheusWebhook = "prometheus_webhook"
	//AlertingChannelPrometheusWebhookFieldReceiver const for the receiver field of the Prometheus Webhook alerting channel
	AlertingChannelPrometheusWebhookFieldReceiver = "receiver"

	//AlertingChannelFieldChannelWebexTeamsWebhook const for schema field of the Webex Teams Webhook channel
	AlertingChannelFieldChannelWebexTeamsWebhook = "webex_teams_webhook"

	//AlertingChannelFieldChannelWatsonAIOpsWebhook const for schema field of the Watson AIOps Webhook channel
	AlertingChannelFieldChannelWatsonAIOpsWebhook = "watson_aiops_webhook"
)
