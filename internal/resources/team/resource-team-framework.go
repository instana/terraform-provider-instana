package team

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewTeamResourceHandleFramework creates the resource handle for RBAC Teams
func NewTeamResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.Team] {
	return &teamResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaTeamFramework,
			Schema:        buildTeamSchema(),
			SchemaVersion: 1,
		},
	}
}

// buildTeamSchema constructs the Terraform schema for the team resource
func buildTeamSchema() schema.Schema {
	return schema.Schema{
		Description: TeamDescResource,
		Attributes: map[string]schema.Attribute{
			TeamFieldID: schema.StringAttribute{
				Computed:    true,
				Description: TeamDescID,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			TeamFieldTag: schema.StringAttribute{
				Required:    true,
				Description: TeamDescTag,
			},
			TeamFieldInfo: schema.SingleNestedAttribute{
				Description: TeamDescInfo,
				Optional:    true,
				Attributes:  buildTeamInfoAttributes(),
			},
			TeamFieldMembers: schema.SetNestedAttribute{
				Description:  TeamDescMembers,
				Optional:     true,
				NestedObject: buildMemberNestedObject(),
			},
			TeamFieldScope: schema.SingleNestedAttribute{
				Description: TeamDescScope,
				Optional:    true,
				Attributes:  buildScopeAttributes(),
			},
		},
	}
}

// buildTeamInfoAttributes constructs the attributes for team info
func buildTeamInfoAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		TeamFieldInfoDescription: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescInfoDescription,
		},
	}
}

// buildMemberNestedObject constructs the nested object schema for team members
func buildMemberNestedObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			TeamFieldMemberUserID: schema.StringAttribute{
				Required:    true,
				Description: TeamDescMemberUserID,
			},
			TeamFieldMemberEmail: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: TeamDescMemberEmail,
			},
			TeamFieldMemberName: schema.StringAttribute{
				Optional:    true,
				Description: TeamDescMemberName,
			},
			TeamFieldMemberRoles: schema.SetNestedAttribute{
				Description:  TeamDescMemberRoles,
				Optional:     true,
				NestedObject: buildMemberRoleNestedObject(),
			},
		},
	}
}

// buildMemberRoleNestedObject constructs the nested object schema for member roles
func buildMemberRoleNestedObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			TeamFieldMemberRoleID: schema.StringAttribute{
				Required:    true,
				Description: TeamDescMemberRoleID,
			},
			TeamFieldMemberRoleName: schema.StringAttribute{
				Optional:    true,
				Description: TeamDescMemberRoleName,
			},
			TeamFieldMemberRoleViaIdP: schema.BoolAttribute{
				Optional:    true,
				Description: TeamDescMemberRoleViaIdP,
			},
		},
	}
}

// buildScopeAttributes constructs the attributes for team scope
func buildScopeAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		TeamFieldScopeAccessPermissions: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeAccessPermissions,
			ElementType: types.StringType,
		},
		TeamFieldScopeApplications: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeApplications,
			ElementType: types.StringType,
		},
		TeamFieldScopeKubernetesClusters: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeKubernetesClusters,
			ElementType: types.StringType,
		},
		TeamFieldScopeKubernetesNamespaces: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeKubernetesNamespaces,
			ElementType: types.StringType,
		},
		TeamFieldScopeMobileApps: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeMobileApps,
			ElementType: types.StringType,
		},
		TeamFieldScopeWebsites: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeWebsites,
			ElementType: types.StringType,
		},
		TeamFieldScopeInfraDFQFilter: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeInfraDFQFilter,
		},
		TeamFieldScopeActionFilter: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeActionFilter,
		},
		TeamFieldScopeLogFilter: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeLogFilter,
		},
		TeamFieldScopeBusinessPerspectives: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeBusinessPerspectives,
			ElementType: types.StringType,
		},
		TeamFieldScopeSloIDs: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeSloIDs,
			ElementType: types.StringType,
		},
		TeamFieldScopeSyntheticTests: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeSyntheticTests,
			ElementType: types.StringType,
		},
		TeamFieldScopeSyntheticCredentials: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeSyntheticCredentials,
			ElementType: types.StringType,
		},
		TeamFieldScopeTagIDs: schema.SetAttribute{
			Optional:    true,
			Description: TeamDescScopeTagIDs,
			ElementType: types.StringType,
		},
		TeamFieldScopeRestrictedApplicationFilter: schema.SingleNestedAttribute{
			Description: TeamDescScopeRestrictedApplicationFilter,
			Optional:    true,
			Attributes:  buildRestrictedApplicationFilterAttributes(),
		},
	}
}

