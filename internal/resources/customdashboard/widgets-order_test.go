package customdashboard

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestWidgetsEquivalent(t *testing.T) {
	a := `[{"id":"b","x":0},{"id":"a","x":1}]`
	b := `[ {"x":1,"id":"a"}, {"x":0,"id":"b"} ]`
	assert.True(t, widgetsEquivalent(a, b))
	assert.False(t, widgetsEquivalent(a, `[{"id":"a","x":2},{"id":"b","x":0}]`))
	assert.False(t, widgetsEquivalent(a, `[{"id":"a","x":1}]`))
	assert.False(t, widgetsEquivalent(a, `not json`))
}

func TestWidgetsKeepStatePlanModifier(t *testing.T) {
	stateRaw := tftypes.NewValue(tftypes.String, "x")
	state := tfsdk.State{Raw: stateRaw}
	run := func(st, plan string, creating bool) types.String {
		req := planmodifier.StringRequest{
			StateValue: types.StringValue(st),
			PlanValue:  types.StringValue(plan),
			State:      state,
		}
		if creating {
			req.State = tfsdk.State{Raw: tftypes.NewValue(tftypes.String, nil)}
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		widgetsKeepStatePlanModifier{}.PlanModifyString(context.Background(), req, resp)
		return resp.PlanValue
	}

	assert.Equal(t, types.StringValue(`[{"id":"a"},{"id":"b"}]`), run(`[{"id":"a"},{"id":"b"}]`, `[{"id":"b"},{"id":"a"}]`, false))
	assert.Equal(t, types.StringValue(`[{"id":"c"}]`), run(`[{"id":"a"}]`, `[{"id":"c"}]`, false))
	assert.Equal(t, types.StringValue(`[{"id":"b"},{"id":"a"}]`), run(`[{"id":"a"},{"id":"b"}]`, `[{"id":"b"},{"id":"a"}]`, true))
}
