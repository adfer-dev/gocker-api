package routes

import (
	"gocker-api/database"
	"gocker-api/models"
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

type UpdateUserBody struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func CreateResponseUser(user models.User) ResponseUser {
	return ResponseUser{ID: user.ID, FirstName: user.FirstName, Email: user.Email}
}

func InitUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users", utils.ParseToHandlerFunc(handleUserRoutes))
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(handleUserParamRoutes))
}

func handleUserRoutes(res http.ResponseWriter, req *http.Request) error {

	if req.Method == "GET" {
		return getAllUsers(res, req)
	} else if req.Method == "POST" {
		return createUser(res, req)
	} else {
		return utils.WriteJSON(res, 400, utils.ApiError{Error: "Method not allowed"})
	}
}

func handleUserParamRoutes(res http.ResponseWriter, req *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	if req.Method == "GET" {
		return getSingleUser(res, req, id)
	} else if req.Method == "PUT" {
		return updateUser(res, req, id)
	} else if req.Method == "DELETE" {
		return deleteUser(res, req, id)
	} else {
		return utils.WriteJSON(res, 400, utils.ApiError{Error: "Method not allowed"})
	}
}

func getAllUsers(res http.ResponseWriter, req *http.Request) error {
	var users []models.User
	var responseUsers []ResponseUser

	database := database.GetInstance().GetDB()
	database.Find(&users)

	for _, value := range users {
		responseUsers = append(responseUsers, CreateResponseUser(value))
	}

	return utils.WriteJSON(res, 200, responseUsers)
}

func getSingleUser(res http.ResponseWriter, req *http.Request, id int) error {
	var user models.User

	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "user not found."})
	} else {
		return utils.WriteJSON(res, 200, CreateResponseUser(user))
	}
}

func createUser(res http.ResponseWriter, req *http.Request) error {
	var user models.User

	if parseErr := utils.ReadJSON(req.Body, user); parseErr != nil {
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

	database := database.GetInstance().GetDB()
	database.Create(&user)

	return utils.WriteJSON(res, 201, CreateResponseUser(user))
}

func updateUser(res http.ResponseWriter, req *http.Request, id int) error {
	var user models.User
	var updatedUser UpdateUserBody

	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "User not found."})
	}

	if parseErr := utils.ReadJSON(req.Body, updatedUser); parseErr != nil {
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
	if updatedUser.FirstName != "" {
		user.FirstName = updatedUser.FirstName
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	if updatedUser.Password != "" {
		user.EncodePassword(updatedUser.Password)
	}

	database.Save(&user)

	return utils.WriteJSON(res, 201, CreateResponseUser(user))
}

func deleteUser(res http.ResponseWriter, req *http.Request, id int) error {
	var user models.User

	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "User not found."})
	} else {
		database.Delete(user)
		return utils.WriteJSON(res, 201, map[string]string{"Success": "User successfully deleted."})
	}
}
