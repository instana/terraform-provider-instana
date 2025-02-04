package instana

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaSloConfig the name of the terraform-provider-instana resource to manage SLI configurations
const ResourceInstanaSloConfig = "instana_slo_config"
const SloConfigFromTerraformIdPrefix = "SLOTF"

const (
	//SloConfigField names for terraform
	SloConfigFieldName                      = "name"
	SloConfigFieldTarget                    = "target"
	SloConfigFieldTags                      = "tags"
	SloConfigFieldLastUpdated               = "last_updated"
	SloConfigFieldCreatedDate               = "created_date"
	SloConfigFieldSloEntity                 = "entity"
	SloConfigFieldSloIndicator              = "indicator"
	SloConfigFieldSloTimeWindow             = "time_window"
	SloConfigFieldApplicationID             = "application_id"
	SloConfigFieldWebsiteID                 = "website_id"
	SloConfigFieldSyntheticTestIDs          = "synthetic_test_ids"
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
	SloConfigFieldDuration                  = "duration"
	SloConfigFieldDurationUnit              = "duration_unit"
	SloConfigFieldStartTimestamp            = "start_timestamp"

	// Slo entity types for terraform
	SloConfigApplicationEntity = "application"
	SloConfigWebsiteEntity     = "website"
	SloConfigSyntheticEntity   = "synthetic"

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
)

const (
	// SloConfigFieldNames and values for API
	SloConfigAPIFieldThreshold       = "threshold"
	SloConfigAPIFieldAggregation     = "aggregation"
	SloConfigAPIFieldDuration        = "duration"
	SloConfigAPIFieldDurationUnit    = "durationUnit"
	SloConfigAPIFieldStartTimestamp  = "startTimestamp"
	SloConfigAPIFieldTrafficType     = "trafficType"
	SloConfigAPIFieldGoodEventFilter = "goodEventFilterExpression"
	SloConfigAPIFieldBadEventFilter  = "badEventFilterExpression"

	SloConfigAPIFieldFilter = "tagFilterExpression"

	SloConfigAPIIndicatorBlueprintLatency      = "latency"
	SloConfigAPIIndicatorBlueprintAvailability = "availability"
	SloConfigAPIIndicatorBlueprintTraffic      = "traffic"
	SloConfigAPIIndicatorBlueprintCustom       = "custom"

	SloConfigAPIFieldBlueprint = "blueprint"
	SloConfigAPIFieldType      = "type"

	SloConfigAPIIndicatorMeasurementTypeTimeBased  = "timeBased"
	SloConfigAPIIndicatorMeasurementTypeEventBased = "eventBased"
	SloConfigAPITrafficIndicatorTypeAll            = "all"
	SloConfigAPITrafficIndicatorTypeErroneous      = "erroneous"
)

var sloConfigSliEntityTypeKeys = []string{
	"entity.0." + SloConfigApplicationEntity,
	"entity.0." + SloConfigWebsiteEntity,
	"entity.0." + SloConfigSyntheticEntity,
}

var sloConfigTimeWindowTypeKeys = []string{
	"time_window.0." + SloConfigRollingTimeWindow,
	"time_window.0." + SloConfigFixedTimeWindow,
}

var sloConfigIndicatorTypeKeys = []string{
	"indicator.0." + SloConfigTimeBasedLatencyIndicator,
	"indicator.0." + SloConfigEventBasedLatencyIndicator,
	"indicator.0." + SloConfigTimeBasedAvailabilityIndicator,
	"indicator.0." + SloConfigEventBasedAvailabilityIndicator,
	"indicator.0." + SloConfigTrafficIndicator,
	"indicator.0." + SloConfigCustomIndicator,
}

