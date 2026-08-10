package roles

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
)

// NewRoleResourceHandle creates the resource handle for RBAC Roles
func NewRoleResourceHandle() resourcehandle.ResourceHandle[*api.Role] {
	return &roleResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaRole,
			Schema:        buildRoleSchema(),
			SchemaVersion: 1,
		},
	}
}

// buildRoleSchema constructs the Terraform schema for the role resource
func buildRoleSchema() schema.Schema {
	return schema.Schema{
		Description: RoleDescResource,
		Attributes: map[string]schema.Attribute{
			RoleFieldID: schema.StringAttribute{
				Computed:    true,
				Description: RoleDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RoleFieldName: schema.StringAttribute{
				Required:    true,
				Description: RoleDescName,
			},
			RoleFieldMembers: schema.SetNestedAttribute{
				Description: RoleDescMembers,
				Optional:    true,
				Computed:    true,
				// UseStateForUnknown: when the user omits the member block, the plan would
				// normally be "unknown". This modifier fills it with the current state value
				// instead, so Terraform can detect drift on subsequent plans.
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: buildMemberNestedObject(),
			},
			RoleFieldPermissions: schema.SetAttribute{
				Required:    true,
				Description: RoleDescPermissions,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(api.SupportedInstanaPermissions.ToStringSlice()...),
					),
				},
			},
		},
	}
}

// buildMemberNestedObject constructs the nested object schema for role members
func buildMemberNestedObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			RoleFieldMemberUserID: schema.StringAttribute{
				Required:    true,
				Description: RoleDescMemberUserID,
			},
		},
	}
}

type roleResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *roleResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *roleResource) GetRestResource(api client.InstanaAPI) rest.RestResource[*api.Role] {
	return api.Roles()
}

func (r *roleResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *roleResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, role *api.Role) diag.Diagnostics {
	var diags diag.Diagnostics

	membersSet, membersDiags := r.mapMembersToSet(ctx, role.Members)
	diags.Append(membersDiags...)
	if diags.HasError() {
		return diags
	}

	model := RoleModel{
		ID:          types.StringValue(role.ID),
		Name:        types.StringValue(role.Name),
		Members:     membersSet,
		Permissions: role.Permissions,
	}
	return state.Set(ctx, model)
}

// mapMembersToSet converts API members to a types.Set.
// The Roles API always returns "members": [] — using types.Set (instead of []RoleMemberModel)
// lets Terraform hold the unknown/null/empty distinction correctly when the field is Optional+Computed.
func (r *roleResource) mapMembersToSet(ctx context.Context, apiMembers []api.APIMember) (types.Set, diag.Diagnostics) {
	memberType := buildMemberNestedObject().Type()

	if len(apiMembers) == 0 {
		return types.SetValueMust(memberType, []attr.Value{}), nil
	}

	elems := make([]attr.Value, len(apiMembers))
	for i, apiMember := range apiMembers {
		obj, diags := types.ObjectValue(
			map[string]attr.Type{RoleFieldMemberUserID: types.StringType},
			map[string]attr.Value{RoleFieldMemberUserID: types.StringValue(apiMember.UserID)},
		)
		if diags.HasError() {
			return types.SetNull(memberType), diags
		}
		elems[i] = obj
	}

	return types.SetValue(memberType, elems)
}

// MapStateToDataObject maps Terraform state/plan to API Role object
func (r *roleResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.Role, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	apiMembers, membersDiags := r.mapSetMembersToAPI(ctx, model.Members)
	diags.Append(membersDiags...)
	if diags.HasError() {
		return nil, diags
	}

	role := &api.Role{
		ID:          r.extractRoleID(model),
		Name:        model.Name.ValueString(),
		Members:     apiMembers,
		Permissions: model.Permissions,
	}

	return role, diags
}

// extractModelFromPlanOrState retrieves the RoleModel from plan or state
func (r *roleResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (RoleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model RoleModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractRoleID extracts the role ID from the model
func (r *roleResource) extractRoleID(model RoleModel) string {
	if model.ID.IsNull() {
		return ""
	}
	return model.ID.ValueString()
}

// mapSetMembersToAPI converts a types.Set of members to API members
func (r *roleResource) mapSetMembersToAPI(ctx context.Context, membersSet types.Set) ([]api.APIMember, diag.Diagnostics) {
	if membersSet.IsNull() || membersSet.IsUnknown() || len(membersSet.Elements()) == 0 {
		return []api.APIMember{}, nil
	}

	var modelMembers []RoleMemberModel
	diags := membersSet.ElementsAs(ctx, &modelMembers, false)
	if diags.HasError() {
		return nil, diags
	}

	apiMembers := make([]api.APIMember, len(modelMembers))
	for i, m := range modelMembers {
		apiMembers[i] = api.APIMember{UserID: m.UserID.ValueString()}
	}
	return apiMembers, nil
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *roleResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: resourcehandle.CreateStateUpgraderForVersion(0),
	}
}
