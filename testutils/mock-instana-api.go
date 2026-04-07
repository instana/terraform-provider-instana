package testutils

import (
	"github.com/instana/instana-go-client/api"
	"github.com/instana/instana-go-client/shared/rest"
)

// MockInstanaAPI is a mock implementation of the InstanaAPI interface for testing purposes.
// It returns nil for all methods by default. Tests can override specific methods by embedding
// this struct and providing custom implementations for the methods they need.
type MockInstanaAPI struct{}

// CustomEventSpecifications mock implementation
func (m *MockInstanaAPI) CustomEventSpecifications() rest.RestResource[*api.CustomEventSpecification] {
	return nil
}

// BuiltinEventSpecifications mock implementation
func (m *MockInstanaAPI) BuiltinEventSpecifications() rest.ReadOnlyRestResource[*api.BuiltinEventSpecification] {
	return nil
}

// APITokens mock implementation
func (m *MockInstanaAPI) APITokens() rest.RestResource[*api.APIToken] {
	return nil
}

// ApplicationConfigs mock implementation
func (m *MockInstanaAPI) ApplicationConfigs() rest.RestResource[*api.ApplicationConfig] {
	return nil
}

// ApplicationAlertConfigs mock implementation
func (m *MockInstanaAPI) ApplicationAlertConfigs() rest.RestResource[*api.ApplicationAlertConfig] {
	return nil
}

// GlobalApplicationAlertConfigs mock implementation
func (m *MockInstanaAPI) GlobalApplicationAlertConfigs() rest.RestResource[*api.ApplicationAlertConfig] {
	return nil
}

// AlertingChannels mock implementation
func (m *MockInstanaAPI) AlertingChannels() rest.RestResource[*api.AlertingChannel] {
	return nil
}

// AlertingConfigurations mock implementation
func (m *MockInstanaAPI) AlertingConfigurations() rest.RestResource[*api.AlertingConfiguration] {
	return nil
}

// SliConfigs mock implementation
func (m *MockInstanaAPI) SliConfigs() rest.RestResource[*api.SliConfig] {
	return nil
}

// SloConfigs mock implementation
func (m *MockInstanaAPI) SloConfigs() rest.RestResource[*api.SloConfig] {
	return nil
}

// SloAlertConfig mock implementation
func (m *MockInstanaAPI) SloAlertConfigs() rest.RestResource[*api.SloAlertConfig] {
	return nil
}

// SloCorrectionConfig mock implementation
func (m *MockInstanaAPI) SloCorrectionConfigs() rest.RestResource[*api.SloCorrectionConfig] {
	return nil
}

// WebsiteMonitoringConfig mock implementation
func (m *MockInstanaAPI) WebsiteMonitoringConfigs() rest.RestResource[*api.WebsiteMonitoringConfig] {
	return nil
}

// WebsiteAlertConfig mock implementation
func (m *MockInstanaAPI) WebsiteAlertConfigs() rest.RestResource[*api.WebsiteAlertConfig] {
	return nil
}

// InfraAlertConfig mock implementation
func (m *MockInstanaAPI) InfraAlertConfigs() rest.RestResource[*api.InfraAlertConfig] {
	return nil
}

// Teams mock implementation
func (m *MockInstanaAPI) Teams() rest.RestResource[*api.Team] {
	return nil
}

// Groups mock implementation
func (m *MockInstanaAPI) Groups() rest.RestResource[*api.Group] {
	return nil
}

// Roles mock implementation
func (m *MockInstanaAPI) Roles() rest.RestResource[*api.Role] {
	return nil
}

// CustomDashboards mock implementation
func (m *MockInstanaAPI) CustomDashboards() rest.RestResource[*api.CustomDashboard] {
	return nil
}

// SyntheticTest mock implementation
func (m *MockInstanaAPI) SyntheticTests() rest.RestResource[*api.SyntheticTest] {
	return nil
}

// SyntheticLocation mock implementation
func (m *MockInstanaAPI) SyntheticLocations() rest.ReadOnlyRestResource[*api.SyntheticLocation] {
	return nil
}

// SyntheticAlertConfigs mock implementation
func (m *MockInstanaAPI) SyntheticAlertConfigs() rest.RestResource[*api.SyntheticAlertConfig] {
	return nil
}

// AutomationActions mock implementation
func (m *MockInstanaAPI) AutomationActions() rest.RestResource[*api.AutomationAction] {
	return nil
}

// AutomationPolicies mock implementation
func (m *MockInstanaAPI) AutomationPolicies() rest.RestResource[*api.AutomationPolicy] {
	return nil
}

// HostAgents mock implementation
func (m *MockInstanaAPI) HostAgents() rest.ReadOnlyRestResource[*api.HostAgent] {
	return nil
}

// Users mock implementation
func (m *MockInstanaAPI) Users() rest.ReadOnlyRestResource[*api.User] {
	return nil
}

// LogAlertConfig mock implementation
func (m *MockInstanaAPI) LogAlertConfigs() rest.RestResource[*api.LogAlertConfig] {
	return nil
}

// MobileAppConfig mock implementation
func (m *MockInstanaAPI) MobileAppConfig() rest.RestResource[*api.MobileAppConfig] {
	return nil
}

// MobileAlertConfig mock implementation
func (m *MockInstanaAPI) MobileAlertConfigs() rest.RestResource[*api.MobileAlertConfig] {
	return nil
}

// MaintenanceWindowConfigs mock implementation
func (m *MockInstanaAPI) MaintenanceWindowConfigs() rest.RestResource[*api.MaintenanceWindow] {
	return nil
}
