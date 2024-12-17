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

func CreateOrder(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.CreateOrder(ctx, order); err != nil {
		http.Error(w, fmt.Sprintf("Sipariş oluşturulamadı: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()
	id := mux.Vars(r)["id"]

	order, err := store.GetOrderByID(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Sipariş bulunamadı: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}
