package automationpolicy

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewAutomationPolicyResourceHandleFramework creates the resource handle for Automation Policies
func NewAutomationPolicyResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.AutomationPolicy] {
	return &automationPolicyResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAutomationPolicyFramework,
			Schema: schema.Schema{
				Description: AutomationPolicyDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: AutomationPolicyDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					shared.AutomationPolicyFieldName: schema.StringAttribute{
						Required:    true,
						Description: AutomationPolicyDescName,
					},
					shared.AutomationPolicyFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: AutomationPolicyDescDescription,
					},
					shared.AutomationPolicyFieldTags: schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: AutomationPolicyDescTags,
					},
					shared.AutomationPolicyFieldTrigger: schema.SingleNestedAttribute{
						Required:    true,
						Description: AutomationPolicyDescTrigger,
						Attributes: map[string]schema.Attribute{
							shared.AutomationPolicyFieldId: schema.StringAttribute{
								Required:    true,
								Description: AutomationPolicyDescTriggerID,
							},
							shared.AutomationPolicyFieldType: schema.StringAttribute{
								Required:    true,
								Description: AutomationPolicyDescTriggerType,
								Validators: []validator.String{
									stringvalidator.OneOf(shared.SupportedTriggerTypes...),
								},
							},
							shared.AutomationPolicyFieldName: schema.StringAttribute{
								Optional:    true,
								Description: AutomationPolicyDescTriggerName,
							},
							shared.AutomationPolicyFieldDescription: schema.StringAttribute{
								Optional:    true,
								Description: AutomationPolicyDescTriggerDescription,
							},
							"scheduling": schema.SingleNestedAttribute{
								Optional:    true,
								Description: AutomationPolicyDescTriggerScheduling,
								Attributes: map[string]schema.Attribute{
									"start_time": schema.Int64Attribute{
										Optional:    true,
										Description: AutomationPolicyDescTriggerSchedulingStartTime,
									},
									"duration": schema.Int64Attribute{
										Optional:    true,
										Description: AutomationPolicyDescTriggerSchedulingDuration,
									},
									"duration_unit": schema.StringAttribute{
										Optional:    true,
										Description: AutomationPolicyDescTriggerSchedulingDurationUnit,
										Validators: []validator.String{
											stringvalidator.OneOf("MINUTE", "HOUR", "DAY"),
										},
									},
									"recurrent_rule": schema.StringAttribute{
										Optional:    true,
										Description: AutomationPolicyDescTriggerSchedulingRecurrentRule,
									},
									"recurrent": schema.BoolAttribute{
										Optional:    true,
										Description: AutomationPolicyDescTriggerSchedulingRecurrent,
									},
								},
							},
						},
					},
					shared.AutomationPolicyFieldTypeConfiguration: schema.ListNestedAttribute{
						Required:    true,
						Description: AutomationPolicyDescTypeConfiguration,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								shared.AutomationPolicyFieldName: schema.StringAttribute{
									Required:    true,
									Description: AutomationPolicyDescTypeConfigurationName,
									Validators: []validator.String{
										stringvalidator.OneOf(shared.SupportedPolicyTypes...),
									},
								},
								shared.AutomationPolicyFieldCondition: schema.SingleNestedAttribute{
									Optional:    true,
									Description: AutomationPolicyDescCondition,
									Attributes: map[string]schema.Attribute{
										shared.AutomationPolicyFieldQuery: schema.StringAttribute{
											Required:    true,
											Description: AutomationPolicyDescConditionQuery,
										},
									},
								},
								shared.AutomationPolicyFieldAction: schema.ListNestedAttribute{
									Required:    true,
									Description: AutomationPolicyDescAction,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.SingleNestedAttribute{
												Required:    true,
												Description: AutomationPolicyDescActionAction,
												Attributes:  shared.GetAutomationActionSchemaAttributes(),
											},
											shared.AutomationPolicyFieldAgentId: schema.StringAttribute{
												Optional:    true,
												Description: AutomationPolicyDescActionAgentID,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			SchemaVersion: 0,
		},
	}
}

type automationPolicyResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *automationPolicyResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *automationPolicyResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AutomationPolicy] {
	return api.AutomationPolicies()
}

func (r *automationPolicyResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *automationPolicyResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, policy *restapi.AutomationPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	var model AutomationPolicyModel

	// Read from plan to preserve user-configured values (especially for optional fields)
	// This is important for fields that might not be returned by the API
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
		if diags.HasError() {
			return diags
		}
	}

	// Update model with values from API response
	model.ID = types.StringValue(policy.ID)
	model.Name = types.StringValue(policy.Name)
	model.Description = types.StringValue(policy.Description)
	// Handle tags
	if policy.Tags != nil {
		tags, d := r.mapTagsToState(ctx, policy.Tags)
		diags.Append(d...)
		if !diags.HasError() {
			model.Tags = tags
		}
	} else {
		model.Tags = []string{}
	}

	// Map trigger
	model.Trigger = r.mapTriggerToState(&policy.Trigger, model.Trigger)

	// Map type configurations
	model.TypeConfiguration = r.mapTypeConfigurationsToState(ctx, policy.TypeConfigurations)

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *automationPolicyResourceFramework) mapTagsToState(ctx context.Context, tags interface{}) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tags == nil {
		return []string{}, diags
	}

	// Handle tags based on their type
	switch v := tags.(type) {
	case []interface{}:
		result := make([]string, len(v))
		for i, tag := range v {
			if strTag, ok := tag.(string); ok {
				result[i] = strTag
			} else {
				diags.AddError(
					AutomationPolicyErrMappingTags,
					fmt.Sprintf(AutomationPolicyErrTagNotString, i),
				)
				return []string{}, diags
			}
		}
		return result, diags
	default:
		diags.AddError(
			AutomationPolicyErrMappingTags,
			AutomationPolicyErrTagsFormat,
		)
		return []string{}, diags
	}
}

func (r *automationPolicyResourceFramework) mapTriggerToState(trigger *restapi.Trigger, triggerModel TriggerModel) TriggerModel {
	if triggerModel.ID.IsNull() || triggerModel.ID.IsUnknown() {
		triggerModel.ID = types.StringValue(trigger.Id)
	}
	if triggerModel.Type.IsNull() || triggerModel.Type.IsUnknown() {
		triggerModel.Type = types.StringValue(trigger.Type)
	}

	if triggerModel.Description.IsNull() || triggerModel.Description.IsUnknown() {
		if trigger.Description != "" {
			triggerModel.Description = types.StringValue(trigger.Description)
		} else {
			triggerModel.Description = types.StringNull()
		}

	}
	if triggerModel.Name.IsNull() || triggerModel.Name.IsUnknown() {
		if trigger.Name != "" {
			triggerModel.Name = types.StringValue(trigger.Name)
		} else {
			triggerModel.Name = types.StringNull()
		}

	}

	// Map scheduling from API response if not already set in the model
	// The scheduling field is preserved from the plan in UpdateState function
	if triggerModel.Scheduling == nil {
		if trigger.Scheduling.StartTime != 0 {
			// Handle duration_unit - set to null if empty
			var durationUnit types.String
			if trigger.Scheduling.DurationUnit != "" {
				durationUnit = types.StringValue(string(trigger.Scheduling.DurationUnit))
			} else {
				durationUnit = types.StringNull()
			}

			// Handle recurrent_rule - set to null if empty
			var recurrentRule types.String
			if trigger.Scheduling.RecurrentRule != "" {
				recurrentRule = types.StringValue(trigger.Scheduling.RecurrentRule)
			} else {
				recurrentRule = types.StringNull()
			}

			triggerModel.Scheduling = &SchedulingModel{
				StartTime:     types.Int64Value(trigger.Scheduling.StartTime),
				Duration:      types.Int64Value(int64(trigger.Scheduling.Duration)),
				DurationUnit:  durationUnit,
				RecurrentRule: recurrentRule,
				Recurrent:     types.BoolValue(trigger.Scheduling.Recurrent),
			}
		}
		// If no scheduling data from API, leave it as nil
	}

	return triggerModel
}

