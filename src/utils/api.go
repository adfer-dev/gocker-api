package utils

import (
	"gocker-api/auth"
	"gocker-api/database"
	"gocker-api/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Error string
}

func Auth(c *gin.Context) {
	fullToken := c.GetHeader("Authorization")

	if c.GetHeader("Authorization") == "" || !strings.HasPrefix(c.GetHeader("Authorization"), "Bearer") {
		c.AbortWithStatusJSON(403, ApiError{Error: "Authorization header must be provided, starting with Bearer"})
	} else {
		tokenString := fullToken[7:]

		if err := auth.ValidateToken(tokenString); err != nil {
			c.AbortWithStatusJSON(403, ApiError{Error: err.Error()})
		} else {
			claims, claimsErr := auth.GetClaims(tokenString)

			if claimsErr != nil {
				c.AbortWithStatusJSON(403, ApiError{Error: claimsErr.Error()})
			} else {
				var user models.User
				database := database.GetInstance().GetDB()

				database.Find(&user, "email LIKE ?", claims["email"])

				if (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE") && user.GetRole() != "ADMIN" {
					c.AbortWithStatusJSON(403, ApiError{Error: "method not allowed."})
				}
			}
		}
	}
}

// Function to validate a request's body.
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

// Middleware to check if the id parameter of an endpoint is a valid number.
func ValidateIdParam(c *gin.Context) {

	if _, err := strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatusJSON(400, ApiError{Error: "Id parameter must be a number."})
	}
}
