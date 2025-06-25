package app

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func SetupDB() *sql.DB {

	db, err := sql.Open("postgres", "user=postgres dbname=billing-engine sslmode=disable password=099052")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
