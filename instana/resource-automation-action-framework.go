package instana

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaAutomationActionFramework the name of the terraform-provider-instana resource to manage automation actions
const ResourceInstanaAutomationActionFramework = "automation_action"

// AutomationActionModel is now defined in resource-automation-action-mapping.go

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

// NewAutomationActionResourceHandleFramework creates the resource handle for Automation Actions
func NewAutomationActionResourceHandleFramework() ResourceHandleFramework[*restapi.AutomationAction] {
	return &automationActionResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaAutomationActionFramework,
			Schema: schema.Schema{
				Description: "This resource manages automation actions in Instana.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the automation action.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AutomationActionFieldName: schema.StringAttribute{
						Required:    true,
						Description: "The name of the automation action.",
					},
					AutomationActionFieldDescription: schema.StringAttribute{
						Required:    true,
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
				},
			},
			SchemaVersion: 0,
		},
	}
}

type automationActionResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *automationActionResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *automationActionResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AutomationAction] {
	return api.AutomationActions()
}

func (r *automationActionResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *automationActionResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, automationAction *restapi.AutomationAction) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API response
	model := AutomationActionModel{
		ID:          types.StringValue(automationAction.ID),
		Name:        types.StringValue(automationAction.Name),
		Description: types.StringValue(automationAction.Description),
	}

	// Handle tags
	if automationAction.Tags != nil {
		tagsList, d := r.mapTagsToState(ctx, automationAction.Tags)
		diags.Append(d...)
		if !diags.HasError() {
			model.Tags = tagsList
		}
	} else {
		model.Tags = types.ListNull(types.StringType)
	}

	// Handle input parameters
	if len(automationAction.InputParameters) > 0 {
		inputParams, d := r.mapInputParametersToState(ctx, automationAction.InputParameters)
		diags.Append(d...)
		if !diags.HasError() {
			model.InputParameter = inputParams
		}
	} else {
		model.InputParameter = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionParameterFieldName:        types.StringType,
				AutomationActionParameterFieldDescription: types.StringType,
				AutomationActionParameterFieldLabel:       types.StringType,
				AutomationActionParameterFieldRequired:    types.BoolType,
				AutomationActionParameterFieldHidden:      types.BoolType,
				AutomationActionParameterFieldType:        types.StringType,
				AutomationActionParameterFieldValue:       types.StringType,
			},
		})
	}

	// Handle action type specific configuration
	d := r.mapActionTypeFieldsToState(ctx, automationAction, &model)
	diags.Append(d...)

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

func (r *automationActionResourceFramework) mapActionTypeFieldsToState(ctx context.Context, action *restapi.AutomationAction, model *AutomationActionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Use the common mapping function
	MapActionTypeFieldsToState(ctx, action, model)

	return diags
}

func (r *automationActionResourceFramework) mapTagsToState(ctx context.Context, tags interface{}) (types.List, diag.Diagnostics) {
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

func (r *automationActionResourceFramework) mapInputParametersToState(ctx context.Context, parameters []restapi.Parameter) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	elements := make([]attr.Value, len(parameters))
	for i, param := range parameters {
		paramObj := map[string]attr.Value{
			AutomationActionParameterFieldName:        types.StringValue(param.Name),
			AutomationActionParameterFieldDescription: types.StringValue(param.Description),
			AutomationActionParameterFieldLabel:       types.StringValue(param.Label),
			AutomationActionParameterFieldRequired:    types.BoolValue(param.Required),
			AutomationActionParameterFieldHidden:      types.BoolValue(param.Hidden),
			AutomationActionParameterFieldType:        types.StringValue(param.Type),
			AutomationActionParameterFieldValue:       types.StringValue(param.Value),
		}

		objValue, d := types.ObjectValue(
			map[string]attr.Type{
				AutomationActionParameterFieldName:        types.StringType,
				AutomationActionParameterFieldDescription: types.StringType,
				AutomationActionParameterFieldLabel:       types.StringType,
				AutomationActionParameterFieldRequired:    types.BoolType,
				AutomationActionParameterFieldHidden:      types.BoolType,
				AutomationActionParameterFieldType:        types.StringType,
				AutomationActionParameterFieldValue:       types.StringType,
			},
			paramObj,
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					AutomationActionParameterFieldName:        types.StringType,
					AutomationActionParameterFieldDescription: types.StringType,
					AutomationActionParameterFieldLabel:       types.StringType,
					AutomationActionParameterFieldRequired:    types.BoolType,
					AutomationActionParameterFieldHidden:      types.BoolType,
					AutomationActionParameterFieldType:        types.StringType,
					AutomationActionParameterFieldValue:       types.StringType,
				},
			}), diags
		}

		elements[i] = objValue
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionParameterFieldName:        types.StringType,
				AutomationActionParameterFieldDescription: types.StringType,
				AutomationActionParameterFieldLabel:       types.StringType,
				AutomationActionParameterFieldRequired:    types.BoolType,
				AutomationActionParameterFieldHidden:      types.BoolType,
				AutomationActionParameterFieldType:        types.StringType,
				AutomationActionParameterFieldValue:       types.StringType,
			},
		},
		elements,
	), diags
}

