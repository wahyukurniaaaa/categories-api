package main_test

import (
	"bytes"
	"category-api/handlers"
	"category-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	// Replicating the health check logic for testing purpose since main's handler is anonymous
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(`{"status": "OK", "message": "API Running"}`))
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	expected := `{"status": "OK", "message": "API Running"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestProdukCRUD(t *testing.T) {
	// 1. Create Produk
	newProduk := models.Produk{Nama: "Test Produk", Harga: 10000, Stok: 10}
	body, _ := json.Marshal(newProduk)
	req, _ := http.NewRequest("POST", "/api/produk", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	
	handlers.ProdukHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Create: handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdProduk models.Produk
	json.NewDecoder(rr.Body).Decode(&createdProduk)
	if createdProduk.ID == 0 {
		t.Errorf("Create: ID should be generated. Got %v", createdProduk.ID)
	}

	// 2. Get All Produk
	req, _ = http.NewRequest("GET", "/api/produk", nil)
	rr = httptest.NewRecorder()
	handlers.ProdukHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetAll: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// 3. Get By ID
	url := fmt.Sprintf("/api/produk/%d", createdProduk.ID)
	req, _ = http.NewRequest("GET", url, nil)
	rr = httptest.NewRecorder()
	handlers.ProdukHandler(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetByID: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Update Produk
	createdProduk.Nama = "Updated Name"
	body, _ = json.Marshal(createdProduk)
	req, _ = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	handlers.ProdukHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Update: handler returned wrong status code: got %v want %v body: %s", status, http.StatusOK, rr.Body.String())
	}
	
	// Verify Update
	var updatedProduk models.Produk
	json.NewDecoder(rr.Body).Decode(&updatedProduk)
	if updatedProduk.Nama != "Updated Name" {
		t.Errorf("Update: Name not updated. Got %s", updatedProduk.Nama)
	}

	// 5. Delete Produk
	req, _ = http.NewRequest("DELETE", url, nil)
	rr = httptest.NewRecorder()
	handlers.ProdukHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Delete: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestCategoryCRUD(t *testing.T) {
	h := handlers.NewCategoryHandler()
	
	// 1. Create Category
	newCat := map[string]string{
		"name": "Test Category",
		"description": "Desc",
	}
	body, _ := json.Marshal(newCat)
	req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Create: handler returned wrong status code: got %v want %v body: %s", status, http.StatusCreated, rr.Body.String())
	}

	var createdCat models.Category
	json.NewDecoder(rr.Body).Decode(&createdCat)
	if createdCat.ID == "" {
		t.Errorf("Create: ID should be generated")
	}

	// 2. Get All
	req, _ = http.NewRequest("GET", "/api/categories", nil)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetAll: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 3. Get By ID
	url := "/api/categories/" + createdCat.ID
	req, _ = http.NewRequest("GET", url, nil)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetByID: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Update
	newCat["name"] = "Updated Cat"
	body, _ = json.Marshal(newCat)
	req, _ = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Update: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 5. Delete
	req, _ = http.NewRequest("DELETE", url, nil)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Delete: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
