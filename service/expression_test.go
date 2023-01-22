package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/repository"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/util"
	"testing"
)

func TestExpression_ExecuteExpression(t *testing.T) {
	testCases := []struct {
		name           string
		databaseMock   repository.Stub
		expectedResult model.Response
		expectedError  error
		urlParams      string
		expressionId   string
	}{
		{
			name:           "should return error, expressionId not a number",
			databaseMock:   repository.Stub{},
			expectedResult: model.Response{},
			expectedError:  errors.New(util.ErrParsingExpressionId),
			urlParams:      "",
			expressionId:   "potato",
		},
		{
			name: "should return error, error recovering from database",
			databaseMock: repository.Stub{
				GetExpressionByIdError: errors.New("error recovering from database"),
			},
			expectedResult: model.Response{},
			expectedError:  errors.New(util.ErrRecoveringExpressionFromDatabase),
			urlParams:      "",
			expressionId:   "5",
		},
		{
			name: "should return error, error creating expression ",
			databaseMock: repository.Stub{
				GetExpressionByIdResponse: model.Expression{
					ID:         1,
					Definition: "&",
				},
			},
			expectedResult: model.Response{},
			expectedError:  errors.New(util.ErrCreatingEvaluableExpression),
			urlParams:      "",
			expressionId:   "5",
		},
		{
			name: "should return error, error evaluating expression",
			databaseMock: repository.Stub{
				GetExpressionByIdResponse: model.Expression{
					ID:         1,
					Definition: "x OR y",
				},
			},
			expectedResult: model.Response{},
			expectedError:  errors.New(util.ErrEvaluatingExpression),
			urlParams:      "",
			expressionId:   "5",
		},
		{
			name: "should return success, operators with capital letters",
			databaseMock: repository.Stub{
				GetExpressionByIdResponse: model.Expression{
					ID:         1,
					Definition: "x OR y",
				},
			},
			expectedResult: model.Response{
				Definition: "x OR y",
				Values:     "x=1,y=0",
				Result:     true,
			},
			expectedError: nil,
			urlParams:     "x=1,y=0",
			expressionId:  "5",
		},
		{
			name: "should return success, operators without capital letters",
			databaseMock: repository.Stub{
				GetExpressionByIdResponse: model.Expression{
					ID:         1,
					Definition: "x or y",
				},
			},
			expectedResult: model.Response{
				Definition: "x or y",
				Values:     "x=1,y=0",
				Result:     true,
			},
			expectedError: nil,
			urlParams:     "x=1,y=0",
			expressionId:  "5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := ExpressionService{expressionRepository: &tc.databaseMock}

			result, err := service.ExecuteExpression(tc.expressionId, tc.urlParams)
			assert.Equal(t, tc.expectedResult, result, "values should be the same")

			if err != nil {
				assert.Equal(t, tc.expectedError.Error(), err.Error(), "values should be the same")
			}
		})
	}
}
