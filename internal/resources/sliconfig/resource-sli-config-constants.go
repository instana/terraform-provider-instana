package sliconfig

// ResourceInstanaSliConfig the name of the terraform-provider-instana resource to manage SLI configurations
const ResourceInstanaSliConfig = "sli_config"

const (
	//SliConfigFieldName constant value for the schema field name
	SliConfigFieldName = "name"
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

	// Resource description constants

	// SliConfigDescResource is the description for the SLI config resource
	SliConfigDescResource = "This resource manages SLI configurations in Instana."
	// SliConfigDescID is the description for the ID field
	SliConfigDescID = "The ID of the SLI configuration."
	// SliConfigDescName is the description for the name field
	SliConfigDescName = "The name of the SLI config"
	// SliConfigDescInitialEvaluationTimestamp is the description for the initial_evaluation_timestamp field
	SliConfigDescInitialEvaluationTimestamp = "Initial evaluation timestamp for the SLI config"
	// SliConfigDescMetricConfiguration is the description for the metric_configuration block
	SliConfigDescMetricConfiguration = "Metric configuration for the SLI config"
	// SliConfigDescMetricName is the description for the metric_name field
	SliConfigDescMetricName = "The metric name for the metric configuration"
	// SliConfigDescAggregation is the description for the aggregation field
	SliConfigDescAggregation = "The aggregation type for the metric configuration (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND)"
	// SliConfigDescThreshold is the description for the threshold field
	SliConfigDescThreshold = "The threshold for the metric configuration"
	// SliConfigDescSliEntity is the description for the sli_entity block
	SliConfigDescSliEntity = "The entity to use for the SLI config."
	// SliConfigDescApplicationTimeBased is the description for the application_time_based block
	SliConfigDescApplicationTimeBased = "The SLI entity of type application to use for the SLI config"
	// SliConfigDescApplicationID is the description for the application_id field
	SliConfigDescApplicationID = "The application ID of the entity"
	// SliConfigDescServiceID is the description for the service_id field
	SliConfigDescServiceID = "The service ID of the entity"
	// SliConfigDescEndpointID is the description for the endpoint_id field
	SliConfigDescEndpointID = "The endpoint ID of the entity"
	// SliConfigDescBoundaryScope is the description for the boundary_scope field
	SliConfigDescBoundaryScope = "The boundary scope for the entity configuration (ALL, INBOUND)"
	// SliConfigDescApplicationEventBased is the description for the application_event_based block
	SliConfigDescApplicationEventBased = "The SLI entity of type availability to use for the SLI config"
	// SliConfigDescBadEventFilterExpression is the description for the bad_event_filter_expression field
	SliConfigDescBadEventFilterExpression = "The tag filter expression for bad events"
	// SliConfigDescGoodEventFilterExpression is the description for the good_event_filter_expression field
	SliConfigDescGoodEventFilterExpression = "The tag filter expression for good events"
	// SliConfigDescIncludeInternal is the description for the include_internal field
	SliConfigDescIncludeInternal = "Optional flag to indicate whether also internal calls are included"
	// SliConfigDescIncludeSynthetic is the description for the include_synthetic field
	SliConfigDescIncludeSynthetic = "Optional flag to indicate whether also synthetic calls are included in the scope or not"
	// SliConfigDescEndpointIDAvailability is the description for the endpoint_id field in availability context
	SliConfigDescEndpointIDAvailability = "Specifies the ID of the Endpoint to be monitored by the availability-based application SLO"
	// SliConfigDescServiceIDAvailability is the description for the service_id field in availability context
	SliConfigDescServiceIDAvailability = "Identifies the service to be monitored by the availability-based application SLO"
	// SliConfigDescWebsiteEventBased is the description for the website_event_based block
	SliConfigDescWebsiteEventBased = "The SLI entity of type websiteEventBased to use for the SLI config"
	// SliConfigDescWebsiteID is the description for the website_id field
	SliConfigDescWebsiteID = "The website ID of the entity"
	// SliConfigDescBeaconType is the description for the beacon_type field
	SliConfigDescBeaconType = "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)"
	// SliConfigDescWebsiteTimeBased is the description for the website_time_based block
	SliConfigDescWebsiteTimeBased = "The SLI entity of type websiteTimeBased to use for the SLI config"
	// SliConfigDescFilterExpression is the description for the filter_expression field
	SliConfigDescFilterExpression = "The tag filter expression"

	// Error message constants

	// SliConfigErrUnsupportedEntityType is the error title for unsupported SLI entity type
	SliConfigErrUnsupportedEntityType = "Unsupported SLI entity type"
	// SliConfigErrUnsupportedEntityTypeMsg is the error message for unsupported SLI entity type
	SliConfigErrUnsupportedEntityTypeMsg = "Unsupported SLI entity type: %s"
	// SliConfigErrMappingGoodEventFilter is the error title for mapping good event filter expression
	SliConfigErrMappingGoodEventFilter = "Error mapping good event filter expression"
	// SliConfigErrMappingGoodEventFilterMsg is the error message for mapping good event filter expression
	SliConfigErrMappingGoodEventFilterMsg = "Failed to map good event filter expression: %s"
	// SliConfigErrMappingBadEventFilter is the error title for mapping bad event filter expression
	SliConfigErrMappingBadEventFilter = "Error mapping bad event filter expression"
	// SliConfigErrMappingBadEventFilterMsg is the error message for mapping bad event filter expression
	SliConfigErrMappingBadEventFilterMsg = "Failed to map bad event filter expression: %s"
	// SliConfigErrMappingFilterExpression is the error title for mapping filter expression
	SliConfigErrMappingFilterExpression = "Error mapping filter expression"
	// SliConfigErrMappingFilterExpressionMsg is the error message for mapping filter expression
	SliConfigErrMappingFilterExpressionMsg = "Failed to map filter expression: %s"
	// SliConfigErrMappingSliEntity is the error title for mapping SLI entity
	SliConfigErrMappingSliEntity = "Error mapping SLI entity"
	// SliConfigErrMappingSliEntityMsg is the error message for mapping SLI entity
	SliConfigErrMappingSliEntityMsg = "Failed to map SLI entity: %s"
	// SLI entity type constants

	// SliEntityTypeApplication represents the application time-based SLI entity type
	SliEntityTypeApplication = "application"
	// SliEntityTypeAvailability represents the application event-based SLI entity type
	SliEntityTypeAvailability = "availability"
	// SliEntityTypeWebsiteEventBased represents the website event-based SLI entity type
	SliEntityTypeWebsiteEventBased = "websiteEventBased"
	// SliEntityTypeWebsiteTimeBased represents the website time-based SLI entity type
	SliEntityTypeWebsiteTimeBased = "websiteTimeBased"

	// Aggregation type constants

	// AggregationTypeSum represents the SUM aggregation type
	AggregationTypeSum = "SUM"
	// AggregationTypeMean represents the MEAN aggregation type
	AggregationTypeMean = "MEAN"
	// AggregationTypeMax represents the MAX aggregation type
	AggregationTypeMax = "MAX"
	// AggregationTypeMin represents the MIN aggregation type
	AggregationTypeMin = "MIN"
	// AggregationTypeP25 represents the P25 percentile aggregation type
	AggregationTypeP25 = "P25"
	// AggregationTypeP50 represents the P50 percentile aggregation type
	AggregationTypeP50 = "P50"
	// AggregationTypeP75 represents the P75 percentile aggregation type
	AggregationTypeP75 = "P75"
	// AggregationTypeP90 represents the P90 percentile aggregation type
	AggregationTypeP90 = "P90"
	// AggregationTypeP95 represents the P95 percentile aggregation type
	AggregationTypeP95 = "P95"
	// AggregationTypeP98 represents the P98 percentile aggregation type
	AggregationTypeP98 = "P98"
	// AggregationTypeP99 represents the P99 percentile aggregation type
	AggregationTypeP99 = "P99"
	// AggregationTypeP99_9 represents the P99.9 percentile aggregation type
	AggregationTypeP99_9 = "P99_9"
	// AggregationTypeP99_99 represents the P99.99 percentile aggregation type
	AggregationTypeP99_99 = "P99_99"
	// AggregationTypeDistribution represents the DISTRIBUTION aggregation type
	AggregationTypeDistribution = "DISTRIBUTION"
	// AggregationTypeDistinctCount represents the DISTINCT_COUNT aggregation type
	AggregationTypeDistinctCount = "DISTINCT_COUNT"
	// AggregationTypeSumPositive represents the SUM_POSITIVE aggregation type
	AggregationTypeSumPositive = "SUM_POSITIVE"
	// AggregationTypePerSecond represents the PER_SECOND aggregation type
	AggregationTypePerSecond = "PER_SECOND"

	// Beacon type constants

	// BeaconTypePageLoad represents the pageLoad beacon type
	BeaconTypePageLoad = "pageLoad"
	// BeaconTypeResourceLoad represents the resourceLoad beacon type
	BeaconTypeResourceLoad = "resourceLoad"
	// BeaconTypeHttpRequest represents the httpRequest beacon type
	BeaconTypeHttpRequest = "httpRequest"
	// BeaconTypeError represents the error beacon type
	BeaconTypeError = "error"
	// BeaconTypeCustom represents the custom beacon type
	BeaconTypeCustom = "custom"
	// BeaconTypePageChange represents the pageChange beacon type
	BeaconTypePageChange = "pageChange"

	// Boundary scope constants

	// BoundaryScopeAll represents the ALL boundary scope
	BoundaryScopeAll = "ALL"
	// BoundaryScopeInbound represents the INBOUND boundary scope
	BoundaryScopeInbound = "INBOUND"

	// Schema field identifier constants

	// SchemaFieldID represents the id field identifier
	SchemaFieldID = "id"
	// SchemaFieldName represents the name field identifier
	SchemaFieldName = "name"
	// SchemaFieldInitialEvaluationTimestamp represents the initial_evaluation_timestamp field identifier
	SchemaFieldInitialEvaluationTimestamp = "initial_evaluation_timestamp"
	// SchemaFieldMetricConfiguration represents the metric_configuration field identifier
	SchemaFieldMetricConfiguration = "metric_configuration"
	// SchemaFieldMetricName represents the metric_name field identifier
	SchemaFieldMetricName = "metric_name"
	// SchemaFieldAggregation represents the aggregation field identifier
	SchemaFieldAggregation = "aggregation"
	// SchemaFieldThreshold represents the threshold field identifier
	SchemaFieldThreshold = "threshold"
	// SchemaFieldSliEntity represents the sli_entity field identifier
	SchemaFieldSliEntity = "sli_entity"
	// SchemaFieldApplicationTimeBased represents the application_time_based field identifier
	SchemaFieldApplicationTimeBased = "application_time_based"
	// SchemaFieldApplicationEventBased represents the application_event_based field identifier
	SchemaFieldApplicationEventBased = "application_event_based"
	// SchemaFieldWebsiteEventBased represents the website_event_based field identifier
	SchemaFieldWebsiteEventBased = "website_event_based"
	// SchemaFieldWebsiteTimeBased represents the website_time_based field identifier
	SchemaFieldWebsiteTimeBased = "website_time_based"
	// SchemaFieldApplicationID represents the application_id field identifier
	SchemaFieldApplicationID = "application_id"
	// SchemaFieldServiceID represents the service_id field identifier
	SchemaFieldServiceID = "service_id"
	// SchemaFieldEndpointID represents the endpoint_id field identifier
	SchemaFieldEndpointID = "endpoint_id"
	// SchemaFieldBoundaryScope represents the boundary_scope field identifier
	SchemaFieldBoundaryScope = "boundary_scope"
	// SchemaFieldBadEventFilterExpression represents the bad_event_filter_expression field identifier
	SchemaFieldBadEventFilterExpression = "bad_event_filter_expression"
	// SchemaFieldGoodEventFilterExpression represents the good_event_filter_expression field identifier
	SchemaFieldGoodEventFilterExpression = "good_event_filter_expression"
	// SchemaFieldIncludeInternal represents the include_internal field identifier
	SchemaFieldIncludeInternal = "include_internal"
	// SchemaFieldIncludeSynthetic represents the include_synthetic field identifier
	SchemaFieldIncludeSynthetic = "include_synthetic"
	// SchemaFieldWebsiteID represents the website_id field identifier
	SchemaFieldWebsiteID = "website_id"
	// SchemaFieldBeaconType represents the beacon_type field identifier
	SchemaFieldBeaconType = "beacon_type"
	// SchemaFieldFilterExpression represents the filter_expression field identifier
	SchemaFieldFilterExpression = "filter_expression"

	// Validation constants

	// NameMaxLength represents the maximum length for the name field
	NameMaxLength = 256
	// NameMinLength represents the minimum length for the name field
	NameMinLength = 0
	// ThresholdMinValue represents the minimum value for the threshold field
	ThresholdMinValue = 0.000001
	// DefaultInitialEvaluationTimestamp represents the default value for initial evaluation timestamp
	DefaultInitialEvaluationTimestamp = 0
	// DefaultIncludeInternalValue represents the default value for include_internal flag
	DefaultIncludeInternalValue = false
	// DefaultIncludeSyntheticValue represents the default value for include_synthetic flag
	DefaultIncludeSyntheticValue = false

	// Error message format constants

	// ErrMsgFailedToParseFilterExpression is the error message format for filter expression parsing failures
	ErrMsgFailedToParseFilterExpression = "failed to parse %s: %v"
	// ErrMsgBadEventFilterContext is the context identifier for bad event filter errors
	ErrMsgBadEventFilterContext = "bad event filter expression"
	// ErrMsgGoodEventFilterContext is the context identifier for good event filter errors
	ErrMsgGoodEventFilterContext = "good event filter expression"
	// ErrMsgFilterExpressionContext is the context identifier for filter expression errors
	ErrMsgFilterExpressionContext = "filter expression"
)
