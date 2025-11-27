package restapi

// RolesResourcePath path to Role resource of Instana RESTful API
const RolesResourcePath = RBACSettingsBasePath + "/roles"

// Role data structure for the Instana API model for roles
type Role struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Members     []APIMember   `json:"members"`
	Permissions []string      `json:"permissions"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (r *Role) GetIDForResourcePath() string {
	return r.ID
}

// Made with Bob
