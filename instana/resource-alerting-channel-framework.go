package instana

import (
	"context"
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResourceInstanaAlertingChannelFramework the name of the terraform-provider-instana resource to manage alerting channels
const ResourceInstanaAlertingChannelFramework = "alerting_channel"

// AlertingChannelModel represents the data model for the alerting channel resource
type AlertingChannelModel struct {
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

// NewAlertingChannelResourceHandleFramework creates the resource handle for Alerting Channels
func NewAlertingChannelResourceHandleFramework() ResourceHandleFramework[*restapi.AlertingChannel] {
	supportedOpsGenieRegions := []string{"EU", "US"}
	return &alertingChannelResourceFramework{
		metaData: ResourceMetaDataFramework{
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
				},
				Blocks: map[string]schema.Block{
					AlertingChannelFieldChannelEmail: schema.ListNestedBlock{
						Description: "The configuration of the Email channel",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								AlertingChannelEmailFieldEmails: schema.SetAttribute{
									Required:    true,
									Description: "The list of emails of the Email alerting channel",
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
					},
					AlertingChannelFieldChannelPageDuty: schema.ListNestedBlock{
						Description: "The configuration of the Pager Duty channel",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								AlertingChannelPagerDutyFieldServiceIntegrationKey: schema.StringAttribute{
									Required:    true,
									Description: "The Service Integration Key of the PagerDuty alerting channel",
								},
							},
						},
					},
					AlertingChannelFieldChannelSlack: schema.ListNestedBlock{
						Description: "The configuration of the Slack channel",
						NestedObject: schema.NestedBlockObject{
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
					},
					AlertingChannelFieldChannelSplunk: schema.ListNestedBlock{
						Description: "The configuration of the Splunk channel",
						NestedObject: schema.NestedBlockObject{
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
					},
					AlertingChannelFieldChannelVictorOps: schema.ListNestedBlock{
						Description: "The configuration of the VictorOps channel",
						NestedObject: schema.NestedBlockObject{
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
					},
					AlertingChannelFieldChannelWebhook: schema.ListNestedBlock{
						Description: "The configuration of the Webhook channel",
						NestedObject: schema.NestedBlockObject{
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
					},
					AlertingChannelFieldChannelOffice365: schema.ListNestedBlock{
						Description: "The configuration of the Office 365 channel",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
									Required:    true,
									Description: "The webhook URL of the Office 365 alerting channel",
								},
							},
						},
					},
					AlertingChannelFieldChannelGoogleChat: schema.ListNestedBlock{
						Description: "The configuration of the Google Chat channel",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								AlertingChannelWebhookBasedFieldWebhookURL: schema.StringAttribute{
									Required:    true,
									Description: "The webhook URL of the Google Chat alerting channel",
								},
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
	metaData ResourceMetaDataFramework
}

func (r *alertingChannelResourceFramework) MetaData() *ResourceMetaDataFramework {
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

	// Create a model and populate it with values from the alerting channel
	model := AlertingChannelModel{
		ID:   types.StringValue(alertingChannel.ID),
		Name: types.StringValue(alertingChannel.Name),
	}

	// Initialize all channel types as null lists with proper attribute types
	model.Email = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelEmailFieldEmails: types.SetType{ElemType: types.StringType},
		},
	})
	model.OpsGenie = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelOpsGenieFieldAPIKey: types.StringType,
			AlertingChannelOpsGenieFieldRegion: types.StringType,
			AlertingChannelOpsGenieFieldTags:   types.ListType{ElemType: types.StringType},
		},
	})
	model.PagerDuty = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelPagerDutyFieldServiceIntegrationKey: types.StringType,
		},
	})
	model.Slack = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelSlackFieldWebhookURL: types.StringType,
			AlertingChannelSlackFieldIconURL:    types.StringType,
			AlertingChannelSlackFieldChannel:    types.StringType,
		},
	})
	model.Splunk = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelSplunkFieldURL:   types.StringType,
			AlertingChannelSplunkFieldToken: types.StringType,
		},
	})
	model.VictorOps = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelVictorOpsFieldAPIKey:     types.StringType,
			AlertingChannelVictorOpsFieldRoutingKey: types.StringType,
		},
	})
	model.Webhook = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelWebhookFieldWebhookURLs: types.SetType{ElemType: types.StringType},
			AlertingChannelWebhookFieldHTTPHeaders: types.MapType{ElemType: types.StringType},
		},
	})
	model.Office365 = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelWebhookBasedFieldWebhookURL: types.StringType,
		},
	})
	model.GoogleChat = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			AlertingChannelWebhookBasedFieldWebhookURL: types.StringType,
		},
	})

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

