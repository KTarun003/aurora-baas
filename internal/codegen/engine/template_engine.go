package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateEngine handles template rendering with custom functions
type TemplateEngine struct {
	funcMap template.FuncMap
}

// NewTemplateEngine creates a new template engine with helper functions
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		funcMap: GetHelperFunctions(),
	}
}

// RenderString renders a template string with the given data
func (e *TemplateEngine) RenderString(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Funcs(e.funcMap).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderFile renders a template file with the given data
func (e *TemplateEngine) RenderFile(templatePath string, data interface{}) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return e.RenderString(string(content), data)
}

// RenderFiles renders multiple template files from a directory
func (e *TemplateEngine) RenderFiles(templateDir string, data interface{}) (map[string]string, error) {
	results := make(map[string]string)

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() || filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Render the template
		content, err := e.RenderFile(path, data)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", path, err)
		}

		// Store with relative path (without .tmpl extension)
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Safety check: ensure path is long enough to contain .tmpl extension
		if len(relPath) < 6 { // 5 for .tmpl + 1 for filename
			return fmt.Errorf("invalid template path %q: too short to contain .tmpl extension", relPath)
		}

		outputPath := relPath[:len(relPath)-5] // Remove .tmpl extension

		results[outputPath] = content
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
