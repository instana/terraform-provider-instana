package shared

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/restapi"
)

// Common models for automation actions

// AutomationActionModel represents the data model for an automation action
type AutomationActionModel struct {
	ID             types.String     `tfsdk:"id"`
	Name           types.String     `tfsdk:"name"`
	Description    types.String     `tfsdk:"description"`
	Tags           types.List       `tfsdk:"tags"`
	Script         *ScriptModel     `tfsdk:"script"`
	Http           *HttpModel       `tfsdk:"http"`
	Manual         *ManualModel     `tfsdk:"manual"`
	Jira           *JiraModel       `tfsdk:"jira"`
	GitHub         *GitHubModel     `tfsdk:"github"`
	DocLink        *DocLinkModel    `tfsdk:"doc_link"`
	GitLab         *GitLabModel     `tfsdk:"gitlab"`
	Ansible        *AnsibleModel    `tfsdk:"ansible"`
	InputParameter []ParameterModel `tfsdk:"input_parameter"`
}

type AnsibleModel struct {
	WorkflowId       types.String `tfsdk:"workflow_id"`
	PlaybookId       types.String `tfsdk:"playbook_id"`
	PlaybookFileName types.String `tfsdk:"playbook_file_name"`
	AnsibleUrl       types.String `tfsdk:"url"`
	HostId           types.String `tfsdk:"host_id"`
}

// ScriptModel represents the script configuration for an automation action
type ScriptModel struct {
	Content     types.String `tfsdk:"content"`
	Interpreter types.String `tfsdk:"interpreter"`
	Timeout     types.String `tfsdk:"timeout"`
	Source      types.String `tfsdk:"source"`
}

// HttpModel represents the HTTP configuration for an automation action
type HttpModel struct {
	Host             types.String `tfsdk:"host"`
	Method           types.String `tfsdk:"method"`
	Body             types.String `tfsdk:"body"`
	Headers          types.Map    `tfsdk:"headers"`
	IgnoreCertErrors types.Bool   `tfsdk:"ignore_certificate_errors"`
	Timeout          types.String `tfsdk:"timeout"`
	Language         types.String `tfsdk:"language"`
	ContentType      types.String `tfsdk:"content_type"`
	Auth             *AuthModel   `tfsdk:"auth"`
}

// AuthModel represents the authentication configuration for HTTP requests
type AuthModel struct {
	BasicAuth *BasicAuthModel   `tfsdk:"basic_auth"`
	Token     *BearerTokenModel `tfsdk:"token"`
	ApiKey    *ApiKeyModel      `tfsdk:"api_key"`
}

// BasicAuthModel represents the basic authentication configuration
type BasicAuthModel struct {
	UserName types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// ApiKeyModel represents the API key authentication configuration
type ApiKeyModel struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	KeyLocation types.String `tfsdk:"key_location"`
}

// BearerTokenModel represents the bearer token authentication configuration
type BearerTokenModel struct {
	BearerToken types.String `tfsdk:"bearer_token"`
}

type ManualModel struct {
	Content types.String `tfsdk:"content"`
}

type JiraModel struct {
	Project     types.String `tfsdk:"project"`
	Operation   types.String `tfsdk:"operation"`
	IssueType   types.String `tfsdk:"issue_type"`
	Description types.String `tfsdk:"description"`
	Assignee    types.String `tfsdk:"assignee"`
	Title       types.String `tfsdk:"title"`
	Labels      types.String `tfsdk:"labels"`
	Comment     types.String `tfsdk:"comment"`
}

type GitHubModel struct {
	Owner     types.String `tfsdk:"owner"`
	Repo      types.String `tfsdk:"repo"`
	Title     types.String `tfsdk:"title"`
	Body      types.String `tfsdk:"body"`
	Operation types.String `tfsdk:"operation"`
	Assignees types.String `tfsdk:"assignees"`
	Labels    types.String `tfsdk:"labels"`
	Comment   types.String `tfsdk:"comment"`
}

type DocLinkModel struct {
	Url types.String `tfsdk:"url"`
}

