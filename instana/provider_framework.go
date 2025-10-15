package instana

import (
	"context"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
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
	resp.DataSourceData = &ProviderMeta{
		InstanaAPI: instanaAPI,
	}
	resp.ResourceData = &ProviderMeta{
		InstanaAPI: instanaAPI,
	}
}

// DataSources defines the data sources implemented in the provider
func (p *InstanaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add data sources here when implemented
	}
}

// Resources defines the resources implemented in the provider
func (p *InstanaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add resources here -
		addResouceHandle(NewAlertingConfigResourceHandleFramework),
		addResouceHandle(NewLogAlertConfigResourceHandleFramework),
		addResouceHandle(NewAlertingChannelResourceHandleFramework),
		addResouceHandle(NewAPITokenResourceHandleFramework),
		addResouceHandle(NewApplicationAlertConfigResourceHandleFramework),
		addResouceHandle(NewGlobalApplicationAlertConfigResourceHandleFramework),
		addResouceHandle(NewApplicationConfigResourceHandleFramework),
		addResouceHandle(NewAutomationActionResourceHandleFramework),
		addResouceHandle(NewAutomationPolicyResourceHandleFramework),
		addResouceHandle(NewCustomDashboardResourceHandleFramework),
		addResouceHandle(NewCustomEventSpecificationResourceHandleFramework),
		addResouceHandle(NewInfraAlertConfigResourceHandleFramework),
		addResouceHandle(NewGroupResourceHandleFramework),
		addResouceHandle(NewSliConfigResourceHandleFramework),
	}
}

// Helper function to wrap resource handles
func addResouceHandle[T restapi.InstanaDataObject](handleFunc func() ResourceHandleFramework[T]) func() resource.Resource {
	return func() resource.Resource {
		return NewTerraformResourceFramework(handleFunc())
	}
}

// Made with Bob
