.PHONY: help api dev build build-legacy test test-coverage test-db test-graph \
        docker-up docker-down generate generate-graphql generate-ent \
        clean format vet tidy

# Default target - show help
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

##@ General

help: ## Show this help message
	@echo '$(BLUE)Available commands:$(NC)'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-18s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

api: ## Run GraphQL API server (dev mode)
	@echo "$(BLUE)Starting GraphQL API server...$(NC)"
	go run ./cmd/graphql

dev: api ## Alias for 'make api'

legacy: ## Run legacy REST API server (reference)
	@echo "$(YELLOW)Starting legacy REST API server...$(NC)"
	go run ./cmd/legacy

##@ Build

build: ## Build GraphQL server for production
	@echo "$(BLUE)Building GraphQL server...$(NC)"
	go build -o ./build/graphql-server ./cmd/graphql
	@echo "$(GREEN)✓ Built: ./build/graphql-server$(NC)"

build-legacy: ## Build legacy REST server
	@echo "$(BLUE)Building legacy REST server...$(NC)"
	go build -o ./build/rest-server ./cmd/legacy
	@echo "$(GREEN)✓ Built: ./build/rest-server$(NC)"

##@ Testing

test: ## Run all tests
	@echo "$(BLUE)Running all tests...$(NC)"
	go test ./...

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test ./... -cover

test-coverage-html: ## Run tests and generate HTML coverage report
	@echo "$(BLUE)Generating HTML coverage report...$(NC)"
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

test-db: ## Run database repository tests only
	@echo "$(BLUE)Running database tests...$(NC)"
	go test ./internal/database -v

test-graph: ## Run GraphQL resolver tests only
	@echo "$(BLUE)Running GraphQL tests...$(NC)"
	go test ./internal/graph -v

##@ Code Generation

generate: generate-ent generate-graphql ## Regenerate all code (EntGo + GraphQL)

generate-ent: ## Regenerate EntGo database code
	@echo "$(BLUE)Regenerating EntGo code...$(NC)"
	go generate ./internal/ent
	@echo "$(GREEN)✓ EntGo code regenerated$(NC)"

generate-graphql: ## Regenerate GraphQL code (gqlgen)
	@echo "$(BLUE)Regenerating GraphQL code...$(NC)"
	go run github.com/99designs/gqlgen generate
	@echo "$(GREEN)✓ GraphQL code regenerated$(NC)"

##@ Docker & Database

docker-up: ## Start PostgreSQL with Docker Compose
	@echo "$(BLUE)Starting PostgreSQL...$(NC)"
	docker-compose up -d
	@echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(NC)"
	@sleep 3
	@echo "$(GREEN)✓ PostgreSQL is ready$(NC)"

docker-down: ## Stop PostgreSQL Docker containers
	@echo "$(BLUE)Stopping PostgreSQL...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ PostgreSQL stopped$(NC)"

docker-logs: ## Show PostgreSQL logs
	docker-compose logs -f postgres

##@ Code Quality

format: ## Format Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	go vet ./...
	@echo "$(GREEN)✓ Vet completed$(NC)"

tidy: ## Tidy go modules
	@echo "$(BLUE)Tidying modules...$(NC)"
	go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

##@ Cleanup

clean: ## Remove build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf ./build
	go clean
	@echo "$(GREEN)✓ Cleaned$(NC)"

clean-all: clean docker-down ## Stop containers and remove all artifacts
	@echo "$(GREEN)✓ Full cleanup complete$(NC)"

##@ Setup

deps: ## Download Go dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

setup: docker-up deps ## Setup development environment (PostgreSQL + deps)
	@echo "$(GREEN)✓ Development environment ready!$(NC)"
	@echo "$(YELLOW)Run 'make api' to start the GraphQL server$(NC)"

##@ Database Operations (EntGo handles migrations automatically)

# Note: EntGo automatically migrates schema on application startup
# No manual migration commands needed!
db-info: ## Show database migration info
	@echo "$(YELLOW)EntGo Auto-Migration Info:$(NC)"
	@echo "  • Migrations run automatically on server startup"
	@echo "  • Schema changes detected from internal/ent/schema/"
	@echo "  • No manual migration files needed"
	@echo ""
	@echo "$(BLUE)To apply schema changes:$(NC)"
	@echo "  1. Edit files in internal/ent/schema/"
	@echo "  2. Run: make generate-ent"
	@echo "  3. Run: make api (migrations auto-apply on startup)"
