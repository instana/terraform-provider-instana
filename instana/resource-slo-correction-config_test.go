package instana

import (
	"context"
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSloCorrectionConfig(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaSloCorrectionConfig + ".example"
	inst := &sloCorrectionConfigTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewSloCorrectionConfigResourceHandle(),
	}
	inst.run(t)
}

type sloCorrectionConfigTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle[*restapi.SloCorrectionConfig]
}

func (test *sloCorrectionConfigTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("%s_HasCorrectResourceName", ResourceInstanaSloCorrectionConfig), test.createTestResourceShouldHaveResourceName())
	t.Run(fmt.Sprintf("%s_HasOneStateUpgrader", ResourceInstanaSloCorrectionConfig), test.createTestResourceShouldHaveOneStateUpgrader())
	t.Run(fmt.Sprintf("%s_HasSchemaVersionOne", ResourceInstanaSloCorrectionConfig), test.createTestResourceShouldHaveSchemaVersionOne())
	t.Run("MapStateToDataObject_MapsAllFieldsCorrectly", test.testMapStateToDataObject())
	t.Run("UpdateState_SetsAllFieldsCorrectly", test.testUpdateState())
	t.Run("StateUpgradeV0_MigratesFullNameToName", test.testSloCorrectionConfigStateUpgradeV0())
	t.Run("SetComputedFields_ReturnsNil", test.testSetComputedFieldsReturnsNil())
	t.Run("StateUpgraders_ReturnsOneUpgraderWithVersionZero", test.testStateUpgraders())
	t.Run("MapStateToDataObject_NoTags_EmptyTagsSlice", test.testMapStateToDataObjectWithNoTags())
	t.Run("UpdateState_NoTags_EmptyTagsSlice", test.testUpdateStateWithNoTags())
	t.Run("MapStateToDataObject_IntStartTime_HandlesIntType", test.testMapStateToDataObjectWithIntStartTime())
}

// creating a test SLO Correction Config.
func testResourceDataForSloCorrectionConfig(t *testing.T) *schema.ResourceData {
	r := NewSloCorrectionConfigResourceHandle()
	resource := r.MetaData().Schema
	d := schema.TestResourceDataRaw(t, resource, map[string]interface{}{
		SloCorrectionConfigFieldName:        "test-correction",
		SloCorrectionConfigFieldDescription: "desc",
		SloCorrectionConfigFieldActive:      true,
		SloCorrectionConfigFieldScheduling: []interface{}{
			map[string]interface{}{
				SloCorrectionConfigFieldSchedulingStartTime:     1741600800000,
				SloCorrectionConfigFieldSchedulingDuration:      60,
				SloCorrectionConfigFieldSchedulingDurationUnit:  "MINUTE",
				SloCorrectionConfigFieldSchedulingRecurrentRule: "FREQ=DAILY",
			},
		},
		SloCorrectionConfigFieldSloIds: []interface{}{"slo-1", "slo-2"},
		SloCorrectionConfigFieldTags:   []interface{}{"tag1", "tag2"},
	})
	d.SetId("test-id")
	return d
}

func (test *sloCorrectionConfigTest) createTestResourceShouldHaveResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, test.resourceHandle.MetaData().ResourceName, "instana_slo_correction_config", "Resource name should match the expected value")
	}
}

func (test *sloCorrectionConfigTest) createTestResourceShouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Len(t, test.resourceHandle.StateUpgraders(), 1)
	}
}

func (test *sloCorrectionConfigTest) createTestResourceShouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, test.resourceHandle.MetaData().SchemaVersion)
	}
}

func (test *sloCorrectionConfigTest) testMapStateToDataObject() func(t *testing.T) {
	return func(t *testing.T) {
		r := NewSloCorrectionConfigResourceHandle()
		d := testResourceDataForSloCorrectionConfig(t)
		obj, err := r.MapStateToDataObject(d)
		require.NoError(t, err)
		require.Equal(t, "test-id", obj.ID)
		require.Equal(t, "test-correction", obj.Name)
		require.Equal(t, "desc", obj.Description)
		require.True(t, obj.Active)
		require.Equal(t, int64(1741600800000), obj.Scheduling.StartTime)
		require.Equal(t, 60, obj.Scheduling.Duration)
		require.Equal(t, restapi.DurationUnit("MINUTE"), obj.Scheduling.DurationUnit)
		require.Equal(t, "FREQ=DAILY", obj.Scheduling.RecurrentRule)
		require.ElementsMatch(t, []string{"slo-1", "slo-2"}, obj.SloIds)
		require.ElementsMatch(t, []string{"tag1", "tag2"}, obj.Tags)
	}
}

