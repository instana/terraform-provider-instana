package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/restapi"
)

// DataSourceInstanaUsers the name of the terraform-provider-instana data source to read multiple users
const DataSourceInstanaUsers = "users"

// UsersDataSourceModel represents the data model for the users data source
type UsersDataSourceModel struct {
	ID     types.String    `tfsdk:"id"`
	Emails []string        `tfsdk:"emails"`
	Users  []UserItemModel `tfsdk:"users"`
}

// UserItemModel represents a single user in the list
type UserItemModel struct {
	ID       types.String `tfsdk:"id"`
	Email    types.String `tfsdk:"email"`
	FullName types.String `tfsdk:"full_name"`
}

// NewUsersDataSource creates a new data source for multiple users
func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

type UsersDataSource struct {
	instanaAPI restapi.InstanaAPI
}

func (d *UsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceInstanaUsers
}

func (d *UsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: UsersDescDataSource,
		Attributes: map[string]schema.Attribute{
			UsersDataSourceFieldID: schema.StringAttribute{
				Description: UsersDescID,
				Computed:    true,
			},
			UsersFieldEmails: schema.ListAttribute{
				Description: UsersDescEmails,
				Required:    true,
				ElementType: types.StringType,
			},
			UsersFieldUsers: schema.ListNestedAttribute{
				Description: UsersDescUsers,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						UsersFieldUserID: schema.StringAttribute{
							Description: UsersDescUserID,
							Computed:    true,
						},
						UsersFieldUserEmail: schema.StringAttribute{
							Description: UsersDescUserEmail,
							Computed:    true,
						},
						UsersFieldUserFullName: schema.StringAttribute{
							Description: UsersDescUserFullName,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerMeta, ok := req.ProviderData.(*restapi.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			UsersErrUnexpectedConfigureType,
			fmt.Sprintf("Expected *restapi.ProviderMeta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.instanaAPI = providerMeta.InstanaAPI
}

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all users from the API
	allUsers, err := d.instanaAPI.Users().GetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			UsersErrReadingUsers,
			fmt.Sprintf("Could not read users: %s", err),
		)
		return
	}

	// Create a map of emails for quick lookup
	emailMap := make(map[string]bool)
	for _, email := range data.Emails {
		emailMap[email] = true
	}

	// Filter users based on the provided emails
	var matchedUsers []UserItemModel
	for _, user := range *allUsers {
		if emailMap[user.Email] {
			matchedUsers = append(matchedUsers, UserItemModel{
				ID:       types.StringValue(user.ID),
				Email:    types.StringValue(user.Email),
				FullName: types.StringValue(user.FullName),
			})
		}
	}

	// Set the results
	data.Users = matchedUsers
	data.ID = types.StringValue(fmt.Sprintf("users-%d", len(matchedUsers)))

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Made with Bob
