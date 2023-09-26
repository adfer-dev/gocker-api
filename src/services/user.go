package services

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
	"os"

	"github.com/joho/godotenv"
)

type UserBody struct {
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UpdateUserBody struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func GetAllUsers() []models.User {
	var users []models.User

	database := database.GetInstance().GetDB()
	database.Find(&users)

	return users
}

func GetUserById(id int) (user models.User, err error) {
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		err = errors.New("user not found")
	}

	return
}

func GetUserByEmail(email string) (user models.User, err error) {
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "email LIKE ?", email); result.RowsAffected == 0 {
		err = errors.New("user not found")
	}

	return
}

func CreateUser(userBody UserBody) (user models.User, err error) {

	// first check that the user email has not already been registered
	if _, notFoundErr := GetUserByEmail(userBody.Email); notFoundErr == nil {
		err = errors.New("email already registered")
		return
	}

	if envErr := godotenv.Load(); envErr != nil {
		err = envErr
		return
	}

	var userRole models.UserRole
	database := database.GetInstance().GetDB()

	// Set user properties
	if userBody.Email == os.Getenv("ADMIN_EMAIL") {
		userRole = models.Admin
	} else {
		userRole = models.Standard
	}

	user = models.User{
		FirstName: userBody.FirstName,
		Email:     userBody.Email,
		Password:  nil,
		Role:      userRole,
	}

	user.EncodePassword(userBody.Password)
	database.Create(&user)

	return
}

func UpdateUser(id int, updatedUser UpdateUserBody) (user models.User, err error) {
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		err = errors.New("user not found")
		return
	}

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

	return
}

func DeleteUser(id int) error {
	var user models.User
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	database.Delete(user)

	return nil
}
