# 🚀 Gin CRUD API - Department & Employee Management

A clean architecture REST API built with Go, Gin framework, and EntGo ORM for type-safe database operations.

## 📋 Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Key Concepts for Go Beginners](#key-concepts-for-go-beginners)
- [Development Guide](#development-guide)

## 🎯 Overview

This project is a REST API for managing departments and employees with:
- **Full CRUD operations** for both entities
- **EntGo ORM**: Type-safe database operations with automatic migrations
- **Clean architecture** with repository pattern
- **PostgreSQL**: Production-ready database with automatic schema management
- **Foreign key constraints**: Database-enforced relationships
- **Thread-safe** operations with connection pooling
- **Docker support** for easy PostgreSQL setup
- **Auto-generated timestamps**: created_at and updated_at handled automatically

### Entity Relationship
```
Department (1) ──────→ (many) Employee
    ├── id (UUID)           ├── id (UUID)
    ├── name                ├── name
    ├── created_at          ├── email
    └── updated_at          ├── department_id (FK)
                            ├── created_at
                            └── updated_at
```

**Note**: Timestamps are automatically managed by EntGo

## 🏗 Architecture

### How the Code Flows

```
HTTP Request
    ↓
Router (Gin) → Defines URL paths (/api/v1/departments, etc.)
    ↓
Handler → Processes request, validates input
    ↓
Repository Interface → Defines what operations are available
    ↓
EntGo Repository → Type-safe database operations
    ↓
EntGo Client → Generates SQL queries automatically
    ↓
PostgreSQL → Stores data with ACID guarantees
```

### Why This Architecture?

1. **Type Safety**: EntGo generates type-safe code at compile time
2. **Automatic Migrations**: Schema changes applied automatically on startup
3. **No Raw SQL**: Write Go code, EntGo generates optimized SQL
4. **Separation of Concerns**: Each layer has one job
5. **Easy Testing**: Can swap real database with mock
6. **Maintainable**: Clear boundaries between components

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

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ installed
- Docker (for PostgreSQL)
- Make (optional, for shortcuts)

### 1. Clone and Install
```bash
git clone <your-repo>
cd gin-crud-api
go mod download
```

### 2. Start PostgreSQL & Run Application

```bash
# 1. Start PostgreSQL with Docker
docker-compose up -d
# OR
make docker-up

# 2. Start the server (migrations run automatically!)
go run cmd/api/main.go
# OR
make run

# Server will start on http://localhost:8080
# EntGo automatically creates/updates database tables on startup
```

**What happens on startup?**
- ✅ Connects to PostgreSQL database
- ✅ Runs automatic schema migrations (creates tables if they don't exist)
- ✅ Updates existing tables if schema changed
- ✅ Starts HTTP server on port 8080

**No manual migrations needed!** EntGo handles everything automatically.

### 3. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Create a department
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Content-Type: application/json" \
  -d '{"name": "Engineering"}'
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Departments

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| GET | `/departments` | List all departments | - |
| GET | `/departments/:id` | Get specific department | - |
| POST | `/departments` | Create department | `{"name": "Engineering"}` |
| PUT | `/departments/:id` | Update department | `{"name": "New Name"}` |
| DELETE | `/departments/:id` | Delete department (cascade) | - |

#### Employees

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| GET | `/employees` | List all employees | - |
| GET | `/employees/:id` | Get specific employee | - |
| POST | `/employees` | Create employee | See below |
| PUT | `/employees/:id` | Update employee | See below |
| DELETE | `/employees/:id` | Delete employee | - |

**Create/Update Employee Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "department_id": "uuid-here"
}
```

### Example Flow

```bash
# 1. Create a department
curl -X POST localhost:8080/api/v1/departments \
  -H "Content-Type: application/json" \
  -d '{"name": "Engineering"}'

# Response: {"id": "123-456", "name": "Engineering"}

# 2. Create an employee in that department
curl -X POST localhost:8080/api/v1/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice",
    "email": "alice@example.com",
    "department_id": "123-456"
  }'

