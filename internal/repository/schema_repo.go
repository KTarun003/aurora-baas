package repository

import (
	"errors"
	"fmt"

	"github.com/ktarun.reddy/baas/internal/domain"
	"gorm.io/gorm"
)

// SchemaRepository handles database operations for schemas
type SchemaRepository struct {
	db *gorm.DB
}

// NewSchemaRepository creates a new schema repository
func NewSchemaRepository(db *gorm.DB) *SchemaRepository {
	return &SchemaRepository{db: db}
}

// Create inserts a new schema into the database
func (r *SchemaRepository) Create(schema *domain.Schema) error {
	if err := r.db.Create(schema).Error; err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

// FindByProjectID retrieves all schemas for a project ordered by version DESC
func (r *SchemaRepository) FindByProjectID(projectID string) ([]*domain.Schema, error) {
	var schemas []*domain.Schema
	err := r.db.Where("project_id = ?", projectID).Order("version DESC").Find(&schemas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find schemas: %w", err)
	}
	return schemas, nil
}

// FindLatestByProjectID retrieves the latest schema version for a project
func (r *SchemaRepository) FindLatestByProjectID(projectID string) (*domain.Schema, error) {
	var schema domain.Schema
	err := r.db.Where("project_id = ?", projectID).Order("version DESC").First(&schema).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schema not found")
		}
		return nil, fmt.Errorf("failed to find latest schema: %w", err)
	}
	return &schema, nil
}
