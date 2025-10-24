package instana

import (
	"context"
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	Threshold           types.List   `tfsdk:"threshold"`
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
	Thresholds        types.List   `tfsdk:"thresholds"`
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
								ApplicationAlertConfigFieldThresholds: schema.ListNestedBlock{
									Description: "Threshold configuration for different severity levels",
									NestedObject: schema.NestedBlockObject{
										Blocks: map[string]schema.Block{
											LogAlertConfigFieldWarning:  StaticAndAdaptiveThresholdBlockSchema(),
											LogAlertConfigFieldCritical: StaticAndAdaptiveThresholdBlockSchema(),
										},
									},
									Validators: []validator.List{
										listvalidator.SizeAtMost(1),
									},
								},
							},
						},
					},
					ResourceFieldThreshold: schema.ListNestedBlock{
						Description: "Threshold configuration for different severity levels",
						NestedObject: schema.NestedBlockObject{
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
										"value": schema.Float64Attribute{
											Description: "The numeric value for the static threshold.",
											Optional:    true,
										},
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
										"deviation_factor": schema.Float64Attribute{
											Description: "The numeric value for the deviation factor.",
											Optional:    true,
										},
										"adaptability": schema.Float64Attribute{
											Description: "The numeric value for the adaptability.",
											Optional:    true,
										},
										"seasonality": schema.StringAttribute{
											Description: "Value for the seasonality.",
											Optional:    true,
										},
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtMost(1),
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
			SkipIDGeneration: true,
			SchemaVersion:    1,
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

func (r *applicationAlertConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, obj *restapi.ApplicationAlertConfig) diag.Diagnostics {
	var diags diag.Diagnostics

	// Delegate to the resource implementation
	resource, resourceDiags := r.NewResource(ctx, nil)
	if resourceDiags.HasError() {
		diags.Append(resourceDiags...)
		return diags
	}

	// Update the state with the object
	return resource.UpdateState(ctx, state, obj)
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
	UpdateState(ctx context.Context, state *tfsdk.State, obj T) diag.Diagnostics
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
		customerPayloadFields, payloadDiags = MapCustomPayloadFieldsToAPIObject(ctx, model.CustomPayloadFields)
		if payloadDiags.HasError() {
			diags.Append(payloadDiags...)
			return nil, diags
		}
		result.CustomerPayloadFields = customerPayloadFields
	}

	// Handle threshold
	if !model.Threshold.IsNull() && !model.Threshold.IsUnknown() {
		// We need to convert from ThresholdRule to Threshold
		// First, get the operator from the schema
		var thresholdElements []types.Object
		diags.Append(model.Threshold.ElementsAs(ctx, &thresholdElements, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(thresholdElements) > 0 {
			// Create a new Threshold object
			threshold := &restapi.Threshold{}
			log.Printf("static threshold elements : %+v\n", thresholdElements)

			// Check for static threshold
			if staticVal, ok := thresholdElements[0].Attributes()[ThresholdFieldStatic]; ok && !staticVal.IsNull() && !staticVal.IsUnknown() {
				log.Printf("static threshold value type: %T\n", staticVal)

				// Extract the operator and value directly from the static object
				var staticObj struct {
					Operator types.String  `tfsdk:"operator"`
					Value    types.Float64 `tfsdk:"value"`
				}

				staticObject, ok := staticVal.(types.Object)
				if ok && !staticObject.IsNull() && !staticObject.IsUnknown() {
					log.Printf("static threshold individual: %+v\n", staticObject)

					diags.Append(staticObject.As(ctx, &staticObj, basetypes.ObjectAsOptions{})...)
					if diags.HasError() {
						return nil, diags
					}

					// Set the threshold type and operator
					threshold.Type = "staticThreshold"
					threshold.Operator = restapi.ThresholdOperator(staticObj.Operator.ValueString())

					// Set the value
					if !staticObj.Value.IsNull() && !staticObj.Value.IsUnknown() {
						value := staticObj.Value.ValueFloat64()
						threshold.Value = &value
					}
					log.Printf("static threshold final: %+v\n", threshold)
				}
			}

			// Check for adaptive baseline threshold
			if adaptiveVal, ok := thresholdElements[0].Attributes()[ThresholdFieldAdaptiveBaseline]; ok && !adaptiveVal.IsNull() && !adaptiveVal.IsUnknown() {
				log.Printf("adaptive threshold value type: %T\n", adaptiveVal)

				// Extract the operator, deviation factor, adaptability, and seasonality directly
				var adaptiveObj struct {
					Operator        types.String  `tfsdk:"operator"`
					DeviationFactor types.Float64 `tfsdk:"deviation_factor"`
					Adaptability    types.Float64 `tfsdk:"adaptability"`
					Seasonality     types.String  `tfsdk:"seasonality"`
				}

				adaptiveObject, ok := adaptiveVal.(types.Object)
				if ok && !adaptiveObject.IsNull() && !adaptiveObject.IsUnknown() {
					log.Printf("adaptive threshold individual: %+v\n", adaptiveObject)

					diags.Append(adaptiveObject.As(ctx, &adaptiveObj, basetypes.ObjectAsOptions{})...)
					if diags.HasError() {
						return nil, diags
					}

					// Set the threshold type and operator
					threshold.Type = "adaptiveBaseline"
					threshold.Operator = restapi.ThresholdOperator(adaptiveObj.Operator.ValueString())

					// Set the deviation factor
					if !adaptiveObj.DeviationFactor.IsNull() && !adaptiveObj.DeviationFactor.IsUnknown() {
						deviationFactor := float32(adaptiveObj.DeviationFactor.ValueFloat64())
						threshold.DeviationFactor = &deviationFactor
					}

					// Set the seasonality
					if !adaptiveObj.Seasonality.IsNull() && !adaptiveObj.Seasonality.IsUnknown() {
						seasonality := restapi.ThresholdSeasonality(adaptiveObj.Seasonality.ValueString())
						threshold.Seasonality = &seasonality
					}

					log.Printf("adaptive threshold final: %+v\n", threshold)
				}
			}

			result.Threshold = threshold
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
					result.Rule.AlertType = ApplicationAlertConfigFieldRuleErrorRate
					result.Rule.MetricName = errorRateRules[0].MetricName.ValueString()
					result.Rule.Aggregation = restapi.Aggregation(errorRateRules[0].Aggregation.ValueString())
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
					result.Rule.AlertType = ApplicationAlertConfigFieldRuleErrors
					result.Rule.MetricName = errorsRules[0].MetricName.ValueString()
					result.Rule.Aggregation = restapi.Aggregation(errorsRules[0].Aggregation.ValueString())
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
					result.Rule.AlertType = ApplicationAlertConfigFieldRuleLogs
					result.Rule.MetricName = logsRules[0].MetricName.ValueString()
					result.Rule.Aggregation = restapi.Aggregation(logsRules[0].Aggregation.ValueString())

					// Set additional fields for logs
					level := restapi.LogLevel(logsRules[0].Level.ValueString())
					result.Rule.Level = &level

					message := logsRules[0].Message.ValueString()
					result.Rule.Message = &message

					operator := restapi.ExpressionOperator(logsRules[0].Operator.ValueString())
					result.Rule.Operator = &operator
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
					result.Rule.AlertType = ApplicationAlertConfigFieldRuleSlowness
					result.Rule.MetricName = slownessRules[0].MetricName.ValueString()
					result.Rule.Aggregation = restapi.Aggregation(slownessRules[0].Aggregation.ValueString())
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
					result.Rule.AlertType = ApplicationAlertConfigFieldRuleStatusCode
					result.Rule.MetricName = statusCodeRules[0].MetricName.ValueString()
					result.Rule.Aggregation = restapi.Aggregation(statusCodeRules[0].Aggregation.ValueString())

					// Set additional fields for status code
					statusCodeStart := int32(statusCodeRules[0].StatusCodeStart.ValueInt64())
					result.Rule.StatusCodeStart = &statusCodeStart

					statusCodeEnd := int32(statusCodeRules[0].StatusCodeEnd.ValueInt64())
					result.Rule.StatusCodeEnd = &statusCodeEnd
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
					result.Rule.AlertType = ApplicationAlertConfigFieldRuleThroughput
					result.Rule.MetricName = throughputRules[0].MetricName.ValueString()
					result.Rule.Aggregation = restapi.Aggregation(throughputRules[0].Aggregation.ValueString())
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
						result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleErrorRate
						result.Rules[i].Rule.MetricName = errorRateRules[0].MetricName.ValueString()
						result.Rules[i].Rule.Aggregation = restapi.Aggregation(errorRateRules[0].Aggregation.ValueString())
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
						result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleErrors
						result.Rules[i].Rule.MetricName = errorsRules[0].MetricName.ValueString()
						result.Rules[i].Rule.Aggregation = restapi.Aggregation(errorsRules[0].Aggregation.ValueString())
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
						result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleLogs
						result.Rules[i].Rule.MetricName = logsRules[0].MetricName.ValueString()
						result.Rules[i].Rule.Aggregation = restapi.Aggregation(logsRules[0].Aggregation.ValueString())

						// Set additional fields for logs
						level := restapi.LogLevel(logsRules[0].Level.ValueString())
						result.Rules[i].Rule.Level = &level

						message := logsRules[0].Message.ValueString()
						result.Rules[i].Rule.Message = &message

						operator := restapi.ExpressionOperator(logsRules[0].Operator.ValueString())
						result.Rules[i].Rule.Operator = &operator
					}
				}

				// Handle slowness rule
				if !rule.Slowness.IsNull() && !rule.Slowness.IsUnknown() {
					log.Printf("Slowness : %+v\n", rule.Slowness)
					var slownessRules []RuleConfigModel
					diags = rule.Slowness.ElementsAs(ctx, &slownessRules, false)
					log.Printf("Slowness mapped : %+v\n", slownessRules)
					if diags.HasError() {
						return nil, diags
					}
					if len(slownessRules) > 0 {
						result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleSlowness
						result.Rules[i].Rule.MetricName = slownessRules[0].MetricName.ValueString()
						result.Rules[i].Rule.Aggregation = restapi.Aggregation(slownessRules[0].Aggregation.ValueString())
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
						result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleStatusCode
						result.Rules[i].Rule.MetricName = statusCodeRules[0].MetricName.ValueString()
						result.Rules[i].Rule.Aggregation = restapi.Aggregation(statusCodeRules[0].Aggregation.ValueString())

						// Set additional fields for status code
						statusCodeStart := int32(statusCodeRules[0].StatusCodeStart.ValueInt64())
						result.Rules[i].Rule.StatusCodeStart = &statusCodeStart

						statusCodeEnd := int32(statusCodeRules[0].StatusCodeEnd.ValueInt64())
						result.Rules[i].Rule.StatusCodeEnd = &statusCodeEnd
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
						result.Rules[i].Rule.AlertType = ApplicationAlertConfigFieldRuleThroughput
						result.Rules[i].Rule.MetricName = throughputRules[0].MetricName.ValueString()
						result.Rules[i].Rule.Aggregation = restapi.Aggregation(throughputRules[0].Aggregation.ValueString())
					}
				}
			}

			// Handle thresholds
			var thresholdMap map[restapi.AlertSeverity]restapi.ThresholdRule
			var thresholdDiags diag.Diagnostics

			if !ruleWithThreshold.Thresholds.IsNull() && !ruleWithThreshold.Thresholds.IsUnknown() {
				thresholdMap, thresholdDiags = MapThresholdsFromState(ctx, ruleWithThreshold.Thresholds)
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
					result.TimeThreshold.Type = "requestImpact"
					result.TimeThreshold.TimeWindow = requestImpacts[0].TimeWindow.ValueInt64()
					result.TimeThreshold.Requests = int(requestImpacts[0].Requests.ValueInt64())

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
					result.TimeThreshold.Type = "violationsInPeriod"
					result.TimeThreshold.TimeWindow = violationsInPeriods[0].TimeWindow.ValueInt64()
					result.TimeThreshold.Violations = int(violationsInPeriods[0].Violations.ValueInt64())
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
					result.TimeThreshold.Type = "violationsInSequence"
					result.TimeThreshold.TimeWindow = violationsInSequences[0].TimeWindow.ValueInt64()
				}
			}
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

func (r *applicationAlertConfigResourceFrameworkImpl) UpdateState(ctx context.Context, state *tfsdk.State, data *restapi.ApplicationAlertConfig) diag.Diagnostics {
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
		model.GracePeriod = types.Int64Value(*data.GracePeriod)
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
	customPayloadFieldsList, payloadDiags := CustomPayloadFieldsToTerraform(ctx, data.CustomerPayloadFields)
	if payloadDiags.HasError() {
		return payloadDiags
	}
	model.CustomPayloadFields = customPayloadFieldsList

	if diags.HasError() {
		return diags
	}
	log.Printf("Before Threshold stage")
	// Handle threshold
	if data.Threshold != nil {
		// Map thresholds
		thresholdObj := map[string]attr.Value{}

		// Create threshold object based on the threshold type
		if data.Threshold != nil {
			// Map static threshold
			if data.Threshold.Type == "staticThreshold" {
				// Create static threshold object
				staticObj := map[string]attr.Value{}

				// Add operator
				staticObj["operator"] = types.StringValue(string(data.Threshold.Operator))

				// Add value
				if data.Threshold.Value != nil {
					staticObj["value"] = types.Float64Value(*data.Threshold.Value)
				} else {
					staticObj["value"] = types.Float64Null()
				}

				// Create static object value
				staticObjVal, staticObjDiags := types.ObjectValue(
					map[string]attr.Type{
						"operator": types.StringType,
						"value":    types.Float64Type,
					},
					staticObj,
				)
				if staticObjDiags.HasError() {
					return staticObjDiags
				}

				// Create static list
				staticList, staticListDiags := types.ListValue(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"operator": types.StringType,
							"value":    types.Float64Type,
						},
					},
					[]attr.Value{staticObjVal},
				)
				if staticListDiags.HasError() {
					return staticListDiags
				}

				thresholdObj[ThresholdFieldStatic] = staticList
			} else if data.Threshold.Type == "adaptiveBaseline" {
				// Create adaptive baseline threshold object
				adaptiveObj := map[string]attr.Value{}

				// Add operator
				adaptiveObj["operator"] = types.StringValue(string(data.Threshold.Operator))

				// Add deviation factor
				if data.Threshold.DeviationFactor != nil {
					adaptiveObj["deviation_factor"] = types.Float64Value(float64(*data.Threshold.DeviationFactor))
				} else {
					adaptiveObj["deviation_factor"] = types.Float64Null()
				}

				// Add adaptability (assuming it's stored in Value for now)
				if data.Threshold.Value != nil {
					adaptiveObj["adaptability"] = types.Float64Value(*data.Threshold.Value)
				} else {
					adaptiveObj["adaptability"] = types.Float64Null()
				}

				// Add seasonality
				if data.Threshold.Seasonality != nil {
					adaptiveObj["seasonality"] = types.StringValue(string(*data.Threshold.Seasonality))
				} else {
					adaptiveObj["seasonality"] = types.StringNull()
				}

				// Create adaptive baseline object value
				adaptiveObjVal, adaptiveObjDiags := types.ObjectValue(
					map[string]attr.Type{
						"operator":         types.StringType,
						"deviation_factor": types.Float64Type,
						"adaptability":     types.Float64Type,
						"seasonality":      types.StringType,
					},
					adaptiveObj,
				)
				if adaptiveObjDiags.HasError() {
					return adaptiveObjDiags
				}

				// Create adaptive baseline list
				adaptiveList, adaptiveListDiags := types.ListValue(
					types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"operator":         types.StringType,
							"deviation_factor": types.Float64Type,
							"adaptability":     types.Float64Type,
							"seasonality":      types.StringType,
						},
					},
					[]attr.Value{adaptiveObjVal},
				)
				if adaptiveListDiags.HasError() {
					return adaptiveListDiags
				}

				thresholdObj[ThresholdFieldAdaptiveBaseline] = adaptiveList
			}
		}

		// Create threshold object value with the appropriate attribute types
		thresholdAttrTypes := map[string]attr.Type{}

		// Add static threshold attribute type if present
		if _, ok := thresholdObj[ThresholdFieldStatic]; ok {
			thresholdAttrTypes[ThresholdFieldStatic] = types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"operator": types.StringType,
						"value":    types.Float64Type,
					},
				},
			}
		}

		// Add adaptive baseline attribute type if present
		if _, ok := thresholdObj[ThresholdFieldAdaptiveBaseline]; ok {
			thresholdAttrTypes[ThresholdFieldAdaptiveBaseline] = types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"operator":         types.StringType,
						"deviation_factor": types.Float64Type,
						"adaptability":     types.Float64Type,
						"seasonality":      types.StringType,
					},
				},
			}
		}

		// Create the object value
		thresholdObjVal, thresholdObjDiags := types.ObjectValue(
			thresholdAttrTypes,
			thresholdObj,
		)
		if thresholdObjDiags.HasError() {
			return thresholdObjDiags
		}

		// Create the list value
		thresholdList, thresholdListDiags := types.ListValue(
			types.ObjectType{
				AttrTypes: thresholdAttrTypes,
			},
			[]attr.Value{thresholdObjVal},
		)
		if thresholdListDiags.HasError() {
			return thresholdListDiags
		}

		model.Threshold = thresholdList
	}

	if diags.HasError() {
		return diags
	}
	log.Printf("Before Rule stage")
	// Handle rule (deprecated but supported for backward compatibility)
	if data.Rule != nil {
		ruleModel := RuleModel{}
		errorRateElements := make([]attr.Value, 0)
		errorsElements := make([]attr.Value, 0)
		logsElements := make([]attr.Value, 0)
		slownessElements := make([]attr.Value, 0)
		statusCodeElements := make([]attr.Value, 0)
		throughputElements := make([]attr.Value, 0)

		// Set rule model fields based on the AlertType
		if data.Rule != nil {
			switch data.Rule.AlertType {
			case ApplicationAlertConfigFieldRuleErrorRate:
				errorRateElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(data.Rule.MetricName),
						"aggregation": types.StringValue(string(data.Rule.Aggregation)),
					}),
				}

			case ApplicationAlertConfigFieldRuleErrors:
				errorsElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(data.Rule.MetricName),
						"aggregation": types.StringValue(string(data.Rule.Aggregation)),
					}),
				}

			case ApplicationAlertConfigFieldRuleLogs:
				// For logs, we need to handle additional fields
				level := ""
				message := ""
				operator := ""

				if data.Rule.Level != nil {
					level = string(*data.Rule.Level)
				}
				if data.Rule.Message != nil {
					message = *data.Rule.Message
				}
				if data.Rule.Operator != nil {
					operator = string(*data.Rule.Operator)
				}

				logsElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
						"level":       types.StringType,
						"message":     types.StringType,
						"operator":    types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(data.Rule.MetricName),
						"aggregation": types.StringValue(string(data.Rule.Aggregation)),
						"level":       types.StringValue(level),
						"message":     types.StringValue(message),
						"operator":    types.StringValue(operator),
					}),
				}

			case ApplicationAlertConfigFieldRuleSlowness:
				slownessElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(data.Rule.MetricName),
						"aggregation": types.StringValue(string(data.Rule.Aggregation)),
					}),
				}

			case ApplicationAlertConfigFieldRuleStatusCode:
				// For status code, we need to handle additional fields
				statusCodeStart := int64(0)
				statusCodeEnd := int64(0)

				if data.Rule.StatusCodeStart != nil {
					statusCodeStart = int64(*data.Rule.StatusCodeStart)
				}
				if data.Rule.StatusCodeEnd != nil {
					statusCodeEnd = int64(*data.Rule.StatusCodeEnd)
				}

				statusCodeElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name":       types.StringType,
						"aggregation":       types.StringType,
						"status_code_start": types.Int64Type,
						"status_code_end":   types.Int64Type,
					}, map[string]attr.Value{
						"metric_name":       types.StringValue(data.Rule.MetricName),
						"aggregation":       types.StringValue(string(data.Rule.Aggregation)),
						"status_code_start": types.Int64Value(statusCodeStart),
						"status_code_end":   types.Int64Value(statusCodeEnd),
					}),
				}

			case ApplicationAlertConfigFieldRuleThroughput:
				throughputElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(data.Rule.MetricName),
						"aggregation": types.StringValue(string(data.Rule.Aggregation)),
					}),
				}

			}
		}

		ruleModel.ErrorRate = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, errorRateElements)

		ruleModel.Errors = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, errorsElements)

		ruleModel.Logs = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
				"level":       types.StringType,
				"message":     types.StringType,
				"operator":    types.StringType,
			},
		}, logsElements)

		ruleModel.Slowness = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, slownessElements)

		ruleModel.StatusCode = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name":       types.StringType,
				"aggregation":       types.StringType,
				"status_code_start": types.Int64Type,
				"status_code_end":   types.Int64Type,
			},
		}, statusCodeElements)

		ruleModel.Throughput = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
			},
		}, throughputElements)

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
			errorRateElements := make([]attr.Value, 0)
			errorsElements := make([]attr.Value, 0)
			logsElements := make([]attr.Value, 0)
			slownessElements := make([]attr.Value, 0)
			statusCodeElements := make([]attr.Value, 0)
			throughputElements := make([]attr.Value, 0)

			// Set rule model fields based on the AlertType
			switch ruleWithThreshold.Rule.AlertType {
			case ApplicationAlertConfigFieldRuleErrorRate:
				errorRateElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
						"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					}),
				}

			case ApplicationAlertConfigFieldRuleErrors:
				errorsElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
						"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					}),
				}

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

				logsElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
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
					}),
				}

			case ApplicationAlertConfigFieldRuleSlowness:
				slownessElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
						"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					}),
				}

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

				statusCodeElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name":       types.StringType,
						"aggregation":       types.StringType,
						"status_code_start": types.Int64Type,
						"status_code_end":   types.Int64Type,
					}, map[string]attr.Value{
						"metric_name":       types.StringValue(ruleWithThreshold.Rule.MetricName),
						"aggregation":       types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
						"status_code_start": types.Int64Value(statusCodeStart),
						"status_code_end":   types.Int64Value(statusCodeEnd),
					}),
				}

			case ApplicationAlertConfigFieldRuleThroughput:
				throughputElements = []attr.Value{
					types.ObjectValueMust(map[string]attr.Type{
						"metric_name": types.StringType,
						"aggregation": types.StringType,
					}, map[string]attr.Value{
						"metric_name": types.StringValue(ruleWithThreshold.Rule.MetricName),
						"aggregation": types.StringValue(string(ruleWithThreshold.Rule.Aggregation)),
					}),
				}

			}

			ruleModel.ErrorRate = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, errorRateElements)

			ruleModel.Errors = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, errorsElements)

			ruleModel.Logs = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
					"level":       types.StringType,
					"message":     types.StringType,
					"operator":    types.StringType,
				},
			}, logsElements)

			ruleModel.Slowness = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, slownessElements)

			ruleModel.StatusCode = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name":       types.StringType,
					"aggregation":       types.StringType,
					"status_code_start": types.Int64Type,
					"status_code_end":   types.Int64Type,
				},
			}, statusCodeElements)

			ruleModel.Throughput = types.ListValueMust(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metric_name": types.StringType,
					"aggregation": types.StringType,
				},
			}, throughputElements)

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

			// Map thresholds
			thresholdObj := map[string]attr.Value{}

			// Map warning threshold
			warningThreshold, isWarningThresholdPresent := ruleWithThreshold.Thresholds[restapi.WarningSeverity]
			warningThresholdList, warningDiags := MapThresholdToState(ctx, isWarningThresholdPresent, &warningThreshold, []string{"static", "adaptiveBaseline"})
			diags.Append(warningDiags...)
			if diags.HasError() {
				return diags
			}
			thresholdObj[LogAlertConfigFieldWarning] = warningThresholdList

			// Map critical threshold
			criticalThreshold, isCriticalThresholdPresent := ruleWithThreshold.Thresholds[restapi.CriticalSeverity]
			criticalThresholdList, criticalDiags := MapThresholdToState(ctx, isCriticalThresholdPresent, &criticalThreshold, []string{"static", "adaptiveBaseline"})
			diags.Append(criticalDiags...)
			if diags.HasError() {
				return diags
			}
			thresholdObj[LogAlertConfigFieldCritical] = criticalThresholdList

			// Create threshold object value
			thresholdObjVal, thresholdObjDiags := types.ObjectValue(
				map[string]attr.Type{
					LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
					LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
				},
				thresholdObj,
			)
			diags.Append(thresholdObjDiags...)
			if diags.HasError() {
				return diags
			}

			// Add threshold to rule
			thresholdList, thresholdListDiags := types.ListValue(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
						LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
					},
				},
				[]attr.Value{thresholdObjVal},
			)
			diags.Append(thresholdListDiags...)
			if diags.HasError() {
				return diags
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
				"thresholds": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
							LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
						},
					},
				},
			}, map[string]attr.Value{
				"rule":               ruleObj,
				"threshold_operator": types.StringValue(ruleWithThreshold.ThresholdOperator),
				"thresholds":         thresholdList,
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
				"thresholds": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							LogAlertConfigFieldWarning:  GetStaticAndAdaptiveThresholdAttrListTypes(),
							LogAlertConfigFieldCritical: GetStaticAndAdaptiveThresholdAttrListTypes(),
						},
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

		violationsInSequenceElements := []attr.Value{}
		requestImpactElements := []attr.Value{}
		violationsInPeriodElements := []attr.Value{}
		// Determine which time threshold to populate based on the Type field
		switch data.TimeThreshold.Type {
		case "requestImpact":
			requestImpactElements = []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"time_window": types.Int64Type,
					"requests":    types.Int64Type,
				}, map[string]attr.Value{
					"time_window": types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
					"requests":    types.Int64Value(int64(data.TimeThreshold.Requests)),
				}),
			}

		case "violationsInPeriod":
			violationsInPeriodElements = []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"time_window": types.Int64Type,
					"violations":  types.Int64Type,
				}, map[string]attr.Value{
					"time_window": types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
					"violations":  types.Int64Value(int64(data.TimeThreshold.Violations)),
				}),
			}

		case "violationsInSequence":
			violationsInSequenceElements = []attr.Value{
				types.ObjectValueMust(map[string]attr.Type{
					"time_window": types.Int64Type,
				}, map[string]attr.Value{
					"time_window": types.Int64Value(int64(data.TimeThreshold.TimeWindow)),
				}),
			}

		}
		timeThresholdModel.RequestImpact = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"time_window": types.Int64Type,
				"requests":    types.Int64Type,
			},
		}, requestImpactElements)

		timeThresholdModel.ViolationsInPeriod = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"time_window": types.Int64Type,
				"violations":  types.Int64Type,
			},
		}, violationsInPeriodElements)
		timeThresholdModel.ViolationsInSequence = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"time_window": types.Int64Type,
			},
		}, violationsInSequenceElements)

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
	if diags.HasError() {
		return diags
	}
	log.Printf("Reached final stage")
	log.Printf("static threshold elements : %+v\n", model)
	return state.Set(ctx, model)
}
