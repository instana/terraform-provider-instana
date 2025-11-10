package instana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaSliConfigFramework the name of the terraform-provider-instana resource to manage SLI configurations
const ResourceInstanaSliConfigFramework = "sli_config"

const (
	//SliConfigFieldName constant value for the schema field name
	SliConfigFieldName = "name"
	//SliConfigFieldFullName constant value for schema field full_name
	SliConfigFieldFullName = "full_name"
	//SliConfigFieldInitialEvaluationTimestamp constant value for the schema field initial_evaluation_timestamp
	SliConfigFieldInitialEvaluationTimestamp = "initial_evaluation_timestamp"
	//SliConfigFieldMetricConfiguration constant value for the schema field metric_configuration
	SliConfigFieldMetricConfiguration = "metric_configuration"
	//SliConfigFieldMetricName constant value for the schema field metric_configuration.metric_name
	SliConfigFieldMetricName = "metric_name"
	//SliConfigFieldMetricAggregation constant value for the schema field metric_configuration.aggregation
	SliConfigFieldMetricAggregation = "aggregation"
	//SliConfigFieldMetricThreshold constant value for the schema field metric_configuration.threshold
	SliConfigFieldMetricThreshold = "threshold"
	//SliConfigFieldSliEntity constant value for the schema field sli_entity
	SliConfigFieldSliEntity = "sli_entity"
	//SliConfigFieldSliEntityApplicationTimeBased constant value for the schema field sli_entity.application
	SliConfigFieldSliEntityApplicationTimeBased = "application_time_based"
	//SliConfigFieldSliEntityApplicationEventBased constant value for the schema field sli_entity.availability
	SliConfigFieldSliEntityApplicationEventBased = "application_event_based"
	//SliConfigFieldSliEntityWebsiteEventBased constant value for the schema field sli_entity.website_event_based
	SliConfigFieldSliEntityWebsiteEventBased = "website_event_based"
	//SliConfigFieldSliEntityWebsiteTimeBased constant value for the schema field sli_entity.website_time_based
	SliConfigFieldSliEntityWebsiteTimeBased = "website_time_based"
	//SliConfigFieldApplicationID constant value for the schema field sli_entity.*.application_id
	SliConfigFieldApplicationID = "application_id"
	//SliConfigFieldServiceID constant value for the schema field sli_entity.*.service_id
	SliConfigFieldServiceID = "service_id"
	//SliConfigFieldEndpointID constant value for the schema field sli_entity.*.endpoint_id
	SliConfigFieldEndpointID = "endpoint_id"
	//SliConfigFieldWebsiteID constant value for the schema field sli_entity.*.website_id
	SliConfigFieldWebsiteID = "website_id"
	//SliConfigFieldBeaconType constant value for the schema field sli_entity.*.beacon_Type
	SliConfigFieldBeaconType = "beacon_type"
	//SliConfigFieldBoundaryScope constant value for the schema field sli_entity.boundary_scope
	SliConfigFieldBoundaryScope = "boundary_scope"
	//SliConfigFieldBadEventFilterExpression constant value for the schema field sli_entity.*.bad_event_filter_expression
	SliConfigFieldBadEventFilterExpression = "bad_event_filter_expression"
	//SliConfigFieldFilterExpression constant value for the schema field sli_entity.*.filter_expression
	SliConfigFieldFilterExpression = "filter_expression"
	//SliConfigFieldGoodEventFilterExpression constant value for the schema field sli_entity.*.good_event_filter_expression
	SliConfigFieldGoodEventFilterExpression = "good_event_filter_expression"
	//SliConfigFieldIncludeInternal constant value for the schema field sli_entity.*.good_event_filter_expression
	SliConfigFieldIncludeInternal = "include_internal"
	//SliConfigFieldIncludeSynthetic constant value for the schema field sli_entity.*.good_event_filter_expression
	SliConfigFieldIncludeSynthetic = "include_synthetic"
)

var sliConfigSliEntityTypeKeys = []string{
	"sli_entity.0.application_event_based",
	"sli_entity.0.application_time_based",
	"sli_entity.0.website_event_based",
	"sli_entity.0.website_time_based",
}