func (r *automationActionResourceFramework) getFieldValue(action *restapi.AutomationAction, fieldName string) string {
	for _, v := range action.Fields {
		if v.Name == fieldName {
			return v.Value
		}
	}
	return ""
}

func (r *automationActionResourceFramework) getBoolFieldValueOrDefault(action *restapi.AutomationAction, fieldName string, defaultValue bool) bool {
	for _, v := range action.Fields {
		if v.Name == fieldName {
			boolValue, err := strconv.ParseBool(v.Value)
			if err != nil {
				return defaultValue
			}
			return boolValue
		}
	}
	return defaultValue
}

func (r *automationActionResourceFramework) mapScriptFieldsToState(ctx context.Context, action *restapi.AutomationAction) (ScriptModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	scriptModel := ScriptModel{
		Content: types.StringValue(r.getFieldValue(action, restapi.ScriptSshFieldName)),
	}

	// Add optional fields if they exist
	interpreter := r.getFieldValue(action, restapi.SubtypeFieldName)
	if interpreter != "" {
		scriptModel.Interpreter = types.StringValue(interpreter)
	} else {
		scriptModel.Interpreter = types.StringNull()
	}

	timeout := r.getFieldValue(action, restapi.TimeoutFieldName)
	if timeout != "" {
		scriptModel.Timeout = types.StringValue(timeout)
	} else {
		scriptModel.Timeout = types.StringNull()
	}

	source := r.getFieldValue(action, "source")
	if source != "" {
		scriptModel.Source = types.StringValue(source)
	} else {
		scriptModel.Source = types.StringNull()
	}

	return scriptModel, diags
}