// buildRestrictedApplicationFilterAttributes constructs the attributes for restricted application filter
func buildRestrictedApplicationFilterAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		TeamFieldScopeRestrictedApplicationFilterLabel: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeRestrictedApplicationFilterLabel,
		},
		TeamFieldScopeRestrictedApplicationFilterRestrictingApplicationID: schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeRestrictedApplicationFilterRestrictingApplicationID,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		TeamFieldScopeRestrictedApplicationFilterScope: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeRestrictedApplicationFilterScope,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(restapi.RestrictedApplicationFilterScopeIncludeNoDownstream),
					string(restapi.RestrictedApplicationFilterScopeIncludeImmediateDownstream),
					string(restapi.RestrictedApplicationFilterScopeIncludeAllDownstream),
				),
			},
		},
		TeamFieldScopeRestrictedApplicationFilterTagFilterExpression: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeRestrictedApplicationFilterTagFilterExpression,
		},
	}
}

type teamResourceFramework struct {
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *teamResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
	return &r.metaData
}

func (r *teamResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.Team] {
	return api.Teams()
}

func (r *teamResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *teamResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, _ *tfsdk.Plan, team *restapi.Team) diag.Diagnostics {
	model, diags := r.buildTeamModelFromAPIResponse(team)
	if diags.HasError() {
		return diags
	}
	return state.Set(ctx, model)
}

// buildTeamModelFromAPIResponse constructs a TeamModel from the API Team response
func (r *teamResourceFramework) buildTeamModelFromAPIResponse(team *restapi.Team) (TeamModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := TeamModel{
		ID:  types.StringValue(team.ID),
		Tag: types.StringValue(team.Tag),
	}

	if team.Info != nil {
		model.Info = r.mapTeamInfoToModel(team.Info)
	}

	if len(team.Members) > 0 {
		model.Members = r.mapMembersToModel(team.Members)
	}

	if team.Scope != nil {
		scopeModel, scopeDiags := r.mapScopeToModel(team.Scope)
		diags.Append(scopeDiags...)
		if !diags.HasError() {
			model.Scope = scopeModel
		}
	}

	return model, diags
}

// mapTeamInfoToModel converts API team info to model team info
func (r *teamResourceFramework) mapTeamInfoToModel(apiInfo *restapi.TeamInfo) *TeamInfoModel {
	return &TeamInfoModel{
		Description: util.SetStringPointerToState(apiInfo.Description),
	}
}

// mapMembersToModel converts API members to model members
func (r *teamResourceFramework) mapMembersToModel(apiMembers []restapi.TeamMember) []TeamMemberModel {
	members := make([]TeamMemberModel, len(apiMembers))
	for i, apiMember := range apiMembers {
		members[i] = TeamMemberModel{
			UserID: types.StringValue(apiMember.UserID),
			Email:  util.SetStringPointerToState(apiMember.Email),
			Name:   util.SetStringPointerToState(apiMember.Name),
			Roles:  r.mapRolesToModel(apiMember.Roles),
		}
	}
	return members
}

// mapRolesToModel converts API roles to model roles
func (r *teamResourceFramework) mapRolesToModel(apiRoles []restapi.TeamRole) []TeamMemberRole {
	if len(apiRoles) == 0 {
		return nil
	}

	roles := make([]TeamMemberRole, len(apiRoles))
	for i, apiRole := range apiRoles {
		roles[i] = TeamMemberRole{
			RoleID:   types.StringValue(apiRole.RoleID),
			RoleName: util.SetStringPointerToState(apiRole.RoleName),
			ViaIdP:   util.SetBoolPointerToState(apiRole.ViaIdP),
		}
	}
	return roles
}

