package tagfilter

import (
	"fmt"

	tag "github.com/instana/instana-go-client/shared/tagfilter"
	"github.com/instana/instana-go-client/shared/types"
)

// ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *tagFilterMapper) ToAPIModel(input *FilterExpression) *tag.TagFilter {
	return m.mapLogicalOrToAPIModel(input.Expression)
}

func (m *tagFilterMapper) mapLogicalOrToAPIModel(input *LogicalOrExpression) *tag.TagFilter {
	left := m.mapLogicalAndToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalOrToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, types.LogicalOr)
		rightElements := m.unwrapExpressionElements(right, types.LogicalOr)
		return tag.NewLogicalOrTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) mapLogicalAndToAPIModel(input *LogicalAndExpression) *tag.TagFilter {
	left := m.mapBracketExpressionToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalAndToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, types.LogicalAnd)
		rightElements := m.unwrapExpressionElements(right, types.LogicalAnd)
		return tag.NewLogicalAndTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) unwrapExpressionElements(element *tag.TagFilter, operator types.LogicalOperatorType) []*tag.TagFilter {
	if element.GetType() == tag.TagFilterExpressionType && element.LogicalOperator != nil && *element.LogicalOperator == operator {
		return element.Elements
	}
	return []*tag.TagFilter{element}
}

func (m *tagFilterMapper) mapBracketExpressionToAPIModel(input *BracketExpression) *tag.TagFilter {
	if input.Bracket != nil {
		return m.mapLogicalOrToAPIModel(input.Bracket)
	}
	return m.mapPrimaryExpressionToAPIModel(input.Primary)
}

func (m *tagFilterMapper) mapPrimaryExpressionToAPIModel(input *PrimaryExpression) *tag.TagFilter {
	if input.UnaryOperation != nil {
		return m.mapUnaryOperatorExpressionToAPIModel(input.UnaryOperation)
	}
	return m.mapComparisonExpressionToAPIModel(input.Comparison)
}

func (m *tagFilterMapper) mapUnaryOperatorExpressionToAPIModel(input *UnaryOperationExpression) *tag.TagFilter {
	origin := EntityOriginDestination.TagFilterEntity()
	if input.Entity.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*input.Entity.Origin).TagFilterEntity()
	}
	return tag.NewUnaryTagFilterWithTagKey(origin, input.Entity.Identifier, input.Entity.TagKey, types.ExpressionOperator(input.Operator))
}

func (m *tagFilterMapper) mapComparisonExpressionToAPIModel(input *ComparisonExpression) *tag.TagFilter {
	origin := EntityOriginDestination.TagFilterEntity()
	if input.Entity.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*input.Entity.Origin).TagFilterEntity()
	}
	if input.Entity.TagKey != nil {
		return tag.NewTagTagFilter(origin, input.Entity.Identifier, types.ExpressionOperator(input.Operator), *input.Entity.TagKey, m.mapValueAsString(input))
	} else if input.NumberValue != nil {
		return tag.NewNumberTagFilter(origin, input.Entity.Identifier, types.ExpressionOperator(input.Operator), *input.NumberValue)
	} else if input.BooleanValue != nil {
		return tag.NewBooleanTagFilter(origin, input.Entity.Identifier, types.ExpressionOperator(input.Operator), *input.BooleanValue)
	}
	return tag.NewStringTagFilter(origin, input.Entity.Identifier, types.ExpressionOperator(input.Operator), *input.StringValue)
}

func (m *tagFilterMapper) mapValueAsString(input *ComparisonExpression) string {
	if input.NumberValue != nil {
		return fmt.Sprintf("%d", *input.NumberValue)
	} else if input.BooleanValue != nil {
		return fmt.Sprintf("%t", *input.BooleanValue)
	}
	return *input.StringValue
}
