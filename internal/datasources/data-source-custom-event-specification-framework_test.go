package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewCustomEventSpecificationDataSourceFramework(t *testing.T) {
	ds := NewCustomEventSpecificationDataSourceFramework()
	require.NotNil(t, ds)
}

func TestCustomEventSpecificationDataSourceMetadata(t *testing.T) {
	ds := NewCustomEventSpecificationDataSourceFramework()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_custom_event_spec", resp.TypeName)
}

func TestCustomEventSpecificationDataSourceSchema(t *testing.T) {
	ds := NewCustomEventSpecificationDataSourceFramework()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, CustomEventSpecificationDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldID)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldName)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldDescription)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldEntityType)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldTriggering)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldEnabled)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldQuery)
	require.Contains(t, resp.Schema.Attributes, CustomEventSpecificationFieldExpirationTime)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[CustomEventSpecificationFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify name field is required
	nameAttr := resp.Schema.Attributes[CustomEventSpecificationFieldName]
	require.True(t, nameAttr.(schema.StringAttribute).Required)

	// Verify entity_type field is required
	entityTypeAttr := resp.Schema.Attributes[CustomEventSpecificationFieldEntityType]
	require.True(t, entityTypeAttr.(schema.StringAttribute).Required)

	// Verify computed fields
	descAttr := resp.Schema.Attributes[CustomEventSpecificationFieldDescription]
	require.True(t, descAttr.(schema.StringAttribute).Computed)

	triggeringAttr := resp.Schema.Attributes[CustomEventSpecificationFieldTriggering]
	require.True(t, triggeringAttr.(schema.BoolAttribute).Computed)

	enabledAttr := resp.Schema.Attributes[CustomEventSpecificationFieldEnabled]
	require.True(t, enabledAttr.(schema.BoolAttribute).Computed)

	queryAttr := resp.Schema.Attributes[CustomEventSpecificationFieldQuery]
	require.True(t, queryAttr.(schema.StringAttribute).Computed)

	expirationTimeAttr := resp.Schema.Attributes[CustomEventSpecificationFieldExpirationTime]
	require.True(t, expirationTimeAttr.(schema.Int64Attribute).Computed)
}

// Made with Bob
