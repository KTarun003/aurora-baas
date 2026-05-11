package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchema_Validate(t *testing.T) {
	tests := []struct {
		name    string
		schema  Schema
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid schema",
			schema: Schema{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				ProjectID: "550e8400-e29b-41d4-a716-446655440000",
				Content:   "entities:\n  - name: User\n    fields:\n      - name: id\n        type: uuid",
				Version:   1,
			},
			wantErr: false,
		},
		{
			name: "missing project ID",
			schema: Schema{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				ProjectID: "",
				Content:   "entities:\n  - name: User",
				Version:   1,
			},
			wantErr: true,
			errMsg:  "project_id is required",
		},
		{
			name: "missing content",
			schema: Schema{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				ProjectID: "550e8400-e29b-41d4-a716-446655440000",
				Content:   "",
				Version:   1,
			},
			wantErr: true,
			errMsg:  "content is required",
		},
		{
			name: "invalid version zero",
			schema: Schema{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				ProjectID: "550e8400-e29b-41d4-a716-446655440000",
				Content:   "entities:\n  - name: User",
				Version:   0,
			},
			wantErr: true,
			errMsg:  "version must be >= 1",
		},
		{
			name: "invalid version negative",
			schema: Schema{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				ProjectID: "550e8400-e29b-41d4-a716-446655440000",
				Content:   "entities:\n  - name: User",
				Version:   -1,
			},
			wantErr: true,
			errMsg:  "version must be >= 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSchema_TableName(t *testing.T) {
	s := Schema{}
	assert.Equal(t, "schemas", s.TableName())
}
