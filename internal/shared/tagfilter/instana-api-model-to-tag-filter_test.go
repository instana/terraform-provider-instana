package tagfilter_test

import (
	"fmt"
	"testing"

	"github.com/instana/terraform-provider-instana/utils"

	"github.com/stretchr/testify/require"

	tag "github.com/instana/instana-go-client/shared/tagfilter"
	common "github.com/instana/instana-go-client/shared/types"
	. "github.com/instana/terraform-provider-instana/internal/shared/tagfilter"
)

const (
	invalidOperator   = "invalid operator"
	tagFilterOperator = "tag filter operator"
	tagFilterName     = "name"
)

func TestShouldMapEmptyTagFilterExpressionFromInstanaAPI(t *testing.T) {
	for _, operator := range common.SupportedLogicalOperatorTypes {
		t.Run(fmt.Sprintf("TestShouldMapEmpty%sTagFilterExpressionFromInstnaAPI", string(operator)), func(t *testing.T) {
			op := operator
			expression := &tag.TagFilter{
				Type:            tag.TagFilterExpressionType,
				LogicalOperator: &op,
			}

			runTestCaseForMappingFromAPI(expression, nil, t)
		})
	}
}

func TestShouldMapStringTagFilterFromInstanaAPI(t *testing.T) {
	value := "value"
	input := tag.NewStringTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(common.EqualsOperator),
		StringValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapNumberTagFilterFromInstanaAPI(t *testing.T) {
	value := int64(1234)
	input := tag.NewNumberTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(common.EqualsOperator),
		NumberValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapBooleanTagFilterFromInstanaAPI(t *testing.T) {
	value := true
	input := tag.NewBooleanTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.EqualsOperator, value)

	comparison := &ComparisonExpression{
		Entity:       &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:     Operator(common.EqualsOperator),
		BooleanValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func TestShouldMapComparisonTagFilterWithTagKeyValueFromInstanaAPI(t *testing.T) {
	key := "key"
	value := "value"
	input := tag.NewTagTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.EqualsOperator, key, value)

	comparison := &ComparisonExpression{
		Entity:      &EntitySpec{Identifier: tagFilterName, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
		Operator:    Operator(common.EqualsOperator),
		StringValue: &value,
	}

	testMappingOfTagFilterFromInstanaApi(input, comparison, t)
}

func testMappingOfTagFilterFromInstanaApi(tagFilter *tag.TagFilter, comparison *ComparisonExpression, t *testing.T) {
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
	for _, v := range common.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), testMappingOfSupportedComparisonOperatorsFromInstanaAPI(v))
	}
}

func testMappingOfSupportedComparisonOperatorsFromInstanaAPI(operator common.ExpressionOperator) func(t *testing.T) {
	return func(t *testing.T) {
		value := "value"
		input := tag.NewStringTagFilter(tag.TagFilterEntityDestination, tagFilterName, operator, value)

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
	input := tag.NewStringTagFilter(tag.TagFilterEntityDestination, tagFilterName, "FOO", value)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldMapAllSupportedUnaryOperationsFromInstanaAPI(t *testing.T) {
	for _, v := range common.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfSupportedUnaryOperationFromInstanaAPI(v))
	}
}

