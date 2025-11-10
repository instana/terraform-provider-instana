package provider

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/internal/datasources"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/alertingchannel"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/alertingconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/apitoken"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/applicationalertconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/applicationconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/automationaction"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/automationpolicy"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/customdashboard"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/cutomeventspec"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/group"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/infralertconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/logalertconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/sliconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/sloalertconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/sloconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/slocorrectionconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/syntheticalertconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/synthetictest"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/websitealertconfig"
	"github.com/gessnerfl/terraform-provider-instana/internal/resources/websitemonitoringconfig"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &InstanaProvider{}
)

// SchemaFieldAPIToken the name of the provider configuration option for the api token
const SchemaFieldAPIToken = "api_token"

// SchemaFieldEndpoint the name of the provider configuration option for the instana endpoint
const SchemaFieldEndpoint = "endpoint"

// SchemaFieldTlsSkipVerify flag to deactivate skip tls verification
const SchemaFieldTlsSkipVerify = "tls_skip_verify"

// InstanaProviderModel describes the provider data model
type InstanaProviderModel struct {
	APIToken      types.String `tfsdk:"api_token"`
	Endpoint      types.String `tfsdk:"endpoint"`
	TLSSkipVerify types.Bool   `tfsdk:"tls_skip_verify"`
}

// InstanaProvider is the provider implementation
type InstanaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New creates a new provider instance
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &InstanaProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name
func (p *InstanaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "instana"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data
func (p *InstanaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Instana provider is used to interact with the Instana monitoring platform.",
		Attributes: map[string]schema.Attribute{
			SchemaFieldAPIToken: schema.StringAttribute{
				Description: "API token used to authenticate with the Instana Backend",
				Required:    true,
				Sensitive:   true,
			},
			SchemaFieldEndpoint: schema.StringAttribute{
				Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io)",
				Required:    true,
			},
			SchemaFieldTlsSkipVerify: schema.BoolAttribute{
				Description: "If set to true, TLS verification will be skipped when calling Instana API",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a Instana API client for data sources and resources
func (p *InstanaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config InstanaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If values are not provided, retrieve from environment variables
	if config.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Instana API Token",
			"The provider cannot create the Instana API client as there is an unknown configuration value for the Instana API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INSTANA_API_TOKEN environment variable.",
		)
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Instana Endpoint",
			"The provider cannot create the Instana API client as there is an unknown configuration value for the Instana endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INSTANA_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override with Terraform configuration value if set
	apiToken := strings.TrimSpace(config.APIToken.ValueString())
	endpoint := strings.TrimSpace(config.Endpoint.ValueString())
	skipTlsVerify := false
	if !config.TLSSkipVerify.IsNull() {
		skipTlsVerify = config.TLSSkipVerify.ValueBool()
	}

	// Create a new Instana client using the configuration values
	instanaAPI := restapi.NewInstanaAPI(apiToken, endpoint, skipTlsVerify)

	// Make the Instana client available during DataSource and Resource Configure methods
	resp.DataSourceData = &restapi.ProviderMeta{
		InstanaAPI: instanaAPI,
	}
	resp.ResourceData = &restapi.ProviderMeta{
		InstanaAPI: instanaAPI,
	}
}

// DataSources defines the data sources implemented in the provider
func (p *InstanaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add data sources here when implemented
		datasources.NewAlertingChannelDataSourceFramework,
		datasources.NewAutomationActionDataSourceFramework,
		datasources.NewBuiltinEventDataSourceFramework,
		datasources.NewCustomEventSpecificationDataSourceFramework,
		datasources.NewHostAgentsDataSourceFramework,
		datasources.NewSyntheticLocationDataSourceFramework,
	}
}

// Resources defines the resources implemented in the provider
func (p *InstanaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add resources here -
		addResouceHandle(alertingconfig.NewAlertingConfigResourceHandleFramework),
		addResouceHandle(logalertconfig.NewLogAlertConfigResourceHandleFramework),
		addResouceHandle(alertingchannel.NewAlertingChannelResourceHandleFramework),
		addResouceHandle(apitoken.NewAPITokenResourceHandleFramework),
		addResouceHandle(applicationalertconfig.NewApplicationAlertConfigResourceHandleFramework),
		addResouceHandle(applicationalertconfig.NewGlobalApplicationAlertConfigResourceHandleFramework),
		addResouceHandle(applicationconfig.NewApplicationConfigResourceHandleFramework),
		addResouceHandle(automationaction.NewAutomationActionResourceHandleFramework),
		addResouceHandle(automationpolicy.NewAutomationPolicyResourceHandleFramework),
		addResouceHandle(customdashboard.NewCustomDashboardResourceHandleFramework),
		addResouceHandle(cutomeventspec.NewCustomEventSpecificationResourceHandleFramework),
		addResouceHandle(infralertconfig.NewInfraAlertConfigResourceHandleFramework),
		addResouceHandle(group.NewGroupResourceHandleFramework),
		addResouceHandle(sliconfig.NewSliConfigResourceHandleFramework),
		addResouceHandle(sloalertconfig.NewSloAlertConfigResourceHandleFramework),
		addResouceHandle(slocorrectionconfig.NewSloCorrectionConfigResourceHandleFramework),
		addResouceHandle(syntheticalertconfig.NewSyntheticAlertConfigResourceHandleFramework),
		addResouceHandle(synthetictest.NewSyntheticTestResourceHandleFramework),
		addResouceHandle(websitealertconfig.NewWebsiteAlertConfigResourceHandleFramework),
		addResouceHandle(websitemonitoringconfig.NewWebsiteMonitoringConfigResourceHandleFramework),
		addResouceHandle(sloconfig.NewSloConfigResourceHandleFramework),
	}
}

// Helper function to wrap resource handles
func addResouceHandle[T restapi.InstanaDataObject](handleFunc func() ResourceHandleFramework[T]) func() resource.Resource {
	return func() resource.Resource {
		return NewTerraformResourceFramework(handleFunc())
	}
}

// Made with Bob
