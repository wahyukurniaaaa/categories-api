package handlers

import (
	"category-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CategoryFile = "categories.json"

type CategoryHandler struct {
	categories map[string]models.Category
	mutex      sync.RWMutex
	filename   string
}

func NewCategoryHandler() *CategoryHandler {
	h := &CategoryHandler{
		categories: make(map[string]models.Category),
		filename:   CategoryFile,
	}
	h.loadFromFile()
	return h
}

func (h *CategoryHandler) loadFromFile() {
	file, err := os.ReadFile(h.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, start with empty map
			return
		}
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// If file is empty, do nothing
	if len(file) == 0 {
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()
	if err := json.Unmarshal(file, &h.categories); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
	}
}

func (h *CategoryHandler) saveToFile() error {
	data, err := json.MarshalIndent(h.categories, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.filename, data, 0644)
}

// GET /
func (h *CategoryHandler) ShowHome(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var categoriesHTML string
	if len(h.categories) == 0 {
		categoriesHTML = "<p><i>No categories available. Use POST /categories to add data.</i></p>"
	} else {
		categoriesHTML = "<ul style='padding-left: 20px;'>"
		for _, cat := range h.categories {
			categoriesHTML += fmt.Sprintf("<li style='margin-bottom: 8px;'><strong>%s</strong> (ID: <code style='background:#eee;padding:2px 4px;border-radius:3px;'>%s</code>)<br><span style='color:#666;'>%s</span></li>", cat.Name, cat.ID, cat.Description)
		}
		categoriesHTML += "</ul>"
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Category API Docs</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 20px; background-color: #f4f4f9; }
        h1 { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; display: flex; align-items: center; gap: 10px; }
        h2 { color: #34495e; margin-top: 30px; border-left: 4px solid #3498db; padding-left: 10px; }
        .endpoint { background: #fff; padding: 15px; border-radius: 8px; margin-bottom: 10px; box-shadow: 0 2px 5px rgba(0,0,0,0.05); transition: transform 0.2s; border: 1px solid #eee; }
        .endpoint:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
        .method { display: inline-block; padding: 4px 8px; border-radius: 4px; color: #fff; font-weight: bold; font-size: 0.85em; width: 60px; text-align: center; text-transform: uppercase; }
        .get { background-color: #3498db; }
        .post { background-color: #2ecc71; }
        .put { background-color: #f39c12; }
        .delete { background-color: #e74c3c; }
        .url { font-family: 'Consolas', 'Monaco', monospace; font-weight: bold; margin-left: 10px; color: #333; font-size: 1.05em; }
        .desc { margin-left: 10px; color: #7f8c8d; font-style: italic; }
        .data-section { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.05); border: 1px solid #eee; }
        code { font-family: 'Consolas', 'Monaco', monospace; background: #e8f4f8; padding: 2px 5px; border-radius: 3px; color: #c0392b; }
    </style>
</head>
<body>
    <h1>📂 Category API Documentation</h1>
    <p>Welcome to the <strong>Category API</strong>. Below are the available endpoints to manage your data.</p>
    
    <div class="endpoint">
        <span class="method get">GET</span> <span class="url">/categories</span> <span class="desc">Retrieve all categories</span>
    </div>
    <div class="endpoint">
        <span class="method post">POST</span> <span class="url">/categories</span> <span class="desc">Create a new category (JSON Body: <code>name</code>, <code>description</code>)</span>
    </div>
    <div class="endpoint">
        <span class="method get">GET</span> <span class="url">/categories/:id</span> <span class="desc">Get details of a specific category</span>
    </div>
    <div class="endpoint">
        <span class="method put">PUT</span> <span class="url">/categories/:id</span> <span class="desc">Update a category (JSON Body: <code>name</code>, <code>description</code>)</span>
    </div>
    <div class="endpoint">
        <span class="method delete">DELETE</span> <span class="url">/categories/:id</span> <span class="desc">Remove a category</span>
    </div>

    <h2>📊 Live Data: Categories</h2>
    <div class="data-section">
        %s
    </div>
    
    <footer style="margin-top: 40px; text-align: center; color: #aaa; font-size: 0.9em;">
        Category API v1.0 &bull; Running with Gin Framework
    </footer>
</body>
</html>
`, categoriesHTML)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GET /categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	categories := make([]models.Category, 0, len(h.categories))
	for _, cat := range h.categories {
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// GET /categories/:id
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if category, exists := h.categories[id]; exists {
		c.JSON(http.StatusOK, category)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
	}
}

// POST /categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCategory := models.Category{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
	}

	h.mutex.Lock()
	h.categories[newCategory.ID] = newCategory
	// Save to file
	if err := h.saveToFile(); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}
	h.mutex.Unlock()

	c.JSON(http.StatusCreated, newCategory)
}

// PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.categories[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	updatedCategory := models.Category{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
	}
	h.categories[id] = updatedCategory

	// Save to file
	if err := h.saveToFile(); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.categories[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	delete(h.categories, id)

	// Save to file
	if err := h.saveToFile(); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
