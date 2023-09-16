package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func MakeHTTPHandleFunc(function apiFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		if err := function(response, request); err != nil {
			WriteJSON(response, 500, ApiError{Error: err.Error()})
		}

	}
}

func WriteJSON(response http.ResponseWriter, status int, value any) error {
	response.WriteHeader(status)
	response.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(response).Encode(value)
}

func ReadJSON(body io.Reader, any interface{}) []ApiError {
	bodyValidator := validator.New()
	unmarshalErr := json.NewDecoder(body).Decode(&any)
	errors := make([]ApiError, 0)

	if unmarshalErr != nil {
		errors = append(errors, ApiError{Error: unmarshalErr.Error()})
	}

	if validationErr := bodyValidator.Struct(any); validationErr != nil {
		for _, validationError := range validationErr.(validator.ValidationErrors) {
			errors = append(errors, ApiError{Error: "field " + validationError.Field() + " must be provided"})
		}
	}

	return errors
}
