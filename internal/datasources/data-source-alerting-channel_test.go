package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stretchr/testify/require"
)

func TestNewAlertingChannelDataSource(t *testing.T) {
	ds := NewAlertingChannelDataSource()
	require.NotNil(t, ds)
}

func TestAlertingChannelDataSourceMetadata(t *testing.T) {
	ds := NewAlertingChannelDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "instana",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	require.Equal(t, "instana_alerting_channel", resp.TypeName)
}

func TestAlertingChannelDataSourceSchema(t *testing.T) {
	ds := NewAlertingChannelDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	require.NotNil(t, resp.Schema)
	require.Equal(t, AlertingChannelDescDataSource, resp.Schema.Description)

	// Verify required fields exist
	require.Contains(t, resp.Schema.Attributes, AlertingChannelDataSourceFieldID)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldName)

	// Verify channel type fields exist
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelEmail)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelOpsGenie)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelPageDuty)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelSlack)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelSplunk)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelVictorOps)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelWebhook)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelOffice365)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelGoogleChat)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelServiceNow)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelServiceNowApplication)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelPrometheusWebhook)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelWebexTeamsWebhook)
	require.Contains(t, resp.Schema.Attributes, AlertingChannelFieldChannelWatsonAIOpsWebhook)

	// Verify ID field is computed
	idAttr := resp.Schema.Attributes[AlertingChannelDataSourceFieldID]
	require.True(t, idAttr.(schema.StringAttribute).Computed)

	// Verify name field is required
	nameAttr := resp.Schema.Attributes[AlertingChannelFieldName]
	require.True(t, nameAttr.(schema.StringAttribute).Required)
}
