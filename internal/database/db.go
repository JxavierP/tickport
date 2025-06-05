package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(path string) (*sql.DB, error)  {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open tickport database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to connect to tickport database: %w", err)
	}

	fmt.Println("Connected to tickport database successfully:", path)
	return db, nil
}

