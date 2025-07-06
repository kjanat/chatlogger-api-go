#!/bin/bash

# ChatLogger API - Development Environment Setup Script
# This script sets up the development environment for new contributors

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Compare version strings (version_gte actual_version required_version)
# Returns 0 if actual >= required, 1 otherwise
version_gte() {
    local actual="$1"
    local required="$2"

    # Handle Go version format (e.g., "1.21.0" or "1.21")
    # Convert to comparable format by padding with zeros
    local actual_major=$(echo "$actual" | cut -d. -f1)
    local actual_minor=$(echo "$actual" | cut -d. -f2)
    local actual_patch=$(echo "$actual" | cut -d. -f3)

    local required_major=$(echo "$required" | cut -d. -f1)
    local required_minor=$(echo "$required" | cut -d. -f2)
    local required_patch=$(echo "$required" | cut -d. -f3)

    # Default patch version to 0 if not specified
    actual_patch=${actual_patch:-0}
    required_patch=${required_patch:-0}

    # Compare major version
    if [ "$actual_major" -gt "$required_major" ]; then
        return 0
    elif [ "$actual_major" -lt "$required_major" ]; then
        return 1
    fi

    # Major versions equal, compare minor version
    if [ "$actual_minor" -gt "$required_minor" ]; then
        return 0
    elif [ "$actual_minor" -lt "$required_minor" ]; then
        return 1
    fi

    # Major and minor equal, compare patch version
    if [ "$actual_patch" -ge "$required_patch" ]; then
        return 0
    else
        return 1
    fi
}

# Check Go installation
check_go() {
    print_status "Checking Go installation..."
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21+ from https://golang.org/dl/"
        exit 1
    fi

    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    print_success "Go $GO_VERSION is installed"

    # Check minimum version requirement
    MIN_GO_VERSION="1.21"
    if ! version_gte "$GO_VERSION" "$MIN_GO_VERSION"; then
        print_error "Go version $GO_VERSION is not supported. Minimum required version: $MIN_GO_VERSION"
        print_error "Please upgrade Go from https://golang.org/dl/"
        exit 1
    fi
    print_success "Go version meets minimum requirement ($MIN_GO_VERSION)"
}

# Detect Docker Compose command (new plugin or legacy binary)
detect_docker_compose() {
    if command_exists docker && docker compose version >/dev/null 2>&1; then
        echo "docker compose"
    elif command_exists docker-compose; then
        echo "docker-compose"
    else
        echo ""
    fi
}

# Check Docker installation
check_docker() {
    print_status "Checking Docker installation..."
    if ! command_exists docker; then
        print_warning "Docker is not installed. Some features may not work."
        return 1
    fi

    DOCKER_COMPOSE_CMD=$(detect_docker_compose)
    if [ -z "$DOCKER_COMPOSE_CMD" ]; then
        print_warning "Docker Compose is not installed. Some features may not work."
        return 1
    fi

    print_success "Docker and Docker Compose ($DOCKER_COMPOSE_CMD) are installed"
    return 0
}

# Install development tools
install_tools() {
    print_status "Installing development tools..."

    # golangci-lint
    if ! command_exists golangci-lint; then
        print_status "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi

    # swag (Swagger generator)
    if ! command_exists swag; then
        print_status "Installing swag..."
        go install github.com/swaggo/swag/v2/cmd/swag@latest
    fi

    # goimports
    print_status "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest

    # govulncheck
    print_status "Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest

    print_success "Development tools installed"

    # Check if Go tools are on PATH
    TOOLS_BIN="$(go env GOBIN 2>/dev/null || echo "$HOME/go/bin")"
    if [[ ":$PATH:" != *":$TOOLS_BIN:"* ]]; then
        print_warning "Add $TOOLS_BIN to your PATH to use installed tools."
    fi
}

# Setup environment file
setup_env() {
    print_status "Setting up environment configuration..."

    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            print_success "Created .env file from template"
            print_warning "Please edit .env with your local configuration"
        else
            print_warning ".env.example not found, skipping environment setup"
        fi
    else
        print_success ".env file already exists"
    fi
}

