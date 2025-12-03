package alertingchannel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// NewAlertingChannelResourceHandle creates the resource handle for Alerting Channels
func NewAlertingChannelResourceHandle() resourcehandle.ResourceHandle[*restapi.AlertingChannel] {
	supportedOpsGenieRegions := []string{OpsGenieRegionEU, OpsGenieRegionUS}
	return &alertingChannelResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannel,
			Schema: schema.Schema{
				Description: AlertingChannelDescResource,
				Attributes: map[string]schema.Attribute{
					AlertingChannelFieldID: schema.StringAttribute{
						Computed:    true,
						Description: AlertingChannelDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AlertingChannelFieldName: schema.StringAttribute{
						Required:    true,
						Description: AlertingChannelDescName,
					},
					AlertingChannelFieldChannelEmail: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescEmail,
						Attributes: map[string]schema.Attribute{
							AlertingChannelEmailFieldEmails: schema.SetAttribute{
								Required:    true,
								Description: AlertingChannelDescEmailEmails,
								ElementType: types.StringType,
							},
						},
					},
					AlertingChannelFieldChannelOpsGenie: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescOpsGenie,
						Attributes: map[string]schema.Attribute{
							AlertingChannelOpsGenieFieldAPIKey: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescOpsGenieAPIKey,
							},
							AlertingChannelOpsGenieFieldTags: schema.ListAttribute{
								Required:    true,
								Description: AlertingChannelDescOpsGenieTags,
								ElementType: types.StringType,
							},
							AlertingChannelOpsGenieFieldRegion: schema.StringAttribute{
								Required:    true,
								Description: fmt.Sprintf(AlertingChannelDescOpsGenieRegion, strings.Join(supportedOpsGenieRegions, ", ")),
								Validators: []validator.String{
									stringvalidator.OneOf(supportedOpsGenieRegions...),
								},
							},
						},
					},
					AlertingChannelFieldChannelPageDuty: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescPagerDuty,
						Attributes: map[string]schema.Attribute{
							AlertingChannelPagerDutyFieldServiceIntegrationKey: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescPagerDutyServiceKey,
							},
						},
					},
					AlertingChannelFieldChannelSlack: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescSlack,
						Attributes: map[string]schema.Attribute{
							AlertingChannelSlackFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSlackWebhookURL,
							},
							AlertingChannelSlackFieldIconURL: schema.StringAttribute{
								Optional:    true,
								Description: AlertingChannelDescSlackIconURL,
							},
							AlertingChannelSlackFieldChannel: schema.StringAttribute{
								Optional:    true,
								Description: AlertingChannelDescSlackChannel,
							},
						},
					},
					AlertingChannelFieldChannelSplunk: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescSplunk,
						Attributes: map[string]schema.Attribute{
							AlertingChannelSplunkFieldURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSplunkURL,
							},
							AlertingChannelSplunkFieldToken: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSplunkToken,
							},
						},
					},
					AlertingChannelFieldChannelVictorOps: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescVictorOps,
						Attributes: map[string]schema.Attribute{
							AlertingChannelVictorOpsFieldAPIKey: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescVictorOpsAPIKey,
							},
							AlertingChannelVictorOpsFieldRoutingKey: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescVictorOpsRoutingKey,
							},
						},
					},
					AlertingChannelFieldChannelWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescWebhook,
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookFieldWebhookURLs: schema.SetAttribute{
								Required:    true,
								Description: AlertingChannelDescWebhookWebhookURLs,
								ElementType: types.StringType,
							},
							AlertingChannelWebhookFieldHTTPHeaders: schema.MapAttribute{
								Optional:    true,
								Description: AlertingChannelDescWebhookHTTPHeaders,
								ElementType: types.StringType,
							},
						},
					},
					AlertingChannelFieldChannelOffice365: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescOffice365,
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescOffice365WebhookURL,
							},
						},
					},
					AlertingChannelFieldChannelGoogleChat: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescGoogleChat,
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescGoogleChatWebhookURL,
							},
						},
					},
					AlertingChannelFieldChannelServiceNow: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescServiceNow,
						Attributes: map[string]schema.Attribute{
							AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescServiceNowURL,
							},
							AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescServiceNowUsername,
							},
							AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowPassword,
								PlanModifiers: []planmodifier.String{
									// When the plan does not include the password, keep the value from state.
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							AlertingChannelServiceNowFieldAutoCloseIncidents: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAutoClose,
							},
						},
					},
					AlertingChannelFieldChannelServiceNowApplication: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescServiceNowApplication,
						Attributes: map[string]schema.Attribute{
							AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescServiceNowAppURL,
							},
							AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescServiceNowAppUsername,
							},
							AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppPassword,
								PlanModifiers: []planmodifier.String{
									// When the plan does not include the password, keep the value from state.
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							AlertingChannelServiceNowApplicationFieldTenant: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescServiceNowAppTenant,
							},
							AlertingChannelServiceNowApplicationFieldUnit: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescServiceNowAppUnit,
							},
							AlertingChannelServiceNowFieldAutoCloseIncidents: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAutoClose,
							},
							AlertingChannelServiceNowApplicationFieldInstanaURL: schema.StringAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppInstanaURL,
							},
							AlertingChannelServiceNowApplicationFieldEnableSendInstanaNotes: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppSendNotes,
							},
							AlertingChannelServiceNowApplicationFieldEnableSendServiceNowActivities: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppSendActivities,
							},
							AlertingChannelServiceNowApplicationFieldEnableSendServiceNowWorkNotes: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppSendWorkNotes,
							},
							AlertingChannelServiceNowApplicationFieldManuallyClosedIncidents: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppManualClose,
							},
							AlertingChannelServiceNowApplicationFieldResolutionOfIncident: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppResolution,
							},
							AlertingChannelServiceNowApplicationFieldSnowStatusOnCloseEvent: schema.Int64Attribute{
								Optional:    true,
								Description: AlertingChannelDescServiceNowAppCloseStatus,
							},
						},
					},
					AlertingChannelFieldChannelPrometheusWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescPrometheusWebhook,
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescPrometheusWebhookURL,
							},
							AlertingChannelPrometheusWebhookFieldReceiver: schema.StringAttribute{
								Optional:    true,
								Description: AlertingChannelDescPrometheusWebhookReceiver,
							},
						},
					},
					AlertingChannelFieldChannelWebexTeamsWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescWebexTeamsWebhook,
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescWebexTeamsWebhookURL,
							},
						},
					},
					AlertingChannelFieldChannelWatsonAIOpsWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescWatsonAIOpsWebhook,
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescWatsonAIOpsWebhookURL,
							},
							AlertingChannelWebhookFieldHTTPHeaders: schema.ListAttribute{
								Optional:    true,
								Description: AlertingChannelDescWatsonAIOpsHTTPHeaders,
								ElementType: types.StringType,
							},
						},
					},
					AlertingChannelFieldChannelSlackApp: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescSlackApp,
						Attributes: map[string]schema.Attribute{
							AlertingChannelSlackAppFieldAppID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSlackAppAppID,
							},
							AlertingChannelSlackAppFieldTeamID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSlackAppTeamID,
							},
							AlertingChannelSlackAppFieldTeamName: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSlackAppTeamName,
							},
							AlertingChannelSlackAppFieldChannelID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSlackAppChannelID,
							},
							AlertingChannelSlackAppFieldChannelName: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescSlackAppChannelName,
							},
							AlertingChannelSlackAppFieldEmojiRendering: schema.BoolAttribute{
								Optional:    true,
								Description: AlertingChannelDescSlackAppEmojiRendering,
							},
						},
					},
					AlertingChannelFieldChannelMsTeamsApp: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AlertingChannelDescMsTeamsApp,
						Attributes: map[string]schema.Attribute{
							AlertingChannelMsTeamsAppFieldAPITokenID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppAPITokenID,
							},
							AlertingChannelMsTeamsAppFieldTeamID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppTeamID,
							},
							AlertingChannelMsTeamsAppFieldTeamName: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppTeamName,
							},
							AlertingChannelMsTeamsAppFieldChannelID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppChannelID,
							},
							AlertingChannelMsTeamsAppFieldChannelName: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppChannelName,
							},
							AlertingChannelMsTeamsAppFieldInstanaURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppInstanaURL,
							},
							AlertingChannelMsTeamsAppFieldServiceURL: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppServiceURL,
							},
							AlertingChannelMsTeamsAppFieldTenantID: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppTenantID,
							},
							AlertingChannelMsTeamsAppFieldTenantName: schema.StringAttribute{
								Required:    true,
								Description: AlertingChannelDescMsTeamsAppTenantName,
							},
						},
					},
				},
			},
			SchemaVersion: 1,
		},
	}
}

type alertingChannelResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *alertingChannelResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for alerting channels
func (r *alertingChannelResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

// SetComputedFields sets computed fields in the plan
func (r *alertingChannelResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with the alerting channel data from the API
func (r *alertingChannelResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, alertingChannel *restapi.AlertingChannel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create base model with common fields
	model := r.createBaseModel(alertingChannel)

	// Map channel-specific data based on channel type
	channelDiags := r.mapChannelTypeToModel(ctx, alertingChannel, &model)
	if channelDiags.HasError() {
		return channelDiags
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// createBaseModel creates the base model with common fields (ID and Name)
func (r *alertingChannelResource) createBaseModel(alertingChannel *restapi.AlertingChannel) AlertingChannelModel {
	return AlertingChannelModel{
		ID:   types.StringValue(alertingChannel.ID),
		Name: types.StringValue(alertingChannel.Name),
	}
}

// mapChannelTypeToModel maps the API channel data to the appropriate model field based on channel type
func (r *alertingChannelResource) mapChannelTypeToModel(ctx context.Context, alertingChannel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	var diags diag.Diagnostics

	switch alertingChannel.Kind {
	case restapi.EmailChannelType:
		return r.mapEmailToModel(ctx, alertingChannel, model)
	case restapi.OpsGenieChannelType:
		return r.mapOpsGenieToModel(ctx, alertingChannel, model)
	case restapi.PagerDutyChannelType:
		return r.mapPagerDutyToModel(ctx, alertingChannel, model)
	case restapi.SlackChannelType:
		return r.mapSlackToModel(ctx, alertingChannel, model)
	case restapi.SplunkChannelType:
		return r.mapSplunkToModel(ctx, alertingChannel, model)
	case restapi.VictorOpsChannelType:
		return r.mapVictorOpsToModel(ctx, alertingChannel, model)
	case restapi.WebhookChannelType:
		return r.mapWebhookToModel(ctx, alertingChannel, model)
	case restapi.Office365ChannelType:
		return r.mapOffice365ToModel(ctx, alertingChannel, model)
	case restapi.GoogleChatChannelType:
		return r.mapGoogleChatToModel(ctx, alertingChannel, model)
	case restapi.ServiceNowChannelType:
		return r.mapServiceNowToModel(ctx, alertingChannel, model)
	case restapi.ServiceNowApplicationChannelType:
		return r.mapServiceNowApplicationToModel(ctx, alertingChannel, model)
	case restapi.PrometheusWebhookChannelType:
		return r.mapPrometheusWebhookToModel(ctx, alertingChannel, model)
	case restapi.WebexTeamsWebhookChannelType:
		return r.mapWebexTeamsWebhookToModel(ctx, alertingChannel, model)
	case restapi.WatsonAIOpsWebhookChannelType:
		return r.mapWatsonAIOpsWebhookToModel(ctx, alertingChannel, model)
	case restapi.SlackAppChannelType:
		return r.mapSlackAppToModel(ctx, alertingChannel, model)
	case restapi.MsTeamsAppChannelType:
		return r.mapMsTeamsAppToModel(ctx, alertingChannel, model)
	default:
		diags.AddError(
			AlertingChannelErrUnsupportedType,
			fmt.Sprintf(AlertingChannelErrUnsupportedTypeMsg, alertingChannel.Kind),
		)
		return diags
	}
}

// Individual channel type mapping methods
func (r *alertingChannelResource) mapEmailToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	emailChannel, diags := shared.MapEmailChannelToState(ctx, channel)
	if !diags.HasError() {
		model.Email = emailChannel
	}
	return diags
}

func (r *alertingChannelResource) mapOpsGenieToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	opsGenieChannel, diags := shared.MapOpsGenieChannelToState(ctx, channel)
	if !diags.HasError() {
		model.OpsGenie = opsGenieChannel
	}
	return diags
}

func (r *alertingChannelResource) mapPagerDutyToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	pagerDutyChannel, diags := shared.MapPagerDutyChannelToState(ctx, channel)
	if !diags.HasError() {
		model.PagerDuty = pagerDutyChannel
	}
	return diags
}

func (r *alertingChannelResource) mapSlackToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	slackChannel, diags := shared.MapSlackChannelToState(ctx, channel)
	if !diags.HasError() {
		model.Slack = slackChannel
	}
	return diags
}

