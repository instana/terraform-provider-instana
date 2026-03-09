package tagfilter_test

import (
	"fmt"
	"testing"

	"github.com/instana/terraform-provider-instana/utils"

	"github.com/stretchr/testify/require"

	"github.com/instana/instana-go-client/instana"
	. "github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
)

const (
	invalidOperator   = "invalid operator"
	tagFilterOperator = "tag filter operator"
	tagFilterName     = "name"
)

func TestShouldMapEmptyTagFilterExpressionFromInstanaAPI(t *testing.T) {
	for _, operator := range instana.SupportedLogicalOperatorTypes {
		t.Run(fmt.Sprintf("TestShouldMapEmpty%sTagFilterExpressionFromInstnaAPI", string(operator)), func(t *testing.T) {
			op := operator
			expression := &instana.TagFilter{
				Type:            instana.TagFilterExpressionType,
				LogicalOperator: &op,
			}

			runTestCaseForMappingFromAPI(expression, nil, t)
		})
	}
}

func TestShouldMapStringTagFilterFromInstanaAPI(t *testing.T) {
	value := "value"
	input := instana.NewStringTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(instana.EqualsOperator),
		StringValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapNumberTagFilterFromInstanaAPI(t *testing.T) {
	value := int64(1234)
	input := instana.NewNumberTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(instana.EqualsOperator),
		NumberValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapBooleanTagFilterFromInstanaAPI(t *testing.T) {
	value := true
	input := instana.NewBooleanTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:       &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:     Operator(instana.EqualsOperator),
		BooleanValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapComparisonTagFilterWithTagKeyValueFromInstanaAPI(t *testing.T) {
	key := "key"
	value := "value"
	input := instana.NewTagTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.EqualsOperator, key, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(instana.EqualsOperator),
		StringValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func testMappingOfTagFilterFromInstanaApi(tagFilter *instana.TagFilter, comparison *ComparisonExpression, t *testing.T) {
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
	for _, v := range instana.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), testMappingOfSupportedComparisonOperatorsFromInstanaAPI(v))
	}
}

func testMappingOfSupportedComparisonOperatorsFromInstanaAPI(operator instana.ExpressionOperator) func(t *testing.T) {
	return func(t *testing.T) {
		value := "value"
		input := instana.NewStringTagFilter(instana.TagFilterEntityDestination, tagFilterName, operator, value)

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
	input := instana.NewStringTagFilter(instana.TagFilterEntityDestination, tagFilterName, "FOO", value)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldMapAllSupportedUnaryOperationsFromInstanaAPI(t *testing.T) {
	for _, v := range instana.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfSupportedUnaryOperationFromInstanaAPI(v))
	}
}

func testMappingOfSupportedUnaryOperationFromInstanaAPI(operator instana.ExpressionOperator) func(t *testing.T) {
	return func(t *testing.T) {
		input := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, operator)

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
	input := instana.NewUnaryTagFilterWithTagKey(instana.TagFilterEntityDestination, tagFilterName, &key, instana.NotEmptyOperator)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(instana.NotEmptyOperator),
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldFailToMapTagFilterFromInstanaAPIWhenUnaryOperationIsNotSupported(t *testing.T) {
	input := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, "FOO")

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldFailToMapTagFilterExpressionElementFromInstanaAPIWhenTypeIsMissing(t *testing.T) {
	name := tagFilterName
	operator := instana.ExpressionOperator("FOO")
	input := &instana.TagFilter{
		Name:     &name,
		Operator: &operator,
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported tag filter expression")
}

func TestShouldMapLogicalAndWithTwoPrimaryExpressionsFromInstanaAPI(t *testing.T) {
	operator := Operator(instana.IsEmptyOperator)
	and := Operator(instana.LogicalAnd)
	primaryExpression1 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name1", instana.IsEmptyOperator)
	primaryExpression2 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name2", instana.IsEmptyOperator)
	input := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression1, primaryExpression2})

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
	operator := Operator(instana.IsEmptyOperator)
	and := Operator(instana.LogicalAnd)
	primaryExpression1 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name1", instana.IsEmptyOperator)
	primaryExpression2 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name2", instana.IsEmptyOperator)
	primaryExpression3 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name3", instana.IsEmptyOperator)
	input := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression1, primaryExpression2, primaryExpression3})

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
	operator := Operator(instana.IsEmptyOperator)
	and := Operator(instana.LogicalAnd)
	primaryExpression1 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	primaryExpression2 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name2", instana.IsEmptyOperator)
	nestedAnd := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression2, primaryExpression2})
	input := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression1, nestedAnd})

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
	operator := Operator(instana.IsEmptyOperator)
	and := Operator(instana.LogicalAnd)
	or := Operator(instana.LogicalOr)
	primaryExpression1 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	primaryExpression2 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name2", instana.IsEmptyOperator)
	nestedOr := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression2, primaryExpression2})
	input := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression1, nestedOr})

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
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	input := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(instana.IsEmptyOperator),
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
	operator := Operator(instana.IsEmptyOperator)
	or := Operator(instana.LogicalOr)
	primaryExpression1 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name1", instana.IsEmptyOperator)
	primaryExpression2 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name2", instana.IsEmptyOperator)
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression1, primaryExpression2})

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
	operator := Operator(instana.IsEmptyOperator)
	or := Operator(instana.LogicalOr)
	primaryExpression1 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name1", instana.IsEmptyOperator)
	primaryExpression2 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name2", instana.IsEmptyOperator)
	primaryExpression3 := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, "name3", instana.IsEmptyOperator)
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression1, primaryExpression2, primaryExpression3})

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
	operator := Operator(instana.IsEmptyOperator)
	or := Operator(instana.LogicalOr)
	and := Operator(instana.LogicalAnd)
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	nestedAnd := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression, primaryExpression})
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{nestedAnd, primaryExpression})

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
	operator := Operator(instana.IsEmptyOperator)
	or := Operator(instana.LogicalOr)
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	nestedOr := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression, primaryExpression})
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression, nestedOr})

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
	operator := Operator(instana.IsEmptyOperator)
	or := Operator(instana.LogicalOr)
	and := Operator(instana.LogicalAnd)
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	nestedAnd := instana.NewLogicalAndTagFilter([]*instana.TagFilter{primaryExpression, primaryExpression})
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression, nestedAnd})

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
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	nestedOr := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression, primaryExpression})
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{nestedOr, primaryExpression})

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for first element")
}

