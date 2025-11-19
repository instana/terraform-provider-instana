package applicationalertconfig

import (
	"context"
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
										shared.LogAlertConfigFieldWarning:  shared.AllThresholdAttributeSchema(),
										shared.LogAlertConfigFieldCritical: shared.AllThresholdAttributeSchema(),
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
										Validators: []validator.Int64{
											int64validator.Between(1, 12),
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

// ============================================================================
// Resource Framework Implementation
// ============================================================================

type applicationAlertConfigResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
	isGlobal bool
}

// MetaData returns the resource metadata
func (r *applicationAlertConfigResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

// GetRestResource returns the appropriate REST resource based on whether this is a global config
func (r *applicationAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig] {
	if r.isGlobal {
		return api.GlobalApplicationAlertConfigs()
	}
	return api.ApplicationAlertConfigs()
}

// MapStateToDataObject converts Terraform state to API data object by delegating to the implementation
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

// UpdateState updates Terraform state with API data by delegating to the implementation
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

// SetComputedFields sets computed fields in the plan (none for this resource)
func (r *applicationAlertConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ============================================================================
// Resource Framework Interface
// ============================================================================

// ResourceFramework interface for application alert config resources
type ResourceFramework[T restapi.InstanaDataObject] interface {
	GetID(data T) string
	SetID(data T, id string)
	MapStateToDataObject(ctx context.Context, state tfsdk.State) (T, diag.Diagnostics)
	UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, obj T) diag.Diagnostics
}

// NewResource creates a new resource implementation instance
func (r *applicationAlertConfigResourceFramework) NewResource(ctx context.Context, api restapi.InstanaAPI) (ResourceFramework[*restapi.ApplicationAlertConfig], diag.Diagnostics) {
	return &applicationAlertConfigResourceFrameworkImpl{
		api:      api,
		isGlobal: r.isGlobal,
	}, nil
}

// ============================================================================
// Resource Implementation
// ============================================================================

// applicationAlertConfigResourceFrameworkImpl implements the ResourceFramework interface for ApplicationAlertConfig
type applicationAlertConfigResourceFrameworkImpl struct {
	api      restapi.InstanaAPI
	isGlobal bool
}

// GetID returns the ID from the API data object
func (r *applicationAlertConfigResourceFrameworkImpl) GetID(data *restapi.ApplicationAlertConfig) string {
	return data.ID
}

// SetID sets the ID in the API data object
func (r *applicationAlertConfigResourceFrameworkImpl) SetID(data *restapi.ApplicationAlertConfig, id string) {
	data.ID = id
}

// ============================================================================
// State to API Mapping
// ============================================================================

// MapStateToDataObject converts Terraform state to API data object
func (r *applicationAlertConfigResourceFrameworkImpl) MapStateToDataObject(ctx context.Context, state tfsdk.State) (*restapi.ApplicationAlertConfig, diag.Diagnostics) {
	var model ApplicationAlertConfigModel
	// Extract model from state
	diags := state.Get(ctx, &model)
	if diags.HasError() {
		return nil, diags
	}

	// Initialize result with basic fields
	result := &restapi.ApplicationAlertConfig{
		ID:               model.ID.ValueString(),
		Name:             model.Name.ValueString(),
		Description:      model.Description.ValueString(),
		BoundaryScope:    restapi.BoundaryScope(model.BoundaryScope.ValueString()),
		EvaluationType:   restapi.ApplicationAlertEvaluationType(model.EvaluationType.ValueString()),
		IncludeInternal:  model.IncludeInternal.ValueBool(),
		IncludeSynthetic: model.IncludeSynthetic.ValueBool(),
		Triggering:       model.Triggering.ValueBool(),
		GracePeriod:      extractGracePeriod(model.GracePeriod),
	}
	// Map granularity if present
	if !model.Granularity.IsNull() && !model.Granularity.IsUnknown() {
		result.Granularity = restapi.Granularity(model.Granularity.ValueInt64())
	}

	// Map optional and complex fields
	if err := r.mapTagFilter(&model, result, &diags); err != nil {
		return nil, diags
	}

	r.mapAlertChannels(ctx, &model, result, &diags)
	if diags.HasError() {
		return nil, diags
	}

	r.mapApplications(&model, result)

	if err := r.mapCustomPayloadFields(ctx, &model, result, &diags); err != nil {
		return nil, diags
	}

	if err := r.mapRules(&model, result, &diags); err != nil {
		return nil, diags
	}

	r.mapTimeThreshold(&model, result)

	return result, diags
}

// mapTagFilter parses and maps tag filter expression
func (r *applicationAlertConfigResourceFrameworkImpl) mapTagFilter(model *ApplicationAlertConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if model.TagFilter.IsNull() || model.TagFilter.IsUnknown() {
		return nil
	}

	tagFilterExpression := model.TagFilter.ValueString()
	if len(tagFilterExpression) == 0 {
		return nil
	}

	parsedExpression, err := tagfilter.ParseExpression(tagFilterExpression)
	if err != nil {
		diags.AddError(
			ErrorMessageFailedToParseTagFilter,
			fmt.Sprintf(ErrorMessageInvalidTagFilter, err.Error()),
		)
		return err
	}
	result.TagFilterExpression = parsedExpression
	return nil
}

// mapAlertChannels converts alert channels map from state to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapAlertChannels(ctx context.Context, model *ApplicationAlertConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) {
	if model.AlertChannels.IsNull() || model.AlertChannels.IsUnknown() {
		return
	}

	alertChannels := make(map[string][]string, len(model.AlertChannels.Elements()))
	for severity, channelSet := range model.AlertChannels.Elements() {
		var channelIDs []string
		*diags = channelSet.(types.Set).ElementsAs(ctx, &channelIDs, false)
		if diags.HasError() {
			return
		}
		alertChannels[severity] = channelIDs
	}
	result.AlertChannels = alertChannels
}

// mapApplications converts application scope configuration
func (r *applicationAlertConfigResourceFrameworkImpl) mapApplications(model *ApplicationAlertConfigModel, result *restapi.ApplicationAlertConfig) {
	// Always initialize as empty map, never nil
	result.Applications = make(map[string]restapi.IncludedApplication)

	if len(model.Applications) == 0 {
		return
	}

	for _, app := range model.Applications {
		appID := app.ApplicationID.ValueString()
		result.Applications[appID] = restapi.IncludedApplication{
			ApplicationID: appID,
			Inclusive:     app.Inclusive.ValueBool(),
			Services:      r.mapServices(app.Services),
		}
	}
}

// mapServices converts service scope configuration
func (r *applicationAlertConfigResourceFrameworkImpl) mapServices(services []ServiceModel) map[string]restapi.IncludedService {
	if len(services) == 0 {
		return nil
	}

	result := make(map[string]restapi.IncludedService, len(services))
	for _, svc := range services {
		svcID := svc.ServiceID.ValueString()
		result[svcID] = restapi.IncludedService{
			ServiceID: svcID,
			Inclusive: svc.Inclusive.ValueBool(),
			Endpoints: r.mapEndpoints(svc.Endpoints),
		}
	}
	return result
}

// mapEndpoints converts endpoint scope configuration
func (r *applicationAlertConfigResourceFrameworkImpl) mapEndpoints(endpoints []EndpointModel) map[string]restapi.IncludedEndpoint {
	if len(endpoints) == 0 {
		return nil
	}

	result := make(map[string]restapi.IncludedEndpoint, len(endpoints))
	for _, ep := range endpoints {
		epID := ep.EndpointID.ValueString()
		result[epID] = restapi.IncludedEndpoint{
			EndpointID: epID,
			Inclusive:  ep.Inclusive.ValueBool(),
		}
	}
	return result
}

// mapCustomPayloadFields converts custom payload fields from state to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapCustomPayloadFields(ctx context.Context, model *ApplicationAlertConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if model.CustomPayloadFields.IsNull() || model.CustomPayloadFields.IsUnknown() {
		return nil
	}

	customerPayloadFields, payloadDiags := shared.MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
	if payloadDiags.HasError() {
		diags.Append(payloadDiags...)
		return errors.New(ErrorMessageFailedToMapCustomPayload)
	}
	result.CustomerPayloadFields = customerPayloadFields
	return nil
}

