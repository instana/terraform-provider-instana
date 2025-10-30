package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResourceInstanaWebsiteAlertConfigFramework the name of the terraform-provider-instana resource to manage website alert configs
const ResourceInstanaWebsiteAlertConfigFramework = "website_alert_config"

// NewWebsiteAlertConfigResourceHandleFramework creates the resource handle for Website Alert Configs
func NewWebsiteAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.WebsiteAlertConfig] {
	return &websiteAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaWebsiteAlertConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages Website Alert Configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the Website Alert Configuration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the Website Alert Configuration.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					"description": schema.StringAttribute{
						Required:    true,
						Description: "The description of the Website Alert Configuration.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 65536),
						},
					},
					"severity": schema.StringAttribute{
						Required:    true,
						Description: "The severity of the alert when triggered.",
						Validators: []validator.String{
							stringvalidator.OneOf("warning", "critical"),
						},
					},
					"triggering": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Flag to indicate whether also an Incident is triggered or not.",
						Default:     booldefault.StaticBool(false),
					},
					"website_id": schema.StringAttribute{
						Required:    true,
						Description: "Unique ID of the website.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 64),
						},
					},
					"tag_filter": schema.StringAttribute{
						Optional:    true,
						Description: "The tag filter expression for the Website Alert Configuration.",
					},
					"alert_channel_ids": schema.SetAttribute{
						Optional:    true,
						Description: "List of IDs of alert channels defined in Instana.",
						ElementType: types.StringType,
					},
					"granularity": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "The evaluation granularity used for detection of violations of the defined threshold.",
						Default:     int64default.StaticInt64(600000),
					},
					ApplicationAlertConfigFieldRules: schema.ListNestedAttribute{
						Description: "A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"operator": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The operator to apply for threshold comparison",
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "<", "<="),
									},
								},
								"rule": schema.SingleNestedAttribute{
									Description: "Indicates the type of rule this alert configuration is about.",
									Attributes: map[string]schema.Attribute{
										"slowness": schema.ListNestedAttribute{
											Description: "Rule based on the slowness of the configured alert configuration target.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"metric_name": schema.StringAttribute{
														Required:    true,
														Description: "The metric name of the website alert rule.",
													},
													"aggregation": schema.StringAttribute{
														Required:    true,
														Description: "The aggregation function of the website alert rule.",
														Validators: []validator.String{
															stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
														},
													},
												},
											},
										},
										"specific_js_error": schema.ListNestedAttribute{
											Description: "Rule based on a specific javascript error of the configured alert configuration target.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"metric_name": schema.StringAttribute{
														Required:    true,
														Description: "The metric name of the website alert rule.",
													},
													"aggregation": schema.StringAttribute{
														Optional:    true,
														Description: "The aggregation function of the website alert rule.",
														Validators: []validator.String{
															stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
														},
													},
													"operator": schema.StringAttribute{
														Required:    true,
														Description: "The operator which will be applied to evaluate this rule.",
														Validators: []validator.String{
															stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
														},
													},
													"value": schema.StringAttribute{
														Optional:    true,
														Description: "The value identify the specific javascript error.",
													},
												},
											},
										},
										"status_code": schema.ListNestedAttribute{
											Description: "Rule based on the HTTP status code of the configured alert configuration target.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"metric_name": schema.StringAttribute{
														Required:    true,
														Description: "The metric name of the website alert rule.",
													},
													"aggregation": schema.StringAttribute{
														Optional:    true,
														Description: "The aggregation function of the website alert rule.",
														Validators: []validator.String{
															stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
														},
													},
													"operator": schema.StringAttribute{
														Required:    true,
														Description: "The operator which will be applied to evaluate this rule.",
														Validators: []validator.String{
															stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
														},
													},
													"value": schema.StringAttribute{
														Required:    true,
														Description: "The value identify the specific http status code.",
													},
												},
											},
										},
										"throughput": schema.ListNestedAttribute{
											Description: "Rule based on the throughput of the configured alert configuration target.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"metric_name": schema.StringAttribute{
														Required:    true,
														Description: "The metric name of the website alert rule.",
													},
													"aggregation": schema.StringAttribute{
														Optional:    true,
														Description: "The aggregation function of the website alert rule.",
														Validators: []validator.String{
															stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
														},
													},
												},
											},
										},
									},
								},
								ApplicationAlertConfigFieldThreshold: schema.SingleNestedAttribute{
									Description: "Threshold configuration for different severity levels",
									Attributes: map[string]schema.Attribute{
										LogAlertConfigFieldWarning:  StaticAndAdaptiveThresholdAttributeSchema(),
										LogAlertConfigFieldCritical: StaticAndAdaptiveThresholdAttributeSchema(),
									},
								},
							},
						},
					},
				},
				Blocks: map[string]schema.Block{
					"custom_payload_fields": GetCustomPayloadFieldsSchema(),
					"threshold": schema.SingleNestedBlock{
						Description: "The threshold configuration for the Website Alert Configuration.",
						Attributes: map[string]schema.Attribute{
							"operator": schema.StringAttribute{
								Required:    true,
								Description: "The operator of the threshold.",
								Validators: []validator.String{
									stringvalidator.OneOf(">", ">=", "<", "<=", "=="),
								},
							},
							"value": schema.Float64Attribute{
								Required:    true,
								Description: "The value of the threshold.",
							},
						},
					},
					"rule": schema.SingleNestedBlock{
						Description: "Indicates the type of rule this alert configuration is about.",
						Blocks: map[string]schema.Block{
							"slowness": schema.SingleNestedBlock{
								Description: "Rule based on the slowness of the configured alert configuration target.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"metric_name": schema.StringAttribute{
											Required:    true,
											Description: "The metric name of the website alert rule.",
										},
										"aggregation": schema.StringAttribute{
											Required:    true,
											Description: "The aggregation function of the website alert rule.",
											Validators: []validator.String{
												stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
											},
										},
									},
								},
							},
							"specific_js_error": schema.SingleNestedBlock{
								Description: "Rule based on a specific javascript error of the configured alert configuration target.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"metric_name": schema.StringAttribute{
											Required:    true,
											Description: "The metric name of the website alert rule.",
										},
										"aggregation": schema.StringAttribute{
											Optional:    true,
											Description: "The aggregation function of the website alert rule.",
											Validators: []validator.String{
												stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
											},
										},
										"operator": schema.StringAttribute{
											Required:    true,
											Description: "The operator which will be applied to evaluate this rule.",
											Validators: []validator.String{
												stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
											},
										},
										"value": schema.StringAttribute{
											Optional:    true,
											Description: "The value identify the specific javascript error.",
										},
									},
								},
							},
							"status_code": schema.SingleNestedBlock{
								Description: "Rule based on the HTTP status code of the configured alert configuration target.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"metric_name": schema.StringAttribute{
											Required:    true,
											Description: "The metric name of the website alert rule.",
										},
										"aggregation": schema.StringAttribute{
											Optional:    true,
											Description: "The aggregation function of the website alert rule.",
											Validators: []validator.String{
												stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
											},
										},
										"operator": schema.StringAttribute{
											Required:    true,
											Description: "The operator which will be applied to evaluate this rule.",
											Validators: []validator.String{
												stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
											},
										},
										"value": schema.StringAttribute{
											Required:    true,
											Description: "The value identify the specific http status code.",
										},
									},
								},
							},
							"throughput": schema.SingleNestedBlock{
								Description: "Rule based on the throughput of the configured alert configuration target.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"metric_name": schema.StringAttribute{
											Required:    true,
											Description: "The metric name of the website alert rule.",
										},
										"aggregation": schema.StringAttribute{
											Optional:    true,
											Description: "The aggregation function of the website alert rule.",
											Validators: []validator.String{
												stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
											},
										},
									},
								},
							},
						},
					},
					"time_threshold": schema.ListNestedBlock{
						Description: "Indicates the type of violation of the defined threshold.",
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"user_impact_of_violations_in_sequence": schema.ListNestedBlock{
									Description: "Time threshold base on user impact of violations in sequence.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"time_window": schema.Int64Attribute{
												Optional:    true,
												Description: "The time window if the time threshold.",
											},
											"impact_measurement_method": schema.StringAttribute{
												Required:    true,
												Description: "The impact method of the time threshold based on user impact of violations in sequence.",
												Validators: []validator.String{
													stringvalidator.OneOf("AGGREGATED", "PER_WINDOW"),
												},
											},
											"user_percentage": schema.Float64Attribute{
												Optional:    true,
												Description: "The percentage of impacted users of the time threshold based on user impact of violations in sequence.",
											},
											"users": schema.Int64Attribute{
												Optional:    true,
												Description: "The number of impacted users of the time threshold based on user impact of violations in sequence.",
											},
										},
									},
								},
								"violations_in_period": schema.ListNestedBlock{
									Description: "Time threshold base on violations in period.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"time_window": schema.Int64Attribute{
												Optional:    true,
												Description: "The time window if the time threshold.",
											},
											"violations": schema.Int64Attribute{
												Optional:    true,
												Description: "The violations appeared in the period.",
											},
										},
									},
								},
								"violations_in_sequence": schema.ListNestedBlock{
									Description: "Time threshold base on violations in sequence.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"time_window": schema.Int64Attribute{
												Optional:    true,
												Description: "The time window if the time threshold.",
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

type websiteAlertConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *websiteAlertConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *websiteAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return api.WebsiteAlertConfig()
}

