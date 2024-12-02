package api

import (
	"github.com/gorilla/mux"
	"github.com/nurovic/hmall/api/routers"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	routers.RegisterUserRoutes(router)
	routers.RegisterProductRoutes(router)
	routers.RegisterOrderRoutes(router)

	return router
}
