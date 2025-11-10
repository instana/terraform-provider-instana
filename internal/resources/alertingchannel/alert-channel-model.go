package alertingchannel

import (
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AlertingChannelModel represents the data model for the alerting channel resource
type AlertingChannelModel struct {
	ID                    types.String                       `tfsdk:"id"`
	Name                  types.String                       `tfsdk:"name"`
	Email                 *shared.EmailModel                 `tfsdk:"email"`
	OpsGenie              *shared.OpsGenieModel              `tfsdk:"ops_genie"`
	PagerDuty             *shared.PagerDutyModel             `tfsdk:"pager_duty"`
	Slack                 *shared.SlackModel                 `tfsdk:"slack"`
	Splunk                *shared.SplunkModel                `tfsdk:"splunk"`
	VictorOps             *shared.VictorOpsModel             `tfsdk:"victor_ops"`
	Webhook               *shared.WebhookModel               `tfsdk:"webhook"`
	Office365             *shared.WebhookBasedModel          `tfsdk:"office_365"`
	GoogleChat            *shared.WebhookBasedModel          `tfsdk:"google_chat"`
	ServiceNow            *shared.ServiceNowModel            `tfsdk:"service_now"`
	ServiceNowApplication *shared.ServiceNowApplicationModel `tfsdk:"service_now_application"`
	PrometheusWebhook     *shared.PrometheusWebhookModel     `tfsdk:"prometheus_webhook"`
	WebexTeamsWebhook     *shared.WebhookBasedModel          `tfsdk:"webex_teams_webhook"`
	WatsonAIOpsWebhook    *shared.WatsonAIOpsWebhookModel    `tfsdk:"watson_aiops_webhook"`
}
