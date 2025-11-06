# Legacy Database Implementation Files

These files represent the original database implementation using raw SQL and PostgreSQL connection pools. They are kept for reference and learning purposes.

## Files in This Directory

### `postgres.go`
- Connection pool setup using `pgx/v5`
- Direct database connection management
- Pool configuration (max connections, timeouts, etc.)

### `postgres_repository.go`
- Raw SQL implementation of repository interfaces
- Manual query construction with prepared statements
- Direct `pgx.Rows` scanning and error handling
- Shows CRUD operations using SQL strings

### `memory_repository.go`
- In-memory storage using Go maps
- Thread-safe with sync.RWMutex
- Used for development/testing without a real database
- Demonstrates repository pattern with no database

### `migrate.go`
- Integration with `golang-migrate/migrate`
- Runs SQL migration files from `migrations/` directory
- Manual migration version control

## Learning Value

Compare these files with the current EntGo implementation to understand:

1. **Query Building**: Raw SQL strings vs. type-safe builders
2. **Error Handling**: Manual scanning vs. automatic conversion
3. **Migrations**: SQL files vs. Go schema definitions
4. **Type Safety**: Runtime errors vs. compile-time checks

## Current Implementation

The new EntGo-based implementation is in:
- `../ent_client.go` - Database client setup
- `../ent_department_repo.go` - Department repository
- `../ent_employee_repo.go` - Employee repository
- `../../ent/schema/` - Schema definitions

**Recommendation**: Study both approaches to understand the tradeoffs between control (raw SQL) and productivity (ORM).
