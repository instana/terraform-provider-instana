package synthetictest

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSyntheticTestResourceHandleFramework(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewSyntheticTestResourceHandleFramework()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSyntheticTestFramework, metadata.ResourceName)
		assert.Equal(t, int64(0), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewSyntheticTestResourceHandleFramework()
		metadata := handle.MetaData()

		schema := metadata.Schema
		assert.NotNil(t, schema.Attributes["id"])
		assert.NotNil(t, schema.Attributes["label"])
		assert.NotNil(t, schema.Attributes["description"])
		assert.NotNil(t, schema.Attributes["active"])
		assert.NotNil(t, schema.Attributes["application_id"])
		assert.NotNil(t, schema.Attributes["applications"])
		assert.NotNil(t, schema.Attributes["mobile_apps"])
		assert.NotNil(t, schema.Attributes["websites"])
		assert.NotNil(t, schema.Attributes["custom_properties"])
		assert.NotNil(t, schema.Attributes["locations"])
		assert.NotNil(t, schema.Attributes["rbac_tags"])
		assert.NotNil(t, schema.Attributes["playback_mode"])
		assert.NotNil(t, schema.Attributes["test_frequency"])
		assert.NotNil(t, schema.Attributes["http_action"])
		assert.NotNil(t, schema.Attributes["http_script"])
		assert.NotNil(t, schema.Attributes["browser_script"])
		assert.NotNil(t, schema.Attributes["dns"])
		assert.NotNil(t, schema.Attributes["ssl_certificate"])
		assert.NotNil(t, schema.Attributes["webpage_action"])
		assert.NotNil(t, schema.Attributes["webpage_script"])
	})
}

func TestMetaData(t *testing.T) {
	t.Run("should return metadata", func(t *testing.T) {
		resource := &syntheticTestResourceFramework{
			metaData: resourcehandle.ResourceMetaDataFramework{
				ResourceName:  ResourceInstanaSyntheticTestFramework,
				SchemaVersion: 0,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSyntheticTestFramework, metadata.ResourceName)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return synthetic test rest resource", func(t *testing.T) {
		resource := &syntheticTestResourceFramework{}

		mockAPI := &mockInstanaAPI{}
		restResource := resource.GetRestResource(mockAPI)

		assert.NotNil(t, restResource)
	})
}

// Mock API for testing
type mockInstanaAPI struct{}

func (m *mockInstanaAPI) CustomEventSpecifications() restapi.RestResource[*restapi.CustomEventSpecification] {
	return nil
}
func (m *mockInstanaAPI) BuiltinEventSpecifications() restapi.ReadOnlyRestResource[*restapi.BuiltinEventSpecification] {
	return nil
}
func (m *mockInstanaAPI) APITokens() restapi.RestResource[*restapi.APIToken] { return nil }
func (m *mockInstanaAPI) ApplicationConfigs() restapi.RestResource[*restapi.ApplicationConfig] {
	return nil
}
func (m *mockInstanaAPI) ApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) GlobalApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) AlertingChannels() restapi.RestResource[*restapi.AlertingChannel] {
	return nil
}
func (m *mockInstanaAPI) AlertingConfigurations() restapi.RestResource[*restapi.AlertingConfiguration] {
	return nil
}
func (m *mockInstanaAPI) SliConfigs() restapi.RestResource[*restapi.SliConfig]          { return nil }
func (m *mockInstanaAPI) SloConfigs() restapi.RestResource[*restapi.SloConfig]          { return nil }
func (m *mockInstanaAPI) SloAlertConfig() restapi.RestResource[*restapi.SloAlertConfig] { return nil }
func (m *mockInstanaAPI) SloCorrectionConfig() restapi.RestResource[*restapi.SloCorrectionConfig] {
	return nil
}
func (m *mockInstanaAPI) WebsiteMonitoringConfig() restapi.RestResource[*restapi.WebsiteMonitoringConfig] {
	return nil
}
func (m *mockInstanaAPI) WebsiteAlertConfig() restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) InfraAlertConfig() restapi.RestResource[*restapi.InfraAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) Groups() restapi.RestResource[*restapi.Group] { return nil }
func (m *mockInstanaAPI) CustomDashboards() restapi.RestResource[*restapi.CustomDashboard] {
	return nil
}
func (m *mockInstanaAPI) SyntheticTest() restapi.RestResource[*restapi.SyntheticTest] {
	return &mockSyntheticTestRestResource{}
}
func (m *mockInstanaAPI) SyntheticLocation() restapi.ReadOnlyRestResource[*restapi.SyntheticLocation] {
	return nil
}
func (m *mockInstanaAPI) SyntheticAlertConfigs() restapi.RestResource[*restapi.SyntheticAlertConfig] {
	return nil
}
func (m *mockInstanaAPI) AutomationActions() restapi.RestResource[*restapi.AutomationAction] {
	return nil
}
func (m *mockInstanaAPI) AutomationPolicies() restapi.RestResource[*restapi.AutomationPolicy] {
	return nil
}
func (m *mockInstanaAPI) HostAgents() restapi.ReadOnlyRestResource[*restapi.HostAgent]  { return nil }
func (m *mockInstanaAPI) LogAlertConfig() restapi.RestResource[*restapi.LogAlertConfig] { return nil }

