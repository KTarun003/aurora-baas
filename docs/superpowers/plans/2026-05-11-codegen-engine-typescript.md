# Code Generation Engine - Core Template System + TypeScript Generator

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a template-based code generation engine that generates complete TypeScript/Express.js REST APIs from Aurora schema definitions.

**Architecture:** Template engine using Go's text/template with custom helpers, generator interface for language-specific implementations, file system abstraction for output management, integration with existing Aurora Core API.

**Tech Stack:** Go 1.22+, text/template, Express.js templates, Prisma ORM, TypeScript 5.x

---

## Context

This is **Phase 2 Week 1-3** of Aurora BaaS platform. After completing this plan, the system will:
- Accept schema definitions and generate complete TypeScript projects
- Produce Express.js REST APIs with Prisma ORM integration
- Generate CRUD endpoints for all tables
- Include extension points for custom code
- Provide downloadable project archives

**Not in scope for this plan:**
- Python generator (separate plan)
- GraphQL support (separate plan)
- Client SDKs (separate plan)
- Authentication generation (Phase 3)

---

## File Structure

```
/Users/ktarun.reddy/Coding/examples/baas/
├── internal/
│   ├── codegen/
│   │   ├── types.go                      # Common types and interfaces
│   │   ├── generator.go                  # Generator interface
│   │   ├── engine/
│   │   │   ├── template_engine.go        # Template rendering engine
│   │   │   ├── helpers.go                # Template helper functions
│   │   │   └── file_writer.go            # File system output management
│   │   ├── typescript/
│   │   │   ├── generator.go              # TypeScript generator implementation
│   │   │   ├── prisma.go                 # Prisma schema generation
│   │   │   ├── express.go                # Express.js code generation
│   │   │   └── models.go                 # TypeScript type generation
│   │   └── archive/
│   │       └── zip.go                    # Project archive creation
│   ├── service/
│   │   └── codegen_service.go            # Code generation business logic
│   └── api/
│       └── handlers/
│           └── codegen.go                # Code generation API endpoints
├── templates/
│   └── typescript/
│       ├── base/
│       │   ├── package.json.tmpl
│       │   ├── tsconfig.json.tmpl
│       │   ├── server.ts.tmpl
│       │   ├── database.ts.tmpl
│       │   └── .env.example.tmpl
│       ├── crud/
│       │   ├── controller.generated.ts.tmpl
│       │   ├── routes.generated.ts.tmpl
│       │   └── model.generated.ts.tmpl
│       ├── extensions/
│       │   ├── custom-routes.ts.tmpl
│       │   └── hooks.ts.tmpl
│       ├── prisma/
│       │   └── schema.prisma.tmpl
│       └── docker/
│           ├── Dockerfile.tmpl
│           └── docker-compose.yml.tmpl
└── tests/
    └── codegen/
        ├── engine_test.go
        ├── typescript_test.go
        └── integration_test.go
```

---

## Task 1: Core Types and Interfaces

**Files:**
- Create: `internal/codegen/types.go`
- Create: `internal/codegen/generator.go`

---

- [ ] **Step 1: Write test for GeneratorContext**

```bash
mkdir -p internal/codegen
cat > internal/codegen/types_test.go << 'EOF'
package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

func TestGeneratorContext_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *GeneratorContext
		wantErr bool
	}{
		{
			name: "valid context",
			ctx: &GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "TestProject",
				Language:    "typescript",
				APIStyle:    "rest",
				Database:    "postgres",
				Schema: &validator.SchemaDefinition{
					Version:  "1.0",
					Database: "postgres",
					Tables:   []validator.Table{},
				},
			},
			wantErr: false,
		},
		{
			name: "missing project ID",
			ctx: &GeneratorContext{
				ProjectName: "TestProject",
				Language:    "typescript",
			},
			wantErr: true,
		},
		{
			name: "invalid language",
			ctx: &GeneratorContext{
				ProjectID: "test",
				Language:  "ruby",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
EOF
```

---

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /Users/ktarun.reddy/Coding/examples/baas
go test ./internal/codegen/... -v
```

Expected: FAIL - GeneratorContext type not defined

---

- [ ] **Step 3: Implement GeneratorContext type**

```bash
cat > internal/codegen/types.go << 'EOF'
package codegen

import (
	"errors"
	"time"

	"github.com/ktarun.reddy/baas/pkg/validator"
)

// GeneratorContext contains all information needed for code generation
type GeneratorContext struct {
	ProjectID   string
	ProjectName string
	Language    string // typescript, python
	APIStyle    string // rest, graphql
	Database    string // postgres, mongodb
	Schema      *validator.SchemaDefinition
	OutputDir   string
	Timestamp   time.Time
	Options     map[string]interface{}
}

// Validate checks if the context has all required fields
func (c *GeneratorContext) Validate() error {
	if c.ProjectID == "" {
		return errors.New("project ID is required")
	}

	if c.ProjectName == "" {
		return errors.New("project name is required")
	}

	validLanguages := map[string]bool{
		"typescript": true,
		"python":     true,
	}
	if !validLanguages[c.Language] {
		return errors.New("invalid language: must be typescript or python")
	}

	validAPIStyles := map[string]bool{
		"rest":    true,
		"graphql": true,
	}
	if !validAPIStyles[c.APIStyle] {
		return errors.New("invalid API style: must be rest or graphql")
	}

	validDatabases := map[string]bool{
		"postgres": true,
		"mongodb":  true,
	}
	if !validDatabases[c.Database] {
		return errors.New("invalid database: must be postgres or mongodb")
	}

	if c.Schema == nil {
		return errors.New("schema is required")
	}

	return nil
}

// GeneratedFile represents a single generated file
type GeneratedFile struct {
	Path     string // Relative path within project
	Content  string // File contents
	Mode     int    // File permissions (0644, 0755, etc)
	IsStatic bool   // True if file should never be regenerated
}

// GeneratedProject represents a complete generated project
type GeneratedProject struct {
	ProjectID string
	Language  string
	Files     []*GeneratedFile
	Metadata  map[string]string
}
EOF
```

---

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/codegen/... -v
```

Expected: PASS

---

- [ ] **Step 5: Write test for Generator interface**

