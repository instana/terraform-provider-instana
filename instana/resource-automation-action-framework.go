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

// AutomationActionModel represents the data model for the automation action resource
type AutomationActionModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Tags           types.List   `tfsdk:"tags"`
	Script         types.List   `tfsdk:"script"`
	Http           types.List   `tfsdk:"http"`
	InputParameter types.List   `tfsdk:"input_parameter"`
}

// ScriptModel represents the script configuration for an automation action
type ScriptModel struct {
	Content     types.String `tfsdk:"content"`
	Interpreter types.String `tfsdk:"interpreter"`
	Timeout     types.String `tfsdk:"timeout"`
}

// HttpModel represents the HTTP configuration for an automation action
type HttpModel struct {
	Host             types.String `tfsdk:"host"`
	Method           types.String `tfsdk:"method"`
	Body             types.String `tfsdk:"body"`
	Headers          types.Map    `tfsdk:"headers"`
	IgnoreCertErrors types.Bool   `tfsdk:"ignore_certificate_errors"`
	Timeout          types.String `tfsdk:"timeout"`
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
				},
				Blocks: map[string]schema.Block{
					AutomationActionFieldScript: schema.ListNestedBlock{
						Description: "Script configuration for the automation action.",
						NestedObject: schema.NestedBlockObject{
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
							},
						},
					},
					AutomationActionFieldHttp: schema.ListNestedBlock{
						Description: "HTTP configuration for the automation action.",
						NestedObject: schema.NestedBlockObject{
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
							},
						},
					},
					AutomationActionFieldInputParameter: schema.ListNestedBlock{
						Description: "Input parameters for the automation action.",
						NestedObject: schema.NestedBlockObject{
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

func (r *automationActionResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, automationAction *restapi.AutomationAction) diag.Diagnostics {
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

	// Handle script configuration
	if automationAction.Type == ActionTypeScript {
		scriptConfig, d := r.mapScriptFieldsToState(ctx, automationAction)
		diags.Append(d...)
		if !diags.HasError() {
			model.Script = scriptConfig
			model.Http = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					AutomationActionFieldHost:             types.StringType,
					AutomationActionFieldMethod:           types.StringType,
					AutomationActionFieldBody:             types.StringType,
					AutomationActionFieldHeaders:          types.MapType{ElemType: types.StringType},
					AutomationActionFieldIgnoreCertErrors: types.BoolType,
					AutomationActionFieldTimeout:          types.StringType,
				},
			})
		}
	} else if automationAction.Type == ActionTypeHttp {
		httpConfig, d := r.mapHttpFieldsToState(ctx, automationAction)
		diags.Append(d...)
		if !diags.HasError() {
			model.Http = httpConfig
			model.Script = types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					AutomationActionFieldContent:     types.StringType,
					AutomationActionFieldInterpreter: types.StringType,
					AutomationActionFieldTimeout:     types.StringType,
				},
			})
		}
	} else {
		model.Script = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionFieldContent:     types.StringType,
				AutomationActionFieldInterpreter: types.StringType,
				AutomationActionFieldTimeout:     types.StringType,
			},
		})
		model.Http = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionFieldHost:             types.StringType,
				AutomationActionFieldMethod:           types.StringType,
				AutomationActionFieldBody:             types.StringType,
				AutomationActionFieldHeaders:          types.MapType{ElemType: types.StringType},
				AutomationActionFieldIgnoreCertErrors: types.BoolType,
				AutomationActionFieldTimeout:          types.StringType,
			},
		})
	}

	// Set the entire model to state
	diags.Append(state.Set(ctx, model)...)
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

func (r *automationActionResourceFramework) mapScriptFieldsToState(ctx context.Context, action *restapi.AutomationAction) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	scriptObj := map[string]attr.Value{
		AutomationActionFieldContent: types.StringValue(r.getFieldValue(action, restapi.ScriptSshFieldName)),
	}

	// Add optional fields if they exist
	interpreter := r.getFieldValue(action, restapi.SubtypeFieldName)
	if interpreter != "" {
		scriptObj[AutomationActionFieldInterpreter] = types.StringValue(interpreter)
	} else {
		scriptObj[AutomationActionFieldInterpreter] = types.StringNull()
	}

	timeout := r.getFieldValue(action, restapi.TimeoutFieldName)
	if timeout != "" {
		scriptObj[AutomationActionFieldTimeout] = types.StringValue(timeout)
	} else {
		scriptObj[AutomationActionFieldTimeout] = types.StringNull()
	}

	objValue, d := types.ObjectValue(
		map[string]attr.Type{
			AutomationActionFieldContent:     types.StringType,
			AutomationActionFieldInterpreter: types.StringType,
			AutomationActionFieldTimeout:     types.StringType,
		},
		scriptObj,
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionFieldContent:     types.StringType,
				AutomationActionFieldInterpreter: types.StringType,
				AutomationActionFieldTimeout:     types.StringType,
			},
		}), diags
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionFieldContent:     types.StringType,
				AutomationActionFieldInterpreter: types.StringType,
				AutomationActionFieldTimeout:     types.StringType,
			},
		},
		[]attr.Value{objValue},
	), diags
}

