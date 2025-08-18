package instana

import (
	"fmt"
	"os"
	"reflect"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

// tf resource -> api
func (r *sloConfigResource) mapSliTimeWindowListFromState(stateObject map[string]interface{}) (any, error) {
	debug(">> mapSliTimeWindowListFromState")

	if len(stateObject) > 0 {
		if details, ok := stateObject[SloConfigRollingTimeWindow]; ok && r.isSet(details) {
			data := details.([]interface{})[0].(map[string]interface{})
			return restapi.SloRollingTimeWindow{
				Type:         SloConfigRollingTimeWindow,
				Duration:     *GetPointerFromMap[int](data, SloConfigFieldDuration),
				DurationUnit: *GetPointerFromMap[string](data, SloConfigFieldDurationUnit),
				Timezone:     *GetPointerFromMap[string](data, SloConfigFieldTimezone),
			}, nil
		}
		if details, ok := stateObject[SloConfigFixedTimeWindow]; ok && r.isSet(details) {
			data := details.([]interface{})[0].(map[string]interface{})
			return restapi.SloFixedTimeWindow{
				Type:         SloConfigFixedTimeWindow,
				Duration:     *GetPointerFromMap[int](data, SloConfigFieldDuration),
				DurationUnit: *GetPointerFromMap[string](data, SloConfigFieldDurationUnit),
				Timezone:     *GetPointerFromMap[string](data, SloConfigFieldTimezone),
				StartTime: *GetPointerFromMap[float64](data, SloConfigFieldStartTimestamp),
			}, nil

		}
	}
	return nil, fmt.Errorf("exactly one time window configuration of type is required")
}

// api -> tf resource
func (r *sloConfigResource) mapSloTimeWindowToState(sloConfig *restapi.SloConfig) (map[string]interface{}, error) {
	sloTimeWindow := sloConfig.TimeWindow
	debug(">> mapSloTimeWindowToState")
	fmt.Fprintln(os.Stderr, obj2json(sloTimeWindow))

	if timeWindow, ok := sloTimeWindow.(map[string]interface{}); ok {
		if timeWindow["type"] == SloConfigRollingTimeWindow {
			result := map[string]interface{}{
				SloConfigRollingTimeWindow: []interface{}{
					map[string]interface{}{
						SloConfigFieldDuration:     timeWindow[SloConfigAPIFieldDuration],
						SloConfigFieldDurationUnit: timeWindow[SloConfigAPIFieldDurationUnit].(string),
						SloConfigFieldTimezone:     timeWindow[SloConfigAPIFieldTimezone].(string),
					},
				},
			}
			return result, nil
		} else if timeWindow["type"] == SloConfigFixedTimeWindow {
			result := map[string]interface{}{
				SloConfigFixedTimeWindow: []interface{}{
					map[string]interface{}{
						SloConfigFieldDuration:       timeWindow[SloConfigAPIFieldDuration],
						SloConfigFieldDurationUnit:   timeWindow[SloConfigAPIFieldDurationUnit].(string),
						SloConfigFieldTimezone:       timeWindow[SloConfigAPIFieldTimezone].(string),
						SloConfigFieldStartTimestamp: timeWindow[SloConfigAPIFieldStartTimestamp].(float64),
					},
				},
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("unsupported time window type %s", reflect.TypeOf(sloTimeWindow).Name())
}
