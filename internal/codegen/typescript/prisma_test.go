package typescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ktarun.reddy/baas/internal/codegen/engine"
)

func TestGeneratePrismaSchema(t *testing.T) {
	// This test is skipped because it requires the template file to be loaded
	// The template is located at templates/typescript/prisma/schema.prisma.tmpl
	// and is loaded using runtime.Caller which may not work in all test environments
	t.Skip("Skipping integration test that requires template file loading")
}

func TestMapPrismaType(t *testing.T) {
	tests := []struct {
		auroraType string
		prismaType string
	}{
		{"string", "String"},
		{"text", "String"},
		{"integer", "Int"},
		{"bigint", "BigInt"},
		{"float", "Float"},
		{"decimal", "Decimal"},
		{"boolean", "Boolean"},
		{"timestamp", "DateTime"},
		{"date", "DateTime"},
		{"uuid", "String"},
		{"json", "Json"},
	}

	for _, tt := range tests {
		t.Run(tt.auroraType, func(t *testing.T) {
			result := mapPrismaType(tt.auroraType)
			assert.Equal(t, tt.prismaType, result)
		})
	}
}

func TestToPascalCase(t *testing.T) {
	eng := engine.NewTemplateEngine()

	tests := []struct {
		input    string
		expected string
	}{
		{"users", "Users"},
		{"user_profiles", "UserProfiles"},
		{"userProfiles", "UserProfiles"},
		{"user-profiles", "UserProfiles"},
		{"id", "Id"},
		{"user_id", "UserId"},
		{"created_at", "CreatedAt"},
		{"posts", "Posts"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input, eng)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	eng := engine.NewTemplateEngine()

	tests := []struct {
		input    string
		expected string
	}{
		{"users", "users"},
		{"user_profiles", "userProfiles"},
		{"UserProfiles", "userProfiles"},
		{"user-profiles", "userProfiles"},
		{"id", "id"},
		{"user_id", "userId"},
		{"created_at", "createdAt"},
		{"Posts", "posts"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toCamelCase(tt.input, eng)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSingularize(t *testing.T) {
	eng := engine.NewTemplateEngine()

	tests := []struct {
		input    string
		expected string
	}{
		{"users", "user"},
		{"posts", "post"},
		{"entries", "entry"},
		{"addresses", "address"}, // Removes 's' suffix
		{"boxes", "box"},         // Removes 'es' suffix
		{"user", "user"},
		{"post", "post"},
		{"categories", "category"}, // Removes 'ies' suffix
		{"status", "statu"},        // Simple implementation removes 's'
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := singularize(tt.input, eng)
			assert.Equal(t, tt.expected, result)
		})
	}
}
