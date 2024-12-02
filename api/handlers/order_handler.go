package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/models"
	"github.com/nurovic/hmall/store"
)

// Sipariş oluşturma
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.CreateOrder(order); err != nil {
		http.Error(w, fmt.Sprintf("Sipariş oluşturulamadı: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// Siparişi ID ile alma
func GetOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	order, err := store.GetOrderByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Sipariş bulunamadı: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}
