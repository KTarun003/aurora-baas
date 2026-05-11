package testutils

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID creates a new UUID for testing
func GenerateUUID() string {
	return uuid.New().String()
}

// FixedTime returns a fixed time for deterministic tests
func FixedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2026-05-11T10:00:00Z")
	return t
}

// AssertNoError fails the test if error is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
