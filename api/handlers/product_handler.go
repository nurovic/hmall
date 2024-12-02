package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/models"
	"github.com/nurovic/hmall/store"
)

// Ürün oluşturma
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.CreateProduct(product); err != nil {
		http.Error(w, fmt.Sprintf("Ürün oluşturulamadı: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// Ürünü ID ile alma
func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	product, err := store.GetProductByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ürün bulunamadı: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}
