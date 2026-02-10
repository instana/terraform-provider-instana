package maintenancewindowconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
)

// NewMaintenanceWindowConfigResourceHandle creates the resource handle for Maintenance Window Configuration
func NewMaintenanceWindowConfigResourceHandle() resourcehandle.ResourceHandle[*restapi.MaintenanceWindowConfig] {
	return &maintenanceWindowConfigResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaMaintenanceWindowConfig,
			Schema: schema.Schema{
				Description: MaintenanceWindowConfigDescResource,
				Attributes: map[string]schema.Attribute{
					MaintenanceWindowConfigFieldID: schema.StringAttribute{
						Computed:    true,
						Description: MaintenanceWindowConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					MaintenanceWindowConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: MaintenanceWindowConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 256),
						},
					},
					MaintenanceWindowConfigFieldQuery: schema.StringAttribute{
						Required:    true,
						Description: MaintenanceWindowConfigDescQuery,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 2048),
						},
					},
					MaintenanceWindowConfigFieldScheduling: schema.SingleNestedAttribute{
						Required:    true,
						Description: MaintenanceWindowConfigDescScheduling,
						Attributes:  buildSchedulingAttributes(),
					},
					MaintenanceWindowConfigFieldTagFilterExpressionEnabled: schema.BoolAttribute{
						Optional:    true,
						Description: MaintenanceWindowConfigDescTagFilterExpressionEnabled,
					},
					MaintenanceWindowConfigFieldTagFilterExpression: schema.StringAttribute{
						Optional:    true,
						Description: MaintenanceWindowConfigDescTagFilterExpression,
					},
				},
			},
			SchemaVersion: 0,
		},
	}
}

// buildSchedulingAttributes constructs the attributes for scheduling
func buildSchedulingAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SchedulingFieldStart: schema.Int64Attribute{
			Required:    true,
			Description: SchedulingDescStart,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		SchedulingFieldDuration: schema.SingleNestedAttribute{
			Required:    true,
			Description: SchedulingDescDuration,
			Attributes:  buildDurationAttributes(),
		},
		SchedulingFieldType: schema.StringAttribute{
			Required:    true,
			Description: SchedulingDescType,
			Validators: []validator.String{
				stringvalidator.OneOf(SupportedSchedulingTypes...),
			},
		},
		SchedulingFieldRrule: schema.StringAttribute{
			Optional:    true,
			Description: SchedulingDescRrule,
		},
		SchedulingFieldTimezoneId: schema.StringAttribute{
			Optional:    true,
			Description: SchedulingDescTimezoneId,
		},
	}
}

// buildDurationAttributes constructs the attributes for duration
func buildDurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		DurationFieldAmount: schema.Int64Attribute{
			Required:    true,
			Description: DurationDescAmount,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		DurationFieldUnit: schema.StringAttribute{
			Required:    true,
			Description: DurationDescUnit,
			Validators: []validator.String{
				stringvalidator.OneOf(SupportedDurationUnits...),
			},
		},
	}
}

// ============================================================================
// Resource Implementation
// ============================================================================

type maintenanceWindowConfigResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *maintenanceWindowConfigResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for maintenance window configurations
func (r *maintenanceWindowConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.MaintenanceWindowConfig] {
	return api.MaintenanceWindowConfigs()
}

// SetComputedFields sets computed fields in the plan
func (r *maintenanceWindowConfigResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ============================================================================
// State Management
// ============================================================================

// UpdateState updates the Terraform state with the maintenance window configuration data from the API
func (r *maintenanceWindowConfigResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, config *restapi.MaintenanceWindowConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create base model
	model := MaintenanceWindowConfigModel{
		ID:    types.StringValue(config.ID),
		Name:  types.StringValue(config.Name),
		Query: types.StringValue(config.Query),
	}

	// Map scheduling
	schedulingDiags := r.mapSchedulingToModel(ctx, config.Scheduling, &model)
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return diags
	}

	// Map tag filter expression
	if config.TagFilterExpressionEnabled != nil {
		model.TagFilterExpressionEnabled = types.BoolValue(*config.TagFilterExpressionEnabled)
	} else {
		model.TagFilterExpressionEnabled = types.BoolNull()
	}

	if config.TagFilterExpression != nil {
		tagFilterString, err := tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			diags.AddError("Failed to map tag filter expression", err.Error())
			return diags
		}
		if tagFilterString != nil {
			model.TagFilterExpression = types.StringValue(*tagFilterString)
		} else {
			model.TagFilterExpression = types.StringNull()
		}
	} else {
		model.TagFilterExpression = types.StringNull()
	}

	// Set state
	setStateDiags := state.Set(ctx, &model)
	diags.Append(setStateDiags...)

	return diags
}

