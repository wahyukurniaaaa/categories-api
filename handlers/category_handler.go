package handlers

import (
	"category-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

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
			return
		}
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

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

// ServeHTTP implements the routing logic for Categories
func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle /api/categories
	if r.URL.Path == "/api/categories" {
		switch r.Method {
		case http.MethodGet:
			h.getAllCategories(w, r)
		case http.MethodPost:
			h.createCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle /api/categories/{id}
	if strings.HasPrefix(r.URL.Path, "/api/categories/") {
		id := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		if id == "" {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getCategoryByID(w, r, id)
		case http.MethodPut:
			h.updateCategory(w, r, id)
		case http.MethodDelete:
			h.deleteCategory(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	
	http.NotFound(w, r)
}

func (h *CategoryHandler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	categories := make([]models.Category, 0, len(h.categories))
	for _, cat := range h.categories {
		categories = append(categories, cat)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request, id string) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if category, exists := h.categories[id]; exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(category)
	} else {
		http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
	}
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	
	if input.Name == "" {
		http.Error(w, `{"error": "Name is required"}`, http.StatusBadRequest)
		return
	}

	newCategory := models.Category{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
	}

	h.mutex.Lock()
	h.categories[newCategory.ID] = newCategory
	if err := h.saveToFile(); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}
	h.mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request, id string) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	
	if input.Name == "" {
		http.Error(w, `{"error": "Name is required"}`, http.StatusBadRequest)
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.categories[id]; !exists {
		http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
		return
	}

	updatedCategory := models.Category{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
	}
	h.categories[id] = updatedCategory

	if err := h.saveToFile(); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCategory)
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request, id string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.categories[id]; !exists {
		http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
		return
	}

	delete(h.categories, id)

	if err := h.saveToFile(); err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Category deleted successfully"}`))
}

