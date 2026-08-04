package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/client"
	"github.com/instana/terraform-provider-instana/internal/resources/team"
	"github.com/instana/terraform-provider-instana/internal/shared"
)

const DataSourceInstanaRbacTeam = "rbac_team"

// RbacTeamDataSourceModel represents the data model for the rbac_team data source
type RbacTeamDataSourceModel struct {
	ID      types.String           `tfsdk:"id"`
	Tag     types.String           `tfsdk:"tag"`
	Info    *team.TeamInfoModel    `tfsdk:"info"`
	Members []team.TeamMemberModel `tfsdk:"member"`
	Scope   *team.TeamScopeModel   `tfsdk:"scope"`
}

// NewRbacTeamDataSource creates a new data source for rbac_team
func NewRbacTeamDataSource() datasource.DataSource {
	return &RbacTeamDataSource{}
}

type RbacTeamDataSource struct {
	instanaAPI client.InstanaAPI
}

func (d *RbacTeamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaRbacTeam
}

func (d *RbacTeamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: RbacTeamDescDataSource,
		Attributes: map[string]schema.Attribute{
			RbacTeamDataSourceFieldID: schema.StringAttribute{
				Description: RbacTeamDescID,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot(RbacTeamDataSourceFieldID), path.MatchRoot(RbacTeamDataSourceFieldTag)),
				},
			},
			RbacTeamDataSourceFieldTag: schema.StringAttribute{
				Description: RbacTeamDescTag,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot(RbacTeamDataSourceFieldID), path.MatchRoot(RbacTeamDataSourceFieldTag)),
				},
			},
			team.TeamFieldInfo: schema.SingleNestedAttribute{
				Description: team.TeamDescInfo,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					team.TeamFieldInfoDescription: schema.StringAttribute{
						Description: team.TeamDescInfoDescription,
						Computed:    true,
					},
				},
			},
			team.TeamFieldMembers: schema.SetNestedAttribute{
				Description: team.TeamDescMembers,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						team.TeamFieldMemberUserID: schema.StringAttribute{
							Description: team.TeamDescMemberUserID,
							Computed:    true,
						},
						team.TeamFieldMemberRoles: schema.SetNestedAttribute{
							Description: team.TeamDescMemberRoles,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									team.TeamFieldMemberRoleID: schema.StringAttribute{
										Description: team.TeamDescMemberRoleID,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			team.TeamFieldScope: schema.SingleNestedAttribute{
				Description: team.TeamDescScope,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					team.TeamFieldScopeAccessPermissions: schema.SetAttribute{
						Description: team.TeamDescScopeAccessPermissions,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeApplications: schema.SetAttribute{
						Description: team.TeamDescScopeApplications,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeKubernetesClusters: schema.SetAttribute{
						Description: team.TeamDescScopeKubernetesClusters,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeKubernetesNamespaces: schema.SetAttribute{
						Description: team.TeamDescScopeKubernetesNamespaces,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeMobileApps: schema.SetAttribute{
						Description: team.TeamDescScopeMobileApps,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeWebsites: schema.SetAttribute{
						Description: team.TeamDescScopeWebsites,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeInfraDFQFilter: schema.StringAttribute{
						Description: team.TeamDescScopeInfraDFQFilter,
						Computed:    true,
					},
					team.TeamFieldScopeActionFilter: schema.StringAttribute{
						Description: team.TeamDescScopeActionFilter,
						Computed:    true,
					},
					team.TeamFieldScopeLogFilter: schema.StringAttribute{
						Description: team.TeamDescScopeLogFilter,
						Computed:    true,
					},
					team.TeamFieldScopeBusinessPerspectives: schema.SetAttribute{
						Description: team.TeamDescScopeBusinessPerspectives,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeSloIDs: schema.SetAttribute{
						Description: team.TeamDescScopeSloIDs,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeSyntheticTests: schema.SetAttribute{
						Description: team.TeamDescScopeSyntheticTests,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeSyntheticCredentials: schema.SetAttribute{
						Description: team.TeamDescScopeSyntheticCredentials,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeTagIDs: schema.SetAttribute{
						Description: team.TeamDescScopeTagIDs,
						Computed:    true,
						ElementType: types.StringType,
					},
					team.TeamFieldScopeRestrictedApplicationFilter: schema.SingleNestedAttribute{
						Description: team.TeamDescScopeRestrictedApplicationFilter,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							team.TeamFieldScopeRestrictedApplicationFilterLabel: schema.StringAttribute{
								Description: team.TeamDescScopeRestrictedApplicationFilterLabel,
								Computed:    true,
							},
							team.TeamFieldScopeRestrictedApplicationFilterScope: schema.StringAttribute{
								Description: team.TeamDescScopeRestrictedApplicationFilterScope,
								Computed:    true,
							},
							team.TeamFieldScopeRestrictedApplicationFilterTagFilterExpression: schema.StringAttribute{
								Description: team.TeamDescScopeRestrictedApplicationFilterTagFilterExpression,
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *RbacTeamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*shared.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *shared.ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *RbacTeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RbacTeamDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	tag := data.Tag.ValueString()
	if id == "" && tag == "" {
		resp.Diagnostics.AddError("Missing Attribute", RbacTeamErrMissingLookupAttribute)
		return
	}

	var teamData *api.Team
	var err error
	if id != "" {
		teamData, err = d.instanaAPI.Teams().GetOne(id)
		if err != nil {
			resp.Diagnostics.AddError("Error reading RBAC Team", fmt.Sprintf(RbacTeamErrReadByID, id, err))
			return
		}
		if teamData == nil {
			resp.Diagnostics.AddError("Not found", fmt.Sprintf(RbacTeamErrNotFoundByID, id))
			return
		}
	} else {
		teams, err := d.instanaAPI.Teams().GetAll()
		if err != nil {
			resp.Diagnostics.AddError("Error reading RBAC Teams", fmt.Sprintf(RbacTeamErrReadAll, err))
			return
		}
		for _, t := range *teams {
			if t.Tag == tag {
				teamData = t
				break
			}
		}
		if teamData == nil {
			resp.Diagnostics.AddError("Not found", fmt.Sprintf(RbacTeamErrNotFoundByTag, tag))
			return
		}
	}

	// Map the API response to the data source model
	data.ID = types.StringValue(teamData.ID)
	data.Tag = types.StringValue(teamData.Tag)

	// Map team info
	if teamData.Info != nil {
		data.Info = &team.TeamInfoModel{
			Description: types.StringPointerValue(teamData.Info.Description),
		}
	}

	// Map members
	if len(teamData.Members) > 0 {
		data.Members = make([]team.TeamMemberModel, len(teamData.Members))
		for i, member := range teamData.Members {
			data.Members[i] = team.TeamMemberModel{
				UserID: types.StringValue(member.UserID),
			}
			if len(member.Roles) > 0 {
				data.Members[i].Roles = make([]team.TeamMemberRole, len(member.Roles))
				for j, role := range member.Roles {
					data.Members[i].Roles[j] = team.TeamMemberRole{
						RoleID: types.StringValue(role.RoleID),
					}
				}
			}
		}
	}

	// Map scope
	if teamData.Scope != nil {
		scopeModel := &team.TeamScopeModel{
			InfraDFQFilter: types.StringPointerValue(teamData.Scope.InfraDFQFilter),
			ActionFilter:   types.StringPointerValue(teamData.Scope.ActionFilter),
			LogFilter:      types.StringPointerValue(teamData.Scope.LogFilter),
		}

		// Convert each []string scope field to types.Set
		type setField struct {
			src    []string
			target *types.Set
		}
		fields := []setField{
			{teamData.Scope.AccessPermissions, &scopeModel.AccessPermissions},
			{teamData.Scope.Applications, &scopeModel.Applications},
			{teamData.Scope.KubernetesClusters, &scopeModel.KubernetesClusters},
			{teamData.Scope.KubernetesNamespaces, &scopeModel.KubernetesNamespaces},
			{teamData.Scope.MobileApps, &scopeModel.MobileApps},
			{teamData.Scope.Websites, &scopeModel.Websites},
			{teamData.Scope.BusinessPerspectives, &scopeModel.BusinessPerspectives},
			{teamData.Scope.SloIDs, &scopeModel.SloIDs},
			{teamData.Scope.SyntheticTests, &scopeModel.SyntheticTests},
			{teamData.Scope.SyntheticCredentials, &scopeModel.SyntheticCredentials},
			{teamData.Scope.TagIDs, &scopeModel.TagIDs},
		}
		for _, f := range fields {
			s, d := types.SetValueFrom(ctx, types.StringType, f.src)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			*f.target = s
		}

		// Map restricted application filter if present
		if teamData.Scope.RestrictedApplicationFilter != nil {
			scopeModel.RestrictedApplicationFilter = &team.TeamRestrictedApplicationFilterModel{
				Label: types.StringPointerValue(teamData.Scope.RestrictedApplicationFilter.Label),
			}
			if teamData.Scope.RestrictedApplicationFilter.Scope != nil {
				scopeModel.RestrictedApplicationFilter.Scope = types.StringValue(string(*teamData.Scope.RestrictedApplicationFilter.Scope))
			}
			scopeModel.RestrictedApplicationFilter.TagFilterExpression = types.StringNull()
		}

		data.Scope = scopeModel
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
