package instana_test

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestSliConfigTest(t *testing.T) {
	unitTest := &sliConfigUnitTest{}
	t.Run("CRUD integration test of with SLI Entity of type applicationTimeBased", sliConfigIntegrationTestWithSliEntityOfTypeApplicationTimeBased().testCRUD())
	t.Run("CRUD integration test of with SLI Entity of type applicationEventBased", sliConfigIntegrationTestWithSliEntityOfTypeApplicationEventBased().testCRUD())
	t.Run("CRUD integration test of with SLI Entity of type WebsiteEventBased", sliConfigIntegrationTestWithSliEntityOfTypeWebsiteEventBased().testCRUD())
	t.Run("CRUD integration test of with SLI Entity of type WebsiteTimeBased", sliConfigIntegrationTestWithSliEntityOfTypeWebsiteTimeBased().testCRUD())
	t.Run("should have valid resource schema", unitTest.shouldHaveValidResourceSchema())
	t.Run("should return correct resource name", unitTest.shouldReturnCorrectResourceNameForSliConfigs())
	t.Run("should have schema version one", unitTest.shouldHaveSchemaVersionOne())
	t.Run("should have on schema state upgrader", unitTest.shouldHaveOneStateUpgrader())
	t.Run("should migrate full name to name when executing first state upgrader and full name is available", unitTest.shouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable())
	t.Run("should migrate do nothing when executing first state upgrader and full name is not available", unitTest.shouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsNotAvailable())
	t.Run("should failt to update resource state for SLI Config with entity type is not supported", unitTest.shouldFailToUpdateResourceStateWhenEntityTypeIsNotSupported())
	t.Run("should update resource state for SLI Config with SLI Entity of type Application", unitTest.shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeApplication())
	t.Run("should update resource state for SLI Config with SLI Entity of type Availability", unitTest.shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeAvailability())
	t.Run("should fail to update resource state for SLI Config with SLI Entity of type Availability when good event filter expression is invalid", unitTest.shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeAvailabilityWhenGoodEventFilterExpressionIsNotValid())
	t.Run("should fail to update resource state for SLI Config with SLI Entity of type Availability when bad event filter expression is invalid", unitTest.shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeAvailabilityWhenBadEventFilterExpressionIsNotValid())
	t.Run("should update resource state for SLI Config with SLI Entity of type WebsiteTimeBased", unitTest.shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeWebsiteTimeBased())
	t.Run("should fail to update resource state for SLI Config with SLI Entity of type WebsiteTimeBased when filter expression is invalid", unitTest.shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeWebsiteTimeBasedWhenFilterExpressionIsNotValid())
	t.Run("should update resource state for SLI Config with SLI Entity of type WebsiteEventBased", unitTest.shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeWebsiteEventBased())
	t.Run("should fail to update resource state for SLI Config with SLI Entity of type WebsiteEventBased when good event filter expression is invalid", unitTest.shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeWebsiteEventBasedWhenGoodEventFilterExpressionIsNotValid())
	t.Run("should fail to update resource state for SLI Config with SLI Entity of type WebsiteEventBased when bad event filter expression is invalid", unitTest.shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeWebsiteEventBasedWhenBadEventFilterExpressionIsNotValid())
	t.Run("should map state of SLI Config with SLI Entity of type Application to Data Model", unitTest.shouldMapStateOfSliConfigWithEntityOfTypeApplicationToDataModel())
	t.Run("should map state of SLI Config with SLI Entity of type Availability to Data Model", unitTest.shouldMapStateOfSliConfigWithEntityOfTypeAvailabilityToDataModel())
	t.Run("should fail to map state of SLI Config with SLI Entity of type Availability to Data Model when good event filter expression is invalid", unitTest.shouldFailToMapStateOfSliConfigWithEntityOfTypeAvailabilityToDataModelWhenGoodEventFilterExpressionIsInvalid())
	t.Run("should fail to map state of SLI Config with SLI Entity of type Availability to Data Model when bad event filter expression is invalid", unitTest.shouldFailToMapStateOfSliConfigWithEntityOfTypeAvailabilityToDataModelWhenBadEventFilterExpressionIsInvalid())
	t.Run("should map state of SLI Config with SLI Entity of type WebsiteEventBased to Data Model", unitTest.shouldMapStateOfSliConfigWithEntityOfTypeWebsiteEventBasedToDataModel())
	t.Run("should fail to map state of SLI Config with SLI Entity of type WebsiteEventBased to Data Model when good event filter expression is invalid", unitTest.shouldFailToMapStateOfSliConfigWithEntityOfTypeWebsiteEventBasedToDataModelWhenGoodEventFilterExpressionIsInvalid())
	t.Run("should fail to map state of SLI Config with SLI Entity of type WebsiteEventBased to Data Model when bad event filter expression is invalid", unitTest.shouldFailToMapStateOfSliConfigWithEntityOfTypeWebsiteEventBasedToDataModelWhenBadEventFilterExpressionIsInvalid())
	t.Run("should map state of SLI Config with SLI Entity of type WebsiteTimeBased to Data Model", unitTest.shouldMapStateOfSliConfigWithEntityOfTypeWebsiteTimeBasedToDataModel())
	t.Run("should fail to map state of SLI Config with SLI Entity of type WebsiteTimeBased to Data Model when filter expression is invalid", unitTest.shouldFailToMapStateOfSliConfigWithEntityOfTypeWebsiteTimeBasedToDataModelWhenFilterExpressionIsInvalid())
	t.Run("should fail to map state of SLI Config with SLI Entity when no Sli Entity is provided", unitTest.shouldFailToMapStateOfSliConfigWhenNoSliEntityIsProvided())
	t.Run("should fail to map state of SLI Config with SLI Entity when no Sli Entity Data is provided", unitTest.shouldFailToMapStateOfSliConfigWhenNoSupportedliEntityIsProvided())
	t.Run("should require metric threshold to be greater than 0", unitTest.shouldRequireMetricConfigurationThresholdToBeGreaterThanZero())
}

