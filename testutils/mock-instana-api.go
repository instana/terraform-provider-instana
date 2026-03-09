package testutils

import "github.com/instana/instana-go-client/instana"

// MockInstanaAPI is a mock implementation of the InstanaAPI interface for testing purposes.
// It returns nil for all methods by default. Tests can override specific methods by embedding
// this struct and providing custom implementations for the methods they need.
type MockInstanaAPI struct{}

// CustomEventSpecifications mock implementation
func (m *MockInstanaAPI) CustomEventSpecifications() instana.RestResource[*instana.CustomEventSpecification] {
	return nil
}

// BuiltinEventSpecifications mock implementation
func (m *MockInstanaAPI) BuiltinEventSpecifications() instana.ReadOnlyRestResource[*instana.BuiltinEventSpecification] {
	return nil
}

// APITokens mock implementation
func (m *MockInstanaAPI) APITokens() instana.RestResource[*instana.APIToken] {
	return nil
}

// ApplicationConfigs mock implementation
func (m *MockInstanaAPI) ApplicationConfigs() instana.RestResource[*instana.ApplicationConfig] {
	return nil
}

// ApplicationAlertConfigs mock implementation
func (m *MockInstanaAPI) ApplicationAlertConfigs() instana.RestResource[*instana.ApplicationAlertConfig] {
	return nil
}

// GlobalApplicationAlertConfigs mock implementation
func (m *MockInstanaAPI) GlobalApplicationAlertConfigs() instana.RestResource[*instana.ApplicationAlertConfig] {
	return nil
}

// AlertingChannels mock implementation
func (m *MockInstanaAPI) AlertingChannels() instana.RestResource[*instana.AlertingChannel] {
	return nil
}

// AlertingConfigurations mock implementation
func (m *MockInstanaAPI) AlertingConfigurations() instana.RestResource[*instana.AlertingConfiguration] {
	return nil
}

// SliConfigs mock implementation
func (m *MockInstanaAPI) SliConfigs() instana.RestResource[*instana.SliConfig] {
	return nil
}

// SloConfigs mock implementation
func (m *MockInstanaAPI) SloConfigs() instana.RestResource[*instana.SloConfig] {
	return nil
}

// SloAlertConfig mock implementation
func (m *MockInstanaAPI) SloAlertConfig() instana.RestResource[*instana.SloAlertConfig] {
	return nil
}

// SloCorrectionConfig mock implementation
func (m *MockInstanaAPI) SloCorrectionConfig() instana.RestResource[*instana.SloCorrectionConfig] {
	return nil
}

// WebsiteMonitoringConfig mock implementation
func (m *MockInstanaAPI) WebsiteMonitoringConfig() instana.RestResource[*instana.WebsiteMonitoringConfig] {
	return nil
}

// WebsiteAlertConfig mock implementation
func (m *MockInstanaAPI) WebsiteAlertConfig() instana.RestResource[*instana.WebsiteAlertConfig] {
	return nil
}

// InfraAlertConfig mock implementation
func (m *MockInstanaAPI) InfraAlertConfig() instana.RestResource[*instana.InfraAlertConfig] {
	return nil
}

// Teams mock implementation
func (m *MockInstanaAPI) Teams() instana.RestResource[*instana.Team] {
	return nil
}

// Groups mock implementation
func (m *MockInstanaAPI) Groups() instana.RestResource[*instana.Group] {
	return nil
}

// Roles mock implementation
func (m *MockInstanaAPI) Roles() instana.RestResource[*instana.Role] {
	return nil
}

// CustomDashboards mock implementation
func (m *MockInstanaAPI) CustomDashboards() instana.RestResource[*instana.CustomDashboard] {
	return nil
}

// SyntheticTest mock implementation
func (m *MockInstanaAPI) SyntheticTest() instana.RestResource[*instana.SyntheticTest] {
	return nil
}

// SyntheticLocation mock implementation
func (m *MockInstanaAPI) SyntheticLocation() instana.ReadOnlyRestResource[*instana.SyntheticLocation] {
	return nil
}

// SyntheticAlertConfigs mock implementation
func (m *MockInstanaAPI) SyntheticAlertConfigs() instana.RestResource[*instana.SyntheticAlertConfig] {
	return nil
}

// AutomationActions mock implementation
func (m *MockInstanaAPI) AutomationActions() instana.RestResource[*instana.AutomationAction] {
	return nil
}

// AutomationPolicies mock implementation
func (m *MockInstanaAPI) AutomationPolicies() instana.RestResource[*instana.AutomationPolicy] {
	return nil
}

// HostAgents mock implementation
func (m *MockInstanaAPI) HostAgents() instana.ReadOnlyRestResource[*instana.HostAgent] {
	return nil
}

// Users mock implementation
func (m *MockInstanaAPI) Users() instana.ReadOnlyRestResource[*instana.User] {
	return nil
}

// LogAlertConfig mock implementation
func (m *MockInstanaAPI) LogAlertConfig() instana.RestResource[*instana.LogAlertConfig] {
	return nil
}

// MobileAppConfig mock implementation
func (m *MockInstanaAPI) MobileAppConfig() restapi.RestResource[*restapi.MobileAppConfig] {
	return nil
}

// MobileAlertConfig mock implementation
func (m *MockInstanaAPI) MobileAlertConfig() instana.RestResource[*instana.MobileAlertConfig] {
	return nil
}

// MaintenanceWindowConfigs mock implementation
func (m *MockInstanaAPI) MaintenanceWindowConfigs() instana.RestResource[*instana.MaintenanceWindowConfig] {
	return nil
}
