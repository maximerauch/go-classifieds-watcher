# Binary output name
BINARY_NAME=bin/watcher

# Database connection string for integration tests (Docker Localhost)
TEST_DB_DSN="postgres://user:password@localhost:5432/watcher_db?sslmode=disable"

.PHONY: all build run test clean fmt lint audit help test-integration cover

# Set default target to help
default: help

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: Run the application via Docker Compose (Build + Up)
run:
	docker-compose up --build

## run-local: Run the application locally (without Docker, using go run)
run-local:
	go run cmd/watcher/main.go

## db-up: Start only the database container (useful for local dev)
db-up:
	docker-compose up -d db

## build: Compile the binary into the bin/ directory
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) cmd/watcher/main.go

## clean: Remove artifacts (binary, containers, coverage reports)
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	docker-compose down --remove-orphans

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## fmt: Format all code (goimports/gofmt)
fmt:
	go fmt ./...

## lint: Run the linter (golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

## test: Run unit tests (fast, no DB dependency)
test:
	@echo "Running unit tests..."
	go test -v -short ./...

## test-integration: Run integration tests (requires active DB container)
test-integration:
	@echo "Running integration tests..."
	DATABASE_URL=$(TEST_DB_DSN) go test -v -tags=integration ./internal/adapters/postgres/...

## cover: Generate and open coverage report (Core logic only)
cover:
	go test -coverprofile=coverage.out ./internal/core/...
	go tool cover -html=coverage.out
	rm coverage.out

## audit: Full check-up (Format + Lint + Unit Tests) - Run before push
audit: fmt lint test

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: Display this help message
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'