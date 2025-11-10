package instana

// import (
// 	"context"
// 	"strings"

// 	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
// )

// const ResourceInstanaSloCorrectionConfig = "instana_slo_correction_config"

// const (
// 	// Slo Correction Config Field names for Terraform
// 	SloCorrectionConfigFieldName                    = "name"
// 	SloCorrectionConfigFieldFullName                = "full_name"
// 	SloCorrectionConfigFieldDescription             = "description"
// 	SloCorrectionConfigFieldActive                  = "active"
// 	SloCorrectionConfigFieldScheduling              = "scheduling"
// 	SloCorrectionConfigFieldSloIds                  = "slo_ids"
// 	SloCorrectionConfigFieldTags                    = "tags"
// 	SloCorrectionConfigFieldSchedulingStartTime     = "start_time"
// 	SloCorrectionConfigFieldSchedulingDuration      = "duration"
// 	SloCorrectionConfigFieldSchedulingDurationUnit  = "duration_unit"
// 	SloCorrectionConfigFieldSchedulingRecurrentRule = "recurrent_rule"
// )

// var (
// 	// SloCorrectionConfigName defines the schema for the 'name' field of instana_slo_correction_config
// 	SloCorrectionConfigName = &schema.Schema{
// 		Type:         schema.TypeString,
// 		Required:     true,
// 		ValidateFunc: validation.StringLenBetween(0, 256),
// 		Description:  "The name of the SLO Correction config.",
// 	}

// 	// SloCorrectionConfigFullName defines the schema for the computed 'full_name' field
// 	SloCorrectionConfigFullName = &schema.Schema{
// 		Type:        schema.TypeString,
// 		Computed:    true,
// 		Description: "The full name of the SLO Correction config, computed based on provider-level default_name_prefix and default_name_suffix.",
// 	}

// 	// SloCorrectionConfigDescription defines the schema for the 'description' field
// 	SloCorrectionConfigDescription = &schema.Schema{
// 		Type:        schema.TypeString,
// 		Required:    true,
// 		Description: "The description of the SLO Correction config.",
// 	}

// 	// SloCorrectionConfigActive defines the schema for the 'active' field
// 	SloCorrectionConfigActive = &schema.Schema{
// 		Type:        schema.TypeBool,
// 		Required:    true,
// 		Description: "Indicates whether the Correction config is active.",
// 	}

// 	// SloCorrectionConfigScheduling defines the schema for the 'scheduling' field
// 	SloCorrectionConfigScheduling = &schema.Schema{
// 		Type:        schema.TypeList,
// 		Required:    true,
// 		Description: "Scheduling configuration for the SLO Correction config.",
// 		Elem: &schema.Resource{
// 			Schema: map[string]*schema.Schema{
// 				SloCorrectionConfigFieldSchedulingStartTime: {
// 					Type:        schema.TypeInt,
// 					Required:    true,
// 					Description: "The start time of the scheduling in Unix timestamp in milliseconds.",
// 				},
// 				SloCorrectionConfigFieldSchedulingDuration: {
// 					Type:        schema.TypeInt,
// 					Required:    true,
// 					Description: "The duration of the scheduling in the specified unit.",
// 				},
// 				SloCorrectionConfigFieldSchedulingDurationUnit: {
// 					Type:         schema.TypeString,
// 					Required:     true,
// 					Description:  "The unit of the duration (e.g.,'MINUTE' 'HOUR', 'DAY').",
// 					ValidateFunc: validation.StringInSlice([]string{"MINUTE", "HOUR", "DAY"}, true),
// 				},
// 				SloCorrectionConfigFieldSchedulingRecurrentRule: {
// 					Type:        schema.TypeString,
// 					Optional:    true,
// 					Description: "Recurrent rule for scheduling, if applicable.",
// 				},
// 			},
// 		},
// 	}

