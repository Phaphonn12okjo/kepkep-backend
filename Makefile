# Load .env if it exists, and pass its values to every command.
-include .env
export

MIGRATIONS_DIR := db/migrations

.PHONY: help tools run build test test-short lint fmt migrate-up migrate-down migrate-status migrate-create sqlc

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

tools: ## Install dev tools (macOS, Homebrew)
	brew install goose sqlc golangci-lint

run: ## Start the API
	go run ./cmd/api

build: ## Build the API binary into bin/
	go build -o bin/api ./cmd/api

test: ## Unit + integration tests (needs Docker for testcontainers)
	go test -race -count=1 ./...

test-short: ## Unit tests only (skips tests that call testing.Short)
	go test -race -short ./...

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	golangci-lint fmt

migrate-up: ## Apply all migrations
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down: ## Roll back the last migration
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status: ## Show migration status
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-create: ## Create a migration: make migrate-create name=create_users
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

sqlc: ## Regenerate query code
	sqlc generate