// mapRules converts alert rules with thresholds from state to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapRules(model *ApplicationAlertConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if len(model.Rules) == 0 {
		return nil
	}

	result.Rules = make([]restapi.ApplicationAlertRuleWithThresholds, len(model.Rules))
	for i, ruleWithThreshold := range model.Rules {
		if err := r.mapSingleRule(i, &ruleWithThreshold, result, diags); err != nil {
			return err
		}

		// Validate that at least one threshold is defined for each rule
		if err := r.validateRuleThresholds(i, &ruleWithThreshold, diags); err != nil {
			return err
		}
	}
	return nil
}

// validateRuleThresholds validates that at least one threshold is defined for a rule
func (r *applicationAlertConfigResourceFrameworkImpl) validateRuleThresholds(index int, ruleWithThreshold *RuleWithThresholdModel, diags *diag.Diagnostics) error {
	if ruleWithThreshold.Thresholds == nil {
		diags.AddError(
			ErrorMessageValidationError,
			fmt.Sprintf("Rule at index %d must have at least one threshold defined (warning or critical)", index),
		)
		return errors.New("missing threshold definition")
	}

	// Check if at least one threshold level is defined
	hasWarning := ruleWithThreshold.Thresholds.Warning != nil &&
		(ruleWithThreshold.Thresholds.Warning.Static != nil ||
			ruleWithThreshold.Thresholds.Warning.AdaptiveBaseline != nil ||
			ruleWithThreshold.Thresholds.Warning.HistoricBaseline != nil)
	hasCritical := ruleWithThreshold.Thresholds.Critical != nil &&
		(ruleWithThreshold.Thresholds.Critical.Static != nil ||
			ruleWithThreshold.Thresholds.Critical.AdaptiveBaseline != nil ||
			ruleWithThreshold.Thresholds.Critical.HistoricBaseline != nil)

	if !hasWarning && !hasCritical {
		diags.AddError(
			ErrorMessageValidationError,
			fmt.Sprintf("Rule at index %d must have at least one threshold defined (warning or critical) with either static, adaptive_baseline, or historic_baseline configuration", index),
		)
		return errors.New("no valid threshold configuration found")
	}

	return nil
}

