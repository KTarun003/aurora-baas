package service

import (
	"testing"

	"github.com/ktarun.reddy/baas/internal/codegen"
	"github.com/ktarun.reddy/baas/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockProjectRepository for testing
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) FindByID(id string) (*domain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) List() ([]*domain.Project, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockSchemaRepository for testing
type MockSchemaRepository struct {
	mock.Mock
}

func (m *MockSchemaRepository) Create(schema *domain.Schema) error {
	args := m.Called(schema)
	return args.Error(0)
}

func (m *MockSchemaRepository) FindLatestByProjectID(projectID string) (*domain.Schema, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Schema), args.Error(1)
}

// MockGenerator for testing
type MockGenerator struct {
	mock.Mock
}

func (m *MockGenerator) Generate(ctx *codegen.GeneratorContext) (*codegen.GeneratedProject, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*codegen.GeneratedProject), args.Error(1)
}

func (m *MockGenerator) Name() string {
	args := m.Called()
	return args.String(0)
}

func TestNewCodegenService(t *testing.T) {
	registry := codegen.NewRegistry()
	projectRepo := &MockProjectRepository{}
	schemaRepo := &MockSchemaRepository{}

	service := NewCodegenService(registry, projectRepo, schemaRepo)
	assert.NotNil(t, service)
}

func TestNewCodegenService_NilRegistry(t *testing.T) {
	projectRepo := &MockProjectRepository{}
	schemaRepo := &MockSchemaRepository{}

	assert.Panics(t, func() {
		NewCodegenService(nil, projectRepo, schemaRepo)
	})
}

func TestNewCodegenService_NilProjectRepo(t *testing.T) {
	registry := codegen.NewRegistry()
	schemaRepo := &MockSchemaRepository{}

	assert.Panics(t, func() {
		NewCodegenService(registry, nil, schemaRepo)
	})
}

func TestNewCodegenService_NilSchemaRepo(t *testing.T) {
	registry := codegen.NewRegistry()
	projectRepo := &MockProjectRepository{}

	assert.Panics(t, func() {
		NewCodegenService(registry, projectRepo, nil)
	})
}

func TestCodegenService_GenerateProject(t *testing.T) {
	// Setup
	registry := codegen.NewRegistry()
	projectRepo := &MockProjectRepository{}
	schemaRepo := &MockSchemaRepository{}

	// Create mock generator
	mockGen := &MockGenerator{}
	mockGen.On("Name").Return("typescript")
	mockGen.On("Generate", mock.Anything).Return(&codegen.GeneratedProject{
		ProjectID: "test-project",
		Language:  "typescript",
		Files:     []*codegen.GeneratedFile{},
	}, nil)

	registry.Register(mockGen)

	service := NewCodegenService(registry, projectRepo, schemaRepo)

	// Setup mock expectations
	project := &domain.Project{
		ID:           "test-project",
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	projectRepo.On("FindByID", "test-project").Return(project, nil)

	schemaContent := `
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
`
	schema := &domain.Schema{
		ProjectID: "test-project",
		Version:   1,
		Content:   schemaContent,
	}
	schemaRepo.On("FindLatestByProjectID", "test-project").Return(schema, nil)

	// Execute
	generatedProject, err := service.GenerateProject("test-project", "/tmp/output")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, generatedProject)
	assert.Equal(t, "test-project", generatedProject.ProjectID)
	assert.Equal(t, "typescript", generatedProject.Language)

	projectRepo.AssertExpectations(t)
	schemaRepo.AssertExpectations(t)
	mockGen.AssertExpectations(t)
}

func TestCodegenService_GenerateProject_ProjectNotFound(t *testing.T) {
	// Setup
	registry := codegen.NewRegistry()
	projectRepo := &MockProjectRepository{}
	schemaRepo := &MockSchemaRepository{}

	service := NewCodegenService(registry, projectRepo, schemaRepo)

	// Setup mock expectations
	projectRepo.On("FindByID", "nonexistent").Return(nil, assert.AnError)

	// Execute
	_, err := service.GenerateProject("nonexistent", "/tmp/output")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get project")

	projectRepo.AssertExpectations(t)
}

func TestCodegenService_GenerateProject_SchemaNotFound(t *testing.T) {
	// Setup
	registry := codegen.NewRegistry()
	projectRepo := &MockProjectRepository{}
	schemaRepo := &MockSchemaRepository{}

	service := NewCodegenService(registry, projectRepo, schemaRepo)

	// Setup mock expectations
	project := &domain.Project{
		ID:           "test-project",
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	projectRepo.On("FindByID", "test-project").Return(project, nil)
	schemaRepo.On("FindLatestByProjectID", "test-project").Return(nil, assert.AnError)

	// Execute
	_, err := service.GenerateProject("test-project", "/tmp/output")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get schema")

	projectRepo.AssertExpectations(t)
	schemaRepo.AssertExpectations(t)
}

func TestCodegenService_GenerateProject_InvalidSchema(t *testing.T) {
	// Setup
	registry := codegen.NewRegistry()
	projectRepo := &MockProjectRepository{}
	schemaRepo := &MockSchemaRepository{}

	service := NewCodegenService(registry, projectRepo, schemaRepo)

	// Setup mock expectations
	project := &domain.Project{
		ID:           "test-project",
		Name:         "Test Project",
		Language:     "typescript",
		DatabaseType: "postgres",
		APIStyle:     "rest",
	}
	projectRepo.On("FindByID", "test-project").Return(project, nil)

	schema := &domain.Schema{
		ProjectID: "test-project",
		Version:   1,
		Content:   "invalid yaml content {{{",
	}
	schemaRepo.On("FindLatestByProjectID", "test-project").Return(schema, nil)

	// Execute
	_, err := service.GenerateProject("test-project", "/tmp/output")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse schema")

	projectRepo.AssertExpectations(t)
	schemaRepo.AssertExpectations(t)
}
