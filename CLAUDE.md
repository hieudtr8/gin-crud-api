# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **GraphQL API** built with Go, gqlgen, and EntGo ORM for managing Departments and Employees. The project demonstrates clean architecture with:
- **Schema-first GraphQL** with gqlgen (single source of truth)
- **Type-safe database operations** with EntGo ORM
- **Repository pattern** with dependency injection
- **Structured logging** with zerolog
- **Comprehensive testing** (78.5% repository coverage)
- **Automatic schema migrations**
- **Legacy REST API** preserved in `internal/legacy/` for reference

### Single Source of Truth Architecture

**GraphQL Schema** (`internal/graph/schema.graphql`) is the **ONLY** source of truth:
- Defines all types, queries, and mutations
- gqlgen auto-generates Go types (`internal/graph/model/models_gen.go`)
- Input types (CreateDepartmentInput, etc.) are auto-generated
- Output types (Department, Employee) are auto-generated
- **These generated types are used throughout the entire application** - in repositories, resolvers, and handlers
- No separate domain models - GraphQL models serve all layers

**Legacy REST DTOs** (`internal/legacy/models.go`) are isolated:
- Only used for REST API request binding (CreateDepartmentRequest, etc.)
- Legacy REST handlers also use GraphQL models for data operations
- Kept for backward compatibility reference

## Commands

All commands are available via Makefile for convenience. Run `make help` to see all available commands.

### Quick Start
```bash
# First time setup
make setup            # Start PostgreSQL + download dependencies
make api              # Run GraphQL server (dev mode)

# Access GraphQL Playground
# Open http://localhost:8081 in your browser
```

### Development
```bash
# Run servers
make api              # Run GraphQL server (primary)
make dev              # Alias for 'make api'
make legacy           # Run legacy REST server (reference)

# Code generation
make generate         # Regenerate all code (EntGo + GraphQL)
make generate-ent     # Regenerate EntGo code only
make generate-graphql # Regenerate GraphQL code only

# Code quality
make format           # Format Go code
make vet              # Run go vet
make tidy             # Tidy go modules
```

### Building
```bash
# Build for production
make build            # Build GraphQL server → ./build/graphql-server
make build-legacy     # Build legacy REST server → ./build/rest-server

# Run built binary
./build/graphql-server
```

### Testing
```bash
# Run tests
make test             # Run all tests
make test-coverage    # Run tests with coverage
make test-db          # Run database repository tests only
make test-graph       # Run GraphQL resolver tests only
```

### Docker & Database
```bash
# PostgreSQL management
make docker-up        # Start PostgreSQL with Docker Compose
make docker-down      # Stop PostgreSQL
make docker-logs      # Show PostgreSQL logs
```

### Utilities
```bash
# Cleanup
make clean            # Remove build artifacts
make clean-all        # Stop containers + remove all artifacts

# Dependencies
make deps             # Download Go dependencies

# Help
make help             # Show all available commands
```

## Architecture

### Simplified Architecture
The application uses GraphQL-generated models throughout all layers:

1. **GraphQL Layer** (`internal/graph/`): Schema definition and auto-generated types
2. **Resolver Layer** (`internal/graph/schema.resolvers.go`): Business logic and validation
3. **Repository Layer** (`internal/database/`): Data access using GraphQL models
4. **EntGo Layer** (`internal/ent/`): Type-safe database operations
5. **Database Layer**: PostgreSQL with automatic schema management

### Dependency Flow
```
GraphQL Schema → gqlgen (generates types) → Used by Resolvers & Repositories → EntGo Client → PostgreSQL
```
- **Single type system**: GraphQL-generated models used everywhere
- No conversion between layers - same types from API to database
- All dependencies injected via constructors
- Repositories work directly with `*model.Department` and `*model.Employee`
- EntGo converts between GraphQL models and database entities
- Automatic database migrations on startup

### GraphQL Architecture

The project uses **gqlgen** (schema-first approach):
- **Schema Definition** (`internal/graph/schema.graphql`): Define GraphQL API contract
- **Code Generation**: Run `go run github.com/99designs/gqlgen generate` to create resolver stubs
- **Resolver Implementation** (`internal/graph/schema.resolvers.go`): Implement business logic
- **Dependency Injection** (`internal/graph/resolver.go`): Inject repositories into resolvers
- **Type Safety**: Compile-time errors for schema mismatches

### Key Design Patterns

1. **Repository Pattern**: All data access goes through repository interfaces (`DepartmentRepository`, `EmployeeRepository`)
2. **Dependency Injection**: Constructor-based injection throughout
3. **Code Generation**: Both gqlgen and EntGo use code generation for type safety
4. **Input Types**: GraphQL Input types (`CreateDepartmentInput`, `UpdateDepartmentInput`, etc.) for mutations
5. **Schema-First GraphQL**: GraphQL schema defined in `.graphql` files, then code generated
6. **Schema-First Database**: Database schema defined in Go code (`internal/ent/schema/`)

### EntGo Architecture

