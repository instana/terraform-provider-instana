package automationaction

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/shared"
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

// ============================================================================
// Resource Factory
// ============================================================================

// NewAutomationActionResourceHandle creates the resource handle for Automation Actions
func NewAutomationActionResourceHandle() resourcehandle.ResourceHandle[*restapi.AutomationAction] {
	return &automationActionResource{
		metaData: resourcehandle.ResourceMetaData{
			ResourceName: ResourceInstanaAutomationAction,
			Schema: schema.Schema{
				Description: AutomationActionDescResource,
				Attributes: map[string]schema.Attribute{
					AutomationActionFieldID: schema.StringAttribute{
						Computed:    true,
						Description: AutomationActionDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					AutomationActionFieldName: schema.StringAttribute{
						Required:    true,
						Description: AutomationActionDescName,
					},
					AutomationActionFieldDescription: schema.StringAttribute{
						Required:    true,
						Description: AutomationActionDescDescription,
					},
					AutomationActionFieldTags: schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: AutomationActionDescTags,
					},
					AutomationActionFieldScript: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescScript,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldContent: schema.StringAttribute{
								Required:    true,
								Description: AutomationActionDescScriptContent,
							},
							AutomationActionFieldInterpreter: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescScriptInterpreter,
							},
							AutomationActionFieldTimeout: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescScriptTimeout,
							},
							AutomationActionFieldSource: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescScriptSource,
							},
						},
					},
					AutomationActionFieldHttp: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescHttp,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldHost: schema.StringAttribute{
								Required:    true,
								Description: AutomationActionDescHttpHost,
							},
							AutomationActionFieldMethod: schema.StringAttribute{
								Required:    true,
								Description: AutomationActionDescHttpMethod,
								Validators: []validator.String{
									stringvalidator.OneOf(HTTPMethodGET, HTTPMethodPOST, HTTPMethodPUT, HTTPMethodDELETE),
								},
							},
							AutomationActionFieldBody: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescHttpBody,
							},
							AutomationActionFieldHeaders: schema.MapAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: AutomationActionDescHttpHeaders,
							},
							AutomationActionFieldIgnoreCertErrors: schema.BoolAttribute{
								Optional:    true,
								Description: AutomationActionDescHttpIgnoreCertErrors,
							},
							AutomationActionFieldTimeout: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescHttpTimeout,
							},
							AutomationActionFieldLanguage: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescHttpLanguage,
							},
							AutomationActionFieldContentType: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescHttpContentType,
							},
							AutomationActionFieldAuth: schema.SingleNestedAttribute{
								Optional:    true,
								Description: AutomationActionDescHttpAuth,
								Attributes: map[string]schema.Attribute{
									AutomationActionFieldBasicAuth: schema.SingleNestedAttribute{
										Optional:    true,
										Description: AutomationActionDescHttpAuthBasic,
										Attributes: map[string]schema.Attribute{
											AutomationActionFieldUsername: schema.StringAttribute{
												Required:    true,
												Description: AutomationActionDescHttpAuthBasicUsername,
											},
											AutomationActionFieldPassword: schema.StringAttribute{
												Required:    true,
												Description: AutomationActionDescHttpAuthBasicPassword,
											},
										},
									},
									AutomationActionFieldToken: schema.SingleNestedAttribute{
										Optional:    true,
										Description: AutomationActionDescHttpAuthToken,
										Attributes: map[string]schema.Attribute{
											AutomationActionFieldBearerToken: schema.StringAttribute{
												Required:    true,
												Description: AutomationActionDescHttpAuthBearerToken,
											},
										},
									},
									AutomationActionFieldApiKey: schema.SingleNestedAttribute{
										Optional:    true,
										Description: AutomationActionDescHttpAuthApiKey,
										Attributes: map[string]schema.Attribute{
											AutomationActionFieldKey: schema.StringAttribute{
												Required:    true,
												Description: AutomationActionDescHttpAuthApiKeyKey,
											},
											AutomationActionFieldValue: schema.StringAttribute{
												Required:    true,
												Description: AutomationActionDescHttpAuthApiKeyValue,
											},
											AutomationActionFieldKeyLocation: schema.StringAttribute{
												Required:    true,
												Description: AutomationActionDescHttpAuthApiKeyLocation,
											},
										},
									},
								},
							},
						},
					},
					AutomationActionFieldManual: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescManual,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldContent: schema.StringAttribute{
								Required:    true,
								Description: AutomationActionDescManualContent,
							},
						},
					},
					AutomationActionFieldJira: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescJira,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldProject: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraProject,
							},
							AutomationActionFieldOperation: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraOperation,
							},
							AutomationActionFieldIssueType: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraIssueType,
							},
							AutomationActionFieldDescription: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraDescription,
							},
							AutomationActionFieldAssignee: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraAssignee,
							},
							AutomationActionFieldTitle: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraTitle,
							},
							AutomationActionFieldLabels: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraLabels,
							},
							AutomationActionFieldComment: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescJiraComment,
							},
						},
					},
					AutomationActionFieldGitHub: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescGitHub,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldOwner: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubOwner,
							},
							AutomationActionFieldRepo: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubRepo,
							},
							AutomationActionFieldTitle: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubTitle,
							},
							AutomationActionFieldBody: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubBody,
							},
							AutomationActionFieldOperation: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubOperation,
							},
							AutomationActionFieldAssignees: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubAssignees,
							},
							AutomationActionFieldLabels: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubLabels,
							},
							AutomationActionFieldComment: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitHubComment,
							},
						},
					},
					AutomationActionFieldDocLink: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescDocLink,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldUrl: schema.StringAttribute{
								Required:    true,
								Description: AutomationActionDescDocLinkUrl,
							},
						},
					},
					AutomationActionFieldGitLab: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescGitLab,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldProjectId: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabProjectId,
							},
							AutomationActionFieldTitle: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabTitle,
							},
							AutomationActionFieldDescription: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabDescription,
							},
							AutomationActionFieldOperation: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabOperation,
							},
							AutomationActionFieldLabels: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabLabels,
							},
							AutomationActionFieldIssueType: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabIssueType,
							},
							AutomationActionFieldComment: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescGitLabComment,
							},
						},
					},
					AutomationActionFieldAnsible: schema.SingleNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescAnsible,
						Attributes: map[string]schema.Attribute{
							AutomationActionFieldWorkflowId: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescAnsibleWorkflowId,
							},
							AutomationActionFieldPlaybookId: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescAnsiblePlaybookId,
							},
							AutomationActionFieldPlaybookFileName: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescAnsiblePlaybookFileName,
							},
							AutomationActionFieldUrl: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescAnsibleUrl,
							},
							AutomationActionFieldHostId: schema.StringAttribute{
								Optional:    true,
								Description: AutomationActionDescAnsibleHostId,
							},
						},
					},
					AutomationActionFieldInputParameter: schema.ListNestedAttribute{
						Optional:    true,
						Description: AutomationActionDescInputParameter,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								AutomationActionParameterFieldName: schema.StringAttribute{
									Required:    true,
									Description: AutomationActionDescParameterName,
								},
								AutomationActionParameterFieldDescription: schema.StringAttribute{
									Required:    true,
									Description: AutomationActionDescParameterDescription,
								},
								AutomationActionParameterFieldLabel: schema.StringAttribute{
									Required:    true,
									Description: AutomationActionDescParameterLabel,
								},
								AutomationActionParameterFieldRequired: schema.BoolAttribute{
									Required:    true,
									Description: AutomationActionDescParameterRequired,
								},
								AutomationActionParameterFieldHidden: schema.BoolAttribute{
									Required:    true,
									Description: AutomationActionDescParameterHidden,
								},
								AutomationActionParameterFieldType: schema.StringAttribute{
									Required:    true,
									Description: AutomationActionDescParameterType,
									Validators: []validator.String{
										stringvalidator.OneOf(
											AutomationActionParameterTypeStatic,
											AutomationActionParameterTypeDynamic,
											AutomationActionParameterTypeVault,
										),
									},
								},
								AutomationActionParameterFieldValue: schema.StringAttribute{
									Required:    true,
									Description: AutomationActionDescParameterValue,
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

// ============================================================================
// Resource Implementation
// ============================================================================

type automationActionResource struct {
	metaData resourcehandle.ResourceMetaData
}

// MetaData returns the resource metadata
func (r *automationActionResource) MetaData() *resourcehandle.ResourceMetaData {
	return &r.metaData
}

// GetRestResource returns the REST resource for automation actions
func (r *automationActionResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AutomationAction] {
	return api.AutomationActions()
}

// SetComputedFields sets computed fields in the plan (none for this resource)
func (r *automationActionResource) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

// ============================================================================
// API to State Mapping
// ============================================================================

// UpdateState converts API data object to Terraform state
func (r *automationActionResource) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, automationAction *restapi.AutomationAction) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API response
	model := shared.AutomationActionModel{
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
		model.InputParameter = shared.MapInputParametersToState(ctx, automationAction.InputParameters)
	} else {
		model.InputParameter = nil
	}

	// Handle action type specific configuration
	d := r.mapActionTypeFieldsToState(ctx, automationAction, &model)
	diags.Append(d...)

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
	return diags
}

