package handlers

import (
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/services"
	"gocker-api/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ResponseUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

func CreateResponseUser(user models.User) ResponseUser {
	return ResponseUser{ID: user.ID, FirstName: user.FirstName, Email: user.Email}
}

func InitUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users", utils.ParseToHandlerFunc(handleGetUsers)).Methods("GET")
	router.HandleFunc("/api/v1/users", utils.ParseToHandlerFunc(handleCreateUser)).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(handleGetUser)).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(handleUpdateUser)).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(handleDeleteUser)).Methods("DELETE")
}

func handleGetUsers(res http.ResponseWriter, req *http.Request) error {
	users := services.GetAllUsers()
	var responseUsers []ResponseUser = make([]ResponseUser, 0)

	for _, value := range users {
		responseUsers = append(responseUsers, CreateResponseUser(value))
	}

	return utils.WriteJSON(res, 200, responseUsers)
}

func handleGetUser(res http.ResponseWriter, req *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	user, notFoundErr := services.GetUserById(id)

	if notFoundErr != nil {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: notFoundErr.Error()})
	}

	return utils.WriteJSON(res, 200, CreateResponseUser(*user))
}

func handleCreateUser(res http.ResponseWriter, req *http.Request) error {
	var userBody services.UserBody

	// Handle body validation
	if parseErr := utils.ReadJSON(req.Body, &userBody); parseErr != nil {
		if validationErrs, ok := parseErr.(validator.ValidationErrors); ok {
			validationErrors := make([]utils.ApiError, 0)

			for _, validationErr := range validationErrs {
				validationErrors = append(validationErrors, utils.ApiError{Error: "Field " + validationErr.Field() + " must be provided"})
			}

			return utils.WriteJSON(res, 400, validationErrors)
		} else {
			return utils.WriteJSON(res, 400, utils.ApiError{Error: "not valid json."})
		}
	}

	user, err := services.CreateUser(userBody)

	if err != nil {
		return utils.WriteJSON(res, 500, err.Error())
	}

	return utils.WriteJSON(res, 201, CreateResponseUser(*user))
}

func handleUpdateUser(res http.ResponseWriter, req *http.Request) error {
	var updatedUser services.UpdateUserBody
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	if parseErr := utils.ReadJSON(req.Body, &updatedUser); parseErr != nil {
		if errors, ok := parseErr.(validator.ValidationErrors); ok {
			validationErrors := make([]utils.ApiError, 0)

			for _, validationErr := range errors {
				validationErrors = append(validationErrors, utils.ApiError{Error: "Field " + validationErr.Field() + " must be provided"})
			}

			return utils.WriteJSON(res, 400, validationErrors)
		} else {
			return utils.WriteJSON(res, 400, utils.ApiError{Error: "not valid json."})
		}
	}

	user, notFoundErr := services.UpdateUser(id, updatedUser)

	if notFoundErr != nil {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "user not found"})
	}

	return utils.WriteJSON(res, 201, CreateResponseUser(*user))
}

func handleDeleteUser(res http.ResponseWriter, req *http.Request) error {
	var user models.User
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "User not found."})
	} else {
		database.Delete(user)
		return utils.WriteJSON(res, 201, map[string]string{"Success": "User successfully deleted."})
	}
}
