package routers

import (
	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/api/handlers"
)

func RegisterProductRoutes(router *mux.Router) {
	productRouter := router.PathPrefix("/products").Subrouter()
	productRouter.HandleFunc("/", handlers.CreateProduct).Methods("POST")
	productRouter.HandleFunc("/{id}", handlers.GetProduct).Methods("GET")

}
