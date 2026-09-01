package synthetictest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/rest"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSyntheticTestResourceHandle(t *testing.T) {
	t.Run("should create resource handle with correct metadata", func(t *testing.T) {
		handle := NewSyntheticTestResourceHandle()

		require.NotNil(t, handle)
		metadata := handle.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSyntheticTest, metadata.ResourceName)
		assert.Equal(t, int64(1), metadata.SchemaVersion)
	})

	t.Run("should have correct schema attributes", func(t *testing.T) {
		handle := NewSyntheticTestResourceHandle()
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
		resource := &syntheticTestResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaSyntheticTest,
				SchemaVersion: 0,
			},
		}
		metadata := resource.MetaData()
		require.NotNil(t, metadata)
		assert.Equal(t, ResourceInstanaSyntheticTest, metadata.ResourceName)
	})
}

func TestGetRestResource(t *testing.T) {
	t.Run("should return synthetic test rest resource", func(t *testing.T) {
		resource := &syntheticTestResource{}

		mockAPI := &mockSyntheticTestAPI{}
		restResource := resource.GetRestResource(mockAPI)

		assert.NotNil(t, restResource)
	})
}

// mockSyntheticTestAPI extends the common mock to provide specific behavior for synthetic test tests
type mockSyntheticTestAPI struct {
	testutils.MockInstanaAPI
}

func (m *mockSyntheticTestAPI) SyntheticTests() rest.RestResource[*api.SyntheticTest] {
	return &mockSyntheticTestRestResource{}
}

// Mock rest resource
type mockSyntheticTestRestResource struct{}

func (m *mockSyntheticTestRestResource) GetAll() (*[]*api.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) GetOne(id string) (*api.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) Create(data *api.SyntheticTest) (*api.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) Update(data *api.SyntheticTest) (*api.SyntheticTest, error) {
	return nil, nil
}

func (m *mockSyntheticTestRestResource) Delete(data *api.SyntheticTest) error {
	return nil
}

func (m *mockSyntheticTestRestResource) DeleteByID(id string) error {
	return nil
}