// Mock rest resource
type mockSyntheticTestRestResource struct{}

func (m *mockSyntheticTestRestResource) GetAll() (*[]*restapi.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) GetOne(id string) (*restapi.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) Create(data *restapi.SyntheticTest) (*restapi.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) Update(data *restapi.SyntheticTest) (*restapi.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) Delete(data *restapi.SyntheticTest) error {
	return nil
}

func (m *mockSyntheticTestRestResource) DeleteByID(id string) error {
	return nil
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		resource := &syntheticTestResourceFramework{
			metaData: resourcehandle.ResourceMetaDataFramework{
				ResourceName:  ResourceInstanaSyntheticTestFramework,
				Schema:        NewSyntheticTestResourceHandleFramework().MetaData().Schema,
				SchemaVersion: 0,
			},
		}
		ctx := context.Background()

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}

		diags := resource.SetComputedFields(ctx, plan)
		assert.False(t, diags.HasError())
	})
}

func TestMapStateToDataObject(t *testing.T) {
	resource := &syntheticTestResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSyntheticTestFramework,
			Schema:        NewSyntheticTestResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should map HTTP Action from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpAction: &HttpActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Operation:         types.StringNull(),
				Headers:           types.MapNull(types.StringType),
				Body:              types.StringNull(),
				ValidationString:  types.StringNull(),
				FollowRedirect:    types.BoolNull(),
				AllowInsecure:     types.BoolNull(),
				ExpectStatus:      types.Int64Null(),
				ExpectMatch:       types.StringNull(),
				ExpectExists:      types.SetNull(types.StringType),
				ExpectNotEmpty:    types.SetNull(types.StringType),
				ExpectJson:        types.StringNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		for _, d := range resultDiags.Errors() {
			t.Logf("Error: %s - %s", d.Summary(), d.Detail())
		}
		for _, d := range resultDiags.Warnings() {
			t.Logf("Warning: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, resultDiags.HasError(), "Expected no errors but got diagnostics")
		assert.NotNil(t, result)
		assert.Equal(t, "HTTPAction", result.Configuration.SyntheticType)
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, "Test", result.Label)
	})

	t.Run("should map HTTP Script from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpScript: &HttpScriptConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				Script:            types.StringValue("console.log('test');"),
				ScriptType:        types.StringNull(),
				FileName:          types.StringNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "HTTPScript", result.Configuration.SyntheticType)
	})

	t.Run("should map Browser Script from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			BrowserScript: &BrowserScriptConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				Script:            types.StringValue("browser script"),
				ScriptType:        types.StringNull(),
				FileName:          types.StringNull(),
				Browser:           types.StringNull(),
				RecordVideo:       types.BoolNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "BrowserScript", result.Configuration.SyntheticType)
	})

	t.Run("should map DNS from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			DNS: &DNSConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				Lookup:            types.StringValue("example.com"),
				Server:            types.StringValue("8.8.8.8"),
				QueryType:         types.StringNull(),
				Port:              types.Int64Null(),
				Transport:         types.StringNull(),
				AcceptCNAME:       types.BoolNull(),
				LookupServerName:  types.BoolNull(),
				RecursiveLookups:  types.BoolNull(),
				ServerRetries:     types.Int64Null(),
				TargetValues: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				}}),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "DNS", result.Configuration.SyntheticType)
	})

	t.Run("should map SSL Certificate from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			SSLCertificate: &SSLCertificateConfigModel{
				MarkSyntheticCall:    types.BoolValue(false),
				Retries:              types.Int64Value(0),
				RetryInterval:        types.Int64Value(1),
				Timeout:              types.StringNull(),
				Hostname:             types.StringValue("example.com"),
				DaysRemainingCheck:   types.Int64Value(30),
				AcceptSelfSignedCert: types.BoolNull(),
				Port:                 types.Int64Null(),
				ValidationRules: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				}}),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "SSLCertificate", result.Configuration.SyntheticType)
	})

	t.Run("should map Webpage Action from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			WebpageAction: &WebpageActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Browser:           types.StringNull(),
				RecordVideo:       types.BoolNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "WebpageAction", result.Configuration.SyntheticType)
	})

	t.Run("should map Webpage Script from state", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			WebpageScript: &WebpageScriptConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				Script:            types.StringValue("webpage script"),
				FileName:          types.StringNull(),
				Browser:           types.StringNull(),
				RecordVideo:       types.BoolNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "WebpageScript", result.Configuration.SyntheticType)
	})

	t.Run("should error when no configuration provided", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.True(t, resultDiags.HasError())
		assert.Nil(t, result)
	})

	t.Run("should error when multiple configurations provided", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpAction: &HttpActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Operation:         types.StringNull(),
				Headers:           types.MapNull(types.StringType),
				Body:              types.StringNull(),
				ValidationString:  types.StringNull(),
				FollowRedirect:    types.BoolNull(),
				AllowInsecure:     types.BoolNull(),
				ExpectStatus:      types.Int64Null(),
				ExpectMatch:       types.StringNull(),
				ExpectExists:      types.SetNull(types.StringType),
				ExpectNotEmpty:    types.SetNull(types.StringType),
				ExpectJson:        types.StringNull(),
			},
			HttpScript: &HttpScriptConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				Script:            types.StringValue("script"),
				ScriptType:        types.StringNull(),
				FileName:          types.StringNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.True(t, resultDiags.HasError())
		assert.Nil(t, result)
	})

	t.Run("should map additional fields", func(t *testing.T) {
		rbacTags := []attr.Value{
			types.ObjectValueMust(
				map[string]attr.Type{
					"name":  types.StringType,
					"value": types.StringType,
				},
				map[string]attr.Value{
					"name":  types.StringValue("dept"),
					"value": types.StringValue("eng"),
				},
			),
		}

		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			PlaybackMode:     types.StringValue("Simultaneous"),
			ApplicationID:    types.StringValue("app-123"),
			TestFrequency:    types.Int64Value(30),
			CustomProperties: types.MapNull(types.StringType),
			RbacTags: types.SetValueMust(
				types.ObjectType{AttrTypes: map[string]attr.Type{
					"name":  types.StringType,
					"value": types.StringType,
				}},
				rbacTags,
			),
			HttpAction: &HttpActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Operation:         types.StringNull(),
				Headers:           types.MapNull(types.StringType),
				Body:              types.StringNull(),
				ValidationString:  types.StringNull(),
				FollowRedirect:    types.BoolNull(),
				AllowInsecure:     types.BoolNull(),
				ExpectStatus:      types.Int64Null(),
				ExpectMatch:       types.StringNull(),
				ExpectExists:      types.SetNull(types.StringType),
				ExpectNotEmpty:    types.SetNull(types.StringType),
				ExpectJson:        types.StringNull(),
			},
		}

		applications, _ := types.SetValueFrom(ctx, types.StringType, []string{"app-1"})
		model.Applications = applications

		mobileApps, _ := types.SetValueFrom(ctx, types.StringType, []string{"mobile-1"})
		model.MobileApps = mobileApps

		websites, _ := types.SetValueFrom(ctx, types.StringType, []string{"web-1"})
		model.Websites = websites

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		for _, d := range resultDiags.Errors() {
			t.Logf("Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, resultDiags.HasError(), "Expected no errors")
		require.NotNil(t, result)
		require.NotNil(t, result.ApplicationID, "ApplicationID should not be nil")
		assert.Equal(t, "app-123", *result.ApplicationID)
		assert.Equal(t, int32(30), *result.TestFrequency)
		assert.Len(t, result.Applications, 1)
		assert.Len(t, result.MobileApps, 1)
		assert.Len(t, result.Websites, 1)
		assert.Len(t, result.RbacTags, 1)
	})

	t.Run("should map HTTP Action with complex fields", func(t *testing.T) {
		expectExists, _ := types.SetValueFrom(ctx, types.StringType, []string{"$.data"})
		expectNotEmpty, _ := types.SetValueFrom(ctx, types.StringType, []string{"$.items"})

		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpAction: &HttpActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Operation:         types.StringValue("POST"),
				Headers:           types.MapNull(types.StringType),
				Body:              types.StringValue(`{"test": "data"}`),
				ValidationString:  types.StringValue("success"),
				FollowRedirect:    types.BoolValue(true),
				AllowInsecure:     types.BoolValue(false),
				ExpectStatus:      types.Int64Value(200),
				ExpectMatch:       types.StringValue(".*success.*"),
				ExpectExists:      expectExists,
				ExpectNotEmpty:    expectNotEmpty,
				ExpectJson:        types.StringNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "HTTPAction", result.Configuration.SyntheticType)
		// Headers is nil because we used MapNull
		assert.Nil(t, result.Configuration.Headers)
		assert.NotNil(t, result.Configuration.Body)
		assert.NotNil(t, result.Configuration.ExpectStatus)
		assert.Equal(t, int32(200), *result.Configuration.ExpectStatus)
	})

	t.Run("should map DNS with filters", func(t *testing.T) {
		targetValues := []attr.Value{
			types.ObjectValueMust(
				map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				},
				map[string]attr.Value{
					"key":      types.StringValue("A"),
					"operator": types.StringValue("EQUALS"),
					"value":    types.StringValue("192.168.1.1"),
				},
			),
		}

		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			DNS: &DNSConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				Lookup:            types.StringValue("example.com"),
				Server:            types.StringValue("8.8.8.8"),
				QueryType:         types.StringValue("A"),
				Port:              types.Int64Value(53),
				Transport:         types.StringValue("UDP"),
				AcceptCNAME:       types.BoolNull(),
				LookupServerName:  types.BoolNull(),
				RecursiveLookups:  types.BoolNull(),
				ServerRetries:     types.Int64Null(),
				QueryTime: &DNSFilterQueryTimeModel{
					Key:      types.StringValue("query_time"),
					Operator: types.StringValue("LESS_THAN"),
					Value:    types.Int64Value(100),
				},
				TargetValues: types.SetValueMust(
					types.ObjectType{AttrTypes: map[string]attr.Type{
						"key":      types.StringType,
						"operator": types.StringType,
						"value":    types.StringType,
					}},
					targetValues,
				),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "DNS", result.Configuration.SyntheticType)
		assert.NotNil(t, result.Configuration.QueryTime)
		assert.Len(t, result.Configuration.TargetValues, 1)
	})

	t.Run("should map HTTP Action with plan parameter", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpAction: &HttpActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Operation:         types.StringNull(),
				Headers:           types.MapNull(types.StringType),
				Body:              types.StringNull(),
				ValidationString:  types.StringNull(),
				FollowRedirect:    types.BoolNull(),
				AllowInsecure:     types.BoolNull(),
				ExpectStatus:      types.Int64Null(),
				ExpectMatch:       types.StringNull(),
				ExpectExists:      types.SetNull(types.StringType),
				ExpectNotEmpty:    types.SetNull(types.StringType),
				ExpectJson:        types.StringNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		diags := plan.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, plan, nil)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "HTTPAction", result.Configuration.SyntheticType)
	})

	t.Run("should map SSL Certificate with validation rules", func(t *testing.T) {
		validationRules := []attr.Value{
			types.ObjectValueMust(
				map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				},
				map[string]attr.Value{
					"key":      types.StringValue("issuer"),
					"operator": types.StringValue("CONTAINS"),
					"value":    types.StringValue("Let's Encrypt"),
				},
			),
		}

		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			SSLCertificate: &SSLCertificateConfigModel{
				MarkSyntheticCall:    types.BoolValue(false),
				Retries:              types.Int64Value(0),
				RetryInterval:        types.Int64Value(1),
				Timeout:              types.StringNull(),
				Hostname:             types.StringValue("example.com"),
				DaysRemainingCheck:   types.Int64Value(30),
				Port:                 types.Int64Value(443),
				AcceptSelfSignedCert: types.BoolValue(false),
				ValidationRules: types.SetValueMust(
					types.ObjectType{AttrTypes: map[string]attr.Type{
						"key":      types.StringType,
						"operator": types.StringType,
						"value":    types.StringType,
					}},
					validationRules,
				),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		for _, d := range diags.Errors() {
			t.Logf("State.Set Error: %s - %s", d.Summary(), d.Detail())
		}
		require.False(t, diags.HasError(), "Expected no errors setting state")

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "SSLCertificate", result.Configuration.SyntheticType)
		assert.Len(t, result.Configuration.ValidationRules, 1)
	})
}

