package handlers

import (
	"category-api/models"
	"category-api/services"
	"encoding/json"
	"net/http"
	"strings"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
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
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request, id string) {
	category, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var input models.Category

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, `{"error": "Name is required"}`, http.StatusBadRequest)
		return
	}

	newCategory, err := h.service.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request, id string) {
	var input models.Category

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, `{"error": "Name is required"}`, http.StatusBadRequest)
		return
	}

	updatedCategory, err := h.service.Update(id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCategory)
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Category deleted successfully"}`))
}
