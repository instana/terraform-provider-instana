package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewUsersDataSource(t *testing.T) {
	ds := NewUsersDataSource()
	require.NotNil(t, ds)
}

func TestUsersDataSourceMetadata(t *testing.T) {
	ds := NewUsersDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_users", resp.TypeName)
}

func TestUsersDataSourceSchema(t *testing.T) {
	ds := NewUsersDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, UsersDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, UsersDataSourceFieldID)
	require.Contains(t, resp.Schema.Attributes, UsersFieldEmails)
	require.Contains(t, resp.Schema.Attributes, UsersFieldUsers)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[UsersDataSourceFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify emails field is required and is a list
	emailsAttr := resp.Schema.Attributes[UsersFieldEmails]
	listAttr, ok := emailsAttr.(schema.ListAttribute)
	require.True(t, ok, "emails should be a ListAttribute")
	require.True(t, listAttr.Required)

	// Verify users field is computed and is a list
	usersAttr := resp.Schema.Attributes[UsersFieldUsers]
	listNestedAttr, ok := usersAttr.(schema.ListNestedAttribute)
	require.True(t, ok, "users should be a ListNestedAttribute")
	require.True(t, listNestedAttr.Computed)

	// Verify nested attributes in users
	nestedAttrs := listNestedAttr.NestedObject.Attributes
	require.Contains(t, nestedAttrs, UsersFieldUserID)
	require.Contains(t, nestedAttrs, UsersFieldUserEmail)
	require.Contains(t, nestedAttrs, UsersFieldUserFullName)

	// Verify all nested fields are computed
	require.True(t, nestedAttrs[UsersFieldUserID].(schema.StringAttribute).Computed)
	require.True(t, nestedAttrs[UsersFieldUserEmail].(schema.StringAttribute).Computed)
	require.True(t, nestedAttrs[UsersFieldUserFullName].(schema.StringAttribute).Computed)
}

// Made with Bob