// SliConfigModel represents the data model for SLI configuration
type SliConfigModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	InitialEvaluationTimestamp types.Int64  `tfsdk:"initial_evaluation_timestamp"`
	MetricConfiguration        types.List   `tfsdk:"metric_configuration"`
	SliEntity                  types.List   `tfsdk:"sli_entity"`
}

// MetricConfigurationModel represents the metric configuration for SLI
type MetricConfigurationModel struct {
	MetricName  types.String  `tfsdk:"metric_name"`
	Aggregation types.String  `tfsdk:"aggregation"`
	Threshold   types.Float64 `tfsdk:"threshold"`
}

// SliEntityModel represents the SLI entity configuration
type SliEntityModel struct {
	ApplicationTimeBased  types.List `tfsdk:"application_time_based"`
	ApplicationEventBased types.List `tfsdk:"application_event_based"`
	WebsiteEventBased     types.List `tfsdk:"website_event_based"`
	WebsiteTimeBased      types.List `tfsdk:"website_time_based"`
}

// ApplicationTimeBasedModel represents the application time based SLI entity
type ApplicationTimeBasedModel struct {
	ApplicationID types.String `tfsdk:"application_id"`
	ServiceID     types.String `tfsdk:"service_id"`
	EndpointID    types.String `tfsdk:"endpoint_id"`
	BoundaryScope types.String `tfsdk:"boundary_scope"`
}

// ApplicationEventBasedModel represents the application event based SLI entity
type ApplicationEventBasedModel struct {
	ApplicationID             types.String `tfsdk:"application_id"`
	BoundaryScope             types.String `tfsdk:"boundary_scope"`
	BadEventFilterExpression  types.String `tfsdk:"bad_event_filter_expression"`
	GoodEventFilterExpression types.String `tfsdk:"good_event_filter_expression"`
	IncludeInternal           types.Bool   `tfsdk:"include_internal"`
	IncludeSynthetic          types.Bool   `tfsdk:"include_synthetic"`
	ServiceID                 types.String `tfsdk:"service_id"`
	EndpointID                types.String `tfsdk:"endpoint_id"`
}

// WebsiteEventBasedModel represents the website event based SLI entity
type WebsiteEventBasedModel struct {
	WebsiteID                 types.String `tfsdk:"website_id"`
	BadEventFilterExpression  types.String `tfsdk:"bad_event_filter_expression"`
	GoodEventFilterExpression types.String `tfsdk:"good_event_filter_expression"`
	BeaconType                types.String `tfsdk:"beacon_type"`
}

// WebsiteTimeBasedModel represents the website time based SLI entity
type WebsiteTimeBasedModel struct {
	WebsiteID        types.String `tfsdk:"website_id"`
	FilterExpression types.String `tfsdk:"filter_expression"`
	BeaconType       types.String `tfsdk:"beacon_type"`
}

type sliConfigResourceFramework struct {
	metaData ResourceMetaDataFramework
}

var applicationTimeBasedObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"application_id": types.StringType,
		"service_id":     types.StringType,
		"endpoint_id":    types.StringType,
		"boundary_scope": types.StringType,
	},
}
var websiteTimeBasedObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"website_id":        types.StringType,
		"filter_expression": types.StringType,
		"beacon_type":       types.StringType,
	},
}
var applicationEventBasedObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"application_id":               types.StringType,
		"boundary_scope":               types.StringType,
		"bad_event_filter_expression":  types.StringType,
		"good_event_filter_expression": types.StringType,
		"include_internal":             types.BoolType,
		"include_synthetic":            types.BoolType,
		"service_id":                   types.StringType,
		"endpoint_id":                  types.StringType,
	},
}
var websiteEventBasedObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"website_id":                   types.StringType,
		"bad_event_filter_expression":  types.StringType,
		"good_event_filter_expression": types.StringType,
		"beacon_type":                  types.StringType,
	},
}

