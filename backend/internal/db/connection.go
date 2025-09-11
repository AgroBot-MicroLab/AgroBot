package db

import (
    "os"
    "log"
	"database/sql"
)

func NewDBConnection() (*sql.DB) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

    return db
}