func (r *websiteAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *websiteAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.WebsiteAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model WebsiteAlertConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Convert severity from Terraform representation to API representation
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(model.Severity.ValueString())
	if err != nil {
		diags.AddError("Error converting severity", err.Error())
		return nil, diags
	}

	// Map tag filter expression
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError("Error parsing tag filter", err.Error())
			return nil, diags
		}
		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
	}

	// Map alert channel IDs
	var alertChannelIDs []string
	if !model.AlertChannelIDs.IsNull() && !model.AlertChannelIDs.IsUnknown() {
		diags.Append(model.AlertChannelIDs.ElementsAs(ctx, &alertChannelIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map custom payload fields
	var customPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		// Skip custom payload fields for now
		customPayloadFields = []restapi.CustomPayloadField[any]{}
	}

	// Map rule
	rule, ruleDiags := r.mapRuleFromModel(ctx, model)
	if ruleDiags.HasError() {
		diags.Append(ruleDiags...)
		return nil, diags
	}

	// Map threshold
	threshold, thresholdDiags := r.mapThresholdFromModel(ctx, model)
	if thresholdDiags.HasError() {
		diags.Append(thresholdDiags...)
		return nil, diags
	}

	// Map time threshold
	timeThreshold, timeThresholdDiags := r.mapTimeThresholdFromModel(ctx, model)
	if timeThresholdDiags.HasError() {
		diags.Append(timeThresholdDiags...)
		return nil, diags
	}

	//Map rules
	rules := make([]restapi.WebsiteAlertRuleWithThresholds, 0)
	if len(model.Rules) != 0 {

		// Skip custom payload fields for now
		for _, i := range model.Rules {
			var websiteAlertRule *restapi.WebsiteAlertRule
			if i.Rule.Slowness != nil {
				websiteAlertRule, diags = mapSlownessRule(ctx, diags, *i.Rule.Slowness)
			} else if i.Rule.SpecificJsError != nil {
				websiteAlertRule, diags = mapSpecificJsErrorRule(ctx, diags, *i.Rule.SpecificJsError)
			} else if i.Rule.StatusCode != nil {
				websiteAlertRule, diags = mapStatusCodeRule(ctx, diags, *i.Rule.StatusCode)
			} else if i.Rule.Throughput != nil {
				websiteAlertRule, diags = mapThroughputRule(ctx, diags, *i.Rule.Throughput)
			}
			var thresholdMap map[restapi.AlertSeverity]restapi.ThresholdRule
			warningThreshold, warningDiags := MapThresholdRulePluginFromState(ctx, i.Thresholds.Warning)
			diags.Append(warningDiags...)
			if diags.HasError() {
				return nil, diags
			}
			criticalThreshold, criticalDiags := MapThresholdRulePluginFromState(ctx, i.Thresholds.Critical)
			diags.Append(criticalDiags...)
			if diags.HasError() {
				return nil, diags
			}
			thresholdMap[restapi.WarningSeverity] = *warningThreshold
			thresholdMap[restapi.CriticalSeverity] = *criticalThreshold

			rules = append(rules, restapi.WebsiteAlertRuleWithThresholds{
				Rule:              websiteAlertRule,
				ThresholdOperator: i.ThresholdOperator.ValueString(),
				Thresholds:        thresholdMap,
			})
		}

	}
	// Create API object
	return &restapi.WebsiteAlertConfig{
		ID:                    model.ID.ValueString(),
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		Severity:              severity,
		Triggering:            model.Triggering.ValueBool(),
		WebsiteID:             model.WebsiteID.ValueString(),
		TagFilterExpression:   tagFilter,
		AlertChannelIDs:       alertChannelIDs,
		Granularity:           restapi.Granularity(model.Granularity.ValueInt64()),
		CustomerPayloadFields: customPayloadFields,
		Rule:                  rule,
		Threshold:             *threshold,
		TimeThreshold:         *timeThreshold,
		Rules:                 rules,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapRuleFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if rule is set
	if model.Rule.IsNull() || model.Rule.IsUnknown() {
		diags.AddError("Rule is required", "Website alert config rule is required")
		return nil, diags
	}

	var ruleModels []WebsiteAlertRuleModel
	diags.Append(model.Rule.ElementsAs(ctx, &ruleModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(ruleModels) != 1 {
		diags.AddError("Invalid rule configuration", "Exactly one rule configuration is required")
		return nil, diags
	}

	ruleModel := ruleModels[0]

	// Check which rule type is set
	if !ruleModel.Slowness.IsNull() && !ruleModel.Slowness.IsUnknown() {
		var slownessModels []WebsiteAlertRuleConfigModel
		diags.Append(ruleModel.Slowness.ElementsAs(ctx, &slownessModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		return r.mapSlownessRule(ctx, diags, slownessModels[0])
	} else if !ruleModel.SpecificJsError.IsNull() && !ruleModel.SpecificJsError.IsUnknown() {
		var specificJsErrorModels []WebsiteAlertRuleConfigCompleteModel
		diags.Append(ruleModel.SpecificJsError.ElementsAs(ctx, &specificJsErrorModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		return r.mapSpecificJsErrorRule(ctx, diags, specificJsErrorModels[0])
	} else if !ruleModel.StatusCode.IsNull() && !ruleModel.StatusCode.IsUnknown() {
		var statusCodeModels []WebsiteAlertRuleConfigCompleteModel
		diags.Append(ruleModel.StatusCode.ElementsAs(ctx, &statusCodeModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		return r.mapStatusCodeRule(ctx, diags, statusCodeModels[0])
	} else if !ruleModel.Throughput.IsNull() && !ruleModel.Throughput.IsUnknown() {
		var throughputModels []WebsiteAlertRuleConfigModel
		diags.Append(ruleModel.Throughput.ElementsAs(ctx, &throughputModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		return r.mapThroughputRule(ctx, diags, throughputModels[0])
	}

	diags.AddError("Invalid rule configuration", "Exactly one rule type configuration is required")
	return nil, diags
}

func (r *websiteAlertConfigResourceFramework) mapThroughputRule(ctx context.Context, diags diag.Diagnostics, throughputModel WebsiteAlertRuleConfigModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {
	var aggregationPtr *restapi.Aggregation
	if !throughputModel.Aggregation.IsNull() && !throughputModel.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(throughputModel.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	return &restapi.WebsiteAlertRule{
		AlertType:   "throughput",
		MetricName:  throughputModel.MetricName.ValueString(),
		Aggregation: aggregationPtr,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapStatusCodeRule(ctx context.Context, diags diag.Diagnostics, statusCodeModel WebsiteAlertRuleConfigCompleteModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {

	var aggregationPtr *restapi.Aggregation
	if !statusCodeModel.Aggregation.IsNull() && !statusCodeModel.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(statusCodeModel.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	operator := restapi.ExpressionOperator(statusCodeModel.Operator.ValueString())
	value := statusCodeModel.Value.ValueString()

	return &restapi.WebsiteAlertRule{
		AlertType:   "statusCode",
		MetricName:  statusCodeModel.MetricName.ValueString(),
		Aggregation: aggregationPtr,
		Operator:    &operator,
		Value:       &value,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapSpecificJsErrorRule(ctx context.Context, diags diag.Diagnostics, specificJsErrorModel WebsiteAlertRuleConfigCompleteModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {

	var aggregationPtr *restapi.Aggregation
	if !specificJsErrorModel.Aggregation.IsNull() && !specificJsErrorModel.Aggregation.IsUnknown() {
		aggregation := restapi.Aggregation(specificJsErrorModel.Aggregation.ValueString())
		aggregationPtr = &aggregation
	}

	operator := restapi.ExpressionOperator(specificJsErrorModel.Operator.ValueString())
	var valuePtr *string
	if !specificJsErrorModel.Value.IsNull() && !specificJsErrorModel.Value.IsUnknown() {
		value := specificJsErrorModel.Value.ValueString()
		valuePtr = &value
	}

	return &restapi.WebsiteAlertRule{
		AlertType:   "specificJsError",
		MetricName:  specificJsErrorModel.MetricName.ValueString(),
		Aggregation: aggregationPtr,
		Operator:    &operator,
		Value:       valuePtr,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapSlownessRule(ctx context.Context, diags diag.Diagnostics, slownessModel WebsiteAlertRuleConfigModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {
	aggregation := restapi.Aggregation(slownessModel.Aggregation.ValueString())

	return &restapi.WebsiteAlertRule{
		AlertType:   "slowness",
		MetricName:  slownessModel.MetricName.ValueString(),
		Aggregation: &aggregation,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.Threshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if threshold is set
	if model.Threshold.IsNull() || model.Threshold.IsUnknown() {
		diags.AddError("Threshold is required", "Website alert config threshold is required")
		return nil, diags
	}

	var thresholdModel WebsiteThresholdModel
	diags.Append(model.Threshold.As(ctx, &thresholdModel, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	value := thresholdModel.Value.ValueFloat64()
	valuePtr := &value
	return &restapi.Threshold{
		Operator: restapi.ThresholdOperator(thresholdModel.Operator.ValueString()),
		Value:    valuePtr,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapTimeThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if time threshold is set
	if model.TimeThreshold.IsNull() || model.TimeThreshold.IsUnknown() {
		diags.AddError("Time threshold is required", "Website alert config time threshold is required")
		return nil, diags
	}

	var timeThresholdModels []WebsiteTimeThresholdModel
	diags.Append(model.TimeThreshold.ElementsAs(ctx, &timeThresholdModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(timeThresholdModels) != 1 {
		diags.AddError("Invalid time threshold", "Exactly one time threshold configuration is required")
		return nil, diags
	}

	timeThresholdModel := timeThresholdModels[0]

	// Check which time threshold type is set
	if !timeThresholdModel.UserImpactOfViolationsInSequence.IsNull() && !timeThresholdModel.UserImpactOfViolationsInSequence.IsUnknown() {
		var userImpactModels []WebsiteUserImpactOfViolationsInSequenceModel
		diags.Append(timeThresholdModel.UserImpactOfViolationsInSequence.ElementsAs(ctx, &userImpactModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(userImpactModels) != 1 {
			diags.AddError("Invalid user impact configuration", "Exactly one user impact of violations in sequence configuration is required")
			return nil, diags
		}

		userImpactModel := userImpactModels[0]
		var timeWindowPtr *int64
		if !userImpactModel.TimeWindow.IsNull() && !userImpactModel.TimeWindow.IsUnknown() {
			timeWindow := userImpactModel.TimeWindow.ValueInt64()
			timeWindowPtr = &timeWindow
		}

		impactMeasurementMethod := restapi.WebsiteImpactMeasurementMethod(userImpactModel.ImpactMeasurementMethod.ValueString())
		var userPercentagePtr *float64
		if !userImpactModel.UserPercentage.IsNull() && !userImpactModel.UserPercentage.IsUnknown() {
			userPercentage := userImpactModel.UserPercentage.ValueFloat64()
			userPercentagePtr = &userPercentage
		}

		var usersPtr *int32
		if !userImpactModel.Users.IsNull() && !userImpactModel.Users.IsUnknown() {
			users := int32(userImpactModel.Users.ValueInt64())
			usersPtr = &users
		}

		return &restapi.WebsiteTimeThreshold{
			Type:                    "userImpactOfViolationsInSequence",
			TimeWindow:              timeWindowPtr,
			ImpactMeasurementMethod: &impactMeasurementMethod,
			UserPercentage:          userPercentagePtr,
			Users:                   usersPtr,
		}, diags
	} else if !timeThresholdModel.ViolationsInPeriod.IsNull() && !timeThresholdModel.ViolationsInPeriod.IsUnknown() {
		var violationsInPeriodModels []WebsiteViolationsInPeriodModel
		diags.Append(timeThresholdModel.ViolationsInPeriod.ElementsAs(ctx, &violationsInPeriodModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(violationsInPeriodModels) != 1 {
			diags.AddError("Invalid violations in period", "Exactly one violations in period configuration is required")
			return nil, diags
		}

		violationsInPeriodModel := violationsInPeriodModels[0]
		var timeWindowPtr *int64
		if !violationsInPeriodModel.TimeWindow.IsNull() && !violationsInPeriodModel.TimeWindow.IsUnknown() {
			timeWindow := violationsInPeriodModel.TimeWindow.ValueInt64()
			timeWindowPtr = &timeWindow
		}

		var violationsPtr *int32
		if !violationsInPeriodModel.Violations.IsNull() && !violationsInPeriodModel.Violations.IsUnknown() {
			violations := int32(violationsInPeriodModel.Violations.ValueInt64())
			violationsPtr = &violations
		}

		return &restapi.WebsiteTimeThreshold{
			Type:       "violationsInPeriod",
			TimeWindow: timeWindowPtr,
			Violations: violationsPtr,
		}, diags
	} else if !timeThresholdModel.ViolationsInSequence.IsNull() && !timeThresholdModel.ViolationsInSequence.IsUnknown() {
		var violationsInSequenceModels []WebsiteViolationsInSequenceModel
		diags.Append(timeThresholdModel.ViolationsInSequence.ElementsAs(ctx, &violationsInSequenceModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(violationsInSequenceModels) != 1 {
			diags.AddError("Invalid violations in sequence", "Exactly one violations in sequence configuration is required")
			return nil, diags
		}

		violationsInSequenceModel := violationsInSequenceModels[0]
		var timeWindowPtr *int64
		if !violationsInSequenceModel.TimeWindow.IsNull() && !violationsInSequenceModel.TimeWindow.IsUnknown() {
			timeWindow := violationsInSequenceModel.TimeWindow.ValueInt64()
			timeWindowPtr = &timeWindow
		}

		return &restapi.WebsiteTimeThreshold{
			Type:       "violationsInSequence",
			TimeWindow: timeWindowPtr,
		}, diags
	}

	diags.AddError("Invalid time threshold configuration", "Exactly one time threshold type configuration is required")
	return nil, diags
}

func (r *websiteAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.WebsiteAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Convert severity from API representation to Terraform representation
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(apiObject.Severity)
	if err != nil {
		diags.AddError("Error converting severity", err.Error())
		return diags
	}

	// Create a model and populate it with values from the config
	model := WebsiteAlertConfigModel{
		ID:          types.StringValue(apiObject.ID),
		Name:        types.StringValue(apiObject.Name),
		Description: types.StringValue(apiObject.Description),
		Severity:    types.StringValue(severity),
		Triggering:  types.BoolValue(apiObject.Triggering),
		WebsiteID:   types.StringValue(apiObject.WebsiteID),
		Granularity: types.Int64Value(int64(apiObject.Granularity)),
	}

	// Map tag filter expression
	if apiObject.TagFilterExpression != nil {
		model.TagFilter = types.StringValue("tag filter expression")
	} else {
		model.TagFilter = types.StringNull()
	}

	// Map alert channel IDs
	if apiObject.AlertChannelIDs != nil && len(apiObject.AlertChannelIDs) > 0 {
		alertChannelIDs := make([]attr.Value, len(apiObject.AlertChannelIDs))
		for i, id := range apiObject.AlertChannelIDs {
			alertChannelIDs[i] = types.StringValue(id)
		}
		model.AlertChannelIDs = types.SetValueMust(types.StringType, alertChannelIDs)
	} else {
		model.AlertChannelIDs = types.SetNull(types.StringType)
	}

	//map rule
	model.Rule = r.mapRuleToState(ctx, apiObject)

	//map rules

	// Set state
	diags.Append(state.Set(ctx, &model)...)
	return diags
}

func (r *websiteAlertConfigResourceFramework) mapRuleToState(ctx context.Context, apiObject *restapi.WebsiteAlertConfig) WebsiteAlertRuleModel {

	rule := apiObject.Rule
	websiteAlertRuleModel := WebsiteAlertRuleModel{}
	switch rule.AlertType {
	case "throughput":
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
		}
		websiteAlertRuleModel.Throughput = &websiteAlertRuleConfigModel
	case "slowness":
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
		}
		websiteAlertRuleModel.Slowness = &websiteAlertRuleConfigModel
	case "statusCode":
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigCompleteModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
		}
		websiteAlertRuleModel.StatusCode = &websiteAlertRuleConfigModel
	case "specificJsError":
		websiteAlertRuleConfigModel := WebsiteAlertRuleConfigCompleteModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(*rule.Aggregation)),
			Operator:    types.StringValue(string(*rule.Operator)),
			Value:       types.StringValue(*rule.Value),
		}
		websiteAlertRuleModel.SpecificJsError = &websiteAlertRuleConfigModel

	}

	return websiteAlertRuleModel

}
