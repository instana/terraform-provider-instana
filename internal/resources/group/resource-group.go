package group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// NewGroupResourceHandle creates the resource handle for RBAC Groups
func NewGroupResourceHandle() resourcehandle.ResourceHandle[*restapi.Group] {
	return &groupResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaGroup,
			Schema:        buildGroupSchema(),
			SchemaVersion: 2,
		},
	}
}

// buildGroupSchema constructs the Terraform schema for the group resource
func buildGroupSchema() schema.Schema {
	return schema.Schema{
		Description: GroupDescResource,
		Attributes: map[string]schema.Attribute{
			GroupFieldID: schema.StringAttribute{
				Computed:    true,
				Description: GroupDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			GroupFieldName: schema.StringAttribute{
				Required:    true,
				Description: GroupDescName,
			},
			GroupFieldPermissionSet: schema.SingleNestedAttribute{
				Description: GroupDescPermissionSet,
				Optional:    true,
				Attributes:  buildPermissionSetAttributes(),
			},
			GroupFieldMembers: schema.SetNestedAttribute{
				Description:  GroupDescMembers,
				Optional:     true,
				NestedObject: buildMemberNestedObject(),
			},
		},
	}
}

// buildPermissionSetAttributes constructs the attributes for the permission set
func buildPermissionSetAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		GroupFieldPermissionSetApplicationIDs: schema.SetAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetApplicationIDs,
			ElementType: types.StringType,
		},
		GroupFieldPermissionSetInfraDFQFilter: schema.StringAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetInfraDFQFilter,
		},
		GroupFieldPermissionSetKubernetesClusterUUIDs: schema.SetAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetKubernetesClusterUUIDs,
			ElementType: types.StringType,
		},
		GroupFieldPermissionSetKubernetesNamespaceUIDs: schema.SetAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetKubernetesNamespaceUIDs,
			ElementType: types.StringType,
		},
		GroupFieldPermissionSetMobileAppIDs: schema.SetAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetMobileAppIDs,
			ElementType: types.StringType,
		},
		GroupFieldPermissionSetWebsiteIDs: schema.SetAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetWebsiteIDs,
			ElementType: types.StringType,
		},
		GroupFieldPermissionSetPermissions: schema.SetAttribute{
			Optional:    true,
			Description: GroupDescPermissionSetPermissions,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf(restapi.SupportedInstanaPermissions.ToStringSlice()...),
				),
			},
		},
	}
}

// buildMemberNestedObject constructs the nested object schema for group members
func buildMemberNestedObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			GroupFieldMemberUserID: schema.StringAttribute{
				Required:    true,
				Description: GroupDescMemberUserID,
			},
			GroupFieldMemberEmail: schema.StringAttribute{
				Optional:    true,
				Description: GroupDescMemberEmail,
			},
		},
	}
}

type groupResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *groupResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *groupResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.Group] {
	return api.Groups()
}

func (r *groupResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *groupResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, group *restapi.Group) diag.Diagnostics {
	var diags diag.Diagnostics
	var model GroupModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	} else {
		model = GroupModel{}
	}
	model = r.buildGroupModelFromAPIResponse(group, model)
	return state.Set(ctx, model)
}

// buildGroupModelFromAPIResponse constructs a GroupModel from the API Group response
func (r *groupResource) buildGroupModelFromAPIResponse(group *restapi.Group, groupModel GroupModel) GroupModel {
	model := GroupModel{
		ID:   types.StringValue(group.ID),
		Name: types.StringValue(group.Name),
	}

	if groupModel.Members == nil || len(groupModel.Members) == 0 {
		model.Members = r.mapMembersToModel(group.Members)
	} else {
		model.Members = groupModel.Members
	}

	if !group.PermissionSet.IsEmpty() {
		model.PermissionSet = r.mapPermissionSetToModel(&group.PermissionSet)
	}

	return model
}

// mapMembersToModel converts API members to model members
func (r *groupResource) mapMembersToModel(apiMembers []restapi.APIMember) []GroupMemberModel {
	if len(apiMembers) == 0 {
		return nil
	}
	members := make([]GroupMemberModel, len(apiMembers))
	for i, apiMember := range apiMembers {
		members[i] = GroupMemberModel{
			UserID: types.StringValue(apiMember.UserID),
			Email:  util.SetStringPointerToState(apiMember.Email),
		}
	}
	return members
}

// mapPermissionSetToModel converts API permission set to model permission set
func (r *groupResource) mapPermissionSetToModel(apiPermissionSet *restapi.APIPermissionSetWithRoles) *GroupPermissionSetModel {
	permissionSetModel := &GroupPermissionSetModel{
		InfraDFQFilter: r.mapInfraDFQFilterToModel(apiPermissionSet.InfraDFQFilter),
	}

	permissionSetModel.ApplicationIDs = r.extractScopeIDs(apiPermissionSet.ApplicationIDs)
	permissionSetModel.KubernetesClusterUUIDs = r.extractScopeIDs(apiPermissionSet.KubernetesClusterUUIDs)
	permissionSetModel.KubernetesNamespaceUIDs = r.extractScopeIDs(apiPermissionSet.KubernetesNamespaceUIDs)
	permissionSetModel.MobileAppIDs = r.extractScopeIDs(apiPermissionSet.MobileAppIDs)
	permissionSetModel.WebsiteIDs = r.extractScopeIDs(apiPermissionSet.WebsiteIDs)
	permissionSetModel.Permissions = r.convertPermissionsToStrings(apiPermissionSet.Permissions)

	return permissionSetModel
}

