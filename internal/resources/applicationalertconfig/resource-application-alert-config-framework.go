package applicationalertconfig

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewApplicationAlertConfigResourceHandleFramework creates the resource handle for Application Alert Configuration
func NewApplicationAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.ApplicationAlertConfig] {
	return &applicationAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaApplicationAlertConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages application alert configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the application alert configuration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					ApplicationAlertConfigFieldAlertChannelIDs: schema.SetAttribute{
						Optional:    true,
						Computed:    true,
						Description: "List of IDs of alert channels defined in Instana. Deprecated: Use alert_channels instead.",
						ElementType: types.StringType,
					},
					ApplicationAlertConfigFieldAlertChannels: schema.MapAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Set of alert channel IDs associated with the severity.",
						ElementType: types.SetType{
							ElemType: types.StringType,
						},
					},
					ApplicationAlertConfigFieldBoundaryScope: schema.StringAttribute{
						Required:    true,
						Description: "The boundary scope of the application alert config",
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedApplicationAlertConfigBoundaryScopes.ToStringSlice()...),
						},
					},
					ApplicationAlertConfigFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: "The description text of the application alert config",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 65536),
						},
					},
					ApplicationAlertConfigFieldEvaluationType: schema.StringAttribute{
						Required:    true,
						Description: "The evaluation type of the application alert config",
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedApplicationAlertEvaluationTypes.ToStringSlice()...),
						},
					},
					ApplicationAlertConfigFieldGracePeriod: schema.Int64Attribute{
						Optional:    true,
						Description: "The duration for which an alert remains open after conditions are no longer violated, with the alert auto-closing once the grace period expires.",
					},
					ApplicationAlertConfigFieldGranularity: schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(int64(restapi.Granularity600000)),
						Description: "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
					},
					ApplicationAlertConfigFieldIncludeInternal: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Optional flag to indicate whether also internal calls are included in the scope or not. The default is false",
					},
					ApplicationAlertConfigFieldIncludeSynthetic: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not. The default is false",
					},
					ApplicationAlertConfigFieldName: schema.StringAttribute{
						Required:    true,
						Description: "Name for the application alert configuration",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					ApplicationAlertConfigFieldSeverity: schema.StringAttribute{
						Optional:    true,
						Description: "The severity of the alert when triggered. Deprecated: Use rules with thresholds instead.",
						Validators: []validator.String{
							stringvalidator.OneOf(restapi.SupportedSeverities.TerraformRepresentations()...),
						},
					},
					ApplicationAlertConfigFieldTagFilter: schema.StringAttribute{
						Optional:    true,
						Description: "The tag filter expression for the application alert config",
					},
					ApplicationAlertConfigFieldTriggering: schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Optional flag to indicate whether also an Incident is triggered or not. The default is false",
					},
					shared.DefaultCustomPayloadFieldsName: shared.GetCustomPayloadFieldsSchema(),
					ApplicationAlertConfigFieldRules: schema.ListNestedAttribute{
						Description: "A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								ApplicationAlertConfigFieldThresholdOperator: schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The operator to apply for threshold comparison",
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "<", "<="),
									},
								},
								ApplicationAlertConfigFieldThreshold: schema.SingleNestedAttribute{
									Description: "Threshold configuration for different severity levels",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										shared.LogAlertConfigFieldWarning:  shared.StaticAndAdaptiveThresholdBlockSchema(),
										shared.LogAlertConfigFieldCritical: shared.StaticAndAdaptiveThresholdBlockSchema(),
									},
								},
								ApplicationAlertConfigFieldRule: schema.SingleNestedAttribute{
									Description: "The rule configuration",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										ApplicationAlertConfigFieldRuleErrorRate: schema.SingleNestedAttribute{
											Description: "Rule based on the error rate of the configured alert configuration target",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The metric name of the application alert rule",
												},
												ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The aggregation function of the application alert rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
													},
												},
											},
										},
										ApplicationAlertConfigFieldRuleErrors: schema.SingleNestedAttribute{
											Description: "Rule based on the number of errors of the configured alert configuration target",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The metric name of the application alert rule",
												},
												ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The aggregation function of the application alert rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
													},
												},
											},
										},
										ApplicationAlertConfigFieldRuleLogs: schema.SingleNestedAttribute{
											Description: "Rule based on logs of the configured alert configuration target",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The metric name of the application alert rule",
												},
												ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The aggregation function of the application alert rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
													},
												},
												ApplicationAlertConfigFieldRuleLogsLevel: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The log level for which this rule applies to",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedLogLevels.ToStringSlice()...),
													},
												},
												ApplicationAlertConfigFieldRuleLogsMessage: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The log message for which this rule applies to",
												},
												ApplicationAlertConfigFieldRuleLogsOperator: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The operator which will be applied to evaluate this rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedExpressionOperators.ToStringSlice()...),
													},
												},
											},
										},
										ApplicationAlertConfigFieldRuleSlowness: schema.SingleNestedAttribute{
											Description: "Rule based on the slowness of the configured alert configuration target",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The metric name of the application alert rule",
												},
												ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The aggregation function of the application alert rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
													},
												},
											},
										},
										ApplicationAlertConfigFieldRuleStatusCode: schema.SingleNestedAttribute{
											Description: "Rule based on the HTTP status code of the configured alert configuration target",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The metric name of the application alert rule",
												},
												ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The aggregation function of the application alert rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
													},
												},
												ApplicationAlertConfigFieldRuleStatusCodeStart: schema.Int64Attribute{
													Optional:    true,
													Computed:    true,
													Description: "minimal HTTP status code applied for this rule",
												},
												ApplicationAlertConfigFieldRuleStatusCodeEnd: schema.Int64Attribute{
													Optional:    true,
													Computed:    true,
													Description: "maximum HTTP status code applied for this rule",
												},
											},
										},
										ApplicationAlertConfigFieldRuleThroughput: schema.SingleNestedAttribute{
											Description: "Rule based on the throughput of the configured alert configuration target",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The metric name of the application alert rule",
												},
												ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
													Optional:    true,
													Computed:    true,
													Description: "The aggregation function of the application alert rule",
													Validators: []validator.String{
														stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
													},
												},
											},
										},
									},
								},
							},
						},
					},
					ApplicationAlertConfigFieldTimeThreshold: schema.SingleNestedAttribute{
						Description: "Indicates the type of violation of the defined threshold.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							ApplicationAlertConfigFieldTimeThresholdRequestImpact: schema.SingleNestedAttribute{
								Description: "Time threshold base on request impact",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									ApplicationAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Required:    true,
										Description: "The time window if the time threshold",
									},
									ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests: schema.Int64Attribute{
										Required:    true,
										Description: "The number of requests in the given window",
									},
								},
							},
							ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: schema.SingleNestedAttribute{
								Description: "Time threshold base on violations in period",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									ApplicationAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Required:    true,
										Description: "The time window if the time threshold",
									},
									ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations: schema.Int64Attribute{
										Required:    true,
										Description: "The violations appeared in the period",
										Validators:  []validator.Int64{
											// TODO: Add validator for range 1-12
										},
									},
								},
							},
							ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: schema.SingleNestedAttribute{
								Description: "Time threshold base on violations in sequence",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									ApplicationAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
										Required:    true,
										Description: "The time window if the time threshold",
									},
								},
							},
						},
					},
					ApplicationAlertConfigFieldApplications: schema.SetNestedAttribute{
						Description: "Selection of applications in scope.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								ApplicationAlertConfigFieldApplicationsApplicationID: schema.StringAttribute{
									Required:    true,
									Description: "ID of the included application",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 64),
									},
								},
								ApplicationAlertConfigFieldApplicationsInclusive: schema.BoolAttribute{
									Required:    true,
									Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
								},
								ApplicationAlertConfigFieldApplicationsServices: schema.SetNestedAttribute{
									Description: "Selection of services in scope.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldApplicationsServicesServiceID: schema.StringAttribute{
												Required:    true,
												Description: "ID of the included service",
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 64),
												},
											},
											ApplicationAlertConfigFieldApplicationsInclusive: schema.BoolAttribute{
												Required:    true,
												Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
											},
											ApplicationAlertConfigFieldApplicationsServicesEndpoints: schema.SetNestedAttribute{
												Description: "Selection of endpoints in scope.",
												Optional:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID: schema.StringAttribute{
															Required:    true,
															Description: "ID of the included endpoint",
															Validators: []validator.String{
																stringvalidator.LengthBetween(0, 64),
															},
														},
														ApplicationAlertConfigFieldApplicationsInclusive: schema.BoolAttribute{
															Required:    true,
															Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			SkipIDGeneration: true,
			SchemaVersion:    1,
		},
	}
}

