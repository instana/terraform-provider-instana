package sliconfig

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
)
