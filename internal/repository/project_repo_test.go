package repository

import (
	"testing"

	"github.com/ktarun.reddy/baas/internal/config"
	"github.com/ktarun.reddy/baas/internal/database"
	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
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
	_, err = sqlDB.Exec("DELETE FROM schemas")
	require.NoError(t, err)
	_, err = sqlDB.Exec("DELETE FROM projects")
	require.NoError(t, err)

	return db
}

func TestNewProjectRepository(t *testing.T) {
	db := setupTestDB(t)

	repo := NewProjectRepository(db)
	assert.NotNil(t, repo)
}

func TestProjectRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	project := &domain.Project{
		Name:         "Test Project",
		Description:  "A test project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := repo.Create(project)
	require.NoError(t, err)
	assert.NotEmpty(t, project.ID)
	assert.NotZero(t, project.CreatedAt)
	assert.NotZero(t, project.UpdatedAt)
}

func TestProjectRepository_Create_WithoutOptionalFields(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	project := &domain.Project{
		Name:         "Minimal Project",
		Language:     "python",
		DatabaseType: "mongodb",
		APIStyle:     "graphql",
	}

	err := repo.Create(project)
	require.NoError(t, err)
	assert.NotEmpty(t, project.ID)
}

func TestProjectRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	// Create a project
	project := &domain.Project{
		Name:         "Test Project",
		Description:  "A test project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := repo.Create(project)
	require.NoError(t, err)

	// Find by ID
	found, err := repo.FindByID(project.ID)
	require.NoError(t, err)
	assert.Equal(t, project.ID, found.ID)
	assert.Equal(t, project.Name, found.Name)
	assert.Equal(t, project.Description, found.Description)
	assert.Equal(t, project.Language, found.Language)
	assert.Equal(t, project.DatabaseType, found.DatabaseType)
	assert.Equal(t, project.APIStyle, found.APIStyle)
}

func TestProjectRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	found, err := repo.FindByID("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Equal(t, "project not found", err.Error())
	assert.Nil(t, found)
}

func TestProjectRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	// Create multiple projects
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
			DatabaseType: "postgres",
			APIStyle:     "rest",
		},
	}

	for _, p := range projects {
		err := repo.Create(p)
		require.NoError(t, err)
	}

	// List all projects
	result, err := repo.List()
	require.NoError(t, err)
	assert.Len(t, result, 3)

	// Verify ordering by created_at DESC (most recent first)
	// The last created should be first in the list
	assert.Equal(t, "Project 3", result[0].Name)
	assert.Equal(t, "Project 2", result[1].Name)
	assert.Equal(t, "Project 1", result[2].Name)
}

func TestProjectRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	result, err := repo.List()
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestProjectRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	// Create a project
	project := &domain.Project{
		Name:         "Original Name",
		Description:  "Original Description",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := repo.Create(project)
	require.NoError(t, err)

	// Update the project
	project.Name = "Updated Name"
	project.Description = "Updated Description"
	err = repo.Update(project)
	require.NoError(t, err)

	// Verify the update
	found, err := repo.FindByID(project.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", found.Name)
	assert.Equal(t, "Updated Description", found.Description)
}

func TestProjectRepository_Update_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	project := &domain.Project{
		ID:           "00000000-0000-0000-0000-000000000000",
		Name:         "Non-existent",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := repo.Update(project)
	assert.Error(t, err)
}

func TestProjectRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	// Create a project
	project := &domain.Project{
		Name:         "Project to Delete",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := repo.Create(project)
	require.NoError(t, err)

	// Delete the project
	err = repo.Delete(project.ID)
	require.NoError(t, err)

	// Verify deletion
	found, err := repo.FindByID(project.ID)
	assert.Error(t, err)
	assert.Equal(t, "project not found", err.Error())
	assert.Nil(t, found)
}

func TestProjectRepository_Delete_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProjectRepository(db)

	err := repo.Delete("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}
