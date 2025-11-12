package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/internal/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedGranularitiesAsIntSlice(t *testing.T) {
	expected := []int{60000, 300000, 600000, 900000, 1200000, 1800000}
	require.Equal(t, expected, SupportedSmartAlertGranularities.ToIntSlice())
}
