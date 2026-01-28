# Agent Commission Management System - Makefile
# Provides common development tasks and build automation

# Variables
APP_NAME=agent-commission
VERSION=1.0.0
BUILD_DIR=./build
MAIN_FILE=main.go
BINARY=$(BUILD_DIR)/$(APP_NAME)

# Database
DB_USER=postgres
DB_PASS=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME_DEV=agent_commission_dev
DB_NAME_TEST=agent_commission_test
DB_DSN_DEV=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_DEV)?sslmode=disable
DB_DSN_TEST=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TEST)?sslmode=disable

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down lint fmt generate

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "$(CYAN)Agent Commission Management System - Make Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

## Build Commands

build: ## Build the application
	@echo "$(CYAN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) $(MAIN_FILE)
	@echo "$(GREEN)✓ Build complete: $(BINARY)$(NC)"

build-prod: ## Build for production with optimizations
	@echo "$(CYAN)Building $(APP_NAME) for production...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o $(BINARY) $(MAIN_FILE)
	@echo "$(GREEN)✓ Production build complete$(NC)"

clean: ## Clean build artifacts
	@echo "$(CYAN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

## Run Commands

run: ## Run the application
	@echo "$(CYAN)Starting $(APP_NAME)...$(NC)"
	@APP_ENV=dev go run $(MAIN_FILE)

run-prod: build-prod ## Run production build
	@echo "$(CYAN)Starting $(APP_NAME) (production mode)...$(NC)"
	@APP_ENV=prod $(BINARY)

## Test Commands

test: ## Run unit tests
	@echo "$(CYAN)Running unit tests...$(NC)"
	@go test -v -race -short ./...

test-coverage: ## Run tests with coverage report
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

test-integration: ## Run integration tests (requires docker services)
	@echo "$(CYAN)Running integration tests...$(NC)"
	@APP_ENV=test go test -v -race -tags=integration ./test/...

test-e2e: ## Run end-to-end tests
	@echo "$(CYAN)Running E2E tests...$(NC)"
	@APP_ENV=test go test -v -tags=e2e ./test/...

benchmark: ## Run benchmarks
	@echo "$(CYAN)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

## Code Quality

lint: ## Run linter
	@echo "$(CYAN)Running linter...$(NC)"
	@golangci-lint run ./...

fmt: ## Format code
	@echo "$(CYAN)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...

tidy: ## Tidy go modules
	@echo "$(CYAN)Tidying go modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

generate: ## Generate code (validators, mocks, etc.)
	@echo "$(CYAN)Generating code...$(NC)"
	@go generate ./...
	@echo "$(GREEN)✓ Code generation complete$(NC)"

## Docker Commands

docker-up: ## Start docker services (PostgreSQL + Temporal + Redis)
	@echo "$(CYAN)Starting docker services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Temporal: localhost:7233"
	@echo "  - Temporal UI: http://localhost:8080"
	@echo "  - Redis: localhost:6379"

docker-up-dev: ## Start docker services with dev tools (includes PgAdmin)
	@echo "$(CYAN)Starting docker services with dev tools...$(NC)"
	@docker-compose --profile dev-tools up -d
	@echo "$(GREEN)✓ Services started with dev tools$(NC)"
	@echo "  - PgAdmin: http://localhost:5050"

docker-down: ## Stop docker services
	@echo "$(CYAN)Stopping docker services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

docker-down-v: ## Stop docker services and remove volumes (DANGER: Deletes all data)
	@echo "$(RED)WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		echo "$(GREEN)✓ Services stopped and volumes removed$(NC)"; \
	fi

docker-logs: ## View docker logs
	@docker-compose logs -f

docker-ps: ## Show running docker containers
	@docker-compose ps

## Database Commands

migrate-up: ## Apply all database migrations
	@echo "$(CYAN)Applying database migrations...$(NC)"
	@migrate -path db/migrations -database "$(DB_DSN_DEV)" up
	@echo "$(GREEN)✓ Migrations applied$(NC)"

migrate-down: ## Rollback last database migration
	@echo "$(CYAN)Rolling back last migration...$(NC)"
	@migrate -path db/migrations -database "$(DB_DSN_DEV)" down 1
	@echo "$(GREEN)✓ Migration rolled back$(NC)"

migrate-version: ## Show current migration version
	@migrate -path db/migrations -database "$(DB_DSN_DEV)" version

migrate-force: ## Force migration version (use carefully!)
	@read -p "Enter version to force: " version; \
	migrate -path db/migrations -database "$(DB_DSN_DEV)" force $$version

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=migration_name$(NC)"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir db/migrations -seq $(NAME)
	@echo "$(GREEN)✓ Migration created$(NC)"

db-reset: ## Reset database (drop all + reapply migrations) - DANGER!
	@echo "$(RED)WARNING: This will delete all database data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		migrate -path db/migrations -database "$(DB_DSN_DEV)" down && \
		migrate -path db/migrations -database "$(DB_DSN_DEV)" up && \
		echo "$(GREEN)✓ Database reset complete$(NC)"; \
	fi

db-shell: ## Open PostgreSQL shell
	@psql $(DB_DSN_DEV)

## Utility Commands

install-tools: ## Install development tools
	@echo "$(CYAN)Installing development tools...$(NC)"
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

deps: ## Download dependencies
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

verify: lint vet test ## Run all verification checks

check-env: ## Check if required environment tools are installed
	@echo "$(CYAN)Checking environment...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)✗ Go is not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ Go $(shell go version | awk '{print $$3}')$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)✗ Docker is not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ Docker $(shell docker --version | awk '{print $$3}')$(NC)"
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(RED)✗ Docker Compose is not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ Docker Compose $(shell docker-compose --version | awk '{print $$4}')$(NC)"
	@echo "$(GREEN)✓ Environment check passed$(NC)"

## Development Workflow

dev-setup: check-env docker-up migrate-up ## Complete development environment setup
	@echo "$(GREEN)✓ Development environment ready!$(NC)"
	@echo "$(CYAN)Next steps:$(NC)"
	@echo "  1. Run 'make run' to start the application"
	@echo "  2. Access http://localhost:8080/health"
	@echo "  3. View Temporal UI at http://localhost:8080"

dev-teardown: docker-down ## Tear down development environment
	@echo "$(GREEN)✓ Development environment torn down$(NC)"

all: clean build test ## Clean, build, and test

.PHONY: install-migrate
install-migrate: ## Install golang-migrate
	@echo "$(CYAN)Installing golang-migrate...$(NC)"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)✓ golang-migrate installed$(NC)"
