package repository

import (
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
)

var _ ExpressionInterface = (*Stub)(nil)

type Stub struct {
	CreateExpressionError      error
	CreateExpressionCalledWith map[string]any

	SaveExpressionError      error
	SaveExpressionCalledWith map[string]any

	GetAllExpressionsResponse []model.Expression
	GetAllExpressionsError    error

	GetExpressionByIdResponse   model.Expression
	GetExpressionByIdError      error
	GetExpressionByIdCalledWith map[string]any

	DeleteExpressionCalledWith map[string]any
	DeleteExpressionError      error
}

func (s *Stub) GetAllExpressions() ([]model.Expression, error) {
	return s.GetAllExpressionsResponse, s.GetAllExpressionsError
}

func (s *Stub) GetExpressionById(expressionId int) (model.Expression, error) {
	s.GetExpressionByIdCalledWith = map[string]any{
		"expressionId": expressionId,
	}

	return s.GetExpressionByIdResponse, s.GetExpressionByIdError
}

func (s *Stub) CreateExpression(definition string) error {
	s.CreateExpressionCalledWith = map[string]any{
		"definition": definition,
	}
	return s.CreateExpressionError
}

func (s *Stub) SaveExpression(expressionId int, definition string) error {
	s.SaveExpressionCalledWith = map[string]any{
		"expressionId": expressionId,
		"definition":   definition,
	}
	return s.SaveExpressionError
}

func (s *Stub) DeleteExpression(expressionId int) error {
	s.DeleteExpressionCalledWith = map[string]any{
		"expressionId": expressionId,
	}
	return s.DeleteExpressionError
}
