package service

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/repository"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/util"
	"regexp"
	"strconv"
	"strings"
)

type ExpressionService struct {
	expressionRepository repository.ExpressionInterface
}

func NewExpressionService(repo repository.Repository, balance int) ExpressionService {
	return ExpressionService{
		expressionRepository: &repo,
	}
}

func (es *ExpressionService) ExecuteExpression(expressionId string, urlParams string) (model.Response, error) {

	expressionIdAsInt, err := strconv.Atoi(expressionId)
	if err != nil {
		return model.Response{}, errors.New(util.ErrParsingExpressionId)
	}

	expression, err := es.expressionRepository.GetExpressionById(expressionIdAsInt)
	if err != nil {
		return model.Response{}, errors.New(util.ErrRecoveringExpressionFromDatabase)
	}

	searchRegexOr := regexp.MustCompile("(?i)" + "or")
	searchRegexAnd := regexp.MustCompile("(?i)" + "and")

	expressionString := expression.Definition
	expressionString = searchRegexOr.ReplaceAllString(expressionString, "||")
	expressionString = searchRegexAnd.ReplaceAllString(expressionString, "&&")

	evaluateExpression, err := govaluate.NewEvaluableExpression(expressionString)
	if err != nil {
		return model.Response{}, errors.New(util.ErrCreatingEvaluableExpression)
	}

	params := strings.Split(urlParams, ",")
	parameters := make(map[string]interface{}, len(params))

	for _, param := range params {
		variable := strings.Split(param, "=")
		if len(variable) > 1 {
			boolValue, err := strconv.ParseBool(variable[1])
			if err == nil {
				parameters[variable[0]] = boolValue
			}
		}
	}

	result, err := evaluateExpression.Evaluate(parameters)
	if err != nil {
		return model.Response{}, errors.New(util.ErrEvaluatingExpression)
	}

	response := model.Response{
		Definition: expression.Definition,
		Values:     urlParams,
		Result:     result.(bool),
	}

	fmt.Println(response)

	return response, nil
}
