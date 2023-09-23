package api

import (
	"gocker-api/auth"
	"gocker-api/utils"
	"net/http"
	"regexp"
	"strconv"

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
			authErr := auth.CheckAuth(res, req)

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
