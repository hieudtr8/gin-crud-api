# 🚀 Gin CRUD API - Department & Employee Management

A clean architecture REST API built with Go and Gin framework, featuring both in-memory and PostgreSQL storage options.

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
- **Two storage options**: In-memory (for development) or PostgreSQL (for production)
- **Clean architecture** with repository pattern
- **Cascade delete**: Deleting a department automatically removes its employees
- **Thread-safe** operations
- **Docker support** for easy PostgreSQL setup

### Entity Relationship
```
Department (1) ──────→ (many) Employee
    ├── id (UUID)           ├── id (UUID)
    └── name                ├── name
                            ├── email
                            └── department_id (FK)
```

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
Repository Implementation → Actual data operations
    ├── InMemory: Uses Go maps + mutex for thread safety
    └── PostgreSQL: Uses database with connection pooling
```

### Why This Architecture?

1. **Separation of Concerns**: Each layer has one job
2. **Easy Testing**: Can swap real database with mock
3. **Flexibility**: Switch storage without changing business logic
4. **Maintainable**: Clear boundaries between components

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

### 2. Choose Your Storage Mode

#### Option A: In-Memory (Simple, No Setup)
```bash
# Using make
make run-memory

# Or directly
STORAGE_TYPE=memory go run cmd/api/main.go
```

#### Option B: PostgreSQL (Persistent, Production-like)
```bash
# 1. Start PostgreSQL with Docker
make docker-up

# 2. Run migrations (create tables)
make migrate-up

# 3. Start server
make run-postgres

# Or without make:
docker-compose up -d
STORAGE_TYPE=postgres go run cmd/api/main.go
```

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
│   │   ├── database.go        # Repository interfaces + in-memory impl
│   │   ├── postgres.go        # PostgreSQL connection management
│   │   ├── postgres_repository.go # PostgreSQL implementations
│   │   └── migrate.go         # Database migration runner
│   │
│   ├── models/
│   │   └── models.go          # Data structures (Department, Employee)
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
├── migrations/                 # SQL migration files
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

## 🛠 Development Guide

### Configuration (.env file)

```bash
# Storage type: "memory" or "postgres"
STORAGE_TYPE=postgres

# PostgreSQL settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_crud_api
```

### Common Commands

```bash
# Development
make help           # Show all commands
make docker-up      # Start PostgreSQL
make migrate-up     # Create database tables
make run-postgres   # Run with PostgreSQL
make run-memory     # Run with in-memory storage
make docker-down    # Stop PostgreSQL

# Testing
make test           # Run tests
./test_api.sh       # Run integration tests

# Build
make build          # Create executable
./gin-crud-api      # Run the executable
```

### How Data Flows Through the Code

1. **Request arrives** at router (`/api/v1/departments`)
2. **Router** calls appropriate handler (`deptHandler.Create`)
3. **Handler** validates input and calls repository
4. **Repository** performs actual data operation
5. **Response** sent back to client

### Storage Differences

| Feature | In-Memory | PostgreSQL |
|---------|-----------|------------|
| **Persistence** | Lost on restart | Saved permanently |
| **Setup** | None | Needs Docker/PostgreSQL |
| **Performance** | Fastest | Fast with indexes |
| **Use Case** | Development/Testing | Production |

### Adding New Features

To add a new field (e.g., `phone` to Employee):

1. **Update Model** (`internal/models/models.go`):
```go
type Employee struct {
    // ... existing fields
    Phone string `json:"phone"`
}
```

2. **Update Database** (create migration):
```sql
ALTER TABLE employees ADD COLUMN phone VARCHAR(50);
```

3. **Update Repository** implementations to handle new field

4. **Update Handler** to accept phone in requests

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
2. **Thread Safety**: In-memory storage uses mutexes for concurrent access
3. **Cascade Delete**: Deleting a department removes all its employees
4. **Validation**: Email must be valid format, all fields are required
5. **Empty Lists**: APIs return `[]` not `null` when no data exists

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

### Migration Failed
```bash
# Check database exists
docker exec -it gin_crud_postgres psql -U postgres -c "\l"

# Reset migrations
make migrate-down
make migrate-up
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