# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a REST API CRUD application built with Go, Gin web framework, and EntGo ORM for managing Departments and Employees. The project demonstrates clean architecture with dependency injection, repository pattern, type-safe database operations with EntGo, and automatic schema migrations. Full CRUD operations are implemented using PostgreSQL as the data store.

## Commands

### Build and Run
```bash
# Build the application
go build -o /build/gin-crud-api ./cmd/api

# Run the application (starts on port 8080)
./gin-crud-api

# Build and run in one step
go run ./cmd/api
```

### Development
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Format code
go fmt ./...

# Vet code
go vet ./...

# Regenerate EntGo code (after schema changes)
go generate ./internal/ent
```

### Testing
```bash
# Run tests (currently no tests exist)
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Architecture

### Layered Architecture
The application follows a clean architecture pattern with these layers:

1. **HTTP Layer** (`internal/router/`): Gin router setup with middleware
2. **Handler Layer** (`internal/department/`, `internal/employee/`): HTTP request handling and validation
3. **Repository Layer** (`internal/database/`): Repository interfaces and EntGo implementations
4. **EntGo Layer** (`internal/ent/`): Generated type-safe database client
5. **Database Layer**: PostgreSQL with automatic schema management

### Dependency Flow
```
main.go → EntGo Client → EntGo Repositories → Handlers → Router
```
- All dependencies are injected via constructors
- Handlers depend on repository interfaces, not concrete implementations
- EntGo client provides type-safe database operations
- Automatic migrations on application startup

### Key Design Patterns

1. **Repository Pattern**: All data access goes through repository interfaces (`DepartmentRepository`, `EmployeeRepository`)
2. **Dependency Injection**: Constructor-based injection throughout
3. **Code Generation**: EntGo generates ~20+ files with type-safe CRUD operations
4. **DTO Pattern**: Separate request DTOs (`Create*Request`, `Update*Request`) with validation tags
5. **Schema-First**: Database schema defined in Go code (`internal/ent/schema/`)

### EntGo Architecture

EntGo uses a code-generation approach:
- **Schema Definition** (`internal/ent/schema/*.go`): Define database structure in Go
- **Code Generation**: Run `go generate ./internal/ent` to create client code
- **Generated Client** (`internal/ent/client.go`): Type-safe database operations
- **Auto-Migrations**: Schema changes applied automatically on startup
- **Type Safety**: Compile-time errors for invalid queries

## Implementation Status

### Fully Implemented Endpoints
- `POST /api/v1/departments` - Create department
- `GET /api/v1/departments` - List all departments
- `GET /api/v1/departments/:id` - Get department by ID
- `PUT /api/v1/departments/:id` - Update department
- `DELETE /api/v1/departments/:id` - Delete department
- `POST /api/v1/employees` - Create employee (validates department exists)
- `GET /api/v1/employees` - List all employees
- `GET /api/v1/employees/:id` - Get employee by ID
- `PUT /api/v1/employees/:id` - Update employee
- `DELETE /api/v1/employees/:id` - Delete employee
- `GET /health` - Health check

All CRUD operations are fully functional with PostgreSQL persistence via EntGo.

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

### Adding a New Field to Existing Entity

1. **Update EntGo Schema** (`internal/ent/schema/employee.go`):
   ```go
   field.String("phone").Optional()
   ```

2. **Regenerate EntGo Code**:
   ```bash
   go generate ./internal/ent
   ```

3. **Update Domain Model** (`internal/models/models.go`):
   ```go
   Phone string `json:"phone,omitempty"`
   ```

4. **Update Repository** (`internal/database/ent_employee_repo.go`):
   Add phone field to Save/Update methods

5. **Restart Application**: EntGo auto-migrates the new column!

### Adding a New Entity

1. **Create Schema** (`internal/ent/schema/project.go`):
   ```bash
   go run entgo.io/ent/cmd/ent new --target internal/ent/schema Project
   ```

2. **Define Fields and Relationships** in the schema file

3. **Regenerate Code**: `go generate ./internal/ent`

4. **Create Repository Interface** in `internal/database/database.go`

5. **Implement Repository** in `internal/database/ent_project_repo.go`

6. **Create Handler** in `internal/project/handler.go`

7. **Add Routes** in `internal/router/router.go`

### Modifying EntGo Schema
- After ANY schema change, ALWAYS run: `go generate ./internal/ent`
- Changes are applied automatically on next app restart
- EntGo compares schema to database and applies diffs

### Adding Validation
- Request validation uses Gin's binding tags in DTOs (`internal/models/models.go`)
- Add validation tags like `binding:"required,email"` to struct fields
- Database constraints defined in EntGo schema (e.g., `.NotEmpty()`, `.Unique()`)

## Project Structure
```
cmd/api/main.go                          - Application entry point, EntGo client initialization
internal/
├── config/config.go                     - Environment configuration loader
├── models/models.go                     - Domain models and request DTOs
├── database/
│   ├── database.go                      - Repository interfaces
│   ├── ent_client.go                    - EntGo client setup & migrations
│   ├── ent_department_repo.go           - Department repository (EntGo)
│   ├── ent_employee_repo.go             - Employee repository (EntGo)
│   └── legacy/                          - Old implementations (for reference)
├── ent/                                 - EntGo generated code (DO NOT EDIT MANUALLY!)
│   ├── schema/
│   │   ├── department.go                - Department schema definition
│   │   └── employee.go                  - Employee schema definition
│   ├── client.go                        - Generated database client
│   ├── department.go, department_*.go   - Generated department code
│   ├── employee.go, employee_*.go       - Generated employee code
│   └── ...                              - Other generated files
├── router/router.go                     - HTTP routing configuration
├── department/handler.go                - Department HTTP handlers
└── employee/handler.go                  - Employee HTTP handlers
```

### Important Directories

- **`internal/ent/schema/`**: Define your database schema here (manually edit these)
- **`internal/ent/`** (everything else): Generated code - DO NOT edit manually!
- **`legacy/`**: Old implementations kept for learning purposes

## Notes

- Vietnamese comments in code explain architecture decisions
- Port configured via environment variable (default: 8080)
- Uses Gin's Logger and Recovery middleware by default
- No authentication/authorization implemented
- No test files exist yet - testing infrastructure needs to be set up
- **IMPORTANT**: Always run `go generate ./internal/ent` after modifying schema files
- Database migrations happen automatically on application startup
- EntGo generates ~2000+ lines of code from schema definitions