func (r *automationActionResourceFramework) mapHttpFieldsToState(ctx context.Context, action *restapi.AutomationAction) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	httpObj := map[string]attr.Value{
		AutomationActionFieldHost:   types.StringValue(r.getFieldValue(action, restapi.HttpHostFieldName)),
		AutomationActionFieldMethod: types.StringValue(r.getFieldValue(action, restapi.HttpMethodFieldName)),
	}

	// Add optional fields if they exist
	body := r.getFieldValue(action, restapi.HttpBodyFieldName)
	if body != "" {
		httpObj[AutomationActionFieldBody] = types.StringValue(body)
	} else {
		httpObj[AutomationActionFieldBody] = types.StringNull()
	}

	httpObj[AutomationActionFieldIgnoreCertErrors] = types.BoolValue(
		r.getBoolFieldValueOrDefault(action, restapi.HttpIgnoreCertErrorsFieldName, false),
	)

	timeout := r.getFieldValue(action, restapi.TimeoutFieldName)
	if timeout != "" {
		httpObj[AutomationActionFieldTimeout] = types.StringValue(timeout)
	} else {
		httpObj[AutomationActionFieldTimeout] = types.StringNull()
	}

	// Handle headers
	headersData := r.getFieldValue(action, restapi.HttpHeaderFieldName)
	if headersData != "" {
		var headersMap map[string]interface{}
		err := json.Unmarshal([]byte(headersData), &headersMap)
		if err != nil {
			diags.AddError(
				"Error unmarshaling HTTP headers",
				fmt.Sprintf("Failed to unmarshal HTTP headers: %s", err),
			)
			return types.ListNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					AutomationActionFieldHost:             types.StringType,
					AutomationActionFieldMethod:           types.StringType,
					AutomationActionFieldBody:             types.StringType,
					AutomationActionFieldHeaders:          types.MapType{ElemType: types.StringType},
					AutomationActionFieldIgnoreCertErrors: types.BoolType,
					AutomationActionFieldTimeout:          types.StringType,
				},
			}), diags
		}

		elements := make(map[string]attr.Value)
		for k, v := range headersMap {
			if strVal, ok := v.(string); ok {
				elements[k] = types.StringValue(strVal)
			} else {
				elements[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
		}

		headersValue, d := types.MapValue(types.StringType, elements)
		diags.Append(d...)
		if !diags.HasError() {
			httpObj[AutomationActionFieldHeaders] = headersValue
		}
	} else {
		httpObj[AutomationActionFieldHeaders] = types.MapNull(types.StringType)
	}

	objValue, d := types.ObjectValue(
		map[string]attr.Type{
			AutomationActionFieldHost:             types.StringType,
			AutomationActionFieldMethod:           types.StringType,
			AutomationActionFieldBody:             types.StringType,
			AutomationActionFieldHeaders:          types.MapType{ElemType: types.StringType},
			AutomationActionFieldIgnoreCertErrors: types.BoolType,
			AutomationActionFieldTimeout:          types.StringType,
		},
		httpObj,
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionFieldHost:             types.StringType,
				AutomationActionFieldMethod:           types.StringType,
				AutomationActionFieldBody:             types.StringType,
				AutomationActionFieldHeaders:          types.MapType{ElemType: types.StringType},
				AutomationActionFieldIgnoreCertErrors: types.BoolType,
				AutomationActionFieldTimeout:          types.StringType,
			},
		}), diags
	}

	return types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				AutomationActionFieldHost:             types.StringType,
				AutomationActionFieldMethod:           types.StringType,
				AutomationActionFieldBody:             types.StringType,
				AutomationActionFieldHeaders:          types.MapType{ElemType: types.StringType},
				AutomationActionFieldIgnoreCertErrors: types.BoolType,
				AutomationActionFieldTimeout:          types.StringType,
			},
		},
		[]attr.Value{objValue},
	), diags
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
	if !model.Script.IsNull() {
		var scriptModels []ScriptModel
		diags.Append(model.Script.ElementsAs(ctx, &scriptModels, false)...)
		if diags.HasError() {
			return "", nil, diags
		}

		if len(scriptModels) > 0 {
			actionType = ActionTypeScript
			scriptModel := scriptModels[0]

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
		}
	} else if !model.Http.IsNull() {
		var httpModels []HttpModel
		diags.Append(model.Http.ElementsAs(ctx, &httpModels, false)...)
		if diags.HasError() {
			return "", nil, diags
		}

		if len(httpModels) > 0 {
			actionType = ActionTypeHttp
			httpModel := httpModels[0]

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
		}
	} else {
		diags.AddError(
			"Invalid action configuration",
			"Either script or http configuration must be provided",
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
