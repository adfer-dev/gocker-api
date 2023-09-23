package api

import (
	"gocker-api/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
}

func (server *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(AuthMiddleware)
	router.Use(ValidateIdParam)
	initRoutes(router)

	return http.ListenAndServe(server.ListenAddress, router)
}

func initRoutes(router *mux.Router) {
	handlers.InitUserRoutes(router)
	handlers.InitAuthRoutes(router)
}
