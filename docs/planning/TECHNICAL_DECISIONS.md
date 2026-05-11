# Aurora Technical Decisions Log

**Purpose:** Document key technical decisions, rationale, trade-offs, and alternatives considered.

**Format:** Architecture Decision Records (ADR) format  
**Status:** Living document - updated as decisions are made  
**Last Updated:** May 11, 2026

---

## Table of Contents

- [How to Use This Document](#how-to-use-this-document)
- [Decision Template](#decision-template)
- [Active Decisions](#active-decisions)
- [Pending Decisions](#pending-decisions)
- [Superseded Decisions](#superseded-decisions)

---

## How to Use This Document

### When to Add a Decision

Document a decision when:
- It affects system architecture
- It involves significant investment (time or money)
- It constrains future options
- It impacts multiple teams or systems
- Future developers will ask "why did we do it this way?"

### Decision Status

- **Proposed:** Under consideration
- **Accepted:** Decision made and implemented
- **Deprecated:** No longer valid but kept for history
- **Superseded:** Replaced by a newer decision

### Review Process

1. Author creates decision record
2. Team reviews in tech sync
3. Discussion and refinement
4. Final decision by tech lead or engineering manager
5. Document updated with final decision

---

## Decision Template

```markdown
## ADR-XXX: [Title]

**Status:** [Proposed | Accepted | Deprecated | Superseded by ADR-YYY]  
**Date:** YYYY-MM-DD  
**Deciders:** [Names]  
**Technical Story:** [Link to issue/epic]

### Context
What is the issue we're facing? What forces are at play?

### Decision
What did we decide to do?

### Rationale
Why did we choose this approach?

### Alternatives Considered
What other options did we evaluate?

### Consequences
What are the positive and negative impacts?

### Implementation Notes
How should this be implemented? Any specific guidelines?

### Revision Points
Under what conditions should we revisit this decision?
```

---

## Active Decisions

### ADR-001: Go as Primary Backend Language

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Engineering Team  
**Technical Story:** Foundation architecture

#### Context

Needed to choose a primary language for the Aurora backend. Requirements:
- High performance for API serving
- Strong concurrency support (for code generation tasks)
- Good ecosystem for web frameworks
- Type safety
- Easy deployment
- Strong standard library
- Growing community

#### Decision

Use **Go (Golang) 1.26+** as the primary backend language.

#### Rationale

**Performance:**
- Compiled language with excellent performance characteristics
- Efficient memory usage
- Fast startup times (important for containers)
- Built-in concurrency with goroutines

**Developer Experience:**
- Simple language with fewer footguns than Rust
- Fast compilation
- Strong tooling (gofmt, go vet, gopls)
- Excellent standard library
- Easy to learn for new team members

**Ecosystem:**
- Mature web frameworks (Gin, Echo, Chi)
- Excellent database libraries (GORM, sqlx)
- Strong Docker/Kubernetes support
- Good for microservices

**Deployment:**
- Single binary deployment
- Small container images
- No runtime dependencies
- Cross-compilation easy

#### Alternatives Considered

**1. Node.js/TypeScript**
- Pros: Large ecosystem, team familiarity, async by default
- Cons: Runtime overhead, single-threaded (workarounds exist), type safety limitations
- Why not: Performance concerns for code generation workloads

**2. Python**
- Pros: Rapid development, great libraries, AI/ML ecosystem
- Cons: Slower performance, GIL limits concurrency, deployment complexity
- Why not: Performance and concurrency limitations

**3. Rust**
- Pros: Blazing fast, memory safe, no GC
- Cons: Steep learning curve, longer development time, smaller ecosystem for web
- Why not: Development speed priority over marginal performance gains

**4. Java/Kotlin**
- Pros: Mature ecosystem, excellent tooling, strong typing
- Cons: Heavy runtime (JVM), verbose (Java), slower startup
- Why not: Startup time and container size concerns

#### Consequences

**Positive:**
- ✅ Excellent performance for API serving
- ✅ Great concurrency for parallel code generation
- ✅ Small, fast containers
- ✅ Strong type safety
- ✅ Easy deployment
- ✅ Team can learn quickly

**Negative:**
- ❌ Smaller ecosystem than Node.js or Python
- ❌ Some libraries less mature
- ❌ Generic support still evolving
- ❌ No automatic memory management (but GC handles most cases)

**Neutral:**
- Hiring pool smaller than JavaScript but healthy
- Community growing rapidly

#### Implementation Notes

**Standards:**
- Use Go 1.26+ (or latest stable)
- Follow official Go style guide
- Use `gofmt` and `goimports`
- Leverage `go mod` for dependencies
- Pin dependencies in production

**Project Structure:**
```
internal/    - Private application code
pkg/        - Public library code
cmd/        - Entry points
```

**Testing:**
- Use standard `testing` package
- Use `testify` for assertions (team preference)
- Aim for 80%+ coverage

#### Revision Points

**Revisit if:**
- Performance becomes an issue (unlikely)
- Need features only available in other languages
- Team expertise changes significantly
- Go ecosystem lacks critical libraries

**Next review:** April 2027 (after 1 year of production use)

---

### ADR-002: PostgreSQL as Primary Database

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Engineering Team  
**Technical Story:** Database selection

#### Context

Need a database for Aurora's control plane (projects, schemas, metadata). Requirements:
- ACID transactions (critical for schema versioning)
- Strong consistency
- Rich query capabilities
- JSON support (for flexibility)
- Good Go support
- Production-proven
- Easy to operate

#### Decision

Use **PostgreSQL 15+** as the primary database.

#### Rationale

**Reliability:**
- Battle-tested in production at scale
- Strong ACID guarantees
- Excellent data integrity
- Proven consistency model

**Features:**
- Native UUID support
- JSON/JSONB for flexible schemas
- Full-text search
- Array and custom types
- Rich indexing options
- Window functions and CTEs

**Ecosystem:**
- Excellent Go support (pgx, GORM)
- Great tooling (pgAdmin, psql)
- Docker support
- Managed offerings (RDS, Cloud SQL)

**Operational:**
- Well-understood operations
- Mature backup/restore tools
- Replication and HA options
- Good monitoring tools

#### Alternatives Considered

**1. MySQL/MariaDB**
- Pros: Widespread adoption, good performance, familiar
- Cons: Less feature-rich than PostgreSQL, JSON support weaker
- Why not: PostgreSQL features (UUID, JSON, arrays) more aligned with needs

**2. MongoDB**
- Pros: Flexible schema, JSON-native, horizontal scaling
- Cons: Weaker consistency guarantees, transaction complexity
- Why not: Need strong consistency for schema versioning

**3. SQLite**
- Pros: Zero configuration, embedded, simple
- Cons: Single-writer limit, not suitable for production multi-user
- Why not: Not suitable for production SaaS

**4. CockroachDB**
- Pros: Distributed, PostgreSQL-compatible, strong consistency
- Cons: Operational complexity, cost, overkill for current scale
- Why not: Can migrate later if needed, PostgreSQL compatibility helps

#### Consequences

**Positive:**
- ✅ Strong consistency for critical operations
- ✅ ACID transactions
- ✅ Rich query capabilities
- ✅ JSON support for flexibility
- ✅ Excellent tooling
- ✅ Easy to find expertise

**Negative:**
- ❌ Vertical scaling limits (can address later)
- ❌ Sharding complexity (if needed)
- ❌ Not "cool" like NoSQL (not actually a problem)

**Neutral:**
- Managed services available when needed
- Can add read replicas for scaling

#### Implementation Notes

**Version:** PostgreSQL 15+

**Schema Management:**
- Use GORM Auto-Migrate for development
- Manual migrations for production (future)
- Version all schema changes

**Connection Pooling:**
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Conventions:**
- UUID primary keys (`gen_random_uuid()`)
- `created_at` and `updated_at` timestamps
- Soft deletes where appropriate
- Use indexes liberally

**Backup Strategy:**
- Daily backups
- 30-day retention
- Point-in-time recovery
- Test restore monthly

#### Revision Points

**Revisit if:**
- Hit scaling limits (unlikely before 100K+ projects)
- Need multi-region (consider CockroachDB)
- Cost becomes prohibitive
- Performance issues (optimize first)

**Next review:** October 2026 (after Phase 4 scaling work)

---

### ADR-003: GORM as ORM Layer

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Backend Team  
**Technical Story:** Database abstraction

#### Context

Need an abstraction layer between Go code and PostgreSQL. Options:
- Raw SQL with database/sql
- Query builder (squirrel, goqu)
- ORM (GORM, ent, sqlc)

Requirements:
- Type safety
- Reduce boilerplate
- Auto-migrations (development)
- Association handling
- Good performance

#### Decision

Use **GORM v2** as the ORM layer.

#### Rationale

**Productivity:**
- Significantly reduces boilerplate
- Auto-migrations for rapid development
- Association handling built-in
- Hooks for lifecycle events

**Type Safety:**
- Compile-time checking
- Strong typing for queries
- IntelliSense support

**Features:**
- Transactions support
- Preloading (eager loading)
- Soft deletes
- Scopes for reusable queries
- Plugin system

**Ecosystem:**
- Well-maintained
- Large community
- Good documentation
- Plugins available (caching, sharding)

#### Alternatives Considered

**1. Raw SQL (database/sql)**
- Pros: Full control, no magic, best performance
- Cons: Lots of boilerplate, manual type mapping, more error-prone
- Why not: Productivity loss not worth slight performance gain

**2. sqlx**
- Pros: Minimal abstraction, good performance, struct scanning
- Cons: Still much boilerplate, no auto-migrations
- Why not: Not enough productivity gain over raw SQL

**3. sqlc (code generation)**
- Pros: Type-safe SQL, good performance, no runtime overhead
- Cons: Requires SQL writing, no auto-migrations, less flexible
- Why not: Defeats purpose of ORM productivity

**4. ent (Facebook)**
- Pros: Code generation, type-safe, graph traversal
- Cons: Steep learning curve, more opinionated, schema-first
- Why not: GORM more accessible to team

#### Consequences

**Positive:**
- ✅ Rapid development
- ✅ Less boilerplate
- ✅ Type safety
- ✅ Auto-migrations (development)
- ✅ Team productivity

**Negative:**
- ❌ Performance overhead (small)
- ❌ Less control over SQL
- ❌ Learning curve for advanced features
- ❌ Complex queries may need raw SQL

**Neutral:**
- Can use raw SQL when needed
- Can optimize queries later

#### Implementation Notes

**Usage Guidelines:**

**Good for:**
- CRUD operations
- Simple queries
- Associations
- Transactions

**Use raw SQL for:**
- Complex joins
- Performance-critical queries
- Bulk operations
- Analytics queries

**Example:**
```go
// Good use of GORM
var project domain.Project
db.First(&project, "id = ?", id)

// Use raw SQL for complex queries
db.Raw("SELECT ... complex query ...").Scan(&results)
```

**Configuration:**
```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Error), // Only errors in production
    NowFunc: func() time.Time {
        return time.Now().UTC() // Always use UTC
    },
})
```

**Best Practices:**
- Always use prepared statements (default)
- Use transactions for multi-step operations
- Preload associations to avoid N+1
- Use Scopes for reusable query logic
- Profile slow queries and optimize

#### Revision Points

**Revisit if:**
- Performance becomes critical bottleneck
- Complex query needs outgrow GORM
- Better alternative emerges
- Team expertise changes

**Next review:** October 2026

---

### ADR-004: Gin as Web Framework

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Backend Team  
**Technical Story:** Web framework selection

#### Context

Need a web framework for Aurora's REST API. Requirements:
- Fast and efficient
- Good middleware support
- Easy routing
- JSON handling
- Active maintenance
- Good documentation

#### Decision

Use **Gin** as the web framework.

#### Rationale

**Performance:**
- One of fastest Go web frameworks
- Built on httprouter (fast routing)
- Efficient JSON encoding
- Low memory overhead

**Developer Experience:**
- Simple, intuitive API
- Similar to Express.js (familiar)
- Good documentation
- Active community

**Features:**
- Middleware support
- Route grouping
- JSON binding and validation
- Error handling
- Static file serving

#### Alternatives Considered

**1. stdlib net/http**
- Pros: No dependencies, full control, standard
- Cons: More boilerplate, no middleware system, manual routing
- Why not: Productivity loss

**2. Echo**
- Pros: Similar to Gin, good performance, feature-rich
- Cons: More opinionated, slightly less popular
- Why not: Gin more familiar to team, larger community

**3. Chi**
- Pros: Stdlib-compatible, lightweight, middleware focus
- Cons: Less features out of box, smaller community
- Why not: Want more batteries included

**4. Fiber**
- Pros: Fastest Go framework, Express-like API
- Cons: Not stdlib-compatible, uses fasthttp (different ecosystem)
- Why not: Stdlib compatibility preferred

#### Consequences

**Positive:**
- ✅ Excellent performance
- ✅ Easy to learn
- ✅ Good middleware ecosystem
- ✅ Active maintenance

**Negative:**
- ❌ Another dependency
- ❌ Slightly opinionated

**Neutral:**
- Can switch to stdlib or Echo later if needed

#### Implementation Notes

**Middleware Stack:**
```go
router := gin.New()
router.Use(gin.Recovery())
router.Use(middleware.Logger())
// router.Use(middleware.Auth())     // Phase 3
// router.Use(middleware.CORS())     // Phase 3
// router.Use(middleware.RateLimit()) // Phase 3
```

**Route Organization:**
```go
v1 := router.Group("/api/v1")
{
    projects := v1.Group("/projects")
    {
        projects.POST("", handler.CreateProject)
        projects.GET("/:id", handler.GetProject)
        // ...
    }
}
```

**Error Handling:**
```go
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
    return
}
```

#### Revision Points

**Revisit if:**
- Performance becomes issue (unlikely)
- Need features Gin doesn't provide
- Gin maintenance stops

**Next review:** Not needed unless issues arise

---

### ADR-005: Clean Architecture Pattern

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Engineering Team  
**Technical Story:** Code organization

#### Context

Need to organize code for:
- Testability
- Maintainability
- Flexibility
- Clear responsibilities
- Scalability

#### Decision

Implement **Clean Architecture** with layered approach:
- API Layer (handlers)
- Service Layer (business logic)
- Repository Layer (data access)
- Domain Layer (entities)

#### Rationale

**Separation of Concerns:**
- Each layer has single responsibility
- Clear boundaries
- Easy to understand

**Testability:**
- Each layer independently testable
- Easy to mock dependencies
- Unit tests are simple

**Flexibility:**
- Easy to swap implementations
- Can change database without affecting business logic
- Can add new API (gRPC) without changing core

**Maintainability:**
- Changes isolated to specific layers
- Easier onboarding for new developers
- Clear where code belongs

#### Alternatives Considered

**1. Monolithic / MVC**
- Pros: Simpler, faster to start, less code
- Cons: Tight coupling, hard to test, difficult to change
- Why not: Doesn't scale well with team/features

**2. Domain-Driven Design (Full)**
- Pros: Rich domain model, very flexible
- Cons: Complex, overkill for CRUD, steep learning curve
- Why not: Too much for current needs

**3. Hexagonal Architecture**
- Pros: Excellent separation, port/adapter pattern
- Cons: More complex, more code
- Why not: Clean Architecture simpler but similar benefits

#### Consequences

**Positive:**
- ✅ Testable
- ✅ Flexible
- ✅ Clear organization
- ✅ Easy to onboard

**Negative:**
- ❌ More files and interfaces
- ❌ More boilerplate
- ❌ Potentially over-engineered for simple CRUD

**Neutral:**
- Worth it as project grows
- Can simplify if needed

#### Implementation Notes

**Layer Dependencies:**
```
API → Service → Repository → Domain
```

Domain has zero dependencies.

**File Structure:**
```
internal/
├── api/          # HTTP handlers
├── service/      # Business logic
├── repository/   # Data access
└── domain/       # Entities
```

**Example:**
```go
// domain/project.go
type Project struct { /* ... */ }

// repository/project_repo.go
type ProjectRepository interface {
    Create(*Project) error
    FindByID(string) (*Project, error)
}

// service/project_service.go
type ProjectService struct {
    repo ProjectRepository
}

// api/handlers/project.go
type ProjectHandler struct {
    service *ProjectService
}
```

#### Revision Points

**Revisit if:**
- Becomes overly complex
- Team finds it too restrictive
- Better pattern emerges

**Next review:** October 2026

---

### ADR-006: Go Templates for Code Generation

**Status:** Accepted  
**Date:** 2026-05-11  
**Deciders:** Backend Team  
**Technical Story:** Phase 2 - Code Generation Engine

#### Context

Phase 2 requires generating code in multiple languages (TypeScript, Python). Need a template system:
- Render code from templates
- Support variables and logic
- Maintain readability
- Good performance
- Easy to maintain

#### Decision

Use **Go's text/template** package for code generation.

#### Rationale

**Built-in:**
- No external dependencies
- Well-documented
- Standard library quality

**Powerful:**
- Variables and expressions
- Conditionals and loops
- Custom functions
- Template composition

**Performance:**
- Compiled templates
- Efficient rendering
- Low memory usage

**Type Safety:**
- Compile-time template validation
- Type-safe data passing

#### Alternatives Considered

**1. Handlebars (Go port)**
- Pros: Simpler syntax, familiar to JS developers
- Cons: External dependency, less powerful, smaller ecosystem
- Why not: Go templates more powerful and native

**2. Pongo2 (Django-like)**
- Pros: Familiar to Python developers, powerful
- Cons: External dependency, learning curve
- Why not: Go templates sufficient and native

**3. Code Generation Library (e.g., Jennifer)**
- Pros: Type-safe code generation, programmatic
- Cons: Only works for Go, different approach per language
- Why not: Need multi-language support, templates more maintainable

**4. AST Manipulation**
- Pros: Ultimate flexibility, type-safe
- Cons: Complex, language-specific, hard to maintain
- Why not: Overkill, templates more maintainable

#### Consequences

**Positive:**
- ✅ No dependencies
- ✅ Good performance
- ✅ Flexible and powerful
- ✅ Standard library

**Negative:**
- ❌ Go template syntax (not as pretty as Handlebars)
- ❌ Limited logic in templates (by design)

**Neutral:**
- Can use custom functions to extend
- Template readability depends on design

#### Implementation Notes

**Template Structure:**
```
templates/
├── typescript/
│   ├── express/
│   │   ├── server.ts.tmpl
│   │   └── routes.ts.tmpl
│   └── common/
│       └── tsconfig.json.tmpl
└── python/
    ├── fastapi/
    │   ├── main.py.tmpl
    │   └── routers.py.tmpl
    └── common/
        └── requirements.txt.tmpl
```

**Custom Functions:**
```go
funcMap := template.FuncMap{
    "toCamelCase":  toCamelCase,
    "toPascalCase": toPascalCase,
    "toSnakeCase":  toSnakeCase,
    "pluralize":    pluralize,
}
```

**Usage:**
```go
tmpl := template.Must(template.New("server").Funcs(funcMap).ParseFiles("server.ts.tmpl"))
err := tmpl.Execute(writer, data)
```

**Best Practices:**
- Keep logic in code, not templates
- Use helper functions for formatting
- Test templates independently
- Version templates with code

#### Revision Points

**Revisit if:**
- Templates become unmaintainable
- Need features Go templates lack
- Performance issues (unlikely)
- Better alternative emerges

**Next review:** After Phase 2 completion (July 2026)

---

### ADR-007: YAML for Schema Definition

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Engineering Team  
**Technical Story:** Schema definition format

#### Context

Users need to define their data schema. Requirements:
- Human-readable
- Easy to write
- Supports comments
- Widely understood
- Good tooling support
- Parseable in Go

#### Decision

Use **YAML** for schema definition.

#### Rationale

**Human-Friendly:**
- Very readable
- Whitespace-based (no braces)
- Supports comments
- Familiar to most developers

**Flexible:**
- Nested structures
- Lists and maps
- Anchors and aliases (reuse)
- Multiple documents

**Tooling:**
- Good editor support
- Syntax highlighting
- Validation tools
- Go libraries (gopkg.in/yaml.v3)

**Precedent:**
- Used by Kubernetes, Docker Compose, GitHub Actions
- Familiar to DevOps community

#### Alternatives Considered

**1. JSON**
- Pros: Ubiquitous, strict, good tooling
- Cons: No comments, verbose, less human-friendly
- Why not: Readability important for schemas

**2. TOML**
- Pros: Simple, readable, good for config
- Cons: Less flexible for nested structures, less common
- Why not: YAML more flexible

**3. HCL (HashiCorp)**
- Pros: Readable, expressive, powerful
- Cons: Less common, Go-specific ecosystem
- Why not: YAML more widely known

**4. DSL (Custom)**
- Pros: Perfect fit for our needs
- Cons: Learning curve, tooling needed, maintenance burden
- Why not: YAML sufficient

#### Consequences

**Positive:**
- ✅ Very readable
- ✅ Widely understood
- ✅ Good tooling
- ✅ Comments supported

**Negative:**
- ❌ Whitespace-sensitive (source of errors)
- ❌ Parsing can be ambiguous
- ❌ Security concerns (anchors, aliases) - mitigated by validation

**Neutral:**
- Can convert to/from JSON easily
- Strict validation prevents issues

#### Implementation Notes

**Example Schema:**
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
        required: true
      - name: "created_at"
        type: "timestamp"
        default: "now()"
```

**Validation:**
- Use JSON Schema for YAML validation
- Strict parsing (reject unknown fields)
- Version validation

**Library:**
```go
import "gopkg.in/yaml.v3"

var schema Schema
err := yaml.Unmarshal(data, &schema)
```

#### Revision Points

**Revisit if:**
- Users struggle with YAML
- Security issues arise
- Better format emerges
- Tooling becomes problematic

**Next review:** Not needed unless issues

---

### ADR-008: UUID for Primary Keys

**Status:** Accepted  
**Date:** 2026-04-01  
**Deciders:** Backend Team  
**Technical Story:** Database design

#### Context

Need to choose primary key strategy for database tables. Requirements:
- Unique across systems
- Secure (not guessable)
- Scalable
- Good performance

#### Decision

Use **UUID v4** (random) for all primary keys.

#### Rationale

**Uniqueness:**
- Globally unique (no collisions in practice)
- No coordination needed
- Safe for distributed systems

**Security:**
- Not sequential (no enumeration attacks)
- Not guessable
- Safe to expose in APIs

**Flexibility:**
- Can generate client-side
- No auto-increment issues
- Easy to merge databases
- Good for sharding (future)

**PostgreSQL Support:**
- Native UUID type
- Built-in generation (`gen_random_uuid()`)
- Efficient storage (16 bytes)
- Good index performance

#### Alternatives Considered

**1. Auto-incrementing integers**
- Pros: Simple, small, efficient, ordered
- Cons: Sequential (security issue), coordination needed, not globally unique
- Why not: Security and distributed system concerns

**2. UUID v1 (timestamp-based)**
- Pros: Has ordering, unique
- Cons: Exposes MAC address, timestamp (security), less random
- Why not: Security concerns

**3. ULID (Universally Unique Lexicographically Sortable ID)**
- Pros: Sortable, UUID-compatible, more compact string representation
- Cons: Newer, less support, minimal benefit over UUID
- Why not: UUID more standard

**4. Snowflake IDs**
- Pros: Sortable, efficient, distributed
- Cons: Coordination needed, custom implementation
- Why not: UUID simpler and sufficient

#### Consequences

**Positive:**
- ✅ Globally unique
- ✅ Secure (not guessable)
- ✅ No coordination
- ✅ Distributed-system ready

**Negative:**
- ❌ Larger than integers (16 bytes vs 8)
- ❌ No natural ordering
- ❌ Slightly slower joins (vs int)
- ❌ Less human-friendly

**Neutral:**
- Performance difference negligible for our scale
- Can add created_at for ordering

#### Implementation Notes

**PostgreSQL:**
```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- ...
);
```

**Go:**
```go
import "github.com/google/uuid"

project := &Project{
    ID: uuid.New().String(),
    // ...
}
```

**API:**
- Expose UUIDs in URLs: `/api/v1/projects/550e8400-e29b-41d4-a716-446655440000`
- Use string type in JSON
- Validate UUID format in handlers

#### Revision Points

**Revisit if:**
- Performance becomes an issue (unlikely)
- Need natural ordering (add timestamp field)
- Better alternative emerges

**Next review:** Not needed

---

## Pending Decisions

### PENDING-001: Authentication Strategy

**Status:** Proposed  
**Date:** 2026-05-11  
**Decision Needed By:** Phase 3 start (August 2026)

#### Context

Need authentication for Aurora API. Users need to:
- Register and login
- Secure API access
- Manage sessions
- Reset passwords

#### Options

**Option 1: JWT (JSON Web Tokens)**
- Pros: Stateless, scalable, widely used
- Cons: Revocation complexity, size

**Option 2: Session-based (with Redis)**
- Pros: Easy revocation, server control
- Cons: Stateful, Redis dependency

**Option 3: Hybrid (JWT + refresh tokens in Redis)**
- Pros: Best of both, flexible
- Cons: More complex

#### Recommendation

**Hybrid approach** - JWT for access tokens (short-lived), refresh tokens in Redis (long-lived)

**Next Steps:**
- Research implementation details
- Evaluate libraries (jwt-go, etc.)
- Security audit
- Final decision in June 2026

---

### PENDING-002: File Storage Solution

**Status:** Proposed  
**Date:** 2026-05-11  
**Decision Needed By:** Phase 3 (September 2026)

#### Context

Phase 3 includes file upload. Need storage for:
- User uploads
- Generated code ZIPs
- Documentation assets

#### Options

**Option 1: S3 (or compatible)**
- Pros: Scalable, reliable, CDN integration
- Cons: Cost, vendor dependency

**Option 2: Local filesystem**
- Pros: Simple, no cost, fast
- Cons: Not scalable, backup complexity

**Option 3: MinIO (self-hosted S3-compatible)**
- Pros: S3-compatible, self-hosted, flexible
- Cons: Operational overhead

#### Recommendation

**S3-compatible API** (either AWS S3 or MinIO) for flexibility

**Next Steps:**
- Cost analysis
- MinIO evaluation
- Migration path planning
- Decision in July 2026

---

### PENDING-003: Real-time Architecture

**Status:** Proposed  
**Date:** 2026-05-11  
**Decision Needed By:** Phase 3 (September 2026)

#### Context

Phase 3 includes real-time features (schema change notifications, collaboration). Need:
- Push updates to clients
- Scalable architecture
- Low latency

#### Options

**Option 1: WebSockets**
- Pros: Full-duplex, low latency, native browser support
- Cons: Sticky sessions, scaling complexity

**Option 2: Server-Sent Events (SSE)**
- Pros: Simple, HTTP-based, auto-reconnect
- Cons: Uni-directional only

**Option 3: Long polling**
- Pros: Simple, works everywhere
- Cons: Higher latency, more overhead

**Option 4: Third-party (Pusher, Ably)**
- Pros: Managed, scalable, no ops
- Cons: Cost, vendor dependency

#### Recommendation

**WebSockets** with Redis Pub/Sub for multi-instance

**Next Steps:**
- Prototype WebSocket implementation
- Load testing
- Evaluate libraries (gorilla/websocket)
- Final decision in August 2026

---

## Superseded Decisions

### ADR-000: Initial Framework Choice (Superseded)

**Status:** Superseded by ADR-004  
**Date:** 2026-03-15  
**Original Decision:** Use Echo framework

**Reason for Change:**
Team more familiar with Gin, larger community support

**Migration:**
Switched to Gin in early April 2026 (Phase 1)

---

## Decision Review Schedule

### Quarterly Reviews

**Q3 2026 (July):**
- ADR-002: PostgreSQL (scaling check)
- ADR-005: Clean Architecture (working well?)
- ADR-006: Go Templates (maintainable?)

**Q4 2026 (October):**
- ADR-002: PostgreSQL (Phase 4 scaling review)
- ADR-003: GORM (performance check)

**Q1 2027 (January):**
- ADR-001: Go language (1 year review)
- All core decisions

### Ad-hoc Reviews

**Trigger events:**
- Performance issues
- Security vulnerabilities
- Team feedback
- Better alternatives emerge
- Scale milestones (10K users, 100K projects, etc.)

---

## How to Propose a New Decision

1. **Copy the template** at the top of this document
2. **Fill out all sections** thoroughly
3. **Research alternatives** - understand trade-offs
4. **Create a PR** with your proposal
5. **Schedule discussion** in tech sync
6. **Gather feedback** from team
7. **Make decision** (tech lead or eng manager)
8. **Update document** with final decision
9. **Communicate** to team

---

## Appendix: Decision-Making Framework

### When to Decide Quickly

- Low impact, easily reversible
- Clear best practice
- Team consensus

### When to Decide Carefully

- High impact, hard to reverse
- Significant investment
- No clear best option
- Team divided

### Decision-Making Levels

**Individual Engineer:**
- Implementation details
- Library choice (low impact)
- Code style (within guidelines)

**Tech Lead:**
- Component architecture
- Library choice (medium impact)
- Performance optimization strategies

**Engineering Manager:**
- System architecture
- Language/framework choice
- Cross-team decisions

**VP Engineering / CTO:**
- Strategic technology direction
- Major infrastructure changes
- Build vs buy decisions

---

**Document Maintainer:** Engineering Manager  
**Last Updated:** May 11, 2026  
**Next Review:** July 2026  
**Questions?** Ask in #aurora-dev or tech-sync meeting