func (r *alertingChannelResourceFramework) mapEmailChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create email set
	emailsSet, emailsDiags := types.SetValueFrom(ctx, types.StringType, channel.Emails)
	if emailsDiags.HasError() {
		diags.Append(emailsDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create email channel object
	emailObj := map[string]attr.Value{
		AlertingChannelEmailFieldEmails: emailsSet,
	}

	emailObjVal, emailObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelEmailFieldEmails: types.SetType{ElemType: types.StringType},
		},
		emailObj,
	)
	diags.Append(emailObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelEmailFieldEmails: types.SetType{ElemType: types.StringType},
			},
		},
		[]attr.Value{emailObjVal},
	)
}

func (r *alertingChannelResourceFramework) mapOpsGenieChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert comma-separated tags to slice
	tags := r.convertCommaSeparatedListToSlice(*channel.Tags)

	// Create tags list
	tagsList, tagsDiags := types.ListValueFrom(ctx, types.StringType, tags)
	if tagsDiags.HasError() {
		diags.Append(tagsDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create OpsGenie channel object
	opsGenieObj := map[string]attr.Value{
		AlertingChannelOpsGenieFieldAPIKey: setStringPointerToState(channel.APIKey),
		AlertingChannelOpsGenieFieldRegion: setStringPointerToState(channel.Region),
		AlertingChannelOpsGenieFieldTags:   tagsList,
	}

	opsGenieObjVal, opsGenieObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelOpsGenieFieldAPIKey: types.StringType,
			AlertingChannelOpsGenieFieldRegion: types.StringType,
			AlertingChannelOpsGenieFieldTags:   types.ListType{ElemType: types.StringType},
		},
		opsGenieObj,
	)
	diags.Append(opsGenieObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelOpsGenieFieldAPIKey: types.StringType,
				AlertingChannelOpsGenieFieldRegion: types.StringType,
				AlertingChannelOpsGenieFieldTags:   types.ListType{ElemType: types.StringType},
			},
		},
		[]attr.Value{opsGenieObjVal},
	)
}

func (r *alertingChannelResourceFramework) mapPagerDutyChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create PagerDuty channel object
	pagerDutyObj := map[string]attr.Value{
		AlertingChannelPagerDutyFieldServiceIntegrationKey: setStringPointerToState(channel.ServiceIntegrationKey),
	}

	pagerDutyObjVal, pagerDutyObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelPagerDutyFieldServiceIntegrationKey: types.StringType,
		},
		pagerDutyObj,
	)
	diags.Append(pagerDutyObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelPagerDutyFieldServiceIntegrationKey: types.StringType,
			},
		},
		[]attr.Value{pagerDutyObjVal},
	)
}

