package synthetictest

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/util"
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
	var applications []string
	if !model.Applications.IsNull() && !model.Applications.IsUnknown() {
		diags.Append(model.Applications.ElementsAs(ctx, &applications, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map mobile apps
	var mobileApps []string
	if !model.MobileApps.IsNull() && !model.MobileApps.IsUnknown() {
		diags.Append(model.MobileApps.ElementsAs(ctx, &mobileApps, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Map websites
	var websites []string
	if !model.Websites.IsNull() && !model.Websites.IsUnknown() {
		diags.Append(model.Websites.ElementsAs(ctx, &websites, false)...)
		if diags.HasError() {
			return nil, diags
		}
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

	// Map RBAC tags
	var rbacTags []restapi.ApiTag
	if !model.RbacTags.IsNull() && !model.RbacTags.IsUnknown() {
		var rbacTagModels []RbacTagModel
		diags.Append(model.RbacTags.ElementsAs(ctx, &rbacTagModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, tagModel := range rbacTagModels {
			rbacTags = append(rbacTags, restapi.ApiTag{
				Name:  tagModel.Name.ValueString(),
				Value: tagModel.Value.ValueString(),
			})
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

		// Map expect exists
		var expectExists []string
		if !httpActionModel.ExpectExists.IsNull() && !httpActionModel.ExpectExists.IsUnknown() {
			diags.Append(httpActionModel.ExpectExists.ElementsAs(ctx, &expectExists, false)...)
			if diags.HasError() {
				return restapi.SyntheticTestConfig{}, diags
			}
		}

		// Map expect not empty
		var expectNotEmpty []string
		if !httpActionModel.ExpectNotEmpty.IsNull() && !httpActionModel.ExpectNotEmpty.IsUnknown() {
			diags.Append(httpActionModel.ExpectNotEmpty.ElementsAs(ctx, &expectNotEmpty, false)...)
			if diags.HasError() {
				return restapi.SyntheticTestConfig{}, diags
			}
		}

		// Map expect json
		var expectJson map[string]interface{}
		if !httpActionModel.ExpectJson.IsNull() && !httpActionModel.ExpectJson.IsUnknown() {
			diags.Append(httpActionModel.ExpectJson.ElementsAs(ctx, &expectJson, false)...)
			if diags.HasError() {
				return restapi.SyntheticTestConfig{}, diags
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
		if !dnsModel.TargetValues.IsNull() && !dnsModel.TargetValues.IsUnknown() {
			var targetValueModels []DNSFilterTargetValueModel
			diags.Append(dnsModel.TargetValues.ElementsAs(ctx, &targetValueModels, false)...)
			if !diags.HasError() {
				for _, tvModel := range targetValueModels {
					config.TargetValues = append(config.TargetValues, restapi.DNSFilterTargetValue{
						Key:      tvModel.Key.ValueString(),
						Operator: tvModel.Operator.ValueString(),
						Value:    tvModel.Value.ValueString(),
					})
				}
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
		if !sslModel.ValidationRules.IsNull() && !sslModel.ValidationRules.IsUnknown() {
			var validationRuleModels []SSLCertificateValidationModel
			diags.Append(sslModel.ValidationRules.ElementsAs(ctx, &validationRuleModels, false)...)
			if !diags.HasError() {
				for _, vrModel := range validationRuleModels {
					config.ValidationRules = append(config.ValidationRules, restapi.SSLCertificateValidation{
						Key:      vrModel.Key.ValueString(),
						Operator: vrModel.Operator.ValueString(),
						Value:    vrModel.Value.ValueString(),
					})
				}
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
		appValues := make([]attr.Value, len(apiObject.Applications))
		for i, app := range apiObject.Applications {
			appValues[i] = types.StringValue(app)
		}
		model.Applications = types.SetValueMust(types.StringType, appValues)
	} else {
		model.Applications = types.SetNull(types.StringType)
	}

	// Map mobile apps
	if apiObject.MobileApps != nil && len(apiObject.MobileApps) > 0 {
		mobileAppValues := make([]attr.Value, len(apiObject.MobileApps))
		for i, app := range apiObject.MobileApps {
			mobileAppValues[i] = types.StringValue(app)
		}
		model.MobileApps = types.SetValueMust(types.StringType, mobileAppValues)
	} else {
		model.MobileApps = types.SetNull(types.StringType)
	}

	// Map websites
	if apiObject.Websites != nil && len(apiObject.Websites) > 0 {
		websiteValues := make([]attr.Value, len(apiObject.Websites))
		for i, website := range apiObject.Websites {
			websiteValues[i] = types.StringValue(website)
		}
		model.Websites = types.SetValueMust(types.StringType, websiteValues)
	} else {
		model.Websites = types.SetNull(types.StringType)
	}

	// Map test frequency
	if apiObject.TestFrequency != nil {
		model.TestFrequency = util.SetInt64PointerToState(apiObject.TestFrequency)
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

	// Map RBAC tags
	if apiObject.RbacTags != nil && len(apiObject.RbacTags) > 0 {
		rbacTagValues := make([]attr.Value, len(apiObject.RbacTags))
		for i, tag := range apiObject.RbacTags {
			tagObj, _ := types.ObjectValue(
				map[string]attr.Type{
					"name":  types.StringType,
					"value": types.StringType,
				},
				map[string]attr.Value{
					"name":  types.StringValue(tag.Name),
					"value": types.StringValue(tag.Value),
				},
			)
			rbacTagValues[i] = tagObj
		}
		model.RbacTags = types.SetValueMust(
			types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}},
			rbacTagValues,
		)
	} else {
		model.RbacTags = types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		}})
	}

	// Map configuration based on synthetic type
	if apiObject.Configuration.SyntheticType == "HTTPAction" {
		httpActionModel := &HttpActionConfigModel{
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
			headersMap := make(map[string]attr.Value)
			for k, v := range apiObject.Configuration.Headers {
				headersMap[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
			httpActionModel.Headers = types.MapValueMust(types.StringType, headersMap)
		} else {
			httpActionModel.Headers = types.MapNull(types.StringType)
		}

		// Map expect exists
		if apiObject.Configuration.ExpectExists != nil && len(apiObject.Configuration.ExpectExists) > 0 {
			expectExistsValues := make([]attr.Value, len(apiObject.Configuration.ExpectExists))
			for i, val := range apiObject.Configuration.ExpectExists {
				expectExistsValues[i] = types.StringValue(val)
			}
			httpActionModel.ExpectExists = types.SetValueMust(types.StringType, expectExistsValues)
		} else {
			httpActionModel.ExpectExists = types.SetNull(types.StringType)
		}

		// Map expect not empty
		if apiObject.Configuration.ExpectNotEmpty != nil && len(apiObject.Configuration.ExpectNotEmpty) > 0 {
			expectNotEmptyValues := make([]attr.Value, len(apiObject.Configuration.ExpectNotEmpty))
			for i, val := range apiObject.Configuration.ExpectNotEmpty {
				expectNotEmptyValues[i] = types.StringValue(val)
			}
			httpActionModel.ExpectNotEmpty = types.SetValueMust(types.StringType, expectNotEmptyValues)
		} else {
			httpActionModel.ExpectNotEmpty = types.SetNull(types.StringType)
		}

		// Map expect json
		if apiObject.Configuration.ExpectJson != nil && len(apiObject.Configuration.ExpectJson) > 0 {
			expectJsonMap := make(map[string]attr.Value)
			for k, v := range apiObject.Configuration.ExpectJson {
				expectJsonMap[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
			httpActionModel.ExpectJson = types.MapValueMust(types.StringType, expectJsonMap)
		} else {
			httpActionModel.ExpectJson = types.MapNull(types.StringType)
		}

		model.HttpAction = httpActionModel
		// Set all other config types to nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil
	} else if apiObject.Configuration.SyntheticType == "HTTPScript" {
		httpScriptModel := &HttpScriptConfigModel{
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
		}

		model.HttpScript = httpScriptModel
		model.HttpAction = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "BrowserScript" {
		browserScriptModel := &BrowserScriptConfigModel{
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
		}

		model.BrowserScript = browserScriptModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "DNS" {
		dnsModel := &DNSConfigModel{
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
		}

		// Map target values
		if apiObject.Configuration.TargetValues != nil && len(apiObject.Configuration.TargetValues) > 0 {
			targetValueObjs := make([]attr.Value, len(apiObject.Configuration.TargetValues))
			for i, tv := range apiObject.Configuration.TargetValues {
				tvObj, _ := types.ObjectValue(
					map[string]attr.Type{
						"key":      types.StringType,
						"operator": types.StringType,
						"value":    types.StringType,
					},
					map[string]attr.Value{
						"key":      types.StringValue(tv.Key),
						"operator": types.StringValue(tv.Operator),
						"value":    types.StringValue(tv.Value),
					},
				)
				targetValueObjs[i] = tvObj
			}
			dnsModel.TargetValues = types.SetValueMust(
				types.ObjectType{AttrTypes: map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				}},
				targetValueObjs,
			)
		} else {
			dnsModel.TargetValues = types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"key":      types.StringType,
				"operator": types.StringType,
				"value":    types.StringType,
			}})
		}

		model.DNS = dnsModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.SSLCertificate = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "SSLCertificate" {
		sslModel := &SSLCertificateConfigModel{
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
			validationRuleObjs := make([]attr.Value, len(apiObject.Configuration.ValidationRules))
			for i, vr := range apiObject.Configuration.ValidationRules {
				vrObj, _ := types.ObjectValue(
					map[string]attr.Type{
						"key":      types.StringType,
						"operator": types.StringType,
						"value":    types.StringType,
					},
					map[string]attr.Value{
						"key":      types.StringValue(vr.Key),
						"operator": types.StringValue(vr.Operator),
						"value":    types.StringValue(fmt.Sprintf("%v", vr.Value)),
					},
				)
				validationRuleObjs[i] = vrObj
			}
			sslModel.ValidationRules = types.SetValueMust(
				types.ObjectType{AttrTypes: map[string]attr.Type{
					"key":      types.StringType,
					"operator": types.StringType,
					"value":    types.StringType,
				}},
				validationRuleObjs,
			)
		} else {
			sslModel.ValidationRules = types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
				"key":      types.StringType,
				"operator": types.StringType,
				"value":    types.StringType,
			}})
		}

		model.SSLCertificate = sslModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.WebpageAction = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "WebpageAction" {
		webpageActionModel := &WebpageActionConfigModel{
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

		model.WebpageAction = webpageActionModel
		model.HttpAction = nil
		model.HttpScript = nil
		model.BrowserScript = nil
		model.DNS = nil
		model.SSLCertificate = nil
		model.WebpageScript = nil

	} else if apiObject.Configuration.SyntheticType == "WebpageScript" {
		webpageScriptModel := &WebpageScriptConfigModel{
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

		model.WebpageScript = webpageScriptModel
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
