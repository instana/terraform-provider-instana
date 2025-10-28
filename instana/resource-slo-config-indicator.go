package instana

// import (
// 	"fmt"
// 	"reflect"

// 	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
// )

// // state -> api
// func (r *sloConfigResource) mapSliIndicatorListFromState(stateObject map[string]interface{}) (any, error) {
// 	debug(">> mapSliIndicatorListFromState")

// 	if len(stateObject) > 0 {
// 		if details, ok := stateObject["time_based_latency"]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			return restapi.SloTimeBasedLatencyIndicator{
// 				Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
// 				Blueprint:   SloConfigAPIIndicatorBlueprintLatency,
// 				Threshold:   *GetPointerFromMap[float64](data, SloConfigFieldThreshold),
// 				Aggregation: *GetPointerFromMap[string](data, SloConfigFieldAggregation),
// 			}, nil
// 		}
// 		if details, ok := stateObject["event_based_latency"]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			return restapi.SloEventBasedLatencyIndicator{
// 				Type:      SloConfigAPIIndicatorMeasurementTypeEventBased,
// 				Blueprint: SloConfigAPIIndicatorBlueprintLatency,
// 				Threshold: *GetPointerFromMap[float64](data, SloConfigFieldThreshold),
// 			}, nil
// 		}
// 		if details, ok := stateObject["time_based_availability"]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			return restapi.SloTimeBasedLatencyIndicator{
// 				Type:        SloConfigAPIIndicatorMeasurementTypeTimeBased,
// 				Blueprint:   SloConfigAPIIndicatorBlueprintAvailability,
// 				Threshold:   *GetPointerFromMap[float64](data, SloConfigFieldThreshold),
// 				Aggregation: *GetPointerFromMap[string](data, SloConfigFieldAggregation),
// 			}, nil
// 		}
// 		if _, ok := stateObject["event_based_availability"]; ok {
// 			return restapi.SloEventBasedLatencyIndicator{
// 				Type:      SloConfigAPIIndicatorMeasurementTypeEventBased,
// 				Blueprint: SloConfigAPIIndicatorBlueprintAvailability,
// 			}, nil
// 		}
// 		if details, ok := stateObject["custom"]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			return restapi.SloTrafficIndicator{
// 				Blueprint:   SloConfigAPIIndicatorBlueprintTraffic,
// 				TrafficType: *GetPointerFromMap[string](data, SloConfigFieldTrafficType),
// 				Threshold:   *GetPointerFromMap[float64](data, SloConfigFieldThreshold),
// 				Aggregation: *GetPointerFromMap[string](data, SloConfigFieldAggregation),
// 			}, nil
// 		}
// 		if details, ok := stateObject["traffic"]; ok && r.isSet(details) {
// 			data := details.([]interface{})[0].(map[string]interface{})
// 			var goodEventFilterExpression *restapi.TagFilter
// 			var badEventFilterExpression *restapi.TagFilter
// 			var err error
// 			if tagFilterString, ok := data[SloConfigFieldGoodEventFilterExpression]; ok {
// 				goodEventFilterExpression, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
// 				if err != nil {
// 					debug(err)
// 					// return nil, err
// 				}
// 			}
// 			if tagFilterString, ok := data[SloConfigFieldBadEventFilterExpression]; ok {
// 				badEventFilterExpression, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
// 				if err != nil {
// 					debug(err)
// 					// return nil, err
// 				}
// 			}

// 			return restapi.SloCustomIndicator{
// 				Type:                      SloConfigAPIIndicatorMeasurementTypeEventBased,
// 				Blueprint:                 SloConfigAPIIndicatorBlueprintCustom,
// 				GoodEventFilterExpression: goodEventFilterExpression,
// 				BadEventFilterExpression:  badEventFilterExpression,
// 			}, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("exactly one indicator configuration of type is required")
// }

// // api -> state
// func (r *sloConfigResource) mapSloIndicatorToState(sloConfig *restapi.SloConfig) (map[string]interface{}, error) {

