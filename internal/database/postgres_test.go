package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ktarun.reddy/baas/internal/config"
)

func TestConnect(t *testing.T) {
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
	require.NotNil(t, db)

	// Verify connection is alive
	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)
}
