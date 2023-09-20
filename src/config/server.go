package config

import (
	"gocker-api/routes"
	"gocker-api/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
}

func (server *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(utils.AuthMiddleware)
	router.Use(utils.ValidateIdParam)
	initRoutes(router)

	return http.ListenAndServe(server.ListenAddress, router)
}

func initRoutes(router *mux.Router) {
	routes.InitUserRoutes(router)
	routes.InitAuthRoutes(router)
}