```bash
cat > internal/codegen/generator_test.go << 'EOF'
package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// MockGenerator for testing
type MockGenerator struct {
	GenerateFunc func(ctx *GeneratorContext) (*GeneratedProject, error)
	NameFunc     func() string
}

func (m *MockGenerator) Generate(ctx *GeneratorContext) (*GeneratedProject, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx)
	}
	return &GeneratedProject{}, nil
}

func (m *MockGenerator) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock"
}

func TestGenerator_Interface(t *testing.T) {
	mock := &MockGenerator{
		GenerateFunc: func(ctx *GeneratorContext) (*GeneratedProject, error) {
			return &GeneratedProject{
				ProjectID: ctx.ProjectID,
				Language:  ctx.Language,
				Files:     []*GeneratedFile{},
			}, nil
		},
		NameFunc: func() string {
			return "mock-generator"
		},
	}

	ctx := &GeneratorContext{
		ProjectID:   "test-project",
		ProjectName: "Test",
		Language:    "typescript",
		APIStyle:    "rest",
		Database:    "postgres",
		Schema: &validator.SchemaDefinition{
			Version:  "1.0",
			Database: "postgres",
			Tables:   []validator.Table{},
		},
	}

	project, err := mock.Generate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "test-project", project.ProjectID)
	assert.Equal(t, "typescript", project.Language)
	assert.Equal(t, "mock-generator", mock.Name())
}
EOF
```

---

- [ ] **Step 6: Run test to verify it fails**

```bash
go test ./internal/codegen/... -v
```

Expected: FAIL - Generator interface not defined

---

- [ ] **Step 7: Implement Generator interface**

```bash
cat > internal/codegen/generator.go << 'EOF'
package codegen

// Generator is the interface that all code generators must implement
type Generator interface {
	// Generate produces a complete project from the given context
	Generate(ctx *GeneratorContext) (*GeneratedProject, error)

	// Name returns the name of this generator (e.g., "typescript", "python")
	Name() string
}

// Registry manages available generators
type Registry struct {
	generators map[string]Generator
}

// NewRegistry creates a new generator registry
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]Generator),
	}
}

// Register adds a generator to the registry
func (r *Registry) Register(gen Generator) {
	r.generators[gen.Name()] = gen
}

// Get retrieves a generator by name
func (r *Registry) Get(name string) (Generator, bool) {
	gen, ok := r.generators[name]
	return gen, ok
}

// List returns all registered generator names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.generators))
	for name := range r.generators {
		names = append(names, name)
	}
	return names
}
EOF
```

---

- [ ] **Step 8: Run test to verify it passes**

```bash
go test ./internal/codegen/... -v
```

Expected: PASS

---

- [ ] **Step 9: Commit core types**

```bash
git add internal/codegen/
git commit -m "feat(codegen): add core types and generator interface"
```

---

## Task 2: Template Engine

**Files:**
- Create: `internal/codegen/engine/template_engine.go`
- Create: `internal/codegen/engine/helpers.go`
- Test: `internal/codegen/engine/template_engine_test.go`

---

- [ ] **Step 1: Write test for template engine**

```bash
mkdir -p internal/codegen/engine
cat > internal/codegen/engine/template_engine_test.go << 'EOF'
package engine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateEngine_RenderString(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		data     interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple variable",
			template: "Hello {{.Name}}",
			data:     map[string]string{"Name": "World"},
			want:     "Hello World",
			wantErr:  false,
		},
		{
			name:     "with helper function",
			template: "{{pascalCase .name}}",
			data:     map[string]string{"name": "user_profile"},
			want:     "UserProfile",
			wantErr:  false,
		},
		{
			name:     "invalid template",
			template: "{{.Missing",
			data:     map[string]string{},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestTemplateEngine_RenderFile(t *testing.T) {
	engine := NewTemplateEngine()

	// Create a temporary template
	tmpl := "package {{.Package}}\n\nconst Version = \"{{.Version}}\""
	data := map[string]string{
		"Package": "main",
		"Version": "1.0.0",
	}

	result, err := engine.RenderString(tmpl, data)
	require.NoError(t, err)
	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "const Version = \"1.0.0\"")
}
EOF
```

---

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/codegen/engine/... -v
```

Expected: FAIL - TemplateEngine not defined

---

- [ ] **Step 3: Implement template engine**

```bash
cat > internal/codegen/engine/template_engine.go << 'EOF'
package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateEngine handles template rendering with custom functions
type TemplateEngine struct {
	funcMap template.FuncMap
}

// NewTemplateEngine creates a new template engine with helper functions
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		funcMap: GetHelperFunctions(),
	}
}

// RenderString renders a template string with the given data
func (e *TemplateEngine) RenderString(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Funcs(e.funcMap).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderFile renders a template file with the given data
func (e *TemplateEngine) RenderFile(templatePath string, data interface{}) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return e.RenderString(string(content), data)
}

// RenderFiles renders multiple template files from a directory
func (e *TemplateEngine) RenderFiles(templateDir string, data interface{}) (map[string]string, error) {
	results := make(map[string]string)

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() || filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Render the template
		content, err := e.RenderFile(path, data)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", path, err)
		}

		// Store with relative path (without .tmpl extension)
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}
		outputPath := relPath[:len(relPath)-5] // Remove .tmpl extension

		results[outputPath] = content
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
EOF
```

---

- [ ] **Step 4: Implement helper functions**

```bash
cat > internal/codegen/engine/helpers.go << 'EOF'
package engine

import (
	"strings"
	"text/template"
	"unicode"
)

// GetHelperFunctions returns template helper functions
func GetHelperFunctions() template.FuncMap {
	return template.FuncMap{
		"pascalCase": pascalCase,
		"camelCase":  camelCase,
		"snakeCase":  snakeCase,
		"kebabCase":  kebabCase,
		"upper":      strings.ToUpper,
		"lower":      strings.Lower,
		"title":      strings.Title,
		"plural":     plural,
		"singular":   singular,
		"join":       strings.Join,
		"split":      strings.Split,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
	}
}

// pascalCase converts a string to PascalCase
func pascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Split by common delimiters
	parts := splitByDelimiters(s)
	
	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Capitalize first letter, lowercase rest
		runes := []rune(part)
		result.WriteRune(unicode.ToUpper(runes[0]))
		result.WriteString(strings.ToLower(string(runes[1:])))
	}

	return result.String()
}