// NewGlobalApplicationAlertConfigResourceHandleFramework creates the resource handle for Global Application Alert Configuration
func NewGlobalApplicationAlertConfigResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.ApplicationAlertConfig] {
	return &applicationAlertConfigResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaGlobalApplicationAlertConfigFramework,
			Schema:        NewApplicationAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
		isGlobal: true,
	}
}

type applicationAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
	isGlobal bool
}

func (r *applicationAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *applicationAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig] {
	if r.isGlobal {
		return api.GlobalApplicationAlertConfigs()
	}
	return api.ApplicationAlertConfigs()
}

func (r *applicationAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.ApplicationAlertConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model ApplicationAlertConfigModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Delegate to the resource implementation
	resource, resourceDiags := r.NewResource(ctx, nil)
	if resourceDiags.HasError() {
		diags.Append(resourceDiags...)
		return nil, diags
	}

	// Create a new state with the model
	var tempState tfsdk.State
	if plan != nil {
		tempState.Schema = plan.Schema
	} else if state != nil {
		tempState.Schema = state.Schema
	}

	// Set the model in the state
	diags.Append(tempState.Set(ctx, model)...)
	if diags.HasError() {
		return nil, diags
	}

	// Map the state to the data object
	return resource.MapStateToDataObject(ctx, tempState)
}

func (r *applicationAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, obj *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Delegate to the resource implementation
	resource, resourceDiags := r.NewResource(ctx, nil)
	if resourceDiags.HasError() {
		diags.Append(resourceDiags...)
		return diags
	}

	// Update the state with the object
	return resource.UpdateState(ctx, state, nil, obj)
}

