package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewApplicationConfigDataSource(t *testing.T) {
	ds := NewApplicationConfigDataSource()
	require.NotNil(t, ds)
	_, ok := ds.(*applicationConfigDataSource)
	require.True(t, ok)
}

func TestApplicationConfigDataSourceMetadata(t *testing.T) {
	ds := NewApplicationConfigDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_application_config", resp.TypeName)
}

func TestApplicationConfigDataSourceSchema(t *testing.T) {
	ds := NewApplicationConfigDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, ApplicationConfigDataSourceDesc, resp.Schema.Description)

	// Verify attributes exist
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldID)
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldName)
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldLabel)
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldScope)
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldBoundaryScope)
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldTagFilter)
	require.Contains(t, resp.Schema.Attributes, ApplicationConfigDataSourceFieldAccessRules)

	// Verify ID field is optional and computed
	idAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Optional)
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify Name field is optional and computed
	nameAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldName]
	require.True(t, nameAttr.(schema.StringAttribute).Optional)
	require.True(t, nameAttr.(schema.StringAttribute).Computed)

	// Verify other fields are computed
	labelAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldLabel]
	require.True(t, labelAttr.(schema.StringAttribute).Computed)

	scopeAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldScope]
	require.True(t, scopeAttr.(schema.StringAttribute).Computed)

	boundaryScopeAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldBoundaryScope]
	require.True(t, boundaryScopeAttr.(schema.StringAttribute).Computed)

	tagFilterAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldTagFilter]
	require.True(t, tagFilterAttr.(schema.StringAttribute).Computed)

	accessRulesAttr := resp.Schema.Attributes[ApplicationConfigDataSourceFieldAccessRules]
	require.True(t, accessRulesAttr.(schema.ListNestedAttribute).Computed)
}
