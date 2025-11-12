package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedAccessTypesAsStringSlice(t *testing.T) {
	expected := []string{"READ", "READ_WRITE"}
	require.Equal(t, expected, SupportedAccessTypes.ToStringSlice())
}