EntGo uses a code-generation approach:
- **Schema Definition** (`internal/ent/schema/*.go`): Define database structure in Go
- **Code Generation**: Run `go generate ./internal/ent` to create client code
- **Generated Client** (`internal/ent/client.go`): Type-safe database operations
- **Auto-Migrations**: Schema changes applied automatically on startup
- **Type Safety**: Compile-time errors for invalid queries

## Implementation Status

### Fully Implemented GraphQL Operations

**Queries:**
- `department(id: ID!)` - Get department by ID
- `departments` - List all departments
- `employee(id: ID!)` - Get employee by ID
- `employees` - List all employees
- `employeesByDepartment(departmentID: ID!)` - Get employees by department
- `health` - Health check

**Mutations:**
- `createDepartment(input: CreateDepartmentInput!)` - Create department
- `updateDepartment(id: ID!, input: UpdateDepartmentInput!)` - Update department
- `deleteDepartment(id: ID!)` - Delete department (cascade deletes employees)
- `createEmployee(input: CreateEmployeeInput!)` - Create employee (validates department exists)
- `updateEmployee(id: ID!, input: UpdateEmployeeInput!)` - Update employee
- `deleteEmployee(id: ID!)` - Delete employee

All CRUD operations are fully functional with PostgreSQL persistence via EntGo.

**GraphQL Playground:** Access at `http://localhost:8081/` to explore the schema and test queries/mutations interactively.

## Data Storage

Uses PostgreSQL with EntGo ORM:
- Data persisted in PostgreSQL database
- EntGo generates type-safe database operations
- Automatic schema migrations on startup
- Foreign key constraints enforced at database level
- Connection pooling with configurable limits
- Timestamps (`created_at`, `updated_at`) managed automatically

### EntGo Features in Use
- **Type-Safe Queries**: Compile-time error checking
- **Auto-Migrations**: No manual SQL migration files needed
- **Relationships**: One-to-Many (Department → Employees) with foreign keys
- **Generated CRUD**: All basic operations auto-generated
- **Query Builders**: Fluent API for complex queries

## Business Rules

1. **Employee Creation**: Must reference an existing department (validated in handler)
2. **IDs**: Generated as UUIDs using google/uuid package
3. **Validation**:
   - Department name is required
   - Employee name and email are required
   - Employee email must be valid format
   - Employee must belong to existing department

## Common Tasks

### Adding a New Field to Existing GraphQL Type

1. **Update GraphQL Schema** (`internal/graph/schema.graphql`):
   ```graphql
   type Employee {
     id: ID!
     name: String!
     email: String!
     phone: String  # Add new field
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

4. **Update Resolvers** (`internal/graph/schema.resolvers.go`):
   Add phone field to conversion logic in resolvers

5. **Restart Application**: EntGo auto-migrates the new column!

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

5. **Create Repository Interface** in `internal/database/database.go`

6. **Implement Repository** in `internal/database/ent_project_repo.go`

7. **Update Resolver Struct** (`internal/graph/resolver.go`) to inject new repository

8. **Implement Resolvers** in `internal/graph/schema.resolvers.go`

### Modifying GraphQL Schema
- After ANY schema change, ALWAYS run: `make generate-graphql`
- gqlgen updates generated code and resolver stubs
- Your resolver implementations are preserved during regeneration

### Modifying EntGo Schema
- After ANY EntGo schema change, ALWAYS run: `make generate-ent`
- Or run `make generate` to regenerate both GraphQL and EntGo code
- Changes are applied automatically on next app restart
- EntGo compares schema to database and applies diffs

### Adding Validation
- GraphQL input validation happens in resolvers (`internal/graph/schema.resolvers.go`)
- Add validation logic at the start of mutation resolvers
- Database constraints defined in EntGo schema (e.g., `.NotEmpty()`, `.Unique()`)

## Project Structure
```
cmd/
├── graphql/main.go                      - GraphQL server entry point
└── legacy/rest_main.go                  - Legacy REST server (for reference)
internal/
├── config/config.go                     - Environment configuration (DB, server, logging)
├── logger/logger.go                     - Structured logging with zerolog
├── middleware/logging.go                - GraphQL logging middleware
├── testutil/testutil.go                 - Test utilities and helpers
├── graph/                               - GraphQL layer
│   ├── schema.graphql                   - GraphQL schema (SINGLE SOURCE OF TRUTH - EDIT THIS!)
│   ├── resolver.go                      - Resolver struct with dependency injection
│   ├── schema.resolvers.go              - Resolver implementations (EDIT THIS!)
│   ├── validation_test.go               - Email validation tests
│   ├── generated.go                     - Generated GraphQL execution code (DO NOT EDIT!)
│   └── model/
│       ├── models_gen.go                - Generated GraphQL models (DO NOT EDIT!)
│       │                                  These types are used EVERYWHERE in the app
│       └── model.go                     - Custom GraphQL models (if needed)
├── database/
│   ├── database.go                      - Repository interfaces & ErrNotFound
│   ├── ent_client.go                    - EntGo client setup & migrations
│   ├── ent_department_repo.go           - Department repository (uses GraphQL models)
│   ├── ent_department_repo_test.go      - Department repository tests (13 tests)
│   ├── ent_employee_repo.go             - Employee repository (uses GraphQL models)
│   ├── ent_employee_repo_test.go        - Employee repository tests (18 tests)
│   └── legacy/                          - Legacy in-memory/postgres repositories
├── ent/                                 - EntGo generated code (DO NOT EDIT MANUALLY!)
│   ├── schema/
│   │   ├── department.go                - Department schema definition (EDIT THIS!)
│   │   └── employee.go                  - Employee schema definition (EDIT THIS!)
│   ├── client.go                        - Generated database client
│   ├── department.go, department_*.go   - Generated department code
│   ├── employee.go, employee_*.go       - Generated employee code
│   └── ...                              - Other generated files (~20+ files)
└── legacy/                              - Legacy code (for reference)
    ├── models.go                        - Legacy REST Request DTOs (only for REST)
    └── rest/                            - Legacy REST API implementation
        ├── department/handler.go        - REST department handlers
        ├── employee/handler.go          - REST employee handlers
        └── router/router.go             - REST routing configuration