// NewSliConfigResourceHandleFramework creates the resource handle for SLI configuration
func NewSliConfigResourceHandleFramework() ResourceHandleFramework[*restapi.SliConfig] {
	return &sliConfigResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSliConfigFramework,
			Schema: schema.Schema{
				Description: "This resource manages SLI configurations in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the SLI configuration.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the SLI config",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 256),
						},
					},
					"initial_evaluation_timestamp": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Initial evaluation timestamp for the SLI config",
					},
				},
				Blocks: map[string]schema.Block{
					"metric_configuration": schema.ListNestedBlock{
						Description: "Metric configuration for the SLI config",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"metric_name": schema.StringAttribute{
									Required:    true,
									Description: "The metric name for the metric configuration",
								},
								"aggregation": schema.StringAttribute{
									Required:    true,
									Description: "The aggregation type for the metric configuration (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND)",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTRIBUTION", "DISTINCT_COUNT", "SUM_POSITIVE", "PER_SECOND",
										),
									},
								},
								"threshold": schema.Float64Attribute{
									Required:    true,
									Description: "The threshold for the metric configuration",
									Validators: []validator.Float64{
										float64validator.AtLeast(0.000001),
									},
								},
							},
						},
					},
					"sli_entity": schema.ListNestedBlock{
						Description: "The entity to use for the SLI config.",
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"application_time_based": schema.ListNestedBlock{
									Description: "The SLI entity of type application to use for the SLI config",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"application_id": schema.StringAttribute{
												Required:    true, // might be optional
												Description: "The application ID of the entity",
											},
											"service_id": schema.StringAttribute{
												Optional:    true,
												Description: "The service ID of the entity",
											},
											"endpoint_id": schema.StringAttribute{
												Optional:    true,
												Description: "The endpoint ID of the entity",
											},
											"boundary_scope": schema.StringAttribute{
												Required:    true,
												Description: "The boundary scope for the entity configuration (ALL, INBOUND)",
												Validators: []validator.String{
													stringvalidator.OneOf("ALL", "INBOUND"),
												},
											},
										},
									},
								},
								"application_event_based": schema.ListNestedBlock{
									Description: "The SLI entity of type availability to use for the SLI config",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"application_id": schema.StringAttribute{
												Required:    true,
												Description: "The application ID of the entity",
											},
											"boundary_scope": schema.StringAttribute{
												Required:    true,
												Description: "The boundary scope for the entity configuration (ALL, INBOUND)",
												Validators: []validator.String{
													stringvalidator.OneOf("ALL", "INBOUND"),
												},
											},
											"bad_event_filter_expression": schema.StringAttribute{
												Required:    true,
												Description: "The tag filter expression for bad events",
											},
											"good_event_filter_expression": schema.StringAttribute{
												Required:    true,
												Description: "The tag filter expression for good events",
											},
											"include_internal": schema.BoolAttribute{
												Optional:    true,
												Description: "Optional flag to indicate whether also internal calls are included",
											},
											"include_synthetic": schema.BoolAttribute{
												Optional:    true,
												Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not",
											},
											"endpoint_id": schema.StringAttribute{
												Optional:    true,
												Description: "Specifies the ID of the Endpoint to be monitored by the availability-based application SLO",
											},
											"service_id": schema.StringAttribute{
												Optional:    true,
												Description: "Identifies the service to be monitored by the availability-based application SLO",
											},
										},
									},
								},
								"website_event_based": schema.ListNestedBlock{
									Description: "The SLI entity of type websiteEventBased to use for the SLI config",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"website_id": schema.StringAttribute{
												Required:    true,
												Description: "The website ID of the entity",
											},
											"bad_event_filter_expression": schema.StringAttribute{
												Required:    true,
												Description: "The tag filter expression for bad events",
											},
											"good_event_filter_expression": schema.StringAttribute{
												Required:    true,
												Description: "The tag filter expression for good events",
											},
											"beacon_type": schema.StringAttribute{
												Required:    true,
												Description: "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)",
												Validators: []validator.String{
													stringvalidator.OneOf("pageLoad", "resourceLoad", "httpRequest", "error", "custom", "pageChange"),
												},
											},
										},
									},
								},
								"website_time_based": schema.ListNestedBlock{
									Description: "The SLI entity of type websiteTimeBased to use for the SLI config",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"website_id": schema.StringAttribute{
												Required:    true,
												Description: "The website ID of the entity",
											},
											"filter_expression": schema.StringAttribute{
												Optional:    true,
												Description: "The tag filter expression",
											},
											"beacon_type": schema.StringAttribute{
												Required:    true,
												Description: "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)",
												Validators: []validator.String{
													stringvalidator.OneOf("pageLoad", "resourceLoad", "httpRequest", "error", "custom", "pageChange"),
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
			SchemaVersion: 1,
			CreateOnly:    true,
		},
	}
}

