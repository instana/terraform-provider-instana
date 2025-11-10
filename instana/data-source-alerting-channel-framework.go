package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaAlertingChannelFramework the name of the terraform-provider-instana data source to read alerting channel
const DataSourceInstanaAlertingChannelFramework = "alerting_channel"

// AlertingChannelDataSourceModel represents the data model for the alerting channel data source
type AlertingChannelDataSourceModel struct {
	ID                    types.String                `tfsdk:"id"`
	Name                  types.String                `tfsdk:"name"`
	Email                 *EmailModel                 `tfsdk:"email"`
	OpsGenie              *OpsGenieModel              `tfsdk:"ops_genie"`
	PagerDuty             *PagerDutyModel             `tfsdk:"pager_duty"`
	Slack                 *SlackModel                 `tfsdk:"slack"`
	Splunk                *SplunkModel                `tfsdk:"splunk"`
	VictorOps             *VictorOpsModel             `tfsdk:"victor_ops"`
	Webhook               *WebhookModel               `tfsdk:"webhook"`
	Office365             *WebhookBasedModel          `tfsdk:"office_365"`
	GoogleChat            *WebhookBasedModel          `tfsdk:"google_chat"`
	ServiceNow            *ServiceNowModel            `tfsdk:"service_now"`
	ServiceNowApplication *ServiceNowApplicationModel `tfsdk:"service_now_application"`
	PrometheusWebhook     *PrometheusWebhookModel     `tfsdk:"prometheus_webhook"`
	WebexTeamsWebhook     *WebhookBasedModel          `tfsdk:"webex_teams_webhook"`
	WatsonAIOpsWebhook    *WatsonAIOpsWebhookModel    `tfsdk:"watson_aiops_webhook"`
}

// NewAlertingChannelDataSourceFramework creates a new data source for alerting channel
func NewAlertingChannelDataSourceFramework() datasource.DataSource {
	return &alertingChannelDataSourceFramework{}
}

type alertingChannelDataSourceFramework struct {
	instanaAPI restapi.InstanaAPI
}

func (d *alertingChannelDataSourceFramework) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaAlertingChannelFramework
}

