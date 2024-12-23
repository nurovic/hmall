package routers

import (
	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/api/handlers"
)

func RegisterUserRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/me", handlers.GetMe).Methods("GET")
	userRouter.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	userRouter.HandleFunc("/", handlers.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", handlers.GetUser).Methods("GET")
}