func TestUpdateState(t *testing.T) {
	resource := &syntheticTestResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSyntheticTestFramework,
			Schema:        NewSyntheticTestResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should update state with HTTP Action", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.Equal(t, "test-id", model.ID.ValueString())
		assert.NotNil(t, model.HttpAction)
	})

	t.Run("should update state with HTTP Script", func(t *testing.T) {
		script := "console.log('test');"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPScript",
				Script:            &script,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.HttpScript)
	})

	t.Run("should update state with Browser Script", func(t *testing.T) {
		script := "browser script"
		browser := "chrome"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "BrowserScript",
				Script:            &script,
				Browser:           &browser,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.BrowserScript)
	})

	t.Run("should update state with DNS", func(t *testing.T) {
		lookup := "example.com"
		server := "8.8.8.8"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "DNS",
				Lookup:            &lookup,
				Server:            &server,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.DNS)
	})

	t.Run("should update state with SSL Certificate", func(t *testing.T) {
		hostname := "example.com"
		days := int32(30)
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall:  false,
				Retries:            0,
				RetryInterval:      1,
				SyntheticType:      "SSLCertificate",
				Hostname:           &hostname,
				DaysRemainingCheck: &days,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.SSLCertificate)
	})

	t.Run("should update state with Webpage Action", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "WebpageAction",
				URL:               &url,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.WebpageAction)
	})

	t.Run("should update state with Webpage Script", func(t *testing.T) {
		script := "webpage script"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "WebpageScript",
				Script:            &script,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.WebpageScript)
	})

	t.Run("should update state with HTTP Script with Scripts", func(t *testing.T) {
		script := "console.log('test');"
		scriptType := "Jest"
		fileName := "test.js"
		bundle := "bundle content"
		scriptFile := "script file content"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPScript",
				Script:            &script,
				ScriptType:        &scriptType,
				FileName:          &fileName,
				Scripts: &restapi.MultipleScriptsConfiguration{
					Bundle:     &bundle,
					ScriptFile: &scriptFile,
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.HttpScript)
		assert.NotNil(t, model.HttpScript.Scripts)
		assert.Equal(t, "bundle content", model.HttpScript.Scripts.Bundle.ValueString())
	})

	t.Run("should update state with Browser Script with Scripts", func(t *testing.T) {
		script := "browser script"
		browser := "chrome"
		recordVideo := true
		bundle := "browser bundle"
		scriptFile := "browser script file"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "BrowserScript",
				Script:            &script,
				Browser:           &browser,
				RecordVideo:       &recordVideo,
				Scripts: &restapi.MultipleScriptsConfiguration{
					Bundle:     &bundle,
					ScriptFile: &scriptFile,
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.BrowserScript)
		assert.NotNil(t, model.BrowserScript.Scripts)
		assert.True(t, model.BrowserScript.RecordVideo.ValueBool())
	})

	t.Run("should update state with DNS with QueryTime", func(t *testing.T) {
		lookup := "example.com"
		server := "8.8.8.8"
		port := int32(53)
		serverRetries := int32(3)
		acceptCNAME := true
		lookupServerName := false
		recursiveLookups := true
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "DNS",
				Lookup:            &lookup,
				Server:            &server,
				Port:              &port,
				ServerRetries:     &serverRetries,
				AcceptCNAME:       &acceptCNAME,
				LookupServerName:  &lookupServerName,
				RecursiveLookups:  &recursiveLookups,
				QueryTime: &restapi.DNSFilterQueryTime{
					Key:      "query_time",
					Operator: "LESS_THAN",
					Value:    100,
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.DNS)
		assert.NotNil(t, model.DNS.QueryTime)
		assert.Equal(t, "query_time", model.DNS.QueryTime.Key.ValueString())
	})

	t.Run("should update state with SSL Certificate with Port", func(t *testing.T) {
		hostname := "example.com"
		days := int32(30)
		port := int32(443)
		acceptSelfSigned := true
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall:    false,
				Retries:              0,
				RetryInterval:        1,
				SyntheticType:        "SSLCertificate",
				Hostname:             &hostname,
				DaysRemainingCheck:   &days,
				SSLPort:              &port,
				AcceptSelfSignedCert: &acceptSelfSigned,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.SSLCertificate)
		assert.Equal(t, int64(443), model.SSLCertificate.Port.ValueInt64())
		assert.True(t, model.SSLCertificate.AcceptSelfSignedCert.ValueBool())
	})

	t.Run("should update state with HTTP Script with empty FileName", func(t *testing.T) {
		script := "console.log('test');"
		fileName := ""
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPScript",
				Script:            &script,
				FileName:          &fileName,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.HttpScript)
		assert.True(t, model.HttpScript.FileName.IsNull())
	})

	t.Run("should update state with Browser Script with empty FileName", func(t *testing.T) {
		script := "browser script"
		fileName := ""
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "BrowserScript",
				Script:            &script,
				FileName:          &fileName,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.BrowserScript)
		assert.True(t, model.BrowserScript.FileName.IsNull())
	})

	t.Run("should update state with Webpage Script with empty FileName", func(t *testing.T) {
		script := "webpage script"
		fileName := ""
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "WebpageScript",
				Script:            &script,
				FileName:          &fileName,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.WebpageScript)
		assert.True(t, model.WebpageScript.FileName.IsNull())
	})

	t.Run("should update state with HTTP Action with all optional fields", func(t *testing.T) {
		url := "https://example.com"
		operation := "POST"
		body := "test body"
		validationString := "success"
		followRedirect := true
		allowInsecure := false
		expectStatus := int32(200)
		expectMatch := ".*success.*"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
				Operation:         &operation,
				Body:              &body,
				ValidationString:  &validationString,
				FollowRedirect:    &followRedirect,
				AllowInsecure:     &allowInsecure,
				ExpectStatus:      &expectStatus,
				ExpectMatch:       &expectMatch,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.HttpAction)
		assert.True(t, model.HttpAction.FollowRedirect.ValueBool())
		assert.False(t, model.HttpAction.AllowInsecure.ValueBool())
	})

	t.Run("should update state with empty arrays", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Applications: []string{},
			MobileApps:   []string{},
			Websites:     []string{},
			RbacTags:     []restapi.ApiTag{},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
				ExpectExists:      []string{},
				ExpectNotEmpty:    []string{},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.True(t, model.Applications.IsNull())
		assert.True(t, model.MobileApps.IsNull())
		assert.True(t, model.Websites.IsNull())
	})

	t.Run("should update state with nil TestFrequency", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:            "test-id",
			Label:         "Test",
			Active:        true,
			PlaybackMode:  "Simultaneous",
			Locations:     []string{"loc-1"},
			TestFrequency: nil,
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.True(t, model.TestFrequency.IsNull())
	})

	t.Run("should update state with empty CustomProperties", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:               "test-id",
			Label:            "Test",
			Active:           true,
			PlaybackMode:     "Simultaneous",
			Locations:        []string{"loc-1"},
			CustomProperties: map[string]string{},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.True(t, model.CustomProperties.IsNull())
	})

	t.Run("should update state with HTTP Action empty Headers and ExpectJson", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
				Headers:           map[string]string{},
				ExpectJson:        nil,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.HttpAction)
		assert.True(t, model.HttpAction.Headers.IsNull())
		assert.True(t, model.HttpAction.ExpectJson.IsNull())
	})

	t.Run("should update state with DNS empty TargetValues", func(t *testing.T) {
		lookup := "example.com"
		server := "8.8.8.8"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "DNS",
				Lookup:            &lookup,
				Server:            &server,
				TargetValues:      []restapi.DNSFilterTargetValue{},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.DNS)
		assert.True(t, model.DNS.TargetValues.IsNull())
	})

	t.Run("should update state with SSL Certificate empty ValidationRules", func(t *testing.T) {
		hostname := "example.com"
		days := int32(30)
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall:  false,
				Retries:            0,
				RetryInterval:      1,
				SyntheticType:      "SSLCertificate",
				Hostname:           &hostname,
				DaysRemainingCheck: &days,
				ValidationRules:    []restapi.SSLCertificateValidation{},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.SSLCertificate)
		assert.True(t, model.SSLCertificate.ValidationRules.IsNull())
	})

	t.Run("should update state with additional fields", func(t *testing.T) {
		appID := "app-123"
		testFreq := int32(30)
		url := "https://example.com"
		desc := "Test Description"

		apiObject := &restapi.SyntheticTest{
			ID:            "test-id",
			Label:         "Test",
			Description:   &desc,
			Active:        true,
			PlaybackMode:  "Simultaneous",
			ApplicationID: &appID,
			Applications:  []string{"app-1"},
			MobileApps:    []string{"mobile-1"},
			Websites:      []string{"web-1"},
			Locations:     []string{"loc-1"},
			TestFrequency: &testFreq,
			CustomProperties: map[string]string{
				"env": "prod",
			},
			RbacTags: []restapi.ApiTag{
				{Name: "dept", Value: "eng"},
			},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.Equal(t, "Test Description", model.Description.ValueString())
		assert.Equal(t, "app-123", model.ApplicationID.ValueString())
		assert.Equal(t, int64(30), model.TestFrequency.ValueInt64())
	})
}