```

### Important Directories

**Files You SHOULD Edit:**
- **`internal/graph/schema.graphql`**: Define your GraphQL API here (manually edit)
- **`internal/graph/schema.resolvers.go`**: Implement GraphQL resolvers here (manually edit)
  - **Special Note**: This file says "will be automatically regenerated" BUT gqlgen uses **intelligent partial regeneration**
  - Your function implementations are preserved across regenerations
  - Only function signatures and new stubs are updated by gqlgen
  - Think of it as: gqlgen manages the signatures, you write the logic
- **`internal/ent/schema/`**: Define your database schema here (manually edit)
- **`internal/graph/resolver.go`**: Dependency injection setup (manually edit)

**Files You Should NEVER Edit:**
- **`internal/graph/generated.go`**: Fully generated by gqlgen - completely rewritten on regeneration
- **`internal/graph/model/models_gen.go`**: Fully generated by gqlgen - completely rewritten
- **`internal/ent/`** (all files except `schema/`): Fully generated by EntGo - completely rewritten
- **`internal/legacy/rest/`**: Old REST implementation kept for learning purposes

## Testing

The project has comprehensive testing infrastructure:

### Test Coverage
- **51 total tests** with **78.5% coverage** of repository layer
- Repository tests: 31 tests (13 department + 18 employee)
- Validation tests: 20 tests (email format validation)
- Fast execution: ~1.5 seconds total

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/database -v
go test ./internal/graph -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Utilities
Use `internal/testutil` for creating test data:
```go
client := testutil.NewTestEntClient(t)  // In-memory SQLite
dept := testutil.SeedTestDepartment(t, client, "Engineering")
emp := testutil.SeedTestEmployee(t, client, "John", "john@example.com", dept.ID)
```

## Logging

The project uses structured logging with zerolog:

### Log Levels
- `DEBUG`: Detailed operation traces (use `LOG_LEVEL=debug`)
- `INFO`: Operation start/completion
- `WARN`: Expected issues (e.g., "not found")
- `ERROR`: System failures requiring attention

### Configuration
Set in `.env`:
```bash
LOG_LEVEL=info       # debug, info, warn, error
LOG_PRETTY=true      # Pretty console output for development
```

### Request Tracing
Every GraphQL operation gets a unique request ID:
- Propagates through middleware → resolver → repository
- Enables tracing in high-concurrency scenarios
- Check logs with: `request_id` field

## Notes

- **GraphQL models are used everywhere** - no separate domain models
- Vietnamese comments in code explain architecture decisions
- GraphQL server port configured via `GRAPHQL_PORT` environment variable (default: 8081)
- Legacy REST server port configured via `SERVER_PORT` environment variable (default: 8080)
- No authentication/authorization implemented yet
- **IMPORTANT Code Generation**:
  - After modifying GraphQL schema: `make generate-graphql`
  - After modifying EntGo schema: `make generate-ent`
  - Or regenerate both: `make generate`
- Database migrations happen automatically on application startup
- Both gqlgen and EntGo use code generation (~3000+ lines generated total)
- GraphQL Playground provides interactive API documentation and testing
- `ErrNotFound` moved to `internal/database` package (used by all repositories)

## GraphQL Examples

### Query Examples (use in Playground at http://localhost:8081)

```graphql
# Get all departments
query {
  departments {
    id
    name
  }
}

# Get department with employees
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

# Get all employees with their department info
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

### Mutation Examples

```graphql
# Create a department
mutation {
  createDepartment(input: {name: "Engineering"}) {
    id
    name
  }
}

# Create an employee
mutation {
  createEmployee(input: {
    name: "John Doe"
    email: "john@example.com"
    departmentID: "your-dept-id"
  }) {
    id
    name
    email
  }
}

# Update an employee
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

# Delete a department (cascades to employees)
mutation {
  deleteDepartment(id: "your-dept-id")
}
```