// mapSingleRule converts a single rule with its thresholds
func (r *applicationAlertConfigResourceFrameworkImpl) mapSingleRule(index int, ruleWithThreshold *RuleWithThresholdModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	result.Rules[index] = restapi.ApplicationAlertRuleWithThresholds{
		ThresholdOperator: ruleWithThreshold.ThresholdOperator.ValueString(),
	}

	if ruleWithThreshold.Rule != nil {
		if err := r.mapRuleConfiguration(index, ruleWithThreshold.Rule, result, diags); err != nil {
			return err
		}
	}

	if ruleWithThreshold.Thresholds != nil {
		r.mapThresholds(index, ruleWithThreshold.Thresholds, result)
	}

	return nil
}

// mapRuleConfiguration maps the rule type and its specific configuration
func (r *applicationAlertConfigResourceFrameworkImpl) mapRuleConfiguration(index int, rule *RuleModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	result.Rules[index].Rule = &restapi.ApplicationAlertRule{}

	// Map each rule type
	if rule.ErrorRate != nil {
		return r.mapErrorRateRule(index, rule.ErrorRate, result, diags)
	}
	if rule.Errors != nil {
		return r.mapErrorsRule(index, rule.Errors, result, diags)
	}
	if rule.Logs != nil {
		return r.mapLogsRule(index, rule.Logs, result, diags)
	}
	if rule.Slowness != nil {
		return r.mapSlownessRule(index, rule.Slowness, result, diags)
	}
	if rule.StatusCode != nil {
		return r.mapStatusCodeRule(index, rule.StatusCode, result, diags)
	}
	if rule.Throughput != nil {
		return r.mapThroughputRule(index, rule.Throughput, result, diags)
	}

	return nil
}