var (
	//SloConfigName schema field definition of instana_slo_config field name
	SloConfigName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(0, 256),
		Description:  "The name of the SLI config",
	}

	//SloConfigFullName schema field definition of instana_slo_config field full_name
	SloConfigFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the SLI config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}

	SloConfigTarget = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
		//	Computed:    true,
		Description: "The full name of the SLI config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}

	SloConfigTags = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "The full name of the SLI config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
		MinItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	//SloConfigInitialEvaluationTimestamp schema field definition of instana_slo_config field initial_evaluation_timestamp
	SloConfigInitialEvaluationTimestamp = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Initial evaluation timestamp for the SLI config",
	}

	SloConfigLastUpdated = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Initial evaluation timestamp for the SLI config",
	}

	SloConfigCreatedDate = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Initial evaluation timestamp for the SLI config",
	}

	SloConfigSchemaAggregation = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTRIBUTION", "DISTINCT_COUNT", "SUM_POSITIVE", "PER_SECOND"}, true),
		Description:  "The aggregation type for the metric configuration (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND)",
	}

	SloConfigSchemaThreshold = &schema.Schema{
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "The threshold for the metric configuration",
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(float64)
			if v <= 0 {
				errs = append(errs, fmt.Errorf("metric threshold must be greater than 0"))
			}
			return
		},
	}

	SloConfigSchemaSliEntityType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"application", "website", "synthetic"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	SloConfigSchemaSliEntityApplicationId = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The application ID of the entity",
	}

	SloConfigSchemaSliEntityWebsiteId = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The website ID of the entity",
	}

	SloConfigSchemaSliEntitySyntheticTestIds = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "The synthetics ID of the entity",
		MinItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	SloConfigSchemaSliEntityServiceId = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The service ID of the entity",
	}
	SloConfigSchemaSliEntityEndpointId = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The endpoint ID of the entity",
	}

	SloConfigSchemaSliEntityBoundaryScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"ALL", "INBOUND"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	SloConfigSchemaSliEntityBeaconType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"pageLoad", "resourceLoad", "httpRequest", "error", "custom", "pageChange"}, true),
		Description:  "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)",
	}
	SloConfigEntityFilter = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Entity filter",
	}

	SloConfigSchemaIncludeInternal = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also internal calls are included",
	}
	SloConfigSchemaIncludeSynthetic = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not",
	}

	SloConfigSchemaTimeWindowType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"fixed", "rolling"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	SloConfigSchemaDuration = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not",
	}

	SloConfigSchemaDurationUnit = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"day", "week"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	SloConfigSchemaStartTime = &schema.Schema{
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "Time window start time",
	}

	SloConfigSchemaBlueprint = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"latency", "availability", "traffic", "custom"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	SloConfigSchemaBlueprintType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"timeBased", "eventBased"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	SloConfigSchemaTrafficType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"all", "erroneous"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}

	//SloConfigSliEntity schema field definition of instana_slo_config field slo_entity
	SloConfigSliEntity = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "The entity to use for the SLI config.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SloConfigApplicationEntity: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Application entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SloConfigFieldApplicationID:    SloConfigSchemaSliEntityApplicationId,
							SloConfigFieldBoundaryScope:    SloConfigSchemaSliEntityBoundaryScope,
							SloConfigFieldFilterExpression: SloConfigEntityFilter,
							SloConfigFieldIncludeInternal:  SloConfigSchemaIncludeInternal,
							SloConfigFieldIncludeSynthetic: SloConfigSchemaIncludeSynthetic,
							SloConfigFieldServiceID:        SloConfigSchemaSliEntityServiceId,
							SloConfigFieldEndpointID:       SloConfigSchemaSliEntityEndpointId,
						},
					},
					ExactlyOneOf: sloConfigSliEntityTypeKeys,
				},
				SloConfigWebsiteEntity: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Website entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SloConfigFieldWebsiteID:        SloConfigSchemaSliEntityWebsiteId,
							SloConfigFieldFilterExpression: SloConfigEntityFilter,
							SloConfigFieldBeaconType:       SloConfigSchemaSliEntityBeaconType,
						},
					},
					ExactlyOneOf: sloConfigSliEntityTypeKeys,
				},
				SloConfigSyntheticEntity: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Synthetic entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SloConfigFieldSyntheticTestIDs: SloConfigSchemaSliEntitySyntheticTestIds,
							SloConfigFieldFilterExpression: SloConfigEntityFilter,
						},
					},
					ExactlyOneOf: sloConfigSliEntityTypeKeys,
				},
			},
		},
	}

	//SloConfigIndicator schema field definition of instana_slo_config field slo_entity
	SloConfigIndicator = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "The entity to use for the SLI config.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"time_based_latency": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Application entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"threshold":   SloConfigSchemaThreshold,
							"aggregation": SloConfigSchemaAggregation,
						},
					},
					ExactlyOneOf: sloConfigIndicatorTypeKeys,
				},
				"event_based_latency": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Website entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"threshold": SloConfigSchemaThreshold,
						},
					},
					ExactlyOneOf: sloConfigIndicatorTypeKeys,
				},
				"time_based_availability": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Application entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"threshold":   SloConfigSchemaThreshold,
							"aggregation": SloConfigSchemaAggregation,
						},
					},
					ExactlyOneOf: sloConfigIndicatorTypeKeys,
				},
				"event_based_availability": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Website entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{},
					},
					ExactlyOneOf: sloConfigIndicatorTypeKeys,
				},
				"traffic": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Application entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"traffic_type": SloConfigSchemaTrafficType,
							"threshold":    SloConfigSchemaThreshold,
							"aggregation":  SloConfigSchemaAggregation,
						},
					},
					ExactlyOneOf: sloConfigIndicatorTypeKeys,
				},
				"custom": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Application entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"good_event_filter_expression": RequiredTagFilterExpressionSchema,
							"bad_event_filter_expression":  OptionalTagFilterExpressionSchema,
						},
					},
					ExactlyOneOf: sloConfigIndicatorTypeKeys,
				},
			},
		},
	}

	//SloConfigTimeWindow schema field definition of instana_slo_config field slo_entity
	SloConfigTimeWindow = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "The entity to use for the SLI config.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rolling": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Application entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"duration":      SloConfigSchemaDuration,
							"duration_unit": SloConfigSchemaDurationUnit,
						},
					},
					ExactlyOneOf: sloConfigTimeWindowTypeKeys,
				},
				"fixed": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Website entity of SLO",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"duration":        SloConfigSchemaDuration,
							"duration_unit":   SloConfigSchemaDurationUnit,
							"start_timestamp": SloConfigSchemaStartTime,
						},
					},
					ExactlyOneOf: sloConfigTimeWindowTypeKeys,
				},
			},
		},
	}
)
