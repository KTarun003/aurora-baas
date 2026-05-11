# Aurora Development Guide

This guide covers day-to-day development tasks, debugging, and common workflows for Aurora contributors.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Running the Application](#running-the-application)
- [Database Management](#database-management)
- [Debugging](#debugging)
- [Testing](#testing)
- [Common Development Tasks](#common-development-tasks)
- [Troubleshooting](#troubleshooting)
- [Performance Tips](#performance-tips)

---

## Development Environment Setup

### Prerequisites

Ensure you have:
- Go 1.26.2+
- Docker Desktop (or Docker Engine + Docker Compose)
- PostgreSQL client tools (optional, for manual DB access)
- Your favorite IDE (VS Code, GoLand, etc.)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/ktarun.reddy/baas.git
cd baas

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Start PostgreSQL
docker-compose up -d

# Verify database is running
docker-compose ps

# Run migrations and start server
go run cmd/server/main.go
```

### IDE Setup

#### Visual Studio Code

**Recommended Extensions:**

```json
{
  "recommendations": [
    "golang.go",              // Go language support
    "ms-azuretools.vscode-docker", // Docker support
    "redhat.vscode-yaml",     // YAML syntax
    "eamodio.gitlens",        // Git integration
    "humao.rest-client"       // Test API endpoints
  ]
}
```

**Settings** (`.vscode/settings.json`):

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "gofmt",
  "go.testFlags": ["-v"],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

**Debug Configuration** (`.vscode/launch.json`):

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/server/main.go",
      "env": {
        "PORT": "8080",
        "ENV": "development"
      },
      "args": []
    },
    {
      "name": "Debug Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}",
      "args": ["-test.run", "^TestName$"]
    }
  ]
}
```

#### GoLand / IntelliJ IDEA

1. **Open Project:** File → Open → Select `baas` directory
2. **Enable Go Modules:** Settings → Go → Go Modules → Enable
3. **Run Configuration:**
   - Run → Edit Configurations
   - Add → Go Build
   - Files: `cmd/server/main.go`
   - Working directory: Project root
   - Environment: Load from `.env`

### Git Hooks Setup

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Run tests
go test ./...

if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Running the Application

### Standard Start

```bash
# Start database
docker-compose up -d

# Run server
go run cmd/server/main.go
```

### With Hot Reload

**Option 1: Using Air** (Recommended)

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
cat > .air.toml << 'EOF'
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server/main.go"
bin = "tmp/main"
include_ext = ["go", "yaml"]
exclude_dir = ["tmp", "vendor", "tests"]
delay = 1000

[log]
time = true
EOF

# Run with hot reload
air
```

**Option 2: Using Nodemon**

```bash
# Install nodemon globally
npm install -g nodemon

# Run
nodemon --exec "go run cmd/server/main.go" --ext go,yaml
```

**Option 3: Using Watch + Loop**

```bash
# Install fswatch (macOS)
brew install fswatch

# Watch and restart
fswatch -o . | xargs -n1 -I{} go run cmd/server/main.go
```

### Build and Run Binary

```bash
# Build
make build

# Run
./bin/baas

# Build for Linux (for deployment)
GOOS=linux GOARCH=amd64 go build -o bin/baas-linux cmd/server/main.go
```

### Running with Different Environments

```bash
# Development (default)
ENV=development go run cmd/server/main.go

# Production mode
ENV=production PORT=3000 DB_SSLMODE=require go run cmd/server/main.go

# Custom database
DB_HOST=prod-db.example.com DB_USER=produser go run cmd/server/main.go
```

---

## Database Management

### Starting/Stopping Database

```bash
# Start PostgreSQL
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs postgres
docker-compose logs -f postgres  # Follow mode

# Stop database
docker-compose down

# Stop and remove data
docker-compose down -v
```

### Accessing Database

**Using psql (command-line):**

```bash
# Connect to database
docker exec -it aurora-postgres psql -U aurora -d aurora_dev

# Common commands
\dt              # List tables
\d projects      # Describe projects table
\d schemas       # Describe schemas table
SELECT * FROM projects;
\q               # Quit
```

**Using GUI clients:**

Connection details:
- Host: `localhost`
- Port: `5432`
- Database: `aurora_dev`
- Username: `aurora`
- Password: `aurora_dev`

**Recommended GUI clients:**
- **pgAdmin** (Free, feature-rich)
- **DBeaver** (Free, multi-database)
- **TablePlus** (Paid, beautiful UI)
- **Postico** (macOS, simple)

### Manual Queries

```sql
-- View all projects
SELECT id, name, language, database_type, created_at 
FROM projects 
ORDER BY created_at DESC;

-- View schemas for a project
SELECT version, created_at 
FROM schemas 
WHERE project_id = 'your-project-id' 
ORDER BY version DESC;

-- Get latest schema version
SELECT MAX(version) as latest_version 
FROM schemas 
WHERE project_id = 'your-project-id';

-- Count projects by language
SELECT language, COUNT(*) 
FROM projects 
GROUP BY language;

-- Find projects with no schemas
SELECT p.id, p.name 
FROM projects p 
LEFT JOIN schemas s ON p.id = s.project_id 
WHERE s.id IS NULL;
```

### Database Migrations

Aurora uses GORM Auto-Migrate, which runs on server startup.

**How it works:**

```go
// internal/database/migrations.go
func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &domain.Project{},
        &domain.Schema{},
    )
}
```

**Adding a new field to Project:**

1. Edit `internal/domain/project.go`:
   ```go
   type Project struct {
       // ... existing fields
       Tags []string `gorm:"type:text[]" json:"tags"`  // New field
   }
   ```

2. Restart server:
   ```bash
   go run cmd/server/main.go
   ```

3. GORM automatically adds the new column!

**Manual migration (if needed):**

```bash
docker exec -it aurora-postgres psql -U aurora -d aurora_dev

ALTER TABLE projects ADD COLUMN tags text[];
```

### Resetting the Database

```bash
# Option 1: Drop and recreate database
docker exec -it aurora-postgres psql -U aurora << EOF
DROP DATABASE aurora_dev;
CREATE DATABASE aurora_dev;
EOF

# Option 2: Delete all data (keeps schema)
docker exec -it aurora-postgres psql -U aurora -d aurora_dev << EOF
TRUNCATE TABLE schemas CASCADE;
TRUNCATE TABLE projects CASCADE;
EOF

# Option 3: Nuclear option (removes all data and containers)
docker-compose down -v
docker-compose up -d
```

### Database Backups

**Backup:**

```bash
# Dump database to file
docker exec aurora-postgres pg_dump -U aurora aurora_dev > backup.sql

# Dump with compression
docker exec aurora-postgres pg_dump -U aurora aurora_dev | gzip > backup.sql.gz
```

**Restore:**

```bash
# Restore from backup
docker exec -i aurora-postgres psql -U aurora -d aurora_dev < backup.sql

# Restore from compressed backup
gunzip -c backup.sql.gz | docker exec -i aurora-postgres psql -U aurora -d aurora_dev
```

---

## Debugging

### Logging

**Current logging:**

```go
// Middleware logs all requests
// internal/api/middleware/logger.go
log.Printf("[%s] %s %s %d %v", method, path, clientIP, statusCode, latency)
```

**Add debug logging to your code:**

```go
import "log"

// In any function
log.Printf("Debug: project ID is %s", projectID)
log.Printf("Debug: schema content: %s", content)
```

**Environment-based logging:**

```go
import "os"

if os.Getenv("ENV") == "development" {
    log.Printf("Debug info: %v", data)
}
```

**Better logging with levels (future):**

```go
// Using zerolog or zap
import "github.com/rs/zerolog/log"

log.Debug().Str("project_id", id).Msg("Fetching project")
log.Info().Msg("Server started")
log.Error().Err(err).Msg("Failed to connect to database")
```

### Debugging with Delve

**Install Delve:**

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Debug server:**

```bash
# Start debugger
dlv debug cmd/server/main.go

# Set breakpoint
break internal/service/project_service.go:26

# Run
continue

# Inspect variables
print project
print err

# Step through
next
step
```

**Debug tests:**

```bash
# Debug specific test
dlv test ./internal/service -- -test.run TestCreateProject
```

### Debugging HTTP Requests

**Using curl:**

```bash
# Create project (verbose)
curl -v -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "language": "typescript",
    "database_type": "postgres",
    "api_style": "rest"
  }'

# Get project
curl -v http://localhost:8080/api/v1/projects/PROJECT_ID

# Apply schema
curl -v -X POST http://localhost:8080/api/v1/projects/PROJECT_ID/schemas \
  -H "Content-Type: application/json" \
  -d @examples/schema.json
```

**Using HTTPie (better formatting):**

```bash
# Install
brew install httpie

# Create project
http POST localhost:8080/api/v1/projects \
  name="Test" \
  language="typescript" \
  database_type="postgres" \
  api_style="rest"

# Get project
http GET localhost:8080/api/v1/projects/PROJECT_ID
```

**Using REST Client (VS Code):**

Create `api-tests.http`:

```http
### Health Check
GET http://localhost:8080/health

### Create Project
POST http://localhost:8080/api/v1/projects
Content-Type: application/json

{
  "name": "My Test Project",
  "language": "typescript",
  "database_type": "postgres",
  "api_style": "rest"
}

### Get Project (replace PROJECT_ID)
GET http://localhost:8080/api/v1/projects/PROJECT_ID

### Apply Schema
POST http://localhost:8080/api/v1/projects/PROJECT_ID/schemas
Content-Type: application/json

{
  "content": "version: '1.0'\ndatabase: 'postgres'\ntables:\n  - name: 'users'\n    columns:\n      - name: 'id'\n        type: 'uuid'\n        primary_key: true"
}
```

### Common Debugging Scenarios

**Issue: Handler returns 500 error**

1. Check server logs for error message
2. Add logging to service layer
3. Check if validation is failing
4. Verify database connectivity

**Issue: Database query returns empty result**

```go
// Add debug logging
func (r *ProjectRepository) FindByID(id string) (*domain.Project, error) {
    var project domain.Project
    log.Printf("Searching for project with ID: %s", id)
    
    result := r.db.First(&project, "id = ?", id)
    log.Printf("Query result: %+v", project)
    log.Printf("Error: %v", result.Error)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return &project, nil
}
```

**Issue: JSON parsing fails**

```go
// Add debug logging
func (h *ProjectHandler) Create(c *gin.Context) {
    var req CreateProjectRequest
    
    // Log raw body
    bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
    log.Printf("Raw body: %s", string(bodyBytes))
    c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
    
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("JSON parse error: %v", err)
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // ...
}
```

---

## Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/service/...

# Specific test
go test -run TestCreateProject ./internal/service/

# With verbose output
go test -v ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Integration tests
make test-integration

# Race detection
go test -race ./...

# Parallel tests
go test -parallel 4 ./...
```

### Writing Unit Tests

**Example: Testing service layer**

```go
// internal/service/project_service_test.go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/ktarun.reddy/baas/internal/domain"
    "github.com/ktarun.reddy/baas/internal/service"
)

// Mock repository
type mockProjectRepository struct {
    createFunc func(*domain.Project) error
    findFunc   func(string) (*domain.Project, error)
}

func (m *mockProjectRepository) Create(p *domain.Project) error {
    if m.createFunc != nil {
        return m.createFunc(p)
    }
    return nil
}

func (m *mockProjectRepository) FindByID(id string) (*domain.Project, error) {
    if m.findFunc != nil {
        return m.findFunc(id)
    }
    return &domain.Project{ID: id}, nil
}

// Test
func TestCreateProject_Success(t *testing.T) {
    // Arrange
    repo := &mockProjectRepository{}
    svc := service.NewProjectService(repo)
    project := &domain.Project{
        Name:         "Test",
        Language:     "typescript",
        DatabaseType: "postgres",
        APIStyle:     "rest",
    }
    
    // Act
    err := svc.CreateProject(project)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, project.ID)
}

func TestCreateProject_InvalidLanguage(t *testing.T) {
    // Arrange
    repo := &mockProjectRepository{}
    svc := service.NewProjectService(repo)
    project := &domain.Project{
        Name:         "Test",
        Language:     "invalid",
        DatabaseType: "postgres",
        APIStyle:     "rest",
    }
    
    // Act
    err := svc.CreateProject(project)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "language")
}
```

### Writing Integration Tests

**Example: Testing API endpoints**

```go
// tests/integration/project_test.go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/ktarun.reddy/baas/internal/api"
    // ... imports
)

func setupTestRouter(t *testing.T) *gin.Engine {
    // Setup test database
    db := setupTestDB(t)
    
    // Initialize dependencies
    projectRepo := repository.NewProjectRepository(db)
    schemaRepo := repository.NewSchemaRepository(db)
    projectService := service.NewProjectService(projectRepo)
    schemaService := service.NewSchemaService(schemaRepo, projectRepo)
    
    // Create router
    return api.NewRouter(projectService, schemaService)
}

func TestCreateProject_EndToEnd(t *testing.T) {
    // Arrange
    router := setupTestRouter(t)
    body := map[string]string{
        "name":          "Test Project",
        "language":      "typescript",
        "database_type": "postgres",
        "api_style":     "rest",
    }
    bodyJSON, _ := json.Marshal(body)
    
    req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(bodyJSON))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // Act
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotEmpty(t, response["id"])
    assert.Equal(t, "Test Project", response["name"])
}
```

### Test Coverage Reports

```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View in browser (HTML)
go tool cover -html=coverage.out

# Coverage by package
go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out | grep -E "^total:"
```

**Target coverage:**
- Domain: 100%
- Service: 90%+
- Repository: 80%+
- API handlers: 80%+

---

## Common Development Tasks

### Adding a New API Endpoint

**Example: Add GET /api/v1/projects/:id/stats endpoint**

**Step 1: Define the handler**

```go
// internal/api/handlers/project.go

func (h *ProjectHandler) GetStats(c *gin.Context) {
    projectID := c.Param("id")
    
    stats, err := h.service.GetProjectStats(projectID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, stats)
}
```

**Step 2: Add the route**

```go
// internal/api/router.go

projectsGroup.GET("/:id/stats", projectHandler.GetStats)
```

**Step 3: Implement service method**

```go
// internal/service/project_service.go

type ProjectStats struct {
    SchemaCount   int `json:"schema_count"`
    LatestVersion int `json:"latest_version"`
}

func (s *ProjectService) GetProjectStats(projectID string) (*ProjectStats, error) {
    // Implementation
    count := s.schemaRepo.CountByProjectID(projectID)
    latest := s.schemaRepo.GetLatestVersion(projectID)
    
    return &ProjectStats{
        SchemaCount:   count,
        LatestVersion: latest,
    }, nil
}
```

**Step 4: Add tests**

```go
// internal/service/project_service_test.go

func TestGetProjectStats(t *testing.T) {
    // Test implementation
}
```

**Step 5: Test manually**

```bash
curl http://localhost:8080/api/v1/projects/PROJECT_ID/stats
```

### Adding a New Domain Entity

**Example: Add "Environment" entity (dev, staging, prod per project)**

**Step 1: Define the entity**

```go
// internal/domain/environment.go

package domain

import "time"

type Environment struct {
    ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    ProjectID string    `gorm:"type:uuid;not null;index" json:"project_id"`
    Name      string    `gorm:"type:varchar(50);not null" json:"name"`
    URL       string    `gorm:"type:varchar(255)" json:"url"`
    CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
    UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (Environment) TableName() string {
    return "environments"
}

func (e *Environment) Validate() error {
    if e.Name == "" {
        return errors.New("name is required")
    }
    return nil
}
```

**Step 2: Add migration**

```go
// internal/database/migrations.go

func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &domain.Project{},
        &domain.Schema{},
        &domain.Environment{},  // Add this
    )
}
```

**Step 3: Create repository**

```go
// internal/repository/environment_repo.go

package repository

type EnvironmentRepository struct {
    db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
    return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) Create(env *domain.Environment) error {
    return r.db.Create(env).Error
}

// ... other CRUD methods
```

**Step 4: Create service**

```go
// internal/service/environment_service.go

package service

type EnvironmentService struct {
    repo    *repository.EnvironmentRepository
    projRepo *repository.ProjectRepository
}

func NewEnvironmentService(repo *repository.EnvironmentRepository, projRepo *repository.ProjectRepository) *EnvironmentService {
    return &EnvironmentService{repo: repo, projRepo: projRepo}
}

func (s *EnvironmentService) CreateEnvironment(env *domain.Environment) error {
    // Validation
    if err := env.Validate(); err != nil {
        return err
    }
    
    // Verify project exists
    if _, err := s.projRepo.FindByID(env.ProjectID); err != nil {
        return errors.New("project not found")
    }
    
    return s.repo.Create(env)
}
```

**Step 5: Create handler and routes**

```go
// internal/api/handlers/environment.go
// internal/api/router.go - add routes
```

**Step 6: Restart server to run migrations**

```bash
go run cmd/server/main.go
```

### Adding Configuration Options

**Step 1: Add to Config struct**

```go
// internal/config/config.go

type ServerConfig struct {
    Port          string
    Env           string
    RequestTimeout int  // New option
}
```

**Step 2: Add environment variable**

```go
// internal/config/config.go - in Load()

server := ServerConfig{
    Port:          getEnv("PORT", "8080"),
    Env:           getEnv("ENV", "development"),
    RequestTimeout: getEnvInt("REQUEST_TIMEOUT", 30),  // New
}
```

**Step 3: Update .env.example**

```bash
# Add to .env.example
REQUEST_TIMEOUT=30
```

**Step 4: Use in code**

```go
// cmd/server/main.go or middleware

timeout := time.Duration(cfg.Server.RequestTimeout) * time.Second
```

---

## Troubleshooting

### Server Won't Start

**Check 1: Is PostgreSQL running?**

```bash
docker-compose ps
# Should show aurora-postgres as "Up"
```

**Check 2: Is port available?**

```bash
lsof -i :8080
# If something is using port 8080, kill it or change PORT in .env
```

**Check 3: Can server connect to database?**

```bash
docker exec -it aurora-postgres psql -U aurora -d aurora_dev -c "SELECT 1"
# Should return "1"
```

**Check 4: Check logs**

```bash
# Run server with verbose logging
go run cmd/server/main.go 2>&1 | tee server.log
```

### Tests Failing

**Issue: "database connection refused"**

```bash
# Start test database
docker-compose up -d

# Wait for it to be ready
sleep 3

# Run tests
make test-integration
```

**Issue: "table doesn't exist"**

```bash
# Migrations might not have run
# Restart server to trigger migrations
go run cmd/server/main.go
# (migrations run on startup)
```

**Issue: Tests pass individually but fail in parallel**

```bash
# Database state conflicts - use transactions in tests
# Or run sequentially:
go test -parallel 1 ./...
```

### Performance Issues

**Issue: API responses are slow**

```bash
# Check database query performance
# Enable GORM logging:
```

```go
// internal/database/postgres.go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),  // Enable query logging
})
```

**Issue: High memory usage**

```bash
# Profile memory
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# Or run server with profiling
go run cmd/server/main.go &
curl localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

---

## Performance Tips

### Database Optimization

**Add indexes for common queries:**

```go
// internal/domain/project.go
type Project struct {
    // ...
    Language string `gorm:"type:varchar(50);not null;index"` // Add index
}
```

**Use connection pooling:**

```go
// internal/database/postgres.go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Batch operations:**

```go
// Instead of:
for _, project := range projects {
    repo.Create(project)
}

// Use:
repo.db.Create(&projects)  // Batch insert
```

### API Optimization

**Enable GZIP compression:**

```go
// internal/api/router.go
import "github.com/gin-contrib/gzip"

router.Use(gzip.Gzip(gzip.DefaultCompression))
```

**Add pagination:**

```go
func (h *ProjectHandler) List(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "20")
    
    projects, total := h.service.ListProjects(page, pageSize)
    
    c.JSON(200, gin.H{
        "data":  projects,
        "total": total,
        "page":  page,
    })
}
```

### Code Optimization

**Use context for timeouts:**

```go
ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
defer cancel()

result, err := repo.FindWithContext(ctx, id)
```

**Profile CPU usage:**

```bash
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

---

## Next Steps

- Read [ARCHITECTURE.md](./ARCHITECTURE.md) for system design
- Read [CONTRIBUTING.md](../CONTRIBUTING.md) for code style
- Explore the codebase starting from `cmd/server/main.go`
- Join development discussions on GitHub

Happy coding!
