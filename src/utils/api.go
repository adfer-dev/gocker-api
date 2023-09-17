package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Error string
}

func ValidateBody(body interface{}) []ApiError {
	newValidator := validator.New()
	errors := make([]ApiError, 0)

	if err := newValidator.Struct(body); err != nil {

		for _, validationError := range err.(validator.ValidationErrors) {
			errors = append(errors, ApiError{Error: "Field " + validationError.Field() + " must be provided."})
		}
	}

	return errors
}

// Middleware to check if the id parameter of an endpoint is a number
func ValidateIdParam(c *gin.Context) {

	if _, err := strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(400, ApiError{Error: "Id parameter must be a number."})
	}

	c.Next()
}
