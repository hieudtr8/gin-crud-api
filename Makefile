.PHONY: help run-memory run-postgres docker-up docker-down migrate-up migrate-down build test clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make run-memory     - Run with in-memory storage"
	@echo "  make run-postgres   - Run with PostgreSQL storage"
	@echo "  make docker-up      - Start PostgreSQL with Docker"
	@echo "  make docker-down    - Stop PostgreSQL Docker"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"

# Run with in-memory storage
run-memory:
	STORAGE_TYPE=memory go run cmd/api/main.go

# Run with PostgreSQL storage
run-postgres:
	STORAGE_TYPE=postgres go run cmd/api/main.go

# Start PostgreSQL with Docker
docker-up:
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3

# Stop PostgreSQL Docker
docker-down:
	docker-compose down

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path migrations \
		-database "postgresql://postgres:postgres@localhost:5432/gin_crud_api?sslmode=disable" \
		up

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path migrations \
		-database "postgresql://postgres:postgres@localhost:5432/gin_crud_api?sslmode=disable" \
		down 1

# Build the application
build:
	go build -o gin-crud-api ./cmd/api

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f gin-crud-api
	go clean

# Development setup (docker + migrations)
dev-setup: docker-up migrate-up
	@echo "Development environment is ready!"

# Full cleanup
dev-cleanup: docker-down clean
	@echo "Cleaned up development environment"