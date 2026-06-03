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
	"github.com/instana/terraform-provider-instana/internal/shared"
)

const DataSourceInstanaRbacRole = "rbac_role"

// RbacRoleDataSourceModel represents the data model for the rbac_role data source
type RbacRoleDataSourceModel struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Permissions []types.String `tfsdk:"permissions"`
}

// NewRbacRoleDataSource creates a new data source for rbac_role
func NewRbacRoleDataSource() datasource.DataSource {
	return &RbacRoleDataSource{}
}

type RbacRoleDataSource struct {
	instanaAPI client.InstanaAPI
}

func (d *RbacRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaRbacRole
}

func (d *RbacRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: RbacRoleDescDataSource,
		Attributes: map[string]schema.Attribute{
			RbacRoleDataSourceFieldID: schema.StringAttribute{
				Description: RbacRoleDescID,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot(RbacRoleDataSourceFieldID), path.MatchRoot(RbacRoleDataSourceFieldName)),
				},
			},
			RbacRoleDataSourceFieldName: schema.StringAttribute{
				Description: RbacRoleDescName,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot(RbacRoleDataSourceFieldID), path.MatchRoot(RbacRoleDataSourceFieldName)),
				},
			},
			RbacRoleDataSourceFieldPermissions: schema.SetAttribute{
				Description: RbacRoleDescPermissions,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *RbacRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RbacRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RbacRoleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	name := data.Name.ValueString()
	if id == "" && name == "" {
		resp.Diagnostics.AddError("Missing Attribute", RbacRoleErrMissingLookupAttribute)
		return
	}

		var role *api.Role
	var err error
	if id != "" {
		role, err = d.instanaAPI.Roles().GetOne(id)
		if err != nil {
			resp.Diagnostics.AddError("Error reading RBAC Role", fmt.Sprintf(RbacRoleErrReadByID, id, err))
			return
		}
		if role == nil {
			resp.Diagnostics.AddError("Not found", fmt.Sprintf(RbacRoleErrNotFoundByID, id))
			return
		}
	} else {
		roles, err := d.instanaAPI.Roles().GetAll()
		if err != nil {
			resp.Diagnostics.AddError("Error reading RBAC Roles", fmt.Sprintf(RbacRoleErrReadAll, err))
			return
		}
		for _, r := range *roles {
			if r.Name == name {
				role = r
				break
			}
		}
		if role == nil {
			resp.Diagnostics.AddError("Not found", fmt.Sprintf(RbacRoleErrNotFoundByName, name))
			return
		}
	}

	data.ID = types.StringValue(role.ID)
	data.Name = types.StringValue(role.Name)
	data.Permissions = make([]types.String, len(role.Permissions))
	for i, p := range role.Permissions {
		data.Permissions[i] = types.StringValue(p)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
