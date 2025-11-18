package synthetictest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/gessnerfl/terraform-provider-instana/internal/resourcehandle"
	"github.com/gessnerfl/terraform-provider-instana/internal/restapi"
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
			ResourceName:  ResourceInstanaSyntheticTestFramework,
			Schema:        buildSchema(),
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
// Schema builder methods

// buildCommonConfigAttributes builds common configuration attributes shared across all test types
func buildCommonConfigAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticTestFieldMarkSyntheticCall: schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescMarkSyntheticCall,
			Default:     booldefault.StaticBool(SyntheticTestDefaultMarkSynthetic),
		},
		SyntheticTestFieldRetries: schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescRetries,
			Default:     int64default.StaticInt64(SyntheticTestDefaultRetries),
			Validators: []validator.Int64{
				int64Validator{
					min: SyntheticTestMinRetries,
					max: SyntheticTestMaxRetries,
				},
			},
		},
		SyntheticTestFieldRetryInterval: schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescRetryInterval,
			Default:     int64default.StaticInt64(SyntheticTestDefaultRetryInterval),
			Validators: []validator.Int64{
				int64Validator{
					min: SyntheticTestMinRetryInterval,
					max: SyntheticTestMaxRetryInterval,
				},
			},
		},
		SyntheticTestFieldTimeout: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescTimeout,
		},
	}
}

// buildScriptAttributes builds script-related attributes
func buildScriptAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticTestFieldScript: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescScript,
		},
		SyntheticTestFieldScriptType: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescScriptType,
			Validators: []validator.String{
				stringvalidator.OneOf(SyntheticTestScriptTypeBasic, SyntheticTestScriptTypeJest),
			},
		},
		SyntheticTestFieldFileName: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescFileName,
		},
		SyntheticTestFieldScripts: schema.SingleNestedAttribute{
			Optional:    true,
			Description: SyntheticTestDescScripts,
			Attributes: map[string]schema.Attribute{
				SyntheticTestFieldBundle: schema.StringAttribute{
					Optional:    true,
					Description: SyntheticTestDescBundle,
				},
				SyntheticTestFieldScriptFile: schema.StringAttribute{
					Optional:    true,
					Description: SyntheticTestDescScriptFile,
				},
			},
		},
	}
}

// buildBrowserAttributes builds browser-related attributes
func buildBrowserAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticTestFieldBrowser: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescBrowser,
			Validators: []validator.String{
				stringvalidator.OneOf(SyntheticTestBrowserChrome, SyntheticTestBrowserFirefox),
			},
		},
		SyntheticTestFieldRecordVideo: schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescRecordVideo,
			Default:     booldefault.StaticBool(SyntheticTestDefaultRecordVideo),
		},
	}
}

// buildHttpActionSchema builds HTTP Action configuration schema
func buildHttpActionSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add HTTP Action specific attributes
	attrs[SyntheticTestFieldURL] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescURL,
		Validators: []validator.String{
			stringvalidator.RegexMatches(regexp.MustCompile(urlRegex), SyntheticTestValidatorURLRegex),
		},
	}
	attrs[SyntheticTestFieldOperation] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescOperation,
		Validators: []validator.String{
			stringvalidator.OneOf(
				SyntheticTestOperationGET,
				SyntheticTestOperationHEAD,
				SyntheticTestOperationOPTIONS,
				SyntheticTestOperationPATCH,
				SyntheticTestOperationPOST,
				SyntheticTestOperationPUT,
				SyntheticTestOperationDELETE,
			),
		},
	}
	attrs[SyntheticTestFieldHeaders] = schema.MapAttribute{
		Optional:    true,
		Description: SyntheticTestDescHeaders,
		ElementType: types.StringType,
	}
	attrs[SyntheticTestFieldBody] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescBody,
	}
	attrs[SyntheticTestFieldValidationString] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescValidationString,
	}
	attrs[SyntheticTestFieldFollowRedirect] = schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: SyntheticTestDescFollowRedirect,
		Default:     booldefault.StaticBool(SyntheticTestDefaultFollowRedirect),
	}
	attrs[SyntheticTestFieldAllowInsecure] = schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: SyntheticTestDescAllowInsecure,
		Default:     booldefault.StaticBool(SyntheticTestDefaultAllowInsecure),
	}
	attrs[SyntheticTestFieldExpectStatus] = schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticTestDescExpectStatus,
	}
	attrs[SyntheticTestFieldExpectMatch] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescExpectMatch,
	}
	attrs[SyntheticTestFieldExpectExists] = schema.SetAttribute{
		Optional:    true,
		Description: SyntheticTestDescExpectExists,
		ElementType: types.StringType,
	}
	attrs[SyntheticTestFieldExpectNotEmpty] = schema.SetAttribute{
		Optional:    true,
		Description: SyntheticTestDescExpectNotEmpty,
		ElementType: types.StringType,
	}
	attrs[SyntheticTestFieldExpectJson] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescExpectJson,
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescHttpAction,
		Attributes:  attrs,
	}
}

// buildHttpScriptSchema builds HTTP Script configuration schema
func buildHttpScriptSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add script attributes
	for k, v := range buildScriptAttributes() {
		attrs[k] = v
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescHttpScript,
		Attributes:  attrs,
	}
}

// buildBrowserScriptSchema builds Browser Script configuration schema
func buildBrowserScriptSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add script attributes
	for k, v := range buildScriptAttributes() {
		attrs[k] = v
	}

	// Add browser attributes
	for k, v := range buildBrowserAttributes() {
		attrs[k] = v
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescBrowserScript,
		Attributes:  attrs,
	}
}

