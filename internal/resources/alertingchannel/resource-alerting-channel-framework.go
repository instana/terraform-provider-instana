package alertingchannel

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewAlertingChannelResourceHandleFramework creates the resource handle for Alerting Channels
func NewAlertingChannelResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.AlertingChannel] {
	supportedOpsGenieRegions := []string{"EU", "US"}
	return &alertingChannelResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAlertingChannelFramework,
			Schema: schema.Schema{
				Description: "This resource manages alerting channels in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the alerting channel.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AlertingChannelFieldName: schema.StringAttribute{
						Required:    true,
						Description: "Configures the name of the alerting channel",
					},
					AlertingChannelFieldChannelEmail: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Email channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelEmailFieldEmails: schema.SetAttribute{
								Required:    true,
								Description: "The list of emails of the Email alerting channel",
								ElementType: types.StringType,
							},
						},
					},
					AlertingChannelFieldChannelOpsGenie: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Ops Genie channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelOpsGenieFieldAPIKey: schema.StringAttribute{
								Required:    true,
								Description: "The OpsGenie API Key of the OpsGenie alerting channel",
							},
							AlertingChannelOpsGenieFieldTags: schema.ListAttribute{
								Required:    true,
								Description: "The OpsGenie tags of the OpsGenie alerting channel",
								ElementType: types.StringType,
							},
							AlertingChannelOpsGenieFieldRegion: schema.StringAttribute{
								Required:    true,
								Description: fmt.Sprintf("The OpsGenie region (%s) of the OpsGenie alerting channel", strings.Join(supportedOpsGenieRegions, ", ")),
								Validators: []validator.String{
									stringvalidator.OneOf(supportedOpsGenieRegions...),
								},
							},
						},
					},
					AlertingChannelFieldChannelPageDuty: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Pager Duty channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelPagerDutyFieldServiceIntegrationKey: schema.StringAttribute{
								Required:    true,
								Description: "The Service Integration Key of the PagerDuty alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelSlack: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Slack channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelSlackFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: "The webhook URL of the Slack alerting channel",
							},
							AlertingChannelSlackFieldIconURL: schema.StringAttribute{
								Optional:    true,
								Description: "The icon URL of the Slack alerting channel",
							},
							AlertingChannelSlackFieldChannel: schema.StringAttribute{
								Optional:    true,
								Description: "The Slack channel of the Slack alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelSplunk: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Splunk channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelSplunkFieldURL: schema.StringAttribute{
								Required:    true,
								Description: "The URL of the Splunk alerting channel",
							},
							AlertingChannelSplunkFieldToken: schema.StringAttribute{
								Required:    true,
								Description: "The token of the Splunk alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelVictorOps: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the VictorOps channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelVictorOpsFieldAPIKey: schema.StringAttribute{
								Required:    true,
								Description: "The API Key of the VictorOps alerting channel",
							},
							AlertingChannelVictorOpsFieldRoutingKey: schema.StringAttribute{
								Required:    true,
								Description: "The Routing Key of the VictorOps alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Webhook channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookFieldWebhookURLs: schema.SetAttribute{
								Required:    true,
								Description: "The list of webhook urls of the Webhook alerting channel",
								ElementType: types.StringType,
							},
							AlertingChannelWebhookFieldHTTPHeaders: schema.MapAttribute{
								Optional:    true,
								Description: "The optional map of HTTP headers of the Webhook alerting channel",
								ElementType: types.StringType,
							},
						},
					},
					AlertingChannelFieldChannelOffice365: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Office 365 channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: "The webhook URL of the Office 365 alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelGoogleChat: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Google Chat channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: "The webhook URL of the Google Chat alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelServiceNow: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the ServiceNow channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
								Required:    true,
								Description: "The ServiceNow URL of the ServiceNow alerting channel",
							},
							AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
								Required:    true,
								Description: "The username of the ServiceNow alerting channel",
							},
							AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
								Optional:    true,
								Description: "The password of the ServiceNow alerting channel",
								PlanModifiers: []planmodifier.String{
									// When the plan does not include the password, keep the value from state.
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							AlertingChannelServiceNowFieldAutoCloseIncidents: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to automatically close incidents in ServiceNow",
							},
						},
					},
					AlertingChannelFieldChannelServiceNowApplication: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the ServiceNow Enhanced (ITSM) channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelServiceNowFieldServiceNowURL: schema.StringAttribute{
								Required:    true,
								Description: "The ServiceNow URL of the ServiceNow Enhanced alerting channel",
							},
							AlertingChannelServiceNowFieldUsername: schema.StringAttribute{
								Required:    true,
								Description: "The username of the ServiceNow Enhanced alerting channel",
							},
							AlertingChannelServiceNowFieldPassword: schema.StringAttribute{
								Optional:    true,
								Description: "The password of the ServiceNow Enhanced alerting channel",
								PlanModifiers: []planmodifier.String{
									// When the plan does not include the password, keep the value from state.
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							AlertingChannelServiceNowApplicationFieldTenant: schema.StringAttribute{
								Required:    true,
								Description: "The tenant of the ServiceNow Enhanced alerting channel",
							},
							AlertingChannelServiceNowApplicationFieldUnit: schema.StringAttribute{
								Required:    true,
								Description: "The unit of the ServiceNow Enhanced alerting channel",
							},
							AlertingChannelServiceNowFieldAutoCloseIncidents: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to automatically close incidents in ServiceNow",
							},
							AlertingChannelServiceNowApplicationFieldInstanaURL: schema.StringAttribute{
								Optional:    true,
								Description: "The Instana URL for the ServiceNow Enhanced alerting channel",
							},
							AlertingChannelServiceNowApplicationFieldEnableSendInstanaNotes: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to send Instana notes to ServiceNow",
							},
							AlertingChannelServiceNowApplicationFieldEnableSendServiceNowActivities: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to send ServiceNow activities",
							},
							AlertingChannelServiceNowApplicationFieldEnableSendServiceNowWorkNotes: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to send ServiceNow work notes",
							},
							AlertingChannelServiceNowApplicationFieldManuallyClosedIncidents: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether incidents are manually closed",
							},
							AlertingChannelServiceNowApplicationFieldResolutionOfIncident: schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to resolve incidents",
							},
							AlertingChannelServiceNowApplicationFieldSnowStatusOnCloseEvent: schema.Int64Attribute{
								Optional:    true,
								Description: "The ServiceNow status code when closing events",
							},
						},
					},
					AlertingChannelFieldChannelPrometheusWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Prometheus Webhook channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: "The webhook URL of the Prometheus Webhook alerting channel",
							},
							AlertingChannelPrometheusWebhookFieldReceiver: schema.StringAttribute{
								Optional:    true,
								Description: "The receiver of the Prometheus Webhook alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelWebexTeamsWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Webex Teams Webhook channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: "The webhook URL of the Webex Teams Webhook alerting channel",
							},
						},
					},
					AlertingChannelFieldChannelWatsonAIOpsWebhook: schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The configuration of the Watson AIOps Webhook channel",
						Attributes: map[string]schema.Attribute{
							AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
								Required:    true,
								Description: "The webhook URL of the Watson AIOps Webhook alerting channel",
							},
							AlertingChannelWebhookFieldHTTPHeaders: schema.ListAttribute{
								Optional:    true,
								Description: "The list of HTTP headers for the Watson AIOps Webhook alerting channel",
								ElementType: types.StringType,
							},
						},
					},
				},
			},
			SchemaVersion: 1,
		},
	}
}

type alertingChannelResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *alertingChannelResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *alertingChannelResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *alertingChannelResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, alertingChannel *restapi.AlertingChannel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get the plan model to preserve optional values in the response
	var planModel AlertingChannelModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &planModel)...)
		if diags.HasError() {
			return diags
		}
	}

	// Create a model and populate it with values from the alerting channel
	model := AlertingChannelModel{
		ID:   types.StringValue(alertingChannel.ID),
		Name: types.StringValue(alertingChannel.Name),
	}

	// Initialize all channel types as nil (they are pointer models)
	// Only the matching channel type will be populated

	// Set the appropriate channel type based on the alerting channel kind
	switch alertingChannel.Kind {
	case restapi.EmailChannelType:
		emailChannel, emailDiags := r.mapEmailChannelToState(ctx, alertingChannel)
		if emailDiags.HasError() {
			diags.Append(emailDiags...)
			return diags
		}
		model.Email = emailChannel
	case restapi.OpsGenieChannelType:
		opsGenieChannel, opsGenieDiags := r.mapOpsGenieChannelToState(ctx, alertingChannel)
		if opsGenieDiags.HasError() {
			diags.Append(opsGenieDiags...)
			return diags
		}
		model.OpsGenie = opsGenieChannel
	case restapi.PagerDutyChannelType:
		pagerDutyChannel, pagerDutyDiags := r.mapPagerDutyChannelToState(ctx, alertingChannel)
		if pagerDutyDiags.HasError() {
			diags.Append(pagerDutyDiags...)
			return diags
		}
		model.PagerDuty = pagerDutyChannel
	case restapi.SlackChannelType:
		slackChannel, slackDiags := r.mapSlackChannelToState(ctx, alertingChannel)
		if slackDiags.HasError() {
			diags.Append(slackDiags...)
			return diags
		}
		model.Slack = slackChannel
	case restapi.SplunkChannelType:
		splunkChannel, splunkDiags := r.mapSplunkChannelToState(ctx, alertingChannel)
		if splunkDiags.HasError() {
			diags.Append(splunkDiags...)
			return diags
		}
		model.Splunk = splunkChannel
	case restapi.VictorOpsChannelType:
		victorOpsChannel, victorOpsDiags := r.mapVictorOpsChannelToState(ctx, alertingChannel)
		if victorOpsDiags.HasError() {
			diags.Append(victorOpsDiags...)
			return diags
		}
		model.VictorOps = victorOpsChannel
	case restapi.WebhookChannelType:
		webhookChannel, webhookDiags := r.mapWebhookChannelToState(ctx, alertingChannel)
		if webhookDiags.HasError() {
			diags.Append(webhookDiags...)
			return diags
		}
		model.Webhook = webhookChannel
	case restapi.Office365ChannelType:
		office365Channel, office365Diags := r.mapWebhookBasedChannelToState(ctx, alertingChannel)
		if office365Diags.HasError() {
			diags.Append(office365Diags...)
			return diags
		}
		model.Office365 = office365Channel
	case restapi.GoogleChatChannelType:
		googleChatChannel, googleChatDiags := r.mapWebhookBasedChannelToState(ctx, alertingChannel)
		if googleChatDiags.HasError() {
			diags.Append(googleChatDiags...)
			return diags
		}
		model.GoogleChat = googleChatChannel
	case restapi.ServiceNowChannelType:
		serviceNowChannel, serviceNowDiags := r.mapServiceNowChannelToState(ctx, alertingChannel)
		if serviceNowDiags.HasError() {
			diags.Append(serviceNowDiags...)
			return diags
		}
		model.ServiceNow = serviceNowChannel
	case restapi.ServiceNowApplicationChannelType:
		serviceNowEnhancedChannel, serviceNowEnhancedDiags := r.mapServiceNowApplicationChannelToState(ctx, alertingChannel)
		if serviceNowEnhancedDiags.HasError() {
			diags.Append(serviceNowEnhancedDiags...)
			return diags
		}
		model.ServiceNowApplication = serviceNowEnhancedChannel
	case restapi.PrometheusWebhookChannelType:
		prometheusWebhookChannel, prometheusWebhookDiags := r.mapPrometheusWebhookChannelToState(ctx, alertingChannel)
		if prometheusWebhookDiags.HasError() {
			diags.Append(prometheusWebhookDiags...)
			return diags
		}
		model.PrometheusWebhook = prometheusWebhookChannel
	case restapi.WebexTeamsWebhookChannelType:
		webexTeamsWebhookChannel, webexTeamsWebhookDiags := r.mapWebhookBasedChannelToState(ctx, alertingChannel)
		if webexTeamsWebhookDiags.HasError() {
			diags.Append(webexTeamsWebhookDiags...)
			return diags
		}
		model.WebexTeamsWebhook = webexTeamsWebhookChannel
	case restapi.WatsonAIOpsWebhookChannelType:
		watsonAIOpsWebhookChannel, watsonAIOpsWebhookDiags := r.mapWatsonAIOpsWebhookChannelToState(ctx, alertingChannel)
		if watsonAIOpsWebhookDiags.HasError() {
			diags.Append(watsonAIOpsWebhookDiags...)
			return diags
		}
		model.WatsonAIOpsWebhook = watsonAIOpsWebhookChannel
	default:
		diags.AddError(
			"Unsupported alerting channel type",
			fmt.Sprintf("Received unsupported alerting channel of type %s", alertingChannel.Kind),
		)
		return diags
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *alertingChannelResourceFramework) mapEmailChannelFromState(ctx context.Context, id string, name string, email *EmailModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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

func (r *alertingChannelResourceFramework) mapOpsGenieChannelFromState(ctx context.Context, id string, name string, opsGenie *OpsGenieModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert tags list to string slice
	var tags []string
	diags.Append(opsGenie.Tags.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Join tags into comma-separated string
	tagsString := strings.Join(tags, ",")

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

func (r *alertingChannelResourceFramework) mapPagerDutyChannelFromState(ctx context.Context, id string, name string, pagerDuty *PagerDutyModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	serviceIntegrationKeyValue := pagerDuty.ServiceIntegrationKey.ValueString()

	return &restapi.AlertingChannel{
		ID:                    id,
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKeyValue,
	}, nil
}

func (r *alertingChannelResourceFramework) mapSlackChannelFromState(ctx context.Context, id string, name string, slack *SlackModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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

func (r *alertingChannelResourceFramework) mapSplunkChannelFromState(ctx context.Context, id string, name string, splunk *SplunkModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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

func (r *alertingChannelResourceFramework) mapVictorOpsChannelFromState(ctx context.Context, id string, name string, victorOps *VictorOpsModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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

func (r *alertingChannelResourceFramework) mapWebhookChannelFromState(ctx context.Context, id string, name string, webhook *WebhookModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}

		result.Headers = headers
	}

	return result, nil
}

func (r *alertingChannelResourceFramework) mapWebhookBasedChannelFromState(ctx context.Context, id string, name string, webhookBased *WebhookBasedModel, channelType restapi.AlertingChannelType) (*restapi.AlertingChannel, diag.Diagnostics) {
	// Create alerting channel
	webhookURLValue := webhookBased.WebhookURL.ValueString()

	return &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       channelType,
		WebhookURL: &webhookURLValue,
	}, nil
}

func (r *alertingChannelResourceFramework) mapServiceNowChannelFromState(ctx context.Context, id string, name string, serviceNow *ServiceNowModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if serviceNow.Password.IsNull() || serviceNow.Password.IsUnknown() {
		diags.AddError("Missing Password", "password must be specified when creating the resource")
		return nil, diags
	}
	log.Printf("passwordValue: %s", serviceNow.Password.ValueString())

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

func (r *alertingChannelResourceFramework) mapServiceNowApplicationChannelFromState(ctx context.Context, id string, name string, serviceNowApp *ServiceNowApplicationModel) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if serviceNowApp.Password.IsNull() || serviceNowApp.Password.IsUnknown() {
		diags.AddError("Missing Password", "password must be specified when creating the resource")
		return nil, diags
	}

	log.Printf("passwordValue: %s", serviceNowApp.Password.ValueString())

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

	log.Printf("Inatna url : user %v", serviceNowApp.InstanaURL.ValueString())
	if !serviceNowApp.InstanaURL.IsNull() && !serviceNowApp.InstanaURL.IsUnknown() {
		instanaURLValue := serviceNowApp.InstanaURL.ValueString()
		result.InstanaURL = &instanaURLValue
	} else {
		diags.AddError(
			"InstanaURL is required",
			"InstanaURL is required when creating the resource",
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

	log.Printf("[DEBUG] mapServiceNowApplicationChannelFromState: %v", result)
	log.Printf("[DEBUG] intana url: %v", *result.InstanaURL)
	return result, nil
}

func (r *alertingChannelResourceFramework) mapPrometheusWebhookChannelFromState(ctx context.Context, id string, name string, prometheusWebhook *PrometheusWebhookModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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

func (r *alertingChannelResourceFramework) mapWatsonAIOpsWebhookChannelFromState(ctx context.Context, id string, name string, watsonAIOpsWebhook *WatsonAIOpsWebhookModel) (*restapi.AlertingChannel, diag.Diagnostics) {
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

func (r *alertingChannelResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AlertingChannelModel

	// Get current state from plan or state
	if plan != nil {
		log.Printf("Model from plan")
		diags.Append(plan.Get(ctx, &model)...)
	} else {
		log.Printf("Model from state")
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map name
	name := model.Name.ValueString()

	// Determine which channel type is configured and map accordingly
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

	diags.AddError(
		"Invalid Alerting Channel Configuration",
		"No valid alerting channel configuration found. Please configure exactly one channel type.",
	)
	return nil, diags
}

// Made with Bob
