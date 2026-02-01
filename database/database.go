package database

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"

	"category-api/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) *sql.DB {
	var dsn string
	if cfg.DBConn != "" {
		// Parse DSN to handle IPv6/IPv4 resolution issues (Render/Supabase specific)
		parsedURL, err := url.Parse(cfg.DBConn)
		if err == nil {
			host := parsedURL.Hostname()
			// Try to resolve multiple IPs and pick IPv4
			ips, err := net.LookupIP(host)
			if err == nil {
				for _, ip := range ips {
					if ip.To4() != nil {
						// Found IPv4, replace hostname with IP
						log.Printf("Resolved %s to IPv4: %s", host, ip.String())
						port := parsedURL.Port()
						if port != "" {
							parsedURL.Host = fmt.Sprintf("%s:%s", ip.String(), port)
						} else {
							parsedURL.Host = ip.String()
						}
						// Force DSN to use this new URL
						dsn = parsedURL.String()
						break
					}
				}
			}
		}

		if dsn == "" { // Fallback if parsing failed or no IPv4 found (unlikely)
			dsn = cfg.DBConn
		}
		log.Println("Using DB_CONN from environment variables")
	} else {
		if cfg.DBHost == "" {
			log.Fatal("Database configuration missing: DB_CONN or DB_HOST must be set")
		}

		// Default SSL mode - require for remote connections
		sslMode := "disable"
		if cfg.DBHost != "localhost" && cfg.DBHost != "127.0.0.1" {
			sslMode = "require"
		}
		
		// Try to resolve IPv4 for the host to avoid IPv6 issues on some cloud platforms
		hostaddr := ""
		if cfg.DBHost != "localhost" && cfg.DBHost != "127.0.0.1" {
			ips, err := net.LookupIP(cfg.DBHost)
			if err == nil {
				for _, ip := range ips {
					if ip.To4() != nil {
						hostaddr = ip.String()
						log.Printf("Resolved %s to IPv4: %s", cfg.DBHost, hostaddr)
						break
					}
				}
			}
		}
		
		// Build DSN - use hostaddr if we found IPv4, keeping host for SSL verification
		if hostaddr != "" {
			dsn = fmt.Sprintf("host=%s hostaddr=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				cfg.DBHost, hostaddr, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)
		} else {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, sslMode)
		}
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