// mapErrorRateRule maps error rate rule configuration to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapErrorRateRule(index int, errorRate *RuleConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if errorRate.MetricName.IsNull() || errorRate.MetricName.IsUnknown() {
		diags.AddError(ErrorMessageValidationError, fmt.Sprintf(ErrorMessageMetricNameRequired, "error rate"))
		return errors.New(ErrorMessageMissingMetricName)
	}
	result.Rules[index].Rule.AlertType = APIAlertTypeErrorRate
	result.Rules[index].Rule.MetricName = errorRate.MetricName.ValueString()
	result.Rules[index].Rule.Aggregation = restapi.Aggregation(errorRate.Aggregation.ValueString())
	return nil
}

// mapErrorsRule maps errors rule configuration to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapErrorsRule(index int, errorsRule *RuleConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if errorsRule.MetricName.IsNull() || errorsRule.MetricName.IsUnknown() {
		diags.AddError(ErrorMessageValidationError, fmt.Sprintf(ErrorMessageMetricNameRequired, "error"))
		return errors.New(ErrorMessageMissingMetricName)
	}
	result.Rules[index].Rule.AlertType = APIAlertTypeErrors
	result.Rules[index].Rule.MetricName = errorsRule.MetricName.ValueString()
	result.Rules[index].Rule.Aggregation = restapi.Aggregation(errorsRule.Aggregation.ValueString())
	return nil
}

// mapLogsRule maps logs rule configuration with additional fields to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapLogsRule(index int, logs *LogsRuleModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if logs.MetricName.IsNull() || logs.MetricName.IsUnknown() ||
		logs.Level.IsNull() || logs.Level.IsUnknown() ||
		logs.Operator.IsNull() || logs.Operator.IsUnknown() {
		diags.AddError(ErrorMessageValidationError, ErrorMessageLogsFieldsRequired)
		return errors.New(ErrorMessageMissingRequiredFields)
	}

	result.Rules[index].Rule.AlertType = APIAlertTypeLogs
	result.Rules[index].Rule.MetricName = logs.MetricName.ValueString()
	result.Rules[index].Rule.Aggregation = restapi.Aggregation(logs.Aggregation.ValueString())

	level := restapi.LogLevel(logs.Level.ValueString())
	result.Rules[index].Rule.Level = &level

	message := logs.Message.ValueString()
	result.Rules[index].Rule.Message = &message

	operator := restapi.ExpressionOperator(logs.Operator.ValueString())
	result.Rules[index].Rule.Operator = &operator

	return nil
}

// mapSlownessRule maps slowness rule configuration to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapSlownessRule(index int, slowness *RuleConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if slowness.MetricName.IsNull() || slowness.MetricName.IsUnknown() {
		diags.AddError(ErrorMessageValidationError, fmt.Sprintf(ErrorMessageMetricNameRequired, "slowness"))
		return errors.New(ErrorMessageMissingMetricName)
	}
	result.Rules[index].Rule.AlertType = APIAlertTypeSlowness
	result.Rules[index].Rule.MetricName = slowness.MetricName.ValueString()
	result.Rules[index].Rule.Aggregation = restapi.Aggregation(slowness.Aggregation.ValueString())
	return nil
}

// mapStatusCodeRule maps status code rule configuration with range to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapStatusCodeRule(index int, statusCode *StatusCodeRuleModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if statusCode.MetricName.IsNull() || statusCode.MetricName.IsUnknown() {
		diags.AddError(ErrorMessageValidationError, fmt.Sprintf(ErrorMessageMetricNameRequired, "status code"))
		return errors.New(ErrorMessageMissingMetricName)
	}

	result.Rules[index].Rule.AlertType = APIAlertTypeStatusCode
	result.Rules[index].Rule.MetricName = statusCode.MetricName.ValueString()
	result.Rules[index].Rule.Aggregation = restapi.Aggregation(statusCode.Aggregation.ValueString())

	statusCodeStart := int32(statusCode.StatusCodeStart.ValueInt64())
	result.Rules[index].Rule.StatusCodeStart = &statusCodeStart

	statusCodeEnd := int32(statusCode.StatusCodeEnd.ValueInt64())
	result.Rules[index].Rule.StatusCodeEnd = &statusCodeEnd

	return nil
}

