package alertingchannel

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/instana/terraform-provider-instana/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlertingChannelResourceHandle(t *testing.T) {
	handle := NewAlertingChannelResourceHandle()
	require.NotNil(t, handle)

	metadata := handle.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, ResourceInstanaAlertingChannel, metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
	assert.NotNil(t, metadata.Schema)
}

func TestMetaData(t *testing.T) {
	resource := &alertingChannelResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 1,
		},
	}

	metadata := resource.MetaData()
	require.NotNil(t, metadata)
	assert.Equal(t, "test_resource", metadata.ResourceName)
	assert.Equal(t, int64(1), metadata.SchemaVersion)
}

func TestSetComputedFields(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	diags := resource.SetComputedFields(ctx, nil)
	assert.False(t, diags.HasError())
}

func TestMapEmailChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	emails := []string{"test1@example.com", "test2@example.com"}
	emailsSet, _ := types.SetValueFrom(ctx, types.StringType, emails)

	emailModel := &shared.EmailModel{
		Emails: emailsSet,
	}

	channel, diags := resource.mapEmailChannelFromState(ctx, "test-id", "test-name", emailModel)
	require.False(t, diags.HasError())
	require.NotNil(t, channel)
	assert.Equal(t, "test-id", channel.ID)
	assert.Equal(t, "test-name", channel.Name)
	assert.Equal(t, restapi.EmailChannelType, channel.Kind)
	assert.Equal(t, emails, channel.Emails)
}

func TestMapOpsGenieChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	tags := []string{"tag1", "tag2", "tag3"}
	tagsList, _ := types.ListValueFrom(ctx, types.StringType, tags)

	opsGenieModel := &shared.OpsGenieModel{
		APIKey: types.StringValue("test-api-key"),
		Region: types.StringValue("EU"),
		Tags:   tagsList,
	}

	channel, diags := resource.mapOpsGenieChannelFromState(ctx, "test-id", "test-name", opsGenieModel)
	require.False(t, diags.HasError())
	require.NotNil(t, channel)
	assert.Equal(t, "test-id", channel.ID)
	assert.Equal(t, "test-name", channel.Name)
	assert.Equal(t, restapi.OpsGenieChannelType, channel.Kind)
	assert.Equal(t, "test-api-key", *channel.APIKey)
	assert.Equal(t, "EU", *channel.Region)
	assert.Equal(t, "tag1,tag2,tag3", *channel.Tags)
}

func TestMapPagerDutyChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	pagerDutyModel := &shared.PagerDutyModel{
		ServiceIntegrationKey: types.StringValue("test-integration-key"),
	}

	channel, diags := resource.mapPagerDutyChannelFromState(ctx, "test-id", "test-name", pagerDutyModel)
	require.False(t, diags.HasError())
	require.NotNil(t, channel)
	assert.Equal(t, "test-id", channel.ID)
	assert.Equal(t, "test-name", channel.Name)
	assert.Equal(t, restapi.PagerDutyChannelType, channel.Kind)
	assert.Equal(t, "test-integration-key", *channel.ServiceIntegrationKey)
}

func TestMapSlackChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("with all fields", func(t *testing.T) {
		slackModel := &shared.SlackModel{
			WebhookURL: types.StringValue("https://hooks.slack.com/test"),
			IconURL:    types.StringValue("https://example.com/icon.png"),
			Channel:    types.StringValue("#alerts"),
		}

		channel, diags := resource.mapSlackChannelFromState(ctx, "test-id", "test-name", slackModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, "test-id", channel.ID)
		assert.Equal(t, "test-name", channel.Name)
		assert.Equal(t, restapi.SlackChannelType, channel.Kind)
		assert.Equal(t, "https://hooks.slack.com/test", *channel.WebhookURL)
		assert.Equal(t, "https://example.com/icon.png", *channel.IconURL)
		assert.Equal(t, "#alerts", *channel.Channel)
	})

	t.Run("with only required fields", func(t *testing.T) {
		slackModel := &shared.SlackModel{
			WebhookURL: types.StringValue("https://hooks.slack.com/test"),
			IconURL:    types.StringNull(),
			Channel:    types.StringNull(),
		}

		channel, diags := resource.mapSlackChannelFromState(ctx, "test-id", "test-name", slackModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.SlackChannelType, channel.Kind)
		assert.Nil(t, channel.IconURL)
		assert.Nil(t, channel.Channel)
	})
}

func TestMapSplunkChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	splunkModel := &shared.SplunkModel{
		URL:   types.StringValue("https://splunk.example.com"),
		Token: types.StringValue("test-token"),
	}

	channel, diags := resource.mapSplunkChannelFromState(ctx, "test-id", "test-name", splunkModel)
	require.False(t, diags.HasError())
	require.NotNil(t, channel)
	assert.Equal(t, "test-id", channel.ID)
	assert.Equal(t, "test-name", channel.Name)
	assert.Equal(t, restapi.SplunkChannelType, channel.Kind)
	assert.Equal(t, "https://splunk.example.com", *channel.URL)
	assert.Equal(t, "test-token", *channel.Token)
}

func TestMapVictorOpsChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	victorOpsModel := &shared.VictorOpsModel{
		APIKey:     types.StringValue("test-api-key"),
		RoutingKey: types.StringValue("test-routing-key"),
	}

	channel, diags := resource.mapVictorOpsChannelFromState(ctx, "test-id", "test-name", victorOpsModel)
	require.False(t, diags.HasError())
	require.NotNil(t, channel)
	assert.Equal(t, "test-id", channel.ID)
	assert.Equal(t, "test-name", channel.Name)
	assert.Equal(t, restapi.VictorOpsChannelType, channel.Kind)
	assert.Equal(t, "test-api-key", *channel.APIKey)
	assert.Equal(t, "test-routing-key", *channel.RoutingKey)
}

func TestMapWebhookChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("with headers", func(t *testing.T) {
		webhookURLs := []string{"https://webhook1.example.com", "https://webhook2.example.com"}
		webhookURLsSet, _ := types.SetValueFrom(ctx, types.StringType, webhookURLs)

		headers := map[string]string{
			"Authorization": "Bearer token",
			"Content-Type":  "application/json",
		}
		headersMap, _ := types.MapValueFrom(ctx, types.StringType, headers)

		webhookModel := &shared.WebhookModel{
			WebhookURLs: webhookURLsSet,
			HTTPHeaders: headersMap,
		}

		channel, diags := resource.mapWebhookChannelFromState(ctx, "test-id", "test-name", webhookModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, "test-id", channel.ID)
		assert.Equal(t, "test-name", channel.Name)
		assert.Equal(t, restapi.WebhookChannelType, channel.Kind)
		assert.Equal(t, webhookURLs, channel.WebhookURLs)
		assert.Len(t, channel.Headers, 2)
	})

	t.Run("without headers", func(t *testing.T) {
		webhookURLs := []string{"https://webhook.example.com"}
		webhookURLsSet, _ := types.SetValueFrom(ctx, types.StringType, webhookURLs)

		webhookModel := &shared.WebhookModel{
			WebhookURLs: webhookURLsSet,
			HTTPHeaders: types.MapNull(types.StringType),
		}

		channel, diags := resource.mapWebhookChannelFromState(ctx, "test-id", "test-name", webhookModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.WebhookChannelType, channel.Kind)
		assert.Nil(t, channel.Headers)
	})
}

func TestMapWebhookBasedChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	webhookBasedModel := &shared.WebhookBasedModel{
		WebhookURL: types.StringValue("https://webhook.example.com"),
	}

	t.Run("Office365", func(t *testing.T) {
		channel, diags := resource.mapWebhookBasedChannelFromState(ctx, "test-id", "test-name", webhookBasedModel, restapi.Office365ChannelType)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.Office365ChannelType, channel.Kind)
		assert.Equal(t, "https://webhook.example.com", *channel.WebhookURL)
	})

	t.Run("GoogleChat", func(t *testing.T) {
		channel, diags := resource.mapWebhookBasedChannelFromState(ctx, "test-id", "test-name", webhookBasedModel, restapi.GoogleChatChannelType)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.GoogleChatChannelType, channel.Kind)
	})

	t.Run("WebexTeamsWebhook", func(t *testing.T) {
		channel, diags := resource.mapWebhookBasedChannelFromState(ctx, "test-id", "test-name", webhookBasedModel, restapi.WebexTeamsWebhookChannelType)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.WebexTeamsWebhookChannelType, channel.Kind)
	})
}

func TestMapServiceNowChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("with all fields", func(t *testing.T) {
		autoClose := true
		serviceNowModel := &shared.ServiceNowModel{
			ServiceNowURL:      types.StringValue("https://servicenow.example.com"),
			Username:           types.StringValue("test-user"),
			Password:           types.StringValue("test-password"),
			AutoCloseIncidents: types.BoolValue(autoClose),
		}

		channel, diags := resource.mapServiceNowChannelFromState(ctx, "test-id", "test-name", serviceNowModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.ServiceNowChannelType, channel.Kind)
		assert.Equal(t, "https://servicenow.example.com", *channel.ServiceNowURL)
		assert.Equal(t, "test-user", *channel.Username)
		assert.Equal(t, "test-password", *channel.Password)
		assert.Equal(t, autoClose, *channel.AutoCloseIncidents)
	})

	t.Run("missing password", func(t *testing.T) {
		serviceNowModel := &shared.ServiceNowModel{
			ServiceNowURL:      types.StringValue("https://servicenow.example.com"),
			Username:           types.StringValue("test-user"),
			Password:           types.StringNull(),
			AutoCloseIncidents: types.BoolNull(),
		}

		channel, diags := resource.mapServiceNowChannelFromState(ctx, "test-id", "test-name", serviceNowModel)
		require.True(t, diags.HasError())
		assert.Nil(t, channel)
		assert.Contains(t, diags[0].Summary(), "Missing Password")
	})
}

func TestMapServiceNowApplicationChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("with all fields", func(t *testing.T) {
		serviceNowAppModel := &shared.ServiceNowApplicationModel{
			ServiceNowURL:                  types.StringValue("https://servicenow.example.com"),
			Username:                       types.StringValue("test-user"),
			Password:                       types.StringValue("test-password"),
			Tenant:                         types.StringValue("test-tenant"),
			Unit:                           types.StringValue("test-unit"),
			AutoCloseIncidents:             types.BoolValue(true),
			InstanaURL:                     types.StringValue("https://instana.example.com"),
			EnableSendInstanaNotes:         types.BoolValue(true),
			EnableSendServiceNowActivities: types.BoolValue(false),
			EnableSendServiceNowWorkNotes:  types.BoolValue(true),
			ManuallyClosedIncidents:        types.BoolValue(false),
			ResolutionOfIncident:           types.BoolValue(true),
			SnowStatusOnCloseEvent:         types.Int64Value(6),
		}

		channel, diags := resource.mapServiceNowApplicationChannelFromState(ctx, "test-id", "test-name", serviceNowAppModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.ServiceNowApplicationChannelType, channel.Kind)
		assert.Equal(t, "https://servicenow.example.com", *channel.ServiceNowURL)
		assert.Equal(t, "test-user", *channel.Username)
		assert.Equal(t, "test-password", *channel.Password)
		assert.Equal(t, "test-tenant", *channel.Tenant)
		assert.Equal(t, "test-unit", *channel.Unit)
		assert.Equal(t, "https://instana.example.com", *channel.InstanaURL)
		assert.Equal(t, true, *channel.EnableSendInstanaNotes)
		assert.Equal(t, 6, *channel.SnowStatusOnCloseEvent)
	})

	t.Run("missing password", func(t *testing.T) {
		serviceNowAppModel := &shared.ServiceNowApplicationModel{
			ServiceNowURL: types.StringValue("https://servicenow.example.com"),
			Username:      types.StringValue("test-user"),
			Password:      types.StringNull(),
			Tenant:        types.StringValue("test-tenant"),
			Unit:          types.StringValue("test-unit"),
			InstanaURL:    types.StringValue("https://instana.example.com"),
		}

		_, diags := resource.mapServiceNowApplicationChannelFromState(ctx, "test-id", "test-name", serviceNowAppModel)
		require.True(t, diags.HasError())
	})

	t.Run("missing InstanaURL", func(t *testing.T) {
		serviceNowAppModel := &shared.ServiceNowApplicationModel{
			ServiceNowURL: types.StringValue("https://servicenow.example.com"),
			Username:      types.StringValue("test-user"),
			Password:      types.StringValue("test-password"),
			Tenant:        types.StringValue("test-tenant"),
			Unit:          types.StringValue("test-unit"),
			InstanaURL:    types.StringNull(),
		}

		_, diags := resource.mapServiceNowApplicationChannelFromState(ctx, "test-id", "test-name", serviceNowAppModel)
		require.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), "InstanaURL is required")
	})
}

func TestMapPrometheusWebhookChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("with receiver", func(t *testing.T) {
		prometheusModel := &shared.PrometheusWebhookModel{
			WebhookURL: types.StringValue("https://prometheus.example.com"),
			Receiver:   types.StringValue("test-receiver"),
		}

		channel, diags := resource.mapPrometheusWebhookChannelFromState(ctx, "test-id", "test-name", prometheusModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.PrometheusWebhookChannelType, channel.Kind)
		assert.Equal(t, "https://prometheus.example.com", *channel.WebhookURL)
		assert.Equal(t, "test-receiver", *channel.Receiver)
	})

	t.Run("without receiver", func(t *testing.T) {
		prometheusModel := &shared.PrometheusWebhookModel{
			WebhookURL: types.StringValue("https://prometheus.example.com"),
			Receiver:   types.StringNull(),
		}

		channel, diags := resource.mapPrometheusWebhookChannelFromState(ctx, "test-id", "test-name", prometheusModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Nil(t, channel.Receiver)
	})
}

func TestMapWatsonAIOpsWebhookChannelFromState(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("with headers", func(t *testing.T) {
		headers := []string{"Authorization: Bearer token", "Content-Type: application/json"}
		headersList, _ := types.ListValueFrom(ctx, types.StringType, headers)

		watsonModel := &shared.WatsonAIOpsWebhookModel{
			WebhookURL:  types.StringValue("https://watson.example.com"),
			HTTPHeaders: headersList,
		}

		channel, diags := resource.mapWatsonAIOpsWebhookChannelFromState(ctx, "test-id", "test-name", watsonModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.WatsonAIOpsWebhookChannelType, channel.Kind)
		assert.Equal(t, "https://watson.example.com", *channel.WebhookURL)
		assert.Equal(t, headers, channel.Headers)
	})

	t.Run("without headers", func(t *testing.T) {
		watsonModel := &shared.WatsonAIOpsWebhookModel{
			WebhookURL:  types.StringValue("https://watson.example.com"),
			HTTPHeaders: types.ListNull(types.StringType),
		}

		channel, diags := resource.mapWatsonAIOpsWebhookChannelFromState(ctx, "test-id", "test-name", watsonModel)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Nil(t, channel.Headers)
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &alertingChannelResource{}
	ctx := context.Background()

	t.Run("Email channel", func(t *testing.T) {
		emails := []string{"test@example.com"}
		emailsSet, _ := types.SetValueFrom(ctx, types.StringType, emails)

		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			Email: &shared.EmailModel{
				Emails: emailsSet,
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.EmailChannelType, channel.Kind)
	})

	t.Run("OpsGenie channel", func(t *testing.T) {
		tags := []string{"tag1"}
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, tags)

		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			OpsGenie: &shared.OpsGenieModel{
				APIKey: types.StringValue("api-key"),
				Region: types.StringValue("EU"),
				Tags:   tagsList,
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.OpsGenieChannelType, channel.Kind)
	})

	t.Run("PagerDuty channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			PagerDuty: &shared.PagerDutyModel{
				ServiceIntegrationKey: types.StringValue("integration-key"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.PagerDutyChannelType, channel.Kind)
	})

	t.Run("Slack channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			Slack: &shared.SlackModel{
				WebhookURL: types.StringValue("https://slack.com/webhook"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.SlackChannelType, channel.Kind)
	})

	t.Run("Splunk channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			Splunk: &shared.SplunkModel{
				URL:   types.StringValue("https://splunk.com"),
				Token: types.StringValue("token"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.SplunkChannelType, channel.Kind)
	})

	t.Run("VictorOps channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			VictorOps: &shared.VictorOpsModel{
				APIKey:     types.StringValue("api-key"),
				RoutingKey: types.StringValue("routing-key"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.VictorOpsChannelType, channel.Kind)
	})

	t.Run("Webhook channel", func(t *testing.T) {
		urls := []string{"https://webhook.com"}
		urlsSet, _ := types.SetValueFrom(ctx, types.StringType, urls)

		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			Webhook: &shared.WebhookModel{
				WebhookURLs: urlsSet,
				HTTPHeaders: types.MapNull(types.StringType),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.WebhookChannelType, channel.Kind)
	})

	t.Run("Office365 channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			Office365: &shared.WebhookBasedModel{
				WebhookURL: types.StringValue("https://office365.com/webhook"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.Office365ChannelType, channel.Kind)
	})

	t.Run("GoogleChat channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			GoogleChat: &shared.WebhookBasedModel{
				WebhookURL: types.StringValue("https://chat.google.com/webhook"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.GoogleChatChannelType, channel.Kind)
	})

	t.Run("ServiceNow channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			ServiceNow: &shared.ServiceNowModel{
				ServiceNowURL: types.StringValue("https://servicenow.com"),
				Username:      types.StringValue("user"),
				Password:      types.StringValue("pass"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.ServiceNowChannelType, channel.Kind)
	})

	t.Run("ServiceNowApplication channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			ServiceNowApplication: &shared.ServiceNowApplicationModel{
				ServiceNowURL: types.StringValue("https://servicenow.com"),
				Username:      types.StringValue("user"),
				Password:      types.StringValue("pass"),
				Tenant:        types.StringValue("tenant"),
				Unit:          types.StringValue("unit"),
				InstanaURL:    types.StringValue("https://instana.com"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.ServiceNowApplicationChannelType, channel.Kind)
	})

	t.Run("PrometheusWebhook channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			PrometheusWebhook: &shared.PrometheusWebhookModel{
				WebhookURL: types.StringValue("https://prometheus.com/webhook"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.PrometheusWebhookChannelType, channel.Kind)
	})

	t.Run("WebexTeamsWebhook channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			WebexTeamsWebhook: &shared.WebhookBasedModel{
				WebhookURL: types.StringValue("https://webex.com/webhook"),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.WebexTeamsWebhookChannelType, channel.Kind)
	})

	t.Run("WatsonAIOpsWebhook channel", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
			WatsonAIOpsWebhook: &shared.WatsonAIOpsWebhookModel{
				WebhookURL:  types.StringValue("https://watson.com/webhook"),
				HTTPHeaders: types.ListNull(types.StringType),
			},
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, channel)
		assert.Equal(t, restapi.WatsonAIOpsWebhookChannelType, channel.Kind)
	})

	t.Run("No channel configured", func(t *testing.T) {
		model := AlertingChannelModel{
			ID:   types.StringValue("test-id"),
			Name: types.StringValue("test-name"),
		}

		state := createMockState(t, ctx, model)
		channel, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.True(t, diags.HasError())
		assert.Nil(t, channel)
		assert.Contains(t, diags[0].Summary(), "Invalid Alerting Channel Configuration")
	})
}

