package team

import (
	"context"

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
	"github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/instana/terraform-provider-instana/internal/util"
)

// stringsToSet converts a []string to a types.Set of strings.
// Both nil (JSON null) and [] (JSON []) are stored as an empty set — the Teams API
// uses null and [] interchangeably to mean "not configured", and the schema fields are
// Optional+Computed with UseStateForUnknown. Storing an empty set for both cases keeps
// the plan value ([] from TF filling an unset nested attr) consistent with state after
// apply, regardless of whether the API echoes back null or [].
func stringsToSet(ctx context.Context, ss []string) (types.Set, diag.Diagnostics) {
	if len(ss) == 0 {
		return types.SetValueMust(types.StringType, []attr.Value{}), nil
	}
	return types.SetValueFrom(ctx, types.StringType, ss)
}

// setToStrings extracts a []string from a types.Set. Returns nil for null/unknown/empty.
func setToStrings(ctx context.Context, s types.Set) ([]string, diag.Diagnostics) {
	if s.IsNull() || s.IsUnknown() || len(s.Elements()) == 0 {
		return nil, nil
	}
	var result []string
	return result, s.ElementsAs(ctx, &result, false)
}

// NewTeamResourceHandle creates the resource handle for RBAC Teams
func NewTeamResourceHandle() resourcehandle.ResourceHandle[*api.Team] {
	return &teamResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaTeam,
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
		},
	}
}

// buildScopeAttributes constructs the attributes for team scope
func buildScopeAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		TeamFieldScopeAccessPermissions: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeAccessPermissions,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeApplications: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeApplications,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeKubernetesClusters: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeKubernetesClusters,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeKubernetesNamespaces: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeKubernetesNamespaces,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeMobileApps: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeMobileApps,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeWebsites: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeWebsites,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
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
			Computed:    true,
			Description: TeamDescScopeBusinessPerspectives,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeSloIDs: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeSloIDs,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeSyntheticTests: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeSyntheticTests,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeSyntheticCredentials: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeSyntheticCredentials,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
		},
		TeamFieldScopeTagIDs: schema.SetAttribute{
			Optional:    true,
			Computed:    true,
			Description: TeamDescScopeTagIDs,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
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
		TeamFieldScopeRestrictedApplicationFilterScope: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeRestrictedApplicationFilterScope,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(api.RestrictedApplicationFilterScopeIncludeNoDownstream),
					string(api.RestrictedApplicationFilterScopeIncludeImmediateDownstream),
					string(api.RestrictedApplicationFilterScopeIncludeAllDownstream),
				),
			},
		},
		TeamFieldScopeRestrictedApplicationFilterTagFilterExpression: schema.StringAttribute{
			Optional:    true,
			Description: TeamDescScopeRestrictedApplicationFilterTagFilterExpression,
		},
	}
}

type teamResource struct {
	metaData resourcehandle.ResourceMetaData
}

