.PHONY: help dev build test test-integration migrate clean swagger

help: ## Show this help message
	@echo "Available commands:"
	@echo "  make dev              - Start development environment"
	@echo "  make build            - Build the binary"
	@echo "  make test             - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make migrate          - Run database migrations"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make clean            - Clean build artifacts"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@/Users/tarun/go/bin/swag init -g cmd/server/main.go

dev: ## Start development environment
	docker-compose up -d
	@echo "Waiting for postgres..."
	@sleep 3
	/usr/local/go/bin/go run cmd/server/main.go

build: ## Build the binary
	@echo "Building BaaS..."
	@/usr/local/go/bin/go build -o bin/baas cmd/server/main.go

test: ## Run unit tests
	@echo "Running tests..."
	@/usr/local/go/bin/go test -v ./...

test-integration: ## Run integration tests
	docker-compose up -d
	@sleep 3
	go test -v -tags=integration ./tests/integration/...
	docker-compose down

migrate: ## Run database migrations
	@echo "Migrations will be run automatically on startup"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