func TestMapStateToDataObjectEdgeCases(t *testing.T) {
	resource := &syntheticTestResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSyntheticTestFramework,
			Schema:        NewSyntheticTestResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should handle error when getting from plan fails", func(t *testing.T) {
		// Create an invalid plan that will cause Get to fail
		plan := &tfsdk.Plan{
			Schema: resource.metaData.Schema,
		}
		// Don't set any data in the plan, which should cause an error when trying to get locations

		result, resultDiags := resource.MapStateToDataObject(ctx, plan, nil)
		// The function should handle the error gracefully
		assert.True(t, resultDiags.HasError() || result != nil)
	})

	t.Run("should map with empty custom properties", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapValueMust(types.StringType, map[string]attr.Value{}),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpAction: &HttpActionConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				URL:               types.StringValue("https://example.com"),
				Operation:         types.StringNull(),
				Headers:           types.MapNull(types.StringType),
				Body:              types.StringNull(),
				ValidationString:  types.StringNull(),
				FollowRedirect:    types.BoolNull(),
				AllowInsecure:     types.BoolNull(),
				ExpectStatus:      types.Int64Null(),
				ExpectMatch:       types.StringNull(),
				ExpectExists:      types.SetNull(types.StringType),
				ExpectNotEmpty:    types.SetNull(types.StringType),
				ExpectJson:        types.StringNull(),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.NotNil(t, result.CustomProperties)
		assert.Equal(t, 0, len(result.CustomProperties))
	})
}

