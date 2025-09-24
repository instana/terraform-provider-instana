package restapi

const (
	//InstanaAPIBasePath path to Instana RESTful API
	InstanaAPIBasePath = "/api"
	//EventsBasePath path to Events resource of Instana RESTful API
	EventsBasePath = InstanaAPIBasePath + "/events"
	//settingsPathElement path element to settings
	settingsPathElement = "/settings"
	//EventSettingsBasePath path to Event Settings resource of Instana RESTful API
	EventSettingsBasePath = EventsBasePath + settingsPathElement
	//SettingsBasePath path to Event Settings resource of Instana RESTful API
	SettingsBasePath = InstanaAPIBasePath + settingsPathElement
	//RBACSettingsBasePath path to Role Based Access Control Settings resources of Instana RESTful API
	RBACSettingsBasePath = SettingsBasePath + "/rbac"
	//WebsiteMonitoringResourcePath path to website monitoring
	WebsiteMonitoringResourcePath = InstanaAPIBasePath + "/website-monitoring"
	//SyntheticSettingsBasePath path to synthetic monitoring
	SyntheticSettingsBasePath = InstanaAPIBasePath + "/synthetics" + settingsPathElement
	//SyntheticTestResourcePath path to synthetic monitoring tests
	SyntheticTestResourcePath = SyntheticSettingsBasePath + "/tests"
	//SyntheticLocationResourcePath path to synthetic monitoring tests
	SyntheticLocationResourcePath = SyntheticSettingsBasePath + "/locations"
	// AutomationBasePath path to Automation resources of Instana RESTful API
	AutomationBasePath = InstanaAPIBasePath + "/automation"
	// AutomationBasePath path to Automation resources of Instana RESTful API
	HostAgentResourcePath = InstanaAPIBasePath + "/host-agent"
)

// InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	CustomEventSpecifications() RestResource[*CustomEventSpecification]
	BuiltinEventSpecifications() ReadOnlyRestResource[*BuiltinEventSpecification]
	APITokens() RestResource[*APIToken]
	ApplicationConfigs() RestResource[*ApplicationConfig]
	ApplicationAlertConfigs() RestResource[*ApplicationAlertConfig]
	GlobalApplicationAlertConfigs() RestResource[*ApplicationAlertConfig]
	AlertingChannels() RestResource[*AlertingChannel]
	AlertingConfigurations() RestResource[*AlertingConfiguration]
	SliConfigs() RestResource[*SliConfig]
	SloConfigs() RestResource[*SloConfig]
	SloAlertConfig() RestResource[*SloAlertConfig]
	SloCorrectionConfig() RestResource[*SloCorrectionConfig]
	WebsiteMonitoringConfig() RestResource[*WebsiteMonitoringConfig]
	WebsiteAlertConfig() RestResource[*WebsiteAlertConfig]
	InfraAlertConfig() RestResource[*InfraAlertConfig]
	Groups() RestResource[*Group]
	CustomDashboards() RestResource[*CustomDashboard]
	SyntheticTest() RestResource[*SyntheticTest]
	SyntheticLocation() ReadOnlyRestResource[*SyntheticLocation]
	SyntheticAlertConfigs() RestResource[*SyntheticAlertConfig]
	AutomationActions() RestResource[*AutomationAction]
	AutomationPolicies() RestResource[*AutomationPolicy]
	HostAgents() ReadOnlyRestResource[*HostAgent]
}

// NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string, skipTlsVerification bool) InstanaAPI {
	client := NewClient(apiToken, endpoint, skipTlsVerification)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client RestClient
}

// CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) CustomEventSpecifications() RestResource[*CustomEventSpecification] {
	return NewCreatePUTUpdatePUTRestResource(CustomEventSpecificationResourcePath, NewDefaultJSONUnmarshaller(&CustomEventSpecification{}), api.client)
}

// BuiltinEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) BuiltinEventSpecifications() ReadOnlyRestResource[*BuiltinEventSpecification] {
	return NewReadOnlyRestResource(BuiltinEventSpecificationResourcePath, NewDefaultJSONUnmarshaller(&BuiltinEventSpecification{}), api.client)
}

// APITokens implementation of InstanaAPI interface
func (api *baseInstanaAPI) APITokens() RestResource[*APIToken] {
	return NewCreatePOSTUpdatePUTRestResource(APITokensResourcePath, NewDefaultJSONUnmarshaller(&APIToken{}), api.client)
}

// ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() RestResource[*ApplicationConfig] {
	return NewCreatePOSTUpdatePUTRestResource(ApplicationConfigsResourcePath, NewDefaultJSONUnmarshaller(&ApplicationConfig{}), api.client)
}

// ApplicationAlertConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationAlertConfigs() RestResource[*ApplicationAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(ApplicationAlertConfigsResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&ApplicationAlertConfig{})), api.client)
}

// GlobalApplicationAlertConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) GlobalApplicationAlertConfigs() RestResource[*ApplicationAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(GlobalApplicationAlertConfigsResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&ApplicationAlertConfig{})), api.client)
}

// AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() RestResource[*AlertingChannel] {
	return NewCreatePUTUpdatePUTRestResource(AlertingChannelsResourcePath, NewDefaultJSONUnmarshaller(&AlertingChannel{}), api.client)
}

// AlertingConfigurations implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingConfigurations() RestResource[*AlertingConfiguration] {
	return NewCreatePUTUpdatePUTRestResource(AlertsResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&AlertingConfiguration{})), api.client)
}

func (api *baseInstanaAPI) SliConfigs() RestResource[*SliConfig] {
	return NewCreatePOSTUpdateNotSupportedRestResource(SliConfigResourcePath, NewDefaultJSONUnmarshaller(&SliConfig{}), api.client)
}

func (api *baseInstanaAPI) SloConfigs() RestResource[*SloConfig] {
	return NewCreatePOSTUpdatePUTRestResource(SloConfigResourcePath, NewDefaultJSONUnmarshaller(&SloConfig{}), api.client)
}

func (api *baseInstanaAPI) SloAlertConfig() RestResource[*SloAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(SloAlertConfigResourcePath, NewDefaultJSONUnmarshaller(&SloAlertConfig{}), api.client)
}

func (api *baseInstanaAPI) SloCorrectionConfig() RestResource[*SloCorrectionConfig] {
	return NewCreatePOSTUpdatePUTRestResource(SloCorrectionConfigResourcePath, NewDefaultJSONUnmarshaller(&SloCorrectionConfig{}), api.client)
}

func (api *baseInstanaAPI) WebsiteMonitoringConfig() RestResource[*WebsiteMonitoringConfig] {
	return NewWebsiteMonitoringConfigRestResource(NewDefaultJSONUnmarshaller(&WebsiteMonitoringConfig{}), api.client)
}

func (api *baseInstanaAPI) WebsiteAlertConfig() RestResource[*WebsiteAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(WebsiteAlertConfigResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&WebsiteAlertConfig{})), api.client)
}

func (api *baseInstanaAPI) InfraAlertConfig() RestResource[*InfraAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(InfraAlertConfigResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&InfraAlertConfig{})), api.client)
}

func (api *baseInstanaAPI) Groups() RestResource[*Group] {
	return NewCreatePOSTUpdatePUTRestResource(GroupsResourcePath, NewDefaultJSONUnmarshaller(&Group{}), api.client)
}

func (api *baseInstanaAPI) CustomDashboards() RestResource[*CustomDashboard] {
	return NewCreatePOSTUpdatePUTRestResource(CustomDashboardsResourcePath, NewDefaultJSONUnmarshaller(&CustomDashboard{}), api.client)
}

func (api *baseInstanaAPI) SyntheticTest() RestResource[*SyntheticTest] {
	return NewSyntheticTestRestResource(NewDefaultJSONUnmarshaller(&SyntheticTest{}), api.client)
}

// SyntheticLocation implementation of InstanaAPI interface
func (api *baseInstanaAPI) SyntheticLocation() ReadOnlyRestResource[*SyntheticLocation] {
	return NewReadOnlyRestResource(SyntheticLocationResourcePath, NewDefaultJSONUnmarshaller(&SyntheticLocation{}), api.client)
}

// SyntheticAlertConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) SyntheticAlertConfigs() RestResource[*SyntheticAlertConfig] {
	return NewCreatePOSTUpdatePOSTRestResource(SyntheticAlertConfigsResourcePath, NewCustomPayloadFieldsUnmarshallerAdapter(NewDefaultJSONUnmarshaller(&SyntheticAlertConfig{})), api.client)
}

func (api *baseInstanaAPI) AutomationActions() RestResource[*AutomationAction] {
	return NewCreatePOSTUpdatePUTRestResource(AutomationActionResourcePath, NewDefaultJSONUnmarshaller(&AutomationAction{}), api.client)
}

func (api *baseInstanaAPI) AutomationPolicies() RestResource[*AutomationPolicy] {
	return NewCreatePOSTUpdatePUTRestResource(AutomationPolicyResourcePath, NewDefaultJSONUnmarshaller(&AutomationPolicy{}), api.client)
}

func (api *baseInstanaAPI) HostAgents() ReadOnlyRestResource[*HostAgent] {
	return NewReadOnlyRestResource(HostAgentResourcePath, NewHostAgentJSONUnmarshaller(&HostAgent{}), api.client)
}
