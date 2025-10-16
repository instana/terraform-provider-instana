package instana

import (
	"context"
	"fmt"

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

// ResourceInstanaApplicationAlertConfigFramework the name of the terraform-provider-instana resource to manage application alert configs
const ResourceInstanaApplicationAlertConfigFramework = "application_alert_config"

// ResourceInstanaGlobalApplicationAlertConfigFramework the name of the terraform-provider-instana resource to manage global application alert configs
const ResourceInstanaGlobalApplicationAlertConfigFramework = "global_application_alert_config"

// ApplicationAlertConfigModel represents the data model for the application alert configuration resource
type ApplicationAlertConfigModel struct {
	ID                  types.String `tfsdk:"id"`
	AlertChannelIDs     types.Set    `tfsdk:"alert_channel_ids"`
	Applications        types.Set    `tfsdk:"application"`
	BoundaryScope       types.String `tfsdk:"boundary_scope"`
	CustomPayloadFields types.List   `tfsdk:"custom_payload_field"`
	Description         types.String `tfsdk:"description"`
	EvaluationType      types.String `tfsdk:"evaluation_type"`
	Granularity         types.Int64  `tfsdk:"granularity"`
	IncludeInternal     types.Bool   `tfsdk:"include_internal"`
	IncludeSynthetic    types.Bool   `tfsdk:"include_synthetic"`
	Name                types.String `tfsdk:"name"`
	Rule                types.List   `tfsdk:"rule"`
	Severity            types.String `tfsdk:"severity"`
	TagFilter           types.String `tfsdk:"tag_filter"`
	Threshold           types.Object `tfsdk:"threshold"`
	TimeThreshold       types.List   `tfsdk:"time_threshold"`
	Triggering          types.Bool   `tfsdk:"triggering"`
}

// ApplicationModel represents an application in the application alert config
type ApplicationModel struct {
	ApplicationID types.String `tfsdk:"application_id"`
	Inclusive     types.Bool   `tfsdk:"inclusive"`
	Services      types.Set    `tfsdk:"service"`
}

// ServiceModel represents a service in the application alert config
type ServiceModel struct {
	ServiceID types.String `tfsdk:"service_id"`
	Inclusive types.Bool   `tfsdk:"inclusive"`
	Endpoints types.Set    `tfsdk:"endpoint"`
}

// EndpointModel represents an endpoint in the application alert config
type EndpointModel struct {
	EndpointID types.String `tfsdk:"endpoint_id"`
	Inclusive  types.Bool   `tfsdk:"inclusive"`
}

// RuleModel represents a rule in the application alert config
type RuleModel struct {
	ErrorRate  types.List `tfsdk:"error_rate"`
	Errors     types.List `tfsdk:"errors"`
	Logs       types.List `tfsdk:"logs"`
	Slowness   types.List `tfsdk:"slowness"`
	StatusCode types.List `tfsdk:"status_code"`
	Throughput types.List `tfsdk:"throughput"`
}

// RuleConfigModel represents the common configuration for rules
type RuleConfigModel struct {
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
}

// LogsRuleModel represents the logs rule configuration
type LogsRuleModel struct {
	MetricName  types.String `tfsdk:"metric_name"`
	Aggregation types.String `tfsdk:"aggregation"`
	Level       types.String `tfsdk:"level"`
	Message     types.String `tfsdk:"message"`
	Operator    types.String `tfsdk:"operator"`
}

// StatusCodeRuleModel represents the status code rule configuration
type StatusCodeRuleModel struct {
	MetricName      types.String `tfsdk:"metric_name"`
	Aggregation     types.String `tfsdk:"aggregation"`
	StatusCodeStart types.Int64  `tfsdk:"status_code_start"`
	StatusCodeEnd   types.Int64  `tfsdk:"status_code_end"`
}

type AppAlertTimeThresholdModel struct {
	RequestImpact        types.List `tfsdk:"request_impact"`
	ViolationsInPeriod   types.List `tfsdk:"violations_in_period"`
	ViolationsInSequence types.List `tfsdk:"violations_in_sequence"`
}

// AppAlertRequestImpactModel represents the request impact time threshold configuration
type AppAlertRequestImpactModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
	Requests   types.Int64 `tfsdk:"requests"`
}

// AppAlertViolationsInPeriodModel represents the violations in period time threshold configuration
type AppAlertViolationsInPeriodModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
	Violations types.Int64 `tfsdk:"violations"`
}

// AppAlertViolationsInSequenceModel represents the violations in sequence time threshold configuration
type AppAlertViolationsInSequenceModel struct {
	TimeWindow types.Int64 `tfsdk:"time_window"`
}

