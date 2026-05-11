package domain

import (
	"fmt"
	"time"
)

// Schema represents a project's schema definition
type Schema struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID string    `gorm:"type:uuid;not null;column:project_id;index" json:"project_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Version   int       `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Schema) TableName() string {
	return "schemas"
}

// Validate checks if the schema has valid field values
func (s *Schema) Validate() error {
	if s.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}

	if s.Content == "" {
		return fmt.Errorf("content is required")
	}

	if s.Version < 1 {
		return fmt.Errorf("version must be >= 1")
	}

	return nil
}
