package routers

import (
	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/api/handlers"
)

func RegisterOrderRoutes(router *mux.Router) {
	// Siparişle ilgili rotalar
	orderRouter := router.PathPrefix("/orders").Subrouter()
	orderRouter.HandleFunc("/", handlers.CreateOrder).Methods("POST")
	orderRouter.HandleFunc("/{id}", handlers.GetOrder).Methods("GET")

}
