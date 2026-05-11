# Aurora Architecture Deep Dive

This document provides a comprehensive overview of Aurora's architecture, design decisions, and system internals.

## Table of Contents

- [System Overview](#system-overview)
- [Architecture Principles](#architecture-principles)
- [Layer Architecture](#layer-architecture)
- [Request Flow](#request-flow)
- [Database Schema](#database-schema)
- [Key Components](#key-components)
- [Design Decisions](#design-decisions)
- [Extension Points](#extension-points)
- [Future Architecture](#future-architecture)

---

## System Overview

Aurora is a Backend as a Service (BaaS) platform designed to simplify backend development by providing:
- Dynamic project creation with configurable languages, databases, and API styles
- Schema definition and versioning through YAML
- Automatic database migrations
- RESTful API for management operations

**High-Level Architecture:**

```
┌───────────────────────────────────────────────────────┐
│                    Client Layer                       │
│           (HTTP Clients, Web UI, CLI)                 │
└───────────────────┬───────────────────────────────────┘
                    │ HTTP/REST
                    ▼
┌───────────────────────────────────────────────────────┐
│                  API Gateway Layer                    │
│         (Gin Router + Middleware)                     │
│  ┌──────────┬──────────┬──────────┬──────────────┐   │
│  │ Logger   │  Auth    │  CORS    │  Rate Limit  │   │
│  │  (✓)     │ (Future) │ (Future) │   (Future)   │   │
│  └──────────┴──────────┴──────────┴──────────────┘   │
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────────┐
│                   Handler Layer                       │
│           (HTTP Request/Response)                     │
│  ┌──────────────┬──────────────┬──────────────────┐  │
│  │   Health     │   Project    │     Schema       │  │
│  │   Handler    │   Handler    │     Handler      │  │
│  └──────────────┴──────────────┴──────────────────┘  │
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────────┐
│                  Service Layer                        │
│              (Business Logic)                         │
│  ┌──────────────────────┬──────────────────────────┐ │
│  │  Project Service     │    Schema Service        │ │
│  │  - Validation        │    - Validation          │ │
│  │  - Business Rules    │    - Versioning Logic    │ │
│  │                      │    - YAML Processing     │ │
│  └──────────────────────┴──────────────────────────┘ │
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────────┐
│                Repository Layer                       │
│               (Data Access)                           │
│  ┌──────────────────────┬──────────────────────────┐ │
│  │  Project Repository  │   Schema Repository      │ │
│  │  - CRUD Operations   │   - CRUD Operations      │ │
│  │  - Query Building    │   - Version Queries      │ │
│  └──────────────────────┴──────────────────────────┘ │
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────────┐
│                 Database Layer                        │
│            (GORM + PostgreSQL)                        │
│  ┌────────────────┬───────────────┬────────────────┐ │
│  │  Connection    │   Migrations  │   Connection   │ │
│  │   Pooling      │   Management  │     Pool       │ │
│  └────────────────┴───────────────┴────────────────┘ │
└───────────────────┬───────────────────────────────────┘
                    │
                    ▼
            ┌───────────────┐
            │  PostgreSQL   │
            │   Database    │
            └───────────────┘
```

---

## Architecture Principles

Aurora follows these core principles:

### 1. Clean Architecture

**Separation of Concerns:** Each layer has a single responsibility:
- **API Layer:** HTTP protocol handling only
- **Service Layer:** Business logic and rules
- **Repository Layer:** Data persistence
- **Domain Layer:** Core business entities

**Dependency Rule:** Dependencies point inward:
```
API → Service → Repository → Domain
```

Domain has no dependencies on outer layers. This makes the domain logic independent of frameworks, databases, and UI.

**Why?**
- **Testability:** Each layer can be tested independently
- **Maintainability:** Changes in one layer don't cascade
- **Flexibility:** Easy to swap implementations (e.g., switch from PostgreSQL to MongoDB)

### 2. Domain-Driven Design (DDD)

**Core Concepts:**
- **Entities:** `Project` and `Schema` are our core entities
- **Value Objects:** Could add for things like `Email`, `APIKey` (future)
- **Repositories:** Abstract data storage
- **Services:** Coordinate business operations

**Why?**
- Business logic lives in the domain, not scattered across layers
- Code reflects business concepts (Project, Schema, Version)
- Easy to discuss with non-technical stakeholders

### 3. SOLID Principles

**Single Responsibility:** Each struct/function has one job
```go
// ProjectService handles project business logic
// SchemaService handles schema business logic
// NOT: ProjectAndSchemaService
```

**Open/Closed:** Open for extension, closed for modification
```go
// Easy to add new database types by implementing the interface
type Repository interface {
    Create(entity interface{}) error
    FindByID(id string) (interface{}, error)
}
```

**Liskov Substitution:** Interfaces define contracts
```go
// Any Repository implementation can be used
service := NewProjectService(postgresRepo)  // or mongoRepo, or redisRepo
```

**Interface Segregation:** Small, focused interfaces
```go
// Not: OneGiantRepository interface
// Instead: Multiple small interfaces
type ProjectFinder interface {
    FindByID(id string) (*Project, error)
}
type ProjectCreator interface {
    Create(project *Project) error
}
```

**Dependency Inversion:** Depend on abstractions, not concrete types
```go
// Service depends on Repository interface, not PostgresRepository
type ProjectService struct {
    repo ProjectRepositoryInterface  // Not *PostgresProjectRepository
}
```

### 4. Convention Over Configuration

**Sensible Defaults:**
- PostgreSQL on localhost:5432
- Server on port 8080
- Automatic migrations on startup

**Why?**
- New developers get running quickly
- Less configuration to maintain
- Easy to override when needed

---

## Layer Architecture

### API Layer (`internal/api/`)

**Responsibilities:**
- Parse HTTP requests
- Validate request format (not business rules)
- Call appropriate service methods
- Format HTTP responses
- Handle HTTP errors

**Components:**

1. **Router** (`router.go`):
   ```go
   // Sets up routes and middleware
   router.GET("/health", healthHandler.Check)
   router.POST("/api/v1/projects", projectHandler.Create)
   ```

2. **Handlers** (`handlers/`):
   - `health.go` - Health check endpoint
   - `project.go` - Project CRUD endpoints
   - `schema.go` - Schema management endpoints

3. **Middleware** (`middleware/`):
   - `logger.go` - Request/response logging
   - Future: `auth.go`, `rate_limiter.go`, `cors.go`

**Design Pattern:** Controller pattern
- Each handler is a "controller" for a resource
- Thin layer - no business logic
- Delegates to services

**Example Flow:**
```go
// Handler receives request
func (h *ProjectHandler) Create(c *gin.Context) {
    // 1. Parse JSON body
    var req CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 2. Convert to domain entity
    project := &domain.Project{
        Name:     req.Name,
        Language: req.Language,
        // ...
    }
    
    // 3. Call service
    if err := h.service.CreateProject(project); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 4. Return response
    c.JSON(201, project)
}
```

### Service Layer (`internal/service/`)

**Responsibilities:**
- Implement business logic
- Validate business rules
- Coordinate multiple repositories
- Handle transactions (future)
- Emit events (future)

**Components:**

1. **ProjectService** (`project_service.go`):
   - Validates project configuration
   - Enforces business rules (valid language, database, API style)
   - Coordinates project CRUD operations

2. **SchemaService** (`schema_service.go`):
   - Validates schema YAML
   - Manages schema versioning
   - Coordinates schema storage
   - Future: Generate migrations

**Design Patterns:**
- **Service pattern:** Each service encapsulates business logic for a domain entity
- **Facade pattern:** Services provide simple interface to complex operations

**Example Logic:**
```go
func (s *SchemaService) ApplySchema(projectID, content string) (*domain.Schema, error) {
    // 1. Verify project exists
    project, err := s.projectRepo.FindByID(projectID)
    if err != nil {
        return nil, err
    }
    
    // 2. Validate schema YAML
    if err := validator.ValidateSchema(content); err != nil {
        return nil, err
    }
    
    // 3. Check database type matches
    schemaDB := extractDatabaseType(content)
    if schemaDB != project.DatabaseType {
        return nil, errors.New("schema database doesn't match project")
    }
    
    // 4. Determine next version
    latestVersion, _ := s.schemaRepo.GetLatestVersion(projectID)
    nextVersion := latestVersion + 1
    
    // 5. Create schema
    schema := &domain.Schema{
        ProjectID: projectID,
        Content:   content,
        Version:   nextVersion,
    }
    
    // 6. Save to repository
    if err := s.schemaRepo.Create(schema); err != nil {
        return nil, err
    }
    
    return schema, nil
}
```

**Why This Layer is Important:**
- Business rules are isolated and testable
- No HTTP details leak in (can use this service from CLI, gRPC, etc.)
- Easy to add complex workflows (multi-step transactions, events)

### Repository Layer (`internal/repository/`)

**Responsibilities:**
- Data persistence operations (CRUD)
- Database query building
- Transaction management
- Connection pooling (via GORM)

**Components:**

1. **ProjectRepository** (`project_repo.go`):
   ```go
   type ProjectRepository struct {
       db *gorm.DB
   }
   
   func (r *ProjectRepository) Create(project *domain.Project) error {
       return r.db.Create(project).Error
   }
   
   func (r *ProjectRepository) FindByID(id string) (*domain.Project, error) {
       var project domain.Project
       if err := r.db.First(&project, "id = ?", id).Error; err != nil {
           return nil, err
       }
       return &project, nil
   }
   ```

2. **SchemaRepository** (`schema_repo.go`):
   - Adds versioning queries
   - `GetLatestVersion(projectID)` - Get highest version number
   - `FindByProjectAndVersion(projectID, version)` - Specific version

**Design Pattern:** Repository pattern
- Abstracts data storage mechanism
- Provides collection-like interface
- Hides database implementation details

**Why Use Repositories?**
- Easy to test (mock the repository)
- Easy to change databases (reimplement repository interface)
- Keeps SQL/GORM out of business logic

### Domain Layer (`internal/domain/`)

**Responsibilities:**
- Define core entities
- Entity validation rules
- Business constraints

**Components:**

1. **Project** (`project.go`):
   ```go
   type Project struct {
       ID           string    `gorm:"type:uuid;primary_key"`
       Name         string    `gorm:"type:varchar(255);not null"`
       Language     string    `gorm:"type:varchar(50);not null"`
       DatabaseType string    `gorm:"type:varchar(50);not null"`
       APIStyle     string    `gorm:"type:varchar(50);not null"`
       CreatedAt    time.Time
       UpdatedAt    time.Time
   }
   
   func (p *Project) Validate() error {
       // Domain validation rules
       if !isValidLanguage(p.Language) {
           return errors.New("invalid language")
       }
       return nil
   }
   ```

2. **Schema** (`schema.go`):
   - Versioned schema definition
   - YAML content storage
   - Project relationship

**Design Pattern:** Entity pattern (DDD)
- Entities have identity (ID)
- Entities encapsulate state and behavior
- Validation is part of the entity

**Why Pure Domain?**
- No framework dependencies
- Easy to test
- Represents business concepts directly
- Can be reused in different contexts (API, CLI, jobs)

### Database Layer (`internal/database/`)

**Responsibilities:**
- Database connection management
- Connection pooling
- Schema migrations
- Health checks

**Components:**

1. **Connection** (`postgres.go`):
   ```go
   func Connect(cfg *Config) (*gorm.DB, error) {
       dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
           cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
       
       db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
       if err != nil {
           return nil, err
       }
       
       // Configure connection pool
       sqlDB, _ := db.DB()
       sqlDB.SetMaxIdleConns(10)
       sqlDB.SetMaxOpenConns(100)
       sqlDB.SetConnMaxLifetime(time.Hour)
       
       return db, nil
   }
   ```

2. **Migrations** (`migrations.go`):
   ```go
   func RunMigrations(db *gorm.DB) error {
       // Auto-migrate tables
       return db.AutoMigrate(
           &domain.Project{},
           &domain.Schema{},
       )
   }
   ```

**Why Separate This?**
- Database setup is environment-specific
- Easy to add health checks
- Connection pooling configuration in one place

---

## Request Flow

### Creating a Project

Detailed flow when client calls `POST /api/v1/projects`:

```
1. Client sends HTTP POST request
   POST /api/v1/projects
   Content-Type: application/json
   {
     "name": "My App",
     "language": "typescript",
     "database_type": "postgres",
     "api_style": "rest"
   }
   
   ↓

2. Gin Router receives request
   - Matches route: POST /api/v1/projects
   - Calls middleware stack:
     * Logger middleware logs request
     * (Future: Auth middleware checks token)
     * (Future: Rate limiter checks quota)
   
   ↓

3. ProjectHandler.Create() is called
   - Parses JSON body using Gin binding
   - Validates JSON structure (required fields present)
   - Converts request DTO to domain.Project
   - Calls: h.service.CreateProject(project)
   
   ↓

4. ProjectService.CreateProject() is called
   - Validates business rules:
     * Name not empty
     * Language in [typescript, python]
     * DatabaseType in [postgres, mongodb]
     * APIStyle in [rest, graphql]
   - If validation fails: return error
   - Calls: s.repo.Create(project)
   
   ↓

5. ProjectRepository.Create() is called
   - Builds SQL via GORM
   - Executes: INSERT INTO projects (id, name, ...) VALUES (...)
   - PostgreSQL generates UUID for id
   - Sets created_at, updated_at timestamps
   - Returns project with generated ID
   
   ↓

6. Database persists the project
   - PostgreSQL stores row in projects table
   - Returns success
   
   ↓

7. Success bubbles back up
   Repository → Service → Handler
   
   ↓

8. Handler sends HTTP response
   HTTP 201 Created
   Content-Type: application/json
   {
     "id": "123e4567-e89b-12d3-a456-426614174000",
     "name": "My App",
     "language": "typescript",
     "database_type": "postgres",
     "api_style": "rest",
     "created_at": "2026-05-11T10:30:00Z",
     "updated_at": "2026-05-11T10:30:00Z"
   }
```

**Error Flow:**

If validation fails at any step:

```
1. Error occurs (e.g., invalid language)
   
   ↓

2. Service returns error
   return fmt.Errorf("language must be one of: typescript, python")
   
   ↓

3. Handler catches error
   if err := h.service.CreateProject(project); err != nil {
       c.JSON(400, gin.H{"error": err.Error()})
       return
   }
   
   ↓

4. HTTP error response
   HTTP 400 Bad Request
   {
     "error": "language must be one of: typescript, python"
   }
```

### Applying a Schema

Flow when client calls `POST /api/v1/projects/:id/schemas`:

```
1. Client sends HTTP POST request
   POST /api/v1/projects/abc-123/schemas
   Content-Type: application/json
   {
     "content": "version: '1.0'\ndatabase: 'postgres'\n..."
   }
   
   ↓

2. Router dispatches to SchemaHandler.Apply()
   
   ↓

3. SchemaHandler validates and calls service
   - Extracts project_id from URL
   - Parses request body
   - Calls: h.service.ApplySchema(projectID, content)
   
   ↓

4. SchemaService validates and versions
   - Verifies project exists (calls projectRepo)
   - Validates YAML syntax (calls validator)
   - Extracts database type from YAML
   - Checks database type matches project
   - Gets latest version (calls schemaRepo)
   - Calculates next version (latest + 1)
   - Creates schema entity
   - Calls: s.schemaRepo.Create(schema)
   
   ↓

5. SchemaRepository persists schema
   - Inserts into schemas table
   - Returns schema with generated ID
   
   ↓

6. Handler returns response
   HTTP 201 Created
   {
     "id": "schema-uuid",
     "project_id": "abc-123",
     "content": "...",
     "version": 2,
     "created_at": "...",
     "updated_at": "..."
   }
```

---

## Database Schema

### Projects Table

```sql
CREATE TABLE projects (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    language     VARCHAR(50) NOT NULL,
    database_type VARCHAR(50) NOT NULL,
    api_style    VARCHAR(50) NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_language ON projects(language);
CREATE INDEX idx_projects_database_type ON projects(database_type);
```

**Field Descriptions:**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `id` | UUID | Unique project identifier | Primary key, auto-generated |
| `name` | VARCHAR(255) | Project name | Not null |
| `description` | TEXT | Optional project description | Nullable |
| `language` | VARCHAR(50) | Target language | Not null, enum: typescript, python |
| `database_type` | VARCHAR(50) | Database backend | Not null, enum: postgres, mongodb |
| `api_style` | VARCHAR(50) | API style | Not null, enum: rest, graphql |
| `created_at` | TIMESTAMP | Creation timestamp | Not null, auto-set |
| `updated_at` | TIMESTAMP | Last update timestamp | Not null, auto-updated |

**Design Decisions:**

- **UUID for ID:** Universally unique, can generate client-side, no auto-increment issues in distributed systems
- **VARCHAR for enums:** Simple, readable, easier to add new options than Postgres enums
- **Separate description field:** Keeps name short, description can be long
- **Timestamps:** Essential for auditing and debugging

### Schemas Table

```sql
CREATE TABLE schemas (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    version    INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    
    -- Unique constraint: one version per project
    UNIQUE(project_id, version)
);

-- Indexes
CREATE INDEX idx_schemas_project_id ON schemas(project_id);
CREATE INDEX idx_schemas_project_version ON schemas(project_id, version DESC);
```

**Field Descriptions:**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `id` | UUID | Unique schema identifier | Primary key, auto-generated |
| `project_id` | UUID | Associated project | Foreign key to projects(id), cascade delete |
| `content` | TEXT | YAML schema definition | Not null |
| `version` | INTEGER | Schema version number | Not null, starts at 1 |
| `created_at` | TIMESTAMP | Creation timestamp | Not null |
| `updated_at` | TIMESTAMP | Last update timestamp | Not null |

**Design Decisions:**

- **Foreign key with CASCADE:** When project is deleted, schemas are auto-deleted
- **TEXT for content:** YAML can be large, TEXT supports up to 1GB
- **Integer version:** Simple, sequential, easy to understand
- **Unique constraint on (project_id, version):** Prevents duplicate versions
- **Index on (project_id, version DESC):** Optimizes "get latest version" query

### Entity Relationship Diagram

```
┌─────────────────────────────┐
│         projects            │
├─────────────────────────────┤
│ id (PK)         UUID        │
│ name            VARCHAR(255)│
│ description     TEXT        │
│ language        VARCHAR(50) │
│ database_type   VARCHAR(50) │
│ api_style       VARCHAR(50) │
│ created_at      TIMESTAMP   │
│ updated_at      TIMESTAMP   │
└────────────┬────────────────┘
             │
             │ 1:N
             │
             ▼
┌─────────────────────────────┐
│          schemas            │
├─────────────────────────────┤
│ id (PK)         UUID        │
│ project_id (FK) UUID        │
│ content         TEXT        │
│ version         INTEGER     │
│ created_at      TIMESTAMP   │
│ updated_at      TIMESTAMP   │
│ UNIQUE(project_id, version) │
└─────────────────────────────┘
```

**Relationships:**

- **Project → Schema:** One-to-Many
  - One project can have many schemas (versions)
  - Each schema belongs to one project
  - Cascade delete: delete project → delete all its schemas

---

## Key Components

### Schema Validator (`pkg/validator/`)

**Purpose:** Validate YAML schema definitions before storing.

**Validation Rules:**

1. **YAML Syntax:** Must be valid YAML
2. **Required Fields:**
   - `version` - Schema version string
   - `database` - Database type (postgres, mongodb)
   - `tables` - Array of table definitions
3. **Table Validation:**
   - Must have `name`
   - Must have at least one `column`
4. **Column Validation:**
   - Must have `name` and `type`
   - Valid types: string, text, integer, bigint, decimal, float, boolean, date, timestamp, uuid, json
   - Primary key: at most one per table
   - Foreign key: must reference existing table.column

**Example Usage:**

```go
content := `
version: "1.0"
database: "postgres"
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
        primary_key: true
`

if err := validator.ValidateSchema(content); err != nil {
    // Handle validation error
}
```

**Future Enhancements:**
- Validate foreign key references
- Check for circular dependencies
- Validate default values match column types
- Suggest indexes for foreign keys

### Configuration Management (`internal/config/`)

**Purpose:** Load configuration from environment variables.

**Configuration Structure:**

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
}

type ServerConfig struct {
    Port string  // HTTP port (default: 8080)
    Env  string  // Environment: development, staging, production
}

type DatabaseConfig struct {
    Host     string  // Database host
    Port     string  // Database port
    User     string  // Database user
    Password string  // Database password
    Name     string  // Database name
    SSLMode  string  // SSL mode: disable, require
}
```

**Loading Priority:**

1. Environment variables (highest priority)
2. `.env` file
3. Default values (lowest priority)

**Example:**

```go
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}

// Use configuration
fmt.Printf("Server starting on port %s\n", cfg.Server.Port)
```

**Why Environment Variables?**
- 12-factor app principle
- Different configs for different environments
- Secure (no passwords in code)
- Easy for containerized deployments

---

## Design Decisions

### Why PostgreSQL for Phase 1?

**Decision:** Use PostgreSQL as the primary database for the Aurora control plane.

**Rationale:**
- **ACID transactions:** Critical for schema versioning
- **UUID support:** Native UUID type for IDs
- **JSON support:** Future flexibility for metadata
- **Proven at scale:** Battle-tested in production
- **Rich ecosystem:** GORM, pgx, great tooling

**Trade-offs:**
- More complex setup than SQLite
- Requires server/container

**Alternatives Considered:**
- **SQLite:** Simpler, but not suitable for production multi-user
- **MongoDB:** Schemaless, but we need strong consistency for schema management
- **MySQL:** Similar to PostgreSQL, less feature-rich

### Why GORM?

**Decision:** Use GORM as the ORM layer.

**Rationale:**
- **Auto-migrations:** Simplifies schema management
- **Type-safe:** Compile-time checking of queries
- **Productivity:** Less boilerplate than raw SQL
- **Associations:** Easy to define relationships
- **Hooks:** Can add lifecycle callbacks

**Trade-offs:**
- Performance overhead vs raw SQL
- Learning curve for advanced features
- Some complex queries need raw SQL

**Alternatives Considered:**
- **sqlx:** More control, but more boilerplate
- **Raw database/sql:** Maximum control, most boilerplate
- **gorp:** Lighter, but less features

### Why Gin Web Framework?

**Decision:** Use Gin for HTTP routing and middleware.

**Rationale:**
- **Performance:** One of the fastest Go web frameworks
- **Familiar:** Similar to Express.js (helps onboarding)
- **Middleware:** Easy to add logging, auth, etc.
- **JSON binding:** Automatic request parsing
- **Active community:** Well-maintained, good docs

**Trade-offs:**
- Not standard library (dependency)
- Opinionated routing

**Alternatives Considered:**
- **stdlib net/http:** More control, more boilerplate
- **Echo:** Similar to Gin, slightly different API
- **Chi:** Lighter, but less features

### Why Clean Architecture?

**Decision:** Separate API, Service, Repository, and Domain layers.

**Rationale:**
- **Testability:** Each layer independently testable
- **Flexibility:** Easy to change implementations
- **Clarity:** Clear responsibilities
- **Scalability:** Easy to add features

**Trade-offs:**
- More files and interfaces
- Learning curve for new developers
- Potential over-engineering for simple features

**When It Pays Off:**
- Adding new API (gRPC alongside REST)
- Switching databases
- Adding complex business logic
- Writing comprehensive tests

### Why Schema Versioning?

**Decision:** Track schema changes as versions instead of migrations.

**Rationale:**
- **History:** Can view any previous schema
- **Rollback:** Can revert to previous version (future)
- **Audit:** Know exactly when schema changed
- **Compare:** Diff between versions (future)

**Implementation:**
- Store complete schema for each version (not just deltas)
- Version number increments sequentially
- One version per project at a time

**Trade-offs:**
- Storage: duplicates schema content
- Complexity: managing versions

**Future:** Generate migrations by diffing versions.

---

## Extension Points

These are designed extension points for future phases:

### 1. Adding New Database Types

**Current:** PostgreSQL, MongoDB (defined in project creation)

**To Add:** MySQL, SQLite, CockroachDB

**Extension Point:**

```go
// internal/domain/project.go
func (p *Project) Validate() error {
    validDatabases := map[string]bool{
        "postgres": true,
        "mongodb":  true,
        // Add here:
        "mysql":    true,
        "sqlite":   true,
    }
    // ...
}
```

**Future Architecture:**
- Database-specific schema validators
- Database-specific migration generators
- Adapter pattern for different SQL dialects

### 2. Adding Authentication

**Current:** No authentication

**Extension Point:**

```go
// internal/api/middleware/auth.go (new file)
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        // Validate token
        user, err := validateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // Store user in context
        c.Set("user", user)
        c.Next()
    }
}

// internal/api/router.go
router.Use(middleware.AuthMiddleware())
```

**Future:**
- JWT tokens
- API keys
- OAuth integration
- User management service

### 3. Code Generation (Phase 2)

**Extension Point:**

```go
// internal/service/codegen_service.go (new file)
type CodeGenService struct {
    schemaService *SchemaService
    projectService *ProjectService
}

func (s *CodeGenService) GenerateSDK(projectID string) ([]byte, error) {
    // Get project
    project := s.projectService.GetProject(projectID)
    
    // Get latest schema
    schema := s.schemaService.GetLatestSchema(projectID)
    
    // Generate code based on project.Language
    switch project.Language {
    case "typescript":
        return s.generateTypeScript(project, schema)
    case "python":
        return s.generatePython(project, schema)
    }
}
```

**Architecture:**
- Generator interface per language
- Template-based code generation
- Downloadable SDK as ZIP

### 4. Real-Time Updates (Phase 3)

**Extension Point:**

```go
// internal/api/websocket/schema_updates.go (new file)
func (h *SchemaHandler) SubscribeToSchemaChanges(c *gin.Context) {
    projectID := c.Param("id")
    
    // Upgrade to WebSocket
    conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
    
    // Subscribe to schema changes
    updates := h.service.WatchSchemaChanges(projectID)
    
    // Stream updates
    for update := range updates {
        conn.WriteJSON(update)
    }
}
```

**Architecture:**
- WebSocket support
- Pub/sub for schema changes
- Event sourcing pattern

### 5. Multi-Tenancy

**Extension Point:**

```go
// internal/domain/project.go
type Project struct {
    // Add:
    OrganizationID string `gorm:"type:uuid;not null;index"`
    // ...
}

// internal/api/middleware/tenant.go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        orgID := extractOrganizationID(c)
        c.Set("org_id", orgID)
        c.Next()
    }
}
```

**Architecture:**
- Organization entity
- Row-level security
- Separate database per tenant (future)

---

## Future Architecture

### Phase 2: Code Generation

**New Components:**

```
┌─────────────────────────────────────┐
│      Code Generation Service        │
│  ┌────────────┬─────────────────┐   │
│  │ TypeScript │     Python      │   │
│  │ Generator  │    Generator    │   │
│  └────────────┴─────────────────┘   │
└─────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│        Template Engine              │
│  (Go templates or custom)           │
└─────────────────────────────────────┘
```

**Flow:**
1. Get project + latest schema
2. Parse schema YAML
3. Generate code from templates
4. Package as ZIP
5. Return download link

### Phase 3: Advanced Features

**New Components:**

```
┌─────────────────────────────────────┐
│       Authentication Service        │
│  (JWT, API Keys, OAuth)             │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│       Authorization Service         │
│  (RBAC, Resource Permissions)       │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│         Rate Limiting               │
│  (Token bucket, Redis)              │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│       WebSocket Service             │
│  (Real-time updates)                │
└─────────────────────────────────────┘
```

### Phase 4: Operations & Scaling

**Architecture Changes:**

```
                ┌─────────────┐
                │Load Balancer│
                └──────┬──────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼───┐     ┌───▼────┐    ┌───▼────┐
   │Aurora  │     │Aurora  │    │Aurora  │
   │Instance│     │Instance│    │Instance│
   └────┬───┘     └───┬────┘    └───┬────┘
        │             │              │
        └─────────────┼──────────────┘
                      │
            ┌─────────▼──────────┐
            │   PostgreSQL       │
            │  (Read Replicas)   │
            └────────────────────┘
```

**New Components:**
- Distributed caching (Redis)
- Message queue (RabbitMQ/Kafka)
- Monitoring (Prometheus + Grafana)
- Distributed tracing (Jaeger)
- Log aggregation (ELK stack)

---

## Conclusion

Aurora's architecture is designed for:
- **Clarity:** Easy to understand for new developers
- **Flexibility:** Easy to extend and modify
- **Testability:** Each component independently testable
- **Scalability:** Ready for future phases

The clean architecture and clear separation of concerns provide a solid foundation for building a production-grade BaaS platform.

**Next Steps:**
- Read [DEVELOPMENT.md](./DEVELOPMENT.md) for day-to-day development
- Read [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines
- Explore the code starting from `cmd/server/main.go`
