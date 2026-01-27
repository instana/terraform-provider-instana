package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/instana/terraform-provider-instana/internal/datasources"
	"github.com/instana/terraform-provider-instana/internal/resourcehandle"
	"github.com/instana/terraform-provider-instana/internal/resources/alertingchannel"
	"github.com/instana/terraform-provider-instana/internal/resources/alertingconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/apitoken"
	"github.com/instana/terraform-provider-instana/internal/resources/applicationalertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/applicationconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/automationaction"
	"github.com/instana/terraform-provider-instana/internal/resources/automationpolicy"
	"github.com/instana/terraform-provider-instana/internal/resources/customdashboard"
	"github.com/instana/terraform-provider-instana/internal/resources/customeventspec"
	"github.com/instana/terraform-provider-instana/internal/resources/group"
	"github.com/instana/terraform-provider-instana/internal/resources/infralertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/logalertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/mobilealertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/roles"
	"github.com/instana/terraform-provider-instana/internal/resources/sliconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/sloalertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/sloconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/slocorrectionconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/syntheticalertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/synthetictest"
	"github.com/instana/terraform-provider-instana/internal/resources/team"
	"github.com/instana/terraform-provider-instana/internal/resources/websitealertconfig"
	"github.com/instana/terraform-provider-instana/internal/resources/websitemonitoringconfig"
	"github.com/instana/terraform-provider-instana/internal/restapi"
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
				Description: "API token used to authenticate with the Instana Backend. Can also be set via INSTANA_API_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			SchemaFieldEndpoint: schema.StringAttribute{
				Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io). Can also be set via INSTANA_ENDPOINT environment variable.",
				Optional:    true,
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

	// Default values to environment variables, but override with Terraform configuration value if set
	apiToken := strings.TrimSpace(config.APIToken.ValueString())
	if apiToken == "" {
		apiToken = os.Getenv("INSTANA_API_TOKEN")
	}

	endpoint := strings.TrimSpace(config.Endpoint.ValueString())
	if endpoint == "" {
		endpoint = os.Getenv("INSTANA_ENDPOINT")
	}

	// Validate that required values are present
	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Instana API Token",
			"The provider cannot create the Instana API client as there is a missing or empty value for the Instana API token. "+
				"Set the api_token value in the configuration or use the INSTANA_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Instana Endpoint",
			"The provider cannot create the Instana API client as there is a missing or empty value for the Instana endpoint. "+
				"Set the endpoint value in the configuration or use the INSTANA_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
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
		datasources.NewAlertingChannelDataSource,
		datasources.NewAutomationActionDataSource,
		datasources.NewBuiltinEventDataSource,
		datasources.NewCustomEventSpecificationDataSource,
		datasources.NewHostAgentsDataSource,
		datasources.NewSyntheticLocationDataSource,
		datasources.NewUserDataSource,
	}
}

// Resources defines the resources implemented in the provider
func (p *InstanaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add resources here -
		addResouceHandle(alertingconfig.NewAlertingConfigResourceHandle),
		addResouceHandle(logalertconfig.NewLogAlertConfigResourceHandle),
		addResouceHandle(alertingchannel.NewAlertingChannelResourceHandle),
		addResouceHandle(apitoken.NewAPITokenResourceHandle),
		addResouceHandle(applicationalertconfig.NewApplicationAlertConfigResourceHandle),
		addResouceHandle(applicationalertconfig.NewGlobalApplicationAlertConfigResourceHandle),
		addResouceHandle(applicationconfig.NewApplicationConfigResourceHandle),
		addResouceHandle(automationaction.NewAutomationActionResourceHandle),
		addResouceHandle(automationpolicy.NewAutomationPolicyResourceHandle),
		addResouceHandle(customdashboard.NewCustomDashboardResourceHandle),
		addResouceHandle(customeventspec.NewCustomEventSpecificationResourceHandle),
		addResouceHandle(infralertconfig.NewInfraAlertConfigResourceHandle),
		addResouceHandle(mobilealertconfig.NewMobileAlertConfigResourceHandle),
		addResouceHandle(group.NewGroupResourceHandle),
		addResouceHandle(team.NewTeamResourceHandle),
		addResouceHandle(roles.NewRoleResourceHandle),
		addResouceHandle(sliconfig.NewSliConfigResourceHandle),
		addResouceHandle(sloalertconfig.NewSloAlertConfigResourceHandle),
		addResouceHandle(slocorrectionconfig.NewSloCorrectionConfigResourceHandle),
		addResouceHandle(syntheticalertconfig.NewSyntheticAlertConfigResourceHandle),
		addResouceHandle(synthetictest.NewSyntheticTestResourceHandle),
		addResouceHandle(websitealertconfig.NewWebsiteAlertConfigResourceHandle),
		addResouceHandle(websitemonitoringconfig.NewWebsiteMonitoringConfigResourceHandle),
		addResouceHandle(sloconfig.NewSloConfigResourceHandle),
	}
}

// Helper function to wrap resource handles
func addResouceHandle[T restapi.InstanaDataObject](handleFunc func() resourcehandle.ResourceHandle[T]) func() resource.Resource {
	return func() resource.Resource {
		return NewTerraformResource(handleFunc())
	}
}
