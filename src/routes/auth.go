package routes

import (
	"gocker-api/auth"
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type UserBody struct {
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UserAuthenticateBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	ID         uint   `json:"id"`
	TokenValue string `json:"token"`
}

func InitAuthRoutes(app *gin.Engine) {
	app.POST("/api/v1/auth/register", registerUser)
	app.POST("/api/v1/auth/authenticate", authenticateUser)
}

func CreateResponseToken(token models.Token) TokenResponse {
	return TokenResponse{ID: token.ID, TokenValue: token.TokenValue}
}

func registerUser(c *gin.Context) {
	var userBody UserBody

	if parseErr := c.BindJSON(&userBody); parseErr != nil {
		c.JSON(400, utils.ApiError{Error: parseErr.Error()})
	} else if validationErrors := utils.ValidateBody(userBody); len(validationErrors) != 0 {
		c.JSON(400, validationErrors)
	} else {
		database := database.GetInstance().GetDB()
		user := models.User{
			FirstName: userBody.FirstName,
			Email:     userBody.Email,
			Password:  userBody.Password,
		}

		if envErr := godotenv.Load(); envErr != nil {
			c.JSON(500, utils.ApiError{Error: envErr.Error()})
		}

		if user.Email == os.Getenv("ADMIN_EMAIL") {
			user.SetRole("ADMIN")
		} else {
			user.SetRole("USER")
		}
		user.EncodePassword(user.Password)
		database.Create(&user)

		token, tokenErr := auth.GenerateToken(user)

		if tokenErr != nil {
			c.JSON(500, utils.ApiError{Error: tokenErr.Error()})
		} else {
			token := models.Token{
				TokenValue: token,
				UserRefer:  user.ID,
			}
			database.Create(&token)
			c.JSON(201, CreateResponseToken(token))
		}
	}
}

func authenticateUser(c *gin.Context) {
	var userAuth UserAuthenticateBody
	var user models.User
	var token models.Token

	if parseErr := c.BindJSON(&userAuth); parseErr != nil {
		c.JSON(400, utils.ApiError{Error: parseErr.Error()})
	} else if validationErrors := utils.ValidateBody(userAuth); len(validationErrors) != 0 {
		c.JSON(400, validationErrors)
	} else {
		database := database.GetInstance().GetDB()
		if result := database.First(&user, "email LIKE ?", userAuth.Email); result.RowsAffected == 0 {
			c.JSON(404, utils.ApiError{Error: "user not found"})
		} else if err := user.ComparePassword(userAuth.Password); err != nil {
			c.JSON(400, utils.ApiError{Error: "wrong password. Please, try again"})
		} else {
			database.Find(&token, "user_refer = ?", user.ID)
			c.JSON(200, CreateResponseToken(token))
		}
	}
}