// NewApplicationAlertConfigResourceHandleFramework creates the resource handle for Application Alert Configuration
func NewApplicationAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.ApplicationAlertConfig] {
	return &applicationAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
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
						Description: "List of IDs of alert channels defined in Instana.",
						ElementType: types.StringType,
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
						Required:    true,
						Description: "The severity of the alert when triggered",
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
				},
				Blocks: map[string]schema.Block{
					ApplicationAlertConfigFieldApplications: schema.SetNestedBlock{
						Description: "Selection of applications in scope.",
						NestedObject: schema.NestedBlockObject{
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
							},
							Blocks: map[string]schema.Block{
								ApplicationAlertConfigFieldApplicationsServices: schema.SetNestedBlock{
									Description: "Selection of services in scope.",
									NestedObject: schema.NestedBlockObject{
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
										},
										Blocks: map[string]schema.Block{
											ApplicationAlertConfigFieldApplicationsServicesEndpoints: schema.SetNestedBlock{
												Description: "Selection of endpoints in scope.",
												NestedObject: schema.NestedBlockObject{
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
					DefaultCustomPayloadFieldsName: GetCustomPayloadFieldsSchema(),
					ApplicationAlertConfigFieldRule: schema.ListNestedBlock{
						Description: "Indicates the type of rule this alert configuration is about.",
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								ApplicationAlertConfigFieldRuleErrorRate: schema.ListNestedBlock{
									Description: "Rule based on the error rate of the configured alert configuration target",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
												Required:    true,
												Description: "The metric name of the application alert rule",
											},
											ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
												Optional:    true,
												Description: "The aggregation function of the application alert rule",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
												},
											},
										},
									},
								},
								ApplicationAlertConfigFieldRuleErrors: schema.ListNestedBlock{
									Description: "Rule based on the number of errors of the configured alert configuration target",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
												Required:    true,
												Description: "The metric name of the application alert rule",
											},
											ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
												Optional:    true,
												Description: "The aggregation function of the application alert rule",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
												},
											},
										},
									},
								},
								ApplicationAlertConfigFieldRuleLogs: schema.ListNestedBlock{
									Description: "Rule based on logs of the configured alert configuration target",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
												Required:    true,
												Description: "The metric name of the application alert rule",
											},
											ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
												Optional:    true,
												Description: "The aggregation function of the application alert rule",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
												},
											},
											ApplicationAlertConfigFieldRuleLogsLevel: schema.StringAttribute{
												Required:    true,
												Description: "The log level for which this rule applies to",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedLogLevels.ToStringSlice()...),
												},
											},
											ApplicationAlertConfigFieldRuleLogsMessage: schema.StringAttribute{
												Optional:    true,
												Description: "The log message for which this rule applies to",
											},
											ApplicationAlertConfigFieldRuleLogsOperator: schema.StringAttribute{
												Required:    true,
												Description: "The operator which will be applied to evaluate this rule",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedExpressionOperators.ToStringSlice()...),
												},
											},
										},
									},
								},
								ApplicationAlertConfigFieldRuleSlowness: schema.ListNestedBlock{
									Description: "Rule based on the slowness of the configured alert configuration target",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
												Required:    true,
												Description: "The metric name of the application alert rule",
											},
											ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
												Required:    true,
												Description: "The aggregation function of the application alert rule",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
												},
											},
										},
									},
								},
								ApplicationAlertConfigFieldRuleStatusCode: schema.ListNestedBlock{
									Description: "Rule based on the HTTP status code of the configured alert configuration target",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
												Required:    true,
												Description: "The metric name of the application alert rule",
											},
											ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
												Optional:    true,
												Description: "The aggregation function of the application alert rule",
												Validators: []validator.String{
													stringvalidator.OneOf(restapi.SupportedAggregations.ToStringSlice()...),
												},
											},
											ApplicationAlertConfigFieldRuleStatusCodeStart: schema.Int64Attribute{
												Optional:    true,
												Description: "minimal HTTP status code applied for this rule",
											},
											ApplicationAlertConfigFieldRuleStatusCodeEnd: schema.Int64Attribute{
												Optional:    true,
												Description: "maximum HTTP status code applied for this rule",
											},
										},
									},
								},
								ApplicationAlertConfigFieldRuleThroughput: schema.ListNestedBlock{
									Description: "Rule based on the throughput of the configured alert configuration target",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldRuleMetricName: schema.StringAttribute{
												Required:    true,
												Description: "The metric name of the application alert rule",
											},
											ApplicationAlertConfigFieldRuleAggregation: schema.StringAttribute{
												Optional:    true,
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
					ResourceFieldThreshold: schema.SingleNestedBlock{
						Description: "The threshold configuration for the alert",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The type of the threshold",
								Validators: []validator.String{
									stringvalidator.OneOf("upperBound", "lowerBound"),
								},
							},
							"value": schema.Float64Attribute{
								Required:    true,
								Description: "The value of the threshold",
							},
						},
					},
					ApplicationAlertConfigFieldTimeThreshold: schema.ListNestedBlock{
						Description: "Indicates the type of violation of the defined threshold.",
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								ApplicationAlertConfigFieldTimeThresholdRequestImpact: schema.ListNestedBlock{
									Description: "Time threshold base on request impact",
									NestedObject: schema.NestedBlockObject{
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
								},
								ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: schema.ListNestedBlock{
									Description: "Time threshold base on violations in period",
									NestedObject: schema.NestedBlockObject{
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
								},
								ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: schema.ListNestedBlock{
									Description: "Time threshold base on violations in sequence",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											ApplicationAlertConfigFieldTimeThresholdTimeWindow: schema.Int64Attribute{
												Required:    true,
												Description: "The time window if the time threshold",
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

// NewGlobalApplicationAlertConfigResourceHandleFramework creates the resource handle for Global Application Alert Configuration
func NewGlobalApplicationAlertConfigResourceHandleFramework() ResourceHandleFramework[*restapi.ApplicationAlertConfig] {
	return &applicationAlertConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaGlobalApplicationAlertConfigFramework,
			Schema:        NewApplicationAlertConfigResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 1,
		},
		isGlobal: true,
	}
}

type applicationAlertConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
	isGlobal bool
}

func (r *applicationAlertConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *applicationAlertConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig] {
	if r.isGlobal {
		return api.GlobalApplicationAlertConfigs()
	}
	return api.ApplicationAlertConfigs()
}

func (r *applicationAlertConfigResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *applicationAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, config *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the config
	model := ApplicationAlertConfigModel{
		ID:               types.StringValue(config.ID),
		Name:             types.StringValue(config.Name),
		Description:      types.StringValue(config.Description),
		BoundaryScope:    types.StringValue(string(config.BoundaryScope)),
		EvaluationType:   types.StringValue(string(config.EvaluationType)),
		Granularity:      types.Int64Value(int64(config.Granularity)),
		IncludeInternal:  types.BoolValue(config.IncludeInternal),
		IncludeSynthetic: types.BoolValue(config.IncludeSynthetic),
		Triggering:       types.BoolValue(config.Triggering),
	}

	// Convert severity from Instana API to Terraform representation
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(config.Severity)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Error converting severity",
				fmt.Sprintf("Failed to convert severity: %s", err),
			),
		}
	}
	model.Severity = types.StringValue(severity)

	// Set alert channel IDs
	alertChannelIDs, diags := types.SetValueFrom(ctx, types.StringType, config.AlertChannelIDs)
	if diags.HasError() {
		return diags
	}
	model.AlertChannelIDs = alertChannelIDs

	// Set tag filter
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression)
		if err != nil {
			return diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Error converting tag filter",
					fmt.Sprintf("Failed to convert tag filter: %s", err),
				),
			}
		}
		if normalizedTagFilterString != nil {
			model.TagFilter = types.StringValue(*normalizedTagFilterString)
		} else {
			model.TagFilter = types.StringNull()
		}
	} else {
		model.TagFilter = types.StringNull()
	}

	// Convert applications to Terraform types
	applications, appDiags := r.mapApplicationsToTerraform(ctx, config.Applications)
	if appDiags.HasError() {
		return appDiags
	}
	model.Applications = applications

	// Convert rule to Terraform types
	rule, ruleDiags := r.mapRuleToTerraform(ctx, &config.Rule)
	if ruleDiags.HasError() {
		return ruleDiags
	}
	model.Rule = rule

	// Convert threshold to Terraform types
	threshold, thresholdDiags := r.mapThresholdToTerraform(ctx, &config.Threshold)
	if thresholdDiags.HasError() {
		return thresholdDiags
	}
	model.Threshold = threshold

	// Convert time threshold to Terraform types
	timeThreshold, timeThresholdDiags := r.mapTimeThresholdToTerraform(ctx, &config.TimeThreshold)
	if timeThresholdDiags.HasError() {
		return timeThresholdDiags
	}
	model.TimeThreshold = timeThreshold

	// Convert custom payload fields to Terraform types
	customPayloadFields, payloadDiags := CustomPayloadFieldsToTerraform(ctx, config.CustomerPayloadFields)
	if payloadDiags.HasError() {
		return payloadDiags
	}
	model.CustomPayloadFields = customPayloadFields

	// Set the entire model to state
	diags = state.Set(ctx, model)
	return diags
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

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map severity
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(model.Severity.ValueString())
	if err != nil {
		diags.AddError(
			"Error converting severity",
			fmt.Sprintf("Failed to convert severity: %s", err),
		)
		return nil, diags
	}

	// Map tag filter
	var tagFilter *restapi.TagFilter
	if !model.TagFilter.IsNull() {
		tagFilterStr := model.TagFilter.ValueString()
		parser := tagfilter.NewParser()
		expr, err := parser.Parse(tagFilterStr)
		if err != nil {
			diags.AddError(
				"Error parsing tag filter",
				fmt.Sprintf("Failed to parse tag filter: %s", err),
			)
			return nil, diags
		}

		mapper := tagfilter.NewMapper()
		tagFilter = mapper.ToAPIModel(expr)
	}

	// Map alert channel IDs
	var alertChannelIDs []string
	if !model.AlertChannelIDs.IsNull() {
		diags.Append(model.AlertChannelIDs.ElementsAs(ctx, &alertChannelIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map applications
	applications, appDiags := r.mapApplicationsFromTerraform(ctx, model.Applications)
	if appDiags.HasError() {
		diags.Append(appDiags...)
		return nil, diags
	}

	// Map rule
	rule, ruleDiags := r.mapRuleFromTerraform(ctx, model.Rule)
	if ruleDiags.HasError() {
		diags.Append(ruleDiags...)
		return nil, diags
	}

	// Map threshold
	threshold, thresholdDiags := r.mapThresholdFromTerraform(ctx, model.Threshold)
	if thresholdDiags.HasError() {
		diags.Append(thresholdDiags...)
		return nil, diags
	}

	// Map time threshold
	timeThreshold, timeThresholdDiags := r.mapTimeThresholdFromTerraform(ctx, model.TimeThreshold)
	if timeThresholdDiags.HasError() {
		diags.Append(timeThresholdDiags...)
		return nil, diags
	}

	// Map custom payload fields
	var customerPayloadFields []restapi.CustomPayloadField[any]
	if !model.CustomPayloadFields.IsNull() {
		var payloadDiags diag.Diagnostics
		customerPayloadFields, payloadDiags = MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
	}

	return &restapi.ApplicationAlertConfig{
		ID:                    id,
		Name:                  model.Name.ValueString(),
		Description:           model.Description.ValueString(),
		Severity:              severity,
		Triggering:            model.Triggering.ValueBool(),
		Applications:          applications,
		BoundaryScope:         restapi.BoundaryScope(model.BoundaryScope.ValueString()),
		TagFilterExpression:   tagFilter,
		IncludeInternal:       model.IncludeInternal.ValueBool(),
		IncludeSynthetic:      model.IncludeSynthetic.ValueBool(),
		EvaluationType:        restapi.ApplicationAlertEvaluationType(model.EvaluationType.ValueString()),
		AlertChannelIDs:       alertChannelIDs,
		Granularity:           restapi.Granularity(model.Granularity.ValueInt64()),
		CustomerPayloadFields: customerPayloadFields,
		Rule:                  rule,
		Threshold:             threshold,
		TimeThreshold:         timeThreshold,
	}, diags
}

// Helper methods for mapping applications

func (r *applicationAlertConfigResourceFramework) mapApplicationsFromTerraform(ctx context.Context, applications types.Set) (map[string]restapi.IncludedApplication, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make(map[string]restapi.IncludedApplication)

	if applications.IsNull() {
		return result, diags
	}

	var appList []ApplicationModel
	diags.Append(applications.ElementsAs(ctx, &appList, false)...)
	if diags.HasError() {
		return result, diags
	}

	for _, app := range appList {
		services, serviceDiags := r.mapServicesFromTerraform(ctx, app.Services)
		if serviceDiags.HasError() {
			diags.Append(serviceDiags...)
			return result, diags
		}

		result[app.ApplicationID.ValueString()] = restapi.IncludedApplication{
			ApplicationID: app.ApplicationID.ValueString(),
			Inclusive:     app.Inclusive.ValueBool(),
			Services:      services,
		}
	}

	return result, diags
}

func (r *applicationAlertConfigResourceFramework) mapServicesFromTerraform(ctx context.Context, services types.Set) (map[string]restapi.IncludedService, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make(map[string]restapi.IncludedService)

	if services.IsNull() {
		return result, diags
	}

	var serviceList []ServiceModel
	diags.Append(services.ElementsAs(ctx, &serviceList, false)...)
	if diags.HasError() {
		return result, diags
	}

	for _, service := range serviceList {
		endpoints, endpointDiags := r.mapEndpointsFromTerraform(ctx, service.Endpoints)
		if endpointDiags.HasError() {
			diags.Append(endpointDiags...)
			return result, diags
		}

		result[service.ServiceID.ValueString()] = restapi.IncludedService{
			ServiceID: service.ServiceID.ValueString(),
			Inclusive: service.Inclusive.ValueBool(),
			Endpoints: endpoints,
		}
	}

	return result, diags
}

func (r *applicationAlertConfigResourceFramework) mapEndpointsFromTerraform(ctx context.Context, endpoints types.Set) (map[string]restapi.IncludedEndpoint, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make(map[string]restapi.IncludedEndpoint)

	if endpoints.IsNull() {
		return result, diags
	}

	var endpointList []EndpointModel
	diags.Append(endpoints.ElementsAs(ctx, &endpointList, false)...)
	if diags.HasError() {
		return result, diags
	}

	for _, endpoint := range endpointList {
		result[endpoint.EndpointID.ValueString()] = restapi.IncludedEndpoint{
			EndpointID: endpoint.EndpointID.ValueString(),
			Inclusive:  endpoint.Inclusive.ValueBool(),
		}
	}

	return result, diags
}

func (r *applicationAlertConfigResourceFramework) mapApplicationsToTerraform(ctx context.Context, applications map[string]restapi.IncludedApplication) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	var appList []ApplicationModel

	for _, app := range applications {
		services, serviceDiags := r.mapServicesToTerraform(ctx, app.Services)
		if serviceDiags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{}}), serviceDiags
		}

		appList = append(appList, ApplicationModel{
			ApplicationID: types.StringValue(app.ApplicationID),
			Inclusive:     types.BoolValue(app.Inclusive),
			Services:      services,
		})
	}

	// Convert the list to a set
	appSet, diags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"application_id": types.StringType,
			"inclusive":      types.BoolType,
			"service":        types.SetType{ElemType: types.ObjectType{}},
		},
	}, appList)

	return appSet, diags
}

