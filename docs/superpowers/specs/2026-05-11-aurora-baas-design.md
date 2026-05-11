# Aurora Backend-as-a-Service Platform - Design Specification

**Date:** 2026-05-11  
**Status:** Draft  
**Version:** 1.0

## Executive Summary

Aurora is a polyglot Backend-as-a-Service platform that automatically generates production-ready microservices from table schemas. Developers define their data models through a Web UI, YAML files, or by importing existing databases, and Aurora generates complete microservices with authentication, API gateway, and deployment pipelines.

**Target Users:** Solo developers and small-to-mid-size teams (1-50 people) who want to accelerate backend development without sacrificing control or customization.

**Deployment Model:** Hybrid SaaS—Aurora control plane runs as a managed service, but generated microservices can deploy to either Aurora's infrastructure or the user's own Kubernetes clusters.

---

## MVP Scope

### Core Capabilities
- **Languages:** TypeScript (Node.js) and Python (FastAPI) - full-featured implementations
- **Databases:** PostgreSQL and MongoDB support
- **API Styles:** REST or GraphQL (developer chooses per service at generation time)
- **Authentication:** JWT, OAuth 2.0/OIDC, Magic Link, Basic Auth, API Keys
- **API Gateway:** All 8 features (routing, auth integration, rate limiting, CORS, logging, caching, transformation, versioning)
- **Deployment:** Container-based (Docker + Kubernetes) with full GitOps (dev/staging/prod branches)
- **Architecture:** One microservice per table with strict service isolation (no shared database access)

### Essential Platform Components (Priority Order)
1. **Code Generation Engine** - Template-based microservice generation
2. **CLI Tool** - Developer experience and local management
3. **Deployment Orchestrator** - K8s deployment automation
4. **Web UI** - Schema designer and project management

---

## System Architecture

### High-Level Overview

Aurora uses a **hybrid architecture**: monolithic core with distributed worker pools for scalability where it matters.

```
┌─────────────────────────────────────┐
│   Aurora Core (Go Monolith)        │
│ • REST API (Gin/Echo)               │
│ • Web UI (React SPA)                │
│ • PostgreSQL (metadata)             │
│ • Job Queue Publisher (Redis)      │
└───────────┬─────────────────────────┘
            ↓ job queue ↓
┌───────────┴─────────┬───────────────┐
↓                     ↓               ↓
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Code Gen     │  │ K8s Deploy   │  │  Build       │
│ Workers (Go) │  │ Workers (Go) │  │ Workers (Go) │
└──────────────┘  └──────────────┘  └──────────────┘
         ↓ generates ↓
┌─────────────────────────────────────┐
│  User's Generated Stack:            │
│  • Microservices (TS/Python)        │
│  • Auth Service                     │
│  • API Gateway (Kong)               │
│  • K8s Manifests                    │
└─────────────────────────────────────┘
```

### Control Plane Components

**1. Aurora Core (Go)**
- Single Go monolith serving both API and Web UI
- REST API using Gin or Echo framework
- React SPA served as static files
- PostgreSQL database for metadata (projects, schemas, users, deployments)
- Redis for job queue (BullMQ-compatible or Go channels)
- Shared code libraries with CLI tool

**2. Worker Pools (Go)**
All workers written in Go for consistency:
- **Code Generator Workers**: Template-based code generation (TypeScript/Python)
- **Deployment Workers**: K8s manifest generation and cluster deployment
- **Build Workers**: Docker image building and registry pushing

Horizontally scalable—workers pick jobs from Redis queue, process independently, report status back to Core.

**3. CLI Tool (Go)**
- Compiled Go binary for cross-platform distribution
- Calls Aurora Core REST API for remote operations
- Manages local dev environment (Docker Compose)
- Shares code/libraries with Aurora Core for consistency

### User's Generated Stack

**Per Project:**
- **Microservices**: One service per table (TypeScript or Python)
- **Auth Service**: Shared authentication service (JWT, OAuth, Magic Link, Basic, API Key)
- **API Gateway**: Kong Gateway with declarative config
- **Git Repository**: Complete source code with extension points
- **K8s Manifests**: Kustomize-based with dev/staging/prod overlays
- **CI/CD Pipelines**: GitHub Actions workflows for GitOps

---

## Schema Definition & Input Methods

Aurora supports three input methods, all converging to a unified internal representation.

### 1. Web UI Schema Builder

Visual drag-and-drop interface:
- Canvas-based table editor
- Column configuration: name, type, constraints (NOT NULL, UNIQUE, DEFAULT)
- Relationship mapper: foreign keys, one-to-many, many-to-many
- Index and constraint configuration
- Real-time validation and error highlighting
- Export to YAML for version control

### 2. Schema Definition Files (YAML)

Declarative schema files checked into version control:

