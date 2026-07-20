package provider

import (
	"sync"
	"testing"

	"github.com/instana/instana-go-client/config"
	"github.com/instana/terraform-provider-instana/internal/shared"
)

// TestAddCorrelationIDToClient_ConcurrentSafe reproduces the scenario in #91:
// many resource instances (as created by NewTerraformResource per
// Create/Read/Update/Delete call) sharing the same *shared.ProviderMeta and
// calling addCorrelationIDToClient concurrently, as happens under
// Terraform's default parallelism when a config fans out many resources of
// the same type. Before the fix, `go test -race` (or enough iterations/CPUs)
// reliably reports a `fatal error: concurrent map writes` panic or a data
// race on Headers.Custom; with the mutex in place, it passes cleanly.
func TestAddCorrelationIDToClient_ConcurrentSafe(t *testing.T) {
	providerMeta := &shared.ProviderMeta{
		ClientConfig: &config.ClientConfig{},
	}

	const goroutines = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(n int) {
			defer wg.Done()
			r := &terraformResourceImpl[fakeInstanaDataObject]{
				providerMeta: providerMeta,
			}
			r.addCorrelationIDToClient("correlation-id")
		}(i)
	}

	wg.Wait()

	if providerMeta.ClientConfig.Headers.Custom[CorrelationIDHeader] != "correlation-id" {
		t.Fatalf("expected correlation ID header to be set after concurrent writes")
	}
}

// fakeInstanaDataObject is a minimal stand-in satisfying the
// client.InstanaDataObject type parameter for terraformResourceImpl in this
// test; none of its methods are exercised by addCorrelationIDToClient.
type fakeInstanaDataObject struct{}

func (fakeInstanaDataObject) GetIDForResourcePath() string { return "" }
