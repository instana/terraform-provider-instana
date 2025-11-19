package alertingchannel

const (
	// ResourceInstanaAlertingChannelFramework the name of the terraform-provider-instana resource to manage alerting channels
	ResourceInstanaAlertingChannelFramework = "alerting_channel"

	//AlertingChannelFieldName constant value for the schema field name
	AlertingChannelFieldName = "name"
	//AlertingChannelFieldID constant value for the schema field id
	AlertingChannelFieldID = "id"

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

	//AlertingChannelFieldChannelSlackApp const for schema field of the Slack App channel
	AlertingChannelFieldChannelSlackApp = "slack_app"
	//AlertingChannelSlackAppFieldAppID const for the appId field of the Slack App alerting channel
	AlertingChannelSlackAppFieldAppID = "app_id"
	//AlertingChannelSlackAppFieldTeamID const for the teamId field of the Slack App alerting channel
	AlertingChannelSlackAppFieldTeamID = "team_id"
	//AlertingChannelSlackAppFieldTeamName const for the teamName field of the Slack App alerting channel
	AlertingChannelSlackAppFieldTeamName = "team_name"
	//AlertingChannelSlackAppFieldChannelID const for the channelId field of the Slack App alerting channel
	AlertingChannelSlackAppFieldChannelID = "channel_id"
	//AlertingChannelSlackAppFieldChannelName const for the channelName field of the Slack App alerting channel
	AlertingChannelSlackAppFieldChannelName = "channel_name"
	//AlertingChannelSlackAppFieldEmojiRendering const for the emojiRendering field of the Slack App alerting channel
	AlertingChannelSlackAppFieldEmojiRendering = "emoji_rendering"

	//AlertingChannelFieldChannelMsTeamsApp const for schema field of the MS Teams App channel
	AlertingChannelFieldChannelMsTeamsApp = "ms_teams_app"
	//AlertingChannelMsTeamsAppFieldAPITokenID const for the apiTokenId field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldAPITokenID = "api_token_id"
	//AlertingChannelMsTeamsAppFieldTeamID const for the teamId field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldTeamID = "team_id"
	//AlertingChannelMsTeamsAppFieldTeamName const for the teamName field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldTeamName = "team_name"
	//AlertingChannelMsTeamsAppFieldChannelID const for the channelId field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldChannelID = "channel_id"
	//AlertingChannelMsTeamsAppFieldChannelName const for the channelName field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldChannelName = "channel_name"
	//AlertingChannelMsTeamsAppFieldInstanaURL const for the instanaUrl field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldInstanaURL = "instana_url"
	//AlertingChannelMsTeamsAppFieldServiceURL const for the serviceUrl field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldServiceURL = "service_url"
	//AlertingChannelMsTeamsAppFieldTenantID const for the tenantId field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldTenantID = "tenant_id"
	//AlertingChannelMsTeamsAppFieldTenantName const for the tenantName field of the MS Teams App alerting channel
	AlertingChannelMsTeamsAppFieldTenantName = "tenant_name"

	// OpsGenie regions
	OpsGenieRegionEU = "EU"
	OpsGenieRegionUS = "US"

	// Separators and formatting
	HeaderSeparator = ": "
	TagSeparator    = ","

	// Schema descriptions
	AlertingChannelDescResource                    = "This resource manages alerting channels in Instana."
	AlertingChannelDescID                          = "The ID of the alerting channel."
	AlertingChannelDescName                        = "Configures the name of the alerting channel"
	AlertingChannelDescEmail                       = "The configuration of the Email channel"
	AlertingChannelDescEmailEmails                 = "The list of emails of the Email alerting channel"
	AlertingChannelDescOpsGenie                    = "The configuration of the Ops Genie channel"
	AlertingChannelDescOpsGenieAPIKey              = "The OpsGenie API Key of the OpsGenie alerting channel"
	AlertingChannelDescOpsGenieTags                = "The OpsGenie tags of the OpsGenie alerting channel"
	AlertingChannelDescOpsGenieRegion              = "The OpsGenie region (%s) of the OpsGenie alerting channel"
	AlertingChannelDescPagerDuty                   = "The configuration of the Pager Duty channel"
	AlertingChannelDescPagerDutyServiceKey         = "The Service Integration Key of the PagerDuty alerting channel"
	AlertingChannelDescSlack                       = "The configuration of the Slack channel"
	AlertingChannelDescSlackWebhookURL             = "The webhook URL of the Slack alerting channel"
	AlertingChannelDescSlackIconURL                = "The icon URL of the Slack alerting channel"
	AlertingChannelDescSlackChannel                = "The Slack channel of the Slack alerting channel"
	AlertingChannelDescSplunk                      = "The configuration of the Splunk channel"
	AlertingChannelDescSplunkURL                   = "The URL of the Splunk alerting channel"
	AlertingChannelDescSplunkToken                 = "The token of the Splunk alerting channel"
	AlertingChannelDescVictorOps                   = "The configuration of the VictorOps channel"
	AlertingChannelDescVictorOpsAPIKey             = "The API Key of the VictorOps alerting channel"
	AlertingChannelDescVictorOpsRoutingKey         = "The Routing Key of the VictorOps alerting channel"
	AlertingChannelDescWebhook                     = "The configuration of the Webhook channel"
	AlertingChannelDescWebhookWebhookURLs          = "The list of webhook urls of the Webhook alerting channel"
	AlertingChannelDescWebhookHTTPHeaders          = "The optional map of HTTP headers of the Webhook alerting channel"
	AlertingChannelDescOffice365                   = "The configuration of the Office 365 channel"
	AlertingChannelDescOffice365WebhookURL         = "The webhook URL of the Office 365 alerting channel"
	AlertingChannelDescGoogleChat                  = "The configuration of the Google Chat channel"
	AlertingChannelDescGoogleChatWebhookURL        = "The webhook URL of the Google Chat alerting channel"
	AlertingChannelDescServiceNow                  = "The configuration of the ServiceNow channel"
	AlertingChannelDescServiceNowURL               = "The ServiceNow URL of the ServiceNow alerting channel"
	AlertingChannelDescServiceNowUsername          = "The username of the ServiceNow alerting channel"
	AlertingChannelDescServiceNowPassword          = "The password of the ServiceNow alerting channel"
	AlertingChannelDescServiceNowAutoClose         = "Whether to automatically close incidents in ServiceNow"
	AlertingChannelDescServiceNowApplication       = "The configuration of the ServiceNow Enhanced (ITSM) channel"
	AlertingChannelDescServiceNowAppURL            = "The ServiceNow URL of the ServiceNow Enhanced alerting channel"
	AlertingChannelDescServiceNowAppUsername       = "The username of the ServiceNow Enhanced alerting channel"
	AlertingChannelDescServiceNowAppPassword       = "The password of the ServiceNow Enhanced alerting channel"
	AlertingChannelDescServiceNowAppTenant         = "The tenant of the ServiceNow Enhanced alerting channel"
	AlertingChannelDescServiceNowAppUnit           = "The unit of the ServiceNow Enhanced alerting channel"
	AlertingChannelDescServiceNowAppInstanaURL     = "The Instana URL for the ServiceNow Enhanced alerting channel"
	AlertingChannelDescServiceNowAppSendNotes      = "Whether to send Instana notes to ServiceNow"
	AlertingChannelDescServiceNowAppSendActivities = "Whether to send ServiceNow activities"
	AlertingChannelDescServiceNowAppSendWorkNotes  = "Whether to send ServiceNow work notes"
	AlertingChannelDescServiceNowAppManualClose    = "Whether incidents are manually closed"
	AlertingChannelDescServiceNowAppResolution     = "Whether to resolve incidents"
	AlertingChannelDescServiceNowAppCloseStatus    = "The ServiceNow status code when closing events"
	AlertingChannelDescPrometheusWebhook           = "The configuration of the Prometheus Webhook channel"
	AlertingChannelDescPrometheusWebhookURL        = "The webhook URL of the Prometheus Webhook alerting channel"
	AlertingChannelDescPrometheusWebhookReceiver   = "The receiver of the Prometheus Webhook alerting channel"
	AlertingChannelDescWebexTeamsWebhook           = "The configuration of the Webex Teams Webhook channel"
	AlertingChannelDescWebexTeamsWebhookURL        = "The webhook URL of the Webex Teams Webhook alerting channel"
	AlertingChannelDescWatsonAIOpsWebhook          = "The configuration of the Watson AIOps Webhook channel"
	AlertingChannelDescWatsonAIOpsWebhookURL       = "The webhook URL of the Watson AIOps Webhook alerting channel"
	AlertingChannelDescWatsonAIOpsHTTPHeaders      = "The list of HTTP headers for the Watson AIOps Webhook alerting channel"
	AlertingChannelDescSlackApp                    = "The configuration of the Slack App (bidirectional) channel"
	AlertingChannelDescSlackAppAppID               = "The App ID of the Slack App alerting channel"
	AlertingChannelDescSlackAppTeamID              = "The Team ID of the Slack App alerting channel"
	AlertingChannelDescSlackAppTeamName            = "The Team Name of the Slack App alerting channel"
	AlertingChannelDescSlackAppChannelID           = "The Channel ID of the Slack App alerting channel"
	AlertingChannelDescSlackAppChannelName         = "The Channel Name of the Slack App alerting channel"
	AlertingChannelDescSlackAppEmojiRendering      = "Whether to enable emoji rendering in the Slack App alerting channel"
	AlertingChannelDescMsTeamsApp                  = "The configuration of the MS Teams App (bidirectional) channel"
	AlertingChannelDescMsTeamsAppAPITokenID        = "The API Token ID of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppTeamID            = "The Team ID of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppTeamName          = "The Team Name of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppChannelID         = "The Channel ID of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppChannelName       = "The Channel Name of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppInstanaURL        = "The Instana URL for the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppServiceURL        = "The Service URL of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppTenantID          = "The Tenant ID of the MS Teams App alerting channel"
	AlertingChannelDescMsTeamsAppTenantName        = "The Tenant Name of the MS Teams App alerting channel"

	// Error messages
	AlertingChannelErrUnsupportedType       = "Unsupported alerting channel type"
	AlertingChannelErrUnsupportedTypeMsg    = "Received unsupported alerting channel of type %s"
	AlertingChannelErrMissingPassword       = "Missing Password"
	AlertingChannelErrMissingPasswordMsg    = "password must be specified when creating the resource"
	AlertingChannelErrInstanaURLRequired    = "InstanaURL is required"
	AlertingChannelErrInstanaURLRequiredMsg = "InstanaURL is required when creating the resource"
	AlertingChannelErrInvalidConfig         = "Invalid Alerting Channel Configuration"
	AlertingChannelErrInvalidConfigMsg      = "No valid alerting channel configuration found. Please configure exactly one channel type."

	// Log messages
	AlertingChannelLogPasswordValue    = "passwordValue: %s"
	AlertingChannelLogInstanaURL       = "Inatna url : user %v"
	AlertingChannelLogMapServiceNowApp = "[DEBUG] mapServiceNowApplicationChannelFromState: %v"
	AlertingChannelLogInstanaURLDebug  = "[DEBUG] intana url: %v"
	AlertingChannelLogModelFromPlan    = "Model from plan"
	AlertingChannelLogModelFromState   = "Model from state"
	AlertingChannelLogDebugCall        = "[DEBUG] Call %s %s\n"
)
