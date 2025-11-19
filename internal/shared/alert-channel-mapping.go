package shared

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type PagerDutyModel struct {
	ServiceIntegrationKey types.String `tfsdk:"service_integration_key"`
}

type SlackModel struct {
	WebhookURL types.String `tfsdk:"webhook_url"`
	IconURL    types.String `tfsdk:"icon_url"`
	Channel    types.String `tfsdk:"channel"`
}

type SplunkModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

type VictorOpsModel struct {
	APIKey     types.String `tfsdk:"api_key"`
	RoutingKey types.String `tfsdk:"routing_key"`
}

type WebhookModel struct {
	WebhookURLs types.Set `tfsdk:"webhook_urls"`
	HTTPHeaders types.Map `tfsdk:"http_headers"`
}

type EmailModel struct {
	Emails types.Set `tfsdk:"emails"`
}

type OpsGenieModel struct {
	APIKey types.String `tfsdk:"api_key"`
	Region types.String `tfsdk:"region"`
	Tags   types.List   `tfsdk:"tags"`
}
type WatsonAIOpsWebhookModel struct {
	WebhookURL  types.String `tfsdk:"webhook_url"`
	HTTPHeaders types.List   `tfsdk:"http_headers"`
}

type PrometheusWebhookModel struct {
	WebhookURL types.String `tfsdk:"webhook_url"`
	Receiver   types.String `tfsdk:"receiver"`
}

type WebhookBasedModel struct {
	WebhookURL types.String `tfsdk:"webhook_url"`
}

const ResourceFieldThresholdRuleWarningSeverity = "warning"
const ResourceFieldThresholdRuleCriticalSeverity = "critical"

type ServiceNowApplicationModel struct {
	ServiceNowURL                  types.String `tfsdk:"service_now_url"`
	Username                       types.String `tfsdk:"username"`
	Password                       types.String `tfsdk:"password"`
	Tenant                         types.String `tfsdk:"tenant"`
	Unit                           types.String `tfsdk:"unit"`
	AutoCloseIncidents             types.Bool   `tfsdk:"auto_close_incidents"`
	InstanaURL                     types.String `tfsdk:"instana_url"`
	EnableSendInstanaNotes         types.Bool   `tfsdk:"enable_send_instana_notes"`
	EnableSendServiceNowActivities types.Bool   `tfsdk:"enable_send_service_now_activities"`
	EnableSendServiceNowWorkNotes  types.Bool   `tfsdk:"enable_send_service_now_work_notes"`
	ManuallyClosedIncidents        types.Bool   `tfsdk:"manually_closed_incidents"`
	ResolutionOfIncident           types.Bool   `tfsdk:"resolution_of_incident"`
	SnowStatusOnCloseEvent         types.Int64  `tfsdk:"snow_status_on_close_event"`
}

