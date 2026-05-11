package validator

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Column represents a database column definition
type Column struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	PrimaryKey bool   `yaml:"primary_key"`
	Unique     bool   `yaml:"unique"`
	Nullable   bool   `yaml:"nullable"`
	Default    string `yaml:"default"`
	ForeignKey string `yaml:"foreign_key"`
}

// Index represents a database index definition
type Index struct {
	Columns []string `yaml:"columns"`
	Unique  bool     `yaml:"unique"`
}

// Table represents a database table definition
type Table struct {
	Name    string    `yaml:"name"`
	Columns []Column  `yaml:"columns"`
	Indexes []Index   `yaml:"indexes"`
}

// SchemaDefinition represents the overall schema definition
type SchemaDefinition struct {
	Version  string  `yaml:"version"`
	Database string  `yaml:"database"`
	Tables   []Table `yaml:"tables"`
}

// validColumnTypes contains all valid column types
var validColumnTypes = map[string]bool{
	"string":    true,
	"text":      true,
	"integer":   true,
	"bigint":    true,
	"decimal":   true,
	"float":     true,
	"boolean":   true,
	"date":      true,
	"timestamp": true,
	"uuid":      true,
	"json":      true,
}

// validDatabases contains all valid database types
var validDatabases = map[string]bool{
	"postgres": true,
	"mongodb":  true,
}

// ValidateSchema validates a YAML schema content and returns a SchemaDefinition or error
func ValidateSchema(yamlContent string) (*SchemaDefinition, error) {
	// Parse YAML
	var schema SchemaDefinition
	err := yaml.Unmarshal([]byte(yamlContent), &schema)
	if err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	// Validate version is required
	if strings.TrimSpace(schema.Version) == "" {
		return nil, fmt.Errorf("version is required")
	}

	// Validate database is required
	if strings.TrimSpace(schema.Database) == "" {
		return nil, fmt.Errorf("database type is required")
	}

	// Validate database type is valid
	if !validDatabases[schema.Database] {
		return nil, fmt.Errorf("database type must be 'postgres' or 'mongodb', got '%s'", schema.Database)
	}

	// Track seen table names for duplicate detection
	seenTableNames := make(map[string]bool)

	// Validate tables
	for _, table := range schema.Tables {
		// Validate table name is required
		if strings.TrimSpace(table.Name) == "" {
			return nil, fmt.Errorf("table name is required")
		}

		// Validate duplicate table names
		if seenTableNames[table.Name] {
			return nil, fmt.Errorf("duplicate table name: '%s'", table.Name)
		}
		seenTableNames[table.Name] = true

		// Validate table has at least one column
		if len(table.Columns) == 0 {
			return nil, fmt.Errorf("table '%s' must have at least one column", table.Name)
		}

		// Track seen column names for duplicate detection within this table
		seenColumnNames := make(map[string]bool)

		// Validate columns
		for _, column := range table.Columns {
			// Validate column name is required
			if strings.TrimSpace(column.Name) == "" {
				return nil, fmt.Errorf("column name is required in table '%s'", table.Name)
			}

			// Validate duplicate column names within table
			if seenColumnNames[column.Name] {
				return nil, fmt.Errorf("table %s: duplicate column name: '%s'", table.Name, column.Name)
			}
			seenColumnNames[column.Name] = true

			// Validate column type is required
			if strings.TrimSpace(column.Type) == "" {
				return nil, fmt.Errorf("column type is required for column '%s' in table '%s'", column.Name, table.Name)
			}

			// Validate column type is valid
			if !validColumnTypes[column.Type] {
				return nil, fmt.Errorf("column type '%s' is invalid for column '%s' in table '%s'. valid types are: string, text, integer, bigint, decimal, float, boolean, date, timestamp, uuid, json", column.Type, column.Name, table.Name)
			}
		}
	}

	return &schema, nil
}