// mapThroughputRule maps throughput rule configuration to API format
func (r *applicationAlertConfigResourceFrameworkImpl) mapThroughputRule(index int, throughput *RuleConfigModel, result *restapi.ApplicationAlertConfig, diags *diag.Diagnostics) error {
	if throughput.MetricName.IsNull() || throughput.MetricName.IsUnknown() {
		diags.AddError(ErrorMessageValidationError, fmt.Sprintf(ErrorMessageMetricNameRequired, "throughput"))
		return errors.New(ErrorMessageMissingMetricName)
	}
	result.Rules[index].Rule.AlertType = APIAlertTypeThroughput
	result.Rules[index].Rule.MetricName = throughput.MetricName.ValueString()
	result.Rules[index].Rule.Aggregation = restapi.Aggregation(throughput.Aggregation.ValueString())
	return nil
}

// mapThresholds converts threshold configurations for warning and critical severity levels using shared mapping
func (r *applicationAlertConfigResourceFrameworkImpl) mapThresholds(index int, thresholds *ApplicationThresholdModel, result *restapi.ApplicationAlertConfig) {
	// Convert ApplicationThresholdModel to shared.ThresholdAllPluginModel
	sharedThresholds := &shared.ThresholdAllPluginModel{}

	if thresholds.Warning != nil {
		sharedThresholds.Warning = &shared.ThresholdAllTypeModel{
			Static:           thresholds.Warning.Static,
			AdaptiveBaseline: thresholds.Warning.AdaptiveBaseline,
			HistoricBaseline: thresholds.Warning.HistoricBaseline,
		}
	}

	if thresholds.Critical != nil {
		sharedThresholds.Critical = &shared.ThresholdAllTypeModel{
			Static:           thresholds.Critical.Static,
			AdaptiveBaseline: thresholds.Critical.AdaptiveBaseline,
			HistoricBaseline: thresholds.Critical.HistoricBaseline,
		}
	}

	// Use shared mapping function
	thresholdMap, _ := shared.MapThresholdsAllPluginFromState(nil, sharedThresholds)
	result.Rules[index].Thresholds = thresholdMap
}

// mapTimeThreshold converts time threshold configuration
func (r *applicationAlertConfigResourceFrameworkImpl) mapTimeThreshold(model *ApplicationAlertConfigModel, result *restapi.ApplicationAlertConfig) {
	if model.TimeThreshold == nil {
		return
	}

	result.TimeThreshold = &restapi.ApplicationAlertTimeThreshold{}

	if model.TimeThreshold.RequestImpact != nil {
		result.TimeThreshold.Type = TimeThresholdTypeRequestImpact
		result.TimeThreshold.TimeWindow = model.TimeThreshold.RequestImpact.TimeWindow.ValueInt64()
		result.TimeThreshold.Requests = int(model.TimeThreshold.RequestImpact.Requests.ValueInt64())
	} else if model.TimeThreshold.ViolationsInPeriod != nil {
		result.TimeThreshold.Type = TimeThresholdTypeViolationsInPeriod
		result.TimeThreshold.TimeWindow = model.TimeThreshold.ViolationsInPeriod.TimeWindow.ValueInt64()
		result.TimeThreshold.Violations = int(model.TimeThreshold.ViolationsInPeriod.Violations.ValueInt64())
	} else if model.TimeThreshold.ViolationsInSequence != nil {
		result.TimeThreshold.Type = TimeThresholdTypeViolationsInSequence
		result.TimeThreshold.TimeWindow = model.TimeThreshold.ViolationsInSequence.TimeWindow.ValueInt64()
	}
}

// extractGracePeriod converts types.Int64 to *int64 for API, handling null/unknown values
func extractGracePeriod(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil // send as null (omitted in JSON)
	}
	val := v.ValueInt64()
	return &val
}

// ============================================================================
// API to State Mapping
// ============================================================================

