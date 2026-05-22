package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Server Server
	DB     DBConfig
}

type Server struct {
	Name    string
	Host    string
	Port    string
	BaseURL string
}

type DBConfig struct {
	DSN string
}

var (
	dbHost     = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbUser     = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
	dbSchema   = os.Getenv("DB_SCHEMA")
)

func buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSchema,
	)
}

func New(repoName string) (*Config, error) {
	name := os.Getenv("SERVER_NAME")
	if name == "" {
		name = repoName
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	return &Config{
		Server: Server{
			Name:    name,
			Host:    host,
			Port:    port,
			BaseURL: os.Getenv("SERVER_BASE_URL"),
		},
		DB: DBConfig{
			DSN: buildDSN(),
		},
	}, nil
}
