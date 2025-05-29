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

const (
	customEventSpecId             = "customEventId1"
	customEventSpecName           = "eventName"
	customEventSpecDescription    = "custom event specification for test"
	customEventSpecEntityType     = "host"
	customEventSpecTriggering     = false
	customEventSpecEnabled        = true
	customEventSpecQuery          = "entity.zone:\"Neutral Zone\""
	customEventSpecExpirationTime = 10
)

type dataSourceCustomEventSpecificationUnitTest struct{}

func TestCustomEventSpecificationDataSource(t *testing.T) {
	unitTest := &dataSourceCustomEventSpecificationUnitTest{}
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("schema version should be 0", unitTest.schemaShouldHaveVersion0)
	t.Run("should successfully read custom event specification", unitTest.shouldSuccessfullyReadCustomEventSpecification)
	t.Run("should fail to read custom event specification when api call fails", unitTest.shouldFailToReadCustomEventSpecificationWhenApiCallFails)
	t.Run("should fail to custom event specification when no event found for name and entity type", unitTest.shouldFailToReadCustomEventSpecificationWhenNoEventIsFound)
}

func (r *dataSourceCustomEventSpecificationUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewCustomEventSpecificationDataSource().CreateResource().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 7)

	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsComputedAndOfTypeBool(CustomEventSpecificationFieldTriggering)
	schemaAssert.AssertSchemaIsComputedAndOfTypeBool(CustomEventSpecificationFieldEnabled)
	schemaAssert.AssertSchemaIsComputedAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
}

func (r *dataSourceCustomEventSpecificationUnitTest) schemaShouldHaveVersion0(t *testing.T) {
	require.Equal(t, 0, NewCustomEventSpecificationResourceHandle().MetaData().SchemaVersion)
}

func (r *dataSourceCustomEventSpecificationUnitTest) shouldSuccessfullyReadCustomEventSpecification(t *testing.T) {
	description := customEventSpecDescription
	query := customEventSpecQuery
	expirationTime := customEventSpecExpirationTime

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.CustomEventSpecification{
			ID:             customEventSpecId,
			Name:           customEventSpecName,
			EntityType:     customEventSpecEntityType,
			Description:    &description,
			Triggering:     customEventSpecTriggering,
			Query:          &query,
			ExpirationTime: &expirationTime,
			Enabled:        customEventSpecEnabled,
		}

		CustomEventSpecificationAPI := mocks.NewMockRestResource[*restapi.CustomEventSpecification](ctrl)
		CustomEventSpecificationAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.CustomEventSpecification{&data}, nil)
		mockInstanaApi.EXPECT().CustomEventSpecifications().Return(CustomEventSpecificationAPI).Times(1)

		sut := NewCustomEventSpecificationDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			CustomEventSpecificationFieldName:       customEventSpecName,
			CustomEventSpecificationFieldEntityType: customEventSpecEntityType,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.Nil(t, diag)
		require.Equal(t, data.ID, resourceData.Id())
		require.Equal(t, data.Name, resourceData.Get(CustomEventSpecificationFieldName))
		require.Equal(t, data.EntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
		require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
		require.Equal(t, query, resourceData.Get(CustomEventSpecificationFieldQuery))
		require.Equal(t, data.Triggering, resourceData.Get(CustomEventSpecificationFieldTriggering))
		require.Equal(t, expirationTime, resourceData.Get(CustomEventSpecificationFieldExpirationTime))
		require.Equal(t, data.Enabled, resourceData.Get(CustomEventSpecificationFieldEnabled))
	})
}

func (r *dataSourceCustomEventSpecificationUnitTest) shouldFailToReadCustomEventSpecificationWhenApiCallFails(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		expectedError := errors.New("test")

		CustomEventSpecificationAPI := mocks.NewMockRestResource[*restapi.CustomEventSpecification](ctrl)
		CustomEventSpecificationAPI.EXPECT().GetAll().Times(1).Return(nil, expectedError)
		mockInstanaApi.EXPECT().CustomEventSpecifications().Return(CustomEventSpecificationAPI).Times(1)

		sut := NewCustomEventSpecificationDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			CustomEventSpecificationFieldName:       customEventSpecName,
			CustomEventSpecificationFieldEntityType: customEventSpecEntityType,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, expectedError.Error())
	})
}

func (r *dataSourceCustomEventSpecificationUnitTest) shouldFailToReadCustomEventSpecificationWhenNoEventIsFound(t *testing.T) {
	description := customEventSpecDescription
	query := customEventSpecQuery
	expirationTime := customEventSpecExpirationTime

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.CustomEventSpecification{
			ID:             customEventSpecId,
			Name:           customEventSpecName,
			EntityType:     customEventSpecEntityType,
			Description:    &description,
			Triggering:     customEventSpecTriggering,
			Query:          &query,
			ExpirationTime: &expirationTime,
			Enabled:        customEventSpecEnabled,
		}

		CustomEventSpecificationAPI := mocks.NewMockRestResource[*restapi.CustomEventSpecification](ctrl)
		CustomEventSpecificationAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.CustomEventSpecification{&data}, nil)
		mockInstanaApi.EXPECT().CustomEventSpecifications().Return(CustomEventSpecificationAPI).Times(1)

		sut := NewCustomEventSpecificationDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			CustomEventSpecificationFieldName:       "customEvent2",
			CustomEventSpecificationFieldEntityType: customEventSpecEntityType,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, "no custom event specification found for name 'customEvent2' and entity type 'host'")
	})
}
