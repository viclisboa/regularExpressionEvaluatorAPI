package service

import (
	"errors"
	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/util"
	"regexp"
	"strconv"
	"strings"
)

type ExpressionService struct {
	Logger log.Entry
}

func (es *ExpressionService) ExecuteExpression(expression model.Expression, urlParams string) (model.Response, error) {
	logger := es.Logger.WithField("expressionId", expression.ID)
	searchRegexOr := regexp.MustCompile("(?i)" + "or")
	searchRegexAnd := regexp.MustCompile("(?i)" + "and")

	expressionString := expression.Definition
	expressionString = searchRegexOr.ReplaceAllString(expressionString, "||")
	expressionString = searchRegexAnd.ReplaceAllString(expressionString, "&&")

	logger.WithFields(log.Fields{
		"expression": expression.Definition,
		"params":     urlParams,
	})

	evaluateExpression, err := govaluate.NewEvaluableExpression(expressionString)
	if err != nil {
		logger.WithField("err", err.Error())
		logger.Error("error creating evaluable expression")
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
		logger.WithField("err", err.Error())
		logger.Error("error evaluating expression")
		return model.Response{}, errors.New(util.ErrEvaluatingExpression)
	}

	response := model.Response{
		Definition: expression.Definition,
		Values:     urlParams,
		Result:     result.(bool),
	}
	logger.WithField("result", result).Info("expression evaluated successfully")

	return response, nil
}
