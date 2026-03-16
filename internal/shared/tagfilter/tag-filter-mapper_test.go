package tagfilter_test

import (
	"testing"

	"github.com/instana/instana-go-client/instana"
	"github.com/instana/instana-go-client/shared/tagfilter"
	. "github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnNilWhenMappingEmptyTagFilterExpressionToNormalizedString(t *testing.T) {
	operatorType := instana.LogicalOr
	input := &instana.TagFilter{
		Type:            instana.TagFilterExpressionType,
		LogicalOperator: &operatorType,
	}

	result, err := MapTagFilterToNormalizedString(input)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestShouldReturnErrorWhenMappingAnInvalidTagFilterExpressionToNormalizedString(t *testing.T) {
	input := &instana.TagFilter{
		Type: instana.TagFilterExpressionElementType("invalid"),
	}

	result, err := MapTagFilterToNormalizedString(input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestShouldReturnStringWhenMappingAValidTagFilterExpressionToNormalizedString(t *testing.T) {
	value := int64(1234)
	input := instana.NewNumberTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, instana.EqualsOperator, value)

	result, err := MapTagFilterToNormalizedString(input)

	require.NoError(t, err)
	require.Equal(t, "name@dest EQUALS 1234", *result)
}
