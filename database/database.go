package database

import (
	"database/sql"
	"fmt"
	"log"

	"category-api/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) *sql.DB {
	var dsn string
	if cfg.DBConn != "" {
		dsn = cfg.DBConn
		log.Println("Using DB_CONN from environment variables")
	} else {
		if cfg.DBHost == "" {
			log.Fatal("Database configuration missing: DB_CONN or DB_HOST must be set")
		}
		
		// Default SSL mode override for local development vs production
		sslMode := "disable"
		if cfg.DBHost != "localhost" && cfg.DBHost != "127.0.0.1" {
			sslMode = "require"
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)
		log.Printf("Using structured DB config (Host: %s)", cfg.DBHost)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db
}