// camelCase converts a string to camelCase
func camelCase(s string) string {
	pascal := pascalCase(s)
	if pascal == "" {
		return ""
	}

	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// snakeCase converts a string to snake_case
func snakeCase(s string) string {
	parts := splitByDelimiters(s)
	return strings.ToLower(strings.Join(parts, "_"))
}

// kebabCase converts a string to kebab-case
func kebabCase(s string) string {
	parts := splitByDelimiters(s)
	return strings.ToLower(strings.Join(parts, "-"))
}

// splitByDelimiters splits a string by common delimiters
func splitByDelimiters(s string) []string {
	// Replace common delimiters with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, ".", " ")

	// Split by spaces and capital letters
	var parts []string
	var current strings.Builder

	for i, r := range s {
		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else if unicode.IsUpper(r) && i > 0 {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			current.WriteRune(r)
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// plural returns the plural form of a word (simple implementation)
func plural(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	return s + "s"
}

// singular returns the singular form of a word (simple implementation)
func singular(s string) string {
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "es") {
		if strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "ches") || strings.HasSuffix(s, "shes") {
			return s[:len(s)-2]
		}
		return s[:len(s)-1]
	}
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}
EOF
```

---

- [ ] **Step 5: Add test for helper functions**

```bash
cat >> internal/codegen/engine/template_engine_test.go << 'EOF'

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		helper   string
		input    string
		expected string
	}{
		{"pascalCase simple", "pascalCase", "user_profile", "UserProfile"},
		{"pascalCase mixed", "pascalCase", "user-profile-data", "UserProfileData"},
		{"camelCase simple", "camelCase", "user_profile", "userProfile"},
		{"snakeCase simple", "snakeCase", "UserProfile", "user_profile"},
		{"kebabCase simple", "kebabCase", "UserProfile", "user-profile"},
		{"plural regular", "plural", "user", "users"},
		{"plural y", "plural", "category", "categories"},
		{"plural s", "plural", "class", "classes"},
		{"singular regular", "singular", "users", "user"},
		{"singular ies", "singular", "categories", "category"},
	}

	engine := NewTemplateEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := fmt.Sprintf("{{%s .input}}", tt.helper)
			data := map[string]string{"input": tt.input}
			
			result, err := engine.RenderString(tmpl, data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
EOF
```

---

- [ ] **Step 6: Run all tests**

```bash
go test ./internal/codegen/engine/... -v
```

Expected: PASS - all template engine tests passing

---

- [ ] **Step 7: Commit template engine**

```bash
git add internal/codegen/engine/
git commit -m "feat(codegen): add template engine with helper functions"
```

---

## Task 3: File Writer

**Files:**
- Create: `internal/codegen/engine/file_writer.go`
- Test: `internal/codegen/engine/file_writer_test.go`

---

- [ ] **Step 1: Write test for file writer**

```bash
cat > internal/codegen/engine/file_writer_test.go << 'EOF'
package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/internal/codegen"
)

func TestFileWriter_WriteFiles(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	files := []*codegen.GeneratedFile{
		{
			Path:    "src/index.ts",
			Content: "console.log('Hello');",
			Mode:    0644,
		},
		{
			Path:    "package.json",
			Content: `{"name": "test"}`,
			Mode:    0644,
		},
	}

	err = writer.WriteFiles(files)
	require.NoError(t, err)

	// Verify files exist
	indexPath := filepath.Join(tmpDir, "src", "index.ts")
	assert.FileExists(t, indexPath)

	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	assert.Equal(t, "console.log('Hello');", string(content))
}

func TestFileWriter_EnsureDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	writer := NewFileWriter(tmpDir)

	deepPath := "src/controllers/v1/users"
	err = writer.EnsureDir(deepPath)
	require.NoError(t, err)

	fullPath := filepath.Join(tmpDir, deepPath)
	info, err := os.Stat(fullPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
EOF
```

---

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/codegen/engine/... -v -run TestFileWriter
```

Expected: FAIL - FileWriter not defined

---

- [ ] **Step 3: Implement file writer**

```bash
cat > internal/codegen/engine/file_writer.go << 'EOF'
package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ktarun.reddy/baas/internal/codegen"
)

// FileWriter handles writing generated files to disk
type FileWriter struct {
	outputDir string
}

// NewFileWriter creates a new file writer
func NewFileWriter(outputDir string) *FileWriter {
	return &FileWriter{
		outputDir: outputDir,
	}
}

// WriteFiles writes all generated files to disk
func (w *FileWriter) WriteFiles(files []*codegen.GeneratedFile) error {
	for _, file := range files {
		if err := w.WriteFile(file); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.Path, err)
		}
	}
	return nil
}

// WriteFile writes a single file to disk
func (w *FileWriter) WriteFile(file *codegen.GeneratedFile) error {
	fullPath := filepath.Join(w.outputDir, file.Path)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file with specified mode
	mode := os.FileMode(file.Mode)
	if mode == 0 {
		mode = 0644 // Default mode
	}

	if err := os.WriteFile(fullPath, []byte(file.Content), mode); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// EnsureDir ensures a directory exists
func (w *FileWriter) EnsureDir(path string) error {
	fullPath := filepath.Join(w.outputDir, path)
	return os.MkdirAll(fullPath, 0755)
}

// Remove removes a file or directory
func (w *FileWriter) Remove(path string) error {
	fullPath := filepath.Join(w.outputDir, path)
	return os.RemoveAll(fullPath)
}

// Exists checks if a file or directory exists
func (w *FileWriter) Exists(path string) bool {
	fullPath := filepath.Join(w.outputDir, path)
	_, err := os.Stat(fullPath)
	return err == nil
}
EOF
```

---

- [ ] **Step 4: Run tests**

```bash
go test ./internal/codegen/engine/... -v
```

Expected: PASS - all file writer tests passing

---

- [ ] **Step 5: Commit file writer**

```bash
git add internal/codegen/engine/file_writer.go internal/codegen/engine/file_writer_test.go
git commit -m "feat(codegen): add file writer for generated code output"
```

---

## Task 4: TypeScript Base Templates

**Files:**
- Create: `templates/typescript/base/package.json.tmpl`
- Create: `templates/typescript/base/tsconfig.json.tmpl`
- Create: `templates/typescript/base/server.ts.tmpl`
- Create: `templates/typescript/base/database.ts.tmpl`
- Create: `templates/typescript/base/.env.example.tmpl`
- Create: `templates/typescript/base/.gitignore.tmpl`
- Create: `templates/typescript/base/README.md.tmpl`

---

- [ ] **Step 1: Create package.json template**

```bash
mkdir -p templates/typescript/base
cat > templates/typescript/base/package.json.tmpl << 'EOF'
{
  "name": "{{kebabCase .ProjectName}}",
  "version": "1.0.0",
  "description": "Generated by Aurora BaaS",
  "main": "dist/server.js",
  "scripts": {
    "dev": "ts-node-dev --respawn --transpile-only src/server.ts",
    "build": "tsc",
    "start": "node dist/server.js",
    "test": "jest",
    "lint": "eslint src --ext .ts",
    "format": "prettier --write \"src/**/*.ts\"",
    "prisma:generate": "prisma generate",
    "prisma:migrate": "prisma migrate dev",
    "prisma:studio": "prisma studio"
  },
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^7.1.0",
    "morgan": "^1.10.0",
    "dotenv": "^16.3.1",
    "@prisma/client": "^5.7.0",
    "zod": "^3.22.4"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/cors": "^2.8.17",
    "@types/morgan": "^1.9.9",
    "@types/node": "^20.10.6",
    "typescript": "^5.3.3",
    "ts-node-dev": "^2.0.0",
    "prisma": "^5.7.0",
    "jest": "^29.7.0",
    "@types/jest": "^29.5.11",
    "ts-jest": "^29.1.1",
    "eslint": "^8.56.0",
    "@typescript-eslint/parser": "^6.17.0",
    "@typescript-eslint/eslint-plugin": "^6.17.0",
    "prettier": "^3.1.1"
  },
  "engines": {
    "node": ">=18.0.0"
  }
}
EOF
```

---

- [ ] **Step 2: Create tsconfig.json template**

```bash
cat > templates/typescript/base/tsconfig.json.tmpl << 'EOF'
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "moduleResolution": "node",
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
EOF
```

---

- [ ] **Step 3: Create server.ts template**

```bash
cat > templates/typescript/base/server.ts.tmpl << 'EOF'
import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import dotenv from 'dotenv';
import { connectDatabase } from './database';
{{range .Schema.Tables}}
import { {{camelCase .Name}}Router } from './routes/{{snakeCase .Name}}.generated';
{{end}}

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());
app.use(morgan('combined'));

// Health check
app.get('/health', (req, res) => {
  res.json({ 
    status: 'healthy',
    timestamp: new Date().toISOString(),
    service: '{{.ProjectName}}'
  });
});

// API Routes
{{range .Schema.Tables}}
app.use('/api/v1/{{kebabCase (plural .Name)}}', {{camelCase .Name}}Router);
{{end}}

// 404 handler
app.use((req, res) => {
  res.status(404).json({ error: 'Not found' });
});

// Error handler
app.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
  console.error('Error:', err);
  res.status(500).json({ 
    error: 'Internal server error',
    message: process.env.NODE_ENV === 'development' ? err.message : undefined
  });
});

// Start server
const startServer = async () => {
  try {
    // Connect to database
    await connectDatabase();
    
    app.listen(PORT, () => {
      console.log(`🚀 Server running on port ${PORT}`);
      console.log(`📝 Health check: http://localhost:${PORT}/health`);
      {{range .Schema.Tables}}
      console.log(`   {{pascalCase .Name}}: http://localhost:${PORT}/api/v1/{{kebabCase (plural .Name)}}`);
      {{end}}
    });
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
};

