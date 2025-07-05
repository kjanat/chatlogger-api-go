# ChatLogger API - Development Makefile

.PHONY: help build test lint clean docker dev setup install-tools

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME := chatlogger-api
SERVER_BINARY := bin/server
WORKER_BINARY := bin/worker
COVERAGE_FILE := coverage.txt
DOCKER_COMPOSE_FILE := docker-compose.yml

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

## help: Display this help message
help:
	@echo "ChatLogger API Development Commands:"
	@echo ""
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//g' | column -t -s ':'

## setup: Install development tools and dependencies
setup: install-tools
	@echo "$(GREEN)Setting up development environment...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Development environment ready!$(NC)"

## install-tools: Install required development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "$(GREEN)Development tools installed!$(NC)"

## build: Build all binaries
build: build-server build-worker

## build-server: Build the API server
build-server:
	@echo "$(GREEN)Building server...$(NC)"
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $(SERVER_BINARY) ./cmd/server
	@echo "$(GREEN)Server built: $(SERVER_BINARY)$(NC)"

## build-worker: Build the background worker
build-worker:
	@echo "$(GREEN)Building worker...$(NC)"
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $(WORKER_BINARY) ./cmd/worker
	@echo "$(GREEN)Worker built: $(WORKER_BINARY)$(NC)"

## run: Run the API server locally
run: build-server
	@echo "$(GREEN)Starting API server...$(NC)"
	./$(SERVER_BINARY)

## run-worker: Run the background worker locally
run-worker: build-worker
	@echo "$(GREEN)Starting background worker...$(NC)"
	./$(WORKER_BINARY)

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -race ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## test-integration: Run integration tests with database
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		echo "$(YELLOW)TEST_DATABASE_URL not set, using Docker PostgreSQL...$(NC)"; \
		docker-compose up -d postgres redis; \
		sleep 5; \
		export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/chatlogger_test?sslmode=disable"; \
	fi
	go test -v -race -tags=integration ./...

## benchmark: Run performance benchmarks
benchmark:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	go test -bench=. -benchmem -run=^$$ ./internal/repository/ ./internal/service/
	@echo "$(GREEN)Benchmarks completed$(NC)"

## race: Run race condition tests
race:
	@echo "$(GREEN)Running race condition tests...$(NC)"
	go test -race -run="Race|Concurrent" -timeout=10m ./...

## lint: Run code linting
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run
	@echo "$(GREEN)Linting completed$(NC)"

## format: Format code
format:
	@echo "$(GREEN)Formatting code...$(NC)"
	gofmt -s -w .
	goimports -w .
	@echo "$(GREEN)Code formatted$(NC)"

## security: Run security checks
security:
	@echo "$(GREEN)Running security checks...$(NC)"
	govulncheck ./...
	@echo "$(GREEN)Security checks completed$(NC)"

## docs: Generate API documentation
docs:
	@echo "$(GREEN)Generating documentation...$(NC)"
	chmod +x ./scripts/docs_generate.sh
	./scripts/docs_generate.sh
	@echo "$(GREEN)Documentation generated$(NC)"

## clean: Clean build artifacts and cache
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf bin/
	rm -f $(COVERAGE_FILE) coverage.html
	go clean -cache -testcache -modcache
	@echo "$(GREEN)Cleaned$(NC)"

## docker: Build and start all services with Docker Compose
docker:
	@echo "$(GREEN)Starting services with Docker Compose...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --build

## docker-down: Stop all Docker services
docker-down:
	@echo "$(GREEN)Stopping Docker services...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

## docker-logs: Show Docker logs
docker-logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

## migration: Run database migrations
migration:
	@echo "$(GREEN)Running database migrations...$(NC)"
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "$(RED)DATABASE_URL not set$(NC)"; \
		exit 1; \
	fi
	psql $(DATABASE_URL) -f migrations/001_initial_schema.sql
	psql $(DATABASE_URL) -f migrations/002_ensure_defaults.sql
	psql $(DATABASE_URL) -f migrations/003_add_exports_table.sql
	@echo "$(GREEN)Migrations completed$(NC)"

## dev: Start development environment
dev: docker docs
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "API Server: http://localhost:8080"
	@echo "API Docs: http://localhost:8080/swagger/index.html"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

## ci: Run full CI pipeline locally
ci: clean lint test-coverage benchmark security
	@echo "$(GREEN)CI pipeline completed successfully!$(NC)"

## release: Prepare for release (format, lint, test, docs)
release: format lint test-coverage docs
	@echo "$(GREEN)Release preparation completed!$(NC)"

# Check if required tools are installed
check-tools:
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed. Run 'make install-tools'$(NC)" && exit 1)
	@which swag > /dev/null || (echo "$(RED)swag not installed. Run 'make install-tools'$(NC)" && exit 1)

# Development dependencies
deps:
	go mod download
	go mod verify
	go mod tidy