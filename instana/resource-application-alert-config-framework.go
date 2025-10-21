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

const (
	//ApplicationAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_application_alert_config
	ApplicationAlertConfigFieldAlertChannelIDs = "alert_channel_ids"
	//ApplicationAlertConfigFieldBoundaryScope constant value for field boundary_scope of resource instana_application_alert_config
	ApplicationAlertConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationAlertConfigFieldDescription constant value for field description of resource instana_application_alert_config
	ApplicationAlertConfigFieldDescription = "description"
	//ApplicationAlertConfigFieldEvaluationType constant value for field evaluation_type of resource instana_application_alert_config
	ApplicationAlertConfigFieldEvaluationType = "evaluation_type"
	//ApplicationAlertConfigFieldGranularity constant value for field granularity of resource instana_application_alert_config
	ApplicationAlertConfigFieldGranularity = "granularity"
	//ApplicationAlertConfigFieldIncludeInternal constant value for field include_internal of resource instana_application_alert_config
	ApplicationAlertConfigFieldIncludeInternal = "include_internal"
	//ApplicationAlertConfigFieldIncludeSynthetic constant value for field include_synthetic of resource instana_application_alert_config
	ApplicationAlertConfigFieldIncludeSynthetic = "include_synthetic"
	//ApplicationAlertConfigFieldName constant value for field name of resource instana_application_alert_config
	ApplicationAlertConfigFieldName = "name"
	//ApplicationAlertConfigFieldFullName constant value for field full_name of resource instana_application_alert_config
	ApplicationAlertConfigFieldFullName = "full_name"
	//ApplicationAlertConfigFieldRule constant value for field rule of resource instana_application_alert_config
	ApplicationAlertConfigFieldRule = "rule"
	//ApplicationAlertConfigFieldRuleMetricName constant value for field rule.*.metric_name of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleMetricName = "metric_name"
	//ApplicationAlertConfigFieldRuleAggregation constant value for field rule.*.aggregation of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleAggregation = "aggregation"
	//ApplicationAlertConfigFieldRuleErrorRate constant value for field rule.error_rate of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleErrorRate = "error_rate"
	//ApplicationAlertConfigFieldRuleErrors constant value for field rule.errors of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleErrors = "errors"
	//ApplicationAlertConfigFieldRuleLogs constant value for field rule.logs of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogs = "logs"
	//ApplicationAlertConfigFieldRuleLogsLevel constant value for field rule.logs.level of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsLevel = "level"
	//ApplicationAlertConfigFieldRuleLogsMessage constant value for field rule.logs.message of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsMessage = "message"
	//ApplicationAlertConfigFieldRuleLogsOperator constant value for field rule.logs.operator of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsOperator = "operator"
	//ApplicationAlertConfigFieldRuleSlowness constant value for field rule.slowness of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleSlowness = "slowness"
	//ApplicationAlertConfigFieldRuleStatusCode constant value for field rule.status_code of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCode = "status_code"
	//ApplicationAlertConfigFieldRuleStatusCodeStart constant value for field rule.status_code.status_code_start of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCodeStart = "status_code_start"
	//ApplicationAlertConfigFieldRuleStatusCodeEnd constant value for field rule.status_code.status_code_end of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCodeEnd = "status_code_end"
	//ApplicationAlertConfigFieldRuleThroughput constant value for field rule.throughput of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleThroughput = "throughput"
	//ApplicationAlertConfigFieldSeverity constant value for field severity of resource instana_application_alert_config
	ApplicationAlertConfigFieldSeverity = "severity"
	//ApplicationAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_application_alert_config
	ApplicationAlertConfigFieldTagFilter = "tag_filter"
	//ApplicationAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThreshold = "time_threshold"
	//ApplicationAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	//ApplicationAlertConfigFieldTimeThresholdRequestImpact constant value for field time_threshold.request_impact of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdRequestImpact = "request_impact"
	//ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests constant value for field time_threshold.request_impact.requests of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests = "requests"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	//ApplicationAlertConfigFieldTriggering constant value for field triggering of resource instana_application_alert_config
	ApplicationAlertConfigFieldTriggering = "triggering"
)

// ResourceFramework defines the interface for resource implementations
type ResourceFramework[T restapi.InstanaDataObject] interface {
	// GetID returns the ID of the data object
	GetID(data T) string

	// SetID sets the ID of the data object
	SetID(data T, id string)

	// MapStateToDataObject maps the current state to the API model
	MapStateToDataObject(ctx context.Context, state tfsdk.State) (T, diag.Diagnostics)

	// UpdateState updates the state with the given data object
	UpdateState(ctx context.Context, data T, state *tfsdk.State) diag.Diagnostics
}

