# amartha-reconciliation-service

-include .env
export

MIGRATION_DIR := schema/reconciliation

DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USERNAME ?= postgres
DB_PASSWORD ?= postgres
DB_NAME     ?= postgres

PSQL := PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USERNAME) -d $(DB_NAME) -v ON_ERROR_STOP=1

.PHONY: all build run test clean watch \
	db-up db-down docker-run docker-down \
	migrate setup env-setup itest

all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

env-setup:
	@if [ ! -f .env ]; then cp .env-dist .env && echo "Created .env from .env-dist"; else echo ".env already exists"; fi

db-up: env-setup
	@if docker compose up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi
	@echo "Postgres starting on $(DB_HOST):$(DB_PORT)..."

db-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

docker-run: db-up

docker-down: db-down

migrate: env-setup
	@echo "Applying migrations from $(MIGRATION_DIR)..."
	@command -v psql >/dev/null 2>&1 || { echo "psql is required; install PostgreSQL client tools"; exit 1; }
	@for f in $$(ls $(MIGRATION_DIR)/*.sql 2>/dev/null | sort); do \
		echo "==> $$f"; \
		$(PSQL) -f "$$f"; \
	done
	@echo "Migrations complete."

setup: db-up
	@echo "Waiting for Postgres..."
	@sleep 3
	@$(MAKE) migrate

db-reset: db-down
	docker volume rm amartha-reconciliation-service_psql_volume_bp 2>/dev/null || true
	$(MAKE) setup

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

clean:
	@echo "Cleaning..."
	@rm -f main

watch:
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Go's 'air' is not installed. Install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
		else \
			exit 1; \
		fi; \
	fi
