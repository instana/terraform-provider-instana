package customdashboard

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// widgetsEquivalent reports whether two widget JSON documents are equal, ignoring
// formatting, object key order and the order of the elements of a top level array.
// The Instana API does not return the widgets in a stable order.
func widgetsEquivalent(a, b string) bool {
	va, err := decodeAndSortWidgets(a)
	if err != nil {
		return false
	}
	vb, err := decodeAndSortWidgets(b)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

// decodeAndSortWidgets decodes the JSON document. A top level array is replaced by the sorted
// list of the canonical JSON of its elements, which makes the comparison order independent.
func decodeAndSortWidgets(input string) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return nil, err
	}
	list, ok := v.([]interface{})
	if !ok {
		return v, nil
	}

	elements := make([]string, len(list))
	for i, e := range list {
		b, err := json.Marshal(e) // map keys are marshalled in sorted order
		if err != nil {
			return nil, err
		}
		elements[i] = string(b)
	}
	sort.Strings(elements)
	return elements, nil
}

// widgetsKeepStatePlanModifier keeps the widgets from the state when the planned value
// only differs from it in formatting or widget order.
type widgetsKeepStatePlanModifier struct{}

func (m widgetsKeepStatePlanModifier) Description(_ context.Context) string {
	return "Keeps the state value when the widgets only differ in formatting or widget order."
}

func (m widgetsKeepStatePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m widgetsKeepStatePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() {
		return // resource is being created
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if widgetsEquivalent(req.StateValue.ValueString(), req.PlanValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}