// Additional application alert config field names
const (
	ApplicationAlertConfigFieldRules             = "rules"
	ApplicationAlertConfigFieldThresholds        = "thresholds"
	ApplicationAlertConfigFieldThresholdOperator = "threshold_operator"
	ApplicationAlertConfigFieldGracePeriod       = "grace_period"
	ApplicationAlertConfigFieldAlertChannels     = "alert_channels"
	ApplicationAlertConfigFieldRuleConfig        = "rule_config"
	ApplicationAlertConfigFieldValue             = "value"

	// Re-define constants from resource-application-alert-config.go for compatibility
	ApplicationAlertConfigFieldApplications                            = "application"
	ApplicationAlertConfigFieldApplicationsApplicationID               = "application_id"
	ApplicationAlertConfigFieldApplicationsInclusive                   = "inclusive"
	ApplicationAlertConfigFieldApplicationsServices                    = "service"
	ApplicationAlertConfigFieldApplicationsServicesServiceID           = "service_id"
	ApplicationAlertConfigFieldApplicationsServicesEndpoints           = "endpoint"
	ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID = "endpoint_id"
)

// ApplicationAlertConfigModel represents the data model for the application alert configuration resource
type ApplicationAlertConfigModel struct {
	ID                  types.String `tfsdk:"id"`
	AlertChannelIDs     types.Set    `tfsdk:"alert_channel_ids"`
	AlertChannels       types.Map    `tfsdk:"alert_channels"`
	Applications        types.Set    `tfsdk:"application"`
	BoundaryScope       types.String `tfsdk:"boundary_scope"`
	CustomPayloadFields types.List   `tfsdk:"custom_payload_field"`
	Description         types.String `tfsdk:"description"`
	EvaluationType      types.String `tfsdk:"evaluation_type"`
	GracePeriod         types.Int64  `tfsdk:"grace_period"`
	Granularity         types.Int64  `tfsdk:"granularity"`
	IncludeInternal     types.Bool   `tfsdk:"include_internal"`
	IncludeSynthetic    types.Bool   `tfsdk:"include_synthetic"`
	Name                types.String `tfsdk:"name"`
	Rule                types.List   `tfsdk:"rule"`
	Rules               types.List   `tfsdk:"rules"`
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

// RuleWithThresholdModel represents a rule with multiple thresholds and severity levels
type RuleWithThresholdModel struct {
	Rule              types.Object `tfsdk:"rule"`
	ThresholdOperator types.String `tfsdk:"threshold_operator"`
	Thresholds        types.Map    `tfsdk:"thresholds"`
}

// ThresholdConfigRuleModel represents a threshold configuration for a rule
type ThresholdConfigRuleModel struct {
	Value types.Float64 `tfsdk:"value"`
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
						Description: "List of IDs of alert channels defined in Instana. Deprecated: Use alert_channels instead.",
						ElementType: types.StringType,
					},
					ApplicationAlertConfigFieldAlertChannels: schema.MapAttribute{
						Optional:    true,
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
					ApplicationAlertConfigFieldRules: schema.ListNestedBlock{
						Description: "A list of rules where each rule is associated with multiple thresholds and their corresponding severity levels.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								ApplicationAlertConfigFieldThresholdOperator: schema.StringAttribute{
									Required:    true,
									Description: "The operator to apply for threshold comparison",
									Validators: []validator.String{
										stringvalidator.OneOf(">", ">=", "<", "<="),
									},
								},
							},
							Blocks: map[string]schema.Block{
								ApplicationAlertConfigFieldRule: schema.SingleNestedBlock{
									Description: "The rule configuration",
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
								ApplicationAlertConfigFieldThresholds: schema.SingleNestedBlock{
									Description: "Map of severity to threshold configurations",
									Attributes: map[string]schema.Attribute{
										"warning": schema.Float64Attribute{
											Optional:    true,
											Description: "The threshold value for warning severity level",
										},
										"critical": schema.Float64Attribute{
											Optional:    true,
											Description: "The threshold value for critical severity level",
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

func (r *applicationAlertConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.ApplicationAlertConfig, diag.Diagnostics) {
	// Delegate to the resource implementation
	resource, diags := r.NewResource(ctx, nil)
	if diags.HasError() {
		return nil, diags
	}

	// If state is provided, use it, otherwise use plan
	if state != nil {
		return resource.MapStateToDataObject(ctx, *state)
	} else if plan != nil {
		// Create a temporary state from the plan
		var tempState tfsdk.State
		tempState.Schema = plan.Schema
		diags = plan.Get(ctx, &tempState.Raw)
		if diags.HasError() {
			return nil, diags
		}
		return resource.MapStateToDataObject(ctx, tempState)
	} else {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Invalid parameters",
				"Either plan or state must be provided",
			),
		}
	}
}

func (r *applicationAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, obj *restapi.ApplicationAlertConfig) diag.Diagnostics {
	// Delegate to the resource implementation
	resource, diags := r.NewResource(ctx, nil)
	if diags.HasError() {
		return diags
	}

	return resource.UpdateState(ctx, obj, state)
}

func (r *applicationAlertConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	// No computed fields to set
	return nil
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
	if !model.GracePeriod.IsNull() && !model.GracePeriod.IsUnknown() {
		result.GracePeriod = int(model.GracePeriod.ValueInt64())
	}

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
		var customPayloadFields []map[string]interface{}
		diags = model.CustomPayloadFields.ElementsAs(ctx, &customPayloadFields, false)
		if diags.HasError() {
			return nil, diags
		}

		fields := make([]restapi.CustomPayloadField[interface{}], len(customPayloadFields))
		for i, field := range customPayloadFields {
			key, ok := field["key"].(string)
			if !ok {
				return nil, diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Invalid custom payload field",
						"Custom payload field key must be a string",
					),
				}
			}
			value, ok := field["value"].(string)
			if !ok {
				return nil, diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Invalid custom payload field",
						"Custom payload field value must be a string",
					),
				}
			}
			payloadType, ok := field["type"].(string)
			if !ok {
				payloadType = string(restapi.StaticStringCustomPayloadType)
			}
			fields[i] = restapi.CustomPayloadField[interface{}]{
				Type:  restapi.CustomPayloadType(payloadType),
				Key:   key,
				Value: value,
			}
		}
		result.CustomerPayloadFields = fields
	}

	// Handle threshold
	if !model.Threshold.IsNull() && !model.Threshold.IsUnknown() {
		var threshold map[string]interface{}
		diags = model.Threshold.As(ctx, &threshold, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, diags
		}

		thresholdType, ok := threshold["type"].(string)
		if !ok {
			return nil, diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Invalid threshold",
					"Threshold type must be a string",
				),
			}
		}
		thresholdValue, ok := threshold["value"].(float64)
		if !ok {
			return nil, diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Invalid threshold",
					"Threshold value must be a number",
				),
			}
		}

		thresholdValuePtr := thresholdValue
		result.Threshold = &restapi.Threshold{
			Type:  thresholdType,
			Value: &thresholdValuePtr,
		}
	}

	// Handle rule (deprecated but supported for backward compatibility)
	if !model.Rule.IsNull() && !model.Rule.IsUnknown() {
		var rules []RuleModel
		diags = model.Rule.ElementsAs(ctx, &rules, false)
		if diags.HasError() {
			return nil, diags
		}

		if len(rules) > 0 {
			rule := rules[0]
			result.Rule = &restapi.ApplicationAlertRule{}

			// Handle error rate rule
			if !rule.ErrorRate.IsNull() && !rule.ErrorRate.IsUnknown() {
				var errorRateRules []RuleConfigModel
				diags = rule.ErrorRate.ElementsAs(ctx, &errorRateRules, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(errorRateRules) > 0 {
					result.Rule.ErrorRate = &restapi.ApplicationAlertRuleErrorRate{
						MetricName:  errorRateRules[0].MetricName.ValueString(),
						Aggregation: errorRateRules[0].Aggregation.ValueString(),
					}
				}
			}

			// Handle errors rule
			if !rule.Errors.IsNull() && !rule.Errors.IsUnknown() {
				var errorsRules []RuleConfigModel
				diags = rule.Errors.ElementsAs(ctx, &errorsRules, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(errorsRules) > 0 {
					result.Rule.Errors = &restapi.ApplicationAlertRuleErrors{
						MetricName:  errorsRules[0].MetricName.ValueString(),
						Aggregation: errorsRules[0].Aggregation.ValueString(),
					}
				}
			}

			// Handle logs rule
			if !rule.Logs.IsNull() && !rule.Logs.IsUnknown() {
				var logsRules []LogsRuleModel
				diags = rule.Logs.ElementsAs(ctx, &logsRules, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(logsRules) > 0 {
					result.Rule.Logs = &restapi.ApplicationAlertRuleLogs{
						MetricName:  logsRules[0].MetricName.ValueString(),
						Aggregation: logsRules[0].Aggregation.ValueString(),
						Level:       logsRules[0].Level.ValueString(),
						Message:     logsRules[0].Message.ValueString(),
						Operator:    logsRules[0].Operator.ValueString(),
					}
				}
			}

			// Handle slowness rule
			if !rule.Slowness.IsNull() && !rule.Slowness.IsUnknown() {
				var slownessRules []RuleConfigModel
				diags = rule.Slowness.ElementsAs(ctx, &slownessRules, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(slownessRules) > 0 {
					result.Rule.Slowness = &restapi.ApplicationAlertRuleSlowness{
						MetricName:  slownessRules[0].MetricName.ValueString(),
						Aggregation: slownessRules[0].Aggregation.ValueString(),
					}
				}
			}

			// Handle status code rule
			if !rule.StatusCode.IsNull() && !rule.StatusCode.IsUnknown() {
				var statusCodeRules []StatusCodeRuleModel
				diags = rule.StatusCode.ElementsAs(ctx, &statusCodeRules, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(statusCodeRules) > 0 {
					result.Rule.StatusCode = &restapi.ApplicationAlertRuleStatusCode{
						MetricName:      statusCodeRules[0].MetricName.ValueString(),
						Aggregation:     statusCodeRules[0].Aggregation.ValueString(),
						StatusCodeStart: int(statusCodeRules[0].StatusCodeStart.ValueInt64()),
						StatusCodeEnd:   int(statusCodeRules[0].StatusCodeEnd.ValueInt64()),
					}
				}
			}

			// Handle throughput rule
			if !rule.Throughput.IsNull() && !rule.Throughput.IsUnknown() {
				var throughputRules []RuleConfigModel
				diags = rule.Throughput.ElementsAs(ctx, &throughputRules, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(throughputRules) > 0 {
					result.Rule.Throughput = &restapi.ApplicationAlertRuleThroughput{
						MetricName:  throughputRules[0].MetricName.ValueString(),
						Aggregation: throughputRules[0].Aggregation.ValueString(),
					}
				}
			}
		}
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
					var errorRateRules []RuleConfigModel
					diags = rule.ErrorRate.ElementsAs(ctx, &errorRateRules, false)
					if diags.HasError() {
						return nil, diags
					}
					if len(errorRateRules) > 0 {
						result.Rules[i].Rule.ErrorRate = &restapi.ApplicationAlertRuleErrorRate{
							MetricName:  errorRateRules[0].MetricName.ValueString(),
							Aggregation: errorRateRules[0].Aggregation.ValueString(),
						}
					}
				}

				// Handle errors rule
				if !rule.Errors.IsNull() && !rule.Errors.IsUnknown() {
					var errorsRules []RuleConfigModel
					diags = rule.Errors.ElementsAs(ctx, &errorsRules, false)
					if diags.HasError() {
						return nil, diags
					}
					if len(errorsRules) > 0 {
						result.Rules[i].Rule.Errors = &restapi.ApplicationAlertRuleErrors{
							MetricName:  errorsRules[0].MetricName.ValueString(),
							Aggregation: errorsRules[0].Aggregation.ValueString(),
						}
					}
				}

				// Handle logs rule
				if !rule.Logs.IsNull() && !rule.Logs.IsUnknown() {
					var logsRules []LogsRuleModel
					diags = rule.Logs.ElementsAs(ctx, &logsRules, false)
					if diags.HasError() {
						return nil, diags
					}
					if len(logsRules) > 0 {
						result.Rules[i].Rule.Logs = &restapi.ApplicationAlertRuleLogs{
							MetricName:  logsRules[0].MetricName.ValueString(),
							Aggregation: logsRules[0].Aggregation.ValueString(),
							Level:       logsRules[0].Level.ValueString(),
							Message:     logsRules[0].Message.ValueString(),
							Operator:    logsRules[0].Operator.ValueString(),
						}
					}
				}

				// Handle slowness rule
				if !rule.Slowness.IsNull() && !rule.Slowness.IsUnknown() {
					var slownessRules []RuleConfigModel
					diags = rule.Slowness.ElementsAs(ctx, &slownessRules, false)
					if diags.HasError() {
						return nil, diags
					}
					if len(slownessRules) > 0 {
						result.Rules[i].Rule.Slowness = &restapi.ApplicationAlertRuleSlowness{
							MetricName:  slownessRules[0].MetricName.ValueString(),
							Aggregation: slownessRules[0].Aggregation.ValueString(),
						}
					}
				}

				// Handle status code rule
				if !rule.StatusCode.IsNull() && !rule.StatusCode.IsUnknown() {
					var statusCodeRules []StatusCodeRuleModel
					diags = rule.StatusCode.ElementsAs(ctx, &statusCodeRules, false)
					if diags.HasError() {
						return nil, diags
					}
					if len(statusCodeRules) > 0 {
						result.Rules[i].Rule.StatusCode = &restapi.ApplicationAlertRuleStatusCode{
							MetricName:      statusCodeRules[0].MetricName.ValueString(),
							Aggregation:     statusCodeRules[0].Aggregation.ValueString(),
							StatusCodeStart: int(statusCodeRules[0].StatusCodeStart.ValueInt64()),
							StatusCodeEnd:   int(statusCodeRules[0].StatusCodeEnd.ValueInt64()),
						}
					}
				}

				// Handle throughput rule
				if !rule.Throughput.IsNull() && !rule.Throughput.IsUnknown() {
					var throughputRules []RuleConfigModel
					diags = rule.Throughput.ElementsAs(ctx, &throughputRules, false)
					if diags.HasError() {
						return nil, diags
					}
					if len(throughputRules) > 0 {
						result.Rules[i].Rule.Throughput = &restapi.ApplicationAlertRuleThroughput{
							MetricName:  throughputRules[0].MetricName.ValueString(),
							Aggregation: throughputRules[0].Aggregation.ValueString(),
						}
					}
				}
			}

			// Handle thresholds
			if !ruleWithThreshold.Thresholds.IsNull() && !ruleWithThreshold.Thresholds.IsUnknown() {
				thresholds := make(map[string]restapi.ThresholdValue)
				for k, v := range ruleWithThreshold.Thresholds.Elements() {
					var thresholdConfig ThresholdConfigRuleModel
					diags = v.(types.Object).As(ctx, &thresholdConfig, basetypes.ObjectAsOptions{})
					if diags.HasError() {
						return nil, diags
					}
					thresholds[k] = restapi.ThresholdValue{
						Value: thresholdConfig.Value.ValueFloat64(),
					}
				}
				result.Rules[i].Thresholds = thresholds
			}
		}
	}

	// Handle time threshold
	if !model.TimeThreshold.IsNull() && !model.TimeThreshold.IsUnknown() {
		var timeThresholds []AppAlertTimeThresholdModel
		diags = model.TimeThreshold.ElementsAs(ctx, &timeThresholds, false)
		if diags.HasError() {
			return nil, diags
		}

		if len(timeThresholds) > 0 {
			timeThreshold := timeThresholds[0]
			result.TimeThreshold = &restapi.ApplicationAlertTimeThreshold{}

			// Handle request impact
			if !timeThreshold.RequestImpact.IsNull() && !timeThreshold.RequestImpact.IsUnknown() {
				var requestImpacts []AppAlertRequestImpactModel
				diags = timeThreshold.RequestImpact.ElementsAs(ctx, &requestImpacts, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(requestImpacts) > 0 {
					result.TimeThreshold.RequestImpact = &restapi.ApplicationAlertTimeThresholdRequestImpact{
						TimeWindow: int(requestImpacts[0].TimeWindow.ValueInt64()),
						Requests:   int(requestImpacts[0].Requests.ValueInt64()),
					}
				}
			}

			// Handle violations in period
			if !timeThreshold.ViolationsInPeriod.IsNull() && !timeThreshold.ViolationsInPeriod.IsUnknown() {
				var violationsInPeriods []AppAlertViolationsInPeriodModel
				diags = timeThreshold.ViolationsInPeriod.ElementsAs(ctx, &violationsInPeriods, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(violationsInPeriods) > 0 {
					result.TimeThreshold.ViolationsInPeriod = &restapi.ApplicationAlertTimeThresholdViolationsInPeriod{
						TimeWindow: int(violationsInPeriods[0].TimeWindow.ValueInt64()),
						Violations: int(violationsInPeriods[0].Violations.ValueInt64()),
					}
				}
			}

			// Handle violations in sequence
			if !timeThreshold.ViolationsInSequence.IsNull() && !timeThreshold.ViolationsInSequence.IsUnknown() {
				var violationsInSequences []AppAlertViolationsInSequenceModel
				diags = timeThreshold.ViolationsInSequence.ElementsAs(ctx, &violationsInSequences, false)
				if diags.HasError() {
					return nil, diags
				}
				if len(violationsInSequences) > 0 {
					result.TimeThreshold.ViolationsInSequence = &restapi.ApplicationAlertTimeThresholdViolationsInSequence{
						TimeWindow: int(violationsInSequences[0].TimeWindow.ValueInt64()),
					}
				}
			}
		}
	}

	return result, nil
}

func (r *applicationAlertConfigResourceFrameworkImpl) UpdateState(ctx context.Context, data *restapi.ApplicationAlertConfig, state *tfsdk.State) diag.Diagnostics {
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
	if data.GracePeriod > 0 {
		model.GracePeriod = types.Int64Value(int64(data.GracePeriod))
	}

	// Handle tag filter
	if data.TagFilterExpression != nil {
		model.TagFilter = types.StringValue(tagfilter.RenderExpression(data.TagFilterExpression))
	}

	// Handle severity (deprecated but supported for backward compatibility)
	if data.Severity > 0 {
		severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(data.Severity)
		if err == nil {
			model.Severity = types.StringValue(severity)
		}
	}

	// Handle alert channel IDs (deprecated but supported for backward compatibility)
	if len(data.AlertChannelIDs) > 0 {
		elements := make([]attr.Value, len(data.AlertChannelIDs))
		for i, id := range data.AlertChannelIDs {
			elements[i] = types.StringValue(id)
		}
		model.AlertChannelIDs = types.SetValueMust(types.StringType, elements)
	}

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
	}

	// Handle applications
	if len(data.Applications) > 0 {
		appElements := make([]attr.Value, len(data.Applications))
		appIndex := 0
		for _, app := range data.Applications {
			appModel := ApplicationModel{
				ApplicationID: types.StringValue(app.ApplicationID),
				Inclusive:     types.BoolValue(app.Inclusive),
			}

			if len(app.Services) > 0 {
				svcElements := make([]attr.Value, len(app.Services))
				svcIndex := 0
				for _, svc := range app.Services {
					svcModel := ServiceModel{
						ServiceID: types.StringValue(svc.ServiceID),
						Inclusive: types.BoolValue(svc.Inclusive),
					}

					if len(svc.Endpoints) > 0 {
						epElements := make([]attr.Value, len(svc.Endpoints))
						epIndex := 0
						for _, ep := range svc.Endpoints {
							epModel := EndpointModel{
								EndpointID: types.StringValue(ep.EndpointID),
								Inclusive:  types.BoolValue(ep.Inclusive),
							}
							epObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
								"endpoint_id": types.StringType,
								"inclusive":   types.BoolType,
							}, epModel)
							if diags.HasError() {
								return diags
							}
							epElements[epIndex] = epObj
							epIndex++
						}
						svcModel.Endpoints = types.SetValueMust(types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"endpoint_id": types.StringType,
								"inclusive":   types.BoolType,
							},
						}, epElements)
					}

					svcObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
						"service_id": types.StringType,
						"inclusive":  types.BoolType,
						"endpoint": types.SetType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"endpoint_id": types.StringType,
								"inclusive":   types.BoolType,
							},
						}},
					}, svcModel)
					if diags.HasError() {
						return diags
					}
					svcElements[svcIndex] = svcObj
					svcIndex++
				}
				appModel.Services = types.SetValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"service_id": types.StringType,
						"inclusive":  types.BoolType,
						"endpoint": types.SetType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"endpoint_id": types.StringType,
								"inclusive":   types.BoolType,
							},
						}},
					},
				}, svcElements)
			}

			appObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"application_id": types.StringType,
				"inclusive":      types.BoolType,
				"service": types.SetType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"service_id": types.StringType,
						"inclusive":  types.BoolType,
						"endpoint": types.SetType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"endpoint_id": types.StringType,
								"inclusive":   types.BoolType,
							},
						}},
					},
				}},
			}, appModel)
			if diags.HasError() {
				return diags
			}
			appElements[appIndex] = appObj
			appIndex++
		}
		model.Applications = types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"application_id": types.StringType,
				"inclusive":      types.BoolType,
				"service": types.SetType{ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"service_id": types.StringType,
						"inclusive":  types.BoolType,
						"endpoint": types.SetType{ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"endpoint_id": types.StringType,
								"inclusive":   types.BoolType,
							},
						}},
					},
				}},
			},
		}, appElements)
	}

	// Handle custom payload fields
	if len(data.CustomerPayloadFields) > 0 {
		elements := make([]attr.Value, len(data.CustomerPayloadFields))
		for i, field := range data.CustomerPayloadFields {
			fieldObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
				"type":  types.StringType,
			}, map[string]interface{}{
				"key":   field.Key,
				"value": field.Value,
				"type":  string(field.Type),
			})
			if diags.HasError() {
				return diags
			}
			elements[i] = fieldObj
		}
		model.CustomPayloadFields = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
				"type":  types.StringType,
			},
		}, elements)
	}

	// Handle threshold
	if data.Threshold != nil {
		thresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"type":  types.StringType,
			"value": types.Float64Type,
		}, map[string]interface{}{
			"type":  data.Threshold.Type,
			"value": data.Threshold.Value,
		})
		if diags.HasError() {
			return diags
		}
		model.Threshold = thresholdObj
	}

	// Handle rule (deprecated but supported for backward compatibility)
	if data.Rule != nil {
		ruleModel := RuleModel{}

		// Handle error rate rule
		if data.Rule.ErrorRate != nil {
			errorRateElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(data.Rule.ErrorRate.MetricName),
					"aggregation": types.StringValue(data.Rule.ErrorRate.Aggregation),
				}),
			}
			ruleModel.ErrorRate = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, errorRateElements)
		}

		// Handle errors rule
		if data.Rule.Errors != nil {
			errorsElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(data.Rule.Errors.MetricName),
					"aggregation": types.StringValue(data.Rule.Errors.Aggregation),
				}),
			}
			ruleModel.Errors = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, errorsElements)
		}

		// Handle logs rule
		if data.Rule.Logs != nil {
			logsElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
					"level":       types.StringType,
					"message":     types.StringType,
					"operator":    types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(data.Rule.Logs.MetricName),
					"aggregation": types.StringValue(data.Rule.Logs.Aggregation),
					"level":       types.StringValue(data.Rule.Logs.Level),
					"message":     types.StringValue(data.Rule.Logs.Message),
					"operator":    types.StringValue(data.Rule.Logs.Operator),
				}),
			}
			ruleModel.Logs = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
					"level":       types.StringType,
					"message":     types.StringType,
					"operator":    types.StringType,
				},
			}, logsElements)
		}

		// Handle slowness rule
		if data.Rule.Slowness != nil {
			slownessElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(data.Rule.Slowness.MetricName),
					"aggregation": types.StringValue(data.Rule.Slowness.Aggregation),
				}),
			}
			ruleModel.Slowness = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, slownessElements)
		}

		// Handle status code rule
		if data.Rule.StatusCode != nil {
			statusCodeElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"metric_name":       types.StringType,
					"aggregation":       types.StringType,
					"status_code_start": types.Int64Type,
					"status_code_end":   types.Int64Type,
				}, map[string]attr.Value{
					"metric_name":       types.StringValue(data.Rule.StatusCode.MetricName),
					"aggregation":       types.StringValue(data.Rule.StatusCode.Aggregation),
					"status_code_start": types.Int64Value(int64(data.Rule.StatusCode.StatusCodeStart)),
					"status_code_end":   types.Int64Value(int64(data.Rule.StatusCode.StatusCodeEnd)),
				}),
			}
			ruleModel.StatusCode = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":       types.StringType,
					"aggregation":       types.StringType,
					"status_code_start": types.Int64Type,
					"status_code_end":   types.Int64Type,
				},
			}, statusCodeElements)
		}

		// Handle throughput rule
		if data.Rule.Throughput != nil {
			throughputElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				}, map[string]attr.Value{
					"metric_name": types.StringValue(data.Rule.Throughput.MetricName),
					"aggregation": types.StringValue(data.Rule.Throughput.Aggregation),
				}),
			}
			ruleModel.Throughput = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, throughputElements)
		}

		ruleObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"error_rate":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
			"errors":      types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
			"logs":        types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}}},
			"slowness":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
			"status_code": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}}},
			"throughput":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
		}, ruleModel)
		if diags.HasError() {
			return diags
		}

		model.Rule = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"error_rate":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
				"errors":      types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
				"logs":        types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}}},
				"slowness":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
				"status_code": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}}},
				"throughput":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
			},
		}, []attr.Value{ruleObj})
	}

	// Handle rules (new format with multiple thresholds and severity levels)
	if len(data.Rules) > 0 {
		rulesElements := make([]attr.Value, len(data.Rules))
		for i, ruleWithThreshold := range data.Rules {
			// Create rule model
			ruleModel := RuleModel{}

			// Handle error rate rule
			if ruleWithThreshold.Rule.ErrorRate != nil {
				errorRateElements := []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.ErrorRate.MetricName),
						"aggregation": types.StringValue(ruleWithThreshold.Rule.ErrorRate.Aggregation),
					}),
				}
				ruleModel.ErrorRate = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					},
				}, errorRateElements)
			}

			// Handle errors rule
			if ruleWithThreshold.Rule.Errors != nil {
				errorsElements := []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.Errors.MetricName),
						"aggregation": types.StringValue(ruleWithThreshold.Rule.Errors.Aggregation),
					}),
				}
				ruleModel.Errors = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					},
				}, errorsElements)
			}

			// Handle logs rule
			if ruleWithThreshold.Rule.Logs != nil {
				logsElements := []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
						"level":       types.StringType,
						"message":     types.StringType,
						"operator":    types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.Logs.MetricName),
						"aggregation": types.StringValue(ruleWithThreshold.Rule.Logs.Aggregation),
						"level":       types.StringValue(ruleWithThreshold.Rule.Logs.Level),
						"message":     types.StringValue(ruleWithThreshold.Rule.Logs.Message),
						"operator":    types.StringValue(ruleWithThreshold.Rule.Logs.Operator),
					}),
				}
				ruleModel.Logs = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
						"level":       types.StringType,
						"message":     types.StringType,
						"operator":    types.StringType,
					},
				}, logsElements)
			}

			// Handle slowness rule
			if ruleWithThreshold.Rule.Slowness != nil {
				slownessElements := []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.Slowness.MetricName),
						"aggregation": types.StringValue(ruleWithThreshold.Rule.Slowness.Aggregation),
					}),
				}
				ruleModel.Slowness = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					},
				}, slownessElements)
			}

			// Handle status code rule
			if ruleWithThreshold.Rule.StatusCode != nil {
				statusCodeElements := []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name":       types.StringType,
						"aggregation":       types.StringType,
						"status_code_start": types.Int64Type,
						"status_code_end":   types.Int64Type,
					}, map[string]attr.Value{
						"metric_name":       types.StringValue(ruleWithThreshold.Rule.StatusCode.MetricName),
						"aggregation":       types.StringValue(ruleWithThreshold.Rule.StatusCode.Aggregation),
						"status_code_start": types.Int64Value(int64(ruleWithThreshold.Rule.StatusCode.StatusCodeStart)),
						"status_code_end":   types.Int64Value(int64(ruleWithThreshold.Rule.StatusCode.StatusCodeEnd)),
					}),
				}
				ruleModel.StatusCode = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name":       types.StringType,
						"aggregation":       types.StringType,
						"status_code_start": types.Int64Type,
						"status_code_end":   types.Int64Type,
					},
				}, statusCodeElements)
			}

			// Handle throughput rule
			if ruleWithThreshold.Rule.Throughput != nil {
				throughputElements := []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.Throughput.MetricName),
						"aggregation": types.StringValue(ruleWithThreshold.Rule.Throughput.Aggregation),
					}),
				}
				ruleModel.Throughput = types.ListValueMust(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					},
				}, throughputElements)
			}

			ruleObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"error_rate":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
				"errors":      types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
				"logs":        types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}}},
				"slowness":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
				"status_code": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}}},
				"throughput":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
			}, ruleModel)
			if diags.HasError() {
				return diags
			}

			// Create thresholds map
			thresholdElements := make(map[string]attr.Value)
			for severity, threshold := range ruleWithThreshold.Thresholds {
				thresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
					"value": types.Float64Type,
				}, map[string]interface{}{
					"value": threshold.Value,
				})
				if diags.HasError() {
					return diags
				}
				thresholdElements[severity] = thresholdObj
			}

			// Create rule with threshold object
			ruleWithThresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"rule": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"error_rate":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
						"errors":      types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
						"logs":        types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}}},
						"slowness":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
						"status_code": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}}},
						"throughput":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
					},
				},
				"threshold_operator": types.StringType,
				"thresholds": types.MapType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"value": types.Float64Type,
						},
					},
				},
			}, map[string]interface{}{
				"rule":               ruleObj,
				"threshold_operator": ruleWithThreshold.ThresholdOperator,
				"thresholds":         types.MapValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"value": types.Float64Type}}, thresholdElements),
			})
			if diags.HasError() {
				return diags
			}
			rulesElements[i] = ruleWithThresholdObj
		}

		model.Rules = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"rule": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"error_rate":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
						"errors":      types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
						"logs":        types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "level": types.StringType, "message": types.StringType, "operator": types.StringType}}},
						"slowness":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
						"status_code": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType, "status_code_start": types.Int64Type, "status_code_end": types.Int64Type}}},
						"throughput":  types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"metric_name": types.StringType, "aggregation": types.StringType}}},
					},
				},
				"threshold_operator": types.StringType,
				"thresholds": types.MapType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"value": types.Float64Type,
						},
					},
				},
			},
		}, rulesElements)
	}

	// Handle time threshold
	if data.TimeThreshold != nil {
		timeThresholdModel := AppAlertTimeThresholdModel{}

		// Handle request impact
		if data.TimeThreshold.RequestImpact != nil {
			requestImpactElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"time_window": types.Int64Type,
					"requests":    types.Int64Type,
				}, map[string]attr.Value{
					"time_window": types.Int64Value(int64(data.TimeThreshold.RequestImpact.TimeWindow)),
					"requests":    types.Int64Value(int64(data.TimeThreshold.RequestImpact.Requests)),
				}),
			}
			timeThresholdModel.RequestImpact = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"time_window": types.Int64Type,
					"requests":    types.Int64Type,
				},
			}, requestImpactElements)
		}

		// Handle violations in period
		if data.TimeThreshold.ViolationsInPeriod != nil {
			violationsInPeriodElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"time_window": types.Int64Type,
					"violations":  types.Int64Type,
				}, map[string]attr.Value{
					"time_window": types.Int64Value(int64(data.TimeThreshold.ViolationsInPeriod.TimeWindow)),
					"violations":  types.Int64Value(int64(data.TimeThreshold.ViolationsInPeriod.Violations)),
				}),
			}
			timeThresholdModel.ViolationsInPeriod = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"time_window": types.Int64Type,
					"violations":  types.Int64Type,
				},
			}, violationsInPeriodElements)
		}

		// Handle violations in sequence
		if data.TimeThreshold.ViolationsInSequence != nil {
			violationsInSequenceElements := []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"time_window": types.Int64Type,
				}, map[string]attr.Value{
					"time_window": types.Int64Value(int64(data.TimeThreshold.ViolationsInSequence.TimeWindow)),
				}),
			}
			timeThresholdModel.ViolationsInSequence = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"time_window": types.Int64Type,
				},
			}, violationsInSequenceElements)
		}

		timeThresholdObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"request_impact":         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type, "requests": types.Int64Type}}},
			"violations_in_period":   types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type, "violations": types.Int64Type}}},
			"violations_in_sequence": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type}}},
		}, timeThresholdModel)
		if diags.HasError() {
			return diags
		}

		model.TimeThreshold = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"request_impact":         types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type, "requests": types.Int64Type}}},
				"violations_in_period":   types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type, "violations": types.Int64Type}}},
				"violations_in_sequence": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"time_window": types.Int64Type}}},
			},
		}, []attr.Value{timeThresholdObj})
	}

	return state.Set(ctx, model)
}