func (r *alertingChannelResource) mapSplunkToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	splunkChannel, diags := shared.MapSplunkChannelToState(ctx, channel)
	if !diags.HasError() {
		model.Splunk = splunkChannel
	}
	return diags
}

func (r *alertingChannelResource) mapVictorOpsToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	victorOpsChannel, diags := shared.MapVictorOpsChannelToState(ctx, channel)
	if !diags.HasError() {
		model.VictorOps = victorOpsChannel
	}
	return diags
}

func (r *alertingChannelResource) mapWebhookToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	webhookChannel, diags := shared.MapWebhookChannelToState(ctx, channel)
	if !diags.HasError() {
		model.Webhook = webhookChannel
	}
	return diags
}

func (r *alertingChannelResource) mapOffice365ToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	office365Channel, diags := shared.MapWebhookBasedChannelToState(ctx, channel)
	if !diags.HasError() {
		model.Office365 = office365Channel
	}
	return diags
}

func (r *alertingChannelResource) mapGoogleChatToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	googleChatChannel, diags := shared.MapWebhookBasedChannelToState(ctx, channel)
	if !diags.HasError() {
		model.GoogleChat = googleChatChannel
	}
	return diags
}

func (r *alertingChannelResource) mapServiceNowToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	serviceNowChannel, diags := shared.MapServiceNowChannelToState(ctx, channel)
	if !diags.HasError() {
		model.ServiceNow = serviceNowChannel
	}
	return diags
}