func TestSetComputedFields(t *testing.T) {
	t.Run("should return nil diagnostics", func(t *testing.T) {
		resource := &syntheticTestResource{
			metaData: resourcehandle.ResourceMetaData{
				ResourceName:  ResourceInstanaSyntheticTest,
				Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
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
	resource := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
					"id":           types.StringType,
					"display_name": types.StringType,
				},
				map[string]attr.Value{
					"id":           types.StringValue("dept"),
					"display_name": types.StringValue("eng"),
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
					"id":           types.StringType,
					"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
	resource := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should update state with HTTP Action", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.SSLCertificate)
	})

	t.Run("should update state with Webpage Action", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.NotNil(t, model.WebpageAction)
	})

	t.Run("should update state with Webpage Script", func(t *testing.T) {
		script := "webpage script"
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "HTTPScript",
				Script:            &script,
				ScriptType:        &scriptType,
				FileName:          &fileName,
				Scripts: &api.MultipleScriptsConfiguration{
					Bundle:     &bundle,
					ScriptFile: &scriptFile,
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "BrowserScript",
				Script:            &script,
				Browser:           &browser,
				RecordVideo:       &recordVideo,
				Scripts: &api.MultipleScriptsConfiguration{
					Bundle:     &bundle,
					ScriptFile: &scriptFile,
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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
				QueryTime: &api.DNSFilterQueryTime{
					Key:      "query_time",
					Operator: "LESS_THAN",
					Value:    100,
				},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Applications: []string{},
			MobileApps:   []string{},
			Websites:     []string{},
			RbacTags:     []api.RbacTag{},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:            "test-id",
			Label:         "Test",
			Active:        true,
			PlaybackMode:  "Simultaneous",
			Locations:     []string{"loc-1"},
			TestFrequency: nil,
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.True(t, model.TestFrequency.IsNull())
	})

	t.Run("should update state with empty CustomProperties", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &api.SyntheticTest{
			ID:               "test-id",
			Label:            "Test",
			Active:           true,
			PlaybackMode:     "Simultaneous",
			Locations:        []string{"loc-1"},
			CustomProperties: map[string]string{},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

		diags := resource.UpdateState(ctx, state, nil, apiObject)
		assert.False(t, diags.HasError())

		var model SyntheticTestModel
		diags = state.Get(ctx, &model)
		assert.False(t, diags.HasError())
		assert.False(t, model.CustomProperties.IsNull())
		assert.Equal(t, 0, len(model.CustomProperties.Elements()))
	})

	t.Run("should update state with HTTP Action empty Headers and ExpectJson", func(t *testing.T) {
		url := "https://example.com"
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "DNS",
				Lookup:            &lookup,
				Server:            &server,
				TargetValues:      []api.DNSFilterTargetValue{},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall:  false,
				Retries:            0,
				RetryInterval:      1,
				SyntheticType:      "SSLCertificate",
				Hostname:           &hostname,
				DaysRemainingCheck: &days,
				ValidationRules:    []api.SSLCertificateValidation{},
			},
		}

		state := &tfsdk.State{
			Schema: resource.metaData.Schema,
		}

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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

		apiObject := &api.SyntheticTest{
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
			RbacTags: []api.RbacTag{
				{ID: "dept", DisplayName: "eng"},
			},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
	resource := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
	resource := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should update state with Webpage Action with RecordVideo", func(t *testing.T) {
		url := "https://example.com"
		browser := "chrome"
		recordVideo := true
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall: false,
				Retries:           0,
				RetryInterval:     1,
				SyntheticType:     "DNS",
				Lookup:            &lookup,
				Server:            &server,
				TargetValues: []api.DNSFilterTargetValue{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall:  false,
				Retries:            0,
				RetryInterval:      1,
				SyntheticType:      "SSLCertificate",
				Hostname:           &hostname,
				DaysRemainingCheck: &days,
				ValidationRules: []api.SSLCertificateValidation{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
		apiObject := &api.SyntheticTest{
			ID:           "test-id",
			Label:        "Test",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-1"},
			Configuration: api.SyntheticTestConfig{
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

		// Initialize state with empty model
		initializeEmptyState(t, ctx, state)

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
	resource := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
				"id":           types.StringType,
				"display_name": types.StringType,
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
	t.Run("getStringPointerFromType should handle null", func(t *testing.T) {
		result := getStringPointerFromType(types.StringNull())
		assert.Nil(t, result)
	})

	t.Run("getStringPointerFromType should handle unknown", func(t *testing.T) {
		result := getStringPointerFromType(types.StringUnknown())
		assert.Nil(t, result)
	})

	t.Run("getStringPointerFromType should return value", func(t *testing.T) {
		result := getStringPointerFromType(types.StringValue("test"))
		assert.NotNil(t, result)
		assert.Equal(t, "test", *result)
	})

	t.Run("getBoolPointerFromType should handle null", func(t *testing.T) {
		result := getBoolPointerFromType(types.BoolNull())
		assert.Nil(t, result)
	})

	t.Run("getBoolPointerFromType should handle unknown", func(t *testing.T) {
		result := getBoolPointerFromType(types.BoolUnknown())
		assert.Nil(t, result)
	})

	t.Run("getBoolPointerFromType should return value", func(t *testing.T) {
		result := getBoolPointerFromType(types.BoolValue(true))
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

// initializeEmptyState initializes the state with an empty model to ensure proper state initialization
func initializeEmptyState(t *testing.T, ctx context.Context, state *tfsdk.State) {
	emptyModel := SyntheticTestModel{
		ID:               types.StringNull(),
		Label:            types.StringNull(),
		Description:      types.StringNull(),
		Active:           types.BoolNull(),
		ApplicationID:    types.StringNull(),
		Applications:     types.SetNull(types.StringType),
		MobileApps:       types.SetNull(types.StringType),
		Websites:         types.SetNull(types.StringType),
		CustomProperties: types.MapNull(types.StringType),
		Locations:        types.SetNull(types.StringType),
		PlaybackMode:     types.StringNull(),
		TestFrequency:    types.Int64Null(),
		RbacTags:         types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType, "display_name": types.StringType}}),
		HttpAction:       nil,
		HttpScript:       nil,
		BrowserScript:    nil,
		DNS:              nil,
		SSLCertificate:   nil,
		WebpageAction:    nil,
		WebpageScript:    nil,
		ICMP:             nil,
	}
	diags := state.Set(ctx, emptyModel)
	require.False(t, diags.HasError(), "Failed to initialize empty state")
}

// ---------------------------------------------------------------------------
// ICMP tests
// ---------------------------------------------------------------------------

func icmpValidationRulesAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"key":      types.StringType,
		"operator": types.StringType,
		"value":    types.Int64Type,
	}
}

func TestMapStateToDataObjectICMP(t *testing.T) {
	resource := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should map ICMP from state with all fields", func(t *testing.T) {
		ruleObj, _ := types.ObjectValue(
			icmpValidationRulesAttrType(),
			map[string]attr.Value{
				"key":      types.StringValue("packetLoss"),
				"operator": types.StringValue("LESS_THAN"),
				"value":    types.Int64Value(5),
			},
		)
		rules, _ := types.SetValue(
			types.ObjectType{AttrTypes: icmpValidationRulesAttrType()},
			[]attr.Value{ruleObj},
		)

		model := SyntheticTestModel{
			ID:               types.StringValue("icmp-id"),
			Label:            types.StringValue("ICMP Test"),
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
				"id":           types.StringType,
				"display_name": types.StringType,
			}}),
			ICMP: &ICMPConfigModel{
				MarkSyntheticCall:   types.BoolValue(true),
				Retries:             types.Int64Value(1),
				RetryInterval:       types.Int64Value(2),
				Timeout:             types.StringValue("30s"),
				TargetHost:          types.StringValue("192.0.2.1"),
				PacketCount:         types.Int64Value(5),
				PacketInterval:      types.StringValue("1s"),
				PacketSize:          types.Int64Value(64),
				PacketTimeout:       types.StringValue("2s"),
				UseDNS:              types.BoolValue(true),
				UseIPv6:             types.BoolValue(false),
				ICMPValidationRules: rules,
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{Schema: resource.metaData.Schema}
		require.False(t, state.Set(ctx, model).HasError())

		result, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, result)

		cfg := result.Configuration
		assert.Equal(t, SyntheticTestTypeICMP, cfg.SyntheticType)
		assert.True(t, cfg.MarkSyntheticCall)
		assert.Equal(t, int32(1), cfg.Retries)
		assert.Equal(t, int32(2), cfg.RetryInterval)
		require.NotNil(t, cfg.TargetHost)
		assert.Equal(t, "192.0.2.1", *cfg.TargetHost)
		require.NotNil(t, cfg.PacketCount)
		assert.Equal(t, int32(5), *cfg.PacketCount)
		require.NotNil(t, cfg.PacketSize)
		assert.Equal(t, int32(64), *cfg.PacketSize)
		require.NotNil(t, cfg.PacketInterval)
		assert.Equal(t, "1s", *cfg.PacketInterval)
		require.NotNil(t, cfg.PacketTimeout)
		assert.Equal(t, "2s", *cfg.PacketTimeout)
		require.NotNil(t, cfg.UseDNS)
		assert.True(t, *cfg.UseDNS)
		require.NotNil(t, cfg.UseIPv6)
		assert.False(t, *cfg.UseIPv6)
		require.Len(t, cfg.ICMPValidationRules, 1)
		assert.Equal(t, "packetLoss", cfg.ICMPValidationRules[0].Key)
		assert.Equal(t, "LESS_THAN", cfg.ICMPValidationRules[0].Operator)
		assert.Equal(t, int64(5), cfg.ICMPValidationRules[0].Value)
	})

	t.Run("should map ICMP from state with only required fields", func(t *testing.T) {
		model := SyntheticTestModel{
			ID:               types.StringValue("icmp-id"),
			Label:            types.StringValue("ICMP Test"),
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
				"id":           types.StringType,
				"display_name": types.StringType,
			}}),
			ICMP: &ICMPConfigModel{
				MarkSyntheticCall: types.BoolValue(false),
				Retries:           types.Int64Value(0),
				RetryInterval:     types.Int64Value(1),
				Timeout:           types.StringNull(),
				TargetHost:        types.StringValue("ping.example.com"),
				PacketCount:       types.Int64Null(),
				PacketInterval:    types.StringNull(),
				PacketSize:        types.Int64Null(),
				PacketTimeout:     types.StringNull(),
				UseDNS:            types.BoolNull(),
				UseIPv6:           types.BoolNull(),
				ICMPValidationRules: types.SetNull(types.ObjectType{
					AttrTypes: icmpValidationRulesAttrType(),
				}),
			},
		}

		locations, _ := types.SetValueFrom(ctx, types.StringType, []string{"loc-1"})
		model.Locations = locations

		state := &tfsdk.State{Schema: resource.metaData.Schema}
		require.False(t, state.Set(ctx, model).HasError())

		result, diags := resource.MapStateToDataObject(ctx, nil, state)
		require.False(t, diags.HasError())
		require.NotNil(t, result)

		cfg := result.Configuration
		assert.Equal(t, SyntheticTestTypeICMP, cfg.SyntheticType)
		require.NotNil(t, cfg.TargetHost)
		assert.Equal(t, "ping.example.com", *cfg.TargetHost)
		assert.Nil(t, cfg.PacketCount)
		assert.Nil(t, cfg.PacketSize)
		assert.Empty(t, cfg.ICMPValidationRules)
	})
}