startServer();
EOF
```

---

- [ ] **Step 4: Create database.ts template**

```bash
cat > templates/typescript/base/database.ts.tmpl << 'EOF'
import { PrismaClient } from '@prisma/client';

let prisma: PrismaClient;

export const connectDatabase = async (): Promise<void> => {
  try {
    prisma = new PrismaClient({
      log: process.env.NODE_ENV === 'development' ? ['query', 'error', 'warn'] : ['error'],
    });

    await prisma.$connect();
    console.log('✅ Database connected successfully');
  } catch (error) {
    console.error('❌ Database connection failed:', error);
    throw error;
  }
};

export const disconnectDatabase = async (): Promise<void> => {
  if (prisma) {
    await prisma.$disconnect();
    console.log('Database disconnected');
  }
};

export const getPrismaClient = (): PrismaClient => {
  if (!prisma) {
    throw new Error('Database not connected. Call connectDatabase() first.');
  }
  return prisma;
};

// Graceful shutdown
process.on('SIGINT', async () => {
  await disconnectDatabase();
  process.exit(0);
});

process.on('SIGTERM', async () => {
  await disconnectDatabase();
  process.exit(0);
});
EOF
```

---

- [ ] **Step 5: Create .env.example template**

```bash
cat > templates/typescript/base/.env.example.tmpl << 'EOF'
# Server Configuration
PORT=3000
NODE_ENV=development

# Database Configuration
{{if eq .Database "postgres"}}
DATABASE_URL="postgresql://user:password@localhost:5432/{{snakeCase .ProjectName}}?schema=public"
{{else if eq .Database "mongodb"}}
DATABASE_URL="mongodb://localhost:27017/{{snakeCase .ProjectName}}"
{{end}}

# CORS Configuration
CORS_ORIGIN=http://localhost:3000

# Logging
LOG_LEVEL=info
EOF
```

---

- [ ] **Step 6: Create .gitignore template**

```bash
cat > templates/typescript/base/.gitignore.tmpl << 'EOF'
# Dependencies
node_modules/
package-lock.json
yarn.lock

# Build output
dist/
build/

# Environment
.env
.env.local
.env.*.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
logs/
*.log
npm-debug.log*

# Testing
coverage/
.nyc_output/

