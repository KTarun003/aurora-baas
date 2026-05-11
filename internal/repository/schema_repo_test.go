package repository

import (
	"testing"

	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSchemaRepository(t *testing.T) {
	db := setupTestDB(t)

	repo := NewSchemaRepository(db)
	assert.NotNil(t, repo)
}

func TestSchemaRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project first
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Create a schema
	schema := &domain.Schema{
		ProjectID: project.ID,
		Content:   `{"entities": [{"name": "User"}]}`,
		Version:   1,
	}

	err = schemaRepo.Create(schema)
	require.NoError(t, err)
	assert.NotEmpty(t, schema.ID)
	assert.NotZero(t, schema.CreatedAt)
	assert.NotZero(t, schema.UpdatedAt)
}

func TestSchemaRepository_Create_MultipleVersions(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project first
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Create multiple schema versions
	schemas := []*domain.Schema{
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}]}`,
			Version:   1,
		},
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}, {"name": "Post"}]}`,
			Version:   2,
		},
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}, {"name": "Post"}, {"name": "Comment"}]}`,
			Version:   3,
		},
	}

	for _, s := range schemas {
		err := schemaRepo.Create(s)
		require.NoError(t, err)
		assert.NotEmpty(t, s.ID)
	}
}

func TestSchemaRepository_FindByProjectID(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Create multiple schemas
	schemas := []*domain.Schema{
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}]}`,
			Version:   1,
		},
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}, {"name": "Post"}]}`,
			Version:   2,
		},
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}, {"name": "Post"}, {"name": "Comment"}]}`,
			Version:   3,
		},
	}

	for _, s := range schemas {
		err := schemaRepo.Create(s)
		require.NoError(t, err)
	}

	// Find all schemas for the project
	result, err := schemaRepo.FindByProjectID(project.ID)
	require.NoError(t, err)
	assert.Len(t, result, 3)

	// Verify ordering by version DESC (highest version first)
	assert.Equal(t, 3, result[0].Version)
	assert.Equal(t, 2, result[1].Version)
	assert.Equal(t, 1, result[2].Version)
}

func TestSchemaRepository_FindByProjectID_Empty(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project without schemas
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Find schemas
	result, err := schemaRepo.FindByProjectID(project.ID)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestSchemaRepository_FindByProjectID_MultipleProjects(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create two projects
	project1 := &domain.Project{
		Name:         "Project 1",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project1)
	require.NoError(t, err)

	project2 := &domain.Project{
		Name:         "Project 2",
		Language:     "python",
		DatabaseType: "mongodb",
		APIStyle:     "graphql",
	}
	err = projectRepo.Create(project2)
	require.NoError(t, err)

	// Create schemas for both projects
	schema1 := &domain.Schema{
		ProjectID: project1.ID,
		Content:   `{"entities": [{"name": "User"}]}`,
		Version:   1,
	}
	err = schemaRepo.Create(schema1)
	require.NoError(t, err)

	schema2 := &domain.Schema{
		ProjectID: project2.ID,
		Content:   `{"entities": [{"name": "Product"}]}`,
		Version:   1,
	}
	err = schemaRepo.Create(schema2)
	require.NoError(t, err)

	// Verify each project only sees its own schemas
	result1, err := schemaRepo.FindByProjectID(project1.ID)
	require.NoError(t, err)
	assert.Len(t, result1, 1)
	assert.Contains(t, result1[0].Content, "User")

	result2, err := schemaRepo.FindByProjectID(project2.ID)
	require.NoError(t, err)
	assert.Len(t, result2, 1)
	assert.Contains(t, result2[0].Content, "Product")
}

func TestSchemaRepository_FindLatestByProjectID(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Create multiple schemas
	schemas := []*domain.Schema{
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}]}`,
			Version:   1,
		},
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}, {"name": "Post"}]}`,
			Version:   2,
		},
		{
			ProjectID: project.ID,
			Content:   `{"entities": [{"name": "User"}, {"name": "Post"}, {"name": "Comment"}]}`,
			Version:   3,
		},
	}

	for _, s := range schemas {
		err := schemaRepo.Create(s)
		require.NoError(t, err)
	}

	// Find latest schema
	latest, err := schemaRepo.FindLatestByProjectID(project.ID)
	require.NoError(t, err)
	assert.NotNil(t, latest)
	assert.Equal(t, 3, latest.Version)
	assert.Contains(t, latest.Content, "Comment")
}

func TestSchemaRepository_FindLatestByProjectID_SingleVersion(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Create single schema
	schema := &domain.Schema{
		ProjectID: project.ID,
		Content:   `{"entities": [{"name": "User"}]}`,
		Version:   1,
	}
	err = schemaRepo.Create(schema)
	require.NoError(t, err)

	// Find latest schema
	latest, err := schemaRepo.FindLatestByProjectID(project.ID)
	require.NoError(t, err)
	assert.NotNil(t, latest)
	assert.Equal(t, 1, latest.Version)
	assert.Equal(t, schema.ID, latest.ID)
}

func TestSchemaRepository_FindLatestByProjectID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	projectRepo := NewProjectRepository(db)
	schemaRepo := NewSchemaRepository(db)

	// Create a project without schemas
	project := &domain.Project{
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	err := projectRepo.Create(project)
	require.NoError(t, err)

	// Try to find latest schema
	latest, err := schemaRepo.FindLatestByProjectID(project.ID)
	assert.Error(t, err)
	assert.Equal(t, "schema not found", err.Error())
	assert.Nil(t, latest)
}

func TestSchemaRepository_FindLatestByProjectID_NonExistentProject(t *testing.T) {
	db := setupTestDB(t)
	schemaRepo := NewSchemaRepository(db)

	// Try to find latest schema for non-existent project
	latest, err := schemaRepo.FindLatestByProjectID("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	assert.Equal(t, "schema not found", err.Error())
	assert.Nil(t, latest)
}