func TestUpdateStateICMP(t *testing.T) {
	r := &syntheticTestResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  ResourceInstanaSyntheticTest,
			Schema:        NewSyntheticTestResourceHandle().MetaData().Schema,
			SchemaVersion: 0,
		},
	}
	ctx := context.Background()

	t.Run("should update state from ICMP API object", func(t *testing.T) {
		targetHost := "192.0.2.1"
		packetCount := int32(3)
		packetSize := int32(32)
		packetInterval := "500ms"
		packetTimeout := "1s"
		useDNS := false
		useIPv6 := true

		apiObj := &api.SyntheticTest{
			ID:           "icmp-api-id",
			Label:        "ICMP from API",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-a"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall: true,
				SyntheticType:     SyntheticTestTypeICMP,
				TargetHost:        &targetHost,
				PacketCount:       &packetCount,
				PacketSize:        &packetSize,
				PacketInterval:    &packetInterval,
				PacketTimeout:     &packetTimeout,
				UseDNS:            &useDNS,
				UseIPv6:           &useIPv6,
				ICMPValidationRules: []api.ICMPValidation{
					{Key: "rtt", Operator: "LESS_THAN_OR_EQUALS", Value: 100},
				},
			},
		}

		state := &tfsdk.State{Schema: r.metaData.Schema}
		initializeEmptyState(t, ctx, state)

		diags := r.UpdateState(ctx, state, nil, apiObj)
		require.False(t, diags.HasError())

		var model SyntheticTestModel
		require.False(t, state.Get(ctx, &model).HasError())

		assert.Equal(t, "icmp-api-id", model.ID.ValueString())
		require.NotNil(t, model.ICMP)

		icmp := model.ICMP
		assert.Equal(t, "192.0.2.1", icmp.TargetHost.ValueString())
		assert.Equal(t, int64(3), icmp.PacketCount.ValueInt64())
		assert.Equal(t, int64(32), icmp.PacketSize.ValueInt64())
		assert.Equal(t, "500ms", icmp.PacketInterval.ValueString())
		assert.Equal(t, "1s", icmp.PacketTimeout.ValueString())
		assert.False(t, icmp.UseDNS.ValueBool())
		assert.True(t, icmp.UseIPv6.ValueBool())

		var ruleModels []ICMPValidationModel
		require.False(t, icmp.ICMPValidationRules.ElementsAs(ctx, &ruleModels, false).HasError())
		require.Len(t, ruleModels, 1)
		assert.Equal(t, "rtt", ruleModels[0].Key.ValueString())
		assert.Equal(t, "LESS_THAN_OR_EQUALS", ruleModels[0].Operator.ValueString())
		assert.Equal(t, int64(100), ruleModels[0].Value.ValueInt64())

		// All other config types must be nil
		assert.Nil(t, model.HttpAction)
		assert.Nil(t, model.HttpScript)
		assert.Nil(t, model.BrowserScript)
		assert.Nil(t, model.DNS)
		assert.Nil(t, model.SSLCertificate)
		assert.Nil(t, model.WebpageAction)
		assert.Nil(t, model.WebpageScript)
	})

	t.Run("should update state from ICMP API object without validation rules", func(t *testing.T) {
		targetHost := "ping.example.com"

		apiObj := &api.SyntheticTest{
			ID:           "icmp-no-rules",
			Label:        "ICMP no rules",
			Active:       true,
			PlaybackMode: "Simultaneous",
			Locations:    []string{"loc-a"},
			Configuration: api.SyntheticTestConfig{
				MarkSyntheticCall: false,
				SyntheticType:     SyntheticTestTypeICMP,
				TargetHost:        &targetHost,
			},
		}

		state := &tfsdk.State{Schema: r.metaData.Schema}
		initializeEmptyState(t, ctx, state)

		diags := r.UpdateState(ctx, state, nil, apiObj)
		require.False(t, diags.HasError())

		var model SyntheticTestModel
		require.False(t, state.Get(ctx, &model).HasError())

		require.NotNil(t, model.ICMP)
		assert.Equal(t, "ping.example.com", model.ICMP.TargetHost.ValueString())
		assert.True(t, model.ICMP.ICMPValidationRules.IsNull())
	})
}

