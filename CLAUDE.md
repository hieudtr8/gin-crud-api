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

## Configuration System

The project uses **Viper** for configuration management with environment-specific YAML files and environment variable overrides.

### Configuration Architecture

1. **YAML Config Files** (`configs/`): Environment-specific default configurations
   - `dev.yaml` - Development environment (localhost, debug logging)
   - `prod.yaml` - Production environment (postgres host, info logging, SSL required)
   - `test.yaml` - Test environment (reduced logging, lower resource limits)

2. **Environment Variables**: Override YAML values with `GINAPI_` prefix
   - Format: `GINAPI_<SECTION>_<KEY>` (e.g., `GINAPI_DATABASE_HOST=localhost`)
   - Automatically bind to nested YAML keys (`.` becomes `_`)

3. **Configuration Priority** (highest to lowest):
   - Environment variables (`GINAPI_*`)
   - YAML config file (`configs/{APP_ENV}.yaml`)
   - Built-in defaults in code

### Environment Selection

Set `APP_ENV` to choose the configuration file:

```bash
# Development (default)
make api                    # Uses configs/dev.yaml
APP_ENV=dev go run cmd/graphql/main.go

# Production
APP_ENV=prod go run cmd/graphql/main.go

# Test
APP_ENV=test go test ./...
```

### Environment Variables Reference

**Format**: All environment variables use the `GINAPI_` prefix with sections separated by underscores.

**Server Configuration:**
- `GINAPI_SERVER_GRAPHQL_PORT` - GraphQL API port (default: 8081)
- `GINAPI_SERVER_REST_PORT` - Legacy REST API port (default: 8080)

**Database Configuration:**
- `GINAPI_DATABASE_HOST` - PostgreSQL host (dev: localhost, prod: postgres)
- `GINAPI_DATABASE_PORT` - PostgreSQL port (default: 5432)
- `GINAPI_DATABASE_USER` - Database username (default: postgres)
- `GINAPI_DATABASE_PASSWORD` - Database password (⚠️ **Always override in production**)
- `GINAPI_DATABASE_DBNAME` - Database name (default: gin_crud_api)
- `GINAPI_DATABASE_SSLMODE` - SSL mode (dev: disable, prod: require)
- `GINAPI_DATABASE_MAX_CONNS` - Max connection pool size (dev: 25, prod: 50)
- `GINAPI_DATABASE_MIN_CONNS` - Min idle connections (dev: 5, prod: 10)

**Logging Configuration:**
- `GINAPI_LOGGING_LEVEL` - Log level: debug, info, warn, error (dev: debug, prod: info)
- `GINAPI_LOGGING_PRETTY` - Pretty console output (dev: true, prod: false)

### Configuration Examples

**Development (using dev.yaml defaults):**
```bash
# Uses configs/dev.yaml - no env vars needed
make api
```

**Development with overrides:**
```bash
# Override specific values
GINAPI_LOGGING_LEVEL=debug \
GINAPI_DATABASE_PASSWORD=secret \
make api
```

**Production (Docker):**
```bash
# In docker-compose.yml
environment:
  APP_ENV: prod                           # Load configs/prod.yaml
  GINAPI_DATABASE_HOST: postgres          # Override for Docker networking
  GINAPI_DATABASE_PASSWORD: ${DB_PASSWORD} # Secure password from env
```

**Local custom configuration:**
```bash
# Create configs/local.yaml (gitignored)
# Then run with:
APP_ENV=local make api
```

For complete details, see `configs/README.md`.

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
# Full stack deployment
make docker-run       # Start all services (PostgreSQL + GraphQL API)
make docker-build     # Build Docker image
make docker-rebuild   # Rebuild and restart everything

# PostgreSQL management
make docker-up        # Start PostgreSQL only
make docker-down      # Stop all containers

# Monitoring
make docker-logs-api  # Show API logs
make docker-logs      # Show PostgreSQL logs
make docker-logs-all  # Show all logs
make docker-ps        # Show running containers
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

