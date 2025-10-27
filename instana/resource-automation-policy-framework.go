package instana

import (
	"context"
	"fmt"

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
)

// ResourceInstanaAutomationPolicyFramework the name of the terraform-provider-instana resource to manage automation policies
const ResourceInstanaAutomationPolicyFramework = "automation_policy"

// AutomationPolicyModel represents the data model for the automation policy resource
type AutomationPolicyModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Tags              types.List   `tfsdk:"tags"`
	Trigger           types.List   `tfsdk:"trigger"`
	TypeConfiguration types.List   `tfsdk:"type_configuration"`
}

// TriggerModel represents a trigger in the automation policy
type TriggerModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

// TypeConfigurationModel represents a type configuration in the automation policy
type TypeConfigurationModel struct {
	Name      types.String `tfsdk:"name"`
	Condition types.List   `tfsdk:"condition"`
	Action    types.List   `tfsdk:"action"`
}

// ConditionModel represents a condition in the automation policy
type ConditionModel struct {
	Query types.String `tfsdk:"query"`
}

// ActionModel represents an action in the automation policy
type ActionModel struct {
	ActionID        types.String `tfsdk:"action_id"`
	AgentID         types.String `tfsdk:"agent_id"`
	InputParameters types.Map    `tfsdk:"input_parameters"`
}

// NewAutomationPolicyResourceHandleFramework creates the resource handle for Automation Policies
func NewAutomationPolicyResourceHandleFramework() ResourceHandleFramework[*restapi.AutomationPolicy] {
	return &automationPolicyResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAutomationPolicyFramework,
			Schema: schema.Schema{
				Description: "This resource manages automation policies in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the automation policy.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AutomationPolicyFieldName: schema.StringAttribute{
						Required:    true,
						Description: "The name of the automation policy.",
					},
					AutomationPolicyFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: "The description of the automation policy.",
					},
					AutomationPolicyFieldTags: schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The tags of the automation policy.",
					},
				},
				Blocks: map[string]schema.Block{
					AutomationPolicyFieldTrigger: schema.ListNestedBlock{
						Description: "The trigger for the automation policy.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								AutomationPolicyFieldId: schema.StringAttribute{
									Required:    true,
									Description: "Trigger (Instana event or Smart Alert) identifier.",
								},
								AutomationPolicyFieldType: schema.StringAttribute{
									Required:    true,
									Description: "Instana event or Smart Alert type.",
									Validators: []validator.String{
										stringvalidator.OneOf(supportedTriggerTypes...),
									},
								},
							},
						},
					},
					AutomationPolicyFieldTypeConfiguration: schema.ListNestedBlock{
						Description: "A list of configurations with the list of actions to run and the mode (automatic or manual).",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								AutomationPolicyFieldName: schema.StringAttribute{
									Required:    true,
									Description: "The policy type.",
									Validators: []validator.String{
										stringvalidator.OneOf(supportedPolicyTypes...),
									},
								},
							},
							Blocks: map[string]schema.Block{
								AutomationPolicyFieldCondition: schema.ListNestedBlock{
									Description: "The condition that selects a list of entities on which the policy is run. Only for automatic policy type.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											AutomationPolicyFieldQuery: schema.StringAttribute{
												Required:    true,
												Description: "Dynamic Focus Query string that selects a list of entities on which the policy is run.",
											},
										},
									},
								},
								AutomationPolicyFieldAction: schema.ListNestedBlock{
									Description: "The configuration for the automation action.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											AutomationPolicyFieldActionId: schema.StringAttribute{
												Required:    true,
												Description: "The identifier for the automation action.",
											},
											AutomationPolicyFieldAgentId: schema.StringAttribute{
												Optional:    true,
												Description: "The identifier of the agent host.",
											},
											AutomationPolicyFieldInputParameters: schema.MapAttribute{
												ElementType: types.StringType,
												Optional:    true,
												Description: "Optional map with input parameters name and value.",
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
	metaData ResourceMetaDataFramework
}

