package instana

import (
	"testing"
)

func TestProviderFramework(t *testing.T) {
	provider := New("test")()
	if provider == nil {
		t.Fatal("provider is nil")
	}
}
