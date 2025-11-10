# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **GraphQL API** CRUD application built with Go, gqlgen, and EntGo ORM for managing Departments and Employees. The project demonstrates clean architecture with dependency injection, repository pattern, schema-first GraphQL development, type-safe database operations with EntGo, and automatic schema migrations. Full CRUD operations (queries and mutations) are implemented using PostgreSQL as the data store.

**Note:** The original REST API implementation has been moved to `internal/legacy/rest/` for reference and learning purposes.

## Commands

### Build and Run GraphQL Server
```bash
# Run the GraphQL server (starts on port 8081)
go run ./cmd/graphql

# Build the GraphQL server
go build -o ./build/graphql-server ./cmd/graphql

# Run the built server
./build/graphql-server

# Access GraphQL Playground
# Open http://localhost:8081 in your browser
```

### Run Legacy REST API (Optional)
```bash
# Run the legacy REST server (starts on port 8080)
go run ./cmd/legacy/rest_main.go
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

# Regenerate EntGo code (after database schema changes)
go generate ./internal/ent

# Regenerate GraphQL code (after GraphQL schema changes)
go run github.com/99designs/gqlgen generate
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

1. **GraphQL Layer** (`internal/graph/`): GraphQL schema, resolvers, and generated code
2. **Resolver Layer** (`internal/graph/schema.resolvers.go`): GraphQL query/mutation handling and validation
3. **Repository Layer** (`internal/database/`): Repository interfaces and EntGo implementations
4. **EntGo Layer** (`internal/ent/`): Generated type-safe database client
5. **Database Layer**: PostgreSQL with automatic schema management

### Dependency Flow
```
cmd/graphql/main.go → EntGo Client → EntGo Repositories → GraphQL Resolvers → gqlgen Server
```
- All dependencies are injected via constructors
- Resolvers depend on repository interfaces, not concrete implementations
- EntGo client provides type-safe database operations
- gqlgen generates type-safe GraphQL execution code
- Automatic database migrations on application startup

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
   go generate ./internal/ent
   go run github.com/99designs/gqlgen generate
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
   go generate ./internal/ent
   go run github.com/99designs/gqlgen generate
   ```

5. **Create Repository Interface** in `internal/database/database.go`

6. **Implement Repository** in `internal/database/ent_project_repo.go`

7. **Update Resolver Struct** (`internal/graph/resolver.go`) to inject new repository

8. **Implement Resolvers** in `internal/graph/schema.resolvers.go`

### Modifying GraphQL Schema
- After ANY schema change, ALWAYS run: `go run github.com/99designs/gqlgen generate`
- gqlgen updates generated code and resolver stubs
- Your resolver implementations are preserved during regeneration

### Modifying EntGo Schema
- After ANY EntGo schema change, ALWAYS run: `go generate ./internal/ent`
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
├── config/config.go                     - Environment configuration loader
├── models/models.go                     - Domain models
├── graph/                               - GraphQL layer
│   ├── schema.graphql                   - GraphQL schema definition (EDIT THIS!)
│   ├── resolver.go                      - Resolver struct with dependency injection
│   ├── schema.resolvers.go              - Resolver implementations (EDIT THIS!)
│   ├── generated.go                     - Generated GraphQL execution code (DO NOT EDIT!)
│   └── model/
│       ├── models_gen.go                - Generated GraphQL models (DO NOT EDIT!)
│       └── model.go                     - Custom GraphQL models (if needed)
├── database/
│   ├── database.go                      - Repository interfaces
│   ├── ent_client.go                    - EntGo client setup & migrations
│   ├── ent_department_repo.go           - Department repository (EntGo)
│   └── ent_employee_repo.go             - Employee repository (EntGo)
├── ent/                                 - EntGo generated code (DO NOT EDIT MANUALLY!)
│   ├── schema/
│   │   ├── department.go                - Department schema definition (EDIT THIS!)
│   │   └── employee.go                  - Employee schema definition (EDIT THIS!)
│   ├── client.go                        - Generated database client
│   ├── department.go, department_*.go   - Generated department code
│   ├── employee.go, employee_*.go       - Generated employee code
│   └── ...                              - Other generated files
└── legacy/rest/                         - Legacy REST API (for reference)
    ├── department/handler.go            - REST department handlers
    ├── employee/handler.go              - REST employee handlers
    └── router/router.go                 - REST routing configuration
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

## Notes

- Vietnamese comments in code explain architecture decisions
- GraphQL server port configured via `GRAPHQL_PORT` environment variable (default: 8081)
- Legacy REST server port configured via `SERVER_PORT` environment variable (default: 8080)
- No authentication/authorization implemented yet
- No test files exist yet - testing infrastructure needs to be set up
- **IMPORTANT Code Generation**:
  - After modifying GraphQL schema: `go run github.com/99designs/gqlgen generate`
  - After modifying EntGo schema: `go generate ./internal/ent`
- Database migrations happen automatically on application startup
- Both gqlgen and EntGo use code generation (~3000+ lines generated total)
- GraphQL Playground provides interactive API documentation and testing

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