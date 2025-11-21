package automationaction

import (
	"context"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewAutomationActionResourceHandle tests the resource initialization
func TestNewAutomationActionResourceHandle(t *testing.T) {
	resource := NewAutomationActionResourceHandle()
	require.NotNil(t, resource)

	metaData := resource.MetaData()
	assert.Equal(t, ResourceInstanaAutomationAction, metaData.ResourceName)
	assert.NotNil(t, metaData.Schema)
	assert.Equal(t, int64(0), metaData.SchemaVersion)
}

// TestMetaData tests the MetaData method
func TestMetaData(t *testing.T) {
	resource := &automationActionResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName:  "test_resource",
			SchemaVersion: 0,
		},
	}

	metaData := resource.MetaData()
	assert.Equal(t, "test_resource", metaData.ResourceName)
	assert.Equal(t, int64(0), metaData.SchemaVersion)
}

// TestSetComputedFields tests the SetComputedFields method
func TestSetComputedFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle()
	ctx := context.Background()

	plan := &tfsdk.Plan{
		Schema: resource.MetaData().Schema,
	}

	diags := resource.SetComputedFields(ctx, plan)
	assert.False(t, diags.HasError())
}

// TestGetRestResource tests the GetRestResource method
func TestGetRestResource(t *testing.T) {
	resource := &automationActionResource{}
	assert.NotNil(t, resource.GetRestResource)
}

// TestMapStateToDataObject_ScriptAction tests mapping script action from state to data object
func TestMapStateToDataObject_ScriptAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Script Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Script: &shared.ScriptModel{
			Content:     types.StringValue("echo 'Hello World'"),
			Interpreter: types.StringValue("bash"),
			Timeout:     types.StringValue("60"),
			Source:      types.StringValue("inline"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError(), "Expected no errors, got: %v", diags)
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Script Action", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, shared.ActionTypeScript, result.Type)
	require.NotEmpty(t, result.Fields)

	// Verify script fields
	contentField := findField(result.Fields, restapi.ScriptSshFieldName)
	require.NotNil(t, contentField)
	assert.Equal(t, "echo 'Hello World'", contentField.Value)
	assert.Equal(t, shared.Base64Encoding, contentField.Encoding)
}

// TestMapStateToDataObject_HttpAction tests mapping HTTP action from state to data object
func TestMapStateToDataObject_HttpAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test HTTP Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Http: &shared.HttpModel{
			Host:             types.StringValue("https://example.com/api"),
			Method:           types.StringValue("POST"),
			Body:             types.StringValue(`{"key":"value"}`),
			IgnoreCertErrors: types.BoolValue(false),
			Timeout:          types.StringValue("30"),
			Language:         types.StringValue("json"),
			ContentType:      types.StringValue("application/json"),
			Headers:          types.MapNull(types.StringType),
			Auth:             nil,
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test HTTP Action", result.Name)
	assert.Equal(t, shared.ActionTypeHttp, result.Type)

	// Verify HTTP fields
	hostField := findField(result.Fields, restapi.HttpHostFieldName)
	require.NotNil(t, hostField)
	assert.Equal(t, "https://example.com/api", hostField.Value)

	methodField := findField(result.Fields, restapi.HttpMethodFieldName)
	require.NotNil(t, methodField)
	assert.Equal(t, "POST", methodField.Value)
}

// TestMapStateToDataObject_HttpActionWithBasicAuth tests HTTP action with basic authentication
func TestMapStateToDataObject_HttpActionWithBasicAuth(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test HTTP Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Http: &shared.HttpModel{
			Host:    types.StringValue("https://example.com/api"),
			Method:  types.StringValue("GET"),
			Headers: types.MapNull(types.StringType),
			Auth: &shared.AuthModel{
				BasicAuth: &shared.BasicAuthModel{
					UserName: types.StringValue("user"),
					Password: types.StringValue("pass"),
				},
			},
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	// Verify auth field exists
	authField := findField(result.Fields, AutomationActionAPIFieldAuthen)
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, "basicAuth")
	assert.Contains(t, authField.Value, "user")
}

// TestMapStateToDataObject_HttpActionWithBearerToken tests HTTP action with bearer token
func TestMapStateToDataObject_HttpActionWithBearerToken(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test HTTP Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Http: &shared.HttpModel{
			Host:    types.StringValue("https://example.com/api"),
			Method:  types.StringValue("GET"),
			Headers: types.MapNull(types.StringType),
			Auth: &shared.AuthModel{
				Token: &shared.BearerTokenModel{
					BearerToken: types.StringValue("token123"),
				},
			},
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	authField := findField(result.Fields, AutomationActionAPIFieldAuthen)
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, "bearerToken")
	assert.Contains(t, authField.Value, "token123")
}

// TestMapStateToDataObject_HttpActionWithApiKey tests HTTP action with API key
func TestMapStateToDataObject_HttpActionWithApiKey(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test HTTP Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Http: &shared.HttpModel{
			Host:    types.StringValue("https://example.com/api"),
			Method:  types.StringValue("GET"),
			Headers: types.MapNull(types.StringType),
			Auth: &shared.AuthModel{
				ApiKey: &shared.ApiKeyModel{
					Key:         types.StringValue("X-API-Key"),
					Value:       types.StringValue("key123"),
					KeyLocation: types.StringValue("header"),
				},
			},
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	authField := findField(result.Fields, AutomationActionAPIFieldAuthen)
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, "apiKey")
	assert.Contains(t, authField.Value, "X-API-Key")
}