func (r *alertingChannelResource) mapServiceNowApplicationToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	serviceNowAppChannel, diags := shared.MapServiceNowApplicationChannelToState(ctx, channel)
	if !diags.HasError() {
		model.ServiceNowApplication = serviceNowAppChannel
	}
	return diags
}

func (r *alertingChannelResource) mapPrometheusWebhookToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	prometheusChannel, diags := shared.MapPrometheusWebhookChannelToState(ctx, channel)
	if !diags.HasError() {
		model.PrometheusWebhook = prometheusChannel
	}
	return diags
}

func (r *alertingChannelResource) mapWebexTeamsWebhookToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	webexChannel, diags := shared.MapWebhookBasedChannelToState(ctx, channel)
	if !diags.HasError() {
		model.WebexTeamsWebhook = webexChannel
	}
	return diags
}

func (r *alertingChannelResource) mapWatsonAIOpsWebhookToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	watsonChannel, diags := shared.MapWatsonAIOpsWebhookChannelToState(ctx, channel)
	if !diags.HasError() {
		model.WatsonAIOpsWebhook = watsonChannel
	}
	return diags
}

func (r *alertingChannelResource) mapSlackAppToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	slackAppChannel, diags := shared.MapSlackAppChannelToState(ctx, channel)
	if !diags.HasError() {
		model.SlackApp = slackAppChannel
	}
	return diags
}

func (r *alertingChannelResource) mapMsTeamsAppToModel(ctx context.Context, channel *restapi.AlertingChannel, model *AlertingChannelModel) diag.Diagnostics {
	msTeamsAppChannel, diags := shared.MapMsTeamsAppChannelToState(ctx, channel)
	if !diags.HasError() {
		model.MsTeamsApp = msTeamsAppChannel
	}
	return diags
}

// ============================================================================
// Channel Mapping Methods: State to API
// These methods convert Terraform state models to API objects
// ============================================================================

