package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewUserDataSource(t *testing.T) {
	ds := NewUserDataSource()
	require.NotNil(t, ds)
}

func TestUserDataSourceMetadata(t *testing.T) {
	ds := NewUserDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_user", resp.TypeName)
}

func TestUserDataSourceSchema(t *testing.T) {
	ds := NewUserDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, UserDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, UserDataSourceFieldID)
	require.Contains(t, resp.Schema.Attributes, UserFieldEmail)
	require.Contains(t, resp.Schema.Attributes, UserFieldFullName)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[UserDataSourceFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify email field is required
	emailAttr := resp.Schema.Attributes[UserFieldEmail]
	require.True(t, emailAttr.(schema.StringAttribute).Required)

	// Verify full_name field is computed
	fullNameAttr := resp.Schema.Attributes[UserFieldFullName]
	require.True(t, fullNameAttr.(schema.StringAttribute).Computed)
}

// Made with Bob