const (
	sliConfigDefinition = "instana_sli_config.example_sli_config"

	sliMetricResourceFieldPattern = "%s.0.%s"
	sliEntityResourceFieldPattern = "%s.0.%s.0.%s"

	sliConfigID                               = "id"
	sliConfigName                             = resourceName
	sliConfigInitialEvaluationTimestamp       = 0
	sliConfigMetricName                       = "metric_name"
	sliConfigMetricAggregation                = "SUM"
	sliConfigMetricThreshold                  = 1.0
	sliConfigEntityApplicationID              = "application_id"
	sliConfigEntityWebsiteID                  = "website_id"
	sliConfigEntityServiceID                  = "service_id"
	sliConfigEntityEndpointID                 = "endpoint_id"
	sliConfigEntityBoundaryScope              = "ALL"
	sliConfigEntityBeaconType                 = "pageLoad"
	sliConfigTagFilterExpressionString        = "request.path@dest EQUALS '/home'"
	invalidSliConfigTagFilterExpressionString = "request.path@dest EQUALS"
)

var (
	tagFilterEntityDestination         = restapi.TagFilterEntityDestination
	invalidTagFilterExpressionName     = "foo"
	invalidTagFilterExpressionOperator = restapi.EqualsOperator
	invalidTagFilterExpression         = &restapi.TagFilter{
		Entity:   &tagFilterEntityDestination,
		Name:     &invalidTagFilterExpressionName,
		Operator: &invalidTagFilterExpressionOperator,
		Type:     "invalid",
	}

	sliConfigTagFilterExpression = restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "/home")
)

func sliConfigIntegrationTestWithSliEntityOfTypeApplicationTimeBased() *sliConfigIntegrationTest {
	resourceTemplate := `
resource "instana_sli_config" "example_sli_config" {
	name = "name %d"
	initial_evaluation_timestamp = 0
	metric_configuration {
		metric_name = "metric_name"
		aggregation = "SUM"
		threshold = 1.0
	}
	sli_entity {
		application_time_based {
			application_id = "application_id"
			service_id     = "service_id"
			endpoint_id    = "endpoint_id"
			boundary_scope = "ALL"
		}
	}
}
`
	serverResponseTemplate := `
{
	"id" : "%s",
	"sliName" : "name %d",
	"initialEvaluationTimestamp": 0,
	"metricConfiguration": {
		"metricName" : "metric_name",
		"metricAggregation"	: "SUM",
		"threshold" : 1.0
	},
	"sliEntity": {
		"sliType" : "application",
		"applicationId"	: "application_id",
		"serviceId" : "service_id",
		"endpointId" : "endpoint_id",
		"boundaryScope"	: "ALL"
	}
}
`
	useCaseSpecificChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), sliConfigMetricName),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), sliConfigMetricAggregation),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationTimeBased, SliConfigFieldApplicationID), sliConfigEntityApplicationID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationTimeBased, SliConfigFieldServiceID), sliConfigEntityServiceID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationTimeBased, SliConfigFieldEndpointID), sliConfigEntityEndpointID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationTimeBased, SliConfigFieldBoundaryScope), sliConfigEntityBoundaryScope),
	}
	return newSliConfigIntegrationTest(resourceTemplate, serverResponseTemplate, useCaseSpecificChecks)
}

