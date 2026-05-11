package repository

import (
	"errors"
	"fmt"

	"github.com/ktarun.reddy/baas/internal/domain"
	"gorm.io/gorm"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create inserts a new project into the database
func (r *ProjectRepository) Create(project *domain.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// FindByID retrieves a project by its ID
func (r *ProjectRepository) FindByID(id string) (*domain.Project, error) {
	var project domain.Project
	err := r.db.Where("id = ?", id).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, fmt.Errorf("failed to find project: %w", err)
	}
	return &project, nil
}

// List retrieves all projects ordered by created_at DESC
func (r *ProjectRepository) List() ([]*domain.Project, error) {
	var projects []*domain.Project
	err := r.db.Order("created_at DESC").Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	return projects, nil
}

// Update updates an existing project in the database
func (r *ProjectRepository) Update(project *domain.Project) error {
	result := r.db.Model(project).Where("id = ?", project.ID).Updates(project)
	if result.Error != nil {
		return fmt.Errorf("failed to update project: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}

// Delete removes a project from the database
func (r *ProjectRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&domain.Project{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}
