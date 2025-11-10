# Legacy REST API Documentation

This directory contains the original REST API implementation, preserved for reference, learning, and quick reversion if needed.

## 📁 Contents

- **`rest/`**: The original REST API implementation
  - `department/handler.go`: Department REST handlers with Gin
  - `employee/handler.go`: Employee REST handlers with Gin
  - `router/router.go`: Gin router configuration

## 🚀 Running the Legacy REST API

You can still run the original REST API:

```bash
go run cmd/legacy/rest_main.go
```

The REST server will start on **port 8080** (configurable via `SERVER_PORT` environment variable).

## 🎯 Purpose

This legacy code is kept for:
1. **Learning**: Compare REST vs GraphQL implementations side-by-side
2. **Reference**: Understand the migration process from REST to GraphQL
3. **Quick Revert**: Ability to quickly switch back if needed
4. **Educational**: See two different API architectures using the same business logic

## 📚 REST API Documentation

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

### Example REST Flow

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

# 4. Get specific employee
curl localhost:8080/api/v1/employees/employee-id

# 5. Update employee
curl -X PUT localhost:8080/api/v1/employees/employee-id \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "email": "alice.smith@example.com",
    "department_id": "123-456"
  }'

# 6. Delete employee
curl -X DELETE localhost:8080/api/v1/employees/employee-id

# 7. Health check
curl localhost:8080/health
```

## 🔄 REST vs GraphQL Comparison

### REST Approach
```bash
# Multiple requests needed to get department with employees
curl localhost:8080/api/v1/departments/dept-id
curl localhost:8080/api/v1/employees?department_id=dept-id
```

### GraphQL Approach
```graphql
# Single request gets everything
query {
  department(id: "dept-id") {
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

### Key Differences

| Aspect | REST | GraphQL |
|--------|------|---------|
| **Endpoints** | Multiple (`/departments`, `/employees`) | Single (`/query`) |
| **Data Fetching** | Fixed response structure | Client specifies fields |
| **Over-fetching** | Common (get more data than needed) | Eliminated |
| **Under-fetching** | Requires multiple requests | Single query gets related data |
| **Documentation** | Separate (Swagger, etc.) | Built-in (schema introspection) |
| **Versioning** | URL versioning (`/v1/`, `/v2/`) | Schema evolution |
| **Testing** | curl, Postman | Interactive Playground |

## 🏗️ REST Implementation Architecture

### How REST Requests Flow

```
HTTP Request
    ↓
Router (Gin) → Matches URL path (/api/v1/departments)
    ↓
Handler → Validates input, calls repository
    ↓
Repository Interface → Defines CRUD operations
    ↓
EntGo Repository → Type-safe database operations
    ↓
PostgreSQL → Stores/retrieves data
    ↓
JSON Response → Fixed structure
```

### REST Handler Example

```go
// internal/legacy/rest/department/handler.go
func (h *Handler) Create(c *gin.Context) {
    var req models.CreateDepartmentRequest

    // Bind and validate JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Create department
    dept := &models.Department{
        ID:   uuid.New().String(),
        Name: req.Name,
    }

    // Save via repository
    if err := h.repo.Save(dept); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // Return fixed JSON structure
    c.JSON(201, dept)
}
```

### REST Router Setup

```go
// internal/legacy/rest/router/router.go
func Setup(deptHandler, empHandler) *gin.Engine {
    router := gin.Default()

    v1 := router.Group("/api/v1")
    {
        // Department routes
        v1.POST("/departments", deptHandler.Create)
        v1.GET("/departments", deptHandler.GetAll)
        v1.GET("/departments/:id", deptHandler.GetByID)
        v1.PUT("/departments/:id", deptHandler.Update)
        v1.DELETE("/departments/:id", deptHandler.Delete)

        // Employee routes
        v1.POST("/employees", empHandler.Create)
        v1.GET("/employees", empHandler.GetAll)
        // ... more routes
    }

    return router
}
```

## 🆚 Why We Migrated to GraphQL

### Problems with REST

1. **Over-fetching**: Getting entire employee object when you only need the name
2. **Under-fetching**: Need multiple requests to get department with employees
3. **Multiple Endpoints**: Have to remember `/departments`, `/employees`, etc.
4. **Fixed Responses**: Can't customize what fields are returned
5. **Versioning Complexity**: Need `/v1/`, `/v2/` when API changes

### GraphQL Solutions

1. **Precise Fetching**: Request exactly the fields you need
2. **Nested Queries**: Get related data in one request
3. **Single Endpoint**: Just `/query` for everything
4. **Flexible Responses**: Client controls the shape
5. **Schema Evolution**: Add fields without breaking existing clients

## 🔧 Configuration for REST API

If running the legacy REST server, use these environment variables:

```bash
# Server settings
SERVER_PORT=8080

# PostgreSQL settings (shared with GraphQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_crud_api
DB_SSLMODE=disable
```

## 🧪 Testing REST API

```bash
# Run the test script (if it exists)
./test_api.sh

# Manual testing with curl
curl -X POST localhost:8080/api/v1/departments \
  -H "Content-Type: application/json" \
  -d '{"name": "Sales"}'

curl localhost:8080/api/v1/departments
```

## 📖 Learning Resources

To understand the migration better, compare:

**REST Handler**
- File: `internal/legacy/rest/department/handler.go`
- Uses: Gin context, JSON binding, HTTP status codes
- Returns: Fixed JSON structure

**GraphQL Resolver**
- File: `internal/graph/schema.resolvers.go`
- Uses: Context, generated input types, GraphQL errors
- Returns: Flexible structure based on query

**Shared Components**
- Repository interfaces: `internal/database/database.go`
- Repository implementations: `internal/database/ent_*_repo.go`
- Domain models: `internal/models/models.go`
- Database layer: `internal/ent/`

## 🚨 Important Notes

- Both REST and GraphQL use the **same database** and **same repositories**
- Only the API layer changed (handlers → resolvers)
- Business logic and validation are identical
- EntGo ORM is shared between both implementations
- Switching between REST and GraphQL doesn't require database changes

## ⚠️ Troubleshooting REST API

### Port Already in Use
```bash
# Find and kill process on port 8080
lsof -i :8080
kill -9 <PID>
```

### Gin Router Not Working
```bash
# Check if you're in the right directory
pwd  # Should be in gin-crud-api root

# Check if dependencies are installed
go mod download
```

---

## 🔄 Returning to GraphQL

The primary API is now GraphQL. To use it:

```bash
# Run GraphQL server (port 8081)
go run ./cmd/graphql

# Open Playground
open http://localhost:8081
```

See the main [README.md](../../README.md) for GraphQL documentation.
