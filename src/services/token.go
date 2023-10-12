package services

import (
	"errors"
	"gocker-api/database"
	"gocker-api/models"
	"gocker-api/storage"
)

var tokenStorage storage.Storage = &storage.TokenStorage{}

func GetTokenById(id int) (*models.Token, error) {
	token, getTokenErr := tokenStorage.Get(id)

	if getTokenErr != nil {
		return nil, getTokenErr
	}

	return token.(*models.Token), nil
}

// Function that gets a token by its value
func GetTokenByValue(tokenString string) (*models.Token, error) {
	var token *models.Token
	database := database.GetInstance().GetDB()
	result := database.Find(&token, "token_value LIKE ?", tokenString)

	if result.RowsAffected == 0 {
		return nil, errors.New("token not found")
	}

	return token, nil
}

// Function that saves a token to the database
func CreateToken(token *models.Token) (*models.Token, error) {
	if createErr := tokenStorage.Create(token); createErr != nil {
		return nil, createErr
	}

	return token, nil
}

// Function that deletes a token from the database
func DeleteToken(token *models.Token) error {
	return tokenStorage.Delete(token)
}
