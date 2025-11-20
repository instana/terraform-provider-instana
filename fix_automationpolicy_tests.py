#!/usr/bin/env python3
import re

# Read the test file
with open('internal/resources/automationpolicy/resource-automation-policy-framework_test.go', 'r') as f:
    content = f.read()

# Helper function to add before the last closing brace
helper_function = '''
// initializeEmptyState initializes the state with an empty AutomationPolicyModel
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := AutomationPolicyModel{
		ID:                types.StringNull(),
		Name:              types.StringNull(),
		Description:       types.StringNull(),
		Tags:              types.ListNull(types.StringType),
		Trigger:           TriggerModel{
			ID:          types.StringNull(),
			Type:        types.StringNull(),
			Name:        types.StringNull(),
			Description: types.StringNull(),
			Scheduling:  nil,
		},
		TypeConfiguration: []TypeConfigurationModel{},
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}

// Made with Bob'''

# Add helper function before the last line
lines = content.split('\n')
# Find the last line (should be "// Made with Bob")
last_line_idx = len(lines) - 1
while last_line_idx >= 0 and lines[last_line_idx].strip() == '':
    last_line_idx -= 1

# Insert helper function before the last line
lines.insert(last_line_idx, helper_function)
content = '\n'.join(lines)

# Pattern 1: For tests using getTestSchema()
# Look for: state := &tfsdk.State{\n\t\tSchema: getTestSchema(),\n\t}
# And add initialization after it
pattern1 = r'(state := &tfsdk\.State\{\s+Schema: getTestSchema\(\),\s+\})\s+(diags := resource\.UpdateState)'

replacement1 = r'\1\n\n\t// Initialize state with empty model\n\tinitializeEmptyState(t, ctx, state)\n\n\t\2'

content = re.sub(pattern1, replacement1, content)

# Pattern 2: For tests using resource.metaData.Schema (if any)
pattern2 = r'(state := &tfsdk\.State\{\s+Schema: resource\.metaData\.Schema,\s+\})\s+(diags := resource\.UpdateState)'

replacement2 = r'\1\n\n\t// Initialize state with empty model\n\tinitializeEmptyState(t, ctx, state)\n\n\t\2'

content = re.sub(pattern2, replacement2, content)

# Write back
with open('internal/resources/automationpolicy/resource-automation-policy-framework_test.go', 'w') as f:
    f.write(content)

print("Fixed automationpolicy test file")

# Made with Bob
