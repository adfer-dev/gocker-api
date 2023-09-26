package services

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
)

// Function that gets a token by its value
func GetTokenByValue(tokenString string) (token models.Token, err error) {
	database := database.GetInstance().GetDB()

	result := database.Find(&token, "token_value LIKE ?", tokenString)

	if result.RowsAffected == 0 {
		err = errors.New("token not found")
	}

	return
}

// Function that saves a token to the database
func CreateToken(token *models.Token) {
	database := database.GetInstance().GetDB()
	database.Create(token)
}

// Function that deletes a token from the database
func DeleteToken(token models.Token) {
	database := database.GetInstance().GetDB()

	database.Delete(token)
}