// 	// SloCorrectionConfigSloIds defines the schema for the 'slo_ids' field
// 	SloCorrectionConfigSloIds = &schema.Schema{
// 		Type:        schema.TypeSet,
// 		Required:    true,
// 		Description: "A set of SLO IDs that this correction config applies to.",
// 		Elem: &schema.Schema{
// 			Type: schema.TypeString,
// 		},
// 	}

// 	// SloCorrectionConfigTags defines the schema for the 'tags' field
// 	SloCorrectionConfigTags = &schema.Schema{
// 		Type:        schema.TypeSet,
// 		Optional:    true,
// 		Description: "A list of tags to be associated with the SLO Correction config.",
// 		Elem: &schema.Schema{
// 			Type: schema.TypeString,
// 		},
// 	}
// )

// // sloCorrectionConfigResource is a struct that implements ResourceHandle for SLO Correction Config
// type sloCorrectionConfigResource struct {
// 	metaData ResourceMetaData
// }

// // NewSloCorrectionConfigResourceHandle creates a new ResourceHandle for SLO Correction Config
// func NewSloCorrectionConfigResourceHandle() ResourceHandle[*restapi.SloCorrectionConfig] {
// 	resource := &sloCorrectionConfigResource{
// 		metaData: ResourceMetaData{
// 			ResourceName: ResourceInstanaSloCorrectionConfig,
// 			Schema: map[string]*schema.Schema{
// 				SloCorrectionConfigFieldName:        SloCorrectionConfigName,
// 				SloCorrectionConfigFieldDescription: SloCorrectionConfigDescription,
// 				SloCorrectionConfigFieldActive:      SloCorrectionConfigActive,
// 				SloCorrectionConfigFieldScheduling:  SloCorrectionConfigScheduling,
// 				SloCorrectionConfigFieldSloIds:      SloCorrectionConfigSloIds,
// 				SloCorrectionConfigFieldTags:        SloCorrectionConfigTags,
// 			},
// 			SchemaVersion:    1,
// 			CreateOnly:       false,
// 			SkipIDGeneration: true,
// 		},
// 	}
// 	return resource
// }

// func (r *sloCorrectionConfigResource) MetaData() *ResourceMetaData {
// 	return &r.metaData
// }

// func (r *sloCorrectionConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloCorrectionConfig] {
// 	return api.SloCorrectionConfig()
// }

// func (r *sloCorrectionConfigResource) SetComputedFields(_ *schema.ResourceData) error {
// 	return nil
// }

// func (r *sloCorrectionConfigResource) StateUpgraders() []schema.StateUpgrader {
// 	return []schema.StateUpgrader{
// 		{
// 			Type:    r.sloCorrectionConfigSchemaV0().CoreConfigSchema().ImpliedType(),
// 			Upgrade: r.sloCorrectionConfigStateUpgradeV0,
// 			Version: 0,
// 		},
// 	}
// }

// func (r *sloCorrectionConfigResource) sloCorrectionConfigSchemaV0() *schema.Resource {
// 	return &schema.Resource{
// 		Schema: map[string]*schema.Schema{
// 			SloCorrectionConfigFieldName:        SloCorrectionConfigName,
// 			SloCorrectionConfigFieldDescription: SloCorrectionConfigDescription,
// 			SloCorrectionConfigFieldActive:      SloCorrectionConfigActive,
// 			SloCorrectionConfigFieldScheduling:  SloCorrectionConfigScheduling,
// 			SloCorrectionConfigFieldSloIds:      SloCorrectionConfigSloIds,
// 			SloCorrectionConfigFieldTags:        SloCorrectionConfigTags,
// 		},
// 	}
// }

// func (r *sloCorrectionConfigResource) sloCorrectionConfigStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
// 	if _, ok := state[SloCorrectionConfigFieldFullName]; ok {
// 		state[SloCorrectionConfigFieldName] = state[SloCorrectionConfigFieldFullName]
// 		delete(state, SloCorrectionConfigFieldFullName)
// 	}
// 	return state, nil
// }