// mapEmailChannelFromState converts Email channel state to API object
func (r *alertingChannelResource) mapEmailChannelFromState(ctx context.Context, id string, name string, email *shared.EmailModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert emails set to string slice
	var emails []string
	diags.Append(email.Emails.ElementsAs(ctx, &emails, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	return &restapi.AlertingChannel{
		ID:     id,
		Name:   name,
		Kind:   restapi.EmailChannelType,
		Emails: emails,
	}, nil
}

// mapOpsGenieChannelFromState converts OpsGenie channel state to API object
func (r *alertingChannelResource) mapOpsGenieChannelFromState(ctx context.Context, id string, name string, opsGenie *shared.OpsGenieModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert tags list to string slice
	var tags []string
	diags.Append(opsGenie.Tags.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Join tags into comma-separated string
	tagsString := strings.Join(tags, TagSeparator)

	// Create alerting channel
	apiKeyValue := opsGenie.APIKey.ValueString()
	regionValue := opsGenie.Region.ValueString()

	return &restapi.AlertingChannel{
		ID:     id,
		Name:   name,
		Kind:   restapi.OpsGenieChannelType,
		APIKey: &apiKeyValue,
		Region: &regionValue,
		Tags:   &tagsString,
	}, nil
}

// mapPagerDutyChannelFromState converts PagerDuty channel state to API object
func (r *alertingChannelResource) mapPagerDutyChannelFromState(ctx context.Context, id string, name string, pagerDuty *shared.PagerDutyModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	serviceIntegrationKeyValue := pagerDuty.ServiceIntegrationKey.ValueString()

	return &restapi.AlertingChannel{
		ID:                    id,
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKeyValue,
	}, nil
}

// mapSlackChannelFromState converts Slack channel state to API object
func (r *alertingChannelResource) mapSlackChannelFromState(ctx context.Context, id string, name string, slack *shared.SlackModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	webhookURLValue := slack.WebhookURL.ValueString()

	result := &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURLValue,
	}

	// Add optional fields if present
	if !slack.IconURL.IsNull() {
		iconURLValue := slack.IconURL.ValueString()
		result.IconURL = &iconURLValue
	}

	if !slack.Channel.IsNull() {
		channelValue := slack.Channel.ValueString()
		result.Channel = &channelValue
	}

	return result, nil
}

// mapSplunkChannelFromState converts Splunk channel state to API object
func (r *alertingChannelResource) mapSplunkChannelFromState(ctx context.Context, id string, name string, splunk *shared.SplunkModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	urlValue := splunk.URL.ValueString()
	tokenValue := splunk.Token.ValueString()

	return &restapi.AlertingChannel{
		ID:    id,
		Name:  name,
		Kind:  restapi.SplunkChannelType,
		URL:   &urlValue,
		Token: &tokenValue,
	}, nil
}

// mapVictorOpsChannelFromState converts VictorOps channel state to API object
func (r *alertingChannelResource) mapVictorOpsChannelFromState(ctx context.Context, id string, name string, victorOps *shared.VictorOpsModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	apiKeyValue := victorOps.APIKey.ValueString()
	routingKeyValue := victorOps.RoutingKey.ValueString()

	return &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKeyValue,
		RoutingKey: &routingKeyValue,
	}, nil
}

// mapWebhookChannelFromState converts Webhook channel state to API object
func (r *alertingChannelResource) mapWebhookChannelFromState(ctx context.Context, id string, name string, webhook *shared.WebhookModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert webhook URLs set to string slice
	var webhookURLs []string
	diags.Append(webhook.WebhookURLs.ElementsAs(ctx, &webhookURLs, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	result := &restapi.AlertingChannel{
		ID:          id,
		Name:        name,
		Kind:        restapi.WebhookChannelType,
		WebhookURLs: webhookURLs,
	}

	// Add HTTP headers if present
	if !webhook.HTTPHeaders.IsNull() && !webhook.HTTPHeaders.IsUnknown() {
		var httpHeaders map[string]string
		diags.Append(webhook.HTTPHeaders.ElementsAs(ctx, &httpHeaders, false)...)
		if diags.HasError() {
			return nil, diags
		}

		// Convert map to header list
		headers := make([]string, 0, len(httpHeaders))
		for key, value := range httpHeaders {
			headers = append(headers, fmt.Sprintf("%s%s%s", key, HeaderSeparator, value))
		}

		result.Headers = headers
	}

	return result, nil
}

// mapWebhookBasedChannelFromState converts webhook-based channel state to API object
// Used for Office365, GoogleChat, and WebexTeams channels
func (r *alertingChannelResource) mapWebhookBasedChannelFromState(ctx context.Context, id string, name string, webhookBased *shared.WebhookBasedModel, channelType restapi.AlertingChannelType) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	webhookURLValue := webhookBased.WebhookURL.ValueString()

	return &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       channelType,
		WebhookURL: &webhookURLValue,
	}, nil
}

// mapServiceNowChannelFromState converts ServiceNow channel state to API object
func (r *alertingChannelResource) mapServiceNowChannelFromState(ctx context.Context, id string, name string, serviceNow *shared.ServiceNowModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if serviceNow.Password.IsNull() || serviceNow.Password.IsUnknown() {
		diags.AddError(AlertingChannelErrMissingPassword, AlertingChannelErrMissingPasswordMsg)
		return nil, diags
	}

	serviceNowURLValue := serviceNow.ServiceNowURL.ValueString()
	usernameValue := serviceNow.Username.ValueString()
	passwordValue := serviceNow.Password.ValueString()

	result := &restapi.AlertingChannel{
		ID:            id,
		Name:          name,
		Kind:          restapi.ServiceNowChannelType,
		ServiceNowURL: &serviceNowURLValue,
		Username:      &usernameValue,
		Password:      &passwordValue,
	}

	if !serviceNow.AutoCloseIncidents.IsNull() {
		autoCloseValue := serviceNow.AutoCloseIncidents.ValueBool()
		result.AutoCloseIncidents = &autoCloseValue
	}

	return result, nil
}

// mapServiceNowApplicationChannelFromState converts ServiceNow Enhanced (ITSM) channel state to API object
func (r *alertingChannelResource) mapServiceNowApplicationChannelFromState(ctx context.Context, id string, name string, serviceNowApp *shared.ServiceNowApplicationModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if serviceNowApp.Password.IsNull() || serviceNowApp.Password.IsUnknown() {
		diags.AddError(AlertingChannelErrMissingPassword, AlertingChannelErrMissingPasswordMsg)
		return nil, diags
	}

	serviceNowURLValue := serviceNowApp.ServiceNowURL.ValueString()
	usernameValue := serviceNowApp.Username.ValueString()
	passwordValue := serviceNowApp.Password.ValueString()
	tenantValue := serviceNowApp.Tenant.ValueString()
	unitValue := serviceNowApp.Unit.ValueString()

	result := &restapi.AlertingChannel{
		ID:            id,
		Name:          name,
		Kind:          restapi.ServiceNowApplicationChannelType,
		ServiceNowURL: &serviceNowURLValue,
		Username:      &usernameValue,
		Password:      &passwordValue,
		Tenant:        &tenantValue,
		Unit:          &unitValue,
	}

	// Add optional fields
	if !serviceNowApp.AutoCloseIncidents.IsNull() {
		autoCloseValue := serviceNowApp.AutoCloseIncidents.ValueBool()
		result.AutoCloseIncidents = &autoCloseValue
	}

	if !serviceNowApp.InstanaURL.IsNull() && !serviceNowApp.InstanaURL.IsUnknown() {
		instanaURLValue := serviceNowApp.InstanaURL.ValueString()
		result.InstanaURL = &instanaURLValue
	} else {
		diags.AddError(
			AlertingChannelErrInstanaURLRequired,
			AlertingChannelErrInstanaURLRequiredMsg,
		)
		return result, diags
	}

	if !serviceNowApp.EnableSendInstanaNotes.IsNull() {
		enableSendInstanaNotesValue := serviceNowApp.EnableSendInstanaNotes.ValueBool()
		result.EnableSendInstanaNotes = &enableSendInstanaNotesValue
	}

	if !serviceNowApp.EnableSendServiceNowActivities.IsNull() {
		enableSendServiceNowActivitiesValue := serviceNowApp.EnableSendServiceNowActivities.ValueBool()
		result.EnableSendServiceNowActivities = &enableSendServiceNowActivitiesValue
	}

	if !serviceNowApp.EnableSendServiceNowWorkNotes.IsNull() {
		enableSendServiceNowWorkNotesValue := serviceNowApp.EnableSendServiceNowWorkNotes.ValueBool()
		result.EnableSendServiceNowWorkNotes = &enableSendServiceNowWorkNotesValue
	}

	if !serviceNowApp.ManuallyClosedIncidents.IsNull() {
		manuallyClosedIncidentsValue := serviceNowApp.ManuallyClosedIncidents.ValueBool()
		result.ManuallyClosedIncidents = &manuallyClosedIncidentsValue
	}

	if !serviceNowApp.ResolutionOfIncident.IsNull() {
		resolutionOfIncidentValue := serviceNowApp.ResolutionOfIncident.ValueBool()
		result.ResolutionOfIncident = &resolutionOfIncidentValue
	}

	if !serviceNowApp.SnowStatusOnCloseEvent.IsNull() {
		snowStatusValue := int(serviceNowApp.SnowStatusOnCloseEvent.ValueInt64())
		result.SnowStatusOnCloseEvent = &snowStatusValue
	}

	return result, nil
}

// mapPrometheusWebhookChannelFromState converts Prometheus Webhook channel state to API object
func (r *alertingChannelResource) mapPrometheusWebhookChannelFromState(ctx context.Context, id string, name string, prometheusWebhook *shared.PrometheusWebhookModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	webhookURLValue := prometheusWebhook.WebhookURL.ValueString()

	result := &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       restapi.PrometheusWebhookChannelType,
		WebhookURL: &webhookURLValue,
	}

	if !prometheusWebhook.Receiver.IsNull() {
		receiverValue := prometheusWebhook.Receiver.ValueString()
		result.Receiver = &receiverValue
	}

	return result, nil
}

