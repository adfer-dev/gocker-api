package services

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/storage"
	"os"
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

var userStorage storage.Storage = &storage.UserStorage{}

func GetAllUsers() []models.User {
	var users []models.User

	database := database.GetInstance().GetDB()
	database.Find(&users)

	return users
}

func GetUserById(id int) (*models.User, error) {

	// don't check the type assertion, since we are sure that the Get method is returning *models.User
	user, err := userStorage.Get(id)

	if err != nil {
		return nil, err
	}

	return user.(*models.User), nil
}

func GetUserByEmail(email string) (user *models.User, err error) {
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "email LIKE ?", email); result.RowsAffected == 0 {
		err = errors.New("user not found")
	}

	return
}

func CreateUser(userBody UserBody) (*models.User, error) {

	// first check that the user email has not already been registered
	if _, notFoundErr := GetUserByEmail(userBody.Email); notFoundErr == nil {
		return nil, errors.New("email already registered")
	}

	var userRole models.UserRole

	// Set user properties
	if userBody.Email == os.Getenv("ADMIN_EMAIL") {
		userRole = models.Admin
	} else {
		userRole = models.Standard
	}

	user := &models.User{
		FirstName: userBody.FirstName,
		Email:     userBody.Email,
		Password:  nil,
		Role:      userRole,
	}

	user.EncodePassword(userBody.Password)
	createErr := userStorage.Create(user)

	return user, createErr
}

func UpdateUser(id int, updatedUser UpdateUserBody) (*models.User, error) {
	var user *models.User
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return nil, errors.New("user not found")
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

	updateErr := userStorage.Update(user)

	return user, updateErr
}

func DeleteUser(id int) (err error) {
	var user *models.User
	database := database.GetInstance().GetDB()

	if result := database.Find(user, "id = ?", id); result.Error != nil {
		err = errors.New("user not found")
		return
	}

	err = userStorage.Delete(user)

	return
}