// mapActionTypeFieldsToState maps action type specific fields from API to state
func (r *automationActionResource) mapActionTypeFieldsToState(ctx context.Context, action *restapi.AutomationAction, model *shared.AutomationActionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Use the common mapping function
	shared.MapActionTypeFieldsToState(ctx, action, model)

	return diags
}

// mapTagsToState converts tags from API format to Terraform state format
func (r *automationActionResource) mapTagsToState(ctx context.Context, tags interface{}) (types.List, diag.Diagnostics) {
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
					AutomationActionErrMappingTags,
					fmt.Sprintf(AutomationActionErrTagNotString, i),
				)
				return types.ListNull(types.StringType), diags
			}
		}
		return types.ListValueMust(types.StringType, elements), diags
	default:
		diags.AddError(
			AutomationActionErrMappingTags,
			AutomationActionErrTagsFormat,
		)
		return types.ListNull(types.StringType), diags
	}
}

// ============================================================================
// State to API Mapping
// ============================================================================

// MapStateToDataObject converts Terraform state to API data object
func (r *automationActionResource) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.AutomationAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model shared.AutomationActionModel

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

// ============================================================================
// Action Type Mappers
// ============================================================================

// mapActionTypeAndFields determines the action type and maps corresponding fields
func (r *automationActionResource) mapActionTypeAndFields(ctx context.Context, model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check each action type and delegate to specific mapper
	if model.Script != nil && !model.Script.Content.IsNull() {
		return r.mapScriptAction(model)
	} else if model.Http != nil && !model.Http.Host.IsNull() {
		return r.mapHttpAction(ctx, model)
	} else if model.Manual != nil && !model.Manual.Content.IsNull() {
		return r.mapManualAction(model)
	} else if model.Jira != nil && !model.Jira.Project.IsNull() {
		return r.mapJiraAction(model)
	} else if model.GitHub != nil && !model.GitHub.Owner.IsNull() {
		return r.mapGitHubAction(model)
	} else if model.DocLink != nil && !model.DocLink.Url.IsNull() {
		return r.mapDocLinkAction(model)
	} else if model.GitLab != nil && !model.GitLab.ProjectId.IsNull() {
		return r.mapGitLabAction(model)
	} else if model.Ansible != nil && (!model.Ansible.WorkflowId.IsNull() || !model.Ansible.PlaybookId.IsNull() || !model.Ansible.PlaybookFileName.IsNull() || !model.Ansible.AnsibleUrl.IsNull() || !model.Ansible.HostId.IsNull()) {
		return r.mapAnsibleAction(model)
	}

	diags.AddError(
		AutomationActionErrInvalidConfig,
		AutomationActionErrInvalidConfigMsg,
	)
	return "", nil, diags
}

