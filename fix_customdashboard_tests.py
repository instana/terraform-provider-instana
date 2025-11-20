#!/usr/bin/env python3
import re

# Read the test file
with open('internal/resources/customdashboard/resource-custom-dashboard-framework_test.go', 'r') as f:
    content = f.read()

# Helper function to add before the last closing brace
helper_function = '''
// initializeEmptyState initializes the state with an empty CustomDashboardModel
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := CustomDashboardModel{
		ID:          types.StringNull(),
		Title:       types.StringNull(),
		AccessRules: []AccessRuleModel{},
		Widgets:     jsontypes.NewNormalizedNull(),
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

# Pattern 1: For tests using resource.metaData.Schema
pattern1 = r'(state := &tfsdk\.State\{\s+Schema: resource\.metaData\.Schema,\s+\})\s+(diags := resource\.UpdateState)'

replacement1 = r'\1\n\n\t\t// Initialize state with empty model\n\t\tinitializeEmptyState(t, ctx, state)\n\n\t\t\2'

content = re.sub(pattern1, replacement1, content)

# Pattern 2: For tests using handle.MetaData().Schema
pattern2 = r'(state := &tfsdk\.State\{\s+Schema: handle\.MetaData\(\)\.Schema,\s+\})\s+(diags := resource\.UpdateState)'

replacement2 = r'\1\n\n\t\t// Initialize state with empty model\n\t\tinitializeEmptyState(t, ctx, state)\n\n\t\t\2'

content = re.sub(pattern2, replacement2, content)

# Write back
with open('internal/resources/customdashboard/resource-custom-dashboard-framework_test.go', 'w') as f:
    f.write(content)

print("Fixed customdashboard test file")

# Made with Bob
