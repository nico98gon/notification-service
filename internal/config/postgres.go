package config

import (
	"database/sql"
	"fmt"
)

func NewPostgresConnection(dbConfig DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Ping successful")
	return db, nil
}
