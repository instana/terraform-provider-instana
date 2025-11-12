package group

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewGroupResourceHandleFramework creates the resource handle for RBAC Groups
func NewGroupResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.Group] {
	return &groupResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName: ResourceInstanaGroupFramework,
			Schema: schema.Schema{
				Description: GroupDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: GroupDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: GroupDescName,
					},
					"permission_set": schema.SingleNestedAttribute{
						Description: GroupDescPermissionSet,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"application_ids": schema.SetAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetApplicationIDs,
								ElementType: types.StringType,
							},
							"infra_dfq_filter": schema.StringAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetInfraDFQFilter,
							},
							"kubernetes_cluster_uuids": schema.SetAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetKubernetesClusterUUIDs,
								ElementType: types.StringType,
							},
							"kubernetes_namespaces_uuids": schema.SetAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetKubernetesNamespaceUIDs,
								ElementType: types.StringType,
							},
							"mobile_app_ids": schema.SetAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetMobileAppIDs,
								ElementType: types.StringType,
							},
							"website_ids": schema.SetAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetWebsiteIDs,
								ElementType: types.StringType,
							},
							"permissions": schema.SetAttribute{
								Optional:    true,
								Description: GroupDescPermissionSetPermissions,
								ElementType: types.StringType,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(restapi.SupportedInstanaPermissions.ToStringSlice()...),
									),
								},
							},
						},
					},
					"member": schema.SetNestedAttribute{
						Description: GroupDescMembers,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"user_id": schema.StringAttribute{
									Required:    true,
									Description: GroupDescMemberUserID,
								},
								"email": schema.StringAttribute{
									Optional:    true,
									Description: GroupDescMemberEmail,
								},
							},
						},
					},
				},
			},
			SchemaVersion: 1,
		},
	}
}

type groupResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *groupResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *groupResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.Group] {
	return api.Groups()
}