func (test *sloCorrectionConfigTest) testUpdateState() func(t *testing.T) {
	return func(t *testing.T) {
		r := NewSloCorrectionConfigResourceHandle()
		d := schema.TestResourceDataRaw(t, r.MetaData().Schema, map[string]interface{}{})
		obj := &restapi.SloCorrectionConfig{
			ID:          "id-123",
			Name:        "name-1",
			Description: "desc-1",
			Active:      false,
			Scheduling: restapi.Scheduling{
				StartTime:     1741600800000,
				Duration:      30,
				DurationUnit:  restapi.DurationUnit("HOUR"),
				RecurrentRule: "FREQ=WEEKLY",
			},
			SloIds: []string{"slo-x"},
			Tags:   []string{"tag-x"},
		}
		err := r.UpdateState(d, obj)
		require.NoError(t, err)
		require.Equal(t, "id-123", d.Id())
		require.Equal(t, "name-1", d.Get(SloCorrectionConfigFieldName))
		require.Equal(t, "desc-1", d.Get(SloCorrectionConfigFieldDescription))
		require.Equal(t, false, d.Get(SloCorrectionConfigFieldActive))
		scheduling := d.Get(SloCorrectionConfigFieldScheduling).([]interface{})[0].(map[string]interface{})
		require.Equal(t, 1741600800000, scheduling[SloCorrectionConfigFieldSchedulingStartTime])
		require.Equal(t, 30, scheduling[SloCorrectionConfigFieldSchedulingDuration])
		require.Equal(t, "HOUR", scheduling[SloCorrectionConfigFieldSchedulingDurationUnit])
		require.Equal(t, "FREQ=WEEKLY", scheduling[SloCorrectionConfigFieldSchedulingRecurrentRule])
		require.Contains(t, d.Get(SloCorrectionConfigFieldSloIds).(*schema.Set).List(), "slo-x")
		require.Contains(t, d.Get(SloCorrectionConfigFieldTags).(*schema.Set).List(), "tag-x")
	}
}

func (test *sloCorrectionConfigTest) testSloCorrectionConfigStateUpgradeV0() func(t *testing.T) {
	return func(t *testing.T) {
		r := &sloCorrectionConfigResource{}
		state := map[string]interface{}{
			SloCorrectionConfigFieldFullName: "full-name",
		}
		newState, err := r.sloCorrectionConfigStateUpgradeV0(context.Background(), state, nil)
		require.NoError(t, err)
		require.Equal(t, "full-name", newState[SloCorrectionConfigFieldName])
		_, exists := newState[SloCorrectionConfigFieldFullName]
		require.False(t, exists)
	}
}

func (test *sloCorrectionConfigTest) testSetComputedFieldsReturnsNil() func(t *testing.T) {
	return func(t *testing.T) {
		r := &sloCorrectionConfigResource{}
		require.NoError(t, r.SetComputedFields(nil))
	}
}

func (test *sloCorrectionConfigTest) testStateUpgraders() func(t *testing.T) {
	return func(t *testing.T) {
		r := &sloCorrectionConfigResource{}
		upgraders := r.StateUpgraders()
		require.Len(t, upgraders, 1)
		require.Equal(t, 0, upgraders[0].Version)
	}
}

func (test *sloCorrectionConfigTest) testMapStateToDataObjectWithNoTags() func(t *testing.T) {
	return func(t *testing.T) {
		r := NewSloCorrectionConfigResourceHandle()
		resource := r.MetaData().Schema
		d := schema.TestResourceDataRaw(t, resource, map[string]interface{}{
			SloCorrectionConfigFieldName:        "no-tags",
			SloCorrectionConfigFieldDescription: "desc",
			SloCorrectionConfigFieldActive:      true,
			SloCorrectionConfigFieldScheduling: []interface{}{
				map[string]interface{}{
					SloCorrectionConfigFieldSchedulingStartTime:     1741600800000,
					SloCorrectionConfigFieldSchedulingDuration:      60,
					SloCorrectionConfigFieldSchedulingDurationUnit:  "MINUTE",
					SloCorrectionConfigFieldSchedulingRecurrentRule: "FREQ=DAILY",
				},
			},
			SloCorrectionConfigFieldSloIds: []interface{}{"slo-1"},
		})
		d.SetId("id-no-tags")
		obj, err := r.MapStateToDataObject(d)
		require.NoError(t, err)
		require.Equal(t, "id-no-tags", obj.ID)
		require.Equal(t, "no-tags", obj.Name)
		require.Equal(t, "desc", obj.Description)
		require.True(t, obj.Active)
		require.Equal(t, int64(1741600800000), obj.Scheduling.StartTime)
		require.Equal(t, 60, obj.Scheduling.Duration)
		require.Equal(t, restapi.DurationUnit("MINUTE"), obj.Scheduling.DurationUnit)
		require.Equal(t, "FREQ=DAILY", obj.Scheduling.RecurrentRule)
		require.ElementsMatch(t, []string{"slo-1"}, obj.SloIds)
		require.Empty(t, obj.Tags)
	}
}