func (r *automationActionResourceFramework) mapHttpFieldsToState(ctx context.Context, action *restapi.AutomationAction) (HttpModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	httpModel := HttpModel{
		Host:   types.StringValue(r.getFieldValue(action, restapi.HttpHostFieldName)),
		Method: types.StringValue(r.getFieldValue(action, restapi.HttpMethodFieldName)),
	}

	// Auth models will be nil by default (pointers)

	// Add optional fields if they exist
	body := r.getFieldValue(action, restapi.HttpBodyFieldName)
	if body != "" {
		httpModel.Body = types.StringValue(body)
	} else {
		httpModel.Body = types.StringNull()
	}

	httpModel.IgnoreCertErrors = types.BoolValue(
		r.getBoolFieldValueOrDefault(action, restapi.HttpIgnoreCertErrorsFieldName, false),
	)

	timeout := r.getFieldValue(action, restapi.TimeoutFieldName)
	if timeout != "" {
		httpModel.Timeout = types.StringValue(timeout)
	} else {
		httpModel.Timeout = types.StringNull()
	}

	// Language
	language := r.getFieldValue(action, "language")
	if language != "" {
		httpModel.Language = types.StringValue(language)
	} else {
		httpModel.Language = types.StringNull()
	}

	// Content Type
	contentType := r.getFieldValue(action, "content_type")
	if contentType != "" {
		httpModel.ContentType = types.StringValue(contentType)
	} else {
		httpModel.ContentType = types.StringNull()
	}

	// Auth - parse JSON to determine auth type and populate appropriate model
	authData := r.getFieldValue(action, "authen")
	if authData != "" {
		var authMap map[string]interface{}
		err := json.Unmarshal([]byte(authData), &authMap)
		if err == nil {
			authType, _ := authMap["type"].(string)

			// Only create AuthModel if it's an actual auth type (not noAuth)
			switch authType {
			case "basicAuth":
				username, _ := authMap["username"].(string)
				password, _ := authMap["password"].(string)
				httpModel.Auth = &AuthModel{
					BasicAuth: &BasicAuthModel{
						UserName: types.StringValue(username),
						Password: types.StringValue(password),
					},
				}
			case "bearerToken":
				bearerToken, _ := authMap["bearerToken"].(string)
				httpModel.Auth = &AuthModel{
					Token: &BearerTokenModel{
						BearerToken: types.StringValue(bearerToken),
					},
				}
			case "apiKey":
				key, _ := authMap["apiKey"].(string)
				value, _ := authMap["apiKeyValue"].(string)
				location, _ := authMap["apiKeyAddTo"].(string)
				httpModel.Auth = &AuthModel{
					ApiKey: &ApiKeyModel{
						Key:         types.StringValue(key),
						Value:       types.StringValue(value),
						KeyLocation: types.StringValue(location),
					},
				}
				// case "noAuth" or any other type: leave httpModel.Auth as nil
			}
		}
	}
	// If no auth data or noAuth type, Auth model remains nil (default for pointer)

	// Handle headers
	headersData := r.getFieldValue(action, restapi.HttpHeaderFieldName)
	if headersData != "" && headersData != "{}" {
		var headersMap map[string]interface{}
		err := json.Unmarshal([]byte(headersData), &headersMap)
		if err != nil {
			diags.AddError(
				"Error unmarshaling HTTP headers",
				fmt.Sprintf("Failed to unmarshal HTTP headers: %s", err),
			)
			httpModel.Headers = types.MapNull(types.StringType)
			return httpModel, diags
		}

		if len(headersMap) > 0 {
			elements := make(map[string]attr.Value, len(headersMap))
			for k, v := range headersMap {
				if strVal, ok := v.(string); ok {
					elements[k] = types.StringValue(strVal)
				} else {
					elements[k] = types.StringValue(fmt.Sprintf("%v", v))
				}
			}

			httpModel.Headers = types.MapValueMust(types.StringType, elements)
		} else {
			httpModel.Headers = types.MapNull(types.StringType)
		}
	} else {
		httpModel.Headers = types.MapNull(types.StringType)
	}

	return httpModel, diags
}

func (r *automationActionResourceFramework) mapManualFieldsToState(ctx context.Context, action *restapi.AutomationAction) (ManualModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	manualModel := ManualModel{
		Content: types.StringValue(r.getFieldValue(action, "content")),
	}

	return manualModel, diags
}

func (r *automationActionResourceFramework) mapJiraFieldsToState(ctx context.Context, action *restapi.AutomationAction) (JiraModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	jiraModel := JiraModel{}

	if project := r.getFieldValue(action, "project"); project != "" {
		jiraModel.Project = types.StringValue(project)
	} else {
		jiraModel.Project = types.StringNull()
	}

	if operation := r.getFieldValue(action, "ticketActionType"); operation != "" {
		jiraModel.Operation = types.StringValue(operation)
	} else {
		jiraModel.Operation = types.StringNull()
	}

	if issueType := r.getFieldValue(action, "issue_type"); issueType != "" {
		jiraModel.IssueType = types.StringValue(issueType)
	} else {
		jiraModel.IssueType = types.StringNull()
	}

	if description := r.getFieldValue(action, "body"); description != "" {
		jiraModel.Description = types.StringValue(description)
	} else {
		jiraModel.Description = types.StringNull()
	}

	if assignee := r.getFieldValue(action, "assignee"); assignee != "" {
		jiraModel.Assignee = types.StringValue(assignee)
	} else {
		jiraModel.Assignee = types.StringNull()
	}

	if title := r.getFieldValue(action, "summary"); title != "" {
		jiraModel.Title = types.StringValue(title)
	} else {
		jiraModel.Title = types.StringNull()
	}

	if labels := r.getFieldValue(action, "labels"); labels != "" {
		jiraModel.Labels = types.StringValue(labels)
	} else {
		jiraModel.Labels = types.StringNull()
	}

	if comment := r.getFieldValue(action, "comment"); comment != "" {
		jiraModel.Comment = types.StringValue(comment)
	} else {
		jiraModel.Comment = types.StringNull()
	}

	return jiraModel, diags
}

