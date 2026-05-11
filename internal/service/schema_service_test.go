package service

import (
	"testing"

	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSchemaService_ApplySchema tests successful schema application
func TestSchemaService_ApplySchema(t *testing.T) {
	projectService, schemaService := setupTestService(t)

	// Create a project first
	project := &domain.Project{
		Name:         "ApplySchema Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)

	validYAML := `
version: "1.0"
database: postgres
tables:
  - name: users
    columns:
      - name: id
        type: uuid
        primary_key: true
      - name: name
        type: string
`

	schema, err := schemaService.ApplySchema(project.ID, validYAML)
	assert.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, 1, schema.Version)
}

// TestSchemaService_ApplySchema_InvalidYAML tests schema application with invalid YAML
func TestSchemaService_ApplySchema_InvalidYAML(t *testing.T) {
	projectService, schemaService := setupTestService(t)

	// Create a project first
	project := &domain.Project{
		Name:         "InvalidYAML Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)

	invalidYAML := "{ invalid yaml :"

	schema, err := schemaService.ApplySchema(project.ID, invalidYAML)
	assert.Error(t, err)
	assert.Nil(t, schema)
}

// TestSchemaService_ApplySchema_ProjectNotFound tests schema application with non-existent project
func TestSchemaService_ApplySchema_ProjectNotFound(t *testing.T) {
	_, schemaService := setupTestService(t)

	projectID := "00000000-0000-0000-0000-000000000099"
	validYAML := `
version: "1.0"
database: postgres
tables:
  - name: users
    columns:
      - name: id
        type: uuid
        primary_key: true
`

	schema, err := schemaService.ApplySchema(projectID, validYAML)
	assert.Error(t, err)
	assert.Nil(t, schema)
	assert.Contains(t, err.Error(), "project not found")
}

// TestSchemaService_GetLatestSchema tests getting the latest schema version
func TestSchemaService_GetLatestSchema(t *testing.T) {
	projectService, schemaService := setupTestService(t)

	// Create a project first
	project := &domain.Project{
		Name:         "Schema Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)

	// Apply 3 schema versions
	validYAML := `
version: "1.0"
database: postgres
tables:
  - name: users
    columns:
      - name: id
        type: uuid
        primary_key: true
      - name: name
        type: string
`

	for i := 0; i < 3; i++ {
		_, err := schemaService.ApplySchema(project.ID, validYAML)
		require.NoError(t, err)
	}

	// Get the latest schema
	latest, err := schemaService.GetLatestSchema(project.ID)
	require.NoError(t, err)
	assert.NotNil(t, latest)
	assert.Equal(t, 3, latest.Version)
}

// TestSchemaService_ListSchemas tests listing all schema versions for a project
func TestSchemaService_ListSchemas(t *testing.T) {
	projectService, schemaService := setupTestService(t)

	// Create a project first
	project := &domain.Project{
		Name:         "List Schemas Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}

	err := projectService.CreateProject(project)
	require.NoError(t, err)

	// Apply 3 schema versions
	validYAML := `
version: "1.0"
database: postgres
tables:
  - name: users
    columns:
      - name: id
        type: uuid
        primary_key: true
      - name: name
        type: string
`

	for i := 0; i < 3; i++ {
		_, err := schemaService.ApplySchema(project.ID, validYAML)
		require.NoError(t, err)
	}

	// List all schemas
	schemas, err := schemaService.ListSchemas(project.ID)
	require.NoError(t, err)
	assert.Len(t, schemas, 3)

	// Verify versions
	assert.Equal(t, 3, schemas[0].Version)
	assert.Equal(t, 2, schemas[1].Version)
	assert.Equal(t, 1, schemas[2].Version)
}