1. **Create New GraphQL Schema File** (`internal/graph/schema/project.graphql`):
   ```graphql
   # Project Schema - Project entity and related operations

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
   This will create `internal/graph/project.resolvers.go` automatically

5. **Create Repository Interface** in `internal/database/database.go`

6. **Implement Repository** in `internal/database/ent_project_repo.go`

7. **Update Resolver Struct** (`internal/graph/resolver.go`) to inject new repository

8. **Implement Resolvers** in `internal/graph/project.resolvers.go` (auto-generated by gqlgen)

### Modifying GraphQL Schema
- Edit schema files in `internal/graph/schema/` directory
- After ANY schema change, ALWAYS run: `make generate-graphql`
- gqlgen updates generated code and resolver stubs
- Your resolver implementations are preserved during regeneration
- Resolver files (e.g., `department.resolvers.go`) are named after schema files (e.g., `department.graphql`)

### Modifying EntGo Schema
- After ANY EntGo schema change, ALWAYS run: `make generate-ent`
- Or run `make generate` to regenerate both GraphQL and EntGo code
- Changes are applied automatically on next app restart
- EntGo compares schema to database and applies diffs

### Adding Validation
- GraphQL input validation happens in resolver files (e.g., `internal/graph/employee.resolvers.go`)
- Add validation logic at the start of mutation resolvers
- Helper functions should go in `internal/graph/validation.go`
- gqlgen will warn you if you put helper functions in resolver files
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
│   ├── schema/                          - GraphQL schema files (EDIT THESE!)
│   │   ├── common.graphql               - Base Query/Mutation types + health check
│   │   ├── department.graphql           - Department schema, inputs, queries, mutations
│   │   └── employee.graphql             - Employee schema, inputs, queries, mutations
│   ├── resolver.go                      - Resolver struct with dependency injection
│   ├── common.resolvers.go              - Common resolver implementations (health check)
│   ├── department.resolvers.go          - Department resolver implementations (EDIT THIS!)
│   ├── employee.resolvers.go            - Employee resolver implementations (EDIT THIS!)
│   ├── validation.go                    - Helper functions (email validation, etc.)
│   ├── validation_test.go               - Validation tests
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
- **`internal/graph/schema/`**: Define your GraphQL API here (manually edit schema files)
  - **`common.graphql`**: Base Query/Mutation types
  - **`department.graphql`**: Department entity schema
  - **`employee.graphql`**: Employee entity schema
  - Use `extend type Query` and `extend type Mutation` for entity-specific operations
- **`internal/graph/*.resolvers.go`**: Implement GraphQL resolvers here (manually edit)
  - **`department.resolvers.go`**: Department operations
  - **`employee.resolvers.go`**: Employee operations
  - **Special Note**: These files say "will be automatically regenerated" BUT gqlgen uses **intelligent partial regeneration**
  - Your function implementations are preserved across regenerations
  - Only function signatures and new stubs are updated by gqlgen
  - Think of it as: gqlgen manages the signatures, you write the logic
- **`internal/graph/validation.go`**: Helper functions like email validation (manually edit)
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
Set in YAML config files (`configs/dev.yaml`, `configs/prod.yaml`) or override with environment variables:
```bash
# Via environment variables
GINAPI_LOGGING_LEVEL=debug   # debug, info, warn, error
GINAPI_LOGGING_PRETTY=true   # Pretty console output for development

# Or in configs/dev.yaml
logging:
  level: debug
  pretty: true
```

### Request Tracing
Every GraphQL operation gets a unique request ID:
- Propagates through middleware → resolver → repository
- Enables tracing in high-concurrency scenarios
- Check logs with: `request_id` field

## Docker Deployment

The project is fully dockerized for easy deployment:

### Quick Deploy
```bash
# Start everything (PostgreSQL + GraphQL API)
make docker-run

# Or directly
docker-compose up -d
```

Access at: http://localhost:8081

### Docker Architecture
- **Multi-stage Dockerfile**: Optimized build (~15MB final image)
- **Docker Compose**: Orchestrates PostgreSQL + GraphQL API
- **Health Checks**: Built-in for both services
- **Auto-restart**: Services recover automatically
- **Networking**: Services communicate via app-network
- **Volumes**: PostgreSQL data persisted

### Environment Variables

**Configuration Selection:**
- `APP_ENV`: Environment selector (dev, prod, test) - loads corresponding YAML file

**Configuration Overrides** (use `GINAPI_` prefix):
- `GINAPI_SERVER_GRAPHQL_PORT`: GraphQL server port (default from YAML)
- `GINAPI_DATABASE_HOST`: Database host (`postgres` for Docker, `localhost` for local)
- `GINAPI_DATABASE_PORT`, `GINAPI_DATABASE_USER`, `GINAPI_DATABASE_PASSWORD`, `GINAPI_DATABASE_DBNAME`: Database config
- `GINAPI_LOGGING_LEVEL`: Logging level (debug, info, warn, error)
- `GINAPI_LOGGING_PRETTY`: Pretty logs (true for dev, false for production)

See the **Configuration System** section above for complete details.

### Docker Files
- `Dockerfile`: Multi-stage build configuration
- `docker-compose.yml`: Service orchestration
- `.dockerignore`: Excludes unnecessary files from build
- `configs/`: YAML configuration files (dev.yaml, prod.yaml, test.yaml)

## Notes

- **GraphQL models are used everywhere** - no separate domain models
- **Docker deployment available** - One command to run entire stack
- **Viper configuration** - Environment-specific YAML files with env var overrides
- Vietnamese comments in code explain architecture decisions
- Configuration via `APP_ENV` variable (dev/prod/test) + YAML files + `GINAPI_*` env vars
- GraphQL server port: `GINAPI_SERVER_GRAPHQL_PORT` (default: 8081 from YAML)
- Legacy REST server port: `GINAPI_SERVER_REST_PORT` (default: 8080 from YAML)
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