package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

// MockGenerator for testing
type MockGenerator struct {
	GenerateFunc func(ctx *GeneratorContext) (*GeneratedProject, error)
	NameFunc     func() string
}

func (m *MockGenerator) Generate(ctx *GeneratorContext) (*GeneratedProject, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx)
	}
	return &GeneratedProject{}, nil
}

func (m *MockGenerator) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock"
}

func TestGenerator_Interface(t *testing.T) {
	mock := &MockGenerator{
		GenerateFunc: func(ctx *GeneratorContext) (*GeneratedProject, error) {
			return &GeneratedProject{
				ProjectID: ctx.ProjectID,
				Language:  ctx.Language,
				Files:     []*GeneratedFile{},
			}, nil
		},
		NameFunc: func() string {
			return "mock-generator"
		},
	}

	ctx := &GeneratorContext{
		ProjectID:   "test-project",
		ProjectName: "Test",
		Language:    "typescript",
		APIStyle:    "rest",
		Database:    "postgres",
		Schema: &validator.SchemaDefinition{
			Version:  "1.0",
			Database: "postgres",
			Tables:   []validator.Table{},
		},
	}

	project, err := mock.Generate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "test-project", project.ProjectID)
	assert.Equal(t, "typescript", project.Language)
	assert.Equal(t, "mock-generator", mock.Name())
}

func TestRegistry_NewRegistry(t *testing.T) {
	registry := NewRegistry()
	assert.NotNil(t, registry)
	assert.Equal(t, 0, len(registry.List()))
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	gen1 := &MockGenerator{
		NameFunc: func() string {
			return "typescript"
		},
	}

	gen2 := &MockGenerator{
		NameFunc: func() string {
			return "python"
		},
	}

	registry.Register(gen1)
	registry.Register(gen2)

	assert.Equal(t, 2, len(registry.List()))
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()

	gen := &MockGenerator{
		NameFunc: func() string {
			return "typescript"
		},
	}

	registry.Register(gen)

	retrieved, ok := registry.Get("typescript")
	assert.True(t, ok)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "typescript", retrieved.Name())

	// Test getting non-existent generator
	missing, ok := registry.Get("python")
	assert.False(t, ok)
	assert.Nil(t, missing)
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()

	generators := []string{"typescript", "python", "go"}

	for _, name := range generators {
		gen := &MockGenerator{
			NameFunc: func() string {
				return name
			},
		}
		registry.Register(gen)
	}

	list := registry.List()
	assert.Equal(t, 3, len(list))

	// Verify all generators are in the list
	listMap := make(map[string]bool)
	for _, name := range list {
		listMap[name] = true
	}

	for _, name := range generators {
		assert.True(t, listMap[name], "generator %s should be in list", name)
	}
}

func TestRegistry_Generate(t *testing.T) {
	registry := NewRegistry()

	// Register a mock generator
	gen := &MockGenerator{
		NameFunc: func() string {
			return "test-gen"
		},
		GenerateFunc: func(ctx *GeneratorContext) (*GeneratedProject, error) {
			return &GeneratedProject{
				ProjectID: ctx.ProjectID,
				Language:  ctx.Language,
				Files:     []*GeneratedFile{},
			}, nil
		},
	}
	registry.Register(gen)

	// Generate using the registry
	ctx := &GeneratorContext{
		ProjectID:   "test-project",
		ProjectName: "Test Project",
		Language:    "typescript",
		APIStyle:    "rest",
		Database:    "postgres",
		OutputDir:   "/tmp/test",
		Schema: &validator.SchemaDefinition{
			Version:  "1.0",
			Database: "postgres",
			Tables:   []validator.Table{},
		},
	}

	project, err := registry.Generate("test-gen", ctx)
	require.NoError(t, err, "expected generation to succeed")
	assert.Equal(t, "test-project", project.ProjectID)
}

func TestRegistry_GenerateNonExistent(t *testing.T) {
	registry := NewRegistry()

	ctx := &GeneratorContext{
		ProjectID: "test-project",
		Language:  "typescript",
	}

	_, err := registry.Generate("nonexistent", ctx)
	assert.Error(t, err, "expected error for non-existent generator")
	assert.Contains(t, err.Error(), "generator not found")
}
