package tf_framework

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                = &syntheticAlertConfigResource{}
	_ resource.ResourceWithConfigure   = &syntheticAlertConfigResource{}
	_ resource.ResourceWithImportState = &syntheticAlertConfigResource{}
)

// generateSyntheticAlertRandomID generates a random ID for resources
func generateSyntheticAlertRandomID() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// NewSyntheticAlertConfigResource creates a new resource for Synthetic Alert Config
func NewSyntheticAlertConfigResource() resource.Resource {
	return &syntheticAlertConfigResource{}
}

// syntheticAlertConfigResource is the resource implementation
type syntheticAlertConfigResource struct {
	instanaAPI restapi.InstanaAPI
}

// Metadata returns the resource type name
func (r *syntheticAlertConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synthetic_alert_config"
}

// Schema defines the schema for the resource
func (r *syntheticAlertConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource manages Synthetic Alert Configurations in Instana.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Synthetic Alert Config.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Synthetic Alert Config.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "The description of the Synthetic Alert Config.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1024),
				},
			},
			"synthetic_test_ids": schema.SetAttribute{
				Required:    true,
				Description: "A set of Synthetic Test IDs that this alert config applies to.",
				ElementType: types.StringType,
			},
			"severity": schema.Int64Attribute{
				Optional:    true,
				Description: "The severity of the alert (5=critical, 10=warning).",
				Validators: []validator.Int64{
					int64validator.OneOf(5, 10),
				},
			},
			"tag_filter": schema.StringAttribute{
				Optional:    true,
				Description: "The tag filter expression used for this synthetic alert.",
			},
			"alert_channel_ids": schema.SetAttribute{
				Required:    true,
				Description: "A set of Alert Channel IDs.",
				ElementType: types.StringType,
			},
			"grace_period": schema.Int64Attribute{
				Optional:    true,
				Description: "The duration in milliseconds for which an alert remains open after conditions are no longer violated.",
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "Configuration for the synthetic alert rule.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alert_type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the alert rule (e.g., failure).",
							Validators: []validator.String{
								stringvalidator.OneOf("failure"),
							},
						},
						"metric_name": schema.StringAttribute{
							Required:    true,
							Description: "The metric name to monitor (e.g., status).",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 256),
							},
						},
						"aggregation": schema.StringAttribute{
							Optional:    true,
							Description: "The aggregation method {SUM,MEAN,MAX,MIN,P25,P50,P75,P90,P95,P98,P99,P99_9,P99_99,DISTINCT_COUNT,SUM_POSITIVE,PER_SECOND,INCREASE}.",
							Validators: []validator.String{
								stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTINCT_COUNT", "SUM_POSITIVE", "PER_SECOND", "INCREASE"),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 1),
				},
			},
			"time_threshold": schema.ListNestedBlock{
				Description: "Configuration for the time threshold.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the time threshold (only violationsInSequence is supported).",
							Validators: []validator.String{
								stringvalidator.OneOf("violationsInSequence"),
							},
						},
						"violations_count": schema.Int64Attribute{
							Required:    true,
							Description: "The number of violations required to trigger the alert (value between 1 and 12).",
							Validators: []validator.Int64{
								int64validator.Between(1, 12),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 1),
				},
			},
			"custom_payload_field": schema.ListNestedBlock{
				Description: "Custom payload fields for the alerting configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "The key of the custom payload field.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of a static string custom payload field.",
						},
						"dynamic_value": schema.ListNestedAttribute{
							Optional:    true,
							Description: "The value of a dynamic custom payload field.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Optional:    true,
										Description: "The key of the dynamic custom payload field.",
									},
									"tag_name": schema.StringAttribute{
										Required:    true,
										Description: "The name of the tag of the dynamic custom payload field.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *syntheticAlertConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(restapi.InstanaAPI)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected restapi.InstanaAPI, got: "+string(rune(len(req.ProviderData.(string)))),
		)
		return
	}

	r.instanaAPI = providerMeta
}

