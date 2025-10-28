package instana

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"os"
// 	// "log"

// 	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
// 	"github.com/gessnerfl/terraform-provider-instana/tfutils"

// 	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// // Debug utils here
// func obj2json(data any) string {
// 	b, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Printf("Error: %s", err)
// 	}
// 	return string(b)
// }

// func debug(statement any) {
// 	fmt.Fprintln(os.Stderr, statement)
// }

// // fiterExpFromAPI is a helper func to convert tag filter expression from terraform state format to API format
// func fiterExpFromAPIModel(filter interface{}) (*tagfilter.FilterExpression, bool, error) {
// 	var t *tagfilter.FilterExpression
// 	var err error
// 	if tagFilterString, ok := filter.(*restapi.TagFilter); ok {
// 		mapper := tagfilter.NewMapper()
// 		t, err = mapper.FromAPIModel(tagFilterString)

// 		if err != nil {
// 			debug(obj2json(err))
// 			return nil, true, err
// 		}
// 	}
// 	return t, false, nil
// }

// func (r *sloConfigResource) mapTagFilterStringToAPIModel(input string) (*restapi.TagFilter, error) {
// 	parser := tagfilter.NewParser()
// 	expr, err := parser.Parse(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mapper := tagfilter.NewMapper()
// 	return mapper.ToAPIModel(expr), nil
// }

// func (r *sloConfigResource) isSet(details interface{}) bool {
// 	list, ok := details.([]interface{})
// 	return ok && len(list) == 1
// }

// // NewSloConfigResourceHandle creates the resource handle for SLI configuration
// func NewSloConfigResourceHandle() ResourceHandle[*restapi.SloConfig] {
// 	cfgResource := &sloConfigResource{
// 		metaData: ResourceMetaData{
// 			ResourceName: ResourceInstanaSloConfig,
// 			Schema: map[string]*schema.Schema{
// 				SloConfigFieldName:          SloConfigName,
// 				SloConfigFieldTarget:        SloConfigTarget,
// 				SloConfigFieldTags:          SloConfigTags,
// 				SloConfigFieldLastUpdated:   SloConfigLastUpdated,
// 				SloConfigFieldCreatedDate:   SloConfigCreatedDate,
// 				SloConfigFieldSloEntity:     SloConfigSliEntity,
// 				SloConfigFieldSloIndicator:  SloConfigIndicator,
// 				SloConfigFieldSloTimeWindow: SloConfigTimeWindow,
// 			},
// 			SchemaVersion:    1,
// 			CreateOnly:       false,
// 			SkipIDGeneration: true,
// 		},
// 	}

// 	return cfgResource
// }

// type sloConfigResource struct {
// 	metaData ResourceMetaData
// }

// func (r *sloConfigResource) MetaData() *ResourceMetaData {
// 	resourceData := &r.metaData
// 	return resourceData
// }

// func (r *sloConfigResource) StateUpgraders() []schema.StateUpgrader {
// 	return []schema.StateUpgrader{
// 		{
// 			Type:    r.sloConfigSchemaV0().CoreConfigSchema().ImpliedType(),
// 			Upgrade: r.sloConfigStateUpgradeV0,
// 			Version: 0,
// 		},
// 	}
// }

// func (r *sloConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SloConfig] {
// 	x := api.SloConfigs()
// 	return x
// }

// func (r *sloConfigResource) SetComputedFields(_ *schema.ResourceData) error {
// 	return nil
// }

// func (r *sloConfigResource) UpdateState(d *schema.ResourceData, sloConfig *restapi.SloConfig) error {
// 	debug(">> UpdateState")
// 	SloEntity, err := r.mapSloEntityToState(sloConfig)
// 	if err != nil {
// 		return err
// 	}

// 	SloIndicator, err := r.mapSloIndicatorToState(sloConfig)
// 	if err != nil {
// 		return err
// 	}

// 	SloTimeWindow, err := r.mapSloTimeWindowToState(sloConfig)
// 	if err != nil {
// 		return err
// 	}

// 	tfData := map[string]interface{}{
// 		SloConfigFieldName:          sloConfig.Name,
// 		SloConfigFieldTarget:        sloConfig.Target,
// 		SloConfigFieldTags:          sloConfig.Tags,
// 		SloConfigFieldSloEntity:     []interface{}{SloEntity},
// 		SloConfigFieldSloIndicator:  []interface{}{SloIndicator},
// 		SloConfigFieldSloTimeWindow: []interface{}{SloTimeWindow},
// 		// SloConfigFieldLastUpdated:   sloConfig.LastUpdated,
// 		// SloConfigFieldCreatedDate:   sloConfig.CreatedDate,
// 	}