```yaml
version: "1.0"
database: postgres
tables:
  - name: users
    columns:
      - name: id
        type: uuid
        primary_key: true
        default: gen_random_uuid()
      - name: email
        type: string
        unique: true
        nullable: false
      - name: created_at
        type: timestamp
        default: now()
    indexes:
      - columns: [email]
        unique: true

  - name: orders
    columns:
      - name: id
        type: uuid
        primary_key: true
      - name: user_id
        type: uuid
        foreign_key: users.id
        nullable: false
      - name: total
        type: decimal
      - name: status
        type: string
        default: "pending"
```

### 3. Database Import (Introspection)

Connect to existing databases and import schemas:
- PostgreSQL: Use `information_schema` queries
- MySQL: Use `SHOW TABLES` and `DESCRIBE` commands
- MongoDB: Sample documents to infer schema
- User reviews imported schema before generation
- Converts to Aurora's canonical format

### Internal Representation

All three methods convert to Go structs:

```go
type SchemaDefinition struct {
    Version  string  `json:"version"`
    Database string  `json:"database"` // postgres, mysql, mongodb
    Tables   []Table `json:"tables"`
}

type Table struct {
    Name      string      `json:"name"`
    Columns   []Column    `json:"columns"`
    Indexes   []Index     `json:"indexes"`
    Relations []Relation  `json:"relations"`
}
```

This ensures consistency regardless of input method and makes code generation deterministic.

---

## Code Generation Engine

### Template-Based Generation

Workers use Go templates to generate service code. Each language has a structured template directory:

```
templates/
├── typescript/
│   ├── base/
│   │   ├── server.ts.tmpl          # Express/Fastify server
│   │   ├── database.ts.tmpl        # DB connection pool
│   │   ├── router.ts.tmpl          # Route registration
│   │   └── middleware.ts.tmpl      # Auth, logging, errors
│   ├── crud/
│   │   ├── controller.ts.tmpl      # CRUD endpoints (protected)
│   │   ├── service.ts.tmpl         # Business logic layer
│   │   ├── repository.ts.tmpl      # DB queries (protected)
│   │   └── model.ts.tmpl           # Type definitions
│   ├── extensions/
│   │   ├── custom-endpoints.ts     # User adds routes
│   │   ├── hooks.ts                # Pre/post CRUD hooks
│   │   └── middleware.custom.ts    # Custom middleware
│   └── tests/
│       └── service.test.ts.tmpl    # Test scaffolding
└── python/
    └── (similar structure with FastAPI)
```

### Generated File Types

**Protected Files (Auto-regenerated):**
- `controller.generated.ts` - CRUD endpoint handlers
- `repository.generated.ts` - Database query layer
- `model.generated.ts` - TypeScript type definitions
- Header: `// AUTO-GENERATED - DO NOT EDIT`

**Extension Files (Created Once, Never Overwritten):**
- `custom-endpoints.ts` - User-defined routes
- `hooks.ts` - Pre/post CRUD lifecycle hooks
- `middleware.custom.ts` - Custom middleware
- Clear interfaces/hooks for user code

### Extension Points Example

```typescript
// hooks.ts (user edits this file)
export const userHooks = {
  beforeCreate: async (data: UserInput) => {
    // Hash password, validate email, etc.
    data.passwordHash = await bcrypt.hash(data.password, 10);
    return data;
  },
  
  afterCreate: async (user: User) => {
    // Send welcome email, create profile, etc.
    await emailService.sendWelcome(user.email);
  },
  
  beforeUpdate: async (id: string, data: Partial<UserInput>) => {
    // Custom validation, authorization checks
    return data;
  },
  
  afterDelete: async (id: string) => {
    // Cleanup related resources
  }
};
```

### Regeneration Strategy

**Schema Change Flow:**
1. Developer updates schema via UI, YAML, or CLI
2. Aurora Core validates changes
3. Code Gen Worker regenerates ONLY protected files
4. Extension files remain untouched
5. Worker uses git diff to detect conflicts
6. If conflicts, reports to user via CLI: `aurora schema sync --resolve`
7. User reviews changes and resolves manually if needed

