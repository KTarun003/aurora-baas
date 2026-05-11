package database

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/internal/config"
)

func TestRunMigrations(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "aurora",
		Password: "aurora_dev",
		Name:     "aurora_dev",
		SSLMode:  "disable",
	}

	db, err := Connect(cfg)
	require.NoError(t, err)

	err = RunMigrations(db)
	require.NoError(t, err)

	// Verify tables exist
	sqlDB, _ := db.DB()
	var tableCount int
	err = sqlDB.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_name IN ('projects', 'schemas')
	`).Scan(&tableCount)
	require.NoError(t, err)
	require.Equal(t, 2, tableCount)
}