func (r *automationPolicyResourceFramework) mapTypeConfigurationsToState(ctx context.Context, typeConfigs []restapi.TypeConfiguration) []TypeConfigurationModel {
	result := make([]TypeConfigurationModel, len(typeConfigs))

	for i, typeConfig := range typeConfigs {
		// Map condition
		var condition *ConditionModel
		if typeConfig.Condition != nil && typeConfig.Condition.Query != "" {
			condition = &ConditionModel{
				Query: types.StringValue(typeConfig.Condition.Query),
			}
		}

		// Map actions
		actions := r.mapActionsToState(ctx, &typeConfig.Runnable)

		result[i] = TypeConfigurationModel{
			Name:      types.StringValue(typeConfig.Name),
			Condition: condition,
			Action:    actions,
		}
	}

	return result
}

func (r *automationPolicyResourceFramework) mapActionsToState(ctx context.Context, runnable *restapi.Runnable) []PolicyActionModel {
	result := make([]PolicyActionModel, len(runnable.RunConfiguration.Actions))

	for i, actionPolicy := range runnable.RunConfiguration.Actions {
		// Map the full automation action from the nested Action field
		tags, _ := shared.MapTagsToState(ctx, actionPolicy.Action.Tags)
		inputParams := shared.MapInputParametersToState(ctx, actionPolicy.Action.InputParameters)

		actionModel := shared.AutomationActionModel{
			ID:             types.StringValue(actionPolicy.Action.ID),
			Name:           types.StringValue(actionPolicy.Action.Name),
			Description:    types.StringValue(actionPolicy.Action.Description),
			Tags:           tags,
			InputParameter: inputParams,
		}

		// Map action type-specific fields using the common function
		shared.MapActionTypeFieldsToState(ctx, &actionPolicy.Action, &actionModel)

		agentID := types.StringNull()
		if actionPolicy.AgentId != "" {
			agentID = types.StringValue(actionPolicy.AgentId)
		}

		result[i] = PolicyActionModel{
			Action:  actionModel,
			AgentID: agentID,
		}
	}

	return result
}

func (r *automationPolicyResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AutomationPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AutomationPolicyModel

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

	// Map trigger
	trigger, d := r.mapTriggerFromState(ctx, model.Trigger)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Map type configurations
	typeConfigurations, d := r.mapTypeConfigurationsFromState(ctx, model.TypeConfiguration)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Map tags
	tags, d := r.mapTagsFromState(ctx, model.Tags)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &restapi.AutomationPolicy{
		ID:                 id,
		Name:               model.Name.ValueString(),
		Description:        model.Description.ValueString(),
		Tags:               tags,
		Trigger:            trigger,
		TypeConfigurations: typeConfigurations,
	}, diags
}

func (r *automationPolicyResourceFramework) mapTriggerFromState(ctx context.Context, triggerModel TriggerModel) (restapi.Trigger, diag.Diagnostics) {
	var diags diag.Diagnostics

	trigger := restapi.Trigger{
		Id:   triggerModel.ID.ValueString(),
		Type: triggerModel.Type.ValueString(),
	}

	// Map optional fields
	if !triggerModel.Name.IsNull() {
		trigger.Name = triggerModel.Name.ValueString()
	}

	if !triggerModel.Description.IsNull() {
		trigger.Description = triggerModel.Description.ValueString()
	}

	// Map scheduling if present
	if triggerModel.Scheduling != nil {
		schedulingModel := triggerModel.Scheduling

		// Only set scheduling if at least start_time is provided
		if !schedulingModel.StartTime.IsNull() && !schedulingModel.StartTime.IsUnknown() {
			trigger.Scheduling = restapi.Scheduling{
				StartTime:    schedulingModel.StartTime.ValueInt64(),
				Duration:     int(schedulingModel.Duration.ValueInt64()),
				DurationUnit: restapi.DurationUnit(schedulingModel.DurationUnit.ValueString()),
				Recurrent:    schedulingModel.Recurrent.ValueBool(),
			}

			if !schedulingModel.RecurrentRule.IsNull() {
				trigger.Scheduling.RecurrentRule = schedulingModel.RecurrentRule.ValueString()
			}
		}
	}

	return trigger, diags
}

