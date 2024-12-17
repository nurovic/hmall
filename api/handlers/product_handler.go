package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/models"
	"github.com/nurovic/hmall/store"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.CreateProduct(ctx, product); err != nil {
		http.Error(w, fmt.Sprintf("Ürün oluşturulamadı: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()

	id := mux.Vars(r)["id"]

	product, err := store.GetProductByID(ctx, id) 
	if err != nil {
		http.Error(w, fmt.Sprintf("Ürün bulunamadı: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}