func sliConfigIntegrationTestWithSliEntityOfTypeApplicationEventBased() *sliConfigIntegrationTest {
	resourceTemplate := `
resource "instana_sli_config" "example_sli_config" {
	name = "name %d"
	initial_evaluation_timestamp = 0
	sli_entity {
		application_event_based {
			application_id               = "application_id"
			boundary_scope               = "ALL"
			include_internal             = true
			include_synthetic            = true
			good_event_filter_expression = "request.path@dest EQUALS '/home'"
			bad_event_filter_expression  = "request.path@dest EQUALS '/404'"
		}
	}
}
`
	serverResponseTemplate := `
{
	"id" : "%s",
	"sliName" : "name %d",
	"initialEvaluationTimestamp": 0,
	"sliEntity": {
		"sliType" : "availability",
		"applicationId" : "application_id",
		"boundaryScope" : "ALL",
        "includeInternal" : true,
        "includeSynthetic" : true,
		"goodEventFilterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/home",
			"value" : "/home"
		},
		"badEventFilterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/404",
			"value" : "/404"
		}
	}
}
`
	useCaseSpecificChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationEventBased, SliConfigFieldApplicationID), sliConfigEntityApplicationID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationEventBased, SliConfigFieldBoundaryScope), sliConfigEntityBoundaryScope),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationEventBased, SliConfigFieldIncludeSynthetic), "true"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationEventBased, SliConfigFieldIncludeInternal), "true"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationEventBased, SliConfigFieldGoodEventFilterExpression), "request.path@dest EQUALS '/home'"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplicationEventBased, SliConfigFieldBadEventFilterExpression), "request.path@dest EQUALS '/404'"),
	}
	return newSliConfigIntegrationTest(resourceTemplate, serverResponseTemplate, useCaseSpecificChecks)
}

func sliConfigIntegrationTestWithSliEntityOfTypeWebsiteEventBased() *sliConfigIntegrationTest {
	resourceTemplate := `
resource "instana_sli_config" "example_sli_config" {
	name = "name %d"
	initial_evaluation_timestamp = 0
	sli_entity {
		website_event_based {
			website_id                   = "website_id"
			beacon_type                  = "pageLoad"
			good_event_filter_expression = "request.path@dest EQUALS '/home'"
			bad_event_filter_expression  = "request.path@dest EQUALS '/404'"
		}
	}
}
`
	serverResponseTemplate := `
{
	"id" : "%s",
	"sliName" : "name %d",
	"initialEvaluationTimestamp": 0,
	"sliEntity": {
		"sliType" : "websiteEventBased",
		"websiteId" : "website_id",
		"beaconType" : "pageLoad",
		"goodEventFilterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/home",
			"value" : "/home"
		},
		"badEventFilterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/404",
			"value" : "/404"
		}
	}
}
`
	useCaseSpecificChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteEventBased, SliConfigFieldWebsiteID), sliConfigEntityWebsiteID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteEventBased, SliConfigFieldBeaconType), sliConfigEntityBeaconType),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteEventBased, SliConfigFieldGoodEventFilterExpression), "request.path@dest EQUALS '/home'"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteEventBased, SliConfigFieldBadEventFilterExpression), "request.path@dest EQUALS '/404'"),
	}
	return newSliConfigIntegrationTest(resourceTemplate, serverResponseTemplate, useCaseSpecificChecks)
}

func sliConfigIntegrationTestWithSliEntityOfTypeWebsiteTimeBased() *sliConfigIntegrationTest {
	resourceTemplate := `
resource "instana_sli_config" "example_sli_config" {
	name = "name %d"
	initial_evaluation_timestamp = 0
	metric_configuration {
		metric_name = "metric_name"
		aggregation = "SUM"
		threshold = 1.0
	}
	sli_entity {
		website_time_based {
			website_id        = "website_id"
			beacon_type       = "pageLoad"
			filter_expression = "request.path@dest EQUALS '/home'"
		}
	}
}
`
	serverResponseTemplate := `
{
	"id" : "%s",
	"sliName" : "name %d",
	"initialEvaluationTimestamp": 0,
	"metricConfiguration": {
		"metricName" : "metric_name",
		"metricAggregation"	: "SUM",
		"threshold"	 : 1.0
	},
	"sliEntity": {
		"sliType" : "websiteTimeBased",
		"websiteId" : "website_id",
		"beaconType" : "pageLoad",
		"filterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/home",
			"value" : "/home"
		}
	}
}
`
	useCaseSpecificChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), sliConfigMetricName),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), sliConfigMetricAggregation),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteTimeBased, SliConfigFieldWebsiteID), sliConfigEntityWebsiteID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteTimeBased, SliConfigFieldBeaconType), sliConfigEntityBeaconType),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityWebsiteTimeBased, SliConfigFieldFilterExpression), "request.path@dest EQUALS '/home'"),
	}
	return newSliConfigIntegrationTest(resourceTemplate, serverResponseTemplate, useCaseSpecificChecks)
}

func newSliConfigIntegrationTest(resourceTemplate string, serverResponseTemplate string, useCaseSpecificChecks []resource.TestCheckFunc) *sliConfigIntegrationTest {
	return &sliConfigIntegrationTest{
		resourceTemplate:       resourceTemplate,
		serverResponseTemplate: serverResponseTemplate,
		useCaseSpecificChecks:  useCaseSpecificChecks,
	}
}

type sliConfigIntegrationTest struct {
	resourceTemplate       string
	serverResponseTemplate string
	useCaseSpecificChecks  []resource.TestCheckFunc
}

