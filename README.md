# 🚀 Gin CRUD API - GraphQL

A clean architecture **GraphQL API** built with Go, gqlgen, and EntGo ORM for type-safe database operations with schema-first development.

## 📋 Table of Contents
- [Overview](#-overview)
- [Quick Start](#-quick-start)
- [GraphQL Examples](#-graphql-examples)
- [Development Commands](#-development-commands)
- [Architecture](#-architecture)
- [Testing](#-testing)

## 🎯 Overview

A **GraphQL API** for managing departments and employees featuring:
- **GraphQL with gqlgen**: Schema-first development with interactive Playground
- **EntGo ORM**: Type-safe database operations with automatic migrations
- **PostgreSQL**: Production-ready database with foreign key constraints
- **Clean architecture**: Repository pattern with dependency injection
- **Modular schema**: Organized by entity for easy scaling
- **Docker support**: One-command deployment

### Entity Relationship
```
Department (1) ──────→ (many) Employee
    ├── id (UUID)           ├── id (UUID)
    ├── name                ├── name
    ├── created_at          ├── email (validated)
    └── updated_at          ├── department_id (FK)
                            ├── created_at
                            └── updated_at
```

## ⚡ Quick Start

### Prerequisites
- Go 1.24.3+
- Docker (for PostgreSQL)
- Make

### Local Development

```bash
# 1. Setup environment
make setup              # Start PostgreSQL + download dependencies

# 2. Run GraphQL server
make api                # Server starts on http://localhost:8081

# 3. Open GraphQL Playground
# Navigate to http://localhost:8081 and try:
query {
  health
}
```

### Docker Deployment

```bash
# One command to start everything (PostgreSQL + API)
make docker-run

# Or use docker-compose directly
docker-compose up -d
```

Access at **http://localhost:8081**

### Environment Configuration

Create `.env` file:
```bash
# GraphQL Server
GRAPHQL_PORT=8081

# PostgreSQL (use "postgres" for Docker, "localhost" for local)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_crud_api
```

**Note**: When using Docker, `DB_HOST` should be `postgres` (handled automatically by docker-compose).

## 🎮 GraphQL Examples

### Create a Department
```graphql
mutation {
  createDepartment(input: {name: "Engineering"}) {
    id
    name
  }
}
```

### Query Employees with Department Info
```graphql
query {
  employees {
    id
    name
    email
    department {
      id
      name
    }
  }
}
```

### Update an Employee
```graphql
mutation {
  updateEmployee(
    id: "your-emp-id"
    input: {
      name: "Jane Doe"
      email: "jane@example.com"
      departmentID: "your-dept-id"
    }
  ) {
    id
    name
    email
  }
}
```

### Delete a Department (Cascades to Employees)
```graphql
mutation {
  deleteDepartment(id: "your-dept-id")
}
```

## 🛠 Development Commands

```bash
# Development
make api              # Run GraphQL server
make dev              # Alias for 'make api'
make test             # Run all tests
make help             # Show all available commands

# Code Generation (after schema changes)
make generate         # Regenerate all code (EntGo + GraphQL)
make generate-graphql # Regenerate GraphQL code only
make generate-ent     # Regenerate EntGo database code only

# Docker & Database
make docker-run       # Start all services (PostgreSQL + API)
make docker-build     # Build Docker image
make docker-down      # Stop all containers
make docker-logs-api  # Show API logs
make docker-ps        # Show running containers

# Build
make build            # Build production binary
make clean            # Remove build artifacts

# Code Quality
make format           # Format Go code
make vet              # Run go vet
make tidy             # Tidy go modules
```

## 🏗 Architecture

### Request Flow
```
GraphQL Query/Mutation
    ↓
gqlgen Server (validates against schema)
    ↓
Resolver (processes operation, validates input)
    ↓
Repository (CRUD operations)
    ↓
EntGo Client (generates SQL)
    ↓
PostgreSQL
```

### Key Technologies

**gqlgen** - GraphQL code generation for Go
- Define API in `.graphql` files
- Auto-generates type-safe Go code
- Built-in GraphQL Playground

**EntGo** - Modern ORM for Go
- Define database schema in Go code
- Type-safe query builders
- Automatic migrations (no manual SQL!)

### Project Structure

```
cmd/
├── graphql/main.go              # GraphQL server entry point ⭐
└── legacy/rest_main.go          # Legacy REST (for reference)

internal/
├── graph/                       # GraphQL Layer ⭐
│   ├── schema/                  # GraphQL schemas by entity
│   │   ├── common.graphql       # Base Query/Mutation types
│   │   ├── department.graphql   # Department schema (EDIT THIS!)
│   │   └── employee.graphql     # Employee schema (EDIT THIS!)
│   ├── department.resolvers.go  # Department resolvers (EDIT THIS!)
│   ├── employee.resolvers.go    # Employee resolvers (EDIT THIS!)
│   ├── common.resolvers.go      # Health check resolver
│   ├── validation.go            # Helper functions (email validation)
│   ├── resolver.go              # Dependency injection
│   ├── generated.go             # Generated code (DO NOT EDIT!)
│   └── model/
│       └── models_gen.go        # Generated GraphQL models (DO NOT EDIT!)
│
├── ent/                         # EntGo ORM Layer
│   ├── schema/
│   │   ├── department.go        # Department schema (EDIT THIS!)
│   │   └── employee.go          # Employee schema (EDIT THIS!)
│   └── ...                      # Generated files (DO NOT EDIT!)
│
├── database/                    # Repository Layer
│   ├── database.go              # Repository interfaces
│   ├── ent_client.go            # EntGo setup & migrations
│   ├── ent_department_repo.go   # Department repository
│   └── ent_employee_repo.go     # Employee repository
│
└── config/                      # Configuration
    └── config.go                # Environment loader

gqlgen.yml                       # gqlgen configuration
docker-compose.yml               # Docker orchestration
```

### Files You Should Edit

**GraphQL Schemas** (split by entity):
- `internal/graph/schema/common.graphql` - Base types
- `internal/graph/schema/department.graphql` - Department entity
- `internal/graph/schema/employee.graphql` - Employee entity

**Resolvers** (auto-generated, implementations preserved):
- `internal/graph/department.resolvers.go` - Department operations
- `internal/graph/employee.resolvers.go` - Employee operations

**Database Schemas**:
- `internal/ent/schema/department.go`
- `internal/ent/schema/employee.go`

**Helper Functions**:
- `internal/graph/validation.go` - Validation helpers

### Files You Should NOT Edit

- `internal/graph/generated.go` - Fully regenerated by gqlgen
- `internal/graph/model/models_gen.go` - Generated GraphQL models
- `internal/ent/*.go` (except `schema/`) - Generated by EntGo

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/database -v
go test ./internal/graph -v
```

### Test Helpers

```go
// Create test database
client := testutil.NewTestEntClient(t)

// Seed test data
dept := testutil.SeedTestDepartment(t, client, "Engineering")
emp := testutil.SeedTestEmployee(t, client, "John", "john@example.com", dept.ID)
```

## 📝 Common Tasks

### Adding a New Field

1. **Update GraphQL Schema** (`internal/graph/schema/employee.graphql`):
```graphql
type Employee {
  id: ID!
  name: String!
  email: String!
  phone: String  # New field
  departmentID: ID!
}

input CreateEmployeeInput {
  name: String!
  email: String!
  phone: String  # Add here too
  departmentID: ID!
}
```

2. **Update EntGo Schema** (`internal/ent/schema/employee.go`):
```go
field.String("phone").Optional()
```

3. **Regenerate & Restart**:
```bash
make generate
make api  # EntGo auto-migrates the new column!
```

### Adding a New Entity

1. **Create GraphQL Schema** (`internal/graph/schema/project.graphql`):
```graphql
type Project {
  id: ID!
  name: String!
}

extend type Query {
  projects: [Project!]!
}

extend type Mutation {
  createProject(input: CreateProjectInput!): Project!
}
```

2. **Create EntGo Schema**:
```bash
go run entgo.io/ent/cmd/ent new --target internal/ent/schema Project
```

3. **Regenerate**:
```bash
make generate  # Auto-creates project.resolvers.go
```

4. **Implement repository** in `internal/database/ent_project_repo.go`

5. **Update resolver** (`internal/graph/resolver.go`) with dependency injection

6. **Implement resolver logic** in auto-generated `internal/graph/project.resolvers.go`

## 🔑 Key Features

- **Schema-first GraphQL**: Define API in `.graphql`, generate type-safe code
- **Modular organization**: Schemas split by entity for scalability
- **Auto-migrations**: Database schema updates automatically
- **Type safety**: Compile-time checking for GraphQL and database
- **Repository pattern**: Clean separation of concerns
- **Dependency injection**: Easy testing with mock repositories
- **Docker deployment**: Multi-stage build (~15MB image)
- **Interactive playground**: Built-in API explorer
- **Comprehensive testing**: 78.5% repository coverage

## 🚨 Important Notes

1. **After schema changes**: Always run `make generate-graphql` or `make generate-ent`
2. **Resolver files**: gqlgen preserves your implementation code during regeneration
3. **Helper functions**: Put them in `validation.go`, not in resolver files
4. **GraphQL models**: Used throughout the entire app (no separate domain models)
5. **Auto-migrations**: EntGo handles database schema changes on startup
6. **Docker vs Local**: Use `DB_HOST=postgres` for Docker, `localhost` for local

## 📚 Resources

- **[gqlgen Documentation](https://gqlgen.com/)** - GraphQL code generation
- **[EntGo Documentation](https://entgo.io/)** - ORM and migrations
- **[GraphQL Playground](http://localhost:8081)** - Interactive API explorer (when server is running)
- **Legacy REST API**: See `internal/legacy/rest/README.md` for comparison

---

**Made with ❤️ using Go, GraphQL, and EntGo**
