package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewAutomationActionDataSource(t *testing.T) {
	ds := NewAutomationActionDataSource()
	require.NotNil(t, ds)
}

func TestAutomationActionDataSourceMetadata(t *testing.T) {
	ds := NewAutomationActionDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_automation_action", resp.TypeName)
}

func TestAutomationActionDataSourceSchema(t *testing.T) {
	ds := NewAutomationActionDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, AutomationActionDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, AutomationActionFieldID)
	require.Contains(t, resp.Schema.Attributes, AutomationActionFieldName)
	require.Contains(t, resp.Schema.Attributes, AutomationActionFieldDescription)
	require.Contains(t, resp.Schema.Attributes, AutomationActionFieldType)
	require.Contains(t, resp.Schema.Attributes, AutomationActionFieldTags)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[AutomationActionFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify name field is required
	nameAttr := resp.Schema.Attributes[AutomationActionFieldName]
	require.True(t, nameAttr.(schema.StringAttribute).Required)

	// Verify type field is required
	typeAttr := resp.Schema.Attributes[AutomationActionFieldType]
	require.True(t, typeAttr.(schema.StringAttribute).Required)

	// Verify description field is computed
	descAttr := resp.Schema.Attributes[AutomationActionFieldDescription]
	require.True(t, descAttr.(schema.StringAttribute).Computed)

	// Verify tags field is computed
	tagsAttr := resp.Schema.Attributes[AutomationActionFieldTags]
	require.True(t, tagsAttr.(schema.ListAttribute).Computed)
}

// Made with Bob
