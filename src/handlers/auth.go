package handlers

import (
	"gocker-api/models"
	"gocker-api/services"
	"gocker-api/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AuthenticationResponse struct {
	TokenValue        string `json:"token"`
	RefreshTokenValue string `json:"refresh-token"`
}

type TokenResponse struct {
	TokenValue string `json:"token"`
}

func InitAuthRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/auth/register", utils.ParseToHandlerFunc(handleRegisterUser)).Methods("POST")
	router.HandleFunc("/api/v1/auth/authenticate", utils.ParseToHandlerFunc(handleAuthenticateUser)).Methods("POST")
	router.HandleFunc("/api/v1/auth/refresh-token", utils.ParseToHandlerFunc(handleRefreshToken)).Methods("POST")
}

func CreateResponseToken(token models.Token) AuthenticationResponse {
	return AuthenticationResponse{TokenValue: token.TokenValue}
}

// Function that creates a new user and returns its JWT token
func handleRegisterUser(res http.ResponseWriter, req *http.Request) error {
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

	accessToken, refreshToken, err := services.RegisterUser(userBody)

	if err != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(res, 201, AuthenticationResponse{TokenValue: accessToken.TokenValue, RefreshTokenValue: refreshToken.TokenValue})
}

// Function that returns a user's JWT token, given its email and password
func handleAuthenticateUser(res http.ResponseWriter, req *http.Request) error {
	var userAuth services.UserAuthenticateBody

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

	accessToken, refreshToken, err := services.AuthenticateUser(userAuth)

	if err != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(res, 200, AuthenticationResponse{TokenValue: accessToken.TokenValue, RefreshTokenValue: refreshToken.TokenValue})
}

func handleRefreshToken(res http.ResponseWriter, req *http.Request) error {
	var refreshTokenRequest services.RefreshTokenRequest

	//Validate request body
	if parseErr := utils.ReadJSON(req.Body, &refreshTokenRequest); parseErr != nil {
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

	accessToken, err := services.RefreshToken(refreshTokenRequest)

	if err != nil {
		return utils.WriteJSON(res, 400, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(res, 201, TokenResponse{TokenValue: accessToken.TokenValue})
}