// buildDNSSchema builds DNS configuration schema
func buildDNSSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add DNS specific attributes
	attrs[SyntheticTestFieldLookup] = schema.StringAttribute{
		Required:    true,
		Description: SyntheticTestDescLookup,
	}
	attrs[SyntheticTestFieldServer] = schema.StringAttribute{
		Required:    true,
		Description: SyntheticTestDescServer,
	}
	attrs[SyntheticTestFieldQueryType] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescQueryType,
		Validators: []validator.String{
			stringvalidator.OneOf(
				SyntheticTestDNSQueryTypeALL,
				SyntheticTestDNSQueryTypeALL_CONDITIONS,
				SyntheticTestDNSQueryTypeANY,
				SyntheticTestDNSQueryTypeA,
				SyntheticTestDNSQueryTypeAAAA,
				SyntheticTestDNSQueryTypeCNAME,
				SyntheticTestDNSQueryTypeNS,
			),
		},
	}
	attrs[SyntheticTestFieldPort] = schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticTestDescPortNumber,
	}
	attrs[SyntheticTestFieldTransport] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescTransport,
		Validators: []validator.String{
			stringvalidator.OneOf(SyntheticTestTransportTCP, SyntheticTestTransportUDP),
		},
	}
	attrs[SyntheticTestFieldAcceptCNAME] = schema.BoolAttribute{
		Optional:    true,
		Description: SyntheticTestDescAcceptCNAME,
	}
	attrs[SyntheticTestFieldLookupServerName] = schema.BoolAttribute{
		Optional:    true,
		Description: SyntheticTestDescLookupServerName,
	}
	attrs[SyntheticTestFieldRecursiveLookups] = schema.BoolAttribute{
		Optional:    true,
		Description: SyntheticTestDescRecursiveLookups,
	}
	attrs[SyntheticTestFieldServerRetries] = schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticTestDescServerRetries,
	}
	attrs[SyntheticTestFieldQueryTime] = schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescQueryTime,
		Attributes: map[string]schema.Attribute{
			SyntheticTestFieldKey: schema.StringAttribute{
				Required:    true,
				Description: SyntheticTestDescFilterKey,
			},
			SyntheticTestFieldOperator: schema.StringAttribute{
				Required:    true,
				Description: SyntheticTestDescFilterOperator,
				Validators: []validator.String{
					stringvalidator.OneOf(
						SyntheticTestOperatorCONTAINS,
						SyntheticTestOperatorEQUALS,
						SyntheticTestOperatorGREATER_THAN,
						SyntheticTestOperatorIS,
						SyntheticTestOperatorLESS_THAN,
						SyntheticTestOperatorMATCHES,
						SyntheticTestOperatorNOT_MATCHES,
					),
				},
			},
			SyntheticTestFieldValue: schema.Int64Attribute{
				Required:    true,
				Description: SyntheticTestDescFilterValue,
			},
		},
	}
	attrs[SyntheticTestFieldTargetValues] = schema.SetNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescTargetValues,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				SyntheticTestFieldKey: schema.StringAttribute{
					Required:    true,
					Description: SyntheticTestDescFilterKey,
					Validators: []validator.String{
						stringvalidator.OneOf(
							SyntheticTestDNSQueryTypeALL,
							SyntheticTestDNSQueryTypeALL_CONDITIONS,
							SyntheticTestDNSQueryTypeANY,
							SyntheticTestDNSQueryTypeA,
							SyntheticTestDNSQueryTypeAAAA,
							SyntheticTestDNSQueryTypeCNAME,
							SyntheticTestDNSQueryTypeNS,
						),
					},
				},
				SyntheticTestFieldOperator: schema.StringAttribute{
					Required:    true,
					Description: SyntheticTestDescFilterOperator,
					Validators: []validator.String{
						stringvalidator.OneOf(
							SyntheticTestOperatorCONTAINS,
							SyntheticTestOperatorEQUALS,
							SyntheticTestOperatorGREATER_THAN,
							SyntheticTestOperatorIS,
							SyntheticTestOperatorLESS_THAN,
							SyntheticTestOperatorMATCHES,
							SyntheticTestOperatorNOT_MATCHES,
						),
					},
				},
				SyntheticTestFieldValue: schema.StringAttribute{
					Required:    true,
					Description: SyntheticTestDescFilterValue,
				},
			},
		},
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescDNS,
		Attributes:  attrs,
	}
}

// buildSSLCertificateSchema builds SSL Certificate configuration schema
func buildSSLCertificateSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add SSL Certificate specific attributes
	attrs[SyntheticTestFieldHostname] = schema.StringAttribute{
		Required:    true,
		Description: SyntheticTestDescHostname,
		Validators: []validator.String{
			stringvalidator.LengthBetween(SyntheticTestMinHostnameLength, SyntheticTestMaxHostnameLength),
		},
	}
	attrs[SyntheticTestFieldDaysRemainingCheck] = schema.Int64Attribute{
		Required:    true,
		Description: SyntheticTestDescDaysRemainingCheck,
		Validators: []validator.Int64{
			int64Validator{
				min: SyntheticTestMinDaysRemainingCheck,
				max: SyntheticTestMaxDaysRemainingCheck,
			},
		},
	}
	attrs[SyntheticTestFieldAcceptSelfSignedCertificate] = schema.BoolAttribute{
		Optional:    true,
		Description: SyntheticTestDescAcceptSelfSignedCert,
	}
	attrs[SyntheticTestFieldPort] = schema.Int64Attribute{
		Optional:    true,
		Description: SyntheticTestDescPortNumber,
	}
	attrs[SyntheticTestFieldValidationRules] = schema.SetNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescValidationRules,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				SyntheticTestFieldKey: schema.StringAttribute{
					Required:    true,
					Description: SyntheticTestDescValidationKey,
				},
				SyntheticTestFieldOperator: schema.StringAttribute{
					Required:    true,
					Description: SyntheticTestDescValidationOperator,
					Validators: []validator.String{
						stringvalidator.OneOf(
							SyntheticTestOperatorCONTAINS,
							SyntheticTestOperatorEQUALS,
							SyntheticTestOperatorGREATER_THAN,
							SyntheticTestOperatorIS,
							SyntheticTestOperatorLESS_THAN,
							SyntheticTestOperatorMATCHES,
							SyntheticTestOperatorNOT_MATCHES,
						),
					},
				},
				SyntheticTestFieldValue: schema.StringAttribute{
					Required:    true,
					Description: SyntheticTestDescValidationValue,
				},
			},
		},
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescSSLCertificate,
		Attributes:  attrs,
	}
}

// buildWebpageActionSchema builds Webpage Action configuration schema
func buildWebpageActionSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add URL attribute
	attrs[SyntheticTestFieldURL] = schema.StringAttribute{
		Required:    true,
		Description: SyntheticTestDescURL,
		Validators: []validator.String{
			stringvalidator.RegexMatches(regexp.MustCompile(urlRegex), SyntheticTestValidatorURLRegex),
		},
	}

	// Add browser attributes
	for k, v := range buildBrowserAttributes() {
		attrs[k] = v
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescWebpageAction,
		Attributes:  attrs,
	}
}

// buildWebpageScriptSchema builds Webpage Script configuration schema
func buildWebpageScriptSchema() schema.SingleNestedAttribute {
	attrs := buildCommonConfigAttributes()

	// Add script attribute (required for webpage script)
	attrs[SyntheticTestFieldScript] = schema.StringAttribute{
		Required:    true,
		Description: SyntheticTestDescScript,
	}
	attrs[SyntheticTestFieldFileName] = schema.StringAttribute{
		Optional:    true,
		Description: SyntheticTestDescFileName,
	}

	// Add browser attributes
	for k, v := range buildBrowserAttributes() {
		attrs[k] = v
	}

	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: SyntheticTestDescWebpageScript,
		Attributes:  attrs,
	}
}