func testMappingOfSupportedUnaryOperationFromInstanaAPI(operator common.ExpressionOperator) func(t *testing.T) {
	return func(t *testing.T) {
		input := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, operator)

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
	input := tag.NewUnaryTagFilterWithTagKey(tag.TagFilterEntityDestination, tagFilterName, &key, common.NotEmptyOperator)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, TagKey: &key, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(common.NotEmptyOperator),
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldFailToMapTagFilterFromInstanaAPIWhenUnaryOperationIsNotSupported(t *testing.T) {
	input := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, "FOO")

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldFailToMapTagFilterExpressionElementFromInstanaAPIWhenTypeIsMissing(t *testing.T) {
	name := tagFilterName
	operator := common.ExpressionOperator("FOO")
	input := &tag.TagFilter{
		Name:     &name,
		Operator: &operator,
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported tag filter expression")
}

func TestShouldMapLogicalAndWithTwoPrimaryExpressionsFromInstanaAPI(t *testing.T) {
	operator := Operator(common.IsEmptyOperator)
	and := Operator(common.LogicalAnd)
	primaryExpression1 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name1", common.IsEmptyOperator)
	primaryExpression2 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name2", common.IsEmptyOperator)
	input := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression1, primaryExpression2})

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
	operator := Operator(common.IsEmptyOperator)
	and := Operator(common.LogicalAnd)
	primaryExpression1 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name1", common.IsEmptyOperator)
	primaryExpression2 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name2", common.IsEmptyOperator)
	primaryExpression3 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name3", common.IsEmptyOperator)
	input := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression1, primaryExpression2, primaryExpression3})

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
	operator := Operator(common.IsEmptyOperator)
	and := Operator(common.LogicalAnd)
	primaryExpression1 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	primaryExpression2 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name2", common.IsEmptyOperator)
	nestedAnd := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression2, primaryExpression2})
	input := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression1, nestedAnd})

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
	operator := Operator(common.IsEmptyOperator)
	and := Operator(common.LogicalAnd)
	or := Operator(common.LogicalOr)
	primaryExpression1 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	primaryExpression2 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name2", common.IsEmptyOperator)
	nestedOr := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression2, primaryExpression2})
	input := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression1, nestedOr})

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
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	input := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: tagFilterName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(common.IsEmptyOperator),
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
	operator := Operator(common.IsEmptyOperator)
	or := Operator(common.LogicalOr)
	primaryExpression1 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name1", common.IsEmptyOperator)
	primaryExpression2 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name2", common.IsEmptyOperator)
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression1, primaryExpression2})

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
	operator := Operator(common.IsEmptyOperator)
	or := Operator(common.LogicalOr)
	primaryExpression1 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name1", common.IsEmptyOperator)
	primaryExpression2 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name2", common.IsEmptyOperator)
	primaryExpression3 := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, "name3", common.IsEmptyOperator)
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression1, primaryExpression2, primaryExpression3})

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
	operator := Operator(common.IsEmptyOperator)
	or := Operator(common.LogicalOr)
	and := Operator(common.LogicalAnd)
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	nestedAnd := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression, primaryExpression})
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{nestedAnd, primaryExpression})

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
	operator := Operator(common.IsEmptyOperator)
	or := Operator(common.LogicalOr)
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	nestedOr := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression, primaryExpression})
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression, nestedOr})

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
	operator := Operator(common.IsEmptyOperator)
	or := Operator(common.LogicalOr)
	and := Operator(common.LogicalAnd)
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	nestedAnd := tag.NewLogicalAndTagFilter([]*tag.TagFilter{primaryExpression, primaryExpression})
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression, nestedAnd})

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
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	nestedOr := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression, primaryExpression})
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{nestedOr, primaryExpression})

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for first element")
}

func TestShouldFailToMapLogicalOrWhenOnlyOneElementIsProvided(t *testing.T) {
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	input := tag.NewLogicalOrTagFilter([]*tag.TagFilter{primaryExpression})

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "at least two elements are expected for logical or")
}

func TestShouldFailToMapTagFilterExpressionFromInstanaAPIWhenLogicalOperatorIsNotValid(t *testing.T) {
	primaryExpression := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
	operator := common.LogicalOperatorType("FOO")
	input := &tag.TagFilter{
		Type:            tag.TagFilterExpressionType,
		LogicalOperator: &operator,
		Elements:        []*tag.TagFilter{primaryExpression, primaryExpression},
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid logical operator")
}

func TestShouldReturnMappingErrorWhenAnyElementOfTagFilterExpressionIsNotValid(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("TestShouldReturnMappingErrorWhenElement%dOfTagFilterExpressionIsNotValid", i), func(t *testing.T) {
			invalidElement := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, "INVALID")
			validElement := tag.NewUnaryTagFilter(tag.TagFilterEntityDestination, tagFilterName, common.IsEmptyOperator)
			elements := make([]*tag.TagFilter, 5)
			for j := 0; j < 5; j++ {
				if j == i {
					elements[j] = invalidElement
				} else {
					elements[j] = validElement
				}
			}
			input := tag.NewLogicalOrTagFilter(elements)

			mapper := NewMapper()
			_, err := mapper.FromAPIModel(input)

			require.NotNil(t, err)
			require.Contains(t, err.Error(), invalidOperator)
			require.Contains(t, err.Error(), tagFilterOperator)
		})
	}
}

func runTestCaseForMappingFromAPI(input *tag.TagFilter, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewMapper()
	result, err := mapper.FromAPIModel(input)

	require.Nil(t, err)
	require.Equal(t, expectedResult, result)
}
