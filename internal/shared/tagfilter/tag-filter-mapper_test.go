package tagfilter_test

import (
	"testing"

	tag "github.com/instana/instana-go-client/shared/tagfilter"
	common "github.com/instana/instana-go-client/shared/types"
	. "github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnNilWhenMappingEmptyTagFilterExpressionToNormalizedString(t *testing.T) {
	operatorType := common.LogicalOr
	input := &tag.TagFilter{
		Type:            tag.TagFilterExpressionType,
		LogicalOperator: &operatorType,
	}

	result, err := MapTagFilterToNormalizedString(input)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestShouldReturnErrorWhenMappingAnInvalidTagFilterExpressionToNormalizedString(t *testing.T) {
	input := &tag.TagFilter{
		Type: tag.TagFilterExpressionElementType("invalid"),
	}

	result, err := MapTagFilterToNormalizedString(input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestShouldReturnStringWhenMappingAValidTagFilterExpressionToNormalizedString(t *testing.T) {
	value := int64(1234)
	input := tag.NewNumberTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.EqualsOperator, value)

	result, err := MapTagFilterToNormalizedString(input)

	require.NoError(t, err)
	require.Equal(t, "name@dest EQUALS 1234", *result)
}