func (r *sliConfigResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *sliConfigResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SliConfig] {
	return api.SliConfigs()
}

func (r *sliConfigResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *sliConfigResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, sliConfig *restapi.SliConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	model := SliConfigModel{
		ID:                         types.StringValue(sliConfig.ID),
		Name:                       types.StringValue(sliConfig.Name),
		InitialEvaluationTimestamp: types.Int64Value(int64(sliConfig.InitialEvaluationTimestamp)),
	}

	// Map metric configuration if present
	if sliConfig.MetricConfiguration != nil {
		metricConfigModel := MetricConfigurationModel{
			MetricName:  types.StringValue(sliConfig.MetricConfiguration.Name),
			Aggregation: types.StringValue(sliConfig.MetricConfiguration.Aggregation),
			Threshold:   types.Float64Value(sliConfig.MetricConfiguration.Threshold),
		}

		metricConfigObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"metric_name": types.StringType,
			"aggregation": types.StringType,
			"threshold":   types.Float64Type,
		}, metricConfigModel)
		if diags.HasError() {
			return diags
		}

		model.MetricConfiguration = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
				"threshold":   types.Float64Type,
			},
		}, []attr.Value{metricConfigObj})
	} else {
		model.MetricConfiguration = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metric_name": types.StringType,
				"aggregation": types.StringType,
				"threshold":   types.Float64Type,
			},
		})
	}

	// Map SLI entity
	sliEntityModel := SliEntityModel{}
	var entityDiags diag.Diagnostics

	sliEntityModel.ApplicationTimeBased = types.ListNull(applicationTimeBasedObjectType)
	sliEntityModel.ApplicationEventBased = types.ListNull(applicationEventBasedObjectType)
	sliEntityModel.WebsiteEventBased = types.ListNull(websiteEventBasedObjectType)
	sliEntityModel.WebsiteTimeBased = types.ListNull(websiteTimeBasedObjectType)

	switch sliConfig.SliEntity.Type {
	case "application":
		sliEntityModel.ApplicationTimeBased, entityDiags = r.mapApplicationTimeBasedToState(ctx, sliConfig.SliEntity)
	case "availability":
		sliEntityModel.ApplicationEventBased, entityDiags = r.mapApplicationEventBasedToState(ctx, sliConfig.SliEntity)
	case "websiteEventBased":
		sliEntityModel.WebsiteEventBased, entityDiags = r.mapWebsiteEventBasedToState(ctx, sliConfig.SliEntity)
	case "websiteTimeBased":
		sliEntityModel.WebsiteTimeBased, entityDiags = r.mapWebsiteTimeBasedToState(ctx, sliConfig.SliEntity)
	default:
		diags.AddError(
			"Unsupported SLI entity type",
			fmt.Sprintf("Unsupported SLI entity type: %s", sliConfig.SliEntity.Type),
		)
		return diags
	}

	if entityDiags.HasError() {
		diags.Append(entityDiags...)
		return diags
	}

	// Create SLI entity object
	sliEntityObj, entityObjDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"application_time_based":  types.ListType{ElemType: applicationTimeBasedObjectType},
		"application_event_based": types.ListType{ElemType: applicationEventBasedObjectType},
		"website_event_based":     types.ListType{ElemType: websiteEventBasedObjectType},
		"website_time_based":      types.ListType{ElemType: websiteTimeBasedObjectType},
	}, sliEntityModel)
	if entityObjDiags.HasError() {
		diags.Append(entityObjDiags...)
		return diags
	}

	model.SliEntity = types.ListValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"application_time_based":  types.ListType{ElemType: applicationTimeBasedObjectType},
			"application_event_based": types.ListType{ElemType: applicationEventBasedObjectType},
			"website_event_based":     types.ListType{ElemType: websiteEventBasedObjectType},
			"website_time_based":      types.ListType{ElemType: websiteTimeBasedObjectType},
		},
	}, []attr.Value{sliEntityObj})

	// Set the state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *sliConfigResourceFramework) mapApplicationTimeBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	appTimeBasedModel := ApplicationTimeBasedModel{
		ApplicationID: setStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: setStringPointerToState(sliEntity.BoundaryScope),
	}

	appTimeBasedModel.ServiceID = setStringPointerToState(sliEntity.ServiceID)

	appTimeBasedModel.EndpointID = setStringPointerToState(sliEntity.EndpointID)

	appTimeBasedObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"application_id": types.StringType,
		"service_id":     types.StringType,
		"endpoint_id":    types.StringType,
		"boundary_scope": types.StringType,
	}, appTimeBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValueMust(applicationTimeBasedObjectType, []attr.Value{appTimeBasedObj}), diags
}