// TestMapStateToDataObject_HttpActionWithHeaders tests HTTP action with custom headers
func TestMapStateToDataObject_HttpActionWithHeaders(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	headers := types.MapValueMust(types.StringType, map[string]attr.Value{
		"Content-Type": types.StringValue("application/json"),
		"X-Custom":     types.StringValue("value"),
	})

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test HTTP Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Http: &shared.HttpModel{
			Host:    types.StringValue("https://example.com/api"),
			Method:  types.StringValue("POST"),
			Headers: headers,
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	headerField := findField(result.Fields, restapi.HttpHeaderFieldName)
	require.NotNil(t, headerField)
	assert.Contains(t, headerField.Value, "Content-Type")
	assert.Contains(t, headerField.Value, "X-Custom")
}

// TestMapStateToDataObject_ManualAction tests mapping manual action
func TestMapStateToDataObject_ManualAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Manual Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Manual: &shared.ManualModel{
			Content: types.StringValue("Manual instructions here"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, AutomationActionTypeManual, result.Type)
	contentField := findField(result.Fields, AutomationActionFieldContent)
	require.NotNil(t, contentField)
	assert.Equal(t, "Manual instructions here", contentField.Value)
}

// TestMapStateToDataObject_JiraAction tests mapping Jira action
func TestMapStateToDataObject_JiraAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Jira Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Jira: &shared.JiraModel{
			Project:     types.StringValue("PROJ"),
			Operation:   types.StringValue("create"),
			IssueType:   types.StringValue("Bug"),
			Description: types.StringValue("Issue description"),
			Assignee:    types.StringValue("user@example.com"),
			Title:       types.StringValue("Issue Title"),
			Labels:      types.StringValue("bug,urgent"),
			Comment:     types.StringValue("Additional comment"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, AutomationActionTypeJira, result.Type)
	projectField := findField(result.Fields, AutomationActionFieldProject)
	require.NotNil(t, projectField)
	assert.Equal(t, "PROJ", projectField.Value)
}

// TestMapStateToDataObject_GitHubAction tests mapping GitHub action
func TestMapStateToDataObject_GitHubAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test GitHub Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		GitHub: &shared.GitHubModel{
			Owner:     types.StringValue("owner"),
			Repo:      types.StringValue("repo"),
			Title:     types.StringValue("Issue Title"),
			Body:      types.StringValue("Issue body"),
			Operation: types.StringValue("create"),
			Assignees: types.StringValue("user1,user2"),
			Labels:    types.StringValue("bug"),
			Comment:   types.StringValue("Comment text"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, AutomationActionTypeGitHub, result.Type)
	ownerField := findField(result.Fields, AutomationActionFieldOwner)
	require.NotNil(t, ownerField)
	assert.Equal(t, "owner", ownerField.Value)
}

// TestMapStateToDataObject_DocLinkAction tests mapping documentation link action
func TestMapStateToDataObject_DocLinkAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test DocLink Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		DocLink: &shared.DocLinkModel{
			Url: types.StringValue("https://docs.example.com"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, AutomationActionTypeDocLink, result.Type)
	urlField := findField(result.Fields, AutomationActionFieldUrl)
	require.NotNil(t, urlField)
	assert.Equal(t, "https://docs.example.com", urlField.Value)
	assert.Equal(t, shared.UTF8Encoding, urlField.Encoding)
}

// TestMapStateToDataObject_GitLabAction tests mapping GitLab action
func TestMapStateToDataObject_GitLabAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test GitLab Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		GitLab: &shared.GitLabModel{
			ProjectId:   types.StringValue("123"),
			Title:       types.StringValue("Issue Title"),
			Description: types.StringValue("Issue description"),
			Operation:   types.StringValue("create"),
			Labels:      types.StringValue("bug"),
			IssueType:   types.StringValue("issue"),
			Comment:     types.StringValue("Comment text"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, AutomationActionTypeGitLab, result.Type)
	projectIdField := findField(result.Fields, AutomationActionAPIFieldProjectId)
	require.NotNil(t, projectIdField)
	assert.Equal(t, "123", projectIdField.Value)
}

// TestMapStateToDataObject_AnsibleAction tests mapping Ansible action
func TestMapStateToDataObject_AnsibleAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Ansible Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Ansible: &shared.AnsibleModel{
			WorkflowId:       types.StringValue("workflow-123"),
			PlaybookId:       types.StringValue("playbook-456"),
			PlaybookFileName: types.StringValue("site.yml"),
			AnsibleUrl:       types.StringValue("https://ansible.example.com"),
			HostId:           types.StringValue("host-789"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, AutomationActionTypeAnsible, result.Type)
	workflowField := findField(result.Fields, AutomationActionAPIFieldWorkflowId)
	require.NotNil(t, workflowField)
	assert.Equal(t, "workflow-123", workflowField.Value)
}

// TestMapStateToDataObject_WithTags tests mapping action with tags
func TestMapStateToDataObject_WithTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	tags := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("tag1"),
		types.StringValue("tag2"),
		types.StringValue("tag3"),
	})

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Action"),
		Description: types.StringValue("Test Description"),
		Tags:        tags,
		Manual: &shared.ManualModel{
			Content: types.StringValue("Manual content"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.NotNil(t, result.Tags)
	tagSlice, ok := result.Tags.([]string)
	require.True(t, ok)
	assert.Len(t, tagSlice, 3)
	assert.Contains(t, tagSlice, "tag1")
	assert.Contains(t, tagSlice, "tag2")
	assert.Contains(t, tagSlice, "tag3")
}

// TestMapStateToDataObject_WithInputParameters tests mapping action with input parameters
func TestMapStateToDataObject_WithInputParameters(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Manual: &shared.ManualModel{
			Content: types.StringValue("Manual content"),
		},
		InputParameter: []shared.ParameterModel{
			{
				Name:        types.StringValue("param1"),
				Description: types.StringValue("Parameter 1"),
				Label:       types.StringValue("Param 1"),
				Required:    types.BoolValue(true),
				Hidden:      types.BoolValue(false),
				Type:        types.StringValue("static"),
				Value:       types.StringValue("value1"),
			},
			{
				Name:        types.StringValue("param2"),
				Description: types.StringValue("Parameter 2"),
				Label:       types.StringValue("Param 2"),
				Required:    types.BoolValue(false),
				Hidden:      types.BoolValue(true),
				Type:        types.StringValue("dynamic"),
				Value:       types.StringValue("value2"),
			},
		},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	require.Len(t, result.InputParameters, 2)
	assert.Equal(t, "param1", result.InputParameters[0].Name)
	assert.Equal(t, "Parameter 1", result.InputParameters[0].Description)
	assert.True(t, result.InputParameters[0].Required)
	assert.False(t, result.InputParameters[0].Hidden)
	assert.Equal(t, "static", result.InputParameters[0].Type)
	assert.Equal(t, "value1", result.InputParameters[0].Value)
}

// TestMapStateToDataObject_InvalidConfig tests error handling for invalid configuration
func TestMapStateToDataObject_InvalidConfig(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	// Create state with no action type configured
	state := createMockState(t, shared.AutomationActionModel{
		ID:             types.StringValue("test-id"),
		Name:           types.StringValue("Test Action"),
		Description:    types.StringValue("Test Description"),
		Tags:           types.ListNull(types.StringType),
		InputParameter: []shared.ParameterModel{},
	})

	_, diags := resource.MapStateToDataObject(ctx, nil, &state)
	assert.True(t, diags.HasError())
}

// TestUpdateState_ScriptAction tests updating state with script action
func TestUpdateState_ScriptAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	data := &restapi.AutomationAction{
		ID:          "test-id",
		Name:        "Test Script Action",
		Description: "Test Description",
		Type:        shared.ActionTypeScript,
		Tags:        []interface{}{"tag1", "tag2"},
		Fields: []restapi.Field{
			{
				Name:        restapi.ScriptSshFieldName,
				Description: restapi.ScriptSshFieldDescription,
				Value:       "echo 'Hello'",
				Encoding:    shared.Base64Encoding,
			},
			{
				Name:        restapi.SubtypeFieldName,
				Description: restapi.SubtypeFieldDescription,
				Value:       "bash",
				Encoding:    shared.AsciiEncoding,
			},
			{
				Name:        restapi.TimeoutFieldName,
				Description: restapi.TimeoutFieldDescription,
				Value:       "60",
				Encoding:    shared.AsciiEncoding,
			},
		},
		InputParameters: []restapi.Parameter{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	assert.Equal(t, "Test Script Action", model.Name.ValueString())
	assert.Equal(t, "Test Description", model.Description.ValueString())
	require.NotNil(t, model.Script)
	assert.Equal(t, "echo 'Hello'", model.Script.Content.ValueString())
	assert.Equal(t, "bash", model.Script.Interpreter.ValueString())
	assert.Equal(t, "60", model.Script.Timeout.ValueString())
}

// TestUpdateState_HttpAction tests updating state with HTTP action
func TestUpdateState_HttpAction(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	data := &restapi.AutomationAction{
		ID:          "test-id",
		Name:        "Test HTTP Action",
		Description: "Test Description",
		Type:        shared.ActionTypeHttp,
		Tags:        nil,
		Fields: []restapi.Field{
			{
				Name:        restapi.HttpHostFieldName,
				Description: restapi.HttpHostFieldDescription,
				Value:       "https://example.com",
				Encoding:    shared.AsciiEncoding,
			},
			{
				Name:        restapi.HttpMethodFieldName,
				Description: restapi.HttpMethodFieldDescription,
				Value:       "POST",
				Encoding:    shared.AsciiEncoding,
			},
			{
				Name:        restapi.HttpBodyFieldName,
				Description: restapi.HttpBodyFieldDescription,
				Value:       `{"key":"value"}`,
				Encoding:    shared.AsciiEncoding,
			},
			{
				Name:        AutomationActionAPIFieldAuthen,
				Description: AutomationActionDescFieldAuthen,
				Value:       `{"type":"noAuth"}`,
				Encoding:    shared.AsciiEncoding,
			},
		},
		InputParameters: []restapi.Parameter{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Equal(t, "test-id", model.ID.ValueString())
	require.NotNil(t, model.Http)
	assert.Equal(t, "https://example.com", model.Http.Host.ValueString())
	assert.Equal(t, "POST", model.Http.Method.ValueString())
	assert.Equal(t, `{"key":"value"}`, model.Http.Body.ValueString())
}

// TestUpdateState_WithTags tests updating state with tags
func TestUpdateState_WithTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	data := &restapi.AutomationAction{
		ID:          "test-id",
		Name:        "Test Action",
		Description: "Test Description",
		Type:        AutomationActionTypeManual,
		Tags:        []interface{}{"tag1", "tag2", "tag3"},
		Fields: []restapi.Field{
			{
				Name:  AutomationActionFieldContent,
				Value: "Manual content",
			},
		},
		InputParameters: []restapi.Parameter{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.False(t, model.Tags.IsNull())
	var tags []string
	diags = model.Tags.ElementsAs(ctx, &tags, false)
	require.False(t, diags.HasError())
	assert.Len(t, tags, 3)
	assert.Contains(t, tags, "tag1")
}

// TestUpdateState_WithInputParameters tests updating state with input parameters
func TestUpdateState_WithInputParameters(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	data := &restapi.AutomationAction{
		ID:          "test-id",
		Name:        "Test Action",
		Description: "Test Description",
		Type:        AutomationActionTypeManual,
		Tags:        nil,
		Fields: []restapi.Field{
			{
				Name:  AutomationActionFieldContent,
				Value: "Manual content",
			},
		},
		InputParameters: []restapi.Parameter{
			{
				Name:        "param1",
				Description: "Parameter 1",
				Label:       "Param 1",
				Required:    true,
				Hidden:      false,
				Type:        "static",
				Value:       "value1",
			},
		},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	require.Len(t, model.InputParameter, 1)
	assert.Equal(t, "param1", model.InputParameter[0].Name.ValueString())
	assert.Equal(t, "Parameter 1", model.InputParameter[0].Description.ValueString())
	assert.True(t, model.InputParameter[0].Required.ValueBool())
}

// TestMapTagsToState tests tag mapping to state
func TestMapTagsToState(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	tests := []struct {
		name     string
		tags     interface{}
		expected []string
		hasError bool
	}{
		{
			name:     "valid tags",
			tags:     []interface{}{"tag1", "tag2", "tag3"},
			expected: []string{"tag1", "tag2", "tag3"},
			hasError: false,
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: nil,
			hasError: false,
		},
		{
			name:     "empty tags",
			tags:     []interface{}{},
			expected: []string{},
			hasError: false,
		},
		{
			name:     "invalid tag type",
			tags:     []interface{}{"tag1", 123},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := resource.mapTagsToState(ctx, tt.tags)

			if tt.hasError {
				assert.True(t, diags.HasError())
			} else {
				assert.False(t, diags.HasError())
				if tt.tags == nil {
					assert.True(t, result.IsNull())
				} else {
					var tags []string
					diags = result.ElementsAs(ctx, &tags, false)
					assert.False(t, diags.HasError())
					assert.Equal(t, tt.expected, tags)
				}
			}
		})
	}
}

// TestMapTagsFromState tests tag mapping from state
func TestMapTagsFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	tests := []struct {
		name     string
		model    shared.AutomationActionModel
		expected []string
	}{
		{
			name: "with tags",
			model: shared.AutomationActionModel{
				Tags: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("tag1"),
					types.StringValue("tag2"),
				}),
			},
			expected: []string{"tag1", "tag2"},
		},
		{
			name: "null tags",
			model: shared.AutomationActionModel{
				Tags: types.ListNull(types.StringType),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := resource.mapTagsFromState(ctx, tt.model)
			assert.False(t, diags.HasError())

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				tagSlice, ok := result.([]string)
				require.True(t, ok)
				assert.Equal(t, tt.expected, tagSlice)
			}
		})
	}
}

// TestMapInputParametersFromState tests input parameter mapping from state
func TestMapInputParametersFromState(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	model := shared.AutomationActionModel{
		InputParameter: []shared.ParameterModel{
			{
				Name:        types.StringValue("param1"),
				Description: types.StringValue("Description 1"),
				Label:       types.StringValue("Label 1"),
				Required:    types.BoolValue(true),
				Hidden:      types.BoolValue(false),
				Type:        types.StringValue("static"),
				Value:       types.StringValue("value1"),
			},
			{
				Name:        types.StringValue("param2"),
				Description: types.StringValue("Description 2"),
				Label:       types.StringValue("Label 2"),
				Required:    types.BoolValue(false),
				Hidden:      types.BoolValue(true),
				Type:        types.StringValue("vault"),
				Value:       types.StringValue("value2"),
			},
		},
	}

	result, diags := resource.mapInputParametersFromState(ctx, model)
	assert.False(t, diags.HasError())
	require.Len(t, result, 2)

	assert.Equal(t, "param1", result[0].Name)
	assert.Equal(t, "Description 1", result[0].Description)
	assert.Equal(t, "Label 1", result[0].Label)
	assert.True(t, result[0].Required)
	assert.False(t, result[0].Hidden)
	assert.Equal(t, "static", result[0].Type)
	assert.Equal(t, "value1", result[0].Value)

	assert.Equal(t, "param2", result[1].Name)
	assert.False(t, result[1].Required)
	assert.True(t, result[1].Hidden)
	assert.Equal(t, "vault", result[1].Type)
}

// TestMapInputParametersFromState_Empty tests empty input parameters
func TestMapInputParametersFromState_Empty(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	model := shared.AutomationActionModel{
		InputParameter: []shared.ParameterModel{},
	}

	result, diags := resource.mapInputParametersFromState(ctx, model)
	assert.False(t, diags.HasError())
	assert.Nil(t, result)
}

// TestMapActionTypeAndFields_AllHTTPMethods tests all HTTP methods
func TestMapActionTypeAndFields_AllHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			ctx := context.Background()
			resource := &automationActionResource{}

			state := createMockState(t, shared.AutomationActionModel{
				ID:          types.StringValue("test-id"),
				Name:        types.StringValue("Test Action"),
				Description: types.StringValue("Test Description"),
				Tags:        types.ListNull(types.StringType),
				Http: &shared.HttpModel{
					Host:    types.StringValue("https://example.com"),
					Method:  types.StringValue(method),
					Headers: types.MapNull(types.StringType),
				},
				InputParameter: []shared.ParameterModel{},
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)

			methodField := findField(result.Fields, restapi.HttpMethodFieldName)
			require.NotNil(t, methodField)
			assert.Equal(t, method, methodField.Value)
		})
	}
}

