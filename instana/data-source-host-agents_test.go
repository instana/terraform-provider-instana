package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

type dataSourceHostAgentsUnitTest struct{}

func TestHostAgentsDataSource(t *testing.T) {
	unitTest := &dataSourceHostAgentsUnitTest{}
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("schema version should be 0", unitTest.schemaShouldHaveVersion0)
	t.Run("should successfully read host agents", unitTest.shouldSuccessfullyReadHostAgents)
}

func (ut *dataSourceHostAgentsUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewHostAgentsDataSource().CreateResource().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 2)

	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(HostAgentFieldFilter)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(HostAgentFieldItems)

	agentSchema := schemaData[HostAgentFieldItems].Elem.(*schema.Resource).Schema
	require.Len(t, agentSchema, 5)

	schemaAssert = testutils.NewTerraformSchemaAssert(agentSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(HostAgentFieldSnapshotId)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(HostAgentFieldLabel)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(HostAgentFieldHost)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(HostAgentFieldPlugin)
}

func (ut *dataSourceHostAgentsUnitTest) schemaShouldHaveVersion0(t *testing.T) {
	require.Equal(t, 0, NewAutomationActionResourceHandle().MetaData().SchemaVersion)
}

func (ut *dataSourceHostAgentsUnitTest) shouldSuccessfullyReadHostAgents(t *testing.T) {
	testHelper := NewTestHelper[*restapi.HostAgent](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.HostAgent{
			SnapshotID: "spanshotId1",
			Label:      "Instana agent @localhost",
			Host:       "hostId1",
			Plugin:     "instanaAgent",
			Tags:       []string{"test", "development"},
		}

		expectedQueryParams := make(map[string]string)
		expectedQueryParams["query"] = "entity.agent.capability:action-script"

		HostAgentAPI := mocks.NewMockReadOnlyRestResource[*restapi.HostAgent](ctrl)
		HostAgentAPI.EXPECT().GetByQuery(expectedQueryParams).Times(1).Return(&[]*restapi.HostAgent{&data}, nil)
		mockInstanaApi.EXPECT().HostAgents().Return(HostAgentAPI).Times(1)

		sut := NewHostAgentsDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			HostAgentFieldFilter: "entity.agent.capability:action-script",
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.Nil(t, diag)
		require.NotNil(t, resourceData.Id())

		agent := resourceData.Get(HostAgentFieldItems).([]interface{})[0].(map[string]interface{})
		require.Equal(t, data.SnapshotID, agent[HostAgentFieldSnapshotId])
		require.Equal(t, data.Label, agent[HostAgentFieldLabel])
		require.Equal(t, data.Host, agent[HostAgentFieldHost])
		require.Equal(t, data.Plugin, agent[HostAgentFieldPlugin])
		require.IsType(t, []interface{}{}, agent[HostAgentFieldTags])
		require.Len(t, agent[HostAgentFieldTags].([]interface{}), 2)
		tags := agent[HostAgentFieldTags].([]interface{})
		require.Len(t, tags, 2)
		require.Contains(t, tags, "test")
		require.Contains(t, tags, "development")
	})
}