func (r *sliConfigResourceFramework) mapApplicationEventBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	appEventBasedModel := ApplicationEventBasedModel{
		ApplicationID: setStringPointerToState(sliEntity.ApplicationID),
		BoundaryScope: setStringPointerToState(sliEntity.BoundaryScope),
		EndpointID:    setStringPointerToState(sliEntity.EndpointID),
		ServiceID:     setStringPointerToState(sliEntity.ServiceID),
	}

	// Map good event filter expression
	if sliEntity.GoodEventFilterExpression != nil {
		goodEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.GoodEventFilterExpression)
		if err != nil {
			diags.AddError(
				"Error mapping good event filter expression",
				fmt.Sprintf("Failed to map good event filter expression: %s", err),
			)
			return types.ListNull(types.ObjectType{}), diags
		}
		appEventBasedModel.GoodEventFilterExpression = setStringPointerToState(goodEventFilterStr)
	} else {
		appEventBasedModel.GoodEventFilterExpression = types.StringNull()
	}

	// Map bad event filter expression
	if sliEntity.BadEventFilterExpression != nil {
		badEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.BadEventFilterExpression)
		if err != nil {
			diags.AddError(
				"Error mapping bad event filter expression",
				fmt.Sprintf("Failed to map bad event filter expression: %s", err),
			)
			return types.ListNull(types.ObjectType{}), diags
		}
		appEventBasedModel.BadEventFilterExpression = setStringPointerToState(badEventFilterStr)
	} else {
		appEventBasedModel.BadEventFilterExpression = types.StringNull()
	}

	// Map include internal and synthetic flags
	if sliEntity.IncludeInternal != nil {
		appEventBasedModel.IncludeInternal = types.BoolValue(*sliEntity.IncludeInternal)
	} else {
		appEventBasedModel.IncludeInternal = types.BoolValue(false)
	}

	if sliEntity.IncludeSynthetic != nil {
		appEventBasedModel.IncludeSynthetic = types.BoolValue(*sliEntity.IncludeSynthetic)
	} else {
		appEventBasedModel.IncludeSynthetic = types.BoolValue(false)
	}

	appEventBasedObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"application_id":               types.StringType,
		"boundary_scope":               types.StringType,
		"bad_event_filter_expression":  types.StringType,
		"good_event_filter_expression": types.StringType,
		"include_internal":             types.BoolType,
		"include_synthetic":            types.BoolType,
		"service_id":                   types.StringType,
		"endpoint_id":                  types.StringType,
	}, appEventBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValueMust(applicationEventBasedObjectType, []attr.Value{appEventBasedObj}), diags
}

