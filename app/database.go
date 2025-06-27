package app

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func SetupDB() *sql.DB {

	conn := os.Getenv("DATABASE_URL")
	if conn == "" {
		conn = "postgres://postgres:099052@localhost:5432/billing-engine?sslmode=disable"
	}

	log.Println(conn)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	log.Println("Database connection established")

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

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

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Database migrated successfully")

	return db
}