func (r *applicationAlertConfigResourceFramework) mapServicesToTerraform(ctx context.Context, services map[string]restapi.IncludedService) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	var serviceList []ServiceModel

	for _, service := range services {
		endpoints, endpointDiags := r.mapEndpointsToTerraform(ctx, service.Endpoints)
		if endpointDiags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{}}), endpointDiags
		}

		serviceList = append(serviceList, ServiceModel{
			ServiceID: types.StringValue(service.ServiceID),
			Inclusive: types.BoolValue(service.Inclusive),
			Endpoints: endpoints,
		})
	}

	// Convert the list to a set
	serviceSet, diags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"service_id": types.StringType,
			"inclusive":  types.BoolType,
			"endpoint":   types.SetType{ElemType: types.ObjectType{}},
		},
	}, serviceList)

	return serviceSet, diags
}

func (r *applicationAlertConfigResourceFramework) mapEndpointsToTerraform(ctx context.Context, endpoints map[string]restapi.IncludedEndpoint) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	var endpointList []EndpointModel

	for _, endpoint := range endpoints {
		endpointList = append(endpointList, EndpointModel{
			EndpointID: types.StringValue(endpoint.EndpointID),
			Inclusive:  types.BoolValue(endpoint.Inclusive),
		})
	}

	// Convert the list to a set
	endpointSet, diags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"endpoint_id": types.StringType,
			"inclusive":   types.BoolType,
		},
	}, endpointList)

	return endpointSet, diags
}

