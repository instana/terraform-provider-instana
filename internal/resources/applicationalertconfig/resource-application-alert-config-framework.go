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
	if len(model.AlertChannelIDs) > 0 {
		alertChannelIDs := make([]string, len(model.AlertChannelIDs))
		for i, id := range model.AlertChannelIDs {
			alertChannelIDs[i] = id.ValueString()
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
	if len(model.Applications) > 0 {
		result.Applications = make(map[string]restapi.IncludedApplication)
		for _, app := range model.Applications {
			appID := app.ApplicationID.ValueString()
			services := make(map[string]restapi.IncludedService)

			if len(app.Services) > 0 {
				for _, svc := range app.Services {
					svcID := svc.ServiceID.ValueString()
					endpoints := make(map[string]restapi.IncludedEndpoint)

					if len(svc.Endpoints) > 0 {
						for _, ep := range svc.Endpoints {
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
	if len(model.Rules) > 0 {
		result.Rules = make([]restapi.ApplicationAlertRuleWithThresholds, len(model.Rules))
		for i, ruleWithThreshold := range model.Rules {
			result.Rules[i] = restapi.ApplicationAlertRuleWithThresholds{
				ThresholdOperator: ruleWithThreshold.ThresholdOperator.ValueString(),
			}

			// Handle rule configuration
			if ruleWithThreshold.Rule != nil {
				rule := ruleWithThreshold.Rule

				result.Rules[i].Rule = &restapi.ApplicationAlertRule{}

				// Handle error rate rule
				if rule.ErrorRate != nil {
					if rule.ErrorRate.MetricName.IsNull() || rule.ErrorRate.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for Error rate rules",
							"MetricName is required for Error rate rules",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleErrorRate
					result.Rules[i].Rule.MetricName = rule.ErrorRate.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(rule.ErrorRate.Aggregation.ValueString())
				}

				// Handle errors rule
				if rule.Errors != nil {
					if rule.Errors.MetricName.IsNull() || rule.Errors.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for Error rules",
							"MetricName is required for Error rules",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleErrors
					result.Rules[i].Rule.MetricName = rule.Errors.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(rule.Errors.Aggregation.ValueString())
				}

				// Handle logs rule
				if rule.Logs != nil {
					if rule.Logs.MetricName.IsNull() || rule.Logs.MetricName.IsUnknown() ||
						rule.Logs.Level.IsNull() || rule.Logs.Level.IsUnknown() ||
						rule.Logs.Operator.IsNull() || rule.Logs.Operator.IsUnknown() {
						diags.AddError(
							"MetricName,log level,log operator are required for log rules",
							"MetricName,log level,log operator are required for log rules",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleLogs
					result.Rules[i].Rule.MetricName = rule.Logs.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(rule.Logs.Aggregation.ValueString())

					// Set additional fields for logs
					level := restapi.LogLevel(rule.Logs.Level.ValueString())
					result.Rules[i].Rule.Level = &level

					message := rule.Logs.Message.ValueString()
					result.Rules[i].Rule.Message = &message

					operator := restapi.ExpressionOperator(rule.Logs.Operator.ValueString())
					result.Rules[i].Rule.Operator = &operator
				}

				// Handle slowness rule
				if rule.Slowness != nil {
					log.Printf("Slowness : %+v\n", rule.Slowness)
					if rule.Slowness.MetricName.IsNull() || rule.Slowness.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for slowness",
							"MetricName is required for slowness",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleSlowness
					result.Rules[i].Rule.MetricName = rule.Slowness.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(rule.Slowness.Aggregation.ValueString())
				}

				// Handle status code rule
				if rule.StatusCode != nil {
					if rule.StatusCode.MetricName.IsNull() || rule.StatusCode.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for status code",
							"MetricName is required for statuc code",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleStatusCode
					result.Rules[i].Rule.MetricName = rule.StatusCode.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(rule.StatusCode.Aggregation.ValueString())

					// Set additional fields for status code
					statusCodeStart := int32(rule.StatusCode.StatusCodeStart.ValueInt64())
					result.Rules[i].Rule.StatusCodeStart = &statusCodeStart

					statusCodeEnd := int32(rule.StatusCode.StatusCodeEnd.ValueInt64())
					result.Rules[i].Rule.StatusCodeEnd = &statusCodeEnd
				}

				// Handle throughput rule
				if rule.Throughput != nil {
					if rule.Throughput.MetricName.IsNull() || rule.Throughput.MetricName.IsUnknown() {
						diags.AddError(
							"MetricName is required for througput",
							"MetricName is required for throughput",
						)
						return nil, diags
					}
					result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleThroughput
					result.Rules[i].Rule.MetricName = rule.Throughput.MetricName.ValueString()
					result.Rules[i].Rule.Aggregation = restapi.Aggregation(rule.Throughput.Aggregation.ValueString())
				}
			}

			// Handle thresholds
			if ruleWithThreshold.Thresholds != nil {
				thresholdMap := make(map[restapi.AlertSeverity]restapi.ThresholdRule)

				// Map warning threshold
				if ruleWithThreshold.Thresholds.Static != nil {
					warningThreshold := &restapi.ThresholdRule{
						Type: "staticThreshold",
					}
					if !ruleWithThreshold.Thresholds.Static.Value.IsNull() && !ruleWithThreshold.Thresholds.Static.Value.IsUnknown() {
						value := float64(ruleWithThreshold.Thresholds.Static.Value.ValueInt64())
						warningThreshold.Value = &value
					}
					thresholdMap[restapi.WarningSeverity] = *warningThreshold
				}

				// Map adaptive baseline threshold
				if ruleWithThreshold.Thresholds.AdaptiveBaseline != nil {
					adaptiveThreshold := &restapi.ThresholdRule{
						Type: "adaptiveBaseline",
					}
					if !ruleWithThreshold.Thresholds.AdaptiveBaseline.DeviationFactor.IsNull() && !ruleWithThreshold.Thresholds.AdaptiveBaseline.DeviationFactor.IsUnknown() {
						deviation := ruleWithThreshold.Thresholds.AdaptiveBaseline.DeviationFactor.ValueFloat32()
						adaptiveThreshold.DeviationFactor = &deviation
					}
					if !ruleWithThreshold.Thresholds.AdaptiveBaseline.Adaptability.IsNull() && !ruleWithThreshold.Thresholds.AdaptiveBaseline.Adaptability.IsUnknown() {
						adaptability := ruleWithThreshold.Thresholds.AdaptiveBaseline.Adaptability.ValueFloat32()
						adaptiveThreshold.Adaptability = &adaptability
					}
					if !ruleWithThreshold.Thresholds.AdaptiveBaseline.Seasonality.IsNull() && !ruleWithThreshold.Thresholds.AdaptiveBaseline.Seasonality.IsUnknown() {
						seasonality := restapi.ThresholdSeasonality(ruleWithThreshold.Thresholds.AdaptiveBaseline.Seasonality.ValueString())
						adaptiveThreshold.Seasonality = &seasonality
					}
					thresholdMap[restapi.CriticalSeverity] = *adaptiveThreshold
				}

				result.Rules[i].Thresholds = thresholdMap
			}
		}
	}

	// Handle time threshold
	if model.TimeThreshold != nil {
		result.TimeThreshold = &restapi.ApplicationAlertTimeThreshold{}

		// Handle request impact
		if model.TimeThreshold.RequestImpact != nil {
			result.TimeThreshold.Type = "requestImpact"
			result.TimeThreshold.TimeWindow = model.TimeThreshold.RequestImpact.TimeWindow.ValueInt64()
			result.TimeThreshold.Requests = int(model.TimeThreshold.RequestImpact.Requests.ValueInt64())
		}

		// Handle violations in period
		if model.TimeThreshold.ViolationsInPeriod != nil {
			result.TimeThreshold.Type = "violationsInPeriod"
			result.TimeThreshold.TimeWindow = model.TimeThreshold.ViolationsInPeriod.TimeWindow.ValueInt64()
			result.TimeThreshold.Violations = int(model.TimeThreshold.ViolationsInPeriod.Violations.ValueInt64())
		}

		// Handle violations in sequence
		if model.TimeThreshold.ViolationsInSequence != nil {
			result.TimeThreshold.Type = "violationsInSequence"
			result.TimeThreshold.TimeWindow = model.TimeThreshold.ViolationsInSequence.TimeWindow.ValueInt64()
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
		model.AlertChannelIDs = make([]types.String, len(data.AlertChannelIDs))
		for i, id := range data.AlertChannelIDs {
			model.AlertChannelIDs[i] = types.StringValue(id)
		}
	} else {
		model.AlertChannelIDs = []types.String{}
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
	if len(data.Applications) > 0 {
		model.Applications = make([]ApplicationModel, 0, len(data.Applications))
		for _, app := range data.Applications {
			appModel := ApplicationModel{
				ApplicationID: types.StringValue(app.ApplicationID),
				Inclusive:     types.BoolValue(app.Inclusive),
			}

			if len(app.Services) > 0 {
				appModel.Services = make([]ServiceModel, 0, len(app.Services))
				for _, svc := range app.Services {
					svcModel := ServiceModel{
						ServiceID: types.StringValue(svc.ServiceID),
						Inclusive: types.BoolValue(svc.Inclusive),
					}

					if len(svc.Endpoints) > 0 {
						svcModel.Endpoints = make([]EndpointModel, 0, len(svc.Endpoints))
						for _, ep := range svc.Endpoints {
							epModel := EndpointModel{
								EndpointID: types.StringValue(ep.EndpointID),
								Inclusive:  types.BoolValue(ep.Inclusive),
							}
							svcModel.Endpoints = append(svcModel.Endpoints, epModel)
						}
					} else {
						svcModel.Endpoints = []EndpointModel{}
					}

					appModel.Services = append(appModel.Services, svcModel)
				}
			} else {
				appModel.Services = []ServiceModel{}
			}

			model.Applications = append(model.Applications, appModel)
		}
	} else {
		model.Applications = []ApplicationModel{}
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
		ruleModels := make([]RuleWithThresholdModel, len(data.Rules))
		for i, ruleWithThreshold := range data.Rules {
			// Create rule model
			ruleModel := RuleModel{}
			// initialize all with nil
			ruleModel.Errors = nil
			ruleModel.Logs = nil
			ruleModel.Slowness = nil
			ruleModel.StatusCode = nil
			ruleModel.Throughput = nil
			ruleModel.ErrorRate = nil

			// Set rule model fields based on the AlertType
			switch ruleWithThreshold.Rule.AlertType {
			case ApplicationAlertConfigFieldRuleErrorRate:
				ruleModel.ErrorRate = &RuleConfigModel{
					MetricName:  types.StringValue(ruleWithThreshold.Rule.MetricName),
					Aggregation: types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				}

			case ApplicationAlertConfigFieldRuleErrors:
				ruleModel.ErrorRate = nil
				ruleModel.Errors = &RuleConfigModel{
					MetricName:  types.StringValue(ruleWithThreshold.Rule.MetricName),
					Aggregation: types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				}

			case ApplicationAlertConfigFieldRuleLogs:
				// For logs, we need to handle additional fields
				var level types.String
				var message types.String
				var operator types.String

				if ruleWithThreshold.Rule.Level != nil {
					level = types.StringValue(string(*ruleWithThreshold.Rule.Level))
				} else {
					level = types.StringNull()
				}
				if ruleWithThreshold.Rule.Message != nil {
					message = types.StringValue(*ruleWithThreshold.Rule.Message)
				} else {
					message = types.StringNull()
				}
				if ruleWithThreshold.Rule.Operator != nil {
					operator = types.StringValue(string(*ruleWithThreshold.Rule.Operator))
				} else {
					operator = types.StringNull()
				}

				ruleModel.Logs = &LogsRuleModel{
					MetricName:  types.StringValue(ruleWithThreshold.Rule.MetricName),
					Aggregation: types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					Level:       level,
					Message:     message,
					Operator:    operator,
				}

			case ApplicationAlertConfigFieldRuleSlowness:

				ruleModel.Slowness = &RuleConfigModel{
					MetricName:  types.StringValue(ruleWithThreshold.Rule.MetricName),
					Aggregation: types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				}

			case ApplicationAlertConfigFieldRuleStatusCode:
				// For status code, we need to handle additional fields
				var statusCodeStart types.Int64
				var statusCodeEnd types.Int64

				if ruleWithThreshold.Rule.StatusCodeStart != nil {
					statusCodeStart = types.Int64Value(int64(*ruleWithThreshold.Rule.StatusCodeStart))
				} else {
					statusCodeStart = types.Int64Null()
				}
				if ruleWithThreshold.Rule.StatusCodeEnd != nil {
					statusCodeEnd = types.Int64Value(int64(*ruleWithThreshold.Rule.StatusCodeEnd))
				} else {
					statusCodeEnd = types.Int64Null()
				}

				ruleModel.StatusCode = &StatusCodeRuleModel{
					MetricName:      types.StringValue(ruleWithThreshold.Rule.MetricName),
					Aggregation:     types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					StatusCodeStart: statusCodeStart,
					StatusCodeEnd:   statusCodeEnd,
				}
			case ApplicationAlertConfigFieldRuleThroughput:
				ruleModel.Throughput = &RuleConfigModel{
					MetricName:  types.StringValue(ruleWithThreshold.Rule.MetricName),
					Aggregation: types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
				}
			}

			log.Printf("before calling rules value initial")
			ruleWithThresholdModel := RuleWithThresholdModel{}
			ruleWithThresholdModel.Rule = &ruleModel
			ruleWithThresholdModel.ThresholdOperator = types.StringValue(ruleWithThreshold.ThresholdOperator)

			// Map thresholds to ApplicationThresholdModel
			if len(ruleWithThreshold.Thresholds) > 0 {
				thresholdModel := &ApplicationThresholdModel{}

				// Check for warning threshold (static)
				if warningThreshold, ok := ruleWithThreshold.Thresholds[restapi.WarningSeverity]; ok {
					if warningThreshold.Type == "staticThreshold" && warningThreshold.Value != nil {
						thresholdModel.Static = &shared.StaticTypeModel{
							Value: types.Int64Value(int64(*warningThreshold.Value)),
						}
					}
				}

				// Check for critical threshold (adaptive baseline)
				if criticalThreshold, ok := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]; ok {
					if criticalThreshold.Type == "adaptiveBaseline" {
						adaptiveModel := &shared.AdaptiveBaselineModel{}
						if criticalThreshold.DeviationFactor != nil {
							adaptiveModel.DeviationFactor = types.Float32Value(*criticalThreshold.DeviationFactor)
						} else {
							adaptiveModel.DeviationFactor = types.Float32Null()
						}
						if criticalThreshold.Adaptability != nil {
							adaptiveModel.Adaptability = types.Float32Value(*criticalThreshold.Adaptability)
						} else {
							adaptiveModel.Adaptability = types.Float32Null()
						}
						if criticalThreshold.Seasonality != nil {
							adaptiveModel.Seasonality = types.StringValue(string(*criticalThreshold.Seasonality))
						} else {
							adaptiveModel.Seasonality = types.StringNull()
						}
						thresholdModel.AdaptiveBaseline = adaptiveModel
					}
				}

				ruleWithThresholdModel.Thresholds = thresholdModel
			}

			ruleModels[i] = ruleWithThresholdModel
		}

		log.Printf("before calling rules value final")
		// Directly assign the slice of models
		model.Rules = ruleModels
	} else {
		model.Rules = []RuleWithThresholdModel{}
	}
	if diags.HasError() {
		return diags
	}
	log.Printf("Before TimeThreshold stage")
	// Handle time threshold
	if data.TimeThreshold != nil {
		timeThresholdModel := &AppAlertTimeThresholdModel{}

		// Determine which time threshold to populate based on the Type field
		switch data.TimeThreshold.Type {
		case "requestImpact":
			timeThresholdModel.RequestImpact = &AppAlertRequestImpactModel{
				TimeWindow: types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
				Requests:   types.Int64Value(int64(data.TimeThreshold.Requests)),
			}
			timeThresholdModel.ViolationsInPeriod = nil
			timeThresholdModel.ViolationsInSequence = nil

		case "violationsInPeriod":
			timeThresholdModel.RequestImpact = nil
			timeThresholdModel.ViolationsInPeriod = &AppAlertViolationsInPeriodModel{
				TimeWindow: types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
				Violations: types.Int64Value(int64(data.TimeThreshold.Violations)),
			}
			timeThresholdModel.ViolationsInSequence = nil

		case "violationsInSequence":
			timeThresholdModel.RequestImpact = nil
			timeThresholdModel.ViolationsInPeriod = nil
			timeThresholdModel.ViolationsInSequence = &AppAlertViolationsInSequenceModel{
				TimeWindow: types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
			}
		}

		model.TimeThreshold = timeThresholdModel
	}
	if diags.HasError() {
		return diags
	}
	log.Printf("Reached final stage")
	log.Printf("static threshold elements : %+v\n", model)
	return state.Set(ctx, model)
}
