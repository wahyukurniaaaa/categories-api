package handlers

import (
	"category-api/models"
	"category-api/services"
	"encoding/json"
	"net/http"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/checkout" {
		switch r.Method {
		case http.MethodPost:
			h.checkout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func (h *TransactionHandler) checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.Checkout(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}