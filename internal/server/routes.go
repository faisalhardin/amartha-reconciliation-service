package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/bank-statement", func(r chi.Router) {
			r.Post("/upload", s.handlers.BankStatementHandler.Upload)
		})

		r.Get("/reconciliation/summary/{fileID}", s.handlers.ReconciliationHandler.Summary)

		r.Route("/reconciliation/{bankCode}", func(r chi.Router) {
			h := s.handlers.ReconciliationHandler
			r.Post("/run", h.Run)
			r.Get("/", h.List)
		})

		r.Route("/transaction/{bankCode}", func(r chi.Router) {
			h := s.handlers.TransactionHandler
			r.Post("/upload", h.UploadCSV)
			r.Post("/", h.Create)
			r.Get("/", h.List)
			r.Get("/{id}", h.GetByID)
		})
	})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