# Prisma
prisma/migrations/*
!prisma/migrations/.gitkeep
EOF
```

---

- [ ] **Step 7: Create README template**

```bash
cat > templates/typescript/base/README.md.tmpl << 'EOF'
# {{.ProjectName}}

Generated by Aurora BaaS

## Getting Started

### Prerequisites

- Node.js >= 18
- {{if eq .Database "postgres"}}PostgreSQL{{else if eq .Database "mongodb"}}MongoDB{{end}}

### Installation

```bash
npm install
```

### Configuration

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

### Database Setup

```bash
# Generate Prisma client
npm run prisma:generate

# Run migrations
npm run prisma:migrate
```

### Development

```bash
npm run dev
```

The server will start on http://localhost:3000

### API Endpoints

{{range .Schema.Tables}}
**{{pascalCase .Name}}**
- `GET /api/v1/{{kebabCase (plural .Name)}}` - List all {{plural .Name}}
- `GET /api/v1/{{kebabCase (plural .Name)}}/:id` - Get {{singular .Name}} by ID
- `POST /api/v1/{{kebabCase (plural .Name)}}` - Create new {{singular .Name}}
- `PUT /api/v1/{{kebabCase (plural .Name)}}/:id` - Update {{singular .Name}}
- `DELETE /api/v1/{{kebabCase (plural .Name)}}/:id` - Delete {{singular .Name}}

{{end}}

### Testing

```bash
npm test
```

### Building for Production

```bash
npm run build
npm start
```

## Project Structure

```
src/
├── server.ts              # Entry point
├── database.ts            # Database connection
├── routes/                # API routes (generated)
├── controllers/           # Request handlers (generated)
├── models/                # Type definitions (generated)
├── extensions/            # Custom code (user-editable)
│   ├── custom-routes.ts
│   └── hooks.ts
└── prisma/
    └── schema.prisma      # Database schema
```

## Generated vs Extension Files

**Generated Files** (auto-regenerated, do not edit):
- `routes/*.generated.ts`
- `controllers/*.generated.ts`
- `models/*.generated.ts`

**Extension Files** (safe to edit):
- `extensions/custom-routes.ts`
- `extensions/hooks.ts`
- Any file without `.generated.` in the name

## License

Generated by Aurora BaaS - https://aurora.io
EOF
```

---

- [ ] **Step 8: Commit base templates**

```bash
git add templates/typescript/base/
git commit -m "feat(codegen): add TypeScript base project templates"
```

---

## Task 5: Prisma Schema Generator

**Files:**
- Create: `internal/codegen/typescript/prisma.go`
- Create: `templates/typescript/prisma/schema.prisma.tmpl`
- Test: `internal/codegen/typescript/prisma_test.go`

---

- [ ] **Step 1: Write test for Prisma schema generation**

```bash
mkdir -p internal/codegen/typescript
cat > internal/codegen/typescript/prisma_test.go << 'EOF'
package typescript

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

func TestGeneratePrismaSchema(t *testing.T) {
	schema := &validator.SchemaDefinition{
		Version:  "1.0",
		Database: "postgres",
		Tables: []validator.Table{
			{
				Name: "users",
				Columns: []validator.Column{
					{Name: "id", Type: "uuid", PrimaryKey: true},
					{Name: "email", Type: "string", Unique: true},
					{Name: "name", Type: "string"},
					{Name: "created_at", Type: "timestamp"},
				},
			},
			{
				Name: "posts",
				Columns: []validator.Column{
					{Name: "id", Type: "uuid", PrimaryKey: true},
					{Name: "user_id", Type: "uuid", ForeignKey: "users.id"},
					{Name: "title", Type: "string"},
					{Name: "content", Type: "text"},
				},
			},
		},
	}

	result, err := GeneratePrismaSchema(schema)
	require.NoError(t, err)

	// Check for datasource
	assert.Contains(t, result, "datasource db")
	assert.Contains(t, result, "provider = \"postgresql\"")

	// Check for generator
	assert.Contains(t, result, "generator client")

	// Check for models
	assert.Contains(t, result, "model User")
	assert.Contains(t, result, "model Post")

	// Check for fields
	assert.Contains(t, result, "id       String")
	assert.Contains(t, result, "email    String")
	assert.Contains(t, result, "userId   String")

	// Check for relations
	assert.Contains(t, result, "@relation")
}

func TestMapPrismaType(t *testing.T) {
	tests := []struct {
		auroraType string
		prismaType string
	}{
		{"string", "String"},
		{"text", "String"},
		{"integer", "Int"},
		{"bigint", "BigInt"},
		{"float", "Float"},
		{"decimal", "Decimal"},
		{"boolean", "Boolean"},
		{"timestamp", "DateTime"},
		{"date", "DateTime"},
		{"uuid", "String"},
		{"json", "Json"},
	}

	for _, tt := range tests {
		t.Run(tt.auroraType, func(t *testing.T) {
			result := mapPrismaType(tt.auroraType)
			assert.Equal(t, tt.prismaType, result)
		})
	}
}
EOF
```

---

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/codegen/typescript/... -v
```

Expected: FAIL - GeneratePrismaSchema not defined

---

- [ ] **Step 3: Implement Prisma schema generator**

```bash
cat > internal/codegen/typescript/prisma.go << 'EOF'
package typescript

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/ktarun.reddy/baas/internal/codegen/engine"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// GeneratePrismaSchema generates a Prisma schema from Aurora schema definition
func GeneratePrismaSchema(schema *validator.SchemaDefinition) (string, error) {
	tmpl := `// Generated by Aurora BaaS
// Do not edit this file manually

generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "{{.Provider}}"
  url      = env("DATABASE_URL")
}

{{range .Models}}
model {{.Name}} {
{{range .Fields}}  {{.Name}} {{.Type}}{{if .Modifier}} {{.Modifier}}{{end}}{{if .Attributes}} {{.Attributes}}{{end}}
{{end}}{{if .Relations}}
{{range .Relations}}  {{.Name}} {{.Type}}{{if .Modifier}} {{.Modifier}}{{end}} {{.Attributes}}
{{end}}{{end}}}

{{end}}`

	data := struct {
		Provider string
		Models   []PrismaModel
	}{
		Provider: mapDatabaseProvider(schema.Database),
		Models:   convertToPrismaModels(schema.Tables),
	}

	t, err := template.New("prisma").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// PrismaModel represents a Prisma model
type PrismaModel struct {
	Name      string
	Fields    []PrismaField
	Relations []PrismaField
}

// PrismaField represents a field in a Prisma model
type PrismaField struct {
	Name       string
	Type       string
	Modifier   string // ?, []
	Attributes string // @id, @unique, @default(), etc.
}

// convertToPrismaModels converts Aurora tables to Prisma models
func convertToPrismaModels(tables []validator.Table) []PrismaModel {
	models := make([]PrismaModel, 0, len(tables))

	for _, table := range tables {
		model := PrismaModel{
			Name:      toPascalCase(table.Name),
			Fields:    []PrismaField{},
			Relations: []PrismaField{},
		}

		// Convert columns to fields
		for _, col := range table.Columns {
			field := PrismaField{
				Name: toCamelCase(col.Name),
				Type: mapPrismaType(col.Type),
			}

			// Build attributes
			var attrs []string

			if col.PrimaryKey {
				attrs = append(attrs, "@id")
				if col.Type == "uuid" {
					attrs = append(attrs, "@default(uuid())")
				} else if col.Type == "integer" || col.Type == "bigint" {
					attrs = append(attrs, "@default(autoincrement())")
				}
			}

			if col.Unique {
				attrs = append(attrs, "@unique")
			}

			if col.Default != "" && !col.PrimaryKey {
				if col.Type == "timestamp" || col.Type == "date" {
					attrs = append(attrs, "@default(now())")
				} else if col.Type == "string" || col.Type == "text" {
					attrs = append(attrs, fmt.Sprintf("@default(\"%s\")", col.Default))
				} else {
					attrs = append(attrs, fmt.Sprintf("@default(%s)", col.Default))
				}
			}

			if !col.Nullable && !col.PrimaryKey {
				// Prisma fields are required by default
			} else if col.Nullable {
				field.Modifier = "?"
			}

			field.Attributes = strings.Join(attrs, " ")

			// Check if this is a foreign key
			if col.ForeignKey != "" {
				// This will be handled as a relation
				field.Attributes += fmt.Sprintf(" @relation(fields: [%s], references: [id])", field.Name)
				model.Fields = append(model.Fields, field)

				// Add the relation field
				parts := strings.Split(col.ForeignKey, ".")
				if len(parts) == 2 {
					relatedTable := parts[0]
					relationField := PrismaField{
						Name: toCamelCase(singularize(relatedTable)),
						Type: toPascalCase(singularize(relatedTable)),
					}
					model.Relations = append(model.Relations, relationField)
				}
			} else {
				model.Fields = append(model.Fields, field)
			}
		}

		models = append(models, model)
	}

	return models
}

// mapDatabaseProvider maps Aurora database type to Prisma provider
func mapDatabaseProvider(dbType string) string {
	switch dbType {
	case "postgres":
		return "postgresql"
	case "mongodb":
		return "mongodb"
	case "mysql":
		return "mysql"
	case "sqlite":
		return "sqlite"
	default:
		return "postgresql"
	}
}

// mapPrismaType maps Aurora column type to Prisma type
func mapPrismaType(auroraType string) string {
	switch auroraType {
	case "string":
		return "String"
	case "text":
		return "String"
	case "integer":
		return "Int"
	case "bigint":
		return "BigInt"
	case "float":
		return "Float"
	case "decimal":
		return "Decimal"
	case "boolean":
		return "Boolean"
	case "timestamp":
		return "DateTime"
	case "date":
		return "DateTime"
	case "uuid":
		return "String"
	case "json":
		return "Json"
	default:
		return "String"
	}
}

// toPascalCase converts string to PascalCase
func toPascalCase(s string) string {
	engine := engine.NewTemplateEngine()
	result, _ := engine.RenderString("{{pascalCase .}}", s)
	return result
}

// toCamelCase converts string to camelCase
func toCamelCase(s string) string {
	engine := engine.NewTemplateEngine()
	result, _ := engine.RenderString("{{camelCase .}}", s)
	return result
}

// singularize converts plural to singular (simple implementation)
func singularize(s string) string {
	engine := engine.NewTemplateEngine()
	result, _ := engine.RenderString("{{singular .}}", s)
	return result
}
EOF
```

---

- [ ] **Step 4: Run tests**

```bash
go test ./internal/codegen/typescript/... -v
```

Expected: PASS

---

- [ ] **Step 5: Create Prisma schema template directory**

```bash
mkdir -p templates/typescript/prisma
cat > templates/typescript/prisma/schema.prisma.tmpl << 'EOF'
// This is a placeholder - schema is generated programmatically
// See internal/codegen/typescript/prisma.go for implementation
EOF
```

---

- [ ] **Step 6: Commit Prisma generator**

```bash
git add internal/codegen/typescript/prisma.go internal/codegen/typescript/prisma_test.go templates/typescript/prisma/
git commit -m "feat(codegen): add Prisma schema generator for TypeScript"
```

---

## Task 6: TypeScript CRUD Templates

**Files:**
- Create: `templates/typescript/crud/controller.generated.ts.tmpl`
- Create: `templates/typescript/crud/routes.generated.ts.tmpl`
- Create: `templates/typescript/crud/model.generated.ts.tmpl`

---

- [ ] **Step 1: Create controller template**

```bash
mkdir -p templates/typescript/crud
cat > templates/typescript/crud/controller.generated.ts.tmpl << 'EOF'
// AUTO-GENERATED - DO NOT EDIT
// Generated by Aurora BaaS
// To customize behavior, use hooks in extensions/{{snakeCase .TableName}}-hooks.ts

import { Request, Response, NextFunction } from 'express';
import { getPrismaClient } from '../database';
import { {{pascalCase .TableName}}CreateInput, {{pascalCase .TableName}}UpdateInput, validate{{pascalCase .TableName}}Create, validate{{pascalCase .TableName}}Update } from '../models/{{snakeCase .TableName}}.generated';
import { get{{pascalCase .TableName}}Hooks } from '../extensions/{{snakeCase .TableName}}-hooks';

const prisma = getPrismaClient();
const hooks = get{{pascalCase .TableName}}Hooks();

/**
 * Get all {{plural .TableName}}
 * GET /api/v1/{{kebabCase (plural .TableName)}}
 */