// 	d.SetId(sloConfig.ID)

// 	debug(">> UpdateState with: " + obj2json(tfData))

// 	return tfutils.UpdateState(d, tfData)
// }

// // tf state -> api
// func (r *sloConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SloConfig, error) {
// 	debug(">> MapStateToDataObject")
// 	debug(obj2json(d))

// 	debug(">> entity")
// 	sloEntitiesStateObject := d.Get(SloConfigFieldSloEntity).([]interface{})
// 	var sloEntity any
// 	var err error
// 	if len(sloEntitiesStateObject) == 1 {
// 		sloEntityObject := sloEntitiesStateObject[0].(map[string]interface{})
// 		sloEntity, err = r.mapSliEntityListFromState(sloEntityObject)
// 		debug(">> MapStateToDataObject")
// 		if err != nil {
// 			debug(obj2json(err))
// 			return nil, err
// 		}
// 	} else {
// 		return nil, errors.New("exactly one entity configuration is required")
// 	}

// 	debug(">> indicator")
// 	sloIndicatorStateObject := d.Get(SloConfigFieldSloIndicator).([]interface{})
// 	debug(obj2json(sloIndicatorStateObject))
// 	var sloIndicator any
// 	if len(sloIndicatorStateObject) == 1 {
// 		sloIndicatorObject := sloIndicatorStateObject[0].(map[string]interface{})
// 		sloIndicator, err = r.mapSliIndicatorListFromState(sloIndicatorObject)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		return nil, errors.New("exactly one indicator configuration is required")
// 	}

// 	debug(">> time window")
// 	sloTimeWindowStateObject := d.Get(SloConfigFieldSloTimeWindow).([]interface{})
// 	var sloTimeWindow any
// 	if len(sloTimeWindowStateObject) == 1 {
// 		sloTimeWindowObject := sloTimeWindowStateObject[0].(map[string]interface{})
// 		sloTimeWindow, err = r.mapSliTimeWindowListFromState(sloTimeWindowObject)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		return nil, errors.New("exactly one time window configuration is required")
// 	}

// 	sid := d.Id()
// 	if len(sid) == 0 {
// 		sid = SloConfigFromTerraformIdPrefix + RandomID()
// 	}

// 	payload := &restapi.SloConfig{
// 		ID:   sid,
// 		Name: d.Get(SloConfigFieldName).(string),
// 		// LastUpdated: d.Get(SloConfigFieldLastUpdated).(int),
// 		// CreatedDate: d.Get(SloConfigFieldCreatedDate).(int),
// 		Target:     d.Get(SloConfigFieldTarget).(float64),
// 		Tags:       d.Get(SloConfigFieldTags),
// 		Entity:     sloEntity,
// 		Indicator:  sloIndicator,
// 		TimeWindow: sloTimeWindow,
// 	}

// 	// debug utils
// 	// payloadJSON, err := json.MarshalIndent(payload, "", "  ")
//     // if err != nil {
//     //     log.Printf("Error marshalling payload to JSON: %v", err)
//     // } else {
//     //     log.Printf("Payload sent to API: %s", string(payloadJSON))
//     // }

// 	return payload, nil
// }

// func (r *sloConfigResource) sloConfigStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
// 	if _, ok := state[SloConfigFieldFullName]; ok {
// 		state[SloConfigFieldName] = state[SloConfigFieldFullName]
// 		delete(state, SloConfigFieldFullName)
// 	}
// 	return state, nil
// }

// // root schema
// func (r *sloConfigResource) sloConfigSchemaV0() *schema.Resource {
// 	return &schema.Resource{
// 		Schema: map[string]*schema.Schema{
// 			SloConfigFieldName    : SloConfigName,
// 			SloConfigFieldFullName: SloConfigFullName,
// 			SloConfigFieldTarget  : SloConfigTarget,
// 			SloConfigFieldTags    : SloConfigTags,
// 			  // SloConfigFieldLastUpdated:   SloConfigLastUpdated,
// 			  // SloConfigFieldCreatedDate:   SloConfigCreatedDate,
// 			SloConfigFieldSloEntity    : SloConfigSliEntity,
// 			SloConfigFieldSloIndicator : SloConfigIndicator,
// 			SloConfigFieldSloTimeWindow: SloConfigTimeWindow,
// 		},
// 	}
// }
