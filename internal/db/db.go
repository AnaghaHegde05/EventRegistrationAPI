package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := fmt.Sprintf(
		"host=localhost port=5432 user=postgres password=%s dbname=eventdb sslmode=disable",
		os.Getenv("DB_PASSWORD"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}