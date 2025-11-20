package restapi_test

import (
	"testing"

	. "github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	groupId                = "group-id"
	groupName              = "group-name"
	groupPermissionScopeID = "group-permission-scope-id"
)

func TestShouldReturnSupportedInstanaPermissionsAsString(t *testing.T) {
	expectedResult := []string{
		"CAN_CONFIGURE_APPLICATIONS",
		"CAN_CONFIGURE_EUM_APPLICATIONS",
		"CAN_CONFIGURE_AGENTS",
		"CAN_VIEW_TRACE_DETAILS",
		"CAN_VIEW_LOGS",
		"CAN_CONFIGURE_SESSION_SETTINGS",
		"CAN_CONFIGURE_INTEGRATIONS",
		"CAN_CONFIGURE_GLOBAL_APPLICATION_SMART_ALERTS",
		"CAN_CONFIGURE_GLOBAL_SYNTHETIC_SMART_ALERTS",
		"CAN_CONFIGURE_GLOBAL_INFRA_SMART_ALERTS",
		"CAN_CONFIGURE_GLOBAL_LOG_SMART_ALERTS",
		"CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD",
		"CAN_CONFIGURE_MOBILE_APP_MONITORING",
		"CAN_CONFIGURE_API_TOKENS",
		"CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS",
		"CAN_CONFIGURE_AUTHENTICATION_METHODS",
		"CAN_CONFIGURE_RELEASES",
		"CAN_VIEW_AUDIT_LOG",
		"CAN_CONFIGURE_EVENTS_AND_ALERTS",
		"CAN_CONFIGURE_MAINTENANCE_WINDOWS",
		"CAN_CONFIGURE_APPLICATION_SMART_ALERTS",
		"CAN_CONFIGURE_WEBSITE_SMART_ALERTS",
		"CAN_CONFIGURE_MOBILE_APP_SMART_ALERTS",
		"CAN_CONFIGURE_AGENT_RUN_MODE",
		"CAN_CONFIGURE_SERVICE_MAPPING",
		"CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS",
		"CAN_CONFIGURE_USERS",
		"CAN_INSTALL_NEW_AGENTS",
		"CAN_CONFIGURE_TEAMS",
		"CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS",
		"CAN_CONFIGURE_LOG_MANAGEMENT",
		"CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION",
		"CAN_VIEW_SYNTHETIC_TESTS",
		"CAN_VIEW_SYNTHETIC_LOCATIONS",
		"CAN_CREATE_THREAD_DUMP",
		"CAN_CREATE_HEAP_DUMP",
		"CAN_CONFIGURE_DATABASE_MANAGEMENT",
		"CAN_CONFIGURE_LOG_RETENTION_PERIOD",
		"CAN_CONFIGURE_PERSONAL_API_TOKENS",
		"ACCESS_INFRASTRUCTURE_ANALYZE",
		"CAN_VIEW_LOG_VOLUME",
		"CAN_RUN_AUTOMATION_ACTIONS",
		"CAN_VIEW_SYNTHETIC_TEST_RESULTS",
		"CAN_INVOKE_ALERT_CHANNEL",
		"CAN_MANUALLY_CLOSE_ISSUE",
		"CAN_DELETE_LOGS",
		"CAN_CONFIGURE_SYNTHETIC_TESTS",
		"CAN_VIEW_BUSINESS_PROCESS_DETAILS",
		"CAN_VIEW_BIZOPS_ALERTS",
		"CAN_USE_SYNTHETIC_CREDENTIALS",
		"CAN_DELETE_AUTOMATION_ACTION_HISTORY",
		"CAN_CONFIGURE_SYNTHETIC_LOCATIONS",
		"CAN_CONFIGURE_SYNTHETIC_CREDENTIALS",
		"CAN_CONFIGURE_SUBTRACES",
		"CAN_CONFIGURE_LLM",
		"CAN_CONFIGURE_BIZOPS",
		"CAN_CONFIGURE_AUTOMATION_POLICIES",
		"CAN_CONFIGURE_AUTOMATION_ACTIONS",
		"LIMITED_APPLICATIONS_SCOPE",
		"LIMITED_LINUX_KVM_HYPERVISOR_SCOPE",
		"LIMITED_VSPHERE_SCOPE",
		"LIMITED_OPENSTACK_SCOPE",
		"LIMITED_ZHMC_SCOPE",
		"LIMITED_XENSERVER_SCOPE",
		"LIMITED_KUBERNETES_SCOPE",
		"LIMITED_POWERVC_SCOPE",
		"LIMITED_SAP_SCOPE",
		"LIMITED_PCF_SCOPE",
		"LIMITED_SYNTHETICS_SCOPE",
		"LIMITED_SERVICE_LEVEL_SCOPE",
		"LIMITED_AUTOMATION_SCOPE",
		"LIMITED_BIZOPS_SCOPE",
		"LIMITED_PHMC_SCOPE",
		"LIMITED_GEN_AI_SCOPE",
		"LIMITED_INFRASTRUCTURE_SCOPE",
		"LIMITED_NUTANIX_SCOPE",
		"LIMITED_WINDOWS_HYPERVISOR_SCOPE",
		"LIMITED_AI_GATEWAY_SCOPE",
		"LIMITED_MOBILE_APPS_SCOPE",
		"LIMITED_WEBSITES_SCOPE",
		"CAN_CONFIGURE_APDEX",
		"CAN_CONFIGURE_CUSTOM_ENTITIES",
		"CAN_CONFIGURE_SERVICE_LEVELS",
		"CAN_CONFIGURE_SERVICE_LEVEL_CORRECTION_WINDOWS",
		"CAN_CONFIGURE_SERVICE_LEVEL_SMART_ALERTS",
		"ACCESS_APPLICATIONS",
		"ACCESS_MOBILE_APPS",
		"ACCESS_SYNTHETICS",
		"ACCESS_WEBSITES",
	}
	assert.Equal(t, expectedResult, SupportedInstanaPermissions.ToStringSlice())
}