func (r *alertingChannelResourceFramework) mapSlackChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create Slack channel object
	slackObj := map[string]attr.Value{
		AlertingChannelSlackFieldWebhookURL: setStringPointerToState(channel.WebhookURL),
	}

	// Add optional fields if present
	if channel.IconURL != nil && *channel.IconURL != "" {
		slackObj[AlertingChannelSlackFieldIconURL] = setStringPointerToState(channel.IconURL)
	} else {
		slackObj[AlertingChannelSlackFieldIconURL] = types.StringNull()
	}

	if channel.Channel != nil && *channel.Channel != "" {
		slackObj[AlertingChannelSlackFieldChannel] = setStringPointerToState(channel.Channel)
	} else {
		slackObj[AlertingChannelSlackFieldChannel] = types.StringNull()
	}

	slackObjVal, slackObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelSlackFieldWebhookURL: types.StringType,
			AlertingChannelSlackFieldIconURL:    types.StringType,
			AlertingChannelSlackFieldChannel:    types.StringType,
		},
		slackObj,
	)
	diags.Append(slackObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelSlackFieldWebhookURL: types.StringType,
				AlertingChannelSlackFieldIconURL:    types.StringType,
				AlertingChannelSlackFieldChannel:    types.StringType,
			},
		},
		[]attr.Value{slackObjVal},
	)
}

func (r *alertingChannelResourceFramework) mapSplunkChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create Splunk channel object
	splunkObj := map[string]attr.Value{
		AlertingChannelSplunkFieldURL:   setStringPointerToState(channel.URL),
		AlertingChannelSplunkFieldToken: setStringPointerToState(channel.Token),
	}

	splunkObjVal, splunkObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelSplunkFieldURL:   types.StringType,
			AlertingChannelSplunkFieldToken: types.StringType,
		},
		splunkObj,
	)
	diags.Append(splunkObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelSplunkFieldURL:   types.StringType,
				AlertingChannelSplunkFieldToken: types.StringType,
			},
		},
		[]attr.Value{splunkObjVal},
	)
}

func (r *alertingChannelResourceFramework) mapVictorOpsChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create VictorOps channel object
	victorOpsObj := map[string]attr.Value{
		AlertingChannelVictorOpsFieldAPIKey:     setStringPointerToState(channel.APIKey),
		AlertingChannelVictorOpsFieldRoutingKey: setStringPointerToState(channel.RoutingKey),
	}

	victorOpsObjVal, victorOpsObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelVictorOpsFieldAPIKey:     types.StringType,
			AlertingChannelVictorOpsFieldRoutingKey: types.StringType,
		},
		victorOpsObj,
	)
	diags.Append(victorOpsObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelVictorOpsFieldAPIKey:     types.StringType,
				AlertingChannelVictorOpsFieldRoutingKey: types.StringType,
			},
		},
		[]attr.Value{victorOpsObjVal},
	)
}

func (r *alertingChannelResourceFramework) mapWebhookChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create webhook URLs set
	webhookURLsSet, webhookURLsDiags := types.SetValueFrom(ctx, types.StringType, channel.WebhookURLs)
	if webhookURLsDiags.HasError() {
		diags.Append(webhookURLsDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create HTTP headers map
	headers := r.createHTTPHeaderMapFromList(channel.Headers)
	headersMap, headersDiags := types.MapValueFrom(ctx, types.StringType, headers)
	if headersDiags.HasError() {
		diags.Append(headersDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	// Create Webhook channel object
	webhookObj := map[string]attr.Value{
		AlertingChannelWebhookFieldWebhookURLs: webhookURLsSet,
		AlertingChannelWebhookFieldHTTPHeaders: headersMap,
	}

	webhookObjVal, webhookObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelWebhookFieldWebhookURLs: types.SetType{ElemType: types.StringType},
			AlertingChannelWebhookFieldHTTPHeaders: types.MapType{ElemType: types.StringType},
		},
		webhookObj,
	)
	diags.Append(webhookObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelWebhookFieldWebhookURLs: types.SetType{ElemType: types.StringType},
				AlertingChannelWebhookFieldHTTPHeaders: types.MapType{ElemType: types.StringType},
			},
		},
		[]attr.Value{webhookObjVal},
	)
}

