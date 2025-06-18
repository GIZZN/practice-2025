package db

import (
	"database/sql"
	"fmt"
	"os"
	_"github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		if dbURL[len(dbURL)-1] != '?' {
			dbURL += "?"
		}
		dbURL += "sslmode=disable"

		var err error
		DB, err = sql.Open("postgres", dbURL)
		if err != nil {
			return fmt.Errorf("error opening database: %v", err)
		}
	} else {
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)

		var err error
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			return fmt.Errorf("error opening database: %v", err)
		}
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	return nil
}