// mapInfraDFQFilterToModel converts the infra DFQ filter to a Terraform string type
func (r *groupResource) mapInfraDFQFilterToModel(infraFilter *restapi.ScopeBinding) types.String {
	if infraFilter != nil && len(infraFilter.ScopeID) > 0 {
		return types.StringValue(infraFilter.ScopeID)
	}
	return types.StringNull()
}

// extractScopeIDs extracts scope IDs from a slice of ScopeBindings
func (r *groupResource) extractScopeIDs(scopeBindings []restapi.ScopeBinding) []string {
	if len(scopeBindings) == 0 {
		return nil
	}

	scopeIDs := make([]string, len(scopeBindings))
	for i, binding := range scopeBindings {
		scopeIDs[i] = binding.ScopeID
	}
	return scopeIDs
}

// convertPermissionsToStrings converts InstanaPermission slice to string slice
func (r *groupResource) convertPermissionsToStrings(permissions []restapi.InstanaPermission) []string {
	if len(permissions) == 0 {
		return nil
	}

	permissionStrings := make([]string, len(permissions))
	for i, permission := range permissions {
		permissionStrings[i] = string(permission)
	}
	return permissionStrings
}

// MapStateToDataObject maps Terraform state/plan to API Group object
func (r *groupResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.Group, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	group := &restapi.Group{
		ID:            r.extractGroupID(model),
		Name:          model.Name.ValueString(),
		Members:       r.mapModelMembersToAPI(model.Members),
		PermissionSet: r.mapModelPermissionSetToAPI(model.PermissionSet),
	}

	return group, diags
}

// extractModelFromPlanOrState retrieves the GroupModel from plan or state
func (r *groupResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (GroupModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model GroupModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractGroupID extracts the group ID from the model
func (r *groupResource) extractGroupID(model GroupModel) string {
	if model.ID.IsNull() {
		return EmptyScopeID
	}
	return model.ID.ValueString()
}

// mapModelMembersToAPI converts model members to API members
func (r *groupResource) mapModelMembersToAPI(modelMembers []GroupMemberModel) []restapi.APIMember {
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

		apiMembers = append(apiMembers, apiMember)
	}

	return apiMembers
}

// mapModelPermissionSetToAPI converts model permission set to API permission set
func (r *groupResource) mapModelPermissionSetToAPI(modelPermissionSet *GroupPermissionSetModel) restapi.APIPermissionSetWithRoles {
	permissionSet := r.initializeEmptyPermissionSet()

	if modelPermissionSet == nil {
		return permissionSet
	}

	permissionSet.ApplicationIDs = r.createScopeBindings(modelPermissionSet.ApplicationIDs)
	permissionSet.InfraDFQFilter = r.mapModelInfraDFQFilterToAPI(modelPermissionSet.InfraDFQFilter)
	permissionSet.KubernetesClusterUUIDs = r.createScopeBindings(modelPermissionSet.KubernetesClusterUUIDs)
	permissionSet.KubernetesNamespaceUIDs = r.createScopeBindings(modelPermissionSet.KubernetesNamespaceUIDs)
	permissionSet.MobileAppIDs = r.createScopeBindings(modelPermissionSet.MobileAppIDs)
	permissionSet.WebsiteIDs = r.createScopeBindings(modelPermissionSet.WebsiteIDs)
	permissionSet.Permissions = r.convertStringsToPermissions(modelPermissionSet.Permissions)

	return permissionSet
}

// initializeEmptyPermissionSet creates an empty permission set with initialized slices
func (r *groupResource) initializeEmptyPermissionSet() restapi.APIPermissionSetWithRoles {
	emptyScopeBindings := make([]restapi.ScopeBinding, 0)
	return restapi.APIPermissionSetWithRoles{
		ApplicationIDs:          emptyScopeBindings,
		KubernetesNamespaceUIDs: emptyScopeBindings,
		KubernetesClusterUUIDs:  emptyScopeBindings,
		WebsiteIDs:              emptyScopeBindings,
		MobileAppIDs:            emptyScopeBindings,
		Permissions:             make([]restapi.InstanaPermission, 0),
	}
}

// createScopeBindings creates ScopeBinding slice from string slice
func (r *groupResource) createScopeBindings(scopeIDs []string) []restapi.ScopeBinding {
	if len(scopeIDs) == 0 {
		return make([]restapi.ScopeBinding, 0)
	}

	scopeBindings := make([]restapi.ScopeBinding, len(scopeIDs))
	for i, scopeID := range scopeIDs {
		scopeBindings[i] = restapi.ScopeBinding{ScopeID: scopeID}
	}
	return scopeBindings
}

// mapModelInfraDFQFilterToAPI converts model infra DFQ filter to API format
func (r *groupResource) mapModelInfraDFQFilterToAPI(infraFilter types.String) *restapi.ScopeBinding {
	if !infraFilter.IsNull() && !infraFilter.IsUnknown() {
		return &restapi.ScopeBinding{
			ScopeID: infraFilter.ValueString(),
		}
	}

	roleID := DefaultScopeRoleID
	return &restapi.ScopeBinding{
		ScopeID:     EmptyScopeID,
		ScopeRoleID: &roleID,
	}
}

// convertStringsToPermissions converts string slice to InstanaPermission slice
func (r *groupResource) convertStringsToPermissions(permissionStrings []string) []restapi.InstanaPermission {
	if len(permissionStrings) == 0 {
		return make([]restapi.InstanaPermission, 0)
	}

	permissions := make([]restapi.InstanaPermission, len(permissionStrings))
	for i, permissionString := range permissionStrings {
		permissions[i] = restapi.InstanaPermission(permissionString)
	}
	return permissions
}

// GetStateUpgraders returns the state upgraders for this resource
func (r *groupResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: resourcehandle.CreateStateUpgraderForVersion(1),
	}
}