// 	sloIndicator := sloConfig.Indicator

// 	if indicator, ok := sloIndicator.(map[string]interface{}); ok {
// 		if indicator[SloConfigAPIFieldType] == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator[SloConfigAPIFieldBlueprint] == SloConfigAPIIndicatorBlueprintLatency {
// 			result := map[string]interface{}{
// 				"time_based_latency": []interface{}{
// 					map[string]interface{}{
// 						SloConfigFieldThreshold:   indicator[SloConfigAPIFieldThreshold].(float64),
// 						SloConfigFieldAggregation: indicator[SloConfigAPIFieldAggregation].(string),
// 					},
// 				},
// 			}
// 			return result, nil
// 		}
// 		if indicator[SloConfigAPIFieldType] == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator[SloConfigAPIFieldBlueprint] == SloConfigAPIIndicatorBlueprintLatency {
// 			result := map[string]interface{}{
// 				"event_based_latency": []interface{}{
// 					map[string]interface{}{
// 						SloConfigFieldThreshold: indicator[SloConfigAPIFieldThreshold].(float64),
// 					},
// 				},
// 			}
// 			return result, nil
// 		}

// 		if indicator[SloConfigAPIFieldType] == SloConfigAPIIndicatorMeasurementTypeTimeBased && indicator[SloConfigAPIFieldBlueprint] == SloConfigAPIIndicatorBlueprintAvailability {
// 			result := map[string]interface{}{
// 				"time_based_latency": []interface{}{
// 					map[string]interface{}{
// 						SloConfigFieldThreshold:   indicator[SloConfigAPIFieldThreshold].(float64),
// 						SloConfigFieldAggregation: indicator[SloConfigAPIFieldAggregation].(string),
// 					},
// 				},
// 			}
// 			return result, nil
// 		}
// 		if indicator[SloConfigAPIFieldType] == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator[SloConfigAPIFieldBlueprint] == SloConfigAPIIndicatorBlueprintAvailability {
// 			result := map[string]interface{}{
// 				"event_based_latency": []interface{}{},
// 			}
// 			return result, nil
// 		}

// 		if indicator[SloConfigAPIFieldBlueprint] == "traffic" {
// 			result := map[string]interface{}{
// 				"traffic": []interface{}{
// 					map[string]interface{}{
// 						SloConfigFieldTrafficType: indicator[SloConfigAPIFieldTrafficType].(string),
// 						SloConfigFieldThreshold:   indicator[SloConfigAPIFieldThreshold].(float64),
// 						SloConfigFieldAggregation: indicator[SloConfigAPIFieldAggregation].(string),
// 					},
// 				},
// 			}
// 			return result, nil
// 		}
// 		if indicator[SloConfigAPIFieldType] == SloConfigAPIIndicatorMeasurementTypeEventBased && indicator[SloConfigAPIFieldBlueprint] == SloConfigAPIIndicatorBlueprintCustom {
// 			goodEventFilterExp, validExp1, err1 := fiterExpFromAPIModel(indicator[SloConfigAPIFieldGoodEventFilter])
// 			if validExp1 {
// 				return nil, err1
// 			}
// 			badEventFilterExp, validExp2, err2 := fiterExpFromAPIModel(indicator[SloConfigAPIFieldBadEventFilter])
// 			if validExp2 {
// 				return nil, err2
// 			}

// 			result := map[string]interface{}{
// 				"custom": []interface{}{
// 					map[string]interface{}{
// 						SloConfigFieldGoodEventFilterExpression: goodEventFilterExp,
// 						SloConfigFieldBadEventFilterExpression:  badEventFilterExp,
// 					},
// 				},
// 			}
// 			return result, nil
// 		}

// 	}
// 	return nil, fmt.Errorf("Unsupported indicator type %s", reflect.TypeOf(sloIndicator).Name())
// }
