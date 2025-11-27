package testutils

import "github.com/instana/terraform-provider-instana/internal/restapi"

// MockInstanaAPI is a mock implementation of the InstanaAPI interface for testing purposes.
// It returns nil for all methods by default. Tests can override specific methods by embedding
// this struct and providing custom implementations for the methods they need.
type MockInstanaAPI struct{}

// CustomEventSpecifications mock implementation
func (m *MockInstanaAPI) CustomEventSpecifications() restapi.RestResource[*restapi.CustomEventSpecification] {
	return nil
}

// BuiltinEventSpecifications mock implementation
func (m *MockInstanaAPI) BuiltinEventSpecifications() restapi.ReadOnlyRestResource[*restapi.BuiltinEventSpecification] {
	return nil
}

// APITokens mock implementation
func (m *MockInstanaAPI) APITokens() restapi.RestResource[*restapi.APIToken] {
	return nil
}

// ApplicationConfigs mock implementation
func (m *MockInstanaAPI) ApplicationConfigs() restapi.RestResource[*restapi.ApplicationConfig] {
	return nil
}

// ApplicationAlertConfigs mock implementation
func (m *MockInstanaAPI) ApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return nil
}

// GlobalApplicationAlertConfigs mock implementation
func (m *MockInstanaAPI) GlobalApplicationAlertConfigs() restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return nil
}

// AlertingChannels mock implementation
func (m *MockInstanaAPI) AlertingChannels() restapi.RestResource[*restapi.AlertingChannel] {
	return nil
}

// AlertingConfigurations mock implementation
func (m *MockInstanaAPI) AlertingConfigurations() restapi.RestResource[*restapi.AlertingConfiguration] {
	return nil
}

// SliConfigs mock implementation
func (m *MockInstanaAPI) SliConfigs() restapi.RestResource[*restapi.SliConfig] {
	return nil
}

// SloConfigs mock implementation
func (m *MockInstanaAPI) SloConfigs() restapi.RestResource[*restapi.SloConfig] {
	return nil
}

// SloAlertConfig mock implementation
func (m *MockInstanaAPI) SloAlertConfig() restapi.RestResource[*restapi.SloAlertConfig] {
	return nil
}

// SloCorrectionConfig mock implementation
func (m *MockInstanaAPI) SloCorrectionConfig() restapi.RestResource[*restapi.SloCorrectionConfig] {
	return nil
}

// WebsiteMonitoringConfig mock implementation
func (m *MockInstanaAPI) WebsiteMonitoringConfig() restapi.RestResource[*restapi.WebsiteMonitoringConfig] {
	return nil
}

// WebsiteAlertConfig mock implementation
func (m *MockInstanaAPI) WebsiteAlertConfig() restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return nil
}

// InfraAlertConfig mock implementation
func (m *MockInstanaAPI) InfraAlertConfig() restapi.RestResource[*restapi.InfraAlertConfig] {
	return nil
}

// Teams mock implementation
func (m *MockInstanaAPI) Teams() restapi.RestResource[*restapi.Team] {
	return nil
}

// Groups mock implementation
func (m *MockInstanaAPI) Groups() restapi.RestResource[*restapi.Group] {
	return nil
}

// Roles mock implementation
func (m *MockInstanaAPI) Roles() restapi.RestResource[*restapi.Role] {
	return nil
}

// CustomDashboards mock implementation
func (m *MockInstanaAPI) CustomDashboards() restapi.RestResource[*restapi.CustomDashboard] {
	return nil
}

// SyntheticTest mock implementation
func (m *MockInstanaAPI) SyntheticTest() restapi.RestResource[*restapi.SyntheticTest] {
	return nil
}

// SyntheticLocation mock implementation
func (m *MockInstanaAPI) SyntheticLocation() restapi.ReadOnlyRestResource[*restapi.SyntheticLocation] {
	return nil
}

// SyntheticAlertConfigs mock implementation
func (m *MockInstanaAPI) SyntheticAlertConfigs() restapi.RestResource[*restapi.SyntheticAlertConfig] {
	return nil
}

// AutomationActions mock implementation
func (m *MockInstanaAPI) AutomationActions() restapi.RestResource[*restapi.AutomationAction] {
	return nil
}

// AutomationPolicies mock implementation
func (m *MockInstanaAPI) AutomationPolicies() restapi.RestResource[*restapi.AutomationPolicy] {
	return nil
}

// HostAgents mock implementation
func (m *MockInstanaAPI) HostAgents() restapi.ReadOnlyRestResource[*restapi.HostAgent] {
	return nil
}

// LogAlertConfig mock implementation
func (m *MockInstanaAPI) LogAlertConfig() restapi.RestResource[*restapi.LogAlertConfig] {
	return nil
}

// Made with Bob
