package domain

import (
	"fmt"
	"time"
)

// Project represents a BaaS project
type Project struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	Language     string    `gorm:"type:varchar(50);not null" json:"language"`
	DatabaseType string    `gorm:"type:varchar(50);not null;column:database_type" json:"database_type"`
	APIStyle     string    `gorm:"type:varchar(50);not null;column:api_style" json:"api_style"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Project) TableName() string {
	return "projects"
}

// Validate checks if the project has valid field values
func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	validLanguages := map[string]bool{
		"typescript": true,
		"python":     true,
	}
	if !validLanguages[p.Language] {
		return fmt.Errorf("language must be one of: typescript, python")
	}

	validDatabases := map[string]bool{
		"postgres": true,
		"mongodb":  true,
	}
	if !validDatabases[p.DatabaseType] {
		return fmt.Errorf("database_type must be one of: postgres, mongodb")
	}

	validAPIStyles := map[string]bool{
		"rest":    true,
		"graphql": true,
	}
	if !validAPIStyles[p.APIStyle] {
		return fmt.Errorf("api_style must be one of: rest, graphql")
	}

	return nil
}