func TestUpdateStateEdgeCases(t *testing.T) {
	resource := &syntheticTestResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSyntheticTestFramework,
			Schema:        NewSyntheticTestResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should update state with Webpage Action with RecordVideo", func(t *testing.T) {
		url := "https://example.com"
		browser := "chrome"
		recordVideo := true
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "WebpageAction",
				URL:               &url,
				Browser:           &browser,
				RecordVideo:       &recordVideo,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.WebpageAction)
		assert.True(t, model.WebpageAction.RecordVideo.ValueBool())
	})

	t.Run("should update state with HTTP Action with populated Headers and ExpectJson", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPAction",
				URL:               &url,
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer token",
				},
				ExpectExists:   []string{"$.data", "$.status"},
				ExpectNotEmpty: []string{"$.items"},
				ExpectJson:     []byte(`{"status":"success","code":200}`),
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.HttpAction)
		assert.False(t, model.HttpAction.Headers.IsNull())
		assert.False(t, model.HttpAction.ExpectJson.IsNull())
		assert.False(t, model.HttpAction.ExpectExists.IsNull())
		assert.False(t, model.HttpAction.ExpectNotEmpty.IsNull())
	})

	t.Run("should update state with DNS with populated TargetValues", func(t *testing.T) {
		lookup := "example.com"
		server := "8.8.8.8"
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "DNS",
				Lookup:            &lookup,
				Server:            &server,
				TargetValues: []restapi.DNSFilterTargetValue{
					{
						Key:      "A",
						Operator: "EQUALS",
						Value:    "192.168.1.1",
					},
					{
						Key:      "AAAA",
						Operator: "CONTAINS",
						Value:    "2001:db8",
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.DNS)
		assert.False(t, model.DNS.TargetValues.IsNull())
	})

	t.Run("should update state with SSL Certificate with populated ValidationRules", func(t *testing.T) {
		hostname := "example.com"
		days := int32(30)
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall:  false,
				Retries:            0,
				RetryInterval:      1,
				SyntheticType:      "SSLCertificate",
				Hostname:           &hostname,
				DaysRemainingCheck: &days,
				ValidationRules: []restapi.SSLCertificateValidation{
					{
						Key:      "issuer",
						Operator: "CONTAINS",
						Value:    "Let's Encrypt",
					},
					{
						Key:      "subject",
						Operator: "EQUALS",
						Value:    "example.com",
					},
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.SSLCertificate)
		assert.False(t, model.SSLCertificate.ValidationRules.IsNull())
	})

	t.Run("should update state with Webpage Script with RecordVideo", func(t *testing.T) {
		script := "webpage script"
		browser := "firefox"
		recordVideo := false
		apiObject := &restapi.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: restapi.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "WebpageScript",
				Script:            &script,
				Browser:           &browser,
				RecordVideo:       &recordVideo,
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.WebpageScript)
		assert.False(t, model.WebpageScript.RecordVideo.ValueBool())
	})
}

