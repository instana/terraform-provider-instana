package sloconfig

// ResourceInstanaSloConfigFramework the name of the terraform-provider-instana resource to manage SLO configurations
const ResourceInstanaSloConfigFramework = "slo_config"

// SloConfigFieldRbacTags is the field name for RBAC tags
const SloConfigFieldRbacTags = "rbac_tags"

// Resource description constants
const (
	// SloConfigDescResource is the description for the SLO config resource
	SloConfigDescResource = "This resource manages SLO Configurations in Instana."
	// SloConfigDescID is the description for the ID field
	SloConfigDescID = "The ID of the SLO configuration"
	// SloConfigDescName is the description for the name field
	SloConfigDescName = "The name of the SLO configuration"
	// SloConfigDescTarget is the description for the target field
	SloConfigDescTarget = "The target of the SLO configuration"
	// SloConfigDescTags is the description for the tags field
	SloConfigDescTags = "The tags of the SLO configuration"
	// SloConfigDescRbacTags is the description for the rbac_tags field
	SloConfigDescRbacTags = "RBAC tags for the SLO configuration"
	// SloConfigDescRbacTagDisplayName is the description for the display_name field in rbac_tags
	SloConfigDescRbacTagDisplayName = "Display name of the RBAC tag"
	// SloConfigDescRbacTagID is the description for the id field in rbac_tags
	SloConfigDescRbacTagID = "ID of the RBAC tag"
	// SloConfigDescEntity is the description for the entity block
	SloConfigDescEntity = "The entity to use for the SLO configuration"
	// SloConfigDescApplicationEntity is the description for the application entity block
	SloConfigDescApplicationEntity = "Application entity of SLO"
	// SloConfigDescApplicationID is the description for the application_id field
	SloConfigDescApplicationID = "The application ID of the entity"
	// SloConfigDescBoundaryScope is the description for the boundary_scope field
	SloConfigDescBoundaryScope = "The boundary scope for the entity configuration (ALL, INBOUND)"
	// SloConfigDescEntityFilter is the description for the filter_expression field
	SloConfigDescEntityFilter = "Entity filter"
	// SloConfigDescIncludeInternal is the description for the include_internal field
	SloConfigDescIncludeInternal = "Optional flag to indicate whether also internal calls are included"
	// SloConfigDescIncludeSynthetic is the description for the include_synthetic field
	SloConfigDescIncludeSynthetic = "Optional flag to indicate whether also synthetic calls are included in the scope or not"
	// SloConfigDescServiceID is the description for the service_id field
	SloConfigDescServiceID = "The service ID of the entity"
	// SloConfigDescEndpointID is the description for the endpoint_id field
	SloConfigDescEndpointID = "The endpoint ID of the entity"
	// SloConfigDescWebsiteEntity is the description for the website entity block
	SloConfigDescWebsiteEntity = "Website entity of SLO"
	// SloConfigDescWebsiteID is the description for the website_id field
	SloConfigDescWebsiteID = "The website ID of the entity"
	// SloConfigDescBeaconType is the description for the beacon_type field
	SloConfigDescBeaconType = "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)"
	// SloConfigDescSyntheticEntity is the description for the synthetic entity block
	SloConfigDescSyntheticEntity = "Synthetic entity of SLO"
	// SloConfigDescSyntheticTestIDs is the description for the synthetic_test_ids field
	SloConfigDescSyntheticTestIDs = "The synthetics ID of the entity"
	// SloConfigDescInfrastructureEntity is the description for the infrastructure entity block
	SloConfigDescInfrastructureEntity = "Infrastructure entity of SLO"
	// SloConfigDescInfraType is the description for the infra_type field
	SloConfigDescInfraType = "The infrastructure type (e.g., kubernetesCluster)"
	// SloConfigDescIndicator is the description for the indicator block
	SloConfigDescIndicator = "The indicator to use for the SLO configuration"
	// SloConfigDescTimeBasedLatency is the description for the time_based_latency indicator
	SloConfigDescTimeBasedLatency = "Time-based latency indicator"
	// SloConfigDescEventBasedLatency is the description for the event_based_latency indicator
	SloConfigDescEventBasedLatency = "Event-based latency indicator"
	// SloConfigDescTimeBasedAvailability is the description for the time_based_availability indicator
	SloConfigDescTimeBasedAvailability = "Time-based availability indicator"
	// SloConfigDescEventBasedAvailability is the description for the event_based_availability indicator
	SloConfigDescEventBasedAvailability = "Event-based availability indicator"
	// SloConfigDescTraffic is the description for the traffic indicator
	SloConfigDescTraffic = "Traffic indicator"
	// SloConfigDescCustom is the description for the custom indicator
	SloConfigDescCustom = "Custom indicator"
	// SloConfigDescSaturation is the description for the saturation indicator
	SloConfigDescSaturation = "Saturation indicator"
	// SloConfigDescMetricName is the description for the metric_name field
	SloConfigDescMetricName = "The metric name for saturation indicator"
	// SloConfigDescThreshold is the description for the threshold field
	SloConfigDescThreshold = "The threshold for the metric configuration"
	// SloConfigDescAggregation is the description for the aggregation field
	SloConfigDescAggregation = "The aggregation type for the metric configuration"
	// SloConfigDescTrafficType is the description for the traffic_type field
	SloConfigDescTrafficType = "The traffic type for the indicator"
	// SloConfigDescOperator is the description for the operator field
	SloConfigDescOperator = "The aggregation type for the metric configuration"
	// SloConfigDescGoodEventFilterExpression is the description for the good_event_filter_expression field
	SloConfigDescGoodEventFilterExpression = "Good event filter expression"
	// SloConfigDescBadEventFilterExpression is the description for the bad_event_filter_expression field
	SloConfigDescBadEventFilterExpression = "Bad event filter expression"
	// SloConfigDescTimeWindow is the description for the time_window block
	SloConfigDescTimeWindow = "The time window to use for the SLO configuration"
	// SloConfigDescRollingTimeWindow is the description for the rolling time window
	SloConfigDescRollingTimeWindow = "Rolling time window"
	// SloConfigDescFixedTimeWindow is the description for the fixed time window
	SloConfigDescFixedTimeWindow = "Fixed time window"
	// SloConfigDescDuration is the description for the duration field
	SloConfigDescDuration = "The duration of the time window"
	// SloConfigDescDurationUnit is the description for the duration_unit field
	SloConfigDescDurationUnit = "The duration unit of the time window (day, week)"
	// SloConfigDescTimezone is the description for the timezone field
	SloConfigDescTimezone = "The timezone for the SLO configuration"
	// SloConfigDescStartTimestamp is the description for the start_timestamp field
	SloConfigDescStartTimestamp = "Time window start time"

	// Error message constants

	// SloConfigErrMappingState is the error title for mapping state to data object
	SloConfigErrMappingState = "Error mapping state to data object"
	// SloConfigErrBothPlanStateNil is the error message when both plan and state are nil
	SloConfigErrBothPlanStateNil = "Both plan and state are nil"
	// SloConfigErrApplicationIDRequired is the error message for missing application_id
	SloConfigErrApplicationIDRequired = "application_id and boundary_scope are required for application entity"
	// SloConfigErrWebsiteIDRequired is the error message for missing website_id
	SloConfigErrWebsiteIDRequired = "website_id and beacon_type are required for website entity"
	// SloConfigErrSyntheticTestIDsRequired is the error message for missing synthetic_test_ids
	SloConfigErrSyntheticTestIDsRequired = "synthetic_test_ids is required for synthetic entity"
	// SloConfigErrMissingEntity is the error title for missing entity configuration
	SloConfigErrMissingEntity = "Missing entity configuration"
	// SloConfigErrExactlyOneEntity is the error message for missing entity configuration
	SloConfigErrExactlyOneEntity = "Exactly one entity configuration is required"
	// SloConfigErrParsingFilterExpression is the error title for parsing filter expression
	SloConfigErrParsingFilterExpression = "Error parsing filter expression"
	// SloConfigErrParsingFilterExpressionMsg is the error message for parsing filter expression
	SloConfigErrParsingFilterExpressionMsg = "Could not parse filter expression: %s"
	// SloConfigErrTimeBasedLatencyRequired is the error message for missing time_based_latency fields
	SloConfigErrTimeBasedLatencyRequired = "threshold and  aggregation are required for time_based_latency indicator"
	// SloConfigErrEventBasedLatencyRequired is the error message for missing event_based_latency fields
	SloConfigErrEventBasedLatencyRequired = "threshold is required for event_based_latency indicator"
	// SloConfigErrTimeBasedAvailabilityRequired is the error message for missing time_based_availability fields
	SloConfigErrTimeBasedAvailabilityRequired = "threshold and  aggregation are required for time_based_availability indicator"
	// SloConfigErrTrafficRequired is the error message for missing traffic fields
	SloConfigErrTrafficRequired = "threshold is required for time_based_latency traffic indicator"
	// SloConfigErrCustomRequired is the error message for missing custom indicator fields
	SloConfigErrCustomRequired = "good_event_filter_expression is required for custom indicator"
	// SloConfigErrMissingIndicator is the error title for missing indicator configuration
	SloConfigErrMissingIndicator = "Missing indicator configuration"
	// SloConfigErrExactlyOneIndicator is the error message for missing indicator configuration
	SloConfigErrExactlyOneIndicator = "Exactly one indicator configuration is required"
	// SloConfigErrRollingTimeWindowRequired is the error message for missing rolling time window fields
	SloConfigErrRollingTimeWindowRequired = "duration and duration_unit are required for rolling time window"
	// SloConfigErrFixedTimeWindowRequired is the error message for missing fixed time window fields
	SloConfigErrFixedTimeWindowRequired = "duration,duration_unit,start_timestamp are required for fixed time window"
	// SloConfigErrMissingTimeWindow is the error title for missing time window configuration
	SloConfigErrMissingTimeWindow = "Missing time window configuration"
	// SloConfigErrExactlyOneTimeWindow is the error message for missing time window configuration
	SloConfigErrExactlyOneTimeWindow = "Exactly one time window configuration is required"
	// SloConfigErrMappingEntityToState is the error title for mapping entity to state
	SloConfigErrMappingEntityToState = "Error mapping entity to state"
	// SloConfigErrUnsupportedEntityType is the error message for unsupported entity type
	SloConfigErrUnsupportedEntityType = "Unsupported entity type: %s"
	// SloConfigErrNormalizingFilterExpression is the error title for normalizing filter expression
	SloConfigErrNormalizingFilterExpression = "Error normalizing filter expression"
	// SloConfigErrNormalizingFilterExpressionMsg is the error message for normalizing filter expression
	SloConfigErrNormalizingFilterExpressionMsg = "Could not normalize filter expression: %s"
	// SloConfigErrNormalizingGoodEventFilter is the error title for normalizing good event filter expression
	SloConfigErrNormalizingGoodEventFilter = "Error normalizing goodEventFilterExpression"
	// SloConfigErrNormalizingGoodEventFilterMsg is the error message for normalizing good event filter expression
	SloConfigErrNormalizingGoodEventFilterMsg = "Could not normalize goodEventFilterExpression: %s"
	// SloConfigErrNormalizingBadEventFilter is the error title for normalizing bad event filter expression
	SloConfigErrNormalizingBadEventFilter = "Error normalizing badEventFilterExpression"
	// SloConfigErrNormalizingBadEventFilterMsg is the error message for normalizing bad event filter expression
	SloConfigErrNormalizingBadEventFilterMsg = "Could not normalize badEventFilterExpression: %s"
	// SloConfigErrMappingIndicatorToState is the error title for mapping indicator to state
	SloConfigErrMappingIndicatorToState = "Error mapping indicator to state"
	// SloConfigErrUnsupportedIndicatorType is the error message for unsupported indicator type
	SloConfigErrUnsupportedIndicatorType = "Unsupported indicator type: %s, blueprint: %s"
	// SloConfigErrMappingTimeWindowToState is the error title for mapping time window to state
	SloConfigErrMappingTimeWindowToState = "Error mapping time window to state"
	// SloConfigErrUnsupportedTimeWindowType is the error message for unsupported time window type
	SloConfigErrUnsupportedTimeWindowType = "Unsupported time window type: %s"
	// SloConfigErrInfraTypeRequired is the error title for missing infrastructure type
	SloConfigErrInfraTypeRequired = "Infrastructure type required"
	// SloConfigErrInfraTypeRequiredMsg is the error message for missing infrastructure type
	SloConfigErrInfraTypeRequiredMsg = "infra_type is required for infrastructure entity"
	// SloConfigErrSaturationRequired is the error title for missing saturation indicator fields
	SloConfigErrSaturationRequired = "Saturation indicator fields required"
	// SloConfigErrSaturationRequiredMsg is the error message for missing saturation indicator fields
	SloConfigErrSaturationRequiredMsg = "threshold and operator are required for saturation indicator"

	// ResourceInstanaSloConfig the name of the terraform-provider-instana resource to manage SLI configurations
	ResourceInstanaSloConfig       = "instana_slo_config"
	SloConfigFromTerraformIdPrefix = "SLOTF"

	//SloConfigField names for terraform
	SloConfigFieldName   = "name"
	SloConfigFieldTarget = "target"
	SloConfigFieldTags                      = "tags"
	SloConfigFieldLastUpdated               = "last_updated"
	SloConfigFieldCreatedDate               = "created_date"
	SloConfigFieldSloEntity                 = "entity"
	SloConfigFieldSloIndicator              = "indicator"
	SloConfigFieldSloTimeWindow             = "time_window"
	SloConfigFieldApplicationID             = "application_id"
	SloConfigFieldWebsiteID                 = "website_id"
	SloConfigFieldSyntheticTestIDs          = "synthetic_test_ids"
	SloConfigFieldInfraType                 = "infra_type"
	SloConfigFieldFilterExpression          = "filter_expression"
	SloConfigFieldServiceID                 = "service_id"
	SloConfigFieldEndpointID                = "endpoint_id"
	SloConfigFieldIncludeInternal           = "include_internal"
	SloConfigFieldIncludeSynthetic          = "include_synthetic"
	SloConfigFieldBeaconType                = "beacon_type"
	SloConfigFieldBoundaryScope             = "boundary_scope"
	SloConfigFieldThreshold                 = "threshold"
	SloConfigFieldAggregation               = "aggregation"
	SloConfigFieldBadEventFilterExpression  = "bad_event_filter_expression"
	SloConfigFieldGoodEventFilterExpression = "good_event_filter_expression"
	SloConfigFieldTrafficType               = "traffic_type"
	SloConfigFieldMetricName                = "metric_name"
	SloConfigFieldDuration                  = "duration"
	SloConfigFieldDurationUnit              = "duration_unit"
	SloConfigFieldTimezone                  = "timezone"
	SloConfigFieldStartTimestamp            = "start_timestamp"

	// Slo entity types for terraform
	SloConfigApplicationEntity    = "application"
	SloConfigWebsiteEntity        = "website"
	SloConfigSyntheticEntity      = "synthetic"
	SloConfigInfrastructureEntity = "infrastructure"

	// Slo time windows types
	SloConfigRollingTimeWindow = "rolling"
	SloConfigFixedTimeWindow   = "fixed"

	// Slo indicator types for terraform
	SloConfigTimeBasedLatencyIndicator       = "time_based_latency"
	SloConfigEventBasedLatencyIndicator      = "event_based_latency"
	SloConfigTimeBasedAvailabilityIndicator  = "time_based_availability"
	SloConfigEventBasedAvailabilityIndicator = "event_based_availability"
	SloConfigTrafficIndicator                = "traffic"
	SloConfigCustomIndicator                 = "custom"

	// SloConfigFieldNames and values for API
	SloConfigAPIFieldThreshold       = "threshold"
	SloConfigAPIFieldAggregation     = "aggregation"
	SloConfigAPIFieldDuration        = "duration"
	SloConfigAPIFieldDurationUnit    = "durationUnit"
	SloConfigAPIFieldTimezone        = "timezone"
	SloConfigAPIFieldStartTimestamp  = "startTimestamp"
	SloConfigAPIFieldTrafficType     = "trafficType"
	SloConfigAPIFieldGoodEventFilter = "goodEventFilterExpression"
	SloConfigAPIFieldBadEventFilter  = "badEventFilterExpression"

	SloConfigAPIFieldFilter = "tagFilterExpression"

	SloConfigAPIIndicatorBlueprintLatency      = "latency"
	SloConfigAPIIndicatorBlueprintAvailability = "availability"
	SloConfigAPIIndicatorBlueprintTraffic      = "traffic"
	SloConfigAPIIndicatorBlueprintCustom       = "custom"
	SloConfigAPIIndicatorBlueprintSaturation   = "saturation"

	SloConfigAPIFieldBlueprint = "blueprint"
	SloConfigAPIFieldType      = "type"

	SloConfigAPIIndicatorMeasurementTypeTimeBased  = "timeBased"
	SloConfigAPIIndicatorMeasurementTypeEventBased = "eventBased"
	SloConfigAPITrafficIndicatorTypeAll            = "all"
	SloConfigAPITrafficIndicatorTypeErroneous      = "erroneous"
	// Schema field identifier constants

	// SchemaFieldID represents the id field identifier
	SchemaFieldID = "id"
	// SchemaFieldDisplayName represents the display_name field identifier
	SchemaFieldDisplayName = "display_name"
	// SchemaFieldRolling represents the rolling field identifier
	SchemaFieldRolling = "rolling"
	// SchemaFieldFixed represents the fixed field identifier
	SchemaFieldFixed = "fixed"
	// SchemaFieldTimeBasedLatency represents the time_based_latency field identifier
	SchemaFieldTimeBasedLatency = "time_based_latency"
	// SchemaFieldEventBasedLatency represents the event_based_latency field identifier
	SchemaFieldEventBasedLatency = "event_based_latency"
	// SchemaFieldTimeBasedAvailability represents the time_based_availability field identifier
	SchemaFieldTimeBasedAvailability = "time_based_availability"
	// SchemaFieldEventBasedAvailability represents the event_based_availability field identifier
	SchemaFieldEventBasedAvailability = "event_based_availability"
	// SchemaFieldTraffic represents the traffic field identifier
	SchemaFieldTraffic = "traffic"
	// SchemaFieldCustom represents the custom field identifier
	SchemaFieldCustom = "custom"
	// SchemaFieldSaturation represents the saturation field identifier
	SchemaFieldSaturation = "saturation"
	// SchemaFieldOperator represents the operator field identifier
	SchemaFieldOperator = "operator"

	// Operator constants

	// OperatorGreaterThan represents the > operator
	OperatorGreaterThan = ">"
	// OperatorGreaterThanOrEqual represents the >= operator
	OperatorGreaterThanOrEqual = ">="
	// OperatorLessThan represents the < operator
	OperatorLessThan = "<"
	// OperatorLessThanOrEqual represents the <= operator
	OperatorLessThanOrEqual = "<="

	// Tag filter constants

	// TagFilterTypeExpression represents the EXPRESSION tag filter type
	TagFilterTypeExpression = "EXPRESSION"
	// LogicalOperatorAnd represents the AND logical operator
	LogicalOperatorAnd = "AND"

	// Default values

	// DefaultAggregation represents the default aggregation type
	DefaultAggregation = "MEAN"
	// EmptyString represents an empty string constant
	EmptyString = ""
)
