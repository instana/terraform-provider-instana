package restapi

// MobileAppConfigResourcePath path to mobile app monitoring config resource of Instana RESTful API
const MobileAppConfigResourcePath = InstanaAPIBasePath + "/mobile-app-monitoring/config"

// MobileAppConfig data structure of a Mobile App Configuration of the Instana API
type MobileAppConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (r *MobileAppConfig) GetIDForResourcePath() string {
	return r.ID
}