func TestMapConfigurationOptionalFields(t *testing.T) {
	resource := &syntheticTestResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
			ResourceName:  ResourceInstanaSyntheticTestFramework,
			Schema:        NewSyntheticTestResourceHandleFramework().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should map HTTP Script with Scripts configuration", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			HttpScript: &HttpScriptConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringValue("30s"),
				Script:            types.StringValue("console.log('test');"),
				ScriptType:        types.StringValue("Jest"),
				FileName:          types.StringValue("test.js"),
				Scripts: &MultipleScriptsModel{
					Bundle:     types.StringValue("bundle content"),
					ScriptFile: types.StringValue("script file content"),
				},
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "HTTPScript", result.Configuration.SyntheticType)
		assert.NotNil(t, result.Configuration.Scripts)
		assert.NotNil(t, result.Configuration.Scripts.Bundle)
		assert.Equal(t, "bundle content", *result.Configuration.Scripts.Bundle)
	})

	t.Run("should map Browser Script with Scripts configuration", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			BrowserScript: &BrowserScriptConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringValue("60s"),
				Script:            types.StringValue("browser script"),
				ScriptType:        types.StringValue("Jest"),
				FileName:          types.StringValue("browser.js"),
				Browser:           types.StringValue("chrome"),
				RecordVideo:       types.BoolValue(true),
				Scripts: &MultipleScriptsModel{
					Bundle:     types.StringValue("browser bundle"),
					ScriptFile: types.StringValue("browser script file"),
				},
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "BrowserScript", result.Configuration.SyntheticType)
		assert.NotNil(t, result.Configuration.Scripts)
		assert.NotNil(t, result.Configuration.RecordVideo)
		assert.True(t, *result.Configuration.RecordVideo)
	})

	t.Run("should map DNS with all optional fields", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			DNS: &DNSConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(1),
				RetryInterval:     types.Int64Value(2),
				Timeout:           types.StringValue("10s"),
				Lookup:            types.StringValue("example.com"),
				Server:            types.StringValue("8.8.8.8"),
				QueryType:         types.StringValue("A"),
				Port:              types.Int64Value(53),
				Transport:         types.StringValue("UDP"),
				AcceptCNAME:       types.BoolValue(true),
				LookupServerName:  types.BoolValue(false),
				RecursiveLookups:  types.BoolValue(true),
				ServerRetries:     types.Int64Value(3),
				QueryTime:         nil,
				TargetValues: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				}}),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "DNS", result.Configuration.SyntheticType)
		assert.NotNil(t, result.Configuration.Port)
		assert.Equal(t, int32(53), *result.Configuration.Port)
		assert.NotNil(t, result.Configuration.ServerRetries)
		assert.Equal(t, int32(3), *result.Configuration.ServerRetries)
	})

	t.Run("should map SSL Certificate with all optional fields", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("test-id"),
			Label:            types.StringValue("Test"),
			Description:      types.StringNull(),
			Active:           types.BoolValue(true),
			ApplicationID:    types.StringNull(),
			Applications:     types.SetNull(types.StringType),
			MobileApps:       types.SetNull(types.StringType),
			Websites:         types.SetNull(types.StringType),
			CustomProperties: types.MapNull(types.StringType),
			PlaybackMode:     types.StringValue("Simultaneous"),
			TestFrequency:    types.Int64Null(),
			RbacTags: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}}),
			SSLCertificate: &SSLCertificateConfigModel{
				MarkSyntheticCall:    types.BoolValue(false),
				Retries:              types.Int64Value(1),
				RetryInterval:        types.Int64Value(2),
				Timeout:              types.StringValue("15s"),
				Hostname:             types.StringValue("example.com"),
				DaysRemainingCheck:   types.Int64Value(30),
				Port:                 types.Int64Value(443),
				AcceptSelfSignedCert: types.BoolValue(true),
				ValidationRules: types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				}}),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}
		diags := state.Set(ctx, model)
		require.False(t, diags.HasError())

		result, resultDiags := resource.MapStateToDataObject(ctx, nil, state)
		assert.False(t, resultDiags.HasError())
		assert.NotNil(t, result)
		assert.Equal(t, "SSLCertificate", result.Configuration.SyntheticType)
		assert.NotNil(t, result.Configuration.SSLPort)
		assert.Equal(t, int32(443), *result.Configuration.SSLPort)
		assert.NotNil(t, result.Configuration.AcceptSelfSignedCert)
		assert.True(t, *result.Configuration.AcceptSelfSignedCert)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("getStringPointerFromFrameworkType should handle null", func(t *testing.T) {
		result := getStringPointerFromFrameworkType(types.StringNull())
		assert.Nil(t, result)
	})

	t.Run("getStringPointerFromFrameworkType should handle unknown", func(t *testing.T) {
		result := getStringPointerFromFrameworkType(types.StringUnknown())
		assert.Nil(t, result)
	})

	t.Run("getStringPointerFromFrameworkType should return value", func(t *testing.T) {
		result := getStringPointerFromFrameworkType(types.StringValue("test"))
		assert.NotNil(t, result)
		assert.Equal(t, "test", *result)
	})

	t.Run("getBoolPointerFromFrameworkType should handle null", func(t *testing.T) {
		result := getBoolPointerFromFrameworkType(types.BoolNull())
		assert.Nil(t, result)
	})

	t.Run("getBoolPointerFromFrameworkType should handle unknown", func(t *testing.T) {
		result := getBoolPointerFromFrameworkType(types.BoolUnknown())
		assert.Nil(t, result)
	})

	t.Run("getBoolPointerFromFrameworkType should return value", func(t *testing.T) {
		result := getBoolPointerFromFrameworkType(types.BoolValue(true))
		assert.NotNil(t, result)
		assert.True(t, *result)
	})
}