// buildBaseSchema builds the base schema attributes
func buildBaseSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		SyntheticTestFieldID: schema.StringAttribute{
			Computed:    true,
			Description: SyntheticTestDescID,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		SyntheticTestFieldLabel: schema.StringAttribute{
			Required:    true,
			Description: SyntheticTestDescLabel,
			Validators: []validator.String{
				stringvalidator.LengthBetween(SyntheticTestMinLabelLength, SyntheticTestMaxLabelLength),
			},
		},
		SyntheticTestFieldDescription: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescDescription,
			Validators: []validator.String{
				stringvalidator.LengthBetween(SyntheticTestMinDescriptionLength, SyntheticTestMaxDescriptionLength),
			},
		},
		SyntheticTestFieldActive: schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescActive,
			Default:     booldefault.StaticBool(SyntheticTestDefaultActive),
		},
		SyntheticTestFieldApplicationID: schema.StringAttribute{
			Optional:    true,
			Description: SyntheticTestDescApplicationID,
		},
		SyntheticTestFieldApplications: schema.SetAttribute{
			Optional:    true,
			Description: SyntheticTestDescApplications,
			ElementType: types.StringType,
		},
		SyntheticTestFieldMobileApps: schema.SetAttribute{
			Optional:    true,
			Description: SyntheticTestDescMobileApps,
			ElementType: types.StringType,
		},
		SyntheticTestFieldWebsites: schema.SetAttribute{
			Optional:    true,
			Description: SyntheticTestDescWebsites,
			ElementType: types.StringType,
		},
		SyntheticTestFieldCustomProperties: schema.MapAttribute{
			Optional:    true,
			Description: SyntheticTestDescCustomProperties,
			ElementType: types.StringType,
		},
		SyntheticTestFieldLocations: schema.SetAttribute{
			Required:    true,
			Description: SyntheticTestDescLocations,
			ElementType: types.StringType,
		},
		SyntheticTestFieldRbacTags: schema.SetNestedAttribute{
			Optional:    true,
			Description: SyntheticTestDescRbacTags,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					SyntheticTestFieldName: schema.StringAttribute{
						Required:    true,
						Description: SyntheticTestDescTagName,
					},
					SyntheticTestFieldValue: schema.StringAttribute{
						Required:    true,
						Description: SyntheticTestDescTagValue,
					},
				},
			},
		},
		SyntheticTestFieldPlaybackMode: schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescPlaybackMode,
			Default:     stringdefault.StaticString(SyntheticTestPlaybackModeSimultaneous),
			Validators: []validator.String{
				stringvalidator.OneOf(SyntheticTestPlaybackModeSimultaneous, SyntheticTestPlaybackModeStaggered),
			},
		},
		SyntheticTestFieldTestFrequency: schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: SyntheticTestDescTestFrequency,
			Default:     int64default.StaticInt64(SyntheticTestDefaultTestFrequency),
			Validators: []validator.Int64{
				int64Validator{
					min: SyntheticTestMinTestFrequency,
					max: SyntheticTestMaxTestFrequency,
				},
			},
		},
	}
}

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
	diags.Append(r.getModelFromState(ctx, plan, state, &model)...)
	if diags.HasError() {
		return nil, diags
	}

	// Map application ID
	applicationID := r.mapApplicationIDFromModel(model)

	// Map collection fields
	applications, appDiags := r.mapApplicationsFromModel(ctx, model)
	diags.Append(appDiags...)

	mobileApps, mobileDiags := r.mapMobileAppsFromModel(ctx, model)
	diags.Append(mobileDiags...)

	websites, webDiags := r.mapWebsitesFromModel(ctx, model)
	diags.Append(webDiags...)

	customProperties, propDiags := r.mapCustomPropertiesFromModel(ctx, model)
	diags.Append(propDiags...)

	locations, locDiags := r.mapLocationsFromModel(ctx, model)
	diags.Append(locDiags...)

	rbacTags, tagDiags := r.mapRbacTagsFromModel(ctx, model)
	diags.Append(tagDiags...)

	if diags.HasError() {
		return nil, diags
	}

	// Map test frequency
	testFrequency := r.mapTestFrequencyFromModel(model)

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

// MapStateToDataObject helper methods

// getModelFromState retrieves the model from plan or state
func (r *syntheticTestResourceFramework) getModelFromState(ctx context.Context, plan *tfsdk.Plan, state *tfsdk.State, model *SyntheticTestModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if plan != nil {
		diags.Append(plan.Get(ctx, model)...)
	} else if state != nil {
		diags.Append(state.Get(ctx, model)...)
	}
	return diags
}

// mapApplicationIDFromModel maps application ID from model
func (r *syntheticTestResourceFramework) mapApplicationIDFromModel(model SyntheticTestModel) *string {
	if !model.ApplicationID.IsNull() && !model.ApplicationID.IsUnknown() {
		appID := model.ApplicationID.ValueString()
		return &appID
	}
	return nil
}

// mapApplicationsFromModel maps applications from model
func (r *syntheticTestResourceFramework) mapApplicationsFromModel(ctx context.Context, model SyntheticTestModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var applications []string
	if !model.Applications.IsNull() && !model.Applications.IsUnknown() {
		diags.Append(model.Applications.ElementsAs(ctx, &applications, false)...)
	}
	return applications, diags
}

// mapMobileAppsFromModel maps mobile apps from model
func (r *syntheticTestResourceFramework) mapMobileAppsFromModel(ctx context.Context, model SyntheticTestModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var mobileApps []string
	if !model.MobileApps.IsNull() && !model.MobileApps.IsUnknown() {
		diags.Append(model.MobileApps.ElementsAs(ctx, &mobileApps, false)...)
	}
	return mobileApps, diags
}

// mapWebsitesFromModel maps websites from model
func (r *syntheticTestResourceFramework) mapWebsitesFromModel(ctx context.Context, model SyntheticTestModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var websites []string
	if !model.Websites.IsNull() && !model.Websites.IsUnknown() {
		diags.Append(model.Websites.ElementsAs(ctx, &websites, false)...)
	}
	return websites, diags
}

// mapCustomPropertiesFromModel maps custom properties from model
func (r *syntheticTestResourceFramework) mapCustomPropertiesFromModel(ctx context.Context, model SyntheticTestModel) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	customProperties := make(map[string]string)
	if !model.CustomProperties.IsNull() && !model.CustomProperties.IsUnknown() {
		diags.Append(model.CustomProperties.ElementsAs(ctx, &customProperties, false)...)
	}
	return customProperties, diags
}