func (r *automationPolicyResourceFramework) MetaData() *ResourceMetaDataFramework {
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

	// Create a model and populate it with values from the policy
	model := AutomationPolicyModel{
		ID:          types.StringValue(policy.ID),
		Name:        types.StringValue(policy.Name),
		Description: types.StringValue(policy.Description),
	}

	// Handle tags
	if policy.Tags != nil {
		tagsList, d := r.mapTagsToState(ctx, policy.Tags)
		diags.Append(d...)
		if !diags.HasError() {
			model.Tags = tagsList
		}
	} else {
		model.Tags = types.ListNull(types.StringType)
	}

	// Map trigger
	trigger, d := r.mapTriggerToState(ctx, &policy.Trigger)
	diags.Append(d...)
	if !diags.HasError() {
		model.Trigger = trigger
	}

	// Map type configurations
	typeConfigs, d := r.mapTypeConfigurationsToState(ctx, policy.TypeConfigurations)
	diags.Append(d...)
	if !diags.HasError() {
		model.TypeConfiguration = typeConfigs
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *automationPolicyResourceFramework) mapTagsToState(ctx context.Context, tags interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tags == nil {
		return types.ListNull(types.StringType), diags
	}

	// Handle tags based on their type
	switch v := tags.(type) {
	case []interface{}:
		elements := make([]attr.Value, len(v))
		for i, tag := range v {
			if strTag, ok := tag.(string); ok {
				elements[i] = types.StringValue(strTag)
			} else {
				diags.AddError(
					"Error mapping tags",
					fmt.Sprintf("Tag at index %d is not a string", i),
				)
				return types.ListNull(types.StringType), diags
			}
		}
		return types.ListValueMust(types.StringType, elements), diags
	default:
		diags.AddError(
			"Error mapping tags",
			"Tags are not in the expected format",
		)
		return types.ListNull(types.StringType), diags
	}
}

func (r *automationPolicyResourceFramework) mapTriggerToState(ctx context.Context, trigger *restapi.Trigger) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	triggerObj := map[string]attr.Value{
		AutomationPolicyFieldId:   types.StringValue(trigger.Id),
		AutomationPolicyFieldType: types.StringValue(trigger.Type),
	}

	objValue, d := types.ObjectValue(
		map[string]attr.Type{
			AutomationPolicyFieldId:   types.StringType,
			AutomationPolicyFieldType: types.StringType,
		},
		triggerObj,
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationPolicyFieldId:   types.StringType,
				AutomationPolicyFieldType: types.StringType,
			},
		}), diags
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationPolicyFieldId:   types.StringType,
				AutomationPolicyFieldType: types.StringType,
			},
		},
		[]attr.Value{objValue},
	), diags
}

func (r *automationPolicyResourceFramework) mapTypeConfigurationsToState(ctx context.Context, typeConfigs []restapi.TypeConfiguration) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make([]attr.Value, len(typeConfigs))

	for i, typeConfig := range typeConfigs {
		// Map condition
		condition, d := r.mapConditionToState(ctx, typeConfig.Condition)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		// Map actions
		actions, d := r.mapActionsToState(ctx, &typeConfig.Runnable)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		typeConfigObj := map[string]attr.Value{
			AutomationPolicyFieldName:      types.StringValue(typeConfig.Name),
			AutomationPolicyFieldCondition: condition,
			AutomationPolicyFieldAction:    actions,
		}

		objValue, d := types.ObjectValue(
			map[string]attr.Type{
				AutomationPolicyFieldName:      types.StringType,
				AutomationPolicyFieldCondition: types.ListType{ElemType: types.ObjectType{}},
				AutomationPolicyFieldAction:    types.ListType{ElemType: types.ObjectType{}},
			},
			typeConfigObj,
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		elements[i] = objValue
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationPolicyFieldName:      types.StringType,
				AutomationPolicyFieldCondition: types.ListType{ElemType: types.ObjectType{}},
				AutomationPolicyFieldAction:    types.ListType{ElemType: types.ObjectType{}},
			},
		},
		elements,
	), diags
}

func (r *automationPolicyResourceFramework) mapConditionToState(ctx context.Context, condition *restapi.Condition) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if condition == nil || condition.Query == "" {
		return types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					AutomationPolicyFieldQuery: types.StringType,
				},
			},
			[]attr.Value{},
		), diags
	}

	conditionObj := map[string]attr.Value{
		AutomationPolicyFieldQuery: types.StringValue(condition.Query),
	}

	objValue, d := types.ObjectValue(
		map[string]attr.Type{
			AutomationPolicyFieldQuery: types.StringType,
		},
		conditionObj,
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationPolicyFieldQuery: types.StringType,
			},
		},
		[]attr.Value{objValue},
	), diags
}

func (r *automationPolicyResourceFramework) mapActionsToState(ctx context.Context, runnable *restapi.Runnable) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make([]attr.Value, len(runnable.RunConfiguration.Actions))

	for i, action := range runnable.RunConfiguration.Actions {
		// Map input parameters
		inputParams, d := r.mapInputParametersToState(ctx, action.InputParameterValues)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		actionObj := map[string]attr.Value{
			AutomationPolicyFieldActionId:        types.StringValue(action.Action.Id),
			AutomationPolicyFieldAgentId:         types.StringValue(action.AgentId),
			AutomationPolicyFieldInputParameters: inputParams,
		}

		objValue, d := types.ObjectValue(
			map[string]attr.Type{
				AutomationPolicyFieldActionId:        types.StringType,
				AutomationPolicyFieldAgentId:         types.StringType,
				AutomationPolicyFieldInputParameters: types.MapType{ElemType: types.StringType},
			},
			actionObj,
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		elements[i] = objValue
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationPolicyFieldActionId:        types.StringType,
				AutomationPolicyFieldAgentId:         types.StringType,
				AutomationPolicyFieldInputParameters: types.MapType{ElemType: types.StringType},
			},
		},
		elements,
	), diags
}

