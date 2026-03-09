package tagfilter

import (
	"fmt"

	"github.com/instana/instana-go-client/instana"
)

// ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *tagFilterMapper) ToAPIModel(input *FilterExpression) *instana.TagFilter {
	return m.mapLogicalOrToAPIModel(input.Expression)
}

func (m *tagFilterMapper) mapLogicalOrToAPIModel(input *LogicalOrExpression) *instana.TagFilter {
	left := m.mapLogicalAndToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalOrToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, instana.LogicalOr)
		rightElements := m.unwrapExpressionElements(right, instana.LogicalOr)
		return instana.NewLogicalOrTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) mapLogicalAndToAPIModel(input *LogicalAndExpression) *instana.TagFilter {
	left := m.mapBracketExpressionToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalAndToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, instana.LogicalAnd)
		rightElements := m.unwrapExpressionElements(right, instana.LogicalAnd)
		return instana.NewLogicalAndTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) unwrapExpressionElements(element *instana.TagFilter, operator instana.LogicalOperatorType) []*instana.TagFilter {
	if element.GetType() == instana.TagFilterExpressionType && element.LogicalOperator != nil && *element.LogicalOperator == operator {
		return element.Elements
	}
	return []*instana.TagFilter{element}
}

func (m *tagFilterMapper) mapBracketExpressionToAPIModel(input *BracketExpression) *instana.TagFilter {
	if input.Bracket != nil {
		return m.mapLogicalOrToAPIModel(input.Bracket)
	}
	return m.mapPrimaryExpressionToAPIModel(input.Primary)
}

func (m *tagFilterMapper) mapPrimaryExpressionToAPIModel(input *PrimaryExpression) *instana.TagFilter {
	if input.UnaryOperation != nil {
		return m.mapUnaryOperatorExpressionToAPIModel(input.UnaryOperation)
	}
	return m.mapComparisonExpressionToAPIModel(input.Comparison)
}

func (m *tagFilterMapper) mapUnaryOperatorExpressionToAPIModel(input *UnaryOperationExpression) *instana.TagFilter {
	origin := EntityOriginDestination.TagFilterEntity()
	if input.Entity.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*input.Entity.Origin).TagFilterEntity()
	}
	return instana.NewUnaryTagFilterWithTagKey(origin, input.Entity.Identifier, input.Entity.TagKey, instana.ExpressionOperator(input.Operator))
}

func (m *tagFilterMapper) mapComparisonExpressionToAPIModel(input *ComparisonExpression) *instana.TagFilter {
	origin := EntityOriginDestination.TagFilterEntity()
	if input.Entity.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*input.Entity.Origin).TagFilterEntity()
	}
	if input.Entity.TagKey != nil {
		return instana.NewTagTagFilter(origin, input.Entity.Identifier, instana.ExpressionOperator(input.Operator), *input.Entity.TagKey, m.mapValueAsString(input))
	} else if input.NumberValue != nil {
		return instana.NewNumberTagFilter(origin, input.Entity.Identifier, instana.ExpressionOperator(input.Operator), *input.NumberValue)
	} else if input.BooleanValue != nil {
		return instana.NewBooleanTagFilter(origin, input.Entity.Identifier, instana.ExpressionOperator(input.Operator), *input.BooleanValue)
	}
	return instana.NewStringTagFilter(origin, input.Entity.Identifier, instana.ExpressionOperator(input.Operator), *input.StringValue)
}

func (m *tagFilterMapper) mapValueAsString(input *ComparisonExpression) string {
	if input.NumberValue != nil {
		return fmt.Sprintf("%d", *input.NumberValue)
	} else if input.BooleanValue != nil {
		return fmt.Sprintf("%t", *input.BooleanValue)
	}
	return *input.StringValue
}
