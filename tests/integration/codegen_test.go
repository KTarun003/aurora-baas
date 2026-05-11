//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

// setupTestServerWithCodegen creates and configures a test HTTP server with codegen support
func setupTestServerWithCodegen(t *testing.T) *httptest.Server {
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

	// Initialize code generation
	codegenRegistry := codegen.NewRegistry()
	tsGenerator := typescript.NewTypeScriptGenerator()
	codegenRegistry.Register(tsGenerator)

	codegenService := service.NewCodegenService(codegenRegistry, projectRepository, schemaRepository)

	// Create router with codegen support
	router := api.NewRouter(projectService, schemaService, codegenService)

	// Return httptest.Server with the router
	return httptest.NewServer(router)
}

// TestCodegen_FullWorkflow tests the complete code generation workflow
func TestCodegen_FullWorkflow(t *testing.T) {
	server := setupTestServerWithCodegen(t)
	defer server.Close()

	var projectID string

	// Step 1: Create project
	t.Run("Create Project", func(t *testing.T) {
		projectPayload := map[string]string{
			"name":          "Codegen Test Project",
			"description":   "Testing code generation",
			"language":      "typescript",
			"database_type": "postgres",
			"api_style":     "rest",
		}

		payload, err := json.Marshal(projectPayload)
		require.NoError(t, err)

		resp, err := http.Post(
			server.URL+"/api/v1/projects",
			"application/json",
			bytes.NewBuffer(payload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		projectID = result["id"].(string)
		assert.NotEmpty(t, projectID)
		t.Logf("Created project with ID: %s", projectID)
	})

	// Step 2: Apply schema
	t.Run("Apply Schema", func(t *testing.T) {
		yamlContent := `
version: 1.0
database: postgres
tables:
  - name: users
    columns:
      - name: id
        type: uuid
        primary_key: true
      - name: email
        type: string
        unique: true
        nullable: false
      - name: name
        type: string
        nullable: false
      - name: created_at
        type: timestamp
        default: "now()"
      - name: updated_at
        type: timestamp
        default: "now()"
  - name: posts
    columns:
      - name: id
        type: uuid
        primary_key: true
      - name: user_id
        type: uuid
        nullable: false
        foreign_key:
          table: users
          column: id
      - name: title
        type: string
        nullable: false
      - name: content
        type: text
        nullable: false
      - name: published
        type: boolean
        default: "false"
      - name: created_at
        type: timestamp
        default: "now()"
      - name: updated_at
        type: timestamp
        default: "now()"
`

		resp, err := http.Post(
			server.URL+"/api/v1/projects/"+projectID+"/schemas",
			"application/yaml",
			bytes.NewBufferString(yamlContent),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, projectID, result["project_id"])
		assert.Equal(t, float64(1), result["version"])
		t.Logf("Applied schema version: %v", result["version"])
	})

	// Step 3: Generate code
	t.Run("Generate Code", func(t *testing.T) {
		// Create temporary output directory
		outputDir, err := os.MkdirTemp("", "baas-codegen-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(outputDir)

		t.Logf("Output directory: %s", outputDir)

		generatePayload := map[string]string{
			"output_dir": outputDir,
		}

		payload, err := json.Marshal(generatePayload)
		require.NoError(t, err)

		resp, err := http.Post(
			server.URL+"/api/v1/projects/"+projectID+"/generate",
			"application/json",
			bytes.NewBuffer(payload),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, projectID, result["project_id"])
		assert.Equal(t, "typescript", result["language"])
		assert.Greater(t, int(result["files_count"].(float64)), 0)

		t.Logf("Generated %v files", result["files_count"])

		// Verify some key files exist
		expectedFiles := []string{
			"package.json",
			"tsconfig.json",
			".gitignore",
			".env.example",
			"README.md",
			"src/server.ts",
			"src/database.ts",
			"prisma/schema.prisma",
			"src/routes/users.generated.ts",
			"src/routes/posts.generated.ts",
			"src/controllers/users.generated.ts",
			"src/controllers/posts.generated.ts",
			"src/models/users.generated.ts",
			"src/models/posts.generated.ts",
			"src/hooks/users.ts",
			"src/hooks/posts.ts",
			"src/routes/custom.ts",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(outputDir, file)
			_, err := os.Stat(filePath)
			assert.NoError(t, err, "expected file to exist: %s", file)

			if err == nil {
				// Read and check file has content
				content, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr, "failed to read file: %s", file)
				assert.NotEmpty(t, content, "expected file to have content: %s", file)
				t.Logf("✓ Verified file: %s (%d bytes)", file, len(content))
			}
		}

		// Verify Prisma schema contains our tables
		prismaPath := filepath.Join(outputDir, "prisma/schema.prisma")
		prismaContent, err := os.ReadFile(prismaPath)
		require.NoError(t, err)

		prismaStr := string(prismaContent)
		assert.Contains(t, prismaStr, "model User", "Prisma schema should contain User model")
		assert.Contains(t, prismaStr, "model Post", "Prisma schema should contain Post model")
		assert.Contains(t, prismaStr, "datasource db", "Prisma schema should contain datasource")
		assert.Contains(t, prismaStr, "generator client", "Prisma schema should contain generator")

		// Verify package.json has required dependencies
		packagePath := filepath.Join(outputDir, "package.json")
		packageContent, err := os.ReadFile(packagePath)
		require.NoError(t, err)

		var packageJSON map[string]interface{}
		err = json.Unmarshal(packageContent, &packageJSON)
		require.NoError(t, err)

		dependencies, ok := packageJSON["dependencies"].(map[string]interface{})
		assert.True(t, ok, "package.json should have dependencies")
		assert.Contains(t, dependencies, "express")
		assert.Contains(t, dependencies, "@prisma/client")
		assert.Contains(t, dependencies, "zod")

		t.Log("✓ All verifications passed!")
	})

	// Step 4: Verify project still accessible
	t.Run("Verify Project After Generation", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/projects/" + projectID)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, projectID, result["id"])
		assert.Equal(t, "Codegen Test Project", result["name"])
	})
}

// TestCodegen_WithoutSchema tests error handling when trying to generate without a schema
func TestCodegen_WithoutSchema(t *testing.T) {
	server := setupTestServerWithCodegen(t)
	defer server.Close()

	// Create project without schema
	projectPayload := map[string]string{
		"name":          "No Schema Project",
		"description":   "Testing error handling",
		"language":      "typescript",
		"database_type": "postgres",
		"api_style":     "rest",
	}

	payload, err := json.Marshal(projectPayload)
	require.NoError(t, err)

	resp, err := http.Post(
		server.URL+"/api/v1/projects",
		"application/json",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	projectID := result["id"].(string)

	// Try to generate without schema
	outputDir, err := os.MkdirTemp("", "baas-codegen-noschema-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	generatePayload := map[string]string{
		"output_dir": outputDir,
	}

	generateReq, err := json.Marshal(generatePayload)
	require.NoError(t, err)

	genResp, err := http.Post(
		server.URL+"/api/v1/projects/"+projectID+"/generate",
		"application/json",
		bytes.NewBuffer(generateReq),
	)
	require.NoError(t, err)
	defer genResp.Body.Close()

	// Should fail with 500 because there's no schema
	assert.Equal(t, http.StatusInternalServerError, genResp.StatusCode)

	var errorResult map[string]interface{}
	err = json.NewDecoder(genResp.Body).Decode(&errorResult)
	require.NoError(t, err)

	assert.Contains(t, errorResult["error"], "failed to get schema")
}
