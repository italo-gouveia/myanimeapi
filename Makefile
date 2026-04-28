# MyAnimeAPI Makefile
# Useful commands for development and testing

.PHONY: help build run test test-unit test-integration test-e2e test-coverage test-mocks clean docker-build docker-run

# Variables
BINARY_NAME=myanimeapi
BUILD_DIR=build
MAIN_PATH=cmd/main.go

# Default command
.DEFAULT_GOAL := help

# Help
help: ## Shows this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build
build: ## Compiles the project
	@echo "Compiling MyAnimeAPI..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Run
run: ## Runs the project
	@echo "Running MyAnimeAPI..."
	go run $(MAIN_PATH)

# Tests
test: ## Runs all tests
	@echo "Running all tests..."
	go test -v ./tests/...

test-unit: ## Runs only unit tests
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

test-integration: ## Runs only integration tests
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-e2e: ## Runs only E2E tests
	@echo "Running E2E tests..."
	go test -v ./tests/e2e/...

test-coverage: ## Runs tests with coverage
	@echo "Running tests with coverage..."
	@TEST_PKGS=$$(go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./... 2>/dev/null | tr '\n' ' '); \
	if [ -z "$$TEST_PKGS" ]; then \
		echo "No test files found, skipping coverage"; \
	else \
		go test -v -coverprofile=coverage.out $$TEST_PKGS; \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "Coverage report generated: coverage.html"; \
	fi

test-mocks: ## Generates mocks for tests
	@echo "Generating mocks..."
	@if command -v mockgen > /dev/null; then \
		mockgen -source=api/repositories/repository.go -destination=tests/mocks/repo_mocks.go; \
		mockgen -source=api/services/service.go -destination=tests/mocks/service_mocks.go; \
		echo "Mocks generated successfully!"; \
	else \
		echo "mockgen not found. Install with: go install github.com/golang/mock/mockgen@latest"; \
	fi

test-parallel: ## Runs tests in parallel
	@echo "Running tests in parallel..."
	go test -v -parallel 4 ./tests/...

test-specific: ## Runs specific test (use TEST_NAME=test_name)
	@if [ -z "$(TEST_NAME)" ]; then \
		echo "Error: TEST_NAME not specified. Use: make test-specific TEST_NAME=TestName"; \
		exit 1; \
	fi
	@echo "Running test: $(TEST_NAME)"
	go test -v -run $(TEST_NAME) ./tests/...

# Cleanup
clean: ## Cleans build and test files
	@echo "Cleaning files..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "Cleanup completed!"

# Docker
docker-build: ## Builds Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .
	@echo "Docker image built: $(BINARY_NAME)"

docker-run: ## Runs Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Dependencies
deps: ## Installs project dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Linting and formatting
lint: ## Runs linter on code
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Formats the code
	@echo "Formatting code..."
	go fmt ./...
	go vet ./...

# Migrations
migrate-up: ## Runs database migrations (via the application itself at startup)
	@echo "Migrations are embedded and run automatically at application startup."
	@echo "To inspect migration files: ls internal/database/migrations/"

# Swagger
swagger: ## Generates Swagger documentation
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/main.go --dir ./cmd,./api/adapters/http,./api/models,./internal/errors -o cmd/docs; \
		echo "Swagger documentation generated!"; \
	else \
		echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Development
dev: ## Runs in development mode with hot reload
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		go run $(MAIN_PATH); \
	fi

# Quality checks
quality: lint fmt test ## Runs quality checks (lint, fmt, test)
	@echo "Quality checks completed!"

# Initial setup
setup: deps test-mocks ## Initial project setup
	@echo "Initial setup completed!"

# Test Environment
test-env-up: ## Starts test environment (Docker)
	@echo "Starting test environment..."
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for services to be ready..."
	@echo "PostgreSQL: localhost:5433"
	@echo "Redis: localhost:6380"
	@echo "MinIO: localhost:9001"

test-env-down: ## Stops test environment (Docker)
	@echo "Stopping test environment..."
	docker-compose -f docker-compose.test.yml down

test-env-logs: ## Shows test environment logs
	@echo "Test environment logs:"
	docker-compose -f docker-compose.test.yml logs

test-setup: test-env-up ## Sets up complete test environment
	@echo "Setting up test environment variables..."
	@echo "export TEST_DB_HOST=localhost"
	@echo "export TEST_DB_PORT=5433"
	@echo "export TEST_DB_USER=test"
	@echo "export TEST_DB_PASSWORD=test"
	@echo "export TEST_DB_NAME=test_db"
	@echo "Test environment configured!"

test-db: ## Starts Postgres and runs all tests (including DB-dependent ones)
	@echo "Starting Postgres test container..."
	docker-compose -f docker-compose.test.yml up -d postgres-test
	@echo "Waiting for Postgres to be healthy..."
	@until docker exec myanimeapi-postgres-test pg_isready -U test -d test_db > /dev/null 2>&1; do \
		echo "  ...waiting"; sleep 2; \
	done
	@echo "Postgres is ready. Running all tests..."
	TEST_DB_HOST=localhost TEST_DB_PORT=5433 TEST_DB_USER=test TEST_DB_PASSWORD=test TEST_DB_NAME=test_db \
		go test -v -count=1 ./tests/...
	@echo "Stopping Postgres test container..."
	docker-compose -f docker-compose.test.yml stop postgres-test

# CI/CD
ci: quality test-coverage ## Runs CI pipeline
	@echo "CI pipeline completed!"
