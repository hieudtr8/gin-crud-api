.PHONY: help api dev build build-legacy test test-coverage test-db test-graph \
        docker-up docker-down docker-build docker-run docker-restart docker-logs \
        docker-logs-api docker-logs-all docker-ps docker-rebuild \
        generate generate-graphql generate-ent clean format vet tidy

# Default target - show help
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

##@ General

help: ## Show this help message
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

docker-up: ## Start PostgreSQL with Docker Compose (database only, local)
	@echo "$(BLUE)Starting PostgreSQL (local)...$(NC)"
	docker-compose -f docker-compose.local.yml up -d postgres
	@echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(NC)"
	@sleep 3
	@echo "$(GREEN)✓ PostgreSQL is ready$(NC)"
	@echo "$(YELLOW)Connection: localhost:5432$(NC)"

docker-down: ## Stop all Docker containers
	@echo "$(BLUE)Stopping containers...$(NC)"
	docker-compose -f docker-compose.local.yml down 2>/dev/null || true
	docker-compose -f docker-compose.prod.yml down 2>/dev/null || true
	@echo "$(GREEN)✓ Containers stopped$(NC)"

docker-logs: ## Show PostgreSQL logs (local)
	docker-compose -f docker-compose.local.yml logs -f postgres

docker-build: ## Build Docker image for GraphQL API (local)
	@echo "$(BLUE)Building Docker image (local)...$(NC)"
	docker-compose -f docker-compose.local.yml build --no-cache graphql-api
	@echo "$(GREEN)✓ Docker image built$(NC)"

docker-build-prod: ## Build Docker image (production)
	@echo "$(BLUE)Building Docker image (production)...$(NC)"
	docker-compose -f docker-compose.prod.yml build --no-cache graphql-api
	@echo "$(GREEN)✓ Production Docker image built$(NC)"

docker-run: ## Start all services - LOCAL (dev environment)
	@echo "$(BLUE)Starting all services (LOCAL)...$(NC)"
	@echo "$(YELLOW)Building fresh image and starting containers...$(NC)"
	docker-compose -f docker-compose.local.yml up --build -d
	@echo "$(GREEN)✓ All services started (LOCAL)$(NC)"
	@echo "$(YELLOW)Environment: dev (configs/dev.yaml)$(NC)"
	@echo "$(YELLOW)GraphQL Playground: http://localhost:8081$(NC)"
	@echo "$(YELLOW)GraphQL API: http://localhost:8081/query$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"

docker-run-prod: ## Start all services - PRODUCTION (prod environment)
	@echo "$(RED)⚠️  PRODUCTION DEPLOYMENT$(NC)"
	@echo "$(BLUE)Starting all services (PRODUCTION)...$(NC)"
	@echo "$(YELLOW)Building fresh image and starting containers...$(NC)"
	docker-compose -f docker-compose.prod.yml up --build -d
	@echo "$(GREEN)✓ All services started (PRODUCTION)$(NC)"
	@echo "$(YELLOW)Environment: prod (configs/prod.yaml)$(NC)"
	@echo "$(RED)⚠️  Ensure DB_PASSWORD is set securely!$(NC)"
	@echo "$(YELLOW)GraphQL API: http://localhost:8081/query$(NC)"

docker-restart: ## Restart all services (local)
	@echo "$(BLUE)Restarting services (local)...$(NC)"
	docker-compose -f docker-compose.local.yml restart
	@echo "$(GREEN)✓ Services restarted$(NC)"

docker-restart-prod: ## Restart all services (production)
	@echo "$(BLUE)Restarting services (production)...$(NC)"
	docker-compose -f docker-compose.prod.yml restart
	@echo "$(GREEN)✓ Services restarted$(NC)"

docker-logs-api: ## Show GraphQL API logs (local)
	docker-compose -f docker-compose.local.yml logs -f graphql-api

docker-logs-api-prod: ## Show GraphQL API logs (production)
	docker-compose -f docker-compose.prod.yml logs -f graphql-api

docker-logs-all: ## Show all container logs (local)
	docker-compose -f docker-compose.local.yml logs -f

docker-ps: ## Show running containers
	@docker ps -a | grep gin_crud || echo "No containers found"

docker-rebuild: docker-down docker-run ## Rebuild and restart (local)

docker-rebuild-prod: docker-down docker-run-prod ## Rebuild and restart (production)

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
