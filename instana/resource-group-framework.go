package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaGroupFramework the name of the terraform-provider-instana resource to manage groups for role based access control
const ResourceInstanaGroupFramework = "instana_rbac_group"

// GroupModel represents the data model for RBAC Group
type GroupModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Members       types.Set    `tfsdk:"member"`
	PermissionSet types.List   `tfsdk:"permission_set"`
}

// GroupMemberModel represents a member in the group
type GroupMemberModel struct {
	UserID types.String `tfsdk:"user_id"`
	Email  types.String `tfsdk:"email"`
}

// GroupPermissionSetModel represents the permission set for the group
type GroupPermissionSetModel struct {
	ApplicationIDs          types.Set    `tfsdk:"application_ids"`
	InfraDFQFilter          types.String `tfsdk:"infra_dfq_filter"`
	KubernetesClusterUUIDs  types.Set    `tfsdk:"kubernetes_cluster_uuids"`
	KubernetesNamespaceUIDs types.Set    `tfsdk:"kubernetes_namespaces_uuids"`
	MobileAppIDs            types.Set    `tfsdk:"mobile_app_ids"`
	WebsiteIDs              types.Set    `tfsdk:"website_ids"`
	Permissions             types.Set    `tfsdk:"permissions"`
}

// NewGroupResourceHandleFramework creates the resource handle for RBAC Groups
func NewGroupResourceHandleFramework() ResourceHandleFramework[*restapi.Group] {
	return &groupResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaGroupFramework,
			Schema: schema.Schema{
				Description: "This resource manages RBAC groups in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the group.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the Group",
					},
				},
				Blocks: map[string]schema.Block{
					"member": schema.SetNestedBlock{
						Description: "The members of the group",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"user_id": schema.StringAttribute{
									Required:    true,
									Description: "The user id of the group member",
								},
								"email": schema.StringAttribute{
									Optional:    true,
									Description: "The email address of the group member",
								},
							},
						},
					},
					"permission_set": schema.ListNestedBlock{
						Description: "The permission set of the group",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"application_ids": schema.SetAttribute{
									Optional:    true,
									Description: "The scope bindings to restrict access to applications",
									ElementType: types.StringType,
								},
								"infra_dfq_filter": schema.StringAttribute{
									Optional:    true,
									Description: "The scope binding for the dynamic filter query to restrict access to infrastructure assets",
								},
								"kubernetes_cluster_uuids": schema.SetAttribute{
									Optional:    true,
									Description: "The scope bindings to restrict access to Kubernetes Clusters",
									ElementType: types.StringType,
								},
								"kubernetes_namespaces_uuids": schema.SetAttribute{
									Optional:    true,
									Description: "The scope bindings to restrict access to Kubernetes namespaces",
									ElementType: types.StringType,
								},
								"mobile_app_ids": schema.SetAttribute{
									Optional:    true,
									Description: "The scope bindings to restrict access to mobile apps",
									ElementType: types.StringType,
								},
								"website_ids": schema.SetAttribute{
									Optional:    true,
									Description: "The scope bindings to restrict access to websites",
									ElementType: types.StringType,
								},
								"permissions": schema.SetAttribute{
									Optional:    true,
									Description: "The permissions assigned which should be assigned to the users of the group",
									ElementType: types.StringType,
									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf(restapi.SupportedInstanaPermissions.ToStringSlice()...),
										),
									},
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
	metaData ResourceMetaDataFramework
}

func (r *groupResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *groupResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.Group] {
	return api.Groups()
}

