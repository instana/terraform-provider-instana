package restapi

// AlertingChannelsResourcePath path to Alerting channels resource of Instana RESTful API
const AlertingChannelsResourcePath = EventSettingsBasePath + "/alertingChannels"

// AlertingChannel is the representation of an alerting channel in Instana
type AlertingChannel struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Kind                  AlertingChannelType `json:"kind"`
	Emails                []string            `json:"emails"`
	WebhookURL            *string             `json:"webhookUrl"`
	APIKey                *string             `json:"apiKey"`
	Tags                  *string             `json:"tags"`
	Region                *string             `json:"region"`
	RoutingKey            *string             `json:"routingKey"`
	ServiceIntegrationKey *string             `json:"serviceIntegrationKey"`
	IconURL               *string             `json:"iconUrl"`
	Channel               *string             `json:"channel"`
	URL                   *string             `json:"url"`
	Token                 *string             `json:"token"`
	WebhookURLs           []string            `json:"webhookUrls"`
	Headers               []string            `json:"headers"`
	// ServiceNow fields
	ServiceNowURL      *string `json:"serviceNowUrl"`
	Username           *string `json:"username"`
	Password           *string `json:"password"`
	AutoCloseIncidents *bool   `json:"autoCloseIncidents"`
	// ServiceNow Enhanced (ITSM) fields
	Tenant                         *string `json:"tenant"`
	Unit                           *string `json:"unit"`
	InstanaURL                     *string `json:"instanaUrl"`
	EnableSendInstanaNotes         *bool   `json:"enableSendInstanaNotes"`
	EnableSendServiceNowActivities *bool   `json:"enableSendServiceNowActivities"`
	EnableSendServiceNowWorkNotes  *bool   `json:"enableSendServiceNowWorkNotes"`
	ManuallyClosedIncidents        *bool   `json:"manuallyClosedIncidents"`
	ResolutionOfIncident           *bool   `json:"resolutionOfIncident"`
	SnowStatusOnCloseEvent         *int    `json:"snowStatusOnCloseEvent"`
	// Prometheus Webhook fields
	Receiver *string `json:"receiver"`
	// Bidirectional Slack App fields
	AppID          *string `json:"appId"`
	TeamID         *string `json:"teamId"`
	TeamName       *string `json:"teamName"`
	ChannelID      *string `json:"channelId"`
	ChannelName    *string `json:"channelName"`
	EmojiRendering *bool   `json:"emojiRendering"`
	// Bidirectional MS Teams App fields
	APITokenID *string `json:"apiTokenId"`
	ServiceURL *string `json:"serviceUrl"`
	TenantID   *string `json:"tenantId"`
	TenantName *string `json:"tenantName"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (r *AlertingChannel) GetIDForResourcePath() string {
	return r.ID
}
