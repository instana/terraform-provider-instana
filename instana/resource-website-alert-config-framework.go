package instana

import (
	"context"
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
						Optional:    true,
						Computed:    true,
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
						Optional:    true,
						Computed:    true,
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
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"slowness": schema.SingleNestedAttribute{
											Description: "Rule based on the slowness of the configured alert configuration target.",
											Optional:    true,
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
										"specific_js_error": schema.SingleNestedAttribute{
											Description: "Rule based on a specific javascript error of the configured alert configuration target.",
											Optional:    true,
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
										"status_code": schema.SingleNestedAttribute{
											Description: "Rule based on the HTTP status code of the configured alert configuration target.",
											Optional:    true,
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
										"throughput": schema.SingleNestedAttribute{
											Description: "Rule based on the throughput of the configured alert configuration target.",
											Optional:    true,
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
								ApplicationAlertConfigFieldThreshold: schema.SingleNestedAttribute{
									Description: "Threshold configuration for different severity levels",
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										LogAlertConfigFieldWarning:  AllThresholdAttributeSchema(),
										LogAlertConfigFieldCritical: AllThresholdAttributeSchema(),
									},
								},
							},
						},
					},
				},
				Blocks: map[string]schema.Block{
					"custom_payload_fields": GetCustomPayloadFieldsSchema(),
					"threshold": schema.SingleNestedBlock{
						Description: "Threshold configuration for different severity levels",
						Blocks: map[string]schema.Block{
							"static": schema.SingleNestedBlock{
								Description: "Static threshold definition.",
								Attributes: map[string]schema.Attribute{
									"operator": schema.StringAttribute{
										Description: "Comparison operator for the static threshold.",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf([]string{">=", ">", "<=", "<", "=="}...),
										},
									},
									"value": schema.Int64Attribute{
										Description: "The numeric value for the static threshold.",
										Optional:    true,
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"adaptive_baseline": schema.SingleNestedBlock{
								Description: "Static threshold definition.",
								Attributes: map[string]schema.Attribute{
									"operator": schema.StringAttribute{
										Description: "Comparison operator for the static threshold.",
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf([]string{">=", ">", "<=", "<", "=="}...),
										},
									},
									"deviation_factor": schema.Float32Attribute{
										Description: "The numeric value for the deviation factor.",
										Optional:    true,
									},
									"adaptability": schema.Float32Attribute{
										Description: "The numeric value for the adaptability.",
										Optional:    true,
									},
									"seasonality": schema.StringAttribute{
										Description: "Value for the seasonality.",
										Optional:    true,
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"historic_baseline": HistoricBaselineBlockSchema(),
						},
					},
					"rule": schema.SingleNestedBlock{
						Description: "Indicates the type of rule this alert configuration is about.",
						Blocks: map[string]schema.Block{
							"slowness": schema.SingleNestedBlock{
								Description: "Rule based on the slowness of the configured alert configuration target.",
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The metric name of the website alert rule.",
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The aggregation function of the website alert rule.",
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"specific_js_error": schema.SingleNestedBlock{
								Description: "Rule based on a specific javascript error of the configured alert configuration target.",
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The metric name of the website alert rule.",
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The aggregation function of the website alert rule.",
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
									"operator": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The operator which will be applied to evaluate this rule.",
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
										},
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The value identify the specific javascript error.",
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"status_code": schema.SingleNestedBlock{
								Description: "Rule based on the HTTP status code of the configured alert configuration target.",
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The metric name of the website alert rule.",
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The aggregation function of the website alert rule.",
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
									"operator": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The operator which will be applied to evaluate this rule.",
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
										},
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The value identify the specific http status code.",
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"throughput": schema.SingleNestedBlock{
								Description: "Rule based on the throughput of the configured alert configuration target.",
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The metric name of the website alert rule.",
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The aggregation function of the website alert rule.",
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"time_threshold": schema.SingleNestedBlock{
						Description: "Indicates the type of violation of the defined threshold.",
						Blocks: map[string]schema.Block{
							"user_impact_of_violations_in_sequence": schema.SingleNestedBlock{
								Description: "Time threshold base on user impact of violations in sequence.",
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The time window if the time threshold.",
									},
									"impact_measurement_method": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: "The impact method of the time threshold based on user impact of violations in sequence.",
										Validators: []validator.String{
											stringvalidator.OneOf("AGGREGATED", "PER_WINDOW"),
										},
									},
									"user_percentage": schema.Float64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The percentage of impacted users of the time threshold based on user impact of violations in sequence.",
									},
									"users": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The number of impacted users of the time threshold based on user impact of violations in sequence.",
									},
								},
							},
							"violations_in_period": schema.SingleNestedBlock{
								Description: "Time threshold base on violations in period.",
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The time window if the time threshold.",
									},
									"violations": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The violations appeared in the period.",
									},
								},
							},
							"violations_in_sequence": schema.SingleNestedBlock{
								Description: "Time threshold base on violations in sequence.",
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "The time window if the time threshold.",
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

	var severity *int
	// Convert severity from Terraform representation to API representation
	if !model.Severity.IsNull() && !model.Severity.IsUnknown() {
		severityVal, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(model.Severity.ValueString())
		severity = &severityVal
		if err != nil {
			diags.AddError("Error converting severity", err.Error())
			return nil, diags
		}
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
	} else {
		operator := restapi.LogicalOperatorType("AND")
		tagFilter = &restapi.TagFilter{
			Type:            "EXPRESSION",
			LogicalOperator: &operator,
			Elements:        []*restapi.TagFilter{},
		}
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
	customPayloadFields := make([]restapi.CustomPayloadField[any], 0)
	// if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
	// 	// Skip custom payload fields for now
	// 	customPayloadFields = []restapi.CustomPayloadField[any]{}
	// }
	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		var payloadDiags diag.Diagnostics
		customPayloadFields, payloadDiags = MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
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

	log.Printf("reached rules section")
	//Map rules
	rules := make([]restapi.WebsiteAlertRuleWithThresholds, 0)
	if !model.Rules.IsNull() && !model.Rules.IsUnknown() {
		var rulesList []RuleWithThresholdPluginModel
		diags.Append(model.Rules.ElementsAs(ctx, &rulesList, false)...)
		if diags.HasError() {
			return nil, diags
		}

		// Process each rule
		for _, i := range rulesList {
			var websiteAlertRule *restapi.WebsiteAlertRule
			if i.Rule != nil {
				log.Printf("reached rule section")

				// Check each rule type and only process if not null/unknown
				if i.Rule.Slowness != nil && !i.Rule.Slowness.MetricName.IsNull() && !i.Rule.Slowness.MetricName.IsUnknown() {
					log.Printf("inside slowness")
					websiteAlertRule, diags = r.mapSlownessRule(ctx, diags, *i.Rule.Slowness)
				} else if i.Rule.SpecificJsError != nil && !i.Rule.SpecificJsError.MetricName.IsNull() && !i.Rule.SpecificJsError.MetricName.IsUnknown() {
					log.Printf("inside SpecificJsError")
					websiteAlertRule, diags = r.mapSpecificJsErrorRule(ctx, diags, *i.Rule.SpecificJsError)
				} else if i.Rule.StatusCode != nil && !i.Rule.StatusCode.MetricName.IsNull() && !i.Rule.StatusCode.MetricName.IsUnknown() {
					log.Printf("inside StatusCode")
					websiteAlertRule, diags = r.mapStatusCodeRule(ctx, diags, *i.Rule.StatusCode)
				} else if i.Rule.Throughput != nil && !i.Rule.Throughput.MetricName.IsNull() && !i.Rule.Throughput.MetricName.IsUnknown() {
					log.Printf("inside Throughput")
					websiteAlertRule, diags = r.mapThroughputRule(ctx, diags, *i.Rule.Throughput)
				}
			}

			// Only process if we have a valid rule
			if websiteAlertRule != nil {
				thresholdMap, thresholdDiags := MapThresholdsAllPluginFromState(ctx, i.Thresholds)
				diags.Append(thresholdDiags...)
				if diags.HasError() {
					return nil, diags
				}

				rules = append(rules, restapi.WebsiteAlertRuleWithThresholds{
					Rule:              websiteAlertRule,
					ThresholdOperator: i.ThresholdOperator.ValueString(),
					Thresholds:        thresholdMap,
				})
			}
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
		Threshold:             threshold,
		TimeThreshold:         *timeThreshold,
		Rules:                 rules,
	}, diags
}

func (r *websiteAlertConfigResourceFramework) mapRuleFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.WebsiteAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if rule is set
	if model.Rule == nil {
		//diags.AddError("Rule is required", "Website alert config rule is required")
		return nil, diags
	}

	ruleModel := model.Rule

	// Check which rule type is set
	if ruleModel.Slowness != nil {
		return r.mapSlownessRule(ctx, diags, *ruleModel.Slowness)
	} else if ruleModel.SpecificJsError != nil {
		return r.mapSpecificJsErrorRule(ctx, diags, *ruleModel.SpecificJsError)
	} else if ruleModel.StatusCode != nil {
		return r.mapStatusCodeRule(ctx, diags, *ruleModel.StatusCode)
	} else if ruleModel.Throughput != nil {
		return r.mapThroughputRule(ctx, diags, *ruleModel.Throughput)
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
	if model.Threshold == nil {
		//diags.AddError("Threshold is required", "Website alert config threshold is required")
		return nil, diags
	}

	threshold := &restapi.Threshold{}

	if model.Threshold.Static != nil {
		value := float64(model.Threshold.Static.Value.ValueInt64())
		threshold.Value = &value
		threshold.Type = "staticThreshold"

		threshold.Operator = restapi.ThresholdOperator(model.Threshold.Static.Operator.ValueString())

	}
	if model.Threshold.AdaptiveBaseline != nil {
		threshold.Type = "adaptiveBaseline"
		threshold.Operator = restapi.ThresholdOperator(model.Threshold.AdaptiveBaseline.Operator.ValueString())
		deviationFactor := model.Threshold.AdaptiveBaseline.DeviationFactor.ValueFloat32()
		threshold.DeviationFactor = &deviationFactor
		adaptability := model.Threshold.AdaptiveBaseline.Adaptability.ValueFloat32()
		threshold.Adaptability = &adaptability
		seasonality := restapi.ThresholdSeasonality(model.Threshold.AdaptiveBaseline.Seasonality.ValueString())
		threshold.Seasonality = &seasonality

	}
	if model.Threshold.HistoricBaseline != nil {
		threshold.Type = "historicBaseline"
		deviationFactor := model.Threshold.HistoricBaseline.Deviation.ValueFloat32()
		threshold.DeviationFactor = &deviationFactor
		seasonality := restapi.ThresholdSeasonality(model.Threshold.HistoricBaseline.Seasonality.ValueString())
		threshold.Seasonality = &seasonality
		baseline, _ := mapBaselineFromState(ctx, model.Threshold.HistoricBaseline.Baseline)
		threshold.Baseline = baseline

	}
	return threshold, diags
}

func (r *websiteAlertConfigResourceFramework) mapTimeThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if time threshold is set
	if model.TimeThreshold == nil {
		diags.AddError("Time threshold is required", "Website alert config time threshold is required")
		return nil, diags
	}
	timeThresholdModel := *model.TimeThreshold

	// Check which time threshold type is set
	if timeThresholdModel.UserImpactOfViolationsInSequence != nil {
		userImpactModel := timeThresholdModel.UserImpactOfViolationsInSequence
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
	} else if timeThresholdModel.ViolationsInPeriod != nil {
		violationsInPeriodModel := timeThresholdModel.ViolationsInPeriod
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
	} else if timeThresholdModel.ViolationsInSequence != nil {
		violationsInSequenceModel := timeThresholdModel.ViolationsInSequence
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
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(*apiObject.Severity)
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
		filterExprStr, err := tagfilter.MapTagFilterToNormalizedString(apiObject.TagFilterExpression)
		if err != nil {
			diags.AddError(
				"Error mapping filter expression",
				fmt.Sprintf("Failed to map filter expression: %s", err),
			)
			return diags
		}
		model.TagFilter = setStringPointerToState(filterExprStr)
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
	model.Rule = r.mapRuleToState(ctx, apiObject.Rule)

	//map rules
	rulesModel := r.mapRulesToState(ctx, apiObject)
	// Define the proper attribute types for RuleWithThresholdPluginModel
	ruleAttrTypes := map[string]attr.Type{
		"rule": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"slowness": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					},
				},
				"specific_js_error": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
						"operator":    types.StringType,
						"value":       types.StringType,
					},
				},
				"status_code": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
						"operator":    types.StringType,
						"value":       types.StringType,
					},
				},
				"throughput": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					},
				},
			},
		},
		"operator": types.StringType,
		"threshold": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"warning": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"static": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"operator": types.StringType,
								"value":    types.Int64Type,
							},
						},
						"adaptive_baseline": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"deviation_factor": types.Float32Type,
								"adaptability":     types.Float32Type,
								"seasonality":      types.StringType,
								"operator":         types.StringType,
							},
						},
						"historic_baseline": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"baseline":         types.ListType{ElemType: types.ListType{ElemType: types.Float64Type}},
								"deviation_factor": types.Float32Type,
								"seasonality":      types.StringType,
							},
						},
					},
				},
				"critical": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"static": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"operator": types.StringType,
								"value":    types.Int64Type,
							},
						},
						"adaptive_baseline": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"deviation_factor": types.Float32Type,
								"adaptability":     types.Float32Type,
								"seasonality":      types.StringType,
								"operator":         types.StringType,
							},
						},
						"historic_baseline": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"baseline":         types.ListType{ElemType: types.ListType{ElemType: types.Float64Type}},
								"deviation_factor": types.Float32Type,
								"seasonality":      types.StringType,
							},
						},
					},
				},
			},
		},
	}

	// Use ListValueFrom to automatically infer types from the slice of structs
	rulesListValue, rulesDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ruleAttrTypes}, rulesModel)
	diags.Append(rulesDiags...)
	if diags.HasError() {
		return diags
	}
	model.Rules = rulesListValue

	//map custom paylaod
	customPayloadFieldsList, payloadDiags := CustomPayloadFieldsToTerraform(ctx, apiObject.CustomerPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return diags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	//map threshold
	model.Threshold = r.mapThresholdToState(ctx, apiObject.Threshold)

	//map time threshold
	model.TimeThreshold = r.mapTimeThresholdToState(apiObject.TimeThreshold)

	log.Printf("website model %v", model)
	// Set state
	diags.Append(state.Set(ctx, &model)...)
	return diags
}