// UpdateState converts API data object to Terraform state
func (r *applicationAlertConfigResourceFrameworkImpl) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	var model ApplicationAlertConfigModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	} else {
		model = ApplicationAlertConfigModel{}
	}
	// Build base model with simple fields
	model.ID = types.StringValue(data.ID)
	model.Name = types.StringValue(data.Name)
	model.Description = types.StringValue(data.Description)
	model.BoundaryScope = types.StringValue(string(data.BoundaryScope))
	model.EvaluationType = types.StringValue(string(data.EvaluationType))
	model.Granularity = types.Int64Value(int64(data.Granularity))
	model.IncludeInternal = types.BoolValue(data.IncludeInternal)
	model.IncludeSynthetic = types.BoolValue(data.IncludeSynthetic)
	model.Triggering = types.BoolValue(data.Triggering)

	// Map complex fields with error handling
	if err := r.updateGracePeriod(&model, data); err != nil {
		diags.Append(err...)
		return diags
	}

	if model.TagFilter.IsNull() || model.TagFilter.IsUnknown() {
		if err := r.updateTagFilter(&model, data); err != nil {
			diags.Append(err...)
			return diags
		}
	}

	if err := r.updateAlertChannels(ctx, &model, data); err != nil {
		diags.Append(err...)
		return diags
	}

	if err := r.updateApplications(&model, data); err != nil {
		diags.Append(err...)
		return diags
	}

	if err := r.updateCustomPayloadFields(ctx, &model, data); err != nil {
		diags.Append(err...)
		return diags
	}

	if err := r.updateRules(&model, data); err != nil {
		diags.Append(err...)
		return diags
	}

	if err := r.updateTimeThreshold(&model, data); err != nil {
		diags.Append(err...)
		return diags
	}

	// Set final state
	return state.Set(ctx, model)
}

// buildBaseModel creates the base model with simple field mappings
func (r *applicationAlertConfigResourceFrameworkImpl) buildBaseModel(data *restapi.ApplicationAlertConfig) ApplicationAlertConfigModel {
	return ApplicationAlertConfigModel{
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
}

// updateGracePeriod handles grace period field mapping
func (r *applicationAlertConfigResourceFrameworkImpl) updateGracePeriod(model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	if data.GracePeriod != nil {
		model.GracePeriod = util.SetInt64PointerToState(data.GracePeriod)
	} else {
		model.GracePeriod = types.Int64Null()
	}
	return nil
}

// updateTagFilter handles tag filter normalization and mapping
func (r *applicationAlertConfigResourceFrameworkImpl) updateTagFilter(model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(data.TagFilterExpression)
		if err != nil {
			diags.AddError(
				ErrorMessageTagFilterNormalizationError,
				fmt.Sprintf(ErrorMessageFailedToNormalizeTagFilter, err.Error()),
			)
			return diags
		}
		model.TagFilter = util.SetStringPointerToState(normalizedTagFilterString)
	} else {
		model.TagFilter = types.StringNull()
	}

	return diags
}

// updateAlertChannels handles alert channels mapping from API to state
func (r *applicationAlertConfigResourceFrameworkImpl) updateAlertChannels(ctx context.Context, model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.AlertChannels == nil {
		model.AlertChannels = types.MapNull(types.SetType{ElemType: types.StringType})
		return diags
	}

	if len(data.AlertChannels) == 0 {
		emptyMap, mapDiags := types.MapValue(types.SetType{ElemType: types.StringType}, map[string]attr.Value{})
		if mapDiags.HasError() {
			diags.Append(mapDiags...)
			return diags
		}
		model.AlertChannels = emptyMap
		return diags
	}

	elements := make(map[string]attr.Value, len(data.AlertChannels))
	for severity, channelIDs := range data.AlertChannels {
		channelElements := make([]attr.Value, len(channelIDs))
		for i, id := range channelIDs {
			channelElements[i] = types.StringValue(id)
		}

		channelSet, setDiags := types.SetValue(types.StringType, channelElements)
		if setDiags.HasError() {
			diags.Append(setDiags...)
			return diags
		}
		elements[severity] = channelSet
	}

	alertChannelsMap, mapDiags := types.MapValue(types.SetType{ElemType: types.StringType}, elements)
	if mapDiags.HasError() {
		diags.Append(mapDiags...)
		return diags
	}

	model.AlertChannels = alertChannelsMap
	return diags
}

