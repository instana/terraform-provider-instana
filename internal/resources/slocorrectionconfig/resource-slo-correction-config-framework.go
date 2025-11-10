package slocorrectionconfig

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
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

// ResourceInstanaSloCorrectionConfigFramework the name of the terraform-provider-instana resource to manage SLO correction configurations
const ResourceInstanaSloCorrectionConfigFramework = "slo_correction_config"

const (
	// Slo Correction Config Field names for Terraform
	SloCorrectionConfigFieldName                    = "name"
	SloCorrectionConfigFieldFullName                = "full_name"
	SloCorrectionConfigFieldDescription             = "description"
	SloCorrectionConfigFieldActive                  = "active"
	SloCorrectionConfigFieldScheduling              = "scheduling"
	SloCorrectionConfigFieldSloIds                  = "slo_ids"
	SloCorrectionConfigFieldTags                    = "tags"
	SloCorrectionConfigFieldSchedulingStartTime     = "start_time"
	SloCorrectionConfigFieldSchedulingDuration      = "duration"
	SloCorrectionConfigFieldSchedulingDurationUnit  = "duration_unit"
	SloCorrectionConfigFieldSchedulingRecurrentRule = "recurrent_rule"
)

// SloCorrectionConfigModel represents the data model for the SLO Correction Config resource
type SloCorrectionConfigModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Active      types.Bool   `tfsdk:"active"`
	Scheduling  types.List   `tfsdk:"scheduling"`
	SloIds      types.Set    `tfsdk:"slo_ids"`
	Tags        types.Set    `tfsdk:"tags"`
}

// SchedulingModel represents the scheduling configuration for SLO Correction Config
type SchedulingModel struct {
	StartTime     types.Int64  `tfsdk:"start_time"`
	Duration      types.Int64  `tfsdk:"duration"`
	DurationUnit  types.String `tfsdk:"duration_unit"`
	RecurrentRule types.String `tfsdk:"recurrent_rule"`
	Recurrent     types.Bool   `tfsdk:"recurrent"`
}

