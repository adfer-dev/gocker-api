package config

import (
	"gocker-api/routes"
	"gocker-api/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	port string
}

func NewAPIServer(port string) *APIServer {
	return &APIServer{
		port,
	}
}

func (server *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/users", utils.MakeHTTPHandleFunc(server.HandleUserRoutes))
	router.HandleFunc("/api/v1/users/{id}", utils.MakeHTTPHandleFunc(server.HandleUserParamRoutes))

	log.Println("API server running on port: ", server.GetPort())

	http.ListenAndServe(server.GetPort(), router)
}

func (server *APIServer) HandleUserRoutes(response http.ResponseWriter, request *http.Request) error {
	if request.Method == "GET" {
		return routes.GetAllUsers(response, request)
	}
	if request.Method == "POST" {
		return routes.CreateUser(response, request)
	}

	return utils.WriteJSON(response, 403, utils.ApiError{Error: "Method not supported."})
}

func (server *APIServer) HandleUserParamRoutes(response http.ResponseWriter, request *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(request)["id"])

	if err != nil {
		return utils.WriteJSON(response, 400, utils.ApiError{Error: "id must be a number."})
	}

	if request.Method == "GET" {
		return routes.GetSingleUser(response, request, id)
	}
	if request.Method == "PUT" {
		return routes.UpdateUser(response, request, id)
	}
	if request.Method == "DELETE" {
		return routes.DeleteUser(response, request, id)
	}

	return utils.WriteJSON(response, 403, utils.ApiError{Error: "Method not supported."})
}

func (server APIServer) GetPort() string {
	return server.port
}