func (r *groupResourceFramework) SetComputedFields(ctx context.Context, plan *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *groupResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, group *restapi.Group) diag.Diagnostics {
	var diags diag.Diagnostics
	model := GroupModel{
		ID:   types.StringValue(group.ID),
		Name: types.StringValue(group.Name),
	}

	// Map members if present
	if len(group.Members) > 0 {
		memberElements := make([]attr.Value, len(group.Members))
		for i, member := range group.Members {
			memberModel := GroupMemberModel{
				UserID: types.StringValue(member.UserID),
			}

			if member.Email != nil {
				memberModel.Email = types.StringValue(*member.Email)
			} else {
				memberModel.Email = types.StringNull()
			}

			memberObj, memberDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"user_id": types.StringType,
				"email":   types.StringType,
			}, memberModel)

			if memberDiags.HasError() {
				diags.Append(memberDiags...)
				return diags
			}

			memberElements[i] = memberObj
		}

		model.Members = types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"user_id": types.StringType,
				"email":   types.StringType,
			},
		}, memberElements)
	} else {
		model.Members = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"user_id": types.StringType,
				"email":   types.StringType,
			},
		})
	}

	// Map permission set if not empty
	if !group.PermissionSet.IsEmpty() {
		permissionSetModel := GroupPermissionSetModel{}

		// Map application IDs
		if len(group.PermissionSet.ApplicationIDs) > 0 {
			appIDElements := make([]attr.Value, len(group.PermissionSet.ApplicationIDs))
			for i, appID := range group.PermissionSet.ApplicationIDs {
				appIDElements[i] = types.StringValue(appID.ScopeID)
			}
			permissionSetModel.ApplicationIDs = types.SetValueMust(types.StringType, appIDElements)
		} else {
			permissionSetModel.ApplicationIDs = types.SetNull(types.StringType)
		}

		// Map infra DFQ filter
		if group.PermissionSet.InfraDFQFilter != nil && len(group.PermissionSet.InfraDFQFilter.ScopeID) > 0 {
			permissionSetModel.InfraDFQFilter = types.StringValue(group.PermissionSet.InfraDFQFilter.ScopeID)
		} else {
			permissionSetModel.InfraDFQFilter = types.StringNull()
		}

		// Map Kubernetes cluster UUIDs
		if len(group.PermissionSet.KubernetesClusterUUIDs) > 0 {
			kubeClusterElements := make([]attr.Value, len(group.PermissionSet.KubernetesClusterUUIDs))
			for i, kubeCluster := range group.PermissionSet.KubernetesClusterUUIDs {
				kubeClusterElements[i] = types.StringValue(kubeCluster.ScopeID)
			}
			permissionSetModel.KubernetesClusterUUIDs = types.SetValueMust(types.StringType, kubeClusterElements)
		} else {
			permissionSetModel.KubernetesClusterUUIDs = types.SetNull(types.StringType)
		}

		// Map Kubernetes namespace UIDs
		if len(group.PermissionSet.KubernetesNamespaceUIDs) > 0 {
			kubeNsElements := make([]attr.Value, len(group.PermissionSet.KubernetesNamespaceUIDs))
			for i, kubeNs := range group.PermissionSet.KubernetesNamespaceUIDs {
				kubeNsElements[i] = types.StringValue(kubeNs.ScopeID)
			}
			permissionSetModel.KubernetesNamespaceUIDs = types.SetValueMust(types.StringType, kubeNsElements)
		} else {
			permissionSetModel.KubernetesNamespaceUIDs = types.SetNull(types.StringType)
		}

		// Map mobile app IDs
		if len(group.PermissionSet.MobileAppIDs) > 0 {
			mobileAppElements := make([]attr.Value, len(group.PermissionSet.MobileAppIDs))
			for i, mobileApp := range group.PermissionSet.MobileAppIDs {
				mobileAppElements[i] = types.StringValue(mobileApp.ScopeID)
			}
			permissionSetModel.MobileAppIDs = types.SetValueMust(types.StringType, mobileAppElements)
		} else {
			permissionSetModel.MobileAppIDs = types.SetNull(types.StringType)
		}

		// Map website IDs
		if len(group.PermissionSet.WebsiteIDs) > 0 {
			websiteElements := make([]attr.Value, len(group.PermissionSet.WebsiteIDs))
			for i, website := range group.PermissionSet.WebsiteIDs {
				websiteElements[i] = types.StringValue(website.ScopeID)
			}
			permissionSetModel.WebsiteIDs = types.SetValueMust(types.StringType, websiteElements)
		} else {
			permissionSetModel.WebsiteIDs = types.SetNull(types.StringType)
		}

		// Map permissions
		if len(group.PermissionSet.Permissions) > 0 {
			permissionElements := make([]attr.Value, len(group.PermissionSet.Permissions))
			for i, permission := range group.PermissionSet.Permissions {
				permissionElements[i] = types.StringValue(string(permission))
			}
			permissionSetModel.Permissions = types.SetValueMust(types.StringType, permissionElements)
		} else {
			permissionSetModel.Permissions = types.SetNull(types.StringType)
		}

		// Convert permission set model to object
		permissionSetObj, permDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"application_ids":             types.SetType{ElemType: types.StringType},
			"infra_dfq_filter":            types.StringType,
			"kubernetes_cluster_uuids":    types.SetType{ElemType: types.StringType},
			"kubernetes_namespaces_uuids": types.SetType{ElemType: types.StringType},
			"mobile_app_ids":              types.SetType{ElemType: types.StringType},
			"website_ids":                 types.SetType{ElemType: types.StringType},
			"permissions":                 types.SetType{ElemType: types.StringType},
		}, permissionSetModel)

		if permDiags.HasError() {
			diags.Append(permDiags...)
			return diags
		}

		model.PermissionSet = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"application_ids":             types.SetType{ElemType: types.StringType},
				"infra_dfq_filter":            types.StringType,
				"kubernetes_cluster_uuids":    types.SetType{ElemType: types.StringType},
				"kubernetes_namespaces_uuids": types.SetType{ElemType: types.StringType},
				"mobile_app_ids":              types.SetType{ElemType: types.StringType},
				"website_ids":                 types.SetType{ElemType: types.StringType},
				"permissions":                 types.SetType{ElemType: types.StringType},
			},
		}, []attr.Value{permissionSetObj})
	} else {
		model.PermissionSet = types.ListNull(types.ObjectType{})
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
	var members []restapi.APIMember
	if !model.Members.IsNull() && !model.Members.IsUnknown() {
		var memberModels []GroupMemberModel
		diags.Append(model.Members.ElementsAs(ctx, &memberModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		members = make([]restapi.APIMember, len(memberModels))
		for i, memberModel := range memberModels {
			member := restapi.APIMember{
				UserID: memberModel.UserID.ValueString(),
			}

			if !memberModel.Email.IsNull() && !memberModel.Email.IsUnknown() {
				email := memberModel.Email.ValueString()
				member.Email = &email
			}

			members[i] = member
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

	if !model.PermissionSet.IsNull() && !model.PermissionSet.IsUnknown() {
		var permissionSetModels []GroupPermissionSetModel
		diags.Append(model.PermissionSet.ElementsAs(ctx, &permissionSetModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if len(permissionSetModels) > 0 {
			permissionSetModel := permissionSetModels[0]

			// Map application IDs
			if !permissionSetModel.ApplicationIDs.IsNull() && !permissionSetModel.ApplicationIDs.IsUnknown() {
				var appIDs []string
				diags.Append(permissionSetModel.ApplicationIDs.ElementsAs(ctx, &appIDs, false)...)
				if diags.HasError() {
					return nil, diags
				}

				scopeBindings := make([]restapi.ScopeBinding, len(appIDs))
				for i, appID := range appIDs {
					scopeBindings[i] = restapi.ScopeBinding{ScopeID: appID}
				}
				permissionSet.ApplicationIDs = scopeBindings
			}

			// Map infra DFQ filter
			if !permissionSetModel.InfraDFQFilter.IsNull() && !permissionSetModel.InfraDFQFilter.IsUnknown() {
				permissionSet.InfraDFQFilter = &restapi.ScopeBinding{
					ScopeID: permissionSetModel.InfraDFQFilter.ValueString(),
				}
			}

			// Map Kubernetes cluster UUIDs
			if !permissionSetModel.KubernetesClusterUUIDs.IsNull() && !permissionSetModel.KubernetesClusterUUIDs.IsUnknown() {
				var kubeClusterUUIDs []string
				diags.Append(permissionSetModel.KubernetesClusterUUIDs.ElementsAs(ctx, &kubeClusterUUIDs, false)...)
				if diags.HasError() {
					return nil, diags
				}

				scopeBindings := make([]restapi.ScopeBinding, len(kubeClusterUUIDs))
				for i, kubeClusterUUID := range kubeClusterUUIDs {
					scopeBindings[i] = restapi.ScopeBinding{ScopeID: kubeClusterUUID}
				}
				permissionSet.KubernetesClusterUUIDs = scopeBindings
			}

			// Map Kubernetes namespace UIDs
			if !permissionSetModel.KubernetesNamespaceUIDs.IsNull() && !permissionSetModel.KubernetesNamespaceUIDs.IsUnknown() {
				var kubeNamespaceUIDs []string
				diags.Append(permissionSetModel.KubernetesNamespaceUIDs.ElementsAs(ctx, &kubeNamespaceUIDs, false)...)
				if diags.HasError() {
					return nil, diags
				}

				scopeBindings := make([]restapi.ScopeBinding, len(kubeNamespaceUIDs))
				for i, kubeNamespaceUID := range kubeNamespaceUIDs {
					scopeBindings[i] = restapi.ScopeBinding{ScopeID: kubeNamespaceUID}
				}
				permissionSet.KubernetesNamespaceUIDs = scopeBindings
			}

			// Map mobile app IDs
			if !permissionSetModel.MobileAppIDs.IsNull() && !permissionSetModel.MobileAppIDs.IsUnknown() {
				var mobileAppIDs []string
				diags.Append(permissionSetModel.MobileAppIDs.ElementsAs(ctx, &mobileAppIDs, false)...)
				if diags.HasError() {
					return nil, diags
				}

				scopeBindings := make([]restapi.ScopeBinding, len(mobileAppIDs))
				for i, mobileAppID := range mobileAppIDs {
					scopeBindings[i] = restapi.ScopeBinding{ScopeID: mobileAppID}
				}
				permissionSet.MobileAppIDs = scopeBindings
			}

			// Map website IDs
			if !permissionSetModel.WebsiteIDs.IsNull() && !permissionSetModel.WebsiteIDs.IsUnknown() {
				var websiteIDs []string
				diags.Append(permissionSetModel.WebsiteIDs.ElementsAs(ctx, &websiteIDs, false)...)
				if diags.HasError() {
					return nil, diags
				}

				scopeBindings := make([]restapi.ScopeBinding, len(websiteIDs))
				for i, websiteID := range websiteIDs {
					scopeBindings[i] = restapi.ScopeBinding{ScopeID: websiteID}
				}
				permissionSet.WebsiteIDs = scopeBindings
			}

			// Map permissions
			if !permissionSetModel.Permissions.IsNull() && !permissionSetModel.Permissions.IsUnknown() {
				var permissionStrings []string
				diags.Append(permissionSetModel.Permissions.ElementsAs(ctx, &permissionStrings, false)...)
				if diags.HasError() {
					return nil, diags
				}

				permissions := make([]restapi.InstanaPermission, len(permissionStrings))
				for i, permissionString := range permissionStrings {
					permissions[i] = restapi.InstanaPermission(permissionString)
				}
				permissionSet.Permissions = permissions
			}
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

// Made with Bob

// Made with Bob
