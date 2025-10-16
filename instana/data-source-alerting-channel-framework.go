package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceInstanaAlertingChannelFramework the name of the terraform-provider-instana data source to read alerting channel
const DataSourceInstanaAlertingChannelFramework = "alerting_channel"

// AlertingChannelDataSourceModel represents the data model for the alerting channel data source
type AlertingChannelDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Email      types.List   `tfsdk:"email"`
	OpsGenie   types.List   `tfsdk:"ops_genie"`
	PagerDuty  types.List   `tfsdk:"pager_duty"`
	Slack      types.List   `tfsdk:"slack"`
	Splunk     types.List   `tfsdk:"splunk"`
	VictorOps  types.List   `tfsdk:"victor_ops"`
	Webhook    types.List   `tfsdk:"webhook"`
	Office365  types.List   `tfsdk:"office_365"`
	GoogleChat types.List   `tfsdk:"google_chat"`
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
		},
		Blocks: map[string]schema.Block{
			AlertingChannelFieldChannelEmail: schema.ListNestedBlock{
				Description: "The configuration of the Email channel",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						AlertingChannelEmailFieldEmails: schema.SetAttribute{
							Description: "The list of emails of the Email alerting channel",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			AlertingChannelFieldChannelOpsGenie: schema.ListNestedBlock{
				Description: "The configuration of the Ops Genie channel",
				NestedObject: schema.NestedBlockObject{
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
			},
			AlertingChannelFieldChannelPageDuty: schema.ListNestedBlock{
				Description: "The configuration of the Pager Duty channel",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						AlertingChannelPagerDutyFieldServiceIntegrationKey: schema.StringAttribute{
							Description: "The Service Integration Key of the PagerDuty alerting channel",
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			AlertingChannelFieldChannelSlack: schema.ListNestedBlock{
				Description: "The configuration of the Slack channel",
				NestedObject: schema.NestedBlockObject{
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
			},
			AlertingChannelFieldChannelSplunk: schema.ListNestedBlock{
				Description: "The configuration of the Splunk channel",
				NestedObject: schema.NestedBlockObject{
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
			},
			AlertingChannelFieldChannelVictorOps: schema.ListNestedBlock{
				Description: "The configuration of the VictorOps channel",
				NestedObject: schema.NestedBlockObject{
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
			},
			AlertingChannelFieldChannelWebhook: schema.ListNestedBlock{
				Description: "The configuration of the Webhook channel",
				NestedObject: schema.NestedBlockObject{
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
			},
			AlertingChannelFieldChannelOffice365: schema.ListNestedBlock{
				Description: "The configuration of the Office 365 channel",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
							Description: "The webhook URL of the Office 365 alerting channel",
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			AlertingChannelFieldChannelGoogleChat: schema.ListNestedBlock{
				Description: "The configuration of the Google Chat channel",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
							Description: "The webhook URL of the Google Chat alerting channel",
							Computed:    true,
							Sensitive:   true,
						},
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

	// Reset all channel types to null
	resourceHandle := NewAlertingChannelResourceHandleFramework().(*alertingChannelResourceFramework)

	// Initialize all channel types as null lists with proper attribute types
	data.Email = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelEmailFieldEmails: types.SetType{ElemType: types.StringType},
		},
	})
	data.OpsGenie = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelOpsGenieFieldAPIKey: types.StringType,
			AlertingChannelOpsGenieFieldRegion: types.StringType,
			AlertingChannelOpsGenieFieldTags:   types.ListType{ElemType: types.StringType},
		},
	})
	data.PagerDuty = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelPagerDutyFieldServiceIntegrationKey: types.StringType,
		},
	})
	data.Slack = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelSlackFieldWebhookURL: types.StringType,
			AlertingChannelSlackFieldIconURL:    types.StringType,
			AlertingChannelSlackFieldChannel:    types.StringType,
		},
	})
	data.Splunk = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelSplunkFieldURL:   types.StringType,
			AlertingChannelSplunkFieldToken: types.StringType,
		},
	})
	data.VictorOps = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelVictorOpsFieldAPIKey:     types.StringType,
			AlertingChannelVictorOpsFieldRoutingKey: types.StringType,
		},
	})
	data.Webhook = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelWebhookFieldWebhookURLs: types.SetType{ElemType: types.StringType},
			AlertingChannelWebhookFieldHTTPHeaders: types.MapType{ElemType: types.StringType},
		},
	})
	data.Office365 = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelWebhookBasedFieldWebhookURL: types.StringType,
		},
	})
	data.GoogleChat = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelWebhookBasedFieldWebhookURL: types.StringType,
		},
	})

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
