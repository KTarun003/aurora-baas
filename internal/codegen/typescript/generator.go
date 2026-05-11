package typescript

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ktarun.reddy/baas/internal/codegen"
	"github.com/ktarun.reddy/baas/internal/codegen/engine"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// TypeScriptGenerator generates TypeScript REST API projects
type TypeScriptGenerator struct {
	templateEngine *engine.TemplateEngine
}

// NewTypeScriptGenerator creates a new TypeScript generator
func NewTypeScriptGenerator() *TypeScriptGenerator {
	return &TypeScriptGenerator{
		templateEngine: engine.NewTemplateEngine(),
	}
}

// Validate checks if the context is valid for TypeScript generation
func (g *TypeScriptGenerator) Validate(ctx *codegen.GeneratorContext) error {
	// First validate the basic context
	if err := ctx.Validate(); err != nil {
		return err
	}

	// TypeScript-specific validations
	if ctx.Language != "typescript" {
		return fmt.Errorf("TypeScript generator only supports typescript language")
	}

	if ctx.APIStyle != "rest" {
		return fmt.Errorf("TypeScript generator only supports rest API style")
	}

	if ctx.Database != "postgres" {
		return fmt.Errorf("TypeScript generator only supports postgres database")
	}

	return nil
}

// Generate orchestrates the complete generation process
func (g *TypeScriptGenerator) Generate(ctx *codegen.GeneratorContext) (*codegen.GeneratedProject, error) {
	// Validate context
	if err := g.Validate(ctx); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Initialize file writer with output directory as a local variable
	fileWriter := engine.NewFileWriter(ctx.OutputDir)

	// Create project structure
	project := &codegen.GeneratedProject{
		ProjectID: ctx.ProjectID,
		Language:  ctx.Language,
		Files:     []*codegen.GeneratedFile{},
		Metadata: map[string]string{
			"api_style": ctx.APIStyle,
			"database":  ctx.Database,
		},
	}

	// Generate all components
	if err := g.generateBaseFiles(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to generate base files: %w", err)
	}

	if err := g.generatePrismaSchema(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to generate Prisma schema: %w", err)
	}

	if err := g.generateCRUDFiles(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to generate CRUD files: %w", err)
	}

	if err := g.generateExtensionFiles(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to generate extension files: %w", err)
	}

	// Write all files to disk
	if err := fileWriter.WriteFiles(project.Files); err != nil {
		return nil, fmt.Errorf("failed to write files: %w", err)
	}

	return project, nil
}

// generateBaseFiles generates package.json, tsconfig.json, etc.
func (g *TypeScriptGenerator) generateBaseFiles(ctx *codegen.GeneratorContext, project *codegen.GeneratedProject) error {
	// Get template directory
	templateDir := g.getTemplateDir("base")

	// Template data - pass full context for templates that need schema
	data := map[string]interface{}{
		"ProjectName": ctx.ProjectName,
		"ProjectID":   ctx.ProjectID,
		"Database":    ctx.Database,
		"Schema":      ctx.Schema,
	}

	// Map of template file to output path
	templateMapping := map[string]string{
		"package.json.tmpl":   "package.json",
		"tsconfig.json.tmpl":  "tsconfig.json",
		".gitignore.tmpl":     ".gitignore",
		".env.example.tmpl":   ".env.example",
		"README.md.tmpl":      "README.md",
		"server.ts.tmpl":      "src/server.ts",
		"database.ts.tmpl":    "src/database.ts",
	}

	// Render each template explicitly
	for templateFile, outputPath := range templateMapping {
		templatePath := filepath.Join(templateDir, templateFile)
		content, err := g.templateEngine.RenderFile(templatePath, data)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", templateFile, err)
		}

		file := &codegen.GeneratedFile{
			Path:     outputPath,
			Content:  content,
			Mode:     0644,
			IsStatic: false,
		}
		project.Files = append(project.Files, file)
	}

	return nil
}

// generatePrismaSchema generates Prisma schema file
func (g *TypeScriptGenerator) generatePrismaSchema(ctx *codegen.GeneratorContext, project *codegen.GeneratedProject) error {
	// Use existing Prisma generator
	schemaContent, err := GeneratePrismaSchema(ctx.Schema)
	if err != nil {
		return fmt.Errorf("failed to generate Prisma schema: %w", err)
	}

	file := &codegen.GeneratedFile{
		Path:     "prisma/schema.prisma",
		Content:  schemaContent,
		Mode:     0644,
		IsStatic: false,
	}

	project.Files = append(project.Files, file)
	return nil
}