func (r *sliConfigResourceFramework) mapWebsiteEventBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	websiteEventBasedModel := WebsiteEventBasedModel{
		WebsiteID:  setStringPointerToState(sliEntity.WebsiteId),
		BeaconType: setStringPointerToState(sliEntity.BeaconType),
	}

	// Map good event filter expression
	if sliEntity.GoodEventFilterExpression != nil {
		goodEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.GoodEventFilterExpression)
		if err != nil {
			diags.AddError(
				"Error mapping good event filter expression",
				fmt.Sprintf("Failed to map good event filter expression: %s", err),
			)
			return types.ListNull(types.ObjectType{}), diags
		}
		websiteEventBasedModel.GoodEventFilterExpression = setStringPointerToState(goodEventFilterStr)
	} else {
		websiteEventBasedModel.GoodEventFilterExpression = types.StringNull()
	}

	// Map bad event filter expression
	if sliEntity.BadEventFilterExpression != nil {
		badEventFilterStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.BadEventFilterExpression)
		if err != nil {
			diags.AddError(
				"Error mapping bad event filter expression",
				fmt.Sprintf("Failed to map bad event filter expression: %s", err),
			)
			return types.ListNull(types.ObjectType{}), diags
		}
		websiteEventBasedModel.BadEventFilterExpression = setStringPointerToState(badEventFilterStr)
	} else {
		websiteEventBasedModel.BadEventFilterExpression = types.StringNull()
	}

	websiteEventBasedObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"website_id":                   types.StringType,
		"bad_event_filter_expression":  types.StringType,
		"good_event_filter_expression": types.StringType,
		"beacon_type":                  types.StringType,
	}, websiteEventBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValueMust(websiteEventBasedObjectType, []attr.Value{websiteEventBasedObj}), diags
}

func (r *sliConfigResourceFramework) mapWebsiteTimeBasedToState(ctx context.Context, sliEntity restapi.SliEntity) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	websiteTimeBasedModel := WebsiteTimeBasedModel{
		WebsiteID:  setStringPointerToState(sliEntity.WebsiteId),
		BeaconType: setStringPointerToState(sliEntity.BeaconType),
	}

	// Map filter expression
	if sliEntity.FilterExpression != nil {
		filterExprStr, err := tagfilter.MapTagFilterToNormalizedString(sliEntity.FilterExpression)
		if err != nil {
			diags.AddError(
				"Error mapping filter expression",
				fmt.Sprintf("Failed to map filter expression: %s", err),
			)
			return types.ListNull(types.ObjectType{}), diags
		}
		websiteTimeBasedModel.FilterExpression = setStringPointerToState(filterExprStr)
	} else {
		websiteTimeBasedModel.FilterExpression = types.StringNull()
	}

	websiteTimeBasedObj, objDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"website_id":        types.StringType,
		"filter_expression": types.StringType,
		"beacon_type":       types.StringType,
	}, websiteTimeBasedModel)
	if objDiags.HasError() {
		diags.Append(objDiags...)
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValueMust(websiteTimeBasedObjectType, []attr.Value{websiteTimeBasedObj}), diags
}

