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
		Description: AlertingChannelDescDataSource,
		Attributes: map[string]schema.Attribute{
			AlertingChannelDataSourceFieldID: schema.StringAttribute{
				Description: AlertingChannelDescID,
				Computed:    true,
			},
			AlertingChannelFieldName: schema.StringAttribute{
				Description: AlertingChannelDescName,
				Required:    true,
			},
			AlertingChannelFieldChannelEmail: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescEmail,
				Attributes: map[string]schema.Attribute{
					AlertingChannelEmailFieldEmails: schema.SetAttribute{
						Description: AlertingChannelDescEmailEmails,
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			AlertingChannelFieldChannelOpsGenie: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescOpsGenie,
				Attributes: map[string]schema.Attribute{
					AlertingChannelOpsGenieFieldAPIKey: schema.StringAttribute{
						Description: AlertingChannelDescOpsGenieAPIKey,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelOpsGenieFieldTags: schema.ListAttribute{
						Description: AlertingChannelDescOpsGenieTags,
						Computed:    true,
						ElementType: types.StringType,
					},
					AlertingChannelOpsGenieFieldRegion: schema.StringAttribute{
						Description: AlertingChannelDescOpsGenieRegion,
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelPageDuty: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescPagerDuty,
				Attributes: map[string]schema.Attribute{
					AlertingChannelPagerDutyFieldServiceIntegrationKey: schema.StringAttribute{
						Description: AlertingChannelDescPagerDutyServiceIntegrationKey,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelSlack: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescSlack,
				Attributes: map[string]schema.Attribute{
					AlertingChannelSlackFieldWebhookURL: schema.StringAttribute{
						Description: AlertingChannelDescSlackWebhookURL,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelSlackFieldIconURL: schema.StringAttribute{
						Description: AlertingChannelDescSlackIconURL,
						Computed:    true,
					},
					AlertingChannelSlackFieldChannel: schema.StringAttribute{
						Description: AlertingChannelDescSlackChannel,
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelSplunk: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescSplunk,
				Attributes: map[string]schema.Attribute{
					AlertingChannelSplunkFieldURL: schema.StringAttribute{
						Description: AlertingChannelDescSplunkURL,
						Computed:    true,
					},
					AlertingChannelSplunkFieldToken: schema.StringAttribute{
						Description: AlertingChannelDescSplunkToken,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelVictorOps: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescVictorOps,
				Attributes: map[string]schema.Attribute{
					AlertingChannelVictorOpsFieldAPIKey: schema.StringAttribute{
						Description: AlertingChannelDescVictorOpsAPIKey,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelVictorOpsFieldRoutingKey: schema.StringAttribute{
						Description: AlertingChannelDescVictorOpsRoutingKey,
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescWebhook,
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookFieldWebhookURLs: schema.SetAttribute{
						Description: AlertingChannelDescWebhookWebhookURLs,
						Computed:    true,
						ElementType: types.StringType,
					},
					AlertingChannelWebhookFieldHTTPHeaders: schema.MapAttribute{
						Description: AlertingChannelDescWebhookHTTPHeaders,
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			AlertingChannelFieldChannelOffice365: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescOffice365,
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: AlertingChannelDescOffice365WebhookURL,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelGoogleChat: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescGoogleChat,
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: AlertingChannelDescGoogleChatWebhookURL,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelServiceNow: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescServiceNow,
				Attributes: map[string]schema.Attribute{
					AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowURL,
						Computed:    true,
					},
					AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowUsername,
						Computed:    true,
					},
					AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowPassword,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelServiceNowFieldAutoCloseIncidents: schema.BoolAttribute{
						Description: AlertingChannelDescServiceNowAutoCloseIncidents,
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelServiceNowApplication: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescServiceNowApplication,
				Attributes: map[string]schema.Attribute{
					AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationURL,
						Computed:    true,
					},
					AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationUsername,
						Computed:    true,
					},
					AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationPassword,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelServiceNowApplicationFieldTenant: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationTenant,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldUnit: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationUnit,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldInstanaURL: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationInstanaURL,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldEnableSendInstanaNotes: schema.BoolAttribute{
						Description: AlertingChannelDescServiceNowApplicationEnableSendInstanaNotes,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldEnableSendServiceNowActivities: schema.BoolAttribute{
						Description: AlertingChannelDescServiceNowApplicationEnableSendServiceNowActivities,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldEnableSendServiceNowWorkNotes: schema.BoolAttribute{
						Description: AlertingChannelDescServiceNowApplicationEnableSendServiceNowWorkNotes,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldManuallyClosedIncidents: schema.BoolAttribute{
						Description: AlertingChannelDescServiceNowApplicationManuallyClosedIncidents,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldResolutionOfIncident: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationResolutionOfIncident,
						Computed:    true,
					},
					AlertingChannelServiceNowApplicationFieldSnowStatusOnCloseEvent: schema.StringAttribute{
						Description: AlertingChannelDescServiceNowApplicationSnowStatusOnCloseEvent,
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelPrometheusWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescPrometheusWebhook,
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: AlertingChannelDescPrometheusWebhookWebhookURL,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelPrometheusWebhookFieldReceiver: schema.StringAttribute{
						Description: AlertingChannelDescPrometheusWebhookReceiver,
						Computed:    true,
					},
				},
			},
			AlertingChannelFieldChannelWebexTeamsWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescWebexTeamsWebhook,
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: AlertingChannelDescWebexTeamsWebhookWebhookURL,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			AlertingChannelFieldChannelWatsonAIOpsWebhook: schema.SingleNestedAttribute{
				Computed:    true,
				Description: AlertingChannelDescWatsonAIOpsWebhook,
				Attributes: map[string]schema.Attribute{
					AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
						Description: AlertingChannelDescWatsonAIOpsWebhookWebhookURL,
						Computed:    true,
						Sensitive:   true,
					},
					AlertingChannelWebhookFieldHTTPHeaders: schema.ListAttribute{
						Description: AlertingChannelDescWatsonAIOpsWebhookHTTPHeaders,
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
			AlertingChannelErrUnexpectedConfigureType,
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
			AlertingChannelErrReadingChannels,
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
			AlertingChannelErrChannelNotFound,
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
			AlertingChannelErrUnsupportedChannelType,
			fmt.Sprintf("Received unsupported alerting channel of type %s", matchingChannel.Kind),
		)
		return
	}

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