// mapScriptAction maps script action configuration to API fields
func (r *automationActionResource) mapScriptAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	scriptModel := *model.Script
	fields := make([]restapi.Field, 0)

	// Content is required
	fields = append(fields, restapi.Field{
		Name:        restapi.ScriptSshFieldName,
		Description: restapi.ScriptSshFieldDescription,
		Value:       scriptModel.Content.ValueString(),
		Encoding:    shared.Base64Encoding,
		Secured:     false,
	})

	// Interpreter is optional
	if !scriptModel.Interpreter.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.SubtypeFieldName,
			Description: restapi.SubtypeFieldDescription,
			Value:       scriptModel.Interpreter.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	// Timeout is optional
	if !scriptModel.Timeout.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.TimeoutFieldName,
			Description: restapi.TimeoutFieldDescription,
			Value:       scriptModel.Timeout.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	// Source is optional
	if !scriptModel.Source.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldSource,
			Description: AutomationActionDescFieldSource,
			Value:       scriptModel.Source.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	return shared.ActionTypeScript, fields, diags
}

// mapHttpAction maps HTTP action configuration to API fields
func (r *automationActionResource) mapHttpAction(ctx context.Context, model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	httpModel := *model.Http
	fields := make([]restapi.Field, 0)

	// Map basic HTTP fields
	fields = r.mapHttpBasicFields(httpModel, fields)

	// Map authentication
	authFields, authDiags := r.mapHttpAuth(httpModel)
	diags.Append(authDiags...)
	if diags.HasError() {
		return "", nil, diags
	}
	fields = append(fields, authFields...)

	// Map headers
	headerFields, headerDiags := r.mapHttpHeaders(ctx, httpModel)
	diags.Append(headerDiags...)
	if diags.HasError() {
		return "", nil, diags
	}
	fields = append(fields, headerFields...)

	return shared.ActionTypeHttp, fields, diags
}