// mapLocationsFromModel maps locations from model
func (r *syntheticTestResourceFramework) mapLocationsFromModel(ctx context.Context, model SyntheticTestModel) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var locations []string
	if !model.Locations.IsNull() && !model.Locations.IsUnknown() {
		diags.Append(model.Locations.ElementsAs(ctx, &locations, false)...)
	}
	return locations, diags
}

// mapRbacTagsFromModel maps RBAC tags from model
func (r *syntheticTestResourceFramework) mapRbacTagsFromModel(ctx context.Context, model SyntheticTestModel) ([]restapi.ApiTag, diag.Diagnostics) {
	var diags diag.Diagnostics
	var rbacTags []restapi.ApiTag
	if !model.RbacTags.IsNull() && !model.RbacTags.IsUnknown() {
		var rbacTagModels []RbacTagModel
		diags.Append(model.RbacTags.ElementsAs(ctx, &rbacTagModels, false)...)
		if !diags.HasError() {
			for _, tagModel := range rbacTagModels {
				rbacTags = append(rbacTags, restapi.ApiTag{
					Name:  tagModel.Name.ValueString(),
					Value: tagModel.Value.ValueString(),
				})
			}
		}
	}
	return rbacTags, diags
}

// mapTestFrequencyFromModel maps test frequency from model
func (r *syntheticTestResourceFramework) mapTestFrequencyFromModel(model SyntheticTestModel) *int32 {
	if !model.TestFrequency.IsNull() && !model.TestFrequency.IsUnknown() {
		freq := int32(model.TestFrequency.ValueInt64())
		return &freq
	}
	return nil
}

// MapStateToDataObject helper methods for mapping model to API objects

// validateSingleConfigType validates that exactly one configuration type is set
func (r *syntheticTestResourceFramework) validateSingleConfigType(model SyntheticTestModel) (int, diag.Diagnostics) {
	var diags diag.Diagnostics
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

	if configCount == 0 {
		diags.AddError(SyntheticTestErrConfigRequired, "Exactly one synthetic test configuration type must be specified")
	} else if configCount > 1 {
		diags.AddError(SyntheticTestErrInvalidConfig, "Only one synthetic test configuration type can be specified")
	}

	return configCount, diags
}

// mapHttpActionFromModel maps HTTP Action model to API configuration
func (r *syntheticTestResourceFramework) mapHttpActionFromModel(ctx context.Context, httpActionModel *HttpActionConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := restapi.SyntheticTestConfig{
		MarkSyntheticCall: httpActionModel.MarkSyntheticCall.ValueBool(),
		Retries:           int32(httpActionModel.Retries.ValueInt64()),
		RetryInterval:     int32(httpActionModel.RetryInterval.ValueInt64()),
		SyntheticType:     SyntheticTestTypeHTTPAction,
		Timeout:           getStringPointerFromFrameworkType(httpActionModel.Timeout),
		URL:               getStringPointerFromFrameworkType(httpActionModel.URL),
		Operation:         getStringPointerFromFrameworkType(httpActionModel.Operation),
		Body:              getStringPointerFromFrameworkType(httpActionModel.Body),
		ValidationString:  getStringPointerFromFrameworkType(httpActionModel.ValidationString),
		FollowRedirect:    getBoolPointerFromFrameworkType(httpActionModel.FollowRedirect),
		AllowInsecure:     getBoolPointerFromFrameworkType(httpActionModel.AllowInsecure),
		ExpectMatch:       getStringPointerFromFrameworkType(httpActionModel.ExpectMatch),
	}

	// Map expect status
	if !httpActionModel.ExpectStatus.IsNull() && !httpActionModel.ExpectStatus.IsUnknown() {
		status := int32(httpActionModel.ExpectStatus.ValueInt64())
		config.ExpectStatus = &status
	}

	// Map headers
	if !httpActionModel.Headers.IsNull() && !httpActionModel.Headers.IsUnknown() {
		var headers map[string]string
		diags.Append(httpActionModel.Headers.ElementsAs(ctx, &headers, false)...)
		if !diags.HasError() {
			config.Headers = headers
		}
	}

	// Map expect exists
	if !httpActionModel.ExpectExists.IsNull() && !httpActionModel.ExpectExists.IsUnknown() {
		var expectExists []string
		diags.Append(httpActionModel.ExpectExists.ElementsAs(ctx, &expectExists, false)...)
		if !diags.HasError() {
			config.ExpectExists = expectExists
		}
	}

	// Map expect not empty
	if !httpActionModel.ExpectNotEmpty.IsNull() && !httpActionModel.ExpectNotEmpty.IsUnknown() {
		var expectNotEmpty []string
		diags.Append(httpActionModel.ExpectNotEmpty.ElementsAs(ctx, &expectNotEmpty, false)...)
		if !diags.HasError() {
			config.ExpectNotEmpty = expectNotEmpty
		}
	}

	// Map expect json
	if !httpActionModel.ExpectJson.IsNull() && !httpActionModel.ExpectJson.IsUnknown() {
		normalizedWidgets := util.NormalizeJSONString(httpActionModel.ExpectJson.ValueString())
		widgets := json.RawMessage(normalizedWidgets)
		config.ExpectJson = widgets
	}

	return config, diags
}

