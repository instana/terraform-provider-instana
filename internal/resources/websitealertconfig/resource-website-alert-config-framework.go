package websitealertconfig

import (
	"context"
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
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

const (
	//WebsiteAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_website_alert_config
	WebsiteAlertConfigFieldAlertChannelIDs = "alert_channel_ids"
	//WebsiteAlertConfigFieldWebsiteID constant value for field websites.website_id of resource instana_website_alert_config
	WebsiteAlertConfigFieldWebsiteID = "website_id"
	//WebsiteAlertConfigFieldDescription constant value for field description of resource instana_website_alert_config
	WebsiteAlertConfigFieldDescription = "description"
	//WebsiteAlertConfigFieldGranularity constant value for field granularity of resource instana_website_alert_config
	WebsiteAlertConfigFieldGranularity = "granularity"
	//WebsiteAlertConfigFieldName constant value for field name of resource instana_website_alert_config
	WebsiteAlertConfigFieldName = "name"
	//WebsiteAlertConfigFieldFullName constant value for field full_name of resource instana_website_alert_config
	WebsiteAlertConfigFieldFullName = "full_name"

	//WebsiteAlertConfigFieldRule constant value for field rule of resource instana_website_alert_config
	WebsiteAlertConfigFieldRule = "rule"
	//WebsiteAlertConfigFieldRuleMetricName constant value for field rule.*.metric_name of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleMetricName = "metric_name"
	//WebsiteAlertConfigFieldRuleAggregation constant value for field rule.*.aggregation of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleAggregation = "aggregation"
	//WebsiteAlertConfigFieldRuleOperator constant value for field rule.*.operator of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleOperator = "operator"
	//WebsiteAlertConfigFieldRuleValue constant value for field rule.*.value of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleValue = "value"
	//WebsiteAlertConfigFieldRuleSlowness constant value for field rule.slowness of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleSlowness = "slowness"
	//WebsiteAlertConfigFieldRuleStatusCode constant value for field rule.status_code of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleStatusCode = "status_code"
	//WebsiteAlertConfigFieldRuleThroughput constant value for field rule.throughput of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleThroughput = "throughput"
	//WebsiteAlertConfigFieldRuleSpecificJsError constant value for field rule.specific_js_error of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleSpecificJsError = "specific_js_error"

	//WebsiteAlertConfigFieldSeverity constant value for field severity of resource instana_website_alert_config
	WebsiteAlertConfigFieldSeverity = "severity"
	//WebsiteAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_website_alert_config
	WebsiteAlertConfigFieldTagFilter = "tag_filter"

	//WebsiteAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThreshold = "time_threshold"
	//WebsiteAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence constant value for field time_threshold.user_impact_of_violations_in_sequence of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence = "user_impact_of_violations_in_sequence"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod constant value for field time_threshold.user_impact_of_violations_in_sequence.impact_measurement_method of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod = "impact_measurement_method"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage constant value for field time_threshold.user_impact_of_violations_in_sequence.user_percentage of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage = "user_percentage"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers constant value for field time_threshold.user_impact_of_violations_in_sequence.users of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers = "users"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	//WebsiteAlertConfigFieldTriggering constant value for field triggering of resource instana_website_alert_config
	WebsiteAlertConfigFieldTriggering = "triggering"
	WebsiteAlertConfigFieldRules      = "rules"
	WebsiteAlertConfigFieldThreshold  = "threshold"
)

// NewWebsiteAlertConfigResourceHandleFramework creates the resource handle for Website Alert Configs
func NewWebsiteAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.WebsiteAlertConfig] {
	return &websiteAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaWebsiteAlertConfigFramework,
			Schema: schema.Schema{
				Description: WebsiteAlertConfigDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: WebsiteAlertConfigDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: WebsiteAlertConfigDescName,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					"description": schema.StringAttribute{
						Required:    true,
						Description: WebsiteAlertConfigDescDescription,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 65536),
						},
					},
					"severity": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: WebsiteAlertConfigDescSeverity,
						Validators: []validator.String{
							stringvalidator.OneOf("warning", "critical"),
						},
					},
					"triggering": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: WebsiteAlertConfigDescTriggering,
						Default:     booldefault.StaticBool(false),
					},
					"website_id": schema.StringAttribute{
						Required:    true,
						Description: WebsiteAlertConfigDescWebsiteID,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 64),
						},
					},
					"tag_filter": schema.StringAttribute{
						Optional:    true,
						Description: WebsiteAlertConfigDescTagFilter,
					},
					"alert_channel_ids": schema.SetAttribute{
						Optional:    true,
						Description: WebsiteAlertConfigDescAlertChannelIDs,
						ElementType: types.StringType,
					},
					"granularity": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: WebsiteAlertConfigDescGranularity,
						Default:     int64default.StaticInt64(600000),
					},
					WebsiteAlertConfigFieldRules: schema.ListNestedAttribute{
						Description: WebsiteAlertConfigDescRules,
						Optional:    true,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"operator": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: WebsiteAlertConfigDescRuleOperator,
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "<", "<="),
									},
								},
								"rule": schema.SingleNestedAttribute{
									Description: WebsiteAlertConfigDescRule,
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"slowness": schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleSlowness,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"metric_name": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												"aggregation": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
													},
												},
											},
										},
										"specific_js_error": schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleSpecificJsError,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"metric_name": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												"aggregation": schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
													},
												},
												"operator": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleOperatorEval,
													Validators: []validator.String{
														stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
													},
												},
												"value": schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleValueJsError,
												},
											},
										},
										"status_code": schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleStatusCode,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"metric_name": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												"aggregation": schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
													},
												},
												"operator": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleOperatorEval,
													Validators: []validator.String{
														stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
													},
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleValueStatusCode,
												},
											},
										},
										"throughput": schema.SingleNestedAttribute{
											Description: WebsiteAlertConfigDescRuleThroughput,
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"metric_name": schema.StringAttribute{
													Required:    true,
													Description: WebsiteAlertConfigDescRuleMetricName,
												},
												"aggregation": schema.StringAttribute{
													Optional:    true,
													Description: WebsiteAlertConfigDescRuleAggregation,
													Validators: []validator.String{
														stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
													},
												},
											},
										},
									},
								},
								WebsiteAlertConfigFieldThreshold: schema.SingleNestedAttribute{
									Description: WebsiteAlertConfigDescThreshold,
									Optional:    true,
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										shared.LogAlertConfigFieldWarning:  shared.AllThresholdAttributeSchema(),
										shared.LogAlertConfigFieldCritical: shared.AllThresholdAttributeSchema(),
									},
								},
							},
						},
					},
					"custom_payload_fields": shared.GetCustomPayloadFieldsSchema(),
				},
				Blocks: map[string]schema.Block{
					"threshold": schema.SingleNestedBlock{
						Description: WebsiteAlertConfigDescThreshold,
						Blocks: map[string]schema.Block{
							"static": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescThresholdStatic,
								Attributes: map[string]schema.Attribute{
									"operator": schema.StringAttribute{
										Description: WebsiteAlertConfigDescThresholdOperator,
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf([]string{">=", ">", "<=", "<", "=="}...),
										},
									},
									"value": schema.Int64Attribute{
										Description: WebsiteAlertConfigDescThresholdValue,
										Optional:    true,
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"adaptive_baseline": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescThresholdAdaptive,
								Attributes: map[string]schema.Attribute{
									"operator": schema.StringAttribute{
										Description: WebsiteAlertConfigDescThresholdOperator,
										Optional:    true,
										Validators: []validator.String{
											stringvalidator.OneOf([]string{">=", ">", "<=", "<", "=="}...),
										},
									},
									"deviation_factor": schema.Float32Attribute{
										Description: WebsiteAlertConfigDescThresholdDeviationFactor,
										Optional:    true,
									},
									"adaptability": schema.Float32Attribute{
										Description: WebsiteAlertConfigDescThresholdAdaptability,
										Optional:    true,
									},
									"seasonality": schema.StringAttribute{
										Description: WebsiteAlertConfigDescThresholdSeasonality,
										Optional:    true,
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"historic_baseline": shared.HistoricBaselineBlockSchema(),
						},
					},
					"rule": schema.SingleNestedBlock{
						Description: WebsiteAlertConfigDescRule,
						Blocks: map[string]schema.Block{
							"slowness": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescRuleSlowness,
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleMetricName,
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleAggregation,
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
								Description: WebsiteAlertConfigDescRuleSpecificJsError,
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleMetricName,
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleAggregation,
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
									"operator": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleOperatorEval,
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
										},
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleValueJsError,
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"status_code": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescRuleStatusCode,
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleMetricName,
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleAggregation,
										Validators: []validator.String{
											stringvalidator.OneOf("SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99"),
										},
									},
									"operator": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleOperatorEval,
										Validators: []validator.String{
											stringvalidator.OneOf("EQUALS", "DOES_NOT_EQUAL", "CONTAINS", "DOES_NOT_CONTAIN"),
										},
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleValueStatusCode,
									},
								},
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
							"throughput": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescRuleThroughput,
								Attributes: map[string]schema.Attribute{
									"metric_name": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleMetricName,
									},
									"aggregation": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescRuleAggregation,
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
						Description: WebsiteAlertConfigDescTimeThreshold,
						Blocks: map[string]schema.Block{
							"user_impact_of_violations_in_sequence": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescTimeThresholdUserImpact,
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdTimeWindow,
									},
									"impact_measurement_method": schema.StringAttribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdImpactMethod,
										Validators: []validator.String{
											stringvalidator.OneOf("AGGREGATED", "PER_WINDOW"),
										},
									},
									"user_percentage": schema.Float64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdUserPercentage,
									},
									"users": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdUsers,
									},
								},
							},
							"violations_in_period": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescTimeThresholdViolationsInPeriod,
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdTimeWindow,
									},
									"violations": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdViolations,
									},
								},
							},
							"violations_in_sequence": schema.SingleNestedBlock{
								Description: WebsiteAlertConfigDescTimeThresholdViolationsInSequence,
								Attributes: map[string]schema.Attribute{
									"time_window": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: WebsiteAlertConfigDescTimeThresholdTimeWindow,
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
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *websiteAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
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
		severityVal, err := util.ConvertSeverityFromTerraformToInstanaAPIRepresentation(model.Severity.ValueString())
		severity = &severityVal
		if err != nil {
			diags.AddError(WebsiteAlertConfigErrConvertSeverity, err.Error())
			return nil, diags
		}
	}

	// Map tag filter expression
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(model.TagFilter.ValueString())
		if err != nil {
			diags.AddError(WebsiteAlertConfigErrParseTagFilter, err.Error())
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
		customPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
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
				thresholdMap, thresholdDiags := shared.MapThresholdsAllPluginFromState(ctx, i.Thresholds)
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

	diags.AddError(WebsiteAlertConfigErrInvalidRuleConfig, WebsiteAlertConfigErrInvalidRuleConfigMsg)
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
		baseline, _ := shared.MapBaselineFromState(ctx, model.Threshold.HistoricBaseline.Baseline)
		threshold.Baseline = baseline

	}
	return threshold, diags
}

