package syntheticalertconfig

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

// NewSyntheticAlertConfigResourceHandleFramework creates the resource handle for Synthetic Alert Configuration
func NewSyntheticAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.SyntheticAlertConfig] {
	return &syntheticAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSyntheticAlertConfigFramework,
			Schema: schema.Schema{
				Description: SyntheticAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					SyntheticAlertConfigFieldID: schema.StringAttribute{
						Computed:    true,
						Description: SyntheticAlertConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					SyntheticAlertConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: SyntheticAlertConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 256),
						},
					},
					SyntheticAlertConfigFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: SyntheticAlertConfigDescDescription,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 1024),
						},
					},
					SyntheticAlertConfigFieldSyntheticTestIds: schema.SetAttribute{
						Required:    true,
						Description: SyntheticAlertConfigDescSyntheticTestIds,
						ElementType: types.StringType,
					},
					SyntheticAlertConfigFieldSeverity: schema.Int64Attribute{
						Optional:    true,
						Description: SyntheticAlertConfigDescSeverity,
						Validators: []validator.Int64{
							int64validator.OneOf(5, 10),
						},
					},
					SyntheticAlertConfigFieldTagFilter: schema.StringAttribute{
						Optional:    true,
						Description: SyntheticAlertConfigDescTagFilter,
					},
					SyntheticAlertConfigFieldAlertChannelIds: schema.SetAttribute{
						Required:    true,
						Description: SyntheticAlertConfigDescAlertChannelIds,
						ElementType: types.StringType,
					},
					SyntheticAlertConfigFieldGracePeriod: schema.Int64Attribute{
						Optional:    true,
						Description: SyntheticAlertConfigDescGracePeriod,
					},
					SyntheticAlertConfigFieldCustomPayloadField: shared.GetCustomPayloadFieldsSchema(),
					SyntheticAlertConfigFieldRule: schema.SingleNestedAttribute{
						Description: SyntheticAlertConfigDescRule,
						Required:    true,
						Attributes: map[string]schema.Attribute{
							SyntheticAlertRuleFieldAlertType: schema.StringAttribute{
								Required:    true,
								Description: SyntheticAlertConfigDescRuleAlertType,
								Validators: []validator.String{
									stringvalidator.OneOf(SyntheticAlertConfigValidAlertType),
								},
							},
							SyntheticAlertRuleFieldMetricName: schema.StringAttribute{
								Required:    true,
								Description: SyntheticAlertConfigDescRuleMetricName,
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 256),
								},
							},
							SyntheticAlertRuleFieldAggregation: schema.StringAttribute{
								Optional:    true,
								Description: SyntheticAlertConfigDescRuleAggregation,
								Validators: []validator.String{
									stringvalidator.OneOf(AggregationTypeSum, AggregationTypeMean, AggregationTypeMax, AggregationTypeMin, AggregationTypeP25, AggregationTypeP50, AggregationTypeP75, AggregationTypeP90, AggregationTypeP95, AggregationTypeP98, AggregationTypeP99, AggregationTypeP99_9, AggregationTypeP99_99, AggregationTypeDistinctCount, AggregationTypeSumPositive, AggregationTypePerSecond, AggregationTypeIncrease),
								},
							},
						},
					},
					SyntheticAlertConfigFieldTimeThreshold: schema.SingleNestedAttribute{
						Description: SyntheticAlertConfigDescTimeThreshold,
						Required:    true,
						Attributes: map[string]schema.Attribute{
							SyntheticAlertTimeThresholdFieldType: schema.StringAttribute{
								Required:    true,
								Description: SyntheticAlertConfigDescTimeThresholdType,
								Validators: []validator.String{
									stringvalidator.OneOf(SyntheticAlertConfigValidTimeThresholdType),
								},
							},
							SyntheticAlertTimeThresholdFieldViolationsCount: schema.Int64Attribute{
								Required:    true,
								Description: SyntheticAlertConfigDescTimeThresholdViolationsCount,
								Validators: []validator.Int64{
									int64validator.Between(1, 12),
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
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *syntheticAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
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
	if !model.Rule.IsNull() && !model.Rule.IsUnknown() {
		var ruleModel SyntheticAlertRuleModel
		diags.Append(model.Rule.As(ctx, &ruleModel, basetypes.ObjectAsOptions{})...)
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

	// Map time threshold
	var timeThreshold restapi.SyntheticAlertTimeThreshold
	if !model.TimeThreshold.IsNull() && !model.TimeThreshold.IsUnknown() {
		var timeThresholdModel SyntheticAlertTimeThresholdModel
		diags.Append(model.TimeThreshold.As(ctx, &timeThresholdModel, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}

		timeThreshold = restapi.SyntheticAlertTimeThreshold{
			Type:            timeThresholdModel.Type.ValueString(),
			ViolationsCount: int(timeThresholdModel.ViolationsCount.ValueInt64()),
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
				SyntheticAlertConfigErrParsingTagFilter,
				SyntheticAlertConfigErrParsingTagFilterDetail+err.Error(),
			)
			return nil, diags
		}
	} else {
		operator := restapi.LogicalOperatorType(TagFilterLogicalOperatorAnd)
		tagFilter = &restapi.TagFilter{
			Type:            TagFilterTypeExpression,
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		}
	}
	// Map custom payload fields
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() {
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
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

func (r *syntheticAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SyntheticAlertConfig) diag.Diagnostics {
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
				SyntheticAlertConfigErrNormalizingTagFilter,
				SyntheticAlertConfigErrNormalizingTagFilterDetail+err.Error(),
			)
			return diags
		}
		if normalizedTagFilterString != nil {
			model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
		} else {
			model.TagFilter = types.StringNull()
		}

	} else {
		model.TagFilter = types.StringNull()
	}

	// Map rule
	ruleObj := map[string]attr.Value{
		SyntheticAlertRuleFieldAlertType:  types.StringValue(apiObject.Rule.AlertType),
		SyntheticAlertRuleFieldMetricName: types.StringValue(apiObject.Rule.MetricName),
	}

	if apiObject.Rule.Aggregation != "" {
		ruleObj[SyntheticAlertRuleFieldAggregation] = types.StringValue(apiObject.Rule.Aggregation)
	} else {
		ruleObj[SyntheticAlertRuleFieldAggregation] = types.StringNull()
	}

	ruleType := map[string]attr.Type{
		SyntheticAlertRuleFieldAlertType:   types.StringType,
		SyntheticAlertRuleFieldMetricName:  types.StringType,
		SyntheticAlertRuleFieldAggregation: types.StringType,
	}

	ruleValue, ruleDiags := types.ObjectValue(ruleType, ruleObj)
	diags.Append(ruleDiags...)
	if diags.HasError() {
		return diags
	}

	model.Rule = ruleValue

	// Map time threshold
	timeThresholdObj := map[string]attr.Value{
		SyntheticAlertTimeThresholdFieldType:            types.StringValue(apiObject.TimeThreshold.Type),
		SyntheticAlertTimeThresholdFieldViolationsCount: types.Int64Value(int64(apiObject.TimeThreshold.ViolationsCount)),
	}

	timeThresholdType := map[string]attr.Type{
		SyntheticAlertTimeThresholdFieldType:            types.StringType,
		SyntheticAlertTimeThresholdFieldViolationsCount: types.Int64Type,
	}

	timeThresholdValue, timeThresholdDiags := types.ObjectValue(timeThresholdType, timeThresholdObj)
	diags.Append(timeThresholdDiags...)
	if diags.HasError() {
		return diags
	}

	model.TimeThreshold = timeThresholdValue

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
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, apiObject.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
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
