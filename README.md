# ChatLogger API (Go)
<!-- ![ChatLogger API](https://raw.githubusercontent.com/kjanat/chatlogger-api-go/master/assets/logo.png) -->

[![Go version](https://img.shields.io/github/go-mod/go-version/kjanat/chatlogger-api-go?logo=Go&logoColor=white)](go.mod)
[![Go Doc](https://godoc.org/github.com/kjanat/chatlogger-api-go?status.svg)][Package documentation]
[![Go Report Card](https://goreportcard.com/badge/github.com/kjanat/chatlogger-api-go)][Go report] <!-- WTF Can't I Capitalize GO?... [![Go Report](https://img.shields.io/badge/Go%20report-A+-brightgreen.svg)][Go report] -->
[![Static Badge](https://img.shields.io/badge/Bump.sh-API--Docs-blue?logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGZpbGw9Im5vbmUiIHZpZXdCb3g9IjAgMCAyNTYgMjU2Ij48ZyBjbGlwLXBhdGg9InVybCgjcHJlZml4X19hKSI%2BPHBhdGggZmlsbD0iI2ZmZiIgZD0iTTI1NC43MzQgMTAwLjQ2NWMtLjE2Ny0zNi4yMS0yNi4xNC01OC4zMS02My4zODItNTguMzEtMzcuMjQzIDAtNjMuMTEzIDIyLjQzNi02My4xNTIgNTguMDkxIDAgMS42MzktLjE0MiAxNC41MDktLjMyMyAzMS4wMzMtLjAxMyAxLjI3OC0xLjgyIDEuNTEtMi4xNTUuMjcxLTYuOTA3LTI2LjIwNS0zMC4xNDMtNDEuNjMyLTYxLjI2Ny00MS42MzItMzcuMzg1IDAtNjMuMTY0IDIyLjQ3NS02My4xNjQgNTguMTk1IDAgLjc0OC0xLjI5MSA2NS43MzItMS4yOTEgNjUuNzMyaDI1NnMtMS4wMDctMTA4LjgzNi0xLjI1Mi0xMTMuMzh6Ii8%2BPC9nPjxkZWZzPjxjbGlwUGF0aCBpZD0icHJlZml4X19hIj48cGF0aCBmaWxsPSIjZmZmIiBkPSJNMCAwaDI1NnYyNTZIMHoiLz48L2NsaXBQYXRoPjwvZGVmcz48L3N2Zz4%3D&link=https%3A%2F%2Fchatlogger-api-docs.kjanat.com%2F)][API Docs]
[![Tag](https://img.shields.io/github/v/tag/kjanat/chatlogger-api-go?sort=semver&label=Tag)](https://github.com/kjanat/chatlogger-api-go/tags)
[![Release Date](https://img.shields.io/github/release-date/kjanat/chatlogger-api-go?label=Release%20date)][Latest release]
[![License: MIT](https://img.shields.io/github/license/kjanat/chatlogger-api-go?label=License)](LICENSE)
[![Commit activity](https://img.shields.io/github/commit-activity/m/kjanat/chatlogger-api-go?label=Commit%20activity)][Commits]
[![Last commit](https://img.shields.io/github/last-commit/kjanat/chatlogger-api-go?label=Last%20commit)][Commits]
[![CI](https://img.shields.io/github/actions/workflow/status/kjanat/chatlogger-api-go/chatlogger-pipeline.yml?logo=github&label=CI)][Build]
[![Codecov](https://img.shields.io/codecov/c/gh/kjanat/chatlogger-api-go?token=QUop5QdCOv&logo=codecov&logoColor=%23F01F7A&label=Coverage)][Codecov] <!-- https://codecov.io/gh/kjanat/chatlogger-api-go/graph/badge.svg?token=QUop5QdCOv -->
[![Issues](https://img.shields.io/github/issues/kjanat/chatlogger-api-go?label=Issues)][Issues]
<!-- [![autofix enabled](https://shields.io/badge/autofix.ci-yes-success?logo=data:image/svg+xml;base64,PHN2ZyBmaWxsPSIjZmZmIiB2aWV3Qm94PSIwIDAgMTI4IDEyOCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCB0cmFuc2Zvcm09InNjYWxlKDAuMDYxLC0wLjA2MSkgdHJhbnNsYXRlKC0yNTAsLTE3NTApIiBkPSJNMTMyNSAtMzQwcS0xMTUgMCAtMTY0LjUgMzIuNXQtNDkuNSAxMTQuNXEwIDMyIDUgNzAuNXQxMC41IDcyLjV0NS41IDU0djIyMHEtMzQgLTkgLTY5LjUgLTE0dC03MS41IC01cS0xMzYgMCAtMjUxLjUgNjJ0LTE5MSAxNjl0LTkyLjUgMjQxcS05MCAxMjAgLTkwIDI2NnEwIDEwOCA0OC41IDIwMC41dDEzMiAxNTUuNXQxODguNSA4MXExNSA5OSAxMDAuNSAxODAuNXQyMTcgMTMwLjV0MjgyLjUgNDlxMTM2IDAgMjU2LjUgLTQ2IHQyMDkgLTEyNy41dDEyOC41IC0xODkuNXExNDkgLTgyIDIyNyAtMjEzLjV0NzggLTI5OS41cTAgLTEzNiAtNTggLTI0NnQtMTY1LjUgLTE4NC41dC0yNTYuNSAtMTAzLjVsLTI0MyAtMzAwdi01MnEwIC0yNyAzLjUgLTU2LjV0Ni41IC01Ny41dDMgLTUycTAgLTg1IC00MS41IC0xMTguNXQtMTU3LjUgLTMzLjV6TTEzMjUgLTI2MHE3NyAwIDk4IDE0LjV0MjEgNTcuNXEwIDI5IC0zIDY4dC02LjUgNzN0LTMuNSA0OHY2NGwyMDcgMjQ5IHEtMzEgMCAtNjAgNS41dC01NCAxMi41bC0xMDQgLTEyM3EtMSAzNCAtMiA2My41dC0xIDU0LjVxMCA2OSA5IDEyM2wzMSAyMDBsLTExNSAtMjhsLTQ2IC0yNzFsLTIwNSAyMjZxLTE5IC0xNSAtNDMgLTI4LjV0LTU1IC0yNi41bDIxOSAtMjQydi0yNzZxMCAtMjAgLTUuNSAtNjB0LTEwLjUgLTc5dC01IC01OHEwIC00MCAzMCAtNTMuNXQxMDQgLTEzLjV6TTEyNjIgNjE2cS0xMTkgMCAtMjI5LjUgMzQuNXQtMTkzLjUgOTYuNWw0OCA2NCBxNzMgLTU1IDE3MC41IC04NXQyMDQuNSAtMzBxMTM3IDAgMjQ5IDQ1LjV0MTc5IDEyMXQ2NyAxNjUuNWg4MHEwIC0xMTQgLTc3LjUgLTIwNy41dC0yMDggLTE0OXQtMjg5LjUgLTU1LjV6TTgwMyA1OTVxODAgMCAxNDkgMjkuNXQxMDggNzIuNWwyMjEgLTY3bDMwOSA4NnE0NyAtMzIgMTA0LjUgLTUwdDExNy41IC0xOHE5MSAwIDE2NSAzOHQxMTguNSAxMDMuNXQ0NC41IDE0Ni41cTAgNzYgLTM0LjUgMTQ5dC05NS41IDEzNHQtMTQzIDk5IHEtMzcgMTA3IC0xMTUuNSAxODMuNXQtMTg2IDExNy41dC0yMzAuNSA0MXEtMTAzIDAgLTE5Ny41IC0yNnQtMTY5IC03Mi41dC0xMTcuNSAtMTA4dC00MyAtMTMxLjVxMCAtMzQgMTQuNSAtNjIuNXQ0MC41IC01MC41bC01NSAtNTlxLTM0IDI5IC01NCA2NS41dC0yNSA4MS41cS04MSAtMTggLTE0NSAtNzB0LTEwMSAtMTI1LjV0LTM3IC0xNTguNXEwIC0xMDIgNDguNSAtMTgwLjV0MTI5LjUgLTEyM3QxNzkgLTQ0LjV6Ii8+PC9zdmc+)](https://autofix.ci) -->

A multi-tenant backend API for logging and managing chat sessions, supporting both authenticated and unauthenticated usage with analytics capabilities.

## 🚀 Features

- **Multi-tenancy**: Separate organizations with isolated data
- **Dual Authentication**: API key for chat plugins and JWT for dashboard users
- **Role-based Access Control**: Superadmin, admin, user, and viewer roles with specific permissions
- **Analytics**: Usage metrics per organization with customizable reporting
- **Export Capabilities**: Export chat data in multiple formats (CSV, JSON) with both sync and async options
- **Clean Architecture**: Separation of concerns with layered design (handler → service → repository)
- **Strategy Pattern**: Used for exporter formats (JSON, CSV) and other pluggable components
- **Dependency Injection**: Constructor-based DI for improved testability and maintainability
- **Asynchronous Jobs**: Background processing with Redis and Asynq for resource-intensive tasks

## Documentation

The API is documented using Swagger/OpenAPI:

1. Generate documentation: `./scripts/docs_generate.sh` or `./scripts/docs_generate.ps1`
2. Start the server: `go run cmd/server/main.go`
3. Open browser: [`http://localhost:8080/openapi/index.html`](http://localhost:8080/openapi/index.html)

> [!TIP]
> View the package documentation for more details on the API endpoints and usage.
> [godoc.org][Package documentation]
> The API documentation is also available online at [chatlogger-api-docs.kjanat.com][API Docs].

## 🛠️ Tech Stack

- **Language**: [Go 1.24.2+](https://github.com/kjanat/chatlogger-api-go/blob/af38ff3217524cae62314ead1a60b9b397504d05/go.mod#L3)
- **Web Framework**: [Gin][Gin]
- **ORM**: [GORM][Gorm] with PostgreSQL
- **Authentication**: JWT using [golang-jwt/jwt][JWT]
- **Password Hashing**: bcrypt
- **Configuration**: Environment-based loading
- **Queue System**: [Asynq][Asynq] with Redis (for async exports)
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 14+
- Redis (optional, for async exports)
- Docker & Docker Compose (optional for containerized deployment)

## 🛠️ Development Environment

### Quick Setup

For new contributors, run the automated setup script:

```bash
# Clone the repository
git clone https://github.com/kjanat/chatlogger-api-go.git
cd chatlogger-api-go

# Run the development setup script
./scripts/dev-setup.sh
```

This script will:
- ✅ Check Go and Docker installation
- ✅ Install development tools (golangci-lint, swag, etc.)
- ✅ Set up environment configuration
- ✅ Download Go dependencies
- ✅ Start PostgreSQL and Redis with Docker
- ✅ Run database migrations
- ✅ Generate API documentation
- ✅ Run basic tests to verify setup

### Development Tools

The project includes comprehensive tooling for development:

```bash
# Show all available commands
make help

# Start development environment
make dev

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run benchmarks
make benchmark

# Generate documentation
make docs

# Clean build artifacts
make clean
```

### VS Code Configuration

The project includes VS Code configuration for optimal development:

- **Extensions**: Recommended extensions for Go development
- **Settings**: Configured for Go formatting, linting, and testing
- **Debug**: Pre-configured debug configurations for server and worker
- **Tasks**: Common development tasks accessible via Ctrl+Shift+P

### Code Quality

Automated code quality checks include:

- **Linting**: golangci-lint with comprehensive rules
- **Formatting**: gofmt and goimports
- **Security**: gosec security scanner
- **Testing**: Race condition detection, benchmarks, coverage
- **Dependencies**: Vulnerability scanning with govulncheck

### CI/CD Pipeline

GitHub Actions workflow includes:

- **Lint**: Code quality and style checks
- **Test**: Comprehensive test suite with PostgreSQL and Redis
- **Security**: Security scanning and vulnerability checks
- **Build**: Multi-architecture Docker builds
- **Coverage**: Code coverage reporting to Codecov

## 🚀 Quick Start

### Running with Docker

The easiest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/kjanat/chatlogger-api-go.git
cd chatlogger-api-go

# Start the application and database
docker-compose up -d
```

The API will be available at `http://localhost:8080`.

### Running Locally

```bash
# Clone the repository
git clone https://github.com/kjanat/chatlogger-api-go.git
cd chatlogger-api-go

# Install dependencies
go mod download

# Setup environment variables (copy example and modify)
cp .env.example .env
# Edit .env with your database credentials and settings

# Run migrations
psql -U postgres -d chatlogger -f migrations/001_initial_schema.sql
psql -U postgres -d chatlogger -f migrations/002_ensure_defaults.sql
psql -U postgres -d chatlogger -f migrations/003_add_exports_table.sql

# Run the server
go run cmd/server/main.go

# Run the worker (optional, for async exports)
go run cmd/worker/main.go
```

## 🔑 Authentication

### Chat Plugin (Public API)

Uses API key authentication via the `x-organization-api-key` header:

```bash
curl -X POST \
  http://localhost:8080/v1/orgs/acme/chats \
  -H 'x-organization-api-key: your-api-key' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Support Chat",
    "tags": ["support", "billing"],
    "metadata": {
      "browser": "Chrome",
      "platform": "Windows"
    }
  }'
```

### Dashboard (User Authentication)

Uses email/password login with JWT authentication:

```bash
# Login and get JWT cookie
curl -X POST \
  http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "admin@example.com",
    "password": "your-password"
  }' \
  -c cookies.txt

# Use the JWT cookie for authenticated requests
curl -X GET \
  http://localhost:8080/users/me \
  -b cookies.txt
```

## 🔒 Role-Based Access Control

The API implements a comprehensive role-based access control system:

| Role         | Description                                        | Capabilities                                             |
| :----------- | :------------------------------------------------- | :------------------------------------------------------- |
| `superadmin` | System-wide administrator with unrestricted access | Full access to all endpoints and organizations           |
| `admin`      | Organization administrator                         | Full access within their organization including API keys |
| `user`       | Regular organization user                          | Access to own chats/messages and basic functionality     |
| `viewer`     | Read-only user                                     | View-only access to permitted resources                  |

RBAC is implemented via middleware that checks the user's role before allowing access to protected resources.

## 📁 Project Structure

```plaintext
/cmd
  /server                → Main API server entry point
  /tools                 → Utility tools
  /worker                → Background job worker for async exports
/internal
  /api                   → Gin router and setup
  /config                → Configuration loading
  /domain                → Core models and interfaces
  /handler               → Request handlers
  /hash                  → Password hashing utilities
  /jobs                  → Queue and processor for async tasks
  /middleware            → Auth, RBAC, logging middleware
  /repository            → Database access layer
  /service               → Business logic layer
  /strategy              → Strategy pattern implementations (exporters)
  /version               → Version information
/migrations              → SQL schema migrations
/scripts                 → Utility scripts
```

## 📊 API Endpoints

### System Endpoints

| Method | Endpoint        | Description              |
| :----- | :-------------- | :----------------------- |
| `GET`  | `/health`       | Health check endpoint    |
| `GET`  | `/version`      | API version information  |
| `GET`  | `/openapi/*any` | Swagger UI documentation |
| `GET`  | `/openapi/*any` | OpenAPI documentation    |

### Auth Endpoints

| Method | Endpoint         | Description                 |
| :----- | :--------------- | :-------------------------- |
| `POST` | `/auth/login`    | Login and get JWT cookie    |
| `POST` | `/auth/register` | Register a new user         |
| `POST` | `/auth/logout`   | Logout and clear JWT cookie |

### Dashboard API (Authenticated with JWT)

| Method   | Endpoint                     | Description                    |
| :------- | :--------------------------- | :----------------------------- |
| `GET`    | `/v1/users/me`               | Get current user profile       |
| `PATCH`  | `/v1/users/me`               | Update current user profile    |
| `POST`   | `/v1/users/me/password`      | Change current user's password |
| `POST`   | `/v1/chats`                  | Create a new chat              |
| `GET`    | `/v1/chats`                  | List user's chats              |
| `GET`    | `/v1/chats/:chatID`          | Get a specific chat            |
| `PATCH`  | `/v1/chats/:chatID`          | Update a chat                  |
| `DELETE` | `/v1/chats/:chatID`          | Delete a chat                  |
| `GET`    | `/v1/chats/:chatID/messages` | Get messages from a chat       |
| `GET`    | `/v1/analytics/messages`     | Get message analytics          |

### Admin-Only Endpoints (JWT + Admin Role)

| Method   | Endpoint                  | Description                |
| :------- | :------------------------ | :------------------------- |
| `GET`    | `/v1/orgs/me/apikeys`     | List organization API keys |
| `POST`   | `/v1/orgs/me/apikeys`     | Generate a new API key     |
| `DELETE` | `/v1/orgs/me/apikeys/:id` | Revoke an API key          |

### Export Endpoints (JWT Auth)

| Method | Endpoint                   | Description               |
| :----- | :------------------------- | :------------------------ |
| `POST` | `/v1/exports`              | Create async export job   |
| `GET`  | `/v1/exports`              | List export jobs          |
| `GET`  | `/v1/exports/:id`          | Get export job status     |
| `GET`  | `/v1/exports/:id/download` | Download completed export |
| `POST` | `/v1/exports/sync`         | Create synchronous export |

### Public API (API Key Auth)

| Method | Endpoint                                | Description               |
| :----- | :-------------------------------------- | :------------------------ |
| `POST` | `/v1/orgs/:slug/chats`                  | Create a new chat session |
| `POST` | `/v1/orgs/:slug/chats/:chatID/messages` | Add a message to a chat   |

## 🔧 Configuration

The application is configured using environment variables:

| Variable       | Description                     | Default          |
| :------------- | :------------------------------ | :--------------- |
| `PORT`         | Server port                     | `8080`           |
| `DATABASE_URL` | PostgreSQL connection string    | Required         |
| `JWT_SECRET`   | Secret for JWT signing          | Required         |
| `REDIS_ADDR`   | Redis address for job queue     | `localhost:6379` |
| `EXPORT_DIR`   | Directory to store export files | `./exports`      |

## ⚙️ Export System

The application supports both synchronous and asynchronous exports:

### Synchronous Exports

- Immediate response with download
- Good for small data sets
- Available at `/exports/orgs/me/sync`
- Works without Redis

### Asynchronous Exports

- Queued processing with status tracking
- Better for large data sets
- Requires Redis and worker process
- Creates a job at `/exports/orgs/me`
- Check status at `/exports/orgs/me/:id`
- Download at `/exports/orgs/me/:id/download`

Both export types support JSON and CSV formats.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linting checks
golangci-lint run
```

## 🚢 Deployment

### Version Information

The application includes version information that can be injected during build:

```bash
# Build with version information
docker build -t chatlogger-api \
  --build-arg VERSION=0.5.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  .
```

### Image Verification

All container images are signed using [Sigstore cosign](https://github.com/sigstore/cosign). You can verify the authenticity of the images using the following command:

```bash
# Verify the container image signature
cosign verify \
    --key=cosign.pub \
    ghcr.io/kjanat/chatlogger-api-server:0.5.0

# Or use the public key URL
cosign verify \
    --key=https://raw.githubusercontent.com/kjanat/chatlogger-api-go/master/cosign.pub \
    ghcr.io/kjanat/chatlogger-api-worker:0.5.0
```

The public key for verification ([`cosign.pub`](cosign.pub)) is available in the root of the repository. For production deployments, we recommend verifying image signatures as part of your CI/CD pipeline to ensure supply chain security.

### CI/CD

The project includes GitHub Actions workflows for CI/CD:

- **Build Workflow**: Builds, tests, and tags versions on pushes to main
- **Release Workflow**: Creates GitHub releases when tags are pushed

## 📝 Next Steps and Enhancements

Here are some potential enhancements planned for future development:

- [x] **Export Features**: Implemented export functionality as a strategy pattern (JSON/CSV)
- [x] **Async Exports**: Added background processing for large exports
- [x] **Documentation**: Generated API documentation with Swagger/OpenAPI
- [ ] **Pagination**: Add cursor-based pagination to list endpoints
- [ ] **Enhanced Testing**: Expand unit and integration test coverage
- [ ] **Monitoring**: Add structured logging and metrics collection
- [ ] **Real-time notifications**: Add WebSocket support for live updates
- [ ] **Multi-factor Authentication**: Add 2FA support for dashboard users
- [ ] **Improved Analytics**: More detailed analytics dashboards and reports
- [ ] **Search Capabilities**: Add full-text search for chat messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

[Commits]: https://github.com/kjanat/chatlogger-api-go/commits/master/
[Latest release]: https://github.com/kjanat/chatlogger-api-go/releases/latest
[Issues]: https://github.com/kjanat/chatlogger-api-go/issues
[Go report]: https://goreportcard.com/report/github.com/kjanat/chatlogger-api-go
[Gin]: https://github.com/gin-gonic/gin
[Gorm]: https://gorm.io/
[JWT]: https://github.com/golang-jwt/jwt/tree/v5
[Asynq]: https://github.com/hibiken/asynq
[Codecov]: https://codecov.io/gh/kjanat/chatlogger-api-go
[Build]: https://github.com/kjanat/chatlogger-api-go/actions/workflows/chatlogger-pipeline.yml
[Package documentation]: https://godoc.org/github.com/kjanat/chatlogger-api-go
[API Docs]: https://chatlogger-api-docs.kjanat.com/
