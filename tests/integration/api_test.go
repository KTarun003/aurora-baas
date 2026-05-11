//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ktarun.reddy/baas/internal/api"
	"github.com/ktarun.reddy/baas/internal/codegen"
	"github.com/ktarun.reddy/baas/internal/codegen/typescript"
	"github.com/ktarun.reddy/baas/internal/config"
	"github.com/ktarun.reddy/baas/internal/database"
	"github.com/ktarun.reddy/baas/internal/repository"
	"github.com/ktarun.reddy/baas/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestServer creates and configures a test HTTP server
func setupTestServer(t *testing.T) *httptest.Server {
	// Create database config
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "aurora",
		Password: "aurora_dev",
		Name:     "aurora_dev",
		SSLMode:  "disable",
	}

	// Connect to database
	db, err := database.Connect(cfg)
	require.NoError(t, err, "failed to connect to database")

	// Run migrations
	err = database.RunMigrations(db)
	require.NoError(t, err, "failed to run migrations")

	// Clean up existing test data
	sqlDB, err := db.DB()
	require.NoError(t, err, "failed to get database instance")

	// Delete schemas first (foreign key constraint)
	sqlDB.Exec("DELETE FROM schemas")
	// Delete projects
	sqlDB.Exec("DELETE FROM projects")

	// Create repositories
	projectRepository := repository.NewProjectRepository(db)
	schemaRepository := repository.NewSchemaRepository(db)

	// Create services
	projectService := service.NewProjectService(projectRepository)
	schemaService := service.NewSchemaService(schemaRepository, projectRepository)

	// Initialize code generation (needed for router, even if not used in basic tests)
	codegenRegistry := codegen.NewRegistry()
	tsGenerator := typescript.NewTypeScriptGenerator()
	codegenRegistry.Register(tsGenerator)
	codegenService := service.NewCodegenService(codegenRegistry, projectRepository, schemaRepository)

	// Create router
	router := api.NewRouter(projectService, schemaService, codegenService)

	// Return httptest.Server with the router
	return httptest.NewServer(router)
}

// TestAPI_FullWorkflow tests the complete API workflow
func TestAPI_FullWorkflow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Step 1: Health check - GET /health, expect 200
	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		require.NoError(t, err, "failed to make request to /health")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 for health check")
	})

	var projectID string

	// Step 2: Create project - POST /api/v1/projects with JSON body, expect 201
	t.Run("Create Project", func(t *testing.T) {
		projectPayload := map[string]string{
			"name":          "Test Project",
			"description":   "A test project for integration testing",
			"language":      "typescript",
			"database_type": "postgres",
			"api_style":     "rest",
		}

		payload, err := json.Marshal(projectPayload)
		require.NoError(t, err, "failed to marshal project payload")

		resp, err := http.Post(
			server.URL+"/api/v1/projects",
			"application/json",
			bytes.NewBuffer(payload),
		)
		require.NoError(t, err, "failed to make POST request to /api/v1/projects")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected status 201 for project creation")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response")

		projectID = result["id"].(string)
		assert.NotEmpty(t, projectID, "expected project ID in response")
		assert.Equal(t, "Test Project", result["name"])
	})

	// Step 3: Get project - GET /api/v1/projects/:id, expect 200
	t.Run("Get Project", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/projects/" + projectID)
		require.NoError(t, err, "failed to make GET request to project endpoint")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 for get project")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response")

		assert.Equal(t, projectID, result["id"])
		assert.Equal(t, "Test Project", result["name"])
	})

	// Step 4: List projects - GET /api/v1/projects, expect 200
	t.Run("List Projects", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/projects")
		require.NoError(t, err, "failed to make GET request to list projects")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 for list projects")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response")

		projects, ok := result["projects"].([]interface{})
		assert.True(t, ok, "expected projects array in response")
		assert.GreaterOrEqual(t, len(projects), 1, "expected at least one project")
	})

	// Step 5: Apply schema - POST /api/v1/projects/:id/schemas with YAML body, expect 201
	t.Run("Apply Schema", func(t *testing.T) {
		yamlContent := `
version: 1.0
name: TestSchema
entities:
  - name: User
    fields:
      - name: id
        type: string
        required: true
      - name: email
        type: string
        required: true
  - name: Post
    fields:
      - name: id
        type: string
        required: true
      - name: title
        type: string
        required: true
`

		resp, err := http.Post(
			server.URL+"/api/v1/projects/"+projectID+"/schemas",
			"application/yaml",
			bytes.NewBufferString(yamlContent),
		)
		require.NoError(t, err, "failed to make POST request to apply schema")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected status 201 for schema creation")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response")

		assert.Equal(t, projectID, result["project_id"])
		assert.Equal(t, float64(1), result["version"], "expected version 1")
	})

	// Step 6: Get latest schema - GET /api/v1/projects/:id/schemas/latest, expect 200
	t.Run("Get Latest Schema", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/projects/" + projectID + "/schemas/latest")
		require.NoError(t, err, "failed to make GET request to get latest schema")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 for get latest schema")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response")

		assert.Equal(t, projectID, result["project_id"])
		assert.Equal(t, float64(1), result["version"])
		assert.NotEmpty(t, result["content"], "expected schema content")
	})

	// Step 7: Delete project - DELETE /api/v1/projects/:id, expect 200
	t.Run("Delete Project", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/projects/"+projectID, nil)
		require.NoError(t, err, "failed to create DELETE request")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err, "failed to make DELETE request")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 for delete project")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response")

		assert.Equal(t, "project deleted", result["message"])
	})
}
