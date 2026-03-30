package tagfilter_test

import (
	"fmt"
	"testing"

	"github.com/instana/terraform-provider-instana/utils"

	"github.com/instana/instana-go-client/shared/tagfilter"
	tag "github.com/instana/instana-go-client/shared/tagfilter"
	common "github.com/instana/instana-go-client/shared/types"
	. "github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
	"github.com/stretchr/testify/require"
)

const (
	entitySpecKey = "key"
)

func TestShouldMapComparisonToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range common.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test comparison of string value using operatore %s", v), createTestShouldMapStringComparisonToRepresentationOfInstanaAPI(v))
		t.Run(fmt.Sprintf("test comparison of number value using operatore of %s", v), createTestShouldMapNumberComparisonToRepresentationOfInstanaAPI(v))
		t.Run(fmt.Sprintf("test comparison of boolean value using operatore of %s", v), createTestShouldMapBooleanComparisonToRepresentationOfInstanaAPI(v))
		t.Run(fmt.Sprintf("test comparison of tag using operatore of %s", v), createTestShouldMapTagComparisonToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapStringComparisonToRepresentationOfInstanaAPI(operator common.ExpressionOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(operator),
								StringValue: utils.StringPtr("value"),
							},
						},
					},
				},
			},
		}

		expectedResult := tag.NewStringTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, operator, "value")
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func createTestShouldMapNumberComparisonToRepresentationOfInstanaAPI(operator common.ExpressionOperator) func(*testing.T) {
	numberValue := int64(1234)
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(operator),
								NumberValue: &numberValue,
							},
						},
					},
				},
			},
		}

		expectedResult := tag.NewNumberTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, operator, numberValue)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func createTestShouldMapBooleanComparisonToRepresentationOfInstanaAPI(operator common.ExpressionOperator) func(*testing.T) {
	boolValue := true
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:       &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:     Operator(operator),
								BooleanValue: &boolValue,
							},
						},
					},
				},
			},
		}

		expectedResult := tag.NewBooleanTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, operator, boolValue)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func createTestShouldMapTagComparisonToRepresentationOfInstanaAPI(operator common.ExpressionOperator) func(*testing.T) {
	key := "key"
	value := "value"
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: entitySpecKey, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(operator),
								StringValue: &value,
							},
						},
					},
				},
			},
		}

		expectedResult := tag.NewTagTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, operator, key, value)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapTagComparisonToRepresentationOfInstanaAPIUsingAStringValue(t *testing.T) {
	key := "key"
	value := "value"
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: entitySpecKey, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(common.EqualsOperator),
							StringValue: &value,
						},
					},
				},
			},
		},
	}

	expectedResult := tag.NewTagTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.EqualsOperator, key, value)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapTagComparisonToRepresentationOfInstanaAPIUsingANumberValue(t *testing.T) {
	key := "key"
	value := int64(1234)
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: entitySpecKey, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(common.EqualsOperator),
							NumberValue: &value,
						},
					},
				},
			},
		},
	}

	expectedResult := tag.NewTagTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.EqualsOperator, key, "1234")
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapTagComparisonToRepresentationOfInstanaAPIUsingABooleanValue(t *testing.T) {
	key := "key"
	value := true
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:       &EntitySpec{Identifier: entitySpecKey, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:     Operator(common.EqualsOperator),
							BooleanValue: &value,
						},
					},
				},
			},
		},
	}

	expectedResult := tag.NewTagTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.EqualsOperator, key, "true")
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range common.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(operatorName common.ExpressionOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator: Operator(operatorName),
							},
						},
					},
				},
			},
		}

		expectedResult := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, operatorName)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapLogicalAndExpression(t *testing.T) {
	logicalAnd := Operator(common.LogicalAnd)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
			Operator: Operator(common.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left:     &BracketExpression{Primary: &primaryExpression},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{Primary: &primaryExpression},
				},
			},
		},
	}

	expectedPrimaryExpression := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.IsEmptyOperator)
	expectedResult := tag.NewLogicalAndTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedPrimaryExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalAndExpressionWithNestedAnd(t *testing.T) {
	logicalAnd := Operator(common.LogicalAnd)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
			Operator: Operator(common.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left:     &BracketExpression{Primary: &primaryExpression},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left:     &BracketExpression{Primary: &primaryExpression},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &BracketExpression{Primary: &primaryExpression},
					},
				},
			},
		},
	}

	expectedPrimaryExpression := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.IsEmptyOperator)
	expectedResult := tag.NewLogicalAndTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedPrimaryExpression, expectedPrimaryExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalAndExpressionWithNestedOrInBrackets(t *testing.T) {
	logicalAnd := Operator(common.LogicalAnd)
	logicalOr := Operator(common.LogicalOr)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
			Operator: Operator(common.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left:     &BracketExpression{Primary: &primaryExpression},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{
						Bracket: &LogicalOrExpression{
							Left:     &LogicalAndExpression{Left: &BracketExpression{Primary: &primaryExpression}},
							Operator: &logicalOr,
							Right: &LogicalOrExpression{
								Left: &LogicalAndExpression{Left: &BracketExpression{Primary: &primaryExpression}},
							},
						},
					},
				},
			},
		},
	}

	expectedPrimaryExpression := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.IsEmptyOperator)
	expectedOrExpression := tag.NewLogicalOrTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedPrimaryExpression})
	expectedResult := tag.NewLogicalAndTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedOrExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalOrExpression(t *testing.T) {
	logicalOr := Operator(common.LogicalOr)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
			Operator: Operator(common.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{Primary: &primaryExpression},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{Primary: &primaryExpression},
				},
			},
		},
	}

	expectedPrimaryExpression := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.IsEmptyOperator)
	expectedResult := tag.NewLogicalOrTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedPrimaryExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalOrExpressionWithNestedOr(t *testing.T) {
	logicalOr := Operator(common.LogicalOr)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
			Operator: Operator(common.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{Primary: &primaryExpression},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{Primary: &primaryExpression},
				},
				Operator: &logicalOr,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &BracketExpression{Primary: &primaryExpression},
					},
				},
			},
		},
	}

	expectedPrimaryExpression := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.IsEmptyOperator)
	expectedResult := tag.NewLogicalOrTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedPrimaryExpression, expectedPrimaryExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalOrExpressionWithNestedAndInBrackets(t *testing.T) {
	logicalOr := Operator(common.LogicalOr)
	logicalAnd := Operator(common.LogicalAnd)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: utils.StringPtr(EntityOriginDestination.Key())},
			Operator: Operator(common.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{Primary: &primaryExpression},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Bracket: &LogicalOrExpression{
							Left: &LogicalAndExpression{
								Left:     &BracketExpression{Primary: &primaryExpression},
								Operator: &logicalAnd,
								Right: &LogicalAndExpression{
									Left: &BracketExpression{Primary: &primaryExpression},
								},
							},
						},
					},
				},
			},
		},
	}

	expectedPrimaryExpression := tag.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, entitySpecKey, common.IsEmptyOperator)
	expectedAndExpression := tag.NewLogicalAndTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedPrimaryExpression})
	expectedResult := tag.NewLogicalOrTagFilter([]*tag.TagFilter{expectedPrimaryExpression, expectedAndExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func runTestCaseForMappingToAPI(input *FilterExpression, expectedResult *tag.TagFilter, t *testing.T) {
	mapper := NewMapper()
	result := mapper.ToAPIModel(input)

	require.Equal(t, expectedResult, result)
}
