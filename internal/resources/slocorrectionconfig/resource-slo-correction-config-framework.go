package slocorrectionconfig

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
					"scheduling": schema.SingleNestedAttribute{
						Required:    true,
						Description: SloCorrectionConfigDescScheduling,
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
	if model.Scheduling != nil {
		scheduling = restapi.Scheduling{
			StartTime:    model.Scheduling.StartTime.ValueInt64(),
			Duration:     int(model.Scheduling.Duration.ValueInt64()),
			DurationUnit: restapi.DurationUnit(strings.ToUpper(model.Scheduling.DurationUnit.ValueString())),
			Recurrent:    model.Scheduling.Recurrent.ValueBool(),
		}

		if !model.Scheduling.RecurrentRule.IsNull() {
			scheduling.RecurrentRule = model.Scheduling.RecurrentRule.ValueString()
		}
	}

	// Map SLO IDs
	sloIds := model.SloIds

	// Map tags
	tags := model.Tags

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
	recurrentRuleValue := types.StringNull()
	if apiObject.Scheduling.RecurrentRule != "" {
		recurrentRuleValue = types.StringValue(apiObject.Scheduling.RecurrentRule)
	}

	model.Scheduling = &SchedulingModel{
		StartTime:     types.Int64Value(apiObject.Scheduling.StartTime),
		Duration:      types.Int64Value(int64(apiObject.Scheduling.Duration)),
		DurationUnit:  types.StringValue(string(apiObject.Scheduling.DurationUnit)),
		RecurrentRule: recurrentRuleValue,
		Recurrent:     types.BoolValue(apiObject.Scheduling.Recurrent),
	}

	// Map SLO IDs
	model.SloIds = apiObject.SloIds

	// Map tags
	if len(apiObject.Tags) > 0 {
		model.Tags = apiObject.Tags
	} else {
		model.Tags = []string{}
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}
