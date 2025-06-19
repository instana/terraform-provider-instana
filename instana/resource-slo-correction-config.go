package instana

import (
	// "fmt"
	// "strconv"
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const ResourceInstanaSloCorrectionConfig = "instana_slo_correction_config"

const (
	// Slo Correction Config Field names for Terraform
	SloCorrectionConfigFieldName                    = "name"
	SloCorrectionConfigFieldFullName                = "full_name"
	SloCorrectionConfigFieldDescription             = "description"
	SloCorrectionConfigFieldActive                  = "active"
	SloCorrectionConfigFieldScheduling              = "scheduling"
	SloCorrectionConfigFieldSloIds                  = "slo_ids"
	SloCorrectionConfigFieldTags                    = "tags"
	SloCorrectionConfigFieldSchedulingStartTime     = "start_time"
	SloCorrectionConfigFieldSchedulingDuration      = "duration"
	SloCorrectionConfigFieldSchedulingDurationUnit  = "duration_unit"
	SloCorrectionConfigFieldSchedulingRecurrentRule = "recurrent_rule"
)

var (
	// SloCorrectionConfigName defines the schema for the 'name' field of instana_slo_correction_config
	SloCorrectionConfigName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(0, 256),
		Description:  "The name of the SLO Correction config.",
	}

	// SloCorrectionConfigFullName defines the schema for the computed 'full_name' field
	SloCorrectionConfigFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the SLO Correction config, computed based on provider-level default_name_prefix and default_name_suffix.",
	}

	// SloCorrectionConfigDescription defines the schema for the 'description' field
	SloCorrectionConfigDescription = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The description of the SLO Correction config.",
	}

	// SloCorrectionConfigActive defines the schema for the 'active' field
	SloCorrectionConfigActive = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Indicates whether the Correction config is active.",
	}

	// SloCorrectionConfigScheduling defines the schema for the 'scheduling' field
	SloCorrectionConfigScheduling = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Scheduling configuration for the SLO Correction config.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SloCorrectionConfigFieldSchedulingStartTime: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The start time of the scheduling in Unix timestamp in milliseconds.",
				},
				SloCorrectionConfigFieldSchedulingDuration: {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The duration of the scheduling in the specified unit.",
				},
				SloCorrectionConfigFieldSchedulingDurationUnit: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The unit of the duration (e.g.,'MINUTE' 'HOUR', 'DAY').",
					ValidateFunc: validation.StringInSlice([]string{"MINUTE", "HOUR", "DAY"}, true),
				},
				SloCorrectionConfigFieldSchedulingRecurrentRule: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Recurrent rule for scheduling, if applicable.",
				},
			},
		},
	}

	// SloCorrectionConfigSloIds defines the schema for the 'slo_ids' field
	SloCorrectionConfigSloIds = &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Description: "A set of SLO IDs that this correction config applies to.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	// SloCorrectionConfigTags defines the schema for the 'tags' field
	SloCorrectionConfigTags = &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "A list of tags to be associated with the SLO Correction config.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
)

// NewSloCorrectionConfigResourceHandle creates a new ResourceHandle for SLO Correction Config
func NewSloCorrectionConfigResourceHandle() ResourceHandle[*restapi.SloCorrectionConfig] {
	resource := &sloCorrectionConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSloCorrectionConfig,
			Schema: map[string]*schema.Schema{
				SloCorrectionConfigFieldName:        SloCorrectionConfigName,
				SloCorrectionConfigFieldDescription: SloCorrectionConfigDescription,
				SloCorrectionConfigFieldActive:      SloCorrectionConfigActive,
				SloCorrectionConfigFieldScheduling:  SloCorrectionConfigScheduling,
				SloCorrectionConfigFieldSloIds:      SloCorrectionConfigSloIds,
				SloCorrectionConfigFieldTags:        SloCorrectionConfigTags,
			},
			SchemaVersion:    1,
			CreateOnly:       false,
			SkipIDGeneration: true,
		},
	}
	return resource
}

// Schema
func (r *sloCorrectionConfigResource) sloCorrectionConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SloCorrectionConfigFieldName:        SloCorrectionConfigName,
			SloCorrectionConfigFieldDescription: SloCorrectionConfigDescription,
			SloCorrectionConfigFieldActive:      SloCorrectionConfigActive,
			SloCorrectionConfigFieldScheduling:  SloCorrectionConfigScheduling,
			SloCorrectionConfigFieldSloIds:      SloCorrectionConfigSloIds,
			SloCorrectionConfigFieldTags:        SloCorrectionConfigTags,
		},
	}
}
