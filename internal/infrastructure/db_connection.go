package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"nilus-challenge-backend/internal/config"
	infr_config "nilus-challenge-backend/internal/infrastructure/config"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func NewDBConnection() *sql.DB {
	portAsString := os.Getenv("DB_PORT")
	fmt.Println(portAsString)
	if portAsString == "" {
		portAsString = "5432"
	}
	port, _ := strconv.Atoi(portAsString)

	newDatabase := config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := config.NewPostgresConnection(newDatabase)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	infr_config.CreateNotificationsTable(db)

	return db
}