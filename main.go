package main

import (
	"category-api/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	categoryHandler := handlers.NewCategoryHandler()

	// 🔗 Endpoint yang Wajib Ada
	// GET / → Halaman Dokumentasi & Data
	r.GET("/", categoryHandler.ShowHome)

	// GET /categories → Ambil semua kategori
	r.GET("/categories", categoryHandler.GetAllCategories)
	// POST /categories → Tambah kategori
	r.POST("/categories", categoryHandler.CreateCategory)
	// GET /categories/{id} → Ambil detail satu kategori
	r.GET("/categories/:id", categoryHandler.GetCategoryByID)
	// PUT /categories/{id} → Update kategori
	r.PUT("/categories/:id", categoryHandler.UpdateCategory)
	// DELETE /categories/{id} → Hapus kategori
	r.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