func (r *sliConfigIntegrationTest) testCRUD() func(t *testing.T) {
	return func(t *testing.T) {
		resourcePath := restapi.SliConfigResourcePath
		serverResponseTemplate := r.serverResponseTemplate
		pathTemplate := resourcePath + "/{id}"
		httpServer := testutils.NewTestHTTPServer()
		id := RandomID()
		responseHandler := func(w http.ResponseWriter, r *http.Request) {
			callCount := getZeroBasedCallCount(httpServer, http.MethodPost, resourcePath)
			json := formatResponseTemplate(serverResponseTemplate, id, callCount)
			w.Header().Set(contentType, r.Header.Get(contentType))
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(json))
			if err != nil {
				log.Fatalf("failed to write response: %s", err)
			}
		}
		httpServer.AddRoute(http.MethodPost, resourcePath, responseHandler)
		httpServer.AddRoute(http.MethodDelete, pathTemplate, responseHandler)
		httpServer.AddRoute(http.MethodGet, pathTemplate, responseHandler)
		httpServer.Start()
		defer httpServer.Close()

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: testProviderFactory,
			Steps: []resource.TestStep{
				r.createTestCheckFunction(httpServer.GetPort(), 0, id),
				testStepImport(sliConfigDefinition),
				r.createUpdateNotSupportedTestCheckFunction(httpServer.GetPort(), 1),
			},
		})
	}
}

func (r *sliConfigIntegrationTest) createTestCheckFunction(httpPort int, iteration int, id string) resource.TestStep {
	defaultChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
		resource.TestCheckResourceAttr(sliConfigDefinition, "id", id),
		resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, formatResourceName(iteration)),
		resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldInitialEvaluationTimestamp, "0"),
	}
	checks := append(defaultChecks, r.useCaseSpecificChecks...)
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(r.resourceTemplate, iteration), httpPort),
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}

func (r *sliConfigIntegrationTest) createUpdateNotSupportedTestCheckFunction(httpPort int, iteration int) resource.TestStep {
	return resource.TestStep{
		Config:      appendProviderConfig(fmt.Sprintf(r.resourceTemplate, iteration), httpPort),
		ExpectError: regexp.MustCompile("update operations not supported for instana_sli_config resources"),
	}
}

type sliConfigUnitTest struct{}

func (r *sliConfigUnitTest) shouldHaveValidResourceSchema() func(t *testing.T) {
	return func(t *testing.T) {
		resourceHandle := NewSliConfigResourceHandle()

		schemaMap := resourceHandle.MetaData().Schema

		schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
		schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldName)
		schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SliConfigFieldInitialEvaluationTimestamp)
		schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfResource(SliConfigFieldSliEntity)

		r.validateMetricsConfig(t, schemaMap)

		r.validateSliEntity(t, schemaMap)
	}
}

func (r *sliConfigUnitTest) validateSliEntity(t *testing.T, schemaMap map[string]*schema.Schema) {
	sliEntitySchemaMap := schemaMap[SliConfigFieldSliEntity].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntitySchemaMap, 4)
	schemaAssert := testutils.NewTerraformSchemaAssert(sliEntitySchemaMap, t)

	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityApplicationTimeBased)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityApplicationEventBased)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityWebsiteEventBased)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityWebsiteTimeBased)

	sliEntityApplicationSchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityApplicationTimeBased].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityApplicationSchemaMap, 4)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityApplicationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBoundaryScope)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldApplicationID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldServiceID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldEndpointID)

	sliEntityAvailabilitySchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityApplicationEventBased].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityAvailabilitySchemaMap, 6)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityAvailabilitySchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldApplicationID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBoundaryScope)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBadEventFilterExpression)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldGoodEventFilterExpression)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SliConfigFieldIncludeInternal, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SliConfigFieldIncludeSynthetic, false)

	sliEntityWebsiteEventBasedSchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityWebsiteEventBased].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityWebsiteEventBasedSchemaMap, 4)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityWebsiteEventBasedSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldWebsiteID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBeaconType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBadEventFilterExpression)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldGoodEventFilterExpression)

	sliEntityWebsiteTimeBasedSchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityWebsiteTimeBased].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityWebsiteTimeBasedSchemaMap, 3)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityWebsiteTimeBasedSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldWebsiteID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBeaconType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldFilterExpression)
}

func (r *sliConfigUnitTest) validateMetricsConfig(t *testing.T, schemaMap map[string]*schema.Schema) {
	metricConfigurationSchemaMap := schemaMap[SliConfigFieldMetricConfiguration].Elem.(*schema.Resource).Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(metricConfigurationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldMetricName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldMetricAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(SliConfigFieldMetricThreshold)
}

func (r *sliConfigUnitTest) shouldReturnCorrectResourceNameForSliConfigs() func(t *testing.T) {
	return func(t *testing.T) {
		name := NewSliConfigResourceHandle().MetaData().ResourceName

		require.Equal(t, "instana_sli_config", name, "Expected resource name to be instana_sli_config")
	}
}

func (r *sliConfigUnitTest) shouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, NewSliConfigResourceHandle().MetaData().SchemaVersion)
	}
}