func (r *applicationAlertConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	// No computed fields to set
	return nil
}

// ResourceFramework interface for application alert config resources
type ResourceFramework[T restapi.InstanaDataObject] interface {
	GetID(data T) string
	SetID(data T, id string)
	MapStateToDataObject(ctx context.Context, state tfsdk.State) (T, diag.Diagnostics)
	UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, obj T) diag.Diagnostics
}

func (r *applicationAlertConfigResourceFramework) NewResource(ctx context.Context, api restapi.InstanaAPI) (ResourceFramework[*restapi.ApplicationAlertConfig], diag.Diagnostics) {
	return &applicationAlertConfigResourceFrameworkImpl{
		api:      api,
		isGlobal: r.isGlobal,
	}, nil
}

// applicationAlertConfigResource implements the ResourceFramework interface for ApplicationAlertConfig
type applicationAlertConfigResourceFrameworkImpl struct {
	api      restapi.InstanaAPI
	isGlobal bool
}

func (r *applicationAlertConfigResourceFrameworkImpl) GetID(data *restapi.ApplicationAlertConfig) string {
	return data.ID
}

func (r *applicationAlertConfigResourceFrameworkImpl) SetID(data *restapi.ApplicationAlertConfig, id string) {
	data.ID = id
}

func (r *applicationAlertConfigResourceFrameworkImpl) MapStateToDataObject(ctx context.Context, state tfsdk.State) (*restapi.ApplicationAlertConfig, diag.Diagnostics) {
	var model ApplicationAlertConfigModel
	diags := state.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}

	result := &restapi.ApplicationAlertConfig{
		ID:               model.ID.ValueString(),
		Name:             model.Name.ValueString(),
		Description:      model.Description.ValueString(),
		BoundaryScope:    restapi.BoundaryScope(model.BoundaryScope.ValueString()),
		EvaluationType:   restapi.ApplicationAlertEvaluationType(model.EvaluationType.ValueString()),
		IncludeInternal:  model.IncludeInternal.ValueBool(),
		IncludeSynthetic: model.IncludeSynthetic.ValueBool(),
		Triggering:       model.Triggering.ValueBool(),
	}

	// Handle grace period

	result.GracePeriod = extractGracePeriod(model.GracePeriod)

	// Handle granularity
	if !model.Granularity.IsNull() && !model.Granularity.IsUnknown() {
		result.Granularity = restapi.Granularity(model.Granularity.ValueInt64())
	}

	// Handle tag filter
	if !model.TagFilter.IsNull() && !model.TagFilter.IsUnknown() {
		tagFilterExpression := model.TagFilter.ValueString()
		if len(tagFilterExpression) > 0 {
			parsedExpression, err := tagfilter.ParseExpression(tagFilterExpression)
			if err != nil {
				return nil, diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Failed to parse tag filter expression",
						fmt.Sprintf("Invalid tag filter expression: %s", err.Error()),
					),
				}
			}
			result.TagFilterExpression = parsedExpression
		}
	}

	// Handle severity (deprecated but supported for backward compatibility)
	if !model.Severity.IsNull() && !model.Severity.IsUnknown() {
		severity := model.Severity.ValueString()
		result.Severity = restapi.SeverityFromTerraformRepresentation(severity)
	}

	// Handle alert channel IDs (deprecated but supported for backward compatibility)
	if !model.AlertChannelIDs.IsNull() && !model.AlertChannelIDs.IsUnknown() {
		var alertChannelIDs []string
		diags = model.AlertChannelIDs.ElementsAs(ctx, &alertChannelIDs, false)
		if diags.HasError() {
			return nil, diags
		}
		result.AlertChannelIDs = alertChannelIDs
	}

	// Handle alert channels (new format as map of severity to channel IDs)
	if !model.AlertChannels.IsNull() && !model.AlertChannels.IsUnknown() {
		alertChannels := make(map[string][]string)
		for k, v := range model.AlertChannels.Elements() {
			var channelIDs []string
			diags = v.(types.Set).ElementsAs(ctx, &channelIDs, false)
			if diags.HasError() {
				return nil, diags
			}
			alertChannels[k] = channelIDs
		}
		result.AlertChannels = alertChannels
	}

	// Handle applications
	if !model.Applications.IsNull() && !model.Applications.IsUnknown() {
		var applications []ApplicationModel
		diags = model.Applications.ElementsAs(ctx, &applications, false)
		if diags.HasError() {
			return nil, diags
		}

		result.Applications = make(map[string]restapi.IncludedApplication)
		for _, app := range applications {
			appID := app.ApplicationID.ValueString()
			services := make(map[string]restapi.IncludedService)

			if !app.Services.IsNull() && !app.Services.IsUnknown() {
				var serviceModels []ServiceModel
				diags = app.Services.ElementsAs(ctx, &serviceModels, false)
				if diags.HasError() {
					return nil, diags
				}

				for _, svc := range serviceModels {
					svcID := svc.ServiceID.ValueString()
					endpoints := make(map[string]restapi.IncludedEndpoint)

					if !svc.Endpoints.IsNull() && !svc.Endpoints.IsUnknown() {
						var endpointModels []EndpointModel
						diags = svc.Endpoints.ElementsAs(ctx, &endpointModels, false)
						if diags.HasError() {
							return nil, diags
						}

						for _, ep := range endpointModels {
							epID := ep.EndpointID.ValueString()
							endpoints[epID] = restapi.IncludedEndpoint{
								EndpointID: epID,
								Inclusive:  ep.Inclusive.ValueBool(),
							}
						}
					}

					services[svcID] = restapi.IncludedService{
						ServiceID: svcID,
						Inclusive: svc.Inclusive.ValueBool(),
						Endpoints: endpoints,
					}
				}
			}

			result.Applications[appID] = restapi.IncludedApplication{
				ApplicationID: appID,
				Inclusive:     app.Inclusive.ValueBool(),
				Services:      services,
			}
		}
	}

	// Handle custom payload fields
	if !model.CustomPayloadFields.IsNull() && !model.CustomPayloadFields.IsUnknown() {
		var customerPayloadFields []restapi.CustomPayloadField[any]
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
		result.CustomerPayloadFields = customerPayloadFields
	}

	// Handle rules (new format with multiple thresholds and severity levels)
	if !model.Rules.IsNull() && !model.Rules.IsUnknown() {
		var ruleWithThresholds []RuleWithThresholdModel
		diags = model.Rules.ElementsAs(ctx, &ruleWithThresholds, false)
		if diags.HasError() {
			return nil, diags
		}

		result.Rules = make([]restapi.ApplicationAlertRuleWithThresholds, len(ruleWithThresholds))
		for i, ruleWithThreshold := range ruleWithThresholds {
			result.Rules[i] = restapi.ApplicationAlertRuleWithThresholds{
				ThresholdOperator: ruleWithThreshold.ThresholdOperator.ValueString(),
			}

			// Handle rule configuration
			if !ruleWithThreshold.Rule.IsNull() && !ruleWithThreshold.Rule.IsUnknown() {
				var rule RuleModel
				diags = ruleWithThreshold.Rule.As(ctx, &rule, basetypes.ObjectAsOptions{})
				if diags.HasError() {
					return nil, diags
				}

				result.Rules[i].Rule = &restapi.ApplicationAlertRule{}

				// Handle error rate rule
				if !rule.ErrorRate.IsNull() && !rule.ErrorRate.IsUnknown() {
					var errorRateRule RuleConfigModel
					diags = rule.ErrorRate.As(ctx, &errorRateRule, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, diags
					}
					if errorRateRule.MetricName.IsNull() || errorRateRule.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for Error rate rules",
							"MetricName is required for Error rate rules",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleErrorRate
					result.Rules[i].Rule.MetricName = errorRateRule.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(errorRateRule.Aggregation.ValueString())
				}

				// Handle errors rule
				if !rule.Errors.IsNull() && !rule.Errors.IsUnknown() {
					var errorsRule RuleConfigModel
					diags = rule.Errors.As(ctx, &errorsRule, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, diags
					}
					if errorsRule.MetricName.IsNull() || errorsRule.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for Error rules",
							"MetricName is required for Error rules",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleErrors
					result.Rules[i].Rule.MetricName = errorsRule.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(errorsRule.Aggregation.ValueString())
				}

				// Handle logs rule
				if !rule.Logs.IsNull() && !rule.Logs.IsUnknown() {
					var logsRule LogsRuleModel
					diags = rule.Logs.As(ctx, &logsRule, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, diags
					}
					if logsRule.MetricName.IsNull() || logsRule.MetricName.IsUnknown() ||
						logsRule.Level.IsNull() || logsRule.Level.IsUnknown() ||
						logsRule.Operator.IsNull() || logsRule.Operator.IsUnknown() {
						diags.AddError(
							"MetricName,log level,log operator are required for log rules",
							"MetricName,log level,log operator are required for log rules",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleLogs
					result.Rules[i].Rule.MetricName = logsRule.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(logsRule.Aggregation.ValueString())

					// Set additional fields for logs
					level := restapi.LogLevel(logsRule.Level.ValueString())
					result.Rules[i].Rule.Level = &level

					message := logsRule.Message.ValueString()
					result.Rules[i].Rule.Message = &message

					operator := restapi.ExpressionOperator(logsRule.Operator.ValueString())
					result.Rules[i].Rule.Operator = &operator
				}

				// Handle slowness rule
				if !rule.Slowness.IsNull() && !rule.Slowness.IsUnknown() {
					log.Printf("Slowness : %+v\n", rule.Slowness)
					var slownessRule RuleConfigModel
					diags = rule.Slowness.As(ctx, &slownessRule, basetypes.ObjectAsOptions{})
					log.Printf("Slowness mapped : %+v\n", slownessRule)
					if diags.HasError() {
						return nil, diags
					}
					if slownessRule.MetricName.IsNull() || slownessRule.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for slowness",
							"MetricName is required for slowness",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleSlowness
					result.Rules[i].Rule.MetricName = slownessRule.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(slownessRule.Aggregation.ValueString())
				}

				// Handle status code rule
				if !rule.StatusCode.IsNull() && !rule.StatusCode.IsUnknown() {
					var statusCodeRule StatusCodeRuleModel
					diags = rule.StatusCode.As(ctx, &statusCodeRule, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, diags
					}
					if statusCodeRule.MetricName.IsNull() || statusCodeRule.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for status code",
							"MetricName is required for statuc code",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleStatusCode
					result.Rules[i].Rule.MetricName = statusCodeRule.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(statusCodeRule.Aggregation.ValueString())

					// Set additional fields for status code
					statusCodeStart := int32(statusCodeRule.StatusCodeStart.ValueInt64())
					result.Rules[i].Rule.StatusCodeStart = &statusCodeStart

					statusCodeEnd := int32(statusCodeRule.StatusCodeEnd.ValueInt64())
					result.Rules[i].Rule.StatusCodeEnd = &statusCodeEnd
				}

				// Handle throughput rule
				if !rule.Throughput.IsNull() && !rule.Throughput.IsUnknown() {
					var throughputRule RuleConfigModel
					diags = rule.Throughput.As(ctx, &throughputRule, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, diags
					}
					if throughputRule.MetricName.IsNull() || throughputRule.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for througput",
							"MetricName is required for throughput",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleThroughput
					result.Rules[i].Rule.MetricName = throughputRule.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(throughputRule.Aggregation.ValueString())
				}
			}

			// Handle thresholds
			var thresholdMap map[restapi.AlertSeverity]restapi.ThresholdRule
			var thresholdDiags diag.Diagnostics

			if !ruleWithThreshold.Thresholds.IsNull() && !ruleWithThreshold.Thresholds.IsUnknown() {
				thresholdMap, thresholdDiags = shared.MapThresholdsFromStateObject(ctx, ruleWithThreshold.Thresholds)
				diags.Append(thresholdDiags...)
				if diags.HasError() {
					return nil, diags
				}
				result.Rules[i].Thresholds = thresholdMap
			}

			// // Handle thresholds
			// if !ruleWithThreshold.Thresholds.IsNull() && !ruleWithThreshold.Thresholds.IsUnknown() {
			// 	thresholds := make(map[string]restapi.ThresholdValue)
			// 	for k, v := range ruleWithThreshold.Thresholds.Elements() {
			// 		var thresholdConfig ThresholdConfigRuleModel
			// 		diags = v.(types.Object).As(ctx, &thresholdConfig, basetypes.ObjectAsOptions{})
			// 		if diags.HasError() {
			// 			return nil, diags
			// 		}
			// 		thresholds[k] = restapi.ThresholdValue{
			// 			Value: thresholdConfig.Value.ValueFloat64(),
			// 		}
			// 	}
			// 	result.Rules[i].Thresholds = thresholdMap
			// }
		}
	}

	// Handle time threshold
	if !model.TimeThreshold.IsNull() && !model.TimeThreshold.IsUnknown() {
		var timeThreshold AppAlertTimeThresholdModel
		diags = model.TimeThreshold.As(ctx, &timeThreshold, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		result.TimeThreshold = &restapi.ApplicationAlertTimeThreshold{}

		// Handle request impact
		if !timeThreshold.RequestImpact.IsNull() && !timeThreshold.RequestImpact.IsUnknown() {
			var requestImpact AppAlertRequestImpactModel
			diags = timeThreshold.RequestImpact.As(ctx, &requestImpact, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return nil, diags
			}
			result.TimeThreshold.Type = "requestImpact"
			result.TimeThreshold.TimeWindow = requestImpact.TimeWindow.ValueInt64()
			result.TimeThreshold.Requests = int(requestImpact.Requests.ValueInt64())
		}

		// Handle violations in period
		if !timeThreshold.ViolationsInPeriod.IsNull() && !timeThreshold.ViolationsInPeriod.IsUnknown() {
			var violationsInPeriod AppAlertViolationsInPeriodModel
			diags = timeThreshold.ViolationsInPeriod.As(ctx, &violationsInPeriod, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return nil, diags
			}
			result.TimeThreshold.Type = "violationsInPeriod"
			result.TimeThreshold.TimeWindow = violationsInPeriod.TimeWindow.ValueInt64()
			result.TimeThreshold.Violations = int(violationsInPeriod.Violations.ValueInt64())
		}

		// Handle violations in sequence
		if !timeThreshold.ViolationsInSequence.IsNull() && !timeThreshold.ViolationsInSequence.IsUnknown() {
			var violationsInSequence AppAlertViolationsInSequenceModel
			diags = timeThreshold.ViolationsInSequence.As(ctx, &violationsInSequence, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return nil, diags
			}
			result.TimeThreshold.Type = "violationsInSequence"
			result.TimeThreshold.TimeWindow = violationsInSequence.TimeWindow.ValueInt64()
		}
	}

	return result, nil
}

func extractGracePeriod(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil // send as null (omitted in JSON)
	}
	val := v.ValueInt64()
	return &val
}

func (r *applicationAlertConfigResourceFrameworkImpl) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	model := ApplicationAlertConfigModel{
		ID:               types.StringValue(data.ID),
		Name:             types.StringValue(data.Name),
		Description:      types.StringValue(data.Description),
		BoundaryScope:    types.StringValue(string(data.BoundaryScope)),
		EvaluationType:   types.StringValue(string(data.EvaluationType)),
		Granularity:      types.Int64Value(int64(data.Granularity)),
		IncludeInternal:  types.BoolValue(data.IncludeInternal),
		IncludeSynthetic: types.BoolValue(data.IncludeSynthetic),
		Triggering:       types.BoolValue(data.Triggering),
	}

	// Handle grace period
	if data.GracePeriod != nil {
		model.GracePeriod = util.SetInt64PointerToState(data.GracePeriod)
	}

	// Handle tag filter
	if data.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(data.TagFilterExpression)
		if err != nil {
			diags.AddError(
				"Error normalizing tag filter",
				"Could not normalize tag filter: "+err.Error(),
			)
			return diags
		}

		model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)

	} else {
		model.TagFilter = types.StringNull()
	}

	// Handle severity (deprecated but supported for backward compatibility)
	if data.Severity > 0 {
		severity, err := util.ConvertSeverityFromInstanaAPIToTerraformRepresentation(data.Severity)
		if err == nil {
			model.Severity = types.StringValue(severity)
		}
	}
	if diags.HasError() {
		return diags
	}
	log.Printf("Before allertchannel id stage")
	// Handle alert channel IDs (deprecated but supported for backward compatibility)
	if len(data.AlertChannelIDs) > 0 {
		elements := make([]attr.Value, len(data.AlertChannelIDs))
		for i, id := range data.AlertChannelIDs {
			elements[i] = types.StringValue(id)
		}
		model.AlertChannelIDs = types.SetValueMust(types.StringType, elements)
	} else {
		model.AlertChannelIDs = types.SetNull(types.StringType)
	}

	log.Printf("Before alert channel stage")
	// Handle alert channels (new format as map of severity to channel IDs)
	if len(data.AlertChannels) > 0 {
		elements := make(map[string]attr.Value)
		for severity, channelIDs := range data.AlertChannels {
			channelElements := make([]attr.Value, len(channelIDs))
			for i, id := range channelIDs {
				channelElements[i] = types.StringValue(id)
			}
			elements[severity] = types.SetValueMust(types.StringType, channelElements)
		}
		model.AlertChannels = types.MapValueMust(types.SetType{ElemType: types.StringType}, elements)
		log.Printf("static threshold elements : %+v\n", model.AlertChannels)

	} else {
		model.AlertChannels = types.MapNull(types.SetType{ElemType: types.StringType})
	}
	if diags.HasError() {
		return diags
	}
	log.Printf("Before applicaion stage")
	// Handle applications
	// Predefine attribute/type maps once
	endpointAttrTypes := map[string]attr.Type{
		"endpoint_id": types.StringType,
		"inclusive":   types.BoolType,
	}
	endpointObjectType := types.ObjectType{AttrTypes: endpointAttrTypes}
	endpointSetType := types.SetType{ElemType: endpointObjectType}

	serviceAttrTypes := map[string]attr.Type{
		"service_id": types.StringType,
		"inclusive":  types.BoolType,
		"endpoint":   endpointSetType, // key must match tfsdk tag on ServiceModel
	}
	serviceObjectType := types.ObjectType{AttrTypes: serviceAttrTypes}
	serviceSetType := types.SetType{ElemType: serviceObjectType}

	applicationAttrTypes := map[string]attr.Type{
		"application_id": types.StringType,
		"inclusive":      types.BoolType,
		"service":        serviceSetType, // key must match tfsdk tag on ApplicationModel
	}
	applicationObjectType := types.ObjectType{AttrTypes: applicationAttrTypes}

	if len(data.Applications) > 0 {
		appElements := make([]attr.Value, 0, len(data.Applications))
		for _, app := range data.Applications {
			appModel := ApplicationModel{
				ApplicationID: types.StringValue(app.ApplicationID),
				Inclusive:     types.BoolValue(app.Inclusive),
			}

			svcElements := make([]attr.Value, 0, len(app.Services))
			for _, svc := range app.Services {
				svcModel := ServiceModel{
					ServiceID: types.StringValue(svc.ServiceID),
					Inclusive: types.BoolValue(svc.Inclusive),
				}

				epElements := make([]attr.Value, 0, len(svc.Endpoints))
				for _, ep := range svc.Endpoints {
					epModel := EndpointModel{
						EndpointID: types.StringValue(ep.EndpointID),
						Inclusive:  types.BoolValue(ep.Inclusive),
					}
					epObj, diags := types.ObjectValueFrom(ctx, endpointAttrTypes, epModel)
					if diags.HasError() {
						return diags
					}
					epElements = append(epElements, epObj)
				}
				svcModel.Endpoints = types.SetValueMust(endpointObjectType, epElements)
				svcObj, diags := types.ObjectValueFrom(ctx, serviceAttrTypes, svcModel)
				if diags.HasError() {
					return diags
				}
				svcElements = append(svcElements, svcObj)

			}
			appModel.Services = types.SetValueMust(serviceObjectType, svcElements)

			appObj, diags := types.ObjectValueFrom(ctx, applicationAttrTypes, appModel)
			if diags.HasError() {
				return diags
			}
			appElements = append(appElements, appObj)
		}

		model.Applications = types.SetValueMust(applicationObjectType, appElements)
	}
	if diags.HasError() {
		return diags
	}

	log.Printf("Before custom payload stage")
	// Handle custom payload fields
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, data.CustomerPayloadFields)
	if payloadDiags.HasError() {
		return payloadDiags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	if diags.HasError() {
		return diags
	}

	log.Printf("Before Rules stage")
	// Handle rules (new format with multiple thresholds and severity levels)
	if len(data.Rules) > 0 {
		rulesElements := make([]attr.Value, len(data.Rules))
		for i, ruleWithThreshold := range data.Rules {
			// Create rule model
			ruleModel := RuleModel{}

			// Set rule model fields based on the AlertType
			switch ruleWithThreshold.Rule.AlertType {
			case ApplicationAlertConfigFieldRuleErrorRate:
				ruleModel.ErrorRate = types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
					"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				})
				ruleModel.Errors = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Logs = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType})
				ruleModel.Slowness = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.StatusCode = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type})
				ruleModel.Throughput = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})

			case ApplicationAlertConfigFieldRuleErrors:
				ruleModel.ErrorRate = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Errors = types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
					"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				})
				ruleModel.Logs = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType})
				ruleModel.Slowness = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.StatusCode = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type})
				ruleModel.Throughput = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})

			case ApplicationAlertConfigFieldRuleLogs:
				// For logs, we need to handle additional fields
				level := ""
				message := ""
				operator := ""

				if ruleWithThreshold.Rule.Level != nil {
					level = string(*ruleWithThreshold.Rule.Level)
				}
				if ruleWithThreshold.Rule.Message != nil {
					message = *ruleWithThreshold.Rule.Message
				}
				if ruleWithThreshold.Rule.Operator != nil {
					operator = string(*ruleWithThreshold.Rule.Operator)
				}

				ruleModel.ErrorRate = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Errors = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Logs = types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
					"level":       types.StringType,
					"message":     types.StringType,
					"operator":    types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
					"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					"level":       types.StringValue(level),
					"message":     types.StringValue(message),
					"operator":    types.StringValue(operator),
				})
				ruleModel.Slowness = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.StatusCode = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type})
				ruleModel.Throughput = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})

			case ApplicationAlertConfigFieldRuleSlowness:
				ruleModel.ErrorRate = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Errors = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Logs = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType})
				ruleModel.Slowness = types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
					"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				})
				ruleModel.StatusCode = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type})
				ruleModel.Throughput = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})

			case ApplicationAlertConfigFieldRuleStatusCode:
				// For status code, we need to handle additional fields
				statusCodeStart := int64(0)
				statusCodeEnd := int64(0)

				if ruleWithThreshold.Rule.StatusCodeStart != nil {
					statusCodeStart = int64(*ruleWithThreshold.Rule.StatusCodeStart)
				}
				if ruleWithThreshold.Rule.StatusCodeEnd != nil {
					statusCodeEnd = int64(*ruleWithThreshold.Rule.StatusCodeEnd)
				}

				ruleModel.ErrorRate = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Errors = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Logs = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType})
				ruleModel.Slowness = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.StatusCode = types.ObjectValueMust(map[string]attr.Type{
					"metric_name":       types.StringType,
					"aggregation":       types.StringType,
					"status_code_start": types.Int64Type,
					"status_code_end":   types.Int64Type,
				}, map[string]attr.Value{
					"metric_name":       types.StringValue(ruleWithThreshold.Rule.MetricName),
					"aggregation":       types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					"status_code_start": types.Int64Value(statusCodeStart),
					"status_code_end":   types.Int64Value(statusCodeEnd),
				})
				ruleModel.Throughput = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})

			case ApplicationAlertConfigFieldRuleThroughput:
				ruleModel.ErrorRate = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Errors = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.Logs = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType})
				ruleModel.Slowness = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType})
				ruleModel.StatusCode = types.ObjectNull(map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type})
				ruleModel.Throughput = types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
					"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				})
			}

			// Define rule attribute types
			ruleAttrTypes := map[string]attr.Type{
				"error_rate":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
				"errors":      types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
				"logs":        types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}},
				"slowness":    types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
				"status_code": types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}},
				"throughput":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
			}

			ruleObj, diags := types.ObjectValueFrom(ctx, ruleAttrTypes, ruleModel)
			if diags.HasError() {
				return diags
			}

			// Create thresholds map

			// Map thresholds
			thresholdObj := map[string]attr.Value{}

			// Map warning threshold
			warningThreshold, isWarningThresholdPresent := ruleWithThreshold.Thresholds[restapi.WarningSeverity]
			warningThresholdList, warningDiags := shared.MapThresholdToState(ctx, isWarningThresholdPresent, &warningThreshold, []string{"static", "adaptiveBaseline"})
			diags.Append(warningDiags...)
			if diags.HasError() {
				return diags
			}
			thresholdObj[shared.LogAlertConfigFieldWarning] = warningThresholdList

			// Map critical threshold
			criticalThreshold, isCriticalThresholdPresent := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]
			criticalThresholdList, criticalDiags := shared.MapThresholdToState(ctx, isCriticalThresholdPresent, &criticalThreshold, []string{"static", "adaptiveBaseline"})
			diags.Append(criticalDiags...)
			if diags.HasError() {
				return diags
			}
			thresholdObj[shared.LogAlertConfigFieldCritical] = criticalThresholdList

			// Create threshold object value
			thresholdObjVal, thresholdObjDiags := types.ObjectValue(
				map[string]attr.Type{
					shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
					shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
				},
				thresholdObj,
			)
			diags.Append(thresholdObjDiags...)
			if diags.HasError() {
				return diags
			}

			log.Printf("before calling rules value initial")
			ruleWithThresholdModel := RuleWithThresholdModel{}
			ruleWithThresholdModel.Rule = ruleObj
			ruleWithThresholdModel.ThresholdOperator = types.StringValue(ruleWithThreshold.ThresholdOperator)
			ruleWithThresholdModel.Thresholds = thresholdObjVal

			// Create rule with threshold object
			ruleWithThresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"rule": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"error_rate":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"errors":      types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"logs":        types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}},
						"slowness":    types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"status_code": types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}},
						"throughput":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
					},
				},
				"threshold_operator": types.StringType,
				"threshold": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
						shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
					},
				},
			}, ruleWithThresholdModel)
			if diags.HasError() {
				return diags
			}
			rulesElements[i] = ruleWithThresholdObj
		}
		log.Printf("before calling rules value final")
		model.Rules = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"rule": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"error_rate":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"errors":      types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"logs":        types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}},
						"slowness":    types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"status_code": types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}},
						"throughput":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
					},
				},
				"threshold_operator": types.StringType,
				"threshold": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
						shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
					},
				},
			},
		})
		model.Rules = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"rule": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"error_rate":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"errors":      types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"logs":        types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}},
						"slowness":    types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
						"status_code": types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}},
						"throughput":  types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}},
					},
				},
				"threshold_operator": types.StringType,
				"threshold": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						shared.LogAlertConfigFieldWarning:  shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
						shared.LogAlertConfigFieldCritical: shared.GetStaticAndAdaptiveThresholdAttrObjectTypes(),
					},
				},
			},
		}, rulesElements)
	}
	if diags.HasError() {
		return diags
	}
	log.Printf("Before TimeThreshold stage")
	// Handle time threshold
	if data.TimeThreshold != nil {
		timeThresholdModel := AppAlertTimeThresholdModel{}

		// Determine which time threshold to populate based on the Type field
		switch data.TimeThreshold.Type {
		case "requestImpact":
			requestImpactObj := types.ObjectValueMust(map[string]attr.Type{
				"time_window": types.Int64Type,
				"requests":    types.Int64Type,
			}, map[string]attr.Value{
				"time_window": types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
				"requests":    types.Int64Value(int64(data.TimeThreshold.Requests)),
			})
			timeThresholdModel.RequestImpact = requestImpactObj
			timeThresholdModel.ViolationsInPeriod = types.ObjectNull(map[string]attr.Type{
				"time_window": types.Int64Type,
				"violations":  types.Int64Type,
			})
			timeThresholdModel.ViolationsInSequence = types.ObjectNull(map[string]attr.Type{
				"time_window": types.Int64Type,
			})

		case "violationsInPeriod":
			violationsInPeriodObj := types.ObjectValueMust(map[string]attr.Type{
				"time_window": types.Int64Type,
				"violations":  types.Int64Type,
			}, map[string]attr.Value{
				"time_window": types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
				"violations":  types.Int64Value(int64(data.TimeThreshold.Violations)),
			})
			timeThresholdModel.RequestImpact = types.ObjectNull(map[string]attr.Type{
				"time_window": types.Int64Type,
				"requests":    types.Int64Type,
			})
			timeThresholdModel.ViolationsInPeriod = violationsInPeriodObj
			timeThresholdModel.ViolationsInSequence = types.ObjectNull(map[string]attr.Type{
				"time_window": types.Int64Type,
			})

		case "violationsInSequence":
			violationsInSequenceObj := types.ObjectValueMust(map[string]attr.Type{
				"time_window": types.Int64Type,
			}, map[string]attr.Value{
				"time_window": types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
			})
			timeThresholdModel.RequestImpact = types.ObjectNull(map[string]attr.Type{
				"time_window": types.Int64Type,
				"requests":    types.Int64Type,
			})
			timeThresholdModel.ViolationsInPeriod = types.ObjectNull(map[string]attr.Type{
				"time_window": types.Int64Type,
				"violations":  types.Int64Type,
			})
			timeThresholdModel.ViolationsInSequence = violationsInSequenceObj
		}

		// Define attribute types for time threshold
		timeThresholdAttrTypes := map[string]attr.Type{
			"request_impact":         types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type, "requests": types.Int64Type}},
			"violations_in_period":   types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type, "violations": types.Int64Type}},
			"violations_in_sequence": types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type}},
		}

		timeThresholdObj, diags := types.ObjectValueFrom(ctx, timeThresholdAttrTypes, timeThresholdModel)
		if diags.HasError() {
			return diags
		}

		model.TimeThreshold = timeThresholdObj
	}
	if diags.HasError() {
		return diags
	}
	log.Printf("Reached final stage")
	log.Printf("static threshold elements : %+v\n", model)
	return state.Set(ctx, model)
}