func (r *websiteAlertConfigResourceFramework) mapTimeThresholdFromModel(ctx context.Context, model WebsiteAlertConfigModel) (*restapi.WebsiteTimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if time threshold is set
	if model.TimeThreshold == nil {
		diags.AddError(WebsiteAlertConfigErrTimeThresholdRequired, WebsiteAlertConfigErrTimeThresholdRequiredMsg)
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

	diags.AddError(WebsiteAlertConfigErrInvalidTimeThresholdConfig, WebsiteAlertConfigErrInvalidTimeThresholdConfigMsg)
	return nil, diags
}

func (r *websiteAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.WebsiteAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Convert severity from API representation to Terraform representation
	severity, err := util.ConvertSeverityFromInstanaAPIToTerraformRepresentation(*apiObject.Severity)
	if err != nil {
		diags.AddError(WebsiteAlertConfigErrConvertSeverity, err.Error())
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
				WebsiteAlertConfigErrMapFilterExpression,
				fmt.Sprintf(WebsiteAlertConfigErrMapFilterExpressionMsg, err),
			)
			return diags
		}
		model.TagFilter = util.SetStringPointerToState(filterExprStr)
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
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, apiObject.CustomerPayloadFields)
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
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
		}
	case "userImpactOfViolationsInSequence":
		websiteTimeThresholdModel.UserImpactOfViolationsInSequence = &WebsiteUserImpactOfViolationsInSequenceModel{
			TimeWindow:              util.SetInt64PointerToState(timeThreshold.TimeWindow),
			ImpactMeasurementMethod: types.StringValue(string(*timeThreshold.ImpactMeasurementMethod)),
			UserPercentage:          util.SetFloat64PointerToState(timeThreshold.UserPercentage),
			Users:                   util.SetInt64PointerToState(timeThreshold.Users),
		}
	case "violationsInPeriod":
		websiteTimeThresholdModel.ViolationsInPeriod = &WebsiteViolationsInPeriodModel{
			TimeWindow: util.SetInt64PointerToState(timeThreshold.TimeWindow),
			Violations: util.SetInt64PointerToState(timeThreshold.Violations),
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
		staticModel := shared.StaticTypeModel{
			Operator: types.StringValue(string(thresholdVal.Operator)),
			Value:    util.SetInt64PointerToState(thresholdVal.Value),
		}
		websiteThresholdModel.Static = &staticModel
	}
	if thresholdVal.Type == "adaptiveBaseline" {
		adaptiveBaselineModel := shared.AdaptiveBaselineModel{
			Operator:        types.StringValue(string(thresholdVal.Operator)),
			DeviationFactor: util.SetFloat32PointerToState(thresholdVal.DeviationFactor),
			Adaptability:    util.SetFloat32PointerToState(thresholdVal.Adaptability),
			Seasonality:     types.StringValue(string(*thresholdVal.Seasonality)),
		}
		websiteThresholdModel.AdaptiveBaseline = &adaptiveBaselineModel
	}
	if thresholdVal.Type == "historicBaseline" {
		historicBaselineModel := shared.HistoricBaselineModel{
			Deviation:   util.SetFloat32PointerToState(thresholdVal.DeviationFactor),
			Seasonality: types.StringValue(string(*thresholdVal.Seasonality)),
		}
		thresholdRules := restapi.ThresholdRule{
			Baseline: thresholdVal.Baseline,
		}
		historicBaselineModel.Baseline, _ = shared.MapBaseline(&thresholdRules)
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
			Value:       util.SetStringPointerToState(rule.Value),
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

		thresholdPluginModel := shared.ThresholdAllPluginModel{
			Warning:  shared.MapAllThresholdPluginToState(ctx, &warningThreshold, isWarningThresholdPresent),
			Critical: shared.MapAllThresholdPluginToState(ctx, &criticalThreshold, isCriticalThresholdPresent),
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
