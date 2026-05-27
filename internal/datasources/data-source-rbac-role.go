package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			       Description: "ID of the RBAC Role.",
			       Required:    true,
		       },
		       "name": schema.StringAttribute{
			       Description: "Name of the RBAC Role.",
			       Computed:    true,
		       },
		       "permissions": schema.SetAttribute{
			       Description: "Permissions assigned to the RBAC Role.",
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
       resp.Diagnostics.Append(req.Config.Get(ctx, &data)...) // Read config into data
       if resp.Diagnostics.HasError() {
	       return
       }

       id := data.ID.ValueString()
       if id == "" {
	       resp.Diagnostics.AddError("Missing ID", "The 'id' attribute must be set.")
	       return
       }

       role, err := d.instanaAPI.Roles().GetOne(id)
       if err != nil {
	       resp.Diagnostics.AddError("Error reading RBAC Role", fmt.Sprintf("Could not read RBAC Role with ID '%s': %s", id, err))
	       return
       }
       if role == nil {
	       resp.Diagnostics.AddError("Not found", fmt.Sprintf("No RBAC Role found with ID '%s'", id))
	       return
       }

       data.ID = types.StringValue(role.ID)
       data.Name = types.StringValue(role.Name)
       data.Permissions = make([]types.String, len(role.Permissions))
       for i, p := range role.Permissions {
	       data.Permissions[i] = types.StringValue(p)
       }

       resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) 
}
