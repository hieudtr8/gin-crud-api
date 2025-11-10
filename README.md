# 🚀 Gin CRUD API - GraphQL Edition

A clean architecture **GraphQL API** built with Go, gqlgen, and EntGo ORM for type-safe database operations with schema-first development.

## 📋 Table of Contents
- [Overview](#overview)
- [Quick Start](#quick-start)
- [GraphQL Examples](#graphql-examples)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Development Guide](#development-guide)
- [Testing](#testing)
- [Key Concepts](#key-concepts)

## 🎯 Overview

This project is a **GraphQL API** for managing departments and employees with:
- **GraphQL with gqlgen**: Schema-first GraphQL code generation
- **Interactive Playground**: Explore and test the API in your browser
- **Full CRUD operations**: Queries and Mutations for both entities
- **EntGo ORM**: Type-safe database operations with automatic migrations
- **Clean architecture**: Repository pattern with dependency injection
- **PostgreSQL**: Production-ready database with automatic schema management
- **Foreign key constraints**: Database-enforced relationships
- **Input validation**: Email format checking and referential integrity
- **Type safety**: Compile-time errors for both GraphQL and database
- **Docker support**: Easy PostgreSQL setup
- **Legacy REST API preserved**: Compare implementations in `internal/legacy/`

### ✨ Key Features

🚀 **GraphQL Advantages:**
- Single endpoint for all operations (`/query`)
- Request exactly the data you need
- No over-fetching or under-fetching
- Interactive documentation via Playground
- Strong typing with auto-generated code
- Nested queries for related data

### Entity Relationship
```
Department (1) ──────→ (many) Employee
    ├── id (UUID)           ├── id (UUID)
    ├── name                ├── name
    ├── created_at          ├── email (unique)
    └── updated_at          ├── department_id (FK)
                            ├── created_at
                            └── updated_at
```

**Note**: Timestamps are automatically managed by EntGo

## ⚡ Quick Command Reference

```bash
# Most common commands
make setup       # First time setup (PostgreSQL + dependencies)
make api         # Run GraphQL server (dev mode)
make test        # Run all tests
make generate    # Regenerate code after schema changes
make help        # Show all available commands
```

For detailed commands, see the [Development Guide](#-development-guide) section.

## 🚀 Quick Start

### Prerequisites
- Go 1.24.3+ installed
- Docker (for PostgreSQL)
- Make (usually pre-installed on macOS/Linux)
- Web browser (for GraphQL Playground)

### 1. Clone and Install
```bash
cd gin-crud-api
make deps
```

### 2. Configure Environment

Create a `.env` file:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_crud_api

# Server Configuration
GRAPHQL_PORT=8081
```

### 3. Start PostgreSQL

```bash
make docker-up
```

<details>
<summary>Alternative: Manual Docker Setup</summary>

```bash
docker run --name postgres-gin \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gin_crud_api \
  -p 5432:5432 \
  -d postgres:15-alpine
```
</details>

### 4. Run the GraphQL Server

```bash
make api
```

Or use the alias:
```bash
make dev
```

You should see:
```
╔═══════════════════════════════════════════════════════════╗
║  GraphQL Server is running!                               ║
╠═══════════════════════════════════════════════════════════╣
║  Playground:  http://localhost:8081/                     ║
║  GraphQL API: http://localhost:8081/query                ║
╠═══════════════════════════════════════════════════════════╣
║  Database:    PostgreSQL (via EntGo)                      ║
║  Schema:      internal/graph/schema.graphql               ║
╚═══════════════════════════════════════════════════════════╝
```

**What happens on startup?**
- ✅ Connects to PostgreSQL database
- ✅ Runs automatic schema migrations
- ✅ Creates/updates database tables
- ✅ Starts GraphQL server on port 8081
- ✅ GraphQL Playground available at http://localhost:8081

**No manual migrations needed!** EntGo handles everything automatically.

### 5. Open GraphQL Playground

Navigate to **http://localhost:8081** in your browser and try this query:

```graphql
query {
  health
}
```

You should see:
```json
{
  "data": {
    "health": "ok"
  }
}
```

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

### Create an Employee

```graphql
mutation {
  createEmployee(input: {
    name: "Alice Smith"
    email: "alice@example.com"
    departmentID: "YOUR_DEPARTMENT_ID_HERE"
  }) {
    id
    name
    email
  }
}
```

### Query All Departments

```graphql
query {
  departments {
    id
    name
  }
}
```

### Query Department with Employees (Nested Query!)

```graphql
query {
  department(id: "your-dept-id") {
    id
    name
    employees {
      id
      name
      email
    }
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

### Get Employees by Department

```graphql
query {
  employeesByDepartment(departmentID: "your-dept-id") {
    id
    name
    email
  }
}
```

## 🏗 Architecture

### How GraphQL Requests Flow

```
GraphQL Query/Mutation
    ↓
gqlgen Server → Parses and validates against schema
    ↓
GraphQL Resolver → Processes operation, validates input
    ↓
Repository Interface → Defines CRUD operations
    ↓
EntGo Repository → Type-safe database operations
    ↓
EntGo Client → Generates SQL queries
    ↓
PostgreSQL → Stores/retrieves data
```

### Why This Architecture?

1. **Schema-First GraphQL**: Define API contract first, generate type-safe code
2. **Double Type Safety**: Both gqlgen and EntGo provide compile-time safety
3. **Automatic Migrations**: Database schema changes applied automatically
4. **No Raw SQL**: Write Go code, EntGo generates optimized SQL
5. **Self-Documenting API**: GraphQL schema serves as live documentation
6. **Separation of Concerns**: Each layer has one specific responsibility
7. **Easy Testing**: Can swap real database with mock repositories
8. **Flexible Queries**: Clients request exactly what they need

### What is gqlgen?

**gqlgen** is a Go library for building GraphQL APIs using schema-first development:
- Define API in `.graphql` files using GraphQL schema language
- Automatically generates type-safe Go code
- Generates resolver interfaces - you just implement the logic
- Built-in support for GraphQL Playground
- Compile-time type checking for all operations

**Benefits**:
- ✅ **Schema-First** - API contract defined in standard GraphQL SDL
- ✅ **Type-Safe Resolvers** - Compile-time errors for mismatches
- ✅ **Auto-Generated Code** - Less boilerplate, more productivity
- ✅ **Built-in Playground** - Interactive API explorer
- ✅ **Flexible Queries** - Clients request exactly what they need

### What is EntGo?

**EntGo** is a modern ORM (Object-Relational Mapping) framework for Go that:
- Defines database schemas in Go code (not SQL)
- Automatically generates type-safe database client code
- Handles migrations without manual SQL scripts
- Provides query builders that catch errors at compile-time
- Manages relationships between entities automatically

**Benefits**:
- ✅ **No SQL writing needed** - Define schemas in Go
- ✅ **Compile-time safety** - Catch bugs before runtime
- ✅ **Auto-migrations** - Schema changes applied automatically
- ✅ **Better IDE support** - Autocomplete for queries
- ✅ **Less boilerplate** - Generated CRUD operations

## 📁 Project Structure

```
cmd/
├── graphql/main.go              # GraphQL server entry point ⭐ PRIMARY
└── legacy/rest_main.go          # Legacy REST server (reference)

internal/
├── graph/                       # GraphQL Layer ⭐
│   ├── schema.graphql           # GraphQL schema (EDIT THIS!)
│   ├── schema.resolvers.go      # Resolver implementations (EDIT THIS!)
│   ├── resolver.go              # Dependency injection
│   ├── generated.go             # Generated GraphQL code (DO NOT EDIT!)
│   └── model/
│       ├── models_gen.go        # Generated GraphQL models (DO NOT EDIT!)
│       └── model.go             # Custom GraphQL models (if needed)
│
├── ent/                         # EntGo ORM Layer
│   ├── schema/
│   │   ├── department.go        # Department schema (EDIT THIS!)
│   │   └── employee.go          # Employee schema (EDIT THIS!)
│   ├── client.go                # Generated database client (DO NOT EDIT!)
│   ├── department*.go           # Generated department code (DO NOT EDIT!)
│   ├── employee*.go             # Generated employee code (DO NOT EDIT!)
│   └── ...                      # Other generated files
│
├── database/                    # Repository Layer
│   ├── database.go              # Repository interfaces & ErrNotFound
│   ├── ent_client.go            # EntGo client setup & migrations
│   ├── ent_department_repo.go   # Department repository (uses GraphQL models)
│   ├── ent_employee_repo.go     # Employee repository (uses GraphQL models)
│   └── legacy/                  # Legacy in-memory/postgres repositories
│
├── config/                      # Configuration
│   └── config.go                # Environment configuration loader
│
└── legacy/rest/                 # Legacy REST API (for reference)
    ├── README.md                # REST API documentation
    ├── department/handler.go    # REST department handlers
    ├── employee/handler.go      # REST employee handlers
    └── router/router.go         # REST routing configuration

gqlgen.yml                       # gqlgen configuration
docker-compose.yml               # PostgreSQL setup
.env                             # Configuration (don't commit!)
```

### Important Files

**Files You SHOULD Edit:**
- **`internal/graph/schema.graphql`**: Define your GraphQL API here (SINGLE SOURCE OF TRUTH!)
- **`internal/graph/schema.resolvers.go`**: Implement resolver logic here (gqlgen preserves your code!)
- **`internal/ent/schema/`**: Define your database schema here
- **`internal/graph/resolver.go`**: Dependency injection setup
- **`internal/database/*_repo.go`**: Repository implementations

**Files You Should NEVER Edit:**
- **`internal/graph/generated.go`**: Fully generated by gqlgen (DO NOT EDIT!)
- **`internal/graph/model/models_gen.go`**: Generated GraphQL models - **USED EVERYWHERE IN APP** (DO NOT EDIT!)
- **`internal/ent/*.go`** (except `/schema/`): Fully generated by EntGo (DO NOT EDIT!)

**Note on `schema.resolvers.go`**: This file has a special regeneration behavior - when you run `gqlgen generate`, it updates function signatures but **preserves your implementation code**. This is different from fully-generated files that get completely rewritten.

**Other Important Files:**
- **`internal/legacy/rest/README.md`**: REST API documentation and comparison
- **`gqlgen.yml`**: gqlgen configuration

## 🛠 Development Guide

### Quick Command Reference

```bash
# Development
make api              # Run GraphQL server (dev mode)
make dev              # Alias for 'make api'
make legacy           # Run legacy REST server (reference)

# Code Generation
make generate         # Regenerate all code (EntGo + GraphQL)
make generate-ent     # Regenerate EntGo database code only
make generate-graphql # Regenerate GraphQL code only

# Docker & Database
make docker-up        # Start PostgreSQL
make docker-down      # Stop PostgreSQL
make docker-logs      # Show PostgreSQL logs

# Testing
make test             # Run all tests
make test-coverage    # Run tests with coverage
make test-db          # Run database tests only
make test-graph       # Run GraphQL tests only

# Build
make build            # Build GraphQL server for production
make build-legacy     # Build legacy REST server

# Code Quality
make format           # Format Go code
make vet              # Run go vet
make tidy             # Tidy go modules

# Utilities
make clean            # Remove build artifacts
make clean-all        # Stop containers + remove artifacts
make setup            # Setup dev environment (PostgreSQL + deps)
make help             # Show all available commands
```

### First Time Setup

```bash
make setup            # Starts PostgreSQL + downloads dependencies
make api              # Run the server
```

### Adding a New Field to GraphQL Type

1. **Update GraphQL Schema** (`internal/graph/schema.graphql`):
   ```graphql
   type Employee {
     id: ID!
     name: String!
     email: String!
     phone: String  # Add new field
     departmentID: ID!
   }

   # Don't forget to update input types too!
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

3. **Regenerate Code**:
   ```bash
   make generate
   ```

4. **Update Repository** (`internal/database/ent_employee_repo.go`):
   Add phone field to EntGo entity conversion:
   ```go
   SetPhone(emp.Phone)  // In Save and Update methods
   ```

5. **Restart Application**:
   ```bash
   make api
   ```
   EntGo auto-migrates the new column on startup!


### Adding a New GraphQL Type/Entity

1. **Define GraphQL Schema** (`internal/graph/schema.graphql`):
   ```graphql
   type Project {
     id: ID!
     name: String!
     description: String
   }

   input CreateProjectInput {
     name: String!
     description: String
   }

   extend type Query {
     project(id: ID!): Project
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

3. **Define Fields in EntGo Schema** (`internal/ent/schema/project.go`)

4. **Regenerate Code**:
   ```bash
   make generate
   ```

5. **Create Repository Interface** in `internal/database/database.go`:
   ```go
   type ProjectRepository interface {
       Save(project *model.Project) error
       FindByID(id string) (*model.Project, error)
       FindAll() ([]*model.Project, error)
       // ... other methods
   }
   ```
   **Note**: Use `*model.Project` (GraphQL type), not a separate domain model!

6. **Implement Repository** in `internal/database/ent_project_repo.go`:
   - Accept and return GraphQL models (`*model.Project`)
   - Convert between GraphQL models and EntGo entities internally

7. **Update Resolver Struct** (`internal/graph/resolver.go`) to inject new repository

8. **Implement Resolvers** in `internal/graph/schema.resolvers.go`:
   - Resolvers work directly with GraphQL models
   - No conversion logic needed - just pass models to/from repository!

### Configuration (.env file)

```bash
# GraphQL Server
GRAPHQL_PORT=8081

# PostgreSQL Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_crud_api
DB_SSLMODE=disable

# Connection Pool
DB_MAX_CONNS=25
DB_MIN_CONNS=5
```

## 🧪 Testing

The project uses Go's standard `testing` package with testify for assertions and enttest for in-memory database testing.

### Test Structure

We have **51 total tests** with **78.5% coverage** of the repository layer:

- **Repository Tests** (31 tests): Test all CRUD operations with real database (in-memory SQLite)
- **Validation Tests** (20 tests): Test email format validation with various edge cases
- **Integration Tests**: Use enttest with SQLite for fast, isolated testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Run specific package tests
make test-db          # Database repository tests
make test-graph       # GraphQL resolver tests

# Run specific test (use go test directly)
go test ./internal/database -run TestEntDepartmentRepo_Save -v
```

### Test Coverage by Package

```
internal/database:  78.5% (repository layer - all CRUD operations)
internal/graph:     0.3%  (validation function only)
Total:              3.1%  (includes EntGo-generated code)
```

**Note**: The low total coverage is due to EntGo-generated code. Our **business logic (repositories) has 78.5% coverage**.

### Test Utilities

We provide test helpers in `internal/testutil/`:

```go
// Create in-memory test database
client := testutil.NewTestEntClient(t)

// Seed test data
dept := testutil.SeedTestDepartment(t, client, "Engineering")
emp := testutil.SeedTestEmployee(t, client, "John", "john@example.com", dept.ID)

// Seed multiple records
depts := testutil.SeedMultipleDepartments(t, client, []string{"Sales", "HR"})
emps := testutil.SeedMultipleEmployees(t, client, deptID, 5)
```

### Writing Tests

Example repository test:

```go
func TestEntDepartmentRepo_Save(t *testing.T) {
    // Setup: Create in-memory database
    client := testutil.NewTestEntClient(t)
    defer client.Close()
    repo := NewEntDepartmentRepo(client)

    // Test: Save a department (using GraphQL model!)
    dept := &model.Department{
        ID:   uuid.New().String(),
        Name: "Engineering",
    }
    err := repo.Save(dept)

    // Assert: No error and department was saved
    require.NoError(t, err)

    // Verify: Can retrieve it (returns GraphQL model)
    saved, err := repo.FindByID(dept.ID)
    require.NoError(t, err)
    assert.Equal(t, dept.Name, saved.Name)
}
```

**Note**: Tests use `model.Department` from `internal/graph/model`, not a separate domain model.

### Test Organization

Tests follow Go conventions:
- Test files: `*_test.go` alongside source files
- Test functions: `Test*` with `*testing.T` parameter
- Table-driven tests for validation with multiple inputs
- Setup/teardown with `defer` for cleanup

### Testing Best Practices Used

1. **In-Memory Database**: Fast tests with real SQL operations using SQLite
2. **Test Isolation**: Each test gets its own database instance
3. **Seeding Helpers**: Reusable functions for creating test data
4. **Table-Driven Tests**: Validation tests cover 20 different email formats
5. **Clear Naming**: Test names describe what they're testing
6. **Assertions**: Using testify for readable assertions

## 🔰 Key Concepts

### 1. Schema-First GraphQL

Define your API in `.graphql` files, then generate type-safe code:

```graphql
# internal/graph/schema.graphql
type Department {
  id: ID!
  name: String!
  employees: [Employee!]
}

input CreateDepartmentInput {
  name: String!
}

type Mutation {
  createDepartment(input: CreateDepartmentInput!): Department!
}
```

Then run: `make generate-graphql`

### 2. Type-Safe Resolvers

gqlgen generates resolver interfaces, you implement the logic:

```go
// internal/graph/schema.resolvers.go
func (r *mutationResolver) CreateDepartment(
    ctx context.Context,
    input model.CreateDepartmentInput,
) (*model.Department, error) {
    // Validate input
    if input.Name == "" {
        return nil, errors.New("department name is required")
    }

    // Create GraphQL model (used throughout entire app!)
    dept := &model.Department{
        ID:   uuid.New().String(),
        Name: input.Name,
    }

    // Save via repository (accepts GraphQL model directly)
    if err := r.DeptRepo.Save(dept); err != nil {
        return nil, err
    }

    // Return the same model - no conversion needed!
    return dept, nil
}
```

**Key difference**: GraphQL models are used everywhere - no conversion between layers!

### 3. Repository Pattern

Repositories use GraphQL models directly (no separate domain models):

```go
// internal/database/database.go
type DepartmentRepository interface {
    Save(dept *model.Department) error        // Uses GraphQL model!
    FindByID(id string) (*model.Department, error)
    FindAll() ([]*model.Department, error)
    Update(dept *model.Department) error
    Delete(id string) error
}

// internal/database/ent_department_repo.go
type EntDepartmentRepo struct {
    client *ent.Client
}

// Repositories convert between GraphQL models and EntGo entities internally
func (r *EntDepartmentRepo) Save(dept *model.Department) error {
    // Convert GraphQL model → EntGo entity → Save to database
    _, err := r.client.Department.Create().
        SetID(uuid.MustParse(dept.ID)).
        SetName(dept.Name).
        Save(ctx)
    return err
}
```

### 4. Dependency Injection

```go
// Inject repositories into resolver
resolver := graph.NewResolver(deptRepo, empRepo)

// This makes testing easy - just pass mock repositories!
```

### 5. EntGo Schema-First Approach

```go
// Define database schema in Go (internal/ent/schema/department.go)
func (Department) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.String("name").NotEmpty(),
        field.Time("created_at").Default(time.Now),
    }
}

// EntGo automatically generates:
// - SQL CREATE TABLE statements
// - Type-safe query builders
// - CRUD methods
```

After changing schemas, run: `make generate-ent`

## 🚨 Important Notes

1. **Two Code Generators**: gqlgen (GraphQL) and EntGo (Database) - both use schema-first approach
2. **UUID IDs**: All IDs are UUIDs (universally unique identifiers)
3. **Auto-Migrations**: EntGo automatically creates/updates database tables on startup
4. **Foreign Key Constraints**: Database enforces relationships
5. **Timestamps**: `created_at` and `updated_at` managed automatically by EntGo
6. **Type Safety**: Compile-time checking for both GraphQL and database operations
7. **Validation**: Email format, required fields, referential integrity all enforced
8. **GraphQL Playground**: Provides interactive API documentation
9. **Legacy REST API**: Available in `internal/legacy/rest/` for comparison

## 🎯 Key Technologies

- **[gqlgen](https://gqlgen.com/)** - GraphQL code generation for Go
- **[EntGo](https://entgo.io/)** - Type-safe database ORM with code generation
- **[PostgreSQL](https://www.postgresql.org/)** - Relational database
- **[Docker](https://www.docker.com/)** - Containerization for PostgreSQL
