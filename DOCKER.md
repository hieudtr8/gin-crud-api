# Docker Deployment Guide

This project provides separate Docker Compose configurations for local development and production deployment.

## 📦 Docker Compose Files

### `docker-compose.local.yml`
**Purpose**: Local development and testing with Docker

**Features**:
- Uses `configs/dev.yaml` (development environment)
- Debug logging (`GINAPI_LOGGING_LEVEL=debug`)
- Pretty console output for easy reading
- PostgreSQL port **5432 exposed** for DB clients (pgAdmin, DBeaver, TablePlus)
- Faster iteration with live rebuilding

**Container Names**:
- `gin_crud_postgres_local`
- `gin_crud_graphql_local`

### `docker-compose.prod.yml`
**Purpose**: Production deployment

**Features**:
- Uses `configs/prod.yaml` (production environment)
- Info-level logging (`GINAPI_LOGGING_LEVEL=info`)
- JSON log output (for log aggregation systems)
- PostgreSQL port **NOT exposed** (internal network only for security)
- SSL required for database connections (`GINAPI_DATABASE_SSLMODE=require`)
- Resource limits configured (CPU and memory)
- Higher connection pool sizes for production load

**Container Names**:
- `gin_crud_postgres_prod`
- `gin_crud_graphql_prod`

## 🚀 Quick Start

### Local Development

```bash
# Start everything (builds fresh image)
make docker-run

# Or use docker-compose directly
docker-compose -f docker-compose.local.yml up --build -d

# Access:
# - GraphQL Playground: http://localhost:8081
# - GraphQL API: http://localhost:8081/query
# - PostgreSQL: localhost:5432
```

### Production Deployment

```bash
# IMPORTANT: Set secure password first!
export DB_PASSWORD=your_secure_password_here

# Start production stack
make docker-run-prod

# Or use docker-compose directly
DB_PASSWORD=secure123 docker-compose -f docker-compose.prod.yml up --build -d

# Access:
# - GraphQL API: http://localhost:8081/query
# - PostgreSQL: NOT exposed (internal network only)
```

**For Cloud Databases (AWS RDS, Google Cloud SQL, Azure Database):**
```bash
# Enable SSL for cloud database connections
DB_PASSWORD=secure123 \
GINAPI_DATABASE_SSLMODE=require \
GINAPI_DATABASE_HOST=your-rds-endpoint.amazonaws.com \
docker-compose -f docker-compose.prod.yml up -d
```

## 🛠️ Common Commands

### Local Development

```bash
# Start services
make docker-run

# View logs
make docker-logs-api       # GraphQL API logs
make docker-logs           # PostgreSQL logs
make docker-logs-all       # All logs

# Restart
make docker-restart

# Stop and remove
make docker-down

# Rebuild everything
make docker-rebuild
```

### Production

```bash
# Start services
make docker-run-prod

# View logs
make docker-logs-api-prod

# Restart
make docker-restart-prod

# Rebuild everything
make docker-rebuild-prod
```

### General

```bash
# View running containers
make docker-ps
docker ps -a | grep gin_crud

# Stop all (both local and prod)
make docker-down
```

## 🔧 Configuration

### Environment Variables

Both deployments use the `GINAPI_` prefix for configuration:

**Database**:
- `GINAPI_DATABASE_HOST` - Database host (auto-set to `postgres` in Docker)
- `GINAPI_DATABASE_PORT` - Database port (default: 5432)
- `GINAPI_DATABASE_USER` - Database username
- `GINAPI_DATABASE_PASSWORD` - ⚠️ **Override in production!**
- `GINAPI_DATABASE_DBNAME` - Database name
- `GINAPI_DATABASE_SSLMODE` - SSL mode (local: disable, prod: require)
- `GINAPI_DATABASE_MAX_CONNS` - Max connections (local: 25, prod: 50)
- `GINAPI_DATABASE_MIN_CONNS` - Min connections (local: 5, prod: 10)

**Logging**:
- `GINAPI_LOGGING_LEVEL` - Log level (local: debug, prod: info)
- `GINAPI_LOGGING_PRETTY` - Pretty output (local: true, prod: false)

**SSL Configuration**:
- Docker PostgreSQL (default): `GINAPI_DATABASE_SSLMODE=disable`
- Cloud Databases (AWS RDS, Cloud SQL): `GINAPI_DATABASE_SSLMODE=require` or `verify-full`
- Standard postgres:16-alpine image does NOT have SSL configured
- For cloud deployments, always use SSL!