// mapHttpBasicFields maps basic HTTP fields (host, method, body, etc.)
func (r *automationActionResource) mapHttpBasicFields(httpModel shared.HttpModel, fields []restapi.Field) []restapi.Field {
	// Host and method are required
	fields = append(fields, restapi.Field{
		Name:        restapi.HttpHostFieldName,
		Description: restapi.HttpHostFieldDescription,
		Value:       httpModel.Host.ValueString(),
		Encoding:    shared.AsciiEncoding,
		Secured:     false,
	})

	fields = append(fields, restapi.Field{
		Name:        restapi.HttpMethodFieldName,
		Description: restapi.HttpMethodFieldDescription,
		Value:       httpModel.Method.ValueString(),
		Encoding:    shared.AsciiEncoding,
		Secured:     false,
	})

	// Body is optional
	if !httpModel.Body.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.HttpBodyFieldName,
			Description: restapi.HttpBodyFieldDescription,
			Value:       httpModel.Body.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	// IgnoreCertErrors is optional
	if !httpModel.IgnoreCertErrors.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.HttpIgnoreCertErrorsFieldName,
			Description: restapi.HttpIgnoreCertErrorsFieldDescription,
			Value:       strconv.FormatBool(httpModel.IgnoreCertErrors.ValueBool()),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	// Timeout is optional
	if !httpModel.Timeout.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        restapi.TimeoutFieldName,
			Description: restapi.TimeoutFieldDescription,
			Value:       httpModel.Timeout.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	// Language is optional
	if !httpModel.Language.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldLanguage,
			Description: AutomationActionDescFieldLanguage,
			Value:       httpModel.Language.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	// Content Type is optional
	if !httpModel.ContentType.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldContentType,
			Description: AutomationActionDescFieldContentType,
			Value:       httpModel.ContentType.ValueString(),
			Encoding:    shared.AsciiEncoding,
			Secured:     false,
		})
	}

	return fields
}

