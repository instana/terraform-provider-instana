package customdashboard

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRelatedIDRequiredValidator(t *testing.T) {
	attrTypes := map[string]attr.Type{
		CustomDashboardFieldAccessRuleAccessType:   types.StringType,
		CustomDashboardFieldAccessRuleRelatedID:    types.StringType,
		CustomDashboardFieldAccessRuleRelationType: types.StringType,
	}

	validate := func(obj types.Object) bool {
		req := validator.ObjectRequest{
			Path:        path.Root(CustomDashboardFieldAccessRule).AtListIndex(0),
			ConfigValue: obj,
		}
		resp := &validator.ObjectResponse{}
		relatedIDRequiredValidator{}.ValidateObject(context.Background(), req, resp)
		return resp.Diagnostics.HasError()
	}

	rule := func(relationType, relatedID types.String) types.Object {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{
			CustomDashboardFieldAccessRuleAccessType:   types.StringValue("READ"),
			CustomDashboardFieldAccessRuleRelationType: relationType,
			CustomDashboardFieldAccessRuleRelatedID:    relatedID,
		})
	}

	t.Run("GLOBAL without related_id is valid", func(t *testing.T) {
		assert.False(t, validate(rule(types.StringValue("GLOBAL"), types.StringNull())))
	})
	t.Run("GLOBAL with related_id is valid", func(t *testing.T) {
		assert.False(t, validate(rule(types.StringValue("GLOBAL"), types.StringValue("x"))))
	})
	for _, rt := range []string{"USER", "API_TOKEN", "ROLE", "TEAM"} {
		t.Run(rt+" with related_id is valid", func(t *testing.T) {
			assert.False(t, validate(rule(types.StringValue(rt), types.StringValue("id-1"))))
		})
		t.Run(rt+" without related_id is invalid", func(t *testing.T) {
			assert.True(t, validate(rule(types.StringValue(rt), types.StringNull())))
		})
		t.Run(rt+" with empty related_id is invalid", func(t *testing.T) {
			assert.True(t, validate(rule(types.StringValue(rt), types.StringValue(""))))
		})
	}
	t.Run("unknown related_id is skipped", func(t *testing.T) {
		assert.False(t, validate(rule(types.StringValue("USER"), types.StringUnknown())))
	})
	t.Run("unknown relation_type is skipped", func(t *testing.T) {
		assert.False(t, validate(rule(types.StringUnknown(), types.StringNull())))
	})
	t.Run("null object is skipped", func(t *testing.T) {
		assert.False(t, validate(types.ObjectNull(attrTypes)))
	})
}