func (r *sliConfigResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SliConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model SliConfigModel

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

	// Map name
	name := model.Name.ValueString()

	// Map initial evaluation timestamp
	initialEvaluationTimestamp := 0
	if !model.InitialEvaluationTimestamp.IsNull() {
		initialEvaluationTimestamp = int(model.InitialEvaluationTimestamp.ValueInt64())
	}

	// Map metric configuration
	var metricConfiguration *restapi.MetricConfiguration
	if !model.MetricConfiguration.IsNull() && !model.MetricConfiguration.IsUnknown() {
		var metricConfigModels []MetricConfigurationModel
		diags.Append(model.MetricConfiguration.ElementsAs(ctx, &metricConfigModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(metricConfigModels) > 0 {
			metricConfigModel := metricConfigModels[0]
			metricConfiguration = &restapi.MetricConfiguration{
				Name:        metricConfigModel.MetricName.ValueString(),
				Aggregation: metricConfigModel.Aggregation.ValueString(),
				Threshold:   metricConfigModel.Threshold.ValueFloat64(),
			}
		}
	}

	// Map SLI entity
	var sliEntity restapi.SliEntity
	var entityErr error

	if !model.SliEntity.IsNull() && !model.SliEntity.IsUnknown() {
		var sliEntityModels []SliEntityModel
		diags.Append(model.SliEntity.ElementsAs(ctx, &sliEntityModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(sliEntityModels) > 0 {
			sliEntityModel := sliEntityModels[0]

			// Check which entity type is set
			if !sliEntityModel.ApplicationTimeBased.IsNull() && !sliEntityModel.ApplicationTimeBased.IsUnknown() {
				sliEntity, entityErr = r.mapApplicationTimeBasedFromState(ctx, sliEntityModel.ApplicationTimeBased)
			} else if !sliEntityModel.ApplicationEventBased.IsNull() && !sliEntityModel.ApplicationEventBased.IsUnknown() {
				sliEntity, entityErr = r.mapApplicationEventBasedFromState(ctx, sliEntityModel.ApplicationEventBased)
			} else if !sliEntityModel.WebsiteEventBased.IsNull() && !sliEntityModel.WebsiteEventBased.IsUnknown() {
				sliEntity, entityErr = r.mapWebsiteEventBasedFromState(ctx, sliEntityModel.WebsiteEventBased)
			} else if !sliEntityModel.WebsiteTimeBased.IsNull() && !sliEntityModel.WebsiteTimeBased.IsUnknown() {
				sliEntity, entityErr = r.mapWebsiteTimeBasedFromState(ctx, sliEntityModel.WebsiteTimeBased)
			}
		}
	}

	if entityErr != nil {
		diags.AddError(
			"Error mapping SLI entity",
			fmt.Sprintf("Failed to map SLI entity: %s", entityErr),
		)
		return nil, diags
	}

	// Create SLI config
	return &restapi.SliConfig{
		ID:                         id,
		Name:                       name,
		InitialEvaluationTimestamp: initialEvaluationTimestamp,
		MetricConfiguration:        metricConfiguration,
		SliEntity:                  sliEntity,
	}, diags
}

func (r *sliConfigResourceFramework) mapApplicationTimeBasedFromState(ctx context.Context, appTimeBasedList types.List) (restapi.SliEntity, error) {
	if appTimeBasedList.IsNull() || appTimeBasedList.IsUnknown() {
		return restapi.SliEntity{}, errors.New("application time based entity is null or unknown")
	}

	var appTimeBasedModels []ApplicationTimeBasedModel
	diags := appTimeBasedList.ElementsAs(ctx, &appTimeBasedModels, false)
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse application time based entity: %v", diags)
	}

	if len(appTimeBasedModels) == 0 {
		return restapi.SliEntity{}, errors.New("application time based entity list is empty")
	}

	appTimeBasedModel := appTimeBasedModels[0]
	applicationID := appTimeBasedModel.ApplicationID.ValueString()
	boundaryScope := appTimeBasedModel.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          "application",
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if !appTimeBasedModel.ServiceID.IsNull() {
		serviceID := appTimeBasedModel.ServiceID.ValueString()
		entity.ServiceID = &serviceID
	}

	if !appTimeBasedModel.EndpointID.IsNull() {
		endpointID := appTimeBasedModel.EndpointID.ValueString()
		entity.EndpointID = &endpointID
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapTagFilterStringToAPIModel(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *sliConfigResourceFramework) mapApplicationEventBasedFromState(ctx context.Context, appEventBasedList types.List) (restapi.SliEntity, error) {
	if appEventBasedList.IsNull() || appEventBasedList.IsUnknown() {
		return restapi.SliEntity{}, errors.New("application event based entity is null or unknown")
	}

	var appEventBasedModels []ApplicationEventBasedModel
	diags := appEventBasedList.ElementsAs(ctx, &appEventBasedModels, false)
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse application event based entity: %v", diags)
	}

	if len(appEventBasedModels) == 0 {
		return restapi.SliEntity{}, errors.New("application event based entity list is empty")
	}

	appEventBasedModel := appEventBasedModels[0]
	applicationID := appEventBasedModel.ApplicationID.ValueString()
	boundaryScope := appEventBasedModel.BoundaryScope.ValueString()

	entity := restapi.SliEntity{
		Type:          "availability",
		ApplicationID: &applicationID,
		BoundaryScope: &boundaryScope,
	}

	if !appEventBasedModel.BadEventFilterExpression.IsNull() {
		badEventFilter, err := r.mapTagFilterStringToAPIModel(appEventBasedModel.BadEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse bad event filter expression: %v", err)
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !appEventBasedModel.GoodEventFilterExpression.IsNull() {
		goodEventFilter, err := r.mapTagFilterStringToAPIModel(appEventBasedModel.GoodEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse good event filter expression: %v", err)
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	if !appEventBasedModel.IncludeInternal.IsNull() {
		includeInternal := appEventBasedModel.IncludeInternal.ValueBool()
		entity.IncludeInternal = &includeInternal
	}

	if !appEventBasedModel.IncludeSynthetic.IsNull() {
		includeSynthetic := appEventBasedModel.IncludeSynthetic.ValueBool()
		entity.IncludeSynthetic = &includeSynthetic
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapWebsiteEventBasedFromState(ctx context.Context, websiteEventBasedList types.List) (restapi.SliEntity, error) {
	if websiteEventBasedList.IsNull() || websiteEventBasedList.IsUnknown() {
		return restapi.SliEntity{}, errors.New("website event based entity is null or unknown")
	}

	var websiteEventBasedModels []WebsiteEventBasedModel
	diags := websiteEventBasedList.ElementsAs(ctx, &websiteEventBasedModels, false)
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse website event based entity: %v", diags)
	}

	if len(websiteEventBasedModels) == 0 {
		return restapi.SliEntity{}, errors.New("website event based entity list is empty")
	}

	websiteEventBasedModel := websiteEventBasedModels[0]
	websiteID := websiteEventBasedModel.WebsiteID.ValueString()
	beaconType := websiteEventBasedModel.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       "websiteEventBased",
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !websiteEventBasedModel.BadEventFilterExpression.IsNull() {
		badEventFilter, err := r.mapTagFilterStringToAPIModel(websiteEventBasedModel.BadEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse bad event filter expression: %v", err)
		}
		entity.BadEventFilterExpression = badEventFilter
	}

	if !websiteEventBasedModel.GoodEventFilterExpression.IsNull() {
		goodEventFilter, err := r.mapTagFilterStringToAPIModel(websiteEventBasedModel.GoodEventFilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse good event filter expression: %v", err)
		}
		entity.GoodEventFilterExpression = goodEventFilter
	}

	return entity, nil
}

func (r *sliConfigResourceFramework) mapWebsiteTimeBasedFromState(ctx context.Context, websiteTimeBasedList types.List) (restapi.SliEntity, error) {
	if websiteTimeBasedList.IsNull() || websiteTimeBasedList.IsUnknown() {
		return restapi.SliEntity{}, errors.New("website time based entity is null or unknown")
	}

	var websiteTimeBasedModels []WebsiteTimeBasedModel
	diags := websiteTimeBasedList.ElementsAs(ctx, &websiteTimeBasedModels, false)
	if diags.HasError() {
		return restapi.SliEntity{}, fmt.Errorf("failed to parse website time based entity: %v", diags)
	}

	if len(websiteTimeBasedModels) == 0 {
		return restapi.SliEntity{}, errors.New("website time based entity list is empty")
	}

	websiteTimeBasedModel := websiteTimeBasedModels[0]
	websiteID := websiteTimeBasedModel.WebsiteID.ValueString()
	beaconType := websiteTimeBasedModel.BeaconType.ValueString()

	entity := restapi.SliEntity{
		Type:       "websiteTimeBased",
		WebsiteId:  &websiteID,
		BeaconType: &beaconType,
	}

	if !websiteTimeBasedModel.FilterExpression.IsNull() {
		filterExpression, err := r.mapTagFilterStringToAPIModel(websiteTimeBasedModel.FilterExpression.ValueString())
		if err != nil {
			return restapi.SliEntity{}, fmt.Errorf("failed to parse filter expression: %v", err)
		}
		entity.FilterExpression = filterExpression
	}

	return entity, nil
}