type GitLabModel struct {
	ProjectId   types.String `tfsdk:"project_id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Operation   types.String `tfsdk:"operation"`
	Labels      types.String `tfsdk:"labels"`
	IssueType   types.String `tfsdk:"issue_type"`
	Comment     types.String `tfsdk:"comment"`
}

// ParameterModel represents an input parameter for an automation action
type ParameterModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Label       types.String `tfsdk:"label"`
	Required    types.Bool   `tfsdk:"required"`
	Hidden      types.Bool   `tfsdk:"hidden"`
	Type        types.String `tfsdk:"type"`
	Value       types.String `tfsdk:"value"`
}

// action types
const ActionTypeScript = "SCRIPT"
const ActionTypeHttp = "HTTP"

// encodings
const AsciiEncoding = "ascii"
const Base64Encoding = "base64"
const UTF8Encoding = "UTF8"

// Common mapping functions

// MapTagsToState maps tags from API to Terraform state
func MapTagsToState(ctx context.Context, tags interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tags == nil {
		return types.ListNull(types.StringType), diags
	}

	// Handle tags based on their type
	switch v := tags.(type) {
	case []interface{}:
		elements := make([]attr.Value, len(v))
		for i, tag := range v {
			if strTag, ok := tag.(string); ok {
				elements[i] = types.StringValue(strTag)
			} else {
				diags.AddError(
					"Error mapping tags",
					fmt.Sprintf("Tag at index %d is not a string", i),
				)
				return types.ListNull(types.StringType), diags
			}
		}
		return types.ListValueMust(types.StringType, elements), diags
	default:
		diags.AddError(
			"Error mapping tags",
			"Tags are not in the expected format",
		)
		return types.ListNull(types.StringType), diags
	}
}

// MapTagsFromState maps tags from Terraform state to API
func MapTagsFromState(ctx context.Context, tagsList types.List) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if tagsList.IsNull() {
		return nil, diags
	}

	var tags []string
	diags.Append(tagsList.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	return tags, diags
}

// MapInputParametersToState maps input parameters from API to Terraform state
func MapInputParametersToState(ctx context.Context, parameters []restapi.Parameter) []ParameterModel {
	if len(parameters) == 0 {
		return nil
	}

	models := make([]ParameterModel, len(parameters))
	for i, param := range parameters {
		models[i] = ParameterModel{
			Name:        types.StringValue(param.Name),
			Description: types.StringValue(param.Description),
			Label:       types.StringValue(param.Label),
			Required:    types.BoolValue(param.Required),
			Hidden:      types.BoolValue(param.Hidden),
			Type:        types.StringValue(param.Type),
			Value:       types.StringValue(param.Value),
		}
	}

	return models
}

// MapInputParametersFromState maps input parameters from Terraform state to API
func MapInputParametersFromState(ctx context.Context, model AutomationActionModel) ([]restapi.Parameter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(model.InputParameter) == 0 {
		return nil, diags
	}

	parameters := make([]restapi.Parameter, len(model.InputParameter))
	for i, param := range model.InputParameter {
		parameters[i] = restapi.Parameter{
			Name:        param.Name.ValueString(),
			Description: param.Description.ValueString(),
			Label:       param.Label.ValueString(),
			Required:    param.Required.ValueBool(),
			Hidden:      param.Hidden.ValueBool(),
			Type:        param.Type.ValueString(),
			Value:       param.Value.ValueString(),
		}
	}

	return parameters, diags
}

// MapInputParametersToMap converts input parameters to a map for policy actions
func MapInputParametersToMap(ctx context.Context, inputParams []restapi.Parameter) types.Map {
	if len(inputParams) == 0 {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value)
	for _, param := range inputParams {
		elements[param.Name] = types.StringValue(param.Value)
	}

	return types.MapValueMust(types.StringType, elements)
}

// MapInputParametersFromMap converts a map to input parameter values for API
func MapInputParametersFromMap(ctx context.Context, inputParamsMap types.Map) ([]restapi.InputParameterValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var inputParams []restapi.InputParameterValue

	if inputParamsMap.IsNull() {
		return inputParams, diags
	}

	elements := make(map[string]string)
	diags.Append(inputParamsMap.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return inputParams, diags
	}

	inputParams = make([]restapi.InputParameterValue, 0, len(elements))
	for name, value := range elements {
		inputParams = append(inputParams, restapi.InputParameterValue{
			Name:  name,
			Value: value,
		})
	}

	return inputParams, diags
}

// GetAutomationActionSchemaAttributes returns the schema attributes for an automation action
// This is used to embed the full action schema in other resources like automation policies
func GetAutomationActionSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Optional:    true,
			Description: "The ID of the automation action.",
		},
		AutomationActionFieldName: schema.StringAttribute{
			Optional:    true,
			Description: "The name of the automation action.",
		},
		AutomationActionFieldDescription: schema.StringAttribute{
			Optional:    true,
			Description: "The description of the automation action.",
		},
		AutomationActionFieldTags: schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "The tags of the automation action.",
		},
		AutomationActionFieldScript: schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Script configuration for the automation action.",
			Attributes: map[string]schema.Attribute{
				AutomationActionFieldContent: schema.StringAttribute{
					Required:    true,
					Description: "The script content.",
				},
				AutomationActionFieldInterpreter: schema.StringAttribute{
					Optional:    true,
					Description: "The script interpreter.",
				},
				AutomationActionFieldTimeout: schema.StringAttribute{
					Optional:    true,
					Description: "The timeout for script execution in seconds.",
				},
				"source": schema.StringAttribute{
					Optional:    true,
					Description: "The source of the script.",
				},
			},
		},
		AutomationActionFieldHttp: schema.SingleNestedAttribute{
			Optional:    true,
			Description: "HTTP configuration for the automation action.",
			Attributes: map[string]schema.Attribute{
				AutomationActionFieldHost: schema.StringAttribute{
					Required:    true,
					Description: "The URL of the HTTP request.",
				},
				AutomationActionFieldMethod: schema.StringAttribute{
					Required:    true,
					Description: "The HTTP method.",
					Validators: []validator.String{
						stringvalidator.OneOf("GET", "POST", "PUT", "DELETE"),
					},
				},
				AutomationActionFieldBody: schema.StringAttribute{
					Optional:    true,
					Description: "The body of the HTTP request.",
				},
				AutomationActionFieldHeaders: schema.MapAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: "The headers of the HTTP request.",
				},
				AutomationActionFieldIgnoreCertErrors: schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to ignore certificate errors for the request.",
				},
				AutomationActionFieldTimeout: schema.StringAttribute{
					Optional:    true,
					Description: "The timeout for HTTP request execution in seconds.",
				},
				"language": schema.StringAttribute{
					Optional:    true,
					Description: "The language for the HTTP request.",
				},
				"content_type": schema.StringAttribute{
					Optional:    true,
					Description: "The content type for the HTTP request.",
				},
				"auth": schema.SingleNestedAttribute{
					Optional:    true,
					Description: "Authentication configuration for the HTTP request.",
					Attributes: map[string]schema.Attribute{
						"basic_auth": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Basic authentication configuration.",
							Attributes: map[string]schema.Attribute{
								"username": schema.StringAttribute{
									Required:    true,
									Description: "The username for basic authentication.",
								},
								"password": schema.StringAttribute{
									Required:    true,
									Description: "The password for basic authentication.",
								},
							},
						},
						"token": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Bearer token authentication configuration.",
							Attributes: map[string]schema.Attribute{
								"bearer_token": schema.StringAttribute{
									Required:    true,
									Description: "The bearer token for authentication.",
								},
							},
						},
						"api_key": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "API key authentication configuration.",
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Required:    true,
									Description: "The API key name.",
								},
								"value": schema.StringAttribute{
									Required:    true,
									Description: "The API key value.",
								},
								"key_location": schema.StringAttribute{
									Required:    true,
									Description: "Where to add the API key (header or query).",
								},
							},
						},
					},
				},
			},
		},
		"manual": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Manual action configuration.",
			Attributes: map[string]schema.Attribute{
				"content": schema.StringAttribute{
					Required:    true,
					Description: "The content for the manual action.",
				},
			},
		},
		"jira": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Jira action configuration.",
			Attributes: map[string]schema.Attribute{
				"project": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira project.",
				},
				"operation": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira operation type.",
				},
				"issue_type": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira issue type.",
				},
				"description": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira issue description.",
				},
				"assignee": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira issue assignee.",
				},
				"title": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira issue title.",
				},
				"labels": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira issue labels.",
				},
				"comment": schema.StringAttribute{
					Optional:    true,
					Description: "The Jira issue comment.",
				},
			},
		},
		"github": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "GitHub action configuration.",
			Attributes: map[string]schema.Attribute{
				"owner": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub repository owner.",
				},
				"repo": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub repository name.",
				},
				"title": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub issue title.",
				},
				"body": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub issue body.",
				},
				"operation": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub operation type.",
				},
				"assignees": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub issue assignees.",
				},
				"labels": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub issue labels.",
				},
				"comment": schema.StringAttribute{
					Optional:    true,
					Description: "The GitHub issue comment.",
				},
			},
		},
		"doc_link": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Documentation link action configuration.",
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					Required:    true,
					Description: "The URL to the documentation.",
				},
			},
		},
		"gitlab": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "GitLab action configuration.",
			Attributes: map[string]schema.Attribute{
				"project_id": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab project ID.",
				},
				"title": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab issue title.",
				},
				"description": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab issue description.",
				},
				"operation": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab operation type.",
				},
				"labels": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab issue labels.",
				},
				"issue_type": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab issue type.",
				},
				"comment": schema.StringAttribute{
					Optional:    true,
					Description: "The GitLab issue comment.",
				},
			},
		},
		"ansible": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Ansible action configuration.",
			Attributes: map[string]schema.Attribute{
				"workflow_id": schema.StringAttribute{
					Optional:    true,
					Description: "The Ansible workflow ID.",
				},
				"playbook_id": schema.StringAttribute{
					Optional:    true,
					Description: "The Ansible playbook ID.",
				},
				"playbook_file_name": schema.StringAttribute{
					Optional:    true,
					Description: "The Ansible playbook file name.",
				},
				"url": schema.StringAttribute{
					Optional:    true,
					Description: "The Ansible URL.",
				},
				"host_id": schema.StringAttribute{
					Optional:    true,
					Description: "The host ID from which this action is created.",
				},
			},
		},
		AutomationActionFieldInputParameter: schema.ListNestedAttribute{
			Optional:    true,
			Description: "Input parameters for the automation action.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					AutomationActionParameterFieldName: schema.StringAttribute{
						Required:    true,
						Description: "The name of the parameter.",
					},
					AutomationActionParameterFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: "The description of the parameter.",
					},
					AutomationActionParameterFieldLabel: schema.StringAttribute{
						Required:    true,
						Description: "The label of the parameter.",
					},
					AutomationActionParameterFieldRequired: schema.BoolAttribute{
						Required:    true,
						Description: "Whether the parameter is required.",
					},
					AutomationActionParameterFieldHidden: schema.BoolAttribute{
						Required:    true,
						Description: "Whether the parameter is hidden.",
					},
					AutomationActionParameterFieldType: schema.StringAttribute{
						Required:    true,
						Description: "The type of the parameter.",
						Validators: []validator.String{
							stringvalidator.OneOf("static", "dynamic", "vault"),
						},
					},
					AutomationActionParameterFieldValue: schema.StringAttribute{
						Required:    true,
						Description: "The value of the parameter.",
					},
				},
			},
		},
	}
}

// Type-specific field mapping functions

// MapActionTypeFieldsToState maps action type-specific fields from API to state
func MapActionTypeFieldsToState(ctx context.Context, action *restapi.AutomationAction, model *AutomationActionModel) {
	switch action.Type {
	case ActionTypeScript:
		scriptConfig := MapScriptFieldsToState(action)
		model.Script = &scriptConfig
	case ActionTypeHttp:
		httpConfig := MapHttpFieldsToState(action)
		model.Http = &httpConfig
	case "MANUAL":
		manualConfig := MapManualFieldsToState(action)
		model.Manual = &manualConfig
	case "JIRA":
		jiraConfig := MapJiraFieldsToState(action)
		model.Jira = &jiraConfig
	case "GITHUB":
		githubConfig := MapGitHubFieldsToState(action)
		model.GitHub = &githubConfig
	case "DOC_LINK":
		docLinkConfig := MapDocLinkFieldsToState(action)
		model.DocLink = &docLinkConfig
	case "GITLAB":
		gitlabConfig := MapGitLabFieldsToState(action)
		model.GitLab = &gitlabConfig
	case "ANSIBLE":
		ansibleConfig := MapAnsibleFieldsToState(action)
		model.Ansible = &ansibleConfig
	}
}

// GetFieldValue retrieves a field value from an automation action
func GetFieldValue(action *restapi.AutomationAction, fieldName string) string {
	for _, v := range action.Fields {
		if v.Name == fieldName {
			return v.Value
		}
	}
	return ""
}

// GetBoolFieldValueOrDefault retrieves a boolean field value or returns default
func GetBoolFieldValueOrDefault(action *restapi.AutomationAction, fieldName string, defaultValue bool) bool {
	valueStr := GetFieldValue(action, fieldName)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr == "true"
}

// MapScriptFieldsToState maps script-specific fields to state
func MapScriptFieldsToState(action *restapi.AutomationAction) ScriptModel {
	scriptModel := ScriptModel{
		Content: types.StringValue(GetFieldValue(action, restapi.ScriptSshFieldName)),
	}

	interpreter := GetFieldValue(action, restapi.SubtypeFieldName)
	if interpreter != "" {
		scriptModel.Interpreter = types.StringValue(interpreter)
	} else {
		scriptModel.Interpreter = types.StringNull()
	}

	timeout := GetFieldValue(action, restapi.TimeoutFieldName)
	if timeout != "" {
		scriptModel.Timeout = types.StringValue(timeout)
	} else {
		scriptModel.Timeout = types.StringNull()
	}

	source := GetFieldValue(action, "source")
	if source != "" {
		scriptModel.Source = types.StringValue(source)
	} else {
		scriptModel.Source = types.StringNull()
	}

	return scriptModel
}

// MapHttpFieldsToState maps HTTP-specific fields to state
func MapHttpFieldsToState(action *restapi.AutomationAction) HttpModel {
	httpModel := HttpModel{
		Host:   types.StringValue(GetFieldValue(action, restapi.HttpHostFieldName)),
		Method: types.StringValue(GetFieldValue(action, restapi.HttpMethodFieldName)),
	}

	body := GetFieldValue(action, restapi.HttpBodyFieldName)
	if body != "" {
		httpModel.Body = types.StringValue(body)
	} else {
		httpModel.Body = types.StringNull()
	}

	// Only set IgnoreCertErrors if the field is present in the API response
	ignoreCertErrorsValue := GetFieldValue(action, restapi.HttpIgnoreCertErrorsFieldName)
	if ignoreCertErrorsValue != "" {
		httpModel.IgnoreCertErrors = types.BoolValue(ignoreCertErrorsValue == "true")
	} else {
		httpModel.IgnoreCertErrors = types.BoolNull()
	}

	timeout := GetFieldValue(action, restapi.TimeoutFieldName)
	if timeout != "" {
		httpModel.Timeout = types.StringValue(timeout)
	} else {
		httpModel.Timeout = types.StringNull()
	}

	language := GetFieldValue(action, "language")
	if language != "" {
		httpModel.Language = types.StringValue(language)
	} else {
		httpModel.Language = types.StringNull()
	}

	contentType := GetFieldValue(action, "content_type")
	if contentType != "" {
		httpModel.ContentType = types.StringValue(contentType)
	} else {
		httpModel.ContentType = types.StringNull()
	}

	// Handle auth
	httpModel.Auth = MapAuthFieldsToState(action)

	// Handle headers
	httpModel.Headers = MapHeadersToState(action)

	return httpModel
}

// MapAuthFieldsToState maps authentication fields to state
func MapAuthFieldsToState(action *restapi.AutomationAction) *AuthModel {
	authData := GetFieldValue(action, "authen")
	if authData == "" {
		return nil
	}

	var authMap map[string]interface{}
	if err := json.Unmarshal([]byte(authData), &authMap); err != nil {
		return nil
	}

	authType, _ := authMap["type"].(string)
	switch authType {
	case "basicAuth":
		username, _ := authMap["username"].(string)
		password, _ := authMap["password"].(string)
		return &AuthModel{
			BasicAuth: &BasicAuthModel{
				UserName: types.StringValue(username),
				Password: types.StringValue(password),
			},
		}
	case "bearerToken":
		bearerToken, _ := authMap["bearerToken"].(string)
		return &AuthModel{
			Token: &BearerTokenModel{
				BearerToken: types.StringValue(bearerToken),
			},
		}
	case "apiKey":
		key, _ := authMap["apiKey"].(string)
		value, _ := authMap["apiKeyValue"].(string)
		location, _ := authMap["apiKeyAddTo"].(string)
		return &AuthModel{
			ApiKey: &ApiKeyModel{
				Key:         types.StringValue(key),
				Value:       types.StringValue(value),
				KeyLocation: types.StringValue(location),
			},
		}
	}
	return nil
}

// MapHeadersToState maps HTTP headers to state
func MapHeadersToState(action *restapi.AutomationAction) types.Map {
	headersData := GetFieldValue(action, restapi.HttpHeaderFieldName)
	if headersData == "" || headersData == "{}" {
		return types.MapNull(types.StringType)
	}

	var headersMap map[string]interface{}
	if err := json.Unmarshal([]byte(headersData), &headersMap); err != nil {
		return types.MapNull(types.StringType)
	}

	if len(headersMap) == 0 {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value, len(headersMap))
	for k, v := range headersMap {
		if strVal, ok := v.(string); ok {
			elements[k] = types.StringValue(strVal)
		} else {
			elements[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
	}

	return types.MapValueMust(types.StringType, elements)
}

// MapManualFieldsToState maps manual action fields to state
func MapManualFieldsToState(action *restapi.AutomationAction) ManualModel {
	return ManualModel{
		Content: types.StringValue(GetFieldValue(action, "content")),
	}
}

// MapJiraFieldsToState maps Jira-specific fields to state
func MapJiraFieldsToState(action *restapi.AutomationAction) JiraModel {
	jiraModel := JiraModel{}

	if project := GetFieldValue(action, "project"); project != "" {
		jiraModel.Project = types.StringValue(project)
	} else {
		jiraModel.Project = types.StringNull()
	}

	if operation := GetFieldValue(action, "ticketActionType"); operation != "" {
		jiraModel.Operation = types.StringValue(operation)
	} else {
		jiraModel.Operation = types.StringNull()
	}

	if issueType := GetFieldValue(action, "issue_type"); issueType != "" {
		jiraModel.IssueType = types.StringValue(issueType)
	} else {
		jiraModel.IssueType = types.StringNull()
	}

	if description := GetFieldValue(action, "body"); description != "" {
		jiraModel.Description = types.StringValue(description)
	} else {
		jiraModel.Description = types.StringNull()
	}

	if assignee := GetFieldValue(action, "assignee"); assignee != "" {
		jiraModel.Assignee = types.StringValue(assignee)
	} else {
		jiraModel.Assignee = types.StringNull()
	}

	if title := GetFieldValue(action, "summary"); title != "" {
		jiraModel.Title = types.StringValue(title)
	} else {
		jiraModel.Title = types.StringNull()
	}

	if labels := GetFieldValue(action, "labels"); labels != "" {
		jiraModel.Labels = types.StringValue(labels)
	} else {
		jiraModel.Labels = types.StringNull()
	}

	if comment := GetFieldValue(action, "comment"); comment != "" {
		jiraModel.Comment = types.StringValue(comment)
	} else {
		jiraModel.Comment = types.StringNull()
	}

	return jiraModel
}

// MapGitHubFieldsToState maps GitHub-specific fields to state
func MapGitHubFieldsToState(action *restapi.AutomationAction) GitHubModel {
	githubModel := GitHubModel{}

	if owner := GetFieldValue(action, "owner"); owner != "" {
		githubModel.Owner = types.StringValue(owner)
	} else {
		githubModel.Owner = types.StringNull()
	}

	if repo := GetFieldValue(action, "repo"); repo != "" {
		githubModel.Repo = types.StringValue(repo)
	} else {
		githubModel.Repo = types.StringNull()
	}

	if title := GetFieldValue(action, "title"); title != "" {
		githubModel.Title = types.StringValue(title)
	} else {
		githubModel.Title = types.StringNull()
	}

	if body := GetFieldValue(action, "body"); body != "" {
		githubModel.Body = types.StringValue(body)
	} else {
		githubModel.Body = types.StringNull()
	}

	if operation := GetFieldValue(action, "ticketType"); operation != "" {
		githubModel.Operation = types.StringValue(operation)
	} else {
		githubModel.Operation = types.StringNull()
	}

	if assignees := GetFieldValue(action, "assignees"); assignees != "" {
		githubModel.Assignees = types.StringValue(assignees)
	} else {
		githubModel.Assignees = types.StringNull()
	}

	if labels := GetFieldValue(action, "labels"); labels != "" {
		githubModel.Labels = types.StringValue(labels)
	} else {
		githubModel.Labels = types.StringNull()
	}

	if comment := GetFieldValue(action, "comment"); comment != "" {
		githubModel.Comment = types.StringValue(comment)
	} else {
		githubModel.Comment = types.StringNull()
	}

	return githubModel
}

// MapDocLinkFieldsToState maps documentation link fields to state
func MapDocLinkFieldsToState(action *restapi.AutomationAction) DocLinkModel {
	return DocLinkModel{
		Url: types.StringValue(GetFieldValue(action, "url")),
	}
}

// MapGitLabFieldsToState maps GitLab-specific fields to state
func MapGitLabFieldsToState(action *restapi.AutomationAction) GitLabModel {
	gitlabModel := GitLabModel{}

	if projectId := GetFieldValue(action, "projectId"); projectId != "" {
		gitlabModel.ProjectId = types.StringValue(projectId)
	} else {
		gitlabModel.ProjectId = types.StringNull()
	}

	if title := GetFieldValue(action, "title"); title != "" {
		gitlabModel.Title = types.StringValue(title)
	} else {
		gitlabModel.Title = types.StringNull()
	}

	if description := GetFieldValue(action, "body"); description != "" {
		gitlabModel.Description = types.StringValue(description)
	} else {
		gitlabModel.Description = types.StringNull()
	}

	if operation := GetFieldValue(action, "ticketActionType"); operation != "" {
		gitlabModel.Operation = types.StringValue(operation)
	} else {
		gitlabModel.Operation = types.StringNull()
	}

	if labels := GetFieldValue(action, "labels"); labels != "" {
		gitlabModel.Labels = types.StringValue(labels)
	} else {
		gitlabModel.Labels = types.StringNull()
	}

	if issueType := GetFieldValue(action, "issue_type"); issueType != "" {
		gitlabModel.IssueType = types.StringValue(issueType)
	} else {
		gitlabModel.IssueType = types.StringNull()
	}

	if comment := GetFieldValue(action, "comment"); comment != "" {
		gitlabModel.Comment = types.StringValue(comment)
	} else {
		gitlabModel.Comment = types.StringNull()
	}

	return gitlabModel
}

// MapAnsibleFieldsToState maps Ansible-specific fields to state
func MapAnsibleFieldsToState(action *restapi.AutomationAction) AnsibleModel {
	ansibleModel := AnsibleModel{}

	if workflowId := GetFieldValue(action, "workflowId"); workflowId != "" {
		ansibleModel.WorkflowId = types.StringValue(workflowId)
	} else {
		ansibleModel.WorkflowId = types.StringNull()
	}

	if ansibleUrl := GetFieldValue(action, "ansibleUrl"); ansibleUrl != "" {
		ansibleModel.AnsibleUrl = types.StringValue(ansibleUrl)
	} else {
		ansibleModel.AnsibleUrl = types.StringNull()
	}

	if hostId := GetFieldValue(action, "hostId"); hostId != "" {
		ansibleModel.HostId = types.StringValue(hostId)
	} else {
		ansibleModel.HostId = types.StringNull()
	}

	if playbookId := GetFieldValue(action, "playbookId"); playbookId != "" {
		ansibleModel.PlaybookId = types.StringValue(playbookId)
	} else {
		ansibleModel.PlaybookId = types.StringNull()
	}

	if playbookFileName := GetFieldValue(action, "playbookFileName"); playbookFileName != "" {
		ansibleModel.PlaybookFileName = types.StringValue(playbookFileName)
	} else {
		ansibleModel.PlaybookFileName = types.StringNull()
	}

	return ansibleModel
}

// MapActionTypeAndFieldsFromState determines action type and maps type-specific fields from state to API
func MapActionTypeAndFieldsFromState(ctx context.Context, model AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	var actionType string
	var fields []restapi.Field

	// Check which action type is configured
	if model.Script != nil && !model.Script.Content.IsNull() {
		actionType = ActionTypeScript
		fields = mapScriptFieldsFromState(*model.Script)
	} else if model.Http != nil && !model.Http.Host.IsNull() {
		actionType = ActionTypeHttp
		var d diag.Diagnostics
		fields, d = mapHttpFieldsFromState(ctx, *model.Http)
		diags.Append(d...)
	} else if model.Manual != nil && !model.Manual.Content.IsNull() {
		actionType = "MANUAL"
		fields = mapManualFieldsFromState(*model.Manual)
	} else if model.Jira != nil && !model.Jira.Project.IsNull() {
		actionType = "JIRA"
		fields = mapJiraFieldsFromState(*model.Jira)
	} else if model.GitHub != nil && !model.GitHub.Owner.IsNull() {
		actionType = "GITHUB"
		fields = mapGitHubFieldsFromState(*model.GitHub)
	} else if model.DocLink != nil && !model.DocLink.Url.IsNull() {
		actionType = "DOC_LINK"
		fields = mapDocLinkFieldsFromState(*model.DocLink)
	} else if model.GitLab != nil && !model.GitLab.ProjectId.IsNull() {
		actionType = "GITLAB"
		fields = mapGitLabFieldsFromState(*model.GitLab)
	} else if model.Ansible != nil && !model.Ansible.WorkflowId.IsNull() {
		actionType = "ANSIBLE"
		fields = mapAnsibleFieldsFromState(*model.Ansible)
	} else {
		diags.AddError(
			"Invalid action configuration",
			"One of script, http, manual, jira, github, doc_link, gitlab, or ansible configuration must be provided",
		)
		return "", nil, diags
	}

	return actionType, fields, diags
}

// Helper functions to map fields from state to API

func mapScriptFieldsFromState(scriptModel ScriptModel) []restapi.Field {
	fields := make([]restapi.Field, 0)

	// Content is required
	fields = append(fields, restapi.Field{
		Name:        restapi.ScriptSshFieldName,
		Description: restapi.ScriptSshFieldDescription,
		Value:       scriptModel.Content.ValueString(),
		Encoding:    Base64Encoding,
		Secured:     false,
	})

	// Interpreter is optional
	if !scriptModel.Interpreter.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.SubtypeFieldName,
			Description: restapi.SubtypeFieldDescription,
			Value:       scriptModel.Interpreter.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// Timeout is optional
	if !scriptModel.Timeout.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.TimeoutFieldName,
			Description: restapi.TimeoutFieldDescription,
			Value:       scriptModel.Timeout.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// Source is optional
	if !scriptModel.Source.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        "source",
			Description: "The source of the script",
			Value:       scriptModel.Source.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	return fields
}

func mapHttpFieldsFromState(ctx context.Context, httpModel HttpModel) ([]restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := make([]restapi.Field, 0)

	// Host and method are required
	fields = append(fields, restapi.Field{
		Name:        restapi.HttpHostFieldName,
		Description: restapi.HttpHostFieldDescription,
		Value:       httpModel.Host.ValueString(),
		Encoding:    AsciiEncoding,
		Secured:     false,
	})

	fields = append(fields, restapi.Field{
		Name:        restapi.HttpMethodFieldName,
		Description: restapi.HttpMethodFieldDescription,
		Value:       httpModel.Method.ValueString(),
		Encoding:    AsciiEncoding,
		Secured:     false,
	})

	// Body is optional
	if !httpModel.Body.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.HttpBodyFieldName,
			Description: restapi.HttpBodyFieldDescription,
			Value:       httpModel.Body.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// IgnoreCertErrors is optional
	if !httpModel.IgnoreCertErrors.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.HttpIgnoreCertErrorsFieldName,
			Description: restapi.HttpIgnoreCertErrorsFieldDescription,
			Value:       fmt.Sprintf("%t", httpModel.IgnoreCertErrors.ValueBool()),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// Timeout is optional
	if !httpModel.Timeout.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.TimeoutFieldName,
			Description: restapi.TimeoutFieldDescription,
			Value:       httpModel.Timeout.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// Language is optional
	if !httpModel.Language.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        "language",
			Description: "The language for the HTTP request",
			Value:       httpModel.Language.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// Content Type is optional
	if !httpModel.ContentType.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        "content_type",
			Description: "The content type for the HTTP request",
			Value:       httpModel.ContentType.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	// Auth is optional - serialize as JSON
	authValue := mapAuthFromState(httpModel.Auth)
	authJson, err := json.Marshal(authValue)
	if err != nil {
		diags.AddError(
			"Error marshaling auth",
			fmt.Sprintf("Failed to marshal auth: %s", err),
		)
		return nil, diags
	}

	fields = append(fields, restapi.Field{
		Name:        "authen",
		Description: "Authentication for the HTTPS request",
		Value:       string(authJson),
		Encoding:    AsciiEncoding,
		Secured:     false,
	})

	// Headers are optional
	if !httpModel.Headers.IsNull() {
		headersMap := make(map[string]string)
		diags.Append(httpModel.Headers.ElementsAs(ctx, &headersMap, false)...)
		if diags.HasError() {
			return nil, diags
		}

		headersJson, err := json.Marshal(headersMap)
		if err != nil {
			diags.AddError(
				"Error marshaling HTTP headers",
				fmt.Sprintf("Failed to marshal HTTP headers: %s", err),
			)
			return nil, diags
		}

		fields = append(fields, restapi.Field{
			Name:        restapi.HttpHeaderFieldName,
			Description: restapi.HttpHeaderFieldDescription,
			Value:       string(headersJson),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	}

	return fields, diags
}

func mapAuthFromState(auth *AuthModel) map[string]string {
	if auth == nil {
		return map[string]string{"type": "noAuth"}
	}

	if auth.BasicAuth != nil && !auth.BasicAuth.UserName.IsNull() {
		return map[string]string{
			"type":     "basicAuth",
			"username": auth.BasicAuth.UserName.ValueString(),
			"password": auth.BasicAuth.Password.ValueString(),
		}
	}

	if auth.Token != nil && !auth.Token.BearerToken.IsNull() {
		return map[string]string{
			"type":        "bearerToken",
			"bearerToken": auth.Token.BearerToken.ValueString(),
		}
	}

	if auth.ApiKey != nil && !auth.ApiKey.Key.IsNull() {
		return map[string]string{
			"type":        "apiKey",
			"apiKey":      auth.ApiKey.Key.ValueString(),
			"apiKeyValue": auth.ApiKey.Value.ValueString(),
			"apiKeyAddTo": auth.ApiKey.KeyLocation.ValueString(),
		}
	}

	return map[string]string{"type": "noAuth"}
}

func mapManualFieldsFromState(manualModel ManualModel) []restapi.Field {
	return []restapi.Field{
		{
			Name:        "content",
			Description: "Content for manual action",
			Value:       manualModel.Content.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		},
	}
}

func mapJiraFieldsFromState(jiraModel JiraModel) []restapi.Field {
	fields := make([]restapi.Field, 0)

	if !jiraModel.Project.IsNull() {
		fields = append(fields, restapi.Field{Name: "project", Description: "jira project", Value: jiraModel.Project.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.Operation.IsNull() {
		fields = append(fields, restapi.Field{Name: "ticketActionType", Description: "jira ticket type", Value: jiraModel.Operation.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.IssueType.IsNull() {
		fields = append(fields, restapi.Field{Name: "issue_type", Description: "jira issue type", Value: jiraModel.IssueType.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.Description.IsNull() {
		fields = append(fields, restapi.Field{Name: "body", Description: "jira issue description", Value: jiraModel.Description.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.Assignee.IsNull() {
		fields = append(fields, restapi.Field{Name: "assignee", Description: "jira issue assignee", Value: jiraModel.Assignee.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.Title.IsNull() {
		fields = append(fields, restapi.Field{Name: "summary", Description: "jira issue summary", Value: jiraModel.Title.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.Labels.IsNull() {
		fields = append(fields, restapi.Field{Name: "labels", Description: "jira issue labels", Value: jiraModel.Labels.ValueString(), Encoding: AsciiEncoding})
	}
	if !jiraModel.Comment.IsNull() {
		fields = append(fields, restapi.Field{Name: "comment", Description: "jira issue comment", Value: jiraModel.Comment.ValueString(), Encoding: AsciiEncoding})
	}

	return fields
}

func mapGitHubFieldsFromState(githubModel GitHubModel) []restapi.Field {
	fields := make([]restapi.Field, 0)

	if !githubModel.Owner.IsNull() {
		fields = append(fields, restapi.Field{Name: "owner", Description: "github issue owner/repo", Value: githubModel.Owner.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Repo.IsNull() {
		fields = append(fields, restapi.Field{Name: "repo", Description: "github issue repo", Value: githubModel.Repo.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Title.IsNull() {
		fields = append(fields, restapi.Field{Name: "title", Description: "github issue title", Value: githubModel.Title.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Body.IsNull() {
		fields = append(fields, restapi.Field{Name: "body", Description: "github issue body", Value: githubModel.Body.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Operation.IsNull() {
		fields = append(fields, restapi.Field{Name: "ticketType", Description: "github issue type", Value: githubModel.Operation.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Assignees.IsNull() {
		fields = append(fields, restapi.Field{Name: "assignees", Description: "github issue assignees", Value: githubModel.Assignees.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Labels.IsNull() {
		fields = append(fields, restapi.Field{Name: "labels", Description: "github issue labels", Value: githubModel.Labels.ValueString(), Encoding: AsciiEncoding})
	}
	if !githubModel.Comment.IsNull() {
		fields = append(fields, restapi.Field{Name: "comment", Description: "github issue comment", Value: githubModel.Comment.ValueString(), Encoding: AsciiEncoding})
	}

	return fields
}

func mapDocLinkFieldsFromState(docLinkModel DocLinkModel) []restapi.Field {
	return []restapi.Field{
		{
			Name:        "url",
			Description: "URL to remediation documentation",
			Value:       docLinkModel.Url.ValueString(),
			Encoding:    UTF8Encoding,
		},
	}
}

func mapGitLabFieldsFromState(gitlabModel GitLabModel) []restapi.Field {
	fields := make([]restapi.Field, 0)

	if !gitlabModel.ProjectId.IsNull() {
		fields = append(fields, restapi.Field{Name: "projectId", Description: "gitlab projectId", Value: gitlabModel.ProjectId.ValueString(), Encoding: AsciiEncoding})
	}
	if !gitlabModel.Title.IsNull() {
		fields = append(fields, restapi.Field{Name: "title", Description: "gitlab issue title", Value: gitlabModel.Title.ValueString(), Encoding: AsciiEncoding})
	}
	if !gitlabModel.Description.IsNull() {
		fields = append(fields, restapi.Field{Name: "body", Description: "gitlab issue description", Value: gitlabModel.Description.ValueString(), Encoding: AsciiEncoding})
	}
	if !gitlabModel.Operation.IsNull() {
		fields = append(fields, restapi.Field{Name: "ticketActionType", Description: "gitlab ticket type", Value: gitlabModel.Operation.ValueString(), Encoding: AsciiEncoding})
	}
	if !gitlabModel.Labels.IsNull() {
		fields = append(fields, restapi.Field{Name: "labels", Description: "gitlab issue labels", Value: gitlabModel.Labels.ValueString(), Encoding: AsciiEncoding})
	}
	if !gitlabModel.IssueType.IsNull() {
		fields = append(fields, restapi.Field{Name: "issue_type", Description: "gitlab issue type", Value: gitlabModel.IssueType.ValueString(), Encoding: AsciiEncoding})
	}
	if !gitlabModel.Comment.IsNull() {
		fields = append(fields, restapi.Field{Name: "comment", Description: "gitlab issue comment", Value: gitlabModel.Comment.ValueString(), Encoding: AsciiEncoding})
	}

	return fields
}

func mapAnsibleFieldsFromState(ansibleModel AnsibleModel) []restapi.Field {
	fields := make([]restapi.Field, 0)

	if !ansibleModel.WorkflowId.IsNull() {
		fields = append(fields, restapi.Field{Name: "workflowId", Description: "The workflow ID", Value: ansibleModel.WorkflowId.ValueString(), Encoding: AsciiEncoding})
	}
	if !ansibleModel.AnsibleUrl.IsNull() {
		fields = append(fields, restapi.Field{Name: "ansibleUrl", Description: "The ansible url", Value: ansibleModel.AnsibleUrl.ValueString(), Encoding: AsciiEncoding})
	}
	if !ansibleModel.HostId.IsNull() {
		fields = append(fields, restapi.Field{Name: "hostId", Description: "The host ID from which this action is created", Value: ansibleModel.HostId.ValueString(), Encoding: AsciiEncoding})
	}
	if !ansibleModel.PlaybookId.IsNull() {
		fields = append(fields, restapi.Field{Name: "playbookId", Description: "The playbook ID", Value: ansibleModel.PlaybookId.ValueString(), Encoding: AsciiEncoding})
	}
	if !ansibleModel.PlaybookFileName.IsNull() {
		fields = append(fields, restapi.Field{Name: "playbookFileName", Description: "The playbook filename", Value: ansibleModel.PlaybookFileName.ValueString(), Encoding: AsciiEncoding})
	}

	return fields
}