// mapSchedulingToModel maps the scheduling configuration to the model
func (r *maintenanceWindowConfigResource) mapSchedulingToModel(ctx context.Context, scheduling *restapi.MaintenanceScheduling, model *MaintenanceWindowConfigModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if scheduling == nil {
		diags.AddError("Invalid scheduling", "Scheduling configuration is required")
		return diags
	}

	// Map duration
	durationAttrs := map[string]attr.Value{
		DurationFieldAmount: types.Int64Value(scheduling.Duration.Amount),
		DurationFieldUnit:   types.StringValue(scheduling.Duration.Unit),
	}

	durationObj, durationDiags := types.ObjectValue(
		map[string]attr.Type{
			DurationFieldAmount: types.Int64Type,
			DurationFieldUnit:   types.StringType,
		},
		durationAttrs,
	)
	diags.Append(durationDiags...)
	if diags.HasError() {
		return diags
	}

	// Map scheduling
	schedulingAttrs := map[string]attr.Value{
		SchedulingFieldStart:    types.Int64Value(scheduling.Start),
		SchedulingFieldDuration: durationObj,
		SchedulingFieldType:     types.StringValue(scheduling.Type),
	}

	if scheduling.Rrule != nil {
		schedulingAttrs[SchedulingFieldRrule] = types.StringValue(*scheduling.Rrule)
	} else {
		schedulingAttrs[SchedulingFieldRrule] = types.StringNull()
	}

	if scheduling.TimezoneId != nil {
		schedulingAttrs[SchedulingFieldTimezoneId] = types.StringValue(*scheduling.TimezoneId)
	} else {
		schedulingAttrs[SchedulingFieldTimezoneId] = types.StringNull()
	}

	schedulingObj, schedulingDiags := types.ObjectValue(
		map[string]attr.Type{
			SchedulingFieldStart: types.Int64Type,
			SchedulingFieldDuration: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					DurationFieldAmount: types.Int64Type,
					DurationFieldUnit:   types.StringType,
				},
			},
			SchedulingFieldType:       types.StringType,
			SchedulingFieldRrule:      types.StringType,
			SchedulingFieldTimezoneId: types.StringType,
		},
		schedulingAttrs,
	)
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return diags
	}

	model.Scheduling = schedulingObj
	return diags
}

// ============================================================================
// Data Object Mapping
// ============================================================================

// MapStateToDataObject maps the Terraform state to the API data object
func (r *maintenanceWindowConfigResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.MaintenanceWindowConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model MaintenanceWindowConfigModel

	// Get model from plan or state
	if plan != nil {
		getDiags := plan.Get(ctx, &model)
		diags.Append(getDiags...)
	} else if state != nil {
		getDiags := state.Get(ctx, &model)
		diags.Append(getDiags...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Create base config
	config := &restapi.MaintenanceWindowConfig{
		ID:    model.ID.ValueString(),
		Name:  model.Name.ValueString(),
		Query: model.Query.ValueString(),
	}

	// Map scheduling
	schedulingDiags := r.mapSchedulingFromModel(ctx, model.Scheduling, config)
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Map tag filter expression
	if !model.TagFilterExpressionEnabled.IsNull() {
		enabled := model.TagFilterExpressionEnabled.ValueBool()
		config.TagFilterExpressionEnabled = &enabled
	}

	if !model.TagFilterExpression.IsNull() && model.TagFilterExpression.ValueString() != "" {
		tagFilterExpr, err := tagfilter.ParseExpression(model.TagFilterExpression.ValueString())
		if err != nil {
			diags.AddError("Failed to parse tag filter expression", err.Error())
			return nil, diags
		}
		config.TagFilterExpression = tagFilterExpr
	}

	return config, diags
}

// mapSchedulingFromModel maps the scheduling from the model to the API object
func (r *maintenanceWindowConfigResource) mapSchedulingFromModel(ctx context.Context, schedulingObj types.Object, config *restapi.MaintenanceWindowConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if schedulingObj.IsNull() {
		diags.AddError("Invalid scheduling", "Scheduling configuration is required")
		return diags
	}

	var schedulingModel MaintenanceSchedulingModel
	schedulingDiags := schedulingObj.As(ctx, &schedulingModel, basetypes.ObjectAsOptions{})
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return diags
	}

	// Map duration
	var durationModel MaintenanceDurationModel
	durationDiags := schedulingModel.Duration.As(ctx, &durationModel, basetypes.ObjectAsOptions{})
	diags.Append(durationDiags...)
	if diags.HasError() {
		return diags
	}

	duration := &restapi.MaintenanceDuration{
		Amount: durationModel.Amount.ValueInt64(),
		Unit:   durationModel.Unit.ValueString(),
	}

	// Create scheduling
	scheduling := &restapi.MaintenanceScheduling{
		Start:    schedulingModel.Start.ValueInt64(),
		Duration: duration,
		Type:     schedulingModel.Type.ValueString(),
	}

	if !schedulingModel.Rrule.IsNull() {
		rrule := schedulingModel.Rrule.ValueString()
		scheduling.Rrule = &rrule
	}

	if !schedulingModel.TimezoneId.IsNull() {
		timezoneId := schedulingModel.TimezoneId.ValueString()
		scheduling.TimezoneId = &timezoneId
	}

	config.Scheduling = scheduling
	return diags
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *maintenanceWindowConfigResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return nil
}

// Made with Bob
