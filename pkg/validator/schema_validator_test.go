package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateSchema_ValidYAML tests validation of a valid schema with users table
func TestValidateSchema_ValidYAML(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
        primary_key: true
      - name: "email"
        type: "string"
        unique: true
      - name: "created_at"
        type: "timestamp"
`
	schema, err := ValidateSchema(yamlContent)
	require.NoError(t, err)
	require.NotNil(t, schema)
	assert.Equal(t, "1.0", schema.Version)
	assert.Equal(t, "postgres", schema.Database)
	assert.Len(t, schema.Tables, 1)
	assert.Equal(t, "users", schema.Tables[0].Name)
	assert.Len(t, schema.Tables[0].Columns, 3)
}

// TestValidateSchema_InvalidYAML tests that malformed YAML is rejected
func TestValidateSchema_InvalidYAML(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
        primary_key: true
      - name: "email"
        type: "string"
        unique: true
      - name: "created_at"
        type: "timestamp"
this is not valid yaml {{{
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid YAML")
}

// TestValidateSchema_MissingVersion tests that missing version field is caught
func TestValidateSchema_MissingVersion(t *testing.T) {
	yamlContent := `
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_MissingDatabase tests that missing database field is caught
func TestValidateSchema_MissingDatabase(t *testing.T) {
	yamlContent := `
version: "1.0"
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_InvalidDatabase tests that invalid database type is rejected
func TestValidateSchema_InvalidDatabase(t *testing.T) {
	yamlContent := `
version: "1.0"
database: oracle
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database type")
}

// TestValidateSchema_EmptyTableName tests that empty table name is rejected
func TestValidateSchema_EmptyTableName(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: ""
    columns:
      - name: "id"
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "table name")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_EmptyColumnName tests that empty column name is rejected
func TestValidateSchema_EmptyColumnName(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: ""
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column name")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_InvalidColumnType tests that invalid column type is rejected
func TestValidateSchema_InvalidColumnType(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "invalid_type"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column type")
	assert.Contains(t, err.Error(), "invalid")
}

// TestValidateSchema_MissingColumnType tests that empty column type is rejected
func TestValidateSchema_MissingColumnType(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: ""
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column type")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_TableWithNoColumns tests that tables with no columns are rejected
func TestValidateSchema_TableWithNoColumns(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns: []
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have at least one column")
}

// TestValidateSchema_DuplicateTableNames tests that duplicate table names are rejected
func TestValidateSchema_DuplicateTableNames(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate table name")
	assert.Contains(t, err.Error(), "users")
}

// TestValidateSchema_DuplicateColumnNames tests that duplicate column names within a table are rejected
func TestValidateSchema_DuplicateColumnNames(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "uuid"
      - name: "id"
        type: "string"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate column name")
	assert.Contains(t, err.Error(), "id")
}

// TestValidateSchema_WhitespaceOnlyTableName tests that whitespace-only table name is rejected
func TestValidateSchema_WhitespaceOnlyTableName(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "   "
    columns:
      - name: "id"
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "table name")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_WhitespaceOnlyColumnName tests that whitespace-only column name is rejected
func TestValidateSchema_WhitespaceOnlyColumnName(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "   "
        type: "uuid"
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column name")
	assert.Contains(t, err.Error(), "required")
}

// TestValidateSchema_WhitespaceOnlyColumnType tests that whitespace-only column type is rejected
func TestValidateSchema_WhitespaceOnlyColumnType(t *testing.T) {
	yamlContent := `
version: "1.0"
database: postgres
tables:
  - name: "users"
    columns:
      - name: "id"
        type: "   "
`
	_, err := ValidateSchema(yamlContent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "column type")
	assert.Contains(t, err.Error(), "required")
}
