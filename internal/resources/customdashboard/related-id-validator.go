package customdashboard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	model "github.com/instana/instana-go-client/shared/types"
)

// relatedIDRequiredValidator validates that related_id is set for every relation_type
// except GLOBAL. The Instana API rejects such an access rule with "Related ID must be defined."
type relatedIDRequiredValidator struct{}

var _ validator.Object = relatedIDRequiredValidator{}

func (v relatedIDRequiredValidator) Description(_ context.Context) string {
	return "related_id must be set unless relation_type is GLOBAL."
}

func (v relatedIDRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v relatedIDRequiredValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()
	relationType, ok := attrs[CustomDashboardFieldAccessRuleRelationType].(types.String)
	if !ok || relationType.IsNull() || relationType.IsUnknown() {
		return
	}
	if model.RelationType(relationType.ValueString()) == model.RelationTypeGlobal {
		return
	}

	relatedID, ok := attrs[CustomDashboardFieldAccessRuleRelatedID].(types.String)
	if !ok || relatedID.IsUnknown() {
		return
	}
	if relatedID.IsNull() || relatedID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName(CustomDashboardFieldAccessRuleRelatedID),
			"Missing related_id",
			CustomDashboardFieldAccessRuleRelatedID+" must be set when "+
				CustomDashboardFieldAccessRuleRelationType+" is "+relationType.ValueString()+
				". It is only optional for "+string(model.RelationTypeGlobal)+".",
		)
	}
}
