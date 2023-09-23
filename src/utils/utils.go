package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
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

// MIDDLEWARES

//UTILITY FUNCTIONS

func WriteJSON(res http.ResponseWriter, status int, value any) error {

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)

	return json.NewEncoder(res).Encode(value)
}

func ReadJSON(reader io.Reader, body interface{}) error {

	if deserializeErr := json.NewDecoder(reader).Decode(body); deserializeErr != nil {
		return deserializeErr
	}

	if validationErr := validateBody(body); validationErr != nil {
		return validationErr
	}

	return nil
}

//Auxiliary functions

// Function to validate a request's body.
func validateBody(body interface{}) error {
	newValidator := validator.New()

	if err := newValidator.Struct(body); err != nil {
		return err
	}

	return nil
}