// mapScopeToModel converts API scope to model scope
func (r *teamResourceFramework) mapScopeToModel(apiScope *restapi.TeamScope) (*TeamScopeModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopeModel := &TeamScopeModel{
		AccessPermissions:    apiScope.AccessPermissions,
		Applications:         apiScope.Applications,
		KubernetesClusters:   apiScope.KubernetesClusters,
		KubernetesNamespaces: apiScope.KubernetesNamespaces,
		MobileApps:           apiScope.MobileApps,
		Websites:             apiScope.Websites,
		InfraDFQFilter:       util.SetStringPointerToState(apiScope.InfraDFQFilter),
		ActionFilter:         util.SetStringPointerToState(apiScope.ActionFilter),
		LogFilter:            util.SetStringPointerToState(apiScope.LogFilter),
		BusinessPerspectives: apiScope.BusinessPerspectives,
		SloIDs:               apiScope.SloIDs,
		SyntheticTests:       apiScope.SyntheticTests,
		SyntheticCredentials: apiScope.SyntheticCredentials,
		TagIDs:               apiScope.TagIDs,
	}

	if apiScope.RestrictedApplicationFilter != nil {
		filterModel, filterDiags := r.mapRestrictedApplicationFilterToModel(apiScope.RestrictedApplicationFilter)
		diags.Append(filterDiags...)
		if !diags.HasError() {
			scopeModel.RestrictedApplicationFilter = filterModel
		}
	}

	return scopeModel, diags
}

// mapRestrictedApplicationFilterToModel converts API restricted application filter to model
func (r *teamResourceFramework) mapRestrictedApplicationFilterToModel(apiFilter *restapi.RestrictedApplicationFilter) (*TeamRestrictedApplicationFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	filterModel := &TeamRestrictedApplicationFilterModel{
		Label:                    util.SetStringPointerToState(apiFilter.Label),
		RestrictingApplicationID: util.SetStringPointerToState(apiFilter.RestrictingApplicationID),
	}

	if apiFilter.Scope != nil {
		filterModel.Scope = types.StringValue(string(*apiFilter.Scope))
	} else {
		filterModel.Scope = types.StringNull()
	}

	if apiFilter.TagFilterExpression != nil {
		tagFilterString, err := tagfilter.MapTagFilterToNormalizedString(apiFilter.TagFilterExpression)
		if err != nil {
			diags.AddError("Failed to map tag filter expression", err.Error())
			return nil, diags
		}
		filterModel.TagFilterExpression = util.SetStringPointerToState(tagFilterString)
	} else {
		filterModel.TagFilterExpression = types.StringNull()
	}

	return filterModel, diags
}

// MapStateToDataObject maps Terraform state/plan to API Team object
func (r *teamResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	team := &restapi.Team{
		ID:  r.extractTeamID(model),
		Tag: model.Tag.ValueString(),
	}

	if model.Info != nil {
		team.Info = r.mapModelInfoToAPI(model.Info)
	}

	if len(model.Members) > 0 {
		team.Members = r.mapModelMembersToAPI(model.Members)
	}

	if model.Scope != nil {
		scopeAPI, scopeDiags := r.mapModelScopeToAPI(model.Scope)
		diags.Append(scopeDiags...)
		if !diags.HasError() {
			team.Scope = scopeAPI
		}
	}

	return team, diags
}

// extractModelFromPlanOrState retrieves the TeamModel from plan or state
func (r *teamResourceFramework) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (TeamModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model TeamModel

	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	return model, diags
}

// extractTeamID extracts the team ID from the model
func (r *teamResourceFramework) extractTeamID(model TeamModel) string {
	if model.ID.IsNull() {
		return ""
	}
	return model.ID.ValueString()
}

// mapModelInfoToAPI converts model team info to API team info
func (r *teamResourceFramework) mapModelInfoToAPI(modelInfo *TeamInfoModel) *restapi.TeamInfo {
	if modelInfo.Description.IsNull() || modelInfo.Description.IsUnknown() {
		return nil
	}

	desc := modelInfo.Description.ValueString()
	return &restapi.TeamInfo{
		Description: &desc,
	}
}