// Helper methods for mapping rules

func (r *applicationAlertConfigResourceFramework) mapRuleFromTerraform(ctx context.Context, rule types.List) (restapi.ApplicationAlertRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result restapi.ApplicationAlertRule

	if rule.IsNull() {
		return result, diags
	}

	var ruleList []RuleModel
	diags.Append(rule.ElementsAs(ctx, &ruleList, false)...)
	if diags.HasError() {
		return result, diags
	}

	if len(ruleList) == 0 {
		return result, diags
	}

	ruleModel := ruleList[0]

	// Check which rule type is set and map accordingly
	if !ruleModel.ErrorRate.IsNull() {
		var errorRateList []RuleConfigModel
		diags.Append(ruleModel.ErrorRate.ElementsAs(ctx, &errorRateList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(errorRateList) > 0 {
			errorRate := errorRateList[0]
			result.AlertType = "errorRate"
			result.MetricName = errorRate.MetricName.ValueString()
			if !errorRate.Aggregation.IsNull() {
				result.Aggregation = restapi.Aggregation(errorRate.Aggregation.ValueString())
			}
		}
	} else if !ruleModel.Errors.IsNull() {
		var errorsList []RuleConfigModel
		diags.Append(ruleModel.Errors.ElementsAs(ctx, &errorsList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(errorsList) > 0 {
			errors := errorsList[0]
			result.AlertType = "errors"
			result.MetricName = errors.MetricName.ValueString()
			if !errors.Aggregation.IsNull() {
				result.Aggregation = restapi.Aggregation(errors.Aggregation.ValueString())
			}
		}
	} else if !ruleModel.Logs.IsNull() {
		var logsList []LogsRuleModel
		diags.Append(ruleModel.Logs.ElementsAs(ctx, &logsList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(logsList) > 0 {
			logs := logsList[0]
			result.AlertType = "logs"
			result.MetricName = logs.MetricName.ValueString()
			if !logs.Aggregation.IsNull() {
				result.Aggregation = restapi.Aggregation(logs.Aggregation.ValueString())
			}
			if !logs.Level.IsNull() {
				level := restapi.LogLevel(logs.Level.ValueString())
				result.Level = &level
			}
			if !logs.Message.IsNull() {
				message := logs.Message.ValueString()
				result.Message = &message
			}
			if !logs.Operator.IsNull() {
				operator := restapi.ExpressionOperator(logs.Operator.ValueString())
				result.Operator = &operator
			}
		}
	} else if !ruleModel.Slowness.IsNull() {
		var slownessList []RuleConfigModel
		diags.Append(ruleModel.Slowness.ElementsAs(ctx, &slownessList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(slownessList) > 0 {
			slowness := slownessList[0]
			result.AlertType = "slowness"
			result.MetricName = slowness.MetricName.ValueString()
			result.Aggregation = restapi.Aggregation(slowness.Aggregation.ValueString())
		}
	} else if !ruleModel.StatusCode.IsNull() {
		var statusCodeList []StatusCodeRuleModel
		diags.Append(ruleModel.StatusCode.ElementsAs(ctx, &statusCodeList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(statusCodeList) > 0 {
			statusCode := statusCodeList[0]
			result.AlertType = "statusCode"
			result.MetricName = statusCode.MetricName.ValueString()
			if !statusCode.Aggregation.IsNull() {
				result.Aggregation = restapi.Aggregation(statusCode.Aggregation.ValueString())
			}
			if !statusCode.StatusCodeStart.IsNull() {
				start := int32(statusCode.StatusCodeStart.ValueInt64())
				result.StatusCodeStart = &start
			}
			if !statusCode.StatusCodeEnd.IsNull() {
				end := int32(statusCode.StatusCodeEnd.ValueInt64())
				result.StatusCodeEnd = &end
			}
		}
	} else if !ruleModel.Throughput.IsNull() {
		var throughputList []RuleConfigModel
		diags.Append(ruleModel.Throughput.ElementsAs(ctx, &throughputList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(throughputList) > 0 {
			throughput := throughputList[0]
			result.AlertType = "throughput"
			result.MetricName = throughput.MetricName.ValueString()
			if !throughput.Aggregation.IsNull() {
				result.Aggregation = restapi.Aggregation(throughput.Aggregation.ValueString())
			}
		}
	}

	return result, diags
}

func (r *applicationAlertConfigResourceFramework) mapRuleToTerraform(ctx context.Context, rule *restapi.ApplicationAlertRule) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	ruleModel := RuleModel{}

	// Create the appropriate rule type based on the alert type
	switch rule.AlertType {
	case "errorRate":
		ruleConfig := RuleConfigModel{
			MetricName: types.StringValue(rule.MetricName),
		}
		if rule.Aggregation != "" {
			ruleConfig.Aggregation = types.StringValue(string(rule.Aggregation))
		} else {
			ruleConfig.Aggregation = types.StringNull()
		}

		errorRateList, errorRateDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, []RuleConfigModel{ruleConfig})
		if errorRateDiags.HasError() {
			return types.ListNull(types.ObjectType{}), errorRateDiags
		}
		ruleModel.ErrorRate = errorRateList
	case "errors":
		ruleConfig := RuleConfigModel{
			MetricName: types.StringValue(rule.MetricName),
		}
		if rule.Aggregation != "" {
			ruleConfig.Aggregation = types.StringValue(string(rule.Aggregation))
		} else {
			ruleConfig.Aggregation = types.StringNull()
		}

		errorsList, errorsDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, []RuleConfigModel{ruleConfig})
		if errorsDiags.HasError() {
			return types.ListNull(types.ObjectType{}), errorsDiags
		}
		ruleModel.Errors = errorsList
	case "logs":
		logsRule := LogsRuleModel{
			MetricName: types.StringValue(rule.MetricName),
		}
		if rule.Aggregation != "" {
			logsRule.Aggregation = types.StringValue(string(rule.Aggregation))
		} else {
			logsRule.Aggregation = types.StringNull()
		}
		if rule.Level != nil {
			logsRule.Level = types.StringValue(string(*rule.Level))
		} else {
			logsRule.Level = types.StringNull()
		}
		if rule.Message != nil {
			logsRule.Message = types.StringValue(*rule.Message)
		} else {
			logsRule.Message = types.StringNull()
		}
		if rule.Operator != nil {
			logsRule.Operator = types.StringValue(string(*rule.Operator))
		} else {
			logsRule.Operator = types.StringNull()
		}

		logsList, logsDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
				"level":       types.StringType,
				"message":     types.StringType,
				"operator":    types.StringType,
			},
		}, []LogsRuleModel{logsRule})
		if logsDiags.HasError() {
			return types.ListNull(types.ObjectType{}), logsDiags
		}
		ruleModel.Logs = logsList
	case "slowness":
		ruleConfig := RuleConfigModel{
			MetricName:  types.StringValue(rule.MetricName),
			Aggregation: types.StringValue(string(rule.Aggregation)),
		}

		slownessList, slownessDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, []RuleConfigModel{ruleConfig})
		if slownessDiags.HasError() {
			return types.ListNull(types.ObjectType{}), slownessDiags
		}
		ruleModel.Slowness = slownessList
	case "statusCode":
		statusCodeRule := StatusCodeRuleModel{
			MetricName: types.StringValue(rule.MetricName),
		}
		if rule.Aggregation != "" {
			statusCodeRule.Aggregation = types.StringValue(string(rule.Aggregation))
		} else {
			statusCodeRule.Aggregation = types.StringNull()
		}
		if rule.StatusCodeStart != nil {
			statusCodeRule.StatusCodeStart = types.Int64Value(int64(*rule.StatusCodeStart))
		} else {
			statusCodeRule.StatusCodeStart = types.Int64Null()
		}
		if rule.StatusCodeEnd != nil {
			statusCodeRule.StatusCodeEnd = types.Int64Value(int64(*rule.StatusCodeEnd))
		} else {
			statusCodeRule.StatusCodeEnd = types.Int64Null()
		}

		statusCodeList, statusCodeDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name":       types.StringType,
				"aggregation":       types.StringType,
				"status_code_start": types.Int64Type,
				"status_code_end":   types.Int64Type,
			},
		}, []StatusCodeRuleModel{statusCodeRule})
		if statusCodeDiags.HasError() {
			return types.ListNull(types.ObjectType{}), statusCodeDiags
		}
		ruleModel.StatusCode = statusCodeList
	case "throughput":
		ruleConfig := RuleConfigModel{
			MetricName: types.StringValue(rule.MetricName),
		}
		if rule.Aggregation != "" {
			ruleConfig.Aggregation = types.StringValue(string(rule.Aggregation))
		} else {
			ruleConfig.Aggregation = types.StringNull()
		}

		throughputList, throughputDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, []RuleConfigModel{ruleConfig})
		if throughputDiags.HasError() {
			return types.ListNull(types.ObjectType{}), throughputDiags
		}
		ruleModel.Throughput = throughputList
	}

	// Convert the rule model to a list
	ruleList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"error_rate":  types.ListType{ElemType: types.ObjectType{}},
			"errors":      types.ListType{ElemType: types.ObjectType{}},
			"logs":        types.ListType{ElemType: types.ObjectType{}},
			"slowness":    types.ListType{ElemType: types.ObjectType{}},
			"status_code": types.ListType{ElemType: types.ObjectType{}},
			"throughput":  types.ListType{ElemType: types.ObjectType{}},
		},
	}, []RuleModel{ruleModel})

	return ruleList, diags
}

