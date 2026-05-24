package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	entityhttp "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/http"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/database"
)

type Server struct {
	port     int
	db       database.Service
	handlers *entityhttp.Handlers
}

func NewServer(handlers *entityhttp.Handlers) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	s := &Server{
		port:     port,
		db:       database.New(),
		handlers: handlers,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
