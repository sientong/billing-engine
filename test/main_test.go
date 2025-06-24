package test

import (
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	// Setup: connect + migrate
	SetupTestDB()

	// Run tests
	code := m.Run()

	// Teardown: drop schema
	TeardownTestDB()

	// Exit with correct code
	os.Exit(code)
}
