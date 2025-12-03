package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestNewBuiltinEventDataSource(t *testing.T) {
	ds := NewBuiltinEventDataSource()
	require.NotNil(t, ds)
}

func TestBuiltinEventDataSourceMetadata(t *testing.T) {
	ds := NewBuiltinEventDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_builtin_event_spec", resp.TypeName)
}

func TestBuiltinEventDataSourceSchema(t *testing.T) {
	ds := NewBuiltinEventDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, shared.BuiltinEventDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventFieldID)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldName)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldDescription)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldShortPluginID)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldSeverity)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldSeverityCode)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldTriggering)
	require.Contains(t, resp.Schema.Attributes, shared.BuiltinEventSpecificationFieldEnabled)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[shared.BuiltinEventFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify name field is required
	nameAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldName]
	require.True(t, nameAttr.(schema.StringAttribute).Required)

	// Verify short_plugin_id field is required
	shortPluginIDAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldShortPluginID]
	require.True(t, shortPluginIDAttr.(schema.StringAttribute).Required)

	// Verify computed fields
	descAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldDescription]
	require.True(t, descAttr.(schema.StringAttribute).Computed)

	severityAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldSeverity]
	require.True(t, severityAttr.(schema.StringAttribute).Computed)

	severityCodeAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldSeverityCode]
	require.True(t, severityCodeAttr.(schema.Int64Attribute).Computed)

	triggeringAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldTriggering]
	require.True(t, triggeringAttr.(schema.BoolAttribute).Computed)

	enabledAttr := resp.Schema.Attributes[shared.BuiltinEventSpecificationFieldEnabled]
	require.True(t, enabledAttr.(schema.BoolAttribute).Computed)
}
