package service

import (
	"fmt"

	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/ktarun.reddy/baas/internal/repository"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	repo *repository.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	if repo == nil {
		panic("project repository cannot be nil")
	}
	return &ProjectService{
		repo: repo,
	}
}

// CreateProject validates and creates a new project
func (s *ProjectService) CreateProject(project *domain.Project) error {
	// Validate the project
	if err := project.Validate(); err != nil {
		return err
	}

	// Create the project in the repository
	if err := s.repo.Create(project); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(id string) (*domain.Project, error) {
	return s.repo.FindByID(id)
}

// ListProjects retrieves all projects
func (s *ProjectService) ListProjects() ([]*domain.Project, error) {
	projects, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// UpdateProject validates and updates an existing project
func (s *ProjectService) UpdateProject(project *domain.Project) error {
	// Validate the project
	if err := project.Validate(); err != nil {
		return err
	}

	// Update the project in the repository
	if err := s.repo.Update(project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

// DeleteProject deletes a project by ID
func (s *ProjectService) DeleteProject(id string) error {
	return s.repo.Delete(id)
}