func TestShouldReturnIDOfGroupAsIDForAPIPaths(t *testing.T) {
	group := Group{
		ID:   groupId,
		Name: groupName,
	}

	assert.Equal(t, groupId, group.GetIDForResourcePath())
}

func TestShouldReturnTrueWhenPermissionSetIsEmpty(t *testing.T) {
	p := APIPermissionSetWithRoles{}

	require.True(t, p.IsEmpty())

	emptyScopeBinding := ScopeBinding{}
	p = APIPermissionSetWithRoles{
		InfraDFQFilter: &emptyScopeBinding,
	}

	require.True(t, p.IsEmpty())

	emptyScopeBindingSlice := make([]ScopeBinding, 0)
	p = APIPermissionSetWithRoles{
		ApplicationIDs:          emptyScopeBindingSlice,
		KubernetesNamespaceUIDs: emptyScopeBindingSlice,
		KubernetesClusterUUIDs:  emptyScopeBindingSlice,
		InfraDFQFilter:          &emptyScopeBinding,
		MobileAppIDs:            emptyScopeBindingSlice,
		WebsiteIDs:              emptyScopeBindingSlice,
		Permissions:             make([]InstanaPermission, 0),
	}
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenApplicationIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		ApplicationIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenKubernetesClusterUUIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		KubernetesClusterUUIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenKubernetesNamespaceUIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		KubernetesNamespaceUIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenMobileAppIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		MobileAppIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenWebsiteIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		WebsiteIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenPermissionsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		Permissions: []InstanaPermission{PermissionCanConfigureApplications},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenInfrastructureDFQIsSet(t *testing.T) {
	scopeBinding := ScopeBinding{ScopeID: groupPermissionScopeID}
	p := APIPermissionSetWithRoles{
		InfraDFQFilter: &scopeBinding,
	}
	require.False(t, p.IsEmpty())
}
