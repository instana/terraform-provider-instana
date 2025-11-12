package sliconfig

import (
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SliConfigModel represents the data model for SLI configuration
type SliConfigModel struct {
	ID                         types.String              `tfsdk:"id"`
	Name                       types.String              `tfsdk:"name"`
	InitialEvaluationTimestamp types.Int64               `tfsdk:"initial_evaluation_timestamp"`
	MetricConfiguration        *MetricConfigurationModel `tfsdk:"metric_configuration"`
	SliEntity                  *SliEntityModel           `tfsdk:"sli_entity"`
}

// MetricConfigurationModel represents the metric configuration for SLI
type MetricConfigurationModel struct {
	MetricName  types.String  `tfsdk:"metric_name"`
	Aggregation types.String  `tfsdk:"aggregation"`
	Threshold   types.Float64 `tfsdk:"threshold"`
}

// SliEntityModel represents the SLI entity configuration
type SliEntityModel struct {
	ApplicationTimeBased  *ApplicationTimeBasedModel  `tfsdk:"application_time_based"`
	ApplicationEventBased *ApplicationEventBasedModel `tfsdk:"application_event_based"`
	WebsiteEventBased     *WebsiteEventBasedModel     `tfsdk:"website_event_based"`
	WebsiteTimeBased      *WebsiteTimeBasedModel      `tfsdk:"website_time_based"`
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
	metaData resourcehandle.ResourceMetaDataFramework
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
