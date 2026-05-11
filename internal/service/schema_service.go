package service

import (
	"fmt"

	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/ktarun.reddy/baas/internal/repository"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// SchemaService handles business logic for schemas
type SchemaService struct {
	schemaRepo   *repository.SchemaRepository
	projectRepo  *repository.ProjectRepository
}

// NewSchemaService creates a new schema service
func NewSchemaService(schemaRepo *repository.SchemaRepository, projectRepo *repository.ProjectRepository) *SchemaService {
	if schemaRepo == nil {
		panic("schema repository cannot be nil")
	}
	if projectRepo == nil {
		panic("project repository cannot be nil")
	}
	return &SchemaService{
		schemaRepo:  schemaRepo,
		projectRepo: projectRepo,
	}
}

// ApplySchema validates project existence, validates YAML, determines version, and creates schema
func (s *SchemaService) ApplySchema(projectID string, yamlContent string) (*domain.Schema, error) {
	// Check if project exists
	_, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Validate YAML schema using the validator
	_, err = validator.ValidateSchema(yamlContent)
	if err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// Determine next version
	nextVersion := 1
	latestSchema, err := s.schemaRepo.FindLatestByProjectID(projectID)
	if err == nil && latestSchema != nil {
		nextVersion = latestSchema.Version + 1
	}

	// Create the new schema
	schema := &domain.Schema{
		ProjectID: projectID,
		Content:   yamlContent,
		Version:   nextVersion,
	}

	// Validate the schema entity
	if err := schema.Validate(); err != nil {
		return nil, fmt.Errorf("schema entity validation failed: %w", err)
	}

	// Save to repository
	if err := s.schemaRepo.Create(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return schema, nil
}

// GetLatestSchema retrieves the latest schema version for a project
func (s *SchemaService) GetLatestSchema(projectID string) (*domain.Schema, error) {
	schema, err := s.schemaRepo.FindLatestByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest schema: %w", err)
	}

	return schema, nil
}

// ListSchemas lists all schema versions for a project
func (s *SchemaService) ListSchemas(projectID string) ([]*domain.Schema, error) {
	schemas, err := s.schemaRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}

	return schemas, nil
}