// // MapStateToDataObject maps the Terraform state to the SloCorrectionConfig data object.
// func (r *sloCorrectionConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SloCorrectionConfig, error) {
// 	//debug(">> MapStateToDataObject")
// 	//debug(obj2json(d))
// 	sid := d.Id()
// 	if len(sid) == 0 {
// 		sid = RandomID()
// 	}

// 	// Construct payload for SLO Correction Config
// 	schedulingData := d.Get(SloCorrectionConfigFieldScheduling).([]interface{})
// 	var scheduling restapi.Scheduling
// 	data := schedulingData[0].(map[string]interface{})
// 	startTimeRaw := data[SloCorrectionConfigFieldSchedulingStartTime]
// 	var startTime int64
// 	switch v := startTimeRaw.(type) {
// 	case int64:
// 		startTime = v
// 	case int:
// 		startTime = int64(v)
// 	}
// 	duration := data[SloCorrectionConfigFieldSchedulingDuration].(int)
// 	durationUnit := restapi.DurationUnit(strings.ToUpper(data[SloCorrectionConfigFieldSchedulingDurationUnit].(string)))
// 	recurrentRule := ""
// 	if v, ok := data[SloCorrectionConfigFieldSchedulingRecurrentRule]; ok && v != nil {
// 		recurrentRule = v.(string)
// 	}
// 	scheduling = restapi.Scheduling{
// 		StartTime:     startTime,
// 		Duration:      duration,
// 		DurationUnit:  durationUnit,
// 		RecurrentRule: recurrentRule,
// 	}

// 	var tags []string
// 	rawTags := d.Get(SloCorrectionConfigFieldTags)
// 	if rawTags != nil {
// 		tags = convertSetToStringSlice(rawTags.(*schema.Set))
// 	} else {
// 		tags = []string{}
// 	}

// 	payload := &restapi.SloCorrectionConfig{
// 		ID:          sid,
// 		Name:        d.Get(SloCorrectionConfigFieldName).(string),
// 		Description: d.Get(SloCorrectionConfigFieldDescription).(string),
// 		Active:      d.Get(SloCorrectionConfigFieldActive).(bool),
// 		Scheduling:  scheduling,
// 		SloIds:      convertSetToStringSlice(d.Get(SloCorrectionConfigFieldSloIds).(*schema.Set)),
// 		Tags:        tags,
// 	}

// 	return payload, nil
// }

// // UpdateState updates the Terraform state from the SloCorrectionConfig data object.
// func (r *sloCorrectionConfigResource) UpdateState(d *schema.ResourceData, obj *restapi.SloCorrectionConfig) error {
// 	d.SetId(obj.ID)
// 	if err := d.Set(SloCorrectionConfigFieldName, obj.Name); err != nil {
// 		return err
// 	}
// 	if err := d.Set(SloCorrectionConfigFieldDescription, obj.Description); err != nil {
// 		return err
// 	}
// 	if err := d.Set(SloCorrectionConfigFieldActive, obj.Active); err != nil {
// 		return err
// 	}
// 	scheduling := map[string]interface{}{
// 		SloCorrectionConfigFieldSchedulingStartTime:     obj.Scheduling.StartTime,
// 		SloCorrectionConfigFieldSchedulingDuration:      obj.Scheduling.Duration,
// 		SloCorrectionConfigFieldSchedulingDurationUnit:  strings.ToUpper(string(obj.Scheduling.DurationUnit)),
// 		SloCorrectionConfigFieldSchedulingRecurrentRule: obj.Scheduling.RecurrentRule,
// 	}
// 	if err := d.Set(SloCorrectionConfigFieldScheduling, []interface{}{scheduling}); err != nil {
// 		return err
// 	}
// 	if err := d.Set(SloCorrectionConfigFieldSloIds, obj.SloIds); err != nil {
// 		return err
// 	}
// 	if err := d.Set(SloCorrectionConfigFieldTags, obj.Tags); err != nil {
// 		return err
// 	}
// 	return nil
// }