func (r *alertingChannelResourceFramework) createHTTPHeaderMapFromList(headers []string) map[string]string {
	result := make(map[string]string)
	for _, header := range headers {
		keyValue := strings.Split(header, ":")
		if len(keyValue) == 2 {
			result[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
		} else {
			result[strings.TrimSpace(keyValue[0])] = ""
		}
	}
	return result
}

func (r *alertingChannelResourceFramework) mapWebhookBasedChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create webhook-based channel object
	webhookBasedObj := map[string]attr.Value{
		AlertingChannelWebhookBasedFieldWebhookURL: setStringPointerToState(channel.WebhookURL),
	}

	webhookBasedObjVal, webhookBasedObjDiags := types.ObjectValue(
		map[string]attr.Type{
			AlertingChannelWebhookBasedFieldWebhookURL: types.StringType,
		},
		webhookBasedObj,
	)
	diags.Append(webhookBasedObjDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AlertingChannelWebhookBasedFieldWebhookURL: types.StringType,
			},
		},
		[]attr.Value{webhookBasedObjVal},
	)
}

func (r *alertingChannelResourceFramework) convertCommaSeparatedListToSlice(csv string) []string {
	entries := strings.Split(csv, ",")
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = strings.TrimSpace(e)
	}
	return result
}