func (test *sloCorrectionConfigTest) testUpdateStateWithNoTags() func(t *testing.T) {
	return func(t *testing.T) {
		r := NewSloCorrectionConfigResourceHandle()
		d := schema.TestResourceDataRaw(t, r.MetaData().Schema, map[string]interface{}{})
		obj := &restapi.SloCorrectionConfig{
			ID:          "id-no-tags",
			Name:        "name-no-tags",
			Description: "desc-no-tags",
			Active:      true,
			Scheduling: restapi.Scheduling{
				StartTime:     1741600800000,
				Duration:      10,
				DurationUnit:  restapi.DurationUnit("DAY"),
				RecurrentRule: "",
			},
			SloIds: []string{"slo-y"},
		}
		err := r.UpdateState(d, obj)
		require.NoError(t, err)
		require.Equal(t, "id-no-tags", d.Id())
		require.Equal(t, "name-no-tags", d.Get(SloCorrectionConfigFieldName))
		require.Equal(t, "desc-no-tags", d.Get(SloCorrectionConfigFieldDescription))
		require.Equal(t, true, d.Get(SloCorrectionConfigFieldActive))
		scheduling := d.Get(SloCorrectionConfigFieldScheduling).([]interface{})[0].(map[string]interface{})
		require.Equal(t, 1741600800000, scheduling[SloCorrectionConfigFieldSchedulingStartTime])
		require.Equal(t, 10, scheduling[SloCorrectionConfigFieldSchedulingDuration])
		require.Equal(t, "DAY", scheduling[SloCorrectionConfigFieldSchedulingDurationUnit])
		require.Equal(t, "", scheduling[SloCorrectionConfigFieldSchedulingRecurrentRule])
		require.Contains(t, d.Get(SloCorrectionConfigFieldSloIds).(*schema.Set).List(), "slo-y")
		require.Empty(t, d.Get(SloCorrectionConfigFieldTags).(*schema.Set).List())
	}
}

func (test *sloCorrectionConfigTest) testMapStateToDataObjectWithIntStartTime() func(t *testing.T) {
	return func(t *testing.T) {
		r := NewSloCorrectionConfigResourceHandle()
		resource := r.MetaData().Schema
		d := schema.TestResourceDataRaw(t, resource, map[string]interface{}{
			SloCorrectionConfigFieldName:        "int-start-time",
			SloCorrectionConfigFieldDescription: "desc",
			SloCorrectionConfigFieldActive:      true,
			SloCorrectionConfigFieldScheduling: []interface{}{
				map[string]interface{}{
					SloCorrectionConfigFieldSchedulingStartTime:     1741600800000,
					SloCorrectionConfigFieldSchedulingDuration:      60,
					SloCorrectionConfigFieldSchedulingDurationUnit:  "MINUTE",
					SloCorrectionConfigFieldSchedulingRecurrentRule: "FREQ=DAILY",
				},
			},
			SloCorrectionConfigFieldSloIds: []interface{}{"slo-1"},
			SloCorrectionConfigFieldTags:   []interface{}{"tag1"},
		})
		d.SetId("id-int-start-time")
		obj, err := r.MapStateToDataObject(d)
		require.NoError(t, err)
		require.Equal(t, "id-int-start-time", obj.ID)
		require.Equal(t, "int-start-time", obj.Name)
		require.Equal(t, "desc", obj.Description)
		require.True(t, obj.Active)
		require.Equal(t, int64(1741600800000), obj.Scheduling.StartTime)
		require.Equal(t, 60, obj.Scheduling.Duration)
		require.Equal(t, restapi.DurationUnit("MINUTE"), obj.Scheduling.DurationUnit)
		require.Equal(t, "FREQ=DAILY", obj.Scheduling.RecurrentRule)
		require.ElementsMatch(t, []string{"slo-1"}, obj.SloIds)
		require.ElementsMatch(t, []string{"tag1"}, obj.Tags)
	}
}