func TestValidateSingleConfigTypeICMP(t *testing.T) {
	r := &syntheticTestResource{}

	t.Run("ICMP alone is valid", func(t *testing.T) {
		model := SyntheticTestModel{
			ICMP: &ICMPConfigModel{TargetHost: types.StringValue("host")},
		}
		count, diags := r.validateSingleConfigType(model)
		assert.Equal(t, 1, count)
		assert.False(t, diags.HasError())
	})

	t.Run("ICMP with another type is invalid", func(t *testing.T) {
		model := SyntheticTestModel{
			ICMP: &ICMPConfigModel{TargetHost: types.StringValue("host")},
			HttpAction: &HttpActionConfigModel{
				URL: types.StringValue("https://example.com"),
			},
		}
		count, diags := r.validateSingleConfigType(model)
		assert.Equal(t, 2, count)
		assert.True(t, diags.HasError())
	})
}

func TestBuildICMPSchema(t *testing.T) {
	t.Run("should include icmp in schema", func(t *testing.T) {
		handle := NewSyntheticTestResourceHandle()
		schema := handle.MetaData().Schema
		assert.NotNil(t, schema.Attributes[SyntheticTestFieldICMP])
	})
}

func TestMapICMPValidationRulesToModel(t *testing.T) {
	r := &syntheticTestResource{}

	t.Run("returns null set when rules are empty", func(t *testing.T) {
		result := r.mapICMPValidationRulesToModel(nil)
		assert.True(t, result.IsNull())
	})

	t.Run("maps rules correctly", func(t *testing.T) {
		rules := []api.ICMPValidation{
			{Key: "packetLoss", Operator: "EQUALS", Value: 0},
			{Key: "rtt", Operator: "LESS_THAN", Value: 50},
		}
		result := r.mapICMPValidationRulesToModel(rules)
		assert.False(t, result.IsNull())
		assert.Equal(t, 2, len(result.Elements()))
	})
}