// mapHttpAuth maps HTTP authentication configuration to API fields
func (r *automationActionResource) mapHttpAuth(httpModel shared.HttpModel) ([]restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	var authValue string

	if httpModel.Auth != nil {
		if httpModel.Auth.BasicAuth != nil && !httpModel.Auth.BasicAuth.UserName.IsNull() {
			authValue = r.marshalBasicAuth(httpModel.Auth.BasicAuth, &diags)
		} else if httpModel.Auth.Token != nil && !httpModel.Auth.Token.BearerToken.IsNull() {
			authValue = r.marshalBearerToken(httpModel.Auth.Token, &diags)
		} else if httpModel.Auth.ApiKey != nil && !httpModel.Auth.ApiKey.Key.IsNull() {
			authValue = r.marshalApiKey(httpModel.Auth.ApiKey, &diags)
		} else {
			authValue = r.marshalNoAuth(&diags)
		}
	} else {
		authValue = r.marshalNoAuth(&diags)
	}

	if diags.HasError() {
		return nil, diags
	}

	return []restapi.Field{{
		Name:        AutomationActionAPIFieldAuthen,
		Description: AutomationActionDescFieldAuthen,
		Value:       authValue,
		Encoding:    shared.AsciiEncoding,
		Secured:     false,
	}}, diags
}

// marshalBasicAuth marshals basic authentication to JSON string
func (r *automationActionResource) marshalBasicAuth(basicAuth *shared.BasicAuthModel, diags *diag.Diagnostics) string {
	authMap := map[string]string{
		AuthJSONFieldType:     AutomationActionAuthTypeBasicAuth,
		AuthJSONFieldUsername: basicAuth.UserName.ValueString(),
		AuthJSONFieldPassword: basicAuth.Password.ValueString(),
	}
	authJson, err := json.Marshal(authMap)
	if err != nil {
		diags.AddError(
			AutomationActionErrMarshalBasicAuth,
			fmt.Sprintf(AutomationActionErrMarshalBasicAuthFailed, err),
		)
		return ""
	}
	return string(authJson)
}

// marshalBearerToken marshals bearer token authentication to JSON string
func (r *automationActionResource) marshalBearerToken(token *shared.BearerTokenModel, diags *diag.Diagnostics) string {
	authMap := map[string]string{
		AuthJSONFieldType:        AutomationActionAuthTypeBearerToken,
		AuthJSONFieldBearerToken: token.BearerToken.ValueString(),
	}
	authJson, err := json.Marshal(authMap)
	if err != nil {
		diags.AddError(
			AutomationActionErrMarshalBearerToken,
			fmt.Sprintf(AutomationActionErrMarshalBearerTokenFailed, err),
		)
		return ""
	}
	return string(authJson)
}

// marshalApiKey marshals API key authentication to JSON string
func (r *automationActionResource) marshalApiKey(apiKey *shared.ApiKeyModel, diags *diag.Diagnostics) string {
	authMap := map[string]string{
		AuthJSONFieldType:        AutomationActionAuthTypeApiKey,
		AuthJSONFieldAPIKey:      apiKey.Key.ValueString(),
		AuthJSONFieldAPIKeyValue: apiKey.Value.ValueString(),
		AuthJSONFieldAPIKeyAddTo: apiKey.KeyLocation.ValueString(),
	}
	authJson, err := json.Marshal(authMap)
	if err != nil {
		diags.AddError(
			AutomationActionErrMarshalApiKey,
			fmt.Sprintf(AutomationActionErrMarshalApiKeyFailed, err),
		)
		return ""
	}
	return string(authJson)
}

// marshalNoAuth marshals no authentication to JSON string
func (r *automationActionResource) marshalNoAuth(diags *diag.Diagnostics) string {
	authMap := map[string]string{
		AuthJSONFieldType: AutomationActionAuthTypeNoAuth,
	}
	authJson, err := json.Marshal(authMap)
	if err != nil {
		diags.AddError(
			AutomationActionErrMarshalNoAuth,
			fmt.Sprintf(AutomationActionErrMarshalNoAuthFailed, err),
		)
		return ""
	}
	return string(authJson)
}