**Merge Strategy:**
- Protected files: Complete overwrite (user shouldn't edit them)
- Extension files: Preserved (never regenerated)
- New tables: Generate complete new service
- Deleted tables: Mark service for manual deletion (safety check)

---

## Auth Service Architecture

Aurora auto-generates one Auth service per project, shared across all microservices.

### Auth Service Structure

```
auth-service/
├── providers/
│   ├── jwt.go              # JWT generation/validation
│   ├── oauth.go            # OAuth 2.0 + OIDC (Google, GitHub)
│   ├── magic-link.go       # Email-based passwordless
│   ├── basic.go            # Username/password (bcrypt)
│   └── apikey.go           # API key generation/validation
├── handlers/
│   ├── register.go         # User registration
│   ├── login.go            # Login (routes to provider)
│   ├── refresh.go          # Token refresh
│   └── verify.go           # Token verification endpoint
├── middleware/
│   └── auth.go             # Middleware for other services
└── database/
    └── migrations/         # User table schema
```

### Database Schema

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255),      -- NULL for OAuth/magic link
  provider VARCHAR(50),             -- 'local', 'google', 'github'
  provider_id VARCHAR(255),         -- External provider user ID
  api_keys JSONB DEFAULT '[]',     -- Array of API keys
  roles JSONB DEFAULT '[]',        -- RBAC roles
  metadata JSONB DEFAULT '{}',     -- Custom user data
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_provider ON users(provider, provider_id);
```

### Integration with Generated Services

**Middleware Integration:**
- Each microservice imports auth middleware from auth-service package
- Middleware validates JWT tokens or API keys
- Injects decoded user claims into request context
- Services use claims for authorization decisions

**API Gateway Integration:**
- Gateway validates tokens at the edge (first line of defense)
- Passes decoded claims as headers to services: `X-User-Id`, `X-User-Roles`
- Services can optionally call `auth-service/verify` for extra validation

**Token Flow (JWT Example):**
1. User logs in → Auth service validates credentials → returns JWT + refresh token
2. Client includes JWT in `Authorization: Bearer <token>` header
3. API Gateway validates JWT signature, expiry, and claims
4. If valid, routes to microservice with decoded user info in headers
5. Microservice checks user roles/permissions for authorization
6. Service performs business logic and returns response

### Supported Auth Strategies

| Strategy | Use Case | Implementation |
|----------|----------|----------------|
| JWT | Stateless API auth | HS256/RS256, 1-hour expiry, refresh tokens |
| OAuth 2.0 | Social login (Google, GitHub) | Authorization code flow with PKCE |
| Magic Link | Passwordless email auth | Time-limited tokens sent via email |
| Basic Auth | Internal services, dev | Username/password over HTTPS |
| API Keys | Service-to-service, public APIs | Long-lived tokens with scopes |

---

## API Gateway Architecture

Aurora uses **Kong Gateway** (open-source) for production deployments with a lightweight custom Go gateway as an alternative for simpler self-hosted setups.

### Deployment Structure

Each project gets a dedicated Gateway instance:

```
project-namespace/
├── auth-service
├── users-service
├── orders-service
├── api-gateway (Kong)
└── redis (for rate limiting & caching)
```

### Kong Configuration (Generated per Project)

```yaml
# kong.yaml
_format_version: "3.0"

services:
  - name: users-service
    url: http://users-service:3000
    routes:
      - name: users-api
        paths:
          - /api/v1/users
        methods: [GET, POST, PUT, DELETE]
    plugins:
      # Authentication
      - name: jwt
        config:
          secret_is_base64: false
          key_claim_name: kid
      
      # Rate Limiting
      - name: rate-limiting
        config:
          minute: 100
          hour: 5000
          policy: redis
          redis_host: redis
      
      # CORS
      - name: cors
        config:
          origins: ["*"]
          credentials: true
          methods: [GET, POST, PUT, DELETE, PATCH]
      
      # Request Transformation
      - name: request-transformer
        config:
          add:
            headers:
              - X-User-Id: {{jwt.sub}}
              - X-User-Roles: {{jwt.roles}}
      
      # Logging & Analytics
      - name: prometheus
        config:
          per_consumer: true
      
      # Response Caching
      - name: proxy-cache
        config:
          strategy: memory
          cache_ttl: 300
          content_type: [application/json]
      
      # API Versioning
      - name: request-transformer
        config:
          replace:
            uri: /v1/users:$(uri_captures[1])