// Create creates a new Synthetic Alert Config
func (r *syntheticAlertConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan SyntheticAlertConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate ID
	id := "synthetic-alert-" + generateSyntheticAlertRandomID()
	plan.ID = types.StringValue(id)

	// Map to API model
	syntheticAlertConfig, mapDiags := r.mapModelToAPIObject(ctx, plan)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Synthetic Alert Config
	createdConfig, err := r.instanaAPI.SyntheticAlertConfigs().Create(syntheticAlertConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Synthetic Alert Config",
			"Could not create Synthetic Alert Config, unexpected error: "+err.Error(),
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
func (r *syntheticAlertConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state SyntheticAlertConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Synthetic Alert Config from Instana
	syntheticAlertConfig, err := r.instanaAPI.SyntheticAlertConfigs().GetOne(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Synthetic Alert Config",
			"Could not read Synthetic Alert Config ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response to model
	mapDiags := r.mapAPIObjectToModel(ctx, syntheticAlertConfig, &state)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *syntheticAlertConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan SyntheticAlertConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map to API model
	syntheticAlertConfig, mapDiags := r.mapModelToAPIObject(ctx, plan)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Synthetic Alert Config
	updatedConfig, err := r.instanaAPI.SyntheticAlertConfigs().Update(syntheticAlertConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Synthetic Alert Config",
			"Could not update Synthetic Alert Config, unexpected error: "+err.Error(),
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
func (r *syntheticAlertConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state SyntheticAlertConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Synthetic Alert Config
	err := r.instanaAPI.SyntheticAlertConfigs().DeleteByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Synthetic Alert Config",
			"Could not delete Synthetic Alert Config, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *syntheticAlertConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper functions for mapping between Terraform and API models

func (r *syntheticAlertConfigResource) mapModelToAPIObject(ctx context.Context, model SyntheticAlertConfigModel) (*restapi.SyntheticAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Map rule
	var rule restapi.SyntheticAlertRule
	if !model.Rule.IsNull() {
		var ruleElements []types.Object
		diags.Append(model.Rule.ElementsAs(ctx, &ruleElements, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(ruleElements) > 0 {
			var ruleModel SyntheticAlertRuleModel
			diags.Append(ruleElements[0].As(ctx, &ruleModel, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}

			rule = restapi.SyntheticAlertRule{
				AlertType:  ruleModel.AlertType.ValueString(),
				MetricName: ruleModel.MetricName.ValueString(),
			}

			if !ruleModel.Aggregation.IsNull() {
				rule.Aggregation = ruleModel.Aggregation.ValueString()
			}
		}
	}

	// Map time threshold
	var timeThreshold restapi.SyntheticAlertTimeThreshold
	if !model.TimeThreshold.IsNull() {
		var timeThresholdElements []types.Object
		diags.Append(model.TimeThreshold.ElementsAs(ctx, &timeThresholdElements, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(timeThresholdElements) > 0 {
			var timeThresholdModel SyntheticAlertTimeThresholdModel
			diags.Append(timeThresholdElements[0].As(ctx, &timeThresholdModel, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}

			timeThreshold = restapi.SyntheticAlertTimeThreshold{
				Type:            timeThresholdModel.Type.ValueString(),
				ViolationsCount: int(timeThresholdModel.ViolationsCount.ValueInt64()),
			}
		}
	}

	// Map synthetic test IDs
	var syntheticTestIds []string
	if !model.SyntheticTestIds.IsNull() {
		diags.Append(model.SyntheticTestIds.ElementsAs(ctx, &syntheticTestIds, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map alert channel IDs
	var alertChannelIds []string
	if !model.AlertChannelIds.IsNull() {
		diags.Append(model.AlertChannelIds.ElementsAs(ctx, &alertChannelIds, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map tag filter
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		var err error
		tagFilter, err = mapTagFilterExpressionFromSchema(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError(
				"Error parsing tag filter",
				"Could not parse tag filter: "+err.Error(),
			)
			return nil, diags
		}
	}

	// Map custom payload fields
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() {
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = BuildCustomPayloadFieldsTyped(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
	}

	// Create API object
	syntheticAlertConfig := &restapi.SyntheticAlertConfig{
		ID:                    model.ID.ValueString(),
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		SyntheticTestIds:      syntheticTestIds,
		TagFilterExpression:   tagFilter,
		Rule:                  rule,
		AlertChannelIds:       alertChannelIds,
		TimeThreshold:         timeThreshold,
		CustomerPayloadFields: customerPayloadFields,
	}

	// Set severity if present
	if !model.Severity.IsNull() {
		syntheticAlertConfig.Severity = int(model.Severity.ValueInt64())
	}

	// Set grace period if present
	if !model.GracePeriod.IsNull() {
		syntheticAlertConfig.GracePeriod = model.GracePeriod.ValueInt64()
	}

	return syntheticAlertConfig, diags
}

func (r *syntheticAlertConfigResource) mapAPIObjectToModel(ctx context.Context, apiObject *restapi.SyntheticAlertConfig, model *SyntheticAlertConfigModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Map basic fields
	model.ID = types.StringValue(apiObject.ID)
	model.Name = types.StringValue(apiObject.Name)
	model.Description = types.StringValue(apiObject.Description)
	model.Severity = types.Int64Value(int64(apiObject.Severity))

	// Map grace period if present
	if apiObject.GracePeriod > 0 {
		model.GracePeriod = types.Int64Value(apiObject.GracePeriod)
	} else {
		model.GracePeriod = types.Int64Null()
	}

	// Map tag filter
	if apiObject.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(apiObject.TagFilterExpression)
		if err != nil {
			diags.AddError(
				"Error normalizing tag filter",
				"Could not normalize tag filter: "+err.Error(),
			)
			return diags
		}
		model.TagFilter = types.StringValue(*normalizedTagFilterString)
	} else {
		model.TagFilter = types.StringNull()
	}

	// Map rule
	ruleObj := map[string]attr.Value{
		"alert_type":  types.StringValue(apiObject.Rule.AlertType),
		"metric_name": types.StringValue(apiObject.Rule.MetricName),
	}

	if apiObject.Rule.Aggregation != "" {
		ruleObj["aggregation"] = types.StringValue(apiObject.Rule.Aggregation)
	} else {
		ruleObj["aggregation"] = types.StringNull()
	}

	ruleType := map[string]attr.Type{
		"alert_type":  types.StringType,
		"metric_name": types.StringType,
		"aggregation": types.StringType,
	}

	ruleValue, ruleDiags := types.ObjectValue(ruleType, ruleObj)
	diags.Append(ruleDiags...)
	if diags.HasError() {
		return diags
	}

	ruleList, ruleListDiags := types.ListValue(
		types.ObjectType{AttrTypes: ruleType},
		[]attr.Value{ruleValue},
	)
	diags.Append(ruleListDiags...)
	if diags.HasError() {
		return diags
	}

	model.Rule = ruleList

	// Map time threshold
	timeThresholdObj := map[string]attr.Value{
		"type":             types.StringValue(apiObject.TimeThreshold.Type),
		"violations_count": types.Int64Value(int64(apiObject.TimeThreshold.ViolationsCount)),
	}

	timeThresholdType := map[string]attr.Type{
		"type":             types.StringType,
		"violations_count": types.Int64Type,
	}

	timeThresholdValue, timeThresholdDiags := types.ObjectValue(timeThresholdType, timeThresholdObj)
	diags.Append(timeThresholdDiags...)
	if diags.HasError() {
		return diags
	}

	timeThresholdList, timeThresholdListDiags := types.ListValue(
		types.ObjectType{AttrTypes: timeThresholdType},
		[]attr.Value{timeThresholdValue},
	)
	diags.Append(timeThresholdListDiags...)
	if diags.HasError() {
		return diags
	}

	model.TimeThreshold = timeThresholdList

	// Map synthetic test IDs
	syntheticTestIdsSet, syntheticTestIdsDiags := types.SetValueFrom(ctx, types.StringType, apiObject.SyntheticTestIds)
	diags.Append(syntheticTestIdsDiags...)
	if diags.HasError() {
		return diags
	}
	model.SyntheticTestIds = syntheticTestIdsSet

	// Map alert channel IDs
	alertChannelIdsSet, alertChannelIdsDiags := types.SetValueFrom(ctx, types.StringType, apiObject.AlertChannelIds)
	diags.Append(alertChannelIdsDiags...)
	if diags.HasError() {
		return diags
	}
	model.AlertChannelIds = alertChannelIdsSet

	// Map custom payload fields
	customPayloadFieldsList, payloadDiags := tfutils.CustomPayloadFieldsToTerraform(ctx, apiObject.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	return diags
}

func mapTagFilterExpressionFromSchema(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

// BuildCustomPayloadFieldsTyped builds custom payload fields from the Terraform model
func BuildCustomPayloadFieldsTyped(ctx context.Context, customPayloadFields types.List) ([]restapi.CustomPayloadField[any], diag.Diagnostics) {
	var diags diag.Diagnostics
	var result []restapi.CustomPayloadField[any]

	if customPayloadFields.IsNull() || customPayloadFields.IsUnknown() {
		return result, diags
	}

	var elements []types.Object
	diags.Append(customPayloadFields.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	for _, element := range elements {
		var field struct {
			Key          types.String `tfsdk:"key"`
			Value        types.String `tfsdk:"value"`
			DynamicValue types.List   `tfsdk:"dynamic_value"`
		}

		diags.Append(element.As(ctx, &field, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		key := field.Key.ValueString()

		if !field.Value.IsNull() && !field.Value.IsUnknown() {
			// Static string value
			result = append(result, restapi.CustomPayloadField[any]{
				Key:   key,
				Value: field.Value.ValueString(),
			})
		} else if !field.DynamicValue.IsNull() && !field.DynamicValue.IsUnknown() {
			// Dynamic value
			var dynamicElements []types.Object
			diags.Append(field.DynamicValue.ElementsAs(ctx, &dynamicElements, false)...)
			if diags.HasError() {
				return nil, diags
			}

			if len(dynamicElements) > 0 {
				var dynamicValue struct {
					Key     types.String `tfsdk:"key"`
					TagName types.String `tfsdk:"tag_name"`
				}

				diags.Append(dynamicElements[0].As(ctx, &dynamicValue, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}

				dynamicField := map[string]interface{}{
					"tagName": dynamicValue.TagName.ValueString(),
				}

				if !dynamicValue.Key.IsNull() && !dynamicValue.Key.IsUnknown() {
					dynamicField["key"] = dynamicValue.Key.ValueString()
				}

				result = append(result, restapi.CustomPayloadField[any]{
					Key:   key,
					Value: dynamicField,
				})
			}
		}
	}

	return result, diags
}

// Made with Bob
