# Contributing to Aurora

Thank you for your interest in contributing to Aurora! This guide will help you understand our development process and coding standards.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Process](#development-process)
- [Code Style Guidelines](#code-style-guidelines)
- [Branch Naming Conventions](#branch-naming-conventions)
- [Commit Message Format](#commit-message-format)
- [Pull Request Process](#pull-request-process)
- [Testing Requirements](#testing-requirements)
- [Code Review Checklist](#code-review-checklist)
- [Communication](#communication)

---

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/baas.git
   cd baas
   ```
3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/ktarun.reddy/baas.git
   ```
4. **Set up development environment:**
   ```bash
   cp .env.example .env
   docker-compose up -d
   go mod download
   make test
   ```

---

## Development Process

### Before Starting Work

1. **Check existing issues/PRs** to avoid duplicate work
2. **Create or comment on an issue** describing what you plan to do
3. **Wait for approval** from maintainers (for large changes)
4. **Sync with upstream:**
   ```bash
   git checkout main
   git fetch upstream
   git merge upstream/main
   ```

### During Development

1. **Create a feature branch** (see [Branch Naming](#branch-naming-conventions))
2. **Make incremental commits** (see [Commit Messages](#commit-message-format))
3. **Write tests** for new functionality
4. **Run tests frequently:**
   ```bash
   make test
   ```
5. **Keep your branch updated:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### After Development

1. **Run full test suite:**
   ```bash
   make test
   make test-integration
   ```
2. **Check code formatting:**
   ```bash
   go fmt ./...
   go vet ./...
   ```
3. **Update documentation** if needed
4. **Create a pull request** (see [PR Process](#pull-request-process))

---

## Code Style Guidelines

### Go Code Standards

We follow standard Go conventions with some additions:

#### Formatting

- **Use `gofmt`** - All code must be formatted with `gofmt`
- **Line length:** Soft limit of 100 characters
- **Imports:** Organize in groups (stdlib, external, internal):
  ```go
  import (
      // Standard library
      "fmt"
      "time"
      
      // External dependencies
      "github.com/gin-gonic/gin"
      "gorm.io/gorm"
      
      // Internal packages
      "github.com/ktarun.reddy/baas/internal/domain"
      "github.com/ktarun.reddy/baas/internal/repository"
  )
  ```

#### Naming Conventions

- **Packages:** Short, lowercase, single word (e.g., `service`, `repository`)
- **Files:** Lowercase with underscores (e.g., `project_service.go`)
- **Types:** PascalCase (e.g., `ProjectService`)
- **Functions/Methods:** PascalCase for exported, camelCase for private
- **Variables:** camelCase (e.g., `projectRepo`)
- **Constants:** PascalCase or ALL_CAPS for exported constants

**Examples:**
```go
// Good
type ProjectService struct {}
func (s *ProjectService) CreateProject() {}
var defaultTimeout = 30 * time.Second

// Bad
type projectService struct {}  // Should be exported if in service package
func (s *ProjectService) create_project() {}  // Don't use snake_case
var DefaultTimeout = 30  // Missing time unit
```

#### Comments

- **Package comments:** Every package needs a doc comment
  ```go
  // Package service provides business logic for Aurora.
  package service
  ```
- **Exported items:** Must have comments starting with the item name
  ```go
  // ProjectService handles business logic for projects.
  // It coordinates between repositories and enforces business rules.
  type ProjectService struct {}
  
  // CreateProject validates and persists a new project.
  // It returns an error if validation fails or database operation fails.
  func (s *ProjectService) CreateProject(project *domain.Project) error {}
  ```
- **Complex logic:** Add inline comments explaining "why", not "what"
  ```go
  // Good
  // Schema versioning starts at 1 to distinguish from unversioned legacy data
  version := latestVersion + 1
  
  // Bad
  // Increment version by 1
  version := latestVersion + 1
  ```

#### Error Handling

- **Always check errors:**
  ```go
  project, err := s.repo.FindByID(id)
  if err != nil {
      return nil, fmt.Errorf("failed to find project %s: %w", id, err)
  }
  ```
- **Wrap errors with context:**
  ```go
  // Good
  return fmt.Errorf("failed to create project %s: %w", project.Name, err)
  
  // Bad
  return err  // Loses context
  ```
- **Use sentinel errors for expected cases:**
  ```go
  var ErrProjectNotFound = errors.New("project not found")
  ```

#### Function Design

- **Keep functions small:** Aim for < 30 lines
- **Single responsibility:** Each function does one thing well
- **Limit parameters:** Max 3-4 parameters, use structs for more
- **Return early:**
  ```go
  // Good
  func (s *Service) DoSomething(id string) error {
      if id == "" {
          return errors.New("id is required")
      }
      
      // Main logic here
      return nil
  }
  
  // Bad
  func (s *Service) DoSomething(id string) error {
      if id != "" {
          // Nested logic
      }
      return nil
  }
  ```

#### Testing Standards

- **Test file naming:** `*_test.go` alongside source files
- **Test function naming:** `TestFunctionName_Scenario`
  ```go
  func TestCreateProject_Success(t *testing.T) {}
  func TestCreateProject_InvalidLanguage(t *testing.T) {}
  func TestCreateProject_DatabaseError(t *testing.T) {}
  ```
- **Use table-driven tests** for multiple scenarios:
  ```go
  func TestProjectValidation(t *testing.T) {
      tests := []struct {
          name    string
          project *domain.Project
          wantErr bool
      }{
          {"valid project", &domain.Project{...}, false},
          {"missing name", &domain.Project{...}, true},
      }
      
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              err := tt.project.Validate()
              if (err != nil) != tt.wantErr {
                  t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
              }
          })
      }
  }
  ```
- **Test structure:** Use Arrange-Act-Assert pattern
  ```go
  func TestCreateProject(t *testing.T) {
      // Arrange
      repo := &mockRepository{}
      service := NewProjectService(repo)
      project := &domain.Project{Name: "Test"}
      
      // Act
      err := service.CreateProject(project)
      
      // Assert
      if err != nil {
          t.Errorf("unexpected error: %v", err)
      }
  }
  ```

### JSON API Conventions

- **Field names:** Use snake_case in JSON
  ```go
  type Project struct {
      DatabaseType string `json:"database_type"`  // Not "databaseType"
      APIStyle     string `json:"api_style"`
  }
  ```
- **Timestamps:** Use ISO 8601 format (automatic with `time.Time`)
- **Error responses:** Consistent structure
  ```json
  {
      "error": "validation failed",
      "details": "language must be one of: typescript, python"
  }
  ```

### File Organization

- **One primary type per file:**
  - `project_service.go` contains `ProjectService`
  - `project_service_test.go` contains tests for `ProjectService`
- **Related functionality together:**
  - All project-related handlers in `handlers/project.go`
  - All schema-related logic in `service/schema_service.go`

---

## Branch Naming Conventions

Use descriptive branch names with prefixes:

| Prefix | Purpose | Example |
|--------|---------|---------|
| `feature/` | New features | `feature/project-tags` |
| `bugfix/` | Bug fixes | `bugfix/schema-validation-crash` |
| `hotfix/` | Urgent production fixes | `hotfix/database-connection-leak` |
| `refactor/` | Code refactoring | `refactor/extract-validation-logic` |
| `docs/` | Documentation only | `docs/update-api-examples` |
| `test/` | Test improvements | `test/add-integration-tests` |
| `chore/` | Maintenance tasks | `chore/update-dependencies` |

**Naming rules:**
- Use lowercase with hyphens
- Be descriptive but concise
- Include issue number if applicable: `feature/123-project-tags`

**Examples:**
```bash
# Good
git checkout -b feature/schema-versioning
git checkout -b bugfix/456-nil-pointer-in-handler
git checkout -b refactor/repository-interface

# Bad
git checkout -b fix  # Too vague
git checkout -b Feature/NewThing  # Wrong case
git checkout -b add_project_tags  # Use hyphens, not underscores
```

---

## Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

Must be one of:

- **feat:** New feature
- **fix:** Bug fix
- **docs:** Documentation changes
- **style:** Code style changes (formatting, missing semicolons, etc.)
- **refactor:** Code refactoring without feature changes
- **perf:** Performance improvements
- **test:** Adding or updating tests
- **chore:** Maintenance tasks (dependencies, build, etc.)
- **ci:** CI/CD changes

### Scope

Optional, specifies the affected area:
- `api` - API handlers
- `service` - Service layer
- `repo` - Repository layer
- `domain` - Domain models
- `config` - Configuration
- `db` - Database layer
- `validator` - Validation logic

### Subject

- Use imperative mood: "add" not "added" or "adds"
- Don't capitalize first letter
- No period at the end
- Max 50 characters

### Body

- Optional, but recommended for non-trivial changes
- Explain **what** and **why**, not **how**
- Wrap at 72 characters
- Separate from subject with blank line

### Footer

- Optional
- Reference issues: `Closes #123` or `Fixes #456`
- Breaking changes: `BREAKING CHANGE: description`

### Examples

**Simple commit:**
```
feat(api): add pagination to project list endpoint
```

**With body:**
```
fix(service): prevent duplicate schema versions

The schema service was not checking for existing versions
before creating a new one, leading to duplicate version
numbers when concurrent requests occurred.

This adds a unique constraint check and proper error handling.

Fixes #234
```

**Breaking change:**
```
refactor(api): change project creation response format

The response now includes the full project object instead of
just the ID, making it consistent with other endpoints.

BREAKING CHANGE: POST /api/v1/projects now returns the full
project object in the response body instead of just {"id": "..."}.
Clients should be updated to handle the new response format.

Closes #123
```

**Multiple commits in a feature:**
```
feat(domain): add tags field to project model
feat(service): add tag validation logic
feat(api): expose tags in project endpoints
test(service): add tests for tag validation
docs(readme): document project tags feature
```

---

## Pull Request Process

### Before Creating a PR

1. **Ensure tests pass:**
   ```bash
   make test
   make test-integration
   ```

2. **Format code:**
   ```bash
   go fmt ./...
   go vet ./...
   ```

3. **Update documentation** if needed:
   - Update README.md for user-facing changes
   - Update code comments for API changes
   - Add/update docs/ files for architecture changes

4. **Rebase on main:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### Creating the PR

1. **Push your branch:**
   ```bash
   git push origin feature/your-feature
   ```

2. **Create PR on GitHub**

3. **Fill out PR template** with:
   - **Title:** Clear, concise description (< 70 chars)
   - **Summary:** What does this PR do?
   - **Motivation:** Why is this change needed?
   - **Changes:** Bullet points of key changes
   - **Testing:** How was this tested?
   - **Screenshots:** For UI changes (if applicable)
   - **Related Issues:** Link to related issues

### PR Title Format

Similar to commit messages:
```
feat(api): add pagination to project list endpoint
fix(service): prevent duplicate schema versions
docs: improve onboarding documentation
```

### PR Description Template

```markdown
## Summary
Brief description of what this PR does.

## Motivation
Why is this change needed? What problem does it solve?

## Changes
- Added X functionality to Y
- Updated Z to handle edge case
- Refactored W for better maintainability

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] Tested with PostgreSQL
- [ ] Tested with MongoDB (if applicable)

## Checklist
- [ ] Code follows style guidelines
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] No breaking changes (or documented if necessary)
- [ ] Commit messages follow convention

## Related Issues
Closes #123
Related to #456
```

### PR Size Guidelines

- **Small PRs preferred:** < 400 lines changed
- **If PR is large:** Break into smaller PRs or explain why it must be large
- **One feature per PR:** Don't mix multiple unrelated changes

### Review Process

1. **Automated checks** must pass:
   - All tests
   - Code formatting (go fmt)
   - Linting (go vet)

2. **Code review** by maintainer:
   - At least 1 approval required
   - Address all review comments

3. **Merge:**
   - Squash commits for clean history (usually)
   - Rebase if preserving commit history is important

### After Your PR is Merged

1. **Delete your branch:**
   ```bash
   git branch -d feature/your-feature
   git push origin --delete feature/your-feature
   ```

2. **Sync your fork:**
   ```bash
   git checkout main
   git fetch upstream
   git merge upstream/main
   git push origin main
   ```

---

## Testing Requirements

### Coverage Requirements

All PRs must maintain or improve test coverage:

- **New features:** Must include tests (aim for 80%+ coverage)
- **Bug fixes:** Must include regression test
- **Refactoring:** Existing tests must still pass

### Test Types

**Unit Tests** (required for all new code):
```go
// Test individual functions/methods
func TestProjectService_CreateProject(t *testing.T) {
    // Use mocks for dependencies
    repo := &mockProjectRepository{}
    service := NewProjectService(repo)
    
    // Test logic
}
```

**Integration Tests** (required for API changes):
```go
// Test complete flows with real database
func TestCreateProjectAPI(t *testing.T) {
    // Setup test database
    // Make HTTP request
    // Assert response
}
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run integration tests
make test-integration

# Run specific test
go test -v -run TestCreateProject ./internal/service/

# Run tests with race detector
go test -race ./...
```

### Test Quality

Good tests are:
- **Fast:** Unit tests < 100ms, integration tests < 5s
- **Isolated:** No dependencies between tests
- **Deterministic:** Same result every time
- **Readable:** Clear arrange-act-assert structure
- **Maintainable:** Easy to update when code changes

---

## Code Review Checklist

### For Authors (Self-Review Before Submitting)

**Functionality:**
- [ ] Code does what it's supposed to do
- [ ] Edge cases are handled
- [ ] Error handling is comprehensive
- [ ] No hardcoded values (use constants/config)

**Code Quality:**
- [ ] Follows Go idioms and conventions
- [ ] Functions are small and focused
- [ ] No code duplication
- [ ] Clear variable and function names
- [ ] Comments explain "why", not "what"

**Testing:**
- [ ] All tests pass
- [ ] New tests cover new code
- [ ] Tests are readable and maintainable
- [ ] Edge cases are tested

**Documentation:**
- [ ] Public APIs are documented
- [ ] README updated if needed
- [ ] Complex logic has comments

**Performance:**
- [ ] No obvious performance issues
- [ ] Database queries are efficient
- [ ] No N+1 query problems

**Security:**
- [ ] Input validation is present
- [ ] No SQL injection vulnerabilities
- [ ] Sensitive data is not logged
- [ ] Errors don't leak sensitive information

### For Reviewers

**High Priority:**
- [ ] Correctness: Does the code do what it claims?
- [ ] Security: Any security vulnerabilities?
- [ ] Performance: Any performance issues?
- [ ] Tests: Are tests adequate?

**Medium Priority:**
- [ ] Architecture: Fits with overall design?
- [ ] Maintainability: Easy to understand and modify?
- [ ] Documentation: Adequate for future developers?

**Low Priority (Nice to Have):**
- [ ] Style: Follows conventions?
- [ ] Naming: Clear and consistent?
- [ ] Comments: Helpful and up-to-date?

### Providing Feedback

**Be constructive:**
```
# Good
"Consider extracting this validation logic into a separate function 
for reusability. What do you think?"

# Bad
"This is wrong."
```

**Use prefixes:**
- **[nit]:** Minor style issue (not blocking)
- **[question]:** Asking for clarification
- **[suggestion]:** Optional improvement
- **[blocker]:** Must be fixed before merge

**Be specific:**
```
# Good
"This function could panic if `project` is nil. 
Add a nil check at the start of the function."

# Bad
"This could crash."
```

---

## Communication

### Where to Ask Questions

- **GitHub Issues:** Feature requests, bug reports
- **GitHub Discussions:** General questions, design discussions
- **Pull Request Comments:** Questions about specific code

### Response Times

- **Maintainers aim to respond within:** 2-3 business days
- **PRs reviewed within:** 1 week
- **Issues triaged within:** 1 week

### Getting Help

If you're stuck:
1. **Read the docs:** README, ARCHITECTURE.md, DEVELOPMENT.md
2. **Search existing issues:** Your question may be answered
3. **Ask in GitHub Discussions:** Don't be shy!
4. **Be patient:** Maintainers are volunteers

---

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes for significant contributions
- Credited in commit history

Thank you for contributing to Aurora!
