# Aurora - Backend as a Service Platform

**Phase 1: Core Foundation**

Aurora is a Backend as a Service (BaaS) platform that enables developers to quickly create and manage scalable backend services. Phase 1 provides the core foundation with project management and dynamic schema definition capabilities.

---

## For New Developers

**Welcome!** This guide will help you get started with Aurora development.

**Quick Links:**
- [Prerequisites](#prerequisites) - What you need installed
- [Quick Start](#quick-start) - Get up and running in 5 minutes
- [Architecture Overview](#architecture-overview) - Understand how it works
- [Development Workflow](#development-workflow) - How to add features
- [Testing Guide](#testing-guide) - How to test your changes
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Code style and contribution guidelines
- [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Deep dive into system design
- [docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md) - Day-to-day development guide

**What is Aurora?**

Aurora is a platform that lets developers create backend services without writing boilerplate code. Think of it as "infrastructure as code" for your backend:

1. **Create a project** - Choose your language (TypeScript/Python), database (PostgreSQL/MongoDB), and API style (REST/GraphQL)
2. **Define your schema** - Write a simple YAML file describing your data models
3. **Get your API** - Aurora generates migrations and provides REST endpoints automatically

**Current Phase:** Core Foundation (Project & Schema Management)
**Next Phase:** Code generation for client SDKs

---

## Overview

Aurora simplifies backend development by providing a unified platform to:
- Create and manage BaaS projects with configurable languages, databases, and API styles
- Define database schemas dynamically using YAML format
- Version schema changes automatically
- Manage database migrations automatically
- Access your backend through a clean REST API

The Phase 1 core foundation provides essential project and schema management capabilities, with support for PostgreSQL and MongoDB databases.

**Architecture at a Glance:**

```
┌─────────────┐
│   Client    │
│  (HTTP API) │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────┐
│           API Layer (Gin)               │
│  ┌─────────┐  ┌─────────┐  ┌────────┐  │
│  │ Health  │  │ Project │  │ Schema │  │
│  │ Handler │  │ Handler │  │ Handler│  │
│  └─────────┘  └─────────┘  └────────┘  │
└───────────────────┬─────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│        Service Layer (Business Logic)   │
│  ┌──────────────┐  ┌──────────────────┐ │
│  │   Project    │  │     Schema       │ │
│  │   Service    │  │     Service      │ │
│  │  (Validation)│  │  (Validation +   │ │
│  │              │  │   Versioning)    │ │
│  └──────────────┘  └──────────────────┘ │
└───────────────────┬─────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│    Repository Layer (Data Access)       │
│  ┌──────────────┐  ┌──────────────────┐ │
│  │   Project    │  │     Schema       │ │
│  │  Repository  │  │   Repository     │ │
│  │   (GORM)     │  │     (GORM)       │ │
│  └──────────────┘  └──────────────────┘ │
└───────────────────┬─────────────────────┘
                    │
                    ▼
           ┌────────────────┐
           │   PostgreSQL   │
           │   (Database)   │
           └────────────────┘
```

## Features

- ✅ Project Management - Create, list, update, and delete BaaS projects
- ✅ Multi-Language Support - Projects can target TypeScript or Python
- ✅ Database Flexibility - Support for PostgreSQL and MongoDB
- ✅ API Style Options - REST or GraphQL API generation
- ✅ Dynamic Schema Definition - Define tables and columns using YAML format
- ✅ Schema Versioning - Automatic version tracking for schema changes
- ✅ Schema Validation - Comprehensive validation of schema definitions
- ✅ Database Migrations - Automatic migration execution on startup
- ✅ Health Checks - Built-in health check endpoint
- ✅ Error Handling - Comprehensive error handling and validation
- ✅ Structured Logging - Request and operation logging

## Tech Stack

**Core**
- Go 1.21+ - Primary language
- Gin Web Framework - HTTP routing and middleware
- GORM - Object-relational mapping

**Database**
- PostgreSQL - Primary datastore (Phase 1)
- MongoDB - Alternative datastore (Phase 1)

**Configuration**
- Environment variables (.env support)
- godotenv - Environment file loading

**Testing**
- Go testing package - Unit tests
- Integration tests with Docker Compose

## Prerequisites

Before starting, ensure you have the following installed:

| Tool | Version | Purpose | Installation |
|------|---------|---------|-------------|
| **Go** | 1.26.2+ | Primary language | [golang.org/doc/install](https://golang.org/doc/install) |
| **Docker** | 20.10+ | Database container | [docs.docker.com/get-docker](https://docs.docker.com/get-docker/) |
| **Docker Compose** | 2.0+ | Multi-container orchestration | Usually included with Docker Desktop |
| **Make** | Any | Build automation | Pre-installed on macOS/Linux; Windows: [gnuwin32.sourceforge.net](http://gnuwin32.sourceforge.net/packages/make.htm) |
| **Git** | 2.0+ | Version control | [git-scm.com](https://git-scm.com/downloads) |
| **curl** or **Postman** | Any | API testing | Pre-installed on macOS/Linux; [postman.com](https://www.postman.com/) for GUI |

**Verify Installation:**

```bash
go version        # Should show go1.26.2 or higher
docker --version  # Should show 20.10 or higher
docker-compose --version  # Should show 2.0 or higher
make --version    # Should show GNU Make
```

---

## Quick Start

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/ktarun.reddy/baas.git
   cd baas
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```

4. **Start the database (using Docker Compose)**
   ```bash
   docker-compose up -d
   ```

5. **Run database migrations** (automatic on server startup)
   ```bash
   go run cmd/server/main.go
   ```

6. **Verify the server is running**
   ```bash
   curl http://localhost:8080/health
   ```

   Expected response:
   ```json
   {"status":"healthy","timestamp":"2026-05-11T..."}
   ```

**Troubleshooting Common Setup Issues:**

| Issue | Solution |
|-------|----------|
| **Port 5432 already in use** | Another PostgreSQL is running. Stop it: `sudo service postgresql stop` (Linux) or `brew services stop postgresql` (macOS) |
| **Port 8080 already in use** | Change PORT in `.env` file to 8081 or kill process: `lsof -ti:8080 \| xargs kill` |
| **`go: connection refused`** | Check your internet connection or use Go proxy: `export GOPROXY=https://proxy.golang.org,direct` |
| **Docker not found** | Start Docker Desktop or install Docker daemon |
| **Permission denied (Docker)** | Add user to docker group: `sudo usermod -aG docker $USER` then logout/login |

### Development

Start the full development environment:
```bash
make dev
```

This will:
1. Start PostgreSQL via Docker Compose (port 5432)
2. Wait for database to be ready
3. Run the server on port 8080
4. Auto-run migrations on startup

**Alternative: Manual startup**
```bash
# Terminal 1: Start database
docker-compose up

# Terminal 2: Start server
go run cmd/server/main.go
```

## API Endpoints

### Health Check
- **GET** `/health` - Server health status

### Projects Management
- **POST** `/api/v1/projects` - Create a new project
- **GET** `/api/v1/projects` - List all projects
- **GET** `/api/v1/projects/:id` - Get project details
- **PUT** `/api/v1/projects/:id` - Update project
- **DELETE** `/api/v1/projects/:id` - Delete project

### Schema Management
- **POST** `/api/v1/projects/:id/schemas` - Apply schema to a project
- **GET** `/api/v1/projects/:id/schemas` - List all schemas for a project
- **GET** `/api/v1/projects/:id/schemas/latest` - Get the latest schema version

## Schema Definition Format

Schemas are defined using YAML format. Each schema must specify the database type, version, and table definitions.

### Example Schema

```yaml
version: "1.0"
database: "postgres"
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
        primary_key: true
      - name: "email"
        type: "string"
        unique: true
        nullable: false
      - name: "name"
        type: "string"
        nullable: false
      - name: "created_at"
        type: "timestamp"
        nullable: false
        default: "now()"
    indexes:
      - columns: ["email"]
        unique: true
  
  - name: "posts"
    columns:
      - name: "id"
        type: "uuid"
        primary_key: true
      - name: "user_id"
        type: "uuid"
        nullable: false
        foreign_key: "users.id"
      - name: "title"
        type: "string"
        nullable: false
      - name: "content"
        type: "text"
        nullable: true
      - name: "published"
        type: "boolean"
        default: "false"
      - name: "created_at"
        type: "timestamp"
        nullable: false
        default: "now()"
    indexes:
      - columns: ["user_id"]
      - columns: ["created_at"]
```

### Column Definition

Each column supports the following attributes:

- `name` - Column name (required)
- `type` - Column data type (required)
- `primary_key` - Whether this is a primary key (boolean)
- `unique` - Enforce unique constraint (boolean)
- `nullable` - Allow null values (boolean)
- `default` - Default value expression
- `foreign_key` - Foreign key reference (table.column format)

## Supported Column Types

Aurora supports the following column types for schema definitions:

- **string** - Variable-length string (maps to VARCHAR)
- **text** - Large text content
- **integer** - 32-bit signed integer
- **bigint** - 64-bit signed integer
- **decimal** - Arbitrary precision decimal number
- **float** - Floating-point number
- **boolean** - True/false boolean value
- **date** - Date value (YYYY-MM-DD)
- **timestamp** - Date and time value with timezone
- **uuid** - Universally unique identifier
- **json** - JSON document storage

---

## Architecture Overview

Aurora follows a **clean architecture** pattern with clear separation of concerns:

### Layers Explained

1. **API Layer** (`internal/api/`)
   - **Purpose:** Handle HTTP requests/responses
   - **Technology:** Gin web framework
   - **Components:**
     - `router.go` - Route definitions and middleware setup
     - `handlers/` - Request handlers (health, project, schema)
     - `middleware/` - Cross-cutting concerns (logging, error handling)
   - **Responsibility:** Parse requests, call services, format responses

2. **Service Layer** (`internal/service/`)
   - **Purpose:** Business logic and validation
   - **Components:**
     - `project_service.go` - Project operations & validation
     - `schema_service.go` - Schema operations, versioning, validation
   - **Responsibility:** Enforce business rules, coordinate repositories
   - **Why?** Keeps business logic separate from HTTP and database concerns

3. **Repository Layer** (`internal/repository/`)
   - **Purpose:** Data access and persistence
   - **Technology:** GORM (ORM)
   - **Components:**
     - `project_repo.go` - CRUD operations for projects
     - `schema_repo.go` - CRUD + versioning for schemas
   - **Responsibility:** Database queries, transactions
   - **Why?** Abstracts database details from business logic

4. **Domain Layer** (`internal/domain/`)
   - **Purpose:** Core business entities
   - **Components:**
     - `project.go` - Project entity and validation
     - `schema.go` - Schema entity and validation
   - **Responsibility:** Define data structures and constraints
   - **Why?** Single source of truth for data models

5. **Database Layer** (`internal/database/`)
   - **Purpose:** Database connection and migrations
   - **Components:**
     - `postgres.go` - Connection management
     - `migrations.go` - Schema migrations (auto-run on startup)
   - **Responsibility:** Database initialization and schema versioning

6. **Package Layer** (`pkg/`)
   - **Purpose:** Reusable utilities
   - **Components:**
     - `validator/schema_validator.go` - YAML schema validation
   - **Responsibility:** Shared logic that could be used by other projects
   - **Why?** `pkg/` is for exportable code, `internal/` is private to this project

### Request Flow Example

When a client creates a project:

```
1. POST /api/v1/projects
   ↓
2. router.go routes to projectHandler.Create
   ↓
3. Handler parses JSON body
   ↓
4. Handler calls projectService.CreateProject()
   ↓
5. Service validates project (language, database, API style)
   ↓
6. Service calls projectRepository.Create()
   ↓
7. Repository inserts into PostgreSQL via GORM
   ↓
8. Success bubbles back up through layers
   ↓
9. Handler returns 201 Created with project JSON
```

**Why this architecture?**
- **Testability:** Each layer can be tested independently
- **Maintainability:** Changes to one layer don't affect others
- **Scalability:** Easy to add new features or swap implementations
- **Clarity:** New developers can understand each layer's purpose

See [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) for detailed design decisions.

---

## Code Organization

### Project Structure

```
baas/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point - wires everything together
│
├── internal/                    # Private application code
│   ├── api/
│   │   ├── router.go            # HTTP routes and middleware setup
│   │   ├── handlers/
│   │   │   ├── health.go        # Health check endpoint
│   │   │   ├── project.go       # Project CRUD endpoints
│   │   │   └── schema.go        # Schema management endpoints
│   │   └── middleware/
│   │       └── logger.go        # Request logging middleware
│   │
│   ├── config/
│   │   ├── config.go            # Configuration loading from env vars
│   │   └── config_test.go       # Config validation tests
│   │
│   ├── database/
│   │   ├── postgres.go          # Database connection management
│   │   └── migrations.go        # Auto-migrations for projects/schemas tables
│   │
│   ├── domain/
│   │   ├── project.go           # Project entity + validation rules
│   │   └── schema.go            # Schema entity + validation rules
│   │
│   ├── repository/
│   │   ├── project_repo.go      # Project data access (CRUD)
│   │   └── schema_repo.go       # Schema data access (CRUD + versioning)
│   │
│   └── service/
│       ├── project_service.go   # Project business logic
│       └── schema_service.go    # Schema business logic + versioning
│
├── pkg/                         # Reusable packages (could be exported)
│   └── validator/
│       └── schema_validator.go  # YAML schema validation logic
│
├── tests/
│   ├── integration/
│   │   └── api_test.go          # End-to-end API tests
│   └── testutils/
│       └── fixtures.go          # Test data fixtures
│
├── docs/                        # Documentation
│   ├── ARCHITECTURE.md          # Deep dive into system design
│   └── DEVELOPMENT.md           # Development guide
│
├── .env.example                 # Example environment variables
├── docker-compose.yaml          # Local PostgreSQL setup
├── go.mod                       # Go module definition
├── Makefile                     # Development commands
├── CONTRIBUTING.md              # Contribution guidelines
└── README.md                    # This file
```

### Key Files Explained

| File | Purpose | When to Modify |
|------|---------|---------------|
| `cmd/server/main.go` | Application bootstrap | When adding new top-level dependencies |
| `internal/api/router.go` | Route definitions | When adding new API endpoints |
| `internal/domain/*.go` | Data models | When adding new entities or fields |
| `internal/service/*.go` | Business logic | When adding validation or business rules |
| `internal/repository/*.go` | Database queries | When adding new database operations |
| `pkg/validator/*.go` | Schema validation | When extending schema validation rules |

---

## Development Workflow

### How to Add a New Feature

Let's walk through adding a new feature: **Project Tags** (allowing users to tag projects).

**Step 1: Define the Domain Model**

```bash
# Edit internal/domain/project.go
```

Add a `Tags` field:
```go
type Project struct {
    // ... existing fields
    Tags []string `gorm:"type:text[]" json:"tags"`
}
```

**Step 2: Update Repository**

```bash
# No changes needed - GORM handles new field automatically
# But if you need custom queries:
# Edit internal/repository/project_repo.go
```

**Step 3: Update Service**

```bash
# Edit internal/service/project_service.go
```

Add validation:
```go
func (s *ProjectService) CreateProject(project *domain.Project) error {
    // Add tag validation
    if len(project.Tags) > 10 {
        return fmt.Errorf("maximum 10 tags allowed")
    }
    // ... rest of validation
}
```

**Step 4: Update API Handler**

```bash
# Edit internal/api/handlers/project.go
```

The handler automatically handles the new field via JSON binding!

**Step 5: Test Your Changes**

```bash
# Add unit tests
vim internal/service/project_service_test.go

# Run tests
make test

# Test manually
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Project",
    "language": "typescript",
    "database_type": "postgres",
    "api_style": "rest",
    "tags": ["web", "production"]
  }'
```

**Step 6: Update Documentation**

Update this README and create a migration if needed.

See [docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md) for more examples.

---

## Testing Guide

### Running Tests

```bash
# Run all unit tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/service/...

# Run integration tests (requires Docker)
make test-integration

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View in browser
```

### Test Organization

**Unit Tests** (`*_test.go` files alongside source):
- Test individual functions/methods
- Mock external dependencies
- Fast execution (< 1 second)
- Example: `internal/service/project_service_test.go`

**Integration Tests** (`tests/integration/`):
- Test complete request/response flows
- Use real database (Docker)
- Slower execution (5-10 seconds)
- Example: `tests/integration/api_test.go`

### Writing Tests

**Unit Test Example:**

```go
// internal/service/project_service_test.go
func TestCreateProject_InvalidLanguage(t *testing.T) {
    // Arrange
    repo := &mockProjectRepository{}
    service := NewProjectService(repo)
    project := &domain.Project{
        Name:     "Test",
        Language: "invalid", // Invalid language
    }

    // Act
    err := service.CreateProject(project)

    // Assert
    if err == nil {
        t.Error("Expected error for invalid language")
    }
}
```

**Integration Test Example:**

```go
// tests/integration/api_test.go
func TestCreateProject_EndToEnd(t *testing.T) {
    // Setup test server
    router := setupTestRouter()
    
    // Make request
    body := `{"name":"Test","language":"typescript","database_type":"postgres","api_style":"rest"}`
    req := httptest.NewRequest("POST", "/api/v1/projects", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, 201, w.Code)
}
```

### Test Coverage Goals

- **Domain layer:** 100% (pure logic, easy to test)
- **Service layer:** 90%+ (critical business logic)
- **Repository layer:** 80%+ (integration tests)
- **API handlers:** 80%+ (integration tests)

---

## Build & Deployment

### Building the Binary

```bash
# Build for current platform
make build

# Output: bin/baas

# Run the binary
./bin/baas

# Build for Linux (for deployment)
GOOS=linux GOARCH=amd64 go build -o bin/baas-linux cmd/server/main.go
```

### Clean Build Artifacts

```bash
make clean  # Removes bin/ directory
```

---

## Common Issues & Solutions

### Issue: "connection refused" when starting server

**Cause:** PostgreSQL isn't running or isn't ready yet.

**Solution:**
```bash
# Check if PostgreSQL container is running
docker ps

# If not running, start it
docker-compose up -d

# Wait for it to be healthy
docker-compose ps  # Look for "healthy" status

# Check logs if issues persist
docker-compose logs postgres
```

### Issue: "database does not exist"

**Cause:** The database wasn't created automatically.

**Solution:**
```bash
# Connect to PostgreSQL
docker exec -it aurora-postgres psql -U aurora

# Create database manually
CREATE DATABASE aurora_dev;
\q

# Restart server to run migrations
go run cmd/server/main.go
```

### Issue: Tests failing with "database connection refused"

**Cause:** Test database isn't running.

**Solution:**
```bash
# Integration tests need Docker
docker-compose up -d

# Wait a few seconds
sleep 3

# Run tests
make test-integration
```

### Issue: Port 8080 already in use

**Solution:**
```bash
# Option 1: Find and kill the process
lsof -ti:8080 | xargs kill

# Option 2: Change port in .env
echo "PORT=8081" >> .env
```

### Issue: Hot reload not working

Aurora doesn't have hot reload by default. **Solutions:**

```bash
# Option 1: Use air (recommended)
go install github.com/cosmtrek/air@latest
air  # Watches files and auto-restarts

# Option 2: Use nodemon
npm install -g nodemon
nodemon --exec "go run cmd/server/main.go" --ext go

# Option 3: Manual restart with watch
while true; do go run cmd/server/main.go; done
```

See [docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md) for hot reload setup.

---

## Environment Configuration

Aurora uses environment variables for configuration. Default values work for local development.

**Configuration Files:**
- `.env.example` - Template with all available options
- `.env` - Your local configuration (not committed to Git)

**Variables Explained:**

```env
# Server Configuration
PORT=8080                # HTTP server port
ENV=development          # Environment: development, staging, production

# Database Configuration
DB_HOST=localhost        # Database hostname
DB_PORT=5432            # Database port
DB_USER=aurora          # Database username
DB_PASSWORD=aurora_dev  # Database password
DB_NAME=aurora_dev      # Database name
DB_SSLMODE=disable      # SSL mode: disable (dev), require (prod)

# Logging
LOG_LEVEL=info          # Log level: debug, info, warn, error
```

**Setup:**
```bash
# Copy example to create your config
cp .env.example .env

# Edit as needed
vim .env

# Verify configuration
go run cmd/server/main.go  # Should start without errors
```

**Production Configuration:**

For production, use secure values:
```env
ENV=production
DB_SSLMODE=require
DB_PASSWORD=<strong-random-password>
LOG_LEVEL=warn
```

## Next Steps: Future Phases

### Phase 2: Code Generation
- Automatic SDK generation for TypeScript/Python
- Client library generation with type safety
- API client customization options

### Phase 3: Advanced Features
- Real-time subscriptions support
- Authentication and authorization
- Rate limiting and quotas
- Advanced schema migrations

### Phase 4: Operations & Scaling
- Multi-region deployment
- High availability setup
- Performance optimization
- Monitoring and observability

### Phase 5: Enterprise Features
- Team collaboration
- Role-based access control
- Audit logging
- Enterprise SLA support

---

## How to Contribute

We welcome contributions! Here's how to get started:

1. **Read the guidelines:** See [CONTRIBUTING.md](./CONTRIBUTING.md) for code style and process
2. **Pick an issue:** Check GitHub issues or create a new one
3. **Create a branch:** `git checkout -b feature/your-feature-name`
4. **Make changes:** Follow the [Development Workflow](#development-workflow)
5. **Write tests:** Ensure test coverage for new code
6. **Run tests:** `make test && make test-integration`
7. **Submit PR:** Create a pull request with a clear description

**Good First Issues:**
- Add input validation for project names (no special characters)
- Add pagination to project list endpoint
- Add database health check to `/health` endpoint
- Improve error messages for schema validation
- Add OpenAPI/Swagger documentation

**Questions?**
- Open a GitHub Discussion
- Check [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) for design context
- See [docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md) for development tips

---

## Additional Resources

- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Code style, PR process, commit conventions
- **[docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)** - System design, database schema, design decisions
- **[docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md)** - Local setup, debugging, common tasks
- **[API Documentation](#api-endpoints)** - Available endpoints (see above)
- **[Schema Format](#schema-definition-format)** - YAML schema syntax (see above)

---

## License

Aurora is provided as an example project. Use and modify as needed.