// mapModelMembersToAPI converts model members to API members
func (r *teamResourceFramework) mapModelMembersToAPI(modelMembers []TeamMemberModel) []restapi.TeamMember {
	if len(modelMembers) == 0 {
		return nil
	}

	apiMembers := make([]restapi.TeamMember, 0, len(modelMembers))
	for _, memberModel := range modelMembers {
		apiMember := restapi.TeamMember{
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

		if len(memberModel.Roles) > 0 {
			apiMember.Roles = r.mapModelRolesToAPI(memberModel.Roles)
		}

		apiMembers = append(apiMembers, apiMember)
	}

	return apiMembers
}

// mapModelRolesToAPI converts model roles to API roles
func (r *teamResourceFramework) mapModelRolesToAPI(modelRoles []TeamMemberRole) []restapi.TeamRole {
	if len(modelRoles) == 0 {
		return nil
	}

	apiRoles := make([]restapi.TeamRole, len(modelRoles))
	for i, roleModel := range modelRoles {
		apiRole := restapi.TeamRole{
			RoleID: roleModel.RoleID.ValueString(),
		}

		if !roleModel.RoleName.IsNull() && !roleModel.RoleName.IsUnknown() {
			roleName := roleModel.RoleName.ValueString()
			apiRole.RoleName = &roleName
		}

		if !roleModel.ViaIdP.IsNull() && !roleModel.ViaIdP.IsUnknown() {
			viaIdP := roleModel.ViaIdP.ValueBool()
			apiRole.ViaIdP = &viaIdP
		}

		apiRoles[i] = apiRole
	}

	return apiRoles
}

// mapModelScopeToAPI converts model scope to API scope
func (r *teamResourceFramework) mapModelScopeToAPI(modelScope *TeamScopeModel) (*restapi.TeamScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiScope := &restapi.TeamScope{
		AccessPermissions:    modelScope.AccessPermissions,
		Applications:         modelScope.Applications,
		KubernetesClusters:   modelScope.KubernetesClusters,
		KubernetesNamespaces: modelScope.KubernetesNamespaces,
		MobileApps:           modelScope.MobileApps,
		Websites:             modelScope.Websites,
		BusinessPerspectives: modelScope.BusinessPerspectives,
		SloIDs:               modelScope.SloIDs,
		SyntheticTests:       modelScope.SyntheticTests,
		SyntheticCredentials: modelScope.SyntheticCredentials,
		TagIDs:               modelScope.TagIDs,
	}

	if !modelScope.InfraDFQFilter.IsNull() && !modelScope.InfraDFQFilter.IsUnknown() {
		filter := modelScope.InfraDFQFilter.ValueString()
		apiScope.InfraDFQFilter = &filter
	}

	if !modelScope.ActionFilter.IsNull() && !modelScope.ActionFilter.IsUnknown() {
		filter := modelScope.ActionFilter.ValueString()
		apiScope.ActionFilter = &filter
	}

	if !modelScope.LogFilter.IsNull() && !modelScope.LogFilter.IsUnknown() {
		filter := modelScope.LogFilter.ValueString()
		apiScope.LogFilter = &filter
	}

	if modelScope.RestrictedApplicationFilter != nil {
		filterAPI, filterDiags := r.mapModelRestrictedApplicationFilterToAPI(modelScope.RestrictedApplicationFilter)
		diags.Append(filterDiags...)
		if !diags.HasError() {
			apiScope.RestrictedApplicationFilter = filterAPI
		}
	}

	return apiScope, diags
}

// mapModelRestrictedApplicationFilterToAPI converts model restricted application filter to API
func (r *teamResourceFramework) mapModelRestrictedApplicationFilterToAPI(modelFilter *TeamRestrictedApplicationFilterModel) (*restapi.RestrictedApplicationFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiFilter := &restapi.RestrictedApplicationFilter{}

	if !modelFilter.Label.IsNull() && !modelFilter.Label.IsUnknown() {
		label := modelFilter.Label.ValueString()
		apiFilter.Label = &label
	}

	if !modelFilter.RestrictingApplicationID.IsNull() && !modelFilter.RestrictingApplicationID.IsUnknown() {
		appID := modelFilter.RestrictingApplicationID.ValueString()
		apiFilter.RestrictingApplicationID = &appID
	}

	if !modelFilter.Scope.IsNull() && !modelFilter.Scope.IsUnknown() {
		scope := restapi.RestrictedApplicationFilterScope(modelFilter.Scope.ValueString())
		apiFilter.Scope = &scope
	}

	if !modelFilter.TagFilterExpression.IsNull() && !modelFilter.TagFilterExpression.IsUnknown() {
		tagFilterString := modelFilter.TagFilterExpression.ValueString()
		mapper := tagfilter.NewMapper()
		parser := tagfilter.NewParser()

		expr, err := parser.Parse(tagFilterString)
		if err != nil {
			diags.AddError("Failed to parse tag filter expression", err.Error())
			return nil, diags
		}

		apiFilter.TagFilterExpression = mapper.ToAPIModel(expr)
	}

	return apiFilter, diags
}

// Made with Bob
