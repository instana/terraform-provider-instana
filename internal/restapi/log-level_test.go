package restapi_test

import (
	"testing"

	. "github.com/instana/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedLogLevelsAsStringSlice(t *testing.T) {
	expected := []string{"WARN", "ERROR", "ANY"}
	require.Equal(t, expected, SupportedLogLevels.ToStringSlice())
}
