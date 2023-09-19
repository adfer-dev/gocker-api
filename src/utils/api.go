package utils

import (
	"encoding/json"
	"errors"
	"gocker-api/auth"
	"gocker-api/database"
	"gocker-api/models"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

type ApiError struct {
	Error string
}

type APIFunc func(res http.ResponseWriter, req *http.Request) error

func ParseToHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if err := f(res, req); err != nil {
			WriteJSON(res, 500, err.Error())
		}

	}
}

//MIDDLEWARES

// Middleware function to check if the auth token provided is correct and has not expired.
func AuthMiddleware(next http.Handler) http.Handler {

	allowedEndpoints := regexp.MustCompile(`/api/v1/auth/*`)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//If the endpoint is not allowed, check its auth token.
		if !allowedEndpoints.MatchString(req.URL.Path) {

			authErr := checkAuth(res, req)

			//If the token is valid, execute the next function. Otherwise, respond with an error.
			if authErr != nil {
				WriteJSON(res, 403, ApiError{Error: authErr.Error()})
			} else {
				next.ServeHTTP(res, req)
			}
		} else {
			next.ServeHTTP(res, req)
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
				WriteJSON(res, 400, ApiError{Error: "Id parameter must be a number."})
			} else {
				next.ServeHTTP(res, req)
			}
		} else {
			next.ServeHTTP(res, req)
		}

	})
}

//UTILITY FUNCTIONS

// Function to validate a request's body.
func ValidateBody(body interface{}) error {
	newValidator := validator.New()

	if err := newValidator.Struct(body); err != nil {
		return echo.ErrBadGateway
	}

	return nil
}

func WriteJSON(res http.ResponseWriter, status int, value any) error {

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)

	return json.NewEncoder(res).Encode(value)
}

func ReadJSON(reader io.Reader, body interface{}) error {

	if deserializeErr := json.NewDecoder(reader).Decode(&body); deserializeErr != nil {
		return deserializeErr
	}

	if validationErr := ValidateBody(body); validationErr != nil {
		return validationErr
	}

	return nil
}

//Auxiliary functions

func checkAuth(res http.ResponseWriter, req *http.Request) error {
	fullToken := req.Header.Get("Authorization")

	if fullToken == "" || !strings.HasPrefix(fullToken, "Bearer") {
		return errors.New("authorization header must be provided, starting with Bearer")
	}

	tokenString := fullToken[7:]

	if err := auth.ValidateToken(tokenString); err != nil {
		return err
	}

	claims, claimsErr := auth.GetClaims(tokenString)

	if claimsErr != nil {
		return claimsErr
	}

	var user models.User
	database := database.GetInstance().GetDB()

	database.Find(&user, "email LIKE ?", claims["email"])

	if (req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE") && user.GetRole() != "ADMIN" {
		return errors.New("method not allowed")
	}

	return nil
}