// Helper methods for mapping threshold

func (r *applicationAlertConfigResourceFramework) mapThresholdFromTerraform(ctx context.Context, threshold types.Object) (restapi.Threshold, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result restapi.Threshold

	if threshold.IsNull() {
		return result, diags
	}

	// Get the threshold as a map
	var thresholdMap map[string]interface{}
	diags.Append(threshold.As(ctx, &thresholdMap, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return result, diags
	}

	// Extract type
	if typeVal, ok := thresholdMap["type"]; ok {
		typeStr := typeVal.(string)
		result.Type = typeStr
	}

	// Extract value
	if valueVal, ok := thresholdMap["value"]; ok {
		value := valueVal.(float64)
		result.Value = &value
	}

	return result, diags
}

func (r *applicationAlertConfigResourceFramework) mapThresholdToTerraform(ctx context.Context, threshold *restapi.Threshold) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create a map for the threshold values
	thresholdMap := map[string]attr.Value{
		"type": types.StringValue(threshold.Type),
	}

	// Add value if it exists
	if threshold.Value != nil {
		thresholdMap["value"] = types.Float64Value(*threshold.Value)
	} else {
		thresholdMap["value"] = types.Float64Null()
	}

	// Create the object
	thresholdObj, diags := types.ObjectValue(
		map[string]attr.Type{
			"type":  types.StringType,
			"value": types.Float64Type,
		},
		thresholdMap,
	)

	return thresholdObj, diags
}

