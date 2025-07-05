# Development Environment Guide

This document provides comprehensive information about the development environment setup for ChatLogger API.

## 📋 Quick Start

```bash
# Clone and setup
git clone https://github.com/kjanat/chatlogger-api-go.git
cd chatlogger-api-go
./scripts/dev-setup.sh

# Start development environment
make dev
```

## 🛠️ Development Tools

### Core Tools
- **Go 1.24.2+**: Primary language
- **golangci-lint**: Code linting and quality checks
- **swag**: Swagger/OpenAPI documentation generation
- **goimports**: Import formatting
- **govulncheck**: Security vulnerability scanning

### Development Environment
- **Air**: Hot reloading for Go applications
- **Docker & Docker Compose**: Containerized development
- **PostgreSQL 16**: Primary database
- **Redis 7**: Job queue and caching

### Editor Support
- **VS Code**: Full configuration included
  - Go extension with debugging
  - Recommended extensions
  - Tasks and launch configurations
  - Code formatting and linting integration

## 🏗️ Project Structure

```
chatlogger-api-go/
├── .github/workflows/     # CI/CD pipeline
├── .vscode/              # VS Code configuration
├── cmd/                  # Application entry points
│   ├── server/          # API server
│   └── worker/          # Background worker
├── internal/            # Private application code
│   ├── domain/         # Core business models
│   ├── handler/        # HTTP handlers
│   ├── service/        # Business logic
│   ├── repository/     # Data access
│   └── middleware/     # HTTP middleware
├── migrations/          # Database migrations
├── scripts/            # Development scripts
├── test/               # Test utilities and fixtures
├── Makefile           # Development commands
└── docker-compose.yml # Container orchestration
```

## 🚀 Development Commands

### Make Commands
```bash
make help              # Show all available commands
make setup             # Install development tools
make dev               # Start development environment
make build             # Build binaries
make test              # Run tests
make test-coverage     # Run tests with coverage
make test-integration  # Run integration tests
make benchmark         # Run performance benchmarks
make race              # Run race condition tests
make lint              # Run code linting
make format            # Format code
make security          # Run security checks
make docs              # Generate documentation
make clean             # Clean build artifacts
make docker            # Start with Docker Compose
make migration         # Run database migrations
make ci                # Run full CI pipeline locally
```

### VS Code Tasks
- **Build Server** (Ctrl+Shift+P → Tasks: Run Task)
- **Run Tests** 
- **Generate Docs**
- **Start Docker**
- **Run Linter**

## 🧪 Testing Strategy

### Test Types
1. **Unit Tests**: Fast, isolated tests for individual components
2. **Integration Tests**: Database and external service integration
3. **Contract Tests**: API endpoint contract validation
4. **Benchmark Tests**: Performance and scalability testing
5. **Race Condition Tests**: Concurrency safety validation

### Test Infrastructure
- **SQLite In-Memory**: Fast unit testing database
- **PostgreSQL**: Integration testing with real database
- **Test Fixtures**: Consistent test data generation
- **Mock Services**: Isolated component testing

### Coverage Goals
- **Target**: >75% code coverage
- **Current**: ~27% (baseline established)
- **Reports**: Generated in `coverage.html`

## 🔧 Configuration

### Environment Variables
Copy `.env.example` to `.env` and configure:

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/chatlogger?sslmode=disable
TEST_DATABASE_URL=postgres://postgres:password@localhost:5432/chatlogger_test?sslmode=disable

# Server
PORT=8080
GIN_MODE=debug
JWT_SECRET=your-secret-key

# Redis (optional)
REDIS_ADDR=localhost:6379

# Development
LOG_LEVEL=debug
DEBUG_ENDPOINTS=true
```

### Docker Development
```bash
# Start all services
docker-compose -f docker-compose.dev.yml up -d

