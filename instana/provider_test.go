package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProviderShouldValidateInternally(t *testing.T) {
	err := Provider().InternalValidate()

	assert.Nil(t, err)
}

func TestProviderShouldContainValidSchemaDefinition(t *testing.T) {
	config := Provider()

	assert.NotNil(t, config.Schema)
	assert.Equal(t, 3, len(config.Schema))

	schemaAssert := testutils.NewTerraformSchemaAssert(config.Schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldAPIToken)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldEndpoint)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SchemaFieldTlsSkipVerify, false)
}

func TestProviderShouldContainValidResourceDefinitions(t *testing.T) {
	config := Provider()

	assert.Equal(t, 20, len(config.ResourcesMap))

	assert.NotNil(t, config.ResourcesMap[ResourceInstanaAPIToken])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaApplicationConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaApplicationAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaGlobalApplicationAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSliConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSloConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSloAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSloCorrectionConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaWebsiteMonitoringConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaWebsiteAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaInfraAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaGroup])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaCustomDashboard])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSyntheticTest])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSyntheticAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaCustomEventSpecification])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaAlertingChannel])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaAlertingConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaAutomationAction])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaAutomationPolicy])
}

func TestProviderShouldContainValidDataSourceDefinitions(t *testing.T) {
	config := Provider()

	assert.Equal(t, 6, len(config.DataSourcesMap))

	assert.NotNil(t, config.DataSourcesMap[DataSourceBuiltinEvent])
	assert.NotNil(t, config.DataSourcesMap[DataSourceCustomEventSpec])
	assert.NotNil(t, config.DataSourcesMap[DataSourceSyntheticLocation])
	assert.NotNil(t, config.DataSourcesMap[DataSourceAlertingChannel])
	assert.NotNil(t, config.DataSourcesMap[DataSourceAutomationAction])
	assert.NotNil(t, config.DataSourcesMap[DataSourceHostAgents])
}