// mapWatsonAIOpsWebhookChannelFromState converts Watson AIOps Webhook channel state to API object
func (r *alertingChannelResource) mapWatsonAIOpsWebhookChannelFromState(ctx context.Context, id string, name string, watsonAIOpsWebhook *shared.WatsonAIOpsWebhookModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	webhookURLValue := watsonAIOpsWebhook.WebhookURL.ValueString()

	result := &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       restapi.WatsonAIOpsWebhookChannelType,
		WebhookURL: &webhookURLValue,
	}

	// Add headers if present
	if !watsonAIOpsWebhook.HTTPHeaders.IsNull() && !watsonAIOpsWebhook.HTTPHeaders.IsUnknown() {
		var headers []string
		diags.Append(watsonAIOpsWebhook.HTTPHeaders.ElementsAs(ctx, &headers, false)...)
		if diags.HasError() {
			return nil, diags
		}
		result.Headers = headers
	}

	return result, nil
}

// mapSlackAppChannelFromState converts Slack App channel state to API object
func (r *alertingChannelResource) mapSlackAppChannelFromState(ctx context.Context, id string, name string, slackApp *shared.SlackAppModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	appIDValue := slackApp.AppID.ValueString()
	teamIDValue := slackApp.TeamID.ValueString()
	teamNameValue := slackApp.TeamName.ValueString()
	channelIDValue := slackApp.ChannelID.ValueString()
	channelNameValue := slackApp.ChannelName.ValueString()

	result := &restapi.AlertingChannel{
		ID:          id,
		Name:        name,
		Kind:        restapi.SlackAppChannelType,
		AppID:       &appIDValue,
		TeamID:      &teamIDValue,
		TeamName:    &teamNameValue,
		ChannelID:   &channelIDValue,
		ChannelName: &channelNameValue,
	}

	// Add optional emoji rendering field
	if !slackApp.EmojiRendering.IsNull() {
		emojiRenderingValue := slackApp.EmojiRendering.ValueBool()
		result.EmojiRendering = &emojiRenderingValue
	}

	return result, nil
}