func TestInt64Validator(t *testing.T) {
	ctx := context.Background()

	t.Run("should validate value within range", func(t *testing.T) {
		v := int64Validator{min: 1, max: 10}

		assert.Equal(t, "Value must be between 1 and 10", v.Description(ctx))
		assert.Equal(t, "Value must be between 1 and 10", v.MarkdownDescription(ctx))
	})

	t.Run("should accept valid value", func(t *testing.T) {
		v := int64Validator{min: 1, max: 10}
		req := validator.Int64Request{
			Path:        path.Root("test"),
			ConfigValue: types.Int64Value(5),
		}
		resp := &validator.Int64Response{}

		v.ValidateInt64(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("should reject value below minimum", func(t *testing.T) {
		v := int64Validator{min: 1, max: 10}
		req := validator.Int64Request{
			Path:        path.Root("test"),
			ConfigValue: types.Int64Value(0),
		}
		resp := &validator.Int64Response{}

		v.ValidateInt64(ctx, req, resp)

		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("should reject value above maximum", func(t *testing.T) {
		v := int64Validator{min: 1, max: 10}
		req := validator.Int64Request{
			Path:        path.Root("test"),
			ConfigValue: types.Int64Value(11),
		}
		resp := &validator.Int64Response{}

		v.ValidateInt64(ctx, req, resp)

		assert.True(t, resp.Diagnostics.HasError())
	})

	t.Run("should accept null value", func(t *testing.T) {
		v := int64Validator{min: 1, max: 10}
		req := validator.Int64Request{
			Path:        path.Root("test"),
			ConfigValue: types.Int64Null(),
		}
		resp := &validator.Int64Response{}

		v.ValidateInt64(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("should accept unknown value", func(t *testing.T) {
		v := int64Validator{min: 1, max: 10}
		req := validator.Int64Request{
			Path:        path.Root("test"),
			ConfigValue: types.Int64Unknown(),
		}
		resp := &validator.Int64Response{}

		v.ValidateInt64(ctx, req, resp)

		assert.False(t, resp.Diagnostics.HasError())
	})
}

// Made with Bob