func (r *websiteAlertConfigResourceFramework) mapTimeThresholdToState(timeThreshold restapi.WebsiteTimeThreshold) *WebsiteTimeThresholdModel {
	websiteTimeThresholdModel := WebsiteTimeThresholdModel{}
	switch timeThreshold.Type {
	case "violationsInSequence":
		websiteTimeThresholdModel.ViolationsInSequence = &WebsiteViolationsInSequenceModel{
			TimeWindow: setInt64PointerToState(timeThreshold.TimeWindow),
		}
	case "userImpactOfViolationsInSequence":
		websiteTimeThresholdModel.UserImpactOfViolationsInSequence = &WebsiteUserImpactOfViolationsInSequenceModel{
			TimeWindow:              setInt64PointerToState(timeThreshold.TimeWindow),
			ImpactMeasurementMethod: types.StringValue(string(*timeThreshold.ImpactMeasurementMethod)),
			UserPercentage:          setFloat64PointerToState(timeThreshold.UserPercentage),
			Users:                   setInt64PointerToState(timeThreshold.Users),
		}
	case "violationsInPeriod":
		websiteTimeThresholdModel.ViolationsInPeriod = &WebsiteViolationsInPeriodModel{
			TimeWindow: setInt64PointerToState(timeThreshold.TimeWindow),
			Violations: setInt64PointerToState(timeThreshold.Violations),
		}
	}
	return &websiteTimeThresholdModel
}

