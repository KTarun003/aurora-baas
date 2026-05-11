package engine

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateEngine_RenderString(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		data     interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple variable",
			template: "Hello {{.Name}}",
			data:     map[string]string{"Name": "World"},
			want:     "Hello World",
			wantErr:  false,
		},
		{
			name:     "with helper function",
			template: "{{pascalCase .name}}",
			data:     map[string]string{"name": "user_profile"},
			want:     "UserProfile",
			wantErr:  false,
		},
		{
			name:     "invalid template",
			template: "{{.Missing",
			data:     map[string]string{},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestTemplateEngine_RenderFile(t *testing.T) {
	engine := NewTemplateEngine()

	// Create a temporary template file
	tmpFile := t.TempDir() + "/template.tmpl"
	tmpl := "package {{.Package}}\n\nconst Version = \"{{.Version}}\""
	err := os.WriteFile(tmpFile, []byte(tmpl), 0644)
	require.NoError(t, err)

	data := map[string]string{
		"Package": "main",
		"Version": "1.0.0",
	}

	result, err := engine.RenderFile(tmpFile, data)
	require.NoError(t, err)
	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "const Version = \"1.0.0\"")
}

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		helper   string
		input    string
		expected string
	}{
		{"pascalCase simple", "pascalCase", "user_profile", "UserProfile"},
		{"pascalCase mixed", "pascalCase", "user-profile-data", "UserProfileData"},
		{"camelCase simple", "camelCase", "user_profile", "userProfile"},
		{"snakeCase simple", "snakeCase", "UserProfile", "user_profile"},
		{"kebabCase simple", "kebabCase", "UserProfile", "user-profile"},
		{"plural regular", "plural", "user", "users"},
		{"plural y", "plural", "category", "categories"},
		{"plural s", "plural", "class", "classes"},
		{"singular regular", "singular", "users", "user"},
		{"singular ies", "singular", "categories", "category"},
	}

	engine := NewTemplateEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := fmt.Sprintf("{{%s .input}}", tt.helper)
			data := map[string]string{"input": tt.input}

			result, err := engine.RenderString(tmpl, data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_Upper(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "hello", "HELLO"},
		{"mixed case", "HeLLo", "HELLO"},
		{"with numbers", "hello123", "HELLO123"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString("{{upper .input}}", map[string]string{"input": tt.input})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_Lower(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "HELLO", "hello"},
		{"mixed case", "HeLLo", "hello"},
		{"with numbers", "HELLO123", "hello123"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString("{{lower .input}}", map[string]string{"input": tt.input})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_Title(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "hello", "Hello"},
		{"multiple words", "hello world", "Hello World"},
		{"already titled", "Hello World", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString("{{title .input}}", map[string]string{"input": tt.input})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_Join(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"simple join", `{{join (split "hello-world" "-") ","}}`, "hello,world"},
		{"with dots", `{{join (split "a.b.c" ".") "/"}}`, "a/b/c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, map[string]string{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_Split(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"split by hyphen", `{{index (split "hello-world" "-") 0}}`, "hello"},
		{"split by comma", `{{index (split "a,b,c" ",") 2}}`, "c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, map[string]string{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_Contains(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"contains true", `{{contains "hello world" "world"}}`, "true"},
		{"contains false", `{{contains "hello" "xyz"}}`, "false"},
		{"contains substring", `{{contains "user_profile" "_"}}`, "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, map[string]string{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_HasPrefix(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"has prefix true", `{{hasPrefix "UserProfile" "User"}}`, "true"},
		{"has prefix false", `{{hasPrefix "Profile" "User"}}`, "false"},
		{"has prefix underscore", `{{hasPrefix "user_name" "user"}}`, "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, map[string]string{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHelperFunction_HasSuffix(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"has suffix true", `{{hasSuffix "UserProfile" "Profile"}}`, "true"},
		{"has suffix false", `{{hasSuffix "User" "Profile"}}`, "false"},
		{"has suffix with underscore", `{{hasSuffix "user_name" "name"}}`, "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.RenderString(tt.template, map[string]string{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateEngine_RenderFiles(t *testing.T) {
	engine := NewTemplateEngine()

	tests := []struct {
		name              string
		setupFunc         func(t *testing.T, tmpDir string) string // Returns template dir
		data              interface{}
		expectedFiles     map[string]string
		shouldErr         bool
		errMsg            string
	}{
		{
			name: "render single template file",
			setupFunc: func(t *testing.T, tmpDir string) string {
				templateDir := tmpDir + "/templates"
				err := os.Mkdir(templateDir, 0755)
				require.NoError(t, err)

				// Create a single template file
				content := "Hello {{.Name}}"
				err = os.WriteFile(templateDir+"/greeting.tmpl", []byte(content), 0644)
				require.NoError(t, err)

				return templateDir
			},
			data: map[string]string{"Name": "World"},
			expectedFiles: map[string]string{
				"greeting": "Hello World",
			},
			shouldErr: false,
		},
		{
			name: "render multiple template files",
			setupFunc: func(t *testing.T, tmpDir string) string {
				templateDir := tmpDir + "/templates"
				err := os.Mkdir(templateDir, 0755)
				require.NoError(t, err)

				// Create multiple template files
				files := map[string]string{
					"greeting.tmpl": "Hello {{.Name}}",
					"farewell.tmpl": "Goodbye {{.Name}}",
					"title.tmpl":    "{{title .title}}",
				}

				for filename, content := range files {
					err := os.WriteFile(templateDir+"/"+filename, []byte(content), 0644)
					require.NoError(t, err)
				}

				return templateDir
			},
			data: map[string]interface{}{
				"Name":  "Alice",
				"title": "hello world",
			},
			expectedFiles: map[string]string{
				"greeting": "Hello Alice",
				"farewell": "Goodbye Alice",
				"title":    "Hello World",
			},
			shouldErr: false,
		},
		{
			name: "render templates with subdirectories",
			setupFunc: func(t *testing.T, tmpDir string) string {
				templateDir := tmpDir + "/templates"
				err := os.Mkdir(templateDir, 0755)
				require.NoError(t, err)

				// Create subdirectory
				subDir := templateDir + "/models"
				err = os.Mkdir(subDir, 0755)
				require.NoError(t, err)

				// Create files in subdirectory
				err = os.WriteFile(subDir+"/user.tmpl", []byte("type User struct { Name: \"{{.Name}}\" }"), 0644)
				require.NoError(t, err)

				// Create file in root
				err = os.WriteFile(templateDir+"/header.tmpl", []byte("// Generated code"), 0644)
				require.NoError(t, err)

				return templateDir
			},
			data: map[string]string{"Name": "John"},
			expectedFiles: map[string]string{
				"header":         "// Generated code",
				"models/user":    "type User struct { Name: \"John\" }",
			},
			shouldErr: false,
		},
		{
			name: "skip non-template files",
			setupFunc: func(t *testing.T, tmpDir string) string {
				templateDir := tmpDir + "/templates"
				err := os.Mkdir(templateDir, 0755)
				require.NoError(t, err)

				// Create template files
				err = os.WriteFile(templateDir+"/template.tmpl", []byte("template"), 0644)
				require.NoError(t, err)

				// Create non-template files (should be ignored)
				err = os.WriteFile(templateDir+"/readme.txt", []byte("not a template"), 0644)
				require.NoError(t, err)

				err = os.WriteFile(templateDir+"/data.json", []byte(`{"key":"value"}`), 0644)
				require.NoError(t, err)

				return templateDir
			},
			data: map[string]string{},
			expectedFiles: map[string]string{
				"template": "template",
			},
			shouldErr: false,
		},
		{
			name: "empty template directory",
			setupFunc: func(t *testing.T, tmpDir string) string {
				templateDir := tmpDir + "/empty"
				err := os.Mkdir(templateDir, 0755)
				require.NoError(t, err)
				return templateDir
			},
			data:          map[string]string{},
			expectedFiles: map[string]string{},
			shouldErr:     false,
		},
		{
			name: "non-existent directory",
			setupFunc: func(t *testing.T, tmpDir string) string {
				return tmpDir + "/does-not-exist"
			},
			data:      map[string]string{},
			shouldErr: true,
		},
		{
			name: "render with helper functions in templates",
			setupFunc: func(t *testing.T, tmpDir string) string {
				templateDir := tmpDir + "/templates"
				err := os.Mkdir(templateDir, 0755)
				require.NoError(t, err)

				// Create template files using helper functions
				err = os.WriteFile(
					templateDir+"/pascal.tmpl",
					[]byte("{{pascalCase .name}}"),
					0644,
				)
				require.NoError(t, err)

				err = os.WriteFile(
					templateDir+"/camel.tmpl",
					[]byte("{{camelCase .name}}"),
					0644,
				)
				require.NoError(t, err)

				err = os.WriteFile(
					templateDir+"/snake.tmpl",
					[]byte("{{snakeCase .name}}"),
					0644,
				)
				require.NoError(t, err)

				return templateDir
			},
			data: map[string]string{"name": "user_profile"},
			expectedFiles: map[string]string{
				"pascal": "UserProfile",
				"camel":  "userProfile",
				"snake":  "user_profile",
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			templateDir := tt.setupFunc(t, tmpDir)

			results, err := engine.RenderFiles(templateDir, tt.data)

			if tt.shouldErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.expectedFiles), len(results))

				for expectedKey, expectedValue := range tt.expectedFiles {
					assert.Contains(t, results, expectedKey)
					assert.Equal(t, expectedValue, results[expectedKey])
				}
			}
		})
	}
}

func TestTemplateEngine_RenderFiles_InvalidTemplatePath(t *testing.T) {
	engine := NewTemplateEngine()
	tmpDir := t.TempDir()
	templateDir := tmpDir + "/templates"
	err := os.Mkdir(templateDir, 0755)
	require.NoError(t, err)

	// Create a file with a name shorter than the .tmpl extension (edge case)
	// This shouldn't happen in practice since filepath.Walk won't match it as .tmpl
	// But we'll test the boundary condition
	shortPath := templateDir + "/a.tmpl" // Minimum valid name

	err = os.WriteFile(shortPath, []byte("content"), 0644)
	require.NoError(t, err)

	// This should work fine since "a" (path with .tmpl removed) is valid
	results, err := engine.RenderFiles(templateDir, map[string]string{})
	require.NoError(t, err)
	assert.Contains(t, results, "a")
	assert.Equal(t, "content", results["a"])
}
