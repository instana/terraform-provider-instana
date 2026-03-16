package tagfilter_test

import (
	"fmt"
	"testing"

	"github.com/instana/terraform-provider-instana/utils"

	"github.com/stretchr/testify/require"

	"github.com/instana/instana-go-client/shared/tagfilter"
	. "github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
)

const (
	invalidOperator   = "invalid operator"
	tagFilterOperator = "tag filter operator"
	tagFilterName     = "name"
)

func TestShouldMapEmptyTagFilterExpressionFromInstanaAPI(t *testing.T) {
	for _, operator := range typ.SupportedLogicalOperatorTypes {
		t.Run(fmt.Sprintf("TestShouldMapEmpty%sTagFilterExpressionFromInstnaAPI", string(operator)), func(t *testing.T) {
			op := operator
			expression := &tagfilter.TagFilter{
				Type:            tagfilter.TagFilterExpressionType,
				LogicalOperator: &op,
			}

			runTestCaseForMappingFromAPI(expression, nil, t)
		})
	}
}

func TestShouldMapStringTagFilterFromInstanaAPI(t *testing.T) {
	value := "value"
	input := tagfilter.NewStringTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(tagfilter.EqualsOperator),
		StringValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapNumberTagFilterFromInstanaAPI(t *testing.T) {
	value := int64(1234)
	input := tagfilter.NewNumberTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(tagfilter.EqualsOperator),
		NumberValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapBooleanTagFilterFromInstanaAPI(t *testing.T) {
	value := true
	input := tagfilter.NewBooleanTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:       &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:     Operator(tagfilter.EqualsOperator),
		BooleanValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapComparisonTagFilterWithTagKeyValueFromInstanaAPI(t *testing.T) {
	key := "key"
	value := "value"
	input := tagfilter.NewTagTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.EqualsOperator, key, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(tagfilter.EqualsOperator),
		StringValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func testMappingOfTagFilterFromInstanaApi(tagFilter *tagfilter.TagFilter, comparison *ComparisonExpression, t *testing.T) {
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{Comparison: comparison},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(tagFilter, expectedResult, t)
}

func TestShouldMapAllSupportedComparisonOperatorsFromInstanaAPI(t *testing.T) {
	for _, v := range tagfilter.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), testMappingOfSupportedComparisonOperatorsFromInstanaAPI(v))
	}
}

func testMappingOfSupportedComparisonOperatorsFromInstanaAPI(operator tagfilter.ExpressionOperator) func(t *testing.T) {
	return func(t *testing.T) {
		value := "value"
		input := tagfilter.NewStringTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, operator, value)

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(operator),
								StringValue: &value,
							},
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldFailToMapTagFilterFromInstanaAPIWhenOperatorIsNotSupported(t *testing.T) {
	value := "value"
	input := tagfilter.NewStringTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, "FOO", value)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldMapAllSupportedUnaryOperationsFromInstanaAPI(t *testing.T) {
	for _, v := range tagfilter.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfSupportedUnaryOperationFromInstanaAPI(v))
	}
}

