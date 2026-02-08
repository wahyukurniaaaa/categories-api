package main

import (
	"category-api/config"
	"category-api/database"
	"category-api/handlers"
	"category-api/repositories"
	"category-api/services"
	"database/sql"
	"log"
	"net/http"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect Database
	db := database.InitDB(cfg)
	defer db.Close()
	
	// Create Tables
	createTables(db)

	// 3. Init Layers (Dependency Injection)
	// Repositories
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Services
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo, productRepo)
	reportService := services.NewReportService(transactionRepo)

	// Handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)

	// 4. Setup Routing
	http.HandleFunc("/api/categories", categoryHandler.ServeHTTP)
	http.HandleFunc("/api/categories/", categoryHandler.ServeHTTP)

	http.HandleFunc("/api/produk", productHandler.ServeHTTP)
	http.HandleFunc("/api/produk/", productHandler.ServeHTTP)

	// Transaction (Checkout)
	http.HandleFunc("/api/checkout", transactionHandler.ServeHTTP)

	// Report
	http.HandleFunc("/api/report/hari-ini", reportHandler.ServeHTTP)

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "OK", "message": "API Running with DB"}`))
	})
	
	// Root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`Welcome to Kasir API (Layered Architecture).`))
	})

	// Run Server
	log.Printf("Server starting on port %s...", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			nama VARCHAR(100) NOT NULL,
			harga INT NOT NULL,
			stok INT NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
	log.Println("Database tables initialized")
}