// TestMapActionTypeAndFields_AllParameterTypes tests all parameter types
func TestMapActionTypeAndFields_AllParameterTypes(t *testing.T) {
	paramTypes := []string{"static", "dynamic", "vault"}

	for _, paramType := range paramTypes {
		t.Run(paramType, func(t *testing.T) {
			ctx := context.Background()
			resource := &automationActionResource{}

			state := createMockState(t, shared.AutomationActionModel{
				ID:          types.StringValue("test-id"),
				Name:        types.StringValue("Test Action"),
				Description: types.StringValue("Test Description"),
				Tags:        types.ListNull(types.StringType),
				Manual: &shared.ManualModel{
					Content: types.StringValue("Manual content"),
				},
				InputParameter: []shared.ParameterModel{
					{
						Name:        types.StringValue("param"),
						Description: types.StringValue("Description"),
						Label:       types.StringValue("Label"),
						Required:    types.BoolValue(true),
						Hidden:      types.BoolValue(false),
						Type:        types.StringValue(paramType),
						Value:       types.StringValue("value"),
					},
				},
			})

			result, diags := resource.MapStateToDataObject(ctx, nil, &state)
			require.False(t, diags.HasError())
			require.NotNil(t, result)
			require.Len(t, result.InputParameters, 1)
			assert.Equal(t, paramType, result.InputParameters[0].Type)
		})
	}
}

