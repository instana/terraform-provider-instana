package automationpolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/shared"
)

// AutomationPolicyModel represents the data model for the automation policy resource
type AutomationPolicyModel struct {
	ID                types.String             `tfsdk:"id"`
	Name              types.String             `tfsdk:"name"`
	Description       types.String             `tfsdk:"description"`
	Tags              types.Set                `tfsdk:"tags"`
	Trigger           *TriggerModel            `tfsdk:"trigger"`
	TypeConfiguration []TypeConfigurationModel `tfsdk:"type_configuration"`
}

// TriggerModel represents a trigger in the automation policy
type TriggerModel struct {
	ID          types.String     `tfsdk:"id"`
	Type        types.String     `tfsdk:"type"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
	Scheduling  *SchedulingModel `tfsdk:"scheduling"`
}

// SchedulingModel represents the scheduling configuration for automation policy trigger
type SchedulingModel struct {
	StartTime     types.Int64  `tfsdk:"start_time"`
	Duration      types.Int64  `tfsdk:"duration"`
	DurationUnit  types.String `tfsdk:"duration_unit"`
	RecurrentRule types.String `tfsdk:"recurrent_rule"`
	Recurrent     types.Bool   `tfsdk:"recurrent"`
}

// TypeConfigurationModel represents a type configuration in the automation policy
type TypeConfigurationModel struct {
	Name      types.String        `tfsdk:"name"`
	Condition *ConditionModel     `tfsdk:"condition"`
	Action    []PolicyActionModel `tfsdk:"action"`
}

// ConditionModel represents a condition in the automation policy
type ConditionModel struct {
	Query types.String `tfsdk:"query"`
}

// PolicyActionModel represents an action reference in the automation policy
// This is different from AutomationActionModel - it only contains the reference and parameters
type PolicyActionModel struct {
	Action  shared.AutomationActionModel `tfsdk:"action"`
	AgentID types.String                 `tfsdk:"agent_id"`
}
