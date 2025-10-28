package instana

// // api -> tf state
// func (r *sloConfigResource) mapSloEntityToState(sloConfig *restapi.SloConfig) (map[string]interface{}, error) {

// 	sloEntity := sloConfig.Entity
// 	debug(">> mapSloEntityToState with: " + obj2json(sloEntity))

// 	if entity, ok := sloEntity.(map[string]interface{}); ok {
// 		debug(entity)

// 		if entity["type"] == SloConfigApplicationEntity {
// 			filterExp, validExp, err := fiterExpFromAPIModel(entity[SloConfigAPIFieldFilter])
// 			if validExp {
// 				return nil, err
// 			} else {
// 				entityAttrs := map[string]interface{}{
// 					SloConfigFieldApplicationID:    entity["applicationId"].(string),
// 					SloConfigFieldBoundaryScope:    entity["boundaryScope"],
// 					SloConfigFieldIncludeInternal:  entity["includeInternal"],
// 					SloConfigFieldIncludeSynthetic: entity["includeSynthetic"],
// 					SloConfigFieldFilterExpression: filterExp,
// 				}
// 				if s, ok := entity["serviceId"]; ok && entity["serviceId"] != nil {
// 					entityAttrs[SloConfigFieldServiceID] = s.(string)
// 				}
// 				if s, ok := entity["endpointId"]; ok && s != nil {
// 					entityAttrs[SloConfigFieldEndpointID] = s.(string)
// 				}

// 				e := map[string]interface{}{
// 					SloConfigApplicationEntity: []interface{}{
// 						entityAttrs,
// 					},
// 				}
// 				debug(">> entity: " + obj2json(e))
// 				return e, nil

// 			}

// 		}

// 		if entity["type"] == SloConfigWebsiteEntity {
// 			filterExp, validExp, err := fiterExpFromAPIModel(entity[SloConfigAPIFieldFilter])
// 			if validExp {
// 				return nil, err
// 			} else {
// 				entityAttrs := map[string]interface{}{
// 					SloConfigFieldWebsiteID:        entity["websiteId"].(string),
// 					SloConfigFieldBeaconType:       entity["beaconType"],
// 					SloConfigFieldFilterExpression: filterExp,
// 				}

// 				e := map[string]interface{}{
// 					SloConfigWebsiteEntity: []interface{}{
// 						entityAttrs,
// 					},
// 				}
// 				debug(">> entity: " + obj2json(e))
// 				return e, nil

// 			}

// 		}

// 		if entity["type"] == SloConfigSyntheticEntity {
// 			filterExp, validExp, err := fiterExpFromAPIModel(entity[SloConfigAPIFieldFilter])
// 			if validExp {
// 				return nil, err
// 			} else {
// 				entityAttrs := map[string]interface{}{
// 					SloConfigFieldSyntheticTestIDs: entity["syntheticTestIds"],
// 					SloConfigFieldFilterExpression: filterExp,
// 				}

// 				e := map[string]interface{}{
// 					SloConfigSyntheticEntity: []interface{}{
// 						entityAttrs,
// 					},
// 				}
// 				debug(">> entity: " + obj2json(e))
// 				return e, nil

// 			}

// 		}

// 		return nil, fmt.Errorf("unsupported sli entity type: %s", entity["type"])
// 	}
// 	return nil, fmt.Errorf("the \"type: %s\" attribute is missed from entity definition.", "type")
// }

// // tf state -> api
// func (r *sloConfigResource) mapSliEntityListFromState(stateObject map[string]interface{}) (any, error) {
// 	debug(">> mapSliEntityListFromState with: " + obj2json(stateObject))

// 	if len(stateObject) > 0 {
// 		if details, ok := stateObject[SloConfigApplicationEntity]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			var tagFilter *restapi.TagFilter
// 			var err error
// 			if tagFilterString, ok := data[SloConfigFieldFilterExpression]; ok {
// 				tagFilter, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
// 				if err != nil {
// 					debug(obj2json(err))
// 					// return nil, err
// 				}
// 			}

// 			return restapi.SloApplicationEntity{
// 				Type:             SloConfigApplicationEntity,
// 				ApplicationID:    GetPointerFromMap[string](data, SloConfigFieldApplicationID),
// 				ServiceID:        GetPointerFromMap[string](data, SloConfigFieldServiceID),
// 				EndpointID:       GetPointerFromMap[string](data, SloConfigFieldEndpointID),
// 				BoundaryScope:    GetPointerFromMap[string](data, SloConfigFieldBoundaryScope),
// 				IncludeInternal:  GetPointerFromMap[bool](data, SloConfigFieldIncludeInternal),
// 				IncludeSynthetic: GetPointerFromMap[bool](data, SloConfigFieldIncludeSynthetic),
// 				FilterExpression: tagFilter,
// 			}, nil
// 		}
// 		if details, ok := stateObject[SloConfigWebsiteEntity]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			var tagFilter *restapi.TagFilter
// 			var err error
// 			if tagFilterString, ok := data[SloConfigFieldFilterExpression]; ok {
// 				tagFilter, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
// 				if err != nil {
// 					debug(obj2json(err))
// 					// return nil, err
// 				}
// 			}

// 			return restapi.SloWebsiteEntity{
// 				Type:             SloConfigWebsiteEntity,
// 				WebsiteId:        GetPointerFromMap[string](data, SloConfigFieldWebsiteID),
// 				BeaconType:       GetPointerFromMap[string](data, SloConfigFieldBeaconType),
// 				FilterExpression: tagFilter,
// 			}, nil
// 		}
// 		if details, ok := stateObject["synthetic"]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			var tagFilter *restapi.TagFilter
// 			var err error
// 			if tagFilterString, ok := data[SloConfigFieldFilterExpression]; ok {
// 				tagFilter, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
// 				if err != nil {
// 					debug(obj2json(err))
// 					// return nil, err
// 				}
// 			}

// 			return restapi.SloSyntheticEntity{
// 				Type:             SloConfigSyntheticEntity,
// 				SyntheticTestIDs: data[SloConfigFieldSyntheticTestIDs].([]interface{}),
// 				FilterExpression: tagFilter,
// 			}, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("exactly one sli entity configuration of type is required")

// }