// generateCRUDFiles generates routes, controllers, and models for each table
func (g *TypeScriptGenerator) generateCRUDFiles(ctx *codegen.GeneratorContext, project *codegen.GeneratedProject) error {
	crudTemplateDir := g.getTemplateDir("crud")

	for _, table := range ctx.Schema.Tables {
		// Prepare data for this table
		tableData := g.prepareTableData(table)

		// Generate routes
		routePath := fmt.Sprintf("src/routes/%s.generated.ts", tableData["KebabName"])
		routeContent, err := g.templateEngine.RenderFile(
			filepath.Join(crudTemplateDir, "routes.generated.ts.tmpl"),
			tableData,
		)
		if err != nil {
			return fmt.Errorf("failed to render routes for %s: %w", table.Name, err)
		}

		project.Files = append(project.Files, &codegen.GeneratedFile{
			Path:     routePath,
			Content:  routeContent,
			Mode:     0644,
			IsStatic: false,
		})

		// Generate controller
		controllerPath := fmt.Sprintf("src/controllers/%s.generated.ts", tableData["KebabName"])
		controllerContent, err := g.templateEngine.RenderFile(
			filepath.Join(crudTemplateDir, "controller.generated.ts.tmpl"),
			tableData,
		)
		if err != nil {
			return fmt.Errorf("failed to render controller for %s: %w", table.Name, err)
		}

		project.Files = append(project.Files, &codegen.GeneratedFile{
			Path:     controllerPath,
			Content:  controllerContent,
			Mode:     0644,
			IsStatic: false,
		})

		// Generate model
		modelPath := fmt.Sprintf("src/models/%s.generated.ts", tableData["KebabName"])
		modelContent, err := g.templateEngine.RenderFile(
			filepath.Join(crudTemplateDir, "model.generated.ts.tmpl"),
			tableData,
		)
		if err != nil {
			return fmt.Errorf("failed to render model for %s: %w", table.Name, err)
		}

		project.Files = append(project.Files, &codegen.GeneratedFile{
			Path:     modelPath,
			Content:  modelContent,
			Mode:     0644,
			IsStatic: false,
		})
	}

	return nil
}