func (r *websiteAlertConfigResourceFramework) mapThresholdToState(ctx context.Context, threshold *restapi.Threshold) *WebsiteThresholdModel {
	if threshold == nil {
		return nil
	}

	websiteThresholdModel := WebsiteThresholdModel{
		Static:           nil,
		AdaptiveBaseline: nil,
	}

	thresholdVal := *threshold

	if thresholdVal.Type == "staticThreshold" {
		staticModel := StaticTypeModel{
			Operator: types.StringValue(string(thresholdVal.Operator)),
			Value:    setInt64PointerToState(thresholdVal.Value),
		}
		websiteThresholdModel.Static = &staticModel
	}
	if thresholdVal.Type == "adaptiveBaseline" {
		adaptiveBaselineModel := AdaptiveBaselineModel{
			Operator:        types.StringValue(string(thresholdVal.Operator)),
			DeviationFactor: setFloat32PointerToState(thresholdVal.DeviationFactor),
			Adaptability:    setFloat32PointerToState(thresholdVal.Adaptability),
			Seasonality:     types.StringValue(string(*thresholdVal.Seasonality)),
		}
		websiteThresholdModel.AdaptiveBaseline = &adaptiveBaselineModel
	}
	if thresholdVal.Type == "historicBaseline" {
		historicBaselineModel := HistoricBaselineModel{
			Deviation:   setFloat32PointerToState(thresholdVal.DeviationFactor),
			Seasonality: types.StringValue(string(*thresholdVal.Seasonality)),
		}
		thresholdRules := restapi.ThresholdRule{
			Baseline: thresholdVal.Baseline,
		}
		historicBaselineModel.Baseline, _ = mapBaseline(&thresholdRules)
		websiteThresholdModel.HistoricBaseline = &historicBaselineModel
	}

	return &websiteThresholdModel

}
func (r *websiteAlertConfigResourceFramework) mapRuleToState(ctx context.Context, rule *restapi.WebsiteAlertRule) *WebsiteAlertRuleModel {
	websiteAlertRuleModel := WebsiteAlertRuleModel{
		Slowness:        nil,
		SpecificJsError: nil,
		StatusCode:      nil,
		Throughput:      nil,
	}
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
			Value:       setStringPointerToState(rule.Value),
		}
		websiteAlertRuleModel.SpecificJsError = &websiteAlertRuleConfigModel

	}

	return &websiteAlertRuleModel

}

func (r *websiteAlertConfigResourceFramework) mapRulesToState(ctx context.Context, apiObject *restapi.WebsiteAlertConfig) []RuleWithThresholdPluginModel {
	rules := apiObject.Rules
	var rulesModel []RuleWithThresholdPluginModel
	for _, i := range rules {
		warningThreshold, isWarningThresholdPresent := i.Thresholds[restapi.WarningSeverity]
		criticalThreshold, isCriticalThresholdPresent := i.Thresholds[restapi.CriticalSeverity]

		thresholdPluginModel := ThresholdAllPluginModel{
			Warning:  MapAllThresholdPluginToState(ctx, &warningThreshold, isWarningThresholdPresent),
			Critical: MapAllThresholdPluginToState(ctx, &criticalThreshold, isCriticalThresholdPresent),
		}
		ruleModel := RuleWithThresholdPluginModel{
			Rule:              r.mapRuleToState(ctx, i.Rule),
			ThresholdOperator: types.StringValue(i.ThresholdOperator),
			Thresholds:        &thresholdPluginModel,
		}
		rulesModel = append(rulesModel, ruleModel)
	}
	return rulesModel

}