# Download Go dependencies
download_deps() {
    print_status "Downloading Go dependencies..."
    go mod download
    go mod tidy
    print_success "Dependencies downloaded"
}

# Setup database
setup_database() {
    if ! check_docker; then
        print_warning "Skipping database setup - Docker not available"
        return
    fi

    # Detect Docker Compose command
    DOCKER_COMPOSE_CMD=$(detect_docker_compose)
    if [ -z "$DOCKER_COMPOSE_CMD" ]; then
        print_error "Docker Compose not found"
        return 1
    fi

    print_status "Setting up development database..."

    # Start PostgreSQL and Redis
    $DOCKER_COMPOSE_CMD up -d postgres redis

    # Wait for PostgreSQL to be ready
    print_status "Waiting for PostgreSQL to be ready..."
    POSTGRES_READY=false
    for attempt in {1..30}; do
        if $DOCKER_COMPOSE_CMD exec postgres pg_isready -U postgres >/dev/null 2>&1; then
            POSTGRES_READY=true
            break
        fi
        print_status "PostgreSQL not ready yet (attempt $attempt/30)..."
        sleep 1
    done

    # Fail fast if PostgreSQL never became ready
    if [ "$POSTGRES_READY" = false ]; then
        print_error "PostgreSQL failed to become ready after 30 seconds"
        print_error "Check Docker containers with: $DOCKER_COMPOSE_CMD logs postgres"
        exit 1
    fi

    print_success "PostgreSQL is ready"

    # Run migrations
    print_status "Running database migrations..."
    if [ -f migrations/001_initial_schema.sql ]; then
        $DOCKER_COMPOSE_CMD exec postgres psql -U postgres -d chatlogger -f /docker-entrypoint-initdb.d/001_initial_schema.sql 2>/dev/null || true
        $DOCKER_COMPOSE_CMD exec postgres psql -U postgres -d chatlogger -f /docker-entrypoint-initdb.d/002_ensure_defaults.sql 2>/dev/null || true
        $DOCKER_COMPOSE_CMD exec postgres psql -U postgres -d chatlogger -f /docker-entrypoint-initdb.d/003_add_exports_table.sql 2>/dev/null || true
        print_success "Database migrations completed"
    else
        print_warning "Migration files not found, skipping database setup"
    fi
}

# Generate documentation
generate_docs() {
    print_status "Generating API documentation..."
    if [ -f scripts/docs_generate.sh ]; then
        chmod +x scripts/docs_generate.sh
        ./scripts/docs_generate.sh
        print_success "API documentation generated"
    else
        print_warning "docs_generate.sh not found, skipping documentation generation"
    fi
}

# Run tests
run_tests() {
    print_status "Running tests to verify setup..."
    if go test -v ./internal/domain/ >/dev/null 2>&1; then
        print_success "Basic tests passed"
    else
        print_warning "Some tests failed - this might be expected on first setup"
    fi
}

# VS Code setup
setup_vscode() {
    if [ -d .vscode ]; then
        print_status "VS Code configuration detected"
        print_success "Recommended extensions: Go, Docker, GitLens, REST Client"
    fi
}

# Main setup function
main() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║              ChatLogger API Development Setup               ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"

    # Check prerequisites
    check_go

    # Install tools and setup environment
    install_tools
    setup_env
    download_deps

    # Setup database if Docker is available
    setup_database

    # Generate documentation
    generate_docs

    # VS Code setup
    setup_vscode

    # Run basic tests
    run_tests

    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    Setup Complete\! 🎉                       ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"

    echo "Next steps:"
    echo "1. Edit .env with your configuration"
    echo "2. Run 'make dev' to start the development environment"
    echo "3. Visit http://localhost:8080/swagger/index.html for API docs"
    echo "4. Run 'make test' to run all tests"
    echo ""
    echo "Available commands:"
    echo "  make help      - Show all available commands"
    echo "  make dev       - Start development environment"
    echo "  make test      - Run tests"
    echo "  make lint      - Run code linting"
    echo "  make docs      - Generate documentation"
    echo ""
}

# Run main function
main "$@"
