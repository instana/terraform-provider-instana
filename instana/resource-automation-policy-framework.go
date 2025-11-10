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

// runnable types
const actionRunnable = "action"

// AutomationPolicyModel represents the data model for the automation policy resource
type AutomationPolicyModel struct {
	ID                types.String             `tfsdk:"id"`
	Name              types.String             `tfsdk:"name"`
	Description       types.String             `tfsdk:"description"`
	Tags              types.List               `tfsdk:"tags"`
	Trigger           TriggerModel             `tfsdk:"trigger"`
	TypeConfiguration []TypeConfigurationModel `tfsdk:"type_configuration"`
}

// TriggerModel represents a trigger in the automation policy
type TriggerModel struct {
	ID          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
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
	Action  AutomationActionModel `tfsdk:"action"`
	AgentID types.String          `tfsdk:"agent_id"`
}

// NewAutomationPolicyResourceHandleFramework creates the resource handle for Automation Policies
func NewAutomationPolicyResourceHandleFramework() ResourceHandleFramework[*restapi.AutomationPolicy] {
	return &automationPolicyResourceFramework{
		metaData: ResourceMetaDataFramework{
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
					AutomationPolicyFieldName: schema.StringAttribute{
						Required:    true,
						Description: AutomationPolicyDescName,
					},
					AutomationPolicyFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: AutomationPolicyDescDescription,
					},
					AutomationPolicyFieldTags: schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: AutomationPolicyDescTags,
					},
					AutomationPolicyFieldTrigger: schema.SingleNestedAttribute{
						Required:    true,
						Description: AutomationPolicyDescTrigger,
						Attributes: map[string]schema.Attribute{
							AutomationPolicyFieldId: schema.StringAttribute{
								Required:    true,
								Description: AutomationPolicyDescTriggerID,
							},
							AutomationPolicyFieldType: schema.StringAttribute{
								Required:    true,
								Description: AutomationPolicyDescTriggerType,
								Validators: []validator.String{
									stringvalidator.OneOf(supportedTriggerTypes...),
								},
							},
							AutomationPolicyFieldName: schema.StringAttribute{
								Optional:    true,
								Description: AutomationPolicyDescTriggerName,
							},
							AutomationPolicyFieldDescription: schema.StringAttribute{
								Optional:    true,
								Description: AutomationPolicyDescTriggerDescription,
							},
						},
					},
					AutomationPolicyFieldTypeConfiguration: schema.ListNestedAttribute{
						Required:    true,
						Description: AutomationPolicyDescTypeConfiguration,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								AutomationPolicyFieldName: schema.StringAttribute{
									Required:    true,
									Description: AutomationPolicyDescTypeConfigurationName,
									Validators: []validator.String{
										stringvalidator.OneOf(supportedPolicyTypes...),
									},
								},
								AutomationPolicyFieldCondition: schema.SingleNestedAttribute{
									Optional:    true,
									Description: AutomationPolicyDescCondition,
									Attributes: map[string]schema.Attribute{
										AutomationPolicyFieldQuery: schema.StringAttribute{
											Required:    true,
											Description: AutomationPolicyDescConditionQuery,
										},
									},
								},
								AutomationPolicyFieldAction: schema.ListNestedAttribute{
									Required:    true,
									Description: AutomationPolicyDescAction,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.SingleNestedAttribute{
												Required:    true,
												Description: AutomationPolicyDescActionAction,
												Attributes:  GetAutomationActionSchemaAttributes(),
											},
											AutomationPolicyFieldAgentId: schema.StringAttribute{
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
		tagsList, d := r.mapTagsToState(ctx, policy.Tags)
		diags.Append(d...)
		if !diags.HasError() {
			model.Tags = tagsList
		}
	} else {
		model.Tags = types.ListNull(types.StringType)
	}

	// Map trigger
	model.Trigger = r.mapTriggerToState(&policy.Trigger, model.Trigger)

	// Map type configurations
	model.TypeConfiguration = r.mapTypeConfigurationsToState(ctx, policy.TypeConfigurations)

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
					AutomationPolicyErrMappingTags,
					fmt.Sprintf(AutomationPolicyErrTagNotString, i),
				)
				return types.ListNull(types.StringType), diags
			}
		}
		return types.ListValueMust(types.StringType, elements), diags
	default:
		diags.AddError(
			AutomationPolicyErrMappingTags,
			AutomationPolicyErrTagsFormat,
		)
		return types.ListNull(types.StringType), diags
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
		triggerModel.Description = types.StringValue(trigger.Description)
	}
	if triggerModel.Name.IsNull() || triggerModel.Name.IsUnknown() {
		triggerModel.Name = types.StringValue(trigger.Name)
	}
	// The existing value from plan is already preserved in UpdateState function

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
		tags, _ := MapTagsToState(ctx, actionPolicy.Action.Tags)
		inputParams, _ := MapInputParametersToState(ctx, actionPolicy.Action.InputParameters)

		actionModel := AutomationActionModel{
			ID:             types.StringValue(actionPolicy.Action.ID),
			Name:           types.StringValue(actionPolicy.Action.Name),
			Description:    types.StringValue(actionPolicy.Action.Description),
			Tags:           tags,
			InputParameter: inputParams,
		}

		// Map action type-specific fields using the common function
		MapActionTypeFieldsToState(ctx, &actionPolicy.Action, &actionModel)

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

func (r *automationPolicyResourceFramework) mapInputParametersToStateMap(ctx context.Context, inputParams []restapi.Parameter) types.Map {
	if len(inputParams) == 0 {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value)
	for _, param := range inputParams {
		elements[param.Name] = types.StringValue(param.Value)
	}

	return types.MapValueMust(types.StringType, elements)
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
		inputParams, d := MapInputParametersFromState(ctx, actionModel)
		diags.Append(d...)
		if diags.HasError() {
			return runnable, diags
		}

		// Map action type and fields
		actionType, fields, d := MapActionTypeAndFieldsFromState(ctx, actionModel)
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
			tags, d := MapTagsFromState(ctx, actionModel.Tags)
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