# 3. List all employees
curl localhost:8080/api/v1/employees
```

## 📁 Project Structure

```
gin-crud-api/
│
├── cmd/api/
│   └── main.go                 # Application entry point
│
├── internal/                   # Private application code
│   ├── config/
│   │   └── config.go          # Environment configuration
│   │
│   ├── database/
│   │   ├── database.go        # Repository interfaces
│   │   ├── ent_client.go      # EntGo client initialization
│   │   ├── ent_department_repo.go  # Department repository (EntGo)
│   │   ├── ent_employee_repo.go    # Employee repository (EntGo)
│   │   └── legacy/            # Old implementations (for reference)
│   │
│   ├── ent/                   # EntGo generated code
│   │   ├── schema/
│   │   │   ├── department.go  # Department schema definition
│   │   │   └── employee.go    # Employee schema definition
│   │   ├── department/        # Generated department code
│   │   ├── employee/          # Generated employee code
│   │   ├── client.go          # Generated database client
│   │   └── ...                # Other generated files
│   │
│   ├── models/
│   │   └── models.go          # DTOs and domain models
│   │
│   ├── department/
│   │   └── handler.go         # HTTP handlers for departments
│   │
│   ├── employee/
│   │   └── handler.go         # HTTP handlers for employees
│   │
│   └── router/
│       └── router.go          # URL routing setup
│
├── legacy/                     # Old implementations (for learning)
│   ├── migrations/            # Manual SQL migrations (replaced by EntGo)
│   └── README.md              # Explanation of legacy code
│
├── docker-compose.yml         # PostgreSQL setup
├── Makefile                   # Development shortcuts
└── .env                       # Configuration (don't commit!)
```

## 🔰 Key Concepts for Go Beginners

### 1. Interfaces (Repository Pattern)
```go
// Interface defines what methods must exist
type DepartmentRepository interface {
    Save(dept *models.Department) error
    FindByID(id string) (*models.Department, error)
    // ... other methods
}

// Different implementations of the same interface
type InMemoryDepartmentRepo struct { }  // Uses maps
type PostgresDepartmentRepo struct { }   // Uses database

// Both implement the same interface, so handlers don't care which one is used!
```

### 2. Pointers (`*` and `&`)
```go
// * in type = "pointer to"
func Save(dept *models.Department)  // Receives pointer, can modify

// & = "get address of"
dept := &models.Department{...}  // Create and get pointer

// Why? Efficiency and sharing
// - Pointers avoid copying large structs
// - Multiple parts can access same data
```

### 3. Error Handling
```go
// Go functions often return (result, error)
dept, err := repository.FindByID(id)
if err != nil {
    // Handle error
    return nil, err
}
// Use result
```

### 4. Dependency Injection
```go
// Handler receives what it needs through constructor
func NewHandler(repo DepartmentRepository) *Handler {
    return &Handler{repo: repo}
}
// This makes testing easy - just pass a mock repository!
```

### 5. EntGo Schema-First Approach
```go
// Define your database schema in Go (internal/ent/schema/department.go)
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

// Use the generated client (no SQL writing needed!)
dept, err := client.Department.
    Query().
    Where(department.NameContains("Engineering")).
    First(ctx)
```

### 6. EntGo Code Generation
```bash
# After modifying schema files, regenerate code:
go generate ./internal/ent

# This creates/updates all database client code automatically
# EntGo generates ~20+ files with type-safe operations
```

## 🛠 Development Guide

### Configuration (.env file)

```bash
# Server settings
SERVER_PORT=8080

# PostgreSQL settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_crud_api
DB_SSLMODE=disable

# Connection pool settings
DB_MAX_CONNS=25
DB_MIN_CONNS=5
```

**Note**: These settings are automatically loaded by the application. EntGo uses these to connect to PostgreSQL and manage the connection pool.

### Common Commands

```bash
# Development
make help           # Show all commands
make docker-up      # Start PostgreSQL
make run            # Run application (migrations automatic!)
make docker-down    # Stop PostgreSQL