func TestUpdateState(t *testing.T) {
	ctx := context.Background()

	t.Run("Email channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		apiChannel := &restapi.AlertingChannel{
			ID:     "test-id",
			Name:   "test-name",
			Kind:   restapi.EmailChannelType,
			Emails: []string{"test@example.com"},
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.Equal(t, "test-name", model.Name.ValueString())
		assert.NotNil(t, model.Email)
	})

	t.Run("OpsGenie channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		apiKey := "test-api-key"
		region := "EU"
		tags := "tag1,tag2"
		apiChannel := &restapi.AlertingChannel{
			ID:     "test-id",
			Name:   "test-name",
			Kind:   restapi.OpsGenieChannelType,
			APIKey: &apiKey,
			Region: &region,
			Tags:   &tags,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.OpsGenie)
	})

	t.Run("PagerDuty channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		serviceKey := "service-key"
		apiChannel := &restapi.AlertingChannel{
			ID:                    "test-id",
			Name:                  "test-name",
			Kind:                  restapi.PagerDutyChannelType,
			ServiceIntegrationKey: &serviceKey,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.PagerDuty)
	})

	t.Run("Slack channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		webhookURL := "https://slack.com/webhook"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.SlackChannelType,
			WebhookURL: &webhookURL,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.Slack)
	})

	t.Run("Splunk channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		url := "https://splunk.com"
		token := "token"
		apiChannel := &restapi.AlertingChannel{
			ID:    "test-id",
			Name:  "test-name",
			Kind:  restapi.SplunkChannelType,
			URL:   &url,
			Token: &token,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.Splunk)
	})

	t.Run("VictorOps channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		apiKey := "api-key"
		routingKey := "routing-key"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.VictorOpsChannelType,
			APIKey:     &apiKey,
			RoutingKey: &routingKey,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.VictorOps)
	})

	t.Run("Webhook channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		apiChannel := &restapi.AlertingChannel{
			ID:          "test-id",
			Name:        "test-name",
			Kind:        restapi.WebhookChannelType,
			WebhookURLs: []string{"https://webhook.com"},
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.Webhook)
	})

	t.Run("Office365 channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		webhookURL := "https://office365.com/webhook"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.Office365ChannelType,
			WebhookURL: &webhookURL,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.Office365)
	})

	t.Run("GoogleChat channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		webhookURL := "https://chat.google.com/webhook"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.GoogleChatChannelType,
			WebhookURL: &webhookURL,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.GoogleChat)
	})

	t.Run("ServiceNow channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		serviceNowURL := "https://servicenow.com"
		username := "user"
		password := "pass"
		apiChannel := &restapi.AlertingChannel{
			ID:            "test-id",
			Name:          "test-name",
			Kind:          restapi.ServiceNowChannelType,
			ServiceNowURL: &serviceNowURL,
			Username:      &username,
			Password:      &password,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.ServiceNow)
	})

	t.Run("ServiceNowApplication channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		serviceNowURL := "https://servicenow.com"
		username := "user"
		password := "pass"
		tenant := "tenant"
		unit := "unit"
		apiChannel := &restapi.AlertingChannel{
			ID:            "test-id",
			Name:          "test-name",
			Kind:          restapi.ServiceNowApplicationChannelType,
			ServiceNowURL: &serviceNowURL,
			Username:      &username,
			Password:      &password,
			Tenant:        &tenant,
			Unit:          &unit,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.ServiceNowApplication)
	})

	t.Run("PrometheusWebhook channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		webhookURL := "https://prometheus.com/webhook"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.PrometheusWebhookChannelType,
			WebhookURL: &webhookURL,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.PrometheusWebhook)
	})

	t.Run("WebexTeamsWebhook channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		webhookURL := "https://webex.com/webhook"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.WebexTeamsWebhookChannelType,
			WebhookURL: &webhookURL,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.WebexTeamsWebhook)
	})

	t.Run("WatsonAIOpsWebhook channel", func(t *testing.T) {
		resource := &alertingChannelResource{}
		webhookURL := "https://watson.com/webhook"
		apiChannel := &restapi.AlertingChannel{
			ID:         "test-id",
			Name:       "test-name",
			Kind:       restapi.WatsonAIOpsWebhookChannelType,
			WebhookURL: &webhookURL,
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.False(t, diags.HasError())

		var model AlertingChannelModel
		diags = state.Get(ctx, &model)
		require.False(t, diags.HasError())
		assert.NotNil(t, model.WatsonAIOpsWebhook)
	})

	t.Run("Unsupported channel type", func(t *testing.T) {
		resource := &alertingChannelResource{}
		apiChannel := &restapi.AlertingChannel{
			ID:   "test-id",
			Name: "test-name",
			Kind: restapi.AlertingChannelType("UNSUPPORTED"),
		}

		handle := NewAlertingChannelResourceHandle()
		state := &tfsdk.State{
			Schema: handle.MetaData().Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiChannel)
		require.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary(), "Unsupported alerting channel type")
	})
}

// Helper function to create a mock state with a model
func createMockState(t *testing.T, ctx context.Context, model AlertingChannelModel) *tfsdk.State {
	handle := NewAlertingChannelResourceHandle()
	state := &tfsdk.State{
		Schema: handle.MetaData().Schema,
	}

	diags := state.Set(ctx, model)
	if diags.HasError() {
		t.Fatalf("Failed to set state: %v", diags)
	}

	return state
}
