package service

import (
	"errors"
	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/repository"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/util"
	"regexp"
	"strconv"
	"strings"
)

type ExpressionService struct {
	ExpressionRepository repository.ExpressionInterface
	Logger               log.Logger
}

func NewExpressionService(repo repository.Repository, balance int) ExpressionService {
	return ExpressionService{
		ExpressionRepository: &repo,
	}
}

func (es *ExpressionService) ExecuteExpression(expressionId string, urlParams string) (model.Response, error) {
	logger := es.Logger.WithField("expressionId", expressionId)

	expressionIdAsInt, err := strconv.Atoi(expressionId)
	if err != nil {
		logger.Error("error parsing expressionId to int")
		return model.Response{}, errors.New(util.ErrParsingExpressionId)
	}

	expression, err := es.ExpressionRepository.GetExpressionById(expressionIdAsInt)
	if err != nil {
		logger.WithField("err", err.Error())
		logger.Error("error recovering expression from database")
		return model.Response{}, errors.New(util.ErrRecoveringExpressionFromDatabase)
	}

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
