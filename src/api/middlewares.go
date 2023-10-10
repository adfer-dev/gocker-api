package api

import (
	"errors"
	"gocker-api/auth"
	"gocker-api/models"
	"gocker-api/services"
	"gocker-api/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

// Middleware function to check if the auth token provided is correct and has not expired.
func AuthMiddleware(next http.Handler) http.Handler {

	allowedEndpoints := regexp.MustCompile(`/api/v1/auth/*`)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//If the endpoint is not allowed, check its auth token.
		if allowedEndpoints.MatchString(req.URL.Path) {
			next.ServeHTTP(res, req)
		} else {
			authErr := checkAuth(req)

			//If the token is valid, execute the next function. Otherwise, respond with an error.
			if authErr == nil {
				next.ServeHTTP(res, req)
			} else {
				utils.WriteJSON(res, 403, utils.ApiError{Error: authErr.Error()})
			}
		}
	})
}

// Middleware to check if the id parameter of an endpoint is a valid number.
func ValidateIdParam(next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		idParam := mux.Vars(req)["id"]

		//If there is not param, just execute the next function
		if idParam != "" {
			//If there is param check if it's a number.
			if _, err := strconv.Atoi(idParam); err != nil {
				utils.WriteJSON(res, 400, utils.ApiError{Error: "Id parameter must be a number."})
			} else {
				next.ServeHTTP(res, req)
			}
		} else {
			next.ServeHTTP(res, req)
		}

	})
}

// AUX FUNCTIONS
// Function that checks if a request is authorized
func checkAuth(req *http.Request) error {
	fullToken := req.Header.Get("Authorization")

	if fullToken == "" || !strings.HasPrefix(fullToken, "Bearer") {
		return errors.New("authorization token must be provided, starting with Bearer")
	}

	tokenString := fullToken[7:]

	//Validate token
	if err := auth.ValidateToken(tokenString); err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return errors.New("token expired. Please, get a new one at /auth/refresh-token")
		} else {
			return errors.New("token not valid")
		}
	}

	//Then check if token is in the database
	if _, tokenNotFoundErr := services.GetTokenByValue(tokenString); tokenNotFoundErr != nil {
		return errors.New("token revoked")
	}

	claims, claimsErr := auth.GetClaims(tokenString)

	if claimsErr != nil {
		return claimsErr
	}

	user, _ := services.GetUserByEmail(claims["email"].(string))

	if (req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE") && user.Role != models.Admin {
		return errors.New("method not allowed")
	}

	return nil
}
