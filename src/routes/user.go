package routes

import (
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ResponseUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

type UserBody struct {
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
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
	router.HandleFunc("/api/v1/users", utils.ParseToHandlerFunc(getAllUsers)).Methods("GET")
	router.HandleFunc("/api/v1/users", utils.ParseToHandlerFunc(createUser)).Methods("POST")
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(getSingleUser)).Methods("GET")
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(updateUser)).Methods("PUT")
	router.HandleFunc("/api/v1/users/{id}", utils.ParseToHandlerFunc(deleteUser)).Methods("DELETE")
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

func getSingleUser(res http.ResponseWriter, req *http.Request) error {
	var user models.User
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "user not found."})
	} else {
		return utils.WriteJSON(res, 200, CreateResponseUser(user))
	}
}

func createUser(res http.ResponseWriter, req *http.Request) error {
	var userBody UserBody

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

	if envErr := godotenv.Load(); envErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: envErr.Error()})
	}

	var userRole models.UserRole
	database := database.GetInstance().GetDB()

	// Set registered user's role
	if userBody.Email == os.Getenv("ADMIN_EMAIL") {
		userRole = models.Admin
	} else {
		userRole = models.Standard
	}

	user := models.User{
		FirstName: userBody.FirstName,
		Email:     userBody.Email,
		Password:  nil,
		Role:      userRole,
	}

	user.EncodePassword(userBody.Password)
	database.Create(&user)

	return utils.WriteJSON(res, 201, CreateResponseUser(user))
}

func updateUser(res http.ResponseWriter, req *http.Request) error {
	var user models.User
	var updatedUser UpdateUserBody
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "User not found."})
	}

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

func deleteUser(res http.ResponseWriter, req *http.Request) error {
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