export const getAll{{pascalCase (plural .TableName)}} = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    // Pagination
    const page = parseInt(req.query.page as string) || 1;
    const limit = parseInt(req.query.limit as string) || 10;
    const skip = (page - 1) * limit;

    // Execute hook (if provided)
    if (hooks.beforeList) {
      await hooks.beforeList(req.query);
    }

    const [items, total] = await Promise.all([
      prisma.{{camelCase .TableName}}.findMany({
        skip,
        take: limit,
        orderBy: { createdAt: 'desc' },
      }),
      prisma.{{camelCase .TableName}}.count(),
    ]);

    // Execute hook (if provided)
    let result = items;
    if (hooks.afterList) {
      result = await hooks.afterList(items);
    }

    res.json({
      data: result,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit),
      },
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get {{singular .TableName}} by ID
 * GET /api/v1/{{kebabCase (plural .TableName)}}/:id
 */
export const get{{pascalCase .TableName}}ById = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const { id } = req.params;

    // Execute hook (if provided)
    if (hooks.beforeGet) {
      await hooks.beforeGet(id);
    }

    const item = await prisma.{{camelCase .TableName}}.findUnique({
      where: { id },
    });

    if (!item) {
      res.status(404).json({ error: '{{pascalCase .TableName}} not found' });
      return;
    }

    // Execute hook (if provided)
    let result = item;
    if (hooks.afterGet) {
      result = await hooks.afterGet(item);
    }

    res.json(result);
  } catch (error) {
    next(error);
  }
};

