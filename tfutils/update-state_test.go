package tfutils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestUpdateStatePlugin(t *testing.T) {
	ctx := context.Background()

	// Create a schema for testing
	testSchema := fwschema.Schema{
		Attributes: map[string]fwschema.Attribute{
			"string_attr": fwschema.StringAttribute{
				Optional: true,
			},
			"bool_attr": fwschema.BoolAttribute{
				Optional: true,
			},
			"int_attr": fwschema.Int64Attribute{
				Optional: true,
			},
			"float_attr": fwschema.Float64Attribute{
				Optional: true,
			},
			"list_attr": fwschema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"set_attr": fwschema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"map_attr": fwschema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}

	// Create a state with the schema
	state := tfsdk.State{
		Schema: testSchema,
		Raw: tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"string_attr": tftypes.String,
				"bool_attr":   tftypes.Bool,
				"int_attr":    tftypes.Number,
				"float_attr":  tftypes.Number,
				"list_attr":   tftypes.List{ElementType: tftypes.String},
				"set_attr":    tftypes.Set{ElementType: tftypes.String},
				"map_attr":    tftypes.Map{ElementType: tftypes.String},
			},
		}, map[string]tftypes.Value{
			"string_attr": tftypes.NewValue(tftypes.String, nil),
			"bool_attr":   tftypes.NewValue(tftypes.Bool, nil),
			"int_attr":    tftypes.NewValue(tftypes.Number, nil),
			"float_attr":  tftypes.NewValue(tftypes.Number, nil),
			"list_attr":   tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
			"set_attr":    tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
			"map_attr":    tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}),
	}

	// Create test data
	testData := map[string]interface{}{
		"string_attr": "test string",
		"bool_attr":   true,
		"int_attr":    int64(42),
		"float_attr":  float64(3.14),
		"list_attr":   []string{"item1", "item2"},
		"set_attr":    []string{"set1", "set2"},
		"map_attr":    map[string]string{"key1": "value1", "key2": "value2"},
	}

	// Update the state
	err := UpdateStatePlugin(ctx, &state, testData)
	assert.NoError(t, err)

	// Verify the state was updated correctly
	var stringVal types.String
	diags := state.GetAttribute(ctx, path.Root("string_attr"), &stringVal)
	assert.False(t, diags.HasError())
	assert.Equal(t, "test string", stringVal.ValueString())

	var boolVal types.Bool
	diags = state.GetAttribute(ctx, path.Root("bool_attr"), &boolVal)
	assert.False(t, diags.HasError())
	assert.Equal(t, true, boolVal.ValueBool())

	var intVal types.Int64
	diags = state.GetAttribute(ctx, path.Root("int_attr"), &intVal)
	assert.False(t, diags.HasError())
	assert.Equal(t, int64(42), intVal.ValueInt64())

	var floatVal types.Float64
	diags = state.GetAttribute(ctx, path.Root("float_attr"), &floatVal)
	assert.False(t, diags.HasError())
	assert.Equal(t, float64(3.14), floatVal.ValueFloat64())

	var listVal types.List
	diags = state.GetAttribute(ctx, path.Root("list_attr"), &listVal)
	assert.False(t, diags.HasError())
	var listItems []string
	diags = listVal.ElementsAs(ctx, &listItems, false)
	assert.False(t, diags.HasError())
	assert.Equal(t, []string{"item1", "item2"}, listItems)

	var setVal types.Set
	diags = state.GetAttribute(ctx, path.Root("set_attr"), &setVal)
	assert.False(t, diags.HasError())
	var setItems []string
	diags = setVal.ElementsAs(ctx, &setItems, false)
	assert.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"set1", "set2"}, setItems)

	var mapVal types.Map
	diags = state.GetAttribute(ctx, path.Root("map_attr"), &mapVal)
	assert.False(t, diags.HasError())
	var mapItems map[string]string
	diags = mapVal.ElementsAs(ctx, &mapItems, false)
	assert.False(t, diags.HasError())
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, mapItems)
}
