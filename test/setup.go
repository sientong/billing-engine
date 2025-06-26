package test

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var TestDB *sql.DB
var conn = "postgres://postgres:099052@localhost:5432/billing-engine-test?sslmode=disable"

func SetupTestDB() *sql.DB {

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var root string
	if filepath.Base(cwd) == "repository" || filepath.Base(cwd) == "service" || filepath.Base(cwd) == "integration" {
		root = filepath.Join(cwd, "../../")
	} else {
		root = cwd
	}

	migrationsPath := "file://" + filepath.Join(root, "db", "migrations")

	m, err := migrate.New(
		migrationsPath,
		conn,
	)

	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	TruncateAllTables(db)

	return db
}

func TeardownTestDB() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var root string
	if filepath.Base(cwd) == "repository" {
		root = filepath.Join(cwd, "../../")
	} else {
		root = cwd
	}

	migrationsPath := "file://" + filepath.Join(root, "db", "migrations")

	m, err := migrate.New(
		migrationsPath,
		conn,
	)

	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	if err := m.Drop(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	_ = TestDB.Close()
}

func TruncateUsersTable(db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
}

func TruncateAllTables(db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE payments, billing_schedules, loans, users RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
}