func (r *automationPolicyResourceFramework) mapTypeConfigurationsFromState(ctx context.Context, typeConfigModels []TypeConfigurationModel) ([]restapi.TypeConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics

	typeConfigurations := make([]restapi.TypeConfiguration, len(typeConfigModels))
	for i, typeConfigModel := range typeConfigModels {
		// Map condition
		condition, d := r.mapConditionFromState(ctx, typeConfigModel.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		// Map runnable (actions)
		runnable, d := r.mapRunnableFromState(ctx, typeConfigModel.Action)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		typeConfigurations[i] = restapi.TypeConfiguration{
			Name:      typeConfigModel.Name.ValueString(),
			Condition: condition,
			Runnable:  runnable,
		}
	}

	return typeConfigurations, diags
}

func (r *automationPolicyResourceFramework) mapConditionFromState(ctx context.Context, conditionModel *ConditionModel) (*restapi.Condition, diag.Diagnostics) {
	var diags diag.Diagnostics

	if conditionModel == nil {
		return &restapi.Condition{Query: ""}, diags
	}

	condition := &restapi.Condition{
		Query: conditionModel.Query.ValueString(),
	}

	return condition, diags
}

func (r *automationPolicyResourceFramework) mapRunnableFromState(ctx context.Context, actionModels []PolicyActionModel) (restapi.Runnable, diag.Diagnostics) {
	var diags diag.Diagnostics
	runnable := restapi.Runnable{
		Type: actionRunnable,
		RunConfiguration: restapi.RunConfiguration{
			Actions: []restapi.AutomationActionPolicy{},
		},
	}

	if len(actionModels) == 0 {
		return runnable, diags
	}

	actions := make([]restapi.AutomationActionPolicy, len(actionModels))
	for i, policyActionModel := range actionModels {
		// Map the automation action from the model
		actionModel := policyActionModel.Action

		// Map input parameters from the action model
		inputParams, d := shared.MapInputParametersFromState(ctx, actionModel)
		diags.Append(d...)
		if diags.HasError() {
			return runnable, diags
		}

		// Map action type and fields
		actionType, fields, d := shared.MapActionTypeAndFieldsFromState(ctx, actionModel)
		diags.Append(d...)
		if diags.HasError() {
			return runnable, diags
		}

		// Create the automation action
		automationAction := restapi.AutomationAction{
			ID:              actionModel.ID.ValueString(),
			Name:            actionModel.Name.ValueString(),
			Description:     actionModel.Description.ValueString(),
			Type:            actionType,
			Fields:          fields,
			InputParameters: inputParams,
		}

		// Map tags
		if !actionModel.Tags.IsNull() {
			tags, d := shared.MapTagsFromState(ctx, actionModel.Tags)
			diags.Append(d...)
			if !diags.HasError() {
				automationAction.Tags = tags
			}
		}

		// Create the action policy
		agentId := ""
		if !policyActionModel.AgentID.IsNull() {
			agentId = policyActionModel.AgentID.ValueString()
		}

		actions[i] = restapi.AutomationActionPolicy{
			Action:  automationAction,
			AgentId: agentId,
		}

		// Set the ID of the first action as the runnable ID
		if i == 0 {
			runnable.Id = automationAction.ID
		}
	}

	runnable.RunConfiguration.Actions = actions
	return runnable, diags
}

func (r *automationPolicyResourceFramework) mapTagsFromState(ctx context.Context, tags []string) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tags == nil || len(tags) == 0 {
		return nil, diags
	}

	return tags, diags
}