func (r *sliConfigUnitTest) shouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, len(NewSliConfigResourceHandle().StateUpgraders()))
	}
}

func (r *sliConfigUnitTest) shouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"full_name": "test",
		}
		result, err := NewSliConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.NotContains(t, result, SliConfigFieldFullName)
		require.Contains(t, result, SliConfigFieldName)
		require.Equal(t, "test", result[SliConfigFieldName])
	}
}

func (r *sliConfigUnitTest) shouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsNotAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"name": "test",
		}
		result, err := NewSliConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Equal(t, input, result)
	}
}

func (r *sliConfigUnitTest) shouldFailToUpdateResourceStateWhenEntityTypeIsNotSupported() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type: "test",
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported sli entity type test")
	}
}

func (r *sliConfigUnitTest) shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeApplication() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		applicationId := "my-application"
		serviceId := "my-service"
		endpointId := "my-endpint"
		boundaryScope := "INBOUND"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:          "application",
				ApplicationID: &applicationId,
				ServiceID:     &serviceId,
				EndpointID:    &endpointId,
				BoundaryScope: &boundaryScope,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, sliConfigID, resourceData.Id())
		require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
		require.IsType(t, []interface{}{}, resourceData.Get(SliConfigFieldSliEntity))
		sliEntitySlice := resourceData.Get(SliConfigFieldSliEntity).([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntitySlice[0])
		sliEntityData := sliEntitySlice[0].(map[string]interface{})
		require.IsType(t, []interface{}{}, sliEntityData[SliConfigFieldSliEntityApplicationTimeBased])
		sliEntityApplicationSlice := sliEntityData[SliConfigFieldSliEntityApplicationTimeBased].([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntityApplicationSlice[0])
		sliEntityApplicationData := sliEntityApplicationSlice[0].(map[string]interface{})
		require.Equal(t, applicationId, sliEntityApplicationData[SliConfigFieldApplicationID])
		require.Equal(t, serviceId, sliEntityApplicationData[SliConfigFieldServiceID])
		require.Equal(t, endpointId, sliEntityApplicationData[SliConfigFieldEndpointID])
		require.Equal(t, boundaryScope, sliEntityApplicationData[SliConfigFieldBoundaryScope])
	}
}

func (r *sliConfigUnitTest) shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeAvailability() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		applicationId := "my-application"
		includeInternal := true
		includeSynthetic := false
		goodEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "good")
		badEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "bad")
		boundaryScope := "INBOUND"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:                      "availability",
				ApplicationID:             &applicationId,
				BoundaryScope:             &boundaryScope,
				IncludeSynthetic:          &includeSynthetic,
				IncludeInternal:           &includeInternal,
				GoodEventFilterExpression: goodEventTagFilterExpression,
				BadEventFilterExpression:  badEventTagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, sliConfigID, resourceData.Id())
		require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
		require.IsType(t, []interface{}{}, resourceData.Get(SliConfigFieldSliEntity))
		sliEntitySlice := resourceData.Get(SliConfigFieldSliEntity).([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntitySlice[0])
		sliEntityData := sliEntitySlice[0].(map[string]interface{})
		require.IsType(t, []interface{}{}, sliEntityData[SliConfigFieldSliEntityApplicationEventBased])
		sliEntityAvailabilitySlice := sliEntityData[SliConfigFieldSliEntityApplicationEventBased].([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntityAvailabilitySlice[0])
		sliEntityAvailabilityData := sliEntityAvailabilitySlice[0].(map[string]interface{})
		require.Equal(t, applicationId, sliEntityAvailabilityData[SliConfigFieldApplicationID])
		require.Equal(t, boundaryScope, sliEntityAvailabilityData[SliConfigFieldBoundaryScope])
		require.Equal(t, includeInternal, sliEntityAvailabilityData[SliConfigFieldIncludeInternal])
		require.Equal(t, includeSynthetic, sliEntityAvailabilityData[SliConfigFieldIncludeSynthetic])
		require.Equal(t, "request.path@dest EQUALS 'good'", sliEntityAvailabilityData[SliConfigFieldGoodEventFilterExpression])
		require.Equal(t, "request.path@dest EQUALS 'bad'", sliEntityAvailabilityData[SliConfigFieldBadEventFilterExpression])
	}
}

func (r *sliConfigUnitTest) shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeAvailabilityWhenGoodEventFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		applicationId := "my-application"
		includeInternal := true
		includeSynthetic := false
		goodEventTagFilterExpression := invalidTagFilterExpression
		badEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "bad")
		boundaryScope := "INBOUND"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:                      "availability",
				ApplicationID:             &applicationId,
				BoundaryScope:             &boundaryScope,
				IncludeSynthetic:          &includeSynthetic,
				IncludeInternal:           &includeInternal,
				GoodEventFilterExpression: goodEventTagFilterExpression,
				BadEventFilterExpression:  badEventTagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported tag filter expression of type invalid")
	}
}

