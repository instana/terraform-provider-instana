package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/require"
)

func TestNewRbacRoleDataSource(t *testing.T) {
	ds := NewRbacRoleDataSource()
	_, ok := ds.(*RbacRoleDataSource)
	require.True(t, ok)
}

func TestRbacRoleDataSourceMetadata(t *testing.T) {
	ds := NewRbacRoleDataSource()
	req := datasource.MetadataRequest{ProviderTypeName: "instana"}
	resp := &datasource.MetadataResponse{}
	ds.Metadata(context.Background(), req, resp)
	require.Equal(t, "instana_rbac_role", resp.TypeName)
}

func TestRbacRoleDataSourceSchema(t *testing.T) {
	ds := NewRbacRoleDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	require.Equal(t, RbacRoleDescDataSource, resp.Schema.Description)
	_, ok := resp.Schema.Attributes[RbacRoleDataSourceFieldID]
	require.True(t, ok)
	_, ok = resp.Schema.Attributes["name"]
	require.True(t, ok)
}
