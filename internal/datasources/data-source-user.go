package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/restapi"
)

// DataSourceInstanaUser the name of the terraform-provider-instana data source to read user
const DataSourceInstanaUser = "user"

// UserDataSourceModel represents the data model for the user data source
type UserDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Email    types.String `tfsdk:"email"`
	FullName types.String `tfsdk:"full_name"`
}

// NewUserDataSource creates a new data source for user
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	instanaAPI restapi.InstanaAPI
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaUser
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: UserDescDataSource,
		Attributes: map[string]schema.Attribute{
			UserDataSourceFieldID: schema.StringAttribute{
				Description: UserDescID,
				Computed:    true,
			},
			UserFieldEmail: schema.StringAttribute{
				Description: UserDescEmail,
				Required:    true,
			},
			UserFieldFullName: schema.StringAttribute{
				Description: UserDescFullName,
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*restapi.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			UserErrUnexpectedConfigureType,
			fmt.Sprintf("Expected *restapi.ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the email from the configuration
	email := data.Email.ValueString()

	// Get all users
	users, err := d.instanaAPI.Users().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			UserErrReadingUsers,
			fmt.Sprintf("Could not read users: %s", err),
		)
		return
	}

	// Find the user with the matching email
	var matchingUser *restapi.User
	for _, user := range *users {
		if user.Email == email {
			matchingUser = user
			break
		}
	}

	if matchingUser == nil {
		resp.Diagnostics.AddError(
			UserErrUserNotFound,
			fmt.Sprintf("No user found with email: %s", email),
		)
		return
	}

	// Update the data model with the user details
	data.ID = types.StringValue(matchingUser.ID)
	data.FullName = types.StringValue(matchingUser.FullName)

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