// Helper methods for mapping time threshold

func (r *applicationAlertConfigResourceFramework) mapTimeThresholdFromTerraform(ctx context.Context, timeThreshold types.List) (restapi.TimeThreshold, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result restapi.TimeThreshold

	if timeThreshold.IsNull() {
		return result, diags
	}

	var timeThresholdList []AppAlertTimeThresholdModel
	diags.Append(timeThreshold.ElementsAs(ctx, &timeThresholdList, false)...)
	if diags.HasError() {
		return result, diags
	}

	if len(timeThresholdList) == 0 {
		return result, diags
	}

	timeThresholdModel := timeThresholdList[0]

	// Check which time threshold type is set and map accordingly
	if !timeThresholdModel.RequestImpact.IsNull() {
		var requestImpactList []AppAlertRequestImpactModel
		diags.Append(timeThresholdModel.RequestImpact.ElementsAs(ctx, &requestImpactList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(requestImpactList) > 0 {
			requestImpact := requestImpactList[0]
			result.Type = "requestImpact"
			result.TimeWindow = requestImpact.TimeWindow.ValueInt64()
			if !requestImpact.Requests.IsNull() {
				requests := int32(requestImpact.Requests.ValueInt64())
				result.Requests = &requests
			}
		}
	} else if !timeThresholdModel.ViolationsInPeriod.IsNull() {
		var violationsInPeriodList []AppAlertViolationsInPeriodModel
		diags.Append(timeThresholdModel.ViolationsInPeriod.ElementsAs(ctx, &violationsInPeriodList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(violationsInPeriodList) > 0 {
			violationsInPeriod := violationsInPeriodList[0]
			result.Type = "violationsInPeriod"
			result.TimeWindow = violationsInPeriod.TimeWindow.ValueInt64()
			if !violationsInPeriod.Violations.IsNull() {
				violations := int32(violationsInPeriod.Violations.ValueInt64())
				result.Violations = &violations
			}
		}
	} else if !timeThresholdModel.ViolationsInSequence.IsNull() {
		var violationsInSequenceList []AppAlertViolationsInSequenceModel
		diags.Append(timeThresholdModel.ViolationsInSequence.ElementsAs(ctx, &violationsInSequenceList, false)...)
		if diags.HasError() {
			return result, diags
		}

		if len(violationsInSequenceList) > 0 {
			violationsInSequence := violationsInSequenceList[0]
			result.Type = "violationsInSequence"
			result.TimeWindow = violationsInSequence.TimeWindow.ValueInt64()
		}
	}

	return result, diags
}

func (r *applicationAlertConfigResourceFramework) mapTimeThresholdToTerraform(ctx context.Context, timeThreshold *restapi.TimeThreshold) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	timeThresholdModel := AppAlertTimeThresholdModel{}

	// Create the appropriate time threshold type based on the type
	switch timeThreshold.Type {
	case "requestImpact":
		requestImpactModel := AppAlertRequestImpactModel{
			TimeWindow: types.Int64Value(timeThreshold.TimeWindow),
		}
		if timeThreshold.Requests != nil {
			requestImpactModel.Requests = types.Int64Value(int64(*timeThreshold.Requests))
		} else {
			requestImpactModel.Requests = types.Int64Null()
		}

		requestImpactList, requestImpactDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"time_window": types.Int64Type,
				"requests":    types.Int64Type,
			},
		}, []AppAlertRequestImpactModel{requestImpactModel})
		if requestImpactDiags.HasError() {
			return types.ListNull(types.ObjectType{}), requestImpactDiags
		}
		timeThresholdModel.RequestImpact = requestImpactList
	case "violationsInPeriod":
		violationsInPeriodModel := AppAlertViolationsInPeriodModel{
			TimeWindow: types.Int64Value(timeThreshold.TimeWindow),
		}
		if timeThreshold.Violations != nil {
			violationsInPeriodModel.Violations = types.Int64Value(int64(*timeThreshold.Violations))
		} else {
			violationsInPeriodModel.Violations = types.Int64Null()
		}

		violationsInPeriodList, violationsInPeriodDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"time_window": types.Int64Type,
				"violations":  types.Int64Type,
			},
		}, []AppAlertViolationsInPeriodModel{violationsInPeriodModel})
		if violationsInPeriodDiags.HasError() {
			return types.ListNull(types.ObjectType{}), violationsInPeriodDiags
		}
		timeThresholdModel.ViolationsInPeriod = violationsInPeriodList
	case "violationsInSequence":
		violationsInSequenceModel := AppAlertViolationsInSequenceModel{
			TimeWindow: types.Int64Value(timeThreshold.TimeWindow),
		}

		violationsInSequenceList, violationsInSequenceDiags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"time_window": types.Int64Type,
			},
		}, []AppAlertViolationsInSequenceModel{violationsInSequenceModel})
		if violationsInSequenceDiags.HasError() {
			return types.ListNull(types.ObjectType{}), violationsInSequenceDiags
		}
		timeThresholdModel.ViolationsInSequence = violationsInSequenceList
	}

	// Convert the time threshold model to a list
	timeThresholdList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"request_impact":         types.ListType{ElemType: types.ObjectType{}},
			"violations_in_period":   types.ListType{ElemType: types.ObjectType{}},
			"violations_in_sequence": types.ListType{ElemType: types.ObjectType{}},
		},
	}, []AppAlertTimeThresholdModel{timeThresholdModel})

	return timeThresholdList, diags
}

// Made with Bob