// mapHttpHeaders maps HTTP headers to API fields
func (r *automationActionResource) mapHttpHeaders(ctx context.Context, httpModel shared.HttpModel) ([]restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics

	if httpModel.Headers.IsNull() {
		return nil, diags
	}

	headersMap := make(map[string]string)
	diags.Append(httpModel.Headers.ElementsAs(ctx, &headersMap, false)...)
	if diags.HasError() {
		return nil, diags
	}

	headersJson, err := json.Marshal(headersMap)
	if err != nil {
		diags.AddError(
			AutomationActionErrMarshalHeaders,
			fmt.Sprintf(AutomationActionErrMarshalHeadersFailed, err),
		)
		return nil, diags
	}

	return []restapi.Field{{
		Name:        restapi.HttpHeaderFieldName,
		Description: restapi.HttpHeaderFieldDescription,
		Value:       string(headersJson),
		Encoding:    shared.AsciiEncoding,
		Secured:     false,
	}}, diags
}

// mapManualAction maps manual action configuration to API fields
func (r *automationActionResource) mapManualAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := []restapi.Field{{
		Name:        AutomationActionFieldContent,
		Description: AutomationActionDescAPIFieldContent,
		Value:       model.Manual.Content.ValueString(),
		Encoding:    shared.AsciiEncoding,
		Secured:     false,
	}}
	return AutomationActionTypeManual, fields, diags
}

// mapJiraAction maps Jira action configuration to API fields
func (r *automationActionResource) mapJiraAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := make([]restapi.Field, 0)

	if !model.Jira.Project.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldProject,
			Description: AutomationActionDescAPIFieldProject,
			Value:       model.Jira.Project.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.Operation.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldTicketActionType,
			Description: AutomationActionDescAPIFieldTicketActionType,
			Value:       model.Jira.Operation.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.IssueType.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldIssueType,
			Description: AutomationActionDescAPIFieldIssueType,
			Value:       model.Jira.IssueType.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.Description.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldBody,
			Description: AutomationActionDescAPIFieldBody,
			Value:       model.Jira.Description.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.Assignee.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldAssignee,
			Description: AutomationActionDescAPIFieldAssignee,
			Value:       model.Jira.Assignee.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.Title.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldSummary,
			Description: AutomationActionDescAPIFieldSummary,
			Value:       model.Jira.Title.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.Labels.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldLabels,
			Description: AutomationActionDescAPIFieldLabels,
			Value:       model.Jira.Labels.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Jira.Comment.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldComment,
			Description: AutomationActionDescAPIFieldComment,
			Value:       model.Jira.Comment.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}

	return AutomationActionTypeJira, fields, diags
}

// mapGitHubAction maps GitHub action configuration to API fields
func (r *automationActionResource) mapGitHubAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := make([]restapi.Field, 0)

	if !model.GitHub.Owner.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldOwner,
			Description: AutomationActionDescAPIFieldOwner,
			Value:       model.GitHub.Owner.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Repo.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldRepo,
			Description: AutomationActionDescAPIFieldRepo,
			Value:       model.GitHub.Repo.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Title.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldTitle,
			Description: AutomationActionDescAPIFieldTitle,
			Value:       model.GitHub.Title.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Body.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldBody,
			Description: AutomationActionDescAPIFieldGitHubBody,
			Value:       model.GitHub.Body.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Operation.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldTicketType,
			Description: AutomationActionDescAPIFieldTicketType,
			Value:       model.GitHub.Operation.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Assignees.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldAssignees,
			Description: AutomationActionDescAPIFieldAssignees,
			Value:       model.GitHub.Assignees.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Labels.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldLabels,
			Description: AutomationActionDescAPIFieldGitHubLabels,
			Value:       model.GitHub.Labels.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitHub.Comment.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldComment,
			Description: AutomationActionDescAPIFieldGitHubComment,
			Value:       model.GitHub.Comment.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}

	return AutomationActionTypeGitHub, fields, diags
}

