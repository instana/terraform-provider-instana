package synthetictest

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tf_framework"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceInstanaSyntheticTestFramework the name of the terraform-provider-instana resource to manage synthetic tests
const ResourceInstanaSyntheticTestFramework = "synthetic_test"

// NewSyntheticTestResourceHandleFramework creates the resource handle for Synthetic Tests
func NewSyntheticTestResourceHandleFramework() ResourceHandleFramework[*restapi.SyntheticTest] {
	return &syntheticTestResourceFramework{
		metaData: ResourceMetaDataFramework{
			ResourceName: ResourceInstanaSyntheticTestFramework,
			Schema: schema.Schema{
				Description: SyntheticTestDescResource,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: SyntheticTestDescID,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"label": schema.StringAttribute{
						Required:    true,
						Description: SyntheticTestDescLabel,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 128),
						},
					},
					"description": schema.StringAttribute{
						Optional:    true,
						Description: SyntheticTestDescDescription,
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 512),
						},
					},
					"active": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: SyntheticTestDescActive,
						Default:     booldefault.StaticBool(true),
					},
					"application_id": schema.StringAttribute{
						Optional:    true,
						Description: SyntheticTestDescApplicationID,
					},
					"custom_properties": schema.MapAttribute{
						Optional:    true,
						Description: SyntheticTestDescCustomProperties,
						ElementType: types.StringType,
					},
					"locations": schema.SetAttribute{
						Required:    true,
						Description: SyntheticTestDescLocations,
						ElementType: types.StringType,
					},
					"playback_mode": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: SyntheticTestDescPlaybackMode,
						Default:     stringdefault.StaticString("Simultaneous"),
						Validators: []validator.String{
							stringvalidator.OneOf("Simultaneous", "Staggered"),
						},
					},
					"test_frequency": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: SyntheticTestDescTestFrequency,
						Default:     int64default.StaticInt64(15),
						Validators: []validator.Int64{
							int64Validator{
								min: 1,
								max: 120,
							},
						},
					},
				},
				Blocks: map[string]schema.Block{
					"http_action": schema.ListNestedBlock{
						Description: SyntheticTestDescHttpAction,
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"mark_synthetic_call": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescMarkSyntheticCall,
									Default:     booldefault.StaticBool(false),
								},
								"retries": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescRetries,
									Default:     int64default.StaticInt64(0),
									Validators: []validator.Int64{
										int64Validator{
											min: 0,
											max: 2,
										},
									},
								},
								"retry_interval": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescRetryInterval,
									Default:     int64default.StaticInt64(1),
									Validators: []validator.Int64{
										int64Validator{
											min: 1,
											max: 10,
										},
									},
								},
								"timeout": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescTimeout,
								},
								"url": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescURL,
									Validators: []validator.String{
										stringvalidator.RegexMatches(regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`), SyntheticTestValidatorURLRegex),
									},
								},
								"operation": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescOperation,
									Validators: []validator.String{
										stringvalidator.OneOf("GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "DELETE"),
									},
								},
								"headers": schema.MapAttribute{
									Optional:    true,
									Description: SyntheticTestDescHeaders,
									ElementType: types.StringType,
								},
								"body": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescBody,
								},
								"validation_string": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescValidationString,
								},
								"follow_redirect": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescFollowRedirect,
									Default:     booldefault.StaticBool(false),
								},
								"allow_insecure": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescAllowInsecure,
									Default:     booldefault.StaticBool(false),
								},
								"expect_status": schema.Int64Attribute{
									Optional:    true,
									Description: SyntheticTestDescExpectStatus,
								},
								"expect_match": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescExpectMatch,
								},
							},
						},
					},
					"http_script": schema.ListNestedBlock{
						Description: SyntheticTestDescHttpScript,
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"mark_synthetic_call": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescMarkSyntheticCall,
									Default:     booldefault.StaticBool(false),
								},
								"retries": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescRetries,
									Default:     int64default.StaticInt64(0),
									Validators: []validator.Int64{
										int64Validator{
											min: 0,
											max: 2,
										},
									},
								},
								"retry_interval": schema.Int64Attribute{
									Optional:    true,
									Computed:    true,
									Description: SyntheticTestDescRetryInterval,
									Default:     int64default.StaticInt64(1),
									Validators: []validator.Int64{
										int64Validator{
											min: 1,
											max: 10,
										},
									},
								},
								"timeout": schema.StringAttribute{
									Optional:    true,
									Description: SyntheticTestDescTimeout,
								},
								"script": schema.StringAttribute{
									Required:    true,
									Description: SyntheticTestDescScript,
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

// int64Validator is a custom validator for int64 values
type int64Validator struct {
	min int64
	max int64
}

func (v int64Validator) Description(ctx context.Context) string {
	return fmt.Sprintf("Value must be between %d and %d", v.min, v.max)
}

func (v int64Validator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Value must be between %d and %d", v.min, v.max)
}

func (v int64Validator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		value := req.ConfigValue.ValueInt64()
		if value < v.min || value > v.max {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid value",
				fmt.Sprintf("Value must be between %d and %d, got: %d", v.min, v.max, value),
			)
		}
	}
}

// urlRegex is a regular expression to validate URLs with HTTP or HTTPS scheme
var urlRegex = `^https?://[^\s/$.?#].[^\s]*$`

type syntheticTestResourceFramework struct {
	metaData ResourceMetaDataFramework
}

func (r *syntheticTestResourceFramework) MetaData() *ResourceMetaDataFramework {
	return &r.metaData
}

func (r *syntheticTestResourceFramework) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SyntheticTest] {
	return api.SyntheticTest()
}

func (r *syntheticTestResourceFramework) SetComputedFields(_ context.Context, _ *tfsdk.Plan) diag.Diagnostics {
	return nil
}

func (r *syntheticTestResourceFramework) MapStateToDataObject(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State) (*restapi.SyntheticTest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model tf_framework.SyntheticTestModel

	// Get current state from plan or state
	if plan != nil {
		diags.Append(plan.Get(ctx, &model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	// Map application ID
	var applicationID *string
	if !model.ApplicationID.IsNull() && !model.ApplicationID.IsUnknown() {
		appID := model.ApplicationID.ValueString()
		applicationID = &appID
	}

	// Map custom properties
	customProperties := make(map[string]interface{})
	if !model.CustomProperties.IsNull() && !model.CustomProperties.IsUnknown() {
		diags.Append(model.CustomProperties.ElementsAs(ctx, &customProperties, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map locations
	var locations []string
	if !model.Locations.IsNull() && !model.Locations.IsUnknown() {
		diags.Append(model.Locations.ElementsAs(ctx, &locations, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map test frequency
	var testFrequency *int32
	if !model.TestFrequency.IsNull() && !model.TestFrequency.IsUnknown() {
		freq := int32(model.TestFrequency.ValueInt64())
		testFrequency = &freq
	}

	// Map configuration
	configuration, configDiags := r.mapConfigurationFromModel(ctx, model)
	if configDiags.HasError() {
		diags.Append(configDiags...)
		return nil, diags
	}

	// Create API object
	return &restapi.SyntheticTest{
		ID:               model.ID.ValueString(),
		Label:            model.Label.ValueString(),
		Description:      getStringPointerFromFrameworkType(model.Description),
		Active:           model.Active.ValueBool(),
		ApplicationID:    applicationID,
		Configuration:    configuration,
		CustomProperties: customProperties,
		Locations:        locations,
		PlaybackMode:     model.PlaybackMode.ValueString(),
		TestFrequency:    testFrequency,
	}, diags
}

func (r *syntheticTestResourceFramework) mapConfigurationFromModel(ctx context.Context, model tf_framework.SyntheticTestModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Check if either http_action or http_script is set
	if (model.HttpAction.IsNull() || model.HttpAction.IsUnknown()) && (model.HttpScript.IsNull() || model.HttpScript.IsUnknown()) {
		diags.AddError(SyntheticTestErrConfigRequired, SyntheticTestErrConfigRequiredMsg)
		return restapi.SyntheticTestConfig{}, diags
	}

	// Check if both http_action and http_script are set
	if (!model.HttpAction.IsNull() && !model.HttpAction.IsUnknown()) && (!model.HttpScript.IsNull() && !model.HttpScript.IsUnknown()) {
		diags.AddError(SyntheticTestErrInvalidConfig, SyntheticTestErrInvalidConfigMsg)
		return restapi.SyntheticTestConfig{}, diags
	}

	// Map HTTP Action configuration
	if !model.HttpAction.IsNull() && !model.HttpAction.IsUnknown() {
		var httpActionModels []tf_framework.HttpActionConfigModel
		diags.Append(model.HttpAction.ElementsAs(ctx, &httpActionModels, false)...)
		if diags.HasError() {
			return restapi.SyntheticTestConfig{}, diags
		}

		if len(httpActionModels) != 1 {
			diags.AddError(SyntheticTestErrInvalidHttpAction, SyntheticTestErrInvalidHttpActionMsg)
			return restapi.SyntheticTestConfig{}, diags
		}

		httpActionModel := httpActionModels[0]

		// Map headers
		var headers map[string]interface{}
		if !httpActionModel.Headers.IsNull() && !httpActionModel.Headers.IsUnknown() {
			diags.Append(httpActionModel.Headers.ElementsAs(ctx, &headers, false)...)
			if diags.HasError() {
				return restapi.SyntheticTestConfig{}, diags
			}
		}

		// Map expect status
		var expectStatus *int32
		if !httpActionModel.ExpectStatus.IsNull() && !httpActionModel.ExpectStatus.IsUnknown() {
			status := int32(httpActionModel.ExpectStatus.ValueInt64())
			expectStatus = &status
		}

		return restapi.SyntheticTestConfig{
			MarkSyntheticCall: httpActionModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(httpActionModel.Retries.ValueInt64()),
			RetryInterval:     int32(httpActionModel.RetryInterval.ValueInt64()),
			SyntheticType:     "HTTPAction",
			Timeout:           getStringPointerFromFrameworkType(httpActionModel.Timeout),
			URL:               getStringPointerFromFrameworkType(httpActionModel.URL),
			Operation:         getStringPointerFromFrameworkType(httpActionModel.Operation),
			Headers:           headers,
			Body:              getStringPointerFromFrameworkType(httpActionModel.Body),
			ValidationString:  getStringPointerFromFrameworkType(httpActionModel.ValidationString),
			FollowRedirect:    getBoolPointerFromFrameworkType(httpActionModel.FollowRedirect),
			AllowInsecure:     getBoolPointerFromFrameworkType(httpActionModel.AllowInsecure),
			ExpectStatus:      expectStatus,
			ExpectMatch:       getStringPointerFromFrameworkType(httpActionModel.ExpectMatch),
		}, diags
	}

	// Map HTTP Script configuration
	if !model.HttpScript.IsNull() && !model.HttpScript.IsUnknown() {
		var httpScriptModels []tf_framework.HttpScriptConfigModel
		diags.Append(model.HttpScript.ElementsAs(ctx, &httpScriptModels, false)...)
		if diags.HasError() {
			return restapi.SyntheticTestConfig{}, diags
		}

		if len(httpScriptModels) != 1 {
			diags.AddError(SyntheticTestErrInvalidHttpScript, SyntheticTestErrInvalidHttpScriptMsg)
			return restapi.SyntheticTestConfig{}, diags
		}

		httpScriptModel := httpScriptModels[0]

		return restapi.SyntheticTestConfig{
			MarkSyntheticCall: httpScriptModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(httpScriptModel.Retries.ValueInt64()),
			RetryInterval:     int32(httpScriptModel.RetryInterval.ValueInt64()),
			SyntheticType:     "HTTPScript",
			Timeout:           getStringPointerFromFrameworkType(httpScriptModel.Timeout),
			Script:            getStringPointerFromFrameworkType(httpScriptModel.Script),
		}, diags
	}

	// This should never happen due to the checks above
	diags.AddError(SyntheticTestErrNoValidConfig, SyntheticTestErrNoValidConfigMsg)
	return restapi.SyntheticTestConfig{}, diags
}

func (r *syntheticTestResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SyntheticTest) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API object
	model := tf_framework.SyntheticTestModel{
		ID:           types.StringValue(apiObject.ID),
		Label:        types.StringValue(apiObject.Label),
		Active:       types.BoolValue(apiObject.Active),
		PlaybackMode: types.StringValue(apiObject.PlaybackMode),
	}

	// Map description
	model.Description = util.setStringPointerToState(apiObject.Description)

	// Map application ID
	model.ApplicationID = util.setStringPointerToState(apiObject.ApplicationID)

	// Map test frequency
	if apiObject.TestFrequency != nil {
		model.TestFrequency = setInt64PointerToState(apiObject.TestFrequency)
	} else {
		model.TestFrequency = types.Int64Null()
	}

	// Map custom properties
	if apiObject.CustomProperties != nil && len(apiObject.CustomProperties) > 0 {
		customPropertiesMap := make(map[string]attr.Value)
		for k, v := range apiObject.CustomProperties {
			customPropertiesMap[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		model.CustomProperties = types.MapValueMust(types.StringType, customPropertiesMap)
	} else {
		model.CustomProperties = types.MapNull(types.StringType)
	}

	// Map locations
	if apiObject.Locations != nil && len(apiObject.Locations) > 0 {
		locationValues := make([]attr.Value, len(apiObject.Locations))
		for i, location := range apiObject.Locations {
			locationValues[i] = types.StringValue(location)
		}
		model.Locations = types.SetValueMust(types.StringType, locationValues)
	} else {
		model.Locations = types.SetNull(types.StringType)
	}

	// Map configuration based on synthetic type
	if apiObject.Configuration.SyntheticType == "HTTPAction" {
		httpActionModel := tf_framework.HttpActionConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		// Map optional fields
		httpActionModel.Timeout = util.setStringPointerToState(apiObject.Configuration.Timeout)

		httpActionModel.URL = util.setStringPointerToState(apiObject.Configuration.URL)

		httpActionModel.Operation = util.setStringPointerToState(apiObject.Configuration.Operation)

		httpActionModel.Body = util.setStringPointerToState(apiObject.Configuration.Body)

		httpActionModel.ValidationString = util.setStringPointerToState(apiObject.Configuration.ValidationString)

		if apiObject.Configuration.FollowRedirect != nil {
			httpActionModel.FollowRedirect = types.BoolValue(*apiObject.Configuration.FollowRedirect)
		}
		if apiObject.Configuration.AllowInsecure != nil {
			httpActionModel.AllowInsecure = types.BoolValue(*apiObject.Configuration.AllowInsecure)
		}
		if apiObject.Configuration.ExpectStatus != nil {
			httpActionModel.ExpectStatus = setInt64PointerToState(apiObject.Configuration.ExpectStatus)
		}
		httpActionModel.ExpectMatch = util.setStringPointerToState(apiObject.Configuration.ExpectMatch)

		// Map headers
		if apiObject.Configuration.Headers != nil && len(apiObject.Configuration.Headers) > 0 {
			headersMap := make(map[string]attr.Value)
			for k, v := range apiObject.Configuration.Headers {
				headersMap[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
			httpActionModel.Headers = types.MapValueMust(types.StringType, headersMap)
		} else {
			httpActionModel.Headers = types.MapNull(types.StringType)
		}

		// Create object for http_action
		httpActionObj, _ := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"mark_synthetic_call": types.BoolType,
			"retries":             types.Int64Type,
			"retry_interval":      types.Int64Type,
			"timeout":             types.StringType,
			"url":                 types.StringType,
			"operation":           types.StringType,
			"headers":             types.MapType{ElemType: types.StringType},
			"body":                types.StringType,
			"validation_string":   types.StringType,
			"follow_redirect":     types.BoolType,
			"allow_insecure":      types.BoolType,
			"expect_status":       types.Int64Type,
			"expect_match":        types.StringType,
		}, httpActionModel)

		model.HttpAction = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"mark_synthetic_call": types.BoolType,
					"retries":             types.Int64Type,
					"retry_interval":      types.Int64Type,
					"timeout":             types.StringType,
					"url":                 types.StringType,
					"operation":           types.StringType,
					"headers":             types.MapType{ElemType: types.StringType},
					"body":                types.StringType,
					"validation_string":   types.StringType,
					"follow_redirect":     types.BoolType,
					"allow_insecure":      types.BoolType,
					"expect_status":       types.Int64Type,
					"expect_match":        types.StringType,
				},
			},
			[]attr.Value{httpActionObj},
		)
		model.HttpScript = types.ListNull(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"mark_synthetic_call": types.BoolType,
					"retries":             types.Int64Type,
					"retry_interval":      types.Int64Type,
					"timeout":             types.StringType,
					"script":              types.StringType,
				},
			},
		)
	} else if apiObject.Configuration.SyntheticType == "HTTPScript" {
		httpScriptModel := tf_framework.HttpScriptConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		// Map optional fields
		httpScriptModel.Timeout = util.setStringPointerToState(apiObject.Configuration.Timeout)

		httpScriptModel.Script = util.setStringPointerToState(apiObject.Configuration.Script)

		// Create object for http_script
		httpScriptObj, _ := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"mark_synthetic_call": types.BoolType,
			"retries":             types.Int64Type,
			"retry_interval":      types.Int64Type,
			"timeout":             types.StringType,
			"script":              types.StringType,
		}, httpScriptModel)

		model.HttpScript = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"mark_synthetic_call": types.BoolType,
					"retries":             types.Int64Type,
					"retry_interval":      types.Int64Type,
					"timeout":             types.StringType,
					"script":              types.StringType,
				},
			},
			[]attr.Value{httpScriptObj},
		)
		model.HttpAction = types.ListNull(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"mark_synthetic_call": types.BoolType,
					"retries":             types.Int64Type,
					"retry_interval":      types.Int64Type,
					"timeout":             types.StringType,
					"url":                 types.StringType,
					"operation":           types.StringType,
					"headers":             types.MapType{ElemType: types.StringType},
					"body":                types.StringType,
					"validation_string":   types.StringType,
					"follow_redirect":     types.BoolType,
					"allow_insecure":      types.BoolType,
					"expect_status":       types.Int64Type,
					"expect_match":        types.StringType,
				},
			},
		)
	}

	// Set state
	diags.Append(state.Set(ctx, &model)...)
	return diags
}

// Helper functions
func getStringPointerFromFrameworkType(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	str := value.ValueString()
	return &str
}

func getBoolPointerFromFrameworkType(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	b := value.ValueBool()
	return &b
}

// Made with Bob