// generateExtensionFiles generates hooks and custom routes templates
func (g *TypeScriptGenerator) generateExtensionFiles(ctx *codegen.GeneratorContext, project *codegen.GeneratedProject) error {
	extensionTemplateDir := g.getTemplateDir("extensions")

	// Generate hooks for each table
	for _, table := range ctx.Schema.Tables {
		tableData := g.prepareTableData(table)

		hookPath := fmt.Sprintf("src/hooks/%s.ts", tableData["KebabName"])
		hookContent, err := g.templateEngine.RenderFile(
			filepath.Join(extensionTemplateDir, "hooks.ts.tmpl"),
			tableData,
		)
		if err != nil {
			return fmt.Errorf("failed to render hooks for %s: %w", table.Name, err)
		}

		project.Files = append(project.Files, &codegen.GeneratedFile{
			Path:     hookPath,
			Content:  hookContent,
			Mode:     0644,
			IsStatic: true, // Extension files should not be overwritten
		})
	}

	// Generate custom routes template
	customRoutesPath := "src/routes/custom.ts"
	customRoutesContent, err := g.templateEngine.RenderFile(
		filepath.Join(extensionTemplateDir, "custom-routes.ts.tmpl"),
		map[string]interface{}{
			"ProjectName": ctx.ProjectName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to render custom routes: %w", err)
	}

	project.Files = append(project.Files, &codegen.GeneratedFile{
		Path:     customRoutesPath,
		Content:  customRoutesContent,
		Mode:     0644,
		IsStatic: true,
	})

	return nil
}

// prepareTableData prepares template data for a table
func (g *TypeScriptGenerator) prepareTableData(table validator.Table) map[string]interface{} {
	// Use template engine to convert names
	pascalName, _ := g.templateEngine.RenderString("{{pascalCase (singular .)}}", table.Name)
	camelName, _ := g.templateEngine.RenderString("{{camelCase (singular .)}}", table.Name)
	kebabName, _ := g.templateEngine.RenderString("{{kebabCase .}}", table.Name)
	pluralName, _ := g.templateEngine.RenderString("{{plural .}}", table.Name)

	// Prepare columns data
	columns := make([]map[string]interface{}, 0, len(table.Columns))
	// Prepare fields data (non-primary key columns for models)
	fields := make([]map[string]interface{}, 0)

	for _, col := range table.Columns {
		colData := map[string]interface{}{
			"Name":       col.Name,
			"Type":       col.Type,
			"PrimaryKey": col.PrimaryKey,
			"Unique":     col.Unique,
			"Nullable":   col.Nullable,
			"Default":    col.Default,
			"ForeignKey": col.ForeignKey,
		}

		// Add case variations for column name
		colPascalName, _ := g.templateEngine.RenderString("{{pascalCase .}}", col.Name)
		colCamelName, _ := g.templateEngine.RenderString("{{camelCase .}}", col.Name)
		colData["PascalName"] = colPascalName
		colData["CamelName"] = colCamelName

		// Add type mappings for TypeScript and Zod
		colData["TypeScriptType"] = mapToTypeScriptType(col.Type)
		colData["ZodType"] = mapToZodType(col.Type)

		columns = append(columns, colData)

		// Add to fields if not primary key and not auto-generated timestamp
		if !col.PrimaryKey && col.Name != "created_at" && col.Name != "updated_at" {
			fields = append(fields, colData)
		}
	}

	// Prepare hook placeholders
	hooksImport := "// import { beforeCreate, afterCreate, beforeUpdate, afterUpdate, beforeDelete, afterDelete } from '../hooks/" + kebabName + "';"
	beforeCreateHook := "// beforeCreate hook"
	afterCreateHook := "// afterCreate(created);"
	endBeforeCreateHook := ""
	beforeReadHook := "// beforeRead hook"
	afterReadHook := "// afterRead(item);"
	endBeforeReadHook := ""
	beforeUpdateHook := "// beforeUpdate hook"
	afterUpdateHook := "// afterUpdate(updated);"
	endBeforeUpdateHook := ""
	beforeDeleteHook := "// beforeDelete hook"
	afterDeleteHook := "// afterDelete hook"
	endBeforeDeleteHook := ""

	// Hook placeholders for model/schema generation
	beforeSchemaHook := "// beforeSchema hook"
	afterSchemaHook := "// afterSchema hook"
	beforeTypeHook := "// beforeType hook"
	afterTypeHook := "// afterType hook"

	return map[string]interface{}{
		"TableName":             table.Name,
		"PascalName":            pascalName,
		"CamelName":             camelName,
		"KebabName":             kebabName,
		"PluralName":            pluralName,
		"Columns":               columns,
		"Fields":                fields,
		"HooksImport":           hooksImport,
		"BeforeCreateHook":      beforeCreateHook,
		"AfterCreateHook":       afterCreateHook,
		"EndBeforeCreateHook":   endBeforeCreateHook,
		"BeforeReadHook":        beforeReadHook,
		"AfterReadHook":         afterReadHook,
		"EndBeforeReadHook":     endBeforeReadHook,
		"BeforeUpdateHook":      beforeUpdateHook,
		"AfterUpdateHook":       afterUpdateHook,
		"EndBeforeUpdateHook":   endBeforeUpdateHook,
		"BeforeDeleteHook":      beforeDeleteHook,
		"AfterDeleteHook":       afterDeleteHook,
		"EndBeforeDeleteHook":   endBeforeDeleteHook,
		"BeforeSchemaHook":      beforeSchemaHook,
		"AfterSchemaHook":       afterSchemaHook,
		"BeforeTypeHook":        beforeTypeHook,
		"AfterTypeHook":         afterTypeHook,
	}
}

// mapToTypeScriptType maps Aurora types to TypeScript types
func mapToTypeScriptType(auroraType string) string {
	switch auroraType {
	case "string", "text", "uuid":
		return "string"
	case "integer", "bigint", "float", "decimal":
		return "number"
	case "boolean":
		return "boolean"
	case "timestamp", "date":
		return "Date"
	case "json":
		return "any"
	default:
		return "any"
	}
}

// mapToZodType maps Aurora types to Zod validators
func mapToZodType(auroraType string) string {
	switch auroraType {
	case "string", "text", "uuid":
		return "string"
	case "integer", "bigint":
		return "number().int"
	case "float", "decimal":
		return "number"
	case "boolean":
		return "boolean"
	case "timestamp", "date":
		return "date"
	case "json":
		return "any"
	default:
		return "string"
	}
}

// getTemplateDir returns the full path to a template directory
func (g *TypeScriptGenerator) getTemplateDir(subdir string) string {
	// Get the path relative to this package
	_, file, _, _ := runtime.Caller(0)
	packageDir := filepath.Dir(file)

	// Navigate to project root and then to templates
	// generator.go is in internal/codegen/typescript/
	// templates are in templates/typescript/{subdir}/
	projectRoot := filepath.Join(packageDir, "..", "..", "..")
	return filepath.Join(projectRoot, "templates", "typescript", subdir)
}

// Name returns the name of this generator
func (g *TypeScriptGenerator) Name() string {
	return "typescript"
}
