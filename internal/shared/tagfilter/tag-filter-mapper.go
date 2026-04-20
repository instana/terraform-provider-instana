package tagfilter

import tag "github.com/instana/instana-go-client/shared/tagfilter"

// NewMapper creates a new instance of the Mapper
func NewMapper() Mapper {
	return &tagFilterMapper{}
}

// Mapper interface of the tag filter expression mapper
type Mapper interface {
	FromAPIModel(input *tag.TagFilter) (*FilterExpression, error)
	ToAPIModel(input *FilterExpression) *tag.TagFilter
}

// struct for the filter expression mapper implementation for tag filter expressions
type tagFilterMapper struct{}

// MapTagFilterToNormalizedString maps a TagFilterExpressionElement to its normalized string. Returns nil in case an empty expression is provided and an error in case of any error occurred during mapping.
func MapTagFilterToNormalizedString(element *tag.TagFilter) (*string, error) {
	mapper := NewMapper()
	expr, err := mapper.FromAPIModel(element)
	if err != nil {
		return nil, err
	}
	if expr != nil {
		renderedExpression := expr.Render()
		return &renderedExpression, nil
	}
	return nil, nil
}
