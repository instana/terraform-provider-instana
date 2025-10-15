package tf_framework

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &sloCorrectionConfigResource{}
	_ resource.ResourceWithConfigure   = &sloCorrectionConfigResource{}
	_ resource.ResourceWithImportState = &sloCorrectionConfigResource{}
)

// generateRandomID generates a random ID for resources
func generateRandomID() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// NewSloCorrectionConfigResource creates a new resource for SLO Correction Config
func NewSloCorrectionConfigResource() resource.Resource {
	return &sloCorrectionConfigResource{}
}

// sloCorrectionConfigResource is the resource implementation
type sloCorrectionConfigResource struct {
	instanaAPI restapi.InstanaAPI
}

// Metadata returns the resource type name
func (r *sloCorrectionConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slo_correction_config"
}

// Schema defines the schema for the resource
func (r *sloCorrectionConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource manages SLO Correction Configurations in Instana.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the SLO Correction Config.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SLO Correction Config.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 256),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "The description of the SLO Correction Config.",
			},
			"active": schema.BoolAttribute{
				Required:    true,
				Description: "Indicates whether the Correction Config is active.",
			},
			"slo_ids": schema.SetAttribute{
				Required:    true,
				Description: "A set of SLO IDs that this correction config applies to.",
				ElementType: types.StringType,
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Description: "A list of tags to be associated with the SLO Correction Config.",
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"scheduling": schema.ListNestedBlock{
				Description: "Scheduling configuration for the SLO Correction Config.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"start_time": schema.Int64Attribute{
							Required:    true,
							Description: "The start time of the scheduling in Unix timestamp in milliseconds.",
						},
						"duration": schema.Int64Attribute{
							Required:    true,
							Description: "The duration of the scheduling in the specified unit.",
						},
						"duration_unit": schema.StringAttribute{
							Required:    true,
							Description: "The unit of the duration (e.g.,'MINUTE' 'HOUR', 'DAY').",
							Validators: []validator.String{
								stringvalidator.OneOf("MINUTE", "HOUR", "DAY"),
							},
						},
						"recurrent_rule": schema.StringAttribute{
							Optional:    true,
							Description: "Recurrent rule for scheduling, if applicable.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *sloCorrectionConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(restapi.InstanaAPI)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *restapi.ProviderMeta, got: "+string(rune(len(req.ProviderData.(string)))),
		)
		return
	}

	r.instanaAPI = providerMeta
}

// Create creates a new SLO Correction Config
func (r *sloCorrectionConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan SloCorrectionConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate ID
	// Generate a random ID
	id := "slo-correction-" + generateRandomID()
	plan.ID = types.StringValue(id)

	// Map to API model
	sloCorrectionConfig, mapDiags := r.mapModelToAPIObject(ctx, plan)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new SLO Correction Config
	createdConfig, err := r.instanaAPI.SloCorrectionConfig().Create(sloCorrectionConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SLO Correction Config",
			"Could not create SLO Correction Config, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to model
	mapDiags = r.mapAPIObjectToModel(ctx, createdConfig, &plan)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *sloCorrectionConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state SloCorrectionConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get SLO Correction Config from Instana
	sloCorrectionConfig, err := r.instanaAPI.SloCorrectionConfig().GetOne(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SLO Correction Config",
			"Could not read SLO Correction Config ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to model
	mapDiags := r.mapAPIObjectToModel(ctx, sloCorrectionConfig, &state)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *sloCorrectionConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan SloCorrectionConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map to API model
	sloCorrectionConfig, mapDiags := r.mapModelToAPIObject(ctx, plan)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update SLO Correction Config
	updatedConfig, err := r.instanaAPI.SloCorrectionConfig().Update(sloCorrectionConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SLO Correction Config",
			"Could not update SLO Correction Config, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to model
	mapDiags = r.mapAPIObjectToModel(ctx, updatedConfig, &plan)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *sloCorrectionConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state SloCorrectionConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete SLO Correction Config
	err := r.instanaAPI.SloCorrectionConfig().DeleteByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SLO Correction Config",
			"Could not delete SLO Correction Config, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *sloCorrectionConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper functions for mapping between Terraform and API models

func (r *sloCorrectionConfigResource) mapModelToAPIObject(ctx context.Context, model SloCorrectionConfigModel) (*restapi.SloCorrectionConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

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

func (r *sloCorrectionConfigResource) mapAPIObjectToModel(ctx context.Context, apiObject *restapi.SloCorrectionConfig, model *SloCorrectionConfigModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Map basic fields
	model.ID = types.StringValue(apiObject.ID)
	model.Name = types.StringValue(apiObject.Name)
	model.Description = types.StringValue(apiObject.Description)
	model.Active = types.BoolValue(apiObject.Active)

	// Map scheduling
	schedulingObj := map[string]attr.Value{
		"start_time":     types.Int64Value(apiObject.Scheduling.StartTime),
		"duration":       types.Int64Value(int64(apiObject.Scheduling.Duration)),
		"duration_unit":  types.StringValue(string(apiObject.Scheduling.DurationUnit)),
		"recurrent_rule": types.StringValue(apiObject.Scheduling.RecurrentRule),
	}

	schedulingType := map[string]attr.Type{
		"start_time":     types.Int64Type,
		"duration":       types.Int64Type,
		"duration_unit":  types.StringType,
		"recurrent_rule": types.StringType,
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

	return diags
}

// Made with Bob