### Override Configuration

**Local Development - Custom Settings**:
```bash
# Override specific values
GINAPI_DATABASE_PASSWORD=mypass \
GINAPI_LOGGING_LEVEL=debug \
docker-compose -f docker-compose.local.yml up -d
```

**Production - Secure Deployment**:
```bash
# Use .env file (recommended for production)
cat > .env.prod <<EOF
DB_PASSWORD=super_secure_password
DB_MAX_CONNS=100
LOG_LEVEL=warn
EOF

# Load and deploy
set -a; source .env.prod; set +a
docker-compose -f docker-compose.prod.yml up -d
```

## 📊 Comparison

| Feature | Local (`docker-compose.local.yml`) | Production (`docker-compose.prod.yml`) |
|---------|-----------------------------------|----------------------------------------|
| **Config File** | `configs/dev.yaml` | `configs/prod.yaml` |
| **Environment** | `APP_ENV=dev` | `APP_ENV=prod` |
| **Log Level** | `debug` | `info` |
| **Log Format** | Pretty console | JSON |
| **PostgreSQL Port** | 5432 exposed | NOT exposed |
| **SSL Mode** | `disable` | `disable` (Docker); `require` (Cloud DBs) |
| **Max Connections** | 25 | 50 |
| **Resource Limits** | None | CPU/Memory limits |
| **Restart Policy** | `unless-stopped` | `always` |
| **Volume Name** | `postgres_data_local` | `postgres_data_prod` |
| **Container Suffix** | `_local` | `_prod` |

## 🔒 Security Notes

### Production Checklist

- [ ] ✅ Set `DB_PASSWORD` to a strong, unique password
- [ ] ✅ Use `GINAPI_DATABASE_SSLMODE=require` or higher
- [ ] ✅ Do NOT expose PostgreSQL port 5432 to the internet
- [ ] ✅ Use environment variables or secrets management for sensitive data
- [ ] ✅ Review resource limits based on your server capacity
- [ ] ✅ Set up log aggregation for JSON logs
- [ ] ✅ Configure firewall rules to restrict access
- [ ] ✅ Enable Docker health checks and monitoring
- [ ] ✅ Regular backups of `postgres_data_prod` volume

### Development Best Practices

- ✅ Use `docker-compose.local.yml` for local testing
- ✅ PostgreSQL port exposed for easy database access with GUI tools
- ✅ Debug logging for troubleshooting
- ✅ Separate volumes prevent data conflicts with production

## 🐛 Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose -f docker-compose.local.yml logs graphql-api

# Check if port is already in use
lsof -i :8081
lsof -i :5432

# Rebuild from scratch
make docker-down
docker system prune -a  # Clean up old images
make docker-run
```

### Database connection errors

```bash
# Check if PostgreSQL is healthy
docker-compose -f docker-compose.local.yml ps

# Check PostgreSQL logs
make docker-logs

# Verify network
docker network ls | grep gin_crud
```

### Configuration not loading

```bash
# Verify config files are copied to container
docker-compose -f docker-compose.local.yml exec graphql-api ls -la /app/configs

# Check environment variables in container
docker-compose -f docker-compose.local.yml exec graphql-api env | grep GINAPI
```

## 📁 Volume Management

### Local Development

```bash
# View volumes
docker volume ls | grep local

# Backup local database
docker run --rm -v gin-crud-api_postgres_data_local:/data -v $(pwd):/backup ubuntu tar czf /backup/backup-local.tar.gz /data

# Restore local database
docker run --rm -v gin-crud-api_postgres_data_local:/data -v $(pwd):/backup ubuntu tar xzf /backup/backup-local.tar.gz -C /
```

### Production

```bash
# Backup production database
docker run --rm -v gin-crud-api_postgres_data_prod:/data -v $(pwd):/backup ubuntu tar czf /backup/backup-prod.tar.gz /data

# Clean old volumes (⚠️ DESTROYS DATA!)
docker volume rm gin-crud-api_postgres_data_local
docker volume rm gin-crud-api_postgres_data_prod
```

## 🔗 Related Documentation

- [Configuration System](configs/README.md) - Detailed configuration reference
- [CLAUDE.md](CLAUDE.md) - Project architecture and development guide
- [README.md](README.md) - Main project documentation
