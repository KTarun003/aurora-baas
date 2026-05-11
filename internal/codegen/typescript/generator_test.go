package typescript

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ktarun.reddy/baas/internal/codegen"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

func TestTypeScriptGenerator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *codegen.GeneratorContext
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid context",
			ctx: &codegen.GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "Test Project",
				Language:    "typescript",
				APIStyle:    "rest",
				Database:    "postgres",
				Schema: &validator.SchemaDefinition{
					Database: "postgres",
					Tables:   []validator.Table{},
				},
				OutputDir: "/tmp/test",
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid language",
			ctx: &codegen.GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "Test Project",
				Language:    "python",
				APIStyle:    "rest",
				Database:    "postgres",
				Schema: &validator.SchemaDefinition{
					Database: "postgres",
					Tables:   []validator.Table{},
				},
				OutputDir: "/tmp/test",
			},
			wantErr: true,
			errMsg:  "TypeScript generator only supports typescript language",
		},
		{
			name: "invalid API style",
			ctx: &codegen.GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "Test Project",
				Language:    "typescript",
				APIStyle:    "graphql",
				Database:    "postgres",
				Schema: &validator.SchemaDefinition{
					Database: "postgres",
					Tables:   []validator.Table{},
				},
				OutputDir: "/tmp/test",
			},
			wantErr: true,
			errMsg:  "TypeScript generator only supports rest API style",
		},
		{
			name: "invalid database",
			ctx: &codegen.GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "Test Project",
				Language:    "typescript",
				APIStyle:    "rest",
				Database:    "mongodb",
				Schema: &validator.SchemaDefinition{
					Database: "mongodb",
					Tables:   []validator.Table{},
				},
				OutputDir: "/tmp/test",
			},
			wantErr: true,
			errMsg:  "TypeScript generator only supports postgres database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewTypeScriptGenerator()
			err := gen.Validate(tt.ctx)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestTypeScriptGenerator_Generate(t *testing.T) {
	// Create temporary directory for output
	tmpDir, err := os.MkdirTemp("", "typescript-gen-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test schema
	schema := &validator.SchemaDefinition{
		Database: "postgres",
		Tables: []validator.Table{
			{
				Name: "users",
				Columns: []validator.Column{
					{
						Name:       "id",
						Type:       "uuid",
						PrimaryKey: true,
						Nullable:   false,
					},
					{
						Name:     "email",
						Type:     "string",
						Unique:   true,
						Nullable: false,
					},
					{
						Name:     "name",
						Type:     "string",
						Nullable: false,
					},
					{
						Name:     "created_at",
						Type:     "timestamp",
						Nullable: false,
						Default:  "now()",
					},
				},
			},
			{
				Name: "posts",
				Columns: []validator.Column{
					{
						Name:       "id",
						Type:       "uuid",
						PrimaryKey: true,
						Nullable:   false,
					},
					{
						Name:     "title",
						Type:     "string",
						Nullable: false,
					},
					{
						Name:     "content",
						Type:     "text",
						Nullable: false,
					},
					{
						Name:        "user_id",
						Type:        "uuid",
						Nullable:    false,
						ForeignKey:  "users.id",
					},
					{
						Name:     "created_at",
						Type:     "timestamp",
						Nullable: false,
						Default:  "now()",
					},
				},
			},
		},
	}

	ctx := &codegen.GeneratorContext{
		ProjectID:   "test-project",
		ProjectName: "Test Project",
		Language:    "typescript",
		APIStyle:    "rest",
		Database:    "postgres",
		Schema:      schema,
		OutputDir:   tmpDir,
		Timestamp:   time.Now(),
		Options:     map[string]interface{}{},
	}

	gen := NewTypeScriptGenerator()
	project, err := gen.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Verify project metadata
	if project.ProjectID != "test-project" {
		t.Errorf("ProjectID = %v, want %v", project.ProjectID, "test-project")
	}
	if project.Language != "typescript" {
		t.Errorf("Language = %v, want %v", project.Language, "typescript")
	}

	// Verify base files are generated
	expectedBaseFiles := []string{
		"package.json",
		"tsconfig.json",
		".gitignore",
		".env.example",
		"README.md",
		"src/server.ts",
		"src/database.ts",
	}

	for _, expectedPath := range expectedBaseFiles {
		found := false
		for _, file := range project.Files {
			if file.Path == expectedPath {
				found = true
				if file.Content == "" {
					t.Errorf("File %s has empty content", expectedPath)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in generated files", expectedPath)
		}
	}

	// Verify Prisma schema is generated
	prismaSchemaFound := false
	for _, file := range project.Files {
		if file.Path == "prisma/schema.prisma" {
			prismaSchemaFound = true
			if file.Content == "" {
				t.Errorf("Prisma schema has empty content")
			}
			// Check for expected models
			if !contains(file.Content, "model User") {
				t.Errorf("Prisma schema missing User model")
			}
			if !contains(file.Content, "model Post") {
				t.Errorf("Prisma schema missing Post model")
			}
			break
		}
	}
	if !prismaSchemaFound {
		t.Errorf("Prisma schema file not found")
	}

	// Verify CRUD files are generated for each table
	expectedCRUDFiles := []struct {
		path     string
		contains []string
	}{
		{
			path:     "src/routes/users.generated.ts",
			contains: []string{"User", "router"},
		},
		{
			path:     "src/controllers/users.generated.ts",
			contains: []string{"User", "create", "update", "delete"},
		},
		{
			path:     "src/models/users.generated.ts",
			contains: []string{"User", "Schema"},
		},
		{
			path:     "src/routes/posts.generated.ts",
			contains: []string{"Post", "router"},
		},
		{
			path:     "src/controllers/posts.generated.ts",
			contains: []string{"Post", "create", "update", "delete"},
		},
		{
			path:     "src/models/posts.generated.ts",
			contains: []string{"Post", "Schema"},
		},
	}

	for _, expected := range expectedCRUDFiles {
		found := false
		for _, file := range project.Files {
			if file.Path == expected.path {
				found = true
				if file.Content == "" {
					t.Errorf("File %s has empty content", expected.path)
				}
				for _, substr := range expected.contains {
					if !contains(file.Content, substr) {
						t.Errorf("File %s missing expected content: %s", expected.path, substr)
					}
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected CRUD file %s not found", expected.path)
		}
	}

	// Verify extension files are generated
	expectedExtensionFiles := []string{
		"src/hooks/users.ts",
		"src/hooks/posts.ts",
		"src/routes/custom.ts",
	}

	for _, expectedPath := range expectedExtensionFiles {
		found := false
		for _, file := range project.Files {
			if file.Path == expectedPath {
				found = true
				if file.IsStatic != true {
					t.Errorf("Extension file %s should be marked as static", expectedPath)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected extension file %s not found", expectedPath)
		}
	}

	// Verify files are written to disk
	for _, file := range project.Files {
		fullPath := filepath.Join(tmpDir, file.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("File %s was not written to disk", file.Path)
		}
	}
}

func TestTypeScriptGenerator_EmptySchema(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "typescript-gen-empty-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create empty schema
	schema := &validator.SchemaDefinition{
		Database: "postgres",
		Tables:   []validator.Table{},
	}

	ctx := &codegen.GeneratorContext{
		ProjectID:   "test-project",
		ProjectName: "Test Project",
		Language:    "typescript",
		APIStyle:    "rest",
		Database:    "postgres",
		Schema:      schema,
		OutputDir:   tmpDir,
		Timestamp:   time.Now(),
	}

	gen := NewTypeScriptGenerator()
	project, err := gen.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Should still generate base files
	if len(project.Files) == 0 {
		t.Error("Expected base files to be generated even with empty schema")
	}

	// Verify base files exist
	baseFileFound := false
	for _, file := range project.Files {
		if file.Path == "package.json" {
			baseFileFound = true
			break
		}
	}
	if !baseFileFound {
		t.Error("Base files not generated for empty schema")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