type ServiceNowModel struct {
	ServiceNowURL      types.String `tfsdk:"service_now_url"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	AutoCloseIncidents types.Bool   `tfsdk:"auto_close_incidents"`
}

type SlackAppModel struct {
	AppID          types.String `tfsdk:"app_id"`
	TeamID         types.String `tfsdk:"team_id"`
	TeamName       types.String `tfsdk:"team_name"`
	ChannelID      types.String `tfsdk:"channel_id"`
	ChannelName    types.String `tfsdk:"channel_name"`
	EmojiRendering types.Bool   `tfsdk:"emoji_rendering"`
}

type MsTeamsAppModel struct {
	APITokenID types.String `tfsdk:"api_token_id"`
	TeamID     types.String `tfsdk:"team_id"`
	TeamName   types.String `tfsdk:"team_name"`
	ChannelID  types.String `tfsdk:"channel_id"`
	ChannelName types.String `tfsdk:"channel_name"`
	InstanaURL types.String `tfsdk:"instana_url"`
	ServiceURL types.String `tfsdk:"service_url"`
	TenantID   types.String `tfsdk:"tenant_id"`
	TenantName types.String `tfsdk:"tenant_name"`
}

func MapAlertChannelsToState(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		ResourceFieldThresholdRuleWarningSeverity:  types.ListType{ElemType: types.StringType},
		ResourceFieldThresholdRuleCriticalSeverity: types.ListType{ElemType: types.StringType},
	}

	if len(alertChannels) == 0 {
		return types.ObjectNull(attrTypes), diags
	}

	alertChannelsObj := map[string]attr.Value{}

	// Map warning severity
	if warningChannels, ok := alertChannels[restapi.WarningSeverity]; ok && len(warningChannels) > 0 {
		warningList, warningDiags := types.ListValueFrom(ctx, types.StringType, warningChannels)
		diags.Append(warningDiags...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}
		alertChannelsObj[ResourceFieldThresholdRuleWarningSeverity] = warningList
	} else {
		alertChannelsObj[ResourceFieldThresholdRuleWarningSeverity] = types.ListNull(types.StringType)
	}

	// Map critical severity
	if criticalChannels, ok := alertChannels[restapi.CriticalSeverity]; ok && len(criticalChannels) > 0 {
		criticalList, criticalDiags := types.ListValueFrom(ctx, types.StringType, criticalChannels)
		diags.Append(criticalDiags...)
		if diags.HasError() {
			return types.ObjectNull(attrTypes), diags
		}
		alertChannelsObj[ResourceFieldThresholdRuleCriticalSeverity] = criticalList
	} else {
		alertChannelsObj[ResourceFieldThresholdRuleCriticalSeverity] = types.ListNull(types.StringType)
	}

	return types.ObjectValue(attrTypes, alertChannelsObj)
}

func MapAlertChannelsFromState(ctx context.Context, alertChannelsObj types.Object) (map[restapi.AlertSeverity][]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannelsMap := make(map[restapi.AlertSeverity][]string)

	if alertChannelsObj.IsNull() || alertChannelsObj.IsUnknown() {
		return alertChannelsMap, diags
	}

	var alertChannels struct {
		Warning  types.List `tfsdk:"warning"`
		Critical types.List `tfsdk:"critical"`
	}

	diags.Append(alertChannelsObj.As(ctx, &alertChannels, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}
	warningChannels := make([]string, 0)
	// Map warning severity
	if !alertChannels.Warning.IsNull() && !alertChannels.Warning.IsUnknown() {

		diags.Append(alertChannels.Warning.ElementsAs(ctx, &warningChannels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		//if len(warningChannels) > 0 {

		//}
	}
	alertChannelsMap[restapi.WarningSeverity] = warningChannels

	criticalChannels := make([]string, 0)
	// Map critical severity
	if !alertChannels.Critical.IsNull() && !alertChannels.Critical.IsUnknown() {

		diags.Append(alertChannels.Critical.ElementsAs(ctx, &criticalChannels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		//if len(criticalChannels) > 0 {

		//}
	}
	alertChannelsMap[restapi.CriticalSeverity] = criticalChannels

	return alertChannelsMap, diags
}

func MapWebhookBasedChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*WebhookBasedModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create and return webhook-based model
	return &WebhookBasedModel{
		WebhookURL: util.SetStringPointerToState(channel.WebhookURL),
	}, diags
}

func MapServiceNowChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*ServiceNowModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create ServiceNow model
	model := &ServiceNowModel{
		ServiceNowURL: util.SetStringPointerToState(channel.ServiceNowURL),
		Username:      util.SetStringPointerToState(channel.Username),
	}

	// Use password from API response if available
	if channel.Password != nil && *channel.Password != "" {
		model.Password = types.StringValue(*channel.Password)
	} else {
		model.Password = types.StringNull()
	}

	// Add optional autoCloseIncidents field
	if channel.AutoCloseIncidents != nil {
		model.AutoCloseIncidents = types.BoolValue(*channel.AutoCloseIncidents)
	} else {
		model.AutoCloseIncidents = types.BoolNull()
	}

	return model, diags
}

func MapServiceNowApplicationChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*ServiceNowApplicationModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create ServiceNow Enhanced model with required fields
	model := &ServiceNowApplicationModel{
		ServiceNowURL: util.SetStringPointerToState(channel.ServiceNowURL),
		Username:      util.SetStringPointerToState(channel.Username),
		Tenant:        util.SetStringPointerToState(channel.Tenant),
		Unit:          util.SetStringPointerToState(channel.Unit),
	}

	// Use password from API response if available
	if channel.Password != nil && *channel.Password != "" {
		model.Password = types.StringValue(*channel.Password)
	} else {
		model.Password = types.StringNull()
	}

	// Add optional boolean fields
	if channel.AutoCloseIncidents != nil {
		model.AutoCloseIncidents = types.BoolValue(*channel.AutoCloseIncidents)
	} else {
		model.AutoCloseIncidents = types.BoolNull()
	}

	if channel.EnableSendInstanaNotes != nil {
		model.EnableSendInstanaNotes = types.BoolValue(*channel.EnableSendInstanaNotes)
	} else {
		model.EnableSendInstanaNotes = types.BoolNull()
	}

	if channel.EnableSendServiceNowActivities != nil {
		model.EnableSendServiceNowActivities = types.BoolValue(*channel.EnableSendServiceNowActivities)
	} else {
		model.EnableSendServiceNowActivities = types.BoolNull()
	}

	if channel.EnableSendServiceNowWorkNotes != nil {
		model.EnableSendServiceNowWorkNotes = types.BoolValue(*channel.EnableSendServiceNowWorkNotes)
	} else {
		model.EnableSendServiceNowWorkNotes = types.BoolNull()
	}

	if channel.ManuallyClosedIncidents != nil {
		model.ManuallyClosedIncidents = types.BoolValue(*channel.ManuallyClosedIncidents)
	} else {
		model.ManuallyClosedIncidents = types.BoolNull()
	}

	if channel.ResolutionOfIncident != nil {
		model.ResolutionOfIncident = types.BoolValue(*channel.ResolutionOfIncident)
	} else {
		model.ResolutionOfIncident = types.BoolNull()
	}

	// Add optional string field
	if channel.InstanaURL != nil && *channel.InstanaURL != "" {
		model.InstanaURL = util.SetStringPointerToState(channel.InstanaURL)
	} else {
		model.InstanaURL = types.StringNull()
	}

	// Add optional int field
	if channel.SnowStatusOnCloseEvent != nil {
		model.SnowStatusOnCloseEvent = types.Int64Value(int64(*channel.SnowStatusOnCloseEvent))
	} else {
		model.SnowStatusOnCloseEvent = types.Int64Null()
	}

	return model, diags
}

func MapPrometheusWebhookChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*PrometheusWebhookModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create Prometheus Webhook model
	model := &PrometheusWebhookModel{
		WebhookURL: util.SetStringPointerToState(channel.WebhookURL),
	}

	// Add optional receiver field
	if channel.Receiver != nil && *channel.Receiver != "" {
		model.Receiver = util.SetStringPointerToState(channel.Receiver)
	} else {
		model.Receiver = types.StringNull()
	}

	return model, diags
}

func MapWatsonAIOpsWebhookChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*WatsonAIOpsWebhookModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create Watson AIOps Webhook model
	model := &WatsonAIOpsWebhookModel{
		WebhookURL: util.SetStringPointerToState(channel.WebhookURL),
	}

	// Add optional headers field
	if channel.Headers != nil && len(channel.Headers) > 0 {
		headersList, headersDiags := types.ListValueFrom(ctx, types.StringType, channel.Headers)
		if headersDiags.HasError() {
			diags.Append(headersDiags...)
			return nil, diags
		}
		model.HTTPHeaders = headersList
	} else {
		model.HTTPHeaders = types.ListNull(types.StringType)
	}

	return model, diags
}
func MapEmailChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*EmailModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create email set
	emailsSet, emailsDiags := types.SetValueFrom(ctx, types.StringType, channel.Emails)
	if emailsDiags.HasError() {
		diags.Append(emailsDiags...)
		return nil, diags
	}

	// Create and return email model
	return &EmailModel{
		Emails: emailsSet,
	}, diags
}

func MapOpsGenieChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*OpsGenieModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert comma-separated tags to slice
	tags := ConvertCommaSeparatedListToSlice(*channel.Tags)

	// Create tags list
	tagsList, tagsDiags := types.ListValueFrom(ctx, types.StringType, tags)
	if tagsDiags.HasError() {
		diags.Append(tagsDiags...)
		return nil, diags
	}

	// Create and return OpsGenie model
	return &OpsGenieModel{
		APIKey: util.SetStringPointerToState(channel.APIKey),
		Region: util.SetStringPointerToState(channel.Region),
		Tags:   tagsList,
	}, diags
}

func MapPagerDutyChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*PagerDutyModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create and return PagerDuty model
	return &PagerDutyModel{
		ServiceIntegrationKey: util.SetStringPointerToState(channel.ServiceIntegrationKey),
	}, diags
}

func MapSlackChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*SlackModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create Slack model
	model := &SlackModel{
		WebhookURL: util.SetStringPointerToState(channel.WebhookURL),
	}

	// Add optional fields if present
	if channel.IconURL != nil && *channel.IconURL != "" {
		model.IconURL = util.SetStringPointerToState(channel.IconURL)
	} else {
		model.IconURL = types.StringNull()
	}

	if channel.Channel != nil && *channel.Channel != "" {
		model.Channel = util.SetStringPointerToState(channel.Channel)
	} else {
		model.Channel = types.StringNull()
	}

	return model, diags
}

func MapSplunkChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*SplunkModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create and return Splunk model
	return &SplunkModel{
		URL:   util.SetStringPointerToState(channel.URL),
		Token: util.SetStringPointerToState(channel.Token),
	}, diags
}

func MapVictorOpsChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*VictorOpsModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create and return VictorOps model
	return &VictorOpsModel{
		APIKey:     util.SetStringPointerToState(channel.APIKey),
		RoutingKey: util.SetStringPointerToState(channel.RoutingKey),
	}, diags
}

func MapWebhookChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*WebhookModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create webhook URLs set
	webhookURLsSet, webhookURLsDiags := types.SetValueFrom(ctx, types.StringType, channel.WebhookURLs)
	if webhookURLsDiags.HasError() {
		diags.Append(webhookURLsDiags...)
		return nil, diags
	}

	// Create HTTP headers map
	headers := CreateHTTPHeaderMapFromList(channel.Headers)
	headersMap, headersDiags := types.MapValueFrom(ctx, types.StringType, headers)
	if headersDiags.HasError() {
		diags.Append(headersDiags...)
		return nil, diags
	}

	// Create and return Webhook model
	return &WebhookModel{
		WebhookURLs: webhookURLsSet,
		HTTPHeaders: headersMap,
	}, diags
}

func ConvertCommaSeparatedListToSlice(csv string) []string {
	entries := strings.Split(csv, ",")
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = strings.TrimSpace(e)
	}
	return result
}

func CreateHTTPHeaderMapFromList(headers []string) map[string]string {
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

func MapSlackAppChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*SlackAppModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create Slack App model
	model := &SlackAppModel{
		AppID:       util.SetStringPointerToState(channel.AppID),
		TeamID:      util.SetStringPointerToState(channel.TeamID),
		TeamName:    util.SetStringPointerToState(channel.TeamName),
		ChannelID:   util.SetStringPointerToState(channel.ChannelID),
		ChannelName: util.SetStringPointerToState(channel.ChannelName),
	}

	// Add optional emoji rendering field
	if channel.EmojiRendering != nil {
		model.EmojiRendering = types.BoolValue(*channel.EmojiRendering)
	} else {
		model.EmojiRendering = types.BoolNull()
	}

	return model, diags
}

func MapMsTeamsAppChannelToState(ctx context.Context, channel *restapi.AlertingChannel) (*MsTeamsAppModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create MS Teams App model
	model := &MsTeamsAppModel{
		APITokenID: util.SetStringPointerToState(channel.APITokenID),
		TeamID:     util.SetStringPointerToState(channel.TeamID),
		TeamName:   util.SetStringPointerToState(channel.TeamName),
		ChannelID:  util.SetStringPointerToState(channel.ChannelID),
		ChannelName: util.SetStringPointerToState(channel.ChannelName),
		InstanaURL: util.SetStringPointerToState(channel.InstanaURL),
		ServiceURL: util.SetStringPointerToState(channel.ServiceURL),
		TenantID:   util.SetStringPointerToState(channel.TenantID),
		TenantName: util.SetStringPointerToState(channel.TenantName),
	}

	return model, diags
}