// mapHttpScriptFromModel maps HTTP Script model to API configuration
func (r *syntheticTestResourceFramework) mapHttpScriptFromModel(httpScriptModel *HttpScriptConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := restapi.SyntheticTestConfig{
		MarkSyntheticCall: httpScriptModel.MarkSyntheticCall.ValueBool(),
		Retries:           int32(httpScriptModel.Retries.ValueInt64()),
		RetryInterval:     int32(httpScriptModel.RetryInterval.ValueInt64()),
		SyntheticType:     SyntheticTestTypeHTTPScript,
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

// mapBrowserScriptFromModel maps Browser Script model to API configuration
func (r *syntheticTestResourceFramework) mapBrowserScriptFromModel(browserScriptModel *BrowserScriptConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := restapi.SyntheticTestConfig{
		MarkSyntheticCall: browserScriptModel.MarkSyntheticCall.ValueBool(),
		Retries:           int32(browserScriptModel.Retries.ValueInt64()),
		RetryInterval:     int32(browserScriptModel.RetryInterval.ValueInt64()),
		SyntheticType:     SyntheticTestTypeBrowserScript,
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

// mapDNSFromModel maps DNS model to API configuration
func (r *syntheticTestResourceFramework) mapDNSFromModel(ctx context.Context, dnsModel *DNSConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := restapi.SyntheticTestConfig{
		MarkSyntheticCall: dnsModel.MarkSyntheticCall.ValueBool(),
		Retries:           int32(dnsModel.Retries.ValueInt64()),
		RetryInterval:     int32(dnsModel.RetryInterval.ValueInt64()),
		SyntheticType:     SyntheticTestTypeDNS,
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

// mapSSLCertificateFromModel maps SSL Certificate model to API configuration
func (r *syntheticTestResourceFramework) mapSSLCertificateFromModel(ctx context.Context, sslModel *SSLCertificateConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := restapi.SyntheticTestConfig{
		MarkSyntheticCall:    sslModel.MarkSyntheticCall.ValueBool(),
		Retries:              int32(sslModel.Retries.ValueInt64()),
		RetryInterval:        int32(sslModel.RetryInterval.ValueInt64()),
		SyntheticType:        SyntheticTestTypeSSLCertificate,
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

// mapWebpageActionFromModel maps Webpage Action model to API configuration
func (r *syntheticTestResourceFramework) mapWebpageActionFromModel(webpageActionModel *WebpageActionConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	return restapi.SyntheticTestConfig{
		MarkSyntheticCall: webpageActionModel.MarkSyntheticCall.ValueBool(),
		Retries:           int32(webpageActionModel.Retries.ValueInt64()),
		RetryInterval:     int32(webpageActionModel.RetryInterval.ValueInt64()),
		SyntheticType:     SyntheticTestTypeWebpageAction,
		Timeout:           getStringPointerFromFrameworkType(webpageActionModel.Timeout),
		URL:               getStringPointerFromFrameworkType(webpageActionModel.URL),
		Browser:           getStringPointerFromFrameworkType(webpageActionModel.Browser),
		RecordVideo:       getBoolPointerFromFrameworkType(webpageActionModel.RecordVideo),
	}, diags
}

// mapWebpageScriptFromModel maps Webpage Script model to API configuration
func (r *syntheticTestResourceFramework) mapWebpageScriptFromModel(webpageScriptModel *WebpageScriptConfigModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	return restapi.SyntheticTestConfig{
		MarkSyntheticCall: webpageScriptModel.MarkSyntheticCall.ValueBool(),
		Retries:           int32(webpageScriptModel.Retries.ValueInt64()),
		RetryInterval:     int32(webpageScriptModel.RetryInterval.ValueInt64()),
		SyntheticType:     SyntheticTestTypeWebpageScript,
		Timeout:           getStringPointerFromFrameworkType(webpageScriptModel.Timeout),
		Script:            getStringPointerFromFrameworkType(webpageScriptModel.Script),
		FileName:          getStringPointerFromFrameworkType(webpageScriptModel.FileName),
		Browser:           getStringPointerFromFrameworkType(webpageScriptModel.Browser),
		RecordVideo:       getBoolPointerFromFrameworkType(webpageScriptModel.RecordVideo),
	}, diags
}

func (r *syntheticTestResourceFramework) mapConfigurationFromModel(ctx context.Context, model SyntheticTestModel) (restapi.SyntheticTestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate exactly one configuration type is set
	_, validationDiags := r.validateSingleConfigType(model)
	if validationDiags.HasError() {
		return restapi.SyntheticTestConfig{}, validationDiags
	}

	// Map configuration based on type
	if model.HttpAction != nil {
		return r.mapHttpActionFromModel(ctx, model.HttpAction)
	}
	if model.HttpScript != nil {
		return r.mapHttpScriptFromModel(model.HttpScript)
	}
	if model.BrowserScript != nil {
		return r.mapBrowserScriptFromModel(model.BrowserScript)
	}
	if model.DNS != nil {
		return r.mapDNSFromModel(ctx, model.DNS)
	}
	if model.SSLCertificate != nil {
		return r.mapSSLCertificateFromModel(ctx, model.SSLCertificate)
	}
	if model.WebpageAction != nil {
		return r.mapWebpageActionFromModel(model.WebpageAction)
	}
	if model.WebpageScript != nil {
		return r.mapWebpageScriptFromModel(model.WebpageScript)
	}

	// This should never happen due to validation above
	diags.AddError(SyntheticTestErrNoValidConfig, "No valid synthetic test configuration found")
	return restapi.SyntheticTestConfig{}, diags
}

func (r *syntheticTestResourceFramework) UpdateState(ctx context.Context, state *tfsdk.State, plan *tfsdk.Plan, apiObject *restapi.SyntheticTest) diag.Diagnostics {
	var diags diag.Diagnostics
	var model SyntheticTestModel

	// Get current state from plan or state
	//diags.Append(r.getModelFromState(ctx, plan, state, &model)...)
	// Map basic fields
	if plan != nil {
		log.Printf("lloading from plan")
		diags.Append(plan.Get(ctx, &model)...)
		model.ID = types.StringValue(apiObject.ID)
		model.Label = types.StringValue(apiObject.Label)
		model.Active = types.BoolValue(apiObject.Active)
		model.PlaybackMode = types.StringValue(apiObject.PlaybackMode)
		model.Description = util.SetStringPointerToState(apiObject.Description)
		model.ApplicationID = util.SetStringPointerToState(apiObject.ApplicationID)
	} else if state != nil {
		diags.Append(state.Get(ctx, &model)...)
		log.Printf("lloading from state")
		model.ID = types.StringValue(apiObject.ID)
		model.Label = types.StringValue(apiObject.Label)
		model.Active = types.BoolValue(apiObject.Active)
		model.PlaybackMode = types.StringValue(apiObject.PlaybackMode)
		model.Description = util.SetStringPointerToState(apiObject.Description)
		model.ApplicationID = util.SetStringPointerToState(apiObject.ApplicationID)
	} else {
		model = r.mapBasicFieldsToModel(apiObject)
	}

	// Map collection fields
	model.Applications = r.mapApplicationsToModel(apiObject)
	model.MobileApps = r.mapMobileAppsToModel(apiObject)
	model.Websites = r.mapWebsitesToModel(apiObject)
	model.TestFrequency = r.mapTestFrequencyToModel(apiObject)
	model.CustomProperties = r.mapCustomPropertiesToModel(apiObject)
	model.Locations = r.mapLocationsToModel(apiObject)
	model.RbacTags = r.mapRbacTagsToModel(apiObject)

	// Map configuration based on synthetic type
	switch apiObject.Configuration.SyntheticType {
	case SyntheticTestTypeHTTPAction:
		model.HttpAction = r.mapHttpActionConfigToModel(apiObject.Configuration, model.HttpAction)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeHTTPAction)
	case SyntheticTestTypeHTTPScript:
		model.HttpScript = r.mapHttpScriptConfigToModel(apiObject.Configuration)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeHTTPScript)
	case SyntheticTestTypeBrowserScript:
		model.BrowserScript = r.mapBrowserScriptConfigToModel(apiObject.Configuration)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeBrowserScript)
	case SyntheticTestTypeDNS:
		model.DNS = r.mapDNSConfigToModel(apiObject.Configuration)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeDNS)
	case SyntheticTestTypeSSLCertificate:
		model.SSLCertificate = r.mapSSLCertificateConfigToModel(apiObject.Configuration)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeSSLCertificate)
	case SyntheticTestTypeWebpageAction:
		model.WebpageAction = r.mapWebpageActionConfigToModel(apiObject.Configuration)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeWebpageAction)
	case SyntheticTestTypeWebpageScript:
		model.WebpageScript = r.mapWebpageScriptConfigToModel(apiObject.Configuration)
		r.clearOtherConfigTypes(&model, SyntheticTestTypeWebpageScript)
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

// UpdateState helper methods for mapping API objects to Terraform state

// mapBasicFieldsToModel maps basic fields from API object to model
func (r *syntheticTestResourceFramework) mapBasicFieldsToModel(apiObject *restapi.SyntheticTest) SyntheticTestModel {
	return SyntheticTestModel{
		ID:            types.StringValue(apiObject.ID),
		Label:         types.StringValue(apiObject.Label),
		Active:        types.BoolValue(apiObject.Active),
		PlaybackMode:  types.StringValue(apiObject.PlaybackMode),
		Description:   util.SetStringPointerToState(apiObject.Description),
		ApplicationID: util.SetStringPointerToState(apiObject.ApplicationID),
	}
}

// mapApplicationsToModel maps applications array to model
func (r *syntheticTestResourceFramework) mapApplicationsToModel(apiObject *restapi.SyntheticTest) types.Set {
	if len(apiObject.Applications) > 0 {
		appValues := make([]attr.Value, len(apiObject.Applications))
		for i, app := range apiObject.Applications {
			appValues[i] = types.StringValue(app)
		}
		return types.SetValueMust(types.StringType, appValues)
	}
	return types.SetNull(types.StringType)
}

// mapMobileAppsToModel maps mobile apps array to model
func (r *syntheticTestResourceFramework) mapMobileAppsToModel(apiObject *restapi.SyntheticTest) types.Set {
	if len(apiObject.MobileApps) > 0 {
		mobileAppValues := make([]attr.Value, len(apiObject.MobileApps))
		for i, app := range apiObject.MobileApps {
			mobileAppValues[i] = types.StringValue(app)
		}
		return types.SetValueMust(types.StringType, mobileAppValues)
	}
	return types.SetNull(types.StringType)
}

// mapWebsitesToModel maps websites array to model
func (r *syntheticTestResourceFramework) mapWebsitesToModel(apiObject *restapi.SyntheticTest) types.Set {
	if len(apiObject.Websites) > 0 {
		websiteValues := make([]attr.Value, len(apiObject.Websites))
		for i, website := range apiObject.Websites {
			websiteValues[i] = types.StringValue(website)
		}
		return types.SetValueMust(types.StringType, websiteValues)
	}
	return types.SetNull(types.StringType)
}

// mapTestFrequencyToModel maps test frequency to model
func (r *syntheticTestResourceFramework) mapTestFrequencyToModel(apiObject *restapi.SyntheticTest) types.Int64 {
	if apiObject.TestFrequency != nil {
		return util.SetInt64PointerToState(apiObject.TestFrequency)
	}
	return types.Int64Null()
}

// mapCustomPropertiesToModel maps custom properties to model
func (r *syntheticTestResourceFramework) mapCustomPropertiesToModel(apiObject *restapi.SyntheticTest) types.Map {
	if len(apiObject.CustomProperties) > 0 {
		customPropertiesMap := make(map[string]attr.Value)
		for k, v := range apiObject.CustomProperties {
			customPropertiesMap[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		return types.MapValueMust(types.StringType, customPropertiesMap)
	}
	return types.MapNull(types.StringType)
}

// mapLocationsToModel maps locations array to model
func (r *syntheticTestResourceFramework) mapLocationsToModel(apiObject *restapi.SyntheticTest) types.Set {
	if len(apiObject.Locations) > 0 {
		locationValues := make([]attr.Value, len(apiObject.Locations))
		for i, location := range apiObject.Locations {
			locationValues[i] = types.StringValue(location)
		}
		return types.SetValueMust(types.StringType, locationValues)
	}
	return types.SetNull(types.StringType)
}

// mapRbacTagsToModel maps RBAC tags to model
func (r *syntheticTestResourceFramework) mapRbacTagsToModel(apiObject *restapi.SyntheticTest) types.Set {
	if len(apiObject.RbacTags) > 0 {
		rbacTagValues := make([]attr.Value, len(apiObject.RbacTags))
		for i, tag := range apiObject.RbacTags {
			tagObj, _ := types.ObjectValue(
				map[string]attr.Type{
					SyntheticTestFieldName:  types.StringType,
					SyntheticTestFieldValue: types.StringType,
				},
				map[string]attr.Value{
					SyntheticTestFieldName:  types.StringValue(tag.Name),
					SyntheticTestFieldValue: types.StringValue(tag.Value),
				},
			)
			rbacTagValues[i] = tagObj
		}
		return types.SetValueMust(
			types.ObjectType{AttrTypes: map[string]attr.Type{
				SyntheticTestFieldName:  types.StringType,
				SyntheticTestFieldValue: types.StringType,
			}},
			rbacTagValues,
		)
	}
	return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
		SyntheticTestFieldName:  types.StringType,
		SyntheticTestFieldValue: types.StringType,
	}})
}

// mapHttpActionConfigToModel maps HTTP Action configuration to model
func (r *syntheticTestResourceFramework) mapHttpActionConfigToModel(config restapi.SyntheticTestConfig, httpModel *HttpActionConfigModel) *HttpActionConfigModel {
	httpActionModel := &HttpActionConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		URL:               util.SetStringPointerToState(config.URL),
		Operation:         util.SetStringPointerToState(config.Operation),
		Body:              util.SetStringPointerToState(config.Body),
		ValidationString:  util.SetStringPointerToState(config.ValidationString),
		ExpectMatch:       util.SetStringPointerToState(config.ExpectMatch),
	}

	if config.FollowRedirect != nil {
		httpActionModel.FollowRedirect = types.BoolValue(*config.FollowRedirect)
	}
	if config.AllowInsecure != nil {
		httpActionModel.AllowInsecure = types.BoolValue(*config.AllowInsecure)
	}
	if config.ExpectStatus != nil {
		httpActionModel.ExpectStatus = util.SetInt64PointerToState(config.ExpectStatus)
	}

	httpActionModel.Headers = r.mapHeadersToModel(config.Headers)
	httpActionModel.ExpectExists = r.mapExpectExistsToModel(config.ExpectExists)
	httpActionModel.ExpectNotEmpty = r.mapExpectNotEmptyToModel(config.ExpectNotEmpty)
	if httpModel != nil {
		httpActionModelCurrent := *httpModel
		httpActionModel.ExpectJson = httpActionModelCurrent.ExpectJson
	} else {
		log.Printf("updating json again")
		httpActionModel.ExpectJson = r.mapExpectJsonToModel(config.ExpectJson)
	}

	return httpActionModel
}

// buildSchema builds the complete schema using builder functions
func buildSchema() schema.Schema {
	attrs := buildBaseSchema()

	// Add configuration type schemas
	attrs[SyntheticTestFieldHttpAction] = buildHttpActionSchema()
	attrs[SyntheticTestFieldHttpScript] = buildHttpScriptSchema()
	attrs[SyntheticTestFieldBrowserScript] = buildBrowserScriptSchema()
	attrs[SyntheticTestFieldDNS] = buildDNSSchema()
	attrs[SyntheticTestFieldSSLCertificate] = buildSSLCertificateSchema()
	attrs[SyntheticTestFieldWebpageAction] = buildWebpageActionSchema()
	attrs[SyntheticTestFieldWebpageScript] = buildWebpageScriptSchema()

	return schema.Schema{
		Description: SyntheticTestDescResource,
		Attributes:  attrs,
	}
}

// mapHeadersToModel maps headers to model
func (r *syntheticTestResourceFramework) mapHeadersToModel(headers map[string]string) types.Map {
	if len(headers) > 0 {
		headersMap := make(map[string]attr.Value)
		for k, v := range headers {
			headersMap[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		return types.MapValueMust(types.StringType, headersMap)
	}
	return types.MapNull(types.StringType)
}

// mapExpectExistsToModel maps expect exists to model
func (r *syntheticTestResourceFramework) mapExpectExistsToModel(expectExists []string) types.Set {
	if len(expectExists) > 0 {
		expectExistsValues := make([]attr.Value, len(expectExists))
		for i, val := range expectExists {
			expectExistsValues[i] = types.StringValue(val)
		}
		return types.SetValueMust(types.StringType, expectExistsValues)
	}
	return types.SetNull(types.StringType)
}

// mapExpectNotEmptyToModel maps expect not empty to model
func (r *syntheticTestResourceFramework) mapExpectNotEmptyToModel(expectNotEmpty []string) types.Set {
	if len(expectNotEmpty) > 0 {
		expectNotEmptyValues := make([]attr.Value, len(expectNotEmpty))
		for i, val := range expectNotEmpty {
			expectNotEmptyValues[i] = types.StringValue(val)
		}
		return types.SetValueMust(types.StringType, expectNotEmptyValues)
	}
	return types.SetNull(types.StringType)
}

// mapExpectJsonToModel maps expect json to model
func (r *syntheticTestResourceFramework) mapExpectJsonToModel(expectJson json.RawMessage) types.String {
	if expectJson == nil || len(expectJson) == 0 {
		return types.StringNull()
	}
	jsonBytes, _ := expectJson.MarshalJSON()
	jsonString := string(jsonBytes)
	// Check if the JSON is null
	if jsonString == "null" {
		return types.StringNull()
	}
	return types.StringValue(util.NormalizeJSONString(jsonString))
}

// mapHttpScriptConfigToModel maps HTTP Script configuration to model
func (r *syntheticTestResourceFramework) mapHttpScriptConfigToModel(config restapi.SyntheticTestConfig) *HttpScriptConfigModel {
	httpScriptModel := &HttpScriptConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		Script:            util.SetStringPointerToState(config.Script),
		ScriptType:        util.SetStringPointerToState(config.ScriptType),
	}

	if config.FileName != nil && *config.FileName != "" {
		httpScriptModel.FileName = types.StringValue(*config.FileName)
	} else {
		httpScriptModel.FileName = types.StringNull()
	}

	if config.Scripts != nil {
		httpScriptModel.Scripts = &MultipleScriptsModel{
			Bundle:     util.SetStringPointerToState(config.Scripts.Bundle),
			ScriptFile: util.SetStringPointerToState(config.Scripts.ScriptFile),
		}
	}

	return httpScriptModel
}

// mapBrowserScriptConfigToModel maps Browser Script configuration to model
func (r *syntheticTestResourceFramework) mapBrowserScriptConfigToModel(config restapi.SyntheticTestConfig) *BrowserScriptConfigModel {
	browserScriptModel := &BrowserScriptConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		Script:            util.SetStringPointerToState(config.Script),
		ScriptType:        util.SetStringPointerToState(config.ScriptType),
		Browser:           util.SetStringPointerToState(config.Browser),
	}

	if config.FileName != nil && *config.FileName != "" {
		browserScriptModel.FileName = types.StringValue(*config.FileName)
	} else {
		browserScriptModel.FileName = types.StringNull()
	}

	if config.RecordVideo != nil {
		browserScriptModel.RecordVideo = types.BoolValue(*config.RecordVideo)
	}

	if config.Scripts != nil {
		browserScriptModel.Scripts = &MultipleScriptsModel{
			Bundle:     util.SetStringPointerToState(config.Scripts.Bundle),
			ScriptFile: util.SetStringPointerToState(config.Scripts.ScriptFile),
		}
	}

	return browserScriptModel
}

// mapDNSConfigToModel maps DNS configuration to model
func (r *syntheticTestResourceFramework) mapDNSConfigToModel(config restapi.SyntheticTestConfig) *DNSConfigModel {
	dnsModel := &DNSConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		Lookup:            util.SetStringPointerToState(config.Lookup),
		Server:            util.SetStringPointerToState(config.Server),
		QueryType:         util.SetStringPointerToState(config.QueryType),
		Transport:         util.SetStringPointerToState(config.Transport),
	}

	if config.Port != nil {
		dnsModel.Port = util.SetInt64PointerToState(config.Port)
	}
	if config.ServerRetries != nil {
		dnsModel.ServerRetries = util.SetInt64PointerToState(config.ServerRetries)
	}
	if config.AcceptCNAME != nil {
		dnsModel.AcceptCNAME = types.BoolValue(*config.AcceptCNAME)
	}
	if config.LookupServerName != nil {
		dnsModel.LookupServerName = types.BoolValue(*config.LookupServerName)
	}
	if config.RecursiveLookups != nil {
		dnsModel.RecursiveLookups = types.BoolValue(*config.RecursiveLookups)
	}

	if config.QueryTime != nil {
		dnsModel.QueryTime = &DNSFilterQueryTimeModel{
			Key:      types.StringValue(config.QueryTime.Key),
			Operator: types.StringValue(config.QueryTime.Operator),
			Value:    types.Int64Value(config.QueryTime.Value),
		}
	}

	dnsModel.TargetValues = r.mapDNSTargetValuesToModel(config.TargetValues)

	return dnsModel
}

// mapDNSTargetValuesToModel maps DNS target values to model
func (r *syntheticTestResourceFramework) mapDNSTargetValuesToModel(targetValues []restapi.DNSFilterTargetValue) types.Set {
	if targetValues != nil && len(targetValues) > 0 {
		targetValueObjs := make([]attr.Value, len(targetValues))
		for i, tv := range targetValues {
			tvObj, _ := types.ObjectValue(
				map[string]attr.Type{
					SyntheticTestFieldKey:      types.StringType,
					SyntheticTestFieldOperator: types.StringType,
					SyntheticTestFieldValue:    types.StringType,
				},
				map[string]attr.Value{
					SyntheticTestFieldKey:      types.StringValue(tv.Key),
					SyntheticTestFieldOperator: types.StringValue(tv.Operator),
					SyntheticTestFieldValue:    types.StringValue(tv.Value),
				},
			)
			targetValueObjs[i] = tvObj
		}
		return types.SetValueMust(
			types.ObjectType{AttrTypes: map[string]attr.Type{
				SyntheticTestFieldKey:      types.StringType,
				SyntheticTestFieldOperator: types.StringType,
				SyntheticTestFieldValue:    types.StringType,
			}},
			targetValueObjs,
		)
	}
	return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
		SyntheticTestFieldKey:      types.StringType,
		SyntheticTestFieldOperator: types.StringType,
		SyntheticTestFieldValue:    types.StringType,
	}})
}

