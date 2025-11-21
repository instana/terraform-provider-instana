package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewHostAgentsDataSource(t *testing.T) {
	ds := NewHostAgentsDataSource()
	require.NotNil(t, ds)
}

func TestHostAgentsDataSourceMetadata(t *testing.T) {
	ds := NewHostAgentsDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_host_agents", resp.TypeName)
}

func TestHostAgentsDataSourceSchema(t *testing.T) {
	ds := NewHostAgentsDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, HostAgentDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, HostAgentFieldID)
	require.Contains(t, resp.Schema.Attributes, HostAgentFieldFilter)
	require.Contains(t, resp.Schema.Attributes, HostAgentFieldItems)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[HostAgentFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify filter field is required
	filterAttr := resp.Schema.Attributes[HostAgentFieldFilter]
	require.True(t, filterAttr.(schema.StringAttribute).Required)

	// Verify items field is computed and is a list
	itemsAttr := resp.Schema.Attributes[HostAgentFieldItems]
	listAttr, ok := itemsAttr.(schema.ListNestedAttribute)
	require.True(t, ok, "items should be a ListNestedAttribute")
	require.True(t, listAttr.Computed)

	// Verify nested attributes in items
	nestedAttrs := listAttr.NestedObject.Attributes
	require.Contains(t, nestedAttrs, HostAgentFieldSnapshotId)
	require.Contains(t, nestedAttrs, HostAgentFieldLabel)
	require.Contains(t, nestedAttrs, HostAgentFieldHost)
	require.Contains(t, nestedAttrs, HostAgentFieldPlugin)
	require.Contains(t, nestedAttrs, HostAgentFieldTags)

	// Verify all nested fields are computed
	require.True(t, nestedAttrs[HostAgentFieldSnapshotId].(schema.StringAttribute).Computed)
	require.True(t, nestedAttrs[HostAgentFieldLabel].(schema.StringAttribute).Computed)
	require.True(t, nestedAttrs[HostAgentFieldHost].(schema.StringAttribute).Computed)
	require.True(t, nestedAttrs[HostAgentFieldPlugin].(schema.StringAttribute).Computed)
	require.True(t, nestedAttrs[HostAgentFieldTags].(schema.ListAttribute).Computed)
}

// Made with Bob