func TestShouldFailToMapLogicalOrWhenOnlyOneElementIsProvided(t *testing.T) {
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	input := instana.NewLogicalOrTagFilter([]*instana.TagFilter{primaryExpression})

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "at least two elements are expected for logical or")
}

func TestShouldFailToMapTagFilterExpressionFromInstanaAPIWhenLogicalOperatorIsNotValid(t *testing.T) {
	primaryExpression := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
	operator := instana.LogicalOperatorType("FOO")
	input := &instana.TagFilter{
		Type:            instana.TagFilterExpressionType,
		LogicalOperator: &operator,
		Elements:        []*instana.TagFilter{primaryExpression, primaryExpression},
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid logical operator")
}

func TestShouldReturnMappingErrorWhenAnyElementOfTagFilterExpressionIsNotValid(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("TestShouldReturnMappingErrorWhenElement%dOfTagFilterExpressionIsNotValid", i), func(t *testing.T) {
			invalidElement := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, "INVALID")
			validElement := instana.NewUnaryTagFilter(instana.TagFilterEntityDestination, tagFilterName, instana.IsEmptyOperator)
			elements := make([]*instana.TagFilter, 5)
			for j := 0; j < 5; j++ {
				if j == i {
					elements[j] = invalidElement
				} else {
					elements[j] = validElement
				}
			}
			input := instana.NewLogicalOrTagFilter(elements)

			mapper := NewMapper()
			_, err := mapper.FromAPIModel(input)

			require.NotNil(t, err)
			require.Contains(t, err.Error(), invalidOperator)
			require.Contains(t, err.Error(), tagFilterOperator)
		})
	}
}

func runTestCaseForMappingFromAPI(input *instana.TagFilter, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewMapper()
	result, err := mapper.FromAPIModel(input)

	require.Nil(t, err)
	require.Equal(t, expectedResult, result)
}