func (d *alertingChannelDataSourceFramework) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for an Instana alerting channel. Alerting channels are used to send notifications when alerts are triggered.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the alerting channel.",
				Computed:    true,
			},
			AlertingChannelFieldName: schema.StringAttribute{
				Description: "The name of the alerting channel.",
				Required:    true,
			},
			AlertingChannelFieldChannelEmail: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Email channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelEmailFieldEmails: schema.SetAttribute{
						Description: "The list of emails of the Email alerting channel",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			AlertingChannelFieldChannelOpsGenie: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Ops Genie channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelOpsGenieFieldAPIKey: schema.StringAttribute{
						Description: "The OpsGenie API Key of the OpsGenie alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelOpsGenieFieldTags: schema.ListAttribute{
						Description: "The OpsGenie tags of the OpsGenie alerting channel",
						Computed:    true,
						ElementType: types.StringType,
					},
					AlertingChannelOpsGenieFieldRegion: schema.StringAttribute{
						Description: "The OpsGenie region of the OpsGenie alerting channel",
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelPageDuty: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Pager Duty channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelPagerDutyFieldServiceIntegrationKey: schema.StringAttribute{
						Description: "The Service Integration Key of the PagerDuty alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelSlack: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Slack channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelSlackFieldWebhookURL: schema.StringAttribute{
						Description: "The webhook URL of the Slack alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelSlackFieldIconURL: schema.StringAttribute{
						Description: "The icon URL of the Slack alerting channel",
						Computed:    true,
					},
					AlertingChannelSlackFieldChannel: schema.StringAttribute{
						Description: "The Slack channel of the Slack alerting channel",
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelSplunk: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Splunk channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelSplunkFieldURL: schema.StringAttribute{
						Description: "The URL of the Splunk alerting channel",
						Computed:    true,
					},
					AlertingChannelSplunkFieldToken: schema.StringAttribute{
						Description: "The token of the Splunk alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelVictorOps: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the VictorOps channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelVictorOpsFieldAPIKey: schema.StringAttribute{
						Description: "The API Key of the VictorOps alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelVictorOpsFieldRoutingKey: schema.StringAttribute{
						Description: "The Routing Key of the VictorOps alerting channel",
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Webhook channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookFieldWebhookURLs: schema.SetAttribute{
						Description: "The list of webhook urls of the Webhook alerting channel",
						Computed:    true,
						ElementType: types.StringType,
					},
					AlertingChannelWebhookFieldHTTPHeaders: schema.MapAttribute{
						Description: "The optional map of HTTP headers of the Webhook alerting channel",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			AlertingChannelFieldChannelOffice365: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Office 365 channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: "The webhook URL of the Office 365 alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelGoogleChat: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Google Chat channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: "The webhook URL of the Google Chat alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelServiceNow: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the ServiceNow channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
						Description: "The ServiceNow URL of the ServiceNow alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
						Description: "The username of the ServiceNow alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
						Description: "The password of the ServiceNow alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelServiceNowFieldAutoCloseIncidents: schema.BoolAttribute{
						Description: "Whether to automatically close incidents when alerts are resolved",
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelServiceNowApplication: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the ServiceNow ITSM (Enhanced) channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
						Description: "The ServiceNow URL of the ServiceNow ITSM (Enhanced) alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
						Description: "The username of the ServiceNow ITSM (Enhanced) alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
						Description: "The password of the ServiceNow ITSM (Enhanced) alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelServiceNowApplicationFieldTenant: schema.StringAttribute{
						Description: "The tenant of the ServiceNow ITSM (Enhanced) alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldUnit: schema.StringAttribute{
						Description: "The unit of the ServiceNow ITSM (Enhanced) alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldInstanaURL: schema.StringAttribute{
						Description: "The Instana URL of the ServiceNow ITSM (Enhanced) alerting channel",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldEnableSendInstanaNotes: schema.BoolAttribute{
						Description: "Whether to enable sending Instana notes",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldEnableSendServiceNowActivities: schema.BoolAttribute{
						Description: "Whether to enable sending ServiceNow activities",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldEnableSendServiceNowWorkNotes: schema.BoolAttribute{
						Description: "Whether to enable sending ServiceNow work notes",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldManuallyClosedIncidents: schema.BoolAttribute{
						Description: "Whether incidents are manually closed",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldResolutionOfIncident: schema.StringAttribute{
						Description: "The resolution of incident",
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldSnowStatusOnCloseEvent: schema.StringAttribute{
						Description: "The ServiceNow status on close event",
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelPrometheusWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Prometheus Webhook channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: "The webhook URL of the Prometheus Webhook alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelPrometheusWebhookFieldReceiver: schema.StringAttribute{
						Description: "The receiver of the Prometheus Webhook alerting channel",
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelWebexTeamsWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the Webex Teams Webhook channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: "The webhook URL of the Webex Teams Webhook alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelWatsonAIOpsWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The configuration of the IBM Cloud Pack (Watson AIOps) Webhook channel",
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: "The webhook URL of the IBM Cloud Pack (Watson AIOps) Webhook alerting channel",
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelWebhookFieldHTTPHeaders: schema.ListAttribute{
						Description: "The list of HTTP headers for the IBM Cloud Pack (Watson AIOps) Webhook alerting channel",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (d *alertingChannelDataSourceFramework) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *alertingChannelDataSourceFramework) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AlertingChannelDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the name from the configuration
	name := data.Name.ValueString()

	// Get all alerting channels
	channels, err := d.instanaAPI.AlertingChannels().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading alerting channels",
			fmt.Sprintf("Could not read alerting channels: %s", err),
		)
		return
	}

	// Find the channel with the matching name
	var matchingChannel *restapi.AlertingChannel
	for _, channel := range *channels {
		if channel.Name == name {
			matchingChannel = channel
			break
		}
	}

	if matchingChannel == nil {
		resp.Diagnostics.AddError(
			"Alerting channel not found",
			fmt.Sprintf("No alerting channel found with name: %s", name),
		)
		return
	}

	// Update the data model with the channel details
	data.ID = types.StringValue(matchingChannel.ID)

	// Get resource handle to reuse mapping functions
	resourceHandle := NewAlertingChannelResourceHandleFramework().(*alertingChannelResourceFramework)

	// All channel types are nil by default (pointer models)
	// Only the matching channel type will be populated

	// Set the appropriate channel type based on the alerting channel kind
	switch matchingChannel.Kind {
	case restapi.EmailChannelType:
		emailChannel, emailDiags := resourceHandle.mapEmailChannelToState(ctx, matchingChannel)
		if emailDiags.HasError() {
			resp.Diagnostics.Append(emailDiags...)
			return
		}
		data.Email = emailChannel
	case restapi.OpsGenieChannelType:
		opsGenieChannel, opsGenieDiags := resourceHandle.mapOpsGenieChannelToState(ctx, matchingChannel)
		if opsGenieDiags.HasError() {
			resp.Diagnostics.Append(opsGenieDiags...)
			return
		}
		data.OpsGenie = opsGenieChannel
	case restapi.PagerDutyChannelType:
		pagerDutyChannel, pagerDutyDiags := resourceHandle.mapPagerDutyChannelToState(ctx, matchingChannel)
		if pagerDutyDiags.HasError() {
			resp.Diagnostics.Append(pagerDutyDiags...)
			return
		}
		data.PagerDuty = pagerDutyChannel
	case restapi.SlackChannelType:
		slackChannel, slackDiags := resourceHandle.mapSlackChannelToState(ctx, matchingChannel)
		if slackDiags.HasError() {
			resp.Diagnostics.Append(slackDiags...)
			return
		}
		data.Slack = slackChannel
	case restapi.SplunkChannelType:
		splunkChannel, splunkDiags := resourceHandle.mapSplunkChannelToState(ctx, matchingChannel)
		if splunkDiags.HasError() {
			resp.Diagnostics.Append(splunkDiags...)
			return
		}
		data.Splunk = splunkChannel
	case restapi.VictorOpsChannelType:
		victorOpsChannel, victorOpsDiags := resourceHandle.mapVictorOpsChannelToState(ctx, matchingChannel)
		if victorOpsDiags.HasError() {
			resp.Diagnostics.Append(victorOpsDiags...)
			return
		}
		data.VictorOps = victorOpsChannel
	case restapi.WebhookChannelType:
		webhookChannel, webhookDiags := resourceHandle.mapWebhookChannelToState(ctx, matchingChannel)
		if webhookDiags.HasError() {
			resp.Diagnostics.Append(webhookDiags...)
			return
		}
		data.Webhook = webhookChannel
	case restapi.Office365ChannelType:
		office365Channel, office365Diags := resourceHandle.mapWebhookBasedChannelToState(ctx, matchingChannel)
		if office365Diags.HasError() {
			resp.Diagnostics.Append(office365Diags...)
			return
		}
		data.Office365 = office365Channel
	case restapi.GoogleChatChannelType:
		googleChatChannel, googleChatDiags := resourceHandle.mapWebhookBasedChannelToState(ctx, matchingChannel)
		if googleChatDiags.HasError() {
			resp.Diagnostics.Append(googleChatDiags...)
			return
		}
		data.GoogleChat = googleChatChannel
	case restapi.ServiceNowChannelType:
		serviceNowChannel, serviceNowDiags := resourceHandle.mapServiceNowChannelToState(ctx, matchingChannel)
		if serviceNowDiags.HasError() {
			resp.Diagnostics.Append(serviceNowDiags...)
			return
		}
		data.ServiceNow = serviceNowChannel
	case restapi.ServiceNowApplicationChannelType:
		serviceNowApplicationChannel, serviceNowApplicationDiags := resourceHandle.mapServiceNowApplicationChannelToState(ctx, matchingChannel)
		if serviceNowApplicationDiags.HasError() {
			resp.Diagnostics.Append(serviceNowApplicationDiags...)
			return
		}
		data.ServiceNowApplication = serviceNowApplicationChannel
	case restapi.PrometheusWebhookChannelType:
		prometheusWebhookChannel, prometheusWebhookDiags := resourceHandle.mapPrometheusWebhookChannelToState(ctx, matchingChannel)
		if prometheusWebhookDiags.HasError() {
			resp.Diagnostics.Append(prometheusWebhookDiags...)
			return
		}
		data.PrometheusWebhook = prometheusWebhookChannel
	case restapi.WebexTeamsWebhookChannelType:
		webexTeamsWebhookChannel, webexTeamsWebhookDiags := resourceHandle.mapWebhookBasedChannelToState(ctx, matchingChannel)
		if webexTeamsWebhookDiags.HasError() {
			resp.Diagnostics.Append(webexTeamsWebhookDiags...)
			return
		}
		data.WebexTeamsWebhook = webexTeamsWebhookChannel
	case restapi.WatsonAIOpsWebhookChannelType:
		watsonAIOpsWebhookChannel, watsonAIOpsWebhookDiags := resourceHandle.mapWatsonAIOpsWebhookChannelToState(ctx, matchingChannel)
		if watsonAIOpsWebhookDiags.HasError() {
			resp.Diagnostics.Append(watsonAIOpsWebhookDiags...)
			return
		}
		data.WatsonAIOpsWebhook = watsonAIOpsWebhookChannel
	default:
		resp.Diagnostics.AddError(
			"Unsupported alerting channel type",
			fmt.Sprintf("Received unsupported alerting channel of type %s", matchingChannel.Kind),
		)
		return
	}

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
