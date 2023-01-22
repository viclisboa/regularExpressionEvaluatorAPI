package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/repository"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/service"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ExpressionHandler struct {
	ExpressionService    service.ExpressionService
	ExpressionRepository repository.ExpressionInterface
	Logger               log.Logger
}

func (eh *ExpressionHandler) EvaluateExpression(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := GetURLParams(r)
	expressionId := params["expressionId"]
	logger := eh.Logger.WithField("expressionId", expressionId)

	expressionIdAsInt, err := strconv.Atoi(expressionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithField("err", err.Error()).Error("error parsing expressionId to int")
		return
	}

	expression, err := eh.ExpressionRepository.GetExpressionById(expressionIdAsInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger.WithField("err", err.Error()).Error("error recovering expression from database")
		return
	}

	result, err := eh.ExpressionService.ExecuteExpression(expression, r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithField("err", err.Error()).Error("error resolving expression")
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error on marshal ", http.StatusInternalServerError)
		logger.WithField("err", err.Error()).Error("error encoding response")
		return
	}
}

func (eh *ExpressionHandler) SaveExpression(w http.ResponseWriter, r *http.Request) {
	params := GetURLParams(r)
	expressionId := params["expressionId"]

	logger := eh.Logger.WithField("expressionId", expressionId)

	expressionIdAsInt, err := strconv.Atoi(expressionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithField("err", err.Error()).Error("error parsing expressionId to int")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error on get body data for expression update")
		http.Error(w, "Error on get body data for expression update", http.StatusInternalServerError)
		return
	}

	var body map[string]any
	err = json.Unmarshal(b, &body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error on unmarshal  payload for update")
		http.Error(w, "Error on unmarshal landing page payload for creation", http.StatusInternalServerError)
		return
	}

	definition, exists := body["definition"]
	if definition == "" || !exists {
		logger.WithField("body", body).Error("missing definition on body")
		http.Error(w, "missing definition in body", http.StatusBadRequest)
		return
	}

	definitionAsString := definition.(string)

	err = eh.ExpressionRepository.SaveExpression(expressionIdAsInt, definitionAsString)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error updating expressiong")
		http.Error(w, "Error updating expression", http.StatusInternalServerError)
		return
	}

	logger.Info("expression updated successfully")
}

func (eh *ExpressionHandler) CreateExpression(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		eh.Logger.WithField("err", err.Error()).Error("Error on get body data for expression update")
		http.Error(w, "Error on get body data for expression update", http.StatusInternalServerError)
		return
	}

	var body map[string]any
	err = json.Unmarshal(b, &body)
	if err != nil {
		eh.Logger.WithField("err", err.Error()).Error("Error on unmarshal  payload for update")
		http.Error(w, "Error on unmarshal landing page payload for creation", http.StatusInternalServerError)
		return
	}

	definition, exists := body["definition"].(string)
	if definition == "" || !exists {
		eh.Logger.WithField("body", body).Error("missing definition on body")
		http.Error(w, "missing definition in body", http.StatusBadRequest)
		return
	}

	err = eh.ExpressionRepository.CreateExpression(definition)
	if err != nil {
		eh.Logger.WithField("err", err.Error()).Error("Error updating expressiong")
		http.Error(w, "Error updating expression", http.StatusInternalServerError)
		return
	}

	eh.Logger.Info("expression created successfully")
}

func (eh *ExpressionHandler) GetAllExpressions(w http.ResponseWriter, r *http.Request) {
	expressions, err := eh.ExpressionRepository.GetAllExpressions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		eh.Logger.WithField("err", err.Error()).Error("error recovering expression from database")
		return
	}

	if err := json.NewEncoder(w).Encode(expressions); err != nil {
		http.Error(w, "Error on marshal ", http.StatusInternalServerError)
		eh.Logger.WithField("err", err.Error()).Error("error encoding response")
		return
	}

	eh.Logger.Info("all expressions recovered successfully")
}

func (eh *ExpressionHandler) DeleteExpression(w http.ResponseWriter, r *http.Request) {
	params := GetURLParams(r)
	expressionId := params["expressionId"]

	logger := eh.Logger.WithField("expressionId", expressionId)

	expressionIdAsInt, err := strconv.Atoi(expressionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.WithField("err", err.Error()).Error("error parsing expressionId to int")
		return
	}

	err = eh.ExpressionRepository.DeleteExpression(expressionIdAsInt)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting expression")
		http.Error(w, "Error updating expression", http.StatusInternalServerError)
		return
	}

	logger.Info("expression deleted successfully")
}

func GetURLParams(r *http.Request) map[string]string {
	rctx := chi.RouteContext(r.Context())
	var urlParams map[string]string
	for k := len(rctx.URLParams.Keys) - 1; k >= 0; k-- {
		switch urlParams {
		case nil:
			urlParams = make(map[string]string)
			fallthrough
		default:
			if rctx.URLParams.Keys[k] != "*" {
				urlParams[rctx.URLParams.Keys[k]] = rctx.URLParams.Values[k]
			}
		}
	}
	return urlParams
}
