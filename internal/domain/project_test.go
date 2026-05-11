package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid project with all fields",
			project: Project{
				ID:           "550e8400-e29b-41d4-a716-446655440000",
				Name:         "Test Project",
				Description:  "A test project",
				Language:     "typescript",
				DatabaseType: "postgres",
				APIStyle:     "rest",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			project: Project{
				ID:           "550e8400-e29b-41d4-a716-446655440000",
				Name:         "",
				Language:     "typescript",
				DatabaseType: "postgres",
				APIStyle:     "rest",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid language",
			project: Project{
				ID:           "550e8400-e29b-41d4-a716-446655440000",
				Name:         "Test Project",
				Language:     "ruby",
				DatabaseType: "postgres",
				APIStyle:     "rest",
			},
			wantErr: true,
			errMsg:  "language must be one of: typescript, python",
		},
		{
			name: "invalid database type",
			project: Project{
				ID:           "550e8400-e29b-41d4-a716-446655440000",
				Name:         "Test Project",
				Language:     "typescript",
				DatabaseType: "oracle",
				APIStyle:     "rest",
			},
			wantErr: true,
			errMsg:  "database_type must be one of: postgres, mongodb",
		},
		{
			name: "invalid API style",
			project: Project{
				ID:           "550e8400-e29b-41d4-a716-446655440000",
				Name:         "Test Project",
				Language:     "typescript",
				DatabaseType: "postgres",
				APIStyle:     "soap",
			},
			wantErr: true,
			errMsg:  "api_style must be one of: rest, graphql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProject_TableName(t *testing.T) {
	p := Project{}
	assert.Equal(t, "projects", p.TableName())
}