func (r *sliConfigUnitTest) shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeAvailabilityWhenBadEventFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		applicationId := "my-application"
		includeInternal := true
		includeSynthetic := false
		goodEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "good")
		badEventTagFilterExpression := invalidTagFilterExpression
		boundaryScope := "INBOUND"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:                      "availability",
				ApplicationID:             &applicationId,
				BoundaryScope:             &boundaryScope,
				IncludeSynthetic:          &includeSynthetic,
				IncludeInternal:           &includeInternal,
				GoodEventFilterExpression: goodEventTagFilterExpression,
				BadEventFilterExpression:  badEventTagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported tag filter expression of type invalid")
	}
}

func (r *sliConfigUnitTest) shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeWebsiteEventBased() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		websiteId := "my-website"
		goodEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "good")
		badEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "bad")
		beaconType := "pageLoad"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:                      "websiteEventBased",
				WebsiteId:                 &websiteId,
				BeaconType:                &beaconType,
				GoodEventFilterExpression: goodEventTagFilterExpression,
				BadEventFilterExpression:  badEventTagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, sliConfigID, resourceData.Id())
		require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
		require.IsType(t, []interface{}{}, resourceData.Get(SliConfigFieldSliEntity))
		sliEntitySlice := resourceData.Get(SliConfigFieldSliEntity).([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntitySlice[0])
		sliEntityData := sliEntitySlice[0].(map[string]interface{})
		require.IsType(t, []interface{}{}, sliEntityData[SliConfigFieldSliEntityWebsiteEventBased])
		sliEntityWebsiteEventBasedSlice := sliEntityData[SliConfigFieldSliEntityWebsiteEventBased].([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntityWebsiteEventBasedSlice[0])
		sliEntityWebsiteEventBasedData := sliEntityWebsiteEventBasedSlice[0].(map[string]interface{})
		require.Equal(t, websiteId, sliEntityWebsiteEventBasedData[SliConfigFieldWebsiteID])
		require.Equal(t, beaconType, sliEntityWebsiteEventBasedData[SliConfigFieldBeaconType])
		require.Equal(t, "request.path@dest EQUALS 'good'", sliEntityWebsiteEventBasedData[SliConfigFieldGoodEventFilterExpression])
		require.Equal(t, "request.path@dest EQUALS 'bad'", sliEntityWebsiteEventBasedData[SliConfigFieldBadEventFilterExpression])
	}
}

func (r *sliConfigUnitTest) shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeWebsiteEventBasedWhenGoodEventFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		websiteId := "my-website"
		goodEventTagFilterExpression := invalidTagFilterExpression
		badEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "bad")
		beaconType := "pageLoad"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:                      "websiteEventBased",
				WebsiteId:                 &websiteId,
				BeaconType:                &beaconType,
				GoodEventFilterExpression: goodEventTagFilterExpression,
				BadEventFilterExpression:  badEventTagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported tag filter expression of type invalid")
	}
}

func (r *sliConfigUnitTest) shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeWebsiteEventBasedWhenBadEventFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		websiteId := "my-website"
		goodEventTagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "good")
		badEventTagFilterExpression := invalidTagFilterExpression
		beaconType := "pageLoad"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:                      "websiteEventBased",
				WebsiteId:                 &websiteId,
				BeaconType:                &beaconType,
				GoodEventFilterExpression: goodEventTagFilterExpression,
				BadEventFilterExpression:  badEventTagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported tag filter expression of type invalid")
	}
}

func (r *sliConfigUnitTest) shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeWebsiteTimeBased() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		websiteId := "my-website"
		tagFilterExpression := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, "request.path", restapi.EqualsOperator, "test")
		beaconType := "pageLoad"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:             "websiteTimeBased",
				WebsiteId:        &websiteId,
				BeaconType:       &beaconType,
				FilterExpression: tagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, sliConfigID, resourceData.Id())
		require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
		require.IsType(t, []interface{}{}, resourceData.Get(SliConfigFieldSliEntity))
		sliEntitySlice := resourceData.Get(SliConfigFieldSliEntity).([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntitySlice[0])
		sliEntityData := sliEntitySlice[0].(map[string]interface{})
		require.IsType(t, []interface{}{}, sliEntityData[SliConfigFieldSliEntityWebsiteTimeBased])
		sliEntityWebsiteTimeBasedSlice := sliEntityData[SliConfigFieldSliEntityWebsiteTimeBased].([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntityWebsiteTimeBasedSlice[0])
		sliEntityWebsiteTimeBasedData := sliEntityWebsiteTimeBasedSlice[0].(map[string]interface{})
		require.Equal(t, websiteId, sliEntityWebsiteTimeBasedData[SliConfigFieldWebsiteID])
		require.Equal(t, beaconType, sliEntityWebsiteTimeBasedData[SliConfigFieldBeaconType])
		require.Equal(t, "request.path@dest EQUALS 'test'", sliEntityWebsiteTimeBasedData[SliConfigFieldFilterExpression])
	}
}

