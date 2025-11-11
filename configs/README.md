# Configuration Files

This directory contains environment-specific configuration files for the GraphQL API.

## Environment Files

- **`dev.yaml`** - Development environment (default)
  - Used when running locally with `make api` or `go run cmd/graphql/main.go`
  - Optimized for local development (debug logging, localhost database)

- **`prod.yaml`** - Production environment
  - Used in Docker containers and production deployments
  - Set `APP_ENV=prod` to use this configuration
  - Optimized for production (info logging, SSL required, higher connection pools)

- **`test.yaml`** - Test environment
  - Used when running tests with `make test`
  - Set `APP_ENV=test` to use this configuration
  - Optimized for testing (reduced logging, lower resource limits)

## How to Use

### Selecting an Environment

The application selects the configuration file based on the `APP_ENV` environment variable:

```bash
# Development (default)
go run cmd/graphql/main.go
# or
APP_ENV=dev go run cmd/graphql/main.go

# Production
APP_ENV=prod go run cmd/graphql/main.go

# Test
APP_ENV=test go run cmd/graphql/main.go
```

### Environment Variable Overrides

All configuration values can be overridden using environment variables with the `GINAPI_` prefix:

```bash
# Override database host
GINAPI_DATABASE_HOST=my-postgres-server

# Override multiple values
GINAPI_DATABASE_HOST=localhost \
GINAPI_DATABASE_PASSWORD=secret \
GINAPI_LOGGING_LEVEL=debug \
go run cmd/graphql/main.go
```

### Naming Convention

Environment variable names follow this pattern:
```
GINAPI_<SECTION>_<KEY>
```

Examples:
- `GINAPI_SERVER_GRAPHQL_PORT` → `server.graphql_port` in YAML
- `GINAPI_DATABASE_HOST` → `database.host` in YAML
- `GINAPI_LOGGING_LEVEL` → `logging.level` in YAML

### Configuration Priority

Values are loaded in this priority order (highest to lowest):

1. **Environment Variables** (highest priority)
   - Example: `GINAPI_DATABASE_PASSWORD=secret`
2. **YAML Config File**
   - Example: `configs/prod.yaml`
3. **Defaults** (built into the application code)

### Docker Deployment

In `docker-compose.yml`, set the environment:

```yaml
services:
  graphql-api:
    environment:
      APP_ENV: prod
      GINAPI_DATABASE_HOST: postgres
      GINAPI_DATABASE_PASSWORD: ${DB_PASSWORD}
```

## Local Customization

Create a `local.yaml` file for personal overrides (this file is gitignored):

```yaml
# configs/local.yaml (not committed to git)
database:
  password: my-local-password

logging:
  level: debug
```

Use it with:
```bash
APP_ENV=local go run cmd/graphql/main.go
```

## Available Configuration Values

### Server Configuration
- `server.graphql_port` - GraphQL API port (default: 8081)
- `server.rest_port` - Legacy REST API port (default: 8080)

### Database Configuration
- `database.host` - PostgreSQL host
- `database.port` - PostgreSQL port
- `database.user` - Database username
- `database.password` - Database password ⚠️ **Always override in production**
- `database.dbname` - Database name
- `database.sslmode` - SSL mode (disable, require, verify-full)
- `database.max_conns` - Maximum connection pool size
- `database.min_conns` - Minimum idle connections

### Logging Configuration
- `logging.level` - Log level (debug, info, warn, error)
- `logging.pretty` - Pretty console output (true/false)
