package customeventspec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// createStateUpgrader creates a state upgrader which converts the state of any prior schema version
// to the current schema (version 2).
//
// Prior state shapes handled:
//   - version 0 (SDK v2 provider, v5.x): rules is a list with a single element and every rule type
//     (entity_count, ..., threshold) is a list; threshold may hold up to 5 elements and metric_pattern is a list.
//   - version 1 (Plugin Framework provider, v6.x - v7.x): rules is an object and every rule type is an object.
//
// In the current schema threshold is a list of objects (up to 5), all other rule types and metric_pattern
// are objects. The upgrade works on the raw JSON state so no prior schema has to be redeclared.
func (r *customEventSpecificationResource) createStateUpgrader() resource.StateUpgrader {
	return resource.StateUpgrader{
		StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
			if req.RawState == nil || req.RawState.JSON == nil {
				resp.Diagnostics.AddError(
					"State Upgrade Error",
					"The prior state of the custom event specification is not available in JSON format and cannot be upgraded.",
				)
				return
			}

			decoder := json.NewDecoder(bytes.NewReader(req.RawState.JSON))
			decoder.UseNumber()
			var state map[string]any
			if err := decoder.Decode(&state); err != nil {
				resp.Diagnostics.AddError("State Upgrade Error", fmt.Sprintf("Failed to parse the prior state: %s", err))
				return
			}

			upgradeRawState(state)

			upgradedJSON, err := json.Marshal(state)
			if err != nil {
				resp.Diagnostics.AddError("State Upgrade Error", fmt.Sprintf("Failed to serialize the upgraded state: %s", err))
				return
			}

			schemaType := createCustomEventSpecificationSchema().Type().TerraformType(ctx)
			value, err := tftypes.ValueFromJSONWithOpts(upgradedJSON, schemaType, tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}) //nolint:staticcheck
			if err != nil {
				resp.Diagnostics.AddError("State Upgrade Error", fmt.Sprintf("Failed to convert the upgraded state: %s", err))
				return
			}

			dynamicValue, err := tfprotov6.NewDynamicValue(schemaType, value)
			if err != nil {
				resp.Diagnostics.AddError("State Upgrade Error", fmt.Sprintf("Failed to encode the upgraded state: %s", err))
				return
			}
			resp.DynamicValue = &dynamicValue
		},
	}
}

// upgradeRawState converts the raw JSON state of a prior schema version to the shape of the current schema
func upgradeRawState(state map[string]any) {
	switch rules := state[CustomEventSpecificationFieldRules].(type) {
	case []any:
		// version 0: rules is a list holding a single element
		if len(rules) == 0 {
			state[CustomEventSpecificationFieldRules] = nil
			return
		}
		if first, ok := rules[0].(map[string]any); ok {
			upgradeRawRules(first)
			state[CustomEventSpecificationFieldRules] = first
		} else {
			state[CustomEventSpecificationFieldRules] = nil
		}
	case map[string]any:
		// version 1: rules is already an object
		upgradeRawRules(rules)
	}
}

// upgradeRawRules converts the rule types of the rules container
func upgradeRawRules(rules map[string]any) {
	for _, ruleType := range []string{
		CustomEventSpecificationFieldEntityCountRule,
		CustomEventSpecificationFieldEntityCountVerificationRule,
		CustomEventSpecificationFieldEntityVerificationRule,
		CustomEventSpecificationFieldHostAvailabilityRule,
		CustomEventSpecificationFieldSystemRule,
	} {
		if list, ok := rules[ruleType].([]any); ok {
			rules[ruleType] = firstOrNil(list)
		}
	}

	var thresholds []any
	switch threshold := rules[CustomEventSpecificationFieldThresholdRule].(type) {
	case []any:
		thresholds = threshold
	case map[string]any:
		thresholds = []any{threshold}
	}

	if len(thresholds) == 0 {
		rules[CustomEventSpecificationFieldThresholdRule] = nil
		return
	}
	for _, threshold := range thresholds {
		if rule, ok := threshold.(map[string]any); ok {
			upgradeRawMetricPattern(rule)
		}
	}
	rules[CustomEventSpecificationFieldThresholdRule] = thresholds
}

// upgradeRawMetricPattern converts the metric pattern of a threshold rule from a list to an object
func upgradeRawMetricPattern(rule map[string]any) {
	if list, ok := rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]any); ok {
		rule[CustomEventSpecificationThresholdRuleFieldMetricPattern] = firstOrNil(list)
	}

	pattern, ok := rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].(map[string]any)
	if !ok {
		return
	}
	// postfix and placeholder default to an empty string in the current schema
	for _, field := range []string{
		CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix,
		CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder,
	} {
		if pattern[field] == nil {
			pattern[field] = CustomEventSpecificationDefaultEmptyString
		}
	}
}

// firstOrNil returns the first element of the list or nil for an empty list
func firstOrNil(list []any) any {
	if len(list) == 0 {
		return nil
	}
	return list[0]
}