// mapSSLCertificateConfigToModel maps SSL Certificate configuration to model
func (r *syntheticTestResourceFramework) mapSSLCertificateConfigToModel(config restapi.SyntheticTestConfig) *SSLCertificateConfigModel {
	sslModel := &SSLCertificateConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		Hostname:          util.SetStringPointerToState(config.Hostname),
	}

	if config.DaysRemainingCheck != nil {
		sslModel.DaysRemainingCheck = util.SetInt64PointerToState(config.DaysRemainingCheck)
	}
	if config.SSLPort != nil {
		sslModel.Port = util.SetInt64PointerToState(config.SSLPort)
	}
	if config.AcceptSelfSignedCert != nil {
		sslModel.AcceptSelfSignedCert = types.BoolValue(*config.AcceptSelfSignedCert)
	}

	sslModel.ValidationRules = r.mapSSLValidationRulesToModel(config.ValidationRules)

	return sslModel
}

// mapSSLValidationRulesToModel maps SSL validation rules to model
func (r *syntheticTestResourceFramework) mapSSLValidationRulesToModel(validationRules []restapi.SSLCertificateValidation) types.Set {
	if validationRules != nil && len(validationRules) > 0 {
		validationRuleObjs := make([]attr.Value, len(validationRules))
		for i, vr := range validationRules {
			vrObj, _ := types.ObjectValue(
				map[string]attr.Type{
					SyntheticTestFieldKey:      types.StringType,
					SyntheticTestFieldOperator: types.StringType,
					SyntheticTestFieldValue:    types.StringType,
				},
				map[string]attr.Value{
					SyntheticTestFieldKey:      types.StringValue(vr.Key),
					SyntheticTestFieldOperator: types.StringValue(vr.Operator),
					SyntheticTestFieldValue:    types.StringValue(fmt.Sprintf("%v", vr.Value)),
				},
			)
			validationRuleObjs[i] = vrObj
		}
		return types.SetValueMust(
			types.ObjectType{AttrTypes: map[string]attr.Type{
				SyntheticTestFieldKey:      types.StringType,
				SyntheticTestFieldOperator: types.StringType,
				SyntheticTestFieldValue:    types.StringType,
			}},
			validationRuleObjs,
		)
	}
	return types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
		SyntheticTestFieldKey:      types.StringType,
		SyntheticTestFieldOperator: types.StringType,
		SyntheticTestFieldValue:    types.StringType,
	}})
}