func (r *automationPolicyResourceFramework) mapInputParametersToState(ctx context.Context, inputParams []restapi.InputParameterValue) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make(map[string]attr.Value)

	for _, param := range inputParams {
		elements[param.Name] = types.StringValue(param.Value)
	}

	return types.MapValueMust(types.StringType, elements), diags
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

func (r *automationPolicyResourceFramework) mapTriggerFromState(ctx context.Context, triggerList types.List) (restapi.Trigger, diag.Diagnostics) {
	var diags diag.Diagnostics
	var trigger restapi.Trigger

	if triggerList.IsNull() {
		return trigger, diags
	}

	var triggerModels []TriggerModel
	diags.Append(triggerList.ElementsAs(ctx, &triggerModels, false)...)
	if diags.HasError() {
		return trigger, diags
	}

	if len(triggerModels) > 0 {
		triggerModel := triggerModels[0]
		trigger.Id = triggerModel.ID.ValueString()
		trigger.Type = triggerModel.Type.ValueString()
	}

	return trigger, diags
}

func (r *automationPolicyResourceFramework) mapTypeConfigurationsFromState(ctx context.Context, typeConfigList types.List) ([]restapi.TypeConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	var typeConfigurations []restapi.TypeConfiguration

	if typeConfigList.IsNull() {
		return typeConfigurations, diags
	}

	var typeConfigModels []TypeConfigurationModel
	diags.Append(typeConfigList.ElementsAs(ctx, &typeConfigModels, false)...)
	if diags.HasError() {
		return typeConfigurations, diags
	}

	typeConfigurations = make([]restapi.TypeConfiguration, len(typeConfigModels))
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

func (r *automationPolicyResourceFramework) mapConditionFromState(ctx context.Context, conditionList types.List) (*restapi.Condition, diag.Diagnostics) {
	var diags diag.Diagnostics
	condition := &restapi.Condition{
		Query: "",
	}

	if conditionList.IsNull() {
		return condition, diags
	}

	var conditionModels []ConditionModel
	diags.Append(conditionList.ElementsAs(ctx, &conditionModels, false)...)
	if diags.HasError() {
		return condition, diags
	}

	if len(conditionModels) > 0 {
		conditionModel := conditionModels[0]
		condition.Query = conditionModel.Query.ValueString()
	}

	return condition, diags
}

func (r *automationPolicyResourceFramework) mapRunnableFromState(ctx context.Context, actionList types.List) (restapi.Runnable, diag.Diagnostics) {
	var diags diag.Diagnostics
	runnable := restapi.Runnable{
		Type: actionRunnable,
		RunConfiguration: restapi.RunConfiguration{
			Actions: []restapi.ActionConfiguration{},
		},
	}

	if actionList.IsNull() {
		return runnable, diags
	}

	var actionModels []ActionModel
	diags.Append(actionList.ElementsAs(ctx, &actionModels, false)...)
	if diags.HasError() {
		return runnable, diags
	}

	actions := make([]restapi.ActionConfiguration, len(actionModels))
	for i, actionModel := range actionModels {
		// Map input parameters
		inputParams, d := r.mapInputParametersFromState(ctx, actionModel.InputParameters)
		diags.Append(d...)
		if diags.HasError() {
			return runnable, diags
		}

		actionId := actionModel.ActionID.ValueString()
		actions[i] = restapi.ActionConfiguration{
			Action: restapi.Action{
				Id: actionId,
			},
			AgentId:              actionModel.AgentID.ValueString(),
			InputParameterValues: inputParams,
		}

		// Set the ID of the first action as the runnable ID
		if i == 0 {
			runnable.Id = actionId
		}
	}

	runnable.RunConfiguration.Actions = actions
	return runnable, diags
}

func (r *automationPolicyResourceFramework) mapInputParametersFromState(ctx context.Context, inputParamsMap types.Map) ([]restapi.InputParameterValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var inputParams []restapi.InputParameterValue

	if inputParamsMap.IsNull() {
		return inputParams, diags
	}

	elements := make(map[string]string)
	diags.Append(inputParamsMap.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return inputParams, diags
	}

	inputParams = make([]restapi.InputParameterValue, 0, len(elements))
	for name, value := range elements {
		inputParams = append(inputParams, restapi.InputParameterValue{
			Name:  name,
			Value: value,
		})
	}

	return inputParams, diags
}

func (r *automationPolicyResourceFramework) mapTagsFromState(ctx context.Context, tagsList types.List) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagsList.IsNull() {
		return nil, diags
	}

	var tags []string
	diags.Append(tagsList.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	return tags, diags
}

// Made with Bob