func (r *groupResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *groupResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, group *restapi.Group) diag.Diagnostics {
	var diags diag.Diagnostics
	model := GroupModel{
		ID:   types.StringValue(group.ID),
		Name: types.StringValue(group.Name),
	}

	// Map members if present
	if len(group.Members) > 0 {
		model.Members = make([]GroupMemberModel, len(group.Members))
		for i, member := range group.Members {
			model.Members[i] = GroupMemberModel{
				UserID: types.StringValue(member.UserID),
				Email:  util.SetStringPointerToState(member.Email),
			}
		}
	}

	// Map permission set if not empty
	if !group.PermissionSet.IsEmpty() {
		permissionSetModel := &GroupPermissionSetModel{}

		// Map application IDs
		if len(group.PermissionSet.ApplicationIDs) > 0 {
			permissionSetModel.ApplicationIDs = make([]string, len(group.PermissionSet.ApplicationIDs))
			for i, appID := range group.PermissionSet.ApplicationIDs {
				permissionSetModel.ApplicationIDs[i] = appID.ScopeID
			}
		}

		// Map infra DFQ filter
		if group.PermissionSet.InfraDFQFilter != nil && len(group.PermissionSet.InfraDFQFilter.ScopeID) > 0 {
			permissionSetModel.InfraDFQFilter = types.StringValue(group.PermissionSet.InfraDFQFilter.ScopeID)
		} else {
			permissionSetModel.InfraDFQFilter = types.StringNull()
		}

		// Map Kubernetes cluster UUIDs
		if len(group.PermissionSet.KubernetesClusterUUIDs) > 0 {
			permissionSetModel.KubernetesClusterUUIDs = make([]string, len(group.PermissionSet.KubernetesClusterUUIDs))
			for i, kubeCluster := range group.PermissionSet.KubernetesClusterUUIDs {
				permissionSetModel.KubernetesClusterUUIDs[i] = kubeCluster.ScopeID
			}
		}

		// Map Kubernetes namespace UIDs
		if len(group.PermissionSet.KubernetesNamespaceUIDs) > 0 {
			permissionSetModel.KubernetesNamespaceUIDs = make([]string, len(group.PermissionSet.KubernetesNamespaceUIDs))
			for i, kubeNs := range group.PermissionSet.KubernetesNamespaceUIDs {
				permissionSetModel.KubernetesNamespaceUIDs[i] = kubeNs.ScopeID
			}
		}

		// Map mobile app IDs
		if len(group.PermissionSet.MobileAppIDs) > 0 {
			permissionSetModel.MobileAppIDs = make([]string, len(group.PermissionSet.MobileAppIDs))
			for i, mobileApp := range group.PermissionSet.MobileAppIDs {
				permissionSetModel.MobileAppIDs[i] = mobileApp.ScopeID
			}
		}

		// Map website IDs
		if len(group.PermissionSet.WebsiteIDs) > 0 {
			permissionSetModel.WebsiteIDs = make([]string, len(group.PermissionSet.WebsiteIDs))
			for i, website := range group.PermissionSet.WebsiteIDs {
				permissionSetModel.WebsiteIDs[i] = website.ScopeID
			}
		}

		// Map permissions
		if len(group.PermissionSet.Permissions) > 0 {
			permissionSetModel.Permissions = make([]string, len(group.PermissionSet.Permissions))
			for i, permission := range group.PermissionSet.Permissions {
				permissionSetModel.Permissions[i] = string(permission)
			}
		}

		model.PermissionSet = permissionSetModel
	}

	// Set the state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *groupResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.Group, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model GroupModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map ID
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	// Map members
	members := make([]restapi.APIMember, 0)
	if len(model.Members) > 0 {
		for _, memberModel := range model.Members {
			member := restapi.APIMember{
				UserID: memberModel.UserID.ValueString(),
			}

			if !memberModel.Email.IsNull() && !memberModel.Email.IsUnknown() {
				email := memberModel.Email.ValueString()
				member.Email = &email
			}

			members = append(members, member)
		}
	}

	// Map permission set
	permissionSet := restapi.APIPermissionSetWithRoles{}
	emptyScopeBinding := make([]restapi.ScopeBinding, 0)
	permissionSet.ApplicationIDs = emptyScopeBinding
	permissionSet.KubernetesNamespaceUIDs = emptyScopeBinding
	permissionSet.KubernetesClusterUUIDs = emptyScopeBinding
	permissionSet.WebsiteIDs = emptyScopeBinding
	permissionSet.MobileAppIDs = emptyScopeBinding
	permissionSet.Permissions = make([]restapi.InstanaPermission, 0)

	if model.PermissionSet != nil {
		// Map application IDs
		if len(model.PermissionSet.ApplicationIDs) > 0 {
			scopeBindings := make([]restapi.ScopeBinding, len(model.PermissionSet.ApplicationIDs))
			for i, appID := range model.PermissionSet.ApplicationIDs {
				scopeBindings[i] = restapi.ScopeBinding{ScopeID: appID}
			}
			permissionSet.ApplicationIDs = scopeBindings
		}

		// Map infra DFQ filter
		if !model.PermissionSet.InfraDFQFilter.IsNull() && !model.PermissionSet.InfraDFQFilter.IsUnknown() {
			permissionSet.InfraDFQFilter = &restapi.ScopeBinding{
				ScopeID: model.PermissionSet.InfraDFQFilter.ValueString(),
			}
		} else {
			roleId := "-1"
			permissionSet.InfraDFQFilter = &restapi.ScopeBinding{
				ScopeID:     "",
				ScopeRoleID: &roleId,
			}
		}

		// Map Kubernetes cluster UUIDs
		if len(model.PermissionSet.KubernetesClusterUUIDs) > 0 {
			scopeBindings := make([]restapi.ScopeBinding, len(model.PermissionSet.KubernetesClusterUUIDs))
			for i, kubeClusterUUID := range model.PermissionSet.KubernetesClusterUUIDs {
				scopeBindings[i] = restapi.ScopeBinding{ScopeID: kubeClusterUUID}
			}
			permissionSet.KubernetesClusterUUIDs = scopeBindings
		}

		// Map Kubernetes namespace UIDs
		if len(model.PermissionSet.KubernetesNamespaceUIDs) > 0 {
			scopeBindings := make([]restapi.ScopeBinding, len(model.PermissionSet.KubernetesNamespaceUIDs))
			for i, kubeNamespaceUID := range model.PermissionSet.KubernetesNamespaceUIDs {
				scopeBindings[i] = restapi.ScopeBinding{ScopeID: kubeNamespaceUID}
			}
			permissionSet.KubernetesNamespaceUIDs = scopeBindings
		}

		// Map mobile app IDs
		if len(model.PermissionSet.MobileAppIDs) > 0 {
			scopeBindings := make([]restapi.ScopeBinding, len(model.PermissionSet.MobileAppIDs))
			for i, mobileAppID := range model.PermissionSet.MobileAppIDs {
				scopeBindings[i] = restapi.ScopeBinding{ScopeID: mobileAppID}
			}
			permissionSet.MobileAppIDs = scopeBindings
		}

		// Map website IDs
		if len(model.PermissionSet.WebsiteIDs) > 0 {
			scopeBindings := make([]restapi.ScopeBinding, len(model.PermissionSet.WebsiteIDs))
			for i, websiteID := range model.PermissionSet.WebsiteIDs {
				scopeBindings[i] = restapi.ScopeBinding{ScopeID: websiteID}
			}
			permissionSet.WebsiteIDs = scopeBindings
		}

		// Map permissions
		if len(model.PermissionSet.Permissions) > 0 {
			permissions := make([]restapi.InstanaPermission, len(model.PermissionSet.Permissions))
			for i, permissionString := range model.PermissionSet.Permissions {
				permissions[i] = restapi.InstanaPermission(permissionString)
			}
			permissionSet.Permissions = permissions
		}
	}

	// Create the API object
	return &restapi.Group{
		ID:            id,
		Name:          model.Name.ValueString(),
		Members:       members,
		PermissionSet: permissionSet,
	}, diags
}