/**
 * Create new {{singular .TableName}}
 * POST /api/v1/{{kebabCase (plural .TableName)}}
 */
export const create{{pascalCase .TableName}} = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    // Validate input
    const validation = validate{{pascalCase .TableName}}Create(req.body);
    if (!validation.success) {
      res.status(400).json({ 
        error: 'Validation failed',
        details: validation.error.errors 
      });
      return;
    }

    let data: {{pascalCase .TableName}}CreateInput = validation.data;

    // Execute hook (if provided)
    if (hooks.beforeCreate) {
      data = await hooks.beforeCreate(data);
    }

    const item = await prisma.{{camelCase .TableName}}.create({
      data,
    });

    // Execute hook (if provided)
    if (hooks.afterCreate) {
      await hooks.afterCreate(item);
    }

    res.status(201).json(item);
  } catch (error) {
    next(error);
  }
};

/**
 * Update {{singular .TableName}}
 * PUT /api/v1/{{kebabCase (plural .TableName)}}/:id
 */
export const update{{pascalCase .TableName}} = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const { id } = req.params;

    // Check if exists
    const existing = await prisma.{{camelCase .TableName}}.findUnique({
      where: { id },
    });

    if (!existing) {
      res.status(404).json({ error: '{{pascalCase .TableName}} not found' });
      return;
    }

    // Validate input
    const validation = validate{{pascalCase .TableName}}Update(req.body);
    if (!validation.success) {
      res.status(400).json({ 
        error: 'Validation failed',
        details: validation.error.errors 
      });
      return;
    }

    let data: {{pascalCase .TableName}}UpdateInput = validation.data;

    // Execute hook (if provided)
    if (hooks.beforeUpdate) {
      data = await hooks.beforeUpdate(id, data);
    }

    const item = await prisma.{{camelCase .TableName}}.update({
      where: { id },
      data,
    });

    // Execute hook (if provided)
    if (hooks.afterUpdate) {
      await hooks.afterUpdate(item);
    }

    res.json(item);
  } catch (error) {
    next(error);
  }
};

/**
 * Delete {{singular .TableName}}
 * DELETE /api/v1/{{kebabCase (plural .TableName)}}/:id
 */
export const delete{{pascalCase .TableName}} = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const { id } = req.params;

    // Check if exists
    const existing = await prisma.{{camelCase .TableName}}.findUnique({
      where: { id },
    });

    if (!existing) {
      res.status(404).json({ error: '{{pascalCase .TableName}} not found' });
      return;
    }

    // Execute hook (if provided)
    if (hooks.beforeDelete) {
      await hooks.beforeDelete(id);
    }

    await prisma.{{camelCase .TableName}}.delete({
      where: { id },
    });

    // Execute hook (if provided)
    if (hooks.afterDelete) {
      await hooks.afterDelete(id);
    }

    res.json({ message: '{{pascalCase .TableName}} deleted successfully' });
  } catch (error) {
    next(error);
  }
};
EOF
```

---

- [ ] **Step 2: Create routes template**

```bash
cat > templates/typescript/crud/routes.generated.ts.tmpl << 'EOF'
// AUTO-GENERATED - DO NOT EDIT
// Generated by Aurora BaaS

import { Router } from 'express';
import {
  getAll{{pascalCase (plural .TableName)}},
  get{{pascalCase .TableName}}ById,
  create{{pascalCase .TableName}},
  update{{pascalCase .TableName}},
  delete{{pascalCase .TableName}},
} from '../controllers/{{snakeCase .TableName}}.generated';

export const {{camelCase .TableName}}Router = Router();

// List all {{plural .TableName}}
{{camelCase .TableName}}Router.get('/', getAll{{pascalCase (plural .TableName)}});

// Get {{singular .TableName}} by ID
{{camelCase .TableName}}Router.get('/:id', get{{pascalCase .TableName}}ById);

// Create new {{singular .TableName}}
{{camelCase .TableName}}Router.post('/', create{{pascalCase .TableName}});

// Update {{singular .TableName}}
{{camelCase .TableName}}Router.put('/:id', update{{pascalCase .TableName}});

// Delete {{singular .TableName}}
{{camelCase .TableName}}Router.delete('/:id', delete{{pascalCase .TableName}});
EOF
```

---

- [ ] **Step 3: Create model template**

```bash
cat > templates/typescript/crud/model.generated.ts.tmpl << 'EOF'
// AUTO-GENERATED - DO NOT EDIT
// Generated by Aurora BaaS

import { z } from 'zod';

/**
 * {{pascalCase .TableName}} type definition
 */
export interface {{pascalCase .TableName}} {
{{range .Columns}}  {{camelCase .Name}}{{if .Nullable}}?{{end}}: {{if eq .Type "string"}}string{{else if eq .Type "text"}}string{{else if eq .Type "integer"}}number{{else if eq .Type "bigint"}}number{{else if eq .Type "float"}}number{{else if eq .Type "decimal"}}number{{else if eq .Type "boolean"}}boolean{{else if eq .Type "timestamp"}}Date{{else if eq .Type "date"}}Date{{else if eq .Type "uuid"}}string{{else if eq .Type "json"}}any{{else}}any{{end}};
{{end}}}

/**
 * Create input schema
 */
export const {{pascalCase .TableName}}CreateSchema = z.object({
{{range .Columns}}{{if not .PrimaryKey}}  {{camelCase .Name}}: {{if eq .Type "string"}}z.string(){{else if eq .Type "text"}}z.string(){{else if eq .Type "integer"}}z.number().int(){{else if eq .Type "bigint"}}z.number().int(){{else if eq .Type "float"}}z.number(){{else if eq .Type "decimal"}}z.number(){{else if eq .Type "boolean"}}z.boolean(){{else if eq .Type "timestamp"}}z.date(){{else if eq .Type "date"}}z.date(){{else if eq .Type "uuid"}}z.string().uuid(){{else if eq .Type "json"}}z.any(){{else}}z.any(){{end}}{{if .Nullable}}.optional(){{end}}{{if and (eq .Type "string") (not .Nullable)}}.min(1){{end}},
{{end}}{{end}}});

