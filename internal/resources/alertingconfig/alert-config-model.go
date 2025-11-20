package alertingconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

// AlertingConfigModel represents the data model for the alerting configuration resource
type AlertingConfigModel struct {
	ID                             types.String `tfsdk:"id"`
	AlertName                      types.String `tfsdk:"alert_name"`
	IntegrationIDs                 types.Set    `tfsdk:"integration_ids"`
	EventFilterQuery               types.String `tfsdk:"event_filter_query"`
	EventFilterEventTypes          types.Set    `tfsdk:"event_filter_event_types"`
	EventFilterRuleIDs             types.Set    `tfsdk:"event_filter_rule_ids"`
	EventFilterApplicationAlertIDs types.Set    `tfsdk:"event_filter_application_alert_ids"`
	CustomPayloadFields            types.List   `tfsdk:"custom_payload_field"`
}