// updateApplications handles application scope mapping from API to state
func (r *applicationAlertConfigResourceFrameworkImpl) updateApplications(model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	if len(data.Applications) == 0 {
		model.Applications = []ApplicationModel{}
		return nil
	}

	model.Applications = make([]ApplicationModel, 0, len(data.Applications))
	for _, app := range data.Applications {
		appModel := ApplicationModel{
			ApplicationID: types.StringValue(app.ApplicationID),
			Inclusive:     types.BoolValue(app.Inclusive),
			Services:      r.mapServicesToModel(app.Services),
		}
		model.Applications = append(model.Applications, appModel)
	}

	return nil
}

// mapServicesToModel converts services from API to model format
func (r *applicationAlertConfigResourceFrameworkImpl) mapServicesToModel(services map[string]restapi.IncludedService) []ServiceModel {
	if len(services) == 0 {
		return []ServiceModel{}
	}

	serviceModels := make([]ServiceModel, 0, len(services))
	for _, svc := range services {
		svcModel := ServiceModel{
			ServiceID: types.StringValue(svc.ServiceID),
			Inclusive: types.BoolValue(svc.Inclusive),
			Endpoints: r.mapEndpointsToModel(svc.Endpoints),
		}
		serviceModels = append(serviceModels, svcModel)
	}

	return serviceModels
}

// mapEndpointsToModel converts endpoints from API to model format
func (r *applicationAlertConfigResourceFrameworkImpl) mapEndpointsToModel(endpoints map[string]restapi.IncludedEndpoint) []EndpointModel {
	if len(endpoints) == 0 {
		return []EndpointModel{}
	}

	endpointModels := make([]EndpointModel, 0, len(endpoints))
	for _, ep := range endpoints {
		epModel := EndpointModel{
			EndpointID: types.StringValue(ep.EndpointID),
			Inclusive:  types.BoolValue(ep.Inclusive),
		}
		endpointModels = append(endpointModels, epModel)
	}

	return endpointModels
}

// updateCustomPayloadFields handles custom payload fields mapping from API to state
func (r *applicationAlertConfigResourceFrameworkImpl) updateCustomPayloadFields(ctx context.Context, model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	customPayloadFieldsList, payloadDiags := shared.CustomPayloadFieldsToTerraform(ctx, data.CustomerPayloadFields)
	if payloadDiags.HasError() {
		return payloadDiags
	}
	model.CustomPayloadFields = customPayloadFieldsList
	return nil
}

// updateRules handles rules mapping with thresholds from API to state
func (r *applicationAlertConfigResourceFrameworkImpl) updateRules(model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(data.Rules) == 0 {
		model.Rules = []RuleWithThresholdModel{}
		return diags
	}

	ruleModels := make([]RuleWithThresholdModel, len(data.Rules))
	for i, ruleWithThreshold := range data.Rules {
		ruleModel, err := r.mapRuleToModel(&ruleWithThreshold)
		if err != nil {
			diags.Append(err...)
			return diags
		}
		ruleModels[i] = ruleModel
	}

	model.Rules = ruleModels
	return diags
}

// mapRuleToModel converts a single rule with thresholds to model format
func (r *applicationAlertConfigResourceFrameworkImpl) mapRuleToModel(ruleWithThreshold *restapi.ApplicationAlertRuleWithThresholds) (RuleWithThresholdModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	ruleModel := r.createRuleModelByType(ruleWithThreshold.Rule)

	ruleWithThresholdModel := RuleWithThresholdModel{
		Rule:              &ruleModel,
		ThresholdOperator: types.StringValue(ruleWithThreshold.ThresholdOperator),
		Thresholds:        r.mapThresholdsToModel(ruleWithThreshold.Thresholds),
	}

	return ruleWithThresholdModel, diags
}