# With optional tools (pgAdmin, Redis Commander)
docker-compose -f docker-compose.dev.yml --profile tools up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f api-dev
```

## 📊 Code Quality

### Linting Rules
- **golangci-lint**: Comprehensive Go linting
- **gofmt**: Code formatting
- **goimports**: Import organization
- **gosec**: Security scanning
- **staticcheck**: Advanced static analysis

### Pre-commit Hooks
Configured via GitHub Actions and Makefile:
- Format checking
- Lint validation
- Security scanning
- Test execution

### CI/CD Pipeline
GitHub Actions workflow includes:
1. **Lint Job**: Code quality validation
2. **Test Job**: Comprehensive testing with real services
3. **Race Condition Tests**: Concurrency validation
4. **Security Scan**: Vulnerability detection
5. **Build Job**: Multi-architecture builds
6. **Docker Build**: Container image creation

## 🐛 Debugging

### Local Debugging
1. **VS Code**: Use F5 to start debugging
2. **Delve**: Direct debugger integration
3. **Logs**: Structured logging with levels

### Debug Configurations
- **Launch Server**: Debug API server
- **Launch Worker**: Debug background worker
- **Debug Tests**: Debug specific tests
- **Attach to Process**: Debug running containers

### Hot Reloading
Air provides automatic rebuilding:
```bash
# Server hot reload
air -c .air.toml

# Worker hot reload  
air -c .air.worker.toml
```

## 📈 Performance Monitoring

### Benchmarks
```bash
# Run all benchmarks
make benchmark

# Specific benchmarks
go test -bench=BenchmarkChatRepository ./internal/repository/
go test -bench=BenchmarkChatService ./internal/service/
```

### Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# Analyze profiles
go tool pprof cpu.prof
```

## 🔒 Security

### Security Scanning
- **gosec**: Static security analysis
- **govulncheck**: Dependency vulnerability scanning
- **Trivy**: Container security scanning
- **CodeQL**: GitHub security analysis

### Best Practices
- Environment variable validation
- Input sanitization
- SQL injection prevention
- JWT token security
- CORS configuration

## 📚 Documentation

### API Documentation
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI Spec**: Generated from code annotations
- **Postman Collection**: Available in docs/

### Code Documentation
- **GoDoc**: Package documentation
- **README**: Project overview
- **CLAUDE.md**: AI assistant guidance
- **Architecture**: Clean architecture documentation

## 🚀 Deployment

### Local Development
```bash
make dev                 # Full development environment
make run                # Run server only
make run-worker         # Run worker only
```

### Production Builds
```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server

# Docker production build
docker build -t chatlogger-api:latest .
```

### Health Checks
- **Health Endpoint**: `/health`
- **Version Endpoint**: `/version`
- **Metrics**: Prometheus-compatible metrics (planned)

## 🤝 Contributing

### Development Workflow
1. Fork the repository
2. Create feature branch
3. Run `./scripts/dev-setup.sh`
4. Develop with `make dev`
5. Test with `make ci`
6. Submit pull request

### Code Standards
- Go formatting with `gofmt`
- Import organization with `goimports`
- Comprehensive test coverage
- Security-first development
- Clean architecture principles

### Pull Request Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Linting passes
- [ ] Security scan clean
- [ ] Benchmarks considered
- [ ] CLAUDE.md updated if needed

## 🆘 Troubleshooting

### Common Issues

**Database Connection Failed**
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check connection string
make migration
```

**Tests Failing**
```bash
# Clean test cache
go clean -testcache

# Check test database
export TEST_DATABASE_URL="postgres://..."
make test
```

**Hot Reload Not Working**
```bash
# Check Air installation
air -v

# Restart Air
pkill air && make dev
```

**Docker Issues**
```bash
# Clean Docker state
docker-compose down -v
docker system prune -f
make docker
```

### Getting Help
- Check GitHub Issues
- Review CI/CD logs
- Use VS Code integrated debugging
- Check Docker container logs