export type {{pascalCase .TableName}}CreateInput = z.infer<typeof {{pascalCase .TableName}}CreateSchema>;

/**
 * Update input schema (all fields optional)
 */
export const {{pascalCase .TableName}}UpdateSchema = z.object({
{{range .Columns}}{{if not .PrimaryKey}}  {{camelCase .Name}}: {{if eq .Type "string"}}z.string(){{else if eq .Type "text"}}z.string(){{else if eq .Type "integer"}}z.number().int(){{else if eq .Type "bigint"}}z.number().int(){{else if eq .Type "float"}}z.number(){{else if eq .Type "decimal"}}z.number(){{else if eq .Type "boolean"}}z.boolean(){{else if eq .Type "timestamp"}}z.date(){{else if eq .Type "date"}}z.date(){{else if eq .Type "uuid"}}z.string().uuid(){{else if eq .Type "json"}}z.any(){{else}}z.any(){{end}}.optional(),
{{end}}{{end}}});

export type {{pascalCase .TableName}}UpdateInput = z.infer<typeof {{pascalCase .TableName}}UpdateSchema>;

/**
 * Validation functions
 */
export const validate{{pascalCase .TableName}}Create = (data: unknown) => {
  return {{pascalCase .TableName}}CreateSchema.safeParse(data);
};

export const validate{{pascalCase .TableName}}Update = (data: unknown) => {
  return {{pascalCase .TableName}}UpdateSchema.safeParse(data);
};
EOF
```

---

- [ ] **Step 4: Commit CRUD templates**

```bash
git add templates/typescript/crud/
git commit -m "feat(codegen): add TypeScript CRUD templates with hooks support"
```

---

## Task 7: Extension Point Templates

**Files:**
- Create: `templates/typescript/extensions/custom-routes.ts.tmpl`
- Create: `templates/typescript/extensions/hooks.ts.tmpl`

---

- [ ] **Step 1: Create custom routes template**

```bash
mkdir -p templates/typescript/extensions
cat > templates/typescript/extensions/custom-routes.ts.tmpl << 'EOF'
/**
 * Custom Routes
 * 
 * This file is safe to edit - it will never be overwritten by Aurora.
 * Add your custom endpoints here.
 * 
 * Example:
 * 
 * import { Router } from 'express';
 * import { getPrismaClient } from '../database';
 * 
 * export const customRouter = Router();
 * const prisma = getPrismaClient();
 * 
 * customRouter.get('/custom-endpoint', async (req, res) => {
 *   // Your custom logic here
 *   res.json({ message: 'Custom endpoint' });
 * });
 */

import { Router } from 'express';

export const customRouter = Router();

// Add your custom routes here
EOF
```

---

- [ ] **Step 2: Create hooks template for each table**

```bash
cat > templates/typescript/extensions/hooks.ts.tmpl << 'EOF'
/**
 * {{pascalCase .TableName}} Hooks
 * 
 * This file is safe to edit - it will never be overwritten by Aurora.
 * Use these hooks to add custom logic before/after CRUD operations.
 * 
 * All hooks are optional. Only implement the ones you need.
 */

import { {{pascalCase .TableName}}, {{pascalCase .TableName}}CreateInput, {{pascalCase .TableName}}UpdateInput } from '../models/{{snakeCase .TableName}}.generated';

export interface {{pascalCase .TableName}}Hooks {
  // Called before listing {{plural .TableName}}
  // Modify query parameters or perform pre-processing
  beforeList?: (query: any) => Promise<any>;
  
  // Called after listing {{plural .TableName}}
  // Modify the results before sending to client
  afterList?: (items: {{pascalCase .TableName}}[]) => Promise<{{pascalCase .TableName}}[]>;
  
  // Called before getting a single {{singular .TableName}}
  // Perform authorization checks or validation
  beforeGet?: (id: string) => Promise<void>;
  
  // Called after getting a single {{singular .TableName}}
  // Modify the result before sending to client
  afterGet?: (item: {{pascalCase .TableName}}) => Promise<{{pascalCase .TableName}}>;
  
  // Called before creating a {{singular .TableName}}
  // Modify input data, hash passwords, etc.
  beforeCreate?: (data: {{pascalCase .TableName}}CreateInput) => Promise<{{pascalCase .TableName}}CreateInput>;
  
  // Called after creating a {{singular .TableName}}
  // Send notifications, create related records, etc.
  afterCreate?: (item: {{pascalCase .TableName}}) => Promise<void>;
  
  // Called before updating a {{singular .TableName}}
  // Modify input data or perform validation
  beforeUpdate?: (id: string, data: {{pascalCase .TableName}}UpdateInput) => Promise<{{pascalCase .TableName}}UpdateInput>;
  
  // Called after updating a {{singular .TableName}}
  // Send notifications, update related records, etc.
  afterUpdate?: (item: {{pascalCase .TableName}}) => Promise<void>;
  
  // Called before deleting a {{singular .TableName}}
  // Perform authorization or cascade deletes
  beforeDelete?: (id: string) => Promise<void>;
  
  // Called after deleting a {{singular .TableName}}
  // Clean up related resources
  afterDelete?: (id: string) => Promise<void>;
}

/**
 * Get {{pascalCase .TableName}} hooks
 * Implement the hooks you need below
 */
export const get{{pascalCase .TableName}}Hooks = (): {{pascalCase .TableName}}Hooks => {
  return {
    // Example: Hash password before creating a user
    // beforeCreate: async (data) => {
    //   if (data.password) {
    //     data.password = await hashPassword(data.password);
    //   }
    //   return data;
    // },
    
    // Example: Send welcome email after creating a user
    // afterCreate: async (item) => {
    //   await sendWelcomeEmail(item.email);
    // },
  };
};
EOF
```

---

- [ ] **Step 3: Commit extension templates**

```bash
git add templates/typescript/extensions/
git commit -m "feat(codegen): add extension point templates for custom code"
```

---

*Due to length constraints, I'll continue in the next response with the remaining tasks: TypeScript Generator Implementation, API Integration, Testing, and Execution Handoff.*

---

This plan continues with Tasks 8-12. Shall I continue with the remaining tasks?
