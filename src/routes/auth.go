package routes

import (
	"gocker-api/auth"
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/utils"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type UserBody struct {
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UserAuthenticateBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	ID         uint   `json:"id"`
	TokenValue string `json:"token"`
}

func InitAuthRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/auth/register", utils.ParseToHandlerFunc(registerUser)).Methods("POST")
	router.HandleFunc("/api/v1/auth/authenticate", utils.ParseToHandlerFunc(authenticateUser)).Methods("POST")
}

func CreateResponseToken(token models.Token) TokenResponse {
	return TokenResponse{ID: token.ID, TokenValue: token.TokenValue}
}

func registerUser(res http.ResponseWriter, req *http.Request) error {
	var userBody UserBody

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

	database := database.GetInstance().GetDB()
	user := models.User{
		FirstName: userBody.FirstName,
		Email:     userBody.Email,
		Password:  userBody.Password,
	}

	if envErr := godotenv.Load(); envErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: envErr.Error()})
	}

	if user.Email == os.Getenv("ADMIN_EMAIL") {
		user.SetRole("ADMIN")
	} else {
		user.SetRole("USER")
	}
	user.EncodePassword(user.Password)
	database.Create(&user)

	tokenString, tokenErr := auth.GenerateToken(user)

	if tokenErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: tokenErr.Error()})
	}

	token := models.Token{
		TokenValue: tokenString,
		UserRefer:  user.ID,
	}
	database.Create(&token)
	return utils.WriteJSON(res, 201, CreateResponseToken(token))

}

func authenticateUser(res http.ResponseWriter, req *http.Request) error {
	var userAuth UserAuthenticateBody
	var user models.User
	var token models.Token

	if parseErr := utils.ReadJSON(req.Body, &userAuth); parseErr != nil {
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

	if result := database.First(&user, "email LIKE ?", userAuth.Email); result.RowsAffected == 0 {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "user not found"})
	} else if err := user.ComparePassword(userAuth.Password); err != nil {
		return utils.WriteJSON(res, 400, utils.ApiError{Error: "wrong password. Please, try again"})
	}

	database.Find(&token, "user_refer = ?", user.ID)
	return utils.WriteJSON(res, 200, CreateResponseToken(token))
}