// mapWebpageActionConfigToModel maps Webpage Action configuration to model
func (r *syntheticTestResourceFramework) mapWebpageActionConfigToModel(config restapi.SyntheticTestConfig) *WebpageActionConfigModel {
	webpageActionModel := &WebpageActionConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		URL:               util.SetStringPointerToState(config.URL),
		Browser:           util.SetStringPointerToState(config.Browser),
	}

	if config.RecordVideo != nil {
		webpageActionModel.RecordVideo = types.BoolValue(*config.RecordVideo)
	}

	return webpageActionModel
}

// mapWebpageScriptConfigToModel maps Webpage Script configuration to model
func (r *syntheticTestResourceFramework) mapWebpageScriptConfigToModel(config restapi.SyntheticTestConfig) *WebpageScriptConfigModel {
	webpageScriptModel := &WebpageScriptConfigModel{
		MarkSyntheticCall: types.BoolValue(config.MarkSyntheticCall),
		Retries:           types.Int64Value(int64(config.Retries)),
		RetryInterval:     types.Int64Value(int64(config.RetryInterval)),
		Timeout:           util.SetStringPointerToState(config.Timeout),
		Script:            util.SetStringPointerToState(config.Script),
		Browser:           util.SetStringPointerToState(config.Browser),
	}

	if config.FileName != nil && *config.FileName != "" {
		webpageScriptModel.FileName = types.StringValue(*config.FileName)
	} else {
		webpageScriptModel.FileName = types.StringNull()
	}

	if config.RecordVideo != nil {
		webpageScriptModel.RecordVideo = types.BoolValue(*config.RecordVideo)
	}

	return webpageScriptModel
}

// clearOtherConfigTypes sets all configuration types except the specified one to nil
func (r *syntheticTestResourceFramework) clearOtherConfigTypes(model *SyntheticTestModel, keepType string) {
	if keepType != SyntheticTestTypeHTTPAction {
		model.HttpAction = nil
	}
	if keepType != SyntheticTestTypeHTTPScript {
		model.HttpScript = nil
	}
	if keepType != SyntheticTestTypeBrowserScript {
		model.BrowserScript = nil
	}
	if keepType != SyntheticTestTypeDNS {
		model.DNS = nil
	}
	if keepType != SyntheticTestTypeSSLCertificate {
		model.SSLCertificate = nil
	}
	if keepType != SyntheticTestTypeWebpageAction {
		model.WebpageAction = nil
	}
	if keepType != SyntheticTestTypeWebpageScript {
		model.WebpageScript = nil
	}
}
