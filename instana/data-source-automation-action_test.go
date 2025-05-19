package instana_test

import (
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

type dataSourceAutomationActionUnitTest struct{}

func TestAutomationActionDataSource(t *testing.T) {
	unitTest := &dataSourceAutomationActionUnitTest{}
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("schema version should be 0", unitTest.schemaShouldHaveVersion0)
	t.Run("should successfully read automation action", unitTest.shouldSuccessfullyReadAction)
	t.Run("should fail to read automation action when api call fails", unitTest.shouldFailToReadActionWhenApiCallFails)
	t.Run("should fail to read automation action when no action found for name and type", unitTest.shouldFailToReadActionWhenNoActionIsFound)
}

func (r *dataSourceAutomationActionUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewAutomationActionDataSource().CreateResource().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 4)

	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AutomationActionFieldType)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AutomationActionFieldDescription)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfStrings(AutomationActionFieldTags)
}

func (r *dataSourceAutomationActionUnitTest) schemaShouldHaveVersion0(t *testing.T) {
	require.Equal(t, 0, NewAutomationActionResourceHandle().MetaData().SchemaVersion)
}

func (r *dataSourceAutomationActionUnitTest) shouldSuccessfullyReadAction(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.AutomationAction{
			ID:          "id1",
			Name:        "action4test",
			Type:        "SCRIPT",
			Description: "test automation action",
			Tags:        []string{"test", "development"},
		}

		AutomationActionAPI := mocks.NewMockRestResource[*restapi.AutomationAction](ctrl)
		AutomationActionAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.AutomationAction{&data}, nil)
		mockInstanaApi.EXPECT().AutomationActions().Return(AutomationActionAPI).Times(1)

		sut := NewAutomationActionDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AutomationActionFieldName: "action4test",
			AutomationActionFieldType: "SCRIPT",
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.Nil(t, diag)
		require.Equal(t, data.ID, resourceData.Id())
		require.Equal(t, data.Name, resourceData.Get(AutomationActionFieldName))
		require.Equal(t, data.Type, resourceData.Get(AutomationActionFieldType))
		require.Equal(t, data.Description, resourceData.Get(AutomationActionFieldDescription))
		require.IsType(t, []interface{}{}, resourceData.Get(AutomationActionFieldTags))
		require.Len(t, resourceData.Get(AutomationActionFieldTags).([]interface{}), 2)
		tags := resourceData.Get(AutomationActionFieldTags).([]interface{})
		require.Len(t, tags, 2)
		require.Contains(t, tags, "test")
		require.Contains(t, tags, "development")
	})
}

func (r *dataSourceAutomationActionUnitTest) shouldFailToReadActionWhenApiCallFails(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		expectedError := errors.New("test")

		AutomationActionAPI := mocks.NewMockRestResource[*restapi.AutomationAction](ctrl)
		AutomationActionAPI.EXPECT().GetAll().Times(1).Return(nil, expectedError)
		mockInstanaApi.EXPECT().AutomationActions().Return(AutomationActionAPI).Times(1)

		sut := NewAutomationActionDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AutomationActionFieldName: "action4test",
			AutomationActionFieldType: "SCRIPT",
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, expectedError.Error())
	})
}

func (r *dataSourceAutomationActionUnitTest) shouldFailToReadActionWhenNoActionIsFound(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AutomationAction](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.AutomationAction{
			ID:          "id1",
			Name:        "action4test",
			Type:        "HTTP",
			Description: "test automation action",
			Tags:        []string{"test", "development"},
		}

		AutomationActionAPI := mocks.NewMockRestResource[*restapi.AutomationAction](ctrl)
		AutomationActionAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.AutomationAction{&data}, nil)
		mockInstanaApi.EXPECT().AutomationActions().Return(AutomationActionAPI).Times(1)

		sut := NewAutomationActionDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AutomationActionFieldName: "action4test",
			AutomationActionFieldType: "SCRIPT",
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, "no automation action found for name 'action4test' and type 'SCRIPT'")
	})
}