func testMappingOfSupportedUnaryOperationFromInstanaAPI(operator tagfilter.ExpressionOperator) func(t *testing.T) {
	return func(t *testing.T) {
		input := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, operator)

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator: Operator(operator),
							},
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldMapUnaryTagFilterWithTagKeyFromInstanaAPI(t *testing.T) {
	key := "key"
	input := tagfilter.NewUnaryTagFilterWithTagKey(tagfilter.TagFilterEntityDestination, tagFilterName, &key, tagfilter.NotEmptyOperator)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(tagfilter.NotEmptyOperator),
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldFailToMapTagFilterFromInstanaAPIWhenUnaryOperationIsNotSupported(t *testing.T) {
	input := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, "FOO")

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldFailToMapTagFilterExpressionElementFromInstanaAPIWhenTypeIsMissing(t *testing.T) {
	name := tagFilterName
	operator := tagfilter.ExpressionOperator("FOO")
	input := &tagfilter.TagFilter{
		Name:     &name,
		Operator: &operator,
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported tag filter expression")
}

func TestShouldMapLogicalAndWithTwoPrimaryExpressionsFromInstanaAPI(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	and := Operator(tagfilter.LogicalAnd)
	primaryExpression1 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name1", tagfilter.IsEmptyOperator)
	primaryExpression2 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name2", tagfilter.IsEmptyOperator)
	input := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression1, primaryExpression2})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									UnaryOperation: &UnaryOperationExpression{
										Entity:   &EntitySpec{Identifier: "name1", Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator: operator,
									},
								},
							},
							Operator: &and,
							Right: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										UnaryOperation: &UnaryOperationExpression{
											Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
											Operator: operator,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalAndWithThreePrimaryExpressionsFromInstanaAPI(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	and := Operator(tagfilter.LogicalAnd)
	primaryExpression1 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name1", tagfilter.IsEmptyOperator)
	primaryExpression2 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name2", tagfilter.IsEmptyOperator)
	primaryExpression3 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name3", tagfilter.IsEmptyOperator)
	input := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression1, primaryExpression2, primaryExpression3})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									UnaryOperation: &UnaryOperationExpression{
										Entity:   &EntitySpec{Identifier: "name1", Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator: operator,
									},
								},
							},
							Operator: &and,
							Right: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										UnaryOperation: &UnaryOperationExpression{
											Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
											Operator: operator,
										},
									},
								},
								Operator: &and,
								Right: &LogicalAndExpression{
									Left: &BracketExpression{
										Primary: &PrimaryExpression{
											UnaryOperation: &UnaryOperationExpression{
												Entity:   &EntitySpec{Identifier: "name3", Origin: utils.StringPtr(EntityOriginDestination.Key())},
												Operator: operator,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalAndWithTwoElementsFromInstanaAPIWhereTheFirstElementIsAPrimaryExpressionAndTheSecondElementIsAnotherLogicalAnd(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	and := Operator(tagfilter.LogicalAnd)
	primaryExpression1 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	primaryExpression2 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name2", tagfilter.IsEmptyOperator)
	nestedAnd := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression2, primaryExpression2})
	input := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression1, nestedAnd})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									UnaryOperation: &UnaryOperationExpression{
										Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator: operator,
									},
								},
							},
							Operator: &and,
							Right: &LogicalAndExpression{
								Left: &BracketExpression{
									Bracket: &LogicalOrExpression{
										Left: &LogicalAndExpression{
											Left: &BracketExpression{
												Primary: &PrimaryExpression{
													UnaryOperation: &UnaryOperationExpression{
														Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
														Operator: operator,
													},
												},
											},
											Operator: &and,
											Right: &LogicalAndExpression{
												Left: &BracketExpression{
													Primary: &PrimaryExpression{
														UnaryOperation: &UnaryOperationExpression{
															Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
															Operator: operator,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalAndWithTwoElementsFromInstanaAPIWhereTheFirstElementIsAPrimaryExpressionAndTheSecondElementIsALogicalOr(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	and := Operator(tagfilter.LogicalAnd)
	or := Operator(tagfilter.LogicalOr)
	primaryExpression1 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	primaryExpression2 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name2", tagfilter.IsEmptyOperator)
	nestedOr := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression2, primaryExpression2})
	input := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression1, nestedOr})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									UnaryOperation: &UnaryOperationExpression{
										Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator: operator,
									},
								},
							},
							Operator: &and,
							Right: &LogicalAndExpression{
								Left: &BracketExpression{
									Bracket: &LogicalOrExpression{
										Left: &LogicalAndExpression{
											Left: &BracketExpression{
												Primary: &PrimaryExpression{
													UnaryOperation: &UnaryOperationExpression{
														Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
														Operator: operator,
													},
												},
											},
										},
										Operator: &or,
										Right: &LogicalOrExpression{
											Left: &LogicalAndExpression{
												Left: &BracketExpression{
													Primary: &PrimaryExpression{
														UnaryOperation: &UnaryOperationExpression{
															Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
															Operator: operator,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldUnwrapLogicalAndFromInstanaAPIWhenOnlyOneElementIsProvided(t *testing.T) {
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	input := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(tagfilter.IsEmptyOperator),
						},
					},
				},
			},
		},
	}

	mapper := NewMapper()
	result, err := mapper.FromAPIModel(input)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestShouldMapLogicalOrWithTwoPrimaryExpressionsFromInstanaAPI(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	or := Operator(tagfilter.LogicalOr)
	primaryExpression1 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name1", tagfilter.IsEmptyOperator)
	primaryExpression2 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name2", tagfilter.IsEmptyOperator)
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression1, primaryExpression2})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									UnaryOperation: &UnaryOperationExpression{
										Entity:   &EntitySpec{Identifier: "name1", Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator: operator,
									},
								},
							},
						},
						Operator: &or,
						Right: &LogicalOrExpression{
							Left: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										UnaryOperation: &UnaryOperationExpression{
											Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
											Operator: operator,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalOrWithThreePrimaryExpressionsFromInstanaAPI(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	or := Operator(tagfilter.LogicalOr)
	primaryExpression1 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name1", tagfilter.IsEmptyOperator)
	primaryExpression2 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name2", tagfilter.IsEmptyOperator)
	primaryExpression3 := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, "name3", tagfilter.IsEmptyOperator)
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression1, primaryExpression2, primaryExpression3})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{Left: &LogicalAndExpression{Left: &BracketExpression{Bracket: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: "name1", Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: operator,
						},
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: "name2", Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator: operator,
							},
						},
					},
				},
				Operator: &or,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &BracketExpression{
							Primary: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Entity:   &EntitySpec{Identifier: "name3", Origin: utils.StringPtr(EntityOriginDestination.Key())},
									Operator: operator,
								},
							},
						},
					},
				},
			},
		}}}},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalOrWithTwoElementsFromInstanaAPIWhereFirstElementIsALogicalAndAndTheOtherElementIsPrimaryExpression(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	or := Operator(tagfilter.LogicalOr)
	and := Operator(tagfilter.LogicalAnd)
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	nestedAnd := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression, primaryExpression})
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{nestedAnd, primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{Left: &LogicalAndExpression{Left: &BracketExpression{Bracket: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									UnaryOperation: &UnaryOperationExpression{
										Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator: operator,
									},
								},
							},
							Operator: &and,
							Right: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										UnaryOperation: &UnaryOperationExpression{
											Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
											Operator: operator,
										},
									},
								},
							},
						},
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator: operator,
							},
						},
					},
				},
			},
		}}}},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalOrWithTwoElementsFromInstanaAPIWhereFirstElementIsAPrimaryExpressionAndTheOtherElementIsALogicalOr(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	or := Operator(tagfilter.LogicalOr)
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	nestedOr := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression, primaryExpression})
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression, nestedOr})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{Left: &LogicalAndExpression{Left: &BracketExpression{Bracket: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: operator,
						},
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{Left: &LogicalAndExpression{Left: &BracketExpression{Bracket: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator: operator,
							},
						},
					},
				},
				Operator: &or,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &BracketExpression{
							Primary: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
									Operator: operator,
								},
							},
						},
					},
				},
			}}}},
		}}}},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalOrWithTwoElementsWhereFirstElementIsAPrimaryExpressionAndTheOtherElementIsALogicalAnd(t *testing.T) {
	operator := Operator(tagfilter.IsEmptyOperator)
	or := Operator(tagfilter.LogicalOr)
	and := Operator(tagfilter.LogicalAnd)
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	nestedAnd := tagfilter.NewLogicalAndTagFilter([]*tagfilter.TagFilter{primaryExpression, primaryExpression})
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression, nestedAnd})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{Left: &LogicalAndExpression{Left: &BracketExpression{Bracket: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: operator,
						},
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{Left: &LogicalAndExpression{Left: &BracketExpression{Bracket: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator: operator,
							},
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &BracketExpression{
							Primary: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
									Operator: operator,
								},
							},
						},
					},
				},
			}}}},
		}}}},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldFailToMapLogicalOrWhenFirstElementIsALogicalOrExpression(t *testing.T) {
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	nestedOr := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression, primaryExpression})
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{nestedOr, primaryExpression})

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for first element")
}