// mapMsTeamsAppChannelFromState converts MS Teams App channel state to API object
func (r *alertingChannelResource) mapMsTeamsAppChannelFromState(ctx context.Context, id string, name string, msTeamsApp *shared.MsTeamsAppModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	apiTokenIDValue := msTeamsApp.APITokenID.ValueString()
	teamIDValue := msTeamsApp.TeamID.ValueString()
	teamNameValue := msTeamsApp.TeamName.ValueString()
	channelIDValue := msTeamsApp.ChannelID.ValueString()
	channelNameValue := msTeamsApp.ChannelName.ValueString()
	instanaURLValue := msTeamsApp.InstanaURL.ValueString()
	serviceURLValue := msTeamsApp.ServiceURL.ValueString()
	tenantIDValue := msTeamsApp.TenantID.ValueString()
	tenantNameValue := msTeamsApp.TenantName.ValueString()

	result := &restapi.AlertingChannel{
		ID:          id,
		Name:        name,
		Kind:        restapi.MsTeamsAppChannelType,
		APITokenID:  &apiTokenIDValue,
		TeamID:      &teamIDValue,
		TeamName:    &teamNameValue,
		ChannelID:   &channelIDValue,
		ChannelName: &channelNameValue,
		InstanaURL:  &instanaURLValue,
		ServiceURL:  &serviceURLValue,
		TenantID:    &tenantIDValue,
		TenantName:  &tenantNameValue,
	}

	return result, nil
}