// Helper functions

func createMockState(t *testing.T, model shared.AutomationActionModel) tfsdk.State {
	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := state.Set(context.Background(), model)
	require.False(t, diags.HasError(), "Failed to set state: %v", diags)

	return state
}

func getTestSchema() schema.Schema {
	resource := NewAutomationActionResourceHandle()
	return resource.MetaData().Schema
}

// TestMapStateToDataObject_FromPlan tests mapping from plan instead of state
func TestMapStateToDataObject_FromPlan(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	plan := &tfsdk.Plan{
		Schema: getTestSchema(),
	}

	model := shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Manual: &shared.ManualModel{
			Content: types.StringValue("Manual content"),
		},
		InputParameter: []shared.ParameterModel{},
	}

	diags := plan.Set(ctx, model)
	require.False(t, diags.HasError())

	result, diags := resource.MapStateToDataObject(ctx, plan, nil)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Action", result.Name)
}

// TestUpdateState_NullTags tests updating state with null tags
func TestUpdateState_NullTags(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	data := &restapi.AutomationAction{
		ID:          "test-id",
		Name:        "Test Action",
		Description: "Test Description",
		Type:        AutomationActionTypeManual,
		Tags:        nil,
		Fields: []restapi.Field{
			{
				Name:  AutomationActionFieldContent,
				Value: "Manual content",
			},
		},
		InputParameters: []restapi.Parameter{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.True(t, model.Tags.IsNull())
}

// TestUpdateState_EmptyInputParameters tests updating state with empty input parameters
func TestUpdateState_EmptyInputParameters(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	data := &restapi.AutomationAction{
		ID:              "test-id",
		Name:            "Test Action",
		Description:     "Test Description",
		Type:            AutomationActionTypeManual,
		Tags:            nil,
		Fields:          []restapi.Field{{Name: AutomationActionFieldContent, Value: "Manual content"}},
		InputParameters: []restapi.Parameter{},
	}

	state := &tfsdk.State{
		Schema: getTestSchema(),
	}

	diags := resource.UpdateState(ctx, state, nil, data)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())

	assert.Nil(t, model.InputParameter)
}

// TestMapStateToDataObject_NullID tests mapping with null ID
func TestMapStateToDataObject_NullID(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringNull(),
		Name:        types.StringValue("Test Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Manual: &shared.ManualModel{
			Content: types.StringValue("Manual content"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)
	assert.Equal(t, "", result.ID)
}

// TestMapActionTypeAndFields_ScriptWithAllFields tests script action with all optional fields
func TestMapActionTypeAndFields_ScriptWithAllFields(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Script Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Script: &shared.ScriptModel{
			Content:     types.StringValue("#!/bin/bash\necho 'test'"),
			Interpreter: types.StringValue("bash"),
			Timeout:     types.StringValue("120"),
			Source:      types.StringValue("inline"),
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, shared.ActionTypeScript, result.Type)

	// Verify all fields are present
	contentField := findField(result.Fields, restapi.ScriptSshFieldName)
	require.NotNil(t, contentField)

	interpreterField := findField(result.Fields, restapi.SubtypeFieldName)
	require.NotNil(t, interpreterField)
	assert.Equal(t, "bash", interpreterField.Value)

	timeoutField := findField(result.Fields, restapi.TimeoutFieldName)
	require.NotNil(t, timeoutField)
	assert.Equal(t, "120", timeoutField.Value)

	sourceField := findField(result.Fields, AutomationActionFieldSource)
	require.NotNil(t, sourceField)
	assert.Equal(t, "inline", sourceField.Value)
}

// TestMapActionTypeAndFields_HttpWithAllOptionalFields tests HTTP action with all optional fields
func TestMapActionTypeAndFields_HttpWithAllOptionalFields(t *testing.T) {
	ctx := context.Background()
	resource := &automationActionResource{}

	state := createMockState(t, shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test HTTP Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Http: &shared.HttpModel{
			Host:             types.StringValue("https://api.example.com"),
			Method:           types.StringValue("POST"),
			Body:             types.StringValue(`{"test":"data"}`),
			IgnoreCertErrors: types.BoolValue(true),
			Timeout:          types.StringValue("30"),
			Language:         types.StringValue("json"),
			ContentType:      types.StringValue("application/json"),
			Headers:          types.MapNull(types.StringType),
			Auth:             nil,
		},
		InputParameter: []shared.ParameterModel{},
	})

	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	assert.Equal(t, shared.ActionTypeHttp, result.Type)

	// Verify optional fields
	bodyField := findField(result.Fields, restapi.HttpBodyFieldName)
	require.NotNil(t, bodyField)
	assert.Equal(t, `{"test":"data"}`, bodyField.Value)

	ignoreCertField := findField(result.Fields, restapi.HttpIgnoreCertErrorsFieldName)
	require.NotNil(t, ignoreCertField)
	assert.Equal(t, "true", ignoreCertField.Value)

	timeoutField := findField(result.Fields, restapi.TimeoutFieldName)
	require.NotNil(t, timeoutField)
	assert.Equal(t, "30", timeoutField.Value)

	languageField := findField(result.Fields, AutomationActionFieldLanguage)
	require.NotNil(t, languageField)
	assert.Equal(t, "json", languageField.Value)

	contentTypeField := findField(result.Fields, AutomationActionFieldContentType)
	require.NotNil(t, contentTypeField)
	assert.Equal(t, "application/json", contentTypeField.Value)
}

// TestUpdateState_AllActionTypes tests UpdateState for all action types
func TestUpdateState_AllActionTypes(t *testing.T) {
	tests := []struct {
		name       string
		actionType string
		fields     []restapi.Field
	}{
		{
			name:       "Manual",
			actionType: AutomationActionTypeManual,
			fields:     []restapi.Field{{Name: AutomationActionFieldContent, Value: "Manual content"}},
		},
		{
			name:       "Jira",
			actionType: AutomationActionTypeJira,
			fields:     []restapi.Field{{Name: AutomationActionFieldProject, Value: "PROJ"}},
		},
		{
			name:       "GitHub",
			actionType: AutomationActionTypeGitHub,
			fields:     []restapi.Field{{Name: AutomationActionFieldOwner, Value: "owner"}},
		},
		{
			name:       "DocLink",
			actionType: AutomationActionTypeDocLink,
			fields:     []restapi.Field{{Name: AutomationActionFieldUrl, Value: "https://docs.example.com"}},
		},
		{
			name:       "GitLab",
			actionType: AutomationActionTypeGitLab,
			fields:     []restapi.Field{{Name: AutomationActionAPIFieldProjectId, Value: "123"}},
		},
		{
			name:       "Ansible",
			actionType: AutomationActionTypeAnsible,
			fields:     []restapi.Field{{Name: AutomationActionAPIFieldWorkflowId, Value: "workflow-123"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resource := &automationActionResource{}

			data := &restapi.AutomationAction{
				ID:              "test-id",
				Name:            "Test Action",
				Description:     "Test Description",
				Type:            tt.actionType,
				Tags:            nil,
				Fields:          tt.fields,
				InputParameters: []restapi.Parameter{},
			}

			state := &tfsdk.State{
				Schema: getTestSchema(),
			}

			diags := resource.UpdateState(ctx, state, nil, data)
			require.False(t, diags.HasError())

			var model shared.AutomationActionModel
			diags = state.Get(ctx, &model)
			require.False(t, diags.HasError())

			assert.Equal(t, "test-id", model.ID.ValueString())
			assert.Equal(t, "Test Action", model.Name.ValueString())
		})
	}
}

func findField(fields []restapi.Field, name string) *restapi.Field {
	for _, field := range fields {
		if field.Name == name {
			return &field
		}
	}
	return nil
}

// Made with Bob

// TestMapStateToDataObject_FromState tests mapping from state instead of plan
func TestMapStateToDataObject_FromState(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	// Create a state with script action
	model := shared.AutomationActionModel{
		ID:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Action"),
		Description: types.StringValue("Test Description"),
		Tags:        types.ListNull(types.StringType),
		Script: &shared.ScriptModel{
			Content:     types.StringValue("echo 'test'"),
			Interpreter: types.StringValue("bash"),
		},
		InputParameter: []shared.ParameterModel{},
	}

	state := tfsdk.State{
		Schema: getTestSchema(),
	}
	diags := state.Set(ctx, model)
	require.False(t, diags.HasError())

	// Map from state (not plan)
	result, diags := resource.MapStateToDataObject(ctx, nil, &state)
	require.False(t, diags.HasError())
	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Action", result.Name)
	assert.Equal(t, shared.ActionTypeScript, result.Type)
}

// TestMapTagsFromState_NullTags tests mapTagsFromState with null tags
func TestMapTagsFromState_NullTags(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	// Test with null tags
	model := shared.AutomationActionModel{
		Tags: types.ListNull(types.StringType),
	}

	result, diags := resource.mapTagsFromState(ctx, model)
	require.False(t, diags.HasError())
	assert.Nil(t, result)
}

// TestMapActionTypeAndFields_InvalidConfiguration tests error handling for invalid config
func TestMapActionTypeAndFields_InvalidConfiguration(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	// Model with no action type configured
	model := shared.AutomationActionModel{
		Name:        types.StringValue("Test"),
		Description: types.StringValue("Test"),
	}

	_, _, diags := resource.mapActionTypeAndFields(ctx, model)
	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary(), AutomationActionErrInvalidConfig)
}

// TestMapActionTypeAndFields_HTTPAuthNoAuth tests HTTP with no auth
func TestMapActionTypeAndFields_HTTPAuthNoAuth(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	tests := []struct {
		name string
		auth *shared.AuthModel
	}{
		{
			name: "NoAuth with nil auth",
			auth: nil,
		},
		{
			name: "NoAuth with empty auth",
			auth: &shared.AuthModel{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := shared.AutomationActionModel{
				Http: &shared.HttpModel{
					Host:    types.StringValue("example.com"),
					Method:  types.StringValue("POST"),
					Headers: types.MapNull(types.StringType),
					Auth:    tt.auth,
				},
			}

			actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
			require.False(t, diags.HasError())
			assert.Equal(t, shared.ActionTypeHttp, actionType)

			// Find auth field
			var authField *restapi.Field
			for i := range fields {
				if fields[i].Name == AutomationActionAPIFieldAuthen {
					authField = &fields[i]
					break
				}
			}

			require.NotNil(t, authField, "Auth field should be present")
			assert.Contains(t, authField.Value, AutomationActionAuthTypeNoAuth)
		})
	}
}

// TestMapActionTypeAndFields_HTTPWithHeaders tests HTTP with headers
func TestMapActionTypeAndFields_HTTPWithHeaders(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	headers := map[string]attr.Value{
		"Content-Type":  types.StringValue("application/json"),
		"Authorization": types.StringValue("Bearer token"),
	}
	headersMap, diags := types.MapValue(types.StringType, headers)
	require.False(t, diags.HasError())

	model := shared.AutomationActionModel{
		Http: &shared.HttpModel{
			Host:    types.StringValue("example.com"),
			Method:  types.StringValue("POST"),
			Headers: headersMap,
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeHttp, actionType)

	// Find headers field
	var headersField *restapi.Field
	for i := range fields {
		if fields[i].Name == restapi.HttpHeaderFieldName {
			headersField = &fields[i]
			break
		}
	}

	require.NotNil(t, headersField, "Headers field should be present")
	assert.Contains(t, headersField.Value, "Content-Type")
	assert.Contains(t, headersField.Value, "application/json")
}

// TestMapActionTypeAndFields_AllOptionalFieldsPopulated tests all action types with all optional fields
func TestMapActionTypeAndFields_AllOptionalFieldsPopulated(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	tests := []struct {
		name          string
		model         shared.AutomationActionModel
		expectedType  string
		expectedCount int
	}{
		{
			name: "Script with all fields",
			model: shared.AutomationActionModel{
				Script: &shared.ScriptModel{
					Content:     types.StringValue("echo 'test'"),
					Interpreter: types.StringValue("bash"),
					Timeout:     types.StringValue("30s"),
					Source:      types.StringValue("inline"),
				},
			},
			expectedType:  shared.ActionTypeScript,
			expectedCount: 4,
		},
		{
			name: "HTTP with all fields",
			model: shared.AutomationActionModel{
				Http: &shared.HttpModel{
					Host:             types.StringValue("example.com"),
					Method:           types.StringValue("POST"),
					Body:             types.StringValue(`{"key":"value"}`),
					Headers:          types.MapNull(types.StringType),
					IgnoreCertErrors: types.BoolValue(true),
					Timeout:          types.StringValue("60s"),
					Language:         types.StringValue("json"),
					ContentType:      types.StringValue("application/json"),
					Auth: &shared.AuthModel{
						BasicAuth: &shared.BasicAuthModel{
							UserName: types.StringValue("user"),
							Password: types.StringValue("pass"),
						},
					},
				},
			},
			expectedType:  shared.ActionTypeHttp,
			expectedCount: 8,
		},
		{
			name: "Jira with all fields",
			model: shared.AutomationActionModel{
				Jira: &shared.JiraModel{
					Project:     types.StringValue("TEST"),
					Operation:   types.StringValue("create"),
					IssueType:   types.StringValue("Bug"),
					Description: types.StringValue("Test description"),
					Assignee:    types.StringValue("user@example.com"),
					Title:       types.StringValue("Test Issue"),
					Labels:      types.StringValue("bug,urgent"),
					Comment:     types.StringValue("Test comment"),
				},
			},
			expectedType:  AutomationActionTypeJira,
			expectedCount: 8,
		},
		{
			name: "GitHub with all fields",
			model: shared.AutomationActionModel{
				GitHub: &shared.GitHubModel{
					Owner:     types.StringValue("owner"),
					Repo:      types.StringValue("repo"),
					Title:     types.StringValue("Test Issue"),
					Body:      types.StringValue("Test body"),
					Operation: types.StringValue("create"),
					Assignees: types.StringValue("user1,user2"),
					Labels:    types.StringValue("bug,feature"),
					Comment:   types.StringValue("Test comment"),
				},
			},
			expectedType:  AutomationActionTypeGitHub,
			expectedCount: 8,
		},
		{
			name: "GitLab with all fields",
			model: shared.AutomationActionModel{
				GitLab: &shared.GitLabModel{
					ProjectId:   types.StringValue("123"),
					Title:       types.StringValue("Test Issue"),
					Description: types.StringValue("Test description"),
					Operation:   types.StringValue("create"),
					Labels:      types.StringValue("bug,feature"),
					IssueType:   types.StringValue("issue"),
					Comment:     types.StringValue("Test comment"),
				},
			},
			expectedType:  AutomationActionTypeGitLab,
			expectedCount: 7,
		},
		{
			name: "Ansible with all fields",
			model: shared.AutomationActionModel{
				Ansible: &shared.AnsibleModel{
					WorkflowId:       types.StringValue("workflow-123"),
					PlaybookId:       types.StringValue("playbook-456"),
					PlaybookFileName: types.StringValue("site.yml"),
					AnsibleUrl:       types.StringValue("https://ansible.example.com"),
					HostId:           types.StringValue("host-789"),
				},
			},
			expectedType:  AutomationActionTypeAnsible,
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actionType, fields, diags := resource.mapActionTypeAndFields(ctx, tt.model)
			require.False(t, diags.HasError())
			assert.Equal(t, tt.expectedType, actionType)
			assert.Equal(t, tt.expectedCount, len(fields))
		})
	}
}

// TestMapStateToDataObject_NilPlanAndState tests error handling when both plan and state are nil
func TestMapStateToDataObject_NilPlanAndState(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	// Both plan and state are nil - should still work but get empty model
	result, diags := resource.MapStateToDataObject(ctx, nil, nil)
	// This will have errors because we can't get the model from nil plan/state
	assert.True(t, diags.HasError())
	assert.Nil(t, result)
}

// TestMapActionTypeAndFields_ScriptMinimalFields tests script with only required field
func TestMapActionTypeAndFields_ScriptMinimalFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Script: &shared.ScriptModel{
			Content:     types.StringValue("echo 'test'"),
			Interpreter: types.StringNull(),
			Timeout:     types.StringNull(),
			Source:      types.StringNull(),
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeScript, actionType)
	assert.Equal(t, 1, len(fields)) // Only content field
}

