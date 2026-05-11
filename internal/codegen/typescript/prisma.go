package typescript

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ktarun.reddy/baas/internal/codegen/engine"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// GeneratePrismaSchema generates a Prisma schema from Aurora schema definition
func GeneratePrismaSchema(schema *validator.SchemaDefinition) (string, error) {
	// Create template engine once for all operations
	eng := engine.NewTemplateEngine()

	data := struct {
		Provider string
		Models   []PrismaModel
	}{
		Provider: mapDatabaseProvider(schema.Database),
		Models:   convertToPrismaModels(schema.Tables, eng),
	}

	// Load template from file
	templatePath := getTemplateFilePath("schema.prisma.tmpl")
	result, err := eng.RenderFile(templatePath, data)
	if err != nil {
		return "", fmt.Errorf("failed to render prisma schema template: %w", err)
	}

	return result, nil
}

// getTemplateFilePath returns the full path to a template file
func getTemplateFilePath(templateName string) string {
	// Get the path relative to this package
	_, file, _, _ := runtime.Caller(0)
	packageDir := filepath.Dir(file)

	// Navigate to project root and then to templates
	// prisma.go is in internal/codegen/typescript/
	// templates are in templates/typescript/prisma/
	projectRoot := filepath.Join(packageDir, "..", "..", "..")
	return filepath.Join(projectRoot, "templates", "typescript", "prisma", templateName)
}

// PrismaModel represents a Prisma model
type PrismaModel struct {
	Name      string
	Fields    []PrismaField
	Relations []PrismaField
}

// PrismaField represents a field in a Prisma model
type PrismaField struct {
	Name       string
	Type       string
	Modifier   string // ?, []
	Attributes string // @id, @unique, @default(), etc.
}

// convertToPrismaModels converts Aurora tables to Prisma models
func convertToPrismaModels(tables []validator.Table, eng *engine.TemplateEngine) []PrismaModel {
	models := make([]PrismaModel, 0, len(tables))

	for _, table := range tables {
		model := PrismaModel{
			Name:      toPascalCase(table.Name, eng),
			Fields:    []PrismaField{},
			Relations: []PrismaField{},
		}

		// Convert columns to fields
		for _, col := range table.Columns {
			field := PrismaField{
				Name: toCamelCase(col.Name, eng),
				Type: mapPrismaType(col.Type),
			}

			// Build attributes
			var attrs []string

			if col.PrimaryKey {
				attrs = append(attrs, "@id")
				if col.Type == "uuid" {
					attrs = append(attrs, "@default(uuid())")
				} else if col.Type == "integer" || col.Type == "bigint" {
					attrs = append(attrs, "@default(autoincrement())")
				}
			}

			if col.Unique {
				attrs = append(attrs, "@unique")
			}

			if col.Default != "" && !col.PrimaryKey {
				if col.Type == "timestamp" || col.Type == "date" {
					attrs = append(attrs, "@default(now())")
				} else if col.Type == "string" || col.Type == "text" {
					attrs = append(attrs, fmt.Sprintf("@default(\"%s\")", col.Default))
				} else {
					attrs = append(attrs, fmt.Sprintf("@default(%s)", col.Default))
				}
			}

			if !col.Nullable && !col.PrimaryKey {
				// Prisma fields are required by default
			} else if col.Nullable {
				field.Modifier = "?"
			}

			field.Attributes = strings.Join(attrs, " ")

			// Check if this is a foreign key
			if col.ForeignKey != "" {
				// This will be handled as a relation
				field.Attributes += fmt.Sprintf(" @relation(fields: [%s], references: [id])", field.Name)
				model.Fields = append(model.Fields, field)

				// Add the relation field
				parts := strings.Split(col.ForeignKey, ".")
				if len(parts) == 2 {
					relatedTable := parts[0]
					relationField := PrismaField{
						Name: toCamelCase(singularize(relatedTable, eng), eng),
						Type: toPascalCase(singularize(relatedTable, eng), eng),
					}
					model.Relations = append(model.Relations, relationField)
				}
			} else {
				model.Fields = append(model.Fields, field)
			}
		}

		models = append(models, model)
	}

	return models
}

// mapDatabaseProvider maps Aurora database type to Prisma provider
func mapDatabaseProvider(dbType string) string {
	switch dbType {
	case "postgres":
		return "postgresql"
	case "mongodb":
		return "mongodb"
	case "mysql":
		return "mysql"
	case "sqlite":
		return "sqlite"
	default:
		return "postgresql"
	}
}

// mapPrismaType maps Aurora column type to Prisma type
func mapPrismaType(auroraType string) string {
	switch auroraType {
	case "string":
		return "String"
	case "text":
		return "String"
	case "integer":
		return "Int"
	case "bigint":
		return "BigInt"
	case "float":
		return "Float"
	case "decimal":
		return "Decimal"
	case "boolean":
		return "Boolean"
	case "timestamp":
		return "DateTime"
	case "date":
		return "DateTime"
	case "uuid":
		return "String"
	case "json":
		return "Json"
	default:
		return "String"
	}
}

// toPascalCase converts string to PascalCase
func toPascalCase(s string, eng *engine.TemplateEngine) string {
	result, err := eng.RenderString("{{pascalCase .}}", s)
	if err != nil {
		// Log error and return original string
		fmt.Printf("error rendering pascalCase template: %v\n", err)
		return s
	}
	return result
}

// toCamelCase converts string to camelCase
func toCamelCase(s string, eng *engine.TemplateEngine) string {
	result, err := eng.RenderString("{{camelCase .}}", s)
	if err != nil {
		// Log error and return original string
		fmt.Printf("error rendering camelCase template: %v\n", err)
		return s
	}
	return result
}

// singularize converts plural to singular (simple implementation)
func singularize(s string, eng *engine.TemplateEngine) string {
	result, err := eng.RenderString("{{singular .}}", s)
	if err != nil {
		// Log error and return original string
		fmt.Printf("error rendering singular template: %v\n", err)
		return s
	}
	return result
}