// ============================================================================
// Main Mapping Method
// ============================================================================

// MapStateToDataObject converts Terraform state to API object
// This method determines which channel type is configured and delegates to the appropriate mapper
func (r *alertingChannelResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Get model from plan or state
	model, diags := r.getModelFromPlanOrState(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	// Extract common fields
	id, name := r.extractCommonFields(model)

	// Map the configured channel type to API object
	return r.mapConfiguredChannelType(ctx, model, id, name)
}

// getModelFromPlanOrState retrieves the model from either plan or state
func (r *alertingChannelResource) getModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (AlertingChannelModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AlertingChannelModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractCommonFields extracts ID and name from the model
func (r *alertingChannelResource) extractCommonFields(model AlertingChannelModel) (string, string) {
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}
	name := model.Name.ValueString()
	return id, name
}

// mapConfiguredChannelType determines which channel type is configured and maps it to API object
func (r *alertingChannelResource) mapConfiguredChannelType(ctx context.Context, model AlertingChannelModel, id, name string) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check each channel type and map accordingly
	if model.Email != nil {
		return r.mapEmailChannelFromState(ctx, id, name, model.Email)
	}
	if model.OpsGenie != nil {
		return r.mapOpsGenieChannelFromState(ctx, id, name, model.OpsGenie)
	}
	if model.PagerDuty != nil {
		return r.mapPagerDutyChannelFromState(ctx, id, name, model.PagerDuty)
	}
	if model.Slack != nil {
		return r.mapSlackChannelFromState(ctx, id, name, model.Slack)
	}
	if model.Splunk != nil {
		return r.mapSplunkChannelFromState(ctx, id, name, model.Splunk)
	}
	if model.VictorOps != nil {
		return r.mapVictorOpsChannelFromState(ctx, id, name, model.VictorOps)
	}
	if model.Webhook != nil {
		return r.mapWebhookChannelFromState(ctx, id, name, model.Webhook)
	}
	if model.Office365 != nil {
		return r.mapWebhookBasedChannelFromState(ctx, id, name, model.Office365, restapi.Office365ChannelType)
	}
	if model.GoogleChat != nil {
		return r.mapWebhookBasedChannelFromState(ctx, id, name, model.GoogleChat, restapi.GoogleChatChannelType)
	}
	if model.ServiceNow != nil {
		return r.mapServiceNowChannelFromState(ctx, id, name, model.ServiceNow)
	}
	if model.ServiceNowApplication != nil {
		return r.mapServiceNowApplicationChannelFromState(ctx, id, name, model.ServiceNowApplication)
	}
	if model.PrometheusWebhook != nil {
		return r.mapPrometheusWebhookChannelFromState(ctx, id, name, model.PrometheusWebhook)
	}
	if model.WebexTeamsWebhook != nil {
		return r.mapWebhookBasedChannelFromState(ctx, id, name, model.WebexTeamsWebhook, restapi.WebexTeamsWebhookChannelType)
	}
	if model.WatsonAIOpsWebhook != nil {
		return r.mapWatsonAIOpsWebhookChannelFromState(ctx, id, name, model.WatsonAIOpsWebhook)
	}
	if model.SlackApp != nil {
		return r.mapSlackAppChannelFromState(ctx, id, name, model.SlackApp)
	}
	if model.MsTeamsApp != nil {
		return r.mapMsTeamsAppChannelFromState(ctx, id, name, model.MsTeamsApp)
	}

	// No valid channel type configured
	diags.AddError(
		AlertingChannelErrInvalidConfig,
		AlertingChannelErrInvalidConfigMsg,
	)
	return nil, diags
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *alertingChannelResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: resourcehandle.CreateStateUpgraderForVersion(0),
	}
}
