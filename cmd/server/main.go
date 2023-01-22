package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/handler"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/repository"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/service"
	"net/http"
	"os"
	"time"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	repo, err := repository.NewRepository("postgres", `host=localhost user=pg password=pass dbname=postgres port=5432 sslmode=disable`)
	if err != nil {
		panic("error initializing database")
	}

	svc := service.ExpressionService{}

	expressionHandler := handler.ExpressionHandler{
		ExpressionService:    svc,
		ExpressionRepository: &repo,
	}

	credentials := map[string]string{
		"testeUser": "testePassword",
	}

	r := chi.NewRouter()
	r.Use(middleware.BasicAuth("", credentials))
	r.Get("/evaluate/{expressionId}", expressionHandler.EvaluateExpression)
	r.Get("/expressions", expressionHandler.GetAllExpressions)
	r.Post("/expressions/{expressionId}", expressionHandler.SaveExpression)
	r.Delete("/expressions/{expressionId}", expressionHandler.SaveExpression)
	r.Post("/expressions", expressionHandler.CreateExpression)

	http.Handle("/", r)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", "8080"),
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("starting server")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("listen and serve died", "err", err)
	}
}
