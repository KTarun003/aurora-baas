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
