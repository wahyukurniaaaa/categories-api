package main

import (
	"category-api/handlers"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Initialize Category Handler
	categoryHandler := handlers.NewCategoryHandler()

	// Handler untuk /api/categories (handles GET, POST) dan /api/categories/ (handles ID operations)
	http.HandleFunc("/api/categories", categoryHandler.ServeHTTP)
	http.HandleFunc("/api/categories/", categoryHandler.ServeHTTP)

	// Handler untuk /api/produk
	http.HandleFunc("/api/produk", handlers.ProdukHandler)
	http.HandleFunc("/api/produk/", handlers.ProdukHandler)

	// Handler untuk Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(`{"status": "OK", "message": "API Running"}`))
	})
	
	// Handler untuk Root / (Documentation link or simple message)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`Welcome to Kasir API. Endpoints: /api/produk, /api/categories`))
	})

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