# EntGo Commands
go generate ./internal/ent    # Regenerate EntGo code after schema changes

# Testing
make test           # Run tests
./test_api.sh       # Run integration tests

# Build
make build          # Create executable
./build/gin-crud-api  # Run the executable
```

### Working with EntGo Schemas

When you need to modify the database schema:

1. **Edit Schema Files** (`internal/ent/schema/*.go`)
   ```go
   // Example: Add a field to Department
   field.String("description").Optional()
   ```

2. **Regenerate Code**
   ```bash
   go generate ./internal/ent
   ```

3. **Restart Application**
   ```bash
   go run cmd/api/main.go
   # EntGo automatically applies schema changes!
   ```

**That's it!** No manual SQL migrations needed.

### How Data Flows Through the Code

1. **Request arrives** at router (`/api/v1/departments`)
2. **Router** calls appropriate handler (`deptHandler.Create`)
3. **Handler** validates input and calls repository
4. **Repository** uses EntGo client for type-safe operations
5. **EntGo** generates and executes SQL queries
6. **PostgreSQL** stores/retrieves data
7. **Response** sent back to client

### Adding New Features with EntGo

**Example: Adding a `phone` field to Employee**

1. **Update EntGo Schema** (`internal/ent/schema/employee.go`):
   ```go
   func (Employee) Fields() []ent.Field {
       return []ent.Field{
           // ... existing fields
           field.String("phone").
               Optional().
               Comment("Employee phone number"),
       }
   }
   ```

2. **Regenerate EntGo Code**:
   ```bash
   go generate ./internal/ent
   ```

3. **Update Domain Model** (`internal/models/models.go`):
   ```go
   type Employee struct {
       // ... existing fields
       Phone string `json:"phone,omitempty"`
   }
   ```

4. **Update Repository** (`internal/database/ent_employee_repo.go`):
   ```go
   // Add phone to Save and Update methods
   SetPhone(emp.Phone)
   ```

5. **Restart Application**:
   ```bash
   go run cmd/api/main.go
   # EntGo automatically adds the phone column!
   ```

**No SQL writing needed!** EntGo handles all database changes automatically.

## 🧪 Testing

```bash
# Run the test script
./test_api.sh

# Manual testing with curl
curl -X POST localhost:8080/api/v1/departments \
  -H "Content-Type: application/json" \
  -d '{"name": "Sales"}'
```

## 🚨 Important Notes

1. **UUID IDs**: All IDs are UUIDs (universally unique identifiers)
2. **Auto-Migrations**: EntGo automatically creates/updates database tables on startup
3. **Foreign Key Constraints**: Database enforces relationships (cannot add employee with invalid department_id)
4. **Timestamps**: `created_at` and `updated_at` managed automatically by EntGo
5. **Type Safety**: Compile-time checking for all database operations
6. **Validation**: Email must be valid format, all required fields enforced by database
7. **Empty Lists**: APIs return `[]` not `null` when no data exists
8. **Code Generation**: Always run `go generate ./internal/ent` after schema changes

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Find and kill process on port 8080
lsof -i :8080
kill -9 <PID>
```

### PostgreSQL Connection Failed
```bash
# Check if PostgreSQL is running
docker ps

# Restart PostgreSQL
make docker-down
make docker-up
```

### EntGo Schema Generation Failed
```bash
# Re-generate EntGo code
go generate ./internal/ent

# If still failing, check schema files for syntax errors
# Look in internal/ent/schema/*.go
```

### Database Schema Out of Sync
```bash
# EntGo migrations run on startup, so just restart the app
go run cmd/api/main.go

# To manually inspect database schema:
docker exec -it gin_crud_postgres psql -U postgres -d gin_crud_api -c "\d departments"
docker exec -it gin_crud_postgres psql -U postgres -d gin_crud_api -c "\d employees"
```

## 📝 License

MIT

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

Built with ❤️ using Go and Gin Framework