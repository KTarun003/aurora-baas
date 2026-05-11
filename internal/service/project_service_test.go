package service

import (
	"testing"

	"github.com/ktarun.reddy/baas/internal/config"
	"github.com/ktarun.reddy/baas/internal/database"
	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/ktarun.reddy/baas/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestService creates a test database and returns both services
func setupTestService(t *testing.T) (*ProjectService, *SchemaService) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "aurora",
		Password: "aurora_dev",
		Name:     "aurora_dev",
		SSLMode:  "disable",
	}

	db, err := database.Connect(cfg)
	require.NoError(t, err)

	err = database.RunMigrations(db)
	require.NoError(t, err)

	// Clean up existing test data
	sqlDB, _ := db.DB()
	sqlDB.Exec("DELETE FROM schemas")
	sqlDB.Exec("DELETE FROM projects")

	// Create repositories
	projectRepo := repository.NewProjectRepository(db)
	schemaRepo := repository.NewSchemaRepository(db)

	// Create services
	projectService := NewProjectService(projectRepo)
	schemaService := NewSchemaService(schemaRepo, projectRepo)

	return projectService, schemaService
}

// TestProjectService_CreateProject tests successful project creation
func TestProjectService_CreateProject(t *testing.T) {
	projectService, _ := setupTestService(t)

	project := &domain.Project{
		Name:         "Test Project",
		Description:  "A test project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)
	assert.NotEmpty(t, project.ID)
}

// TestProjectService_CreateProject_ValidationFails tests validation error on empty name
func TestProjectService_CreateProject_ValidationFails(t *testing.T) {
	projectService, _ := setupTestService(t)

	project := &domain.Project{
		Name:         "",
		Description:  "A test project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

// TestProjectService_GetProject tests project retrieval by ID
func TestProjectService_GetProject(t *testing.T) {
	projectService, _ := setupTestService(t)

	// Create a project first
	project := &domain.Project{
		Name:         "Test Project",
		Description:  "A test project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)

	// Retrieve it
	retrieved, err := projectService.GetProject(project.ID)
	require.NoError(t, err)
	require.Equal(t, project.ID, retrieved.ID)
	require.Equal(t, "Test Project", retrieved.Name)
}

// TestProjectService_ListProjects tests listing all projects
func TestProjectService_ListProjects(t *testing.T) {
	projectService, _ := setupTestService(t)

	// Create 3 projects
	projects := []*domain.Project{
		{
			Name:         "Project 1",
			Language:     "typescript",
			DatabaseType: "postgres",
			APIStyle:     "rest",
		},
		{
			Name:         "Project 2",
			Language:     "python",
			DatabaseType: "mongodb",
			APIStyle:     "graphql",
		},
		{
			Name:         "Project 3",
			Language:     "typescript",
			DatabaseType: "mongodb",
			APIStyle:     "rest",
		},
	}

	for _, p := range projects {
		err := projectService.CreateProject(p)
		require.NoError(t, err)
	}

	// List projects
	listed, err := projectService.ListProjects()
	require.NoError(t, err)
	assert.Len(t, listed, 3)
}

// TestProjectService_DeleteProject tests project deletion
func TestProjectService_DeleteProject(t *testing.T) {
	projectService, _ := setupTestService(t)

	// Create a project
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)

	// Delete it
	err = projectService.DeleteProject(project.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = projectService.GetProject(project.ID)
	assert.Error(t, err)
}

// TestProjectService_UpdateProject_ValidationFails tests validation error on empty name during update
func TestProjectService_UpdateProject_ValidationFails(t *testing.T) {
	projectService, _ := setupTestService(t)

	// Create a valid project first
	project := &domain.Project{
		Name:         "test-project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectService.CreateProject(project)
	require.NoError(t, err)

	// Try to update with invalid data
	project.Name = ""
	err = projectService.UpdateProject(project)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}
