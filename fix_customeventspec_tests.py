#!/usr/bin/env python3
import re

# Read the test file
with open('internal/resources/customeventspec/resource-custom-event-specification-framework_test.go', 'r') as f:
    content = f.read()

# Helper function to add before the last closing brace
helper_function = '''
// initializeEmptyState initializes the state with an empty CustomEventSpecificationModel
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := CustomEventSpecificationModel{
		ID:                  types.StringNull(),
		Name:                types.StringNull(),
		EntityType:          types.StringNull(),
		Query:               types.StringNull(),
		Triggering:          types.BoolNull(),
		Description:         types.StringNull(),
		ExpirationTime:      types.Int64Null(),
		Enabled:             types.BoolNull(),
		RuleLogicalOperator: types.StringNull(),
		Rules:               nil,
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
pattern1 = r'(state := &tfsdk\.State\{\s+Schema: getTestSchema\(\),\s+\})\s+(diags := resource\.UpdateState)'

replacement1 = r'\1\n\n\t// Initialize state with empty model\n\tinitializeEmptyState(t, ctx, state)\n\n\t\2'

content = re.sub(pattern1, replacement1, content)

# Pattern 2: For tests using resource.metaData.Schema (if any)
pattern2 = r'(state := &tfsdk\.State\{\s+Schema: resource\.metaData\.Schema,\s+\})\s+(diags := resource\.UpdateState)'

replacement2 = r'\1\n\n\t// Initialize state with empty model\n\tinitializeEmptyState(t, ctx, state)\n\n\t\2'

content = re.sub(pattern2, replacement2, content)

# Write back
with open('internal/resources/customeventspec/resource-custom-event-specification-framework_test.go', 'w') as f:
    f.write(content)

print("Fixed customeventspec test file")

# Made with Bob
