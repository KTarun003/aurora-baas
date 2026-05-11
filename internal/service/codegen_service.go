package service

import (
	"fmt"

	"github.com/ktarun.reddy/baas/internal/codegen"
	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// ProjectRepository interface for accessing project data
type ProjectRepository interface {
	FindByID(id string) (*domain.Project, error)
}

// SchemaRepository interface for accessing schema data
type SchemaRepository interface {
	FindLatestByProjectID(projectID string) (*domain.Schema, error)
}

// CodegenService handles business logic for code generation
type CodegenService struct {
	registry          *codegen.Registry
	projectRepo       ProjectRepository
	schemaRepo        SchemaRepository
}

// NewCodegenService creates a new codegen service
func NewCodegenService(
	registry *codegen.Registry,
	projectRepo ProjectRepository,
	schemaRepo SchemaRepository,
) *CodegenService {
	if registry == nil {
		panic("codegen registry cannot be nil")
	}
	if projectRepo == nil {
		panic("project repository cannot be nil")
	}
	if schemaRepo == nil {
		panic("schema repository cannot be nil")
	}

	return &CodegenService{
		registry:    registry,
		projectRepo: projectRepo,
		schemaRepo:  schemaRepo,
	}
}

// GenerateProject generates code for a project
func (s *CodegenService) GenerateProject(projectID, outputDir string) (*codegen.GeneratedProject, error) {
	// Get the project
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Get the latest schema for the project
	schema, err := s.schemaRepo.FindLatestByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %w", err)
	}

	// Parse the schema content
	parsedSchema, err := validator.ValidateSchema(schema.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Create generator context
	ctx := &codegen.GeneratorContext{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		Language:    project.Language,
		APIStyle:    project.APIStyle,
		Database:    project.DatabaseType,
		Schema:      parsedSchema,
		OutputDir:   outputDir,
	}

	// Generate the project using the registry
	generatedProject, err := s.registry.Generate(project.Language, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project: %w", err)
	}

	return generatedProject, nil
}
