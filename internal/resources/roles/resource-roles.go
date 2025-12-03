package roles

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// NewRoleResourceHandle creates the resource handle for RBAC Roles
func NewRoleResourceHandle() resourcehandle.ResourceHandle[*restapi.Role] {
	return &roleResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaRole,
			Schema:        buildRoleSchema(),
			SchemaVersion: 0,
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
				Description:  RoleDescMembers,
				Required:     true,
				NestedObject: buildMemberNestedObject(),
			},
			RoleFieldPermissions: schema.SetAttribute{
				Required:    true,
				Description: RoleDescPermissions,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(restapi.SupportedInstanaPermissions.ToStringSlice()...),
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
			RoleFieldMemberEmail: schema.StringAttribute{
				Optional:    true,
				Description: RoleDescMemberEmail,
			},
			RoleFieldMemberName: schema.StringAttribute{
				Optional:    true,
				Description: RoleDescMemberName,
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

func (r *roleResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.Role] {
	return api.Roles()
}

func (r *roleResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *roleResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, role *restapi.Role) diag.Diagnostics {
	// Get existing state/plan to preserve optional fields
	var existingModel RoleModel
	var diags diag.Diagnostics

	if plan != nil {
		diags.Append(plan.Get(ctx, &existingModel)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &existingModel)...)
	}

	if diags.HasError() {
		return diags
	}

	model := r.buildRoleModelFromAPIResponse(role, existingModel.Members)
	return state.Set(ctx, model)
}

// buildRoleModelFromAPIResponse constructs a RoleModel from the API Role response
func (r *roleResource) buildRoleModelFromAPIResponse(role *restapi.Role, existingMembers []RoleMemberModel) RoleModel {
	model := RoleModel{
		ID:          types.StringValue(role.ID),
		Name:        types.StringValue(role.Name),
		Members:     r.mapMembersToModel(role.Members, existingMembers),
		Permissions: role.Permissions,
	}

	return model
}

// mapMembersToModel converts API members to model members, preserving optional fields from existing state
func (r *roleResource) mapMembersToModel(apiMembers []restapi.APIMember, existingMembers []RoleMemberModel) []RoleMemberModel {
	if len(apiMembers) == 0 {
		return make([]RoleMemberModel, 0)
	}

	// Create a map of existing members by UserID for quick lookup
	existingMemberMap := make(map[string]RoleMemberModel)
	for _, existingMember := range existingMembers {
		if !existingMember.UserID.IsNull() {
			existingMemberMap[existingMember.UserID.ValueString()] = existingMember
		}
	}

	members := make([]RoleMemberModel, len(apiMembers))
	for i, apiMember := range apiMembers {
		member := RoleMemberModel{
			UserID: types.StringValue(apiMember.UserID),
		}

		// Check if we have existing data for this member
		if existingMember, exists := existingMemberMap[apiMember.UserID]; exists {
			// Preserve email from existing state if API doesn't return it or if it was set in plan
			if apiMember.Email != nil && *apiMember.Email != "" {
				member.Email = types.StringValue(*apiMember.Email)
			} else if !existingMember.Email.IsNull() {
				member.Email = existingMember.Email
			} else {
				member.Email = types.StringNull()
			}

			// Preserve name from existing state if API doesn't return it or if it was set in plan
			if apiMember.Name != nil && *apiMember.Name != "" {
				member.Name = types.StringValue(*apiMember.Name)
			} else if !existingMember.Name.IsNull() {
				member.Name = existingMember.Name
			} else {
				member.Name = types.StringNull()
			}
		} else {
			// New member, use API values
			member.Email = util.SetStringPointerToState(apiMember.Email)
			member.Name = util.SetStringPointerToState(apiMember.Name)
		}

		members[i] = member
	}
	return members
}

// MapStateToDataObject maps Terraform state/plan to API Role object
func (r *roleResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.Role, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	role := &restapi.Role{
		ID:          r.extractRoleID(model),
		Name:        model.Name.ValueString(),
		Members:     r.mapModelMembersToAPI(model.Members),
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

// mapModelMembersToAPI converts model members to API members
func (r *roleResource) mapModelMembersToAPI(modelMembers []RoleMemberModel) []restapi.APIMember {
	if len(modelMembers) == 0 {
		return make([]restapi.APIMember, 0)
	}

	apiMembers := make([]restapi.APIMember, 0, len(modelMembers))
	for _, memberModel := range modelMembers {
		apiMember := restapi.APIMember{
			UserID: memberModel.UserID.ValueString(),
		}

		if !memberModel.Email.IsNull() && !memberModel.Email.IsUnknown() {
			email := memberModel.Email.ValueString()
			apiMember.Email = &email
		}

		if !memberModel.Name.IsNull() && !memberModel.Name.IsUnknown() {
			name := memberModel.Name.ValueString()
			apiMember.Name = &name
		}

		apiMembers = append(apiMembers, apiMember)
	}

	return apiMembers
}

// Made with Bob

// GetStateUpgraders returns the state upgraders for this resource
func (r *roleResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return nil
}
