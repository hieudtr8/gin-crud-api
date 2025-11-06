# Legacy Implementation Files

This directory contains old implementation files kept for learning and reference purposes.

## Contents

### `migrations/`
Manual SQL migration files used with `golang-migrate/migrate` library.
- These are now replaced by EntGo's automatic schema migrations
- EntGo generates migrations from Go code instead of writing SQL manually

### Reference
For the current EntGo implementation, see:
- `internal/ent/schema/` - Schema definitions
- `internal/database/ent_*` - EntGo repository implementations

## Why Keep These?

These files demonstrate different approaches to database management in Go:

1. **Manual SQL Migrations** (`migrations/`)
   - Shows how to write raw SQL for schema changes
   - Uses golang-migrate for version control
   - More control but more verbose

2. **Raw SQL Implementation** (`internal/database/legacy/postgres_repository.go`)
   - Direct SQL queries with prepared statements
   - Manual error handling and type conversion
   - Good for understanding SQL fundamentals

3. **In-Memory Storage** (`internal/database/legacy/memory_repository.go`)
   - Simple map-based storage for development/testing
   - Shows how to implement repository interface without a database

## Comparison: Old vs New

| Aspect | Legacy (Raw SQL) | Current (EntGo) |
|--------|------------------|-----------------|
| Schema Changes | Manual SQL files | Go code (auto-generated SQL) |
| Queries | Hand-written SQL strings | Type-safe query builders |
| Migrations | golang-migrate | EntGo automatic migrations |
| Type Safety | Runtime errors | Compile-time errors |
| Code Volume | ~400 lines | ~200 lines |
| Learning Curve | SQL knowledge required | Go structs & methods |

**Recommendation**: Start with the current EntGo implementation for production code, but study the legacy files to understand what's happening under the hood.
