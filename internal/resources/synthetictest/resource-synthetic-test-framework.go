package synthetictest

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

// NewSyntheticTestResourceHandleFramework creates the resource handle for Synthetic Tests
func NewSyntheticTestResourceHandleFramework() resourcehandle.ResourceHandleFramework[*restapi.SyntheticTest] {
	return &syntheticTestResourceFramework{
		metaData: resourcehandle.ResourceMetaDataFramework{
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
					"applications": schema.SetAttribute{
						Optional:    true,
						Description: "Array of application IDs",
						ElementType: types.StringType,
					},
					"mobile_apps": schema.SetAttribute{
						Optional:    true,
						Description: "Array of mobile app IDs",
						ElementType: types.StringType,
					},
					"websites": schema.SetAttribute{
						Optional:    true,
						Description: "Array of website IDs",
						ElementType: types.StringType,
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
					"rbac_tags": schema.SetNestedAttribute{
						Optional:    true,
						Description: "RBAC tags for access control",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Tag name",
								},
								"value": schema.StringAttribute{
									Required:    true,
									Description: "Tag value",
								},
							},
						},
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
					"http_action": schema.SingleNestedAttribute{
						Optional:    true,
						Description: SyntheticTestDescHttpAction,
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
							"expect_exists": schema.SetAttribute{
								Optional:    true,
								Description: "JSON paths that must exist in the response",
								ElementType: types.StringType,
							},
							"expect_not_empty": schema.SetAttribute{
								Optional:    true,
								Description: "JSON paths that must not be empty in the response",
								ElementType: types.StringType,
							},
							"expect_json": schema.MapAttribute{
								Optional:    true,
								Description: "Expected JSON structure in the response",
								ElementType: types.StringType,
							},
						},
					},
					"http_script": schema.SingleNestedAttribute{
						Optional:    true,
						Description: SyntheticTestDescHttpScript,
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
								Optional:    true,
								Description: SyntheticTestDescScript,
							},
							"script_type": schema.StringAttribute{
								Optional:    true,
								Description: "Script type (Basic or Jest)",
								Validators: []validator.String{
									stringvalidator.OneOf("Basic", "Jest"),
								},
							},
							"file_name": schema.StringAttribute{
								Optional:    true,
								Description: "Script file name",
							},
							"scripts": schema.SingleNestedAttribute{
								Optional:    true,
								Description: "Multiple scripts configuration for Jest",
								Attributes: map[string]schema.Attribute{
									"bundle": schema.StringAttribute{
										Optional:    true,
										Description: "Bundle content",
									},
									"script_file": schema.StringAttribute{
										Optional:    true,
										Description: "Script file content",
									},
								},
							},
						},
					},
					"browser_script": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Browser script configuration",
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
								Optional:    true,
								Description: SyntheticTestDescScript,
							},
							"script_type": schema.StringAttribute{
								Optional:    true,
								Description: "Script type (Basic or Jest)",
								Validators: []validator.String{
									stringvalidator.OneOf("Basic", "Jest"),
								},
							},
							"file_name": schema.StringAttribute{
								Optional:    true,
								Description: "Script file name",
							},
							"scripts": schema.SingleNestedAttribute{
								Optional:    true,
								Description: "Multiple scripts configuration for Jest",
								Attributes: map[string]schema.Attribute{
									"bundle": schema.StringAttribute{
										Optional:    true,
										Description: "Bundle content",
									},
									"script_file": schema.StringAttribute{
										Optional:    true,
										Description: "Script file content",
									},
								},
							},
							"browser": schema.StringAttribute{
								Optional:    true,
								Description: "Browser type (chrome or firefox)",
								Validators: []validator.String{
									stringvalidator.OneOf("chrome", "firefox"),
								},
							},
							"record_video": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Record video of the test execution",
								Default:     booldefault.StaticBool(false),
							},
						},
					},
					"dns": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "DNS test configuration",
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
							"lookup": schema.StringAttribute{
								Required:    true,
								Description: "Domain name to lookup",
							},
							"server": schema.StringAttribute{
								Required:    true,
								Description: "DNS server to query",
							},
							"query_type": schema.StringAttribute{
								Optional:    true,
								Description: "DNS query type",
								Validators: []validator.String{
									stringvalidator.OneOf("ALL", "ALL_CONDITIONS", "ANY", "A", "AAAA", "CNAME", "NS"),
								},
							},
							"port": schema.Int64Attribute{
								Optional:    true,
								Description: "DNS server port",
							},
							"transport": schema.StringAttribute{
								Optional:    true,
								Description: "Transport protocol (TCP or UDP)",
								Validators: []validator.String{
									stringvalidator.OneOf("TCP", "UDP"),
								},
							},
							"accept_cname": schema.BoolAttribute{
								Optional:    true,
								Description: "Accept CNAME records",
							},
							"lookup_server_name": schema.BoolAttribute{
								Optional:    true,
								Description: "Lookup server name",
							},
							"recursive_lookups": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable recursive lookups",
							},
							"server_retries": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of server retries",
							},
							"query_time": schema.SingleNestedAttribute{
								Optional:    true,
								Description: "Query time filter",
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required:    true,
										Description: "Filter key",
									},
									"operator": schema.StringAttribute{
										Required:    true,
										Description: "Filter operator",
										Validators: []validator.String{
											stringvalidator.OneOf("CONTAINS", "EQUALS", "GREATER_THAN", "IS", "LESS_THAN", "MATCHES", "NOT_MATCHES"),
										},
									},
									"value": schema.Int64Attribute{
										Required:    true,
										Description: "Filter value",
									},
								},
							},
							"target_values": schema.SetNestedAttribute{
								Optional:    true,
								Description: "Target value filters",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Required:    true,
											Description: "Filter key",
											Validators: []validator.String{
												stringvalidator.OneOf("ALL", "ALL_CONDITIONS", "ANY", "A", "AAAA", "CNAME", "NS"),
											},
										},
										"operator": schema.StringAttribute{
											Required:    true,
											Description: "Filter operator",
											Validators: []validator.String{
												stringvalidator.OneOf("CONTAINS", "EQUALS", "GREATER_THAN", "IS", "LESS_THAN", "MATCHES", "NOT_MATCHES"),
											},
										},
										"value": schema.StringAttribute{
											Required:    true,
											Description: "Filter value",
										},
									},
								},
							},
						},
					},
					"ssl_certificate": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "SSL certificate test configuration",
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
							"hostname": schema.StringAttribute{
								Required:    true,
								Description: "Hostname to check SSL certificate",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 2047),
								},
							},
							"days_remaining_check": schema.Int64Attribute{
								Required:    true,
								Description: "Minimum days remaining before certificate expiration",
								Validators: []validator.Int64{
									int64Validator{
										min: 1,
										max: 365,
									},
								},
							},
							"accept_self_signed_certificate": schema.BoolAttribute{
								Optional:    true,
								Description: "Accept self-signed certificates",
							},
							"port": schema.Int64Attribute{
								Optional:    true,
								Description: "Port number",
							},
							"validation_rules": schema.SetNestedAttribute{
								Optional:    true,
								Description: "SSL certificate validation rules",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Required:    true,
											Description: "Validation key",
										},
										"operator": schema.StringAttribute{
											Required:    true,
											Description: "Validation operator",
											Validators: []validator.String{
												stringvalidator.OneOf("CONTAINS", "EQUALS", "GREATER_THAN", "IS", "LESS_THAN", "MATCHES", "NOT_MATCHES"),
											},
										},
										"value": schema.StringAttribute{
											Required:    true,
											Description: "Validation value",
										},
									},
								},
							},
						},
					},
					"webpage_action": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Webpage action test configuration",
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
								Required:    true,
								Description: SyntheticTestDescURL,
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`), SyntheticTestValidatorURLRegex),
								},
							},
							"browser": schema.StringAttribute{
								Optional:    true,
								Description: "Browser type (chrome or firefox)",
								Validators: []validator.String{
									stringvalidator.OneOf("chrome", "firefox"),
								},
							},
							"record_video": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Record video of the test execution",
								Default:     booldefault.StaticBool(false),
							},
						},
					},
					"webpage_script": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Webpage script test configuration",
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
							"file_name": schema.StringAttribute{
								Optional:    true,
								Description: "Script file name",
							},
							"browser": schema.StringAttribute{
								Optional:    true,
								Description: "Browser type (chrome or firefox)",
								Validators: []validator.String{
									stringvalidator.OneOf("chrome", "firefox"),
								},
							},
							"record_video": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Record video of the test execution",
								Default:     booldefault.StaticBool(false),
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
	metaData resourcehandle.ResourceMetaDataFramework
}

func (r *syntheticTestResourceFramework) MetaData() *resourcehandle.ResourceMetaDataFramework {
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
	var model SyntheticTestModel

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

	// Map applications
	applications := model.Applications

	// Map mobile apps
	mobileApps := model.MobileApps

	// Map websites
	websites := model.Websites

	// Map custom properties
	customProperties := make(map[string]interface{})
	for k, v := range model.CustomProperties {
		customProperties[k] = v
	}

	// Map locations
	locations := model.Locations

	// Map RBAC tags
	var rbacTags []restapi.ApiTag
	for _, tagModel := range model.RbacTags {
		rbacTags = append(rbacTags, restapi.ApiTag{
			Name:  tagModel.Name.ValueString(),
			Value: tagModel.Value.ValueString(),
		})
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
		Applications:     applications,
		MobileApps:       mobileApps,
		Websites:         websites,
		Configuration:    configuration,
		CustomProperties: customProperties,
		Locations:        locations,
		PlaybackMode:     model.PlaybackMode.ValueString(),
		TestFrequency:    testFrequency,
		RbacTags:         rbacTags,
	}, diags
}

func (r *syntheticTestResourceFramework) mapConfigurationFromModel(ctx context.Context, model SyntheticTestModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Count how many configuration types are set
	configCount := 0
	if model.HttpAction != nil {
		configCount++
	}
	if model.HttpScript != nil {
		configCount++
	}
	if model.BrowserScript != nil {
		configCount++
	}
	if model.DNS != nil {
		configCount++
	}
	if model.SSLCertificate != nil {
		configCount++
	}
	if model.WebpageAction != nil {
		configCount++
	}
	if model.WebpageScript != nil {
		configCount++
	}

	// Validate exactly one configuration type is set
	if configCount == 0 {
		diags.AddError(SyntheticTestErrConfigRequired, "Exactly one synthetic test configuration type must be specified")
		return restapi.SyntheticTestConfig{}, diags
	}
	if configCount > 1 {
		diags.AddError(SyntheticTestErrInvalidConfig, "Only one synthetic test configuration type can be specified")
		return restapi.SyntheticTestConfig{}, diags
	}

	// Map HTTP Action configuration
	if model.HttpAction != nil {
		httpActionModel := model.HttpAction

		// Map headers
		var headers map[string]interface{}
		if httpActionModel.Headers != nil && len(httpActionModel.Headers) > 0 {
			headers = make(map[string]interface{}, len(httpActionModel.Headers))
			for k, v := range httpActionModel.Headers {
				headers[k] = v
			}
		}

		// Map expect status
		var expectStatus *int32
		if !httpActionModel.ExpectStatus.IsNull() && !httpActionModel.ExpectStatus.IsUnknown() {
			status := int32(httpActionModel.ExpectStatus.ValueInt64())
			expectStatus = &status
		}

		// Map expect exists
		var expectExists []string
		if len(httpActionModel.ExpectExists) > 0 {
			expectExists = httpActionModel.ExpectExists
		}

		// Map expect not empty
		var expectNotEmpty []string
		if len(httpActionModel.ExpectNotEmpty) > 0 {
			expectNotEmpty = httpActionModel.ExpectNotEmpty
		}

		// Map expect json
		var expectJson map[string]interface{}
		if len(httpActionModel.ExpectJson) > 0 {
			expectJson = make(map[string]interface{}, len(httpActionModel.ExpectJson))
			for k, v := range httpActionModel.ExpectJson {
				expectJson[k] = v
			}
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
			ExpectExists:      expectExists,
			ExpectNotEmpty:    expectNotEmpty,
			ExpectJson:        expectJson,
		}, diags
	}

	// Map HTTP Script configuration
	if model.HttpScript != nil {
		httpScriptModel := model.HttpScript

		config := restapi.SyntheticTestConfig{
			MarkSyntheticCall: httpScriptModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(httpScriptModel.Retries.ValueInt64()),
			RetryInterval:     int32(httpScriptModel.RetryInterval.ValueInt64()),
			SyntheticType:     "HTTPScript",
			Timeout:           getStringPointerFromFrameworkType(httpScriptModel.Timeout),
			Script:            getStringPointerFromFrameworkType(httpScriptModel.Script),
			ScriptType:        getStringPointerFromFrameworkType(httpScriptModel.ScriptType),
			FileName:          getStringPointerFromFrameworkType(httpScriptModel.FileName),
		}

		// Map scripts if present
		if httpScriptModel.Scripts != nil {
			config.Scripts = &restapi.MultipleScriptsConfiguration{
				Bundle:     getStringPointerFromFrameworkType(httpScriptModel.Scripts.Bundle),
				ScriptFile: getStringPointerFromFrameworkType(httpScriptModel.Scripts.ScriptFile),
			}
		}

		return config, diags
	}

	// Map Browser Script configuration
	if model.BrowserScript != nil {
		browserScriptModel := model.BrowserScript

		config := restapi.SyntheticTestConfig{
			MarkSyntheticCall: browserScriptModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(browserScriptModel.Retries.ValueInt64()),
			RetryInterval:     int32(browserScriptModel.RetryInterval.ValueInt64()),
			SyntheticType:     "BrowserScript",
			Timeout:           getStringPointerFromFrameworkType(browserScriptModel.Timeout),
			Script:            getStringPointerFromFrameworkType(browserScriptModel.Script),
			ScriptType:        getStringPointerFromFrameworkType(browserScriptModel.ScriptType),
			FileName:          getStringPointerFromFrameworkType(browserScriptModel.FileName),
			Browser:           getStringPointerFromFrameworkType(browserScriptModel.Browser),
			RecordVideo:       getBoolPointerFromFrameworkType(browserScriptModel.RecordVideo),
		}

		// Map scripts if present
		if browserScriptModel.Scripts != nil {
			config.Scripts = &restapi.MultipleScriptsConfiguration{
				Bundle:     getStringPointerFromFrameworkType(browserScriptModel.Scripts.Bundle),
				ScriptFile: getStringPointerFromFrameworkType(browserScriptModel.Scripts.ScriptFile),
			}
		}

		return config, diags
	}

	// Map DNS configuration
	if model.DNS != nil {
		dnsModel := model.DNS

		config := restapi.SyntheticTestConfig{
			MarkSyntheticCall: dnsModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(dnsModel.Retries.ValueInt64()),
			RetryInterval:     int32(dnsModel.RetryInterval.ValueInt64()),
			SyntheticType:     "DNS",
			Timeout:           getStringPointerFromFrameworkType(dnsModel.Timeout),
			Lookup:            getStringPointerFromFrameworkType(dnsModel.Lookup),
			Server:            getStringPointerFromFrameworkType(dnsModel.Server),
			QueryType:         getStringPointerFromFrameworkType(dnsModel.QueryType),
			Transport:         getStringPointerFromFrameworkType(dnsModel.Transport),
			AcceptCNAME:       getBoolPointerFromFrameworkType(dnsModel.AcceptCNAME),
			LookupServerName:  getBoolPointerFromFrameworkType(dnsModel.LookupServerName),
			RecursiveLookups:  getBoolPointerFromFrameworkType(dnsModel.RecursiveLookups),
		}

		// Map port
		if !dnsModel.Port.IsNull() && !dnsModel.Port.IsUnknown() {
			port := int32(dnsModel.Port.ValueInt64())
			config.Port = &port
		}

		// Map server retries
		if !dnsModel.ServerRetries.IsNull() && !dnsModel.ServerRetries.IsUnknown() {
			retries := int32(dnsModel.ServerRetries.ValueInt64())
			config.ServerRetries = &retries
		}

		// Map query time
		if dnsModel.QueryTime != nil {
			config.QueryTime = &restapi.DNSFilterQueryTime{
				Key:      dnsModel.QueryTime.Key.ValueString(),
				Operator: dnsModel.QueryTime.Operator.ValueString(),
				Value:    dnsModel.QueryTime.Value.ValueInt64(),
			}
		}

		// Map target values
		if len(dnsModel.TargetValues) > 0 {
			for _, tvModel := range dnsModel.TargetValues {
				config.TargetValues = append(config.TargetValues, restapi.DNSFilterTargetValue{
					Key:      tvModel.Key.ValueString(),
					Operator: tvModel.Operator.ValueString(),
					Value:    tvModel.Value.ValueString(),
				})
			}
		}

		return config, diags
	}

	// Map SSL Certificate configuration
	if model.SSLCertificate != nil {
		sslModel := model.SSLCertificate

		config := restapi.SyntheticTestConfig{
			MarkSyntheticCall:    sslModel.MarkSyntheticCall.ValueBool(),
			Retries:              int32(sslModel.Retries.ValueInt64()),
			RetryInterval:        int32(sslModel.RetryInterval.ValueInt64()),
			SyntheticType:        "SSLCertificate",
			Timeout:              getStringPointerFromFrameworkType(sslModel.Timeout),
			Hostname:             getStringPointerFromFrameworkType(sslModel.Hostname),
			AcceptSelfSignedCert: getBoolPointerFromFrameworkType(sslModel.AcceptSelfSignedCert),
		}

		// Map days remaining check
		if !sslModel.DaysRemainingCheck.IsNull() && !sslModel.DaysRemainingCheck.IsUnknown() {
			days := int32(sslModel.DaysRemainingCheck.ValueInt64())
			config.DaysRemainingCheck = &days
		}

		// Map port
		if !sslModel.Port.IsNull() && !sslModel.Port.IsUnknown() {
			port := int32(sslModel.Port.ValueInt64())
			config.SSLPort = &port
		}

		// Map validation rules
		if len(sslModel.ValidationRules) > 0 {
			for _, vrModel := range sslModel.ValidationRules {
				config.ValidationRules = append(config.ValidationRules, restapi.SSLCertificateValidation{
					Key:      vrModel.Key.ValueString(),
					Operator: vrModel.Operator.ValueString(),
					Value:    vrModel.Value.ValueString(),
				})
			}
		}

		return config, diags
	}

	// Map Webpage Action configuration
	if model.WebpageAction != nil {
		webpageActionModel := model.WebpageAction

		return restapi.SyntheticTestConfig{
			MarkSyntheticCall: webpageActionModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(webpageActionModel.Retries.ValueInt64()),
			RetryInterval:     int32(webpageActionModel.RetryInterval.ValueInt64()),
			SyntheticType:     "WebpageAction",
			Timeout:           getStringPointerFromFrameworkType(webpageActionModel.Timeout),
			URL:               getStringPointerFromFrameworkType(webpageActionModel.URL),
			Browser:           getStringPointerFromFrameworkType(webpageActionModel.Browser),
			RecordVideo:       getBoolPointerFromFrameworkType(webpageActionModel.RecordVideo),
		}, diags
	}

	// Map Webpage Script configuration
	if model.WebpageScript != nil {
		webpageScriptModel := model.WebpageScript

		return restapi.SyntheticTestConfig{
			MarkSyntheticCall: webpageScriptModel.MarkSyntheticCall.ValueBool(),
			Retries:           int32(webpageScriptModel.Retries.ValueInt64()),
			RetryInterval:     int32(webpageScriptModel.RetryInterval.ValueInt64()),
			SyntheticType:     "WebpageScript",
			Timeout:           getStringPointerFromFrameworkType(webpageScriptModel.Timeout),
			Script:            getStringPointerFromFrameworkType(webpageScriptModel.Script),
			FileName:          getStringPointerFromFrameworkType(webpageScriptModel.FileName),
			Browser:           getStringPointerFromFrameworkType(webpageScriptModel.Browser),
			RecordVideo:       getBoolPointerFromFrameworkType(webpageScriptModel.RecordVideo),
		}, diags
	}

	// This should never happen due to the checks above
	diags.AddError(SyntheticTestErrNoValidConfig, "No valid synthetic test configuration found")
	return restapi.SyntheticTestConfig{}, diags
}

func (r *syntheticTestResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SyntheticTest) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create a model and populate it with values from the API object
	model := SyntheticTestModel{
		ID:           types.StringValue(apiObject.ID),
		Label:        types.StringValue(apiObject.Label),
		Active:       types.BoolValue(apiObject.Active),
		PlaybackMode: types.StringValue(apiObject.PlaybackMode),
	}

	// Map description
	model.Description = util.SetStringPointerToState(apiObject.Description)

	// Map application ID
	model.ApplicationID = util.SetStringPointerToState(apiObject.ApplicationID)

	// Map applications
	if apiObject.Applications != nil && len(apiObject.Applications) > 0 {
		model.Applications = apiObject.Applications
	} else {
		model.Applications = nil
	}

	// Map mobile apps
	if apiObject.MobileApps != nil && len(apiObject.MobileApps) > 0 {
		model.MobileApps = apiObject.MobileApps
	} else {
		model.MobileApps = nil
	}

	// Map websites
	if apiObject.Websites != nil && len(apiObject.Websites) > 0 {
		model.Websites = apiObject.Websites
	} else {
		model.Websites = nil
	}

	// Map test frequency
	if apiObject.TestFrequency != nil {
		model.TestFrequency = util.SetInt64PointerToState(apiObject.TestFrequency)
	} else {
		model.TestFrequency = types.Int64Null()
	}

	// Map custom properties
	if apiObject.CustomProperties != nil && len(apiObject.CustomProperties) > 0 {
		customPropertiesMap := make(map[string]string)
		for k, v := range apiObject.CustomProperties {
			customPropertiesMap[k] = fmt.Sprintf("%v", v)
		}
		model.CustomProperties = customPropertiesMap
	} else {
		model.CustomProperties = nil
	}

	// Map locations
	if apiObject.Locations != nil && len(apiObject.Locations) > 0 {
		model.Locations = apiObject.Locations
	} else {
		model.Locations = nil
	}

	// Map RBAC tags
	if apiObject.RbacTags != nil && len(apiObject.RbacTags) > 0 {
		rbacTags := make([]RbacTagModel, len(apiObject.RbacTags))
		for i, tag := range apiObject.RbacTags {
			rbacTags[i] = RbacTagModel{
				Name:  types.StringValue(tag.Name),
				Value: types.StringValue(tag.Value),
			}
		}
		model.RbacTags = rbacTags
	} else {
		model.RbacTags = nil
	}

	// Map configuration based on synthetic type
	if apiObject.Configuration.SyntheticType == "HTTPAction" {
		httpActionModel := HttpActionConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		// Map optional fields
		httpActionModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		httpActionModel.URL = util.SetStringPointerToState(apiObject.Configuration.URL)
		httpActionModel.Operation = util.SetStringPointerToState(apiObject.Configuration.Operation)
		httpActionModel.Body = util.SetStringPointerToState(apiObject.Configuration.Body)
		httpActionModel.ValidationString = util.SetStringPointerToState(apiObject.Configuration.ValidationString)

		if apiObject.Configuration.FollowRedirect != nil {
			httpActionModel.FollowRedirect = types.BoolValue(*apiObject.Configuration.FollowRedirect)
		}
		if apiObject.Configuration.AllowInsecure != nil {
			httpActionModel.AllowInsecure = types.BoolValue(*apiObject.Configuration.AllowInsecure)
		}
		if apiObject.Configuration.ExpectStatus != nil {
			httpActionModel.ExpectStatus = util.SetInt64PointerToState(apiObject.Configuration.ExpectStatus)
		}
		httpActionModel.ExpectMatch = util.SetStringPointerToState(apiObject.Configuration.ExpectMatch)

		// Map headers
		if apiObject.Configuration.Headers != nil && len(apiObject.Configuration.Headers) > 0 {
			headersMap := make(map[string]string)
			for k, v := range apiObject.Configuration.Headers {
				headersMap[k] = fmt.Sprintf("%v", v)
			}
			httpActionModel.Headers = headersMap
		} else {
			httpActionModel.Headers = nil
		}

		// Map expect exists
		if apiObject.Configuration.ExpectExists != nil && len(apiObject.Configuration.ExpectExists) > 0 {
			httpActionModel.ExpectExists = apiObject.Configuration.ExpectExists
		} else {
			httpActionModel.ExpectExists = nil
		}

		// Map expect not empty
		if apiObject.Configuration.ExpectNotEmpty != nil && len(apiObject.Configuration.ExpectNotEmpty) > 0 {
			httpActionModel.ExpectNotEmpty = apiObject.Configuration.ExpectNotEmpty
		} else {
			httpActionModel.ExpectNotEmpty = nil
		}

		// Map expect json
		if apiObject.Configuration.ExpectJson != nil && len(apiObject.Configuration.ExpectJson) > 0 {
			expectJsonMap := make(map[string]string)
			for k, v := range apiObject.Configuration.ExpectJson {
				expectJsonMap[k] = fmt.Sprintf("%v", v)
			}
			httpActionModel.ExpectJson = expectJsonMap
		} else {
			httpActionModel.ExpectJson = nil
		}

		// Set http_action and null out other config types
		model.HttpAction = &httpActionModel
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil
	} else if apiObject.Configuration.SyntheticType == "HTTPScript" {
		httpScriptModel := HttpScriptConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		// Map optional fields
		httpScriptModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		httpScriptModel.Script = util.SetStringPointerToState(apiObject.Configuration.Script)
		httpScriptModel.ScriptType = util.SetStringPointerToState(apiObject.Configuration.ScriptType)
		if apiObject.Configuration.FileName != nil && *apiObject.Configuration.FileName != "" {
			httpScriptModel.FileName = types.StringValue(*apiObject.Configuration.FileName)
		} else {
			httpScriptModel.FileName = types.StringNull()
		}

		// Map scripts if present
		if apiObject.Configuration.Scripts != nil {
			httpScriptModel.Scripts = &MultipleScriptsModel{
				Bundle:     util.SetStringPointerToState(apiObject.Configuration.Scripts.Bundle),
				ScriptFile: util.SetStringPointerToState(apiObject.Configuration.Scripts.ScriptFile),
			}
		} else {
			httpScriptModel.Scripts = nil
		}

		// Set http_script and null out other config types
		model.HttpScript = &httpScriptModel
		model.HttpAction = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "BrowserScript" {
		browserScriptModel := BrowserScriptConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		browserScriptModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		browserScriptModel.Script = util.SetStringPointerToState(apiObject.Configuration.Script)
		browserScriptModel.ScriptType = util.SetStringPointerToState(apiObject.Configuration.ScriptType)
		if apiObject.Configuration.FileName != nil && *apiObject.Configuration.FileName != "" {
			browserScriptModel.FileName = types.StringValue(*apiObject.Configuration.FileName)
		} else {
			browserScriptModel.FileName = types.StringNull()
		}
		browserScriptModel.Browser = util.SetStringPointerToState(apiObject.Configuration.Browser)

		if apiObject.Configuration.RecordVideo != nil {
			browserScriptModel.RecordVideo = types.BoolValue(*apiObject.Configuration.RecordVideo)
		}

		if apiObject.Configuration.Scripts != nil {
			browserScriptModel.Scripts = &MultipleScriptsModel{
				Bundle:     util.SetStringPointerToState(apiObject.Configuration.Scripts.Bundle),
				ScriptFile: util.SetStringPointerToState(apiObject.Configuration.Scripts.ScriptFile),
			}
		} else {
			browserScriptModel.Scripts = nil
		}

		// Set browser_script and null out other config types
		model.BrowserScript = &browserScriptModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "DNS" {
		dnsModel := DNSConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		dnsModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		dnsModel.Lookup = util.SetStringPointerToState(apiObject.Configuration.Lookup)
		dnsModel.Server = util.SetStringPointerToState(apiObject.Configuration.Server)
		dnsModel.QueryType = util.SetStringPointerToState(apiObject.Configuration.QueryType)
		dnsModel.Transport = util.SetStringPointerToState(apiObject.Configuration.Transport)

		if apiObject.Configuration.Port != nil {
			dnsModel.Port = util.SetInt64PointerToState(apiObject.Configuration.Port)
		}
		if apiObject.Configuration.ServerRetries != nil {
			dnsModel.ServerRetries = util.SetInt64PointerToState(apiObject.Configuration.ServerRetries)
		}
		if apiObject.Configuration.AcceptCNAME != nil {
			dnsModel.AcceptCNAME = types.BoolValue(*apiObject.Configuration.AcceptCNAME)
		}
		if apiObject.Configuration.LookupServerName != nil {
			dnsModel.LookupServerName = types.BoolValue(*apiObject.Configuration.LookupServerName)
		}
		if apiObject.Configuration.RecursiveLookups != nil {
			dnsModel.RecursiveLookups = types.BoolValue(*apiObject.Configuration.RecursiveLookups)
		}

		// Map query time
		if apiObject.Configuration.QueryTime != nil {
			dnsModel.QueryTime = &DNSFilterQueryTimeModel{
				Key:      types.StringValue(apiObject.Configuration.QueryTime.Key),
				Operator: types.StringValue(apiObject.Configuration.QueryTime.Operator),
				Value:    types.Int64Value(apiObject.Configuration.QueryTime.Value),
			}
		} else {
			dnsModel.QueryTime = nil
		}

		// Map target values
		if apiObject.Configuration.TargetValues != nil && len(apiObject.Configuration.TargetValues) > 0 {
			targetValues := make([]DNSFilterTargetValueModel, len(apiObject.Configuration.TargetValues))
			for i, tv := range apiObject.Configuration.TargetValues {
				targetValues[i] = DNSFilterTargetValueModel{
					Key:      types.StringValue(tv.Key),
					Operator: types.StringValue(tv.Operator),
					Value:    types.StringValue(tv.Value),
				}
			}
			dnsModel.TargetValues = targetValues
		} else {
			dnsModel.TargetValues = nil
		}

		// Set dns and null out other config types
		model.DNS = &dnsModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "SSLCertificate" {
		sslModel := SSLCertificateConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		sslModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		sslModel.Hostname = util.SetStringPointerToState(apiObject.Configuration.Hostname)

		if apiObject.Configuration.DaysRemainingCheck != nil {
			sslModel.DaysRemainingCheck = util.SetInt64PointerToState(apiObject.Configuration.DaysRemainingCheck)
		}
		if apiObject.Configuration.SSLPort != nil {
			sslModel.Port = util.SetInt64PointerToState(apiObject.Configuration.SSLPort)
		}
		if apiObject.Configuration.AcceptSelfSignedCert != nil {
			sslModel.AcceptSelfSignedCert = types.BoolValue(*apiObject.Configuration.AcceptSelfSignedCert)
		}

		// Map validation rules
		if apiObject.Configuration.ValidationRules != nil && len(apiObject.Configuration.ValidationRules) > 0 {
			validationRules := make([]SSLCertificateValidationModel, len(apiObject.Configuration.ValidationRules))
			for i, vr := range apiObject.Configuration.ValidationRules {
				validationRules[i] = SSLCertificateValidationModel{
					Key:      types.StringValue(vr.Key),
					Operator: types.StringValue(vr.Operator),
					Value:    types.StringValue(fmt.Sprintf("%v", vr.Value)),
				}
			}
			sslModel.ValidationRules = validationRules
		} else {
			sslModel.ValidationRules = nil
		}

		// Set ssl and null out other config types
		model.SSLCertificate = &sslModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "WebpageAction" {
		webpageActionModel := WebpageActionConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		webpageActionModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		webpageActionModel.URL = util.SetStringPointerToState(apiObject.Configuration.URL)
		webpageActionModel.Browser = util.SetStringPointerToState(apiObject.Configuration.Browser)

		if apiObject.Configuration.RecordVideo != nil {
			webpageActionModel.RecordVideo = types.BoolValue(*apiObject.Configuration.RecordVideo)
		}

		// Set webpage action and null out other config types
		model.WebpageAction = &webpageActionModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "WebpageScript" {
		webpageScriptModel := WebpageScriptConfigModel{
			MarkSyntheticCall: types.BoolValue(apiObject.Configuration.MarkSyntheticCall),
			Retries:           types.Int64Value(int64(apiObject.Configuration.Retries)),
			RetryInterval:     types.Int64Value(int64(apiObject.Configuration.RetryInterval)),
		}

		webpageScriptModel.Timeout = util.SetStringPointerToState(apiObject.Configuration.Timeout)
		webpageScriptModel.Script = util.SetStringPointerToState(apiObject.Configuration.Script)
		if apiObject.Configuration.FileName != nil && *apiObject.Configuration.FileName != "" {
			webpageScriptModel.FileName = types.StringValue(*apiObject.Configuration.FileName)
		} else {
			webpageScriptModel.FileName = types.StringNull()
		}
		webpageScriptModel.Browser = util.SetStringPointerToState(apiObject.Configuration.Browser)

		if apiObject.Configuration.RecordVideo != nil {
			webpageScriptModel.RecordVideo = types.BoolValue(*apiObject.Configuration.RecordVideo)
		}

		// Set webpage script and null out other config types
		model.WebpageScript = &webpageScriptModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
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
