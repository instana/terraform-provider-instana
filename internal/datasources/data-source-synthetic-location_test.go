package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewSyntheticLocationDataSource(t *testing.T) {
	ds := NewSyntheticLocationDataSource()
	require.NotNil(t, ds)
}

func TestSyntheticLocationDataSourceMetadata(t *testing.T) {
	ds := NewSyntheticLocationDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_synthetic_location", resp.TypeName)
}

func TestSyntheticLocationDataSourceSchema(t *testing.T) {
	ds := NewSyntheticLocationDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, SyntheticLocationDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, SyntheticLocationFieldID)
	require.Contains(t, resp.Schema.Attributes, SyntheticLocationFieldLabel)
	require.Contains(t, resp.Schema.Attributes, SyntheticLocationFieldDescription)
	require.Contains(t, resp.Schema.Attributes, SyntheticLocationFieldLocationType)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[SyntheticLocationFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify label field is optional
	labelAttr := resp.Schema.Attributes[SyntheticLocationFieldLabel]
	require.True(t, labelAttr.(schema.StringAttribute).Optional)

	// Verify description field is optional and computed
	descAttr := resp.Schema.Attributes[SyntheticLocationFieldDescription]
	require.True(t, descAttr.(schema.StringAttribute).Optional)
	require.True(t, descAttr.(schema.StringAttribute).Computed)

	// Verify location_type field is optional
	locationTypeAttr := resp.Schema.Attributes[SyntheticLocationFieldLocationType]
	require.True(t, locationTypeAttr.(schema.StringAttribute).Optional)
}

// Made with Bob