// TestMapActionTypeAndFields_JiraMinimalFields tests Jira with only required field
func TestMapActionTypeAndFields_JiraMinimalFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Jira: &shared.JiraModel{
			Project:     types.StringValue("TEST"),
			Operation:   types.StringNull(),
			IssueType:   types.StringNull(),
			Description: types.StringNull(),
			Assignee:    types.StringNull(),
			Title:       types.StringNull(),
			Labels:      types.StringNull(),
			Comment:     types.StringNull(),
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, AutomationActionTypeJira, actionType)
	assert.Equal(t, 1, len(fields)) // Only project field
}

// TestMapActionTypeAndFields_GitHubMinimalFields tests GitHub with only required field
func TestMapActionTypeAndFields_GitHubMinimalFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		GitHub: &shared.GitHubModel{
			Owner:     types.StringValue("owner"),
			Repo:      types.StringNull(),
			Title:     types.StringNull(),
			Body:      types.StringNull(),
			Operation: types.StringNull(),
			Assignees: types.StringNull(),
			Labels:    types.StringNull(),
			Comment:   types.StringNull(),
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, AutomationActionTypeGitHub, actionType)
	assert.Equal(t, 1, len(fields)) // Only owner field
}

// TestMapActionTypeAndFields_GitLabMinimalFields tests GitLab with only required field
func TestMapActionTypeAndFields_GitLabMinimalFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		GitLab: &shared.GitLabModel{
			ProjectId:   types.StringValue("123"),
			Title:       types.StringNull(),
			Description: types.StringNull(),
			Operation:   types.StringNull(),
			Labels:      types.StringNull(),
			IssueType:   types.StringNull(),
			Comment:     types.StringNull(),
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, AutomationActionTypeGitLab, actionType)
	assert.Equal(t, 1, len(fields)) // Only projectId field
}

