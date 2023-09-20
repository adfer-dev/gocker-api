package routes

import (
	"gocker-api/auth"
	"gocker-api/models"
	"gocker-api/services"
	"gocker-api/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

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

// Function that creates a new user and returns its JWT token
func registerUser(res http.ResponseWriter, req *http.Request) error {
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

	tokenString, tokenErr := auth.GenerateToken(user)

	if tokenErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: tokenErr.Error()})
	}

	token := models.Token{
		TokenValue: tokenString,
		UserRefer:  user.ID,
	}

	services.CreateToken(&token)

	return utils.WriteJSON(res, 201, CreateResponseToken(token))
}

// Function that returns a user's JWT token, given its email and password
func authenticateUser(res http.ResponseWriter, req *http.Request) error {
	var userAuth UserAuthenticateBody

	//Validate user auth body
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

	//Checking if user exists and if password matches
	user, notFoundErr := services.GetUserByEmail(userAuth.Email)

	if notFoundErr != nil {
		return utils.WriteJSON(res, 404, utils.ApiError{Error: "user not found"})
	} else if wrongPasswordErr := user.ComparePassword(userAuth.Password); wrongPasswordErr != nil {
		return utils.WriteJSON(res, 400, utils.ApiError{Error: wrongPasswordErr.Error()})
	}

	//Not checking err since we have previously checked that user exists
	token, _ := services.GetTokenByUserRefer(user.ID)

	return utils.WriteJSON(res, 200, CreateResponseToken(token))
}
