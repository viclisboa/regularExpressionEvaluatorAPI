package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/repository"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveExpression(t *testing.T) {
	testCases := []struct {
		name         string
		databaseMock repository.Stub
		requestBody  map[string]any
		httpStatus   int
	}{
		{
			name: "should return 200",
			databaseMock: repository.Stub{
				SaveExpressionError: nil,
			},
			requestBody: map[string]any{
				"definition": "a or b",
			},
			httpStatus: http.StatusOK,
		},
		{
			name:         "should return 400, missing definition in json",
			databaseMock: repository.Stub{},
			requestBody: map[string]any{
				"sdsadasd": "dasdsa",
			},
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "should return 500, error saving in database",
			databaseMock: repository.Stub{
				SaveExpressionError: errors.New("not found"),
			},
			requestBody: map[string]any{
				"definition": "dasdsa",
			},
			httpStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := ExpressionHandler{
				ExpressionService:    service.ExpressionService{},
				ExpressionRepository: &tc.databaseMock,
				Logger:               log.Entry{},
			}

			r := chi.NewRouter()
			r.HandleFunc("/expressions/{expressionId}", handler.SaveExpression)
			ts := httptest.NewServer(r)
			defer ts.Close()

			url := ts.URL + "/expressions/1"

			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(tc.requestBody)

			response, _ := http.Post(url, "application/json", &buf)

			assert.Equal(t, tc.httpStatus, response.StatusCode)
		})
	}
}

func TestCreateExpression(t *testing.T) {
	testCases := []struct {
		name         string
		databaseMock repository.Stub
		requestBody  map[string]any
		httpStatus   int
	}{
		{
			name: "should return 200",
			databaseMock: repository.Stub{
				CreateExpressionError: nil,
			},
			requestBody: map[string]any{
				"definition": "a or b",
			},
			httpStatus: http.StatusOK,
		},
		{
			name:         "should return 400, missing definition in json",
			databaseMock: repository.Stub{},
			requestBody: map[string]any{
				"sdsadasd": "dasdsa",
			},
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "should return 500, error saving in database",
			databaseMock: repository.Stub{
				CreateExpressionError: errors.New("not found"),
			},
			requestBody: map[string]any{
				"definition": "dasdsa",
			},
			httpStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := ExpressionHandler{
				ExpressionService:    service.ExpressionService{},
				ExpressionRepository: &tc.databaseMock,
				Logger:               log.Entry{},
			}

			r := chi.NewRouter()
			r.HandleFunc("/expressions", handler.CreateExpression)
			ts := httptest.NewServer(r)
			defer ts.Close()

			url := ts.URL + "/expressions"

			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(tc.requestBody)

			response, _ := http.Post(url, "application/json", &buf)

			assert.Equal(t, tc.httpStatus, response.StatusCode)
		})
	}
}

func TestEvaluateExpression(t *testing.T) {
	testCases := []struct {
		name         string
		databaseMock repository.Stub
		httpStatus   int
		expectedBody model.Response
		queryString  string
	}{
		{
			name: "should return 200",
			databaseMock: repository.Stub{
				GetExpressionByIdResponse: model.Expression{
					ID:         10,
					Definition: "x or y",
				},
				GetExpressionByIdError: nil,
			},
			httpStatus: http.StatusOK,
			expectedBody: model.Response{
				Definition: "x or y",
				Values:     "x=1,y=0",
				Result:     true,
			},
			queryString: "?x=1,y=0",
		},
		{
			name: "should return 404, expression not found",
			databaseMock: repository.Stub{
				GetExpressionByIdError: errors.New("not found"),
			},
			httpStatus:   http.StatusNotFound,
			expectedBody: model.Response{},
			queryString:  "?x=1,y=0,z=1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := ExpressionHandler{
				ExpressionService:    service.ExpressionService{},
				ExpressionRepository: &tc.databaseMock,
				Logger:               log.Entry{},
			}

			r := chi.NewRouter()
			r.HandleFunc("/expressions/{expressionId}", handler.EvaluateExpression)
			ts := httptest.NewServer(r)
			defer ts.Close()

			url := ts.URL + "/expressions/10" + tc.queryString

			response, _ := http.Get(url)

			var parsedResponse model.Response

			_ = json.NewDecoder(response.Body).Decode(&parsedResponse)

			assert.Equal(t, tc.httpStatus, response.StatusCode)
			assert.Equal(t, tc.expectedBody, parsedResponse)
		})
	}
}
