package shared

import "github.com/instana/instana-go-client/client"

// ProviderMeta data structure for the metadata which is configured and provided to the resources by this provider
type ProviderMeta struct {
	InstanaAPI client.InstanaAPI
}
