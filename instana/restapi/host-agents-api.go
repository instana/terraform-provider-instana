package restapi

// HostAgent is the representation of a host agent in Instana
type HostAgent struct {
	SnapshotID string   `json:"snapshotId"`
	Label      string   `json:"label"`
	Host       string   `json:"host"`
	Plugin     string   `json:"plugin"`
	Tags       []string `json:"tags"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *HostAgent) GetIDForResourcePath() string {
	return spec.SnapshotID
}