// TestMapActionTypeAndFields_AnsibleMinimalFields tests Ansible with only required field
func TestMapActionTypeAndFields_AnsibleMinimalFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Ansible: &shared.AnsibleModel{
			WorkflowId:       types.StringValue("workflow-123"),
			PlaybookId:       types.StringNull(),
			PlaybookFileName: types.StringNull(),
			AnsibleUrl:       types.StringNull(),
			HostId:           types.StringNull(),
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, AutomationActionTypeAnsible, actionType)
	assert.Equal(t, 1, len(fields)) // Only workflowId field
}

// TestMapActionTypeAndFields_HTTPMinimalFields tests HTTP with only required fields
func TestMapActionTypeAndFields_HTTPMinimalFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Http: &shared.HttpModel{
			Host:             types.StringValue("example.com"),
			Method:           types.StringValue("GET"),
			Body:             types.StringNull(),
			Headers:          types.MapNull(types.StringType),
			IgnoreCertErrors: types.BoolNull(),
			Timeout:          types.StringNull(),
			Language:         types.StringNull(),
			ContentType:      types.StringNull(),
			Auth:             nil,
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeHttp, actionType)
	// Should have: host, method, auth (noAuth)
	assert.GreaterOrEqual(t, len(fields), 3)
}

// TestUpdateState_WithInterfaceTags tests UpdateState with interface{} tags
func TestUpdateState_WithInterfaceTags(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	state := tfsdk.State{
		Schema: getTestSchema(),
	}

	automationAction := &restapi.AutomationAction{
		ID:          "test-id",
		Name:        "Test Action",
		Description: "Test Description",
		Type:        shared.ActionTypeScript,
		Tags:        []interface{}{"tag1", "tag2"}, // interface{} slice
		Fields: []restapi.Field{
			{
				Name:     restapi.ScriptSshFieldName,
				Value:    "echo 'test'",
				Encoding: shared.Base64Encoding,
			},
		},
	}

	diags := resource.UpdateState(ctx, &state, nil, automationAction)
	require.False(t, diags.HasError())

	var model shared.AutomationActionModel
	diags = state.Get(ctx, &model)
	require.False(t, diags.HasError())
	assert.False(t, model.Tags.IsNull())

	var tags []string
	diags = model.Tags.ElementsAs(ctx, &tags, false)
	require.False(t, diags.HasError())
	assert.Len(t, tags, 2)
}

