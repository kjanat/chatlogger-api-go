# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ChatLogger API is a multi-tenant backend API for logging and managing chat sessions. It's built with Go using clean architecture principles with separation of concerns through layered design (handler → service → repository).

## Key Architecture Decisions

1. **Clean Architecture Pattern**: Strict separation between handlers, services, and repositories
2. **Dependency Injection**: Constructor-based DI for testability
3. **Strategy Pattern**: Used for exporters (JSON/CSV) and other pluggable components
4. **Multi-tenancy**: Organizations are isolated with separate data access
5. **Dual Authentication**: API keys for chat plugins, JWT for dashboard users
6. **Async Job Processing**: Redis + Asynq for background tasks (exports)

## Development Commands

### Building
```bash
# Build server
./scripts/build.sh
# Or on Windows: ./scripts/build.ps1

# Build with Docker
docker-compose build
```

### Running Locally
```bash
# Start all services (PostgreSQL, Redis, API server, worker)
docker-compose up -d

# Run just the API server (requires PostgreSQL and Redis running)
go run cmd/server/main.go

# Run the worker for async exports
go run cmd/worker/main.go
```

### Documentation Generation
```bash
# Generate Swagger/OpenAPI documentation
./scripts/docs_generate.sh
# Or on Windows: ./scripts/docs_generate.ps1

# Generate OpenAPI v3.1 documentation
./scripts/docs_generate.sh --v31
```

### Testing
```bash
# Run all tests
go test -v -race ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt ./...

# Run linting (uses golangci-lint)
golangci-lint run
```

### Database Migrations
```bash
# Apply migrations (when using local PostgreSQL)
psql -U postgres -d chatlogger -f migrations/001_initial_schema.sql
psql -U postgres -d chatlogger -f migrations/002_ensure_defaults.sql
psql -U postgres -d chatlogger -f migrations/003_add_exports_table.sql
```

## Code Organization

- `/cmd/server`: Main API server entry point
- `/cmd/worker`: Background job worker for async exports
- `/internal/handler`: HTTP request handlers - implement REST endpoints
- `/internal/service`: Business logic layer - orchestrates operations
- `/internal/repository`: Data access layer - database operations only
- `/internal/middleware`: Auth, RBAC, versioning middleware
- `/internal/strategy`: Strategy pattern implementations (exporters)
- `/internal/jobs`: Async job definitions and processors

## Important Patterns

### Adding New Endpoints
1. Define handler method in appropriate handler file (`/internal/handler/`)
2. Implement business logic in service layer (`/internal/service/`)
3. Add repository methods if needed (`/internal/repository/`)
4. Register route in `/internal/api/router.go`
5. Add Swagger annotations to handler method

### Authentication Flow
- Public API: Check `x-organization-api-key` header via `APIKeyAuth` middleware
- Dashboard: JWT cookie authentication via `JWTAuth` middleware
- Superadmin endpoints: Additional `RequireSuperadmin` middleware

### Export System
- Synchronous exports: Direct response, good for small datasets
- Asynchronous exports: Queued via Redis/Asynq, status tracking, better for large datasets
- Strategy pattern allows easy addition of new export formats

### Error Handling
- Use `domain.NewError()` for domain errors with proper HTTP status codes
- Repository errors should be wrapped with context
- Services return domain errors that handlers can directly use

## Environment Configuration

Required environment variables:
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret for JWT signing
- `REDIS_ADDR`: Redis address (default: localhost:6379)
- `EXPORT_DIR`: Directory for export files (default: ./exports)
- `PORT`: Server port (default: 8080)

## Testing Approach

When adding new features:
1. Write repository tests with actual database (integration tests)
2. Write service tests with mocked repositories (unit tests)
3. Write handler tests with mocked services (unit tests)
4. Consider end-to-end tests for critical workflows

## Version Management

Version is managed in `/internal/version/version.go`. The CI/CD pipeline handles:
- Automatic version tagging on master branch
- Manual version bumping via GitHub Actions workflow
- Version injection during build process
