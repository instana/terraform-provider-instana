package restapi

// TeamsResourcePath path to Team resource of Instana RESTful API
const TeamsResourcePath = RBACSettingsBasePath + "/teams"

// TeamInfo data structure for the Instana API model for team information
type TeamInfo struct {
	Description *string `json:"description,omitempty"`
}

// TeamRole data structure for the Instana API model for team member roles
type TeamRole struct {
	RoleID   string  `json:"roleId"`
	RoleName *string `json:"roleName,omitempty"`
	ViaIdP   *bool   `json:"viaIdP,omitempty"`
}

// TeamMember data structure for the Instana API model for team members
type TeamMember struct {
	UserID string     `json:"userId"`
	Email  *string    `json:"email,omitempty"`
	Name   *string    `json:"name,omitempty"`
	Roles  []TeamRole `json:"roles,omitempty"`
}

// RestrictedApplicationFilterScope type for the scope of restricted application filter
type RestrictedApplicationFilterScope string

const (
	// RestrictedApplicationFilterScopeIncludeNoDownstream includes no downstream services
	RestrictedApplicationFilterScopeIncludeNoDownstream = RestrictedApplicationFilterScope("INCLUDE_NO_DOWNSTREAM")
	// RestrictedApplicationFilterScopeIncludeImmediateDownstream includes immediate downstream database and messaging
	RestrictedApplicationFilterScopeIncludeImmediateDownstream = RestrictedApplicationFilterScope("INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING")
	// RestrictedApplicationFilterScopeIncludeAllDownstream includes all downstream services
	RestrictedApplicationFilterScopeIncludeAllDownstream = RestrictedApplicationFilterScope("INCLUDE_ALL_DOWNSTREAM")
)

// RestrictedApplicationFilter data structure for the Instana API model for restricted application filter
type RestrictedApplicationFilter struct {
	Label                    *string                           `json:"label,omitempty"`
	RestrictingApplicationID *string                           `json:"restrictingApplicationId,omitempty"`
	Scope                    *RestrictedApplicationFilterScope `json:"scope,omitempty"`
	TagFilterExpression      *TagFilter                        `json:"tagFilterExpression,omitempty"`
}

// TeamScope data structure for the Instana API model for team scope
type TeamScope struct {
	AccessPermissions           []string                     `json:"accessPermissions,omitempty"`
	Applications                []string                     `json:"applications,omitempty"`
	KubernetesClusters          []string                     `json:"kubernetesClusters,omitempty"`
	KubernetesNamespaces        []string                     `json:"kubernetesNamespaces,omitempty"`
	MobileApps                  []string                     `json:"mobileApps,omitempty"`
	Websites                    []string                     `json:"websites,omitempty"`
	InfraDFQFilter              *string                      `json:"infraDfqFilter,omitempty"`
	ActionFilter                *string                      `json:"actionFilter,omitempty"`
	LogFilter                   *string                      `json:"logFilter,omitempty"`
	BusinessPerspectives        []string                     `json:"businessPerspectives,omitempty"`
	SloIDs                      []string                     `json:"sloIds,omitempty"`
	SyntheticTests              []string                     `json:"syntheticTests,omitempty"`
	SyntheticCredentials        []string                     `json:"syntheticCredentials,omitempty"`
	TagIDs                      []string                     `json:"tagIds,omitempty"`
	RestrictedApplicationFilter *RestrictedApplicationFilter `json:"restrictedApplicationFilter,omitempty"`
}

// Team data structure for the Instana API model for teams
type Team struct {
	ID      string       `json:"id,omitempty"`
	Tag     string       `json:"tag"`
	Info    *TeamInfo    `json:"info,omitempty"`
	Members []TeamMember `json:"members,omitempty"`
	Scope   *TeamScope   `json:"scope,omitempty"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (t *Team) GetIDForResourcePath() string {
	return t.ID
}

// Made with Bob