func (r *automationActionResourceFramework) mapGitHubFieldsToState(ctx context.Context, action *restapi.AutomationAction) (GitHubModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	githubModel := GitHubModel{}

	if owner := r.getFieldValue(action, "owner"); owner != "" {
		githubModel.Owner = types.StringValue(owner)
	} else {
		githubModel.Owner = types.StringNull()
	}

	if repo := r.getFieldValue(action, "repo"); repo != "" {
		githubModel.Repo = types.StringValue(repo)
	} else {
		githubModel.Repo = types.StringNull()
	}

	if title := r.getFieldValue(action, "title"); title != "" {
		githubModel.Title = types.StringValue(title)
	} else {
		githubModel.Title = types.StringNull()
	}

	if body := r.getFieldValue(action, "body"); body != "" {
		githubModel.Body = types.StringValue(body)
	} else {
		githubModel.Body = types.StringNull()
	}

	if operation := r.getFieldValue(action, "ticketActionType"); operation != "" {
		githubModel.Operation = types.StringValue(operation)
	} else {
		githubModel.Operation = types.StringNull()
	}

	if assignees := r.getFieldValue(action, "assignees"); assignees != "" {
		githubModel.Assignees = types.StringValue(assignees)
	} else {
		githubModel.Assignees = types.StringNull()
	}

	if labels := r.getFieldValue(action, "labels"); labels != "" {
		githubModel.Labels = types.StringValue(labels)
	} else {
		githubModel.Labels = types.StringNull()
	}

	if comment := r.getFieldValue(action, "comment"); comment != "" {
		githubModel.Comment = types.StringValue(comment)
	} else {
		githubModel.Comment = types.StringNull()
	}

	return githubModel, diags
}

func (r *automationActionResourceFramework) mapDocLinkFieldsToState(ctx context.Context, action *restapi.AutomationAction) (DocLinkModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	docLinkModel := DocLinkModel{
		Url: types.StringValue(r.getFieldValue(action, "url")),
	}

	return docLinkModel, diags
}

func (r *automationActionResourceFramework) mapGitLabFieldsToState(ctx context.Context, action *restapi.AutomationAction) (GitLabModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	gitlabModel := GitLabModel{}

	if projectId := r.getFieldValue(action, "projectId"); projectId != "" {
		gitlabModel.ProjectId = types.StringValue(projectId)
	} else {
		gitlabModel.ProjectId = types.StringNull()
	}

	if title := r.getFieldValue(action, "title"); title != "" {
		gitlabModel.Title = types.StringValue(title)
	} else {
		gitlabModel.Title = types.StringNull()
	}

	if description := r.getFieldValue(action, "body"); description != "" {
		gitlabModel.Description = types.StringValue(description)
	} else {
		gitlabModel.Description = types.StringNull()
	}

	if operation := r.getFieldValue(action, "ticketActionType"); operation != "" {
		gitlabModel.Operation = types.StringValue(operation)
	} else {
		gitlabModel.Operation = types.StringNull()
	}

	if labels := r.getFieldValue(action, "labels"); labels != "" {
		gitlabModel.Labels = types.StringValue(labels)
	} else {
		gitlabModel.Labels = types.StringNull()
	}

	if issueType := r.getFieldValue(action, "issue_type"); issueType != "" {
		gitlabModel.IssueType = types.StringValue(issueType)
	} else {
		gitlabModel.IssueType = types.StringNull()
	}

	if comment := r.getFieldValue(action, "comment"); comment != "" {
		gitlabModel.Comment = types.StringValue(comment)
	} else {
		gitlabModel.Comment = types.StringNull()
	}

	return gitlabModel, diags
}

func (r *automationActionResourceFramework) mapAnsibleFieldsToState(ctx context.Context, action *restapi.AutomationAction) (AnsibleModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	ansibleModel := AnsibleModel{}

	if workflowId := r.getFieldValue(action, "workflowId"); workflowId != "" {
		ansibleModel.WorkflowId = types.StringValue(workflowId)
	} else {
		ansibleModel.WorkflowId = types.StringNull()
	}

	if ansibleUrl := r.getFieldValue(action, "ansibleUrl"); ansibleUrl != "" {
		ansibleModel.AnsibleUrl = types.StringValue(ansibleUrl)
	} else {
		ansibleModel.AnsibleUrl = types.StringNull()
	}

	if hostId := r.getFieldValue(action, "hostId"); hostId != "" {
		ansibleModel.HostId = types.StringValue(hostId)
	} else {
		ansibleModel.HostId = types.StringNull()
	}

	if playbookId := r.getFieldValue(action, "playbookId"); playbookId != "" {
		ansibleModel.PlaybookId = types.StringValue(playbookId)
	} else {
		ansibleModel.PlaybookId = types.StringNull()
	}

	if playbookFileName := r.getFieldValue(action, "playbookFileName"); playbookFileName != "" {
		ansibleModel.PlaybookFileName = types.StringValue(playbookFileName)
	} else {
		ansibleModel.PlaybookFileName = types.StringNull()
	}

	return ansibleModel, diags
}

