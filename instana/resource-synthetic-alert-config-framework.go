package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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

// ResourceInstanaSyntheticAlertConfigFramework the name of the terraform-provider-instana resource to manage synthetic alert configurations
const ResourceInstanaSyntheticAlertConfigFramework = "synthetic_alert_config"

// SyntheticAlertConfigModel represents the data model for the Synthetic Alert Config resource
type SyntheticAlertConfigModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	SyntheticTestIds    types.Set    `tfsdk:"synthetic_test_ids"`
	Severity            types.Int64  `tfsdk:"severity"`
	TagFilter           types.String `tfsdk:"tag_filter"`
	Rule                types.List   `tfsdk:"rule"`
	AlertChannelIds     types.Set    `tfsdk:"alert_channel_ids"`
	TimeThreshold       types.List   `tfsdk:"time_threshold"`
	GracePeriod         types.Int64  `tfsdk:"grace_period"`
	CustomPayloadFields types.List   `tfsdk:"custom_payload_field"`
}

// SyntheticAlertRuleModel represents the rule configuration for synthetic alerts
type SyntheticAlertRuleModel struct {
	AlertType   types.String `tfsdk:"alert_type"`
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
}

// SyntheticAlertTimeThresholdModel represents the time threshold configuration for synthetic alerts
type SyntheticAlertTimeThresholdModel struct {
	Type            types.String `tfsdk:"type"`
	ViolationsCount types.Int64  `tfsdk:"violations_count"`
}

// NewSyntheticAlertConfigResourceHandleFramework creates the resource handle for Synthetic Alert Configuration
func NewSyntheticAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.SyntheticAlertConfig] {
	return &syntheticAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSyntheticAlertConfigFramework,
			Schema: schema.Schema{
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
			},
			SchemaVersion: 1,
		},
	}
}

type syntheticAlertConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *syntheticAlertConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *syntheticAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SyntheticAlertConfig] {
	return api.SyntheticAlertConfigs()
}

func (r *syntheticAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *syntheticAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SyntheticAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SyntheticAlertConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

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

func (r *syntheticAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, apiObject *restapi.SyntheticAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Map basic fields
	model := SyntheticAlertConfigModel{
		ID:          types.StringValue(apiObject.ID),
		Name:        types.StringValue(apiObject.Name),
		Description: types.StringValue(apiObject.Description),
		Severity:    types.Int64Value(int64(apiObject.Severity)),
	}

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

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// Made with Bob