// mapDocLinkAction maps documentation link action configuration to API fields
func (r *automationActionResource) mapDocLinkAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := []restapi.Field{{
		Name:        AutomationActionFieldUrl,
		Description: AutomationActionDescFieldUrl,
		Value:       model.DocLink.Url.ValueString(),
		Encoding:    shared.UTF8Encoding,
		Secured:     false,
	}}
	return AutomationActionTypeDocLink, fields, diags
}

// mapGitLabAction maps GitLab action configuration to API fields
func (r *automationActionResource) mapGitLabAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := make([]restapi.Field, 0)

	if !model.GitLab.ProjectId.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldProjectId,
			Description: AutomationActionDescAPIFieldProjectId,
			Value:       model.GitLab.ProjectId.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitLab.Title.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldTitle,
			Description: AutomationActionDescAPIFieldGitLabTitle,
			Value:       model.GitLab.Title.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitLab.Description.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldBody,
			Description: AutomationActionDescAPIFieldGitLabBody,
			Value:       model.GitLab.Description.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitLab.Operation.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldTicketActionType,
			Description: AutomationActionDescAPIFieldGitLabTicketActionType,
			Value:       model.GitLab.Operation.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitLab.Labels.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldLabels,
			Description: AutomationActionDescAPIFieldGitLabLabels,
			Value:       model.GitLab.Labels.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitLab.IssueType.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldIssueType,
			Description: AutomationActionDescAPIFieldGitLabIssueType,
			Value:       model.GitLab.IssueType.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.GitLab.Comment.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionFieldComment,
			Description: AutomationActionDescAPIFieldGitLabComment,
			Value:       model.GitLab.Comment.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}

	return AutomationActionTypeGitLab, fields, diags
}

// mapAnsibleAction maps Ansible action configuration to API fields
func (r *automationActionResource) mapAnsibleAction(model shared.AutomationActionModel) (string, []restapi.Field, diag.Diagnostics) {
	var diags diag.Diagnostics
	fields := make([]restapi.Field, 0)

	if !model.Ansible.WorkflowId.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldWorkflowId,
			Description: AutomationActionDescAPIFieldWorkflowId,
			Value:       model.Ansible.WorkflowId.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Ansible.AnsibleUrl.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldAnsibleUrl,
			Description: AutomationActionDescAPIFieldAnsibleUrl,
			Value:       model.Ansible.AnsibleUrl.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Ansible.HostId.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldHostId,
			Description: AutomationActionDescAPIFieldHostId,
			Value:       model.Ansible.HostId.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Ansible.PlaybookId.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldPlaybookId,
			Description: AutomationActionDescAPIFieldPlaybookId,
			Value:       model.Ansible.PlaybookId.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}
	if !model.Ansible.PlaybookFileName.IsNull() {
		fields = append(fields, restapi.Field{
			Name:        AutomationActionAPIFieldPlaybookFileName,
			Description: AutomationActionDescAPIFieldPlaybookFileName,
			Value:       model.Ansible.PlaybookFileName.ValueString(),
			Encoding:    shared.AsciiEncoding,
		})
	}

	return AutomationActionTypeAnsible, fields, diags
}

// ============================================================================
// Helper Methods
// ============================================================================

// mapInputParametersFromState maps input parameters from state to API format
func (r *automationActionResource) mapInputParametersFromState(ctx context.Context, model shared.AutomationActionModel) ([]restapi.Parameter, diag.Diagnostics) {
	return shared.MapInputParametersFromState(ctx, model)
}

// mapTagsFromState converts tags from state to API format
func (r *automationActionResource) mapTagsFromState(ctx context.Context, model shared.AutomationActionModel) (interface{}, diag.Diagnostics) {
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
