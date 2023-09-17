package routes

import (
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/utils"
	"net/http"
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

func GetAllUsers(response http.ResponseWriter, request *http.Request) error {
	var users []models.User
	var responseUsers []ResponseUser

	database := database.GetInstance().GetDB()
	database.Find(&users)

	for _, value := range users {
		responseUsers = append(responseUsers, CreateResponseUser(value))
	}

	return utils.WriteJSON(response, 200, responseUsers)
}

func GetSingleUser(response http.ResponseWriter, request *http.Request, id int) error {
	var user models.User
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(response, 404, utils.ApiError{Error: "user not found."})
	}

	return utils.WriteJSON(response, 200, CreateResponseUser(user))
}

func CreateUser(response http.ResponseWriter, request *http.Request) error {
	var user models.User

	if errors := utils.ReadJSON(request.Body, &user); len(errors) > 0 {
		return utils.WriteJSON(response, 400, errors)
	}

	user.EncodePassword(user.Password)
	database := database.GetInstance().GetDB()
	database.Create(&user)

	return utils.WriteJSON(response, 201, CreateResponseUser(user))
}

func UpdateUser(response http.ResponseWriter, request *http.Request, id int) error {
	var user models.User
	var updatedUser UpdateUserBody
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(response, 404, utils.ApiError{Error: "User not found."})
	}

	if errors := utils.ReadJSON(request.Body, &updatedUser); len(errors) > 0 {
		return utils.WriteJSON(response, 400, errors)
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

	return utils.WriteJSON(response, 201, CreateResponseUser(user))
}

func DeleteUser(response http.ResponseWriter, request *http.Request, id int) error {
	var user models.User
	database := database.GetInstance().GetDB()

	if result := database.Find(&user, "id = ?", id); result.RowsAffected == 0 {
		return utils.WriteJSON(response, 404, utils.ApiError{Error: "User not found."})
	}

	database.Delete(user)

	return utils.WriteJSON(response, 201, map[string]string{"Success": "User successfully deleted."})
}
