package handlers

import (
	"gocker-api/auth"
	"gocker-api/database"
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

type AuthenticationResponse struct {
	TokenValue        string `json:"token"`
	RefreshTokenValue string `json:"refresh-token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
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

	user, err := services.CreateUser(userBody)

	if err != nil {
		return utils.WriteJSON(res, 500, err.Error())
	}

	tokenString, tokenErr := auth.GenerateToken(user, models.Bearer)

	if tokenErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: tokenErr.Error()})
	}

	token := models.Token{
		TokenValue: tokenString,
		UserRefer:  user.ID,
		Kind:       models.Bearer,
	}

	services.CreateToken(&token)

	refreshTokenString, refreshTokenErr := auth.GenerateToken(user, models.Refresh)

	if refreshTokenErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: refreshTokenErr.Error()})
	}

	refreshToken := models.Token{
		TokenValue: refreshTokenString,
		UserRefer:  user.ID,
		Kind:       models.Refresh,
	}

	services.CreateToken(&refreshToken)

	return utils.WriteJSON(res, 201, AuthenticationResponse{TokenValue: token.TokenValue, RefreshTokenValue: refreshToken.TokenValue})
}

// Function that returns a user's JWT token, given its email and password
func handleAuthenticateUser(res http.ResponseWriter, req *http.Request) error {
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
	token, _ := services.GetTokensByUserReferAndKind(user.ID, models.Bearer)
	refreshToken, _ := services.GetTokensByUserReferAndKind(user.ID, models.Refresh)

	return utils.WriteJSON(res, 200, AuthenticationResponse{TokenValue: token.TokenValue, RefreshTokenValue: refreshToken.TokenValue})
}

func handleRefreshToken(res http.ResponseWriter, req *http.Request) error {
	var refreshTokenRequest RefreshTokenRequest

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

	//Get the refresh-token's user
	claims, claimsErr := auth.GetClaims(refreshTokenRequest.RefreshToken)

	if claimsErr != nil {
		return utils.WriteJSON(res, 500, utils.ApiError{Error: claimsErr.Error()})
	}
	user, notFoundErr := services.GetUserByEmail(claims["email"].(string))

	if notFoundErr != nil {
		utils.WriteJSON(res, 404, utils.ApiError{Error: notFoundErr.Error()})
	}

	//Get the user's bearer token and refresh it
	var token models.Token
	database := database.GetInstance().GetDB()

	database.Find(&token, "user_refer = ? AND kind = ?", user.ID, models.Bearer)
	newTokenString, tokenErr := auth.GenerateToken(user, models.Bearer)

	if tokenErr != nil {
		return utils.WriteJSON(res, 500, tokenErr)
	}

	token.TokenValue = newTokenString
	database.Save(&token)

	return utils.WriteJSON(res, 201, TokenResponse{TokenValue: token.TokenValue})
}