func (r *sliConfigUnitTest) shouldFailToUpateResourceStateForSliConfigWithSliEntityOfTypeWebsiteTimeBasedWhenFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		websiteId := "my-website"
		tagFilterExpression := invalidTagFilterExpression
		beaconType := "pageLoad"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:             "websiteTimeBased",
				WebsiteId:        &websiteId,
				BeaconType:       &beaconType,
				FilterExpression: tagFilterExpression,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported tag filter expression of type invalid")
	}
}

func (r *sliConfigUnitTest) shouldMapStateOfSliConfigWithEntityOfTypeApplicationToDataModel() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityApplicationTimeBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldApplicationID: sliConfigEntityApplicationID,
						SliConfigFieldServiceID:     sliConfigEntityServiceID,
						SliConfigFieldEndpointID:    sliConfigEntityEndpointID,
						SliConfigFieldBoundaryScope: sliConfigEntityBoundaryScope,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		model, err := resourceHandle.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
		require.Equal(t, sliConfigID, model.GetIDForResourcePath())
		require.Equal(t, sliConfigName, model.Name, "name should be equal to name")
		require.Equal(t, sliConfigInitialEvaluationTimestamp, model.InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
		require.Equal(t, sliConfigMetricName, model.MetricConfiguration.Name)
		require.Equal(t, sliConfigMetricAggregation, model.MetricConfiguration.Aggregation)
		require.Equal(t, sliConfigMetricThreshold, model.MetricConfiguration.Threshold)
		require.Equal(t, "application", model.SliEntity.Type)
		require.Equal(t, sliConfigEntityApplicationID, *model.SliEntity.ApplicationID)
		require.Equal(t, sliConfigEntityServiceID, *model.SliEntity.ServiceID)
		require.Equal(t, sliConfigEntityEndpointID, *model.SliEntity.EndpointID)
		require.Equal(t, sliConfigEntityBoundaryScope, *model.SliEntity.BoundaryScope)
	}
}

func (r *sliConfigUnitTest) shouldMapStateOfSliConfigWithEntityOfTypeAvailabilityToDataModel() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityApplicationEventBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldApplicationID:             sliConfigEntityApplicationID,
						SliConfigFieldIncludeInternal:           true,
						SliConfigFieldIncludeSynthetic:          true,
						SliConfigFieldBoundaryScope:             sliConfigEntityBoundaryScope,
						SliConfigFieldGoodEventFilterExpression: sliConfigTagFilterExpressionString,
						SliConfigFieldBadEventFilterExpression:  sliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		model, err := resourceHandle.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
		require.Equal(t, sliConfigID, model.GetIDForResourcePath())
		require.Equal(t, sliConfigName, model.Name, "name should be equal to name")
		require.Equal(t, sliConfigInitialEvaluationTimestamp, model.InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
		require.Equal(t, sliConfigMetricName, model.MetricConfiguration.Name)
		require.Equal(t, sliConfigMetricAggregation, model.MetricConfiguration.Aggregation)
		require.Equal(t, sliConfigMetricThreshold, model.MetricConfiguration.Threshold)
		require.Equal(t, "availability", model.SliEntity.Type)
		require.Equal(t, sliConfigEntityApplicationID, *model.SliEntity.ApplicationID)
		require.Equal(t, sliConfigEntityBoundaryScope, *model.SliEntity.BoundaryScope)
		require.True(t, *model.SliEntity.IncludeInternal)
		require.True(t, *model.SliEntity.IncludeSynthetic)
		require.Equal(t, sliConfigTagFilterExpression, model.SliEntity.GoodEventFilterExpression)
		require.Equal(t, sliConfigTagFilterExpression, model.SliEntity.BadEventFilterExpression)
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWithEntityOfTypeAvailabilityToDataModelWhenGoodEventFilterExpressionIsInvalid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityApplicationEventBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldApplicationID:             sliConfigEntityApplicationID,
						SliConfigFieldIncludeInternal:           true,
						SliConfigFieldIncludeSynthetic:          true,
						SliConfigFieldBoundaryScope:             sliConfigEntityBoundaryScope,
						SliConfigFieldGoodEventFilterExpression: invalidSliConfigTagFilterExpressionString,
						SliConfigFieldBadEventFilterExpression:  sliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected token")
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWithEntityOfTypeAvailabilityToDataModelWhenBadEventFilterExpressionIsInvalid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityApplicationEventBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldApplicationID:             sliConfigEntityApplicationID,
						SliConfigFieldIncludeInternal:           true,
						SliConfigFieldIncludeSynthetic:          true,
						SliConfigFieldBoundaryScope:             sliConfigEntityBoundaryScope,
						SliConfigFieldGoodEventFilterExpression: sliConfigTagFilterExpressionString,
						SliConfigFieldBadEventFilterExpression:  invalidSliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected token")
	}
}