func (r *teamResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

func (r *teamResource) GetRestResource(api client.InstanaAPI) rest.RestResource[*api.Team] {
	return api.Teams()
}

func (r *teamResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// UpdateState updates the Terraform state with data from the API response
func (r *teamResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, team *api.Team) diag.Diagnostics {
	var diags diag.Diagnostics
	var model TeamModel
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}
	m, diags := r.buildTeamModelFromAPIResponse(ctx, team, model)
	if diags.HasError() {
		return diags
	}
	return state.Set(ctx, m)
}

// buildTeamModelFromAPIResponse constructs a TeamModel from the API Team response
func (r *teamResource) buildTeamModelFromAPIResponse(ctx context.Context, team *api.Team, planModel TeamModel) (TeamModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := TeamModel{
		ID:  types.StringValue(team.ID),
		Tag: types.StringValue(team.Tag),
	}

	if planModel.Info != nil {
		model.Info = r.mapTeamInfoToModel(team.Info)
	}

	if len(team.Members) > 0 {
		model.Members = r.mapMembersToModel(team.Members)
	}

	if team.Scope != nil {
		scopeModel, scopeDiags := r.mapScopeToModel(ctx, team.Scope)
		diags.Append(scopeDiags...)
		if !diags.HasError() {
			model.Scope = scopeModel
		}
	}

	return model, diags
}

// mapTeamInfoToModel converts API team info to model team info
func (r *teamResource) mapTeamInfoToModel(apiInfo *api.TeamInfo) *TeamInfoModel {
	return &TeamInfoModel{
		Description: util.SetStringPointerToState(apiInfo.Description),
	}
}

// mapMembersToModel converts API members to model members
func (r *teamResource) mapMembersToModel(apiMembers []api.TeamMember) []TeamMemberModel {
	members := make([]TeamMemberModel, len(apiMembers))
	for i, apiMember := range apiMembers {
		members[i] = TeamMemberModel{
			UserID: types.StringValue(apiMember.UserID),
			Roles:  r.mapRolesToModel(apiMember.Roles),
		}
	}
	return members
}

// mapRolesToModel converts API roles to model roles
func (r *teamResource) mapRolesToModel(apiRoles []api.TeamRole) []TeamMemberRole {
	if len(apiRoles) == 0 {
		return nil
	}

	roles := make([]TeamMemberRole, len(apiRoles))
	for i, apiRole := range apiRoles {
		roles[i] = TeamMemberRole{
			RoleID: types.StringValue(apiRole.RoleID),
		}
	}
	return roles
}

// mapScopeToModel converts API scope to model scope.
// Each []string scope field from the API (null or a list) is stored as types.Set so that
// the Optional+Computed schema attribute can correctly distinguish unknown/null/empty during plan.
func (r *teamResource) mapScopeToModel(ctx context.Context, apiScope *api.TeamScope) (*TeamScopeModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	scopeModel := &TeamScopeModel{
		InfraDFQFilter: util.SetStringPointerToState(apiScope.InfraDFQFilter),
		ActionFilter:   util.SetStringPointerToState(apiScope.ActionFilter),
		LogFilter:      util.SetStringPointerToState(apiScope.LogFilter),
	}

	setFields := []struct {
		apiVal []string
		target *types.Set
	}{
		{apiScope.AccessPermissions, &scopeModel.AccessPermissions},
		{apiScope.Applications, &scopeModel.Applications},
		{apiScope.KubernetesClusters, &scopeModel.KubernetesClusters},
		{apiScope.KubernetesNamespaces, &scopeModel.KubernetesNamespaces},
		{apiScope.MobileApps, &scopeModel.MobileApps},
		{apiScope.Websites, &scopeModel.Websites},
		{apiScope.BusinessPerspectives, &scopeModel.BusinessPerspectives},
		{apiScope.SloIDs, &scopeModel.SloIDs},
		{apiScope.SyntheticTests, &scopeModel.SyntheticTests},
		{apiScope.SyntheticCredentials, &scopeModel.SyntheticCredentials},
		{apiScope.TagIDs, &scopeModel.TagIDs},
	}

	for _, f := range setFields {
		s, d := stringsToSet(ctx, f.apiVal)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		*f.target = s
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
func (r *teamResource) mapRestrictedApplicationFilterToModel(apiFilter *api.RestrictedApplicationFilter) (*TeamRestrictedApplicationFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	filterModel := &TeamRestrictedApplicationFilterModel{
		Label: util.SetStringPointerToState(apiFilter.Label),
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
func (r *teamResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*api.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	model, modelDiags := r.extractModelFromPlanOrState(ctx, plan, state)
	diags.Append(modelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	team := &api.Team{
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
		scopeAPI, scopeDiags := r.mapModelScopeToAPI(ctx, model.Scope)
		diags.Append(scopeDiags...)
		if !diags.HasError() {
			team.Scope = scopeAPI
		}
	}

	return team, diags
}

// extractModelFromPlanOrState retrieves the TeamModel from plan or state
func (r *teamResource) extractModelFromPlanOrState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (TeamModel, diag.Diagnostics) {
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
func (r *teamResource) extractTeamID(model TeamModel) string {
	if model.ID.IsNull() {
		return ""
	}
	return model.ID.ValueString()
}

// mapModelInfoToAPI converts model team info to API team info
func (r *teamResource) mapModelInfoToAPI(modelInfo *TeamInfoModel) *api.TeamInfo {
	if modelInfo.Description.IsNull() || modelInfo.Description.IsUnknown() {
		return nil
	}

	desc := modelInfo.Description.ValueString()
	return &api.TeamInfo{
		Description: &desc,
	}
}

// mapModelMembersToAPI converts model members to API members
func (r *teamResource) mapModelMembersToAPI(modelMembers []TeamMemberModel) []api.TeamMember {
	if len(modelMembers) == 0 {
		return nil
	}

	apiMembers := make([]api.TeamMember, 0, len(modelMembers))
	for _, memberModel := range modelMembers {
		apiMember := api.TeamMember{
			UserID: memberModel.UserID.ValueString(),
		}

		if len(memberModel.Roles) > 0 {
			apiMember.Roles = r.mapModelRolesToAPI(memberModel.Roles)
		}

		apiMembers = append(apiMembers, apiMember)
	}

	return apiMembers
}

// mapModelRolesToAPI converts model roles to API roles
func (r *teamResource) mapModelRolesToAPI(modelRoles []TeamMemberRole) []api.TeamRole {
	if len(modelRoles) == 0 {
		return nil
	}

	apiRoles := make([]api.TeamRole, len(modelRoles))
	for i, roleModel := range modelRoles {
		apiRoles[i] = api.TeamRole{
			RoleID: roleModel.RoleID.ValueString(),
		}
	}

	return apiRoles
}

// mapModelScopeToAPI converts model scope to API scope
func (r *teamResource) mapModelScopeToAPI(ctx context.Context, modelScope *TeamScopeModel) (*api.TeamScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiScope := &api.TeamScope{}

	// Extract each types.Set field back to []string for the API call.
	setFields := []struct {
		src    types.Set
		target *[]string
	}{
		{modelScope.AccessPermissions, &apiScope.AccessPermissions},
		{modelScope.Applications, &apiScope.Applications},
		{modelScope.KubernetesClusters, &apiScope.KubernetesClusters},
		{modelScope.KubernetesNamespaces, &apiScope.KubernetesNamespaces},
		{modelScope.MobileApps, &apiScope.MobileApps},
		{modelScope.Websites, &apiScope.Websites},
		{modelScope.BusinessPerspectives, &apiScope.BusinessPerspectives},
		{modelScope.SloIDs, &apiScope.SloIDs},
		{modelScope.SyntheticTests, &apiScope.SyntheticTests},
		{modelScope.SyntheticCredentials, &apiScope.SyntheticCredentials},
		{modelScope.TagIDs, &apiScope.TagIDs},
	}

	for _, f := range setFields {
		ss, d := setToStrings(ctx, f.src)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		*f.target = ss
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
func (r *teamResource) mapModelRestrictedApplicationFilterToAPI(modelFilter *TeamRestrictedApplicationFilterModel) (*api.RestrictedApplicationFilter, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiFilter := &api.RestrictedApplicationFilter{}

	if !modelFilter.Label.IsNull() && !modelFilter.Label.IsUnknown() {
		label := modelFilter.Label.ValueString()
		apiFilter.Label = &label
	}

	if !modelFilter.Scope.IsNull() && !modelFilter.Scope.IsUnknown() {
		scope := api.RestrictedApplicationFilterScope(modelFilter.Scope.ValueString())
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

// GetStateUpgraders returns the state upgraders for this resource
func (r *teamResource) GetStateUpgraders(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: resourcehandle.CreateStateUpgraderForVersion(0),
	}
}
