package handlers

import (
	"category-api/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	produkData = []models.Produk{}
	produkMu   sync.Mutex
	nextID     = 1
)

func ProdukHandler(w http.ResponseWriter, r *http.Request) {
	// Handle /api/produk
	if r.URL.Path == "/api/produk" {
		switch r.Method {
		case http.MethodGet:
			getAllProduk(w, r)
		case http.MethodPost:
			createProduk(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Handle /api/produk/{id}
	// We expect path to be like /api/produk/123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProdukByID(w, r, id)
	case http.MethodPut:
		updateProduk(w, r, id)
	case http.MethodDelete:
		deleteProduk(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produkData)
}

func createProduk(w http.ResponseWriter, r *http.Request) {
	var p models.Produk
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	produkMu.Lock()
	p.ID = nextID
	nextID++
	produkData = append(produkData, p)
	produkMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getProdukByID(w http.ResponseWriter, r *http.Request, id int) {
	produkMu.Lock()
	defer produkMu.Unlock()

	for _, p := range produkData {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk not found", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request, id int) {
	var updatedData models.Produk
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	produkMu.Lock()
	defer produkMu.Unlock()

	for i, p := range produkData {
		if p.ID == id {
			produkData[i].Nama = updatedData.Nama
			produkData[i].Harga = updatedData.Harga
			produkData[i].Stok = updatedData.Stok
			
			// Return the updated object with ID preserved
			produkData[i].ID = id // Ensure ID hasn't successfully changed if passed in body
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produkData[i])
			return
		}
	}

	http.Error(w, "Produk not found", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request, id int) {
	produkMu.Lock()
	defer produkMu.Unlock()

	for i, p := range produkData {
		if p.ID == id {
			produkData = append(produkData[:i], produkData[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Produk deleted"}`))
			return
		}
	}

	http.Error(w, "Produk not found", http.StatusNotFound)
}