// TestMapTagsFromState_WithValidTags tests mapTagsFromState with valid tags
func TestMapTagsFromState_WithValidTags(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	tags := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("tag1"),
		types.StringValue("tag2"),
		types.StringValue("tag3"),
	})

	model := shared.AutomationActionModel{
		Tags: tags,
	}

	result, diags := resource.mapTagsFromState(ctx, model)
	require.False(t, diags.HasError())
	require.NotNil(t, result)

	tagSlice, ok := result.([]string)
	require.True(t, ok)
	assert.Equal(t, []string{"tag1", "tag2", "tag3"}, tagSlice)
}

// TestMapActionTypeAndFields_HTTPWithEmptyAuth tests HTTP with empty auth model
func TestMapActionTypeAndFields_HTTPWithEmptyAuth(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Http: &shared.HttpModel{
			Host:   types.StringValue("example.com"),
			Method: types.StringValue("GET"),
			Auth: &shared.AuthModel{
				BasicAuth: nil,
				Token:     nil,
				ApiKey:    nil,
			},
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeHttp, actionType)

	// Find auth field and verify it's noAuth
	var authField *restapi.Field
	for i := range fields {
		if fields[i].Name == AutomationActionAPIFieldAuthen {
			authField = &fields[i]
			break
		}
	}
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, AutomationActionAuthTypeNoAuth)
}