func (r *alertingChannelResourceFramework) mapEmailChannelFromState(ctx context.Context, id string, name string, emailList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract email object from list
	if len(emailList.Elements()) == 0 {
		diags.AddError("Invalid Email Channel Configuration", "Email channel configuration is empty")
		return nil, diags
	}

	var emailElements []types.Object
	diags.Append(emailList.ElementsAs(ctx, &emailElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract emails set from object
	var emailObj struct {
		Emails types.Set `tfsdk:"emails"`
	}

	diags.Append(emailElements[0].As(ctx, &emailObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert emails set to string slice
	var emails []string
	diags.Append(emailObj.Emails.ElementsAs(ctx, &emails, false)...)
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

func (r *alertingChannelResourceFramework) mapOpsGenieChannelFromState(ctx context.Context, id string, name string, opsGenieList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract OpsGenie object from list
	if len(opsGenieList.Elements()) == 0 {
		diags.AddError("Invalid OpsGenie Channel Configuration", "OpsGenie channel configuration is empty")
		return nil, diags
	}

	var opsGenieElements []types.Object
	diags.Append(opsGenieList.ElementsAs(ctx, &opsGenieElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract fields from object
	var opsGenieObj struct {
		APIKey types.String `tfsdk:"api_key"`
		Region types.String `tfsdk:"region"`
		Tags   types.List   `tfsdk:"tags"`
	}

	diags.Append(opsGenieElements[0].As(ctx, &opsGenieObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert tags list to string slice
	var tags []string
	diags.Append(opsGenieObj.Tags.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Join tags into comma-separated string
	tagsString := strings.Join(tags, ",")

	// Create alerting channel
	apiKeyValue := opsGenieObj.APIKey.ValueString()
	regionValue := opsGenieObj.Region.ValueString()

	return &restapi.AlertingChannel{
		ID:     id,
		Name:   name,
		Kind:   restapi.OpsGenieChannelType,
		APIKey: &apiKeyValue,
		Region: &regionValue,
		Tags:   &tagsString,
	}, nil
}

func (r *alertingChannelResourceFramework) mapPagerDutyChannelFromState(ctx context.Context, id string, name string, pagerDutyList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract PagerDuty object from list
	if len(pagerDutyList.Elements()) == 0 {
		diags.AddError("Invalid PagerDuty Channel Configuration", "PagerDuty channel configuration is empty")
		return nil, diags
	}

	var pagerDutyElements []types.Object
	diags.Append(pagerDutyList.ElementsAs(ctx, &pagerDutyElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract service integration key from object
	var pagerDutyObj struct {
		ServiceIntegrationKey types.String `tfsdk:"service_integration_key"`
	}

	diags.Append(pagerDutyElements[0].As(ctx, &pagerDutyObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	serviceIntegrationKeyValue := pagerDutyObj.ServiceIntegrationKey.ValueString()

	return &restapi.AlertingChannel{
		ID:                    id,
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKeyValue,
	}, nil
}

func (r *alertingChannelResourceFramework) mapSlackChannelFromState(ctx context.Context, id string, name string, slackList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract Slack object from list
	if len(slackList.Elements()) == 0 {
		diags.AddError("Invalid Slack Channel Configuration", "Slack channel configuration is empty")
		return nil, diags
	}

	var slackElements []types.Object
	diags.Append(slackList.ElementsAs(ctx, &slackElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract fields from object
	var slackObj struct {
		WebhookURL types.String `tfsdk:"webhook_url"`
		IconURL    types.String `tfsdk:"icon_url"`
		Channel    types.String `tfsdk:"channel"`
	}

	diags.Append(slackElements[0].As(ctx, &slackObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	webhookURLValue := slackObj.WebhookURL.ValueString()

	result := &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURLValue,
	}

	// Add optional fields if present
	if !slackObj.IconURL.IsNull() {
		iconURLValue := slackObj.IconURL.ValueString()
		result.IconURL = &iconURLValue
	}

	if !slackObj.Channel.IsNull() {
		channelValue := slackObj.Channel.ValueString()
		result.Channel = &channelValue
	}

	return result, nil
}

func (r *alertingChannelResourceFramework) mapSplunkChannelFromState(ctx context.Context, id string, name string, splunkList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract Splunk object from list
	if len(splunkList.Elements()) == 0 {
		diags.AddError("Invalid Splunk Channel Configuration", "Splunk channel configuration is empty")
		return nil, diags
	}

	var splunkElements []types.Object
	diags.Append(splunkList.ElementsAs(ctx, &splunkElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract fields from object
	var splunkObj struct {
		URL   types.String `tfsdk:"url"`
		Token types.String `tfsdk:"token"`
	}

	diags.Append(splunkElements[0].As(ctx, &splunkObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	urlValue := splunkObj.URL.ValueString()
	tokenValue := splunkObj.Token.ValueString()

	return &restapi.AlertingChannel{
		ID:    id,
		Name:  name,
		Kind:  restapi.SplunkChannelType,
		URL:   &urlValue,
		Token: &tokenValue,
	}, nil
}

func (r *alertingChannelResourceFramework) mapVictorOpsChannelFromState(ctx context.Context, id string, name string, victorOpsList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract VictorOps object from list
	if len(victorOpsList.Elements()) == 0 {
		diags.AddError("Invalid VictorOps Channel Configuration", "VictorOps channel configuration is empty")
		return nil, diags
	}

	var victorOpsElements []types.Object
	diags.Append(victorOpsList.ElementsAs(ctx, &victorOpsElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract fields from object
	var victorOpsObj struct {
		APIKey     types.String `tfsdk:"api_key"`
		RoutingKey types.String `tfsdk:"routing_key"`
	}

	diags.Append(victorOpsElements[0].As(ctx, &victorOpsObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	apiKeyValue := victorOpsObj.APIKey.ValueString()
	routingKeyValue := victorOpsObj.RoutingKey.ValueString()

	return &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKeyValue,
		RoutingKey: &routingKeyValue,
	}, nil
}

func (r *alertingChannelResourceFramework) mapWebhookChannelFromState(ctx context.Context, id string, name string, webhookList types.List) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract Webhook object from list
	if len(webhookList.Elements()) == 0 {
		diags.AddError("Invalid Webhook Channel Configuration", "Webhook channel configuration is empty")
		return nil, diags
	}

	var webhookElements []types.Object
	diags.Append(webhookList.ElementsAs(ctx, &webhookElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract fields from object
	var webhookObj struct {
		WebhookURLs types.Set `tfsdk:"webhook_urls"`
		HTTPHeaders types.Map `tfsdk:"http_headers"`
	}

	diags.Append(webhookElements[0].As(ctx, &webhookObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert webhook URLs set to string slice
	var webhookURLs []string
	diags.Append(webhookObj.WebhookURLs.ElementsAs(ctx, &webhookURLs, false)...)
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
	if !webhookObj.HTTPHeaders.IsNull() && !webhookObj.HTTPHeaders.IsUnknown() {
		var httpHeaders map[string]string
		diags.Append(webhookObj.HTTPHeaders.ElementsAs(ctx, &httpHeaders, false)...)
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

func (r *alertingChannelResourceFramework) mapWebhookBasedChannelFromState(ctx context.Context, id string, name string, webhookBasedList types.List, channelType restapi.AlertingChannelType) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract webhook-based object from list
	if len(webhookBasedList.Elements()) == 0 {
		diags.AddError("Invalid Webhook-based Channel Configuration", "Webhook-based channel configuration is empty")
		return nil, diags
	}

	var webhookBasedElements []types.Object
	diags.Append(webhookBasedList.ElementsAs(ctx, &webhookBasedElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Extract webhook URL from object
	var webhookBasedObj struct {
		WebhookURL types.String `tfsdk:"webhook_url"`
	}

	diags.Append(webhookBasedElements[0].As(ctx, &webhookBasedObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Create alerting channel
	webhookURLValue := webhookBasedObj.WebhookURL.ValueString()

	return &restapi.AlertingChannel{
		ID:         id,
		Name:       name,
		Kind:       channelType,
		WebhookURL: &webhookURLValue,
	}, nil
}

func (r *alertingChannelResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AlertingChannel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AlertingChannelModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
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
	if !model.Email.IsNull() && !model.Email.IsUnknown() && len(model.Email.Elements()) > 0 {
		return r.mapEmailChannelFromState(ctx, id, name, model.Email)
	}

	if !model.OpsGenie.IsNull() && !model.OpsGenie.IsUnknown() && len(model.OpsGenie.Elements()) > 0 {
		return r.mapOpsGenieChannelFromState(ctx, id, name, model.OpsGenie)
	}

	if !model.PagerDuty.IsNull() && !model.PagerDuty.IsUnknown() && len(model.PagerDuty.Elements()) > 0 {
		return r.mapPagerDutyChannelFromState(ctx, id, name, model.PagerDuty)
	}

	if !model.Slack.IsNull() && !model.Slack.IsUnknown() && len(model.Slack.Elements()) > 0 {
		return r.mapSlackChannelFromState(ctx, id, name, model.Slack)
	}

	if !model.Splunk.IsNull() && !model.Splunk.IsUnknown() && len(model.Splunk.Elements()) > 0 {
		return r.mapSplunkChannelFromState(ctx, id, name, model.Splunk)
	}

	if !model.VictorOps.IsNull() && !model.VictorOps.IsUnknown() && len(model.VictorOps.Elements()) > 0 {
		return r.mapVictorOpsChannelFromState(ctx, id, name, model.VictorOps)
	}

	if !model.Webhook.IsNull() && !model.Webhook.IsUnknown() && len(model.Webhook.Elements()) > 0 {
		return r.mapWebhookChannelFromState(ctx, id, name, model.Webhook)
	}

	if !model.Office365.IsNull() && !model.Office365.IsUnknown() && len(model.Office365.Elements()) > 0 {
		return r.mapWebhookBasedChannelFromState(ctx, id, name, model.Office365, restapi.Office365ChannelType)
	}

	if !model.GoogleChat.IsNull() && !model.GoogleChat.IsUnknown() && len(model.GoogleChat.Elements()) > 0 {
		return r.mapWebhookBasedChannelFromState(ctx, id, name, model.GoogleChat, restapi.GoogleChatChannelType)
	}

	diags.AddError(
		"Invalid Alerting Channel Configuration",
		"No valid alerting channel configuration found. Please configure exactly one channel type.",
	)
	return nil, diags
}

// Made with Bob
