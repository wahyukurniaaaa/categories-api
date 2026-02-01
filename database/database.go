package database

import (
	"database/sql"
	"log"

	"category-api/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) *sql.DB {
	if cfg.DBConn == "" {
		log.Fatal("DB_CONN environment variable is required")
	}

	log.Println("Connecting to database...")

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db
}
