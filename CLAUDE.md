# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a REST API CRUD application built with Go and the Gin web framework for managing Departments and Employees. The project demonstrates clean architecture with dependency injection, repository pattern, and in-memory data storage. Currently, only Create and Get operations are implemented.

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
3. **Repository Layer** (`internal/database/`): Data access interfaces and implementations
4. **Model Layer** (`internal/models/`): Domain models and DTOs

### Dependency Flow
```
main.go → InMemoryStore → Repositories → Handlers → Router
```
- All dependencies are injected via constructors
- Handlers depend on repository interfaces, not concrete implementations
- This allows swapping in-memory storage for a real database without changing handler code

### Key Design Patterns

1. **Repository Pattern**: All data access goes through repository interfaces (`DepartmentRepository`, `EmployeeRepository`)
2. **Dependency Injection**: Constructor-based injection throughout
3. **Thread Safety**: Separate RWMutexes for departments and employees allow concurrent reads
4. **DTO Pattern**: Separate request DTOs (`Create*Request`, `Update*Request`) with validation tags

## Implementation Status

### Implemented Endpoints
- `POST /api/v1/departments` - Create department
- `GET /api/v1/departments/:id` - Get department by ID
- `POST /api/v1/employees` - Create employee (validates department exists)
- `GET /api/v1/employees/:id` - Get employee by ID
- `GET /health` - Health check

### Not Implemented
The following operations are partially scaffolded but not complete:
- List all departments (`FindAll()` returns nil)
- Update department
- Delete department
- List all employees
- Update employee
- Delete employee

Routes for these operations are commented out in `internal/router/router.go`.

## Data Storage

Currently uses in-memory storage (`InMemoryStore`):
- Data is stored in Go maps
- Thread-safe with RWMutex protection
- **Important**: All data is lost when the server restarts
- No persistence mechanism implemented

To add persistence, implement the repository interfaces with a real database.

## Business Rules

1. **Employee Creation**: Must reference an existing department (validated in handler)
2. **IDs**: Generated as UUIDs using google/uuid package
3. **Validation**:
   - Department name is required
   - Employee name and email are required
   - Employee email must be valid format
   - Employee must belong to existing department

## Common Tasks

### Adding a New Endpoint
1. Add handler method to appropriate handler file
2. Uncomment or add route in `internal/router/router.go`
3. Implement corresponding repository method if needed

### Switching to Real Database
1. Create new repository implementations (e.g., `PostgresDepartmentRepo`)
2. Implement all methods from the repository interfaces
3. Update `main.go` to use new implementations instead of `InMemoryStore`

### Adding Validation
- Request validation uses Gin's binding tags in DTOs (`internal/models/models.go`)
- Add validation tags like `binding:"required,email"` to struct fields

## Project Structure
```
cmd/api/main.go                 - Application entry point, dependency wiring
internal/models/models.go       - Domain models and request DTOs
internal/database/database.go   - Repository interfaces and in-memory implementation
internal/router/router.go       - HTTP routing configuration
internal/department/handler.go  - Department HTTP handlers
internal/employee/handler.go    - Employee HTTP handlers
```

## Notes

- Vietnamese comments in code explain architecture decisions
- Port 8080 is hardcoded in main.go
- Uses Gin's Logger and Recovery middleware by default
- No authentication/authorization implemented
- No test files exist yet - testing infrastructure needs to be set up