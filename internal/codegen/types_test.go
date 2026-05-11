package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ktarun.reddy/baas/pkg/validator"
)

func TestGeneratorContext_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *GeneratorContext
		wantErr bool
	}{
		{
			name: "valid context",
			ctx: &GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "TestProject",
				Language:    "typescript",
				APIStyle:    "rest",
				Database:    "postgres",
				Schema: &validator.SchemaDefinition{
					Version:  "1.0",
					Database: "postgres",
					Tables:   []validator.Table{},
				},
			},
			wantErr: false,
		},
		{
			name: "missing project ID",
			ctx: &GeneratorContext{
				ProjectName: "TestProject",
				Language:    "typescript",
			},
			wantErr: true,
		},
		{
			name: "invalid language",
			ctx: &GeneratorContext{
				ProjectID: "test",
				Language:  "ruby",
			},
			wantErr: true,
		},
		{
			name: "invalid API style",
			ctx: &GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "TestProject",
				Language:    "typescript",
				APIStyle:    "soap",
			},
			wantErr: true,
		},
		{
			name: "invalid database",
			ctx: &GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "TestProject",
				Language:    "typescript",
				APIStyle:    "rest",
				Database:    "mysql",
			},
			wantErr: true,
		},
		{
			name: "missing project name",
			ctx: &GeneratorContext{
				ProjectID: "test-project",
				Language:  "typescript",
				APIStyle:  "rest",
				Database:  "postgres",
				Schema: &validator.SchemaDefinition{
					Version:  "1.0",
					Database: "postgres",
					Tables:   []validator.Table{},
				},
			},
			wantErr: true,
		},
		{
			name: "missing schema",
			ctx: &GeneratorContext{
				ProjectID:   "test-project",
				ProjectName: "TestProject",
				Language:    "typescript",
				APIStyle:    "rest",
				Database:    "postgres",
				Schema:      nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
