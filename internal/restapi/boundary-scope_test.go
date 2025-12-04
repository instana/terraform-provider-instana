package restapi_test

import (
	"testing"

	. "github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnStringRepresentationOfSupportedApplicationConfigBoundaryScopes(t *testing.T) {
	require.Equal(t, []string{"ALL", "INBOUND", "DEFAULT"}, SupportedApplicationConfigBoundaryScopes.ToStringSlice())
}

func TestShouldReturnStringRepresentationOfSupportedApplicationAlertConfigBoundaryScopes(t *testing.T) {
	require.Equal(t, []string{"ALL", "INBOUND"}, SupportedApplicationAlertConfigBoundaryScopes.ToStringSlice())
}