func (r *automationActionResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AutomationAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model AutomationActionModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Determine action type and map fields
	actionType, fields, d := r.mapActionTypeAndFields(ctx, model)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Map input parameters
	inputParameters, d := r.mapInputParametersFromState(ctx, model)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Map tags
	tags, d := r.mapTagsFromState(ctx, model)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Create the API model
	id := ""
	if !model.ID.IsNull() {
		id = model.ID.ValueString()
	}

	return &restapi.AutomationAction{
		ID:              id,
		Name:            model.Name.ValueString(),
		Description:     model.Description.ValueString(),
		Type:            actionType,
		Tags:            tags,
		Fields:          fields,
		InputParameters: inputParameters,
	}, diags
}

func (r *automationActionResourceFramework) mapActionTypeAndFields(ctx context.Context, model AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	var actionType string
	var fields []restapi.Field

	// Check if script configuration is provided
	if model.Script != nil && !model.Script.Content.IsNull() {
		actionType = ActionTypeScript
		scriptModel := *model.Script

		// Map script fields
		fields = make([]restapi.Field, 0)

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
	} else if model.Http != nil && !model.Http.Host.IsNull() {
		actionType = ActionTypeHttp
		httpModel := *model.Http

		// Map HTTP fields
		fields = make([]restapi.Field, 0)

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
				Value:       strconv.FormatBool(httpModel.IgnoreCertErrors.ValueBool()),
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

		// Auth is optional - serialize as JSON based on auth type
		var authValue string
		if httpModel.Auth != nil {
			if httpModel.Auth.BasicAuth != nil && !httpModel.Auth.BasicAuth.UserName.IsNull() {
				// Basic Auth format: {"type":"basicAuth","username":"@@userName@@","password":"@@password@@"}
				authMap := map[string]string{
					"type":     "basicAuth",
					"username": httpModel.Auth.BasicAuth.UserName.ValueString(),
					"password": httpModel.Auth.BasicAuth.Password.ValueString(),
				}
				authJson, err := json.Marshal(authMap)
				if err != nil {
					diags.AddError(
						"Error marshaling basic auth",
						fmt.Sprintf("Failed to marshal basic auth: %s", err),
					)
					return "", nil, diags
				}
				authValue = string(authJson)
			} else if httpModel.Auth.Token != nil && !httpModel.Auth.Token.BearerToken.IsNull() {
				// Bearer Token format: {"type":"bearerToken","bearerToken":"@@bearerToken@@"}
				authMap := map[string]string{
					"type":        "bearerToken",
					"bearerToken": httpModel.Auth.Token.BearerToken.ValueString(),
				}
				authJson, err := json.Marshal(authMap)
				if err != nil {
					diags.AddError(
						"Error marshaling bearer token",
						fmt.Sprintf("Failed to marshal bearer token: %s", err),
					)
					return "", nil, diags
				}
				authValue = string(authJson)
			} else if httpModel.Auth.ApiKey != nil && !httpModel.Auth.ApiKey.Key.IsNull() {
				// API Key format: {"type":"apiKey","apiKey":"authorization","apiKeyValue":"somKey","apiKeyAddTo":"header"}
				authMap := map[string]string{
					"type":        "apiKey",
					"apiKey":      httpModel.Auth.ApiKey.Key.ValueString(),
					"apiKeyValue": httpModel.Auth.ApiKey.Value.ValueString(),
					"apiKeyAddTo": httpModel.Auth.ApiKey.KeyLocation.ValueString(),
				}
				authJson, err := json.Marshal(authMap)
				if err != nil {
					diags.AddError(
						"Error marshaling API key",
						fmt.Sprintf("Failed to marshal API key: %s", err),
					)
					return "", nil, diags
				}
				authValue = string(authJson)
			} else {
				// No Auth format: {"type":"noAuth"}
				authMap := map[string]string{
					"type": "noAuth",
				}
				authJson, err := json.Marshal(authMap)
				if err != nil {
					diags.AddError(
						"Error marshaling no auth",
						fmt.Sprintf("Failed to marshal no auth: %s", err),
					)
					return "", nil, diags
				}
				authValue = string(authJson)
			}
		} else {
			// No Auth format: {"type":"noAuth"}
			authMap := map[string]string{
				"type": "noAuth",
			}
			authJson, err := json.Marshal(authMap)
			if err != nil {
				diags.AddError(
					"Error marshaling no auth",
					fmt.Sprintf("Failed to marshal no auth: %s", err),
				)
				return "", nil, diags
			}
			authValue = string(authJson)
		}

		fields = append(fields, restapi.Field{
			Name:        "authen",
			Description: "Authentication for the HTTPS request",
			Value:       authValue,
			Encoding:    AsciiEncoding,
			Secured:     false,
		})

		// Headers are optional
		if !httpModel.Headers.IsNull() {
			headersMap := make(map[string]string)
			diags.Append(httpModel.Headers.ElementsAs(ctx, &headersMap, false)...)
			if diags.HasError() {
				return "", nil, diags
			}

			headersJson, err := json.Marshal(headersMap)
			if err != nil {
				diags.AddError(
					"Error marshaling HTTP headers",
					fmt.Sprintf("Failed to marshal HTTP headers: %s", err),
				)
				return "", nil, diags
			}

			fields = append(fields, restapi.Field{
				Name:        restapi.HttpHeaderFieldName,
				Description: restapi.HttpHeaderFieldDescription,
				Value:       string(headersJson),
				Encoding:    AsciiEncoding,
				Secured:     false,
			})
		}
	} else if model.Manual != nil && !model.Manual.Content.IsNull() {
		actionType = "MANUAL"
		fields = make([]restapi.Field, 0)
		fields = append(fields, restapi.Field{
			Name:        "content",
			Description: "Content for manual action",
			Value:       model.Manual.Content.ValueString(),
			Encoding:    AsciiEncoding,
			Secured:     false,
		})
	} else if model.Jira != nil && !model.Jira.Project.IsNull() {
		actionType = "JIRA"
		fields = make([]restapi.Field, 0)
		if !model.Jira.Project.IsNull() {
			fields = append(fields, restapi.Field{Name: "project", Description: "jira project", Value: model.Jira.Project.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.Operation.IsNull() {
			fields = append(fields, restapi.Field{Name: "ticketActionType", Description: "jira ticket type", Value: model.Jira.Operation.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.IssueType.IsNull() {
			fields = append(fields, restapi.Field{Name: "issue_type", Description: "jira issue type", Value: model.Jira.IssueType.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.Description.IsNull() {
			fields = append(fields, restapi.Field{Name: "body", Description: "jira issue description", Value: model.Jira.Description.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.Assignee.IsNull() {
			fields = append(fields, restapi.Field{Name: "assignee", Description: "jira issue assignee", Value: model.Jira.Assignee.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.Title.IsNull() {
			fields = append(fields, restapi.Field{Name: "summary", Description: "jira issue summary", Value: model.Jira.Title.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.Labels.IsNull() {
			fields = append(fields, restapi.Field{Name: "labels", Description: "jira issue labels", Value: model.Jira.Labels.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Jira.Comment.IsNull() {
			fields = append(fields, restapi.Field{Name: "comment", Description: "jira issue comment", Value: model.Jira.Comment.ValueString(), Encoding: AsciiEncoding})
		}
	} else if model.GitHub != nil && !model.GitHub.Owner.IsNull() {
		actionType = "GITHUB"
		fields = make([]restapi.Field, 0)
		if !model.GitHub.Owner.IsNull() {
			fields = append(fields, restapi.Field{Name: "owner", Description: "github issue owner/repo", Value: model.GitHub.Owner.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Repo.IsNull() {
			fields = append(fields, restapi.Field{Name: "repo", Description: "github issue repo", Value: model.GitHub.Repo.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Title.IsNull() {
			fields = append(fields, restapi.Field{Name: "title", Description: "github issue title", Value: model.GitHub.Title.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Body.IsNull() {
			fields = append(fields, restapi.Field{Name: "body", Description: "github issue body", Value: model.GitHub.Body.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Operation.IsNull() {
			fields = append(fields, restapi.Field{Name: "ticketType", Description: "github issue type", Value: model.GitHub.Operation.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Assignees.IsNull() {
			fields = append(fields, restapi.Field{Name: "assignees", Description: "github issue assignees", Value: model.GitHub.Assignees.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Labels.IsNull() {
			fields = append(fields, restapi.Field{Name: "labels", Description: "github issue labels", Value: model.GitHub.Labels.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitHub.Comment.IsNull() {
			fields = append(fields, restapi.Field{Name: "comment", Description: "github issue comment", Value: model.GitHub.Comment.ValueString(), Encoding: AsciiEncoding})
		}
	} else if model.DocLink != nil && !model.DocLink.Url.IsNull() {
		actionType = "DOC_LINK"
		fields = make([]restapi.Field, 0)
		fields = append(fields, restapi.Field{
			Name:        "url",
			Description: "URL to remediation documentation",
			Value:       model.DocLink.Url.ValueString(),
			Encoding:    UTF8Encoding,
		})
	} else if model.GitLab != nil && !model.GitLab.ProjectId.IsNull() {
		actionType = "GITLAB"
		fields = make([]restapi.Field, 0)
		if !model.GitLab.ProjectId.IsNull() {
			fields = append(fields, restapi.Field{Name: "projectId", Description: "gitlab projectId", Value: model.GitLab.ProjectId.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitLab.Title.IsNull() {
			fields = append(fields, restapi.Field{Name: "title", Description: "gitlab issue title", Value: model.GitLab.Title.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitLab.Description.IsNull() {
			fields = append(fields, restapi.Field{Name: "body", Description: "gitlab issue description", Value: model.GitLab.Description.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitLab.Operation.IsNull() {
			fields = append(fields, restapi.Field{Name: "ticketActionType", Description: "gitlab ticket type", Value: model.GitLab.Operation.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitLab.Labels.IsNull() {
			fields = append(fields, restapi.Field{Name: "labels", Description: "github issue labels", Value: model.GitLab.Labels.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitLab.IssueType.IsNull() {
			fields = append(fields, restapi.Field{Name: "issue_type", Description: "gitlab issue type", Value: model.GitLab.IssueType.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.GitLab.Labels.IsNull() {
			fields = append(fields, restapi.Field{Name: "comment", Description: "gitlab issue comment", Value: model.GitLab.Comment.ValueString(), Encoding: AsciiEncoding})
		}
	} else if model.Ansible != nil && !model.Ansible.WorkflowId.IsNull() {
		actionType = "ANSIBLE"
		fields = make([]restapi.Field, 0)
		if !model.Ansible.WorkflowId.IsNull() {
			fields = append(fields, restapi.Field{Name: "workflowId", Description: "The workflow ID", Value: model.Ansible.WorkflowId.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Ansible.AnsibleUrl.IsNull() {
			fields = append(fields, restapi.Field{Name: "ansibleUrl", Description: "The ansible url", Value: model.Ansible.AnsibleUrl.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Ansible.HostId.IsNull() {
			fields = append(fields, restapi.Field{Name: "hostId", Description: "The host ID from which this action is created", Value: model.Ansible.HostId.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Ansible.PlaybookId.IsNull() {
			fields = append(fields, restapi.Field{Name: "playbookId", Description: "The playbook ID", Value: model.Ansible.PlaybookId.ValueString(), Encoding: AsciiEncoding})
		}
		if !model.Ansible.PlaybookFileName.IsNull() {
			fields = append(fields, restapi.Field{Name: "playbookFileName", Description: "The playbook filename", Value: model.Ansible.PlaybookFileName.ValueString(), Encoding: AsciiEncoding})
		}
	} else {
		diags.AddError(
			"Invalid action configuration",
			"One of script, http, manual, jira, github, doclink, gitlab, or ansible configuration must be provided",
		)
		return "", nil, diags
	}

	return actionType, fields, diags
}

func (r *automationActionResourceFramework) mapInputParametersFromState(ctx context.Context, model AutomationActionModel) ([]restapi.Parameter, diag.Diagnostics) {
	var diags diag.Diagnostics
	var parameters []restapi.Parameter

	if model.InputParameter.IsNull() {
		return parameters, diags
	}

	var parameterModels []ParameterModel
	diags.Append(model.InputParameter.ElementsAs(ctx, &parameterModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	parameters = make([]restapi.Parameter, len(parameterModels))
	for i, param := range parameterModels {
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

func (r *automationActionResourceFramework) mapTagsFromState(ctx context.Context, model AutomationActionModel) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model.Tags.IsNull() {
		return nil, diags
	}

	var tags []string
	diags.Append(model.Tags.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	return tags, diags
}

// Made with Bob
