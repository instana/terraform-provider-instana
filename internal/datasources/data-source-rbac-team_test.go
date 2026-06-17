package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/instana/terraform-provider-instana/internal/resources/team"
	"github.com/stretchr/testify/require"
)

func TestNewRbacTeamDataSource(t *testing.T) {
	ds := NewRbacTeamDataSource()
	require.NotNil(t, ds)
}

func TestRbacTeamDataSourceMetadata(t *testing.T) {
	ds := NewRbacTeamDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_rbac_team", resp.TypeName)
}

func TestRbacTeamDataSourceSchema(t *testing.T) {
	ds := NewRbacTeamDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, RbacTeamDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, RbacTeamDataSourceFieldID)
	require.Contains(t, resp.Schema.Attributes, RbacTeamDataSourceFieldTag)
	require.Contains(t, resp.Schema.Attributes, team.TeamFieldInfo)
	require.Contains(t, resp.Schema.Attributes, team.TeamFieldMembers)
	require.Contains(t, resp.Schema.Attributes, team.TeamFieldScope)

	// Verify ID field is optional and computed
	idAttr := resp.Schema.Attributes[RbacTeamDataSourceFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Optional)
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify tag field is optional and computed
	tagAttr := resp.Schema.Attributes[RbacTeamDataSourceFieldTag]
	require.True(t, tagAttr.(schema.StringAttribute).Optional)
	require.True(t, tagAttr.(schema.StringAttribute).Computed)
}