func TestShouldFailToMapLogicalOrWhenOnlyOneElementIsProvided(t *testing.T) {
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	input := tagfilter.NewLogicalOrTagFilter([]*tagfilter.TagFilter{primaryExpression})

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "at least two elements are expected for logical or")
}

func TestShouldFailToMapTagFilterExpressionFromInstanaAPIWhenLogicalOperatorIsNotValid(t *testing.T) {
	primaryExpression := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
	operator := tagfilter.LogicalOperatorType("FOO")
	input := &tagfilter.TagFilter{
		Type:            tagfilter.TagFilterExpressionType,
		LogicalOperator: &operator,
		Elements:        []*tagfilter.TagFilter{primaryExpression, primaryExpression},
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid logical operator")
}

func TestShouldReturnMappingErrorWhenAnyElementOfTagFilterExpressionIsNotValid(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("TestShouldReturnMappingErrorWhenElement%dOfTagFilterExpressionIsNotValid", i), func(t *testing.T) {
			invalidElement := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, "INVALID")
			validElement := tagfilter.NewUnaryTagFilter(tagfilter.TagFilterEntityDestination, tagFilterName, tagfilter.IsEmptyOperator)
			elements := make([]*tagfilter.TagFilter, 5)
			for j := 0; j < 5; j++ {
				if j == i {
					elements[j] = invalidElement
				} else {
					elements[j] = validElement
				}
			}
			input := tagfilter.NewLogicalOrTagFilter(elements)

			mapper := NewMapper()
			_, err := mapper.FromAPIModel(input)

			require.NotNil(t, err)
			require.Contains(t, err.Error(), invalidOperator)
			require.Contains(t, err.Error(), tagFilterOperator)
		})
	}
}

func runTestCaseForMappingFromAPI(input *tagfilter.TagFilter, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewMapper()
	result, err := mapper.FromAPIModel(input)

	require.Nil(t, err)
	require.Equal(t, expectedResult, result)
}