```

### Gateway Features Implementation

| Feature | Implementation | Notes |
|---------|----------------|-------|
| **Routing & Load Balancing** | Kong service discovery + K8s DNS | Health checks, circuit breakers |
| **Auth Integration** | JWT/API key validation plugins | Validates before routing to services |
| **Rate Limiting** | Redis-backed per user/IP limits | Configurable per endpoint |
| **CORS** | Kong CORS plugin | Configured per route, supports preflight |
| **Logging & Analytics** | Prometheus metrics + ELK integration | Request/response logs, performance metrics |
| **Response Caching** | Proxy-cache plugin for GET requests | TTL configurable, cache invalidation support |
| **Transformation** | Request/response transformer plugins | Header injection, payload modification |
| **API Versioning** | Path-based versioning (`/api/v1`, `/api/v2`) | Side-by-side version support |

### Alternative: Custom Go Gateway

For users who want simpler deployment without Kong:
- Lightweight Go reverse proxy
- Same middleware features (auth, rate limiting, CORS)
- Redis for shared state
- Easier to customize but fewer advanced features

---

## GitOps & Deployment Pipeline

Full GitOps workflow with three environment branches and automated deployment.

### Git Repository Structure

```
user-project-repo/
├── services/
│   ├── auth-service/
│   │   ├── src/
│   │   ├── tests/
│   │   ├── Dockerfile
│   │   └── package.json
│   ├── users-service/
│   └── orders-service/
├── gateway/
│   └── kong.yaml
├── k8s/
│   ├── base/                    # Base manifests (Kustomize)
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── ingress.yaml
│   │   └── kustomization.yaml
│   └── overlays/
│       ├── dev/
│       │   ├── kustomization.yaml
│       │   └── patches/
│       ├── staging/
│       └── prod/
├── .github/
│   └── workflows/
│       ├── dev.yaml             # Auto-deploy to dev
│       ├── staging.yaml         # Auto-deploy to staging
│       └── prod.yaml            # Auto-deploy to prod
├── aurora.yaml                  # Aurora project config
└── docker-compose.yaml          # Local dev environment
```

### Branch → Environment Mapping

| Branch | Environment | Trigger | Deployment Target |
|--------|-------------|---------|-------------------|
| `dev` | Development | Push | Auto-deploy to dev K8s cluster |
| `staging` | Staging | Push | Auto-deploy to staging K8s cluster |
| `prod` | Production | Push | Auto-deploy to prod K8s cluster |

### CI/CD Workflow (GitHub Actions)

```yaml
# .github/workflows/dev.yaml
name: Deploy to Dev
on:
  push:
    branches: [dev]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker Images
        run: |
          for service in services/*; do
            name=$(basename $service)
            docker build -t registry.aurora.io/$PROJECT_ID/$name:${{ github.sha }} $service
            docker push registry.aurora.io/$PROJECT_ID/$name:${{ github.sha }}
          done
      
      - name: Update K8s Manifests
        run: |
          cd k8s/overlays/dev
          for service in services/*; do
            name=$(basename $service)
            kustomize edit set image $name=registry.aurora.io/$PROJECT_ID/$name:${{ github.sha }}
          done
      
      - name: Trigger Aurora Deployment
        run: |
          curl -X POST https://api.aurora.io/deployments \
            -H "Authorization: Bearer ${{ secrets.AURORA_TOKEN }}" \
            -H "Content-Type: application/json" \
            -d '{
              "project_id": "${{ secrets.PROJECT_ID }}",
              "environment": "dev",
              "commit_sha": "${{ github.sha }}",
              "services": ["auth-service", "users-service", "orders-service"]
            }'
      
      - name: Wait for Deployment
        run: |
          # Poll deployment status
          while true; do
            status=$(curl -s https://api.aurora.io/deployments/$DEPLOYMENT_ID/status)
            if [[ "$status" == "success" ]]; then
              echo "Deployment successful!"
              exit 0
            elif [[ "$status" == "failed" ]]; then
              echo "Deployment failed!"
              exit 1
            fi
            sleep 10
          done
```

### Deployment Worker Process

1. Receives deployment job from Redis queue
2. Clones git repository at specific commit SHA
3. Runs `kustomize build k8s/overlays/{env}` to generate final manifests
4. Applies manifests to target K8s cluster using kubectl
5. Watches rollout status for all deployments
6. Runs post-deployment health checks
7. Reports status back to Aurora Core
8. Sends webhook notification to user (Slack, email, etc.)

### Database Migration Strategy

Each service includes versioned migrations:

```
users-service/
├── migrations/
│   ├── 001_initial_schema.sql
│   ├── 002_add_email_index.sql
│   └── 003_add_roles_column.sql
```

**Migration Execution:**
- Init container runs before service starts
- Uses `golang-migrate` or `Flyway`
- Migrations run sequentially, once per deployment
- Failed migration prevents service from starting
- Rollback support for reversible migrations

---

## Database Architecture & Management

Flexible isolation levels based on environment.

### Database Provisioning Strategy

**Development Environment:**
- **Default:** Shared PostgreSQL/MongoDB cluster, separate database per service
- **Example:** `postgres://cluster:5432/dev_users`, `postgres://cluster:5432/dev_orders`
- **Benefits:** Reduced local resource usage, faster startup
- **Implementation:** Docker Compose with single DB instance

**Staging Environment:**
- **Recommended:** Shared cluster, separate schemas per service
- **Example:** PostgreSQL schemas: `staging_users`, `staging_orders`
- **Benefits:** Cost-effective, logical separation, easier management

**Production Environment:**
- **Recommended:** Dedicated database instances per service
- **Example:** Separate RDS/CloudSQL instances or K8s StatefulSets
- **Benefits:** True isolation, independent scaling, better security
- **Trade-off:** Higher cost, more operational overhead

### Configuration (aurora.yaml)

```yaml
project:
  id: "my-awesome-app"
  name: "My Awesome App"
  
database:
  type: postgres              # postgres, mysql, mongodb, sqlite
  isolation:
    dev: shared-cluster       # One DB, multiple schemas
    staging: shared-cluster
    prod: dedicated-instances
  
  connection_pooling:
    enabled: true
    max_connections: 20
    idle_timeout: 300
  
  backup:
    enabled: true
    schedule: "0 2 * * *"     # Daily at 2 AM
    retention_days: 30
```

### Service-to-Service Communication

**Strict Rule:** Services NEVER share database access. Communication happens via:

**1. REST API Calls (Synchronous)**
```typescript
// orders-service needs user data
const user = await httpClient.get(`http://users-service/api/v1/users/${userId}`);
```

**2. Event Bus (Asynchronous - Optional)**
```typescript
// User created event
eventBus.publish('user.created', { id: user.id, email: user.email });

// Orders service subscribes
eventBus.subscribe('user.created', async (event) => {
  // Create user profile in orders context
});
```

### Credentials Management

- Stored in K8s Secrets
- Injected as environment variables
- Generated services read from env: `DATABASE_URL`
- CLI command to rotate: `aurora db rotate-credentials --env prod`
- Automatic rotation support (optional)

---

## CLI Tool Design

The Aurora CLI is a compiled Go binary providing local development and remote management.

### Command Structure

```bash
# Project Management
aurora init                          # Initialize new Aurora project
aurora project list                  # List all projects
aurora project info <id>             # Show project details
aurora project delete <id>           # Delete project

# Schema Management
aurora schema apply <schema.yaml>    # Apply schema changes
aurora schema sync --resolve         # Sync and resolve conflicts
aurora schema export                 # Export current schema to YAML
aurora schema import --db postgres://...  # Import from database
aurora schema validate <schema.yaml> # Validate schema file

# Development
aurora dev start                     # Start local environment (Docker Compose)
aurora dev stop                      # Stop local environment
aurora dev logs <service>            # Tail service logs
aurora dev test <service>            # Run service tests
aurora dev shell <service>           # Open shell in service container

# Deployment
aurora deploy --env dev              # Deploy to environment
aurora deploy status                 # Check deployment status
aurora deploy list                   # List recent deployments
aurora rollback --env staging        # Rollback to previous version

# Database
aurora db migrate --env prod         # Run pending migrations
aurora db rollback --env prod        # Rollback last migration
aurora db seed --env dev             # Seed with test data
aurora db rotate-credentials --env prod  # Rotate DB credentials
aurora db backup --env prod          # Manual backup

# Auth Management
aurora auth add-user --email user@example.com --role admin
aurora auth list-users
aurora auth generate-apikey --name "mobile-app" --scopes "read,write"
aurora auth revoke-apikey <key-id>

# Monitoring & Logs
aurora logs --env prod --service users --tail 100 --follow
aurora metrics --env prod            # Open metrics dashboard URL
aurora status --env prod             # Show all services status
```

### CLI Architecture

**Components:**
- Single Go binary (cross-compiled for macOS, Linux, Windows)
- Calls Aurora Core REST API for most operations
- Local operations use Docker Compose API
- Configuration stored in:
  - Global: `~/.aurora/config.yaml`
  - Project: `.aurora/project.yaml`

**Developer Experience Features:**
- Interactive prompts for complex operations (with `--yes` flag to skip)
- Progress bars for long-running tasks (code generation, deployment)
- Colored output for better readability
- Clear error messages with suggested fixes
- Auto-completion scripts for bash/zsh/fish
- Built-in help with examples: `aurora help deploy`
- Update notifications: "New version available: v1.2.3"

---

## Web UI Design

React SPA providing visual tools for non-CLI users.

### Key Pages

**1. Dashboard**
- Projects grid with status badges (active, deploying, failed)
- Quick stats per project: service count, last deployment, uptime %
- Recent activity feed (deployments, schema changes, alerts)
- "Create New Project" CTA

**2. Schema Designer**
- Canvas-based visual editor (React Flow)
- Left sidebar: Table palette, properties panel
- Center: Draggable tables with columns visible
- Right sidebar: Relationship visualizer
- Toolbar: Zoom, pan, undo/redo, export YAML, generate services
- Live validation with error highlighting
- Auto-layout button for clean arrangement

**3. Project Settings**
- **General:** Name, description, git repository URL
- **Database:** Type (Postgres/MySQL/MongoDB), isolation levels per env
- **Code Generation:** Language (TypeScript/Python), API style (REST/GraphQL)
- **Auth:** Enable/disable auth strategies, OAuth provider credentials
- **Environment Variables:** Manage per-environment secrets
- **Deployment:** Target clusters (Aurora-hosted or custom K8s)

**4. Deployments**
- Timeline view of deployments per environment
- Status indicators: pending → building → deploying → success/failed
- Expandable logs viewer (real-time streaming)
- Rollback button (with confirmation)
- Environment tabs (dev/staging/prod)
- Deployment diff viewer (compare commits)

**5. Monitoring Dashboard**
- Service health grid (up/down with color coding)
- Request rate, latency (p50/p95/p99), error rate graphs (Grafana embeds)
- Recent logs aggregated across services (with search/filter)
- Database connection pool stats
- Alerts feed (active alerts, history)

**6. API Documentation**
- Auto-generated from OpenAPI specs
- Interactive API explorer (try-it-out sandbox)
- Authentication examples (curl, JavaScript, Python)
- SDK generation downloads (TypeScript, Python, Go clients)
- Webhook documentation

### Tech Stack

- **Frontend:** React 18 + TypeScript
- **Styling:** Tailwind CSS
- **State Management:** TanStack Query (React Query)
- **Schema Designer:** React Flow
- **Charts:** Recharts or Chart.js
- **Build Tool:** Vite
- **Testing:** Vitest + React Testing Library

---

## Testing Strategy

Comprehensive testing at all levels to ensure quality.

### Generated Service Tests

Each microservice includes a complete test suite:

```
users-service/
├── tests/
│   ├── unit/
│   │   ├── repository.test.ts    # DB layer mocks
│   │   ├── service.test.ts       # Business logic tests
│   │   └── validation.test.ts    # Input validation
│   ├── integration/
│   │   ├── api.test.ts           # Full API endpoint tests
│   │   ├── auth.test.ts          # Auth integration
│   │   └── database.test.ts      # Real DB tests
│   └── e2e/
│       └── flows.test.ts         # End-to-end user flows
└── test-utils/
    ├── fixtures.ts               # Test data factories
    ├── mocks.ts                  # Mock dependencies
    └── setup.ts                  # Test environment setup
```

**Coverage Requirements:**
- Unit tests: 80% minimum for generated code
- Integration tests: All CRUD endpoints covered
- E2E tests: Happy path for each service
- Tests run in CI before every deployment
- Failed tests block deployment

**Test Execution:**
```bash
# Run all tests
npm test

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test -- users.test.ts

# Run in watch mode
npm test -- --watch
```

### Aurora Platform Tests

**1. Code Generator Tests (Go)**
```go
func TestGenerateTypeScriptService(t *testing.T) {
    schema := loadTestSchema("users.yaml")
    code, err := generator.Generate(schema, "typescript", "rest")
    require.NoError(t, err)
    
    // Verify syntax is valid TypeScript
    assert.NoError(t, validateTypeScript(code))
    
    // Verify extension points exist
    assert.Contains(t, code["hooks.ts"], "beforeCreate")
    assert.Contains(t, code["custom-endpoints.ts"], "export const customRoutes")
    
    // Verify protected file headers
    assert.Contains(t, code["controller.generated.ts"], "AUTO-GENERATED")
}

func TestRegenerationPreservesExtensions(t *testing.T) {
    // Generate initial code
    code1 := generator.Generate(schema1, "typescript", "rest")
    
    // User edits extension file
    code1["hooks.ts"] = "// Custom user code"
    
    // Update schema and regenerate
    schema2 := addColumn(schema1, "phone", "string")
    code2 := generator.Regenerate(schema2, code1, "typescript", "rest")
    
    // Extension file preserved, protected files updated
    assert.Equal(t, "// Custom user code", code2["hooks.ts"])
    assert.NotEqual(t, code1["controller.generated.ts"], code2["controller.generated.ts"])
}
```

**2. Deployment Worker Tests**
```go
func TestDeploymentWorker(t *testing.T) {
    // Mock K8s API server
    k8sServer := httptest.NewServer(mockK8sHandler())
    defer k8sServer.Close()
    
    // Create deployment job
    job := &DeploymentJob{
        ProjectID:   "test-project",
        Environment: "dev",
        CommitSHA:   "abc123",
    }
    
    // Execute deployment
    worker := NewDeploymentWorker(k8sServer.URL)
    result := worker.Process(job)
    
    // Verify success
    assert.Equal(t, "success", result.Status)
    assert.NotEmpty(t, result.DeploymentURL)
}
```

**3. Integration Tests**
```go
func TestFullProjectGeneration(t *testing.T) {
    // Spin up test Aurora instance
    aurora := startTestAurora(t)
    defer aurora.Shutdown()
    
    // Create project via API
    project := aurora.CreateProject("test-app", "typescript", "postgres")
    
    // Apply schema
    schema := loadTestSchema("users.yaml")
    aurora.ApplySchema(project.ID, schema)
    
    // Wait for code generation
    waitForJobCompletion(t, project.ID, "code-generation")
    
    // Verify git repo created
    repo := aurora.GetGitRepo(project.ID)
    assert.FileExists(t, repo.Path("services/users-service/src/controller.generated.ts"))
    
    // Deploy to test cluster
    deployment := aurora.Deploy(project.ID, "dev")
    waitForDeployment(t, deployment.ID)
    
    // Verify service responds
    resp := httpGet(t, deployment.URL + "/api/v1/users")
    assert.Equal(t, 200, resp.StatusCode)
}
```

**4. Load Tests**
```go
func TestWorkerPoolScaling(t *testing.T) {
    // Create 100 concurrent code generation jobs
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            job := createCodeGenJob(fmt.Sprintf("project-%d", i))
            result := submitJob(job)
            assert.Equal(t, "success", result.Status)
        }(i)
    }
    
    // All jobs should complete within reasonable time
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        // Success
    case <-time.After(5 * time.Minute):
        t.Fatal("Worker pool failed to handle load")
    }
}
```

---

## Monitoring & Observability

Built-in observability for all generated services and the Aurora platform.

### Metrics (Prometheus + Grafana)

**Service-Level Metrics:**
- Request rate (requests/second per endpoint)
- Latency percentiles (p50, p90, p95, p99)
- Error rate (4xx, 5xx)
- Active connections

**Database Metrics:**
- Connection pool usage (active/idle/max)
- Query duration (p50, p95, p99)
- Transaction counts
- Deadlocks and slow queries

**Gateway Metrics:**
- Total requests per route
- Auth failures and rate limit hits
- Cache hit/miss ratio
- Upstream response times

**Worker Metrics:**
- Job queue depth (pending jobs)
- Job processing time
- Worker CPU/memory usage
- Job failure rate

**Automatic Instrumentation:**
All generated services include Prometheus client libraries with automatic metric collection.

### Logging (ELK Stack or Loki)

**Structured JSON Logging:**
```json
{
  "timestamp": "2026-05-11T10:30:00Z",
  "level": "info",
  "service": "users-service",
  "trace_id": "abc123",
  "span_id": "def456",
  "method": "GET",
  "path": "/api/v1/users/123",
  "status": 200,
  "duration_ms": 45,
  "user_id": "user-789"
}
```

**Log Levels by Environment:**
- **Development:** DEBUG (verbose)
- **Staging:** INFO (standard)
- **Production:** INFO (standard), ERROR (always)

**Centralized Aggregation:**
- All service logs flow to centralized ELK/Loki
- Correlation IDs trace requests across services
- Full-text search and filtering
- Log retention: 30 days (configurable)

### Distributed Tracing (Jaeger/Tempo)

**Automatic Trace Propagation:**
- W3C Trace Context headers
- Service-to-service call visualization
- Database query tracing
- Full trace: Gateway → Service → Database

**Example Trace:**
```
Trace ID: abc123
├─ api-gateway: GET /api/v1/orders (150ms)
   ├─ orders-service: GET /api/v1/orders (145ms)
      ├─ DB Query: SELECT * FROM orders (30ms)
      └─ users-service: GET /api/v1/users/456 (80ms)
         └─ DB Query: SELECT * FROM users WHERE id=456 (15ms)
```

### Alerting

**Pre-configured Alerts:**
1. **Service Down:** No successful requests in 5 minutes
2. **High Error Rate:** 5xx errors > 5% for 10 minutes
3. **High Latency:** p95 latency > 1s for 15 minutes
4. **Database Issues:** Connection pool exhausted
5. **Deployment Failed:** Rollout status = failed

**Notification Channels:**
- Webhook integration (Slack, PagerDuty, Discord)
- Email notifications
- SMS (via Twilio) for critical alerts

**Alert Configuration (aurora.yaml):**
```yaml
alerts:
  enabled: true
  channels:
    - type: slack
      webhook_url: https://hooks.slack.com/...
    - type: email
      recipients: [ops@company.com]
  
  rules:
    - name: high_error_rate
      condition: error_rate > 0.05
      duration: 10m
      severity: warning
    
    - name: service_down
      condition: up == 0
      duration: 5m
      severity: critical
```

### Health Checks

**Every Service Exposes:**
- `/health` - Liveness probe (is process running?)
- `/ready` - Readiness probe (can accept traffic?)

**Health Check Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-05-11T10:30:00Z",
  "checks": {
    "database": "ok",
    "redis": "ok",
    "disk_space": "ok"
  },
  "uptime_seconds": 86400
}
```

**K8s Probe Configuration:**
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 3000
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 3000
  initialDelaySeconds: 5
  periodSeconds: 5
```

---

## Implementation Phases

### Phase 1: Foundation (Weeks 1-4)
- Aurora Core API skeleton (Go + Gin)
- PostgreSQL schema for metadata
- Basic Web UI (React shell)
- CLI tool foundation (project init, config)
- Schema definition format (YAML spec)

### Phase 2: Code Generation (Weeks 5-8)
- TypeScript template system
- Python template system
- Code generation worker (basic CRUD)
- Extension points implementation
- Git repository creation

### Phase 3: Auth & Gateway (Weeks 9-12)
- Auth service generation (JWT + Basic Auth)
- Kong Gateway configuration generation
- API Gateway integration
- OAuth 2.0 provider integration

### Phase 4: Deployment (Weeks 13-16)
- K8s manifest generation (Kustomize)
- Deployment worker
- GitOps workflow (GitHub Actions)
- Docker image building
- Environment management (dev/staging/prod)

### Phase 5: UI & DX (Weeks 17-20)
- Schema designer (React Flow)
- Deployment dashboard
- Monitoring dashboard (Grafana embeds)
- CLI enhancements (dev environment, logs)
- Documentation site

### Phase 6: Polish & Launch (Weeks 21-24)
- Testing automation (unit, integration, e2e)
- Load testing and performance tuning
- Security audit
- Documentation completion
- Beta launch

---

## Security Considerations

### Generated Service Security
- **Input Validation:** All endpoints validate request bodies using Zod/Pydantic
- **SQL Injection Prevention:** Parameterized queries only, no string concatenation
- **XSS Prevention:** Output encoding, Content-Security-Policy headers
- **CSRF Protection:** CSRF tokens for state-changing operations
- **Rate Limiting:** Per-user and per-IP limits at gateway
- **Auth Token Security:** JWT with short expiry (1 hour), secure refresh tokens

### Platform Security
- **API Authentication:** All Aurora API calls require JWT tokens
- **Secrets Management:** K8s Secrets, never in git
- **Network Policies:** Services can only communicate as defined in policies
- **Container Security:** Non-root users, read-only root filesystems
- **Dependency Scanning:** Automated CVE scanning in CI
- **Audit Logging:** All sensitive operations logged (deployments, auth changes)

### User Data Isolation
- **Database Isolation:** Production services use dedicated DB instances
- **K8s Namespaces:** Each project in its own namespace
- **Network Segmentation:** Projects cannot access each other's services
- **Secret Encryption:** K8s secrets encrypted at rest

---

## Success Metrics

### User Metrics
- **Time to First Deploy:** < 10 minutes from signup to live API
- **Code Generation Time:** < 30 seconds for typical 5-table project
- **Deployment Time:** < 5 minutes for full stack (dev environment)
- **Developer Satisfaction:** NPS > 50

### Platform Metrics
- **Service Uptime:** 99.9% for Aurora control plane
- **Generated Service Uptime:** 99.5% on Aurora infrastructure
- **Worker Queue Latency:** < 1 minute from job submission to processing
- **API Response Time:** p95 < 200ms for Aurora Core API

### Business Metrics
- **Active Projects:** 100+ within 6 months
- **User Retention:** 60% monthly active users
- **Deployment Frequency:** Average 10+ deploys per project per month

---

## Open Questions & Future Enhancements

### MVP Scope Questions
- Should we support GraphQL subscriptions (real-time) or just queries/mutations?
- Do we need a visual query builder in the UI, or is the auto-generated API documentation sufficient?
- Should the CLI support offline mode for schema editing?

### Post-MVP Features
- **Additional Languages:** Rust, Java, C# (Phases 7-9)
- **Additional Databases:** DynamoDB, CockroachDB, Cassandra
- **Event Bus Integration:** Kafka, RabbitMQ for async workflows
- **Service Mesh:** Istio/Linkerd for advanced traffic management
- **Multi-region Deployment:** Global load balancing
- **Collaboration Features:** Team workspaces, role-based project access
- **Marketplace:** Community-contributed templates and plugins
- **AI-Powered Features:** Schema suggestion from natural language

---

## Conclusion

Aurora provides a complete Backend-as-a-Service solution that balances automation with developer control. By generating production-ready microservices with clear extension points, full GitOps workflows, and comprehensive observability, Aurora accelerates backend development without creating vendor lock-in or sacrificing customization.

The hybrid architecture (monolithic core + worker pools) ensures the platform can scale to handle many concurrent code generation and deployment jobs while remaining simple to operate. The focus on developer experience—through both a powerful CLI and intuitive Web UI—makes Aurora accessible to solo developers while meeting the needs of growing teams.

With support for multiple languages, databases, and deployment targets, Aurora can adapt to diverse technical requirements while maintaining consistent patterns and best practices across all generated code.