// NewSloCorrectionConfigResourceHandleFramework creates the resource handle for SLO Correction Config
func NewSloCorrectionConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.SloCorrectionConfig] {
	return &sloCorrectionConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSloCorrectionConfigFramework,
			Schema: schema.Schema{
				Description: SloCorrectionConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: SloCorrectionConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: SloCorrectionConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					"description": schema.StringAttribute{
						Required:    true,
						Description: SloCorrectionConfigDescDescription,
					},
					"active": schema.BoolAttribute{
						Required:    true,
						Description: SloCorrectionConfigDescActive,
					},
					"slo_ids": schema.SetAttribute{
						Required:    true,
						Description: SloCorrectionConfigDescSloIds,
						ElementType: types.StringType,
					},
					"tags": schema.SetAttribute{
						Optional:    true,
						Description: SloCorrectionConfigDescTags,
						ElementType: types.StringType,
					},
				},
				Blocks: map[string]schema.Block{
					"scheduling": schema.ListNestedBlock{
						Description: SloCorrectionConfigDescScheduling,
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"start_time": schema.Int64Attribute{
									Required:    true,
									Description: SloCorrectionConfigDescStartTime,
								},
								"duration": schema.Int64Attribute{
									Required:    true,
									Description: SloCorrectionConfigDescDuration,
								},
								"duration_unit": schema.StringAttribute{
									Required:    true,
									Description: SloCorrectionConfigDescDurationUnit,
									Validators: []validator.String{
										stringvalidator.OneOf("millisecond", "second", "minute", "hour", "day", "week", "month"),
									},
								},
								"recurrent_rule": schema.StringAttribute{
									Optional:    true,
									Description: SloCorrectionConfigDescRecurrentRule,
								},
								"recurrent": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Description: SloCorrectionConfigDescRecurrent,
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

type sloCorrectionConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *sloCorrectionConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *sloCorrectionConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloCorrectionConfig] {
	return api.SloCorrectionConfig()
}

func (r *sloCorrectionConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sloCorrectionConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SloCorrectionConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SloCorrectionConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map scheduling
	var scheduling restapi.Scheduling
	if !model.Scheduling.IsNull() {
		var schedulingElements []types.Object
		diags.Append(model.Scheduling.ElementsAs(ctx, &schedulingElements, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(schedulingElements) > 0 {
			var schedulingModel SchedulingModel
			diags.Append(schedulingElements[0].As(ctx, &schedulingModel, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}

			scheduling = restapi.Scheduling{
				StartTime:    schedulingModel.StartTime.ValueInt64(),
				Duration:     int(schedulingModel.Duration.ValueInt64()),
				DurationUnit: restapi.DurationUnit(strings.ToUpper(schedulingModel.DurationUnit.ValueString())),
				Recurrent:    schedulingModel.Recurrent.ValueBool(),
			}

			if !schedulingModel.RecurrentRule.IsNull() {
				scheduling.RecurrentRule = schedulingModel.RecurrentRule.ValueString()
			}
		}
	}

	// Map SLO IDs
	var sloIds []string
	if !model.SloIds.IsNull() {
		diags.Append(model.SloIds.ElementsAs(ctx, &sloIds, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map tags
	var tags []string
	if !model.Tags.IsNull() {
		diags.Append(model.Tags.ElementsAs(ctx, &tags, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Create API object
	sloCorrectionConfig := &restapi.SloCorrectionConfig{
		ID:          model.ID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Active:      model.Active.ValueBool(),
		Scheduling:  scheduling,
		SloIds:      sloIds,
		Tags:        tags,
	}

	return sloCorrectionConfig, diags
}

func (r *sloCorrectionConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SloCorrectionConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the config
	model := SloCorrectionConfigModel{
		ID:          types.StringValue(apiObject.ID),
		Name:        types.StringValue(apiObject.Name),
		Description: types.StringValue(apiObject.Description),
		Active:      types.BoolValue(apiObject.Active),
	}

	// Map scheduling
	schedulingObj := map[string]attr.Value{
		"start_time":     types.Int64Value(apiObject.Scheduling.StartTime),
		"duration":       types.Int64Value(int64(apiObject.Scheduling.Duration)),
		"duration_unit":  types.StringValue(string(apiObject.Scheduling.DurationUnit)),
		"recurrent_rule": types.StringValue(apiObject.Scheduling.RecurrentRule),
		"recurrent":      types.BoolValue(apiObject.Scheduling.Recurrent),
	}

	schedulingType := map[string]attr.Type{
		"start_time":     types.Int64Type,
		"duration":       types.Int64Type,
		"duration_unit":  types.StringType,
		"recurrent_rule": types.StringType,
		"recurrent":      types.BoolType,
	}

	schedulingValue, schedulingDiags := types.ObjectValue(schedulingType, schedulingObj)
	diags.Append(schedulingDiags...)
	if diags.HasError() {
		return diags
	}

	schedulingList, schedulingListDiags := types.ListValue(
		types.ObjectType{AttrTypes: schedulingType},
		[]attr.Value{schedulingValue},
	)
	diags.Append(schedulingListDiags...)
	if diags.HasError() {
		return diags
	}

	model.Scheduling = schedulingList

	// Map SLO IDs
	sloIdsSet, sloIdsDiags := types.SetValueFrom(ctx, types.StringType, apiObject.SloIds)
	diags.Append(sloIdsDiags...)
	if diags.HasError() {
		return diags
	}
	model.SloIds = sloIdsSet

	// Map tags
	if len(apiObject.Tags) > 0 {
		tagsSet, tagsDiags := types.SetValueFrom(ctx, types.StringType, apiObject.Tags)
		diags.Append(tagsDiags...)
		if diags.HasError() {
			return diags
		}
		model.Tags = tagsSet
	} else {
		model.Tags = types.SetNull(types.StringType)
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}