// createRuleModelByType creates appropriate rule model based on alert type
func (r *applicationAlertConfigResourceFrameworkImpl) createRuleModelByType(rule *restapi.ApplicationAlertRule) RuleModel {
	ruleModel := RuleModel{}

	switch rule.AlertType {
	case APIAlertTypeErrorRate:
		ruleModel.ErrorRate = &RuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(rule.Aggregation)),
		}

	case APIAlertTypeErrors:
		ruleModel.Errors = &RuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(rule.Aggregation)),
		}

	case APIAlertTypeLogs:
		ruleModel.Logs = &LogsRuleModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(rule.Aggregation)),
			Level:       util.SetStringPointerToState((*string)(rule.Level)),
			Message:     util.SetStringPointerToState(rule.Message),
			Operator:    util.SetStringPointerToState((*string)(rule.Operator)),
		}

	case APIAlertTypeSlowness:
		ruleModel.Slowness = &RuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(rule.Aggregation)),
		}

	case APIAlertTypeStatusCode:
		ruleModel.StatusCode = &StatusCodeRuleModel{
			MetricName:      types.StringValue(rule.MetricName),
			Aggregation:     types.StringValue(string(rule.Aggregation)),
			StatusCodeStart: util.SetInt32PointerToInt64State(rule.StatusCodeStart),
			StatusCodeEnd:   util.SetInt32PointerToInt64State(rule.StatusCodeEnd),
		}

	case APIAlertTypeThroughput:
		ruleModel.Throughput = &RuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(rule.Aggregation)),
		}
	}

	return ruleModel
}

// mapThresholdsToModel converts threshold map to model format using shared mapping
func (r *applicationAlertConfigResourceFrameworkImpl) mapThresholdsToModel(thresholds map[restapi.AlertSeverity]restapi.ThresholdRule) *ApplicationThresholdModel {
	if len(thresholds) == 0 {
		return nil
	}

	thresholdModel := &ApplicationThresholdModel{}

	if warningThreshold, ok := thresholds[restapi.WarningSeverity]; ok {
		// Use shared mapping function
		sharedWarning := shared.MapAllThresholdPluginToState(nil, &warningThreshold, true)
		if sharedWarning != nil {
			thresholdModel.Warning = &ThresholdLevelModel{
				Static:           sharedWarning.Static,
				AdaptiveBaseline: sharedWarning.AdaptiveBaseline,
				HistoricBaseline: sharedWarning.HistoricBaseline,
			}
		}
	}

	if criticalThreshold, ok := thresholds[restapi.CriticalSeverity]; ok {
		// Use shared mapping function
		sharedCritical := shared.MapAllThresholdPluginToState(nil, &criticalThreshold, true)
		if sharedCritical != nil {
			thresholdModel.Critical = &ThresholdLevelModel{
				Static:           sharedCritical.Static,
				AdaptiveBaseline: sharedCritical.AdaptiveBaseline,
				HistoricBaseline: sharedCritical.HistoricBaseline,
			}
		}
	}

	return thresholdModel
}

// updateTimeThreshold handles time threshold mapping from API to state
func (r *applicationAlertConfigResourceFrameworkImpl) updateTimeThreshold(model *ApplicationAlertConfigModel, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
	if data.TimeThreshold == nil {
		model.TimeThreshold = nil
		return nil
	}

	timeThresholdModel := &AppAlertTimeThresholdModel{}

	switch data.TimeThreshold.Type {
	case TimeThresholdTypeRequestImpact:
		timeThresholdModel.RequestImpact = &AppAlertRequestImpactModel{
			TimeWindow: types.Int64Value(data.TimeThreshold.TimeWindow),
			Requests:   types.Int64Value(int64(data.TimeThreshold.Requests)),
		}

	case TimeThresholdTypeViolationsInPeriod:
		timeThresholdModel.ViolationsInPeriod = &AppAlertViolationsInPeriodModel{
			TimeWindow: types.Int64Value(data.TimeThreshold.TimeWindow),
			Violations: types.Int64Value(int64(data.TimeThreshold.Violations)),
		}

	case TimeThresholdTypeViolationsInSequence:
		timeThresholdModel.ViolationsInSequence = &AppAlertViolationsInSequenceModel{
			TimeWindow: types.Int64Value(data.TimeThreshold.TimeWindow),
		}
	}

	model.TimeThreshold = timeThresholdModel
	return nil
}