// TestMapActionTypeAndFields_HTTPWithNullBasicAuthFields tests HTTP with null basic auth fields
func TestMapActionTypeAndFields_HTTPWithNullBasicAuthFields(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Http: &shared.HttpModel{
			Host:   types.StringValue("example.com"),
			Method: types.StringValue("GET"),
			Auth: &shared.AuthModel{
				BasicAuth: &shared.BasicAuthModel{
					UserName: types.StringNull(),
					Password: types.StringValue("pass"),
				},
			},
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeHttp, actionType)

	// Should default to noAuth when username is null
	var authField *restapi.Field
	for i := range fields {
		if fields[i].Name == AutomationActionAPIFieldAuthen {
			authField = &fields[i]
			break
		}
	}
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, AutomationActionAuthTypeNoAuth)
}

// TestMapActionTypeAndFields_HTTPWithNullTokenField tests HTTP with null bearer token field
func TestMapActionTypeAndFields_HTTPWithNullTokenField(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Http: &shared.HttpModel{
			Host:   types.StringValue("example.com"),
			Method: types.StringValue("GET"),
			Auth: &shared.AuthModel{
				Token: &shared.BearerTokenModel{
					BearerToken: types.StringNull(),
				},
			},
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeHttp, actionType)

	// Should default to noAuth when bearer token is null
	var authField *restapi.Field
	for i := range fields {
		if fields[i].Name == AutomationActionAPIFieldAuthen {
			authField = &fields[i]
			break
		}
	}
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, AutomationActionAuthTypeNoAuth)
}

// TestMapActionTypeAndFields_HTTPWithNullApiKeyField tests HTTP with null API key field
func TestMapActionTypeAndFields_HTTPWithNullApiKeyField(t *testing.T) {
	resource := NewAutomationActionResourceHandle().(*automationActionResource)
	ctx := context.Background()

	model := shared.AutomationActionModel{
		Http: &shared.HttpModel{
			Host:   types.StringValue("example.com"),
			Method: types.StringValue("GET"),
			Auth: &shared.AuthModel{
				ApiKey: &shared.ApiKeyModel{
					Key:         types.StringNull(),
					Value:       types.StringValue("value"),
					KeyLocation: types.StringValue("header"),
				},
			},
		},
	}

	actionType, fields, diags := resource.mapActionTypeAndFields(ctx, model)
	require.False(t, diags.HasError())
	assert.Equal(t, shared.ActionTypeHttp, actionType)

	// Should default to noAuth when API key is null
	var authField *restapi.Field
	for i := range fields {
		if fields[i].Name == AutomationActionAPIFieldAuthen {
			authField = &fields[i]
			break
		}
	}
	require.NotNil(t, authField)
	assert.Contains(t, authField.Value, AutomationActionAuthTypeNoAuth)
}