func (r *sliConfigUnitTest) shouldMapStateOfSliConfigWithEntityOfTypeWebsiteEventBasedToDataModel() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityWebsiteEventBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldWebsiteID:                 sliConfigEntityWebsiteID,
						SliConfigFieldBeaconType:                sliConfigEntityBeaconType,
						SliConfigFieldGoodEventFilterExpression: sliConfigTagFilterExpressionString,
						SliConfigFieldBadEventFilterExpression:  sliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		model, err := resourceHandle.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
		require.Equal(t, sliConfigID, model.GetIDForResourcePath())
		require.Equal(t, sliConfigName, model.Name, "name should be equal to name")
		require.Equal(t, sliConfigInitialEvaluationTimestamp, model.InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
		require.Equal(t, sliConfigMetricName, model.MetricConfiguration.Name)
		require.Equal(t, sliConfigMetricAggregation, model.MetricConfiguration.Aggregation)
		require.Equal(t, sliConfigMetricThreshold, model.MetricConfiguration.Threshold)
		require.Equal(t, "websiteEventBased", model.SliEntity.Type)
		require.Equal(t, sliConfigEntityWebsiteID, *model.SliEntity.WebsiteId)
		require.Equal(t, sliConfigEntityBeaconType, *model.SliEntity.BeaconType)
		require.Equal(t, sliConfigTagFilterExpression, model.SliEntity.GoodEventFilterExpression)
		require.Equal(t, sliConfigTagFilterExpression, model.SliEntity.BadEventFilterExpression)
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWithEntityOfTypeWebsiteEventBasedToDataModelWhenGoodEventFilterExpressionIsInvalid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityWebsiteEventBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldWebsiteID:                 sliConfigEntityWebsiteID,
						SliConfigFieldBeaconType:                sliConfigEntityBeaconType,
						SliConfigFieldGoodEventFilterExpression: invalidSliConfigTagFilterExpressionString,
						SliConfigFieldBadEventFilterExpression:  sliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected token")
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWithEntityOfTypeWebsiteEventBasedToDataModelWhenBadEventFilterExpressionIsInvalid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityWebsiteEventBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldWebsiteID:                 sliConfigEntityWebsiteID,
						SliConfigFieldBeaconType:                sliConfigEntityBeaconType,
						SliConfigFieldGoodEventFilterExpression: sliConfigTagFilterExpressionString,
						SliConfigFieldBadEventFilterExpression:  invalidSliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected token")
	}
}

func (r *sliConfigUnitTest) shouldMapStateOfSliConfigWithEntityOfTypeWebsiteTimeBasedToDataModel() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityWebsiteTimeBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldWebsiteID:        sliConfigEntityWebsiteID,
						SliConfigFieldBeaconType:       sliConfigEntityBeaconType,
						SliConfigFieldFilterExpression: sliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		model, err := resourceHandle.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
		require.Equal(t, sliConfigID, model.GetIDForResourcePath())
		require.Equal(t, sliConfigName, model.Name, "name should be equal to name")
		require.Equal(t, sliConfigInitialEvaluationTimestamp, model.InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
		require.Equal(t, sliConfigMetricName, model.MetricConfiguration.Name)
		require.Equal(t, sliConfigMetricAggregation, model.MetricConfiguration.Aggregation)
		require.Equal(t, sliConfigMetricThreshold, model.MetricConfiguration.Threshold)
		require.Equal(t, "websiteTimeBased", model.SliEntity.Type)
		require.Equal(t, sliConfigEntityWebsiteID, *model.SliEntity.WebsiteId)
		require.Equal(t, sliConfigEntityBeaconType, *model.SliEntity.BeaconType)
		require.Equal(t, sliConfigTagFilterExpression, model.SliEntity.FilterExpression)
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWithEntityOfTypeWebsiteTimeBasedToDataModelWhenFilterExpressionIsInvalid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityWebsiteTimeBased: []interface{}{
					map[string]interface{}{
						SliConfigFieldWebsiteID:        sliConfigEntityWebsiteID,
						SliConfigFieldBeaconType:       sliConfigEntityBeaconType,
						SliConfigFieldFilterExpression: invalidSliConfigTagFilterExpressionString,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "unexpected token")
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWhenNoSliEntityIsProvided() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "exactly one sli entity configuration is required")
	}
}

func (r *sliConfigUnitTest) shouldFailToMapStateOfSliConfigWhenNoSupportedliEntityIsProvided() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		sliEntityStateObject := []interface{}{map[string]interface{}{
			SliConfigFieldSliEntityApplicationTimeBased: []interface{}{},
		}}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		_, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.ErrorContains(t, err, "exactly one sli entity configuration of type")
	}
}

func (r *sliConfigUnitTest) shouldRequireMetricConfigurationThresholdToBeGreaterThanZero() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		resourceData.SetId(sliConfigID)
		setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
		setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

		metricConfigurationStateObject := []map[string]interface{}{
			{
				SliConfigFieldMetricName:        sliConfigMetricName,
				SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
				SliConfigFieldMetricThreshold:   0.0,
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

		_, metricThresholdIsOK := resourceData.GetOk("metric_configuration.0.threshold")
		require.False(t, metricThresholdIsOK)
	}
}
