package routes

import (
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResponseUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

type UpdateUserBody struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func CreateResponseUser(user models.User) ResponseUser {
	return ResponseUser{ID: user.ID, FirstName: user.FirstName, Email: user.Email}
}

func InitUserRoutes(app *gin.Engine) {
	app.GET("/api/v1/users", getAllUsers)
	app.GET("/api/v1/users/:id", utils.ValidateIdParam, getSingleUser)
	app.POST("/api/v1/users", createUser)
	app.PUT("/api/v1/users/:id", utils.ValidateIdParam, updateUser)
	app.DELETE("/api/v1/users/:id", utils.ValidateIdParam, deleteUser)
}

func getAllUsers(c *gin.Context) {
	var users []models.User
	var responseUsers []ResponseUser

	database := database.GetInstance().GetDB()
	database.Find(&users)

	for _, value := range users {
		responseUsers = append(responseUsers, CreateResponseUser(value))
	}

	c.JSON(200, responseUsers)
}

func getSingleUser(c *gin.Context) {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(404, utils.ApiError{Error: "user not found."})
	} else {
		c.JSON(200, CreateResponseUser(user))
	}
}

func createUser(c *gin.Context) {
	var user models.User

	if parseErr := c.BindJSON(&user); parseErr != nil {
		c.JSON(400, utils.ApiError{Error: parseErr.Error()})
	} else if validationErrors := utils.ValidateBody(user); len(validationErrors) != 0 {
		c.JSON(400, validationErrors)
	} else {
		user.EncodePassword(user.Password)
		database := database.GetInstance().GetDB()
		database.Create(&user)

		c.JSON(201, CreateResponseUser(user))
	}
}

func updateUser(c *gin.Context) {
	var user models.User
	var updatedUser UpdateUserBody
	id, _ := strconv.Atoi(c.Param("id"))
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(404, utils.ApiError{Error: "User not found."})
	} else {

		if parseErr := c.BindJSON(&user); parseErr != nil {
			c.JSON(400, utils.ApiError{Error: parseErr.Error()})
		} else if validationErrors := utils.ValidateBody(user); len(validationErrors) > 0 {
			c.JSON(400, validationErrors)
		} else {

			if updatedUser.FirstName != "" {
				user.FirstName = updatedUser.FirstName
			}
			if updatedUser.Email != "" {
				user.Email = updatedUser.Email
			}
			if updatedUser.Password != "" {
				user.EncodePassword(updatedUser.Password)
			}

			database.Save(&user)

			c.JSON(201, CreateResponseUser(user))
		}
	}
}

func deleteUser(c *gin.Context) {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		c.JSON(404, utils.ApiError{Error: "User not found."})
	} else {
		database.Delete(user)
		c.JSON(201, map[string]string{"Success": "User successfully deleted."})
	}
}
